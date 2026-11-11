package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPharmacyRoutes(
	r chi.Router,
	h *PharmacyHandler,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route("/api/pharmacy", func(pr chi.Router) {
		pr.Use(authMiddleware)

		// 1. Incoming Prescriptions & Dispensing Lifecycle
		pr.Get("/prescriptions", h.HandleListPrescriptions)
		pr.Get("/prescriptions/{id}", h.HandleGetPrescription)
		pr.Post("/prescriptions/{id}/status", h.HandleStartPreparing)
		pr.Post("/prescriptions/{id}/dispense", h.HandleDispense)
		pr.Post("/prescriptions/{id}/cancel", h.HandleCancel)

		// 2. Drugs & FEFO Batches
		pr.Get("/drugs", h.HandleListDrugs)
		pr.Post("/drugs", h.HandleCreateDrug)
		pr.Post("/drugs/{id}/restock", h.HandleRestock)

		// 3. Stock Opname CSV Export & Import
		pr.Get("/opname/export", h.HandleExportOpnameCSV)
		pr.Post("/opname/import", h.HandleImportOpnameCSV)
		pr.Get("/opname/logs", h.HandleListOpnameLogs)
	})
}
