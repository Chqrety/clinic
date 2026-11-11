# Implementasi Lengkap Sistem Akuntansi ERP Multi-Unit
## Kidinest (Poli Anak & Tumbuh Kembang) & Semarang Skin Center (SSC - Estetika)

### 1. Ringkasan Eksekutif
Sistem akuntansi terintegrasi ERP untuk dual-brand klinik mandiri (**Kidinest** dan **Semarang Skin Center**) telah selesai diimplementasikan secara menyeluruh dengan mengadopsi standar **Domain-Driven Design (DDD) Hexagonal Architecture**, prinsip akuntansi berterima umum (PABU / PSAK), *Double-Entry Bookkeeping*, dan *Accrual Basis*.

---

### 2. Modul & Komponen yang Diimplementasikan

#### A. Multi-Unit Cost Centers & Chart of Accounts (COA)
- **Cost Centers**:
  - `UNIT-KIDINEST`: Pusat Biaya Poli Anak, Tumbuh Kembang, Vaksinasi, dan Farmasi Anak.
  - `UNIT-SSC`: Pusat Biaya Layanan Estetika, Kulit, Treatment Facial, Laser, dan Skincare.
  - `SHARED-BACKOFFICE`: Treasury Kas/Bank, Kantor Pusat, dan Administrasi Bersama.
- **Dimensi Transaksi**: Seluruh baris jurnal (`journal_lines`) terhubung dengan `cost_center_id` untuk memungkinkan pelaporan tersegmentasi per unit dan konsolidasi.
- **COA 4-Level**: Mengakomodasi klasifikasi akun terpisah antara Kidinest (`4-1xxx`, `5-1xxx`) dan Semarang Skin Center (`4-2xxx`, `5-2xxx`).

#### B. Billing & Multi-Journal Atomik 3-Kluster
Pada saat pelunasan tagihan kasir (`billing.Service.PayInvoice`):
1. **Kluster 1 (Pendapatan & Penerimaan)**:
   - Debit: Kas Kasir Kidinest (`1-1111`) / Kas Kasir SSC (`1-1112`) / Bank BCA (`1-1121`) / Kliring QRIS (`1-1130`) / EDC (`1-1131`) / Piutang (`1-1210`).
   - Kredit: Rincian pendapatan per unit:
     - Kidinest: Pendaftaran (`4-1110`), Konsultasi Sp.A (`4-1120`), Farmasi Anak (`4-1150`).
     - SSC: Konsultasi Sp.DVE (`4-2110`), Facial (`4-2120`), Laser (`4-2130`), Skincare (`4-2140`), Admin (`4-2150`).
2. **Kluster 2 (HPP & Pengurangan Persediaan)**:
   - Debit: HPP Obat Kidinest (`5-1110`) atau HPP Skincare SSC (`5-2110`).
   - Kredit: Persediaan Obat (`1-1310`) atau Persediaan Skincare (`1-1330`).
3. **Kluster 3 (Akrual Jasa Medis Dokter)**:
   - Debit: Beban Jasa Medis Dokter (`5-1210` / `5-2210`).
   - Kredit: Hutang Jasa Medis Dokter (`2-1210` / `2-1220`).

#### C. Doctor Share & Monthly Settlement Engine
- **Aturan Pembagian Hasil**: Persentase atau fixed fee per dokter dan kategori tindakan medis (misal: 70% Dokter : 30% Klinik).
- **Kalkulasi Pajak PMK 168/2023 & UU HPP Pasal 17**:
  - $\text{DPP} = 50\% \times \text{Penghasilan Bruto}$.
  - Tarif Progresif Akumulatif:
    - Lapisan 1: s.d Rp 60 Juta (5%)
    - Lapisan 2: > Rp 60 Juta s.d Rp 250 Juta (15%)
    - Lapisan 3: > Rp 250 Juta s.d Rp 500 Juta (25%)
    - Lapisan 4: > Rp 500 Juta s.d Rp 5 Miliar (30%)
    - Lapisan 5: > Rp 5 Miliar (35%)
- **Otomatisasi Jurnal JKK & Bukti Potong**:
  - Debit: Hutang Jasa Medis Dokter (`2-1210` / `2-1220`) sebesar Bruto.
  - Kredit: Bank BCA Operasional (`1-1121`) sebesar Net Fee.
  - Kredit: Hutang PPh 21 Dokter (`2-1501.02`) sebesar PPh 21 Terutang.
  - Otomatis mencatat rekaman *Tax Withholding* untuk kebutuhan e-Bupot DJP.

