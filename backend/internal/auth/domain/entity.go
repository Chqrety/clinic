package domain

import (
	"context"
	"time"
)

type StaffAuth struct {
	ID             int64     `json:"id"`
	EmployeeNumber string    `json:"employee_number"`
	Name           string    `json:"name"`
	Role           string    `json:"role"`
	ClinicUnit     string    `json:"clinic_unit"`
	PasswordHash   string    `json:"-"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LoginRequest struct {
	EmployeeNumber string `json:"employee_number"`
	Password       string `json:"password"`
}

type UserProfile struct {
	ID             int64  `json:"id"`
	EmployeeNumber string `json:"employee_number"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	ClinicUnit     string `json:"clinic_unit"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type AuthRepository interface {
	FindByEmployeeNumber(ctx context.Context, empNum string) (*StaffAuth, error)
	FindByID(ctx context.Context, id int64) (*StaffAuth, error)
}
