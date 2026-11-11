package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"clinic-api/internal/cashier/domain"
	"clinic-api/internal/cashier/usecase"
	queueDomain "clinic-api/internal/queue/domain"
)

type mockInvoiceRepo struct {
	invoices []*domain.Invoice
	summary  *domain.ShiftSummary
}

func (m *mockInvoiceRepo) Create(ctx context.Context, inv *domain.Invoice) error {
	inv.ID = int64(len(m.invoices) + 1)
	inv.CreatedAt = time.Now()
	m.invoices = append(m.invoices, inv)
	return nil
}

func (m *mockInvoiceRepo) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	for _, inv := range m.invoices {
		if inv.ID == id {
			return inv, nil
		}
	}
	return nil, domain.ErrInvoiceNotFound
}

func (m *mockInvoiceRepo) GetByQueueID(ctx context.Context, queueID int64) (*domain.Invoice, error) {
	for _, inv := range m.invoices {
		if inv.QueueID == queueID {
			return inv, nil
		}
	}
	return nil, domain.ErrInvoiceNotFound
}

func (m *mockInvoiceRepo) List(ctx context.Context, limit, offset int) ([]*domain.Invoice, error) {
	return m.invoices, nil
}

func (m *mockInvoiceRepo) UpdateRefund(ctx context.Context, id int64, bankName, bankAccount, accountName, reason string, refundedAt time.Time) error {
	for _, inv := range m.invoices {
		if inv.ID == id {
			inv.Status = domain.InvoiceStatusRefunded
			inv.RefundBankName = bankName
			inv.RefundBankAccount = bankAccount
			inv.RefundAccountName = accountName
			inv.RefundReason = reason
			inv.RefundedAt = &refundedAt
			return nil
		}
	}
	return domain.ErrInvoiceNotFound
}

func (m *mockInvoiceRepo) GetShiftSummary(ctx context.Context, date string) (*domain.ShiftSummary, error) {
	if m.summary != nil {
		return m.summary, nil
	}
	return &domain.ShiftSummary{
		ShiftDate:         date,
		TotalTransactions: len(m.invoices),
		TotalCash:         100000,
		TotalTransfer:     50000,
		TotalEDC:          50000,
		TotalVoucher:      20000,
		TotalRefund:       0,
		NetTotal:          200000,
	}, nil
}

func (m *mockInvoiceRepo) GenerateInvoiceNumber(ctx context.Context) (string, error) {
	return fmt.Sprintf("INV-20261111-%04d", len(m.invoices)+1), nil
}

type mockShiftRepo struct {
	reports []*domain.CashierShiftReport
}

func (m *mockShiftRepo) CreateReport(ctx context.Context, report *domain.CashierShiftReport) error {
	report.ID = int64(len(m.reports) + 1)
	m.reports = append(m.reports, report)
	return nil
}

func (m *mockShiftRepo) ListReports(ctx context.Context, limit int) ([]*domain.CashierShiftReport, error) {
	return m.reports, nil
}

type mockCashierQRepo struct {
	queues map[int64]*domain.PendingBill
}

func (m *mockCashierQRepo) GetQueueByID(ctx context.Context, id int64) (*domain.PendingBill, error) {
	if q, ok := m.queues[id]; ok {
		return q, nil
	}
	return nil, domain.ErrQueueNotFound
}

func (m *mockCashierQRepo) ListPendingBills(ctx context.Context) ([]*domain.PendingBill, error) {
	var list []*domain.PendingBill
	for _, q := range m.queues {
		list = append(list, q)
	}
	return list, nil
}

func (m *mockCashierQRepo) UpdatePaymentAndQueueStatus(ctx context.Context, queueID int64, paymentStatus, queueStatus string) error {
	if q, ok := m.queues[queueID]; ok {
		q.PaymentStatus = paymentStatus
		q.QueueStatus = queueStatus
	}
	return nil
}