#### D. Cashier Shift & Cash Control Engine (Treasury)
- **Open Shift**: Merekam modal awal kas kasir (`opening_cash`).
- **Close Shift**: Rekonsiliasi fisik berbasis pecahan uang (denominasi: Rp 100.000, 50.000, 20.000, 10.000, 5.000, 2.000, 1.000, 500).
- **Otomatisasi Jurnal Selisih Kas (Adjusting Journal JC)**:
  - *Overage* (Lebih Kas): Debit Kas Kasir (`1-1111` / `1-1112`), Kredit Pendapatan Selisih Kas Lebih (`4-9200`).
  - *Shortage* (Kurang Kas): Debit Beban Selisih Kas Kurang (`6-1700`), Kredit Kas Kasir (`1-1111` / `1-1112`).
- **Deposit Cash to Bank**:
  - Transfer uang fisik kasir ke Bank BCA (`1-1121`).
  - Auto-posting Jurnal Bank (JB): Debit Bank BCA (`1-1121`), Kredit Kas Kasir (`1-1111`/`1-1112`).

#### E. Tax Management Engine (`clinic_pajak`) & e-Bupot 21/26
- Tracking Bukti Potong PPh 21 Dokter, PPh 21 Karyawan, PPh 23, PPh 4(2), PPN Keluaran.
- Generator CSV resmi format **DJP e-Bupot 21/26** untuk pelaporan SPT Masa PPh 21:
  - Header: `MasaPajak,TahunPajak,Pembetulan,NomorBuktiPotong,TanggalBuktiPotong,NPWP,NIK,NamaWajibPajak,KodeObjekPajak,JumlahPenghasilanBruto,JumlahDPP,Tarif,JumlahPPh`.
- Ringkasan rekapitulasi beban pajak per periode.

#### F. Fixed Asset Management & Straight-Line Depreciation
- Registrasi aset tetap: Alat medis anak (`1-2100`), mesin laser & facial estetika (`1-2200`), renovasi gedung (`1-2300`).
- Formula penyusutan garis lurus bulanan:
  $$\text{Depresiasi Bulanan} = \frac{\text{Harga Perolehan} - \text{Nilai Residu}}{\text{Masa Manfaat (Bulan)}}$$
- Auto-posting Jurnal Penyesuaian (JC) bulanan:
  - Debit: Beban Penyusutan Aset Tetap (`6-1300`).
  - Kredit: Akumulasi Penyusutan Aset Tetap (`1-2900`).

#### G. Pelaporan Keuangan Lengkap
1. **Buku Besar (General Ledger)**: Saldo awal, transaksi kronologis debit/kredit, saldo berjalan (*running balance*).
2. **Neraca Saldo 6-Kolom (Trial Balance)**: Saldo Awal, Mutasi Periode, Saldo Akhir (Debit & Kredit) dengan validasi keseimbangan.
3. **Laporan Laba Rugi (Income Statement)**: Pendapatan Operasional - HPP = Laba Kotor - Beban Operasional = Laba Bersih Periode Berjalan.
4. **Laporan Neraca (Balance Sheet)**: Aset = Liabilitas + Ekuitas (termasuk Laba Periode Berjalan).
5. **Laporan Arus Kas (Cash Flow Statement)**:
   - Arus Kas Aktivitas Operasi (Penerimaan pasien/AR dikurangi pengeluaran operasional/supplier/dokter).
   - Arus Kas Aktivitas Investasi (Capex perolehan aset tetap medis/estetika).
   - Arus Kas Aktivitas Pendanaan (Setoran modal/ekuitas).
   - Rekonsiliasi Saldo Kas Awal vs Saldo Kas Akhir.

#### H. Prepaid Expense & Amortization Engine (`1-1400`)
- **Registrasi Biaya Dibayar di Muka**: Sewa gedung klinik (`6-1200`), asuransi malpraktik, lisensi SIMRS/cloud.
- **Auto-Posting Amortisasi Bulanan (Adjusting Journal JC)**:
  - Debit: `6-1200` Beban Sewa Gedung & Pemeliharaan
  - Kredit: `1-1400` Biaya Dibayar di Muka
- Pengurangan otomatis saldo sisa beban dibayar di muka hingga status `FULLY_AMORTIZED`.

#### I. Budgeting & Variance Analysis Engine
- Penetapan plafon anggaran tahunan & bulanan per akun dan per Cost Center (`UNIT-KIDINEST`, `UNIT-SSC`, `SHARED-BACKOFFICE`).
- Workflow persetujuan anggaran (`DRAFT` -> `APPROVED`).
- Laporan komparasi otomatis *Actual vs Budget Variance*:
  - Realisasi ditarik langsung dari baris jurnal `POSTED` per akun.
  - Menghitung deviasi nominal dan persentase varians.

