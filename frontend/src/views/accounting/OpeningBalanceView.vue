<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { BusinessUnit, AccountType, NormalBalance } from '@/types/accounting'
import {
  BookmarkCheck,
  Scale,
  Calendar,
  Building2,
  FileText,
  Search,
  RotateCcw,
  Sparkles,
  Send,
  AlertTriangle,
  CheckCircle2,
  Printer,
  ArrowRight,
  Database,
  Layers,
  HelpCircle,
  FolderOpen,
  Plus,
  Trash2,
  TableProperties,
  ListPlus
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  initializeAllCoaRows()
})

// Header Form
const cutOffDate = ref('2026-01-01')
const referenceNo = ref('OP-20260101-BCS')
const selectedUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const description = ref('Saldo Awal Pembukuan Akuntansi & Neraca Migrasi PT. Brayo Cipta Sejahtera (Kidinest & SSC) per 1 Januari 2026')

// Mode Switch: Full COA Matrix vs Quick Custom Rows
const activeMode = ref<'COA_MATRIX' | 'CUSTOM_ROWS'>('COA_MATRIX')
const matrixSearch = ref('')
const selectedTypeFilter = ref<string>('ALL')

// Input Row Structure
interface BalanceRow {
  code: string
  name: string
  type: AccountType
  normalBalance: NormalBalance
  unit: BusinessUnit
  debit: number | string
  credit: number | string
  notes: string
}

// Full COA Matrix Rows (all leaf accounts from 106 BCS COA)
const matrixRows = ref<BalanceRow[]>([])
// Custom Rows (dynamic manual rows)
const customRows = ref<BalanceRow[]>([])

function initializeAllCoaRows() {
  const leafAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)

  matrixRows.value = leafAccounts.map(acc => {
    const cur = acc.currentBalance || 0
    return {
      code: acc.code,
      name: acc.name,
      type: acc.type,
      normalBalance: acc.normalBalance,
      unit: acc.businessUnit || 'ALL',
      debit: acc.normalBalance === 'DEBIT' ? (cur > 0 ? cur : 0) : 0,
      credit: acc.normalBalance === 'CREDIT' ? (cur > 0 ? cur : 0) : 0,
      notes: cur > 0 ? 'Saldo migrasi berjalan' : ''
    }
  })

  // Also setup default custom rows
  if (customRows.value.length === 0) {
    customRows.value = [
      { code: '1101.001', name: 'Kas Kasir Kidinest', type: 'ASSET', normalBalance: 'DEBIT', unit: 'KIDINEST', debit: 45200000, credit: 0, notes: 'Kas fisik kasir awal' },
      { code: '1102.001', name: 'Bank BCA Operasional (Giro & EDC)', type: 'ASSET', normalBalance: 'DEBIT', unit: 'ALL', debit: 680964487.55, credit: 0, notes: 'Rekening giro operasional' },
      { code: '3101', name: 'Modal Disetor Pemilik Klinik Kidinest & SSC', type: 'EQUITY', normalBalance: 'CREDIT', unit: 'ALL', debit: 0, credit: 726164487.55, notes: 'Penyertaan modal awal' }
    ]
  }
}

// Active Rows based on selected mode
const activeWorkingRows = computed(() => {
  return activeMode.value === 'COA_MATRIX' ? matrixRows.value : customRows.value
})

// Filtered Matrix Rows for search & category
const filteredMatrixRows = computed(() => {
  return matrixRows.value.filter(r => {
    if (selectedTypeFilter.value !== 'ALL' && r.type !== selectedTypeFilter.value) return false
    if (!matrixSearch.value) return true
    const q = matrixSearch.value.toLowerCase()
    return r.code.toLowerCase().includes(q) || r.name.toLowerCase().includes(q)
  })
})

