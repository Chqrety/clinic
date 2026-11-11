# Rencana Implementasi: Bagian 2 — Portal Pasien QR & Mesin Antrian FIFO (Patient Portal & Queue Engine)

Dokumen ini adalah spesifikasi teknis dan rencana implementasi modul pendaftaran mandiri pasien, alokasi nomor antrian berkeadilan (FIFO), serta layar monitor antrian publik untuk unit **Kidinest** dan **Semarang Skin Center (SSC)**.

---

## 1. Tujuan & Ruang Lingkup

### A. Sasaran
1. **Portal Pendaftaran Mobile-First Pasien:** Pasien scan QR di meja registrasi/banner, memilih unit klinik (Kidinest / SSC), mengisi nama & nomor WhatsApp (tanpa NIK), dan memilih kategori layanan.
2. **Aturan Invarian Antrian FIFO:**
   - **Kidinest:** Nomor antrian (`A-001`, `B-001`, dll.) baru diterbitkan secara urut setelah status pembayaran registrasi terverifikasi `PAID` (mencegah penyerobotan antrian oleh draft pendaftaran yang belum bayar).
   - **SSC:** Nomor antrian (`S-001`) diterbitkan langsung saat *booking* jadwal (biaya registrasi Rp0, pelunasan di kasir setelah tindakan).
3. **Monitor Antrian Publik (TV Display):** Tampilan layar antrian ruang tunggu dengan penyamaran nama (*privacy masking*: `Budi P****`), nomor yang sedang dipanggil, estimasi antrian berikutnya, dan pemisahan filter dokter/unit.
4. **Adapter Tiket WhatsApp:** Mengirimkan pesan konfirmasi nomor antrian dan tautan pemantauan live ke nomor WhatsApp pasien.

### B. Batasan (Out of Scope untuk Bagian 2)
- Tidak ada validasi NIK Dukcapil atau integrasi IHS SATUSEHAT (sesuai prinsip *Privacy-First*).
- Tidak ada hardware printer karcis fisik (sistem 100% *paperless*, tiket ada di WhatsApp dan web pasien).
- Pemanggilan suara dokter dan penginputan resep dikerjakan pada Bagian 4.

---

## 2. Struktur Data & Migrasi PostgreSQL

### A. Tabel `patients` (Privacy-First)
```sql
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_patients_phone ON patients(phone);
```

### B. Tabel `queues`
```sql
CREATE TABLE IF NOT EXISTS queues (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC')),
    category VARCHAR(30) NOT NULL CHECK (category IN (
        'A_VAKSIN', 'B_SAKIT', 'C_TUMBANG', 'D_PAKET', 'E_WALKIN', 'SSC_TREATMENT'
    )),
    doctor_id INTEGER REFERENCES staff_users(id),
    doctor_name VARCHAR(150) NOT NULL,
    queue_number VARCHAR(20) NOT NULL,           -- Misal: A-001, B-005, S-002
    registration_fee INTEGER NOT NULL DEFAULT 0,  -- Kidinest: Rp 50.000 (contoh); SSC: 0
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN (
        'QRIS', 'BANK_TRANSFER', 'CASH', 'UNPAID_BOOKING'
    )),
    payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('PENDING', 'PAID', 'REFUNDED')),
    queue_status VARCHAR(30) NOT NULL CHECK (queue_status IN (
        'WAITING', 'CALLING', 'IN_EXAM', 'PHARMACY_WAITING', 'COMPLETED', 'CANCELLED'
    )),
    paid_at TIMESTAMP WITH TIME ZONE,
    called_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_queues_unit_status ON queues(clinic_unit, queue_status);
CREATE INDEX IF NOT EXISTS idx_queues_created_at ON queues(created_at);
```

### C. Sequence Counter Harian (Pencegahan Race Condition)
```sql
CREATE TABLE IF NOT EXISTS daily_queue_counters (
    counter_date DATE NOT NULL,
    clinic_unit VARCHAR(20) NOT NULL,
    category_prefix VARCHAR(5) NOT NULL,
    last_seq INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (counter_date, clinic_unit, category_prefix)
);
```
*Menggunakan `SELECT last_seq FROM daily_queue_counters WHERE ... FOR UPDATE` saat pembayaran lunas agar nomor urut dijamin bebas loncat/duplikat saat diakses bersamaan.*

