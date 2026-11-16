package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

func TestStandardPagination_ParsingAndMetadata(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/journals?page=2&limit=15&search=Kidinest&sort_by=transaction_date&sort_order=asc&status=POSTED", nil)
	cfg := pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"transaction_date": "transaction_date",
			"created_at":       "created_at",
			"journal_number":   "journal_number",
		},
		DefaultSortBy:    "transaction_date",
		DefaultSortOrder: "DESC",
	}

	q, err := pagination.Parse(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error parsing pagination: %v", err)
	}

	if q.Page != 2 {
		t.Errorf("expected page 2, got %d", q.Page)
	}
	if q.Limit != 15 {
		t.Errorf("expected limit 15, got %d", q.Limit)
	}
	if q.Offset != 15 {
		t.Errorf("expected offset 15, got %d", q.Offset)
	}
	if q.Search != "Kidinest" {
		t.Errorf("expected search 'Kidinest', got %s", q.Search)
	}
	if q.SortBy != "transaction_date" {
		t.Errorf("expected sort_by 'transaction_date', got %s", q.SortBy)
	}
	if q.SortOrder != "ASC" {
		t.Errorf("expected sort_order 'ASC', got %s", q.SortOrder)
	}

	// Verify metadata calculation
	totalItems := int64(45)
	meta := q.ToMeta(totalItems)
	if meta.CurrentPage != 2 {
		t.Errorf("expected current_page 2, got %d", meta.CurrentPage)
	}
	if meta.PerPage != 15 {
		t.Errorf("expected per_page 15, got %d", meta.PerPage)
	}
	if meta.TotalItems != 45 {
		t.Errorf("expected total_items 45, got %d", meta.TotalItems)
	}
	if meta.TotalPages != 3 {
		t.Errorf("expected total_pages 3, got %d", meta.TotalPages)
	}
}

func TestStandardPagination_SecurityAndInjectionPrevention(t *testing.T) {
	cfg := pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"due_date":   "due_date",
			"created_at": "created_at",
		},
		DefaultSortBy:    "due_date",
		DefaultSortOrder: "ASC",
	}

	maliciousURLs := []struct {
		name string
		url  string
	}{
		{"SQL injection in sort_by", "/api/v1/invoices?sort_by=due_date;DROP+TABLE+invoices--"},
		{"SQL injection with union", "/api/v1/invoices?sort_by=due_date+UNION+SELECT+1,2,3--"},
		{"Non-whitelisted column", "/api/v1/invoices?sort_by=password_hash"},
		{"Invalid sort direction", "/api/v1/invoices?sort_by=due_date&sort_order=INVALID"},
		{"Negative page", "/api/v1/invoices?page=-5"},
		{"Invalid UUID for branch", "/api/v1/invoices?branch_id=not-a-uuid"},
		{"Start date after end date", "/api/v1/invoices?start_date=2026-12-31&end_date=2026-01-01"},
	}

	for _, tc := range maliciousURLs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			_, err := pagination.Parse(req, cfg)
			if err == nil {
				t.Fatalf("expected error for URL %s, but got nil", tc.url)
			}
		})
	}
}

func TestStandardPagination_DeterministicTieBreaker(t *testing.T) {
	cfg := pagination.Config{
		AllowedSortColumns: map[string]string{
			"issue_date": "r.issue_date",
		},
		DefaultSortBy:    "r.issue_date",
		DefaultSortOrder: "ASC",
		TieBreakerColumn: "r.id",
	}

	req := httptest.NewRequest("GET", "/api/v1/receivables?sort_by=issue_date&sort_order=desc", nil)
	q, err := pagination.Parse(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	orderBy := q.BuildOrderBy("r.")
	expected := "ORDER BY r.issue_date DESC, r.id DESC"
	if orderBy != expected {
		t.Errorf("expected '%s', got '%s'", expected, orderBy)
	}
}

func TestStandardPagination_AuthorizationValidation(t *testing.T) {
	authorizedUnitID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	otherUnitID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	session := &middleware.UserSession{
		UserID:               uuid.New(),
		AllowedCostCenterIDs: []uuid.UUID{authorizedUnitID},
		IsSuperAdmin:         false,
	}

	ctx := middleware.SetUserSession(t.Context(), session)

	// 1. Authorized cost center request should succeed (return true)
	if !middleware.ValidateUnitAuthorization(ctx, authorizedUnitID) {
		t.Fatalf("expected authorized unit to succeed, got false")
	}

	// 2. Unauthorized cost center request should be rejected (return false)
	if middleware.ValidateUnitAuthorization(ctx, otherUnitID) {
		t.Fatalf("expected unauthorized unit to fail, got true")
	}

	// 3. SuperAdmin bypass should succeed for any unit
	superAdminSession := &middleware.UserSession{
		UserID:       uuid.New(),
		IsSuperAdmin: true,
	}
	superCtx := middleware.SetUserSession(t.Context(), superAdminSession)
	if !middleware.ValidateUnitAuthorization(superCtx, otherUnitID) {
		t.Fatalf("expected super admin to access any unit, got false")
	}
}

func TestStandardPagination_JSONResponseFormat(t *testing.T) {
	rec := httptest.NewRecorder()
	items := []string{"account1", "account2", "account3"}
	meta := response.PaginationMeta{
		CurrentPage: 1,
		PerPage:     20,
		TotalItems:  3,
		TotalPages:  1,
	}

	response.OK(rec, "Items retrieved successfully", items, meta)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res response.SuccessResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}

	metaBytes, _ := json.Marshal(res.Meta)
	var decodedMeta response.PaginationMeta
	_ = json.Unmarshal(metaBytes, &decodedMeta)

	if decodedMeta.CurrentPage != 1 || decodedMeta.TotalItems != 3 || decodedMeta.TotalPages != 1 {
		t.Errorf("unexpected meta values: %+v", decodedMeta)
	}
}

func TestStandardPagination_DateFilterConsistency(t *testing.T) {
	cfg := pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"created_at": "created_at",
		},
		DefaultSortBy: "created_at",
	}

	req := httptest.NewRequest("GET", "/api/v1/invoices?start_date=2026-01-01&end_date=2026-01-31", nil)
	q, err := pagination.Parse(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.StartDate == nil || q.EndDate == nil {
		t.Fatalf("expected parsed start and end dates, got nil")
	}

	expectedStart, _ := time.Parse("2006-01-02", "2026-01-01")
	if !q.StartDate.Equal(expectedStart) {
		t.Errorf("expected start date %v, got %v", expectedStart, *q.StartDate)
	}

	// End date should be set to end of day 23:59:59
	if q.EndDate.Hour() != 23 || q.EndDate.Minute() != 59 || q.EndDate.Second() != 59 {
		t.Errorf("expected end date at 23:59:59, got %v", *q.EndDate)
	}
}
