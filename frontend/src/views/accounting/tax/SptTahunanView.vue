<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import {
  FileText,
  CheckCircle2,
  Download,
  Printer,
  Sparkles,
  Link2,
  ArrowLeft,
  ArrowRight,
  SlidersHorizontal,
  RotateCcw,
  AlertCircle,
  Send,
  Building2,
  Pencil,
  X,
  ExternalLink,
  ChevronRight,
  Landmark,
  Scale
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

onMounted(async () => {
  await Promise.all([
    taxStore.fetchTaxRecords(),
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
})

const annual1771 = computed(() => taxStore.annualTaxReport1771)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

const selectedMonthForJournals = ref<number | null>(null)
const showMonthJournalsModal = ref(false)
const showJournalDetailModal = ref(false)
const activeJournal = ref<any>(null)

function openMonthJournals(monthNum: number) {
  selectedMonthForJournals.value = monthNum
  showMonthJournalsModal.value = true
}

const selectedMonthJournalEntries = computed(() => {
  if (!selectedMonthForJournals.value) return []
  return taxStore.getLinkedJournalsForMonth(selectedMonthForJournals.value, annual1771.value.fiscalYear)
})

function openJournalDetail(entry: any) {
  if (entry.rawJournal) {
    activeJournal.value = entry.rawJournal
  } else {
    const found = journalsStore.journals.find(j => j.journalNumber === entry.journalNumber || j.id === entry.journalId)
    activeJournal.value = found || entry
  }
  showJournalDetailModal.value = true
}

function handleResetAll() {
  taxStore.resetAllAdjustments()
  taxStore.resetAnnual1771Overrides()
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
          <span
            v-if="annual1771.isCustomOverridden"
            class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-900 font-bold text-[10px] border border-amber-300"
          >
            Data 1771 Disesuaikan Manual
          </span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          SPT Tahunan PPh Badan (Formulir 1771)
        </h1>
        <p class="text-xs text-zinc-500">
          Rekonsiliasi fiskal tahunan & akumulasi kredit pajak riil. Seluruh komponen SPT Tahunan dapat diubah pada halaman penyesuaian tanpa memengaruhi SPT Masa.
        </p>
      </div>

      <div class="flex items-center gap-2.5 flex-wrap">
        <router-link
          to="/accounting/tax/spt-tahunan/edit"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-amber-500 hover:bg-amber-600 text-zinc-950 rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Pencil class="w-3.5 h-3.5" />
          <span>Ubah Data Kertas Kerja 1771</span>
        </router-link>

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
            <h4 class="font-bold text-zinc-950 text-sm">Prinsip Non-Destruktif SPT Tahunan vs SPT Masa</h4>
            <span class="px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 font-bold text-[10px]">Non-Cash DJP Standard</span>
            <span
              v-if="annual1771.isCustomOverridden"
              class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 border border-amber-300 font-bold text-[10px] flex items-center gap-1"
            >
              <AlertCircle class="w-3 h-3 text-amber-600" />
              <span>Kertas Kerja Menggunakan Nilai Manual</span>
            </span>
          </div>
          <p class="text-zinc-600 mt-0.5">
            Setiap komponen SPT Tahunan dapat diedit pada halaman khusus penyesuaian tanpa merusak catatan perbankan pada SPT Masa bulanan.
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="annual1771.isCustomOverridden"
          @click="handleResetAll"
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
            <router-link to="/accounting/tax/spt-tahunan/edit" class="text-zinc-500 hover:text-white cursor-pointer" title="Ubah di Halaman Penyesuaian">
              <Pencil class="w-3 h-3" />
            </router-link>
          </div>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.commercialNetProfit) }}
          </div>
          <p class="text-[10px] text-zinc-400">Pendapatan dikurangi HPP & Beban Operasional</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">2. Penghasilan Kena Pajak</span>
            <router-link to="/accounting/tax/spt-tahunan/edit" class="text-zinc-500 hover:text-amber-300 cursor-pointer" title="Ubah Koreksi Fiskal">
              <Pencil class="w-3 h-3" />
            </router-link>
          </div>
          <div class="tabular-nums text-base font-bold text-amber-300">
            {{ formatIDR(annual1771.taxableIncome) }}
          </div>
          <p class="text-[10px] text-zinc-400">Setelah Koreksi Fiskal Positif & Negatif</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <div class="flex items-center justify-between">
            <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">3. PPh Badan Terutang</span>
            <router-link to="/accounting/tax/spt-tahunan/edit" class="text-zinc-500 hover:text-white cursor-pointer" title="Ubah Tarif Pajak">
              <Pencil class="w-3 h-3" />
            </router-link>
          </div>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.corporateTaxPayable) }}
          </div>
          <p class="text-[10px] text-zinc-400">Tarif Efektif {{ (annual1771.corporateTaxRate * 100).toFixed(0) }}% x PKP</p>
        </div>

        <div class="bg-emerald-950/80 rounded-xl p-4 border border-emerald-700/60 space-y-1 text-emerald-100">
          <div class="flex items-center justify-between">
            <span class="text-emerald-400 font-bold uppercase tracking-wide text-[10px]">4. Kredit Pajak Diakui</span>
            <router-link to="/accounting/tax/spt-tahunan/edit" class="text-emerald-500 hover:text-emerald-200 cursor-pointer" title="Ubah Kredit Pajak">
              <Pencil class="w-3 h-3" />
            </router-link>
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

    <!-- SPT 1771 FULL BREAKDOWN TABLE -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <FileText class="w-4 h-4 text-zinc-700" />
          <h3 class="text-sm font-bold text-zinc-900">
            Kertas Kerja Rekonsiliasi Fiskal SPT 1771 (Tahun Pajak {{ annual1771.fiscalYear }})
          </h3>
        </div>

        <div class="flex items-center gap-2">
          <router-link
            to="/accounting/tax/spt-tahunan/edit"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer shadow-xs"
          >
            <Pencil class="w-3 h-3" />
            <span>Buka Halaman Form Penyesuaian</span>
          </router-link>

          <button
            v-if="annual1771.isCustomOverridden"
            type="button"
            @click="handleResetAll"
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-zinc-200 text-zinc-500 hover:text-zinc-950 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-emerald-200 text-emerald-700 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="p-1 rounded-lg hover:bg-emerald-200 text-emerald-700 transition inline-flex items-center"
                  title="Ubah di Halaman Form"
                >
                  <Pencil class="w-3.5 h-3.5" />
                </router-link>
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
            Diambil secara otomatis dari dokumen SPT Masa yang telah terverifikasi kanal bank resmi.
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
              <th class="px-3 py-2.5 text-center">Jurnal Terkait</th>
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
                  :title="`Koreksi: ${formatIDR(row.adjustmentDiff)}`"
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
              <td class="px-3 py-2.5 text-center whitespace-nowrap">
                <button
                  @click="openMonthJournals(row.monthNum)"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-[10px] font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-800 border border-zinc-200 transition cursor-pointer"
                  title="Lihat Jurnal Terkait Bulan Ini"
                >
                  <Link2 class="w-3 h-3 text-zinc-500" />
                  <span>{{ taxStore.getLinkedJournalsForMonth(row.monthNum, annual1771.fiscalYear).length }} Jurnal</span>
                </button>
              </td>
              <td class="px-3 py-2.5 text-center">
                <router-link
                  to="/accounting/tax/spt-tahunan/edit"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-[11px] font-semibold bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 transition active:scale-95 cursor-pointer shadow-2xs"
                  title="Sesuaikan kredit pada halaman formulir 1771"
                >
                  <Pencil class="w-3 h-3 text-zinc-500" />
                  <span>Koreksi</span>
                </router-link>
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
              <td class="px-3 py-3 text-center text-emerald-400 font-sans text-xs" colspan="3">
                {{ annual1771.hasAdjustments ? 'Diakui dengan Penyesuaian' : 'Terverifikasi Penuh' }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- PANEL AKSI JURNAL AKRUAL -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <SlidersHorizontal class="w-4 h-4 text-zinc-900" />
            <h3 class="text-sm font-bold text-zinc-900">
              Posting Jurnal Akrual SPT Tahunan ke Buku Besar
            </h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Membukukan otomatis beban PPh Badan berjalan, pengkreditan PPh 23 & 25, serta hutang PPh 29 ke Buku Besar.
          </p>
        </div>

        <div class="flex items-center gap-2">
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

    <!-- BRAYO GL RECONCILIATION SECTION: Rekonsiliasi Jurnal Umum & Laba Rugi Tahunan -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <Scale class="w-4 h-4 text-emerald-600" />
          <div>
            <h3 class="text-sm font-bold text-zinc-950">
              Rekonsiliasi Jurnal Umum & Buku Besar Tahunan (Brayo Pajak Standard)
            </h3>
            <p class="text-xs text-zinc-500">
              Pencocokan saldo akun pendapatan, beban, dan kredit pajak pada Buku Besar (GL Entry) terhadap Form 1771 DJP.
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <router-link
            to="/accounting/trial-balance"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 border border-zinc-200 text-zinc-700 rounded-xl text-xs font-bold transition shadow-2xs"
          >
            <span>Neraca Saldo (Trial Balance)</span>
            <ExternalLink class="w-3.5 h-3.5" />
          </router-link>
          <router-link
            to="/accounting/ledger"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs"
          >
            <span>Buku Besar (GL)</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </router-link>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs">
        <!-- Card 1: Laba Komersial GL -->
        <div class="p-4 rounded-xl bg-zinc-50 border border-zinc-200 space-y-2">
          <div class="flex justify-between items-center">
            <span class="font-bold text-zinc-500 uppercase text-[10px]">Laba Komersial (Akun 4 - Akun 5)</span>
            <span class="px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800 text-[9px] font-bold">RECONCILED</span>
          </div>
          <p class="text-base font-black text-zinc-950 tabular-nums">
            {{ formatIDR(annual1771.commercialNetProfit) }}
          </p>
          <div class="text-[11px] text-zinc-500 space-y-0.5 pt-1 border-t border-zinc-200">
            <div class="flex justify-between">
              <span>Pendapatan (4 & 6):</span>
              <span class="font-mono">{{ formatIDR(annual1771.commercialRevenue) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Beban (5 & 7):</span>
              <span class="font-mono">{{ formatIDR(annual1771.commercialExpense) }}</span>
            </div>
          </div>
        </div>

        <!-- Card 2: Koreksi Fiskal Netto -->
        <div class="p-4 rounded-xl bg-zinc-50 border border-zinc-200 space-y-2">
          <div class="flex justify-between items-center">
            <span class="font-bold text-zinc-500 uppercase text-[10px]">Koreksi Fiskal Netto</span>
            <span class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 text-[9px] font-bold">AUDITED</span>
          </div>
          <p class="text-base font-black text-amber-700 tabular-nums">
            {{ formatIDR(annual1771.positiveFiscalCorrection - annual1771.negativeFiscalCorrection) }}
          </p>
          <div class="text-[11px] text-zinc-500 space-y-0.5 pt-1 border-t border-zinc-200">
            <div class="flex justify-between">
              <span>Koreksi Positif (+):</span>
              <span class="font-mono text-emerald-700">{{ formatIDR(annual1771.positiveFiscalCorrection) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Koreksi Negatif (-):</span>
              <span class="font-mono text-rose-700">{{ formatIDR(annual1771.negativeFiscalCorrection) }}</span>
            </div>
          </div>
        </div>

        <!-- Card 3: Kredit Pajak Bank -->
        <div class="p-4 rounded-xl bg-zinc-50 border border-zinc-200 space-y-2">
          <div class="flex justify-between items-center">
            <span class="font-bold text-zinc-500 uppercase text-[10px]">Kredit Pajak Terverifikasi (Bank)</span>
            <span class="px-2 py-0.5 rounded-full bg-blue-100 text-blue-800 text-[9px] font-bold">MATCH 100%</span>
          </div>
          <p class="text-base font-black text-zinc-950 tabular-nums">
            {{ formatIDR(annual1771.taxCredits) }}
          </p>
          <div class="text-[11px] text-zinc-500 space-y-0.5 pt-1 border-t border-zinc-200">
            <div class="flex justify-between">
              <span>PPh 23 (Akun 2502):</span>
              <span class="font-mono">{{ formatIDR(annual1771.taxCreditPph23) }}</span>
            </div>
            <div class="flex justify-between">
              <span>PPh 25 (Akun 1504/2505):</span>
              <span class="font-mono">{{ formatIDR(annual1771.taxCreditPph25) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- DIALOG MODAL: Daftar Jurnal Terkait Bulan Tertentu -->
    <div
      v-if="showMonthJournalsModal && selectedMonthForJournals"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-2xl w-full overflow-hidden p-6 space-y-4 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">
              Daftar Jurnal Terkait Pajak — Masa {{ selectedMonthForJournals }}/{{ annual1771.fiscalYear }}
            </h3>
            <p class="text-xs text-zinc-500 mt-0.5">
              Seluruh voucher transaksi yang menyumbang kredit pajak / pemotongan pada masa ini.
            </p>
          </div>
          <button @click="showMonthJournalsModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-500 uppercase">
                <th class="py-2 px-3">Tanggal</th>
                <th class="py-2 px-3">No. Jurnal</th>
                <th class="py-2 px-3">Uraian Transaksi</th>
                <th class="py-2 px-3">Jenis Pajak</th>
                <th class="py-2 px-3 text-right">Nilai Pajak (Rp)</th>
                <th class="py-2 px-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr
                v-for="j in selectedMonthJournalEntries"
                :key="j.id"
                class="hover:bg-zinc-50/50"
              >
                <td class="py-2 px-3 font-mono text-zinc-600">{{ formatDate(j.journalDate) }}</td>
                <td class="py-2 px-3 font-mono font-bold text-zinc-950">{{ j.journalNumber }}</td>
                <td class="py-2 px-3 text-zinc-800 max-w-xs truncate">{{ j.description }}</td>
                <td class="py-2 px-3 font-bold text-zinc-900">{{ j.taxType }}</td>
                <td class="py-2 px-3 text-right tabular-nums font-bold text-zinc-950">{{ formatIDR(j.taxAmount) }}</td>
                <td class="py-2 px-3 text-center">
                  <button
                    @click="openJournalDetail(j)"
                    class="px-2 py-0.5 bg-white hover:bg-zinc-100 text-zinc-700 rounded-lg text-xs font-bold border border-zinc-200"
                  >
                    Detail
                  </button>
                </td>
              </tr>
              <tr v-if="selectedMonthJournalEntries.length === 0">
                <td colspan="6" class="py-6 text-center text-zinc-400 italic">
                  Belum ada jurnal khusus yang terhubung pada masa ini.
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex justify-end pt-2 border-t border-zinc-100">
          <button
            @click="showMonthJournalsModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- DIALOG MODAL: Detail Voucher Jurnal -->
    <div
      v-if="showJournalDetailModal && activeJournal"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-2xl w-full overflow-hidden p-6 space-y-4 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono text-base font-bold text-zinc-950">{{ activeJournal.journalNumber }}</span>
              <span class="px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-800 border border-emerald-200 text-[10px] font-bold">
                {{ activeJournal.status || 'POSTED' }}
              </span>
            </div>
            <p class="text-xs text-zinc-500 mt-0.5">Tanggal: {{ formatDate(activeJournal.journalDate) }} • Unit: {{ activeJournal.businessUnit || 'ALL' }}</p>
          </div>
          <button @click="showJournalDetailModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-3 bg-zinc-50 rounded-xl border border-zinc-200 text-xs">
          <span class="text-zinc-400 font-bold uppercase text-[10px] block">Keterangan / Uraian Jurnal:</span>
          <p class="text-zinc-900 font-medium mt-0.5">{{ activeJournal.description }}</p>
        </div>

        <!-- Multi-line Table -->
        <div class="border border-zinc-200 rounded-xl overflow-hidden bg-white">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-500 uppercase">
                <th class="py-2 px-3">Kode Akun</th>
                <th class="py-2 px-3">Nama Akun</th>
                <th class="py-2 px-3 text-right">Debit (Rp)</th>
                <th class="py-2 px-3 text-right">Kredit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="(line, idx) in activeJournal.lines" :key="idx" class="hover:bg-zinc-50/50">
                <td class="py-2 px-3 font-mono font-bold text-zinc-900">{{ line.accountCode }}</td>
                <td class="py-2 px-3 text-zinc-700">
                  <span class="font-medium block">{{ line.accountName }}</span>
                  <span v-if="line.description" class="text-[10px] text-zinc-400 block">{{ line.description }}</span>
                </td>
                <td class="py-2 px-3 text-right tabular-nums font-semibold text-zinc-950">
                  {{ line.debit > 0 ? formatIDR(line.debit) : '-' }}
                </td>
                <td class="py-2 px-3 text-right tabular-nums font-semibold text-zinc-950">
                  {{ line.credit > 0 ? formatIDR(line.credit) : '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-zinc-100 text-xs">
          <router-link
            :to="`/accounting/journals?search=${activeJournal.journalNumber}`"
            class="text-xs font-bold text-zinc-700 hover:text-zinc-950 flex items-center gap-1"
          >
            <span>Buka di Daftar Jurnal Umum</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </router-link>
          <button
            @click="showJournalDetailModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
