package dto

import (
	"time"

	"clinic-api/internal/domain/payable"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreatePayableRequest struct {
	VendorType       payable.VendorType `json:"vendor_type"`
	VendorID         *uuid.UUID         `json:"vendor_id,omitempty"`
	VendorName       string             `json:"vendor_name"`
	InvoiceReference string             `json:"invoice_reference"`
	TotalAmount      decimal.Decimal    `json:"total_amount"`
	IssueDate        time.Time          `json:"issue_date"`
	DueDate          time.Time          `json:"due_date"`
	AccountCode      string             `json:"account_code"` // e.g. 2-1110 (PBF), 2-1120 (Skincare), 2-1210 (Doctor)
	Notes            string             `json:"notes,omitempty"`
}

type PayPayableRequest struct {
	Amount              decimal.Decimal `json:"amount"`
	PaymentDate         time.Time       `json:"payment_date"`
	CashBankAccountCode string          `json:"cash_bank_account_code"` // e.g. 1-1121 (BCA), 1-1113 (Petty Cash)
	ReferenceNo         string          `json:"reference_no,omitempty"`
	Notes               string          `json:"notes,omitempty"`
}

