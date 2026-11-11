package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/domain/treasury"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresTreasuryRepo struct {
	db *sql.DB
}

func NewPostgresTreasuryRepo(db *sql.DB) *PostgresTreasuryRepo {
	return &PostgresTreasuryRepo{db: db}
}

func (r *PostgresTreasuryRepo) OpenShift(ctx context.Context, shift *treasury.CashierShift) error {
	exec := GetExecutor(ctx, r.db)
	if shift.ID == uuid.Nil {
		shift.ID = uuid.New()
	}

	query := `
		INSERT INTO cashier_shifts (
			id, branch_id, cost_center_id, cashier_user_id, cashier_name,
			cash_account_code, start_time, opening_cash, system_cash_in,
			system_cash_out, expected_cash, actual_cash, difference_amount,
			difference_type, status, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13,
			$14, $15, $16, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		shift.ID, shift.BranchID, shift.CostCenterID, shift.CashierUserID, shift.CashierName,
		shift.CashAccountCode, shift.StartTime, shift.OpeningCash, shift.SystemCashIn,
		shift.SystemCashOut, shift.ExpectedCash, shift.ActualCash, shift.DifferenceAmount,
		shift.DifferenceType, shift.Status, shift.Notes,
	)
	if err != nil {
		return fmt.Errorf("open cashier shift: %w", err)
	}
	return nil
}

func (r *PostgresTreasuryRepo) GetActiveShift(ctx context.Context, branchID uuid.UUID, cashAccountCode string) (*treasury.CashierShift, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, cashier_user_id, cashier_name,
		       cash_account_code, start_time, end_time, opening_cash, system_cash_in,
		       system_cash_out, expected_cash, actual_cash, difference_amount,
		       difference_type, reconciliation_journal_id, deposit_journal_id,
		       status, notes, created_at, updated_at
		FROM cashier_shifts
		WHERE branch_id = $1 AND cash_account_code = $2 AND status = 'OPEN'
		ORDER BY start_time DESC
		LIMIT 1
	`
	row := exec.QueryRowContext(ctx, query, branchID, cashAccountCode)
	return r.scanShift(row)
}

func (r *PostgresTreasuryRepo) GetShiftByID(ctx context.Context, id uuid.UUID) (*treasury.CashierShift, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, cashier_user_id, cashier_name,
		       cash_account_code, start_time, end_time, opening_cash, system_cash_in,
		       system_cash_out, expected_cash, actual_cash, difference_amount,
		       difference_type, reconciliation_journal_id, deposit_journal_id,
		       status, notes, created_at, updated_at
		FROM cashier_shifts
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	shift, err := r.scanShift(row)
	if err != nil {
		return nil, err
	}

	// Load denominations
	denomsQuery := `
		SELECT id, shift_id, denomination_value, quantity, subtotal
		FROM cashier_shift_denominations
		WHERE shift_id = $1
		ORDER BY denomination_value DESC
	`
	rows, err := exec.QueryContext(ctx, denomsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var d treasury.ShiftDenomination
			var subtotalStr string
			if err := rows.Scan(&d.ID, &d.ShiftID, &d.DenominationValue, &d.Quantity, &subtotalStr); err == nil {
				d.Subtotal, _ = decimal.NewFromString(subtotalStr)
				shift.Denominations = append(shift.Denominations, d)
			}
		}
	}

	return shift, nil
}

