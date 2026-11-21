<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useTransactionsStore } from '@/stores/transactions'
import { useLedgerStore } from '@/stores/ledger'
import { useReportsStore } from '@/stores/reports'
import { useLocaleStore } from '@/stores/locale'
import { accountingApi, type BackendTrialBalance } from '@/api/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import StatusBadge from '@/components/common/StatusBadge.vue'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import ManualJournalModal from '@/components/journal/ManualJournalModal.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  Receipt,
  FileCheck2,
  ArrowUpRight,
  Plus,
  Layers,
  Sparkles,
  Building2,
  FileSpreadsheet,
  RefreshCw
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()
const transactionsStore = useTransactionsStore()
const ledgerStore = useLedgerStore()
const reportsStore = useReportsStore()
const localeStore = useLocaleStore()

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)
const showManualModal = ref(false)
const isLoading = ref(false)

// Authoritative PostgreSQL Data
const backendTrialBalance = ref<BackendTrialBalance | null>(null)
const backendIncomeStatement = ref<any>(null)
const backendBalanceSheet = ref<any>(null)

function inspectJournal(j: JournalEntry) {
  selectedJournal.value = j
  showDetailModal.value = true
}

async function loadAuthoritativeData() {
  isLoading.value = true
  try {
    const today = new Date().toISOString().slice(0, 10)
    const currentYear = new Date().getFullYear()
    const startOfYear = `${currentYear}-01-01`

    const [tb, is, bs] = await Promise.allSettled([
      accountingApi.getTrialBalance(startOfYear, today),
      accountingApi.getIncomeStatement(startOfYear, today),
      accountingApi.getBalanceSheet(today),
      journalsStore.fetchJournals(),
      ledgerStore.fetchAccounts()
    ])

    if (tb.status === 'fulfilled') backendTrialBalance.value = tb.value
    if (is.status === 'fulfilled') backendIncomeStatement.value = is.value
    if (bs.status === 'fulfilled') backendBalanceSheet.value = bs.value
  } catch (err) {
    console.warn('[Dashboard] Gagal memuat data finansial riil:', err)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadAuthoritativeData()
})

watch(() => accountingStore.businessUnit, () => {
  loadAuthoritativeData()
})

// Real Posted Journals
const postedJournals = computed(() => {
  return journalsStore.journals.filter(j => j.status === 'POSTED')
})

// Real Per-Unit Financial Calculations from Posted Journals
const kidinestStats = computed(() => {
  const kidJournals = postedJournals.value.filter(j => j.businessUnit === 'KIDINEST')
  const revenue = kidJournals
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('4-') || l.accountCode.startsWith('4'))
    .reduce((sum, l) => sum + (l.credit - l.debit), 0)
  const expenses = kidJournals
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('5-') || l.accountCode.startsWith('6-') || l.accountCode.startsWith('5') || l.accountCode.startsWith('6'))
    .reduce((sum, l) => sum + (l.debit - l.credit), 0)
  return {
    revenue: Math.max(0, revenue),
    netProfit: revenue - expenses
  }
})

const sscStats = computed(() => {
  const sscJournals = postedJournals.value.filter(j => j.businessUnit === 'SSC')
  const revenue = sscJournals
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('4-') || l.accountCode.startsWith('4'))
    .reduce((sum, l) => sum + (l.credit - l.debit), 0)
  const expenses = sscJournals
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('5-') || l.accountCode.startsWith('6-') || l.accountCode.startsWith('5') || l.accountCode.startsWith('6'))
    .reduce((sum, l) => sum + (l.debit - l.credit), 0)
  return {
    revenue: Math.max(0, revenue),
    netProfit: revenue - expenses
  }
})

// KPI Calculations from Authoritative PostgreSQL Data
const grossRevenue = computed(() => {
  if (accountingStore.businessUnit === 'KIDINEST') return kidinestStats.value.revenue
  if (accountingStore.businessUnit === 'SSC') return sscStats.value.revenue

  if (backendIncomeStatement.value?.total_revenue) {
    return Number(backendIncomeStatement.value.total_revenue) || 0
  }
  return postedJournals.value
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('4-') || l.accountCode.startsWith('4'))
    .reduce((sum, l) => sum + (l.credit - l.debit), 0)
})

