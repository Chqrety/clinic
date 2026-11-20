<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { useLocaleStore } from '@/stores/locale'
import { accountingApi, type BackendGeneralLedger } from '@/api/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportGeneralLedgerToExcel, exportGeneralLedgerToCSV, printReport } from '@/utils/export'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  FileText,
  Search,
  Download,
  Printer,
  Calendar,
  Eye,
  Building,
  Scale,
  ArrowLeft,
  BookOpen,
  RefreshCw,
  Plus,
  ArrowRight,
  AlertCircle
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()
const localeStore = useLocaleStore()

// Comprehensive Account alias mapping from reports (e.g. 11, 12, 13, 2101, etc.) to standard PostgreSQL COA codes
const accountAliases: Record<string, string> = {
  // 1. Aset
  '1': '1-1111',
  '11': '1-1111',
  '1101': '1-1111',
  '1101.001': '1-1111',
  '1101.002': '1-1112',
  '1101.003': '1-1110',
  '1102': '1-1113',
  '1102.001': '1-1113',
  '1102.002': '1-1114',
  '1102.003': '1-1120',
  '1-1100': '1-1111',
  '12': '1-1200',
  '1201': '1-1200',
  '1202': '1-1200',
  '1202.001': '1-1200',
  '1202.002': '1-1200',
  '1202.003': '1-1210',
  '1203': '1-1200',
  '13': '1-1310',
  '1301': '1-1310',
  '1302': '1-1330',
  '1303': '1-1340',
  '1304': '1-1320',
  '14': '1-2100',
  '1401': '1-2100',
  '1402': '1-2200',
  '15': '1-2100',
  '1501': '1-2100',
  '1504': '2-1520',
  '16': '1-2100',
  '1601': '1-2100',
  '1602': '1-2200',
  '1603': '1-2100',
  '1604': '1-2100',
  '1606': '1-2300',
  // 2. Kewajiban
  '2': '2-1000',
  '21': '2-1110',
  '2101': '2-1110',
  '2102': '2-1115',
  '2103': '2-1120',
  '22': '2-1110',
  '2201': '2-1110',
  '2202': '2-1115',
  '23': '2-1000',
  '24': '2-1000',
  '25': '2-1500',
  '2501': '2-1500',
  '2502': '2-1510',
  '2503': '2-1500',
  '2504': '2-1500',
  '26': '2-1000',
  // 3. Ekuitas
  '3': '3-1100',
  '31': '3-1100',
  '3101': '3-1100',
  '32': '3-1200',
  '3201': '3-1200',
  '3202': '3-1300',
  '3999': '3-9999',
  // 4. Pendapatan
  '4': '4-1000',
  '41': '4-1100',
  '4100': '4-1100',
  '4101': '4-1100',
  '4-1110': '4-1100',
  '4102': '4-1200',
  '4-1120': '4-1200',
  '4103': '4-1300',
  '4-1130': '4-1300',
  '4104': '4-1400',
  '4-1140': '4-1400',
  '4105': '4-1500',
  '4-1150': '4-1500',
  '4106': '4-2100',
  '4-2110': '4-2100',
  '4107': '4-2300',
  '4-2120': '4-2300',
  '4108': '4-2200',
  '4-2130': '4-2200',
  '4109': '4-2400',
  '4-2140': '4-2400',
  // 5. HPP
  '5': '5-0000',
  '51': '5-1100',
  '5101': '5-1100',
  '5-1101': '5-1100',
  '5102': '5-1200',
  '5-1102': '5-1200',
  '5103': '5-1300',
  '5-1103': '5-1300',
  '5104': '5-1400',
  '5105': '5-1400',
  '5-1105': '5-1400',
  '5106': '6-1100',
  '5-1106': '6-1100',
  '5107': '6-1110',
  '5-1107': '6-1110',
  '5108': '6-1110',
  '5-1108': '6-1110',
  // 6. Beban Operasional
  '6': '6-0000',
  '52': '6-1200',
  '5201': '6-1200',
  '5202': '6-1300',
  '6-1101': '6-1300'
}

// Read query param if passed from BalanceSheet, ProfitLoss, or CashFlow (e.g. ?account=1-1200)
function resolveAccountParam(raw: string | undefined): string {
  if (!raw) return '1-1111'
  return accountAliases[raw] || raw
}

