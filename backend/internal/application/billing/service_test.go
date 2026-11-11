package billing_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	billingApp "clinic-api/internal/application/billing"
	accountingDomain "clinic-api/internal/domain/accounting"
	billingDomain "clinic-api/internal/domain/billing"
	"clinic-api/internal/adapters/payment"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockBillingRepo struct {
	invoices map[uuid.UUID]*billingDomain.Invoice
	payments []*billingDomain.PaymentRecord
}

func newMockBillingRepo() *mockBillingRepo {
	return &mockBillingRepo{
		invoices: make(map[uuid.UUID]*billingDomain.Invoice),
	}
}

func (m *mockBillingRepo) Create(ctx context.Context, inv *billingDomain.Invoice) error {
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	m.invoices[inv.ID] = inv
	return nil
}

func (m *mockBillingRepo) GetByID(ctx context.Context, id uuid.UUID) (*billingDomain.Invoice, error) {
	inv, ok := m.invoices[id]
	if !ok {
		return nil, billingApp.ErrInvoiceNotFound
	}
	return inv, nil
}

func (m *mockBillingRepo) GetByNumber(ctx context.Context, number string) (*billingDomain.Invoice, error) {
	for _, inv := range m.invoices {
		if inv.InvoiceNumber == number {
			return inv, nil
		}
	}
	return nil, billingApp.ErrInvoiceNotFound
}

func (m *mockBillingRepo) GetByQueueID(ctx context.Context, queueID int64) ([]billingDomain.Invoice, error) {
	var list []billingDomain.Invoice
	for _, inv := range m.invoices {
		if inv.QueueID != nil && *inv.QueueID == queueID {
			list = append(list, *inv)
		}
	}
	return list, nil
}

func (m *mockBillingRepo) List(ctx context.Context, filter billingDomain.InvoiceFilter) ([]billingDomain.Invoice, int64, error) {
	var list []billingDomain.Invoice
	for _, inv := range m.invoices {
		list = append(list, *inv)
	}
	return list, int64(len(list)), nil
}

func (m *mockBillingRepo) UpdatePayment(ctx context.Context, id uuid.UUID, paidAmount decimal.Decimal, status billingDomain.PaymentStatus, journalID *uuid.UUID) error {
	inv, ok := m.invoices[id]
	if !ok {
		return billingApp.ErrInvoiceNotFound
	}
	inv.PaidAmount = paidAmount
	inv.PaymentStatus = status
	inv.JournalID = journalID
	return nil
}

func (m *mockBillingRepo) AddPaymentRecord(ctx context.Context, record *billingDomain.PaymentRecord) error {
	m.payments = append(m.payments, record)
	return nil
}

// Mock Accounting Repositories
type mockAccountRepo struct {
	accounts map[string]*accountingDomain.Account
}

func newMockAccountRepo() *mockAccountRepo {
	r := &mockAccountRepo{accounts: make(map[string]*accountingDomain.Account)}
	seed := []struct {
		code, name string
		accType    accountingDomain.AccountType
		normal     accountingDomain.NormalBalance
	}{
		{"1-1111", "Kas Kasir", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1112", "Kas Kasir SSC", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1113", "Bank BCA", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1121", "Bank BCA Operasional", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1200", "Piutang Pasien", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1210", "Piutang Pasien Perawatan", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"4-1100", "Pendapatan Pendaftaran", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1110", "Pendapatan Pendaftaran Kidinest", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1120", "Pendapatan Konsultasi Dokter Anak", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1150", "Pendapatan Farmasi Anak", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1200", "Pendapatan Konsultasi", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"4-1300", "Pendapatan Obat", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
	}
	for _, s := range seed {
		acc := &accountingDomain.Account{
			Code:          s.code,
			Name:          s.name,
			Type:          s.accType,
			NormalBalance: s.normal,
			IsActive:      true,
			IsHeader:      false,
		}
		acc.ID = uuid.New()
		r.accounts[s.code] = acc
	}
	return r
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.Account, error) {
	for _, a := range m.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, accountingApp.ErrAccountNotFound
}

func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*accountingDomain.Account, error) {
	a, ok := m.accounts[code]
	if !ok {
		return nil, accountingApp.ErrAccountNotFound
	}
	return a, nil
}

func (m *mockAccountRepo) List(ctx context.Context, filter accountingDomain.AccountFilter) ([]accountingDomain.Account, int64, error) {
	var list []accountingDomain.Account
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func (m *mockAccountRepo) Create(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.Code] = acc
	return nil
}

func (m *mockAccountRepo) Update(ctx context.Context, acc *accountingDomain.Account) error {
	m.accounts[acc.Code] = acc
	return nil
}

