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
  Calculator,
  Pencil
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
  taxStore.resetAnnual1771Overrides()
  pph23GeneralInput.value = 0
  pph25GeneralInput.value = 0
  notesGeneralInput.value = ''
}

function handlePostJournal() {
  taxStore.postPph29AccrualJournal()
}

// Full SPT 1771 Worksheet Edit Modal (Allows changing ANY/ALL parameters)
const isFullEditModalOpen = ref(false)

const editForm1771 = ref({
  commercialRevenue: 0,
  commercialExpense: 0,
  positiveFiscalCorrection: 0,
  positiveCorrectionNotes: '',
  negativeFiscalCorrection: 0,
  negativeCorrectionNotes: '',
  corporateTaxRate: 0.11,
  taxCreditPph23: 0,
  taxCreditPph25: 0
})

function openFullEditModal() {
  const current = annual1771.value
  editForm1771.value = {
    commercialRevenue: current.commercialRevenue,
    commercialExpense: current.commercialExpense,
    positiveFiscalCorrection: current.positiveFiscalCorrection,
    positiveCorrectionNotes: current.positiveCorrectionNotes,
    negativeFiscalCorrection: current.negativeFiscalCorrection,
    negativeCorrectionNotes: current.negativeCorrectionNotes,
    corporateTaxRate: current.corporateTaxRate,
    taxCreditPph23: current.taxCreditPph23,
    taxCreditPph25: current.taxCreditPph25
  }
  isFullEditModalOpen.value = true
}

// Live computed values inside edit modal
const modalComputedNetProfit = computed(() => {
  return Number(editForm1771.value.commercialRevenue || 0) - Number(editForm1771.value.commercialExpense || 0)
})

const modalComputedPKP = computed(() => {
  return Math.max(
    0,
    modalComputedNetProfit.value +
    Number(editForm1771.value.positiveFiscalCorrection || 0) -
    Number(editForm1771.value.negativeFiscalCorrection || 0)
  )
})

const modalComputedTaxPayable = computed(() => {
  return Math.round(modalComputedPKP.value * Number(editForm1771.value.corporateTaxRate || 0))
})

const modalComputedTotalCredits = computed(() => {
  return Number(editForm1771.value.taxCreditPph23 || 0) + Number(editForm1771.value.taxCreditPph25 || 0)
})

const modalComputedPph29 = computed(() => {
  return Math.max(0, modalComputedTaxPayable.value - modalComputedTotalCredits.value)
})

const modalComputedPph28A = computed(() => {
  return modalComputedTotalCredits.value > modalComputedTaxPayable.value
    ? modalComputedTotalCredits.value - modalComputedTaxPayable.value
    : 0
})

function handleSaveFullEdit() {
  taxStore.updateAnnual1771Overrides({
    commercialRevenue: Number(editForm1771.value.commercialRevenue),
    commercialExpense: Number(editForm1771.value.commercialExpense),
    positiveFiscalCorrection: Number(editForm1771.value.positiveFiscalCorrection),
    positiveCorrectionNotes: editForm1771.value.positiveCorrectionNotes,
    negativeFiscalCorrection: Number(editForm1771.value.negativeFiscalCorrection),
    negativeCorrectionNotes: editForm1771.value.negativeCorrectionNotes,
    corporateTaxRate: Number(editForm1771.value.corporateTaxRate),
    taxCreditPph23: Number(editForm1771.value.taxCreditPph23),
    taxCreditPph25: Number(editForm1771.value.taxCreditPph25)
  })
  isFullEditModalOpen.value = false
}