#### J. Year-End Closing (Tutup Buku Tahunan) Engine
- Validasi pemeriksaan tidak adanya jurnal draft yang tertinggal.
- Penutupan otomatis seluruh akun nominal:
  - Debit seluruh akun Pendapatan (`4-xxxx`).
  - Kredit seluruh akun HPP (`5-xxxx`) dan Beban Operasional (`6-xxxx`).
  - Mentransfer laba/rugi bersih tahun berjalan ke akun `3-1200 Laba Ditahan (Retained Earnings)`.
- Mengunci seluruh periode buku tahun lama.

#### K. Financial Audit Trail Logger (`financial_audit_logs`)
- Pencatatan otomatis riwayat aksi finansial krusial:
  `JOURNAL_POST`, `JOURNAL_REVERSE`, `PERIOD_CLOSE`, `PERIOD_REOPEN`, `YEAR_END_CLOSE`, `POLICY_CHANGE`, `SHIFT_CLOSE`.
- Menyimpan snapshot payload `before_payload` dan `after_payload` (JSONB), User ID, IP Address, dan User Agent.

---

### 3. Matriks Endpoint API Baru

| Domain | Method | Endpoint | Keterangan |
| :--- | :--- | :--- | :--- |
| **Doctor Share** | POST | `/api/v1/doctor-shares/rules` | Buat aturan persentase bagi hasil dokter |
| | GET | `/api/v1/doctor-shares/rules` | Daftar aturan bagi hasil |
| | POST | `/api/v1/doctor-shares/settlements/generate` | Hitung settlement dokter (PMK 168 & PPh 21) |
| | GET | `/api/v1/doctor-shares/settlements` | Daftar settlement bulanan dokter |
| | GET | `/api/v1/doctor-shares/settlements/{id}` | Detail settlement dokter & invoice |
| | POST | `/api/v1/doctor-shares/settlements/{id}/approve` | Persetujuan manajer keuangan |
| | POST | `/api/v1/doctor-shares/settlements/{id}/pay` | Pembayaran via BCA & auto-posting JKK |
| **Treasury** | POST | `/api/v1/treasury/shifts/open` | Buka shift kasir & modal awal |
| | GET | `/api/v1/treasury/shifts/active` | Cek status shift kasir aktif |
| | POST | `/api/v1/treasury/shifts/{id}/close` | Tutup shift, hitung fisik, auto-jurnal selisih |
| | POST | `/api/v1/treasury/shifts/{id}/deposit` | Setor fisik uang ke Bank BCA |
| | GET | `/api/v1/treasury/shifts` | Riwayat pergantian shift kasir |
| **Tax** | POST | `/api/v1/taxes/withholdings` | Buat bukti potong pajak |
| | GET | `/api/v1/taxes/withholdings` | Daftar bukti potong pajak |
| | GET | `/api/v1/taxes/summary` | Rekapitulasi pajak per masa |
| | GET | `/api/v1/taxes/ebupot21/export` | Download CSV format resmi e-Bupot DJP |
| | POST | `/api/v1/taxes/withholdings/{id}/report` | Flag pelaporan SPT |
| **Fixed Asset** | POST | `/api/v1/assets` | Registrasi aset medis/laser & kapitalisasi |
| | GET | `/api/v1/assets` | Daftar aset tetap aktif |
| | POST | `/api/v1/assets/{id}/depreciate` | Eksekusi penyusutan bulanan per aset |
| | POST | `/api/v1/assets/depreciate-batch` | Batch penyusutan seluruh aset aktif |
| | GET | `/api/v1/assets/{id}/depreciations` | Riwayat penyusutan & nilai buku |
| **Prepaid** | POST | `/api/v1/prepaid-expenses` | Registrasi biaya sewa/asuransi dibayar di muka |
| | GET | `/api/v1/prepaid-expenses` | Daftar biaya dibayar di muka aktif |
| | POST | `/api/v1/prepaid-expenses/{id}/amortize` | Amortisasi bulanan per item |
| | POST | `/api/v1/prepaid-expenses/amortize-batch` | Batch amortisasi bulanan seluruh biaya di muka |
| | GET | `/api/v1/prepaid-expenses/{id}/amortizations`| Riwayat jadwal amortisasi |
| **Budget** | POST | `/api/v1/budgets` | Buat anggaran per cost center |
| | GET | `/api/v1/budgets` | Daftar anggaran tahunan |
| | POST | `/api/v1/budgets/{id}/approve` | Persetujuan anggaran |
| | GET | `/api/v1/budgets/variance` | Laporan analisis Actual vs Budget Variance |
| **Year-End Closing**| POST | `/api/v1/periods/year-end-close` | Tutup buku tahunan ke Laba Ditahan |
| **Audit Logs** | GET | `/api/v1/audit-logs` | Telusur audit trail perubahan data keuangan |
| **Reporting** | GET | `/api/v1/reports/cash-flow` | Laporan Arus Kas (Operating, Investing, Financing) |

