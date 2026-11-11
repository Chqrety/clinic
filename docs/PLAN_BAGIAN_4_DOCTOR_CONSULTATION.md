# Rencana Implementasi: Bagian 4 — Ruang Pemeriksaan Dokter & E-Resep Cepat (Doctor Consultation & E-Prescription)

Dokumen ini adalah spesifikasi teknis dan rencana eksekusi modul Ruang Pemeriksaan Dokter (*Doctor Consultation Room*) dan Peresepan Elektronik Cepat (*Fast E-Prescription Engine*) untuk unit **Kidinest** dan **Semarang Skin Center (SSC)** sesuai arsitektur DDD (Domain-Driven Design) di Go dan estetika Brandlos/Nordic di Vue 3.

---

## 1. Tujuan & Ruang Lingkup

### A. Sasaran Fungsional (Sesuai Proposal & PRD)
1. **Antrian Pasien Adil & Pemanggilan FIFO (Fitur 3.1):**
   - Dokter memantau daftar antrian pasien berurutan secara adil (*First-In, First-Out* / FIFO).
   - Pasien hanya muncul bila telah menyelesaikan verifikasi pendaftaran (atau pendaftaran reservasi jadwal SSC).
   - Tombol 1-klik "Panggil Pasien" (*Calling*) yang membunyikan bel antrian live dan mengubah status display monitor TV.
   - Tombol 1-klik "Mulai Pemeriksaan" (*In Exam*).
2. **Filter Antrian Khusus Dokter Jaga (Fitur 3.2):**
   - Dokter yang sedang aktif login (misal `DOC-KID-01` dr. Hendra Pratama, Sp.A atau `DOC-SSC-01` dr. Nadia Kusuma, Sp.KK) otomatis melihat daftar pasien yang memilih dirinya.
   - Terdapat pemilih dokter jaga cepat (*Quick Doctor Switcher*) untuk demonstrasi di depan klien.
3. **Penulisan E-Resep Cepat & Fleksibel (Fitur 3.3):**
   - Tabel resep elektronik interaktif:
     - Pencarian obat instan (*typeahead dropdown search*) berdasarkan katalog aktif klinik.
     - Indikator ketersediaan stok fisik apotek (`physical_stock - reserved_stock = available_stock`).
     - Tampilan harga satuan transparan dan subtotal otomatis.
     - **Opsi Input Manual Diketik (Bebas Racikan):** Dokter dapat mengetik bebas nama obat racikan puyer / sirup custom anak atau salep kombinasi SSC non-katalog beserta harga satuan dan aturan pakainya.
4. **Input Komponen Tindakan Medis Kidinest & SSC (Fitur 3.4 & 3.5):**
   - Dokter menginput anamnesis keluhan, diagnosa klinis, kode ICD-10 (standar medis), catatan tindakan medis (*action notes*), dan biaya jasa tindakan/prosedur (*action fee*):
     - Kidinest: Jasa dokter anak, suntik vaksinasi, paket skrining tumbuh kembang.
     - SSC: Jasa konsultasi dermatologi, chemical peeling, facial rejuvenation, laser, dll.
5. **Notifikasi Otomatis Resep ke Apotek & Penguncian Stok Obat (Fitur 3.6):**
   - Saat dokter menekan "Kirim E-Resep & Selesai Periksa":
     - Resep diterbitkan dengan status `RESERVED`.
     - Stok obat terkait otomatis terkunci (*reserved stock*) di database apotek agar tidak diserobot dokter lain.
     - Status antrian pasien beralih ke `PHARMACY_WAITING`.
     - Sistem memicu sinyal chime ke modul farmasi apotek.

### B. Batasan (Out of Scope untuk Bagian 4)
- Penyerahan fisik obat dan peracikan puyer dilakukan oleh apoteker di Bagian 5 (Modul Apotek & Farmasi).
- Pelunasan pembayaran akhir biaya tindakan dan obat dilayani di kasir meja depan (Modul Bagian 3).

---

## 2. Struktur Data & Migrasi PostgreSQL

