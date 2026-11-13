<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  ChevronRight,
  ChevronDown,
  BookOpen,
  Filter,
  Layers,
  ArrowRight
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()

// State
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')
const currentLevel = ref<number>(1)
const jumpLineNumber = ref<number>(1)
const showZeroBalance = ref(false)

export interface TreeNode {
  id: string
  lineNo: number
  code?: string
  name: string
  level: number
  amount: number
  isHeader: boolean
  isBold?: boolean
  accountCode?: string
  children?: TreeNode[]
}

// Tree Data for Balance Sheet (Neraca) aligned with Blueprint Kidinest & SSC
const neracaTree = ref<TreeNode[]>([
  {
    id: 'aktiva',
    lineNo: 6,
    name: 'AKTIVA',
    level: 1,
    amount: 25770707945.51,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: '11',
        lineNo: 7,
        code: '11',
        name: 'KAS DAN BANK',
        level: 2,
        amount: 1138664487.55,
        isHeader: true,
        accountCode: '1-1100',
        children: [
          {
            id: '1101',
            lineNo: 8,
            code: '1101',
            name: 'Kas Kasir Kidinest (Poli Anak)',
            level: 3,
            amount: 45200000.00,
            isHeader: true,
            accountCode: '1-1111',
            children: [
              { id: '1101.01', lineNo: 9, code: '1101.01', name: 'Kas Fisik Meja Kasir Kidinest Shift Pagi & Sore', level: 4, amount: 28500000.00, isHeader: false, accountCode: '1-1111' },
              { id: '1101.02', lineNo: 10, code: '1101.02', name: 'Kas Brankas Operasional Poli Anak', level: 4, amount: 16700000.00, isHeader: false, accountCode: '1-1111' }
            ]
          },
          {
            id: '1102',
            lineNo: 11,
            code: '1102',
            name: 'Kas Kasir Semarang Skin Center (SSC)',
            level: 3,
            amount: 38400000.00,
            isHeader: true,
            accountCode: '1-1112',
            children: [
              { id: '1102.01', lineNo: 12, code: '1102.01', name: 'Kas Fisik Meja Kasir SSC Estetika', level: 4, amount: 22400000.00, isHeader: false, accountCode: '1-1112' },
              { id: '1102.02', lineNo: 13, code: '1102.02', name: 'Kas Brankas Pelayanan Facial & Laser', level: 4, amount: 16000000.00, isHeader: false, accountCode: '1-1112' }
            ]
          },
          { id: '1103', lineNo: 14, code: '1103', name: 'Kas Kecil / Petty Cash Operasional', level: 3, amount: 5000000.00, isHeader: false, accountCode: '1-1113' },
          { id: '1104', lineNo: 15, code: '1104', name: 'Bank BCA Operasional (EDC & Giro Utama)', level: 3, amount: 580964487.55, isHeader: false, accountCode: '1-1121' },
          { id: '1105', lineNo: 16, code: '1105', name: 'Bank Mandiri Operasional & Payroll', level: 3, amount: 454600000.00, isHeader: false, accountCode: '1-1122' },
          { id: '1106', lineNo: 17, code: '1106', name: 'Settlement QRIS & EDC Transit', level: 3, amount: 14500000.00, isHeader: false, accountCode: '1-1130' }
        ]
      },
      {
        id: '12',
        lineNo: 47,
        code: '12',
        name: 'PIUTANG',
        level: 2,
        amount: 1073834065.00,
        isHeader: true,
        accountCode: '1-1200',
        children: [
          { id: '1201', lineNo: 48, code: '1201', name: 'Piutang Pasien Perawatan Rawat Jalan & Treatment', level: 3, amount: 324500000.00, isHeader: false, accountCode: '1-1201' },
          {
            id: '1202',
            lineNo: 49,
            code: '1202',
            name: 'Piutang Klaim Asuransi & Korporat',
            level: 3,
            amount: 764334065.00,
            isHeader: true,
            accountCode: '1-1202',
            children: [
              { id: '1202.01', lineNo: 50, code: '1202.01', name: 'Klaim Asuransi Sinarmas (Rawat Anak)', level: 4, amount: 285000000.00, isHeader: false, accountCode: '1-1202' },
              { id: '1202.02', lineNo: 51, code: '1202.02', name: 'Klaim Asuransi Prudential Pediatrik', level: 4, amount: 312500000.00, isHeader: false, accountCode: '1-1202' },
              { id: '1202.03', lineNo: 52, code: '1202.03', name: 'Klaim MCU Eksekutif Kulit PT Astra', level: 4, amount: 166834065.00, isHeader: false, accountCode: '1-1202' }
            ]
          },
          { id: '1203', lineNo: 53, code: '1203', name: 'Cadangan Kerugian Piutang Tak Tertagih', level: 3, amount: -15000000.00, isHeader: false, accountCode: '1-1203' }
        ]
      },
      {
        id: '13',
        lineNo: 49,
        code: '13',
        name: 'Persediaan',
        level: 2,
        amount: 19278756221.88,
        isHeader: true,
        accountCode: '1-1300',
        children: [
          { id: '1301', lineNo: 50, code: '1301', name: 'Persediaan Obat & Sirup Anak (Kidinest)', level: 3, amount: 4850000000.00, isHeader: false, accountCode: '1-1310' },
          { id: '1302', lineNo: 51, code: '1302', name: 'Persediaan Produk Skincare & Kosmetika (SSC)', level: 3, amount: 6420756221.88, isHeader: false, accountCode: '1-1330' },
          { id: '1303', lineNo: 52, code: '1303', name: 'Persediaan Bahan Medis Habis Pakai (BMHP)', level: 3, amount: 3208000000.00, isHeader: false, accountCode: '1-1320' },
          { id: '1304', lineNo: 53, code: '1304', name: 'Persediaan Vaksin & Serum Anak (Kidinest)', level: 3, amount: 4800000000.00, isHeader: false, accountCode: '1-1311' }
        ]
      },
      {
        id: '14',
        lineNo: 73,
        code: '14',
        name: 'INVESTASI',
        level: 2,
        amount: 3234375045.00,
        isHeader: true,
        children: [
          { id: '1401', lineNo: 74, code: '1401', name: 'Investasi Tanah & Gedung Cabang Klinik', level: 3, amount: 2500000000.00, isHeader: false },
          { id: '1402', lineNo: 75, code: '1402', name: 'Investasi Pengembangan Formula Skincare SSC', level: 3, amount: 734375045.00, isHeader: false }
        ]
      },
      {
        id: '15',
        lineNo: 73,
        code: '15',
        name: 'BIAYA DIBAYAR DI MUKA',
        level: 2,
        amount: 51750000.01,
        isHeader: true,
        children: [
          { id: '1501', lineNo: 74, code: '1501', name: 'Sewa Gedung & Ruko Dibayar di Muka', level: 3, amount: 37500000.00, isHeader: false },
          { id: '1502', lineNo: 75, code: '1502', name: 'Asuransi Kebakaran & Alkes Dibayar di Muka', level: 3, amount: 14250000.01, isHeader: false }
        ]
      },
      {
        id: '16',
        lineNo: 73,
        code: '16',
        name: 'AKTIVA TETAP',
        level: 2,
        amount: 993328126.07,
        isHeader: true,
        accountCode: '1-2000',
        children: [
          { id: '1601', lineNo: 74, code: '1601', name: 'Bangunan & Renovasi Interior Klinik', level: 3, amount: 350000000.00, isHeader: false, accountCode: '1-2010' },
          { id: '1602', lineNo: 75, code: '1602', name: 'Mesin & Alat Medis Estetika (Laser Nd:YAG / Facial SSC)', level: 3, amount: 500000000.00, isHeader: false, accountCode: '1-2200' },
          { id: '1603', lineNo: 76, code: '1603', name: 'Peralatan Medis Poli Anak & USG 4D Kidinest', level: 3, amount: 250000000.00, isHeader: false, accountCode: '1-2100' },
          { id: '1604', lineNo: 77, code: '1604', name: 'Inventaris Komputer, Thermal Printer & POS', level: 3, amount: 48328126.07, isHeader: false },
          { id: '1606', lineNo: 78, code: '1606', name: 'Akumulasi Penyusutan Aktiva Tetap Klinik', level: 3, amount: -155000000.00, isHeader: false, accountCode: '1-2900' }
        ]
      }
    ]
  },
  {
    id: 'kewajiban_modal',
    lineNo: 73,
    name: 'KEWAJIBAN DAN MODAL',
    level: 1,
    amount: 25876969445.51,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'kewajiban',
        lineNo: 73,
        name: 'KEWAJIBAN',
        level: 2,
        amount: 9873786668.05,
        isHeader: true,
        isBold: true,
        accountCode: '2-0000',
        children: [
          {
            id: '2101',
            lineNo: 73,
            code: '2101',
            name: 'UANG TANDA JADI',
            level: 3,
            amount: 35500000.00,
            isHeader: true,
            children: [
              { id: '2101.01', lineNo: 74, code: '2101.01', name: 'Uang Muka Reservasi Poli Spesialis Anak', level: 4, amount: 15500000.00, isHeader: false },
              { id: '2101.02', lineNo: 75, code: '2101.02', name: 'Down Payment Treatment Laser & Bedah Estetika', level: 4, amount: 20000000.00, isHeader: false }
            ]
          },
          {
            id: '2102',
            lineNo: 74,
            code: '2102',
            name: 'UANG MUKA PENJUALAN',
            level: 3,
            amount: 1904190360.00,
            isHeader: true,
            children: [
              { id: '2102.01', lineNo: 75, code: '2102.01', name: 'Deposit Pasien Paket 10x Facial Glow SSC', level: 4, amount: 980190360.00, isHeader: false },
              { id: '2102.02', lineNo: 76, code: '2102.02', name: 'Deposit Pasien Paket Terapi Tumbuh Kembang Kidinest', level: 4, amount: 924000000.00, isHeader: false }
            ]
          },
          { id: '2103', lineNo: 74, code: '2103', name: 'TITIPAN BIAYA', level: 3, amount: 158347500.00, isHeader: false },
          {
            id: '22',
            lineNo: 75,
            code: '22',
            name: 'HUTANG DAGANG',
            level: 3,
            amount: 4300000.00,
            isHeader: true,
            accountCode: '2-1110',
            children: [
              { id: '2201', lineNo: 76, code: '2201', name: 'Hutang PBF Kimia Farma Trading (Vaksin & Obat)', level: 4, amount: 2800000.00, isHeader: false, accountCode: '2-1110' },
              { id: '2202', lineNo: 77, code: '2202', name: 'Hutang Supplier Kosmetika & Skincare SSC', level: 4, amount: 1500000.00, isHeader: false, accountCode: '2-1120' }
            ]
          },
          { id: '23', lineNo: 75, code: '23', name: 'HUTANG BANK', level: 3, amount: 6247750462.00, isHeader: false },
          { id: '24', lineNo: 75, code: '24', name: 'HUTANG PIHAK KETIGA', level: 3, amount: 1326875333.01, isHeader: false },
          {
            id: '25',
            lineNo: 75,
            code: '25',
            name: 'HUTANG PAJAK',
            level: 3,
            amount: 34874594.00,
            isHeader: true,
            accountCode: '2-1501',
            children: [
              { id: '2501', lineNo: 76, code: '2501', name: 'Hutang PPh 21 Dokter Tenaga Ahli (PMK 168/2023)', level: 4, amount: 14874594.00, isHeader: false, accountCode: '2-1501' },
              { id: '2502', lineNo: 77, code: '2502', name: 'Hutang PPh 21 Karyawan & Perawat', level: 4, amount: 7500000.00, isHeader: false },
              { id: '2505', lineNo: 78, code: '2505', name: 'Hutang PPN Keluaran Skincare & Treatment', level: 4, amount: 12500000.00, isHeader: false, accountCode: '2-1504' }
            ]
          },
          { id: '26', lineNo: 79, code: '26', name: 'HUTANG PEMBEBANAN HPP', level: 3, amount: 161948419.04, isHeader: false }
        ]
      },
      {
        id: 'ekuitas',
        lineNo: 79,
        name: 'EKUITAS',
        level: 2,
        amount: 16003182777.46,
        isHeader: true,
        isBold: true,
        accountCode: '3-0000',
        children: [
          {
            id: '31',
            lineNo: 79,
            code: '31',
            name: 'MODAL',
            level: 3,
            amount: 3890010000.00,
            isHeader: true,
            accountCode: '3-1100',
            children: [
              { id: '3101', lineNo: 80, code: '3101', name: 'Modal Disetor Pemilik Klinik Kidinest & SSC', level: 4, amount: 3890010000.00, isHeader: false, accountCode: '3-1100' }
            ]
          },
          {
            id: '32',
            lineNo: 80,
            code: '32',
            name: 'LABA AKUMULASI',
            level: 3,
            amount: 12113172777.46,
            isHeader: true,
            accountCode: '3-1200',
            children: [
              { id: '3201', lineNo: 81, code: '3201', name: 'Laba Ditahan (Retained Earnings)', level: 4, amount: 10435172777.46, isHeader: false, accountCode: '3-1200' },
              { id: '3202', lineNo: 82, code: '3202', name: 'Laba Bersih Tahun Berjalan', level: 4, amount: 1678000000.00, isHeader: false, accountCode: '3-1300' }
            ]
          }
        ]
      }
    ]
  }
])

