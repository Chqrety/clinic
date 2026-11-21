import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Locale = 'id' | 'en'

export interface LocaleTranslations {
  nav: Record<string, string>
  common: Record<string, string>
  settings: Record<string, string>
  dashboard: Record<string, string>
  coa: Record<string, string>
  journals: Record<string, string>
  cashBank: Record<string, string>
  ledger: Record<string, string>
  trialBalance: Record<string, string>
  profitLoss: Record<string, string>
  balanceSheet: Record<string, string>
  cashFlow: Record<string, string>
  openingBalance: Record<string, string>
  approvals: Record<string, string>
  periods: Record<string, string>
  auditTrail: Record<string, string>
  reports: Record<string, string>
  tax: Record<string, string>
  units: Record<string, string>
  [key: string]: Record<string, string>
}

const accountNameTranslations: Record<string, string> = {
  // Level 1 Categories
  'AKTIVA / ASSETS': 'ASSETS',
  'KEWAJIBAN / LIABILITAS': 'LIABILITIES',
  'MODAL / EKUITAS': 'EQUITY',
  'PENDAPATAN KLINIK / REVENUE': 'CLINICAL REVENUE',
  'HARGA POKOK PENJUALAN (HPP)': 'COST OF GOODS SOLD (COGS)',
  'BEBAN OPERASIONAL / EXPENSES': 'OPERATING EXPENSES',

  // Level 2 Sub-Indexes
  'Aset Lancar': 'Current Assets',
  'Aset Tetap & Peralatan Medis': 'Fixed Assets & Medical Equipment',
  'Kewajiban Jangka Pendek': 'Current Liabilities',
  'Hutang Usaha & Biaya Masih Harus Dibayar': 'Accounts Payable & Accrued Expenses',
  'Hutang Pajak (PPh & PPN)': 'Tax Liabilities (Income Tax & VAT)',
  'Modal Pemilik & Saham Disetor': "Owner's Equity & Paid-in Capital",
  'Pendapatan Layanan Medis Kidinest': 'Kidinest Medical Service Revenue',
  'Pendapatan Layanan Medis SSC': 'SSC Dermatology Medical Service Revenue',
  'HPP Farmasi & Medis': 'Pharmacy & Medical COGS',
  'Beban Operasional Umum': 'General Operating Expenses',

  // Level 3 Sub-Kelompok
  'Kas & Setara Kas': 'Cash & Cash Equivalents',
  'Piutang Pasien & Klaim Asuransi': 'Patient Receivables & Insurance Claims',
  'Persediaan Farmasi & Medis': 'Pharmacy & Medical Supplies Inventory',
  'Peralatan Medis Klinik': 'Clinical Medical Equipment',
  'Hutang Usaha Obat & BMHP': 'Medicine & Consumables Vendor Payables',
  'Hutang Pajak DJP': 'Tax Authority (DJP) Liabilities',
  'Modal & Laba Ditahan': 'Capital & Retained Earnings',
  'Layanan Medis Anak (Kidinest)': 'Pediatric Medical Care (Kidinest)',
  'Layanan Medis Estetika (SSC)': 'Aesthetic Dermatology Care (SSC)',
  'HPP Unit Kidinest': 'Cost of Goods Sold - Kidinest Unit',
  'HPP Semarang Skin Center': 'Cost of Goods Sold - SSC Unit',
  'Beban Honor Jasa Medis Dokter': 'Doctor Medical Honorarium Expense',
  'Beban Ketenagakerjaan & Kantor': 'Employment & Office Expenses',

  // Level 4 Detail Leaf Accounts
  'Kas Kecil (Petty Cash Operasional)': 'Petty Cash (Operational)',
  'Kas Kasir Kidinest': 'Kidinest Cashier Cash',
  'Kas Kasir Semarang Skin Center (SSC)': 'SSC Cashier Cash',
  'Bank BCA Operasional (Giro & EDC)': 'Bank BCA Operational (Checking & EDC)',
  'Bank Mandiri Operasional & Payroll': 'Bank Mandiri Operational & Payroll',
  'Kliring Payment Gateway (QRIS/Tripay/Midtrans)': 'Payment Gateway Clearing (QRIS/Tripay/Midtrans)',
  'Piutang Pasien Umum & Asuransi Swasta': 'General Patient & Private Insurance Receivables',
  'Piutang Klaim BPJS Kesehatan': 'BPJS Health Insurance Receivables',
  'Persediaan Obat & Sirup Anak (Kidinest)': 'Pediatric Medicine & Syrups Inventory (Kidinest)',
  'Persediaan Vaksin & Imunisasi Anak (Kidinest)': 'Vaccines & Immunization Inventory (Kidinest)',
  'Persediaan Produk Skincare & Sunscreen (SSC)': 'Skincare & Sunscreen Products Inventory (SSC)',
  'Persediaan BMHP Estetika & Facial (SSC)': 'Aesthetic & Facial Consumables Inventory (SSC)',
  'Peralatan Medis & USG 4D (Kidinest)': 'Medical Equipment & 4D Ultrasound (Kidinest)',
  'Mesin Laser Nd:YAG & Alkes Estetika (SSC)': 'Nd:YAG Laser Machine & Aesthetic Equipment (SSC)',
  'Akumulasi Penyusutan Peralatan Medis': 'Accumulated Depreciation - Medical Equipment',
  'Hutang Pembelian Obat & Vaksin (Distributor)': 'Medicine & Vaccine Supplier Payables',
  'Hutang Bahan Medis & Skincare (Distributor SSC)': 'Medical & Skincare Supplies Payables (SSC)',
  'Hutang Jasa Medis Dokter Spesialis (Fee Dokter)': 'Specialist Doctor Medical Honorarium Payable',
  'Hutang Gaji & Tunjangan Karyawan': 'Salaries & Employee Benefits Payable',
  'Hutang PPh Pasal 21 Karyawan & Dokter': 'Income Tax Article 21 Payable (Staff & Doctors)',
  'Hutang PPh Pasal 23 Jasa Medis': 'Income Tax Article 23 Payable (Medical Services)',
  'Hutang PPh Final Pasal 4 Ayat 2 (Sewa/Bunga)': 'Final Income Tax Article 4(2) Payable (Rent/Interest)',
  'Hutang PPN Masukan / Keluaran Klinik': 'VAT Payable / Input-Output Tax',
  'Modal Pendirian Klinik (dr. Ratna & Rekan)': 'Clinic Founding Capital (Dr. Ratna & Partners)',
  'Saldo Laba Ditahan (Retained Earnings)': 'Retained Earnings',
  'Laba / Rugi Periode Berjalan': 'Current Period Profit / Loss',
  'Pendapatan Konsultasi Dokter Spesialis Anak': 'Pediatric Specialist Consultation Revenue',
  'Pendapatan Tindakan Medis & Nebulizer Anak': 'Pediatric Medical Procedures & Nebulizer Revenue',
  'Pendapatan Penjualan Obat & Vaksin Anak': 'Pediatric Medicine & Vaccine Sales Revenue',
  'Pendapatan Konsultasi Dokter Spesialis Kulit': 'Dermatology Specialist Consultation Revenue',
  'Pendapatan Treatment Facial & Estetika Medis': 'Facial Treatment & Medical Aesthetic Revenue',
  'Pendapatan Tindakan Laser, Cauter & Peeling': 'Laser, Cautery & Peeling Procedure Revenue',
  'Pendapatan Penjualan Produk Skincare & Sunscreen': 'Skincare & Sunscreen Sales Revenue',
  'HPP Obat & Racikan Anak (Kidinest)': 'Pediatric Medicine & Compounding COGS (Kidinest)',
  'HPP Vaksin & Imunisasi Anak (Kidinest)': 'Pediatric Vaccines & Immunization COGS (Kidinest)',
  'HPP Produk Skincare & Homecare (SSC)': 'Skincare & Homecare Products COGS (SSC)',
  'HPP Bahan BMHP & Treatment Estetika (SSC)': 'Aesthetic Treatment Consumables COGS (SSC)',
  'Beban Honor Jasa Dokter Spesialis Anak (Sp.A)': 'Pediatric Specialist (Sp.A) Doctor Fee Expense',
  'Beban Honor Jasa Dokter Kulit & Terapis Estetika': 'Dermatologist & Aesthetic Therapist Fee Expense',
  'Beban Gaji Karyawan, Perawat & Staf': 'Salaries Expense (Nurses, Staff & Admin)',
  'Beban Listrik, Air & Telepon': 'Utilities Expense (Electricity, Water, Telecom)',
  'Beban Administrasi Bank & EDC': 'Bank Charges & EDC Merchant Discount Rate (MDR)',
  'Beban Pajak Penghasilan Badan (PPh Badan Berjalan)': 'Corporate Income Tax Expense'
}

