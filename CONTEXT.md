# System Context & Architecture

## 1. Project Overview
Aplikasi **Smart Clinic Management System (SCMS)** adalah solusi digitalisasi operasional kustom terpadu untuk klinik modern, dirancang secara khusus untuk melayani dua unit bisnis:
1. **Kidinest:** Layanan poli anak, vaksinasi, dan tumbuh kembang dengan antrian berkeadilan (FIFO) dan pembayaran biaya registrasi terverifikasi di awal.
2. **Semarang Skin Center (SSC):** Layanan estetika medis, perawatan wajah (*facial*), dan tindakan kulit dengan skema reservasi jadwal (*booking* tanpa bayar di awal / Rp0) dan pelunasan di kasir meja depan setelah tindakan selesai.

Dibangun menggunakan arsitektur **Modular Monolith** dengan backend **Go (Golang)** dan frontend **Vue 3 SPA**, berfokus murni pada **keunggulan operasional medis (*Clinical Operations Excellence*)**.

> **Prinsip Utama Domain:**
> * **Privacy-First:** Pendaftaran pasien tidak meminta Nomor KTP/NIK. Identitas di monitor antrian publik disamarkan otomatis (`Budi P****`).
> * **Paperless Queue:** Tiket antrian dikirimkan via WhatsApp dan dipantau live melalui web. Pasien *walk-in* tetap didaftarkan lewat web oleh GRO meja depan tanpa struk fisik.
> * **Invarian Stok Obat:** `available_stock = physical_stock - reserved_stock`. Stok otomatis terkunci (*reserved*) saat dokter meresepkan dan terpotong permanen saat diserahkan apoteker.
> * **Bebas Embel-Embel Akuntansi & SATUSEHAT:** Modul akuntansi umum (jurnal/pajak) dan integrasi bridging Kemenkes SATUSEHAT resmi berada di luar lingkup rilis ini (*Out of Scope*).

---

## 2. Technology Stack

* **Backend:**
  * Runtime: Go 1.22+
  * Framework/Router: Go Chi / Standard Library `net/http`
  * Database: PostgreSQL 16+ (ACID compliant, row-level locking `SELECT ... FOR UPDATE` untuk integritas transaksi antrian & reservasi stok obat, native concurrent connection pool)
  * Data Layer: `database/sql` + `github.com/jackc/pgx/v5/stdlib` (atau `github.com/lib/pq`)
  * Containerization (Local Dev): Podman / Docker via `docker-compose.yml` (`podman compose up -d`)
* **Frontend:**
  * Framework: Vue 3 (Composition API, `<script setup>`)
  * Bundler: Vite
  * Styling: Tailwind CSS (estetika Brandlos / Nordic Modern, warm off-white canvas, pill badges)
  * State Store: Pinia
  * UI Icons: Lucide Icons
* **Real-Time Communication:** Server-Sent Events (SSE) / WebSocket untuk sinkronisasi monitor antrian live dan notifikasi audio chime apotek (*zero page refresh*).
* **Integrasi Eksternal:** WhatsApp Gateway Adapter (tiket & link antrian) & Payment Gateway Adapter (QRIS / transfer verifikasi otomatis).

---

## 3. Core Domain Entities & Relationships

### A. Patient & Queue
* **Patient:** `id`, `name`, `phone`, `created_at` *(Privacy-first, tanpa NIK)*
* **Queue:** 
  * `id`, `patient_id`, `clinic_unit` (`KIDINEST`, `SSC`)
  * `category` (`A_VAKSIN`, `B_SAKIT`, `C_TUMBANG`, `D_PAKET`, `E_WALKIN`, `SSC_TREATMENT`)
  * `doctor_id`, `doctor_name`
  * `queue_number` (Format: `A-001`, `B-005`, `E-012`, `S-003`)
  * `registration_fee` (Kidinest: nominal registrasi; SSC: Rp 0)
  * `payment_method` (`QRIS`, `TRANSFER`, `CASH`, `UNPAID_BOOKING`)
  * `payment_status` (`PENDING`, `PAID`, `REFUNDED`)
  * `queue_status` (`WAITING`, `CALLING`, `IN_EXAM`, `PHARMACY_WAITING`, `COMPLETED`, `CANCELLED`)
  * `created_at`, `called_at`, `completed_at`

