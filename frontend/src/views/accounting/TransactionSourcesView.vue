<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import StatusBadge from '@/components/common/StatusBadge.vue'
import TransactionDetailModal from '@/components/transaction/TransactionDetailModal.vue'
import type { OperationalTransaction, SourceModule } from '@/types/accounting'
import {
  Layers,
  Search,
  Filter,
  Sparkles,
  Eye,
  CheckCircle2,
  AlertCircle
} from 'lucide-vue-next'

const transactionsStore = useTransactionsStore()
const accountingStore = useAccountingStore()

const searchQuery = ref('')
const selectedModule = ref<string>('ALL')
const selectedStatus = ref<string>('ALL')

const selectedTx = ref<OperationalTransaction | null>(null)
const showDetailModal = ref(false)

function inspectTx(tx: OperationalTransaction) {
  selectedTx.value = tx
  showDetailModal.value = true
}

const moduleOptions: { value: string; label: string }[] = [
  { value: 'ALL', label: 'Semua Modul Sumber' },
  { value: 'GRO_CASHIER', label: 'GRO / Cashier Closing' },
  { value: 'PHARMACY', label: 'Farmasi & Resep' },
  { value: 'PURCHASING', label: 'Purchasing (PO Obat)' },
  { value: 'DOCTOR_SHARE', label: 'Jasa Medis Dokter' },
  { value: 'PAYROLL', label: 'Payroll & Gaji' },
  { value: 'BANK', label: 'Mutasi & Setoran Bank' },
  { value: 'FIXED_ASSET', label: 'Penyusutan Aset Tetap' }
]

const filteredList = computed(() => {
  return transactionsStore.filteredTransactions.filter(tx => {
    if (selectedModule.value !== 'ALL' && tx.sourceModule !== selectedModule.value) {
      return false
    }
    if (selectedStatus.value !== 'ALL' && tx.journalStatus !== selectedStatus.value) {
      return false
    }
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchNo = tx.transactionNumber.toLowerCase().includes(q)
      const matchRef = tx.referenceNumber.toLowerCase().includes(q)
      const matchDesc = tx.description.toLowerCase().includes(q)
      if (!matchNo && !matchRef && !matchDesc) return false
    }
    return true
  })
})

