package usecase

import (
	"context"
	"strings"

	"clinic-api/internal/staff/domain"
	"golang.org/x/crypto/bcrypt"
)

type StaffUsecase struct {
	repo domain.StaffRepository
}

func NewStaffUsecase(repo domain.StaffRepository) *StaffUsecase {
	return &StaffUsecase{repo: repo}
}

func (u *StaffUsecase) ListStaff(ctx context.Context, filter domain.StaffFilter) ([]domain.StaffUser, error) {
	return u.repo.List(ctx, filter)
}

func (u *StaffUsecase) GetStaff(ctx context.Context, id int64) (*domain.StaffUser, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *StaffUsecase) CreateStaff(ctx context.Context, input domain.CreateStaffInput) (*domain.StaffUser, error) {
	empNum := strings.TrimSpace(input.EmployeeNumber)
	if empNum == "" {
		return nil, domain.ErrEmptyEmployeeNumber
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, domain.ErrEmptyName
	}

	pass := strings.TrimSpace(input.InitialPassword)
	if pass == "" {
		return nil, domain.ErrEmptyPassword
	}

	switch input.Role {
	case domain.RoleOwner, domain.RoleDoctor, domain.RoleCashier, domain.RolePharmacist:
	default:
		return nil, domain.ErrInvalidRole
	}

	switch input.ClinicUnit {
	case domain.UnitKidinest, domain.UnitSSC, domain.UnitAll:
	default:
		return nil, domain.ErrInvalidClinicUnit
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	staff := &domain.StaffUser{
		EmployeeNumber: empNum,
		Name:           name,
		Role:           input.Role,
		ClinicUnit:     input.ClinicUnit,
		PasswordHash:   string(hash),
		IsActive:       true,
	}

	if err := u.repo.Create(ctx, staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (u *StaffUsecase) UpdateStaff(ctx context.Context, id int64, input domain.UpdateStaffInput) (*domain.StaffUser, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, domain.ErrEmptyName
	}

	switch input.Role {
	case domain.RoleOwner, domain.RoleDoctor, domain.RoleCashier, domain.RolePharmacist:
	default:
		return nil, domain.ErrInvalidRole
	}

	switch input.ClinicUnit {
	case domain.UnitKidinest, domain.UnitSSC, domain.UnitAll:
	default:
		return nil, domain.ErrInvalidClinicUnit
	}

	cleanInput := domain.UpdateStaffInput{
		Name:       name,
		Role:       input.Role,
		ClinicUnit: input.ClinicUnit,
		IsActive:   input.IsActive,
	}

	return u.repo.Update(ctx, id, cleanInput)
}

func (u *StaffUsecase) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	pass := strings.TrimSpace(newPassword)
	if pass == "" {
		return domain.ErrEmptyPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.ResetPassword(ctx, id, string(hash))
}
