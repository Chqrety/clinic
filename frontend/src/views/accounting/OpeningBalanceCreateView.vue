<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { api } from '@/api/client'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import CurrencyInput from '@/components/common/CurrencyInput.vue'
import type { BusinessUnit, AccountType, NormalBalance } from '@/types/accounting'
import {
  Folder,
  ArrowLeft,
  BookOpen,
  Scale,
  Calendar,
  Building2,
  FileText,
  Search,
  RotateCcw,
  Sparkles,
  Send,
  Save,
  AlertTriangle,
  CheckCircle2,
  Clock,
  ExternalLink,
  Plus,
  Trash2,
  ListPlus,
  Check,
  ShieldCheck,
  Layers
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()

// Active Tab matching Brayo / SPT Masa layout
const activeTab = ref<'DETAILS' | 'RECONCILIATION' | 'AUDIT'>('DETAILS')

// Header Form State
const cutOffDate = ref('2026-01-01')
const referenceNo = ref('OP-20260101-BCS')
const selectedUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const description = ref('Saldo Awal Pembukuan Akuntansi & Neraca Migrasi PT. Semarang Medical Group (Kidinest & SSC) per 1 Januari 2026')
const sourceBadge = ref('Input Manual / Penyesuaian')

// Balance Row Structure
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

const rows = ref<BalanceRow[]>([])
const isSubmitting = ref(false)

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  initializeRows()
})

// Initialize rows: load existing accounts with balances if available, otherwise start with 2 clean blank rows
function initializeRows() {
  const accountsWithBalance = ledgerStore.accounts.filter(
    a => !a.isHeader && a.currentBalance && a.currentBalance !== 0
  )

  if (accountsWithBalance.length > 0) {
    sourceBadge.value = 'Data Migrasi / Buku Besar Aktif'
    rows.value = accountsWithBalance.map(acc => {
      const bal = Number(acc.currentBalance) || 0
      return {
        code: acc.code,
        name: acc.name,
        type: acc.type,
        normalBalance: acc.normalBalance,
        unit: (acc.businessUnit as BusinessUnit) || 'ALL',
        debit: acc.normalBalance === 'DEBIT' ? (bal > 0 ? bal : 0) : (bal < 0 ? Math.abs(bal) : 0),
        credit: acc.normalBalance === 'CREDIT' ? (bal > 0 ? bal : 0) : (bal < 0 ? Math.abs(bal) : 0),
        notes: `Saldo awal buku besar ${acc.name}`
      }
    })
  } else {
    rows.value = [
      { code: '', name: '', type: 'ASSET', normalBalance: 'DEBIT', unit: 'ALL', debit: 0, credit: 0, notes: '' },
      { code: '', name: '', type: 'LIABILITY', normalBalance: 'CREDIT', unit: 'ALL', debit: 0, credit: 0, notes: '' }
    ]
  }
}

// Available accounts dropdown options
const availableAccountOptions = computed(() => {
  return ledgerStore.accounts
    .filter(a => !a.isHeader && a.isActive !== false)
    .map(a => ({
      value: a.code,
      label: `[${a.code}] ${a.name}`,
      sublabel: `${a.type} • Saldo Normal: ${a.normalBalance}`
    }))
})

