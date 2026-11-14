<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
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
  FolderOpen
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()

// Opening Balance Header Form
const cutOffDate = ref('2026-01-01')
const referenceNo = ref('OP-20260101-BCS')
const selectedUnit = ref<BusinessUnit>('ALL')
const description = ref('Saldo Awal Pembukuan Akuntansi & Neraca Migrasi PT Brayo Cipta Sejahtera per 1 Januari 2026')

// Search & Filter
const searchQuery = ref('')
const selectedCategory = ref<'ALL' | 'ASSET' | 'LIABILITY' | 'EQUITY' | 'NOMINAL'>('ALL')

// Input Rows State: Map of accountCode -> { debit: number, credit: number, notes?: string }
interface OpeningBalanceRowState {
  code: string
  name: string
  type: AccountType
  normalBalance: NormalBalance
  level: number
  isHeader: boolean
  parentCode?: string
  debit: number
  credit: number
  notes: string
}

const balanceRows = ref<OpeningBalanceRowState[]>([])

// Populate rows from ledgerStore.accounts
function initializeRows() {
  balanceRows.value = ledgerStore.accounts
    .filter(a => !a.isHeader) // Only leaf accounts can take opening balance
    .map(a => {
      const cur = a.currentBalance || 0
      let d = 0
      let c = 0
      if (a.normalBalance === 'DEBIT') {
        d = cur >= 0 ? cur : 0
        c = cur < 0 ? Math.abs(cur) : 0
      } else {
        c = cur >= 0 ? cur : 0
        d = cur < 0 ? Math.abs(cur) : 0
      }

      return {
        code: a.code,
        name: a.name,
        type: a.type,
        normalBalance: a.normalBalance,
        level: a.level,
        isHeader: a.isHeader,
        parentCode: a.parentCode,
        debit: d,
        credit: c,
        notes: ''
      }
    })
}

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  initializeRows()
})

