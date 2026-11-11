package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-api/internal/domain/asset"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresAssetRepo struct {
	db *sql.DB
}

func NewPostgresAssetRepo(db *sql.DB) *PostgresAssetRepo {
	return &PostgresAssetRepo{db: db}
}

func (r *PostgresAssetRepo) CreateAsset(ctx context.Context, a *asset.FixedAsset) error {
	exec := GetExecutor(ctx, r.db)
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	query := `
		INSERT INTO fixed_assets (
			id, branch_id, cost_center_id, asset_code, asset_name,
			category, acquisition_date, acquisition_cost, salvage_value,
			useful_life_months, depreciation_method, asset_account_code,
			accum_account_code, expense_account_code, status,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12,
			$13, $14, $15,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		a.ID, a.BranchID, a.CostCenterID, a.AssetCode, a.AssetName,
		a.Category, a.AcquisitionDate, a.AcquisitionCost, a.SalvageValue,
		a.UsefulLifeMonths, a.DepreciationMethod, a.AssetAccountCode,
		a.AccumAccountCode, a.ExpenseAccountCode, a.Status,
	)
	if err != nil {
		return fmt.Errorf("insert fixed asset: %w", err)
	}
	return nil
}

func (r *PostgresAssetRepo) GetAssetByID(ctx context.Context, id uuid.UUID) (*asset.FixedAsset, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, asset_code, asset_name,
		       category, acquisition_date, acquisition_cost, salvage_value,
		       useful_life_months, depreciation_method, asset_account_code,
		       accum_account_code, expense_account_code, status,
		       created_at, updated_at
		FROM fixed_assets
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	return r.scanAsset(row)
}

func (r *PostgresAssetRepo) ListActiveAssets(ctx context.Context, branchID *uuid.UUID) ([]asset.FixedAsset, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, asset_code, asset_name,
		       category, acquisition_date, acquisition_cost, salvage_value,
		       useful_life_months, depreciation_method, asset_account_code,
		       accum_account_code, expense_account_code, status,
		       created_at, updated_at
		FROM fixed_assets
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND status = 'ACTIVE'
		ORDER BY asset_code ASC
	`
	rows, err := exec.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("list active assets: %w", err)
	}
	defer rows.Close()

	var assets []asset.FixedAsset
	for rows.Next() {
		var a asset.FixedAsset
		var costStr, salvageStr string
		err := rows.Scan(
			&a.ID, &a.BranchID, &a.CostCenterID, &a.AssetCode, &a.AssetName,
			&a.Category, &a.AcquisitionDate, &costStr, &salvageStr,
			&a.UsefulLifeMonths, &a.DepreciationMethod, &a.AssetAccountCode,
			&a.AccumAccountCode, &a.ExpenseAccountCode, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan fixed asset: %w", err)
		}
		a.AcquisitionCost, _ = decimal.NewFromString(costStr)
		a.SalvageValue, _ = decimal.NewFromString(salvageStr)
		assets = append(assets, a)
	}
	return assets, nil
}

