package cashbank

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeCash AccountType = "CASH"
	AccountTypeBank AccountType = "BANK"
)

type TransactionType string

const (
	TransactionTypeCashIn   TransactionType = "CASH_IN"
	TransactionTypeCashOut  TransactionType = "CASH_OUT"
	TransactionTypeTransfer TransactionType = "TRANSFER"
	TransactionTypeBankFee  TransactionType = "BANK_FEE"
)

// CashBankAccount represents a physical cash drawer or registered bank account.
type CashBankAccount struct {
	shared.BaseEntity
	shared.BranchScoped
	AccountID     uuid.UUID       `json:"account_id"` // Linked to Chart of Accounts (COA)
	AccountCode   string          `json:"account_code"`
	Code          string          `json:"code"` // e.g. "KASIR-01", "PETTY-CASH", "BCA-OPS"
	Name          string          `json:"name"`
	Type          AccountType     `json:"type"`
	BankName      string          `json:"bank_name,omitempty"`
	AccountNumber string          `json:"account_number,omitempty"`
	CurrentBalance decimal.Decimal `json:"current_balance"`
	IsActive      bool            `json:"is_active"`
}

// CashBankTransaction represents an operational cash in, cash out, or transfer event.
type CashBankTransaction struct {
	shared.BaseEntity
	shared.BranchScoped
	CashBankAccountID    uuid.UUID       `json:"cash_bank_account_id"`
	CashBankAccountCode  string          `json:"cash_bank_account_code,omitempty"`
	CashBankAccountName  string          `json:"cash_bank_account_name,omitempty"`
	TransactionNumber    string          `json:"transaction_number"`
	TransactionType      TransactionType `json:"transaction_type"`
	Amount               decimal.Decimal `json:"amount"`
	TransactionDate      time.Time       `json:"transaction_date"`
	CounterpartAccountID uuid.UUID       `json:"counterpart_account_id"` // Lawan akun (e.g. Beban Listrik, Pendapatan Lain, atau Kas Tujuan)
	CounterpartCode      string          `json:"counterpart_code,omitempty"`
	CounterpartName      string          `json:"counterpart_name,omitempty"`
	Description          string          `json:"description"`
	ReferenceNo          string          `json:"reference_no,omitempty"`
	JournalID            *uuid.UUID      `json:"journal_id,omitempty"` // Audit trail drill-down link
}

type TransactionFilter struct {
	BranchID          *uuid.UUID
	CashBankAccountID *uuid.UUID
	TransactionType   *TransactionType
	StartDate         *time.Time
	EndDate           *time.Time
	Limit             int
	Offset            int
}

type Repository interface {
	CreateAccount(ctx context.Context, acc *CashBankAccount) error
	GetAccountByID(ctx context.Context, id uuid.UUID) (*CashBankAccount, error)
	GetAccountByCode(ctx context.Context, code string) (*CashBankAccount, error)
	ListAccounts(ctx context.Context, branchID *uuid.UUID, accType *AccountType) ([]CashBankAccount, error)
	CreateTransaction(ctx context.Context, tx *CashBankTransaction) error
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*CashBankTransaction, error)
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]CashBankTransaction, int64, error)
}

