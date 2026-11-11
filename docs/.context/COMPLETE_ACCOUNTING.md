# ADDITIONAL MASTER PROMPT — COMPLETE CLINIC ACCOUNTING ENGINE

## 1. ROLE

Bertindak sebagai:
- Senior Accounting System Architect.
- Senior Golang Backend Engineer.
- ERP Financial Management Specialist.
- Financial Controller.
- Senior Accountant.
- Database Architect.
- Tax System Analyst.
- Internal Control and Audit Specialist.

Saya sedang mengembangkan sistem ERP akuntansi klinik berbasis Golang dan PostgreSQL.

Gunakan seluruh master prompt dan blueprint yang sudah diberikan sebelumnya sebagai spesifikasi utama.

Tugas kamu adalah melakukan audit terhadap desain dan implementasi Accounting Module yang sudah ada, menemukan bagian yang belum lengkap, lalu menyempurnakannya.

**Jangan langsung membuat ulang modul yang sudah tersedia. Jangan menganggap modul sudah selesai hanya karena folder, entity, atau endpoint-nya sudah dibuat. Periksa apakah business logic, database, jurnal, validasi, laporan, dan pengujiannya benar-benar berfungsi.**

Sistem harus mendukung dua unit bisnis:
- Kidinest: klinik spesialis anak, tumbuh kembang, vaksinasi, farmasi pediatrik.
- Semarang Skin Center (SSC): klinik estetika, dermatologi, facial, laser, dan produk skincare.

Sistem harus dapat mengelola transaksi kedua unit secara terpisah dan secara konsolidasi, termasuk shared back-office.

---

# 2. ACCOUNTING ENGINE FOUNDATION

Bangun Accounting Engine yang mendukung:

1. Double-entry bookkeeping.
2. Accrual basis accounting.
3. Automatic journal generation.
4. Configurable account mapping.
5. Journal validation.
6. Journal approval.
7. Journal posting.
8. Journal reversal.
9. Adjusting journal.
10. Recurring journal.
11. General Ledger.
12. Subsidiary Ledger.
13. Trial Balance.
14. Financial Statement.
15. Period Closing.
16. Year-End Closing.
17. Opening Balance.
18. Audit Trail.
19. Accounting Reconciliation.
20. Multi-branch and multi-business-unit accounting.

Semua proses akuntansi harus memiliki sumber data yang dapat ditelusuri.

Jangan membuat saldo akun yang dapat dimodifikasi secara bebas.

Saldo akun harus berasal dari transaksi jurnal yang sudah posted.

## 2.1 Accounting Policy Configuration

Buat konfigurasi kebijakan akuntansi yang dapat diatur tanpa mengubah business logic.

Konfigurasi minimal:
- Accounting basis.
- Fiscal year start.
- Fiscal period.
- Inventory valuation method.
- Depreciation method.
- Revenue recognition policy.
- Expense recognition policy.
- Tax configuration.
- Rounding policy.
- Journal numbering.
- Materiality threshold jika diperlukan.
- Retained earnings account.
- Default suspense account.
- Default cash and bank accounts.

Simpan perubahan kebijakan beserta:
- Nilai lama.
- Nilai baru.
- Effective date.
- Changed by.
- Approval status.
- Reason.

Perubahan konfigurasi tidak boleh mengubah jurnal historis yang sudah posted.

---

# 3. COMPLETE CHART OF ACCOUNTS

Kembangkan COA menjadi sistem hierarki yang benar-benar dapat digunakan untuk pembukuan.

Dukung:
- Account Code.
- Account Name.
- Parent Account.
- Account Level.
- Account Type.
- Account Subtype.
- Normal Balance.
- Is Header.
- Is Posting Account.
- Is Active.
- Is System Account.
- Allow Manual Journal.
- Branch Scope.
- Unit Scope.
- Effective Date.
- Deactivation Date.

Jenis akun:
- Asset.
- Liability.
- Equity.
- Revenue.
- Cost of Goods Sold.
- Expense.
- Other Income.
- Other Expense.

Aturan:
1. Akun header tidak boleh digunakan untuk posting.
2. Akun posting harus memiliki tipe dan saldo normal yang valid.
3. Akun yang sudah digunakan tidak boleh dihapus.
4. Akun dapat dinonaktifkan tanpa menghilangkan histori.
5. Perubahan struktur COA harus tercatat di audit trail.
6. Akun sistem tidak boleh sembarang dihapus atau diubah.
7. COA harus mendukung pemisahan akun berdasarkan unit bisnis jika memang dibutuhkan.
8. Jangan membuat akun baru secara otomatis setiap kali transaksi dibuat.

Buat COA template untuk:
- Kidinest.
- Semarang Skin Center.
- Shared Back-office.
- Consolidated Reporting.

Buat account mapping untuk:
- Cash.
- Bank.
- Payment Gateway Clearing.
- Accounts Receivable.
- Accounts Payable.
- Inventory.
- Revenue.
- COGS.
- Doctor Fee.
- Therapist Commission.
- Tax.
- Fixed Assets.
- Accumulated Depreciation.
- Depreciation Expense.
- Patient Deposit.
- Unearned Revenue.
- Discount.
- Refund.
- Cash Shortage.
- Cash Overage.

Semua mapping harus configurable dan memiliki validasi akun aktif serta tipe akun yang sesuai.

---

# 4. ACCOUNTING DIMENSIONS

Buat sistem dimensi akuntansi agar laporan bisa dianalisis berdasarkan unit bisnis dan pusat biaya.

Dimensi minimal:

1. Legal Entity.
2. Business Unit.
3. Branch.
4. Department.
5. Cost Center.
6. Service Category.
7. Doctor.
8. Employee.
9. Warehouse.
10. Project atau Program jika dibutuhkan.

