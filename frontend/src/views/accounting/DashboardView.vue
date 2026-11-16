<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useTransactionsStore } from '@/stores/transactions'
import { useReportsStore } from '@/stores/reports'
import { formatIDR, compactIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
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
            Akuntansi & Keuangan
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Ikhtisar Finansial & Jurnal Akuntansi
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="showManualModal = true"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Buat Jurnal Manual</span>
        </button>

        <button
          @click="router.push('/accounting/reports/profit-loss')"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 text-zinc-600" />
          <span>Laporan Laba Rugi</span>
        </button>
      </div>
    </header>

    <!-- 1-Row Compact Hairline KPI Ticker Strip (Anti AI-Slop, Calm & Harmonious) -->
    <div
      class="grid grid-cols-2 sm:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Pendapatan Bruto</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(accountingStore.businessUnit === 'KIDINEST' ? summary.kidinestRev : accountingStore.businessUnit === 'SSC' ? summary.sscRev : summary.totalRev) }}
        </div>
        <div class="text-[10px] text-emerald-600 font-semibold mt-0.5">+14.2% vs Feb 2026</div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Kas & Bank Likuid</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(balance.kasOperasional + balance.bankOperasional) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">BCA & Kasir Operasional</div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Piutang Pasien / AR</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(balance.piutangPasien) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">Asuransi & Pasien Umum</div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Hutang & Fee Dokter</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(balance.hutangDistributor + balance.hutangDokter) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">Distributor & Jasa Medis</div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex flex-col justify-center">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Laba Bersih (Net)</div>
        <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
          {{ formatIDR(accountingStore.businessUnit === 'KIDINEST' ? summary.kidinestNetProfit : accountingStore.businessUnit === 'SSC' ? summary.sscNetProfit : summary.netProfit) }}
        </div>
        <div class="text-[10px] text-emerald-700 font-bold mt-0.5">28.4% Net Margin</div>
      </div>
    </div>

    <!-- Dual Clinic Unit Comparison & Workflow Health Section -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
      <!-- Unit Performance Card (Kidinest vs SSC) -->
      <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <Building2 class="w-4 h-4 text-zinc-700" />
            <h3 class="text-sm font-bold text-zinc-950">Perbandingan Unit Bisnis</h3>
          </div>
          <span class="text-[11px] font-semibold text-zinc-500">Maret 2026</span>
        </div>

        <!-- Kidinest (Warm Mustard / Amber) -->
        <div class="p-3.5 rounded-xl bg-amber-50/70 border border-amber-200/80 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
              <span class="font-bold text-xs text-amber-950">Kidinest Pediatric Care</span>
            </div>
            <span class="text-xs font-bold tabular-nums text-amber-950">{{ formatIDR(summary.kidinestRev) }}</span>
          </div>
          <div class="w-full bg-amber-200/60 h-1.5 rounded-full overflow-hidden">
            <div
              class="bg-amber-500 h-full rounded-full"
              :style="{ width: `${(summary.kidinestRev / summary.totalRev) * 100}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-amber-900 font-medium">
            <span>Laba Bersih: <strong class="tabular-nums">{{ formatIDR(summary.kidinestNetProfit) }}</strong></span>
            <span>Kontribusi: <strong>{{ Math.round((summary.kidinestRev / summary.totalRev) * 100) }}%</strong></span>
          </div>
        </div>

        <!-- SSC (Dusty Rose / Mauve) -->
        <div class="p-3.5 rounded-xl bg-rose-50/70 border border-rose-200/80 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-rose-500"></span>
              <span class="font-bold text-xs text-rose-950">Semarang Skin Center (SSC)</span>
            </div>
            <span class="text-xs font-bold tabular-nums text-rose-950">{{ formatIDR(summary.sscRev) }}</span>
          </div>
          <div class="w-full bg-rose-200/60 h-1.5 rounded-full overflow-hidden">
            <div
              class="bg-rose-500 h-full rounded-full"
              :style="{ width: `${(summary.sscRev / summary.totalRev) * 100}%` }"
            ></div>
          </div>
          <div class="flex justify-between text-[11px] text-rose-900 font-medium">
            <span>Laba Bersih: <strong class="tabular-nums">{{ formatIDR(summary.sscNetProfit) }}</strong></span>
            <span>Kontribusi: <strong>{{ Math.round((summary.sscRev / summary.totalRev) * 100) }}%</strong></span>
          </div>
        </div>

        <p class="text-[11px] text-zinc-500 leading-relaxed italic">
          * Catatan akuntansi: Pendapatan SSC mencakup tindakan medis laser spesialis kulit & penjualan produk skincare estetika.
        </p>
      </div>

      <!-- Accounting Workflow Status Pipeline -->
      <div class="lg:col-span-2 bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2">
            <FileCheck2 class="w-4 h-4 text-emerald-600" />
            <h3 class="text-sm font-bold text-zinc-950">Alur Pemrosesan Jurnal Akuntansi</h3>
          </div>
          <span class="text-xs text-emerald-700 font-bold bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
            Audit Ready
          </span>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div
            @click="router.push('/accounting/transactions')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-500 uppercase">1. Transaksi Sumber</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ transactionsStore.transactions.length }}
              </span>
              <span class="text-[10px] font-bold text-amber-700 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200">
                {{ transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').length }} Belum
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">Kasir, PO, Farmasi, Fee</p>
          </div>

          <div
            @click="router.push('/accounting/journals?status=DRAFT')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-500 uppercase">2. Draft Jurnal</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ journalsStore.stats.draftCount + journalsStore.stats.submittedCount }}
              </span>
              <span class="text-[10px] font-bold text-zinc-700 bg-zinc-200/80 px-1.5 py-0.5 rounded">
                Maker
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">Staf Akuntansi</p>
          </div>

          <div
            @click="router.push('/accounting/journals/approvals')"
            class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 hover:bg-zinc-100 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-zinc-600 uppercase">3. Review & Approval</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-zinc-900">
                {{ journalsStore.stats.reviewedCount + journalsStore.stats.approvedCount }}
              </span>
              <span class="text-[10px] font-bold text-amber-700 bg-amber-100 px-1.5 py-0.5 rounded">
                Perlu Aksi
              </span>
            </div>
            <p class="text-[11px] text-zinc-500">Supervisor & Manager</p>
          </div>

          <div
            @click="router.push('/accounting/journals?status=POSTED')"
            class="p-3.5 rounded-xl border border-emerald-200 bg-emerald-50/50 hover:bg-emerald-100/60 cursor-pointer transition-colors space-y-1"
          >
            <span class="text-[10px] font-bold text-emerald-800 uppercase">4. Posted ke Ledger</span>
            <div class="flex items-baseline justify-between">
              <span class="text-xl font-bold tabular-nums text-emerald-950">
                {{ journalsStore.stats.postedCount }}
              </span>
              <span class="text-[10px] font-bold text-emerald-700 bg-emerald-100 px-1.5 py-0.5 rounded">
                Resmi
              </span>
            </div>
            <p class="text-[11px] text-emerald-700">Mempengaruhi Laporan</p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200/90 text-xs flex items-center justify-between">
          <div class="flex items-center gap-2 text-zinc-700 font-medium">
            <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
            <span>Total Debet Resmi Terposting ke Buku Besar Periode Ini:</span>
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
            <h3 class="text-sm font-bold text-zinc-950">Transaksi Operasional Masuk</h3>
          </div>
          <button
            @click="router.push('/accounting/transactions')"
            class="text-xs font-bold text-zinc-800 hover:text-zinc-950 flex items-center gap-1 cursor-pointer"
          >
            <span>Lihat Semua</span>
            <ArrowUpRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <div class="space-y-2.5">
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
                <span>Buat Draft</span>
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
            <h3 class="text-sm font-bold text-zinc-950">Jurnal Akuntansi Terkini</h3>
          </div>
          <button
            @click="router.push('/accounting/journals')"
            class="text-xs font-bold text-zinc-800 hover:text-zinc-950 flex items-center gap-1 cursor-pointer"
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

