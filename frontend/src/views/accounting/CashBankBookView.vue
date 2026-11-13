<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  Coins,
  Building,
  ArrowDownLeft,
  ArrowUpRight,
  Wallet,
  Landmark
} from 'lucide-vue-next'

const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

const activeTab = ref<'CASH' | 'BANK'>('CASH')

// Cash accounts: 1-1100 (Kas Kidinest) and 1-1101 (Kas SSC)
// Bank accounts: 1-1200 (Bank Mandiri) and 1-1201 (Bank BCA)
const cashAccounts = ['1-1100', '1-1101']
const bankAccounts = ['1-1200', '1-1201']

const selectedAccount = ref<string>('1-1100')

const relevantAccounts = computed(() => {
  const codes = activeTab.value === 'CASH' ? cashAccounts : bankAccounts
  return ledgerStore.accounts.filter(a => codes.includes(a.code))
})

// Set initial account when switching tabs
function setTab(tab: 'CASH' | 'BANK') {
  activeTab.value = tab
  selectedAccount.value = tab === 'CASH' ? '1-1100' : '1-1201'
}

const ledgerData = computed(() => {
  return ledgerStore.getAccountLedger(selectedAccount.value)
})

// Total Cash across both clinics
const totalCash = computed(() => {
  const k1 = ledgerStore.getAccountLedger('1-1100').endingBalance
  const k2 = ledgerStore.getAccountLedger('1-1101').endingBalance
  return k1 + k2
})

// Total Bank across both accounts
const totalBank = computed(() => {
  const b1 = ledgerStore.getAccountLedger('1-1200').endingBalance
  const b2 = ledgerStore.getAccountLedger('1-1201').endingBalance
  return b1 + b2
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Treasury & Liquidity
          </span>
          <span class="text-xs text-slate-400">&bull; Real-time Cashflow</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Buku Kas & Rekening Bank
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Pemantauan buku kas kasir operasional klinik dan mutasi rekening bank Mandiri & BCA secara terintegrasi.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="px-4 py-2 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-900 text-xs">
          <span class="text-[10px] uppercase font-bold text-emerald-600 block">Total Likuiditas Aktif:</span>
          <span class="font-bold text-base font-mono">{{ formatIDR(totalCash + totalBank) }}</span>
        </div>
      </div>
    </div>

    <!-- 2 Main Cards: Kas vs Bank Overview -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Kas Card -->
      <div
        @click="setTab('CASH')"
        :class="[
          'p-5 rounded-2xl border cursor-pointer transition-all space-y-3',
          activeTab === 'CASH'
            ? 'bg-blue-900 text-white border-blue-900 shadow-md ring-2 ring-blue-500/20'
            : 'bg-white text-slate-900 border-slate-200/80 hover:border-slate-300'
        ]"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div :class="['w-9 h-9 rounded-xl flex items-center justify-center', activeTab === 'CASH' ? 'bg-blue-800 text-white' : 'bg-blue-50 text-blue-600']">
              <Wallet class="w-5 h-5" />
            </div>
            <div>
              <h3 class="font-bold text-sm">Buku Kas Operasional</h3>
              <p :class="['text-[11px]', activeTab === 'CASH' ? 'text-blue-200' : 'text-slate-500']">
                Kasir Fisik Kidinest & SSC
              </p>
            </div>
          </div>
          <span class="font-mono text-lg font-bold">
            {{ formatIDR(totalCash) }}
          </span>
        </div>
      </div>

      <!-- Bank Card -->
      <div
        @click="setTab('BANK')"
        :class="[
          'p-5 rounded-2xl border cursor-pointer transition-all space-y-3',
          activeTab === 'BANK'
            ? 'bg-blue-900 text-white border-blue-900 shadow-md ring-2 ring-blue-500/20'
            : 'bg-white text-slate-900 border-slate-200/80 hover:border-slate-300'
        ]"
      >
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div :class="['w-9 h-9 rounded-xl flex items-center justify-center', activeTab === 'BANK' ? 'bg-blue-800 text-white' : 'bg-blue-50 text-blue-600']">
              <Landmark class="w-5 h-5" />
            </div>
            <div>
              <h3 class="font-bold text-sm">Buku Rekening Bank</h3>
              <p :class="['text-[11px]', activeTab === 'BANK' ? 'text-blue-200' : 'text-slate-500']">
                Bank BCA & Mandiri Operasional
              </p>
            </div>
          </div>
          <span class="font-mono text-lg font-bold">
            {{ formatIDR(totalBank) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Sub-Account Selector -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex items-center justify-between text-xs">
      <div class="flex items-center gap-3">
        <span class="font-bold text-slate-700">Pilih Rekening / Buku Kas:</span>
        <div class="flex items-center gap-2">
          <button
            v-for="acc in relevantAccounts"
            :key="acc.code"
            @click="selectedAccount = acc.code"
            :class="[
              'px-3.5 py-1.5 rounded-lg font-semibold transition-all',
              selectedAccount === acc.code
                ? 'bg-blue-600 text-white shadow-2xs'
                : 'bg-slate-100 text-slate-700 hover:bg-slate-200'
            ]"
          >
            {{ acc.code }} - {{ acc.name }}
          </button>
        </div>
      </div>

      <div class="text-xs font-mono font-bold text-slate-900">
        Saldo Saat Ini: {{ formatIDR(ledgerData.endingBalance) }}
      </div>
    </div>

    <!-- Mutations Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">No Jurnal</th>
              <th class="px-4 py-3">Sumber Mutasi</th>
              <th class="px-4 py-3">Uraian Transaksi</th>
              <th class="px-4 py-3 text-right text-emerald-700">Penerimaan / Masuk (Dr)</th>
              <th class="px-4 py-3 text-right text-rose-700">Pengeluaran / Keluar (Cr)</th>
              <th class="px-4 py-3 text-right">Saldo Kas/Bank</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <!-- Opening -->
            <tr class="bg-slate-50/40 font-semibold text-slate-600">
              <td class="px-4 py-3 font-mono">01/03/2026</td>
              <td colspan="3" class="px-4 py-3 italic">Saldo Awal Buku Kas / Bank</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono font-bold text-slate-900">
                {{ formatIDR(ledgerData.openingBalance) }}
              </td>
            </tr>

            <!-- Transactions -->
            <tr
              v-for="entry in ledgerData.entries"
              :key="entry.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ formatDate(entry.date) }}
              </td>
              <td class="px-4 py-3.5 font-mono font-bold text-blue-700 whitespace-nowrap">
                {{ entry.journalNumber }}
              </td>
              <td class="px-4 py-3.5 font-semibold text-slate-700">
                {{ entry.source }}
              </td>
              <td class="px-4 py-3.5 max-w-sm text-slate-800">
                {{ entry.description }}
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
            </tr>

            <tr v-if="ledgerData.entries.length === 0">
              <td colspan="7" class="px-4 py-8 text-center text-slate-400">
                Belum ada transaksi kas/bank berstatus POSTED pada periode ini.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
