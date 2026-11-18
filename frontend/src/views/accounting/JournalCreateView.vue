<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useJournalsStore } from '@/stores/journals'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { BusinessUnit, AccountType, NormalBalance } from '@/types/accounting'

export type JournalType =
  | 'GENERAL'
  | 'ADJUSTMENT'
  | 'CASH_RECEIPT'
  | 'CASH_PAYMENT'
  | 'BANK_PAYMENT'
  | 'OPENING_BALANCE'
  | 'CLOSING'
  | 'REVERSING'
import {
  ArrowLeft,
  BookOpen,
  Plus,
  Trash2,
  Save,
  CheckCircle2,
  AlertTriangle,
  Sparkles,
  Calendar,
  Building2,
  FileText,
  RotateCcw,
  Layers,
  Send,
  HelpCircle
} from 'lucide-vue-next'

const router = useRouter()
const journalsStore = useJournalsStore()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await Promise.all([
    ledgerStore.fetchAccounts(),
    journalsStore.fetchJournals()
  ])
  setupDefaultRows()
})

// Header Form State
const postingDate = ref<string>(new Date().toISOString().slice(0, 10))
const journalType = ref<JournalType>('GENERAL')
const referenceNo = ref<string>('')
const businessUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const description = ref<string>('')

// Available active leaf accounts
const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

const accountSelectOptions = computed(() => {
  return availableAccounts.value.map(acc => ({
    value: acc.code,
    label: `[${acc.code}] ${acc.name}`,
    sublabel: `${acc.type} • Normal: ${acc.normalBalance}`
  }))
})

// Multi-line Journal Rows
interface JournalRowState {
  accountCode: string
  accountName: string
  accountType?: AccountType
  normalBalance?: NormalBalance
  unit: BusinessUnit
  description: string
  debit: number | string
  credit: number | string
}

const rows = ref<JournalRowState[]>([])

function setupDefaultRows() {
  if (rows.value.length === 0) {
    rows.value = [
      {
        accountCode: '',
        accountName: '',
        unit: businessUnit.value,
        description: '',
        debit: 0,
        credit: 0
      },
      {
        accountCode: '',
        accountName: '',
        unit: businessUnit.value,
        description: '',
        debit: 0,
        credit: 0
      }
    ]
  }
}

function onAccountSelect(index: number, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc && rows.value[index]) {
    rows.value[index].accountCode = acc.code
    rows.value[index].accountName = acc.name
    rows.value[index].accountType = acc.type
    rows.value[index].normalBalance = acc.normalBalance
    if (!rows.value[index].description) {
      rows.value[index].description = description.value || acc.name
    }
  }
}

function addRow() {
  rows.value.push({
    accountCode: '',
    accountName: '',
    unit: businessUnit.value,
    description: description.value || '',
    debit: 0,
    credit: 0
  })
}

function removeRow(index: number) {
  if (rows.value.length > 2) {
    rows.value.splice(index, 1)
  } else {
    rows.value[index] = {
      accountCode: '',
      accountName: '',
      unit: businessUnit.value,
      description: '',
      debit: 0,
      credit: 0
    }
  }
}

