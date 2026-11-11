package usecase

import (
	"context"
	"testing"
	"time"

	"clinic-api/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	staffs map[string]*domain.StaffAuth
}

func (m *mockAuthRepo) FindByEmployeeNumber(ctx context.Context, empNum string) (*domain.StaffAuth, error) {
	s, ok := m.staffs[empNum]
	if !ok {
		return nil, domain.ErrInvalidCredentials
	}
	return s, nil
}

func (m *mockAuthRepo) FindByID(ctx context.Context, id int64) (*domain.StaffAuth, error) {
	for _, s := range m.staffs {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, domain.ErrStaffNotFound
}

func TestAuthUsecase_Login(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("klinik123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &mockAuthRepo{
		staffs: map[string]*domain.StaffAuth{
			"OWNER-01": {
				ID:             1,
				EmployeeNumber: "OWNER-01",
				Name:           "dr. Ratna",
				Role:           "SUPERADMIN_OWNER",
				ClinicUnit:     "ALL",
				PasswordHash:   string(hashedPassword),
				IsActive:       true,
			},
			"DOC-INACTIVE": {
				ID:             2,
				EmployeeNumber: "DOC-INACTIVE",
				Name:           "dr. Budi Nonaktif",
				Role:           "DOCTOR",
				ClinicUnit:     "KIDINEST",
				PasswordHash:   string(hashedPassword),
				IsActive:       false,
			},
		},
	}

	usecase := NewAuthUsecase(repo, "test-secret-32-bytes-long-key!!", 1*time.Hour)

	t.Run("Login sukses mengembalikan token dan user profile", func(t *testing.T) {
		resp, err := usecase.Login(context.Background(), domain.LoginRequest{
			EmployeeNumber: "OWNER-01",
			Password:       "klinik123",
		})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if resp.Token == "" {
			t.Errorf("expected non-empty token")
		}
		if resp.User.EmployeeNumber != "OWNER-01" {
			t.Errorf("expected OWNER-01, got %s", resp.User.EmployeeNumber)
		}
	})

	t.Run("Login gagal jika password salah", func(t *testing.T) {
		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			EmployeeNumber: "OWNER-01",
			Password:       "salahPassword",
		})
		if err != domain.ErrInvalidCredentials {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("Login gagal jika akun non-aktif", func(t *testing.T) {
		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			EmployeeNumber: "DOC-INACTIVE",
			Password:       "klinik123",
		})
		if err != domain.ErrStaffInactive {
			t.Fatalf("expected ErrStaffInactive, got %v", err)
		}
	})

	t.Run("Login gagal jika NIK tidak terdaftar", func(t *testing.T) {
		_, err := usecase.Login(context.Background(), domain.LoginRequest{
			EmployeeNumber: "UNKNOWN-99",
			Password:       "klinik123",
		})
		if err != domain.ErrInvalidCredentials {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}