type mockJournalRepo struct {
	journals map[uuid.UUID]*accountingDomain.JournalEntry
}

func (m *mockJournalRepo) Create(ctx context.Context, entry *accountingDomain.JournalEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	m.journals[entry.ID] = entry
	return nil
}

func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.JournalEntry, error) {
	j, ok := m.journals[id]
	if !ok {
		return nil, accountingApp.ErrInvalidJournalStatus
	}
	return j, nil
}

func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.JournalNumber == number {
			return j, nil
		}
	}
	return nil, accountingApp.ErrInvalidJournalStatus
}

func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.ReferenceType == refType && j.ReferenceID == refID {
			return j, nil
		}
	}
	return nil, accountingApp.ErrInvalidJournalStatus
}

func (m *mockJournalRepo) List(ctx context.Context, filter accountingDomain.JournalFilter) ([]accountingDomain.JournalEntry, int64, error) {
	var list []accountingDomain.JournalEntry
	for _, j := range m.journals {
		list = append(list, *j)
	}
	return list, int64(len(list)), nil
}

func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accountingDomain.JournalStatus, postedBy *uuid.UUID) error {
	j, ok := m.journals[id]
	if !ok {
		return accountingApp.ErrInvalidJournalStatus
	}
	j.Status = status
	return nil
}

type mockFiscalPeriodRepo struct{}

func (m *mockFiscalPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accountingDomain.FiscalPeriod, error) {
	return &accountingDomain.FiscalPeriod{Status: "OPEN"}, nil
}
func (m *mockFiscalPeriodRepo) ClosePeriod(ctx context.Context, period *accountingDomain.FiscalPeriod) error {
	return nil
}
func (m *mockFiscalPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return nil
}
func (m *mockFiscalPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return false, nil // Open
}

// --- UNIT TESTS ---

func TestCreateInvoice_Success(t *testing.T) {
	billingRepo := newMockBillingRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}
	uow := &mockUoW{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	cashProc := payment.NewCashProcessor("1-1111")
	bankProc := payment.NewBankTransferProcessor("1-1113")
	billingSvc := billingApp.NewService(uow, billingRepo, accSvc, cashProc, bankProc)

	ctx := context.Background()
	branchID := uuid.New()

	cmd := billingApp.CreateInvoiceCommand{
		BranchID:        branchID,
		PatientName:     "Budi Santoso",
		PatientNIK:      "3201012345678901",
		InvoiceType:     billingDomain.InvoiceTypeCombined,
		RegistrationFee: decimal.FromInt(25000),
		ConsultationFee: decimal.FromInt(50000),
		MedicationFee:   decimal.FromInt(75000),
	}

	inv, err := billingSvc.CreateInvoice(ctx, cmd)
	if err != nil {
		t.Fatalf("create invoice failed: %v", err)
	}

	expectedTotal := decimal.FromInt(150000)
	if !inv.TotalAmount.Equal(expectedTotal) {
		t.Errorf("expected total %s, got %s", expectedTotal, inv.TotalAmount)
	}
	if inv.PaymentStatus != billingDomain.PaymentStatusUnpaid {
		t.Errorf("expected status UNPAID, got %s", inv.PaymentStatus)
	}
}

