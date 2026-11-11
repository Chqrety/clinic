package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/consultation/domain"
)

type PostgresConsultationRepository struct {
	db *sql.DB
}

func NewPostgresConsultationRepository(db *sql.DB) *PostgresConsultationRepository {
	return &PostgresConsultationRepository{db: db}
}

func (r *PostgresConsultationRepository) ListDoctorQueues(ctx context.Context, doctorName, clinicUnit string) ([]*domain.DoctorQueuePatient, error) {
	query := `
		SELECT 
			q.id, q.queue_number, q.patient_id, p.name, p.phone,
			q.clinic_unit, q.category, q.doctor_id, q.doctor_name,
			q.registration_fee, q.payment_status, q.queue_status,
			q.called_at, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE (
			-- Only paid queues or SSC booking
			(q.payment_status = 'PAID' OR q.payment_method = 'UNPAID_BOOKING')
			AND q.queue_status NOT IN ('COMPLETED', 'CANCELLED')
		)
	`
	var args []any
	argIdx := 1

	if clinicUnit != "" && clinicUnit != "ALL" {
		query += fmt.Sprintf(" AND q.clinic_unit = $%d", argIdx)
		args = append(args, clinicUnit)
		argIdx++
	}

	if doctorName != "" {
		query += fmt.Sprintf(" AND (q.doctor_name = $%d OR q.doctor_name = '')", argIdx)
		args = append(args, doctorName)
		argIdx++
	}

	// FIFO Order
	query += " ORDER BY q.created_at ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.DoctorQueuePatient
	for rows.Next() {
		var p domain.DoctorQueuePatient
		var doctorID sql.NullInt64
		err := rows.Scan(
			&p.QueueID, &p.QueueNumber, &p.PatientID, &p.PatientName, &p.PatientPhone,
			&p.ClinicUnit, &p.Category, &doctorID, &p.DoctorName,
			&p.RegistrationFee, &p.PaymentStatus, &p.QueueStatus,
			&p.CalledAt, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if doctorID.Valid {
			val := doctorID.Int64
			p.DoctorID = &val
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *PostgresConsultationRepository) GetQueuePatientByID(ctx context.Context, queueID int64) (*domain.DoctorQueuePatient, error) {
	query := `
		SELECT 
			q.id, q.queue_number, q.patient_id, p.name, p.phone,
			q.clinic_unit, q.category, q.doctor_id, q.doctor_name,
			q.registration_fee, q.payment_status, q.queue_status,
			q.called_at, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE q.id = $1
	`
	var p domain.DoctorQueuePatient
	var doctorID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, queueID).Scan(
		&p.QueueID, &p.QueueNumber, &p.PatientID, &p.PatientName, &p.PatientPhone,
		&p.ClinicUnit, &p.Category, &doctorID, &p.DoctorName,
		&p.RegistrationFee, &p.PaymentStatus, &p.QueueStatus,
		&p.CalledAt, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrQueueNotFound
	}
	if err != nil {
		return nil, err
	}
	if doctorID.Valid {
		val := doctorID.Int64
		p.DoctorID = &val
	}
	return &p, nil
}

func (r *PostgresConsultationRepository) UpdateQueueStatus(ctx context.Context, queueID int64, status string) error {
	var query string
	if status == "CALLING" {
		query = `UPDATE queues SET queue_status = $1, called_at = CURRENT_TIMESTAMP WHERE id = $2`
	} else {
		query = `UPDATE queues SET queue_status = $1 WHERE id = $2`
	}
	_, err := r.db.ExecContext(ctx, query, status, queueID)
	return err
}

func (r *PostgresConsultationRepository) GetConsultationByQueueID(ctx context.Context, queueID int64) (*domain.Consultation, error) {
	query := `
		SELECT 
			id, queue_id, doctor_id, doctor_name, anamnesis, diagnosis,
			icd10_code, COALESCE(action_notes, ''), action_fee, status,
			created_at, updated_at
		FROM consultations
		WHERE queue_id = $1
		ORDER BY id DESC LIMIT 1
	`
	var c domain.Consultation
	var docID sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, queueID).Scan(
		&c.ID, &c.QueueID, &docID, &c.DoctorName, &c.Anamnesis, &c.Diagnosis,
		&c.ICD10Code, &c.ActionNotes, &c.ActionFee, &c.Status,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrConsultationNotFound
	}
	if err != nil {
		return nil, err
	}
	if docID.Valid {
		v := docID.Int64
		c.DoctorID = &v
	}

	// Fetch prescription if any
	rxQuery := `
		SELECT id, consultation_id, queue_id, status, COALESCE(notes, ''), created_at, updated_at
		FROM prescriptions
		WHERE consultation_id = $1
		ORDER BY id DESC LIMIT 1
	`
	var rx domain.Prescription
	err = r.db.QueryRowContext(ctx, rxQuery, c.ID).Scan(
		&rx.ID, &rx.ConsultationID, &rx.QueueID, &rx.Status, &rx.Notes, &rx.CreatedAt, &rx.UpdatedAt,
	)
	if err == nil {
		// Fetch items
		itemQuery := `
			SELECT id, prescription_id, drug_id, item_name, is_manual, quantity, unit_price, subtotal, dosage_instruction, created_at
			FROM prescription_items
			WHERE prescription_id = $1
			ORDER BY id ASC
		`
		iRows, iErr := r.db.QueryContext(ctx, itemQuery, rx.ID)
		if iErr == nil {
			defer iRows.Close()
			for iRows.Next() {
				var item domain.PrescriptionItem
				var drugID sql.NullInt64
				if err := iRows.Scan(
					&item.ID, &item.PrescriptionID, &drugID, &item.ItemName,
					&item.IsManual, &item.Quantity, &item.UnitPrice, &item.Subtotal,
					&item.DosageInstruction, &item.CreatedAt,
				); err == nil {
					if drugID.Valid {
						dVal := drugID.Int64
						item.DrugID = &dVal
					}
					rx.Items = append(rx.Items, item)
				}
			}
		}
		c.Prescription = &rx
	}

	return &c, nil
}

func (r *PostgresConsultationRepository) SaveConsultationTx(ctx context.Context, c *domain.Consultation, items []domain.PrescriptionItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert consultation
	queryCons := `
		INSERT INTO consultations (
			queue_id, doctor_id, doctor_name, anamnesis, diagnosis,
			icd10_code, action_notes, action_fee, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		) RETURNING id, created_at, updated_at
	`
	err = tx.QueryRowContext(
		ctx, queryCons,
		c.QueueID, c.DoctorID, c.DoctorName, c.Anamnesis, c.Diagnosis,
		c.ICD10Code, c.ActionNotes, c.ActionFee, c.Status,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal simpan konsultasi dokter: %w", err)
	}

	// 2. Insert prescription if items exist
	if len(items) > 0 {
		queryRx := `
			INSERT INTO prescriptions (
				consultation_id, queue_id, status, notes, created_at, updated_at
			) VALUES (
				$1, $2, 'RESERVED', $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			) RETURNING id, created_at, updated_at
		`
		var rxID int64
		var rxCreatedAt, rxUpdatedAt time.Time
		notes := "E-Resep dari " + c.DoctorName
		err = tx.QueryRowContext(ctx, queryRx, c.ID, c.QueueID, notes).Scan(&rxID, &rxCreatedAt, &rxUpdatedAt)
		if err != nil {
			return fmt.Errorf("gagal simpan e-resep: %w", err)
		}

		c.Prescription = &domain.Prescription{
			ID:             rxID,
			ConsultationID: c.ID,
			QueueID:        c.QueueID,
			Status:         "RESERVED",
			Notes:          notes,
			CreatedAt:      rxCreatedAt,
			UpdatedAt:      rxUpdatedAt,
		}

		// 3. Insert items and reserve stock
		queryItem := `
			INSERT INTO prescription_items (
				prescription_id, drug_id, item_name, is_manual,
				quantity, unit_price, subtotal, dosage_instruction, created_at
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8, CURRENT_TIMESTAMP
			) RETURNING id, created_at
		`
		for i := range items {
			items[i].PrescriptionID = rxID
			items[i].Subtotal = items[i].Quantity * items[i].UnitPrice

			err = tx.QueryRowContext(
				ctx, queryItem,
				rxID, items[i].DrugID, items[i].ItemName, items[i].IsManual,
				items[i].Quantity, items[i].UnitPrice, items[i].Subtotal, items[i].DosageInstruction,
			).Scan(&items[i].ID, &items[i].CreatedAt)
			if err != nil {
				return fmt.Errorf("gagal simpan item resep %s: %w", items[i].ItemName, err)
			}

			// Invariant: Reserve physical stock for non-manual catalog drugs
			if !items[i].IsManual && items[i].DrugID != nil {
				// Atomically increment reserved_stock on the latest active batch
				resQuery := `
					UPDATE drug_batches
					SET reserved_stock = reserved_stock + $1
					WHERE id = (
						SELECT id FROM drug_batches
						WHERE drug_id = $2 AND (physical_stock - reserved_stock) >= $1
						ORDER BY expired_date ASC LIMIT 1
					)
				`
				res, err := tx.ExecContext(ctx, resQuery, items[i].Quantity, *items[i].DrugID)
				if err != nil {
					return fmt.Errorf("gagal update reserved stock obat: %w", err)
				}
				rowsAff, _ := res.RowsAffected()
				if rowsAff == 0 {
					return domain.ErrInsufficientStock
				}
			}
		}
		c.Prescription.Items = items
	}

	// 4. Update queue status to PHARMACY_WAITING
	_, err = tx.ExecContext(ctx, `UPDATE queues SET queue_status = 'PHARMACY_WAITING' WHERE id = $1`, c.QueueID)
	if err != nil {
		return fmt.Errorf("gagal perbarui status antrian ke apotek: %w", err)
	}

	return tx.Commit()
}

func (r *PostgresConsultationRepository) ListActiveDrugs(ctx context.Context, clinicUnit string) ([]*domain.CatalogDrug, error) {
	query := `
		SELECT 
			d.id, d.clinic_unit, d.name, COALESCE(d.kfa_code, ''),
			d.unit, d.category, d.sell_price,
			COALESCE(SUM(b.physical_stock), 0) AS physical_stock,
			COALESCE(SUM(b.reserved_stock), 0) AS reserved_stock,
			d.is_active
		FROM drugs d
		LEFT JOIN drug_batches b ON d.id = b.drug_id
		WHERE d.is_active = TRUE
	`
	var args []any
	if clinicUnit != "" && clinicUnit != "ALL" {
		query += " AND (d.clinic_unit = $1 OR d.clinic_unit = 'ALL')"
		args = append(args, clinicUnit)
	}
	query += " GROUP BY d.id ORDER BY d.name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.CatalogDrug
	for rows.Next() {
		var d domain.CatalogDrug
		err := rows.Scan(
			&d.ID, &d.ClinicUnit, &d.Name, &d.KFACode,
			&d.Unit, &d.Category, &d.SellPrice,
			&d.PhysicalStock, &d.ReservedStock, &d.IsActive,
		)
		if err != nil {
			return nil, err
		}
		d.AvailableStock = d.PhysicalStock - d.ReservedStock
		if d.AvailableStock < 0 {
			d.AvailableStock = 0
		}
		list = append(list, &d)
	}
	return list, nil
}

func (r *PostgresConsultationRepository) GetDrugByID(ctx context.Context, id int64) (*domain.CatalogDrug, error) {
	query := `
		SELECT 
			d.id, d.clinic_unit, d.name, COALESCE(d.kfa_code, ''),
			d.unit, d.category, d.sell_price,
			COALESCE(SUM(b.physical_stock), 0) AS physical_stock,
			COALESCE(SUM(b.reserved_stock), 0) AS reserved_stock,
			d.is_active
		FROM drugs d
		LEFT JOIN drug_batches b ON d.id = b.drug_id
		WHERE d.id = $1
		GROUP BY d.id
	`
	var d domain.CatalogDrug
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.ClinicUnit, &d.Name, &d.KFACode,
		&d.Unit, &d.Category, &d.SellPrice,
		&d.PhysicalStock, &d.ReservedStock, &d.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrDrugNotFound
	}
	if err != nil {
		return nil, err
	}
	d.AvailableStock = d.PhysicalStock - d.ReservedStock
	if d.AvailableStock < 0 {
		d.AvailableStock = 0
	}
	return &d, nil
}
