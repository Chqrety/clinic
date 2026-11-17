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
  Trash2,
  X,
  Layers,
  Sparkles,
  RefreshCw,
  BookOpen,
  Check,
  CreditCard
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

// Header Form State (Metadata Dokumen SPT Masa)
const masaMonth = ref<number>(accountingStore.fiscalMonth || (new Date().getMonth() + 1))
const masaYear = ref<number>(accountingStore.fiscalYear || new Date().getFullYear())
const unit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const taxType = ref<SptMasaType>('PPH_23')
const bankAccountCode = ref<string>('1-1113')
const ntpnNumber = ref<string>('')
const notes = ref<string>('')
const status = ref<'DRAFT' | 'PAID' | 'REPORTED'>('REPORTED')

// Multi-line Voucher Table Grid State (Identik dengan Jurnal Umum)
export interface SptMasaRow {
  id: string
  selected: boolean
  isFromJournal: boolean
  journalId?: string
  journalNumber?: string
  date: string
  counterpartCode: string
  counterpartName: string
  taxAccountCode: string
  taxAccountName: string
  description: string
  grossAmount: number
  dppAmount: number
  taxRate: string
  taxAmount: number
  rawJournal?: any
}

const rows = ref<SptMasaRow[]>([])
const showJournalPreviewModal = ref<boolean>(false)
const activeJournal = ref<any>(null)

const taxTypeOptions = [
  { value: 'PPH_21', label: 'PPh Pasal 21', sublabel: 'Tenaga Medis Dokter & Karyawan' },
  { value: 'PPH_23', label: 'PPh Pasal 23', sublabel: 'Jasa Medis, Laboratorium, Limbah B3' },
  { value: 'PPH_4_2', label: 'PPh Final Pasal 4(2)', sublabel: 'Sewa Gedung / Fasilitas / Bunga Giro' },
  { value: 'PPH_25', label: 'PPh Pasal 25', sublabel: 'Angsuran Pajak Penghasilan Badan Bulanan' },
  { value: 'PPN', label: 'PPN', sublabel: 'Pajak Pertambahan Nilai (Keluaran/Masukan)' }
]

// COA accounts for searchable select in table
const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

const accountSelectOptions = computed(() => {
  return availableAccounts.value.map(acc => ({
    value: acc.code,
    label: `[${acc.code}] ${acc.name}`,
    sublabel: `${acc.type} • Saldo Normal: ${acc.normalBalance}`
  }))
})

function getDefaultTaxAccount(tType: SptMasaType): { code: string; name: string } {
  if (tType === 'PPH_23') return { code: '2-1510', name: 'Kredit Pajak PPh 23' }
  if (tType === 'PPH_25') return { code: '2-1520', name: 'Uang Muka PPh 25 Badan' }
  if (tType === 'PPN') return { code: '2-1500', name: 'Hutang Pajak (PPN)' }
  if (tType === 'PPH_4_2') return { code: '2-1500', name: 'Hutang Pajak (PPh Final 4(2))' }
  return { code: '2-1500', name: 'Hutang Pajak (PPh 21 Dokter/Karyawan)' }
}

function getDefaultRate(tType: SptMasaType): string {
  if (tType === 'PPH_23') return '2%'
  if (tType === 'PPH_4_2') return '10%'
  if (tType === 'PPN') return '11%'
  if (tType === 'PPH_25') return 'Tetap'
  return 'PMK 168 / 5%'
}

// Matching journals in store for current period and tax type
const matchingJournalsInStore = computed(() => {
  return taxStore.getLinkedJournalsForMonth(masaMonth.value, masaYear.value, taxType.value)
})

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

  // Auto-pull on initial mount
  pullDataFromGeneralJournal(false)
})

// Watch filters to pull data or reset
watch([masaMonth, masaYear, taxType], () => {
  pullDataFromGeneralJournal(false)
})

