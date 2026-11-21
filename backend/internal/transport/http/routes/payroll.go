package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterPayrollRoutes registers all HTTP endpoints for Staff Payroll & Payslips management.
func RegisterPayrollRoutes(r chi.Router, h *handler.PayrollHandler) {
	if h == nil {
		return
	}
	r.Route("/payroll", func(pr chi.Router) {
		// Employee management
		pr.Get("/employees", h.HandleListEmployees)
		pr.Post("/employees", h.HandleCreateEmployee)
		pr.Get("/employees/{id}", h.HandleGetEmployee)
		pr.Put("/employees/{id}", h.HandleUpdateEmployee)

		// Payslip batch & individual management
		pr.Get("/slips", h.HandleListPayslips)
		pr.Post("/slips", h.HandleCreatePayslip)
		pr.Post("/slips/generate", h.HandleGenerateMonthlyPayroll)
		pr.Get("/slips/{id}", h.HandleGetPayslip)
		pr.Put("/slips/{id}", h.HandleUpdatePayslip)
		pr.Delete("/slips/{id}", h.HandleDeletePayslip)

		// Approval workflows
		pr.Post("/slips/{id}/approve", h.HandleApprovePayslip)
		pr.Post("/slips/approve-all", h.HandleApproveAllPayslips)

		// Disbursement / payment workflows
		pr.Post("/slips/{id}/pay", h.HandlePayPayslip)
		pr.Post("/slips/pay-all", h.HandlePayAllPayslips)

		// General Ledger Posting
		pr.Post("/post-ledger", h.HandlePostPayrollLedger)
	})
}
