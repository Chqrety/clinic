<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useCashBankStore } from '@/stores/cashbank'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import type { CashBankTransaction, BusinessUnit } from '@/types/accounting'
import {
  Wallet,
  Landmark,
  Plus,
  ArrowDownLeft,
  ArrowUpRight,
  ArrowRightLeft,
  Search,
  Printer,
  X,
  CheckCircle2,
  FolderPlus,
  Calendar,
  Building2,
  FileText,
  Database,
  BookmarkCheck
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()
const cashBankStore = useCashBankStore()

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    cashBankStore.fetchCashBankData(),
    journalsStore.fetchJournals()
  ])
})

const activeTab = ref<'CASH' | 'BANK'>('CASH')
const selectedAccountCode = ref<string>('1101.001')
const searchQuery = ref('')
const selectedTypeFilter = ref<'ALL' | 'CASH_IN' | 'CASH_OUT' | 'TRANSFER'>('ALL')

// Modals
const showAddAccountModal = ref(false)
const showCashInModal = ref(false)
const showCashOutModal = ref(false)
const showTransferModal = ref(false)
const showVoucherModal = ref(false)
const activeVoucher = ref<CashBankTransaction | null>(null)

// Forms
const newAccountForm = ref({
  category: 'CASH' as 'CASH' | 'BANK',
  code: '1101.004',
  name: '',
  businessUnit: 'ALL' as BusinessUnit,
  openingBalance: 0,
  description: ''
})

const cashInForm = ref({
  accountCode: '1101.001',
  counterpartAccountCode: '4101',
  amount: 0,
  transactionDate: new Date().toISOString().slice(0, 10),
  recipientOrPayer: '',
  description: '',
  referenceNo: ''
})

const cashOutForm = ref({
  accountCode: '1101.003',
  counterpartAccountCode: '5202',
  amount: 0,
  transactionDate: new Date().toISOString().slice(0, 10),
  recipientOrPayer: '',
  description: '',
  referenceNo: ''
})

const transferForm = ref({
  sourceAccountCode: '1101.001',
  targetAccountCode: '1102.001',
  amount: 0,
  adminFee: 0,
  transactionDate: new Date().toISOString().slice(0, 10),
  description: 'Setoran fisik uang kasir ke rekening Bank BCA',
  referenceNo: ''
})

// Relevant accounts per active tab
const relevantAccounts = computed(() => {
  return activeTab.value === 'CASH' ? cashBankStore.cashAccounts : cashBankStore.bankAccounts
})

// Current account details
const selectedAccount = computed(() => {
  return ledgerStore.accounts.find(a => a.code === selectedAccountCode.value)
})

// Ledger data (mutations & running balance)
const ledgerData = computed(() => {
  return ledgerStore.getAccountLedger(selectedAccountCode.value)
})

// Total Cash
const totalCash = computed(() => {
  return cashBankStore.cashAccounts.reduce((sum, acc) => {
    return sum + (ledgerStore.getAccountLedger(acc.code).endingBalance || 0)
  }, 0)
})

// Total Bank
const totalBank = computed(() => {
  return cashBankStore.bankAccounts.reduce((sum, acc) => {
    return sum + (ledgerStore.getAccountLedger(acc.code).endingBalance || 0)
  }, 0)
})

// Total Debits & Credits for selected account
const totalDebitMutation = computed(() => {
  return ledgerData.value.entries.reduce((sum, e) => sum + e.debit, 0)
})

const totalCreditMutation = computed(() => {
  return ledgerData.value.entries.reduce((sum, e) => sum + e.credit, 0)
})

// Accounts for counterpart selection in forms (non-header accounts)
const counterpartAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader)
})

// Auto switch selected account when switching tab
function setTab(tab: 'CASH' | 'BANK') {
  activeTab.value = tab
  const list = tab === 'CASH' ? cashBankStore.cashAccounts : cashBankStore.bankAccounts
  if (list.length > 0) {
    selectedAccountCode.value = list[0].code
  }
}

// Watch new account category to generate recommended code
watch(
  () => newAccountForm.value.category,
  newCat => {
    if (newCat === 'CASH') {
      const nextIdx = cashBankStore.cashAccounts.length + 1
      newAccountForm.value.code = `1101.${String(nextIdx).padStart(3, '0')}`
    } else {
      const nextIdx = cashBankStore.bankAccounts.length + 1
      newAccountForm.value.code = `1102.${String(nextIdx).padStart(3, '0')}`
    }
  }
)

// Combined transactions (from generalLedgerEntries + cashBankStore metadata)
const combinedEntries = computed(() => {
  const gl = ledgerData.value.entries
  return gl.map(entry => {
    const matchedTx = cashBankStore.transactions.find(
      t => t.journalNumber === entry.journalNumber || t.transactionNumber === entry.journalNumber
    )
    return {
      id: entry.id,
      date: entry.date,
      journalNumber: entry.journalNumber,
      source: entry.source,
      description: entry.description,
      debit: entry.debit,
      credit: entry.credit,
      runningBalance: entry.runningBalance,
      transactionNumber: matchedTx?.transactionNumber || entry.journalNumber,
      type: matchedTx?.type || (entry.debit > 0 ? 'CASH_IN' : 'CASH_OUT'),
      counterpart: matchedTx?.counterpartAccountName || (entry.debit > 0 ? 'Kasir / Lawan Akun' : 'Operasional'),
      counterpartyName: matchedTx?.recipientOrPayer || (entry.debit > 0 ? 'Pasien / Kasir' : 'Vendor / Staf'),
      rawTx: matchedTx
    }
  })
})