func (r *PostgresTreasuryRepo) CloseShift(ctx context.Context, shift *treasury.CashierShift, denoms []treasury.ShiftDenomination) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE cashier_shifts
		SET end_time = $1, system_cash_in = $2, system_cash_out = $3,
		    expected_cash = $4, actual_cash = $5, difference_amount = $6,
		    difference_type = $7, reconciliation_journal_id = $8,
		    status = $9, notes = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $11
	`
	_, err := exec.ExecContext(ctx, query,
		shift.EndTime, shift.SystemCashIn, shift.SystemCashOut,
		shift.ExpectedCash, shift.ActualCash, shift.DifferenceAmount,
		shift.DifferenceType, shift.ReconciliationJournalID,
		shift.Status, shift.Notes, shift.ID,
	)
	if err != nil {
		return fmt.Errorf("close cashier shift: %w", err)
	}

	// Insert denominations
	for _, d := range denoms {
		denomQuery := `
			INSERT INTO cashier_shift_denominations (
				id, shift_id, denomination_value, quantity, subtotal
			) VALUES (
				$1, $2, $3, $4, $5
			)
		`
		denomID := uuid.New()
		_, err := exec.ExecContext(ctx, denomQuery, denomID, shift.ID, d.DenominationValue, d.Quantity, d.Subtotal)
		if err != nil {
			return fmt.Errorf("insert shift denomination: %w", err)
		}
	}

	return nil
}

func (r *PostgresTreasuryRepo) RecordDeposit(ctx context.Context, shiftID uuid.UUID, depositJournalID uuid.UUID) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE cashier_shifts
		SET deposit_journal_id = $1, status = 'DEPOSITED', updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := exec.ExecContext(ctx, query, depositJournalID, shiftID)
	if err != nil {
		return fmt.Errorf("record deposit in shift: %w", err)
	}
	return nil
}

func (r *PostgresTreasuryRepo) ListShifts(ctx context.Context, branchID *uuid.UUID, startDate, endDate *time.Time) ([]treasury.CashierShift, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, cashier_user_id, cashier_name,
		       cash_account_code, start_time, end_time, opening_cash, system_cash_in,
		       system_cash_out, expected_cash, actual_cash, difference_amount,
		       difference_type, reconciliation_journal_id, deposit_journal_id,
		       status, notes, created_at, updated_at
		FROM cashier_shifts
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND ($2::timestamptz IS NULL OR start_time >= $2)
		  AND ($3::timestamptz IS NULL OR start_time <= $3)
		ORDER BY start_time DESC
	`
	rows, err := exec.QueryContext(ctx, query, branchID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("list shifts: %w", err)
	}
	defer rows.Close()

	var shifts []treasury.CashierShift
	for rows.Next() {
		shift, err := r.scanShiftRow(rows)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, *shift)
	}
	return shifts, nil
}

func (r *PostgresTreasuryRepo) scanShift(row *sql.Row) (*treasury.CashierShift, error) {
	var s treasury.CashierShift
	var openStr, sysInStr, sysOutStr, expStr, actStr, diffStr string
	err := row.Scan(
		&s.ID, &s.BranchID, &s.CostCenterID, &s.CashierUserID, &s.CashierName,
		&s.CashAccountCode, &s.StartTime, &s.EndTime, &openStr, &sysInStr,
		&sysOutStr, &expStr, &actStr, &diffStr,
		&s.DifferenceType, &s.ReconciliationJournalID, &s.DepositJournalID,
		&s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cashier shift not found")
		}
		return nil, fmt.Errorf("scan shift: %w", err)
	}
	s.OpeningCash, _ = decimal.NewFromString(openStr)
	s.SystemCashIn, _ = decimal.NewFromString(sysInStr)
	s.SystemCashOut, _ = decimal.NewFromString(sysOutStr)
	s.ExpectedCash, _ = decimal.NewFromString(expStr)
	s.ActualCash, _ = decimal.NewFromString(actStr)
	s.DifferenceAmount, _ = decimal.NewFromString(diffStr)
	return &s, nil
}

func (r *PostgresTreasuryRepo) scanShiftRow(rows *sql.Rows) (*treasury.CashierShift, error) {
	var s treasury.CashierShift
	var openStr, sysInStr, sysOutStr, expStr, actStr, diffStr string
	err := rows.Scan(
		&s.ID, &s.BranchID, &s.CostCenterID, &s.CashierUserID, &s.CashierName,
		&s.CashAccountCode, &s.StartTime, &s.EndTime, &openStr, &sysInStr,
		&sysOutStr, &expStr, &actStr, &diffStr,
		&s.DifferenceType, &s.ReconciliationJournalID, &s.DepositJournalID,
		&s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan shift row: %w", err)
	}
	s.OpeningCash, _ = decimal.NewFromString(openStr)
	s.SystemCashIn, _ = decimal.NewFromString(sysInStr)
	s.SystemCashOut, _ = decimal.NewFromString(sysOutStr)
	s.ExpectedCash, _ = decimal.NewFromString(expStr)
	s.ActualCash, _ = decimal.NewFromString(actStr)
	s.DifferenceAmount, _ = decimal.NewFromString(diffStr)
	return &s, nil
}
