package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain/payable"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresPayableRepo struct {
	db *sql.DB
}

func NewPostgresPayableRepo(db *sql.DB) *PostgresPayableRepo {
	return &PostgresPayableRepo{db: db}
}

func (r *PostgresPayableRepo) Create(ctx context.Context, p *payable.Payable) error {
	exec := GetExecutor(ctx, r.db)
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	query := `
		INSERT INTO payables (
			id, branch_id, payable_number, vendor_type, vendor_id, vendor_name,
			invoice_reference, total_amount, outstanding_amount, issue_date, due_date,
			status, account_id, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		p.ID, p.BranchID, p.PayableNumber, p.VendorType, p.VendorID, p.VendorName,
		p.InvoiceReference, p.TotalAmount, p.OutstandingAmount, p.IssueDate, p.DueDate,
		p.Status, p.AccountID, p.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert payable: %w", err)
	}
	return nil
}

func (r *PostgresPayableRepo) GetByID(ctx context.Context, id uuid.UUID) (*payable.Payable, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT p.id, p.branch_id, p.payable_number, p.vendor_type, p.vendor_id, p.vendor_name,
		       p.invoice_reference, p.total_amount, p.outstanding_amount, p.issue_date, p.due_date,
		       p.status, p.account_id, a.code as account_code, COALESCE(p.notes, ''), p.created_at, p.updated_at
		FROM payables p
		JOIN accounts a ON p.account_id = a.id
		WHERE p.id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)

	var p payable.Payable
	var vendorID *uuid.UUID
	var total, outstanding decimal.Decimal

	err := row.Scan(
		&p.ID, &p.BranchID, &p.PayableNumber, &p.VendorType, &vendorID, &p.VendorName,
		&p.InvoiceReference, &total, &outstanding, &p.IssueDate, &p.DueDate,
		&p.Status, &p.AccountID, &p.AccountCode, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payable not found: %w", err)
		}
		return nil, fmt.Errorf("query payable: %w", err)
	}

	p.VendorID = vendorID
	p.TotalAmount = total
	p.OutstandingAmount = outstanding

	payments, err := r.getPaymentsByPayableID(ctx, exec, p.ID)
	if err != nil {
		return nil, err
	}
	p.Payments = payments

	return &p, nil
}

func (r *PostgresPayableRepo) GetByNumber(ctx context.Context, number string) (*payable.Payable, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT p.id, p.branch_id, p.payable_number, p.vendor_type, p.vendor_id, p.vendor_name,
		       p.invoice_reference, p.total_amount, p.outstanding_amount, p.issue_date, p.due_date,
		       p.status, p.account_id, a.code as account_code, COALESCE(p.notes, ''), p.created_at, p.updated_at
		FROM payables p
		JOIN accounts a ON p.account_id = a.id
		WHERE p.payable_number = $1
	`
	row := exec.QueryRowContext(ctx, query, number)

	var p payable.Payable
	var vendorID *uuid.UUID
	var total, outstanding decimal.Decimal

	err := row.Scan(
		&p.ID, &p.BranchID, &p.PayableNumber, &p.VendorType, &vendorID, &p.VendorName,
		&p.InvoiceReference, &total, &outstanding, &p.IssueDate, &p.DueDate,
		&p.Status, &p.AccountID, &p.AccountCode, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payable with number %s not found: %w", number, err)
		}
		return nil, fmt.Errorf("query payable by number: %w", err)
	}

	p.VendorID = vendorID
	p.TotalAmount = total
	p.OutstandingAmount = outstanding

	payments, err := r.getPaymentsByPayableID(ctx, exec, p.ID)
	if err != nil {
		return nil, err
	}
	p.Payments = payments

	return &p, nil
}

func (r *PostgresPayableRepo) List(ctx context.Context, filter payable.Filter) ([]payable.Payable, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `
		FROM payables p
		JOIN accounts a ON p.account_id = a.id
		WHERE 1=1
	`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND p.branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.VendorType != nil {
		baseQuery += fmt.Sprintf(" AND p.vendor_type = $%d", argIdx)
		args = append(args, *filter.VendorType)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND p.status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND p.issue_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND p.issue_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(p.vendor_name) LIKE $%d OR LOWER(p.invoice_reference) LIKE $%d OR LOWER(p.payable_number) LIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, searchParam, searchParam, searchParam)
		argIdx += 3
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count payables: %w", err)
	}

	allowedSorts := map[string]string{
		"due_date":           "p.due_date",
		"issue_date":         "p.issue_date",
		"created_at":         "p.created_at",
		"payable_number":     "p.payable_number",
		"vendor_name":        "p.vendor_name",
		"total_amount":       "p.total_amount",
		"outstanding_amount": "p.outstanding_amount",
		"status":             "p.status",
	}
	sortCol := "p.due_date"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortDir = "DESC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, p.id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT p.id, p.branch_id, p.payable_number, p.vendor_type, p.vendor_id, p.vendor_name,
		       p.invoice_reference, p.total_amount, p.outstanding_amount, p.issue_date, p.due_date,
		       p.status, p.account_id, a.code as account_code, COALESCE(p.notes, ''), p.created_at, p.updated_at
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
		return nil, 0, fmt.Errorf("query payables: %w", err)
	}
	defer rows.Close()

	var list []payable.Payable
	for rows.Next() {
		var p payable.Payable
		var vendorID *uuid.UUID
		var total, outstanding decimal.Decimal

		if err := rows.Scan(
			&p.ID, &p.BranchID, &p.PayableNumber, &p.VendorType, &vendorID, &p.VendorName,
			&p.InvoiceReference, &total, &outstanding, &p.IssueDate, &p.DueDate,
			&p.Status, &p.AccountID, &p.AccountCode, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan payable row: %w", err)
		}
		p.VendorID = vendorID
		p.TotalAmount = total
		p.OutstandingAmount = outstanding
		list = append(list, p)
	}

	return list, totalCount, nil
}

