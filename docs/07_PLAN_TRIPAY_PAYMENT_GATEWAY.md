# 07. Rencana Integrasi Tripay Payment Gateway (QRIS & Webhook Callback)

## 1. Tujuan & Latar Belakang
Mengganti simulasi bayar instan pada pendaftaran mandiri (Kidinest) dengan integrasi nyata **Tripay Payment Gateway (Closed Transaction - QRIS)**:
* Pasien menerima QR Code QRIS dinamis resmi dari Tripay.
* Saat pasien membayar lewat e-wallet/mobile banking, Tripay mengirimkan notifikasi *Webhook Callback* ke server klinik.
* Server memvalidasi signature HMAC-SHA256, memperbarui status menjadi `PAID`, dan langsung menerbitkan Nomor Antrian resmi FIFO.

---

## 2. Invarian Domain & Keamanan
1. **Validasi Signature Wajib:** Setiap request callback dari Tripay wajib divalidasi keasliannya menggunakan `HMAC-SHA256(private_key, raw_json_body)`. Request dengan signature tidak valid langsung ditolak (`401 Unauthorized`).
2. **Idempotency Webhook:** Jika Tripay mengirimkan callback berulang untuk transaksi yang sama, server tidak boleh menerbitkan tiket antrian ganda.
3. **Graceful Dev/Simulator Fallback:** Jika `TRIPAY_API_KEY` tidak diatur di environment (misal saat testing offline), sistem secara cerdas beralih ke *Simulated Mode* agar alur pengujian tidak terblokir.

---

## 3. Langkah Implementasi Teknis

### Tahap 1: Konfigurasi Environment & Model Data
* Tambahkan variabel konfigurasi di backend (`.env`):
  ```env
  TRIPAY_API_KEY=your_api_key
  TRIPAY_PRIVATE_KEY=your_private_key
  TRIPAY_MERCHANT_CODE=your_merchant_code
  TRIPAY_IS_SANDBOX=true
  TRIPAY_CALLBACK_URL=https://your-domain.com/api/webhooks/tripay
  ```
* Pastikan tabel `queues` memiliki kolom:
  * `payment_reference VARCHAR(100)` (kode referensi transaksi Tripay, contoh: `T123456789`).
  * `payment_url TEXT` (URL QR Code / payload QR string dari Tripay).

### Tahap 2: Backend Tripay Adapter (`backend/internal/shared/tripay`)
* Buat package `tripay` dengan fungsionalitas:
  1. `CreateClosedTransaction(req CreateTransactionRequest) (*TransactionResponse, error)`:
     - URL Sandbox: `https://tripay.co.id/api-sandbox/transaction/create`
     - URL Production: `https://tripay.co.id/api/transaction/create`
     - Payload: `method = "QRIS"`, `amount`, `merchant_ref`, `customer_name`, `customer_phone`, `order_items`.
     - Request Header: `Authorization: Bearer <API_KEY>`.
     - Signature Generator: `HMAC-SHA256(merchant_code + merchant_ref + amount, private_key)`.
  2. `ValidateCallbackSignature(rawBody []byte, callbackSignature string) bool`:
     - Menghitung `HMAC-SHA256(rawBody, private_key)` dan mencocokkannya dengan header `X-Callback-Signature`.

### Tahap 3: Endpoint Webhook Callback
* Buat route publik: `POST /api/webhooks/tripay`
* Alur eksekusi di handler:
  1. Baca raw body dan header `X-Callback-Signature`.
  2. Verifikasi signature via adapter Tripay.
  3. Parse status transaksi: jika `status == "PAID"`, panggil `queueUsecase.VerifyPayment(ctx, reference, "TRIPAY")`.
  4. Terbitkan nomor antrian FIFO secara atomik (`SELECT ... FOR UPDATE`).
  5. Kembalikan response JSON `{"success": true}` ke Tripay.

### Tahap 4: Integrasi Frontend Pendaftaran
* Di [LandingKidinest.vue](file:///home/alpi/Projects/clinic/frontend/src/views/LandingKidinest.vue):
  * Saat pasien memilih bayar QRIS, tampilkan dialog/tampilan modal berisi gambar QR Code dinamis Tripay dan batas waktu pembayaran (*countdown timer*).
  * Pasang interval polling singkat (`GET /api/public/queue/status/:id` setiap 3 detik) atau dengarkan SSE event agar saat webhook masuk, layar pasien seketika beralih ke tiket antrian resmi.

---

## 4. Rencana Pengujian (Verification)
1. **Unit Test:**
   - Test kalkulasi signature HMAC-SHA256 dengan fixture resmi Tripay.
   - Test handling callback `PAID` dan penolakan callback ber-signature palsu.
2. **Integration Test (Sandbox):**
   - Buat transaksi QRIS uji coba ke Tripay Sandbox.
   - Trigger simulasi bayar lewat Tripay Sandbox Merchant Dashboard.
   - Verifikasi nomor antrian terbit otomatis di database PostgreSQL.
