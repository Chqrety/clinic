# MASTER PROMPT FULL VERSION
# INTEGRATED ACCOUNTING & FINANCIAL MANAGEMENT SYSTEM FOR CLINIC
## Backend Golang — PostgreSQL — Automatic Accounting — Inventory — Payment Gateway

---

# 1. ROLE DAN RESPONSIBILITY

Bertindak sebagai:

- Senior Golang Backend Engineer.
- Software Architect.
- Database Architect.
- Accounting System Analyst.
- ERP Developer.
- Financial Reporting Specialist.
- Inventory Management System Developer.
- Payment Gateway Integration Engineer.

Saya ingin membangun backend sistem akuntansi terintegrasi untuk operasional klinik menggunakan Golang.

Sistem harus dirancang seperti aplikasi ERP akuntansi profesional yang menghubungkan seluruh transaksi operasional klinik dengan proses pencatatan akuntansi dan laporan keuangan.

Sistem tidak boleh bergantung pada input laporan keuangan manual atau template Excel tertentu.

Semua laporan harus dihasilkan otomatis dari transaksi yang tersimpan di database melalui Accounting Engine dan General Ledger.

Gunakan prinsip Double-Entry Bookkeeping, Accrual Basis Accounting, Audit Trail, Financial Period Closing, dan Automatic Journal Posting.

Prioritaskan:
1. Akurasi perhitungan akuntansi.
2. Konsistensi data transaksi.
3. Integrasi seluruh modul.
4. Keamanan transaksi keuangan.
5. Otomatisasi laporan.
6. Kemudahan pengembangan.
7. Kemudahan audit.
8. Skalabilitas sistem.

Jangan langsung membuat seluruh source code sekaligus. Lakukan analisis, desain arsitektur, ERD, accounting flow, dan database schema terlebih dahulu.

---

# 2. TECHNOLOGY STACK

Gunakan teknologi berikut:

Backend:
- Golang.
- Gin atau Chi.
- REST API.
- Clean Architecture.
- Modular Monolith.

Database:
- PostgreSQL.
- GORM atau SQLC.
- golang-migrate.

Authentication:
- JWT Access Token.
- Refresh Token.
- RBAC.
- Password Hashing.

Infrastructure:
- Redis.
- Asynq untuk background jobs.
- Docker.
- Docker Compose.

Documentation:
- Swagger / OpenAPI.

Logging:
- Zap atau Zerolog.

Validation:
- go-playground/validator.

Testing:
- Go Testing.
- Testcontainers.
- Integration Testing.
- Unit Testing.

Export:
- PDF.
- Excel.
- CSV.

Gunakan struktur kode yang mudah dipelihara dan tidak mencampurkan business logic dengan HTTP Handler.

Gunakan interface untuk integrasi eksternal.

Jangan menggunakan float32 atau float64 untuk perhitungan uang. Gunakan PostgreSQL NUMERIC/DECIMAL dan tipe data decimal yang sesuai di Golang.

---

# 3. SYSTEM ARCHITECTURE

Gunakan pendekatan Modular Monolith dengan pembagian domain berikut:

1. Authentication Module.
2. User Management Module.
3. Branch Management Module.
4. Master Data Module.
5. Inventory Module.
6. Purchase Module.
7. Sales Module.
8. Payment Module.
9. Cash Management Module.
10. Bank Management Module.
11. Payment Gateway Module.
12. Account Receivable Module.
13. Account Payable Module.
14. Accounting Module.
15. Journal Management Module.
16. General Ledger Module.
17. Financial Reporting Module.
18. Fixed Asset Module.
19. Tax Management Module.
20. Audit Trail Module.
21. Financial Period Module.
22. Notification Module.
23. Export and Document Module.

Setiap modul memiliki:
- Domain.
- Application.
- Repository.
- Infrastructure.
- Delivery / Handler.

Buat dependency flow yang jelas.

Modul operasional tidak boleh melakukan manipulasi langsung terhadap tabel General Ledger. Semua pencatatan harus melalui Accounting Engine.

---

# 4. MASTER DATA MANAGEMENT

Buat master data berikut:

## 4.1 User Management
- User.
- Role.
- Permission.
- User Role.
- Branch Access.
- User Status.
- Login History.

## 4.2 Branch Management
- Branch Code.
- Branch Name.
- Address.
- Contact.
- Tax Identity jika relevan.
- Branch Status.

Sistem harus mendukung beberapa cabang klinik.

Setiap transaksi harus memiliki branch_id jika transaksi terkait dengan cabang.

## 4.3 Master Patient
- Patient Number.
- Patient Name.
- Contact.
- Address.
- Patient Category.
- Insurance Reference jika relevan.

## 4.4 Master Customer
- Customer Code.
- Customer Name.
- Customer Type.
- Customer Address.
- Payment Terms.

## 4.5 Master Supplier
- Supplier Code.
- Supplier Name.
- Supplier Type.
- Contact.
- Address.
- Payment Terms.
- Tax Information.

## 4.6 Master Product

Kategori:
- Obat.
- Alat Kesehatan.
- Bahan Medis Habis Pakai.
- Barang Operasional.
- Barang Non-Medis.

Field:
- Product Code.
- Product Name.
- SKU.
- Category.
- Unit.
- Purchase Price.
- Selling Price.
- Minimum Stock.
- Barcode.
- Active Status.

Dukung satuan dan konversi satuan.

---

# 5. ACCOUNTING CORE ENGINE

Accounting Engine adalah pusat utama sistem.

Semua transaksi keuangan harus melalui engine ini.

