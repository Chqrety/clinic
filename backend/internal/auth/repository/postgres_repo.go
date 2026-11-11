package repository

import (
	"context"
	"database/sql"
	"errors"

	"clinic-api/internal/auth/domain"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) FindByEmployeeNumber(ctx context.Context, empNum string) (*domain.StaffAuth, error) {
	query := `
	SELECT id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at
	FROM staff_users
	WHERE employee_number = $1
	LIMIT 1;
	`
	row := r.db.QueryRowContext(ctx, query, empNum)

	var s domain.StaffAuth
	err := row.Scan(
		&s.ID,
		&s.EmployeeNumber,
		&s.Name,
		&s.Role,
		&s.ClinicUnit,
		&s.PasswordHash,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	return &s, nil
}

func (r *PostgresAuthRepository) FindByID(ctx context.Context, id int64) (*domain.StaffAuth, error) {
	query := `
	SELECT id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at
	FROM staff_users
	WHERE id = $1
	LIMIT 1;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.StaffAuth
	err := row.Scan(
		&s.ID,
		&s.EmployeeNumber,
		&s.Name,
		&s.Role,
		&s.ClinicUnit,
		&s.PasswordHash,
		&s.IsActive,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrStaffNotFound
		}
		return nil, err
	}

	return &s, nil
}
