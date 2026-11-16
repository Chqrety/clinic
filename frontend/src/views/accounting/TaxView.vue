<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, type SptMasaType } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
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
  BookOpen,
  Plus,
  X,
  Printer,
  Sparkles,
  Link2,
  Landmark,
  Database
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()

onMounted(async () => {
  await Promise.all([
    taxStore.fetchTaxRecords(),
    ledgerStore.fetchAccounts()
  ])
})

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
const sptMasaList = computed(() => taxStore.filteredSptMasa)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

// Subview under SPT Masa: 'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU'
const sptMasaSubTab = ref<'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU'>('DOKUMEN_SPT')

// Inline Table Input & Modals
const showInlineTaxForm = ref(false)
const showCreateSptMasaModal = ref(false)
const showDetailModal = ref(false)
const showPaymentBankModal = ref(false)
const selectedSptDetail = ref<any>(null)
const selectedBankForPayment = ref<string>('1102.001')

// Form for new SPT Masa (Inline Table Entry & Strictly Bank Channel Only)
const newSptMasaForm = ref({
  masaMonth: 11,
  masaYear: 2026,
  unit: 'KIDINEST' as 'KIDINEST' | 'SSC',
  taxType: 'PPH_23' as SptMasaType,
  bankAccountCode: '1102.001',
  grossAmount: 150000000,
  dppAmount: 150000000,
  taxRate: '2%',
  taxAmount: 3000000,
  ntpnNumber: '',
  notes: '',
  status: 'REPORTED' as 'DRAFT' | 'PAID' | 'REPORTED'
})

// Auto update suggested tax amount in modal
function updateSuggestedTaxAmount() {
  const dpp = Number(newSptMasaForm.value.dppAmount) || 0
  if (newSptMasaForm.value.taxType === 'PPH_23') {
    newSptMasaForm.value.taxRate = '2%'
    newSptMasaForm.value.taxAmount = Math.round(dpp * 0.02)
  } else if (newSptMasaForm.value.taxType === 'PPH_4_2') {
    newSptMasaForm.value.taxRate = '10%'
    newSptMasaForm.value.taxAmount = Math.round(dpp * 0.10)
  } else if (newSptMasaForm.value.taxType === 'PPH_25') {
    newSptMasaForm.value.taxRate = 'Angsuran Tetap'
    newSptMasaForm.value.taxAmount = 13000000
  } else if (newSptMasaForm.value.taxType === 'PPN') {
    newSptMasaForm.value.taxRate = '11%'
    newSptMasaForm.value.taxAmount = Math.round(dpp * 0.11)
  }
}

function handleSaveNewSptMasa() {
  taxStore.createSptMasa({
    masaMonth: Number(newSptMasaForm.value.masaMonth),
    masaYear: Number(newSptMasaForm.value.masaYear),
    taxType: newSptMasaForm.value.taxType,
    bankAccountCode: newSptMasaForm.value.bankAccountCode,
    grossAmount: Number(newSptMasaForm.value.grossAmount),
    dppAmount: Number(newSptMasaForm.value.dppAmount),
    taxRate: newSptMasaForm.value.taxRate,
    taxAmount: Number(newSptMasaForm.value.taxAmount),
    notes: `[Unit ${newSptMasaForm.value.unit}] ${newSptMasaForm.value.notes || ''}`.trim(),
    ntpnNumber: newSptMasaForm.value.ntpnNumber,
    status: newSptMasaForm.value.status
  })
  showInlineTaxForm.value = false
  showCreateSptMasaModal.value = false
  accountingStore.showToast(`SPT Masa ${newSptMasaForm.value.taxType} Unit ${newSptMasaForm.value.unit} berhasil dicatat ke tabel!`, 'success')
}

// Selected Account Metadata
const selectedAccountMeta = computed(() => {
  return taxStore.taxAccountOptions.find(o => o.code === taxStore.selectedTaxAccountCode) || taxStore.taxAccountOptions[0]
})

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/accounting/ledger', query: { account: accountCode } })
  } else {
    router.push('/accounting/ledger')
  }
}

function handleExecutePayment() {
  const code = taxStore.selectedTaxAccountCode === 'ALL' ? '2501' : taxStore.selectedTaxAccountCode
  const name = selectedAccountMeta.value.name
  const amt = summary.value.totalTax || 39249649
  taxStore.generateTaxPaymentJournal(code, name, amt, selectedBankForPayment.value)
  showPaymentBankModal.value = false
}

function viewSptDetail(item: any) {
  selectedSptDetail.value = item
  showDetailModal.value = true
}
</script>

