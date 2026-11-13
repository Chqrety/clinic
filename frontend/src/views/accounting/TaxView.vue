<script setup lang="ts">
import { computed } from 'vue'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  Calculator,
  Download,
  FileSpreadsheet,
  CheckCircle2,
  FileText,
  Building,
  UserCheck
} from 'lucide-vue-next'

const taxStore = useTaxStore()
const accountingStore = useAccountingStore()

const doctorTaxRecords = computed(() => taxStore.filteredDoctorTax)
const ppn = computed(() => taxStore.ppnSummary)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Perpajakan & Kepatuhan Fiskal
          </span>
          <span class="text-xs text-slate-400">&bull; PMK 168/2023 & UU HPP</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Kalkulasi PPh 21 Dokter & Rekap PPN
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Perhitungan withholding tax PPh 21 dokter spesialis dengan tarif progresif berkesinambungan serta rekapitulasi PPN Masukan & Keluaran.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="taxStore.exportEbupotCSV()"
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-4 h-4" />
          <span>Ekspor e-Bupot 21/26 (DJP CSV)</span>
        </button>
      </div>
    </div>

    <!-- PPN Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-slate-500 text-xs">
          <span class="font-bold uppercase">PPN Keluaran (Output)</span>
          <span class="px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-bold text-[10px]">Tarif 11%</span>
        </div>
        <div class="font-mono text-xl font-bold text-slate-900">
          {{ formatIDR(ppn.ppnOutputTax) }}
        </div>
        <p class="text-[11px] text-slate-400">
          DPP: {{ formatIDR(ppn.ppnOutputGross) }} (Penjualan Obat Non-Resep & Skincare)
        </p>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-slate-500 text-xs">
          <span class="font-bold uppercase">PPN Masukan (Input)</span>
          <span class="px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 font-bold text-[10px]">Faktur Pajak</span>
        </div>
        <div class="font-mono text-xl font-bold text-emerald-700">
          {{ formatIDR(ppn.ppnInputTax) }}
        </div>
        <p class="text-[11px] text-slate-400">
          DPP: {{ formatIDR(ppn.ppnInputGross) }} (Pembelian Supplier Distributor)
        </p>
      </div>

      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-slate-500 text-xs">
          <span class="font-bold uppercase">PPN Kurang Bayar (Net SPT)</span>
          <span class="px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-bold text-[10px]">Setor Kas Negara</span>
        </div>
        <div class="font-mono text-xl font-bold text-amber-800">
          {{ formatIDR(ppn.ppnPayableNet) }}
        </div>
        <p class="text-[11px] text-slate-400">
          Kewajiban setor akhir bulan masa pajak Maret 2026
        </p>
      </div>
    </div>

    <!-- PPh 21 Dokter Progressive Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden space-y-4 p-5">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
        <div class="flex items-center gap-2">
          <UserCheck class="w-4 h-4 text-blue-600" />
          <h3 class="text-sm font-bold text-slate-900">Perhitungan PPh 21 Dokter Spesialis (PMK 168/2023)</h3>
        </div>
        <span class="text-xs font-semibold text-slate-500">
          Metode Non-Pegawai: 50% x Bruto x Tarif Progresif Pasal 17
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-3 py-2.5">Nama Dokter & Spesialisasi</th>
              <th class="px-3 py-2.5">NPWP / NIK</th>
              <th class="px-3 py-2.5 text-center">Unit</th>
              <th class="px-3 py-2.5 text-right">Fee Bruto (Bulan Ini)</th>
              <th class="px-3 py-2.5 text-right">Kumulatif Bruto (YTD)</th>
              <th class="px-3 py-2.5 text-right">DPP (50% Bruto)</th>
              <th class="px-3 py-2.5 text-center">Lapisan Tarif</th>
              <th class="px-3 py-2.5 text-right text-rose-700 font-bold">Potongan PPh 21</th>
              <th class="px-3 py-2.5 text-right text-emerald-800 font-bold">Fee Netto Dokter</th>
              <th class="px-3 py-2.5 text-center">No e-Bupot</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="d in doctorTaxRecords"
              :key="d.id"
              class="hover:bg-slate-50/50 transition-colors"
            >
              <td class="px-3 py-3">
                <span class="font-bold text-slate-900 block">{{ d.doctorName }}</span>
                <span class="text-[11px] text-slate-400">{{ d.specialty }}</span>
              </td>
              <td class="px-3 py-3 font-mono text-[11px] text-slate-600">
                <span>{{ d.npwp }}</span>
                <span class="block text-slate-400 text-[10px]">NIK: {{ d.nik }}</span>
              </td>
              <td class="px-3 py-3 text-center">
                <StatusBadge :status="d.businessUnit" size="sm" />
              </td>
              <td class="px-3 py-3 text-right font-mono font-semibold text-slate-800">
                {{ formatIDR(d.grossIncomeCurrentMonth) }}
              </td>
              <td class="px-3 py-3 text-right font-mono text-slate-500">
                {{ formatIDR(d.grossIncomeCumulativeYear) }}
              </td>
              <td class="px-3 py-3 text-right font-mono text-blue-900 font-semibold bg-blue-50/20">
                {{ formatIDR(d.dppAmount) }}
              </td>
              <td class="px-3 py-3 text-center">
                <span class="px-2 py-0.5 rounded bg-slate-100 text-slate-700 font-bold text-[10px] border border-slate-200">
                  {{ d.taxBracketApplied }}
                </span>
              </td>
              <td class="px-3 py-3 text-right font-mono font-bold text-rose-700 bg-rose-50/30">
                {{ formatIDR(d.pph21Amount) }}
              </td>
              <td class="px-3 py-3 text-right font-mono font-bold text-emerald-800 bg-emerald-50/30">
                {{ formatIDR(d.netPayable) }}
              </td>
              <td class="px-3 py-3 text-center font-mono text-[10px] text-slate-500">
                {{ d.ebupotNumber }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
              <td colspan="3" class="px-3 py-2.5 text-right uppercase text-slate-700">
                TOTAL:
              </td>
              <td class="px-3 py-2.5 text-right font-mono">
                {{ formatIDR(doctorTaxRecords.reduce((s, d) => s + d.grossIncomeCurrentMonth, 0)) }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono text-slate-500">
                -
              </td>
              <td class="px-3 py-2.5 text-right font-mono">
                {{ formatIDR(doctorTaxRecords.reduce((s, d) => s + d.dppAmount, 0)) }}
              </td>
              <td></td>
              <td class="px-3 py-2.5 text-right font-mono text-rose-800">
                {{ formatIDR(doctorTaxRecords.reduce((s, d) => s + d.pph21Amount, 0)) }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono text-emerald-800">
                {{ formatIDR(doctorTaxRecords.reduce((s, d) => s + d.netPayable, 0)) }}
              </td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

