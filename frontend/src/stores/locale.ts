import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Locale = 'id' | 'en'

export interface LocaleTranslations {
  nav: Record<string, string>
  common: Record<string, string>
  settings: Record<string, string>
  reports: Record<string, string>
  tax: Record<string, string>
  units: Record<string, string>
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
      cashierTransactions: 'Transaksi Kasir Otomatis',
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
      languageSettings: 'Pengaturan Bahasa & Regional'
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
      closingBalance: 'Saldo Akhir'
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
      cashierTransactions: 'Automated Cashier Sync',
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
      languageSettings: 'Language & Regional Settings'
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
      closingBalance: 'Closing Balance'
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

  const isIndonesian = computed(() => currentLocale.value === 'id')
  const isEnglish = computed(() => currentLocale.value === 'en')

  return {
    currentLocale,
    isIndonesian,
    isEnglish,
    setLocale,
    toggleLocale,
    t,
    translations
  }
})
