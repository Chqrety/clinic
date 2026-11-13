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
  BookOpen
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()

// State
const startDate = ref('01-09-2026')
const endDate = ref('30-09-2026')
const currentLevel = ref<number>(4)
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
              { id: '4101', lineNo: 9, code: '4101', name: 'Pendapatan Registrasi & Skrining Anak', level: 4, amount: 24500000.00, isHeader: false, accountCode: '4-1110' },
              { id: '4102', lineNo: 10, code: '4102', name: 'Pendapatan Jasa Konsultasi Spesialis Anak (Sp.A)', level: 4, amount: 345000000.00, isHeader: false, accountCode: '4-1120' },
              { id: '4103', lineNo: 11, code: '4103', name: 'Pendapatan Tindakan Vaksinasi & Imunisasi Anak', level: 4, amount: 96000000.00, isHeader: false, accountCode: '4-1130' },
              { id: '4104', lineNo: 12, code: '4104', name: 'Pendapatan Layanan Terapi & Tumbuh Kembang Anak', level: 4, amount: 54000000.00, isHeader: false, accountCode: '4-1140' },
              { id: '4105', lineNo: 13, code: '4105', name: 'Pendapatan Penjualan Obat Resep Farmasi Kidinest', level: 4, amount: 186000000.00, isHeader: false, accountCode: '4-1150' },
              { id: '4106', lineNo: 14, code: '4106', name: 'Pendapatan Konsultasi Dokter Estetika & Kulit (Sp.DVE)', level: 4, amount: 195000000.00, isHeader: false, accountCode: '4-2110' },
              { id: '4107', lineNo: 15, code: '4107', name: 'Pendapatan Tindakan Medis Laser & Cauter SSC', level: 4, amount: 435000000.00, isHeader: false, accountCode: '4-2120' },
              { id: '4108', lineNo: 16, code: '4108', name: 'Pendapatan Tindakan Facial & Treatment Terapis', level: 4, amount: 104800000.00, isHeader: false, accountCode: '4-2130' },
              { id: '4109', lineNo: 17, code: '4109', name: 'Pendapatan Penjualan Produk Skincare & Kosmetika SSC', level: 4, amount: 104000000.00, isHeader: false, accountCode: '4-2140' }
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
              { id: '5101', lineNo: 22, code: '5101', name: 'HPP Obat Resep & Paten Farmasi Kidinest', level: 4, amount: 72000000.00, isHeader: false, accountCode: '5-1101' },
              { id: '5102', lineNo: 23, code: '5102', name: 'HPP Vaksin & Serum Imunisasi Kidinest', level: 4, amount: 24000000.00, isHeader: false, accountCode: '5-1102' },
              { id: '5103', lineNo: 24, code: '5103', name: 'HPP Produk Skincare & Kosmetika Terjual SSC', level: 4, amount: 49500000.00, isHeader: false, accountCode: '5-1103' },
              { id: '5104', lineNo: 25, code: '5104', name: 'HPP Bahan Tindakan Facial & Peeling SSC', level: 4, amount: 21250000.00, isHeader: false },
              { id: '5105', lineNo: 26, code: '5105', name: 'Beban Pemakaian Bahan Medis Habis Pakai (BMHP)', level: 4, amount: 29400000.00, isHeader: false, accountCode: '5-1105' },
              { id: '5106', lineNo: 27, code: '5106', name: 'Beban Bagi Hasil Dokter Spesialis Anak (Kidinest)', level: 4, amount: 435000000.00, isHeader: false, accountCode: '5-1106' },
              { id: '5107', lineNo: 28, code: '5107', name: 'Beban Bagi Hasil Dokter Estetika & Kulit SSC', level: 4, amount: 483195239.18, isHeader: false, accountCode: '5-1107' },
              { id: '5108', lineNo: 29, code: '5108', name: 'Beban Komisi Tindakan Beautician / Terapis SSC', level: 4, amount: 25500000.00, isHeader: false, accountCode: '5-1108' }
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
              { id: '5201', lineNo: 52, code: '5201', name: 'Biaya Marketing, Iklan Meta & Edukasi Medis', level: 4, amount: 43637768.00, isHeader: false },
              { id: '5202', lineNo: 99, code: '5202', name: 'Beban Umum dan Administrasi (Gaji, Listrik, Sewa, Limbah B3, Servis Laser)', level: 4, amount: 197974242.87, isHeader: false, accountCode: '6-1101' }
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

