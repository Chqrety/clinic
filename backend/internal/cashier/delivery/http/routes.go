package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterCashierRoutes(
	r chi.Router,
	h *CashierHandler,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route("/api/cashier", func(cr chi.Router) {
		cr.Use(authMiddleware)

		// 1. Walk-In Registration
		cr.Post("/walkin", h.HandleRegisterWalkIn)

		// 2. Pending Bills & Billing
		cr.Get("/pending-bills", h.HandleListPendingBills)
		cr.Post("/invoices", h.HandlePayBill)
		cr.Get("/invoices", h.HandleListInvoices)
		cr.Get("/invoices/{id}", h.HandleGetInvoice)
		cr.Get("/invoices/queue/{queueId}", h.HandleGetInvoiceByQueue)

		// 3. Refund
		cr.Post("/invoices/{id}/refund", h.HandleRefund)

		// 4. End-of-Shift Reports
		cr.Get("/shift-summary", h.HandleGetShiftSummary)
		cr.Post("/close-shift", h.HandleCloseShift)
		cr.Get("/shift-reports", h.HandleListShiftReports)
	})
}
