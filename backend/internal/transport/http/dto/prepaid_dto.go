package dto

import (
	"time"

	"clinic-api/internal/domain/prepaid"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type RegisterPrepaidRequest struct {
	BranchID           uuid.UUID               `json:"branch_id"`
	CostCenterID       uuid.UUID               `json:"cost_center_id"`
	Name               string                  `json:"name"`
	Category           prepaid.PrepaidCategory `json:"category"`
	TotalAmount        decimal.Decimal         `json:"total_amount"`
	PrepaidAccountCode string                  `json:"prepaid_account_code"`
	ExpenseAccountCode string                  `json:"expense_account_code"`
	StartDate          time.Time               `json:"start_date"`
	EndDate            time.Time               `json:"end_date"`
	PaymentAccountCode string                  `json:"payment_account_code,omitempty"`
	AutoPostPayment    bool                    `json:"auto_post_payment"`
}

type AmortizePrepaidRequest struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
}

type RunBatchAmortizationRequest struct {
	PeriodYear  int `json:"period_year"`
	PeriodMonth int `json:"period_month"`
}