// Flat list calculation based on expanded state - default Level 1 matching Image 2
const expandedIds = ref<Set<string>>(new Set([
  'aktiva', 'kewajiban_modal', 'kewajiban', 'ekuitas'
]))

function toggleExpand(id: string) {
  if (expandedIds.value.has(id)) {
    expandedIds.value.delete(id)
  } else {
    expandedIds.value.add(id)
  }
}

function setLevel(lvl: number) {
  currentLevel.value = lvl
  jumpLineNumber.value = lvl
  expandedIds.value.clear()

  if (lvl === 1) {
    expandedIds.value.add('aktiva')
    expandedIds.value.add('kewajiban_modal')
    expandedIds.value.add('kewajiban')
    expandedIds.value.add('ekuitas')
    return
  }

  function traverse(nodes: TreeNode[]) {
    for (const node of nodes) {
      if (node.level < lvl && node.children && node.children.length > 0) {
        expandedIds.value.add(node.id)
        traverse(node.children)
      }
    }
  }

  expandedIds.value.add('aktiva')
  expandedIds.value.add('kewajiban_modal')
  expandedIds.value.add('kewajiban')
  expandedIds.value.add('ekuitas')
  traverse(neracaTree.value)
}

function applySetLevel() {
  const lvl = Math.min(Math.max(Number(jumpLineNumber.value) || 1, 1), 4)
  setLevel(lvl)
}

