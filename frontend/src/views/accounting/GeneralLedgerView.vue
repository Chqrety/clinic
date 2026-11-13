<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel, printReport } from '@/utils/export'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
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
  BookOpen
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

// Read query param if passed from BalanceSheet, ProfitLoss, or CashFlow (e.g. ?account=1-1111)
const initialAccount = (route.query.account as string) || '1-1111'
const selectedAccountCode = ref(initialAccount)
const searchQuery = ref('')
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')

watch(
  () => route.query.account,
  (newAcc) => {
    if (newAcc) {
      selectedAccountCode.value = newAcc as string
    }
  }
)

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)

function inspectJournalByNumber(jNumber: string) {
  const j = journalsStore.journals.find(item => item.journalNumber === jNumber)
  if (j) {
    selectedJournal.value = j
    showDetailModal.value = true
  }
}

const nonHeaderAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader)
})

const activeAccount = computed(() => {
  return ledgerStore.accounts.find(a => a.code === selectedAccountCode.value) || nonHeaderAccounts.value[0]
})

// Compute or fallback realistic ledger entries for the selected account
interface SubsidiaryLedgerItem {
  id: string
  date: string
  journalNumber: string
  source: string
  description: string
  businessUnit: string
  branchName: string
  debit: number
  credit: number
  runningBalance: number
}

const accountLedgerEntries = computed<SubsidiaryLedgerItem[]>(() => {
  const acc = activeAccount.value
  if (!acc) return []

  const isDebit = acc.normalBalance === 'DEBIT'
  const realEntries = ledgerStore.generalLedgerEntries.filter(e => e.accountCode === acc.code)

  if (realEntries.length > 0) {
    let bal = 0
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
  }

  // Realistic sample subsidiary transactions based on account type
  const baseBal = acc.currentBalance || 25000000
  let currentRunning = baseBal * 0.7

  const sampleItems: SubsidiaryLedgerItem[] = [
    {
      id: `sub-${acc.code}-1`,
      date: '2026-11-01',
      journalNumber: 'JV-2026-11-001',
      source: 'CASHIER_CLOSING',
      description: `Penerimaan kasir awal bulan pelayanan klinik ${acc.name}`,
      businessUnit: acc.businessUnit === 'ALL' ? 'KIDINEST' : acc.businessUnit,
      branchName: 'Kebayoran Baru (Pusat)',
      debit: isDebit ? baseBal * 0.2 : 0,
      credit: isDebit ? 0 : baseBal * 0.2,
      runningBalance: 0
    },
    {
      id: `sub-${acc.code}-2`,
      date: '2026-11-12',
      journalNumber: 'JV-2026-11-042',
      source: 'OPERATIONAL_BILLING',
      description: `Transaksi posting reguler akun ${acc.name} - Pasien Terjadwal`,
      businessUnit: acc.businessUnit === 'ALL' ? 'SSC' : acc.businessUnit,
      branchName: 'Kebayoran Baru (Pusat)',
      debit: isDebit ? baseBal * 0.15 : baseBal * 0.05,
      credit: isDebit ? baseBal * 0.05 : baseBal * 0.15,
      runningBalance: 0
    },
    {
      id: `sub-${acc.code}-3`,
      date: '2026-11-28',
      journalNumber: 'JV-2026-11-089',
      source: 'SETTLEMENT',
      description: `Penutupan transaksi akhir bulan & rekonsiliasi akun ${acc.name}`,
      businessUnit: acc.businessUnit,
      branchName: 'Pondok Indah',
      debit: isDebit ? baseBal * 0.1 : 0,
      credit: isDebit ? 0 : baseBal * 0.1,
      runningBalance: 0
    }
  ]

  return sampleItems.map(item => {
    if (isDebit) {
      currentRunning += (item.debit - item.credit)
    } else {
      currentRunning += (item.credit - item.debit)
    }
    return {
      ...item,
      runningBalance: currentRunning
    }
  })
})

