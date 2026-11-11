package asset

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/asset"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrAssetNotFound          = errors.New("fixed asset not found")
	ErrAssetAlreadyDepreciated = errors.New("asset is already fully depreciated")
	ErrInvalidUsefulLife      = errors.New("useful life must be greater than 0 months")
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	assetRepo         asset.Repository
	accountingService *accountingApp.Service
}

func NewService(
	uow UnitOfWork,
	assetRepo asset.Repository,
	accountingService *accountingApp.Service,
) *Service {
	return &Service{
		uow:               uow,
		assetRepo:         assetRepo,
		accountingService: accountingService,
	}
}

type RegisterAssetCommand struct {
	BranchID           uuid.UUID       `json:"branch_id"`
	CostCenterID       uuid.UUID       `json:"cost_center_id"`
	AssetCode          string          `json:"asset_code"`
	AssetName          string          `json:"asset_name"`
	Category           string          `json:"category"` // 'MEDICAL_EQUIPMENT', 'AESTHETIC_MACHINE', 'RENOVATION', 'OFFICE_IT'
	AcquisitionDate    time.Time       `json:"acquisition_date"`
	AcquisitionCost    decimal.Decimal `json:"acquisition_cost"`
	SalvageValue       decimal.Decimal `json:"salvage_value"`
	UsefulLifeMonths   int             `json:"useful_life_months"`
	CreditAccountCode  string          `json:"credit_account_code,omitempty"` // e.g. '1-1121' Bank BCA or '2-1110' Hutang
	AutoPostAcquisition bool           `json:"auto_post_acquisition"`
}

func (s *Service) RegisterAsset(ctx context.Context, cmd RegisterAssetCommand) (*asset.FixedAsset, *accountingDomain.JournalEntry, error) {
	if cmd.UsefulLifeMonths <= 0 {
		return nil, nil, ErrInvalidUsefulLife
	}

	assetAccountCode := "1-2100" // Peralatan Medis Poli Anak
	switch cmd.Category {
	case "AESTHETIC_MACHINE":
		assetAccountCode = "1-2200" // Peralatan Estetika SSC (Laser, Cauter)
	case "RENOVATION":
		assetAccountCode = "1-2300" // Renovasi Bangunan & Interior Klinik
	case "OFFICE_IT":
		assetAccountCode = "1-2100"
	}

	accumAccountCode := "1-2900"   // Akumulasi Penyusutan Aset Tetap
	expenseAccountCode := "6-1300" // Beban Penyusutan Aset Tetap

	fixedAsset := &asset.FixedAsset{
		BaseEntity:         shared.BaseEntity{ID: uuid.New()},
		BranchScoped:       shared.BranchScoped{BranchID: cmd.BranchID},
		CostCenterID:       cmd.CostCenterID,
		AssetCode:          cmd.AssetCode,
		AssetName:          cmd.AssetName,
		Category:           cmd.Category,
		AcquisitionDate:    cmd.AcquisitionDate,
		AcquisitionCost:    cmd.AcquisitionCost,
		SalvageValue:       cmd.SalvageValue,
		UsefulLifeMonths:   cmd.UsefulLifeMonths,
		DepreciationMethod: "STRAIGHT_LINE",
		AssetAccountCode:   assetAccountCode,
		AccumAccountCode:   accumAccountCode,
		ExpenseAccountCode: expenseAccountCode,
		Status:             asset.AssetStatusActive,
	}

	var journalEntry *accountingDomain.JournalEntry

	err := s.uow.Do(ctx, func(txCtx context.Context) error {
		if err := s.assetRepo.CreateAsset(txCtx, fixedAsset); err != nil {
			return err
		}

		if cmd.AutoPostAcquisition && cmd.CreditAccountCode != "" {
			jLines := []accountingApp.CreateLineCommand{
				{
					AccountCode: assetAccountCode,
					Debit:       cmd.AcquisitionCost,
					Credit:      decimal.Zero,
					Description: fmt.Sprintf("Perolehan Aset Tetap: %s (%s)", fixedAsset.AssetName, fixedAsset.AssetCode),
				},
				{
					AccountCode: cmd.CreditAccountCode,
					Debit:       decimal.Zero,
					Credit:      cmd.AcquisitionCost,
					Description: fmt.Sprintf("Pembayaran/Hutang Perolehan Aset: %s", fixedAsset.AssetName),
				},
			}

			jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
				BranchID:        cmd.BranchID,
				JournalType:     accountingDomain.JournalTypeGeneral,
				TransactionDate: cmd.AcquisitionDate,
				ReferenceType:   "FIXED_ASSET_ACQUISITION",
				ReferenceID:     fixedAsset.AssetCode,
				Description:     fmt.Sprintf("Kapitalisasi Aset Tetap %s", fixedAsset.AssetName),
				AutoPost:        true,
				Lines:           jLines,
			})
			if err != nil {
				return fmt.Errorf("post acquisition journal: %w", err)
			}
			journalEntry = jEntry
		}

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("register asset in uow: %w", err)
	}

	return fixedAsset, journalEntry, nil
}

