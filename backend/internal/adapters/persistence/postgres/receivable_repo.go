package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/domain/receivable"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresReceivableRepo struct {
	db *sql.DB
}

func NewPostgresReceivableRepo(db *sql.DB) *PostgresReceivableRepo {
	return &PostgresReceivableRepo{db: db}
}

func (r *PostgresReceivableRepo) Create(ctx context.Context, rec *receivable.Receivable) error {
	exec := GetExecutor(ctx, r.db)
	if rec.ID == uuid.Nil {
		rec.ID = uuid.New()
	}

	query := `
		INSERT INTO receivables (
			id, branch_id, receivable_number, customer_type, customer_id, customer_name,
			invoice_reference, total_amount, outstanding_amount, issue_date, due_date,
			status, account_id, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		rec.ID, rec.BranchID, rec.ReceivableNumber, rec.CustomerType, rec.CustomerID, rec.CustomerName,
		rec.InvoiceReference, rec.TotalAmount, rec.OutstandingAmount, rec.IssueDate, rec.DueDate,
		rec.Status, rec.AccountID, rec.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert receivable: %w", err)
	}
	return nil
}

func (r *PostgresReceivableRepo) GetByID(ctx context.Context, id uuid.UUID) (*receivable.Receivable, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT r.id, r.branch_id, r.receivable_number, r.customer_type, r.customer_id, r.customer_name,
		       r.invoice_reference, r.total_amount, r.outstanding_amount, r.issue_date, r.due_date,
		       r.status, r.account_id, a.code as account_code, COALESCE(r.notes, ''), r.created_at, r.updated_at
		FROM receivables r
		JOIN accounts a ON r.account_id = a.id
		WHERE r.id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)

	var rec receivable.Receivable
	var customerID *uuid.UUID
	var total, outstanding decimal.Decimal

	err := row.Scan(
		&rec.ID, &rec.BranchID, &rec.ReceivableNumber, &rec.CustomerType, &customerID, &rec.CustomerName,
		&rec.InvoiceReference, &total, &outstanding, &rec.IssueDate, &rec.DueDate,
		&rec.Status, &rec.AccountID, &rec.AccountCode, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("receivable not found: %w", err)
		}
		return nil, fmt.Errorf("query receivable: %w", err)
	}

	rec.CustomerID = customerID
	rec.TotalAmount = total
	rec.OutstandingAmount = outstanding

	payments, err := r.getPaymentsByReceivableID(ctx, exec, rec.ID)
	if err != nil {
		return nil, err
	}
	rec.Payments = payments

	return &rec, nil
}

func (r *PostgresReceivableRepo) GetByNumber(ctx context.Context, number string) (*receivable.Receivable, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT r.id, r.branch_id, r.receivable_number, r.customer_type, r.customer_id, r.customer_name,
		       r.invoice_reference, r.total_amount, r.outstanding_amount, r.issue_date, r.due_date,
		       r.status, r.account_id, a.code as account_code, COALESCE(r.notes, ''), r.created_at, r.updated_at
		FROM receivables r
		JOIN accounts a ON r.account_id = a.id
		WHERE r.receivable_number = $1
	`
	row := exec.QueryRowContext(ctx, query, number)

	var rec receivable.Receivable
	var customerID *uuid.UUID
	var total, outstanding decimal.Decimal

	err := row.Scan(
		&rec.ID, &rec.BranchID, &rec.ReceivableNumber, &rec.CustomerType, &customerID, &rec.CustomerName,
		&rec.InvoiceReference, &total, &outstanding, &rec.IssueDate, &rec.DueDate,
		&rec.Status, &rec.AccountID, &rec.AccountCode, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("receivable with number %s not found: %w", number, err)
		}
		return nil, fmt.Errorf("query receivable: %w", err)
	}

	rec.CustomerID = customerID
	rec.TotalAmount = total
	rec.OutstandingAmount = outstanding

	payments, err := r.getPaymentsByReceivableID(ctx, exec, rec.ID)
	if err != nil {
		return nil, err
	}
	rec.Payments = payments

	return &rec, nil
}

func (r *PostgresReceivableRepo) List(ctx context.Context, filter receivable.Filter) ([]receivable.Receivable, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `
		FROM receivables r
		JOIN accounts a ON r.account_id = a.id
		WHERE 1=1
	`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND r.branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.CustomerType != nil {
		baseQuery += fmt.Sprintf(" AND r.customer_type = $%d", argIdx)
		args = append(args, *filter.CustomerType)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND r.status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND r.issue_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND r.issue_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count receivables: %w", err)
	}

	selectQuery := `
		SELECT r.id, r.branch_id, r.receivable_number, r.customer_type, r.customer_id, r.customer_name,
		       r.invoice_reference, r.total_amount, r.outstanding_amount, r.issue_date, r.due_date,
		       r.status, r.account_id, a.code as account_code, COALESCE(r.notes, ''), r.created_at, r.updated_at
	` + baseQuery + " ORDER BY r.due_date ASC, r.created_at DESC"

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
		return nil, 0, fmt.Errorf("query receivables: %w", err)
	}
	defer rows.Close()

	var list []receivable.Receivable
	for rows.Next() {
		var rec receivable.Receivable
		var customerID *uuid.UUID
		var total, outstanding decimal.Decimal

		if err := rows.Scan(
			&rec.ID, &rec.BranchID, &rec.ReceivableNumber, &rec.CustomerType, &customerID, &rec.CustomerName,
			&rec.InvoiceReference, &total, &outstanding, &rec.IssueDate, &rec.DueDate,
			&rec.Status, &rec.AccountID, &rec.AccountCode, &rec.Notes, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan receivable row: %w", err)
		}
		rec.CustomerID = customerID
		rec.TotalAmount = total
		rec.OutstandingAmount = outstanding
		list = append(list, rec)
	}

	return list, totalCount, nil
}

