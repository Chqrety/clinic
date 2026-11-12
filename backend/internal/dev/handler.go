package dev

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type DevHandler struct {
	db *sql.DB
}

func NewDevHandler(db *sql.DB) *DevHandler {
	return &DevHandler{db: db}
}

func (h *DevHandler) ResetDemo(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		h.jsonError(w, http.StatusServiceUnavailable, "Database belum terhubung")
		return
	}
	if err := SeedInitialDemo(h.db); err != nil {
		h.jsonError(w, http.StatusInternalServerError, "Gagal reset demo: "+err.Error())
		return
	}
	h.jsonSuccess(w, "Berhasil me-reset data awal demo (8 antrian terpasang)")
}

func (h *DevHandler) RushHour(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		h.jsonError(w, http.StatusServiceUnavailable, "Database belum terhubung")
		return
	}
	if err := SeedRushHour(h.db); err != nil {
		h.jsonError(w, http.StatusInternalServerError, "Gagal seed rush hour: "+err.Error())
		return
	}
	h.jsonSuccess(w, "Berhasil memuat skenario jam sibuk (50 antrian realistis lintas Kidinest & SSC)")
}

func (h *DevHandler) WipeTransactions(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		h.jsonError(w, http.StatusServiceUnavailable, "Database belum terhubung")
		return
	}
	if err := WipeTransactions(h.db); err != nil {
		h.jsonError(w, http.StatusInternalServerError, "Gagal membersihkan transaksi: "+err.Error())
		return
	}
	h.jsonSuccess(w, "Seluruh transaksi antrian, invoice, dan resep berhasil dibersihkan (Zero State)")
}

func (h *DevHandler) jsonSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": msg,
	})
}

func (h *DevHandler) jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   msg,
	})
}

// RegisterDevRoutes registers dev endpoints with strict guard (disabled in production).
func RegisterDevRoutes(r chi.Router, h *DevHandler) {
	// Guard: Never mount in production
	if os.Getenv("ENV") == "production" {
		return
	}

	r.Route("/api/dev/seeds", func(dr chi.Router) {
		dr.Post("/reset", h.ResetDemo)
		dr.Post("/rush-hour", h.RushHour)
		dr.Post("/wipe-transactions", h.WipeTransactions)
	})
}
