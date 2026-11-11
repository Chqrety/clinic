package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"clinic-api/internal/cashier/domain"
)

type PostgresCashierRepository struct {
	db *sql.DB
}

func NewPostgresCashierRepository(db *sql.DB) *PostgresCashierRepository {
	return &PostgresCashierRepository{db: db}
}

// GenerateInvoiceNumber formats INV-YYYYMMDD-XXXX atomically
func (r *PostgresCashierRepository) GenerateInvoiceNumber(ctx context.Context) (string, error) {
	today := time.Now().Format("20060102")
	var count int
	query := `SELECT COUNT(*) FROM invoices WHERE invoice_number LIKE $1`
	err := r.db.QueryRowContext(ctx, query, fmt.Sprintf("INV-%s-%%", today)).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("gagal menghitung nomor invoice: %w", err)
	}
	return fmt.Sprintf("INV-%s-%04d", today, count+1), nil
}

func (r *PostgresCashierRepository) Create(ctx context.Context, inv *domain.Invoice) error {
	query := `
		INSERT INTO invoices (
			queue_id, patient_id, invoice_number,
			subtotal_registration, subtotal_action, subtotal_medication,
			voucher_discount, total_amount, payment_method,
			cash_received, cash_change, payment_reference,
			cashier_staff_id, cashier_name, status, notes,
			created_at, updated_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9,
			$10, $11, $12,
			$13, $14, $15, $16,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		inv.QueueID, inv.PatientID, inv.InvoiceNumber,
		inv.SubtotalRegistration, inv.SubtotalAction, inv.SubtotalMedication,
		inv.VoucherDiscount, inv.TotalAmount, string(inv.PaymentMethod),
		inv.CashReceived, inv.CashChange, inv.PaymentReference,
		inv.CashierStaffID, inv.CashierName, string(inv.Status), inv.Notes,
	).Scan(&inv.ID, &inv.CreatedAt, &inv.UpdatedAt)
}

func (r *PostgresCashierRepository) GetByID(ctx context.Context, id int64) (*domain.Invoice, error) {
	query := `
		SELECT 
			i.id, i.queue_id, i.patient_id, i.invoice_number,
			p.name, p.phone, q.clinic_unit, q.queue_number,
			i.subtotal_registration, i.subtotal_action, i.subtotal_medication,
			i.voucher_discount, i.total_amount, i.payment_method,
			i.cash_received, i.cash_change, COALESCE(i.payment_reference, ''),
			COALESCE(i.cashier_staff_id, 0), i.cashier_name, i.status,
			COALESCE(i.refund_bank_name, ''), COALESCE(i.refund_bank_account, ''),
			COALESCE(i.refund_account_name, ''), COALESCE(i.refund_reason, ''),
			i.refunded_at, COALESCE(i.notes, ''),
			i.created_at, i.updated_at
		FROM invoices i
		JOIN patients p ON i.patient_id = p.id
		JOIN queues q ON i.queue_id = q.id
		WHERE i.id = $1
	`
	var inv domain.Invoice
	var method, status string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.QueueID, &inv.PatientID, &inv.InvoiceNumber,
		&inv.PatientName, &inv.PatientPhone, &inv.ClinicUnit, &inv.QueueNumber,
		&inv.SubtotalRegistration, &inv.SubtotalAction, &inv.SubtotalMedication,
		&inv.VoucherDiscount, &inv.TotalAmount, &method,
		&inv.CashReceived, &inv.CashChange, &inv.PaymentReference,
		&inv.CashierStaffID, &inv.CashierName, &status,
		&inv.RefundBankName, &inv.RefundBankAccount,
		&inv.RefundAccountName, &inv.RefundReason,
		&inv.RefundedAt, &inv.Notes,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	inv.PaymentMethod = domain.PaymentMethod(method)
	inv.Status = domain.InvoiceStatus(status)
	return &inv, nil
}

func (r *PostgresCashierRepository) GetByQueueID(ctx context.Context, queueID int64) (*domain.Invoice, error) {
	query := `
		SELECT 
			i.id, i.queue_id, i.patient_id, i.invoice_number,
			p.name, p.phone, q.clinic_unit, q.queue_number,
			i.subtotal_registration, i.subtotal_action, i.subtotal_medication,
			i.voucher_discount, i.total_amount, i.payment_method,
			i.cash_received, i.cash_change, COALESCE(i.payment_reference, ''),
			COALESCE(i.cashier_staff_id, 0), i.cashier_name, i.status,
			COALESCE(i.refund_bank_name, ''), COALESCE(i.refund_bank_account, ''),
			COALESCE(i.refund_account_name, ''), COALESCE(i.refund_reason, ''),
			i.refunded_at, COALESCE(i.notes, ''),
			i.created_at, i.updated_at
		FROM invoices i
		JOIN patients p ON i.patient_id = p.id
		JOIN queues q ON i.queue_id = q.id
		WHERE i.queue_id = $1
		ORDER BY i.id DESC LIMIT 1
	`
	var inv domain.Invoice
	var method, status string
	err := r.db.QueryRowContext(ctx, query, queueID).Scan(
		&inv.ID, &inv.QueueID, &inv.PatientID, &inv.InvoiceNumber,
		&inv.PatientName, &inv.PatientPhone, &inv.ClinicUnit, &inv.QueueNumber,
		&inv.SubtotalRegistration, &inv.SubtotalAction, &inv.SubtotalMedication,
		&inv.VoucherDiscount, &inv.TotalAmount, &method,
		&inv.CashReceived, &inv.CashChange, &inv.PaymentReference,
		&inv.CashierStaffID, &inv.CashierName, &status,
		&inv.RefundBankName, &inv.RefundBankAccount,
		&inv.RefundAccountName, &inv.RefundReason,
		&inv.RefundedAt, &inv.Notes,
		&inv.CreatedAt, &inv.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInvoiceNotFound
	}
	if err != nil {
		return nil, err
	}
	inv.PaymentMethod = domain.PaymentMethod(method)
	inv.Status = domain.InvoiceStatus(status)
	return &inv, nil
}

func (r *PostgresCashierRepository) List(ctx context.Context, limit, offset int) ([]*domain.Invoice, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT 
			i.id, i.queue_id, i.patient_id, i.invoice_number,
			p.name, p.phone, q.clinic_unit, q.queue_number,
			i.subtotal_registration, i.subtotal_action, i.subtotal_medication,
			i.voucher_discount, i.total_amount, i.payment_method,
			i.cash_received, i.cash_change, COALESCE(i.payment_reference, ''),
			COALESCE(i.cashier_staff_id, 0), i.cashier_name, i.status,
			COALESCE(i.refund_bank_name, ''), COALESCE(i.refund_bank_account, ''),
			COALESCE(i.refund_account_name, ''), COALESCE(i.refund_reason, ''),
			i.refunded_at, COALESCE(i.notes, ''),
			i.created_at, i.updated_at
		FROM invoices i
		JOIN patients p ON i.patient_id = p.id
		JOIN queues q ON i.queue_id = q.id
		ORDER BY i.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Invoice
	for rows.Next() {
		var inv domain.Invoice
		var method, status string
		err := rows.Scan(
			&inv.ID, &inv.QueueID, &inv.PatientID, &inv.InvoiceNumber,
			&inv.PatientName, &inv.PatientPhone, &inv.ClinicUnit, &inv.QueueNumber,
			&inv.SubtotalRegistration, &inv.SubtotalAction, &inv.SubtotalMedication,
			&inv.VoucherDiscount, &inv.TotalAmount, &method,
			&inv.CashReceived, &inv.CashChange, &inv.PaymentReference,
			&inv.CashierStaffID, &inv.CashierName, &status,
			&inv.RefundBankName, &inv.RefundBankAccount,
			&inv.RefundAccountName, &inv.RefundReason,
			&inv.RefundedAt, &inv.Notes,
			&inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inv.PaymentMethod = domain.PaymentMethod(method)
		inv.Status = domain.InvoiceStatus(status)
		result = append(result, &inv)
	}
	return result, nil
}

func (r *PostgresCashierRepository) UpdateRefund(ctx context.Context, id int64, bankName, bankAccount, accountName, reason string, refundedAt time.Time) error {
	query := `
		UPDATE invoices
		SET status = 'REFUNDED',
		    refund_bank_name = $1,
		    refund_bank_account = $2,
		    refund_account_name = $3,
		    refund_reason = $4,
		    refunded_at = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`
	res, err := r.db.ExecContext(ctx, query, bankName, bankAccount, accountName, reason, refundedAt, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrInvoiceNotFound
	}
	return nil
}

func (r *PostgresCashierRepository) GetShiftSummary(ctx context.Context, date string) (*domain.ShiftSummary, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	query := `
		SELECT 
			COUNT(id) AS total_transactions,
			COALESCE(SUM(CASE WHEN payment_method = 'CASH' AND status = 'PAID' THEN total_amount ELSE 0 END), 0) AS total_cash,
			COALESCE(SUM(CASE WHEN payment_method = 'BANK_TRANSFER' AND status = 'PAID' THEN total_amount ELSE 0 END), 0) AS total_transfer,
			COALESCE(SUM(CASE WHEN payment_method = 'EDC' AND status = 'PAID' THEN total_amount ELSE 0 END), 0) AS total_edc,
			COALESCE(SUM(voucher_discount), 0) AS total_voucher,
			COALESCE(SUM(CASE WHEN status = 'REFUNDED' THEN total_amount ELSE 0 END), 0) AS total_refund
		FROM invoices
		WHERE DATE(created_at) = $1
	`
	var sum domain.ShiftSummary
	sum.ShiftDate = date
	err := r.db.QueryRowContext(ctx, query, date).Scan(
		&sum.TotalTransactions,
		&sum.TotalCash,
		&sum.TotalTransfer,
		&sum.TotalEDC,
		&sum.TotalVoucher,
		&sum.TotalRefund,
	)
	if err != nil {
		return nil, err
	}
	sum.NetTotal = (sum.TotalCash + sum.TotalTransfer + sum.TotalEDC) - sum.TotalRefund
	return &sum, nil
}

func (r *PostgresCashierRepository) CreateReport(ctx context.Context, report *domain.CashierShiftReport) error {
	query := `
		INSERT INTO cashier_shift_reports (
			staff_id, staff_name, shift_date,
			total_transactions, total_cash, total_transfer,
			total_edc, total_voucher, total_refund, net_total,
			notes, sent_to_owner_at, created_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9, $10,
			$11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		) RETURNING id, sent_to_owner_at, created_at
	`
	return r.db.QueryRowContext(
		ctx, query,
		report.StaffID, report.StaffName, report.ShiftDate,
		report.TotalTransactions, report.TotalCash, report.TotalTransfer,
		report.TotalEDC, report.TotalVoucher, report.TotalRefund, report.NetTotal,
		report.Notes,
	).Scan(&report.ID, &report.SentToOwnerAt, &report.CreatedAt)
}

func (r *PostgresCashierRepository) ListReports(ctx context.Context, limit int) ([]*domain.CashierShiftReport, error) {
	if limit <= 0 {
		limit = 30
	}
	query := `
		SELECT 
			id, staff_id, staff_name, shift_date,
			total_transactions, total_cash, total_transfer,
			total_edc, total_voucher, total_refund, net_total,
			COALESCE(notes, ''), sent_to_owner_at, created_at
		FROM cashier_shift_reports
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*domain.CashierShiftReport
	for rows.Next() {
		var rep domain.CashierShiftReport
		err := rows.Scan(
			&rep.ID, &rep.StaffID, &rep.StaffName, &rep.ShiftDate,
			&rep.TotalTransactions, &rep.TotalCash, &rep.TotalTransfer,
			&rep.TotalEDC, &rep.TotalVoucher, &rep.TotalRefund, &rep.NetTotal,
			&rep.Notes, &rep.SentToOwnerAt, &rep.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &rep)
	}
	return reports, nil
}

