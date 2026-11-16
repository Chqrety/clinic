package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"clinic-api/internal/queue/domain"
	"clinic-api/internal/queue/usecase"
	"clinic-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type QueueHandler struct {
	usecase *usecase.QueueUsecase
}

func NewQueueHandler(u *usecase.QueueUsecase) *QueueHandler {
	return &QueueHandler{usecase: u}
}

func (h *QueueHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	result, err := h.usecase.RegisterQueue(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyPatientName),
			errors.Is(err, domain.ErrEmptyPatientPhone),
			errors.Is(err, domain.ErrInvalidClinicUnit),
			errors.Is(err, domain.ErrInvalidCategory),
			errors.Is(err, domain.ErrInvalidPaymentMethod):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memproses pendaftaran antrian: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *QueueHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QueueID int64 `json:"queue_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	queueItem, err := h.usecase.VerifyPayment(r.Context(), body.QueueID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrQueueNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrPaymentAlreadyProcessed):
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal verifikasi pembayaran: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusOK, queueItem)
}

func (h *QueueHandler) HandleTripayWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "gagal membaca body webhook")
		return
	}
	defer r.Body.Close()

	sig := r.Header.Get("X-Callback-Signature")

	if err := h.usecase.ProcessTripayWebhook(r.Context(), rawBody, sig); err != nil {
		if strings.Contains(err.Error(), "signature") {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *QueueHandler) HandleMidtransWebhook(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "gagal membaca body webhook midtrans")
		return
	}
	defer r.Body.Close()

	if err := h.usecase.ProcessMidtransWebhook(r.Context(), rawBody); err != nil {
		if strings.Contains(err.Error(), "signature") {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"})
}

func (h *QueueHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "parameter ID tiket tidak valid")
		return
	}

	ticket, err := h.usecase.GetTicket(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrQueueNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, ticket)
}

func (h *QueueHandler) GetMonitor(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	data, err := h.usecase.GetMonitorData(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil data monitor antrian: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, data)
}

func (h *QueueHandler) StreamMonitor(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		response.Error(w, http.StatusInternalServerError, "streaming sse tidak didukung oleh browser/koneksi")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	unit := r.URL.Query().Get("unit")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial push
	if data, err := h.usecase.GetMonitorData(r.Context(), unit); err == nil {
		payload, _ := json.Marshal(data)
		if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
			return
		}
		flusher.Flush()
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			data, err := h.usecase.GetMonitorData(r.Context(), unit)
			if err != nil {
				continue
			}
			payload, _ := json.Marshal(data)
			if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func (h *QueueHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID antrian tidak valid")
		return
	}

	queueItem, err := h.usecase.CheckIn(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, queueItem)
}

func (h *QueueHandler) ListReservations(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	list, err := h.usecase.ListReservations(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat daftar reservasi: "+err.Error())
		return
	}

	if list == nil {
		list = []domain.QueueItem{}
	}
	response.Success(w, http.StatusOK, list)
}

func (h *QueueHandler) ListRefundRequests(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	list, err := h.usecase.ListRefundRequests(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat daftar pengajuan refund: "+err.Error())
		return
	}

	if list == nil {
		list = []domain.QueueItem{}
	}
	response.Success(w, http.StatusOK, list)
}

func (h *QueueHandler) GetDoctorRoomsStatus(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.usecase.GetDoctorRoomsStatus(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat status ruangan dokter: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, rooms)
}

func (h *QueueHandler) LookupPatient(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if strings.TrimSpace(query) == "" {
		response.Error(w, http.StatusBadRequest, "parameter query (telepon atau nama) wajib diisi")
		return
	}

	patient, err := h.usecase.LookupPatient(r.Context(), query)
	if err != nil {
		if errors.Is(err, domain.ErrPatientNotFound) {
			response.Error(w, http.StatusNotFound, "data pasien tidak ditemukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal mencari data pasien: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, patient)
}

func (h *QueueHandler) Reschedule(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QueueID int64  `json:"queue_id"`
		NewDate string `json:"new_date"`
		NewTime string `json:"new_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	if body.QueueID <= 0 || strings.TrimSpace(body.NewDate) == "" || strings.TrimSpace(body.NewTime) == "" {
		response.Error(w, http.StatusBadRequest, "queue_id, new_date, dan new_time wajib diisi")
		return
	}

	updated, err := h.usecase.RescheduleQueue(r.Context(), body.QueueID, body.NewDate, body.NewTime)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, updated)
}

func (h *QueueHandler) RequestRefund(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QueueID       int64  `json:"queue_id"`
		Bank          string `json:"bank"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		Notes         string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	if body.QueueID <= 0 || strings.TrimSpace(body.Bank) == "" || strings.TrimSpace(body.AccountNumber) == "" || strings.TrimSpace(body.AccountName) == "" {
		response.Error(w, http.StatusBadRequest, "queue_id, bank, nomor rekening, dan nama pemilik rekening wajib diisi")
		return
	}

	updated, err := h.usecase.RequestRefund(r.Context(), body.QueueID, body.Bank, body.AccountNumber, body.AccountName, body.Notes)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, updated)
}

func (h *QueueHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QueueID int64  `json:"queue_id"`
		Status  string `json:"status"` // "APPROVED" | "REJECTED"
		Notes   string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload json tidak valid")
		return
	}

	if body.QueueID <= 0 || (body.Status != "APPROVED" && body.Status != "REJECTED") {
		response.Error(w, http.StatusBadRequest, "queue_id dan status ('APPROVED'/'REJECTED') wajib diisi")
		return
	}

	updated, err := h.usecase.ProcessRefund(r.Context(), body.QueueID, body.Status, body.Notes)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, updated)
}