function addRow() {
  rows.value.push({
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

function removeRow(index: number) {
  if (rows.value.length > 1) {
    rows.value.splice(index, 1)
  } else {
    rows.value[0].code = ''
    rows.value[0].name = ''
    rows.value[0].debit = 0
    rows.value[0].credit = 0
    rows.value[0].notes = ''
  }
}

function onAccountSelect(index: number, code: string) {
  const acc = ledgerStore.accounts.find(a => a.code === code)
  if (acc && rows.value[index]) {
    rows.value[index].code = acc.code
    rows.value[index].name = acc.name
    rows.value[index].type = acc.type
    rows.value[index].normalBalance = acc.normalBalance
    rows.value[index].unit = (acc.businessUnit as BusinessUnit) || selectedUnit.value || 'ALL'
    if (!rows.value[index].notes) {
      rows.value[index].notes = `Saldo awal ${acc.name}`
    }
  }
}

// Load entire active leaf COA into table
function loadAllLeafAccounts() {
  const leafAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
  if (leafAccounts.length === 0) {
    accountingStore.showToast('Bagan akun (COA) belum tersedia!', 'info')
    return
  }

  rows.value = leafAccounts.map(acc => {
    const cur = Number(acc.currentBalance) || 0
    return {
      code: acc.code,
      name: acc.name,
      type: acc.type,
      normalBalance: acc.normalBalance,
      unit: (acc.businessUnit as BusinessUnit) || 'ALL',
      debit: acc.normalBalance === 'DEBIT' ? (cur > 0 ? cur : 0) : 0,
      credit: acc.normalBalance === 'CREDIT' ? (cur > 0 ? cur : 0) : 0,
      notes: cur > 0 ? 'Saldo migrasi berjalan' : ''
    }
  })

  accountingStore.showToast(`Seluruh ${leafAccounts.length} akun COA berhasil dimuat ke tabel saldo awal!`, 'info')
}

// ⚡ Automatic Pull: Tarik Data Saldo dari Kasir & Buku Besar
async function pullDataFromCashierAndJournals() {
  try {
    sourceBadge.value = 'Data Terintegrasi dari Kasir & Buku Besar'
    const newRows: BalanceRow[] = []

    // 1. Fetch live accounts & existing balances
    await ledgerStore.fetchAccounts()
    const activeAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)

    // 2. Fetch cashier invoices/reports
    let cashierCash = 0
    let cashierBank = 0
    let cashierReceivable = 0
    let cashierRevenue = 0

    try {
      const invoices = await api.getCashierInvoices(100, 0)
      if (Array.isArray(invoices) && invoices.length > 0) {
        invoices.forEach(inv => {
          if (inv.status === 'PAID') {
            const amount = Number(inv.total_amount) || 0
            cashierRevenue += amount
            if (inv.payment_method === 'CASH') {
              cashierCash += amount
            } else {
              cashierBank += amount
            }
          }
        })
      }
    } catch {
      // Cashier API fallback if offline or no shift today
    }

    // 3. Populate rows for active balance accounts
    activeAccounts.forEach(acc => {
      const cur = Number(acc.currentBalance) || 0
      let d = acc.normalBalance === 'DEBIT' ? (cur > 0 ? cur : 0) : 0
      let c = acc.normalBalance === 'CREDIT' ? (cur > 0 ? cur : 0) : 0

      // Add cashier shift integration if relevant
      if ((acc.code === '1101.001' || acc.code === '1-1111') && cashierCash > 0) {
        d += cashierCash
      } else if ((acc.code === '1102.001' || acc.code === '1-1113') && cashierBank > 0) {
        d += cashierBank
      } else if ((acc.code === '1201' || acc.code === '1-1120') && cashierReceivable > 0) {
        d += cashierReceivable
      } else if ((acc.code.startsWith('4-') || acc.code.startsWith('41')) && cashierRevenue > 0) {
        c += cashierRevenue
      }

      if (d > 0 || c > 0) {
        newRows.push({
          code: acc.code,
          name: acc.name,
          type: acc.type,
          normalBalance: acc.normalBalance,
          unit: (acc.businessUnit as BusinessUnit) || 'ALL',
          debit: d,
          credit: c,
          notes: `Saldo terintegrasi mutasi kasir & buku besar`
        })
      }
    })

    if (newRows.length > 0) {
      rows.value = newRows
      autoBalanceToEquity()
      accountingStore.showToast(`Berhasil menarik ${newRows.length} akun saldo dari transaksi kasir & buku besar!`, 'success')
    } else {
      // If system has no transactions yet, load clean standard rows
      loadAllLeafAccounts()
      accountingStore.showToast('Belum ada riwayat mutasi kasir. Menampilkan daftar akun untuk pengisian awal.', 'info')
    }
  } catch (err: any) {
    accountingStore.showToast(`Gagal menarik data kasir: ${err.message}`, 'error')
  }
}

// Computations
const totalDebit = computed(() => {
  return rows.value.reduce((sum, r) => sum + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return rows.value.reduce((sum, r) => sum + (Number(r.credit) || 0), 0)
})

const difference = computed(() => {
  return Math.round((totalDebit.value - totalCredit.value) * 100) / 100
})

const isBalanced = computed(() => {
  return Math.abs(difference.value) < 0.01 && totalDebit.value > 0
})

// Accounting Equation: Asset = Liability + Equity
const totalAsset = computed(() => {
  return rows.value
    .filter(r => r.type === 'ASSET')
    .reduce((sum, r) => sum + (Number(r.debit) || 0) - (Number(r.credit) || 0), 0)
})

const totalLiability = computed(() => {
  return rows.value
    .filter(r => r.type === 'LIABILITY')
    .reduce((sum, r) => sum + (Number(r.credit) || 0) - (Number(r.debit) || 0), 0)
})

const totalEquity = computed(() => {
  return rows.value
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

  let eqRow = rows.value.find(r => r.code === '3999' || r.code === '3-9999' || r.code.includes('9999'))

  if (!eqRow) {
    const acc = ledgerStore.accounts.find(a => a.code === '3999' || a.code === '3-9999') ||
                ledgerStore.accounts.find(a => a.type === 'EQUITY' && !a.isHeader)
    eqRow = {
      code: acc?.code || '3999',
      name: acc?.name || 'Opening Balance Equity (Modal Penyeimbang)',
      type: 'EQUITY',
      normalBalance: 'CREDIT',
      unit: 'ALL',
      debit: 0,
      credit: 0,
      notes: 'Penyeimbang Saldo Awal Neraca'
    }
    rows.value.push(eqRow)
  }

  if (diff > 0) {
    // Debit is higher -> add to Credit of Equity
    eqRow.credit = (Number(eqRow.credit) || 0) + diff
  } else {
    // Credit is higher -> add to Debit of Equity
    eqRow.debit = (Number(eqRow.debit) || 0) + Math.abs(diff)
  }

  accountingStore.showToast(
    `Selisih Rp ${formatIDR(Math.abs(diff))} dialokasikan ke akun Modal Penyeimbang (${eqRow.code} - ${eqRow.name})!`,
    'success'
  )
}

function resetInputs() {
  if (confirm('Kosongkan seluruh baris input saldo awal?')) {
    rows.value = [
      { code: '', name: '', type: 'ASSET', normalBalance: 'DEBIT', unit: 'ALL', debit: 0, credit: 0, notes: '' },
      { code: '', name: '', type: 'LIABILITY', normalBalance: 'CREDIT', unit: 'ALL', debit: 0, credit: 0, notes: '' }
    ]
    sourceBadge.value = 'Input Manual / Penyesuaian'
    accountingStore.showToast('Tabel saldo awal telah dikosongkan.', 'info')
  }
}

// Submit Opening Balance
async function handleSaveOpeningBalance(autoPost: boolean = true) {
  const validLines = rows.value.filter(r => r.code && (Number(r.debit) > 0 || Number(r.credit) > 0))

  if (validLines.length < 2) {
    accountingStore.showToast('Saldo Awal minimal memiliki 2 baris akun yang valid!', 'error')
    return
  }

  if (!isBalanced.value) {
    accountingStore.showToast(
      `Saldo Awal belum seimbang! Total Debit (${formatIDR(totalDebit.value)}) != Total Kredit (${formatIDR(totalCredit.value)}). Selisih: ${formatIDR(Math.abs(difference.value))}. Klik "Seimbangkan Selisih" terlebih dahulu.`,
      'error'
    )
    return
  }

  try {
    isSubmitting.value = true

    await ledgerStore.saveOpeningBalance({
      cutOffDate: cutOffDate.value,
      referenceNo: referenceNo.value || `OP-${cutOffDate.value.replace(/-/g, '')}`,
      businessUnit: selectedUnit.value,
      description: description.value || `Saldo Awal Pembukuan Akuntansi Klinik per ${cutOffDate.value}`,
      lines: validLines.map(r => ({
        accountCode: r.code,
        accountName: r.name,
        debit: Number(r.debit) || 0,
        credit: Number(r.credit) || 0,
        description: r.notes || description.value
      }))
    })

    accountingStore.showToast(
      `Saldo Awal Akuntansi berhasil ${autoPost ? 'resmi diposting ke Buku Besar!' : 'disimpan sebagai draf review!'}`,
      'success'
    )

    router.push('/accounting/opening-balance')
  } catch (err: any) {
    accountingStore.showToast(`Gagal menyimpan Saldo Awal: ${err.message}`, 'error')
  } finally {
    isSubmitting.value = false
  }
}

// Audit trail: opening balance journals list
const openingJournals = computed(() => {
  return journalsStore.journals.filter(
    j => j.sourceReference?.startsWith('OP-') || j.journalNumber?.startsWith('OP-') || j.description?.toLowerCase().includes('saldo awal')
  )
})
</script>

<template>
  <div class="font-['Plus_Jakarta_Sans',sans-serif] max-w-7xl mx-auto pb-16">
    <div class="bg-white text-zinc-800 rounded-2xl border border-zinc-200/90 p-5 sm:p-7 shadow-2xs space-y-6">
      
      <!-- 1. Header: Breadcrumb & Top Actions -->
      <header class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-3 border-b border-zinc-200">
        <!-- Breadcrumb & Back -->
        <div class="flex items-center gap-2 text-xs text-zinc-500 flex-wrap">
          <router-link
            to="/accounting/opening-balance"
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-zinc-500 hover:text-zinc-900 transition mr-2"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Daftar</span>
          </router-link>
          <span class="text-zinc-300">/</span>
          <router-link to="/accounting/opening-balance" class="hover:text-zinc-900 transition flex items-center gap-1.5 text-zinc-500">
            <Folder class="w-4 h-4 text-zinc-400" />
            <span>Saldo Awal</span>
          </router-link>
          <span>/</span>
          <span class="font-mono font-bold text-zinc-900">{{ referenceNo || 'OP-NEW' }}</span>
          <span
            :class="[
              'ml-1.5 px-2 py-0.5 rounded text-[11px] font-bold border',
              isBalanced ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-amber-50 text-amber-700 border-amber-200'
            ]"
          >
            {{ isBalanced ? 'Seimbang' : 'Belum Seimbang' }}
          </span>
        </div>

        <!-- Top Right Actions -->
        <div class="flex items-center gap-2 flex-wrap">
          <!-- Ambil Saldo Kasir & Jurnal Otomatis -->
          <button
            type="button"
            @click="pullDataFromCashierAndJournals"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold rounded-xl text-xs border border-amber-200 transition active:scale-95 cursor-pointer"
            title="Tarik data mutasi kasir dan buku besar otomatis ke draf saldo awal"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Ambil Saldo Kasir & Jurnal</span>
          </button>

          <!-- Seimbangkan Selisih ke Ekuitas 3999 -->
          <button
            v-if="!isBalanced && difference !== 0"
            type="button"
            @click="autoBalanceToEquity"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 font-bold rounded-xl text-xs border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
            title="Seimbangkan selisih ke akun modal/ekuitas penyeimbang (3999)"
          >
            <Scale class="w-3.5 h-3.5 text-zinc-600" />
            <span>Seimbangkan Selisih</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting || !isBalanced"
            @click="handleSaveOpeningBalance(false)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-white hover:bg-zinc-50 disabled:opacity-50 text-zinc-800 font-bold rounded-xl text-xs border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Save class="w-3.5 h-3.5 text-zinc-600" />
            <span>Simpan Draft</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting || !isBalanced"
            @click="handleSaveOpeningBalance(true)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white font-bold rounded-xl text-xs shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Send class="w-3.5 h-3.5 text-emerald-400" />
            <span>Posting ke Buku Besar</span>
          </button>

          <router-link
            to="/accounting/opening-balance"
            class="inline-flex items-center px-3 py-1.5 text-zinc-500 hover:text-zinc-800 text-xs font-semibold transition"
          >
            Cancel
          </router-link>
        </div>
      </header>

      <!-- 2. Tabs Navigation -->
      <div class="border-b border-zinc-200">
        <nav class="flex space-x-6 text-xs font-medium">
          <button
            type="button"
            @click="activeTab = 'DETAILS'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'DETAILS' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Formulir Saldo Awal</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'RECONCILIATION'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'RECONCILIATION' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <Scale class="w-3.5 h-3.5" />
            <span>Rekonsiliasi Neraca & Status</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'AUDIT'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'AUDIT' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <Clock class="w-3.5 h-3.5" />
            <span>Audit Trail & Riwayat Posting</span>
          </button>
        </nav>
      </div>

      <!-- 3. TAB 1: DETAILS (Formulir Saldo Awal) -->
      <div v-if="activeTab === 'DETAILS'" class="space-y-7">
        
        <!-- Section A: Informasi Dokumen Saldo Awal -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Informasi Dokumen & Entitas Saldo Awal
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Tanggal Cut-Off Saldo Awal *</label>
              <input
                v-model="cutOffDate"
                type="date"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">No. Dokumen / Referensi *</label>
              <input
                v-model="referenceNo"
                type="text"
                placeholder="e.g. OP-20260101-BCS"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Unit Bisnis / Entitas *</label>
              <select
                v-model="selectedUnit"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Unit Kidinest (Klinik Anak)</option>
                <option value="SSC">Unit SSC (Semarang Skin Center)</option>
              </select>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Sumber Integrasi Data</label>
              <input
                type="text"
                readonly
                :value="sourceBadge"
                class="w-full px-3 py-2 bg-zinc-100/80 border border-zinc-200 rounded-xl text-zinc-600 font-semibold cursor-not-allowed"
              />
            </div>

            <div class="md:col-span-2 lg:col-span-4">
              <label class="block text-zinc-600 mb-1 font-medium">Catatan / Keterangan Pembukuan Saldo Awal</label>
              <input
                v-model="description"
                type="text"
                placeholder="Tuliskan keterangan lengkap maksud pembukuan saldo awal migrasi ini..."
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-medium focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>
          </div>
        </div>

        <!-- Section B: Status Keseimbangan (Balance Banner) -->
        <div
          :class="[
            'p-4 rounded-2xl border transition flex flex-col md:flex-row md:items-center justify-between gap-4 text-xs',
            isBalanced
              ? 'bg-emerald-50/60 border-emerald-200 text-emerald-900'
              : 'bg-amber-50/70 border-amber-200 text-amber-900'
          ]"
        >
          <div class="flex items-center gap-2.5">
            <CheckCircle2 v-if="isBalanced" class="w-5 h-5 text-emerald-600 shrink-0" />
            <AlertTriangle v-else class="w-5 h-5 text-amber-600 shrink-0" />
            <div>
              <p class="font-bold">
                {{ isBalanced ? 'STATUS: SALDO AWAL SEIMBANG (BALANCED)' : 'STATUS: BELUM SEIMBANG (UNBALANCED)' }}
              </p>
              <p class="text-[11px] opacity-80">
                {{ isBalanced ? 'Total Debet sama dengan Total Kredit. Siap diposting ke Buku Besar.' : `Total Debet dan Kredit belum sama persis (Selisih: ${formatIDR(Math.abs(difference))}). Klik "Seimbangkan Selisih" untuk menyeimbangkan ke akun Modal/Ekuitas (3999).` }}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-6 font-mono self-end md:self-center shrink-0">
            <div>
              <span class="block text-[10px] text-zinc-500 font-sans">TOTAL DEBET</span>
              <span class="text-sm font-bold text-emerald-700">{{ formatIDR(totalDebit) }}</span>
            </div>
            <div>
              <span class="block text-[10px] text-zinc-500 font-sans">TOTAL KREDIT</span>
              <span class="text-sm font-bold text-rose-700">{{ formatIDR(totalCredit) }}</span>
            </div>
            <div>
              <span class="block text-[10px] text-zinc-500 font-sans">SELISIH</span>
              <span :class="['text-sm font-extrabold', isBalanced ? 'text-emerald-700' : 'text-amber-700']">
                {{ formatIDR(Math.abs(difference)) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Section C: Tabel Posisi Akun Saldo Awal -->
        <div>
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-3">
            <div>
              <h3 class="text-sm font-bold text-zinc-900">
                Daftar Posisi Akun Saldo Awal ({{ rows.length }} Baris Akuntansi)
              </h3>
              <span class="text-xs text-zinc-400">
                Masukkan nilai saldo debet dan kredit akun buku besar. Baris terintegrasi langsung dengan mutasi kasir dan jurnal umum klinik.
              </span>
            </div>

            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="loadAllLeafAccounts"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
                title="Tampilkan seluruh daftar bagan akun (COA) aktif"
              >
                <ListPlus class="w-3.5 h-3.5 text-zinc-500" />
                <span>Tampilkan Seluruh COA</span>
              </button>

              <button
                type="button"
                @click="addRow"
                class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
              >
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Baris</span>
              </button>
            </div>
          </div>

          <!-- Clean Table Container -->
          <div class="border border-zinc-200/90 rounded-2xl overflow-hidden bg-white shadow-2xs">
            <div class="overflow-x-auto">
              <table class="w-full text-xs text-left">
                <thead class="bg-zinc-50/90 text-zinc-500 font-bold border-b border-zinc-200 text-[11px] uppercase tracking-wider">
                  <tr>
                    <th class="px-3 py-3 w-10 text-center text-zinc-400">#</th>
                    <th class="px-4 py-3 min-w-[320px]">Akun Buku Besar (COA BCS)</th>
                    <th class="px-3 py-3 w-24 text-center">Unit</th>
                    <th class="px-4 py-3 min-w-[200px]">Uraian Baris (Line Remark)</th>
                    <th class="px-4 py-3 w-48 text-right whitespace-nowrap">Saldo Debet (Rp)</th>
                    <th class="px-4 py-3 w-48 text-right whitespace-nowrap">Saldo Kredit (Rp)</th>
                    <th class="px-2 py-3 w-10 text-center"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-zinc-100 text-zinc-800">
                  <tr v-for="(row, idx) in rows" :key="idx" class="hover:bg-zinc-50/80 transition">
                    <td class="px-3 py-2 text-center text-zinc-400 font-mono">{{ idx + 1 }}</td>
                    <td class="px-4 py-2">
                      <SearchableSelect
                        :model-value="row.code"
                        :options="availableAccountOptions"
                        placeholder="-- Pilih Akun COA --"
                        @update:model-value="(code: any) => onAccountSelect(idx, String(code))"
                      />
                    </td>
                    <td class="px-3 py-2 text-center">
                      <select
                        v-model="row.unit"
                        class="px-2 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs font-bold text-zinc-800 focus:outline-none cursor-pointer"
                      >
                        <option value="ALL">ALL</option>
                        <option value="KIDINEST">Kidinest</option>
                        <option value="SSC">SSC</option>
                      </select>
                    </td>
                    <td class="px-4 py-2">
                      <input
                        v-model="row.notes"
                        type="text"
                        placeholder="Keterangan saldo baris..."
                        class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white"
                      />
                    </td>
                    <td class="px-4 py-2 text-right">
                      <CurrencyInput
                        v-model="row.debit"
                        prefix="Rp"
                        input-class="w-full pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                        placeholder="0"
                      />
                    </td>
                    <td class="px-4 py-2 text-right">
                      <CurrencyInput
                        v-model="row.credit"
                        prefix="Rp"
                        input-class="w-full pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                        placeholder="0"
                      />
                    </td>
                    <td class="px-2 py-2 text-center">
                      <button
                        type="button"
                        @click="removeRow(idx)"
                        class="p-1 text-zinc-400 hover:text-rose-600 transition cursor-pointer rounded"
                        title="Hapus baris"
                      >
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    </td>
                  </tr>

                  <tr v-if="rows.length === 0">
                    <td colspan="7" class="px-4 py-12 text-center text-zinc-400">
                      <FileText class="w-8 h-8 mx-auto mb-2 opacity-40" />
                      <p class="font-medium text-xs">Belum ada akun saldo awal yang dimasukkan</p>
                      <p class="text-[11px] text-zinc-400 mt-0.5">Klik "Ambil Saldo Kasir & Jurnal" untuk menarik data mutasi kasir atau "Tambah Baris" untuk entri manual.</p>
                    </td>
                  </tr>
                </tbody>

                <tfoot class="bg-zinc-50/90 border-t border-zinc-200 font-bold text-xs">
                  <tr>
                    <td colspan="4" class="px-4 py-3 text-right uppercase text-zinc-500 tracking-wider">
                      Total Saldo Awal:
                    </td>
                    <td class="px-4 py-3 text-right font-mono text-emerald-800 text-sm">
                      {{ formatIDR(totalDebit) }}
                    </td>
                    <td class="px-4 py-3 text-right font-mono text-rose-800 text-sm">
                      {{ formatIDR(totalCredit) }}
                    </td>
                    <td></td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>
        </div>

        <!-- Section D: Ringkasan Nilai Neraca Saldo Awal -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3">
            Ringkasan Persamaan Dasar Neraca Saldo Awal
          </h3>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs font-mono">
            <div class="p-4 bg-zinc-50/80 border border-zinc-200 rounded-2xl space-y-1">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block">TOTAL ASET (AKTIVA BERSIH)</span>
              <span class="text-base font-extrabold text-emerald-700 block">{{ formatIDR(totalAsset) }}</span>
              <span class="text-[10px] font-sans text-zinc-400 block">Kas, Bank, Piutang Pasien, Persediaan Obat & Alkes</span>
            </div>

            <div class="p-4 bg-zinc-50/80 border border-zinc-200 rounded-2xl space-y-1">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block">TOTAL KEWAJIBAN (LIABILITAS)</span>
              <span class="text-base font-extrabold text-rose-700 block">{{ formatIDR(totalLiability) }}</span>
              <span class="text-[10px] font-sans text-zinc-400 block">Hutang Obat PBF, Deposit Pasien & Hutang Pajak</span>
            </div>

            <div class="p-4 bg-zinc-50/80 border border-zinc-200 rounded-2xl space-y-1">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block">TOTAL MODAL (EKUITAS BERSIH)</span>
              <span class="text-base font-extrabold text-teal-700 block">{{ formatIDR(totalEquity) }}</span>
              <span class="text-[10px] font-sans text-zinc-400 block">Modal Disetor Pemilik Klinik & Laba Ditahan</span>
            </div>
          </div>
        </div>

        <!-- Bottom Action Bar -->
        <div class="pt-4 border-t border-zinc-200 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs">
          <div class="flex items-center gap-4">
            <router-link
              to="/accounting/opening-balance"
              class="text-zinc-500 hover:text-zinc-800 font-semibold transition flex items-center gap-1.5 cursor-pointer"
            >
              <ArrowLeft class="w-3.5 h-3.5" />
              <span>Kembali ke Daftar Saldo Awal</span>
            </router-link>

            <button
              type="button"
              @click="resetInputs"
              class="text-zinc-500 hover:text-rose-600 font-semibold transition flex items-center gap-1.5 cursor-pointer"
            >
              <RotateCcw class="w-3.5 h-3.5" />
              <span>Kosongkan Seluruh Baris</span>
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              v-if="!isBalanced && difference !== 0"
              type="button"
              @click="autoBalanceToEquity"
              class="inline-flex items-center gap-1.5 px-3 py-2 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold rounded-xl border border-amber-200 transition active:scale-95 cursor-pointer"
            >
              <Scale class="w-3.5 h-3.5 text-amber-600" />
              <span>Seimbangkan Selisih</span>
            </button>

            <button
              type="button"
              :disabled="isSubmitting || !isBalanced"
              @click="handleSaveOpeningBalance(false)"
              class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 disabled:opacity-50 text-zinc-800 font-bold rounded-xl border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
            >
              <Save class="w-3.5 h-3.5 text-zinc-600" />
              <span>Simpan Draft</span>
            </button>

            <button
              type="button"
              :disabled="isSubmitting || !isBalanced"
              @click="handleSaveOpeningBalance(true)"
              class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white font-bold rounded-xl shadow-2xs transition active:scale-95 cursor-pointer"
            >
              <Send class="w-3.5 h-3.5 text-emerald-400" />
              <span>Posting ke Buku Besar</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 4. TAB 2: RECONCILIATION & STATUS PERSAMAAN NERACA -->
      <div v-else-if="activeTab === 'RECONCILIATION'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Status Rekonsiliasi Neraca Saldo Awal</h3>
          <p class="text-xs text-zinc-500">Validasi kepatuhan Standar Akuntansi Keuangan (PSAK): Persamaan Dasar Akuntansi (Aset = Kewajiban + Ekuitas).</p>
        </div>

        <div class="space-y-3 text-xs font-mono">
          <div class="flex items-center justify-between p-4 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div>
              <span class="font-sans font-bold text-zinc-900 block">TOTAL ASET / AKTIVA (DEBET)</span>
              <span class="font-sans text-[11px] text-zinc-500">Total sumber daya ekonomi yang dimiliki klinik per cut-off</span>
            </div>
            <span class="text-base font-bold text-emerald-700">{{ formatIDR(totalAsset) }}</span>
          </div>

          <div class="flex items-center justify-between p-4 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div>
              <span class="font-sans font-bold text-zinc-900 block">TOTAL KEWAJIBAN & MODAL (KREDIT)</span>
              <span class="font-sans text-[11px] text-zinc-500">Jumlah klaim kewajiban pihak ketiga dan ekuitas pemilik klinik</span>
            </div>
            <span class="text-base font-bold text-zinc-900">{{ formatIDR(totalLiability + totalEquity) }}</span>
          </div>

          <div
            :class="[
              'flex items-center justify-between p-4 rounded-xl border',
              isBalanced ? 'bg-emerald-50 border-emerald-200 text-emerald-950' : 'bg-amber-50 border-amber-200 text-amber-950'
            ]"
          >
            <div class="flex items-center gap-2">
              <ShieldCheck v-if="isBalanced" class="w-4 h-4 text-emerald-600" />
              <AlertTriangle v-else class="w-4 h-4 text-amber-600" />
              <span class="font-sans font-bold">SELISIH PERSAMAAN DASAR NERACA</span>
            </div>
            <span class="text-base font-black">
              {{ isBalanced ? 'Rp 0 (SEIMBANG 100%)' : formatIDR(Math.abs(difference)) }}
            </span>
          </div>
        </div>
      </div>

      <!-- 5. TAB 3: AUDIT TRAIL & RIWAYAT POSTING SALDO AWAL -->
      <div v-else class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Riwayat Jurnal Saldo Awal Terdaftar</h3>
          <p class="text-xs text-zinc-500">Daftar transaksi pembukuan saldo awal resmi yang telah disimpan atau diposting di sistem.</p>
        </div>

        <div v-if="openingJournals.length > 0" class="border border-zinc-200 rounded-2xl overflow-hidden bg-white shadow-2xs">
          <table class="w-full text-xs text-left">
            <thead class="bg-zinc-50 text-zinc-500 font-bold border-b border-zinc-200 text-[11px] uppercase">
              <tr>
                <th class="px-4 py-3">No. Jurnal</th>
                <th class="px-4 py-3">Tanggal Cut-Off</th>
                <th class="px-4 py-3">Keterangan</th>
                <th class="px-4 py-3 text-right">Total Saldo</th>
                <th class="px-4 py-3 text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100 text-zinc-800 font-mono">
              <tr v-for="j in openingJournals" :key="j.id" class="hover:bg-zinc-50 transition">
                <td class="px-4 py-3 font-bold text-blue-600">{{ j.journalNumber }}</td>
                <td class="px-4 py-3 text-zinc-600">{{ formatDate(j.journalDate) }}</td>
                <td class="px-4 py-3 font-sans text-zinc-800">{{ j.description }}</td>
                <td class="px-4 py-3 text-right font-bold text-emerald-700">{{ formatIDR(j.totalDebit) }}</td>
                <td class="px-4 py-3 text-center">
                  <span
                    :class="[
                      'px-2 py-0.5 rounded text-[10px] font-bold uppercase',
                      j.status === 'POSTED' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-amber-50 text-amber-700 border border-amber-200'
                    ]"
                  >
                    {{ j.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-else class="text-center py-12 bg-zinc-50 border border-zinc-200 rounded-2xl text-zinc-400">
          <Layers class="w-8 h-8 mx-auto mb-2 opacity-40" />
          <p class="font-medium text-xs">Belum ada riwayat posting saldo awal</p>
          <p class="text-[11px] text-zinc-400 mt-0.5">Setelah saldo awal diposting, riwayat jurnal voucher akan otomatis tercatat di sini.</p>
        </div>
      </div>
    </div>
  </div>
</template>
