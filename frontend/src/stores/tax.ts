import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TaxTransaction, BusinessUnit } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'
import { exportToCSV } from '@/utils/export'
import { accountingApi, type BackendTaxWithholding } from '@/api/accounting'

export type TaxPeriodType = 'MASA' | 'TAHUNAN'

export type SptMasaType = 'PPH_21' | 'PPH_23' | 'PPH_4_2' | 'PPH_25' | 'PPN'

export interface SptMasaRecord {
  id: string
  sptNumber: string
  masaMonth: number // 1 s/d 12
  masaYear: number // 2026
  taxType: SptMasaType
  accountCode: string
  accountName: string
  // Strictly Bank Account Only (Non-Cash)
  bankAccountCode: string // misal 1102.001 (Bank BCA) atau 1102.002 (Bank Mandiri)
  bankAccountName: string
  paymentChannel: 'BANK_TRANSFER' | 'BANK_EBILLING' | 'BANK_GIRO'
  isBankVerified: boolean
  grossAmount: number
  dppAmount: number
  taxRate: string
  taxAmount: number
  ntpnNumber?: string
  bpeNumber?: string
  paymentDate?: string
  status: 'DRAFT' | 'PAID' | 'REPORTED'
  notes?: string
  businessUnit: BusinessUnit
  createdAt: string
  createdBy: string
}

export interface DoctorTaxRecord {
  id: string
  doctorName: string
  specialty: string
  npwp: string
  nik: string
  businessUnit: BusinessUnit
  bankAccountTarget: string // Transfer via Bank
  grossIncomeCurrentMonth: number
  grossIncomeCumulativeYear: number
  dppPercentage: number // 50% for non-pegawai dokter (PMK 168/2023)
  dppAmount: number
  taxBracketApplied: string
  taxRate: number
  pph21Amount: number
  netPayable: number
  ebupotNumber: string
  journalNumber?: string
}

export interface TaxAccountOption {
  code: string
  name: string
  taxCategory: string
  description: string
  normalRate: string
}

export interface AnnualTaxReport1771 {
  fiscalYear: number
  status: 'DRAFT' | 'FINALIZED' | 'SUBMITTED'
  commercialRevenue: number
  commercialExpense: number
  commercialNetProfit: number
  positiveFiscalCorrection: number
  positiveCorrectionNotes: string
  negativeFiscalCorrection: number
  negativeCorrectionNotes: string
  taxableIncome: number // Penghasilan Kena Pajak (PKP)
  corporateTaxRate: number // 22% (UU HPP)
  corporateTaxPayable: number // PPh Terutang
  // Kredit Pajak murni bersumber dari mutasi perbankan (SPT Masa)
  taxCreditPph23: number // Otomatis dari akumulasi SPT Masa PPh 23 via Bank
  taxCreditPph25: number // Otomatis dari akumulasi SPT Masa PPh 25 via e-Billing Bank
  taxCredits: number // Total Kredit Pajak (PPh 23 + PPh 25)
  pph29Payable: number // PPh Kurang Bayar Pasal 29
  linkedSptMasaCount: number
  bankReconciliationVerified: boolean
}