func (r *PostgresAssetRepo) List(ctx context.Context, filter asset.AssetFilter) ([]asset.FixedAsset, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM fixed_assets WHERE 1=1`
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
	if filter.Category != "" {
		baseQuery += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND acquisition_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND acquisition_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(asset_code) LIKE $%d OR LOWER(asset_name) LIKE $%d OR LOWER(category) LIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, searchParam, searchParam, searchParam)
		argIdx += 3
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count fixed assets: %w", err)
	}

	allowedSorts := map[string]string{
		"asset_code":       "asset_code",
		"asset_name":       "asset_name",
		"category":         "category",
		"acquisition_date": "acquisition_date",
		"acquisition_cost": "acquisition_cost",
		"status":           "status",
		"created_at":       "created_at",
	}
	sortCol := "asset_code"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortDir = "DESC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, cost_center_id, asset_code, asset_name,
		       category, acquisition_date, acquisition_cost, salvage_value,
		       useful_life_months, depreciation_method, asset_account_code,
		       accum_account_code, expense_account_code, status,
		       created_at, updated_at
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
		return nil, 0, fmt.Errorf("list fixed assets: %w", err)
	}
	defer rows.Close()

	var assets []asset.FixedAsset
	for rows.Next() {
		var a asset.FixedAsset
		var costStr, salvageStr string
		err := rows.Scan(
			&a.ID, &a.BranchID, &a.CostCenterID, &a.AssetCode, &a.AssetName,
			&a.Category, &a.AcquisitionDate, &costStr, &salvageStr,
			&a.UsefulLifeMonths, &a.DepreciationMethod, &a.AssetAccountCode,
			&a.AccumAccountCode, &a.ExpenseAccountCode, &a.Status,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan fixed asset: %w", err)
		}
		a.AcquisitionCost, _ = decimal.NewFromString(costStr)
		a.SalvageValue, _ = decimal.NewFromString(salvageStr)
		assets = append(assets, a)
	}
	return assets, totalCount, nil
}

func (r *PostgresAssetRepo) RecordDepreciation(ctx context.Context, dep *asset.AssetDepreciation) error {
	exec := GetExecutor(ctx, r.db)
	if dep.ID == uuid.Nil {
		dep.ID = uuid.New()
	}

	query := `
		INSERT INTO fixed_asset_depreciations (
			id, asset_id, period_year, period_month,
			depreciation_amount, accumulated_amount, book_value,
			journal_id, created_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		dep.ID, dep.AssetID, dep.PeriodYear, dep.PeriodMonth,
		dep.DepreciationAmount, dep.AccumulatedAmount, dep.BookValue,
		dep.JournalID,
	)
	if err != nil {
		return fmt.Errorf("insert fixed asset depreciation: %w", err)
	}
	return nil
}

func (r *PostgresAssetRepo) GetDepreciationHistory(ctx context.Context, assetID uuid.UUID) ([]asset.AssetDepreciation, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, asset_id, period_year, period_month,
		       depreciation_amount, accumulated_amount, book_value,
		       journal_id, created_at
		FROM fixed_asset_depreciations
		WHERE asset_id = $1
		ORDER BY period_year DESC, period_month DESC
	`
	rows, err := exec.QueryContext(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("list depreciation history: %w", err)
	}
	defer rows.Close()

	var history []asset.AssetDepreciation
	for rows.Next() {
		var d asset.AssetDepreciation
		var depStr, accumStr, bookStr string
		err := rows.Scan(
			&d.ID, &d.AssetID, &d.PeriodYear, &d.PeriodMonth,
			&depStr, &accumStr, &bookStr,
			&d.JournalID, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan depreciation history: %w", err)
		}
		d.DepreciationAmount, _ = decimal.NewFromString(depStr)
		d.AccumulatedAmount, _ = decimal.NewFromString(accumStr)
		d.BookValue, _ = decimal.NewFromString(bookStr)
		history = append(history, d)
	}
	return history, nil
}

func (r *PostgresAssetRepo) scanAsset(row *sql.Row) (*asset.FixedAsset, error) {
	var a asset.FixedAsset
	var costStr, salvageStr string
	err := row.Scan(
		&a.ID, &a.BranchID, &a.CostCenterID, &a.AssetCode, &a.AssetName,
		&a.Category, &a.AcquisitionDate, &costStr, &salvageStr,
		&a.UsefulLifeMonths, &a.DepreciationMethod, &a.AssetAccountCode,
		&a.AccumAccountCode, &a.ExpenseAccountCode, &a.Status,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("fixed asset not found")
		}
		return nil, fmt.Errorf("scan fixed asset: %w", err)
	}
	a.AcquisitionCost, _ = decimal.NewFromString(costStr)
	a.SalvageValue, _ = decimal.NewFromString(salvageStr)
	return &a, nil
}

