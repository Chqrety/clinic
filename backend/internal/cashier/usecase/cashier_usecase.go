package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"clinic-api/internal/cashier/domain"
	queueDomain "clinic-api/internal/queue/domain"
)

type WhatsAppSender interface {
	SendMessage(phone, message string) error
}

type QueueNumberGenerator interface {
	GenerateNextQueueNumber(ctx context.Context, date time.Time, unit queueDomain.ClinicUnit, category queueDomain.QueueCategory) (string, error)
}

type CashierUsecase struct {
	invoiceRepo   domain.InvoiceRepository
	shiftRepo     domain.ShiftRepository
	cashierQRepo  domain.CashierQueueRepository
	counterRepo   QueueNumberGenerator
	waAdapter     WhatsAppSender
	publicBaseURL string
}

func NewCashierUsecase(
	invoiceRepo domain.InvoiceRepository,
	shiftRepo domain.ShiftRepository,
	cashierQRepo domain.CashierQueueRepository,
	counterRepo QueueNumberGenerator,
	waAdapter WhatsAppSender,
	publicBaseURL string,
) *CashierUsecase {
	if publicBaseURL == "" {
		publicBaseURL = "http://localhost:5173"
	}
	return &CashierUsecase{
		invoiceRepo:   invoiceRepo,
		shiftRepo:     shiftRepo,
		cashierQRepo:  cashierQRepo,
		counterRepo:   counterRepo,
		waAdapter:     waAdapter,
		publicBaseURL: publicBaseURL,
	}
}

func (uc *CashierUsecase) RegisterWalkIn(ctx context.Context, req *domain.WalkInRequest, cashierStaffID int64, cashierName string) (*domain.PendingBill, *domain.Invoice, error) {
	req.PatientName = strings.TrimSpace(req.PatientName)
	req.PatientPhone = strings.TrimSpace(req.PatientPhone)

	if req.PatientName == "" {
		return nil, nil, domain.ErrPatientNameRequired
	}
	if req.PatientPhone == "" {
		return nil, nil, domain.ErrPatientPhoneRequired
	}

	var prefix string
	var regFee int

	if req.ClinicUnit == "SSC" {
		prefix = "S"
		regFee = 0
	} else {
		req.ClinicUnit = "KIDINEST"
		prefix = "E"
		regFee = 50000
	}

	// Generate queue number
	var qNum string
	if uc.counterRepo != nil {
		cat := queueDomain.CategoryWalkin
		if req.ClinicUnit == "SSC" {
			cat = queueDomain.CategoryTreatment
		}
		num, err := uc.counterRepo.GenerateNextQueueNumber(ctx, time.Now(), queueDomain.ClinicUnit(req.ClinicUnit), cat)
		if err != nil {
			return nil, nil, fmt.Errorf("gagal generate nomor antrian walk-in: %w", err)
		}
		qNum = num
	} else {
		qNum = fmt.Sprintf("%s-%03d", prefix, time.Now().Unix()%900+100)
	}

	queueID, _, err := uc.cashierQRepo.RegisterWalkIn(ctx, req, qNum, regFee)
	if err != nil {
		return nil, nil, err
	}

	bill, err := uc.cashierQRepo.GetQueueByID(ctx, queueID)
	if err != nil {
		return nil, nil, err
	}

	// If paid immediate (Cash at frontdesk), automatically generate invoice
	var createdInvoice *domain.Invoice
	if req.IsPaidImmediate && regFee > 0 {
		invNum, err := uc.invoiceRepo.GenerateInvoiceNumber(ctx)
		if err != nil {
			invNum = fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%1000)
		}

		inv := &domain.Invoice{
			QueueID:              queueID,
			PatientID:            bill.PatientID,
			InvoiceNumber:        invNum,
			PatientName:          req.PatientName,
			PatientPhone:         req.PatientPhone,
			ClinicUnit:           req.ClinicUnit,
			QueueNumber:          qNum,
			SubtotalRegistration: regFee,
			SubtotalAction:       0,
			SubtotalMedication:   0,
			VoucherDiscount:      0,
			TotalAmount:          regFee,
			PaymentMethod:        domain.PaymentMethodCash,
			CashReceived:         regFee,
			CashChange:           0,
			PaymentReference:     "WALKIN-CASH-FRONTDESK",
			CashierStaffID:       cashierStaffID,
			CashierName:          cashierName,
			Status:               domain.InvoiceStatusPaid,
			Notes:                "Pendaftaran Walk-In Lunas di Meja Depan",
		}
		if err := uc.invoiceRepo.Create(ctx, inv); err == nil {
			createdInvoice = inv
		}
	}

	// Send WhatsApp ticket notification
	if uc.waAdapter != nil {
		ticketURL := fmt.Sprintf("%s/ticket/%d", uc.publicBaseURL, queueID)
		msg := fmt.Sprintf(
			"Halo %s, pendaftaran walk-in Anda di %s berhasil! Nomor Antrian Anda: *%s*. Dokter: %s. Pantau antrian live di: %s",
			req.PatientName, req.ClinicUnit, qNum, req.DoctorName, ticketURL,
		)
		_ = uc.waAdapter.SendMessage(req.PatientPhone, msg)
	}

	return bill, createdInvoice, nil
}