func (r *PostgresCashierRepository) GetQueueByID(ctx context.Context, id int64) (*domain.PendingBill, error) {
	query := `
		SELECT 
			q.id, q.queue_number, q.patient_id, p.name, p.phone,
			q.clinic_unit, q.category, q.doctor_name,
			q.registration_fee, q.queue_status, q.payment_status, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE q.id = $1
	`
	var bill domain.PendingBill
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bill.QueueID, &bill.QueueNumber, &bill.PatientID, &bill.PatientName, &bill.PatientPhone,
		&bill.ClinicUnit, &bill.Category, &bill.DoctorName,
		&bill.SubtotalRegistration, &bill.QueueStatus, &bill.PaymentStatus, &bill.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrQueueNotFound
	}
	if err != nil {
		return nil, err
	}
	bill.TotalDue = bill.SubtotalRegistration
	if bill.PaymentStatus == "PENDING" {
		bill.BillType = "REGISTRATION_ONLY"
	} else {
		bill.BillType = "FULL_VISIT"
	}
	return &bill, nil
}

func (r *PostgresCashierRepository) ListPendingBills(ctx context.Context) ([]*domain.PendingBill, error) {
	// Queues that either have pending registration payment or have completed exam / dispensed pharmacy needing billing
	query := `
		SELECT 
			q.id, q.queue_number, q.patient_id, p.name, p.phone,
			q.clinic_unit, q.category, q.doctor_name,
			q.registration_fee, q.queue_status, q.payment_status, q.created_at
		FROM queues q
		JOIN patients p ON q.patient_id = p.id
		WHERE (
			-- 1. Pending registration (Cash walk-in / pending cashier)
			(q.payment_status = 'PENDING' AND q.queue_status != 'CANCELLED')
			OR
			-- 2. Post-treatment / completed waiting for final settlement if no paid invoice exists
			(q.queue_status IN ('IN_EXAM', 'PHARMACY_WAITING', 'COMPLETED') AND q.id NOT IN (
				SELECT queue_id FROM invoices WHERE status = 'PAID'
			))
		)
		ORDER BY q.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.PendingBill
	for rows.Next() {
		var bill domain.PendingBill
		err := rows.Scan(
			&bill.QueueID, &bill.QueueNumber, &bill.PatientID, &bill.PatientName, &bill.PatientPhone,
			&bill.ClinicUnit, &bill.Category, &bill.DoctorName,
			&bill.SubtotalRegistration, &bill.QueueStatus, &bill.PaymentStatus, &bill.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Defaults for estimation:
		// If SSC, standard facial/treatment consultation estimate
		if bill.ClinicUnit == "SSC" {
			bill.SubtotalAction = 150000
			bill.SubtotalMedication = 85000
			bill.BillType = "FULL_VISIT"
		} else if bill.PaymentStatus == "PENDING" {
			bill.BillType = "REGISTRATION_ONLY"
		} else {
			bill.SubtotalAction = 100000
			bill.SubtotalMedication = 45000
			bill.BillType = "FULL_VISIT"
		}
		bill.TotalDue = bill.SubtotalRegistration + bill.SubtotalAction + bill.SubtotalMedication
		list = append(list, &bill)
	}
	return list, nil
}

func (r *PostgresCashierRepository) UpdatePaymentAndQueueStatus(ctx context.Context, queueID int64, paymentStatus, queueStatus string) error {
	var query string
	if paymentStatus == "PAID" {
		query = `
			UPDATE queues
			SET payment_status = $1, queue_status = $2, paid_at = CURRENT_TIMESTAMP
			WHERE id = $3
		`
	} else {
		query = `
			UPDATE queues
			SET payment_status = $1, queue_status = $2
			WHERE id = $3
		`
	}
	_, err := r.db.ExecContext(ctx, query, paymentStatus, queueStatus, queueID)
	return err
}

func (r *PostgresCashierRepository) RegisterWalkIn(ctx context.Context, req *domain.WalkInRequest, queueNumber string, fee int) (int64, string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, "", err
	}
	defer tx.Rollback()

	// 1. Patient find or create
	var patientID int64
	err = tx.QueryRowContext(ctx, `SELECT id FROM patients WHERE phone = $1 LIMIT 1`, req.PatientPhone).Scan(&patientID)
	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(ctx, `
			INSERT INTO patients (name, phone, created_at, updated_at)
			VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id
		`, req.PatientName, req.PatientPhone).Scan(&patientID)
		if err != nil {
			return 0, "", fmt.Errorf("gagal membuat data pasien walk-in: %w", err)
		}
	} else if err != nil {
		return 0, "", err
	}

	category := "E_WALKIN"
	if req.ClinicUnit == "SSC" {
		category = "SSC_TREATMENT"
	}

	paymentStatus := "PENDING"
	if req.IsPaidImmediate {
		paymentStatus = "PAID"
	}

	queueStatus := "WAITING"

	insertQueueQuery := `
		INSERT INTO queues (
			patient_id, clinic_unit, category, doctor_id, doctor_name,
			queue_number, registration_fee, payment_method, payment_status,
			queue_status, paid_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, CASE WHEN $9 = 'PAID' THEN CURRENT_TIMESTAMP ELSE NULL END,
			CURRENT_TIMESTAMP
		) RETURNING id
	`
	var queueID int64
	var doctorID *int64
	if req.DoctorID > 0 {
		doctorID = &req.DoctorID
	}
	err = tx.QueryRowContext(
		ctx, insertQueueQuery,
		patientID, req.ClinicUnit, category, doctorID, req.DoctorName,
		queueNumber, fee, req.PaymentMethod, paymentStatus, queueStatus,
	).Scan(&queueID)
	if err != nil {
		return 0, "", fmt.Errorf("gagal insert antrian walk-in: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, "", err
	}
	return queueID, queueNumber, nil
}