Contoh:
- UNIT-KIDINEST.
- UNIT-SSC.
- SHARED-BACKOFFICE.
- COSTCENTER-ADMIN.
- COSTCENTER-MARKETING.
- COSTCENTER-FARMASI.
- COSTCENTER-MEDICAL.

Setiap journal line harus dapat menyimpan dimensi yang relevan.

Jangan mewajibkan semua dimensi pada semua transaksi. Gunakan konfigurasi berdasarkan jenis akun dan jenis transaksi.

Sediakan:
- Profit and Loss by Business Unit.
- Profit and Loss by Branch.
- Expense by Cost Center.
- Revenue by Service Category.
- Doctor Revenue Analysis.
- Consolidated Financial Statements.

## 4.1 Shared Cost Allocation

Buat fitur alokasi biaya bersama antara Kidinest dan SSC.

Contoh biaya:
- Sewa gedung.
- Listrik.
- Internet.
- Administrasi.
- Marketing.
- Gaji back-office.
- Biaya operasional bersama.

Dukung metode alokasi:
- Fixed percentage.
- Revenue proportion.
- Transaction volume.
- Employee headcount.
- Area usage.
- Manual allocation dengan approval.

Simpan:
- Source expense.
- Allocation method.
- Allocation basis.
- Allocation percentage.
- Allocation amount.
- Target business unit.
- Allocation period.
- Approval.
- Journal reference.

Pastikan total hasil alokasi sama dengan nilai sumber. Jangan menggandakan total beban konsolidasi.

Jika kedua unit berada dalam satu entitas hukum, jangan membuat jurnal intercompany fiktif. Jika berada dalam entitas hukum berbeda, desain mekanisme intercompany yang sesuai.

---

# 5. COMPLETE JOURNAL ENGINE

Lengkapi seluruh siklus jurnal:

Draft → Submitted → Approved → Posted

Dukung status:
- Draft.
- Pending Approval.
- Approved.
- Rejected.
- Posted.
- Reversed.
- Cancelled.

Setiap jurnal memiliki:
- Journal Number.
- Journal Date.
- Posting Date.
- Reference Number.
- Source Module.
- Source Transaction ID.
- Description.
- Currency.
- Exchange Rate jika relevan.
- Total Debit.
- Total Credit.
- Branch.
- Business Unit.
- Created By.
- Approved By.
- Posted By.
- Posted At.
- Reversal Reference.

Journal line memiliki:
- Account ID.
- Description.
- Debit.
- Credit.
- Branch.
- Business Unit.
- Cost Center.
- Department.
- Doctor ID jika relevan.
- Patient ID jika relevan.
- Supplier ID jika relevan.
- Source Reference.

## 5.1 Journal Validation

Validasi:
- Total debit harus sama dengan total kredit.
- Nilai debit dan kredit tidak boleh negatif.
- Satu baris tidak boleh memiliki debit dan kredit sekaligus, kecuali desain khusus yang tervalidasi.
- Jurnal harus memiliki minimal dua baris.
- Account harus aktif dan dapat diposting.
- Posting date harus berada dalam periode yang terbuka.
- Reference transaksi tidak boleh menghasilkan posting duplikat.
- Jurnal harus memiliki sumber transaksi yang dapat ditelusuri.
- Dimensi wajib harus terisi.
- Tidak boleh menggunakan akun yang tidak sesuai dengan konfigurasi transaksi.

Gunakan decimal, bukan float.

## 5.2 Journal Posting

Posting jurnal harus:
- Menggunakan database transaction.
- Mengunci data yang relevan untuk mencegah race condition.
- Menggunakan idempotency key.
- Memastikan periode masih terbuka.
- Memastikan jurnal belum pernah diposting.
- Menyimpan posting reference.
- Memperbarui atau menyediakan ledger projection secara konsisten.

Jangan mengubah jurnal posted secara langsung.

Koreksi harus menggunakan reversal dan jurnal pengganti.

## 5.3 Journal Types

Lengkapi:
- General Journal.
- Cash Receipt Journal.
- Cash Disbursement Journal.
- Bank Receipt Journal.
- Bank Disbursement Journal.
- Sales Journal.
- Purchase Journal.
- Payment Gateway Journal.
- Inventory Adjustment Journal.
- Cost of Goods Sold Journal.
- Doctor Fee Accrual Journal.
- Tax Journal.
- Payroll Journal.
- Depreciation Journal.
- Prepaid Expense Journal.
- Accrued Expense Journal.
- Accrued Revenue Journal.
- Adjusting Journal.
- Reclassification Journal.
- Closing Journal.
- Reversal Journal.

Semua jenis jurnal harus tetap masuk ke satu sumber General Ledger resmi. Jangan membuat pembukuan ganda.

---

# 6. SUBSIDIARY LEDGER AND RECONCILIATION

General Ledger harus dapat direkonsiliasi dengan buku pembantu.

Buat subsidiary ledger untuk:

1. Accounts Receivable.
2. Accounts Payable.
3. Inventory.
4. Fixed Assets.
5. Doctor Payable.
6. Tax Payable.
7. Cash.
8. Bank.
9. Payment Gateway Clearing.
10. Patient Deposit.

Sediakan proses rekonsiliasi:

- Total AR sub-ledger versus AR General Ledger.
- Total AP sub-ledger versus AP General Ledger.
- Inventory valuation versus Inventory General Ledger.
- Fixed asset register versus Fixed Asset General Ledger.
- Doctor payable versus Doctor Payable General Ledger.
- Tax payable versus Tax Payable General Ledger.
- Cash balance versus cash transaction records.
- Bank balance versus bank transaction records.
- Gateway clearing versus provider settlement.

Jika terdapat selisih:
- Tampilkan nominal selisih.
- Tampilkan akun terkait.
- Tampilkan transaksi pembentuk selisih.
- Tampilkan periode.
- Tampilkan kemungkinan penyebab.
- Sediakan proses investigasi.
- Sediakan adjustment melalui jurnal yang disetujui.