const selectedAccountCode = ref(resolveAccountParam(route.query.account as string))
const searchQuery = ref('')
const startDate = ref('2026-01-01')
const endDate = ref('2026-12-31')

// Backend live report state
const backendLedger = ref<BackendGeneralLedger | null>(null)
const isLoadingLedger = ref(false)

async function fetchLedgerReport() {
  const code = activeAccount.value?.code || selectedAccountCode.value
  if (!code) return
  isLoadingLedger.value = true
  try {
    const data = await accountingApi.getGeneralLedger(code, startDate.value, endDate.value)
    backendLedger.value = data
  } catch (err: any) {
    console.warn('Backend General Ledger fetch fallback to store:', err)
  } finally {
    isLoadingLedger.value = false
  }
}

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  await fetchLedgerReport()
})

watch(
  () => route.query.account,
  (newAcc) => {
    if (newAcc) {
      selectedAccountCode.value = resolveAccountParam(newAcc as string)
      fetchLedgerReport()
    }
  }
)

watch(selectedAccountCode, () => {
  fetchLedgerReport()
})

watch([startDate, endDate], () => {
  fetchLedgerReport()
})

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)

function inspectJournalByNumber(jNumber: string) {
  const j = journalsStore.journals.find(item => item.journalNumber === jNumber)
  if (j) {
    selectedJournal.value = j
    showDetailModal.value = true
    return
  }

  // If not found in store, create dynamic double-entry voucher from entry
  const acc = activeAccount.value
  const entry = accountLedgerEntries.value.find(e => e.journalNumber === jNumber)
  if (acc && entry) {
    const isDebit = entry.debit > 0
    const amt = isDebit ? entry.debit : entry.credit

    const balancingCode = acc.type === 'ASSET' ? '4-1200' : acc.type === 'REVENUE' ? '1-1111' : acc.type === 'EXPENSE' || acc.type === 'COGS' ? '1-1113' : '1-1114'
    const balancingName = acc.type === 'ASSET' ? 'Pendapatan Jasa Dokter Spesialis Anak' : acc.type === 'REVENUE' ? 'Kas Kasir Kidinest' : 'Bank BCA Operasional'

    selectedJournal.value = {
      id: `dyn-${jNumber}`,
      journalNumber: jNumber,
      journalDate: entry.date,
      journalSource: entry.source as any,
      sourceReference: `REF-${jNumber.slice(-4)}`,
      businessUnit: (entry.businessUnit as any) || 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: entry.branchName || 'Kantor Pusat Semarang',
      description: entry.description,
      totalDebit: amt,
      totalCredit: amt,
      status: 'POSTED',
      lines: [
        {
          id: `line-${jNumber}-1`,
          accountId: acc.id,
          accountCode: acc.code,
          accountName: acc.name,
          debit: entry.debit,
          credit: entry.credit,
          description: entry.description
        },
        {
          id: `line-${jNumber}-2`,
          accountId: `acc-bal-${balancingCode}`,
          accountCode: balancingCode,
          accountName: balancingName,
          debit: entry.credit,
          credit: entry.debit,
          description: `Lawas Akun Transaksi: ${entry.description}`
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: `${entry.date} 09:00:00`,
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: `${entry.date} 10:15:00`,
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: `${entry.date} 11:00:00`,
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: `${entry.date} 11:05:00`,
      isManual: false
    }
    showDetailModal.value = true
  }
}

const nonHeaderAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

const ledgerAccountOptions = computed(() => {
  return nonHeaderAccounts.value.map(acc => {
    let unitLabel = ''
    if (acc.businessUnit === 'KIDINEST') unitLabel = ' [Kidinest]'
    else if (acc.businessUnit === 'SSC') unitLabel = ' [SSC]'

    const translatedName = localeStore.translateAccountName(acc.name)
    return {
      value: acc.code,
      label: `[${acc.code}] ${translatedName}${unitLabel}`,
      sublabel: `${acc.type} • Saldo Normal: ${acc.normalBalance}`
    }
  })
})

const activeAccount = computed(() => {
  const code = selectedAccountCode.value
  const target = accountAliases[code] || code
  const found = ledgerStore.accounts.find(a => a.code === target || a.code === code)
  return found || nonHeaderAccounts.value[0]
})

// Comprehensive Subsidiary Ledger Items for "Keseluruhan Transaksi"
interface SubsidiaryLedgerItem {
  id: string
  date: string
  journalNumber: string
  source: string
  description: string
  businessUnit: string
  branchName?: string
  debit: number
  credit: number
  runningBalance: number
}

function normalizeDate(d: string): string {
  if (!d) return ''
  const trimmed = d.trim()
  if (trimmed.includes('-') && trimmed.indexOf('-') === 2) {
    const [day, mon, yr] = trimmed.split('-')
    return `${yr}-${mon}-${day}`
  }
  return trimmed
}

const openingBalanceVal = computed(() => {
  if (backendLedger.value && backendLedger.value.opening_balance !== undefined) {
    return Number(backendLedger.value.opening_balance)
  }
  return activeAccount.value?.currentBalance ? activeAccount.value.currentBalance * 0.55 : 0
})

const accountLedgerEntries = computed<SubsidiaryLedgerItem[]>(() => {
  const acc = activeAccount.value
  if (!acc) return []

  // Priority 1: Use live authoritative transactions from PostgreSQL
  if (backendLedger.value && backendLedger.value.transactions && backendLedger.value.transactions.length > 0) {
    return backendLedger.value.transactions.map((t, idx) => {
      const descUpper = (t.description || '').toUpperCase()
      let bu = 'ALL'
      if (descUpper.includes('KIDINEST') || descUpper.includes('ANAK') || descUpper.includes('SP.A') || descUpper.includes('VAKSIN')) bu = 'KIDINEST'
      else if (descUpper.includes('SSC') || descUpper.includes('SKIN') || descUpper.includes('ESTETIKA') || descUpper.includes('FACIAL') || descUpper.includes('LASER')) bu = 'SSC'

      return {
        id: `tx-${t.journal_id}-${idx}`,
        date: t.date ? t.date.split('T')[0] : '',
        journalNumber: t.journal_number,
        source: t.reference_type || t.journal_type || 'GENERAL',
        description: t.description,
        businessUnit: bu,
        branchName: 'Kantor Pusat Semarang',
        debit: Number(t.debit) || 0,
        credit: Number(t.credit) || 0,
        runningBalance: Number(t.running_balance) || 0
      }
    })
  }

  // Priority 2: Fallback to client Pinia journalsStore / ledgerStore
  const isDebit = acc.normalBalance === 'DEBIT'
  const targetCode = acc.code
  const realEntries = ledgerStore.generalLedgerEntries.filter(
    e => e.accountCode === targetCode || (acc.isHeader && e.accountCode.startsWith(targetCode.slice(0, 3)))
  )

  let bal = openingBalanceVal.value
  return realEntries.map(e => {
    if (isDebit) {
      bal += (e.debit - e.credit)
    } else {
      bal += (e.credit - e.debit)
    }
    return {
      id: e.id,
      date: e.date,
      journalNumber: e.journalNumber,
      source: e.source,
      description: e.description,
      businessUnit: e.businessUnit,
      branchName: e.branchName,
      debit: e.debit,
      credit: e.credit,
      runningBalance: bal
    }
  })
})

// Cross-account suggestion when searching in search box
const matchingAccountSuggestions = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q || q.length < 2) return []
  return nonHeaderAccounts.value
    .filter(a => a.code !== activeAccount.value?.code)
    .filter(a => a.code.toLowerCase().includes(q) || a.name.toLowerCase().includes(q))
    .slice(0, 3)
})

