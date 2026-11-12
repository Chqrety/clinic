# 06. Plan: God Mode Seeder & Stress-Test Playground (Dev/Demo Only)

Dokumen spesifikasi dan rencana implementasi modul **God Mode Seeder** untuk mempermudah demonstrasi sistem klinik, pengujian skenario jam sibuk (*rush hour* 50 antrian), reset demo, dan pembersihan transaksi uji coba.

---

## 1. Latar Belakang & Batasan Keamanan (Anti Yes-Man)

Saat mendemokan aplikasi kepada direksi/klien klinik, sering dibutuhkan:
1. Kondisi **"Jam Sibuk" (50 antrian)** untuk membuktikan performa sistem, paginasi, pencarian, dan display TV antrian yang ramai.
2. Kondisi **"Reset Bersih"** untuk mengulang demonstrasi alur dari awal tanpa meninggalkan sampah testing.
3. Kondisi **"Zero State" (Jam 07:00 Pagi)** saat klinik baru buka untuk memperlihatkan antrian kosong.

> [!CAUTION] Batasan Keamanan Ketat (Production Guard)
> Modul God Mode ini **HARAM** aktif di lingkungan Production. Rekam Medis Elektronik (RME) dan transaksi keuangan klinik dilindungi regulasi hukum (Permenkes 24/2022).
> - Di backend: Seluruh rute `/api/dev/*` dilindungi middleware yang otomatis mengembalikan **404 Not Found** jika `ENV=production`.
> - Di frontend: UI God Mode hanya ter-render jika `import.meta.env.DEV` atau flag demo aktif.

---

## 2. Tiga Operasi Inti God Mode

| No | Operasi | Endpoint | Deskripsi Aksi |
| :--- | :--- | :--- | :--- |
| 1 | **Reset Demo Initial** | `POST /api/dev/seeds/reset` | Membersihkan transaksi kotor dan mengembalikan **8 antrian demo default** + stok obat standar. Tidak menghapus akun staf. |
| 2 | **Simulate Rush Hour (50 Antrian)** | `POST /api/dev/seeds/rush-hour` | Men-generate **50 antrian pasien realistis** lintas Kidinest & SSC dengan variasi status, invoice lunas, resep aktif, dan panggilan TV. |
| 3 | **Wipe All Transactions** | `POST /api/dev/seeds/wipe-transactions` | Mengosongkan seluruh antrian, resep, dan invoice (`TRUNCATE` tabel transaksional). Menyisakan master data staf & master katalog obat. |

---

## 3. Komposisi Data 50 Antrian (Rush Hour Invariants)

Untuk menjamin seluruh modul (Kasir, Ruang Dokter, Apotek, Billing Shift, Monitor TV) terlihat hidup dan konsisten, 50 antrian dibagi sesuai distribusi jam sibuk nyata:

### A. Distribusi Unit & Nomor Antrian
* **Kidinest (Poli Anak): 30 Pasien**
  - Kategori A (Vaksinasi): `A-001` s/d `A-008` (8 pasien)
  - Kategori B (Anak Sakit): `B-001` s/d `B-012` (12 pasien)
  - Kategori C (Tumbuh Kembang): `C-001` s/d `C-005` (5 pasien)
  - Kategori D (Paket Khusus): `D-001` s/d `D-003` (3 pasien)
  - Kategori E (Walk-in Meja Depan): `E-001` s/d `E-002` (2 pasien)
* **Semarang Skin Center (SSC): 20 Pasien**
  - Kategori SSC (Dermatologi & Estetika): `SSC-001` s/d `SSC-020` (20 pasien)

### B. Distribusi Status & Relasi Modul
1. **`COMPLETED` (22 Pasien):**
   - Memiliki rekam medis selesai.
   - Memiliki invoice status `PAID` dengan metode pembayaran acak realistis (`CASH`, `QRIS`, `EDC`, `BANK_TRANSFER`).
   - Otomatis mengisi rekap pendapatan harian di modul **Billing & Tutup Shift** (~Rp 15.000.000 - Rp 25.000.000).
