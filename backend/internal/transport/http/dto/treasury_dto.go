package dto

import (
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type OpenShiftRequest struct {
	BranchID        uuid.UUID       `json:"branch_id"`
	CostCenterID    uuid.UUID       `json:"cost_center_id"`
	CashierUserID   *uuid.UUID      `json:"cashier_user_id,omitempty"`
	CashierName     string          `json:"cashier_name"`
	CashAccountCode string          `json:"cash_account_code"` // '1-1111' Kidinest, '1-1112' SSC
	OpeningCash     decimal.Decimal `json:"opening_cash"`
	Notes           string          `json:"notes,omitempty"`
}

type DenominationItemRequest struct {
	DenominationValue int `json:"denomination_value"`
	Quantity          int `json:"quantity"`
}

type CloseShiftRequest struct {
	SystemCashIn  decimal.Decimal           `json:"system_cash_in"`
	SystemCashOut decimal.Decimal           `json:"system_cash_out"`
	Denominations []DenominationItemRequest `json:"denominations"`
	Notes         string                    `json:"notes,omitempty"`
}

type DepositCashRequest struct {
	DepositAmount  decimal.Decimal `json:"deposit_amount"`
	BankAccCode    string          `json:"bank_account_code"` // Default '1-1121' Bank BCA
	DepositedBy    *uuid.UUID      `json:"deposited_by,omitempty"`
	DepositSlipRef string          `json:"deposit_slip_ref,omitempty"`
}