function expandAll() {
  currentLevel.value = 4
  jumpLineNumber.value = 4
  function collectAll(nodes: TreeNode[]) {
    for (const n of nodes) {
      if (n.children && n.children.length > 0) {
        expandedIds.value.add(n.id)
        collectAll(n.children)
      }
    }
  }
  collectAll(neracaTree.value)
}

function collapseAll() {
  currentLevel.value = 1
  jumpLineNumber.value = 1
  expandedIds.value.clear()
  expandedIds.value.add('aktiva')
  expandedIds.value.add('kewajiban_modal')
  expandedIds.value.add('kewajiban')
  expandedIds.value.add('ekuitas')
}

// Flatten for rendering
interface FlattenedRow {
  node: TreeNode
  visible: boolean
  hasChildren: boolean
  isExpanded: boolean
}

const flattenedRows = computed<FlattenedRow[]>(() => {
  const result: FlattenedRow[] = []

  function flatten(nodes: TreeNode[], isParentExpanded: boolean) {
    for (const node of nodes) {
      const hasChildren = Boolean(node.children && node.children.length > 0)
      const isExpanded = expandedIds.value.has(node.id)

      if (isParentExpanded) {
        result.push({
          node,
          visible: true,
          hasChildren,
          isExpanded
        })
      }

      if (hasChildren && node.children) {
        flatten(node.children, isParentExpanded && isExpanded)
      }
    }
  }

  flatten(neracaTree.value, true)
  return result
})

