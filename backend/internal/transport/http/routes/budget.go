package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterBudgetRoutes registers endpoints for cost center budgeting and variance reporting.
func RegisterBudgetRoutes(r chi.Router, h *handler.BudgetHandler) {
	if h == nil {
		return
	}
	r.Route("/budgets", func(br chi.Router) {
		br.Post("/", h.HandleCreateBudget)
		br.Get("/", h.HandleListBudgets)
		br.Get("/variance", h.HandleGetVarianceReport)
		br.Get("/{id}", h.HandleGetBudgetByID)
		br.Post("/{id}/approve", h.HandleApproveBudget)
	})
}

