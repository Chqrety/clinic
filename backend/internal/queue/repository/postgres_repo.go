package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/queue/domain"
)

type PostgresQueueRepository struct {
	db *sql.DB
}

func NewPostgresQueueRepository(db *sql.DB) *PostgresQueueRepository {
	return &PostgresQueueRepository{db: db}
}

// ==================== PATIENT REPOSITORY ====================

func (r *PostgresQueueRepository) FindOrCreateByPhone(ctx context.Context, name, phone string) (*domain.Patient, error) {
	selectQuery := `
	SELECT id, name, phone, created_at, updated_at
	FROM patients
	WHERE phone = $1
	LIMIT 1;
	`
	var p domain.Patient
	err := r.db.QueryRowContext(ctx, selectQuery, phone).Scan(
		&p.ID,
		&p.Name,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == nil {
		// Update name if changed
		if p.Name != name {
			_, _ = r.db.ExecContext(ctx, "UPDATE patients SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", name, p.ID)
			p.Name = name
		}
		return &p, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	insertQuery := `
	INSERT INTO patients (name, phone)
	VALUES ($1, $2)
	RETURNING id, name, phone, created_at, updated_at;
	`
	err = r.db.QueryRowContext(ctx, insertQuery, name, phone).Scan(
		&p.ID,
		&p.Name,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PostgresQueueRepository) GetPatientByID(ctx context.Context, id int64) (*domain.Patient, error) {
	query := `SELECT id, name, phone, created_at, updated_at FROM patients WHERE id = $1 LIMIT 1;`
	var p domain.Patient
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPatientNotFound
		}
		return nil, err
	}
	return &p, nil
}

// ==================== QUEUE REPOSITORY ====================

func (r *PostgresQueueRepository) GenerateNextQueueNumber(ctx context.Context, date time.Time, unit domain.ClinicUnit, category domain.QueueCategory) (string, error) {
	prefix := "A"
	if unit == domain.UnitSSC {
		prefix = "S"
	} else {
		switch category {
		case domain.CategoryVaksin:
			prefix = "A"
		case domain.CategorySakit:
			prefix = "B"
		case domain.CategoryTumbang:
			prefix = "C"
		case domain.CategoryPaket:
			prefix = "D"
		case domain.CategoryWalkin:
			prefix = "E"
		default:
			prefix = "A"
		}
	}

	dateOnly := date.Format("2006-01-02")

	upsertQuery := `
	INSERT INTO daily_queue_counters (counter_date, clinic_unit, category_prefix, last_seq)
	VALUES ($1, $2, $3, 1)
	ON CONFLICT (counter_date, clinic_unit, category_prefix)
	DO UPDATE SET last_seq = daily_queue_counters.last_seq + 1
	RETURNING last_seq;
	`

	var lastSeq int
	err := r.db.QueryRowContext(ctx, upsertQuery, dateOnly, string(unit), prefix).Scan(&lastSeq)
	if err != nil {
		return "", fmt.Errorf("gagal alokasi nomor urut harian: %w", err)
	}

	return fmt.Sprintf("%s-%03d", prefix, lastSeq), nil
}

func (r *PostgresQueueRepository) Create(ctx context.Context, q *domain.QueueItem) error {
	query := `
	INSERT INTO queues (
		patient_id, clinic_unit, category, doctor_id, doctor_name,
		queue_number, registration_fee, payment_method, payment_status, queue_status,
		paid_at, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP)
	RETURNING id, created_at;
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		q.PatientID,
		string(q.ClinicUnit),
		string(q.Category),
		q.DoctorID,
		q.DoctorName,
		q.QueueNumber,
		q.RegistrationFee,
		string(q.PaymentMethod),
		string(q.PaymentStatus),
		string(q.QueueStatus),
		q.PaidAt,
	).Scan(&q.ID, &q.CreatedAt)
}

func (r *PostgresQueueRepository) GetByID(ctx context.Context, id int64) (*domain.QueueItem, error) {
	query := `
	SELECT 
		q.id, q.patient_id, p.name, p.phone, q.clinic_unit, q.category,
		q.doctor_id, q.doctor_name, q.queue_number, q.registration_fee,
		q.payment_method, q.payment_status, q.queue_status,
		q.paid_at, q.called_at, q.completed_at, q.created_at
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.id = $1
	LIMIT 1;
	`
	var q domain.QueueItem
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&q.ID,
		&q.PatientID,
		&q.PatientName,
		&q.PatientPhone,
		&q.ClinicUnit,
		&q.Category,
		&q.DoctorID,
		&q.DoctorName,
		&q.QueueNumber,
		&q.RegistrationFee,
		&q.PaymentMethod,
		&q.PaymentStatus,
		&q.QueueStatus,
		&q.PaidAt,
		&q.CalledAt,
		&q.CompletedAt,
		&q.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQueueNotFound
		}
		return nil, err
	}
	return &q, nil
}

func (r *PostgresQueueRepository) UpdatePaymentSuccess(ctx context.Context, id int64, queueNumber string) (*domain.QueueItem, error) {
	query := `
	UPDATE queues
	SET payment_status = 'PAID', queue_status = 'WAITING', queue_number = $1, paid_at = CURRENT_TIMESTAMP
	WHERE id = $2
	RETURNING 
		id, patient_id, clinic_unit, category, doctor_id, doctor_name,
		queue_number, registration_fee, payment_method, payment_status, queue_status,
		paid_at, called_at, completed_at, created_at;
	`
	var q domain.QueueItem
	err := r.db.QueryRowContext(ctx, query, queueNumber, id).Scan(
		&q.ID,
		&q.PatientID,
		&q.ClinicUnit,
		&q.Category,
		&q.DoctorID,
		&q.DoctorName,
		&q.QueueNumber,
		&q.RegistrationFee,
		&q.PaymentMethod,
		&q.PaymentStatus,
		&q.QueueStatus,
		&q.PaidAt,
		&q.CalledAt,
		&q.CompletedAt,
		&q.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQueueNotFound
		}
		return nil, err
	}

	// Fetch patient name & phone
	var pName, pPhone string
	_ = r.db.QueryRowContext(ctx, "SELECT name, phone FROM patients WHERE id = $1", q.PatientID).Scan(&pName, &pPhone)
	q.PatientName = pName
	q.PatientPhone = pPhone

	return &q, nil
}

func (r *PostgresQueueRepository) UpdateQueueStatus(ctx context.Context, id int64, status domain.QueueStatus) error {
	var query string
	switch status {
	case domain.QueueCalling:
		query = "UPDATE queues SET queue_status = $1, called_at = CURRENT_TIMESTAMP WHERE id = $2"
	case domain.QueueCompleted:
		query = "UPDATE queues SET queue_status = $1, completed_at = CURRENT_TIMESTAMP WHERE id = $2"
	default:
		query = "UPDATE queues SET queue_status = $1 WHERE id = $2"
	}

	res, err := r.db.ExecContext(ctx, query, string(status), id)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return domain.ErrQueueNotFound
	}
	return nil
}

func (r *PostgresQueueRepository) GetCallingQueue(ctx context.Context, unit string) (*domain.QueueItem, error) {
	query := `
	SELECT 
		q.id, q.patient_id, p.name, p.phone, q.clinic_unit, q.category,
		q.doctor_id, q.doctor_name, q.queue_number, q.registration_fee,
		q.payment_method, q.payment_status, q.queue_status,
		q.paid_at, q.called_at, q.completed_at, q.created_at
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status = 'CALLING'
	`
	var args []any
	if unit != "" && unit != "ALL" {
		query += " AND q.clinic_unit = $1"
		args = append(args, unit)
	}
	query += " ORDER BY q.called_at DESC NULLS LAST, q.id DESC LIMIT 1;"

	var q domain.QueueItem
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&q.ID,
		&q.PatientID,
		&q.PatientName,
		&q.PatientPhone,
		&q.ClinicUnit,
		&q.Category,
		&q.DoctorID,
		&q.DoctorName,
		&q.QueueNumber,
		&q.RegistrationFee,
		&q.PaymentMethod,
		&q.PaymentStatus,
		&q.QueueStatus,
		&q.PaidAt,
		&q.CalledAt,
		&q.CompletedAt,
		&q.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No queue currently calling
		}
		return nil, err
	}
	return &q, nil
}

func (r *PostgresQueueRepository) ListUpcomingQueues(ctx context.Context, unit string, limit int) ([]domain.QueueItem, error) {
	if limit <= 0 {
		limit = 5
	}
	query := `
	SELECT 
		q.id, q.patient_id, p.name, p.phone, q.clinic_unit, q.category,
		q.doctor_id, q.doctor_name, q.queue_number, q.registration_fee,
		q.payment_method, q.payment_status, q.queue_status,
		q.paid_at, q.called_at, q.completed_at, q.created_at
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status = 'WAITING'
	`
	var args []any
	argIdx := 1
	if unit != "" && unit != "ALL" {
		query += fmt.Sprintf(" AND q.clinic_unit = $%d", argIdx)
		args = append(args, unit)
		argIdx++
	}
	query += fmt.Sprintf(" ORDER BY q.id ASC LIMIT $%d;", argIdx)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.QueueItem
	for rows.Next() {
		var q domain.QueueItem
		if err := rows.Scan(
			&q.ID,
			&q.PatientID,
			&q.PatientName,
			&q.PatientPhone,
			&q.ClinicUnit,
			&q.Category,
			&q.DoctorID,
			&q.DoctorName,
			&q.QueueNumber,
			&q.RegistrationFee,
			&q.PaymentMethod,
			&q.PaymentStatus,
			&q.QueueStatus,
			&q.PaidAt,
			&q.CalledAt,
			&q.CompletedAt,
			&q.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, q)
	}

	return list, rows.Err()
}
