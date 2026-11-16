package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterDoctorShareRoutes registers all routes for doctor sharing rules and monthly settlements.
func RegisterDoctorShareRoutes(r chi.Router, h *handler.DoctorShareHandler) {
	if h == nil {
		return
	}
	r.Route("/doctor-shares", func(dr chi.Router) {
		dr.Post("/rules", h.HandleCreateRule)
		dr.Get("/rules", h.HandleListRules)

		dr.Post("/settlements/generate", h.HandleGenerateSettlement)
		dr.Get("/settlements", h.HandleListSettlements)
		dr.Get("/settlements/{id}", h.HandleGetSettlement)
		dr.Post("/settlements/{id}/approve", h.HandleApproveSettlement)
		dr.Post("/settlements/{id}/pay", h.HandlePaySettlement)
	})
}

