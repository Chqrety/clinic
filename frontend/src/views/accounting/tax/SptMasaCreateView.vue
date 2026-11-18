<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTaxStore, type SptMasaType } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import type { BusinessUnit } from '@/types/accounting'
import {
  Folder,
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
  CreditCard,
  Download,
  Share2,
  Tag,
  Paperclip,
  UserPlus,
  MessageSquare,
  Clock,
  Pencil,
  FileSpreadsheet,
  Printer
} from 'lucide-vue-next'
import {
  exportBrayoSptMasaDocumentToExcel,
  exportBrayoSptMasaToCSV,
  getEntityProfile
} from '@/utils/export'

const router = useRouter()
const route = useRoute()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()

// Active Tab in Brayo Pajak layout
const activeTab = ref<'DETAILS' | 'RECONCILIATION' | 'AUDIT'>('DETAILS')

// Header Form State (Metadata Dokumen Pelaporan Pajak Brayo)
const reportNumber = ref<string>(`PJK-${new Date().getFullYear()}-${String(Math.floor(Math.random() * 90000) + 10000)}`)
const reportName = ref<string>('PPH 21 DOKTER & KARYAWAN')
const masaMonth = ref<number>(accountingStore.fiscalMonth || (new Date().getMonth() + 1))
const masaYear = ref<number>(accountingStore.fiscalYear || new Date().getFullYear())
const unit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const currentEntity = computed(() => getEntityProfile(unit.value))
const companyName = ref<string>(currentEntity.value.title)
const createdDate = ref<string>(new Date().toISOString().slice(0, 10))
const sourceTransactionFilter = ref<string>('Semua Transaksi')
const taxType = ref<SptMasaType>('PPH_21')
const selectedTaxAccountCode = ref<string>('2-1500')
const bankAccountCode = ref<string>('1-1113')
const ntpnNumber = ref<string>('')
const notes = ref<string>('')
const status = ref<'DATA_DIAMBIL' | 'DRAFT' | 'DILAPORKAN'>('DATA_DIAMBIL')

// Multi-line Voucher Table Grid State (Brayo Detail Transaksi Pajak)
export interface BrayoTaxDetailRow {
  id: string
  selected: boolean
  isFromJournal: boolean
  journalId?: string
  journalNumber?: string
  date: string
  voucherType: string // e.g. "Journal Entry"
  project: string // e.g. "KIDINEST", "SSC", "KONSOLIDASI"
  sourceTransaction: string // e.g. "Neraca", "Laba Rugi", "Buku Besar"
  counterpartCode: string
  counterpartName: string
  taxAccountCode: string
  taxAccountName: string
  description: string
  sourceStatus: string // e.g. "Active", "Posted"
  dppOriginal: number
  pajakOriginal: number
  dppAdjustment: number
  pajakAdjustment: number
  dppFinal: number
  pajakFinal: number
  isEditing: boolean
  rawJournal?: any
}

const rows = ref<BrayoTaxDetailRow[]>([])
const showJournalPreviewModal = ref<boolean>(false)
const activeJournal = ref<any>(null)
const commentText = ref<string>('')

const taxAccountCatalog = [
  { code: '2-1500', taxType: 'PPH_21' as SptMasaType, label: '2-1500 - HUTANG PPH 21 DOKTER & KARYAWAN', sublabel: 'Withholding Pajak Tenaga Medis Ahli & Karyawan' },
  { code: '2-1510', taxType: 'PPH_23' as SptMasaType, label: '2-1510 - KREDIT/HUTANG PPH 23 JASA & LIMBAH MEDIS', sublabel: 'Pemotongan PPh 23 Jasa Teknik & Limbah B3 2%' },
  { code: '2-1500.42', taxType: 'PPH_4_2' as SptMasaType, label: '2-1500 - HUTANG PPH FINAL 4(2) SEWA FASILITAS', sublabel: 'PPh Final Persewaan Gedung / Bangunan 10%' },
  { code: '2-1520', taxType: 'PPH_25' as SptMasaType, label: '2-1520 - UANG MUKA PPH 25 ANGSURAN BADAN BULANAN', sublabel: 'Setoran Kas Negara Angsuran Pajak Penghasilan Badan' },
  { code: '2-1500.PPN', taxType: 'PPN' as SptMasaType, label: '2-1500 - PPN KELUARAN & MASUKAN 11%', sublabel: 'Pajak Pertambahan Nilai Penjualan & Pembelian' },
  { code: '2501.001', taxType: 'PPH_21' as SptMasaType, label: '2501.001 - HUTANG PPH 21 NOTARIS & TENAGA AHLI', sublabel: 'Format Standar Brayo Pajak PMK 168' },
  { code: '2502', taxType: 'PPH_23' as SptMasaType, label: '2502 - HUTANG PPH 23 JASA & SEWA (BRAYO)', sublabel: 'Format Standar Brayo Pajak PPh 23' },
  { code: '2503', taxType: 'PPH_4_2' as SptMasaType, label: '2503 - HUTANG PPH FINAL 4(2) SEWA (BRAYO)', sublabel: 'Format Standar Brayo Pajak PPh 4(2)' },
  { code: '2505', taxType: 'PPH_25' as SptMasaType, label: '2505 - PPH 25 ANGSURAN BADAN (BRAYO)', sublabel: 'Format Standar Brayo Pajak PPh 25' },
  { code: '2504', taxType: 'PPN' as SptMasaType, label: '2504 - PPN KELUARAN & MASUKAN (BRAYO)', sublabel: 'Format Standar Brayo Pajak PPN 11%' }
]

const activeTaxAccountCatalog = computed(() => {
  return taxAccountCatalog.find(c => c.code === selectedTaxAccountCode.value) || taxAccountCatalog[0]
})

// COA accounts for dropdowns
const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

// Sync taxType when selectedTaxAccountCode changes
watch(selectedTaxAccountCode, (code) => {
  const found = taxAccountCatalog.find(c => c.code === code)
  if (found) {
    taxType.value = found.taxType
    if (found.taxType === 'PPH_21') reportName.value = 'PPH 21 DOKTER & KARYAWAN'
    else if (found.taxType === 'PPH_23') reportName.value = 'PPH 23 JASA & LIMBAH MEDIS'
    else if (found.taxType === 'PPH_4_2') reportName.value = 'PPH FINAL PASAL 4(2) SEWA'
    else if (found.taxType === 'PPH_25') reportName.value = 'PPH 25 ANGSURAN BADAN'
    else if (found.taxType === 'PPN') reportName.value = 'PPN KELUARAN & MASUKAN'
  }
})

