<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useTransactionsStore } from '@/stores/transactions'
import { useReportsStore } from '@/stores/reports'
import { formatIDR, compactIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import KpiCard from '@/components/common/KpiCard.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import ManualJournalModal from '@/components/journal/ManualJournalModal.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  TrendingUp,
  CreditCard,
  Wallet,
  Receipt,
  FileCheck2,
  ArrowUpRight,
  Plus,
  Layers,
  Sparkles,
  Building2,
  FileSpreadsheet
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()
const journalsStore = useJournalsStore()
const transactionsStore = useTransactionsStore()
const reportsStore = useReportsStore()

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)
const showManualModal = ref(false)

function inspectJournal(j: JournalEntry) {
  selectedJournal.value = j
  showDetailModal.value = true
}

// Financial KPI summaries from reports
const summary = computed(() => reportsStore.profitLossSummary)
const balance = computed(() => reportsStore.balanceSheet)

// Recent posted journals
const recentJournals = computed(() => {
  return journalsStore.journals.slice(0, 5)
})

// Unprocessed transactions preview
const unprocessedTransactions = computed(() => {
  return transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').slice(0, 4)
})
</script>

<template>
  <div class="space-y-6">
    <!-- Top Greeting & Header Bar -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Accounting Command Center
          </span>
          <span class="text-xs text-slate-400">&bull; Periode: Maret 2026</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Ikhtisar Keuangan & Jurnal Akuntansi
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Monitoring terintegrasi antara transaksi operasional, posting buku besar, laba rugi, dan perpajakan klinik.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="showManualModal = true"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Plus class="w-4 h-4" />
          <span>Buat Jurnal Manual</span>
        </button>

        <button
          @click="router.push('/reports/profit-loss')"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-xs font-semibold border border-slate-200 transition-colors"
        >
          <FileSpreadsheet class="w-4 h-4 text-slate-600" />
          <span>Laporan Laba Rugi</span>
        </button>
      </div>
    </div>

    <!-- 5 Main Accounting KPI Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
      <KpiCard
        title="Pendapatan Bruto"
        :amount="accountingStore.businessUnit === 'KIDINEST' ? summary.kidinestRev : accountingStore.businessUnit === 'SSC' ? summary.sscRev : summary.totalRev"
        trend="up"
        trend-value="+14.2%"
        subtitle="Maret 2026"
        icon-color="blue"
      >
        <template #icon>
          <TrendingUp class="w-5 h-5" />
        </template>
      </KpiCard>

      <KpiCard
        title="Laba Bersih (Net)"
        :amount="accountingStore.businessUnit === 'KIDINEST' ? summary.kidinestNetProfit : accountingStore.businessUnit === 'SSC' ? summary.sscNetProfit : summary.netProfit"
        trend="up"
        trend-value="28.4% Margin"
        subtitle="Net Profit Margin"
        icon-color="emerald"
      >
        <template #icon>
          <Sparkles class="w-5 h-5" />
        </template>
      </KpiCard>

      <KpiCard
        title="Kas & Bank Likuid"
        :amount="balance.kasOperasional + balance.bankOperasional"
        subtitle="BCA & Mandiri Operasional"
        icon-color="indigo"
      >
        <template #icon>
          <Wallet class="w-5 h-5" />
        </template>
      </KpiCard>

      <KpiCard
        title="Piutang Pasien & Asuransi"
        :amount="balance.piutangPasien"
        trend="neutral"
        trend-value="Sinarmas & Prudential"
        subtitle="AR Aging Terkontrol"
        icon-color="amber"
      >
        <template #icon>
          <CreditCard class="w-5 h-5" />
        </template>
      </KpiCard>

      <KpiCard
        title="Hutang Distributor & Fee"
        :amount="balance.hutangDistributor + balance.hutangDokter"
        subtitle="Kimia Farma & Jasa Dokter"
        icon-color="rose"
      >
        <template #icon>
          <Receipt class="w-5 h-5" />
        </template>
      </KpiCard>
    </div>

    <!-- Dual Clinic Unit Comparison & Workflow Health Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Unit Performance Card (Kidinest vs SSC) -->
      <div class="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <Building2 class="w-4 h-4 text-blue-600" />
            <h3 class="text-sm font-bold text-slate-900">Perbandingan Unit Bisnis</h3>
          </div>
          <span class="text-[11px] font-semibold text-slate-500">Maret 2026</span>
        </div>

        <!-- Kidinest -->
        <div class="p-3.5 rounded-xl bg-blue-50/60 border border-blue-100 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-blue-600"></span>
              <span class="font-bold text-xs text-blue-950">Kidinest (Klinik Anak)</span>
            </div>
            <span class="font-mono text-xs font-bold text-blue-900">{{ formatIDR(summary.kidinestRev) }}</span>
          </div>
          <div class="w-full bg-blue-200/60 h-2 rounded-full overflow-hidden">
            <div
              class="bg-blue-600 h-full rounded-full"
              :style="{ width: `${(summary.kidinestRev / summary.totalRev) * 100}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-blue-800">
            <span>Laba Bersih: <strong>{{ formatIDR(summary.kidinestNetProfit) }}</strong></span>
            <span>Kontribusi: <strong>{{ Math.round((summary.kidinestRev / summary.totalRev) * 100) }}%</strong></span>
          </div>
        </div>

        <!-- SSC -->
        <div class="p-3.5 rounded-xl bg-rose-50/60 border border-rose-100 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-rose-600"></span>
              <span class="font-bold text-xs text-rose-950">SSC (Skin & Slimming)</span>
            </div>
            <span class="font-mono text-xs font-bold text-rose-900">{{ formatIDR(summary.sscRev) }}</span>
          </div>
          <div class="w-full bg-rose-200/60 h-2 rounded-full overflow-hidden">
            <div
              class="bg-rose-600 h-full rounded-full"
              :style="{ width: `${(summary.sscRev / summary.totalRev) * 100}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-rose-800">
            <span>Laba Bersih: <strong>{{ formatIDR(summary.sscNetProfit) }}</strong></span>
            <span>Kontribusi: <strong>{{ Math.round((summary.sscRev / summary.totalRev) * 100) }}%</strong></span>
          </div>
        </div>

        <p class="text-[11px] text-slate-500 leading-relaxed italic">
          * Catatan akuntansi: Pendapatan SSC mencakup tindakan medis laser spesialis kulit & penjualan produk skincare estetika.
        </p>
      </div>

      <!-- Accounting Workflow Status Pipeline -->
      <div class="lg:col-span-2 bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <FileCheck2 class="w-4 h-4 text-emerald-600" />
            <h3 class="text-sm font-bold text-slate-900">Alur Pemrosesan Jurnal Akuntansi</h3>
          </div>
          <span class="text-xs text-emerald-600 font-semibold bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
            Audit Ready
          </span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div
            @click="router.push('/transactions')"
            class="p-3.5 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-slate-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-slate-500 uppercase">1. Transaksi Sumber</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold font-mono text-slate-900">
                {{ transactionsStore.transactions.length }}
              </span>
              <span class="text-[10px] font-bold text-amber-600 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200">
                {{ transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').length }} Belum
              </span>
            </div>
            <p class="text-[11px] text-slate-500">Kasir, PO, Farmasi, Fee</p>
          </div>

          <div
            @click="router.push('/journals?status=DRAFT')"
            class="p-3.5 rounded-xl border border-slate-200 bg-slate-50/70 hover:bg-slate-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-slate-500 uppercase">2. Draft Jurnal</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold font-mono text-slate-900">
                {{ journalsStore.stats.draftCount + journalsStore.stats.submittedCount }}
              </span>
              <span class="text-[10px] font-bold text-blue-600 bg-blue-50 px-1.5 py-0.5 rounded border border-blue-200">
                Maker
              </span>
            </div>
            <p class="text-[11px] text-slate-500">Staf Akuntansi</p>
          </div>

          <div
            @click="router.push('/journals/approvals')"
            class="p-3.5 rounded-xl border border-sky-200 bg-sky-50/50 hover:bg-sky-100/60 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-sky-700 uppercase">3. Review & Approval</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold font-mono text-sky-900">
                {{ journalsStore.stats.reviewedCount + journalsStore.stats.approvedCount }}
              </span>
              <span class="text-[10px] font-bold text-sky-700 bg-sky-100 px-1.5 py-0.5 rounded">
                Perlu Aksi
              </span>
            </div>
            <p class="text-[11px] text-sky-700">Supervisor & Manager</p>
          </div>

          <div
            @click="router.push('/journals?status=POSTED')"
            class="p-3.5 rounded-xl border border-emerald-200 bg-emerald-50/50 hover:bg-emerald-100/60 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-emerald-700 uppercase">4. Posted ke Ledger</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold font-mono text-emerald-900">
                {{ journalsStore.stats.postedCount }}
              </span>
              <span class="text-[10px] font-bold text-emerald-700 bg-emerald-100 px-1.5 py-0.5 rounded">
                Resmi
              </span>
            </div>
            <p class="text-[11px] text-emerald-700">Mempengaruhi Laporan</p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-slate-50 border border-slate-200 text-xs flex items-center justify-between">
          <div class="flex items-center gap-2 text-slate-700 font-medium">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            <span>Total Debet Resmi Terposting ke Buku Besar Periode Ini:</span>
          </div>
          <span class="font-mono font-bold text-sm text-slate-900">
            {{ formatIDR(journalsStore.stats.totalPostedDebit) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Two Columns: Unprocessed Source Transactions & Recent Journals -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Unprocessed Source Transactions -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <Layers class="w-4 h-4 text-blue-600" />
            <h3 class="text-sm font-bold text-slate-900">Transaksi Operasional Masuk</h3>
          </div>
          <button
            @click="router.push('/transactions')"
            class="text-xs font-semibold text-blue-600 hover:text-blue-800 flex items-center gap-1"
          >
            <span>Lihat Semua</span>
            <ArrowUpRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="space-y-2.5">
          <div
            v-for="tx in unprocessedTransactions"
            :key="tx.id"
            class="p-3 rounded-xl border border-slate-200/90 hover:border-slate-300 bg-slate-50/40 hover:bg-slate-50 transition-colors flex items-center justify-between"
          >
            <div class="space-y-0.5 min-w-0 pr-3">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-slate-800 text-xs">{{ tx.transactionNumber }}</span>
                <StatusBadge :status="tx.businessUnit" size="sm" />
                <span class="text-[10px] font-semibold text-slate-500 uppercase">{{ tx.sourceModule }}</span>
              </div>
              <p class="text-xs text-slate-700 truncate font-medium">{{ tx.description }}</p>
              <p class="text-[11px] text-slate-400">{{ tx.transactionDate }} &bull; {{ tx.branchName }}</p>
            </div>

            <div class="text-right shrink-0">
              <div class="font-mono font-bold text-xs text-slate-900">{{ formatIDR(tx.amount) }}</div>
              <button
                @click="transactionsStore.generateDraftJournal(tx.id)"
                class="mt-1 px-2.5 py-1 bg-blue-50 text-blue-700 hover:bg-blue-600 hover:text-white rounded text-[11px] font-semibold border border-blue-200 transition-colors flex items-center gap-1"
              >
                <Sparkles class="w-3 h-3" />
                <span>Buat Draft</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Journal Entries List -->
      <div class="bg-white rounded-2xl border border-slate-200/80 p-5 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <Receipt class="w-4 h-4 text-blue-600" />
            <h3 class="text-sm font-bold text-slate-900">Jurnal Akuntansi Terkini</h3>
          </div>
          <button
            @click="router.push('/journals')"
            class="text-xs font-semibold text-blue-600 hover:text-blue-800 flex items-center gap-1"
          >
            <span>Daftar Lengkap</span>
            <ArrowUpRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="space-y-2.5">
          <div
            v-for="j in recentJournals"
            :key="j.id"
            @click="inspectJournal(j)"
            class="p-3 rounded-xl border border-slate-200/90 hover:border-blue-300 bg-white hover:bg-blue-50/20 cursor-pointer transition-all flex items-center justify-between"
          >
            <div class="space-y-0.5 min-w-0 pr-3">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-blue-700 text-xs">{{ j.journalNumber }}</span>
                <StatusBadge :status="j.status" size="sm" />
                <StatusBadge :status="j.businessUnit" size="sm" />
              </div>
              <p class="text-xs text-slate-800 font-medium truncate">{{ j.description }}</p>
              <p class="text-[11px] text-slate-400">
                {{ formatDate(j.journalDate) }} &bull; {{ j.lines.length }} Pos Akun &bull; {{ j.branchName }}
              </p>
            </div>

            <div class="text-right shrink-0">
              <span class="font-mono font-bold text-xs text-slate-900">
                {{ formatIDR(j.totalDebit) }}
              </span>
              <p class="text-[10px] text-slate-400 mt-1">
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
