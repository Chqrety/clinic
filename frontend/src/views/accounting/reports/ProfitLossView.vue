<script setup lang="ts">
import { computed } from 'vue'
import { useReportsStore } from '@/stores/reports'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import {
  Receipt,
  Download,
  Printer,
  Building2,
  TrendingUp,
  Percent
} from 'lucide-vue-next'

const reportsStore = useReportsStore()
const accountingStore = useAccountingStore()

const summary = computed(() => reportsStore.profitLossSummary)
const revItems = computed(() => reportsStore.revenueItems)
const cogsItems = computed(() => reportsStore.cogsItems)
const opexItems = computed(() => reportsStore.operationalExpenseItems)

function handleExportExcel() {
  const exportRows: any[] = []

  exportRows.push({ 'Kategori': 'PENDAPATAN', 'Kode': '', 'Nama': '', 'Kidinest': '', 'SSC': '', 'Total': '' })
  revItems.value.forEach(i => {
    exportRows.push({
      'Kategori': 'Pendapatan',
      'Kode': i.code,
      'Nama': i.name,
      'Kidinest': i.kidinestAmount,
      'SSC': i.sscAmount,
      'Total': i.totalAmount
    })
  })
  exportRows.push({ 'Kategori': 'TOTAL PENDAPATAN', 'Kode': '', 'Nama': '', 'Kidinest': summary.value.kidinestRev, 'SSC': summary.value.sscRev, 'Total': summary.value.totalRev })

  exportRows.push({ 'Kategori': 'BEBAN POKOK PENDAPATAN (HPP)', 'Kode': '', 'Nama': '', 'Kidinest': '', 'SSC': '', 'Total': '' })
  cogsItems.value.forEach(i => {
    exportRows.push({
      'Kategori': 'HPP',
      'Kode': i.code,
      'Nama': i.name,
      'Kidinest': i.kidinestAmount,
      'SSC': i.sscAmount,
      'Total': i.totalAmount
    })
  })
  exportRows.push({ 'Kategori': 'LABA KOTOR (GROSS PROFIT)', 'Kode': '', 'Nama': '', 'Kidinest': summary.value.kidinestGrossProfit, 'SSC': summary.value.sscGrossProfit, 'Total': summary.value.grossProfit })

  exportRows.push({ 'Kategori': 'BEBAN OPERASIONAL', 'Kode': '', 'Nama': '', 'Kidinest': '', 'SSC': '', 'Total': '' })
  opexItems.value.forEach(i => {
    exportRows.push({
      'Kategori': 'Beban Operasional',
      'Kode': i.code,
      'Nama': i.name,
      'Kidinest': i.kidinestAmount,
      'SSC': i.sscAmount,
      'Total': i.totalAmount
    })
  })
  exportRows.push({ 'Kategori': 'LABA BERSIH (NET PROFIT)', 'Kode': '', 'Nama': '', 'Kidinest': summary.value.kidinestNetProfit, 'SSC': summary.value.sscNetProfit, 'Total': summary.value.netProfit })

  exportToExcel(exportRows, `Laporan_Laba_Rugi_${accountingStore.businessUnit}_2026-03`, 'Laba Rugi')
  accountingStore.showToast('Laporan Laba Rugi berhasil diekspor ke Excel (.xlsx)', 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs print:hidden">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Laporan Keuangan Resmi
          </span>
          <span class="text-xs text-slate-400">&bull; Statement of Profit & Loss</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Laporan Laba Rugi Komparatif
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Perbandingan performa pendapatan, HPP, beban operasional, dan margin laba bersih antara unit <strong>Kidinest</strong> dan <strong>SSC</strong>.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-xs font-semibold border border-slate-200 transition-colors"
        >
          <Printer class="w-4 h-4 text-slate-600" />
          <span>Cetak Laporan</span>
        </button>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-4 h-4" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>
      </div>
    </div>

    <!-- Official Financial Statement Presentation Sheet -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs p-8 max-w-5xl mx-auto space-y-6">
      <!-- Report Header -->
      <div class="text-center pb-6 border-b border-slate-200 space-y-1">
        <h2 class="text-lg font-extrabold uppercase tracking-wide text-slate-900">
          KLINIK TERPADU KIDINEST & SKIN SLIMMING CLINIC (SSC)
        </h2>
        <h3 class="text-sm font-bold text-blue-600 uppercase tracking-wider">
          LAPORAN LABA RUGI KOMPREHENSIF (KOMPARATIF)
        </h3>
        <p class="text-xs text-slate-500">
          Untuk Periode yang Berakhir pada 31 Maret 2026 (Dinyatakan dalam Rupiah Indonesia)
        </p>
      </div>

      <!-- Financial Table -->
      <div class="overflow-x-auto text-xs">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50 border-b-2 border-slate-300 text-[11px] font-bold text-slate-700 uppercase">
              <th class="py-2.5 px-3 w-28">Kode Akun</th>
              <th class="py-2.5 px-3">Uraian Akun Finansial</th>
              <th class="py-2.5 px-3 text-right w-40 text-blue-900">Kidinest (Rp)</th>
              <th class="py-2.5 px-3 text-right w-40 text-rose-900">SSC (Rp)</th>
              <th class="py-2.5 px-3 text-right w-44 font-extrabold text-slate-900 bg-slate-100/60">Total Konsolidasi</th>
              <th class="py-2.5 px-3 text-right w-20">% Omset</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <!-- 1. REVENUE SECTION -->
            <tr class="bg-slate-100/80 font-bold text-slate-900">
              <td colspan="6" class="py-2 px-3 text-xs uppercase tracking-wider text-blue-900">
                1. PENDAPATAN OPERASIONAL KLINIK
              </td>
            </tr>
            <tr v-for="item in revItems" :key="item.code" class="hover:bg-slate-50/50">
              <td class="py-2 px-3 font-mono text-slate-500">{{ item.code }}</td>
              <td class="py-2 px-3 pl-6 font-medium text-slate-800">{{ item.name }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.kidinestAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.sscAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono font-bold text-slate-900 bg-slate-50/40">{{ formatIDR(item.totalAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-500">{{ item.percentage }}%</td>
            </tr>
            <!-- Total Revenue Row -->
            <tr class="bg-blue-50/40 font-bold border-t border-b border-blue-200">
              <td colspan="2" class="py-2.5 px-3 text-right text-blue-950 uppercase font-bold">
                TOTAL PENDAPATAN OPERASIONAL:
              </td>
              <td class="py-2.5 px-3 text-right font-mono text-blue-900 font-bold">{{ formatIDR(summary.kidinestRev) }}</td>
              <td class="py-2.5 px-3 text-right font-mono text-rose-900 font-bold">{{ formatIDR(summary.sscRev) }}</td>
              <td class="py-2.5 px-3 text-right font-mono font-extrabold text-slate-950 bg-blue-100/50">{{ formatIDR(summary.totalRev) }}</td>
              <td class="py-2.5 px-3 text-right font-mono font-bold text-blue-900">100.0%</td>
            </tr>

            <!-- 2. COGS SECTION -->
            <tr class="bg-slate-100/80 font-bold text-slate-900">
              <td colspan="6" class="py-2 px-3 text-xs uppercase tracking-wider text-blue-900 pt-4">
                2. BEBAN POKOK PENDAPATAN (HPP & SHARE DOKTER)
              </td>
            </tr>
            <tr v-for="item in cogsItems" :key="item.code" class="hover:bg-slate-50/50">
              <td class="py-2 px-3 font-mono text-slate-500">{{ item.code }}</td>
              <td class="py-2 px-3 pl-6 font-medium text-slate-800">{{ item.name }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.kidinestAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.sscAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono font-bold text-slate-900 bg-slate-50/40">{{ formatIDR(item.totalAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-500">{{ item.percentage }}%</td>
            </tr>
            <!-- Total COGS Row -->
            <tr class="bg-amber-50/30 font-semibold border-t border-b border-amber-200">
              <td colspan="2" class="py-2 px-3 text-right text-amber-950 uppercase">
                TOTAL BEBAN POKOK PENDAPATAN (HPP):
              </td>
              <td class="py-2 px-3 text-right font-mono text-amber-900">({{ formatIDR(summary.kidinestCogs, false) }})</td>
              <td class="py-2 px-3 text-right font-mono text-amber-900">({{ formatIDR(summary.sscCogs, false) }})</td>
              <td class="py-2 px-3 text-right font-mono font-bold text-amber-950 bg-amber-50/60">({{ formatIDR(summary.totalCogs, false) }})</td>
              <td class="py-2 px-3 text-right font-mono text-amber-900">35.9%</td>
            </tr>

            <!-- GROSS PROFIT ROW -->
            <tr class="bg-emerald-50/50 font-extrabold border-t-2 border-b-2 border-emerald-300">
              <td colspan="2" class="py-3 px-3 text-right text-emerald-950 uppercase text-xs">
                LABA KOTOR (GROSS PROFIT):
              </td>
              <td class="py-3 px-3 text-right font-mono text-emerald-900">{{ formatIDR(summary.kidinestGrossProfit) }}</td>
              <td class="py-3 px-3 text-right font-mono text-emerald-900">{{ formatIDR(summary.sscGrossProfit) }}</td>
              <td class="py-3 px-3 text-right font-mono font-extrabold text-emerald-950 bg-emerald-100/50">{{ formatIDR(summary.grossProfit) }}</td>
              <td class="py-3 px-3 text-right font-mono text-emerald-900">64.1%</td>
            </tr>

            <!-- 3. OPEX SECTION -->
            <tr class="bg-slate-100/80 font-bold text-slate-900">
              <td colspan="6" class="py-2 px-3 text-xs uppercase tracking-wider text-blue-900 pt-4">
                3. BEBAN OPERASIONAL & UMUM (OPEX)
              </td>
            </tr>
            <tr v-for="item in opexItems" :key="item.code" class="hover:bg-slate-50/50">
              <td class="py-2 px-3 font-mono text-slate-500">{{ item.code }}</td>
              <td class="py-2 px-3 pl-6 font-medium text-slate-800">{{ item.name }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.kidinestAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-700">{{ formatIDR(item.sscAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono font-bold text-slate-900 bg-slate-50/40">{{ formatIDR(item.totalAmount, false) }}</td>
              <td class="py-2 px-3 text-right font-mono text-slate-500">{{ item.percentage }}%</td>
            </tr>
            <!-- Total OPEX Row -->
            <tr class="bg-slate-50 font-semibold border-t border-b border-slate-300">
              <td colspan="2" class="py-2.5 px-3 text-right text-slate-800 uppercase">
                TOTAL BEBAN OPERASIONAL:
              </td>
              <td class="py-2.5 px-3 text-right font-mono text-slate-800">({{ formatIDR(summary.kidinestOpex, false) }})</td>
              <td class="py-2.5 px-3 text-right font-mono text-slate-800">({{ formatIDR(summary.sscOpex, false) }})</td>
              <td class="py-2.5 px-3 text-right font-mono font-bold text-slate-900 bg-slate-100/50">({{ formatIDR(summary.totalOpex, false) }})</td>
              <td class="py-2.5 px-3 text-right font-mono text-slate-700">23.8%</td>
            </tr>

            <!-- NET PROFIT ROW -->
            <tr class="bg-blue-900 text-white font-extrabold border-t-2 border-b-4 border-slate-900">
              <td colspan="2" class="py-3.5 px-3 text-right text-white uppercase text-xs tracking-wider">
                LABA BERSIH TAHUN BERJALAN (NET PROFIT):
              </td>
              <td class="py-3.5 px-3 text-right font-mono text-sm text-blue-200">{{ formatIDR(summary.kidinestNetProfit) }}</td>
              <td class="py-3.5 px-3 text-right font-mono text-sm text-blue-200">{{ formatIDR(summary.sscNetProfit) }}</td>
              <td class="py-3.5 px-3 text-right font-mono text-sm font-extrabold text-white bg-blue-950">{{ formatIDR(summary.netProfit) }}</td>
              <td class="py-3.5 px-3 text-right font-mono text-sm text-emerald-300">{{ summary.netMarginPercent }}%</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Signatures Footer for Formal Printing -->
      <div class="pt-10 grid grid-cols-3 gap-6 text-center text-xs text-slate-700 border-t border-slate-200">
        <div class="space-y-12">
          <p class="font-bold">Dibuat Oleh,</p>
          <div>
            <p class="font-bold border-b border-slate-400 pb-1 max-w-[180px] mx-auto">Rian Pratama, S.Ak.</p>
            <p class="text-[11px] text-slate-500">Staf Akuntansi & Pajak</p>
          </div>
        </div>

        <div class="space-y-12">
          <p class="font-bold">Direview Oleh,</p>
          <div>
            <p class="font-bold border-b border-slate-400 pb-1 max-w-[180px] mx-auto">Dewi Lestari, SE, Ak.</p>
            <p class="text-[11px] text-slate-500">Accounting Supervisor</p>
          </div>
        </div>

        <div class="space-y-12">
          <p class="font-bold">Disetujui Oleh,</p>
          <div>
            <p class="font-bold border-b border-slate-400 pb-1 max-w-[180px] mx-auto">Budi Santoso, SE, Ak., CA</p>
            <p class="text-[11px] text-slate-500">Finance & Accounting Manager</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
