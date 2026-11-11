package asset

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type AssetStatus string

const (
	AssetStatusActive           AssetStatus = "ACTIVE"
	AssetStatusFullyDepreciated AssetStatus = "FULLY_DEPRECIATED"
	AssetStatusDisposed         AssetStatus = "DISPOSED"
)

type FixedAsset struct {
	shared.BaseEntity
	shared.BranchScoped
	CostCenterID       uuid.UUID       `json:"cost_center_id"`
	AssetCode          string          `json:"asset_code"`
	AssetName          string          `json:"asset_name"`
	Category           string          `json:"category"` // 'MEDICAL_EQUIPMENT', 'AESTHETIC_MACHINE', 'RENOVATION', 'OFFICE_IT'
	AcquisitionDate    time.Time       `json:"acquisition_date"`
	AcquisitionCost    decimal.Decimal `json:"acquisition_cost"`
	SalvageValue       decimal.Decimal `json:"salvage_value"`
	UsefulLifeMonths   int             `json:"useful_life_months"`
	DepreciationMethod string          `json:"depreciation_method"` // 'STRAIGHT_LINE'
	AssetAccountCode   string          `json:"asset_account_code"`   // '1-2100' or '1-2200'
	AccumAccountCode   string          `json:"accum_account_code"`   // '1-2900'
	ExpenseAccountCode string          `json:"expense_account_code"` // '6-1300'
	Status             AssetStatus     `json:"status"`
}

type AssetDepreciation struct {
	ID                 uuid.UUID       `json:"id"`
	AssetID            uuid.UUID       `json:"asset_id"`
	PeriodYear         int             `json:"period_year"`
	PeriodMonth        int             `json:"period_month"`
	DepreciationAmount decimal.Decimal `json:"depreciation_amount"`
	AccumulatedAmount  decimal.Decimal `json:"accumulated_amount"`
	BookValue          decimal.Decimal `json:"book_value"`
	JournalID          *uuid.UUID      `json:"journal_id,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
}

// CalculateStraightLineMonthly calculates monthly depreciation: (Cost - Salvage) / UsefulLifeMonths
func (a *FixedAsset) CalculateStraightLineMonthly() decimal.Decimal {
	if a.UsefulLifeMonths <= 0 {
		return decimal.Zero
	}
	depreciable := a.AcquisitionCost.Sub(a.SalvageValue)
	if !depreciable.IsPositive() {
		return decimal.Zero
	}
	return depreciable.Div(decimal.FromInt(int64(a.UsefulLifeMonths)))
}

type Repository interface {
	CreateAsset(ctx context.Context, asset *FixedAsset) error
	GetAssetByID(ctx context.Context, id uuid.UUID) (*FixedAsset, error)
	ListActiveAssets(ctx context.Context, branchID *uuid.UUID) ([]FixedAsset, error)
	RecordDepreciation(ctx context.Context, dep *AssetDepreciation) error
	GetDepreciationHistory(ctx context.Context, assetID uuid.UUID) ([]AssetDepreciation, error)
}

