package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clinic-api/internal/schedule/domain"
	"clinic-api/internal/schedule/usecase"
	"clinic-api/internal/shared/response"

	"github.com/go-chi/chi/v5"
)

type ScheduleHandler struct {
	usecase *usecase.ScheduleUsecase
}

func NewScheduleHandler(uc *usecase.ScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{usecase: uc}
}

func (h *ScheduleHandler) GetDoctors(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	doctors, err := h.usecase.GetDoctors(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil daftar dokter: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, doctors)
}

func (h *ScheduleHandler) GetDoctorSchedules(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	doctorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID dokter tidak valid")
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	availDates, err := h.usecase.GetDoctorAvailableSchedule(r.Context(), doctorID, days)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menghitung jadwal dokter: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, availDates)
}

func (h *ScheduleHandler) AddDoctorLeave(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DoctorID int64  `json:"doctor_id"`
		Date     string `json:"date"`
		Reason   string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	if err := h.usecase.AddDoctorLeave(r.Context(), body.DoctorID, body.Date, body.Reason); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menambah cuti dokter: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "jadwal cuti dokter berhasil disimpan"})
}

func (h *ScheduleHandler) DeleteDoctorLeave(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DoctorID int64  `json:"doctor_id"`
		Date     string `json:"date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	if err := h.usecase.DeleteDoctorLeave(r.Context(), body.DoctorID, body.Date); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal membatalkan cuti dokter: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "cuti dokter berhasil dibatalkan"})
}

func (h *ScheduleHandler) GetDoctorLeaveList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	doctorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID dokter tidak valid")
		return
	}

	leaves, err := h.usecase.GetDoctorLeaveList(r.Context(), doctorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil daftar cuti: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, leaves)
}

func (h *ScheduleHandler) GetDoctorWeeklySchedules(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	doctorID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID dokter tidak valid")
		return
	}

	schedules, err := h.usecase.GetDoctorWeeklySchedules(r.Context(), doctorID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil jadwal mingguan: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, schedules)
}

func (h *ScheduleHandler) UpdateDoctorSchedule(w http.ResponseWriter, r *http.Request) {
	var body domain.DoctorSchedule
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	if body.DoctorID <= 0 || body.StartTime == "" || body.EndTime == "" {
		response.Error(w, http.StatusBadRequest, "doctor_id, start_time, dan end_time wajib diisi")
		return
	}

	if err := h.usecase.UpdateDoctorSchedule(r.Context(), &body); err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memperbarui jadwal dokter: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "jadwal shift dokter berhasil diperbarui"})
}

func (h *ScheduleHandler) GetPublicSchedulesSummary(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	summaries, err := h.usecase.GetPublicSchedulesSummary(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil ringkasan jadwal dokter: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, summaries)
}


