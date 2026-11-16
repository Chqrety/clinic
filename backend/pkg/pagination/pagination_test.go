package pagination

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestParsePagination_Defaults(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/journals", nil)
	cfg := Config{
		DefaultSortBy: "created_at",
		AllowedSortColumns: map[string]string{
			"created_at":       "created_at",
			"transaction_date": "transaction_date",
		},
	}

	q, err := Parse(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Page != 1 {
		t.Errorf("expected page 1, got %d", q.Page)
	}
	if q.Limit != 20 {
		t.Errorf("expected limit 20, got %d", q.Limit)
	}
	if q.Offset != 0 {
		t.Errorf("expected offset 0, got %d", q.Offset)
	}
	if q.SortBy != "created_at" {
		t.Errorf("expected sort_by 'created_at', got %s", q.SortBy)
	}
	if q.SortOrder != "DESC" {
		t.Errorf("expected sort_order 'DESC', got %s", q.SortOrder)
	}
}

func TestParsePagination_CustomAndClamping(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/journals?page=3&limit=500&sort_by=transaction_date&sort_order=asc", nil)
	cfg := Config{
		DefaultSortBy: "created_at",
		AllowedSortColumns: map[string]string{
			"created_at":       "created_at",
			"transaction_date": "transaction_date",
		},
		MaxLimit: 100,
	}

	q, err := Parse(req, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Page != 3 {
		t.Errorf("expected page 3, got %d", q.Page)
	}
	if q.Limit != 100 { // Clamped to MaxLimit 100
		t.Errorf("expected limit clamped to 100, got %d", q.Limit)
	}
	if q.Offset != 200 { // (3 - 1) * 100
		t.Errorf("expected offset 200, got %d", q.Offset)
	}
	if q.SortBy != "transaction_date" {
		t.Errorf("expected sort_by 'transaction_date', got %s", q.SortBy)
	}
	if q.SortOrder != "ASC" {
		t.Errorf("expected sort_order 'ASC', got %s", q.SortOrder)
	}
}

func TestParsePagination_ValidationErrors(t *testing.T) {
	cfg := Config{
		DefaultSortBy: "created_at",
		AllowedSortColumns: map[string]string{
			"created_at": "created_at",
		},
	}

	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{
			name:    "invalid page negative",
			url:     "/api/v1/journals?page=-1",
			wantErr: "invalid 'page' parameter",
		},
		{
			name:    "invalid page string",
			url:     "/api/v1/journals?page=abc",
			wantErr: "invalid 'page' parameter",
		},
		{
			name:    "invalid limit zero",
			url:     "/api/v1/journals?limit=0",
			wantErr: "invalid 'limit' parameter",
		},
		{
			name:    "disallowed sort_by column injection attempt",
			url:     "/api/v1/journals?sort_by=id;DROP+TABLE+users;--",
			wantErr: "invalid 'sort_by' parameter",
		},
		{
			name:    "non whitelisted sort_by column",
			url:     "/api/v1/journals?sort_by=secret_password_hash",
			wantErr: "invalid 'sort_by' parameter",
		},
		{
			name:    "invalid sort_order",
			url:     "/api/v1/journals?sort_order=SIDEWAYS",
			wantErr: "invalid 'sort_order' parameter",
		},
		{
			name:    "start date after end date",
			url:     "/api/v1/journals?start_date=2026-12-01&end_date=2026-01-01",
			wantErr: "'start_date' cannot be after 'end_date'",
		},
		{
			name:    "invalid business unit UUID",
			url:     "/api/v1/journals?business_unit_id=not-a-uuid",
			wantErr: "invalid 'business_unit_id' format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			_, err := Parse(req, cfg)
			if err == nil {
				t.Fatalf("expected error containing '%s', got nil", tt.wantErr)
			}
		})
	}
}

func TestBuildOrderBy_TieBreaker(t *testing.T) {
	q := &Query{
		SortBy:     "transaction_date",
		SortOrder:  "DESC",
		TieBreaker: "id",
	}

	clause := q.BuildOrderBy("j")
	expected := "ORDER BY j.transaction_date DESC, j.id DESC"
	if clause != expected {
		t.Errorf("expected '%s', got '%s'", expected, clause)
	}

	// Tie-breaker when same column
	qSame := &Query{
		SortBy:     "id",
		SortOrder:  "ASC",
		TieBreaker: "id",
	}
	clauseSame := qSame.BuildOrderBy("j")
	expectedSame := "ORDER BY j.id ASC"
	if clauseSame != expectedSame {
		t.Errorf("expected '%s', got '%s'", expectedSame, clauseSame)
	}
}

func TestCursorEncodeDecode(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Nanosecond)
	id := uuid.New()

	cursor := EncodeCursor(now, id)
	if cursor == "" {
		t.Fatal("expected non-empty cursor string")
	}

	decodedTime, decodedID, err := DecodeCursor(cursor)
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if decodedTime.UnixNano() != now.UnixNano() {
		t.Errorf("expected timestamp %v, got %v", now, decodedTime)
	}
	if decodedID != id {
		t.Errorf("expected ID %v, got %v", id, decodedID)
	}
}

func TestToMeta(t *testing.T) {
	q := &Query{Page: 2, Limit: 20}
	meta := q.ToMeta(95)

	if meta.CurrentPage != 2 {
		t.Errorf("expected current page 2, got %d", meta.CurrentPage)
	}
	if meta.PerPage != 20 {
		t.Errorf("expected per page 20, got %d", meta.PerPage)
	}
	if meta.TotalItems != 95 {
		t.Errorf("expected total items 95, got %d", meta.TotalItems)
	}
	if meta.TotalPages != 5 {
		t.Errorf("expected total pages 5, got %d", meta.TotalPages)
	}
}
