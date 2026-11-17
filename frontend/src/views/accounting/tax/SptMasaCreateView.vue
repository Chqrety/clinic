<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, type SptMasaType, type LinkedTaxJournalEntry } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { BusinessUnit } from '@/types/accounting'
import {
  ArrowLeft,
  Calculator,
  Save,
  CheckCircle2,
  Building2,
  Calendar,
  Landmark,
  FileText,
  ShieldCheck,
  AlertCircle,
  ExternalLink,
  Plus,
  X,
  Layers
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

// State
const masaMonth = ref<number>(accountingStore.fiscalMonth || (new Date().getMonth() + 1))
const masaYear = ref<number>(accountingStore.fiscalYear || new Date().getFullYear())
const unit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const taxType = ref<SptMasaType>('PPH_21')
const bankAccountCode = ref<string>('')
const grossAmount = ref<number>(0)
const dppAmount = ref<number>(0)
const taxRate = ref<string>('PMK 168 / 5%')
const taxAmount = ref<number>(0)
const ntpnNumber = ref<string>('')
const notes = ref<string>('')
const status = ref<'DRAFT' | 'PAID' | 'REPORTED'>('REPORTED')

// Journal linking state
const selectedJournalIds = ref<Set<string>>(new Set())
const showJournalPreviewModal = ref<boolean>(false)
const activeJournal = ref<any>(null)

const taxTypeOptions = [
  { value: 'PPH_21', label: 'PPh Pasal 21', sublabel: 'Tenaga Medis Dokter & Karyawan' },
  { value: 'PPH_23', label: 'PPh Pasal 23', sublabel: 'Jasa Medis, Laboratorium, Limbah B3' },
  { value: 'PPH_4_2', label: 'PPh Final Pasal 4(2)', sublabel: 'Sewa Gedung / Fasilitas / Bunga Giro' },
  { value: 'PPH_25', label: 'PPh Pasal 25', sublabel: 'Angsuran Pajak Penghasilan Badan Bulanan' },
  { value: 'PPN', label: 'PPN', sublabel: 'Pajak Pertambahan Nilai (Keluaran/Masukan)' }
]

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals(),
    taxStore.fetchTaxRecords()
  ])

  // Select default bank account if available
  if (taxStore.availableBankAccounts.length > 0) {
    const defaultBank = taxStore.availableBankAccounts.find(b => b.code === '1-1113') || taxStore.availableBankAccounts[0]
    bankAccountCode.value = defaultBank.code
  }
})

// Watch availableBankAccounts to set initial bank code if not set
watch(() => taxStore.availableBankAccounts, (banks) => {
  if (!bankAccountCode.value && banks.length > 0) {
    const defaultBank = banks.find(b => b.code === '1-1113') || banks[0]
    bankAccountCode.value = defaultBank.code
  }
}, { immediate: true })

// Matching journals for selected Masa, Year, and Tax Type
const matchingJournals = computed<LinkedTaxJournalEntry[]>(() => {
  return taxStore.getLinkedJournalsForMonth(masaMonth.value, masaYear.value, taxType.value)
})

// When matching journals change, auto-select all and recalculate
watch(matchingJournals, (journals) => {
  selectedJournalIds.value = new Set(journals.map(j => j.id))
  recalculateFromSelectedJournals()
}, { immediate: true })

function recalculateFromSelectedJournals() {
  const selected = matchingJournals.value.filter(j => selectedJournalIds.value.has(j.id))
  if (selected.length > 0) {
    grossAmount.value = selected.reduce((sum, j) => sum + j.grossAmount, 0)
    dppAmount.value = selected.reduce((sum, j) => sum + j.dppAmount, 0)
    taxAmount.value = selected.reduce((sum, j) => sum + j.taxAmount, 0)
    
    if (taxType.value === 'PPH_21') taxRate.value = 'PMK 168 / TER'
    else if (taxType.value === 'PPH_23') taxRate.value = '2%'
    else if (taxType.value === 'PPH_4_2') taxRate.value = '10%'
    else if (taxType.value === 'PPN') taxRate.value = '11%'
    else if (taxType.value === 'PPH_25') taxRate.value = 'Angsuran Tetap'

    notes.value = `SPT Masa ${taxType.value} Masa ${masaMonth.value}/${masaYear.value} terintegrasi dari ${selected.length} transaksi Jurnal (${selected.map(j => j.journalNumber).join(', ')})`
  } else if (matchingJournals.value.length === 0) {
    // Fallback when no journals in DB for this tax type
    grossAmount.value = 0
    dppAmount.value = 0
    taxAmount.value = 0
    updateSuggestedTaxAmount()
  }
}

