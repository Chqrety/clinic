<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
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
    router.push({ path: '/ledger', query: { account: accountCode } })
  } else {
    router.push('/ledger')
  }
}

function handleExportExcel() {
  const data = cashFlowRows.value.map(r => ({
    'Line': r.lineNo,
    'Uraian': r.name,
    'November 2026 (IDR)': r.amountCurrent,
    'November 2025 (IDR)': r.amountPrevious
  }))

  exportToExcel(data, `Laporan_Arus_Kas_${accountingStore.businessUnit}_${endDate.value}`, 'Arus Kas')
  accountingStore.showToast('Laporan Arus Kas berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Top Header Tabs matching Image 3: Penyusutan dan Amortisasi / Arus Kas -->
    <div class="flex items-center gap-2 border-b border-slate-800 pb-2">
      <button
        @click="activeTab = 'depreciation'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-lg transition-all',
          activeTab === 'depreciation'
            ? 'bg-slate-800 text-white shadow-xs'
            : 'text-slate-400 hover:text-slate-200'
        ]"
      >
        Penyusutan dan Amortisasi
      </button>
      <button
        @click="activeTab = 'cashflow'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-lg transition-all',
          activeTab === 'cashflow'
            ? 'bg-blue-600 text-white shadow-xs'
            : 'text-slate-400 hover:text-slate-200'
        ]"
      >
        Arus Kas
      </button>
    </div>

    <!-- Header matching Image 3 -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-4 text-white flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xl">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-slate-400">
          <span class="hover:text-white cursor-pointer font-medium">Accounting</span>
          <span>/</span>
          <span class="text-blue-400 font-semibold">Arus Kas</span>
        </div>
        <h1 class="text-base font-bold text-white tracking-wide">
          PT. BRAYO CIPTA SEJAHTERA (BCS) &bull; {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SSC)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER' }}
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

    <!-- Tree Table Container (High-Contrast Dark SaaS Look matching Image 3) -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl overflow-hidden shadow-2xl text-slate-200">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-slate-800 bg-[#0F172A] text-[11px] font-bold text-slate-400">
              <th class="py-2.5 px-4 w-16 text-slate-500">#</th>
              <th class="py-2.5 px-4">Uraian</th>
              <th class="py-2.5 px-6 text-right w-56">November 2026</th>
              <th class="py-2.5 px-6 text-right w-56 text-slate-400">November 2025</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60 font-mono">
            <tr
              v-for="row in cashFlowRows"
              :key="row.lineNo"
              :class="[
                'transition-colors group',
                row.isSectionHeader
                  ? 'bg-[#182234] font-extrabold text-white text-xs uppercase tracking-wide'
                  : row.isBold
                  ? 'bg-[#141D2E] font-bold text-slate-100'
                  : 'hover:bg-slate-800/40 text-slate-300'
              ]"
            >
              <!-- Line Number -->
              <td class="py-2 px-4 text-slate-500 text-[11px] select-none font-mono">
                {{ row.lineNo }}
              </td>

              <!-- Uraian -->
              <td class="py-2 px-4 font-sans">
                <div class="flex items-center gap-1.5" :class="{ 'pl-4': !row.isSectionHeader && !row.isBold }">
                  <!-- Clickable Blue Link Style matching Image 3 -->
                  <span
                    v-if="row.isLink"
                    @click="row.accountCode && openLedgerForAccount(row.accountCode)"
                    class="text-blue-400 hover:underline cursor-pointer font-medium flex items-center gap-1"
                  >
                    <span>{{ row.name }}</span>
                    <BookOpen v-if="row.accountCode" class="w-3 h-3 text-slate-500 group-hover:text-blue-400 inline" />
                  </span>
                  <span
                    v-else
                    :class="[
                      row.isSectionHeader ? 'font-bold uppercase tracking-wider text-slate-100' : '',
                      row.isBold && !row.isSectionHeader ? 'font-bold text-white' : ''
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
                  row.isBold ? 'font-bold text-white' : 'font-medium text-slate-300'
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountCurrent) }}
              </td>

              <!-- November 2025 -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight text-slate-400',
                  row.isSectionHeader ? 'text-transparent' : ''
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountPrevious) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Bottom Status Note matching Image 3 -->
      <div class="px-4 py-3 bg-[#0F172A] border-t border-slate-800 flex items-center justify-between text-[11px] text-slate-400">
        <span class="font-mono">
          For comparison, use &gt;5, &lt;10 or =324. For ranges, use 5:10 (for values between 5 &amp; 10).
        </span>
        <span class="text-emerald-400 font-bold font-mono">
          Kenaikan Kas Bersih: {{ formatIDR(106261500) }}
        </span>
      </div>
    </div>

    <!-- Official Report Signature Section matching Image 1 -->
    <ReportSignatures />
  </div>
</template>
