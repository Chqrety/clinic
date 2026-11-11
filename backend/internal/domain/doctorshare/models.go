package doctorshare

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type SettlementStatus string

const (
	StatusDraft    SettlementStatus = "DRAFT"
	StatusApproved SettlementStatus = "APPROVED"
	StatusPaid     SettlementStatus = "PAID"
)

type DoctorShareRule struct {
	ID                uuid.UUID       `json:"id"`
	BranchID          uuid.UUID       `json:"branch_id"`
	CostCenterID      uuid.UUID       `json:"cost_center_id"`
	DoctorID          *uuid.UUID      `json:"doctor_id,omitempty"`
	DoctorName        string          `json:"doctor_name"`
	ServiceCategory   string          `json:"service_category"` // 'CONSULTATION', 'VACCINE_INJECTION', 'GROWTH_PACKAGE', 'FACIAL_TREATMENT', 'LASER_PEELING'
	SharingPercentage decimal.Decimal `json:"sharing_percentage"` // e.g. 70.00
	FixedFee          decimal.Decimal `json:"fixed_fee"`
	IsActive          bool            `json:"is_active"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type DoctorFeeItem struct {
	ID                 uuid.UUID       `json:"id"`
	SettlementID       uuid.UUID       `json:"settlement_id"`
	InvoiceID          uuid.UUID       `json:"invoice_id"`
	ServiceDate        time.Time       `json:"service_date"`
	PatientName        string          `json:"patient_name"`
	ServiceDescription string          `json:"service_description"`
	PatientFee         decimal.Decimal `json:"patient_fee"`
	DoctorShareAmount  decimal.Decimal `json:"doctor_share_amount"`
	CreatedAt          time.Time       `json:"created_at"`
}

type DoctorFeeSettlement struct {
	shared.BaseEntity
	shared.BranchScoped
	DoctorID    *uuid.UUID        `json:"doctor_id,omitempty"`
	DoctorName  string            `json:"doctor_name"`
	DoctorNPWP  string            `json:"doctor_npwp,omitempty"`
	PeriodYear  int               `json:"period_year"`
	PeriodMonth int               `json:"period_month"`
	GrossFee    decimal.Decimal   `json:"gross_fee"`
	DPPAmount   decimal.Decimal   `json:"dpp_amount"` // 50% of gross
	TaxPPh21    decimal.Decimal   `json:"tax_pph21"`  // PMK 168/2023 progressive bracket
	NetFee      decimal.Decimal   `json:"net_fee"`
	Status      SettlementStatus  `json:"status"`
	JournalID   *uuid.UUID        `json:"journal_id,omitempty"`
	PaidAt      *time.Time        `json:"paid_at,omitempty"`
	Items       []DoctorFeeItem   `json:"items,omitempty"`
}

type Repository interface {
	CreateRule(ctx context.Context, rule *DoctorShareRule) error
	GetRule(ctx context.Context, branchID uuid.UUID, doctorName, serviceCategory string) (*DoctorShareRule, error)
	ListRules(ctx context.Context, branchID *uuid.UUID) ([]DoctorShareRule, error)
	CreateSettlement(ctx context.Context, s *DoctorFeeSettlement) error
	GetSettlementByID(ctx context.Context, id uuid.UUID) (*DoctorFeeSettlement, error)
	ListSettlements(ctx context.Context, branchID *uuid.UUID, year, month int) ([]DoctorFeeSettlement, error)
	UpdateSettlementStatus(ctx context.Context, id uuid.UUID, status SettlementStatus, paidAt *time.Time, journalID *uuid.UUID) error
}
