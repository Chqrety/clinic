package pagination

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

const (
	DefaultPage      = 1
	DefaultLimit     = 20
	MaxLimit         = 100
	DefaultSortOrder = "DESC"
)

// Config defines rules for parsing and validating pagination & sorting query parameters.
type Config struct {
	DefaultLimit       int
	MaxLimit           int
	DefaultSortBy      string
	DefaultSortOrder   string
	AllowedSortColumns map[string]string // Maps client query param -> safe DB column expression
	TieBreakerColumn   string            // Deterministic tie-breaker column, e.g. "id"
}

// Query contains validated, sanitized parameters for pagination, filtering, searching, and sorting.
type Query struct {
	Page           int
	Limit          int
	Offset         int
	Search         string
	SortBy         string // Verified safe DB column
	SortOrder      string // "ASC" or "DESC"
	TieBreaker     string // Safe tie-breaker column
	StartDate      *time.Time
	EndDate        *time.Time
	BusinessUnitID *uuid.UUID
	BranchID       *uuid.UUID
	Status         string
	Cursor         string
}

// Parse extracts, sanitizes, and validates standard pagination and filter parameters from HTTP request.
func Parse(r *http.Request, cfg Config) (*Query, error) {
	if cfg.DefaultLimit <= 0 {
		cfg.DefaultLimit = DefaultLimit
	}
	if cfg.MaxLimit <= 0 {
		cfg.MaxLimit = MaxLimit
	}
	if cfg.DefaultSortOrder == "" {
		cfg.DefaultSortOrder = DefaultSortOrder
	}
	if cfg.TieBreakerColumn == "" {
		cfg.TieBreakerColumn = "id"
	}

	q := r.URL.Query()

	// 1. Page validation (minimum 1, default 1)
	page := DefaultPage
	if pageStr := q.Get("page"); pageStr != "" {
		val, err := strconv.Atoi(pageStr)
		if err != nil || val < 1 {
			return nil, fmt.Errorf("invalid 'page' parameter: must be integer >= 1")
		}
		page = val
	}

	// 2. Limit validation (default 20, max 100)
	limit := cfg.DefaultLimit
	if limitStr := q.Get("limit"); limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil || val < 1 {
			return nil, fmt.Errorf("invalid 'limit' parameter: must be integer >= 1")
		}
		if val > cfg.MaxLimit {
			limit = cfg.MaxLimit
		} else {
			limit = val
		}
	}

	offset := (page - 1) * limit

	// 3. Search query (trimmed)
	search := strings.TrimSpace(q.Get("search"))

	// 4. Sort column whitelist validation
	sortBy := cfg.DefaultSortBy
	clientSort := strings.TrimSpace(q.Get("sort_by"))
	if clientSort == "" && strings.Contains(r.URL.RawQuery, "sort_by=") {
		return nil, fmt.Errorf("invalid 'sort_by' parameter: malformed or disallowed characters in sort_by")
	}
	if clientSort != "" {
		if dbCol, ok := cfg.AllowedSortColumns[clientSort]; ok {
			sortBy = dbCol
		} else {
			return nil, fmt.Errorf("invalid 'sort_by' parameter '%s': allowed values are [%s]",
				clientSort, strings.Join(getMapKeys(cfg.AllowedSortColumns), ", "))
		}
	}

	// 5. Sort order validation ("ASC" or "DESC")
	sortOrder := strings.ToUpper(strings.TrimSpace(q.Get("sort_order")))
	if sortOrder == "" {
		sortOrder = strings.ToUpper(cfg.DefaultSortOrder)
	}
	if sortOrder != "ASC" && sortOrder != "DESC" {
		return nil, fmt.Errorf("invalid 'sort_order' parameter '%s': must be 'ASC' or 'DESC'", sortOrder)
	}

	// 6. Date range parsing (start_date, end_date)
	var startDate, endDate *time.Time
	if s := q.Get("start_date"); s != "" {
		t, err := parseDate(s, false)
		if err != nil {
			return nil, fmt.Errorf("invalid 'start_date' format: expected YYYY-MM-DD or RFC3339")
		}
		startDate = &t
	}
	if e := q.Get("end_date"); e != "" {
		t, err := parseDate(e, true)
		if err != nil {
			return nil, fmt.Errorf("invalid 'end_date' format: expected YYYY-MM-DD or RFC3339")
		}
		endDate = &t
	}
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, fmt.Errorf("'start_date' cannot be after 'end_date'")
	}

	// 7. Business Unit / Cost Center ID validation
	var businessUnitID *uuid.UUID
	unitParam := q.Get("business_unit_id")
	if unitParam == "" {
		unitParam = q.Get("cost_center_id")
	}
	if unitParam != "" {
		uid, err := uuid.Parse(unitParam)
		if err != nil {
			return nil, fmt.Errorf("invalid 'business_unit_id' format: must be valid UUID")
		}
		businessUnitID = &uid
	}

	// 8. Branch ID validation
	var branchID *uuid.UUID
	if b := q.Get("branch_id"); b != "" {
		bid, err := uuid.Parse(b)
		if err != nil {
			return nil, fmt.Errorf("invalid 'branch_id' format: must be valid UUID")
		}
		branchID = &bid
	}

	// 9. Status filter
	status := strings.TrimSpace(q.Get("status"))

	// 10. Cursor for keyset pagination
	cursor := strings.TrimSpace(q.Get("cursor"))

	return &Query{
		Page:           page,
		Limit:          limit,
		Offset:         offset,
		Search:         search,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		TieBreaker:     cfg.TieBreakerColumn,
		StartDate:      startDate,
		EndDate:        endDate,
		BusinessUnitID: businessUnitID,
		BranchID:       branchID,
		Status:         status,
		Cursor:         cursor,
	}, nil
}

