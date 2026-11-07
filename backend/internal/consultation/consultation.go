package consultation

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
	r.Get("/consultations/queue/{queueId}", s.HandleGetByQueue)
	r.Post("/consultations", s.HandleSaveConsultation)
	r.Get("/prescriptions", s.HandleGetPrescriptions)
	r.Post("/prescriptions/{id}/dispense", s.HandleDispense)
	r.Post("/prescriptions/{id}/cancel", s.HandleCancel)
}

func (s *Service) HandleGetByQueue(w http.ResponseWriter, r *http.Request) {
	qStr := chi.URLParam(r, "queueId")
	queueID, err := strconv.ParseInt(qStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	cons, err := s.GetOrCreateConsultation(queueID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: cons})
}

type PrescriptionItemInput struct {
	DrugID            int64  `json:"drug_id"`
	Quantity          int    `json:"quantity"`
	DosageInstruction string `json:"dosage_instruction"`
}

type SaveConsultationRequest struct {
	QueueID      int64                   `json:"queue_id"`
	DoctorName   string                  `json:"doctor_name"`
	Anamnesis    string                  `json:"anamnesis"`
	Diagnosis    string                  `json:"diagnosis"`
	ICD10Code    string                  `json:"icd10_code"`
	Prescription []PrescriptionItemInput `json:"prescription"`
}

func (s *Service) HandleSaveConsultation(w http.ResponseWriter, r *http.Request) {
	var req SaveConsultationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cons, err := s.SaveConsultation(req)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Data:    cons,
		Message: "Pemeriksaan dan E-Resep berhasil disimpan (Stok obat di-reserve)",
	})
}

func (s *Service) HandleGetPrescriptions(w http.ResponseWriter, r *http.Request) {
	list, err := s.GetPrescriptions()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: list})
}

func (s *Service) HandleDispense(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid prescription ID")
		return
	}

	if err := s.DispensePrescription(id); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Obat berhasil diserahkan ke pasien, stok fisik resmi terpotong",
	})
}

func (s *Service) HandleCancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid prescription ID")
		return
	}

	if err := s.CancelPrescription(id); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Resep dibatalkan, reserved stock telah dikembalikan",
	})
}

// Logic & DB Implementations
func (s *Service) GetOrCreateConsultation(queueID int64) (*domain.Consultation, error) {
	var c domain.Consultation
	err := s.db.QueryRow(`
		SELECT id, queue_id, doctor_name, coalesce(anamnesis, ''), coalesce(diagnosis, ''), coalesce(icd10_code, ''), status, created_at
		FROM consultations WHERE queue_id = ?
	`, queueID).Scan(&c.ID, &c.QueueID, &c.DoctorName, &c.Anamnesis, &c.Diagnosis, &c.ICD10Code, &c.Status, &c.CreatedAt)

	if err == sql.ErrNoRows {
		// create draft
		docName := "dr. Hendra Pratama, Sp.PD"
		res, err := s.db.Exec(`
			INSERT INTO consultations (queue_id, doctor_name, status)
			VALUES (?, ?, 'DRAFT')
		`, queueID, docName)
		if err != nil {
			return nil, err
		}
		c.ID, _ = res.LastInsertId()
		c.QueueID = queueID
		c.DoctorName = docName
		c.Status = "DRAFT"
		return &c, nil
	} else if err != nil {
		return nil, err
	}

	// load prescription if any
	c.Prescription, _ = s.getPrescriptionByConsultation(c.ID)
	return &c, nil
}

