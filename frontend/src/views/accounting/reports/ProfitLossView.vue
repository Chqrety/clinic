<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { accountingApi } from '@/api/accounting'
import { formatIDR } from '@/utils/currency'
import { exportProfitLossToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  ChevronRight,
  ChevronDown,
  BookOpen,
  Loader2
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

// State
const startDate = ref('01-01-2026')
const endDate = ref('31-12-2026')
const currentLevel = ref<number>(1)
const jumpLineNumber = ref<number>(1)
const hideStockReceipt = ref(true)
const isLoading = ref<boolean>(false)

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

function parseToISODate(dateStr: string): string {
  if (!dateStr) return new Date().toISOString().slice(0, 10)
  const trimmed = dateStr.trim()
  if (/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) return trimmed
  const parts = trimmed.split(/[-/]/)
  if (parts.length === 3) {
    if (parts[0].length === 2 && parts[2].length === 4) {
      return `${parts[2]}-${parts[1].padStart(2, '0')}-${parts[0].padStart(2, '0')}`
    }
  }
  return trimmed
}

function buildEmptyPLTree(): PLTreeNode[] {
  return [
    {
      id: 'sec_pendapatan',
      lineNo: 1,
      name: 'PENDAPATAN',
      level: 1,
      amount: 0,
      isHeader: true,
      isBold: true,
      children: [
        {
          id: 'grp_pendapatan',
          lineNo: 2,
          name: 'PENDAPATAN OPERASIONAL KLINIK',
          level: 2,
          amount: 0,
          isHeader: true,
          children: []
        }
      ]
    },
    {
      id: 'sec_biaya',
      lineNo: 3,
      name: 'BIAYA & BEBAN OPERASIONAL',
      level: 1,
      amount: 0,
      isHeader: true,
      isBold: true,
      children: [
        {
          id: 'grp_cogs',
          lineNo: 4,
          code: '51',
          name: 'HARGA POKOK PENJUALAN (HPP)',
          level: 2,
          amount: 0,
          isHeader: true,
          children: []
        },
        {
          id: 'grp_beban',
          lineNo: 5,
          code: '52',
          name: 'BEBAN OPERASIONAL KLINIK',
          level: 2,
          amount: 0,
          isHeader: true,
          children: []
        }
      ]
    },
    {
      id: 'laba_sebelum_pajak',
      lineNo: 6,
      name: 'LABA (RUGI) BERSIH OPERASIONAL',
      level: 1,
      amount: 0,
      isHeader: false,
      isBold: true
    }
  ]
}

function buildPLTree(data: any): PLTreeNode[] {
  if (!data) return buildEmptyPLTree()

  let lineCounter = 1

  // 1. Revenue
  const revChildren: PLTreeNode[] = []
  const revCats = data.revenue_categories || []
  for (const cat of revCats) {
    const catItems: PLTreeNode[] = (cat.items || []).map((it: any) => ({
      id: `rev-${it.account_code}`,
      lineNo: ++lineCounter,
      code: it.account_code,
      name: it.account_name,
      level: 4,
      amount: Number(it.amount) || 0,
      isHeader: false,
      accountCode: it.account_code
    }))

    revChildren.push({
      id: `grp_rev_${cat.category_name}`,
      lineNo: ++lineCounter,
      name: cat.category_name,
      level: 3,
      amount: Number(cat.total) || 0,
      isHeader: true,
      children: catItems
    })
  }

  const totalRev = Number(data.total_revenue) || 0
  const secPendapatan: PLTreeNode = {
    id: 'sec_pendapatan',
    lineNo: 1,
    name: 'PENDAPATAN',
    level: 1,
    amount: totalRev,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_pendapatan',
        lineNo: 2,
        name: 'PENDAPATAN OPERASIONAL KLINIK',
        level: 2,
        amount: totalRev,
        isHeader: true,
        children: revChildren
      }
    ]
  }

  // 2. Expenses (COGS + Opex)
  const cogsChildren: PLTreeNode[] = []
  const cogsCats = data.cogs_categories || []
  for (const cat of cogsCats) {
    for (const it of (cat.items || [])) {
      cogsChildren.push({
        id: `cogs-${it.account_code}`,
        lineNo: ++lineCounter,
        code: it.account_code,
        name: it.account_name,
        level: 4,
        amount: Number(it.amount) || 0,
        isHeader: false,
        accountCode: it.account_code
      })
    }
  }

  const opexChildren: PLTreeNode[] = []
  const opexCats = data.expense_categories || []
  for (const cat of opexCats) {
    for (const it of (cat.items || [])) {
      opexChildren.push({
        id: `opex-${it.account_code}`,
        lineNo: ++lineCounter,
        code: it.account_code,
        name: it.account_name,
        level: 4,
        amount: Number(it.amount) || 0,
        isHeader: false,
        accountCode: it.account_code
      })
    }
  }

  const totalCogs = Number(data.total_cogs) || 0
  const totalExp = Number(data.total_expenses) || 0

  const secBiaya: PLTreeNode = {
    id: 'sec_biaya',
    lineNo: ++lineCounter,
    name: 'BIAYA & BEBAN OPERASIONAL',
    level: 1,
    amount: totalCogs + totalExp,
    isHeader: true,
    isBold: true,
    children: [
      {
        id: 'grp_cogs',
        lineNo: ++lineCounter,
        code: '51',
        name: 'HARGA POKOK PENJUALAN (HPP)',
        level: 2,
        amount: totalCogs,
        isHeader: true,
        children: cogsChildren
      },
      {
        id: 'grp_beban',
        lineNo: ++lineCounter,
        code: '52',
        name: 'BEBAN OPERASIONAL KLINIK',
        level: 2,
        amount: totalExp,
        isHeader: true,
        children: opexChildren
      }
    ]
  }

  // 3. Net Profit / Laba Bersih
  const netIncome = Number(data.net_income) || 0
  const secLaba: PLTreeNode = {
    id: 'laba_sebelum_pajak',
    lineNo: ++lineCounter,
    name: 'LABA (RUGI) BERSIH OPERASIONAL',
    level: 1,
    amount: netIncome,
    isHeader: false,
    isBold: true
  }

  return [secPendapatan, secBiaya, secLaba]
}

const plTree = ref<PLTreeNode[]>(buildEmptyPLTree())

async function fetchProfitLoss() {
  isLoading.value = true
  try {
    const sDate = parseToISODate(startDate.value)
    const eDate = parseToISODate(endDate.value)
    const data = await accountingApi.getIncomeStatement(sDate, eDate)
    plTree.value = buildPLTree(data)
  } catch (err) {
    console.warn('[ProfitLossView] Gagal mengambil data Laba Rugi dari backend:', err)
    plTree.value = buildEmptyPLTree()
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  ledgerStore.fetchAccounts()
  journalsStore.fetchJournals()
  fetchProfitLoss()
})

watch([startDate, endDate, () => accountingStore.businessUnit], () => {
  fetchProfitLoss()
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