Buat sistem yang mendukung:

- Automatic Journal Generation.
- Double Entry.
- Account Mapping.
- Journal Validation.
- Journal Posting.
- Journal Reversal.
- Journal Adjustment.
- General Ledger.
- Trial Balance.
- Period Closing.
- Audit Trail.

## 5.1 Chart of Accounts (COA)

Buat Chart of Accounts fleksibel dan configurable.

Kelompok akun:

### A. ASSETS

Aset Lancar:
- Kas Kecil.
- Kasir.
- Kas Operasional.
- Bank BCA.
- Bank Mandiri.
- Bank Lainnya.
- Clearing Payment Gateway.
- Piutang Pasien.
- Piutang Perusahaan.
- Piutang Asuransi.
- Piutang BPJS jika relevan.
- Persediaan Obat.
- Persediaan Alat Kesehatan.
- Persediaan BMHP.
- Uang Muka Pembelian.
- Pajak Dibayar Dimuka.

Aset Tidak Lancar:
- Tanah.
- Bangunan.
- Peralatan Medis.
- Peralatan Kantor.
- Kendaraan.
- Komputer.
- Akumulasi Penyusutan.
- Aset Tetap Lainnya.

### B. LIABILITIES

Liabilitas Jangka Pendek:
- Hutang Supplier Obat.
- Hutang Supplier Alat Medis.
- Hutang Dokter.
- Hutang Gaji.
- Hutang Pajak.
- Hutang Operasional.
- Uang Muka Pasien.
- Pendapatan Diterima Dimuka.
- Hutang Lainnya.

Liabilitas Jangka Panjang:
- Pinjaman Bank.
- Pinjaman Jangka Panjang.
- Kewajiban Jangka Panjang Lainnya.

### C. EQUITY

- Modal Disetor.
- Tambahan Modal.
- Laba Ditahan.
- Laba Tahun Berjalan.
- Prive atau Distribusi Pemilik.

### D. REVENUE

- Pendapatan Konsultasi.
- Pendapatan Pemeriksaan.
- Pendapatan Tindakan Medis.
- Pendapatan Laboratorium.
- Pendapatan Penjualan Obat.
- Pendapatan Penjualan Alat Kesehatan.
- Pendapatan Medical Check Up.
- Pendapatan Perawatan.
- Pendapatan Lainnya.

### E. COST OF GOODS SOLD

- HPP Obat.
- HPP Alat Kesehatan Habis Pakai.
- HPP BMHP.

### F. EXPENSES

- Beban Gaji.
- Beban Jasa Dokter.
- Beban Sewa.
- Beban Listrik.
- Beban Air.
- Beban Internet.
- Beban Administrasi Payment Gateway.
- Beban Administrasi Bank.
- Beban Pemasaran.
- Beban Perlengkapan.
- Beban Pemeliharaan.
- Beban Penyusutan.
- Beban Pajak.
- Beban Operasional.
- Beban Lainnya.

COA harus mendukung:
- Parent Account.
- Child Account.
- Account Code.
- Account Name.
- Account Type.
- Normal Balance.
- Is Active.
- Is System Account.
- Allow Manual Journal.
- Branch Scope.
- Created At.
- Updated At.

Jangan hardcode seluruh akun ke dalam business logic.

---

# 6. JOURNAL MANAGEMENT

Buat modul khusus pengelolaan jurnal akuntansi.

Modul ini sangat penting karena menjadi sumber pencatatan resmi untuk General Ledger.

Jenis jurnal yang wajib tersedia:

1. General Journal / Jurnal Umum.
2. Cash Receipt Journal / Jurnal Kas Masuk.
3. Cash Disbursement Journal / Jurnal Kas Keluar.
4. Bank Journal / Jurnal Bank.
5. Payment Gateway Journal.
6. Purchase Journal / Jurnal Pembelian.
7. Sales Journal / Jurnal Penjualan.
8. Adjusting Journal / Jurnal Penyesuaian.
9. Closing Journal / Jurnal Penutup.
10. Reversal Journal / Jurnal Pembalik.

## 6.1 General Journal

Digunakan untuk seluruh pencatatan debit dan kredit.

Field:
- Journal ID.
- Journal Number.
- Journal Date.
- Reference Number.
- Source Module.
- Source Transaction ID.
- Description.
- Total Debit.
- Total Credit.
- Status.
- Created By.
- Approved By.
- Posted At.

Journal Lines:
- Account ID.
- Description.
- Debit.
- Credit.
- Branch ID.
- Cost Center.
- Department.

Validasi:
- Debit harus sama dengan kredit.
- Journal Number harus unik.
- Tidak boleh ada jurnal duplikat.
- Jurnal posted tidak boleh diedit langsung.
- Koreksi menggunakan reversal.
- Periode yang sudah ditutup tidak boleh menerima posting biasa.
- Semua proses menggunakan database transaction.

## 6.2 Cash Receipt Journal

Mencatat seluruh penerimaan kas.

Contoh sumber:
- Pembayaran pasien.
- Pelunasan piutang.
- Pendapatan tunai.
- Pengembalian uang muka.
- Penerimaan kas lainnya.

Informasi:
- Transaction Date.
- Receipt Number.
- Cash Account.
- Source Transaction.
- Customer / Patient.
- Description.
- Amount.
- Payment Method.
- Reference.
- Journal Reference.

Jurnal otomatis harus masuk ke General Ledger.

## 6.3 Cash Disbursement Journal

Mencatat seluruh pengeluaran kas.

