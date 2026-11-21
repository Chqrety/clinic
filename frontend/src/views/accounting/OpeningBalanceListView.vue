<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel, exportToCSV } from '@/utils/export'
import type { JournalEntry, BusinessUnit } from '@/types/accounting'
import {
  Folder,
  Plus,
  Search,
  Download,
  Filter,
  CheckCircle2,
  AlertTriangle,
  Scale,
  Calendar,
  Layers,
  ArrowRight,
  Eye,
  FileText,
  Building,
  RefreshCw,
  Sparkles,
  ExternalLink,
  ShieldCheck,
  X
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

const loading = ref(false)
const searchQuery = ref('')
const selectedUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const selectedStatusTab = ref<'ALL' | 'POSTED' | 'DRAFT'>('ALL')
const activeViewTab = ref<'DOCUMENTS' | 'COA_BALANCES'>('DOCUMENTS')

// Modal preview
const showPreviewModal = ref(false)
const selectedJournal = ref<JournalEntry | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      ledgerStore.fetchAccounts(),
      journalsStore.fetchJournals()
    ])
  } finally {
    loading.value = false
  }
})

// Filtered Opening Balance Journals
const openingJournals = computed(() => {
  return journalsStore.journals.filter(j => 
    j.sourceReference?.startsWith('OP-') || 
    j.journalNumber?.startsWith('OP-') || 
    j.description?.toLowerCase().includes('saldo awal')
  )
})

// Helper for opening balance
function getAccOpening(a: any): number {
  return Number(a.openingBalance ?? a.currentBalance ?? 0)
}

// COA accounts with active opening balances
const accountsWithOpeningBalance = computed(() => {
  return ledgerStore.accounts
    .filter(a => !a.isHeader && getAccOpening(a) !== 0)
    .sort((a, b) => a.code.localeCompare(b.code))
})

// Filtered Documents
const filteredDocuments = computed(() => {
  let docs = openingJournals.value

  if (selectedUnit.value !== 'ALL') {
    docs = docs.filter(j => j.businessUnit === selectedUnit.value || j.businessUnit === 'ALL')
  }

  if (selectedStatusTab.value !== 'ALL') {
    docs = docs.filter(j => j.status === selectedStatusTab.value)
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    docs = docs.filter(j => 
      j.journalNumber.toLowerCase().includes(q) ||
      (j.sourceReference && j.sourceReference.toLowerCase().includes(q)) ||
      j.description.toLowerCase().includes(q) ||
      j.lines.some(l => l.accountName.toLowerCase().includes(q) || l.accountCode.toLowerCase().includes(q))
    )
  }

  return docs
})

// Filtered COA Accounts
const filteredAccounts = computed(() => {
  let accs = accountsWithOpeningBalance.value

  if (selectedUnit.value !== 'ALL') {
    accs = accs.filter(a => a.businessUnit === selectedUnit.value || a.businessUnit === 'ALL')
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    accs = accs.filter(a => 
      a.code.toLowerCase().includes(q) ||
      a.name.toLowerCase().includes(q) ||
      a.type.toLowerCase().includes(q)
    )
  }

  return accs
})

// Summary Calculations
const totalOpeningAsset = computed(() => {
  return accountsWithOpeningBalance.value
    .filter(a => a.type === 'ASSET')
    .reduce((sum, a) => sum + getAccOpening(a), 0)
})

const totalOpeningLiability = computed(() => {
  return accountsWithOpeningBalance.value
    .filter(a => a.type === 'LIABILITY')
    .reduce((sum, a) => sum + getAccOpening(a), 0)
})

const totalOpeningEquity = computed(() => {
  return accountsWithOpeningBalance.value
    .filter(a => a.type === 'EQUITY')
    .reduce((sum, a) => sum + getAccOpening(a), 0)
})

const totalPasiva = computed(() => totalOpeningLiability.value + totalOpeningEquity.value)
const isBalanced = computed(() => Math.abs(totalOpeningAsset.value - totalPasiva.value) < 1)
const difference = computed(() => totalOpeningAsset.value - totalPasiva.value)

