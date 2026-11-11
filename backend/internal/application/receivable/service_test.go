package receivable_test

import (
	"context"
	"errors"
	"testing"
	"time"

	journalApp "clinic-api/internal/application/journal"
	receivableApp "clinic-api/internal/application/receivable"
	accountingDomain "clinic-api/internal/domain/accounting"
	receivableDomain "clinic-api/internal/domain/receivable"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockAccountRepo struct {
	accounts map[string]*accountingDomain.Account
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.Account, error) {
	for _, a := range m.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("account not found")
}

func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*accountingDomain.Account, error) {
	if a, ok := m.accounts[code]; ok {
		return a, nil
	}
	return nil, errors.New("account not found")
}

func (m *mockAccountRepo) List(ctx context.Context, filter accountingDomain.AccountFilter) ([]accountingDomain.Account, error) {
	var list []accountingDomain.Account
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, nil
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
	m.journals[entry.ID] = entry
	return nil
}

func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.JournalEntry, error) {
	if j, ok := m.journals[id]; ok {
		return j, nil
	}
	return nil, errors.New("journal not found")
}

func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.JournalNumber == number {
			return j, nil
		}
	}
	return nil, errors.New("journal not found")
}

func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accountingDomain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.ReferenceType == refType && j.ReferenceID == refID {
			return j, nil
		}
	}
	return nil, errors.New("journal not found")
}

func (m *mockJournalRepo) List(ctx context.Context, filter accountingDomain.JournalFilter) ([]accountingDomain.JournalEntry, int64, error) {
	var list []accountingDomain.JournalEntry
	for _, j := range m.journals {
		list = append(list, *j)
	}
	return list, int64(len(list)), nil
}

func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accountingDomain.JournalStatus, postedBy *uuid.UUID) error {
	if j, ok := m.journals[id]; ok {
		j.Status = status
		return nil
	}
	return errors.New("journal not found")
}

type mockPeriodRepo struct{}

func (m *mockPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accountingDomain.FiscalPeriod, error) {
	return &accountingDomain.FiscalPeriod{Status: "OPEN"}, nil
}
func (m *mockPeriodRepo) ClosePeriod(ctx context.Context, period *accountingDomain.FiscalPeriod) error {
	return nil
}
func (m *mockPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return nil
}
func (m *mockPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return false, nil
}

type mockReceivableRepo struct {
	receivables map[uuid.UUID]*receivableDomain.Receivable
	payments    []*receivableDomain.Payment
}

func (m *mockReceivableRepo) Create(ctx context.Context, r *receivableDomain.Receivable) error {
	m.receivables[r.ID] = r
	return nil
}

func (m *mockReceivableRepo) GetByID(ctx context.Context, id uuid.UUID) (*receivableDomain.Receivable, error) {
	if r, ok := m.receivables[id]; ok {
		return r, nil
	}
	return nil, errors.New("receivable not found")
}

func (m *mockReceivableRepo) GetByNumber(ctx context.Context, number string) (*receivableDomain.Receivable, error) {
	for _, r := range m.receivables {
		if r.ReceivableNumber == number {
			return r, nil
		}
	}
	return nil, errors.New("receivable not found")
}

func (m *mockReceivableRepo) List(ctx context.Context, filter receivableDomain.Filter) ([]receivableDomain.Receivable, int64, error) {
	var list []receivableDomain.Receivable
	for _, r := range m.receivables {
		list = append(list, *r)
	}
	return list, int64(len(list)), nil
}

func (m *mockReceivableRepo) UpdateOutstanding(ctx context.Context, id uuid.UUID, outstanding decimal.Decimal, status receivableDomain.Status) error {
	if r, ok := m.receivables[id]; ok {
		r.OutstandingAmount = outstanding
		r.Status = status
		return nil
	}
	return errors.New("receivable not found")
}

func (m *mockReceivableRepo) AddPayment(ctx context.Context, p *receivableDomain.Payment) error {
	m.payments = append(m.payments, p)
	return nil
}

func (m *mockReceivableRepo) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*receivableDomain.AgingReport, error) {
	return &receivableDomain.AgingReport{
		AsOfDate:   asOfDate,
		GrandTotal: decimal.Zero,
	}, nil
}

