package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain/doctorshare"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresDoctorShareRepo struct {
	db *sql.DB
}

func NewPostgresDoctorShareRepo(db *sql.DB) *PostgresDoctorShareRepo {
	return &PostgresDoctorShareRepo{db: db}
}

func (r *PostgresDoctorShareRepo) CreateRule(ctx context.Context, rule *doctorshare.DoctorShareRule) error {
	exec := GetExecutor(ctx, r.db)
	if rule.ID == uuid.Nil {
		rule.ID = uuid.New()
	}

	query := `
		INSERT INTO doctor_share_rules (
			id, branch_id, cost_center_id, doctor_id, doctor_name,
			service_category, sharing_percentage, fixed_fee, is_active,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		rule.ID, rule.BranchID, rule.CostCenterID, rule.DoctorID, rule.DoctorName,
		rule.ServiceCategory, rule.SharingPercentage, rule.FixedFee, rule.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert doctor share rule: %w", err)
	}
	return nil
}

func (r *PostgresDoctorShareRepo) GetRule(ctx context.Context, branchID uuid.UUID, doctorName, serviceCategory string) (*doctorshare.DoctorShareRule, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, cost_center_id, doctor_id, doctor_name,
		       service_category, sharing_percentage, fixed_fee, is_active,
		       created_at, updated_at
		FROM doctor_share_rules
		WHERE branch_id = $1 AND is_active = TRUE
		  AND (doctor_name ILIKE $2 OR doctor_name = '*')
		  AND (service_category ILIKE $3 OR service_category = '*')
		ORDER BY CASE WHEN doctor_name != '*' THEN 0 ELSE 1 END,
		         CASE WHEN service_category != '*' THEN 0 ELSE 1 END
		LIMIT 1
	`
	row := exec.QueryRowContext(ctx, query, branchID, doctorName, serviceCategory)
	var rule doctorshare.DoctorShareRule
	var fixedFeeStr, sharingPctStr string
	err := row.Scan(
		&rule.ID, &rule.BranchID, &rule.CostCenterID, &rule.DoctorID, &rule.DoctorName,
		&rule.ServiceCategory, &sharingPctStr, &fixedFeeStr, &rule.IsActive,
		&rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor share rule not found")
		}
		return nil, fmt.Errorf("query doctor share rule: %w", err)
	}
	rule.SharingPercentage, _ = decimal.NewFromString(sharingPctStr)
	rule.FixedFee, _ = decimal.NewFromString(fixedFeeStr)
	return &rule, nil
}

func (r *PostgresDoctorShareRepo) ListRules(ctx context.Context, filter doctorshare.RuleFilter) ([]doctorshare.DoctorShareRule, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM doctor_share_rules WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if strings.TrimSpace(filter.Search) != "" {
		searchParam := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		baseQuery += fmt.Sprintf(" AND (LOWER(doctor_name) LIKE $%d OR LOWER(service_category) LIKE $%d)", argIdx, argIdx+1)
		args = append(args, searchParam, searchParam)
		argIdx += 2
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count doctor share rules: %w", err)
	}

	allowedSorts := map[string]string{
		"doctor_name":        "doctor_name",
		"service_category":   "service_category",
		"sharing_percentage": "sharing_percentage",
		"fixed_fee":          "fixed_fee",
		"is_active":          "is_active",
		"created_at":         "created_at",
	}
	sortCol := "doctor_name"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "ASC"
	if strings.ToUpper(filter.SortOrder) == "DESC" {
		sortDir = "DESC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, service_category %s, id %s", sortCol, sortDir, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, cost_center_id, doctor_id, doctor_name,
		       service_category, sharing_percentage, fixed_fee, is_active,
		       created_at, updated_at
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
		return nil, 0, fmt.Errorf("list doctor share rules: %w", err)
	}
	defer rows.Close()

	var rules []doctorshare.DoctorShareRule
	for rows.Next() {
		var rule doctorshare.DoctorShareRule
		var fixedFeeStr, sharingPctStr string
		err := rows.Scan(
			&rule.ID, &rule.BranchID, &rule.CostCenterID, &rule.DoctorID, &rule.DoctorName,
			&rule.ServiceCategory, &sharingPctStr, &fixedFeeStr, &rule.IsActive,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan doctor share rule: %w", err)
		}
		rule.SharingPercentage, _ = decimal.NewFromString(sharingPctStr)
		rule.FixedFee, _ = decimal.NewFromString(fixedFeeStr)
		rules = append(rules, rule)
	}
	return rules, totalCount, nil
}

