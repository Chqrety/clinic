# Comprehensive Codebase & System Audit Report

**Date:** 2026-11-13  
**Auditor:** Senior Infrastructure & Backend Engineer (Ponytail Mode)  
**Target:** Smart Clinic Management System (Frontend & Backend Monorepo)  
**Status:** Audit findings resolved & verified.

---

### What this repo does
This repository runs a multi-tenant clinic management system for outpatient healthcare facilities (Kidinest Pediatric & Semarang Skin Center). It coordinates public patient registration kiosks, TV queue monitors, doctor consultations, e-prescriptions, FEFO pharmacy dispensing, and multi-channel cashier billing (Cash, EDC, QRIS, Bank Transfer). The assumed load is a busy outpatient clinic processing 50–200 visits per day across 6 staff workstations simultaneously, with real-time TV displays and mobile patient queue trackers.

---

### Must fix (Bugs, Security, Data Loss, Concurrency)

#### 1. [RESOLVED] Tripay Webhook Race Condition Overwrites Queue Numbers
- **File:** `backend/internal/queue/usecase/queue_usecase.go` & `backend/internal/queue/repository/postgres_repo.go`
- **Resolution:** Added `AND payment_status != 'PAID'` to atomic update query; duplicate webhook delivery exits gracefully as an idempotent no-op.

#### 2. [RESOLVED] Un-transactional Cashier Invoice Settlement
- **File:** `backend/internal/cashier/usecase/cashier_usecase.go` & `backend/internal/cashier/repository/postgres_repo.go`
- **Resolution:** Created `SettleInvoice(ctx, invoice, queueID)` wrapping invoice insert and queue status update (`PAID` / `COMPLETED`) inside a single atomic database transaction.

#### 3. [RESOLVED] Batch ID Disconnect in Prescription Reserving & Dispensing
- **File:** `backend/internal/consultation/repository/postgres_repo.go` & `backend/internal/pharmacy/repository/postgres_repo.go`
- **Resolution:** Captured selected `batch_id` upon reservation in `prescription_items`, and updated `DispensePrescriptionTx` / `CancelPrescriptionTx` to deduct/release directly using that batch ID, preventing stock desync across batches.

#### 4. [RESOLVED] Hardcoded Fallback JWT Secret in Production
- **File:** `backend/cmd/api/main.go`
- **Resolution:** Added startup check executing `log.Fatalf` if `ENV == "production"` and `JWT_SECRET` is unset or less than 32 characters.

#### 5. [RESOLVED] Dev Seeder Routes Unguarded by Role or Auth
- **File:** `backend/internal/dev/handler.go` & `backend/internal/dev/seeder.go`
- **Resolution:** Dev seeder endpoints require explicit dev mode configuration and are disabled in production mode.

---

### Should fix (Performance, Duplication, Dead Code, Missing Tests)

#### 6. [RESOLVED] SSE Monitor Polling Overhead & Header Cleanup
- **File:** `backend/internal/queue/delivery/http/handler.go`
- **Resolution:** Removed redundant `Access-Control-Allow-Origin: *` header (already handled by router CORS) and added immediate stream exit on client disconnect error.

#### 7. [RESOLVED] Orphan Unused Views in Frontend
- **File:** `frontend/src/views/PatientKiosk.vue` & `frontend/src/views/LandingPage.vue`
- **Resolution:** Both obsolete prototype files deleted (-822 lines).

#### 8. [RESOLVED] Dead Legacy Client Endpoints
- **File:** `frontend/src/api/client.ts`
- **Resolution:** Removed pre-DDD monolithic helper methods pointing to deprecated `/drugs` and `/prescriptions` endpoints.

#### 9. [DEFERRED] Duplicate Legacy SQLite Fallback Services in Backend
- **File:** `backend/internal/queue/queue.go`, `backend/internal/pharmacy/pharmacy.go`, `backend/pkg/database/db.go`
- **Status:** Preserved as offline fallback for standalone SQLite mode until full Docker containerization is standardized.

#### 10. [RESOLVED] Cashier Refund and Voucher Edge Cases
- **File:** `backend/internal/cashier/usecase/cashier_usecase_test.go`
- **Resolution:** Added unit test coverage for 100% voucher discount (Rp 0 net payable).

#### 11. [RESOLVED] Untested SATUSEHAT Sync Payload Error Handling
- **File:** `backend/internal/satusehat/satusehat_test.go`
- **Resolution:** Added comprehensive unit tests for SATUSEHAT encounters, FHIR bundle generation, and sync logs.

---

### Nice to have (Speed, Polish, Shorter Code)

#### 12. [RESOLVED] Redundant CORS Header in SSE Handler
- **Resolution:** Merged into Finding #6.

---

### Summary & Verdict
- **Verdict:** All critical bugs (race conditions, un-transactional invoice settlements, batch ID reservation disconnects, JWT fallback) and frontend dead code have been completely fixed and verified.
- **Verification:**
  - Backend: `cd backend && go test ./...` -> 100% PASS.
  - Frontend: `npm --prefix frontend run build` -> 0 errors (78 kB CSS / 444 kB JS).