function toggleSelectJournal(id: string) {
  if (selectedJournalIds.value.has(id)) {
    selectedJournalIds.value.delete(id)
  } else {
    selectedJournalIds.value.add(id)
  }
  recalculateFromSelectedJournals()
}

const isAllSelected = computed(() => {
  return matchingJournals.value.length > 0 && selectedJournalIds.value.size === matchingJournals.value.length
})

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedJournalIds.value.clear()
  } else {
    selectedJournalIds.value = new Set(matchingJournals.value.map(j => j.id))
  }
  recalculateFromSelectedJournals()
}

function openJournalPreview(entry: LinkedTaxJournalEntry) {
  if (entry.rawJournal) {
    activeJournal.value = entry.rawJournal
  } else {
    const found = journalsStore.journals.find(j => j.journalNumber === entry.journalNumber || j.id === entry.journalId)
    activeJournal.value = found || {
      journalNumber: entry.journalNumber,
      journalDate: entry.journalDate,
      description: entry.description,
      lines: [
        { accountCode: entry.taxAccountCode, accountName: entry.taxAccountName, debit: 0, credit: entry.taxAmount },
        { accountCode: entry.counterpartAccountCode, accountName: entry.counterpartAccountName, debit: entry.taxAmount, credit: 0 }
      ]
    }
  }
  showJournalPreviewModal.value = true
}

function updateSuggestedTaxAmount() {
  const dpp = Number(dppAmount.value) || 0
  if (taxType.value === 'PPH_23') {
    taxRate.value = '2%'
    taxAmount.value = Math.round(dpp * 0.02)
  } else if (taxType.value === 'PPH_4_2') {
    taxRate.value = '10%'
    taxAmount.value = Math.round(dpp * 0.10)
  } else if (taxType.value === 'PPH_25') {
    taxRate.value = 'Angsuran Tetap'
    taxAmount.value = 4500000
  } else if (taxType.value === 'PPN') {
    taxRate.value = '11%'
    taxAmount.value = Math.round(dpp * 0.11)
  } else if (taxType.value === 'PPH_21') {
    taxRate.value = 'PMK 168 / TER'
    taxAmount.value = Math.round(dpp * 0.05)
  }
}

