# PROPOSAL PROYEK PENGEMBANGAN PERANGKAT LUNAK
# SISTEM INFORMASI OPERASIONAL KLINIK TERPADU
### Solusi Digitalisasi Operasional Kustom Kidinest & Semarang Skin Center (SSC)

---

## HALAMAN SAMPUL (COVER PAGE)

* **Nama Proyek:** Pengembangan Smart Clinic Management System (SCMS) — Edisi Operasional Klinik
* **Klien / Organisasi:** Kidinest (Poli Anak & Tumbuh Kembang) & Semarang Skin Center / SSC (Klinik Estetika & Spesialis Kulit)
* **Penyusun Proposal:** Tim Pengembang Perangkat Lunak (Software Engineering & Solution Team)
* **Tanggal Penyusunan:** Sabtu, 7 November 2026 (Revisi: Minggu, 8 November 2026)
* **Versi Dokumen:** 2.0 (Fokus Operasional Klinik — Di Luar Modul Akuntansi)
* **Target Demo Awal:** 16 November 2026
* **Target Peluncuran Penuh (Go-Live):** Pertengahan Desember 2026

---

## PERNYATAAN KERAHASIAAN (STATEMENT OF CONFIDENTIALITY)

Seluruh informasi, data teknis, diagram arsitektur, rincian fitur, dan estimasi biaya yang tercantum dalam dokumen proposal ini bersifat rahasia (*confidential*) dan semata-mata ditujukan untuk evaluasi internal manajemen **Kidinest & Semarang Skin Center (SSC)**. Informasi di dalam dokumen ini tidak boleh disalin, diperbanyak, disebarluaskan, atau dialihkan kepada pihak ketiga tanpa persetujuan tertulis dari pihak penyusun proposal.

---

