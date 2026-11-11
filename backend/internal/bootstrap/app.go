package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clinic-api/internal/config"
	"clinic-api/internal/transport/http/routes"
)

type Application struct {
	Config    *config.Config
	Logger    *Logger
	Container *Container
	Server    *http.Server
}

func NewApplication() (*Application, error) {
	cfg := config.LoadConfig()
	logger := InitLogger(cfg.App.Name, cfg.App.LogLevel)

	db, err := InitDatabase(cfg.Database)
	if err != nil {
		logger.Printf("⚠️ Warning: PostgreSQL not reachable on %s:%d: %v", cfg.Database.Host, cfg.Database.Port, err)
		logger.Println("💡 Note: Running in offline/disconnected mode until PostgreSQL is started.")
	}

	container := BuildContainer(db)
	router := routes.SetupRouter(container.Handlers)

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Application{
		Config:    cfg,
		Logger:    logger,
		Container: container,
		Server:    server,
	}, nil
}

// Run starts the HTTP server with graceful shutdown on interrupt.
func (a *Application) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		a.Logger.Printf("🚀 %s running on port %d (%s mode)", a.Config.App.Name, a.Config.App.Port, a.Config.App.Env)
		a.Logger.Printf("   Routes: /health, /api/v1/accounts, /api/v1/journals, /api/v1/cash, /api/v1/bank, /api/v1/reports")
		serverErrors <- a.Server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		a.Logger.Printf("🛑 Signal %v received, shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.Server.Shutdown(ctx); err != nil {
			_ = a.Server.Close()
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		if a.Container.DB != nil {
			_ = a.Container.DB.Close()
			a.Logger.Println("🔌 Database connection pool closed.")
		}
	}

	a.Logger.Println("👋 Application stopped gracefully.")
	return nil
}

