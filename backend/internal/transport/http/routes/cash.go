package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterCashRoutes registers routes for Cash Drawer, Petty Cash, Cash-in, and Cash-out.
func RegisterCashRoutes(r chi.Router, h *handler.CashBankHandler) {
	if h == nil {
		return
	}
	r.Route("/cash", func(c chi.Router) {
		c.Get("/accounts", h.HandleListAccounts)
		c.Post("/accounts", h.HandleCreateAccount)
		c.Post("/in", h.HandleRecordCashIn)
		c.Post("/out", h.HandleRecordCashOut)
		c.Post("/transfer", h.HandleTransfer)
		c.Get("/transactions", h.HandleListTransactions)
	})
}

