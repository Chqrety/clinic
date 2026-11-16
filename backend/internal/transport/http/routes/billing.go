package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterBillingRoutes mounts patient invoicing and cashier billing endpoints under /api/v1/billing.
func RegisterBillingRoutes(r chi.Router, h *handler.BillingHandler) {
	if h == nil {
		return
	}
	r.Route("/billing", func(b chi.Router) {
		h.RegisterRoutes(b)
	})
}

