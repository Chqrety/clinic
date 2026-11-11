package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-api/internal/domain/prepaid"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresPrepaidRepo struct {
	db *sql.DB
}

func NewPostgresPrepaidRepo(db *sql.DB) *PostgresPrepaidRepo {
	return &PostgresPrepaidRepo{db: db}
}

func (r *PostgresPrepaidRepo) CreatePrepaid(ctx context.Context, p *prepaid.PrepaidExpense) error {
	exec := GetExecutor(ctx, r.db)
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	query := `
		INSERT INTO prepaid_expenses (
			id, branch_id, cost_center_id, name, category,
			total_amount, prepaid_account_code, expense_account_code,
			start_date, end_date, period_months, monthly_amount,
			remaining_amount, status, journal_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		p.ID, p.BranchID, p.CostCenterID, p.Name, p.Category,
		p.TotalAmount, p.PrepaidAccountCode, p.ExpenseAccountCode,
		p.StartDate, p.EndDate, p.PeriodMonths, p.MonthlyAmount,
		p.RemainingAmount, p.Status, p.JournalID,
	)
	if err != nil {
		return fmt.Errorf("insert prepaid expense: %w", err)
	}
	return nil
}

func (r *PostgresPrepaidRepo) GetPrepaidByID(ctx context.Context, id uuid.UUID) (*prepaid.PrepaidExpense, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, name, category,
		       total_amount, prepaid_account_code, expense_account_code,
		       start_date, end_date, period_months, monthly_amount,
		       remaining_amount, status, journal_id, created_at, updated_at
		FROM prepaid_expenses
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	var p prepaid.PrepaidExpense
	var totalStr, monthlyStr, remStr string
	err := row.Scan(
		&p.ID, &p.BranchID, &p.CostCenterID, &p.Name, &p.Category,
		&totalStr, &p.PrepaidAccountCode, &p.ExpenseAccountCode,
		&p.StartDate, &p.EndDate, &p.PeriodMonths, &monthlyStr,
		&remStr, &p.Status, &p.JournalID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("prepaid expense not found")
		}
		return nil, fmt.Errorf("query prepaid expense: %w", err)
	}
	p.TotalAmount, _ = decimal.NewFromString(totalStr)
	p.MonthlyAmount, _ = decimal.NewFromString(monthlyStr)
	p.RemainingAmount, _ = decimal.NewFromString(remStr)
	return &p, nil
}

func (r *PostgresPrepaidRepo) ListActivePrepaid(ctx context.Context, branchID *uuid.UUID) ([]prepaid.PrepaidExpense, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, name, category,
		       total_amount, prepaid_account_code, expense_account_code,
		       start_date, end_date, period_months, monthly_amount,
		       remaining_amount, status, journal_id, created_at, updated_at
		FROM prepaid_expenses
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND status = 'ACTIVE'
		ORDER BY start_date ASC
	`
	rows, err := exec.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("list active prepaid: %w", err)
	}
	defer rows.Close()

	var list []prepaid.PrepaidExpense
	for rows.Next() {
		var p prepaid.PrepaidExpense
		var totalStr, monthlyStr, remStr string
		err := rows.Scan(
			&p.ID, &p.BranchID, &p.CostCenterID, &p.Name, &p.Category,
			&totalStr, &p.PrepaidAccountCode, &p.ExpenseAccountCode,
			&p.StartDate, &p.EndDate, &p.PeriodMonths, &monthlyStr,
			&remStr, &p.Status, &p.JournalID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan prepaid expense: %w", err)
		}
		p.TotalAmount, _ = decimal.NewFromString(totalStr)
		p.MonthlyAmount, _ = decimal.NewFromString(monthlyStr)
		p.RemainingAmount, _ = decimal.NewFromString(remStr)
		list = append(list, p)
	}
	return list, nil
}

