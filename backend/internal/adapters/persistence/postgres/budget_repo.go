package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-api/internal/domain/budget"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresBudgetRepo struct {
	db *sql.DB
}

func NewPostgresBudgetRepo(db *sql.DB) *PostgresBudgetRepo {
	return &PostgresBudgetRepo{db: db}
}

func (r *PostgresBudgetRepo) CreateBudget(ctx context.Context, b *budget.Budget) error {
	exec := GetExecutor(ctx, r.db)
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	query := `
		INSERT INTO budgets (
			id, branch_id, cost_center_id, period_year, name, status,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		b.ID, b.BranchID, b.CostCenterID, b.PeriodYear, b.Name, b.Status,
	)
	if err != nil {
		return fmt.Errorf("insert budget: %w", err)
	}

	for _, item := range b.Items {
		itemID := item.ID
		if itemID == uuid.Nil {
			itemID = uuid.New()
		}
		itemQuery := `
			INSERT INTO budget_items (
				id, budget_id, account_code, period_month, budgeted_amount, notes
			) VALUES (
				$1, $2, $3, $4, $5, $6
			)
		`
		_, err := exec.ExecContext(ctx, itemQuery,
			itemID, b.ID, item.AccountCode, item.PeriodMonth, item.BudgetedAmount, item.Notes,
		)
		if err != nil {
			return fmt.Errorf("insert budget item: %w", err)
		}
	}

	return nil
}

func (r *PostgresBudgetRepo) GetBudgetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, period_year, name, status, created_at, updated_at
		FROM budgets
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	var b budget.Budget
	err := row.Scan(
		&b.ID, &b.BranchID, &b.CostCenterID, &b.PeriodYear, &b.Name, &b.Status,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("budget not found")
		}
		return nil, fmt.Errorf("query budget: %w", err)
	}

	itemsQuery := `
		SELECT id, budget_id, account_code, period_month, budgeted_amount, notes
		FROM budget_items
		WHERE budget_id = $1
		ORDER BY account_code ASC, period_month ASC
	`
	rows, err := exec.QueryContext(ctx, itemsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it budget.BudgetItem
			var amountStr string
			var notes sql.NullString
			if err := rows.Scan(&it.ID, &it.BudgetID, &it.AccountCode, &it.PeriodMonth, &amountStr, &notes); err == nil {
				it.BudgetedAmount, _ = decimal.NewFromString(amountStr)
				if notes.Valid {
					it.Notes = notes.String
				}
				b.Items = append(b.Items, it)
			}
		}
	}

	return &b, nil
}

func (r *PostgresBudgetRepo) GetBudgetByYear(ctx context.Context, branchID uuid.UUID, costCenterID uuid.UUID, year int) (*budget.Budget, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, period_year, name, status, created_at, updated_at
		FROM budgets
		WHERE branch_id = $1 AND cost_center_id = $2 AND period_year = $3
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := exec.QueryRowContext(ctx, query, branchID, costCenterID, year)
	var b budget.Budget
	err := row.Scan(
		&b.ID, &b.BranchID, &b.CostCenterID, &b.PeriodYear, &b.Name, &b.Status,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("budget not found for specified year and cost center")
		}
		return nil, fmt.Errorf("query budget by year: %w", err)
	}

	return r.GetBudgetByID(ctx, b.ID)
}

func (r *PostgresBudgetRepo) UpdateBudgetStatus(ctx context.Context, id uuid.UUID, status budget.BudgetStatus) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE budgets
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := exec.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update budget status: %w", err)
	}
	return nil
}

func (r *PostgresBudgetRepo) ListBudgets(ctx context.Context, filter budget.BudgetFilter) ([]budget.Budget, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM budgets WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.CostCenterID != nil {
		baseQuery += fmt.Sprintf(" AND cost_center_id = $%d", argIdx)
		args = append(args, *filter.CostCenterID)
		argIdx++
	}
	if filter.Year > 0 {
		baseQuery += fmt.Sprintf(" AND period_year = $%d", argIdx)
		args = append(args, filter.Year)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, string(*filter.Status))
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND LOWER(name) LIKE $%d", argIdx)
		args = append(args, searchParam)
		argIdx++
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count budgets: %w", err)
	}

	allowedSorts := map[string]string{
		"period_year": "period_year",
		"name":        "name",
		"status":      "status",
		"created_at":  "created_at",
	}
	sortCol := "period_year"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, created_at DESC, id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, cost_center_id, period_year, name, status, created_at, updated_at
	` + baseQuery + orderBy

	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		selectQuery += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := exec.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list budgets: %w", err)
	}
	defer rows.Close()

	var list []budget.Budget
	for rows.Next() {
		var b budget.Budget
		err := rows.Scan(
			&b.ID, &b.BranchID, &b.CostCenterID, &b.PeriodYear, &b.Name, &b.Status,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan budget: %w", err)
		}
		list = append(list, b)
	}
	return list, totalCount, nil
}

func (r *PostgresBudgetRepo) GetActualAmounts(ctx context.Context, branchID *uuid.UUID, costCenterID *uuid.UUID, year, month int) (map[string]decimal.Decimal, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT a.code,
		       COALESCE(SUM(CASE 
		           WHEN a.type = 'REVENUE' THEN (jl.credit - jl.debit)
		           WHEN a.type IN ('COGS', 'EXPENSE') THEN (jl.debit - jl.credit)
		           ELSE (jl.debit - jl.credit)
		       END), 0) as actual_amount
		FROM journal_lines jl
		JOIN journal_entries je ON jl.journal_entry_id = je.id
		JOIN accounts a ON jl.account_id = a.id
		WHERE je.status = 'POSTED'
		  AND EXTRACT(YEAR FROM je.transaction_date) = $1
		  AND EXTRACT(MONTH FROM je.transaction_date) = $2
		  AND ($3::uuid IS NULL OR je.branch_id = $3)
		  AND ($4::uuid IS NULL OR jl.cost_center_id = $4)
		GROUP BY a.code
	`
	rows, err := exec.QueryContext(ctx, query, year, month, branchID, costCenterID)
	if err != nil {
		return nil, fmt.Errorf("query actual amounts for budget: %w", err)
	}
	defer rows.Close()

	results := make(map[string]decimal.Decimal)
	for rows.Next() {
		var code string
		var actualStr string
		if err := rows.Scan(&code, &actualStr); err == nil {
			d, _ := decimal.NewFromString(actualStr)
			results[code] = d
		}
	}
	return results, nil
}

