<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, type SptMasaType } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
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
  Database
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()

onMounted(async () => {
  await Promise.all([
    taxStore.fetchTaxRecords(),
    ledgerStore.fetchAccounts()
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

// Subview under SPT Masa: 'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU'
const sptMasaSubTab = ref<'DOKUMEN_SPT' | 'EBUPOT_DOKTER' | 'BUKU_PEMBANTU'>('DOKUMEN_SPT')

const showDetailModal = ref(false)
const showPaymentBankModal = ref(false)
const selectedSptDetail = ref<any>(null)
const selectedBankForPayment = ref<string>('1102.001')

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
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
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
          <span>+ Buat SPT Masa Baru</span>
        </router-link>

        <button
          @click="showPaymentBankModal = true"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Landmark class="w-3.5 h-3.5 text-zinc-600" />
          <span>Penyetoran via Bank (NTPN)</span>
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

        <div class="text-xs text-zinc-500 font-medium">
          Ditemukan <strong>{{ sptMasaList.length }}</strong> dokumen SPT Masa
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
              <tr v-for="spt in sptMasaList" :key="spt.id" class="hover:bg-zinc-50/70 transition">
                <td class="px-4 py-2.5 font-mono font-bold text-zinc-900">
                  {{ spt.sptNumber }}
                </td>
                <td class="px-4 py-2.5">
                  <span class="font-bold text-zinc-950">{{ spt.taxType }}</span>
                  <span class="block text-[10px] text-zinc-500 truncate max-w-xs">{{ spt.accountName }}</span>
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
                    class="px-2.5 py-1 bg-white hover:bg-zinc-100 text-zinc-700 hover:text-zinc-950 rounded-lg text-xs font-bold border border-zinc-200 transition"
                  >
                    Detail
                  </button>
                </td>
              </tr>
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
    <div v-else class="space-y-4">
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-900">Buku Pembantu Mutasi Transaksi Pajak</h3>
            <p class="text-xs text-zinc-500 mt-0.5">Rincian invoice operasional yang memicu timbulnya hutang pajak klinik.</p>
          </div>
          <div class="w-80">
            <SearchableSelect
              v-model="taxStore.selectedTaxAccountCode"
              :options="taxAccountOptions"
              placeholder="Cari akun pajak..."
              input-class="!py-1.5"
            />
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

    <!-- DIALOG MODAL 1: Detail SPT Masa -->
    <div
      v-if="showDetailModal && selectedSptDetail"
      class="fixed inset-0 z-50 overflow-y-auto bg-zinc-950/60 backdrop-blur-xs flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl shadow-2xl border border-zinc-200/80 max-w-lg w-full overflow-hidden animate-in fade-in zoom-in-95 duration-150 p-6 space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-3">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">{{ selectedSptDetail.sptNumber }}</h3>
            <span class="text-xs text-zinc-500">Masa Pajak {{ selectedSptDetail.masaMonth }}/{{ selectedSptDetail.masaYear }}</span>
          </div>
          <button @click="showDetailModal = false" class="text-zinc-400 hover:text-zinc-600">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-2 text-xs">
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
          <div class="flex justify-between py-1 border-b border-zinc-50">
            <span class="text-zinc-500">Status:</span>
            <span class="font-bold text-emerald-700">{{ selectedSptDetail.status }}</span>
          </div>
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
  </div>
</template>
