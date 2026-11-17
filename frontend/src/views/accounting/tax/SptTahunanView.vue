<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import {
  FileText,
  CheckCircle2,
  Download,
  Printer,
  Sparkles,
  Link2,
  ArrowRight,
  ArrowLeft,
  SlidersHorizontal,
  RotateCcw,
  AlertCircle,
  Edit3,
  Send,
  Info,
  X,
  ShieldCheck,
  Building2,
  Calculator
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await taxStore.fetchTaxRecords()
})

const annual1771 = computed(() => taxStore.annualTaxReport1771)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

// Month Level Editing Modal
const isEditModalOpen = ref(false)
const editingMonth = ref<{
  monthNum: number
  monthName: string
  taxType: 'PPH_23' | 'PPH_25'
  originalAmount: number
  adjustedAmount: number
  claimStatus: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED'
  reason: string
} | null>(null)

function openEditMonth(
  monthNum: number,
  monthName: string,
  taxType: 'PPH_23' | 'PPH_25',
  original: number,
  currentAdjusted: number,
  status: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED',
  reason: string
) {
  editingMonth.value = {
    monthNum,
    monthName,
    taxType,
    originalAmount: original,
    adjustedAmount: currentAdjusted,
    claimStatus: status,
    reason: reason || ''
  }
  isEditModalOpen.value = true
}

function handleStatusChange(status: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED') {
  if (!editingMonth.value) return
  editingMonth.value.claimStatus = status
  if (status === 'CLAIM_FULL') {
    editingMonth.value.adjustedAmount = editingMonth.value.originalAmount
  } else if (status === 'EXCLUDED') {
    editingMonth.value.adjustedAmount = 0
  }
}

function handleSaveMonthAdjustment() {
  if (!editingMonth.value) return
  taxStore.setMonthAdjustment({
    monthNum: editingMonth.value.monthNum,
    taxType: editingMonth.value.taxType,
    originalAmount: editingMonth.value.originalAmount,
    adjustedAmount: Number(editingMonth.value.adjustedAmount) || 0,
    claimStatus: editingMonth.value.claimStatus,
    reason: editingMonth.value.reason
  })
  isEditModalOpen.value = false
  editingMonth.value = null
}

function handleResetMonthAdjustment() {
  if (!editingMonth.value) return
  taxStore.resetMonthAdjustment(editingMonth.value.monthNum, editingMonth.value.taxType)
  isEditModalOpen.value = false
  editingMonth.value = null
}

// Global Adjustment Form
const pph23GeneralInput = ref(taxStore.annualGeneralAdjustment.pph23Adjustment)
const pph25GeneralInput = ref(taxStore.annualGeneralAdjustment.pph25Adjustment)
const notesGeneralInput = ref(taxStore.annualGeneralAdjustment.notes)

function applyGeneralAdjustment() {
  taxStore.setAnnualGeneralAdjustment({
    pph23Adjustment: Number(pph23GeneralInput.value) || 0,
    pph25Adjustment: Number(pph25GeneralInput.value) || 0,
    notes: notesGeneralInput.value
  })
  accountingStore.showToast('Penyesuaian global kredit pajak SPT 1771 tersimpan!', 'success')
}

function applyPresetOverwithholding() {
  pph23GeneralInput.value = -2400000
  notesGeneralInput.value = 'Koreksi kelebihan potong PPh 23 Masa Januari vendor PT Wastec (bukti potong direvisi non-kredit)'
  applyGeneralAdjustment()
}

function handleResetAll() {
  taxStore.resetAllAdjustments()
  pph23GeneralInput.value = 0
  pph25GeneralInput.value = 0
  notesGeneralInput.value = ''
}

function handlePostJournal() {
  taxStore.postPph29AccrualJournal()
}

function handleExportExcel() {
  accountingStore.showToast('Kertas Kerja SPT Tahunan PPh Badan 1771 berhasil diekspor ke Excel', 'success')
}

