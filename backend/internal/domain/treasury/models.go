package treasury

import (
	"context"
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type ShiftStatus string

const (
	ShiftStatusOpen      ShiftStatus = "OPEN"
	ShiftStatusClosed    ShiftStatus = "CLOSED"
	ShiftStatusDeposited ShiftStatus = "DEPOSITED"
)

type DifferenceType string

const (
	DifferenceNone     DifferenceType = "NONE"
	DifferenceOverage  DifferenceType = "OVERAGE"  // Cash surplus -> 4-9200
	DifferenceShortage DifferenceType = "SHORTAGE" // Cash shortage -> 6-1700
)

type ShiftDenomination struct {
	ID                 uuid.UUID       `json:"id"`
	ShiftID            uuid.UUID       `json:"shift_id"`
	DenominationValue  int             `json:"denomination_value"` // 100000, 50000, 20000, 10000, 5000, 2000, 1000, 500
	Quantity           int             `json:"quantity"`
	Subtotal           decimal.Decimal `json:"subtotal"`
}

type CashierShift struct {
	shared.BaseEntity
	shared.BranchScoped
	CostCenterID            uuid.UUID           `json:"cost_center_id"`
	CashierUserID           *uuid.UUID          `json:"cashier_user_id,omitempty"`
	CashierName             string              `json:"cashier_name"`
	CashAccountCode         string              `json:"cash_account_code"` // '1-1111' Kidinest or '1-1112' SSC
	StartTime               time.Time           `json:"start_time"`
	EndTime                 *time.Time          `json:"end_time,omitempty"`
	OpeningCash             decimal.Decimal     `json:"opening_cash"`
	SystemCashIn            decimal.Decimal     `json:"system_cash_in"`
	SystemCashOut           decimal.Decimal     `json:"system_cash_out"`
	ExpectedCash            decimal.Decimal     `json:"expected_cash"`
	ActualCash              decimal.Decimal     `json:"actual_cash"`
	DifferenceAmount        decimal.Decimal     `json:"difference_amount"`
	DifferenceType          DifferenceType      `json:"difference_type"`
	ReconciliationJournalID *uuid.UUID          `json:"reconciliation_journal_id,omitempty"`
	DepositJournalID        *uuid.UUID          `json:"deposit_journal_id,omitempty"`
	Status                  ShiftStatus         `json:"status"`
	Notes                   string              `json:"notes,omitempty"`
	Denominations           []ShiftDenomination `json:"denominations,omitempty"`
}

type ShiftFilter struct {
	BranchID        *uuid.UUID
	CostCenterID    *uuid.UUID
	CashAccountCode string
	Status          *ShiftStatus
	StartDate       *time.Time
	EndDate         *time.Time
	Search          string
	SortBy          string
	SortOrder       string
	Limit           int
	Offset          int
}

type Repository interface {
	OpenShift(ctx context.Context, shift *CashierShift) error
	GetActiveShift(ctx context.Context, branchID uuid.UUID, cashAccountCode string) (*CashierShift, error)
	GetShiftByID(ctx context.Context, id uuid.UUID) (*CashierShift, error)
	CloseShift(ctx context.Context, shift *CashierShift, denoms []ShiftDenomination) error
	RecordDeposit(ctx context.Context, shiftID uuid.UUID, depositJournalID uuid.UUID) error
	ListShifts(ctx context.Context, filter ShiftFilter) ([]CashierShift, int64, error)
}

