# System Context & Architecture

## 1. Project Overview
Aplikasi **Smart Clinic Management System** adalah solusi digital untuk klinik swasta modern yang mengintegrasikan pendaftaran mandiri (self-service queue), e-resep dokter, manajemen inventori apotek berbasis batch & expired date, single-invoice pasien dengan audit internal, serta kesiapan integrasi **SATUSEHAT Kemenkes RI**.

Dibangun menggunakan arsitektur **Modular Monolith** dengan backend **Go** dan frontend **Vue 3 SPA**.

---

## 2. Technology Stack

* **Backend:**
  * Runtime: Go 1.22+
  * Framework/Router: Go Chi / Fiber (REST API)
  * Database: SQLite (local development & demo zero-config, portable ke PostgreSQL)
  * ORM / Data Layer: GORM atau modernc.org/sqlite + database/sql
* **Frontend:**
  * Framework: Vue 3 (Composition API, `<script setup>`)
  * Bundler: Vite
  * Styling: Tailwind CSS
  * State Store: Pinia
  * UI Helper: Lucide Icons, Headless UI / Tailwind UI
* **Communication:** JSON REST API + Server-Sent Events (SSE) / Polling ringan untuk live updates antrian & notifikasi apotek.

---

## 3. Core Domain Entities & Relationships

### A. Patient & Queue
* **Patient:** `id`, `nik`, `name`, `phone`, `birth_date`, `ihs_number` (SATUSEHAT ID)
* **Queue:** `id`, `patient_id`, `queue_number` (A-001), `payment_method` (QRIS, CASH), `payment_status` (PENDING, PAID), `queue_status` (WAITING, CALLING, IN_EXAM, COMPLETED, CANCELLED), `created_at`

### B. Consultation & Prescription
* **Consultation:** `id`, `queue_id`, `doctor_name`, `anamnesis`, `diagnosis`, `icd10_code`, `status` (DRAFT, COMPLETED)
* **Prescription:** `id`, `consultation_id`, `status` (RESERVED, PREPARING, DISPENSED, CANCELLED)
* **PrescriptionItem:** `id`, `prescription_id`, `drug_id`, `quantity`, `dosage_instruction`, `price`

### C. Pharmacy & Inventory
* **Drug:** `id`, `name`, `kfa_code` (Kode Farmasi Kemenkes), `unit` (Tablet, Botol, Strip), `min_stock_alert`, `sell_price`
* **DrugBatch:** `id`, `drug_id`, `batch_number`, `expired_date`, `physical_stock`, `reserved_stock`, `buy_price`

> **Formula Stok Obat:**  
> `available_stock = physical_stock - reserved_stock`

### D. Billing & Invoice
* **Invoice:** `id`, `queue_id`, `patient_id`, `invoice_number`, `total_amount`, `payment_method`, `status` (PAID)
* **InvoiceItem:** `id`, `invoice_id`, `item_type` (REGISTRATION, CONSULTATION, MEDICATION, ACTION), `item_name`, `amount`, `doctor_share`, `clinic_share`, `tax_share`

### E. SATUSEHAT Integration Log
* **SatusehatLog:** `id`, `resource_type` (Patient, Encounter, Condition, MedicationRequest), `entity_id`, `fhir_payload` (JSON), `status` (READY, SYNCED)

---

## 4. State Transitions & Lifecycle

### Antrian & Pembayaran (FIFO)
```text
[Pasien Input Data]
       │
       ├─ (Pilih QRIS) ──> [Simulasi Bayar Lunas] ──> [Nomor Antrian Aktif] ──> [Masuk Urutan FIFO]
       │
       └─ (Pilih Tunai) ─> [Status Menunggu Kasir] ──> [Kasir Konfirmasi Tunai] ──> [Nomor Antrian Aktif]
```

### Siklus E-Resep & Reserved Stock
```text
[Dokter Rangkai Resep] 
       │
       ▼ (Validasi Available Stock > 0)
[Dokter Kirim Resep] ──> [Reserved Stock +Qty] ──> [Notif Masuk ke Apotek]
       │
       ├─ (Apoteker Serahkan Obat) ─> [Physical Stock -Qty] & [Reserved Stock -Qty] ─> (Selesai)
       │
       └─ (Resep Dibatalkan) ───────> [Reserved Stock -Qty] ─> (Batal, Fisik Utuh)
```

---

## 5. Directory Layout & Architecture Rules

### A. Architectural Principles
1. **Backend: Package by Feature (Domain-Driven Package)**
   * DILARANG menggunakan *Global Layered Architecture* (seperti `pkg/controllers/`, `pkg/services/`, `pkg/repositories/` terpisah di tingkat atas) karena anti-pattern di Go dan memicu *circular dependency*.
   * Setiap modul domain (`queue`, `pharmacy`, dll.) memiliki `handler.go`, `service.go`, dan `repository.go` sendiri di dalam foldernya.
2. **Frontend: Role/Feature Views + Cohesive UI Components**
   * DILARANG menggunakan *Atomic Design* (`atoms/`, `molecules/`, `organisms/`) untuk mencegah *wrapper hell* dan debat kategorisasi tak berujung.
   * Gunakan `src/views/` berbasis Role pengguna klinik dan `src/components/ui/` hanya untuk komponen yang benar-benar dipakai berulang.

### B. Tree Layout

```text
clinic/
├── docs/
│   ├── PRD.md
│   └── ...
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── domain/           # Shared models & interfaces
│   │   ├── queue/            # handler.go, service.go, repo.go (antrian & payment)
│   │   ├── pharmacy/         # handler.go, service.go, repo.go (obat, batch, reserved stock)
│   │   ├── consultation/     # handler.go, service.go, repo.go (rekam medis & e-resep)
│   │   ├── billing/          # handler.go, service.go, repo.go (invoice & breakdown)
│   │   └── satusehat/        # handler.go, service.go (FHIR R4 converter)
│   └── pkg/database/         # SQLite connection & migration
├── frontend/
│   ├── src/
│   │   ├── views/            # Halaman berbasis Role Pengguna
│   │   │   ├── PatientKiosk.vue   # Ambil antrian & bayar QRIS/Cash
│   │   │   ├── CashierView.vue    # Layar kasir (konfirmasi cash & monitor)
│   │   │   ├── DoctorRoom.vue     # Layar periksa & e-resep
│   │   │   ├── PharmacyRoom.vue   # Layar racik obat & potong stok
│   │   │   ├── InventoryView.vue  # Layar stok batch & ED alert
│   │   │   └── BillingView.vue    # Layar invoice pasien & audit internal
│   │   ├── components/
│   │   │   ├── ui/                # PillBadge.vue, ModalDialog.vue, StatCard.vue
│   │   │   └── RoleSwitcher.vue   # Bar floating navigasi demo
│   │   ├── stores/                # Pinia global states (queue.ts, pharmacy.ts)
│   │   └── api/                   # Fetch client helper
│   ├── index.html
│   └── vite.config.ts
├── CONTEXT.md
├── AGENT.md
├── DESIGN.md
└── GEMINI.md
```
