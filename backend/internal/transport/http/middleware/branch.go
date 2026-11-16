package middleware

import (
	"context"
	"net/http"

	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

type contextKey string

const (
	BranchIDKey    contextKey = "branch_id"
	UserIDKey      contextKey = "user_id"
	RoleKey        contextKey = "user_role"
	UserSessionKey contextKey = "user_session"
)

// UserSession holds authenticated session details for branch & cost center scoping.
type UserSession struct {
	UserID               uuid.UUID
	Role                 string
	AllowedBranchIDs     []uuid.UUID
	AllowedCostCenterIDs []uuid.UUID
	IsSuperAdmin         bool
}

var DefaultBranchID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

// BranchContext extracts and validates the X-Branch-ID header for multi-branch accounting isolation.
func BranchContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		branchHeader := r.Header.Get("X-Branch-ID")
		if branchHeader == "" {
			// In development or single-branch clinic, proceed with default branch context
			ctx := context.WithValue(r.Context(), BranchIDKey, DefaultBranchID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		branchID, err := uuid.Parse(branchHeader)
		if err != nil {
			response.BadRequest(w, "Invalid X-Branch-ID header format (UUID required)", nil)
			return
		}

		// Enforce branch access authorization if user session is present in context
		if sess, ok := GetUserSession(r.Context()); ok && !sess.IsSuperAdmin {
			allowed := false
			for _, b := range sess.AllowedBranchIDs {
				if b == branchID {
					allowed = true
					break
				}
			}
			if !allowed && len(sess.AllowedBranchIDs) > 0 {
				response.Forbidden(w, "Access denied: you do not have permission to view or manage this branch")
				return
			}
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

// SetUserSession attaches an authenticated user session to context (e.g. from JWT middleware).
func SetUserSession(ctx context.Context, session *UserSession) context.Context {
	return context.WithValue(ctx, UserSessionKey, session)
}

// GetUserSession retrieves user session from context.
func GetUserSession(ctx context.Context) (*UserSession, bool) {
	val, ok := ctx.Value(UserSessionKey).(*UserSession)
	return val, ok
}

// ValidateUnitAuthorization verifies if the current user is permitted to query the given cost center / unit.
func ValidateUnitAuthorization(ctx context.Context, unitID uuid.UUID) bool {
	sess, ok := GetUserSession(ctx)
	if !ok || sess.IsSuperAdmin || len(sess.AllowedCostCenterIDs) == 0 {
		return true // Allow if superadmin or no restrictive scope configured
	}
	for _, u := range sess.AllowedCostCenterIDs {
		if u == unitID {
			return true
		}
	}
	return false
}
