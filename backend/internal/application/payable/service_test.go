package payable_test

import (
	"context"
	"errors"
	"testing"
	"time"

	journalApp "clinic-api/internal/application/journal"
	payableApp "clinic-api/internal/application/payable"
	accountingDomain "clinic-api/internal/domain/accounting"
	payableDomain "clinic-api/internal/domain/payable"
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

type mockPayableRepo struct {
	payables map[uuid.UUID]*payableDomain.Payable
	payments []*payableDomain.Payment
}

func (m *mockPayableRepo) Create(ctx context.Context, p *payableDomain.Payable) error {
	m.payables[p.ID] = p
	return nil
}

func (m *mockPayableRepo) GetByID(ctx context.Context, id uuid.UUID) (*payableDomain.Payable, error) {
	if p, ok := m.payables[id]; ok {
		return p, nil
	}
	return nil, errors.New("payable not found")
}

func (m *mockPayableRepo) GetByNumber(ctx context.Context, number string) (*payableDomain.Payable, error) {
	for _, p := range m.payables {
		if p.PayableNumber == number {
			return p, nil
		}
	}
	return nil, errors.New("payable not found")
}

func (m *mockPayableRepo) List(ctx context.Context, filter payableDomain.Filter) ([]payableDomain.Payable, int64, error) {
	var list []payableDomain.Payable
	for _, p := range m.payables {
		list = append(list, *p)
	}
	return list, int64(len(list)), nil
}

func (m *mockPayableRepo) UpdateOutstanding(ctx context.Context, id uuid.UUID, outstanding decimal.Decimal, status payableDomain.Status) error {
	if p, ok := m.payables[id]; ok {
		p.OutstandingAmount = outstanding
		p.Status = status
		return nil
	}
	return errors.New("payable not found")
}

func (m *mockPayableRepo) AddPayment(ctx context.Context, p *payableDomain.Payment) error {
	m.payments = append(m.payments, p)
	return nil
}

func (m *mockPayableRepo) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*payableDomain.AgingReport, error) {
	return &payableDomain.AgingReport{
		AsOfDate:   asOfDate,
		GrandTotal: decimal.Zero,
	}, nil
}

func setupPayableService() (*payableApp.Service, *mockPayableRepo, *mockJournalRepo) {
	uow := &mockUoW{}
	accRepo := &mockAccountRepo{accounts: make(map[string]*accountingDomain.Account)}
	accRepo.accounts["2-1110"] = &accountingDomain.Account{
		ID:            uuid.New(),
		Code:          "2-1110",
		Name:          "Hutang Supplier Farmasi & PBF",
		Type:          accountingDomain.AccountTypeLiability,
		NormalBalance: accountingDomain.NormalBalanceCredit,
		IsActive:      true,
		IsHeader:      false,
	}
	accRepo.accounts["1-1121"] = &accountingDomain.Account{
		ID:            uuid.New(),
		Code:          "1-1121",
		Name:          "Bank BCA Operasional",
		Type:          accountingDomain.AccountTypeAsset,
		NormalBalance: accountingDomain.NormalBalanceDebit,
		IsActive:      true,
		IsHeader:      false,
	}

	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockPeriodRepo{}
	jSvc := journalApp.NewService(uow, accRepo, jRepo, pRepo)

	payRepo := &mockPayableRepo{payables: make(map[uuid.UUID]*payableDomain.Payable)}
	paySvc := payableApp.NewService(uow, payRepo, accRepo, jSvc)

	return paySvc, payRepo, jRepo
}

func TestPayableLifecycle(t *testing.T) {
	paySvc, _, jRepo := setupPayableService()
	ctx := context.Background()
	branchID := uuid.New()

	// 1. Record Payable for PBF Supplier (Rp 5.000.000)
	totalAmt := decimal.FromInt(5000000)
	p, err := paySvc.CreatePayable(ctx, payableApp.CreatePayableCommand{
		BranchScoped:     shared.BranchScoped{BranchID: branchID},
		VendorType:       payableDomain.VendorTypeSupplierMedicine,
		VendorName:       "PT Kimia Farma Trading & Distribution",
		InvoiceReference: "FP-KF-202611-098",
		TotalAmount:      totalAmt,
		IssueDate:        time.Now(),
		AccountCode:      "2-1110",
	})
	if err != nil {
		t.Fatalf("failed to create payable: %v", err)
	}

	if p.Status != payableDomain.StatusOpen {
		t.Errorf("expected OPEN status, got %s", p.Status)
	}

	// 2. Pay Partial Amount (Rp 2.000.000 via Bank BCA)
	payAmt1 := decimal.FromInt(2000000)
	payRecord1, err := paySvc.PayPayable(ctx, payableApp.PayPayableCommand{
		PayableID:           p.ID,
		Amount:              payAmt1,
		PaymentDate:         time.Now(),
		CashBankAccountCode: "1-1121",
	})
	if err != nil {
		t.Fatalf("partial payment failed: %v", err)
	}

	if payRecord1.JournalID == nil {
		t.Errorf("expected JKK journal to be generated")
	}

	updatedP, _ := paySvc.GetPayableByID(ctx, p.ID)
	if updatedP.Status != payableDomain.StatusPartial {
		t.Errorf("expected PARTIAL status, got %s", updatedP.Status)
	}
	expectedRemaining := decimal.FromInt(3000000)
	if !updatedP.OutstandingAmount.Equal(expectedRemaining) {
		t.Errorf("expected remaining %s, got %s", expectedRemaining, updatedP.OutstandingAmount)
	}

	// 3. Settle Remaining Balance (Rp 3.000.000) -> Status PAID
	_, err = paySvc.PayPayable(ctx, payableApp.PayPayableCommand{
		PayableID:           p.ID,
		Amount:              expectedRemaining,
		PaymentDate:         time.Now(),
		CashBankAccountCode: "1-1121",
	})
	if err != nil {
		t.Fatalf("final settlement failed: %v", err)
	}

	finalP, _ := paySvc.GetPayableByID(ctx, p.ID)
	if finalP.Status != payableDomain.StatusPaid {
		t.Errorf("expected PAID status, got %s", finalP.Status)
	}
	if !finalP.OutstandingAmount.IsZero() {
		t.Errorf("expected zero outstanding, got %s", finalP.OutstandingAmount)
	}

	// 4. Verify JKK Journals Posted
	if len(jRepo.journals) != 2 {
		t.Errorf("expected 2 JKK journals, got %d", len(jRepo.journals))
	}
}