### B. Staff & Security (Role-Based Access)
* **StaffUser:**
  * `id`, `employee_number` (Nomor Induk Karyawan resmi klinik)
  * `name`, `role` (`SUPERADMIN_OWNER`, `DOCTOR`, `GRO_CASHIER`, `PHARMACIST`)
  * `clinic_unit` (`KIDINEST`, `SSC`, `ALL`)
  * `password_hash`, `is_active`, `created_at`
  *(Akun staf diinput & dikelola langsung oleh Atasan/Owner)*

### C. Consultation & E-Prescription
* **Consultation:** 
  * `id`, `queue_id`, `doctor_id`, `doctor_name`
  * `anamnesis`, `diagnosis`, `action_notes`
  * `action_fee` (Jasa tindakan dokter/terapis: Kidees vaksin/tindakan; SSC facial/peeling/laser)
  * `status` (`DRAFT`, `COMPLETED`)
* **Prescription:** 
  * `id`, `consultation_id`, `status` (`RESERVED`, `PREPARING`, `DISPENSED`, `CANCELLED`)
* **PrescriptionItem:** 
  * `id`, `prescription_id`
  * `drug_id` (nullable untuk obat racikan/item manual)
  * `item_name`, `is_manual` (boolean)
  * `unit_price`, `quantity`, `subtotal`
  * `dosage_instruction` (Signa / aturan pakai)

### D. Pharmacy & Stock Opname
* **Drug:** 
  * `id`, `clinic_unit` (`KIDINEST`, `SSC`, `ALL`)
  * `name`, `unit` (`Tablet`, `Botol`, `Tube`, `Ampul`, `Strip`)
  * `buy_price`, `sell_price`, `min_stock_alert`
* **DrugBatch:** 
  * `id`, `drug_id`, `batch_number`, `expired_date`
  * `physical_stock`, `reserved_stock`
  * Formula Invarian: `available_stock = physical_stock - reserved_stock`
* **StockOpnameLog:** 
  * `id`, `drug_id`, `batch_id`, `system_stock`, `actual_stock`, `discrepancy`, `adjusted_by`, `notes`, `created_at`
  *(Mendukung ekspor & impor massal via file Excel/CSV)*

### E. Billing & Kasir Meja Depan (GRO)
* **Invoice:** 
  * `id`, `queue_id`, `patient_id`, `invoice_number`
  * `subtotal_action`, `subtotal_medication`, `subtotal_registration`
  * `voucher_discount` (input nominal potongan voucher)
  * `total_amount` (setelah diskon)
  * `payment_method` (`CASH`, `BANK_TRANSFER`, `EDC`)
  * `cashier_staff_id`, `status` (`PAID`, `REFUNDED`)
  * `refund_bank_account`, `refund_notes`
  * `created_at`
* **CashierShiftReport:**
  * `id`, `staff_id`, `shift_date`, `total_cash`, `total_transfer`, `total_voucher`, `total_refund`, `sent_to_owner_at`

---

## 4. State Transitions & Lifecycle

### A. Alur Pendaftaran & Antrian (FIFO)
```text
[Pasien Scan QR / Dibantu GRO]
       │
       ├─ Kidinest (Bayar Pendaftaran Awal)
       │       ├─ (QRIS / Transfer) ─> [Auto Verifikasi Lunas] ──┐
       │       └─ (Tunai di Kasir) ──> [Kasir Terima Uang] ─────┼─> [Nomor Antrian Aktif]
       │                                                         │   (Kategori A/B/C/D/E)
       └─ Semarang Skin Center (Booking Saja)                    │   (Masuk Antrian Dokter FIFO)
               └─ [Tanpa Bayar Awal / Rp0] ──────────────────────┘
```