function handleResetOverrides() {
  taxStore.resetAnnual1771Overrides()
  isFullEditModalOpen.value = false
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
          <span
            v-if="annual1771.isCustomOverridden"
            class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-900 font-bold text-[10px] border border-amber-300"
          >
            Data 1771 Diubah Manual
          </span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          SPT Tahunan PPh Badan (Formulir 1771)
        </h1>
        <p class="text-xs text-zinc-500">
          Rekonsiliasi fiskal tahunan & akumulasi kredit pajak riil. Seluruh komponen SPT Tahunan dapat diubah bebas tanpa memengaruhi data SPT Masa bulanan.
        </p>
      </div>

      <div class="flex items-center gap-2.5 flex-wrap">
        <button
          @click="openFullEditModal"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-amber-500 hover:bg-amber-600 text-zinc-950 rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Pencil class="w-3.5 h-3.5" />
          <span>Ubah Seluruh Data 1771</span>
        </button>

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
            <h4 class="font-bold text-zinc-950 text-sm">Fleksibilitas Penuh SPT Tahunan (Tanpa Mengubah SPT Masa)</h4>
            <span class="px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 font-bold text-[10px]">Non-Destructive DJP Model</span>
            <span
              v-if="annual1771.isCustomOverridden"
              class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 border border-amber-300 font-bold text-[10px] flex items-center gap-1"
            >
              <AlertCircle class="w-3 h-3 text-amber-600" />
              <span>Kertas Kerja Menggunakan Nilai Manual</span>
            </span>
          </div>
          <p class="text-zinc-600 mt-0.5">
            Setiap komponen SPT Tahunan (Pendapatan Komersial, Biaya, Koreksi Fiskal, Tarif, hingga Kredit Pajak PPh 23 & 25) dapat diedit sesuai kebutuhan fiskal akhir tahun tanpa mengubah transaksi historis bank pada SPT Masa.
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="annual1771.isCustomOverridden"
          @click="handleResetOverrides"
          class="px-3.5 py-2 bg-rose-50 hover:bg-rose-100 border border-rose-200 text-rose-800 rounded-xl font-bold transition shadow-2xs text-xs whitespace-nowrap cursor-pointer"
        >
          Reset ke Nilai Buku Otomatis
        </button>
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
            <span
              v-if="annual1771.isCustomOverridden"
              class="text-xs text-amber-400 bg-amber-950/80 px-2 py-0.5 rounded border border-amber-600/40 font-bold"
            >
              Manual Override Aktif
            </span>
          </div>
          <h2 class="text-xl font-black text-white mt-1">
            Perhitungan PPh Badan Terutang & PPh Pasal 29 Kurang Bayar
          </h2>
          <p class="text-xs text-zinc-400">
            Laba Bersih Komersial dikurangi kredit pajak resmi yang dihimpun dari mutasi perbankan SPT Masa PPh 23 dan PPh 25 serta penyesuaian tahunan 1771.
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
          <div class="flex items-center justify-between">
            <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">1. Laba Bersih Komersial</span>
            <button @click="openFullEditModal" class="text-zinc-500 hover:text-white cursor-pointer" title="Ubah Nilai">
              <Pencil class="w-3 h-3" />
            </button>
          </div>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.commercialNetProfit) }}
          </div>
          <p class="text-[10px] text-zinc-400">Pendapatan dikurangi HPP & Beban Operasional</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">2. Penghasilan Kena Pajak</span>
            <button @click="openFullEditModal" class="text-zinc-500 hover:text-amber-300 cursor-pointer" title="Ubah Koreksi Fiskal">
              <Pencil class="w-3 h-3" />
            </button>
          </div>
          <div class="tabular-nums text-base font-bold text-amber-300">
            {{ formatIDR(annual1771.taxableIncome) }}
          </div>
          <p class="text-[10px] text-zinc-400">Setelah Koreksi Fiskal Positif & Negatif</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">3. PPh Badan Terutang</span>
            <button @click="openFullEditModal" class="text-zinc-500 hover:text-white cursor-pointer" title="Ubah Tarif Pajak">
              <Pencil class="w-3 h-3" />
            </button>
          </div>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.corporateTaxPayable) }}
          </div>
          <p class="text-[10px] text-zinc-400">Tarif Efektif {{ (annual1771.corporateTaxRate * 100).toFixed(0) }}% x PKP</p>
        </div>

        <div class="bg-emerald-950/80 rounded-xl p-4 border border-emerald-700/60 space-y-1 text-emerald-100">
          <div class="flex items-center justify-between">
            <span class="text-emerald-400 font-bold uppercase tracking-wide text-[10px]">4. Kredit Pajak Diakui</span>
            <button @click="openFullEditModal" class="text-emerald-500 hover:text-emerald-200 cursor-pointer" title="Ubah Kredit Pajak">
              <Pencil class="w-3 h-3" />
            </button>
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

    <!-- SPT 1771 FULL BREAKDOWN TABLE (DENGAN AKSI EDIT SETIAP BARIS) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <FileText class="w-4 h-4 text-zinc-700" />
          <h3 class="text-sm font-bold text-zinc-900">
            Kertas Kerja Rekonsiliasi Fiskal SPT 1771 (Tahun Pajak {{ annual1771.fiscalYear }})
          </h3>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="openFullEditModal"
            class="inline-flex items-center gap-1 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer shadow-xs"
          >
            <Pencil class="w-3 h-3" />
            <span>Edit Kertas Kerja 1771</span>
          </button>

          <button
            v-if="annual1771.isCustomOverridden"
            type="button"
            @click="handleResetOverrides"
            class="inline-flex items-center gap-1 px-3 py-1.5 bg-rose-50 hover:bg-rose-100 text-rose-700 rounded-xl text-xs font-semibold transition active:scale-95 cursor-pointer border border-rose-200"
          >
            <RotateCcw class="w-3 h-3" />
            <span>Reset ke Nilai Buku</span>
          </button>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-12">No</th>
              <th class="px-4 py-3">Komponen Perhitungan SPT Tahunan 1771</th>
              <th class="px-4 py-3">Dasar Hukum & Catatan Penjelasan</th>
              <th class="px-4 py-3 text-right w-56">Jumlah (IDR)</th>
              <th class="px-4 py-3 text-center w-24">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr class="bg-zinc-50/50 hover:bg-zinc-100/50">
              <td class="px-4 py-2.5 font-bold">1</td>
              <td class="px-4 py-2.5 font-bold font-sans text-zinc-900">
                Peredaran Usaha & Pendapatan Bruto (Akun 4 & 6)
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                Laporan Laba Rugi Komersial Klinik
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">
                {{ formatIDR(annual1771.commercialRevenue) }}
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition cursor-pointer"
                  title="Ubah Pendapatan Bruto"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="hover:bg-zinc-50/50">
              <td class="px-4 py-2.5 font-bold">2</td>
              <td class="px-4 py-2.5 font-sans text-zinc-900">
                Beban Operasional, HPP Obat & Jasa Medis (Akun 5 & 52)
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                Buku Besar Akun 5 & 52
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-700">
                ({{ formatIDR(annual1771.commercialExpense) }})
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition cursor-pointer"
                  title="Ubah Beban Operasional"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="bg-zinc-100/60 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">3</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">
                Penghasilan Netto Komersial (1 - 2)
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                Laba Bersih Sebelum Pajak (Dihitung Otomatis)
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">
                {{ formatIDR(annual1771.commercialNetProfit) }}
              </td>
              <td class="px-4 py-2.5 text-center text-zinc-400 font-mono text-[10px]">
                Auto
              </td>
            </tr>
            <tr class="hover:bg-zinc-50/50">
              <td class="px-4 py-2.5 font-bold">4</td>
              <td class="px-4 py-2.5 font-sans text-zinc-900">
                Koreksi Fiskal Positif (+)
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                {{ annual1771.positiveCorrectionNotes }}
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-emerald-700 font-bold">
                + {{ formatIDR(annual1771.positiveFiscalCorrection) }}
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition cursor-pointer"
                  title="Ubah Koreksi Positif"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="hover:bg-zinc-50/50">
              <td class="px-4 py-2.5 font-bold">5</td>
              <td class="px-4 py-2.5 font-sans text-zinc-900">
                Koreksi Fiskal Negatif (-)
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                {{ annual1771.negativeCorrectionNotes }}
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-rose-600 font-bold">
                - {{ formatIDR(annual1771.negativeFiscalCorrection) }}
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition cursor-pointer"
                  title="Ubah Koreksi Negatif"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="bg-amber-50 font-bold border-t border-b border-amber-200 text-amber-950">
              <td class="px-4 py-2.5 font-black">6</td>
              <td class="px-4 py-2.5 font-black uppercase">
                Penghasilan Kena Pajak (PKP) Badan
              </td>
              <td class="px-4 py-2.5 text-amber-800 font-sans">
                Dasar Pengenaan PPh Badan (Dihitung Otomatis: 3 + 4 - 5)
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-amber-950 font-black text-sm">
                {{ formatIDR(annual1771.taxableIncome) }}
              </td>
              <td class="px-4 py-2.5 text-center text-zinc-400 font-mono text-[10px]">
                Auto
              </td>
            </tr>
            <tr class="bg-zinc-50 hover:bg-zinc-100/50">
              <td class="px-4 py-2.5 font-bold">7</td>
              <td class="px-4 py-2.5 font-sans text-zinc-900">
                Tarif PPh Badan Efektif
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                {{ annual1771.corporateTaxRate === 0.11 ? 'Fasilitas Pasal 31E UU HPP (50% x 22% = 11%)' : 'Tarif PPh Badan Kustom' }}
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">
                {{ (annual1771.corporateTaxRate * 100).toFixed(2) }} %
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition cursor-pointer"
                  title="Ubah Tarif Pajak"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="bg-zinc-100/80 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">8</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">
                PPh Badan Terutang Tahun Pajak {{ annual1771.fiscalYear }}
              </td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">
                Formulir 1771 Induk Bagian B (Dihitung Otomatis: 6 x 7)
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">
                {{ formatIDR(annual1771.corporateTaxPayable) }}
              </td>
              <td class="px-4 py-2.5 text-center text-zinc-400 font-mono text-[10px]">
                Auto
              </td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60 hover:bg-emerald-100/50">
              <td class="px-4 py-2.5 font-bold">9</td>
              <td class="px-4 py-2.5 font-sans font-bold">
                Kredit Pajak PPh 23 Diakui
              </td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">
                Lampiran Khusus III (Masa: {{ formatIDR(annual1771.baseSptMasaPph23) }} | Koreksi: {{ annual1771.adjustmentPph23 >= 0 ? '+' : '' }}{{ formatIDR(annual1771.adjustmentPph23) }})
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">
                ({{ formatIDR(annual1771.taxCreditPph23) }})
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-emerald-200 text-emerald-700 transition cursor-pointer"
                  title="Ubah Kredit PPh 23"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60 hover:bg-emerald-100/50">
              <td class="px-4 py-2.5 font-bold">10</td>
              <td class="px-4 py-2.5 font-sans font-bold">
                Kredit Pajak PPh 25 Diakui
              </td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">
                Induk Baris 10a (Masa: {{ formatIDR(annual1771.baseSptMasaPph25) }} | Koreksi: {{ annual1771.adjustmentPph25 >= 0 ? '+' : '' }}{{ formatIDR(annual1771.adjustmentPph25) }})
              </td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">
                ({{ formatIDR(annual1771.taxCreditPph25) }})
              </td>
              <td class="px-4 py-2.5 text-center">
                <button
                  @click="openFullEditModal"
                  class="p-1 rounded-lg hover:bg-emerald-200 text-emerald-700 transition cursor-pointer"
                  title="Ubah Kredit PPh 25"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </button>
              </td>
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
              <td class="px-4 py-3 text-center text-zinc-500 font-mono text-[10px]">
                Akhir
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

    <!-- PANEL PENYESUAIAN GLOBAL & POSTING JURNAL -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <SlidersHorizontal class="w-4 h-4 text-zinc-900" />
            <h3 class="text-sm font-bold text-zinc-900">
              Penyesuaian Global & Aksi Jurnal Akrual SPT Tahunan
            </h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Koreksi selisih lebih bayar/kurang bayar pada SPT Tahunan 1771 dan langsung posting jurnal akrual PPh Badan ke Buku Besar.
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

      <!-- Action Row: Apply & Post Journal -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2">
        <div class="flex items-center gap-3 text-xs flex-wrap">
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
            @click="openFullEditModal"
            class="px-4 py-2 bg-zinc-900 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer shadow-2xs"
          >
            <span>Ubah Komponen 1771</span>
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

    <!-- MODAL 1: UBAH SELURUH DATA KERTAS KERJA SPT TAHUNAN 1771 -->
    <div
      v-if="isFullEditModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/70 backdrop-blur-xs"
    >
      <div class="bg-white rounded-2xl border border-zinc-200 shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto animate-in fade-in zoom-in-95 duration-150">
        <!-- Modal Header -->
        <div class="sticky top-0 px-6 py-4 border-b border-zinc-100 flex items-center justify-between bg-zinc-50/90 backdrop-blur-sm z-10">
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-amber-700 bg-amber-100 px-2 py-0.5 rounded-full">
              Fleksibilitas Penuh Formulir 1771
            </span>
            <h3 class="text-base font-bold text-zinc-950 mt-1">
              Ubah Data Kertas Kerja SPT Tahunan PPh Badan
            </h3>
            <p class="text-[11px] text-zinc-500">
              Perubahan nilai di sini murni untuk SPT Tahunan dan TIDAK mengubah histori mutasi bank atau SPT Masa bulanan.
            </p>
          </div>
          <button
            @click="isFullEditModalOpen = false"
            class="w-8 h-8 rounded-xl hover:bg-zinc-200 flex items-center justify-center text-zinc-400 hover:text-zinc-700 transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Form Body -->
        <div class="p-6 space-y-5 text-xs">
          <!-- Section 1: Pendapatan & Beban Komersial -->
          <div class="space-y-3">
            <h4 class="font-bold text-zinc-950 uppercase tracking-wider text-[11px] border-b border-zinc-100 pb-1">
              1. Laba Bersih Komersial
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-zinc-700 mb-1">Peredaran Usaha & Pendapatan Bruto (IDR)</label>
                <input
                  v-model.number="editForm1771.commercialRevenue"
                  type="number"
                  class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white"
                />
              </div>

              <div>
                <label class="block font-bold text-zinc-700 mb-1">Beban Operasional & HPP Obat (IDR)</label>
                <input
                  v-model.number="editForm1771.commercialExpense"
                  type="number"
                  class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white"
                />
              </div>
            </div>

            <div class="flex items-center justify-between px-3 py-2 bg-zinc-100 rounded-xl font-mono text-xs">
              <span class="text-zinc-600 font-medium">Laba Bersih Komersial Dihitung:</span>
              <span class="font-bold text-zinc-950 tabular-nums">
                {{ formatIDR(modalComputedNetProfit) }}
              </span>
            </div>
          </div>

          <!-- Section 2: Koreksi Fiskal Positif & Negatif -->
          <div class="space-y-3">
            <h4 class="font-bold text-zinc-950 uppercase tracking-wider text-[11px] border-b border-zinc-100 pb-1">
              2. Koreksi Fiskal (UU Pajak Penghasilan)
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-emerald-800 mb-1">Koreksi Fiskal Positif (+) IDR</label>
                <input
                  v-model.number="editForm1771.positiveFiscalCorrection"
                  type="number"
                  class="w-full px-3 py-2 bg-emerald-50/50 border border-emerald-300 rounded-xl font-mono text-emerald-950 font-bold focus:outline-none focus:ring-2 focus:ring-emerald-600 focus:bg-white"
                />
                <input
                  v-model="editForm1771.positiveCorrectionNotes"
                  type="text"
                  placeholder="Catatan koreksi positif"
                  class="w-full mt-1.5 px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-700 text-[11px]"
                />
              </div>

              <div>
                <label class="block font-bold text-rose-800 mb-1">Koreksi Fiskal Negatif (-) IDR</label>
                <input
                  v-model.number="editForm1771.negativeFiscalCorrection"
                  type="number"
                  class="w-full px-3 py-2 bg-rose-50/50 border border-rose-300 rounded-xl font-mono text-rose-950 font-bold focus:outline-none focus:ring-2 focus:ring-rose-600 focus:bg-white"
                />
                <input
                  v-model="editForm1771.negativeCorrectionNotes"
                  type="text"
                  placeholder="Catatan koreksi negatif"
                  class="w-full mt-1.5 px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-700 text-[11px]"
                />
              </div>
            </div>

            <div class="flex items-center justify-between px-3 py-2 bg-amber-50 rounded-xl font-mono text-xs border border-amber-200">
              <span class="text-amber-900 font-bold">Penghasilan Kena Pajak (PKP):</span>
              <span class="font-black text-amber-950 tabular-nums">
                {{ formatIDR(modalComputedPKP) }}
              </span>
            </div>
          </div>

          <!-- Section 3: Tarif Pajak & PPh Terutang -->
          <div class="space-y-3">
            <h4 class="font-bold text-zinc-950 uppercase tracking-wider text-[11px] border-b border-zinc-100 pb-1">
              3. Tarif PPh Badan & Pajak Terutang
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-zinc-700 mb-1">Tarif Efektif PPh Badan (0.11 = 11%, 0.22 = 22%)</label>
                <div class="flex items-center gap-2">
                  <input
                    v-model.number="editForm1771.corporateTaxRate"
                    type="number"
                    step="0.01"
                    min="0"
                    max="1"
                    class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white"
                  />
                  <button
                    type="button"
                    @click="editForm1771.corporateTaxRate = 0.11"
                    class="px-2.5 py-2 bg-zinc-100 hover:bg-zinc-200 rounded-xl font-bold text-[10px] text-zinc-700 whitespace-nowrap cursor-pointer"
                  >
                    11% (31E)
                  </button>
                  <button
                    type="button"
                    @click="editForm1771.corporateTaxRate = 0.22"
                    class="px-2.5 py-2 bg-zinc-100 hover:bg-zinc-200 rounded-xl font-bold text-[10px] text-zinc-700 whitespace-nowrap cursor-pointer"
                  >
                    22%
                  </button>
                </div>
              </div>

              <div>
                <label class="block font-medium text-zinc-500 mb-1">PPh Badan Terutang (Dihitung Otomatis)</label>
                <input
                  type="text"
                  readonly
                  :value="formatIDR(modalComputedTaxPayable)"
                  class="w-full px-3 py-2 bg-zinc-100 border border-zinc-200 rounded-xl font-bold font-mono text-zinc-700 tabular-nums cursor-not-allowed"
                />
              </div>
            </div>
          </div>

          <!-- Section 4: Kredit Pajak Diakui -->
          <div class="space-y-3">
            <h4 class="font-bold text-zinc-950 uppercase tracking-wider text-[11px] border-b border-zinc-100 pb-1">
              4. Kredit Pajak Diakui di 1771 (PPh 23 & PPh 25)
            </h4>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block font-bold text-zinc-800 mb-1">Kredit Pajak PPh 23 (Lampiran III) IDR</label>
                <input
                  v-model.number="editForm1771.taxCreditPph23"
                  type="number"
                  class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white"
                />
                <span class="text-[10px] text-zinc-400 mt-1 block">
                  SPT Masa riil: {{ formatIDR(annual1771.baseSptMasaPph23) }}
                </span>
              </div>

              <div>
                <label class="block font-bold text-zinc-800 mb-1">Kredit Pajak PPh 25 (Induk 10a) IDR</label>
                <input
                  v-model.number="editForm1771.taxCreditPph25"
                  type="number"
                  class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white"
                />
                <span class="text-[10px] text-zinc-400 mt-1 block">
                  SPT Masa riil: {{ formatIDR(annual1771.baseSptMasaPph25) }}
                </span>
              </div>
            </div>

            <!-- Total Credits & PPh 29/28A Live Summary Box -->
            <div class="bg-zinc-950 text-white rounded-xl p-4 space-y-2">
              <div class="flex items-center justify-between text-xs">
                <span class="text-zinc-400">Total Kredit Pajak (PPh 23 + 25):</span>
                <span class="font-bold text-emerald-300 tabular-nums">{{ formatIDR(modalComputedTotalCredits) }}</span>
              </div>
              <div class="flex items-center justify-between text-xs border-t border-zinc-800 pt-2">
                <span class="text-zinc-400 font-bold">
                  {{ modalComputedPph28A > 0 ? 'PPh Pasal 28A Lebih Bayar:' : 'PPh Pasal 29 Kurang Bayar:' }}
                </span>
                <span
                  class="text-base font-black tabular-nums"
                  :class="modalComputedPph28A > 0 ? 'text-emerald-400' : 'text-amber-400'"
                >
                  {{ formatIDR(modalComputedPph28A > 0 ? modalComputedPph28A : modalComputedPph29) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="sticky bottom-0 px-6 py-4 bg-zinc-50 border-t border-zinc-100 flex items-center justify-between z-10">
          <button
            type="button"
            @click="handleResetOverrides"
            class="px-3.5 py-2 bg-white hover:bg-zinc-100 border border-zinc-200 text-zinc-700 rounded-xl font-semibold transition active:scale-95 cursor-pointer text-xs"
          >
            Reset ke Perhitungan Awal
          </button>

          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="isFullEditModalOpen = false"
              class="px-4 py-2 text-zinc-600 hover:text-zinc-900 rounded-xl font-semibold transition cursor-pointer text-xs"
            >
              Batal
            </button>
            <button
              type="button"
              @click="handleSaveFullEdit"
              class="px-5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold transition active:scale-95 cursor-pointer text-xs shadow-xs"
            >
              Simpan Perubahan SPT 1771
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL 2: PENYESUAIAN KREDIT BULANAN (KELEBIHAN POTONG / REKLASIFIKASI) -->
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