// Filtered entries
const filteredEntries = computed(() => {
  return combinedEntries.value.filter(item => {
    if (selectedTypeFilter.value !== 'ALL' && item.type !== selectedTypeFilter.value) {
      return false
    }
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase()
      const matchDesc = item.description.toLowerCase().includes(q)
      const matchNo = item.journalNumber.toLowerCase().includes(q) || item.transactionNumber.toLowerCase().includes(q)
      const matchCp = item.counterpartyName.toLowerCase().includes(q)
      if (!matchDesc && !matchNo && !matchCp) return false
    }
    return true
  })
})

// Open Print Voucher modal
function openVoucher(entry: any) {
  if (entry.rawTx) {
    activeVoucher.value = entry.rawTx
  } else {
    // Construct virtual voucher from entry
    activeVoucher.value = {
      id: entry.id,
      transactionNumber: entry.transactionNumber,
      date: entry.date,
      type: entry.debit > 0 ? 'CASH_IN' : 'CASH_OUT',
      accountCode: selectedAccountCode.value,
      accountName: selectedAccount.value?.name || selectedAccountCode.value,
      counterpartAccountCode: '-',
      counterpartAccountName: entry.counterpart,
      amount: entry.debit > 0 ? entry.debit : entry.credit,
      description: entry.description,
      recipientOrPayer: entry.counterpartyName,
      businessUnit: 'ALL',
      journalNumber: entry.journalNumber,
      createdAt: entry.date,
      createdBy: 'Rian Pratama, S.Ak.'
    }
  }
  showVoucherModal.value = true
}

// Handlers for Submitting Forms
function handleSaveNewAccount() {
  if (!newAccountForm.value.code || !newAccountForm.value.name) {
    accountingStore.showToast('Harap lengkapi kode dan nama buku kas/bank', 'error')
    return
  }
  cashBankStore.createCashBankAccount({
    category: newAccountForm.value.category,
    code: newAccountForm.value.code,
    name: newAccountForm.value.name,
    businessUnit: newAccountForm.value.businessUnit,
    openingBalance: Number(newAccountForm.value.openingBalance) || 0,
    description: newAccountForm.value.description
  })
  selectedAccountCode.value = newAccountForm.value.code
  showAddAccountModal.value = false
  // Reset form
  newAccountForm.value.name = ''
  newAccountForm.value.openingBalance = 0
  newAccountForm.value.description = ''
}

function handleSaveCashIn() {
  if (cashInForm.value.amount <= 0) {
    accountingStore.showToast('Jumlah nominal penerimaan harus lebih dari Rp 0', 'error')
    return
  }
  cashBankStore.recordCashIn({
    accountCode: cashInForm.value.accountCode,
    counterpartAccountCode: cashInForm.value.counterpartAccountCode,
    amount: Number(cashInForm.value.amount),
    transactionDate: cashInForm.value.transactionDate,
    recipientOrPayer: cashInForm.value.recipientOrPayer || 'Pasien / Mitra',
    description: cashInForm.value.description || 'Penerimaan Kas/Bank Operasional',
    referenceNo: cashInForm.value.referenceNo
  })
  selectedAccountCode.value = cashInForm.value.accountCode
  showCashInModal.value = false
  // Reset form
  cashInForm.value.amount = 0
  cashInForm.value.description = ''
  cashInForm.value.recipientOrPayer = ''
  cashInForm.value.referenceNo = ''
}

function handleSaveCashOut() {
  if (cashOutForm.value.amount <= 0) {
    accountingStore.showToast('Jumlah nominal pengeluaran harus lebih dari Rp 0', 'error')
    return
  }
  cashBankStore.recordCashOut({
    accountCode: cashOutForm.value.accountCode,
    counterpartAccountCode: cashOutForm.value.counterpartAccountCode,
    amount: Number(cashOutForm.value.amount),
    transactionDate: cashOutForm.value.transactionDate,
    recipientOrPayer: cashOutForm.value.recipientOrPayer || 'Vendor / Penerima',
    description: cashOutForm.value.description || 'Pengeluaran Kas/Bank Operasional',
    referenceNo: cashOutForm.value.referenceNo
  })
  selectedAccountCode.value = cashOutForm.value.accountCode
  showCashOutModal.value = false
  // Reset form
  cashOutForm.value.amount = 0
  cashOutForm.value.description = ''
  cashOutForm.value.recipientOrPayer = ''
  cashOutForm.value.referenceNo = ''
}

function handleSaveTransfer() {
  if (transferForm.value.amount <= 0) {
    accountingStore.showToast('Jumlah nominal transfer harus lebih dari Rp 0', 'error')
    return
  }
  if (transferForm.value.sourceAccountCode === transferForm.value.targetAccountCode) {
    accountingStore.showToast('Rekening sumber dan rekening tujuan tidak boleh sama', 'error')
    return
  }
  cashBankStore.recordTransfer({
    sourceAccountCode: transferForm.value.sourceAccountCode,
    targetAccountCode: transferForm.value.targetAccountCode,
    amount: Number(transferForm.value.amount),
    adminFee: Number(transferForm.value.adminFee) || 0,
    transactionDate: transferForm.value.transactionDate,
    description: transferForm.value.description || 'Mutasi transfer antar rekening kas dan bank',
    referenceNo: transferForm.value.referenceNo
  })
  selectedAccountCode.value = transferForm.value.sourceAccountCode
  showTransferModal.value = false
  // Reset form
  transferForm.value.amount = 0
  transferForm.value.adminFee = 0
  transferForm.value.referenceNo = ''
}