### B. Alur E-Resep & Reserved Stock Apotek
```text
[Dokter Rangkai Resep di Tabel]
  (Dropdown Search + Harga Satuan + Opsi Manual Racikan)
       │
       ▼
[Dokter Klik "Kirim Resep"] ────> [Audio Chime Berbunyi di Komputer Apotek]
       │                          [Stok Obat Masuk Status RESERVED]
       │
       ├─ (Apoteker Serahkan Obat) ─> [Physical Stock Terpotong, Reserved Rilis] ─> [Selesai]
       │
       └─ (Resep Dibatalkan / Edit) ─> [Reserved Stock Dirilis Kembali]
```

### C. Alur Kasir & Cetak Invoice
```text
[Pasien Selesai Konsultasi / Tindakan]
       │
       ▼
[Kasir Panggil Tagihan Pasien]
  (Total: Jasa Tindakan + Obat/Vaksin - Diskon Voucher)
       │
       ├─ [Bayar: Tunai / Transfer / EDC] ──> [Cetak Struk Thermal 80mm / PDF]
       │
       └─ [Batal Kunjungan Berbayar] ───────> [Catat Norek Refund Manual]
```

---

## 5. Directory Layout & Architecture Rules

### A. Architectural Principles
1. **Backend: Domain-Driven Design (DDD) per Bounded Context**
   * Setiap modul bisnis dikelompokkan ke dalam bounded context di `backend/internal/<context>/`.
   * Setiap context menerapkan pemisahan layer DDD yang bersih dan terisolasi:
     * `domain/`: Pure business entity, value objects, domain error, dan interface repository (bebas dependensi framework/SQL).
     * `usecase/`: Application logic / use cases yang mengorkestrasi alur bisnis domain.
     * `repository/`: Implementasi persistensi data ke PostgreSQL (SQL query, database mapping).
     * `delivery/http/`: HTTP handler (Go Chi), request DTO validation, JSON response, dan route registration.
   * `internal/shared/`: Shared kernel untuk komponen lintas domain (database connection pool, HTTP response wrapper, auth middleware/token).
2. **Frontend: Role/Feature Views + Cohesive Components**
   * Halaman diatur berdasarkan Role pengguna klinik (`src/views/`).
   * Gunakan `src/components/ui/` hanya untuk komponen yang benar-benar digunakan berulang.
3. **Frontend Layout: Prioritas Sidebar untuk Admin & Staf**
   * Antarmuka kerja staf dan admin (Atasan/Owner, Dokter, Kasir GRO, Apotek, Gudang) diprioritaskan menggunakan **Sidebar Navigation (Navigasi Samping Vertikal)** bergaya Nordic Modern yang ringkas dan fungsional, menggantikan tab atas, guna memaksimalkan keleluasaan layar untuk tabel resep, monitor antrian, dan form transaksi.
   * *(Catatan: Pembaruan UI sidebar ini dijadwalkan setelah seluruh modul rencana fungsional inti 1–5 selesai dibangun).*

### B. Tree Layout

