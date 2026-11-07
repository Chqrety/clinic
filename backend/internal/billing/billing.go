package billing

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
	r.Get("/invoices", s.HandleGetInvoices)
	r.Get("/invoices/queue/{queueId}", s.HandleGetByQueue)
}

func (s *Service) HandleGetInvoices(w http.ResponseWriter, r *http.Request) {
	list, err := s.GetInvoices()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: list})
}

func (s *Service) HandleGetByQueue(w http.ResponseWriter, r *http.Request) {
	qStr := chi.URLParam(r, "queueId")
	queueID, err := strconv.ParseInt(qStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	inv, err := s.GetOrCreateInvoice(queueID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: inv})
}

func (s *Service) GetInvoices() ([]domain.Invoice, error) {
	// Auto-generate invoice for any completed queue missing from invoices table
	cRows, err := s.db.Query(`
		SELECT q.id FROM queues q
		WHERE q.queue_status = 'COMPLETED'
		  AND q.id NOT IN (SELECT queue_id FROM invoices)
	`)
	if err == nil {
		var qIDs []int64
		for cRows.Next() {
			var qid int64
			if err := cRows.Scan(&qid); err == nil {
				qIDs = append(qIDs, qid)
			}
		}
		cRows.Close()
		for _, qid := range qIDs {
			_, _ = s.GetOrCreateInvoice(qid)
		}
	}

	query := `
		SELECT i.id, i.queue_id, i.invoice_number, pat.name, pat.nik,
		       i.registration_fee, i.consultation_fee, i.medication_fee, i.total_amount,
		       i.payment_method, i.payment_status, i.created_at
		FROM invoices i
		JOIN queues q ON i.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		ORDER BY i.id DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.PatientName, &inv.PatientNIK,
			&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
			&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}
		// Calculate internal breakdown
		inv.DoctorShare = int(float64(inv.ConsultationFee) * 0.70)
		inv.TaxShare = int(float64(inv.TotalAmount) * 0.11)
		inv.ClinicShare = inv.TotalAmount - inv.DoctorShare - inv.TaxShare

		list = append(list, inv)
	}
	return list, nil
}

func (s *Service) GetOrCreateInvoice(queueID int64) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := s.db.QueryRow(`
		SELECT i.id, i.queue_id, i.invoice_number, pat.name, pat.nik,
		       i.registration_fee, i.consultation_fee, i.medication_fee, i.total_amount,
		       i.payment_method, i.payment_status, i.created_at
		FROM invoices i
		JOIN queues q ON i.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		WHERE i.queue_id = ?
	`, queueID).Scan(
		&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.PatientName, &inv.PatientNIK,
		&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
		&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
	)

	if err == nil {
		// Update medication fee if prescription was dispensed after invoice init
		var currentMedFee int
		_ = s.db.QueryRow(`
			SELECT coalesce(sum(pi.subtotal), 0)
			FROM prescription_items pi
			JOIN prescriptions p ON pi.prescription_id = p.id
			JOIN consultations c ON p.consultation_id = c.id
			WHERE c.queue_id = ? AND p.status != 'CANCELLED'
		`, queueID).Scan(&currentMedFee)

		if currentMedFee != inv.MedicationFee {
			inv.MedicationFee = currentMedFee
			inv.TotalAmount = inv.RegistrationFee + inv.ConsultationFee + currentMedFee
			_, _ = s.db.Exec(`UPDATE invoices SET medication_fee = ?, total_amount = ? WHERE id = ?`,
				inv.MedicationFee, inv.TotalAmount, inv.ID)
		}

		inv.DoctorShare = int(float64(inv.ConsultationFee) * 0.70)
		inv.TaxShare = int(float64(inv.TotalAmount) * 0.11)
		inv.ClinicShare = inv.TotalAmount - inv.DoctorShare - inv.TaxShare
		return &inv, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Generate new invoice
	var patName, patNIK, payMethod string
	err = s.db.QueryRow(`
		SELECT p.name, p.nik, q.payment_method
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE q.id = ?
	`, queueID).Scan(&patName, &patNIK, &payMethod)
	if err != nil {
		return nil, fmt.Errorf("queue not found: %w", err)
	}

	// Calculate fees
	regFee := 35000
	consultFee := 50000

	// Sum dispensed medication if any
	medFee := 0
	_ = s.db.QueryRow(`
		SELECT coalesce(sum(pi.subtotal), 0)
		FROM prescription_items pi
		JOIN prescriptions p ON pi.prescription_id = p.id
		JOIN consultations c ON p.consultation_id = c.id
		WHERE c.queue_id = ? AND p.status != 'CANCELLED'
	`, queueID).Scan(&medFee)

	total := regFee + consultFee + medFee
	invNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), queueID)

	res, err := s.db.Exec(`
		INSERT INTO invoices (queue_id, invoice_number, registration_fee, consultation_fee, medication_fee, total_amount, payment_method, payment_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'PAID')
	`, queueID, invNumber, regFee, consultFee, medFee, total, payMethod)
	if err != nil {
		return nil, err
	}

	invID, _ := res.LastInsertId()
	inv = domain.Invoice{
		ID:              invID,
		QueueID:         queueID,
		InvoiceNumber:   invNumber,
		PatientName:     patName,
		PatientNIK:      patNIK,
		RegistrationFee: regFee,
		ConsultationFee: consultFee,
		MedicationFee:   medFee,
		TotalAmount:     total,
		PaymentMethod:   payMethod,
		PaymentStatus:   "PAID",
		DoctorShare:     int(float64(consultFee) * 0.70),
		TaxShare:        int(float64(total) * 0.11),
		ClinicShare:     total - int(float64(consultFee)*0.70) - int(float64(total)*0.11),
		CreatedAt:       time.Now(),
	}
	return &inv, nil
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}
