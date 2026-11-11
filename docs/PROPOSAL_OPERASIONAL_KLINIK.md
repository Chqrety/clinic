# PROPOSAL PROYEK PENGEMBANGAN PERANGKAT LUNAK
# SISTEM INFORMASI MANAJEMEN OPERASIONAL KLINIK
### Digitalisasi Alur Pelayanan Kidinest & Semarang Skin Center (SSC)

---

## HALAMAN SAMPUL (COVER PAGE)

* **Judul Proyek:** Pengembangan Smart Clinic Management System (SCMS)
* **Klien / Instansi:** Kidinest (Poli Anak & Tumbuh Kembang) & Semarang Skin Center / SSC (Klinik Estetika & Spesialis Kulit)
* **Penyusun Dokumen:** Tim Pengembang Perangkat Lunak (Software Engineering Team)
* **Tanggal Penyusunan:** Sabtu, 7 November 2026 (Diperbarui: Minggu, 8 November 2026)
* **Target Demo Sistem di Klinik:** 16 November 2026
* **Target Peluncuran Operasional Penuh (Go-Live):** Pertengahan Desember 2026

---

## PERNYATAAN KERAHASIAAN (STATEMENT OF CONFIDENTIALITY)

Seluruh informasi, data alur operasional, diagram teknis, spesifikasi fitur, dan estimasi biaya yang tercantum dalam dokumen ini bersifat rahasia (*confidential*) dan semata-mata ditujukan untuk evaluasi internal manajemen **Kidinest & Semarang Skin Center (SSC)**. Informasi di dalam dokumen ini tidak boleh disalin, diperbanyak, disebarluaskan, atau dialihkan kepada pihak ketiga tanpa persetujuan tertulis dari pihak penyusun.

---