// Filtered Accounts
const filteredRows = computed(() => {
  return balanceRows.value.filter(r => {
    // Category filter
    if (selectedCategory.value === 'ASSET' && r.type !== 'ASSET') return false
    if (selectedCategory.value === 'LIABILITY' && r.type !== 'LIABILITY') return false
    if (selectedCategory.value === 'EQUITY' && r.type !== 'EQUITY') return false
    if (selectedCategory.value === 'NOMINAL' && r.type !== 'REVENUE' && r.type !== 'EXPENSE' && r.type !== 'COGS') return false

    // Search query
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      const matchCode = r.code.toLowerCase().includes(q)
      const matchName = r.name.toLowerCase().includes(q)
      if (!matchCode && !matchName) return false
    }

    return true
  })
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

// Action: Muat Angka Standar Saldo Awal BCS
function loadBcsStandardBalances() {
  initializeRows()
  accountingStore.showToast('Template Saldo Awal Standar PT. Brayo Cipta Sejahtera berhasil dimuat!', 'info')
}

// Action: Seimbangkan Otomatis ke Akun Ekuitas Saldo Awal (3999)
function autoBalanceToEquity() {
  const diff = difference.value
  if (Math.abs(diff) < 0.01) {
    accountingStore.showToast('Saldo awal sudah seimbang!', 'info')
    return
  }

  // Find 3999 - Opening Balance Equity - BCS
  let eqRow = balanceRows.value.find(r => r.code === '3999')
  if (!eqRow) {
    // Fallback to 1999 or 3201
    eqRow = balanceRows.value.find(r => r.code === '1999' || r.code === '3201')
  }

  if (!eqRow) {
    accountingStore.showToast('Akun penyeimbang saldo awal (3999 / 1999) tidak ditemukan', 'error')
    return
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
  if (confirm('Apakah Anda yakin ingin mengosongkan seluruh input saldo awal?')) {
    balanceRows.value.forEach(r => {
      r.debit = 0
      r.credit = 0
      r.notes = ''
    })
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
      router.push('/ledger')
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
  <div class="space-y-6">
    <!-- Header Page -->
    <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2.5 py-0.5 rounded border border-blue-200">
            Setup Buku Besar & Migrasi Neraca
          </span>
          <span class="text-xs font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200 flex items-center gap-1">
            <Database class="w-3.5 h-3.5 text-emerald-600" />
            <span>PostgreSQL Live Persistence</span>
          </span>
        </div>
        <h1 class="text-2xl font-black text-slate-900 mt-1">
          Pembuatan & Penyetelan Saldo Awal Akuntansi
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Penetapan angka saldo awal (cut-off neraca) untuk seluruh akun COA PT. Brayo Cipta Sejahtera per awal periode pembukuan.
        </p>
      </div>

      <!-- Quick Actions -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          @click="loadBcsStandardBalances"
          class="flex items-center gap-1.5 px-3 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-bold transition border border-slate-300"
          title="Muat nilai saldo awal standar BCS"
        >
          <FolderOpen class="w-3.5 h-3.5 text-slate-600" />
          <span>Muat Data Standar BCS</span>
        </button>

        <button
          @click="autoBalanceToEquity"
          :disabled="isBalanced || difference === 0"
          :class="[
            'flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold transition shadow-xs',
            !isBalanced && difference !== 0
              ? 'bg-amber-600 hover:bg-amber-700 text-white cursor-pointer'
              : 'bg-slate-100 text-slate-400 cursor-not-allowed'
          ]"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>Seimbangkan ke Akun Ekuitas (3999)</span>
        </button>

        <button
          @click="handleSaveOpeningBalance"
          :disabled="!isBalanced || isSubmitting"
          :class="[
            'flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold transition shadow-xs',
            isBalanced && !isSubmitting
              ? 'bg-blue-600 hover:bg-blue-700 text-white cursor-pointer shadow-blue-500/20'
              : 'bg-slate-200 text-slate-400 cursor-not-allowed'
          ]"
        >
          <Send class="w-3.5 h-3.5" />
          <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan & Posting Saldo Awal' }}</span>
        </button>

        <button
          @click="handlePrint"
          class="flex items-center gap-1 px-3 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-xl text-xs font-bold transition border border-slate-300 shadow-2xs"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Configuration Form Bar -->
    <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
      <div>
        <label class="block font-bold text-slate-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-blue-600" />
          <span>Tanggal Cut-Off Saldo Awal:</span>
        </label>
        <input
          v-model="cutOffDate"
          type="date"
          class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg font-bold text-slate-800 focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div>
        <label class="block font-bold text-slate-700 mb-1 flex items-center gap-1">
          <FileText class="w-3.5 h-3.5 text-purple-600" />
          <span>Nomor Dokumen Bukti:</span>
        </label>
        <input
          v-model="referenceNo"
          type="text"
          placeholder="OP-20260101-BCS"
          class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg font-mono font-bold text-slate-800 focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div>
        <label class="block font-bold text-slate-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-emerald-600" />
          <span>Entitas Bisnis:</span>
        </label>
        <select
          v-model="selectedUnit"
          class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg font-bold text-slate-800 focus:ring-2 focus:ring-blue-500"
        >
          <option value="ALL">Kidinest + SSC (Konsolidasi BCS)</option>
          <option value="KIDINEST">Unit Kidinest (Pediatrik)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-slate-700 mb-1">Keterangan Dokumen Saldo Awal:</label>
        <input
          v-model="description"
          type="text"
          class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg text-slate-800 focus:ring-2 focus:ring-blue-500"
        />
      </div>
    </div>

    <!-- Real-time Balance Status Bar -->
    <div
      :class="[
        'p-5 rounded-2xl border transition-all shadow-xs',
        isBalanced
          ? 'bg-emerald-50/80 border-emerald-300 text-emerald-950'
          : 'bg-amber-50/90 border-amber-300 text-amber-950'
      ]"
    >
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <!-- Status Indicator -->
        <div class="flex items-center gap-3">
          <div
            :class="[
              'w-10 h-10 rounded-xl flex items-center justify-center font-bold text-white shrink-0',
              isBalanced ? 'bg-emerald-600 shadow-md shadow-emerald-600/30' : 'bg-amber-600 shadow-md shadow-amber-600/30'
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
                class="px-2 py-0.5 rounded bg-emerald-200/80 text-emerald-800 font-bold text-[10px]"
              >
                Siap Diposting ke Buku Besar
              </span>
              <span
                v-else
                class="px-2 py-0.5 rounded bg-amber-200 text-amber-800 font-bold text-[10px]"
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
        <div class="flex flex-wrap items-center gap-4 lg:gap-6 bg-white/70 px-4 py-2.5 rounded-xl border border-white/60">
          <div>
            <span class="text-[10px] font-bold text-slate-500 uppercase block">Total Debit:</span>
            <span class="font-mono text-sm font-black text-blue-700">Rp {{ formatIDR(totalDebit) }}</span>
          </div>
          <div class="h-8 w-px bg-slate-200"></div>
          <div>
            <span class="text-[10px] font-bold text-slate-500 uppercase block">Total Kredit:</span>
            <span class="font-mono text-sm font-black text-purple-700">Rp {{ formatIDR(totalCredit) }}</span>
          </div>
          <div class="h-8 w-px bg-slate-200"></div>
          <div>
            <span class="text-[10px] font-bold text-slate-500 uppercase block">Selisih (Balance):</span>
            <span
              :class="[
                'font-mono text-sm font-black',
                isBalanced ? 'text-emerald-700' : 'text-rose-600'
              ]"
            >
              Rp {{ formatIDR(Math.abs(difference)) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Table Section -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <!-- Toolbar Filter & Search -->
      <div class="p-4 border-b border-slate-200 flex flex-col md:flex-row md:items-center justify-between gap-3">
        <!-- Category Tabs -->
        <div class="flex flex-wrap items-center gap-1">
          <button
            v-for="cat in [
              { id: 'ALL', label: 'Semua Akun' },
              { id: 'ASSET', label: '1 - Aset' },
              { id: 'LIABILITY', label: '2 - Kewajiban' },
              { id: 'EQUITY', label: '3 - Ekuitas' },
              { id: 'NOMINAL', label: '4-7 Nominal' }
            ]"
            :key="cat.id"
            @click="selectedCategory = cat.id as any"
            :class="[
              'px-3 py-1.5 rounded-lg text-xs font-bold transition',
              selectedCategory === cat.id
                ? 'bg-blue-600 text-white shadow-xs'
                : 'text-slate-600 hover:bg-slate-100'
            ]"
          >
            {{ cat.label }}
          </button>
        </div>

        <!-- Search Input -->
        <div class="relative w-full md:w-64">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari kode atau nama akun..."
            class="w-full pl-9 pr-3 py-1.5 bg-slate-50 border border-slate-300 rounded-lg text-xs text-slate-800 focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- Main Input Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-slate-50 text-slate-600 font-bold border-b border-slate-200 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-4 py-3 w-32">Kode Akun</th>
              <th class="px-4 py-3 min-w-[240px]">Nama Akun COA</th>
              <th class="px-3 py-3 w-28 text-center">Tipe / Saldo</th>
              <th class="px-4 py-3 w-52 text-right">Saldo Awal Debit (Rp)</th>
              <th class="px-4 py-3 w-52 text-right">Saldo Awal Kredit (Rp)</th>
              <th class="px-4 py-3 min-w-[180px]">Catatan / Keterangan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 text-slate-700">
            <tr
              v-for="row in filteredRows"
              :key="row.code"
              :class="[
                'hover:bg-blue-50/40 transition',
                (row.debit > 0 || row.credit > 0) ? 'bg-slate-50/60 font-semibold' : ''
              ]"
            >
              <!-- Kode Akun -->
              <td class="px-4 py-2.5 font-mono font-bold text-slate-900">
                {{ row.code }}
              </td>

              <!-- Nama Akun -->
              <td class="px-4 py-2.5">
                <div class="flex items-center gap-2">
                  <span
                    v-if="row.code === '3999' || row.code === '1999'"
                    class="px-1.5 py-0.5 rounded bg-amber-100 text-amber-800 font-bold text-[9px]"
                  >
                    PENYEIMBANG
                  </span>
                  <span>{{ row.name }}</span>
                </div>
              </td>

              <!-- Tipe Akun & Normal Balance -->
              <td class="px-3 py-2.5 text-center">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px] font-bold',
                    row.normalBalance === 'DEBIT'
                      ? 'bg-blue-50 text-blue-700 border border-blue-200'
                      : 'bg-purple-50 text-purple-700 border border-purple-200'
                  ]"
                >
                  {{ row.normalBalance }}
                </span>
              </td>

              <!-- Input Debit -->
              <td class="px-4 py-2">
                <div class="relative">
                  <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 font-mono text-[10px]">Rp</span>
                  <input
                    v-model.number="row.debit"
                    type="number"
                    step="any"
                    min="0"
                    placeholder="0"
                    :class="[
                      'w-full pl-8 pr-2.5 py-1.5 rounded-lg border font-mono text-right text-xs transition',
                      row.debit > 0
                        ? 'border-blue-400 bg-blue-50/50 font-bold text-blue-900 ring-1 ring-blue-300'
                        : 'border-slate-200 bg-white text-slate-600 focus:border-blue-500'
                    ]"
                  />
                </div>
              </td>

              <!-- Input Kredit -->
              <td class="px-4 py-2">
                <div class="relative">
                  <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400 font-mono text-[10px]">Rp</span>
                  <input
                    v-model.number="row.credit"
                    type="number"
                    step="any"
                    min="0"
                    placeholder="0"
                    :class="[
                      'w-full pl-8 pr-2.5 py-1.5 rounded-lg border font-mono text-right text-xs transition',
                      row.credit > 0
                        ? 'border-purple-400 bg-purple-50/50 font-bold text-purple-900 ring-1 ring-purple-300'
                        : 'border-slate-200 bg-white text-slate-600 focus:border-purple-500'
                    ]"
                  />
                </div>
              </td>

              <!-- Catatan -->
              <td class="px-4 py-2">
                <input
                  v-model="row.notes"
                  type="text"
                  placeholder="Keterangan saldo awal..."
                  class="w-full px-2.5 py-1.5 rounded-lg border border-slate-200 bg-white text-slate-600 text-xs focus:ring-1 focus:ring-blue-500"
                />
              </td>
            </tr>

            <tr v-if="filteredRows.length === 0">
              <td colspan="6" class="p-8 text-center text-slate-400 italic">
                Tidak ada akun yang sesuai dengan filter pencarian.
              </td>
            </tr>
          </tbody>

          <!-- Table Footer Summary -->
          <tfoot class="bg-slate-100 font-bold border-t-2 border-slate-300 text-slate-900">
            <tr>
              <td colspan="3" class="px-4 py-3 text-right uppercase tracking-wider text-[11px]">
                Total Saldo Awal:
              </td>
              <td class="px-4 py-3 text-right font-mono text-sm text-blue-800">
                Rp {{ formatIDR(totalDebit) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-sm text-purple-800">
                Rp {{ formatIDR(totalCredit) }}
              </td>
              <td class="px-4 py-3 text-right">
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
    </div>
  </div>
</template>
