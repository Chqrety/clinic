package dto

import (
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type BudgetItemRequest struct {
	AccountCode    string          `json:"account_code"`
	PeriodMonth    int             `json:"period_month"`
	BudgetedAmount decimal.Decimal `json:"budgeted_amount"`
	Notes          string          `json:"notes,omitempty"`
}

type CreateBudgetRequest struct {
	BranchID     uuid.UUID           `json:"branch_id"`
	CostCenterID uuid.UUID           `json:"cost_center_id"`
	PeriodYear   int                 `json:"period_year"`
	Name         string              `json:"name"`
	Items        []BudgetItemRequest `json:"items"`
}