const translations: Record<Locale, LocaleTranslations> = {
  id: {
    nav: {
      dashboardGroup: 'DASHBOARD',
      dashboard: 'Ikhtisar Keuangan',
      booksGroup: 'BUKU & TRANSAKSI',
      chartOfAccounts: 'Bagan Akun (COA)',
      openingBalance: 'Saldo Awal Akuntansi',
      generalJournal: 'Jurnal Umum (Tabel)',
      cashBank: 'Buku Kas & Bank',
      generalLedger: 'Buku Besar (Ledger)',
      cashierTransactions: 'Transaksi Kasir',
      payroll: 'Slip Gaji Karyawan',
      reportsGroup: 'LAPORAN KEUANGAN (PSAK)',
      trialBalance: 'Neraca Saldo (Trial Balance)',
      profitLoss: 'Laporan Laba Rugi',
      balanceSheet: 'Laporan Neraca',
      cashFlow: 'Arus Kas (Cash Flow)',
      aging: 'Umur Piutang & Hutang',
      taxGroup: 'PERPAJAKAN KLINIK (DJP)',
      sptMasa: 'SPT Masa (Pajak Bulanan)',
      sptTahunan: 'SPT Tahunan Badan (1771)',
      auditGroup: 'KONTROL & AUDIT',
      approvals: 'Approval Bertingkat',
      periods: 'Tutup Buku Periode',
      auditTrail: 'Jejak Audit Sistem',
      settingsGroup: 'PENGATURAN SISTEM',
      languageSettings: 'Pengaturan Bahasa & Regional',
      selectEntity: 'Pilih Entitas Klinik',
      notifications: 'Pemberitahuan Sistem',
      markAllRead: 'Tandai Dibaca',
      logout: 'Keluar'
    },
    common: {
      save: 'Simpan',
      cancel: 'Batal',
      export: 'Ekspor',
      print: 'Cetak',
      refresh: 'Segarkan',
      search: 'Cari',
      filter: 'Filter',
      back: 'Kembali',
      status: 'Status',
      date: 'Tanggal',
      period: 'Periode',
      unit: 'Unit Bisnis',
      branch: 'Cabang',
      actions: 'Aksi',
      detail: 'Rincian',
      close: 'Tutup',
      total: 'Total',
      notes: 'Keterangan',
      loading: 'Memuat data...',
      success: 'Berhasil',
      error: 'Terjadi Kesalahan',
      allUnits: 'Semua Unit (Konsolidasi)',
      active: 'Aktif',
      inactive: 'Nonaktif',
      debit: 'Debet (Dr)',
      credit: 'Kredit (Cr)',
      runningBalance: 'Saldo Berjalan',
      openingBalance: 'Saldo Awal',
      closingBalance: 'Saldo Akhir',
      settings: 'Pengaturan',
      currentMonth: 'Nov 2026 (Berjalan)'
    },
    settings: {
      title: 'Pengaturan Bahasa & Regional',
      subtitle: 'Sesuaikan bahasa antarmuka, terminologi akuntansi (PSAK vs IFRS), dan format pelaporan klinik.',
      langSectionTitle: '1. Pilihan Bahasa Antarmuka (Language Preference)',
      langSectionDesc: 'Pilih bahasa operasional yang diterapkan pada menu navigasi, formulir entri, dan tabel laporan keuangan.',
      idName: 'Bahasa Indonesia',
      idDesc: 'Standar akuntansi Indonesia (PSAK / SAK EP) dan regulasi pajak DJP (PPh 21, 23, Final 4(2), PPN, Bukti Potong, dan formulir 1771).',
      idBadge: 'Standar Nasional (Default)',
      enName: 'English (US / UK)',
      enDesc: 'International Financial Reporting Standards (IFRS / GAAP) for corporate audits, executive dashboards, and international stakeholders.',
      enBadge: 'International Standard',
      previewTitle: '2. Perbandingan Istilah Akuntansi (Terminology Comparison)',
      previewDesc: 'Contoh padanan istilah akuntansi klinik dalam Bahasa Indonesia dan English:',
      indonesianTerm: 'Istilah Bahasa Indonesia (PSAK)',
      englishTerm: 'International English Term (IFRS / GAAP)',
      formatTitle: '3. Format Tanggal & Mata Uang',
      formatDesc: 'Format representasi angka dan penanggalan yang diterapkan pada cetakan dokumen dan export Excel.',
      currencyLabel: 'Mata Uang Pelaporan',
      currencyDesc: 'Rupiah Indonesia (IDR / Rp) dengan pemisah ribuan titik atau koma.',
      dateFormatLabel: 'Format Tanggal',
      saveButton: 'Terapkan Pengaturan Bahasa',
      toastSuccess: 'Bahasa sistem berhasil diubah ke Bahasa Indonesia!'
    },
    dashboard: {
      category: 'Akuntansi & Keuangan',
      title: 'Ikhtisar Finansial & Jurnal Akuntansi',
      createManualJournal: 'Buat Jurnal Manual',
      profitLossReport: 'Laporan Laba Rugi',
      grossRevenue: 'Pendapatan Bruto',
      cashBankLiquid: 'Kas & Bank Likuid',
      liquidSubtitle: 'BCA & Kasir Operasional',
      receivables: 'Piutang Pasien / AR',
      receivablesSubtitle: 'Asuransi & Pasien Umum',
      payables: 'Hutang & Fee Dokter',
      payablesSubtitle: 'Distributor & Jasa Medis',
      netProfit: 'Laba Bersih (Net)',
      netMargin: 'Net Margin',
      unitComparison: 'Perbandingan Unit Bisnis',
      kidinestTitle: 'Kidinest Pediatric Care',
      sscTitle: 'Semarang Skin Center (SSC)',
      netProfitLabel: 'Laba Bersih',
      contribution: 'Kontribusi',
      sscNote: '* Catatan akuntansi: Pendapatan SSC mencakup tindakan medis laser spesialis kulit & penjualan produk skincare estetika.',
      workflowTitle: 'Alur Pemrosesan Jurnal Akuntansi',
      auditReady: 'Audit Ready',
      stepSource: '1. TRANSAKSI SUMBER',
      stepDraft: '2. DRAFT JURNAL',
      stepReview: '3. REVIEW & APPROVAL',
      stepPosted: '4. POSTED KE LEDGER',
      sourceSubtitle: 'Kasir, PO, Farmasi, Fee',
      makerSubtitle: 'Staf Akuntansi',
      checkerSubtitle: 'Supervisor & Manager',
      postedSubtitle: 'Mempengaruhi Laporan',
      unprocessedBadge: 'Belum',
      actionRequiredBadge: 'Perlu Aksi',
      officialBadge: 'Resmi',
      totalDebetPosted: 'Total Debet Resmi Terposting ke Buku Besar Periode Ini:',
      operationalInflowTitle: 'Transaksi Operasional Masuk',
      viewAll: 'Lihat Semua',
      createDraft: 'Buat Draft',
      recentJournalsTitle: 'Jurnal Akuntansi Terkini',
      viewFullList: 'Daftar Lengkap'
    },
    coa: {
      category: 'Bagan Akun Standar',
      subtitle: 'PSAK & Hierarki Akuntansi Klinik',
      title: 'Chart of Accounts (COA)',
      setOpeningBalance: 'Atur Saldo Awal',
      addNewAccount: 'Tambah Akun Baru',
      searchPlaceholder: 'Cari kode akun, nama, atau jenis pajak...',
      allCategories: 'Semua Kategori',
      allSubCategories: 'Semua Sub-Kategori',
      levelAll: 'Semua',
      level1: 'Lvl 1 (Index)',
      level2: 'Lvl 2 (Sub-Index)',
      level3: 'Lvl 3',
      level4: 'Lvl 4 (Detail)',
      statusAll: 'Semua',
      statusActive: 'Aktif',
      statusInactive: 'Nonaktif',
      expandAll: 'Buka Semua',
      collapseAll: 'Ciutkan Semua',
      showingAccounts: 'Menampilkan {count} akun',
      colCode: 'Kode Akun',
      colName: 'Nama Akun (Hierarki COA)',
      colLevel: 'Level',
      colType: 'Klasifikasi Tipe',
      colNormalBalance: 'Saldo Normal',
      colUnit: 'Unit Bisnis',
      colBalance: 'Saldo Saat Ini (Rp)',
      colTax: 'Integrasi Pajak',
      colActions: 'Aksi',
      actionLedger: 'Buku Besar',
      actionEdit: 'Edit Akun',
      actionDeactivate: 'Nonaktifkan Akun',
      actionActivate: 'Aktifkan Akun',
      modalAddTitle: 'Tambah Akun COA Baru',
      modalEditTitle: 'Edit Akun COA',
      fieldCategory: 'Kategori Utama',
      fieldParent: 'Kelompok / Header Induk',
      fieldIsHeader: 'Jadikan Header / Kelompok Akun (Bukan Akun Transaksi)',
      fieldCode: 'Kode Akun Otomatis',
      fieldName: 'Nama Akun',
      fieldType: 'Tipe Akun',
      fieldNormalBalance: 'Saldo Normal',
      fieldUnit: 'Unit Bisnis',
      fieldDesc: 'Deskripsi / Catatan Tambahan',
      btnCancel: 'Batal',
      btnSave: 'Simpan Akun',
      badgeInactive: 'Nonaktif'
    },
    journals: {
      category: 'Buku & Jurnal Akuntansi',
      subtitle: 'Double-Entry General Journal',
      title: 'Daftar Jurnal Umum',
      createManual: 'Buat Jurnal Manual',
      exportExcel: 'Ekspor Excel',
      searchPlaceholder: 'Cari nomor jurnal, referensi, atau keterangan...',
      colJournalNo: 'No. Jurnal',
      colDate: 'Tanggal',
      colRef: 'Referensi',
      colDesc: 'Keterangan Jurnal',
      colEntity: 'Entitas Unit',
      colTotalDebit: 'Total Debet',
      colTotalCredit: 'Total Kredit',
      colStatus: 'Status',
      colMaker: 'Pembuat',
      colActions: 'Aksi',
      actionView: 'Lihat Detail',
      actionApprove: 'Setujui (Approve)',
      actionReject: 'Tolak (Reject)',
      actionVoid: 'Batalkan (Void)'
    },
    cashBank: {
      category: 'Buku Kas & Bank',
      subtitle: 'Pencatatan Rekening & Mutasi Fisik',
      title: 'Buku Kas & Bank',
      recordCashIn: 'Penerimaan Kas (Kas Masuk)',
      recordCashOut: 'Pengeluaran Kas (Kas Keluar)',
      bankTransfer: 'Transfer Kas / Antar Bank',
      searchPlaceholder: 'Cari transaksi kas atau penerima...',
      colDate: 'Tanggal',
      colRef: 'No. Bukti / Reff',
      colAccount: 'Akun Kas / Bank',
      colDesc: 'Keterangan Transaksi',
      colIn: 'Penerimaan (In)',
      colOut: 'Pengeluaran (Out)',
      colBalance: 'Saldo',
      colActions: 'Aksi'
    },
    ledger: {
      category: 'Buku Besar Akun',
      subtitle: 'Rincian Mutasi & Saldo Akun Spesifik',
      title: 'Buku Besar (General Ledger)',
      selectAccount: 'Pilih Akun Buku Besar',
      exportExcel: 'Ekspor Excel',
      normalBalance: 'Saldo Normal',
      beginningBalance: 'Saldo Awal Periode',
      totalDebit: 'Total Debet Periode Ini',
      totalCredit: 'Total Kredit Periode Ini',
      endingBalance: 'Saldo Akhir Periode',
      colDate: 'Tanggal',
      colJournalNo: 'No. Jurnal',
      colRef: 'No. Bukti / Referensi',
      colDesc: 'Keterangan Transaksi',
      colDebit: 'Debet (Rp)',
      colCredit: 'Kredit (Rp)',
      colRunningBalance: 'Saldo Berjalan (Rp)',
      emptyMessage: 'Belum ada transaksi mutasi buku besar untuk akun ini pada periode yang dipilih.',
      createNewJournal: 'Buat Jurnal Baru Untuk Akun Ini'
    },
    trialBalance: {
      category: 'Laporan Keuangan',
      standard: 'PSAK & Kepatuhan Fiskal',
      title: 'Neraca Saldo (Trial Balance)',
      exportExcel: 'Ekspor Excel',
      printReport: 'Cetak',
      balancedStatus: 'Neraca Saldo Seimbang (Balanced) - Total Debet dan Kredit Sesuai Sempurna',
      outOfBalanceStatus: 'Peringatan: Terdapat Selisih Debet & Kredit!',
      colCode: 'Kode Akun',
      colName: 'Nama Akun',
      colType: 'Tipe',
      colUnit: 'Unit',
      colOpening: 'Saldo Awal',
      colMutation: 'Mutasi Periode Berjalan',
      colEnding: 'Saldo Akhir',
      debit: 'Debet',
      credit: 'Kredit',
      total: 'Total Keseluruhan',
      hideZero: 'Sembunyikan Saldo Nol',
      filterUnit: 'Pilih Unit Bisnis',
      filterType: 'Pilih Tipe Akun',
      startDate: 'Dari Tanggal',
      endDate: 'Sampai Tanggal'
    },
    profitLoss: {
      category: 'Laporan Keuangan',
      standard: 'Laporan Laba Rugi Komprehensif',
      title: 'Laporan Laba Rugi (Profit & Loss)',
      exportExcel: 'Ekspor Excel',
      print: 'Cetak Laporan',
      revenue: 'Pendapatan Usaha (Revenue)',
      cogs: 'Harga Pokok Penjualan (HPP)',
      grossProfit: 'Laba Kotor (Gross Profit)',
      operatingExpenses: 'Beban Operasional & Kantor',
      doctorFees: 'Beban Honor Jasa Medis Dokter',
      operatingIncome: 'Laba Operasional (Operating Income)',
      taxes: 'Beban Pajak Penghasilan (PPh)',
      netProfit: 'Laba Bersih Tahun/Periode Berjalan',
      grossMargin: 'Gross Profit Margin',
      netMargin: 'Net Profit Margin',
      period: 'Periode Pelaporan'
    },
    balanceSheet: {
      category: 'Laporan Posisi Keuangan',
      standard: 'PSAK Standar Akuntansi Entitas',
      title: 'Laporan Neraca (Balance Sheet)',
      exportExcel: 'Ekspor Excel',
      print: 'Cetak Neraca',
      assets: 'Aktiva (Aset)',
      currentAssets: 'Aset Lancar',
      fixedAssets: 'Aset Tetap & Peralatan Medis',
      totalAssets: 'Total Aktiva (Aset)',
      liabilities: 'Kewajiban (Liabilitas)',
      currentLiabilities: 'Kewajiban Jangka Pendek',
      longTermLiabilities: 'Kewajiban Jangka Panjang',
      totalLiabilities: 'Total Kewajiban',
      equity: 'Ekuitas (Modal)',
      retainedEarnings: 'Saldo Laba Ditahan',
      currentNetIncome: 'Laba Periode Berjalan',
      totalEquity: 'Total Ekuitas',
      totalLiabilitiesEquity: 'Total Kewajiban & Ekuitas',
      balanced: 'Neraca Seimbang (Aktiva = Kewajiban + Ekuitas)'
    },
    cashFlow: {
      category: 'Laporan Arus Kas',
      standard: 'PSAK 2 - Metode Langsung / Tidak Langsung',
      title: 'Laporan Arus Kas (Cash Flow Statement)',
      exportExcel: 'Ekspor Excel',
      operatingActivities: 'Arus Kas dari Aktivitas Operasi',
      investingActivities: 'Arus Kas dari Aktivitas Investasi',
      financingActivities: 'Arus Kas dari Aktivitas Pendanaan',
      netCashIncrease: 'Kenaikan (Penurunan) Bersih Kas & Setara Kas',
      beginningCash: 'Kas & Setara Kas Awal Periode',
      endingCash: 'Kas & Setara Kas Akhir Periode'
    },
    openingBalance: {
      category: 'Konfigurasi Akuntansi',
      subtitle: 'Pengaturan Saldo Neraca Awal Periode',
      title: 'Saldo Awal Akuntansi (Opening Balance)',
      saveButton: 'Simpan Saldo Awal',
      totalDebit: 'Total Debet Awal',
      totalCredit: 'Total Kredit Awal',
      difference: 'Selisih',
      balancedNote: 'Debet dan Kredit Seimbang (Sesuai Standar Akuntansi)',
      outOfBalanceNote: 'Debet dan Kredit Belum Seimbang!',
      colCode: 'Kode Akun',
      colName: 'Nama Akun COA',
      colType: 'Tipe',
      colNormalBalance: 'Saldo Normal',
      colDebit: 'Saldo Debet (Rp)',
      colCredit: 'Saldo Kredit (Rp)'
    },
    approvals: {
      title: 'Approval Bertingkat',
      subtitle: 'Verifikasi & Otorisasi Transaksi Akuntansi',
      pending: 'Menunggu Persetujuan',
      approved: 'Disetujui',
      rejected: 'Ditolak'
    },
    periods: {
      title: 'Tutup Buku Periode',
      subtitle: 'Kunci Transaksi & Tutup Buku Bulanan/Tahunan',
      statusOpen: 'Terbuka',
      statusClosed: 'Tutup Buku',
      statusLocked: 'Terkunci'
    },
    auditTrail: {
      title: 'Jejak Audit Sistem',
      subtitle: 'Log Keamanan & Riwayat Aktivitas Keuangan'
    },
    reports: {
      generalLedgerTitle: 'Rincian Buku Besar Akun',
      selectAccount: 'Pilih Akun Buku Besar',
      searchPlaceholder: 'Cari transaksi / no jurnal / nominal...',
      normalBalance: 'Saldo Normal',
      noTransactions: 'Belum ada transaksi mutasi buku besar untuk akun ini pada periode yang dipilih.',
      createNewJournal: 'Buat Jurnal Baru Untuk Akun Ini'
    },
    tax: {
      sptMasaTitle: 'SPT Masa (Pajak Bulanan Klinik)',
      sptTahunanTitle: 'SPT Tahunan PPh Badan (Formulir 1771)',
      recordNewSpt: 'Rekam SPT Masa Baru',
      ebupotLinked: 'Terintegrasi e-Bupot & Jurnal Umum'
    },
    units: {
      consolidation: 'Konsolidasi (Seluruh Unit)',
      kidinest: 'Kidinest (Klinik Spesialis Anak)',
      ssc: 'Semarang Skin Center (SSC Dermatologi)'
    }
  },
  en: {
    nav: {
      dashboardGroup: 'DASHBOARD',
      dashboard: 'Financial Overview',
      booksGroup: 'BOOKS & TRANSACTIONS',
      chartOfAccounts: 'Chart of Accounts (COA)',
      openingBalance: 'Accounting Opening Balance',
      generalJournal: 'General Journal Entries',
      cashBank: 'Cash & Bank Ledger',
      generalLedger: 'General Ledger',
      cashierTransactions: 'Cashier Transactions',
      payroll: 'Staff Payroll & Payslips',
      reportsGroup: 'FINANCIAL STATEMENTS (IFRS)',
      trialBalance: 'Trial Balance',
      profitLoss: 'Profit & Loss (Income Statement)',
      balanceSheet: 'Balance Sheet',
      cashFlow: 'Statement of Cash Flows',
      aging: 'A/R & A/P Aging Report',
      taxGroup: 'CLINIC TAX COMPLIANCE',
      sptMasa: 'Monthly Tax Returns (SPT Masa)',
      sptTahunan: 'Annual Corporate Tax (Form 1771)',
      auditGroup: 'CONTROLS & AUDIT',
      approvals: 'Multi-Tier Approvals',
      periods: 'Fiscal Period Closing',
      auditTrail: 'System Audit Trail',
      settingsGroup: 'SYSTEM SETTINGS',
      languageSettings: 'Language & Regional Settings',
      medicalWorkstation: 'Medical Workstation',
      selectEntity: 'Select Clinic Entity',
      notifications: 'System Notifications',
      markAllRead: 'Mark All Read',
      logout: 'Logout'
    },
    common: {
      save: 'Save',
      cancel: 'Cancel',
      export: 'Export',
      print: 'Print',
      refresh: 'Refresh',
      search: 'Search',
      filter: 'Filter',
      back: 'Back',
      status: 'Status',
      date: 'Date',
      period: 'Period',
      unit: 'Business Unit',
      branch: 'Branch',
      actions: 'Actions',
      detail: 'Details',
      close: 'Close',
      total: 'Total',
      notes: 'Notes / Memo',
      loading: 'Loading data...',
      success: 'Success',
      error: 'An Error Occurred',
      allUnits: 'Consolidated (All Units)',
      active: 'Active',
      inactive: 'Inactive',
      debit: 'Debit (Dr)',
      credit: 'Credit (Cr)',
      runningBalance: 'Running Balance',
      openingBalance: 'Opening Balance',
      closingBalance: 'Closing Balance',
      settings: 'Settings',
      currentMonth: 'Nov 2026 (Active)'
    },
    settings: {
      title: 'Language & Regional Settings',
      subtitle: 'Customize interface language, accounting terminology (PSAK vs IFRS), and regional clinic reporting formats.',
      langSectionTitle: '1. Interface Language Selection',
      langSectionDesc: 'Select the operational language applied across navigation menus, voucher forms, and financial report tables.',
      idName: 'Indonesian (Bahasa Indonesia)',
      idDesc: 'Standard Indonesian PSAK/SAK EP and Tax Authority (DJP) compliance (PPh 21, 23, Final 4(2), VAT, e-Bupot, Form 1771).',
      idBadge: 'National Standard (Default)',
      enName: 'English (International)',
      enDesc: 'International Financial Reporting Standards (IFRS / GAAP) for executive dashboards, external audits, and global stakeholders.',
      enBadge: 'International Standard',
      previewTitle: '2. Accounting Terminology Comparison',
      previewDesc: 'Side-by-side comparison of clinical ERP accounting terms in Indonesian and English:',
      indonesianTerm: 'Indonesian Term (PSAK)',
      englishTerm: 'International English Term (IFRS / GAAP)',
      formatTitle: '3. Date & Number Formatting',
      formatDesc: 'Display formats used in on-screen tables, printable vouchers, and Excel export files.',
      currencyLabel: 'Reporting Currency',
      currencyDesc: 'Indonesian Rupiah (IDR / Rp) with customizable thousands separator.',
      dateFormatLabel: 'Date Format',
      saveButton: 'Apply Language Settings',
      toastSuccess: 'System language successfully updated to English!'
    },
    dashboard: {
      category: 'Accounting & Finance',
      title: 'Financial Overview & Accounting Journals',
      createManualJournal: 'Create Manual Journal',
      profitLossReport: 'Profit & Loss Statement',
      grossRevenue: 'Gross Revenue',
      cashBankLiquid: 'Liquid Cash & Bank',
      liquidSubtitle: 'BCA & Cashier Operational',
      receivables: 'Patient Accounts Receivable (A/R)',
      receivablesSubtitle: 'Insurance & Outpatient',
      payables: 'Accounts Payable & Doctor Fees',
      payablesSubtitle: 'Vendors & Medical Specialists',
      netProfit: 'Net Income / Profit',
      netMargin: 'Net Margin',
      unitComparison: 'Business Unit Comparison',
      kidinestTitle: 'Kidinest Pediatric Care',
      sscTitle: 'Semarang Skin Center (SSC)',
      netProfitLabel: 'Net Profit',
      contribution: 'Contribution',
      sscNote: '* Accounting note: SSC revenue includes dermatology laser procedures & aesthetic skincare product sales.',
      workflowTitle: 'Accounting Journal Workflow Pipeline',
      auditReady: 'Audit Ready',
      stepSource: '1. SOURCE TRANSACTIONS',
      stepDraft: '2. DRAFT JOURNAL',
      stepReview: '3. REVIEW & APPROVAL',
      stepPosted: '4. POSTED TO LEDGER',
      sourceSubtitle: 'Cashier, PO, Pharmacy, Fees',
      makerSubtitle: 'Accounting Staff (Maker)',
      checkerSubtitle: 'Supervisor & Finance Manager',
      postedSubtitle: 'Impacts Financial Statements',
      unprocessedBadge: 'Pending',
      actionRequiredBadge: 'Action Required',
      officialBadge: 'Official',
      totalDebetPosted: 'Total Official Debit Posted to General Ledger This Period:',
      operationalInflowTitle: 'Incoming Operational Transactions',
      viewAll: 'View All',
      createDraft: 'Create Draft',
      recentJournalsTitle: 'Recent Accounting Journals',
      viewFullList: 'Full Ledger List'
    },
    coa: {
      category: 'Standard Chart of Accounts',
      subtitle: 'IFRS / GAAP & Clinical Account Hierarchy',
      title: 'Chart of Accounts (COA)',
      setOpeningBalance: 'Set Opening Balance',
      addNewAccount: 'Add New Account',
      searchPlaceholder: 'Search account code, name, or tax category...',
      allCategories: 'All Categories',
      allSubCategories: 'All Sub-Categories',
      levelAll: 'All',
      level1: 'Lvl 1 (Index)',
      level2: 'Lvl 2 (Sub-Index)',
      level3: 'Lvl 3',
      level4: 'Lvl 4 (Detail)',
      statusAll: 'All',
      statusActive: 'Active',
      statusInactive: 'Inactive',
      expandAll: 'Expand All',
      collapseAll: 'Collapse All',
      showingAccounts: 'Showing {count} accounts',
      colCode: 'Account Code',
      colName: 'Account Name (COA Hierarchy)',
      colLevel: 'Level',
      colType: 'Classification Type',
      colNormalBalance: 'Normal Balance',
      colUnit: 'Business Unit',
      colBalance: 'Current Balance (IDR)',
      colTax: 'Tax Integration',
      colActions: 'Actions',
      actionLedger: 'General Ledger',
      actionEdit: 'Edit Account',
      actionDeactivate: 'Deactivate Account',
      actionActivate: 'Activate Account',
      modalAddTitle: 'Add New COA Account',
      modalEditTitle: 'Edit COA Account',
      fieldCategory: 'Main Category',
      fieldParent: 'Parent Group / Header',
      fieldIsHeader: 'Make Header / Group Account (Non-Transactional)',
      fieldCode: 'Account Code (Auto)',
      fieldName: 'Account Name',
      fieldType: 'Account Type',
      fieldNormalBalance: 'Normal Balance',
      fieldUnit: 'Business Unit',
      fieldDesc: 'Description / Remarks',
      btnCancel: 'Cancel',
      btnSave: 'Save Account',
      badgeInactive: 'Inactive'
    },
    journals: {
      category: 'Books & Accounting Journals',
      subtitle: 'Double-Entry General Journal',
      title: 'General Journal Entries',
      createManual: 'Create Manual Journal',
      exportExcel: 'Export Excel',
      searchPlaceholder: 'Search journal number, reference, or description...',
      colJournalNo: 'Journal #',
      colDate: 'Date',
      colRef: 'Reference',
      colDesc: 'Description',
      colEntity: 'Unit Entity',
      colTotalDebit: 'Total Debit',
      colTotalCredit: 'Total Credit',
      colStatus: 'Status',
      colMaker: 'Maker',
      colActions: 'Actions',
      actionView: 'View Details',
      actionApprove: 'Approve',
      actionReject: 'Reject',
      actionVoid: 'Void Journal'
    },
    cashBank: {
      category: 'Cash & Bank Books',
      subtitle: 'Account Balances & Mutation Records',
      title: 'Cash & Bank Ledger',
      recordCashIn: 'Receive Cash (Cash In)',
      recordCashOut: 'Disburse Cash (Cash Out)',
      bankTransfer: 'Cash / Inter-Bank Transfer',
      searchPlaceholder: 'Search cash transaction or payee...',
      colDate: 'Date',
      colRef: 'Voucher / Ref #',
      colAccount: 'Cash / Bank Account',
      colDesc: 'Transaction Description',
      colIn: 'Cash In',
      colOut: 'Cash Out',
      colBalance: 'Balance',
      colActions: 'Actions'
    },
    ledger: {
      category: 'General Ledger',
      subtitle: 'Account Specific Mutations & Balances',
      title: 'General Ledger Account Detail',
      selectAccount: 'Select Ledger Account',
      exportExcel: 'Export Excel',
      normalBalance: 'Normal Balance',
      beginningBalance: 'Beginning Balance',
      totalDebit: 'Total Debits This Period',
      totalCredit: 'Total Credits This Period',
      endingBalance: 'Ending Balance',
      colDate: 'Date',
      colJournalNo: 'Journal #',
      colRef: 'Reference #',
      colDesc: 'Transaction Description',
      colDebit: 'Debit (IDR)',
      colCredit: 'Credit (IDR)',
      colRunningBalance: 'Running Balance (IDR)',
      emptyMessage: 'No general ledger transactions found for this account in the selected period.',
      createNewJournal: 'Create New Journal For This Account'
    },
    trialBalance: {
      category: 'Financial Statements',
      standard: 'IFRS & Fiscal Compliance',
      title: 'Trial Balance',
      exportExcel: 'Export Excel',
      printReport: 'Print',
      balancedStatus: 'Trial Balance is Balanced - Total Debits and Credits Match Perfectly',
      outOfBalanceStatus: 'Warning: Debit and Credit Discrepancy Detected!',
      colCode: 'Account Code',
      colName: 'Account Name',
      colType: 'Type',
      colUnit: 'Unit',
      colOpening: 'Opening Balance',
      colMutation: 'Period Mutation',
      colEnding: 'Ending Balance',
      debit: 'Debit',
      credit: 'Credit',
      total: 'Grand Total',
      hideZero: 'Hide Zero Balances',
      filterUnit: 'Select Business Unit',
      filterType: 'Select Account Type',
      startDate: 'From Date',
      endDate: 'To Date'
    },
    profitLoss: {
      category: 'Financial Statements',
      standard: 'Comprehensive Income Statement',
      title: 'Profit & Loss Statement (Income Statement)',
      exportExcel: 'Export Excel',
      print: 'Print Statement',
      revenue: 'Operating Revenue',
      cogs: 'Cost of Goods Sold (COGS)',
      grossProfit: 'Gross Profit',
      operatingExpenses: 'Operating & Administrative Expenses',
      doctorFees: 'Doctor & Specialist Medical Honorarium',
      operatingIncome: 'Operating Income',
      taxes: 'Corporate Income Tax (CIT)',
      netProfit: 'Net Income / Profit for the Period',
      grossMargin: 'Gross Profit Margin',
      netMargin: 'Net Profit Margin',
      period: 'Reporting Period'
    },
    balanceSheet: {
      category: 'Statement of Financial Position',
      standard: 'IFRS International Accounting Standard',
      title: 'Balance Sheet (Statement of Financial Position)',
      exportExcel: 'Export Excel',
      print: 'Print Balance Sheet',
      assets: 'Assets',
      currentAssets: 'Current Assets',
      fixedAssets: 'Fixed Assets & Medical Equipment',
      totalAssets: 'Total Assets',
      liabilities: 'Liabilities',
      currentLiabilities: 'Current Liabilities',
      longTermLiabilities: 'Long-Term Liabilities',
      totalLiabilities: 'Total Liabilities',
      equity: 'Equity',
      retainedEarnings: 'Retained Earnings',
      currentNetIncome: 'Current Period Net Income',
      totalEquity: 'Total Equity',
      totalLiabilitiesEquity: 'Total Liabilities & Equity',
      balanced: 'Balance Sheet is Balanced (Assets = Liabilities + Equity)'
    },
    cashFlow: {
      category: 'Statement of Cash Flows',
      standard: 'IAS 7 - Cash Flow Statement',
      title: 'Statement of Cash Flows',
      exportExcel: 'Export Excel',
      operatingActivities: 'Cash Flows from Operating Activities',
      investingActivities: 'Cash Flows from Investing Activities',
      financingActivities: 'Cash Flows from Financing Activities',
      netCashIncrease: 'Net Increase (Decrease) in Cash & Cash Equivalents',
      beginningCash: 'Beginning Cash & Cash Equivalents',
      endingCash: 'Ending Cash & Cash Equivalents'
    },
    openingBalance: {
      category: 'Accounting Setup',
      subtitle: 'Initial Balance Configuration',
      title: 'Accounting Opening Balance',
      saveButton: 'Save Opening Balance',
      totalDebit: 'Total Opening Debit',
      totalCredit: 'Total Opening Credit',
      difference: 'Difference',
      balancedNote: 'Debits and Credits are Balanced (GAAP Compliant)',
      outOfBalanceNote: 'Debits and Credits are Out of Balance!',
      colCode: 'Account Code',
      colName: 'Account Name',
      colType: 'Type',
      colNormalBalance: 'Normal Balance',
      colDebit: 'Debit Balance (IDR)',
      colCredit: 'Credit Balance (IDR)'
    },
    approvals: {
      title: 'Multi-Tier Approvals',
      subtitle: 'Verification & Authorization of Accounting Entries',
      pending: 'Pending Approval',
      approved: 'Approved',
      rejected: 'Rejected'
    },
    periods: {
      title: 'Fiscal Period Closing',
      subtitle: 'Lock Transactions & Close Monthly/Annual Books',
      statusOpen: 'Open',
      statusClosed: 'Closed',
      statusLocked: 'Locked'
    },
    auditTrail: {
      title: 'System Audit Trail',
      subtitle: 'Security Log & Financial Activity History'
    },
    reports: {
      generalLedgerTitle: 'General Ledger Account Detail',
      selectAccount: 'Select Ledger Account',
      searchPlaceholder: 'Search transactions / journal # / amount...',
      normalBalance: 'Normal Balance',
      noTransactions: 'No general ledger transactions found for this account in the selected period.',
      createNewJournal: 'Create New Journal For This Account'
    },
    tax: {
      sptMasaTitle: 'Monthly Tax Returns (SPT Masa)',
      sptTahunanTitle: 'Annual Corporate Income Tax (Form 1771)',
      recordNewSpt: 'Record New Monthly Return',
      ebupotLinked: 'e-Bupot & General Journal Integrated'
    },
    units: {
      consolidation: 'Consolidated (All Clinic Units)',
      kidinest: 'Kidinest Pediatric Center',
      ssc: 'Semarang Skin Center (SSC Dermatology)'
    }
  }
}

