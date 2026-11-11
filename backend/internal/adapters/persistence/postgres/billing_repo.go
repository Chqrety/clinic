package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain/billing"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresBillingRepo struct {
	db *sql.DB
}

func NewPostgresBillingRepo(db *sql.DB) *PostgresBillingRepo {
	return &PostgresBillingRepo{db: db}
}

func (r *PostgresBillingRepo) Create(ctx context.Context, inv *billing.Invoice) error {
	exec := GetExecutor(ctx, r.db)

	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}

	headerQuery := `
		INSERT INTO invoices (
			id, branch_id, invoice_number, patient_id, patient_name, patient_nik,
			queue_id, invoice_type, registration_fee, consultation_fee, medication_fee,
			total_amount, paid_amount, payment_status, payment_method, journal_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`

	_, err := exec.ExecContext(ctx, headerQuery,
		inv.ID, inv.BranchID, inv.InvoiceNumber, inv.PatientID, inv.PatientName, inv.PatientNIK,
		inv.QueueID, inv.InvoiceType, inv.RegistrationFee, inv.ConsultationFee, inv.MedicationFee,
		inv.TotalAmount, inv.PaidAmount, inv.PaymentStatus, inv.PaymentMethod, inv.JournalID,
	)
	if err != nil {
		return fmt.Errorf("insert invoice header: %w", err)
	}

	itemQuery := `
		INSERT INTO invoice_items (
			id, invoice_id, item_type, item_id, description, quantity, unit_price, total_price, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP
		)
	`

	for _, item := range inv.Items {
		itemID := item.ID
		if itemID == uuid.Nil {
			itemID = uuid.New()
		}
		_, err := exec.ExecContext(ctx, itemQuery,
			itemID, inv.ID, item.ItemType, item.ItemID, item.Description,
			item.Quantity, item.UnitPrice, item.TotalPrice,
		)
		if err != nil {
			return fmt.Errorf("insert invoice item: %w", err)
		}
	}

	return nil
}

func (r *PostgresBillingRepo) GetByID(ctx context.Context, id uuid.UUID) (*billing.Invoice, error) {
	query := `
		SELECT id, branch_id, invoice_number, patient_id, patient_name, patient_nik,
		       queue_id, invoice_type, registration_fee, consultation_fee, medication_fee,
		       total_amount, paid_amount, payment_status, payment_method, journal_id,
		       created_at, updated_at
		FROM invoices
		WHERE id = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)

	var inv billing.Invoice
	var patientID, journalID *uuid.UUID
	var queueID *int64
	var method sql.NullString

	err := row.Scan(
		&inv.ID, &inv.BranchID, &inv.InvoiceNumber, &patientID, &inv.PatientName, &inv.PatientNIK,
		&queueID, &inv.InvoiceType, &inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee,
		&inv.TotalAmount, &inv.PaidAmount, &inv.PaymentStatus, &method, &journalID,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invoice not found: %w", err)
		}
		return nil, fmt.Errorf("query invoice by id: %w", err)
	}

	inv.PatientID = patientID
	inv.QueueID = queueID
	inv.JournalID = journalID
	if method.Valid {
		inv.PaymentMethod = method.String
	}

	items, err := r.getItemsByInvoiceID(ctx, exec, inv.ID)
	if err != nil {
		return nil, err
	}
	inv.Items = items

	payments, err := r.getPaymentsByInvoiceID(ctx, exec, inv.ID)
	if err != nil {
		return nil, err
	}
	inv.Payments = payments

	return &inv, nil
}

func (r *PostgresBillingRepo) GetByNumber(ctx context.Context, number string) (*billing.Invoice, error) {
	query := `
		SELECT id, branch_id, invoice_number, patient_id, patient_name, patient_nik,
		       queue_id, invoice_type, registration_fee, consultation_fee, medication_fee,
		       total_amount, paid_amount, payment_status, payment_method, journal_id,
		       created_at, updated_at
		FROM invoices
		WHERE invoice_number = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, number)

	var inv billing.Invoice
	var patientID, journalID *uuid.UUID
	var queueID *int64
	var method sql.NullString

	err := row.Scan(
		&inv.ID, &inv.BranchID, &inv.InvoiceNumber, &patientID, &inv.PatientName, &inv.PatientNIK,
		&queueID, &inv.InvoiceType, &inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee,
		&inv.TotalAmount, &inv.PaidAmount, &inv.PaymentStatus, &method, &journalID,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invoice with number %s not found: %w", number, err)
		}
		return nil, fmt.Errorf("query invoice by number: %w", err)
	}

	inv.PatientID = patientID
	inv.QueueID = queueID
	inv.JournalID = journalID
	if method.Valid {
		inv.PaymentMethod = method.String
	}

	items, err := r.getItemsByInvoiceID(ctx, exec, inv.ID)
	if err != nil {
		return nil, err
	}
	inv.Items = items

	payments, err := r.getPaymentsByInvoiceID(ctx, exec, inv.ID)
	if err != nil {
		return nil, err
	}
	inv.Payments = payments

	return &inv, nil
}