Contoh:
- Pembayaran supplier.
- Pembayaran gaji.
- Pembayaran jasa dokter.
- Pembelian perlengkapan.
- Pengeluaran operasional.
- Pengeluaran lainnya.

Informasi:
- Payment Number.
- Payment Date.
- Cash Account.
- Expense Category.
- Supplier / Employee.
- Description.
- Amount.
- Reference.
- Journal Reference.

## 6.4 Bank Journal

Mencatat transaksi rekening bank.

Dukung:
- Transfer masuk.
- Transfer keluar.
- Pembayaran supplier.
- Penerimaan pasien.
- Biaya administrasi bank.
- Transfer antar rekening.
- Mutasi rekening.
- Bank Reconciliation.

Setiap rekening bank harus memiliki mapping akun COA tersendiri.

Transfer antar rekening tidak boleh dianggap sebagai pendapatan atau beban.

## 6.5 Payment Gateway Journal

Mencatat:
- QRIS.
- Virtual Account.
- E-Wallet.
- Debit Card.
- Credit Card.
- Payment Gateway lainnya.

Pisahkan:
- Gross Payment.
- Gateway Fee.
- Net Settlement.
- Settlement Date.
- Clearing Balance.
- Actual Bank Receipt.

## 6.6 Purchase Journal

Mencatat pembelian barang dan jasa.

Dukung:
- Pembelian tunai.
- Pembelian kredit.
- Purchase Invoice.
- Purchase Return.
- Purchase Discount.
- Tax.

## 6.7 Sales Journal

Mencatat seluruh transaksi penjualan.

Dukung:
- Penjualan obat.
- Penjualan alat medis.
- Pendapatan jasa.
- Penjualan tunai.
- Penjualan kredit.
- Diskon.
- Refund.
- Tax.

## 6.8 Adjusting Journal

Mendukung:
- Penyesuaian persediaan.
- Akrual beban.
- Akrual pendapatan.
- Penyusutan.
- Koreksi transaksi.
- Penyesuaian saldo awal.
- Penyesuaian lainnya.

## 6.9 Closing Journal

Digunakan saat tutup buku periode.

Dukung:
- Penutupan akun pendapatan.
- Penutupan akun beban.
- Perhitungan laba rugi.
- Transfer laba atau rugi ke akun ekuitas sesuai kebijakan.
- Validasi periode.
- Lock Period.

## 6.10 Reversal Journal

Koreksi jurnal harus melalui reversal.

Simpan:
- Original Journal ID.
- Reversal Journal ID.
- Reversal Reason.
- Created By.
- Created At.

Jangan menghapus histori transaksi finansial.

## 6.11 Ketentuan penting jurnal

Jurnal Kas Masuk, Jurnal Kas Keluar, Jurnal Bank, dan jurnal khusus lainnya tidak boleh menciptakan posting debit-kredit kedua.

Gunakan satu sumber posting akuntansi resmi.

Implementasikan journal projection atau view untuk menampilkan jurnal berdasarkan kategori tanpa menduplikasi General Ledger.

---

# 7. GENERAL LEDGER

Buat buku besar otomatis dari seluruh jurnal yang berstatus posted.

Informasi:
- Account Code.
- Account Name.
- Date.
- Journal Number.
- Reference.
- Description.
- Debit.
- Credit.
- Running Balance.
- Opening Balance.
- Closing Balance.

Dukung:
- Filter akun.
- Filter periode.
- Filter cabang.
- Filter sumber transaksi.
- Drill-down.
- Export PDF.
- Export Excel.
- Export CSV.

Saldo berjalan harus dihitung berdasarkan urutan posting dan tanggal yang konsisten.

---

# 8. TRIAL BALANCE

Buat Neraca Saldo otomatis.

Informasi:
- Account Code.
- Account Name.
- Opening Debit.
- Opening Credit.
- Period Debit.
- Period Credit.
- Ending Debit.
- Ending Credit.

Validasi:
Total debit harus sama dengan total kredit.

Sediakan:
- Unadjusted Trial Balance.
- Adjusted Trial Balance.
- Post-Closing Trial Balance.

---

# 9. INVENTORY MANAGEMENT

Buat sistem persediaan khusus klinik.

Pisahkan:
1. Obat.
2. Alat Kesehatan Habis Pakai.
3. BMHP.
4. Alat Medis Reusable.
5. Barang Operasional.

Fitur:
- Product Master.
- Category.
- Supplier.
- Warehouse.
- Purchase Order.
- Goods Receipt.
- Purchase Invoice.
- Stock In.
- Stock Out.
- Stock Adjustment.
- Stock Opname.
- Stock Transfer.
- Batch Number.
- Expired Date.
- Minimum Stock.
- Stock Valuation.
- Inventory Movement.
- Inventory History.

Metode:
- FIFO.
- Weighted Average.

Metode penilaian persediaan harus dapat dikonfigurasi sesuai kebijakan akuntansi.

Gunakan FEFO untuk prioritas pengeluaran obat berdasarkan tanggal kedaluwarsa.

Setiap pergerakan stok harus menghasilkan Inventory Movement Record.

Jangan mengizinkan perubahan stok langsung tanpa histori.

Integrasikan inventory dengan Accounting Engine.

Contoh pembelian obat kredit:

Debit: Persediaan Obat.
Kredit: Hutang Supplier.

Contoh penjualan obat:

Debit: Kas / Bank / Clearing / Piutang.
Kredit: Pendapatan Penjualan Obat.

Jurnal HPP:

Debit: HPP Obat.
Kredit: Persediaan Obat.

