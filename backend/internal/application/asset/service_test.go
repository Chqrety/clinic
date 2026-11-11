package asset_test

import (
	"context"
	"testing"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	assetApp "clinic-api/internal/application/asset"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/asset"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockUoW struct{}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockAssetRepo struct {
	assets        map[uuid.UUID]*asset.FixedAsset
	depreciations map[uuid.UUID][]asset.AssetDepreciation
}

func newMockAssetRepo() *mockAssetRepo {
	return &mockAssetRepo{
		assets:        make(map[uuid.UUID]*asset.FixedAsset),
		depreciations: make(map[uuid.UUID][]asset.AssetDepreciation),
	}
}

func (m *mockAssetRepo) CreateAsset(ctx context.Context, a *asset.FixedAsset) error {
	m.assets[a.ID] = a
	return nil
}

func (m *mockAssetRepo) GetAssetByID(ctx context.Context, id uuid.UUID) (*asset.FixedAsset, error) {
	if a, ok := m.assets[id]; ok {
		return a, nil
	}
	return nil, assetApp.ErrAssetNotFound
}

func (m *mockAssetRepo) ListActiveAssets(ctx context.Context, branchID *uuid.UUID) ([]asset.FixedAsset, error) {
	var list []asset.FixedAsset
	for _, a := range m.assets {
		if a.Status == asset.AssetStatusActive {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockAssetRepo) List(ctx context.Context, filter asset.AssetFilter) ([]asset.FixedAsset, int64, error) {
	var list []asset.FixedAsset
	for _, a := range m.assets {
		list = append(list, *a)
	}
	return list, int64(len(list)), nil
}

func (m *mockAssetRepo) RecordDepreciation(ctx context.Context, dep *asset.AssetDepreciation) error {
	m.depreciations[dep.AssetID] = append(m.depreciations[dep.AssetID], *dep)
	return nil
}

func (m *mockAssetRepo) GetDepreciationHistory(ctx context.Context, assetID uuid.UUID) ([]asset.AssetDepreciation, error) {
	return m.depreciations[assetID], nil
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
		{"1-2200", "Peralatan Estetika SSC", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceDebit},
		{"1-2900", "Akumulasi Penyusutan Aset Tetap", accountingDomain.AccountTypeAsset, accountingDomain.NormalBalanceCredit},
		{"6-1300", "Beban Penyusutan Aset Tetap", accountingDomain.AccountTypeExpense, accountingDomain.NormalBalanceDebit},
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

func TestFixedAsset_RegisterAndMonthlyDepreciation(t *testing.T) {
	uow := &mockUoW{}
	assetRepo := newMockAssetRepo()
	accRepo := newMockAccountRepo()
	jRepo := &mockJournalRepo{journals: make(map[uuid.UUID]*accountingDomain.JournalEntry)}
	pRepo := &mockFiscalPeriodRepo{}

	accSvc := accountingApp.NewService(uow, accRepo, jRepo, pRepo)
	svc := assetApp.NewService(uow, assetRepo, accSvc)

	ctx := context.Background()
	branchID := uuid.New()
	costCenterID := uuid.New()

	// 1. Register Aesthetic Machine for Semarang Skin Center (SSC):
	// Mesin Laser CO2: Cost = 480,000,000, Salvage = 0, Useful Life = 48 months (4 years)
	// Monthly depreciation = 480M / 48 = 10,000,000 / month
	fixedAsset, _, err := svc.RegisterAsset(ctx, assetApp.RegisterAssetCommand{
		BranchID:         branchID,
		CostCenterID:     costCenterID,
		AssetCode:        "AST-SSC-LASER-001",
		AssetName:        "Mesin Laser CO2 Fractional SSC",
		Category:         "AESTHETIC_MACHINE",
		AcquisitionDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		AcquisitionCost:  decimal.FromInt(480000000),
		SalvageValue:     decimal.Zero,
		UsefulLifeMonths: 48,
	})
	if err != nil {
		t.Fatalf("failed to register asset: %v", err)
	}

	if fixedAsset.AssetAccountCode != "1-2200" {
		t.Errorf("expected asset account 1-2200, got %s", fixedAsset.AssetAccountCode)
	}

	// 2. Run monthly straight-line depreciation for month 1 (2026-01)
	dep, journal, err := svc.DepreciateAsset(ctx, assetApp.DepreciateAssetCommand{
		AssetID:     fixedAsset.ID,
		PeriodYear:  2026,
		PeriodMonth: 1,
	})
	if err != nil {
		t.Fatalf("failed to depreciate asset: %v", err)
	}

	expectedDep := decimal.FromInt(10000000) // 10 Million
	expectedBookValue := decimal.FromInt(470000000)

	if !dep.DepreciationAmount.Equal(expectedDep) {
		t.Errorf("expected monthly dep %s, got %s", expectedDep, dep.DepreciationAmount)
	}
	if !dep.BookValue.Equal(expectedBookValue) {
		t.Errorf("expected book value %s, got %s", expectedBookValue, dep.BookValue)
	}

	// 3. Verify Adjusting Journal
	// Debit: 6-1300 Beban Penyusutan (10,000,000)
	// Credit: 1-2900 Akumulasi Penyusutan (10,000,000)
	if journal == nil {
		t.Fatalf("expected depreciation journal to be generated")
	}
	if !journal.TotalDebit.Equal(expectedDep) || !journal.TotalCredit.Equal(expectedDep) {
		t.Errorf("expected balanced 10M, got dr: %s, cr: %s", journal.TotalDebit, journal.TotalCredit)
	}
}