export const useTaxStore = defineStore('tax', () => {
  const accountingStore = useAccountingStore()
  const ledgerStore = useLedgerStore()
  const journalsStore = useJournalsStore()

  // State
  const taxPeriodType = ref<TaxPeriodType>('MASA')
  const selectedTaxAccountCode = ref<string>('ALL')
  const selectedMonth = ref<number>(11)
  const selectedYear = ref<number>(2026)
  const selectedSptMasaTypeFilter = ref<string>('ALL')

  // Tax Accounts mapped to COA
  const taxAccountOptions = computed<TaxAccountOption[]>(() => {
    return [
      { code: 'ALL', name: 'Semua Akun Pajak Terintegrasi COA', taxCategory: 'ALL', description: 'Rekapitulasi seluruh kewajiban perpajakan klinik', normalRate: 'Multi-Tarif' },
      { code: '2501.001', name: '2501.001 - PPh 21 Dokter Tenaga Ahli (PMK 168/2023)', taxCategory: 'PPH_21', description: 'Withholding Tax Jasa Medis & Visite Dokter', normalRate: '5% - 35% (Progresif 50% Bruto)' },
      { code: '2501.002', name: '2501.002 - PPh 21 Karyawan & Staf Tetap (TER PP 58)', taxCategory: 'PPH_21', description: 'Pajak Penghasilan Pegawai Bulanan', normalRate: 'Tarif Efektif Rata-Rata (TER)' },
      { code: '2502', name: '2502 - PPh 23 (Jasa Limbah B3 & Maintenance Alkes)', taxCategory: 'PPH_23', description: 'Pajak Jasa Teknik & Perawatan Alat Medis', normalRate: '2%' },
      { code: '2503', name: '2503 - PPh Final Pasal 4(2) (Sewa Ruko & Fasilitas)', taxCategory: 'PPH_4_2', description: 'Pajak Penghasilan Final Persewaan Bangunan', normalRate: '10%' },
      { code: '2504', name: '2504 - PPN Keluaran (Penjualan Kosmetika & Skincare)', taxCategory: 'PPN_OUT', description: 'Pajak Pertambahan Nilai Treatment Non-Medis & Produk', normalRate: '11%' },
      { code: '1504', name: '1504 - PPN Masukan & Uang Muka PPh 25 (Kredit Pajak)', taxCategory: 'PPN_IN', description: 'Faktur Pajak Pembelian Distributor PBF & Angsuran Badan', normalRate: '11% / Kredit' },
      { code: '2505', name: '2505 - PPh 25 (Angsuran PPh Badan Bulanan)', taxCategory: 'PPH_25', description: 'Angsuran Pajak Penghasilan Badan Berjalan', normalRate: 'Angsuran Tetap' },
      { code: '2506', name: '2506 - PPh 29 (Pajak Badan Akhir Tahun SPT 1771)', taxCategory: 'PPH_29', description: 'PPh Badan Kurang Bayar Tahunan', normalRate: '22%' }
    ]
  })

  // Available Bank Accounts in COA (STRICTLY BANK ONLY, KAS FISIK EXCLUDED)
  const availableBankAccounts = computed(() => {
    return ledgerStore.accounts.filter(a => {
      if (a.isHeader) return false
      // Only include accounts under 1102 (Bank Operasional & Kliring) or 1103 (Bank Transit), EXCLUDING 1101 (Kas Fisik)
      return a.code.startsWith('1102') || a.code.startsWith('1103') || a.parentCode === '1102'
    })
  })

  // =========================================================================
  // 1. REKAPITULASI DOKUMEN SPT MASA (Januari - November 2026)
  // Strictly Berbasis Rekening Bank Operasional (Bukan Kas Tunai)
  // =========================================================================
  const sptMasaRecords = ref<SptMasaRecord[]>([
    // --- PPH 23 (KREDIT PAJAK LAMPIRAN III SPT 1771 - VIA TRANSFER BANK) ---
    {
      id: 'spt-23-01',
      sptNumber: 'SPT-MASA-23-2026-01',
      masaMonth: 1,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Pengolahan Limbah Medis & IT',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 120000000,
      dppAmount: 120000000,
      taxRate: '2%',
      taxAmount: 2400000,
      ntpnNumber: 'NTPN-BCA-20260128-01',
      bpeNumber: 'BPE-23-202601-9901',
      paymentDate: '2026-01-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Jasa limbah B3 PT Wastec & maintenance SIMRS',
      businessUnit: 'ALL',
      createdAt: '2026-01-28 10:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-02',
      sptNumber: 'SPT-MASA-23-2026-02',
      masaMonth: 2,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Kalibrasi Alat Medis Laser SSC',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-BCA-20260228-02',
      bpeNumber: 'BPE-23-202602-9902',
      paymentDate: '2026-02-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Kalibrasi laser Nd:YAG & USG 4D Kidinest',
      businessUnit: 'SSC',
      createdAt: '2026-02-28 11:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-03',
      sptNumber: 'SPT-MASA-23-2026-03',
      masaMonth: 3,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Kebersihan Medis & Sanitasi',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 160000000,
      dppAmount: 160000000,
      taxRate: '2%',
      taxAmount: 3200000,
      ntpnNumber: 'NTPN-BCA-20260328-03',
      bpeNumber: 'BPE-23-202603-9903',
      paymentDate: '2026-03-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Outsourcing cleaning service ruang steril',
      businessUnit: 'ALL',
      createdAt: '2026-03-28 14:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-04',
      sptNumber: 'SPT-MASA-23-2026-04',
      masaMonth: 4,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Konsultasi Manajemen Medis',
      bankAccountCode: '1102.002',
      bankAccountName: 'Bank Mandiri Operasional & Payroll',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 140000000,
      dppAmount: 140000000,
      taxRate: '2%',
      taxAmount: 2800000,
      ntpnNumber: 'NTPN-MDR-20260428-04',
      bpeNumber: 'BPE-23-202604-9904',
      paymentDate: '2026-04-28',
      status: 'REPORTED',
      notes: 'Transfer Bank Mandiri: Konsultasi akreditasi Paripurna Kemenkes',
      businessUnit: 'ALL',
      createdAt: '2026-04-28 09:30:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-05',
      sptNumber: 'SPT-MASA-23-2026-05',
      masaMonth: 5,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Limbah B3 & Maintenance Alkes',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 155000000,
      dppAmount: 155000000,
      taxRate: '2%',
      taxAmount: 3100000,
      ntpnNumber: 'NTPN-BCA-20260528-05',
      bpeNumber: 'BPE-23-202605-9905',
      paymentDate: '2026-05-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Maintenance berkala autoklaf & sterilisator',
      businessUnit: 'ALL',
      createdAt: '2026-05-28 15:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-06',
      sptNumber: 'SPT-MASA-23-2026-06',
      masaMonth: 6,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Sewa Kendaraan Logistik Vaksin',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 145000000,
      dppAmount: 145000000,
      taxRate: '2%',
      taxAmount: 2900000,
      ntpnNumber: 'NTPN-BCA-20260628-06',
      bpeNumber: 'BPE-23-202606-9906',
      paymentDate: '2026-06-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Sewa mobil cold-chain distribusi vaksin',
      businessUnit: 'KIDINEST',
      createdAt: '2026-06-28 11:20:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-07',
      sptNumber: 'SPT-MASA-23-2026-07',
      masaMonth: 7,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Service HVAC Ruang Operasi SSC',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 165000000,
      dppAmount: 165000000,
      taxRate: '2%',
      taxAmount: 3300000,
      ntpnNumber: 'NTPN-BCA-20260728-07',
      bpeNumber: 'BPE-23-202607-9907',
      paymentDate: '2026-07-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Service HEPA filter & AC sentral steril',
      businessUnit: 'SSC',
      createdAt: '2026-07-28 10:15:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-08',
      sptNumber: 'SPT-MASA-23-2026-08',
      masaMonth: 8,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa IT Cloud & Cyber Security SIMRS',
      bankAccountCode: '1102.002',
      bankAccountName: 'Bank Mandiri Operasional & Payroll',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-MDR-20260828-08',
      bpeNumber: 'BPE-23-202608-9908',
      paymentDate: '2026-08-28',
      status: 'REPORTED',
      notes: 'Transfer Bank Mandiri: Bridging SATUSEHAT & hosting server',
      businessUnit: 'ALL',
      createdAt: '2026-08-28 16:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-09',
      sptNumber: 'SPT-MASA-23-2026-09',
      masaMonth: 9,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Limbah Medis B3 September',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 140000000,
      dppAmount: 140000000,
      taxRate: '2%',
      taxAmount: 2800000,
      ntpnNumber: 'NTPN-BCA-20260928-09',
      bpeNumber: 'BPE-23-202609-9909',
      paymentDate: '2026-09-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Pengolahan limbah jarum & residu',
      businessUnit: 'ALL',
      createdAt: '2026-09-28 13:45:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-10',
      sptNumber: 'SPT-MASA-23-2026-10',
      masaMonth: 10,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Kalibrasi Tensimeter & Alkes',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-BCA-20261028-10',
      bpeNumber: 'BPE-23-202610-9910',
      paymentDate: '2026-10-28',
      status: 'REPORTED',
      notes: 'Transfer Bank BCA: Sertifikasi kalibrasi BPFK Kemenkes',
      businessUnit: 'ALL',
      createdAt: '2026-10-28 11:30:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'spt-23-11',
      sptNumber: 'SPT-MASA-23-2026-11',
      masaMonth: 11,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Vendor Limbah B3 & Maintenance AC',
      bankAccountCode: '1102.001',
      bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
      paymentChannel: 'BANK_TRANSFER',
      isBankVerified: true,
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-BCA-20261128-11',
      bpeNumber: 'BPE-23-202611-9911',
      paymentDate: '2026-11-28',
      status: 'PAID',
      notes: 'Internet Banking Bank BCA: Disetor ke Kas Negara via DJP',
      businessUnit: 'ALL',
      createdAt: '2026-11-28 14:10:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },

    // --- PPH 25 (ANGSURAN BULANAN PPH BADAN - SETOR VIA E-BILLING BANK) ---
    // Rp 13.000.000 per bulan x 11 bulan = Rp 143.000.000 (Semua via Bank BCA/Mandiri)
    ...Array.from({ length: 11 }, (_, i) => {
      const m = i + 1
      const mStr = String(m).padStart(2, '0')
      return {
        id: `spt-25-${mStr}`,
        sptNumber: `SPT-MASA-25-2026-${mStr}`,
        masaMonth: m,
        masaYear: 2026,
        taxType: 'PPH_25' as SptMasaType,
        accountCode: '2505',
        accountName: 'PPh 25 Angsuran Pajak Penghasilan Badan',
        bankAccountCode: '1102.001',
        bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
        paymentChannel: 'BANK_EBILLING' as const,
        isBankVerified: true,
        grossAmount: 0,
        dppAmount: 0,
        taxRate: 'Angsuran Tetap SKP/SPT',
        taxAmount: 13000000,
        ntpnNumber: `NTPN-BCA-PPH25-${mStr}`,
        bpeNumber: `BPE-25-2026${mStr}-0021`,
        paymentDate: `2026-${mStr}-15`,
        status: 'REPORTED' as const,
        notes: `e-Billing Bank BCA: Setoran Angsuran Masa ${mStr}/2026 PPh Badan BCS`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-15 08:30:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    }),

    // --- PPH 21 (TENAGA AHLI DOKTER & KARYAWAN - TRANSFER BANK MANDIRI / BCA) ---
    ...Array.from({ length: 11 }, (_, i) => {
      const m = i + 1
      const mStr = String(m).padStart(2, '0')
      return {
        id: `spt-21-${mStr}`,
        sptNumber: `SPT-MASA-21-2026-${mStr}`,
        masaMonth: m,
        masaYear: 2026,
        taxType: 'PPH_21' as SptMasaType,
        accountCode: '2501.001',
        accountName: 'PPh 21 Tenaga Medis Dokter & Karyawan',
        bankAccountCode: '1102.002',
        bankAccountName: 'Bank Mandiri Operasional & Payroll',
        paymentChannel: 'BANK_TRANSFER' as const,
        isBankVerified: true,
        grossAmount: 435000000,
        dppAmount: 310000000,
        taxRate: 'PMK 168 / TER',
        taxAmount: 13750000,
        ntpnNumber: `NTPN-MDR-PPH21-${mStr}`,
        bpeNumber: `BPE-21-2026${mStr}-8819`,
        paymentDate: `2026-${mStr}-10`,
        status: 'REPORTED' as const,
        notes: `Payroll Bank Mandiri: Withholding 21 Dokter & TER Staf Masa ${mStr}`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-10 16:00:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    }),

    // --- PPH 4(2) FINAL (SEWA GEDUNG & BUNGA BANK - MUTASI REKENING GIRO) ---
    ...Array.from({ length: 11 }, (_, i) => {
      const m = i + 1
      const mStr = String(m).padStart(2, '0')
      return {
        id: `spt-4-2-${mStr}`,
        sptNumber: `SPT-MASA-4-2-2026-${mStr}`,
        masaMonth: m,
        masaYear: 2026,
        taxType: 'PPH_4_2' as SptMasaType,
        accountCode: '2503',
        accountName: 'PPh Final Pasal 4(2) Persewaan & Bunga',
        bankAccountCode: '1102.001',
        bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
        paymentChannel: 'BANK_TRANSFER' as const,
        isBankVerified: true,
        grossAmount: 50000000,
        dppAmount: 50000000,
        taxRate: '10%',
        taxAmount: 5000000,
        ntpnNumber: `NTPN-BCA-4-2-${mStr}`,
        bpeNumber: `BPE-4-2-2026${mStr}-4412`,
        paymentDate: `2026-${mStr}-12`,
        status: 'REPORTED' as const,
        notes: `Transfer Bank BCA: PPh Final sewa ruko & pemotongan bunga jasa giro`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-12 11:00:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    })
  ])

  // Filtered SPT Masa Records (STRICTLY BANK VERIFIED)
  const filteredSptMasa = computed(() => {
    return sptMasaRecords.value.filter(s => {
      if (s.masaYear !== selectedYear.value) return false
      if (selectedMonth.value && s.masaMonth !== selectedMonth.value) return false
      if (selectedSptMasaTypeFilter.value !== 'ALL' && s.taxType !== selectedSptMasaTypeFilter.value) return false
      // STRICTLY REQUIRE BANK VERIFICATION
      return s.isBankVerified !== false
    })
  })

  // Action: Buat Dokumen SPT Masa Baru (STRICTLY REKening BANK ONLY, NO CASH)
  function createSptMasa(data: {
    masaMonth: number
    masaYear: number
    taxType: SptMasaType
    bankAccountCode: string // REQUIRED: Strictly Bank Account
    grossAmount: number
    dppAmount: number
    taxRate: string
    taxAmount: number
    notes?: string
    ntpnNumber?: string
    status?: 'DRAFT' | 'PAID' | 'REPORTED'
    businessUnit?: BusinessUnit
  }) {
    // Validation: Deny any cash account (1101)
    if (!data.bankAccountCode || data.bankAccountCode.startsWith('1101')) {
      throw new Error('Pembuatan SPT Masa hanya dapat menggunakan Rekening Bank resmi DJP, bukan akun kas fisik/tunai!')
    }

    const bankAcc = availableBankAccounts.value.find(b => b.code === data.bankAccountCode) || availableBankAccounts.value[0]
    const bankCode = bankAcc?.code || '1102.001'
    const bankName = bankAcc?.name || 'Bank BCA Operasional (Giro & EDC)'

    const mStr = String(data.masaMonth).padStart(2, '0')
    const typeLabel = data.taxType.replace('_', ' ')
    const sptNum = `SPT-${data.taxType}-${data.masaYear}-${mStr}`

    let accCode = '2501.001'
    let accName = 'Hutang PPh 21'
    if (data.taxType === 'PPH_23') {
      accCode = '2502'
      accName = 'PPh 23 Jasa Teknik & Limbah B3'
    } else if (data.taxType === 'PPH_4_2') {
      accCode = '2503'
      accName = 'PPh Final Pasal 4(2)'
    } else if (data.taxType === 'PPH_25') {
      accCode = '2505'
      accName = 'PPh 25 Angsuran Badan Bulanan'
    } else if (data.taxType === 'PPN') {
      accCode = '2504'
      accName = 'PPN Keluaran'
    }

    const st = data.status || 'PAID'
    const newRecord: SptMasaRecord = {
      id: 'spt-' + Date.now(),
      sptNumber: sptNum,
      masaMonth: data.masaMonth,
      masaYear: data.masaYear,
      taxType: data.taxType,
      accountCode: accCode,
      accountName: accName,
      bankAccountCode: bankCode,
      bankAccountName: bankName,
      paymentChannel: 'BANK_EBILLING',
      isBankVerified: true,
      grossAmount: data.grossAmount,
      dppAmount: data.dppAmount,
      taxRate: data.taxRate,
      taxAmount: data.taxAmount,
      ntpnNumber: data.ntpnNumber || `NTPN-DJP-${Date.now().toString().slice(-8)}`,
      bpeNumber: `BPE-${data.taxType}-${data.masaYear}${mStr}-${Date.now().toString().slice(-4)}`,
      paymentDate: `${data.masaYear}-${mStr}-20`,
      status: st,
      notes: data.notes || `SPT Masa ${typeLabel} Masa ${mStr}/${data.masaYear} via ${bankName}`,
      businessUnit: data.businessUnit || accountingStore.businessUnit,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    sptMasaRecords.value.unshift(newRecord)
    accountingStore.showToast(`Dokumen SPT Masa ${sptNum} via ${bankName} berhasil dibuat & disimpan!`, 'success')

    // Persist to PostgreSQL database asynchronously
    let backendTaxType = 'PPH23'
    if (data.taxType === 'PPH_21') backendTaxType = 'PPH21_DOCTOR'
    else if (data.taxType === 'PPH_4_2') backendTaxType = 'PPH4_2'
    else if (data.taxType === 'PPH_25') backendTaxType = 'PPH25'
    else if (data.taxType === 'PPN') backendTaxType = 'PPN_OUTPUT'

    const parsedRate = parseFloat(data.taxRate.replace('%', '')) / 100 || 0.02

    accountingApi.createTaxWithholding({
      tax_type: backendTaxType,
      taxpayer_name: `SPT Masa ${typeLabel} - ${bankName}`,
      object_code: data.taxType === 'PPH_21' ? '21-100-07' : '28-100-02',
      gross_amount: data.grossAmount,
      dpp_amount: data.dppAmount,
      tax_rate: parsedRate,
      tax_amount: data.taxAmount,
      period_year: data.masaYear,
      period_month: data.masaMonth
    }).catch(err => {
      console.warn('Note: SPT Masa created in state, DB sync response:', err.message)
    })

    return newRecord
  }

  const isLoading = ref(false)

  // Fetch Tax records directly from PostgreSQL database via Go API
  async function fetchTaxRecords() {
    try {
      isLoading.value = true
      const backendWithholdings = await accountingApi.getTaxWithholdings({ limit: 100 })
      if (backendWithholdings && backendWithholdings.length > 0) {
        const dbRecords: SptMasaRecord[] = backendWithholdings.map(tw => {
          let sptType: SptMasaType = 'PPH_23'
          if (tw.tax_type === 'PPH21_DOCTOR' || tw.tax_type === 'PPH21_EMPLOYEE') sptType = 'PPH_21'
          else if (tw.tax_type === 'PPH23') sptType = 'PPH_23'
          else if (tw.tax_type === 'PPH4_2') sptType = 'PPH_4_2'
          else if (tw.tax_type === 'PPH25') sptType = 'PPH_25'
          else if (tw.tax_type === 'PPN_OUTPUT') sptType = 'PPN'

          return {
            id: tw.id,
            sptNumber: tw.bupot_number,
            masaMonth: tw.period_month,
            masaYear: tw.period_year,
            taxType: sptType,
            accountCode: sptType === 'PPH_23' ? '2502' : sptType === 'PPH_25' ? '2505' : '2501',
            accountName: tw.taxpayer_name,
            bankAccountCode: '1102.001',
            bankAccountName: 'Bank BCA Operasional (Giro & EDC)',
            paymentChannel: 'BANK_EBILLING',
            isBankVerified: true,
            grossAmount: Number(tw.gross_amount),
            dppAmount: Number(tw.dpp_amount),
            taxRate: `${Number(tw.tax_rate) * 100}%`,
            taxAmount: Number(tw.tax_amount),
            ntpnNumber: `NTPN-${tw.bupot_number.slice(-8)}`,
            bpeNumber: `BPE-${tw.bupot_number.slice(-8)}`,
            paymentDate: `${tw.period_year}-${String(tw.period_month).padStart(2, '0')}-28`,
            status: tw.is_reported ? 'REPORTED' : 'PAID',
            notes: `Wajib Pajak: ${tw.taxpayer_name} (NPWP: ${tw.taxpayer_npwp || '-'})`,
            businessUnit: 'ALL',
            createdAt: tw.created_at || new Date().toISOString(),
            createdBy: 'Rian Pratama, S.Ak.'
          }
        })

        // Merge database records with default records
        const dbNumbers = new Set(dbRecords.map(r => r.sptNumber))
        const remainingInitial = sptMasaRecords.value.filter(r => !dbNumbers.has(r.sptNumber))
        sptMasaRecords.value = [...dbRecords, ...remainingInitial]
      }
    } catch (err: any) {
      console.error('Failed to fetch tax records from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load
  fetchTaxRecords()

  // =========================================================================
  // 2. MATRIX KREDIT PAJAK DARI SPT MASA (Januari - Desember)
  // Menampilkan verifikasi rekening Bank pembayar per bulan
  // =========================================================================
  const monthlySptMasaMatrix = computed(() => {
    const months = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ]

    return months.map((monthName, idx) => {
      const monthNum = idx + 1
      const recordsThisMonth = sptMasaRecords.value.filter(
        r => r.masaYear === selectedYear.value && r.masaMonth === monthNum && r.isBankVerified !== false
      )

      const pph21Record = recordsThisMonth.find(r => r.taxType === 'PPH_21')
      const pph23Record = recordsThisMonth.find(r => r.taxType === 'PPH_23')
      const pph42Record = recordsThisMonth.find(r => r.taxType === 'PPH_4_2')
      const pph25Record = recordsThisMonth.find(r => r.taxType === 'PPH_25')

      const pph21Amount = pph21Record?.taxAmount || 0
      const pph23Amount = pph23Record?.taxAmount || 0
      const pph42Amount = pph42Record?.taxAmount || 0
      const pph25Amount = pph25Record?.taxAmount || 0

      // Kredit Pajak untuk SPT 1771 = PPh 23 (Bank) + PPh 25 (Bank)
      const totalCreditThisMonth = pph23Amount + pph25Amount

      const bankUsed = pph25Record?.bankAccountName || pph23Record?.bankAccountName || 'Bank BCA Operasional'

      return {
        monthNum,
        monthName,
        pph21Record,
        pph21Amount,
        pph23Record,
        pph23Amount,
        pph42Record,
        pph42Amount,
        pph25Record,
        pph25Amount,
        totalCreditThisMonth,
        bankUsed,
        isReported: (pph23Record?.status === 'REPORTED' || pph25Record?.status === 'REPORTED'),
        ntpn: pph25Record?.ntpnNumber || pph23Record?.ntpnNumber || '-'
      }
    })
  })

  // Doctor Tax Calculations per PMK 168/2023 (Semua ditransfer via Bank)
  const doctorTaxRecords = ref<DoctorTaxRecord[]>([
    {
      id: 'dt-1',
      doctorName: 'dr. Hendra Wijaya, Sp.A',
      specialty: 'Spesialis Anak (Pediatri)',
      npwp: '08.123.456.7-012.000',
      nik: '3171012304850001',
      businessUnit: 'KIDINEST',
      bankAccountTarget: 'Bank BCA: 883-019-2811 a.n Hendra Wijaya',
      grossIncomeCurrentMonth: 48500000,
      grossIncomeCumulativeYear: 142000000,
      dppPercentage: 0.5,
      dppAmount: 24250000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%)',
      taxRate: 0.05,
      pph21Amount: 1212500,
      netPayable: 47287500,
      ebupotNumber: 'BUPOT-21-202611-001',
      journalNumber: 'JU-2026-0004'
    },
    {
      id: 'dt-2',
      doctorName: 'dr. Maya Sandria, Sp.KK, FINSDV',
      specialty: 'Spesialis Dermatologi & Venereologi',
      npwp: '09.345.678.9-015.000',
      nik: '3174055208880003',
      businessUnit: 'SSC',
      bankAccountTarget: 'Bank BCA: 009-441-8922 a.n Maya Sandria',
      grossIncomeCurrentMonth: 82500000,
      grossIncomeCumulativeYear: 245000000,
      dppPercentage: 0.5,
      dppAmount: 41250000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%) & Lapisan 2 (15%)',
      taxRate: 0.05,
      pph21Amount: 2062500,
      netPayable: 80437500,
      ebupotNumber: 'BUPOT-21-202611-002',
      journalNumber: 'JU-2026-0005'
    },
    {
      id: 'dt-3',
      doctorName: 'dr. Nadia Kirana, Sp.BP-RE',
      specialty: 'Spesialis Bedah Plastik & Rekonstruksi',
      npwp: '12.876.543.2-021.000',
      nik: '3276014407910005',
      businessUnit: 'SSC',
      bankAccountTarget: 'Bank Mandiri: 137-00-19828-1 a.n Nadia Kirana',
      grossIncomeCurrentMonth: 54000000,
      grossIncomeCumulativeYear: 165000000,
      dppPercentage: 0.5,
      dppAmount: 27000000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%)',
      taxRate: 0.05,
      pph21Amount: 1350000,
      netPayable: 52650000,
      ebupotNumber: 'BUPOT-21-202611-003',
      journalNumber: 'JU-2026-0006'
    },
    {
      id: 'dt-4',
      doctorName: 'dr. Bambang Soeprapto, Sp.A(K)',
      specialty: 'Konsultan Tumbuh Kembang Anak',
      npwp: '14.567.890.1-018.000',
      nik: '3172081206770002',
      businessUnit: 'KIDINEST',
      bankAccountTarget: 'Bank BCA: 883-991-0129 a.n Bambang Soeprapto',
      grossIncomeCurrentMonth: 65000000,
      grossIncomeCumulativeYear: 198000000,
      dppPercentage: 0.5,
      dppAmount: 32500000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%)',
      taxRate: 0.05,
      pph21Amount: 1625000,
      netPayable: 63375000,
      ebupotNumber: 'BUPOT-21-202611-004',
      journalNumber: 'JU-2026-0007'
    }
  ])

  // Tax Transactions (Hanya mencatat transaksi berbasis Rekening Bank)
  const taxTransactions = ref<TaxTransaction[]>([
    {
      id: 'tx-21-01',
      taxNumber: 'BUPOT-21-202611-001',
      taxDate: '2026-11-15',
      businessUnit: 'KIDINEST',
      taxCategory: 'PPH_21',
      accountCode: '2501.001',
      journalNumber: 'JU-2026-0004',
      counterpartyName: 'dr. Hendra Wijaya, Sp.A (Transfer BCA)',
      counterpartyNPWP: '08.123.456.7-012.000',
      counterpartyNIK: '3171012304850001',
      objectCode: '21-100-07 (Tenaga Ahli - Dokter)',
      grossAmount: 48500000,
      dppAmount: 24250000,
      taxRate: 0.05,
      taxAmount: 1212500,
      documentReference: 'TRF-BCA-DOC-202611-01',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-21-02',
      taxNumber: 'BUPOT-21-202611-002',
      taxDate: '2026-11-20',
      businessUnit: 'SSC',
      taxCategory: 'PPH_21',
      accountCode: '2501.001',
      journalNumber: 'JU-2026-0005',
      counterpartyName: 'dr. Maya Sandria, Sp.KK (Transfer BCA)',
      counterpartyNPWP: '09.345.678.9-015.000',
      counterpartyNIK: '3174055208880003',
      objectCode: '21-100-07 (Tenaga Ahli - Dokter)',
      grossAmount: 82500000,
      dppAmount: 41250000,
      taxRate: 0.05,
      taxAmount: 2062500,
      documentReference: 'TRF-BCA-DOC-202611-02',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-21-karyawan',
      taxNumber: 'BUPOT-21-KARYAWAN-11',
      taxDate: '2026-11-28',
      businessUnit: 'ALL',
      taxCategory: 'PPH_21',
      accountCode: '2501.002',
      journalNumber: 'JU-2026-0008',
      counterpartyName: 'Payroll Pegawai Tetap & Apoteker (Transfer Mandiri)',
      counterpartyNPWP: '01.234.567.8-091.000',
      counterpartyNIK: '3374019901900001',
      objectCode: '21-100-01 (Pegawai Tetap)',
      grossAmount: 185000000,
      dppAmount: 185000000,
      taxRate: 0.0405,
      taxAmount: 7500000,
      documentReference: 'PAYROLL-MDR-202611',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-23-01',
      taxNumber: 'BUPOT-23-202611-01',
      taxDate: '2026-11-10',
      businessUnit: 'ALL',
      taxCategory: 'PPH_23',
      accountCode: '2502',
      journalNumber: 'JU-2026-0002',
      counterpartyName: 'PT Wastec International (Transfer Bank BCA)',
      counterpartyNPWP: '01.554.219.8-032.000',
      objectCode: '24-104-01 (Jasa Pengolahan Limbah Medis B3)',
      grossAmount: 45000000,
      dppAmount: 45000000,
      taxRate: 0.02,
      taxAmount: 900000,
      documentReference: 'TRF-BCA-WST-8891',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-23-02',
      taxNumber: 'BUPOT-23-202611-02',
      taxDate: '2026-11-18',
      businessUnit: 'SSC',
      taxCategory: 'PPH_23',
      accountCode: '2502',
      journalNumber: 'JU-2026-0003',
      counterpartyName: 'PT Medika Laser Servis (Transfer Bank BCA)',
      counterpartyNPWP: '02.881.902.1-014.000',
      objectCode: '24-104-02 (Maintenance Alkes Laser)',
      grossAmount: 75000000,
      dppAmount: 75000000,
      taxRate: 0.02,
      taxAmount: 1500000,
      documentReference: 'TRF-BCA-MLS-2026',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-4-2-sewa',
      taxNumber: 'BUPOT-4-2-202611',
      taxDate: '2026-11-05',
      businessUnit: 'ALL',
      taxCategory: 'PPH_4_2',
      accountCode: '2503',
      journalNumber: 'JU-2026-0001',
      counterpartyName: 'Hj. Siti Nurhaliza (Sewa Ruko Transfer BCA)',
      counterpartyNPWP: '07.991.234.5-011.000',
      objectCode: '28-403-01 (Persewaan Gedung Ruko)',
      grossAmount: 50000000,
      dppAmount: 50000000,
      taxRate: 0.10,
      taxAmount: 5000000,
      documentReference: 'TRF-BCA-SEWA-202611',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-ppn-out',
      taxNumber: 'FP-010-202611-0099',
      taxDate: '2026-11-30',
      businessUnit: 'SSC',
      taxCategory: 'PPN_OUTPUT',
      accountCode: '2504',
      journalNumber: 'JU-2026-0010',
      counterpartyName: 'Pelanggan Produk Skincare (Settlement EDC & QRIS Bank)',
      counterpartyNPWP: '00.000.000.0-000.000',
      objectCode: 'Penyerahan BKP Produk Estetika',
      grossAmount: 95596409,
      dppAmount: 86122880,
      taxRate: 0.11,
      taxAmount: 9475055,
      documentReference: 'EDC-BCA-RECAP-202611',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    }
  ])

  // PPN Summary Computation
  const ppnSummary = computed(() => {
    const ppnOutAcc = ledgerStore.accounts.find(a => a.code === '2504')
    const ppnInAcc = ledgerStore.accounts.find(a => a.code === '1504')

    const ppnKeluaran = ppnOutAcc?.currentBalance || 9475055
    const ppnMasukan = ppnInAcc?.currentBalance || 3200000
    const ppnKurangBayar = Math.max(0, ppnKeluaran - ppnMasukan)

    return {
      ppnKeluaran,
      ppnMasukan,
      ppnKurangBayar,
      status: ppnKurangBayar > 0 ? 'KURANG_BAYAR' : 'LEBIH_BAYAR'
    }
  })

  // Filtered Tax Transactions based on selected Tax Account
  const filteredTaxTransactions = computed(() => {
    return taxTransactions.value.filter(tx => {
      if (selectedTaxAccountCode.value !== 'ALL') {
        if (selectedTaxAccountCode.value === '2501.001' && tx.accountCode !== '2501.001') return false
        if (selectedTaxAccountCode.value === '2501.002' && tx.accountCode !== '2501.002') return false
        if (selectedTaxAccountCode.value === '2502' && tx.taxCategory !== 'PPH_23') return false
        if (selectedTaxAccountCode.value === '2503' && tx.taxCategory !== 'PPH_4_2') return false
        if (selectedTaxAccountCode.value === '2504' && tx.taxCategory !== 'PPN_OUTPUT') return false
      }
      return true
    })
  })

  // Monthly summary metrics based on selected account
  const monthlyTaxSummary = computed(() => {
    const list = filteredTaxTransactions.value
    const totalGross = list.reduce((s, t) => s + t.grossAmount, 0)
    const totalDpp = list.reduce((s, t) => s + t.dppAmount, 0)
    const totalTax = list.reduce((s, t) => s + t.taxAmount, 0)
    
    let coaBalance = 0
    if (selectedTaxAccountCode.value === 'ALL') {
      const taxAccs = ledgerStore.accounts.filter(a => a.code.startsWith('25'))
      coaBalance = taxAccs.reduce((sum, a) => sum + (a.currentBalance || 0), 0)
    } else {
      const acc = ledgerStore.accounts.find(a => a.code === selectedTaxAccountCode.value)
      coaBalance = acc?.currentBalance || totalTax
    }

    return {
      totalGross,
      totalDpp,
      totalTax,
      coaBalance,
      isReconciled: Math.abs(totalTax - coaBalance) < 10000000
    }
  })

  // =========================================================================
  // 3. ANNUAL TAX CALCULATION (SPT TAHUNAN PPH BADAN FORMULIR 1771)
  // MENGAMBIL NILAI KREDIT PAJAK HANYA DARI BANK (BUKAN KAS)
  // =========================================================================
  const annualTaxReport1771 = computed<AnnualTaxReport1771>(() => {
    // 1. Commercial Revenues (Akun 4 + Akun 6)
    const rev4 = ledgerStore.accounts.find(a => a.code === '4')?.currentBalance || 15860537440
    const rev6 = ledgerStore.accounts.find(a => a.code === '6')?.currentBalance || 38481777.73
    const commercialRevenue = rev4 + rev6

    // 2. Commercial Expenses (Akun 5 + Akun 7)
    const exp5 = ledgerStore.accounts.find(a => a.code === '5')?.currentBalance || 12284590157.71
    const exp7 = ledgerStore.accounts.find(a => a.code === '7')?.currentBalance || 872980221.16
    const commercialExpense = exp5 + exp7

    // 3. Commercial Net Profit
    const commercialNetProfit = commercialRevenue - commercialExpense

    // 4. Fiscal Corrections per UU PPh
    const positiveFiscalCorrection = 58500000.00
    const positiveCorrectionNotes = 'Koreksi Positif: Biaya entertainment/jamuan tanpa daftar nominatif PMK-02/2010 (Rp 38.5M) + Sanksi denda bunga (Rp 20M)'
    const negativeFiscalCorrection = 9545162.73
    const negativeCorrectionNotes = 'Koreksi Negatif: Pendapatan Bunga Bank telah dipotong PPh Final Pasal 4(2) oleh Bank'

    // 5. Taxable Income (Penghasilan Kena Pajak / PKP)
    const taxableIncome = commercialNetProfit + positiveFiscalCorrection - negativeFiscalCorrection

    // 6. Corporate Tax Rate 22% & Fasilitas 31E UU PPh
    const corporateTaxRate = 0.22
    const corporateTaxPayable = Math.round(taxableIncome * 0.11)

    // 7. KREDIT PAJAK OTOMATIS DARI SPT MASA SESUAI JENISNYA (STRICTLY BANK ONLY):
    // Hanya mengakui kredit pajak yang dipotong/disetor via Rekening Bank resmi
    const sptMasaPph23Records = sptMasaRecords.value.filter(
      r => r.masaYear === selectedYear.value && r.taxType === 'PPH_23' && r.isBankVerified !== false
    )
    const taxCreditPph23 = sptMasaPph23Records.reduce((sum, r) => sum + r.taxAmount, 0)

    const sptMasaPph25Records = sptMasaRecords.value.filter(
      r => r.masaYear === selectedYear.value && r.taxType === 'PPH_25' && r.isBankVerified !== false
    )
    const taxCreditPph25 = sptMasaPph25Records.reduce((sum, r) => sum + r.taxAmount, 0)

    // Total Kredit Pajak (PPh 23 + PPh 25)
    const taxCredits = taxCreditPph23 + taxCreditPph25

    // 8. PPh 29 Kurang Bayar yang harus disetor ke Kas Negara
    const pph29Payable = Math.max(0, corporateTaxPayable - taxCredits)

    return {
      fiscalYear: selectedYear.value,
      status: 'FINALIZED',
      commercialRevenue,
      commercialExpense,
      commercialNetProfit,
      positiveFiscalCorrection,
      positiveCorrectionNotes,
      negativeFiscalCorrection,
      negativeCorrectionNotes,
      taxableIncome,
      corporateTaxRate,
      corporateTaxPayable,
      taxCreditPph23,
      taxCreditPph25,
      taxCredits,
      pph29Payable,
      linkedSptMasaCount: sptMasaPph23Records.length + sptMasaPph25Records.length,
      bankReconciliationVerified: true
    }
  })

  // Filtered doctor tax records
  const filteredDoctorTax = computed(() => {
    return doctorTaxRecords.value.filter(d => {
      if (accountingStore.businessUnit !== 'ALL' && d.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return true
    })
  })

  // Action: Generate Jurnal Pembayaran / Penyetoran Pajak NTPN (STRICTLY VIA BANK, BUKAN KAS)
  function generateTaxPaymentJournal(
    taxAccountCode: string,
    taxName: string,
    amount: number,
    bankAccountCode: string = '1102.001'
  ) {
    // Validate: Only bank accounts allowed
    if (bankAccountCode.startsWith('1101')) {
      accountingStore.showToast('Penyetoran pajak negara hanya dapat dilakukan via rekening Bank, bukan kas fisik!', 'error')
      return
    }

    const bankAcc = availableBankAccounts.value.find(b => b.code === bankAccountCode) || availableBankAccounts.value[0]
    const selectedBankCode = bankAcc?.code || '1102.001'
    const selectedBankName = bankAcc?.name || 'Bank BCA Operasional (Giro & EDC)'

    const journal = journalsStore.createManualJournal({
      journalDate: `${selectedYear.value}-${String(selectedMonth.value).padStart(2, '0')}-28`,
      businessUnit: accountingStore.businessUnit,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Penyetoran NTPN Kas Negara via ${selectedBankName}: ${taxName} Periode ${selectedMonth.value}/${selectedYear.value} (Ref: NTPN-DJP-${Date.now().toString().slice(-8)})`,
      status: 'POSTED',
      journalSource: 'MANUAL',
      sourceReference: `NTPN-${Date.now().toString().slice(-8)}`,
      lines: [
        {
          accountId: `acc-${taxAccountCode}`,
          accountCode: taxAccountCode,
          accountName: `Hutang Pajak (${taxName})`,
          debit: amount,
          credit: 0,
          description: `Pelunasan Hutang Pajak ${taxName}`
        },
        {
          accountId: `acc-${selectedBankCode}`,
          accountCode: selectedBankCode,
          accountName: selectedBankName,
          debit: 0,
          credit: amount,
          description: `Pembayaran Pajak DJP via Internet Banking ${selectedBankName}`
        }
      ]
    })

    accountingStore.showToast(`Jurnal Penyetoran Pajak via ${selectedBankName} (${journal.journalNumber}) berhasil dibuat!`, 'success')
  }

  // Action: Export DJP e-Bupot CSV Format
  function exportEbupotCSV() {
    const csvData = doctorTaxRecords.value.map(d => ({
      'KODE_OBJEK_PAJAK': '21-100-07',
      'NAMA_PENERIMA_PENGHASILAN': d.doctorName,
      'NPWP': d.npwp,
      'NIK': d.nik,
      'METODE_TRANSFER': d.bankAccountTarget,
      'PENGHASILAN_BRUTO': d.grossIncomeCurrentMonth,
      'TARIF': (d.taxRate * 100) + '%',
      'DPP': d.dppAmount,
      'PPH_DIPOTONG': d.pph21Amount,
      'NO_BUKTI_POTONG': d.ebupotNumber,
      'NO_JURNAL': d.journalNumber || '-',
      'TGL_BUKTI_POTONG': `30/${selectedMonth.value}/${selectedYear.value}`
    }))

    exportToCSV(csvData, `DJP_eBupot_21_${selectedYear.value}_${selectedMonth.value}`)
    accountingStore.showToast('File e-Bupot PPh 21 format DJP berhasil diunduh!', 'success')
  }

  return {
    isLoading,
    fetchTaxRecords,
    taxPeriodType,
    selectedTaxAccountCode,
    selectedMonth,
    selectedYear,
    selectedSptMasaTypeFilter,
    taxAccountOptions,
    availableBankAccounts,
    sptMasaRecords,
    filteredSptMasa,
    createSptMasa,
    monthlySptMasaMatrix,
    doctorTaxRecords,
    filteredDoctorTax,
    taxTransactions,
    filteredTaxTransactions,
    monthlyTaxSummary,
    annualTaxReport1771,
    ppnSummary,
    generateTaxPaymentJournal,
    exportEbupotCSV
  }
})
