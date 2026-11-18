<script setup lang="ts">
import { ref, computed, onMounted, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import { exportProfitLossToExcel, exportToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  ChevronRight,
  ChevronDown,
  BookOpen
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

onMounted(() => {
  ledgerStore.fetchAccounts()
  journalsStore.fetchJournals()
})

// State
const startDate = ref('01-09-2026')
const endDate = ref('30-11-2026')
const currentLevel = ref<number>(1)
const jumpLineNumber = ref<number>(1)
const hideStockReceipt = ref(true)

export interface PLTreeNode {
  id: string
  lineNo: number
  code?: string
  name: string
  level: number
  amount: number
  isHeader: boolean
  isBold?: boolean
  accountCode?: string
  children?: PLTreeNode[]
}

// Tree Data for Profit & Loss matching Image 4 and Blueprint
const plTree = ref<PLTreeNode[]>([
  {
    id: 'sec_pendapatan',
    lineNo: 6,
    name: 'PENDAPATAN',
    level: 1,
    amount: 1544300000.00,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_pendapatan',
        lineNo: 7,
        name: 'PENDAPATAN',
        level: 2,
        amount: 1544300000.00,
        isHeader: true,
        children: [
          {
            id: '4100',
            lineNo: 8,
            code: '41',
            name: 'PENDAPATAN OPERASIONAL KLINIK',
            level: 3,
            amount: 1544300000.00,
            isHeader: true,
            accountCode: '4-1000',
            children: [
              { id: '4101', lineNo: 9, code: '4101', name: 'Pendapatan Registrasi & Skrining Anak', level: 4, amount: 24500000.00, isHeader: false, accountCode: '4-1100' },
              { id: '4102', lineNo: 10, code: '4102', name: 'Pendapatan Jasa Konsultasi Spesialis Anak (Sp.A)', level: 4, amount: 345000000.00, isHeader: false, accountCode: '4-1200' },
              { id: '4103', lineNo: 11, code: '4103', name: 'Pendapatan Tindakan Vaksinasi & Imunisasi Anak', level: 4, amount: 96000000.00, isHeader: false, accountCode: '4-1300' },
              { id: '4104', lineNo: 12, code: '4104', name: 'Pendapatan Layanan Terapi & Tumbuh Kembang Anak', level: 4, amount: 54000000.00, isHeader: false, accountCode: '4-1400' },
              { id: '4105', lineNo: 13, code: '4105', name: 'Pendapatan Penjualan Obat Resep Farmasi Kidinest', level: 4, amount: 186000000.00, isHeader: false, accountCode: '4-1500' },
              { id: '4106', lineNo: 14, code: '4106', name: 'Pendapatan Konsultasi Dokter Estetika & Kulit (Sp.DVE)', level: 4, amount: 195000000.00, isHeader: false, accountCode: '4-2100' },
              { id: '4107', lineNo: 15, code: '4107', name: 'Pendapatan Tindakan Medis Laser & Cauter SSC', level: 4, amount: 435000000.00, isHeader: false, accountCode: '4-2300' },
              { id: '4108', lineNo: 16, code: '4108', name: 'Pendapatan Tindakan Facial & Treatment Terapis', level: 4, amount: 104800000.00, isHeader: false, accountCode: '4-2200' },
              { id: '4109', lineNo: 17, code: '4109', name: 'Pendapatan Penjualan Produk Skincare & Kosmetika SSC', level: 4, amount: 104000000.00, isHeader: false, accountCode: '4-2400' }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'sec_biaya',
    lineNo: 19,
    name: 'BIAYA',
    level: 1,
    amount: 1381457250.05,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_beban',
        lineNo: 20,
        name: 'BEBAN',
        level: 2,
        amount: 1381457250.05,
        isHeader: true,
        children: [
          {
            id: '51',
            lineNo: 21,
            code: '51',
            name: 'HARGA POKOK PENJUALAN',
            level: 3,
            amount: 1139845239.18,
            isHeader: true,
            accountCode: '5-0000',
            children: [
              { id: '5101', lineNo: 22, code: '5101', name: 'HPP Obat Resep & Paten Farmasi Kidinest', level: 4, amount: 72000000.00, isHeader: false, accountCode: '5-1100' },
              { id: '5102', lineNo: 23, code: '5102', name: 'HPP Vaksin & Serum Imunisasi Kidinest', level: 4, amount: 24000000.00, isHeader: false, accountCode: '5-1200' },
              { id: '5103', lineNo: 24, code: '5103', name: 'HPP Produk Skincare & Kosmetika Terjual SSC', level: 4, amount: 49500000.00, isHeader: false, accountCode: '5-1300' },
              { id: '5104', lineNo: 25, code: '5104', name: 'HPP Bahan Tindakan Facial & Peeling SSC', level: 4, amount: 21250000.00, isHeader: false, accountCode: '5-1400' },
              { id: '5105', lineNo: 26, code: '5105', name: 'Beban Pemakaian Bahan Medis Habis Pakai (BMHP)', level: 4, amount: 29400000.00, isHeader: false, accountCode: '5-1400' },
              { id: '5106', lineNo: 27, code: '5106', name: 'Beban Honor Dokter Spesialis Anak (Sp.A)', level: 4, amount: 435000000.00, isHeader: false, accountCode: '6-1100' },
              { id: '5107', lineNo: 28, code: '5107', name: 'Beban Honor Dokter Kulit & Terapis Estetika', level: 4, amount: 483195239.18, isHeader: false, accountCode: '6-1110' },
              { id: '5108', lineNo: 29, code: '5108', name: 'Beban Komisi Tindakan Beautician / Terapis SSC', level: 4, amount: 25500000.00, isHeader: false, accountCode: '6-1110' }
            ]
          },
          {
            id: '52',
            lineNo: 51,
            code: '52',
            name: 'BEBAN OPERASIONAL',
            level: 3,
            amount: 241612010.87,
            isHeader: true,
            accountCode: '6-0000',
            children: [
              { id: '5201', lineNo: 52, code: '5201', name: 'Beban Gaji Karyawan, Perawat & Staf', level: 4, amount: 43637768.00, isHeader: false, accountCode: '6-1200' },
              { id: '5202', lineNo: 99, code: '5202', name: 'Beban Operasional (Listrik, Air, Sewa Gedung & Limbah B3)', level: 4, amount: 197974242.87, isHeader: false, accountCode: '6-1300' }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'sec_pend_lain',
    lineNo: 308,
    name: 'PENDAPATAN LAIN-LAIN',
    level: 1,
    amount: 7519886.94,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_pend_non_usaha',
        lineNo: 309,
        name: 'PENDAPATAN NON USAHA',
        level: 2,
        amount: 7519886.94,
        isHeader: true,
        children: [
          { id: '6101', lineNo: 310, code: '6101', name: 'Pendapatan Bunga Bank & Jasa Giro', level: 3, amount: 794106.94, isHeader: false },
          { id: '6106', lineNo: 324, code: '6106', name: 'Pendapatan Lain-Lain Non Usaha & Selisih Lebih Kasir', level: 3, amount: 6725780.00, isHeader: false }
        ]
      }
    ]
  },
  {
    id: 'sec_biaya_lain',
    lineNo: 328,
    name: 'BIAYA LAIN-LAIN',
    level: 1,
    amount: 63575480.42,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_beban_non_usaha',
        lineNo: 329,
        name: 'BEBAN NON USAHA',
        level: 2,
        amount: 63575480.42,
        isHeader: true,
        children: [
          { id: '7103', lineNo: 330, code: '7103', name: 'Beban Biaya Administrasi Bank', level: 3, amount: 1056900.00, isHeader: false },
          { id: '7106', lineNo: 443, code: '7106', name: 'Selisih Pembulatan & Penyesuaian Kasir', level: 3, amount: 16.00, isHeader: false },
          { id: '7107', lineNo: 447, code: '7107', name: 'Beban Bunga Pinjaman & Sewa Guna Usaha', level: 3, amount: 62359740.67, isHeader: false },
          { id: '7104', lineNo: 453, code: '7104', name: 'Pajak Jasa Giro Perbankan', level: 3, amount: 158823.75, isHeader: false }
        ]
      }
    ]
  },
  {
    id: 'laba_sebelum_pajak',
    lineNo: 468,
    name: 'LABA(RUGI) SEBELUM PAJAK',
    level: 1,
    amount: 106787156.47,
    isHeader: false,
    isBold: true
  }
])

// Dynamically sync P&L amounts with ledger trial balance
watchEffect(() => {
  if (!ledgerStore.trialBalanceRows || ledgerStore.trialBalanceRows.length === 0) return

  const tbMap = new Map<string, number>()
  for (const row of ledgerStore.trialBalanceRows) {
    const bal = row.accountType === 'REVENUE' ? (row.endingCredit - row.endingDebit) : (row.endingDebit - row.endingCredit)
    tbMap.set(row.accountCode, Math.max(bal, 0))
  }

  function updatePLNode(node: PLTreeNode): number {
    if (node.children && node.children.length > 0) {
      let sum = 0
      for (const child of node.children) {
        sum += updatePLNode(child)
      }
      node.amount = sum
      return sum
    }
    const code = node.accountCode || node.code || ''
    if (tbMap.has(code)) {
      node.amount = tbMap.get(code)!
    }
    return node.amount
  }

  for (const root of plTree.value) {
    if (root.children) {
      updatePLNode(root)
    }
  }
})

const expandedIds = ref<Set<string>>(new Set([
  'sec_pendapatan', 'grp_pendapatan',
  'sec_biaya', 'grp_beban',
  'sec_pend_lain', 'grp_pend_non_usaha',
  'sec_biaya_lain', 'grp_beban_non_usaha'
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
    expandedIds.value.add('sec_pendapatan')
    expandedIds.value.add('grp_pendapatan')
    expandedIds.value.add('sec_biaya')
    expandedIds.value.add('grp_beban')
    expandedIds.value.add('sec_pend_lain')
    expandedIds.value.add('grp_pend_non_usaha')
    expandedIds.value.add('sec_biaya_lain')
    expandedIds.value.add('grp_beban_non_usaha')
    return
  }

  function traverse(nodes: PLTreeNode[]) {
    for (const node of nodes) {
      if (node.level < lvl && node.children && node.children.length > 0) {
        expandedIds.value.add(node.id)
        traverse(node.children)
      }
    }
  }

  expandedIds.value.add('sec_pendapatan')
  expandedIds.value.add('grp_pendapatan')
  expandedIds.value.add('sec_biaya')
  expandedIds.value.add('grp_beban')
  expandedIds.value.add('sec_pend_lain')
  expandedIds.value.add('grp_pend_non_usaha')
  expandedIds.value.add('sec_biaya_lain')
  expandedIds.value.add('grp_beban_non_usaha')
  traverse(plTree.value)
}

function applySetLevel() {
  const lvl = Math.min(Math.max(Number(jumpLineNumber.value) || 1, 1), 4)
  setLevel(lvl)
}

function expandAll() {
  currentLevel.value = 4
  jumpLineNumber.value = 4
  function collectAll(nodes: PLTreeNode[]) {
    for (const n of nodes) {
      if (n.children && n.children.length > 0) {
        expandedIds.value.add(n.id)
        collectAll(n.children)
      }
    }
  }
  collectAll(plTree.value)
}

function collapseAll() {
  currentLevel.value = 1
  jumpLineNumber.value = 1
  expandedIds.value.clear()
  expandedIds.value.add('sec_pendapatan')
  expandedIds.value.add('grp_pendapatan')
  expandedIds.value.add('sec_biaya')
  expandedIds.value.add('grp_beban')
  expandedIds.value.add('sec_pend_lain')
  expandedIds.value.add('grp_pend_non_usaha')
  expandedIds.value.add('sec_biaya_lain')
  expandedIds.value.add('grp_beban_non_usaha')
}

interface FlattenedPLRow {
  node: PLTreeNode
  visible: boolean
  hasChildren: boolean
  isExpanded: boolean
}

const flattenedRows = computed<FlattenedPLRow[]>(() => {
  const result: FlattenedPLRow[] = []

  function flatten(nodes: PLTreeNode[], isParentExpanded: boolean) {
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

  flatten(plTree.value, true)
  return result
})

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/accounting/ledger', query: { account: accountCode } })
  } else {
    router.push('/accounting/ledger')
  }
}

function handleExportExcel() {
  exportProfitLossToExcel(plTree.value, {
    businessUnit: accountingStore.businessUnit,
    asOfDate: endDate.value,
    startDate: startDate.value
  })
  accountingStore.showToast('Laporan Laba Rugi Komprehensif (Format LHP SPHC) berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Header matching Clinical Light Design -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl p-4 sm:p-5 text-zinc-900 flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-2xs">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-zinc-400 font-medium">
          <span class="hover:text-zinc-700 cursor-pointer">Accounting</span>
          <span>/</span>
          <span class="text-zinc-900 font-bold">Laporan Laba Rugi</span>
        </div>
        <h1 class="text-base sm:text-lg font-bold text-zinc-950 tracking-tight">
          {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SEMARANG SKIN CENTER)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER (SSC)' }}
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-2.5 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-zinc-50 px-3 py-1.5 rounded-xl border border-zinc-200 text-zinc-700 font-medium">
          <span class="text-zinc-500">Periode:</span>
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
          <span class="text-zinc-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor</span>
        </button>

        <button
          @click="printReport"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Checkbox bar -->
    <div class="flex items-center gap-2 text-xs text-zinc-600 px-1">
      <input
        v-model="hideStockReceipt"
        type="checkbox"
        id="hideStock"
        class="rounded bg-white border-zinc-300 text-zinc-900 focus:ring-0 cursor-pointer"
      />
      <label for="hideStock" class="cursor-pointer select-none font-medium">
        Sembunyikan Penerimaan Stok (Material Receipt)
      </label>
    </div>

    <!-- Tree Table Container (Clean Light Theme) -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl overflow-hidden shadow-2xs text-zinc-800">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-zinc-200 bg-zinc-50/90 text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              <th class="py-2.5 px-4 w-16 text-zinc-400">#</th>
              <th class="py-2.5 px-4">Uraian</th>
              <th class="py-2.5 px-6 text-right w-64">Jumlah</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 font-mono">
            <tr
              v-for="row in flattenedRows"
              :key="row.node.id"
              :class="[
                'transition-colors group',
                row.node.level === 1
                  ? 'bg-zinc-100/90 font-extrabold text-zinc-950 text-sm border-b border-zinc-200'
                  : row.node.level === 2
                  ? 'bg-zinc-50/70 font-bold text-zinc-900 text-xs border-b border-zinc-100'
                  : 'hover:bg-zinc-50/80 text-zinc-700'
              ]"
            >
              <!-- Line Number -->
              <td class="py-2 px-4 text-zinc-400 text-[11px] select-none font-mono">
                {{ row.node.lineNo }}
              </td>

              <!-- Uraian with Caret & Book Icon -->
              <td class="py-2 px-4 font-sans">
                <div
                  class="flex items-center gap-1.5"
                  :style="{ paddingLeft: `${(row.node.level - 1) * 20}px` }"
                >
                  <!-- Caret for Parent Nodes -->
                  <button
                    v-if="row.hasChildren"
                    @click="toggleExpand(row.node.id)"
                    class="p-0.5 rounded text-zinc-400 hover:text-zinc-800 hover:bg-zinc-200/60 transition-colors cursor-pointer"
                  >
                    <ChevronDown v-if="row.isExpanded" class="w-3.5 h-3.5 text-zinc-600" />
                    <ChevronRight v-else class="w-3.5 h-3.5 text-zinc-400" />
                  </button>
                  <span v-else class="w-4 inline-block"></span>

                  <!-- Account Code and Name (blue link style) -->
                  <span
                    :class="[
                      'tracking-wide',
                      row.node.level === 1 ? 'font-black tracking-wider text-zinc-950 uppercase' : '',
                      row.node.level === 2 ? 'font-bold text-zinc-900 uppercase' : '',
                      row.node.level === 3 ? 'font-semibold text-zinc-800' : '',
                      row.node.level === 4 ? 'font-semibold text-blue-600 hover:underline cursor-pointer' : ''
                    ]"
                    @click="row.node.accountCode && openLedgerForAccount(row.node.accountCode)"
                  >
                    {{ row.node.code ? `${row.node.code} - ` : '' }}{{ row.node.name }}
                  </span>

                  <!-- Subsidiary Ledger Book Icon (📖) -->
                  <button
                    v-if="row.node.accountCode"
                    @click.stop="openLedgerForAccount(row.node.accountCode)"
                    class="ml-1 text-zinc-400 hover:text-blue-600 transition-colors p-0.5 rounded cursor-pointer"
                    title="Buka Buku Besar untuk akun ini"
                  >
                    <BookOpen class="w-3 h-3 text-zinc-400 group-hover:text-blue-600" />
                  </button>
                </div>
              </td>

              <!-- Jumlah Column -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight',
                  row.node.level === 1
                    ? 'font-black text-sm text-zinc-950'
                    : row.node.level === 2
                    ? 'font-bold text-zinc-900 text-xs'
                    : 'font-medium text-zinc-800'
                ]"
              >
                {{ formatIDR(row.node.amount) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Bottom Level Bar (Clean Light Theme) -->
      <div class="px-4 py-3 bg-zinc-50/80 border-t border-zinc-200 flex flex-wrap items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-2">
          <!-- Jump to Line / Set Level Input -->
          <input
            v-model.number="jumpLineNumber"
            type="number"
            min="1"
            max="4"
            class="w-12 px-2 py-1 bg-white border border-zinc-200 rounded-lg text-center text-xs text-zinc-900 font-mono font-bold shadow-2xs"
            @keyup.enter="applySetLevel"
          />

          <!-- Level Buttons: Lvl 1, Lvl 2, Lvl 3, Lvl 4 (Detail) -->
          <div class="flex items-center bg-zinc-200/70 p-0.5 rounded-xl border border-zinc-300/80">
            <button
              @click="setLevel(1)"
              :class="[
                'px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer',
                currentLevel === 1 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'
              ]"
            >
              Lvl 1
            </button>
            <button
              @click="setLevel(2)"
              :class="[
                'px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer',
                currentLevel === 2 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'
              ]"
            >
              Lvl 2
            </button>
            <button
              @click="setLevel(3)"
              :class="[
                'px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer',
                currentLevel === 3 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'
              ]"
            >
              Lvl 3
            </button>
            <button
              @click="setLevel(4)"
              :class="[
                'px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer',
                currentLevel === 4 ? 'bg-zinc-950 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'
              ]"
            >
              Lvl 4 (Detail)
            </button>
          </div>

          <button
            @click="applySetLevel"
            class="px-3 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-900 rounded-xl font-semibold border border-zinc-200 shadow-2xs transition-colors cursor-pointer"
          >
            Set Level
          </button>

          <button
            @click="expandAll"
            class="px-3 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-900 rounded-xl font-semibold border border-zinc-200 shadow-2xs transition-colors cursor-pointer"
          >
            Expand All
          </button>

          <button
            @click="collapseAll"
            class="px-3 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-900 rounded-xl font-semibold border border-zinc-200 shadow-2xs transition-colors cursor-pointer"
          >
            Collapse All
          </button>
        </div>

        <div class="flex items-center gap-3 text-zinc-500 text-xs">
          <span>Menampilkan <strong>{{ flattenedRows.length }}</strong> pos laba rugi</span>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-emerald-700 font-bold bg-emerald-50 px-2 py-0.5 rounded-md border border-emerald-200">Laba Bersih: {{ formatIDR(106787156.47) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