func setupReceivableService() (*receivableApp.Service, *mockReceivableRepo, *mockJournalRepo) {
	uow := &mockUoW{}
	accRepo := &mockAccountRepo{accounts: make(map[string]*accountingDomain.Account)}
	accRepo.accounts["1-1200"] = &accountingDomain.Account{
		ID:            uuid.New(),
		Code:          "1-1200",
		Name:          "Piutang Pasien & Klaim",
		Type:          accountingDomain.AccountTypeAsset,
		NormalBalance: accountingDomain.NormalBalanceDebit,
		IsActive:      true,
		IsHeader:      false,
	}
	accRepo.accounts["1-1111"] = &accountingDomain.Account{
		ID:            uuid.New(),
		Code:          "1-1111",
		Name:          "Kas Kasir Kidinest",
		Type:          accountingDomain.AccountTypeAsset,
		NormalBalance: accountingDomain.NormalBalanceDebit,
		IsActive:      true,
		IsHeader:      false,
	}

	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockPeriodRepo{}
	jSvc := journalApp.NewService(uow, accRepo, jRepo, pRepo)

	recRepo := &mockReceivableRepo{receivables: make(map[uuid.UUID]*receivableDomain.Receivable)}
	recSvc := receivableApp.NewService(uow, recRepo, accRepo, jSvc)

	return recSvc, recRepo, jRepo
}

func TestReceivableLifecycle(t *testing.T) {
	recSvc, _, jRepo := setupReceivableService()
	ctx := context.Background()
	branchID := uuid.New()

	// 1. Create Receivable for Insurance Claim (Rp 1.000.000)
	totalAmt := decimal.FromInt(1000000)
	rec, err := recSvc.CreateReceivable(ctx, receivableApp.CreateReceivableCommand{
		BranchScoped:     shared.BranchScoped{BranchID: branchID},
		CustomerType:     receivableDomain.CustomerTypeInsurance,
		CustomerName:     "Asuransi Mandiri Inhealth",
		InvoiceReference: "INV-202611-0001",
		TotalAmount:      totalAmt,
		IssueDate:        time.Now(),
		AccountCode:      "1-1200",
	})
	if err != nil {
		t.Fatalf("failed to create receivable: %v", err)
	}

	if rec.Status != receivableDomain.StatusOpen {
		t.Errorf("expected OPEN status, got %s", rec.Status)
	}
	if !rec.OutstandingAmount.Equal(totalAmt) {
		t.Errorf("expected outstanding %s, got %s", totalAmt, rec.OutstandingAmount)
	}

	// 2. Partial Payment (Rp 400.000)
	payAmt1 := decimal.FromInt(400000)
	payRecord1, err := recSvc.PayReceivable(ctx, receivableApp.PayReceivableCommand{
		ReceivableID:        rec.ID,
		Amount:              payAmt1,
		PaymentDate:         time.Now(),
		CashBankAccountCode: "1-1111",
	})
	if err != nil {
		t.Fatalf("partial payment failed: %v", err)
	}

	if payRecord1.JournalID == nil {
		t.Errorf("expected JKM journal to be generated")
	}

	// Verify updated receivable status
	updatedRec, _ := recSvc.GetReceivableByID(ctx, rec.ID)
	if updatedRec.Status != receivableDomain.StatusPartial {
		t.Errorf("expected PARTIAL status, got %s", updatedRec.Status)
	}
	expectedRemaining := decimal.FromInt(600000)
	if !updatedRec.OutstandingAmount.Equal(expectedRemaining) {
		t.Errorf("expected remaining %s, got %s", expectedRemaining, updatedRec.OutstandingAmount)
	}

	// 3. Attempt Payment Exceeding Outstanding -> Must Fail
	_, err = recSvc.PayReceivable(ctx, receivableApp.PayReceivableCommand{
		ReceivableID:        rec.ID,
		Amount:              decimal.FromInt(700000), // Exceeds 600k!
		PaymentDate:         time.Now(),
		CashBankAccountCode: "1-1111",
	})
	if !errors.Is(err, receivableApp.ErrInvalidPaymentAmount) {
		t.Errorf("expected ErrInvalidPaymentAmount, got %v", err)
	}

	// 4. Settle Remaining Balance (Rp 600.000) -> Status PAID
	_, err = recSvc.PayReceivable(ctx, receivableApp.PayReceivableCommand{
		ReceivableID:        rec.ID,
		Amount:              expectedRemaining,
		PaymentDate:         time.Now(),
		CashBankAccountCode: "1-1111",
	})
	if err != nil {
		t.Fatalf("final settlement failed: %v", err)
	}

	finalRec, _ := recSvc.GetReceivableByID(ctx, rec.ID)
	if finalRec.Status != receivableDomain.StatusPaid {
		t.Errorf("expected PAID status, got %s", finalRec.Status)
	}
	if !finalRec.OutstandingAmount.IsZero() {
		t.Errorf("expected zero outstanding, got %s", finalRec.OutstandingAmount)
	}

	// 5. Verify Journals Posted in General Ledger
	if len(jRepo.journals) != 2 {
		t.Errorf("expected 2 JKM journals, got %d", len(jRepo.journals))
	}
}

