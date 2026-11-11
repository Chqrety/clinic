package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	ActionJournalPost    ActionType = "JOURNAL_POST"
	ActionJournalReverse ActionType = "JOURNAL_REVERSE"
	ActionPeriodClose    ActionType = "PERIOD_CLOSE"
	ActionPeriodReopen   ActionType = "PERIOD_REOPEN"
	ActionYearEndClose   ActionType = "YEAR_END_CLOSE"
	ActionPolicyChange   ActionType = "POLICY_CHANGE"
	ActionShiftClose     ActionType = "SHIFT_CLOSE"
	ActionDoctorPayment  ActionType = "DOCTOR_PAYMENT"
	ActionAssetRegister  ActionType = "ASSET_REGISTER"
)

type FinancialAuditLog struct {
	ID            uuid.UUID  `json:"id"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	Action        ActionType `json:"action"`
	EntityType    string     `json:"entity_type"`
	EntityID      string     `json:"entity_id"`
	BeforePayload []byte     `json:"before_payload,omitempty"`
	AfterPayload  []byte     `json:"after_payload,omitempty"`
	IPAddress     string     `json:"ip_address,omitempty"`
	UserAgent     string     `json:"user_agent,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Repository interface {
	Log(ctx context.Context, entry *FinancialAuditLog) error
	ListLogs(ctx context.Context, branchID *uuid.UUID, action *ActionType, entityType string, limit, offset int) ([]FinancialAuditLog, int64, error)
}

