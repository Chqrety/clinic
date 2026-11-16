package bootstrap

import (
	"database/sql"

	accountingApp "clinic-api/internal/application/accounting"
	assetApp "clinic-api/internal/application/asset"
	billingApp "clinic-api/internal/application/billing"
	cashbankApp "clinic-api/internal/application/cashbank"
	doctorshareApp "clinic-api/internal/application/doctorshare"
	journalApp "clinic-api/internal/application/journal"
	payableApp "clinic-api/internal/application/payable"
	receivableApp "clinic-api/internal/application/receivable"
	reportingApp "clinic-api/internal/application/reporting"
	taxApp "clinic-api/internal/application/tax"
	treasuryApp "clinic-api/internal/application/treasury"
	prepaidApp "clinic-api/internal/application/prepaid"
	budgetApp "clinic-api/internal/application/budget"
	auditApp "clinic-api/internal/application/audit"

	"clinic-api/internal/adapters/payment"
	"clinic-api/internal/adapters/persistence/postgres"
	"clinic-api/internal/transport/http/handler"
	"clinic-api/internal/transport/http/routes"
)

type Container struct {
	DB             *sql.DB
	Handlers       *routes.HandlerContainer
	AccountingSvc  *accountingApp.Service
	JournalSvc     *journalApp.Service
	CashBankSvc    *cashbankApp.Service
	ReportSvc      *accountingApp.ReportService
	ReportingSvc   *reportingApp.Service
	BillingSvc     *billingApp.Service
	ReceivableSvc  *receivableApp.Service
	PayableSvc     *payableApp.Service
	DoctorShareSvc *doctorshareApp.Service
	TreasurySvc    *treasuryApp.Service
	TaxSvc         *taxApp.Service
	AssetSvc       *assetApp.Service
	PrepaidSvc     *prepaidApp.Service
	BudgetSvc      *budgetApp.Service
	ClosingSvc     *accountingApp.ClosingService
	AuditSvc       *auditApp.Service
}

// BuildContainer performs pure idiomatic dependency injection without magic reflection.
func BuildContainer(db *sql.DB) *Container {
	// 1. Transaction Unit of Work
	uow := postgres.NewUnitOfWork(db)

	// 2. Repositories
	accountRepo := postgres.NewPostgresAccountRepo(db)
	journalRepo := postgres.NewPostgresJournalRepo(db)
	periodRepo := postgres.NewPostgresFiscalPeriodRepo(db)
	cashbankRepo := postgres.NewPostgresCashBankRepo(db)
	reportRepo := postgres.NewPostgresReportRepo(db)
	billingRepo := postgres.NewPostgresBillingRepo(db)
	receivableRepo := postgres.NewPostgresReceivableRepo(db)
	payableRepo := postgres.NewPostgresPayableRepo(db)
	doctorShareRepo := postgres.NewPostgresDoctorShareRepo(db)
	treasuryRepo := postgres.NewPostgresTreasuryRepo(db)
	taxRepo := postgres.NewPostgresTaxRepo(db)
	assetRepo := postgres.NewPostgresAssetRepo(db)
	prepaidRepo := postgres.NewPostgresPrepaidRepo(db)
	budgetRepo := postgres.NewPostgresBudgetRepo(db)
	auditRepo := postgres.NewPostgresAuditRepo(db)

	// 3. Application Services
	journalSvc := journalApp.NewService(uow, accountRepo, journalRepo, periodRepo)
	accountingSvc := accountingApp.NewService(uow, accountRepo, journalRepo, periodRepo)
	cashbankSvc := cashbankApp.NewService(uow, cashbankRepo, accountingSvc)
	reportSvc := accountingApp.NewReportService(reportRepo)
	reportingSvc := reportingApp.NewService(reportRepo)

	cashProc := payment.NewCashProcessor("1-1111")
	bankProc := payment.NewBankTransferProcessor("1-1121")
	billingSvc := billingApp.NewService(uow, billingRepo, accountingSvc, cashProc, bankProc)

	receivableSvc := receivableApp.NewService(uow, receivableRepo, accountRepo, journalSvc)
	payableSvc := payableApp.NewService(uow, payableRepo, accountRepo, journalSvc)

	taxSvc := taxApp.NewService(taxRepo)
	doctorShareSvc := doctorshareApp.NewService(uow, doctorShareRepo, taxRepo, accountingSvc)
	treasurySvc := treasuryApp.NewService(uow, treasuryRepo, accountingSvc)
	assetSvc := assetApp.NewService(uow, assetRepo, accountingSvc)
	prepaidSvc := prepaidApp.NewService(uow, prepaidRepo, accountingSvc)
	budgetSvc := budgetApp.NewService(budgetRepo)
	auditSvc := auditApp.NewService(auditRepo)
	closingSvc := accountingApp.NewClosingService(db, uow, accountingSvc, periodRepo, accountRepo)

	// 4. HTTP Handlers
	accountingHandler := handler.NewAccountingHandler(accountingSvc)
	cashbankHandler := handler.NewCashBankHandler(cashbankSvc, accountingSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	billingHandler := handler.NewBillingHandler(billingSvc)
	receivableHandler := handler.NewReceivableHandler(receivableSvc)
	payableHandler := handler.NewPayableHandler(payableSvc)
	doctorShareHandler := handler.NewDoctorShareHandler(doctorShareSvc)
	treasuryHandler := handler.NewTreasuryHandler(treasurySvc)
	taxHandler := handler.NewTaxHandler(taxSvc)
	assetHandler := handler.NewAssetHandler(assetSvc)
	prepaidHandler := handler.NewPrepaidHandler(prepaidSvc)
	budgetHandler := handler.NewBudgetHandler(budgetSvc)
	closingHandler := handler.NewClosingHandler(closingSvc)
	auditHandler := handler.NewAuditHandler(auditSvc)

	handlerContainer := &routes.HandlerContainer{
		Accounting:  accountingHandler,
		CashBank:    cashbankHandler,
		Reporting:   reportHandler,
		Billing:     billingHandler,
		Receivable:  receivableHandler,
		Payable:     payableHandler,
		DoctorShare: doctorShareHandler,
		Treasury:    treasuryHandler,
		Tax:         taxHandler,
		Asset:       assetHandler,
		Prepaid:     prepaidHandler,
		Budget:      budgetHandler,
		Closing:     closingHandler,
		Audit:       auditHandler,
	}

	return &Container{
		DB:             db,
		Handlers:       handlerContainer,
		AccountingSvc:  accountingSvc,
		JournalSvc:     journalSvc,
		CashBankSvc:    cashbankSvc,
		ReportSvc:      reportSvc,
		ReportingSvc:   reportingSvc,
		BillingSvc:     billingSvc,
		ReceivableSvc:  receivableSvc,
		PayableSvc:     payableSvc,
		DoctorShareSvc: doctorShareSvc,
		TreasurySvc:    treasurySvc,
		TaxSvc:         taxSvc,
		AssetSvc:       assetSvc,
		PrepaidSvc:     prepaidSvc,
		BudgetSvc:      budgetSvc,
		ClosingSvc:     closingSvc,
		AuditSvc:       auditSvc,
	}
}