// Calculations
const totalDebit = computed(() => {
  return rows.value.reduce((sum, r) => sum + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return rows.value.reduce((sum, r) => sum + (Number(r.credit) || 0), 0)
})

const difference = computed(() => {
  return Math.round((totalDebit.value - totalCredit.value) * 100) / 100
})

const isBalanced = computed(() => {
  return Math.abs(difference.value) < 0.01 && totalDebit.value > 0
})

// Auto Balance Helper
function autoBalanceDifference() {
  const diff = difference.value
  if (Math.abs(diff) < 0.01) {
    accountingStore.showToast('Jurnal sudah seimbang!', 'info')
    return
  }

  // Look for an empty row or append new row
  let targetRow = rows.value.find(r => !r.accountCode && !r.debit && !r.credit)
  if (!targetRow) {
    targetRow = {
      accountCode: '',
      accountName: '',
      unit: businessUnit.value,
      description: 'Penyeimbang Jurnal',
      debit: 0,
      credit: 0
    }
    rows.value.push(targetRow)
  }

  if (diff > 0) {
    // Debit is higher -> add to Credit
    targetRow.credit = diff
    targetRow.debit = 0
  } else {
    // Credit is higher -> add to Debit
    targetRow.debit = Math.abs(diff)
    targetRow.credit = 0
  }

  accountingStore.showToast(`Selisih Rp ${formatIDR(Math.abs(diff))} dialokasikan ke baris baru. Pilih akun penyeimbang.`, 'info')
}

// Preset Templates
function loadTemplate(templateType: 'DOCTOR_SHARE' | 'PHARMACY_PURCHASE' | 'PETTY_CASH_REFILL') {
  if (templateType === 'DOCTOR_SHARE') {
    journalType.value = 'ADJUSTMENT'
    description.value = 'Akrual Bagi Hasil Jasa Medis Dokter Spesialis & Withholding PPh 21'
    referenceNo.value = `DOC-SETTLE-${new Date().toISOString().slice(0, 7)}`
    rows.value = [
      {
        accountCode: '5201',
        accountName: 'Beban Bagi Hasil Dokter Spesialis Anak & Kulit',
        unit: 'ALL',
        description: 'Beban Honor Jasa Medis Dokter',
        debit: 25000000,
        credit: 0
      },
      {
        accountCode: '2501.001',
        accountName: 'PPh 21 Dokter Tenaga Medis Ahli (PMK 168/2023)',
        unit: 'ALL',
        description: 'Potongan Withholding PPh 21 Dokter 5%',
        debit: 0,
        credit: 625000
      },
      {
        accountCode: '2301',
        accountName: 'Hutang Bagi Hasil Dokter',
        unit: 'ALL',
        description: 'Hutang Bersih Jasa Dokter Siap Transfer',
        debit: 0,
        credit: 24375000
      }
    ]
  } else if (templateType === 'PHARMACY_PURCHASE') {
    journalType.value = 'GENERAL'
    description.value = 'Pembelian Obat Resep & Vaksin Anak dari PBF Distributor Farmasi (Tempo)'
    referenceNo.value = `PO-KF-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}`
    rows.value = [
      {
        accountCode: '1301',
        accountName: 'Persediaan Obat Resep & Paten Farmasi Kidinest',
        unit: 'KIDINEST',
        description: 'Penerimaan Stok Obat dari Kimia Farma',
        debit: 12500000,
        credit: 0
      },
      {
        accountCode: '1504',
        accountName: 'PPN Masukan (Faktur Pajak Pembelian B3/BMHP)',
        unit: 'KIDINEST',
        description: 'PPN Masukan 11% Faktur Pajak PBF',
        debit: 1375000,
        credit: 0
      },
      {
        accountCode: '2201',
        accountName: 'Hutang PBF Distributor Farmasi & Obat',
        unit: 'KIDINEST',
        description: 'Hutang Usaha Kimia Farma (Tempo 30 Hari)',
        debit: 0,
        credit: 13875000
      }
    ]
  } else if (templateType === 'PETTY_CASH_REFILL') {
    journalType.value = 'CASH_PAYMENT'
    description.value = 'Pengisian Kembali Kas Kecil (Petty Cash Imprest Fund) dari Bank BCA'
    referenceNo.value = `TRF-BCA-PC-${new Date().toISOString().slice(0, 10).replace(/-/g, '')}`
    rows.value = [
      {
        accountCode: '1101.003',
        accountName: 'Kas Kecil (Petty Cash Operasional)',
        unit: 'ALL',
        description: 'Penerimaan Dana Kas Kecil',
        debit: 5000000,
        credit: 0
      },
      {
        accountCode: '1102.001',
        accountName: 'Bank BCA Operasional (Giro & EDC)',
        unit: 'ALL',
        description: 'Penarikan Tunai Giro Bank BCA',
        debit: 0,
        credit: 5000000
      }
    ]
  }
  accountingStore.showToast(`Template '${templateType}' berhasil dimuat!`, 'success')
}

// Submit Form
const isSubmitting = ref(false)

async function handleSubmit(autoPost: boolean = false) {
  const validRows = rows.value.filter(r => r.accountCode && (Number(r.debit) > 0 || Number(r.credit) > 0))

  if (validRows.length < 2) {
    accountingStore.showToast('Jurnal minimal memiliki 2 baris akun yang valid!', 'error')
    return
  }

  if (!isBalanced.value) {
    accountingStore.showToast(
      `Jurnal belum seimbang! Total Debit (Rp ${formatIDR(totalDebit.value)}) != Total Kredit (Rp ${formatIDR(totalCredit.value)}). Selisih: Rp ${formatIDR(Math.abs(difference.value))}`,
      'error'
    )
    return
  }

  try {
    isSubmitting.value = true

    await journalsStore.createManualJournal({
      journalDate: postingDate.value,
      businessUnit: businessUnit.value,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: description.value || `Jurnal ${journalType.value} tanggal ${postingDate.value}`,
      journalSource: 'MANUAL',
      sourceReference: referenceNo.value || `REF-${postingDate.value.replace(/-/g, '')}`,
      status: autoPost ? 'POSTED' : 'DRAFT',
      lines: validRows.map(r => ({
        accountCode: r.accountCode,
        accountName: r.accountName,
        debit: Number(r.debit) || 0,
        credit: Number(r.credit) || 0,
        description: r.description || description.value
      }))
    })

    accountingStore.showToast(
      `Jurnal umum berhasil dibuat & ${autoPost ? 'langsung diposting ke Buku Besar!' : 'disimpan sebagai draf review!'}`,
      'success'
    )

    router.push('/accounting/journals')
  } catch (err: any) {
    accountingStore.showToast(`Gagal menyimpan jurnal: ${err.message}`, 'error')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-24">
    <!-- Header with Back Button -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="flex items-center gap-3">
        <button
          @click="router.push('/accounting/journals')"
          class="p-2 rounded-xl bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          title="Kembali ke Daftar Jurnal"
        >
          <ArrowLeft class="w-4 h-4" />
        </button>
        <div>
          <div class="flex items-center gap-2">
            <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
              Buku & Transaksi
            </span>
            <span class="text-zinc-300">•</span>
            <span class="text-xs text-zinc-400 font-medium">Formulir Voucher Jurnal Standar Brayo</span>
          </div>
          <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
            Buat Jurnal Umum Baru (Journal Entry)
          </h1>
        </div>
      </div>

      <!-- Quick Template Presets -->
      <div class="flex flex-wrap items-center gap-2">
        <div class="relative group">
          <button
            type="button"
            class="inline-flex items-center gap-1.5 px-3 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold transition border border-zinc-200 shadow-2xs cursor-pointer"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-500" />
            <span>Muat Template Transaksi</span>
          </button>
          <div class="absolute right-0 top-full mt-1 w-64 bg-white border border-zinc-200 rounded-xl shadow-lg p-1.5 hidden group-hover:block z-30 text-xs">
            <button
              @click="loadTemplate('DOCTOR_SHARE')"
              class="w-full text-left px-3 py-2 rounded-lg hover:bg-zinc-100 font-medium text-zinc-800 block"
            >
              🩺 Bagi Hasil Dokter & PPh 21
            </button>
            <button
              @click="loadTemplate('PHARMACY_PURCHASE')"
              class="w-full text-left px-3 py-2 rounded-lg hover:bg-zinc-100 font-medium text-zinc-800 block"
            >
              💊 Pembelian Obat PBF + PPN
            </button>
            <button
              @click="loadTemplate('PETTY_CASH_REFILL')"
              class="w-full text-left px-3 py-2 rounded-lg hover:bg-zinc-100 font-medium text-zinc-800 block"
            >
              💵 Pengisian Kas Kecil dari Bank
            </button>
          </div>
        </div>

        <button
          @click="autoBalanceDifference"
          :disabled="isBalanced || difference === 0"
          :class="[
            'inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95',
            !isBalanced && difference !== 0
              ? 'bg-amber-500 hover:bg-amber-600 text-white cursor-pointer'
              : 'bg-zinc-100 text-zinc-400 border border-zinc-200 cursor-not-allowed'
          ]"
        >
          <Sparkles class="w-3.5 h-3.5" />
          <span>Seimbangkan Selisih</span>
        </button>
      </div>
    </header>

    <!-- Header Form Card (Voucher Metadata) -->
    <div class="bg-white p-5 rounded-2xl border border-zinc-200/90 shadow-2xs grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Tanggal Posting (Posting Date):</span>
        </label>
        <input
          v-model="postingDate"
          type="date"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <BookOpen class="w-3.5 h-3.5 text-zinc-500" />
          <span>Tipe Entri Jurnal:</span>
        </label>
        <select
          v-model="journalType"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="GENERAL">GENERAL (Jurnal Umum Standar)</option>
          <option value="ADJUSTMENT">ADJUSTMENT (Jurnal Penyesuaian)</option>
          <option value="CASH_RECEIPT">CASH_RECEIPT (Penerimaan Kasir/BKM)</option>
          <option value="CASH_PAYMENT">CASH_PAYMENT (Pengeluaran Kas/BKK)</option>
          <option value="BANK_PAYMENT">BANK_PAYMENT (Mutasi/Pembayaran Bank)</option>
          <option value="OPENING_BALANCE">OPENING_BALANCE (Saldo Awal)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <FileText class="w-3.5 h-3.5 text-zinc-500" />
          <span>No. Referensi / Bukti:</span>
        </label>
        <input
          v-model="referenceNo"
          type="text"
          placeholder="e.g. INV-2026-001, MEMO-01"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-zinc-500" />
          <span>Unit Bisnis / Cabang:</span>
        </label>
        <select
          v-model="businessUnit"
          class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Semua Unit (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Poli Anak)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div class="md:col-span-4">
        <label class="block font-bold text-zinc-700 mb-1">Catatan Pengguna / Keterangan Umum Transaksi:</label>
        <input
          v-model="description"
          type="text"
          placeholder="Tuliskan keterangan lengkap maksud transaksi jurnal ini..."
          class="w-full px-3.5 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
        />
      </div>
    </div>

    <!-- Balance Status Bar (Seimbang / Tidak Seimbang) -->
    <div
      :class="[
        'p-4 rounded-2xl border transition-all shadow-2xs',
        isBalanced
          ? 'bg-emerald-50/70 border-emerald-200 text-emerald-950'
          : 'bg-amber-50/70 border-amber-200 text-amber-950'
      ]"
    >
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs">
        <div class="flex items-center gap-3">
          <div
            :class="[
              'w-9 h-9 rounded-xl flex items-center justify-center font-bold text-white shrink-0',
              isBalanced ? 'bg-emerald-600 shadow-2xs' : 'bg-amber-600 shadow-2xs'
            ]"
          >
            <CheckCircle2 v-if="isBalanced" class="w-5 h-5" />
            <AlertTriangle v-else class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-extrabold uppercase tracking-wider text-[11px]">
                {{ isBalanced ? 'STATUS: JURNAL SEIMBANG (BALANCED)' : 'STATUS: BELUM SEIMBANG' }}
              </span>
              <span
                v-if="isBalanced"
                class="px-2 py-0.5 rounded-md bg-emerald-100 text-emerald-800 font-bold text-[10px]"
              >
                Total Debet == Total Kredit
              </span>
              <span
                v-else
                class="px-2 py-0.5 rounded-md bg-amber-100 text-amber-800 font-bold text-[10px]"
              >
                Selisih: {{ formatIDR(Math.abs(difference)) }}
              </span>
            </div>
            <p class="text-[11px] mt-0.5 opacity-80">
              {{ isBalanced ? 'Jurnal memenuhi kaidah akuntansi double-entry dan siap disimpan atau diposting.' : 'Total Debet dan Kredit harus bernilai sama persis sebelum dapat diposting.' }}
            </p>
          </div>
        </div>

        <!-- Metric Numbers -->
        <div class="flex items-center gap-4 bg-white/90 px-4 py-2 rounded-xl border border-zinc-200/80 shadow-2xs self-end sm:self-auto font-mono">
          <div>
            <span class="text-[10px] text-zinc-500 font-sans block font-bold">TOTAL DEBET</span>
            <span class="font-extrabold text-zinc-950 text-sm">{{ formatIDR(totalDebit) }}</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] text-zinc-500 font-sans block font-bold">TOTAL KREDIT</span>
            <span class="font-extrabold text-zinc-950 text-sm">{{ formatIDR(totalCredit) }}</span>
          </div>
          <div class="h-6 w-px bg-zinc-200"></div>
          <div>
            <span class="text-[10px] text-zinc-500 font-sans block font-bold">SELISIH</span>
            <span :class="['font-extrabold text-sm', isBalanced ? 'text-emerald-700' : 'text-rose-600']">
              {{ formatIDR(Math.abs(difference)) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Accounting Entries Table (Tabel Rincian Baris Jurnal) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="p-4 border-b border-zinc-200/90 flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-zinc-50/50">
        <div class="flex items-center gap-2">
          <BookOpen class="w-4 h-4 text-zinc-700" />
          <span class="text-xs font-bold text-zinc-900">
            Daftar Posisi Akun Jurnal ({{ rows.length }} Baris Akuntansi)
          </span>
          <span class="text-[11px] text-zinc-500 hidden sm:inline">
            &bull; Masukkan akun debit dan kredit berpasangan seimbang
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="addRow"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Tambah Baris</span>
          </button>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-xs text-left">
          <thead class="bg-zinc-50 text-zinc-600 font-bold border-b border-zinc-200/90 uppercase tracking-wider text-[11px]">
            <tr>
              <th class="px-3 py-3 w-10 text-center">#</th>
              <th class="px-4 py-3 min-w-[280px]">Akun Buku Besar (COA BCS)</th>
              <th class="px-3 py-3 w-28 text-center">Unit</th>
              <th class="px-4 py-3 min-w-[200px]">Uraian Baris (Line Remark)</th>
              <th class="px-4 py-3 w-48 text-right">Debet (Rp)</th>
              <th class="px-4 py-3 w-48 text-right">Kredit (Rp)</th>
              <th class="px-3 py-3 w-12 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 text-zinc-700">
            <tr
              v-for="(row, idx) in rows"
              :key="idx"
              class="hover:bg-zinc-50/70 transition"
            >
              <td class="px-3 py-3 text-center font-bold text-zinc-400">
                {{ idx + 1 }}
              </td>

              <!-- Account Picker -->
              <td class="px-4 py-2.5">
                <SearchableSelect
                  :model-value="row.accountCode"
                  :options="accountSelectOptions"
                  placeholder="-- Pilih Akun COA --"
                  @update:model-value="(code: any) => onAccountSelect(idx, String(code))"
                />
              </td>

              <!-- Unit Selector -->
              <td class="px-3 py-2.5 text-center">
                <select
                  v-model="row.unit"
                  class="w-full px-2 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-[11px] font-bold text-zinc-800 focus:outline-none"
                >
                  <option value="ALL">ALL</option>
                  <option value="KIDINEST">Kidinest</option>
                  <option value="SSC">SSC</option>
                </select>
              </td>

              <!-- Line Description -->
              <td class="px-4 py-2.5">
                <input
                  v-model="row.description"
                  type="text"
                  placeholder="Keterangan baris transaksi..."
                  class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>

              <!-- Debit Input -->
              <td class="px-4 py-2.5 text-right">
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                  <input
                    v-model.number="row.debit"
                    type="number"
                    min="0"
                    step="any"
                    class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                    placeholder="0"
                  />
                </div>
              </td>

              <!-- Credit Input -->
              <td class="px-4 py-2.5 text-right">
                <div class="relative">
                  <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                  <input
                    v-model.number="row.credit"
                    type="number"
                    min="0"
                    step="any"
                    class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                    placeholder="0"
                  />
                </div>
              </td>

              <!-- Delete Button -->
              <td class="px-3 py-2.5 text-center">
                <button
                  type="button"
                  @click="removeRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition"
                  title="Hapus baris"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          </tbody>
          <tfoot class="bg-zinc-50 border-t-2 border-zinc-200 font-bold text-xs">
            <tr>
              <td colspan="4" class="px-4 py-3 text-right uppercase text-zinc-600 tracking-wider">
                Total Jurnal:
              </td>
              <td class="px-4 py-3 text-right font-mono text-emerald-800 text-sm">
                {{ formatIDR(totalDebit) }}
              </td>
              <td class="px-4 py-3 text-right font-mono text-rose-800 text-sm">
                {{ formatIDR(totalCredit) }}
              </td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Action Buttons Footer -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pt-2">
      <button
        type="button"
        @click="router.push('/accounting/journals')"
        class="px-4 py-2.5 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition"
      >
        Batal / Kembali
      </button>

      <div class="flex items-center gap-2">
        <button
          type="button"
          :disabled="isSubmitting || !isBalanced"
          @click="handleSubmit(false)"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 bg-zinc-100 hover:bg-zinc-200 disabled:opacity-50 text-zinc-800 rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
        >
          <Save class="w-4 h-4 text-zinc-600" />
          <span>Simpan sebagai Draft</span>
        </button>

        <button
          type="button"
          :disabled="isSubmitting || !isBalanced"
          @click="handleSubmit(true)"
          class="inline-flex items-center gap-1.5 px-5 py-2.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white rounded-xl text-xs font-bold transition shadow-md cursor-pointer active:scale-95"
        >
          <Send class="w-4 h-4 text-emerald-400" />
          <span>Simpan & Posting ke Buku Besar</span>
        </button>
      </div>
    </div>
  </div>
</template>