Contoh pemakaian BMHP untuk tindakan medis:

Debit: Beban Bahan Medis.
Kredit: Persediaan BMHP.

Pastikan tidak terjadi pengurangan stok dua kali.

---

# 10. PURCHASE MANAGEMENT

Buat modul pembelian:

- Purchase Request.
- Purchase Order.
- Goods Receipt.
- Purchase Invoice.
- Supplier Invoice.
- Purchase Return.
- Purchase Payment.
- Purchase History.

Dukung pembelian tunai dan kredit.

Alur:

Purchase Request
→ Purchase Order
→ Goods Receipt
→ Purchase Invoice
→ Automatic Journal
→ Payment
→ Accounts Payable Update.

Dukung partial receipt dan partial payment.

Validasi kuantitas, harga, diskon, pajak, dan referensi supplier.

---

# 11. SALES MANAGEMENT

Buat modul penjualan klinik.

Sumber pendapatan:
- Konsultasi.
- Pemeriksaan.
- Tindakan medis.
- Laboratorium.
- Penjualan obat.
- Penjualan alat kesehatan.
- Layanan kecantikan.
- Layanan lainnya.

Dukung:
- Sales Invoice.
- Sales Item.
- Service Item.
- Product Item.
- Discount.
- Tax.
- Refund.
- Partial Payment.
- Credit Sales.
- Payment Allocation.

Setiap transaksi penjualan harus otomatis terintegrasi dengan:
- Pendapatan.
- Piutang.
- Kas.
- Bank.
- Payment Gateway.
- Persediaan.
- HPP.
- General Ledger.

---

# 12. PAYMENT MANAGEMENT

Buat payment management terpusat.

Metode pembayaran:
- Cash.
- Bank Transfer.
- QRIS.
- Virtual Account.
- E-Wallet.
- Debit Card.
- Credit Card.
- Corporate Payment.
- Insurance Payment.

Fitur:
- Payment Creation.
- Payment Verification.
- Payment Allocation.
- Partial Payment.
- Multiple Payment.
- Refund.
- Payment Cancellation.
- Payment History.
- Payment Reconciliation.

Setiap pembayaran wajib memiliki:
- Payment Number.
- Transaction Reference.
- Amount.
- Payment Method.
- Payment Date.
- Payment Status.
- Source.
- Journal Reference.

Jangan menganggap invoice lunas hanya karena pembayaran dibuat.

Status lunas harus berdasarkan jumlah pembayaran valid yang telah dialokasikan ke invoice.

---

# 13. PAYMENT GATEWAY INTEGRATION

Buat adapter agar dapat mengintegrasikan berbagai provider.

Jangan mengikat Accounting Engine langsung ke provider tertentu.

Buat interface:

PaymentProvider:
- CreatePayment.
- GetPaymentStatus.
- VerifyWebhook.
- RefundPayment.
- GetSettlement.
- ReconcilePayment.

Dukung:
- Webhook.
- Signature Verification.
- Idempotency.
- Duplicate Event Prevention.
- Retry.
- Event Logging.
- Settlement Matching.

Jangan mempercayai redirect frontend sebagai bukti pembayaran berhasil.

Gunakan webhook tervalidasi atau verifikasi API provider.

## Accounting Flow

Ketika pembayaran berhasil:

Debit: Clearing Payment Gateway.
Kredit: Piutang atau akun pendapatan sesuai kondisi pengakuan transaksi.

Ketika settlement diterima:

Debit: Bank.
Debit: Beban Administrasi Payment Gateway.
Kredit: Clearing Payment Gateway.

Jika ada refund, buat jurnal koreksi sesuai sumber transaksi dan perlakuan akuntansinya.

Sediakan rekonsiliasi transaksi pembayaran dengan settlement provider dan mutasi rekening bank.

---

# 14. CASH MANAGEMENT

Buat modul pengelolaan kas.

Fitur:
- Cash Account.
- Cash Register.
- Opening Balance.
- Cash In.
- Cash Out.
- Cash Transfer.
- Cash Opname.
- Cash Adjustment.
- Cash Reconciliation.
- Cash History.

Dukung banyak kas:
- Kasir.
- Kas Kecil.
- Kas Operasional.
- Kas Cabang.

Sediakan laporan:
- Cash Book.
- Cash Receipt.
- Cash Disbursement.
- Daily Cash Report.
- Cash Position.
- Cash Reconciliation.

Setiap transaksi kas harus terhubung dengan Jurnal Umum dan General Ledger.

Transfer antar kas hanya memindahkan saldo dan bukan pendapatan atau beban.

---

# 15. BANK MANAGEMENT

Buat modul rekening bank.

Fitur:
- Bank Account.
- Bank Transaction.
- Bank Transfer.
- Bank Statement Import.
- Bank Reconciliation.
- Bank Fee.
- Bank Balance.

Dukung beberapa rekening bank dan cabang.

Sediakan proses rekonsiliasi:
- Matching berdasarkan nominal.
- Matching berdasarkan tanggal.
- Matching berdasarkan reference.
- Manual confirmation jika diperlukan.

Hasil rekonsiliasi harus dapat ditelusuri dan diaudit.

---

# 16. ACCOUNT RECEIVABLE (PIUTANG)

Buat sistem piutang untuk:

- Pasien.
- Perusahaan.
- Asuransi.
- BPJS jika relevan.
- Penjamin lainnya.

Fitur:
- Invoice.
- Due Date.
- Outstanding Balance.
- Partial Payment.
- Payment Allocation.
- Aging Report.
- Overdue.
- Collection History.
- Write-off dengan persetujuan.