// Matching journals in store for current period and tax type
const matchingJournalsInStore = computed(() => {
  const allJournals = taxStore.getLinkedJournalsForMonth(masaMonth.value, masaYear.value, taxType.value)
  if (unit.value === 'ALL') return allJournals
  return allJournals.filter(entry => {
    if (!entry.businessUnit || entry.businessUnit === 'ALL') return true
    return entry.businessUnit === unit.value
  })
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

// Watch filters & selected unit to update company name and re-pull data
watch([masaMonth, masaYear, taxType, selectedTaxAccountCode, unit], () => {
  companyName.value = currentEntity.value.title
  pullDataFromGeneralJournal(false)
})

/**
 * ⚡ Hero Automation: Ambil Data Otomatis dari Jurnal Umum & Buku Besar (Brayo Standard)
 */
function pullDataFromGeneralJournal(showNotice = true) {
  const journals = matchingJournalsInStore.value
  const curCatalog = taxAccountCatalog.find(c => c.code === selectedTaxAccountCode.value) || taxAccountCatalog[0]

  if (journals.length > 0) {
    rows.value = journals.map((j, idx) => {
      let bu = j.businessUnit || (unit.value === 'ALL' ? 'KONSOLIDASI' : unit.value)
      const descUpper = (j.description || '').toUpperCase()
      if (descUpper.includes('KIDINEST') || descUpper.includes('ANAK') || descUpper.includes('SP.A')) bu = 'KIDINEST'
      else if (descUpper.includes('SSC') || descUpper.includes('SKIN') || descUpper.includes('ESTETIKA')) bu = 'SSC'

      return {
        id: `row-brayo-${j.id}-${idx}`,
        selected: true,
        isFromJournal: true,
        journalId: j.journalId,
        journalNumber: j.journalNumber,
        date: j.journalDate ? j.journalDate.slice(0, 10) : new Date().toISOString().slice(0, 10),
        voucherType: 'Journal Entry',
        project: bu,
        sourceTransaction: 'Buku Besar',
        counterpartCode: j.counterpartAccountCode || '6-1100',
        counterpartName: j.counterpartAccountName || 'Beban Operasional Terkait',
        taxAccountCode: j.taxAccountCode || curCatalog.code,
        taxAccountName: j.taxAccountName || curCatalog.label,
        description: j.description || `Pemotongan ${curCatalog.taxType} Masa ${masaMonth.value}/${masaYear.value}`,
        sourceStatus: 'Active',
        dppOriginal: j.dppAmount,
        pajakOriginal: j.taxAmount,
        dppAdjustment: 0,
        pajakAdjustment: 0,
        dppFinal: j.dppAmount,
        pajakFinal: j.taxAmount,
        isEditing: false,
        rawJournal: j.rawJournal
      }
    })

    status.value = 'DATA_DIAMBIL'
    notes.value = `Pelaporan Pajak ${reportName.value} Masa ${masaMonth.value}/${masaYear.value} terintegrasi dari ${rows.value.length} voucher Jurnal Umum (${rows.value.map(r => r.journalNumber).join(', ')})`
    if (showNotice) {
      accountingStore.showToast(`Berhasil menarik ${rows.value.length} transaksi pajak dari Jurnal Umum & Buku Besar!`, 'success')
    }
  } else {
    // If no matching journals, create 1 empty manual row matching Brayo layout
    rows.value = [
      {
        id: `row-manual-1`,
        selected: true,
        isFromJournal: false,
        date: new Date().toISOString().slice(0, 10),
        voucherType: 'Journal Entry',
        project: unit.value === 'ALL' ? 'KONSOLIDASI' : unit.value,
        sourceTransaction: 'Buku Besar',
        counterpartCode: '6-1100',
        counterpartName: 'Beban Operasional Klinik',
        taxAccountCode: curCatalog.code,
        taxAccountName: curCatalog.label,
        description: `Pemotongan Pajak ${reportName.value} Periode Masa ${masaMonth.value}/${masaYear.value}`,
        sourceStatus: 'Draft',
        dppOriginal: 0,
        pajakOriginal: 0,
        dppAdjustment: 0,
        pajakAdjustment: 0,
        dppFinal: 0,
        pajakFinal: 0,
        isEditing: true
      }
    ]
    status.value = 'DRAFT'
    notes.value = `Pelaporan Pajak ${reportName.value} Masa ${masaMonth.value}/${masaYear.value} (Input Manual)`
    if (showNotice) {
      accountingStore.showToast(`Belum ada jurnal untuk akun ${curCatalog.code} pada Masa ${masaMonth.value}/${masaYear.value}. Formulir siap diinput manual.`, 'info')
    }
  }
}

// Add manual row like Brayo
function addRow() {
  const curCatalog = taxAccountCatalog.find(c => c.code === selectedTaxAccountCode.value) || taxAccountCatalog[0]
  rows.value.push({
    id: `row-manual-${Date.now()}`,
    selected: true,
    isFromJournal: false,
    date: new Date().toISOString().slice(0, 10),
    voucherType: 'Journal Entry',
    project: unit.value === 'ALL' ? 'KONSOLIDASI' : unit.value,
    sourceTransaction: 'Buku Besar',
    counterpartCode: '6-1200',
    counterpartName: 'Beban Operasional Klinik',
    taxAccountCode: curCatalog.code,
    taxAccountName: curCatalog.label,
    description: `Transaksi Pajak Tambahan ${reportName.value}`,
    sourceStatus: 'Draft',
    dppOriginal: 0,
    pajakOriginal: 0,
    dppAdjustment: 0,
    pajakAdjustment: 0,
    dppFinal: 0,
    pajakFinal: 0,
    isEditing: true
  })
}

function removeRow(index: number) {
  if (rows.value.length > 1) {
    rows.value.splice(index, 1)
  } else {
    rows.value[0].dppFinal = 0
    rows.value[0].pajakFinal = 0
    rows.value[0].description = ''
  }
}

function onDppFinalChange(row: BrayoTaxDetailRow) {
  row.dppAdjustment = Number(row.dppFinal || 0) - Number(row.dppOriginal || 0)
  // Recalculate suggested tax based on taxType
  const dpp = Number(row.dppFinal || 0)
  if (taxType.value === 'PPH_21') {
    row.pajakFinal = Math.round(dpp * 0.05)
  } else if (taxType.value === 'PPH_23') {
    row.pajakFinal = Math.round(dpp * 0.02)
  } else if (taxType.value === 'PPH_4_2') {
    row.pajakFinal = Math.round(dpp * 0.10)
  } else if (taxType.value === 'PPN') {
    row.pajakFinal = Math.round(dpp * 0.11)
  }
  row.pajakAdjustment = Number(row.pajakFinal || 0) - Number(row.pajakOriginal || 0)
}

function onPajakFinalChange(row: BrayoTaxDetailRow) {
  row.pajakAdjustment = Number(row.pajakFinal || 0) - Number(row.pajakOriginal || 0)
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

// Brayo 3-Column Summary Totals (Original vs Adjustment vs Final)
const totalDppOriginal = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.dppOriginal) || 0), 0)
})

const totalPajakOriginal = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.pajakOriginal) || 0), 0)
})

const totalDppAdjustment = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.dppAdjustment) || 0), 0)
})

const totalPajakAdjustment = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.pajakAdjustment) || 0), 0)
})

const totalDppFinal = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.dppFinal) || 0), 0)
})

const totalPajakFinal = computed(() => {
  return selectedRows.value.reduce((sum, r) => sum + (Number(r.pajakFinal) || 0), 0)
})

// Bank account info for settlement journal preview
const activeBankAccount = computed(() => {
  return taxStore.availableBankAccounts.find(b => b.code === bankAccountCode.value) || {
    code: '1-1113',
    name: 'Bank BCA Operasional'
  }
})

