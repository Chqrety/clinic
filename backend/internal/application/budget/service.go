package budget

import (
	"context"
	"errors"
	"fmt"

	"clinic-api/internal/domain/budget"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrBudgetNotFound        = errors.New("budget not found")
	ErrBudgetAlreadyApproved = errors.New("budget is already approved")
)

type Service struct {
	budgetRepo budget.Repository
}

func NewService(budgetRepo budget.Repository) *Service {
	return &Service{
		budgetRepo: budgetRepo,
	}
}

type BudgetItemInput struct {
	AccountCode    string          `json:"account_code"`
	PeriodMonth    int             `json:"period_month"` // 1 - 12
	BudgetedAmount decimal.Decimal `json:"budgeted_amount"`
	Notes          string          `json:"notes,omitempty"`
}

type CreateBudgetCommand struct {
	BranchID     uuid.UUID         `json:"branch_id"`
	CostCenterID uuid.UUID         `json:"cost_center_id"`
	PeriodYear   int               `json:"period_year"`
	Name         string            `json:"name"`
	Items        []BudgetItemInput `json:"items"`
}

func (s *Service) CreateBudget(ctx context.Context, cmd CreateBudgetCommand) (*budget.Budget, error) {
	b := &budget.Budget{
		BaseEntity:   shared.BaseEntity{ID: uuid.New()},
		BranchScoped: shared.BranchScoped{BranchID: cmd.BranchID},
		CostCenterID: cmd.CostCenterID,
		PeriodYear:   cmd.PeriodYear,
		Name:         cmd.Name,
		Status:       budget.StatusDraft,
	}

	for _, it := range cmd.Items {
		b.Items = append(b.Items, budget.BudgetItem{
			ID:             uuid.New(),
			BudgetID:       b.ID,
			AccountCode:    it.AccountCode,
			PeriodMonth:    it.PeriodMonth,
			BudgetedAmount: it.BudgetedAmount,
			Notes:          it.Notes,
		})
	}

	if err := s.budgetRepo.CreateBudget(ctx, b); err != nil {
		return nil, fmt.Errorf("create budget: %w", err)
	}

	return b, nil
}

func (s *Service) ApproveBudget(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	b, err := s.budgetRepo.GetBudgetByID(ctx, id)
	if err != nil {
		return nil, ErrBudgetNotFound
	}
	if b.Status == budget.StatusApproved {
		return nil, ErrBudgetAlreadyApproved
	}

	if err := s.budgetRepo.UpdateBudgetStatus(ctx, id, budget.StatusApproved); err != nil {
		return nil, fmt.Errorf("approve budget: %w", err)
	}
	b.Status = budget.StatusApproved
	return b, nil
}

func (s *Service) GetBudgetVarianceReport(ctx context.Context, branchID *uuid.UUID, costCenterID uuid.UUID, year, month int) (*budget.BudgetVarianceReport, error) {
	var bID uuid.UUID
	if branchID != nil {
		bID = *branchID
	}

	b, err := s.budgetRepo.GetBudgetByYear(ctx, bID, costCenterID, year)
	if err != nil {
		return nil, fmt.Errorf("lookup budget: %w", err)
	}

	actuals, err := s.budgetRepo.GetActualAmounts(ctx, branchID, &costCenterID, year, month)
	if err != nil {
		return nil, fmt.Errorf("get actual amounts: %w", err)
	}

	totalBudgeted := decimal.Zero
	totalActual := decimal.Zero
	var items []budget.BudgetVarianceItem

	for _, item := range b.Items {
		if item.PeriodMonth != month {
			continue
		}

		actual := actuals[item.AccountCode]
		variance := actual.Sub(item.BudgetedAmount)

		var pct decimal.Decimal
		if item.BudgetedAmount.IsPositive() {
			pct = variance.Div(item.BudgetedAmount).Mul(decimal.FromInt(100))
		}

		items = append(items, budget.BudgetVarianceItem{
			AccountCode:        item.AccountCode,
			AccountName:        item.AccountCode,
			PeriodMonth:        month,
			BudgetedAmount:     item.BudgetedAmount,
			ActualAmount:       actual,
			VarianceAmount:     variance,
			VariancePercentage: pct,
		})

		totalBudgeted = totalBudgeted.Add(item.BudgetedAmount)
		totalActual = totalActual.Add(actual)
	}

	totalVariance := totalActual.Sub(totalBudgeted)

	return &budget.BudgetVarianceReport{
		BranchID:      branchID,
		CostCenterID:  costCenterID,
		PeriodYear:    year,
		PeriodMonth:   month,
		TotalBudgeted: totalBudgeted,
		TotalActual:   totalActual,
		TotalVariance: totalVariance,
		Items:         items,
	}, nil
}

func (s *Service) GetBudgetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	return s.budgetRepo.GetBudgetByID(ctx, id)
}

func (s *Service) ListBudgets(ctx context.Context, branchID *uuid.UUID, year int) ([]budget.Budget, error) {
	return s.budgetRepo.ListBudgets(ctx, branchID, year)
}

