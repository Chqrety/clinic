package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"clinic-api/internal/domain/audit"

	"github.com/google/uuid"
)

type PostgresAuditRepo struct {
	db *sql.DB
}

func NewPostgresAuditRepo(db *sql.DB) *PostgresAuditRepo {
	return &PostgresAuditRepo{db: db}
}

func (r *PostgresAuditRepo) Log(ctx context.Context, entry *audit.FinancialAuditLog) error {
	exec := GetExecutor(ctx, r.db)
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}

	query := `
		INSERT INTO financial_audit_logs (
			id, branch_id, user_id, action, entity_type, entity_id,
			before_payload, after_payload, ip_address, user_agent,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		entry.ID, entry.BranchID, entry.UserID, entry.Action, entry.EntityType, entry.EntityID,
		entry.BeforePayload, entry.AfterPayload, entry.IPAddress, entry.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("insert financial audit log: %w", err)
	}
	return nil
}

func (r *PostgresAuditRepo) ListLogs(ctx context.Context, branchID *uuid.UUID, action *audit.ActionType, entityType string, limit, offset int) ([]audit.FinancialAuditLog, int64, error) {
	exec := GetExecutor(ctx, r.db)

	if limit <= 0 {
		limit = 50
	}

	var actionStr *string
	if action != nil {
		s := string(*action)
		actionStr = &s
	}

	var count int64
	countQuery := `
		SELECT COUNT(*)
		FROM financial_audit_logs
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND ($2::varchar IS NULL OR action = $2)
		  AND ($3 = '' OR entity_type = $3)
	`
	_ = exec.QueryRowContext(ctx, countQuery, branchID, actionStr, entityType).Scan(&count)

	query := `
		SELECT id, branch_id, user_id, action, entity_type, entity_id,
		       before_payload, after_payload, ip_address, user_agent, created_at
		FROM financial_audit_logs
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND ($2::varchar IS NULL OR action = $2)
		  AND ($3 = '' OR entity_type = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := exec.QueryContext(ctx, query, branchID, actionStr, entityType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []audit.FinancialAuditLog
	for rows.Next() {
		var l audit.FinancialAuditLog
		err := rows.Scan(
			&l.ID, &l.BranchID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID,
			&l.BeforePayload, &l.AfterPayload, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan audit log: %w", err)
		}
		logs = append(logs, l)
	}

	return logs, count, nil
}