## DAFTAR ISI
1. [Ringkasan Eksekutif (Executive Summary)](#1-ringkasan-eksekutif-executive-summary)
2. [Tujuan Proyek (Project Objectives)](#2-tujuan-proyek-project-objectives)
3. [Ruang Lingkup Kerja (Scope of Work)](#3-ruang-lingkup-kerja-scope-of-work)
   - 3.1. Lingkup Fungsional Operasional Klinik (27 Fitur Terstruktur dalam 5 Bagian)
   - 3.2. Batasan di Luar Lingkup (Out of Scope — Pemisahan Sistem Akuntansi)
4. [Pendekatan Teknis & Arsitektur (Technical Approach & Stack)](#4-pendekatan-teknis--arsitektur-technical-approach--stack)
5. [Kebutuhan Non-Fungsional (Non-Functional Requirements)](#5-kebutuhan-non-fungsional-non-functional-requirements)
6. [Susunan Tim & Peran (Team Composition & Roles)](#6-susunan-tim--peran-team-composition--roles)
7. [Jadwal Pelaksanaan & Tonggak Capaian (Project Timeline & Milestones)](#7-jadwal-pelaksanaan--tonggak-capaian-project-timeline--milestones)
8. [Analisis Risiko & Mitigasi (Risk Analysis & Mitigation)](#8-analisis-risiko--mitigasi-risk-analysis--mitigation)
9. [Estimasi Biaya & Skema Pembayaran (Cost Estimate & Pricing Terms)](#9-estimasi-biaya--skema-pembayaran-cost-estimate--pricing-terms)
10. [Lembar Persetujuan (Sign-off & Approval)](#10-lembar-persetujuan-sign-off--approval)

---

## 1. RINGKASAN EKSEKUTIF (EXECUTIVE SUMMARY)

Klinik modern menghadapi tantangan operasional harian yang membutuhkan penanganan cepat dan akurat: penumpukan antrian di ruang tunggu, beban resepsionis yang kewalahan membalas chat pendaftaran WhatsApp secara manual, risiko salah baca resep obat fisik antara dokter dan apoteker, potensi selisih stok obat di apotek, serta perlunya rekapitulasi setoran kasir harian yang rapi dan transparan.

Proyek ini bertujuan mengembangkan **Sistem Informasi Manajemen Operasional Klinik Terpadu** yang dirancang khusus dan dikustomisasi sesuai alur operasional dua unit bisnis klinik:
1. **Kidinest:** Layanan poli anak, vaksinasi, dan tumbuh kembang dengan alur antrian adil (FIFO) dan pembayaran biaya registrasi terverifikasi di awal.
2. **Semarang Skin Center (SSC):** Layanan estetika medis, perawatan wajah (*facial*), dan tindakan kulit dengan skema reservasi jadwal (*booking* tanpa bayar di awal / Rp0) dan pelunasan di kasir setelah tindakan medis selesai.

> **Catatan Arsitektur (Pemisahan Sistem Akuntansi):**  
> Sesuai arahan dan strategi operasional manajemen, **sistem akuntansi umum, jurnal pembukuan, dan perpajakan dipisahkan dari aplikasi ini** dan dikelola melalui perangkat lunak akuntansi khusus terpisah milik klinik. Aplikasi yang dibangun ini berfokus murni pada **keunggulan operasional medis (*Clinical Operations Excellence*)**: pendaftaran mandiri pasien, alur antrian cerdas, kasir transaksi & invoice, resep elektronik dokter, pengendalian stok obat apotek, serta pengelolaan akun staf oleh Atasan.

---

## 2. TUJUAN PROYEK (PROJECT OBJECTIVES)

Mengacu pada prinsip **SMART Goals** (*Specific, Measurable, Achievable, Relevant, Time-bound*):

* **S (Specific):** Membangun aplikasi web operasional klinik modular yang mencakup pendaftaran mandiri QR, meja kasir GRO, ruang periksa dokter, farmasi apotek & stok opnam, serta kontrol akun staf oleh pimpinan.
* **M (Measurable):** 
  - Memangkas waktu tunggu registrasi pasien di meja depan dari 10 menit menjadi di bawah 1 menit.
  - Mengeliminasi selisih stok obat fisik (*stock variance*) menjadi 0% melalui sistem reservasi stok otomatis.
  - Menghilangkan 100% ketergantungan kertas resep fisik dokter ke apotek (*paperless prescription*).
* **A (Achievable):** Menggunakan arsitektur *Modular Monolith* berbasis teknologi Go dan Vue 3 yang terbukti sangat cepat, ringan, stabil, dan hemat sumber daya komputasi.
* **R (Relevant):** Mengatasi secara langsung titik kendala staf harian (kasir/GRO, perawat, dokter, apoteker) dan memberikan kenyamanan maksimal bagi pasien.
* **T (Time-bound):** Sesi demo operasional siap diuji coba langsung di klinik pada **16 November 2026** dan peluncuran penuh (*go-live*) tercapai pada **pertengahan Desember 2026**.

---

## 3. RUANG LINGKUP KERJA (SCOPE OF WORK)

### 3.1. Lingkup Fungsional Operasional Klinik (27 Fitur Terstruktur dalam 5 Bagian)

Seluruh fitur fungsional operasional klinik dikelompokkan ke dalam 5 bagian peran pengguna:

#### Bagian 1: Portal Pasien & Pendaftaran Mandiri (Self-Service QR)
Bagian ini digunakan langsung oleh pasien untuk mendaftar secara mandiri melalui ponsel (scan QR / tautan web) tanpa perlu antri di meja resepsionis.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **1.1** | Landing Page Unit Kidinest | Halaman pendaftaran digital poli anak & vaksinasi. Tampilan minimalis dan fokus (hanya logo klinik dan form pendaftaran), langsung terhubung ke metode pembayaran online. |
| **1.2** | Landing Page Unit Semarang Skin Center (SSC) | Halaman reservasi khusus layanan estetika/facial. Menggunakan sistem reservasi jadwal (*booking* tanpa bayar di awal/Rp0). Pembayaran dilakukan di kasir setelah tindakan. |
| **1.3** | Form Pendaftaran Mandiri (Privacy-First) | Menjaga privasi penuh pelanggan (*privacy-first*): pendaftaran tidak meminta Nomor KTP / NIK. Pasien cukup mengisi Nama Lengkap Pasien/Anak dan Nomor WhatsApp aktif untuk konfirmasi tiket antrian. |
| **1.4** | Pemilihan Kategori Antrian Kidinest | Penomoran antrian otomatis berdasarkan jenis kedatangan:<br>• **Kategori A:** Layanan Vaksinasi Anak.<br>• **Kategori B:** Pasien Sakit (Poli Umum Anak).<br>• **Kategori C:** Pemantauan Tumbuh Kembang.<br>• **Kategori D:** Paket Vaksinasi & Tumbuh Kembang.<br>• **Kategori E:** Pasien Datang Langsung / Walk-In. |
| **1.5** | Pemilihan Dokter Jaga | Pasien dapat memilih dokter spesialis anak atau dokter tindakan yang diinginkan sesuai jadwal jaga aktif (tersedia opsi untuk 2–3 dokter). |
| **1.6** | Integrasi Pembayaran Online Pendaftaran | Pembayaran biaya registrasi Kidinest via QRIS / transfer otomatis. Setelah dana diverifikasi, nomor tiket antrian resmi langsung terbit seketika. |
| **1.7** | Notifikasi Antrian Otomatis via WhatsApp | Tiket antrian, nama dokter yang dituju, estimasi waktu, serta tautan langsung untuk memantau antrian otomatis terkirim ke WhatsApp pasien. |
| **1.8** | Layar Pantau Antrian Real-Time (Live Queue Web) | Pasien dapat memantau pergerakan antrian langsung dari layar HP miliknya atau melalui TV monitor ruang tunggu klinik. |
| **1.9** | Perlindungan Privasi Pasien (Sensor Nama) | Nama pasien di layar monitor ruang tunggu publik disamarkan secara otomatis (contoh: *Budi P\*\*\*\** atau *Siti A\*\*\*\*) demi kerahasiaan data medis. |
| **1.10** | Tampilan Antrian Berikutnya (Next Queue Preview) | Menampilkan nomor antrian yang akan dipanggil berikutnya (*preview next queue*), sehingga pasien berikutnya dapat bersiap. |
| **1.11** | Pemilahan Antrian Berdasarkan Dokter | Memungkinkan pasien dan staf memfilter tampilan monitor antrian per masing-masing dokter yang sedang praktek. |

#### Bagian 2: Meja Depan & Kasir (GRO / Frontdesk)
Bagian ini digunakan oleh petugas meja depan (GRO) dan kasir untuk melayani pembayaran, mencetak bukti transaksi, dan menangani pasien darurat.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **2.1** | Bantuan Reservasi Pasien Datang Langsung (Walk-In) | Pasien yang datang langsung atau kondisi darurat tetap diproses melalui sistem reservasi online dengan dibantu oleh staf GRO di meja depan. Bebas kertas (*paperless* / tanpa struk fisik) karena nomor antrian otomatis terkirim langsung ke WhatsApp pasien. |
| **2.2** | Manajemen Pembayaran Multi-Metode di Kasir | Kasir dapat memproses transaksi melalui berbagai metode: Tunai (*Cash*), Mesin EDC, Transfer Bank, hingga pemotongan kupon/Voucher Diskon dengan nominal fleksibel. |
| **2.3** | Cetak Invoice & Struk Transaksi Fisik | Mencetak bukti pembayaran resmi: struk thermal kasir 80mm untuk pasien di tempat atau invoice format PDF ber-kop klinik resmi untuk keperluan klaim asuransi mandiri. |
| **2.4** | Fitur Pelaporan Kasir Harian ke Pimpinan | Tombol tutup shift kasir (*end-of-shift*). Kasir cukup menekan satu tombol untuk mengirim rekapitulasi setoran fisik (tunai, transfer, voucher) langsung ke pimpinan klinik. |
| **2.5** | Pencatatan Pengembalian Dana (Refund Manual) | Fitur resmi jika ada pasien yang membatalkan tindakan berbayar: kasir mencatat nomor rekening tujuan pengembalian, status transaksi menjadi batal (*refunded*), dan alokasi obat dibatalkan. |

#### Bagian 3: Ruang Dokter (Pemeriksaan & E-Resep)
Bagian ini digunakan oleh dokter di ruang praktek untuk memanggil pasien, mencatat rekam medis, dan menerbitkan resep elektronik.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **3.1** | Antrian Pasien Adil (Sistem FIFO Ketat) | Daftar antrian di ruang dokter otomatis memprioritaskan pasien yang menyelesaikan pembayaran pendaftaran lebih awal (*First In, First Out*), mencegah saling serobot antrian. |
| **3.2** | Filter Pasien Khusus Dokter Jaga | Setiap dokter hanya melihat antrian pasien yang memilih dirinya sendiri, lengkap dengan indikator keluhan dan kategori antrian (A, B, C, D, E [Walk-In], atau tindakan SSC). |
| **3.3** | Penulisan E-Resep Cepat (Tabel, Harga & Input Manual) | Dokter menulis resep dalam format tabel interaktif: dilengkapi *Dropdown Search* instan (ketik 2 huruf nama obat muncul), indikator stok fisik apotek, tampilan harga satuan obat secara transparan, serta opsi **input manual diketik** (fleksibel untuk obat racikan khusus atau item non-database beserta harga/nominalnya). |
| **3.4** | Input Komponen Tindakan Medis Kidinest | Form rincian biaya medis anak: input jasa dokter, suntik vaksinasi, dan daftar obat racikan maupun obat paten. |
| **3.5** | Input Komponen Tindakan Medis SSC | Form rincian tindakan estetika: input jasa dokter/terapis (facial, chemical peeling, dll) dan paket produk skincare klinik. |
| **3.6** | Notifikasi Otomatis Resep ke Apotek | Begitu dokter menekan tombol "Kirim Resep", sistem otomatis membunyikan tanda suara/chime di komputer apotek agar racikan obat bisa segera dipersiapkan. |

#### Bagian 4: Farmasi & Pengelolaan Gudang Obat (Apotek)
Bagian ini digunakan oleh staf farmasi untuk meracik obat, memantau ketersediaan stok fisik, dan melakukan audit berkala (*stok opnam*).

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **4.1** | Layar Antrian Resep Masuk Real-Time | Petugas apotek menerima resep dari ruang dokter secara langsung dan memprosesnya dengan alur status transparan: *Diterima*, *Diracik*, dan *Siap Diserahkan*. |
| **4.2** | Penguncian & Pemotongan Stok Otomatis (Anti Selisih) | Obat otomatis berstatus "Dipesan" (*Reserved*) saat dokter membuat resep, dan terpotong permanen saat obat diserahkan ke pasien. Menghilangkan risiko selisih hitung fisik. |
| **4.3** | Master Data Obat & Manajemen Harga | Database master obat klinik: pencatatan nama, satuan (tablet, botol, ampul, tube), harga dasar, harga jual, nomor batch, dan tanggal kadaluarsa (*expired date*). |
| **4.4** | Fitur Ekspor Data Stok ke Excel/CSV | Apoteker dapat mengunduh seluruh data saldo obat ke dalam file Excel hanya dengan 1 klik untuk persiapan audit fisik lemari obat (*stok opnam*). |
| **4.5** | Fitur Impor Data Stok dari Excel/CSV | Memperbarui jumlah stok fisik hasil hitung opnam secara massal dengan mengunggah kembali file Excel ke sistem tanpa perlu input manual satu per satu. |
| **4.6** | Modul Migrasi & Cadangan Data Sistem Lama | Fasilitas impor data riwayat pasien dan daftar obat dari aplikasi lama klinik ke dalam sistem baru agar riwayat klinik tidak hilang. |

#### Bagian 5: Keamanan Data & Manajemen Staf oleh Atasan
Bagian pengelolaan akses untuk memastikan privasi data pasien, keamanan data klinik, dan kontrol penuh pimpinan terhadap akun staf.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **5.1** | Login Staf Menggunakan Nomor Induk Karyawan | Seluruh staf (kasir, dokter, perawat, apoteker) masuk ke sistem menggunakan Nomor Karyawan resmi klinik dan kata sandi. Tidak bergantung pada akun email Google pribadi staf. |
| **5.2** | Pembatasan Hak Akses Berjenjang (Role-Based Access) | Setiap staf hanya dapat mengakses modul kerjanya masing-masing: kasir tidak dapat melihat riwayat medis, dokter fokus pada pasien, apotek fokus pada resep, dan pimpinan memegang akses operasional penuh. |
| **5.3** | Modul Pengelolaan Staf oleh Atasan (Super Admin / Pimpinan) | Menu eksklusif bagi Atasan/Pimpinan untuk mendaftarkan dan mengelola akun staf: input data karyawan baru, penerbitan Nomor Karyawan resmi, penentuan peran/hak akses (kasir, dokter, farmasi), serta kontrol reset kata sandi staf. |

---

### 3.2. Batasan di Luar Lingkup (Out of Scope)
Untuk menjaga fokus penyelesaian proyek tepat waktu, hal-hal berikut secara tegas berada di luar lingkup sistem operasional ini:

1. **Sistem Akuntansi Buku Besar & Perpajakan Umum (Accounting & Tax):**
   * Aplikasi ini berfokus murni pada operasional klinik (antrian, pendaftaran, invoice kasir, e-resep dokter, dan stok farmasi).
   * Modul akuntansi umum (jurnal debit-kredit akuntansi, buku besar, neraca saldo, laba rugi akuntansi komprehensif, dan SPT pajak badan) dipisahkan dan dikelola melalui software akuntansi terpisah milik klinik.
   * *Catatan Kemudahan:* Sistem klinik tetap menyediakan fitur ekspor data rekapitulasi transaksi kasir dan mutasi obat ke format Excel/CSV untuk mempermudah bagian keuangan klinik mengimpor data ke software akuntansi mereka.
2. **Pengadaan Perangkat Keras Fisik (Hardware):**
   * Pengadaan fisik komputer/laptop kasir, printer thermal 80mm, layar TV display monitor antrian, dan router internet (disediakan oleh pihak klinik).
3. **Integrasi Eksternal BPJS & SATUSEHAT Tahap 2:**
   * Integrasi bridging API BPJS Kesehatan atau SATUSEHAT Kemenkes RI dialokasikan pada fase roadmap pengembangan lanjutan.
4. **Modul Rawat Inap (Inpatient Care):**
   * Sistem difokuskan sepenuhnya untuk klinik rawat jalan (*outpatient*), poli anak, spesialis, dan unit tindakan estetika.

---

## 4. PENDEKATAN TEKNIS & ARSITEKTUR (TECHNICAL APPROACH & STACK)

Sistem dirancang menggunakan arsitektur **Modular Monolith** yang menggabungkan performa tinggi, kemudahan pemeliharaan, serta kesederhanaan infrastruktur tanpa dependensi yang berlebihan.

```
+-------------------------------------------------------------------------+
|                  PORTAL PASIEN (MOBILE WEB / SCAN QR)                   |
|            Kidinest (Poli Anak)  |  Semarang Skin Center (SSC)          |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
|                     GATEWAY & FRONTEND LAYER (VUE 3)                    |
|       - Antarmuka Kasir (GRO)        - Antarmuka Dokter (E-Resep)       |
|       - Antarmuka Apotek (Live)      - Dashboard Atasan (Kelola Staf)   |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
|                     CORE BACKEND ENGINE (GOLANG)                        |
|  [Queue Engine FIFO] [E-Prescription] [Stock Reserve] [Invoice Engine]  |
|  [WhatsApp Gateway Adapter]          [Payment Gateway / QRIS Adapter]   |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
|                    PERSISTENCE & STORAGE LAYER                          |
|         Relational Database (SQLite / PostgreSQL) + Audit Log           |
+-------------------------------------------------------------------------+
```

### Rincian Tumpukan Teknologi (*Tech Stack*):
* **Backend:** **Go (Golang)** — Performa eksekusi tinggi, konsumsi memori sangat hemat, dan pemrosesan antrian konkuren (*goroutines*) yang handal.
* **Frontend:** **Vue 3 (Vite + Composition API)** — Antarmuka web responsif, ringan, intuitif, dan nyaman digunakan di ponsel pasien, tablet dokter, maupun PC kasir.
* **Komunikasi Data Real-Time:** **Server-Sent Events (SSE) / WebSocket** — Memastikan monitor antrian pasien dan layar apotek memperbarui data seketika tanpa perlu me-refresh halaman (*zero refresh*).
* **Penyimpanan Data:** **Relational Database (SQLite / PostgreSQL)** — Transaksi atomik (*ACID*) untuk menjamin kuantitas stok obat dan riwayat invoice kasir selalu presisi dan konsisten.
* **Integrasi Eksternal:**
  - **WhatsApp Service API:** Pengiriman notifikasi nomor tiket antrian dan tautan live monitor.
  - **Payment Gateway Adapter:** Otomasi verifikasi pembayaran QRIS dan transfer bank.

---

## 5. KEBUTUHAN NON-FUNGSIONAL (NON-FUNCTIONAL REQUIREMENTS)

1. **Keamanan & Perlindungan Privasi (Security & Privacy):**
   - Penerapan prinsip *Privacy-First*: Tidak meminta dan tidak menyimpan Nomor KTP/NIK pasien.
   - Penyamaran otomatis nama pasien di monitor ruang tunggu publik.
   - Enkripsi kata sandi akun staf menggunakan algoritma standar industri (*bcrypt*).
   - Otorisasi berbasis peran (*Role-Based Access Control*) pada setiap modul sistem.
2. **Kecepatan & Performa (Performance):**
   - Waktu respons rata-rata sistem di bawah 200 milidetik per request.
   - Pembaruan layar monitor antrian live (*real-time sync*) dengan latensi di bawah 1 detik.
3. **Ketersediaan & Keandalan (Availability & Reliability):**
   - Target ketersediaan (*uptime*) 99.9% selama jam operasional klinik berlangsung.
   - Fitur toleransi jaringan kasir: pencatatan transaksi lokal saat koneksi internet sesaat terputus.
4. **Kemudahan Penggunaan (Usability):**
   - Pasien dapat langsung mendaftar tanpa perlu mengunduh aplikasi atau membuat akun baru.
   - Antarmuka dokter dan kasir dirancang dengan alur klik minimal (*fast-entry design*).

---

## 6. SUSUNAN TIM & PERAN (TEAM COMPOSITION & ROLES)

| Peran Tim | Tanggung Jawab Utama |
| :--- | :--- |
| **Project Manager & Solution Architect** | Mengarahkan koordinasi proyek, memastikan kepatuhan ruang lingkup operasional klinik, dan mengawasi jadwal pelaksanaan. |
| **Lead Full-Stack Developer (Go & Vue 3)** | Mengembangkan logika inti sistem backend, mesin antrian FIFO, transaksi database stok obat, serta integrasi API. |
| **Frontend & UI/UX Specialist** | Merancang desain antarmuka pengguna (estetika modern, mobile-first QR, layar monitor antrian, dan tampilan kasir). |
| **QA Engineer & Implementation Lead** | Melakukan pengujian sistem (*stress testing* & fungsional), migrasi data master obat lama, serta memandu sesi pelatihan staf klinik. |

---

## 7. JADWAL PELAKSANAAN & TONGGAK CAPAIAN (PROJECT TIMELINE & MILESTONES)

Rencana kerja terstruktur dalam 4 tahapan (*sprints*):

```
Nov 7          Nov 16                        Awal Des        Pertengahan Des 2026
  |---------------|-----------------------------|---------------------|
Discovery &     Demo Sistem                 Stok Opnam Excel      Go-Live Penuh
Arsitektur      (QR, Kasir, Dokter, Live)   & Training Staf       (Peluncuran)
```

| Periode | Milestone / Capaian Kerja | Target Deliverable |
| :--- | :--- | :--- |
| **7 November 2026** | **Milestone 1: Discovery & Architecture Sign-Off** | Dokumen spesifikasi fungsional disepakati, arsitektur teknis dan alur Excalidraw tuntas. |
| **8 – 15 November 2026** | **Sprint 1: Core System & Demo Development** | Pembangunan portal pendaftaran QR (Kidinest & SSC), antrian FIFO, layar monitor live, form kasir GRO, dan e-resep dokter. |
| **16 November 2026** | **Milestone 2: Kunjungan & Sesi Demo di Klinik** | **Uji coba operasional langsung di depan manajemen klinik:** simulasi scan QR, notifikasi WA, pemanggilan dokter, dan cetak invoice. |
| **17 – 30 November 2026** | **Sprint 2: Pharmacy, Stock Opname & Staff Management** | Integrasi notifikasi apotek live chime, stok opnam ekspor/impor Excel, dan modul Atasan menginput serta mengelola akun staf. |
| **1 – 10 Desember 2026** | **Sprint 3: Data Migration & Staff Training** | Migrasi master obat lama, pengujian stabilitas sistem, dan pelatihan langsung staf klinik (kasir, dokter, apotek). |
| **Pertengahan Desember 2026**| **Milestone 3: Go-Live & Peluncuran Resmi** | Sistem aktif sepenuhnya untuk melayani operasional harian Kidinest dan Semarang Skin Center. |

---

## 8. ANALISIS RISIKO & MITIGASI (RISK ANALYSIS & MITIGATION)

| Potensi Risiko | Tingkat Dampak | Rencana Mitigasi Pencegahan |
| :--- | :---: | :--- |
| **Format data obat dari sistem lama tidak teratur** | Sedang | Tim menyediakan template baku Excel (*format header standar*) sejak awal agar staf apotek mudah mengisi saldo fisik obat. |
| **Kendala koneksi internet di lokasi klinik** | Sedang | Sistem web dirancang ringan dengan aset lokal (*local asset caching*) sehingga kasir tetap dapat melayani pasien. |
| **Salah ketik nomor WhatsApp oleh pasien saat daftar** | Rendah | Nomor antrian tetap dapat dilihat langsung pada layar konfirmasi web di HP pasien dan terpampang di monitor TV ruang tunggu klinik. |
| **Perubahan kebutuhan staf di tengah jalan (*scope creep*)** | Sedang | Mengacu ketat pada dokumen lingkup kerja fungsional yang telah disepakati bersama. |

---

## 9. ESTIMASI BIAYA & SKEMA PEMBAYARAN (COST ESTIMATE & PRICING TERMS)

### 9.1. Rincian Paket Pengembangan Perangkat Lunak Operasional Klinik
Pengembangan sistem mencakup seluruh 5 modul operasional terpadu (Portal QR, Kasir GRO, Ruang Dokter, Farmasi & Stok Opnam, Modul Atasan):

* **Paket Lisensi Pengembangan Software Operasional Klinik:** Rp 22.000.000,- (Dua Puluh Dua Juta Rupiah)
* **Layanan Sudah Termasuk (*Included*):**
  - Pengembangan seluruh 27 fitur fungsional operasional klinik sesuai proposal.
  - Sesi presentasi demo langsung di klinik pada 16 November 2026.
  - Modul impor & migrasi data master obat lama (Excel/CSV).
  - Pelatihan (*training*) penggunaan sistem untuk seluruh staf klinik.
  - Garansi pemeliharaan (*bug fix & support*) selama 3 bulan pasca *go-live*.

### 9.2. Skema Termin Pembayaran (*Milestone Payment Schedule*)
Pembayaran dilakukan secara bertahap berbasis pencapaian tonggak kerja nyata:

| Termin Pembayaran | Persentase | Waktu Pembayaran |
| :--- | :---: | :--- |
| **Termin 1 (Uang Muka / DP)** | **30%** | Saat proposal disetujui dan kontrak pengerjaan ditandatangani. |
| **Termin 2 (Pasca Demo Sukses)** | **40%** | Setelah sesi demo langsung berhasil dipresentasikan di klinik (16 Nov 2026). |
| **Termin 3 (Peluncuran / Final)** | **30%** | Setelah sistem selesai diuji coba (*UAT*) dan siap *go-live* (Pertengahan Des 2026). |

---

## 10. LEMBAR PERSETUJUAN (SIGN-OFF & APPROVAL)

Proposal ini diajukan untuk menjadi pedoman resmi pelaksanaan kerja sama pengembangan perangkat lunak operasional klinik antara Tim Pengembang dan Manajemen Klinik.

| Diajukan Oleh: | Disetujui Oleh: |
| :---: | :---: |
| <br><br><br>____________________________________<br>**Lead Project Manager**<br>Tim Pengembang Sistem | <br><br><br>____________________________________<br>**Pimpinan / Direktur Operasional**<br>Kidinest & Semarang Skin Center (SSC) |
| Tanggal: ........................................ | Tanggal: ........................................ |