// Journal Preview modal
function openJournalPreview(row: BrayoTaxDetailRow) {
  if (row.rawJournal) {
    activeJournal.value = row.rawJournal
  } else {
    const found = journalsStore.journals.find(j => j.journalNumber === row.journalNumber || j.id === row.journalId)
    activeJournal.value = found || {
      journalNumber: row.journalNumber || 'JV-PREVIEW',
      journalDate: row.date,
      description: row.description,
      lines: [
        { accountCode: row.counterpartCode, accountName: row.counterpartName, debit: row.pajakFinal, credit: 0 },
        { accountCode: row.taxAccountCode, accountName: row.taxAccountName, debit: 0, credit: row.pajakFinal }
      ]
    }
  }
  showJournalPreviewModal.value = true
}

// Save action
function handleSave() {
  if (totalPajakFinal.value <= 0) {
    accountingStore.showToast('Total Pajak Final harus lebih dari Rp 0 untuk dapat menyimpan pelaporan SPT Masa', 'error')
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
      grossAmount: r.dppFinal,
      dppAmount: r.dppFinal,
      taxAmount: r.pajakFinal,
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
    grossAmount: totalDppFinal.value,
    dppAmount: totalDppFinal.value,
    taxRate: 'DJP / Withholding',
    taxAmount: totalPajakFinal.value,
    notes: `[Brayo ${reportNumber.value}] ${notes.value || ''}`.trim(),
    ntpnNumber: ntpnNumber.value.trim(),
    status: status.value === 'DILAPORKAN' ? 'REPORTED' : (status.value === 'DATA_DIAMBIL' ? 'PAID' : 'DRAFT'),
    businessUnit: unit.value,
    linkedJournals: linkedJournalsList
  })

  // Post payment journal if status is final
  if (status.value === 'DILAPORKAN' || status.value === 'DATA_DIAMBIL') {
    journalsStore.createManualJournal({
      journalDate: `${masaYear.value}-${String(masaMonth.value).padStart(2, '0')}-20`,
      journalSource: 'MANUAL',
      description: `Penyetoran Pajak ${reportName.value} (${reportNumber.value}) via ${activeBankAccount.value.name} (NTPN: ${ntpnNumber.value || 'DJP-SETOR'})`,
      businessUnit: unit.value,
      lines: [
        {
          accountCode: selectedTaxAccountCode.value.startsWith('2-15') ? selectedTaxAccountCode.value : '2-1500',
          accountName: `Hutang Pajak ${reportName.value}`,
          description: `Pelunasan Hutang Pajak ${reportName.value}`,
          debit: totalPajakFinal.value,
          credit: 0
        },
        {
          accountCode: activeBankAccount.value.code,
          accountName: activeBankAccount.value.name,
          description: `Penyetoran Kas Negara via Bank (NTPN: ${ntpnNumber.value || 'DJP-SETOR'})`,
          debit: 0,
          credit: totalPajakFinal.value
        }
      ]
    })
  }

  accountingStore.showToast(`Dokumen SPT Masa ${reportNumber.value} berhasil disimpan dan disinkronkan ke Buku Besar!`, 'success')
  router.push('/accounting/tax/spt-masa')
}

// 📥 Export Per Jurnal ke Excel (.xlsx) sesuai format Brayo & Dokumen Laporan_Pajak_PJK
function handleExportExcel() {
  const curCatalog = taxAccountCatalog.find(c => c.code === selectedTaxAccountCode.value) || taxAccountCatalog[0]
  exportBrayoSptMasaDocumentToExcel({
    reportNumber: reportNumber.value,
    reportName: reportName.value,
    masaMonth: Number(masaMonth.value),
    masaYear: Number(masaYear.value),
    companyName: companyName.value,
    createdDate: createdDate.value,
    project: unit.value === 'ALL' ? 'KONSOLIDASI' : unit.value,
    sourceTransaction: sourceTransactionFilter.value,
    taxAccountCode: selectedTaxAccountCode.value,
    taxAccountName: curCatalog.label,
    bankAccount: `[${activeBankAccount.value.code}] ${activeBankAccount.value.name}`,
    ntpnNumber: ntpnNumber.value,
    status: status.value === 'DATA_DIAMBIL' ? 'Data Diambil' : (status.value === 'DILAPORKAN' ? 'Dilaporkan' : 'Draft'),
    notes: notes.value,
    rows: rows.value,
    totalDppOriginal: totalDppOriginal.value,
    totalPajakOriginal: totalPajakOriginal.value,
    totalDppAdjustment: totalDppAdjustment.value,
    totalPajakAdjustment: totalPajakAdjustment.value,
    totalDppFinal: totalDppFinal.value,
    totalPajakFinal: totalPajakFinal.value
  }, `Laporan_Pajak_${reportNumber.value}`)
  accountingStore.showToast(`Dokumen ${reportNumber.value} berhasil diekspor ke Excel per transaksi jurnal!`, 'success')
}

// 📥 Export Per Jurnal ke CSV
function handleExportCSV() {
  const curCatalog = taxAccountCatalog.find(c => c.code === selectedTaxAccountCode.value) || taxAccountCatalog[0]
  exportBrayoSptMasaToCSV({
    reportNumber: reportNumber.value,
    reportName: reportName.value,
    masaMonth: Number(masaMonth.value),
    masaYear: Number(masaYear.value),
    companyName: companyName.value,
    createdDate: createdDate.value,
    project: unit.value === 'ALL' ? 'KONSOLIDASI' : unit.value,
    sourceTransaction: sourceTransactionFilter.value,
    taxAccountCode: selectedTaxAccountCode.value,
    taxAccountName: curCatalog.label,
    bankAccount: `[${activeBankAccount.value.code}] ${activeBankAccount.value.name}`,
    ntpnNumber: ntpnNumber.value,
    status: status.value === 'DATA_DIAMBIL' ? 'Data Diambil' : (status.value === 'DILAPORKAN' ? 'Dilaporkan' : 'Draft'),
    notes: notes.value,
    rows: rows.value,
    totalDppOriginal: totalDppOriginal.value,
    totalPajakOriginal: totalPajakOriginal.value,
    totalDppAdjustment: totalDppAdjustment.value,
    totalPajakAdjustment: totalPajakAdjustment.value,
    totalDppFinal: totalDppFinal.value,
    totalPajakFinal: totalPajakFinal.value
  }, `Laporan_Pajak_${reportNumber.value}`)
  accountingStore.showToast(`Data transaksi pajak ${reportNumber.value} berhasil diekspor ke CSV!`, 'success')
}

