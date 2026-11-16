package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterReportingRoutes registers routes for General Ledger, Trial Balance, Income Statement, Balance Sheet, and Cash Flow.
func RegisterReportingRoutes(r chi.Router, h *handler.ReportHandler) {
	if h == nil {
		return
	}
	r.Route("/reports", func(rep chi.Router) {
		rep.Get("/general-ledger", h.HandleGetGeneralLedger)
		rep.Get("/trial-balance", h.HandleGetTrialBalance)
		rep.Get("/income-statement", h.HandleGetIncomeStatement)
		rep.Get("/balance-sheet", h.HandleGetBalanceSheet)
		rep.Get("/cash-flow", h.HandleGetCashFlowStatement)
	})
}

