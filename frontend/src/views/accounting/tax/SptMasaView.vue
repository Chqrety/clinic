<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, type SptMasaType } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportSptMasaToExcel, getEntityProfile } from '@/utils/export'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import {
  Calculator,
  Download,
  FileSpreadsheet,
  CheckCircle2,
  FileText,
  Building,
  UserCheck,
  Receipt,
  Layers,
  ArrowRight,
  Send,
  Calendar,
  AlertCircle,
  ExternalLink,
  ShieldCheck,
  BookOpen,
  Plus,
  X,
  Printer,
  Sparkles,
  Link2,
  Landmark,
  Database,
  ChevronDown,
  ChevronRight
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

onMounted(async () => {
  await Promise.all([
    taxStore.fetchTaxRecords(),
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
})

// Read query params if any
const route = router.currentRoute.value
if (route.query.taxCategory) {
  const cat = String(route.query.taxCategory)
  const matched = taxStore.taxAccountOptions.find(opt => opt.taxCategory === cat)
  if (matched) taxStore.selectedTaxAccountCode = matched.code
}

const doctorTaxRecords = computed(() => taxStore.filteredDoctorTax)
const ppn = computed(() => taxStore.ppnSummary)
const txList = computed(() => taxStore.filteredTaxTransactions)
const summary = computed(() => taxStore.monthlyTaxSummary)
const currentEntity = computed(() => getEntityProfile(accountingStore.businessUnit))

const taxFilterOptions = [
  { value: 'ALL', label: 'Semua Jenis SPT Masa' },
  { value: 'PPH_21', label: 'PPh Pasal 21', sublabel: 'Dokter & Staf' },
  { value: 'PPH_23', label: 'PPh Pasal 23', sublabel: 'Jasa Medis & Limbah B3' },
  { value: 'PPH_4_2', label: 'PPh Final Pasal 4(2)', sublabel: 'Sewa & Bunga' },
  { value: 'PPH_25', label: 'PPh Pasal 25', sublabel: 'Angsuran Badan' },
  { value: 'PPN', label: 'PPN', sublabel: 'Keluaran & Masukan' }
]

const taxAccountOptions = computed(() => {
  return [
    { value: 'ALL', label: 'Semua Akun Pajak (2500 - 2504)' },
    ...taxStore.taxAccountOptions.map(opt => ({
      value: opt.code,
      label: `[${opt.code}] ${opt.name}`,
      sublabel: opt.description
    }))
  ]
})
const sptMasaList = computed(() => taxStore.filteredSptMasa)

// Subview under SPT Masa: 'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU' | 'JURNAL_TERKAIT'
const sptMasaSubTab = ref<'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU' | 'JURNAL_TERKAIT'>('DOKUMEN_SPT')

const showDetailModal = ref(false)
const showPaymentBankModal = ref(false)
const showJournalPreviewModal = ref(false)
const selectedSptDetail = ref<any>(null)
const activeJournal = ref<any>(null)
const selectedBankForPayment = ref<string>('1102.001')

// Expandable SPT Row state for viewing individual journal entries
const expandedSptIds = ref<Set<string>>(new Set())

const isAllExpanded = computed(() => {
  return sptMasaList.value.length > 0 && expandedSptIds.value.size >= sptMasaList.value.length
})

function toggleExpandAll() {
  if (isAllExpanded.value) {
    expandedSptIds.value.clear()
  } else {
    expandedSptIds.value = new Set(sptMasaList.value.map(s => s.id))
  }
}

function toggleExpandSpt(id: string) {
  if (expandedSptIds.value.has(id)) {
    expandedSptIds.value.delete(id)
  } else {
    expandedSptIds.value.add(id)
  }
}

// Selected Account Metadata
const selectedAccountMeta = computed(() => {
  return taxStore.taxAccountOptions.find(o => o.code === taxStore.selectedTaxAccountCode) || taxStore.taxAccountOptions[0]
})

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/accounting/ledger', query: { account: accountCode } })
  } else {
    router.push('/accounting/ledger')
  }
}

function handleExecutePayment() {
  const code = taxStore.selectedTaxAccountCode === 'ALL' ? '2501' : taxStore.selectedTaxAccountCode
  const name = selectedAccountMeta.value.name
  const amt = summary.value.totalTax || 39249649
  taxStore.generateTaxPaymentJournal(code, name, amt, selectedBankForPayment.value)
  showPaymentBankModal.value = false
}

function viewSptDetail(item: any) {
  selectedSptDetail.value = item
  showDetailModal.value = true
}

function openJournalDetail(entry: any) {
  if (entry.rawJournal) {
    activeJournal.value = entry.rawJournal
  } else {
    // Look up in journalsStore
    const found = journalsStore.journals.find(j => j.journalNumber === entry.journalNumber || j.id === entry.journalId)
    if (found) {
      activeJournal.value = found
    } else {
      activeJournal.value = {
        journalNumber: entry.journalNumber,
        journalDate: entry.journalDate,
        description: entry.description,
        businessUnit: entry.businessUnit,
        status: entry.status,
        lines: [
          {
            accountCode: entry.taxAccountCode,
            accountName: entry.taxAccountName,
            description: entry.description,
            debit: entry.direction === 'PAID' ? entry.taxAmount : 0,
            credit: entry.direction === 'WITHHELD' ? entry.taxAmount : 0
          },
          {
            accountCode: entry.counterpartAccountCode,
            accountName: entry.counterpartAccountName,
            description: entry.description,
            debit: entry.direction === 'WITHHELD' ? entry.taxAmount : 0,
            credit: entry.direction === 'PAID' ? entry.taxAmount : 0
          }
        ]
      }
    }
  }
  showJournalPreviewModal.value = true
}

function drillDownToJournals(month: number, taxType: string) {
  taxStore.selectedMonth = month
  taxStore.selectedSptMasaTypeFilter = taxType
  sptMasaSubTab.value = 'JURNAL_TERKAIT'
}

function handleExportExcel() {
  exportSptMasaToExcel({
    sptList: sptMasaList.value,
    doctorList: doctorTaxRecords.value,
    txList: txList.value,
    summary: summary.value,
    linkedJournals: taxStore.filteredLinkedTaxJournalEntries
  }, {
    businessUnit: accountingStore.businessUnit,
    masaMonth: taxStore.selectedMonth,
    fiscalYear: taxStore.selectedYear
  })
  accountingStore.showToast('Dokumen & Kertas Kerja SPT Masa (Beserta Rincian Per Jurnal) berhasil diekspor ke Excel', 'success')
}

function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Screen View (Hidden during print) -->
    <div class="print:hidden space-y-5">
      <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Kepatuhan Fiskal DJP
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">BCA & Mandiri Verified</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          SPT Masa (Pajak Bulanan Klinik)
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <!-- Ke Halaman Baru Buat SPT Masa -->
        <router-link
          to="/accounting/tax/spt-masa/create"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Buat SPT Masa Baru</span>
        </router-link>

        <button
          @click="showPaymentBankModal = true"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Landmark class="w-3.5 h-3.5 text-zinc-600" />
          <span>Penyetoran via Bank (NTPN)</span>
        </button>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Ekspor Excel SPT Masa PPh & PPN Lengkap beserta Detail Per Jurnal"
        >
          <FileSpreadsheet class="w-3.5 h-3.5 text-emerald-600" />
          <span>Ekspor SPT (.xlsx)</span>
        </button>

        <button
          @click="handlePrint"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Cetak Formulir Resmi Bukti Pelaporan SPT Masa DJP"
        >
          <Printer class="w-3.5 h-3.5 text-blue-600" />
          <span>Cetak Formulir</span>
        </button>

        <button
          @click="taxStore.exportEbupotCSV"
          class="inline-flex items-center gap-1.5 px-3 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Ekspor CSV Bukti Potong DJP e-Bupot PPh 21 Dokter"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>CSV e-Bupot</span>
        </button>
      </div>
    </header>

    <!-- Sub-tab Navigation: Dokumen SPT Masa vs e-Bupot Dokter vs Buku Pembantu Pajak -->
    <div class="flex flex-wrap items-center justify-between gap-3 border-b border-zinc-200 pb-2">
      <div class="flex items-center gap-1.5 bg-zinc-100 p-1 rounded-xl">
        <button
          @click="sptMasaSubTab = 'DOKUMEN_SPT'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all',
            sptMasaSubTab === 'DOKUMEN_SPT'
              ? 'bg-white text-zinc-950 shadow-2xs'
              : 'text-zinc-600 hover:text-zinc-950'
          ]"
        >
          Daftar Dokumen SPT Masa
        </button>

        <button
          @click="sptMasaSubTab = 'EBUPOT_DOKTER'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all',
            sptMasaSubTab === 'EBUPOT_DOKTER'
              ? 'bg-white text-zinc-950 shadow-2xs'
              : 'text-zinc-600 hover:text-zinc-950'
          ]"
        >
          e-Bupot PPh 21 Dokter Spesialis
        </button>

        <button
          @click="sptMasaSubTab = 'BUKU_PEMBANTU'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all',
            sptMasaSubTab === 'BUKU_PEMBANTU'
              ? 'bg-white text-zinc-950 shadow-2xs'
              : 'text-zinc-600 hover:text-zinc-950'
          ]"
        >
          Buku Pembantu Hutang Pajak
        </button>

        <button
          @click="sptMasaSubTab = 'JURNAL_TERKAIT'"
          :class="[
            'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all flex items-center gap-1.5',
            sptMasaSubTab === 'JURNAL_TERKAIT'
              ? 'bg-white text-zinc-950 shadow-2xs'
              : 'text-zinc-600 hover:text-zinc-950'
          ]"
        >
          <FileText class="w-3.5 h-3.5 text-zinc-500" />
          <span>Transaksi & Jurnal Terkait</span>
          <span class="px-1.5 py-0.2 bg-zinc-200 text-zinc-800 rounded-full text-[10px] font-bold">
            {{ taxStore.filteredLinkedTaxJournalEntries.length }}
          </span>
        </button>
      </div>

      <!-- Quick Link to SPT Tahunan 1771 -->
      <router-link
        to="/accounting/tax/spt-tahunan"
        class="inline-flex items-center gap-1.5 text-xs font-bold text-zinc-700 hover:text-zinc-950 transition bg-white px-3 py-1.5 rounded-xl border border-zinc-200 shadow-2xs"
      >
        <span>Lihat Akumulasi ke SPT Tahunan (1771)</span>
        <ArrowRight class="w-3.5 h-3.5" />
      </router-link>
    </div>

    <!-- VIEW 1: DAFTAR DOKUMEN SPT MASA -->
    <div v-if="sptMasaSubTab === 'DOKUMEN_SPT'" class="space-y-4">
      <!-- Filter Bar -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs">
        <div class="flex flex-wrap items-center gap-3 text-xs">
          <div>
            <label class="block text-[10px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Masa Pajak</label>
            <select
              v-model.number="taxStore.selectedMonth"
              class="px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-800 focus:outline-none cursor-pointer"
            >
              <option :value="1">01 - Januari</option>
              <option :value="2">02 - Februari</option>
              <option :value="3">03 - Maret</option>
              <option :value="4">04 - April</option>
              <option :value="5">05 - Mei</option>
              <option :value="6">06 - Juni</option>
              <option :value="7">07 - Juli</option>
              <option :value="8">08 - Agustus</option>
              <option :value="9">09 - September</option>
              <option :value="10">10 - Oktober</option>
              <option :value="11">11 - November</option>
              <option :value="12">12 - Desember</option>
            </select>
          </div>

          <div class="w-56">
            <label class="block text-[10px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Jenis SPT</label>
            <SearchableSelect
              :model-value="taxStore.selectedSptMasaTypeFilter"
              :options="taxFilterOptions"
              placeholder="Cari jenis SPT..."
              input-class="!py-1.5"
              @update:model-value="(val: any) => taxStore.selectedSptMasaTypeFilter = val"
            />
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="toggleExpandAll"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-800 rounded-xl text-xs font-bold transition cursor-pointer"
            title="Buka atau tutup seluruh rincian transaksi jurnal umum di bawah"
          >
            <Layers class="w-3.5 h-3.5 text-zinc-600" />
            <span>{{ isAllExpanded ? 'Tutup Rincian' : 'Buka Semua Rincian Jurnal' }}</span>
          </button>
          <span class="text-xs text-zinc-500 font-medium">
            <strong>{{ sptMasaList.length }}</strong> dokumen SPT Masa
          </span>
        </div>
      </div>

      <!-- SPT Masa Table -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-4 py-3">No. Dokumen SPT</th>
                <th class="px-4 py-3">Jenis Pajak</th>
                <th class="px-4 py-3 text-center">Masa</th>
                <th class="px-4 py-3">Rekening Bank Penyetor</th>
                <th class="px-4 py-3 text-right">DPP (Rp)</th>
                <th class="px-4 py-3 text-right">Pajak Terutang (Rp)</th>
                <th class="px-4 py-3 text-center">Status / NTPN</th>
                <th class="px-4 py-3 text-center w-24">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <template v-for="spt in sptMasaList" :key="spt.id">
                <tr class="hover:bg-zinc-50/70 transition">
                  <td class="px-4 py-2.5 font-mono font-bold text-zinc-900">
                    {{ spt.sptNumber }}
                  </td>
                  <td class="px-4 py-2.5">
                    <div class="flex items-center gap-1.5">
                      <span class="font-bold text-zinc-950">{{ spt.taxType }}</span>
                      <span class="px-1.5 py-0.2 bg-zinc-100 text-zinc-700 rounded text-[10px] font-mono font-bold">
                        {{ (spt.linkedJournals || []).length }} Jurnal
                      </span>
                    </div>
                    <span class="block text-[10px] text-zinc-500 truncate max-w-xs">{{ spt.accountName }}</span>
                    <div class="mt-1 flex items-center gap-2">
                      <button
                        @click.stop="toggleExpandSpt(spt.id)"
                        class="inline-flex items-center gap-1 text-[10px] font-bold text-emerald-800 bg-emerald-50 hover:bg-emerald-100 px-2 py-0.5 rounded-md border border-emerald-200 transition cursor-pointer"
                        :title="expandedSptIds.has(spt.id) ? 'Tutup rincian transaksi' : 'Buka rincian voucher Jurnal Umum pajak masa ini'"
                      >
                        <ChevronDown v-if="expandedSptIds.has(spt.id)" class="w-3 h-3 text-emerald-700" />
                        <ChevronRight v-else class="w-3 h-3 text-emerald-700" />
                        <span>{{ expandedSptIds.has(spt.id) ? 'Tutup Rincian' : 'Rincian Jurnal' }} ({{ (spt.linkedJournals || []).length }})</span>
                      </button>
                      <button
                        @click.stop="drillDownToJournals(spt.masaMonth, spt.taxType)"
                        class="text-[10px] text-zinc-500 hover:text-zinc-900 underline"
                        title="Buka tab Jurnal Terkait lengkap"
                      >
                        Tab Jurnal &rarr;
                      </button>
                    </div>
                  </td>
                  <td class="px-4 py-2.5 text-center font-bold font-mono">
                    Masa {{ spt.masaMonth }}/{{ spt.masaYear }}
                  </td>
                  <td class="px-4 py-2.5">
                    <span class="font-bold text-zinc-900 flex items-center gap-1">
                      <Landmark class="w-3 h-3 text-zinc-500" />
                      <span>{{ spt.bankAccountName }}</span>
                    </span>
                    <span class="text-[10px] text-zinc-400 font-mono">{{ spt.bankAccountCode }}</span>
                  </td>
                  <td class="px-4 py-2.5 text-right tabular-nums text-zinc-600">
                    {{ formatIDR(spt.dppAmount) }}
                  </td>
                  <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">
                    {{ formatIDR(spt.taxAmount) }}
                  </td>
                  <td class="px-4 py-2.5 text-center">
                    <span
                      :class="[
                        'px-2 py-0.5 rounded-full text-[10px] font-bold',
                        spt.status === 'REPORTED'
                          ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                          : spt.status === 'PAID'
                          ? 'bg-blue-50 text-blue-800 border border-blue-200'
                          : 'bg-zinc-100 text-zinc-600 border border-zinc-200'
                      ]"
                    >
                      {{ spt.status }}
                    </span>
                    <span v-if="spt.ntpnNumber" class="block font-mono text-[9px] text-zinc-400 mt-0.5">
                      {{ spt.ntpnNumber }}
                    </span>
                  </td>
                  <td class="px-4 py-2.5 text-center">
                    <button
                      @click="viewSptDetail(spt)"
                      class="px-2.5 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-950 rounded-lg text-xs font-bold border border-zinc-200 transition cursor-pointer"
                    >
                      Detail
                    </button>
                  </td>
                </tr>

                <!-- Sub-table of individual tax general journals for this SPT Masa -->
                <tr v-if="expandedSptIds.has(spt.id)" class="bg-zinc-50/90 border-b border-zinc-200/90">
                  <td colspan="8" class="p-3.5">
                    <div class="bg-white rounded-xl border border-zinc-200 shadow-2xs overflow-hidden">
                      <div class="px-4 py-2.5 bg-zinc-100/70 border-b border-zinc-200 flex items-center justify-between">
                        <div class="flex items-center gap-2">
                          <FileText class="w-4 h-4 text-emerald-700" />
                          <h4 class="font-bold text-zinc-950 text-xs">
                            Rincian Transaksi Jurnal Umum — {{ spt.taxType }} Masa {{ spt.masaMonth }}/{{ spt.masaYear }}
                          </h4>
                          <span class="px-2 py-0.5 bg-emerald-50 text-emerald-800 text-[10px] font-bold rounded-full border border-emerald-200">
                            {{ (spt.linkedJournals || []).length }} Transaksi
                          </span>
                        </div>
                        <router-link
                          to="/accounting/journals/create"
                          class="text-[11px] font-bold text-blue-600 hover:text-blue-800 flex items-center gap-1 transition"
                        >
                          <Plus class="w-3 h-3" />
                          <span>Buat Jurnal Pajak Baru</span>
                        </router-link>
                      </div>

                      <div v-if="(spt.linkedJournals || []).length > 0" class="overflow-x-auto">
                        <table class="w-full text-left border-collapse text-xs">
                          <thead>
                            <tr class="bg-zinc-50/80 border-b border-zinc-200 text-[10px] font-bold text-zinc-500 uppercase tracking-wider">
                              <th class="px-3 py-2 w-8 text-center">#</th>
                              <th class="px-3 py-2 w-28">Tanggal</th>
                              <th class="px-3 py-2 w-36">No. Jurnal</th>
                              <th class="px-3 py-2">Uraian / Memo Jurnal</th>
                              <th class="px-3 py-2">Akun Lawan</th>
                              <th class="px-3 py-2 text-right w-32">DPP (Rp)</th>
                              <th class="px-3 py-2 text-right w-32">Pajak (Rp)</th>
                              <th class="px-3 py-2 text-center w-24">Voucher</th>
                            </tr>
                          </thead>
                          <tbody class="divide-y divide-zinc-100">
                            <tr v-for="(j, idx) in spt.linkedJournals" :key="j.id" class="hover:bg-zinc-50/60 transition">
                              <td class="px-3 py-2 text-center text-zinc-400 font-mono text-[11px]">{{ idx + 1 }}</td>
                              <td class="px-3 py-2 font-mono text-zinc-600">{{ formatDate(j.journalDate) }}</td>
                              <td class="px-3 py-2 font-mono font-bold text-blue-600 cursor-pointer hover:underline" @click="openJournalDetail(j)">
                                {{ j.journalNumber }}
                              </td>
                              <td class="px-3 py-2 font-medium text-zinc-900">{{ j.description }}</td>
                              <td class="px-3 py-2 text-zinc-600">
                                <span class="block font-medium truncate max-w-xs">{{ j.counterpartAccountName }}</span>
                                <span class="text-[10px] font-mono text-zinc-400">{{ j.counterpartAccountCode }}</span>
                              </td>
                              <td class="px-3 py-2 text-right font-mono tabular-nums text-zinc-700">{{ formatIDR(j.dppAmount) }}</td>
                              <td class="px-3 py-2 text-right font-mono tabular-nums font-bold text-emerald-700">{{ formatIDR(j.taxAmount) }}</td>
                              <td class="px-3 py-2 text-center">
                                <button
                                  @click="openJournalDetail(j)"
                                  class="px-2 py-0.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-800 rounded font-bold text-[10px] transition cursor-pointer"
                                >
                                  Buka
                                </button>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>

                      <div v-else class="p-6 text-center text-xs text-zinc-500">
                        <p>Belum ada transaksi Jurnal Umum pajak untuk jenis {{ spt.taxType }} di Masa {{ spt.masaMonth }}/{{ spt.masaYear }}.</p>
                        <router-link
                          to="/accounting/journals/create"
                          class="mt-2 inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition shadow-2xs"
                        >
                          <Plus class="w-3.5 h-3.5" />
                          <span>Input Jurnal Pajak di Jurnal Umum &rarr;</span>
                        </router-link>
                      </div>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- VIEW 2: EBUPOT DOKTER -->
    <div v-else-if="sptMasaSubTab === 'EBUPOT_DOKTER'" class="space-y-4">
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-900">Withholding PPh 21 Dokter Spesialis & Umum (PMK 168/2023)</h3>
            <p class="text-xs text-zinc-500 mt-0.5">Penghitungan DPP 50% non-pegawai dokter spesialis anak (Sp.A) & spesialis kulit (Sp.DVE).</p>
          </div>
          <button
            @click="taxStore.exportEbupotCSV"
            class="px-3 py-1.5 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition"
          >
            Ekspor Format DJP e-Bupot 21 (.csv)
          </button>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Nama Dokter</th>
                <th class="px-3 py-2.5">Spesialisasi / Unit</th>
                <th class="px-3 py-2.5">NPWP / NIK</th>
                <th class="px-3 py-2.5 text-right">Penghasilan Bruto (Rp)</th>
                <th class="px-3 py-2.5 text-right">DPP 50% (Rp)</th>
                <th class="px-3 py-2.5 text-center">Tarif TER</th>
                <th class="px-3 py-2.5 text-right">PPh 21 Terutang (Rp)</th>
                <th class="px-3 py-2.5 text-right">Net Honor Dokter (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="d in doctorTaxRecords" :key="d.id" class="hover:bg-zinc-50/50">
                <td class="px-3 py-2.5 font-bold text-zinc-950">{{ d.doctorName }}</td>
                <td class="px-3 py-2.5">
                  <span class="font-medium text-zinc-800">{{ d.specialty }}</span>
                  <span class="block text-[10px] text-zinc-400 font-bold uppercase">{{ d.businessUnit }}</span>
                </td>
                <td class="px-3 py-2.5 font-mono text-zinc-600">{{ d.npwp || d.nik }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums text-zinc-800">{{ formatIDR(d.grossIncomeCurrentMonth) }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums text-zinc-600">{{ formatIDR(d.dppAmount) }}</td>
                <td class="px-3 py-2.5 text-center font-bold text-zinc-900">{{ d.taxBracketApplied }} ({{ d.taxRate }}%)</td>
                <td class="px-3 py-2.5 text-right tabular-nums font-bold text-rose-600">{{ formatIDR(d.pph21Amount) }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums font-bold text-emerald-700">{{ formatIDR(d.netPayable) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- VIEW 3: BUKU PEMBANTU PAJAK -->
    <div v-else-if="sptMasaSubTab === 'BUKU_PEMBANTU'" class="space-y-4">
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-900">Buku Pembantu Mutasi Transaksi Pajak</h3>
            <p class="text-xs text-zinc-500 mt-0.5">Rincian invoice operasional yang memicu timbulnya hutang pajak klinik.</p>
          </div>
          <div class="flex items-center gap-2">
            <select
              v-model.number="taxStore.selectedMonth"
              class="px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-xs font-bold text-zinc-800 focus:outline-none cursor-pointer"
            >
              <option :value="1">01 - Jan</option>
              <option :value="2">02 - Feb</option>
              <option :value="3">03 - Mar</option>
              <option :value="4">04 - Apr</option>
              <option :value="5">05 - Mei</option>
              <option :value="6">06 - Jun</option>
              <option :value="7">07 - Jul</option>
              <option :value="8">08 - Agu</option>
              <option :value="9">09 - Sep</option>
              <option :value="10">10 - Okt</option>
              <option :value="11">11 - Nov</option>
              <option :value="12">12 - Des</option>
            </select>
            <div class="w-72">
              <SearchableSelect
                v-model="taxStore.selectedTaxAccountCode"
                :options="taxAccountOptions"
                placeholder="Cari akun pajak..."
                input-class="!py-1.5"
              />
            </div>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-3 py-2.5">Tanggal</th>
                <th class="px-3 py-2.5">No. Bukti / Invoice</th>
                <th class="px-3 py-2.5">Pihak Lawan / Vendor</th>
                <th class="px-3 py-2.5">Uraian Transaksi</th>
                <th class="px-3 py-2.5 text-right">Nilai Bruto (Rp)</th>
                <th class="px-3 py-2.5 text-right">DPP (Rp)</th>
                <th class="px-3 py-2.5 text-right">Hutang Pajak (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="t in txList" :key="t.id" class="hover:bg-zinc-50/50">
                <td class="px-3 py-2.5 font-mono text-zinc-600">{{ t.taxDate }}</td>
                <td class="px-3 py-2.5 font-mono font-bold text-zinc-900">{{ t.documentReference || t.taxNumber }}</td>
                <td class="px-3 py-2.5 font-bold text-zinc-800">{{ t.counterpartyName }}</td>
                <td class="px-3 py-2.5 text-zinc-600 truncate max-w-xs">{{ t.taxCategory }} - {{ t.objectCode }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums text-zinc-600">{{ formatIDR(t.grossAmount) }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums text-zinc-600">{{ formatIDR(t.dppAmount) }}</td>
                <td class="px-3 py-2.5 text-right tabular-nums font-bold text-zinc-950">{{ formatIDR(t.taxAmount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- VIEW 4: TRANSAKSI & JURNAL UMUM TERKAIT (TERIKAT KE JURNAL ENTRi) -->
    <div v-else-if="sptMasaSubTab === 'JURNAL_TERKAIT'" class="space-y-4">
      <!-- Filter Bar & Quick Create -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs">
        <div class="flex flex-wrap items-center gap-3 text-xs">
          <div>
            <label class="block text-[10px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Masa Pajak</label>
            <select
              v-model.number="taxStore.selectedMonth"
              class="px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-800 focus:outline-none cursor-pointer"
            >
              <option :value="1">01 - Januari</option>
              <option :value="2">02 - Februari</option>
              <option :value="3">03 - Maret</option>
              <option :value="4">04 - April</option>
              <option :value="5">05 - Mei</option>
              <option :value="6">06 - Juni</option>
              <option :value="7">07 - Juli</option>
              <option :value="8">08 - Agustus</option>
              <option :value="9">09 - September</option>
              <option :value="10">10 - Oktober</option>
              <option :value="11">11 - November</option>
              <option :value="12">12 - Desember</option>
            </select>
          </div>

          <div class="w-56">
            <label class="block text-[10px] font-bold text-zinc-500 uppercase tracking-wider mb-1">Jenis Pajak</label>
            <SearchableSelect
              :model-value="taxStore.selectedSptMasaTypeFilter"
              :options="taxFilterOptions"
              placeholder="Cari jenis pajak..."
              input-class="!py-1.5"
              @update:model-value="(val: any) => taxStore.selectedSptMasaTypeFilter = val"
            />
          </div>
        </div>

        <div class="flex items-center gap-2">
          <router-link
            to="/accounting/journals/create"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Buat Jurnal Pajak Baru</span>
          </router-link>
          <span class="text-xs text-zinc-500 font-medium ml-2">
            <strong>{{ taxStore.filteredLinkedTaxJournalEntries.length }}</strong> Jurnal Terkait
          </span>
        </div>
      </div>

      <!-- Linked Journal Entries Table -->
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200/90 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
                <th class="px-4 py-3">Tanggal</th>
                <th class="px-4 py-3">No. Jurnal / Ref</th>
                <th class="px-4 py-3">Unit</th>
                <th class="px-4 py-3">Uraian Transaksi</th>
                <th class="px-4 py-3">Akun Pajak Terkait</th>
                <th class="px-4 py-3">Akun Lawan (Kontra)</th>
                <th class="px-4 py-3 text-right">DPP (Rp)</th>
                <th class="px-4 py-3 text-right">Nilai Pajak (Rp)</th>
                <th class="px-4 py-3">Kanal Bank Penyetor</th>
                <th class="px-4 py-3 text-center">Status</th>
                <th class="px-4 py-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr
                v-for="entry in taxStore.filteredLinkedTaxJournalEntries"
                :key="entry.id"
                class="hover:bg-zinc-50/70 transition"
              >
                <td class="px-4 py-2.5 font-mono text-zinc-600 whitespace-nowrap">
                  {{ formatDate(entry.journalDate) }}
                </td>
                <td class="px-4 py-2.5 whitespace-nowrap">
                  <span class="font-mono font-bold text-zinc-950 block">{{ entry.journalNumber }}</span>
                  <span class="text-[10px] text-zinc-400 font-mono">Masa {{ entry.month }}/{{ entry.year }}</span>
                </td>
                <td class="px-4 py-2.5 whitespace-nowrap font-bold text-[10px] uppercase text-zinc-700">
                  {{ entry.businessUnit }}
                </td>
                <td class="px-4 py-2.5 max-w-xs text-zinc-900 font-medium truncate">
                  {{ entry.description }}
                </td>
                <td class="px-4 py-2.5 whitespace-nowrap">
                  <span class="font-bold text-zinc-950 block font-mono text-[11px]">{{ entry.taxAccountCode }}</span>
                  <span class="text-[10px] text-zinc-500 block truncate max-w-[180px]">{{ entry.taxAccountName }}</span>
                </td>
                <td class="px-4 py-2.5 whitespace-nowrap">
                  <span class="font-bold text-zinc-700 block font-mono text-[11px]">{{ entry.counterpartAccountCode }}</span>
                  <span class="text-[10px] text-zinc-400 block truncate max-w-[180px]">{{ entry.counterpartAccountName }}</span>
                </td>
                <td class="px-4 py-2.5 text-right tabular-nums text-zinc-600 whitespace-nowrap">
                  {{ formatIDR(entry.dppAmount) }}
                </td>
                <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950 whitespace-nowrap">
                  {{ formatIDR(entry.taxAmount) }}
                </td>
                <td class="px-4 py-2.5 whitespace-nowrap text-zinc-700">
                  <span class="flex items-center gap-1 font-medium">
                    <Landmark class="w-3 h-3 text-zinc-400" />
                    <span>{{ entry.bankAccount }}</span>
                  </span>
                </td>
                <td class="px-4 py-2.5 text-center whitespace-nowrap">
                  <span
                    :class="[
                      'px-2 py-0.5 rounded-full text-[10px] font-bold',
                      entry.status === 'POSTED'
                        ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                        : 'bg-zinc-100 text-zinc-700 border border-zinc-200'
                    ]"
                  >
                    {{ entry.status }}
                  </span>
                </td>
                <td class="px-4 py-2.5 text-center whitespace-nowrap">
                  <button
                    @click="openJournalDetail(entry)"
                    class="px-2.5 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-950 rounded-lg text-xs font-bold border border-zinc-200 transition cursor-pointer"
                  >
                    Lihat Jurnal
                  </button>
                </td>
              </tr>
              <tr v-if="taxStore.filteredLinkedTaxJournalEntries.length === 0">
                <td colspan="11" class="px-4 py-8 text-center text-zinc-400 italic">
                  Belum ada transaksi Jurnal Umum terkait pajak untuk masa dan filter yang dipilih.
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- DIALOG MODAL 1: Detail SPT Masa -->
    <div
      v-if="showDetailModal && selectedSptDetail"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-2xl w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150 p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">{{ selectedSptDetail.sptNumber }}</h3>
            <span class="text-xs text-zinc-500">Masa Pajak {{ selectedSptDetail.masaMonth }}/{{ selectedSptDetail.masaYear }}</span>
          </div>
          <button @click="showDetailModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="grid grid-cols-2 gap-x-6 gap-y-2 text-xs">
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Jenis SPT:</span>
            <span class="font-bold text-zinc-950">{{ selectedSptDetail.taxType }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Rekening Bank:</span>
            <span class="font-bold text-zinc-950">{{ selectedSptDetail.bankAccountName }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Dasar Pengenaan Pajak (DPP):</span>
            <span class="font-mono text-zinc-950">{{ formatIDR(selectedSptDetail.dppAmount) }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Tarif:</span>
            <span class="font-bold text-zinc-950">{{ selectedSptDetail.taxRate }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Nilai Pajak Terutang:</span>
            <span class="font-bold text-amber-600 text-sm font-mono">{{ formatIDR(selectedSptDetail.taxAmount) }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Nomor NTPN:</span>
            <span class="font-mono font-bold text-zinc-950">{{ selectedSptDetail.ntpnNumber || '-' }}</span>
          </div>
          <div class="flex justify-between py-1 border-b border-zinc-50 col-span-2">
            <span class="text-zinc-500">Status Penyetoran:</span>
            <span class="font-bold text-emerald-700">{{ selectedSptDetail.status }}</span>
          </div>
        </div>

        <!-- Rincian Transaksi Jurnal Umum Terkait -->
        <div class="pt-2 border-t border-zinc-100 space-y-2">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="text-xs font-bold text-zinc-900">Rincian Transaksi Jurnal Umum</span>
              <span class="px-2 py-0.5 rounded-full bg-blue-50 text-blue-700 font-mono text-[10px] font-bold">
                {{ (selectedSptDetail.linkedJournals || []).length }} Transaksi
              </span>
            </div>
            <router-link
              to="/accounting/journals/create"
              class="text-[11px] font-bold text-blue-600 hover:text-blue-800"
            >
              + Input Jurnal Baru &rarr;
            </router-link>
          </div>

          <div v-if="(selectedSptDetail.linkedJournals || []).length > 0" class="max-h-56 overflow-y-auto border border-zinc-200 rounded-xl">
            <table class="w-full text-left border-collapse text-xs">
              <thead class="bg-zinc-50 sticky top-0 border-b border-zinc-200 text-[10px] font-bold text-zinc-500 uppercase">
                <tr>
                  <th class="px-2.5 py-1.5 w-24">Tanggal</th>
                  <th class="px-2.5 py-1.5 w-32">No. Jurnal</th>
                  <th class="px-2.5 py-1.5">Uraian / Memo</th>
                  <th class="px-2.5 py-1.5 text-right w-24">DPP (Rp)</th>
                  <th class="px-2.5 py-1.5 text-right w-24">Pajak (Rp)</th>
                  <th class="px-2.5 py-1.5 text-center w-16">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-zinc-100">
                <tr v-for="j in selectedSptDetail.linkedJournals" :key="j.id" class="hover:bg-zinc-50">
                  <td class="px-2.5 py-1.5 font-mono text-zinc-500">{{ formatDate(j.journalDate) }}</td>
                  <td
                    class="px-2.5 py-1.5 font-mono font-bold text-blue-600 cursor-pointer hover:underline"
                    @click="showDetailModal = false; openJournalDetail(j)"
                  >
                    {{ j.journalNumber }}
                  </td>
                  <td class="px-2.5 py-1.5 text-zinc-700 truncate max-w-[180px]">{{ j.description }}</td>
                  <td class="px-2.5 py-1.5 text-right font-mono tabular-nums text-zinc-600">{{ formatIDR(j.dppAmount) }}</td>
                  <td class="px-2.5 py-1.5 text-right font-mono tabular-nums font-bold text-emerald-700">{{ formatIDR(j.taxAmount) }}</td>
                  <td class="px-2.5 py-1.5 text-center">
                    <button
                      @click="showDetailModal = false; openJournalDetail(j)"
                      class="px-2 py-0.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-800 rounded font-bold text-[10px] transition cursor-pointer"
                    >
                      Buka
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <p v-else class="text-xs text-zinc-400 italic py-2 text-center bg-zinc-50 rounded-xl border border-zinc-100">
            Belum ada transaksi Jurnal Umum pajak untuk masa ini. Nilai di atas mengacu pada akumulasi mutasi akun.
          </p>
        </div>

        <div class="pt-2 flex justify-end">
          <button
            @click="showDetailModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- DIALOG MODAL 2: Penyetoran Pajak via Bank -->
    <div
      v-if="showPaymentBankModal"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-md w-full overflow-hidden p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div class="flex items-center gap-2">
            <Landmark class="w-4 h-4 text-zinc-950" />
            <h3 class="text-sm font-bold text-zinc-950">Penyetoran Kas Negara (NTPN Bank)</h3>
          </div>
          <button @click="showPaymentBankModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <p class="text-xs text-zinc-600">
          Sesuai aturan PMK DJP, penyetoran pajak negara wajib melalui pendebetan rekening bank operasional resmi klinik:
        </p>

        <div class="space-y-3 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Pilih Rekening Bank Debet</label>
            <select
              v-model="selectedBankForPayment"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900"
            >
              <option
                v-for="b in taxStore.availableBankAccounts"
                :key="b.code"
                :value="b.code"
              >
                {{ b.code }} - {{ b.name }}
              </option>
            </select>
          </div>

          <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200 space-y-1">
            <span class="text-zinc-500 text-[10px] block">Estimasi Jumlah Yang Disetor:</span>
            <span class="text-base font-bold tabular-nums text-zinc-950">
              Rp {{ formatIDR(summary.totalTax || 39249649) }}
            </span>
          </div>
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <button
            @click="showPaymentBankModal = false"
            class="px-4 py-2 bg-white border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-700"
          >
            Batal
          </button>
          <button
            @click="handleExecutePayment"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition"
          >
            Eksekusi Pendebetan & Posting Jurnal
          </button>
        </div>
      </div>
    </div>

    <!-- DIALOG MODAL 3: Detail Voucher Jurnal Umum Terkait -->
    <div
      v-if="showJournalPreviewModal && activeJournal"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-2xl w-full overflow-hidden p-6 space-y-4 animate-in fade-in zoom-in-95 duration-150">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono text-base font-bold text-zinc-950">{{ activeJournal.journalNumber }}</span>
              <span class="px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-800 border border-emerald-200 text-[10px] font-bold">
                {{ activeJournal.status || 'POSTED' }}
              </span>
            </div>
            <p class="text-xs text-zinc-500 mt-0.5">Tanggal: {{ formatDate(activeJournal.journalDate) }} • Unit: {{ activeJournal.businessUnit || 'ALL' }}</p>
          </div>
          <button @click="showJournalPreviewModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="p-3 bg-zinc-50 rounded-xl border border-zinc-200 text-xs">
          <span class="text-zinc-400 font-bold uppercase text-[10px] block">Keterangan / Uraian Jurnal:</span>
          <p class="text-zinc-900 font-medium mt-0.5">{{ activeJournal.description }}</p>
        </div>

        <!-- Multi-line Table -->
        <div class="border border-zinc-200 rounded-xl overflow-hidden bg-white">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-500 uppercase">
                <th class="py-2 px-3">Kode Akun</th>
                <th class="py-2 px-3">Nama Akun</th>
                <th class="py-2 px-3 text-right">Debit (Rp)</th>
                <th class="py-2 px-3 text-right">Kredit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="(line, idx) in activeJournal.lines" :key="idx" class="hover:bg-zinc-50/50">
                <td class="py-2 px-3 font-mono font-bold text-zinc-900">{{ line.accountCode }}</td>
                <td class="py-2 px-3 text-zinc-700">
                  <span class="font-medium block">{{ line.accountName }}</span>
                  <span v-if="line.description" class="text-[10px] text-zinc-400 block">{{ line.description }}</span>
                </td>
                <td class="py-2 px-3 text-right tabular-nums font-semibold text-zinc-950">
                  {{ line.debit > 0 ? formatIDR(line.debit) : '-' }}
                </td>
                <td class="py-2 px-3 text-right tabular-nums font-semibold text-zinc-950">
                  {{ line.credit > 0 ? formatIDR(line.credit) : '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex items-center justify-between pt-2 border-t border-zinc-100 text-xs">
          <router-link
            :to="`/accounting/journals?search=${activeJournal.journalNumber}`"
            class="text-xs font-bold text-zinc-700 hover:text-zinc-950 flex items-center gap-1"
          >
            <span>Buka di Daftar Jurnal Umum</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </router-link>
          <button
            @click="showJournalPreviewModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
    </div>
    <!-- END OF SCREEN VIEW (print:hidden) -->

    <!-- ========================================================================= -->
    <!-- CETAK RESMI FORMULIR REKAPITULASI SPT MASA & BUKU PEMBANTU PAJAK          -->
    <!-- HANYA TAMPIL SAAT CETAK / EXPORT PDF (PRINT MEDIA ONLY)                   -->
    <!-- ========================================================================= -->
    <div class="hidden print:block print-document font-['Times_New_Roman',serif] text-black bg-white">
      <!-- Kop Surat Resmi -->
      <div class="border-b-2 border-black pb-3 mb-4">
        <div class="flex items-center justify-between">
          <div class="space-y-1">
            <h2 class="text-xl font-bold uppercase tracking-wide">{{ currentEntity.title }}</h2>
            <h3 class="text-xs font-semibold text-gray-800">{{ currentEntity.subtitle }}</h3>
            <p class="text-[10px] text-gray-600">{{ currentEntity.address }}</p>
          </div>
          <div class="text-right border border-black p-2 rounded text-xs font-mono">
            <div class="font-bold text-sm">SPT MASA BULANAN</div>
            <div class="text-[9px] text-gray-600">MASA: {{ String(taxStore.selectedMonth).padStart(2, '0') }} / {{ taxStore.selectedYear }}</div>
            <div class="text-[9px] font-bold text-gray-800 mt-0.5">NPWP: {{ currentEntity.npwp }}</div>
          </div>
        </div>
      </div>

      <!-- Judul Dokumen -->
      <div class="text-center my-4 space-y-0.5">
        <h1 class="text-base font-black uppercase underline tracking-wider">
          SURAT PEMBERITAHUAN (SPT) MASA PAJAK & REKAPITULASI BUKU PEMBANTU
        </h1>
        <p class="text-[11px] italic text-gray-700">
          Rekapitulasi Pemotongan Pajak Bulanan Terikat Jurnal Umum & Mutasi Kas Negara (PMK DJP)
        </p>
      </div>

      <!-- Identitas & Parameter -->
      <div class="border border-black p-2.5 text-[11px] mb-4 bg-gray-50/50">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1">
            <div class="flex"><span class="w-36 font-bold">Wajib Pajak:</span><span class="font-bold">{{ currentEntity.title }}</span></div>
            <div class="flex"><span class="w-36 font-bold">NPWP Badan:</span><span class="font-mono font-bold">{{ currentEntity.npwp }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Masa / Tahun Pajak:</span><span>Masa {{ String(taxStore.selectedMonth).padStart(2, '0') }} / {{ taxStore.selectedYear }}</span></div>
          </div>
          <div class="space-y-1">
            <div class="flex"><span class="w-36 font-bold">Entitas Unit Bisnis:</span><span>{{ currentEntity.unitLabel }}</span></div>
            <div class="flex"><span class="w-36 font-bold">Metode Penyetoran:</span><span>Kanal Perbankan Resmi Non-Tunai</span></div>
            <div class="flex"><span class="w-36 font-bold">Total Pajak Disetor:</span><span class="font-bold font-mono">Rp {{ formatIDR(summary.totalTax) }}</span></div>
          </div>
        </div>
      </div>

      <!-- Tabel I: Rekapitulasi Dokumen SPT Masa -->
      <div class="mb-4">
        <h4 class="text-[11px] font-bold uppercase mb-1.5">I. REKAPITULASI DOKUMEN SPT MASA PAJAK (BULANAN)</h4>
        <table class="w-full text-[10px] border-collapse border border-black">
          <thead>
            <tr class="bg-gray-200 text-center font-bold">
              <th class="border border-black px-2 py-1 w-7">No</th>
              <th class="border border-black px-2 py-1 w-28">No. Dokumen SPT</th>
              <th class="border border-black px-2 py-1 w-20">Jenis Pajak</th>
              <th class="border border-black px-2 py-1">Akun Pembukuan & Bank Penyetor</th>
              <th class="border border-black px-2 py-1 w-24 text-right">DPP (Rp)</th>
              <th class="border border-black px-2 py-1 w-24 text-right">Pajak Disetor (Rp)</th>
              <th class="border border-black px-2 py-1 w-28">NTPN DJP</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(spt, idx) in sptMasaList" :key="spt.id">
              <td class="border border-black px-2 py-1 text-center">{{ idx + 1 }}</td>
              <td class="border border-black px-2 py-1 font-mono font-bold">{{ spt.sptNumber }}</td>
              <td class="border border-black px-2 py-1 font-bold text-center">{{ spt.taxType }}</td>
              <td class="border border-black px-2 py-1">
                <div>[{{ spt.accountCode }}] {{ spt.accountName }}</div>
                <div class="text-[9px] text-gray-600">via [{{ spt.bankAccountCode }}] {{ spt.bankAccountName }}</div>
              </td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(spt.dppAmount) }}</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(spt.taxAmount) }}</td>
              <td class="border border-black px-2 py-1 font-mono text-center text-[9px]">{{ spt.ntpnNumber || '-' }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-gray-100 font-bold text-[10px]">
              <td colspan="4" class="border border-black px-3 py-1 text-right uppercase">Total Pajak Masa:</td>
              <td class="border border-black px-2 py-1 text-right font-mono">Rp {{ formatIDR(summary.totalDpp) }}</td>
              <td class="border border-black px-2 py-1 text-right font-mono">Rp {{ formatIDR(summary.totalTax) }}</td>
              <td class="border border-black px-2 py-1 text-center text-[9px]">Lunas Terverifikasi</td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Tabel II: Detail Rincian Per Voucher Jurnal Umum -->
      <div class="mb-4">
        <h4 class="text-[11px] font-bold uppercase mb-1.5">II. DETAIL TRANSAKSI VOUCHER JURNAL UMUM TERIKAT (PER JURNAL)</h4>
        <table class="w-full text-[10px] border-collapse border border-black">
          <thead>
            <tr class="bg-gray-200 text-center font-bold">
              <th class="border border-black px-2 py-1 w-7">No</th>
              <th class="border border-black px-2 py-1 w-20">Tanggal</th>
              <th class="border border-black px-2 py-1 w-28">No. Jurnal</th>
              <th class="border border-black px-2 py-1 w-16">Jenis</th>
              <th class="border border-black px-2 py-1 w-16">Unit</th>
              <th class="border border-black px-2 py-1">Akun Lawan & Uraian Transaksi</th>
              <th class="border border-black px-2 py-1 w-24 text-right">DPP (Rp)</th>
              <th class="border border-black px-2 py-1 w-24 text-right">Pajak (Rp)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(j, idx) in taxStore.filteredLinkedTaxJournalEntries" :key="j.id">
              <td class="border border-black px-2 py-1 text-center">{{ idx + 1 }}</td>
              <td class="border border-black px-2 py-1 font-mono text-center">{{ j.journalDate ? j.journalDate.slice(0, 10) : '-' }}</td>
              <td class="border border-black px-2 py-1 font-mono font-bold">{{ j.journalNumber }}</td>
              <td class="border border-black px-2 py-1 text-center font-bold">{{ j.taxType }}</td>
              <td class="border border-black px-2 py-1 text-center">{{ j.businessUnit }}</td>
              <td class="border border-black px-2 py-1">
                <div class="font-bold text-[9px]">[{{ j.counterpartAccountCode }}] {{ j.counterpartAccountName }}</div>
                <div class="text-[9px] text-gray-700 italic">{{ j.description }}</div>
              </td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(j.dppAmount) }}</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(j.taxAmount) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Lembar Pengesahan Tanda Tangan -->
      <div class="grid grid-cols-2 text-center text-xs pt-4 border-t border-black break-inside-avoid">
        <div class="space-y-14">
          <div>
            <p>Menyetujui & Mengesahkan,</p>
            <p class="font-bold">{{ currentEntity.directorRole }}</p>
          </div>
          <div>
            <p class="font-bold underline">{{ currentEntity.director }}</p>
            <p class="text-[10px] text-gray-600">{{ currentEntity.sip }}</p>
          </div>
        </div>

        <div class="space-y-14">
          <div>
            <p>Semarang, {{ new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) }}</p>
            <p class="font-bold">Finance & Tax Manager</p>
          </div>
          <div>
            <p class="font-bold underline">{{ accountingStore.currentUserName }}</p>
            <p class="text-[10px] text-gray-600">Register Akuntan Negara: CA-78291</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  .break-inside-avoid {
    break-inside: avoid !important;
    page-break-inside: avoid !important;
  }
}
</style>
