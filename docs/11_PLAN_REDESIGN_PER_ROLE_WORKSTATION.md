# 11. Rencana Redesign UI Per-Role Workstation Berbasis Standar #cashier (Nordic Modern / Clinical Editorial)

## 1. Latar Belakang & Filosofi Redesign
Sistem kasir dan frontdesk ([CashierView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/CashierView.vue) / `#cashier`) telah terbukti menjadi standar emas (*golden benchmark*) antarmuka klinik:
1. **Bebas dari AI-Slop:** Tidak ada headline raksasa 72px yang berteriak, tidak ada kartu tebal bertumpuk dengan bayangan pekat (`shadow-2xl`), dan tidak ada gradasi neon ungu/pink murahan.
2. **Kerapian & Kepadatan Klinis:** Menggunakan kanvas proporsional `max-w-[1600px]`, font `Plus Jakarta Sans`, 1-row compact hairline KPI ticker strip, pill tabs berpenghitung monospaced, dan tata letak master-detail dua kolom yang ergonomis.
3. **Ergonomi Staf:** Staf tidak perlu berpindah halaman atau menggulung layar secara berlebihan untuk menyelesaikan siklus kerja (FIFO antrian, aksi cepat, dan status feedback).

Sesuai pembaruan [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md) Bagian 5 & 6, dokumen ini merinci rencana transformasi per-role untuk menyelaraskan 6 workstation internal lainnya dengan standar `#cashier`.

---

