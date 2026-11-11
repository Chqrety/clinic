package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresTaxRepo struct {
	db *sql.DB
}

func NewPostgresTaxRepo(db *sql.DB) *PostgresTaxRepo {
	return &PostgresTaxRepo{db: db}
}

func (r *PostgresTaxRepo) CreateWithholding(ctx context.Context, tw *tax.TaxWithholding) error {
	exec := GetExecutor(ctx, r.db)
	if tw.ID == uuid.Nil {
		tw.ID = uuid.New()
	}

	query := `
		INSERT INTO tax_withholdings (
			id, branch_id, tax_type, bupot_number, taxpayer_nik,
			taxpayer_npwp, taxpayer_name, object_code, gross_amount,
			dpp_amount, tax_rate, tax_amount, period_year,
			period_month, is_reported, journal_id, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		tw.ID, tw.BranchID, tw.TaxType, tw.BupotNumber, tw.TaxpayerNIK,
		tw.TaxpayerNPWP, tw.TaxpayerName, tw.ObjectCode, tw.GrossAmount,
		tw.DPPAmount, tw.TaxRate, tw.TaxAmount, tw.PeriodYear,
		tw.PeriodMonth, tw.IsReported, tw.JournalID,
	)
	if err != nil {
		return fmt.Errorf("insert tax withholding: %w", err)
	}
	return nil
}

func (r *PostgresTaxRepo) GetWithholdingByID(ctx context.Context, id uuid.UUID) (*tax.TaxWithholding, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, tax_type, bupot_number, taxpayer_nik,
		       taxpayer_npwp, taxpayer_name, object_code, gross_amount,
		       dpp_amount, tax_rate, tax_amount, period_year,
		       period_month, is_reported, journal_id, created_at
		FROM tax_withholdings
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	var tw tax.TaxWithholding
	var grossStr, dppStr, rateStr, taxStr string
	err := row.Scan(
		&tw.ID, &tw.BranchID, &tw.TaxType, &tw.BupotNumber, &tw.TaxpayerNIK,
		&tw.TaxpayerNPWP, &tw.TaxpayerName, &tw.ObjectCode, &grossStr,
		&dppStr, &rateStr, &taxStr, &tw.PeriodYear,
		&tw.PeriodMonth, &tw.IsReported, &tw.JournalID, &tw.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("tax withholding not found")
		}
		return nil, fmt.Errorf("query tax withholding: %w", err)
	}
	tw.GrossAmount, _ = decimal.NewFromString(grossStr)
	tw.DPPAmount, _ = decimal.NewFromString(dppStr)
	tw.TaxRate, _ = decimal.NewFromString(rateStr)
	tw.TaxAmount, _ = decimal.NewFromString(taxStr)
	return &tw, nil
}

func (r *PostgresTaxRepo) GetCumulativeDPP(ctx context.Context, branchID uuid.UUID, taxpayerNIK string, year int) (decimal.Decimal, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT COALESCE(SUM(dpp_amount), 0)
		FROM tax_withholdings
		WHERE branch_id = $1 AND taxpayer_nik = $2 AND period_year = $3
	`
	var sumStr string
	err := exec.QueryRowContext(ctx, query, branchID, taxpayerNIK, year).Scan(&sumStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("get cumulative dpp: %w", err)
	}
	d, _ := decimal.NewFromString(sumStr)
	return d, nil
}

func (r *PostgresTaxRepo) ListWithholdings(ctx context.Context, branchID *uuid.UUID, year, month int, taxType *tax.TaxType) ([]tax.TaxWithholding, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, tax_type, bupot_number, taxpayer_nik,
		       taxpayer_npwp, taxpayer_name, object_code, gross_amount,
		       dpp_amount, tax_rate, tax_amount, period_year,
		       period_month, is_reported, journal_id, created_at
		FROM tax_withholdings
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND ($2 = 0 OR period_year = $2)
		  AND ($3 = 0 OR period_month = $3)
		  AND ($4::varchar IS NULL OR tax_type = $4)
		ORDER BY created_at DESC
	`
	var taxTypeStr *string
	if taxType != nil {
		s := string(*taxType)
		taxTypeStr = &s
	}

	rows, err := exec.QueryContext(ctx, query, branchID, year, month, taxTypeStr)
	if err != nil {
		return nil, fmt.Errorf("list tax withholdings: %w", err)
	}
	defer rows.Close()

	var results []tax.TaxWithholding
	for rows.Next() {
		var tw tax.TaxWithholding
		var grossStr, dppStr, rateStr, taxStr string
		err := rows.Scan(
			&tw.ID, &tw.BranchID, &tw.TaxType, &tw.BupotNumber, &tw.TaxpayerNIK,
			&tw.TaxpayerNPWP, &tw.TaxpayerName, &tw.ObjectCode, &grossStr,
			&dppStr, &rateStr, &taxStr, &tw.PeriodYear,
			&tw.PeriodMonth, &tw.IsReported, &tw.JournalID, &tw.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tax withholding: %w", err)
		}
		tw.GrossAmount, _ = decimal.NewFromString(grossStr)
		tw.DPPAmount, _ = decimal.NewFromString(dppStr)
		tw.TaxRate, _ = decimal.NewFromString(rateStr)
		tw.TaxAmount, _ = decimal.NewFromString(taxStr)
		results = append(results, tw)
	}
	return results, nil
}

func (r *PostgresTaxRepo) MarkAsReported(ctx context.Context, id uuid.UUID) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE tax_withholdings
		SET is_reported = TRUE
		WHERE id = $1
	`
	_, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark tax withholding as reported: %w", err)
	}
	return nil
}