func (m *mockCashierQRepo) RegisterWalkIn(ctx context.Context, req *domain.WalkInRequest, queueNumber string, fee int) (int64, string, error) {
	newID := int64(len(m.queues) + 1)
	bill := &domain.PendingBill{
		QueueID:              newID,
		QueueNumber:          queueNumber,
		PatientID:            101,
		PatientName:          req.PatientName,
		PatientPhone:         req.PatientPhone,
		ClinicUnit:           req.ClinicUnit,
		Category:             "E_WALKIN",
		DoctorName:           req.DoctorName,
		SubtotalRegistration: fee,
		TotalDue:             fee,
		QueueStatus:          "WAITING",
		PaymentStatus:        "PENDING",
		CreatedAt:            time.Now(),
	}
	if req.IsPaidImmediate {
		bill.PaymentStatus = "PAID"
	}
	m.queues[newID] = bill
	return newID, queueNumber, nil
}

type mockCounterRepo struct {
	count int
}

func (m *mockCounterRepo) GenerateNextQueueNumber(ctx context.Context, date time.Time, unit queueDomain.ClinicUnit, category queueDomain.QueueCategory) (string, error) {
	m.count++
	prefix := "E"
	if unit == "SSC" {
		prefix = "S"
	}
	return fmt.Sprintf("%s-%03d", prefix, m.count), nil
}

func setupCashierUsecase() (*usecase.CashierUsecase, *mockInvoiceRepo, *mockShiftRepo, *mockCashierQRepo) {
	invRepo := &mockInvoiceRepo{invoices: make([]*domain.Invoice, 0)}
	shiftRepo := &mockShiftRepo{reports: make([]*domain.CashierShiftReport, 0)}
	qRepo := &mockCashierQRepo{queues: make(map[int64]*domain.PendingBill)}
	counterRepo := &mockCounterRepo{}

	uc := usecase.NewCashierUsecase(invRepo, shiftRepo, qRepo, counterRepo, nil, "http://localhost:5173")
	return uc, invRepo, shiftRepo, qRepo
}

func TestRegisterWalkIn(t *testing.T) {
	uc, _, _, _ := setupCashierUsecase()
	ctx := context.Background()

	// 1. Missing name
	_, _, err := uc.RegisterWalkIn(ctx, &domain.WalkInRequest{
		PatientName:  "",
		PatientPhone: "08123456789",
	}, 1, "Kasir Siti")
	if err != domain.ErrPatientNameRequired {
		t.Errorf("expected ErrPatientNameRequired, got %v", err)
	}

	// 2. Successful Kidinest WalkIn paid immediate
	bill, inv, err := uc.RegisterWalkIn(ctx, &domain.WalkInRequest{
		PatientName:     "Ananda Budi",
		PatientPhone:    "081234567890",
		ClinicUnit:      "KIDINEST",
		DoctorID:        1,
		DoctorName:      "dr. Hendra Pratama, Sp.A",
		PaymentMethod:   "CASH",
		IsPaidImmediate: true,
	}, 4, "Siti Sarah")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bill.QueueNumber != "E-001" {
		t.Errorf("expected queue number E-001, got %s", bill.QueueNumber)
	}
	if inv == nil {
		t.Fatalf("expected invoice to be auto-created for paid immediate walkin")
	}
	if inv.TotalAmount != 50000 {
		t.Errorf("expected total amount 50000, got %d", inv.TotalAmount)
	}
}

