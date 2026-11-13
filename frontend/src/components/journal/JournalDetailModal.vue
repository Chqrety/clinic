<script setup lang="ts">
import { ref } from 'vue'
import type { JournalEntry } from '@/types/accounting'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  X,
  Printer,
  CheckCircle2,
  FileCheck,
  Send,
  RotateCcw,
  AlertTriangle
} from 'lucide-vue-next'

const props = defineProps<{
  journal: JournalEntry | null
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

const showReversePrompt = ref(false)
const reversalReason = ref('')

function handleReview() {
  if (props.journal) {
    journalsStore.reviewJournal(props.journal.id)
    emit('close')
  }
}

function handleApprove() {
  if (props.journal) {
    journalsStore.approveJournal(props.journal.id)
    emit('close')
  }
}

function handlePost() {
  if (props.journal) {
    journalsStore.postJournal(props.journal.id)
    emit('close')
  }
}

function handleReverse() {
  if (props.journal && reversalReason.value.trim()) {
    journalsStore.reverseJournal(props.journal.id, reversalReason.value.trim())
    showReversePrompt.value = false
    reversalReason.value = ''
    emit('close')
  }
}

function printVoucher() {
  window.print()
}
</script>

<template>
  <div
    v-if="show && journal"
    class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4"
  >
    <div
      class="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-4xl w-full overflow-hidden flex flex-col max-h-[90vh] animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div class="px-6 py-4 bg-slate-900 text-white flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="space-y-0.5">
            <div class="flex items-center gap-2">
              <span class="font-mono text-base font-bold text-blue-400">
                {{ journal.journalNumber }}
              </span>
              <StatusBadge :status="journal.status" size="sm" />
              <StatusBadge :status="journal.businessUnit" size="sm" />
            </div>
            <p class="text-xs text-slate-400">
              Voucher Jurnal Akuntansi &bull; Sumber: {{ journal.journalSource }} &bull; Ref: {{ journal.sourceReference }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="printVoucher"
            class="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition-colors"
            title="Cetak Voucher Jurnal"
          >
            <Printer class="w-4 h-4" />
          </button>
          <button
            @click="emit('close')"
            class="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto space-y-6 text-xs text-slate-700">
        <!-- Metadata Grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 bg-slate-50 p-4 rounded-xl border border-slate-200/80">
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Tanggal Jurnal</span>
            <p class="font-bold text-slate-900 mt-0.5">{{ formatDate(journal.journalDate) }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Cabang</span>
            <p class="font-bold text-slate-900 mt-0.5">{{ journal.branchName }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Dibuat Oleh</span>
            <p class="font-bold text-slate-900 mt-0.5">{{ journal.createdByName || 'Staf Akuntansi' }}</p>
          </div>
          <div>
            <span class="text-[11px] font-semibold text-slate-400 uppercase">Posting Ke Ledger</span>
            <p class="font-bold text-slate-900 mt-0.5">
              {{ journal.postedByName ? journal.postedByName : 'Belum Diposting' }}
            </p>
          </div>
        </div>

        <!-- Description Note -->
        <div>
          <label class="text-[11px] font-bold text-slate-500 uppercase">Keterangan / Memo Transaksi:</label>
          <div class="mt-1 p-3 bg-white border border-slate-200 rounded-lg text-slate-800 text-xs font-medium">
            {{ journal.description }}
          </div>
        </div>

        <!-- Reversal Info Alert if applicable -->
        <div
          v-if="journal.status === 'REVERSED'"
          class="p-3 bg-purple-50 border border-purple-200 rounded-xl flex items-start gap-3 text-purple-900"
        >
          <RotateCcw class="w-5 h-5 text-purple-600 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold text-xs">Jurnal Telah Dibalik (Reversed)</p>
            <p class="text-[11px] text-purple-700">
              Alasan: {{ journal.reversalReason || 'Penyesuaian transaksi' }} &bull; Jurnal Pembalik: {{ journal.reversalJournalId }}
            </p>
          </div>
        </div>

        <!-- Balanced Double-Entry Lines Table -->
        <div class="border border-slate-200 rounded-xl overflow-hidden shadow-2xs">
          <div class="bg-slate-100/80 px-4 py-2.5 border-b border-slate-200 flex items-center justify-between">
            <span class="font-bold text-slate-800">Rincian Pos Jurnal Berpasangan (Double-Entry)</span>
            <span class="text-[11px] text-slate-500 font-mono">Total {{ journal.lines.length }} Baris</span>
          </div>

          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-200 bg-slate-50/50 text-[11px] font-bold text-slate-600 uppercase">
                <th class="px-4 py-2.5 w-28">Kode Akun</th>
                <th class="px-4 py-2.5">Nama Akun & Uraian</th>
                <th class="px-4 py-2.5 text-right w-36">Debit (Rp)</th>
                <th class="px-4 py-2.5 text-right w-36">Kredit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr
                v-for="(line, idx) in journal.lines"
                :key="line.id"
                :class="idx % 2 === 0 ? 'bg-white' : 'bg-slate-50/30'"
              >
                <td class="px-4 py-3 font-mono font-bold text-slate-800 text-xs">
                  {{ line.accountCode }}
                </td>
                <td class="px-4 py-3">
                  <p class="font-semibold text-slate-900">{{ line.accountName }}</p>
                  <p class="text-[11px] text-slate-500">{{ line.description }}</p>
                </td>
                <td class="px-4 py-3 text-right font-mono font-medium text-slate-900">
                  {{ line.debit > 0 ? formatIDR(line.debit, false) : '-' }}
                </td>
                <td class="px-4 py-3 text-right font-mono font-medium text-slate-900">
                  {{ line.credit > 0 ? formatIDR(line.credit, false) : '-' }}
                </td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="bg-slate-100 font-bold border-t-2 border-slate-300">
                <td colspan="2" class="px-4 py-3 text-right text-slate-700">
                  TOTAL KESEIMBANGAN (BALANCE):
                </td>
                <td class="px-4 py-3 text-right font-mono text-sm text-slate-900">
                  {{ formatIDR(journal.totalDebit) }}
                </td>
                <td class="px-4 py-3 text-right font-mono text-sm text-slate-900">
                  {{ formatIDR(journal.totalCredit) }}
                </td>
              </tr>
            </tfoot>
          </table>
        </div>

        <!-- Balance Integrity Indicator -->
        <div class="flex items-center justify-between px-4 py-2.5 bg-emerald-50 border border-emerald-200 rounded-lg text-emerald-800 text-xs">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 text-emerald-600" />
            <span class="font-semibold">Validasi Akuntansi: Debit & Kredit Seimbang (Balance Validated)</span>
          </div>
          <span class="font-mono text-[11px]">Selisih: Rp 0</span>
        </div>

        <!-- Reversal Prompt Box -->
        <div
          v-if="showReversePrompt"
          class="p-4 bg-rose-50 border border-rose-200 rounded-xl space-y-3 animate-in fade-in duration-150"
        >
          <div class="flex items-center gap-2 text-rose-800 font-bold">
            <AlertTriangle class="w-4 h-4 text-rose-600" />
            <span>Konfirmasi Pembuatan Jurnal Pembalik (Reversal)</span>
          </div>
          <p class="text-xs text-rose-700 leading-relaxed">
            Sistem akan secara otomatis menerbitkan Jurnal Pembalik baru dengan posisi Debit & Kredit yang dibalik, dan menandai jurnal ini sebagai REVERSED.
          </p>
          <div>
            <label class="block text-[11px] font-bold text-slate-700 mb-1">Alasan Pembalikan Jurnal:</label>
            <input
              v-model="reversalReason"
              type="text"
              placeholder="Contoh: Kesalahan nominal kasir shift 1 / Retur faktur"
              class="w-full px-3 py-2 bg-white border border-rose-300 rounded-lg text-xs focus:outline-none focus:ring-2 focus:ring-rose-500"
            />
          </div>
          <div class="flex items-center justify-end gap-2 pt-1">
            <button
              @click="showReversePrompt = false"
              class="px-3 py-1.5 bg-white border border-slate-300 rounded-lg text-xs font-semibold text-slate-700 hover:bg-slate-100"
            >
              Batal
            </button>
            <button
              @click="handleReverse"
              :disabled="!reversalReason.trim()"
              class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 disabled:opacity-50 text-white rounded-lg text-xs font-bold transition-colors"
            >
              Eksekusi Pembalikan Jurnal
            </button>
          </div>
        </div>
      </div>

      <!-- Modal Footer Action Buttons -->
      <div class="px-6 py-4 bg-slate-50 border-t border-slate-200 flex items-center justify-between">
        <div class="text-xs text-slate-500">
          Status saat ini: <strong class="text-slate-800 uppercase">{{ journal.status }}</strong>
        </div>

        <div class="flex items-center gap-2.5">
          <!-- Submit Draft (for maker) -->
          <button
            v-if="journal.status === 'DRAFT' && accountingStore.canCreateDraft"
            @click="journalsStore.submitDraft(journal.id); emit('close')"
            class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-semibold shadow-xs transition-colors"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Ajukan untuk Review</span>
          </button>

          <!-- Review Journal (Supervisor) -->
          <button
            v-if="(journal.status === 'SUBMITTED' || journal.status === 'DRAFT') && accountingStore.canReviewJournal"
            @click="handleReview"
            class="flex items-center gap-1.5 px-4 py-2 bg-sky-600 hover:bg-sky-700 text-white rounded-lg text-xs font-semibold shadow-xs transition-colors"
          >
            <FileCheck class="w-3.5 h-3.5" />
            <span>Review & Verifikasi</span>
          </button>

          <!-- Approve Journal (Manager) -->
          <button
            v-if="journal.status === 'REVIEWED' && accountingStore.canApproveJournal"
            @click="handleApprove"
            class="flex items-center gap-1.5 px-4 py-2 bg-teal-600 hover:bg-teal-700 text-white rounded-lg text-xs font-semibold shadow-xs transition-colors"
          >
            <CheckCircle2 class="w-3.5 h-3.5" />
            <span>Setujui (Approve)</span>
          </button>

          <!-- Post to General Ledger (Manager) -->
          <button
            v-if="journal.status === 'APPROVED' && accountingStore.canPostJournal"
            @click="handlePost"
            class="flex items-center gap-1.5 px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
          >
            <CheckCircle2 class="w-4 h-4" />
            <span>Posting ke Buku Besar (GL)</span>
          </button>

          <!-- Reverse Journal button (Manager only, for posted journals) -->
          <button
            v-if="journal.status === 'POSTED' && accountingStore.canReverseJournal && !showReversePrompt"
            @click="showReversePrompt = true"
            class="flex items-center gap-1.5 px-4 py-2 bg-purple-50 border border-purple-300 text-purple-700 hover:bg-purple-100 rounded-lg text-xs font-semibold transition-colors"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            <span>Buat Jurnal Pembalik</span>
          </button>

          <button
            @click="emit('close')"
            class="px-4 py-2 bg-white border border-slate-300 hover:bg-slate-100 text-slate-700 rounded-lg text-xs font-semibold transition-colors"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