Sistem harus menggunakan basis akrual untuk pengakuan pendapatan.

Jangan mengakui pendapatan hanya karena uang masuk apabila transaksi tersebut sebenarnya merupakan uang muka pasien.

Sediakan laporan:
- Daftar Piutang.
- Umur Piutang.
- Piutang Jatuh Tempo.
- Riwayat Pembayaran.
- Piutang per Customer.
- Piutang per Cabang.

---

# 17. ACCOUNT PAYABLE (HUTANG)

Buat sistem hutang:

- Hutang Supplier.
- Hutang Dokter.
- Hutang Gaji.
- Hutang Operasional.
- Hutang Pajak.
- Hutang Lainnya.

Fitur:
- Supplier Invoice.
- Due Date.
- Outstanding Balance.
- Partial Payment.
- Payment Allocation.
- Aging Report.
- Overdue.
- Payment History.

Pembayaran hutang harus otomatis mengurangi saldo hutang.

Sediakan:
- Daftar Hutang.
- Umur Hutang.
- Hutang Jatuh Tempo.
- Hutang per Supplier.
- Riwayat Pembayaran.

---

# 18. FIXED ASSET MANAGEMENT

Buat modul aset tetap.

Fitur:
- Asset Register.
- Asset Category.
- Acquisition Cost.
- Acquisition Date.
- Useful Life.
- Residual Value.
- Depreciation Method.
- Depreciation Schedule.
- Accumulated Depreciation.
- Asset Disposal.
- Asset Maintenance.

Dukung metode penyusutan yang sesuai kebijakan akuntansi.

Jurnal penyusutan otomatis:

Debit: Beban Penyusutan.
Kredit: Akumulasi Penyusutan.

Jangan langsung mengklasifikasikan seluruh pembelian alat medis sebagai aset tetap.

Gunakan konfigurasi kebijakan kapitalisasi aset.

---

# 19. FINANCIAL REPORTING

Semua laporan keuangan harus dihasilkan otomatis dari General Ledger.

Tidak boleh ada tabel laporan yang membutuhkan input angka manual.

## 19.1 Income Statement / Laporan Laba Rugi

Tampilkan:

Pendapatan:
- Pendapatan Konsultasi.
- Pendapatan Pemeriksaan.
- Pendapatan Tindakan.
- Pendapatan Laboratorium.
- Pendapatan Penjualan Obat.
- Pendapatan Lainnya.

HPP:
- HPP Obat.
- HPP Alat Kesehatan Habis Pakai.
- HPP BMHP.

Perhitungan:

Laba Kotor = Pendapatan - HPP.

Beban:
- Gaji.
- Jasa Dokter.
- Sewa.
- Listrik.
- Air.
- Internet.
- Administrasi Bank.
- Payment Gateway Fee.
- Penyusutan.
- Operasional Lainnya.

Laba Operasional = Laba Kotor - Beban Operasional.

Laba Bersih = Laba Sebelum Pajak - Beban Pajak.

Dukung laporan bulanan, triwulanan, semesteran, dan tahunan.

## 19.2 Balance Sheet / Laporan Neraca

Tampilkan:

Aset:
- Aset Lancar.
- Kas.
- Bank.
- Piutang.
- Persediaan.
- Aset Tidak Lancar.
- Aset Tetap.
- Akumulasi Penyusutan.

Liabilitas:
- Hutang Supplier.
- Hutang Dokter.
- Hutang Gaji.
- Hutang Pajak.
- Hutang Lainnya.

Ekuitas:
- Modal.
- Laba Ditahan.
- Laba Tahun Berjalan.

Persamaan wajib:

ASSET = LIABILITY + EQUITY.

Laporan harus mendukung tanggal tertentu (as-of date).

## 19.3 Cash Flow Statement / Laporan Arus Kas

Dukung:

- Metode Langsung.
- Metode Tidak Langsung.

Kelompok:

Operating Activities:
- Penerimaan pasien.
- Penerimaan piutang.
- Pembayaran supplier.
- Pembayaran gaji.
- Pembayaran operasional.

Investing Activities:
- Pembelian aset tetap.
- Penjualan aset tetap.

Financing Activities:
- Setoran modal.
- Pinjaman.
- Pembayaran pokok pinjaman.
- Distribusi pemilik.

Jangan menganggap transaksi akrual sebagai arus kas.

Gunakan cash flow classification mapping untuk mengelompokkan setiap akun atau transaksi.

Transaksi nonkas tidak boleh dicatat sebagai arus kas.

## 19.4 General Ledger / Buku Besar

Tampilkan:
- Saldo Awal.
- Debit.
- Kredit.
- Saldo Berjalan.
- Saldo Akhir.

Dukung drill-down ke jurnal sumber.

## 19.5 Trial Balance / Neraca Saldo

Tampilkan:
- Saldo Awal Debit.
- Saldo Awal Kredit.
- Mutasi Debit.
- Mutasi Kredit.
- Saldo Akhir Debit.
- Saldo Akhir Kredit.

## 19.6 Financial Reports Tambahan

Buat:
- Jurnal Umum.
- Jurnal Kas Masuk.
- Jurnal Kas Keluar.
- Jurnal Bank.
- Jurnal Pembelian.
- Jurnal Penjualan.
- Jurnal Penyesuaian.
- Jurnal Penutup.
- Buku Kas.
- Buku Bank.
- Neraca Saldo.
- Perubahan Ekuitas.
- Laporan Persediaan.
- Laporan HPP.
- Laporan Piutang.
- Laporan Hutang.
- Aging Piutang.
- Aging Hutang.
- Rekonsiliasi Bank.
- Rekonsiliasi Payment Gateway.
- Laporan Pendapatan.
- Laporan Beban.
- Laporan Penyusutan.
- Laporan Arus Kas Harian.

