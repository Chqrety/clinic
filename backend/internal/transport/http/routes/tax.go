package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterTaxRoutes registers tax withholding and e-Bupot 21/26 DJP reporting endpoints.
func RegisterTaxRoutes(r chi.Router, h *handler.TaxHandler) {
	if h == nil {
		return
	}
	r.Route("/taxes", func(tr chi.Router) {
		tr.Post("/withholdings", h.HandleCreateWithholding)
		tr.Get("/withholdings", h.HandleListWithholdings)
		tr.Get("/summary", h.HandleGetTaxSummary)
		tr.Get("/ebupot21/export", h.HandleExportEBupot21CSV)
		tr.Post("/withholdings/{id}/report", h.HandleMarkAsReported)
	})
}

