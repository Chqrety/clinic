package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/pharmacy/domain"
)

type PostgresPharmacyRepository struct {
	db *sql.DB
}

func NewPostgresPharmacyRepository(db *sql.DB) *PostgresPharmacyRepository {
	return &PostgresPharmacyRepository{db: db}
}

func (r *PostgresPharmacyRepository) ListIncomingPrescriptions(ctx context.Context, clinicUnit string) ([]*domain.PharmacyPrescription, error) {
	query := `
		SELECT 
			pr.id, pr.consultation_id, pr.queue_id, q.queue_number,
			p.name, q.clinic_unit, c.doctor_name, pr.status,
			COALESCE(pr.notes, ''), pr.created_at, pr.updated_at
		FROM prescriptions pr
		JOIN queues q ON pr.queue_id = q.id
		JOIN patients p ON q.patient_id = p.id
		JOIN consultations c ON pr.consultation_id = c.id
		WHERE pr.status IN ('RESERVED', 'PREPARING', 'DISPENSED')
	`
	var args []any
	if clinicUnit != "" && clinicUnit != "ALL" {
		query += " AND q.clinic_unit = $1"
		args = append(args, clinicUnit)
	}
	query += " ORDER BY CASE pr.status WHEN 'RESERVED' THEN 1 WHEN 'PREPARING' THEN 2 ELSE 3 END, pr.created_at DESC LIMIT 50"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.PharmacyPrescription
	for rows.Next() {
		var rx domain.PharmacyPrescription
		err := rows.Scan(
			&rx.ID, &rx.ConsultationID, &rx.QueueID, &rx.QueueNumber,
			&rx.PatientName, &rx.ClinicUnit, &rx.DoctorName, &rx.Status,
			&rx.Notes, &rx.CreatedAt, &rx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Fetch items for each prescription
		itemQuery := `
			SELECT id, prescription_id, drug_id, item_name, is_manual, quantity, unit_price, subtotal, dosage_instruction
			FROM prescription_items
			WHERE prescription_id = $1
			ORDER BY id ASC
		`
		iRows, iErr := r.db.QueryContext(ctx, itemQuery, rx.ID)
		if iErr == nil {
			for iRows.Next() {
				var it domain.PrescriptionItemDetail
				var drugID sql.NullInt64
				if err := iRows.Scan(
					&it.ID, &it.PrescriptionID, &drugID, &it.ItemName,
					&it.IsManual, &it.Quantity, &it.UnitPrice, &it.Subtotal, &it.DosageInstruction,
				); err == nil {
					if drugID.Valid {
						v := drugID.Int64
						it.DrugID = &v
					}
					rx.Items = append(rx.Items, it)
				}
			}
			iRows.Close()
		}
		list = append(list, &rx)
	}
	return list, nil
}

func (r *PostgresPharmacyRepository) GetPrescriptionByID(ctx context.Context, id int64) (*domain.PharmacyPrescription, error) {
	query := `
		SELECT 
			pr.id, pr.consultation_id, pr.queue_id, q.queue_number,
			p.name, q.clinic_unit, c.doctor_name, pr.status,
			COALESCE(pr.notes, ''), pr.created_at, pr.updated_at
		FROM prescriptions pr
		JOIN queues q ON pr.queue_id = q.id
		JOIN patients p ON q.patient_id = p.id
		JOIN consultations c ON pr.consultation_id = c.id
		WHERE pr.id = $1
	`
	var rx domain.PharmacyPrescription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rx.ID, &rx.ConsultationID, &rx.QueueID, &rx.QueueNumber,
		&rx.PatientName, &rx.ClinicUnit, &rx.DoctorName, &rx.Status,
		&rx.Notes, &rx.CreatedAt, &rx.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPrescriptionNotFound
	}
	if err != nil {
		return nil, err
	}

	// Fetch items
	itemQuery := `
		SELECT id, prescription_id, drug_id, item_name, is_manual, quantity, unit_price, subtotal, dosage_instruction
		FROM prescription_items
		WHERE prescription_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, itemQuery, rx.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it domain.PrescriptionItemDetail
			var drugID sql.NullInt64
			if err := rows.Scan(
				&it.ID, &it.PrescriptionID, &drugID, &it.ItemName,
				&it.IsManual, &it.Quantity, &it.UnitPrice, &it.Subtotal, &it.DosageInstruction,
			); err == nil {
				if drugID.Valid {
					v := drugID.Int64
					it.DrugID = &v
				}
				rx.Items = append(rx.Items, it)
			}
		}
	}

	return &rx, nil
}

func (r *PostgresPharmacyRepository) UpdatePrescriptionStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE prescriptions SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// DispensePrescriptionTx: permanently deducts physical_stock and releases reserved_stock
func (r *PostgresPharmacyRepository) DispensePrescriptionTx(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStatus string
	var queueID int64
	err = tx.QueryRowContext(ctx, `SELECT status, queue_id FROM prescriptions WHERE id = $1 FOR UPDATE`, id).Scan(&currentStatus, &queueID)
	if err == sql.ErrNoRows {
		return domain.ErrPrescriptionNotFound
	}
	if err != nil {
		return err
	}
	if currentStatus == "DISPENSED" {
		return domain.ErrAlreadyDispensed
	}
	if currentStatus == "CANCELLED" {
		return domain.ErrAlreadyCancelled
	}

	// Fetch items
	rows, err := tx.QueryContext(ctx, `SELECT drug_id, is_manual, quantity FROM prescription_items WHERE prescription_id = $1`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	type rxDeduct struct {
		drugID   int64
		quantity int
	}
	var deducts []rxDeduct

	for rows.Next() {
		var drugID sql.NullInt64
		var isManual bool
		var qty int
		if err := rows.Scan(&drugID, &isManual, &qty); err == nil {
			if !isManual && drugID.Valid {
				deducts = append(deducts, rxDeduct{drugID: drugID.Int64, quantity: qty})
			}
		}
	}
	rows.Close()

	// Invariant deduction on drug_batches
	for _, d := range deducts {
		// Cut physical and release reserved
		updateQuery := `
			UPDATE drug_batches
			SET physical_stock = physical_stock - $1,
			    reserved_stock = reserved_stock - $1
			WHERE id = (
				SELECT id FROM drug_batches
				WHERE drug_id = $2 AND physical_stock >= $1
				ORDER BY expired_date ASC LIMIT 1
			)
		`
		res, err := tx.ExecContext(ctx, updateQuery, d.quantity, d.drugID)
		if err != nil {
			return fmt.Errorf("gagal potong stok fisik obat ID %d: %w", d.drugID, err)
		}
		rowsAff, _ := res.RowsAffected()
		if rowsAff == 0 {
			// Fallback: update any active batch for this drug
			fallbackQuery := `
				UPDATE drug_batches
				SET physical_stock = GREATEST(physical_stock - $1, 0),
				    reserved_stock = GREATEST(reserved_stock - $1, 0)
				WHERE drug_id = $2
			`
			_, _ = tx.ExecContext(ctx, fallbackQuery, d.quantity, d.drugID)
		}
	}

	// Update prescription status
	_, err = tx.ExecContext(ctx, `UPDATE prescriptions SET status = 'DISPENSED', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// Update queue status: patient ready for cashier settlement or completed
	_, err = tx.ExecContext(ctx, `
		UPDATE queues 
		SET queue_status = CASE 
			WHEN payment_status = 'PAID' THEN 'COMPLETED'
			ELSE 'PHARMACY_WAITING' 
		END,
		completed_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, queueID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CancelPrescriptionTx: releases reserved_stock without cutting physical_stock
func (r *PostgresPharmacyRepository) CancelPrescriptionTx(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStatus string
	err = tx.QueryRowContext(ctx, `SELECT status FROM prescriptions WHERE id = $1 FOR UPDATE`, id).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return domain.ErrPrescriptionNotFound
	}
	if err != nil {
		return err
	}
	if currentStatus == "DISPENSED" {
		return domain.ErrAlreadyDispensed
	}
	if currentStatus == "CANCELLED" {
		return domain.ErrAlreadyCancelled
	}

	// Fetch items
	rows, err := tx.QueryContext(ctx, `SELECT drug_id, is_manual, quantity FROM prescription_items WHERE prescription_id = $1`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var drugID sql.NullInt64
		var isManual bool
		var qty int
		if err := rows.Scan(&drugID, &isManual, &qty); err == nil {
			if !isManual && drugID.Valid {
				// Release reserved stock only
				releaseQuery := `
					UPDATE drug_batches
					SET reserved_stock = GREATEST(reserved_stock - $1, 0)
					WHERE drug_id = $2
				`
				_, _ = tx.ExecContext(ctx, releaseQuery, qty, drugID.Int64)
			}
		}
	}
	rows.Close()

	// Update prescription status
	_, err = tx.ExecContext(ctx, `UPDATE prescriptions SET status = 'CANCELLED', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresPharmacyRepository) ListDrugs(ctx context.Context, clinicUnit string) ([]*domain.Drug, error) {
	query := `
		SELECT 
			d.id, d.clinic_unit, d.name, COALESCE(d.kfa_code, ''),
			d.unit, d.category, d.min_stock_alert, d.buy_price, d.sell_price,
			COALESCE(SUM(b.physical_stock), 0) AS physical_stock,
			COALESCE(SUM(b.reserved_stock), 0) AS reserved_stock,
			d.is_active, d.created_at, d.updated_at
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

	now := time.Now()
	ninetyDaysLater := now.AddDate(0, 3, 0)

	var drugs []*domain.Drug
	for rows.Next() {
		var d domain.Drug
		err := rows.Scan(
			&d.ID, &d.ClinicUnit, &d.Name, &d.KFACode,
			&d.Unit, &d.Category, &d.MinStockAlert, &d.BuyPrice, &d.SellPrice,
			&d.PhysicalStock, &d.ReservedStock, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		d.AvailableStock = d.PhysicalStock - d.ReservedStock
		if d.AvailableStock < 0 {
			d.AvailableStock = 0
		}
		d.IsLowStock = d.AvailableStock <= d.MinStockAlert

		// Fetch batches for this drug (FEFO tracking)
		batchQuery := `
			SELECT id, drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at
			FROM drug_batches
			WHERE drug_id = $1
			ORDER BY expired_date ASC
		`
		bRows, bErr := r.db.QueryContext(ctx, batchQuery, d.ID)
		if bErr == nil {
			for bRows.Next() {
				var b domain.DrugBatch
				var expDate time.Time
				if err := bRows.Scan(
					&b.ID, &b.DrugID, &b.BatchNumber, &expDate,
					&b.PhysicalStock, &b.ReservedStock, &b.CreatedAt,
				); err == nil {
					b.ExpiredDate = expDate.Format("2006-01-02")
					b.AvailableStock = b.PhysicalStock - b.ReservedStock
					if b.AvailableStock < 0 {
						b.AvailableStock = 0
					}
					b.IsExpiringSoon = expDate.Before(ninetyDaysLater)
					if b.IsExpiringSoon {
						d.HasExpiringSoon = true
					}
					d.Batches = append(d.Batches, b)
				}
			}
			bRows.Close()
		}

		drugs = append(drugs, &d)
	}
	return drugs, nil
}

func (r *PostgresPharmacyRepository) GetDrugByID(ctx context.Context, id int64) (*domain.Drug, error) {
	query := `
		SELECT 
			d.id, d.clinic_unit, d.name, COALESCE(d.kfa_code, ''),
			d.unit, d.category, d.min_stock_alert, d.buy_price, d.sell_price,
			COALESCE(SUM(b.physical_stock), 0) AS physical_stock,
			COALESCE(SUM(b.reserved_stock), 0) AS reserved_stock,
			d.is_active, d.created_at, d.updated_at
		FROM drugs d
		LEFT JOIN drug_batches b ON d.id = b.drug_id
		WHERE d.id = $1
		GROUP BY d.id
	`
	var d domain.Drug
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.ClinicUnit, &d.Name, &d.KFACode,
		&d.Unit, &d.Category, &d.MinStockAlert, &d.BuyPrice, &d.SellPrice,
		&d.PhysicalStock, &d.ReservedStock, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
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
	d.IsLowStock = d.AvailableStock <= d.MinStockAlert
	return &d, nil
}

func (r *PostgresPharmacyRepository) CreateDrug(ctx context.Context, in *domain.CreateDrugInput) (*domain.Drug, error) {
	if in.MinStockAlert <= 0 {
		in.MinStockAlert = 10
	}
	if in.ClinicUnit == "" {
		in.ClinicUnit = "ALL"
	}
	query := `
		INSERT INTO drugs (
			clinic_unit, name, kfa_code, unit, category,
			min_stock_alert, buy_price, sell_price, is_active,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, TRUE,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		) RETURNING id, created_at, updated_at
	`
	var d domain.Drug
	d.ClinicUnit = in.ClinicUnit
	d.Name = in.Name
	d.KFACode = in.KFACode
	d.Unit = in.Unit
	d.Category = in.Category
	d.MinStockAlert = in.MinStockAlert
	d.BuyPrice = in.BuyPrice
	d.SellPrice = in.SellPrice
	d.IsActive = true

	err := r.db.QueryRowContext(
		ctx, query,
		in.ClinicUnit, in.Name, in.KFACode, in.Unit, in.Category,
		in.MinStockAlert, in.BuyPrice, in.SellPrice,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *PostgresPharmacyRepository) RestockBatch(ctx context.Context, in *domain.RestockInput) error {
	expTime, err := time.Parse("2006-01-02", in.ExpiredDate)
	if err != nil {
		return domain.ErrInvalidExpiredDate
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO drug_batches (
			drug_id, batch_number, expired_date, physical_stock, reserved_stock, created_at
		) VALUES ($1, $2, $3, $4, 0, CURRENT_TIMESTAMP)
	`
	_, err = tx.ExecContext(ctx, query, in.DrugID, in.BatchNumber, expTime, in.Quantity)
	if err != nil {
		return fmt.Errorf("gagal insert batch restock: %w", err)
	}

	if in.BuyPrice > 0 {
		_, _ = tx.ExecContext(ctx, `UPDATE drugs SET buy_price = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, in.BuyPrice, in.DrugID)
	}

	return tx.Commit()
}

func (r *PostgresPharmacyRepository) ApplyStockOpnameTx(ctx context.Context, items []domain.StockOpnameItem, adjustedBy string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateBatchQuery := `UPDATE drug_batches SET physical_stock = $1 WHERE id = $2`
	insertLogQuery := `
		INSERT INTO stock_opname_logs (
			drug_id, batch_id, system_stock, actual_stock, discrepancy, adjusted_by, notes, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
	`

	for _, it := range items {
		// Update physical stock on the batch
		if it.BatchID > 0 {
			_, err = tx.ExecContext(ctx, updateBatchQuery, it.ActualStock, it.BatchID)
			if err != nil {
				return fmt.Errorf("gagal update stok fisik batch %s: %w", it.BatchNumber, err)
			}
		}

		// Insert audit trail
		var batchID *int64
		if it.BatchID > 0 {
			batchID = &it.BatchID
		}
		notes := it.Notes
		if notes == "" {
			notes = "Penyesuaian stok opnam fisik apotek"
		}

		_, err = tx.ExecContext(
			ctx, insertLogQuery,
			it.DrugID, batchID, it.SystemStock, it.ActualStock, it.Discrepancy, adjustedBy, notes,
		)
		if err != nil {
			return fmt.Errorf("gagal catat log audit opnam obat ID %d: %w", it.DrugID, err)
		}
	}

	return tx.Commit()
}

func (r *PostgresPharmacyRepository) ListOpnameLogs(ctx context.Context, limit int) ([]*domain.StockOpnameLog, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT 
			l.id, l.drug_id, d.name, l.batch_id, COALESCE(b.batch_number, ''),
			l.system_stock, l.actual_stock, l.discrepancy, l.adjusted_by,
			COALESCE(l.notes, ''), l.created_at
		FROM stock_opname_logs l
		JOIN drugs d ON l.drug_id = d.id
		LEFT JOIN drug_batches b ON l.batch_id = b.id
		ORDER BY l.created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.StockOpnameLog
	for rows.Next() {
		var l domain.StockOpnameLog
		var batchID sql.NullInt64
		err := rows.Scan(
			&l.ID, &l.DrugID, &l.DrugName, &batchID, &l.BatchNumber,
			&l.SystemStock, &l.ActualStock, &l.Discrepancy, &l.AdjustedBy,
			&l.Notes, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if batchID.Valid {
			v := batchID.Int64
			l.BatchID = &v
		}
		logs = append(logs, &l)
	}
	return logs, nil
}
