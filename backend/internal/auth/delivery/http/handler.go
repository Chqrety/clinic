package http

import (
	"encoding/json"
	"errors"
	"net/http"

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
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "clinic_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	response.Success(w, http.StatusOK, map[string]string{
		"message": "berhasil logout dari sistem klinik",
	})
}
