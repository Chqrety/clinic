package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterAssetRoutes registers endpoints for fixed assets and monthly depreciation schedules.
func RegisterAssetRoutes(r chi.Router, h *handler.AssetHandler) {
	if h == nil {
		return
	}
	r.Route("/assets", func(ar chi.Router) {
		ar.Post("/", h.HandleRegisterAsset)
		ar.Get("/", h.HandleListAssets)
		ar.Post("/depreciate-batch", h.HandleRunMonthlyDepreciation)
		ar.Get("/{id}", h.HandleGetAssetByID)
		ar.Post("/{id}/depreciate", h.HandleDepreciateAsset)
		ar.Get("/{id}/depreciations", h.HandleGetDepreciationHistory)
	})
}