Jangan otomatis menghilangkan selisih dengan mengubah saldo.

---

# 7. OPENING BALANCE AND MIGRATION

Buat modul saldo awal ketika klinik mulai menggunakan sistem.

Dukung:
- Opening Cash.
- Opening Bank.
- Opening AR.
- Opening AP.
- Opening Inventory.
- Opening Fixed Assets.
- Opening Accumulated Depreciation.
- Opening Tax Balance.
- Opening Patient Deposit.
- Opening Doctor Payable.
- Opening Equity.
- Other Opening Balances.

Fitur:
- Import Excel/CSV.
- Validasi COA.
- Validasi dimensi.
- Validasi debit dan kredit.
- Preview sebelum posting.
- Approval.
- Posting opening journal.
- Opening sub-ledger detail.
- Import history.
- Error report.
- Reconciliation report.

Opening balance tidak boleh hanya mengisi saldo akun tanpa detail pendukung.

Untuk persediaan, simpan detail batch, tanggal kedaluwarsa, kuantitas, dan nilai perolehan.

Untuk AR/AP, simpan detail per invoice dan tanggal jatuh tempo.

Untuk aset tetap, simpan harga perolehan, tanggal akuisisi, umur manfaat, nilai residu, dan akumulasi penyusutan.

---

# 8. ACCRUAL, PREPAID, AND ADJUSTING ENTRIES

Buat modul pengelolaan penyesuaian akuntansi.

## 8.1 Accrued Expense

Contoh:
- Gaji yang belum dibayar.
- Jasa dokter yang sudah menjadi hak dokter tetapi belum dibayar.
- Tagihan listrik yang belum diterima.
- Biaya jasa vendor yang sudah terjadi.

Dukung:
- Accrual date.
- Amount.
- Account mapping.
- Supporting document.
- Reversal date jika diperlukan.
- Approval.
- Journal reference.

## 8.2 Prepaid Expense

Contoh:
- Sewa dibayar di muka.
- Asuransi dibayar di muka.
- Lisensi software.
- Biaya kontrak tahunan.

Dukung:
- Initial payment.
- Coverage start.
- Coverage end.
- Allocation period.
- Monthly amortization.
- Remaining balance.
- Automatic journal.
- Early termination adjustment.

## 8.3 Accrued Revenue

Dukung pendapatan yang sudah menjadi hak klinik tetapi belum ditagihkan atau dibayar, sesuai kebijakan pengakuan pendapatan.

## 8.4 Adjusting Entries

Dukung:
- Accrual.
- Prepaid amortization.
- Depreciation.
- Inventory write-off.
- Inventory shrinkage.
- Bad debt allowance.
- Revenue adjustment.
- Expense reclassification.
- Tax adjustment.
- Other approved adjustments.

Setiap penyesuaian harus memiliki alasan dan bukti pendukung.

---

# 9. REVENUE RECOGNITION AND PATIENT DEPOSIT

Buat aturan pengakuan pendapatan berdasarkan jenis layanan.

Jenis transaksi:
- Konsultasi.
- Pemeriksaan.
- Vaksinasi.
- Laboratorium.
- Tindakan medis.
- Facial.
- Laser.
- Treatment.
- Penjualan obat.
- Penjualan skincare.
- Paket treatment.
- Deposit pasien.
- Corporate billing.
- Insurance billing.

Bedakan:
1. Invoice issued.
2. Revenue earned.
3. Payment received.
4. Receivable recognized.
5. Advance payment.
6. Unearned revenue.
7. Revenue recognized over time jika relevan.

Jangan mengakui pendapatan hanya karena pembayaran diterima.

## 9.1 Patient Deposit

Dukung:
- Deposit creation.
- Deposit usage.
- Deposit refund.
- Deposit expiry sesuai kebijakan dan aturan yang berlaku.
- Deposit balance.
- Deposit allocation.
- Deposit reconciliation.

Contoh penerimaan uang muka:

Debit: Kas/Bank/Clearing.
Kredit: Uang Muka Pasien.

Ketika layanan telah diberikan dan memenuhi kriteria pengakuan:

Debit: Uang Muka Pasien atau Piutang.
Kredit: Pendapatan terkait.

Buat alur terpisah untuk:
- Pembayaran penuh.
- Pembayaran sebagian.
- Deposit.
- Refund.
- Pembatalan.
- Paket layanan.

## 9.2 Treatment Package and Voucher

Dukung:
- Package sales.
- Package usage.
- Remaining session.
- Expiry date.
- Package discount.
- Refund.
- Package revenue recognition.
- Liability for unused services.

Jangan langsung mengakui seluruh nilai paket sebagai pendapatan apabila jasa belum diberikan dan kebijakan akuntansi mensyaratkan pengakuan bertahap.

---

# 10. COMPLETE DOCTOR SHARING AND THERAPIST COMMISSION

Buat modul Doctor Sharing Fee dan Therapist Commission yang terhubung ke accounting.

Dukung skema:
- Percentage sharing.
- Fixed fee.
- Tiered percentage.
- Per service.
- Per procedure.
- Per product.
- Per session.
- Custom agreement.

Contoh:
- Dokter Kidinest 70:30.
- Dokter SSC 60:40.
- Terapis berdasarkan tindakan.
- Fee tetap per konsultasi.

Semua persentase harus configurable berdasarkan:
- Dokter.
- Unit bisnis.
- Jenis layanan.
- Tindakan.
- Produk.
- Periode kontrak.
- Effective date.

## 10.1 Fee Calculation

Hitung fee berdasarkan basis yang disepakati:
- Gross service revenue.
- Net service revenue.
- Revenue after discount.
- Exclude tax jika sesuai kebijakan.
- Exclude refunded transaction.
- Adjust cancellation.
- Adjust promotional transaction.

