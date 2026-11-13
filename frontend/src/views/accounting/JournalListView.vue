<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel } from '@/utils/export'
import StatusBadge from '@/components/common/StatusBadge.vue'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import ManualJournalModal from '@/components/journal/ManualJournalModal.vue'
import type { JournalEntry, JournalStatus } from '@/types/accounting'
import {
  BookOpen,
  Search,
  Plus,
  Download,
  Eye,
  CheckCircle2,
  FileCheck2,
  Send,
  RotateCcw
} from 'lucide-vue-next'

const route = useRoute()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

// URL query support (e.g. ?status=DRAFT)
const initialStatus = (route.query.status as string) || 'ALL'
const selectedStatusTab = ref<string>(initialStatus)
const searchQuery = ref('')
const selectedSource = ref<string>('ALL')

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)
const showManualModal = ref(false)

function inspectJournal(j: JournalEntry) {
  selectedJournal.value = j
  showDetailModal.value = true
}

// Status tab options
const tabs = [
  { id: 'ALL', label: 'Semua Jurnal' },
  { id: 'DRAFT', label: 'Draft' },
  { id: 'SUBMITTED', label: 'Submitted' },
  { id: 'REVIEWED', label: 'Reviewed' },
  { id: 'APPROVED', label: 'Approved' },
  { id: 'POSTED', label: 'Posted' },
  { id: 'REVERSED', label: 'Reversed' }
]

const filteredJournals = computed(() => {
  return journalsStore.filteredJournals.filter(j => {
    // Status tab filter
    if (selectedStatusTab.value !== 'ALL' && j.status !== selectedStatusTab.value) {
      return false
    }
    // Source filter
    if (selectedSource.value !== 'ALL' && j.journalSource !== selectedSource.value) {
      return false
    }
    // Search query
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      const matchNo = j.journalNumber.toLowerCase().includes(q)
      const matchDesc = j.description.toLowerCase().includes(q)
      const matchRef = j.sourceReference.toLowerCase().includes(q)
      const matchLines = j.lines.some(l => l.accountName.toLowerCase().includes(q) || l.accountCode.toLowerCase().includes(q))
      if (!matchNo && !matchDesc && !matchRef && !matchLines) return false
    }
    return true
  })
})

