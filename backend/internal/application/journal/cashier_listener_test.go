package journal_test

import (
	"context"
	"testing"
	"time"

	"clinic-api/internal/application/journal"
	cashierDomain "clinic-api/internal/cashier/domain"
	"clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockAccountRepo struct {
	accounts map[string]*accounting.Account
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accounting.Account, error) {
	for _, acc := range m.accounts {
		if acc.ID == id {
			return acc, nil
		}
	}
	return nil, journal.ErrAccountNotFound
}

func (m *mockAccountRepo) GetByCode(ctx context.Context, code string) (*accounting.Account, error) {
	acc, ok := m.accounts[code]
	if !ok {
		return nil, journal.ErrAccountNotFound
	}
	return acc, nil
}

func (m *mockAccountRepo) List(ctx context.Context, filter accounting.AccountFilter) ([]accounting.Account, int64, error) {
	return nil, 0, nil
}
func (m *mockAccountRepo) Create(ctx context.Context, acc *accounting.Account) error { return nil }
func (m *mockAccountRepo) Update(ctx context.Context, acc *accounting.Account) error { return nil }

type mockJournalRepo struct {
	saved []*accounting.JournalEntry
}

func (m *mockJournalRepo) Create(ctx context.Context, entry *accounting.JournalEntry) error {
	m.saved = append(m.saved, entry)
	return nil
}
func (m *mockJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accounting.JournalEntry, error) {
	return nil, nil
}
func (m *mockJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accounting.JournalEntry, error) {
	return nil, nil
}
func (m *mockJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accounting.JournalEntry, error) {
	return nil, nil
}
func (m *mockJournalRepo) List(ctx context.Context, filter accounting.JournalFilter) ([]accounting.JournalEntry, int64, error) {
	return nil, 0, nil
}
func (m *mockJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accounting.JournalStatus, postedBy *uuid.UUID) error {
	return nil
}

type mockPeriodRepo struct{}

func (m *mockPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accounting.FiscalPeriod, error) {
	return nil, nil
}
func (m *mockPeriodRepo) ClosePeriod(ctx context.Context, period *accounting.FiscalPeriod) error {
	return nil
}
func (m *mockPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	return nil
}
func (m *mockPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	return false, nil
}

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestCashierJournalListener_OnInvoiceSettled(t *testing.T) {
	accRepo := &mockAccountRepo{
		accounts: map[string]*accounting.Account{
			"1-1111": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "1-1111", Name: "Kas Kidinest", IsActive: true},
			"1-1112": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "1-1112", Name: "Kas SSC", IsActive: true},
			"1-1113": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "1-1113", Name: "Bank BCA", IsActive: true},
			"4-1100": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "4-1100", Name: "Pendapatan Registrasi", IsActive: true},
			"4-1200": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "4-1200", Name: "Pendapatan Jasa Dokter Sp.A", IsActive: true},
			"4-1500": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "4-1500", Name: "Pendapatan Obat Anak", IsActive: true},
			"4-2200": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "4-2200", Name: "Pendapatan Treatment SSC", IsActive: true},
			"4-2400": {BaseEntity: shared.BaseEntity{ID: uuid.New()}, Code: "4-2400", Name: "Pendapatan Skincare SSC", IsActive: true},
		},
	}
	jRepo := &mockJournalRepo{}
	pRepo := &mockPeriodRepo{}
	uow := &mockUoW{}

	svc := journal.NewService(uow, accRepo, jRepo, pRepo)
	listener := journal.NewCashierJournalListener(svc)

	inv := &cashierDomain.Invoice{
		InvoiceNumber:        "INV-20261117-001",
		ClinicUnit:           "KIDINEST",
		PatientName:          "Ananda Rayyan",
		PaymentMethod:        cashierDomain.PaymentMethodCash,
		SubtotalRegistration: 50000,
		SubtotalAction:       150000,
		SubtotalMedication:   100000,
		TotalAmount:          300000,
	}

	err := listener.OnInvoiceSettled(context.Background(), inv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(jRepo.saved) != 1 {
		t.Fatalf("expected 1 journal entry, got %d", len(jRepo.saved))
	}

	saved := jRepo.saved[0]
	if !saved.TotalDebit.Equal(decimal.FromInt(300000)) {
		t.Errorf("expected total debit 300000, got %s", saved.TotalDebit)
	}
	if !saved.TotalCredit.Equal(decimal.FromInt(300000)) {
		t.Errorf("expected total credit 300000, got %s", saved.TotalCredit)
	}
	if saved.ReferenceID != "INV-20261117-001" {
		t.Errorf("expected ref INV-20261117-001, got %s", saved.ReferenceID)
	}
}
