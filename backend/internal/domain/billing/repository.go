package billing

import (
	"context"
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type InvoiceFilter struct {
	BranchID      *uuid.UUID
	PaymentStatus *PaymentStatus
	InvoiceType   *InvoiceType
	StartDate     *time.Time
	EndDate       *time.Time
	Limit         int
	Offset        int
}

type Repository interface {
	Create(ctx context.Context, inv *Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	GetByNumber(ctx context.Context, number string) (*Invoice, error)
	GetByQueueID(ctx context.Context, queueID int64) ([]Invoice, error)
	List(ctx context.Context, filter InvoiceFilter) ([]Invoice, int64, error)
	UpdatePayment(ctx context.Context, id uuid.UUID, paidAmount decimal.Decimal, status PaymentStatus, journalID *uuid.UUID) error
	AddPaymentRecord(ctx context.Context, record *PaymentRecord) error
}

