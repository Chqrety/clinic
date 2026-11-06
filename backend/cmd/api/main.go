package main

import (
	"log"
	"net/http"
	"os"

	"clinic-api/internal/billing"
	"clinic-api/internal/consultation"
	"clinic-api/internal/domain"
	"clinic-api/internal/pharmacy"
	"clinic-api/internal/queue"
	"clinic-api/internal/satusehat"
	"clinic-api/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "clinic.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[SERVER] Initializing database at %s...", dbPath)
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("[SERVER] Database initialization failed: %v", err)
	}
	defer db.Close()

	// Initialize services
	queueSvc := queue.NewService(db)
	pharmacySvc := pharmacy.NewService(db)
	consultationSvc := consultation.NewService(db)
	billingSvc := billing.NewService(db)
	satusehatSvc := satusehat.NewService(db)

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

	// Mount API v1 routes
	r.Route("/api", func(api chi.Router) {
		queue.RegisterRoutes(api, queueSvc)
		pharmacy.RegisterRoutes(api, pharmacySvc)
		consultation.RegisterRoutes(api, consultationSvc)
		billing.RegisterRoutes(api, billingSvc)
		satusehat.RegisterRoutes(api, satusehatSvc)
	})

	log.Printf("==================================================")
	log.Printf("🚀 Smart Clinic API Server running on port %s", port)
	log.Printf("   URL: http://localhost:%s", port)
	log.Printf("   Health: http://localhost:%s/health", port)
	log.Printf("   API Base: http://localhost:%s/api", port)
	log.Printf("==================================================")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("[SERVER] Server crashed: %v", err)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.WriteHeader(status)
	return nil
}
