# Product Requirement Document (PRD)
## Smart Clinic Management System (Demo & MVP)

---

## 1. Executive Summary
Aplikasi manajemen klinik terpadu (Modular Monolith) berbasis **Go** dan **Vue 3** yang dirancang untuk mengatasi masalah operasional harian klinik swasta:
1. Antrian menumpuk dan chat WhatsApp overload pendaftaran malam hari.
2. Miskomunikasi resep manual antara dokter dan apoteker.
3. Kebocoran atau kekacauan stok obat tanpa pencatatan Batch & Expired Date.
4. Kesiapan regulasi integrasi **SATUSEHAT Kemenkes RI**.
5. Kebutuhan 1 invoice sederhana untuk pasien dengan pencatatan pembagian internal klinik.

Dokumen ini menjadi acuan pengembangan prototipe / aplikasi demo untuk presentasi ke klien sebelum ekspansi ke rilis penuh.

---

## 2. Problem Statement & User Pain Points
* **Pasien:** Antri lama tanpa kepastian nomor dan estimasi waktu; kerepotan mendaftar manual via WhatsApp yang lambat dibalas.
* **Resepsionis / Kasir:** Terbebani membaca ratusan chat WA semalam untuk nomor antrian; kebingungan mencocokkan pembayaran tunai vs digital.
* **Dokter:** Menulis resep kertas yang rawan salah baca oleh apoteker; tidak tahu apakah stok obat di apotek sedang kosong atau ada.
* **Apoteker:** Menunggu resep fisik diantar; risiko memberikan obat mendekati kadaluarsa (*Expired Date*); salah hitung stok fisik.
* **Manajemen Klinik:** Waspada terhadap regulasi wajib SATUSEHAT Kemenkes dan perlunya audit internal pendapatan tanpa membingungkan pasien.

---

## 3. Product Scope & Functional Modules

### Modul 1: Antrian Pasien Mandiri & Kasir (Queue & Payment)
* **Pasien Self-Service (Mobile Web / QR Kiosk):**
  * Pasien scan QR atau buka link klinik dari HP.
  * Input NIK, Nama Lengkap, Nomor HP, dan Keluhan Singkat.
  * Pilihan metode pembayaran biaya pendaftaran/konsultasi (contoh: Rp 35.000):
    * **QRIS (Mock Payment Gateway):** Simulasi bayar instan, langsung menerbitkan Nomor Antrian FIFO.
    * **Tunai / Cash:** Pasien mendapat status *Pending Kasir*, nomor antrian aktif setelah disetujui kasir.
  * Layar Live Queue: Nomor antrian pasien saat ini, estimasi antrian di depannya, status (Menunggu, Dipanggil, Selesai).
* **Layar Kasir / Resepsionis:**
  * Daftar antrian pending (pembayaran tunai).
  * Tombol 1-klik "Terima Tunai & Terbitkan Antrian".
  * Monitor antrian FIFO keseluruhan.

### Modul 2: Ruang Dokter & E-Resep (Doctor Consultation)
* **Daftar Pasien Menunggu:** Menampilkan urutan FIFO pasien yang sudah lunas bayar pendaftaran.
* **Pemeriksaan & Diagnosa:**
  * Input Diagnosa Dokter (dilengkapi saran kode ICD-10 untuk integrasi SATUSEHAT).
  * Catatan klinis / anamnesis singkat.
* **E-Prescription (Resep Elektronik):**
  * Pencarian obat dengan indikator stok tersedia (*available stock*).
  * Dokter menentukan jumlah & aturan pakai (signa).
  * Saat dokter menyimpan resep: Stok obat otomatis masuk status **Reserved** (tidak bisa diserobot dokter lain).
  * Notifikasi instan terkirim ke modul Apotek.

### Modul 3: Apotek & Reserved Stock Management
* **Antrian Resep Masuk:**
  * Layar live yang memperbarui otomatis ketika ada resep baru dari dokter (dilengkapi visual/audio chime cue).
  * Menampilkan detail pasien, dokter perujuk, dan daftar obat.
* **Mekanisme Reserved Stock (Shadow Deduction):**
  * Stok Fisik (`physical_stock`): Total stok nyata di apotek.
  * Stok Dipesan (`reserved_stock`): Total yang sedang ada di resep aktif dokter.
  * Stok Tersedia (`available_stock` = `physical_stock - reserved_stock`).
  * Saat apoteker menekan **"Selesai Racik / Serahkan Obat"**:
    * `physical_stock` berkurang sesuai jumlah obat.
    * `reserved_stock` dilepas.
  * Jika resep dibatalkan: `reserved_stock` dikurangi tanpa mengurangi `physical_stock`.

### Modul 4: Inventori & Gudang Apotek
* **Pencatatan 1 Gudang Aktif:** Fokus pada efisiensi stok apotek tunggal.
* **Manajemen Batch & FEFO (First Expired, First Out):**
  * Setiap batch obat wajib memiliki Nomor Batch dan Tanggal Kadaluarsa (Expired Date / ED).
  * Indikator warna kartu obat:
    * 🔴 **Merah:** ED < 3 bulan (peringatan segera retur/habiskan).
    * 🟡 **Kuning:** Stok total < batas stok minimum (*Low Stock Alert*).
    * 🟢 **Hijau:** Stok dan masa berlaku aman.
* **Fitur Tambah Stok (Restock):** Input faktur masuk, nomor batch baru, tanggal ED, jumlah botol/strip, dan harga beli.

### Modul 5: Billing & Invoicing (Pasien vs Internal)
* **Single Patient Invoice:**
  * Pasien hanya menerima 1 struk/PDF ringkas: Biaya Konsultasi Dokter + Biaya Tindakan + Total Biaya Obat.
  * Dilengkapi QR status lunas dan cetak struk ramah printer thermal.
* **Internal Breakdown (Laporan Manajemen):**
  * Rincian internal pembagian pendapatan: Jasa Medis Dokter, Biaya Sarana/Klinik, Margin Penjualan Farmasi, dan Pajak (PPN).

### Modul 6: Kesiapan SATUSEHAT Kemenkes RI
* **Indikator Kesiapan (Ready Badge):**
  * Status verifikasi NIK Pasien (mock IHS Patient ID: `P01293847561`).
  * Mapping Kode KFA (*Kamus Farmasi dan Alkes*) pada master obat.
  * Mapping Kode ICD-10 pada diagnosa dokter.
  * Log Aktivitas: JSON Payload FHIR R4 (`Encounter`, `Condition`, `MedicationRequest`) yang valid dan siap dikirim ke API SATUSEHAT saat production.

---

## 4. Tech Stack & Architecture

* **Arsitektur:** Modular Monolith (1 Backend Server, 1 Frontend SPA).
* **Backend:**
  * Bahasa: **Go (Golang)**
  * HTTP Router: **Chi** atau **Fiber**
  * Database: **SQLite** (menggunakan driver CGO-free/modernc atau GORM untuk kemudahan porting ke PostgreSQL di production).
* **Frontend:**
  * Framework: **Vue 3** (Composition API, `<script setup>`)
  * Build Tool: **Vite**
  * Styling: **Tailwind CSS**
  * State Management: **Pinia**
  * Icons: **Lucide Vue**

---

## 5. Non-Goals untuk Fase Demo
* Multi-cabang klinik yang kompleks (cukup 1 klinik & 1 apotek).
* Integrasi langsung live webhook dengan bank payment gateway (menggunakan interactive simulator).
* Request live SSL/OAuth handshake ke server production Kemenkes DTO (menggunakan simulator FHIR payload & IHS verify mock).
