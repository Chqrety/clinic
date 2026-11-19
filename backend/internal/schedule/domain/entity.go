package domain

import (
	"context"
	"time"
)

type DoctorSchedule struct {
	ID           int64     `json:"id"`
	DoctorID     int64     `json:"doctor_id"`
	DoctorName   string    `json:"doctor_name"`
	ClinicUnit   string    `json:"clinic_unit"`
	DayOfWeek    int       `json:"day_of_week"` // 0=Minggu, 1=Senin, ..., 6=Sabtu
	ShiftName    string    `json:"shift_name"`  // "PAGI" | "SORE"
	StartTime    string    `json:"start_time"`  // "08:30"
	EndTime      string    `json:"end_time"`    // "12:00"
	SlotMinutes  int       `json:"slot_minutes"`
	QuotaPerSlot int       `json:"quota_per_slot"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ScheduleOverride struct {
	ID           int64     `json:"id"`
	DoctorID     int64     `json:"doctor_id"`
	OverrideDate string    `json:"override_date"` // YYYY-MM-DD
	IsCancelled  bool      `json:"is_cancelled"`
	Reason       string    `json:"reason"`
	CreatedAt    time.Time `json:"created_at"`
}

type AvailableSlot struct {
	TimeSlot    string `json:"time_slot"`   // misal "08:30 - 09:00"
	StartTime   string `json:"start_time"`  // "08:30"
	EndTime     string `json:"end_time"`    // "09:00"
	TotalQuota  int    `json:"total_quota"` // 3
	Booked      int    `json:"booked"`      // 1
	Available   int    `json:"available"`   // 2
	IsAvailable bool   `json:"is_available"`
}

type AvailableDate struct {
	Date        string          `json:"date"`        // "2026-11-20"
	DayName     string          `json:"day_name"`    // "Jumat"
	DayOfWeek   int             `json:"day_of_week"` // 5
	ShiftName   string          `json:"shift_name"`  // "PAGI"
	IsAvailable bool            `json:"is_available"`
	Slots       []AvailableSlot `json:"slots"`
}

type DoctorInfo struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	ClinicUnit string `json:"clinic_unit"`
	RoomName   string `json:"room_name,omitempty"`
}

type DoctorScheduleSummary struct {
	DoctorID          int64            `json:"doctor_id"`
	DoctorName        string           `json:"doctor_name"`
	Role              string           `json:"role"`
	ClinicUnit        string           `json:"clinic_unit"`
	RoomName          string           `json:"room_name"`
	Specialty         string           `json:"specialty"`
	ShiftName         string           `json:"shift_name"`
	DaysText          string           `json:"days_text"`
	TimeText          string           `json:"time_text"`
	StartTime         string           `json:"start_time"`
	EndTime           string           `json:"end_time"`
	SlotMinutes       int              `json:"slot_minutes"`
	QuotaPerSlot      int              `json:"quota_per_slot"`
	IsPracticingToday bool             `json:"is_practicing_today"`
	IsOnLeaveToday    bool             `json:"is_on_leave_today"`
	TodayStatus       string           `json:"today_status"`
	WeeklySchedules   []DoctorSchedule `json:"weekly_schedules"`
	AvailableDates    []AvailableDate  `json:"available_dates,omitempty"`
}

type ScheduleRepository interface {
	GetDoctorsByUnit(ctx context.Context, unit string) ([]DoctorInfo, error)
	GetDoctorSchedules(ctx context.Context, doctorID int64) ([]DoctorSchedule, error)
	GetDoctorOverrides(ctx context.Context, doctorID int64, fromDate, toDate string) ([]ScheduleOverride, error)
	GetBookedSlotsCount(ctx context.Context, doctorID int64, date string) (map[string]int, error)
	UpsertDoctorSchedule(ctx context.Context, s *DoctorSchedule) error
	AddDoctorOverride(ctx context.Context, o *ScheduleOverride) error
	DeleteDoctorOverride(ctx context.Context, doctorID int64, date string) error
}
