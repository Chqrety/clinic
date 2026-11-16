package prepaid_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	prepaidApp "clinic-api/internal/application/prepaid"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/prepaid"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockPrepaidRepo struct {
	items         map[uuid.UUID]*prepaid.PrepaidExpense
	amortizations map[uuid.UUID][]prepaid.PrepaidAmortization
}

func newMockPrepaidRepo() *mockPrepaidRepo {
	return &mockPrepaidRepo{
		items:         make(map[uuid.UUID]*prepaid.PrepaidExpense),
		amortizations: make(map[uuid.UUID][]prepaid.PrepaidAmortization),
	}
}

func (m *mockPrepaidRepo) CreatePrepaid(ctx context.Context, p *prepaid.PrepaidExpense) error {
	m.items[p.ID] = p
	return nil
}

func (m *mockPrepaidRepo) GetPrepaidByID(ctx context.Context, id uuid.UUID) (*prepaid.PrepaidExpense, error) {
	if p, ok := m.items[id]; ok {
		return p, nil
	}
	return nil, prepaidApp.ErrPrepaidNotFound
}

func (m *mockPrepaidRepo) ListActivePrepaid(ctx context.Context, branchID *uuid.UUID) ([]prepaid.PrepaidExpense, error) {
	var list []prepaid.PrepaidExpense
	for _, p := range m.items {
		if p.Status == prepaid.StatusActive {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *mockPrepaidRepo) List(ctx context.Context, filter prepaid.PrepaidFilter) ([]prepaid.PrepaidExpense, int64, error) {
	var list []prepaid.PrepaidExpense
	for _, p := range m.items {
		list = append(list, *p)
	}
	return list, int64(len(list)), nil
}

func (m *mockPrepaidRepo) RecordAmortization(ctx context.Context, a *prepaid.PrepaidAmortization, newRemaining decimal.Decimal, newStatus prepaid.PrepaidStatus) error {
	m.amortizations[a.PrepaidExpenseID] = append(m.amortizations[a.PrepaidExpenseID], *a)
	if p, ok := m.items[a.PrepaidExpenseID]; ok {
		p.RemainingAmount = newRemaining
		p.Status = newStatus
	}
	return nil
}

func (m *mockPrepaidRepo) GetAmortizationHistory(ctx context.Context, prepaidID uuid.UUID) ([]prepaid.PrepaidAmortization, error) {
	return m.amortizations[prepaidID], nil
}

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
		{"1-1121", "Bank BCA Operasional", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1400", "Biaya Dibayar di Muka", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"6-1200", "Beban Sewa Gedung & Pemeliharaan", accountingDomain.AccountTypeExpense, accountingDomain.NormalBalanceDebit},
	}
	for _, s := range seed {
		acc := &accountingDomain.Account{
			Code:          s.code,
			Name:          s.name,
			Type:          s.accType,
			NormalBalance: s.normal,
			IsActive:      true,
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
	return nil, nil
}

func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*accountingDomain.Account, error) {
	if a, ok := m.accounts[code]; ok {
		return a, nil
	}
	return nil, nil
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
	m.journals[entry.ID] = entry
	return nil
}

func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accountingDomain.JournalEntry, error) {
	return m.journals[id], nil
}

func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accountingDomain.JournalEntry, error) {
	return nil, nil
}

func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accountingDomain.JournalEntry, error) {
	return nil, nil
}

func (m *mockJournalRepo) List(ctx context.Context, filter accountingDomain.JournalFilter) ([]accountingDomain.JournalEntry, int64, error) {
	return nil, 0, nil
}

func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accountingDomain.JournalStatus, postedBy *uuid.UUID) error {
	return nil
}

type mockFiscalPeriodRepo struct{}

func (m *mockFiscalPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accountingDomain.FiscalPeriod, error) {
	return nil, nil
}
func (m *mockFiscalPeriodRepo) ClosePeriod(ctx context.Context, period *accountingDomain.FiscalPeriod) error {
	return nil
}
func (m *mockFiscalPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return nil
}
func (m *mockFiscalPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return false, nil
}

func TestPrepaidExpense_RegisterAndMonthlyAmortization(t *testing.T) {
	uow := &mockUoW{}
	repo := newMockPrepaidRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	svc := prepaidApp.NewService(uow, repo, accSvc)

	ctx := context.Background()
	branchID := uuid.New()
	costCenterID := uuid.New()

	// 1. Register Prepaid Rent: 120,000,000 for 12 months (Rp 10,000,000 / month)
	// Paid via Bank BCA 1-1121
	p, payJournal, err := svc.RegisterPrepaid(ctx, prepaidApp.RegisterPrepaidCommand{
		BranchID:           branchID,
		CostCenterID:       costCenterID,
		Name:               "Sewa Gedung Klinik Semarang 2026",
		Category:           prepaid.CategoryRent,
		TotalAmount:        decimal.FromInt(120000000),
		PrepaidAccountCode: "1-1400",
		ExpenseAccountCode: "6-1200",
		StartDate:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:            time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		PaymentAccountCode: "1-1121",
		AutoPostPayment:    true,
	})
	if err != nil {
		t.Fatalf("failed to register prepaid expense: %v", err)
	}

	if p.PeriodMonths != 12 {
		t.Errorf("expected 12 months, got %d", p.PeriodMonths)
	}
	expectedMonthly := decimal.FromInt(10000000)
	if !p.MonthlyAmount.Equal(expectedMonthly) {
		t.Errorf("expected monthly %s, got %s", expectedMonthly, p.MonthlyAmount)
	}

	// Verify Payment Journal (Debit 1-1400, Credit 1-1121)
	if payJournal == nil {
		t.Fatalf("expected payment journal to be generated")
	}
	if !payJournal.TotalDebit.Equal(decimal.FromInt(120000000)) {
		t.Errorf("expected 120M payment journal, got %s", payJournal.TotalDebit)
	}

	// 2. Amortize Month 1 (2026-01)
	amort, amortJournal, err := svc.AmortizeMonthly(ctx, prepaidApp.AmortizeCommand{
		PrepaidID:   p.ID,
		PeriodYear:  2026,
		PeriodMonth: 1,
	})
	if err != nil {
		t.Fatalf("failed to amortize: %v", err)
	}

	if !amort.Amount.Equal(expectedMonthly) {
		t.Errorf("expected amort amount %s, got %s", expectedMonthly, amort.Amount)
	}
	if amortJournal == nil {
		t.Fatalf("expected amortization journal to be generated")
	}
	// Verify Adjusting Journal: Debit 6-1200, Credit 1-1400
	if !amortJournal.TotalDebit.Equal(expectedMonthly) || !amortJournal.TotalCredit.Equal(expectedMonthly) {
		t.Errorf("expected balanced 10M, got dr: %s, cr: %s", amortJournal.TotalDebit, amortJournal.TotalCredit)
	}

	// Verify remaining balance (120M - 10M = 110M)
	updated, _ := svc.GetPrepaidByID(ctx, p.ID)
	expectedRemaining := decimal.FromInt(110000000)
	if !updated.RemainingAmount.Equal(expectedRemaining) {
		t.Errorf("expected remaining %s, got %s", expectedRemaining, updated.RemainingAmount)
	}
}

