<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel } from '@/utils/export'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  FileText,
  Search,
  Download,
  Calendar,
  Eye,
  Building,
  Scale
} from 'lucide-vue-next'

const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

// Default account: Kas Operasional Kidinest
const selectedAccountCode = ref('1-1100')
const searchQuery = ref('')

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)

function inspectJournalByNumber(jNumber: string) {
  const j = journalsStore.journals.find(item => item.journalNumber === jNumber)
  if (j) {
    selectedJournal.value = j
    showDetailModal.value = true
  }
}

const activeLedgerData = computed(() => {
  return ledgerStore.getAccountLedger(selectedAccountCode.value)
})

const nonHeaderAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader)
})

function handleExportExcel() {
  const acc = activeLedgerData.value.account
  const data = activeLedgerData.value.entries.map(e => ({
    'Tanggal': e.date,
    'No Jurnal': e.journalNumber,
    'Keterangan': e.description,
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
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Buku Besar Utama (General Ledger)
          </span>
          <span class="text-xs text-slate-400">&bull; Posted Transactions Only</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Rincian Mutasi Buku Besar
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Hanya jurnal berstatus <strong>POSTED</strong> yang masuk ke dalam buku besar dan membentuk saldo resmi akun keuangan.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-semibold border border-slate-300 shadow-2xs transition-colors"
        >
          <Download class="w-4 h-4 text-slate-600" />
          <span>Ekspor Buku Besar</span>
        </button>
      </div>
    </div>

    <!-- Account Picker Toolbar -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <span class="font-bold text-slate-700">Pilih Akun Buku Besar:</span>
        <select
          v-model="selectedAccountCode"
          class="px-3.5 py-2 bg-slate-50 border border-slate-200 rounded-lg font-bold text-xs text-slate-900 focus:ring-2 focus:ring-blue-500 focus:outline-none min-w-[320px] cursor-pointer"
        >
          <option
            v-for="acc in nonHeaderAccounts"
            :key="acc.code"
            :value="acc.code"
          >
            {{ acc.code }} &bull; {{ acc.name }} ({{ acc.normalBalance }})
          </option>
        </select>
      </div>

      <!-- Quick Balance Badge -->
      <div
        v-if="activeLedgerData.account"
        class="flex items-center gap-3 bg-slate-50 px-4 py-2 rounded-xl border border-slate-200"
      >
        <span class="text-slate-500 text-xs">Saldo Akhir Akun:</span>
        <span class="font-mono font-bold text-base text-slate-900">
          {{ formatIDR(activeLedgerData.endingBalance) }}
        </span>
        <span class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-blue-100 text-blue-800">
          {{ activeLedgerData.account.normalBalance }}
        </span>
      </div>
    </div>

    <!-- Account Metadata Overview Card -->
    <div
      v-if="activeLedgerData.account"
      class="bg-slate-900 text-white p-5 rounded-2xl shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-4"
    >
      <div class="space-y-1">
        <div class="flex items-center gap-2">
          <span class="font-mono font-bold text-lg text-blue-400">{{ activeLedgerData.account.code }}</span>
          <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-800 text-slate-300 border border-slate-700">
            Tipe: {{ activeLedgerData.account.type }}
          </span>
          <StatusBadge :status="activeLedgerData.account.businessUnit" size="sm" />
        </div>
        <h2 class="text-base font-bold text-white">{{ activeLedgerData.account.name }}</h2>
        <p class="text-xs text-slate-400">
          Saldo Normal: {{ activeLedgerData.account.normalBalance }} &bull; Level: {{ activeLedgerData.account.level }}
        </p>
      </div>

      <div class="flex items-center gap-4 text-right">
        <div>
          <span class="text-[10px] font-semibold text-slate-400 uppercase">Total Mutasi Debet</span>
          <p class="font-mono font-bold text-sm text-emerald-400">
            {{ formatIDR(activeLedgerData.entries.reduce((s, e) => s + e.debit, 0)) }}
          </p>
        </div>
        <div class="h-8 w-px bg-slate-800"></div>
        <div>
          <span class="text-[10px] font-semibold text-slate-400 uppercase">Total Mutasi Kredit</span>
          <p class="font-mono font-bold text-sm text-rose-400">
            {{ formatIDR(activeLedgerData.entries.reduce((s, e) => s + e.credit, 0)) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Ledger Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">No Jurnal</th>
              <th class="px-4 py-3">Sumber Modul</th>
              <th class="px-4 py-3">Unit & Cabang</th>
              <th class="px-4 py-3">Uraian Transaksi</th>
              <th class="px-4 py-3 text-right">Debet (Rp)</th>
              <th class="px-4 py-3 text-right">Kredit (Rp)</th>
              <th class="px-4 py-3 text-right">Saldo Berjalan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <!-- Opening Balance Row -->
            <tr class="bg-slate-50/40 font-semibold text-slate-600">
              <td class="px-4 py-3 text-slate-500 font-mono">01/03/2026</td>
              <td class="px-4 py-3 italic" colspan="4">Saldo Awal Buku Besar Periode Maret 2026</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono">-</td>
              <td class="px-4 py-3 text-right font-mono font-bold text-slate-800">
                {{ formatIDR(activeLedgerData.openingBalance) }}
              </td>
            </tr>

            <!-- Entry Rows -->
            <tr
              v-for="entry in activeLedgerData.entries"
              :key="entry.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ formatDate(entry.date) }}
              </td>
              <td class="px-4 py-3.5 font-mono font-bold text-blue-700 whitespace-nowrap">
                <button
                  @click="inspectJournalByNumber(entry.journalNumber)"
                  class="hover:underline flex items-center gap-1"
                >
                  <span>{{ entry.journalNumber }}</span>
                  <Eye class="w-3 h-3 text-blue-500" />
                </button>
              </td>
              <td class="px-4 py-3.5 font-semibold text-slate-700">
                {{ entry.source }}
              </td>
              <td class="px-4 py-3.5">
                <StatusBadge :status="entry.businessUnit" size="sm" />
                <span class="block text-[10px] text-slate-400">{{ entry.branchName }}</span>
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

            <tr v-if="activeLedgerData.entries.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-slate-400">
                Tidak ada mutasi jurnal berstatus POSTED untuk akun ini pada periode berjalan.
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
              <td colspan="5" class="px-4 py-3 text-right text-slate-700">
                SALDO AKHIR (ENDING BALANCE):
              </td>
              <td class="px-4 py-3 text-right font-mono text-emerald-800">
                {{ formatIDR(activeLedgerData.entries.reduce((s, e) => s + e.debit, 0)) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-rose-800">
                {{ formatIDR(activeLedgerData.entries.reduce((s, e) => s + e.credit, 0)) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-sm text-slate-900">
                {{ formatIDR(activeLedgerData.endingBalance) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Modal Drilldown -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />
  </div>
</template>