2. **`PHARMACY_WAITING` (8 Pasien):**
   - Sudah selesai konsultasi dokter, resepnya sedang antre diracik di modul **Farmasi & Apotek**.
3. **`IN_EXAM` (4 Pasien):**
   - 2 pasien di Kidinest (dr. Hendra Sp.A di Poli 1, dr. Patricia Sp.A di Poli 2).
   - 2 pasien di SSC (dr. Dian Sp.D.V.E di Ruang 1, dr. Kevin Sp.D.V.E di Ruang 2).
4. **`CALLING` (2 Pasien):**
   - 1 pasien Kidinest (`B-005`) & 1 pasien SSC (`SSC-008`) dalam status panggilan aktif (siap bunyi lonceng di Monitor TV).
5. **`WAITING` (14 Pasien):**
   - Menunggu di ruang tunggu, siap dipanggil oleh dokter berikutnya.

### C. Nama Pasien Realistis Khas Indonesia
- Anak-anak (Kidinest): *Kenzo Alvaro (4 th), Aisyah Putri (2 th), Rafathar Malik (6 th), Clarissa Aurelia (1 th), Bima Sakti (3 th), Naura Zahrani (5 th), Arsya Hermansyah (8 bln), dll.*
- Dewasa (SSC): *Nadia Saphira (28 th), Jessica Mila (31 th), Kevin Sanjaya (29 th), Laura Basuki (33 th), Nicholas Saputra (35 th), dll.*

---

## 4. Arsitektur Teknis (Filosofi Ponytail)

### A. Backend (`backend/internal/shared/database/godmode.go` & `backend/internal/dev/`)
- Cukup 1 file service helper sederhana yang mengeksekusi query batch SQL di dalam database transaction (`db.BeginTx`).
- Middleware `DevOnly`:
  ```go
  if os.Getenv("ENV") == "production" {
      http.NotFound(w, r)
      return
  }
  ```
- Tidak menambah dependensi baru (menggunakan `database/sql` bawaan).

### B. Frontend (`frontend/src/components/GodModeModal.vue`)
- Desain minimalis Nordic / Brandlos, warna dark zinc dengan aksen amber/emerald.
- **Cara Akses:**
  1. Pintasan Keyboard: `Ctrl + Shift + G` (atau `Cmd + Shift + G` di macOS).
  2. URL Query: Menambahkan `?godmode=1` di browser.
  3. Floating Pill Button kecil di pojok kiri bawah: `[⚡ DEV GOD MODE]` (hanya muncul di local dev/demo).
- **Indikator Loading & Feedback:**
  - Status spinner saat query seeding dijalankan.
  - Toast notifikasi instan saat selesai.
  - Emit event untuk me-refresh antrian di workstation aktif.

---

## 5. Rencana Langkah Kerja (Checklist Eksekusi)

- [ ] **Langkah 1 (Backend):** Buat handler seeding (`reset`, `rush-hour`, `wipe`) dan daftarkan rute `/api/dev/seeds` dengan guard `ENV != production`.
- [ ] **Langkah 2 (Frontend):** Buat komponen `GodModeModal.vue` dengan 3 tombol aksi dan keyboard shortcut `Ctrl + Shift + G`.
- [ ] **Langkah 3 (Integrasi):** Mount `GodModeModal.vue` di `App.vue` agar dapat diakses dari halaman mana pun.
- [ ] **Langkah 4 (Pengujian & Validasi):**
  - Jalankan build backend & frontend.
  - Uji eksekusi `rush-hour` 50 antrian.
  - Periksa responsivitas tabel antrian, monitor TV, ruang dokter, dan billing shift.
  - Uji eksekusi `reset` dan `wipe-transactions`.
- [ ] **Langkah 5 (Git Commit):** Simpan perubahan dengan pesan commit deskriptif.
