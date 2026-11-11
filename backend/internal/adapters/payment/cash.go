package payment

import (
	"context"
	"fmt"
)

type CashProcessor struct {
	DefaultCashAccountCode string // e.g., "1-1111" (Kas Kasir)
}

func NewCashProcessor(accountCode string) *CashProcessor {
	if accountCode == "" {
		accountCode = "1-1111"
	}
	return &CashProcessor{DefaultCashAccountCode: accountCode}
}

func (p *CashProcessor) ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error) {
	return &PaymentResult{
		TransactionRef: fmt.Sprintf("CASH-%s", req.PaymentID.String()[:8]),
		Status:         "SUCCESS",
		AccountCode:    p.DefaultCashAccountCode,
		Notes:          "Pembayaran tunai kasir terverifikasi",
	}, nil
}

