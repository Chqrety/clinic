<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { accountingApi } from '@/api/accounting'
import { formatIDR } from '@/utils/currency'
import { exportCashFlowToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  BookOpen,
  Loader2
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()

// State
const activeTab = ref<'depreciation' | 'cashflow'>('cashflow')
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')
const isLoading = ref(false)
const netCashChange = ref<number>(0)

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

const cashFlowRows = ref<CashFlowRow[]>([])

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

function buildEmptyRows(): CashFlowRow[] {
  return [
    { lineNo: 1, name: 'ARUS KAS DARI AKTIVITAS OPERASIONAL:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
    { lineNo: 2, name: 'Penerimaan Kas Layanan & Piutang Pasien', amountCurrent: 0, amountPrevious: 0 },
    { lineNo: 3, name: 'Pengeluaran Kas Operasional, Obat & Dokter', amountCurrent: 0, amountPrevious: 0 },
    { lineNo: 4, name: 'Arus Kas Bersih yang Diperoleh dari Aktivitas Operasional', amountCurrent: 0, amountPrevious: 0, isBold: true },
    { lineNo: 6, name: 'ARUS KAS DARI AKTIVITAS INVESTASI:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
    { lineNo: 7, name: 'Perolehan Mesin Medis & Aset Tetap', amountCurrent: 0, amountPrevious: 0 },
    { lineNo: 8, name: 'Kas Bersih Digunakan untuk Aktivitas Investasi', amountCurrent: 0, amountPrevious: 0, isBold: true },
    { lineNo: 10, name: 'ARUS KAS DARI AKTIVITAS PENDANAAN:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
    { lineNo: 11, name: 'Setoran Modal & Pembiayaan', amountCurrent: 0, amountPrevious: 0 },
    { lineNo: 12, name: 'Kas Bersih Diperoleh dari Aktivitas Pendanaan', amountCurrent: 0, amountPrevious: 0, isBold: true },
    { lineNo: 14, name: 'KENAIKAN (PENURUNAN) BERSIH KAS DAN SETARA KAS', amountCurrent: 0, amountPrevious: 0, isBold: true },
    { lineNo: 15, name: 'KAS DAN SETARA KAS PADA AWAL PERIODE', amountCurrent: 0, amountPrevious: 0, isBold: true },
    { lineNo: 16, name: 'KAS DAN SETARA KAS PADA AKHIR PERIODE (KAS & BANK)', amountCurrent: 0, amountPrevious: 0, isBold: true }
  ]
}

async function fetchCashFlow() {
  isLoading.value = true
  try {
    const sDate = parseToISODate(startDate.value)
    const eDate = parseToISODate(endDate.value)
    const data = await accountingApi.getCashFlowStatement(sDate, eDate)

    if (!data) {
      cashFlowRows.value = buildEmptyRows()
      netCashChange.value = 0
      return
    }

    const rows: CashFlowRow[] = []
    let line = 1

    // 1. Operating Activities
    rows.push({
      lineNo: line++,
      name: 'ARUS KAS DARI AKTIVITAS OPERASIONAL:',
      amountCurrent: 0,
      amountPrevious: 0,
      isSectionHeader: true,
      isBold: true
    })

    const opItems = data.operating_activities?.items || []
    if (opItems.length > 0) {
      for (const item of opItems) {
        rows.push({
          lineNo: line++,
          name: item.account_name,
          amountCurrent: Number(item.amount) || 0,
          amountPrevious: 0,
          isLink: true,
          accountCode: item.account_code
        })
      }
    } else {
      rows.push({
        lineNo: line++,
        name: 'Penerimaan Kas Layanan & Piutang Pasien',
        amountCurrent: 0,
        amountPrevious: 0
      })
      rows.push({
        lineNo: line++,
        name: 'Pengeluaran Kas Operasional, Obat & Dokter',
        amountCurrent: 0,
        amountPrevious: 0
      })
    }

    rows.push({
      lineNo: line++,
      name: 'Arus Kas Bersih yang Diperoleh dari Aktivitas Operasional',
      amountCurrent: Number(data.net_operating_cash) || 0,
      amountPrevious: 0,
      isBold: true
    })

    // 2. Investing Activities
    line++
    rows.push({
      lineNo: line++,
      name: 'ARUS KAS DARI AKTIVITAS INVESTASI:',
      amountCurrent: 0,
      amountPrevious: 0,
      isSectionHeader: true,
      isBold: true
    })

    const invItems = data.investing_activities?.items || []
    if (invItems.length > 0) {
      for (const item of invItems) {
        rows.push({
          lineNo: line++,
          name: item.account_name,
          amountCurrent: Number(item.amount) || 0,
          amountPrevious: 0,
          isLink: true,
          accountCode: item.account_code
        })
      }
    } else {
      rows.push({
        lineNo: line++,
        name: 'Perolehan Mesin Medis & Aset Tetap',
        amountCurrent: 0,
        amountPrevious: 0
      })
    }

    rows.push({
      lineNo: line++,
      name: 'Kas Bersih Digunakan untuk Aktivitas Investasi',
      amountCurrent: Number(data.net_investing_cash) || 0,
      amountPrevious: 0,
      isBold: true
    })

    // 3. Financing Activities
    line++
    rows.push({
      lineNo: line++,
      name: 'ARUS KAS DARI AKTIVITAS PENDANAAN:',
      amountCurrent: 0,
      amountPrevious: 0,
      isSectionHeader: true,
      isBold: true
    })

    const finItems = data.financing_activities?.items || []
    if (finItems.length > 0) {
      for (const item of finItems) {
        rows.push({
          lineNo: line++,
          name: item.account_name,
          amountCurrent: Number(item.amount) || 0,
          amountPrevious: 0,
          isLink: true,
          accountCode: item.account_code
        })
      }
    } else {
      rows.push({
        lineNo: line++,
        name: 'Setoran Modal & Pembiayaan',
        amountCurrent: 0,
        amountPrevious: 0
      })
    }

    rows.push({
      lineNo: line++,
      name: 'Kas Bersih Diperoleh dari Aktivitas Pendanaan',
      amountCurrent: Number(data.net_financing_cash) || 0,
      amountPrevious: 0,
      isBold: true
    })

    // 4. Summaries
    line++
    rows.push({
      lineNo: line++,
      name: 'KENAIKAN (PENURUNAN) BERSIH KAS DAN SETARA KAS',
      amountCurrent: Number(data.net_cash_change) || 0,
      amountPrevious: 0,
      isBold: true
    })
    rows.push({
      lineNo: line++,
      name: 'KAS DAN SETARA KAS PADA AWAL PERIODE',
      amountCurrent: Number(data.beginning_cash) || 0,
      amountPrevious: 0,
      isBold: true
    })
    rows.push({
      lineNo: line++,
      name: 'KAS DAN SETARA KAS PADA AKHIR PERIODE (KAS & BANK)',
      amountCurrent: Number(data.ending_cash) || 0,
      amountPrevious: 0,
      isBold: true
    })

    cashFlowRows.value = rows
    netCashChange.value = Number(data.net_cash_change) || 0
  } catch (err) {
    console.warn('[CashFlowView] Gagal mengambil arus kas dari backend:', err)
    cashFlowRows.value = buildEmptyRows()
    netCashChange.value = 0
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchCashFlow()
})

watch([startDate, endDate, () => accountingStore.businessUnit], () => {
  fetchCashFlow()
})

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
  accountingStore.showToast('Laporan Arus Kas berhasil diekspor ke Excel', 'success')
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
          Kenaikan Kas Bersih: {{ formatIDR(netCashChange) }}
        </span>
      </div>
    </div>
  </div>
</template>
