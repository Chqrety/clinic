package queue

import (
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/domain"
)

type CreateQueueRequest struct {
	NIK           string `json:"nik"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	BirthDate     string `json:"birth_date"`
	Gender        string `json:"gender"`
	Complaint     string `json:"complaint"`
	PaymentMethod string `json:"payment_method"` // "QRIS" or "CASH"
}

type UpdateStatusRequest struct {
	Status string `json:"status"` // "WAITING", "CALLING", "IN_EXAM", "COMPLETED", "CANCELLED"
}

type Service struct {
	repo Repository
}

// NewService creates a queue Service using standard database/sql.DB.
func NewService(db *sql.DB) *Service {
	return &Service{repo: NewSQLRepository(db)}
}

// NewServiceWithRepo creates a queue Service using a custom Repository.
func NewServiceWithRepo(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetQueues() ([]domain.QueueItem, error) {
	return s.repo.GetQueues()
}

func (s *Service) GetQueueByID(id int64) (*domain.QueueItem, error) {
	return s.repo.GetQueueByID(id)
}

func (s *Service) CreateQueue(req CreateQueueRequest) (*domain.QueueItem, error) {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Get or create patient record
	patientID, ihsNumber, err := s.repo.FindPatientByNIK(tx, req.NIK)
	if err == sql.ErrNoRows {
		// Mock SATUSEHAT IHS identifier: P + timestamp
		ihsNumber = fmt.Sprintf("P%d", time.Now().UnixNano()%100000000000)
		patientID, err = s.repo.CreatePatient(tx, req.NIK, req.Name, req.Phone, req.BirthDate, req.Gender, ihsNumber)
		if err != nil {
			return nil, fmt.Errorf("create patient: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("find patient by NIK: %w", err)
	}

	// 2. Generate next sequential queue number (A-001, A-002, ...)
	maxID, err := s.repo.GetMaxQueueID(tx)
	if err != nil {
		maxID = 0
	}
	queueNumber := fmt.Sprintf("A-%03d", maxID+1)

	// Payment status & initial queue status
	paymentStatus := "PENDING"
	queueStatus := "WAITING"
	var paidAt *time.Time

	if req.PaymentMethod == "QRIS" {
		// Instant mock payment
		paymentStatus = "PAID"
		now := time.Now()
		paidAt = &now
	}

	queueID, err := s.repo.InsertQueue(tx, patientID, queueNumber, req.Complaint, req.PaymentMethod, paymentStatus, queueStatus, paidAt)
	if err != nil {
		return nil, fmt.Errorf("insert queue: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit queue transaction: %w", err)
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
		PaidAt:        paidAt,
		CreatedAt:     time.Now(),
	}, nil
}

func (s *Service) ConfirmPayment(queueID int64) error {
	now := time.Now()
	return s.repo.UpdatePaymentStatus(queueID, "PAID", now)
}

func (s *Service) UpdateStatus(queueID int64, status string) error {
	return s.repo.UpdateQueueStatus(queueID, status)
}

