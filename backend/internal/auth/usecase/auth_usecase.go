package usecase

import (
	"context"
	"strings"
	"time"

	"clinic-api/internal/auth/domain"
	"clinic-api/internal/shared/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo      domain.AuthRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthUsecase(repo domain.AuthRepository, jwtSecret string, tokenTTL time.Duration) *AuthUsecase {
	if tokenTTL <= 0 {
		tokenTTL = 24 * time.Hour
	}
	return &AuthUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, input domain.LoginRequest) (*domain.LoginResponse, error) {
	empNum := strings.TrimSpace(input.EmployeeNumber)
	if empNum == "" || strings.TrimSpace(input.Password) == "" {
		return nil, domain.ErrInvalidCredentials
	}

	staff, err := u.repo.FindByEmployeeNumber(ctx, empNum)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !staff.IsActive {
		return nil, domain.ErrStaffInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := middleware.GenerateToken(
		u.jwtSecret,
		staff.ID,
		staff.EmployeeNumber,
		staff.Name,
		staff.Role,
		staff.ClinicUnit,
		u.tokenTTL,
	)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User: domain.UserProfile{
			ID:             staff.ID,
			EmployeeNumber: staff.EmployeeNumber,
			Name:           staff.Name,
			Role:           staff.Role,
			ClinicUnit:     staff.ClinicUnit,
		},
	}, nil
}

func (u *AuthUsecase) GetProfile(ctx context.Context, id int64) (*domain.UserProfile, error) {
	staff, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !staff.IsActive {
		return nil, domain.ErrStaffInactive
	}
	return &domain.UserProfile{
		ID:             staff.ID,
		EmployeeNumber: staff.EmployeeNumber,
		Name:           staff.Name,
		Role:           staff.Role,
		ClinicUnit:     staff.ClinicUnit,
	}, nil
}