// BuildOrderBy produces a safe, deterministic SQL ORDER BY clause with a tie-breaker.
func (q *Query) BuildOrderBy(tablePrefix string) string {
	sortCol := q.SortBy
	tieBreaker := q.TieBreaker
	if tieBreaker == "" {
		tieBreaker = "id"
	}

	if tablePrefix != "" {
		tablePrefix = strings.TrimSuffix(tablePrefix, ".") + "."
		if !strings.Contains(sortCol, ".") {
			sortCol = tablePrefix + sortCol
		}
		if !strings.Contains(tieBreaker, ".") {
			tieBreaker = tablePrefix + tieBreaker
		}
	}

	if sortCol == tieBreaker {
		return fmt.Sprintf("ORDER BY %s %s", sortCol, q.SortOrder)
	}

	return fmt.Sprintf("ORDER BY %s %s, %s %s", sortCol, q.SortOrder, tieBreaker, q.SortOrder)
}

// ToMeta builds standardized response.PaginationMeta from total items count.
func (q *Query) ToMeta(totalItems int64) response.PaginationMeta {
	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((totalItems + int64(q.Limit) - 1) / int64(q.Limit))
	}
	return response.PaginationMeta{
		CurrentPage: q.Page,
		PerPage:     q.Limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}
}

// EncodeCursor creates an opaque base64 string from a timestamp and UUID for cursor/keyset pagination.
func EncodeCursor(t time.Time, id uuid.UUID) string {
	raw := fmt.Sprintf("%d,%s", t.UnixNano(), id.String())
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor decodes an opaque cursor string into timestamp and UUID.
func DecodeCursor(cursor string) (time.Time, uuid.UUID, error) {
	bytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}

	parts := strings.Split(string(bytes), ",")
	if len(parts) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor structure")
	}

	nano, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor timestamp")
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor id")
	}

	return time.Unix(0, nano), id, nil
}

func parseDate(val string, endOfDay bool) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, val); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", val); err == nil {
		if endOfDay {
			return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC), nil
		}
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("invalid date")
}

func getMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
