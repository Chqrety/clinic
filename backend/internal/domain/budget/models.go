package budget

import (
	"context"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type BudgetStatus string

const (
	StatusDraft    BudgetStatus = "DRAFT"
	StatusApproved BudgetStatus = "APPROVED"
	StatusRevised  BudgetStatus = "REVISED"
)

type BudgetItem struct {
	ID             uuid.UUID       `json:"id"`
	BudgetID       uuid.UUID       `json:"budget_id"`
	AccountCode    string          `json:"account_code"`
	PeriodMonth    int             `json:"period_month"` // 1 - 12
	BudgetedAmount decimal.Decimal `json:"budgeted_amount"`
	Notes          string          `json:"notes,omitempty"`
}

type Budget struct {
	shared.BaseEntity
	shared.BranchScoped
	CostCenterID uuid.UUID    `json:"cost_center_id"`
	PeriodYear   int          `json:"period_year"`
	Name         string       `json:"name"`
	Status       BudgetStatus `json:"status"`
	Items        []BudgetItem `json:"items,omitempty"`
}

type BudgetVarianceItem struct {
	AccountCode        string          `json:"account_code"`
	AccountName        string          `json:"account_name"`
	PeriodMonth        int             `json:"period_month"`
	BudgetedAmount     decimal.Decimal `json:"budgeted_amount"`
	ActualAmount       decimal.Decimal `json:"actual_amount"`
	VarianceAmount     decimal.Decimal `json:"variance_amount"` // Actual - Budget
	VariancePercentage decimal.Decimal `json:"variance_percentage"`
}

type BudgetVarianceReport struct {
	BranchID      *uuid.UUID           `json:"branch_id,omitempty"`
	CostCenterID  uuid.UUID            `json:"cost_center_id"`
	PeriodYear    int                  `json:"period_year"`
	PeriodMonth   int                  `json:"period_month"`
	TotalBudgeted decimal.Decimal      `json:"total_budgeted"`
	TotalActual   decimal.Decimal      `json:"total_actual"`
	TotalVariance decimal.Decimal      `json:"total_variance"`
	Items         []BudgetVarianceItem `json:"items"`
}

type BudgetFilter struct {
	BranchID     *uuid.UUID
	CostCenterID *uuid.UUID
	Year         int
	Status       *BudgetStatus
	Search       string
	SortBy       string
	SortOrder    string
	Limit        int
	Offset       int
}

type Repository interface {
	CreateBudget(ctx context.Context, b *Budget) error
	GetBudgetByID(ctx context.Context, id uuid.UUID) (*Budget, error)
	GetBudgetByYear(ctx context.Context, branchID uuid.UUID, costCenterID uuid.UUID, year int) (*Budget, error)
	UpdateBudgetStatus(ctx context.Context, id uuid.UUID, status BudgetStatus) error
	ListBudgets(ctx context.Context, filter BudgetFilter) ([]Budget, int64, error)
	GetActualAmounts(ctx context.Context, branchID *uuid.UUID, costCenterID *uuid.UUID, year, month int) (map[string]decimal.Decimal, error)
}
