<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
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
    code: '1',
    name: 'AKTIVA - BCS',
    level: 1,
    amount: 26461564553.84,
    isHeader: true,
    isBold: true,
    accountCode: '1',
    children: [
      {
        id: '11',
        lineNo: 7,
        code: '11',
        name: 'KAS DAN BANK - BCS',
        level: 2,
        amount: 1277929020.88,
        isHeader: true,
        accountCode: '11',
        children: [
          {
            id: '1101',
            lineNo: 8,
            code: '1101',
            name: 'Kas Kasir & Operasional - BCS',
            level: 3,
            amount: 88600000.00,
            isHeader: true,
            accountCode: '1101',
            children: [
              { id: '1101.001', lineNo: 9, code: '1101.001', name: 'Kas Kasir Kidinest', level: 4, amount: 45200000.00, isHeader: false, accountCode: '1101.001' },
              { id: '1101.002', lineNo: 10, code: '1101.002', name: 'Kas Kasir Semarang Skin Center (SSC)', level: 4, amount: 38400000.00, isHeader: false, accountCode: '1101.002' },
              { id: '1101.003', lineNo: 11, code: '1101.003', name: 'Kas Kecil (Petty Cash Operasional)', level: 4, amount: 5000000.00, isHeader: false, accountCode: '1101.003' }
            ]
          },
          {
            id: '1102',
            lineNo: 12,
            code: '1102',
            name: 'Bank Operasional & Kliring - BCS',
            level: 3,
            amount: 1189329020.88,
            isHeader: true,
            accountCode: '1102',
            children: [
              { id: '1102.001', lineNo: 13, code: '1102.001', name: 'Bank BCA Operasional (Giro & EDC)', level: 4, amount: 680964487.55, isHeader: false, accountCode: '1102.001' },
              { id: '1102.002', lineNo: 14, code: '1102.002', name: 'Bank Mandiri Operasional & Payroll', level: 4, amount: 493864533.33, isHeader: false, accountCode: '1102.002' },
              { id: '1102.003', lineNo: 15, code: '1102.003', name: 'Settlement QRIS & Payment Gateway Transit', level: 4, amount: 14500000.00, isHeader: false, accountCode: '1102.003' }
            ]
          }
        ]
      },
      {
        id: '12',
        lineNo: 47,
        code: '12',
        name: 'PIUTANG - BCS',
        level: 2,
        amount: 1568434065.00,
        isHeader: true,
        accountCode: '12',
        children: [
          { id: '1201', lineNo: 48, code: '1201', name: 'Piutang Pasien Perawatan Rawat Jalan & Treatment', level: 3, amount: 485600000.00, isHeader: false, accountCode: '1201' },
          {
            id: '1202',
            lineNo: 49,
            code: '1202',
            name: 'Piutang Klaim Asuransi & Korporat',
            level: 3,
            amount: 1097834065.00,
            isHeader: true,
            accountCode: '1202',
            children: [
              { id: '1202.001', lineNo: 50, code: '1202.001', name: 'Piutang Asuransi Sinarmas', level: 4, amount: 420000000.00, isHeader: false, accountCode: '1202.001' },
              { id: '1202.002', lineNo: 51, code: '1202.002', name: 'Piutang Asuransi Prudential', level: 4, amount: 465000000.00, isHeader: false, accountCode: '1202.002' },
              { id: '1202.003', lineNo: 52, code: '1202.003', name: 'Piutang MCU Eksekutif Korporat', level: 4, amount: 212834065.00, isHeader: false, accountCode: '1202.003' }
            ]
          },
          { id: '1203', lineNo: 53, code: '1203', name: 'Cadangan Kerugian Piutang Tak Tertagih', level: 3, amount: -15000000.00, isHeader: false, accountCode: '1203' }
        ]
      },
      {
        id: '13',
        lineNo: 49,
        code: '13',
        name: 'Persediaan - BCS',
        level: 2,
        amount: 19335748296.88,
        isHeader: true,
        accountCode: '13',
        children: [
          { id: '1301', lineNo: 50, code: '1301', name: 'Persediaan Obat Resep & Paten Farmasi Kidinest', level: 3, amount: 4885000000.00, isHeader: false, accountCode: '1301' },
          { id: '1302', lineNo: 51, code: '1302', name: 'Persediaan Produk Skincare & Kosmetika SSC', level: 3, amount: 6442748296.88, isHeader: false, accountCode: '1302' },
          { id: '1303', lineNo: 52, code: '1303', name: 'Persediaan Bahan Medis Habis Pakai (BMHP)', level: 3, amount: 3208000000.00, isHeader: false, accountCode: '1303' },
          { id: '1304', lineNo: 53, code: '1304', name: 'Persediaan Vaksin & Serum Imunisasi Anak', level: 3, amount: 4800000000.00, isHeader: false, accountCode: '1304' }
        ]
      },
      {
        id: '14',
        lineNo: 73,
        code: '14',
        name: 'INVESTASI - BCS',
        level: 2,
        amount: 3234375045.00,
        isHeader: true,
        accountCode: '14',
        children: [
          { id: '1401', lineNo: 74, code: '1401', name: 'Investasi Tanah & Bangunan Cabang Baru', level: 3, amount: 2500000000.00, isHeader: false, accountCode: '1401' },
          { id: '1402', lineNo: 75, code: '1402', name: 'Investasi Pengembangan Riset Formulasi Skincare', level: 3, amount: 734375045.00, isHeader: false, accountCode: '1402' }
        ]
      },
      {
        id: '15',
        lineNo: 73,
        code: '15',
        name: 'BIAYA DIBAYAR DI MUKA - BCS',
        level: 2,
        amount: 2250000.00,
        isHeader: true,
        accountCode: '15',
        children: [
          { id: '1501', lineNo: 74, code: '1501', name: 'Sewa Ruko & Fasilitas Dibayar di Muka', level: 3, amount: 2250000.00, isHeader: false, accountCode: '1501' },
          { id: '1504', lineNo: 75, code: '1504', name: 'Uang Muka PPh 25 & Pajak Masukan (Kredit Pajak)', level: 3, amount: 0.00, isHeader: false, accountCode: '1504' }
        ]
      },
      {
        id: '16',
        lineNo: 73,
        code: '16',
        name: 'AKTIVA TETAP - BCS',
        level: 2,
        amount: 1042828126.08,
        isHeader: true,
        accountCode: '16',
        children: [
          { id: '1601', lineNo: 74, code: '1601', name: 'Bangunan & Renovasi Interior Klinik', level: 3, amount: 350000000.00, isHeader: false, accountCode: '1601' },
          { id: '1602', lineNo: 75, code: '1602', name: 'Mesin Laser Nd:YAG & Alkes Estetika SSC', level: 3, amount: 500000000.00, isHeader: false, accountCode: '1602' },
          { id: '1603', lineNo: 76, code: '1603', name: 'Peralatan USG 4D & Poli Anak Kidinest', level: 3, amount: 250000000.00, isHeader: false, accountCode: '1603' },
          { id: '1604', lineNo: 77, code: '1604', name: 'Inventaris Komputer, POS Kasir & Jaringan', level: 3, amount: 97828126.08, isHeader: false, accountCode: '1604' },
          { id: '1606', lineNo: 78, code: '1606', name: 'Akumulasi Penyusutan Aktiva Tetap', level: 3, amount: -155000000.00, isHeader: false, accountCode: '1606' }
        ]
      },
      {
        id: '1999',
        lineNo: 79,
        code: '1999',
        name: 'Temporary Opening - BCS',
        level: 2,
        amount: 0.00,
        isHeader: false,
        accountCode: '1999'
      }
    ]
  },
  {
    id: 'kewajiban_modal',
    lineNo: 73,
    name: 'KEWAJIBAN DAN MODAL',
    level: 1,
    amount: 26461564553.84,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'kewajiban',
        lineNo: 73,
        code: '2',
        name: 'KEWAJIBAN - BCS',
        level: 2,
        amount: 9962073223.05,
        isHeader: true,
        isBold: true,
        accountCode: '2',
        children: [
          {
            id: '21',
            lineNo: 73,
            code: '21',
            name: 'UANG MUKA & TITIPAN - BCS',
            level: 3,
            amount: 2147099360.00,
            isHeader: true,
            accountCode: '21',
            children: [
              {
                id: '2101',
                lineNo: 73,
                code: '2101',
                name: 'UANG TANDA JADI',
                level: 4,
                amount: 35500000.00,
                isHeader: false,
                accountCode: '2101'
              },
              {
                id: '2102',
                lineNo: 74,
                code: '2102',
                name: 'UANG MUKA PENJUALAN',
                level: 4,
                amount: 1904190360.00,
                isHeader: false,
                accountCode: '2102'
              },
              {
                id: '2103',
                lineNo: 75,
                code: '2103',
                name: 'TITIPAN BIAYA',
                level: 4,
                amount: 207409000.00,
                isHeader: false,
                accountCode: '2103'
              }
            ]
          },
          {
            id: '22',
            lineNo: 75,
            code: '22',
            name: 'HUTANG DAGANG - BCS',
            level: 3,
            amount: 4300000.00,
            isHeader: true,
            accountCode: '22',
            children: [
              { id: '2201', lineNo: 76, code: '2201', name: 'Hutang PBF Distributor Farmasi & Obat', level: 4, amount: 2800000.00, isHeader: false, accountCode: '2201' },
              { id: '2202', lineNo: 77, code: '2202', name: 'Hutang Supplier Kosmetika & Skincare SSC', level: 4, amount: 1500000.00, isHeader: false, accountCode: '2202' }
            ]
          },
          { id: '23', lineNo: 75, code: '23', name: 'HUTANG BANK - BCS', level: 3, amount: 6245750462.00, isHeader: false, accountCode: '23' },
          { id: '24', lineNo: 75, code: '24', name: 'HUTANG PIHAK KETIGA - BCS', level: 3, amount: 1363725333.01, isHeader: false, accountCode: '24' },
          {
            id: '25',
            lineNo: 75,
            code: '25',
            name: 'HUTANG PAJAK - BCS',
            level: 3,
            amount: 39249649.00,
            isHeader: true,
            accountCode: '25',
            children: [
              { id: '2501', lineNo: 76, code: '2501', name: 'Hutang PPh 21 (Dokter & Karyawan)', level: 4, amount: 22374594.00, isHeader: false, accountCode: '2501' },
              { id: '2502', lineNo: 77, code: '2502', name: 'Hutang PPh 23 (Jasa Medis, B3 & Maintenance)', level: 4, amount: 2400000.00, isHeader: false, accountCode: '2502' },
              { id: '2503', lineNo: 78, code: '2503', name: 'Hutang PPh Final Pasal 4(2) (Sewa Gedung Ruko)', level: 4, amount: 5000000.00, isHeader: false, accountCode: '2503' },
              { id: '2504', lineNo: 79, code: '2504', name: 'Hutang PPN Keluaran (Skincare & Treatment)', level: 4, amount: 9475055.00, isHeader: false, accountCode: '2504' }
            ]
          },
          { id: '26', lineNo: 79, code: '26', name: 'HUTANG PEMBEBANAN HPP - BCS', level: 3, amount: 161948419.04, isHeader: false, accountCode: '26' }
        ]
      },
      {
        id: 'ekuitas',
        lineNo: 79,
        code: '3',
        name: 'EKUITAS - BCS',
        level: 2,
        amount: 16499491330.79,
        isHeader: true,
        isBold: true,
        accountCode: '3',
        children: [
          {
            id: '31',
            lineNo: 79,
            code: '31',
            name: 'MODAL - BCS',
            level: 3,
            amount: 2367382942.15,
            isHeader: true,
            accountCode: '31',
            children: [
              { id: '3101', lineNo: 80, code: '3101', name: 'Modal Disetor Pemilik Klinik Kidinest & SSC', level: 4, amount: 2367382942.15, isHeader: false, accountCode: '3101' }
            ]
          },
          {
            id: '32',
            lineNo: 80,
            code: '32',
            name: 'LABA AKUMULASI - BCS',
            level: 3,
            amount: 14132108388.64,
            isHeader: true,
            accountCode: '32',
            children: [
              { id: '3201', lineNo: 81, code: '3201', name: 'Laba Ditahan (Retained Earnings)', level: 4, amount: 11496921049.78, isHeader: false, accountCode: '3201' },
              { id: '3202', lineNo: 82, code: '3202', name: 'Laba Bersih Periode Berjalan', level: 4, amount: 2635187338.86, isHeader: false, accountCode: '3202' }
            ]
          },
          { id: '3999', lineNo: 83, code: '3999', name: 'Opening Balance Equity - BCS', level: 3, amount: 0.00, isHeader: false, accountCode: '3999' }
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

    <!-- Official Report Signature Section matching Image 1 -->
    <ReportSignatures />
  </div>
</template>
