package doctorshare_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	docApp "clinic-api/internal/application/doctorshare"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/doctorshare"
	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockDoctorShareRepo struct {
	rules       map[string]*doctorshare.DoctorShareRule
	settlements map[uuid.UUID]*doctorshare.DoctorFeeSettlement
}

func newMockDoctorShareRepo() *mockDoctorShareRepo {
	return &mockDoctorShareRepo{
		rules:       make(map[string]*doctorshare.DoctorShareRule),
		settlements: make(map[uuid.UUID]*doctorshare.DoctorFeeSettlement),
	}
}

func (m *mockDoctorShareRepo) CreateRule(ctx context.Context, rule *doctorshare.DoctorShareRule) error {
	m.rules[rule.DoctorName+rule.ServiceCategory] = rule
	return nil
}

func (m *mockDoctorShareRepo) GetRule(ctx context.Context, branchID uuid.UUID, doctorName, serviceCategory string) (*doctorshare.DoctorShareRule, error) {
	if r, ok := m.rules[doctorName+serviceCategory]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *mockDoctorShareRepo) ListRules(ctx context.Context, branchID *uuid.UUID) ([]doctorshare.DoctorShareRule, error) {
	var list []doctorshare.DoctorShareRule
	for _, r := range m.rules {
		list = append(list, *r)
	}
	return list, nil
}

func (m *mockDoctorShareRepo) CreateSettlement(ctx context.Context, s *doctorshare.DoctorFeeSettlement) error {
	m.settlements[s.ID] = s
	return nil
}

func (m *mockDoctorShareRepo) GetSettlementByID(ctx context.Context, id uuid.UUID) (*doctorshare.DoctorFeeSettlement, error) {
	if s, ok := m.settlements[id]; ok {
		return s, nil
	}
	return nil, docApp.ErrSettlementNotFound
}

func (m *mockDoctorShareRepo) ListSettlements(ctx context.Context, branchID *uuid.UUID, year, month int) ([]doctorshare.DoctorFeeSettlement, error) {
	var list []doctorshare.DoctorFeeSettlement
	for _, s := range m.settlements {
		list = append(list, *s)
	}
	return list, nil
}

func (m *mockDoctorShareRepo) UpdateSettlementStatus(ctx context.Context, id uuid.UUID, status doctorshare.SettlementStatus, paidAt *time.Time, journalID *uuid.UUID) error {
	if s, ok := m.settlements[id]; ok {
		s.Status = status
		s.PaidAt = paidAt
		s.JournalID = journalID
	}
	return nil
}

type mockTaxRepo struct {
	withholdings []tax.TaxWithholding
}

func (m *mockTaxRepo) CreateWithholding(ctx context.Context, tw *tax.TaxWithholding) error {
	m.withholdings = append(m.withholdings, *tw)
	return nil
}