```sql
-- 1. Master Obat (Drugs)
CREATE TABLE IF NOT EXISTS drugs (
    id SERIAL PRIMARY KEY,
    clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC', 'ALL')),
    name VARCHAR(150) NOT NULL,
    kfa_code VARCHAR(50),
    unit VARCHAR(30) NOT NULL, -- Botol, Strip, Tablet, Ampul, Tube
    category VARCHAR(50) NOT NULL,
    min_stock_alert INTEGER NOT NULL DEFAULT 10,
    buy_price INTEGER NOT NULL DEFAULT 0,
    sell_price INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Batch & Stok Obat (Drug Batches)
CREATE TABLE IF NOT EXISTS drug_batches (
    id SERIAL PRIMARY KEY,
    drug_id INTEGER NOT NULL REFERENCES drugs(id) ON DELETE CASCADE,
    batch_number VARCHAR(50) NOT NULL,
    expired_date DATE NOT NULL,
    physical_stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Rekam Konsultasi & Tindakan Medis Dokter
CREATE TABLE IF NOT EXISTS consultations (
    id SERIAL PRIMARY KEY,
    queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
    doctor_id INTEGER REFERENCES staff_users(id),
    doctor_name VARCHAR(150) NOT NULL,
    anamnesis TEXT,
    diagnosis TEXT,
    icd10_code VARCHAR(50),
    action_notes TEXT,
    action_fee INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'COMPLETED', -- DRAFT, COMPLETED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. E-Resep Dokter
CREATE TABLE IF NOT EXISTS prescriptions (
    id SERIAL PRIMARY KEY,
    consultation_id INTEGER NOT NULL REFERENCES consultations(id) ON DELETE CASCADE,
    queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
    status VARCHAR(30) NOT NULL DEFAULT 'RESERVED', -- RESERVED, PREPARING, DISPENSED, CANCELLED
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Item Rincian E-Resep (Mendukung Obat Katalog & Racikan Manual)
CREATE TABLE IF NOT EXISTS prescription_items (
    id SERIAL PRIMARY KEY,
    prescription_id INTEGER NOT NULL REFERENCES prescriptions(id) ON DELETE CASCADE,
    drug_id INTEGER REFERENCES drugs(id),
    item_name VARCHAR(150) NOT NULL,
    is_manual BOOLEAN NOT NULL DEFAULT FALSE,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price INTEGER NOT NULL DEFAULT 0,
    subtotal INTEGER NOT NULL DEFAULT 0,
    dosage_instruction VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## 3. Desain Arsitektur Backend DDD (`internal/consultation/`)

```text
backend/internal/consultation/
├── domain/
│   ├── entity.go             # Consultation, Prescription, PrescriptionItem, Master Drug
│   └── errors.go             # Domain errors (ErrPatientNotInExam, ErrEmptyPrescription, dll)
├── repository/
│   └── postgres_repo.go      # Persistensi PostgreSQL (Query FIFO dokter, simpan resep, lock stok)
├── usecase/
│   ├── consultation_usecase.go # Logika konsultasi dokter, racikan manual, status antrian
│   └── consultation_usecase_test.go # Unit testing FIFO, validasi racikan, dan penguncian stok
└── delivery/
    └── http/
        ├── handler.go        # Chi HTTP handlers
        └── routes.go         # Endpoint /api/consultations & /api/consultations/call
```

---

## 4. Desain Antarmuka Frontend (`DoctorRoom.vue`)

- **Top Bar:** Identitas dokter aktif (`dr. Hendra Pratama, Sp.A` / `dr. Nadia Kusuma, Sp.KK`), pemilih unit spesialisasi, dan status ruang periksa.
- **Daftar Antrian FIFO Sisi Kiri:**
  - Badge kategori antrian (`A_VAKSIN`, `B_SAKIT`, `C_TUMBANG`, `D_PAKET`, `E_WALKIN`, `SSC_TREATMENT`).
  - Tombol aksi: "Panggil Pasien" (status `CALLING`) & "Mulai Periksa" (status `IN_EXAM`).
- **Form Konsultasi & Rekam Medis Tengah/Kanan:**
  - Anamnesis / Keluhan pasien.
  - Diagnosa Dokter + Presets Cepat ICD-10.
  - Input Tindakan Medis & Jasa Dokter / Prosedur (Kidinest Vaksinasi vs SSC Facial/Peeling).
  - **Tabel E-Resep Cepat:**
    - Tombol tambah dari Katalog Obat (dengan stok fisik real-time).
    - Tombol tambah Racikan Manual (bebas ketik nama obat racikan & harga).
    - Kolom Aturan Pakai (Signa).
  - Tombol "Kirim E-Resep & Selesai Periksa".