/**
 * ⚡ Hero Automation: Ambil Data Otomatis dari Jurnal Umum (Brayo-Style Withholding Tax Extraction)
 */
function pullDataFromGeneralJournal(showNotice = true) {
  const journals = matchingJournalsInStore.value
  const defTax = getDefaultTaxAccount(taxType.value)
  const defRate = getDefaultRate(taxType.value)

  if (journals.length > 0) {
    rows.value = journals.map((j, idx) => ({
      id: `row-ju-${j.id}-${idx}`,
      selected: true,
      isFromJournal: true,
      journalId: j.journalId,
      journalNumber: j.journalNumber,
      date: j.journalDate ? j.journalDate.slice(0, 10) : new Date().toISOString().slice(0, 10),
      counterpartCode: j.counterpartAccountCode || '6-1300',
      counterpartName: j.counterpartAccountName || 'Beban Operasional Terkait',
      taxAccountCode: j.taxAccountCode || defTax.code,
      taxAccountName: j.taxAccountName || defTax.name,
      description: j.description || `Pemotongan ${taxType.value} Masa ${masaMonth.value}/${masaYear.value}`,
      grossAmount: j.grossAmount || j.dppAmount,
      dppAmount: j.dppAmount,
      taxRate: defRate,
      taxAmount: j.taxAmount,
      rawJournal: j.rawJournal
    }))

    notes.value = `SPT Masa ${taxType.value} Masa ${masaMonth.value}/${masaYear.value} terintegrasi dari ${rows.value.length} voucher Jurnal Umum (${rows.value.map(r => r.journalNumber).join(', ')})`
    if (showNotice) {
      accountingStore.showToast(`Berhasil menarik ${rows.value.length} transaksi pajak dari Jurnal Umum!`, 'success')
    }
  } else {
    // If no matching journals, create 1 empty manual row so user can input directly
    rows.value = [
      {
        id: `row-manual-1`,
        selected: true,
        isFromJournal: false,
        date: new Date().toISOString().slice(0, 10),
        counterpartCode: '6-1300',
        counterpartName: 'Beban Listrik, Air & Telepon',
        taxAccountCode: defTax.code,
        taxAccountName: defTax.name,
        description: `Pemotongan ${taxType.value} Periode Masa ${masaMonth.value}/${masaYear.value}`,
        grossAmount: 0,
        dppAmount: 0,
        taxRate: defRate,
        taxAmount: 0
      }
    ]
    notes.value = `SPT Masa ${taxType.value} Masa ${masaMonth.value}/${masaYear.value} (Input Manual Voucher)`
    if (showNotice) {
      accountingStore.showToast(`Belum ada jurnal ${taxType.value} pada Masa ${masaMonth.value}/${masaYear.value}. Formulir siap diinput manual.`, 'info')
    }
  }
}

// Add manual row like JournalCreateView
function addRow() {
  const defTax = getDefaultTaxAccount(taxType.value)
  const defRate = getDefaultRate(taxType.value)
  rows.value.push({
    id: `row-manual-${Date.now()}`,
    selected: true,
    isFromJournal: false,
    date: new Date().toISOString().slice(0, 10),
    counterpartCode: '6-1200',
    counterpartName: 'Beban Operasional Klinik',
    taxAccountCode: defTax.code,
    taxAccountName: defTax.name,
    description: `Transaksi Pajak Tambahan ${taxType.value}`,
    grossAmount: 0,
    dppAmount: 0,
    taxRate: defRate,
    taxAmount: 0
  })
}

function removeRow(index: number) {
  if (rows.value.length > 1) {
    rows.value.splice(index, 1)
  } else {
    rows.value[0].grossAmount = 0
    rows.value[0].dppAmount = 0
    rows.value[0].taxAmount = 0
    rows.value[0].description = ''
  }
}