// 🖨️ Cetak Resmi Formulir SPT Masa (DJP / Brayo Layout)
function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="font-['Plus_Jakarta_Sans',sans-serif] max-w-7xl mx-auto pb-16">
    <!-- ========================================================================= -->
    <!-- SCREEN UI: BRAYO PAJAK MODERN DARK INTERFACE (apptax.brayo.id)             -->
    <!-- Meniru 1:1 Tampilan Sederhana & Bersih Sesuai Screenshot User              -->
    <!-- ========================================================================= -->
    <div class="print:hidden bg-[#18191c] text-zinc-300 min-h-screen -mx-4 -mt-4 p-5 sm:p-7 rounded-2xl shadow-xl space-y-6">
      
      <!-- 1. Header: Breadcrumb & Brayo Action Bar -->
      <header class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-3 border-b border-[#282a31]">
        <!-- Breadcrumb -->
        <div class="flex items-center gap-2 text-xs text-zinc-400 flex-wrap">
          <router-link to="/accounting/tax/spt-masa" class="hover:text-zinc-200 transition flex items-center gap-1.5 text-zinc-400">
            <Folder class="w-4 h-4 text-zinc-500" />
            <span>Brayo Pajak</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-400">Pelaporan Pajak</span>
          <span>/</span>
          <span class="font-mono font-bold text-zinc-100">{{ reportNumber }}</span>
          <span class="ml-1.5 px-2 py-0.5 rounded text-[11px] font-bold bg-[#143c44] text-[#2dd4bf] border border-[#195663]">
            {{ status === 'DATA_DIAMBIL' ? 'Data Diambil' : (status === 'DILAPORKAN' ? 'Dilaporkan' : 'Draft') }}
          </span>
        </div>

        <!-- Top Right Actions matching Brayo Screenshot -->
        <div class="flex items-center gap-2 flex-wrap">
          <button
            type="button"
            @click="handleExportExcel"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-zinc-200 text-zinc-950 font-bold rounded-lg text-xs shadow-sm transition active:scale-95 cursor-pointer"
            title="Export Excel (.xlsx)"
          >
            <FileSpreadsheet class="w-3.5 h-3.5 text-emerald-700" />
            <span>Export Excel (.xlsx)</span>
          </button>

          <button
            type="button"
            @click="handleExportCSV"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#25272e] hover:bg-[#2d3038] text-zinc-300 font-semibold rounded-lg text-xs border border-[#373a43] transition active:scale-95 cursor-pointer"
            title="Export CSV"
          >
            <Download class="w-3.5 h-3.5 text-zinc-400" />
            <span>Export CSV</span>
          </button>

          <button
            type="button"
            @click="handlePrint"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-[#25272e] hover:bg-[#2d3038] text-zinc-300 font-semibold rounded-lg text-xs border border-[#373a43] transition active:scale-95 cursor-pointer"
            title="Cetak Formulir Resmi DJP"
          >
            <Printer class="w-3.5 h-3.5 text-blue-400" />
            <span>Cetak Formulir</span>
          </button>

          <button
            type="button"
            @click="() => pullDataFromGeneralJournal(true)"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-500/20 hover:bg-amber-500/30 text-amber-300 font-bold rounded-lg text-xs border border-amber-500/40 transition active:scale-95 cursor-pointer"
            title="Tarik data transaksi terbaru dari Jurnal Umum & Buku Besar"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-400" />
            <span>⚡ Ambil Data Otomatis</span>
          </button>

          <button
            type="button"
            @click="handleSave"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white font-bold rounded-lg text-xs shadow-sm transition active:scale-95 cursor-pointer"
          >
            <Save class="w-3.5 h-3.5" />
            <span>Simpan</span>
          </button>

          <router-link
            to="/accounting/tax/spt-masa"
            class="inline-flex items-center px-3 py-1.5 text-zinc-400 hover:text-zinc-200 text-xs font-semibold transition"
          >
            Cancel
          </router-link>
        </div>
      </header>

      <!-- 2. Brayo Tabs Navigation -->
      <div class="border-b border-[#282a31]">
        <nav class="flex space-x-6 text-xs font-medium">
          <button
            type="button"
            @click="activeTab = 'DETAILS'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'DETAILS' ? 'border-b-2 border-white text-white font-bold' : 'border-b-2 border-transparent text-zinc-400 hover:text-zinc-200'"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Details</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'RECONCILIATION'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'RECONCILIATION' ? 'border-b-2 border-white text-white font-bold' : 'border-b-2 border-transparent text-zinc-400 hover:text-zinc-200'"
          >
            <Calculator class="w-3.5 h-3.5" />
            <span>Rekonsiliasi & Status</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'AUDIT'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'AUDIT' ? 'border-b-2 border-white text-white font-bold' : 'border-b-2 border-transparent text-zinc-400 hover:text-zinc-200'"
          >
            <Clock class="w-3.5 h-3.5" />
            <span>Audit Trail & Sync</span>
          </button>
        </nav>
      </div>

      <!-- 3. TAB 1: DETAILS (Brayo 2-Column Split: Left Content + Right Sidebar) -->
      <div v-if="activeTab === 'DETAILS'" class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        
        <!-- LEFT MAIN CONTENT (9 Cols) -->
        <div class="lg:col-span-9 space-y-7">
          
          <!-- Section A: Periode Pelaporan (Form 2 Kolom Bersih Sesuai Screenshot) -->
          <div>
            <h3 class="text-sm font-semibold text-zinc-200 mb-3.5">
              Periode Pelaporan
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-3.5 text-xs">
              <!-- Left Form Col -->
              <div class="space-y-3">
                <div>
                  <label class="block text-zinc-400 mb-1">Nama Laporan / Transaksi *</label>
                  <input
                    v-model="reportName"
                    type="text"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none focus:border-zinc-500 transition"
                  />
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Company *</label>
                  <input
                    v-model="companyName"
                    type="text"
                    readonly
                    class="w-full px-3 py-2 bg-[#1b1c20] border border-[#272930] rounded-lg text-zinc-400 font-semibold cursor-not-allowed"
                  />
                  <p class="text-[10px] text-zinc-500 mt-1 font-mono">NPWP: {{ currentEntity.npwp }} • {{ currentEntity.subtitle }}</p>
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Project</label>
                  <select
                    v-model="unit"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none focus:border-zinc-500 cursor-pointer transition"
                  >
                    <option value="ALL">Semua Unit (Konsolidasi)</option>
                    <option value="KIDINEST">Unit Kidinest (Klinik Anak)</option>
                    <option value="SSC">Unit SSC (Semarang Skin Center)</option>
                  </select>
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Sumber Transaksi *</label>
                  <select
                    v-model="sourceTransactionFilter"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none focus:border-zinc-500 cursor-pointer transition"
                  >
                    <option value="Semua Transaksi">Semua Transaksi</option>
                    <option value="Buku Besar">Buku Besar (General Ledger)</option>
                    <option value="Jurnal Umum">Jurnal Umum (Journal Entry)</option>
                    <option value="Kasir Closing">Transaksi Kasir Otomatis</option>
                  </select>
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Jenis Pajak</label>
                  <select
                    v-model="selectedTaxAccountCode"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none focus:border-zinc-500 cursor-pointer transition"
                  >
                    <option
                      v-for="c in taxAccountCatalog"
                      :key="c.code"
                      :value="c.code"
                    >
                      {{ c.label }}
                    </option>
                  </select>
                </div>

                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-zinc-400 mb-1">Bulan *</label>
                    <select
                      v-model.number="masaMonth"
                      class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none cursor-pointer"
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
                    <label class="block text-zinc-400 mb-1">Tahun *</label>
                    <input
                      v-model.number="masaYear"
                      type="number"
                      class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg font-mono text-zinc-100 font-semibold focus:outline-none"
                    />
                  </div>
                </div>
              </div>

              <!-- Right Form Col -->
              <div class="space-y-3">
                <div>
                  <label class="block text-zinc-400 mb-1">Status</label>
                  <input
                    type="text"
                    readonly
                    :value="status === 'DATA_DIAMBIL' ? 'Data Diambil' : (status === 'DILAPORKAN' ? 'Dilaporkan' : 'Draft')"
                    class="w-full px-3 py-2 bg-[#1b1c20] border border-[#272930] rounded-lg text-teal-400 font-semibold cursor-not-allowed"
                  />
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Tanggal Dibuat</label>
                  <input
                    v-model="createdDate"
                    type="date"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg font-mono text-zinc-100 focus:outline-none"
                  />
                </div>

                <div>
                  <label class="block text-zinc-400 mb-1">Rekening Bank Penyetor</label>
                  <select
                    v-model="bankAccountCode"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-zinc-100 font-semibold focus:outline-none cursor-pointer"
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
                  <label class="block text-zinc-400 mb-1">Nomor NTPN / Kode Billing DJP</label>
                  <input
                    v-model="ntpnNumber"
                    type="text"
                    placeholder="Contoh: NTPN-BCA-20261120-9901"
                    class="w-full px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg font-mono text-zinc-100 placeholder-zinc-600 focus:outline-none"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Section B: Detail Transaksi Pajak (Brayo Clean Dark Table) -->
          <div>
            <div class="flex items-center justify-between mb-2">
              <div>
                <h3 class="text-sm font-semibold text-zinc-200">
                  Detail Transaksi Pajak
                </h3>
                <span class="text-xs text-zinc-500">Detail</span>
              </div>

              <div class="flex items-center gap-2">
                <button
                  type="button"
                  @click="toggleSelectAll"
                  class="text-[11px] font-medium text-cyan-400 hover:text-cyan-300 transition px-2 py-1 rounded cursor-pointer"
                >
                  {{ isAllSelected ? 'Batal Pilih Semua' : 'Pilih Semua' }}
                </button>
                <button
                  type="button"
                  @click="addRow"
                  class="inline-flex items-center gap-1 px-2.5 py-1 bg-[#24262c] hover:bg-[#2e3138] text-zinc-200 rounded-lg text-xs font-semibold border border-[#373a43] transition cursor-pointer"
                >
                  <Plus class="w-3 h-3" />
                  <span>Tambah Baris</span>
                </button>
              </div>
            </div>

            <!-- Brayo Clean Dark Table -->
            <div class="border border-[#282a31] rounded-xl overflow-hidden bg-[#18191d]">
              <div class="overflow-x-auto">
                <table class="w-full text-xs text-left">
                  <thead class="bg-[#1e2025] text-zinc-400 font-semibold border-b border-[#282a31] text-[11px]">
                    <tr>
                      <th class="px-3 py-2.5 w-10 text-center">
                        <input
                          type="checkbox"
                          :checked="isAllSelected"
                          @change="toggleSelectAll"
                          class="w-3.5 h-3.5 rounded bg-[#25272e] border-zinc-600 text-cyan-500 focus:ring-0 cursor-pointer"
                        />
                      </th>
                      <th class="px-2 py-2.5 w-8 text-center text-zinc-500">No.</th>
                      <th class="px-3 py-2.5 w-24">Tanggal</th>
                      <th class="px-3 py-2.5 w-28">Voucher Type</th>
                      <th class="px-3 py-2.5 w-36">No Bukti</th>
                      <th class="px-3 py-2.5 w-28">Project</th>
                      <th class="px-3 py-2.5 w-28">Sumber Transaksi</th>
                      <th class="px-3 py-2.5 min-w-[180px]">Akun Pajak</th>
                      <th class="px-3 py-2.5 w-20 text-center">Source Status</th>
                      <th class="px-3 py-2.5 w-32 text-right">DPP Final</th>
                      <th class="px-3 py-2.5 w-32 text-right">Pajak Final</th>
                      <th class="px-2 py-2.5 w-8 text-center"></th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-[#23252c] text-zinc-300">
                    <tr
                      v-for="(row, idx) in rows"
                      :key="row.id"
                      :class="[
                        'hover:bg-[#1f2127] transition',
                        row.selected ? 'bg-transparent' : 'opacity-40'
                      ]"
                    >
                      <!-- Checkbox -->
                      <td class="px-3 py-2 text-center">
                        <input
                          v-model="row.selected"
                          type="checkbox"
                          class="w-3.5 h-3.5 rounded bg-[#25272e] border-zinc-600 text-cyan-500 focus:ring-0 cursor-pointer"
                        />
                      </td>

                      <!-- No -->
                      <td class="px-2 py-2 text-center text-zinc-500 font-mono">
                        {{ idx + 1 }}
                      </td>

                      <!-- Tanggal -->
                      <td class="px-3 py-2 font-mono text-zinc-300">
                        <input
                          v-model="row.date"
                          type="date"
                          class="bg-transparent border-b border-transparent hover:border-zinc-600 text-xs text-zinc-200 focus:outline-none"
                        />
                      </td>

                      <!-- Voucher Type -->
                      <td class="px-3 py-2 text-zinc-400">
                        {{ row.voucherType }}
                      </td>

                      <!-- No Bukti -->
                      <td class="px-3 py-2 font-mono">
                        <button
                          v-if="row.journalNumber"
                          type="button"
                          @click="openJournalPreview(row)"
                          class="text-cyan-400 hover:text-cyan-300 hover:underline flex items-center gap-1 cursor-pointer font-bold"
                          title="Klik pratinjau voucher jurnal"
                        >
                          <span>{{ row.journalNumber }}</span>
                          <ExternalLink class="w-2.5 h-2.5 text-zinc-500" />
                        </button>
                        <span v-else class="text-zinc-500 italic">[Manual]</span>
                      </td>

                      <!-- Project -->
                      <td class="px-3 py-2">
                        <span class="px-1.5 py-0.5 rounded text-[10px] bg-[#24262c] text-zinc-300 font-medium">
                          {{ row.project }}
                        </span>
                      </td>

                      <!-- Sumber Transaksi -->
                      <td class="px-3 py-2 text-zinc-400">
                        {{ row.sourceTransaction }}
                      </td>

                      <!-- Akun Pajak -->
                      <td class="px-3 py-2 font-mono">
                        <span class="text-zinc-200 font-semibold">{{ row.taxAccountCode }}</span>
                        <span class="text-zinc-400 text-[11px] block truncate max-w-[180px]">{{ row.taxAccountName }}</span>
                      </td>

                      <!-- Source Status -->
                      <td class="px-3 py-2 text-center">
                        <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-[#163628] text-[#4ade80]">
                          {{ row.sourceStatus }}
                        </span>
                      </td>

                      <!-- DPP Final -->
                      <td class="px-3 py-2 text-right font-mono">
                        <input
                          v-model.number="row.dppFinal"
                          @input="() => onDppFinalChange(row)"
                          type="number"
                          class="w-28 px-2 py-1 bg-[#202227] border border-[#2f323a] rounded text-right font-mono text-zinc-100 focus:outline-none focus:border-zinc-500"
                        />
                      </td>

                      <!-- Pajak Final with Pen Icon -->
                      <td class="px-3 py-2 text-right font-mono">
                        <div class="flex items-center justify-end gap-1">
                          <input
                            v-model.number="row.pajakFinal"
                            @input="() => onPajakFinalChange(row)"
                            type="number"
                            class="w-28 px-2 py-1 bg-[#202227] border border-[#2f323a] rounded text-right font-mono font-bold text-teal-400 focus:outline-none focus:border-zinc-500"
                          />
                          <Pencil class="w-3 h-3 text-zinc-500 shrink-0" />
                        </div>
                      </td>

                      <!-- Delete -->
                      <td class="px-2 py-2 text-center">
                        <button
                          type="button"
                          @click="removeRow(idx)"
                          class="text-zinc-600 hover:text-rose-400 transition cursor-pointer"
                        >
                          <Trash2 class="w-3 h-3" />
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- Section C: Ringkasan Nilai Pajak (3 Kolom Bersih Sesuai Screenshot) -->
          <div>
            <h3 class="text-sm font-semibold text-zinc-200 mb-3">
              Ringkasan Nilai Pajak
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 text-xs bg-[#1a1b1f] border border-[#282a31] rounded-xl p-4">
              <!-- Col 1: Original -->
              <div class="space-y-3">
                <div>
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total DPP Original</span>
                  <span class="font-mono text-sm text-zinc-100 font-bold">
                    Rp {{ formatIDR(totalDppOriginal) }}
                  </span>
                </div>
                <div class="pt-2 border-t border-[#26282e]">
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total Pajak Original</span>
                  <span class="font-mono text-sm text-zinc-100 font-bold">
                    Rp {{ formatIDR(totalPajakOriginal) }}
                  </span>
                </div>
              </div>

              <!-- Col 2: Adjustment -->
              <div class="space-y-3">
                <div>
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total DPP Adjustment</span>
                  <span class="font-mono text-sm text-zinc-100 font-bold">
                    Rp {{ formatIDR(totalDppAdjustment) }}
                  </span>
                </div>
                <div class="pt-2 border-t border-[#26282e]">
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total Pajak Adjustment</span>
                  <span class="font-mono text-sm text-zinc-100 font-bold">
                    Rp {{ formatIDR(totalPajakAdjustment) }}
                  </span>
                </div>
              </div>

              <!-- Col 3: Final -->
              <div class="space-y-3">
                <div>
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total DPP Final</span>
                  <span class="font-mono text-sm text-teal-400 font-bold">
                    Rp {{ formatIDR(totalDppFinal) }}
                  </span>
                </div>
                <div class="pt-2 border-t border-[#26282e]">
                  <span class="text-zinc-400 block text-[11px] mb-0.5">Total Pajak Final</span>
                  <span class="font-mono text-sm text-teal-400 font-extrabold">
                    Rp {{ formatIDR(totalPajakFinal) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Section D: Comments (Sesuai Screenshot Brayo) -->
          <div class="pt-2">
            <h3 class="text-sm font-semibold text-zinc-200 mb-3">
              Comments
            </h3>

            <div class="flex items-start gap-3">
              <div class="w-8 h-8 rounded-full bg-pink-700 text-white font-bold text-xs flex items-center justify-center shrink-0">
                AP
              </div>
              <div class="flex-1 flex items-center gap-2">
                <input
                  v-model="commentText"
                  type="text"
                  placeholder="Type a reply / comment..."
                  class="flex-1 px-3 py-2 bg-[#202227] border border-[#2f323a] rounded-lg text-xs text-zinc-100 placeholder-zinc-500 focus:outline-none focus:border-zinc-500"
                  @keyup.enter="accountingStore.showToast('Komentar berhasil ditambahkan!', 'success')"
                />
                <button
                  type="button"
                  @click="accountingStore.showToast('Komentar berhasil ditambahkan!', 'success')"
                  class="px-3 py-2 bg-[#2a2c33] hover:bg-[#343740] text-zinc-200 rounded-lg text-xs font-semibold transition cursor-pointer"
                >
                  Send
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- RIGHT SIDEBAR (3 Cols - Matching Brayo Screenshot) -->
        <div class="lg:col-span-3 space-y-5 lg:border-l lg:border-[#282a31] lg:pl-6">
          <!-- Document Header -->
          <div class="space-y-1">
            <h4 class="text-sm font-bold text-zinc-100 uppercase">
              {{ reportName }}
            </h4>
            <div class="flex items-center gap-2">
              <span class="font-mono text-xs text-zinc-400">{{ reportNumber }}</span>
              <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-[#143c44] text-[#2dd4bf]">
                {{ status === 'DATA_DIAMBIL' ? 'Data Diambil' : status }}
              </span>
            </div>
          </div>

          <!-- Brayo Action Items with '+' icon -->
          <div class="space-y-3.5 text-xs text-zinc-300">
            <!-- Assign -->
            <div class="flex items-center justify-between py-1 border-b border-[#24262c]">
              <span class="flex items-center gap-2 text-zinc-400">
                <UserPlus class="w-3.5 h-3.5 text-zinc-500" />
                <span>Assign</span>
              </span>
              <button type="button" class="text-zinc-500 hover:text-zinc-300 font-bold cursor-pointer">+</button>
            </div>

            <!-- Attachments -->
            <div class="space-y-2 py-1 border-b border-[#24262c]">
              <div class="flex items-center justify-between">
                <span class="flex items-center gap-2 text-zinc-400">
                  <Paperclip class="w-3.5 h-3.5 text-zinc-500" />
                  <span>Attachments</span>
                </span>
                <button type="button" @click="handleExportExcel" class="text-zinc-500 hover:text-zinc-300 font-bold cursor-pointer">+</button>
              </div>

              <!-- Attached file chip (matching screenshot) -->
              <div
                @click="handleExportExcel"
                class="flex items-center gap-2.5 p-2 bg-[#202227] hover:bg-[#252830] border border-[#2c2f38] rounded-lg cursor-pointer transition"
                title="Klik untuk mengunduh Excel"
              >
                <FileSpreadsheet class="w-4 h-4 text-emerald-400 shrink-0" />
                <div class="min-w-0 flex-1">
                  <div class="text-[11px] font-semibold text-zinc-200 truncate">
                    Laporan_Pajak_{{ reportNumber }}.xlsx
                  </div>
                  <div class="text-[10px] text-zinc-500 font-mono">
                    XLSX • 6.4 KB
                  </div>
                </div>
              </div>
            </div>

            <!-- Tags -->
            <div class="space-y-2 py-1 border-b border-[#24262c]">
              <div class="flex items-center justify-between">
                <span class="flex items-center gap-2 text-zinc-400">
                  <Tag class="w-3.5 h-3.5 text-zinc-500" />
                  <span>Tags</span>
                </span>
                <button type="button" class="text-zinc-500 hover:text-zinc-300 font-bold cursor-pointer">+</button>
              </div>
              <div class="flex flex-wrap gap-1.5">
                <span class="px-2 py-0.5 rounded text-[10px] bg-[#24262c] text-zinc-400 border border-[#2e3038]">
                  SPT Masa
                </span>
                <span class="px-2 py-0.5 rounded text-[10px] bg-[#24262c] text-zinc-400 border border-[#2e3038]">
                  {{ unit }}
                </span>
                <span class="px-2 py-0.5 rounded text-[10px] bg-[#24262c] text-zinc-400 border border-[#2e3038]">
                  PMK DJP
                </span>
              </div>
            </div>

            <!-- Share -->
            <div class="flex items-center justify-between py-1 border-b border-[#24262c]">
              <span class="flex items-center gap-2 text-zinc-400">
                <Share2 class="w-3.5 h-3.5 text-zinc-500" />
                <span>Share</span>
              </span>
              <button type="button" class="text-zinc-500 hover:text-zinc-300 font-bold cursor-pointer">+</button>
            </div>
          </div>

          <!-- Metadata timestamps (matching Brayo screenshot) -->
          <div class="space-y-1 text-[11px] text-zinc-500 pt-2">
            <div>You last edited this · Just now</div>
            <div>You created this · Today</div>
          </div>
        </div>
      </div>

      <!-- TAB 2: RECONCILIATION & STATUS -->
      <div v-else-if="activeTab === 'RECONCILIATION'" class="space-y-6">
        <div class="bg-[#1e2025] border border-[#282a31] rounded-xl p-5 space-y-4">
          <div class="border-b border-[#282a31] pb-3 flex items-center justify-between">
            <div>
              <h3 class="text-sm font-bold text-zinc-100">
                Rekonsiliasi Saldo Akun Pajak di Buku Besar vs SPT Masa
              </h3>
              <p class="text-xs text-zinc-400">
                Memverifikasi bahwa seluruh pemotongan pada Jurnal Umum telah tercakup 100% pada pelaporan SPT Masa tanpa selisih.
              </p>
            </div>
            <span class="px-2.5 py-1 rounded-full text-xs font-bold bg-[#143c44] text-[#2dd4bf] border border-[#195663]">
              Status: Tereconciled
            </span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-xs font-mono">
            <div class="p-4 rounded-xl bg-[#18191d] border border-[#282a31]">
              <span class="text-zinc-400 font-sans block font-semibold text-[11px]">SALDO AKUN BUKU BESAR ({{ selectedTaxAccountCode }})</span>
              <span class="text-lg font-bold text-zinc-100 block mt-1">Rp {{ formatIDR(totalPajakOriginal) }}</span>
            </div>
            <div class="p-4 rounded-xl bg-[#18191d] border border-[#282a31]">
              <span class="text-zinc-400 font-sans block font-semibold text-[11px]">TOTAL PELAPORAN SPT MASA FINAL</span>
              <span class="text-lg font-bold text-teal-400 block mt-1">Rp {{ formatIDR(totalPajakFinal) }}</span>
            </div>
            <div class="p-4 rounded-xl bg-[#163628] border border-[#1e5239]">
              <span class="text-emerald-300 font-sans block font-semibold text-[11px]">VARIANCE / SELISIH REKONSILIASI</span>
              <span class="text-lg font-black text-emerald-400 block mt-1">Rp 0,00 (SEIMBANG)</span>
            </div>
          </div>
        </div>
      </div>

      <!-- TAB 3: AUDIT TRAIL & SYNC -->
      <div v-else class="space-y-6">
        <div class="bg-[#1e2025] border border-[#282a31] rounded-xl p-5 space-y-4">
          <div class="border-b border-[#282a31] pb-3">
            <h3 class="text-sm font-bold text-zinc-100">
              Audit Trail & Log Sinkronisasi Jurnal
            </h3>
            <p class="text-xs text-zinc-400">
              Catatan riwayat penarikan data dari Jurnal Umum & Buku Besar klinik.
            </p>
          </div>

          <div class="space-y-3 text-xs">
            <div class="flex items-start gap-3 p-3 bg-[#18191d] rounded-xl border border-[#282a31]">
              <CheckCircle2 class="w-4 h-4 text-emerald-400 shrink-0 mt-0.5" />
              <div>
                <p class="font-bold text-zinc-100">Data Diambil dari Jurnal Umum</p>
                <p class="text-zinc-400 text-[11px]">
                  Menarik {{ rows.length }} voucher transaksi jurnal untuk akun {{ selectedTaxAccountCode }} Masa {{ masaMonth }}/{{ masaYear }}.
                </p>
                <span class="text-[10px] text-zinc-500 font-mono">{{ createdDate }} • Operator: {{ accountingStore.currentUserName }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Preview Voucher Jurnal Sumber -->
      <div
        v-if="showJournalPreviewModal"
        class="fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4"
        @click.self="showJournalPreviewModal = false"
      >
        <div class="bg-[#1e2025] border border-[#2f323a] rounded-2xl max-w-2xl w-full p-5 space-y-4 shadow-2xl text-zinc-200">
          <div class="flex items-center justify-between border-b border-[#2f323a] pb-3">
            <div class="flex items-center gap-2">
              <Layers class="w-4 h-4 text-cyan-400" />
              <h3 class="font-bold text-sm text-zinc-100">
                Voucher Jurnal Umum Sumber: {{ activeJournal?.journalNumber }}
              </h3>
            </div>
            <button @click="showJournalPreviewModal = false" class="text-zinc-400 hover:text-white transition cursor-pointer">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="text-xs space-y-1 text-zinc-400">
            <p><strong class="text-zinc-200">Tanggal Posting:</strong> {{ formatDate(activeJournal?.journalDate || activeJournal?.transaction_date) }}</p>
            <p><strong class="text-zinc-200">Keterangan:</strong> {{ activeJournal?.description }}</p>
          </div>

          <div class="border border-[#2f323a] rounded-xl overflow-hidden text-xs">
            <table class="w-full text-left">
              <thead class="bg-[#18191d] border-b border-[#2f323a] text-[10px] font-bold text-zinc-400 uppercase">
                <tr>
                  <th class="px-3 py-2">Akun Buku Besar (COA)</th>
                  <th class="px-3 py-2 text-right w-28">Debet (Rp)</th>
                  <th class="px-3 py-2 text-right w-28">Kredit (Rp)</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#282a31] text-zinc-300 font-mono">
                <tr v-for="(line, lIdx) in activeJournal?.lines || []" :key="lIdx">
                  <td class="px-3 py-2 font-sans">
                    <span class="font-bold text-zinc-100">[{{ line.accountCode }}]</span> {{ line.accountName }}
                  </td>
                  <td class="px-3 py-2 text-right">{{ line.debit > 0 ? formatIDR(line.debit) : '-' }}</td>
                  <td class="px-3 py-2 text-right">{{ line.credit > 0 ? formatIDR(line.credit) : '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="flex justify-end pt-2">
            <button
              @click="showJournalPreviewModal = false"
              class="px-4 py-2 bg-zinc-800 text-white rounded-xl text-xs font-bold hover:bg-zinc-700 transition cursor-pointer"
            >
              Tutup
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- END OF SCREEN UI CONTAINER (print:hidden) -->

    <!-- ========================================================================= -->
    <!-- CETAK RESMI FORMULIR PELAPORAN SPT MASA (DJP / BRAYO AUDIT GRADE)         -->
    <!-- HANYA MUNCUL PADA SAAT PRINT (PRINT MEDIA ONLY)                           -->
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
            <div class="font-bold text-sm">{{ reportNumber }}</div>
            <div class="text-[9px] text-gray-600">STATUS: {{ status === 'DATA_DIAMBIL' ? 'TERVERIFIKASI JURNAL' : status }}</div>
            <div class="text-[9px] font-bold text-gray-800 mt-0.5">NPWP: {{ currentEntity.npwp }}</div>
          </div>
        </div>
      </div>

      <!-- Judul Dokumen Resmi -->
      <div class="text-center my-4 space-y-0.5">
        <h1 class="text-base font-black uppercase underline tracking-wider">
          SURAT PEMBERITAHUAN (SPT) MASA PAJAK & BUKTI REKONSILIASI JURNAL
        </h1>
        <p class="text-[11px] italic text-gray-700">
          Formulir Resmi Pelaporan Pajak Bulanan Terintegrasi Jurnal Umum & Buku Besar
        </p>
      </div>

      <!-- Identitas & Parameter Dokumen (2 Kolom) -->
      <div class="grid grid-cols-2 gap-4 text-[11px] mb-4 border border-black p-2.5 bg-gray-50/50">
        <div class="space-y-1">
          <div class="flex"><span class="w-36 font-bold">Nama Dokumen / Pajak:</span><span>{{ reportName }}</span></div>
          <div class="flex"><span class="w-36 font-bold">Kode Akun Pajak:</span><span class="font-mono">[{{ selectedTaxAccountCode }}] {{ activeTaxAccountCatalog.label }}</span></div>
          <div class="flex"><span class="w-36 font-bold">Masa / Tahun Pajak:</span><span>Masa {{ String(masaMonth).padStart(2, '0') }} / {{ masaYear }}</span></div>
          <div class="flex"><span class="w-36 font-bold">Unit Bisnis / Entitas:</span><span>{{ currentEntity.unitLabel }}</span></div>
        </div>
        <div class="space-y-1">
          <div class="flex"><span class="w-36 font-bold">Kanal Bank Penyetor:</span><span>{{ activeBankAccount.name }} ([{{ activeBankAccount.code }}])</span></div>
          <div class="flex"><span class="w-36 font-bold">Nomor NTPN / DJP:</span><span class="font-mono font-bold">{{ ntpnNumber || '(NTPN-BCA-TERVERIFIKASI)' }}</span></div>
          <div class="flex"><span class="w-36 font-bold">Tanggal Pelaporan:</span><span>{{ createdDate }}</span></div>
          <div class="flex"><span class="w-36 font-bold">Sumber Transaksi:</span><span>{{ sourceTransactionFilter }} (Buku Besar)</span></div>
        </div>
      </div>

      <!-- Tabel Rincian Transaksi Per Jurnal -->
      <div class="mb-4">
        <h4 class="text-[11px] font-bold uppercase mb-1.5">I. DAFTAR RINCIAN TRANSAKSI PER VOUCHER JURNAL UMUM</h4>
        <table class="w-full text-[10px] border-collapse border border-black">
          <thead>
            <tr class="bg-gray-200 text-center font-bold">
              <th class="border border-black px-2 py-1 w-7">No</th>
              <th class="border border-black px-2 py-1 w-20">Tanggal</th>
              <th class="border border-black px-2 py-1 w-28">No. Bukti Jurnal</th>
              <th class="border border-black px-2 py-1 w-16">Unit</th>
              <th class="border border-black px-2 py-1">Akun Lawan & Uraian Transaksi</th>
              <th class="border border-black px-2 py-1 w-24 text-right">DPP (Rp)</th>
              <th class="border border-black px-2 py-1 w-24 text-right">Pajak (Rp)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in rows" :key="row.id">
              <td class="border border-black px-2 py-1 text-center">{{ idx + 1 }}</td>
              <td class="border border-black px-2 py-1 font-mono text-center">{{ row.date }}</td>
              <td class="border border-black px-2 py-1 font-mono font-bold">{{ row.journalNumber || '-' }}</td>
              <td class="border border-black px-2 py-1 text-center">{{ row.project }}</td>
              <td class="border border-black px-2 py-1">
                <div class="font-bold text-[9px]">[{{ row.counterpartCode }}] {{ row.counterpartName }}</div>
                <div class="text-[9px] text-gray-700 italic">{{ row.description }}</div>
              </td>
              <td class="border border-black px-2 py-1 text-right font-mono">{{ formatIDR(row.dppFinal) }}</td>
              <td class="border border-black px-2 py-1 text-right font-mono font-bold">{{ formatIDR(row.pajakFinal) }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-gray-100 font-bold text-[10px]">
              <td colspan="5" class="border border-black px-3 py-1 text-right uppercase">
                Total Nilai Pelaporan Pajak Final:
              </td>
              <td class="border border-black px-2 py-1 text-right font-mono">
                Rp {{ formatIDR(totalDppFinal) }}
              </td>
              <td class="border border-black px-2 py-1 text-right font-mono">
                Rp {{ formatIDR(totalPajakFinal) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Ringkasan Rekonsiliasi 3 Kolom & Ayat Jurnal Bank -->
      <div class="grid grid-cols-2 gap-4 text-[11px] mb-6">
        <div class="border border-black p-2.5 space-y-1">
          <h5 class="font-bold uppercase text-[10px] border-b border-black pb-0.5">Rekonsiliasi Nilai Pajak (Brayo Standard)</h5>
          <div class="flex justify-between"><span>Total DPP Original:</span><span class="font-mono">Rp {{ formatIDR(totalDppOriginal) }}</span></div>
          <div class="flex justify-between"><span>Penyesuaian (Adjustment):</span><span class="font-mono">Rp {{ formatIDR(totalDppAdjustment) }}</span></div>
          <div class="flex justify-between font-bold border-t border-black/30 pt-0.5"><span>Total DPP Final:</span><span class="font-mono">Rp {{ formatIDR(totalDppFinal) }}</span></div>
          <div class="flex justify-between font-bold text-xs pt-0.5 border-t border-black"><span>Total Pajak Terutang:</span><span class="font-mono">Rp {{ formatIDR(totalPajakFinal) }}</span></div>
        </div>

        <div class="border border-black p-2.5 space-y-1">
          <h5 class="font-bold uppercase text-[10px] border-b border-black pb-0.5">Ayat Jurnal Penyetoran Bank Kas Negara</h5>
          <div class="text-[10px] space-y-0.5">
            <div class="flex justify-between"><span class="font-bold">(Dr) [{{ selectedTaxAccountCode }}] {{ activeTaxAccountCatalog.label }}</span><span class="font-mono">Rp {{ formatIDR(totalPajakFinal) }}</span></div>
            <div class="flex justify-between pl-3"><span class="font-bold">(Cr) [{{ activeBankAccount.code }}] {{ activeBankAccount.name }}</span><span class="font-mono">Rp {{ formatIDR(totalPajakFinal) }}</span></div>
            <div class="text-gray-600 italic text-[9px] pt-1">NTPN Resmi: {{ ntpnNumber || '(Nomor NTPN DJP Terlampir)' }}</div>
          </div>
        </div>
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

<style>
@media print {
  /* Hide all app chrome during print */
  aside,
  header.app-navbar,
  nav.app-navbar,
  .no-print,
  button,
  .print\:hidden {
    display: none !important;
  }

  body, html, main, #app {
    background: white !important;
    color: black !important;
    padding: 0 !important;
    margin: 0 !important;
    height: auto !important;
    overflow: visible !important;
  }

  .print-document {
    display: block !important;
    width: 100% !important;
    margin: 0 auto !important;
    padding: 10mm 15mm !important;
  }

  .break-inside-avoid {
    break-inside: avoid !important;
    page-break-inside: avoid !important;
  }
}
</style>