## DAFTAR ISI
1. [Ringkasan Eksekutif (Executive Summary)](#1-ringkasan-eksekutif-executive-summary)
2. [Tujuan Proyek (Project Objectives)](#2-tujuan-proyek-project-objectives)
3. [Ruang Lingkup Kerja (Scope of Work)](#3-ruang-lingkup-kerja-scope-of-work)
   - 3.1. Lingkup Fungsional Operasional (27 Fitur dalam 5 Bagian Layanan)
   - 3.2. Batasan di Luar Lingkup (Out of Scope)
4. [Pendekatan Teknis & Arsitektur (Technical Approach & Stack)](#4-pendekatan-teknis--arsitektur-technical-approach--stack)
5. [Kebutuhan Non-Fungsional (Non-Functional Requirements)](#5-kebutuhan-non-fungsional-non-functional-requirements)
6. [Susunan Tim & Peran (Team Composition & Roles)](#6-susunan-tim--peran-team-composition--roles)
7. [Jadwal Pelaksanaan & Tonggak Capaian (Project Timeline & Milestones)](#7-jadwal-pelaksanaan--tonggak-capaian-project-timeline--milestones)
8. [Analisis Risiko & Mitigasi (Risk Analysis & Mitigation)](#8-analisis-risiko--mitigasi-risk-analysis--mitigation)
9. [Estimasi Biaya & Skema Pembayaran (Cost Estimate & Pricing Terms)](#9-estimasi-biaya--skema-pembayaran-cost-estimate--pricing-terms)
10. [Lembar Persetujuan (Sign-off & Approval)](#10-lembar-persetujuan-sign-off--approval)

---

## 1. RINGKASAN EKSEKUTIF (EXECUTIVE SUMMARY)

Klinik dengan multi-layanan menghadapi tantangan operasional harian yang membutuhkan penanganan cepat dan presisi: penumpukan antrian di ruang tunggu, beban resepsionis yang kewalahan membalas chat pendaftaran WhatsApp secara manual, risiko salah baca resep obat fisik antara dokter dan apoteker, potensi selisih stok obat di apotek, serta perlunya rekapitulasi setoran kasir harian yang rapi dan dapat dipertanggungjawabkan kepada pimpinan.

Proyek ini bertujuan mengembangkan **Sistem Informasi Manajemen Operasional Klinik** yang dirancang secara kustom mengikuti standar operasional dua unit bisnis klinik:
1. **Kidinest:** Layanan poli anak, vaksinasi, dan tumbuh kembang dengan alur antrian adil (*First In, First Out / FIFO*) dan penyelesaian pembayaran pendaftaran di awal melalui integrasi pembayaran digital.
2. **Semarang Skin Center (SSC):** Layanan estetika medis, perawatan wajah (*facial*), dan tindakan kulit dengan skema reservasi jadwal (*booking* tanpa bayar di awal / Rp0) dan pelunasan di kasir meja depan setelah tindakan selesai dilakukan.

### Keunggulan Utama Sistem:
* **Pendaftaran Mandiri & Bebas Kertas (Self-Service QR & Paperless):** Pasien mendaftar mandiri via scan QR dari ponsel, menerima tiket antrian langsung via WhatsApp, dan memantau giliran periksa secara live tanpa perlu mencetak kertas antrian fisik.
* **Perlindungan Privasi Pasien (Privacy-First):** Pendaftaran tidak meminta Nomor KTP/NIK pasien demi kenyamanan privasi, dan nama pasien di layar monitor ruang tunggu disamarkan secara otomatis (*Budi P\*\*\*\**).
* **Integrasi E-Resep & Gudang Apotek Real-Time:** Dokter menulis resep melalui tabel interaktif cepat (lengkap dengan pencarian nama obat, info stok fisik, nominal harga obat, dan opsi input manual racikan) yang otomatis membunyikan notifikasi suara di komputer apotek serta mengunci stok obat secara aman.
* **Tertib Kasir & Kontrol Pimpinan:** Kasir dapat memproses multi-metode pembayaran, mencetak bukti transaksi, menangani pengembalian dana (*refund manual*), dan mengirimkan rekapitulasi setoran harian ke pimpinan klinik dalam satu kali klik. Atasan memegang kendali penuh atas akun dan hak akses seluruh staf.

---

## 2. TUJUAN PROYEK (PROJECT OBJECTIVES)

Mengacu pada prinsip **SMART Goals** (*Specific, Measurable, Achievable, Relevant, Time-bound*):

* **S (Specific):** Membangun aplikasi web operasional klinik modular yang mencakup pendaftaran mandiri QR pasien, kasir meja depan (GRO), ruang periksa dokter, farmasi apotek & stok opnam, serta kontrol akun staf oleh Atasan/Pimpinan.
* **M (Measurable):** 
  - Memangkas waktu tunggu registrasi pasien di meja depan dari rata-rata 10 menit menjadi di bawah 1 menit.
  - Mengeliminasi selisih stok obat fisik (*stock variance*) menjadi 0% melalui sistem reservasi stok otomatis.
  - Menghilangkan 100% ketergantungan kertas resep fisik dokter ke apotek (*paperless prescription*).
* **A (Achievable):** Menggunakan arsitektur *Modular Monolith* berbasis teknologi Go dan Vue 3 yang terbukti sangat cepat, ringan, stabil, dan hemat sumber daya komputasi.
* **R (Relevant):** Mengatasi secara langsung titik kendala kerja harian staf (GRO, perawat, dokter, apoteker) dan memberikan kenyamanan maksimal bagi pasien.
* **T (Time-bound):** Sesi demo operasional siap diuji coba langsung di klinik pada **16 November 2026** dan peluncuran operasional penuh (*go-live*) tercapai pada **pertengahan Desember 2026**.

---

## 3. RUANG LINGKUP KERJA (SCOPE OF WORK)

### 3.1. Lingkup Fungsional Operasional (27 Fitur dalam 5 Bagian Layanan)

Seluruh fitur operasional klinik dikelompokkan ke dalam 5 bagian alur kerja:

#### Bagian 1: Portal Pasien & Pendaftaran Mandiri (Self-Service QR)
Bagian ini digunakan langsung oleh pasien untuk mendaftar secara mandiri melalui ponsel (scan QR / tautan web) tanpa perlu mengantri di meja resepsionis.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **1.1** | Landing Page Unit Kidinest | Halaman pendaftaran digital poli anak & vaksinasi. Desain minimalis dan fokus (hanya logo klinik dan form pendaftaran), langsung terhubung ke metode pembayaran online. |
| **1.2** | Landing Page Unit Semarang Skin Center (SSC) | Halaman reservasi khusus layanan estetika/facial. Menggunakan sistem reservasi jadwal (*booking* tanpa bayar di awal/Rp0). Pembayaran diselesaikan di kasir setelah tindakan. |
| **1.3** | Form Pendaftaran Mandiri (Privacy-First) | Menjaga privasi penuh pelanggan (*privacy-first*): pendaftaran tidak meminta Nomor KTP / NIK. Pasien cukup mengisi Nama Lengkap Pasien/Anak dan Nomor WhatsApp aktif untuk konfirmasi tiket antrian. |
| **1.4** | Pemilihan Kategori Antrian Kidinest | Penomoran antrian otomatis berdasarkan jenis kedatangan:<br>• **Kategori A:** Layanan Vaksinasi Anak.<br>• **Kategori B:** Pasien Sakit (Poli Umum Anak).<br>• **Kategori C:** Pemantauan Tumbuh Kembang.<br>• **Kategori D:** Paket Vaksinasi & Tumbuh Kembang.<br>• **Kategori E:** Pasien Datang Langsung / Walk-In. |
| **1.5** | Pemilihan Dokter Jaga | Pasien dapat memilih dokter spesialis anak atau dokter tindakan yang diinginkan sesuai jadwal jaga aktif (tersedia opsi untuk 2–3 dokter). |
| **1.6** | Integrasi Pembayaran Online Pendaftaran | Pembayaran biaya registrasi Kidinest via QRIS / transfer otomatis. Setelah dana diverifikasi, nomor tiket antrian resmi langsung terbit seketika. |
| **1.7** | Notifikasi Antrian Otomatis via WhatsApp | Tiket antrian, nama dokter yang dituju, estimasi jam pelayanan, serta tautan langsung untuk memantau antrian otomatis terkirim ke WhatsApp pasien. |
| **1.8** | Layar Pantau Antrian Real-Time (Live Queue Web) | Pasien dapat memantau pergerakan antrian langsung dari layar HP miliknya atau melalui TV monitor ruang tunggu klinik secara langsung (*real-time*). |
| **1.9** | Perlindungan Privasi Pasien (Sensor Nama) | Nama pasien di layar monitor ruang tunggu publik disamarkan secara otomatis (contoh: *Budi P\*\*\*\** atau *Siti A\*\*\*\*) demi menjaga kerahasiaan identitas. |
| **1.10** | Tampilan Antrian Berikutnya (Next Queue Preview) | Menampilkan nomor antrian yang akan dipanggil berikutnya (*preview next queue*), sehingga pasien berikutnya dapat bersiap. |
| **1.11** | Pemilahan Antrian Berdasarkan Dokter | Memungkinkan pasien dan staf memfilter tampilan monitor antrian per masing-masing dokter yang sedang praktek. |

#### Bagian 2: Meja Depan & Kasir (GRO / Frontdesk)
Bagian ini digunakan oleh petugas meja depan (GRO) dan kasir untuk melayani pembayaran, mencetak bukti transaksi, dan menangani pasien darurat.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **2.1** | Bantuan Reservasi Pasien Datang Langsung (Walk-In) | Pasien yang datang langsung atau kondisi darurat tetap diproses melalui sistem reservasi online dengan dibantu oleh staf GRO di meja depan. Bebas kertas (*paperless* / tanpa struk fisik) karena nomor antrian otomatis terkirim langsung ke WhatsApp pasien. |
| **2.2** | Manajemen Pembayaran Multi-Metode di Kasir | Kasir dapat memproses transaksi melalui berbagai metode: Tunai (*Cash*), Mesin EDC, Transfer Bank, hingga pemotongan kupon/Voucher Diskon dengan input nominal yang fleksibel. |
| **2.3** | Cetak Invoice & Struk Transaksi Fisik | Mencetak bukti pembayaran resmi: struk thermal kasir 80mm untuk pasien di tempat atau invoice format PDF ber-kop klinik resmi untuk keperluan klaim asuransi mandiri. |
| **2.4** | Fitur Pelaporan Kasir Harian ke Pimpinan | Tombol tutup shift kasir (*end-of-shift*). Kasir cukup menekan satu tombol untuk mengirim rekapitulasi setoran fisik (tunai, transfer, voucher) langsung ke pimpinan klinik. |
| **2.5** | Pencatatan Pengembalian Dana (Refund Manual) | Fitur resmi jika ada pasien yang membatalkan tindakan berbayar: kasir mencatat nomor rekening tujuan pengembalian dana, status transaksi diperbarui menjadi batal (*refunded*), dan alokasi obat dibatalkan. |

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
| **4.3** | Master Data Obat & Manajemen Harga | Database master obat klinik: pencatatan nama, satuan (tablet, botol, ampul, tube), harga dasar, harga jual ke pasien, nomor batch, dan tanggal kadaluarsa (*expired date*). |
| **4.4** | Fitur Ekspor Data Stok ke Excel/CSV | Apoteker dapat mengunduh seluruh data saldo obat ke dalam file Excel hanya dengan 1 klik untuk persiapan audit fisik lemari obat (*stok opnam*). |
| **4.5** | Fitur Impor Data Stok dari Excel/CSV | Memperbarui jumlah stok fisik hasil hitung opnam secara massal dengan mengunggah kembali file Excel ke sistem tanpa perlu input manual satu per satu. |
| **4.6** | Modul Migrasi & Cadangan Data Sistem Lama | Fasilitas impor data riwayat pasien dan daftar obat dari aplikasi lama klinik ke dalam sistem baru agar riwayat operasional masa lalu tidak hilang. |

#### Bagian 5: Keamanan Data & Manajemen Staf oleh Atasan
Bagian pengelolaan akses untuk memastikan privasi data, kelancaran operasional klinik, dan kontrol penuh pimpinan terhadap akun staf.

| No | Item Fitur | Keterangan & Manfaat Fungsional |
| :---: | :--- | :--- |
| **5.1** | Login Staf Menggunakan Nomor Induk Karyawan | Seluruh staf (kasir, dokter, perawat, apoteker) masuk ke sistem menggunakan Nomor Karyawan resmi klinik dan kata sandi. Bebas dari ketergantungan akun email Google pribadi staf. |
| **5.2** | Pembatasan Hak Akses Berjenjang (Role-Based Access) | Setiap staf hanya dapat mengakses modul kerjanya masing-masing: kasir tidak dapat melihat riwayat medis, dokter fokus pada pasien, apotek fokus pada resep, dan pimpinan memegang akses operasional penuh. |
| **5.3** | Modul Pengelolaan Staf oleh Atasan (Super Admin / Pimpinan) | Menu eksklusif bagi Atasan/Pimpinan untuk mendaftarkan dan mengelola akun staf: input data karyawan baru, penerbitan Nomor Karyawan resmi, penentuan peran/hak akses (kasir, dokter, farmasi), serta kontrol reset kata sandi staf. |

---

### 3.2. Batasan di Luar Lingkup (Out of Scope)
Untuk menjaga agar pengerjaan proyek tepat waktu, terukur, dan fokus pada keunggulan operasional klinik, batasan lingkup kerja disepakati sebagai berikut:

1. **Pengadaan Perangkat Keras Fisik (Hardware & Infrastructure):**
   * Pengadaan unit fisik komputer/laptop kerja, printer thermal kasir 80mm, layar TV display monitor antrian, dan router jaringan internet lokal (disediakan mandiri oleh pihak klinik).
2. **Layanan Rawat Inap & Tindakan Operasi Bedah Mayor:**
   * Sistem berfokus sepenuhnya untuk alur pelayanan rawat jalan (*outpatient*), konsultasi dokter anak & vaksinasi, serta tindakan estetika klinik.
3. **Integrasi Ekosistem Eksternal Lanjutan (Fase Roadmap Mendatang):**
   * Integrasi bridging ke sistem asuransi pihak ketiga atau jaringan platform eksternal tambahan dialokasikan pada fase pengembangan lanjutan setelah sistem operasional dasar beroperasi stabil.
4. **Permintaan Fitur di Luar Spesifikasi yang Disepakati:**
   * Modul atau alur kerja di luar 27 daftar fungsional operasional yang tertera pada proposal ini dapat didiskusikan kembali pada tahapan pengembangan berikutnya (*Change Request*).

---

## 4. PENDEKATAN TEKNIS & ARSITEKTUR (TECHNICAL APPROACH & STACK)

Sistem dirancang menggunakan arsitektur **Modular Monolith** yang menggabungkan performa eksekusi ultra-cepat, kemudahan pemeliharaan, serta kesederhanaan operasional server tanpa kerumitan infrastruktur yang berlebihan.

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
* **Backend Engine:** **Go (Golang)** — Performa komputasi sangat tinggi, konsumsi RAM sangat hemat, dan pemrosesan antrian konkuren (*goroutines*) yang terbukti tangguh menangani lonjakan pasien.
* **Frontend UI:** **Vue 3 (Vite + Composition API)** — Antarmuka web yang reaktif, ringan, cepat, dan nyaman dibuka dari berbagai perangkat (ponsel pasien, tablet dokter, maupun PC kasir).
* **Komunikasi Data Real-Time:** **Server-Sent Events (SSE) / WebSocket** — Memastikan monitor antrian pasien dan komputer apotek memperbarui status seketika tanpa perlu memuat ulang halaman (*zero refresh*).
* **Penyimpanan Data:** **Relational Database (SQLite / PostgreSQL)** — Transaksi atomik (*ACID*) untuk menjamin saldo stok obat dan pencatatan transaksi kasir selalu presisi tanpa risiko selisih.
* **Integrasi Layanan:**
  - **WhatsApp Service API:** Otomasi pengiriman nomor tiket antrian dan tautan live monitor langsung ke nomor ponsel pasien.
  - **Payment Gateway Adapter:** Otomasi verifikasi pembayaran QRIS dan transfer bank pendaftaran.

---

## 5. KEBUTUHAN NON-FUNGSIONAL (NON-FUNCTIONAL REQUIREMENTS)

1. **Keamanan & Perlindungan Privasi (Security & Privacy):**
   - Penerapan prinsip *Privacy-First*: Tidak meminta dan tidak menyimpan Nomor KTP/NIK pasien.
   - Penyamaran otomatis nama pasien di monitor ruang tunggu publik.
   - Enkripsi kata sandi akun staf menggunakan standar industri (*bcrypt*).
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
| **7 November 2026** | **Milestone 1: Discovery & Architecture Sign-Off** | Dokumen spesifikasi fungsional disepakati, arsitektur teknis dan alur operasional tuntas. |
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
| **Kendala koneksi internet di lokasi klinik** | Sedang | Sistem web dirancang ringan dengan aset lokal (*local asset caching*) sehingga kasir tetap dapat melayani transaksi. |
| **Salah ketik nomor WhatsApp oleh pasien saat daftar** | Rendah | Nomor antrian tetap dapat dilihat langsung pada layar konfirmasi web di HP pasien dan terpampang di monitor TV ruang tunggu klinik. |
| **Perubahan kebutuhan staf di tengah jalan (*scope creep*)** | Sedang | Mengacu ketat pada dokumen lingkup kerja operasional fungsional yang telah disepakati bersama. |

---

## 9. ESTIMASI BIAYA & SKEMA PEMBAYARAN (COST ESTIMATE & PRICING TERMS)

### 9.1. Rincian Paket Pengembangan Perangkat Lunak Operasional Klinik
Pengembangan sistem mencakup seluruh 5 modul operasional terpadu (Portal QR, Kasir GRO, Ruang Dokter, Farmasi & Stok Opnam, Modul Atasan):

* **Total Fee Paket Pengembangan Software Operasional Klinik:** Rp 12.000.000,- (Dua Belas Juta Rupiah)
* **Layanan Sudah Termasuk (*Included*):**
  - Pengembangan seluruh 27 fitur fungsional operasional klinik sesuai proposal.
  - Sesi presentasi demo langsung di klinik pada 16 November 2026.
  - Modul impor & migrasi data master obat lama (Excel/CSV).
  - Pelatihan (*training*) penggunaan sistem untuk seluruh staf klinik.
  - Garansi pemeliharaan (*bug fix & support*) selama 3 bulan pasca *go-live*.

### 9.2. Skema Termin Pembayaran (*Milestone Payment Schedule*)
Pembayaran dilakukan secara bertahap berbasis pencapaian tonggak kerja nyata:

| Termin Pembayaran | Persentase | Nominal | Waktu Pembayaran |
| :--- | :---: | :---: | :--- |
| **Termin 1 (Uang Muka / DP)** | **30%** | **Rp 3.600.000,-** | Saat proposal disetujui dan kontrak pengerjaan ditandatangani. |
| **Termin 2 (Pasca Demo Sukses)** | **40%** | **Rp 4.800.000,-** | Setelah sesi demo langsung berhasil dipresentasikan di klinik (16 Nov 2026). |
| **Termin 3 (Peluncuran / Final)** | **30%** | **Rp 3.600.000,-** | Setelah sistem selesai diuji coba (*UAT*) dan siap *go-live* (Pertengahan Des 2026). |

### 9.3. Ketentuan Masa Uji Coba & Tarif Pengembangan Lanjutan

Sistem diberikan **Masa Evaluasi Operasional (Trial Period) selama 6 (enam) bulan** pasca peluncuran resmi (*go-live*). Periode ini mencakup pemeliharaan teknis (*bug fix*) serta penyesuaian **maksimal 2 (dua) fitur minor** secara cuma-cuma (contoh: penyesuaian kolom input form, perapian template struk thermal, atau redaksi notifikasi WhatsApp).

Segala bentuk penambahan fitur baru di luar kesepakatan awal maupun implementasi sistem untuk cabang/unit baru akan diproses secara resmi (*Change Request*) dengan acuan tarif terstandarisasi berikut:

| Kategori Layanan Tambahan | Rentang Tarif | Estimasi Durasi | Cakupan Pekerjaan |
| :--- | :---: | :---: | :--- |
| **Masa Trial Operasional** | **Termasuk (Rp 0)** | 6 Bulan | Garansi teknis, pemeliharaan sistem, dan maksimal 2 penyesuaian fitur minor. |
| **Fitur Tambahan Skala Kecil** | **Rp 2.000.000 – Rp 3.500.000** | 2 – 3 Hari Kerja | Modifikasi fungsi pada modul yang ada tanpa pembuatan tabel database baru. |
| **Fitur Baru Skala Menengah** | **Rp 3.500.000 – Rp 5.000.000** | 4 – 7 Hari Kerja | Pengembangan fitur mandiri baru dengan tabel database dan antarmuka tersendiri. |
| **Kloning Sistem Cabang Baru** | **Rp 5.000.000 – Rp 8.000.000** | 5 – 10 Hari Kerja | *Re-skin* branding, setup database mandiri cabang, koneksi gateway, dan *deployment*. |

---

## 10. LEMBAR PERSETUJUAN (SIGN-OFF & APPROVAL)

Proposal ini diajukan untuk menjadi pedoman resmi pelaksanaan kerja sama pengembangan perangkat lunak operasional klinik antara Tim Pengembang dan Manajemen Klinik.

| Diajukan Oleh: | Disetujui Oleh: |
| :---: | :---: |
| <br><br><br>____________________________________<br>**Lead Project Manager**<br>Tim Pengembang Sistem Operasional | <br><br><br>____________________________________<br>**Pimpinan / Direktur Operasional**<br>Kidinest & Semarang Skin Center (SSC) |
| Tanggal: ........................................ | Tanggal: ........................................ |
