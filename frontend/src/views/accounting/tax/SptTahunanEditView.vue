<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import {
  ArrowLeft,
  Save,
  RotateCcw,
  Sparkles,
  Info,
  CheckCircle2,
  AlertCircle,
  Calculator,
  SlidersHorizontal,
  Landmark,
  Building2,
  FileText
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await taxStore.fetchTaxRecords()
  initFormValues()
})

const annual1771 = computed(() => taxStore.annualTaxReport1771)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

// Edit Form State
const form = ref({
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

function initFormValues() {
  const current = annual1771.value
  form.value = {
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
}

// Live Computations on Form
const computedNetProfit = computed(() => {
  return Number(form.value.commercialRevenue || 0) - Number(form.value.commercialExpense || 0)
})

const computedPKP = computed(() => {
  return Math.max(
    0,
    computedNetProfit.value +
    Number(form.value.positiveFiscalCorrection || 0) -
    Number(form.value.negativeFiscalCorrection || 0)
  )
})

const computedTaxPayable = computed(() => {
  return Math.round(computedPKP.value * Number(form.value.corporateTaxRate || 0))
})

const computedTotalCredits = computed(() => {
  return Number(form.value.taxCreditPph23 || 0) + Number(form.value.taxCreditPph25 || 0)
})

const computedPph29 = computed(() => {
  return Math.max(0, computedTaxPayable.value - computedTotalCredits.value)
})

const computedPph28A = computed(() => {
  return computedTotalCredits.value > computedTaxPayable.value
    ? computedTotalCredits.value - computedTaxPayable.value
    : 0
})

// Monthly row claim adjustment directly on page
function updateMonthClaim(monthNum: number, taxType: 'PPH_23' | 'PPH_25', original: number, status: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED', adjustedVal?: number) {
  let val = original
  if (status === 'EXCLUDED') val = 0
  else if (status === 'CLAIM_PARTIAL') val = adjustedVal !== undefined ? adjustedVal : Math.round(original * 0.8)

  taxStore.setMonthAdjustment({
    monthNum,
    taxType,
    originalAmount: original,
    adjustedAmount: val,
    claimStatus: status,
    reason: status === 'CLAIM_FULL' ? '' : 'Penyesuaian kelebihan potong pada formulir 1771'
  })

  // Sync with form credits
  form.value.taxCreditPph23 = annual1771.value.taxCreditPph23
  form.value.taxCreditPph25 = annual1771.value.taxCreditPph25
}

function handleMonthAmountInput(monthNum: number, taxType: 'PPH_23' | 'PPH_25', original: number, newVal: number) {
  taxStore.setMonthAdjustment({
    monthNum,
    taxType,
    originalAmount: original,
    adjustedAmount: Number(newVal) || 0,
    claimStatus: Number(newVal) === original ? 'CLAIM_FULL' : (Number(newVal) === 0 ? 'EXCLUDED' : 'CLAIM_PARTIAL'),
    reason: 'Koreksi nominal di formulir 1771'
  })

  form.value.taxCreditPph23 = annual1771.value.taxCreditPph23
  form.value.taxCreditPph25 = annual1771.value.taxCreditPph25
}

function applyPresetOverwithholding() {
  form.value.taxCreditPph23 = Math.max(0, form.value.taxCreditPph23 - 2400000)
  accountingStore.showToast('Simulasi kelebihan potong PPh 23 (-Rp 2.400.000) diterapkan pada formulir', 'info')
}

function handleSaveAll() {
  taxStore.updateAnnual1771Overrides({
    commercialRevenue: Number(form.value.commercialRevenue),
    commercialExpense: Number(form.value.commercialExpense),
    positiveFiscalCorrection: Number(form.value.positiveFiscalCorrection),
    positiveCorrectionNotes: form.value.positiveCorrectionNotes,
    negativeFiscalCorrection: Number(form.value.negativeFiscalCorrection),
    negativeCorrectionNotes: form.value.negativeCorrectionNotes,
    corporateTaxRate: Number(form.value.corporateTaxRate),
    taxCreditPph23: Number(form.value.taxCreditPph23),
    taxCreditPph25: Number(form.value.taxCreditPph25)
  })

  router.push('/accounting/tax/spt-tahunan')
}

function handleResetDefaults() {
  taxStore.resetAnnual1771Overrides()
  taxStore.resetAllAdjustments()
  initFormValues()
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-16 max-w-6xl mx-auto">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <router-link
            to="/accounting/tax/spt-tahunan"
            class="text-xs font-semibold text-zinc-500 hover:text-zinc-900 transition flex items-center gap-1"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke SPT Tahunan 1771</span>
          </router-link>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Formulir Khusus 1771</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-1">
          Penyesuaian & Reklasifikasi SPT Tahunan PPh Badan 1771
        </h1>
        <p class="text-xs text-zinc-500">
          Ubah seluruh komponen SPT Tahunan sesuai kebutuhan fiskal akhir tahun tanpa mengubah transaksi historis bank pada SPT Masa.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="handleResetDefaults"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <RotateCcw class="w-3.5 h-3.5 text-zinc-500" />
          <span>Reset ke Nilai Buku</span>
        </button>

        <button
          type="button"
          @click="handleSaveAll"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan & Kembali ke 1771</span>
        </button>
      </div>
    </header>

    <!-- Info Callout -->
    <div class="bg-amber-50/80 border border-amber-200 rounded-2xl p-4 flex items-start gap-3 text-xs text-amber-900">
      <Info class="w-5 h-5 text-amber-700 shrink-0 mt-0.5" />
      <div>
        <h4 class="font-bold text-amber-950">Prinsip Independensi SPT Tahunan vs SPT Masa</h4>
        <p class="text-amber-800 mt-0.5">
          Perubahan pada formulir ini hanya berlaku untuk SPT Tahunan Formulir 1771 (termasuk jurnal akrual akhir tahun). Seluruh SPT Masa bulanan yang telah terverifikasi kanal perbankan tetap utuh sebagai bukti potong dan mutasi kas bank resmi.
        </p>
      </div>
    </div>

    <!-- Live Floating Summary Bar -->
    <div class="bg-zinc-950 text-white rounded-2xl p-5 shadow-xl border border-zinc-800 space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-800 pb-3">
        <div class="flex items-center gap-2">
          <Calculator class="w-4 h-4 text-amber-400" />
          <span class="text-xs font-bold uppercase tracking-wider text-zinc-300">
            Preview Perhitungan Fiskal Real-Time
          </span>
        </div>

        <div class="flex items-center gap-2">
          <span class="text-xs text-zinc-400">
            {{ computedPph28A > 0 ? 'Status: PPh 28A Lebih Bayar' : 'Status: PPh 29 Kurang Bayar' }}
          </span>
          <span
            class="text-xl font-black tabular-nums"
            :class="computedPph28A > 0 ? 'text-emerald-400' : 'text-amber-400'"
          >
            {{ formatIDR(computedPph28A > 0 ? computedPph28A : computedPph29) }}
          </span>
        </div>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
        <div class="bg-zinc-900/80 p-3 rounded-xl border border-zinc-800/60">
          <span class="text-[10px] text-zinc-400 uppercase font-bold block">1. Laba Komersial</span>
          <span class="text-sm font-bold text-white tabular-nums">{{ formatIDR(computedNetProfit) }}</span>
        </div>
        <div class="bg-zinc-900/80 p-3 rounded-xl border border-zinc-800/60">
          <span class="text-[10px] text-zinc-400 uppercase font-bold block">2. Penghasilan Kena Pajak</span>
          <span class="text-sm font-bold text-amber-300 tabular-nums">{{ formatIDR(computedPKP) }}</span>
        </div>
        <div class="bg-zinc-900/80 p-3 rounded-xl border border-zinc-800/60">
          <span class="text-[10px] text-zinc-400 uppercase font-bold block">3. PPh Terutang</span>
          <span class="text-sm font-bold text-white tabular-nums">{{ formatIDR(computedTaxPayable) }}</span>
        </div>
        <div class="bg-emerald-950/80 p-3 rounded-xl border border-emerald-700/60 text-emerald-100">
          <span class="text-[10px] text-emerald-400 uppercase font-bold block">4. Total Kredit Pajak</span>
          <span class="text-sm font-bold text-emerald-300 tabular-nums">{{ formatIDR(computedTotalCredits) }}</span>
        </div>
      </div>
    </div>

    <!-- Form Section Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs">
      <!-- Section 1: Pendapatan & Beban Komersial -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
        <div class="border-b border-zinc-100 pb-2">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <Building2 class="w-4 h-4 text-zinc-700" />
            <span>1. Laba Bersih Komersial</span>
          </h3>
          <p class="text-[11px] text-zinc-500">Angka bersumber dari Laporan Laba Rugi Komersial klinik.</p>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">
              Peredaran Usaha & Pendapatan Bruto (Akun 4 & 6)
            </label>
            <input
              v-model.number="form.commercialRevenue"
              type="number"
              class="w-full px-3.5 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">
              Beban Operasional & HPP Obat (Akun 5 & 52)
            </label>
            <input
              v-model.number="form.commercialExpense"
              type="number"
              class="w-full px-3.5 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>

          <div class="flex items-center justify-between p-3 bg-zinc-100 rounded-xl font-mono text-xs">
            <span class="text-zinc-600 font-medium">Laba Bersih Komersial:</span>
            <span class="font-bold text-zinc-950 tabular-nums">
              {{ formatIDR(computedNetProfit) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Section 2: Rekonsiliasi Fiskal -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
        <div class="border-b border-zinc-100 pb-2">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <SlidersHorizontal class="w-4 h-4 text-zinc-700" />
            <span>2. Koreksi Fiskal UU PPh</span>
          </h3>
          <p class="text-[11px] text-zinc-500">Penyesuaian beda tetap & temporer antara komersial vs fiskal.</p>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-emerald-800 mb-1">
              Koreksi Fiskal Positif (+) IDR
            </label>
            <input
              v-model.number="form.positiveFiscalCorrection"
              type="number"
              class="w-full px-3.5 py-2.5 bg-emerald-50/40 border border-emerald-300 rounded-xl font-mono text-emerald-950 font-bold focus:outline-none focus:ring-2 focus:ring-emerald-600 focus:bg-white text-xs"
            />
            <input
              v-model="form.positiveCorrectionNotes"
              type="text"
              placeholder="Catatan: Biaya entertainment tanpa daftar nominatif, denda pajak, dll"
              class="w-full mt-1.5 px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-700 text-[11px]"
            />
          </div>

          <div>
            <label class="block font-bold text-rose-800 mb-1">
              Koreksi Fiskal Negatif (-) IDR
            </label>
            <input
              v-model.number="form.negativeFiscalCorrection"
              type="number"
              class="w-full px-3.5 py-2.5 bg-rose-50/40 border border-rose-300 rounded-xl font-mono text-rose-950 font-bold focus:outline-none focus:ring-2 focus:ring-rose-600 focus:bg-white text-xs"
            />
            <input
              v-model="form.negativeCorrectionNotes"
              type="text"
              placeholder="Catatan: Pendapatan bunga bank dipotong PPh Final, dll"
              class="w-full mt-1.5 px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-700 text-[11px]"
            />
          </div>

          <div class="flex items-center justify-between p-3 bg-amber-50 rounded-xl font-mono text-xs border border-amber-200">
            <span class="text-amber-900 font-bold">Penghasilan Kena Pajak (PKP):</span>
            <span class="font-black text-amber-950 tabular-nums">
              {{ formatIDR(computedPKP) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Section 3: Tarif PPh Badan -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
        <div class="border-b border-zinc-100 pb-2">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <Landmark class="w-4 h-4 text-zinc-700" />
            <span>3. Tarif PPh Badan Efektif</span>
          </h3>
          <p class="text-[11px] text-zinc-500">Pilih skema tarif pajak fasilitas UU HPP atau tarif umum.</p>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-zinc-700 mb-1.5">
              Skema Tarif Pajak Penghasilan Badan
            </label>
            <div class="grid grid-cols-2 gap-2 mb-2">
              <button
                type="button"
                @click="form.corporateTaxRate = 0.11"
                class="px-3 py-2 rounded-xl text-center border font-bold text-[11px] transition cursor-pointer"
                :class="form.corporateTaxRate === 0.11 ? 'bg-zinc-950 text-white border-zinc-950' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                11% (Fasilitas 31E UU HPP)
              </button>
              <button
                type="button"
                @click="form.corporateTaxRate = 0.22"
                class="px-3 py-2 rounded-xl text-center border font-bold text-[11px] transition cursor-pointer"
                :class="form.corporateTaxRate === 0.22 ? 'bg-zinc-950 text-white border-zinc-950' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                22% (Tarif Umum Badan)
              </button>
            </div>
            <div class="relative">
              <input
                v-model.number="form.corporateTaxRate"
                type="number"
                step="0.01"
                min="0"
                max="1"
                class="w-full px-3.5 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
              />
            </div>
          </div>

          <div class="flex items-center justify-between p-3 bg-zinc-100 rounded-xl font-mono text-xs">
            <span class="text-zinc-600 font-medium">PPh Badan Terutang (PKP x Tarif):</span>
            <span class="font-bold text-zinc-950 tabular-nums">
              {{ formatIDR(computedTaxPayable) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Section 4: Kredit Pajak Diakui di 1771 -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-2">
          <div>
            <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span>4. Kredit Pajak Diakui di 1771</span>
            </h3>
            <p class="text-[11px] text-zinc-500">Nilai yang diakui sebagai pengurang PPh Badan.</p>
          </div>
          <button
            type="button"
            @click="applyPresetOverwithholding"
            class="px-2.5 py-1 bg-amber-50 hover:bg-amber-100 text-amber-800 border border-amber-200 rounded-lg text-[10px] font-bold transition cursor-pointer"
          >
            Simulasi Kelebihan Potong
          </button>
        </div>

        <div class="space-y-3">
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="font-bold text-zinc-800">Kredit Pajak PPh 23 (Lampiran III) IDR</label>
              <span class="text-[10px] text-zinc-400">SPT Masa Asli: {{ formatIDR(annual1771.baseSptMasaPph23) }}</span>
            </div>
            <input
              v-model.number="form.taxCreditPph23"
              type="number"
              class="w-full px-3.5 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>

          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="font-bold text-zinc-800">Kredit Pajak PPh 25 Angsuran (Induk 10a) IDR</label>
              <span class="text-[10px] text-zinc-400">SPT Masa Asli: {{ formatIDR(annual1771.baseSptMasaPph25) }}</span>
            </div>
            <input
              v-model.number="form.taxCreditPph25"
              type="number"
              class="w-full px-3.5 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-2 focus:ring-zinc-950 focus:bg-white text-xs"
            />
          </div>

          <div class="flex items-center justify-between p-3 bg-emerald-50 rounded-xl font-mono text-xs border border-emerald-200">
            <span class="text-emerald-900 font-bold">Total Kredit Diakui (23 + 25):</span>
            <span class="font-black text-emerald-950 tabular-nums">
              {{ formatIDR(computedTotalCredits) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Section 5: Detail Penyesuaian Bulanan (Tabel Langsung di Halaman Tanpa Modal) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <FileText class="w-4 h-4 text-zinc-700" />
          <h3 class="text-sm font-bold text-zinc-900">
            Detail Penyesuaian Status Klaim Kredit Pajak Per Bulan (Tanpa Mengubah SPT Masa)
          </h3>
        </div>
        <p class="text-xs text-zinc-500 mt-0.5">
          Atur status klaim masing-masing bulan jika terdapat bukti potong yang cacat atau kelebihan bayar. Nilai pada SPT Masa tetap utuh tersimpan di basis data.
        </p>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-3 py-2.5">Masa Pajak</th>
              <th class="px-3 py-2.5 text-right">SPT Masa PPh 23</th>
              <th class="px-3 py-2.5 text-right">SPT Masa PPh 25</th>
              <th class="px-3 py-2.5 text-center">Status Pengakuan 1771</th>
              <th class="px-3 py-2.5 text-right w-44">Nilai Diakui di 1771 (IDR)</th>
              <th class="px-3 py-2.5 text-right">Selisih Koreksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr v-for="row in sptMatrix" :key="row.monthNum" class="hover:bg-zinc-50/50">
              <td class="px-3 py-2.5 font-bold font-sans text-zinc-800">
                {{ row.monthName }} (Masa {{ row.monthNum }})
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-zinc-700">
                {{ row.pph23Original > 0 ? formatIDR(row.pph23Original, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-emerald-700 font-semibold">
                {{ row.pph25Original > 0 ? formatIDR(row.pph25Original, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-center">
                <div v-if="row.pph23Original > 0 || row.pph25Original > 0" class="inline-flex rounded-lg bg-zinc-100 p-0.5 text-[10px] font-bold">
                  <button
                    type="button"
                    @click="updateMonthClaim(row.monthNum, row.pph23Original > 0 ? 'PPH_23' : 'PPH_25', row.pph23Original > 0 ? row.pph23Original : row.pph25Original, 'CLAIM_FULL')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="(!row.isAdjusted || row.pph23Status === 'CLAIM_FULL') ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    100%
                  </button>
                  <button
                    type="button"
                    @click="updateMonthClaim(row.monthNum, row.pph23Original > 0 ? 'PPH_23' : 'PPH_25', row.pph23Original > 0 ? row.pph23Original : row.pph25Original, 'CLAIM_PARTIAL')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="row.pph23Status === 'CLAIM_PARTIAL' ? 'bg-amber-500 text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    Sebagian
                  </button>
                  <button
                    type="button"
                    @click="updateMonthClaim(row.monthNum, row.pph23Original > 0 ? 'PPH_23' : 'PPH_25', row.pph23Original > 0 ? row.pph23Original : row.pph25Original, 'EXCLUDED')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="row.pph23Status === 'EXCLUDED' ? 'bg-rose-600 text-white shadow-2xs' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    Tidak
                  </button>
                </div>
                <span v-else class="text-zinc-300">-</span>
              </td>
              <td class="px-3 py-2.5 text-right">
                <input
                  v-if="row.pph23Original > 0 || row.pph25Original > 0"
                  type="number"
                  :value="row.pph23Original > 0 ? row.pph23Amount : row.pph25Amount"
                  @input="e => handleMonthAmountInput(row.monthNum, row.pph23Original > 0 ? 'PPH_23' : 'PPH_25', row.pph23Original > 0 ? row.pph23Original : row.pph25Original, Number((e.target as HTMLInputElement).value))"
                  class="w-full px-2 py-1 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-zinc-950 text-xs focus:bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
                />
                <span v-else class="text-zinc-300">-</span>
              </td>
              <td class="px-3 py-2.5 text-right font-mono text-xs tabular-nums">
                <span
                  v-if="row.isAdjusted"
                  :class="row.adjustmentDiff < 0 ? 'text-rose-600 font-bold' : 'text-emerald-600 font-bold'"
                >
                  {{ row.adjustmentDiff >= 0 ? '+' : '' }}{{ formatIDR(row.adjustmentDiff, false) }}
                </span>
                <span v-else class="text-zinc-400">0</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Bottom Actions Bar -->
    <div class="flex items-center justify-between pt-4 border-t border-zinc-200">
      <router-link
        to="/accounting/tax/spt-tahunan"
        class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition"
      >
        Batal & Kembali
      </router-link>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="handleResetDefaults"
          class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition cursor-pointer"
        >
          Reset ke Nilai Buku
        </button>

        <button
          type="button"
          @click="handleSaveAll"
          class="inline-flex items-center gap-1.5 px-6 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan Perubahan Formulir 1771</span>
        </button>
      </div>
    </div>
  </div>
</template>