const stats = computed(() => {
  const all = transactionsStore.filteredTransactions
  return {
    total: all.length,
    unprocessed: all.filter(t => t.journalStatus === 'UNPROCESSED').length,
    draftGenerated: all.filter(t => t.journalStatus === 'DRAFT_GENERATED').length
  }
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Sumber Transaksi Eksternal
          </span>
          <span class="text-xs text-slate-400">&bull; Integration Hub</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Transaksi Operasional Masuk
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Data transaksi dari kasir, farmasi, purchasing, penggajian, dan jasa dokter untuk divalidasi dan diolah menjadi draft jurnal akuntansi.
        </p>
      </div>

      <!-- Quick Stats Pills -->
      <div class="flex items-center gap-3">
        <div class="px-3.5 py-2 rounded-xl bg-amber-50 border border-amber-200 text-amber-800 text-xs">
          <span class="block text-[10px] uppercase font-bold text-amber-600">Menunggu Proses:</span>
          <span class="font-bold text-base font-mono">{{ stats.unprocessed }} Transaksi</span>
        </div>
        <div class="px-3.5 py-2 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs">
          <span class="block text-[10px] uppercase font-bold text-emerald-600">Draft Diterbitkan:</span>
          <span class="font-bold text-base font-mono">{{ stats.draftGenerated }} Jurnal</span>
        </div>
      </div>
    </div>

    <!-- Filters & Search Toolbar -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search Input -->
        <div class="relative w-full md:w-72">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari no transaksi, referensi, memo..."
            class="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <!-- Module Selector -->
        <div class="relative">
          <select
            v-model="selectedModule"
            class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium text-slate-800 focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer"
          >
            <option
              v-for="opt in moduleOptions"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
        </div>

        <!-- Status Filter -->
        <div class="relative">
          <select
            v-model="selectedStatus"
            class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium text-slate-800 focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer"
          >
            <option value="ALL">Semua Status Pemrosesan</option>
            <option value="UNPROCESSED">Belum Diproses</option>
            <option value="DRAFT_GENERATED">Draft Diterbitkan</option>
          </select>
        </div>
      </div>

      <div class="text-xs text-slate-500 font-medium">
        Menampilkan <strong>{{ filteredList.length }}</strong> transaksi
      </div>
    </div>

    <!-- Data Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">No Transaksi</th>
              <th class="px-4 py-3">Waktu</th>
              <th class="px-4 py-3">Modul Sumber</th>
              <th class="px-4 py-3">Unit</th>
              <th class="px-4 py-3">Deskripsi Transaksi</th>
              <th class="px-4 py-3 text-right">Nominal (Rp)</th>
              <th class="px-4 py-3 text-center">Status Pemetaan</th>
              <th class="px-4 py-3 text-center">Status Jurnal</th>
              <th class="px-4 py-3 text-center w-36">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="tx in filteredList"
              :key="tx.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 font-mono font-bold text-slate-900">
                {{ tx.transactionNumber }}
                <span class="block text-[10px] font-normal text-slate-400">Ref: {{ tx.referenceNumber }}</span>
              </td>
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ tx.transactionDate }}
              </td>
              <td class="px-4 py-3.5">
                <span class="px-2 py-0.5 rounded font-bold text-[10px] bg-slate-100 text-slate-700 border border-slate-200">
                  {{ tx.sourceModule }}
                </span>
              </td>
              <td class="px-4 py-3.5">
                <StatusBadge :status="tx.businessUnit" size="sm" />
              </td>
              <td class="px-4 py-3.5 text-slate-800 max-w-xs truncate">
                {{ tx.description }}
                <span class="block text-[10px] text-slate-400 font-mono">
                  Dr: {{ tx.suggestedDebitAccount.split(' ')[0] }} &bull; Cr: {{ tx.suggestedCreditAccount.split(' ')[0] }}
                </span>
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-slate-900 whitespace-nowrap">
                {{ formatIDR(tx.amount) }}
              </td>
              <td class="px-4 py-3.5 text-center">
                <StatusBadge :status="tx.mappingStatus" size="sm" />
              </td>
              <td class="px-4 py-3.5 text-center">
                <span
                  v-if="tx.journalStatus === 'DRAFT_GENERATED'"
                  class="inline-flex items-center gap-1 text-[11px] font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full border border-emerald-200"
                >
                  <CheckCircle2 class="w-3 h-3" /> {{ tx.generatedJournalId }}
                </span>
                <span
                  v-else
                  class="inline-flex items-center gap-1 text-[11px] font-bold text-amber-700 bg-amber-50 px-2 py-0.5 rounded-full border border-amber-200"
                >
                  <AlertCircle class="w-3 h-3" /> Menunggu Draft
                </span>
              </td>
              <td class="px-4 py-3.5 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    @click="inspectTx(tx)"
                    class="p-1.5 rounded-lg text-slate-500 hover:text-blue-600 hover:bg-blue-50 transition-colors"
                    title="Lihat Rincian Transaksi"
                  >
                    <Eye class="w-4 h-4" />
                  </button>

                  <button
                    v-if="tx.journalStatus === 'UNPROCESSED'"
                    @click="transactionsStore.generateDraftJournal(tx.id)"
                    class="flex items-center gap-1 px-2.5 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-[11px] font-bold shadow-2xs transition-colors"
                    title="Generate Draft Jurnal Akuntansi"
                  >
                    <Sparkles class="w-3 h-3" />
                    <span>Buat Draft</span>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredList.length === 0">
              <td colspan="9" class="px-4 py-8 text-center text-slate-400">
                Tidak ada transaksi yang cocok dengan kriteria pencarian atau filter.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Detail Modal -->
    <TransactionDetailModal
      :show="showDetailModal"
      :transaction="selectedTx"
      @close="showDetailModal = false; selectedTx = null"
    />
  </div>
</template>