Jangan mengasumsikan semua jenis layanan menggunakan basis perhitungan yang sama.

Simpan:
- Doctor ID.
- Service ID.
- Invoice ID.
- Transaction ID.
- Fee scheme.
- Fee base.
- Percentage.
- Gross amount.
- Adjustment.
- Tax withholding.
- Net payable.
- Calculation period.

## 10.2 Monthly Settlement

Buat:
- Doctor fee statement.
- Therapist commission statement.
- Approval.
- Monthly settlement.
- Partial payout.
- Payment history.
- Outstanding payable.
- Reconciliation with GL.

Contoh akrual:

Debit: Beban Bagi Hasil Dokter.
Kredit: Hutang Jasa Medis Dokter.

Contoh pembayaran setelah pemotongan pajak:

Debit: Hutang Jasa Medis Dokter.
Kredit: Hutang Pajak.
Kredit: Bank.

Besaran potongan pajak harus dihitung oleh tax engine berdasarkan status penerima dan peraturan yang berlaku.

Jangan menggunakan satu rumus pajak untuk semua dokter atau semua jenis hubungan kerja.

---

# 11. TAX MANAGEMENT ENGINE

Buat Tax Management Module yang configurable dan versioned.

Jenis pajak yang perlu dipertimbangkan:
- PPh 21 karyawan.
- PPh 21 dokter atau tenaga ahli.
- PPh 23.
- PPh Final Pasal 4 ayat (2).
- PPN.
- Pajak lain yang relevan.

Buat:
- Tax Type.
- Tax Code.
- Tax Rate.
- Tax Base.
- Tax Calculation Method.
- Tax Category.
- Tax Treatment.
- Tax Effective Date.
- Tax Expiry Date.
- Tax Mapping COA.
- Tax Document.
- Tax Withholding.
- Tax Payable.
- Tax Reconciliation.

## 11.1 Tax Classification

Setiap produk dan layanan dapat memiliki klasifikasi pajak sendiri.

Contoh:
- Pelayanan medis.
- Tindakan medis.
- Estetika non-medis.
- Produk skincare.
- Obat.
- Sewa gedung.
- Jasa vendor.
- Jasa dokter.

Perlakuan pajak tidak boleh hanya ditentukan berdasarkan nama produk atau nama layanan.

Gunakan tax classification dan konfigurasi yang dapat diperbarui.

## 11.2 Tax Calculation

Sediakan:
- Tax inclusive.
- Tax exclusive.
- Withholding tax.
- Multiple tax components jika diperlukan.
- Tax rounding.
- Tax adjustment.
- Tax refund.
- Tax reconciliation.

## 11.3 Tax Reporting

Sediakan:
- Tax Transaction Report.
- PPh 21 Report.
- PPh 23 Report.
- PPh Final Report.
- VAT Input Report.
- VAT Output Report.
- Tax Payable Report.
- Tax Reconciliation Report.
- Export data untuk kebutuhan pelaporan pajak.

Penting:
- Jangan mengunci tarif atau rumus pajak dalam source code.
- Setiap perubahan peraturan harus memiliki effective date.
- Simpan histori perhitungan pajak.
- Verifikasi peraturan dengan sumber resmi DJP dan peraturan pemerintah terbaru.
- Jangan menganggap angka, tarif, kode objek pajak, atau perlakuan pajak dalam blueprint sebagai kebenaran final tanpa verifikasi.
- Jika aturan belum terkonfirmasi, tandai sebagai konfigurasi yang membutuhkan persetujuan konsultan pajak atau akuntan.

---

# 12. CASHIER SHIFT AND CASH CONTROL

Buat modul Cashier Shift Management.

Fitur:
- Open Shift.
- Opening Cash.
- Cashier Assignment.
- Cash Receipt.
- Cash Disbursement.
- Payment Summary.
- Cash Transfer.
- Cash Count.
- Shift Closing.
- Shift Approval.
- Cash Deposit.
- Shift Reconciliation.

Setiap shift menyimpan:
- Shift Number.
- Cashier ID.
- Branch ID.
- Business Unit.
- Cash Account.
- Opening Cash.
- Total Cash Receipt.
- Total Cash Disbursement.
- Expected Cash.
- Actual Cash.
- Cash Difference.
- QRIS Total.
- EDC Total.
- Bank Transfer Total.
- Refund Total.
- Closing Date.
- Approved By.

Rumus:

Expected Cash =
Opening Cash
+ Cash Receipts
- Cash Disbursements
- Cash Refunds
- Cash Deposits.

Actual Cash = hasil perhitungan fisik kasir.

Cash Difference = Actual Cash - Expected Cash.

Selisih harus masuk ke proses investigasi dan approval.

Jangan otomatis membuat jurnal selisih tanpa kebijakan dan persetujuan.

## 12.1 Cashier Accounting

Cash overage dan cash shortage harus menggunakan akun yang dapat dikonfigurasi.

Simpan:
- Cashier.
- Shift.
- Difference.
- Investigation Reason.
- Supporting Evidence.
- Approval.
- Journal Reference.

## 12.2 Cash Deposit to Bank

Dukung:
- Deposit slip.
- Deposit date.
- Bank account.
- Amount.
- Cash source.
- Bank confirmation.
- Outstanding deposit.
- Bank reconciliation.

Transfer kas ke bank bukan pendapatan baru.

---

# 13. PAYMENT GATEWAY AND CLEARING RECONCILIATION

Lengkapi Payment Gateway Clearing.

Dukung:
- QRIS.
- EDC.
- Virtual Account.
- E-Wallet.
- Credit Card.
- Debit Card.
- Payment Gateway Provider.

Pisahkan:
- Gross transaction.
- Gateway fee.
- Tax on fee jika relevan.
- Refund.
- Chargeback.
- Net settlement.
- Settlement date.
- Actual bank receipt.
- Clearing balance.

