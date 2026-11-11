# Rencana Implementasi: Bagian 3 — Meja Depan & Kasir GRO (Frontdesk & Cashier Operations)

Dokumen ini adalah spesifikasi teknis dan rencana eksekusi modul Meja Depan (*Frontdesk / Guest Relations Officer*) dan Kasir klinik untuk unit **Kidinest** dan **Semarang Skin Center (SSC)** sesuai arsitektur DDD (Domain-Driven Design) di Go dan estetika Brandlos/Nordic di Vue 3.

---

## 1. Tujuan & Ruang Lingkup

### A. Sasaran Fungsional (Sesuai Proposal & PRD)
1. **Bantuan Registrasi Walk-In Paperless (Fitur 2.1):**
   - Petugas Meja Depan (GRO) dapat mendaftarkan pasien yang datang langsung tanpa smartphone / darurat.
   - Kategori otomatis dialokasikan ke `E_WALKIN` (Kidinest) atau `SSC_TREATMENT` (Semarang Skin Center).
   - Tiket digital langsung terkirim via WhatsApp pasien dan nomor antrian terbit seketika bila pembayaran pendaftaran tunai diterima di loket.
2. **Manajemen Pembayaran Multi-Metode di Kasir (Fitur 2.2):**
   - Kasir melayani pelunasan billing: biaya registrasi, jasa tindakan/konsultasi dokter, dan obat/vaksin.
   - Metode pembayaran lengkap:
     - **Tunai (*Cash*):** Dilengkapi kalkulator uang diterima dan nominal kembalian (*change*).
     - **EDC:** Input nomor referensi kartu debit/kredit.
     - **Transfer Bank:** Input rekening bank tujuan dan nomor referensi.
     - **Voucher Diskon:** Input nominal potongan harga (rupiah) yang memotong total tagihan secara fleksibel.
3. **Cetak Invoice Struk Thermal 80mm & PDF (Fitur 2.3):**
   - *Single Patient Invoice* (sesuai PRD Modul 5): Struk ringkas 1 halaman menggabungkan biaya tindakan dokter + obat + registrasi dikurangi diskon voucher.
   - Format cetak ramah printer thermal 80mm (*receipt*) serta tampilan dokumen invoice kop resmi klinik.
4. **Fitur Pelaporan Kasir Harian ke Pimpinan (Fitur 2.4):**
   - Fitur tutup shift (*End-of-Shift*).
   - Kasir cukup menekan satu tombol untuk mengunci dan mengirim rekapitulasi setoran fisik (tunai, transfer, EDC, voucher, refund) langsung ke database pimpinan klinik.
5. **Pencatatan Pengembalian Dana (Refund Manual) (Fitur 2.5):**
   - Jika pasien membatalkan kunjungan atau tindakan berbayar:
   - Kasir mencatat nomor rekening pengembalian, nama bank, atas nama pemilik rekening, dan alasan refund.
   - Status invoice menjadi `REFUNDED` dan status antrian menjadi `CANCELLED`.
   - Nominal refund tercatat otomatis sebagai pengurang dalam rekapitulasi kasir shift harian.

### B. Batasan (Out of Scope untuk Bagian 3)
- Tidak ada modul general ledger/buku besar akuntansi atau pelaporan pajak PPh/PPN (Zero Accounting invariant).
- Integrasi EDC menggunakan konfirmasi manual auth code/approval code (bukan direct serial kabel POS EDC hardware).

---

## 2. Struktur Data & Migrasi PostgreSQL