func (m *mockTaxRepo) GetWithholdingByID(ctx context.Context, id uuid.UUID) (*tax.TaxWithholding, error) {
	for _, w := range m.withholdings {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, nil
}

func (m *mockTaxRepo) GetCumulativeDPP(ctx context.Context, branchID uuid.UUID, taxpayerNIK string, year int) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

func (m *mockTaxRepo) ListWithholdings(ctx context.Context, branchID *uuid.UUID, year, month int, taxType *tax.TaxType) ([]tax.TaxWithholding, error) {
	return m.withholdings, nil
}

func (m *mockTaxRepo) MarkAsReported(ctx context.Context, id uuid.UUID) error {
	return nil
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
		{"2-1210", "Hutang Jasa Medis Dokter Kidinest", accountingDomain.AccountTypeLiability, accountingDomain.NormalBalanceCredit},
		{"2-1501.02", "Hutang PPh 21 Dokter", accountingDomain.AccountTypeLiability, accountingDomain.NormalBalanceCredit},
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

func TestDoctorShareLifecycle_PMK168AndPaymentJournal(t *testing.T) {
	uow := &mockUoW{}
	docRepo := newMockDoctorShareRepo()
	taxRepo := &mockTaxRepo{}
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	svc := docApp.NewService(uow, docRepo, taxRepo, accSvc)

	ctx := context.Background()
	branchID := uuid.New()

	// 1. Generate Settlement for Dr. Spesialis Anak
	// Patient Fee = 10,000,000, Gross doctor share = 7,000,000 (70%)
	settlement, err := svc.GenerateSettlement(ctx, docApp.GenerateSettlementCommand{
		BranchID:    branchID,
		DoctorName:  "dr. Budi Sp.A",
		DoctorNPWP:  "09.123.456.7-001.000",
		TaxpayerNIK: "3374010101800001",
		PeriodYear:  2026,
		PeriodMonth: 11,
		Items: []docApp.FeeItemInput{
			{
				InvoiceID:          uuid.New(),
				ServiceDate:        time.Now(),
				PatientName:        "Ananda Kevin",
				ServiceDescription: "Konsultasi Spesialis Anak & Denver II",
				PatientFee:         decimal.FromInt(10000000),
				DoctorShareAmount:  decimal.FromInt(7000000),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to generate settlement: %v", err)
	}

	// 2. Verify PMK 168/2023 Tax Calculation:
	// Gross = 7,000,000
	// DPP = 50% * 7,000,000 = 3,500,000
	// Tax PPh 21 = 5% * 3,500,000 = 175,000
	// Net Fee = 7,000,000 - 175,000 = 6,825,000
	expectedGross := decimal.FromInt(7000000)
	expectedDPP := decimal.FromInt(3500000)
	expectedTax := decimal.FromInt(175000)
	expectedNet := decimal.FromInt(6825000)

	if !settlement.GrossFee.Equal(expectedGross) {
		t.Errorf("expected gross %s, got %s", expectedGross, settlement.GrossFee)
	}
	if !settlement.DPPAmount.Equal(expectedDPP) {
		t.Errorf("expected DPP %s, got %s", expectedDPP, settlement.DPPAmount)
	}
	if !settlement.TaxPPh21.Equal(expectedTax) {
		t.Errorf("expected tax %s, got %s", expectedTax, settlement.TaxPPh21)
	}
	if !settlement.NetFee.Equal(expectedNet) {
		t.Errorf("expected net %s, got %s", expectedNet, settlement.NetFee)
	}

	// 3. Approve Settlement
	appSettlement, err := svc.ApproveSettlement(ctx, settlement.ID)
	if err != nil {
		t.Fatalf("failed to approve settlement: %v", err)
	}
	if appSettlement.Status != doctorshare.StatusApproved {
		t.Errorf("expected status APPROVED, got %s", appSettlement.Status)
	}

	// 4. Pay Settlement via Bank BCA
	paidSettlement, journal, err := svc.PaySettlement(ctx, docApp.PaySettlementCommand{
		SettlementID: settlement.ID,
		BankAccCode:  "1-1121",
		UnitCode:     "KIDINEST",
	})
	if err != nil {
		t.Fatalf("failed to pay settlement: %v", err)
	}

	if paidSettlement.Status != doctorshare.StatusPaid {
		t.Errorf("expected status PAID, got %s", paidSettlement.Status)
	}

	// 5. Verify Journal Double-Entry Balance
	// Debit: 2-1210 Hutang Jasa Medis (7,000,000)
	// Credit: 1-1121 Bank BCA (6,825,000)
	// Credit: 2-1501.02 Hutang PPh 21 Dokter (175,000)
	if journal == nil {
		t.Fatalf("expected journal to be returned")
	}
	if !journal.TotalDebit.Equal(expectedGross) || !journal.TotalCredit.Equal(expectedGross) {
		t.Errorf("expected balanced %s, got dr: %s, cr: %s", expectedGross, journal.TotalDebit, journal.TotalCredit)
	}

	// 6. Verify e-Bupot Tax Withholding Record Created
	if len(taxRepo.withholdings) != 1 {
		t.Fatalf("expected 1 tax withholding recorded, got %d", len(taxRepo.withholdings))
	}
	tw := taxRepo.withholdings[0]
	if !tw.TaxAmount.Equal(expectedTax) {
		t.Errorf("expected tax withholding amount %s, got %s", expectedTax, tw.TaxAmount)
	}
}

