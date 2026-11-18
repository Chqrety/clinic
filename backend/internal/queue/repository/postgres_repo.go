package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
		if p.Name != name && name != "" {
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

func (r *PostgresQueueRepository) LookupPatient(ctx context.Context, query string) (*domain.PatientLookupResult, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, domain.ErrPatientNotFound
	}
	sqlQuery := `
	SELECT id, name, phone, created_at, updated_at
	FROM patients
	WHERE phone = $1 OR LOWER(name) = LOWER($1)
	ORDER BY id DESC
	LIMIT 1;
	`
	var p domain.Patient
	err := r.db.QueryRowContext(ctx, sqlQuery, q).Scan(
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

	result := &domain.PatientLookupResult{
		Patient: p,
	}

	// Cek apakah pasien memiliki tiket aktif (BOOKED, WAITING, CALLING)
	activeTicketQuery := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.patient_id = $1 AND q.queue_status IN ('BOOKED', 'WAITING', 'CALLING')
	ORDER BY q.id DESC
	LIMIT 1;
	`, queueSelectFields)

	var activeTicket domain.QueueItem
	if scanErr := scanQueueRow(r.db.QueryRowContext(ctx, activeTicketQuery, p.ID), &activeTicket); scanErr == nil {
		result.ActiveTicket = &activeTicket
	}

	return result, nil
}

// ==================== QUEUE REPOSITORY ====================

const queueSelectFields = `
	q.id, q.patient_id, p.name, p.phone, q.clinic_unit, q.category,
	q.doctor_id, q.doctor_name, q.queue_number, q.registration_fee,
	q.payment_method, q.payment_status, q.queue_status,
	COALESCE(q.payment_reference, ''), COALESCE(q.payment_url, ''),
	q.paid_at, q.called_at, q.completed_at, q.created_at,
	TO_CHAR(q.appointment_date, 'YYYY-MM-DD'), q.appointment_time, COALESCE(q.patient_type, 'BARU'),
	COALESCE(q.reschedule_count, 0), q.refund_status, q.refund_notes, q.refund_bank, q.refund_account_number, q.refund_account_name
`

func scanQueueRow(scanner interface{ Scan(dest ...any) error }, q *domain.QueueItem) error {
	var appDate, appTime, refundStatus, refundNotes, refundBank, refundAccNum, refundAccName sql.NullString
	err := scanner.Scan(
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
		&q.PaymentReference,
		&q.PaymentURL,
		&q.PaidAt,
		&q.CalledAt,
		&q.CompletedAt,
		&q.CreatedAt,
		&appDate,
		&appTime,
		&q.PatientType,
		&q.RescheduleCount,
		&refundStatus,
		&refundNotes,
		&refundBank,
		&refundAccNum,
		&refundAccName,
	)
	if err != nil {
		return err
	}
	if appDate.Valid {
		q.AppointmentDate = &appDate.String
	}
	if appTime.Valid {
		q.AppointmentTime = &appTime.String
	}
	if refundStatus.Valid {
		q.RefundStatus = &refundStatus.String
	}
	if refundNotes.Valid {
		q.RefundNotes = &refundNotes.String
	}
	if refundBank.Valid {
		q.RefundBank = &refundBank.String
	}
	if refundAccNum.Valid {
		q.RefundAccountNumber = &refundAccNum.String
	}
	if refundAccName.Valid {
		q.RefundAccountName = &refundAccName.String
	}
	return nil
}

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
		payment_reference, payment_url, paid_at,
		appointment_date, appointment_time, patient_type, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, CURRENT_TIMESTAMP)
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
		q.PaymentReference,
		q.PaymentURL,
		q.PaidAt,
		q.AppointmentDate,
		q.AppointmentTime,
		q.PatientType,
	).Scan(&q.ID, &q.CreatedAt)
}

func (r *PostgresQueueRepository) GetByID(ctx context.Context, id int64) (*domain.QueueItem, error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.id = $1
	LIMIT 1;
	`, queueSelectFields)

	var q domain.QueueItem
	err := scanQueueRow(r.db.QueryRowContext(ctx, query, id), &q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQueueNotFound
		}
		return nil, err
	}
	return &q, nil
}

func (r *PostgresQueueRepository) UpdatePaymentDetails(ctx context.Context, id int64, reference, paymentURL string) error {
	query := `UPDATE queues SET payment_reference = $1, payment_url = $2 WHERE id = $3;`
	_, err := r.db.ExecContext(ctx, query, reference, paymentURL, id)
	return err
}

