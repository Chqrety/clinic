package routes

import (
	"net/http"

	"clinic-api/internal/transport/http/handler"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type HandlerContainer struct {
	Accounting  *handler.AccountingHandler
	CashBank    *handler.CashBankHandler
	Reporting   *handler.ReportHandler
	Billing     *handler.BillingHandler
	Receivable  *handler.ReceivableHandler
	Payable     *handler.PayableHandler
	DoctorShare *handler.DoctorShareHandler
	Treasury    *handler.TreasuryHandler
	Tax         *handler.TaxHandler
	Asset       *handler.AssetHandler
	Prepaid     *handler.PrepaidHandler
	Budget      *handler.BudgetHandler
	Closing     *handler.ClosingHandler
	Audit       *handler.AuditHandler
}

// SetupRouter initializes the root Chi Mux with global middlewares and domain route modules.
func SetupRouter(h *HandlerContainer) *chi.Mux {
	r := chi.NewRouter()

	// 1. Global Middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Branch-ID", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 2. Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, "Clinic Accounting System is healthy & operational", map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		}, nil)
	})

	// 3. Mount Separated Accounting API Routes under /api/accounting
	MountAccountingRoutes(r, h)

	return r
}

// MountAccountingRoutes mounts all accounting submodules under /api/accounting with optional JWT auth protection.
func MountAccountingRoutes(r chi.Router, h *HandlerContainer, authMiddlewares ...func(http.Handler) http.Handler) {
	r.Route("/api/accounting", func(acc chi.Router) {
		for _, mw := range authMiddlewares {
			if mw != nil {
				acc.Use(mw)
			}
		}
		acc.Use(middleware.BranchContext)

		if h != nil {
			// A. Chart of Accounts (COA)
			RegisterAccountingRoutes(acc, h.Accounting)

			// B. Journal Management
			RegisterJournalRoutes(acc, h.Accounting)

			// C. Cash Management
			RegisterCashRoutes(acc, h.CashBank)

			// D. Bank Management
			RegisterBankRoutes(acc, h.CashBank)

			// E. Financial Reporting (GL, Trial Balance, Laba Rugi, Neraca)
			RegisterReportingRoutes(acc, h.Reporting)

			// F. Fiscal Periods & Year-End Closing Management
			acc.Route("/periods", func(p chi.Router) {
				if h.CashBank != nil {
					p.Post("/close", h.CashBank.HandleClosePeriod)
					p.Post("/{id}/reopen", h.CashBank.HandleReopenPeriod)
				}
				if h.Closing != nil {
					p.Post("/year-end-close", h.Closing.HandleYearEndClosing)
				}
			})

			// G. Billing & Kasir POS
			RegisterBillingRoutes(acc, h.Billing)

			// H. Account Receivable (Piutang)
			RegisterReceivableRoutes(acc, h.Receivable)

			// I. Account Payable (Hutang)
			RegisterPayableRoutes(acc, h.Payable)

			// J. Doctor Sharing & Settlements
			RegisterDoctorShareRoutes(acc, h.DoctorShare)

			// K. Treasury & Cashier Shifts
			RegisterTreasuryRoutes(acc, h.Treasury)

			// L. Tax & e-Bupot 21/26
			RegisterTaxRoutes(acc, h.Tax)

			// M. Fixed Assets & Depreciation
			RegisterAssetRoutes(acc, h.Asset)

			// N. Prepaid Expenses & Amortization
			RegisterPrepaidRoutes(acc, h.Prepaid)

			// O. Cost Center Budgeting & Variance Analysis
			RegisterBudgetRoutes(acc, h.Budget)

			// P. Financial Audit Logs
			RegisterAuditRoutes(acc, h.Audit)
		}
	})
}
