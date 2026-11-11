package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterBankRoutes registers routes for Bank Accounts (BCA, Mandiri) and Bank Mutations.
func RegisterBankRoutes(r chi.Router, h *handler.CashBankHandler) {
	if h == nil {
		return
	}
	r.Route("/bank", func(b chi.Router) {
		b.Get("/accounts", h.HandleListAccounts)
		b.Post("/accounts", h.HandleCreateAccount)
		b.Post("/transfer", h.HandleTransfer)
		b.Get("/transactions", h.HandleListTransactions)
	})
}