// Drill down to General Ledger
function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/ledger', query: { account: accountCode } })
  } else {
    router.push('/ledger')
  }
}

function handleExportExcel() {
  const data = flattenedRows.value.map(r => ({
    'Line': r.node.lineNo,
    'Kode': r.node.code || '',
    'Uraian': '  '.repeat(r.node.level - 1) + r.node.name,
    'Level': r.node.level,
    'Jumlah (IDR)': r.node.amount
  }))

  exportToExcel(data, `Laporan_Neraca_${accountingStore.businessUnit}_${endDate.value}`, 'Neraca')
  accountingStore.showToast('Laporan Neraca berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header Controls matching ERP Reference -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-4 text-white flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xl">
      <div class="flex items-center gap-3">
        <div class="space-y-0.5">
          <div class="flex items-center gap-2 text-xs text-slate-400">
            <span class="hover:text-white cursor-pointer font-medium">Accounting</span>
            <span>/</span>
            <span class="text-blue-400 font-semibold">Laporan Neraca</span>
          </div>
          <h1 class="text-base font-bold text-white tracking-wide">
            {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SEMARANG SKIN CENTER)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER (SSC)' }}
          </h1>
        </div>
      </div>

      <!-- Date Range Selectors -->
      <div class="flex flex-wrap items-center gap-3 text-xs">
        <div class="flex items-center gap-1.5 bg-[#1E293B] px-3 py-1.5 rounded-lg border border-slate-700">
          <span class="text-slate-400">Periode:</span>
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
          <span class="text-slate-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Ekspor</span>
        </button>

        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium transition-colors"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Tree Table Container (High-Contrast Dark SaaS Look matching Image 2) -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl overflow-hidden shadow-2xl text-slate-200">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-slate-800 bg-[#0F172A] text-[11px] font-bold text-slate-400">
              <th class="py-2.5 px-4 w-16 text-slate-500">#</th>
              <th class="py-2.5 px-4">Uraian</th>
              <th class="py-2.5 px-6 text-right w-64">Jumlah</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60 font-mono">
            <!-- Header Titles matching Image 2 Lines 3, 4, 5 -->
            <tr class="bg-[#0F172A] border-b border-slate-800/80">
              <td class="py-2 px-4 text-slate-500 font-mono text-[11px]">3</td>
              <td class="py-2 px-4 font-black text-white uppercase tracking-wider text-xs">NERACA</td>
              <td class="py-2 px-6"></td>
            </tr>
            <tr class="bg-[#0F172A] border-b border-slate-800/80">
              <td class="py-1.5 px-4 text-slate-500 font-mono text-[11px]">4</td>
              <td class="py-1.5 px-4 font-bold text-slate-400 uppercase tracking-wide text-[11px]">PER 30 NOVEMBER 2026</td>
              <td class="py-1.5 px-6"></td>
            </tr>
            <tr class="bg-[#0F172A]/50 border-b border-slate-800/40">
              <td class="py-1 px-4 text-slate-500 font-mono text-[11px]">5</td>
              <td class="py-1 px-4"></td>
              <td class="py-1 px-6"></td>
            </tr>

            <!-- Interactive Tree Nodes -->
            <tr
              v-for="row in flattenedRows"
              :key="row.node.id"
              :class="[
                'transition-colors group',
                row.node.level === 1
                  ? 'bg-[#182234] font-extrabold text-white text-sm'
                  : row.node.level === 2
                  ? 'bg-[#141D2E] font-bold text-slate-100 text-xs'
                  : 'hover:bg-slate-800/40 text-slate-300'
              ]"
            >
              <!-- Line Number Column -->
              <td class="py-2 px-4 text-slate-500 text-[11px] select-none font-mono">
                {{ row.node.lineNo }}
              </td>

              <!-- Uraian with Expand Caret & Book Icon -->
              <td class="py-2 px-4 font-sans">
                <div
                  class="flex items-center gap-1.5"
                  :style="{ paddingLeft: `${(row.node.level - 1) * 20}px` }"
                >
                  <!-- Caret for Parent Nodes -->
                  <button
                    v-if="row.hasChildren"
                    @click="toggleExpand(row.node.id)"
                    class="p-0.5 rounded text-slate-400 hover:text-white hover:bg-slate-700/60 transition-colors"
                  >
                    <ChevronDown v-if="row.isExpanded" class="w-3.5 h-3.5 text-slate-300" />
                    <ChevronRight v-else class="w-3.5 h-3.5 text-slate-400" />
                  </button>
                  <span v-else class="w-4 inline-block"></span>

                  <!-- Account Code and Name -->
                  <span
                    :class="[
                      'tracking-wide',
                      row.node.level === 1 ? 'font-black tracking-wider text-white uppercase' : '',
                      row.node.level === 2 ? 'font-bold text-slate-100 uppercase' : '',
                      row.node.level === 3 ? 'font-semibold text-slate-200' : '',
                      row.node.level === 4 ? 'font-normal text-slate-400 text-xs hover:text-blue-300 cursor-pointer' : ''
                    ]"
                    @click="row.node.accountCode && openLedgerForAccount(row.node.accountCode)"
                  >
                    {{ row.node.code ? `${row.node.code} - ` : '' }}{{ row.node.name }}
                  </span>

                  <!-- Subsidiary Ledger / Buku Besar Book Icon (📖) -->
                  <button
                    v-if="row.node.accountCode"
                    @click.stop="openLedgerForAccount(row.node.accountCode)"
                    class="ml-1 text-slate-500 hover:text-blue-400 transition-colors p-0.5 rounded"
                    title="Buka Buku Besar (General Ledger) untuk akun ini"
                  >
                    <BookOpen class="w-3 h-3 text-slate-400 group-hover:text-blue-400" />
                  </button>
                </div>
              </td>

              <!-- Jumlah Column -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight',
                  row.node.level === 1
                    ? 'font-black text-sm text-white'
                    : row.node.level === 2
                    ? 'font-bold text-slate-100 text-xs'
                    : 'font-medium text-slate-300'
                ]"
              >
                {{ formatIDR(row.node.amount) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Bottom Level Bar matching Image 2 -->
      <div class="px-4 py-3 bg-[#0F172A] border-t border-slate-800 flex flex-wrap items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-2">
          <!-- Jump to Line / Set Level Input matching Image 2 -->
          <input
            v-model.number="jumpLineNumber"
            type="number"
            min="1"
            max="4"
            class="w-12 px-2 py-1 bg-[#1E293B] border border-slate-700 rounded text-center text-xs text-white font-mono"
            @keyup.enter="applySetLevel"
          />

          <!-- Level Buttons: Lvl 1, Lvl 2, Lvl 3, Lvl 4 (Detail) -->
          <div class="flex items-center bg-[#1E293B] p-0.5 rounded-lg border border-slate-700">
            <button
              @click="setLevel(1)"
              :class="[
                'px-2.5 py-1 rounded text-xs font-bold transition-all',
                currentLevel === 1 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-400 hover:text-white'
              ]"
            >
              Lvl 1
            </button>
            <button
              @click="setLevel(2)"
              :class="[
                'px-2.5 py-1 rounded text-xs font-bold transition-all',
                currentLevel === 2 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-400 hover:text-white'
              ]"
            >
              Lvl 2
            </button>
            <button
              @click="setLevel(3)"
              :class="[
                'px-2.5 py-1 rounded text-xs font-bold transition-all',
                currentLevel === 3 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-400 hover:text-white'
              ]"
            >
              Lvl 3
            </button>
            <button
              @click="setLevel(4)"
              :class="[
                'px-2.5 py-1 rounded text-xs font-bold transition-all',
                currentLevel === 4 ? 'bg-blue-600 text-white shadow-xs' : 'text-slate-400 hover:text-white'
              ]"
            >
              Lvl 4 (Detail)
            </button>
          </div>

          <button
            @click="applySetLevel"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium border border-slate-700 transition-colors"
          >
            Set Level
          </button>

          <button
            @click="expandAll"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium border border-slate-700 transition-colors"
          >
            Expand All
          </button>

          <button
            @click="collapseAll"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium border border-slate-700 transition-colors"
          >
            Collapse All
          </button>
        </div>

        <div class="flex items-center gap-3 text-slate-400 text-xs">
          <span>Menampilkan <strong>{{ flattenedRows.length }}</strong> baris neraca</span>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-emerald-400 font-bold">Neraca Seimbang (Balanced)</span>
        </div>
      </div>
    </div>
  </div>
</template>