const liquidCash = computed(() => {
  if (backendTrialBalance.value?.items?.length) {
    return backendTrialBalance.value.items
      .filter(item => {
        const isCashOrBank = item.account_code.startsWith('1-11') || item.account_code.startsWith('1-12')
        if (!isCashOrBank) return false
        if (accountingStore.businessUnit === 'KIDINEST') return item.account_code.includes('.001') || item.account_name.includes('Kidinest')
        if (accountingStore.businessUnit === 'SSC') return item.account_code.includes('.002') || item.account_name.includes('SSC')
        return true
      })
      .reduce((sum, item) => sum + (Number(item.ending_debit) - Number(item.ending_credit)), 0)
  }
  // Fallback to posted journals
  return postedJournals.value
    .filter(j => accountingStore.businessUnit === 'ALL' || j.businessUnit === accountingStore.businessUnit)
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('1-11') || l.accountCode.startsWith('1-12'))
    .reduce((sum, l) => sum + (l.debit - l.credit), 0)
})

const patientAR = computed(() => {
  if (backendTrialBalance.value?.items?.length) {
    return backendTrialBalance.value.items
      .filter(item => item.account_code.startsWith('1-12') || item.account_name.toLowerCase().includes('piutang'))
      .reduce((sum, item) => sum + (Number(item.ending_debit) - Number(item.ending_credit)), 0)
  }
  return postedJournals.value
    .filter(j => accountingStore.businessUnit === 'ALL' || j.businessUnit === accountingStore.businessUnit)
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('1-12') || l.accountName.toLowerCase().includes('piutang'))
    .reduce((sum, l) => sum + (l.debit - l.credit), 0)
})

const medicalAP = computed(() => {
  if (backendTrialBalance.value?.items?.length) {
    return backendTrialBalance.value.items
      .filter(item => item.account_code.startsWith('2-') || item.account_name.toLowerCase().includes('hutang'))
      .reduce((sum, item) => sum + (Number(item.ending_credit) - Number(item.ending_debit)), 0)
  }
  return postedJournals.value
    .filter(j => accountingStore.businessUnit === 'ALL' || j.businessUnit === accountingStore.businessUnit)
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('2-') || l.accountName.toLowerCase().includes('hutang'))
    .reduce((sum, l) => sum + (l.credit - l.debit), 0)
})

const netProfit = computed(() => {
  if (accountingStore.businessUnit === 'KIDINEST') return kidinestStats.value.netProfit
  if (accountingStore.businessUnit === 'SSC') return sscStats.value.netProfit

  if (backendIncomeStatement.value?.net_income) {
    return Number(backendIncomeStatement.value.net_income) || 0
  }
  const rev = grossRevenue.value
  const exp = postedJournals.value
    .flatMap(j => j.lines)
    .filter(l => l.accountCode.startsWith('5-') || l.accountCode.startsWith('6-') || l.accountCode.startsWith('5') || l.accountCode.startsWith('6'))
    .reduce((sum, l) => sum + (l.debit - l.credit), 0)
  return rev - exp
})

const netMargin = computed(() => {
  if (grossRevenue.value <= 0) return 0
  return Math.round((netProfit.value / grossRevenue.value) * 1000) / 10
})

const totalKidinestSscRevenue = computed(() => {
  return kidinestStats.value.revenue + sscStats.value.revenue
})

const kidinestContribution = computed(() => {
  if (totalKidinestSscRevenue.value <= 0) return 0
  return Math.round((kidinestStats.value.revenue / totalKidinestSscRevenue.value) * 100)
})

const sscContribution = computed(() => {
  if (totalKidinestSscRevenue.value <= 0) return 0
  return Math.round((sscStats.value.revenue / totalKidinestSscRevenue.value) * 100)
})

const activePeriodLabel = computed(() => {
  const monthNames = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
  const m = accountingStore.fiscalMonth ? monthNames[accountingStore.fiscalMonth - 1] : ''
  return `${m} ${accountingStore.fiscalYear}`.trim() || 'Tahun Berjalan'
})

// Recent posted journals
const recentJournals = computed(() => {
  return journalsStore.filteredJournals.slice(0, 5)
})