Sediakan:
- Transaction matching.
- Settlement matching.
- Bank statement matching.
- Unsettled transaction report.
- Outstanding clearing report.
- Settlement discrepancy report.
- Duplicate transaction report.

Setiap transaksi harus memiliki idempotency key dan provider reference.

Jangan mengakui uang sudah masuk ke bank sebelum settlement benar-benar terkonfirmasi.

---

# 14. BANK RECONCILIATION AND TREASURY

Buat Bank Reconciliation yang lengkap.

Dukung:
- Bank statement import.
- CSV/Excel import.
- Manual bank transaction.
- Automatic matching.
- Exact amount matching.
- Reference matching.
- Date tolerance.
- Partial matching jika didukung.
- Unmatched transaction.
- Duplicate detection.
- Bank fee recognition.
- Outstanding deposit.
- Outstanding payment.
- Reconciliation approval.

Simpan:
- Statement ID.
- Bank Transaction ID.
- Internal Transaction ID.
- Match Status.
- Match Amount.
- Difference.
- Reconciliation Date.
- Approved By.

Bank reconciliation tidak boleh mengubah jurnal posted secara langsung.

Selisih harus diselesaikan melalui transaksi atau jurnal yang sah.

---

# 15. ACCOUNTS RECEIVABLE AND CREDIT CONTROL

Lengkapi AR untuk:
- Pasien.
- Perusahaan.
- Asuransi.
- BPJS jika relevan.
- Corporate client.

Fitur:
- Invoice.
- Due Date.
- Payment Terms.
- Credit Limit.
- Partial Payment.
- Payment Allocation.
- Advance Payment.
- Credit Note.
- Debit Note.
- Refund.
- Write-off.
- Collection History.
- Aging Report.
- Overdue Notification.

Aging bucket:
- Current.
- 1–30 days.
- 31–60 days.
- 61–90 days.
- More than 90 days.

Sediakan:
- AR Subledger.
- AR Control Account.
- AR Reconciliation.
- Collection Report.
- Bad Debt Analysis.
- Allowance for Doubtful Accounts jika sesuai kebijakan.

Write-off harus memerlukan approval.

Jangan menghapus piutang hanya untuk membuat saldo menjadi nol.

---

# 16. ACCOUNTS PAYABLE AND PAYMENT CONTROL

Lengkapi AP untuk:
- Supplier.
- Dokter.
- Terapis.
- Karyawan.
- Vendor.
- Pajak.
- Operasional.

Fitur:
- Supplier Invoice.
- Due Date.
- Payment Terms.
- Partial Payment.
- Payment Allocation.
- Credit Note.
- Debit Note.
- Advance Payment.
- Payment Request.
- Payment Approval.
- Payment Batch.
- Payment History.
- AP Aging.
- Outstanding Payable.

Sediakan:
- AP Subledger.
- AP Control Account.
- AP Reconciliation.
- Vendor Statement Reconciliation.
- Duplicate Invoice Detection.
- Payment Due Report.

Dukung approval berdasarkan nominal, unit bisnis, cabang, dan role.

---

# 17. INVENTORY ACCOUNTING INTEGRATION

Hubungkan Inventory Ledger dengan General Ledger.

Dukung:
- FIFO.
- Weighted Average jika dipilih.
- Batch.
- Expired Date.
- FEFO.
- Stock Transfer.
- Stock Adjustment.
- Stock Opname.
- Purchase Receipt.
- Sales Issue.
- Internal Consumption.
- Damaged Stock.
- Expired Stock.
- Stock Return.

Setiap mutasi harus memiliki quantity dan valuation.

Accounting event:
- Purchase receipt.
- Purchase invoice.
- Sales issue.
- COGS.
- Internal consumption.
- Expired inventory.
- Damaged inventory.
- Stock adjustment.
- Purchase return.
- Sales return.

Buat rekonsiliasi:
Inventory Quantity vs Inventory Movement.
Inventory Valuation vs Inventory GL.

Jangan mengurangi stok dua kali akibat invoice dan pembayaran.

Untuk pemakaian BMHP tindakan medis, pastikan stok keluar hanya satu kali dan nilai pemakaian masuk ke akun biaya yang sesuai.

---

# 18. FIXED ASSET ACCOUNTING

Lengkapi Fixed Asset Module.

Fitur:
- Asset Register.
- Asset Category.
- Acquisition.
- Capitalization Policy.
- Useful Life.
- Residual Value.
- Depreciation Method.
- Depreciation Schedule.
- Depreciation Posting.
- Asset Transfer.
- Asset Maintenance.
- Asset Impairment jika relevan.
- Asset Disposal.
- Gain/Loss on Disposal.

Dukung:
- Straight Line.
- Declining Balance.
- Other configurable methods sesuai kebijakan.

Sediakan:
- Fixed Asset Register.
- Asset Movement Report.
- Depreciation Report.
- Asset Disposal Report.
- Fixed Asset Reconciliation.

Pembelian alat medis tidak otomatis menjadi aset tetap. Gunakan capitalization policy dan approval.

---

# 19. PAYROLL AND EMPLOYEE ACCOUNTING

Buat integrasi accounting untuk:
- Salary.
- Allowance.
- Overtime.
- Bonus.
- Employee Deduction.
- BPJS atau kewajiban terkait jika relevan.
- PPh 21.
- Payroll Payable.
- Payroll Payment.

Dukung pemisahan:
- Beban gaji.
- Hutang gaji.
- Hutang pajak.
- Kewajiban karyawan lainnya.

Sediakan payroll journal dan rekonsiliasi payroll dengan General Ledger.

Jangan menyamakan skema pajak dokter dengan skema pajak karyawan.

---

# 20. FINANCIAL REPORTING ENGINE

Semua laporan keuangan harus dihasilkan dari General Ledger dan sub-ledger yang telah direkonsiliasi.

