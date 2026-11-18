<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useJournalsStore } from '@/stores/journals'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportJournalListToExcel, exportJournalListToCSV } from '@/utils/export'
import StatusBadge from '@/components/common/StatusBadge.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import type { JournalEntry, JournalStatus, BusinessUnit } from '@/types/accounting'
import {
  BookOpen,
  Search,
  Plus,
  Download,
  Eye,
  CheckCircle2,
  FileCheck2,
  Send,
  RotateCcw,
  Database,
  Trash2,
  AlertCircle,
  X,
  Sparkles
} from 'lucide-vue-next'

const route = useRoute()
const journalsStore = useJournalsStore()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await Promise.all([
    journalsStore.fetchJournals(),
    ledgerStore.fetchAccounts()
  ])
})

// URL query support (e.g. ?status=DRAFT)
const initialStatus = (route.query.status as string) || 'ALL'
const selectedStatusTab = ref<string>(initialStatus)
const searchQuery = ref('')
const selectedSource = ref<string>('ALL')

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)
const showInlineForm = ref(false)

// Inline Table Entry State (Brayo-Style Direct Table Input)
interface InlineRow {
  unit: 'KIDINEST' | 'SSC'
  accountCode: string
  accountName: string
  description: string
  debit: number
  credit: number
}

const newJournalDate = ref(new Date().toISOString().split('T')[0])
const newJournalUnit = ref<BusinessUnit>(accountingStore.businessUnit === 'SSC' ? 'SSC' : 'KIDINEST')
const newJournalRef = ref('')
const newJournalDesc = ref('')

watch(
  () => accountingStore.businessUnit,
  (unit) => {
    if (unit === 'KIDINEST' || unit === 'SSC') {
      newJournalUnit.value = unit
    }
  }
)

const inlineRows = ref<InlineRow[]>([
  {
    unit: 'KIDINEST',
    accountCode: '1-1100',
    accountName: 'Kas Operasional Kidinest',
    description: '',
    debit: 0,
    credit: 0
  },
  {
    unit: 'KIDINEST',
    accountCode: '4-1100',
    accountName: 'Pendapatan Jasa Medis Kidinest',
    description: '',
    debit: 0,
    credit: 0
  }
])

const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

const accountSelectOptions = computed(() => {
  return availableAccounts.value.map(acc => ({
    value: acc.code,
    label: `[${acc.code}] ${acc.name}`,
    sublabel: `${acc.type} • Normal: ${acc.normalBalance}`
  }))
})

function onAccountSelect(index: number, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc) {
    inlineRows.value[index].accountCode = acc.code
    inlineRows.value[index].accountName = acc.name
  }
}

function addInlineRow() {
  inlineRows.value.push({
    unit: newJournalUnit.value === 'SSC' ? 'SSC' : 'KIDINEST',
    accountCode: '',
    accountName: '',
    description: '',
    debit: 0,
    credit: 0
  })
}

function removeInlineRow(index: number) {
  if (inlineRows.value.length > 2) {
    inlineRows.value.splice(index, 1)
  } else {
    accountingStore.showToast('Jurnal minimal memiliki 2 baris (Debit & Kredit)', 'error')
  }
}

const totalInlineDebit = computed(() => {
  return inlineRows.value.reduce((s, r) => s + (Number(r.debit) || 0), 0)
})

const totalInlineCredit = computed(() => {
  return inlineRows.value.reduce((s, r) => s + (Number(r.credit) || 0), 0)
})

const isInlineBalanced = computed(() => {
  return totalInlineDebit.value > 0 && Math.abs(totalInlineDebit.value - totalInlineCredit.value) < 0.01
})

const inlineBalanceDiff = computed(() => {
  return Math.abs(totalInlineDebit.value - totalInlineCredit.value)
})

function autoBalanceCreditRow() {
  if (inlineRows.value.length >= 2 && totalInlineDebit.value > 0) {
    const lastIdx = inlineRows.value.length - 1
    const otherCredit = inlineRows.value.reduce((s, r, idx) => idx === lastIdx ? s : s + (Number(r.credit) || 0), 0)
    const target = totalInlineDebit.value - otherCredit
    if (target > 0) {
      inlineRows.value[lastIdx].credit = target
      accountingStore.showToast(`Baris kredit terakhir disesuaikan menjadi ${formatIDR(target)}`, 'info')
    }
  }
}

