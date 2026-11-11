package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"clinic-api/internal/shared/response"
	"clinic-api/internal/staff/domain"
	"clinic-api/internal/staff/usecase"
	"github.com/go-chi/chi/v5"
)

type StaffHandler struct {
	usecase *usecase.StaffUsecase
}

func NewStaffHandler(u *usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{usecase: u}
}

func (h *StaffHandler) List(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	role := r.URL.Query().Get("role")
	activeParam := r.URL.Query().Get("is_active")

	var isActive *bool
	if activeParam != "" {
		b := activeParam == "true" || activeParam == "1"
		isActive = &b
	}

	filter := domain.StaffFilter{
		ClinicUnit: unit,
		Role:       role,
		IsActive:   isActive,
	}

	staffs, err := h.usecase.ListStaff(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil daftar staf: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, staffs)
}

func (h *StaffHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "parameter ID tidak valid")
		return
	}

	staff, err := h.usecase.GetStaff(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrStaffNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, staff)
}

func (h *StaffHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	staff, err := h.usecase.CreateStaff(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateEmployeeNum),
			errors.Is(err, domain.ErrInvalidRole),
			errors.Is(err, domain.ErrInvalidClinicUnit),
			errors.Is(err, domain.ErrEmptyEmployeeNumber),
			errors.Is(err, domain.ErrEmptyName),
			errors.Is(err, domain.ErrEmptyPassword):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal mendaftarkan staf: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusCreated, staff)
}

func (h *StaffHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "parameter ID tidak valid")
		return
	}

	var input domain.UpdateStaffInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	staff, err := h.usecase.UpdateStaff(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrStaffNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrEmptyName),
			errors.Is(err, domain.ErrInvalidRole),
			errors.Is(err, domain.ErrInvalidClinicUnit):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memperbarui data staf: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusOK, staff)
}

func (h *StaffHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "parameter ID tidak valid")
		return
	}

	var input domain.ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	if err := h.usecase.ResetPassword(r.Context(), id, input.NewPassword); err != nil {
		switch {
		case errors.Is(err, domain.ErrStaffNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrEmptyPassword):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal mereset password staf: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusOK, map[string]string{
		"message": "password staf berhasil direset oleh atasan",
	})
}
