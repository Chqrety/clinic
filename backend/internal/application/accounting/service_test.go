package accounting_test

import (
	"context"
	"testing"
	"time"

	app "clinic-api/internal/application/accounting"
	domain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

// Mock UnitOfWork
type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// Mock AccountRepo
type mockAccountRepo struct {
	accounts map[uuid.UUID]*domain.Account
	byCode   map[string]*domain.Account
}

func newMockAccountRepo() *mockAccountRepo {
	r := &mockAccountRepo{
		accounts: make(map[uuid.UUID]*domain.Account),
		byCode:   make(map[string]*domain.Account),
	}
	// Seed sample accounts
	kasID := uuid.New()
	kas := &domain.Account{
		Code:          "1-1111",
		Name:          "Kas Kasir",
		Type:          domain.AccountTypeAsset,
		NormalBalance: domain.NormalBalanceDebit,
		IsActive:      true,
		IsHeader:      false,
	}
	kas.ID = kasID
	r.accounts[kasID] = kas
	r.byCode[kas.Code] = kas

	revID := uuid.New()
	rev := &domain.Account{
		Code:          "4-1100",
		Name:          "Pendapatan Konsultasi",
		Type:          domain.AccountTypeRevenue,
		NormalBalance: domain.NormalBalanceCredit,
		IsActive:      true,
		IsHeader:      false,
	}
	rev.ID = revID
	r.accounts[revID] = rev
	r.byCode[rev.Code] = rev

	return r
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, app.ErrAccountNotFound
	}
	return acc, nil
}

func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*domain.Account, error) {
	acc, ok := m.byCode[code]
	if !ok {
		return nil, app.ErrAccountNotFound
	}
	return acc, nil
}

func (m *mockAccountRepo) List(ctx context.Context, filter domain.AccountFilter) ([]domain.Account, int64, error) {
	var list []domain.Account
	for _, a := range m.accounts {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func (m *mockAccountRepo) Create(ctx context.Context, acc *domain.Account) error {
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}
	m.accounts[acc.ID] = acc
	m.byCode[acc.Code] = acc
	return nil
}

func (m *mockAccountRepo) Update(ctx context.Context, acc *domain.Account) error {
	m.accounts[acc.ID] = acc
	m.byCode[acc.Code] = acc
	return nil
}

// Mock JournalRepo
type mockJournalRepo struct {
	journals map[uuid.UUID]*domain.JournalEntry
}

func newMockJournalRepo() *mockJournalRepo {
	return &mockJournalRepo{
		journals: make(map[uuid.UUID]*domain.JournalEntry),
	}
}

func (m *mockJournalRepo) Create(ctx context.Context, entry *domain.JournalEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	m.journals[entry.ID] = entry
	return nil
}

func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.JournalEntry, error) {
	j, ok := m.journals[id]
	if !ok {
		return nil, app.ErrInvalidJournalStatus
	}
	return j, nil
}

func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*domain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.JournalNumber == number {
			return j, nil
		}
	}
	return nil, app.ErrInvalidJournalStatus
}

func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*domain.JournalEntry, error) {
	for _, j := range m.journals {
		if j.ReferenceType == refType && j.ReferenceID == refID {
			return j, nil
		}
	}
	return nil, app.ErrInvalidJournalStatus
}

func (m *mockJournalRepo) List(ctx context.Context, filter domain.JournalFilter) ([]domain.JournalEntry, int64, error) {
	var list []domain.JournalEntry
	for _, j := range m.journals {
		list = append(list, *j)
	}
	return list, int64(len(list)), nil
}

func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JournalStatus, postedBy *uuid.UUID) error {
	j, ok := m.journals[id]
	if !ok {
		return app.ErrInvalidJournalStatus
	}
	j.Status = status
	j.PostedBy = postedBy
	return nil
}

// Mock FiscalPeriodRepo
type mockFiscalPeriodRepo struct {
	isLocked bool
}

func (m *mockFiscalPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*domain.FiscalPeriod, error) {
	return &domain.FiscalPeriod{
		BranchID: branchID,
		Year:     year,
		Month:    month,
		Status:   "OPEN",
	}, nil
}

func (m *mockFiscalPeriodRepo) ClosePeriod(ctx context.Context, period *domain.FiscalPeriod) error {
	m.isLocked = true
	return nil
}

func (m *mockFiscalPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	m.isLocked = false
	return nil
}

func (m *mockFiscalPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return m.isLocked, nil
}

// --- UNIT TESTS ---

func TestCreateJournal_UnbalancedRejected(t *testing.T) {
	accRepo := newMockAccountRepo()
	jRepo := newMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{isLocked: false}
	svc := app.NewService(&mockUoW{}, accRepo, jRepo, pRepo)

	ctx := context.Background()
	branchID := uuid.New()

	cmd := app.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     domain.JournalTypeCashReceipt,
		TransactionDate: time.Now(),
		ReferenceType:   "INVOICE",
		ReferenceID:     "INV-001",
		Description:     "Pembayaran Konsultasi",
		Lines: []app.CreateLineCommand{
			{AccountCode: "1-1111", Debit: decimal.FromInt(50000), Credit: decimal.Zero},
			{AccountCode: "4-1100", Debit: decimal.Zero, Credit: decimal.FromInt(45000)}, // Unbalanced! 50k != 45k
		},
	}

	_, err := svc.CreateJournal(ctx, cmd)
	if err == nil {
		t.Fatalf("expected error for unbalanced journal, got nil")
	}
}

