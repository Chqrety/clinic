# Plan: Digitalisasi Pembayaran On-The-Spot (EDC Fisik & QRIS Meja via Trace Entry)

## 1. Latar Belakang & Masalah
* Klinik menggunakan **Mesin EDC Fisik (BCA / BRI)** dan **Papan Akrilik QRIS Statis Meja** untuk transaksi *on-the-spot* di kasir agar terbebas dari potongan biaya gateway online dan dana langsung masuk ke rekening giro bank klinik.
* Mesin EDC fisik dan papan QRIS meja berdiri sendiri (*offline/standalone*), tidak memiliki kabel atau koneksi API langsung ke browser sistem kasir.
* Diperlukan mekanisme pencatatan digital yang **cepat (2-3 detik)**, **anti-salah catat**, dan memiliki **kode unik** (*Trace Number* / *Approval Code* / *Reference ID*) untuk memudahkan rekonsiliasi harian dan audit pembukuan klinik.

---

## 2. Cakupan Perubahan (Scope of Work)

### A. Frontend Kasir (`CashierView.vue`)
1. **Opsi Metode Pembayaran 4 Saluran:**
   * `CASH` (Tunai dengan kalkulator kembalian)
   * `EDC` (Mesin EDC Fisik Bank)
   * `QRIS` (Papan Akrilik QRIS Meja / E-Wallet)
   * `BANK_TRANSFER` (Transfer Bank Manual)
2. **Sub-Pilihan Bank & Input Unik Otomatis:**
   * Jika **EDC**: Tombol cepat Bank (`BCA`, `BRI`, `Mandiri`), input field `Trace Number / App Code (6 digit slip kertas EDC)`.
   * Jika **QRIS Meja**: Tombol cepat Provider (`BCA QRIS`, `GoPay / Shopee`, `Mandiri`), input field `No. Referensi Bukti HP Pasien`.
   * Jika **Transfer**: Input field `Nama Pengirim & No. Ref Transfer`.
3. **Format Penyimpanan Referensi:**
   * Format terstandarisasi otomatis: `EDC-BCA-000124` atau `QRIS-BCA-98421`.
4. **Visibilitas Audit & Cetak Struk:**
   * Tabel Riwayat Invoice menampilkan badge metode pembayaran dan nomor Trace unik.
   * Struk cetak kasir (slip thermal 80mm) mencantumkan baris `Ref/Trace: <No. Ref>`.
   * Tab Rekap Shift Kasir mencantumkan rincian transaksi EDC & QRIS untuk pencocokan (*matching*) dengan kertas *Batch Settlement* EDC.

### B. Backend Verifikasi & Database
* Kolom `payment_reference` sudah ada di tabel `invoices` dan struct `PayBillRequest`.
* Memastikan `PaymentMethodQRIS` ("QRIS") diproses dengan benar tanpa error constraint.

---

## 3. Rencana Eksekusi
1. Update `CashierView.vue`:
   - State `selectedBank` dan `traceNumberInput`.
   - Selector 4 metode pembayaran di panel kasir kanan.
   - Kolom tabel riwayat invoice dan struk cetak 80mm.
2. Verifikasi kompilasi frontend (`npm run build`).
3. Verifikasi backend test (`go test ./...`).
