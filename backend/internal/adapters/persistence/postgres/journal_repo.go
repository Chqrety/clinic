package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresJournalRepo struct {
	db *sql.DB
}

func NewPostgresJournalRepo(db *sql.DB) *PostgresJournalRepo {
	return &PostgresJournalRepo{db: db}
}

func (r *PostgresJournalRepo) Create(ctx context.Context, entry *accounting.JournalEntry) error {
	exec := GetExecutor(ctx, r.db)

	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}

	headerQuery := `
		INSERT INTO journal_entries (
			id, branch_id, journal_number, journal_type, transaction_date,
			posting_date, status, reference_type, reference_id,
			total_debit, total_credit, description, posted_by, created_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`

	_, err := exec.ExecContext(ctx, headerQuery,
		entry.ID, entry.BranchID, entry.JournalNumber, entry.JournalType, entry.TransactionDate,
		entry.PostingDate, entry.Status, entry.ReferenceType, entry.ReferenceID,
		entry.TotalDebit, entry.TotalCredit, entry.Description, entry.PostedBy, entry.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert journal header: %w", err)
	}

	lineQuery := `
		INSERT INTO journal_lines (
			id, journal_entry_id, account_id, debit, credit, description, line_order, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP
		)
	`

	for idx, line := range entry.Lines {
		lineID := line.ID
		if lineID == uuid.Nil {
			lineID = uuid.New()
		}
		_, err := exec.ExecContext(ctx, lineQuery,
			lineID, entry.ID, line.AccountID, line.Debit, line.Credit, line.Description, idx+1,
		)
		if err != nil {
			return fmt.Errorf("insert journal line [%d]: %w", idx, err)
		}
	}

	return nil
}

func (r *PostgresJournalRepo) GetByID(ctx context.Context, id uuid.UUID) (*accounting.JournalEntry, error) {
	query := `
		SELECT id, branch_id, journal_number, journal_type, transaction_date,
		       posting_date, status, reference_type, reference_id,
		       total_debit, total_credit, description, posted_by, created_by, created_at, updated_at
		FROM journal_entries
		WHERE id = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, id)

	var entry accounting.JournalEntry
	var postingDate *time.Time
	var postedBy, createdBy *uuid.UUID

	err := row.Scan(
		&entry.ID, &entry.BranchID, &entry.JournalNumber, &entry.JournalType, &entry.TransactionDate,
		&postingDate, &entry.Status, &entry.ReferenceType, &entry.ReferenceID,
		&entry.TotalDebit, &entry.TotalCredit, &entry.Description, &postedBy, &createdBy,
		&entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("journal entry not found: %w", err)
		}
		return nil, fmt.Errorf("query journal header: %w", err)
	}

	entry.PostingDate = postingDate
	entry.PostedBy = postedBy
	entry.CreatedBy = createdBy

	lines, err := r.getLinesByEntryID(ctx, exec, entry.ID)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines

	return &entry, nil
}

func (r *PostgresJournalRepo) GetByJournalNumber(ctx context.Context, number string) (*accounting.JournalEntry, error) {
	query := `
		SELECT id, branch_id, journal_number, journal_type, transaction_date,
		       posting_date, status, reference_type, reference_id,
		       total_debit, total_credit, description, posted_by, created_by, created_at, updated_at
		FROM journal_entries
		WHERE journal_number = $1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, number)

	var entry accounting.JournalEntry
	var postingDate *time.Time
	var postedBy, createdBy *uuid.UUID

	err := row.Scan(
		&entry.ID, &entry.BranchID, &entry.JournalNumber, &entry.JournalType, &entry.TransactionDate,
		&postingDate, &entry.Status, &entry.ReferenceType, &entry.ReferenceID,
		&entry.TotalDebit, &entry.TotalCredit, &entry.Description, &postedBy, &createdBy,
		&entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("journal with number %s not found: %w", number, err)
		}
		return nil, fmt.Errorf("query journal by number: %w", err)
	}

	entry.PostingDate = postingDate
	entry.PostedBy = postedBy
	entry.CreatedBy = createdBy

	lines, err := r.getLinesByEntryID(ctx, exec, entry.ID)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines

	return &entry, nil
}

