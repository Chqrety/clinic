package dto

import (
	"time"

	"clinic-api/internal/domain/cashbank"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateCashBankAccountRequest struct {
	AccountID     uuid.UUID            `json:"account_id"` // COA Account UUID
	Code          string               `json:"code"`
	Name          string               `json:"name"`
	Type          cashbank.AccountType `json:"type"` // CASH, BANK
	BankName      string               `json:"bank_name,omitempty"`
	AccountNumber string               `json:"account_number,omitempty"`
}

type CashBankAccountResponse struct {
	ID            uuid.UUID            `json:"id"`
	BranchID      uuid.UUID            `json:"branch_id"`
	AccountID     uuid.UUID            `json:"account_id"`
	AccountCode   string               `json:"account_code"`
	Code          string               `json:"code"`
	Name          string               `json:"name"`
	Type          cashbank.AccountType `json:"type"`
	BankName      string               `json:"bank_name,omitempty"`
	AccountNumber string               `json:"account_number,omitempty"`
	IsActive      bool                 `json:"is_active"`
	CreatedAt     time.Time            `json:"created_at"`
}

type RecordCashInRequest struct {
	CashBankAccountID    uuid.UUID       `json:"cash_bank_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	TransactionDate      string          `json:"transaction_date"` // YYYY-MM-DD
	CounterpartAccountID uuid.UUID       `json:"counterpart_account_id"`
	Description          string          `json:"description"`
	ReferenceNo          string          `json:"reference_no,omitempty"`
}

type RecordCashOutRequest struct {
	CashBankAccountID    uuid.UUID       `json:"cash_bank_account_id"`
	Amount               decimal.Decimal `json:"amount"`
	TransactionDate      string          `json:"transaction_date"` // YYYY-MM-DD
	CounterpartAccountID uuid.UUID       `json:"counterpart_account_id"`
	Description          string          `json:"description"`
	ReferenceNo          string          `json:"reference_no,omitempty"`
}

type TransferCashBankRequest struct {
	SourceAccountID uuid.UUID       `json:"source_account_id"`
	TargetAccountID uuid.UUID       `json:"target_account_id"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionDate string          `json:"transaction_date"` // YYYY-MM-DD
	Description     string          `json:"description"`
	ReferenceNo     string          `json:"reference_no,omitempty"`
}

type CashBankTransactionResponse struct {
	ID                   uuid.UUID                `json:"id"`
	BranchID             uuid.UUID                `json:"branch_id"`
	CashBankAccountID    uuid.UUID                `json:"cash_bank_account_id"`
	CashBankAccountCode  string                   `json:"cash_bank_account_code,omitempty"`
	CashBankAccountName  string                   `json:"cash_bank_account_name,omitempty"`
	TransactionNumber    string                   `json:"transaction_number"`
	TransactionType      cashbank.TransactionType `json:"transaction_type"`
	Amount               decimal.Decimal          `json:"amount"`
	TransactionDate      string                   `json:"transaction_date"`
	CounterpartAccountID uuid.UUID                `json:"counterpart_account_id"`
	CounterpartCode      string                   `json:"counterpart_code,omitempty"`
	CounterpartName      string                   `json:"counterpart_name,omitempty"`
	Description          string                   `json:"description"`
	ReferenceNo          string                   `json:"reference_no,omitempty"`
	JournalID            *uuid.UUID               `json:"journal_id,omitempty"`
	CreatedAt            time.Time                `json:"created_at"`
}

type ClosePeriodRequest struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Notes string `json:"notes,omitempty"`
}

type ReopenPeriodRequest struct {
	Reason string `json:"reason"`
}

