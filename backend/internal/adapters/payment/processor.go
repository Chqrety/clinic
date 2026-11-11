package payment

import (
	"context"
	"errors"
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

var (
	ErrUnsupportedPaymentMethod = errors.New("unsupported payment method")
	ErrPaymentGatewayDeferred   = errors.New("payment gateway integration is currently deferred")
)

type Method string

const (
	MethodCash         Method = "CASH"
	MethodBankTransfer Method = "BANK_TRANSFER"
	MethodEDC          Method = "EDC"
	MethodReceivable   Method = "RECEIVABLE" // Piutang Pasien / Asuransi
	MethodGatewayQRIS  Method = "GATEWAY_QRIS" // [DEFERRED]
	MethodGatewayVA    Method = "GATEWAY_VA"   // [DEFERRED]
)

type PaymentRequest struct {
	PaymentID   uuid.UUID
	BranchID    uuid.UUID
	InvoiceID   uuid.UUID
	Amount      decimal.Decimal
	Method      Method
	ReferenceNo string
	ReceivedBy  uuid.UUID
	PaidAt      time.Time
}

type PaymentResult struct {
	TransactionRef string
	Status         string // SUCCESS, PENDING, FAILED
	AccountCode    string // Account to debit (e.g. Kas Kasir: 1-1111, Bank BCA: 1-1113)
	Notes          string
}

// PaymentProcessor defines the contract for processing payments.
type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
}

