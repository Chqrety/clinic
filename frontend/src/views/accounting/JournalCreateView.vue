<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useJournalsStore } from '@/stores/journals'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import CurrencyInput from '@/components/common/CurrencyInput.vue'
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
  Folder,
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
  HelpCircle,
  Calculator,
  Clock,
  Check
} from 'lucide-vue-next'

const router = useRouter()
const journalsStore = useJournalsStore()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

// Active tab matching Brayo SPT Masa layout
const activeTab = ref<'DETAILS' | 'RECONCILIATION' | 'AUDIT'>('DETAILS')

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
  <div class="font-['Plus_Jakarta_Sans',sans-serif] max-w-7xl mx-auto pb-16">
    <div class="bg-white text-zinc-800 rounded-2xl border border-zinc-200/90 p-5 sm:p-7 shadow-2xs space-y-6">
      
      <!-- 1. Header: Breadcrumb & Action Bar -->
      <header class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-3 border-b border-zinc-200">
        <!-- Breadcrumb -->
        <div class="flex items-center gap-2 text-xs text-zinc-500 flex-wrap">
          <router-link to="/accounting/journals" class="hover:text-zinc-900 transition flex items-center gap-1.5 text-zinc-500">
            <Folder class="w-4 h-4 text-zinc-400" />
            <span>Buku & Transaksi</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-500">Jurnal Umum</span>
          <span>/</span>
          <span class="font-mono font-bold text-zinc-900">{{ referenceNo || 'JV-BARU' }}</span>
          <span
            :class="[
              'ml-1.5 px-2 py-0.5 rounded text-[11px] font-bold border',
              isBalanced ? 'bg-emerald-50 text-emerald-700 border-emerald-200' : 'bg-amber-50 text-amber-700 border-amber-200'
            ]"
          >
            {{ isBalanced ? 'Seimbang' : 'Belum Seimbang' }}
          </span>
        </div>

        <!-- Top Right Actions -->
        <div class="flex items-center gap-2 flex-wrap">
          <button
            v-if="!isBalanced && difference !== 0"
            type="button"
            @click="autoBalanceDifference"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold rounded-xl text-xs border border-amber-200 transition active:scale-95 cursor-pointer"
            title="Seimbangkan selisih baris jurnal secara otomatis"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Seimbangkan Selisih</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting || !isBalanced"
            @click="handleSubmit(false)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-white hover:bg-zinc-50 disabled:opacity-50 text-zinc-800 font-bold rounded-xl text-xs border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Save class="w-3.5 h-3.5 text-zinc-600" />
            <span>Simpan Draft</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting || !isBalanced"
            @click="handleSubmit(true)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white font-bold rounded-xl text-xs shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Send class="w-3.5 h-3.5 text-emerald-400" />
            <span>Simpan & Posting</span>
          </button>

          <router-link
            to="/accounting/journals"
            class="inline-flex items-center px-3 py-1.5 text-zinc-500 hover:text-zinc-800 text-xs font-semibold transition"
          >
            Cancel
          </router-link>
        </div>
      </header>

      <!-- 2. Tabs Navigation -->
      <div class="border-b border-zinc-200">
        <nav class="flex space-x-6 text-xs font-medium">
          <button
            type="button"
            @click="activeTab = 'DETAILS'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'DETAILS' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Details</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'RECONCILIATION'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'RECONCILIATION' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <Calculator class="w-3.5 h-3.5" />
            <span>Rekonsiliasi & Status</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'AUDIT'"
            class="pb-2.5 transition flex items-center gap-1.5 cursor-pointer"
            :class="activeTab === 'AUDIT' ? 'border-b-2 border-zinc-950 text-zinc-950 font-bold' : 'border-b-2 border-transparent text-zinc-500 hover:text-zinc-800'"
          >
            <Clock class="w-3.5 h-3.5" />
            <span>Audit Trail & Log</span>
          </button>
        </nav>
      </div>

      <!-- 3. TAB 1: DETAILS -->
      <div v-if="activeTab === 'DETAILS'" class="space-y-7">
        
        <!-- Section A: Informasi Voucher Jurnal -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Informasi Voucher Jurnal
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Tanggal Posting *</label>
              <input
                v-model="postingDate"
                type="date"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Tipe Entri Jurnal *</label>
              <select
                v-model="journalType"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
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
              <label class="block text-zinc-600 mb-1 font-medium">No. Referensi / Bukti</label>
              <input
                v-model="referenceNo"
                type="text"
                placeholder="e.g. INV-2026-001, MEMO-01"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Unit Bisnis / Cabang *</label>
              <select
                v-model="businessUnit"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
              >
                <option value="ALL">Semua Unit (Konsolidasi)</option>
                <option value="KIDINEST">Unit Kidinest (Klinik Spesialis Anak)</option>
                <option value="SSC">Unit SSC (Semarang Skin Center)</option>
              </select>
            </div>

            <div class="md:col-span-2 lg:col-span-4">
              <label class="block text-zinc-600 mb-1 font-medium">Catatan / Keterangan Umum Transaksi</label>
              <input
                v-model="description"
                type="text"
                placeholder="Tuliskan keterangan lengkap maksud transaksi jurnal ini..."
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-medium focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>
          </div>
        </div>

        <!-- Section B: Detail Transaksi Baris Jurnal -->
        <div>
          <div class="flex items-center justify-between mb-3.5">
            <div>
              <h3 class="text-sm font-bold text-zinc-900">
                Detail Transaksi Jurnal ({{ rows.length }} Baris Akuntansi)
              </h3>
              <p class="text-xs text-zinc-500">
                Masukkan akun debit dan kredit berpasangan seimbang (kaidah double-entry bookkeeping).
              </p>
            </div>

            <button
              type="button"
              @click="addRow"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Tambah Baris</span>
            </button>
          </div>

          <div class="border border-zinc-200 rounded-xl overflow-hidden shadow-2xs">
            <div class="overflow-x-auto">
              <table class="w-full text-xs text-left">
                <thead class="bg-zinc-50 text-zinc-500 font-bold border-b border-zinc-200 uppercase tracking-wider text-[11px]">
                  <tr>
                    <th class="px-3 py-2.5 w-10 text-center">#</th>
                    <th class="px-4 py-2.5 min-w-[280px]">Akun Buku Besar (COA BCS)</th>
                    <th class="px-3 py-2.5 w-28 text-center">Unit</th>
                    <th class="px-4 py-2.5 min-w-[200px]">Uraian Baris (Line Remark)</th>
                    <th class="px-4 py-2.5 w-44 text-right">Debet (Rp)</th>
                    <th class="px-4 py-2.5 w-44 text-right">Kredit (Rp)</th>
                    <th class="px-3 py-2.5 w-12 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-zinc-100 text-zinc-700 bg-white">
                  <tr
                    v-for="(row, idx) in rows"
                    :key="idx"
                    class="hover:bg-zinc-50/70 transition"
                  >
                    <td class="px-3 py-2.5 text-center font-bold text-zinc-400">
                      {{ idx + 1 }}
                    </td>

                    <td class="px-4 py-2">
                      <SearchableSelect
                        :model-value="row.accountCode"
                        :options="accountSelectOptions"
                        placeholder="-- Pilih Akun COA --"
                        @update:model-value="(code: any) => onAccountSelect(idx, String(code))"
                      />
                    </td>

                    <td class="px-3 py-2 text-center">
                      <select
                        v-model="row.unit"
                        class="w-full px-2 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-[11px] font-semibold text-zinc-800 focus:outline-none"
                      >
                        <option value="ALL">ALL</option>
                        <option value="KIDINEST">Kidinest</option>
                        <option value="SSC">SSC</option>
                      </select>
                    </td>

                    <td class="px-4 py-2">
                      <input
                        v-model="row.description"
                        type="text"
                        placeholder="Keterangan baris transaksi..."
                        class="w-full px-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-xs text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white"
                      />
                    </td>

                    <td class="px-4 py-2 text-right">
                      <CurrencyInput
                        v-model="row.debit"
                        prefix="Rp"
                        input-class="w-full pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                        placeholder="0"
                      />
                    </td>

                    <td class="px-4 py-2 text-right">
                      <CurrencyInput
                        v-model="row.credit"
                        prefix="Rp"
                        input-class="w-full pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                        placeholder="0"
                      />
                    </td>

                    <td class="px-3 py-2 text-center">
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
                <tfoot class="bg-zinc-50/90 border-t border-zinc-200 font-bold text-xs">
                  <tr>
                    <td colspan="4" class="px-4 py-3 text-right uppercase text-zinc-500 tracking-wider">
                      Total Jurnal:
                    </td>
                    <td class="px-4 py-3 text-right font-mono text-emerald-800 text-sm">
                      Rp {{ formatIDR(totalDebit) }}
                    </td>
                    <td class="px-4 py-3 text-right font-mono text-rose-800 text-sm">
                      Rp {{ formatIDR(totalCredit) }}
                    </td>
                    <td></td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>
        </div>

        <!-- Section C: Ringkasan Nilai Jurnal -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Ringkasan & Keseimbangan Jurnal
          </h3>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs font-mono">
            <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block mb-1">TOTAL DEBET</span>
              <span class="text-base font-extrabold text-emerald-700">Rp {{ formatIDR(totalDebit) }}</span>
            </div>

            <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block mb-1">TOTAL KREDIT</span>
              <span class="text-base font-extrabold text-rose-700">Rp {{ formatIDR(totalCredit) }}</span>
            </div>

            <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block mb-1">STATUS SELISIH</span>
              <span :class="['text-base font-extrabold', isBalanced ? 'text-emerald-700' : 'text-amber-700']">
                {{ isBalanced ? 'SEIMBANG (Rp 0)' : `Selisih Rp ${formatIDR(Math.abs(difference))}` }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. TAB 2: RECONCILIATION & STATUS -->
      <div v-else-if="activeTab === 'RECONCILIATION'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Status Rekonsiliasi & Aturan Pembukuan</h3>
          <p class="text-xs text-zinc-500">Validasi kepatuhan standar akuntansi keuangan (PSAK / Double-Entry Bookkeeping).</p>
        </div>

        <div class="space-y-3 text-xs">
          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 v-if="isBalanced" class="w-4 h-4 text-emerald-600" />
              <AlertTriangle v-else class="w-4 h-4 text-amber-600" />
              <span class="font-semibold text-zinc-900">Keseimbangan Debet & Kredit</span>
            </div>
            <span :class="['font-bold font-mono', isBalanced ? 'text-emerald-700' : 'text-amber-700']">
              {{ isBalanced ? 'Seimbang (Rp 0)' : `Selisih Rp ${formatIDR(Math.abs(difference))}` }}
            </span>
          </div>

          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 v-if="rows.filter(r => r.accountCode).length >= 2" class="w-4 h-4 text-emerald-600" />
              <AlertTriangle v-else class="w-4 h-4 text-amber-600" />
              <span class="font-semibold text-zinc-900">Jumlah Akun Terlibat</span>
            </div>
            <span class="font-bold text-zinc-700">{{ rows.filter(r => r.accountCode).length }} Akun Terpilih</span>
          </div>

          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span class="font-semibold text-zinc-900">Status Entitas / Unit Bisnis</span>
            </div>
            <span class="font-bold text-zinc-700">{{ businessUnit === 'ALL' ? 'Konsolidasi' : businessUnit }}</span>
          </div>
        </div>
      </div>

      <!-- 5. TAB 3: AUDIT TRAIL & LOG -->
      <div v-else-if="activeTab === 'AUDIT'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Audit Trail & Jejak Pencatatan</h3>
          <p class="text-xs text-zinc-500">Informasi pengguna, tanggal pembuatan, dan status otorisasi transaksi.</p>
        </div>

        <div class="space-y-3 text-xs">
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Petugas Pembuat:</span>
            <span class="font-bold text-zinc-900">Budi Santoso, SE, Ak., CA (Finance & Accounting)</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Waktu Pembuatan:</span>
            <span class="font-bold text-zinc-900 font-mono">{{ new Date().toLocaleString('id-ID') }}</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Metode Entri:</span>
            <span class="font-bold text-zinc-900">Input Manual Voucher Standar</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