## 20.1 Financial Statements

Wajib tersedia:

1. Income Statement.
2. Balance Sheet.
3. Cash Flow Statement.
4. Statement of Changes in Equity.
5. Trial Balance.
6. General Ledger.
7. Journal Listing.
8. Account Ledger.
9. Cash Book.
10. Bank Book.

## 20.2 Income Statement

Dukung:
- Revenue.
- Sales Discount.
- Sales Return.
- Net Revenue.
- COGS.
- Gross Profit.
- Doctor Sharing Fee.
- Therapist Commission.
- Operating Expense.
- Other Income.
- Other Expense.
- Tax Expense.
- Net Profit.

Dukung laporan berdasarkan:
- Periode.
- Bulan.
- Tahun.
- Business Unit.
- Branch.
- Department.
- Cost Center.
- Service Category.

## 20.3 Balance Sheet

Dukung:
- Current Assets.
- Non-current Assets.
- Current Liabilities.
- Long-term Liabilities.
- Equity.
- Current Year Profit.
- Retained Earnings.

Pastikan:

Assets = Liabilities + Equity.

## 20.4 Cash Flow Statement

Dukung:
- Operating Activities.
- Investing Activities.
- Financing Activities.

Sediakan metode:
- Direct Method.
- Indirect Method.

Cash flow harus didasarkan pada pergerakan kas dan setara kas yang memenuhi definisi kebijakan akuntansi.

Transfer antar rekening sendiri bukan arus kas eksternal dan tidak boleh menggandakan cash flow.

## 20.5 Comparative Reporting

Sediakan:
- Current Period vs Previous Period.
- Current Month vs Previous Month.
- Current Year vs Previous Year.
- Actual vs Budget.
- Unit vs Consolidated.
- Branch Comparison.

## 20.6 Drill-down

Dari angka laporan, pengguna harus bisa menelusuri:

Financial Statement
→ Account
→ General Ledger
→ Journal Entry
→ Source Transaction
→ Invoice / Payment / Inventory Movement.

Terapkan RBAC dan branch access pada setiap tahap.

---

# 21. BUDGETING AND FINANCIAL ANALYSIS

Tambahkan modul budgeting.

Dukung:
- Annual Budget.
- Monthly Budget.
- Revenue Budget.
- Expense Budget.
- Department Budget.
- Cost Center Budget.
- Business Unit Budget.
- Budget Revision.
- Budget Approval.
- Actual vs Budget.
- Budget Variance.

Sediakan financial analysis:
- Revenue Growth.
- Gross Profit Margin.
- Net Profit Margin.
- Operating Expense Ratio.
- AR Turnover.
- AP Outstanding.
- Inventory Turnover.
- Cash Position.
- Revenue per Doctor.
- Revenue per Branch.
- Revenue per Service Category.

Semua rasio harus memiliki formula yang terdokumentasi dan dapat diuji.

---

# 22. FINANCIAL PERIOD CLOSING

Buat proses closing yang terkontrol.

## 22.1 Monthly Closing Checklist

Sebelum periode ditutup, sistem harus memeriksa:

1. Seluruh transaksi operasional sudah diproses.
2. Tidak ada jurnal draft yang tertinggal tanpa keputusan.
3. Seluruh jurnal wajib sudah posted.
4. AR sudah direkonsiliasi.
5. AP sudah direkonsiliasi.
6. Inventory sudah direkonsiliasi.
7. Cash sudah direkonsiliasi.
8. Bank sudah direkonsiliasi.
9. Payment Gateway Clearing sudah direkonsiliasi.
10. Doctor Payable sudah direkonsiliasi.
11. Tax Payable sudah direkonsiliasi.
12. Fixed Asset dan depreciation sudah diperiksa.
13. Accrual dan prepaid sudah diproses.
14. Trial Balance sudah seimbang.
15. Suspense account sudah diperiksa.
16. Laporan keuangan sudah ditinjau.

## 22.2 Closing Workflow

Draft Closing
→ Validation
→ Reconciliation
→ Review
→ Approval
→ Closing Journal jika diperlukan
→ Period Lock.

Dukung:
- Closing Checklist.
- Closing Status.
- Closing Approval.
- Closing Journal.
- Period Lock.
- Reopening Request.
- Reopening Approval.
- Closing History.

Periode yang sudah ditutup tidak boleh menerima posting biasa.

Pembukaan kembali periode harus memerlukan otorisasi khusus dan audit trail.

## 22.3 Year-End Closing

Dukung:
- Closing Revenue Accounts.
- Closing Expense Accounts.
- Profit/Loss Calculation.
- Retained Earnings Transfer.
- Opening New Fiscal Year.
- Carry Forward Balance Sheet Accounts.
- Post-Closing Trial Balance.

Jangan menutup akun neraca seperti akun pendapatan dan beban. Pastikan penutupan akun nominal dan pemindahan hasil periode mengikuti kebijakan akuntansi yang disetujui.

---

# 23. APPROVAL AND INTERNAL CONTROL

Buat configurable approval workflow.

Minimal untuk:
- Manual Journal.
- Journal Reversal.
- Cash Disbursement.
- Bank Payment.
- Refund.
- Purchase Invoice.
- Stock Adjustment.
- Doctor Fee Settlement.
- Tax Payment.
- Fixed Asset Disposal.
- Bad Debt Write-off.
- Period Closing.
- Period Reopening.

Dukung:
- Maker.
- Checker.
- Approver.
- Approval Limit.
- Multi-level Approval.
- Branch-based Approval.
- Business Unit Approval.
- Segregation of Duties.

Buat kontrol agar pembuat transaksi tidak dapat menyetujui transaksi yang sama jika kebijakan segregation of duties melarangnya.

---

# 24. DATABASE DESIGN

Audit dan lengkapi database accounting.

Pastikan tersedia tabel atau struktur yang setara untuk:

