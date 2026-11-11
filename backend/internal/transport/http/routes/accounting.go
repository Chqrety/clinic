package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterAccountingRoutes registers routes for Chart of Accounts (COA) management.
func RegisterAccountingRoutes(r chi.Router, h *handler.AccountingHandler) {
	if h == nil {
		return
	}
	r.Route("/accounts", func(acc chi.Router) {
		acc.Get("/", h.HandleListAccounts)
		acc.Post("/", h.HandleCreateAccount)
		acc.Get("/{id}", h.HandleGetAccountByID)
	})
}

