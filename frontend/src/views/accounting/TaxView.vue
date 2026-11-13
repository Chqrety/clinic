<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { formatIDR } from '@/utils/currency'
import StatusBadge from '@/components/common/StatusBadge.vue'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import {
  Calculator,
  Download,
  FileSpreadsheet,
  CheckCircle2,
  FileText,
  Building,
  UserCheck,
  Receipt,
  Layers,
  ArrowRight,
  Send,
  Calendar,
  AlertCircle,
  ExternalLink,
  ShieldCheck,
  BookOpen
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()

// Read query params if any
const route = router.currentRoute.value
if (route.query.taxCategory) {
  const cat = String(route.query.taxCategory)
  const matched = taxStore.taxAccountOptions.find(opt => opt.taxCategory === cat)
  if (matched) taxStore.selectedTaxAccountCode = matched.code
}

const doctorTaxRecords = computed(() => taxStore.filteredDoctorTax)
const ppn = computed(() => taxStore.ppnSummary)
const txList = computed(() => taxStore.filteredTaxTransactions)
const summary = computed(() => taxStore.monthlyTaxSummary)
const annual1771 = computed(() => taxStore.annualTaxReport1771)

// Selected Account Metadata
const selectedAccountMeta = computed(() => {
  return taxStore.taxAccountOptions.find(o => o.code === taxStore.selectedTaxAccountCode) || taxStore.taxAccountOptions[0]
})

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/ledger', query: { account: accountCode } })
  } else {
    router.push('/ledger')
  }
}