func (r *PostgresJournalRepo) GetByReference(ctx context.Context, refType, refID string) (*accounting.JournalEntry, error) {
	query := `
		SELECT id, branch_id, journal_number, journal_type, transaction_date,
		       posting_date, status, reference_type, reference_id,
		       total_debit, total_credit, description, posted_by, created_by, created_at, updated_at
		FROM journal_entries
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	exec := GetExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, query, refType, refID)

	var entry accounting.JournalEntry
	var postingDate *time.Time
	var postedBy, createdBy *uuid.UUID

	err := row.Scan(
		&entry.ID, &entry.BranchID, &entry.JournalNumber, &entry.JournalType, &entry.TransactionDate,
		&postingDate, &entry.Status, &entry.ReferenceType, &entry.ReferenceID,
		&entry.TotalDebit, &entry.TotalCredit, &entry.Description, &postedBy, &createdBy,
		&entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("journal with reference %s:%s not found: %w", refType, refID, err)
		}
		return nil, fmt.Errorf("query journal by reference: %w", err)
	}

	entry.PostingDate = postingDate
	entry.PostedBy = postedBy
	entry.CreatedBy = createdBy

	lines, err := r.getLinesByEntryID(ctx, exec, entry.ID)
	if err != nil {
		return nil, err
	}
	entry.Lines = lines

	return &entry, nil
}

func (r *PostgresJournalRepo) List(ctx context.Context, filter accounting.JournalFilter) ([]accounting.JournalEntry, int64, error) {
	baseQuery := `FROM journal_entries WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.Type != nil {
		baseQuery += fmt.Sprintf(" AND journal_type = $%d", argIdx)
		args = append(args, *filter.Type)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND transaction_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND transaction_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}

	exec := GetExecutor(ctx, r.db)

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count journals: %w", err)
	}

	selectQuery := `
		SELECT id, branch_id, journal_number, journal_type, transaction_date,
		       posting_date, status, reference_type, reference_id,
		       total_debit, total_credit, description, posted_by, created_by, created_at, updated_at
	` + baseQuery + " ORDER BY transaction_date DESC, created_at DESC"

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
		return nil, 0, fmt.Errorf("query journals: %w", err)
	}
	defer rows.Close()

	var entries []accounting.JournalEntry
	for rows.Next() {
		var entry accounting.JournalEntry
		var postingDate *time.Time
		var postedBy, createdBy *uuid.UUID

		if err := rows.Scan(
			&entry.ID, &entry.BranchID, &entry.JournalNumber, &entry.JournalType, &entry.TransactionDate,
			&postingDate, &entry.Status, &entry.ReferenceType, &entry.ReferenceID,
			&entry.TotalDebit, &entry.TotalCredit, &entry.Description, &postedBy, &createdBy,
			&entry.CreatedAt, &entry.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan journal row: %w", err)
		}
		entry.PostingDate = postingDate
		entry.PostedBy = postedBy
		entry.CreatedBy = createdBy
		entries = append(entries, entry)
	}

	return entries, totalCount, nil
}

func (r *PostgresJournalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status accounting.JournalStatus, postedBy *uuid.UUID) error {
	exec := GetExecutor(ctx, r.db)
	var query string
	var args []any

	if status == accounting.JournalStatusPosted {
		query = `
			UPDATE journal_entries
			SET status = $1, posting_date = CURRENT_TIMESTAMP, posted_by = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $3
		`
		args = []any{status, postedBy, id}
	} else {
		query = `
			UPDATE journal_entries
			SET status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
		`
		args = []any{status, id}
	}

	res, err := exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update journal status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("journal entry not found")
	}
	return nil
}

func (r *PostgresJournalRepo) getLinesByEntryID(ctx context.Context, exec QueryExecutor, entryID uuid.UUID) ([]accounting.JournalLine, error) {
	query := `
		SELECT jl.id, jl.journal_entry_id, jl.account_id, a.code, a.name,
		       jl.debit, jl.credit, COALESCE(jl.description, ''), jl.line_order
		FROM journal_lines jl
		JOIN accounts a ON jl.account_id = a.id
		WHERE jl.journal_entry_id = $1
		ORDER BY jl.line_order ASC
	`
	rows, err := exec.QueryContext(ctx, query, entryID)
	if err != nil {
		return nil, fmt.Errorf("query journal lines: %w", err)
	}
	defer rows.Close()

	var lines []accounting.JournalLine
	for rows.Next() {
		var line accounting.JournalLine
		var debit, credit decimal.Decimal
		if err := rows.Scan(
			&line.ID, &line.JournalEntryID, &line.AccountID, &line.AccountCode, &line.AccountName,
			&debit, &credit, &line.Description, &line.LineOrder,
		); err != nil {
			return nil, fmt.Errorf("scan journal line: %w", err)
		}
		line.Debit = debit
		line.Credit = credit
		lines = append(lines, line)
	}

	return lines, nil
}
