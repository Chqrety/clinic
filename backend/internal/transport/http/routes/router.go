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

	// 3. Mount Separated API v1 Routes
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(middleware.BranchContext)

		if h != nil {
			// A. Chart of Accounts (COA)
			RegisterAccountingRoutes(v1, h.Accounting)

			// B. Journal Management
			RegisterJournalRoutes(v1, h.Accounting)

			// C. Cash Management
			RegisterCashRoutes(v1, h.CashBank)

			// D. Bank Management
			RegisterBankRoutes(v1, h.CashBank)

			// E. Financial Reporting (GL, Trial Balance, Laba Rugi, Neraca)
			RegisterReportingRoutes(v1, h.Reporting)

			// F. Fiscal Periods & Year-End Closing Management
			v1.Route("/periods", func(p chi.Router) {
				if h.CashBank != nil {
					p.Post("/close", h.CashBank.HandleClosePeriod)
					p.Post("/{id}/reopen", h.CashBank.HandleReopenPeriod)
				}
				if h.Closing != nil {
					p.Post("/year-end-close", h.Closing.HandleYearEndClosing)
				}
			})

			// G. Billing & Kasir POS
			RegisterBillingRoutes(v1, h.Billing)

			// H. Account Receivable (Piutang)
			RegisterReceivableRoutes(v1, h.Receivable)

			// I. Account Payable (Hutang)
			RegisterPayableRoutes(v1, h.Payable)

			// J. Doctor Sharing & Settlements
			RegisterDoctorShareRoutes(v1, h.DoctorShare)

			// K. Treasury & Cashier Shifts
			RegisterTreasuryRoutes(v1, h.Treasury)

			// L. Tax & e-Bupot 21/26
			RegisterTaxRoutes(v1, h.Tax)

			// M. Fixed Assets & Depreciation
			RegisterAssetRoutes(v1, h.Asset)

			// N. Prepaid Expenses & Amortization
			RegisterPrepaidRoutes(v1, h.Prepaid)

			// O. Cost Center Budgeting & Variance Analysis
			RegisterBudgetRoutes(v1, h.Budget)

			// P. Financial Audit Logs
			RegisterAuditRoutes(v1, h.Audit)
		}
	})

	return r
}
