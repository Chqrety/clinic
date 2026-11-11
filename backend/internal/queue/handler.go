package queue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"clinic-api/internal/domain"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, s *Service) {
	r.Get("/queues", s.HandleGetQueues)
	r.Post("/queues", s.HandleCreateQueue)
	r.Post("/queues/{id}/pay", s.HandlePayQueue)
	r.Post("/queues/{id}/status", s.HandleUpdateStatus)
}

func (s *Service) HandleGetQueues(w http.ResponseWriter, r *http.Request) {
	queues, err := s.GetQueues()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: queues})
}

func (s *Service) HandleCreateQueue(w http.ResponseWriter, r *http.Request) {
	var req CreateQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.NIK == "" || req.Name == "" {
		sendError(w, http.StatusBadRequest, "NIK dan Nama wajib diisi")
		return
	}

	queueItem, err := s.CreateQueue(req)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, domain.ApiResponse{
		Success: true,
		Data:    queueItem,
		Message: "Antrian berhasil didaftarkan",
	})
}

func (s *Service) HandlePayQueue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	if err := s.ConfirmPayment(id); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Pembayaran terverifikasi, nomor antrian aktif",
	})
}

func (s *Service) HandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid queue ID")
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.UpdateStatus(id, req.Status); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: fmt.Sprintf("Status antrian diubah menjadi %s", req.Status),
	})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}

