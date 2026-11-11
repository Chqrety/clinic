package dto

import (
	"time"

	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateDoctorShareRuleRequest struct {
	BranchID          uuid.UUID       `json:"branch_id"`
	CostCenterID      uuid.UUID       `json:"cost_center_id"`
	DoctorID          *uuid.UUID      `json:"doctor_id,omitempty"`
	DoctorName        string          `json:"doctor_name"`
	ServiceCategory   string          `json:"service_category"` // 'CONSULTATION', 'VACCINE_INJECTION', 'GROWTH_PACKAGE', 'FACIAL_TREATMENT', 'LASER_PEELING'
	SharingPercentage decimal.Decimal `json:"sharing_percentage"`
	FixedFee          decimal.Decimal `json:"fixed_fee"`
	IsActive          bool            `json:"is_active"`
}

type FeeItemRequest struct {
	InvoiceID          uuid.UUID       `json:"invoice_id"`
	ServiceDate        time.Time       `json:"service_date"`
	PatientName        string          `json:"patient_name"`
	ServiceDescription string          `json:"service_description"`
	PatientFee         decimal.Decimal `json:"patient_fee"`
	DoctorShareAmount  decimal.Decimal `json:"doctor_share_amount"`
}

type GenerateDoctorSettlementRequest struct {
	BranchID    uuid.UUID        `json:"branch_id"`
	DoctorID    *uuid.UUID       `json:"doctor_id,omitempty"`
	DoctorName  string           `json:"doctor_name"`
	DoctorNPWP  string           `json:"doctor_npwp,omitempty"`
	TaxpayerNIK string           `json:"taxpayer_nik,omitempty"`
	PeriodYear  int              `json:"period_year"`
	PeriodMonth int              `json:"period_month"`
	Items       []FeeItemRequest `json:"items"`
}

type PayDoctorSettlementRequest struct {
	BankAccCode string     `json:"bank_account_code"` // e.g. '1-1121' Bank BCA
	UnitCode    string     `json:"unit_code"`         // 'KIDINEST' or 'SSC'
	PaidBy      *uuid.UUID `json:"paid_by,omitempty"`
}

