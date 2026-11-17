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
  Trash2
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()

// Opening Balance Header Form
const cutOffDate = ref('2026-01-01')
const referenceNo = ref('OP-20260101-CLINIC')
const selectedUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const description = ref('Saldo Awal Pembukuan Akuntansi & Neraca Migrasi Klinik Kidinest & SSC per 1 Januari 2026')

watch(() => accountingStore.businessUnit, (unit) => {
  if (unit) selectedUnit.value = unit
})

// Available active leaf accounts for SearchableSelect
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

// Input Rows State: Dynamic Row Entry (Opsi A)
interface OpeningBalanceRowState {
  code: string
  name: string
  type?: AccountType
  normalBalance?: NormalBalance
  unit?: BusinessUnit
  debit: number | string
  credit: number | string
  notes: string
}

const balanceRows = ref<OpeningBalanceRowState[]>([])

// Populate rows: only accounts with balances, or clean start rows
function initializeRows() {
  const existing = ledgerStore.accounts
    .filter(a => !a.isHeader && a.isActive !== false && (a.currentBalance && a.currentBalance !== 0))
    .map(a => {
      const cur = a.currentBalance || 0
      return {
        code: a.code,
        name: a.name,
        type: a.type,
        normalBalance: a.normalBalance,
        unit: a.businessUnit || 'ALL',
        debit: a.normalBalance === 'DEBIT' ? (cur > 0 ? cur : 0) : (cur < 0 ? Math.abs(cur) : 0),
        credit: a.normalBalance === 'CREDIT' ? (cur > 0 ? cur : 0) : (cur < 0 ? Math.abs(cur) : 0),
        notes: 'Saldo pembukuan migrasi'
      }
    })

  if (existing.length > 0) {
    balanceRows.value = existing
  } else {
    // Default 2 empty rows ready to be filled
    balanceRows.value = [
      { code: '', name: '', debit: 0, credit: 0, notes: '', unit: selectedUnit.value || 'ALL' },
      { code: '', name: '', debit: 0, credit: 0, notes: '', unit: selectedUnit.value || 'ALL' }
    ]
  }
}

function onAccountSelect(index: number, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc && balanceRows.value[index]) {
    balanceRows.value[index].code = acc.code
    balanceRows.value[index].name = acc.name
    balanceRows.value[index].type = acc.type
    balanceRows.value[index].normalBalance = acc.normalBalance
    balanceRows.value[index].unit = acc.businessUnit || selectedUnit.value || 'ALL'
  }
}

function addRow() {
  balanceRows.value.push({
    code: '',
    name: '',
    debit: 0,
    credit: 0,
    notes: '',
    unit: selectedUnit.value || 'ALL'
  })
}

function removeRow(index: number) {
  if (balanceRows.value.length > 1) {
    balanceRows.value.splice(index, 1)
  } else {
    balanceRows.value[0] = { code: '', name: '', debit: 0, credit: 0, notes: '', unit: selectedUnit.value || 'ALL' }
  }
}

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  initializeRows()
})

