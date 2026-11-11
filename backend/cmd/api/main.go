package main

import (
	"log"
	"net/http"
	"os"
	"time"

	authHttp "clinic-api/internal/auth/delivery/http"
	authRepo "clinic-api/internal/auth/repository"
	authUsecase "clinic-api/internal/auth/usecase"
	"clinic-api/internal/billing"
	cashierHttp "clinic-api/internal/cashier/delivery/http"
	cashierRepo "clinic-api/internal/cashier/repository"
	cashierUsecase "clinic-api/internal/cashier/usecase"
	legacyConsultation "clinic-api/internal/consultation"
	consHttp "clinic-api/internal/consultation/delivery/http"
	consRepo "clinic-api/internal/consultation/repository"
	consUsecase "clinic-api/internal/consultation/usecase"
	"clinic-api/internal/domain"
	legacyPharmacy "clinic-api/internal/pharmacy"
	pharmacyHttp "clinic-api/internal/pharmacy/delivery/http"
	pharmacyRepo "clinic-api/internal/pharmacy/repository"
	pharmacyUsecase "clinic-api/internal/pharmacy/usecase"
	legacyQueue "clinic-api/internal/queue"
	queueHttp "clinic-api/internal/queue/delivery/http"
	queueRepo "clinic-api/internal/queue/repository"
	queueUsecase "clinic-api/internal/queue/usecase"
	"clinic-api/internal/satusehat"
	sharedDb "clinic-api/internal/shared/database"
	sharedMw "clinic-api/internal/shared/middleware"
	staffHttp "clinic-api/internal/staff/delivery/http"
	staffRepo "clinic-api/internal/staff/repository"
	staffUsecase "clinic-api/internal/staff/usecase"
	legacyDb "clinic-api/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "clinic-dev-secret-key-at-least-32-chars-long"
	}

	// 1. Inisialisasi PostgreSQL (Database Utama SCMS)
	pgCfg := sharedDb.LoadConfig()
	log.Printf("[SERVER] Menghubungkan ke PostgreSQL di %s:%s (DB: %s)...", pgCfg.Host, pgCfg.Port, pgCfg.DBName)
	pgDB, pgErr := sharedDb.InitPostgres(pgCfg)
	if pgErr != nil {
		log.Printf("[SERVER] Peringatan: PostgreSQL belum aktif (%v). Jalankan 'podman compose up -d' untuk modul staf & auth.", pgErr)
	} else {
		log.Printf("[SERVER] PostgreSQL berhasil terhubung & skema staff_users termigrasi.")
		defer pgDB.Close()
	}

	// 2. Inisialisasi fallback database untuk modul legacy yang belum dimigrasi
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "clinic.db"
	}
	legacySqlite, err := legacyDb.InitDB(dbPath)
	if err != nil {
		log.Printf("[SERVER] Peringatan: Legacy SQLite db tidak dapat diakses: %v", err)
	} else {
		defer legacySqlite.Close()
	}

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = jsonResponse(w, http.StatusOK, domain.ApiResponse{
			Success: true,
			Message: "Smart Clinic API is healthy and operational",
		})
	})

	// Auth & Staff Middlewares
	authMiddleware := sharedMw.Authenticate(jwtSecret)
	ownerMiddleware := sharedMw.RequireRole("SUPERADMIN_OWNER")

	// 3. Register Bounded Contexts: Auth & Staff (DDD) jika PG aktif
	if pgDB != nil {
		authRepository := authRepo.NewPostgresAuthRepository(pgDB)
		authUc := authUsecase.NewAuthUsecase(authRepository, jwtSecret, 24*time.Hour)
		authHandler := authHttp.NewAuthHandler(authUc)
		authHttp.RegisterAuthRoutes(r, authHandler, authMiddleware)

		staffRepository := staffRepo.NewPostgresStaffRepository(pgDB)
		staffUc := staffUsecase.NewStaffUsecase(staffRepository)
		staffHandler := staffHttp.NewStaffHandler(staffUc)
		staffHttp.RegisterStaffRoutes(r, staffHandler, authMiddleware, ownerMiddleware)

		// Bounded Context: Queue & Patient Portal (DDD)
		queueRepository := queueRepo.NewPostgresQueueRepository(pgDB)
		waAdapter := queueUsecase.NewMockWhatsAppAdapter()
		queueUc := queueUsecase.NewQueueUsecase(queueRepository, queueRepository, waAdapter, "http://localhost:5173")
		queueHandler := queueHttp.NewQueueHandler(queueUc)
		queueHttp.RegisterPublicQueueRoutes(r, queueHandler)

		// Bounded Context: Cashier & Frontdesk GRO (DDD)
		cashierRepository := cashierRepo.NewPostgresCashierRepository(pgDB)
		cashierUc := cashierUsecase.NewCashierUsecase(cashierRepository, cashierRepository, cashierRepository, queueRepository, waAdapter, "http://localhost:5173")
		cashierHandler := cashierHttp.NewCashierHandler(cashierUc)
		cashierHttp.RegisterCashierRoutes(r, cashierHandler, authMiddleware)

		// Bounded Context: Doctor Consultation & Fast E-Prescription (DDD)
		consRepository := consRepo.NewPostgresConsultationRepository(pgDB)
		consUc := consUsecase.NewConsultationUsecase(consRepository)
		consHandler := consHttp.NewConsultationHandler(consUc)
		consHttp.RegisterConsultationRoutes(r, consHandler, authMiddleware)

		// Bounded Context: Pharmacy & Stock Opname (DDD)
		pharmRepository := pharmacyRepo.NewPostgresPharmacyRepository(pgDB)
		pharmUc := pharmacyUsecase.NewPharmacyUsecase(pharmRepository)
		pharmHandler := pharmacyHttp.NewPharmacyHandler(pharmUc)
		pharmacyHttp.RegisterPharmacyRoutes(r, pharmHandler, authMiddleware)
	}

	// 4. Register Legacy Modules (Fase sebelum Bagian 3-5 dimigrasikan)
	if legacySqlite != nil {
		queueSvc := legacyQueue.NewService(legacySqlite)
		pharmacySvc := legacyPharmacy.NewService(legacySqlite)
		consultationSvc := legacyConsultation.NewService(legacySqlite)
		billingSvc := billing.NewService(legacySqlite)
		satusehatSvc := satusehat.NewService(legacySqlite)

		r.Route("/api", func(api chi.Router) {
			legacyQueue.RegisterRoutes(api, queueSvc)
			legacyPharmacy.RegisterRoutes(api, pharmacySvc)
			legacyConsultation.RegisterRoutes(api, consultationSvc)
			billing.RegisterRoutes(api, billingSvc)
			satusehat.RegisterRoutes(api, satusehatSvc)
		})
	}

	log.Printf("==================================================")
	log.Printf("🚀 Smart Clinic API Server running on port %s", port)
	log.Printf("   URL: http://localhost:%s", port)
	log.Printf("   Health: http://localhost:%s/health", port)
	log.Printf("   Auth: http://localhost:%s/api/auth/login", port)
	log.Printf("   Staff Admin: http://localhost:%s/api/admin/staff", port)
	log.Printf("==================================================")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("[SERVER] Server crashed: %v", err)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	return nil
}
