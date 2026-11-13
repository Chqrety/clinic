import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TaxTransaction, BusinessUnit } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'
import { exportToCSV } from '@/utils/export'

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
  taxCreditPph23: number // Otomatis dari akumulasi SPT Masa PPh 23
  taxCreditPph25: number // Otomatis dari akumulasi SPT Masa PPh 25
  taxCredits: number // Total Kredit Pajak (PPh 23 + PPh 25)
  pph29Payable: number // PPh Kurang Bayar Pasal 29
  linkedSptMasaCount: number
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

  // =========================================================================
  // 1. REKAPITULASI DOKUMEN SPT MASA (Januari - November 2026)
  // Dipisah pembuatannya dari SPT Tahunan!
  // =========================================================================
  const sptMasaRecords = ref<SptMasaRecord[]>([
    // --- PPH 23 (KREDIT PAJAK LAMPIRAN III SPT 1771) ---
    {
      id: 'spt-23-01',
      sptNumber: 'SPT-MASA-23-2026-01',
      masaMonth: 1,
      masaYear: 2026,
      taxType: 'PPH_23',
      accountCode: '2502',
      accountName: 'PPh 23 Jasa Pengolahan Limbah Medis & IT',
      grossAmount: 120000000,
      dppAmount: 120000000,
      taxRate: '2%',
      taxAmount: 2400000,
      ntpnNumber: 'NTPN-DJP-20260128-01',
      bpeNumber: 'BPE-23-202601-9901',
      paymentDate: '2026-01-28',
      status: 'REPORTED',
      notes: 'Jasa limbah B3 PT Wastec & maintenance server SIMKlinik',
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
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-DJP-20260228-02',
      bpeNumber: 'BPE-23-202602-9902',
      paymentDate: '2026-02-28',
      status: 'REPORTED',
      notes: 'Kalibrasi laser Nd:YAG & USG 4D Kidinest',
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
      grossAmount: 160000000,
      dppAmount: 160000000,
      taxRate: '2%',
      taxAmount: 3200000,
      ntpnNumber: 'NTPN-DJP-20260328-03',
      bpeNumber: 'BPE-23-202603-9903',
      paymentDate: '2026-03-28',
      status: 'REPORTED',
      notes: 'Outsourcing cleaning service ruang steril & poli anak',
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
      accountName: 'PPh 23 Jasa Konsultasi Hukum & Manajemen Medis',
      grossAmount: 140000000,
      dppAmount: 140000000,
      taxRate: '2%',
      taxAmount: 2800000,
      ntpnNumber: 'NTPN-DJP-20260428-04',
      bpeNumber: 'BPE-23-202604-9904',
      paymentDate: '2026-04-28',
      status: 'REPORTED',
      notes: 'Konsultasi akreditasi klinik Paripurna Kemenkes',
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
      grossAmount: 155000000,
      dppAmount: 155000000,
      taxRate: '2%',
      taxAmount: 3100000,
      ntpnNumber: 'NTPN-DJP-20260528-05',
      bpeNumber: 'BPE-23-202605-9905',
      paymentDate: '2026-05-28',
      status: 'REPORTED',
      notes: 'Maintenance rutin autoklaf & sterilisator bedah',
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
      grossAmount: 145000000,
      dppAmount: 145000000,
      taxRate: '2%',
      taxAmount: 2900000,
      ntpnNumber: 'NTPN-DJP-20260628-06',
      bpeNumber: 'BPE-23-202606-9906',
      paymentDate: '2026-06-28',
      status: 'REPORTED',
      notes: 'Sewa mobil cold chain distribusi vaksin Kidinest',
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
      accountName: 'PPh 23 Jasa Service & Maintenance HVAC Ruang Operasi',
      grossAmount: 165000000,
      dppAmount: 165000000,
      taxRate: '2%',
      taxAmount: 3300000,
      ntpnNumber: 'NTPN-DJP-20260728-07',
      bpeNumber: 'BPE-23-202607-9907',
      paymentDate: '2026-07-28',
      status: 'REPORTED',
      notes: 'Service HEPA filter & AC sentral ruang tindakan medis',
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
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-DJP-20260828-08',
      bpeNumber: 'BPE-23-202608-9908',
      paymentDate: '2026-08-28',
      status: 'REPORTED',
      notes: 'Integrasi bridging SATUSEHAT Kemenkes & BPJS',
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
      accountName: 'PPh 23 Jasa Limbah Medis B3 Periode September',
      grossAmount: 140000000,
      dppAmount: 140000000,
      taxRate: '2%',
      taxAmount: 2800000,
      ntpnNumber: 'NTPN-DJP-20260928-09',
      bpeNumber: 'BPE-23-202609-9909',
      paymentDate: '2026-09-28',
      status: 'REPORTED',
      notes: 'Pengolahan jarum suntik, spuit & residu klinis',
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
      accountName: 'PPh 23 Jasa Kalibrasi Tensimeter & Alkes Diagnostik',
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-DJP-20261028-10',
      bpeNumber: 'BPE-23-202610-9910',
      paymentDate: '2026-10-28',
      status: 'REPORTED',
      notes: 'Sertifikasi kalibrasi tahunan BPFK Kemenkes',
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
      grossAmount: 150000000,
      dppAmount: 150000000,
      taxRate: '2%',
      taxAmount: 3000000,
      ntpnNumber: 'NTPN-DJP-20261128-11',
      bpeNumber: 'BPE-23-202611-9911',
      paymentDate: '2026-11-28',
      status: 'PAID',
      notes: 'Disetor ke Kas Negara via BCA Internet Banking',
      businessUnit: 'ALL',
      createdAt: '2026-11-28 14:10:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },

    // --- PPH 25 (ANGSURAN BULANAN PPH BADAN - INDUK 10A SPT 1771) ---
    // Rp 13.000.000 per bulan x 11 bulan = Rp 143.000.000
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
        grossAmount: 0,
        dppAmount: 0,
        taxRate: 'Angsuran Tetap SKP/SPT',
        taxAmount: 13000000,
        ntpnNumber: `NTPN-PPH25-2026${mStr}15`,
        bpeNumber: `BPE-25-2026${mStr}-0021`,
        paymentDate: `2026-${mStr}-15`,
        status: 'REPORTED' as const,
        notes: `Setoran Angsuran Masa ${mStr}/2026 PPh Badan BCS`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-15 08:30:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    }),

    // --- PPH 21 (TENAGA AHLI DOKTER & KARYAWAN) ---
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
        grossAmount: 435000000,
        dppAmount: 310000000,
        taxRate: 'PMK 168 / TER',
        taxAmount: 13750000,
        ntpnNumber: `NTPN-PPH21-2026${mStr}10`,
        bpeNumber: `BPE-21-2026${mStr}-8819`,
        paymentDate: `2026-${mStr}-10`,
        status: 'REPORTED' as const,
        notes: `Withholding e-Bupot 21 Dokter & TER Staf Masa ${mStr}`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-10 16:00:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    }),

    // --- PPH 4(2) FINAL (SEWA GEDUNG & BUNGA BANK) ---
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
        grossAmount: 50000000,
        dppAmount: 50000000,
        taxRate: '10%',
        taxAmount: 5000000,
        ntpnNumber: `NTPN-4-2-2026${mStr}12`,
        bpeNumber: `BPE-4-2-2026${mStr}-4412`,
        paymentDate: `2026-${mStr}-12`,
        status: 'REPORTED' as const,
        notes: `PPh Final sewa ruko lantai 2 & bunga jasa giro`,
        businessUnit: 'ALL' as BusinessUnit,
        createdAt: `2026-${mStr}-12 11:00:00`,
        createdBy: 'Rian Pratama, S.Ak.'
      }
    })
  ])

  // Filtered SPT Masa Records
  const filteredSptMasa = computed(() => {
    return sptMasaRecords.value.filter(s => {
      if (s.masaYear !== selectedYear.value) return false
      if (selectedMonth.value && s.masaMonth !== selectedMonth.value) return false
      if (selectedSptMasaTypeFilter.value !== 'ALL' && s.taxType !== selectedSptMasaTypeFilter.value) return false
      return true
    })
  })

  // Action: Buat Dokumen SPT Masa Baru
  function createSptMasa(data: {
    masaMonth: number
    masaYear: number
    taxType: SptMasaType
    grossAmount: number
    dppAmount: number
    taxRate: string
    taxAmount: number
    notes?: string
    ntpnNumber?: string
    status?: 'DRAFT' | 'PAID' | 'REPORTED'
    businessUnit?: BusinessUnit
  }) {
    const mStr = String(data.masaMonth).padStart(2, '0')
    const typeLabel = data.taxType.replace('_', ' ')
    const count = sptMasaRecords.value.filter(r => r.taxType === data.taxType && r.masaYear === data.masaYear).length + 1
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
      grossAmount: data.grossAmount,
      dppAmount: data.dppAmount,
      taxRate: data.taxRate,
      taxAmount: data.taxAmount,
      ntpnNumber: data.ntpnNumber || `NTPN-DJP-${Date.now().toString().slice(-8)}`,
      bpeNumber: `BPE-${data.taxType}-${data.masaYear}${mStr}-${Date.now().toString().slice(-4)}`,
      paymentDate: `${data.masaYear}-${mStr}-20`,
      status: st,
      notes: data.notes || `SPT Masa ${typeLabel} Masa ${mStr}/${data.masaYear}`,
      businessUnit: data.businessUnit || accountingStore.businessUnit,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    sptMasaRecords.value.unshift(newRecord)
    accountingStore.showToast(`Dokumen SPT Masa ${sptNum} berhasil dibuat & disimpan!`, 'success')
    return newRecord
  }

  // =========================================================================
  // 2. MATRIX KREDIT PAJAK DARI SPT MASA (Januari - Desember)
  // Menghubungkan secara transparan nilai SPT Masa ke SPT Tahunan 1771
  // =========================================================================
  const monthlySptMasaMatrix = computed(() => {
    const months = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ]

    return months.map((monthName, idx) => {
      const monthNum = idx + 1
      const recordsThisMonth = sptMasaRecords.value.filter(
        r => r.masaYear === selectedYear.value && r.masaMonth === monthNum
      )

      const pph21Record = recordsThisMonth.find(r => r.taxType === 'PPH_21')
      const pph23Record = recordsThisMonth.find(r => r.taxType === 'PPH_23')
      const pph42Record = recordsThisMonth.find(r => r.taxType === 'PPH_4_2')
      const pph25Record = recordsThisMonth.find(r => r.taxType === 'PPH_25')

      const pph21Amount = pph21Record?.taxAmount || 0
      const pph23Amount = pph23Record?.taxAmount || 0
      const pph42Amount = pph42Record?.taxAmount || 0
      const pph25Amount = pph25Record?.taxAmount || 0

      // Kredit Pajak untuk SPT 1771 = PPh 23 + PPh 25
      const totalCreditThisMonth = pph23Amount + pph25Amount

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
        isReported: (pph23Record?.status === 'REPORTED' || pph25Record?.status === 'REPORTED'),
        ntpn: pph25Record?.ntpnNumber || pph23Record?.ntpnNumber || '-'
      }
    })
  })

  // Doctor Tax Calculations per PMK 168/2023
  const doctorTaxRecords = ref<DoctorTaxRecord[]>([
    {
      id: 'dt-1',
      doctorName: 'dr. Hendra Wijaya, Sp.A',
      specialty: 'Spesialis Anak (Pediatri)',
      npwp: '08.123.456.7-012.000',
      nik: '3171012304850001',
      businessUnit: 'KIDINEST',
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

  // Comprehensive Tax Transactions List linked with COA Account Codes & General Ledger
  const taxTransactions = ref<TaxTransaction[]>([
    {
      id: 'tx-21-01',
      taxNumber: 'BUPOT-21-202611-001',
      taxDate: '2026-11-15',
      businessUnit: 'KIDINEST',
      taxCategory: 'PPH_21',
      accountCode: '2501.001',
      journalNumber: 'JU-2026-0004',
      counterpartyName: 'dr. Hendra Wijaya, Sp.A',
      counterpartyNPWP: '08.123.456.7-012.000',
      counterpartyNIK: '3171012304850001',
      objectCode: '21-100-07 (Tenaga Ahli - Dokter)',
      grossAmount: 48500000,
      dppAmount: 24250000,
      taxRate: 0.05,
      taxAmount: 1212500,
      documentReference: 'FEE-DOC-KID-202611-01',
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
      counterpartyName: 'dr. Maya Sandria, Sp.KK, FINSDV',
      counterpartyNPWP: '09.345.678.9-015.000',
      counterpartyNIK: '3174055208880003',
      objectCode: '21-100-07 (Tenaga Ahli - Dokter)',
      grossAmount: 82500000,
      dppAmount: 41250000,
      taxRate: 0.05,
      taxAmount: 2062500,
      documentReference: 'FEE-DOC-SSC-202611-02',
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
      counterpartyName: 'Rekapitulasi 24 Staf, Perawat & Apoteker (TER PP 58)',
      counterpartyNPWP: '01.234.567.8-091.000',
      counterpartyNIK: '3374019901900001',
      objectCode: '21-100-01 (Pegawai Tetap)',
      grossAmount: 185000000,
      dppAmount: 185000000,
      taxRate: 0.0405,
      taxAmount: 7500000,
      documentReference: 'PAYROLL-RECAP-202611',
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
      counterpartyName: 'PT Wastec International (Jasa Limbah Medis B3)',
      counterpartyNPWP: '01.554.219.8-032.000',
      objectCode: '24-104-01 (Jasa Pengolahan Limbah)',
      grossAmount: 45000000,
      dppAmount: 45000000,
      taxRate: 0.02,
      taxAmount: 900000,
      documentReference: 'INV-WST-8891',
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
      counterpartyName: 'PT Medika Laser Servis (Maintenance Alkes Estetika)',
      counterpartyNPWP: '02.881.902.1-014.000',
      objectCode: '24-104-02 (Jasa Perawatan Mesin)',
      grossAmount: 75000000,
      dppAmount: 75000000,
      taxRate: 0.02,
      taxAmount: 1500000,
      documentReference: 'INV-MLS-2026',
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
      counterpartyName: 'Hj. Siti Nurhaliza (Sewa Ruko Klinik Kidinest)',
      counterpartyNPWP: '07.991.234.5-011.000',
      objectCode: '28-403-01 (Persewaan Tanah dan/atau Bangunan)',
      grossAmount: 50000000,
      dppAmount: 50000000,
      taxRate: 0.10,
      taxAmount: 5000000,
      documentReference: 'KWT-SEWA-202611',
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
      counterpartyName: 'Pelanggan Skincare & Kosmetika SSC (Rekap Faktur Digunggung)',
      counterpartyNPWP: '00.000.000.0-000.000',
      objectCode: 'Penyerahan BKP Produk Estetika',
      grossAmount: 95596409,
      dppAmount: 86122880,
      taxRate: 0.11,
      taxAmount: 9475055,
      documentReference: 'FP-RECAP-202611',
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
    
    // Check COA balance for selected account
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
  // MENGAMBIL NILAI KREDIT PAJAK DARI SPT MASA SESUAI JENISNYA!
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
    const corporateTaxPayable = Math.round(taxableIncome * 0.11) // Efektif fasilitas 31E (50% x 22% = 11%)

    // 7. KREDIT PAJAK OTOMATIS DARI SPT MASA SESUAI JENISNYA:
    // a. Kredit Pajak PPh 23 = Sum of SPT Masa PPh 23 untuk tahun berjalan
    const sptMasaPph23Records = sptMasaRecords.value.filter(
      r => r.masaYear === selectedYear.value && r.taxType === 'PPH_23'
    )
    const taxCreditPph23 = sptMasaPph23Records.reduce((sum, r) => sum + r.taxAmount, 0)

    // b. Kredit Pajak PPh 25 = Sum of SPT Masa PPh 25 (Angsuran Badan) untuk tahun berjalan
    const sptMasaPph25Records = sptMasaRecords.value.filter(
      r => r.masaYear === selectedYear.value && r.taxType === 'PPH_25'
    )
    const taxCreditPph25 = sptMasaPph25Records.reduce((sum, r) => sum + r.taxAmount, 0)

    // c. Total Kredit Pajak (PPh 23 + PPh 25)
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
      linkedSptMasaCount: sptMasaPph23Records.length + sptMasaPph25Records.length
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

  // Action: Generate Jurnal Pembayaran / Penyetoran Pajak NTPN
  function generateTaxPaymentJournal(taxAccountCode: string, taxName: string, amount: number) {
    const journal = journalsStore.createManualJournal({
      journalDate: `${selectedYear.value}-${String(selectedMonth.value).padStart(2, '0')}-28`,
      businessUnit: accountingStore.businessUnit,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Penyetoran NTPN Kas Negara: ${taxName} Periode ${selectedMonth.value}/${selectedYear.value} (Ref: NTPN-DJP-${Date.now().toString().slice(-8)})`,
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
          accountId: 'acc-1102.001',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: amount,
          description: `Pembayaran via Internet Banking Pajak DJP`
        }
      ]
    })

    accountingStore.showToast(`Draft Jurnal Pembayaran Pajak ${journal.journalNumber} berhasil dibuat!`, 'success')
  }

  // Action: Export DJP e-Bupot CSV Format
  function exportEbupotCSV() {
    const csvData = doctorTaxRecords.value.map(d => ({
      'KODE_OBJEK_PAJAK': '21-100-07',
      'NAMA_PENERIMA_PENGHASILAN': d.doctorName,
      'NPWP': d.npwp,
      'NIK': d.nik,
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
    taxPeriodType,
    selectedTaxAccountCode,
    selectedMonth,
    selectedYear,
    selectedSptMasaTypeFilter,
    taxAccountOptions,
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