// Real-time Calculations
const totalDebit = computed(() => {
  return balanceRows.value.reduce((sum, r) => sum + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return balanceRows.value.reduce((sum, r) => sum + (Number(r.credit) || 0), 0)
})

const difference = computed(() => {
  return totalDebit.value - totalCredit.value
})

const isBalanced = computed(() => {
  return Math.abs(difference.value) < 0.01 && totalDebit.value > 0
})

// Action: Muat Angka Standar Saldo Awal Klinik (Simulasi Demo)
function loadClinicStandardBalances() {
  const sampleMap: Record<string, { debit?: number; credit?: number; notes?: string }> = {
    '1-1110': { debit: 5000000, notes: 'Kas Kecil Operasional' },
    '1-1111': { debit: 12500000, notes: 'Kas Kasir Kidinest' },
    '1-1112': { debit: 18500000, notes: 'Kas Kasir SSC' },
    '1-1113': { debit: 145000000, notes: 'Bank BCA Operasional' },
    '1-1114': { debit: 85000000, notes: 'Bank Mandiri Payroll' },
    '1-1310': { debit: 42000000, notes: 'Stok Obat & Vaksin Anak Kidinest' },
    '1-1330': { debit: 58000000, notes: 'Stok Produk Skincare SSC' },
    '1-2100': { debit: 250000000, notes: 'Alkes USG 4D Kidinest' },
    '1-2200': { debit: 380000000, notes: 'Mesin Laser Estetika SSC' },
    '2-1110': { credit: 24500000, notes: 'Hutang Distributor Obat PBF' },
    '2-1115': { credit: 31000000, notes: 'Hutang Supplier Skincare SSC' },
    '3-1100': { credit: 940500000, notes: 'Modal Disetor Pemilik Klinik' }
  }

  balanceRows.value = Object.entries(sampleMap).map(([code, data]) => {
    const acc = ledgerStore.accounts.find(a => a.code === code)
    return {
      code,
      name: acc?.name || code,
      type: acc?.type,
      normalBalance: acc?.normalBalance,
      unit: acc?.businessUnit || 'ALL',
      debit: data.debit || 0,
      credit: data.credit || 0,
      notes: data.notes || 'Migrasi Standar Klinik'
    }
  })
  accountingStore.showToast('Template Saldo Awal Standar Klinik berhasil dimuat!', 'info')
}

// Action: Seimbangkan Otomatis ke Akun Ekuitas Saldo Awal (3999)
function autoBalanceToEquity() {
  const diff = difference.value
  if (Math.abs(diff) < 0.01) {
    accountingStore.showToast('Saldo awal sudah seimbang!', 'info')
    return
  }

  // Find 3-9999 / 3999 - Opening Balance Equity
  let eqRow = balanceRows.value.find(r => r.code === '3-9999' || r.code === '3999' || r.code.includes('9999'))
  if (!eqRow) {
    const acc = ledgerStore.accounts.find(a => a.code === '3-9999' || a.code === '3999') ||
                ledgerStore.accounts.find(a => a.type === 'EQUITY' && !a.isHeader)
    if (!acc) {
      accountingStore.showToast('Akun penyeimbang saldo awal (3-9999) tidak ditemukan', 'error')
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
    balanceRows.value.push(eqRow)
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

  accountingStore.showToast(`Selisih Rp ${formatIDR(Math.abs(diff))} dialokasikan ke ${eqRow.code} - ${eqRow.name}!`, 'success')
}

// Action: Reset Saldo Awal
function resetAllBalances() {
  if (confirm('Kosongkan seluruh baris input saldo awal?')) {
    balanceRows.value = [
      { code: '', name: '', debit: 0, credit: 0, notes: '', unit: selectedUnit.value || 'ALL' },
      { code: '', name: '', debit: 0, credit: 0, notes: '', unit: selectedUnit.value || 'ALL' }
    ]
    accountingStore.showToast('Seluruh input saldo awal dikosongkan', 'info')
  }
}

// Action: Simpan & Posting Saldo Awal ke Database
const isSubmitting = ref(false)
async function handleSaveOpeningBalance() {
  if (!isBalanced.value) {
    accountingStore.showToast(`Saldo Awal belum seimbang! Terdapat selisih Rp ${formatIDR(Math.abs(difference.value))}. Klik 'Seimbangkan Otomatis' atau koreksi input.`, 'error')
    return
  }

  try {
    isSubmitting.value = true

    const lines = balanceRows.value
      .filter(r => (Number(r.debit) > 0 || Number(r.credit) > 0))
      .map(r => ({
        accountCode: r.code,
        accountName: r.name,
        debit: Number(r.debit) || 0,
        credit: Number(r.credit) || 0,
        description: r.notes || `Saldo Awal ${r.code} - ${r.name}`
      }))

    await ledgerStore.saveOpeningBalance({
      cutOffDate: cutOffDate.value,
      referenceNo: referenceNo.value,
      businessUnit: selectedUnit.value,
      description: description.value,
      lines
    })

    // Navigate to General Ledger or Journals to see result
    setTimeout(() => {
      router.push('/accounting/ledger')
    }, 1500)
  } catch (err: any) {
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

// Print / Cetak Saldo Awal
function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif] pb-24">
    <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Setup Buku Besar & Migrasi
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Saldo Awal Akuntansi (Cut-Off)
        </h1>
      </div>

      <!-- Quick Actions -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="loadClinicStandardBalances"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold transition-all border border-zinc-200 shadow-2xs cursor-pointer active:scale-95"
          title="Muat nilai saldo awal standar Klinik"
        >
          <FolderOpen class="w-3.5 h-3.5 text-zinc-600" />
          <span>Muat Data Standar Klinik</span>
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
          @click="handlePrint"
          class="inline-flex items-center gap-1 px-3 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold transition border border-zinc-200 shadow-2xs cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5 text-zinc-600" />
          <span>Cetak</span>
        </button>
      </div>
    </header>

    <!-- Configuration Form Bar -->
    <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-600" />
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
          <FileText class="w-3.5 h-3.5 text-zinc-600" />
          <span>Nomor Dokumen Bukti:</span>
        </label>
        <input
          v-model="referenceNo"
          type="text"
          placeholder="OP-20260101-CLINIC"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-zinc-600" />
          <span>Unit Bisnis:</span>
        </label>
        <select
          v-model="selectedUnit"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Semua Unit (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Klinik Spesialis Anak)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1">Keterangan Dokumen Saldo Awal:</label>
        <input
          v-model="description"
          type="text"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>
    </div>

    <!-- Real-time Balance Status Bar -->
    <div
      :class="[
        'p-5 rounded-2xl border transition-all shadow-2xs',
        isBalanced
          ? 'bg-emerald-50/70 border-emerald-200 text-emerald-950'
          : 'bg-amber-50/70 border-amber-200 text-amber-950'
      ]"
    >
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <!-- Status Indicator -->
        <div class="flex items-center gap-3">
          <div
            :class="[
              'w-10 h-10 rounded-xl flex items-center justify-center font-bold text-white shrink-0',
              isBalanced ? 'bg-emerald-600 shadow-2xs' : 'bg-amber-600 shadow-2xs'
            ]"
          >
            <CheckCircle2 v-if="isBalanced" class="w-5 h-5" />
            <AlertTriangle v-else class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="text-xs font-black uppercase tracking-wider">
                {{ isBalanced ? 'STATUS: SEIMBANG (BALANCED)' : 'STATUS: BELUM SEIMBANG' }}
              </span>
              <span
                v-if="isBalanced"
                class="px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-800 font-bold text-[10px]"
              >
                Siap Diposting ke Buku Besar
              </span>
              <span
                v-else
                class="px-2 py-0.5 rounded-md bg-amber-100 text-amber-800 font-bold text-[10px]"
              >
                Debit & Kredit Wajib Sama
              </span>
            </div>
            <p class="text-xs mt-0.5 opacity-90">
              {{
                isBalanced
                  ? 'Total Debit sama persis dengan Total Kredit. Jurnal Saldo Awal dapat diposting ke database.'
                  : `Terdapat selisih sebesar Rp ${formatIDR(Math.abs(difference))}. Gunakan tombol "Seimbangkan ke Akun Ekuitas" untuk alokasi otomatis.`
              }}
            </p>
          </div>
        </div>

        <!-- Metric Numbers -->
        <div class="flex flex-wrap items-center gap-4 lg:gap-6 bg-white/90 px-4 py-2.5 rounded-xl border border-zinc-200/80 shadow-2xs">
          <div>
            <span class="text-[10px] font-bold text-zinc-500 uppercase block">Total Debit:</span>
            <span class="tabular-nums text-sm font-extrabold text-zinc-950">Rp {{ formatIDR(totalDebit) }}</span>
          </div>
          <div class="h-8 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] font-bold text-zinc-500 uppercase block">Total Kredit:</span>
            <span class="tabular-nums text-sm font-extrabold text-zinc-950">Rp {{ formatIDR(totalCredit) }}</span>
          </div>
          <div class="h-8 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] font-bold text-zinc-500 uppercase block">Selisih:</span>
            <span
              :class="[
                'tabular-nums text-sm font-extrabold',
                isBalanced ? 'text-emerald-700' : 'text-rose-600'
              ]"
            >
              Rp {{ formatIDR(Math.abs(difference)) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Dynamic Row Input Table (Opsi A: Minimalis & Cepat) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <!-- Toolbar Header -->
      <div class="p-4 border-b border-zinc-200/90 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-zinc-50/50">
        <div class="flex items-center gap-2">
          <Scale class="w-4 h-4 text-zinc-700" />
          <span class="text-xs font-bold text-zinc-900">
            Daftar Baris Akun Bersaldo ({{ balanceRows.length }} Baris Dimasukkan)
          </span>
          <span class="text-[11px] text-zinc-500 hidden sm:inline">
            &bull; Masukkan hanya akun yang memiliki saldo awal (seperti Jurnal Migrasi)
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="addRow"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95 cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Akun</span>
          </button>
          <button
            type="button"
            @click="resetAllBalances"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-100 text-zinc-700 rounded-xl text-xs font-bold transition-all border border-zinc-200 shadow-2xs cursor-pointer"
          >
            <RotateCcw class="w-3.5 h-3.5 text-zinc-500" />
            <span>Reset</span>
          </button>
        </div>
      </div>

      <!-- Main Input Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200/90 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-12 text-center">#</th>
              <th class="px-4 py-3 min-w-[280px]">Akun Saldo Awal (Bagan Akun COA)</th>
              <th class="px-3 py-3 w-28 text-center">Unit Bisnis</th>
              <th class="px-3 py-3 w-24 text-center">Saldo Normal</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Awal Debit (Rp)</th>
              <th class="px-4 py-3 w-48 text-right">Saldo Awal Kredit (Rp)</th>
              <th class="px-4 py-3 min-w-[180px]">Catatan / Keterangan</th>
              <th class="px-3 py-3 w-12 text-center">Hapus</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr
              v-for="(row, idx) in balanceRows"
              :key="idx"
              :class="[
                'hover:bg-zinc-50/70 transition',
                (Number(row.debit) > 0 || Number(row.credit) > 0) ? 'bg-zinc-50/40' : ''
              ]"
            >
              <!-- Index -->
              <td class="px-3 py-2 text-center font-mono font-bold text-zinc-400">
                {{ idx + 1 }}
              </td>

              <!-- Akun Saldo Awal with SearchableSelect -->
              <td class="px-4 py-2">
                <SearchableSelect
                  :model-value="row.code"
                  :options="accountSelectOptions"
                  placeholder="Pilih atau cari akun COA..."
                  @update:model-value="(val: any) => onAccountSelect(idx, String(val))"
                />
              </td>

              <!-- Unit Bisnis -->
              <td class="px-3 py-2 text-center">
                <select
                  v-model="row.unit"
                  class="w-full px-2 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
                >
                  <option value="ALL">ALL</option>
                  <option value="KIDINEST">Kidinest</option>
                  <option value="SSC">SSC</option>
                </select>
              </td>

              <!-- Normal Balance Badge -->
              <td class="px-3 py-2 text-center">
                <span
                  v-if="row.normalBalance"
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    row.normalBalance === 'DEBIT'
                      ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                      : 'bg-rose-50 text-rose-800 border border-rose-200'
                  ]"
                >
                  {{ row.normalBalance }}
                </span>
                <span v-else class="text-zinc-400">-</span>
              </td>

              <!-- Input Debit -->
              <td class="px-4 py-2">
                <div class="relative">
                  <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-[10px]">Rp</span>
                  <input
                    v-model.number="row.debit"
                    type="number"
                    step="any"
                    min="0"
                    placeholder="0"
                    :class="[
                      'w-full pl-8 pr-2.5 py-1.5 rounded-xl border tabular-nums text-right text-xs transition focus:outline-none',
                      Number(row.debit) > 0
                        ? 'border-zinc-900 bg-white font-bold text-zinc-950 ring-1 ring-zinc-900'
                        : 'border-zinc-200 bg-white text-zinc-600 focus:border-zinc-900'
                    ]"
                  />
                </div>
              </td>

              <!-- Input Kredit -->
              <td class="px-4 py-2">
                <div class="relative">
                  <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-[10px]">Rp</span>
                  <input
                    v-model.number="row.credit"
                    type="number"
                    step="any"
                    min="0"
                    placeholder="0"
                    :class="[
                      'w-full pl-8 pr-2.5 py-1.5 rounded-xl border tabular-nums text-right text-xs transition focus:outline-none',
                      Number(row.credit) > 0
                        ? 'border-zinc-900 bg-white font-bold text-zinc-950 ring-1 ring-zinc-900'
                        : 'border-zinc-200 bg-white text-zinc-600 focus:border-zinc-900'
                    ]"
                  />
                </div>
              </td>

              <!-- Catatan -->
              <td class="px-4 py-2">
                <input
                  v-model="row.notes"
                  type="text"
                  placeholder="Keterangan saldo..."
                  class="w-full px-2.5 py-1.5 rounded-xl border border-zinc-200 bg-white text-zinc-700 text-xs focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>

              <!-- Hapus Row -->
              <td class="px-3 py-2 text-center">
                <button
                  type="button"
                  @click="removeRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                  title="Hapus baris ini"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>

            <tr v-if="balanceRows.length === 0">
              <td colspan="8" class="p-8 text-center text-zinc-400 italic">
                Belum ada baris akun. Klik "Tambah Akun" atau "Muat Data Standar Klinik".
              </td>
            </tr>
          </tbody>

          <!-- Table Footer Summary -->
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-zinc-900">
            <tr>
              <td colspan="4" class="px-4 py-3 text-right uppercase tracking-wider text-[11px] text-zinc-500">
                Total Saldo Awal ({{ balanceRows.length }} Baris):
              </td>
              <td class="px-4 py-3 text-right tabular-nums text-sm font-extrabold text-zinc-950">
                Rp {{ formatIDR(totalDebit) }}
              </td>
              <td class="px-4 py-3 text-right tabular-nums text-sm font-extrabold text-zinc-950">
                Rp {{ formatIDR(totalCredit) }}
              </td>
              <td colspan="2" class="px-4 py-3 text-right">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-black uppercase',
                    isBalanced
                      ? 'bg-emerald-100 text-emerald-800 border border-emerald-300'
                      : 'bg-rose-100 text-rose-800 border border-rose-300'
                  ]"
                >
                  {{ isBalanced ? 'SEIMBANG' : `SELISIH: Rp ${formatIDR(Math.abs(difference))}` }}
                </span>
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Add Row Bottom Action -->
      <div class="p-3 bg-zinc-50/50 border-t border-zinc-200/90 flex justify-center">
        <button
          type="button"
          @click="addRow"
          class="inline-flex items-center gap-2 px-4 py-2 bg-white hover:bg-zinc-100 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all cursor-pointer active:scale-95"
        >
          <Plus class="w-4 h-4 text-zinc-600" />
          <span>Tambah Baris Akun Saldo Awal</span>
        </button>
      </div>
    </div>

    <!-- Sticky Bottom Balance Bar (Clinical Fast-Action Dock) -->
    <div class="fixed bottom-4 left-68 right-8 z-30 bg-white/95 backdrop-blur-md rounded-2xl border border-zinc-300/90 p-4 shadow-xl flex flex-col sm:flex-row items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-5 text-xs">
        <div>
          <span class="text-[10px] uppercase font-bold text-zinc-500">Total Debit</span>
          <div class="text-sm font-extrabold tabular-nums text-zinc-950">Rp {{ formatIDR(totalDebit) }}</div>
        </div>
        <div class="w-px h-7 bg-zinc-200"></div>
        <div>
          <span class="text-[10px] uppercase font-bold text-zinc-500">Total Kredit</span>
          <div class="text-sm font-extrabold tabular-nums text-zinc-950">Rp {{ formatIDR(totalCredit) }}</div>
        </div>
        <div class="w-px h-7 bg-zinc-200"></div>
        <div>
          <span class="text-[10px] uppercase font-bold text-zinc-500">Keseimbangan Neraca</span>
          <div class="mt-0.5">
            <span
              :class="[
                'px-2.5 py-0.5 rounded-full text-[11px] font-bold inline-flex items-center gap-1.5',
                isBalanced
                  ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                  : 'bg-rose-50 text-rose-800 border border-rose-200'
              ]"
            >
              <span :class="['w-1.5 h-1.5 rounded-full', isBalanced ? 'bg-emerald-500' : 'bg-rose-500']"></span>
              <span>{{ isBalanced ? 'SEIMBANG (Rp 0)' : `SELISIH: Rp ${formatIDR(Math.abs(difference))}` }}</span>
            </span>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="!isBalanced && difference !== 0"
          @click="autoBalanceToEquity"
          class="px-3.5 py-2 rounded-xl text-xs font-bold bg-amber-500 hover:bg-amber-600 text-white transition-all shadow-2xs active:scale-95 cursor-pointer"
        >
          Seimbangkan Otomatis
        </button>
        <button
          @click="handleSaveOpeningBalance"
          :disabled="!isBalanced || isSubmitting"
          :class="[
            'inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95 cursor-pointer',
            isBalanced && !isSubmitting
              ? 'bg-zinc-950 hover:bg-zinc-800 text-white'
              : 'bg-zinc-100 text-zinc-400 border border-zinc-200 cursor-not-allowed'
          ]"
        >
          <Send class="w-3.5 h-3.5" />
          <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan & Posting Saldo Awal' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>