ACCOUNTING CORE:
- accounts.
- account_groups.
- accounting_mappings.
- journal_entries.
- journal_lines.
- journal_sequences.
- fiscal_years.
- fiscal_periods.
- opening_balances.
- accounting_policies.
- accounting_dimensions.
- cost_centers.
- journal_approvals.
- journal_reversals.

SUBLEDGER:
- receivables.
- receivable_invoices.
- receivable_allocations.
- payables.
- payable_invoices.
- payable_allocations.
- inventory_valuations.
- fixed_assets.
- asset_depreciations.
- doctor_fee_accruals.
- tax_transactions.

CASH AND BANK:
- cash_accounts.
- cashier_shifts.
- cash_transactions.
- cash_counts.
- cash_deposits.
- bank_accounts.
- bank_transactions.
- bank_statements.
- bank_reconciliations.

TAX:
- tax_types.
- tax_codes.
- tax_rates.
- tax_rules.
- tax_configurations.
- tax_withholdings.
- tax_reports.

DOCTOR SHARING:
- doctor_contracts.
- doctor_fee_schemes.
- doctor_fee_details.
- doctor_fee_settlements.
- therapist_commissions.

PERIOD CLOSING:
- closing_checklists.
- closing_tasks.
- period_closings.
- period_reopenings.

AUDIT:
- audit_logs.
- approval_histories.
- journal_posting_logs.
- reconciliation_logs.

Tabel boleh disesuaikan berdasarkan desain domain, tetapi seluruh kebutuhan data harus terpenuhi.

Gunakan:
- Foreign Key.
- Unique Constraints.
- Check Constraints.
- Indexes.
- Transaction Isolation yang sesuai.
- Created At.
- Updated At.
- Created By.
- Updated By.
- Branch ID jika relevan.
- Business Unit ID jika relevan.

Jangan menyimpan nilai uang menggunakan floating point.

Jangan mengandalkan cache sebagai sumber saldo resmi.

---

# 25. API ENDPOINT REQUIREMENTS

Sediakan endpoint yang lengkap dan konsisten.

ACCOUNTING:
- GET /api/v1/accounting/accounts
- POST /api/v1/accounting/accounts
- PUT /api/v1/accounting/accounts/{id}
- GET /api/v1/accounting/account-mappings
- GET /api/v1/accounting/opening-balances
- POST /api/v1/accounting/opening-balances/import
- POST /api/v1/accounting/opening-balances/validate
- POST /api/v1/accounting/opening-balances/post

JOURNAL:
- GET /api/v1/journals
- POST /api/v1/journals
- GET /api/v1/journals/{id}
- POST /api/v1/journals/{id}/submit
- POST /api/v1/journals/{id}/approve
- POST /api/v1/journals/{id}/reject
- POST /api/v1/journals/{id}/post
- POST /api/v1/journals/{id}/reverse

LEDGER:
- GET /api/v1/reports/general-ledger
- GET /api/v1/reports/account-ledger
- GET /api/v1/reports/trial-balance
- GET /api/v1/reports/subledger-reconciliation

FINANCIAL STATEMENTS:
- GET /api/v1/reports/income-statement
- GET /api/v1/reports/balance-sheet
- GET /api/v1/reports/cash-flow
- GET /api/v1/reports/changes-in-equity

TREASURY:
- POST /api/v1/cashier-shifts/open
- POST /api/v1/cashier-shifts/{id}/count
- POST /api/v1/cashier-shifts/{id}/close
- GET /api/v1/cashier-shifts/{id}/report
- POST /api/v1/cash-deposits
- GET /api/v1/bank-reconciliations

DOCTOR FEE:
- GET /api/v1/doctor-fees
- POST /api/v1/doctor-fees/calculate
- POST /api/v1/doctor-fees/{id}/approve
- POST /api/v1/doctor-fees/{id}/settle

TAX:
- GET /api/v1/tax/configurations
- POST /api/v1/tax/configurations
- GET /api/v1/tax/reports
- GET /api/v1/tax/reconciliation

CLOSING:
- GET /api/v1/accounting/periods
- POST /api/v1/accounting/periods/{id}/validate-closing
- POST /api/v1/accounting/periods/{id}/close
- POST /api/v1/accounting/periods/{id}/reopen

Sesuaikan endpoint dengan REST convention, RBAC, validasi, pagination, dan response standard yang sudah ditentukan.

---

# 26. ACCOUNTING TESTING

Buat unit test dan integration test.

ACCOUNTING:
- Debit equals credit.
- Invalid account rejected.
- Header account cannot be posted.
- Duplicate journal prevented.
- Journal posting atomic.
- Journal reversal correct.
- Posted journal immutable.
- Closed period rejects posting.
- Opening balance balanced.
- Trial Balance balanced.
- Balance Sheet equation correct.
- Closing journal correct.

SUBLEDGER:
- AR agrees with GL.
- AP agrees with GL.
- Inventory valuation agrees with GL.
- Fixed Asset agrees with GL.
- Doctor Payable agrees with GL.
- Tax Payable agrees with GL.

TREASURY:
- Cash expected balance correct.
- Cash difference correct.
- Deposit to bank not recognized as revenue.
- Bank reconciliation matching correct.
- Gateway settlement matches gross, fee, and net.

DOCTOR FEE:
- Fee percentage correct.
- Fixed fee correct.
- Refund adjustment correct.
- Partial payout correct.
- Tax withholding calculation uses correct effective tax configuration.

REPORTING:
- Income Statement correct.
- Balance Sheet balances.
- Cash Flow correct.
- Comparative report correct.
- Consolidated report correct.
- Drill-down totals match source transactions.

Gunakan skenario test dengan angka nyata dalam Rupiah dan tanggal transaksi yang jelas.

---

# 27. ACCOUNTING ACCEPTANCE CRITERIA