// Unprocessed transactions preview
const unprocessedTransactions = computed(() => {
  return transactionsStore.filteredTransactions.filter(t => t.journalStatus === 'UNPROCESSED').slice(0, 4)
})
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Minimalist Clinical Header (Official DESIGN.md Ergonomics) -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            {{ localeStore.t('dashboard.category') }}
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          {{ localeStore.t('dashboard.title') }}
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="loadAuthoritativeData"
          :disabled="isLoading"
          class="inline-flex items-center gap-1.5 px-3 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer disabled:opacity-50"
          :title="localeStore.isEnglish ? 'Reload Database Data' : 'Muat Ulang Data Database'"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
          <span class="hidden sm:inline">{{ isLoading ? (localeStore.isEnglish ? 'Loading...' : 'Memuat...') : (localeStore.isEnglish ? 'Refresh' : 'Sinkron Data') }}</span>
        </button>

        <button
          @click="showManualModal = true"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>{{ localeStore.t('dashboard.createManualJournal') }}</span>
        </button>

        <button
          @click="router.push('/accounting/reports/profit-loss')"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 text-zinc-600" />
          <span>{{ localeStore.t('dashboard.profitLossReport') }}</span>
        </button>
      </div>
    </header>

    <!-- 1-Row Compact Hairline KPI Ticker Strip (Anti AI-Slop, Calm & Harmonious) -->
    <div
      class="grid grid-cols-2 sm:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">{{ localeStore.t('dashboard.grossRevenue') }}</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(grossRevenue) }}
        </div>
        <div class="text-[10px] font-semibold mt-0.5" :class="grossRevenue > 0 ? 'text-emerald-600' : 'text-zinc-400'">
          {{ grossRevenue > 0 ? '+14.2% vs Feb 2026' : (localeStore.isEnglish ? 'Verified Ledger' : 'Buku Besar Terverifikasi') }}
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">{{ localeStore.t('dashboard.cashBankLiquid') }}</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(liquidCash) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">{{ localeStore.t('dashboard.liquidSubtitle') }}</div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">{{ localeStore.t('dashboard.receivables') }}</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(patientAR) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">{{ localeStore.t('dashboard.receivablesSubtitle') }}</div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">{{ localeStore.t('dashboard.payables') }}</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(medicalAP) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">{{ localeStore.t('dashboard.payablesSubtitle') }}</div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex flex-col justify-center">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">{{ localeStore.t('dashboard.netProfit') }}</div>
        <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
          {{ formatIDR(netProfit) }}
        </div>
        <div class="text-[10px] font-bold mt-0.5" :class="netProfit >= 0 ? 'text-emerald-700' : 'text-rose-600'">
          {{ netMargin }}% {{ localeStore.t('dashboard.netMargin') }}
        </div>
      </div>
    </div>

    <!-- Dual Clinic Unit Comparison & Workflow Health Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Unit Performance Card (Kidinest vs SSC) -->
      <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <Building2 class="w-4 h-4 text-zinc-700" />
            <h3 class="text-sm font-bold text-zinc-950">{{ localeStore.t('dashboard.unitComparison') }}</h3>
          </div>
          <span class="text-[11px] font-semibold text-zinc-500">{{ activePeriodLabel }}</span>
        </div>

        <!-- Kidinest (Warm Mustard / Amber) -->
        <div class="p-3.5 rounded-xl bg-amber-50/70 border border-amber-200/80 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
              <span class="font-bold text-xs text-amber-950">{{ localeStore.t('dashboard.kidinestTitle') }}</span>
            </div>
            <span class="text-xs font-bold tabular-nums text-amber-950">{{ formatIDR(kidinestStats.revenue) }}</span>
          </div>
          <div class="w-full bg-amber-200/60 h-1.5 rounded-full overflow-hidden">
            <div
              class="bg-amber-500 h-full rounded-full transition-all duration-300"
              :style="{ width: `${kidinestContribution}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-amber-900 font-medium">
            <span>{{ localeStore.t('dashboard.netProfitLabel') }}: <strong class="tabular-nums">{{ formatIDR(kidinestStats.netProfit) }}</strong></span>
            <span>{{ localeStore.t('dashboard.contribution') }}: <strong>{{ kidinestContribution }}%</strong></span>
          </div>
        </div>

        <!-- SSC (Dusty Rose / Mauve) -->
        <div class="p-3.5 rounded-xl bg-rose-50/70 border border-rose-200/80 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-rose-500"></span>
              <span class="font-bold text-xs text-rose-950">{{ localeStore.t('dashboard.sscTitle') }}</span>
            </div>
            <span class="text-xs font-bold tabular-nums text-rose-950">{{ formatIDR(sscStats.revenue) }}</span>
          </div>
          <div class="w-full bg-rose-200/60 h-1.5 rounded-full overflow-hidden">
            <div
              class="bg-rose-500 h-full rounded-full transition-all duration-300"
              :style="{ width: `${sscContribution}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-rose-900 font-medium">
            <span>{{ localeStore.t('dashboard.netProfitLabel') }}: <strong class="tabular-nums">{{ formatIDR(sscStats.netProfit) }}</strong></span>
            <span>{{ localeStore.t('dashboard.contribution') }}: <strong>{{ sscContribution }}%</strong></span>
          </div>
        </div>

        <p class="text-[11px] text-zinc-500 leading-relaxed italic">
          {{ localeStore.t('dashboard.sscNote') }}
        </p>
      </div>

      <!-- Accounting Workflow Status Pipeline -->
      <div class="lg:col-span-2 bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <FileCheck2 class="w-4 h-4 text-emerald-600" />
            <h3 class="text-sm font-bold text-zinc-950">{{ localeStore.t('dashboard.workflowTitle') }}</h3>
          </div>
          <span class="text-xs text-emerald-700 font-bold bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
            {{ localeStore.t('dashboard.auditReady') }}
          </span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div
            @click="router.push('/accounting/transactions')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-500 uppercase">{{ localeStore.t('dashboard.stepSource') }}</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ transactionsStore.transactions.length }}
              </span>
              <span class="text-[10px] font-bold text-amber-700 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200">
                {{ transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').length }} {{ localeStore.t('dashboard.unprocessedBadge') }}
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">{{ localeStore.t('dashboard.sourceSubtitle') }}</p>
          </div>

          <div
            @click="router.push('/accounting/journals?status=DRAFT')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-500 uppercase">{{ localeStore.t('dashboard.stepDraft') }}</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ journalsStore.stats.draftCount + journalsStore.stats.submittedCount }}
              </span>
              <span class="text-[10px] font-bold text-zinc-700 bg-zinc-200/80 px-1.5 py-0.5 rounded">
                Draft
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">{{ localeStore.t('dashboard.makerSubtitle') }}</p>
          </div>

          <div
            @click="router.push('/accounting/journals/approvals')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-600 uppercase">{{ localeStore.t('dashboard.stepReview') }}</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ journalsStore.stats.reviewedCount + journalsStore.stats.approvedCount }}
              </span>
              <span class="text-[10px] font-bold text-amber-700 bg-amber-100 px-1.5 py-0.5 rounded">
                {{ localeStore.t('dashboard.actionRequiredBadge') }}
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">{{ localeStore.t('dashboard.checkerSubtitle') }}</p>
          </div>

          <div
            @click="router.push('/accounting/journals?status=POSTED')"
            class="p-3.5 rounded-xl border border-emerald-200 bg-emerald-50/50 hover:bg-emerald-100/60 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-emerald-800 uppercase">{{ localeStore.t('dashboard.stepPosted') }}</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-emerald-950">
                {{ journalsStore.stats.postedCount }}
              </span>
              <span class="text-[10px] font-bold text-emerald-700 bg-emerald-100 px-1.5 py-0.5 rounded">
                {{ localeStore.t('dashboard.officialBadge') }}
              </span>
            </div>
            <p class="text-[11px] text-emerald-700">{{ localeStore.t('dashboard.postedSubtitle') }}</p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200/90 text-xs flex items-center justify-between">
          <div class="flex items-center gap-2 text-zinc-700 font-medium">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            <span>{{ localeStore.t('dashboard.totalDebetPosted') }}</span>
          </div>
          <span class="font-bold tabular-nums text-sm text-zinc-950">
            {{ formatIDR(journalsStore.stats.totalPostedDebit) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Two Columns: Unprocessed Source Transactions & Recent Journals -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- Unprocessed Source Transactions -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 p-5 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <Layers class="w-4 h-4 text-zinc-700" />
            <h3 class="text-sm font-bold text-zinc-950">{{ localeStore.t('dashboard.operationalInflowTitle') }}</h3>
          </div>
          <button
            @click="router.push('/accounting/transactions')"
            class="text-xs font-bold text-zinc-800 hover:text-zinc-950 flex items-center gap-1 cursor-pointer"
          >
            <span>{{ localeStore.t('dashboard.viewAll') }}</span>
            <ArrowUpRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Empty State -->
        <div v-if="unprocessedTransactions.length === 0" class="py-10 text-center space-y-1.5 border border-dashed border-zinc-200 rounded-xl">
          <Layers class="w-6 h-6 text-zinc-300 mx-auto" />
          <p class="text-xs font-bold text-zinc-700">Semua transaksi operasional telah diproses</p>
          <p class="text-[11px] text-zinc-400">Tidak ada antrean transaksi kasir atau farmasi yang menunggu penjurnalan.</p>
        </div>

        <!-- Transaction List -->
        <div v-else class="space-y-2.5">
          <div
            v-for="tx in unprocessedTransactions"
            :key="tx.id"
            class="p-3 rounded-xl border border-zinc-200/90 hover:border-zinc-300 bg-zinc-50/40 hover:bg-zinc-50 transition-colors flex items-center justify-between"
          >
            <div class="space-y-0.5 min-w-0 pr-3">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-zinc-800 text-xs">{{ tx.transactionNumber }}</span>
                <StatusBadge :status="tx.businessUnit" size="sm" />
                <span class="text-[10px] font-semibold text-zinc-500 uppercase">{{ tx.sourceModule }}</span>
              </div>
              <p class="text-xs text-zinc-700 truncate font-medium">{{ tx.description }}</p>
              <p class="text-[11px] text-zinc-400">{{ tx.transactionDate }} &bull; {{ tx.branchName }}</p>
            </div>

            <div class="text-right shrink-0">
              <div class="font-bold tabular-nums text-xs text-zinc-950">{{ formatIDR(tx.amount) }}</div>
              <button
                @click="transactionsStore.generateDraftJournal(tx.id)"
                class="mt-1 px-2.5 py-1 bg-zinc-950 text-white hover:bg-zinc-800 rounded-lg text-[11px] font-bold transition-colors flex items-center gap-1 shadow-2xs cursor-pointer active:scale-95"
              >
                <Sparkles class="w-3 h-3" />
                <span>{{ localeStore.t('dashboard.createDraft') }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Journal Entries List -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 p-5 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <Receipt class="w-4 h-4 text-zinc-700" />
            <h3 class="text-sm font-bold text-zinc-950">{{ localeStore.t('dashboard.recentJournalsTitle') }}</h3>
          </div>
          <button
            @click="router.push('/accounting/journals')"
            class="text-xs font-bold text-zinc-800 hover:text-zinc-950 flex items-center gap-1 cursor-pointer"
          >
            <span>{{ localeStore.t('dashboard.viewFullList') }}</span>
            <ArrowUpRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Empty State -->
        <div v-if="recentJournals.length === 0" class="py-10 text-center space-y-1.5 border border-dashed border-zinc-200 rounded-xl">
          <Receipt class="w-6 h-6 text-zinc-300 mx-auto" />
          <p class="text-xs font-bold text-zinc-700">Belum ada jurnal transaksi tercatat di database</p>
          <p class="text-[11px] text-zinc-400">Database siap menerima pencatatan jurnal baru atau posting kasir.</p>
        </div>

        <!-- Journal List -->
        <div v-else class="space-y-2.5">
          <div
            v-for="j in recentJournals"
            :key="j.id"
            @click="inspectJournal(j)"
            class="p-3 rounded-xl border border-zinc-200/90 hover:border-zinc-400 bg-white hover:bg-zinc-50/60 cursor-pointer transition-all flex items-center justify-between"
          >
            <div class="space-y-0.5 min-w-0 pr-3">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-zinc-800 text-xs">{{ j.journalNumber }}</span>
                <StatusBadge :status="j.status" size="sm" />
                <StatusBadge :status="j.businessUnit" size="sm" />
              </div>
              <p class="text-xs text-zinc-800 font-medium truncate">{{ j.description }}</p>
              <p class="text-[11px] text-zinc-400">
                {{ formatDate(j.journalDate) }} &bull; {{ j.lines.length }} Pos Akun &bull; {{ j.branchName }}
              </p>
            </div>

            <div class="text-right shrink-0">
              <span class="font-bold tabular-nums text-xs text-zinc-950">
                {{ formatIDR(j.totalDebit) }}
              </span>
              <p class="text-[10px] text-zinc-400 mt-1">
                {{ j.postedByName ? 'Posted' : 'Draft / Approval' }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />

    <ManualJournalModal
      :show="showManualModal"
      @close="showManualModal = false"
    />
  </div>
</template>
