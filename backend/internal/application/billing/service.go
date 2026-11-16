package billing

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/billing"
	"clinic-api/internal/domain/shared"
	"clinic-api/internal/adapters/payment"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrInvoiceNotFound      = errors.New("invoice not found")
	ErrInvoiceAlreadyPaid   = errors.New("invoice is already fully paid")
	ErrInvalidPaymentAmount = errors.New("payment amount must be positive")
	ErrOverpayment          = errors.New("payment amount exceeds remaining unpaid balance")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	billingRepo       billing.Repository
	accountingService *accountingApp.Service
	cashProcessor     payment.PaymentProcessor
	bankProcessor     payment.PaymentProcessor
}

func NewService(
	uow UnitOfWork,
	billingRepo billing.Repository,
	accountingService *accountingApp.Service,
	cashProcessor payment.PaymentProcessor,
	bankProcessor payment.PaymentProcessor,
) *Service {
	return &Service{
		uow:               uow,
		billingRepo:       billingRepo,
		accountingService: accountingService,
		cashProcessor:     cashProcessor,
		bankProcessor:     bankProcessor,
	}
}

type CreateItemCommand struct {
	ItemType    billing.ItemType `json:"item_type"`
	ItemID      *uuid.UUID       `json:"item_id,omitempty"`
	Description string           `json:"description"`
	Quantity    int              `json:"quantity"`
	UnitPrice   decimal.Decimal  `json:"unit_price"`
}

type CreateInvoiceCommand struct {
	BranchID        uuid.UUID           `json:"branch_id"`
	PatientID       *uuid.UUID          `json:"patient_id,omitempty"`
	PatientName     string              `json:"patient_name"`
	PatientNIK      string              `json:"patient_nik"`
	QueueID         *int64              `json:"queue_id,omitempty"`
	InvoiceType     billing.InvoiceType `json:"invoice_type"`
	RegistrationFee decimal.Decimal     `json:"registration_fee"`
	ConsultationFee decimal.Decimal     `json:"consultation_fee"`
	MedicationFee   decimal.Decimal     `json:"medication_fee"`
	Items           []CreateItemCommand `json:"items"`
}

func (s *Service) CreateInvoice(ctx context.Context, cmd CreateInvoiceCommand) (*billing.Invoice, error) {
	total := cmd.RegistrationFee.Add(cmd.ConsultationFee).Add(cmd.MedicationFee)

	var items []billing.InvoiceItem
	if len(cmd.Items) > 0 {
		itemsTotal := decimal.Zero
		for _, it := range cmd.Items {
			qty := it.Quantity
			if qty <= 0 {
				qty = 1
			}
			lineTotal := it.UnitPrice.Mul(decimal.FromInt(int64(qty)))
			itemsTotal = itemsTotal.Add(lineTotal)

			items = append(items, billing.InvoiceItem{
				ItemType:    it.ItemType,
				ItemID:      it.ItemID,
				Description: it.Description,
				Quantity:    qty,
				UnitPrice:   it.UnitPrice,
				TotalPrice:  lineTotal,
			})
		}
		if total.IsZero() {
			total = itemsTotal
		}
	}

	invoiceNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	inv := &billing.Invoice{
		BaseEntity:      shared.BaseEntity{ID: uuid.New()},
		BranchScoped:    shared.BranchScoped{BranchID: cmd.BranchID},
		InvoiceNumber:   invoiceNumber,
		PatientID:       cmd.PatientID,
		PatientName:     cmd.PatientName,
		PatientNIK:      cmd.PatientNIK,
		QueueID:         cmd.QueueID,
		InvoiceType:     cmd.InvoiceType,
		RegistrationFee: cmd.RegistrationFee,
		ConsultationFee: cmd.ConsultationFee,
		MedicationFee:   cmd.MedicationFee,
		TotalAmount:     total,
		PaidAmount:      decimal.Zero,
		PaymentStatus:   billing.PaymentStatusUnpaid,
		Items:           items,
	}

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.billingRepo.Create(txCtx, inv)
	})
	if err != nil {
		return nil, fmt.Errorf("create invoice in uow: %w", err)
	}

	return inv, nil
}

