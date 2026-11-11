package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"clinic-api/internal/domain/audit"

	"github.com/google/uuid"
)

type Service struct {
	auditRepo audit.Repository
}

func NewService(auditRepo audit.Repository) *Service {
	return &Service{
		auditRepo: auditRepo,
	}
}

type LogCommand struct {
	BranchID      *uuid.UUID       `json:"branch_id,omitempty"`
	UserID        *uuid.UUID       `json:"user_id,omitempty"`
	Action        audit.ActionType `json:"action"`
	EntityType    string           `json:"entity_type"`
	EntityID      string           `json:"entity_id"`
	BeforePayload any              `json:"before_payload,omitempty"`
	AfterPayload  any              `json:"after_payload,omitempty"`
	IPAddress     string           `json:"ip_address,omitempty"`
	UserAgent     string           `json:"user_agent,omitempty"`
}

func (s *Service) Log(ctx context.Context, cmd LogCommand) error {
	var beforeBytes, afterBytes []byte
	var err error

	if cmd.BeforePayload != nil {
		beforeBytes, err = json.Marshal(cmd.BeforePayload)
		if err != nil {
			return fmt.Errorf("marshal before payload: %w", err)
		}
	}
	if cmd.AfterPayload != nil {
		afterBytes, err = json.Marshal(cmd.AfterPayload)
		if err != nil {
			return fmt.Errorf("marshal after payload: %w", err)
		}
	}

	entry := &audit.FinancialAuditLog{
		ID:            uuid.New(),
		BranchID:      cmd.BranchID,
		UserID:        cmd.UserID,
		Action:        cmd.Action,
		EntityType:    cmd.EntityType,
		EntityID:      cmd.EntityID,
		BeforePayload: beforeBytes,
		AfterPayload:  afterBytes,
		IPAddress:     cmd.IPAddress,
		UserAgent:     cmd.UserAgent,
		CreatedAt:     time.Now(),
	}

	return s.auditRepo.Log(ctx, entry)
}

func (s *Service) ListLogs(ctx context.Context, branchID *uuid.UUID, action *audit.ActionType, entityType string, limit, offset int) ([]audit.FinancialAuditLog, int64, error) {
	return s.auditRepo.ListLogs(ctx, branchID, action, entityType, limit, offset)
}

