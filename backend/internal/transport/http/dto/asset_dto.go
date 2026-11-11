package dto

import (
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type RegisterAssetRequest struct {
	BranchID            uuid.UUID       `json:"branch_id"`
	CostCenterID        uuid.UUID       `json:"cost_center_id"`
	AssetCode           string          `json:"asset_code"`
	AssetName           string          `json:"asset_name"`
	Category            string          `json:"category"` // 'MEDICAL_EQUIPMENT', 'AESTHETIC_MACHINE', 'RENOVATION', 'OFFICE_IT'
	AcquisitionDate     time.Time       `json:"acquisition_date"`
	AcquisitionCost     decimal.Decimal `json:"acquisition_cost"`
	SalvageValue        decimal.Decimal `json:"salvage_value"`
	UsefulLifeMonths    int             `json:"useful_life_months"`
	CreditAccountCode   string          `json:"credit_account_code,omitempty"`
	AutoPostAcquisition bool            `json:"auto_post_acquisition"`
}

type DepreciateAssetRequest struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
}

type RunMonthlyDepreciationRequest struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
}