func (r *PostgresBillingRepo) GetByQueueID(ctx context.Context, queueID int64) ([]billing.Invoice, error) {
	query := `
		SELECT id, branch_id, invoice_number, patient_id, patient_name, patient_nik,
		       queue_id, invoice_type, registration_fee, consultation_fee, medication_fee,
		       total_amount, paid_amount, payment_status, payment_method, journal_id,
		       created_at, updated_at
		FROM invoices
		WHERE queue_id = $1
		ORDER BY created_at ASC
	`
	exec := GetExecutor(ctx, r.db)
	rows, err := exec.QueryContext(ctx, query, queueID)
	if err != nil {
		return nil, fmt.Errorf("query invoices by queue id: %w", err)
	}
	defer rows.Close()

	var list []billing.Invoice
	for rows.Next() {
		var inv billing.Invoice
		var patientID, journalID *uuid.UUID
		var qID *int64
		var method sql.NullString

		if err := rows.Scan(
			&inv.ID, &inv.BranchID, &inv.InvoiceNumber, &patientID, &inv.PatientName, &inv.PatientNIK,
			&qID, &inv.InvoiceType, &inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee,
			&inv.TotalAmount, &inv.PaidAmount, &inv.PaymentStatus, &method, &journalID,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan invoice row: %w", err)
		}
		inv.PatientID = patientID
		inv.QueueID = qID
		inv.JournalID = journalID
		if method.Valid {
			inv.PaymentMethod = method.String
		}
		list = append(list, inv)
	}

	return list, nil
}

func (r *PostgresBillingRepo) List(ctx context.Context, filter billing.InvoiceFilter) ([]billing.Invoice, int64, error) {
	baseQuery := `FROM invoices WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.PaymentStatus != nil {
		baseQuery += fmt.Sprintf(" AND payment_status = $%d", argIdx)
		args = append(args, *filter.PaymentStatus)
		argIdx++
	}
	if filter.InvoiceType != nil {
		baseQuery += fmt.Sprintf(" AND invoice_type = $%d", argIdx)
		args = append(args, *filter.InvoiceType)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(invoice_number) LIKE $%d OR LOWER(patient_name) LIKE $%d OR LOWER(patient_nik) LIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, searchParam, searchParam, searchParam)
		argIdx += 3
	}

	exec := GetExecutor(ctx, r.db)

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count invoices: %w", err)
	}

	allowedSorts := map[string]string{
		"created_at":     "created_at",
		"invoice_number": "invoice_number",
		"patient_name":   "patient_name",
		"total_amount":   "total_amount",
		"paid_amount":    "paid_amount",
		"payment_status": "payment_status",
	}
	sortCol := "created_at"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, invoice_number, patient_id, patient_name, patient_nik,
		       queue_id, invoice_type, registration_fee, consultation_fee, medication_fee,
		       total_amount, paid_amount, payment_status, payment_method, journal_id,
		       created_at, updated_at
	` + baseQuery + orderBy

	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		selectQuery += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := exec.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []billing.Invoice
	for rows.Next() {
		var inv billing.Invoice
		var patientID, journalID *uuid.UUID
		var queueID *int64
		var method sql.NullString

		if err := rows.Scan(
			&inv.ID, &inv.BranchID, &inv.InvoiceNumber, &patientID, &inv.PatientName, &inv.PatientNIK,
			&queueID, &inv.InvoiceType, &inv.RegistrationFee, &inv.ConsultationFee, &inv.MedicationFee,
			&inv.TotalAmount, &inv.PaidAmount, &inv.PaymentStatus, &method, &journalID,
			&inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan invoice row: %w", err)
		}
		inv.PatientID = patientID
		inv.QueueID = queueID
		inv.JournalID = journalID
		if method.Valid {
			inv.PaymentMethod = method.String
		}
		invoices = append(invoices, inv)
	}

	return invoices, totalCount, nil
}

