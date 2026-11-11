package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-api/internal/staff/domain"
)

type PostgresStaffRepository struct {
	db *sql.DB
}

func NewPostgresStaffRepository(db *sql.DB) *PostgresStaffRepository {
	return &PostgresStaffRepository{db: db}
}

func (r *PostgresStaffRepository) List(ctx context.Context, filter domain.StaffFilter) ([]domain.StaffUser, error) {
	query := `
	SELECT id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at
	FROM staff_users
	WHERE 1=1
	`
	var args []any
	argIdx := 1

	if filter.ClinicUnit != "" {
		query += fmt.Sprintf(" AND clinic_unit = $%d", argIdx)
		args = append(args, filter.ClinicUnit)
		argIdx++
	}

	if filter.Role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, filter.Role)
		argIdx++
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}

	query += " ORDER BY id ASC;"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.StaffUser
	for rows.Next() {
		var s domain.StaffUser
		if err := rows.Scan(
			&s.ID,
			&s.EmployeeNumber,
			&s.Name,
			&s.Role,
			&s.ClinicUnit,
			&s.PasswordHash,
			&s.IsActive,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	return list, rows.Err()
}

func (r *PostgresStaffRepository) GetByID(ctx context.Context, id int64) (*domain.StaffUser, error) {
	query := `
	SELECT id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at
	FROM staff_users
	WHERE id = $1
	LIMIT 1;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.StaffUser
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

func (r *PostgresStaffRepository) GetByEmployeeNumber(ctx context.Context, empNum string) (*domain.StaffUser, error) {
	query := `
	SELECT id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at
	FROM staff_users
	WHERE employee_number = $1
	LIMIT 1;
	`
	row := r.db.QueryRowContext(ctx, query, empNum)

	var s domain.StaffUser
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

func (r *PostgresStaffRepository) Create(ctx context.Context, staff *domain.StaffUser) error {
	query := `
	INSERT INTO staff_users (employee_number, name, role, clinic_unit, password_hash, is_active)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		staff.EmployeeNumber,
		staff.Name,
		staff.Role,
		staff.ClinicUnit,
		staff.PasswordHash,
		staff.IsActive,
	).Scan(&staff.ID, &staff.CreatedAt, &staff.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrDuplicateEmployeeNum
		}
		return err
	}

	return nil
}

func (r *PostgresStaffRepository) Update(ctx context.Context, id int64, input domain.UpdateStaffInput) (*domain.StaffUser, error) {
	query := `
	UPDATE staff_users
	SET name = $1, role = $2, clinic_unit = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
	RETURNING id, employee_number, name, role, clinic_unit, password_hash, is_active, created_at, updated_at;
	`
	row := r.db.QueryRowContext(ctx, query, input.Name, input.Role, input.ClinicUnit, input.IsActive, id)

	var s domain.StaffUser
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

func (r *PostgresStaffRepository) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	query := `
	UPDATE staff_users
	SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
	WHERE id = $2;
	`
	res, err := r.db.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrStaffNotFound
	}
	return nil
}