---

# 20. AUTOMATIC TRANSACTION-TO-JOURNAL MAPPING

Buat sistem mapping transaksi ke akun COA.

Contoh:

A. Penjualan jasa tunai:

Debit: Kas.
Kredit: Pendapatan Jasa.

B. Penjualan jasa melalui gateway:

Debit: Clearing Payment Gateway.
Kredit: Pendapatan Jasa atau Piutang yang telah diakui.

C. Penjualan kredit:

Debit: Piutang.
Kredit: Pendapatan Jasa.

D. Pembelian obat kredit:

Debit: Persediaan Obat.
Kredit: Hutang Supplier.

E. Pembayaran hutang:

Debit: Hutang Supplier.
Kredit: Bank.

F. Penjualan obat:

Debit: Kas / Bank / Clearing / Piutang.
Kredit: Pendapatan Penjualan Obat.

Jurnal HPP:

Debit: HPP Obat.
Kredit: Persediaan Obat.

G. Settlement Payment Gateway:

Debit: Bank.
Debit: Beban Gateway.
Kredit: Clearing Gateway.

H. Penerimaan piutang:

Debit: Kas / Bank.
Kredit: Piutang.

I. Pembayaran gaji:

Debit: Beban Gaji.
Kredit: Kas / Bank.

J. Pembelian aset tetap:

Debit: Aset Tetap.
Kredit: Kas / Bank / Hutang.

K. Penyusutan:

Debit: Beban Penyusutan.
Kredit: Akumulasi Penyusutan.

L. Pembayaran uang muka pasien:

Debit: Kas / Bank.
Kredit: Uang Muka Pasien.

Saat jasa diberikan dan pendapatan telah memenuhi syarat pengakuan:

Debit: Uang Muka Pasien.
Kredit: Pendapatan Jasa.

M. Pembelian barang tunai:

Debit: Persediaan / Aset / Beban sesuai klasifikasi.
Kredit: Kas / Bank.

Setiap mapping harus configurable dan memiliki validasi.

Jangan mengasumsikan seluruh transaksi pembayaran merupakan pendapatan.

---

# 21. FINANCIAL PERIOD MANAGEMENT

Buat pengelolaan periode akuntansi.

Fitur:
- Fiscal Year.
- Accounting Period.
- Open Period.
- Close Period.
- Reopen Period dengan otorisasi.
- Period Lock.
- Closing Checklist.
- Closing Journal.
- Audit Trail.

Ketika periode sudah ditutup:
- Transaksi baru tidak boleh diposting ke periode tersebut.
- Jurnal tidak dapat diedit.
- Koreksi harus melalui periode terbuka dan mekanisme reversal sesuai kebijakan.
- Seluruh laporan periode tersebut harus tetap dapat diakses.

---

# 22. TRANSACTION LINK DAN REPORT DRILL-DOWN

Semua transaksi harus memiliki URL detail.

Contoh:

GET /api/v1/transactions/{id}

GET /api/v1/sales/{id}

GET /api/v1/purchases/{id}

GET /api/v1/payments/{id}

GET /api/v1/journals/{id}

GET /api/v1/inventory/movements/{id}

GET /api/v1/reports/income-statement

GET /api/v1/reports/balance-sheet

GET /api/v1/reports/cash-flow

GET /api/v1/reports/general-ledger

GET /api/v1/reports/trial-balance

GET /api/v1/reports/cash-journal

GET /api/v1/reports/bank-journal

Setiap laporan wajib mendukung:

- Filter tanggal.
- Filter periode.
- Filter cabang.
- Filter akun.
- Filter kategori.
- Filter transaksi.
- Pagination.
- Sorting.
- Summary.
- Total.
- Subtotal.
- Drill-down.
- Export.

Ketika akuntan membuka angka pendapatan pada laporan laba rugi, sistem harus menampilkan daftar transaksi pembentuk angka tersebut.

Dari transaksi, akuntan dapat membuka jurnal.

Dari jurnal, akuntan dapat membuka detail pembayaran atau transaksi operasional.

Gunakan authenticated URL dan RBAC. Jangan membuat tautan publik yang dapat membuka data keuangan tanpa otorisasi.

---

# 23. DATABASE DESIGN

Rancang PostgreSQL database.

Tabel minimal:

AUTHENTICATION:
- users
- roles
- permissions
- user_roles
- role_permissions
- refresh_tokens
- audit_logs

MASTER DATA:
- branches
- patients
- customers
- suppliers
- products
- product_categories
- product_units
- warehouses
- departments
- cost_centers

ACCOUNTING:
- accounts
- accounting_mappings
- journal_entries
- journal_lines
- fiscal_years
- fiscal_periods
- journal_sequences
- cash_flow_mappings

INVENTORY:
- inventory_batches
- inventory_movements
- inventory_balances
- stock_adjustments
- stock_opnames
- stock_transfers

PURCHASE:
- purchase_requests
- purchase_orders
- purchase_order_items
- goods_receipts
- purchase_invoices
- purchase_invoice_items
- purchase_returns

SALES:
- sales
- sales_items
- sales_invoices
- sales_invoice_items
- sales_returns

