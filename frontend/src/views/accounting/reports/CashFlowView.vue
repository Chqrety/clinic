<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportCashFlowToExcel, exportToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  BookOpen,
  Layers,
  Calculator
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()

// State
const activeTab = ref<'depreciation' | 'cashflow'>('cashflow')
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')

interface CashFlowRow {
  lineNo: number
  name: string
  amountCurrent: number
  amountPrevious: number
  isSectionHeader?: boolean
  isBold?: boolean
  isLink?: boolean
  accountCode?: string
}

// Cash Flow Statement matching Image 3 and Clinic ERP Context
const cashFlowRows = ref<CashFlowRow[]>([
  { lineNo: 9, name: 'Penyesuaian untuk:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 10, name: 'Penyusutan dan Amortisasi', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '1606' },
  { lineNo: 11, name: 'Laba Operasi Sebelum Perubahan Modal Kerja', amountCurrent: 0.00, amountPrevious: 0.00, isBold: true },
  { lineNo: 13, name: 'Perubahan Modal Kerja:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 14, name: 'Penurunan (Kenaikan) Persediaan', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '13' },
  { lineNo: 15, name: 'Penurunan (Kenaikan) Piutang Usaha', amountCurrent: -411190000.00, amountPrevious: 0.00, isLink: true, accountCode: '12' },
  { lineNo: 16, name: 'Penurunan (Kenaikan) Biaya Dibayar di Muka', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '15' },
  { lineNo: 17, name: 'Penurunan (Kenaikan) Investasi', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '14' },
  { lineNo: 18, name: 'Kenaikan (Penurunan) Hutang Uang Muka', amountCurrent: 517451500.00, amountPrevious: 0.00, isLink: true, accountCode: '2102' },
  { lineNo: 19, name: 'Kenaikan (Penurunan) Hutang Dagang', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '22' },
  { lineNo: 20, name: 'Kenaikan (Penurunan) Hutang Bank', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '23' },
  { lineNo: 21, name: 'Kenaikan (Penurunan) Hutang Pembebanan HPP', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '26' },
  { lineNo: 22, name: 'Kenaikan (Penurunan) Hutang Pajak', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '25' },
  { lineNo: 23, name: 'Kas Bersih dari Aktivitas Operasi', amountCurrent: 106261500.00, amountPrevious: 0.00, isBold: true },

  { lineNo: 25, name: 'ARUS KAS DARI AKTIVITAS INVESTASI:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 26, name: 'Pembelian Aset Tetap', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '16' },
  { lineNo: 27, name: 'Kas Bersih dari Aktivitas Investasi', amountCurrent: 0.00, amountPrevious: 0.00, isBold: true },

  { lineNo: 29, name: 'ARUS KAS DARI AKTIVITAS PENDANAAN:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 30, name: 'Kenaikan (Penurunan) Modal', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '31' },
  { lineNo: 31, name: 'Kenaikan (Penurunan) Hutang Pihak Ketiga', amountCurrent: 0.00, amountPrevious: 0.00, isLink: true, accountCode: '24' },
  { lineNo: 32, name: 'Kas Bersih dari Aktivitas Pendanaan', amountCurrent: 0.00, amountPrevious: 0.00, isBold: true },

  { lineNo: 34, name: 'KENAIKAN (PENURUNAN) KAS DAN SETARA KAS', amountCurrent: 106261500.00, amountPrevious: 0.00, isBold: true },

  { lineNo: 36, name: 'PERUBAHAN DALAM KOMPONEN KAS DAN SETARA KAS:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 37, name: 'Kas dan Setara Kas pada awal periode', amountCurrent: 1138664487.55, amountPrevious: 0.00, isLink: true, accountCode: '11' },
  { lineNo: 38, name: 'Kas dan Setara Kas pada akhir periode', amountCurrent: 1244925987.55, amountPrevious: 0.00, isLink: true, accountCode: '11' },
  { lineNo: 39, name: 'Kenaikan (Penurunan) Kas', amountCurrent: 106261500.00, amountPrevious: 0.00, isBold: true }
])

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/accounting/ledger', query: { account: accountCode } })
  } else {
    router.push('/accounting/ledger')
  }
}

function handleExportExcel() {
  exportCashFlowToExcel(cashFlowRows.value, {
    businessUnit: accountingStore.businessUnit,
    asOfDate: endDate.value,
    startDate: startDate.value
  })
  accountingStore.showToast('Laporan Arus Kas Audited (Format LHP SPHC) berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Top Header Tabs matching Clinical Light Design -->
    <div class="flex items-center gap-2 border-b border-zinc-200 pb-2">
      <button
        @click="activeTab = 'depreciation'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-xl transition-all cursor-pointer',
          activeTab === 'depreciation'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'
        ]"
      >
        Penyusutan dan Amortisasi
      </button>
      <button
        @click="activeTab = 'cashflow'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-xl transition-all cursor-pointer',
          activeTab === 'cashflow'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'
        ]"
      >
        Arus Kas
      </button>
    </div>

    <!-- Header matching Clinical Light Design -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl p-4 sm:p-5 text-zinc-900 flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-2xs">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-zinc-400 font-medium">
          <span class="hover:text-zinc-700 cursor-pointer">Accounting</span>
          <span>/</span>
          <span class="text-zinc-900 font-bold">Laporan Arus Kas</span>
        </div>
        <h1 class="text-base sm:text-lg font-bold text-zinc-950 tracking-tight">
          {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SSC)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER (SSC)' }}
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

    <!-- Tree Table Container (Clean Light Theme) -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl overflow-hidden shadow-2xs text-zinc-800">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-zinc-200 bg-zinc-50/90 text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              <th class="py-2.5 px-4 w-16 text-zinc-400">#</th>
              <th class="py-2.5 px-4">Uraian</th>
              <th class="py-2.5 px-6 text-right w-56">November 2026</th>
              <th class="py-2.5 px-6 text-right w-56 text-zinc-400">November 2025</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 font-mono">
            <tr
              v-for="row in cashFlowRows"
              :key="row.lineNo"
              :class="[
                'transition-colors group',
                row.isSectionHeader
                  ? 'bg-zinc-100/90 font-extrabold text-zinc-950 text-xs uppercase tracking-wide border-b border-zinc-200'
                  : row.isBold
                  ? 'bg-zinc-50/80 font-bold text-zinc-900 border-b border-zinc-100'
                  : 'hover:bg-zinc-50/80 text-zinc-700'
              ]"
            >
              <!-- Line Number -->
              <td class="py-2 px-4 text-zinc-400 text-[11px] select-none font-mono">
                {{ row.lineNo }}
              </td>

              <!-- Uraian -->
              <td class="py-2 px-4 font-sans">
                <div class="flex items-center gap-1.5" :class="{ 'pl-4': !row.isSectionHeader && !row.isBold }">
                  <!-- Clickable Blue Link Style -->
                  <span
                    v-if="row.isLink"
                    @click="row.accountCode && openLedgerForAccount(row.accountCode)"
                    class="text-blue-600 hover:underline cursor-pointer font-medium flex items-center gap-1"
                  >
                    <span>{{ row.name }}</span>
                    <BookOpen v-if="row.accountCode" class="w-3 h-3 text-zinc-400 group-hover:text-blue-600 inline" />
                  </span>
                  <span
                    v-else
                    :class="[
                      row.isSectionHeader ? 'font-black uppercase tracking-wider text-zinc-950' : '',
                      row.isBold && !row.isSectionHeader ? 'font-bold text-zinc-900' : ''
                    ]"
                  >
                    {{ row.name }}
                  </span>
                </div>
              </td>

              <!-- November 2026 -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight',
                  row.isSectionHeader ? 'text-transparent' : '',
                  row.isBold ? 'font-bold text-zinc-950' : 'font-medium text-zinc-800'
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountCurrent) }}
              </td>

              <!-- November 2025 -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight text-zinc-400',
                  row.isSectionHeader ? 'text-transparent' : ''
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountPrevious) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Bottom Status Note (Clean Light Theme) -->
      <div class="px-4 py-3 bg-zinc-50/80 border-t border-zinc-200 flex items-center justify-between text-[11px] text-zinc-500">
        <span class="font-mono">
          Perbandingan arus kas tahun berjalan vs periode sebelumnya.
        </span>
        <span class="text-emerald-700 font-bold font-mono bg-emerald-50 px-2 py-0.5 rounded-md border border-emerald-200">
          Kenaikan Kas Bersih: {{ formatIDR(106261500) }}
        </span>
      </div>
    </div>
  </div>
</template>