## 2. Blueprint Komponen Bersama (*The Cashier Standard*)
Setiap workstation akan mengadopsi 4 blok arsitektur identik:
```
┌────────────────────────────────────────────────────────────────────────┐
│ 1. Minimalist Clinical Header (Breadcrumb, Title, Duty Badge, Sync)     │
├────────────────────────────────────────────────────────────────────────┤
│ 2. Compact Operational Ticker Strip (1-Row Hairline Grid, 4-5 Metrics) │
├────────────────────────────────────────────────────────────────────────┤
│ 3. Navigation Tabs / Segmented Filter Pills (All / Unit / Status)      │
├────────────────────────────────────────────────────────────────────────┤
│ 4. Master-Detail Two-Pane Grid (lg:grid-cols-12 gap-6)                 │
│   ├── Left Pane: Queue / Catalog Feed (5 atau 7 Cols)                  │
│   └── Right Pane: Execution & Action Drawer (7 atau 5 Cols Sticky)     │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Rencana Redesign Per-Role Workstation

### Role 1: Ruang Periksa Dokter (`#doctor` - `DoctorRoom.vue`)
* **Masalah Saat Ini:**
  - Terkunci di `max-w-7xl` sehingga sempit di layar monitor klinik.
  - Masih menggunakan badge biru terang AI (`bg-blue-50 text-blue-900 border-blue-200`) dan card berbayang tanpa keseragaman.
  - Belum memiliki 1-row hairline metric ticker strip (antrian menunggu, pasien selesai, rata-rata durasi konsultasi).
  - Kolom resep obat dan tindakan terasa terpisah dari data anamnesis pasien.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Ruang Periksa Dokter • Semarang Medical Group`.
     - Selector dokter jaga yang ramping dan duty badge aktif.
     - Tombol `Sync` antrian dan switch panggilan cepat.
  2. **1-Row KPI Metric Ticker Strip:**
     - `Antre Konsultasi` (angka dan estimasi tunggu).
     - `Sedang Diperiksa` (nama & nomor antrian aktif).
     - `Selesai Hari Ini` (jumlah pasien yang telah dilayani).
     - `E-Resep Terkirim` (status resep masuk ke apotek).
  3. **Master-Detail Two-Pane Grid (`lg:grid-cols-12`):**
     - **Left Pane (`lg:col-span-5`):** Antrian Pasien FIFO dengan filter status (`MENUNGGU`, `DIPANGGIL`, `SELESAI`), tombol Panggil Suara terpadu dengan Web Speech TTS, serta status alert triase.
     - **Right Pane (`lg:col-span-7`, Sticky):** Rekam Medis & E-Resep Terpadu:
       - Ringkasan Pasien & Riwayat Kunjungan.
       - Form Anamnesis & Diagnosa (pencarian ICD-10 instan).
       - Tabel E-Resep Cepat: Baris obat katalog/racikan manual dengan kalkulasi harga transparan dan dosis terstruktur.
       - Tombol Aksi: `Kirim ke Apotek & Selesai Konsultasi` (warna `bg-zinc-950 hover:bg-zinc-800`).

---

### Role 2: Apotek & Farmasi Terpadu (`#pharmacy` - `PharmacyRoom.vue`)
* **Masalah Saat Ini:**
  - Card resep berderet vertikal satu kolom panjang dengan banyak badge warna-warni yang membingungkan mata apoteker.
  - Tidak ada panel verifikasi resep racikan yang fokus saat menyiapkan obat.
  - Belum memiliki 1-row hairline metric ticker strip.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Instalasi Farmasi & Apotek • Semarang Medical Group`.
     - Duty badge Apoteker Penanggung Jawab.
     - Aksi cepat: `Unduh Template Opnam (CSV)`, `Unggah CSV Opnam`, dan tombol `Sync`.
  2. **1-Row KPI Metric Ticker Strip:**
     - `Resep Masuk (Reserved)` (stok telah teralokasi).
     - `Sedang Diracik (Preparing)` (dalam proses dispensing).
     - `Siap Diserahkan (Ready)` (menunggu pasien dipanggil).
     - `Total Penyerahan Selesai (Dispensed)`.
     - `Nilai Obat Terlayani` (ringkasan omset obat hari ini).
  3. **Master-Detail Two-Pane Grid (`lg:grid-cols-12`):**
     - **Left Pane (`lg:col-span-5`):** Feed Resep Real-Time terbagi tab status (`Semua`, `Menunggu Racik`, `Siap Serah`). Dilengkapi badge unit (`KIDINEST` / `SSC`), nama dokter peresep, dan timer tunggu resep.
     - **Right Pane (`lg:col-span-7`, Sticky):** Workstation Dispensing & Telaah Obat:
       - Rincian resep terpilih: Diagnosa dokter, aturan pakai / etiket.
       - Checklist Verifikasi 7 Benar Obat (Pasien, Obat, Dosis, Rute, Waktu, Dokumentasi, Informasi).
       - Tombol Aksi Utama:
         - `Mulai Racik (PREPARING)`.
         - `Serahkan Obat ke Pasien (DISPENSED)` $\rightarrow$ otomatis merilis reservasi dan memotong stok fisik.
         - `Batal / Koreksi Resep`.

---

### Role 3: Gudang Farmasi & Stok Opnam (`#inventory` - `InventoryView.vue`)
* **Masalah Saat Ini:**
  - Filter dan tombol aksi tersebar di banyak posisi.
  - Tampilan tabel stok obat masih memiliki whitespace berlebih dan kurang tabular.
  - Tab navigasi belum menggunakan format pill sleek standar kasir.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Gudang Logistik Medis & Farmasi • Semarang Medical Group`.
     - Duty badge Petugas Gudang / Farmasis.
     - Aksi tombol: `+ Tambah Obat Baru`, `Unduh CSV Opnam`, `Unggah CSV Opnam`, `Sync`.
  2. **1-Row KPI Metric Ticker Strip:**
     - `Total Item Obat` (SKU aktif di katalog).
     - `Stok Menipis (Alert)` (di bawah batas minimum stok).
     - `Mendekati Kadaluarsa (FEFO 90 Hari)` (peringatan batch).
     - `Total Nilai Valuasi Fisik` (HPP stok riil).
  3. **Tab Navigasi Sleek:**
     - `[ Katalog Stok Obat (FEFO) ]` | `[ Rekonsiliasi Stok Opnam ]` | `[ Riwayat Mutasi & Audit ]`.
  4. **High-Density Data Table:**
     - Tampilan berdensitas tinggi dengan angka tabular monospaced (`font-mono`).
     - Badge peringatan FEFO (Kedaluwarsa hijau/amber/rose).
     - Rincian batch dan pergerakan stok per-lot dapat diexpand secara seamless.

---

### Role 4: Billing & Tutup Shift (`#billing` - `BillingView.vue`)
* **Masalah Saat Ini:**
  - Memisahkan billing dengan gaya UI lama yang belum terhubung harmonis dengan rekapitulasi shift kasir.
  - Tombol aksi mode pasien vs internal kurang kontras dan ergonomis.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Billing Kasir & Rekap Keuangan • Semarang Medical Group`.
     - Selector unit dan tanggal rekapitulasi.
     - Tombol cetak laporan shift dan sync data invoice.
  2. **1-Row KPI Metric Ticker Strip:**
     - Menggunakan ticker strip yang sama persis dengan `CashierView.vue`:
       `Tunai Laci` | `EDC & Transfer` | `Diskon Voucher` | `Total Refund` | `Net Pendapatan Shift`.
  3. **Master-Detail Two-Pane Grid (`lg:grid-cols-12`):**
     - **Left Pane (`lg:col-span-5`):** Daftar Invoice Resmi (Single Invoice) dengan indikator status pembayaran (`LUNAS`, `REFUND`), filter metode bayar, dan search pasien/nomor invoice.
     - **Right Pane (`lg:col-span-7`, Sticky):** Preview Struk Thermal / Lembar Tagihan Medis:
       - Header resmi Semarang Medical Group (Kidinest & SSC).
       - Breakdown terpisah: Registrasi, Jasa Medis Dokter, Obat Apotek, Diskon Voucher.
       - Aksi: `Cetak Thermal Struk (58/80mm)`, `Cetak Invoice PDF A4`, `Ajukan Pengembalian Dana (Refund)`.

---

### Role 5: Owner & Staff Admin (`#owner` - `StaffAdmin.vue`)
* **Masalah Saat Ini:**
  - Desain form modal dan kartu staff belum rapi, masih menggunakan card berbayang lembut generik.
  - Filter peran staf kurang terintegrasi.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Manajemen Staf & Hak Akses • Semarang Medical Group`.
     - Duty badge: `dr. Ratna Paramita, M.Kes (Superadmin / Owner)`.
     - Tombol `+ Tambah Akun Staf` dan `Sync`.
  2. **1-Row KPI Metric Ticker Strip:**
     - `Dokter Spesialis Aktif`.
     - `Apoteker & Asisten Farmasi`.
     - `Petugas Kasir & Frontdesk`.
     - `Total Akun Terdaftar`.
  3. **Data Grid & Management Table:**
     - High-density staff table dengan avatar inisial monospaced, badge role spesifik (`bg-stone-100 text-zinc-900 border-zinc-200`), dan switch status aktif instan.
     - Modal dialog tambah staf & reset password sesuai standar modal minimalis `#cashier` (`rounded-3xl border border-zinc-200/80 shadow-2xl p-6`).