func (s *Service) SaveConsultation(req SaveConsultationRequest) (*domain.Consultation, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Get or create consultation
	var consID int64
	err = tx.QueryRow("SELECT id FROM consultations WHERE queue_id = ?", req.QueueID).Scan(&consID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec(`
			INSERT INTO consultations (queue_id, doctor_name, anamnesis, diagnosis, icd10_code, status)
			VALUES (?, ?, ?, ?, ?, 'COMPLETED')
		`, req.QueueID, req.DoctorName, req.Anamnesis, req.Diagnosis, req.ICD10Code)
		if err != nil {
			return nil, err
		}
		consID, _ = res.LastInsertId()
	} else if err != nil {
		return nil, err
	} else {
		_, err = tx.Exec(`
			UPDATE consultations
			SET doctor_name = ?, anamnesis = ?, diagnosis = ?, icd10_code = ?, status = 'COMPLETED'
			WHERE id = ?
		`, req.DoctorName, req.Anamnesis, req.Diagnosis, req.ICD10Code, consID)
		if err != nil {
			return nil, err
		}
	}

	// 2. Process e-prescription if items provided
	if len(req.Prescription) > 0 {
		// Create prescription record
		res, err := tx.Exec(`
			INSERT INTO prescriptions (consultation_id, status)
			VALUES (?, 'RESERVED')
		`, consID)
		if err != nil {
			return nil, err
		}
		rxID, _ := res.LastInsertId()

		// Validate & Reserve stock batch by batch (FEFO)
		for _, item := range req.Prescription {
			var drugName string
			var unitPrice int
			err := tx.QueryRow("SELECT name, sell_price FROM drugs WHERE id = ?", item.DrugID).Scan(&drugName, &unitPrice)
			if err != nil {
				return nil, fmt.Errorf("obat id %d tidak ditemukan", item.DrugID)
			}

			// Check total available stock across batches
			var totalAvail int
			err = tx.QueryRow(`
				SELECT coalesce(sum(physical_stock - reserved_stock), 0)
				FROM drug_batches WHERE drug_id = ?
			`, item.DrugID).Scan(&totalAvail)
			if err != nil {
				return nil, err
			}

			if totalAvail < item.Quantity {
				return nil, fmt.Errorf("stok %s tidak mencukupi (tersedia: %d, diminta: %d)", drugName, totalAvail, item.Quantity)
			}

			// Pick batch with earliest expiry (FEFO) that has available stock
			var batchID int64
			var batchNum string
			err = tx.QueryRow(`
				SELECT id, batch_number
				FROM drug_batches
				WHERE drug_id = ? AND (physical_stock - reserved_stock) >= ?
				ORDER BY expired_date ASC
				LIMIT 1
			`, item.DrugID, item.Quantity).Scan(&batchID, &batchNum)

			if err != nil {
				// Fallback to any batch with remaining available stock
				err = tx.QueryRow(`
					SELECT id, batch_number
					FROM drug_batches
					WHERE drug_id = ? AND (physical_stock - reserved_stock) > 0
					ORDER BY expired_date ASC
					LIMIT 1
				`, item.DrugID).Scan(&batchID, &batchNum)
				if err != nil {
					return nil, fmt.Errorf("tidak ada batch aktif untuk %s", drugName)
				}
			}

			// Atomic update: increment reserved_stock
			_, err = tx.Exec(`
				UPDATE drug_batches
				SET reserved_stock = reserved_stock + ?
				WHERE id = ?
			`, item.Quantity, batchID)
			if err != nil {
				return nil, err
			}

			subtotal := unitPrice * item.Quantity
			_, err = tx.Exec(`
				INSERT INTO prescription_items (prescription_id, drug_id, batch_id, quantity, dosage_instruction, unit_price, subtotal)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, rxID, item.DrugID, batchID, item.Quantity, item.DosageInstruction, unitPrice, subtotal)
			if err != nil {
				return nil, err
			}
		}
	}

	// Update queue status: if prescription exists, wait in pharmacy; else complete consultation
	if len(req.Prescription) > 0 {
		_, _ = tx.Exec("UPDATE queues SET queue_status = 'IN_EXAM' WHERE id = ?", req.QueueID)
	} else {
		_, _ = tx.Exec("UPDATE queues SET queue_status = 'COMPLETED' WHERE id = ?", req.QueueID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetOrCreateConsultation(req.QueueID)
}

func (s *Service) GetPrescriptions() ([]domain.Prescription, error) {
	query := `
		SELECT p.id, p.consultation_id, c.queue_id, pat.name, pat.nik, c.doctor_name, p.status, p.created_at
		FROM prescriptions p
		JOIN consultations c ON p.consultation_id = c.id
		JOIN queues q ON c.queue_id = q.id
		JOIN patients pat ON q.patient_id = pat.id
		ORDER BY 
			CASE WHEN p.status = 'RESERVED' THEN 1
			     WHEN p.status = 'PREPARING' THEN 2
			     ELSE 3 END,
			p.id DESC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Prescription
	for rows.Next() {
		var rx domain.Prescription
		if err := rows.Scan(
			&rx.ID, &rx.ConsultationID, &rx.QueueID, &rx.PatientName, &rx.PatientNIK,
			&rx.DoctorName, &rx.Status, &rx.CreatedAt,
		); err != nil {
			return nil, err
		}

		// fetch items
		items, total, _ := s.getPrescriptionItems(rx.ID)
		rx.Items = items
		rx.TotalAmount = total
		list = append(list, rx)
	}
	return list, nil
}

func (s *Service) DispensePrescription(rxID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get items
	rows, err := tx.Query(`
		SELECT batch_id, quantity FROM prescription_items WHERE prescription_id = ?
	`, rxID)
	if err != nil {
		return err
	}
	defer rows.Close()

	type itemD struct {
		batchID sql.NullInt64
		qty     int
	}
	var items []itemD
	for rows.Next() {
		var i itemD
		if err := rows.Scan(&i.batchID, &i.qty); err != nil {
			return err
		}
		items = append(items, i)
	}
	rows.Close()

	// Commit deduction: physical_stock = physical_stock - qty, reserved_stock = reserved_stock - qty
	for _, it := range items {
		if it.batchID.Valid {
			_, err = tx.Exec(`
				UPDATE drug_batches
				SET physical_stock = physical_stock - ?,
				    reserved_stock = reserved_stock - ?
				WHERE id = ?
			`, it.qty, it.qty, it.batchID.Int64)
			if err != nil {
				return err
			}
		}
	}

	// Set status to DISPENSED
	_, err = tx.Exec("UPDATE prescriptions SET status = 'DISPENSED' WHERE id = ?", rxID)
	if err != nil {
		return err
	}

	// Update queue status to COMPLETED
	var queueID int64
	_ = tx.QueryRow(`
		SELECT c.queue_id FROM consultations c
		JOIN prescriptions p ON p.consultation_id = c.id
		WHERE p.id = ?
	`, rxID).Scan(&queueID)

	if queueID > 0 {
		_, _ = tx.Exec("UPDATE queues SET queue_status = 'COMPLETED' WHERE id = ?", queueID)

		// Auto-generate or update invoice in billing
		var hasInv int
		_ = tx.QueryRow("SELECT count(*) FROM invoices WHERE queue_id = ?", queueID).Scan(&hasInv)

		var medFee int
		_ = tx.QueryRow(`
			SELECT coalesce(sum(pi.subtotal), 0)
			FROM prescription_items pi
			JOIN prescriptions p ON pi.prescription_id = p.id
			JOIN consultations c ON p.consultation_id = c.id
			WHERE c.queue_id = ? AND p.status != 'CANCELLED'
		`, queueID).Scan(&medFee)

		regFee := 35000
		consultFee := 50000

		if hasInv == 0 {
			var patName, patNIK, payMethod string
			_ = tx.QueryRow(`
				SELECT p.name, p.nik, q.payment_method
				FROM queues q
				JOIN patients p ON q.patient_id = p.id
				WHERE q.id = ?
			`, queueID).Scan(&patName, &patNIK, &payMethod)

			total := regFee + consultFee + medFee
			invNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), queueID)
			_, _ = tx.Exec(`
				INSERT INTO invoices (queue_id, invoice_number, registration_fee, consultation_fee, medication_fee, total_amount, payment_method, payment_status)
				VALUES (?, ?, ?, ?, ?, ?, ?, 'PAID')
			`, queueID, invNumber, regFee, consultFee, medFee, total, payMethod)
		} else {
			_, _ = tx.Exec(`
				UPDATE invoices
				SET medication_fee = ?, total_amount = registration_fee + consultation_fee + ?
				WHERE queue_id = ?
			`, medFee, medFee, queueID)
		}
	}

	return tx.Commit()
}

func (s *Service) CancelPrescription(rxID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Release reserved stock without touching physical stock
	rows, err := tx.Query(`
		SELECT batch_id, quantity FROM prescription_items WHERE prescription_id = ?
	`, rxID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bID sql.NullInt64
		var qty int
		if err := rows.Scan(&bID, &qty); err != nil {
			return err
		}
		if bID.Valid {
			_, _ = tx.Exec(`
				UPDATE drug_batches
				SET reserved_stock = reserved_stock - ?
				WHERE id = ?
			`, qty, bID.Int64)
		}
	}
	rows.Close()

	_, err = tx.Exec("UPDATE prescriptions SET status = 'CANCELLED' WHERE id = ?", rxID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) getPrescriptionByConsultation(consID int64) (*domain.Prescription, error) {
	var rx domain.Prescription
	err := s.db.QueryRow(`
		SELECT id, consultation_id, status, created_at
		FROM prescriptions WHERE consultation_id = ?
		ORDER BY id DESC LIMIT 1
	`, consID).Scan(&rx.ID, &rx.ConsultationID, &rx.Status, &rx.CreatedAt)
	if err != nil {
		return nil, err
	}

	items, total, err := s.getPrescriptionItems(rx.ID)
	if err != nil {
		return nil, err
	}
	rx.Items = items
	rx.TotalAmount = total
	return &rx, nil
}

func (s *Service) getPrescriptionItems(rxID int64) ([]domain.PrescriptionItem, int, error) {
	rows, err := s.db.Query(`
		SELECT pi.id, pi.prescription_id, pi.drug_id, d.name, pi.batch_id, coalesce(b.batch_number, '-'),
		       pi.quantity, pi.dosage_instruction, pi.unit_price, pi.subtotal
		FROM prescription_items pi
		JOIN drugs d ON pi.drug_id = d.id
		LEFT JOIN drug_batches b ON pi.batch_id = b.id
		WHERE pi.prescription_id = ?
	`, rxID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []domain.PrescriptionItem
	total := 0
	for rows.Next() {
		var it domain.PrescriptionItem
		var bID sql.NullInt64
		if err := rows.Scan(
			&it.ID, &it.PrescriptionID, &it.DrugID, &it.DrugName, &bID, &it.BatchNumber,
			&it.Quantity, &it.DosageInstruction, &it.UnitPrice, &it.Subtotal,
		); err != nil {
			return nil, 0, err
		}
		if bID.Valid {
			it.BatchID = &bID.Int64
		}
		total += it.Subtotal
		items = append(items, it)
	}
	return items, total, nil
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}