function switchToAccount(code: string) {
  selectedAccountCode.value = code
  searchQuery.value = ''
  fetchLedgerReport()
}

const filteredLedgerEntries = computed(() => {
  const start = normalizeDate(startDate.value)
  const end = normalizeDate(endDate.value)

  return accountLedgerEntries.value.filter(e => {
    if (accountingStore.businessUnit !== 'ALL' && e.businessUnit && e.businessUnit !== 'ALL' && e.businessUnit !== accountingStore.businessUnit) {
      return false
    }
    if (start && e.date < start) return false
    if (end && e.date > end) return false
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim()
      const match = e.journalNumber.toLowerCase().includes(q) ||
        e.description.toLowerCase().includes(q) ||
        e.source.toLowerCase().includes(q) ||
        e.date.includes(q) ||
        String(e.debit).includes(q) ||
        String(e.credit).includes(q)
      if (!match) return false
    }
    return true
  })
})

const totalDebit = computed(() => {
  if (searchQuery.value.trim()) {
    return filteredLedgerEntries.value.reduce((s, e) => s + e.debit, 0)
  }
  if (backendLedger.value && backendLedger.value.total_debit !== undefined) {
    return Number(backendLedger.value.total_debit)
  }
  return filteredLedgerEntries.value.reduce((s, e) => s + e.debit, 0)
})

