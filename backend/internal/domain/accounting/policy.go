package accounting

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CostCenter struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"` // 'UNIT-KIDINEST', 'UNIT-SSC', 'SHARED-BACKOFFICE'
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PolicyDataType string

const (
	PolicyDataTypeString  PolicyDataType = "STRING"
	PolicyDataTypeBoolean PolicyDataType = "BOOLEAN"
	PolicyDataTypeNumber  PolicyDataType = "NUMBER"
	PolicyDataTypeJSON    PolicyDataType = "JSON"
)

type AccountingPolicy struct {
	ID            uuid.UUID      `json:"id"`
	BranchID      *uuid.UUID     `json:"branch_id,omitempty"`
	PolicyKey     string         `json:"policy_key"`
	PolicyValue   string         `json:"policy_value"`
	DataType      PolicyDataType `json:"data_type"`
	Description   string         `json:"description,omitempty"`
	EffectiveDate time.Time      `json:"effective_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type PolicyRepository interface {
	GetByKey(ctx context.Context, branchID *uuid.UUID, key string) (*AccountingPolicy, error)
	SetPolicy(ctx context.Context, policy *AccountingPolicy) error
	List(ctx context.Context, branchID *uuid.UUID) ([]AccountingPolicy, error)
	GetCostCenterByCode(ctx context.Context, code string) (*CostCenter, error)
	ListCostCenters(ctx context.Context) ([]CostCenter, error)
}

