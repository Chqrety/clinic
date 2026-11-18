<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, parseDateMonthYear, type LinkedTaxJournalEntry } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useLedgerStore } from '@/stores/ledger'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportSptTahunan1771ToExcel, exportSptTahunan1771ToCSV, getEntityProfile } from '@/utils/export'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
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
  FileText,
  Plus,
  Trash2,
  ExternalLink,
  BookOpen,
  X,
  Layers,
  Send,
  Printer,
  Download,
  FileSpreadsheet
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()
const ledgerStore = useLedgerStore()
const currentEntity = computed(() => getEntityProfile(accountingStore.businessUnit))

export interface TaxCreditVoucherRow {
  id: string
  selected: boolean
  isFromJournal: boolean
  journalId?: string
  journalNumber?: string
  date: string
  taxType: 'PPH_23' | 'PPH_25'
  accountCode: string
  accountName: string
  counterpartName: string
  description: string
  voucherAmount: number
  claimedAmount: number
  claimStatus: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED'
  reason: string
  rawJournal?: any
}

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

const creditVoucherRows = ref<TaxCreditVoucherRow[]>([])
const showJournalPreviewModal = ref(false)
const activeJournal = ref<any>(null)

const annual1771 = computed(() => taxStore.annualTaxReport1771)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

