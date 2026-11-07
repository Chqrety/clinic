# Plan: Patient-Facing Landing Page Polish
## Website Resmi Pelayanan Pasien — Klinik Sehat Modern

Dokumen ini memuat rencana kerja pemolesan [LandingPage.vue](file:///home/alpi/Projects/clinic/frontend/src/views/LandingPage.vue) sebagai **website resmi klinik yang ditujukan untuk pasien dan masyarakat umum** yang ingin berobat (bukan halaman penjualan B2B software, karena penawaran ke pemilik klinik menggunakan PPT terpisah).

---

## 1. Tujuan Utama (*Core Objectives*)

1. **Memberikan Pengalaman Berobat yang Modern & Transparan bagi Pasien:**
   - Menghilangkan stres antri: Pasien bisa memantau antrian aktif dari rumah sebelum berangkat ke klinik.
   - Pendaftaran mandiri instan via web/ponsel tanpa perlu mengantri panjang di depan loket.
2. **Bahasa Ramah Pasien (Empatik & Lugas):**
   - Menghapus semua istilah programmer/IT teknis (*FHIR*, *KFA*, *disposisi*).
   - Menggunakan bahasa sederhana: jam buka, nama dokter jaga, cara ambil karcis, cara tebus obat, dan lokasi klinik.
3. **Membuktikan Kualitas Sistem kepada Calon Klien:**
   - Saat pemilik klinik melihat website ini saat demo, mereka langsung terbayang betapa profesionalnya citra klinik mereka di mata pasien.

---

## 2. Rincian Desain & Konten per Bagian

### A. Hero Section (Kenyamanan Berobat)
* **Badge / Pill:**  
  `Klinik Sehat Modern • Layanan Rawat Jalan & Farmasi Terpadu`
* **Headline Utama (Brandlos Editorial):**  
  > *"Berobat lebih tenang, pantau antrian langsung dari ponsel."*
* **Deskripsi Emosional & Praktis:**  
  Kesehatan Anda adalah prioritas kami. Nikmati pelayanan dokter umum yang ramah tanpa perlu berdesakan di ruang tunggu. Cek nomor antrian yang sedang dipanggil secara live, ambil tiket karcis mandiri, dan tebus obat di apotek tanpa repot membawa kertas resep.
* **Call to Action (CTA Pasien):**  
  Tombol mencolok: `[ Ambil Antrian Mandiri Sekarang → ]`

---

### B. Widget Live Antrian Klinik (Real-Time Monitor)
Fitur interaktif yang memberi kepastian waktu bagi pasien:
1. **Nomor yang Sedang Diperiksa:**  
   Menampilkan nomor berjalan (contoh: `A-001 - dr. Hendra Pratama, Sp.PD`).
2. **Pasien Menunggu:**  
   Jumlah antrian di ruang tunggu (contoh: `2 Pasien Menunggu • Estimasi ±8 menit/pasien`).
3. **Status Apotek:**  
   `Farmasi Buka • Obat Disiapkan Langsung Selesai Periksa`.
4. **Jam Operasional & Izin Resmi:**  
   `Buka Setiap Hari: 08.00 – 21.00 WIB • Izin Faskes: 445/0912/DINKES/2024`.

---

### C. Tiga Kemudahan Pasien (Patient Journey)

#### 01 / Ambil Antrian dari Mana Saja
* **Poin:** Pasien cukup mengisi data ringkas (NIK & nama), memilih metode pembayaran (QRIS instan atau tunai saat tiba), dan nomor karcis langsung aktif.

#### 02 / E-Resep Terhubung Langsung ke Apotek
* **Poin:** Tidak ada lagi resep kertas yang kusut atau tulisan tangan yang sulit dibaca. Begitu dokter selesai memeriksa, resep sudah terkirim ke meja apoteker untuk segera disiapkan.

#### 03 / Struk Biaya Transparan & Resmi
* **Poin:** Pasien memegang struk pembayaran kasir 80mm yang jelas memisahkan biaya karcis pendaftaran dan rincian obat yang ditebus.

---

### D. Informasi Kunjungan & Footer (*Visit Block*)
* **Alamat:** Jl. Kesehatan No. 88, Jakarta Pusat (Dekat Halte & Area Parkir Luas).
* **Kontak Darurat & WhatsApp:** (021) 555-0199 / WA: 0812-8877-6655.
* **Akses Staf Medis:** Tombol minimalis di pojok navigasi atau footer: `[ Portal Staf Medis ]` agar dokter dan apoteker tetap bisa masuk ke ruang kerja mereka.

---

## 3. Checklist Eksekusi (*Execution Steps*)

- [ ] Edit file `frontend/src/views/LandingPage.vue`.
- [ ] Ubah teks copy menjadi sudut pandang kenyamanan pasien berobat.
- [ ] Pertahankan integrasi live query `api.getQueues()` untuk nomor panggilan berjalan.
- [ ] Verifikasi tombol CTA mengarahkan pasien langsung ke Kiosk pendaftaran.
- [ ] Jalankan `npm --prefix frontend run build` untuk memverifikasi kelulusan build.
