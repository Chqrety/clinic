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

	invoices, err := s.GetInvoicesByQueue(queueID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: invoices})
}

func (s *Service) GetInvoices() ([]domain.Invoice, error) {
	// 1. Auto-create REGISTRATION invoices for paid queues without registration invoice
	regQRows, err := s.db.Query(`
		SELECT q.id FROM queues q
		WHERE q.payment_status = 'PAID'
		  AND q.id NOT IN (SELECT queue_id FROM invoices WHERE invoice_type = 'REGISTRATION')
	`)
	if err == nil {
		var qIDs []int64
		for regQRows.Next() {
			var qid int64
			if err := regQRows.Scan(&qid); err == nil {
				qIDs = append(qIDs, qid)
			}
		}
		regQRows.Close()
		for _, qid := range qIDs {
			_, _ = s.GetOrCreateRegistrationInvoice(qid)
		}
	}

	// 2. Auto-create PHARMACY invoices for dispensed prescriptions without pharmacy invoice
	pharmRows, err := s.db.Query(`
		SELECT c.queue_id FROM prescriptions p
		JOIN consultations c ON p.consultation_id = c.id
		WHERE p.status = 'DISPENSED'
		  AND c.queue_id NOT IN (SELECT queue_id FROM invoices WHERE invoice_type = 'PHARMACY')
	`)
	if err == nil {
		var qIDs []int64
		for pharmRows.Next() {
			var qid int64
			if err := pharmRows.Scan(&qid); err == nil {
				qIDs = append(qIDs, qid)
			}
		}
		pharmRows.Close()
		for _, qid := range qIDs {
			_, _ = s.GetOrCreatePharmacyInvoice(qid)
		}
	}

	query := `
		SELECT i.id, i.queue_id, i.invoice_number, coalesce(i.invoice_type, 'REGISTRATION'), pat.name, pat.nik,
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
			&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.InvoiceType, &inv.PatientName, &inv.PatientNIK,
			&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
			&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}

		s.calculateShares(&inv)
		if inv.InvoiceType == "PHARMACY" {
			inv.Items, _ = s.getPrescriptionItemsForQueue(inv.QueueID)
		}

		list = append(list, inv)
	}
	return list, nil
}

func (s *Service) GetInvoicesByQueue(queueID int64) ([]domain.Invoice, error) {
	// Ensure registration invoice exists if paid
	var payStatus string
	_ = s.db.QueryRow("SELECT payment_status FROM queues WHERE id = ?", queueID).Scan(&payStatus)
	if payStatus == "PAID" {
		_, _ = s.GetOrCreateRegistrationInvoice(queueID)
	}

	// Ensure pharmacy invoice exists if dispensed
	var hasDispensed int
	_ = s.db.QueryRow(`
		SELECT count(*) FROM prescriptions p
		JOIN consultations c ON p.consultation_id = c.id
		WHERE c.queue_id = ? AND p.status = 'DISPENSED'
	`, queueID).Scan(&hasDispensed)
	if hasDispensed > 0 {
		_, _ = s.GetOrCreatePharmacyInvoice(queueID)
	}

	query := `
		SELECT i.id, i.queue_id, i.invoice_number, coalesce(i.invoice_type, 'REGISTRATION'), pat.name, pat.nik,
		       i.registration_fee, i.consultation_fee, i.medication_fee, i.total_amount,
		       i.payment_method, i.payment_status, i.created_at
		FROM invoices i
		JOIN queues q ON i.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		WHERE i.queue_id = ?
		ORDER BY i.id ASC
	`
	rows, err := s.db.Query(query, queueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		if err := rows.Scan(
			&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.InvoiceType, &inv.PatientName, &inv.PatientNIK,
			&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
			&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
		); err != nil {
			return nil, err
		}
		s.calculateShares(&inv)
		if inv.InvoiceType == "PHARMACY" {
			inv.Items, _ = s.getPrescriptionItemsForQueue(inv.QueueID)
		}
		list = append(list, inv)
	}
	return list, nil
}

func (s *Service) GetOrCreateRegistrationInvoice(queueID int64) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := s.db.QueryRow(`
		SELECT i.id, i.queue_id, i.invoice_number, coalesce(i.invoice_type, 'REGISTRATION'), pat.name, pat.nik,
		       i.registration_fee, i.consultation_fee, i.medication_fee, i.total_amount,
		       i.payment_method, i.payment_status, i.created_at
		FROM invoices i
		JOIN queues q ON i.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		WHERE i.queue_id = ? AND i.invoice_type = 'REGISTRATION'
	`, queueID).Scan(
		&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.InvoiceType, &inv.PatientName, &inv.PatientNIK,
		&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
		&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
	)
	if err == nil {
		s.calculateShares(&inv)
		return &inv, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

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

	regFee := 35000
	invNumber := fmt.Sprintf("INV-REG-%s-%04d", time.Now().Format("20060102"), queueID)

	res, err := s.db.Exec(`
		INSERT INTO invoices (queue_id, invoice_number, invoice_type, registration_fee, consultation_fee, medication_fee, total_amount, payment_method, payment_status)
		VALUES (?, ?, 'REGISTRATION', ?, 0, 0, ?, ?, 'PAID')
	`, queueID, invNumber, regFee, regFee, payMethod)
	if err != nil {
		return nil, err
	}

	invID, _ := res.LastInsertId()
	inv = domain.Invoice{
		ID:              invID,
		QueueID:         queueID,
		InvoiceNumber:   invNumber,
		InvoiceType:     "REGISTRATION",
		PatientName:     patName,
		PatientNIK:      patNIK,
		RegistrationFee: regFee,
		ConsultationFee: 0,
		MedicationFee:   0,
		TotalAmount:     regFee,
		PaymentMethod:   payMethod,
		PaymentStatus:   "PAID",
		CreatedAt:       time.Now(),
	}
	s.calculateShares(&inv)
	return &inv, nil
}

func (s *Service) GetOrCreatePharmacyInvoice(queueID int64) (*domain.Invoice, error) {
	var inv domain.Invoice
	err := s.db.QueryRow(`
		SELECT i.id, i.queue_id, i.invoice_number, coalesce(i.invoice_type, 'PHARMACY'), pat.name, pat.nik,
		       i.registration_fee, i.consultation_fee, i.medication_fee, i.total_amount,
		       i.payment_method, i.payment_status, i.created_at
		FROM invoices i
		JOIN queues q ON i.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		WHERE i.queue_id = ? AND i.invoice_type = 'PHARMACY'
	`, queueID).Scan(
		&inv.ID, &inv.QueueID, &inv.InvoiceNumber, &inv.InvoiceType, &inv.PatientName, &inv.PatientNIK,
		&inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee, &inv.TotalAmount,
		&inv.PaymentMethod, &inv.PaymentStatus, &inv.CreatedAt,
	)
	if err == nil {
		s.calculateShares(&inv)
		inv.Items, _ = s.getPrescriptionItemsForQueue(queueID)
		return &inv, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	var medFee int
	_ = s.db.QueryRow(`
		SELECT coalesce(sum(pi.subtotal), 0)
		FROM prescription_items pi
		JOIN prescriptions p ON pi.prescription_id = p.id
		JOIN consultations c ON p.consultation_id = c.id
		WHERE c.queue_id = ? AND p.status != 'CANCELLED'
	`, queueID).Scan(&medFee)

	if medFee <= 0 {
		return nil, nil // No medication prescribed
	}

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

	invNumber := fmt.Sprintf("INV-FAR-%s-%04d", time.Now().Format("20060102"), queueID)

	res, err := s.db.Exec(`
		INSERT INTO invoices (queue_id, invoice_number, invoice_type, registration_fee, consultation_fee, medication_fee, total_amount, payment_method, payment_status)
		VALUES (?, ?, 'PHARMACY', 0, 0, ?, ?, ?, 'PAID')
	`, queueID, invNumber, medFee, medFee, payMethod)
	if err != nil {
		return nil, err
	}

	invID, _ := res.LastInsertId()
	items, _ := s.getPrescriptionItemsForQueue(queueID)
	inv = domain.Invoice{
		ID:              invID,
		QueueID:         queueID,
		InvoiceNumber:   invNumber,
		InvoiceType:     "PHARMACY",
		PatientName:     patName,
		PatientNIK:      patNIK,
		RegistrationFee: 0,
		ConsultationFee: 0,
		MedicationFee:   medFee,
		TotalAmount:     medFee,
		PaymentMethod:   payMethod,
		PaymentStatus:   "PAID",
		Items:           items,
		CreatedAt:       time.Now(),
	}
	s.calculateShares(&inv)
	return &inv, nil
}

func (s *Service) calculateShares(inv *domain.Invoice) {
	if inv.InvoiceType == "PHARMACY" {
		inv.DoctorShare = 0
		inv.TaxShare = int(float64(inv.TotalAmount) * 0.11)
		inv.ClinicShare = inv.TotalAmount - inv.TaxShare
	} else {
		// REGISTRATION invoice: konsultasi (jika ada) 70% ke dokter, sisanya operasional klinik
		inv.DoctorShare = int(float64(inv.ConsultationFee) * 0.70)
		inv.TaxShare = int(float64(inv.TotalAmount) * 0.11)
		inv.ClinicShare = inv.TotalAmount - inv.DoctorShare - inv.TaxShare
	}
}

func (s *Service) getPrescriptionItemsForQueue(queueID int64) ([]domain.PrescriptionItem, error) {
	rows, err := s.db.Query(`
		SELECT pi.id, pi.prescription_id, pi.drug_id, d.name, pi.batch_id, coalesce(b.batch_number, '-'),
		       pi.quantity, pi.dosage_instruction, pi.unit_price, pi.subtotal
		FROM prescription_items pi
		JOIN prescriptions p ON pi.prescription_id = p.id
		JOIN consultations c ON p.consultation_id = c.id
		JOIN drugs d ON pi.drug_id = d.id
		LEFT JOIN drug_batches b ON pi.batch_id = b.id
		WHERE c.queue_id = ? AND p.status != 'CANCELLED'
		ORDER BY pi.id ASC
	`, queueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.PrescriptionItem
	for rows.Next() {
		var it domain.PrescriptionItem
		var bID sql.NullInt64
		if err := rows.Scan(
			&it.ID, &it.PrescriptionID, &it.DrugID, &it.DrugName, &bID, &it.BatchNumber,
			&it.Quantity, &it.DosageInstruction, &it.UnitPrice, &it.Subtotal,
		); err != nil {
			return nil, err
		}
		if bID.Valid {
			it.BatchID = &bID.Int64
		}
		items = append(items, it)
	}
	return items, nil
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}