---

## 3. Desain Arsitektur Backend (Domain-Driven Design / DDD)

Bounded context antrian dan pasien disusun menggunakan struktur DDD modular di `backend/internal/queue/`:

```text
backend/internal/queue/
├── domain/
│   ├── entity.go             # Entity Patient, Queue, DailyQueueCounter, Category & Payment enums
│   ├── repository.go         # Interface QueueRepository & PatientRepository
│   └── errors.go             # Domain errors (ErrQueueNotFound, ErrPaymentAlreadyProcessed)
├── usecase/
│   ├── queue_usecase.go      # Logika pendaftaran, verifikasi bayar, FIFO number generation
│   └── wa_adapter.go         # Port/Adapter interface pengiriman notifikasi tiket WhatsApp
├── repository/
│   └── postgres_repo.go      # Implementasi PostgreSQL: SELECT ... FOR UPDATE daily_queue_counters
└── delivery/http/
    ├── handler.go            # HTTP Handler pendaftaran, simulasi verifikasi bayar, tiket web
    ├── monitor_handler.go    # HTTP Handler monitor TV & SSE stream antrian
    └── routes.go             # Rute /api/public/queue/*
```

### A. Endpoint Publik & Pasien (`internal/queue/delivery/http`)
1. **`POST /api/public/queue/register`**
   - Payload:
     ```json
     {
       "clinic_unit": "KIDINEST",
       "patient_name": "Budi Pratama",
       "patient_phone": "081234567890",
       "category": "A_VAKSIN",
       "doctor_id": 2,
       "payment_method": "QRIS"
     }
     ```
   - Respon (Kidinest Pending Bayar): Mengembalikan tagihan pendaftaran + data QRIS dinamis/instruksi transfer.
   - Respon (SSC Booking): Langsung menerbitkan tiket antrian `S-001` (status `WAITING`).
2. **`POST /api/public/queue/verify-payment`** *(Webhook / Simulasi Kasir)*
   - Payload: `{ "queue_id": 12, "mock_success": true }`
   - Melakukan locking counter harian, menerbitkan nomor urut `A-001`, mengubah status ke `PAID` + `WAITING`, dan men-trigger pengiriman tiket WhatsApp.
3. **`GET /api/public/queue/ticket/{id}`**
   - Mengambil status tiket pendaftaran pasien untuk layar smartphone (posisi antrian di depan, estimasi waktu).

### B. Endpoint Monitor & Real-Time Antrian
1. **`GET /api/public/queue/monitor`**
   - Query: `?unit=KIDINEST` atau `?unit=SSC` atau tanpa param (gabungan).
   - Mengembalikan daftar:
     - `currently_calling`: Pasien yang sedang dipanggil / di dalam ruang periksa.
     - `upcoming_queue`: Maksimal 5 antrian berikutnya.
     - Nama otomatis disamarkan (`Budi P****`).
2. **`GET /api/public/queue/stream` (Server-Sent Events / SSE)**
   - Mengalirkan event instan setiap kali ada panggilan antrian baru atau pergeseran status tanpa perlu reload layar TV.

---

## 4. Desain Tampilan Frontend (Vue 3 + Tailwind CSS)

### A. Tampilan Pendaftaran Mandiri Pasien (`src/views/PatientPortal.vue`)
- **Tampilan Mobile Kiosk:**
  - Pilihan Unit Klinik: Tab/Card responsif antara **Kidinest (Poli Anak & Vaksinasi)** dan **Semarang Skin Center (Estetika & Kulit)**.
  - Form Minimalis: Nama Pasien & Nomor WhatsApp (disertai kode negara `+62`).
  - Pemilihan Kategori Layanan (Pill buttons interaktif):
    - *Kidinest:* Vaksinasi (A), Sakit (B), Tumbuh Kembang (C), Paket (D).
    - *SSC:* Perawatan Wajah / Konsultasi Dokter Kulit (S).
  - Pilihan Dokter yang berpraktek hari ini.
  - Layar Pembayaran (khusus Kidinest): Modal QRIS / Info Rekening Bank + tombol simulasi bayar instan untuk demo.
  - Layar Tiket Sukses: Kartu tiket bergaya boarding pass digital dengan nomor besar (`A-005`), nama tersamar, nama dokter, dan barcode/QR pemantauan.