function inspectJournal(j: JournalEntry) {
  selectedJournal.value = j
  showPreviewModal.value = true
}

// Export functions
function handleExportExcel() {
  if (activeViewTab.value === 'DOCUMENTS') {
    const exportData = filteredDocuments.value.map((d, i) => ({
      No: i + 1,
      'No. Dokumen / Jurnal': d.journalNumber,
      'No. Referensi': d.sourceReference || '-',
      'Tanggal Cut-Off': formatDate(d.journalDate),
      'Unit Bisnis': d.businessUnit,
      'Keterangan': d.description,
      'Total Debet (Rp)': d.totalDebit,
      'Total Kredit (Rp)': d.totalCredit,
      'Status': d.status
    }))
    exportToExcel(exportData, `Daftar_Dokumen_Saldo_Awal_${accountingStore.businessUnit}`)
    accountingStore.showToast('Daftar dokumen saldo awal berhasil diekspor ke Excel (.xlsx)', 'success')
  } else {
    const exportData = filteredAccounts.value.map((a, i) => ({
      No: i + 1,
      'Kode Akun': a.code,
      'Nama Akun': a.name,
      'Kategori COA': a.type,
      'Unit': a.businessUnit,
      'Posisi Normal': a.normalBalance,
      'Saldo Debet (Rp)': a.normalBalance === 'DEBIT' ? getAccOpening(a) : 0,
      'Saldo Kredit (Rp)': a.normalBalance === 'CREDIT' ? getAccOpening(a) : 0
    }))
    exportToExcel(exportData, `Daftar_Saldo_Akun_COA_${accountingStore.businessUnit}`)
    accountingStore.showToast('Daftar saldo awal akun COA berhasil diekspor ke Excel (.xlsx)', 'success')
  }
}