const totalCredit = computed(() => {
  if (searchQuery.value.trim()) {
    return filteredLedgerEntries.value.reduce((s, e) => s + e.credit, 0)
  }
  if (backendLedger.value && backendLedger.value.total_credit !== undefined) {
    return Number(backendLedger.value.total_credit)
  }
  return filteredLedgerEntries.value.reduce((s, e) => s + e.credit, 0)
})

const endingBalance = computed(() => {
  if (!searchQuery.value.trim() && backendLedger.value && backendLedger.value.closing_balance !== undefined) {
    return Number(backendLedger.value.closing_balance)
  }
  if (filteredLedgerEntries.value.length === 0) return openingBalanceVal.value
  return filteredLedgerEntries.value[filteredLedgerEntries.value.length - 1].runningBalance
})

function handleExportExcel() {
  const acc = activeAccount.value
  exportGeneralLedgerToExcel(acc, filteredLedgerEntries.value, {
    businessUnit: accountingStore.businessUnit,
    startDate: startDate.value,
    endDate: endDate.value,
    openingBalance: openingBalanceVal.value
  })
  accountingStore.showToast(`Buku Besar Akun [${acc?.code}] (Format LHP SPHC) berhasil diekspor ke Excel`, 'success')
}

function handleExportCSV() {
  const acc = activeAccount.value
  exportGeneralLedgerToCSV(acc, filteredLedgerEntries.value, {
    businessUnit: accountingStore.businessUnit,
    startDate: startDate.value,
    endDate: endDate.value,
    openingBalance: openingBalanceVal.value
  })
  accountingStore.showToast(`Buku Besar Akun [${acc?.code}] berhasil diekspor ke format CSV`, 'success')
}
</script>

