package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-api/internal/domain/cashbank"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresCashBankRepo struct {
	db *sql.DB
}

func NewPostgresCashBankRepo(db *sql.DB) *PostgresCashBankRepo {
	return &PostgresCashBankRepo{db: db}
}

func (r *PostgresCashBankRepo) CreateAccount(ctx context.Context, acc *cashbank.CashBankAccount) error {
	exec := GetExecutor(ctx, r.db)
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}

	query := `
		INSERT INTO cash_bank_accounts (
			id, branch_id, account_id, code, name, type, bank_name, account_number, is_active,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		acc.ID, acc.BranchID, acc.AccountID, acc.Code, acc.Name,
		acc.Type, acc.BankName, acc.AccountNumber, acc.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert cash_bank_account: %w", err)
	}
	return nil
}

func (r *PostgresCashBankRepo) GetAccountByID(ctx context.Context, id uuid.UUID) (*cashbank.CashBankAccount, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT cba.id, cba.branch_id, cba.account_id, a.code as account_code,
		       cba.code, cba.name, cba.type, COALESCE(cba.bank_name, ''), COALESCE(cba.account_number, ''),
		       cba.is_active, cba.created_at, cba.updated_at
		FROM cash_bank_accounts cba
		JOIN accounts a ON cba.account_id = a.id
		WHERE cba.id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)

	var acc cashbank.CashBankAccount
	err := row.Scan(
		&acc.ID, &acc.BranchID, &acc.AccountID, &acc.AccountCode,
		&acc.Code, &acc.Name, &acc.Type, &acc.BankName, &acc.AccountNumber,
		&acc.IsActive, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cash bank account not found: %w", err)
		}
		return nil, fmt.Errorf("query cash bank account by id: %w", err)
	}

	return &acc, nil
}

func (r *PostgresCashBankRepo) GetAccountByCode(ctx context.Context, code string) (*cashbank.CashBankAccount, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT cba.id, cba.branch_id, cba.account_id, a.code as account_code,
		       cba.code, cba.name, cba.type, COALESCE(cba.bank_name, ''), COALESCE(cba.account_number, ''),
		       cba.is_active, cba.created_at, cba.updated_at
		FROM cash_bank_accounts cba
		JOIN accounts a ON cba.account_id = a.id
		WHERE cba.code = $1
	`
	row := exec.QueryRowContext(ctx, query, code)

	var acc cashbank.CashBankAccount
	err := row.Scan(
		&acc.ID, &acc.BranchID, &acc.AccountID, &acc.AccountCode,
		&acc.Code, &acc.Name, &acc.Type, &acc.BankName, &acc.AccountNumber,
		&acc.IsActive, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cash bank account with code %s not found: %w", code, err)
		}
		return nil, fmt.Errorf("query cash bank account by code: %w", err)
	}

	return &acc, nil
}

func (r *PostgresCashBankRepo) ListAccounts(ctx context.Context, filter cashbank.CashAccountFilter) ([]cashbank.CashBankAccount, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `
		FROM cash_bank_accounts cba
		JOIN accounts a ON cba.account_id = a.id
		WHERE cba.is_active = TRUE
	`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND cba.branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.Type != nil {
		baseQuery += fmt.Sprintf(" AND cba.type = $%d", argIdx)
		args = append(args, *filter.Type)
		argIdx++
	}
	if filter.Search != "" {
		pat := "%" + strings.ToLower(filter.Search) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(cba.code) LIKE $%d OR LOWER(cba.name) LIKE $%d OR LOWER(COALESCE(cba.bank_name, '')) LIKE $%d OR LOWER(COALESCE(cba.account_number, '')) LIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, pat)
		argIdx++
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count cash bank accounts: %w", err)
	}

	sortCol := "cba.code"
	sortOrder := "ASC"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "code", "cba.code":
			sortCol = "cba.code"
		case "name", "cba.name":
			sortCol = "cba.name"
		case "type", "cba.type":
			sortCol = "cba.type"
		case "created_at", "cba.created_at":
			sortCol = "cba.created_at"
		}
	}
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	query := `
		SELECT cba.id, cba.branch_id, cba.account_id, a.code as account_code,
		       cba.code, cba.name, cba.type, COALESCE(cba.bank_name, ''), COALESCE(cba.account_number, ''),
		       cba.is_active, cba.created_at, cba.updated_at
	` + baseQuery + fmt.Sprintf(" ORDER BY %s %s, cba.id %s", sortCol, sortOrder, sortOrder)

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list cash bank accounts: %w", err)
	}
	defer rows.Close()

	var list []cashbank.CashBankAccount
	for rows.Next() {
		var acc cashbank.CashBankAccount
		if err := rows.Scan(
			&acc.ID, &acc.BranchID, &acc.AccountID, &acc.AccountCode,
			&acc.Code, &acc.Name, &acc.Type, &acc.BankName, &acc.AccountNumber,
			&acc.IsActive, &acc.CreatedAt, &acc.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan cash bank account: %w", err)
		}
		list = append(list, acc)
	}

	return list, totalCount, nil
}

