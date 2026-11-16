# Laporan Integrasi & Analisis Konflik: Branch Klinik vs Branch Accounting

**Tanggal:** 16 November 2026  
**Branch Asal Klinik:** `feat/clinic-system` (Commit: `d7ca6e4`)  
**Branch Asal Accounting:** `origin/accounting` (Commit: `41e8e9e`)  
**Branch Integrasi (Merge):** `merge/clinic-accounting`

---

## 1. Ringkasan Eksekutif

Pekerjaan modul klinik telah diamankan 100% ke dalam branch baru **`feat/clinic-system`**.  
Proses merge ke branch integrasi **`merge/clinic-accounting`** menghasilkan **4 file konflik**, dengan rincian:
- **2 Konflik Dependensi (Trivial):** `backend/go.mod` dan `backend/go.sum`.
- **2 Konflik Arsitektur Kritis (Butuh Review & Approval):** `backend/cmd/api/main.go` dan `frontend/src/App.vue`.

> [!IMPORTANT]
> **Kebutuhan Utama Klien:** *"Accounting menggunakan infra yang sama, tapi tampilan dan routes berbeda, serta **harus disembunyikan dari pihak klinik** (hanya yang berwenang akuntansi yang tahu)."*  
> Implementasi di branch `accounting` saat ini mengambil alih rute *root* `/` untuk dashboard akuntansi. Ini bertentangan dengan kebutuhan isolasi jika tidak disesuaikan.

---

## 2. Rincian & Analisis Konflik

### A. Frontend: `frontend/src/App.vue` (Kritis)

#### Kondisi Saat Ini:
1. **Branch Klinik (`feat/clinic-system`):**
   * Mengatur tampilan publik dan operasional klinik:
     - `/` & `/kidinest` $\rightarrow$ Landing Page Pediatri Kidinest
     - `/ssc` $\rightarrow$ Landing Page Estetika Semarang Skin Center
     - `/kiosk` / `/kiosk-kidinest` / `/kiosk-ssc` $\rightarrow$ Portal Pendaftaran Pasien 4-Step Wizard & Booking Jadwal Dokter
     - `/login` $\rightarrow$ Multi-role Staff Login (Dokter, Kasir, Farmasi, GRO)
     - `/admin` $\rightarrow$ Workstation Staf (Pemeriksaan Dokter, Kasir POS, Farmasi & E-Prescription)
     - `/monitor` $\rightarrow$ TV Display Antrian Pasien
2. **Branch Accounting (`origin/accounting`):**
   * Mengganti seluruh isi `App.vue` hanya menjadi `<router-view />`.
   * Di file `frontend/src/router/index.ts`, branch accounting menyetel:
     - Path `/` $\rightarrow$ `AccountingLayout.vue`
     - Path `''` $\rightarrow$ redirect ke `/dashboard` (Ikhtisar Keuangan)
     - Path `*` (Catch-all) $\rightarrow$ redirect ke `/dashboard`

#### Dampak Jika Langsung Ditimpa:
* Semua halaman klinik (pasien, dokter, kasir, farmasi) **akan hilang/tertimpa** dan setiap orang yang membuka website klinik akan langsung melihat modul akuntansi (Bagan Akun COA, Neraca, Jurnal, Pajak).

#### Rekomendasi Solusi Isolasi:
1. **Gunakan Vue Router Terpadu:**
   * Jadikan `/` tetap untuk Beranda Kidinest, `/ssc` untuk Beranda SSC, `/kiosk` untuk Pasien, `/admin` untuk Workstation Medis, `/monitor` untuk TV Display.
2. **Pindahkan Accounting ke Prefix Khusus/Tersembunyi:**
   * Ubah rute akuntansi di `frontend/src/router/index.ts` dari `/` menjadi `/accounting/*` (misal: `/accounting/dashboard`, `/accounting/journals`, `/accounting/reports/*`).
3. **Sembunyikan Navigasi:**
   * Di UI klinik, tidak ada tombol atau menu yang mengarah ke `/accounting`. Hanya staf akuntansi yang mengetahui URL atau login ke modul tersebut.

---

### B. Backend: `backend/cmd/api/main.go` (Kritis)

