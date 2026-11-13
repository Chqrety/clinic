import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TaxTransaction, BusinessUnit } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'
import { exportToCSV } from '@/utils/export'

export type TaxPeriodType = 'MASA' | 'TAHUNAN'

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
  taxCredits: number // Kredit Pajak (PPh 23 & PPh 25)
  pph29Payable: number // PPh Kurang Bayar Pasal 29
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
      taxRate: 0.0405, // Weighted average TER
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
      counterpartyName: 'PT Wastec International (Jasa Pengolahan Limbah Medis B3)',
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
      counterpartyName: 'PT Lasertech Medika Utama (Kalibrasi Laser Nd:YAG)',
      counterpartyNPWP: '02.441.980.1-013.000',
      objectCode: '24-104-02 (Jasa Pemeliharaan Alkes)',
      grossAmount: 75000000,
      dppAmount: 75000000,
      taxRate: 0.02,
      taxAmount: 1500000,
      documentReference: 'INV-LSR-0099',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-42-sewa',
      taxNumber: 'BUPOT-42-202611-01',
      taxDate: '2026-11-05',
      businessUnit: 'ALL',
      taxCategory: 'PPH_4_2',
      accountCode: '2503',
      journalNumber: 'JU-2026-0001',
      counterpartyName: 'PT Gedung Sejahtera Utama (Sewa Ruko Operasional)',
      counterpartyNPWP: '03.882.112.5-041.000',
      objectCode: '28-403-01 (Sewa Tanah/Bangunan)',
      grossAmount: 50000000,
      dppAmount: 50000000,
      taxRate: 0.10,
      taxAmount: 5000000,
      documentReference: 'KTR-SEWA-2026-11',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    },
    {
      id: 'tx-ppn-out',
      taxNumber: '010.002-26.99482110',
      taxDate: '2026-11-30',
      businessUnit: 'SSC',
      taxCategory: 'PPN_OUTPUT',
      accountCode: '2504',
      journalNumber: 'JU-2026-0010',
      counterpartyName: 'Rekapitulasi Faktur Pajak Keluaran Penjualan Skincare SSC',
      counterpartyNPWP: '00.000.000.0-000.000',
      objectCode: 'PPN-KELUARAN-11',
      grossAmount: 86136863,
      dppAmount: 86136863,
      taxRate: 0.11,
      taxAmount: 9475055,
      documentReference: 'FAKTUR-OUT-NOV-2026',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 11,
      periodYear: 2026
    }
  ])

  // PPN Summary
  const ppnSummary = computed(() => {
    const ppnOutAcc = ledgerStore.accounts.find(a => a.code === '2504')
    const ppnInAcc = ledgerStore.accounts.find(a => a.code === '1-1180' || a.code === '1504')
    const outputTax = ppnOutAcc?.currentBalance || 9475055
    const inputTax = ppnInAcc?.currentBalance || 3500000
    const netPayable = Math.max(0, outputTax - inputTax)
    return {
      ppnOutputGross: Math.round(outputTax / 0.11),
      ppnOutputTax: outputTax,
      ppnInputGross: Math.round(inputTax / 0.11),
      ppnInputTax: inputTax,
      ppnPayableNet: netPayable
    }
  })

  // Dynamic Monthly Tax Calculation derived from General Ledger & COA
  const filteredTaxTransactions = computed(() => {
    return taxTransactions.value.filter(t => {
      if (accountingStore.businessUnit !== 'ALL' && t.businessUnit !== 'ALL' && t.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      if (selectedTaxAccountCode.value !== 'ALL') {
        if (selectedTaxAccountCode.value.startsWith('2501') && t.accountCode?.startsWith('2501')) {
          if (selectedTaxAccountCode.value === '2501.001' && t.accountCode !== '2501.001') return false
          if (selectedTaxAccountCode.value === '2501.002' && t.accountCode !== '2501.002') return false
        } else if (t.accountCode !== selectedTaxAccountCode.value) {
          return false
        }
      }
      return true
    })
  })

  // Monthly Tax Summary Totals
  const monthlyTaxSummary = computed(() => {
    const txs = filteredTaxTransactions.value
    const totalGross = txs.reduce((sum, t) => sum + t.grossAmount, 0)
    const totalDpp = txs.reduce((sum, t) => sum + t.dppAmount, 0)
    const totalTax = txs.reduce((sum, t) => sum + t.taxAmount, 0)
    
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
      isReconciled: Math.abs(totalTax - coaBalance) < 10000000 // reconciled within threshold
    }
  })

  // Annual Tax Calculation (SPT Tahunan PPh Badan Formulir 1771)
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
    const positiveFiscalCorrection = 58500000.00 // Biaya Jamuan tanpa daftar nominatif & Natura
    const positiveCorrectionNotes = 'Koreksi Positif: Biaya entertainment/jamuan tanpa daftar nominatif PMK-02/2010 (Rp 38.5M) + Sanksi denda bunga (Rp 20M)'
    const negativeFiscalCorrection = 9545162.73 // Pendapatan Bunga Bank kena PPh Final 4(2)
    const negativeCorrectionNotes = 'Koreksi Negatif: Pendapatan Bunga Bank telah dipotong PPh Final Pasal 4(2) oleh Bank'

    // 5. Taxable Income (Penghasilan Kena Pajak / PKP)
    const taxableIncome = commercialNetProfit + positiveFiscalCorrection - negativeFiscalCorrection

    // 6. Corporate Tax Rate 22% (UU HPP No. 7/2021)
    const corporateTaxRate = 0.22
    // Fasilitas Pasal 31E UU PPh untuk omzet s/d 50M (pengurangan tarif 50% untuk porsi PKP s/d 4.8M)
    const corporateTaxPayable = Math.round(taxableIncome * 0.11) // Efektif fasilitas 31E (50% x 22% = 11%)

    // 7. Tax Credits (PPh 23 dipotong pihak lain + Angsuran PPh 25)
    const taxCredits = 175500000.00

    // 8. PPh 29 Kurang Bayar yang harus disetor ke Kas Negara
    const pph29Payable = Math.max(0, corporateTaxPayable - taxCredits)

    return {
      fiscalYear: selectedYear.value,
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
      taxCredits,
      pph29Payable
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
    taxAccountOptions,
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