### B. Monitor Antrian Publik Ruang Tunggu (`src/views/QueueMonitor.vue`)
- **Tampilan TV Fullscreen (1920x1080 optimized):**
  - Header: Jam digital live, logo klinik (Kidinest & SSC), dan status unit.
  - Panel Utama (Kiri): Kotak sorotan nomor antrian yang sedang dipanggil dengan typography ultra-bold (`text-7xl`), nama dokter, dan nama ruang periksa.
  - Panel Samping (Kanan): Daftar antrian berikutnya (Urutan 1 s.d. 5) dengan badge kategori dan status `Menunggu`.
  - Micro-interaksi: Animasi pulsa halus saat nomor baru dipanggil + text-to-speech / audio chime opsional.

---

## 5. Rencana Langkah Kerja Eksekusi (Checklist DDD)

1. [x] **Step 1: Skema Database & Counter FIFO PostgreSQL**
   - Buat tabel `patients`, `queues`, dan `daily_queue_counters`.
   - Implementasikan fungsi aman konkurensi `GenerateNextQueueNumber(tx, unit, category)` dengan row locking.
2. [x] **Step 2: Bounded Context Queue (`internal/queue/`)**
   - Buat `domain/entity.go` & `repository.go` (definisi entity antrian & contract repository).
   - Buat `repository/postgres_repo.go` (query database antrian & PostgreSQL counter).
   - Buat `usecase/queue_usecase.go` (logika bisnis pendaftaran, masking nama pasien, dan verifikasi bayar).
   - Buat `delivery/http/handler.go` & `monitor_handler.go` (HTTP router publik, streaming SSE).
3. [x] **Step 3: Adapter Pesan WhatsApp (`internal/queue/usecase/wa_adapter.go`)**
   - Buat interface pengiriman pesan WA (mock adapter log untuk mode development/demo).
4. [x] **Step 4: Frontend Patient Portal (`src/views/PatientPortal.vue`)**
   - Komponen UI pemilihan unit Kidinest & SSC.
   - Form pendaftaran cepat, pop-up QRIS dummy, dan tampilan tiket digital.
5. [x] **Step 5: Frontend Monitor Antrian TV (`src/views/QueueMonitor.vue`)**
   - Desain layar monitor TV fullscreen sesuai estetika Brandlos/Nordic Modern.
   - Sambungkan ke stream SSE / interval polling data antrian.
6. [x] **Step 6: Verifikasi End-to-End**
   - Pasien mendaftar Kidinest -> berstatus PENDING dan belum dapat nomor antrian resmi sebelum bayar.
   - Verifikasi bayar berhasil -> nomor antrian urut FIFO `A-001` terbit seketika.
   - Pasien mendaftar SSC -> langsung terbit nomor reservasi `S-001` tanpa biaya awal (Rp0).
   - Layar monitor antrian TV menampilkan nomor baru dengan nama disamarkan secara otomatis (`Budi P****`).

---

## 6. Skenario Pengujian & Akseptansi
- [x] Pendaftaran Kidinest tidak menghasilkan nomor antrian saat pembayaran masih `PENDING`.
- [x] Pendaftaran simultan oleh 2 pasien menghasilkan nomor urut berurutan tanpa duplikasi (`A-001` dan `A-002`).
- [x] Pasien SSC langsung mendapatkan nomor urut `S-001` tanpa melalui tagihan registrasi.
- [x] Tampilan monitor TV menampilkan nama pasien dengan format sensor (`Rina W****`).
- [x] Pesan log tiket WhatsApp tercetak dengan format tautan pelacakan yang valid.