#### Kondisi Saat Ini:
1. **Branch Klinik (`feat/clinic-system`):**
   * Menjalankan server Chi HTTP di port `8080`.
   * Menginisialisasi koneksi PostgreSQL klinik (`sharedDb.InitPostgres`) untuk tabel:
     - `staff_users`, `patients`, `queues`, `consultations`, `pharmacy_medicines`, `doctor_schedules`, `schedule_overrides`, `satusehat_mappings`.
   * Mendaftarkan rute operasional:
     - `/api/auth/*` (Login Staf & JWT)
     - `/api/public/queue/*` & `/api/public/doctors/*` (Portal Pasien & Jadwal Dokter)
     - `/api/queues/*` (Frontdesk & Kasir)
     - `/api/doctor/*` (Konsultasi & Pengaturan Shift Dokter)
     - `/api/pharmacy/*` (Farmasi)
     - `/api/satusehat/*` (Integrasi Kemenkes)
2. **Branch Accounting (`origin/accounting`):**
   * Mengubah `main.go` menjadi runner ringkas:
     ```go
     app, err := bootstrap.NewApplication()
     app.Run()
     ```
   * Menjalankan modul akuntansi di bawah router `routes.SetupRouter(container.Handlers)`.
   * Mendaftarkan rute:
     - `/api/v1/accounts/*` (Chart of Accounts)
     - `/api/v1/journals/*` (Jurnal Umum & Approvals)
     - `/api/v1/cash/*` & `/api/v1/bank/*` (Buku Kas & Bank)
     - `/api/v1/reports/*` (Laba Rugi, Neraca, Arus Kas, General Ledger)
     - `/api/v1/tax/*`, `/api/v1/periods/*`, `/api/v1/receivables/*`, dll.

#### Dampak Jika Langsung Ditimpa:
* Jika memilih versi accounting murni: semua endpoint klinik mati (login, antrian, resep, dokter, kasir tidak bisa diakses).
* Jika memilih versi klinik murni: semua endpoint akuntansi mati.

#### Rekomendasi Solusi Integrasi:
* **Dual-Mounted Architecture di 1 Server:**
  * Pertahankan port `8080`.
  * Rute klinik tetap di-mount pada prefix `/api/*`.
  * Rute akuntansi di-mount pada prefix `/api/v1/*` (sesuai struktur bawaan branch accounting).
  * Kedua modul berjalan di atas koneksi database PostgreSQL yang sama tanpa saling mengganggu.

---

### C. Backend Dependencies: `backend/go.mod` & `backend/go.sum` (Ringan)

* **Branch Klinik** menambahkan: `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto`.
* **Branch Accounting** menambahkan: `github.com/shopspring/decimal`, `modernc.org/sqlite`.
* **Solusi:** Satukan seluruh dependensi ke dalam satu `go.mod` dan jalankan `go mod tidy`.

---

## 3. Analisis Database & Skema Tabel

Kabar baik: **Tidak ada tabrakan nama tabel (*zero table collision*)**.
* **Tabel Klinik:** `patients`, `queues`, `consultations`, `pharmacy_medicines`, `pharmacy_batches`, `doctor_schedules`, `doctor_schedule_overrides`, `satusehat_mappings`.
* **Tabel Accounting (via `backend/migrations/`):** `branches`, `accounts`, `fiscal_periods`, `journal_entries`, `journal_lines`, `invoices` (UUID), `cash_bank_books`, `tax_transactions`, `assets`, dll.

Kedua skema dapat hidup berdampingan secara aman di dalam satu database PostgreSQL yang sama.

---

## 4. Opsi Langkah Selanjutnya untuk Approval

| Pilihan | Rencana Aksi | Kelebihan | Catatan |
|---|---|---|---|
| **Opsi 1 (Direkomendasikan)** | **Isolasi Rute Akuntansi (`/accounting`) & Dual-Route Backend** | • Klinik dan Akuntansi aktif bersamaan di 1 repo.<br>• Akuntansi 100% tersembunyi dari pasien & dokter klinik.<br>• Tidak merusak fitur wizard/jadwal klinik yang baru selesai. | Perlu penyesuaian router Vue di `frontend/src/router/index.ts` dan mounting rute di `main.go`. |
| **Opsi 2** | **Fokus Penuh Accounting (Klinik Dinonaktifkan Sementara)** | • Langsung melihat dan menguji sistem accounting teman. | Halaman pendaftaran pasien dan modul dokter sementara tidak dapat diakses. |
| **Opsi 3** | **Pisahkan Sub-domain / Port Berbeda** | • Isolasi total di level port/domain (misal port 8080 klinik, port 8081 accounting). | Memerlukan konfigurasi infrastruktur & CORS tambahan. |

---

*Laporan dibuat otomatis oleh asisten sebagai bahan telaah teknis sebelum aksi merge diselesaikan.*
