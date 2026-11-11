package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"clinic-api/internal/shared/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmpNumKey contextKey = "user_emp_num"
	UserNameKey  contextKey = "user_name"
	UserRoleKey  contextKey = "user_role"
	UserUnitKey  contextKey = "user_unit"
)

type StaffClaims struct {
	UserID         int64  `json:"user_id"`
	EmployeeNumber string `json:"employee_number"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	ClinicUnit     string `json:"clinic_unit"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, userID int64, empNum, name, role, unit string, duration time.Duration) (string, error) {
	claims := StaffClaims{
		UserID:         userID,
		EmployeeNumber: empNum,
		Name:           name,
		Role:           role,
		ClinicUnit:     unit,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   empNum,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, tokenString string) (*StaffClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StaffClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*StaffClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	return claims, nil
}

func Authenticate(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			var tokenStr string

			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			} else if cookie, err := r.Cookie("clinic_token"); err == nil {
				tokenStr = cookie.Value
			}

			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "akses ditolak: sesi login tidak ditemukan")
				return
			}

			claims, err := ValidateToken(secret, tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "sesi login tidak valid atau kadaluarsa: "+err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmpNumKey, claims.EmployeeNumber)
			ctx = context.WithValue(ctx, UserNameKey, claims.Name)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			ctx = context.WithValue(ctx, UserUnitKey, claims.ClinicUnit)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(UserRoleKey)
			if roleVal == nil {
				response.Error(w, http.StatusUnauthorized, "sesi pengguna tidak teridentifikasi")
				return
			}

			userRole, ok := roleVal.(string)
			if !ok {
				response.Error(w, http.StatusForbidden, "tipe peran pengguna tidak valid")
				return
			}

			allowed := false
			for _, r := range allowedRoles {
				if r == userRole {
					allowed = true
					break
				}
			}

			if !allowed {
				response.Error(w, http.StatusForbidden, "akses ditolak: peran '"+userRole+"' tidak memiliki izin untuk tindakan ini")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
