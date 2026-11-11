package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterPayableRoutes mounts Account Payable (AP) endpoints under /api/v1/payables.
func RegisterPayableRoutes(r chi.Router, h *handler.PayableHandler) {
	if h == nil {
		return
	}
	r.Route("/payables", func(ap chi.Router) {
		ap.Get("/", h.HandleListPayables)
		ap.Post("/", h.HandleCreatePayable)
		ap.Get("/aging", h.HandleGetAgingReport)
		ap.Get("/{id}", h.HandleGetPayableByID)
		ap.Post("/{id}/pay", h.HandlePayPayable)
	})
}