---

### Role 6: SATUSEHAT Regulatory Bridge (`#satusehat` - `SatusehatView.vue`)
* **Masalah Saat Ini:**
  - Menggunakan card tebal `shadow-[0_2px_12px_rgba(0,0,0,0.03)]` dengan tombol amber besar dan banner status terpisah yang berantakan.
  - JSON inspector kurang nyaman dibaca dan tidak memadukan master-detail.
* **Transformasi Berbasis #cashier:**
  1. **Clinical Top Bar:**
     - Breadcrumb: `Interoperabilitas Kemenkes • SATUSEHAT FHIR R4`.
     - Status Faskes: `TERVALIDASI (Org ID: 100028919)`.
     - Tombol Aksi: `Sinkronisasi Sekarang (DTO API)`.
  2. **1-Row KPI Metric Ticker Strip:**
     - `Pasien IHS Terverifikasi` (rasio NIK valid Dukcapil).
     - `Master Obat Terpetakan (KFA)`.
     - `Encounter Terkirim`.
     - `MedicationDispense Terkirim`.
  3. **Master-Detail Two-Pane Grid (`lg:grid-cols-12`):**
     - **Left Pane (`lg:col-span-5`):** Log Transaksi Sinkronisasi Kunjungan & Resep dengan status kode HTTP (`200 OK`, `201 Created`, `Failed`).
     - **Right Pane (`lg:col-span-7`, Sticky):** FHIR R4 Inspector & Developer JSON Viewer:
       - Format syntax highlighting rapi dengan monospaced font.
       - Validasi Resource: `Patient`, `Encounter`, `Condition`, `MedicationRequest`, `MedicationDispense`.

---

## 4. Urutan Eksekusi / Roadmap Implementasi

| Fase | Target Workstation | File Terkait | Estimasi Aksi | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Fase 1** | Ruang Dokter (`#doctor`) | [DoctorRoom.vue](file:///home/alpi/Projects/clinic/frontend/src/views/DoctorRoom.vue) | Refactor layout `max-w-[1600px]`, 1-row ticker strip, master-detail consultation split. | **SELESAI (Verified)** |
| **Fase 2** | Apotek (`#pharmacy`) | [PharmacyRoom.vue](file:///home/alpi/Projects/clinic/frontend/src/views/PharmacyRoom.vue) | Refactor layout `max-w-[1600px]`, 1-row ticker strip, dispensing split & checklist telaah. | **SELESAI (Verified)** |
| **Fase 3** | Gudang Obat (`#inventory`) | [InventoryView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/InventoryView.vue) | High-density FEFO table, 1-row ticker strip, unified search & export actions. | **SELESAI (Verified)** |
| **Fase 4** | Billing & Shift (`#billing`) | [BillingView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/BillingView.vue) | Ticker strip financial, single-invoice preview & thermal printer integration. | **SELESAI (Verified)** |
| **Fase 5** | Staff Admin (`#owner`) | [StaffAdmin.vue](file:///home/alpi/Projects/clinic/frontend/src/views/StaffAdmin.vue) | High-density user table, 1-row staff stats strip, sleek modal styling. | **SELESAI (Verified)** |
| **Fase 6** | SATUSEHAT Bridge (`#satusehat`) | [SatusehatView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/SatusehatView.vue) | Clean developer FHIR payload inspector & Kemenkes sync metrics strip. | **SELESAI (Verified)** |

---

## 5. Kriteria Keberhasilan & Verifikasi
1. **Zero Build Error:** Seluruh perubahan lulus kompilasi `npm run build` tanpa peringatan tipe TypeScript.
2. **Visual Consistency:** Seluruh workstation memiliki:
   - Lebar kontainer konsisten `max-w-[1600px]`.
   - Header klinis dengan breadcrumb `• Semarang Medical Group`.
   - 1-row compact hairline KPI metric ticker strip di bagian atas.
   - Master-detail 2-pane grid berjarak `gap-6` dengan hairline borders `border-zinc-200/90 shadow-2xs`.
   - Tipografi terpadu `font-['Plus_Jakarta_Sans',sans-serif]`.
3. **Invarian Domain Terlindungi:** Tidak ada logika bisnis (FIFO antrian, Single Invoice, FEFO batch, atau reservasi stok obat) yang terganggu oleh perubahan tampilan visual.