function printReport() {
  window.print()
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-12">
    <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Perpajakan Badan Usaha
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Formulir 1771 DJP</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          SPT Tahunan PPh Badan (Formulir 1771)
        </h1>
        <p class="text-xs text-zinc-500">
          Rekonsiliasi fiskal tahunan & akumulasi kredit pajak riil yang dihimpun otomatis dari SPT Masa PPh 23 & PPh 25 via kanal bank.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <router-link
          to="/accounting/tax/spt-masa"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Lihat Data SPT Masa</span>
        </router-link>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor 1771 (.xlsx)</span>
        </button>

        <button
          @click="printReport"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak Formulir 1771</span>
        </button>
      </div>
    </header>

    <!-- Info Banner: Connected with SPT Masa (Strictly Bank Only) & Audit Note -->
    <div class="bg-zinc-100 border border-zinc-200 rounded-2xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-zinc-950 text-white flex items-center justify-center font-bold shrink-0 shadow-xs">
          <Link2 class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2 flex-wrap">
            <h4 class="font-bold text-zinc-950 text-sm">Kredit Pajak SPT 1771 Bersumber dari Histori SPT Masa</h4>
            <span class="px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 font-bold text-[10px]">Non-Cash DJP Standard</span>
            <span
              v-if="annual1771.hasAdjustments"
              class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 border border-amber-300 font-bold text-[10px] flex items-center gap-1"
            >
              <AlertCircle class="w-3 h-3 text-amber-600" />
              <span>Ada Penyesuaian Tahunan / Kelebihan Potong</span>
            </span>
          </div>
          <p class="text-zinc-600 mt-0.5">
            Kredit Pajak PPh 23 dan PPh 25 diakumulasi dari SPT Masa. Jika terjadi kelebihan pemotongan atau salah klasifikasi, lakukan penyesuaian di SPT Tahunan tanpa mengubah bukti transaksi perbankan SPT Masa bulanan.
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <router-link
          to="/accounting/tax/spt-masa"
          class="px-3.5 py-2 bg-white hover:bg-zinc-50 border border-zinc-300 text-zinc-800 rounded-xl font-bold transition shadow-2xs text-xs whitespace-nowrap"
        >
          <span>Kelola SPT Masa Bulanan</span>
        </router-link>
      </div>
    </div>

    <!-- Annual Tax Summary Box -->
    <div class="bg-zinc-950 text-white rounded-2xl p-6 shadow-xl border border-zinc-800 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-400/30 font-bold text-[11px] flex items-center gap-1">
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
              <span>Validasi Terverifikasi (Kanal Bank Resmi)</span>
            </span>
            <span class="text-xs text-zinc-400">Tahun Pajak {{ annual1771.fiscalYear }}</span>
          </div>
          <h2 class="text-xl font-black text-white mt-1">
            Perhitungan PPh Badan Terutang & PPh Pasal 29 Kurang Bayar
          </h2>
          <p class="text-xs text-zinc-400">
            Laba Bersih Komersial dikurangi kredit pajak resmi yang dihimpun murni dari mutasi perbankan SPT Masa PPh 23 dan PPh 25 serta penyesuaian tahunan.
          </p>
        </div>

        <div class="text-right">
          <span class="text-xs text-zinc-400 block font-medium">
            {{ annual1771.pph28aOverpaid > 0 ? 'PPh Pasal 28A Lebih Bayar:' : 'PPh Pasal 29 Kurang Bayar:' }}
          </span>
          <span
            class="text-2xl font-black tabular-nums"
            :class="annual1771.pph28aOverpaid > 0 ? 'text-emerald-400' : 'text-amber-400'"
          >
            {{ formatIDR(annual1771.pph28aOverpaid > 0 ? annual1771.pph28aOverpaid : annual1771.pph29Payable) }}
          </span>
        </div>
      </div>

      <!-- 4 Pillars: Komersial -> PKP -> PPh Terutang -> Kredit Pajak Bank -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">1. Laba Bersih Komersial</span>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.commercialNetProfit) }}
          </div>
          <p class="text-[10px] text-zinc-400">Pendapatan dikurangi HPP & Beban Operasional</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">2. Penghasilan Kena Pajak</span>
          <div class="tabular-nums text-base font-bold text-amber-300">
            {{ formatIDR(annual1771.taxableIncome) }}
          </div>
          <p class="text-[10px] text-zinc-400">Setelah Koreksi Fiskal Netto (+Rp 48.95M)</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">3. PPh Badan Terutang</span>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.corporateTaxPayable) }}
          </div>
          <p class="text-[10px] text-zinc-400">Fasilitas Pasal 31E UU PPh (11% x PKP)</p>
        </div>

        <div class="bg-emerald-950/80 rounded-xl p-4 border border-emerald-700/60 space-y-1 text-emerald-100">
          <div class="flex items-center justify-between">
            <span class="text-emerald-400 font-bold uppercase tracking-wide text-[10px]">4. Kredit Pajak Diakui</span>
            <span v-if="annual1771.hasAdjustments" class="text-[9px] font-bold text-amber-300 bg-amber-950/80 px-1.5 py-0.5 rounded border border-amber-600/40">
              Ada Koreksi
            </span>
          </div>
          <div class="tabular-nums text-base font-bold text-emerald-300">
            {{ formatIDR(annual1771.taxCredits) }}
          </div>
          <p class="text-[10px] text-emerald-400 font-medium">
            PPh 23: {{ formatIDR(annual1771.taxCreditPph23) }} | PPh 25: {{ formatIDR(annual1771.taxCreditPph25) }}
          </p>
        </div>
      </div>
    </div>

    <!-- CROSS-RECONCILIATION TABLE: Rekapitulasi Aliran Kredit Pajak dari SPT Masa & Aksi Penyesuaian -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <Sparkles class="w-4 h-4 text-amber-500" />
            <h3 class="text-sm font-bold text-zinc-900">Rekapitulasi Bulanan SPT Masa Sebagai Sumber Kredit Pajak SPT Tahunan 1771</h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Diambil secara otomatis dari dokumen SPT Masa yang telah terverifikasi kanal bank resmi. Anda dapat menyesuaikan klaim per bulan jika terjadi kelebihan pemotongan.
          </p>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-xs tabular-nums font-bold text-zinc-950 bg-zinc-100 px-3 py-1.5 rounded-xl border border-zinc-200">
            Total Kredit Netto: {{ formatIDR(annual1771.taxCredits) }}
          </span>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-3 py-2.5">Masa Pajak</th>
              <th class="px-3 py-2.5">Kanal Rekening Bank Pembayar</th>
              <th class="px-3 py-2.5 text-right text-zinc-800">Kredit PPh 23 (Lamp. III)</th>
              <th class="px-3 py-2.5 text-right text-emerald-800">Angsuran PPh 25 (Induk 10a)</th>
              <th class="px-3 py-2.5 text-right text-zinc-600">PPh 21 Medis</th>
              <th class="px-3 py-2.5 text-right text-zinc-600">PPh 4(2) Final</th>
              <th class="px-3 py-2.5 text-right font-bold text-zinc-950 bg-zinc-50/50">Total Kredit 1771</th>
              <th class="px-3 py-2.5 text-center">Status Klaim 1771</th>
              <th class="px-3 py-2.5 text-center">Aksi Penyesuaian</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr v-for="row in sptMatrix" :key="row.monthNum" class="hover:bg-zinc-50/50">
              <td class="px-3 py-2.5 font-bold font-sans text-zinc-800">
                {{ row.monthName }} (Masa {{ row.monthNum }})
              </td>
              <td class="px-3 py-2.5 font-sans font-medium text-zinc-700">
                <span class="text-zinc-950 font-bold">{{ row.bankUsed }}</span>
                <span class="block text-[10px] text-zinc-400">Verifikasi Mutasi Rekening Bank</span>
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-zinc-900">
                <div v-if="row.pph23Amount > 0 || row.pph23Original > 0">
                  <span :class="row.pph23Amount !== row.pph23Original ? 'text-amber-600 font-black' : ''">
                    {{ formatIDR(row.pph23Amount, false) }}
                  </span>
                  <span v-if="row.pph23Amount !== row.pph23Original" class="block text-[10px] text-zinc-400 line-through">
                    SPT Masa: {{ formatIDR(row.pph23Original, false) }}
                  </span>
                </div>
                <span v-else class="text-zinc-400">-</span>
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-emerald-700">
                <div v-if="row.pph25Amount > 0 || row.pph25Original > 0">
                  <span :class="row.pph25Amount !== row.pph25Original ? 'text-amber-600 font-black' : ''">
                    {{ formatIDR(row.pph25Amount, false) }}
                  </span>
                  <span v-if="row.pph25Amount !== row.pph25Original" class="block text-[10px] text-zinc-400 line-through">
                    SPT Masa: {{ formatIDR(row.pph25Original, false) }}
                  </span>
                </div>
                <span v-else class="text-zinc-400">-</span>
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-zinc-500">
                {{ row.pph21Amount > 0 ? formatIDR(row.pph21Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-zinc-500">
                {{ row.pph42Amount > 0 ? formatIDR(row.pph42Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-zinc-950 bg-zinc-50/50">
                {{ row.totalCreditThisMonth > 0 ? formatIDR(row.totalCreditThisMonth, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-center font-sans">
                <span
                  v-if="row.isAdjusted"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-800 border border-amber-200"
                  :title="`Koreksi: ${formatIDR(row.adjustmentDiff)} (${row.pph23Reason || row.pph25Reason || 'Kelebihan potong'})`"
                >
                  Disesuaikan ({{ formatIDR(row.adjustmentDiff) }})
                </span>
                <span
                  v-else-if="row.isReported"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
                >
                  Klaim Penuh (100%)
                </span>
                <span
                  v-else
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-500 border border-zinc-200"
                >
                  Belum Ada
                </span>
              </td>
              <td class="px-3 py-2.5 text-center">
                <button
                  v-if="row.pph23Original > 0 || row.pph25Original > 0"
                  @click="openEditMonth(
                    row.monthNum,
                    row.monthName,
                    row.pph23Original > 0 ? 'PPH_23' : 'PPH_25',
                    row.pph23Original > 0 ? row.pph23Original : row.pph25Original,
                    row.pph23Original > 0 ? row.pph23Amount : row.pph25Amount,
                    row.pph23Original > 0 ? row.pph23Status : row.pph25Status,
                    row.pph23Original > 0 ? row.pph23Reason : row.pph25Reason
                  )"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-[11px] font-semibold bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 transition active:scale-95 cursor-pointer shadow-2xs"
                  title="Sesuaikan kredit bulan ini karena kelebihan potong"
                >
                  <Edit3 class="w-3 h-3 text-zinc-500" />
                  <span>{{ row.isAdjusted ? 'Edit Koreksi' : 'Koreksi' }}</span>
                </button>
                <span v-else class="text-zinc-300">-</span>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-zinc-950 text-white font-bold border-t-2 border-zinc-800">
              <td class="px-3 py-3 uppercase font-sans text-amber-400">TOTAL KREDIT PAJAK TAHUNAN:</td>
              <td class="px-3 py-3 font-sans text-zinc-400 text-xs">Rekap Kanal Bank & Koreksi 1771</td>
              <td class="px-3 py-3 text-right tabular-nums text-white">{{ formatIDR(annual1771.taxCreditPph23) }}</td>
              <td class="px-3 py-3 text-right tabular-nums text-emerald-400">{{ formatIDR(annual1771.taxCreditPph25) }}</td>
              <td class="px-3 py-3 text-right text-zinc-400">-</td>
              <td class="px-3 py-3 text-right text-zinc-400">-</td>
              <td class="px-3 py-3 text-right tabular-nums text-amber-300 text-sm font-black bg-zinc-900">{{ formatIDR(annual1771.taxCredits) }}</td>
              <td class="px-3 py-3 text-center text-emerald-400 font-sans text-xs" colspan="2">
                {{ annual1771.hasAdjustments ? 'Diakui dengan Penyesuaian' : 'Terverifikasi Penuh' }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- PANEL PENYESUAIAN & REKLASIFIKASI KREDIT PAJAK TAHUNAN (KELEBIHAN POTONG) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <SlidersHorizontal class="w-4 h-4 text-zinc-900" />
            <h3 class="text-sm font-bold text-zinc-900">
              Penyesuaian & Reklasifikasi Kredit Pajak Tahunan (Kelebihan Potong / Koreksi Bukti Potong)
            </h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Gunakan panel ini jika terdapat kelebihan pemotongan pada SPT Masa atau bukti potong yang tidak diakui fiskus, sehingga harus disesuaikan di SPT Tahunan 1771 tanpa merusak data perbankan bulanan.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="applyPresetOverwithholding"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-700 rounded-xl text-xs font-semibold transition active:scale-95 cursor-pointer"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Simulasi Kelebihan Potong PPh 23 (-Rp 2,4M)</span>
          </button>

          <button
            v-if="annual1771.hasAdjustments"
            type="button"
            @click="handleResetAll"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 rounded-xl text-xs font-semibold transition active:scale-95 cursor-pointer border border-rose-200"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            <span>Reset Semua Penyesuaian</span>
          </button>
        </div>
      </div>

      <!-- Form Inputs Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
        <div>
          <label class="block font-bold text-zinc-700 mb-1">
            Penyesuaian Kredit PPh 23 Tahunan (IDR)
          </label>
          <div class="relative">
            <input
              v-model.number="pph23GeneralInput"
              type="number"
              placeholder="0 (misal -2400000 untuk kelebihan potong)"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>
          <p class="text-[10px] text-zinc-400 mt-1">
            Gunakan nilai negatif (-) jika terdapat kelebihan potong yang tidak dapat dikreditkan.
          </p>
        </div>

        <div>
          <label class="block font-bold text-zinc-700 mb-1">
            Penyesuaian Angsuran PPh 25 Tahunan (IDR)
          </label>
          <div class="relative">
            <input
              v-model.number="pph25GeneralInput"
              type="number"
              placeholder="0"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>
          <p class="text-[10px] text-zinc-400 mt-1">
            Koreksi jika ada setoran PPh 25 tambahan atau penyesuaian SKP/STP berjalan.
          </p>
        </div>

        <div>
          <label class="block font-bold text-zinc-700 mb-1">
            Alasan / Keterangan Penyesuaian Fiskal
          </label>
          <input
            v-model="notesGeneralInput"
            type="text"
            placeholder="Contoh: Kelebihan potong vendor limbah B3 PT Wastec"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
          />
          <p class="text-[10px] text-zinc-400 mt-1">
            Wajib dicantumkan untuk lampiran penjelasan SPT 1771 Induk.
          </p>
        </div>
      </div>

      <!-- Action Row: Apply & Post Journal -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2 border-t border-zinc-100">
        <div class="flex items-center gap-3 text-xs">
          <div class="bg-zinc-50 px-3 py-2 rounded-xl border border-zinc-200">
            <span class="text-zinc-500 font-medium">Asli SPT Masa:</span>
            <span class="font-bold text-zinc-900 ml-1.5 tabular-nums">
              {{ formatIDR(annual1771.baseSptMasaPph23 + annual1771.baseSptMasaPph25) }}
            </span>
          </div>

          <div class="bg-zinc-50 px-3 py-2 rounded-xl border border-zinc-200">
            <span class="text-zinc-500 font-medium">Total Koreksi:</span>
            <span
              class="font-bold ml-1.5 tabular-nums"
              :class="annual1771.adjustmentPph23 + annual1771.adjustmentPph25 < 0 ? 'text-rose-600' : 'text-emerald-600'"
            >
              {{ annual1771.adjustmentPph23 + annual1771.adjustmentPph25 >= 0 ? '+' : '' }}{{ formatIDR(annual1771.adjustmentPph23 + annual1771.adjustmentPph25) }}
            </span>
          </div>

          <div class="bg-zinc-950 text-white px-3 py-2 rounded-xl">
            <span class="text-zinc-400 font-medium">Kredit Bersih 1771:</span>
            <span class="font-bold text-amber-300 ml-1.5 tabular-nums">
              {{ formatIDR(annual1771.taxCredits) }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="applyGeneralAdjustment"
            class="px-4 py-2 bg-zinc-900 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer shadow-2xs"
          >
            <span>Simpan Koreksi Global</span>
          </button>

          <button
            type="button"
            @click="handlePostJournal"
            class="inline-flex items-center gap-1.5 px-4 py-2 bg-amber-500 hover:bg-amber-600 text-zinc-950 rounded-xl text-xs font-black transition active:scale-95 cursor-pointer shadow-xs"
            title="Mencatat jurnal akrual Beban PPh Badan vs Kredit Pajak vs Hutang PPh 29 ke Buku Besar"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Posting Jurnal Akrual PPh 29</span>
          </button>
        </div>
      </div>
    </div>

    <!-- SPT 1771 Full Breakdown Table -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <h3 class="text-sm font-bold text-zinc-900 border-b border-zinc-100 pb-2 flex items-center gap-2">
        <FileText class="w-4 h-4 text-zinc-700" />
        <span>Kertas Kerja Rekonsiliasi Fiskal SPT 1771 (Tahun Pajak {{ annual1771.fiscalYear }})</span>
      </h3>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-16">No</th>
              <th class="px-4 py-3">Komponen Perhitungan SPT Tahunan 1771</th>
              <th class="px-4 py-3">Dasar Hukum & Sumber Data Akuntansi</th>
              <th class="px-4 py-3 text-right w-64">Jumlah (IDR)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr class="bg-zinc-50/50">
              <td class="px-4 py-2.5 font-bold">1</td>
              <td class="px-4 py-2.5 font-bold font-sans">Peredaran Usaha & Pendapatan Bruto (Akun 4 & 6)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Laporan Laba Rugi Komersial Klinik</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">{{ formatIDR(annual1771.commercialRevenue) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">2</td>
              <td class="px-4 py-2.5 font-sans">Beban Operasional, HPP Obat & Jasa Medis (Akun 5 & 52)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Buku Besar Akun 5 & 52</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-700">({{ formatIDR(annual1771.commercialExpense) }})</td>
            </tr>
            <tr class="bg-zinc-100/60 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">3</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">Penghasilan Netto Komersial (1 - 2)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Laba Bersih Sebelum Pajak</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">{{ formatIDR(annual1771.commercialNetProfit) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">4</td>
              <td class="px-4 py-2.5 font-sans">Koreksi Fiskal Positif (+)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.positiveCorrectionNotes }}</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-emerald-700 font-bold">+ {{ formatIDR(annual1771.positiveFiscalCorrection) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">5</td>
              <td class="px-4 py-2.5 font-sans">Koreksi Fiskal Negatif (-)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.negativeCorrectionNotes }}</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-rose-600 font-bold">- {{ formatIDR(annual1771.negativeFiscalCorrection) }}</td>
            </tr>
            <tr class="bg-amber-50 font-bold border-t border-b border-amber-200 text-amber-950">
              <td class="px-4 py-2.5 font-black">6</td>
              <td class="px-4 py-2.5 font-black uppercase">Penghasilan Kena Pajak (PKP) Badan</td>
              <td class="px-4 py-2.5 text-amber-800 font-sans">Dasar Pengenaan PPh Badan</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-amber-950 font-black text-sm">{{ formatIDR(annual1771.taxableIncome) }}</td>
            </tr>
            <tr class="bg-zinc-50">
              <td class="px-4 py-2.5 font-bold">7</td>
              <td class="px-4 py-2.5 font-sans">Tarif PPh Badan Efektif (Fasilitas Ps 31E UU HPP)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">50% x 22% = 11% (Peredaran &lt; Rp 50M)</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">11,00 %</td>
            </tr>
            <tr class="bg-zinc-100/80 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">8</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">PPh Badan Terutang Tahun Pajak {{ annual1771.fiscalYear }}</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Formulir 1771 Induk Bagian B</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">{{ formatIDR(annual1771.corporateTaxPayable) }}</td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60">
              <td class="px-4 py-2.5 font-bold">9</td>
              <td class="px-4 py-2.5 font-sans font-bold">Kredit Pajak PPh 23 Diakui (SPT Masa & Koreksi Tahunan)</td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">
                Lampiran Khusus III (Masa: {{ formatIDR(annual1771.baseSptMasaPph23) }} | Koreksi: {{ annual1771.adjustmentPph23 >= 0 ? '+' : '' }}{{ formatIDR(annual1771.adjustmentPph23) }})
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">({{ formatIDR(annual1771.taxCreditPph23) }})</td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60">
              <td class="px-4 py-2.5 font-bold">10</td>
              <td class="px-4 py-2.5 font-sans font-bold">Kredit Pajak PPh 25 Diakui (SPT Masa & Koreksi Tahunan)</td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">
                Induk Baris 10a (Masa: {{ formatIDR(annual1771.baseSptMasaPph25) }} | Koreksi: {{ annual1771.adjustmentPph25 >= 0 ? '+' : '' }}{{ formatIDR(annual1771.adjustmentPph25) }})
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">({{ formatIDR(annual1771.taxCreditPph25) }})</td>
            </tr>
            <tr class="bg-zinc-950 text-white font-black border-t-2 border-zinc-800">
              <td class="px-4 py-3 font-black text-amber-400">11</td>
              <td class="px-4 py-3 uppercase text-amber-400 text-sm">
                {{ annual1771.pph28aOverpaid > 0 ? 'PPH PASAL 28A LEBIH BAYAR (9 + 10 - 8)' : 'PPH PASAL 29 KURANG BAYAR (8 - 9 - 10)' }}
              </td>
              <td class="px-4 py-3 text-zinc-400 font-sans text-xs">
                {{ annual1771.pph28aOverpaid > 0 ? 'Dapat direstitusi atau dikompensasikan ke tahun depan' : 'Wajib disetor sebelum pelaporan SPT 1771' }}
              </td>
              <td
                class="px-4 py-3 text-right tabular-nums font-black text-base"
                :class="annual1771.pph28aOverpaid > 0 ? 'text-emerald-400' : 'text-amber-300'"
              >
                {{ formatIDR(annual1771.pph28aOverpaid > 0 ? annual1771.pph28aOverpaid : annual1771.pph29Payable) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Formal Clinical Signature Section -->
      <ReportSignatures
        preparer-role="Supervisor Pajak & Akuntansi"
        preparer-name="Rian Pratama, S.Ak."
        reviewer-role="Tax Manager & Konsultan Pajak"
        reviewer-name="Dewi Lestari, SE, Ak., BKP"
        approver-role="Direktur Keuangan (CFO)"
        approver-name="Budi Santoso, SE, Ak., CA"
      />
    </div>

    <!-- MODAL PENYESUAIAN KREDIT BULANAN (KELEBIHAN POTONG / REKLASIFIKASI) -->
    <div
      v-if="isEditModalOpen && editingMonth"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/60 backdrop-blur-xs"
    >
      <div class="bg-white rounded-2xl border border-zinc-200 shadow-2xl max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 border-b border-zinc-100 flex items-center justify-between bg-zinc-50/50">
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Penyesuaian Kredit SPT Tahunan</span>
            <h3 class="text-sm font-bold text-zinc-950">
              Koreksi {{ editingMonth.taxType === 'PPH_23' ? 'PPh 23' : 'PPh 25' }} - {{ editingMonth.monthName }} (Masa {{ editingMonth.monthNum }})
            </h3>
          </div>
          <button
            @click="isEditModalOpen = false"
            class="w-7 h-7 rounded-lg hover:bg-zinc-200 flex items-center justify-center text-zinc-400 hover:text-zinc-700 transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-6 space-y-4 text-xs">
          <!-- Information Banner -->
          <div class="bg-amber-50 border border-amber-200 rounded-xl p-3 text-amber-900 space-y-1">
            <div class="flex items-center gap-1.5 font-bold text-xs">
              <Info class="w-4 h-4 text-amber-700" />
              <span>Koreksi ini hanya memengaruhi SPT Tahunan 1771</span>
            </div>
            <p class="text-[11px] text-amber-800">
              Nilai SPT Masa asli (Rp {{ formatIDR(editingMonth.originalAmount) }}) tetap tercatat sebagai mutasi audit kanal bank dan tidak dihapus.
            </p>
          </div>

          <!-- Claim Status Selector -->
          <div>
            <label class="block font-bold text-zinc-700 mb-1.5">Status Pengakuan Kredit di SPT 1771</label>
            <div class="grid grid-cols-3 gap-2">
              <button
                type="button"
                @click="handleStatusChange('CLAIM_FULL')"
                class="px-2.5 py-2 rounded-xl text-center border font-bold text-[11px] transition cursor-pointer"
                :class="editingMonth.claimStatus === 'CLAIM_FULL' ? 'bg-zinc-950 text-white border-zinc-950' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                Klaim Penuh (100%)
              </button>
              <button
                type="button"
                @click="handleStatusChange('CLAIM_PARTIAL')"
                class="px-2.5 py-2 rounded-xl text-center border font-bold text-[11px] transition cursor-pointer"
                :class="editingMonth.claimStatus === 'CLAIM_PARTIAL' ? 'bg-amber-500 text-zinc-950 border-amber-500' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                Klaim Sebagian
              </button>
              <button
                type="button"
                @click="handleStatusChange('EXCLUDED')"
                class="px-2.5 py-2 rounded-xl text-center border font-bold text-[11px] transition cursor-pointer"
                :class="editingMonth.claimStatus === 'EXCLUDED' ? 'bg-rose-600 text-white border-rose-600' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                Tidak Diklaim
              </button>
            </div>
          </div>

          <!-- Nominal Inputs -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-medium text-zinc-500 mb-1">Nilai Asli SPT Masa</label>
              <input
                type="text"
                readonly
                :value="formatIDR(editingMonth.originalAmount)"
                class="w-full px-3 py-2 bg-zinc-100 border border-zinc-200 rounded-xl text-zinc-600 font-bold tabular-nums cursor-not-allowed"
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-800 mb-1">Nilai Diakui di 1771</label>
              <input
                v-model.number="editingMonth.adjustedAmount"
                type="number"
                class="w-full px-3 py-2 bg-white border border-zinc-300 rounded-xl font-bold font-mono text-zinc-950 focus:outline-none focus:ring-2 focus:ring-zinc-950"
              />
            </div>
          </div>

          <!-- Difference preview -->
          <div class="flex items-center justify-between px-3 py-2 bg-zinc-50 rounded-xl border border-zinc-200 font-mono text-[11px]">
            <span class="text-zinc-500">Selisih Koreksi:</span>
            <span
              class="font-bold tabular-nums"
              :class="editingMonth.adjustedAmount - editingMonth.originalAmount < 0 ? 'text-rose-600' : 'text-emerald-600'"
            >
              {{ editingMonth.adjustedAmount - editingMonth.originalAmount >= 0 ? '+' : '' }}{{ formatIDR(editingMonth.adjustedAmount - editingMonth.originalAmount) }}
            </span>
          </div>

          <!-- Reason Input -->
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Alasan Penyesuaian / Reklasifikasi</label>
            <textarea
              v-model="editingMonth.reason"
              rows="2"
              placeholder="Contoh: Kelebihan potong vendor PT Wastec, belum ada bukti setor revisi"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs resize-none"
            ></textarea>
          </div>
        </div>

        <div class="px-6 py-3.5 bg-zinc-50 border-t border-zinc-100 flex items-center justify-between">
          <button
            type="button"
            @click="handleResetMonthAdjustment"
            class="px-3 py-1.5 bg-white hover:bg-zinc-100 border border-zinc-200 text-zinc-700 rounded-xl font-semibold transition active:scale-95 cursor-pointer text-xs"
          >
            Kembalikan ke Nilai Awal
          </button>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="isEditModalOpen = false"
              class="px-3.5 py-1.5 text-zinc-600 hover:text-zinc-900 rounded-xl font-semibold transition cursor-pointer text-xs"
            >
              Batal
            </button>
            <button
              type="button"
              @click="handleSaveMonthAdjustment"
              class="px-4 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold transition active:scale-95 cursor-pointer text-xs shadow-xs"
            >
              Simpan Koreksi
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
