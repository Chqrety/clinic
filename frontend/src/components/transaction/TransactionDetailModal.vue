<script setup lang="ts">
import type { OperationalTransaction } from '@/types/accounting'
import { useTransactionsStore } from '@/stores/transactions'
import { formatIDR } from '@/utils/currency'
import StatusBadge from '@/components/common/StatusBadge.vue'
import { X, ArrowRight, Sparkles, CheckCircle2 } from 'lucide-vue-next'

const props = defineProps<{
  transaction: OperationalTransaction | null
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const transactionsStore = useTransactionsStore()

function handleGenerateDraft() {
  if (props.transaction) {
    transactionsStore.generateDraftJournal(props.transaction.id)
    emit('close')
  }
}
</script>

<template>
  <div
    v-if="show && transaction"
    class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-3xl w-full overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div class="px-6 py-4 bg-slate-900 text-white flex items-center justify-between">
        <div>
          <div class="flex items-center gap-2">
            <span class="font-mono font-bold text-base text-blue-400">
              {{ transaction.transactionNumber }}
            </span>
            <StatusBadge :status="transaction.businessUnit" size="sm" />
            <StatusBadge :status="transaction.mappingStatus" size="sm" />
          </div>
          <p class="text-xs text-slate-400 mt-0.5">
            Sumber Sistem: <strong class="text-slate-200">{{ transaction.sourceModule }}</strong> &bull; Ref: {{ transaction.referenceNumber }}
          </p>
        </div>
        <button
          @click="emit('close')"
          class="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto space-y-5 text-xs text-slate-700">
        <!-- Overview Grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 bg-slate-50 border border-slate-200/80 rounded-xl">
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Waktu Transaksi</span>
            <p class="font-bold text-slate-900 mt-0.5">{{ transaction.transactionDate }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Cabang Pelayanan</span>
            <p class="font-bold text-slate-900 mt-0.5">{{ transaction.branchName }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Nominal Bruto</span>
            <p class="font-bold text-slate-900 font-mono text-sm mt-0.5">{{ formatIDR(transaction.amount) }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Status Akuntansi</span>
            <p class="font-bold text-slate-900 mt-0.5">
              <span
                v-if="transaction.journalStatus === 'DRAFT_GENERATED'"
                class="text-emerald-600 font-semibold flex items-center gap-1"
              >
                <CheckCircle2 class="w-3.5 h-3.5" /> Draft {{ transaction.generatedJournalId }}
              </span>
              <span v-else class="text-amber-600 font-semibold">
                Menunggu Pemrosesan Jurnal
              </span>
            </p>
          </div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Deskripsi Transaksi Sumber:</label>
          <div class="p-3 bg-white border border-slate-200 rounded-lg text-slate-800 font-medium">
            {{ transaction.description }}
          </div>
        </div>

        <!-- Raw Operational Details Breakdown -->
        <div>
          <label class="block text-[11px] font-bold text-slate-500 uppercase mb-1">Data Rincian Operasional:</label>
          <div class="p-3 bg-slate-50 border border-slate-200 rounded-lg font-mono text-[11px] text-slate-800 space-y-1">
            <div v-for="(val, key) in transaction.details" :key="key" class="flex justify-between border-b border-slate-200/60 pb-1 last:border-0 last:pb-0">
              <span class="text-slate-500 capitalize">{{ String(key).replace(/([A-Z])/g, ' $1') }}:</span>
              <span class="font-bold text-slate-900">{{ typeof val === 'number' && val > 1000 ? formatIDR(val) : val }}</span>
            </div>
          </div>
        </div>

        <!-- Accounting Mapping Rule Preview -->
        <div class="p-4 bg-blue-50/60 border border-blue-200 rounded-xl space-y-2">
          <div class="flex items-center gap-1.5 text-blue-900 font-bold text-xs">
            <Sparkles class="w-4 h-4 text-blue-600" />
            <span>Aturan Pemetaan Akun Terotomasi (Accounting Rule Mapping)</span>
          </div>
          <p class="text-slate-600 text-[11px]">
            Sistem secara otomatis menghubungkan modul <strong>{{ transaction.sourceModule }}</strong> ke struktur akun Bagan Akun (COA):
          </p>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2">
            <div class="p-2.5 bg-white border border-blue-200/80 rounded-lg">
              <span class="text-[10px] font-bold text-emerald-700 uppercase block">Posisi Debit:</span>
              <span class="font-bold text-slate-800 text-xs">{{ transaction.suggestedDebitAccount }}</span>
            </div>
            <div class="p-2.5 bg-white border border-blue-200/80 rounded-lg">
              <span class="text-[10px] font-bold text-blue-700 uppercase block">Posisi Kredit:</span>
              <span class="font-bold text-slate-800 text-xs">{{ transaction.suggestedCreditAccount }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-200 flex items-center justify-between">
        <button
          @click="emit('close')"
          type="button"
          class="px-4 py-2 bg-white border border-slate-300 hover:bg-slate-100 text-slate-700 rounded-lg text-xs font-semibold"
        >
          Tutup
        </button>

        <button
          v-if="transaction.journalStatus === 'UNPROCESSED'"
          @click="handleGenerateDraft"
          type="button"
          class="flex items-center gap-2 px-5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>Buat Draft Jurnal Akuntansi</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