PAYMENT:
- payments
- payment_allocations
- payment_transactions
- payment_webhook_events
- payment_settlements
- payment_refunds

CASH AND BANK:
- cash_accounts
- cash_transactions
- bank_accounts
- bank_transactions
- bank_statements
- bank_reconciliations

RECEIVABLE:
- receivables
- receivable_payments
- receivable_aging_snapshots jika diperlukan

PAYABLE:
- payables
- payable_payments
- payable_aging_snapshots jika diperlukan

FIXED ASSETS:
- fixed_assets
- asset_categories
- depreciation_records
- asset_disposals

REPORTING:
- report_jobs jika asynchronous generation diperlukan.
- report_export_files.

Buat relasi foreign key yang tepat.

Tambahkan:
- Indexes.
- Unique Constraints.
- Check Constraints.
- Created At.
- Updated At.
- Created By.
- Updated By.
- Branch ID.
- Soft Delete jika relevan.

Jangan menyimpan saldo keuangan sebagai data yang dapat diubah secara bebas.

Saldo harus dapat ditelusuri ke jurnal dan transaksi sumber.

Jika menggunakan cached balance atau materialized view, harus tersedia proses rekonsiliasi dengan sumber data utama.

---

# 24. API RESPONSE STANDARD

Gunakan response JSON yang konsisten.

Success:

{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {},
  "meta": {}
}

Error:

{
  "success": false,
  "message": "Validation failed",
  "errors": []
}

Gunakan HTTP status code yang sesuai.

Sediakan Swagger untuk semua endpoint.

Gunakan versioning:

/api/v1/

---

# 25. SECURITY DAN AUDIT TRAIL

Implementasikan:

- JWT.
- Refresh Token Rotation.
- RBAC.
- Password Hashing.
- Rate Limiting.
- Input Validation.
- SQL Injection Protection.
- Webhook Signature Validation.
- Idempotency.
- Audit Logging.
- Financial Approval.
- Branch Access Restriction.
- Period Locking.

Setiap aktivitas finansial harus dapat ditelusuri:

- Siapa yang membuat.
- Siapa yang mengubah.
- Siapa yang menyetujui.
- Kapan transaksi terjadi.
- Kapan jurnal diposting.
- Apa sumber transaksinya.
- Apa alasan koreksinya.

Jangan menghapus histori jurnal yang sudah diposting.

---

# 26. APPROVAL WORKFLOW

Buat approval workflow configurable.

Contoh:

- Purchase Order.
- Purchase Invoice.
- Payment Request.
- Cash Disbursement.
- Journal Manual.
- Stock Adjustment.
- Refund.
- Period Reopening.

Dukung:
- Pending.
- Approved.
- Rejected.
- Cancelled.

Batas persetujuan dapat dikonfigurasi berdasarkan nominal, cabang, dan role.

---

# 27. TESTING DAN ACCOUNTING VALIDATION

Buat test untuk semua proses.

ACCOUNTING:
- Debit harus sama dengan kredit.
- Tidak ada jurnal ganda.
- Jurnal posted tidak dapat diubah langsung.
- Reversal harus membalik jurnal dengan benar.
- Periode tertutup menolak posting biasa.

INVENTORY:
- Stok bertambah setelah penerimaan barang.
- Stok berkurang setelah penjualan.
- HPP dihitung sesuai metode penilaian.
- FEFO berfungsi.
- Stock opname menghasilkan penyesuaian yang dapat diaudit.

PAYMENT:
- Pembayaran berhasil menghasilkan jurnal yang benar.
- Webhook duplikat tidak menghasilkan jurnal duplikat.
- Refund menghasilkan jurnal yang sesuai.
- Settlement sesuai dengan nominal bruto, fee, dan net.
- Partial payment menghitung outstanding balance dengan benar.

REPORTING:
- Neraca seimbang.
- Laba rugi sesuai General Ledger.
- Arus kas hanya mencatat transaksi kas.
- Trial Balance seimbang.
- Buku besar sesuai jurnal.
- Piutang dan hutang sesuai saldo akun kontrol.
- Nilai persediaan sesuai inventory valuation.

Buat integration testing menggunakan database PostgreSQL sungguhan melalui Testcontainers.

---

# 28. PROJECT FOLDER STRUCTURE

Gunakan struktur seperti berikut:

cmd/
  api/
    main.go
  worker/
    main.go

internal/
  domain/
    accounting/
    inventory/
    sales/
    purchase/
    payment/
    cash/
    bank/
    receivable/
    payable/
    asset/

  application/
    accounting/
    inventory/
    sales/
    purchase/
    payment/
    reporting/

  infrastructure/
    database/
    redis/
    payment_gateway/
    storage/
    logger/

  delivery/
    http/
      handlers/
      middleware/
      routes/
      dto/

pkg/
  decimal/
  response/
  validator/
  pagination/

migrations/
docs/
tests/
scripts/
deployments/

Buat dependency injection yang jelas.

Hindari circular dependency.

---

# 29. DEVELOPMENT PHASE

Implementasikan sistem bertahap.

PHASE 1 — FOUNDATION

- Project Setup.
- PostgreSQL.
- Migration.
- Authentication.
- RBAC.
- Master Data.
- Branch Management.
- Chart of Accounts.
- Accounting Mapping.
- Journal Entry.
- Journal Lines.
- General Ledger.
- Trial Balance.

PHASE 2 — JOURNAL MANAGEMENT

