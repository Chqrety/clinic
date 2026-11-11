package usecase

import (
	"context"
	"testing"
	"time"

	"clinic-api/internal/staff/domain"
	"golang.org/x/crypto/bcrypt"
)

type mockStaffRepo struct {
	staffs map[int64]*domain.StaffUser
	nextID int64
}

func newMockStaffRepo() *mockStaffRepo {
	return &mockStaffRepo{
		staffs: make(map[int64]*domain.StaffUser),
		nextID: 1,
	}
}

func (m *mockStaffRepo) List(ctx context.Context, filter domain.StaffFilter) ([]domain.StaffUser, error) {
	var list []domain.StaffUser
	for _, s := range m.staffs {
		if filter.ClinicUnit != "" && string(s.ClinicUnit) != filter.ClinicUnit {
			continue
		}
		if filter.Role != "" && string(s.Role) != filter.Role {
			continue
		}
		if filter.IsActive != nil && s.IsActive != *filter.IsActive {
			continue
		}
		list = append(list, *s)
	}
	return list, nil
}

func (m *mockStaffRepo) GetByID(ctx context.Context, id int64) (*domain.StaffUser, error) {
	s, ok := m.staffs[id]
	if !ok {
		return nil, domain.ErrStaffNotFound
	}
	return s, nil
}

func (m *mockStaffRepo) GetByEmployeeNumber(ctx context.Context, empNum string) (*domain.StaffUser, error) {
	for _, s := range m.staffs {
		if s.EmployeeNumber == empNum {
			return s, nil
		}
	}
	return nil, domain.ErrStaffNotFound
}

func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.StaffUser) error {
	for _, s := range m.staffs {
		if s.EmployeeNumber == staff.EmployeeNumber {
			return domain.ErrDuplicateEmployeeNum
		}
	}
	staff.ID = m.nextID
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()
	m.staffs[staff.ID] = staff
	m.nextID++
	return nil
}

func (m *mockStaffRepo) Update(ctx context.Context, id int64, input domain.UpdateStaffInput) (*domain.StaffUser, error) {
	s, ok := m.staffs[id]
	if !ok {
		return nil, domain.ErrStaffNotFound
	}
	s.Name = input.Name
	s.Role = input.Role
	s.ClinicUnit = input.ClinicUnit
	s.IsActive = input.IsActive
	s.UpdatedAt = time.Now()
	return s, nil
}

func (m *mockStaffRepo) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	s, ok := m.staffs[id]
	if !ok {
		return domain.ErrStaffNotFound
	}
	s.PasswordHash = passwordHash
	s.UpdatedAt = time.Now()
	return nil
}

func TestStaffUsecase_CRUD(t *testing.T) {
	repo := newMockStaffRepo()
	uc := NewStaffUsecase(repo)

	t.Run("Create staff sukses", func(t *testing.T) {
		staff, err := uc.CreateStaff(context.Background(), domain.CreateStaffInput{
			EmployeeNumber:  "DOC-01",
			Name:            "dr. Anisa",
			Role:            domain.RoleDoctor,
			ClinicUnit:      domain.UnitKidinest,
			InitialPassword: "password123",
		})
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if staff.ID == 0 {
			t.Errorf("expected non-zero ID")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte("password123")); err != nil {
			t.Errorf("expected bcrypt password to match")
		}
	})

	t.Run("Create staff duplikat NIK gagal", func(t *testing.T) {
		_, err := uc.CreateStaff(context.Background(), domain.CreateStaffInput{
			EmployeeNumber:  "DOC-01",
			Name:            "dr. Budi",
			Role:            domain.RoleDoctor,
			ClinicUnit:      domain.UnitKidinest,
			InitialPassword: "password123",
		})
		if err != domain.ErrDuplicateEmployeeNum {
			t.Fatalf("expected ErrDuplicateEmployeeNum, got %v", err)
		}
	})

	t.Run("Reset password berhasil", func(t *testing.T) {
		err := uc.ResetPassword(context.Background(), 1, "newPass456")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		s, _ := repo.GetByID(context.Background(), 1)
		if err := bcrypt.CompareHashAndPassword([]byte(s.PasswordHash), []byte("newPass456")); err != nil {
			t.Errorf("expected new password to match")
		}
	})
}