// Auto recalculate tax amount on row DPP or Rate change
function onRowDppChange(row: SptMasaRow) {
  const dpp = Number(row.dppAmount) || 0
  if (!row.grossAmount || row.grossAmount < dpp) {
    row.grossAmount = dpp
  }
  calculateRowTax(row)
}

function calculateRowTax(row: SptMasaRow) {
  const dpp = Number(row.dppAmount) || 0
  const rateStr = (row.taxRate || '').toLowerCase().trim()

  if (rateStr.includes('2%') || rateStr === '2') {
    row.taxAmount = Math.round(dpp * 0.02)
  } else if (rateStr.includes('10%') || rateStr === '10') {
    row.taxAmount = Math.round(dpp * 0.10)
  } else if (rateStr.includes('11%') || rateStr === '11') {
    row.taxAmount = Math.round(dpp * 0.11)
  } else if (rateStr.includes('5%') || rateStr === '5') {
    row.taxAmount = Math.round(dpp * 0.05)
  } else if (rateStr.includes('tetap')) {
    row.taxAmount = row.taxAmount || 4500000
  }
}

function onCounterpartSelect(row: SptMasaRow, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc) {
    row.counterpartCode = acc.code
    row.counterpartName = acc.name
  }
}

function onTaxAccountSelect(row: SptMasaRow, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc) {
    row.taxAccountCode = acc.code
    row.taxAccountName = acc.name
  }
}

// Selections
const selectedRows = computed(() => rows.value.filter(r => r.selected))

const isAllSelected = computed(() => {
  return rows.value.length > 0 && rows.value.every(r => r.selected)
})

function toggleSelectAll() {
  const target = !isAllSelected.value
  rows.value.forEach(r => { r.selected = target })
}

// Totals
const totalGross = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.grossAmount) || 0), 0)
})

const totalDpp = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.dppAmount) || 0), 0)
})

const totalTax = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.taxAmount) || 0), 0)
})

// Bank account info for settlement journal preview
const activeBankAccount = computed(() => {
  return taxStore.availableBankAccounts.find(b => b.code === bankAccountCode.value) || {
    code: '1-1113',
    name: 'Bank BCA Operasional'
  }
})

const primaryTaxAccount = computed(() => {
  const first = selectedRows.value[0]
  if (first?.taxAccountCode) {
    return { code: first.taxAccountCode, name: first.taxAccountName }
  }
  return getDefaultTaxAccount(taxType.value)
})

// Journal Preview modal
function openJournalPreview(row: SptMasaRow) {
  if (row.rawJournal) {
    activeJournal.value = row.rawJournal
  } else {
    const found = journalsStore.journals.find(j => j.journalNumber === row.journalNumber || j.id === row.journalId)
    activeJournal.value = found || {
      journalNumber: row.journalNumber || 'JV-PREVIEW',
      journalDate: row.date,
      description: row.description,
      lines: [
        { accountCode: row.counterpartCode, accountName: row.counterpartName, debit: row.taxAmount, credit: 0 },
        { accountCode: row.taxAccountCode, accountName: row.taxAccountName, debit: 0, credit: row.taxAmount }
      ]
    }
  }
  showJournalPreviewModal.value = true
}