function handleSave() {
  if (dppAmount.value <= 0 || taxAmount.value <= 0) {
    accountingStore.showToast('Nilai DPP dan Pajak Terutang harus lebih dari Rp 0', 'error')
    return
  }

  const selectedJournalsList = matchingJournals.value.filter(j => selectedJournalIds.value.has(j.id))

  taxStore.createSptMasa({
    masaMonth: Number(masaMonth.value),
    masaYear: Number(masaYear.value),
    taxType: taxType.value,
    bankAccountCode: bankAccountCode.value || '1-1113',
    grossAmount: Number(grossAmount.value),
    dppAmount: Number(dppAmount.value),
    taxRate: taxRate.value,
    taxAmount: Number(taxAmount.value),
    notes: `[Unit ${unit.value}] ${notes.value || ''}`.trim(),
    ntpnNumber: ntpnNumber.value.trim(),
    status: status.value,
    businessUnit: unit.value,
    linkedJournals: selectedJournalsList
  })

  accountingStore.showToast(`SPT Masa ${taxType.value} Unit ${unit.value} Masa ${masaMonth.value}/${masaYear.value} berhasil disimpan!`, 'success')
  router.push('/accounting/tax/spt-masa')
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif] max-w-4xl mx-auto pb-12">
    <!-- Breadcrumb & Back Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
          <router-link to="/accounting/tax/spt-masa" class="hover:text-zinc-950 transition flex items-center gap-1">
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>SPT Masa</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-950">Formulir Rekam Baru</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight">
          Perekaman Dokumen SPT Masa Baru
        </h1>
        <p class="text-xs text-zinc-500">
          Mencatat kepatuhan pajak bulanan klinik (PPh 21, PPh 23, PPh 4(2), PPh 25, & PPN) terverifikasi kanal perbankan resmi dan terikat langsung ke Jurnal Umum.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <router-link
          to="/accounting/tax/spt-masa"
          class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 transition"
        >
          Batal
        </router-link>
        <button
          @click="handleSave"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan Dokumen SPT Masa</span>
        </button>
      </div>
    </header>

    <!-- Form Container -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-6">
      <!-- Section 1: Periode & Klasifikasi -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Calendar class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            1. Periode Masa & Unit Bisnis
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Masa Pajak (Bulan)</label>
            <select
              v-model.number="masaMonth"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
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

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Tahun Pajak</label>
            <input
              v-model.number="masaYear"
              type="number"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold font-mono text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Unit Bisnis Klinik</label>
            <select
              v-model="unit"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option value="ALL">Klinik Terpadu (Konsolidasi)</option>
              <option value="KIDINEST">Unit Kidinest (Klinik Anak)</option>
              <option value="SSC">Unit SSC (Semarang Skin Center)</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Section 2: Transaksi Jurnal Umum & Bukti Potong Terkait -->
      <div class="space-y-3">
        <div class="flex items-center justify-between border-b border-zinc-100 pb-2">
          <div class="flex items-center gap-2">
            <FileText class="w-4 h-4 text-emerald-700" />
            <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
              2. Transaksi Jurnal Umum Pajak Masa Ini (Terintegrasi Otomatis)
            </h3>
            <span class="px-2 py-0.5 bg-emerald-50 text-emerald-800 text-[10px] font-bold rounded-full border border-emerald-200">
              {{ matchingJournals.length }} Transaksi
            </span>
          </div>

          <div class="flex items-center gap-2">
            <button
              v-if="matchingJournals.length > 0"
              type="button"
              @click="toggleSelectAll"
              class="text-[11px] font-bold text-blue-600 hover:text-blue-800 transition cursor-pointer"
            >
              {{ isAllSelected ? 'Batal Pilih Semua' : 'Pilih Semua Transaksi' }}
            </button>
            <router-link
              to="/accounting/journals/create"
              target="_blank"
              class="text-[11px] font-bold text-zinc-600 hover:text-zinc-950 inline-flex items-center gap-1"
            >
              <Plus class="w-3 h-3" />
              <span>+ Buat Jurnal Baru</span>
            </router-link>
          </div>
        </div>

        <!-- Table of matching journals -->
        <div v-if="matchingJournals.length > 0" class="border border-zinc-200 rounded-xl overflow-hidden shadow-2xs">
          <table class="w-full text-left text-xs">
            <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase tracking-wider">
              <tr>
                <th class="px-3 py-2 w-10 text-center">Pilih</th>
                <th class="px-3 py-2 w-24">Tanggal</th>
                <th class="px-3 py-2 w-36">No. Jurnal</th>
                <th class="px-3 py-2">Uraian / Memo Jurnal</th>
                <th class="px-3 py-2">Akun Lawan / Beban</th>
                <th class="px-3 py-2 text-right w-28">DPP (Rp)</th>
                <th class="px-3 py-2 text-right w-28">Pajak (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr
                v-for="j in matchingJournals"
                :key="j.id"
                :class="[
                  'hover:bg-zinc-50/80 transition cursor-pointer',
                  selectedJournalIds.has(j.id) ? 'bg-blue-50/20' : ''
                ]"
                @click="toggleSelectJournal(j.id)"
              >
                <td class="px-3 py-2.5 text-center" @click.stop="toggleSelectJournal(j.id)">
                  <input
                    type="checkbox"
                    :checked="selectedJournalIds.has(j.id)"
                    class="w-4 h-4 rounded text-zinc-950 focus:ring-0 cursor-pointer"
                  />
                </td>
                <td class="px-3 py-2.5 font-mono text-zinc-600">{{ formatDate(j.journalDate) }}</td>
                <td
                  class="px-3 py-2.5 font-mono font-bold text-blue-600 hover:underline flex items-center gap-1"
                  @click.stop="openJournalPreview(j)"
                  title="Klik untuk melihat voucher jurnal lengkap"
                >
                  <span>{{ j.journalNumber }}</span>
                  <ExternalLink class="w-2.5 h-2.5 text-zinc-400" />
                </td>
                <td class="px-3 py-2.5 font-medium text-zinc-900">{{ j.description }}</td>
                <td class="px-3 py-2.5 text-zinc-600 text-[11px]">
                  <span class="block font-medium truncate max-w-[200px]">{{ j.counterpartAccountName }}</span>
                  <span class="font-mono text-zinc-400 text-[10px]">{{ j.counterpartAccountCode }}</span>
                </td>
                <td class="px-3 py-2.5 text-right font-mono tabular-nums text-zinc-700">
                  {{ formatIDR(j.dppAmount, false) }}
                </td>
                <td class="px-3 py-2.5 text-right font-mono tabular-nums font-bold text-emerald-700">
                  {{ formatIDR(j.taxAmount, false) }}
                </td>
              </tr>
            </tbody>
            <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs">
              <tr>
                <td colspan="5" class="px-3 py-2 text-right uppercase tracking-wider text-zinc-600">
                  Total Terpilih ({{ selectedJournalIds.size }} Transaksi):
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums text-zinc-950 font-bold">
                  {{ formatIDR(dppAmount, false) }}
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums text-emerald-800 font-extrabold">
                  {{ formatIDR(taxAmount, false) }}
                </td>
              </tr>
            </tfoot>
          </table>
        </div>

        <!-- Alert when no journals found for this type in this period -->
        <div v-else class="p-4 rounded-xl bg-amber-50/70 border border-amber-200 text-amber-900 flex items-start gap-3 text-xs">
          <AlertCircle class="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
          <div class="space-y-1">
            <p class="font-bold">
              Belum ada voucher Jurnal Umum dengan akun {{ taxType }} pada Masa {{ masaMonth }}/{{ masaYear }}.
            </p>
            <p class="text-[11px] text-amber-800">
              Nilai DPP dan Pajak Terutang dapat diisi secara manual di bawah, atau Anda dapat membuat Jurnal Umum terlebih dahulu agar terhubung ke Buku Besar.
            </p>
          </div>
        </div>
      </div>

      <!-- Section 3: Jenis Pajak & Perhitungan -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Calculator class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            3. Jenis Pajak & Nilai Terutang
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Jenis SPT Masa</label>
            <SearchableSelect
              :model-value="taxType"
              :options="taxTypeOptions"
              placeholder="Pilih atau cari jenis SPT..."
              @update:model-value="(val: any) => { taxType = val; }"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Tarif Pajak Efektif / Persentase</label>
            <input
              v-model="taxRate"
              type="text"
              placeholder="Contoh: 2%, 10%, 11%, TER"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Nilai Bruto / Objek Transaksi (Rp)</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-xs">Rp</span>
              <input
                v-model.number="grossAmount"
                type="number"
                class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">DPP (Dasar Pengenaan Pajak) (Rp)</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-xs">Rp</span>
              <input
                v-model.number="dppAmount"
                @input="updateSuggestedTaxAmount"
                type="number"
                class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              />
            </div>
          </div>
        </div>

        <!-- Calculated Tax Amount Banner -->
        <div class="p-4 rounded-xl bg-zinc-900 text-white flex items-center justify-between shadow-2xs">
          <div>
            <span class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider block">
              Total Pajak Terutang (SPT Masa)
              <span v-if="selectedJournalIds.size > 0" class="text-amber-300 font-normal ml-1">
                ({{ selectedJournalIds.size }} transaksi jurnal dipilih)
              </span>
            </span>
            <span class="text-xl font-bold tabular-nums text-amber-300 font-mono">
              {{ formatIDR(taxAmount) }}
            </span>
          </div>
          <div class="w-48">
            <label class="block text-[10px] font-semibold text-zinc-400 mb-1">Koreksi Manual Nominal (Rp):</label>
            <input
              v-model.number="taxAmount"
              type="number"
              class="w-full px-2.5 py-1 bg-zinc-800 border border-zinc-700 rounded-lg text-xs font-mono font-bold text-white focus:outline-none"
            />
          </div>
        </div>
      </div>

      <!-- Section 4: Kanal Bank & NTPN -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Landmark class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            4. Kanal Perbankan & NTPN (Strictly Non-Cash)
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div class="sm:col-span-2">
            <label class="block font-bold text-zinc-700 mb-1">Rekening Bank Penyetor Resmi</label>
            <select
              v-model="bankAccountCode"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option
                v-for="b in taxStore.availableBankAccounts"
                :key="b.code"
                :value="b.code"
              >
                [{{ b.code }}] {{ b.name }}
              </option>
            </select>
            <span class="text-[10px] text-zinc-400 mt-0.5 block">
              Sesuai PMK DJP: Penyetoran SPT Masa wajib melalui rekening bank operasional klinik, bukan kas tunai.
            </span>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Status Pelaporan</label>
            <select
              v-model="status"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option value="DRAFT">DRAFT (Belum Bayar)</option>
              <option value="PAID">PAID (Sudah Terdebit Bank)</option>
              <option value="REPORTED">REPORTED (Lengkap NTPN / BPE)</option>
            </select>
          </div>

          <div class="sm:col-span-2">
            <label class="block font-bold text-zinc-700 mb-1">Nomor NTPN / Kode Billing DJP</label>
            <input
              v-model="ntpnNumber"
              type="text"
              placeholder="Contoh: NTPN-BCA-20261120-9901 atau e-Billing DJP"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div class="sm:col-span-3">
            <label class="block font-bold text-zinc-700 mb-1">Keterangan / Uraian Transaksi</label>
            <input
              v-model="notes"
              type="text"
              placeholder="Contoh: Pemotongan Pajak Masa November Terintegrasi Jurnal"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Preview Voucher Jurnal -->
    <div
      v-if="showJournalPreviewModal"
      class="fixed inset-0 z-50 bg-black/50 backdrop-blur-xs flex items-center justify-center p-4"
      @click.self="showJournalPreviewModal = false"
    >
      <div class="bg-white rounded-2xl max-w-2xl w-full p-5 space-y-4 shadow-xl">
        <div class="flex items-center justify-between border-b border-zinc-200 pb-3">
          <div class="flex items-center gap-2">
            <Layers class="w-4 h-4 text-blue-600" />
            <h3 class="font-bold text-sm text-zinc-950">
              Voucher Jurnal Umum: {{ activeJournal?.journalNumber }}
            </h3>
          </div>
          <button @click="showJournalPreviewModal = false" class="text-zinc-400 hover:text-zinc-950 transition cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="text-xs space-y-1 text-zinc-600">
          <p><strong class="text-zinc-900">Tanggal:</strong> {{ formatDate(activeJournal?.journalDate || activeJournal?.transaction_date) }}</p>
          <p><strong class="text-zinc-900">Keterangan:</strong> {{ activeJournal?.description }}</p>
        </div>

        <div class="border border-zinc-200 rounded-xl overflow-hidden text-xs">
          <table class="w-full text-left">
            <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase">
              <tr>
                <th class="px-3 py-2">Akun</th>
                <th class="px-3 py-2 text-right w-28">Debet (Rp)</th>
                <th class="px-3 py-2 text-right w-28">Kredit (Rp)</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="(l, idx) in (activeJournal?.lines || [])" :key="idx">
                <td class="px-3 py-2">
                  <span class="font-bold font-mono text-zinc-900">{{ l.accountCode || l.account_code }}</span> -
                  <span class="text-zinc-700">{{ l.accountName || l.account_name }}</span>
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums">
                  {{ Number(l.debit) > 0 ? formatIDR(Number(l.debit), false) : '-' }}
                </td>
                <td class="px-3 py-2 text-right font-mono tabular-nums">
                  {{ Number(l.credit) > 0 ? formatIDR(Number(l.credit), false) : '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="flex justify-end pt-2">
          <button
            @click="showJournalPreviewModal = false"
            class="px-4 py-2 bg-zinc-950 text-white rounded-xl text-xs font-bold hover:bg-zinc-800 transition cursor-pointer"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
