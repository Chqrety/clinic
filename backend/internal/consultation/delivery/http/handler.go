package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"clinic-api/internal/consultation/domain"
	"clinic-api/internal/consultation/usecase"
	"clinic-api/internal/shared/middleware"
	"clinic-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type ConsultationHandler struct {
	usecase *usecase.ConsultationUsecase
}

func NewConsultationHandler(u *usecase.ConsultationUsecase) *ConsultationHandler {
	return &ConsultationHandler{usecase: u}
}

func (h *ConsultationHandler) HandleListQueues(w http.ResponseWriter, r *http.Request) {
	doctorName := r.URL.Query().Get("doctor_name")
	unit := r.URL.Query().Get("unit")

	// If logged in as doctor, default doctorName if not specified
	if doctorName == "" {
		if roleVal := r.Context().Value(middleware.UserRoleKey); roleVal == "DOCTOR" {
			if nameVal := r.Context().Value(middleware.UserNameKey); nameVal != nil {
				if name, ok := nameVal.(string); ok {
					doctorName = name
				}
			}
		}
	}

	queues, err := h.usecase.ListDoctorQueues(r.Context(), doctorName, unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat antrian dokter: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, queues)
}

func (h *ConsultationHandler) HandleCallPatient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID antrian tidak valid")
		return
	}

	if err := h.usecase.CallPatient(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"queue_id": id, "status": "CALLING"})
}

func (h *ConsultationHandler) HandleStartExam(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID antrian tidak valid")
		return
	}

	if err := h.usecase.StartExam(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"queue_id": id, "status": "IN_EXAM"})
}

func (h *ConsultationHandler) HandleGetByQueue(w http.ResponseWriter, r *http.Request) {
	qStr := chi.URLParam(r, "queueId")
	qID, err := strconv.ParseInt(qStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID antrian tidak valid")
		return
	}

	cons, err := h.usecase.GetConsultation(r.Context(), qID)
	if err != nil {
		if errors.Is(err, domain.ErrConsultationNotFound) {
			response.Success(w, http.StatusOK, nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal memuat konsultasi: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, cons)
}

func (h *ConsultationHandler) HandleListDrugs(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	drugs, err := h.usecase.ListActiveDrugs(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat katalog obat: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, drugs)
}

func (h *ConsultationHandler) HandleSubmitConsultation(w http.ResponseWriter, r *http.Request) {
	var input domain.SaveConsultationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload konsultasi tidak valid")
		return
	}

	// Auto-fill doctor info from context if present and not provided
	if input.DoctorName == "" {
		if nameVal := r.Context().Value(middleware.UserNameKey); nameVal != nil {
			if name, ok := nameVal.(string); ok {
				input.DoctorName = name
			}
		}
	}
	if input.DoctorID == nil {
		if idVal := r.Context().Value(middleware.UserIDKey); idVal != nil {
			if id, ok := idVal.(int64); ok {
				input.DoctorID = &id
			}
		}
	}

	cons, err := h.usecase.SubmitConsultation(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDoctorNameRequired),
			errors.Is(err, domain.ErrDiagnosisRequired),
			errors.Is(err, domain.ErrInvalidQuantity),
			errors.Is(err, domain.ErrInsufficientStock):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal menyimpan konsultasi dan resep: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusCreated, cons)
}
