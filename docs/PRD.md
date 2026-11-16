# Product Requirement Document (PRD)
## Smart Clinic Management System (Production-Ready)

---

## 1. Executive Summary & Project Status
Sistem Informasi Manajemen Klinik Terpadu (**Smart Clinic Management System**) berbasis arsitektur **Modular Monolith** dengan teknologi backend **Go (Golang)** dan frontend **Vue 3 SPA**. 

Sistem ini telah resmi disetujui (*Deal & Approved Contract*) untuk diimplementasikan secara penuh (*Production Deployment*) di dua unit klinik operasional:
1. **Kidinest:** Layanan poli anak, pemantauan tumbuh kembang, dan vaksinasi dengan antrian FIFO ketat dan pembayaran pendaftaran online terverifikasi di awal.
2. **Semarang Skin Center (SSC):** Layanan estetika medis dan perawatan kulit dengan sistem reservasi jadwal (*booking* tanpa bayar di muka / Rp 0) dan pelunasan transaksi di kasir setelah tindakan medis selesai.

Dokumen ini merupakan spesifikasi kebutuhan fungsional dan teknis **tahap produksi nyata (*Production Scope*)**, menggantikan versi rancangan demo/prototipe sebelumnya.

---

## 2. Prinsip Utama Domain & Regulasi
* **Privacy-First Patient Data:** Sistem tidak meminta dan tidak menyimpan Nomor Induk Kependudukan (NIK) / KTP pasien demi kenyamanan dan perlindungan data pribadi. Nama pasien pada layar monitor ruang tunggu publik disamarkan secara otomatis (*data masking*, contoh: `Budi P****`).
* **Paperless & Zero Friction Queue:** Tiket antrian diterbitkan secara digital, dikirim otomatis via WhatsApp Gateway, dan dipantau secara langsung (*real-time*) oleh pasien melalui browser ponsel tanpa perlu menginstal aplikasi.
* **Invarian Stok Obat Ketat (Anti-Leakage & Anti-Variance):**
  $$\text{available\_stock} = \text{physical\_stock} - \text{reserved\_stock}$$
  Stok otomatis terkunci (*reserved*) seketika saat dokter menerbitkan resep elektronik (*E-Prescription*), dan terpotong permanen dari `physical_stock` saat apoteker menekan serah terima obat.
* **Single Patient Invoice:** Pasien hanya menerima 1 lembar tagihan ringkas (Registrasi + Jasa Tindakan Medis + Total Obat). Rincian pembagian internal (Jasa Medis, Sarana Klinik, Farmasi) dicatat terpisah untuk laporan pimpinan.
* **Batas Lingkup Tegas (Out of Scope):** Modul akuntansi buku besar/jurnal umum, modul rawat inap, dan integrasi BPJS/SATUSEHAT bridging berada di luar lingkup rilis ini. Sistem menyediakan fasilitas ekspor data transaksi dan mutasi stok ke Excel/CSV untuk diimpor ke software akuntansi eksternal klinik.

---

## 3. Tech Stack & Infrastruktur Produksi

* **Backend Engine:**
  * Bahasa: **Go (Golang 1.22+)**
  * Router & Middleware: **Go Chi** / standard library `net/http`
  * Database: **PostgreSQL 16+** (ACID compliant, transaksi atomik dengan *row-level locking* `SELECT ... FOR UPDATE` untuk penerbitan nomor antrian berurutan dan alokasi stok obat concurrent).
  * Driver: `database/sql` + `github.com/jackc/pgx/v5`
* **Frontend Application:**
  * Framework: **Vue 3** (Composition API, `<script setup>`)
  * Build Tool: **Vite**
  * Styling: **Tailwind CSS** (Estetika Brandlos / Nordic Modern, warm canvas, pill badges, micro-interaction responsif)
  * State Management: **Pinia**
  * Icon Pack: **Lucide Vue Next**
* **Real-Time Synchronization:**
  * **Server-Sent Events (SSE) / Event Polling & WebSocket** untuk sinkronisasi monitor antrian publik dan pemanggilan suara (*audio chime*) di apotek & ruang dokter tanpa *full reload*.
