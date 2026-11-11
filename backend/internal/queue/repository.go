package queue

import (
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/domain"
)

// Repository defines data access operations for queue and patient registration.
type Repository interface {
	GetQueues() ([]domain.QueueItem, error)
	GetQueueByID(id int64) (*domain.QueueItem, error)
	FindPatientByNIK(tx *sql.Tx, nik string) (int64, string, error)
	CreatePatient(tx *sql.Tx, nik, name, phone, birthDate, gender, ihs string) (int64, error)
	GetMaxQueueID(tx *sql.Tx) (int, error)
	InsertQueue(tx *sql.Tx, patientID int64, queueNumber, complaint, paymentMethod, paymentStatus, queueStatus string, paidAt *time.Time) (int64, error)
	UpdatePaymentStatus(queueID int64, paymentStatus string, paidAt time.Time) error
	UpdateQueueStatus(queueID int64, status string) error
	BeginTx() (*sql.Tx, error)
}

type sqlRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *sqlRepository) GetQueues() ([]domain.QueueItem, error) {
	query := `
		SELECT q.id, q.patient_id, p.name, p.nik, p.phone, coalesce(p.ihs_number, ''),
		       q.queue_number, coalesce(q.complaint, ''), q.payment_method, q.payment_status,
		       q.queue_status, q.paid_at, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		ORDER BY 
			CASE WHEN q.queue_status = 'CALLING' THEN 1
			     WHEN q.queue_status = 'IN_EXAM' THEN 2
			     WHEN q.queue_status = 'WAITING' THEN 3
			     ELSE 4 END,
			q.id ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query queues: %w", err)
	}
	defer rows.Close()

	var items []domain.QueueItem
	for rows.Next() {
		var q domain.QueueItem
		var paidAt sql.NullTime
		if err := rows.Scan(
			&q.ID, &q.PatientID, &q.PatientName, &q.PatientNIK, &q.PatientPhone, &q.PatientIHS,
			&q.QueueNumber, &q.Complaint, &q.PaymentMethod, &q.PaymentStatus,
			&q.QueueStatus, &paidAt, &q.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan queue item: %w", err)
		}
		if paidAt.Valid {
			q.PaidAt = &paidAt.Time
		}
		items = append(items, q)
	}
	return items, nil
}

func (r *sqlRepository) GetQueueByID(id int64) (*domain.QueueItem, error) {
	query := `
		SELECT q.id, q.patient_id, p.name, p.nik, p.phone, coalesce(p.ihs_number, ''),
		       q.queue_number, coalesce(q.complaint, ''), q.payment_method, q.payment_status,
		       q.queue_status, q.paid_at, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE q.id = ?
	`
	var q domain.QueueItem
	var paidAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&q.ID, &q.PatientID, &q.PatientName, &q.PatientNIK, &q.PatientPhone, &q.PatientIHS,
		&q.QueueNumber, &q.Complaint, &q.PaymentMethod, &q.PaymentStatus,
		&q.QueueStatus, &paidAt, &q.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get queue by id %d: %w", id, err)
	}
	if paidAt.Valid {
		q.PaidAt = &paidAt.Time
	}
	return &q, nil
}

func (r *sqlRepository) FindPatientByNIK(tx *sql.Tx, nik string) (int64, string, error) {
	var patientID int64
	var ihsNumber string
	query := "SELECT id, coalesce(ihs_number, '') FROM patients WHERE nik = ?"
	var err error
	if tx != nil {
		err = tx.QueryRow(query, nik).Scan(&patientID, &ihsNumber)
	} else {
		err = r.db.QueryRow(query, nik).Scan(&patientID, &ihsNumber)
	}
	return patientID, ihsNumber, err
}

func (r *sqlRepository) CreatePatient(tx *sql.Tx, nik, name, phone, birthDate, gender, ihs string) (int64, error) {
	query := `
		INSERT INTO patients (nik, name, phone, birth_date, gender, ihs_number)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, nik, name, phone, birthDate, gender, ihs)
	} else {
		res, err = r.db.Exec(query, nik, name, phone, birthDate, gender, ihs)
	}
	if err != nil {
		return 0, fmt.Errorf("insert patient: %w", err)
	}
	return res.LastInsertId()
}

func (r *sqlRepository) GetMaxQueueID(tx *sql.Tx) (int, error) {
	query := "SELECT coalesce(max(id), 0) FROM queues"
	var maxID int
	var err error
	if tx != nil {
		err = tx.QueryRow(query).Scan(&maxID)
	} else {
		err = r.db.QueryRow(query).Scan(&maxID)
	}
	return maxID, err
}

func (r *sqlRepository) InsertQueue(tx *sql.Tx, patientID int64, queueNumber, complaint, paymentMethod, paymentStatus, queueStatus string, paidAt *time.Time) (int64, error) {
	query := `
		INSERT INTO queues (patient_id, queue_number, complaint, payment_method, payment_status, queue_status, paid_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.Exec(query, patientID, queueNumber, complaint, paymentMethod, paymentStatus, queueStatus, paidAt)
	} else {
		res, err = r.db.Exec(query, patientID, queueNumber, complaint, paymentMethod, paymentStatus, queueStatus, paidAt)
	}
	if err != nil {
		return 0, fmt.Errorf("insert queue: %w", err)
	}
	return res.LastInsertId()
}

func (r *sqlRepository) UpdatePaymentStatus(queueID int64, paymentStatus string, paidAt time.Time) error {
	query := `
		UPDATE queues
		SET payment_status = ?, paid_at = ?
		WHERE id = ?
	`
	res, err := r.db.Exec(query, paymentStatus, paidAt, queueID)
	if err != nil {
		return fmt.Errorf("update queue payment status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("queue not found or not modified")
	}
	return nil
}

func (r *sqlRepository) UpdateQueueStatus(queueID int64, status string) error {
	query := `
		UPDATE queues
		SET queue_status = ?
		WHERE id = ?
	`
	res, err := r.db.Exec(query, status, queueID)
	if err != nil {
		return fmt.Errorf("update queue status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("queue not found or not modified")
	}
	return nil
}