<template>
  <div class="space-y-6">
    <!-- Minimalist Clinical Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-200/90 pb-5">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="text-[11px] font-bold text-zinc-400 tracking-wider uppercase">
            Akuntansi & Keuangan &bull; Semarang Medical Group
          </span>
          <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-700 border border-zinc-200">
            <ShieldCheck class="w-3 h-3 text-emerald-600" />
            Strictly Bank Validated
          </span>
        </div>
        <h1 class="text-2xl font-black text-zinc-950 tracking-tight">
          Manajemen Pajak & SPT Masa
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Rekapitulasi SPT Masa, bukti potong e-Bupot 21 medis, dan rekonsiliasi fiskal kredit pajak SPT Tahunan 1771 berbasis kanal bank operasional.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <button
          v-if="taxStore.taxPeriodType === 'MASA'"
          @click="showInlineTaxForm = !showInlineTaxForm"
          :class="[
            'flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold transition shadow-xs',
            showInlineTaxForm ? 'bg-zinc-800 text-white' : 'bg-zinc-950 hover:bg-zinc-800 text-white'
          ]"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>{{ showInlineTaxForm ? 'Tutup Form Tabel' : '+ Entri SPT Masa (Tabel)' }}</span>
        </button>

        <button
          @click="taxStore.exportEbupotCSV()"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 border border-zinc-200 text-zinc-700 rounded-xl text-xs font-bold shadow-2xs transition"
        >
          <Download class="w-3.5 h-3.5 text-zinc-500" />
          <span>Ekspor e-Bupot CSV</span>
        </button>

        <button
          v-if="taxStore.taxPeriodType === 'MASA'"
          @click="showPaymentBankModal = true"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-xs transition"
        >
          <Send class="w-3.5 h-3.5" />
          <span>Jurnal Setor NTPN (via Bank)</span>
        </button>
      </div>
    </div>

    <!-- Mode Selector Tabs: SPT Masa vs SPT Tahunan 1771 -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-3 rounded-2xl border border-zinc-200/90 shadow-2xs">
      <div class="flex items-center gap-2">
        <div class="flex items-center bg-zinc-100 p-1 rounded-xl border border-zinc-200">
          <button
            @click="taxStore.taxPeriodType = 'MASA'"
            :class="[
              'px-4 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-2',
              taxStore.taxPeriodType === 'MASA'
                ? 'bg-zinc-950 text-white shadow-xs'
                : 'text-zinc-600 hover:text-zinc-900'
            ]"
          >
            <Calendar class="w-3.5 h-3.5" />
            <span>1. Pembuatan SPT Masa (Bulanan via Bank)</span>
          </button>
          <button
            @click="taxStore.taxPeriodType = 'TAHUNAN'"
            :class="[
              'px-4 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-2',
              taxStore.taxPeriodType === 'TAHUNAN'
                ? 'bg-zinc-950 text-white shadow-xs'
                : 'text-zinc-600 hover:text-zinc-900'
            ]"
          >
            <Building class="w-3.5 h-3.5" />
            <span>2. SPT Tahunan 1771 (Auto-Credit Bank Mutasi)</span>
          </button>
        </div>
      </div>

      <!-- Context Selectors -->
      <div v-if="taxStore.taxPeriodType === 'MASA'" class="flex flex-wrap items-center gap-2">
        <div class="flex items-center gap-1.5 text-xs">
          <span class="text-zinc-500 font-bold">Masa:</span>
          <select
            v-model.number="taxStore.selectedMonth"
            class="px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 text-zinc-900 rounded-xl text-xs font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          >
            <option v-for="m in 12" :key="m" :value="m">Masa {{ m }}</option>
          </select>
        </div>

        <div class="flex items-center gap-1.5 text-xs">
          <span class="text-zinc-500 font-bold">Tahun:</span>
          <select
            v-model.number="taxStore.selectedYear"
            class="px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 text-zinc-900 rounded-xl text-xs font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          >
            <option :value="2026">2026</option>
            <option :value="2025">2025</option>
          </select>
        </div>
      </div>

      <div v-else class="flex items-center gap-2 text-xs">
        <span class="text-zinc-500 font-bold">Tahun Pajak Fiskal:</span>
        <select
          v-model.number="taxStore.selectedYear"
          class="px-3 py-1.5 bg-zinc-50 border border-zinc-200 text-zinc-900 rounded-xl text-xs font-bold focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        >
          <option :value="2026">Tahun Pajak 2026</option>
          <option :value="2025">Tahun Pajak 2025</option>
        </select>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- VIEW A: PEMBUATAN & PENGELOLAAN SPT MASA (BERBASIS BANK) -->
    <!-- ============================================================== -->
    <div v-if="taxStore.taxPeriodType === 'MASA'" class="space-y-6">
      <!-- Sub-Tabs under SPT Masa -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-zinc-200 pb-3">
        <div class="flex items-center gap-2">
          <button
            @click="sptMasaSubTab = 'DOKUMEN_SPT'"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition flex items-center gap-1.5',
              sptMasaSubTab === 'DOKUMEN_SPT'
                ? 'bg-zinc-950 text-white shadow-xs'
                : 'text-zinc-600 hover:bg-zinc-100'
            ]"
          >
            <FileText class="w-3.5 h-3.5" />
            <span>Dokumen SPT Masa ({{ sptMasaList.length }})</span>
          </button>

          <button
            @click="sptMasaSubTab = 'EBUPOT_DOKTER'"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition flex items-center gap-1.5',
              sptMasaSubTab === 'EBUPOT_DOKTER'
                ? 'bg-zinc-950 text-white shadow-xs'
                : 'text-zinc-600 hover:bg-zinc-100'
            ]"
          >
            <UserCheck class="w-3.5 h-3.5" />
            <span>e-Bupot 21 Dokter (Transfer Bank)</span>
          </button>

          <button
            @click="sptMasaSubTab = 'BUKU_PEMBANTU'"
            :class="[
              'px-3.5 py-1.5 rounded-xl text-xs font-bold transition flex items-center gap-1.5',
              sptMasaSubTab === 'BUKU_PEMBANTU'
                ? 'bg-zinc-950 text-white shadow-xs'
                : 'text-zinc-600 hover:bg-zinc-100'
            ]"
          >
            <FileSpreadsheet class="w-3.5 h-3.5" />
            <span>Buku Pembantu Pajak (Mutasi Bank)</span>
          </button>
        </div>

        <!-- Filter Jenis Pajak -->
        <div class="flex items-center gap-2 text-xs">
          <span class="text-zinc-500 font-bold">Filter Jenis:</span>
          <select
            v-model="taxStore.selectedSptMasaTypeFilter"
            class="px-3 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          >
            <option value="ALL">Semua Jenis Pajak</option>
            <option value="PPH_21">PPh Pasal 21 (Dokter & Staf via Bank)</option>
            <option value="PPH_23">PPh Pasal 23 (Kredit 1771 via Bank)</option>
            <option value="PPH_4_2">PPh Final 4(2) (Giro & Persewaan)</option>
            <option value="PPH_25">PPh Pasal 25 (e-Billing Bank)</option>
            <option value="PPN">PPN (EDC / Kliring)</option>
          </select>
        </div>
      </div>

      <!-- 1-Row Compact Hairline KPI Ticker Strip -->
      <div class="grid grid-cols-2 md:grid-cols-4 divide-y md:divide-y-0 md:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
        <div class="p-4 space-y-1">
          <span class="text-[11px] font-bold text-zinc-500 uppercase tracking-wider block">PPh 21 Dokter (Payroll Bank)</span>
          <div class="tabular-nums text-lg font-bold text-zinc-950">
            {{ formatIDR(sptMasaList.filter(s => s.taxType === 'PPH_21').reduce((sum, s) => sum + s.taxAmount, 0)) }}
          </div>
          <span class="text-[10px] text-zinc-400 block">Transfer Mandiri & BCA Payroll</span>
        </div>

        <div class="p-4 space-y-1">
          <span class="text-[11px] font-bold text-zinc-600 uppercase tracking-wider block">PPh 23 Terpotong (Bank BCA)</span>
          <div class="tabular-nums text-lg font-bold text-zinc-950">
            {{ formatIDR(sptMasaList.filter(s => s.taxType === 'PPH_23').reduce((sum, s) => sum + s.taxAmount, 0)) }}
          </div>
          <span class="text-[10px] text-emerald-600 font-semibold block">&bull; Kredit Lampiran III SPT 1771</span>
        </div>

        <div class="p-4 space-y-1">
          <span class="text-[11px] font-bold text-zinc-600 uppercase tracking-wider block">PPh 25 Disetor (e-Billing)</span>
          <div class="tabular-nums text-lg font-bold text-zinc-950">
            {{ formatIDR(sptMasaList.filter(s => s.taxType === 'PPH_25').reduce((sum, s) => sum + s.taxAmount, 0)) }}
          </div>
          <span class="text-[10px] text-emerald-600 font-semibold block">&bull; Angsuran Induk 10a SPT 1771</span>
        </div>

        <div class="p-4 space-y-1">
          <span class="text-[11px] font-bold text-zinc-500 uppercase tracking-wider block">PPh Final 4(2) & PPN</span>
          <div class="tabular-nums text-lg font-bold text-zinc-950">
            {{ formatIDR(sptMasaList.filter(s => s.taxType === 'PPH_4_2' || s.taxType === 'PPN').reduce((sum, s) => sum + s.taxAmount, 0)) }}
          </div>
          <span class="text-[10px] text-zinc-400 block">Jasa Giro Bank & Sewa Gedung</span>
        </div>
      </div>

      <!-- INLINE TABLE ENTRY: Input SPT Masa Model Baris Tabel (Brayo-Style) -->
      <div
        v-if="showInlineTaxForm"
        class="bg-white rounded-2xl border border-zinc-300 shadow-md p-5 space-y-4 animate-in fade-in slide-in-from-top-2 duration-200"
      >
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <div class="p-1.5 bg-zinc-100 text-zinc-900 rounded-lg">
              <Plus class="w-4 h-4" />
            </div>
            <div>
              <h2 class="text-sm font-bold text-zinc-900">Entri Pajak & SPT Masa (Model Tabel Langsung)</h2>
              <p class="text-[11px] text-zinc-500">
                Pencatatan SPT Masa berbasis tabel baris per Unit Layanan dan kanal Rekening Bank.
              </p>
            </div>
          </div>
          <button
            @click="showInlineTaxForm = false"
            class="p-1 text-zinc-400 hover:text-zinc-600 rounded-lg hover:bg-zinc-100 transition"
            title="Tutup Form"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Inline Table Form -->
        <form @submit.prevent="handleSaveNewSptMasa" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-4 gap-3 bg-zinc-50/80 p-3.5 rounded-xl border border-zinc-200 text-xs">
            <div>
              <label class="block font-bold text-zinc-700 mb-1">Unit Layanan:</label>
              <select
                v-model="newSptMasaForm.unit"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              >
                <option value="KIDINEST">Unit Kidinest (Klinik Spesialis Anak)</option>
                <option value="SSC">Unit Semarang Skin Center (SSC)</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Masa & Tahun Pajak:</label>
              <div class="grid grid-cols-2 gap-2">
                <select
                  v-model.number="newSptMasaForm.masaMonth"
                  class="w-full px-2 py-1.5 bg-white border border-zinc-300 rounded-xl font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
                >
                  <option v-for="m in 12" :key="m" :value="m">Masa {{ m }}</option>
                </select>
                <input
                  v-model.number="newSptMasaForm.masaYear"
                  type="number"
                  class="w-full px-2 py-1.5 bg-white border border-zinc-300 rounded-xl font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
                />
              </div>
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Jenis Pajak:</label>
              <select
                v-model="newSptMasaForm.taxType"
                @change="updateSuggestedTaxAmount"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              >
                <option value="PPH_23">PPh 23 (Kredit Pajak SPT 1771)</option>
                <option value="PPH_25">PPh 25 (Angsuran Bulanan)</option>
                <option value="PPH_21">PPh 21 (Dokter & Karyawan)</option>
                <option value="PPH_4_2">PPh Final 4(2)</option>
                <option value="PPN">PPN (Pajak Pertambahan Nilai)</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Kanal Rekening Bank:</label>
              <select
                v-model="newSptMasaForm.bankAccountCode"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              >
                <option
                  v-for="b in taxStore.availableBankAccounts"
                  :key="b.code"
                  :value="b.code"
                >
                  {{ b.code }} - {{ b.name }}
                </option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-4 gap-3 bg-zinc-50/80 p-3.5 rounded-xl border border-zinc-200 text-xs">
            <div>
              <label class="block font-bold text-zinc-700 mb-1">Nilai DPP / Bruto (Rp):</label>
              <input
                v-model.number="newSptMasaForm.dppAmount"
                @input="updateSuggestedTaxAmount"
                type="number"
                min="0"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl tabular-nums font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Tarif Pajak (%):</label>
              <input
                v-model="newSptMasaForm.taxRate"
                type="text"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl font-mono font-bold text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Nilai Pajak Terutang (Rp):</label>
              <input
                v-model.number="newSptMasaForm.taxAmount"
                type="number"
                min="0"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl tabular-nums font-bold text-rose-700 focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">No. Bukti NTPN / Ref Bank:</label>
              <input
                v-model="newSptMasaForm.ntpnNumber"
                type="text"
                placeholder="Contoh: 1234567890123456"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl font-mono text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              />
            </div>
          </div>

          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-2">
            <div class="flex-1 max-w-lg">
              <input
                v-model="newSptMasaForm.notes"
                type="text"
                placeholder="Uraian objek pajak (contoh: Jasa Perawatan Mesin Laser / Sewa)"
                class="w-full px-3 py-1.5 bg-white border border-zinc-300 rounded-xl text-xs text-zinc-800 focus:ring-1 focus:ring-zinc-950"
              />
            </div>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="showInlineTaxForm = false"
                class="px-3.5 py-1.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-700 rounded-xl text-xs font-bold transition"
              >
                Batal
              </button>
              <button
                type="submit"
                class="flex items-center gap-1.5 px-4 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-xs transition cursor-pointer"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Simpan ke Rekapitulasi Tabel</span>
              </button>
            </div>
          </div>
        </form>
      </div>

      <!-- TAB 1: Daftar Dokumen SPT Masa -->
      <div v-if="sptMasaSubTab === 'DOKUMEN_SPT'" class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
        <div class="p-4 bg-zinc-50/70 border-b border-zinc-200 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div class="flex items-center gap-2">
            <FileText class="w-4 h-4 text-zinc-700" />
            <h3 class="text-xs font-bold text-zinc-800">
              Dokumen SPT Masa Terbitan Resmi (Masa {{ taxStore.selectedMonth }}/{{ taxStore.selectedYear }}) &bull; Berbasis Rekening Bank
            </h3>
          </div>
          <span class="text-[11px] text-zinc-500 font-medium">Hanya mengakui pembayaran/pemotongan non-cash</span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-4 py-3">No. Dokumen SPT</th>
                <th class="px-4 py-3">Masa / Tahun</th>
                <th class="px-4 py-3">Jenis Pajak</th>
                <th class="px-4 py-3">Kanal Rekening Bank (Non-Cash)</th>
                <th class="px-4 py-3">Keterangan / Objek</th>
                <th class="px-4 py-3 text-right">DPP / Bruto</th>
                <th class="px-4 py-3 text-right font-bold text-rose-700">Pajak Terutang</th>
                <th class="px-4 py-3">No. NTPN / Ref Bank</th>
                <th class="px-4 py-3 text-center">Status</th>
                <th class="px-4 py-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr
                v-for="spt in sptMasaList"
                :key="spt.id"
                class="hover:bg-zinc-50/60 transition-colors"
              >
                <td class="px-4 py-3 font-mono font-bold text-zinc-900 whitespace-nowrap">
                  {{ spt.sptNumber }}
                </td>
                <td class="px-4 py-3 font-mono whitespace-nowrap text-zinc-600">
                  Masa {{ spt.masaMonth }} / {{ spt.masaYear }}
                </td>
                <td class="px-4 py-3 whitespace-nowrap">
                  <span
                    v-if="spt.taxType === 'PPH_23'"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-800 border border-zinc-200"
                  >
                    PPh 23 (Kredit 1771)
                  </span>
                  <span
                    v-else-if="spt.taxType === 'PPH_25'"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
                  >
                    PPh 25 (Angsuran 1771)
                  </span>
                  <span
                    v-else-if="spt.taxType === 'PPH_21'"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-purple-50 text-purple-800 border border-purple-200"
                  >
                    PPh 21 Medis & Staf
                  </span>
                  <span
                    v-else
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-800 border border-zinc-200"
                  >
                    {{ spt.taxType }}
                  </span>
                </td>
                <!-- STRICTLY BANK ACCOUNT COLUMN -->
                <td class="px-4 py-3 font-medium text-zinc-800 whitespace-nowrap">
                  <div class="flex items-center gap-1.5">
                    <Landmark class="w-3.5 h-3.5 text-zinc-600 shrink-0" />
                    <span class="font-bold text-[11px]">{{ spt.bankAccountName || 'Bank BCA Operasional' }}</span>
                  </div>
                  <span class="text-[10px] text-zinc-400 font-mono block">Kode COA: {{ spt.bankAccountCode || '1102.001' }}</span>
                </td>
                <td class="px-4 py-3 text-zinc-700 max-w-xs truncate">
                  {{ spt.notes || '-' }}
                </td>
                <td class="px-4 py-3 text-right tabular-nums font-bold text-zinc-800">
                  {{ formatIDR(spt.dppAmount || spt.grossAmount) }}
                </td>
                <td class="px-4 py-3 text-right tabular-nums font-bold text-rose-700">
                  {{ formatIDR(spt.taxAmount) }}
                </td>
                <td class="px-4 py-3 font-mono text-[11px] text-zinc-600 whitespace-nowrap">
                  {{ spt.ntpnNumber || '-' }}
                </td>
                <td class="px-4 py-3 text-center whitespace-nowrap">
                  <span
                    v-if="spt.status === 'REPORTED'"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
                  >
                    Dilaporkan (BPE)
                  </span>
                  <span
                    v-else-if="spt.status === 'PAID'"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-800 border border-zinc-200"
                  >
                    Disetor (NTPN)
                  </span>
                  <span
                    v-else
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-800 border border-amber-200"
                  >
                    Draft
                  </span>
                </td>
                <td class="px-4 py-3 text-center whitespace-nowrap">
                  <button
                    @click="viewSptDetail(spt)"
                    class="px-2.5 py-1 bg-zinc-100 hover:bg-zinc-200 text-zinc-700 rounded-lg font-bold text-[11px] transition"
                  >
                    Detail
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- TAB 2: Bukti Potong Dokter e-Bupot PMK 168/2023 -->
      <div v-else-if="sptMasaSubTab === 'EBUPOT_DOKTER'" class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-5 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div>
            <h3 class="text-sm font-bold text-zinc-900">Perhitungan e-Bupot 21 Tenaga Medis Dokter (PMK 168/2023)</h3>
            <p class="text-xs text-zinc-500 mt-0.5">Penghasilan bruto ditransfer langsung ke rekening bank dokter (DPP 50% Bruto)</p>
          </div>
          <button
            @click="taxStore.exportEbupotCSV()"
            class="px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-xs flex items-center gap-1.5 transition"
          >
            <Download class="w-3.5 h-3.5" />
            <span>Ekspor e-Bupot Format DJP</span>
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Nama Dokter</th>
                <th class="px-3 py-2.5">NPWP / NIK</th>
                <th class="px-3 py-2.5">Kanal Transfer Bank</th>
                <th class="px-3 py-2.5 text-right">Bruto Bulan Ini</th>
                <th class="px-3 py-2.5 text-right">DPP (50%)</th>
                <th class="px-3 py-2.5">Lapisan Tarif</th>
                <th class="px-3 py-2.5 text-right font-bold text-rose-700">PPh 21 Dipotong</th>
                <th class="px-3 py-2.5 text-right font-bold text-emerald-700">Take Home Pay</th>
                <th class="px-3 py-2.5 text-center">No. e-Bupot</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="d in doctorTaxRecords" :key="d.id" class="hover:bg-zinc-50/50">
                <td class="px-3 py-3 font-bold text-zinc-900">{{ d.doctorName }}</td>
                <td class="px-3 py-3 font-mono text-[11px] text-zinc-500">{{ d.npwp }}</td>
                <td class="px-3 py-3 text-[11px] text-zinc-700 font-semibold">{{ d.bankAccountTarget }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-zinc-800">{{ formatIDR(d.grossIncomeCurrentMonth) }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-zinc-900 bg-zinc-50/50">{{ formatIDR(d.dppAmount) }}</td>
                <td class="px-3 py-3 text-[11px] text-zinc-600">{{ d.taxBracketApplied }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-rose-700 bg-rose-50/30">{{ formatIDR(d.pph21Amount) }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-emerald-700 bg-emerald-50/30">{{ formatIDR(d.netPayable) }}</td>
                <td class="px-3 py-3 text-center font-mono text-[10px] text-zinc-500">{{ d.ebupotNumber }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- TAB 3: Buku Pembantu Pajak (Mutasi Bank & GL) -->
      <div v-else class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-5 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div>
            <h3 class="text-sm font-bold text-zinc-900">Buku Pembantu Pajak (Verifikasi Mutasi Bank & Buku Besar)</h3>
            <p class="text-xs text-zinc-500 mt-0.5">Seluruh pemotongan pajak terverifikasi via mutasi Bank BCA & Bank Mandiri</p>
          </div>
          <button
            @click="openLedgerForAccount('25')"
            class="px-3.5 py-1.5 bg-white border border-zinc-200 hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold shadow-2xs flex items-center gap-1.5 transition"
          >
            <BookOpen class="w-3.5 h-3.5 text-zinc-500" />
            <span>Lihat Buku Besar Akun 25 (Hutang Pajak)</span>
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Tanggal</th>
                <th class="px-3 py-2.5">No. Jurnal</th>
                <th class="px-3 py-2.5">Akun Pajak</th>
                <th class="px-3 py-2.5">Lawan Transaksi & Kanal Bank</th>
                <th class="px-3 py-2.5">Objek Pajak</th>
                <th class="px-3 py-2.5 text-right">Nilai Bruto</th>
                <th class="px-3 py-2.5 text-right font-bold text-rose-700">Pajak Terutang</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="tx in txList" :key="tx.id" class="hover:bg-zinc-50/50">
                <td class="px-3 py-3 font-mono text-zinc-600">{{ tx.taxDate }}</td>
                <td class="px-3 py-3 font-mono font-bold text-zinc-900">{{ tx.journalNumber }}</td>
                <td class="px-3 py-3 font-mono font-bold text-zinc-700">{{ tx.accountCode }}</td>
                <td class="px-3 py-3 font-medium text-zinc-900">{{ tx.counterpartyName }}</td>
                <td class="px-3 py-3 text-zinc-600">{{ tx.objectCode }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-zinc-800">{{ formatIDR(tx.grossAmount) }}</td>
                <td class="px-3 py-3 text-right tabular-nums font-bold text-rose-700">{{ formatIDR(tx.taxAmount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- VIEW B: PAJAK TAHUNAN (SPT TAHUNAN PPH BADAN FORMULIR 1771) -->
    <!-- MENGAMBIL NILAI KREDIT PAJAK HANYA DARI MUTASI BANK! -->
    <!-- ============================================================== -->
    <div v-else class="space-y-6">
      <!-- Info Banner: Connected with SPT Masa (Strictly Bank Only) -->
      <div class="bg-zinc-100 border border-zinc-200 rounded-2xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-zinc-950 text-white flex items-center justify-center font-bold shrink-0 shadow-xs">
            <Link2 class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h4 class="font-bold text-zinc-950 text-sm">Kredit Pajak SPT 1771 Bersumber 100% dari Kanal Bank Resmi</h4>
              <span class="px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 font-bold text-[10px]">Non-Cash DJP Standard</span>
            </div>
            <p class="text-zinc-600 mt-0.5">
              Hanya dokumen SPT Masa yang transaksinya terverifikasi melalui <strong>Rekening Bank Operasional (BCA & Mandiri)</strong> yang diakui sebagai kredit pajak PPh 23 & PPh 25. Kas fisik/tunai tanpa rekening koran dikecualikan dari kredit pajak.
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="taxStore.taxPeriodType = 'MASA'"
            class="px-3.5 py-2 bg-white hover:bg-zinc-50 border border-zinc-300 text-zinc-800 rounded-xl font-bold transition shadow-2xs"
          >
            <span>Buka & Kelola SPT Masa Bank</span>
          </button>
        </div>
      </div>

      <!-- Annual Tax Summary Box -->
      <div class="bg-zinc-950 text-white rounded-2xl p-6 shadow-xl border border-zinc-800 space-y-4">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2.5 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-400/30 font-bold text-[11px] flex items-center gap-1">
                <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
                <span>Validasi Perbankan (100% Bank Transferred)</span>
              </span>
              <span class="text-xs text-zinc-400">Tahun Pajak {{ annual1771.fiscalYear }}</span>
            </div>
            <h2 class="text-xl font-black text-white mt-1">
              Perhitungan PPh Badan Terutang & PPh Pasal 29 Kurang Bayar
            </h2>
            <p class="text-xs text-zinc-400">
              Laba Bersih Komersial dikurangi kredit pajak resmi yang dihimpun murni dari mutasi perbankan SPT Masa PPh 23 dan PPh 25.
            </p>
          </div>

          <div class="text-right">
            <span class="text-xs text-zinc-400 block font-medium">PPh Pasal 29 Kurang Bayar:</span>
            <span class="text-2xl font-black tabular-nums text-amber-400">
              {{ formatIDR(annual1771.pph29Payable) }}
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
            <span class="text-emerald-400 font-bold uppercase tracking-wide text-[10px]">4. Kredit Pajak via Bank</span>
            <div class="tabular-nums text-base font-bold text-emerald-300">
              {{ formatIDR(annual1771.taxCredits) }}
            </div>
            <p class="text-[10px] text-emerald-400 font-medium">PPh 23 (Rp 32.5M) + PPh 25 (Rp 143M)</p>
          </div>
        </div>
      </div>

      <!-- CROSS-RECONCILIATION TABLE: Rekapitulasi Aliran Kredit Pajak dari SPT Masa (Strictly Bank Accounts) -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 border-b border-zinc-100 pb-3">
          <div>
            <div class="flex items-center gap-2">
              <Sparkles class="w-4 h-4 text-amber-500" />
              <h3 class="text-sm font-bold text-zinc-900">Rekapitulasi Bulanan SPT Masa Sebagai Sumber Kredit Pajak SPT Tahunan 1771</h3>
            </div>
            <p class="text-xs text-zinc-500 mt-0.5">
              Hanya mengakui penyetoran & pemotongan yang diverifikasi melalui rekening bank resmi (BCA & Mandiri), kas fisik kasir dikecualikan.
            </p>
          </div>
          <span class="text-xs tabular-nums font-bold text-zinc-950 bg-zinc-100 px-3 py-1.5 rounded-xl border border-zinc-200">
            Kredit Bank Terverifikasi: {{ formatIDR(annual1771.taxCredits) }}
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Masa Pajak</th>
                <th class="px-3 py-2.5">Kanal Rekening Bank Pembayar</th>
                <th class="px-3 py-2.5 text-right text-zinc-800">Kredit PPh 23 (Lampiran III)</th>
                <th class="px-3 py-2.5 text-right text-emerald-800">Angsuran PPh 25 (Induk 10a)</th>
                <th class="px-3 py-2.5 text-right text-zinc-600">PPh 21 Medis</th>
                <th class="px-3 py-2.5 text-right text-zinc-600">PPh 4(2) Final</th>
                <th class="px-3 py-2.5 text-right font-bold text-zinc-950 bg-zinc-50/50">Total Kredit SPT 1771</th>
                <th class="px-3 py-2.5 text-center">Status NTPN</th>
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
                  {{ row.pph23Amount > 0 ? formatIDR(row.pph23Amount, false) : '-' }}
                </td>
                <td class="px-3 py-2.5 text-right tabular-nums font-bold text-emerald-700">
                  {{ row.pph25Amount > 0 ? formatIDR(row.pph25Amount, false) : '-' }}
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
                    v-if="row.isReported"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
                  >
                    Valid (NTPN Bank)
                  </span>
                  <span
                    v-else
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-500 border border-zinc-200"
                  >
                    Belum Ada
                  </span>
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-zinc-950 text-white font-bold border-t-2 border-zinc-800">
                <td class="px-3 py-3 uppercase font-sans text-amber-400">TOTAL KREDIT PAJAK TAHUNAN:</td>
                <td class="px-3 py-3 font-sans text-zinc-400 text-xs">Rekap 100% Kanal Bank</td>
                <td class="px-3 py-3 text-right tabular-nums text-white">{{ formatIDR(annual1771.taxCreditPph23) }}</td>
                <td class="px-3 py-3 text-right tabular-nums text-emerald-400">{{ formatIDR(annual1771.taxCreditPph25) }}</td>
                <td class="px-3 py-3 text-right text-zinc-400">-</td>
                <td class="px-3 py-3 text-right text-zinc-400">-</td>
                <td class="px-3 py-3 text-right tabular-nums text-amber-300 text-sm font-black bg-zinc-900">{{ formatIDR(annual1771.taxCredits) }}</td>
                <td class="px-3 py-3 text-center text-emerald-400 font-sans text-xs">Terverifikasi</td>
              </tr>
            </tfoot>
          </table>
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
                <td class="px-4 py-2.5 text-zinc-500 font-sans">Buku Besar: Penjualan Medis & Jasa Giro Bank</td>
                <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-900">{{ formatIDR(annual1771.commercialRevenue) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">2</td>
                <td class="px-4 py-2.5 font-sans">Harga Pokok Penjualan (HPP Akun 51)</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">Persediaan Obat Farmasi & Biaya Medis Langsung</td>
                <td class="px-4 py-2.5 text-right tabular-nums text-zinc-700">({{ formatIDR(9807287775.40) }})</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">3</td>
                <td class="px-4 py-2.5 font-sans">Biaya Usaha & Operasional (Akun 52 & 7)</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">Gaji, Listrik, Sewa, Limbah B3 & Bunga Pinjaman</td>
                <td class="px-4 py-2.5 text-right tabular-nums text-zinc-700">({{ formatIDR(3350282603.47) }})</td>
              </tr>
              <tr class="bg-zinc-50 font-bold">
                <td class="px-4 py-2.5">4</td>
                <td class="px-4 py-2.5 font-sans">Penghasilan Netto Komersial Sebelum Pajak</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">Laba Bersih Buku Besar Komersial</td>
                <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950">{{ formatIDR(annual1771.commercialNetProfit) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">5</td>
                <td class="px-4 py-2.5 font-sans text-emerald-700 font-semibold">Koreksi Fiskal Positif</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.positiveCorrectionNotes }}</td>
                <td class="px-4 py-2.5 text-right tabular-nums text-emerald-700 font-bold">+ {{ formatIDR(annual1771.positiveFiscalCorrection) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">6</td>
                <td class="px-4 py-2.5 font-sans text-rose-700 font-semibold">Koreksi Fiskal Negatif</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.negativeCorrectionNotes }}</td>
                <td class="px-4 py-2.5 text-right tabular-nums text-rose-700 font-bold">- {{ formatIDR(annual1771.negativeFiscalCorrection) }}</td>
              </tr>
              <tr class="bg-amber-50/40 font-bold border-y-2 border-amber-200">
                <td class="px-4 py-3">7</td>
                <td class="px-4 py-3 font-sans text-amber-950 font-bold">Penghasilan Kena Pajak (PKP)</td>
                <td class="px-4 py-3 text-zinc-600 font-sans font-normal">Dasar Pengenaan Pajak Penghasilan Badan</td>
                <td class="px-4 py-3 text-right tabular-nums text-amber-900 text-sm font-black">{{ formatIDR(annual1771.taxableIncome) }}</td>
              </tr>
              <tr>
                <td class="px-4 py-2.5">8</td>
                <td class="px-4 py-2.5 font-sans font-bold">PPh Badan Terutang (Tahun Berjalan)</td>
                <td class="px-4 py-2.5 text-zinc-500 font-sans">Fasilitas Pasal 31E UU PPh (11% x PKP)</td>
                <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">{{ formatIDR(annual1771.corporateTaxPayable) }}</td>
              </tr>
              <!-- LINE 9: KREDIT PAJAK HANYA DARI BANK -->
              <tr class="bg-emerald-50/50 border border-emerald-200">
                <td class="px-4 py-3 text-emerald-900 font-bold">9</td>
                <td class="px-4 py-3 font-sans text-emerald-950">
                  <div class="font-bold">Kredit Pajak Dalam Negeri (100% Berbasis Bank)</div>
                  <div class="text-[11px] text-emerald-700 mt-0.5">
                    &bull; PPh 23 dari SPT Masa via Bank: {{ formatIDR(annual1771.taxCreditPph23) }} <br/>
                    &bull; PPh 25 Angsuran Badan via e-Billing Bank: {{ formatIDR(annual1771.taxCreditPph25) }}
                  </div>
                </td>
                <td class="px-4 py-3 text-emerald-800 font-sans">
                  Total akumulasi 11 masa pajak PPh 23 & PPh 25 yang disetor via rekening Bank BCA/Mandiri (Kas fisik dikecualikan)
                </td>
                <td class="px-4 py-3 text-right tabular-nums font-bold text-emerald-800 text-sm">({{ formatIDR(annual1771.taxCredits) }})</td>
              </tr>
              <!-- LINE 10: PPH 29 KURANG BAYAR -->
              <tr class="bg-zinc-950 text-white font-bold text-sm">
                <td class="px-4 py-3 text-amber-400">10</td>
                <td class="px-4 py-3 font-sans text-white font-bold">PPh Kurang Bayar (PPh Pasal 29) - Setor Kas Negara via Bank</td>
                <td class="px-4 py-3 text-zinc-300 font-sans font-normal">Wajib disetor via Rekening Bank sebelum lapor SPT 1771 (Akun COA 2506)</td>
                <td class="px-4 py-3 text-right tabular-nums text-amber-400 font-black">{{ formatIDR(annual1771.pph29Payable) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="pt-4 flex items-center justify-end gap-3">
          <button
            @click="showPaymentBankModal = true"
            class="px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-xs flex items-center gap-1.5"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Generate Jurnal Pelunasan PPh 29 (via Bank)</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Official Report Signature Section matching Image 1 -->
    <ReportSignatures />

    <!-- MODAL 1: Buat SPT Masa Baru (Strictly Bank Accounts Only) -->
    <div
      v-if="showCreateSptMasaModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-zinc-950 text-white flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Plus class="w-4 h-4 text-zinc-400" />
            <h3 class="text-sm font-bold text-white">Buat Dokumen SPT Masa Baru (Kanal Bank)</h3>
          </div>
          <button @click="showCreateSptMasaModal = false" class="text-zinc-400 hover:text-white transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveNewSptMasa" class="p-6 space-y-4 text-xs">
          <!-- Mandatory Bank Account Selector -->
          <div class="p-3 bg-zinc-50 border border-zinc-200 rounded-2xl space-y-1">
            <label class="block font-bold text-zinc-900">Rekening Bank Pembayar / Persepsi DJP (Wajib Bank, Bukan Kas)</label>
            <select
              v-model="newSptMasaForm.bankAccountCode"
              class="w-full px-3 py-2 bg-white border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 font-bold text-xs"
              required
            >
              <option
                v-for="b in taxStore.availableBankAccounts"
                :key="b.code"
                :value="b.code"
              >
                {{ b.code }} - {{ b.name }}
              </option>
            </select>
            <p class="text-[10px] text-zinc-500">
              &bull; Sesuai ketentuan DJP, penyetoran & pemotongan pajak hanya diakui jika diproses melalui Rekening Bank resmi. Kas fisik kasir/petty cash dilarang digunakan.
            </p>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-700 mb-1">Jenis Pajak</label>
              <select
                v-model="newSptMasaForm.taxType"
                @change="updateSuggestedTaxAmount"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 font-bold"
              >
                <option value="PPH_23">PPh Pasal 23 (Kredit Pajak SPT 1771)</option>
                <option value="PPH_25">PPh Pasal 25 (Angsuran PPh Badan)</option>
                <option value="PPH_21">PPh Pasal 21 (Dokter & Karyawan)</option>
                <option value="PPH_4_2">PPh Final Pasal 4(2)</option>
                <option value="PPN">PPN (Pajak Pertambahan Nilai)</option>
              </select>
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Masa Pajak / Bulan</label>
              <select
                v-model.number="newSptMasaForm.masaMonth"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950"
              >
                <option v-for="m in 12" :key="m" :value="m">Masa {{ m }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-700 mb-1">Tahun Pajak</label>
              <input
                v-model.number="newSptMasaForm.masaYear"
                type="number"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 font-mono font-bold"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Status Dokumen</label>
              <select
                v-model="newSptMasaForm.status"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 font-semibold"
              >
                <option value="REPORTED">Dilaporkan (BPE)</option>
                <option value="PAID">Disetor (NTPN)</option>
                <option value="DRAFT">Draft</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-zinc-700 mb-1">DPP / Nilai Bruto (Rp)</label>
              <input
                v-model.number="newSptMasaForm.dppAmount"
                @input="updateSuggestedTaxAmount"
                type="number"
                min="0"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 tabular-nums font-bold"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-zinc-700 mb-1">Pajak Terutang (Rp)</label>
              <input
                v-model.number="newSptMasaForm.taxAmount"
                type="number"
                min="0"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 tabular-nums font-bold text-rose-700"
                required
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Nomor NTPN Bank / Bukti Setor DJP (Opsional)</label>
            <input
              v-model="newSptMasaForm.ntpnNumber"
              type="text"
              placeholder="Contoh: NTPN-BCA-202611-9988"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950 font-mono"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Uraian / Keterangan Objek Pajak</label>
            <input
              v-model="newSptMasaForm.notes"
              type="text"
              placeholder="Contoh: Jasa kalibrasi alat alkes / Jasa limbah B3"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl focus:ring-1 focus:ring-zinc-950"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-zinc-200">
            <button
              @click="showCreateSptMasaModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-zinc-300 rounded-xl font-semibold text-zinc-700 hover:bg-zinc-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-xs transition"
            >
              Simpan SPT Masa (via Bank)
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 2: Dialog Pilih Rekening Bank untuk Penyetoran Pajak NTPN -->
    <div
      v-if="showPaymentBankModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-md w-full overflow-hidden p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-200 pb-3">
          <div class="flex items-center gap-2">
            <Landmark class="w-4 h-4 text-emerald-600" />
            <h3 class="text-sm font-bold text-zinc-950">Pilih Rekening Bank Penyetoran Pajak NTPN</h3>
          </div>
          <button @click="showPaymentBankModal = false" class="text-zinc-400 hover:text-zinc-600 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Rekening Bank Sumber Dana (Kredit):</label>
            <select
              v-model="selectedBankForPayment"
              class="w-full px-3 py-2.5 bg-zinc-50 border border-zinc-300 rounded-xl font-bold text-xs focus:ring-1 focus:ring-zinc-950"
            >
              <option
                v-for="b in taxStore.availableBankAccounts"
                :key="b.code"
                :value="b.code"
              >
                {{ b.code }} - {{ b.name }}
              </option>
            </select>
          </div>

          <div class="p-3 bg-zinc-50 border border-zinc-200 rounded-2xl space-y-1 text-zinc-800">
            <span class="font-bold block text-[11px] text-zinc-950">Validasi Rekening Bank:</span>
            <p class="text-[11px] text-zinc-600">
              Jurnal pengeluaran pelunasan pajak akan mengkredit rekening Bank di atas dan mendebit akun Hutang Pajak. Akun kas fisik kasir dikecualikan.
            </p>
          </div>
        </div>

        <div class="flex justify-end gap-2 pt-3 border-t border-zinc-200">
          <button
            @click="showPaymentBankModal = false"
            class="px-4 py-2 border border-zinc-300 rounded-xl font-bold text-zinc-700 text-xs hover:bg-zinc-100 transition"
          >
            Batal
          </button>
          <button
            @click="handleExecutePayment"
            class="px-5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold text-xs shadow-xs transition"
          >
            Proses Jurnal Setor Pajak
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 3: Detail Dokumen SPT Masa -->
    <div
      v-if="showDetailModal && selectedSptDetail"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-lg w-full overflow-hidden p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-200 pb-3">
          <div>
            <span class="text-[10px] font-bold text-zinc-500 uppercase tracking-wider">Dokumen SPT Masa DJP (Bank Channel)</span>
            <h3 class="text-base font-black text-zinc-950">{{ selectedSptDetail.sptNumber }}</h3>
          </div>
          <button @click="showDetailModal = false" class="text-zinc-400 hover:text-zinc-600 transition">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="grid grid-cols-2 gap-3 text-xs">
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">Jenis Pajak:</span>
            <span class="font-bold text-zinc-800">{{ selectedSptDetail.taxType }}</span>
          </div>
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">Masa / Tahun:</span>
            <span class="font-bold text-zinc-800">Masa {{ selectedSptDetail.masaMonth }} / {{ selectedSptDetail.masaYear }}</span>
          </div>
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">Kanal Bank:</span>
            <span class="font-mono font-bold text-zinc-900">{{ selectedSptDetail.bankAccountName || 'Bank BCA Operasional' }}</span>
          </div>
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">Status:</span>
            <span class="font-bold text-emerald-700">{{ selectedSptDetail.status }}</span>
          </div>
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">No. NTPN:</span>
            <span class="font-mono font-bold text-zinc-800">{{ selectedSptDetail.ntpnNumber || '-' }}</span>
          </div>
          <div>
            <span class="text-zinc-400 text-[10px] uppercase font-bold block">No. BPE:</span>
            <span class="font-mono text-zinc-600">{{ selectedSptDetail.bpeNumber || '-' }}</span>
          </div>
        </div>

        <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-2xl space-y-1.5 text-xs">
          <div class="flex justify-between font-mono">
            <span class="text-zinc-600">Nilai Bruto / DPP:</span>
            <span class="tabular-nums font-bold text-zinc-900">{{ formatIDR(selectedSptDetail.dppAmount || selectedSptDetail.grossAmount) }}</span>
          </div>
          <div class="flex justify-between font-mono">
            <span class="text-zinc-600">Tarif Pajak:</span>
            <span class="font-bold text-zinc-900">{{ selectedSptDetail.taxRate }}</span>
          </div>
          <div class="flex justify-between font-mono pt-1.5 border-t border-zinc-200 text-sm">
            <span class="font-bold text-zinc-800">Pajak Terutang:</span>
            <span class="tabular-nums font-black text-rose-700">{{ formatIDR(selectedSptDetail.taxAmount) }}</span>
          </div>
        </div>

        <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-2xl text-xs text-zinc-800">
          <span class="font-bold block text-[11px] text-zinc-900">Keterangan & Catatan Bank:</span>
          <p class="mt-0.5 text-zinc-600">{{ selectedSptDetail.notes || 'Pembayaran diverifikasi via mutasi rekening bank.' }}</p>
        </div>

        <div class="flex justify-end pt-2 border-t border-zinc-200">
          <button
            @click="showDetailModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl font-bold text-xs hover:bg-zinc-800 transition"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
