package gateway

import (
	"context"
	"errors"

	"clinic-api/internal/adapters/payment"
)

var (
	ErrPaymentGatewayPostponed = errors.New("payment gateway integration is postponed in this phase")
)

// PaymentGatewayProvider defines the pluggable contract for payment gateways
// (e.g. Midtrans, Xendit, QRIS Network) to be implemented in future phase.
type PaymentGatewayProvider interface {
	CreateTransaction(ctx context.Context, req payment.PaymentRequest) (redirectURL string, snapToken string, err error)
	VerifyWebhookSignature(payload []byte, signature string) bool
	CheckStatus(ctx context.Context, transactionID string) (status string, err error)
}

// DeferredGatewayProcessor acts as a stub while payment gateway is postponed.
type DeferredGatewayProcessor struct{}

func NewDeferredGatewayProcessor() *DeferredGatewayProcessor {
	return &DeferredGatewayProcessor{}
}

func (d *DeferredGatewayProcessor) ProcessPayment(ctx context.Context, req payment.PaymentRequest) (*payment.PaymentResult, error) {
	// Returns error indicating payment gateway is deferred per project phase specification.
	return nil, ErrPaymentGatewayPostponed
}

