# Rencana Implementasi: Bagian 5 — Farmasi Apotek, Manajemen Stok Invarian & Stok Opnam Excel (Pharmacy & Stock Opname)

Dokumen ini adalah spesifikasi teknis dan rencana eksekusi modul Apotek Farmasi (*Pharmacy Dispensing Room*), Manajemen Stok Obat Terkunci (*Reserved Stock Shadow Deduction*), serta Rekonsiliasi Audit Stok Opnam Excel/CSV untuk unit **Kidinest** dan **Semarang Skin Center (SSC)** sesuai arsitektur DDD (Domain-Driven Design) di Go dan estetika Brandlos/Nordic di Vue 3.

---

## 1. Tujuan & Ruang Lingkup

### A. Sasaran Fungsional (Sesuai Proposal & PRD)
1. **Layar Antrian E-Resep Masuk Real-Time (Fitur 4.1):**
   - Apoteker menerima resep elektronik langsung dari ruang periksa dokter anak (Kidinest) dan dokter kulit (SSC) tanpa lembar kertas fisik.
   - Alur status transparan:
     - `RESERVED` (Resep Baru Masuk / Menunggu Racik).
     - `PREPARING` (Sedang Diracik / Disiapkan).
     - `DISPENSED` (Selesai Diracik & Obat Diserahkan ke Pasien).
   - Notifikasi suara/audio chime visual saat ada resep baru dari dokter.
2. **Penguncian & Pemotongan Stok Invarian Otomatis (Fitur 4.2):**
   - Formula Invarian Anti-Selisih: `available_stock = physical_stock - reserved_stock`.
   - Saat resep terbit: `reserved_stock` bertambah (terkunci agar tidak diserobot dokter lain).
   - Saat apoteker menekan **"Selesai Racik / Serahkan Obat"** (`DISPENSED`):
     - `physical_stock` terpotong permanen sesuai jumlah obat.
     - `reserved_stock` dilepaskan (`reserved_stock -= quantity`).
     - Status antrian pasien beralih ke meja kasir untuk penyelesaian single invoice.
   - Jika resep dibatalkan (`CANCELLED`):
     - `reserved_stock` dilepas kembali tanpa memotong `physical_stock`.
3. **Master Data Obat & Manajemen Batch FEFO (Fitur 4.3):**
   - Pencatatan master obat: nama, unit kemasan (botol, strip, tablet, ampul, tube), kategori, harga beli, harga jual, dan batas peringatan stok minimum (*Min Stock Alert*).
   - Manajemen Batch & FEFO (*First Expired, First Out*): nomor batch dan tanggal kadaluarsa (*Expired Date*).
   - Indikator peringatan visual:
     - 🔴 **Merah:** ED < 3 bulan (segera retur/habiskan).
     - 🟡 **Kuning:** Stok total < batas minimum (*Low Stock Alert*).
     - 🟢 **Hijau:** Stok dan masa berlaku aman.
   - Form Restock Batch Baru: input nomor batch faktur, tanggal ED, jumlah unit masuk, dan harga beli.
4. **Fitur Ekspor Saldo Stok ke Excel/CSV (Fitur 4.4):**
   - 1-klik download format spreadsheet rekonsiliasi lemari apotek: ID obat, nama obat, unit, batch, tanggal ED, stok sistem komputer, kolom kosong hitung fisik nyata, dan catatan pemeriksa.
5. **Fitur Impor Hasil Stok Opnam dari Excel/CSV (Fitur 4.5):**
   - Apoteker mengunggah kembali file spreadsheet hasil hitung fisik lemari apotek.
   - Sistem memvalidasi, menghitung selisih (`discrepancy = actual_stock - system_stock`), mencatat log audit ke tabel `stock_opname_logs`, dan memperbarui saldo fisik secara massal.
6. **Modul Migrasi & Cadangan Data Sistem Lama (Fitur 4.6):**
   - Fasilitas cadangan / ekspor riwayat mutasi obat dan kartu stok klinik ke file CSV terpadu.

### B. Batasan (Out of Scope untuk Bagian 5)
- Tidak ada modul pembelian/Purchase Order (PO) multi-tahap dengan vendor distributor PBF (pencatatan restock fokus pada penerimaan batch langsung di apotek tunggal).
- Tidak ada bridging live Kamus Farmasi Alkes (KFA) Kemenkes SATUSEHAT (hanya mapping kode KFA master obat).

---

## 2. Struktur Data & Migrasi PostgreSQL

