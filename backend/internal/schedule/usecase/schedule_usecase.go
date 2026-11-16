package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/schedule/domain"
)

type ScheduleUsecase struct {
	repo domain.ScheduleRepository
}

func NewScheduleUsecase(repo domain.ScheduleRepository) *ScheduleUsecase {
	return &ScheduleUsecase{repo: repo}
}

func (u *ScheduleUsecase) GetDoctors(ctx context.Context, unit string) ([]domain.DoctorInfo, error) {
	return u.repo.GetDoctorsByUnit(ctx, unit)
}

var indonesianDayNames = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

func (u *ScheduleUsecase) GetDoctorAvailableSchedule(ctx context.Context, doctorID int64, daysAhead int) ([]domain.AvailableDate, error) {
	if daysAhead <= 0 || daysAhead > 90 {
		daysAhead = 30
	}

	schedules, err := u.repo.GetDoctorSchedules(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil jadwal dokter: %w", err)
	}

	schedByDay := make(map[int]domain.DoctorSchedule)
	for _, s := range schedules {
		schedByDay[s.DayOfWeek] = s
	}

	today := time.Now()
	startDate := today.Format("2006-01-02")
	endDate := today.AddDate(0, 0, daysAhead).Format("2006-01-02")

	overrides, err := u.repo.GetDoctorOverrides(ctx, doctorID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data libur dokter: %w", err)
	}

	overrideMap := make(map[string]domain.ScheduleOverride)
	for _, o := range overrides {
		overrideMap[o.OverrideDate] = o
	}

	var results []domain.AvailableDate

	for i := 0; i < daysAhead; i++ {
		targetDate := today.AddDate(0, 0, i)
		dateStr := targetDate.Format("2006-01-02")
		weekday := int(targetDate.Weekday())

		availDate := domain.AvailableDate{
			Date:        dateStr,
			DayName:     indonesianDayNames[targetDate.Weekday()],
			DayOfWeek:   weekday,
			IsAvailable: false,
			Slots:       []domain.AvailableSlot{},
		}

		// 1. Cek apakah dokter cuti/libur di tanggal ini
		if ovr, exists := overrideMap[dateStr]; exists && ovr.IsCancelled {
			results = append(results, availDate)
			continue
		}

		// 2. Cek apakah ada jadwal rutin di hari ini
		sched, hasSched := schedByDay[weekday]
		if !hasSched || !sched.IsActive {
			results = append(results, availDate)
			continue
		}

		availDate.ShiftName = sched.ShiftName

		// 3. Hitung slot waktu & cek kuota terisi
		bookedCounts, err := u.repo.GetBookedSlotsCount(ctx, doctorID, dateStr)
		if err != nil {
			bookedCounts = make(map[string]int)
		}

		slots := generateSlots(sched.StartTime, sched.EndTime, sched.SlotMinutes, sched.QuotaPerSlot, bookedCounts)
		availDate.Slots = slots

		for _, s := range slots {
			if s.IsAvailable {
				availDate.IsAvailable = true
				break
			}
		}

		results = append(results, availDate)
	}

	return results, nil
}

func generateSlots(startTimeStr, endTimeStr string, slotMinutes int, quota int, booked map[string]int) []domain.AvailableSlot {
	if slotMinutes <= 0 {
		slotMinutes = 30
	}
	if quota <= 0 {
		quota = 1
	}

	tLayout := "15:04"
	cur, err := time.Parse(tLayout, startTimeStr)
	if err != nil {
		cur, _ = time.Parse(tLayout, "09:00")
	}
	end, err := time.Parse(tLayout, endTimeStr)
	if err != nil {
		end, _ = time.Parse(tLayout, "12:00")
	}

	var slots []domain.AvailableSlot
	duration := time.Duration(slotMinutes) * time.Minute

	for cur.Before(end) {
		next := cur.Add(duration)
		if next.After(end) {
			break
		}

		slotLabel := fmt.Sprintf("%s - %s", cur.Format("15:04"), next.Format("15:04"))
		bookedNum := booked[slotLabel]
		if bookedNum == 0 {
			bookedNum = booked[cur.Format("15:04")]
		}
		availNum := quota - bookedNum
		if availNum < 0 {
			availNum = 0
		}

		slots = append(slots, domain.AvailableSlot{
			TimeSlot:    slotLabel,
			StartTime:   cur.Format("15:04"),
			EndTime:     next.Format("15:04"),
			TotalQuota:  quota,
			Booked:      bookedNum,
			Available:   availNum,
			IsAvailable: availNum > 0,
		})

		cur = next
	}

	return slots
}

func (u *ScheduleUsecase) AddDoctorLeave(ctx context.Context, doctorID int64, date string, reason string) error {
	return u.repo.AddDoctorOverride(ctx, &domain.ScheduleOverride{
		DoctorID:     doctorID,
		OverrideDate: date,
		IsCancelled:  true,
		Reason:       reason,
	})
}

func (u *ScheduleUsecase) DeleteDoctorLeave(ctx context.Context, doctorID int64, date string) error {
	return u.repo.DeleteDoctorOverride(ctx, doctorID, date)
}

func (u *ScheduleUsecase) GetDoctorLeaveList(ctx context.Context, doctorID int64) ([]domain.ScheduleOverride, error) {
	today := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 90).Format("2006-01-02")
	return u.repo.GetDoctorOverrides(ctx, doctorID, today, end)
}

func (u *ScheduleUsecase) GetDoctorWeeklySchedules(ctx context.Context, doctorID int64) ([]domain.DoctorSchedule, error) {
	return u.repo.GetDoctorSchedules(ctx, doctorID)
}

func (u *ScheduleUsecase) UpdateDoctorSchedule(ctx context.Context, s *domain.DoctorSchedule) error {
	if s.DoctorID <= 0 {
		return errors.New("doctor_id wajib diisi")
	}
	if s.StartTime == "" || s.EndTime == "" {
		return errors.New("start_time dan end_time wajib diisi")
	}
	if s.SlotMinutes <= 0 {
		s.SlotMinutes = 30
	}
	if s.QuotaPerSlot <= 0 {
		s.QuotaPerSlot = 1
	}
	return u.repo.UpsertDoctorSchedule(ctx, s)
}


func (u *ScheduleUsecase) RemoveDoctorLeave(ctx context.Context, doctorID int64, date string) error {
	return u.repo.DeleteDoctorOverride(ctx, doctorID, date)
}