* **Integrasi Pihak Ketiga (Production Gateways):**
  * **Payment Gateway (Tripay Open API):** Pembuatan tagihan QRIS dinamis & Virtual Account otomatis, verifikasi status lunas melalui webhook callback dengan validasi signature HMAC-SHA256.
  * **WhatsApp Notification Gateway:** Pengiriman otomatis pesan konfirmasi registrasi, nomor tiket, estimasi waktu panggil, dan tautan live monitor antrian.

---

## 4. Ruang Lingkup 27 Fitur Fungsional Produksi (5 Modul Utama)

### Bagian 1: Portal Pasien & Pendaftaran Mandiri (Self-Service)
1. **Landing Page Unit Kidinest:** Antarmuka pendaftaran khusus poli anak & vaksinasi, terhubung langsung ke gateway pembayaran pendaftaran.
2. **Landing Page Unit Semarang Skin Center (SSC):** Antarmuka reservasi perawatan estetika/kulit (*booking fee* Rp 0, pelunasan pasca tindakan).
3. **Form Pendaftaran Privacy-First:** Form cepat ramah seluler (Nama Lengkap, Nomor WhatsApp, Keluhan Singkat) tanpa syarat KTP/NIK.
4. **Kategori Antrian Kidinest:**
   * Kategori A: Vaksinasi Anak
   * Kategori B: Pasien Sakit (Poli Umum Anak)
   * Kategori C: Tumbuh Kembang
   * Kategori D: Paket Vaksin & Tumbuh Kembang
   * Kategori E: Pasien Datang Langsung (Walk-In)
5. **Pemilihan Dokter Jaga:** Pasien memilih dokter spesialis anak atau dokter tindakan sesuai jadwal aktif.
6. **Integrasi Pembayaran Pendaftaran (Tripay):** Penerbitan QRIS/VA otomatis. Setelah notifikasi webhook `PAID` diterima, nomor antrian resmi langsung terbit seketika.
7. **Notifikasi WhatsApp Otomatis:** Tiket digital berisi nama pasien, nomor antrian, dokter tujuan, dan tautan live monitor dikirim ke nomor WhatsApp pasien.
8. **Live Queue Monitor Web:** Halaman pemantau antrian interaktif di HP pasien dan layar TV ruang tunggu klinik.
9. **Penyamaran Nama Pasien:** Sensor privasi nama publik otomatis (contoh: `Ahmad R****`).
10. **Preview Antrian Berikutnya:** Notifikasi "1 antrian di depan Anda" agar pasien bersiap di depan poliklinik.
11. **Filter Antrian Per Dokter:** Pemilahan antrian pada display monitor per masing-masing dokter praktek.

### Bagian 2: Meja Depan & Kasir (Frontdesk / GRO)
12. **Bantuan Registrasi Walk-In:** Resepsionis/GRO mendaftarkan pasien darurat atau lansia melalui antarmuka kasir tanpa kertas (*paperless ticket* langsung ke WhatsApp pasien).
13. **Multi-Metode Pembayaran Kasir:** Menerima Tunai (Cash), Mesin EDC, Transfer Bank, QRIS Kasir, serta pemotongan kupon/Voucher Diskon nominal fleksibel.
14. **Cetak Bukti Transaksi Resmi:** Pencetakan struk kasir thermal 80mm ramah printer POS dan opsi export faktur resmi PDF ber-kop klinik.
15. **Tutup Shift & Pelaporan Pimpinan:** Rekapitulasi setoran kasir harian (tunai fisik laci, transaksi digital EDC/Transfer/QRIS, potongan voucher, dan refund) dalam 1 tombol laporan ke pimpinan.
16. **Pencatatan Refund Manual:** Pengembalian dana resmi untuk pembatalan tindakan: pencatatan rekening tujuan, pembaruan status `REFUNDED`, dan pengembalian stok obat ke status awal.

### Bagian 3: Ruang Dokter (Pemeriksaan & E-Resep)
17. **Antrian FIFO Ketat:** Dokter memanggil pasien berdasarkan urutan pembayaran lunas yang adil (*First In, First Out*).
18. **Filter Pasien Dokter Jaga:** Dokter hanya melihat daftar pasien yang memilih dirinya pada sesi hari tersebut.
19. **Tabel E-Resep Interaktif:**
    * *Dropdown Search* cepat obat dengan indikator stok fisik apotek real-time.
    * Tampilan transparan harga satuan obat.
    * Opsi **input manual diketik** beserta input nominal harga untuk racikan khusus atau item non-database.