// Save action
function handleSave() {
  if (totalTax.value <= 0) {
    accountingStore.showToast('Total Pajak Terutang harus lebih dari Rp 0 untuk dapat merekam SPT Masa', 'error')
    return
  }

  const linkedJournalsList = selectedRows.value
    .filter(r => r.isFromJournal)
    .map(r => ({
      id: r.id,
      journalId: r.journalId || '',
      journalNumber: r.journalNumber || '',
      journalDate: r.date,
      description: r.description,
      businessUnit: unit.value,
      status: 'POSTED',
      taxType: taxType.value,
      taxAccountCode: r.taxAccountCode,
      taxAccountName: r.taxAccountName,
      counterpartAccountCode: r.counterpartCode,
      counterpartAccountName: r.counterpartName,
      grossAmount: r.grossAmount,
      dppAmount: r.dppAmount,
      taxAmount: r.taxAmount,
      direction: 'WITHHELD' as const,
      bankAccount: `[${activeBankAccount.value.code}] ${activeBankAccount.value.name}`,
      month: masaMonth.value,
      year: masaYear.value,
      rawJournal: r.rawJournal
    }))

  taxStore.createSptMasa({
    masaMonth: Number(masaMonth.value),
    masaYear: Number(masaYear.value),
    taxType: taxType.value,
    bankAccountCode: bankAccountCode.value || '1-1113',
    grossAmount: totalGross.value,
    dppAmount: totalDpp.value,
    taxRate: getDefaultRate(taxType.value),
    taxAmount: totalTax.value,
    notes: `[Unit ${unit.value}] ${notes.value || ''}`.trim(),
    ntpnNumber: ntpnNumber.value.trim(),
    status: status.value,
    businessUnit: unit.value,
    linkedJournals: linkedJournalsList
  })

  // Also auto-post settlement journal if status is PAID or REPORTED
  if (status.value === 'REPORTED' || status.value === 'PAID') {
    journalsStore.createManualJournal({
      journalDate: `${masaYear.value}-${String(masaMonth.value).padStart(2, '0')}-20`,
      journalSource: 'MANUAL',
      description: `Penyetoran SPT Masa ${taxType.value} Masa ${masaMonth.value}/${masaYear.value} via ${activeBankAccount.value.name} (NTPN: ${ntpnNumber.value || 'DJP-SETOR'})`,
      businessUnit: unit.value,
      lines: [
        {
          accountCode: primaryTaxAccount.value.code,
          accountName: primaryTaxAccount.value.name,
          description: `Pelunasan Hutang Pajak SPT Masa ${taxType.value}`,
          debit: totalTax.value,
          credit: 0
        },
        {
          accountCode: activeBankAccount.value.code,
          accountName: activeBankAccount.value.name,
          description: `Penyetoran Pajak via Bank (NTPN: ${ntpnNumber.value || 'DJP-SETOR'})`,
          debit: 0,
          credit: totalTax.value
        }
      ]
    })
  }

  accountingStore.showToast(`SPT Masa ${taxType.value} Masa ${masaMonth.value}/${masaYear.value} berhasil disimpan dan terekam di Jurnal Umum!`, 'success')
  router.push('/accounting/tax/spt-masa')
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] max-w-7xl mx-auto pb-16">
    <!-- Header & Action Navigation -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
          <router-link to="/accounting/tax/spt-masa" class="hover:text-zinc-950 transition flex items-center gap-1">
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>SPT Masa</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-950">Formulir Voucher Rekam Baru</span>
          <span class="text-zinc-300">•</span>
          <span class="text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-full text-[10px] font-bold border border-emerald-200">
            Standar Brayo ERP
          </span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight flex items-center gap-2">
          <span>Perekaman Dokumen SPT Masa & Voucher Pajak</span>
        </h1>
        <p class="text-xs text-zinc-500">
          Format input baris voucher Jurnal Umum dengan fitur penarikan otomatis dari transaksi Jurnal Umum (Withholding Tax Brayo/ERPNext).
        </p>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <router-link
          to="/accounting/tax/spt-masa"
          class="px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 transition"
        >
          Batal
        </router-link>

        <button
          type="button"
          @click="() => pullDataFromGeneralJournal(true)"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-amber-500 hover:bg-amber-600 text-zinc-950 rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Sparkles class="w-3.5 h-3.5 text-zinc-950" />
          <span>⚡ Ambil Data Otomatis dari Jurnal Umum</span>
          <span v-if="matchingJournalsInStore.length > 0" class="px-1.5 py-0.2 bg-black/10 rounded-full text-[10px]">
            {{ matchingJournalsInStore.length }}
          </span>
        </button>

        <button
          type="button"
          @click="handleSave"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan Dokumen SPT Masa</span>
        </button>
      </div>
    </header>

    <!-- Header Form Card (Voucher Metadata) -->
    <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 text-xs">
      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Masa Pajak (Bulan):</span>
        </label>
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
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Tahun Pajak:</span>
        </label>
        <input
          v-model.number="masaYear"
          type="number"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-zinc-500" />
          <span>Unit Bisnis Klinik:</span>
        </label>
        <select
          v-model="unit"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Klinik Terpadu (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Klinik Anak)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <FileText class="w-3.5 h-3.5 text-zinc-500" />
          <span>Jenis SPT Masa:</span>
        </label>
        <select
          v-model="taxType"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="PPH_21">PPh Pasal 21 (Dokter & Karyawan)</option>
          <option value="PPH_23">PPh Pasal 23 (Jasa Limbah & Medis)</option>
          <option value="PPH_4_2">PPh Final Pasal 4(2) (Sewa Gedung)</option>
          <option value="PPH_25">PPh Pasal 25 (Angsuran Badan Bulanan)</option>
          <option value="PPN">PPN (Pajak Pertambahan Nilai)</option>
        </select>
      </div>

      <!-- Second Row Metadata -->
      <div class="sm:col-span-2">
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Landmark class="w-3.5 h-3.5 text-zinc-500" />
          <span>Rekening Bank Penyetor Resmi (PMK DJP Non-Cash):</span>
        </label>
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
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <CreditCard class="w-3.5 h-3.5 text-zinc-500" />
          <span>Nomor NTPN / Kode Billing DJP:</span>
        </label>
        <input
          v-model="ntpnNumber"
          type="text"
          placeholder="e.g. NTPN-BCA-20261120-9901"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <ShieldCheck class="w-3.5 h-3.5 text-zinc-500" />
          <span>Status Pelaporan:</span>
        </label>
        <select
          v-model="status"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="DRAFT">DRAFT (Belum Bayar)</option>
          <option value="PAID">PAID (Sudah Terdebit Bank)</option>
          <option value="REPORTED">REPORTED (Lengkap NTPN / BPE)</option>
        </select>
      </div>

      <div class="sm:col-span-4">
        <label class="block font-bold text-zinc-700 mb-1">Catatan Dokumen / Uraian Pelaporan Pajak:</label>
        <input
          v-model="notes"
          type="text"
          placeholder="Tuliskan keterangan lengkap maksud transaksi pelaporan pajak ini..."
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>
    </div>

    <!-- Balance Status Bar (Seimbang / Double-Entry Settlement Voucher Preview) -->
    <div class="bg-emerald-50/70 border border-emerald-200 p-4 rounded-2xl shadow-2xs text-emerald-950">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-emerald-600 shadow-2xs flex items-center justify-center font-bold text-white shrink-0">
            <CheckCircle2 class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-extrabold uppercase tracking-wider text-[11px]">
                STATUS: JURNAL PENYETORAN SEIMBANG (BALANCED)
              </span>
              <span class="px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-800 font-bold text-[10px]">
                Debet: {{ primaryTaxAccount.code }} == Kredit: {{ activeBankAccount.code }}
              </span>
            </div>
            <p class="text-[11px] mt-0.5 opacity-80">
              Penyetoran SPT Masa tercatat ke rekening kas negara via perbankan resmi (NTPN) dan mengeliminasi hutang pajak klinik.
            </p>
          </div>
        </div>

        <!-- Metric Numbers -->
        <div class="flex items-center gap-4 bg-white/90 px-4 py-2 rounded-xl border border-zinc-200/80 shadow-2xs self-end sm:self-auto font-mono">
          <div>
            <span class="text-[10px] text-zinc-500 font-sans block font-bold">TRANSAKSI TERPILIH</span>
            <span class="font-extrabold text-zinc-950 text-sm">{{ selectedRows.length }} Baris</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] text-zinc-500 font-sans block font-bold">TOTAL DPP</span>
            <span class="font-extrabold text-zinc-950 text-sm">Rp {{ formatIDR(totalDpp) }}</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] text-emerald-700 font-sans block font-bold">TOTAL PAJAK TERUTANG</span>
            <span class="font-extrabold text-emerald-700 text-sm">Rp {{ formatIDR(totalTax) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Voucher Table Form (Tabel Grid Input Baris Pajak - Format Jurnal Umum) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <!-- Toolbar Header -->
      <div class="p-4 border-b border-zinc-200/90 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-zinc-50/50">
        <div class="flex items-center gap-2">
          <BookOpen class="w-4 h-4 text-zinc-700" />
          <span class="text-xs font-bold text-zinc-900">
            Rincian Transaksi Jurnal Pajak ({{ rows.length }} Baris Transaksi)
          </span>
          <span class="text-[11px] text-zinc-500 hidden sm:inline">
            &bull; Ambil otomatis dari Jurnal Umum atau input transaksi baru secara fleksibel
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="toggleSelectAll"
            class="text-[11px] font-bold text-blue-600 hover:text-blue-800 transition px-2 py-1 bg-blue-50 hover:bg-blue-100 rounded-lg cursor-pointer"
          >
            {{ isAllSelected ? 'Batal Pilih Semua' : 'Pilih Semua' }}
          </button>

          <button
            type="button"
            @click="addRow"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Baris Manual</span>
          </button>
        </div>
      </div>

      <!-- Grid Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200/90 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-10 text-center">Pilih</th>
              <th class="px-3 py-3 w-28">Tanggal</th>
              <th class="px-3 py-3 w-36">No. Jurnal Sumber</th>
              <th class="px-4 py-3 min-w-[220px]">Akun Lawan / Beban (COA)</th>
              <th class="px-4 py-3 min-w-[200px]">Akun Pajak Terutang (COA)</th>
              <th class="px-4 py-3 min-w-[200px]">Uraian Transaksi / Rekanan</th>
              <th class="px-3 py-3 w-36 text-right">DPP (Rp)</th>
              <th class="px-3 py-3 w-24 text-center">Tarif (%)</th>
              <th class="px-4 py-3 w-36 text-right">Pajak Terutang (Rp)</th>
              <th class="px-3 py-3 w-12 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr
              v-for="(row, idx) in rows"
              :key="row.id"
              :class="[
                'hover:bg-zinc-50/70 transition',
                row.selected ? 'bg-white' : 'opacity-50 bg-zinc-50/30'
              ]"
            >
              <!-- Checkbox -->
              <td class="px-3 py-3 text-center">
                <input
                  v-model="row.selected"
                  type="checkbox"
                  class="w-4 h-4 rounded text-zinc-950 focus:ring-0 cursor-pointer"
                />
              </td>

              <!-- Date -->
              <td class="px-3 py-2.5">
                <input
                  v-model="row.date"
                  type="date"
                  class="w-full px-2 py-1 bg-zinc-50 border border-zinc-200 rounded-lg text-xs font-mono font-bold text-zinc-900 focus:outline-none"
                />
              </td>

              <!-- Journal Source Voucher Number -->
              <td class="px-3 py-2.5 font-mono">
                <div v-if="row.isFromJournal && row.journalNumber" class="flex items-center gap-1">
                  <button
                    type="button"
                    @click="openJournalPreview(row)"
                    class="px-2 py-0.5 rounded-md bg-blue-50 text-blue-700 font-bold hover:bg-blue-100 transition inline-flex items-center gap-1 cursor-pointer border border-blue-200"
                    title="Klik untuk membuka voucher Jurnal Umum"
                  >
                    <span>{{ row.journalNumber }}</span>
                    <ExternalLink class="w-2.5 h-2.5" />
                  </button>
                </div>
                <div v-else class="text-zinc-400 text-[10px] italic">
                  [Manual Input]
                </div>
              </td>

              <!-- Counterpart Account (Beban/Lawan) -->
              <td class="px-4 py-2.5">
                <SearchableSelect
                  :model-value="row.counterpartCode"
                  :options="accountSelectOptions"
                  placeholder="-- Pilih Akun Beban --"
                  @update:model-value="(code: any) => onCounterpartSelect(row, String(code))"
                />
              </td>

              <!-- Tax Account -->
              <td class="px-4 py-2.5">
                <SearchableSelect
                  :model-value="row.taxAccountCode"
                  :options="accountSelectOptions"
                  placeholder="-- Akun Pajak --"
                  @update:model-value="(code: any) => onTaxAccountSelect(row, String(code))"
                />
              </td>

              <!-- Remark / Description -->
              <td class="px-4 py-2.5">
                <input
                  v-model="row.description"
                  type="text"
                  placeholder="Uraian pemotongan / nama rekanan..."
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>

              <!-- DPP Input -->
              <td class="px-3 py-2.5 text-right">
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                  <input
                    v-model.number="row.dppAmount"
                    @input="() => onRowDppChange(row)"
                    type="number"
                    min="0"
                    step="any"
                    class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                    placeholder="0"
                  />
                </div>
              </td>

              <!-- Tax Rate Input -->
              <td class="px-3 py-2.5 text-center">
                <input
                  v-model="row.taxRate"
                  @change="() => calculateRowTax(row)"
                  type="text"
                  class="w-full px-2 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-center font-bold text-xs text-zinc-800 focus:outline-none"
                />
              </td>

              <!-- Tax Amount Input -->
              <td class="px-4 py-2.5 text-right">
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                  <input
                    v-model.number="row.taxAmount"
                    type="number"
                    min="0"
                    step="any"
                    class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                    placeholder="0"
                  />
                </div>
              </td>

              <!-- Remove Action -->
              <td class="px-3 py-2.5 text-center">
                <button
                  type="button"
                  @click="removeRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition cursor-pointer"
                  title="Hapus Baris"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>

          <!-- Table Footer Total -->
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs text-zinc-900">
            <tr>
              <td colspan="6" class="px-4 py-3 text-right uppercase tracking-wider text-zinc-600">
                Total Terpilih ({{ selectedRows.length }} Baris):
              </td>
              <td class="px-3 py-3 text-right font-mono tabular-nums text-zinc-950 font-extrabold">
                Rp {{ formatIDR(totalDpp) }}
              </td>
              <td class="px-3 py-3 text-center text-zinc-400 font-mono">
                &mdash;
              </td>
              <td class="px-4 py-3 text-right font-mono tabular-nums text-emerald-700 font-extrabold text-sm">
                Rp {{ formatIDR(totalTax) }}
              </td>
              <td class="px-3 py-3"></td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Table Bottom Quick Bar -->
      <div class="p-3 bg-zinc-50/80 border-t border-zinc-200/90 flex flex-wrap items-center justify-between gap-2 text-xs">
        <div class="flex items-center gap-2 text-zinc-500">
          <Sparkles class="w-3.5 h-3.5 text-amber-500" />
          <span>Setiap baris transaksi dengan nomor jurnal terikat langsung ke voucher Buku Besar.</span>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="addRow"
            class="px-3 py-1 bg-white hover:bg-zinc-100 text-zinc-800 rounded-lg border border-zinc-200 font-bold text-[11px] transition cursor-pointer"
          >
            + Baris Baru
          </button>
          <button
            type="button"
            @click="() => pullDataFromGeneralJournal(true)"
            class="px-3 py-1 bg-amber-50 hover:bg-amber-100 text-amber-900 rounded-lg border border-amber-200 font-bold text-[11px] transition cursor-pointer"
          >
            ⚡ Tarik Ulang dari Jurnal
          </button>
        </div>
      </div>
    </div>

    <!-- Settlement Journal Entry Preview (Voucher Jurnal Penyetoran Pajak) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-5 space-y-3">
      <div class="flex items-center justify-between border-b border-zinc-100 pb-2.5">
        <div class="flex items-center gap-2">
          <Landmark class="w-4 h-4 text-emerald-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            Pratinjau Voucher Jurnal Penyetoran Pajak ke Kas Negara (Double-Entry Balanced)
          </h3>
        </div>
        <span class="text-[11px] font-bold text-emerald-700 bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
          Auto-Posting saat status PAID / REPORTED
        </span>
      </div>

      <div class="border border-zinc-200 rounded-xl overflow-hidden text-xs">
        <table class="w-full text-left">
          <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase">
            <tr>
              <th class="px-4 py-2.5">Posisi Akun Akuntansi</th>
              <th class="px-4 py-2.5">Uraian / Memo Baris</th>
              <th class="px-4 py-2.5 text-right w-40">Debet (Rp)</th>
              <th class="px-4 py-2.5 text-right w-40">Kredit (Rp)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr>
              <td class="px-4 py-2.5">
                <span class="font-bold font-mono text-zinc-950">[{{ primaryTaxAccount.code }}]</span>
                <span class="ml-1 font-semibold text-zinc-800">{{ primaryTaxAccount.name }}</span>
              </td>
              <td class="px-4 py-2.5 text-zinc-600 text-[11px]">
                Eliminasi Hutang Pajak SPT Masa {{ taxType }} Masa {{ masaMonth }}/{{ masaYear }}
              </td>
              <td class="px-4 py-2.5 text-right font-mono font-bold text-emerald-700 tabular-nums">
                Rp {{ formatIDR(totalTax) }}
              </td>
              <td class="px-4 py-2.5 text-right font-mono text-zinc-400">
                -
              </td>
            </tr>
            <tr>
              <td class="px-4 py-2.5">
                <span class="font-bold font-mono text-zinc-950">[{{ activeBankAccount.code }}]</span>
                <span class="ml-1 font-semibold text-zinc-800">{{ activeBankAccount.name }}</span>
              </td>
              <td class="px-4 py-2.5 text-zinc-600 text-[11px]">
                Penyetoran Resmi via NTPN: {{ ntpnNumber || '(Nomor NTPN DJP)' }}
              </td>
              <td class="px-4 py-2.5 text-right font-mono text-zinc-400">
                -
              </td>
              <td class="px-4 py-2.5 text-right font-mono font-bold text-emerald-700 tabular-nums">
                Rp {{ formatIDR(totalTax) }}
              </td>
            </tr>
          </tbody>
          <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200 text-xs">
            <tr>
              <td colspan="2" class="px-4 py-2 text-right uppercase tracking-wider text-zinc-600">
                Total Jurnal Penyetoran (Seimbang):
              </td>
              <td class="px-4 py-2 text-right font-mono font-extrabold text-emerald-800">
                Rp {{ formatIDR(totalTax) }}
              </td>
              <td class="px-4 py-2 text-right font-mono font-extrabold text-emerald-800">
                Rp {{ formatIDR(totalTax) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Modal Preview Voucher Jurnal Sumber -->
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
              Voucher Jurnal Umum Sumber: {{ activeJournal?.journalNumber }}
            </h3>
          </div>
          <button @click="showJournalPreviewModal = false" class="text-zinc-400 hover:text-zinc-950 transition cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="text-xs space-y-1 text-zinc-600">
          <p><strong class="text-zinc-900">Tanggal Posting:</strong> {{ formatDate(activeJournal?.journalDate || activeJournal?.transaction_date) }}</p>
          <p><strong class="text-zinc-900">Keterangan:</strong> {{ activeJournal?.description }}</p>
        </div>

        <div class="border border-zinc-200 rounded-xl overflow-hidden text-xs">
          <table class="w-full text-left">
            <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase">
              <tr>
                <th class="px-3 py-2">Akun Buku Besar (COA)</th>
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