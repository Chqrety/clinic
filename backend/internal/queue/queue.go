package queue

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"clinic-api/internal/domain"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func RegisterRoutes(r chi.Router, s *Service) {
	r.Get("/queues", s.HandleGetQueues)
	r.Post("/queues", s.HandleCreateQueue)
	r.Post("/queues/{id}/pay", s.HandlePayQueue)
	r.Post("/queues/{id}/status", s.HandleUpdateStatus)
}

func (s *Service) HandleGetQueues(w http.ResponseWriter, r *http.Request) {
	queues, err := s.GetQueues()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: queues})
}

type CreateQueueRequest struct {
	NIK           string `json:"nik"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	BirthDate     string `json:"birth_date"`
	Gender        string `json:"gender"`
	Complaint     string `json:"complaint"`
	PaymentMethod string `json:"payment_method"` // "QRIS" or "CASH"
}

func (s *Service) HandleCreateQueue(w http.ResponseWriter, r *http.Request) {
	var req CreateQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.NIK == "" || req.Name == "" {
		sendError(w, http.StatusBadRequest, "NIK dan Nama wajib diisi")
		return
	}

	queueItem, err := s.CreateQueue(req)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, domain.ApiResponse{
		Success: true,
		Data:    queueItem,
		Message: "Antrian berhasil didaftarkan",
	})
}

func (s *Service) HandlePayQueue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	if err := s.ConfirmPayment(id); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Pembayaran terverifikasi, nomor antrian aktif",
	})
}

type UpdateStatusRequest struct {
	Status string `json:"status"` // "WAITING", "CALLING", "IN_EXAM", "COMPLETED", "CANCELLED"
}

func (s *Service) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.UpdateStatus(id, req.Status); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: fmt.Sprintf("Status antrian diubah menjadi %s", req.Status),
	})
}

// Service Logic & DB
func (s *Service) GetQueues() ([]domain.QueueItem, error) {
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
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if paidAt.Valid {
			q.PaidAt = &paidAt.Time
		}
		items = append(items, q)
	}
	return items, nil
}

func (s *Service) CreateQueue(req CreateQueueRequest) (*domain.QueueItem, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Get or create patient
	var patientID int64
	var ihsNumber string
	err = tx.QueryRow("SELECT id, coalesce(ihs_number, '') FROM patients WHERE nik = ?", req.NIK).Scan(&patientID, &ihsNumber)
	if err == sql.ErrNoRows {
		// Mock SATUSEHAT IHS identifier: P + last 11 digits of NIK or timestamp
		ihsNumber = fmt.Sprintf("P%d", time.Now().UnixNano()%100000000000)
		res, err := tx.Exec(`
			INSERT INTO patients (nik, name, phone, birth_date, gender, ihs_number)
			VALUES (?, ?, ?, ?, ?, ?)
		`, req.NIK, req.Name, req.Phone, req.BirthDate, req.Gender, ihsNumber)
		if err != nil {
			return nil, err
		}
		patientID, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// 2. Next queue number (A-001, A-002, ...)
	var maxID int
	_ = tx.QueryRow("SELECT coalesce(max(id), 0) FROM queues").Scan(&maxID)
	queueNumber := fmt.Sprintf("A-%03d", maxID+1)

	// Payment status & queue status
	paymentStatus := "PENDING"
	queueStatus := "WAITING"
	var paidAt interface{} = nil

	if req.PaymentMethod == "QRIS" {
		// Mock instant QRIS payment success
		paymentStatus = "PAID"
		paidAt = time.Now()
	}

	res, err := tx.Exec(`
		INSERT INTO queues (patient_id, queue_number, complaint, payment_method, payment_status, queue_status, paid_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, patientID, queueNumber, req.Complaint, req.PaymentMethod, paymentStatus, queueStatus, paidAt)
	if err != nil {
		return nil, err
	}

	queueID, _ := res.LastInsertId()
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var paidAtTime *time.Time
	if req.PaymentMethod == "QRIS" {
		t := time.Now()
		paidAtTime = &t
	}

	return &domain.QueueItem{
		ID:            queueID,
		PatientID:     patientID,
		PatientName:   req.Name,
		PatientNIK:    req.NIK,
		PatientPhone:  req.Phone,
		PatientIHS:    ihsNumber,
		QueueNumber:   queueNumber,
		Complaint:     req.Complaint,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: paymentStatus,
		QueueStatus:   queueStatus,
		PaidAt:        paidAtTime,
		CreatedAt:     time.Now(),
	}, nil
}

func (s *Service) ConfirmPayment(queueID int64) error {
	_, err := s.db.Exec(`
		UPDATE queues
		SET payment_status = 'PAID', paid_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, queueID)
	return err
}

func (s *Service) UpdateStatus(queueID int64, status string) error {
	_, err := s.db.Exec(`
		UPDATE queues
		SET queue_status = ?
		WHERE id = ?
	`, status, queueID)
	return err
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}