func TestPayInvoice_Cash_AutoJournalPosting(t *testing.T) {
	billingRepo := newMockBillingRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}
	uow := &mockUoW{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	cashProc := payment.NewCashProcessor("1-1111")
	bankProc := payment.NewBankTransferProcessor("1-1113")
	billingSvc := billingApp.NewService(uow, billingRepo, accSvc, cashProc, bankProc)

	ctx := context.Background()
	branchID := uuid.New()

	// 1. Create Invoice
	inv, err := billingSvc.CreateInvoice(ctx, billingApp.CreateInvoiceCommand{
		BranchID:        branchID,
		PatientName:     "Siti Rahma",
		PatientNIK:      "3201019876543210",
		InvoiceType:     billingDomain.InvoiceTypeCombined,
		RegistrationFee: decimal.FromInt(25000),
		ConsultationFee: decimal.FromInt(50000),
		MedicationFee:   decimal.FromInt(75000),
	})
	if err != nil {
		t.Fatalf("create invoice failed: %v", err)
	}

	// 2. Pay Invoice via Cash
	paidInv, journalEntry, err := billingSvc.PayInvoice(ctx, billingApp.PayInvoiceCommand{
		InvoiceID:     inv.ID,
		Amount:        decimal.FromInt(150000),
		PaymentMethod: "CASH",
		ReferenceNo:   "CASH-001",
	})
	if err != nil {
		t.Fatalf("pay invoice failed: %v", err)
	}

	// 3. Verify Invoice State
	if paidInv.PaymentStatus != billingDomain.PaymentStatusPaid {
		t.Errorf("expected invoice status PAID, got %s", paidInv.PaymentStatus)
	}
	if paidInv.JournalID == nil {
		t.Fatalf("expected invoice to have linked journal_id")
	}

	// 4. Verify Accounting Journal Entry (Double-Entry Balance & Accounts)
	if journalEntry == nil {
		t.Fatalf("expected journal entry to be generated")
	}
	if journalEntry.JournalType != accountingDomain.JournalTypeCashReceipt {
		t.Errorf("expected CASH_RECEIPT journal type, got %s", journalEntry.JournalType)
	}
	if journalEntry.Status != accountingDomain.JournalStatusPosted {
		t.Errorf("expected status POSTED, got %s", journalEntry.Status)
	}
	if !journalEntry.TotalDebit.Equal(decimal.FromInt(150000)) || !journalEntry.TotalCredit.Equal(decimal.FromInt(150000)) {
		t.Errorf("expected balanced 150000, got debit: %s, credit: %s", journalEntry.TotalDebit, journalEntry.TotalCredit)
	}

	// Verify line count (1 Debit Kas Kasir + 3 Credits Pendapatan)
	if len(journalEntry.Lines) != 4 {
		t.Fatalf("expected 4 journal lines, got %d", len(journalEntry.Lines))
	}
	if !journalEntry.Lines[0].Debit.Equal(decimal.FromInt(150000)) {
		t.Errorf("expected first line debit 150000, got %s", journalEntry.Lines[0].Debit)
	}
}

func TestPayInvoice_BankTransfer_Success(t *testing.T) {
	billingRepo := newMockBillingRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}
	uow := &mockUoW{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	cashProc := payment.NewCashProcessor("1-1111")
	bankProc := payment.NewBankTransferProcessor("1-1113")
	billingSvc := billingApp.NewService(uow, billingRepo, accSvc, cashProc, bankProc)

	ctx := context.Background()
	branchID := uuid.New()

	inv, _ := billingSvc.CreateInvoice(ctx, billingApp.CreateInvoiceCommand{
		BranchID:        branchID,
		PatientName:     "Ahmad Fauzi",
		PatientNIK:      "3201011122334455",
		InvoiceType:     billingDomain.InvoiceTypeConsultation,
		ConsultationFee: decimal.FromInt(100000),
	})

	paidInv, journalEntry, err := billingSvc.PayInvoice(ctx, billingApp.PayInvoiceCommand{
		InvoiceID:     inv.ID,
		Amount:        decimal.FromInt(100000),
		PaymentMethod: "BANK_TRANSFER",
		ReferenceNo:   "TRF-BCA-987654",
	})
	if err != nil {
		t.Fatalf("pay invoice via bank transfer failed: %v", err)
	}

	if paidInv.PaymentStatus != billingDomain.PaymentStatusPaid {
		t.Errorf("expected PAID, got %s", paidInv.PaymentStatus)
	}
	if journalEntry.JournalType != accountingDomain.JournalTypeBank {
		t.Errorf("expected BANK journal type, got %s", journalEntry.JournalType)
	}
	if !journalEntry.TotalDebit.Equal(decimal.FromInt(100000)) {
		t.Errorf("expected debit 100000, got %s", journalEntry.TotalDebit)
	}
}

func TestPayInvoice_OverpaymentRejected(t *testing.T) {
	billingRepo := newMockBillingRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}
	uow := &mockUoW{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	cashProc := payment.NewCashProcessor("1-1111")
	bankProc := payment.NewBankTransferProcessor("1-1113")
	billingSvc := billingApp.NewService(uow, billingRepo, accSvc, cashProc, bankProc)

	ctx := context.Background()
	branchID := uuid.New()

	inv, _ := billingSvc.CreateInvoice(ctx, billingApp.CreateInvoiceCommand{
		BranchID:        branchID,
		PatientName:     "Test Pasien",
		InvoiceType:     billingDomain.InvoiceTypeRegistration,
		RegistrationFee: decimal.FromInt(35000),
	})

	_, _, err := billingSvc.PayInvoice(ctx, billingApp.PayInvoiceCommand{
		InvoiceID:     inv.ID,
		Amount:        decimal.FromInt(50000), // Exceeds 35000!
		PaymentMethod: "CASH",
	})
	if err == nil {
		t.Fatalf("expected overpayment error, got nil")
	}
}