```sql
-- 1. Tabel Invoices
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    queue_id INTEGER NOT NULL REFERENCES queues(id) ON DELETE RESTRICT,
    patient_id INTEGER NOT NULL REFERENCES patients(id) ON DELETE RESTRICT,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    subtotal_registration INTEGER NOT NULL DEFAULT 0,
    subtotal_action INTEGER NOT NULL DEFAULT 0,
    subtotal_medication INTEGER NOT NULL DEFAULT 0,
    voucher_discount INTEGER NOT NULL DEFAULT 0,
    total_amount INTEGER NOT NULL DEFAULT 0,
    payment_method VARCHAR(30) NOT NULL CHECK (payment_method IN ('CASH', 'BANK_TRANSFER', 'EDC', 'QRIS')),
    cash_received INTEGER NOT NULL DEFAULT 0,
    cash_change INTEGER NOT NULL DEFAULT 0,
    payment_reference VARCHAR(100),
    cashier_staff_id INTEGER REFERENCES staff_users(id),
    cashier_name VARCHAR(150) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('PAID', 'REFUNDED')),
    refund_bank_name VARCHAR(50),
    refund_bank_account VARCHAR(50),
    refund_account_name VARCHAR(150),
    refund_reason TEXT,
    refunded_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_invoices_queue_id ON invoices(queue_id);
CREATE INDEX IF NOT EXISTS idx_invoices_created_at ON invoices(created_at);

-- 2. Tabel Cashier Shift Reports (Tutup Shift Kasir)
CREATE TABLE IF NOT EXISTS cashier_shift_reports (
    id SERIAL PRIMARY KEY,
    staff_id INTEGER NOT NULL REFERENCES staff_users(id),
    staff_name VARCHAR(150) NOT NULL,
    shift_date DATE NOT NULL,
    total_transactions INTEGER NOT NULL DEFAULT 0,
    total_cash INTEGER NOT NULL DEFAULT 0,
    total_transfer INTEGER NOT NULL DEFAULT 0,
    total_edc INTEGER NOT NULL DEFAULT 0,
    total_voucher INTEGER NOT NULL DEFAULT 0,
    total_refund INTEGER NOT NULL DEFAULT 0,
    net_total INTEGER NOT NULL DEFAULT 0,
    notes TEXT,
    sent_to_owner_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shift_reports_date ON cashier_shift_reports(shift_date);
```

---

## 3. Desain Arsitektur Backend DDD (`internal/cashier/`)

```text
backend/internal/cashier/
├── domain/
│   ├── entity.go          # Invoice, ShiftReport, WalkInRequest, DTO, Repository Interfaces
│   └── errors.go          # Domain errors (ErrInvoiceNotFound, ErrNegativeAmount, dll)
├── repository/
│   └── postgres_repo.go   # Implementasi persistensi PostgreSQL dengan query ACID
├── usecase/
│   ├── cashier_usecase.go # Logika transaksi POS, diskon voucher, refund, & shift closing
│   └── cashier_usecase_test.go # Unit testing kalkulasi pembayaran & invariant refund
└── delivery/
    └── http/
        ├── handler.go     # Chi HTTP handlers
        └── routes.go      # Route endpoints terlindungi JWT role GRO_CASHIER & SUPERADMIN_OWNER
```

---

## 4. Desain Antarmuka Frontend (`CashierView.vue`)

- **Header Staf & Status Shift:** Menampilkan nama kasir aktif, status shift, dan jam operasional.
- **3 Tab Operasional:**
  1. **Tagihan & Pembayaran Kasir (POS):**
     - Daftar antrian yang siap dibayar (menunggu verifikasi tunai pendaftaran atau pelunasan tagihan medis/obat).
     - Rincian item tagihan (biaya registrasi + tindakan dokter + obat apotek).
     - Input voucher diskon interaktif dengan kalkulasi instan.
     - Pemilihan metode pembayaran (Tunai, EDC, Transfer Bank) + kalkulator kembalian uang tunai.
     - Modal struk cetak thermal 80mm & Invoice PDF ber-kop klinik.
     - Fitur Form Refund Manual terpadu jika terjadi pembatalan kunjungan.
  2. **Bantuan Registrasi Walk-In (GRO):**
     - Pendaftaran cepat pasien di meja resepsionis (Kidinest `E_WALKIN` / SSC).
     - Input Nama Pasien, Nomor WhatsApp, Pilihan Dokter Jaga, dan Opsi Terima Tunai langsung di loket.
     - Notifikasi WhatsApp otomatis ke HP pasien & penerbitan nomor antrian FIFO.
  3. **Laporan Tutup Shift Kasir (End-of-Shift):**
     - Statistik rekapitulasi setoran fisik hari ini: Total Cash, Transfer, EDC, Potongan Voucher, Total Refund, dan Net Setoran Kasir.
     - Tabel riwayat transaksi shift berjalan.
     - Tombol 1-klik "Tutup Shift & Kirim Rekap ke Pimpinan".