func (r *PostgresQueueRepository) GetByPaymentReference(ctx context.Context, reference string) (*domain.QueueItem, error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.payment_reference = $1 OR q.id = CASE WHEN $1 ~ '^KID-([0-9]+)' THEN SUBSTRING($1 FROM 'KID-([0-9]+)')::integer ELSE 0 END
	LIMIT 1;
	`, queueSelectFields)

	var q domain.QueueItem
	err := scanQueueRow(r.db.QueryRowContext(ctx, query, reference), &q)
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
	SET payment_status = 'PAID',
	    queue_status = CASE 
	        WHEN appointment_date IS NOT NULL AND appointment_date > CURRENT_DATE THEN 'BOOKED'
	        ELSE 'WAITING'
	    END,
	    queue_number = $1,
	    paid_at = CURRENT_TIMESTAMP
	WHERE id = $2 AND payment_status != 'PAID';
	`
	res, err := r.db.ExecContext(ctx, query, queueNumber, id)
	if err != nil {
		return nil, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if ra == 0 {
		return nil, domain.ErrQueueNotFound
	}

	return r.GetByID(ctx, id)
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
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status = 'CALLING'
	`, queueSelectFields)

	var args []any
	if unit != "" && unit != "ALL" {
		query += " AND q.clinic_unit = $1"
		args = append(args, unit)
	}
	query += " ORDER BY q.called_at DESC NULLS LAST, q.id DESC LIMIT 1;"

	var q domain.QueueItem
	err := scanQueueRow(r.db.QueryRowContext(ctx, query, args...), &q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *PostgresQueueRepository) ListUpcomingQueues(ctx context.Context, unit string, limit int) ([]domain.QueueItem, error) {
	if limit <= 0 {
		limit = 5
	}
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status = 'WAITING'
	`, queueSelectFields)

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
		if err := scanQueueRow(rows, &q); err != nil {
			return nil, err
		}
		list = append(list, q)
	}

	return list, rows.Err()
}

func (r *PostgresQueueRepository) ListReservations(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status = 'BOOKED'
	`, queueSelectFields)

	var args []any
	if unit != "" && unit != "ALL" {
		query += " AND q.clinic_unit = $1"
		args = append(args, unit)
	}
	query += " ORDER BY q.id ASC;"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.QueueItem
	for rows.Next() {
		var q domain.QueueItem
		if err := scanQueueRow(rows, &q); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

func (r *PostgresQueueRepository) ListActiveDoctorQueues(ctx context.Context) ([]domain.QueueItem, error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.queue_status IN ('WAITING', 'CALLING', 'IN_EXAM')
	ORDER BY q.id ASC;
	`, queueSelectFields)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.QueueItem
	for rows.Next() {
		var q domain.QueueItem
		if err := scanQueueRow(rows, &q); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

func (r *PostgresQueueRepository) Reschedule(ctx context.Context, id int64, newDate, newTime string) error {
	// Guard: pastikan slot waktu belum diambil pasien aktif lain pada dokter yang sama
	var conflictCount int
	conflictCheck := `
	SELECT COUNT(*) FROM queues 
	WHERE doctor_id = (SELECT doctor_id FROM queues WHERE id = $3) 
	  AND appointment_date = $1 
	  AND appointment_time = $2 
	  AND queue_status IN ('BOOKED', 'WAITING', 'CALLING') 
	  AND id != $3;
	`
	if err := r.db.QueryRowContext(ctx, conflictCheck, newDate, newTime, id).Scan(&conflictCount); err == nil && conflictCount > 0 {
		return fmt.Errorf("slot waktu %s jam %s sudah terisi oleh pasien lain", newDate, newTime)
	}

	query := `
	UPDATE queues
	SET appointment_date = $1, 
	    appointment_time = $2, 
	    reschedule_count = reschedule_count + 1,
	    queue_status = CASE 
	        WHEN $1 <= TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD') AND payment_status = 'PAID' THEN 'WAITING'
	        WHEN $1 > TO_CHAR(CURRENT_DATE, 'YYYY-MM-DD') THEN 'BOOKED'
	        ELSE queue_status
	    END
	WHERE id = $3;
	`
	res, err := r.db.ExecContext(ctx, query, newDate, newTime, id)
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

func (r *PostgresQueueRepository) RequestRefund(ctx context.Context, id int64, bank, accNumber, accName, notes string) error {
	query := `
	UPDATE queues
	SET refund_status = 'REQUESTED', refund_bank = $1, refund_account_number = $2, refund_account_name = $3, refund_notes = $4
	WHERE id = $5;
	`
	res, err := r.db.ExecContext(ctx, query, bank, accNumber, accName, notes, id)
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

func (r *PostgresQueueRepository) ProcessRefund(ctx context.Context, id int64, status string, notes string) error {
	query := `
	UPDATE queues
	SET refund_status = $1,
	    refund_notes = COALESCE(NULLIF($2, ''), refund_notes),
	    queue_status = CASE WHEN $1 = 'APPROVED' THEN 'CANCELLED' ELSE queue_status END,
	    payment_status = CASE WHEN $1 = 'APPROVED' THEN 'REFUNDED' ELSE payment_status END
	WHERE id = $3;
	`
	res, err := r.db.ExecContext(ctx, query, status, notes, id)
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

func (r *PostgresQueueRepository) ListRefundRequests(ctx context.Context, unit string) ([]domain.QueueItem, error) {
	query := fmt.Sprintf(`
	SELECT %s
	FROM queues q
	JOIN patients p ON q.patient_id = p.id
	WHERE q.refund_status IS NOT NULL AND q.refund_status != 'NONE'
	`, queueSelectFields)

	var args []any
	if unit != "" && unit != "ALL" {
		query += " AND q.clinic_unit = $1"
		args = append(args, unit)
	}
	query += " ORDER BY q.id DESC;"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.QueueItem
	for rows.Next() {
		var q domain.QueueItem
		if err := scanQueueRow(rows, &q); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}
