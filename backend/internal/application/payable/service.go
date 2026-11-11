package payable

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"clinic-api/internal/application/common"
	journalApp "clinic-api/internal/application/journal"
	"clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/payable"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrPayableNotFound       = errors.New("payable record not found")
	ErrPayableAlreadyPaid    = errors.New("payable is already fully settled")
	ErrInvalidPaymentAmount  = errors.New("payment amount must be positive and cannot exceed outstanding amount")
	ErrZeroAmount            = errors.New("payable total amount must be positive")
)

type CreatePayableCommand struct {
	shared.BranchScoped
	VendorType       payable.VendorType `json:"vendor_type"`
	VendorID         *uuid.UUID         `json:"vendor_id,omitempty"`
	VendorName       string             `json:"vendor_name"`
	InvoiceReference string             `json:"invoice_reference"`
	TotalAmount      decimal.Decimal    `json:"total_amount"`
	IssueDate        time.Time          `json:"issue_date"`
	DueDate          time.Time          `json:"due_date"`
	AccountCode      string             `json:"account_code"` // e.g. 2-1110 (PBF), 2-1120 (Skincare), 2-1210 (Doctor)
	Notes            string             `json:"notes,omitempty"`
}

type PayPayableCommand struct {
	PayableID           uuid.UUID       `json:"payable_id"`
	Amount              decimal.Decimal `json:"amount"`
	PaymentDate         time.Time       `json:"payment_date"`
	CashBankAccountCode string          `json:"cash_bank_account_code"` // e.g. 1-1121 (Bank BCA), 1-1113 (Petty Cash)
	ReferenceNo         string          `json:"reference_no,omitempty"`
	Notes               string          `json:"notes,omitempty"`
	PaidBy              *uuid.UUID      `json:"paid_by,omitempty"`
}

type Service struct {
	uow            common.UnitOfWork
	payableRepo    payable.Repository
	accountRepo    accounting.AccountRepository
	journalService *journalApp.Service
	seqCounter     uint64
}

func NewService(
	uow common.UnitOfWork,
	payableRepo payable.Repository,
	accountRepo accounting.AccountRepository,
	journalService *journalApp.Service,
) *Service {
	return &Service{
		uow:            uow,
		payableRepo:    payableRepo,
		accountRepo:    accountRepo,
		journalService: journalService,
		seqCounter:     uint64(time.Now().UnixNano() % 100000),
	}
}

// CreatePayable records an obligation to pay a pharmacy PBF supplier, skincare vendor, or doctor.
func (s *Service) CreatePayable(ctx context.Context, cmd CreatePayableCommand) (*payable.Payable, error) {
	if !cmd.TotalAmount.IsPositive() {
		return nil, ErrZeroAmount
	}
	if cmd.IssueDate.IsZero() {
		cmd.IssueDate = time.Now()
	}
	if cmd.DueDate.IsZero() {
		cmd.DueDate = cmd.IssueDate.AddDate(0, 0, 30) // Default net 30 days
	}

	accCode := cmd.AccountCode
	if accCode == "" {
		accCode = "2-1110" // Default Hutang Supplier Farmasi & PBF
	}

	acc, err := s.accountRepo.GetByCode(ctx, accCode)
	if err != nil {
		return nil, fmt.Errorf("lookup payable account %s: %w", accCode, err)
	}

	now := time.Now()
	payNumber := s.generatePayableNumber(cmd.IssueDate)

	p := &payable.Payable{
		BaseEntity: shared.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		BranchScoped: shared.BranchScoped{
			BranchID: cmd.BranchID,
		},
		PayableNumber:     payNumber,
		VendorType:        cmd.VendorType,
		VendorID:          cmd.VendorID,
		VendorName:        cmd.VendorName,
		InvoiceReference:  cmd.InvoiceReference,
		TotalAmount:       cmd.TotalAmount,
		OutstandingAmount: cmd.TotalAmount,
		IssueDate:         cmd.IssueDate,
		DueDate:           cmd.DueDate,
		Status:            payable.StatusOpen,
		AccountID:         acc.ID,
		AccountCode:       acc.Code,
		Notes:             cmd.Notes,
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		return s.payableRepo.Create(txCtx, p)
	})
	if err != nil {
		return nil, fmt.Errorf("persist payable: %w", err)
	}

	return p, nil
}

