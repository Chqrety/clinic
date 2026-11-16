# Panduan Integrasi: Tripay Payment Gateway & WhatsApp Gateway (Fonnte)

Dokumen ini berisi panduan langkah demi langkah cara membuat akun, mengambil kredensial API, dan menghubungkan gateway pembayaran **Tripay (QRIS)** serta gateway notifikasi **WhatsApp (Fonnte)** ke backend Smart Clinic Management System (SCMS).

---

## 1. Integrasi Tripay (Payment Gateway QRIS)

Tripay digunakan untuk pembayaran non-tunai otomatis (khususnya **QRIS Dinamis / Closed Payment**) di Kiosk Mandiri dan Kasir.

### Langkah 1: Registrasi Akun Sandbox Tripay (Tahap Pengujian)
1. Buka browser dan kunjungi: **[https://tripay.co.id/simulator](https://tripay.co.id/simulator)** atau langsung daftar di **[https://tripay.co.id/member](https://tripay.co.id/member)**.
2. Untuk pengujian (development), buka portal Sandbox: **[https://tripay.co.id/member-sandbox](https://tripay.co.id/member-sandbox)** (atau klik toggle **Sandbox** di pojok kanan atas dashboard).
3. Buat akun baru dengan email dan password Anda (di Sandbox tidak perlu verifikasi KTP/NPWP).

### Langkah 2: Buat Merchant & Ambil API Key
1. Masuk ke dashboard Sandbox: **[https://tripay.co.id/member-sandbox](https://tripay.co.id/member-sandbox)**.
2. Di menu sidebar kiri, pilih **Pengaturan** > **Merchant**.
3. Klik tombol **Tambah Merchant**:
   * **Nama Merchant:** `Klinik Utama Semarang` (atau nama klinik Anda).
   * **Website / URL:** `http://localhost:5173` (atau domain klinik Anda).
   * **Whitelist IP:** Kosongkan terlebih dahulu saat testing lokal.
4. Setelah merchant berhasil dibuat, klik nama merchant tersebut, lalu buka tab **API & Webhook**:
   * **Kode Merchant (Merchant Code):** Formatnya huruf T diikuti angka (contoh: `T12345`).
   * **API Key:** String acak panjang (contoh: `DEV-abcde12345...`).
   * **Private Key:** Kunci privat untuk signature HMAC-SHA256 (contoh: `xyz98765...`).
   * **URL Callback / Webhook:** Isi dengan endpoint webhook backend Anda:
     ```text
     https://<domain-kamu-atau-ngrok>/api/webhooks/tripay
     ```
     *(Catatan: Jika testing di localhost, Tripay tidak bisa menembus `localhost:8080`. Gunakan ngrok seperti `https://abcd-123.ngrok-free.app/api/webhooks/tripay`)*.

### Langkah 3: Pasang Kredensial di `backend/.env`
Buka file `backend/.env` pada repositori ini dan isi nilainya:

```ini
# Tripay Payment Gateway (Closed QRIS)
TRIPAY_API_KEY=DEV-abcde12345xxxxxxxxxxxxxxxx
TRIPAY_PRIVATE_KEY=xyz98765xxxxxxxxxxxxxxxxxxxxxx
TRIPAY_MERCHANT_CODE=T12345
TRIPAY_IS_SANDBOX=true
TRIPAY_CALLBACK_URL=https://your-ngrok-subdomain.ngrok-free.app/api/webhooks/tripay
TRIPAY_RETURN_URL=http://localhost:5173/kiosk
```

### Langkah 4: Cara Testing Simulasi Pembayaran
1. Pasien memilih pembayaran **QRIS** di Kiosk atau Kasir.
2. Sistem akan memanggil Tripay API dan menampilkan gambar QRIS resmi dari Tripay.
3. Untuk menyimulasikan pembayaran tanpa uang asli:
   * Buka simulator Tripay: **[https://tripay.co.id/simulator](https://tripay.co.id/simulator)**.
   * Masukkan **Nomor Referensi Transaksi** yang muncul di layar.
   * Klik tombol **Bayar Sekarang**.
   * Webhook Tripay akan mengirim notifikasi ke server backend Anda, mengubah invoice menjadi **LUNAS**, dan memutar audio chime di kasir.

> **Transisi ke Production (Live):**
> Jika sudah siap beroperasi nyata, ajukan verifikasi merchant di dashboard Live Tripay (`https://tripay.co.id/member`), ubah `TRIPAY_IS_SANDBOX=false`, lalu ganti ketiga key di atas dengan key dari dashboard Live.

---

## 2. Integrasi WhatsApp Gateway (Fonnte)

Sistem menggunakan provider **Fonnte** untuk mengirim tiket antrian paperless, notifikasi nomor antrian dipanggil, dan resep siap ambil ke nomor WhatsApp pasien secara instan.

### Langkah 1: Registrasi Akun Fonnte
1. Kunjungi situs resmi: **[https://fonnte.com](https://fonnte.com)**.
2. Klik tombol **Daftar** / **Register** dan selesaikan verifikasi email.
3. Masuk ke dashboard Fonnte: **[https://md.fonnte.com](https://md.fonnte.com)**.

### Langkah 2: Hubungkan Nomor WhatsApp Klinik (Scan QR)
1. Di dashboard Fonnte, buka menu **Device**.
2. Klik tombol **Tambah Device** / **Add Device**:
   * Masukkan nama device (contoh: `WA Klinik Semarang`).
3. Klik tombol **Connect** / icon QR Code pada device yang baru dibuat.
4. Buka aplikasi WhatsApp di HP klinik (bisa WA reguler atau WA Business):
   * Ketuk menu titik tiga (Android) atau Pengaturan (iOS) > **Perangkat Tertaut (Linked Devices)**.
   * Pilih **Tautkan Perangkat** dan scan QR Code yang tampil di layar Fonnte.
5. Tunggu hingga status device di Fonnte berubah menjadi **Connected** (hijau).

### Langkah 3: Ambil API Token Fonnte
1. Masih di menu **Device** pada device yang sudah *Connected*.
2. Lihat kolom **Token** atau klik tombol **Salin Token**.
3. Format token Fonnte adalah string kombinasi angka dan huruf (contoh: `aB9#kL2@mZ19xQ...`).

### Langkah 4: Pasang Kredensial di `backend/.env`
Buka file `backend/.env` dan perbarui baris berikut:

```ini
# WhatsApp Notification Gateway (Paperless Queue Ticket)
WA_GATEWAY_PROVIDER=fonnte
WA_GATEWAY_API_URL=https://api.fonnte.com/send
WA_GATEWAY_API_TOKEN=aB9#kL2@mZ19xQ_token_fonnte_kamu
```

### Langkah 5: Cara Verifikasi Pengiriman WhatsApp
1. Buka Kiosk Mandiri (`/kiosk`) atau Meja Depan Kasir (`/admin#cashier`).
2. Daftarkan pasien baru dengan nomor HP yang valid (format `08...`, `628...`, atau `+628...` otomatis distandarisasi oleh backend).
3. Setelah tiket dibuat:
   * Backend akan langsung menembak API Fonnte.
   * Pasien menerima pesan WhatsApp otomatis berisi rincian:
     * Nama Klinik & Unit Praktik (Kidinest / Semarang Skin Center).
     * Nomor Antrian (contoh: `A-01`).
     * Jumlah antrian di depan pasien.
     * Tautan live status monitor antrian.
4. Jika `WA_GATEWAY_PROVIDER=mock`, pesan tidak akan dikirim ke WhatsApp fisik melainkan hanya tercatat di terminal/console backend sebagai log debug.

---

## 3. Rangkuman File Konfigurasi Lengkap (`backend/.env`)

Berikut adalah template lengkap file `backend/.env` siap pakai:

```ini
# Database Configuration (PostgreSQL 16)
DB_HOST=localhost
DB_PORT=5432
DB_USER=clinic
DB_PASSWORD=clinic_secret
DB_NAME=clinic_db
DB_SSLMODE=disable

# Application Server
PORT=8080
ENV=development
JWT_SECRET=clinic-dev-secret-key-at-least-32-chars-long

# -------------------------------------------------------------
# 1. TRIPAY PAYMENT GATEWAY (QRIS Closed Payment)
# -------------------------------------------------------------
# Dapatkan dari: https://tripay.co.id/member-sandbox (Tab API & Webhook)
TRIPAY_API_KEY=DEV-xxxxxxxxx
TRIPAY_PRIVATE_KEY=yyyyyyyyy
TRIPAY_MERCHANT_CODE=T12345
TRIPAY_IS_SANDBOX=true
TRIPAY_CALLBACK_URL=https://your-tunnel.ngrok-free.app/api/webhooks/tripay
TRIPAY_RETURN_URL=http://localhost:5173/kiosk

# -------------------------------------------------------------
# 2. WHATSAPP GATEWAY (Fonnte API)
# -------------------------------------------------------------
# Dapatkan dari: https://md.fonnte.com (Menu Device > Salin Token)
WA_GATEWAY_PROVIDER=fonnte
WA_GATEWAY_API_URL=https://api.fonnte.com/send
WA_GATEWAY_API_TOKEN=your_fonnte_token_here
```

---

## 4. Troubleshooting Cepat

| Kendala | Penyebab Umum | Solusi |
| :--- | :--- | :--- |
| **QRIS Tripay gagal dibuat (Error 401)** | `API_KEY` atau `MERCHANT_CODE` salah ketik. | Cek kembali kesesuaian antara environment (`IS_SANDBOX=true` vs dashboard sandbox). |
| **Status tidak berubah lunas otomatis saat test** | Webhook Tripay tidak bisa mengakses URL localhost. | Gunakan ngrok (`ngrok http 8080`), lalu pasang URL ngrok tersebut pada `TRIPAY_CALLBACK_URL`. |
| **WA tidak masuk ke pasien** | Device di Fonnte disconnect atau token salah. | Buka dashboard Fonnte, pastikan status device berwarna hijau (*Connected*), dan pastikan kuota pesan Fonnte masih tersedia. |
| **Nomor HP pasien format luar/salah** | Format acak (misal `0812-3456-7890`). | Backend sudah dilengkapi fungsi normalisasi otomatis menjadi `6281234567890`. |