function handleSaveInlineJournal(status: JournalStatus = 'POSTED') {
  if (!isInlineBalanced.value) {
    accountingStore.showToast('Jurnal belum seimbang! Total Debit harus sama persis dengan Total Kredit.', 'error')
    return
  }
  if (!newJournalDesc.value.trim()) {
    accountingStore.showToast('Harap isi keterangan memo umum jurnal.', 'error')
    return
  }
  for (const r of inlineRows.value) {
    if (!r.accountCode) {
      accountingStore.showToast('Harap pilih akun COA untuk semua baris.', 'error')
      return
    }
  }

  try {
    journalsStore.createManualJournal({
      journalDate: newJournalDate.value,
      businessUnit: newJournalUnit.value,
      description: newJournalDesc.value.trim(),
      sourceReference: newJournalRef.value.trim() || 'MANUAL-ENTRY',
      status,
      lines: inlineRows.value.map(r => ({
        accountCode: r.accountCode,
        accountName: r.accountName,
        debit: Number(r.debit) || 0,
        credit: Number(r.credit) || 0,
        description: r.description.trim() || newJournalDesc.value.trim()
      }))
    })

    accountingStore.showToast(
      `Jurnal berhasil dibuat (${status === 'POSTED' ? 'Langsung Terposting ke Buku Besar' : 'Draft Tersimpan'})!`,
      'success'
    )

    // Reset Form
    newJournalDesc.value = ''
    newJournalRef.value = ''
    inlineRows.value = [
      { unit: 'KIDINEST', accountCode: '1-1100', accountName: 'Kas Operasional Kidinest', description: '', debit: 0, credit: 0 },
      { unit: 'KIDINEST', accountCode: '4-1100', accountName: 'Pendapatan Jasa Medis Kidinest', description: '', debit: 0, credit: 0 }
    ]
    showInlineForm.value = false
  } catch (err: any) {
    accountingStore.showToast(err.message || 'Gagal menyimpan jurnal', 'error')
  }
}

function inspectJournal(j: JournalEntry) {
  selectedJournal.value = j
  showDetailModal.value = true
}

// Status tab options
const tabs = [
  { id: 'ALL', label: 'Semua Jurnal' },
  { id: 'DRAFT', label: 'Draft' },
  { id: 'SUBMITTED', label: 'Submitted' },
  { id: 'REVIEWED', label: 'Reviewed' },
  { id: 'APPROVED', label: 'Approved' },
  { id: 'POSTED', label: 'Posted' },
  { id: 'REVERSED', label: 'Reversed' }
]

const filteredJournals = computed(() => {
  return journalsStore.filteredJournals.filter(j => {
    // Status tab filter
    if (selectedStatusTab.value !== 'ALL' && j.status !== selectedStatusTab.value) {
      return false
    }
    // Source filter
    if (selectedSource.value !== 'ALL' && j.journalSource !== selectedSource.value) {
      return false
    }
    // Search query
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchNo = j.journalNumber.toLowerCase().includes(q)
      const matchDesc = j.description.toLowerCase().includes(q)
      const matchRef = j.sourceReference.toLowerCase().includes(q)
      const matchLines = j.lines.some(l => l.accountName.toLowerCase().includes(q) || l.accountCode.toLowerCase().includes(q))
      if (!matchNo && !matchDesc && !matchRef && !matchLines) return false
    }
    return true
  })
})

function handleExportExcel() {
  exportJournalListToExcel(filteredJournals.value, {
    businessUnit: accountingStore.businessUnit,
    fiscalPeriod: `Masa ${String(accountingStore.fiscalMonth).padStart(2, '0')}/2026`
  })
  accountingStore.showToast('Buku Jurnal Umum (Multi-Line Double Entry) berhasil diekspor ke Excel (.xlsx)', 'success')
}

