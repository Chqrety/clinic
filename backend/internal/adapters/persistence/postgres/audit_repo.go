package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (r *PostgresAuditRepo) ListLogs(ctx context.Context, filter audit.AuditFilter) ([]audit.FinancialAuditLog, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM financial_audit_logs WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.Action != nil {
		baseQuery += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, string(*filter.Action))
		argIdx++
	}
	if filter.EntityType != "" {
		baseQuery += fmt.Sprintf(" AND entity_type = $%d", argIdx)
		args = append(args, filter.EntityType)
		argIdx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *filter.StartDate)
		argIdx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *filter.EndDate)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(entity_id) LIKE $%d OR LOWER(COALESCE(ip_address, '')) LIKE $%d OR LOWER(COALESCE(user_agent, '')) LIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, searchParam, searchParam, searchParam)
		argIdx += 3
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count financial audit logs: %w", err)
	}

	allowedSorts := map[string]string{
		"created_at":  "created_at",
		"action":      "action",
		"entity_type": "entity_type",
		"entity_id":   "entity_id",
	}
	sortCol := "created_at"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, id %s", sortCol, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, user_id, action, entity_type, entity_id,
		       before_payload, after_payload, ip_address, user_agent, created_at
	` + baseQuery + orderBy

	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		selectQuery += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
		argIdx++
	}

	rows, err := exec.QueryContext(ctx, selectQuery, args...)
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

	return logs, totalCount, nil
}