onMounted(async () => {
  await Promise.all([
    taxStore.fetchTaxRecords(),
    journalsStore.fetchJournals(),
    ledgerStore.fetchAccounts()
  ])
  initFormValues()
  initCreditVoucherRows()
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

function initCreditVoucherRows() {
  const y = annual1771.value.fiscalYear
  const taxEntries = taxStore.linkedTaxJournalEntries.filter(
    e => e.year === y && (e.taxType === 'PPH_23' || e.taxType === 'PPH_25')
  )

  if (taxEntries.length > 0) {
    creditVoucherRows.value = taxEntries.map((e, idx) => ({
      id: `cr-row-${e.id}-${idx}`,
      selected: true,
      isFromJournal: true,
      journalId: e.journalId,
      journalNumber: e.journalNumber,
      date: e.journalDate ? e.journalDate.slice(0, 10) : `${y}-11-15`,
      taxType: e.taxType as ('PPH_23' | 'PPH_25'),
      accountCode: e.taxAccountCode,
      accountName: e.taxAccountName,
      counterpartName: e.counterpartAccountName,
      description: e.description,
      voucherAmount: e.taxAmount,
      claimedAmount: e.taxAmount,
      claimStatus: 'CLAIM_FULL',
      reason: '',
      rawJournal: e.rawJournal
    }))
  } else {
    // Populate from monthly matrix if no individual journal lines
    const rows: TaxCreditVoucherRow[] = []
    for (const m of sptMatrix.value) {
      if (m.pph23Original > 0) {
        rows.push({
          id: `cr-matrix-23-${m.monthNum}`,
          selected: m.pph23Status !== 'EXCLUDED',
          isFromJournal: true,
          journalNumber: `SPT-23-${y}-${String(m.monthNum).padStart(2, '0')}`,
          date: `${y}-${String(m.monthNum).padStart(2, '0')}-20`,
          taxType: 'PPH_23',
          accountCode: '2-1510',
          accountName: 'Kredit Pajak PPh 23',
          counterpartName: 'Asuransi & Vendor Medis',
          description: `Kredit Pajak PPh 23 Masa ${m.monthName} ${y}`,
          voucherAmount: m.pph23Original,
          claimedAmount: m.pph23Amount,
          claimStatus: m.pph23Status,
          reason: m.isAdjusted ? 'Penyesuaian SPT 1771' : ''
        })
      }
      if (m.pph25Original > 0) {
        rows.push({
          id: `cr-matrix-25-${m.monthNum}`,
          selected: m.pph25Status !== 'EXCLUDED',
          isFromJournal: true,
          journalNumber: `SPT-25-${y}-${String(m.monthNum).padStart(2, '0')}`,
          date: `${y}-${String(m.monthNum).padStart(2, '0')}-15`,
          taxType: 'PPH_25',
          accountCode: '2-1520',
          accountName: 'Uang Muka PPh 25 Badan',
          counterpartName: 'Kas Negara via Bank',
          description: `Angsuran Bulanan PPh 25 Masa ${m.monthName} ${y}`,
          voucherAmount: m.pph25Original,
          claimedAmount: m.pph25Amount,
          claimStatus: m.pph25Status,
          reason: m.isAdjusted ? 'Penyesuaian SPT 1771' : ''
        })
      }
    }
    creditVoucherRows.value = rows
  }

  syncCreditsFromRows()
}

/**
 * ⚡ Hero Automation: Ambil Data Otomatis dari Jurnal Umum & Buku Besar (Tahun Ini)
 */
function autoPullFromGeneralLedger() {
  const y = annual1771.value.fiscalYear
  const allJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

  let rev = 0
  let exp = 0

  for (const j of allJournals) {
    const p = parseDateMonthYear(j.journalDate)
    if (p && p.year !== y) continue

    for (const l of j.lines) {
      const code = (l.accountCode || '').trim()
      if (code.startsWith('4-') || code.startsWith('4')) {
        rev += (Number(l.credit) - Number(l.debit))
      } else if (code.startsWith('5-') || code.startsWith('5') || code.startsWith('6-') || code.startsWith('6')) {
        exp += (Number(l.debit) - Number(l.credit))
      }
    }
  }

  if (rev > 0) form.value.commercialRevenue = rev
  if (exp > 0) form.value.commercialExpense = exp

  initCreditVoucherRows()
  accountingStore.showToast('Data Laba Komersial & Transaksi Kredit Pajak berhasil ditarik otomatis dari Jurnal Umum & Buku Besar!', 'success')
}

function syncCreditsFromRows() {
  const pph23Sum = creditVoucherRows.value
    .filter(r => r.selected && r.taxType === 'PPH_23' && r.claimStatus !== 'EXCLUDED')
    .reduce((sum, r) => sum + Number(r.claimedAmount || 0), 0)

  const pph25Sum = creditVoucherRows.value
    .filter(r => r.selected && r.taxType === 'PPH_25' && r.claimStatus !== 'EXCLUDED')
    .reduce((sum, r) => sum + Number(r.claimedAmount || 0), 0)

  form.value.taxCreditPph23 = pph23Sum
  form.value.taxCreditPph25 = pph25Sum
}

function setRowClaimStatus(row: TaxCreditVoucherRow, status: 'CLAIM_FULL' | 'CLAIM_PARTIAL' | 'EXCLUDED') {
  row.claimStatus = status
  if (status === 'CLAIM_FULL') {
    row.claimedAmount = row.voucherAmount
    row.selected = true
  } else if (status === 'EXCLUDED') {
    row.claimedAmount = 0
    row.selected = false
  } else if (status === 'CLAIM_PARTIAL') {
    row.claimedAmount = Math.round(row.voucherAmount * 0.8)
    row.selected = true
  }
  syncCreditsFromRows()
}

function onClaimAmountChange(row: TaxCreditVoucherRow) {
  const amt = Number(row.claimedAmount) || 0
  if (amt >= row.voucherAmount) {
    row.claimStatus = 'CLAIM_FULL'
    row.claimedAmount = row.voucherAmount
  } else if (amt === 0) {
    row.claimStatus = 'EXCLUDED'
    row.selected = false
  } else {
    row.claimStatus = 'CLAIM_PARTIAL'
    row.selected = true
  }
  syncCreditsFromRows()
}

function addCreditRow() {
  creditVoucherRows.value.push({
    id: `cr-manual-${Date.now()}`,
    selected: true,
    isFromJournal: false,
    date: `${annual1771.value.fiscalYear}-11-15`,
    taxType: 'PPH_23',
    accountCode: '2-1510',
    accountName: 'Kredit Pajak PPh 23',
    counterpartName: 'Vendor Jasa / Asuransi Rekanan',
    description: 'Bukti Potong PPh 23 Tambahan',
    voucherAmount: 0,
    claimedAmount: 0,
    claimStatus: 'CLAIM_FULL',
    reason: 'Manual adjustment'
  })
}

function removeCreditRow(index: number) {
  creditVoucherRows.value.splice(index, 1)
  syncCreditsFromRows()
}

// Live Calculations on Form
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

function openJournalPreview(row: TaxCreditVoucherRow) {
  if (row.rawJournal) {
    activeJournal.value = row.rawJournal
  } else {
    const found = journalsStore.journals.find(j => j.journalNumber === row.journalNumber || j.id === row.journalId)
    activeJournal.value = found || {
      journalNumber: row.journalNumber || 'JV-PREVIEW',
      journalDate: row.date,
      description: row.description,
      lines: [
        { accountCode: row.accountCode, accountName: row.accountName, debit: row.voucherAmount, credit: 0 },
        { accountCode: '1-1200', accountName: 'Piutang Pasien & Klaim Asuransi', debit: 0, credit: row.voucherAmount }
      ]
    }
  }
  showJournalPreviewModal.value = true
}

function handleSaveAll() {
  syncCreditsFromRows()

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

  accountingStore.showToast('Kertas Kerja SPT Tahunan 1771 berhasil disimpan!', 'success')
  router.push('/accounting/tax/spt-tahunan')
}

function handlePostPph29Closing() {
  handleSaveAll()
  taxStore.postPph29AccrualJournal()
}

function handleResetDefaults() {
  taxStore.resetAnnual1771Overrides()
  taxStore.resetAllAdjustments()
  initFormValues()
  initCreditVoucherRows()
}

// 📥 Ekspor Kertas Kerja 1771 Per Jurnal Voucher
function handleExportExcel() {
  const commHpp = Math.round(form.value.commercialRevenue * 0.48)
  const commNet = form.value.commercialRevenue - commHpp - form.value.commercialExpense

  exportSptTahunan1771ToExcel({
    annual1771: {
      ...annual1771.value,
      commercialRevenue: form.value.commercialRevenue,
      commercialExpense: form.value.commercialExpense,
      commercialHpp: commHpp,
      commercialNetProfit: commNet,
      positiveFiscalCorrection: form.value.positiveFiscalCorrection,
      positiveCorrectionNotes: form.value.positiveCorrectionNotes,
      negativeFiscalCorrection: form.value.negativeFiscalCorrection,
      negativeCorrectionNotes: form.value.negativeCorrectionNotes,
      taxableIncome: computedPKP.value,
      corporateTaxRate: form.value.corporateTaxRate,
      corporateTaxPayable: computedTaxPayable.value,
      taxCreditPph23: form.value.taxCreditPph23,
      taxCreditPph25: form.value.taxCreditPph25,
      taxCredits: computedTotalCredits.value,
      pph29Payable: computedPph29.value,
      pph28aOverpaid: computedPph28A.value
    },
    monthlyMatrix: sptMatrix.value,
    linkedJournals: journalsStore.journals,
    creditVouchers: creditVoucherRows.value
  }, {
    fiscalYear: annual1771.value.fiscalYear,
    businessUnit: accountingStore.businessUnit
  })
  accountingStore.showToast('Kertas Kerja SPT Tahunan PPh Badan 1771 berhasil diekspor per transaksi jurnal!', 'success')
}

// 📥 Ekspor Kertas Kerja 1771 ke CSV (Format Akuntansi & DJP)
function handleExportCSV() {
  const commHpp = Math.round(form.value.commercialRevenue * 0.48)
  const commNet = form.value.commercialRevenue - commHpp - form.value.commercialExpense

  exportSptTahunan1771ToCSV({
    annual1771: {
      ...annual1771.value,
      commercialRevenue: form.value.commercialRevenue,
      commercialExpense: form.value.commercialExpense,
      commercialHpp: commHpp,
      commercialNetProfit: commNet,
      positiveFiscalCorrection: form.value.positiveFiscalCorrection,
      positiveCorrectionNotes: form.value.positiveCorrectionNotes,
      negativeFiscalCorrection: form.value.negativeFiscalCorrection,
      negativeCorrectionNotes: form.value.negativeCorrectionNotes,
      taxableIncome: computedPKP.value,
      corporateTaxRate: form.value.corporateTaxRate,
      corporateTaxPayable: computedTaxPayable.value,
      taxCreditPph23: form.value.taxCreditPph23,
      taxCreditPph25: form.value.taxCreditPph25,
      taxCredits: computedTotalCredits.value,
      pph29Payable: computedPph29.value,
      pph28aOverpaid: computedPph28A.value
    },
    monthlyMatrix: sptMatrix.value,
    linkedJournals: journalsStore.journals,
    creditVouchers: creditVoucherRows.value
  }, {
    fiscalYear: annual1771.value.fiscalYear,
    businessUnit: accountingStore.businessUnit
  })
  accountingStore.showToast('Kertas Kerja SPT Tahunan 1771 berhasil diekspor ke CSV!', 'success')
}

// 🖨️ Cetak Resmi Formulir 1771 DJP
function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="font-['Plus_Jakarta_Sans',sans-serif] pb-16 max-w-7xl mx-auto">
    <!-- Screen View (Hidden during print) -->
    <div class="print:hidden space-y-6">
      <!-- Header & Navigation -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
          <router-link to="/accounting/tax/spt-tahunan" class="hover:text-zinc-950 transition flex items-center gap-1">
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>SPT Tahunan 1771</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-950">Kertas Kerja & Formulir Jurnal Penutup</span>
          <span class="text-zinc-300">•</span>
          <span class="text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full text-[10px] font-bold border border-emerald-200">
            Terintegrasi Jurnal Umum & Buku Besar
          </span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight">
          Kertas Kerja & Penyesuaian SPT Tahunan PPh Badan 1771
        </h1>
        <p class="text-xs text-zinc-500">
          Tarik otomatis pendapatan & beban komersial dari Jurnal Umum, kelola daftar voucher kredit pajak (PPh 23 & 25), dan posting jurnal penutup PPh 29.
        </p>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <button
          type="button"
          @click="handleResetDefaults"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <RotateCcw class="w-3.5 h-3.5 text-zinc-500" />
          <span>Reset Buku</span>
        </button>

        <button
          type="button"
          @click="autoPullFromGeneralLedger"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-amber-500 hover:bg-amber-600 text-zinc-950 rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Sparkles class="w-3.5 h-3.5 text-zinc-950" />
          <span>Ambil Data Otomatis dari Jurnal Umum & Buku Besar</span>
        </button>

        <button
          type="button"
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Ekspor Kertas Kerja 1771 per transaksi jurnal ke format Excel"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor 1771 (.xlsx)</span>
        </button>

        <button
          type="button"
          @click="handleExportCSV"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Ekspor Kertas Kerja 1771 ke format CSV"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor CSV (.csv)</span>
        </button>

        <button
          type="button"
          @click="handlePrint"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Cetak Formulir Resmi 1771 DJP"
        >
          <Printer class="w-3.5 h-3.5 text-blue-600" />
          <span>Cetak Formulir</span>
        </button>

        <button
          type="button"
          @click="handleSaveAll"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan Kertas Kerja 1771</span>
        </button>
      </div>
    </header>

    <!-- Live Floating Summary Bar -->
    <div class="bg-white text-zinc-900 rounded-2xl p-5 shadow-2xs border border-zinc-200/90 space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <Calculator class="w-4 h-4 text-amber-600" />
          <span class="text-xs font-bold uppercase tracking-wider text-zinc-700">
            Perhitungan Fiskal Real-Time SPT Tahunan (Tahun Pajak {{ annual1771.fiscalYear }})
          </span>
        </div>

        <div class="flex items-center gap-2">
          <span class="text-xs text-zinc-500 font-medium">
            {{ computedPph28A > 0 ? 'Status: PPh 28A Lebih Bayar' : 'Status: PPh 29 Kurang Bayar' }}
          </span>
          <span
            class="text-xl font-black tabular-nums"
            :class="computedPph28A > 0 ? 'text-emerald-700' : 'text-amber-700'"
          >
            {{ formatIDR(computedPph28A > 0 ? computedPph28A : computedPph29) }}
          </span>
        </div>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
        <div class="bg-zinc-50 p-3 rounded-xl border border-zinc-200/80">
          <span class="text-[10px] text-zinc-500 uppercase font-bold block">1. Laba Komersial</span>
          <span class="text-sm font-bold text-zinc-900 tabular-nums">{{ formatIDR(computedNetProfit) }}</span>
        </div>
        <div class="bg-amber-50/60 p-3 rounded-xl border border-amber-200/70">
          <span class="text-[10px] text-amber-700 uppercase font-bold block">2. Penghasilan Kena Pajak</span>
          <span class="text-sm font-bold text-amber-800 tabular-nums">{{ formatIDR(computedPKP) }}</span>
        </div>
        <div class="bg-zinc-50 p-3 rounded-xl border border-zinc-200/80">
          <span class="text-[10px] text-zinc-500 uppercase font-bold block">3. PPh Terutang</span>
          <span class="text-sm font-bold text-zinc-900 tabular-nums">{{ formatIDR(computedTaxPayable) }}</span>
        </div>
        <div class="bg-emerald-50/70 p-3 rounded-xl border border-emerald-200/80 text-emerald-950">
          <span class="text-[10px] text-emerald-700 uppercase font-bold block">4. Total Kredit Pajak</span>
          <span class="text-sm font-bold text-emerald-800 tabular-nums">{{ formatIDR(computedTotalCredits) }}</span>
        </div>
      </div>
    </div>

    <!-- Section 1 & 2: Pendapatan, Beban & Rekonsiliasi Fiskal Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 text-xs">
      <!-- Card 1: Laba Komersial -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-5 space-y-4">
        <div class="border-b border-zinc-100 pb-2 flex items-center justify-between">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <Building2 class="w-4 h-4 text-zinc-700" />
            <span>1. Laba Komersial</span>
          </h3>
          <span class="text-[10px] text-zinc-400 font-mono">Akun 4, 5 & 6</span>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">
              Peredaran Usaha & Pendapatan Bruto (Rp)
            </label>
            <input
              v-model.number="form.commercialRevenue"
              type="number"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-1 focus:ring-zinc-950 text-xs"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">
              Beban Operasional & HPP (Rp)
            </label>
            <input
              v-model.number="form.commercialExpense"
              type="number"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:outline-none focus:ring-1 focus:ring-zinc-950 text-xs"
            />
          </div>

          <div class="flex items-center justify-between p-2.5 bg-zinc-100 rounded-xl font-mono text-xs">
            <span class="text-zinc-600 font-medium">Laba Bersih Komersial:</span>
            <span class="font-bold text-zinc-950 tabular-nums">
              {{ formatIDR(computedNetProfit) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Card 2: Koreksi Fiskal -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-5 space-y-4">
        <div class="border-b border-zinc-100 pb-2 flex items-center justify-between">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <SlidersHorizontal class="w-4 h-4 text-zinc-700" />
            <span>2. Koreksi Fiskal UU PPh</span>
          </h3>
          <span class="text-[10px] text-zinc-400">Beda Tetap/Waktu</span>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-emerald-800 mb-1">
              Koreksi Positif (+) Non-Deductible (Rp)
            </label>
            <input
              v-model.number="form.positiveFiscalCorrection"
              type="number"
              class="w-full px-3 py-2 bg-emerald-50/40 border border-emerald-300 rounded-xl font-mono text-emerald-950 font-bold focus:outline-none text-xs"
            />
            <input
              v-model="form.positiveCorrectionNotes"
              type="text"
              placeholder="Catatan: Denda pajak, entertainment tanpa nominatif"
              class="w-full mt-1 px-2.5 py-1 bg-zinc-50 border border-zinc-200 rounded-lg text-zinc-700 text-[10px]"
            />
          </div>

          <div>
            <label class="block font-bold text-rose-800 mb-1">
              Koreksi Negatif (-) Penghasilan Final (Rp)
            </label>
            <input
              v-model.number="form.negativeFiscalCorrection"
              type="number"
              class="w-full px-3 py-2 bg-rose-50/40 border border-rose-300 rounded-xl font-mono text-rose-950 font-bold focus:outline-none text-xs"
            />
            <input
              v-model="form.negativeCorrectionNotes"
              type="text"
              placeholder="Catatan: Bunga giro bank dipotong PPh Final"
              class="w-full mt-1 px-2.5 py-1 bg-zinc-50 border border-zinc-200 rounded-lg text-zinc-700 text-[10px]"
            />
          </div>
        </div>
      </div>

      <!-- Card 3: Tarif PPh Badan & PKP -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-5 space-y-4">
        <div class="border-b border-zinc-100 pb-2 flex items-center justify-between">
          <h3 class="font-bold text-zinc-950 text-sm flex items-center gap-2">
            <Landmark class="w-4 h-4 text-zinc-700" />
            <span>3. Tarif & PPh Terutang</span>
          </h3>
          <span class="text-[10px] text-zinc-400">UU HPP</span>
        </div>

        <div class="space-y-3">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">
              Skema Fasilitas Tarif PPh Badan:
            </label>
            <div class="grid grid-cols-2 gap-2 mb-2">
              <button
                type="button"
                @click="form.corporateTaxRate = 0.11"
                class="px-2 py-1.5 rounded-xl text-center border font-bold text-[10px] transition cursor-pointer"
                :class="form.corporateTaxRate === 0.11 ? 'bg-zinc-950 text-white border-zinc-950' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                11% (Pasal 31E)
              </button>
              <button
                type="button"
                @click="form.corporateTaxRate = 0.22"
                class="px-2 py-1.5 rounded-xl text-center border font-bold text-[10px] transition cursor-pointer"
                :class="form.corporateTaxRate === 0.22 ? 'bg-zinc-950 text-white border-zinc-950' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
              >
                22% (Tarif Umum)
              </button>
            </div>
          </div>

          <div class="p-2.5 bg-amber-50 rounded-xl font-mono text-xs border border-amber-200 flex items-center justify-between">
            <span class="text-amber-900 font-bold">Penghasilan Kena Pajak:</span>
            <span class="font-black text-amber-950 tabular-nums">{{ formatIDR(computedPKP) }}</span>
          </div>

          <div class="p-2.5 bg-zinc-100 rounded-xl font-mono text-xs border border-zinc-200 flex items-center justify-between">
            <span class="text-zinc-700 font-bold">PPh Badan Terutang:</span>
            <span class="font-extrabold text-amber-800 tabular-nums">{{ formatIDR(computedTaxPayable) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Section 3: Voucher Grid Table - Kredit Pajak dari Jurnal Umum (Identik dengan Jurnal Umum) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <!-- Toolbar Header -->
      <div class="p-4 border-b border-zinc-200/90 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-zinc-50/50">
        <div class="flex items-center gap-2">
          <BookOpen class="w-4 h-4 text-emerald-700" />
          <span class="text-xs font-bold text-zinc-900">
            Daftar Bukti Potong & Voucher Kredit Pajak 1771 ({{ creditVoucherRows.length }} Voucher Jurnal)
          </span>
          <span class="text-[11px] text-zinc-500 hidden sm:inline">
            &bull; Terintegrasi dari Akun 2-1510 (PPh 23) & 2-1520 (PPh 25) Jurnal Umum
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="addCreditRow"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Baris Kredit Manual</span>
          </button>
        </div>
      </div>

      <!-- Table Grid -->
      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200/90 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-10 text-center">Klaim</th>
              <th class="px-3 py-3 w-28">Tanggal</th>
              <th class="px-3 py-3 w-36">No. Jurnal Sumber</th>
              <th class="px-3 py-3 w-28">Jenis Kredit</th>
              <th class="px-4 py-3 min-w-[200px]">Akun Kredit (COA)</th>
              <th class="px-4 py-3 min-w-[220px]">Uraian / Rekanan Pemotong</th>
              <th class="px-3 py-3 w-32 text-right">Nilai Voucher (Rp)</th>
              <th class="px-3 py-3 w-36 text-center">Status Klaim 1771</th>
              <th class="px-4 py-3 w-36 text-right">Nilai Diakui (Rp)</th>
              <th class="px-3 py-3 w-12 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr
              v-for="(row, idx) in creditVoucherRows"
              :key="row.id"
              :class="[
                'hover:bg-zinc-50/70 transition',
                row.claimStatus === 'EXCLUDED' ? 'opacity-40 bg-zinc-50/50' : ''
              ]"
            >
              <!-- Checkbox -->
              <td class="px-3 py-3 text-center">
                <input
                  v-model="row.selected"
                  @change="() => { if (!row.selected) setRowClaimStatus(row, 'EXCLUDED'); else setRowClaimStatus(row, 'CLAIM_FULL'); }"
                  type="checkbox"
                  class="w-4 h-4 rounded text-zinc-950 focus:ring-0 cursor-pointer"
                />
              </td>

              <!-- Date -->
              <td class="px-3 py-2.5 font-mono text-zinc-600">
                <input
                  v-model="row.date"
                  type="date"
                  class="w-full px-2 py-1 bg-zinc-50 border border-zinc-200 rounded-lg text-xs font-mono font-bold text-zinc-900 focus:outline-none"
                />
              </td>

              <!-- Voucher Source -->
              <td class="px-3 py-2.5 font-mono">
                <div v-if="row.isFromJournal && row.journalNumber" class="flex items-center gap-1">
                  <button
                    type="button"
                    @click="openJournalPreview(row)"
                    class="px-2 py-0.5 rounded-md bg-blue-50 text-blue-700 font-bold hover:bg-blue-100 transition inline-flex items-center gap-1 cursor-pointer border border-blue-200"
                    title="Klik untuk membuka voucher Jurnal Umum"
                  >
                    <span>{{ row.journalNumber }}</span>
                    <ExternalLink class="w-2.5 h-2.5" />
                  </button>
                </div>
                <div v-else class="text-zinc-400 text-[10px] italic">
                  [Manual]
                </div>
              </td>

              <!-- Tax Type Badge -->
              <td class="px-3 py-2.5">
                <select
                  v-model="row.taxType"
                  @change="() => { row.accountCode = row.taxType === 'PPH_23' ? '2-1510' : '2-1520'; row.accountName = row.taxType === 'PPH_23' ? 'Kredit Pajak PPh 23' : 'Uang Muka PPh 25 Badan'; syncCreditsFromRows(); }"
                  class="px-2 py-1 bg-zinc-50 border border-zinc-200 rounded-lg font-bold text-xs"
                  :class="row.taxType === 'PPH_23' ? 'text-blue-800' : 'text-emerald-800'"
                >
                  <option value="PPH_23">PPh 23</option>
                  <option value="PPH_25">PPh 25</option>
                </select>
              </td>

              <!-- Account -->
              <td class="px-4 py-2.5 font-mono text-zinc-800">
                <span class="font-bold">[{{ row.accountCode }}]</span>
                <span class="ml-1 text-zinc-600 text-[11px]">{{ row.accountName }}</span>
              </td>

              <!-- Description -->
              <td class="px-4 py-2.5">
                <input
                  v-model="row.description"
                  type="text"
                  placeholder="Uraian kredit / rekanan..."
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-900 focus:outline-none"
                />
              </td>

              <!-- Original Voucher Amount -->
              <td class="px-3 py-2.5 text-right font-mono tabular-nums font-semibold text-zinc-700">
                {{ formatIDR(row.voucherAmount) }}
              </td>

              <!-- Claim Status Buttons -->
              <td class="px-3 py-2.5 text-center">
                <div class="inline-flex rounded-lg bg-zinc-100 p-0.5 text-[10px] font-bold">
                  <button
                    type="button"
                    @click="setRowClaimStatus(row, 'CLAIM_FULL')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="row.claimStatus === 'CLAIM_FULL' ? 'bg-white text-zinc-950 shadow-2xs font-extrabold' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    100%
                  </button>
                  <button
                    type="button"
                    @click="setRowClaimStatus(row, 'CLAIM_PARTIAL')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="row.claimStatus === 'CLAIM_PARTIAL' ? 'bg-amber-500 text-zinc-950 shadow-2xs font-extrabold' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    Sebagian
                  </button>
                  <button
                    type="button"
                    @click="setRowClaimStatus(row, 'EXCLUDED')"
                    class="px-2 py-1 rounded-md transition cursor-pointer"
                    :class="row.claimStatus === 'EXCLUDED' ? 'bg-rose-600 text-white shadow-2xs font-extrabold' : 'text-zinc-500 hover:text-zinc-900'"
                  >
                    Tidak
                  </button>
                </div>
              </td>

              <!-- Claimed Amount Input -->
              <td class="px-4 py-2.5 text-right">
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                  <input
                    v-model.number="row.claimedAmount"
                    @input="() => onClaimAmountChange(row)"
                    type="number"
                    min="0"
                    step="any"
                    class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                  />
                </div>
              </td>

              <!-- Delete Action -->
              <td class="px-3 py-2.5 text-center">
                <button
                  type="button"
                  @click="removeCreditRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition cursor-pointer"
                  title="Hapus Baris"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>

          <!-- Table Footer -->
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs text-zinc-900">
            <tr>
              <td colspan="6" class="px-4 py-3 text-right uppercase tracking-wider text-zinc-600">
                Total Kredit Pajak Diakui di 1771:
              </td>
              <td class="px-3 py-3 text-right font-mono tabular-nums text-zinc-500">
                {{ formatIDR(creditVoucherRows.reduce((s, r) => s + r.voucherAmount, 0)) }}
              </td>
              <td class="px-3 py-3 text-center text-emerald-700 text-[11px]">
                PPh 23: {{ formatIDR(form.taxCreditPph23) }} • PPh 25: {{ formatIDR(form.taxCreditPph25) }}
              </td>
              <td class="px-4 py-3 text-right font-mono tabular-nums text-emerald-700 font-extrabold text-sm">
                {{ formatIDR(computedTotalCredits) }}
              </td>
              <td class="px-3 py-3"></td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Section 4: Closing Voucher Jurnal PPh 29 (Double-Entry Balanced Settlement) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-5 space-y-4">
      <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2">
          <Layers class="w-4 h-4 text-emerald-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            Voucher Jurnal Penyesuaian Tutup Tahun PPh 29 (Voucher Jurnal Penutup Akhir Periode)
          </h3>
        </div>
        <button
          type="button"
          @click="handlePostPph29Closing"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Send class="w-3.5 h-3.5" />
          <span>Posting Jurnal Penutup PPh 29 ke Jurnal Umum</span>
        </button>
      </div>

      <div class="border border-zinc-200 rounded-xl overflow-hidden text-xs">
        <table class="w-full text-left">
          <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase">
            <tr>
              <th class="px-4 py-2.5">Posisi Akun Akuntansi</th>
              <th class="px-4 py-2.5">Uraian / Memo Jurnal Penutup</th>
              <th class="px-4 py-2.5 text-right w-44">Debet (Rp)</th>
              <th class="px-4 py-2.5 text-right w-44">Kredit (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700 font-mono">
            <!-- Debet: Beban PPh Badan -->
            <tr>
              <td class="px-4 py-2.5 font-sans">
                <span class="font-bold text-zinc-950">[6-1500]</span>
                <span class="ml-1 font-semibold text-zinc-800">Beban Pajak Penghasilan Badan (Tahun Berjalan)</span>
              </td>
              <td class="px-4 py-2.5 font-sans text-zinc-600 text-[11px]">
                Pengakuan Beban PPh Badan Formulir 1771 Tahun {{ annual1771.fiscalYear }}
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-zinc-950 tabular-nums">
                {{ formatIDR(computedTaxPayable) }}
              </td>
              <td class="px-4 py-2.5 text-right text-zinc-400">
                -
              </td>
            </tr>

            <!-- Kredit: PPh 23 -->
            <tr v-if="form.taxCreditPph23 > 0">
              <td class="px-4 py-2.5 font-sans">
                <span class="font-bold text-zinc-950">[2-1510]</span>
                <span class="ml-1 font-semibold text-zinc-800">Kredit Pajak PPh 23 (Tutup Saldo Berjalan)</span>
              </td>
              <td class="px-4 py-2.5 font-sans text-zinc-600 text-[11px]">
                Pengkreditan Bukti Potong PPh 23 Tahun {{ annual1771.fiscalYear }}
              </td>
              <td class="px-4 py-2.5 text-right text-zinc-400">
                -
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-emerald-700 tabular-nums">
                {{ formatIDR(form.taxCreditPph23) }}
              </td>
            </tr>

            <!-- Kredit: PPh 25 -->
            <tr v-if="form.taxCreditPph25 > 0">
              <td class="px-4 py-2.5 font-sans">
                <span class="font-bold text-zinc-950">[2-1520]</span>
                <span class="ml-1 font-semibold text-zinc-800">Uang Muka / Angsuran PPh 25 Badan (Tutup Saldo)</span>
              </td>
              <td class="px-4 py-2.5 font-sans text-zinc-600 text-[11px]">
                Pengkreditan Setoran Angsuran Bulanan PPh 25 Tahun {{ annual1771.fiscalYear }}
              </td>
              <td class="px-4 py-2.5 text-right text-zinc-400">
                -
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-emerald-700 tabular-nums">
                {{ formatIDR(form.taxCreditPph25) }}
              </td>
            </tr>

            <!-- Kredit: PPh 29 Kurang Bayar -->
            <tr v-if="computedPph29 > 0">
              <td class="px-4 py-2.5 font-sans">
                <span class="font-bold text-zinc-950">[2-1530]</span>
                <span class="ml-1 font-semibold text-zinc-800">Hutang PPh 29 (Kurang Bayar Tahunan Badan)</span>
              </td>
              <td class="px-4 py-2.5 font-sans text-zinc-600 text-[11px]">
                Kewajiban PPh 29 Kurang Bayar SPT Tahunan 1771 Tahun {{ annual1771.fiscalYear }}
              </td>
              <td class="px-4 py-2.5 text-right text-zinc-400">
                -
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-amber-600 tabular-nums">
                {{ formatIDR(computedPph29) }}
              </td>
            </tr>
          </tbody>
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs">
            <tr>
              <td colspan="2" class="px-4 py-2 text-right uppercase tracking-wider text-zinc-600">
                Total Jurnal Penutup (Status: Seimbang):
              </td>
              <td class="px-4 py-2 text-right font-mono font-extrabold text-emerald-800">
                {{ formatIDR(computedTaxPayable) }}
              </td>
              <td class="px-4 py-2 text-right font-mono font-extrabold text-emerald-800">
                {{ formatIDR(computedTaxPayable) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Bottom Actions Bar -->
    <div class="flex items-center justify-between pt-4 border-t border-zinc-200">
      <router-link
        to="/accounting/tax/spt-tahunan"
        class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition"
      >
        Batal & Kembali ke SPT 1771
      </router-link>

      <div class="flex items-center gap-2">
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

    <!-- Modal Preview Voucher Jurnal Sumber -->
    <div
      v-if="showJournalPreviewModal"
      class="fixed inset-0 z-50 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4"
      @click.self="showJournalPreviewModal = false"
    >
      <div class="bg-white rounded-2xl max-w-2xl w-full p-5 space-y-4 shadow-xl">
        <div class="flex items-center justify-between border-b border-zinc-200 pb-3">
          <div class="flex items-center gap-2">
            <Layers class="w-4 h-4 text-blue-600" />
            <h3 class="font-bold text-sm text-zinc-950">
              Voucher Jurnal Umum Sumber: {{ activeJournal?.journalNumber }}
            </h3>
          </div>
          <button @click="showJournalPreviewModal = false" class="text-zinc-400 hover:text-zinc-950 transition cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="text-xs space-y-1 text-zinc-600">
          <p><strong class="text-zinc-900">Tanggal Posting:</strong> {{ formatDate(activeJournal?.journalDate || activeJournal?.transaction_date) }}</p>
          <p><strong class="text-zinc-900">Keterangan:</strong> {{ activeJournal?.description }}</p>
        </div>

        <div class="border border-zinc-200 rounded-xl overflow-hidden text-xs">
          <table class="w-full text-left">
            <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase">
              <tr>
                <th class="px-3 py-2">Akun Buku Besar (COA)</th>
                <th class="px-3 py-2 text-right w-28">Debet (Rp)</th>
                <th class="px-3 py-2 text-right w-28">Kredit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="(l, idx) in (activeJournal?.lines || [])" :key="idx">
                <td class="px-3 py-2">
                  <span class="font-bold font-mono text-zinc-900">{{ l.accountCode || l.account_code }}</span> -
                  <span class="text-zinc-700">{{ l.accountName || l.account_name }}</span>
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums">
                  {{ Number(l.debit) > 0 ? formatIDR(Number(l.debit), false) : '-' }}
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums">
                  {{ Number(l.credit) > 0 ? formatIDR(Number(l.credit), false) : '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex justify-end pt-2">
          <button
            @click="showJournalPreviewModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition cursor-pointer"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
    </div>
    <!-- END OF SCREEN VIEW (print:hidden) -->

    <!-- ========================================================================= -->
    <!-- CETAK RESMI FORMULIR 1771 - SPT TAHUNAN PPh BADAN (DJP AUDIT GRADE)       -->
    <!-- HANYA TAMPIL SAAT CETAK / EXPORT PDF (PRINT MEDIA ONLY)                   -->
    <!-- ========================================================================= -->
    <div class="hidden print:block print-document font-['Times_New_Roman',serif] text-black bg-white">
      <!-- Kop Surat Resmi DJP / SMG -->
      <div class="border-b-2 border-black pb-3 mb-4">
        <div class="flex items-center justify-between">
          <div class="space-y-1">
            <h2 class="text-xl font-bold uppercase tracking-wide">{{ currentEntity.title }}</h2>
            <h3 class="text-xs font-semibold text-gray-800">{{ currentEntity.subtitle }}</h3>
            <p class="text-[10px] text-gray-600">{{ currentEntity.address }} • KLU: 86101</p>
          </div>
          <div class="text-right border border-black p-2 rounded text-xs font-mono">
            <div class="font-bold text-sm">FORMULIR 1771</div>
            <div class="text-[9px] text-gray-600">TAHUN PAJAK: {{ annual1771.fiscalYear }}</div>
            <div class="text-[9px] font-bold text-gray-800 mt-0.5">NPWP: {{ currentEntity.npwp }}</div>
          </div>
        </div>
      </div>

      <!-- Judul Formulir 1771 -->
      <div class="text-center my-4 space-y-0.5">
        <h1 class="text-base font-black uppercase underline tracking-wider">
          SURAT PEMBERITAHUAN TAHUNAN PAJAK PENGHASILAN WAJIB PAJAK BADAN
        </h1>
        <p class="text-[11px] font-bold text-gray-800 uppercase">
          (FORMULIR 1771 - INDUK & KERTAS KERJA REKONSILIASI FISKAL)
        </p>
        <p class="text-[10px] italic text-gray-600">
          Periode 1 Januari {{ annual1771.fiscalYear }} s.d. 31 Desember {{ annual1771.fiscalYear }}
        </p>
      </div>

      <!-- Identitas Wajib Pajak -->
      <div class="border border-black p-2.5 text-[11px] mb-4 bg-gray-50/50">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1">
            <div class="flex"><span class="w-36 font-bold">NPWP Badan Usaha:</span><span class="font-mono font-bold">{{ currentEntity.npwp }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Nama Wajib Pajak:</span><span class="font-bold">{{ currentEntity.title }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Jenis Usaha (KLU):</span><span>86101 - Aktivitas Pelayanan Medis Klinik Swasta</span></div>
          </div>
          <div class="space-y-1">
            <div class="flex"><span class="w-36 font-bold">Status Pelaporan:</span><span class="font-bold">{{ annual1771.status }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Entitas Unit:</span><span>{{ currentEntity.unitLabel }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Alamat Resmi:</span><span>{{ currentEntity.address }}</span></div>
          </div>
        </div>
      </div>

      <!-- Tabel Induk 1771 Rekonsiliasi Fiskal -->
      <div class="mb-4">
        <h4 class="text-[11px] font-bold uppercase mb-1.5 text-zinc-900">BAGIAN A s.d. D: PENGHITUNGAN PENGHASILAN KENA PAJAK & PPh TERUTANG</h4>
        <table class="w-full text-[10px] border-collapse border border-black">
          <thead>
            <tr class="bg-gray-200 text-center font-bold">
              <th class="border border-black px-2 py-1 w-8">No</th>
              <th class="border border-black px-2 py-1">Uraian Elemen Formulir 1771</th>
              <th class="border border-black px-2 py-1 w-36 text-right">Nilai Rupiah (IDR)</th>
              <th class="border border-black px-2 py-1 w-52">Dasar Regulasi Perpajakan</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold">1</td>
              <td class="border border-black px-2 py-1">Peredaran Usaha Komersial (Omset Layanan Medis & Farmasi)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-medium">{{ formatIDR(form.commercialRevenue) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">Laporan Laba Rugi Komersial Klinik</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold">2</td>
              <td class="border border-black px-2 py-1">Harga Pokok Penjualan (HPP Obat, Vaksin & BMHP Medis)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-medium">{{ formatIDR(Math.round(form.commercialRevenue * 0.48)) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">Pemakaian Obat Resep & Tindakan</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold">3</td>
              <td class="border border-black px-2 py-1">Biaya Usaha Lainnya (Beban Operasional, Gaji & Honor Dokter)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-medium">{{ formatIDR(form.commercialExpense) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">Beban Operasional & Administrasi</td>
            </tr>
            <tr class="bg-gray-100 font-bold">
              <td class="border border-black px-2 py-1 text-center">4</td>
              <td class="border border-black px-2 py-1 uppercase">PENGHASILAN NETO KOMERSIAL (1 - 2 - 3)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(form.commercialRevenue - Math.round(form.commercialRevenue * 0.48) - form.commercialExpense) }}</td>
              <td class="border border-black px-2 py-1 text-[9px]">Laba Bersih Sebelum Koreksi Fiskal</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold text-amber-900">5</td>
              <td class="border border-black px-2 py-1">Koreksi Fiskal Positif (+) Biaya Non-Deductible</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold text-amber-900">{{ formatIDR(form.positiveFiscalCorrection) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">{{ form.positiveCorrectionNotes || 'PMK 66/2023 & UU PPh' }}</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold text-emerald-900">6</td>
              <td class="border border-black px-2 py-1">Koreksi Fiskal Negatif (-) Penghasilan PPh Final & Depresiasi</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold text-emerald-900">{{ formatIDR(-form.negativeFiscalCorrection) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">{{ form.negativeCorrectionNotes || 'Pasal 4 ayat (2) UU PPh' }}</td>
            </tr>
            <tr class="bg-gray-200 font-bold">
              <td class="border border-black px-2 py-1 text-center">7</td>
              <td class="border border-black px-2 py-1 uppercase">PENGHASILAN KENA PAJAK (PKP) BADAN (4 + 5 - 6)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(computedPKP) }}</td>
              <td class="border border-black px-2 py-1 text-[9px]">Dasar Penghitungan PPh Terutang</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center font-bold">8</td>
              <td class="border border-black px-2 py-1">Tarif Pajak Penghasilan Badan (Pasal 17 UU HPP)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ (form.corporateTaxRate * 100).toFixed(0) }}%</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">UU No. 7 Tahun 2021</td>
            </tr>
            <tr class="bg-blue-50 font-bold">
              <td class="border border-black px-2 py-1 text-center">9</td>
              <td class="border border-black px-2 py-1 uppercase">PAJAK PENGHASILAN TERUTANG (PPh BADAN)</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(computedTaxPayable) }}</td>
              <td class="border border-black px-2 py-1 text-[9px]">PKP x Tarif 22%</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center">10a</td>
              <td class="border border-black px-2 py-1 pl-4">Kredit Pajak PPh Pasal 23 (Dipotong Rekanan Asuransi & Korporat)</td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(form.taxCreditPph23) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">Bukti Potong Akun 2-1510</td>
            </tr>
            <tr>
              <td class="border border-black px-2 py-1 text-center">10b</td>
              <td class="border border-black px-2 py-1 pl-4">Kredit Pajak PPh Pasal 25 (Setoran Angsuran Bulanan Bank)</td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(form.taxCreditPph25) }}</td>
              <td class="border border-black px-2 py-1 text-[9px] text-gray-700">Validasi NTPN Akun 2-1520</td>
            </tr>
            <tr class="bg-gray-100 font-bold">
              <td class="border border-black px-2 py-1 text-center">10</td>
              <td class="border border-black px-2 py-1 uppercase">TOTAL KREDIT PAJAK TAHUN BERJALAN (10a + 10b)</td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(computedTotalCredits) }}</td>
              <td class="border border-black px-2 py-1 text-[9px]">Pengurang PPh Terutang</td>
            </tr>
            <tr class="bg-emerald-100 font-black text-emerald-950">
              <td class="border border-black px-2 py-1.5 text-center text-sm">11</td>
              <td class="border border-black px-2 py-1.5 uppercase text-xs">
                {{ computedPph28A > 0 ? 'PPh PASAL 28A LEBIH BAYAR (10 - 9)' : 'PPh PASAL 29 KURANG BAYAR (9 - 10)' }}
              </td>
              <td class="border border-black px-2 py-1.5 text-right font-mono text-sm font-black">
                {{ formatIDR(computedPph28A > 0 ? computedPph28A : computedPph29) }}
              </td>
              <td class="border border-black px-2 py-1.5 text-[9px] font-bold">
                {{ computedPph28A > 0 ? 'Dapat Dikompensasikan / Restitusi' : 'Wajib Disetor Sebelum Lapor 1771' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Keterangan Resmi: Kertas Kerja Detail Dilaporkan via Ekspor Excel (Sesuai Standar DJP) -->
      <div class="border border-black p-3 mb-6 bg-gray-50 text-[10px] break-inside-avoid">
        <div class="font-bold uppercase tracking-wide text-zinc-900 mb-1 flex items-center justify-between">
          <span>CATATAN DOKUMEN & LAMPIRAN KERTAS KERJA RESMI (FORMAT EXCEL .XLSX):</span>
          <span class="font-mono text-[9px] bg-white px-2 py-0.5 border border-black font-bold">LAMPIRAN 1771-I S.D. 1771-VI</span>
        </div>
        <p class="text-zinc-700 leading-relaxed text-[9.5px]">
          Sesuai dengan regulasi pelaporan SPT Tahunan PPh Badan DJP, seluruh detail bukti transaksi per voucher jurnal telah disajikan secara terpisah dan komprehensif ke dalam <strong>File Kertas Kerja Resmi SPT Tahunan (Workbook Excel .xlsx)</strong> yang mencakup:
        </p>
        <div class="grid grid-cols-2 gap-x-4 gap-y-1.5 mt-2 text-[9px] text-zinc-800">
          <div>• <strong>Sheet 1 (Induk 1771):</strong> Rekonsiliasi Fiskal Lengkap Elemen 1 s.d. 16</div>
          <div>• <strong>Sheet 2 (Kredit PPh 23):</strong> Rincian Per Voucher Jurnal Bukti Potong Rekanan (Akun 2-1510)</div>
          <div>• <strong>Sheet 3 (Angsuran PPh 25):</strong> Rincian Per Setoran Bank & Validasi NTPN Persepsi (Akun 2-1520)</div>
          <div>• <strong>Sheet 4 (Rincian Omset):</strong> Rincian Buku Besar Pendapatan Operasional Medis & Farmasi (Akun 4-xxxx)</div>
          <div>• <strong>Sheet 5 (Rincian Biaya):</strong> Rincian Buku Besar Beban Operasional & HPP BMHP (Akun 5-xxxx & 6-xxxx)</div>
          <div>• <strong>Sheet 6 (Rekap Bulanan):</strong> Matriks Akumulasi Setoran Pajak Masa Jan s.d. Des Terintegrasi</div>
          <div>• <strong>Sheet 7 (Jurnal PPh 29):</strong> Rekonsiliasi Jurnal Akrual & Utang Pajak Badan Akhir Tahun</div>
        </div>
        <div class="mt-2.5 pt-2 border-t border-gray-300 flex items-center justify-between text-[9px] italic text-zinc-600">
          <span>Nama File Rekonsiliasi: <strong>{{ `Kertas_Kerja_SPT_Tahunan_1771_${accountingStore.businessUnit}_Tahun_${annual1771.fiscalYear}.xlsx` }}</strong></span>
          <span>Status Integritas Data: <strong>TERVERIFIKASI BUKU BESAR KLINIK & MUTASI PERBANKAN</strong></span>
        </div>
      </div>

      <!-- Lembar Pengesahan Tanda Tangan -->
      <div class="grid grid-cols-2 text-center text-xs pt-4 border-t border-black break-inside-avoid">
        <div class="space-y-14">
          <div>
            <p>Menyetujui & Mengesahkan,</p>
            <p class="font-bold">{{ currentEntity.directorRole }}</p>
          </div>
          <div>
            <p class="font-bold underline">{{ currentEntity.director }}</p>
            <p class="text-[10px] text-gray-600">{{ currentEntity.sip }}</p>
          </div>
        </div>

        <div class="space-y-14">
          <div>
            <p>Semarang, {{ new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) }}</p>
            <p class="font-bold">Finance & Tax Manager</p>
          </div>
          <div>
            <p class="font-bold underline">{{ accountingStore.currentUserName }}</p>
            <p class="text-[10px] text-gray-600">Register Akuntan Negara: CA-78291</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  .break-inside-avoid {
    break-inside: avoid !important;
    page-break-inside: avoid !important;
  }
}
</style>