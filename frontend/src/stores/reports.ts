import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAccountingStore } from './accounting'

export interface ProfitLossItem {
  code: string
  name: string
  kidinestAmount: number
  sscAmount: number
  totalAmount: number
  percentage: number
}

export interface AgingBucket {
  id: string
  entityName: string
  invoiceNumber: string
  invoiceDate: string
  dueDate: string
  current: number
  days1to30: number
  days31to60: number
  days61to90: number
  over90: number
  total: number
  status: 'CURRENT' | 'OVERDUE'
}

export const useReportsStore = defineStore('reports', () => {
  const accountingStore = useAccountingStore()

  // Profit & Loss Data Structure
  const revenueItems = ref<ProfitLossItem[]>([
    {
      code: '4-1100',
      name: 'Pendapatan Jasa Medis & Konsultasi',
      kidinestAmount: 165000000,
      sscAmount: 245000000,
      totalAmount: 410000000,
      percentage: 75.9
    },
    {
      code: '4-1200',
      name: 'Pendapatan Farmasi & Obat-obatan',
      kidinestAmount: 62000000,
      sscAmount: 20000000,
      totalAmount: 82000000,
      percentage: 15.2
    },
    {
      code: '4-1300',
      name: 'Pendapatan Skincare & Produk Slimming',
      kidinestAmount: 0,
      sscAmount: 48000000,
      totalAmount: 48000000,
      percentage: 8.9
    }
  ])

  const cogsItems = ref<ProfitLossItem[]>([
    {
      code: '5-1100',
      name: 'HPP Farmasi & Bahan Medis Habis Pakai',
      kidinestAmount: 32000000,
      sscAmount: 16500000,
      totalAmount: 48500000,
      percentage: 9.0
    },
    {
      code: '5-1200',
      name: 'Beban Jasa Medis Dokter (Share Medis)',
      kidinestAmount: 58000000,
      sscAmount: 87000000,
      totalAmount: 145000000,
      percentage: 26.9
    }
  ])

  const operationalExpenseItems = ref<ProfitLossItem[]>([
    {
      code: '6-1100',
      name: 'Gaji, Tunjangan & Upah Karyawan',
      kidinestAmount: 34500000,
      sscAmount: 33500000,
      totalAmount: 68000000,
      percentage: 12.6
    },
    {
      code: '6-1200',
      name: 'Beban Sewa Gedung Klinik',
      kidinestAmount: 12500000,
      sscAmount: 12500000,
      totalAmount: 25000000,
      percentage: 4.6
    },
    {
      code: '6-1300',
      name: 'Beban Listrik, Air, Gas & Internet',
      kidinestAmount: 4800000,
      sscAmount: 5000000,
      totalAmount: 9800000,
      percentage: 1.8
    },
    {
      code: '6-1400',
      name: 'Beban Penyusutan Aset & Peralatan Medis',
      kidinestAmount: 4166667,
      sscAmount: 6250000,
      totalAmount: 10416667,
      percentage: 1.9
    },
    {
      code: '6-1500',
      name: 'Beban Perlengkapan & Pemeliharaan Medis',
      kidinestAmount: 6200000,
      sscAmount: 8600000,
      totalAmount: 14800000,
      percentage: 2.7
    },
    {
      code: '6-1900',
      name: 'Beban Administrasi Bank & EDC Merchant Fee',
      kidinestAmount: 180000,
      sscAmount: 270000,
      totalAmount: 450000,
      percentage: 0.1
    }
  ])

  // Computed P&L Totals
  const profitLossSummary = computed(() => {
    const totalRev = revenueItems.value.reduce((s, i) => s + i.totalAmount, 0)
    const kidinestRev = revenueItems.value.reduce((s, i) => s + i.kidinestAmount, 0)
    const sscRev = revenueItems.value.reduce((s, i) => s + i.sscAmount, 0)

    const totalCogs = cogsItems.value.reduce((s, i) => s + i.totalAmount, 0)
    const kidinestCogs = cogsItems.value.reduce((s, i) => s + i.kidinestAmount, 0)
    const sscCogs = cogsItems.value.reduce((s, i) => s + i.sscAmount, 0)

    const grossProfit = totalRev - totalCogs
    const kidinestGrossProfit = kidinestRev - kidinestCogs
    const sscGrossProfit = sscRev - sscCogs

    const totalOpex = operationalExpenseItems.value.reduce((s, i) => s + i.totalAmount, 0)
    const kidinestOpex = operationalExpenseItems.value.reduce((s, i) => s + i.kidinestAmount, 0)
    const sscOpex = operationalExpenseItems.value.reduce((s, i) => s + i.sscAmount, 0)

    const netProfit = grossProfit - totalOpex
    const kidinestNetProfit = kidinestGrossProfit - kidinestOpex
    const sscNetProfit = sscGrossProfit - sscOpex

    const netMarginPercent = totalRev > 0 ? (netProfit / totalRev) * 100 : 0

    return {
      totalRev,
      kidinestRev,
      sscRev,
      totalCogs,
      kidinestCogs,
      sscCogs,
      grossProfit,
      kidinestGrossProfit,
      sscGrossProfit,
      totalOpex,
      kidinestOpex,
      sscOpex,
      netProfit,
      kidinestNetProfit,
      sscNetProfit,
      netMarginPercent: Math.round(netMarginPercent * 10) / 10
    }
  })

  // Balance Sheet Data
  const balanceSheet = computed(() => {
    // Current Assets
    const kasOperasional = 83600000 // 45.2jt + 38.4jt
    const bankOperasional = 427600000 // 185.6jt + 242jt
    const piutangPasien = 64200000
    const persediaan = 98500000
    const ppnMasukan = 8400000
    const totalAsetLancar = kasOperasional + bankOperasional + piutangPasien + persediaan + ppnMasukan

    // Fixed Assets
    const asetTetapGross = 750000000
    const akumulasiPenyusutan = 125000000
    const totalAsetTetapNet = asetTetapGross - akumulasiPenyusutan

    const totalAset = totalAsetLancar + totalAsetTetapNet

    // Liabilities
    const hutangDistributor = 58200000
    const hutangDokter = 42100000
    const hutangPajak = 22750000 // PPh 21 (7.85) + PPh 23 (2.4) + PPN (12.5)
    const totalKewajiban = hutangDistributor + hutangDokter + hutangPajak

    // Equity
    const modalDisetor = 800000000
    const labaDitahan = 215450000
    const labaTahunBerjalan = profitLossSummary.value.netProfit
    const totalEkuitas = modalDisetor + labaDitahan + labaTahunBerjalan

    const totalKewajibanDanEkuitas = totalKewajiban + totalEkuitas
    const balanceDifference = totalAset - totalKewajibanDanEkuitas

    return {
      kasOperasional,
      bankOperasional,
      piutangPasien,
      persediaan,
      ppnMasukan,
      totalAsetLancar,
      asetTetapGross,
      akumulasiPenyusutan,
      totalAsetTetapNet,
      totalAset,
      hutangDistributor,
      hutangDokter,
      hutangPajak,
      totalKewajiban,
      modalDisetor,
      labaDitahan,
      labaTahunBerjalan,
      totalEkuitas,
      totalKewajibanDanEkuitas,
      isBalanced: Math.abs(balanceDifference) < 1,
      balanceDifference
    }
  })

  // AR Aging
  const arAgingList = ref<AgingBucket[]>([
    {
      id: 'ar-1',
      entityName: 'Asuransi Sinarmas (Klaim Pasien Anak)',
      invoiceNumber: 'INV-CLM-2026-088',
      invoiceDate: '2026-02-15',
      dueDate: '2026-03-15',
      current: 18500000,
      days1to30: 0,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 18500000,
      status: 'CURRENT'
    },
    {
      id: 'ar-2',
      entityName: 'Asuransi Prudential (Klaim Rawat Inap & Observasi)',
      invoiceNumber: 'INV-CLM-2026-074',
      invoiceDate: '2026-01-20',
      dueDate: '2026-02-20',
      current: 0,
      days1to30: 24300000,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 24300000,
      status: 'OVERDUE'
    },
    {
      id: 'ar-3',
      entityName: 'Corporate PT Astra Sedaya (MCU Eksekutif SSC)',
      invoiceNumber: 'INV-CRP-2026-012',
      invoiceDate: '2026-02-28',
      dueDate: '2026-03-30',
      current: 21400000,
      days1to30: 0,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 21400000,
      status: 'CURRENT'
    }
  ])

  // AP Aging
  const apAgingList = ref<AgingBucket[]>([
    {
      id: 'ap-1',
      entityName: 'PT Kimia Farma Trading & Distribution',
      invoiceNumber: 'INV-KF-99881',
      invoiceDate: '2026-02-10',
      dueDate: '2026-03-12',
      current: 15540000,
      days1to30: 0,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 15540000,
      status: 'CURRENT'
    },
    {
      id: 'ap-2',
      entityName: 'PT Kalbe Farma Tbk (Divisi Farma & Nutrisi)',
      invoiceNumber: 'INV-KLB-7712',
      invoiceDate: '2026-01-25',
      dueDate: '2026-02-25',
      current: 0,
      days1to30: 22200000,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 22200000,
      status: 'OVERDUE'
    },
    {
      id: 'ap-3',
      entityName: 'PT Parazelsus Indonesia (Serum & Benang Estetika)',
      invoiceNumber: 'INV-PZ-3401',
      invoiceDate: '2026-02-22',
      dueDate: '2026-03-24',
      current: 20460000,
      days1to30: 0,
      days31to60: 0,
      days61to90: 0,
      over90: 0,
      total: 20460000,
      status: 'CURRENT'
    }
  ])

  return {
    revenueItems,
    cogsItems,
    operationalExpenseItems,
    profitLossSummary,
    balanceSheet,
    arAgingList,
    apAgingList
  }
})
