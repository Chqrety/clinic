package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
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

func getDoctorSpecialty(name, unit string) string {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "patricia") {
		return "Spesialis Anak & Vaksinasi Pediatri"
	}
	if strings.Contains(lower, "hendra") {
		return "Konsultan Tumbuh Kembang & Nutrisi"
	}
	if strings.Contains(lower, "dian") {
		return "Spesialis Dermatologi & Venereologi"
	}
	if strings.Contains(lower, "kevin") {
		return "Spesialis Kulit & Estetika Medis"
	}
	if unit == "KIDINEST" {
		return "Spesialis Anak (Pediatri)"
	}
	return "Dokter Spesialis"
}

func formatDaysRange(days []int) string {
	if len(days) == 0 {
		return "Belum Ada Jadwal"
	}
	sort.Ints(days)
	dayNames := map[int]string{
		0: "Minggu",
		1: "Senin",
		2: "Selasa",
		3: "Rabu",
		4: "Kamis",
		5: "Jumat",
		6: "Sabtu",
	}
	isConsecutive := true
	for i := 1; i < len(days); i++ {
		if days[i] != days[i-1]+1 {
			isConsecutive = false
			break
		}
	}
	if isConsecutive && len(days) >= 3 {
		return fmt.Sprintf("%s – %s", dayNames[days[0]], dayNames[days[len(days)-1]])
	}
	var names []string
	for _, d := range days {
		if name, ok := dayNames[d]; ok {
			names = append(names, name)
		}
	}
	return strings.Join(names, ", ")
}

func (u *ScheduleUsecase) GetPublicSchedulesSummary(ctx context.Context, unit string) ([]domain.DoctorScheduleSummary, error) {
	doctors, err := u.repo.GetDoctorsByUnit(ctx, unit)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar dokter: %w", err)
	}

	today := time.Now()
	todayStr := today.Format("2006-01-02")
	todayDayOfWeek := int(today.Weekday()) // 0=Sunday, 1=Monday...

	var summaries []domain.DoctorScheduleSummary
	for _, doc := range doctors {
		schedules, err := u.repo.GetDoctorSchedules(ctx, doc.ID)
		if err != nil {
			continue
		}

		overrides, _ := u.repo.GetDoctorOverrides(ctx, doc.ID, todayStr, todayStr)
		isOnLeave := false
		for _, o := range overrides {
			if o.OverrideDate == todayStr && o.IsCancelled {
				isOnLeave = true
				break
			}
		}

		var activeDays []int
		shiftName := "PAGI"
		startTime := "08:30"
		endTime := "12:00"
		slotMin := 30
		quota := 1
		isPracticingToday := false

		for _, s := range schedules {
			if s.IsActive {
				activeDays = append(activeDays, s.DayOfWeek)
				shiftName = s.ShiftName
				startTime = s.StartTime
				endTime = s.EndTime
				slotMin = s.SlotMinutes
				quota = s.QuotaPerSlot
				if s.DayOfWeek == todayDayOfWeek {
					isPracticingToday = true
				}
			}
		}

		daysText := formatDaysRange(activeDays)
		timeText := fmt.Sprintf("%s – %s WIB", startTime, endTime)

		todayStatus := "Libur Hari Ini"
		if isOnLeave {
			todayStatus = "Cuti Hari Ini"
			isPracticingToday = false
		} else if isPracticingToday {
			todayStatus = fmt.Sprintf("Praktik Hari Ini (%s – %s WIB)", startTime, endTime)
		} else if len(activeDays) > 0 {
			todayStatus = fmt.Sprintf("Jadwal: %s (%s)", daysText, timeText)
		}

		availDates, _ := u.GetDoctorAvailableSchedule(ctx, doc.ID, 7)

		summaries = append(summaries, domain.DoctorScheduleSummary{
			DoctorID:          doc.ID,
			DoctorName:        doc.Name,
			Role:              doc.Role,
			ClinicUnit:        doc.ClinicUnit,
			RoomName:          doc.RoomName,
			Specialty:         getDoctorSpecialty(doc.Name, doc.ClinicUnit),
			ShiftName:         shiftName,
			DaysText:          daysText,
			TimeText:          timeText,
			StartTime:         startTime,
			EndTime:           endTime,
			SlotMinutes:       slotMin,
			QuotaPerSlot:      quota,
			IsPracticingToday: isPracticingToday,
			IsOnLeaveToday:    isOnLeave,
			TodayStatus:       todayStatus,
			WeeklySchedules:   schedules,
			AvailableDates:    availDates,
		})
	}

	return summaries, nil
}