type PayInvoiceCommand struct {
	InvoiceID         uuid.UUID       `json:"invoice_id"`
	Amount            decimal.Decimal `json:"amount"`
	PaymentMethod     string          `json:"payment_method"` // CASH, BANK_TRANSFER, RECEIVABLE, QRIS, EDC
	UnitCode          string          `json:"unit_code,omitempty"` // 'KIDINEST' or 'SSC'
	ReferenceNo       string          `json:"reference_no,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	ReceivedBy        *uuid.UUID      `json:"received_by,omitempty"`
	DoctorName        string          `json:"doctor_name,omitempty"`
	DoctorShareAmount decimal.Decimal `json:"doctor_share_amount,omitempty"`
	COGSAmount        decimal.Decimal `json:"cogs_amount,omitempty"`
}

// PayInvoice processes patient payment, updates invoice status, and posts double-entry accounting journal automatically.
func (s *Service) PayInvoice(ctx context.Context, cmd PayInvoiceCommand) (*billing.Invoice, *accountingDomain.JournalEntry, error) {
	inv, err := s.billingRepo.GetByID(ctx, cmd.InvoiceID)
	if err != nil {
		return nil, nil, ErrInvoiceNotFound
	}

	if inv.PaymentStatus == billing.PaymentStatusPaid {
		return nil, nil, ErrInvoiceAlreadyPaid
	}

	if !cmd.Amount.IsPositive() {
		return nil, nil, ErrInvalidPaymentAmount
	}

	remaining := inv.TotalAmount.Sub(inv.PaidAmount)
	if cmd.Amount.GreaterThan(remaining) {
		return nil, nil, fmt.Errorf("%w (Remaining: %s, Paid: %s)", ErrOverpayment, remaining, cmd.Amount)
	}

	newPaidAmount := inv.PaidAmount.Add(cmd.Amount)
	newStatus := billing.PaymentStatusPartial
	if newPaidAmount.GreaterThanOrEqual(inv.TotalAmount) {
		newStatus = billing.PaymentStatusPaid
	}

	// 1. Process payment via active channel
	payReq := payment.PaymentRequest{
		PaymentID:   uuid.New(),
		BranchID:    inv.BranchID,
		InvoiceID:   inv.ID,
		Amount:      cmd.Amount,
		Method:      payment.Method(cmd.PaymentMethod),
		ReferenceNo: cmd.ReferenceNo,
		PaidAt:      time.Now(),
	}
	if cmd.ReceivedBy != nil {
		payReq.ReceivedBy = *cmd.ReceivedBy
	}

	var debitAccountCode string
	var journalType accountingDomain.JournalType

	isSSC := cmd.UnitCode == "SSC"

	switch cmd.PaymentMethod {
	case "CASH":
		if s.cashProcessor != nil {
			res, err := s.cashProcessor.ProcessPayment(ctx, payReq)
			if err != nil {
				return nil, nil, fmt.Errorf("cash payment processing failed: %w", err)
			}
			debitAccountCode = res.AccountCode
		} else {
			if isSSC {
				debitAccountCode = "1-1112" // Kas Kasir Semarang Skin Center (SSC)
			} else {
				debitAccountCode = "1-1111" // Kas Kasir Kidinest (Poli Anak)
			}
		}
		journalType = accountingDomain.JournalTypeCashReceipt

	case "BANK_TRANSFER":
		if s.bankProcessor != nil {
			res, err := s.bankProcessor.ProcessPayment(ctx, payReq)
			if err != nil {
				return nil, nil, fmt.Errorf("bank payment processing failed: %w", err)
			}
			debitAccountCode = res.AccountCode
		} else {
			debitAccountCode = "1-1121" // Bank BCA Operasional
		}
		journalType = accountingDomain.JournalTypeBank

	case "QRIS", "MIDTRANS", "XENDIT":
		debitAccountCode = "1-1130" // Kliring Payment Gateway (QRIS/Settlement Transit)
		journalType = accountingDomain.JournalTypeBank

	case "EDC":
		debitAccountCode = "1-1131" // Piutang EDC Mesin Belum Settlement
		journalType = accountingDomain.JournalTypeBank

	case "RECEIVABLE":
		debitAccountCode = "1-1210" // Piutang Pasien Perawatan
		journalType = accountingDomain.JournalTypeSales

	default:
		return nil, nil, fmt.Errorf("unsupported payment method: %s", cmd.PaymentMethod)
	}

	// 2. Prepare Balanced Journal Lines: Cluster 1 (Revenue & Receipt)
	journalLines := []accountingApp.CreateLineCommand{
		{
			AccountCode: debitAccountCode,
			Debit:       cmd.Amount,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Penerimaan pembayaran %s %s", cmd.PaymentMethod, inv.InvoiceNumber),
		},
	}

	// Allocate revenue credit breakdown
	if cmd.Amount.Equal(inv.TotalAmount) {
		if isSSC {
			// Semarang Skin Center (SSC) Revenue Accounts
			if inv.RegistrationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-2150", // Pendapatan Administrasi & Layanan Lain SSC
					Debit:       decimal.Zero,
					Credit:      inv.RegistrationFee,
					Description: "Pendapatan Administrasi SSC",
				})
			}
			if inv.ConsultationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-2110", // Pendapatan Konsultasi Dokter Estetika / Kulit
					Debit:       decimal.Zero,
					Credit:      inv.ConsultationFee,
					Description: "Pendapatan Konsultasi Dokter Estetika",
				})
			}
			if inv.MedicationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-2140", // Pendapatan Penjualan Produk Skincare & Kosmetika
					Debit:       decimal.Zero,
					Credit:      inv.MedicationFee,
					Description: "Pendapatan Produk Skincare & Kosmetika",
				})
			}
		} else {
			// Kidinest Revenue Accounts
			if inv.RegistrationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-1110", // Pendapatan Registrasi & Administrasi Pasien Anak
					Debit:       decimal.Zero,
					Credit:      inv.RegistrationFee,
					Description: "Pendapatan Pendaftaran Pasien Anak",
				})
			}
			if inv.ConsultationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-1120", // Pendapatan Jasa Konsultasi Spesialis Anak
					Debit:       decimal.Zero,
					Credit:      inv.ConsultationFee,
					Description: "Pendapatan Konsultasi Spesialis Anak",
				})
			}
			if inv.MedicationFee.IsPositive() {
				journalLines = append(journalLines, accountingApp.CreateLineCommand{
					AccountCode: "4-1150", // Pendapatan Penjualan Obat Resep Farmasi Anak
					Debit:       decimal.Zero,
					Credit:      inv.MedicationFee,
					Description: "Pendapatan Obat Farmasi Anak",
				})
			}
		}
	} else {
		// Partial payment allocation
		revAcc := "4-1110"
		if isSSC {
			revAcc = "4-2150"
		}
		journalLines = append(journalLines, accountingApp.CreateLineCommand{
			AccountCode: revAcc,
			Debit:       decimal.Zero,
			Credit:      cmd.Amount,
			Description: fmt.Sprintf("Pendapatan Sebagian %s", inv.InvoiceNumber),
		})
	}

	// Cluster 2: COGS & Inventory Deduction (if COGS amount specified)
	if cmd.COGSAmount.IsPositive() {
		cogsAcc := "5-1110" // HPP Obat & Sirup Anak
		invAcc := "1-1310"  // Persediaan Obat Paten & Sirup Anak
		if isSSC {
			cogsAcc = "5-2110" // HPP Produk Skincare & Kosmetika
			invAcc = "1-1330"  // Persediaan Produk Skincare & Kosmetika
		}

		journalLines = append(journalLines,
			accountingApp.CreateLineCommand{
				AccountCode: cogsAcc,
				Debit:       cmd.COGSAmount,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Beban Pokok Penjualan (HPP) %s", inv.InvoiceNumber),
			},
			accountingApp.CreateLineCommand{
				AccountCode: invAcc,
				Debit:       decimal.Zero,
				Credit:      cmd.COGSAmount,
				Description: fmt.Sprintf("Pengurangan Persediaan Farmasi %s", inv.InvoiceNumber),
			},
		)
	}

	// Cluster 3: Doctor Sharing Accrual (if Doctor Share amount specified)
	if cmd.DoctorShareAmount.IsPositive() {
		docExpenseAcc := "5-1210" // Beban Jasa Medis Dokter Spesialis Anak
		docPayableAcc := "2-1210" // Hutang Jasa Medis Dokter Kidinest
		if isSSC {
			docExpenseAcc = "5-2210" // Beban Jasa Medis Dokter Estetika & Kulit
			docPayableAcc = "2-1220" // Hutang Jasa Medis Dokter SSC
		}

		journalLines = append(journalLines,
			accountingApp.CreateLineCommand{
				AccountCode: docExpenseAcc,
				Debit:       cmd.DoctorShareAmount,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Akrual Beban Jasa Medis Dokter (%s) %s", cmd.DoctorName, inv.InvoiceNumber),
			},
			accountingApp.CreateLineCommand{
				AccountCode: docPayableAcc,
				Debit:       decimal.Zero,
				Credit:      cmd.DoctorShareAmount,
				Description: fmt.Sprintf("Akrual Hutang Jasa Medis Dokter (%s) %s", cmd.DoctorName, inv.InvoiceNumber),
			},
		)
	}

	// 3. Post Automatic Journal and Update Invoice in an Atomic Unit of Work
	var journalEntry *accountingDomain.JournalEntry

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// Post journal via Accounting Engine
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        inv.BranchID,
			JournalType:     journalType,
			TransactionDate: time.Now(),
			ReferenceType:   "INVOICE_PAYMENT",
			ReferenceID:     inv.InvoiceNumber,
			Description:     fmt.Sprintf("Pelunasan invoice %s pasien %s", inv.InvoiceNumber, inv.PatientName),
			CreatedBy:       cmd.ReceivedBy,
			AutoPost:        true,
			Lines:           journalLines,
		})
		if err != nil {
			return fmt.Errorf("auto-post journal failed: %w", err)
		}
		journalEntry = jEntry

		// Record payment history
		paymentNumber := fmt.Sprintf("PAY-%s-%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
		payRec := &billing.PaymentRecord{
			InvoiceID:     inv.ID,
			PaymentNumber: paymentNumber,
			Amount:        cmd.Amount,
			PaymentMethod: cmd.PaymentMethod,
			ReferenceNo:   cmd.ReferenceNo,
			Notes:         cmd.Notes,
			ReceivedBy:    cmd.ReceivedBy,
			JournalID:     &jEntry.ID,
			PaidAt:        time.Now(),
		}
		if err := s.billingRepo.AddPaymentRecord(txCtx, payRec); err != nil {
			return fmt.Errorf("add payment record: %w", err)
		}

		// Update invoice status & journal link
		return s.billingRepo.UpdatePayment(txCtx, inv.ID, newPaidAmount, newStatus, &jEntry.ID)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("execute pay invoice uow: %w", err)
	}

	inv.PaidAmount = newPaidAmount
	inv.PaymentStatus = newStatus
	inv.PaymentMethod = cmd.PaymentMethod
	inv.JournalID = &journalEntry.ID

	return inv, journalEntry, nil
}

func (s *Service) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*billing.Invoice, error) {
	return s.billingRepo.GetByID(ctx, id)
}

func (s *Service) GetInvoicesByQueue(ctx context.Context, queueID int64) ([]billing.Invoice, error) {
	return s.billingRepo.GetByQueueID(ctx, queueID)
}

func (s *Service) ListInvoices(ctx context.Context, filter billing.InvoiceFilter) ([]billing.Invoice, int64, error) {
	return s.billingRepo.List(ctx, filter)
}

