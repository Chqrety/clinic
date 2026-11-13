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

  // Profit & Loss Data Structure based on Blueprint Kidinest & SSC
  const revenueItems = ref<ProfitLossItem[]>([
    {
      code: '4-1110',
      name: 'Pendapatan Jasa Konsultasi Spesialis Anak (Sp.A)',
      kidinestAmount: 115000000,
      sscAmount: 0,
      totalAmount: 115000000,
      percentage: 21.3
    },
    {
      code: '4-1120',
      name: 'Pendapatan Tindakan Vaksinasi & Imunisasi Anak',
      kidinestAmount: 32000000,
      sscAmount: 0,
      totalAmount: 32000000,
      percentage: 5.9
    },
    {
      code: '4-1130',
      name: 'Pendapatan Layanan Terapi & Tumbuh Kembang Anak',
      kidinestAmount: 18000000,
      sscAmount: 0,
      totalAmount: 18000000,
      percentage: 3.3
    },
    {
      code: '4-1140',
      name: 'Pendapatan Penjualan Obat Resep Farmasi Anak',
      kidinestAmount: 62000000,
      sscAmount: 0,
      totalAmount: 62000000,
      percentage: 11.5
    },
    {
      code: '4-2110',
      name: 'Pendapatan Konsultasi Dokter Estetika & Kulit (Sp.DVE)',
      kidinestAmount: 0,
      sscAmount: 65000000,
      totalAmount: 65000000,
      percentage: 12.0
    },
    {
      code: '4-2120',
      name: 'Pendapatan Tindakan Medis Laser & Cauter SSC',
      kidinestAmount: 0,
      sscAmount: 145000000,
      totalAmount: 145000000,
      percentage: 26.9
    },
    {
      code: '4-2130',
      name: 'Pendapatan Tindakan Facial & Treatment Terapis',
      kidinestAmount: 0,
      sscAmount: 35000000,
      totalAmount: 35000000,
      percentage: 6.5
    },
    {
      code: '4-2140',
      name: 'Pendapatan Penjualan Produk Skincare & Kosmetika',
      kidinestAmount: 0,
      sscAmount: 68000000,
      totalAmount: 68000000,
      percentage: 12.6
    }
  ])

  const cogsItems = ref<ProfitLossItem[]>([
    {
      code: '5-1101',
      name: 'HPP Obat Resep & Paten Farmasi Kidinest',
      kidinestAmount: 24000000,
      sscAmount: 0,
      totalAmount: 24000000,
      percentage: 4.4
    },
    {
      code: '5-1102',
      name: 'HPP Vaksin & Serum Imunisasi Kidinest',
      kidinestAmount: 8000000,
      sscAmount: 0,
      totalAmount: 8000000,
      percentage: 1.5
    },
    {
      code: '5-1103',
      name: 'HPP Produk Skincare & Kosmetika Terjual SSC',
      kidinestAmount: 0,
      sscAmount: 16500000,
      totalAmount: 16500000,
      percentage: 3.1
    },
    {
      code: '5-1105',
      name: 'Beban Pemakaian Bahan Medis Habis Pakai (BMHP)',
      kidinestAmount: 3800000,
      sscAmount: 6000000,
      totalAmount: 9800000,
      percentage: 1.8
    },
    {
      code: '5-1106',
      name: 'Beban Bagi Hasil Dokter Spesialis Anak (Kidinest)',
      kidinestAmount: 58000000,
      sscAmount: 0,
      totalAmount: 58000000,
      percentage: 10.7
    },
    {
      code: '5-1107',
      name: 'Beban Bagi Hasil Dokter Estetika & Kulit (SSC)',
      kidinestAmount: 0,
      sscAmount: 78500000,
      totalAmount: 78500000,
      percentage: 14.5
    },
    {
      code: '5-1108',
      name: 'Beban Komisi Tindakan Beautician / Terapis SSC',
      kidinestAmount: 0,
      sscAmount: 8500000,
      totalAmount: 8500000,
      percentage: 1.6
    }
  ])

  const operationalExpenseItems = ref<ProfitLossItem[]>([
    {
      code: '6-1101',
      name: 'Beban Gaji Staf, Perawat, Apoteker & Kasir',
      kidinestAmount: 34500000,
      sscAmount: 33500000,
      totalAmount: 68000000,
      percentage: 12.6
    },
    {
      code: '6-1102',
      name: 'Beban Listrik PLN, Air PDAM & Internet Klinik',
      kidinestAmount: 4800000,
      sscAmount: 5000000,
      totalAmount: 9800000,
      percentage: 1.8
    },
    {
      code: '6-1103',
      name: 'Beban Pengolahan Limbah Medis B3',
      kidinestAmount: 1800000,
      sscAmount: 2400000,
      totalAmount: 4200000,
      percentage: 0.8
    },
    {
      code: '6-1104',
      name: 'Beban Kalibrasi & Servis Laser / Alkes Medis',
      kidinestAmount: 1600000,
      sscAmount: 4000000,
      totalAmount: 5600000,
      percentage: 1.0
    },
    {
      code: '6-1105',
      name: 'Beban Sewa Gedung & Ruko Klinik',
      kidinestAmount: 12500000,
      sscAmount: 12500000,
      totalAmount: 25000000,
      percentage: 4.6
    },
    {
      code: '6-1107',
      name: 'Beban Penyusutan Aktiva Tetap Klinik',
      kidinestAmount: 4166667,
      sscAmount: 6250000,
      totalAmount: 10416667,
      percentage: 1.9
    },
    {
      code: '6-1201',
      name: 'Beban Administrasi Bank & MDR Digital Payment',
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

  // Balance Sheet Data per Blueprint
  const balanceSheet = computed(() => {
    // Current Assets
    const kasOperasional = 88600000 // 45.2jt (Kidinest) + 38.4jt (SSC) + 5jt (Petty Cash)
    const bankOperasional = 442100000 // 242jt (BCA) + 185.6jt (Mandiri) + 14.5jt (Kliring)
    const piutangPasien = 64200000 // 21.4jt (Pasien) + 42.8jt (Asuransi)
    const persediaan = 98500000 // 38.2jt (Obat) + 24.5jt (Vaksin) + 12.4jt (BMHP) + 23.4jt (Skincare)
    const ppnMasukan = 8400000
    const totalAsetLancar = kasOperasional + bankOperasional + piutangPasien + persediaan + ppnMasukan

    // Fixed Assets
    const asetTetapGross = 750000000 // 250jt (Poli Anak) + 500jt (Laser SSC)
    const akumulasiPenyusutan = 125000000
    const totalAsetTetapNet = asetTetapGross - akumulasiPenyusutan

    const totalAset = totalAsetLancar + totalAsetTetapNet

    // Liabilities
    const hutangDistributor = 58200000 // 38.2jt (PBF) + 20jt (Skincare)
    const hutangDokter = 42100000 // 18.4jt (Dokter Anak) + 23.7jt (Dokter Kulit)
    const hutangPajak = 22750000 // PPh 21 Dokter (7.85jt) + PPh 23 (2.4jt) + PPN Keluaran (12.5jt)
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
      entityName: 'Asuransi Sinarmas (Klaim Rawat Jalan Anak)',
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
      entityName: 'Asuransi Prudential (Klaim Observasi Pediatrik)',
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
      entityName: 'Corporate PT Astra Sedaya (MCU Eksekutif Kulit SSC)',
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
      entityName: 'PT Kimia Farma Trading & Distribution (PBF Farmasi)',
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
      entityName: 'PT Kalbe Farma Tbk (Nutrisi & Vaksin Anak)',
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
      entityName: 'PT Parazelsus Indonesia (Serum & Kosmetika SSC)',
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