func (uc *CashierUsecase) ListPendingBills(ctx context.Context) ([]*domain.PendingBill, error) {
	return uc.cashierQRepo.ListPendingBills(ctx)
}

func (uc *CashierUsecase) PayBill(ctx context.Context, req *domain.PayBillRequest, cashierStaffID int64, cashierName string) (*domain.Invoice, error) {
	if req.QueueID <= 0 {
		return nil, domain.ErrQueueNotFound
	}

	bill, err := uc.cashierQRepo.GetQueueByID(ctx, req.QueueID)
	if err != nil {
		return nil, err
	}

	grossTotal := req.SubtotalRegistration + req.SubtotalAction + req.SubtotalMedication
	if grossTotal < 0 {
		return nil, domain.ErrInvalidAmount
	}

	if req.VoucherDiscount > grossTotal {
		return nil, domain.ErrVoucherExceedsTotal
	}
	if req.VoucherDiscount < 0 {
		req.VoucherDiscount = 0
	}

	netTotal := grossTotal - req.VoucherDiscount

	// Cash validation
	cashChange := 0
	if req.PaymentMethod == domain.PaymentMethodCash {
		if req.CashReceived < netTotal {
			return nil, domain.ErrInsufficientCash
		}
		cashChange = req.CashReceived - netTotal
	} else {
		req.CashReceived = netTotal
		cashChange = 0
	}

	invNum, err := uc.invoiceRepo.GenerateInvoiceNumber(ctx)
	if err != nil {
		invNum = fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), time.Now().Unix()%1000)
	}

	invoice := &domain.Invoice{
		QueueID:              req.QueueID,
		PatientID:            bill.PatientID,
		InvoiceNumber:        invNum,
		PatientName:          bill.PatientName,
		PatientPhone:         bill.PatientPhone,
		ClinicUnit:           bill.ClinicUnit,
		QueueNumber:          bill.QueueNumber,
		SubtotalRegistration: req.SubtotalRegistration,
		SubtotalAction:       req.SubtotalAction,
		SubtotalMedication:   req.SubtotalMedication,
		VoucherDiscount:      req.VoucherDiscount,
		TotalAmount:          netTotal,
		PaymentMethod:        req.PaymentMethod,
		CashReceived:         req.CashReceived,
		CashChange:           cashChange,
		PaymentReference:     req.PaymentReference,
		CashierStaffID:       cashierStaffID,
		CashierName:          cashierName,
		Status:               domain.InvoiceStatusPaid,
		Notes:                req.Notes,
	}

	if err := uc.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("gagal membuat invoice kasir: %w", err)
	}

	// Update queue status
	nextStatus := "WAITING"
	if bill.QueueStatus == "IN_EXAM" || bill.QueueStatus == "PHARMACY_WAITING" || bill.QueueStatus == "COMPLETED" {
		nextStatus = "COMPLETED"
	}
	_ = uc.cashierQRepo.UpdatePaymentAndQueueStatus(ctx, req.QueueID, "PAID", nextStatus)

	return invoice, nil
}

