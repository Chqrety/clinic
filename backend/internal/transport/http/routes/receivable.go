package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterReceivableRoutes mounts Account Receivable (AR) endpoints under /api/v1/receivables.
func RegisterReceivableRoutes(r chi.Router, h *handler.ReceivableHandler) {
	if h == nil {
		return
	}
	r.Route("/receivables", func(ar chi.Router) {
		ar.Get("/", h.HandleListReceivables)
		ar.Post("/", h.HandleCreateReceivable)
		ar.Get("/aging", h.HandleGetAgingReport)
		ar.Get("/{id}", h.HandleGetReceivableByID)
		ar.Post("/{id}/pay", h.HandlePayReceivable)
	})
}

