package prepaid

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PrepaidStatus string

const (
	StatusActive          PrepaidStatus = "ACTIVE"
	StatusFullyAmortized  PrepaidStatus = "FULLY_AMORTIZED"
	StatusTerminated      PrepaidStatus = "TERMINATED"
)

type PrepaidCategory string

const (
	CategoryRent            PrepaidCategory = "RENT"             // Sewa Gedung Klinik
	CategoryInsurance       PrepaidCategory = "INSURANCE"        // Asuransi Malpraktik / Kebakaran
	CategorySoftwareLicense PrepaidCategory = "SOFTWARE_LICENSE" // Lisensi Cloud / SIMRS
	CategoryMaintenance     PrepaidCategory = "MAINTENANCE"      // Pemeliharaan Mesin Medis / Laser
)

type PrepaidExpense struct {
	shared.BaseEntity
	shared.BranchScoped
	CostCenterID       uuid.UUID       `json:"cost_center_id"`
	Name               string          `json:"name"`
	Category           PrepaidCategory `json:"category"`
	TotalAmount        decimal.Decimal `json:"total_amount"`
	PrepaidAccountCode string          `json:"prepaid_account_code"` // '1-1400'
	ExpenseAccountCode string          `json:"expense_account_code"` // '6-1200'
	StartDate          time.Time       `json:"start_date"`
	EndDate            time.Time       `json:"end_date"`
	PeriodMonths       int             `json:"period_months"`
	MonthlyAmount      decimal.Decimal `json:"monthly_amount"`
	RemainingAmount    decimal.Decimal `json:"remaining_amount"`
	Status             PrepaidStatus   `json:"status"`
	JournalID          *uuid.UUID      `json:"journal_id,omitempty"`
}

type PrepaidAmortization struct {
	ID               uuid.UUID       `json:"id"`
	PrepaidExpenseID uuid.UUID       `json:"prepaid_expense_id"`
	PeriodYear       int             `json:"period_year"`
	PeriodMonth      int             `json:"period_month"`
	Amount           decimal.Decimal `json:"amount"`
	JournalID        *uuid.UUID      `json:"journal_id,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

type Repository interface {
	CreatePrepaid(ctx context.Context, p *PrepaidExpense) error
	GetPrepaidByID(ctx context.Context, id uuid.UUID) (*PrepaidExpense, error)
	ListActivePrepaid(ctx context.Context, branchID *uuid.UUID) ([]PrepaidExpense, error)
	RecordAmortization(ctx context.Context, a *PrepaidAmortization, newRemaining decimal.Decimal, newStatus PrepaidStatus) error
	GetAmortizationHistory(ctx context.Context, prepaidID uuid.UUID) ([]PrepaidAmortization, error)
}