// Calculations based on active mode
const totalDebit = computed(() => {
  return activeWorkingRows.value.reduce((sum, r) => sum + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return activeWorkingRows.value.reduce((sum, r) => sum + (Number(r.credit) || 0), 0)
})

const difference = computed(() => {
  return Math.round((totalDebit.value - totalCredit.value) * 100) / 100
})

const isBalanced = computed(() => {
  return Math.abs(difference.value) < 0.01 && totalDebit.value > 0
})

// Asset, Liability, Equity Breakdown
const totalAsset = computed(() => {
  return activeWorkingRows.value
    .filter(r => r.type === 'ASSET')
    .reduce((sum, r) => sum + (Number(r.debit) || 0) - (Number(r.credit) || 0), 0)
})

const totalLiability = computed(() => {
  return activeWorkingRows.value
    .filter(r => r.type === 'LIABILITY')
    .reduce((sum, r) => sum + (Number(r.credit) || 0) - (Number(r.debit) || 0), 0)
})

const totalEquity = computed(() => {
  return activeWorkingRows.value
    .filter(r => r.type === 'EQUITY')
    .reduce((sum, r) => sum + (Number(r.credit) || 0) - (Number(r.debit) || 0), 0)
})

// Auto-Balance to Equity (3999 Opening Balance Equity)
function autoBalanceToEquity() {
  const diff = difference.value
  if (Math.abs(diff) < 0.01) {
    accountingStore.showToast('Saldo Awal sudah seimbang!', 'info')
    return
  }

  const targetList = activeMode.value === 'COA_MATRIX' ? matrixRows.value : customRows.value

  // Look for 3999 or 3-9999
  let eqRow = targetList.find(r => r.code === '3999' || r.code === '3-9999' || r.code.includes('9999'))

  if (!eqRow) {
    const acc = ledgerStore.accounts.find(a => a.code === '3999' || a.code === '3-9999') ||
                ledgerStore.accounts.find(a => a.type === 'EQUITY' && !a.isHeader)
    if (!acc) {
      accountingStore.showToast('Akun penyeimbang modal (3999 Opening Balance Equity) tidak ditemukan!', 'error')
      return
    }
    eqRow = {
      code: acc.code,
      name: acc.name,
      type: acc.type,
      normalBalance: acc.normalBalance,
      unit: 'ALL',
      debit: 0,
      credit: 0,
      notes: 'Penyeimbang Saldo Awal'
    }
    targetList.push(eqRow)
  }

  if (diff > 0) {
    // Debit is greater than Credit -> Add diff to Credit of Equity
    eqRow.credit = (Number(eqRow.credit) || 0) + diff
    eqRow.notes = 'Penyeimbang selisih debit saldo awal'
  } else {
    // Credit is greater than Debit -> Add diff to Debit of Equity
    eqRow.debit = (Number(eqRow.debit) || 0) + Math.abs(diff)
    eqRow.notes = 'Penyeimbang selisih kredit saldo awal'
  }

  accountingStore.showToast(`Selisih Rp ${formatIDR(Math.abs(diff))} dialokasikan ke akun ${eqRow.code} - ${eqRow.name}!`, 'success')
}

// Load Standard Template
function loadBcsTemplate() {
  // Preset standard realistic clinic balances
  const templateMap: Record<string, { debit?: number; credit?: number; notes: string }> = {
    '1101.001': { debit: 45200000, notes: 'Kas Kasir Kidinest' },
    '1101.002': { debit: 38400000, notes: 'Kas Kasir SSC' },
    '1101.003': { debit: 5000000, notes: 'Kas Kecil Petty Cash' },
    '1102.001': { debit: 680964487.55, notes: 'Bank BCA Giro & EDC' },
    '1102.002': { debit: 493864533.33, notes: 'Bank Mandiri Operasional' },
    '1102.003': { debit: 14500000, notes: 'QRIS Transit Settlement' },
    '1201': { debit: 485600000, notes: 'Piutang Pasien Perawatan' },
    '1202.001': { debit: 420000000, notes: 'Piutang Klaim Sinarmas' },
    '1202.002': { debit: 465000000, notes: 'Piutang Klaim Prudential' },
    '1301': { debit: 4885000000, notes: 'Persediaan Obat Resep Kidinest' },
    '1302': { debit: 6442748296.88, notes: 'Persediaan Produk Skincare SSC' },
    '1303': { debit: 3208000000, notes: 'Persediaan Bahan Medis (BMHP)' },
    '1304': { debit: 4800000000, notes: 'Persediaan Vaksin Anak' },
    '1601': { debit: 350000000, notes: 'Renovasi Interior Klinik' },
    '1602': { debit: 500000000, notes: 'Mesin Laser Nd:YAG SSC' },
    '1603': { debit: 250000000, notes: 'Peralatan USG 4D Poli Anak' },
    '1606': { credit: 155000000, notes: 'Akumulasi Penyusutan Alkes' },
    '2101.001': { credit: 35500000, notes: 'Uang Muka Reservasi Poli Anak' },
    '2102.001': { credit: 1904190360, notes: 'Deposit Paket Skincare Pasien' },
    '2201': { credit: 2800000, notes: 'Hutang PBF Kimia Farma' },
    '2202': { credit: 1500000, notes: 'Hutang Supplier Kosmetika SSC' },
    '2301': { credit: 85000000, notes: 'Hutang Bagi Hasil Dokter' },
    '2501.001': { credit: 22374594, notes: 'Hutang PPh 21 Dokter' },
    '2502': { credit: 2400000, notes: 'Hutang PPh 23 Limbah Medis B3' },
    '2503': { credit: 5000000, notes: 'Hutang PPh Final 4(2) Sewa' },
    '2504': { credit: 9475055, notes: 'Hutang PPN Keluaran' },
    '3101': { credit: 2367382942.15, notes: 'Modal Disetor Pemilik Klinik' },
    '3201': { credit: 14132108388.64, notes: 'Laba Ditahan Periode Lalu' }
  }

  matrixRows.value.forEach(row => {
    const t = templateMap[row.code]
    if (t) {
      row.debit = t.debit || 0
      row.credit = t.credit || 0
      row.notes = t.notes
    } else {
      row.debit = 0
      row.credit = 0
    }
  })

  autoBalanceToEquity()
  accountingStore.showToast('Template Saldo Awal Standar BCS berhasil dimuat!', 'success')
}

// Reset Balances
function resetBalances() {
  if (confirm('Kosongkan seluruh nilai saldo awal pada tabel aktif?')) {
    if (activeMode.value === 'COA_MATRIX') {
      matrixRows.value.forEach(r => {
        r.debit = 0
        r.credit = 0
        r.notes = ''
      })
    } else {
      customRows.value = [
        { code: '', name: '', type: 'ASSET', normalBalance: 'DEBIT', unit: 'ALL', debit: 0, credit: 0, notes: '' }
      ]
    }
    accountingStore.showToast('Seluruh input saldo awal berhasil dikosongkan.', 'info')
  }
}

// Custom Row Helpers
function addCustomRow() {
  customRows.value.push({
    code: '',
    name: '',
    type: 'ASSET',
    normalBalance: 'DEBIT',
    unit: selectedUnit.value || 'ALL',
    debit: 0,
    credit: 0,
    notes: ''
  })
}

function removeCustomRow(index: number) {
  if (customRows.value.length > 1) {
    customRows.value.splice(index, 1)
  }
}

function onCustomAccountSelect(index: number, code: string) {
  const acc = ledgerStore.accounts.find(a => a.code === code)
  if (acc && customRows.value[index]) {
    customRows.value[index].code = acc.code
    customRows.value[index].name = acc.name
    customRows.value[index].type = acc.type
    customRows.value[index].normalBalance = acc.normalBalance
    customRows.value[index].unit = acc.businessUnit || selectedUnit.value || 'ALL'
  }
}

// Submit & Post Opening Balance
const isSubmitting = ref(false)

async function handleSaveOpeningBalance() {
  if (!isBalanced.value) {
    accountingStore.showToast(`Saldo Awal belum seimbang! Selisih Rp ${formatIDR(Math.abs(difference.value))}. Klik 'Seimbangkan ke Ekuitas' terlebih dahulu.`, 'error')
    return
  }

  const validLines = activeWorkingRows.value
    .filter(r => (Number(r.debit) > 0 || Number(r.credit) > 0))
    .map(r => ({
      accountCode: r.code,
      accountName: r.name,
      debit: Number(r.debit) || 0,
      credit: Number(r.credit) || 0,
      description: r.notes || `Saldo Awal ${r.code} - ${r.name}`
    }))

  if (validLines.length === 0) {
    accountingStore.showToast('Harap masukkan minimal satu baris akun dengan nominal saldo!', 'error')
    return
  }

  try {
    isSubmitting.value = true

    await ledgerStore.saveOpeningBalance({
      cutOffDate: cutOffDate.value,
      referenceNo: referenceNo.value || `OP-${cutOffDate.value.replace(/-/g, '')}`,
      businessUnit: selectedUnit.value,
      description: description.value,
      lines: validLines
    })

    accountingStore.showToast(
      `Saldo Awal Akuntansi per ${cutOffDate.value} berhasil diposting ke database PostgreSQL!`,
      'success'
    )

    router.push('/accounting/chart-of-accounts')
  } catch (err: any) {
    accountingStore.showToast(`Gagal menyimpan saldo awal: ${err.message}`, 'error')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-24">
    <!-- Header with Breadcrumb -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="flex items-center gap-3">
        <button
          @click="router.push('/accounting/chart-of-accounts')"
          class="p-2 rounded-xl bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Kembali ke Bagan Akun"
        >
          <ArrowLeft class="w-4 h-4" />
        </button>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
              Setup Buku Besar & Migrasi
            </span>
            <span class="text-zinc-300">•</span>
            <span class="text-xs text-zinc-400 font-medium">Bagan Akun BCS 106 Akun</span>
          </div>
          <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
            Saldo Awal Akuntansi (Opening Balance Engine)
          </h1>
        </div>
      </div>

      <!-- Action Presets -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="loadBcsTemplate"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold transition border border-zinc-200 shadow-2xs cursor-pointer active:scale-95"
          title="Muat nilai saldo awal standar BCS"
        >
          <FolderOpen class="w-3.5 h-3.5 text-zinc-600" />
          <span>Muat Template Standar BCS</span>
        </button>

        <button
          @click="autoBalanceToEquity"
          :disabled="isBalanced || difference === 0"
          :class="[
            'inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95',
            !isBalanced && difference !== 0
              ? 'bg-amber-500 hover:bg-amber-600 text-white cursor-pointer'
              : 'bg-zinc-100 text-zinc-400 border border-zinc-200 cursor-not-allowed'
          ]"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>Seimbangkan ke Ekuitas (3999)</span>
        </button>

        <button
          @click="resetBalances"
          class="inline-flex items-center gap-1.5 px-3 py-2 bg-white hover:bg-zinc-100 text-zinc-700 rounded-xl text-xs font-bold transition border border-zinc-200 shadow-2xs cursor-pointer"
        >
          <RotateCcw class="w-3.5 h-3.5 text-zinc-500" />
          <span>Reset Input</span>
        </button>
      </div>
    </header>

    <!-- Header Metadata Configuration -->
    <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Tanggal Cut-Off Saldo Awal:</span>
        </label>
        <input
          v-model="cutOffDate"
          type="date"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <FileText class="w-3.5 h-3.5 text-zinc-500" />
          <span>Nomor Dokumen Bukti:</span>
        </label>
        <input
          v-model="referenceNo"
          type="text"
          placeholder="OP-20260101-BCS"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-zinc-500" />
          <span>Unit Bisnis:</span>
        </label>
        <select
          v-model="selectedUnit"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Semua Unit (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Poli Anak)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1">Keterangan Memo Saldo Awal:</label>
        <input
          v-model="description"
          type="text"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>
    </div>

    <!-- Accounting Equation Summary Bar (Aktiva = Pasiva) -->
    <div class="grid grid-cols-1 sm:grid-cols-4 gap-3 text-xs">
      <!-- Total Aktiva (Aset) -->
      <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs">
        <span class="text-[10px] font-bold text-zinc-500 uppercase block">1. Total Aktiva (Aset)</span>
        <span class="text-base font-extrabold text-emerald-700 font-mono block mt-1">
          Rp {{ formatIDR(totalAsset) }}
        </span>
      </div>

      <!-- Total Kewajiban (Liability) -->
      <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs">
        <span class="text-[10px] font-bold text-zinc-500 uppercase block">2. Total Kewajiban (Hutang)</span>
        <span class="text-base font-extrabold text-rose-700 font-mono block mt-1">
          Rp {{ formatIDR(totalLiability) }}
        </span>
      </div>

      <!-- Total Ekuitas (Modal) -->
      <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs">
        <span class="text-[10px] font-bold text-zinc-500 uppercase block">3. Total Ekuitas (Modal)</span>
        <span class="text-base font-extrabold text-blue-700 font-mono block mt-1">
          Rp {{ formatIDR(totalEquity) }}
        </span>
      </div>

      <!-- Balance Status Card -->
      <div
        :class="[
          'p-4 rounded-2xl border flex flex-col justify-between shadow-2xs',
          isBalanced
            ? 'bg-emerald-50 border-emerald-200 text-emerald-950'
            : 'bg-amber-50 border-amber-200 text-amber-950'
        ]"
      >
        <div class="flex items-center justify-between">
          <span class="text-[10px] font-black uppercase tracking-wider">
            {{ isBalanced ? '✓ NERACA SEIMBANG' : '⚠ BELUM SEIMBANG' }}
          </span>
          <span
            :class="[
              'w-2 h-2 rounded-full',
              isBalanced ? 'bg-emerald-500' : 'bg-amber-500 animate-pulse'
            ]"
          ></span>
        </div>
        <div class="mt-1 font-mono">
          <span class="text-[10px] opacity-75 block">Selisih Debit - Kredit:</span>
          <span :class="['text-sm font-extrabold', isBalanced ? 'text-emerald-800' : 'text-rose-700']">
            Rp {{ formatIDR(Math.abs(difference)) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Mode Selector Tabs: COA Matrix vs Custom Rows -->
    <div class="flex items-center justify-between gap-4 border-b border-zinc-200 pb-2">
      <div class="flex items-center gap-2">
        <button
          @click="activeMode = 'COA_MATRIX'"
          :class="[
            'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition cursor-pointer',
            activeMode === 'COA_MATRIX'
              ? 'bg-zinc-950 text-white shadow-2xs'
              : 'bg-white text-zinc-600 border border-zinc-200 hover:bg-zinc-50'
          ]"
        >
          <TableProperties class="w-4 h-4" />
          <span>Matriks Seluruh Bagan Akun COA ({{ matrixRows.length }} Akun)</span>
        </button>

        <button
          @click="activeMode = 'CUSTOM_ROWS'"
          :class="[
            'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition cursor-pointer',
            activeMode === 'CUSTOM_ROWS'
              ? 'bg-zinc-950 text-white shadow-2xs'
              : 'bg-white text-zinc-600 border border-zinc-200 hover:bg-zinc-50'
          ]"
        >
          <ListPlus class="w-4 h-4" />
          <span>Entri Baris Kustom Cepat ({{ customRows.length }} Baris)</span>
        </button>
      </div>

      <!-- Quick Search Bar (For Matrix Mode) -->
      <div v-if="activeMode === 'COA_MATRIX'" class="flex items-center gap-2">
        <select
          v-model="selectedTypeFilter"
          class="px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs font-bold text-zinc-800 focus:outline-none"
        >
          <option value="ALL">Semua Kategori</option>
          <option value="ASSET">Aset (1)</option>
          <option value="LIABILITY">Kewajiban (2)</option>
          <option value="EQUITY">Ekuitas (3)</option>
        </select>
        <div class="relative w-48 sm:w-64">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
          <input
            v-model="matrixSearch"
            type="text"
            placeholder="Cari kode / nama akun..."
            class="w-full pl-8 pr-3 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs text-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-950"
          />
        </div>
      </div>
    </div>

    <!-- MODE 1: MATRIKS SELURUH BAGAN AKUN (COA MATRIX) -->
    <div v-if="activeMode === 'COA_MATRIX'" class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto max-h-[600px] overflow-y-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200 sticky top-0 z-10 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-10 text-center">#</th>
              <th class="px-4 py-3 min-w-[280px]">Bagan Akun (Kode & Nama)</th>
              <th class="px-3 py-3 w-28 text-center">Kategori</th>
              <th class="px-3 py-3 w-24 text-center">Saldo Normal</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Awal Debet (Rp)</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Awal Kredit (Rp)</th>
              <th class="px-4 py-3 min-w-[180px]">Keterangan / Catatan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr
              v-for="(row, idx) in filteredMatrixRows"
              :key="row.code"
              :class="[
                'hover:bg-zinc-50/70 transition',
                (Number(row.debit) > 0 || Number(row.credit) > 0) ? 'bg-amber-50/20 font-medium' : ''
              ]"
            >
              <td class="px-3 py-2 text-center text-zinc-400 font-mono text-[11px]">
                {{ idx + 1 }}
              </td>

              <td class="px-4 py-2">
                <span class="font-mono font-bold text-zinc-900 mr-1.5">{{ row.code }}</span>
                <span>{{ row.name }}</span>
              </td>

              <td class="px-3 py-2 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                    row.type === 'ASSET' ? 'bg-emerald-50 text-emerald-800' :
                    row.type === 'LIABILITY' ? 'bg-rose-50 text-rose-800' :
                    'bg-blue-50 text-blue-800'
                  ]"
                >
                  {{ row.type }}
                </span>
              </td>

              <td class="px-3 py-2 text-center font-mono text-[11px] font-bold text-zinc-500">
                {{ row.normalBalance }}
              </td>

              <!-- Debit Input -->
              <td class="px-4 py-1.5 text-right">
                <input
                  v-model.number="row.debit"
                  type="number"
                  min="0"
                  step="any"
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                  placeholder="0"
                />
              </td>

              <!-- Credit Input -->
              <td class="px-4 py-1.5 text-right">
                <input
                  v-model.number="row.credit"
                  type="number"
                  min="0"
                  step="any"
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                  placeholder="0"
                />
              </td>

              <!-- Notes -->
              <td class="px-4 py-1.5">
                <input
                  v-model="row.notes"
                  type="text"
                  placeholder="Catatan saldo awal..."
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODE 2: ENTRI BARIS KUSTOM CEPAT (CUSTOM ROWS) -->
    <div v-else class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="p-4 border-b border-zinc-200/90 flex items-center justify-between bg-zinc-50/50">
        <span class="text-xs font-bold text-zinc-900">
          Entri Akun Bersaldo Spesifik ({{ customRows.length }} Baris)
        </span>
        <button
          type="button"
          @click="addCustomRow"
          class="inline-flex items-center gap-1 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tambah Baris</span>
        </button>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-10 text-center">#</th>
              <th class="px-4 py-3 min-w-[280px]">Pilih Akun COA</th>
              <th class="px-3 py-3 w-28 text-center">Unit</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Debet (Rp)</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Kredit (Rp)</th>
              <th class="px-4 py-3 min-w-[180px]">Catatan Baris</th>
              <th class="px-3 py-3 w-12 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr v-for="(row, idx) in customRows" :key="idx" class="hover:bg-zinc-50/70 transition">
              <td class="px-3 py-3 text-center text-zinc-400 font-mono">{{ idx + 1 }}</td>
              <td class="px-4 py-2">
                <SearchableSelect
                  :model-value="row.code"
                  :options="ledgerStore.accounts.filter(a => !a.isHeader).map(a => ({ value: a.code, label: `[${a.code}] ${a.name}`, sublabel: a.type }))"
                  placeholder="-- Pilih Akun --"
                  @update:model-value="(code: any) => onCustomAccountSelect(idx, String(code))"
                />
              </td>
              <td class="px-3 py-2 text-center">
                <select v-model="row.unit" class="px-2 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs font-bold text-zinc-800">
                  <option value="ALL">ALL</option>
                  <option value="KIDINEST">Kidinest</option>
                  <option value="SSC">SSC</option>
                </select>
              </td>
              <td class="px-4 py-2 text-right">
                <input v-model.number="row.debit" type="number" min="0" step="any" class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:outline-none" placeholder="0" />
              </td>
              <td class="px-4 py-2 text-right">
                <input v-model.number="row.credit" type="number" min="0" step="any" class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:outline-none" placeholder="0" />
              </td>
              <td class="px-4 py-2">
                <input v-model="row.notes" type="text" placeholder="Catatan..." class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-800 focus:outline-none" />
              </td>
              <td class="px-3 py-2 text-center">
                <button type="button" @click="removeCustomRow(idx)" class="p-1.5 text-zinc-400 hover:text-rose-600 rounded-lg transition">
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Bottom Submission Bar -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2">
      <div class="text-xs text-zinc-500 flex items-center gap-2">
        <span class="font-bold text-zinc-700">Total Akun Terisi:</span>
        <span class="font-mono font-bold text-zinc-950">{{ activeWorkingRows.filter(r => (Number(r.debit) > 0 || Number(r.credit) > 0)).length }} Akun</span>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          :disabled="isSubmitting || !isBalanced"
          @click="handleSaveOpeningBalance"
          class="inline-flex items-center gap-1.5 px-6 py-2.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white rounded-xl text-xs font-bold transition shadow-md cursor-pointer active:scale-95"
        >
          <Send class="w-4 h-4 text-emerald-400" />
          <span>Simpan & Posting Saldo Awal ke PostgreSQL</span>
        </button>
      </div>
    </div>
  </div>
</template>
