package receivable

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CustomerType string

const (
	CustomerTypePatient   CustomerType = "PATIENT"
	CustomerTypeInsurance CustomerType = "INSURANCE"
	CustomerTypeCorporate CustomerType = "CORPORATE"
	CustomerTypeBPJS      CustomerType = "BPJS"
)

type Status string

const (
	StatusOpen     Status = "OPEN"
	StatusPartial  Status = "PARTIAL"
	StatusPaid     Status = "PAID"
	StatusOverdue  Status = "OVERDUE"
	StatusWriteOff Status = "WRITE_OFF"
)

type Receivable struct {
	shared.BaseEntity
	shared.BranchScoped
	ReceivableNumber  string          `json:"receivable_number"`
	CustomerType      CustomerType    `json:"customer_type"`
	CustomerID        *uuid.UUID      `json:"customer_id,omitempty"`
	CustomerName      string          `json:"customer_name"`
	InvoiceReference  string          `json:"invoice_reference"`
	TotalAmount       decimal.Decimal `json:"total_amount"`
	OutstandingAmount decimal.Decimal `json:"outstanding_amount"`
	IssueDate         time.Time       `json:"issue_date"`
	DueDate           time.Time       `json:"due_date"`
	Status            Status          `json:"status"`
	AccountID         uuid.UUID       `json:"account_id"`
	AccountCode       string          `json:"account_code,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	Payments          []Payment       `json:"payments,omitempty"`
}

type Payment struct {
	ID                  uuid.UUID       `json:"id"`
	ReceivableID        uuid.UUID       `json:"receivable_id"`
	PaymentNumber       string          `json:"payment_number"`
	Amount              decimal.Decimal `json:"amount"`
	PaymentDate         time.Time       `json:"payment_date"`
	CashBankAccountCode string          `json:"cash_bank_account_code"` // e.g. 1-1111 (Kasir), 1-1113 (BCA)
	ReferenceNo         string          `json:"reference_no,omitempty"`
	Notes               string          `json:"notes,omitempty"`
	JournalID           *uuid.UUID      `json:"journal_id,omitempty"`
	CreatedBy           *uuid.UUID      `json:"created_by,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
}

type AgingBucket struct {
	CustomerName      string          `json:"customer_name"`
	Current           decimal.Decimal `json:"current"` // Not yet due
	Days1To30         decimal.Decimal `json:"days_1_30"`
	Days31To60        decimal.Decimal `json:"days_31_60"`
	Days61To90        decimal.Decimal `json:"days_61_90"`
	DaysOver90        decimal.Decimal `json:"days_over_90"`
	TotalOutstanding  decimal.Decimal `json:"total_outstanding"`
}

type AgingReport struct {
	BranchID         *uuid.UUID      `json:"branch_id,omitempty"`
	AsOfDate         time.Time       `json:"as_of_date"`
	TotalCurrent     decimal.Decimal `json:"total_current"`
	Total1To30       decimal.Decimal `json:"total_1_30"`
	Total31To60      decimal.Decimal `json:"total_31_60"`
	Total61To90      decimal.Decimal `json:"total_61_90"`
	TotalOver90      decimal.Decimal `json:"total_over_90"`
	GrandTotal       decimal.Decimal `json:"grand_total"`
	Buckets          []AgingBucket   `json:"buckets"`
}

type Filter struct {
	BranchID     *uuid.UUID
	CustomerType *CustomerType
	Status       *Status
	StartDate    *time.Time
	EndDate      *time.Time
	Search       string
	SortBy       string
	SortOrder    string
	Limit        int
	Offset       int
}

type Repository interface {
	Create(ctx context.Context, r *Receivable) error
	GetByID(ctx context.Context, id uuid.UUID) (*Receivable, error)
	GetByNumber(ctx context.Context, number string) (*Receivable, error)
	List(ctx context.Context, filter Filter) ([]Receivable, int64, error)
	UpdateOutstanding(ctx context.Context, id uuid.UUID, outstanding decimal.Decimal, status Status) error
	AddPayment(ctx context.Context, p *Payment) error
	GetAgingReport(ctx context.Context, branchID *uuid.UUID, asOfDate time.Time) (*AgingReport, error)
}