Accounting Engine belum boleh dianggap selesai sebelum:

1. COA dapat dikelola secara dinamis.
2. Account Mapping dapat dikonfigurasi.
3. Opening Balance dapat diimpor dan direkonsiliasi.
4. Semua transaksi penting menghasilkan jurnal yang benar.
5. Journal Posting menggunakan transaksi database yang atomik.
6. Tidak ada jurnal duplikat.
7. General Ledger benar.
8. Subsidiary Ledger dapat direkonsiliasi dengan GL.
9. Trial Balance seimbang.
10. Income Statement benar.
11. Balance Sheet seimbang.
12. Cash Flow Statement benar.
13. Multi-unit reporting berjalan.
14. Shared Cost Allocation berjalan.
15. Doctor Sharing Fee dapat dihitung dan direkonsiliasi.
16. Tax Engine menggunakan konfigurasi bertanggal efektif.
17. Cashier Shift dapat ditutup dan direkonsiliasi.
18. Bank Reconciliation berjalan.
19. Payment Gateway Clearing dapat direkonsiliasi.
20. AR dan AP aging benar.
21. Inventory valuation cocok dengan GL.
22. Fixed Asset dan depreciation terintegrasi.
23. Monthly Closing berjalan.
24. Year-End Closing berjalan.
25. Approval dan audit trail berjalan.
26. Semua laporan dapat drill-down.
27. Seluruh test accounting dan integration berhasil.

---

# 28. IMPLEMENTATION INSTRUCTIONS

Lakukan pekerjaan secara bertahap.

## PHASE 1 — ACCOUNTING GAP ANALYSIS

Audit source code dan database yang sudah ada.

Buat tabel:

| Modul | Sudah Ada | Sudah Berfungsi | Belum Lengkap | Risiko |
|---|---|---|---|---|

Jangan menganggap implementasi selesai hanya karena file atau endpoint sudah ada.

Identifikasi:
- Missing entity.
- Missing repository.
- Missing use case.
- Missing database migration.
- Missing journal mapping.
- Missing validation.
- Missing report.
- Missing test.
- Missing accounting control.

## PHASE 2 — ACCOUNTING DESIGN

Berikan:
1. Accounting domain diagram.
2. Accounting ERD.
3. COA lengkap.
4. Journal mapping matrix.
5. Transaction-to-journal flow.
6. Subsidiary ledger design.
7. Tax engine design.
8. Doctor sharing design.
9. Cashier shift design.
10. Financial reporting formula.
11. Period closing workflow.
12. API specification.

## PHASE 3 — CORE ACCOUNTING IMPLEMENTATION

Prioritaskan:
1. COA.
2. Account Mapping.
3. Journal Engine.
4. Journal Posting.
5. Journal Reversal.
6. General Ledger.
7. Trial Balance.
8. Opening Balance.
9. Subledger Reconciliation.

## PHASE 4 — CLINIC ACCOUNTING

Implementasikan:
1. Revenue Recognition.
2. AR.
3. AP.
4. Inventory Accounting.
5. Doctor Sharing.
6. Tax Management.
7. Cashier Shift.
8. Cash Management.
9. Bank Management.
10. Payment Gateway Clearing.

## PHASE 5 — FINANCIAL CONTROL

Implementasikan:
1. Fixed Asset.
2. Accrual.
3. Prepaid Expense.
4. Budgeting.
5. Financial Reporting.
6. Monthly Closing.
7. Year-End Closing.
8. Approval Workflow.
9. Audit Trail.

## PHASE 6 — VALIDATION AND PRODUCTION

Lakukan:
- Integration testing.
- Accounting scenario testing.
- Data reconciliation testing.
- Concurrency testing.
- Security testing.
- Performance testing.
- Export testing.
- Migration testing.

---

# 29. FINAL RULES

1. Jangan membuat implementasi berdasarkan asumsi akuntansi yang belum disetujui.
2. Jangan mengubah histori jurnal yang sudah posted.
3. Jangan membuat jurnal ganda.
4. Jangan membuat saldo akun manual.
5. Jangan menggunakan float untuk nilai uang.
6. Jangan mengakui pendapatan hanya berdasarkan penerimaan uang.
7. Jangan menganggap transfer antar kas atau bank sebagai pendapatan atau beban.
8. Jangan membuat laporan keuangan berdasarkan data hardcode.
9. Jangan membuat laporan dari tabel transaksi operasional secara langsung jika laporan tersebut seharusnya bersumber dari General Ledger.
10. Jangan membuat pajak dengan tarif hardcoded tanpa effective date.
11. Jangan menganggap semua dokter memiliki skema fee dan pajak yang sama.
12. Jangan menggabungkan laporan unit bisnis tanpa dimensi yang jelas.
13. Jangan menghapus data finansial historis.
14. Jangan melakukan posting di periode yang terkunci.
15. Semua transaksi keuangan harus memiliki audit trail.
16. Semua perubahan saldo harus berasal dari transaksi yang sah.
17. Semua proses finansial kritis harus atomic dan idempotent.
18. Semua angka laporan harus dapat ditelusuri sampai sumber transaksi.
19. Jangan membuat pseudo-code jika diminta implementasi.
20. Jangan mengklaim fitur selesai sebelum migration, business logic, API, dan test berhasil.

Jika ada keputusan akuntansi, pajak, atau kebijakan klinik yang belum diketahui, buat daftar pertanyaan dan pilihan konfigurasi. Jangan mengarang jawabannya.

**MULAI DENGAN PHASE 1: AUDIT IMPLEMENTASI DAN GAP ANALYSIS.**

Jangan langsung melakukan perubahan source code sebelum memberikan hasil audit dan meminta persetujuan.

Target akhir adalah Accounting Engine ERP klinik yang akurat, terintegrasi, dapat diaudit, dan siap dikembangkan untuk Kidinest serta Semarang Skin Center.