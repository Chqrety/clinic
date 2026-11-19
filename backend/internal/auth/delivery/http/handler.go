package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"clinic-api/internal/auth/domain"
	"clinic-api/internal/auth/usecase"
	"clinic-api/internal/shared/middleware"
	"clinic-api/internal/shared/response"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler(u *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: u}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	res, err := h.usecase.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrStaffInactive):
			response.Error(w, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal server: "+err.Error())
		}
		return
	}

	// Set httpOnly cookie for convenience
	http.SetCookie(w, &http.Cookie{
		Name:     "clinic_token",
		Value:    res.Token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response.Success(w, http.StatusOK, res)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(middleware.UserIDKey)
	if userIDVal == nil {
		response.Error(w, http.StatusUnauthorized, "pengguna tidak terautentikasi")
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "data identitas sesi tidak valid")
		return
	}

	profile, err := h.usecase.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, http.StatusOK, profile)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil token dari Authorization header atau cookie untuk di-blacklist
	authHeader := r.Header.Get("Authorization")
	var tokenStr string
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := r.Cookie("clinic_token"); err == nil {
		tokenStr = cookie.Value
	}

	// 2. Blacklist / Cabut token di sistem backend
	if tokenStr != "" {
		middleware.RevokeToken(tokenStr)
	}

	// 3. Clear cookie clinic_token secara tuntas
	http.SetCookie(w, &http.Cookie{
		Name:     "clinic_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	response.Success(w, http.StatusOK, map[string]string{
		"message": "berhasil logout dari sistem klinik, token dan sesi backend telah dicabut",
	})
}
