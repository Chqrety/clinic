package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/domain/accounting"

	"github.com/google/uuid"
)

type PostgresFiscalPeriodRepo struct {
	db *sql.DB
}

func NewPostgresFiscalPeriodRepo(db *sql.DB) *PostgresFiscalPeriodRepo {
	return &PostgresFiscalPeriodRepo{db: db}
}

func (r *PostgresFiscalPeriodRepo) GetPeriod(ctx context.Context, branchID uuid.UUID, year, month int) (*accounting.FiscalPeriod, error) {
	query := `
		SELECT id, branch_id, year, month, status, closed_at, closed_by, notes, created_at, updated_at
		FROM fiscal_periods
		WHERE branch_id = $1 AND year = $2 AND month = $3
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, branchID, year, month)

	var p accounting.FiscalPeriod
	var closedAt *time.Time
	var closedBy *uuid.UUID
	var notes sql.NullString

	err := row.Scan(
		&p.ID, &p.BranchID, &p.Year, &p.Month, &p.Status,
		&closedAt, &closedBy, &notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Period not yet registered -> default to OPEN
			return &accounting.FiscalPeriod{
				BranchID: branchID,
				Year:     year,
				Month:    month,
				Status:   "OPEN",
			}, nil
		}
		return nil, fmt.Errorf("get fiscal period: %w", err)
	}

	p.ClosedAt = closedAt
	p.ClosedBy = closedBy
	if notes.Valid {
		p.Notes = notes.String
	}

	return &p, nil
}

func (r *PostgresFiscalPeriodRepo) ClosePeriod(ctx context.Context, period *accounting.FiscalPeriod) error {
	query := `
		INSERT INTO fiscal_periods (id, branch_id, year, month, status, closed_at, closed_by, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'CLOSED', CURRENT_TIMESTAMP, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (branch_id, year, month)
		DO UPDATE SET status = 'CLOSED', closed_at = CURRENT_TIMESTAMP, closed_by = EXCLUDED.closed_by, notes = EXCLUDED.notes, updated_at = CURRENT_TIMESTAMP
	`
	if period.ID == uuid.Nil {
		period.ID = uuid.New()
	}

	exec := GetExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, query,
		period.ID, period.BranchID, period.Year, period.Month,
		period.ClosedBy, period.Notes,
	)
	if err != nil {
		return fmt.Errorf("close fiscal period: %w", err)
	}
	return nil
}

func (r *PostgresFiscalPeriodRepo) ReopenPeriod(ctx context.Context, periodID uuid.UUID, reason string) error {
	query := `
		UPDATE fiscal_periods
		SET status = 'OPEN', closed_at = NULL, closed_by = NULL, notes = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	exec := GetExecutor(ctx, r.db)
	res, err := exec.ExecContext(ctx, query, reason, periodID)
	if err != nil {
		return fmt.Errorf("reopen fiscal period: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("fiscal period not found")
	}
	return nil
}

func (r *PostgresFiscalPeriodRepo) IsPeriodLocked(ctx context.Context, branchID uuid.UUID, date time.Time) (bool, error) {
	year := date.Year()
	month := int(date.Month())

	query := `
		SELECT status
		FROM fiscal_periods
		WHERE branch_id = $1 AND year = $2 AND month = $3
	`
	exec := GetExecutor(ctx, r.db)
	var status string
	err := exec.QueryRowContext(ctx, query, branchID, year, month).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Open if no period entry exists
		}
		return false, fmt.Errorf("check period lock status: %w", err)
	}

	return status == "CLOSED" || status == "LOCKED", nil
}