func (r *PostgresDoctorShareRepo) CreateSettlement(ctx context.Context, s *doctorshare.DoctorFeeSettlement) error {
	exec := GetExecutor(ctx, r.db)
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	query := `
		INSERT INTO doctor_fee_settlements (
			id, branch_id, doctor_id, doctor_name, doctor_npwp,
			period_year, period_month, gross_fee, dpp_amount, tax_pph21,
			net_fee, status, journal_id, paid_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		s.ID, s.BranchID, s.DoctorID, s.DoctorName, s.DoctorNPWP,
		s.PeriodYear, s.PeriodMonth, s.GrossFee, s.DPPAmount, s.TaxPPh21,
		s.NetFee, s.Status, s.JournalID, s.PaidAt,
	)
	if err != nil {
		return fmt.Errorf("insert doctor fee settlement: %w", err)
	}

	for i := range s.Items {
		item := &s.Items[i]
		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.SettlementID = s.ID
		itemQuery := `
			INSERT INTO doctor_fee_items (
				id, settlement_id, invoice_id, service_date,
				patient_name, service_description, patient_fee, doctor_share_amount,
				created_at
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				CURRENT_TIMESTAMP
			)
		`
		_, err := exec.ExecContext(ctx, itemQuery,
			item.ID, item.SettlementID, item.InvoiceID, item.ServiceDate,
			item.PatientName, item.ServiceDescription, item.PatientFee, item.DoctorShareAmount,
		)
		if err != nil {
			return fmt.Errorf("insert doctor fee item: %w", err)
		}
	}

	return nil
}

func (r *PostgresDoctorShareRepo) GetSettlementByID(ctx context.Context, id uuid.UUID) (*doctorshare.DoctorFeeSettlement, error) {
	exec := GetExecutor(ctx, r.db)
	query := `
		SELECT id, branch_id, doctor_id, doctor_name, doctor_npwp,
		       period_year, period_month, gross_fee, dpp_amount, tax_pph21,
		       net_fee, status, journal_id, paid_at, created_at, updated_at
		FROM doctor_fee_settlements
		WHERE id = $1
	`
	row := exec.QueryRowContext(ctx, query, id)
	var s doctorshare.DoctorFeeSettlement
	var grossStr, dppStr, taxStr, netStr string
	err := row.Scan(
		&s.ID, &s.BranchID, &s.DoctorID, &s.DoctorName, &s.DoctorNPWP,
		&s.PeriodYear, &s.PeriodMonth, &grossStr, &dppStr, &taxStr,
		&netStr, &s.Status, &s.JournalID, &s.PaidAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor fee settlement not found")
		}
		return nil, fmt.Errorf("query doctor fee settlement: %w", err)
	}
	s.GrossFee, _ = decimal.NewFromString(grossStr)
	s.DPPAmount, _ = decimal.NewFromString(dppStr)
	s.TaxPPh21, _ = decimal.NewFromString(taxStr)
	s.NetFee, _ = decimal.NewFromString(netStr)

	// Fetch items
	itemsQuery := `
		SELECT id, settlement_id, invoice_id, service_date,
		       patient_name, service_description, patient_fee, doctor_share_amount,
		       created_at
		FROM doctor_fee_items
		WHERE settlement_id = $1
		ORDER BY service_date ASC, created_at ASC
	`
	rows, err := exec.QueryContext(ctx, itemsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it doctorshare.DoctorFeeItem
			var patFeeStr, docFeeStr string
			if err := rows.Scan(
				&it.ID, &it.SettlementID, &it.InvoiceID, &it.ServiceDate,
				&it.PatientName, &it.ServiceDescription, &patFeeStr, &docFeeStr,
				&it.CreatedAt,
			); err == nil {
				it.PatientFee, _ = decimal.NewFromString(patFeeStr)
				it.DoctorShareAmount, _ = decimal.NewFromString(docFeeStr)
				s.Items = append(s.Items, it)
			}
		}
	}

	return &s, nil
}

func (r *PostgresDoctorShareRepo) ListSettlements(ctx context.Context, filter doctorshare.SettlementFilter) ([]doctorshare.DoctorFeeSettlement, int64, error) {
	exec := GetExecutor(ctx, r.db)
	baseQuery := `FROM doctor_fee_settlements WHERE 1=1`
	var args []any
	argIdx := 1

	if filter.BranchID != nil {
		baseQuery += fmt.Sprintf(" AND branch_id = $%d", argIdx)
		args = append(args, *filter.BranchID)
		argIdx++
	}
	if filter.DoctorID != nil {
		baseQuery += fmt.Sprintf(" AND doctor_id = $%d", argIdx)
		args = append(args, *filter.DoctorID)
		argIdx++
	}
	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.Year > 0 {
		baseQuery += fmt.Sprintf(" AND period_year = $%d", argIdx)
		args = append(args, filter.Year)
		argIdx++
	}
	if filter.Month > 0 {
		baseQuery += fmt.Sprintf(" AND period_month = $%d", argIdx)
		args = append(args, filter.Month)
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
		baseQuery += fmt.Sprintf(" AND (LOWER(doctor_name) LIKE $%d OR LOWER(COALESCE(doctor_npwp, '')) LIKE $%d)", argIdx, argIdx+1)
		args = append(args, searchParam, searchParam)
		argIdx += 2
	}

	var totalCount int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := exec.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf("count settlements: %w", err)
	}

	allowedSorts := map[string]string{
		"doctor_name":  "doctor_name",
		"period_year":  "period_year",
		"period_month": "period_month",
		"gross_fee":    "gross_fee",
		"net_fee":      "net_fee",
		"status":       "status",
		"created_at":   "created_at",
	}
	sortCol := "period_year"
	if col, ok := allowedSorts[strings.ToLower(filter.SortBy)]; ok {
		sortCol = col
	}
	sortDir := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortDir = "ASC"
	}
	orderBy := fmt.Sprintf(" ORDER BY %s %s, period_month %s, doctor_name ASC, id %s", sortCol, sortDir, sortDir, sortDir)

	selectQuery := `
		SELECT id, branch_id, doctor_id, doctor_name, doctor_npwp,
		       period_year, period_month, gross_fee, dpp_amount, tax_pph21,
		       net_fee, status, journal_id, paid_at, created_at, updated_at
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
		return nil, 0, fmt.Errorf("list settlements: %w", err)
	}
	defer rows.Close()

	var settlements []doctorshare.DoctorFeeSettlement
	for rows.Next() {
		var s doctorshare.DoctorFeeSettlement
		var grossStr, dppStr, taxStr, netStr string
		err := rows.Scan(
			&s.ID, &s.BranchID, &s.DoctorID, &s.DoctorName, &s.DoctorNPWP,
			&s.PeriodYear, &s.PeriodMonth, &grossStr, &dppStr, &taxStr,
			&netStr, &s.Status, &s.JournalID, &s.PaidAt, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan settlement: %w", err)
		}
		s.GrossFee, _ = decimal.NewFromString(grossStr)
		s.DPPAmount, _ = decimal.NewFromString(dppStr)
		s.TaxPPh21, _ = decimal.NewFromString(taxStr)
		s.NetFee, _ = decimal.NewFromString(netStr)
		settlements = append(settlements, s)
	}
	return settlements, totalCount, nil
}

func (r *PostgresDoctorShareRepo) UpdateSettlementStatus(ctx context.Context, id uuid.UUID, status doctorshare.SettlementStatus, paidAt *time.Time, journalID *uuid.UUID) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE doctor_fee_settlements
		SET status = $1, paid_at = $2, journal_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	_, err := exec.ExecContext(ctx, query, status, paidAt, journalID, id)
	if err != nil {
		return fmt.Errorf("update settlement status: %w", err)
	}
	return nil
}