const expandedIds = ref<Set<string>>(new Set([
  'sec_pendapatan', 'grp_pendapatan', '4100',
  'sec_biaya', 'grp_beban', '51', '52',
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
  expandedIds.value.clear()

  function traverse(nodes: PLTreeNode[]) {
    for (const node of nodes) {
      if (node.level < lvl && node.children && node.children.length > 0) {
        expandedIds.value.add(node.id)
        traverse(node.children)
      }
    }
  }

  traverse(plTree.value)
}

function expandAll() {
  currentLevel.value = 4
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
  expandedIds.value.clear()
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

  exportToExcel(data, `Laporan_Laba_Rugi_${accountingStore.businessUnit}_${endDate.value}`, 'Laba Rugi')
  accountingStore.showToast('Laporan Laba Rugi berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header matching Image 4 -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-4 text-white flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xl">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-slate-400">
          <span class="hover:text-white cursor-pointer font-medium">Accounting</span>
          <span>/</span>
          <span class="text-blue-400 font-semibold">Laba Rugi</span>
        </div>
        <h1 class="text-base font-bold text-white tracking-wide">
          {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SEMARANG SKIN CENTER)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER (SSC)' }}
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-3 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-[#1E293B] px-3 py-1.5 rounded-lg border border-slate-700">
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

    <!-- Checkbox bar -->
    <div class="flex items-center gap-2 text-xs text-slate-400 px-1">
      <input
        v-model="hideStockReceipt"
        type="checkbox"
        id="hideStock"
        class="rounded bg-[#1E293B] border-slate-700 text-blue-600 focus:ring-0"
      />
      <label for="hideStock" class="cursor-pointer select-none">
        Sembunyikan Penerimaan Stok (Material Receipt)
      </label>
    </div>

    <!-- Tree Table Container (High-Contrast Dark SaaS Look matching Image 4) -->
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
              <!-- Line Number -->
              <td class="py-2 px-4 text-slate-500 text-[11px] select-none font-mono">
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
                    class="p-0.5 rounded text-slate-400 hover:text-white hover:bg-slate-700/60 transition-colors"
                  >
                    <ChevronDown v-if="row.isExpanded" class="w-3.5 h-3.5 text-slate-300" />
                    <ChevronRight v-else class="w-3.5 h-3.5 text-slate-400" />
                  </button>
                  <span v-else class="w-4 inline-block"></span>

                  <!-- Account Code and Name (blue link style matching Image 4) -->
                  <span
                    :class="[
                      'tracking-wide',
                      row.node.level === 1 ? 'font-black tracking-wider text-white uppercase' : '',
                      row.node.level === 2 ? 'font-bold text-slate-100 uppercase' : '',
                      row.node.level === 3 ? 'font-semibold text-slate-200' : '',
                      row.node.level === 4 ? 'font-semibold text-blue-400 hover:underline cursor-pointer' : ''
                    ]"
                    @click="row.node.accountCode && openLedgerForAccount(row.node.accountCode)"
                  >
                    {{ row.node.code ? `${row.node.code} - ` : '' }}{{ row.node.name }}
                  </span>

                  <!-- Subsidiary Ledger Book Icon (📖) -->
                  <button
                    v-if="row.node.accountCode"
                    @click.stop="openLedgerForAccount(row.node.accountCode)"
                    class="ml-1 text-slate-500 hover:text-blue-400 transition-colors p-0.5 rounded"
                    title="Buka Buku Besar untuk akun ini"
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

      <!-- Bottom Level Bar matching Image 4 -->
      <div class="px-4 py-3 bg-[#0F172A] border-t border-slate-800 flex flex-wrap items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-2">
          <!-- Jump to Line -->
          <input
            v-model.number="jumpLineNumber"
            type="number"
            class="w-12 px-2 py-1 bg-[#1E293B] border border-slate-700 rounded text-center text-xs text-white font-mono"
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
            @click="expandAll"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium border border-slate-700"
          >
            Expand All
          </button>

          <button
            @click="collapseAll"
            class="px-3 py-1 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium border border-slate-700"
          >
            Collapse All
          </button>
        </div>

        <div class="flex items-center gap-3 text-slate-400 text-xs">
          <span>Menampilkan <strong>{{ flattenedRows.length }}</strong> pos laba rugi</span>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-emerald-400 font-bold">Laba Bersih: {{ formatIDR(106787156.47) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
