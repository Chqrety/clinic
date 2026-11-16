<script setup lang="ts">
import { ref, computed } from 'vue'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import StatusBadge from '@/components/common/StatusBadge.vue'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  FileCheck2,
  CheckCircle2,
  Check,
  Eye,
  Send,
  AlertTriangle
} from 'lucide-vue-next'

const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)

function inspect(j: JournalEntry) {
  selectedJournal.value = j
  showDetailModal.value = true
}

// Queue list: SUBMITTED, REVIEWED, APPROVED
const queueList = computed(() => {
  return journalsStore.filteredJournals.filter(j =>
    ['SUBMITTED', 'REVIEWED', 'APPROVED'].includes(j.status)
  )
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-sky-600 uppercase tracking-wider bg-sky-50 px-2 py-0.5 rounded border border-sky-200">
            Authorization & Segregation of Duties
          </span>
          <span class="text-xs text-slate-400">&bull; Maker-Checker-Approver</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Antrian Review & Approval Jurnal
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Memastikan prinsip dual control: jurnal yang dibuat oleh staf akuntansi harus direview oleh supervisor dan disetujui oleh Finance Manager sebelum diposting.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="px-4 py-2 bg-sky-50 border border-sky-200 rounded-xl text-sky-900 text-xs text-right">
          <span class="text-[10px] uppercase font-bold text-sky-600 block">Total Antrian:</span>
          <span class="font-bold text-lg font-mono">{{ queueList.length }} Jurnal</span>
        </div>
      </div>
    </div>

    <!-- Active User Role Indicator -->
    <div class="p-4 bg-blue-50/70 border border-blue-200 rounded-xl flex items-center justify-between text-xs text-blue-900">
      <div class="flex items-center gap-2.5">
        <FileCheck2 class="w-5 h-5 text-blue-600" />
        <div>
          <span class="font-bold">Otorisasi Login: {{ accountingStore.currentUserName }} ({{ accountingStore.currentUserRole.replace('_', ' ').toUpperCase() }})</span>
          <p class="text-[11px] text-blue-700">
            {{ accountingStore.canApproveJournal ? 'Anda berwenang melakukan Approval dan Posting ke Buku Besar.' : 'Anda memiliki akses Review / Maker.' }}
          </p>
        </div>
      </div>

      <span class="px-2.5 py-1 bg-white rounded-lg border border-blue-300 font-bold text-[11px] text-blue-800 shadow-2xs">
        Level Otoritas Aktif
      </span>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">No Jurnal</th>
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">Sumber & Ref</th>
              <th class="px-4 py-3">Unit Bisnis</th>
              <th class="px-4 py-3">Uraian / Memo Jurnal</th>
              <th class="px-4 py-3 text-right">Total Nilai (IDR)</th>
              <th class="px-4 py-3 text-center">Status Saat Ini</th>
              <th class="px-4 py-3 text-center w-48">Tindakan Persetujuan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="j in queueList"
              :key="j.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 font-mono font-bold text-blue-700">
                {{ j.journalNumber }}
              </td>
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ formatDate(j.journalDate) }}
              </td>
              <td class="px-4 py-3.5">
                <span class="font-semibold text-slate-800">{{ j.journalSource }}</span>
                <span class="block text-[10px] text-slate-400 font-mono">{{ j.sourceReference }}</span>
              </td>
              <td class="px-4 py-3.5">
                <StatusBadge :status="j.businessUnit" size="sm" />
              </td>
              <td class="px-4 py-3.5 max-w-sm">
                <p class="font-medium text-slate-800 truncate">{{ j.description }}</p>
                <p class="text-[10px] text-slate-400">Maker: {{ j.createdByName }}</p>
              </td>
              <td class="px-4 py-3.5 text-right font-mono font-bold text-slate-900 whitespace-nowrap">
                {{ formatIDR(j.totalDebit) }}
              </td>
              <td class="px-4 py-3.5 text-center">
                <StatusBadge :status="j.status" size="sm" />
              </td>
              <td class="px-4 py-3.5 text-center">
                <div class="flex items-center justify-center gap-1.5">
                  <button
                    @click="inspect(j)"
                    class="p-1.5 rounded-lg text-slate-500 hover:text-blue-600 hover:bg-blue-50 transition-colors"
                    title="Periksa Pos Jurnal"
                  >
                    <Eye class="w-4 h-4" />
                  </button>

                  <!-- Step 1: Review -->
                  <button
                    v-if="j.status === 'SUBMITTED' && accountingStore.canReviewJournal"
                    @click="journalsStore.reviewJournal(j.id)"
                    class="px-2.5 py-1 bg-sky-600 hover:bg-sky-700 text-white rounded-md text-[11px] font-bold shadow-2xs transition-colors flex items-center gap-1"
                  >
                    <Check class="w-3 h-3" />
                    <span>Review</span>
                  </button>

                  <!-- Step 2: Approve -->
                  <button
                    v-if="j.status === 'REVIEWED' && accountingStore.canApproveJournal"
                    @click="journalsStore.approveJournal(j.id)"
                    class="px-2.5 py-1 bg-teal-600 hover:bg-teal-700 text-white rounded-md text-[11px] font-bold shadow-2xs transition-colors flex items-center gap-1"
                  >
                    <CheckCircle2 class="w-3 h-3" />
                    <span>Approve</span>
                  </button>

                  <!-- Step 3: Post -->
                  <button
                    v-if="j.status === 'APPROVED' && accountingStore.canPostJournal"
                    @click="journalsStore.postJournal(j.id)"
                    class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md text-[11px] font-bold shadow-2xs transition-colors flex items-center gap-1"
                  >
                    <CheckCircle2 class="w-3 h-3" />
                    <span>Post ke GL</span>
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="queueList.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-slate-400">
                Tidak ada jurnal yang sedang menunggu review atau approval saat ini.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />
  </div>
</template>