- General Journal.
- Cash Receipt Journal.
- Cash Disbursement Journal.
- Bank Journal.
- Sales Journal.
- Purchase Journal.
- Payment Gateway Journal.
- Adjusting Journal.
- Reversal Journal.
- Journal Approval.
- Journal Posting.

PHASE 3 — INVENTORY AND PURCHASING

- Product.
- Warehouse.
- Supplier.
- Purchase Order.
- Goods Receipt.
- Purchase Invoice.
- Stock Movement.
- Batch.
- Expired Date.
- Stock Opname.
- Inventory Valuation.
- Automatic HPP.

PHASE 4 — SALES AND PAYMENT

- Sales.
- Sales Invoice.
- Cash Payment.
- Bank Transfer.
- Payment Gateway.
- Webhook.
- Settlement.
- Refund.
- Cash Management.
- Bank Management.
- Reconciliation.

PHASE 5 — RECEIVABLE AND PAYABLE

- Accounts Receivable.
- Accounts Payable.
- Partial Payment.
- Payment Allocation.
- Aging Report.
- Collection.
- Supplier Payment.

PHASE 6 — FINANCIAL REPORTING

- General Ledger.
- Trial Balance.
- Income Statement.
- Balance Sheet.
- Cash Flow.
- Cash Journal.
- Bank Journal.
- Equity Report.
- Inventory Report.
- Receivable Report.
- Payable Report.

PHASE 7 — FIXED ASSET AND CLOSING

- Asset Register.
- Depreciation.
- Journal Closing.
- Fiscal Period.
- Period Lock.
- Audit Trail.
- Approval Workflow.

PHASE 8 — EXPORT AND PRODUCTION

- PDF.
- Excel.
- CSV.
- Swagger.
- Docker.
- Security Testing.
- Performance Testing.
- Integration Testing.
- Production Deployment.

---

# 30. REQUIRED OUTPUT FROM AI CODING AGENT

Pada tahap awal, berikan:

1. Analisis kebutuhan.
2. System Architecture Diagram.
3. ERD lengkap.
4. Database Schema.
5. Chart of Accounts klinik.
6. Transaction Flow.
7. Accounting Flow.
8. Journal Mapping.
9. API Endpoint Specification.
10. Project Folder Structure.

Setelah desain disetujui, lanjutkan implementasi source code Golang.

Setiap implementasi wajib mencakup:
- File path.
- Source code lengkap.
- Migration.
- Handler.
- Service.
- Repository.
- Domain Model.
- Request DTO.
- Response DTO.
- Unit Test.
- Integration Test.
- Contoh API Request.
- Contoh API Response.

Jangan hanya memberikan pseudo-code.

Jangan memberikan implementasi yang tidak dapat dikompilasi.

Jangan mengarang aturan akuntansi yang belum dikonfirmasi.

Jika terdapat keputusan akuntansi yang belum diketahui, berikan pilihan konfigurasi dan jelaskan dampaknya.

Jangan membuat seluruh modul dalam satu file.

Jangan membuat jurnal ganda.

Jangan membuat laporan keuangan berdasarkan data hardcode.

Jangan mengizinkan perubahan langsung pada saldo akun tanpa transaksi yang dapat diaudit.

---

# 31. ACCEPTANCE CRITERIA

Sistem dianggap memenuhi kebutuhan apabila:

1. Semua transaksi operasional terhubung dengan Accounting Engine.
2. Semua transaksi yang memerlukan jurnal menghasilkan jurnal otomatis.
3. Total debit sama dengan total kredit.
4. Jurnal Umum menjadi sumber posting utama.
5. Jurnal Kas Masuk dan Keluar terintegrasi tanpa duplikasi.
6. Jurnal Bank dan Payment Gateway dapat direkonsiliasi.
7. Inventory terintegrasi dengan HPP.
8. Piutang dan hutang terintegrasi dengan pembayaran.
9. General Ledger dihasilkan otomatis.
10. Trial Balance seimbang.
11. Laporan Laba Rugi dihasilkan otomatis.
12. Laporan Neraca dihasilkan otomatis.
13. Laporan Arus Kas dihasilkan otomatis.
14. Setiap angka laporan dapat ditelusuri sampai transaksi sumber.
15. Transaksi dan laporan memiliki endpoint yang terdokumentasi.
16. Laporan dapat diekspor PDF, Excel, dan CSV.
17. Semua perubahan keuangan memiliki audit trail.
18. Periode akuntansi dapat ditutup dan dikunci.
19. Sistem mendukung multi-cabang.
20. Seluruh test accounting dan integration berhasil.

---

# 32. FINAL INSTRUCTION

Mulai pekerjaan dengan melakukan analisis kebutuhan secara menyeluruh.

Jangan langsung melakukan coding.

Pertama tampilkan:

1. Arsitektur sistem.
2. ERD database.
3. Chart of Accounts.
4. Daftar seluruh jenis jurnal.
5. Transaction-to-Journal Mapping.
6. Contoh jurnal debit-kredit.
7. Alur integrasi Inventory, Payment, dan Accounting.
8. Desain API.
9. Roadmap implementasi.

Pastikan desain dapat mengakomodasi kebutuhan klinik yang saat ini belum memiliki struktur piutang, hutang, dan pendapatan yang final.

Gunakan konfigurasi agar COA, mapping akun, kategori transaksi, dan perlakuan akuntansi dapat dikembangkan tanpa mengubah keseluruhan source code.

Prioritaskan ketepatan akuntansi, integritas data, kemudahan audit, dan otomatisasi laporan.

Setelah seluruh desain selesai, tunggu persetujuan sebelum memulai implementasi Phase 1.