func (r *PostgresPayableRepo) UpdateOutstanding(ctx context.Context, id uuid.UUID, outstanding decimal.Decimal, status payable.Status) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payables
		SET outstanding_amount = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	res, err := exec.ExecContext(ctx, query, outstanding, status, id)
	if err != nil {
		return fmt.Errorf("update payable status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("payable not found or no changes made")
	}
	return nil
}

func (r *PostgresPayableRepo) AddPayment(ctx context.Context, p *payable.Payment) error {
	exec := GetExecutor(ctx, r.db)
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	query := `
		INSERT INTO payable_payments (
			id, payable_id, payment_number, amount, payment_date,
			cash_bank_account_code, reference_no, notes, journal_id, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		p.ID, p.PayableID, p.PaymentNumber, p.Amount, p.PaymentDate,
		p.CashBankAccountCode, p.ReferenceNo, p.Notes, p.JournalID, p.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert payable payment: %w", err)
	}
	return nil
}

func (r *PostgresPayableRepo) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*payable.AgingReport, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT 
			vendor_name,
			COALESCE(SUM(CASE WHEN due_date >= $1 THEN outstanding_amount ELSE 0 END), 0) as current_amount,
			COALESCE(SUM(CASE WHEN $1 - due_date BETWEEN 1 AND 30 THEN outstanding_amount ELSE 0 END), 0) as days_1_30,
			COALESCE(SUM(CASE WHEN $1 - due_date BETWEEN 31 AND 60 THEN outstanding_amount ELSE 0 END), 0) as days_31_60,
			COALESCE(SUM(CASE WHEN $1 - due_date BETWEEN 61 AND 90 THEN outstanding_amount ELSE 0 END), 0) as days_61_90,
			COALESCE(SUM(CASE WHEN $1 - due_date > 90 THEN outstanding_amount ELSE 0 END), 0) as days_over_90,
			COALESCE(SUM(outstanding_amount), 0) as total_outstanding
		FROM payables
		WHERE status IN ('OPEN', 'PARTIAL', 'OVERDUE')
		  AND ($2::uuid IS NULL OR branch_id = $2)
		GROUP BY vendor_name
		ORDER BY total_outstanding DESC
	`
	rows, err := exec.QueryContext(ctx, query, asOfDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query payable aging: %w", err)
	}
	defer rows.Close()

	var buckets []payable.AgingBucket
	totalCurrent := decimal.Zero
	total1To30 := decimal.Zero
	total31To60 := decimal.Zero
	total61To90 := decimal.Zero
	totalOver90 := decimal.Zero
	grandTotal := decimal.Zero

	for rows.Next() {
		var b payable.AgingBucket
		if err := rows.Scan(
			&b.VendorName, &b.Current, &b.Days1To30, &b.Days31To60, &b.Days61To90, &b.DaysOver90, &b.TotalOutstanding,
		); err != nil {
			return nil, fmt.Errorf("scan payable aging row: %w", err)
		}
		totalCurrent = totalCurrent.Add(b.Current)
		total1To30 = total1To30.Add(b.Days1To30)
		total31To60 = total31To60.Add(b.Days31To60)
		total61To90 = total61To90.Add(b.Days61To90)
		totalOver90 = totalOver90.Add(b.DaysOver90)
		grandTotal = grandTotal.Add(b.TotalOutstanding)

		buckets = append(buckets, b)
	}

	return &payable.AgingReport{
		BranchID:         branchID,
		AsOfDate:         asOfDate,
		TotalCurrent:     totalCurrent,
		Total1To30:       total1To30,
		Total31To60:      total31To60,
		Total61To90:      total61To90,
		TotalOver90:      totalOver90,
		GrandTotal:       grandTotal,
		Buckets:          buckets,
	}, nil
}

func (r *PostgresPayableRepo) getPaymentsByPayableID(ctx context.Context, exec QueryExecutor, payableID uuid.UUID) ([]payable.Payment, error) {
	query := `
		SELECT id, payable_id, payment_number, amount, payment_date,
		       cash_bank_account_code, COALESCE(reference_no, ''), COALESCE(notes, ''),
		       journal_id, created_by, created_at
		FROM payable_payments
		WHERE payable_id = $1
		ORDER BY payment_date ASC, created_at ASC
	`
	rows, err := exec.QueryContext(ctx, query, payableID)
	if err != nil {
		return nil, fmt.Errorf("query payable payments: %w", err)
	}
	defer rows.Close()

	var payments []payable.Payment
	for rows.Next() {
		var p payable.Payment
		var jID, createdBy *uuid.UUID
		var amt decimal.Decimal

		if err := rows.Scan(
			&p.ID, &p.PayableID, &p.PaymentNumber, &amt, &p.PaymentDate,
			&p.CashBankAccountCode, &p.ReferenceNo, &p.Notes,
			&jID, &createdBy, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payable payment row: %w", err)
		}
		p.Amount = amt
		p.JournalID = jID
		p.CreatedBy = createdBy
		payments = append(payments, p)
	}
	return payments, nil
}