func TestCreateJournal_PeriodLockedRejected(t *testing.T) {
	accRepo := newMockAccountRepo()
	jRepo := newMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{isLocked: true} // Period locked!
	svc := app.NewService(&mockUoW{}, accRepo, jRepo, pRepo)

	ctx := context.Background()
	branchID := uuid.New()

	cmd := app.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     domain.JournalTypeCashReceipt,
		TransactionDate: time.Now(),
		ReferenceType:   "INVOICE",
		ReferenceID:     "INV-001",
		Description:     "Pembayaran Konsultasi",
		Lines: []app.CreateLineCommand{
			{AccountCode: "1-1111", Debit: decimal.FromInt(50000), Credit: decimal.Zero},
			{AccountCode: "4-1100", Debit: decimal.Zero, Credit: decimal.FromInt(50000)},
		},
	}

	_, err := svc.CreateJournal(ctx, cmd)
	if err == nil {
		t.Fatalf("expected error when fiscal period is locked, got nil")
	}
}

func TestCreateJournal_Success(t *testing.T) {
	accRepo := newMockAccountRepo()
	jRepo := newMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{isLocked: false}
	svc := app.NewService(&mockUoW{}, accRepo, jRepo, pRepo)

	ctx := context.Background()
	branchID := uuid.New()

	cmd := app.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     domain.JournalTypeCashReceipt,
		TransactionDate: time.Now(),
		ReferenceType:   "INVOICE",
		ReferenceID:     "INV-001",
		Description:     "Pembayaran Konsultasi",
		AutoPost:        false,
		Lines: []app.CreateLineCommand{
			{AccountCode: "1-1111", Debit: decimal.FromInt(50000), Credit: decimal.Zero},
			{AccountCode: "4-1100", Debit: decimal.Zero, Credit: decimal.FromInt(50000)},
		},
	}

	entry, err := svc.CreateJournal(ctx, cmd)
	if err != nil {
		t.Fatalf("expected journal creation to succeed, got %v", err)
	}

	if entry.Status != domain.JournalStatusDraft {
		t.Errorf("expected status DRAFT, got %s", entry.Status)
	}
	if !entry.TotalDebit.Equal(decimal.FromInt(50000)) {
		t.Errorf("expected total debit 50000, got %s", entry.TotalDebit)
	}
}

func TestPostAndReverseJournal_Success(t *testing.T) {
	accRepo := newMockAccountRepo()
	jRepo := newMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{isLocked: false}
	svc := app.NewService(&mockUoW{}, accRepo, jRepo, pRepo)

	ctx := context.Background()
	branchID := uuid.New()
	userID := uuid.New()

	// 1. Create Draft Journal
	cmd := app.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     domain.JournalTypeGeneral,
		TransactionDate: time.Now(),
		ReferenceType:   "MANUAL",
		ReferenceID:     "ADJ-001",
		Description:     "Manual Adjustment",
		AutoPost:        false,
		Lines: []app.CreateLineCommand{
			{AccountCode: "1-1111", Debit: decimal.FromInt(100000), Credit: decimal.Zero},
			{AccountCode: "4-1100", Debit: decimal.Zero, Credit: decimal.FromInt(100000)},
		},
	}
	entry, err := svc.CreateJournal(ctx, cmd)
	if err != nil {
		t.Fatalf("create journal failed: %v", err)
	}

	// 2. Post Journal
	posted, err := svc.PostJournal(ctx, entry.ID, &userID)
	if err != nil {
		t.Fatalf("post journal failed: %v", err)
	}
	if posted.Status != domain.JournalStatusPosted {
		t.Errorf("expected status POSTED, got %s", posted.Status)
	}

	// 3. Reverse Journal
	reversed, err := svc.ReverseJournal(ctx, entry.ID, "Koreksi salah input akun", &userID)
	if err != nil {
		t.Fatalf("reverse journal failed: %v", err)
	}
	if reversed.JournalType != domain.JournalTypeReversal {
		t.Errorf("expected type REVERSAL, got %s", reversed.JournalType)
	}
	// Verify lines were flipped
	if !reversed.Lines[0].Credit.Equal(decimal.FromInt(100000)) {
		t.Errorf("expected first line to be credit 100000, got %s", reversed.Lines[0].Credit)
	}
}

func TestCloseAndReopenPeriod_Success(t *testing.T) {
	accRepo := newMockAccountRepo()
	jRepo := newMockJournalRepo()
	pRepo := &mockFiscalPeriodRepo{isLocked: false}
	svc := app.NewService(&mockUoW{}, accRepo, jRepo, pRepo)

	ctx := context.Background()
	branchID := uuid.New()
	userID := uuid.New()

	// 1. Close and Lock Period
	period, err := svc.ClosePeriod(ctx, branchID, 2026, 11, &userID, "Tutup buku bulanan November 2026")
	if err != nil {
		t.Fatalf("close period failed: %v", err)
	}
	if period.Status != "LOCKED" {
		t.Errorf("expected LOCKED, got %s", period.Status)
	}

	// 2. Reopen Period
	err = svc.ReopenPeriod(ctx, period.ID, "Koreksi transaksi terlambat")
	if err != nil {
		t.Fatalf("reopen period failed: %v", err)
	}
}