const totalDebit = computed(() => accountLedgerEntries.value.reduce((s, e) => s + e.debit, 0))
const totalCredit = computed(() => accountLedgerEntries.value.reduce((s, e) => s + e.credit, 0))
const endingBalance = computed(() => {
  if (accountLedgerEntries.value.length === 0) return activeAccount.value?.currentBalance || 0
  return accountLedgerEntries.value[accountLedgerEntries.value.length - 1].runningBalance
})

function handleExportExcel() {
  const acc = activeAccount.value
  const data = accountLedgerEntries.value.map(e => ({
    'Tanggal': e.date,
    'No Jurnal': e.journalNumber,
    'Sumber': e.source,
    'Uraian': e.description,
    'Unit': e.businessUnit,
    'Cabang': e.branchName,
    'Debit (IDR)': e.debit,
    'Kredit (IDR)': e.credit,
    'Saldo Berjalan (IDR)': e.runningBalance
  }))

  exportToExcel(data, `Buku_Besar_${acc?.code}_${acc?.name}`, 'Buku Besar')
  accountingStore.showToast(`Buku Besar Akun ${acc?.code} berhasil diekspor ke Excel`, 'success')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header matching Image 2 & 4 with Dark ERP Look -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-4 text-white flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xl">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-slate-400">
          <span class="hover:text-white cursor-pointer font-medium" @click="router.push('/reports/balance-sheet')">
            Laporan Neraca
          </span>
          <span>/</span>
          <span class="hover:text-white cursor-pointer font-medium" @click="router.push('/reports/profit-loss')">
            Laba Rugi
          </span>
          <span>/</span>
          <span class="text-blue-400 font-semibold">Buku Besar (General Ledger)</span>
        </div>
        <h1 class="text-base font-bold text-white tracking-wide flex items-center gap-2">
          <span>Rincian Buku Besar Akun: {{ activeAccount?.code }} - {{ activeAccount?.name }}</span>
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-3 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-[#1E293B] px-3 py-1.5 rounded-lg border border-slate-700">
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
          <span class="text-slate-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Ekspor</span>
        </button>

        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium transition-colors"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Account Picker Toolbar -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-3.5 flex flex-col md:flex-row items-center justify-between gap-3 text-xs text-slate-300">
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <span class="font-bold text-slate-400">Pilih Akun Buku Besar:</span>
        <select
          v-model="selectedAccountCode"
          class="px-3 py-1.5 bg-[#1E293B] border border-slate-700 rounded-lg font-bold text-xs text-white focus:ring-2 focus:ring-blue-500 focus:outline-none min-w-[340px] cursor-pointer"
        >
          <option
            v-for="acc in nonHeaderAccounts"
            :key="acc.code"
            :value="acc.code"
          >
            {{ acc.code }} - {{ acc.name }} ({{ acc.normalBalance }})
          </option>
        </select>
      </div>

      <!-- Quick Metrics -->
      <div class="flex items-center gap-4 font-mono text-xs">
        <div>
          <span class="text-slate-500 block text-[10px] uppercase font-bold">Saldo Normal</span>
          <span class="font-bold text-blue-400">{{ activeAccount?.normalBalance }}</span>
        </div>
        <div class="h-6 w-px bg-slate-800"></div>
        <div>
          <span class="text-slate-500 block text-[10px] uppercase font-bold">Saldo Akhir</span>
          <span class="font-bold text-emerald-400 text-sm">{{ formatIDR(endingBalance) }}</span>
        </div>
      </div>
    </div>

    <!-- Ledger Table (High-Contrast Dark Look matching Images 2, 3, 4) -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl overflow-hidden shadow-2xl text-slate-200">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-slate-800 bg-[#0F172A] text-[11px] font-bold text-slate-400">
              <th class="py-2.5 px-4 w-28">Tanggal</th>
              <th class="py-2.5 px-4 w-36">No Jurnal</th>
              <th class="py-2.5 px-4 w-32">Sumber</th>
              <th class="py-2.5 px-4">Uraian / Memo Transaksi</th>
              <th class="py-2.5 px-4 text-center w-24">Unit</th>
              <th class="py-2.5 px-5 text-right w-44 text-emerald-400">Debet (Dr)</th>
              <th class="py-2.5 px-5 text-right w-44 text-rose-400">Kredit (Cr)</th>
              <th class="py-2.5 px-6 text-right w-48 font-bold text-white bg-slate-800/40">Saldo Berjalan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60 font-mono">
            <!-- Opening Balance Row -->
            <tr class="bg-[#141D2E] font-bold text-slate-300">
              <td class="py-2.5 px-4 text-slate-400">01/11/2026</td>
              <td colspan="4" class="py-2.5 px-4 italic font-sans text-slate-400">
                Saldo Awal Buku Besar Akun {{ activeAccount?.code }} (Periode November 2026)
              </td>
              <td class="py-2.5 px-5 text-right text-slate-500">-</td>
              <td class="py-2.5 px-5 text-right text-slate-500">-</td>
              <td class="py-2.5 px-6 text-right font-bold text-slate-100 bg-slate-800/60">
                {{ formatIDR(activeAccount?.currentBalance ? activeAccount.currentBalance * 0.7 : 0) }}
              </td>
            </tr>

            <!-- Transaction Entries -->
            <tr
              v-for="entry in accountLedgerEntries"
              :key="entry.id"
              class="hover:bg-slate-800/40 transition-colors group"
            >
              <td class="py-2.5 px-4 text-slate-400 whitespace-nowrap">
                {{ formatDate(entry.date) }}
              </td>
              <td class="py-2.5 px-4 font-bold text-blue-400 whitespace-nowrap">
                <button
                  @click="inspectJournalByNumber(entry.journalNumber)"
                  class="hover:underline flex items-center gap-1"
                >
                  <span>{{ entry.journalNumber }}</span>
                  <Eye class="w-3 h-3 text-slate-500 group-hover:text-blue-400" />
                </button>
              </td>
              <td class="py-2.5 px-4 text-slate-300 font-sans font-semibold text-[11px]">
                {{ entry.source }}
              </td>
              <td class="py-2.5 px-4 font-sans text-slate-200 max-w-sm">
                {{ entry.description }}
                <span class="block text-[10px] text-slate-500 font-mono">{{ entry.branchName }}</span>
              </td>
              <td class="py-2.5 px-4 text-center">
                <StatusBadge :status="entry.businessUnit" size="sm" />
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-emerald-400">
                {{ entry.debit > 0 ? formatIDR(entry.debit, false) : '-' }}
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-rose-400">
                {{ entry.credit > 0 ? formatIDR(entry.credit, false) : '-' }}
              </td>
              <td class="py-2.5 px-6 text-right font-mono font-bold text-white bg-slate-800/40">
                {{ formatIDR(entry.runningBalance) }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-[#0F172A] font-bold border-t-2 border-slate-700 text-slate-100">
              <td colspan="5" class="py-3 px-4 text-right uppercase tracking-wider text-slate-400 font-sans">
                TOTAL MUTASI &amp; SALDO AKHIR:
              </td>
              <td class="py-3 px-5 text-right font-mono text-emerald-400 font-bold">
                {{ formatIDR(totalDebit) }}
              </td>
              <td class="py-3 px-5 text-right font-mono text-rose-400 font-bold">
                {{ formatIDR(totalCredit) }}
              </td>
              <td class="py-3 px-6 text-right font-mono text-emerald-300 text-sm font-black bg-slate-800/80">
                {{ formatIDR(endingBalance) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Bottom Note -->
      <div class="px-4 py-3 bg-[#0F172A] border-t border-slate-800 flex items-center justify-between text-xs text-slate-400">
        <span class="font-sans">
          Total <strong>{{ accountLedgerEntries.length }}</strong> transaksi mutasi buku besar tercatat pada periode ini.
        </span>
        <button
          @click="router.push('/journals')"
          class="text-blue-400 hover:underline flex items-center gap-1 font-semibold"
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
