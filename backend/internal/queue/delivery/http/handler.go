package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")

	unit := r.URL.Query().Get("unit")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial push
	if data, err := h.usecase.GetMonitorData(r.Context(), unit); err == nil {
		payload, _ := json.Marshal(data)
		fmt.Fprintf(w, "data: %s\n\n", payload)
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
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}