func (r *PostgresPrepaidRepo) List(ctx context.Context, filter prepaid.PrepaidFilter) ([]prepaid.PrepaidExpense, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM prepaid_expenses WHERE 1=1`
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
	if filter.Category != nil {
		baseQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, string(*filter.Category))
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, string(*filter.Status))
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND start_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND end_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(prepaid_account_code) LIKE $%d OR LOWER(expense_account_code) LIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, searchParam, searchParam, searchParam)
		argIdx += 3
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count prepaid expenses: %w", err)
	}

	allowedSorts := map[string]string{
		"name":             "name",
		"category":         "category",
		"total_amount":     "total_amount",
		"remaining_amount": "remaining_amount",
		"start_date":       "start_date",
		"end_date":         "end_date",
		"status":           "status",
		"created_at":       "created_at",
	}
	sortCol := "start_date"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, cost_center_id, name, category,
		       total_amount, prepaid_account_code, expense_account_code,
		       start_date, end_date, period_months, monthly_amount,
		       remaining_amount, status, journal_id, created_at, updated_at
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
		return nil, 0, fmt.Errorf("list prepaid expenses: %w", err)
	}
	defer rows.Close()

	var list []prepaid.PrepaidExpense
	for rows.Next() {
		var p prepaid.PrepaidExpense
		var totalStr, monthlyStr, remStr string
		err := rows.Scan(
			&p.ID, &p.BranchID, &p.CostCenterID, &p.Name, &p.Category,
			&totalStr, &p.PrepaidAccountCode, &p.ExpenseAccountCode,
			&p.StartDate, &p.EndDate, &p.PeriodMonths, &monthlyStr,
			&remStr, &p.Status, &p.JournalID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan prepaid expense: %w", err)
		}
		p.TotalAmount, _ = decimal.NewFromString(totalStr)
		p.MonthlyAmount, _ = decimal.NewFromString(monthlyStr)
		p.RemainingAmount, _ = decimal.NewFromString(remStr)
		list = append(list, p)
	}
	return list, totalCount, nil
}

func (r *PostgresPrepaidRepo) RecordAmortization(ctx context.Context, a *prepaid.PrepaidAmortization, newRemaining decimal.Decimal, newStatus prepaid.PrepaidStatus) error {
	exec := GetExecutor(ctx, r.db)
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	queryAmort := `
		INSERT INTO prepaid_amortizations (
			id, prepaid_expense_id, period_year, period_month,
			amount, journal_id, created_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, queryAmort,
		a.ID, a.PrepaidExpenseID, a.PeriodYear, a.PeriodMonth,
		a.Amount, a.JournalID,
	)
	if err != nil {
		return fmt.Errorf("insert prepaid amortization: %w", err)
	}

	queryUpdate := `
		UPDATE prepaid_expenses
		SET remaining_amount = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err = exec.ExecContext(ctx, queryUpdate, newRemaining, newStatus, a.PrepaidExpenseID)
	if err != nil {
		return fmt.Errorf("update remaining prepaid amount: %w", err)
	}

	return nil
}

func (r *PostgresPrepaidRepo) GetAmortizationHistory(ctx context.Context, prepaidID uuid.UUID) ([]prepaid.PrepaidAmortization, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, prepaid_expense_id, period_year, period_month,
		       amount, journal_id, created_at
		FROM prepaid_amortizations
		WHERE prepaid_expense_id = $1
		ORDER BY period_year DESC, period_month DESC
	`
	rows, err := exec.QueryContext(ctx, query, prepaidID)
	if err != nil {
		return nil, fmt.Errorf("list amortization history: %w", err)
	}
	defer rows.Close()

	var history []prepaid.PrepaidAmortization
	for rows.Next() {
		var a prepaid.PrepaidAmortization
		var amountStr string
		err := rows.Scan(
			&a.ID, &a.PrepaidExpenseID, &a.PeriodYear, &a.PeriodMonth,
			&amountStr, &a.JournalID, &a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan amortization: %w", err)
		}
		a.Amount, _ = decimal.NewFromString(amountStr)
		history = append(history, a)
	}
	return history, nil
}