func (r *PostgresCashBankRepo) CreateTransaction(ctx context.Context, tx *cashbank.CashBankTransaction) error {
	exec := GetExecutor(ctx, r.db)
	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}

	query := `
		INSERT INTO cash_bank_transactions (
			id, branch_id, cash_bank_account_id, transaction_number, transaction_type,
			amount, transaction_date, counterpart_account_id, description, reference_no,
			journal_id, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		tx.ID, tx.BranchID, tx.CashBankAccountID, tx.TransactionNumber, tx.TransactionType,
		tx.Amount, tx.TransactionDate, tx.CounterpartAccountID, tx.Description, tx.ReferenceNo,
		tx.JournalID, tx.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("insert cash bank transaction: %w", err)
	}
	return nil
}

func (r *PostgresCashBankRepo) GetTransactionByID(ctx context.Context, id uuid.UUID) (*cashbank.CashBankTransaction, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT t.id, t.branch_id, t.cash_bank_account_id, cba.code as cba_code, cba.name as cba_name,
		       t.transaction_number, t.transaction_type, t.amount, t.transaction_date,
		       t.counterpart_account_id, ca.code as counterpart_code, ca.name as counterpart_name,
		       t.description, COALESCE(t.reference_no, ''), t.journal_id, t.created_by, t.created_at
		FROM cash_bank_transactions t
		JOIN cash_bank_accounts cba ON t.cash_bank_account_id = cba.id
		JOIN accounts ca ON t.counterpart_account_id = ca.id
		WHERE t.id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)

	var tx cashbank.CashBankTransaction
	var amount decimal.Decimal
	var journalID, createdBy *uuid.UUID

	err := row.Scan(
		&tx.ID, &tx.BranchID, &tx.CashBankAccountID, &tx.CashBankAccountCode, &tx.CashBankAccountName,
		&tx.TransactionNumber, &tx.TransactionType, &amount, &tx.TransactionDate,
		&tx.CounterpartAccountID, &tx.CounterpartCode, &tx.CounterpartName,
		&tx.Description, &tx.ReferenceNo, &journalID, &createdBy, &tx.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found: %w", err)
		}
		return nil, fmt.Errorf("query transaction by id: %w", err)
	}

	tx.Amount = amount
	tx.JournalID = journalID
	tx.CreatedBy = createdBy

	return &tx, nil
}

func (r *PostgresCashBankRepo) ListTransactions(ctx context.Context, filter cashbank.TransactionFilter) ([]cashbank.CashBankTransaction, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `
		FROM cash_bank_transactions t
		JOIN cash_bank_accounts cba ON t.cash_bank_account_id = cba.id
		JOIN accounts ca ON t.counterpart_account_id = ca.id
		WHERE 1=1
	`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND t.branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.CashBankAccountID != nil {
		baseQuery += fmt.Sprintf(" AND t.cash_bank_account_id = $%d", argIdx)
		args = append(args, *filter.CashBankAccountID)
		argIdx++
	}
	if filter.TransactionType != nil {
		baseQuery += fmt.Sprintf(" AND t.transaction_type = $%d", argIdx)
		args = append(args, *filter.TransactionType)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND t.transaction_date >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND t.transaction_date <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if filter.Search != "" {
		pat := "%" + strings.ToLower(filter.Search) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(t.transaction_number) LIKE $%d OR LOWER(t.description) LIKE $%d OR LOWER(COALESCE(t.reference_no, '')) LIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, pat)
		argIdx++
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count cash bank transactions: %w", err)
	}

	sortCol := "t.transaction_date"
	sortOrder := "DESC"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "transaction_date", "t.transaction_date":
			sortCol = "t.transaction_date"
		case "transaction_number", "t.transaction_number":
			sortCol = "t.transaction_number"
		case "amount", "t.amount":
			sortCol = "t.amount"
		case "created_at", "t.created_at":
			sortCol = "t.created_at"
		}
	}
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	selectQuery := `
		SELECT t.id, t.branch_id, t.cash_bank_account_id, cba.code as cba_code, cba.name as cba_name,
		       t.transaction_number, t.transaction_type, t.amount, t.transaction_date,
		       t.counterpart_account_id, ca.code as counterpart_code, ca.name as counterpart_name,
		       t.description, COALESCE(t.reference_no, ''), t.journal_id, t.created_by, t.created_at
	` + baseQuery + fmt.Sprintf(" ORDER BY %s %s, t.id %s", sortCol, sortOrder, sortOrder)

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
		return nil, 0, fmt.Errorf("query cash bank transactions: %w", err)
	}
	defer rows.Close()

	var list []cashbank.CashBankTransaction
	for rows.Next() {
		var tx cashbank.CashBankTransaction
		var amount decimal.Decimal
		var journalID, createdBy *uuid.UUID

		if err := rows.Scan(
			&tx.ID, &tx.BranchID, &tx.CashBankAccountID, &tx.CashBankAccountCode, &tx.CashBankAccountName,
			&tx.TransactionNumber, &tx.TransactionType, &amount, &tx.TransactionDate,
			&tx.CounterpartAccountID, &tx.CounterpartCode, &tx.CounterpartName,
			&tx.Description, &tx.ReferenceNo, &journalID, &createdBy, &tx.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan transaction row: %w", err)
		}
		tx.Amount = amount
		tx.JournalID = journalID
		tx.CreatedBy = createdBy
		list = append(list, tx)
	}

	return list, totalCount, nil
}