```sql
-- 1. Master Obat (Drugs) - Sudah dimigrasikan di Plan 4
CREATE TABLE IF NOT EXISTS drugs (
    id SERIAL PRIMARY KEY,
    clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC', 'ALL')),
    name VARCHAR(150) NOT NULL,
    kfa_code VARCHAR(50),
    unit VARCHAR(30) NOT NULL,
    category VARCHAR(50) NOT NULL,
    min_stock_alert INTEGER NOT NULL DEFAULT 10,
    buy_price INTEGER NOT NULL DEFAULT 0,
    sell_price INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Batch & Stok Obat (Drug Batches) - FEFO Tracking
CREATE TABLE IF NOT EXISTS drug_batches (
    id SERIAL PRIMARY KEY,
    drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE CASCADE,
    batch_number VARCHAR(50) NOT NULL,
    expired_date DATE NOT NULL,
    physical_stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Log Audit Stok Opnam & Rekonsiliasi Excel (Stock Opname Logs)
CREATE TABLE IF NOT EXISTS stock_opname_logs (
    id SERIAL PRIMARY KEY,
    drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE RESTRICT,
    batch_id INTEGER REFERENCES drug_batches(id) ON DELETE SET NULL,
    system_stock INTEGER NOT NULL,
    actual_stock INTEGER NOT NULL,
    discrepancy INTEGER NOT NULL, -- actual - system
    adjusted_by VARCHAR(150) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_stock_opname_drug_id ON stock_opname_logs(drug_id);
CREATE INDEX IF NOT EXISTS idx_stock_opname_created_at ON stock_opname_logs(created_at);
```

---

## 3. Desain Arsitektur Backend DDD (`internal/pharmacy/`)

```text
backend/internal/pharmacy/
├── domain/
│   ├── entity.go             # Drug, Batch, PrescriptionQueue, StockOpnameLog, Repository Interfaces
│   └── errors.go             # Domain errors (ErrBatchExpired, ErrInvalidOpnameQty, dll)
├── repository/
│   └── postgres_repo.go      # Persistensi PostgreSQL (Query FEFO, mutasi stok fisik, ekspor/impor batch)
├── usecase/
│   ├── pharmacy_usecase.go   # Alur dispense, pembatalan resep, restock, generate CSV opnam, parse CSV impor
│   └── pharmacy_usecase_test.go # Unit testing invarian stok (shadow deduction) & rekonsiliasi opnam
└── delivery/
    └── http/
        ├── handler.go        # Chi HTTP handlers (dispense, cancel, restock, download/upload CSV)
        └── routes.go         # Endpoint /api/pharmacy/...
```

---

## 4. Rincian Endpoint HTTP API

| Method | Endpoint | Deskripsi | Hak Akses |
|---|---|---|---|
| `GET` | `/api/pharmacy/prescriptions` | Mengambil daftar e-resep masuk dokter (`RESERVED`, `PREPARING`) | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `POST` | `/api/pharmacy/prescriptions/{id}/status` | Mengubah status racik (`PREPARING`) | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `POST` | `/api/pharmacy/prescriptions/{id}/dispense` | Selesai racik: potong `physical_stock` permanen, rilis `reserved_stock` | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `POST` | `/api/pharmacy/prescriptions/{id}/cancel` | Batalkan resep: rilis `reserved_stock` tanpa memotong fisik | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `GET` | `/api/pharmacy/drugs` | Mengambil master obat lengkap dengan agregat stok dan indikator FEFO | `PHARMACIST`, `SUPERADMIN_OWNER`, `DOCTOR` |
| `POST` | `/api/pharmacy/drugs` | Menambahkan jenis master obat baru | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `POST` | `/api/pharmacy/drugs/{id}/restock` | Menambah stok masuk batch baru faktur | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `GET` | `/api/pharmacy/opname/export` | Download spreadsheet CSV template stok opnam lemari apotek | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `POST` | `/api/pharmacy/opname/import` | Upload file CSV hasil hitung fisik opnam & rekonsiliasi massal | `PHARMACIST`, `SUPERADMIN_OWNER` |
| `GET` | `/api/pharmacy/opname/logs` | Riwayat audit selisih penyesuaian stok opnam | `PHARMACIST`, `SUPERADMIN_OWNER` |

---

## 5. Desain Antarmuka Frontend (`PharmacyRoom.vue` & `InventoryView.vue`)

1. **Layar Antrian E-Resep Masuk (`PharmacyRoom.vue`):**
   - Kartu antrian live dengan badge unit (Kidinest vs SSC) dan nama dokter perujuk.
   - Rincian item obat paten & racikan manual puyer / sirup lengkap dengan aturan pakai (*Signa*).
   - Tombol 1-klik "Mulai Racik" (`PREPARING`) dan "Selesai Racik & Serahkan Obat" (`DISPENSED`).
   - Tombol "Batalkan Resep" (mengembalikan alokasi stok terkunci ke apotek).
2. **Layar Gudang & Stok Opnam Excel (`InventoryView.vue`):**
   - Statistik Cepat: Total SKU Obat, Nilai Aset Fisik Apotek, Total Obat Kadaluarsa Segera (< 3 bulan), dan Peringatan Stok Kritis.
   - Tabel Master Obat & Batch: Indikator warna kartu FEFO (Merah, Kuning, Hijau).
   - Tombol "Tambah Obat Baru" & Tombol "Restock Batch Faktur Masuk".
   - **Fitur Ekspor & Impor Stok Opnam:**
     - Tombol 1-klik "Unduh Lembar Stok Opnam (CSV/Excel)".
     - Modal Unggah File "Impor Hasil Stok Opnam Fisik" dengan preview perbandingan selisih sebelum konfirmasi simpan.
     - Tabel riwayat log penyesuaian stok opnam (*Audit Trail*).
