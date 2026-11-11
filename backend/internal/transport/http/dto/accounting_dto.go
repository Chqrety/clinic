package dto

import (
	"time"

	"clinic-api/internal/domain/accounting"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	Code          string                   `json:"code"`
	Name          string                   `json:"name"`
	Type          accounting.AccountType   `json:"type"`
	NormalBalance accounting.NormalBalance `json:"normal_balance"`
	ParentID      *uuid.UUID               `json:"parent_id,omitempty"`
	Level         int                      `json:"level"`
	IsHeader      bool                     `json:"is_header"`
	Description   string                   `json:"description,omitempty"`
}

type AccountResponse struct {
	ID            uuid.UUID                `json:"id"`
	Code          string                   `json:"code"`
	Name          string                   `json:"name"`
	Type          accounting.AccountType   `json:"type"`
	NormalBalance accounting.NormalBalance `json:"normal_balance"`
	ParentID      *uuid.UUID               `json:"parent_id,omitempty"`
	Level         int                      `json:"level"`
	IsHeader      bool                     `json:"is_header"`
	IsActive      bool                     `json:"is_active"`
	Description   string                   `json:"description,omitempty"`
}

type JournalLineRequest struct {
	AccountID   uuid.UUID       `json:"account_id,omitempty"`
	AccountCode string          `json:"account_code,omitempty"`
	Debit       decimal.Decimal `json:"debit"`
	Credit      decimal.Decimal `json:"credit"`
	Description string          `json:"description,omitempty"`
}

type CreateJournalRequest struct {
	JournalType     accounting.JournalType `json:"journal_type"`
	TransactionDate string                 `json:"transaction_date"` // YYYY-MM-DD
	ReferenceType   string                 `json:"reference_type"`
	ReferenceID     string                 `json:"reference_id"`
	Description     string                 `json:"description"`
	AutoPost        bool                   `json:"auto_post"`
	Lines           []JournalLineRequest   `json:"lines"`
}

type ReverseJournalRequest struct {
	Reason string `json:"reason"`
}

type JournalLineResponse struct {
	ID          uuid.UUID       `json:"id"`
	AccountID   uuid.UUID       `json:"account_id"`
	AccountCode string          `json:"account_code"`
	AccountName string          `json:"account_name"`
	Debit       decimal.Decimal `json:"debit"`
	Credit      decimal.Decimal `json:"credit"`
	Description string          `json:"description"`
	LineOrder   int             `json:"line_order"`
}

type JournalResponse struct {
	ID              uuid.UUID              `json:"id"`
	BranchID        uuid.UUID              `json:"branch_id"`
	JournalNumber   string                 `json:"journal_number"`
	JournalType     accounting.JournalType `json:"journal_type"`
	TransactionDate string                 `json:"transaction_date"`
	PostingDate     *time.Time             `json:"posting_date,omitempty"`
	Status          accounting.JournalStatus `json:"status"`
	ReferenceType   string                 `json:"reference_type"`
	ReferenceID     string                 `json:"reference_id"`
	Description     string                 `json:"description"`
	TotalDebit      decimal.Decimal        `json:"total_debit"`
	TotalCredit     decimal.Decimal        `json:"total_credit"`
	PostedBy        *uuid.UUID             `json:"posted_by,omitempty"`
	CreatedBy       *uuid.UUID             `json:"created_by,omitempty"`
	Lines           []JournalLineResponse  `json:"lines,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