// Terbilang Rupiah Helper
function terbilang(n: number): string {
  const bilangan = ['', 'Satu', 'Dua', 'Tiga', 'Empat', 'Lima', 'Enam', 'Tujuh', 'Delapan', 'Sembilan', 'Sepuluh', 'Sebelas']
  if (n < 12) return bilangan[n]
  if (n < 20) return terbilang(n - 10) + ' Belas'
  if (n < 100) return terbilang(Math.floor(n / 10)) + ' Puluh ' + terbilang(n % 10)
  if (n < 200) return 'Seratus ' + terbilang(n - 100)
  if (n < 1000) return terbilang(Math.floor(n / 100)) + ' Ratus ' + terbilang(n % 100)
  if (n < 2000) return 'Seribu ' + terbilang(n - 1000)
  if (n < 1000000) return terbilang(Math.floor(n / 1000)) + ' Ribu ' + terbilang(n % 1000)
  if (n < 1000000000) return terbilang(Math.floor(n / 1000000)) + ' Juta ' + terbilang(n % 1000000)
  return terbilang(Math.floor(n / 1000000000)) + ' Miliar ' + terbilang(n % 1000000000)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Treasury & Cash Management
          </span>
          <span class="text-xs text-slate-400">&bull; Standar Akuntansi Klinik BCS</span>
        </div>
        <h1 class="text-2xl font-black text-slate-900 mt-1">
          Buku Kas & Rekening Bank
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Pencatatan langsung transaksi kas & bank operasional klinik, penerimaan (BKM), pengeluaran (BKK), dan transfer dana tanpa harus menunggu closing kasir / apotek.
        </p>
      </div>

      <!-- Quick Action Buttons -->
      <div class="flex flex-wrap items-center gap-2">
        <!-- Atur Saldo Awal -->
        <button
          @click="router.push('/opening-balance')"
          class="flex items-center gap-1.5 px-3 py-2 bg-purple-700 hover:bg-purple-800 text-white rounded-xl text-xs font-bold transition shadow-xs"
        >
          <BookmarkCheck class="w-3.5 h-3.5" />
          <span>Atur Saldo Awal</span>
        </button>

        <!-- Buat Buku Kas / Rekening Baru -->
        <button
          @click="showAddAccountModal = true"
          class="flex items-center gap-1.5 px-3 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-xl text-xs font-bold transition shadow-xs"
        >
          <FolderPlus class="w-3.5 h-3.5" />
          <span>+ Buat Buku Kas / Bank</span>
        </button>

        <!-- Catat Kas Masuk (BKM) -->
        <button
          @click="showCashInModal = true"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-xl text-xs font-bold transition shadow-xs"
        >
          <ArrowDownLeft class="w-3.5 h-3.5" />
          <span>+ Kas/Bank Masuk (BKM)</span>
        </button>

        <!-- Catat Kas Keluar (BKK) -->
        <button
          @click="showCashOutModal = true"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-xl text-xs font-bold transition shadow-xs"
        >
          <ArrowUpRight class="w-3.5 h-3.5" />
          <span>- Kas/Bank Keluar (BKK)</span>
        </button>

        <!-- Transfer Antar Kas & Bank -->
        <button
          @click="showTransferModal = true"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold transition shadow-xs"
        >
          <ArrowRightLeft class="w-3.5 h-3.5" />
          <span>⇄ Transfer Kas / Bank</span>
        </button>
      </div>
    </div>

    <!-- 2 Overview Cards: Kas Fisik vs Rekening Bank -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Kas Card -->
      <div
        @click="setTab('CASH')"
        :class="[
          'p-5 rounded-2xl border cursor-pointer transition-all space-y-3',
          activeTab === 'CASH'
            ? 'bg-slate-900 text-white border-slate-900 shadow-md ring-2 ring-blue-500/20'
            : 'bg-white text-slate-900 border-slate-200/80 hover:border-slate-300'
        ]"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div :class="['w-10 h-10 rounded-xl flex items-center justify-center', activeTab === 'CASH' ? 'bg-blue-600 text-white' : 'bg-blue-50 text-blue-600']">
              <Wallet class="w-5 h-5" />
            </div>
            <div>
              <h3 class="font-bold text-sm">Buku Kas Operasional & Kasir</h3>
              <p :class="['text-[11px]', activeTab === 'CASH' ? 'text-slate-300' : 'text-slate-500']">
                Kasir Fisik Kidinest, SSC, dan Kas Kecil (Petty Cash)
              </p>
            </div>
          </div>
          <div class="text-right">
            <span class="text-[10px] uppercase font-bold block opacity-70">Total Kas:</span>
            <span class="font-mono text-lg font-bold text-emerald-400">
              {{ formatIDR(totalCash) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Bank Card -->
      <div
        @click="setTab('BANK')"
        :class="[
          'p-5 rounded-2xl border cursor-pointer transition-all space-y-3',
          activeTab === 'BANK'
            ? 'bg-slate-900 text-white border-slate-900 shadow-md ring-2 ring-blue-500/20'
            : 'bg-white text-slate-900 border-slate-200/80 hover:border-slate-300'
        ]"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div :class="['w-10 h-10 rounded-xl flex items-center justify-center', activeTab === 'BANK' ? 'bg-indigo-600 text-white' : 'bg-indigo-50 text-indigo-600']">
              <Landmark class="w-5 h-5" />
            </div>
            <div>
              <h3 class="font-bold text-sm">Buku Rekening Bank Operasional</h3>
              <p :class="['text-[11px]', activeTab === 'BANK' ? 'text-slate-300' : 'text-slate-500']">
                Bank BCA Giro, Mandiri Payroll, dan Settlement Transit
              </p>
            </div>
          </div>
          <div class="text-right">
            <span class="text-[10px] uppercase font-bold block opacity-70">Total Bank:</span>
            <span class="font-mono text-lg font-bold text-emerald-400">
              {{ formatIDR(totalBank) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Sub-Account Selector Pills -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs space-y-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-slate-700">Pilih Buku Kas / Rekening Aktif:</span>
          <span class="text-[11px] text-slate-500">({{ relevantAccounts.length }} rekening terdaftar)</span>
        </div>

        <div class="text-xs font-mono font-bold text-slate-800 bg-slate-50 px-3 py-1.5 rounded-lg border border-slate-200">
          Saldo Buku: <span class="text-blue-700 text-sm">{{ formatIDR(ledgerData.endingBalance) }}</span>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-2 pt-1 border-t border-slate-100">
        <button
          v-for="acc in relevantAccounts"
          :key="acc.code"
          @click="selectedAccountCode = acc.code"
          :class="[
            'px-3.5 py-2 rounded-xl text-xs font-bold transition-all flex items-center gap-2',
            selectedAccountCode === acc.code
              ? 'bg-blue-600 text-white shadow-xs'
              : 'bg-slate-100 text-slate-700 hover:bg-slate-200'
          ]"
        >
          <span class="font-mono text-[11px] opacity-80">{{ acc.code }}</span>
          <span>{{ acc.name }}</span>
          <span :class="['px-1.5 py-0.5 rounded text-[10px] font-mono', selectedAccountCode === acc.code ? 'bg-blue-700 text-blue-100' : 'bg-slate-200 text-slate-600']">
            {{ formatIDR(ledgerStore.getAccountLedger(acc.code).endingBalance, false) }}
          </span>
        </button>
      </div>
    </div>

    <!-- Active Account KPI Summary Bar -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs text-xs">
      <div class="p-3 bg-slate-50 rounded-lg border border-slate-200">
        <span class="text-[10px] uppercase font-bold text-slate-500 block">Saldo Awal Periode</span>
        <span class="font-mono font-bold text-sm text-slate-900 block mt-0.5">
          {{ formatIDR(ledgerData.openingBalance) }}
        </span>
      </div>

      <div class="p-3 bg-emerald-50 rounded-lg border border-emerald-200 text-emerald-950">
        <span class="text-[10px] uppercase font-bold text-emerald-700 block">Total Masuk (Debit)</span>
        <span class="font-mono font-bold text-sm text-emerald-700 block mt-0.5">
          + {{ formatIDR(totalDebitMutation) }}
        </span>
      </div>

      <div class="p-3 bg-rose-50 rounded-lg border border-rose-200 text-rose-950">
        <span class="text-[10px] uppercase font-bold text-rose-700 block">Total Keluar (Kredit)</span>
        <span class="font-mono font-bold text-sm text-rose-700 block mt-0.5">
          - {{ formatIDR(totalCreditMutation) }}
        </span>
      </div>

      <div class="p-3 bg-blue-50 rounded-lg border border-blue-200 text-blue-950">
        <span class="text-[10px] uppercase font-bold text-blue-700 block">Saldo Akhir Berjalan</span>
        <span class="font-mono font-bold text-sm text-blue-700 block mt-0.5">
          {{ formatIDR(ledgerData.endingBalance) }}
        </span>
      </div>
    </div>

    <!-- Filter and Search Bar -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 bg-white p-3 rounded-xl border border-slate-200/80 shadow-2xs text-xs">
      <div class="flex items-center gap-2">
        <span class="font-bold text-slate-600">Filter Mutasi:</span>
        <div class="flex items-center bg-slate-100 p-0.5 rounded-lg border border-slate-200">
          <button
            @click="selectedTypeFilter = 'ALL'"
            :class="['px-3 py-1 rounded-md font-bold transition text-xs', selectedTypeFilter === 'ALL' ? 'bg-white shadow-xs text-blue-600' : 'text-slate-600 hover:text-slate-900']"
          >
            Semua ({{ combinedEntries.length }})
          </button>
          <button
            @click="selectedTypeFilter = 'CASH_IN'"
            :class="['px-3 py-1 rounded-md font-bold transition text-xs', selectedTypeFilter === 'CASH_IN' ? 'bg-white shadow-xs text-emerald-600' : 'text-slate-600 hover:text-slate-900']"
          >
            Penerimaan / Masuk
          </button>
          <button
            @click="selectedTypeFilter = 'CASH_OUT'"
            :class="['px-3 py-1 rounded-md font-bold transition text-xs', selectedTypeFilter === 'CASH_OUT' ? 'bg-white shadow-xs text-rose-600' : 'text-slate-600 hover:text-slate-900']"
          >
            Pengeluaran / Keluar
          </button>
          <button
            @click="selectedTypeFilter = 'TRANSFER'"
            :class="['px-3 py-1 rounded-md font-bold transition text-xs', selectedTypeFilter === 'TRANSFER' ? 'bg-white shadow-xs text-blue-600' : 'text-slate-600 hover:text-slate-900']"
          >
            Transfer Antar Rekening
          </button>
        </div>
      </div>

      <div class="relative w-full md:w-72">
        <Search class="w-3.5 h-3.5 absolute left-3 top-2.5 text-slate-400" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari uraian, nomor bukti, pihak ketiga..."
          class="w-full pl-8 pr-3 py-1.5 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
        />
      </div>
    </div>

    <!-- Mutations Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">No. Bukti / Jurnal</th>
              <th class="px-4 py-3">Tipe</th>
              <th class="px-4 py-3">Pihak Terkait / Sumber</th>
              <th class="px-4 py-3">Uraian Transaksi</th>
              <th class="px-4 py-3 text-right text-emerald-700">Penerimaan (Dr)</th>
              <th class="px-4 py-3 text-right text-rose-700">Pengeluaran (Cr)</th>
              <th class="px-4 py-3 text-right">Saldo Kas/Bank</th>
              <th class="px-4 py-3 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <!-- Opening Balance Row -->
            <tr class="bg-slate-50/40 font-semibold text-slate-600">
              <td class="px-4 py-3 font-mono">01/03/2026</td>
              <td class="px-4 py-3 font-mono text-slate-400">-</td>
              <td class="px-4 py-3">
                <span class="px-2 py-0.5 bg-slate-200 text-slate-700 rounded text-[10px] font-bold">SALDO AWAL</span>
              </td>
              <td colspan="2" class="px-4 py-3 italic text-slate-500">Saldo Awal Buku Kas / Rekening Bank BCS</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono font-bold text-slate-900 bg-slate-50/40">
                {{ formatIDR(ledgerData.openingBalance) }}
              </td>
              <td class="px-4 py-3 text-center text-slate-400">-</td>
            </tr>

            <!-- Transactions Rows -->
            <tr
              v-for="entry in filteredEntries"
              :key="entry.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap font-mono text-[11px]">
                {{ formatDate(entry.date) }}
              </td>
              <td class="px-4 py-3.5 whitespace-nowrap">
                <span class="font-mono font-bold text-blue-700 block text-[11px]">
                  {{ entry.transactionNumber }}
                </span>
                <span class="font-mono text-[10px] text-slate-400 block">
                  {{ entry.journalNumber }}
                </span>
              </td>
              <td class="px-4 py-3.5 whitespace-nowrap">
                <span
                  v-if="entry.type === 'CASH_IN'"
                  class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-100 text-emerald-800 border border-emerald-200"
                >
                  Kas Masuk
                </span>
                <span
                  v-else-if="entry.type === 'CASH_OUT'"
                  class="px-2 py-0.5 rounded text-[10px] font-bold bg-rose-100 text-rose-800 border border-rose-200"
                >
                  Kas Keluar
                </span>
                <span
                  v-else
                  class="px-2 py-0.5 rounded text-[10px] font-bold bg-blue-100 text-blue-800 border border-blue-200"
                >
                  Transfer Mutasi
                </span>
              </td>
              <td class="px-4 py-3.5 text-slate-700 whitespace-nowrap font-medium max-w-xs truncate">
                {{ entry.counterpartyName }}
              </td>
              <td class="px-4 py-3.5 max-w-md text-slate-800">
                <p class="font-medium text-slate-900">{{ entry.description }}</p>
                <p class="text-[10px] text-slate-400 font-mono mt-0.5">Kontra: {{ entry.counterpart }}</p>
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-emerald-700 whitespace-nowrap">
                {{ entry.debit > 0 ? formatIDR(entry.debit, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-rose-700 whitespace-nowrap">
                {{ entry.credit > 0 ? formatIDR(entry.credit, false) : '-' }}
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-slate-900 whitespace-nowrap bg-slate-50/30">
                {{ formatIDR(entry.runningBalance) }}
              </td>
              <td class="px-4 py-3.5 text-center whitespace-nowrap">
                <button
                  @click="openVoucher(entry)"
                  title="Cetak Bukti Transaksi Kas / Bank (Voucher)"
                  class="p-1.5 text-slate-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition"
                >
                  <Printer class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>

            <tr v-if="filteredEntries.length === 0">
              <td colspan="9" class="px-4 py-8 text-center text-slate-400 italic">
                Belum ada transaksi kas/bank yang sesuai dengan filter pencarian.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Official Report Signatures (Gambar 1 Standard) -->
    <ReportSignatures />

    <!-- MODAL 1: Buat Buku Kas / Rekening Bank Baru -->
    <div
      v-if="showAddAccountModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-md w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-slate-900 text-white flex items-center justify-between">
          <div class="flex items-center gap-2">
            <FolderPlus class="w-4 h-4 text-blue-400" />
            <h3 class="text-sm font-bold text-white">Buat Buku Kas / Rekening Bank Baru</h3>
          </div>
          <button @click="showAddAccountModal = false" class="text-slate-400 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveNewAccount" class="p-6 space-y-4 text-xs">
          <div>
            <label class="block font-bold text-slate-700 mb-1">Kategori Buku Kas / Bank</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                type="button"
                @click="newAccountForm.category = 'CASH'"
                :class="[
                  'py-2 px-3 rounded-lg border font-bold flex items-center justify-center gap-1.5 transition',
                  newAccountForm.category === 'CASH'
                    ? 'bg-blue-50 border-blue-600 text-blue-700'
                    : 'bg-white border-slate-200 text-slate-600'
                ]"
              >
                <Wallet class="w-3.5 h-3.5" />
                <span>Buku Kas Fisik</span>
              </button>
              <button
                type="button"
                @click="newAccountForm.category = 'BANK'"
                :class="[
                  'py-2 px-3 rounded-lg border font-bold flex items-center justify-center gap-1.5 transition',
                  newAccountForm.category === 'BANK'
                    ? 'bg-blue-50 border-blue-600 text-blue-700'
                    : 'bg-white border-slate-200 text-slate-600'
                ]"
              >
                <Landmark class="w-3.5 h-3.5" />
                <span>Rekening Bank</span>
              </button>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Kode Akun COA</label>
              <input
                v-model="newAccountForm.code"
                type="text"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono font-bold"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Unit Bisnis</label>
              <select
                v-model="newAccountForm.businessUnit"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Kidinest</option>
                <option value="SSC">Semarang Skin Center</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Nama Buku Kas / Rekening Bank</label>
            <input
              v-model="newAccountForm.name"
              type="text"
              placeholder="Contoh: Kas Kecil Farmasi Cabang 2 / Bank BRI Giro"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Saldo Awal Pembukaan (Rp)</label>
            <input
              v-model.number="newAccountForm.openingBalance"
              type="number"
              min="0"
              placeholder="0"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono font-bold"
            />
            <p class="text-[10px] text-slate-400 mt-1">Saldo awal akan otomatis dijurnal penyeimbang ke Temporary Opening.</p>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Deskripsi Tambahan</label>
            <input
              v-model="newAccountForm.description"
              type="text"
              placeholder="Catatan peruntukan buku kas / bank..."
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              @click="showAddAccountModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-slate-300 rounded-lg font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs"
            >
              Simpan Buku Kas/Bank
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 2: Catat Kas/Bank Masuk (BKM / BBM) -->
    <div
      v-if="showCashInModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-emerald-700 text-white flex items-center justify-between">
          <div class="flex items-center gap-2">
            <ArrowDownLeft class="w-4 h-4 text-emerald-200" />
            <h3 class="text-sm font-bold text-white">Catat Penerimaan Kas/Bank Masuk (BKM / BBM)</h3>
          </div>
          <button @click="showCashInModal = false" class="text-emerald-200 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveCashIn" class="p-6 space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Tanggal Transaksi</label>
              <input
                v-model="cashInForm.transactionDate"
                type="date"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500 font-mono"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Rekening / Kas Penerima (Dr)</label>
              <select
                v-model="cashInForm.accountCode"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500 font-bold"
              >
                <optgroup label="Buku Kas Operasional">
                  <option v-for="acc in cashBankStore.cashAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
                <optgroup label="Rekening Bank">
                  <option v-for="acc in cashBankStore.bankAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Diterima Dari (Payer / Sumber Dana)</label>
            <input
              v-model="cashInForm.recipientOrPayer"
              type="text"
              placeholder="Contoh: Asuransi Sinarmas / Pasien Rahmawati / Bank Bunga"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
              required
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Akun Lawan Kredit (Sumber Pengakuan)</label>
            <select
              v-model="cashInForm.counterpartAccountCode"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500 font-medium"
            >
              <option v-for="acc in counterpartAccounts" :key="acc.code" :value="acc.code">
                {{ acc.code }} - {{ acc.name }} ({{ acc.type }})
              </option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Jumlah Nominal (Rp)</label>
              <input
                v-model.number="cashInForm.amount"
                type="number"
                min="1"
                placeholder="0"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500 font-mono font-bold text-emerald-700"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Nomor Referensi / Cek / Ref Bank</label>
              <input
                v-model="cashInForm.referenceNo"
                type="text"
                placeholder="Contoh: TRF-BCA-99120"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Uraian / Keterangan Transaksi</label>
            <textarea
              v-model="cashInForm.description"
              rows="2"
              placeholder="Jelaskan tujuan dan rincian penerimaan kas/bank ini..."
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
              required
            ></textarea>
          </div>

          <!-- Jurnal Preview -->
          <div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl space-y-1.5">
            <span class="text-[10px] font-bold text-emerald-800 uppercase block">Otomatisasi Jurnal Umum (POSTED):</span>
            <div class="flex justify-between font-mono text-[11px] text-emerald-900">
              <span>(Dr) {{ cashInForm.accountCode }}</span>
              <span>{{ formatIDR(cashInForm.amount || 0) }}</span>
            </div>
            <div class="flex justify-between font-mono text-[11px] text-emerald-900 pl-4">
              <span>(Cr) {{ cashInForm.counterpartAccountCode }}</span>
              <span>{{ formatIDR(cashInForm.amount || 0) }}</span>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              @click="showCashInModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-slate-300 rounded-lg font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg font-bold shadow-xs"
            >
              Simpan Penerimaan Kas (BKM)
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 3: Catat Kas/Bank Keluar (BKK / BBK / Petty Cash Voucher) -->
    <div
      v-if="showCashOutModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-rose-700 text-white flex items-center justify-between">
          <div class="flex items-center gap-2">
            <ArrowUpRight class="w-4 h-4 text-rose-200" />
            <h3 class="text-sm font-bold text-white">Catat Pengeluaran Kas/Bank Keluar (BKK / BBK)</h3>
          </div>
          <button @click="showCashOutModal = false" class="text-rose-200 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveCashOut" class="p-6 space-y-4 text-xs">
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Tanggal Transaksi</label>
              <input
                v-model="cashOutForm.transactionDate"
                type="date"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500 font-mono"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Buku Kas / Rekening Sumber (Cr)</label>
              <select
                v-model="cashOutForm.accountCode"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500 font-bold"
              >
                <optgroup label="Buku Kas Operasional (Petty Cash / Kasir)">
                  <option v-for="acc in cashBankStore.cashAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
                <optgroup label="Rekening Bank">
                  <option v-for="acc in cashBankStore.bankAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
              </select>
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Dibayarkan Kepada (Penerima / Vendor)</label>
            <input
              v-model="cashOutForm.recipientOrPayer"
              type="text"
              placeholder="Contoh: Toko ATK / Vendor Limbah Medis B3 / Staf Konsumsi"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500"
              required
            />
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Akun Lawan Beban / Hutang (Dr)</label>
            <select
              v-model="cashOutForm.counterpartAccountCode"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500 font-medium"
            >
              <option v-for="acc in counterpartAccounts" :key="acc.code" :value="acc.code">
                {{ acc.code }} - {{ acc.name }} ({{ acc.type }})
              </option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Jumlah Nominal (Rp)</label>
              <input
                v-model.number="cashOutForm.amount"
                type="number"
                min="1"
                placeholder="0"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500 font-mono font-bold text-rose-700"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Nomor Nota / Kwitansi / Ref</label>
              <input
                v-model="cashOutForm.referenceNo"
                type="text"
                placeholder="Contoh: NOTA-1102 / KW-009"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Uraian / Keperluan Pengeluaran</label>
            <textarea
              v-model="cashOutForm.description"
              rows="2"
              placeholder="Jelaskan kebutuhan pengeluaran dana kas/bank ini..."
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-rose-500"
              required
            ></textarea>
          </div>

          <!-- Jurnal Preview -->
          <div class="p-3 bg-rose-50 border border-rose-200 rounded-xl space-y-1.5">
            <span class="text-[10px] font-bold text-rose-800 uppercase block">Otomatisasi Jurnal Umum (POSTED):</span>
            <div class="flex justify-between font-mono text-[11px] text-rose-900">
              <span>(Dr) {{ cashOutForm.counterpartAccountCode }}</span>
              <span>{{ formatIDR(cashOutForm.amount || 0) }}</span>
            </div>
            <div class="flex justify-between font-mono text-[11px] text-rose-900 pl-4">
              <span>(Cr) {{ cashOutForm.accountCode }}</span>
              <span>{{ formatIDR(cashOutForm.amount || 0) }}</span>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              @click="showCashOutModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-slate-300 rounded-lg font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-rose-600 hover:bg-rose-700 text-white rounded-lg font-bold shadow-xs"
            >
              Simpan Pengeluaran Kas (BKK)
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 4: Transfer Antar Kas & Bank (Mutasi / Pemindahbukuan) -->
    <div
      v-if="showTransferModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150">
        <div class="px-6 py-4 bg-blue-600 text-white flex items-center justify-between">
          <div class="flex items-center gap-2">
            <ArrowRightLeft class="w-4 h-4 text-blue-200" />
            <h3 class="text-sm font-bold text-white">Transfer Dana Antar Rekening Kas & Bank</h3>
          </div>
          <button @click="showTransferModal = false" class="text-blue-200 hover:text-white">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleSaveTransfer" class="p-6 space-y-4 text-xs">
          <div>
            <label class="block font-bold text-slate-700 mb-1">Tanggal Transfer</label>
            <input
              v-model="transferForm.transactionDate"
              type="date"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono"
              required
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Rekening Asal / Sumber Dana (Cr)</label>
              <select
                v-model="transferForm.sourceAccountCode"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-bold"
              >
                <optgroup label="Buku Kas Operasional">
                  <option v-for="acc in cashBankStore.cashAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
                <optgroup label="Rekening Bank">
                  <option v-for="acc in cashBankStore.bankAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
              </select>
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Rekening Tujuan / Penerima (Dr)</label>
              <select
                v-model="transferForm.targetAccountCode"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-bold"
              >
                <optgroup label="Rekening Bank">
                  <option v-for="acc in cashBankStore.bankAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
                <optgroup label="Buku Kas Operasional">
                  <option v-for="acc in cashBankStore.cashAccounts" :key="acc.code" :value="acc.code">
                    {{ acc.code }} - {{ acc.name }}
                  </option>
                </optgroup>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block font-bold text-slate-700 mb-1">Nominal Transfer (Rp)</label>
              <input
                v-model.number="transferForm.amount"
                type="number"
                min="1"
                placeholder="0"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono font-bold text-blue-700"
                required
              />
            </div>

            <div>
              <label class="block font-bold text-slate-700 mb-1">Biaya Administrasi Bank (Opsional)</label>
              <input
                v-model.number="transferForm.adminFee"
                type="number"
                min="0"
                placeholder="0"
                class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 font-mono"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-slate-700 mb-1">Keterangan / Uraian Transfer</label>
            <input
              v-model="transferForm.description"
              type="text"
              placeholder="Contoh: Setor fisik penerimaan kasir shift 1 ke BCA / Top up Petty Cash"
              class="w-full px-3 py-2 bg-slate-50 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          <!-- Jurnal Preview -->
          <div class="p-3 bg-blue-50 border border-blue-200 rounded-xl space-y-1.5">
            <span class="text-[10px] font-bold text-blue-800 uppercase block">Otomatisasi Jurnal Umum (POSTED):</span>
            <div class="flex justify-between font-mono text-[11px] text-blue-900">
              <span>(Dr) {{ transferForm.targetAccountCode }}</span>
              <span>{{ formatIDR(transferForm.amount || 0) }}</span>
            </div>
            <div v-if="transferForm.adminFee > 0" class="flex justify-between font-mono text-[11px] text-blue-900">
              <span>(Dr) 7103 - BIAYA ADMINISTRASI BANK</span>
              <span>{{ formatIDR(transferForm.adminFee) }}</span>
            </div>
            <div class="flex justify-between font-mono text-[11px] text-blue-900 pl-4">
              <span>(Cr) {{ transferForm.sourceAccountCode }}</span>
              <span>{{ formatIDR((transferForm.amount || 0) + (transferForm.adminFee || 0)) }}</span>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-3 border-t border-slate-200">
            <button
              @click="showTransferModal = false"
              type="button"
              class="px-4 py-2 bg-white border border-slate-300 rounded-lg font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              type="submit"
              class="px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs"
            >
              Proses Mutasi Transfer
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL 5: Cetak Voucher Bukti Kas Resmi (Print Preview) -->
    <div
      v-if="showVoucherModal && activeVoucher"
      class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 print:p-0"
    >
      <div class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-2xl w-full overflow-hidden p-6 space-y-6">
        <!-- Voucher Header -->
        <div class="flex items-center justify-between border-b pb-4">
          <div>
            <h2 class="text-lg font-black text-slate-900 uppercase tracking-tight">PT. BRAYO CIPTA SEJAHTERA (BCS)</h2>
            <p class="text-[11px] text-slate-500">Klinik Kidinest & Semarang Skin Center (SSC)</p>
          </div>
          <div class="text-right">
            <h3 class="text-base font-black text-blue-700 uppercase">
              {{
                activeVoucher.type === 'CASH_IN'
                  ? 'BUKTI KAS/BANK MASUK (BKM)'
                  : activeVoucher.type === 'CASH_OUT'
                    ? 'BUKTI KAS/BANK KELUAR (BKK)'
                    : 'BUKTI MUTASI TRANSFER (TRF)'
              }}
            </h3>
            <p class="font-mono text-xs font-bold text-slate-700">{{ activeVoucher.transactionNumber }}</p>
          </div>
        </div>

        <!-- Voucher Details -->
        <div class="grid grid-cols-2 gap-4 text-xs">
          <div>
            <span class="text-slate-400 block text-[10px] uppercase font-bold">Tanggal Transaksi:</span>
            <span class="font-bold text-slate-800">{{ formatDate(activeVoucher.date) }}</span>
          </div>
          <div>
            <span class="text-slate-400 block text-[10px] uppercase font-bold">Buku Kas / Rekening:</span>
            <span class="font-bold text-slate-800">{{ activeVoucher.accountCode }} - {{ activeVoucher.accountName }}</span>
          </div>
          <div>
            <span class="text-slate-400 block text-[10px] uppercase font-bold">
              {{ activeVoucher.type === 'CASH_IN' ? 'Diterima Dari:' : 'Dibayarkan Kepada:' }}
            </span>
            <span class="font-bold text-slate-900 text-sm">{{ activeVoucher.recipientOrPayer || '-' }}</span>
          </div>
          <div>
            <span class="text-slate-400 block text-[10px] uppercase font-bold">Nomor Referensi:</span>
            <span class="font-mono font-bold text-slate-700">{{ activeVoucher.referenceNo || '-' }}</span>
          </div>
        </div>

        <div class="p-3 bg-slate-50 border border-slate-200 rounded-xl space-y-2">
          <div class="flex justify-between items-center text-xs">
            <span class="font-bold text-slate-600">Keperluan / Keterangan:</span>
            <span class="font-mono text-[10px] text-slate-400">Ref Jurnal: {{ activeVoucher.journalNumber }}</span>
          </div>
          <p class="text-xs text-slate-800 font-medium">{{ activeVoucher.description }}</p>
          <div class="pt-2 border-t border-slate-200 flex items-center justify-between">
            <span class="text-xs font-bold text-slate-700">Total Nominal:</span>
            <span class="font-mono font-black text-lg text-slate-900">{{ formatIDR(activeVoucher.amount) }}</span>
          </div>
          <div class="text-[11px] italic text-slate-600 bg-white p-2 rounded border border-slate-200">
            Terbilang: <span class="font-semibold">{{ terbilang(activeVoucher.amount) }} Rupiah</span>
          </div>
        </div>

        <!-- Voucher 4-Signature Columns -->
        <div class="grid grid-cols-4 gap-2 pt-4 border-t border-slate-200 text-center text-xs">
          <div class="space-y-12">
            <p class="text-[10px] font-bold text-slate-500 uppercase">Disiapkan Oleh,</p>
            <div>
              <p class="font-bold text-slate-800 underline">{{ activeVoucher.createdBy || 'Rian Pratama, S.Ak.' }}</p>
              <p class="text-[9px] text-slate-400">Kasir / Treasury</p>
            </div>
          </div>
          <div class="space-y-12">
            <p class="text-[10px] font-bold text-slate-500 uppercase">Diperiksa Oleh,</p>
            <div>
              <p class="font-bold text-slate-800 underline">Dewi Lestari, SE, Ak.</p>
              <p class="text-[9px] text-slate-400">Supervisor Akuntansi</p>
            </div>
          </div>
          <div class="space-y-12">
            <p class="text-[10px] font-bold text-slate-500 uppercase">Disetujui Oleh,</p>
            <div>
              <p class="font-bold text-slate-800 underline">Budi Santoso, SE, CA</p>
              <p class="text-[9px] text-slate-400">Manajer Keuangan</p>
            </div>
          </div>
          <div class="space-y-12">
            <p class="text-[10px] font-bold text-slate-500 uppercase">
              {{ activeVoucher.type === 'CASH_IN' ? 'Penyetor,' : 'Penerima,' }}
            </p>
            <div>
              <p class="font-bold text-slate-800 underline">({{ activeVoucher.recipientOrPayer || '........................' }})</p>
              <p class="text-[9px] text-slate-400">Pihak Ketiga</p>
            </div>
          </div>
        </div>

        <!-- Voucher Buttons -->
        <div class="flex items-center justify-end gap-2 pt-4 border-t print:hidden">
          <button
            @click="showVoucherModal = false"
            class="px-4 py-2 border border-slate-300 rounded-lg font-bold text-slate-700 hover:bg-slate-100 text-xs"
          >
            Tutup
          </button>
          <button
            onclick="window.print()"
            class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold text-xs shadow-xs"
          >
            <Printer class="w-3.5 h-3.5" />
            <span>Cetak Voucher</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