function handleExportCSV() {
  if (activeViewTab.value === 'DOCUMENTS') {
    const exportData = filteredDocuments.value.map((d, i) => ({
      No: i + 1,
      No_Dokumen: d.journalNumber,
      No_Referensi: d.sourceReference || '-',
      Tanggal_Cutoff: d.journalDate,
      Unit_Bisnis: d.businessUnit,
      Keterangan: d.description,
      Total_Debet: d.totalDebit,
      Total_Kredit: d.totalCredit,
      Status: d.status
    }))
    exportToCSV(exportData, `Daftar_Dokumen_Saldo_Awal_${accountingStore.businessUnit}`)
    accountingStore.showToast('Daftar dokumen saldo awal berhasil diekspor ke CSV (.csv)', 'success')
  } else {
    const exportData = filteredAccounts.value.map((a, i) => ({
      No: i + 1,
      Kode_Akun: a.code,
      Nama_Akun: a.name,
      Kategori: a.type,
      Unit: a.businessUnit,
      Posisi_Normal: a.normalBalance,
      Debet: a.normalBalance === 'DEBIT' ? getAccOpening(a) : 0,
      Kredit: a.normalBalance === 'CREDIT' ? getAccOpening(a) : 0
    }))
    exportToCSV(exportData, `Daftar_Saldo_Akun_COA_${accountingStore.businessUnit}`)
    accountingStore.showToast('Daftar saldo awal akun COA berhasil diekspor ke CSV (.csv)', 'success')
  }
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- 1. Header & Actions -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2 text-xs text-zinc-500 mb-1">
          <router-link to="/accounting/chart-of-accounts" class="hover:text-zinc-900 transition flex items-center gap-1.5">
            <Folder class="w-4 h-4 text-zinc-400" />
            <span>Buku & Transaksi</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-900 font-semibold">Saldo Awal Akuntansi</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight">
          Daftar Saldo Awal Akuntansi
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Pusat kendali dokumen saldo awal pembukuan akuntansi, neraca awal migrasi, dan mutasi saldo berjalan per unit klinik.
        </p>
      </div>

      <div class="flex items-center gap-2.5 flex-wrap">
        <!-- Button: Buat / Atur Saldo Awal Baru -->
        <router-link
          to="/accounting/opening-balance/create"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs active:scale-95 cursor-pointer"
        >
          <Plus class="w-4 h-4 text-emerald-400" />
          <span>Atur / Buat Saldo Awal Baru</span>
        </router-link>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Ekspor ke spreadsheet Excel"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>

        <button
          @click="handleExportCSV"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Ekspor ke file teks CSV"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor CSV</span>
        </button>
      </div>
    </header>

    <!-- 2. KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1: Total Aset / Aktiva -->
      <div class="p-4 bg-white border border-zinc-200/90 rounded-2xl shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-xs text-zinc-500">
          <span class="font-medium">Total Aset (Aktiva Bersih)</span>
          <span class="p-1 rounded-lg bg-emerald-50 text-emerald-700">
            <Scale class="w-3.5 h-3.5" />
          </span>
        </div>
        <div class="text-xl font-bold font-mono text-emerald-700">
          {{ formatIDR(totalOpeningAsset) }}
        </div>
        <p class="text-[11px] text-zinc-400">Kas, Bank, Piutang Pasien & Persediaan Obat</p>
      </div>

      <!-- Card 2: Total Kewajiban & Ekuitas -->
      <div class="p-4 bg-white border border-zinc-200/90 rounded-2xl shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-xs text-zinc-500">
          <span class="font-medium">Total Pasiva (Hutang + Modal)</span>
          <span class="p-1 rounded-lg bg-teal-50 text-teal-700">
            <Layers class="w-3.5 h-3.5" />
          </span>
        </div>
        <div class="text-xl font-bold font-mono text-zinc-900">
          {{ formatIDR(totalPasiva) }}
        </div>
        <p class="text-[11px] text-zinc-400">
          Hutang: {{ formatIDR(totalOpeningLiability) }} • Modal: {{ formatIDR(totalOpeningEquity) }}
        </p>
      </div>

      <!-- Card 3: Status Keseimbangan Neraca -->
      <div
        :class="[
          'p-4 border rounded-2xl shadow-2xs space-y-1 transition',
          isBalanced ? 'bg-emerald-50/60 border-emerald-200' : 'bg-amber-50/60 border-amber-200'
        ]"
      >
        <div class="flex items-center justify-between text-xs">
          <span :class="['font-bold', isBalanced ? 'text-emerald-900' : 'text-amber-900']">Status Keseimbangan Neraca</span>
          <ShieldCheck v-if="isBalanced" class="w-4 h-4 text-emerald-600" />
          <AlertTriangle v-else class="w-4 h-4 text-amber-600" />
        </div>
        <div :class="['text-base font-black', isBalanced ? 'text-emerald-800' : 'text-amber-800']">
          {{ isBalanced ? '100% SEIMBANG (BALANCE)' : 'BELUM SEIMBANG' }}
        </div>
        <p :class="['text-[11px]', isBalanced ? 'text-emerald-700' : 'text-amber-700 font-semibold']">
          {{ isBalanced ? 'Total Debet sama dengan Total Kredit' : `Selisih: ${formatIDR(Math.abs(difference))}` }}
        </p>
      </div>

      <!-- Card 4: Status Cut-off & Akun Terisi -->
      <div class="p-4 bg-white border border-zinc-200/90 rounded-2xl shadow-2xs space-y-1">
        <div class="flex items-center justify-between text-xs text-zinc-500">
          <span class="font-medium">Cut-Off & Rekening Aktif</span>
          <span class="p-1 rounded-lg bg-blue-50 text-blue-700">
            <Calendar class="w-3.5 h-3.5" />
          </span>
        </div>
        <div class="text-xl font-bold text-zinc-900 font-mono">
          01 Jan 2026
        </div>
        <p class="text-[11px] text-zinc-400">
          {{ accountsWithOpeningBalance.length }} Akun COA memiliki saldo awal aktif
        </p>
      </div>
    </div>

    <!-- 3. Navigation Sub-Tabs -->
    <div class="border-b border-zinc-200 flex items-center justify-between">
      <nav class="flex space-x-6 text-xs font-medium">
        <button
          type="button"
          @click="activeViewTab = 'DOCUMENTS'"
          class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
          :class="activeViewTab === 'DOCUMENTS' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
        >
          <FileText class="w-3.5 h-3.5" />
          <span>Daftar Dokumen Saldo Awal ({{ openingJournals.length }})</span>
        </button>

        <button
          type="button"
          @click="activeViewTab = 'COA_BALANCES'"
          class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
          :class="activeViewTab === 'COA_BALANCES' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
        >
          <Layers class="w-3.5 h-3.5" />
          <span>Rincian Saldo Akun COA Aktif ({{ accountsWithOpeningBalance.length }})</span>
        </button>
      </nav>

      <!-- Fast Link to Create Form -->
      <router-link
        to="/accounting/opening-balance/create"
        class="hidden sm:inline-flex items-center gap-1 text-xs font-bold text-blue-600 hover:text-blue-800 pb-2.5 transition"
      >
        <span>Buka Formulir Input Saldo Awal</span>
        <ArrowRight class="w-3.5 h-3.5" />
      </router-link>
    </div>

    <!-- 4. Filter and Search Bar -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-2 flex-1 max-w-md">
        <div class="relative w-full">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari no. dokumen, keterangan, nama akun..."
            class="w-full pl-8 pr-3 py-2 bg-white border border-zinc-200 rounded-xl text-zinc-800 placeholder-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-950 transition"
          />
        </div>
      </div>

      <div class="flex items-center gap-2">
        <!-- Unit Filter -->
        <select
          v-model="selectedUnit"
          class="px-3 py-2 bg-white border border-zinc-200 rounded-xl font-medium text-zinc-800 focus:outline-none focus:ring-1 focus:ring-zinc-950 cursor-pointer"
        >
          <option value="ALL">Semua Unit (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Klinik Spesialis Anak)</option>
          <option value="SSC">Unit SSC (Semarang Skin Center)</option>
        </select>

        <!-- Status Filter (for Documents tab) -->
        <select
          v-if="activeViewTab === 'DOCUMENTS'"
          v-model="selectedStatusTab"
          class="px-3 py-2 bg-white border border-zinc-200 rounded-xl font-medium text-zinc-800 focus:outline-none focus:ring-1 focus:ring-zinc-950 cursor-pointer"
        >
          <option value="ALL">Semua Status</option>
          <option value="POSTED">POSTED (Terposting Buku Besar)</option>
          <option value="DRAFT">DRAFT (Dalam Review)</option>
        </select>
      </div>
    </div>

    <!-- 5. TAB 1: DAFTAR DOKUMEN SALDO AWAL -->
    <div v-if="activeViewTab === 'DOCUMENTS'" class="space-y-4">
      <div class="border border-zinc-200 rounded-2xl overflow-hidden bg-white shadow-2xs">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-500 font-bold border-b border-zinc-200 text-[11px] uppercase tracking-wider">
            <tr>
              <th class="px-4 py-3 w-12 text-center">#</th>
              <th class="px-4 py-3 w-44">No. Dokumen / Jurnal</th>
              <th class="px-4 py-3 w-32">Tanggal Cut-Off</th>
              <th class="px-4 py-3 w-28 text-center">Unit Bisnis</th>
              <th class="px-4 py-3 min-w-[240px]">Keterangan Pembukuan</th>
              <th class="px-4 py-3 w-36 text-right">Total Debet</th>
              <th class="px-4 py-3 w-36 text-right">Total Kredit</th>
              <th class="px-4 py-3 w-28 text-center">Status</th>
              <th class="px-4 py-3 w-36 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-800">
            <!-- Row Items -->
            <tr
              v-for="(doc, idx) in filteredDocuments"
              :key="doc.id"
              class="hover:bg-zinc-50/80 transition"
            >
              <td class="px-4 py-3 text-center text-zinc-400 font-mono">
                {{ idx + 1 }}
              </td>

              <td class="px-4 py-3 font-mono">
                <button
                  type="button"
                  @click="inspectJournal(doc)"
                  class="font-bold text-blue-600 hover:text-blue-800 hover:underline flex items-center gap-1 cursor-pointer"
                  title="Lihat Rincian Voucher"
                >
                  <span>{{ doc.journalNumber }}</span>
                  <ExternalLink class="w-3 h-3 text-zinc-400" />
                </button>
                <span v-if="doc.sourceReference" class="text-[10px] text-zinc-400 block font-mono">
                  Ref: {{ doc.sourceReference }}
                </span>
              </td>

              <td class="px-4 py-3 font-mono text-zinc-600 whitespace-nowrap">
                {{ formatDate(doc.journalDate) }}
              </td>

              <td class="px-4 py-3 text-center whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-zinc-100 text-zinc-700 border border-zinc-200">
                  {{ doc.businessUnit }}
                </span>
              </td>

              <td class="px-4 py-3 font-sans text-zinc-800">
                <div class="font-medium">{{ doc.description }}</div>
                <div class="text-[10px] text-zinc-400 font-mono mt-0.5">
                  {{ (doc.lines || []).length }} Baris Entri Double-Entry
                </div>
              </td>

              <td class="px-4 py-3 text-right font-mono font-bold text-emerald-700 whitespace-nowrap">
                {{ formatIDR(doc.totalDebit) }}
              </td>

              <td class="px-4 py-3 text-right font-mono font-bold text-rose-700 whitespace-nowrap">
                {{ formatIDR(doc.totalCredit) }}
              </td>

              <td class="px-4 py-3 text-center whitespace-nowrap">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                    doc.status === 'POSTED'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : 'bg-amber-50 text-amber-700 border border-amber-200'
                  ]"
                >
                  {{ doc.status }}
                </span>
              </td>

              <td class="px-4 py-3 text-center whitespace-nowrap">
                <div class="flex items-center justify-center gap-1.5">
                  <router-link
                    to="/accounting/opening-balance/create"
                    class="px-2.5 py-1 bg-zinc-100 hover:bg-zinc-200 text-zinc-800 font-semibold rounded-lg text-xs transition cursor-pointer"
                    title="Buka Formulir untuk Menyesuaikan Saldo Awal"
                  >
                    Buka / Edit
                  </router-link>
                  <button
                    type="button"
                    @click="inspectJournal(doc)"
                    class="p-1 text-zinc-500 hover:text-zinc-900 rounded-lg hover:bg-zinc-100 transition cursor-pointer"
                    title="Pratinjau Voucher Jurnal"
                  >
                    <Eye class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>

            <!-- Fallback Master Row if no custom journal entry is found yet -->
            <tr v-if="filteredDocuments.length === 0 && accountsWithOpeningBalance.length > 0" class="hover:bg-zinc-50/80 transition">
              <td class="px-4 py-3 text-center text-zinc-400 font-mono">1</td>
              <td class="px-4 py-3 font-mono">
                <span class="font-bold text-zinc-900">OP-20260101-BCS</span>
                <span class="text-[10px] text-zinc-400 block font-mono">Migrasi Saldo Awal</span>
              </td>
              <td class="px-4 py-3 font-mono text-zinc-600 whitespace-nowrap">01/01/2026</td>
              <td class="px-4 py-3 text-center whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-zinc-100 text-zinc-700 border border-zinc-200">
                  {{ selectedUnit }}
                </span>
              </td>
              <td class="px-4 py-3 font-sans text-zinc-800">
                <div class="font-medium">Saldo Awal Pembukuan Akuntansi & Neraca Migrasi PT. Semarang Medical Group</div>
                <div class="text-[10px] text-zinc-400 font-mono mt-0.5">
                  {{ accountsWithOpeningBalance.length }} Akun COA Aktif
                </div>
              </td>
              <td class="px-4 py-3 text-right font-mono font-bold text-emerald-700 whitespace-nowrap">
                {{ formatIDR(totalOpeningAsset) }}
              </td>
              <td class="px-4 py-3 text-right font-mono font-bold text-rose-700 whitespace-nowrap">
                {{ formatIDR(totalPasiva) }}
              </td>
              <td class="px-4 py-3 text-center whitespace-nowrap">
                <span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase bg-emerald-50 text-emerald-700 border border-emerald-200">
                  POSTED
                </span>
              </td>
              <td class="px-4 py-3 text-center whitespace-nowrap">
                <router-link
                  to="/accounting/opening-balance/create"
                  class="px-2.5 py-1 bg-zinc-950 hover:bg-zinc-800 text-white font-bold rounded-lg text-xs transition cursor-pointer"
                >
                  Buka / Edit
                </router-link>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-if="filteredDocuments.length === 0 && accountsWithOpeningBalance.length === 0">
              <td colspan="9" class="px-4 py-12 text-center text-zinc-400">
                <Layers class="w-8 h-8 mx-auto mb-2 opacity-40 text-zinc-400" />
                <p class="font-semibold text-zinc-700 text-xs">Belum ada dokumen saldo awal akuntansi</p>
                <p class="text-[11px] text-zinc-400 mt-1 max-w-sm mx-auto">
                  Mulai dengan membuat dokumen saldo awal pertama untuk memigrasikan neraca kasir dan buku besar.
                </p>
                <router-link
                  to="/accounting/opening-balance/create"
                  class="mt-3 inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition"
                >
                  <Plus class="w-3.5 h-3.5 text-emerald-400" />
                  <span>Buat Saldo Awal Sekarang</span>
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 6. TAB 2: RINCIAN SALDO AKUN COA AKTIF -->
    <div v-else class="space-y-4">
      <div class="border border-zinc-200 rounded-2xl overflow-hidden bg-white shadow-2xs">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-500 font-bold border-b border-zinc-200 text-[11px] uppercase tracking-wider">
            <tr>
              <th class="px-4 py-3 w-12 text-center">#</th>
              <th class="px-4 py-3 w-32 font-mono">Kode COA</th>
              <th class="px-4 py-3 min-w-[240px]">Nama Akun Buku Besar</th>
              <th class="px-4 py-3 w-28 text-center">Kategori</th>
              <th class="px-4 py-3 w-24 text-center">Unit</th>
              <th class="px-4 py-3 w-28 text-center">Saldo Normal</th>
              <th class="px-4 py-3 w-36 text-right">Saldo Debet (Rp)</th>
              <th class="px-4 py-3 w-36 text-right">Saldo Kredit (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-800 font-mono">
            <tr
              v-for="(acc, idx) in filteredAccounts"
              :key="acc.id || acc.code"
              class="hover:bg-zinc-50 transition"
            >
              <td class="px-4 py-2.5 text-center text-zinc-400">{{ idx + 1 }}</td>
              <td class="px-4 py-2.5 font-bold text-zinc-900">{{ acc.code }}</td>
              <td class="px-4 py-2.5 font-sans font-medium text-zinc-900">{{ acc.name }}</td>
              <td class="px-4 py-2.5 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                    acc.type === 'ASSET' ? 'bg-emerald-50 text-emerald-700' :
                    acc.type === 'LIABILITY' ? 'bg-rose-50 text-rose-700' :
                    'bg-teal-50 text-teal-700'
                  ]"
                >
                  {{ acc.type }}
                </span>
              </td>
              <td class="px-4 py-2.5 text-center text-zinc-600">{{ acc.businessUnit }}</td>
              <td class="px-4 py-2.5 text-center">
                <span class="text-[10px] font-semibold text-zinc-500">{{ acc.normalBalance }}</span>
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-emerald-700">
                {{ acc.normalBalance === 'DEBIT' ? formatIDR(getAccOpening(acc)) : '-' }}
              </td>
              <td class="px-4 py-2.5 text-right font-bold text-rose-700">
                {{ acc.normalBalance === 'CREDIT' ? formatIDR(getAccOpening(acc)) : '-' }}
              </td>
            </tr>

            <tr v-if="filteredAccounts.length === 0">
              <td colspan="8" class="px-4 py-10 text-center text-zinc-400 font-sans">
                Tidak ada rekening COA yang cocok dengan filter.
              </td>
            </tr>
          </tbody>
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs">
            <tr>
              <td colspan="6" class="px-4 py-3 text-right uppercase text-zinc-500 font-sans">
                Total Saldo Awal Terdata:
              </td>
              <td class="px-4 py-3 text-right font-mono text-emerald-800">
                {{ formatIDR(totalOpeningAsset) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-rose-800">
                {{ formatIDR(totalPasiva) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- 7. Detail Preview Modal -->
    <div
      v-if="showPreviewModal && selectedJournal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/50 backdrop-blur-xs animate-in fade-in"
    >
      <div class="bg-white rounded-2xl border border-zinc-200 shadow-xl max-w-3xl w-full p-6 space-y-4 max-h-[90vh] flex flex-col">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div>
            <span class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider block">Voucher Jurnal Saldo Awal</span>
            <h3 class="text-base font-bold text-zinc-950 font-mono">{{ selectedJournal.journalNumber }}</h3>
          </div>
          <button
            type="button"
            @click="showPreviewModal = false"
            class="p-1.5 rounded-lg text-zinc-400 hover:text-zinc-700 hover:bg-zinc-100 transition cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 bg-zinc-50 p-3 rounded-xl border border-zinc-200 text-xs font-mono">
          <div>
            <span class="text-[10px] text-zinc-500 block font-sans">Tanggal:</span>
            <span class="font-bold text-zinc-900">{{ formatDate(selectedJournal.journalDate) }}</span>
          </div>
          <div>
            <span class="text-[10px] text-zinc-500 block font-sans">Unit:</span>
            <span class="font-bold text-zinc-900">{{ selectedJournal.businessUnit }}</span>
          </div>
          <div>
            <span class="text-[10px] text-zinc-500 block font-sans">Status:</span>
            <span class="font-bold text-emerald-700">{{ selectedJournal.status }}</span>
          </div>
          <div>
            <span class="text-[10px] text-zinc-500 block font-sans">No. Referensi:</span>
            <span class="font-bold text-zinc-900">{{ selectedJournal.sourceReference || '-' }}</span>
          </div>
        </div>

        <div class="overflow-y-auto flex-1 border border-zinc-200 rounded-xl">
          <table class="w-full text-xs text-left">
            <thead class="bg-zinc-50 text-zinc-500 font-bold border-b border-zinc-200 text-[10px] uppercase">
              <tr>
                <th class="px-3 py-2">Akun COA</th>
                <th class="px-3 py-2">Uraian</th>
                <th class="px-3 py-2 text-right">Debet</th>
                <th class="px-3 py-2 text-right">Kredit</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100 font-mono text-zinc-800">
              <tr v-for="(line, lidx) in selectedJournal.lines" :key="lidx">
                <td class="px-3 py-2 font-bold">{{ line.accountCode }} - {{ line.accountName }}</td>
                <td class="px-3 py-2 font-sans text-zinc-600">{{ line.description }}</td>
                <td class="px-3 py-2 text-right text-emerald-700 font-bold">{{ line.debit > 0 ? formatIDR(line.debit) : '-' }}</td>
                <td class="px-3 py-2 text-right text-rose-700 font-bold">{{ line.credit > 0 ? formatIDR(line.credit) : '-' }}</td>
              </tr>
            </tbody>
            <tfoot class="bg-zinc-50 font-mono font-bold text-xs border-t border-zinc-200">
              <tr>
                <td colspan="2" class="px-3 py-2 text-right font-sans uppercase text-zinc-500">Total Voucher:</td>
                <td class="px-3 py-2 text-right text-emerald-700">{{ formatIDR(selectedJournal.totalDebit) }}</td>
                <td class="px-3 py-2 text-right text-rose-700">{{ formatIDR(selectedJournal.totalCredit) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>

        <div class="flex items-center justify-between pt-3 border-t border-zinc-100 text-xs">
          <router-link
            to="/accounting/opening-balance/create"
            class="px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white font-bold rounded-xl transition"
          >
            Buka di Formulir Edit
          </router-link>
          <button
            type="button"
            @click="showPreviewModal = false"
            class="px-3.5 py-1.5 border border-zinc-200 rounded-xl text-zinc-700 font-semibold hover:bg-zinc-50 transition"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