func TestPayBillCalculationAndVoucher(t *testing.T) {
	uc, _, _, qRepo := setupCashierUsecase()
	ctx := context.Background()

	qRepo.queues[10] = &domain.PendingBill{
		QueueID:              10,
		QueueNumber:          "B-005",
		PatientID:            55,
		PatientName:          "Ibu Kartika",
		PatientPhone:         "081299998888",
		ClinicUnit:           "KIDINEST",
		DoctorName:           "dr. Hendra",
		SubtotalRegistration: 50000,
		QueueStatus:          "COMPLETED",
		PaymentStatus:        "PENDING",
	}

	// 1. Voucher exceeds total
	_, err := uc.PayBill(ctx, &domain.PayBillRequest{
		QueueID:              10,
		SubtotalRegistration: 50000,
		SubtotalAction:       100000,
		SubtotalMedication:   50000,
		VoucherDiscount:      250000, // exceeds gross 200,000
		PaymentMethod:        domain.PaymentMethodCash,
		CashReceived:         250000,
	}, 4, "Siti Sarah")
	if err != domain.ErrVoucherExceedsTotal {
		t.Errorf("expected ErrVoucherExceedsTotal, got %v", err)
	}

	// 2. Insufficient cash
	_, err = uc.PayBill(ctx, &domain.PayBillRequest{
		QueueID:              10,
		SubtotalRegistration: 50000,
		SubtotalAction:       100000,
		SubtotalMedication:   50000,
		VoucherDiscount:      20000, // net = 180,000
		PaymentMethod:        domain.PaymentMethodCash,
		CashReceived:         150000, // less than 180,000
	}, 4, "Siti Sarah")
	if err != domain.ErrInsufficientCash {
		t.Errorf("expected ErrInsufficientCash, got %v", err)
	}

	// 3. Valid payment with correct change
	inv, err := uc.PayBill(ctx, &domain.PayBillRequest{
		QueueID:              10,
		SubtotalRegistration: 50000,
		SubtotalAction:       100000,
		SubtotalMedication:   50000,
		VoucherDiscount:      20000, // net = 180,000
		PaymentMethod:        domain.PaymentMethodCash,
		CashReceived:         200000, // change = 20,000
	}, 4, "Siti Sarah")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.TotalAmount != 180000 {
		t.Errorf("expected total amount 180000, got %d", inv.TotalAmount)
	}
	if inv.CashChange != 20000 {
		t.Errorf("expected cash change 20000, got %d", inv.CashChange)
	}
	if inv.Status != domain.InvoiceStatusPaid {
		t.Errorf("expected status PAID, got %s", inv.Status)
	}
}

func TestProcessRefund(t *testing.T) {
	uc, invRepo, _, qRepo := setupCashierUsecase()
	ctx := context.Background()

	qRepo.queues[20] = &domain.PendingBill{
		QueueID:       20,
		QueueNumber:   "A-001",
		PaymentStatus: "PAID",
		QueueStatus:   "WAITING",
	}

	invRepo.invoices = append(invRepo.invoices, &domain.Invoice{
		ID:            1,
		QueueID:       20,
		TotalAmount:   50000,
		PaymentMethod: domain.PaymentMethodCash,
		Status:        domain.InvoiceStatusPaid,
	})

	// Missing refund details
	_, err := uc.ProcessRefund(ctx, 1, &domain.RefundRequest{
		RefundBankName: "",
	}, 4, "Siti Sarah")
	if err != domain.ErrRefundReasonRequired {
		t.Errorf("expected ErrRefundReasonRequired, got %v", err)
	}

	// Successful refund
	refunded, err := uc.ProcessRefund(ctx, 1, &domain.RefundRequest{
		RefundBankName:    "BCA",
		RefundBankAccount: "1234567890",
		RefundAccountName: "Budi Santoso",
		RefundReason:      "Pasien mendadak demam tinggi harus rujukan RS",
	}, 4, "Siti Sarah")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refunded.Status != domain.InvoiceStatusRefunded {
		t.Errorf("expected status REFUNDED, got %s", refunded.Status)
	}

	// Double refund attempt
	_, err = uc.ProcessRefund(ctx, 1, &domain.RefundRequest{
		RefundBankName:    "BCA",
		RefundBankAccount: "1234567890",
		RefundReason:      "Retur lagi",
	}, 4, "Siti Sarah")
	if err != domain.ErrAlreadyRefunded {
		t.Errorf("expected ErrAlreadyRefunded, got %v", err)
	}
}

func TestCloseShift(t *testing.T) {
	uc, _, shiftRepo, _ := setupCashierUsecase()
	ctx := context.Background()

	report, err := uc.CloseShift(ctx, 4, "Siti Sarah", &domain.CloseShiftRequest{
		Notes: "Shift pagi selesai lancar tanpa selisih fisik",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.NetTotal != 200000 {
		t.Errorf("expected net total 200000, got %d", report.NetTotal)
	}
	if len(shiftRepo.reports) != 1 {
		t.Errorf("expected 1 report saved, got %d", len(shiftRepo.reports))
	}
}