func (uc *CashierUsecase) GetInvoice(ctx context.Context, id int64) (*domain.Invoice, error) {
	return uc.invoiceRepo.GetByID(ctx, id)
}

func (uc *CashierUsecase) GetInvoiceByQueue(ctx context.Context, queueID int64) (*domain.Invoice, error) {
	return uc.invoiceRepo.GetByQueueID(ctx, queueID)
}

func (uc *CashierUsecase) ListInvoices(ctx context.Context, limit, offset int) ([]*domain.Invoice, error) {
	return uc.invoiceRepo.List(ctx, limit, offset)
}

func (uc *CashierUsecase) ProcessRefund(ctx context.Context, invoiceID int64, req *domain.RefundRequest, staffID int64, staffName string) (*domain.Invoice, error) {
	req.RefundBankName = strings.TrimSpace(req.RefundBankName)
	req.RefundBankAccount = strings.TrimSpace(req.RefundBankAccount)
	req.RefundAccountName = strings.TrimSpace(req.RefundAccountName)
	req.RefundReason = strings.TrimSpace(req.RefundReason)

	if req.RefundBankName == "" || req.RefundBankAccount == "" || req.RefundReason == "" {
		return nil, domain.ErrRefundReasonRequired
	}

	inv, err := uc.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	if inv.Status == domain.InvoiceStatusRefunded {
		return nil, domain.ErrAlreadyRefunded
	}

	now := time.Now()
	err = uc.invoiceRepo.UpdateRefund(
		ctx, invoiceID,
		req.RefundBankName, req.RefundBankAccount, req.RefundAccountName, req.RefundReason,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal memproses refund: %w", err)
	}

	// Cancel queue
	_ = uc.cashierQRepo.UpdatePaymentAndQueueStatus(ctx, inv.QueueID, "REFUNDED", "CANCELLED")

	log.Printf("[CASHIER] Refund diproses untuk Invoice %s senilai Rp %d. Tujuan: %s %s a/n %s. Alasan: %s",
		inv.InvoiceNumber, inv.TotalAmount, req.RefundBankName, req.RefundBankAccount, req.RefundAccountName, req.RefundReason,
	)

	return uc.invoiceRepo.GetByID(ctx, invoiceID)
}

func (uc *CashierUsecase) GetShiftSummary(ctx context.Context, date string) (*domain.ShiftSummary, error) {
	return uc.invoiceRepo.GetShiftSummary(ctx, date)
}

func (uc *CashierUsecase) CloseShift(ctx context.Context, staffID int64, staffName string, req *domain.CloseShiftRequest) (*domain.CashierShiftReport, error) {
	today := time.Now().Format("2006-01-02")
	summary, err := uc.invoiceRepo.GetShiftSummary(ctx, today)
	if err != nil {
		return nil, fmt.Errorf("gagal merangkum shift kasir: %w", err)
	}

	report := &domain.CashierShiftReport{
		StaffID:           staffID,
		StaffName:         staffName,
		ShiftDate:         today,
		TotalTransactions: summary.TotalTransactions,
		TotalCash:         summary.TotalCash,
		TotalTransfer:     summary.TotalTransfer,
		TotalEDC:          summary.TotalEDC,
		TotalVoucher:      summary.TotalVoucher,
		TotalRefund:       summary.TotalRefund,
		NetTotal:          summary.NetTotal,
		Notes:             req.Notes,
	}

	if err := uc.shiftRepo.CreateReport(ctx, report); err != nil {
		return nil, fmt.Errorf("gagal menyimpan laporan tutup shift: %w", err)
	}

	log.Printf("[CASHIER] Tutup shift berhasil dikirim ke Pimpinan. Kasir: %s, Net Total Setoran: Rp %d",
		staffName, report.NetTotal,
	)

	return report, nil
}

func (uc *CashierUsecase) ListShiftReports(ctx context.Context, limit int) ([]*domain.CashierShiftReport, error) {
	return uc.shiftRepo.ListReports(ctx, limit)
}