---

### 4. Hasil Verifikasi & Pengujian
- **Unit & Scenario Testing**:
  - `TestComprehensiveClinicScenario_KidinestAndSSC_EndToEnd`: **PASSED** (Simulasi siklus penuh transaksi pasien Kidinest & SSC, mutasi kasir, penyusutan laser, amortisasi sewa gedung, fee dokter PMK 168, audit balance $\sum \text{Dr} = \sum \text{Cr}$, dan Year-End Closing).
  - `application/prepaid`: **100% Passed** (Register, Capitalization, Monthly Amortization).
  - `application/budget`: **100% Passed** (Budget Creation, Approval, Variance & Favorable/Unfavorable calculations).
  - `application/accounting`: **100% Passed** (GL, Trial Balance, Income Statement, Balance Sheet, Year-End Closing).
  - `application/billing`: **100% Passed** (Dual-unit billing, 3-cluster journal).
  - `application/doctorshare`: **100% Passed** (PMK 168/2023 progressive brackets, payment JKK).
  - `application/treasury`: **100% Passed** (Cashier shift, denominations, overage/shortage JC, bank deposit JB).
  - `application/tax`: **100% Passed** (Withholding, DJP e-Bupot CSV).
  - `application/asset`: **100% Passed** (Straight-line depreciation).
  - `application/receivable`: **100% Passed** (AR aging).
  - `application/payable`: **100% Passed** (AP aging).
- **Kompilasi Biner**:
  - `bin/api` (12 MB): Sukses terkompilasi
  - `bin/worker` (1.6 MB): Sukses terkompilasi
  - `bin/migrate` (1.6 MB): Sukses terkompilasi
| | GET | `/api/v1/treasury/shifts` | Riwayat pergantian shift kasir |
| **Tax** | POST | `/api/v1/taxes/withholdings` | Buat bukti potong pajak |
| | GET | `/api/v1/taxes/withholdings` | Daftar bukti potong pajak |
| | GET | `/api/v1/taxes/summary` | Rekapitulasi pajak per masa |
| | GET | `/api/v1/taxes/ebupot21/export` | Download CSV format resmi e-Bupot DJP |
| | POST | `/api/v1/taxes/withholdings/{id}/report` | Flag pelaporan SPT |
| **Fixed Asset** | POST | `/api/v1/assets` | Registrasi aset medis/laser & kapitalisasi |
| | GET | `/api/v1/assets` | Daftar aset tetap aktif |
| | POST | `/api/v1/assets/{id}/depreciate` | Eksekusi penyusutan bulanan per aset |
| | POST | `/api/v1/assets/depreciate-batch` | Batch penyusutan seluruh aset aktif |
| | GET | `/api/v1/assets/{id}/depreciations` | Riwayat penyusutan & nilai buku |
| **Reporting** | GET | `/api/v1/reports/cash-flow` | Laporan Arus Kas (Operating, Investing, Financing) |

---

### 4. Hasil Verifikasi & Pengujian
- **Unit Testing**:
  - `application/accounting`: **100% Passed** (General Ledger, Trial Balance, Income Statement, Balance Sheet, Journal Posting, Reversal, Fiscal Period Locking).
  - `application/billing`: **100% Passed** (Invoice Creation, Cash POS, Bank Transfer, Overpayment Check).
  - `application/doctorshare`: **100% Passed** (PMK 168/2023 progressive bracket, DPP 50%, Net Fee, Double-Entry Payment JKK).
  - `application/treasury`: **100% Passed** (Open Shift, Physical Cash Denominations, Overage/Shortage Adjusting Journal, Bank Deposit).
  - `application/tax`: **100% Passed** (Withholding Tracking, Summary, e-Bupot 21/26 CSV Generator).
  - `application/asset`: **100% Passed** (Asset Register, Straight-Line Depreciation, Adjusting Depreciation Journal).
  - `application/cashbank`: **100% Passed** (BKM, BKK, TRF).
  - `application/receivable`: **100% Passed** (AR Lifecycle, Aging Report).
  - `application/payable`: **100% Passed** (AP Lifecycle, Aging Report).
- **Kompilasi Biner**:
  - `bin/api` (12 MB): Sukses terkompilasi
  - `bin/worker` (1.6 MB): Sukses terkompilasi
  - `bin/migrate` (1.6 MB): Sukses terkompilasi