function handleExportCSV() {
  exportJournalListToCSV(filteredJournals.value, {
    businessUnit: accountingStore.businessUnit,
    fiscalPeriod: `Masa ${String(accountingStore.fiscalMonth).padStart(2, '0')}/2026`
  })
  accountingStore.showToast('Buku Jurnal Umum (Format CSV Akuntansi) berhasil diunduh!', 'success')
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Buku Jurnal Umum
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Double-Entry Ledger Engine</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Daftar Jurnal Akuntansi
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <router-link
          to="/accounting/journals/create"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition-all shadow-md active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5 text-emerald-400" />
          <span>Buat Jurnal Baru (New Page)</span>
        </router-link>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Ekspor Buku Jurnal Umum Lengkap beserta Rincian Debet Kredit Akun ke Excel"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>

        <button
          @click="handleExportCSV"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Ekspor Buku Jurnal Umum Double-Entry ke Format CSV"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor CSV (.csv)</span>
        </button>
      </div>
    </header>

    <!-- Form Entri Jurnal Umum Langsung (Spreadsheet-Style Table) -->
    <div
      v-if="showInlineForm"
      class="bg-white rounded-2xl border border-zinc-300 shadow-lg p-5 space-y-4 animate-in fade-in slide-in-from-top-2 duration-200"
    >
      <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
        <div class="flex items-center gap-2.5">
          <div class="p-1.5 bg-zinc-100 text-zinc-900 rounded-lg">
            <Plus class="w-4 h-4" />
          </div>
          <div>
            <h2 class="text-sm font-bold text-zinc-950">Entri Jurnal Umum Cepat (Mode Spreadsheet)</h2>
            <p class="text-[11px] text-zinc-500">
              Input double-entry langsung pada baris tabel. Total Debit wajib sama persis dengan Kredit.
            </p>
          </div>
        </div>
        <button
          @click="showInlineForm = false"
          class="p-1.5 text-zinc-400 hover:text-zinc-700 rounded-lg hover:bg-zinc-100 transition-colors cursor-pointer"
          title="Tutup Form"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Header Transaksi Form -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-3 bg-zinc-50 p-3.5 rounded-xl border border-zinc-200/90 text-xs">
        <div>
          <label class="block font-bold text-zinc-700 mb-1">Tanggal Transaksi:</label>
          <input
            v-model="newJournalDate"
            type="date"
            class="w-full px-3 py-1.5 bg-white border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          />
        </div>
        <div>
          <label class="block font-bold text-zinc-700 mb-1">Unit Layanan Utama:</label>
          <select
            v-model="newJournalUnit"
            class="w-full px-3 py-1.5 bg-white border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
          >
            <option value="KIDINEST">Kidinest (Klinik Spesialis Anak)</option>
            <option value="SSC">Semarang Skin Center (SSC)</option>
            <option value="ALL">Semua Unit (Konsolidasi)</option>
          </select>
        </div>
        <div>
          <label class="block font-bold text-zinc-700 mb-1">No. Bukti / Dokumen:</label>
          <input
            v-model="newJournalRef"
            type="text"
            placeholder="Contoh: BKK-2026/03/001 atau INV-092"
            class="w-full px-3 py-1.5 bg-white border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          />
        </div>
        <div>
          <label class="block font-bold text-zinc-700 mb-1">Keterangan Memo Umum:</label>
          <input
            v-model="newJournalDesc"
            type="text"
            placeholder="Contoh: Pembelian BHP medis operasional"
            class="w-full px-3 py-1.5 bg-white border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none font-medium"
          />
        </div>
      </div>

      <!-- Spreadsheet-Style Lines Table -->
      <div class="border border-zinc-200/90 rounded-xl overflow-hidden shadow-2xs">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="py-2.5 px-3 w-10 text-center">#</th>
              <th class="py-2.5 px-3 w-36">Unit</th>
              <th class="py-2.5 px-3 w-72">Akun Rekening (COA)</th>
              <th class="py-2.5 px-3">Uraian Baris (Opsional)</th>
              <th class="py-2.5 px-3 w-44 text-right">Debit (Rp)</th>
              <th class="py-2.5 px-3 w-44 text-right">Kredit (Rp)</th>
              <th class="py-2.5 px-3 w-12 text-center">Hapus</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="(row, idx) in inlineRows"
              :key="idx"
              class="hover:bg-zinc-50/70"
            >
              <td class="py-2 px-3 text-center text-zinc-400 font-mono font-bold">{{ idx + 1 }}</td>
              <td class="py-2 px-3">
                <select
                  v-model="row.unit"
                  class="w-full px-2 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
                >
                  <option value="KIDINEST">Kidinest</option>
                  <option value="SSC">SSC</option>
                </select>
              </td>
              <td class="py-2 px-3 min-w-[260px]">
                <SearchableSelect
                  :model-value="row.accountCode"
                  :options="accountSelectOptions"
                  placeholder="Pilih atau cari akun COA..."
                  @update:model-value="(val: any) => onAccountSelect(idx, String(val))"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model="row.description"
                  type="text"
                  placeholder="Keterangan baris (opsional)"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model.number="row.debit"
                  type="number"
                  min="0"
                  step="1000"
                  placeholder="0"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs tabular-nums font-bold text-right text-zinc-950 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model.number="row.credit"
                  type="number"
                  min="0"
                  step="1000"
                  placeholder="0"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs tabular-nums font-bold text-right text-zinc-950 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3 text-center">
                <button
                  @click="removeInlineRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                  title="Hapus baris ini"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          </tbody>
          <!-- Table Footer Running Totals -->
          <tfoot>
            <tr class="bg-zinc-50 border-t border-zinc-200/90 font-bold text-zinc-900">
              <td colspan="4" class="py-3 px-4">
                <div class="flex items-center justify-between">
                  <button
                    @click="addInlineRow"
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-100 text-zinc-800 rounded-xl border border-zinc-200 font-bold shadow-2xs transition-colors cursor-pointer"
                  >
                    <Plus class="w-3.5 h-3.5" />
                    <span>Tambah Baris (+Row)</span>
                  </button>
                  <span class="text-xs uppercase tracking-wider text-zinc-500 font-bold">Total Nilai Jurnal:</span>
                </div>
              </td>
              <td class="py-3 px-3 text-right tabular-nums font-extrabold text-zinc-950 bg-zinc-100/60">
                {{ formatIDR(totalInlineDebit) }}
              </td>
              <td class="py-3 px-3 text-right tabular-nums font-extrabold text-zinc-950 bg-zinc-100/60">
                {{ formatIDR(totalInlineCredit) }}
              </td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Action & Balance Verification Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-2">
        <div class="flex items-center gap-2.5">
          <div
            v-if="isInlineBalanced"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
          >
            <CheckCircle2 class="w-4 h-4 text-emerald-600" />
            <span>SEIMBANG ({{ formatIDR(totalInlineDebit) }})</span>
          </div>
          <div
            v-else
            class="inline-flex items-center gap-2"
          >
            <div class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold bg-rose-50 text-rose-800 border border-rose-200">
              <AlertCircle class="w-4 h-4 text-rose-600" />
              <span>BELUM SEIMBANG (Selisih: {{ formatIDR(inlineBalanceDiff) }})</span>
            </div>
            <button
              v-if="totalInlineDebit > 0 && !isInlineBalanced"
              @click="autoBalanceCreditRow"
              class="px-2.5 py-1 text-xs font-bold bg-amber-500 hover:bg-amber-600 text-white rounded-xl shadow-2xs cursor-pointer transition-all"
            >
              Seimbangkan Kredit Otomatis
            </button>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="showInlineForm = false"
            class="px-3 py-2 bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 rounded-xl text-xs font-bold transition-colors shadow-2xs cursor-pointer"
          >
            Batal
          </button>
          <button
            @click="handleSaveInlineJournal('DRAFT')"
            :disabled="!isInlineBalanced"
            :class="[
              'px-3.5 py-2 rounded-xl text-xs font-bold transition-colors shadow-2xs',
              isInlineBalanced
                ? 'bg-amber-500 hover:bg-amber-600 text-white cursor-pointer'
                : 'bg-zinc-100 text-zinc-400 border border-zinc-200 cursor-not-allowed'
            ]"
          >
            Simpan Draft
          </button>
          <button
            @click="handleSaveInlineJournal('POSTED')"
            :disabled="!isInlineBalanced"
            :class="[
              'inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold transition-colors shadow-2xs active:scale-95',
              isInlineBalanced
                ? 'bg-zinc-950 hover:bg-zinc-800 text-white cursor-pointer'
                : 'bg-zinc-100 text-zinc-400 border border-zinc-200 cursor-not-allowed'
            ]"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Simpan & Posting ke Buku Besar</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Status Sub-Tabs -->
    <div class="flex items-center bg-zinc-100 p-0.5 rounded-xl border border-zinc-200/90 w-fit">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="selectedStatusTab = tab.id"
        :class="[
          'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all whitespace-nowrap',
          selectedStatusTab === tab.id
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950'
        ]"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Toolbar Filters & Search -->
    <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-80">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari no jurnal, uraian, akun, referensi..."
            class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
          />
        </div>

        <!-- Source Selector -->
        <div class="relative">
          <select
            v-model="selectedSource"
            class="px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-medium text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
          >
            <option value="ALL">Semua Sumber Jurnal</option>
            <option value="CASHIER_CLOSING">Closing Kasir</option>
            <option value="PHARMACY_PURCHASE">Pembelian Farmasi (PO)</option>
            <option value="DOCTOR_FEE">Jasa Medis Dokter</option>
            <option value="FIXED_ASSET">Penyusutan Aset</option>
            <option value="INVENTORY_CONSUMPTION">Pemakaian BHP</option>
            <option value="MANUAL">Jurnal Penyesuaian Manual</option>
          </select>
        </div>
      </div>

      <div class="text-xs text-zinc-500 font-medium">
        Menampilkan <strong>{{ filteredJournals.length }}</strong> jurnal
      </div>
    </div>

    <!-- Journals Data Table -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3">No Jurnal</th>
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">Sumber & Ref</th>
              <th class="px-4 py-3">Unit Bisnis</th>
              <th class="px-4 py-3">Keterangan Memo</th>
              <th class="px-4 py-3 text-right">Total Nilai (Rp)</th>
              <th class="px-4 py-3 text-center">Status</th>
              <th class="px-4 py-3 text-center w-32">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="j in filteredJournals"
              :key="j.id"
              class="hover:bg-zinc-50/70 transition-colors"
            >
              <td class="px-4 py-3.5 font-mono font-bold text-zinc-950 whitespace-nowrap">
                {{ j.journalNumber }}
              </td>
              <td class="px-4 py-3.5 text-zinc-600 whitespace-nowrap">
                {{ formatDate(j.journalDate) }}
              </td>
              <td class="px-4 py-3.5">
                <span class="block font-semibold text-zinc-900">{{ j.journalSource }}</span>
                <span class="text-[10px] text-zinc-400 font-mono">{{ j.sourceReference }}</span>
              </td>
              <td class="px-4 py-3.5">
                <StatusBadge :status="j.businessUnit" size="sm" />
              </td>
              <td class="px-4 py-3.5 max-w-sm">
                <p class="font-medium text-zinc-900 truncate">{{ j.description }}</p>
                <p class="text-[10px] text-zinc-400">
                  {{ j.lines.length }} Baris Akun &bull; Dibuat oleh {{ j.createdByName }}
                </p>
              </td>
              <td class="px-4 py-3.5 text-right tabular-nums font-bold text-zinc-950 whitespace-nowrap">
                {{ formatIDR(j.totalDebit) }}
              </td>
              <td class="px-4 py-3.5 text-center">
                <StatusBadge :status="j.status" size="sm" />
              </td>
              <td class="px-4 py-3.5 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    @click="inspectJournal(j)"
                    class="px-2.5 py-1 bg-white hover:bg-zinc-100 text-zinc-700 rounded-lg text-[11px] font-bold border border-zinc-200 shadow-2xs transition-colors flex items-center gap-1 cursor-pointer"
                  >
                    <Eye class="w-3.5 h-3.5 text-zinc-500" />
                    <span>Periksa</span>
                  </button>

                  <button
                    v-if="j.status === 'APPROVED' && accountingStore.canPostJournal"
                    @click="journalsStore.postJournal(j.id)"
                    class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-[11px] font-bold shadow-2xs transition-colors cursor-pointer"
                    title="Posting ke Buku Besar"
                  >
                    Post
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredJournals.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-zinc-400">
                Tidak ada jurnal akuntansi yang ditemukan pada kategori ini.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modals -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />
  </div>
</template>

