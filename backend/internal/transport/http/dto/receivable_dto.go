package dto

import (
	"time"

	"clinic-api/internal/domain/receivable"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateReceivableRequest struct {
	CustomerType     receivable.CustomerType `json:"customer_type"`
	CustomerID       *uuid.UUID              `json:"customer_id,omitempty"`
	CustomerName     string                  `json:"customer_name"`
	InvoiceReference string                  `json:"invoice_reference"`
	TotalAmount      decimal.Decimal         `json:"total_amount"`
	IssueDate        time.Time               `json:"issue_date"`
	DueDate          time.Time               `json:"due_date"`
	AccountCode      string                  `json:"account_code"` // e.g. 1-1200 / 1-1210 / 1-1220
	Notes            string                  `json:"notes,omitempty"`
}

type PayReceivableRequest struct {
	Amount              decimal.Decimal `json:"amount"`
	PaymentDate         time.Time       `json:"payment_date"`
	CashBankAccountCode string          `json:"cash_bank_account_code"` // e.g. 1-1111 (Kasir), 1-1121 (BCA)
	ReferenceNo         string          `json:"reference_no,omitempty"`
	Notes               string          `json:"notes,omitempty"`
}