```text
clinic/
├── docker-compose.yml          # PostgreSQL 16 local container (Podman / Docker)
├── .env.example                # Template konfigurasi DB & Port
├── docs/
│   ├── PRD.md
│   ├── PROPOSAL_OPERASIONAL_KLINIK.md
│   ├── PLAN_BAGIAN_1_STAFF_AUTH.md
│   ├── PLAN_BAGIAN_2_PATIENT_QUEUE.md
│   ├── clinic-flowchart-spec.excalidraw
│   └── ...
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go         # Dependency injection & router bootstrap
│   ├── internal/
│   │   ├── shared/             # Shared kernel lintas domain
│   │   │   ├── database/       # Postgres connection & migration runner
│   │   │   ├── response/       # Standard JSON response helper
│   │   │   └── middleware/     # Auth JWT & CORS middleware
│   │   ├── auth/               # Bounded Context: Autentikasi Staf & NIK Karyawan
│   │   │   ├── domain/         # StaffAuth entity, Token claims, Repo interface
│   │   │   ├── usecase/        # Login & verify logic
│   │   │   ├── repository/     # Postgres query auth
│   │   │   └── delivery/http/  # Login handler & auth routes
│   │   ├── staff/              # Bounded Context: Manajemen Staf (Atasan/Owner)
│   │   │   ├── domain/         # Staff entity, Role/Unit value types, Repo interface
│   │   │   ├── usecase/        # CreateStaff, UpdateStaff, ResetPassword
│   │   │   ├── repository/     # Postgres CRUD staff
│   │   │   └── delivery/http/  # Staff admin handler & routes
│   │   ├── queue/              # Bounded Context: Pendaftaran & Antrian FIFO
│   │   │   ├── domain/         # Patient, Queue, DailyCounter entity & repo interface
│   │   │   ├── usecase/        # RegisterQueue, VerifyPayment, NextQueue
│   │   │   ├── repository/     # Postgres FIFO locking & queries
│   │   │   └── delivery/http/  # Public registration, SSE stream, TV monitor handler
│   │   ├── consultation/       # Bounded Context: Konsultasi Dokter & E-Resep
│   │   │   ├── domain/         # Consultation, Prescription, Item entity
│   │   │   ├── usecase/        # SubmitConsultation, ReserveMedication
│   │   │   ├── repository/     # Postgres consultation queries
│   │   │   └── delivery/http/  # Doctor room handler & routes
│   │   ├── pharmacy/           # Bounded Context: Farmasi & Gudang Obat
│   │   │   ├── domain/         # Drug, Batch, StockOpname entity
│   │   │   ├── usecase/        # DispenseDrug, StockOpnameImportExport
│   │   │   ├── repository/     # Postgres drug batch & atomic stock queries
│   │   │   └── delivery/http/  # Pharmacy room & inventory handlers
│   │   └── cashier/            # Bounded Context: Kasir GRO & Tagihan
│   │       ├── domain/         # Invoice, Payment, ShiftReport entity
│   │       ├── usecase/        # ProcessPayment, CloseShiftReport, Refund
│   │       ├── repository/     # Postgres invoice & payment queries
│   │       └── delivery/http/  # Cashier frontdesk handler & thermal receipt
├── frontend/
│   ├── src/
│   │   ├── views/              # Halaman berbasis Peran Pengguna
│   │   │   ├── StaffLogin.vue    # Login NIK Karyawan & password
│   │   │   ├── StaffAdmin.vue    # Layar Atasan: input data staf & kontrol password
│   │   │   ├── PatientPortal.vue # Landing page QR (Kidinest & SSC) + Tiket
│   │   │   ├── QueueMonitor.vue  # Layar TV monitor antrian publik (disamarkan)
│   │   │   ├── CashierView.vue   # Meja depan GRO: walk-in, bayar, struk, tutup shift
│   │   │   ├── DoctorRoom.vue    # Layar periksa dokter: antrian FIFO & e-resep tabel
│   │   │   ├── PharmacyRoom.vue  # Layar apotek: resep live chime & dispensing
│   │   │   └── InventoryView.vue # Gudang obat: master obat & stok opnam Excel
│   │   ├── components/
│   │   │   ├── ui/               # PillBadge.vue, ModalDialog.vue, StatCard.vue
│   │   │   └── RoleSwitcher.vue  # Floating bar demo switcher
│   │   ├── stores/               # Pinia global states (auth, queue, pharmacy)
│   │   └── api/                  # Fetch client helper
│   ├── index.html
│   └── vite.config.ts
├── CONTEXT.md
├── AGENT.md
├── DESIGN.md
└── GEMINI.md
├── AGENT.md
├── DESIGN.md
└── GEMINI.md
```
