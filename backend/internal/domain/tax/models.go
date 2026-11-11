package tax

import (
	"context"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type TaxType string

const (
	TaxTypePPh21Doctor   TaxType = "PPH21_DOCTOR"
	TaxTypePPh21Employee TaxType = "PPH21_EMPLOYEE"
	TaxTypePPh23         TaxType = "PPH23"
	TaxTypePPh4_2        TaxType = "PPH4_2"
	TaxTypePPNOutput     TaxType = "PPN_OUTPUT"
)

type TaxWithholding struct {
	shared.BaseEntity
	shared.BranchScoped
	TaxType       TaxType         `json:"tax_type"`
	BupotNumber   string          `json:"bupot_number"`
	TaxpayerNIK   string          `json:"taxpayer_nik"`
	TaxpayerNPWP  string          `json:"taxpayer_npwp"`
	TaxpayerName  string          `json:"taxpayer_name"`
	ObjectCode    string          `json:"object_code"` // e.g. "21-100-07" for Doctor
	GrossAmount   decimal.Decimal `json:"gross_amount"`
	DPPAmount     decimal.Decimal `json:"dpp_amount"`
	TaxRate       decimal.Decimal `json:"tax_rate"` // e.g. 5.00
	TaxAmount     decimal.Decimal `json:"tax_amount"`
	PeriodYear    int             `json:"period_year"`
	PeriodMonth   int             `json:"period_month"`
	IsReported    bool            `json:"is_reported"`
	JournalID     *uuid.UUID      `json:"journal_id,omitempty"`
}

// CalculateDoctorPPh21 applies PMK 168/2023 & UU HPP Pasal 17 progressive brackets:
// DPP = 50% * Gross
// Bracket 1: 0 - 60M @ 5%
// Bracket 2: 60M - 250M @ 15%
// Bracket 3: 250M - 500M @ 25%
// Bracket 4: 500M - 5B @ 30%
// Bracket 5: > 5B @ 35%
func CalculateDoctorPPh21(grossFee decimal.Decimal, priorYearDPP decimal.Decimal) (dpp decimal.Decimal, tax decimal.Decimal) {
	// DPP = 50% * Gross
	dpp = grossFee.Mul(decimal.FromFloat(0.50))
	if !dpp.IsPositive() {
		return decimal.Zero, decimal.Zero
	}

	b1 := decimal.FromInt(60000000)
	b2 := decimal.FromInt(250000000)
	b3 := decimal.FromInt(500000000)
	b4 := decimal.FromInt(5000000000)

	startDPP := priorYearDPP
	endDPP := priorYearDPP.Add(dpp)

	tax = decimal.Zero

	computeSlice := func(tierStart, tierEnd decimal.Decimal, rate float64) {
		if endDPP.GreaterThan(tierStart) && startDPP.LessThan(tierEnd) {
			s := startDPP
			if s.LessThan(tierStart) {
				s = tierStart
			}
			e := endDPP
			if e.GreaterThan(tierEnd) {
				e = tierEnd
			}
			sliceAmount := e.Sub(s)
			if sliceAmount.IsPositive() {
				tax = tax.Add(sliceAmount.Mul(decimal.FromFloat(rate)))
			}
		}
	}

	computeSlice(decimal.Zero, b1, 0.05)
	computeSlice(b1, b2, 0.15)
	computeSlice(b2, b3, 0.25)
	computeSlice(b3, b4, 0.30)
	computeSlice(b4, decimal.FromInt(999999999999), 0.35)

	return dpp, tax
}

type Repository interface {
	CreateWithholding(ctx context.Context, tw *TaxWithholding) error
	GetWithholdingByID(ctx context.Context, id uuid.UUID) (*TaxWithholding, error)
	GetCumulativeDPP(ctx context.Context, branchID uuid.UUID, taxpayerNIK string, year int) (decimal.Decimal, error)
	ListWithholdings(ctx context.Context, branchID *uuid.UUID, year, month int, taxType *TaxType) ([]TaxWithholding, error)
	MarkAsReported(ctx context.Context, id uuid.UUID) error
}
