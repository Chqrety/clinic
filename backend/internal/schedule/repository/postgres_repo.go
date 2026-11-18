package repository

import (
	"context"
	"database/sql"
	"fmt"

	"clinic-api/internal/schedule/domain"
)

type PostgresScheduleRepository struct {
	db *sql.DB
}

func NewPostgresScheduleRepository(db *sql.DB) *PostgresScheduleRepository {
	return &PostgresScheduleRepository{db: db}
}

func (r *PostgresScheduleRepository) GetDoctorsByUnit(ctx context.Context, unit string) ([]domain.DoctorInfo, error) {
	query := `
	SELECT id, name, role, clinic_unit
	FROM staff_users
	WHERE role = 'DOCTOR' AND is_active = TRUE
	  AND ($1 = '' OR clinic_unit = $1 OR clinic_unit = 'ALL')
	ORDER BY id ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, unit)
	if err != nil {
		return nil, fmt.Errorf("query doctors by unit: %w", err)
	}
	defer rows.Close()

	var doctors []domain.DoctorInfo
	for rows.Next() {
		var d domain.DoctorInfo
		if err := rows.Scan(&d.ID, &d.Name, &d.Role, &d.ClinicUnit); err != nil {
			return nil, fmt.Errorf("scan doctor: %w", err)
		}
		// Format nama ruang poli
		if d.ClinicUnit == "KIDINEST" {
			if d.ID%2 == 1 {
				d.RoomName = "Ruang Poli Anak 1 (Lantai 1)"
			} else {
				d.RoomName = "Ruang Poli Anak 2 (Lantai 1)"
			}
		} else {
			if d.ID%2 == 1 {
				d.RoomName = "Ruang Laser & Estetika 1 (Lantai 2)"
			} else {
				d.RoomName = "Ruang Tindakan Bedah Kulit 2 (Lantai 2)"
			}
		}
		doctors = append(doctors, d)
	}
	return doctors, nil
}

func (r *PostgresScheduleRepository) GetDoctorSchedules(ctx context.Context, doctorID int64) ([]domain.DoctorSchedule, error) {
	query := `
	SELECT id, doctor_id, doctor_name, clinic_unit, day_of_week, shift_name,
	       start_time, end_time, slot_minutes, quota_per_slot, is_active, created_at, updated_at
	FROM doctor_schedules
	WHERE doctor_id = $1 AND is_active = TRUE
	ORDER BY day_of_week ASC, start_time ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("query doctor schedules: %w", err)
	}
	defer rows.Close()

	schedules := make([]domain.DoctorSchedule, 0)
	for rows.Next() {
		var s domain.DoctorSchedule
		if err := rows.Scan(
			&s.ID, &s.DoctorID, &s.DoctorName, &s.ClinicUnit, &s.DayOfWeek, &s.ShiftName,
			&s.StartTime, &s.EndTime, &s.SlotMinutes, &s.QuotaPerSlot, &s.IsActive,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *PostgresScheduleRepository) GetDoctorOverrides(ctx context.Context, doctorID int64, fromDate, toDate string) ([]domain.ScheduleOverride, error) {
	query := `
	SELECT id, doctor_id, TO_CHAR(override_date, 'YYYY-MM-DD'), is_cancelled, COALESCE(reason, ''), created_at
	FROM doctor_schedule_overrides
	WHERE doctor_id = $1 AND override_date BETWEEN $2::date AND $3::date;
	`
	rows, err := r.db.QueryContext(ctx, query, doctorID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("query doctor overrides: %w", err)
	}
	defer rows.Close()

	overrides := make([]domain.ScheduleOverride, 0)
	for rows.Next() {
		var o domain.ScheduleOverride
		if err := rows.Scan(&o.ID, &o.DoctorID, &o.OverrideDate, &o.IsCancelled, &o.Reason, &o.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan override: %w", err)
		}
		overrides = append(overrides, o)
	}
	return overrides, nil
}

func (r *PostgresScheduleRepository) GetBookedSlotsCount(ctx context.Context, doctorID int64, date string) (map[string]int, error) {
	query := `
	SELECT appointment_time, COUNT(*)
	FROM queues
	WHERE doctor_id = $1 AND appointment_date = $2::date AND queue_status != 'CANCELLED'
	GROUP BY appointment_time;
	`
	rows, err := r.db.QueryContext(ctx, query, doctorID, date)
	if err != nil {
		return nil, fmt.Errorf("query booked slots: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var slot string
		var count int
		if err := rows.Scan(&slot, &count); err != nil {
			return nil, fmt.Errorf("scan booked slot count: %w", err)
		}
		counts[slot] = count
	}
	return counts, nil
}

func (r *PostgresScheduleRepository) UpsertDoctorSchedule(ctx context.Context, s *domain.DoctorSchedule) error {
	slotMinutes := s.SlotMinutes
	if slotMinutes <= 0 {
		slotMinutes = 30
	}
	quota := s.QuotaPerSlot
	if quota <= 0 {
		quota = 1
	}

	if s.ID > 0 {
		updateByID := `
		UPDATE doctor_schedules
		SET shift_name = $1, start_time = $2, end_time = $3, is_active = $4, slot_minutes = $5, quota_per_slot = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7;
		`
		res, err := r.db.ExecContext(ctx, updateByID, s.ShiftName, s.StartTime, s.EndTime, s.IsActive, slotMinutes, quota, s.ID)
		if err != nil {
			return err
		}
		if ra, _ := res.RowsAffected(); ra > 0 {
			return nil
		}
	}

	updateByDocDay := `
	UPDATE doctor_schedules
	SET shift_name = $1, start_time = $2, end_time = $3, is_active = $4, slot_minutes = $5, quota_per_slot = $6, updated_at = CURRENT_TIMESTAMP
	WHERE doctor_id = $7 AND day_of_week = $8;
	`
	res, err := r.db.ExecContext(ctx, updateByDocDay, s.ShiftName, s.StartTime, s.EndTime, s.IsActive, slotMinutes, quota, s.DoctorID, s.DayOfWeek)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra > 0 {
		return nil
	}

	insertQuery := `
	INSERT INTO doctor_schedules (
		doctor_id, doctor_name, clinic_unit, day_of_week, shift_name,
		start_time, end_time, slot_minutes, quota_per_slot, is_active, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP);
	`
	_, err = r.db.ExecContext(ctx, insertQuery,
		s.DoctorID, s.DoctorName, s.ClinicUnit, s.DayOfWeek, s.ShiftName,
		s.StartTime, s.EndTime, slotMinutes, quota, s.IsActive,
	)
	return err
}

func (r *PostgresScheduleRepository) AddDoctorOverride(ctx context.Context, o *domain.ScheduleOverride) error {
	query := `
	INSERT INTO doctor_schedule_overrides (doctor_id, override_date, is_cancelled, reason)
	VALUES ($1, $2::date, $3, $4);
	`
	_, err := r.db.ExecContext(ctx, query, o.DoctorID, o.OverrideDate, o.IsCancelled, o.Reason)
	return err
}

func (r *PostgresScheduleRepository) DeleteDoctorOverride(ctx context.Context, doctorID int64, date string) error {
	query := `DELETE FROM doctor_schedule_overrides WHERE doctor_id = $1 AND override_date = $2::date;`
	_, err := r.db.ExecContext(ctx, query, doctorID, date)
	return err
}
