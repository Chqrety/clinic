package treasury_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	treasuryApp "clinic-api/internal/application/treasury"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/treasury"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockTreasuryRepo struct {
	shifts map[uuid.UUID]*treasury.CashierShift
}

func newMockTreasuryRepo() *mockTreasuryRepo {
	return &mockTreasuryRepo{
		shifts: make(map[uuid.UUID]*treasury.CashierShift),
	}
}

func (m *mockTreasuryRepo) OpenShift(ctx context.Context, shift *treasury.CashierShift) error {
	m.shifts[shift.ID] = shift
	return nil
}

func (m *mockTreasuryRepo) GetActiveShift(ctx context.Context, branchID uuid.UUID, cashAccountCode string) (*treasury.CashierShift, error) {
	for _, s := range m.shifts {
		if s.BranchID == branchID && s.CashAccountCode == cashAccountCode && s.Status == treasury.ShiftStatusOpen {
			return s, nil
		}
	}
	return nil, nil
}

func (m *mockTreasuryRepo) GetShiftByID(ctx context.Context, id uuid.UUID) (*treasury.CashierShift, error) {
	if s, ok := m.shifts[id]; ok {
		return s, nil
	}
	return nil, treasuryApp.ErrShiftNotFound
}

func (m *mockTreasuryRepo) CloseShift(ctx context.Context, shift *treasury.CashierShift, denoms []treasury.ShiftDenomination) error {
	m.shifts[shift.ID] = shift
	return nil
}

func (m *mockTreasuryRepo) RecordDeposit(ctx context.Context, shiftID uuid.UUID, depositJournalID uuid.UUID) error {
	if s, ok := m.shifts[shiftID]; ok {
		s.DepositJournalID = &depositJournalID
		s.Status = treasury.ShiftStatusDeposited
	}
	return nil
}

func (m *mockTreasuryRepo) ListShifts(ctx context.Context, branchID *uuid.UUID, startDate, endDate *time.Time) ([]treasury.CashierShift, error) {
	var list []treasury.CashierShift
	for _, s := range m.shifts {
		list = append(list, *s)
	}
	return list, nil
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
		{"1-1111", "Kas Kasir Kidinest", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-1121", "Bank BCA Operasional", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"4-9200", "Pendapatan Selisih Kas Lebih", accountingDomain.AccountTypeRevenue, accountingDomain.NormalBalanceCredit},
		{"6-1700", "Beban Selisih Kas Kurang", accountingDomain.AccountTypeExpense, accountingDomain.NormalBalanceDebit},
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

func TestCashierShift_OpenCloseWithOverageAndDeposit(t *testing.T) {
	uow := &mockUoW{}
	treasuryRepo := newMockTreasuryRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	svc := treasuryApp.NewService(uow, treasuryRepo, accSvc)

	ctx := context.Background()
	branchID := uuid.New()
	costCenterID := uuid.New()

	// 1. Open Shift with Rp 500,000 modal awal
	shift, err := svc.OpenShift(ctx, treasuryApp.OpenShiftCommand{
		BranchID:        branchID,
		CostCenterID:    costCenterID,
		CashierName:     "Kurniawan Kasir Kidinest",
		CashAccountCode: "1-1111",
		OpeningCash:     decimal.FromInt(500000),
	})
	if err != nil {
		t.Fatalf("failed to open shift: %v", err)
	}
	if shift.Status != treasury.ShiftStatusOpen {
		t.Errorf("expected shift status OPEN, got %s", shift.Status)
	}

	// 2. Close Shift with Physical Cash Count:
	// System: In = 1,500,000, Out = 200,000 -> Expected = 500,000 + 1,500,000 - 200,000 = 1,800,000
	// Physical Count: 19 lembar Rp 100,000 = 1,900,000 (Overage of Rp 100,000!)
	closedShift, recJournal, err := svc.CloseShift(ctx, treasuryApp.CloseShiftCommand{
		ShiftID:       shift.ID,
		SystemCashIn:  decimal.FromInt(1500000),
		SystemCashOut: decimal.FromInt(200000),
		Denominations: []treasuryApp.DenominationInput{
			{DenominationValue: 100000, Quantity: 19},
		},
		Notes: "Rekonsiliasi shift sore",
	})
	if err != nil {
		t.Fatalf("failed to close shift: %v", err)
	}

	if closedShift.Status != treasury.ShiftStatusClosed {
		t.Errorf("expected shift status CLOSED, got %s", closedShift.Status)
	}
	if closedShift.DifferenceType != treasury.DifferenceOverage {
		t.Errorf("expected difference OVERAGE, got %s", closedShift.DifferenceType)
	}
	if !closedShift.DifferenceAmount.Equal(decimal.FromInt(100000)) {
		t.Errorf("expected difference amount 100000, got %s", closedShift.DifferenceAmount)
	}

	// Verify Adjusting Journal for Overage (Debit 1-1111, Credit 4-9200)
	if recJournal == nil {
		t.Fatalf("expected reconciliation journal to be generated")
	}
	if !recJournal.TotalDebit.Equal(decimal.FromInt(100000)) || !recJournal.TotalCredit.Equal(decimal.FromInt(100000)) {
		t.Errorf("expected balanced 100000, got dr: %s, cr: %s", recJournal.TotalDebit, recJournal.TotalCredit)
	}

	// 3. Deposit physical cash to Bank BCA (Rp 1,500,000)
	depositedShift, depJournal, err := svc.DepositCashToBank(ctx, treasuryApp.DepositCommand{
		ShiftID:        closedShift.ID,
		DepositAmount:  decimal.FromInt(1500000),
		BankAccCode:    "1-1121",
		DepositSlipRef: "SETOR-BCA-20261111-01",
	})
	if err != nil {
		t.Fatalf("failed to deposit cash to bank: %v", err)
	}

	if depositedShift.Status != treasury.ShiftStatusDeposited {
		t.Errorf("expected status DEPOSITED, got %s", depositedShift.Status)
	}
	if depJournal == nil {
		t.Fatalf("expected deposit journal to be posted")
	}
	// Verify Deposit Journal (Debit 1-1121 Bank BCA, Credit 1-1111 Kas Kasir)
	if !depJournal.TotalDebit.Equal(decimal.FromInt(1500000)) || !depJournal.TotalCredit.Equal(decimal.FromInt(1500000)) {
		t.Errorf("expected balanced deposit journal 1500000, got dr: %s, cr: %s", depJournal.TotalDebit, depJournal.TotalCredit)
	}
}

