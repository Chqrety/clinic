package payment

import (
	"context"
	"fmt"
)

type BankTransferProcessor struct {
	DefaultBankAccountCode string // e.g., "1-1113" (Bank Operasional)
}

func NewBankTransferProcessor(accountCode string) *BankTransferProcessor {
	if accountCode == "" {
		accountCode = "1-1113"
	}
	return &BankTransferProcessor{DefaultBankAccountCode: accountCode}
}

func (p *BankTransferProcessor) ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error) {
	return &PaymentResult{
		TransactionRef: fmt.Sprintf("TRF-%s", req.ReferenceNo),
		Status:         "SUCCESS",
		AccountCode:    p.DefaultBankAccountCode,
		Notes:          "Pembayaran transfer bank / EDC tervalidasi manual",
	}, nil
}

