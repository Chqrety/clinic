package accounting

import (
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeCOGS      AccountType = "COGS"
	AccountTypeExpense   AccountType = "EXPENSE"
)

type NormalBalance string

const (
	NormalBalanceDebit  NormalBalance = "DEBIT"
	NormalBalanceCredit NormalBalance = "CREDIT"
)

type JournalType string

const (
	JournalTypeGeneral          JournalType = "GENERAL"
	JournalTypeCashReceipt      JournalType = "CASH_RECEIPT"
	JournalTypeCashDisbursement JournalType = "CASH_DISBURSEMENT"
	JournalTypeBank             JournalType = "BANK"
	JournalTypeSales            JournalType = "SALES"
	JournalTypePurchase         JournalType = "PURCHASE"
	JournalTypeAdjusting        JournalType = "ADJUSTING"
	JournalTypeClosing          JournalType = "CLOSING"
	JournalTypeReversal         JournalType = "REVERSAL"
)

type JournalStatus string

const (
	JournalStatusDraft    JournalStatus = "DRAFT"
	JournalStatusPosted   JournalStatus = "POSTED"
	JournalStatusReversed JournalStatus = "REVERSED"
	JournalStatusVoid     JournalStatus = "VOID"
)

// Account represents the Chart of Accounts (COA) entity.
type Account struct {
	shared.BaseEntity
	Code          string        `json:"code"`
	Name          string        `json:"name"`
	Type          AccountType   `json:"type"`
	NormalBalance NormalBalance `json:"normal_balance"`
	ParentID      *uuid.UUID    `json:"parent_id,omitempty"`
	Level         int           `json:"level"`
	IsHeader      bool          `json:"is_header"`
	IsActive      bool          `json:"is_active"`
	Description   string        `json:"description,omitempty"`
}

// JournalEntry represents a double-entry transaction header.
type JournalEntry struct {
	shared.BaseEntity
	shared.BranchScoped
	JournalNumber   string          `json:"journal_number"`
	JournalType     JournalType     `json:"journal_type"`
	TransactionDate time.Time       `json:"transaction_date"`
	PostingDate     *time.Time      `json:"posting_date,omitempty"`
	Status          JournalStatus   `json:"status"`
	ReferenceType   string          `json:"reference_type"` // e.g. "INVOICE", "PURCHASE", "CASH_IN"
	ReferenceID     string          `json:"reference_id"`   // External source UUID or code
	Description     string          `json:"description"`
	TotalDebit      decimal.Decimal `json:"total_debit"`
	TotalCredit     decimal.Decimal `json:"total_credit"`
	PostedBy        *uuid.UUID      `json:"posted_by,omitempty"`
	Lines           []JournalLine   `json:"lines"`
}

// JournalLine represents each single debit/credit line item in a journal entry.
type JournalLine struct {
	ID             uuid.UUID       `json:"id"`
	JournalEntryID uuid.UUID       `json:"journal_entry_id"`
	AccountID      uuid.UUID       `json:"account_id"`
	AccountCode    string          `json:"account_code,omitempty"`
	AccountName    string          `json:"account_name,omitempty"`
	Debit          decimal.Decimal `json:"debit"`
	Credit         decimal.Decimal `json:"credit"`
	Description    string          `json:"description"`
	LineOrder      int             `json:"line_order"`
}

const (
	PeriodStatusOpen   = "OPEN"
	PeriodStatusClosed = "CLOSED"
	PeriodStatusLocked = "LOCKED"
)

// FiscalPeriod represents the accounting period state for period locking.
type FiscalPeriod struct {
	shared.BaseEntity
	shared.BranchScoped
	Year      int        `json:"year"`
	Month     int        `json:"month"`
	Status    string     `json:"status"` // "OPEN", "CLOSED", "LOCKED"
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	ClosedBy  *uuid.UUID `json:"closed_by,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

// AccountingMapping defines dynamic transaction-to-account rule mappings.
type AccountingMapping struct {
	shared.BaseEntity
	TransactionCode string    `json:"transaction_code"` // e.g., "REV_CONSULTATION", "REV_PHARMACY", "HPP_MEDICINE"
	DebitAccountID  uuid.UUID `json:"debit_account_id"`
	CreditAccountID uuid.UUID `json:"credit_account_id"`
	Description     string    `json:"description"`
	IsActive        bool      `json:"is_active"`
}

