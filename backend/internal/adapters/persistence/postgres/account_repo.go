package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"clinic-api/internal/domain/accounting"

	"github.com/google/uuid"
)

type PostgresAccountRepo struct {
	db *sql.DB
}

func NewPostgresAccountRepo(db *sql.DB) *PostgresAccountRepo {
	return &PostgresAccountRepo{db: db}
}

func (r *PostgresAccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*accounting.Account, error) {
	query := `
		SELECT id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)

	var acc accounting.Account
	var parentID *uuid.UUID
	var desc sql.NullString

	err := row.Scan(
		&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.NormalBalance,
		&parentID, &acc.Level, &acc.IsHeader, &acc.IsActive, &desc,
		&acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("query account by id: %w", err)
	}

	acc.ParentID = parentID
	if desc.Valid {
		acc.Description = desc.String
	}

	return &acc, nil
}

func (r *PostgresAccountRepo) GetByCode(ctx context.Context, code string) (*accounting.Account, error) {
	query := `
		SELECT id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description, created_at, updated_at
		FROM accounts
		WHERE code = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, code)

	var acc accounting.Account
	var parentID *uuid.UUID
	var desc sql.NullString

	err := row.Scan(
		&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.NormalBalance,
		&parentID, &acc.Level, &acc.IsHeader, &acc.IsActive, &desc,
		&acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account with code %s not found: %w", code, err)
		}
		return nil, fmt.Errorf("query account by code: %w", err)
	}

	acc.ParentID = parentID
	if desc.Valid {
		acc.Description = desc.String
	}

	return &acc, nil
}

func (r *PostgresAccountRepo) List(ctx context.Context, filter accounting.AccountFilter) ([]accounting.Account, error) {
	query := `
		SELECT id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description, created_at, updated_at
		FROM accounts
		WHERE 1=1
	`
	var args []any
	argIdx := 1

	if filter.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, *filter.Type)
		argIdx++
	}
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}
	if filter.IsHeader != nil {
		query += fmt.Sprintf(" AND is_header = $%d", argIdx)
		args = append(args, *filter.IsHeader)
		argIdx++
	}

	query += " ORDER BY code ASC"

	exec := GetExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []accounting.Account
	for rows.Next() {
		var acc accounting.Account
		var parentID *uuid.UUID
		var desc sql.NullString

		if err := rows.Scan(
			&acc.ID, &acc.Code, &acc.Name, &acc.Type, &acc.NormalBalance,
			&parentID, &acc.Level, &acc.IsHeader, &acc.IsActive, &desc,
			&acc.CreatedAt, &acc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan account row: %w", err)
		}
		acc.ParentID = parentID
		if desc.Valid {
			acc.Description = desc.String
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (r *PostgresAccountRepo) Create(ctx context.Context, acc *accounting.Account) error {
	query := `
		INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}

	exec := GetExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, query,
		acc.ID, acc.Code, acc.Name, acc.Type, acc.NormalBalance,
		acc.ParentID, acc.Level, acc.IsHeader, acc.IsActive, acc.Description,
	)
	if err != nil {
		return fmt.Errorf("insert account: %w", err)
	}
	return nil
}

func (r *PostgresAccountRepo) Update(ctx context.Context, acc *accounting.Account) error {
	query := `
		UPDATE accounts
		SET name = $1, is_active = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	exec := GetExecutor(ctx, r.db)
	res, err := exec.ExecContext(ctx, query, acc.Name, acc.IsActive, acc.Description, acc.ID)
	if err != nil {
		return fmt.Errorf("update account: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("account not found or no changes made")
	}
	return nil
}