func (r *PostgresReceivableRepo) UpdateOutstanding(ctx context.Context, id uuid.UUID, outstanding decimal.Decimal, status receivable.Status) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE receivables
		SET outstanding_amount = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	res, err := exec.ExecContext(ctx, query, outstanding, status, id)
	if err != nil {
		return fmt.Errorf("update receivable outstanding: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("receivable not found or not modified")
	}
	return nil
}

func (r *PostgresReceivableRepo) AddPayment(ctx context.Context, p *receivable.Payment) error {
	exec := GetExecutor(ctx, r.db)
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	query := `
		INSERT INTO receivable_payments (
			id, receivable_id, payment_number, amount, payment_date,
			cash_bank_account_code, reference_no, notes, journal_id, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		p.ID, p.ReceivableID, p.PaymentNumber, p.Amount, p.PaymentDate,
		p.CashBankAccountCode, p.ReferenceNo, p.Notes, p.JournalID, p.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert receivable payment: %w", err)
	}
	return nil
}

func (r *PostgresReceivableRepo) GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*receivable.AgingReport, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT 
			customer_name,
			COALESCE(SUM(CASE WHEN due_date >= $1 THEN outstanding_amount ELSE 0 END), 0) as current_amount,
			COALESCE(SUM(CASE WHEN ($1 - due_date) BETWEEN 1 AND 30 THEN outstanding_amount ELSE 0 END), 0) as days_1_30,
			COALESCE(SUM(CASE WHEN ($1 - due_date) BETWEEN 31 AND 60 THEN outstanding_amount ELSE 0 END), 0) as days_31_60,
			COALESCE(SUM(CASE WHEN ($1 - due_date) BETWEEN 61 AND 90 THEN outstanding_amount ELSE 0 END), 0) as days_61_90,
			COALESCE(SUM(CASE WHEN ($1 - due_date) > 90 THEN outstanding_amount ELSE 0 END), 0) as days_over_90,
			COALESCE(SUM(outstanding_amount), 0) as total_outstanding
		FROM receivables
		WHERE status IN ('OPEN', 'PARTIAL', 'OVERDUE')
		  AND ($2::uuid IS NULL OR branch_id = $2)
		GROUP BY customer_name
		ORDER BY customer_name ASC
	`
	rows, err := exec.QueryContext(ctx, query, asOfDate, branchID)
	if err != nil {
		return nil, fmt.Errorf("query receivable aging: %w", err)
	}
	defer rows.Close()

	var buckets []receivable.AgingBucket
	totalCurrent := decimal.Zero
	total1To30 := decimal.Zero
	total31To60 := decimal.Zero
	total61To90 := decimal.Zero
	totalOver90 := decimal.Zero
	grandTotal := decimal.Zero

	for rows.Next() {
		var b receivable.AgingBucket
		var cur, d30, d60, d90, dOver, tot decimal.Decimal

		if err := rows.Scan(&b.CustomerName, &cur, &d30, &d60, &d90, &dOver, &tot); err != nil {
			return nil, fmt.Errorf("scan aging bucket: %w", err)
		}
		b.Current = cur
		b.Days1To30 = d30
		b.Days31To60 = d60
		b.Days61To90 = d90
		b.DaysOver90 = dOver
		b.TotalOutstanding = tot

		totalCurrent = totalCurrent.Add(cur)
		total1To30 = total1To30.Add(d30)
		total31To60 = total31To60.Add(d60)
		total61To90 = total61To90.Add(d90)
		totalOver90 = totalOver90.Add(dOver)
		grandTotal = grandTotal.Add(tot)

		buckets = append(buckets, b)
	}

	return &receivable.AgingReport{
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

func (r *PostgresReceivableRepo) getPaymentsByReceivableID(ctx context.Context, exec QueryExecutor, receivableID uuid.UUID) ([]receivable.Payment, error) {
	query := `
		SELECT id, receivable_id, payment_number, amount, payment_date,
		       cash_bank_account_code, COALESCE(reference_no, ''), COALESCE(notes, ''),
		       journal_id, created_by, created_at
		FROM receivable_payments
		WHERE receivable_id = $1
		ORDER BY payment_date ASC
	`
	rows, err := exec.QueryContext(ctx, query, receivableID)
	if err != nil {
		return nil, fmt.Errorf("query receivable payments: %w", err)
	}
	defer rows.Close()

	var payments []receivable.Payment
	for rows.Next() {
		var p receivable.Payment
		var amount decimal.Decimal
		var jID, createdBy *uuid.UUID

		if err := rows.Scan(
			&p.ID, &p.ReceivableID, &p.PaymentNumber, &amount, &p.PaymentDate,
			&p.CashBankAccountCode, &p.ReferenceNo, &p.Notes,
			&jID, &createdBy, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan receivable payment: %w", err)
		}
		p.Amount = amount
		p.JournalID = jID
		p.CreatedBy = createdBy
		payments = append(payments, p)
	}
	return payments, nil
}