// PayPayable records a payment against a supplier invoice or doctor fee and auto-posts a Cash Disbursement Journal (JKK).
func (s *Service) PayPayable(ctx context.Context, cmd PayPayableCommand) (*payable.Payment, error) {
	if !cmd.Amount.IsPositive() {
		return nil, ErrInvalidPaymentAmount
	}
	if cmd.PaymentDate.IsZero() {
		cmd.PaymentDate = time.Now()
	}

	p, err := s.payableRepo.GetByID(ctx, cmd.PayableID)
	if err != nil {
		return nil, ErrPayableNotFound
	}

	if p.Status == payable.StatusPaid || p.OutstandingAmount.IsZero() {
		return nil, ErrPayableAlreadyPaid
	}

	if cmd.Amount.GreaterThan(p.OutstandingAmount) {
		return nil, fmt.Errorf("%w: payment %s exceeds outstanding %s",
			ErrInvalidPaymentAmount, cmd.Amount.StringFixed(2), p.OutstandingAmount.StringFixed(2))
	}

	newOutstanding := p.OutstandingAmount.Sub(cmd.Amount)
	newStatus := payable.StatusPartial
	if newOutstanding.IsZero() {
		newStatus = payable.StatusPaid
	}

	payRecordNumber := s.generatePaymentRecordNumber(cmd.PaymentDate)

	cashAccountCode := cmd.CashBankAccountCode
	if cashAccountCode == "" {
		cashAccountCode = "1-1121" // Bank BCA Operasional
	}

	payment := &payable.Payment{
		ID:                  uuid.New(),
		PayableID:           p.ID,
		PaymentNumber:       payRecordNumber,
		Amount:              cmd.Amount,
		PaymentDate:         cmd.PaymentDate,
		CashBankAccountCode: cashAccountCode,
		ReferenceNo:         cmd.ReferenceNo,
		Notes:               cmd.Notes,
		CreatedBy:           cmd.PaidBy,
		CreatedAt:           time.Now(),
	}

	// Double-Entry Journal:
	// Debit: Hutang Supplier / Hutang Dokter (e.g. 2-1110, 2-1120, 2-1210)
	// Kredit: Bank BCA / Kasir (e.g. 1-1121, 1-1111)
	journalCmd := journalApp.CreateJournalCommand{
		BranchScoped:    shared.BranchScoped{BranchID: p.BranchID},
		JournalType:     accounting.JournalTypeCashDisbursement,
		TransactionDate: cmd.PaymentDate,
		ReferenceType:   "PAYABLE_PAYMENT",
		ReferenceID:     p.PayableNumber,
		Description:     fmt.Sprintf("Pembayaran Hutang %s - Vendor: %s", p.PayableNumber, p.VendorName),
		CreatedBy:       cmd.PaidBy,
		Lines: []journalApp.CreateLineCommand{
			{
				AccountCode: p.AccountCode,
				Debit:       cmd.Amount,
				Credit:      decimal.Zero,
				Description: fmt.Sprintf("Pelunasan Hutang Faktur %s", p.InvoiceReference),
			},
			{
				AccountCode: cashAccountCode,
				Debit:       decimal.Zero,
				Credit:      cmd.Amount,
				Description: fmt.Sprintf("Pengeluaran Kas/Bank untuk %s", p.VendorName),
			},
		},
	}

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		// 1. Post JKK
		jEntry, err := s.journalService.CreateJournal(txCtx, journalCmd)
		if err != nil {
			return fmt.Errorf("post payable cash disbursement journal: %w", err)
		}
		payment.JournalID = &jEntry.ID

		// 2. Persist Payment Record
		if err := s.payableRepo.AddPayment(txCtx, payment); err != nil {
			return fmt.Errorf("record payable payment: %w", err)
		}

		// 3. Update Outstanding & Status
		if err := s.payableRepo.UpdateOutstanding(txCtx, p.ID, newOutstanding, newStatus); err != nil {
			return fmt.Errorf("update payable status: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("execute payable payment unit of work: %w", err)
	}

	return payment, nil
}

func (s *Service) GetPayableByID(ctx context.Context, id uuid.UUID) (*payable.Payable, error) {
	return s.payableRepo.GetByID(ctx, id)
}

func (s *Service) GetPayableByNumber(ctx context.Context, number string) (*payable.Payable, error) {
	return s.payableRepo.GetByNumber(ctx, number)
}

func (s *Service) ListPayables(ctx context.Context, filter payable.Filter) ([]payable.Payable, int64, error) {
	return s.payableRepo.List(ctx, filter)
}

func (s *Service) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*payable.AgingReport, error) {
	if asOfDate.IsZero() {
		asOfDate = time.Now()
	}
	return s.payableRepo.GetAgingReport(ctx, branchID, asOfDate)
}

func (s *Service) generatePayableNumber(date time.Time) string {
	seq := atomic.AddUint64(&s.seqCounter, 1) % 10000
	return fmt.Sprintf("AP-%s-%04d", date.Format("20060102"), seq)
}

func (s *Service) generatePaymentRecordNumber(date time.Time) string {
	seq := atomic.AddUint64(&s.seqCounter, 1) % 10000
	return fmt.Sprintf("PAY-AP-%s-%04d", date.Format("20060102"), seq)
}

