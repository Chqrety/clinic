package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterPrepaidRoutes registers endpoints for prepaid expenses and monthly amortizations.
func RegisterPrepaidRoutes(r chi.Router, h *handler.PrepaidHandler) {
	if h == nil {
		return
	}
	r.Route("/prepaid-expenses", func(pr chi.Router) {
		pr.Post("/", h.HandleRegisterPrepaid)
		pr.Get("/", h.HandleListActivePrepaid)
		pr.Post("/amortize-batch", h.HandleRunBatchAmortization)
		pr.Get("/{id}", h.HandleGetPrepaidByID)
		pr.Post("/{id}/amortize", h.HandleAmortize)
		pr.Get("/{id}/amortizations", h.HandleGetAmortizationHistory)
	})
}