20. **Komponen Tindakan Medis Kidinest:** Input jasa konsultasi dokter, suntik vaksinasi, dan tindakan poli anak.
21. **Komponen Tindakan Medis SSC:** Input tindakan estetika (facial, peeling, laser) dan produk perawatan kulit klinik.
22. **Notifikasi Otomatis ke Apotek:** Trigger suara (*chime cue*) dan kemunculan resep baru seketika di layar petugas farmasi saat dokter mengirim resep.

### Bagian 4: Farmasi & Gudang Obat (Apotek)
23. **Layar Antrian Resep Real-Time:** Apoteker memproses resep masuk dengan alur status: *Diterima* $\rightarrow$ *Diracik* $\rightarrow$ *Siap Diserahkan*.
24. **Penguncian & Pemotongan Stok Otomatis:** Mekanisme *reserved stock* mencegah selisih fisik obat saat beberapa dokter meresepkan obat yang sama secara bersamaan.
25. **Master Obat & Manajemen Batch FEFO:** Pencatatan nama obat, satuan, harga modal, harga jual, nomor batch, dan tanggal kadaluarsa (*First Expired, First Out*).
26. **Ekspor & Impor Stok Opnam (Excel/CSV):** Fasilitas unduh saldo stok obat untuk audit fisik lemari obat, serta unggah kembali hasil hitung opnam massal dalam 1 klik.
27. **Migrasi Data Sistem Lama:** Fitur impor data riwayat pasien dan daftar obat lama klinik via Excel/CSV.

### Bagian 5: Keamanan & Manajemen Akses Staf
* **Login Menggunakan Nomor Induk Karyawan (NIK Staf):** Otentikasi staf menggunakan Nomor Karyawan resmi klinik dan kata sandi terenkripsi bcrypt (bukan Google OAuth/email pribadi).
* **Role-Based Access Control (RBAC):** Pembatasan modul kerja:
  * `GRO_CASHIER`: Pendaftaran walk-in, pembayaran kasir, refund, tutup shift.
  * `DOCTOR`: Rekam medis pasien terdaftar, input tindakan, e-resep.
  * `PHARMACIST`: Antrian resep farmasi, serah terima obat, stok opnam, master obat.
  * `SUPERADMIN_OWNER`: Akses menyeluruh, pendaftaran akun staf baru, pengaturan tarif, dan rekapitulasi finansial pimpinan.
* **Manajemen Staf oleh Pimpinan:** Menu khusus Super Admin untuk membuat akun staf baru, penentuan unit klinik (`KIDINEST`, `SSC`, `ALL`), dan reset password staf yang lupa kata sandi.

---

## 5. Non-Goals & Batasan di Luar Lingkup (Out of Scope)
1. **Sistem Akuntansi Buku Besar Umum & Pajak (GL / Tax):** Jurnal debit-kredit akuntansi, neraca komprehensif, dan SPT Badan dikelola software akuntansi eksternal klinik (disediakan ekspor data transaksi kasir format CSV/Excel).
2. **Pengadaan Hardware Fisik:** Komputer, printer thermal 80mm, barcode scanner, dan TV display ruang tunggu disediakan oleh pihak klinik.
3. **Integrasi BPJS Kesehatan:** Tidak termasuk dalam rilis operasional klinik swasta non-BPJS ini.
4. **Layanan Rawat Inap (Inpatient):** Seluruh arsitektur dikhususkan untuk klinik rawat jalan (*outpatient ambulatory care*).

---

## 6. Rencana Implementasi & Roadmap Tahap Produksi
* **Milestone 1:** Selesai pengembangan modul inti operasional (Pendaftaran, Kasir, Dokter, Apotek, PostgreSQL).
* **Milestone 2 (16 November 2026):** Uji coba operasional menyeluruh dan presentasi langsung di klinik.
* **Milestone 3 (Pertengahan Desember 2026):** Integrasi live Tripay Payment Gateway (QRIS/VA) & WhatsApp Gateway, pelatihan staf klinik (*training*), uji coba penerimaan pengguna (*UAT*), dan peluncuran resmi (*Go-Live*).
