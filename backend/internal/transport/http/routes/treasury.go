package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterTreasuryRoutes registers cashier shifts and treasury physical cash reconciliation endpoints.
func RegisterTreasuryRoutes(r chi.Router, h *handler.TreasuryHandler) {
	if h == nil {
		return
	}
	r.Route("/treasury", func(tr chi.Router) {
		tr.Post("/shifts/open", h.HandleOpenShift)
		tr.Get("/shifts/active", h.HandleGetActiveShift)
		tr.Get("/shifts", h.HandleListShifts)
		tr.Get("/shifts/{id}", h.HandleGetShift)
		tr.Post("/shifts/{id}/close", h.HandleCloseShift)
		tr.Post("/shifts/{id}/deposit", h.HandleDepositCash)
	})
}

