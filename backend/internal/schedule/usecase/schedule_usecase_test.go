package usecase

import (
	"context"
	"testing"
	"time"

	"clinic-api/internal/schedule/domain"
)

type mockScheduleRepo struct {
	doctors   []domain.DoctorInfo
	schedules []domain.DoctorSchedule
	overrides []domain.ScheduleOverride
	counts    map[string]int
}

func (m *mockScheduleRepo) GetDoctorsByUnit(ctx context.Context, unit string) ([]domain.DoctorInfo, error) {
	return m.doctors, nil
}

func (m *mockScheduleRepo) GetDoctorSchedules(ctx context.Context, doctorID int64) ([]domain.DoctorSchedule, error) {
	return m.schedules, nil
}

func (m *mockScheduleRepo) GetDoctorOverrides(ctx context.Context, doctorID int64, startDate, endDate string) ([]domain.ScheduleOverride, error) {
	return m.overrides, nil
}

func (m *mockScheduleRepo) GetBookedSlotsCount(ctx context.Context, doctorID int64, date string) (map[string]int, error) {
	return m.counts, nil
}

func (m *mockScheduleRepo) AddDoctorOverride(ctx context.Context, override *domain.ScheduleOverride) error {
	m.overrides = append(m.overrides, *override)
	return nil
}

func (m *mockScheduleRepo) DeleteDoctorOverride(ctx context.Context, doctorID int64, date string) error {
	return nil
}

func (m *mockScheduleRepo) UpsertDoctorSchedule(ctx context.Context, schedule *domain.DoctorSchedule) error {
	return nil
}

func TestScheduleUsecase_GetDoctorAvailableSchedule(t *testing.T) {
	today := time.Now()
	todayDayOfWeek := int(today.Weekday())

	repo := &mockScheduleRepo{
		doctors: []domain.DoctorInfo{
			{ID: 1, Name: "dr. Hendra Pratama, Sp.A", Role: "DOCTOR", ClinicUnit: "KIDINEST"},
		},
		schedules: []domain.DoctorSchedule{
			{
				DoctorID:     1,
				DayOfWeek:    todayDayOfWeek,
				StartTime:    "08:00",
				EndTime:      "09:00",
				SlotMinutes:  30,
				QuotaPerSlot: 1,
				IsActive:     true,
			},
		},
		overrides: []domain.ScheduleOverride{},
		counts: map[string]int{
			"08:00": 0, // Belum ada booking -> available
			"08:30": 1, // Sudah 1 booking -> Full
		},
	}

	uc := NewScheduleUsecase(repo)

	dates, err := uc.GetDoctorAvailableSchedule(context.Background(), 1, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dates) != 7 {
		t.Fatalf("expected 7 dates, got %d", len(dates))
	}

	firstDate := dates[0]
	if !firstDate.IsAvailable {
		t.Fatalf("expected today to be available")
	}

	if len(firstDate.Slots) != 2 {
		t.Fatalf("expected 2 slots (08:00, 08:30), got %d", len(firstDate.Slots))
	}

	// Slot 1: 08:00 (quota 2, booked 1 -> remaining 1, available true)
	if firstDate.Slots[0].Available != 1 || !firstDate.Slots[0].IsAvailable {
		t.Errorf("slot 08:00 available expected 1 and available=true, got available=%d is_available=%v",
			firstDate.Slots[0].Available, firstDate.Slots[0].IsAvailable)
	}

	// Slot 2: 08:30 (quota 2, booked 2 -> remaining 0, available false)
	if firstDate.Slots[1].Available != 0 || firstDate.Slots[1].IsAvailable {
		t.Errorf("slot 08:30 available expected 0 and is_available=false, got available=%d is_available=%v",
			firstDate.Slots[1].Available, firstDate.Slots[1].IsAvailable)
	}
}

func TestScheduleUsecase_UpdateDoctorSchedule(t *testing.T) {
	repo := &mockScheduleRepo{}
	uc := NewScheduleUsecase(repo)

	// Valid input
	err := uc.UpdateDoctorSchedule(context.Background(), &domain.DoctorSchedule{
		DoctorID:  1,
		StartTime: "09:00",
		EndTime:   "12:00",
	})
	if err != nil {
		t.Fatalf("expected nil error on valid schedule, got %v", err)
	}

	// Invalid input (empty start_time)
	err = uc.UpdateDoctorSchedule(context.Background(), &domain.DoctorSchedule{
		DoctorID: 1,
	})
	if err == nil {
		t.Fatalf("expected error on empty times, got nil")
	}
}
