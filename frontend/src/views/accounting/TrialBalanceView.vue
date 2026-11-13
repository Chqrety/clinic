<script setup lang="ts">
import { computed } from 'vue'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel } from '@/utils/export'
import {
  Scale,
  Download,
  CheckCircle2,
  AlertCircle
} from 'lucide-vue-next'

const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

const rows = computed(() => ledgerStore.trialBalanceRows)
const totals = computed(() => ledgerStore.trialBalanceTotals)

const isBalanced = computed(() => {
  return totals.value.endingDebit === totals.value.endingCredit
})

function handleExportExcel() {
  const data = rows.value.map(r => ({
    'Kode Akun': r.accountCode,
    'Nama Akun': r.accountName,
    'Tipe': r.accountType,
    'Saldo Awal Debet': r.openingDebit,
    'Saldo Awal Kredit': r.openingCredit,
    'Mutasi Debet': r.periodDebit,
    'Mutasi Kredit': r.periodCredit,
    'Saldo Akhir Debet': r.endingDebit,
    'Saldo Akhir Kredit': r.endingCredit
  }))

  exportToExcel(data, `Neraca_Saldo_Trial_Balance_${accountingStore.businessUnit}_2026-03`, 'Trial Balance')
  accountingStore.showToast('Neraca Saldo berhasil diekspor ke Excel (.xlsx)', 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Financial Health Verification
          </span>
          <span class="text-xs text-slate-400">&bull; Neraca Percobaan</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Neraca Saldo (Trial Balance)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Memverifikasi keseimbangan matematis buku besar dari saldo awal, mutasi jurnal periode berjalan, hingga saldo akhir seluruh akun.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-semibold border border-slate-300 shadow-2xs transition-colors"
        >
          <Download class="w-4 h-4 text-slate-600" />
          <span>Ekspor Neraca Saldo</span>
        </button>
      </div>
    </div>

    <!-- Balance Verification Banner -->
    <div
      v-if="isBalanced"
      class="p-4 bg-emerald-50 border border-emerald-200 rounded-2xl flex items-center justify-between text-emerald-900 text-xs"
    >
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-emerald-100 flex items-center justify-center text-emerald-600">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <div>
          <h4 class="font-bold text-sm">Neraca Saldo Seimbang (Trial Balance is in Balance)</h4>
          <p class="text-emerald-700 text-[11px]">
            Total Debet dan Total Kredit pada Saldo Akhir bernilai sama persis. Siap untuk proses penyusunan Laporan Keuangan.
          </p>
        </div>
      </div>
      <div class="text-right font-mono font-bold text-sm">
        {{ formatIDR(totals.endingDebit) }}
      </div>
    </div>
    <div
      v-else
      class="p-4 bg-rose-50 border border-rose-200 rounded-2xl flex items-center justify-between text-rose-900 text-xs"
    >
      <div class="flex items-center gap-3">
        <AlertCircle class="w-6 h-6 text-rose-600" />
        <div>
          <h4 class="font-bold text-sm">Terdapat Selisih Neraca Saldo!</h4>
          <p class="text-rose-700 text-[11px]">
            Selisih: {{ formatIDR(Math.abs(totals.endingDebit - totals.endingCredit)) }}. Periksa jurnal penyesuaian yang belum seimbang.
          </p>
        </div>
      </div>
    </div>

    <!-- 6-Column Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-100/80 border-b border-slate-300 text-[11px] font-bold text-slate-700 uppercase">
              <th rowspan="2" class="px-4 py-2 border-r border-slate-200 w-28">Kode Akun</th>
              <th rowspan="2" class="px-4 py-2 border-r border-slate-200">Nama Akun (COA)</th>
              <th colspan="2" class="px-4 py-2 text-center border-r border-slate-200 bg-slate-50">Saldo Awal</th>
              <th colspan="2" class="px-4 py-2 text-center border-r border-slate-200 bg-blue-50/50">Mutasi Periode</th>
              <th colspan="2" class="px-4 py-2 text-center bg-emerald-50/50">Saldo Akhir</th>
            </tr>
            <tr class="bg-slate-50 border-b border-slate-300 text-[10px] font-bold text-slate-600 uppercase">
              <th class="px-3 py-1.5 text-right w-28 border-r border-slate-200">Debet</th>
              <th class="px-3 py-1.5 text-right w-28 border-r border-slate-200">Kredit</th>
              <th class="px-3 py-1.5 text-right w-28 border-r border-slate-200">Debet</th>
              <th class="px-3 py-1.5 text-right w-28 border-r border-slate-200">Kredit</th>
              <th class="px-3 py-1.5 text-right w-32 border-r border-slate-200">Debet</th>
              <th class="px-3 py-1.5 text-right w-32">Kredit</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="row in rows"
              :key="row.accountCode"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-2.5 font-mono font-bold text-slate-800 border-r border-slate-100">
                {{ row.accountCode }}
              </td>
              <td class="px-4 py-2.5 font-medium text-slate-900 border-r border-slate-100">
                {{ row.accountName }}
                <span class="text-[10px] text-slate-400 block">{{ row.accountType }}</span>
              </td>
              <!-- Opening -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-slate-100 text-slate-700">
                {{ row.openingDebit > 0 ? formatIDR(row.openingDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono border-r border-slate-100 text-slate-700">
                {{ row.openingCredit > 0 ? formatIDR(row.openingCredit, false) : '-' }}
              </td>
              <!-- Movement -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-slate-100 text-emerald-700 font-semibold">
                {{ row.periodDebit > 0 ? formatIDR(row.periodDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono border-r border-slate-100 text-rose-700 font-semibold">
                {{ row.periodCredit > 0 ? formatIDR(row.periodCredit, false) : '-' }}
              </td>
              <!-- Ending -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-slate-100 font-bold text-slate-900 bg-slate-50/30">
                {{ row.endingDebit > 0 ? formatIDR(row.endingDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono font-bold text-slate-900 bg-slate-50/30">
                {{ row.endingCredit > 0 ? formatIDR(row.endingCredit, false) : '-' }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-slate-100 font-bold border-t-2 border-slate-400">
              <td colspan="2" class="px-4 py-3 text-right text-slate-800 uppercase">
                TOTAL KESEIMBANGAN:
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-slate-200">
                {{ formatIDR(totals.openingDebit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-slate-200">
                {{ formatIDR(totals.openingCredit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-slate-200 text-emerald-800">
                {{ formatIDR(totals.periodDebit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-slate-200 text-rose-800">
                {{ formatIDR(totals.periodCredit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono text-sm border-r border-slate-200 text-slate-900">
                {{ formatIDR(totals.endingDebit) }}
              </td>
              <td class="px-3 py-3 text-right font-mono text-sm text-slate-900">
                {{ formatIDR(totals.endingCredit) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

