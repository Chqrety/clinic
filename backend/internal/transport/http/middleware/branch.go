package middleware

import (
	"context"
	"net/http"

	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

type contextKey string

const (
	BranchIDKey contextKey = "branch_id"
	UserIDKey   contextKey = "user_id"
	RoleKey     contextKey = "user_role"
)

// BranchContext extracts and validates the X-Branch-ID header for multi-branch accounting isolation.
func BranchContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		branchHeader := r.Header.Get("X-Branch-ID")
		if branchHeader == "" {
			// In development or if single-branch clinic, allow default fallback or require header
			next.ServeHTTP(w, r)
			return
		}

		branchID, err := uuid.Parse(branchHeader)
		if err != nil {
			response.BadRequest(w, "Invalid X-Branch-ID header format (UUID required)", nil)
			return
		}

		ctx := context.WithValue(r.Context(), BranchIDKey, branchID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetBranchID retrieves the branch UUID from request context.
func GetBranchID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(BranchIDKey).(uuid.UUID)
	return val, ok
}