function handleExportExcel() {
  const data = filteredJournals.value.map(j => ({
    'Nomor Jurnal': j.journalNumber,
    'Tanggal': j.journalDate,
    'Unit Bisnis': j.businessUnit,
    'Cabang': j.branchName,
    'Sumber': j.journalSource,
    'Referensi': j.sourceReference,
    'Uraian': j.description,
    'Total Nilai (IDR)': j.totalDebit,
    'Status': j.status,
    'Dibuat Oleh': j.createdByName,
    'Diposting Oleh': j.postedByName || '-'
  }))

  exportToExcel(data, `Daftar_Jurnal_${accountingStore.businessUnit}_2026-03`, 'Jurnal Umum')
  accountingStore.showToast('Data jurnal berhasil diekspor ke Excel (.xlsx)', 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Buku Jurnal Umum (General Journal)
          </span>
          <span class="text-xs text-slate-400">&bull; Periode Berjalan</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Daftar Jurnal Akuntansi
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Kelola pencatatan double-entry, periksa status review & approval, serta posting jurnal ke Buku Besar.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="showManualModal = true"
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Plus class="w-4 h-4" />
          <span>Buat Jurnal Manual</span>
        </button>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-semibold border border-slate-300 shadow-2xs transition-colors"
        >
          <Download class="w-4 h-4 text-slate-600" />
          <span>Ekspor Excel</span>
        </button>
      </div>
    </div>

    <!-- Status Sub-Tabs -->
    <div class="flex items-center gap-1 overflow-x-auto pb-1 border-b border-slate-200 text-xs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        @click="selectedStatusTab = tab.id"
        :class="[
          'px-4 py-2 rounded-t-lg font-bold transition-all whitespace-nowrap',
          selectedStatusTab === tab.id
            ? 'bg-white text-blue-600 border-t-2 border-t-blue-600 border-x border-slate-200 shadow-2xs'
            : 'text-slate-500 hover:text-slate-900 hover:bg-slate-100/60'
        ]"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Toolbar Filters & Search -->
    <div class="bg-white p-4 rounded-xl border border-slate-200/80 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2.5 w-full md:w-auto">
        <!-- Search -->
        <div class="relative w-full md:w-80">
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari no jurnal, uraian, akun, referensi..."
            class="w-full pl-9 pr-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <!-- Source Selector -->
        <div class="relative">
          <select
            v-model="selectedSource"
            class="px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium text-slate-800 focus:ring-2 focus:ring-blue-500 focus:outline-none cursor-pointer"
          >
            <option value="ALL">Semua Sumber Jurnal</option>
            <option value="CASHIER_CLOSING">Closing Kasir</option>
            <option value="PHARMACY_PURCHASE">Pembelian Farmasi (PO)</option>
            <option value="DOCTOR_FEE">Jasa Medis Dokter</option>
            <option value="FIXED_ASSET">Penyusutan Aset</option>
            <option value="INVENTORY_CONSUMPTION">Pemakaian BHP</option>
            <option value="MANUAL">Jurnal Penyesuaian Manual</option>
          </select>
        </div>
      </div>

      <div class="text-xs text-slate-500 font-medium">
        Menampilkan <strong>{{ filteredJournals.length }}</strong> jurnal
      </div>
    </div>

    <!-- Journals Data Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">No Jurnal</th>
              <th class="px-4 py-3">Tanggal</th>
              <th class="px-4 py-3">Sumber & Ref</th>
              <th class="px-4 py-3">Unit & Cabang</th>
              <th class="px-4 py-3">Keterangan Memo</th>
              <th class="px-4 py-3 text-right">Total Debit / Kredit</th>
              <th class="px-4 py-3 text-center">Status</th>
              <th class="px-4 py-3 text-center w-32">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="j in filteredJournals"
              :key="j.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5 font-mono font-bold text-blue-700 whitespace-nowrap">
                {{ j.journalNumber }}
              </td>
              <td class="px-4 py-3.5 text-slate-600 whitespace-nowrap">
                {{ formatDate(j.journalDate) }}
              </td>
              <td class="px-4 py-3.5">
                <span class="block font-semibold text-slate-800">{{ j.journalSource }}</span>
                <span class="text-[10px] text-slate-400 font-mono">{{ j.sourceReference }}</span>
              </td>
              <td class="px-4 py-3.5">
                <StatusBadge :status="j.businessUnit" size="sm" />
                <span class="block text-[10px] text-slate-500 mt-0.5">{{ j.branchName }}</span>
              </td>
              <td class="px-4 py-3.5 max-w-sm">
                <p class="font-medium text-slate-800 truncate">{{ j.description }}</p>
                <p class="text-[10px] text-slate-400">
                  {{ j.lines.length }} Baris Akun &bull; Dibuat oleh {{ j.createdByName }}
                </p>
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
                    @click="inspectJournal(j)"
                    class="px-2.5 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-md text-[11px] font-semibold border border-slate-200 transition-colors flex items-center gap-1"
                  >
                    <Eye class="w-3.5 h-3.5 text-slate-500" />
                    <span>Periksa</span>
                  </button>

                  <button
                    v-if="j.status === 'APPROVED' && accountingStore.canPostJournal"
                    @click="journalsStore.postJournal(j.id)"
                    class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md text-[11px] font-bold shadow-2xs transition-colors"
                    title="Posting ke Buku Besar"
                  >
                    Post
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredJournals.length === 0">
              <td colspan="8" class="px-4 py-8 text-center text-slate-400">
                Tidak ada jurnal akuntansi yang ditemukan pada kategori ini.
              </td>
            </tr>
          </tbody>
        </table>
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

