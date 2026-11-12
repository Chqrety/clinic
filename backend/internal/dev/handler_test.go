package dev

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRegisterDevRoutes_ProductionGuard(t *testing.T) {
	// Set ENV to production
	origEnv := os.Getenv("ENV")
	defer os.Setenv("ENV", origEnv)

	os.Setenv("ENV", "production")

	r := chi.NewRouter()
	handler := NewDevHandler(nil)
	RegisterDevRoutes(r, handler)

	req := httptest.NewRequest("POST", "/api/dev/seeds/reset", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// In production, the route MUST NOT exist and return 404
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 in production, got %d", rec.Code)
	}
}

func TestRegisterDevRoutes_DevelopmentMount(t *testing.T) {
	origEnv := os.Getenv("ENV")
	defer os.Setenv("ENV", origEnv)

	os.Setenv("ENV", "development")

	r := chi.NewRouter()
	handler := NewDevHandler(nil)
	RegisterDevRoutes(r, handler)

	// Request an endpoint without db; it reaches the route
	req := httptest.NewRequest("POST", "/api/dev/seeds/reset", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Since db is nil, it executes and returns 500 error (or panic if unhandled), NOT 404
	if rec.Code == http.StatusNotFound {
		t.Fatalf("expected route to be mounted in development, but got 404")
	}
}
