package budget_test

import (
	"context"
	"testing"

	budgetApp "clinic-api/internal/application/budget"
	"clinic-api/internal/domain/budget"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type mockBudgetRepo struct {
	budgets map[uuid.UUID]*budget.Budget
	actuals map[string]decimal.Decimal
}

func newMockBudgetRepo() *mockBudgetRepo {
	return &mockBudgetRepo{
		budgets: make(map[uuid.UUID]*budget.Budget),
		actuals: make(map[string]decimal.Decimal),
	}
}

func (m *mockBudgetRepo) CreateBudget(ctx context.Context, b *budget.Budget) error {
	m.budgets[b.ID] = b
	return nil
}

func (m *mockBudgetRepo) GetBudgetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	if b, ok := m.budgets[id]; ok {
		return b, nil
	}
	return nil, budgetApp.ErrBudgetNotFound
}

func (m *mockBudgetRepo) GetBudgetByYear(ctx context.Context, branchID uuid.UUID, costCenterID uuid.UUID, year int) (*budget.Budget, error) {
	for _, b := range m.budgets {
		if b.BranchID == branchID && b.CostCenterID == costCenterID && b.PeriodYear == year {
			return b, nil
		}
	}
	return nil, budgetApp.ErrBudgetNotFound
}

func (m *mockBudgetRepo) UpdateBudgetStatus(ctx context.Context, id uuid.UUID, status budget.BudgetStatus) error {
	if b, ok := m.budgets[id]; ok {
		b.Status = status
	}
	return nil
}

func (m *mockBudgetRepo) ListBudgets(ctx context.Context, filter budget.BudgetFilter) ([]budget.Budget, int64, error) {
	var list []budget.Budget
	for _, b := range m.budgets {
		list = append(list, *b)
	}
	return list, int64(len(list)), nil
}

func (m *mockBudgetRepo) GetActualAmounts(ctx context.Context, branchID *uuid.UUID, costCenterID *uuid.UUID, year, month int) (map[string]decimal.Decimal, error) {
	return m.actuals, nil
}

func TestBudget_CreateApproveAndVarianceReport(t *testing.T) {
	repo := newMockBudgetRepo()
	svc := budgetApp.NewService(repo)

	ctx := context.Background()
	branchID := uuid.New()
	costCenterID := uuid.New()

	// 1. Create Annual Budget for Kidinest Poli Anak
	// Budget Pendapatan 4-1120 Konsultasi: 50,000,000 / month
	// Budget Beban 6-1100 Gaji Perawat: 15,000,000 / month
	b, err := svc.CreateBudget(ctx, budgetApp.CreateBudgetCommand{
		BranchID:     branchID,
		CostCenterID: costCenterID,
		PeriodYear:   2026,
		Name:         "Rencana Anggaran Kidinest 2026",
		Items: []budgetApp.BudgetItemInput{
			{
				AccountCode:    "4-1120",
				PeriodMonth:    11,
				BudgetedAmount: decimal.FromInt(50000000),
				Notes:          "Target Pendapatan Sp.A",
			},
			{
				AccountCode:    "6-1100",
				PeriodMonth:    11,
				BudgetedAmount: decimal.FromInt(15000000),
				Notes:          "Plafon Gaji Tenaga Medis",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create budget: %v", err)
	}

	if b.Status != budget.StatusDraft {
		t.Errorf("expected DRAFT status, got %s", b.Status)
	}

	// 2. Approve Budget
	appBudget, err := svc.ApproveBudget(ctx, b.ID)
	if err != nil {
		t.Fatalf("failed to approve budget: %v", err)
	}
	if appBudget.Status != budget.StatusApproved {
		t.Errorf("expected APPROVED status, got %s", appBudget.Status)
	}

	// 3. Set Mock Actuals:
	// Actual Pendapatan 4-1120 = 60,000,000 (Favorable variance +10M)
	// Actual Beban 6-1100 = 14,000,000 (Favorable variance -1M)
	repo.actuals["4-1120"] = decimal.FromInt(60000000)
	repo.actuals["6-1100"] = decimal.FromInt(14000000)

	// 4. Generate Budget Variance Report for Month 11
	report, err := svc.GetBudgetVarianceReport(ctx, &branchID, costCenterID, 2026, 11)
	if err != nil {
		t.Fatalf("failed to get variance report: %v", err)
	}

	if len(report.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(report.Items))
	}

	expectedTotalBudget := decimal.FromInt(65000000)
	expectedTotalActual := decimal.FromInt(74000000)
	if !report.TotalBudgeted.Equal(expectedTotalBudget) {
		t.Errorf("expected total budget %s, got %s", expectedTotalBudget, report.TotalBudgeted)
	}
	if !report.TotalActual.Equal(expectedTotalActual) {
		t.Errorf("expected total actual %s, got %s", expectedTotalActual, report.TotalActual)
	}
}

