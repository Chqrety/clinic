# Panduan Setup Midtrans Sandbox (QRIS & Payment Gateway)

Integrasi Midtrans telah siap di backend dengan dukungan:
- **Direct Dynamic QRIS (Core API `/v2/charge`):** Otomatis menghasilkan gambar QRIS GoPay/ShopeePay/BCA untuk scan di Kiosk/Portal Pasien.
- **Snap Token & Redirect URL:** Dukungan popup/checkout Midtrans Snap.
- **Webhook Web Callback (`/api/webhooks/midtrans`):** Idempotent, verifikasi signature SHA-512, auto-generate nomor antrean FIFO, dan kirim tiket via WhatsApp.

---

### 1. Ambil Kredensial di Midtrans Sandbox

1. Buka [https://dashboard.sandbox.midtrans.com](https://dashboard.sandbox.midtrans.com) dan login.
2. Pastikan badge di pojok kiri atas bertuliskan **SANDBOX** (bukan Production).
3. Buka menu **Settings** $\rightarrow$ **Access Keys**.
4. Salin:
   - **Server Key** (contoh: `SB-Mid-server-xxxxxxxxx`)
   - **Client Key** (contoh: `SB-Mid-client-xxxxxxxxx`)

---

### 2. Masukkan ke `.env` Backend

Buka file [backend/.env](file:///home/alpi/Projects/clinic/backend/.env), isi variabel berikut:

```env
# Midtrans Payment Gateway (Core API QRIS & Snap)
MIDTRANS_SERVER_KEY=SB-Mid-server-xxxxxxxxx
MIDTRANS_CLIENT_KEY=SB-Mid-client-xxxxxxxxx
MIDTRANS_IS_PRODUCTION=false
```

---

### 3. Setting Webhook / Notification URL di Midtrans

Agar status pembayaran otomatis ter-update dan tiket antrean terbit saat pasien scan QR:

1. Di Dashboard Midtrans Sandbox, buka **Settings** $\rightarrow$ **Configuration**.
2. Pada field **Payment Notification URL**, isi URL publik backend:
   ```text
   https://<domain-kamu-atau-ngrok>/api/webhooks/midtrans
   ```
   *(Untuk pengujian lokal di mesin dev, gunakan ngrok: `ngrok http 8080`, lalu pasang URL ngrok tersebut).*
3. Centang opsi notifikasi yang diinginkan dan klik **Update**.

---

### 4. Menguji Pembayaran via Simulator Midtrans

1. Buka Portal Pasien di browser dan daftarkan antrean dengan metode pembayaran **QRIS**.
2. Salin kode QR string atau URL QRIS.
3. Buka [Midtrans QRIS Simulator](https://simulator.sandbox.midtrans.com/qris/index).
4. Masukkan QR string atau scan QR code, lalu klik **Pay**.
5. Backend akan menerima notifikasi webhook, memvalidasi signature SHA-512, menerbitkan nomor antrian resmi (misal `A-001`), dan mengirim tiket via WhatsApp.
