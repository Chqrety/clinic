package satusehat_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clinic-api/internal/domain"
	"clinic-api/internal/satusehat"
	"clinic-api/pkg/database"

	"github.com/go-chi/chi/v5"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init in-memory db: %v", err)
	}
	return db
}

func TestSatusehat_StatsAndEncounter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := satusehat.NewService(db)
	r := chi.NewRouter()
	satusehat.RegisterRoutes(r, svc)

	// 1. Test Stats
	req := httptest.NewRequest(http.MethodGet, "/satusehat/stats", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var statsResp domain.ApiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &statsResp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if !statsResp.Success {
		t.Errorf("expected success true")
	}

	// 2. Test Encounter FHIR
	req = httptest.NewRequest(http.MethodGet, "/satusehat/encounter/1", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var encounterResp domain.ApiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &encounterResp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !encounterResp.Success {
		t.Errorf("expected success true")
	}
}

func TestSatusehat_Sync(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := satusehat.NewService(db)
	r := chi.NewRouter()
	satusehat.RegisterRoutes(r, svc)

	// Test Sync with queueId = 1
	req := httptest.NewRequest(http.MethodPost, "/satusehat/sync/1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var syncResp domain.ApiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &syncResp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !syncResp.Success {
		t.Errorf("expected success true, got false")
	}

	// Verify logged in satusehat_logs
	var count int
	_ = db.QueryRow("SELECT count(*) FROM satusehat_logs WHERE entity_id = 1").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 synced log, got %d", count)
	}
}
