package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterAuditRoutes registers endpoints for financial audit trails.
func RegisterAuditRoutes(r chi.Router, h *handler.AuditHandler) {
	if h == nil {
		return
	}
	r.Route("/audit-logs", func(ar chi.Router) {
		ar.Get("/", h.HandleListAuditLogs)
	})
}