type DepreciateAssetCommand struct {
	AssetID     uuid.UUID `json:"asset_id"`
	PeriodYear  int       `json:"period_year"`
	PeriodMonth int       `json:"period_month"`
}

func (s *Service) DepreciateAsset(ctx context.Context, cmd DepreciateAssetCommand) (*asset.AssetDepreciation, *accountingDomain.JournalEntry, error) {
	a, err := s.assetRepo.GetAssetByID(ctx, cmd.AssetID)
	if err != nil {
		return nil, nil, ErrAssetNotFound
	}
	if a.Status != asset.AssetStatusActive {
		return nil, nil, ErrAssetAlreadyDepreciated
	}

	monthlyDep := a.CalculateStraightLineMonthly()
	if !monthlyDep.IsPositive() {
		return nil, nil, errors.New("depreciation amount calculated is zero or negative")
	}

	history, err := s.assetRepo.GetDepreciationHistory(ctx, a.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("get depreciation history: %w", err)
	}

	currentAccum := decimal.Zero
	for _, h := range history {
		if h.PeriodYear == cmd.PeriodYear && h.PeriodMonth == cmd.PeriodMonth {
			return nil, nil, fmt.Errorf("asset has already been depreciated for period %d-%02d", cmd.PeriodYear, cmd.PeriodMonth)
		}
		currentAccum = currentAccum.Add(h.DepreciationAmount)
	}

	maxDepreciable := a.AcquisitionCost.Sub(a.SalvageValue)
	remainingDepreciable := maxDepreciable.Sub(currentAccum)

	actualDep := monthlyDep
	if actualDep.GreaterThan(remainingDepreciable) {
		actualDep = remainingDepreciable
	}

	newAccum := currentAccum.Add(actualDep)
	bookValue := a.AcquisitionCost.Sub(newAccum)

	now := time.Now()
	jLines := []accountingApp.CreateLineCommand{
		{
			AccountCode: a.ExpenseAccountCode, // 6-1300 Beban Penyusutan
			Debit:       actualDep,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Beban Penyusutan %s Periode %d-%02d", a.AssetName, cmd.PeriodYear, cmd.PeriodMonth),
		},
		{
			AccountCode: a.AccumAccountCode, // 1-2900 Akumulasi Penyusutan
			Debit:       decimal.Zero,
			Credit:      actualDep,
			Description: fmt.Sprintf("Akumulasi Penyusutan %s Periode %d-%02d", a.AssetName, cmd.PeriodYear, cmd.PeriodMonth),
		},
	}

	depRecord := &asset.AssetDepreciation{
		ID:                 uuid.New(),
		AssetID:            a.ID,
		PeriodYear:         cmd.PeriodYear,
		PeriodMonth:        cmd.PeriodMonth,
		DepreciationAmount: actualDep,
		AccumulatedAmount:  newAccum,
		BookValue:          bookValue,
		CreatedAt:          now,
	}

	var journalEntry *accountingDomain.JournalEntry

	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        a.BranchID,
			JournalType:     accountingDomain.JournalTypeAdjusting,
			TransactionDate: now,
			ReferenceType:   "FIXED_ASSET_DEPRECIATION",
			ReferenceID:     fmt.Sprintf("%s-%d%02d", a.AssetCode, cmd.PeriodYear, cmd.PeriodMonth),
			Description:     fmt.Sprintf("Penyusutan Aset %s (%d-%02d)", a.AssetName, cmd.PeriodYear, cmd.PeriodMonth),
			AutoPost:        true,
			Lines:           jLines,
		})
		if err != nil {
			return fmt.Errorf("create depreciation journal: %w", err)
		}
		journalEntry = jEntry
		depRecord.JournalID = &jEntry.ID

		return s.assetRepo.RecordDepreciation(txCtx, depRecord)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("record depreciation in uow: %w", err)
	}

	return depRecord, journalEntry, nil
}

func (s *Service) RunMonthlyDepreciationAll(ctx context.Context, branchID *uuid.UUID, year, month int) ([]asset.AssetDepreciation, error) {
	activeAssets, err := s.assetRepo.ListActiveAssets(ctx, branchID)
	if err != nil {
		return nil, fmt.Errorf("list active assets: %w", err)
	}

	var results []asset.AssetDepreciation
	for _, a := range activeAssets {
		dep, _, err := s.DepreciateAsset(ctx, DepreciateAssetCommand{
			AssetID:     a.ID,
			PeriodYear:  year,
			PeriodMonth: month,
		})
		if err == nil && dep != nil {
			results = append(results, *dep)
		}
	}
	return results, nil
}

func (s *Service) GetAssetByID(ctx context.Context, id uuid.UUID) (*asset.FixedAsset, error) {
	return s.assetRepo.GetAssetByID(ctx, id)
}

func (s *Service) ListAssets(ctx context.Context, filter asset.AssetFilter) ([]asset.FixedAsset, int64, error) {
	return s.assetRepo.List(ctx, filter)
}

func (s *Service) GetDepreciationHistory(ctx context.Context, assetID uuid.UUID) ([]asset.AssetDepreciation, error) {
	return s.assetRepo.GetDepreciationHistory(ctx, assetID)
}

