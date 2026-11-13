<script setup lang="ts">
import { ref, computed } from 'vue'
import { useReportsStore } from '@/stores/reports'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel } from '@/utils/export'
import {
  History,
  Download,
  CreditCard,
  Receipt,
  AlertCircle,
  CheckCircle2
} from 'lucide-vue-next'

const reportsStore = useReportsStore()
const accountingStore = useAccountingStore()

const activeTab = ref<'AR' | 'AP'>('AR')

const currentList = computed(() => {
  return activeTab.value === 'AR' ? reportsStore.arAgingList : reportsStore.apAgingList
})

const totalAmount = computed(() => {
  return currentList.value.reduce((s, item) => s + item.total, 0)
})

const totalCurrent = computed(() => {
  return currentList.value.reduce((s, item) => s + item.current, 0)
})

const totalOverdue = computed(() => {
  return currentList.value.reduce(
    (s, item) => s + (item.days1to30 + item.days31to60 + item.days61to90 + item.over90),
    0
  )
})

function handleExportExcel() {
  const isAR = activeTab.value === 'AR'
  const title = isAR ? 'Laporan_Umur_Piutang_AR' : 'Laporan_Umur_Hutang_AP'

  const data = currentList.value.map(item => ({
    'Entitas': item.entityName,
    'No Invoice / Klaim': item.invoiceNumber,
    'Tgl Invoice': item.invoiceDate,
    'Jatuh Tempo': item.dueDate,
    'Lancar (Current)': item.current,
    '1 - 30 Hari': item.days1to30,
    '31 - 60 Hari': item.days31to60,
    '61 - 90 Hari': item.days61to90,
    '> 90 Hari': item.over90,
    'Total (IDR)': item.total,
    'Status': item.status
  }))

  exportToExcel(data, `${title}_${accountingStore.businessUnit}_2026-03`, isAR ? 'AR Aging' : 'AP Aging')
  accountingStore.showToast(`Laporan ${isAR ? 'Piutang' : 'Hutang'} berhasil diekspor ke Excel`, 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Working Capital & Credit Management
          </span>
          <span class="text-xs text-slate-400">&bull; Analisis Umur Hutang Piutang</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Laporan Umur Piutang (AR) & Hutang (AP)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Monitoring jatuh tempo klaim asuransi pasien dan kewajiban pembayaran distributor farmasi berdasarkan rentang hari penuaan.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-4 h-4" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>
      </div>
    </div>

    <!-- Toggle Buttons: Piutang (AR) vs Hutang (AP) -->
    <div class="flex items-center gap-3">
      <button
        @click="activeTab = 'AR'"
        :class="[
          'flex items-center gap-2 px-4 py-2.5 rounded-xl font-bold text-xs transition-all',
          activeTab === 'AR'
            ? 'bg-blue-600 text-white shadow-sm'
            : 'bg-white border border-slate-200 text-slate-700 hover:bg-slate-50'
        ]"
      >
        <CreditCard class="w-4 h-4" />
        <span>Umur Piutang Pasien & Asuransi (AR)</span>
      </button>

      <button
        @click="activeTab = 'AP'"
        :class="[
          'flex items-center gap-2 px-4 py-2.5 rounded-xl font-bold text-xs transition-all',
          activeTab === 'AP'
            ? 'bg-blue-600 text-white shadow-sm'
            : 'bg-white border border-slate-200 text-slate-700 hover:bg-slate-50'
        ]"
      >
        <Receipt class="w-4 h-4" />
        <span>Umur Hutang Usaha & Distributor (AP)</span>
      </button>
    </div>

    <!-- Overview Metrics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs space-y-1">
        <span class="text-[11px] font-bold text-slate-500 uppercase">
          Total Saldo Terutang:
        </span>
        <div class="font-mono text-xl font-bold text-slate-900">
          {{ formatIDR(totalAmount) }}
        </div>
        <p class="text-[11px] text-slate-400">Total eksposur kredit aktif</p>
      </div>

      <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs space-y-1">
        <span class="text-[11px] font-bold text-emerald-600 uppercase flex items-center gap-1">
          <CheckCircle2 class="w-3.5 h-3.5" /> Lancar (Current / Belum JT):
        </span>
        <div class="font-mono text-xl font-bold text-emerald-700">
          {{ formatIDR(totalCurrent) }}
        </div>
        <p class="text-[11px] text-slate-400">
          {{ totalAmount > 0 ? Math.round((totalCurrent / totalAmount) * 100) : 0 }}% dalam batas tempo normal
        </p>
      </div>

      <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs space-y-1">
        <span class="text-[11px] font-bold text-rose-600 uppercase flex items-center gap-1">
          <AlertCircle class="w-3.5 h-3.5" /> Lewat Jatuh Tempo (Overdue):
        </span>
        <div class="font-mono text-xl font-bold text-rose-700">
          {{ formatIDR(totalOverdue) }}
        </div>
        <p class="text-[11px] text-slate-400">Memerlukan follow-up penagihan / pembayaran</p>
      </div>
    </div>

    <!-- Aging Bucket Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-700 uppercase tracking-wider">
              <th class="px-4 py-3">Nama Mitra / Entitas</th>
              <th class="px-4 py-3">No Invoice / Klaim</th>
              <th class="px-4 py-3">Jatuh Tempo</th>
              <th class="px-4 py-3 text-right text-emerald-800">Lancar (Current)</th>
              <th class="px-4 py-3 text-right">1 - 30 Hari</th>
              <th class="px-4 py-3 text-right">31 - 60 Hari</th>
              <th class="px-4 py-3 text-right">61 - 90 Hari</th>
              <th class="px-4 py-3 text-right text-rose-700">> 90 Hari</th>
              <th class="px-4 py-3 text-right font-extrabold">Total Terutang</th>
              <th class="px-4 py-3 text-center w-24">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="item in currentList"
              :key="item.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 font-bold text-slate-900">
                {{ item.entityName }}
              </td>
              <td class="px-4 py-3.5 font-mono text-blue-700">
                {{ item.invoiceNumber }}
              </td>
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ formatDate(item.dueDate) }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono text-slate-700">
                {{ item.current > 0 ? formatIDR(item.current, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono text-amber-700 font-semibold">
                {{ item.days1to30 > 0 ? formatIDR(item.days1to30, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono text-amber-700 font-semibold">
                {{ item.days31to60 > 0 ? formatIDR(item.days31to60, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono text-rose-700 font-semibold">
                {{ item.days61to90 > 0 ? formatIDR(item.days61to90, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono text-rose-700 font-bold">
                {{ item.over90 > 0 ? formatIDR(item.over90, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-slate-900 bg-slate-50/50">
                {{ formatIDR(item.total) }}
              </td>
              <td class="px-4 py-3.5 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    item.status === 'CURRENT'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : 'bg-rose-50 text-rose-700 border border-rose-200'
                  ]"
                >
                  {{ item.status }}
                </span>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
              <td colspan="3" class="px-4 py-3 text-right uppercase text-slate-800">
                TOTAL KESELURUHAN:
              </td>
              <td class="px-4 py-3 text-right font-mono text-emerald-900">
                {{ formatIDR(totalCurrent, false) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-amber-900">
                {{ formatIDR(currentList.reduce((s, i) => s + i.days1to30, 0), false) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-amber-900">
                {{ formatIDR(currentList.reduce((s, i) => s + i.days31to60, 0), false) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-rose-900">
                {{ formatIDR(currentList.reduce((s, i) => s + i.days61to90, 0), false) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-rose-900">
                {{ formatIDR(currentList.reduce((s, i) => s + i.over90, 0), false) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-slate-950 text-sm">
                {{ formatIDR(totalAmount) }}
              </td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>
