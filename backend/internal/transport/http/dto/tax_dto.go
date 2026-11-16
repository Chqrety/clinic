package dto

import (
	"clinic-api/internal/domain/tax"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateWithholdingRequest struct {
	BranchID     uuid.UUID       `json:"branch_id"`
	TaxType      tax.TaxType     `json:"tax_type"`
	TaxpayerNIK  string          `json:"taxpayer_nik"`
	TaxpayerNPWP string          `json:"taxpayer_npwp"`
	TaxpayerName string          `json:"taxpayer_name"`
	ObjectCode   string          `json:"object_code"`
	GrossAmount  decimal.Decimal `json:"gross_amount"`
	DPPAmount    decimal.Decimal `json:"dpp_amount"`
	TaxRate      decimal.Decimal `json:"tax_rate"`
	TaxAmount    decimal.Decimal `json:"tax_amount"`
	PeriodYear   int             `json:"period_year"`
	PeriodMonth  int             `json:"period_month"`
	JournalID    *uuid.UUID      `json:"journal_id,omitempty"`
}

