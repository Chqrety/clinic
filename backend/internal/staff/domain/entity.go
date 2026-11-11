package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleOwner      Role = "SUPERADMIN_OWNER"
	RoleDoctor     Role = "DOCTOR"
	RoleCashier    Role = "GRO_CASHIER"
	RolePharmacist Role = "PHARMACIST"
)

type ClinicUnit string

const (
	UnitKidinest ClinicUnit = "KIDINEST"
	UnitSSC      ClinicUnit = "SSC"
	UnitAll      ClinicUnit = "ALL"
)

type StaffUser struct {
	ID             int64      `json:"id"`
	EmployeeNumber string     `json:"employee_number"`
	Name           string     `json:"name"`
	Role           Role       `json:"role"`
	ClinicUnit     ClinicUnit `json:"clinic_unit"`
	PasswordHash   string     `json:"-"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateStaffInput struct {
	EmployeeNumber  string     `json:"employee_number"`
	Name            string     `json:"name"`
	Role            Role       `json:"role"`
	ClinicUnit      ClinicUnit `json:"clinic_unit"`
	InitialPassword string     `json:"initial_password"`
}

type UpdateStaffInput struct {
	Name       string     `json:"name"`
	Role       Role       `json:"role"`
	ClinicUnit ClinicUnit `json:"clinic_unit"`
	IsActive   bool       `json:"is_active"`
}

type ResetPasswordInput struct {
	NewPassword string `json:"new_password"`
}

type StaffFilter struct {
	ClinicUnit string
	Role       string
	IsActive   *bool
}

type StaffRepository interface {
	List(ctx context.Context, filter StaffFilter) ([]StaffUser, error)
	GetByID(ctx context.Context, id int64) (*StaffUser, error)
	GetByEmployeeNumber(ctx context.Context, empNum string) (*StaffUser, error)
	Create(ctx context.Context, staff *StaffUser) error
	Update(ctx context.Context, id int64, input UpdateStaffInput) (*StaffUser, error)
	ResetPassword(ctx context.Context, id int64, passwordHash string) error
}