function handleGenerateNTPNJournal() {
  const code = taxStore.selectedTaxAccountCode === 'ALL' ? '2501' : taxStore.selectedTaxAccountCode
  const name = selectedAccountMeta.value.name
  const amt = summary.value.totalTax || 39249649
  taxStore.generateTaxPaymentJournal(code, name, amt)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Perpajakan & Kepatuhan Fiskal Terintegrasi COA
          </span>
          <span class="text-xs text-slate-400">&bull; PMK 168/2023, PP 58/2023 & UU HPP</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Manajemen Pajak Masa & Pajak Tahunan (SPT 1771)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Perhitungan pajak otomatis dari Chart of Accounts (COA kelompok 25 - Hutang Pajak), Buku Besar (General Ledger), dan Jurnal Transaksi.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2.5">
        <button
          @click="taxStore.exportEbupotCSV()"
          class="flex items-center gap-1.5 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Ekspor e-Bupot CSV (DJP)</span>
        </button>

        <button
          @click="handleGenerateNTPNJournal"
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Send class="w-3.5 h-3.5" />
          <span>Generate Jurnal Setor NTPN</span>
        </button>
      </div>
    </div>

    <!-- Mode Selector Tabs: Pajak Masa vs Pajak Tahunan (SPT 1771) -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-slate-900 text-white p-4 rounded-xl shadow-lg border border-slate-800">
      <div class="flex items-center gap-2">
        <div class="flex items-center bg-slate-800 p-1 rounded-lg border border-slate-700">
          <button
            @click="taxStore.taxPeriodType = 'MASA'"
            :class="[
              'px-4 py-1.5 rounded-md text-xs font-bold transition-all flex items-center gap-1.5',
              taxStore.taxPeriodType === 'MASA'
                ? 'bg-blue-600 text-white shadow-xs'
                : 'text-slate-400 hover:text-white'
            ]"
          >
            <Calendar class="w-3.5 h-3.5" />
            <span>Pajak Masa (Bulanan)</span>
          </button>
          <button
            @click="taxStore.taxPeriodType = 'TAHUNAN'"
            :class="[
              'px-4 py-1.5 rounded-md text-xs font-bold transition-all flex items-center gap-1.5',
              taxStore.taxPeriodType === 'TAHUNAN'
                ? 'bg-blue-600 text-white shadow-xs'
                : 'text-slate-400 hover:text-white'
            ]"
          >
            <Building class="w-3.5 h-3.5" />
            <span>Pajak Tahunan (SPT Badan 1771)</span>
          </button>
        </div>
      </div>

      <!-- COA Account Selector Dropdown (When in Pajak Masa Mode) -->
      <div v-if="taxStore.taxPeriodType === 'MASA'" class="flex items-center gap-3">
        <span class="text-xs text-slate-300 font-medium">Pilih Akun Pajak di COA:</span>
        <select
          v-model="taxStore.selectedTaxAccountCode"
          class="px-3 py-1.5 bg-slate-800 border border-slate-700 text-white rounded-lg text-xs font-semibold focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer max-w-xs"
        >
          <option
            v-for="opt in taxStore.taxAccountOptions"
            :key="opt.code"
            :value="opt.code"
          >
            {{ opt.name }}
          </option>
        </select>
      </div>

      <!-- Fiscal Year (When in Pajak Tahunan Mode) -->
      <div v-else class="flex items-center gap-3">
        <span class="text-xs text-slate-300 font-medium">Tahun Pajak (Fiskal):</span>
        <select
          v-model="taxStore.selectedYear"
          class="px-3 py-1.5 bg-slate-800 border border-slate-700 text-white rounded-lg text-xs font-bold focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer"
        >
          <option :value="2026">Tahun Pajak 2026</option>
          <option :value="2025">Tahun Pajak 2025</option>
        </select>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- VIEW A: PAJAK MASA (BULANAN) BERBASIS COA & BUKU BESAR -->
    <!-- ============================================================== -->
    <div v-if="taxStore.taxPeriodType === 'MASA'" class="space-y-6">
      <!-- Info Banner of Selected Tax Account -->
      <div class="bg-blue-50/70 border border-blue-200 rounded-xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-blue-600 text-white flex items-center justify-center shrink-0 font-bold">
            <Receipt class="w-5 h-5" />
          </div>
          <div>
            <h4 class="font-bold text-slate-900 text-sm">{{ selectedAccountMeta.name }}</h4>
            <p class="text-slate-600 mt-0.5">{{ selectedAccountMeta.description }} &bull; Standar Tarif: <strong>{{ selectedAccountMeta.normalRate }}</strong></p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="openLedgerForAccount(taxStore.selectedTaxAccountCode !== 'ALL' ? taxStore.selectedTaxAccountCode : '25')"
            class="px-3 py-1.5 bg-white hover:bg-slate-50 border border-slate-200 text-blue-700 rounded-lg font-bold flex items-center gap-1.5 transition-colors shadow-2xs"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Lihat Buku Besar Akun Ini</span>
          </button>
        </div>
      </div>

      <!-- KPI Summary Cards for Selected Tax Category -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <span class="text-xs font-bold text-slate-500 uppercase">Total Bruto Objek Pajak</span>
          <div class="font-mono text-xl font-bold text-slate-900">
            {{ formatIDR(summary.totalGross) }}
          </div>
          <p class="text-[11px] text-slate-400">
            Nilai transaksi sebelum pemotongan
          </p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <span class="text-xs font-bold text-slate-500 uppercase">Dasar Pengenaan Pajak (DPP)</span>
          <div class="font-mono text-xl font-bold text-blue-700">
            {{ formatIDR(summary.totalDpp) }}
          </div>
          <p class="text-[11px] text-slate-400">
            Dasar kalkulasi perkalian tarif
          </p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <span class="text-xs font-bold text-slate-500 uppercase">Pajak Terutang (Buku Besar)</span>
          <div class="font-mono text-xl font-bold text-rose-700">
            {{ formatIDR(summary.coaBalance) }}
          </div>
          <p class="text-[11px] text-slate-400">
            Saldo Akun COA Kelompok 25
          </p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <div class="flex items-center justify-between text-slate-500 text-xs">
            <span class="font-bold uppercase">Status Rekonsiliasi</span>
            <span class="px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 font-bold text-[10px]">Sinkron</span>
          </div>
          <div class="font-mono text-xl font-bold text-emerald-700 flex items-center gap-1.5">
            <CheckCircle2 class="w-5 h-5 text-emerald-600 inline" />
            <span>100% Sesuai</span>
          </div>
          <p class="text-[11px] text-slate-400">
            Jurnal entri cocok dengan buku besar
          </p>
        </div>
      </div>

      <!-- Specific PPN Summary Section if PPN is selected or ALL -->
      <div v-if="taxStore.selectedTaxAccountCode === 'ALL' || taxStore.selectedTaxAccountCode === '2504' || taxStore.selectedTaxAccountCode === '1504'" class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <div class="flex items-center justify-between text-slate-500 text-xs">
            <span class="font-bold uppercase">PPN Keluaran (Akun COA 2504)</span>
            <span class="px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-bold text-[10px]">Tarif 11%</span>
          </div>
          <div class="font-mono text-xl font-bold text-slate-900">
            {{ formatIDR(ppn.ppnOutputTax) }}
          </div>
          <p class="text-[11px] text-slate-400">
            DPP: {{ formatIDR(ppn.ppnOutputGross) }} (Penjualan Produk Skincare & Estetika)
          </p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <div class="flex items-center justify-between text-slate-500 text-xs">
            <span class="font-bold uppercase">PPN Masukan (Akun COA 1504)</span>
            <span class="px-2 py-0.5 rounded bg-emerald-50 text-emerald-700 font-bold text-[10px]">Faktur Masukan</span>
          </div>
          <div class="font-mono text-xl font-bold text-emerald-700">
            {{ formatIDR(ppn.ppnInputTax) }}
          </div>
          <p class="text-[11px] text-slate-400">
            DPP: {{ formatIDR(ppn.ppnInputGross) }} (Faktur Pembelian PBF Distributor)
          </p>
        </div>

        <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-1">
          <div class="flex items-center justify-between text-slate-500 text-xs">
            <span class="font-bold uppercase">PPN Kurang Bayar (Setoran SPT)</span>
            <span class="px-2 py-0.5 rounded bg-amber-50 text-amber-700 font-bold text-[10px]">NTPN DJP</span>
          </div>
          <div class="font-mono text-xl font-bold text-amber-800">
            {{ formatIDR(ppn.ppnPayableNet) }}
          </div>
          <p class="text-[11px] text-slate-400">
            Kewajiban setor ke Kas Negara sebelum akhir bulan
          </p>
        </div>
      </div>

      <!-- Specific PPh 21 Dokter Progressive Table if 2501 or ALL -->
      <div v-if="taxStore.selectedTaxAccountCode === 'ALL' || taxStore.selectedTaxAccountCode === '2501.001'" class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden space-y-4 p-5">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <UserCheck class="w-4 h-4 text-blue-600" />
            <h3 class="text-sm font-bold text-slate-900">Rincian Withholding PPh 21 Dokter Tenaga Ahli (PMK 168/2023)</h3>
          </div>
          <span class="text-xs font-semibold text-slate-500">
            Akun COA: 2501.001 &bull; Rumus: 50% Bruto x Tarif Progresif Pasal 17 UU HPP
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Nama Dokter & Spesialisasi</th>
                <th class="px-3 py-2.5">NPWP / NIK</th>
                <th class="px-3 py-2.5 text-center">Unit</th>
                <th class="px-3 py-2.5 text-right">Fee Bruto</th>
                <th class="px-3 py-2.5 text-right">DPP (50%)</th>
                <th class="px-3 py-2.5 text-center">Lapisan Tarif</th>
                <th class="px-3 py-2.5 text-right text-rose-700 font-bold">Potongan PPh 21</th>
                <th class="px-3 py-2.5 text-right text-emerald-800 font-bold">Netto Diterima</th>
                <th class="px-3 py-2.5 text-center">No Jurnal</th>
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
                <td class="px-3 py-3 text-center font-mono text-xs">
                  <span class="text-blue-600 hover:underline cursor-pointer font-bold" @click="router.push('/journals')">
                    {{ d.journalNumber }}
                  </span>
                </td>
                <td class="px-3 py-3 text-center font-mono text-[10px] text-slate-500">
                  {{ d.ebupotNumber }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Transaction Line Details directly extracted from General Ledger & Journal Lines -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden space-y-4 p-5">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <FileSpreadsheet class="w-4 h-4 text-emerald-600" />
            <h3 class="text-sm font-bold text-slate-900">Buku Pembantu Pajak (Kalkulasi dari Jurnal Entri & Buku Besar)</h3>
          </div>
          <span class="text-xs text-slate-500">
            Menampilkan <strong>{{ txList.length }}</strong> transaksi kena pajak periode November 2026
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Tanggal</th>
                <th class="px-3 py-2.5">No. Jurnal</th>
                <th class="px-3 py-2.5">Akun Pajak COA</th>
                <th class="px-3 py-2.5">Lawan Transaksi</th>
                <th class="px-3 py-2.5">NPWP / NIK</th>
                <th class="px-3 py-2.5">Objek Pajak</th>
                <th class="px-3 py-2.5 text-right">Nilai Bruto</th>
                <th class="px-3 py-2.5 text-right">DPP</th>
                <th class="px-3 py-2.5 text-center">Tarif</th>
                <th class="px-3 py-2.5 text-right font-bold text-rose-700">Pajak Terutang</th>
                <th class="px-3 py-2.5 text-center">No. Bukti / Faktur</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="tx in txList"
                :key="tx.id"
                class="hover:bg-slate-50/50 transition-colors"
              >
                <td class="px-3 py-3 font-mono text-[11px] text-slate-600 whitespace-nowrap">
                  {{ tx.taxDate }}
                </td>
                <td class="px-3 py-3 font-mono text-[11px] font-bold text-blue-600 hover:underline cursor-pointer" @click="router.push('/journals')">
                  {{ tx.journalNumber }}
                </td>
                <td class="px-3 py-3 font-mono font-bold text-purple-700">
                  {{ tx.accountCode }}
                </td>
                <td class="px-3 py-3 font-medium text-slate-900">
                  {{ tx.counterpartyName }}
                </td>
                <td class="px-3 py-3 font-mono text-[11px] text-slate-500 whitespace-nowrap">
                  {{ tx.counterpartyNPWP }}
                </td>
                <td class="px-3 py-3 text-[11px] text-slate-600">
                  {{ tx.objectCode }}
                </td>
                <td class="px-3 py-3 text-right font-mono font-medium text-slate-800">
                  {{ formatIDR(tx.grossAmount) }}
                </td>
                <td class="px-3 py-3 text-right font-mono font-medium text-blue-800 bg-blue-50/20">
                  {{ formatIDR(tx.dppAmount) }}
                </td>
                <td class="px-3 py-3 text-center font-mono font-bold text-slate-700">
                  {{ (tx.taxRate * 100).toFixed(1) }}%
                </td>
                <td class="px-3 py-3 text-right font-mono font-bold text-rose-700 bg-rose-50/30">
                  {{ formatIDR(tx.taxAmount) }}
                </td>
                <td class="px-3 py-3 text-center font-mono text-[10px] text-slate-500">
                  {{ tx.taxNumber }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
                <td colspan="6" class="px-3 py-2.5 text-right uppercase text-slate-700">
                  TOTAL KESELURUHAN:
                </td>
                <td class="px-3 py-2.5 text-right font-mono">
                  {{ formatIDR(txList.reduce((s, t) => s + t.grossAmount, 0)) }}
                </td>
                <td class="px-3 py-2.5 text-right font-mono text-blue-800">
                  {{ formatIDR(txList.reduce((s, t) => s + t.dppAmount, 0)) }}
                </td>
                <td></td>
                <td class="px-3 py-2.5 text-right font-mono text-rose-800">
                  {{ formatIDR(txList.reduce((s, t) => s + t.taxAmount, 0)) }}
                </td>
                <td></td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- VIEW B: PAJAK TAHUNAN (SPT TAHUNAN PPH BADAN FORMULIR 1771) -->
    <!-- ============================================================== -->
    <div v-else class="space-y-6">
      <!-- Annual Tax Summary Box -->
      <div class="bg-gradient-to-r from-slate-900 to-blue-950 text-white rounded-2xl p-6 shadow-xl border border-slate-800 space-y-4">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-800 pb-4">
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2.5 py-0.5 rounded-full bg-blue-500/20 text-blue-300 border border-blue-400/30 font-bold text-[11px]">
                SPT Tahunan PPh Badan Formulir 1771
              </span>
              <span class="text-xs text-slate-400">Tahun Pajak {{ annual1771.fiscalYear }}</span>
            </div>
            <h2 class="text-lg font-black text-white mt-1">
              Rekonsiliasi Fiskal & Perhitungan PPh Badan Terutang (Pasal 29)
            </h2>
            <p class="text-xs text-slate-400">
              Dihitung otomatis dari akumulasi Buku Besar Pendapatan (Akun 4 & 6) dan Beban (Akun 5 & 7).
            </p>
          </div>

          <div class="text-right">
            <span class="text-xs text-slate-400 block">PPh Pasal 29 Kurang Bayar:</span>
            <span class="text-2xl font-black font-mono text-amber-400">
              {{ formatIDR(annual1771.pph29Payable) }}
            </span>
          </div>
        </div>

        <!-- 3 Pillars: Komersial -> Koreksi Fiskal -> Pajak Terutang -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
          <div class="bg-slate-800/80 rounded-xl p-4 border border-slate-700/60 space-y-1.5">
            <span class="text-slate-400 font-bold uppercase tracking-wide">1. Laba Bersih Komersial</span>
            <div class="font-mono text-lg font-bold text-white">
              {{ formatIDR(annual1771.commercialNetProfit) }}
            </div>
            <p class="text-[11px] text-slate-400">
              Pendapatan: {{ formatIDR(annual1771.commercialRevenue) }} <br/>
              Beban: {{ formatIDR(annual1771.commercialExpense) }}
            </p>
          </div>

          <div class="bg-slate-800/80 rounded-xl p-4 border border-slate-700/60 space-y-1.5">
            <span class="text-slate-400 font-bold uppercase tracking-wide">2. Koreksi Fiskal Bersih</span>
            <div class="font-mono text-lg font-bold text-emerald-400">
              + {{ formatIDR(annual1771.positiveFiscalCorrection - annual1771.negativeFiscalCorrection) }}
            </div>
            <p class="text-[11px] text-slate-400">
              Koreksi Positif: +{{ formatIDR(annual1771.positiveFiscalCorrection) }} <br/>
              Koreksi Negatif: -{{ formatIDR(annual1771.negativeFiscalCorrection) }}
            </p>
          </div>

          <div class="bg-slate-800/80 rounded-xl p-4 border border-slate-700/60 space-y-1.5">
            <span class="text-slate-400 font-bold uppercase tracking-wide">3. Penghasilan Kena Pajak (PKP)</span>
            <div class="font-mono text-lg font-bold text-amber-300">
              {{ formatIDR(annual1771.taxableIncome) }}
            </div>
            <p class="text-[11px] text-slate-400">
              Tarif Efektif Fasilitas 31E: 11% (50% x 22%)
            </p>
          </div>
        </div>
      </div>

      <!-- SPT 1771 Breakdown Table -->
      <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden p-6 space-y-4">
        <h3 class="text-sm font-bold text-slate-900 border-b border-slate-100 pb-2 flex items-center gap-2">
          <FileText class="w-4 h-4 text-blue-600" />
          <span>Kertas Kerja Rekonsiliasi Fiskal SPT 1771 (Tahun Pajak {{ annual1771.fiscalYear }})</span>
        </h3>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
                <th class="px-4 py-3 w-16">No</th>
                <th class="px-4 py-3">Komponen Perhitungan SPT Tahunan</th>
                <th class="px-4 py-3">Dasar Hukum & Catatan Akuntansi</th>
                <th class="px-4 py-3 text-right w-64">Jumlah (IDR)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 font-mono">
              <tr class="bg-slate-50/50">
                <td class="px-4 py-2.5 font-bold">1</td>
                <td class="px-4 py-2.5 font-bold font-sans">Peredaran Usaha & Pendapatan Bruto (Akun 4 & 6)</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Penjualan Rumah / Pelayanan Medis + Bunga Bank</td>
                <td class="px-4 py-2.5 text-right font-bold text-slate-900">{{ formatIDR(annual1771.commercialRevenue) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">2</td>
                <td class="px-4 py-2.5 font-sans">Harga Pokok Penjualan (HPP Akun 51)</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Perolehan Tanah / Bahan Baku & Tindakan Medis</td>
                <td class="px-4 py-2.5 text-right text-slate-700">({{ formatIDR(9807287775.40) }})</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">3</td>
                <td class="px-4 py-2.5 font-sans">Biaya Usaha & Operasional (Akun 52 & 7)</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Gaji, Listrik, Sewa, Limbah B3 & Bunga Pinjaman</td>
                <td class="px-4 py-2.5 text-right text-slate-700">({{ formatIDR(3350282603.47) }})</td>
              </tr>
              <tr class="bg-blue-50/30 font-bold">
                <td class="px-4 py-2.5">4</td>
                <td class="px-4 py-2.5 font-sans">Penghasilan Netto Komersial Sebelum Pajak</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Laba Buku Besar Komersial</td>
                <td class="px-4 py-2.5 text-right text-blue-900">{{ formatIDR(annual1771.commercialNetProfit) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">5</td>
                <td class="px-4 py-2.5 font-sans text-emerald-700 font-semibold">Koreksi Fiskal Positif</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">{{ annual1771.positiveCorrectionNotes }}</td>
                <td class="px-4 py-2.5 text-right text-emerald-700 font-bold">+ {{ formatIDR(annual1771.positiveFiscalCorrection) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">6</td>
                <td class="px-4 py-2.5 font-sans text-rose-700 font-semibold">Koreksi Fiskal Negatif</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">{{ annual1771.negativeCorrectionNotes }}</td>
                <td class="px-4 py-2.5 text-right text-rose-700 font-bold">- {{ formatIDR(annual1771.negativeFiscalCorrection) }}</td>
              </tr>
              <tr class="bg-amber-50/40 font-bold border-y-2 border-amber-200">
                <td class="px-4 py-3">7</td>
                <td class="px-4 py-3 font-sans text-amber-950 font-bold">Penghasilan Kena Pajak (PKP)</td>
                <td class="px-4 py-3 text-slate-600 font-sans font-normal">Dasar Pengenaan Pajak Penghasilan Badan</td>
                <td class="px-4 py-3 text-right text-amber-900 text-sm font-black">{{ formatIDR(annual1771.taxableIncome) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">8</td>
                <td class="px-4 py-2.5 font-sans font-bold">PPh Badan Terutang (Tahun Berjalan)</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Fasilitas Pasal 31E UU PPh (11% x PKP)</td>
                <td class="px-4 py-2.5 text-right font-bold text-slate-900">{{ formatIDR(annual1771.corporateTaxPayable) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">9</td>
                <td class="px-4 py-2.5 font-sans text-emerald-700">Kredit Pajak Dalam Negeri (PPh 23 & PPh 25)</td>
                <td class="px-4 py-2.5 text-slate-500 font-sans">Kredit Pajak Dibayar di Muka (Akun COA 1504)</td>
                <td class="px-4 py-2.5 text-right font-bold text-emerald-700">({{ formatIDR(annual1771.taxCredits) }})</td>
              </tr>
              <tr class="bg-slate-900 text-white font-bold text-sm">
                <td class="px-4 py-3 text-amber-400">10</td>
                <td class="px-4 py-3 font-sans text-white font-bold">PPh Kurang Bayar (PPh Pasal 29) - Setor Kas Negara</td>
                <td class="px-4 py-3 text-slate-300 font-sans font-normal">Wajib disetor sebelum lapor SPT 1771 (Akun COA 2506)</td>
                <td class="px-4 py-3 text-right font-mono text-amber-400 font-black">{{ formatIDR(annual1771.pph29Payable) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="pt-4 flex items-center justify-end gap-3">
          <button
            @click="taxStore.generateTaxPaymentJournal('2506', 'PPh 29 Tahunan Badan', annual1771.pph29Payable)"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold transition-colors shadow-xs flex items-center gap-1.5"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Generate Jurnal Pelunasan PPh 29 (Akun 2506)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Official Report Signature Section matching Image 1 -->
    <ReportSignatures />
  </div>
</template>