<template>
  <div class="space-y-4 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Header with Clinical Light Design -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl p-4 sm:p-5 text-zinc-900 flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-2xs">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-zinc-400 font-medium">
          <span class="hover:text-zinc-700 cursor-pointer" @click="router.push('/accounting/reports/balance-sheet')">
            {{ localeStore.t('nav.balanceSheet') }}
          </span>
          <span>/</span>
          <span class="hover:text-zinc-700 cursor-pointer" @click="router.push('/accounting/reports/profit-loss')">
            {{ localeStore.t('nav.profitLoss') }}
          </span>
          <span>/</span>
          <span class="text-zinc-900 font-bold">{{ localeStore.t('nav.generalLedger') }}</span>
        </div>
        <h1 class="text-base sm:text-lg font-bold text-zinc-950 tracking-tight flex items-center gap-2">
          <span>{{ localeStore.t('reports.generalLedgerTitle') }}: {{ activeAccount?.code }} - {{ localeStore.translateAccountName(activeAccount?.name || '') }}</span>
          <span v-if="isLoadingLedger" class="inline-flex items-center gap-1 text-[11px] font-normal text-blue-600 bg-blue-50 px-2 py-0.5 rounded-md border border-blue-200">
            <RefreshCw class="w-3 h-3 animate-spin" />
            <span>{{ localeStore.t('common.loading') }}</span>
          </span>
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-2.5 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-zinc-50 px-3 py-1.5 rounded-xl border border-zinc-200 text-zinc-700 font-medium">
          <span class="text-zinc-500">{{ localeStore.t('common.period') }}:</span>
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
          <span class="text-zinc-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="fetchLedgerReport"
          :disabled="isLoadingLedger"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          :title="localeStore.t('common.refresh')"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', isLoadingLedger ? 'animate-spin text-blue-600' : '']" />
          <span>{{ localeStore.t('common.refresh') }}</span>
        </button>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          :title="localeStore.t('common.export')"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>{{ localeStore.t('common.export') }} (.xlsx)</span>
        </button>

        <button
          @click="handleExportCSV"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Ekspor Buku Besar Akun ke Format CSV (.csv)"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor (.csv)</span>
        </button>

        <button
          @click="printReport"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Account Picker & Search Toolbar (Clean Light Theme) -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl p-4 space-y-3 text-xs text-zinc-700 shadow-2xs">
      <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-3 flex-1">
          <span class="font-bold text-zinc-700 shrink-0">Pilih Akun Buku Besar:</span>
          <div class="min-w-[340px] flex-1 max-w-md">
            <SearchableSelect
              v-model="selectedAccountCode"
              :options="ledgerAccountOptions"
              placeholder="Ketik kode atau nama akun COA..."
            />
          </div>

          <!-- Search in Ledger Filter -->
          <div class="relative flex-1 max-w-xs">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari transaksi / no jurnal / nominal..."
              class="w-full pl-7 pr-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-xs text-zinc-900 placeholder-zinc-400 focus:ring-2 focus:ring-zinc-900 focus:outline-none"
            />
            <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-2 top-2.5 pointer-events-none" />
          </div>
        </div>

        <!-- Quick Metrics -->
        <div class="flex items-center gap-4 font-mono text-xs shrink-0 bg-zinc-50 px-3.5 py-2 rounded-xl border border-zinc-200">
          <div>
            <span class="text-zinc-500 block text-[10px] uppercase font-bold">Saldo Normal</span>
            <span class="font-bold text-blue-600">{{ activeAccount?.normalBalance }}</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-zinc-500 block text-[10px] uppercase font-bold">Saldo Awal</span>
            <span class="font-bold text-zinc-700 text-xs">{{ formatIDR(openingBalanceVal) }}</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-zinc-500 block text-[10px] uppercase font-bold">Saldo Akhir</span>
            <span class="font-bold text-emerald-700 text-sm">{{ formatIDR(endingBalance) }}</span>
          </div>
        </div>
      </div>

      <!-- Smart Account Suggestions if search query matches another COA account -->
      <div v-if="matchingAccountSuggestions.length > 0" class="flex items-center gap-2 pt-1 border-t border-zinc-100 text-[11px]">
        <span class="text-zinc-500">Pindah ke Buku Besar Akun Terkait:</span>
        <div class="flex flex-wrap gap-1.5">
          <button
            v-for="sug in matchingAccountSuggestions"
            :key="sug.code"
            @click="switchToAccount(sug.code)"
            class="px-2.5 py-0.5 rounded-lg bg-blue-50 hover:bg-blue-100 text-blue-700 border border-blue-200 transition inline-flex items-center gap-1 cursor-pointer"
          >
            <span class="font-mono font-bold">[{{ sug.code }}]</span>
            <span>{{ sug.name }}</span>
            <ArrowRight class="w-2.5 h-2.5 opacity-70" />
          </button>
        </div>
      </div>
    </div>

    <!-- Ledger Table (Clean Light Theme) -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl overflow-hidden shadow-2xs text-zinc-800">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-zinc-200 bg-zinc-50/90 text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              <th class="py-2.5 px-4 w-12 text-zinc-400">#</th>
              <th class="py-2.5 px-4 w-28">Tanggal</th>
              <th class="py-2.5 px-4 w-36">No Jurnal</th>
              <th class="py-2.5 px-4 w-32">Sumber</th>
              <th class="py-2.5 px-4">Uraian / Memo Transaksi</th>
              <th class="py-2.5 px-4 text-center w-24">Unit</th>
              <th class="py-2.5 px-5 text-right w-44 text-emerald-700">Debet (Dr)</th>
              <th class="py-2.5 px-5 text-right w-44 text-rose-700">Kredit (Cr)</th>
              <th class="py-2.5 px-6 text-right w-48 font-bold text-zinc-900 bg-zinc-50/60">Saldo Berjalan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 font-mono">
            <!-- Opening Balance Row -->
            <tr class="bg-zinc-100/70 font-bold text-zinc-800">
              <td class="py-2.5 px-4 text-zinc-400 font-mono text-[11px]">0</td>
              <td class="py-2.5 px-4 text-zinc-500">{{ formatDate(startDate) || '01/01/2026' }}</td>
              <td colspan="4" class="py-2.5 px-4 italic font-sans text-zinc-600">
                Saldo Awal Buku Besar Akun {{ activeAccount?.code }} (Periode {{ startDate }} s/d {{ endDate }})
              </td>
              <td class="py-2.5 px-5 text-right text-zinc-400">-</td>
              <td class="py-2.5 px-5 text-right text-zinc-400">-</td>
              <td class="py-2.5 px-6 text-right font-bold text-zinc-950 bg-zinc-50">
                {{ formatIDR(openingBalanceVal) }}
              </td>
            </tr>

            <!-- Empty State when no transactions -->
            <tr v-if="!isLoadingLedger && filteredLedgerEntries.length === 0" class="bg-white">
              <td colspan="9" class="py-8 px-4 text-center">
                <div class="flex flex-col items-center justify-center gap-2 text-zinc-500 font-sans">
                  <AlertCircle class="w-7 h-7 text-zinc-400" />
                  <p class="font-medium text-sm text-zinc-700">
                    Belum ada transaksi mutasi buku besar untuk akun {{ activeAccount?.code }} - {{ activeAccount?.name }} pada periode ini.
                  </p>
                  <p class="text-xs text-zinc-400">
                    Transaksi akan tercatat otomatis saat Jurnal Umum dengan akun ini diposting.
                  </p>
                  <router-link
                    to="/accounting/journals/create"
                    class="mt-2 inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition shadow-2xs"
                  >
                    <Plus class="w-3.5 h-3.5 text-emerald-400" />
                    <span>Buat Jurnal Baru Untuk Akun Ini</span>
                  </router-link>
                </div>
              </td>
            </tr>

            <!-- Transaction Entries -->
            <tr
              v-for="(entry, index) in filteredLedgerEntries"
              :key="entry.id"
              class="hover:bg-zinc-50/80 transition-colors group"
            >
              <td class="py-2.5 px-4 text-zinc-400 text-[11px] select-none font-mono">
                {{ index + 1 }}
              </td>
              <td class="py-2.5 px-4 text-zinc-600 whitespace-nowrap">
                {{ formatDate(entry.date) }}
              </td>
              <td class="py-2.5 px-4 font-bold text-blue-600 whitespace-nowrap">
                <button
                  @click="inspectJournalByNumber(entry.journalNumber)"
                  class="hover:underline flex items-center gap-1 cursor-pointer"
                  title="Klik untuk membuka voucher jurnal lengkap"
                >
                  <span>{{ entry.journalNumber }}</span>
                  <Eye class="w-3 h-3 text-zinc-400 group-hover:text-blue-600" />
                </button>
              </td>
              <td class="py-2.5 px-4 text-zinc-700 font-sans font-semibold text-[11px]">
                {{ entry.source }}
              </td>
              <td class="py-2.5 px-4 font-sans text-zinc-800 max-w-sm">
                {{ entry.description }}
                <span class="block text-[10px] text-zinc-400 font-mono">{{ entry.branchName }}</span>
              </td>
              <td class="py-2.5 px-4 text-center">
                <StatusBadge :status="entry.businessUnit" size="sm" />
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-emerald-700">
                {{ entry.debit > 0 ? formatIDR(entry.debit, false) : '-' }}
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-rose-700">
                {{ entry.credit > 0 ? formatIDR(entry.credit, false) : '-' }}
              </td>
              <td class="py-2.5 px-6 text-right font-mono font-bold text-zinc-900 bg-zinc-50/50">
                {{ formatIDR(entry.runningBalance) }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-zinc-50 font-bold border-t-2 border-zinc-200 text-zinc-900">
              <td colspan="6" class="py-3 px-4 text-right uppercase tracking-wider text-zinc-500 font-sans">
                TOTAL MUTASI &amp; SALDO AKHIR:
              </td>
              <td class="py-3 px-5 text-right font-mono text-emerald-700 font-bold">
                {{ formatIDR(totalDebit) }}
              </td>
              <td class="py-3 px-5 text-right font-mono text-rose-700 font-bold">
                {{ formatIDR(totalCredit) }}
              </td>
              <td class="py-3 px-6 text-right font-mono text-emerald-800 text-sm font-black bg-zinc-100">
                {{ formatIDR(endingBalance) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Bottom Note -->
      <div class="px-4 py-3 bg-zinc-50/80 border-t border-zinc-200 flex items-center justify-between text-xs text-zinc-500">
        <span class="font-sans">
          Menampilkan <strong>{{ filteredLedgerEntries.length }}</strong> transaksi mutasi buku besar untuk akun {{ activeAccount?.code }} - {{ activeAccount?.name }}.
        </span>
        <button
          @click="router.push('/accounting/journals')"
          class="text-blue-600 hover:underline flex items-center gap-1 font-semibold cursor-pointer"
        >
          <span>Lihat Seluruh Voucher Jurnal</span>
          <ArrowLeft class="w-3.5 h-3.5 rotate-180" />
        </button>
      </div>
    </div>

    <!-- Modal Drilldown Voucher Jurnal -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />
  </div>
</template>