func (r *PostgresBillingRepo) UpdatePayment(ctx context.Context, id uuid.UUID, paidAmount decimal.Decimal, status billing.PaymentStatus, journalID *uuid.UUID) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE invoices
		SET paid_amount = $1, payment_status = $2, journal_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	res, err := exec.ExecContext(ctx, query, paidAmount, status, journalID, id)
	if err != nil {
		return fmt.Errorf("update invoice payment status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("invoice not found or not updated")
	}
	return nil
}

func (r *PostgresBillingRepo) AddPaymentRecord(ctx context.Context, record *billing.PaymentRecord) error {
	exec := GetExecutor(ctx, r.db)

	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}
	if record.PaidAt.IsZero() {
		record.PaidAt = time.Now()
	}

	query := `
		INSERT INTO payment_records (
			id, invoice_id, payment_number, amount, payment_method,
			reference_no, notes, received_by, journal_id, paid_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10
		)
	`
	_, err := exec.ExecContext(ctx, query,
		record.ID, record.InvoiceID, record.PaymentNumber, record.Amount, record.PaymentMethod,
		record.ReferenceNo, record.Notes, record.ReceivedBy, record.JournalID, record.PaidAt,
	)
	if err != nil {
		return fmt.Errorf("insert payment record: %w", err)
	}
	return nil
}

func (r *PostgresBillingRepo) getItemsByInvoiceID(ctx context.Context, exec QueryExecutor, invoiceID uuid.UUID) ([]billing.InvoiceItem, error) {
	query := `
		SELECT id, invoice_id, item_type, item_id, description, quantity, unit_price, total_price
		FROM invoice_items
		WHERE invoice_id = $1
	`
	rows, err := exec.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("query invoice items: %w", err)
	}
	defer rows.Close()

	var items []billing.InvoiceItem
	for rows.Next() {
		var item billing.InvoiceItem
		var itemID *uuid.UUID
		if err := rows.Scan(
			&item.ID, &item.InvoiceID, &item.ItemType, &itemID, &item.Description,
			&item.Quantity, &item.UnitPrice, &item.TotalPrice,
		); err != nil {
			return nil, fmt.Errorf("scan invoice item: %w", err)
		}
		item.ItemID = itemID
		items = append(items, item)
	}
	return items, nil
}

func (r *PostgresBillingRepo) getPaymentsByInvoiceID(ctx context.Context, exec QueryExecutor, invoiceID uuid.UUID) ([]billing.PaymentRecord, error) {
	query := `
		SELECT id, invoice_id, payment_number, amount, payment_method, reference_no, notes, received_by, journal_id, paid_at
		FROM payment_records
		WHERE invoice_id = $1
		ORDER BY paid_at ASC
	`
	rows, err := exec.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("query payment records: %w", err)
	}
	defer rows.Close()

	var records []billing.PaymentRecord
	for rows.Next() {
		var rec billing.PaymentRecord
		var recBy, jID *uuid.UUID
		var refNo, notes sql.NullString

		if err := rows.Scan(
			&rec.ID, &rec.InvoiceID, &rec.PaymentNumber, &rec.Amount, &rec.PaymentMethod,
			&refNo, &notes, &recBy, &jID, &rec.PaidAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment record: %w", err)
		}
		rec.ReceivedBy = recBy
		rec.JournalID = jID
		if refNo.Valid {
			rec.ReferenceNo = refNo.String
		}
		if notes.Valid {
			rec.Notes = notes.String
		}
		records = append(records, rec)
	}
	return records, nil
}