export const useLocaleStore = defineStore('locale', () => {
  // Read saved locale or default to 'id'
  const savedLocale = (typeof localStorage !== 'undefined' ? localStorage.getItem('clinic_locale') : null) as Locale | null
  const currentLocale = ref<Locale>(savedLocale === 'en' ? 'en' : 'id')

  // Set HTML document lang tag
  if (typeof document !== 'undefined') {
    document.documentElement.lang = currentLocale.value
  }

  function setLocale(newLocale: Locale) {
    currentLocale.value = newLocale
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('clinic_locale', newLocale)
    }
    if (typeof document !== 'undefined') {
      document.documentElement.lang = newLocale
    }
  }

  function toggleLocale() {
    setLocale(currentLocale.value === 'id' ? 'en' : 'id')
  }

  // Reactive translation lookup: t('nav.generalLedger') or t('common.save')
  function t(path: string, fallback?: string): string {
    const parts = path.split('.')
    let current: any = translations[currentLocale.value]

    for (const part of parts) {
      if (current && typeof current === 'object' && part in current) {
        current = current[part]
      } else {
        return fallback || path
      }
    }

    return typeof current === 'string' ? current : fallback || path
  }

  // Multi-word phrase regex rules for dynamic auto-translation
  const dynamicPhrases: [RegExp, string][] = [
    [/Harga Pokok Penjualan/gi, 'Cost of Goods Sold'],
    [/Beban Operasional/gi, 'Operating Expenses'],
    [/Kas & Setara Kas/gi, 'Cash & Cash Equivalents'],
    [/Piutang Pasien/gi, 'Patient Receivables'],
    [/Hutang Usaha/gi, 'Accounts Payable'],
    [/Hutang Pajak/gi, 'Tax Liabilities'],
    [/Laba Bersih/gi, 'Net Income'],
    [/Laba Kotor/gi, 'Gross Profit'],
    [/Laba Ditahan/gi, 'Retained Earnings'],
    [/Jasa Medis/gi, 'Medical Services'],
    [/Dokter Spesialis Anak/gi, 'Pediatric Specialist'],
    [/Dokter Spesialis Kulit/gi, 'Dermatology Specialist'],
    [/Persediaan Farmasi/gi, 'Pharmacy Inventory'],
    [/Peralatan Medis/gi, 'Medical Equipment'],
    [/Transaksi Kasir/gi, 'Cashier Transactions'],
    [/Kasir Operasional/gi, 'Cashier Operations'],
    [/Tutup Buku/gi, 'Period Closing'],
    [/Saldo Awal/gi, 'Beginning Balance'],
    [/Saldo Akhir/gi, 'Ending Balance'],
    [/Saldo Berjalan/gi, 'Running Balance'],
    [/Buku Besar/gi, 'General Ledger'],
    [/Neraca Saldo/gi, 'Trial Balance'],
    [/Jurnal Umum/gi, 'General Journal'],
    [/Beban Honor/gi, 'Honorarium Expense'],
    [/Beban Gaji/gi, 'Salaries Expense'],
    [/Beban Listrik/gi, 'Electricity Expense'],
    [/Beban Sewa/gi, 'Rent Expense']
  ]

  // Individual token glossary for dynamic automatic translation
  const dynamicGlossary: Record<string, string> = {
    'Aktiva': 'Assets',
    'Aset': 'Assets',
    'Kewajiban': 'Liabilities',
    'Hutang': 'Payables',
    'Utang': 'Payables',
    'Piutang': 'Receivables',
    'Modal': 'Equity',
    'Ekuitas': 'Equity',
    'Pendapatan': 'Revenue',
    'Beban': 'Expense',
    'Biaya': 'Expense',
    'Kas': 'Cash',
    'Bank': 'Bank',
    'Persediaan': 'Inventory',
    'Peralatan': 'Equipment',
    'Penyusutan': 'Depreciation',
    'Akumulasi': 'Accumulated',
    'Obat': 'Medicine',
    'Vaksin': 'Vaccines',
    'Gaji': 'Salaries',
    'Karyawan': 'Staff',
    'Perawat': 'Nurses',
    'Dokter': 'Doctor',
    'Pasien': 'Patient',
    'Klinik': 'Clinic',
    'Cabang': 'Branch',
    'Umum': 'General',
    'Lancar': 'Current',
    'Tetap': 'Fixed',
    'Pajak': 'Tax',
    'Penjualan': 'Sales',
    'Pembelian': 'Purchases',
    'Operasional': 'Operational',
    'Klaim': 'Claims',
    'Asuransi': 'Insurance',
    'Debet': 'Debit',
    'Kredit': 'Credit',
    'Saldo': 'Balance',
    'Mutasi': 'Mutation',
    'Total': 'Total',
    'Konsolidasi': 'Consolidation',
    'Semua': 'All',
    'Aktif': 'Active',
    'Nonaktif': 'Inactive',
    'Kasir': 'Cashier',
    'Transaksi': 'Transactions'
  }

  // Smart client-side dynamic auto-translation (zero server load, instant execution)
  function autoTranslate(text: string): string {
    if (currentLocale.value === 'id' || !text) return text

    // 1. Exact match from pre-defined master accounts
    if (accountNameTranslations[text]) {
      return accountNameTranslations[text]
    }

    let result = text

    // 2. Multi-word phrases translation
    for (const [regex, replacement] of dynamicPhrases) {
      result = result.replace(regex, replacement)
    }

    // 3. Word-by-word token replacement
    result = result.replace(/\b([a-zA-ZÀ-ÿ]+)\b/g, (match) => {
      const cap = match.charAt(0).toUpperCase() + match.slice(1).toLowerCase()
      if (dynamicGlossary[match]) return dynamicGlossary[match]
      if (dynamicGlossary[cap]) return dynamicGlossary[cap]
      return match
    })

    return result
  }

  // Helper to translate COA account names and categories dynamically
  function translateAccountName(originalName: string): string {
    return autoTranslate(originalName)
  }

  const isIndonesian = computed(() => currentLocale.value === 'id')
  const isEnglish = computed(() => currentLocale.value === 'en')

  return {
    currentLocale,
    isIndonesian,
    isEnglish,
    setLocale,
    toggleLocale,
    t,
    autoTranslate,
    translateAccountName,
    translations
  }
})
