<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCashBankStore } from '@/stores/cashbank'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { formatIDR } from '@/utils/currency'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { BusinessUnit, CashBankTransactionType } from '@/types/accounting'
import {
  Folder,
  ArrowLeft,
  ArrowDownLeft,
  ArrowUpRight,
  ArrowLeftRight,
  Plus,
  Trash2,
  Save,
  CheckCircle2,
  AlertCircle,
  Sparkles,
  Wallet,
  Landmark,
  Building2,
  FileText,
  BookOpen,
  Calculator,
  Clock,
  Send
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const cashBankStore = useCashBankStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()

// Active tab matching Brayo SPT Masa layout
const activeTab = ref<'DETAILS' | 'RECONCILIATION' | 'AUDIT'>('DETAILS')

onMounted(async () => {
  await cashBankStore.fetchCashBankData()

  // Initialize query param type if provided
  const qType = route.query.type as string
  if (qType === 'CASH_IN' || qType === 'CASH_OUT' || qType === 'TRANSFER') {
    voucherType.value = qType
  }
  setupDefaultRows()
})

// Voucher Form Header State
const voucherType = ref<CashBankTransactionType>('CASH_OUT')
const primaryAccountCode = ref<string>('1101.003') // Default Kas Kecil Operasional
const transactionDate = ref<string>(new Date().toISOString().slice(0, 10))
const recipientOrPayer = ref<string>('')
const referenceNo = ref<string>('')
const businessUnit = ref<BusinessUnit>('KIDINEST')
const description = ref<string>('')

// Available Cash & Bank accounts for primary dropdown
const availableCashBankAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => {
    if (a.isHeader) return false
    if (a.isActive === false) return false
    return (
      a.code.startsWith('1101') ||
      a.code.startsWith('1102') ||
      a.code.startsWith('1103') ||
      a.parentCode === '1101' ||
      a.parentCode === '1102'
    )
  })
})

// Current balance of selected primary account
const selectedPrimaryAccount = computed(() => {
  return ledgerStore.accounts.find(a => a.code === primaryAccountCode.value)
})

// Available COA accounts for counterpart lines (Excluding headers and inactive accounts)
const availableAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)
})

const primaryAccountOptions = computed(() => {
  return availableCashBankAccounts.value.map(acc => ({
    value: acc.code,
    label: `[${acc.code}] ${acc.name}`,
    sublabel: `Saldo: ${formatIDR(acc.currentBalance || 0)}`
  }))
})

const accountSelectOptions = computed(() => {
  return availableAccounts.value.map(acc => ({
    value: acc.code,
    label: `[${acc.code}] ${acc.name}`,
    sublabel: `${acc.type} • Normal: ${acc.normalBalance}`
  }))
})

// Multi-line Voucher Rows
interface VoucherRow {
  accountCode: string
  accountName: string
  description: string
  debit: number | string
  credit: number | string
  unit: BusinessUnit
}

const voucherRows = ref<VoucherRow[]>([])

function setupDefaultRows() {
  const primaryAcc = selectedPrimaryAccount.value || availableCashBankAccounts.value[0]
  if (primaryAcc) {
    primaryAccountCode.value = primaryAcc.code
  }

  if (voucherType.value === 'CASH_OUT') {
    // Standard BKK: Counterpart (Expense/Payable) on Debit, Cash/Bank on Credit
    voucherRows.value = [
      {
        accountCode: '',
        accountName: '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      }
    ]
  } else if (voucherType.value === 'CASH_IN') {
    // Standard BKM: Cash/Bank on Debit, Counterpart (Revenue/Receivable) on Credit
    voucherRows.value = [
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: '',
        accountName: '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      }
    ]
  } else {
    // Standard Transfer: Destination Bank on Debit, Source Bank/Cash on Credit
    voucherRows.value = [
      {
        accountCode: '',
        accountName: '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || '',
        description: '',
        debit: 0,
        credit: 0,
        unit: businessUnit.value
      }
    ]
  }
}

// Watch voucher type change
function onVoucherTypeChange(newType: CashBankTransactionType) {
  voucherType.value = newType
  setupDefaultRows()
}

// Watch primary account change
function onPrimaryAccountChange(code: string) {
  primaryAccountCode.value = code
  const acc = ledgerStore.accounts.find(a => a.code === code)
  if (!acc) return

  // Update corresponding cash/bank row in the table
  if (voucherType.value === 'CASH_OUT') {
    // Look for credit row that is a cash/bank account
    const cbRow = voucherRows.value.find(r => Number(r.credit) > 0)
    if (cbRow) {
      cbRow.accountCode = acc.code
      cbRow.accountName = acc.name
    }
  } else if (voucherType.value === 'CASH_IN') {
    // Look for debit row that is a cash/bank account
    const cbRow = voucherRows.value.find(r => Number(r.debit) > 0)
    if (cbRow) {
      cbRow.accountCode = acc.code
      cbRow.accountName = acc.name
    }
  }
}

function onAccountSelect(idx: number, code: string) {
  const acc = ledgerStore.accounts.find(a => a.code === code)
  if (acc && voucherRows.value[idx]) {
    voucherRows.value[idx].accountCode = acc.code
    voucherRows.value[idx].accountName = acc.name
  }
}

function addRow() {
  voucherRows.value.push({
    accountCode: '',
    accountName: '',
    description: '',
    debit: 0,
    credit: 0,
    unit: businessUnit.value
  })
}

function removeRow(idx: number) {
  if (voucherRows.value.length <= 2) {
    accountingStore.showToast('Voucher minimal memiliki 2 baris (Debit & Kredit)', 'info')
    return
  }
  voucherRows.value.splice(idx, 1)
}

// Totals & Balance Computations
const totalDebit = computed(() => {
  return voucherRows.value.reduce((s, r) => s + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return voucherRows.value.reduce((s, r) => s + (Number(r.credit) || 0), 0)
})

const isBalanced = computed(() => {
  return totalDebit.value > 0 && Math.abs(totalDebit.value - totalCredit.value) < 0.01
})

const balanceDiff = computed(() => {
  return Math.abs(totalDebit.value - totalCredit.value)
})

// Auto-balance button helper
function autoBalanceRows() {
  if (voucherRows.value.length < 2) return

  if (voucherType.value === 'CASH_OUT') {
    // Sum debits (expenses), balance the primary cash/bank row on Credit
    const debitSum = voucherRows.value.reduce((s, r) => s + (Number(r.debit) || 0), 0)
    const primaryRow = voucherRows.value.find(r => r.accountCode === primaryAccountCode.value)
    if (primaryRow) {
      primaryRow.credit = debitSum
      primaryRow.debit = 0
      accountingStore.showToast(`Baris Kas/Bank (${primaryAccountCode.value}) disesuaikan menjadi ${formatIDR(debitSum)}`, 'info')
    } else {
      voucherRows.value[voucherRows.value.length - 1].credit = debitSum
    }
  } else {
    // Sum credits, balance the primary cash/bank row on Debit
    const creditSum = voucherRows.value.reduce((s, r) => s + (Number(r.credit) || 0), 0)
    const primaryRow = voucherRows.value.find(r => r.accountCode === primaryAccountCode.value)
    if (primaryRow) {
      primaryRow.debit = creditSum
      primaryRow.credit = 0
      accountingStore.showToast(`Baris Kas/Bank (${primaryAccountCode.value}) disesuaikan menjadi ${formatIDR(creditSum)}`, 'info')
    } else {
      voucherRows.value[0].debit = creditSum
    }
  }
}

// Indonesian Terbilang Utility
function terbilang(n: number): string {
  const angka = Math.floor(Math.abs(n))
  const bilangan = ['', 'Satu', 'Dua', 'Tiga', 'Empat', 'Lima', 'Enam', 'Tujuh', 'Delapan', 'Sembilan', 'Sepuluh', 'Sebelas']
  if (angka < 12) return bilangan[angka]
  if (angka < 20) return terbilang(angka - 10) + ' Belas'
  if (angka < 100) return terbilang(Math.floor(angka / 10)) + ' Puluh ' + terbilang(angka % 10)
  if (angka < 200) return 'Seratus ' + terbilang(angka - 100)
  if (angka < 1000) return terbilang(Math.floor(angka / 100)) + ' Ratus ' + terbilang(angka % 100)
  if (angka < 2000) return 'Seribu ' + terbilang(angka - 1000)
  if (angka < 1000000) return terbilang(Math.floor(angka / 1000)) + ' Ribu ' + terbilang(angka % 1000)
  if (angka < 1000000000) return terbilang(Math.floor(angka / 1000000)) + ' Juta ' + terbilang(angka % 1000000)
  if (angka < 1000000000000) return terbilang(Math.floor(angka / 1000000000)) + ' Miliar ' + terbilang(angka % 1000000000)
  return String(angka)
}

function handleSaveVoucher() {
  if (!isBalanced.value) {
    accountingStore.showToast('Voucher belum seimbang! Total Debit harus sama persis dengan Total Kredit.', 'error')
    return
  }

  if (!description.value.trim()) {
    accountingStore.showToast('Harap isi keterangan memo umum voucher.', 'error')
    return
  }

  for (const r of voucherRows.value) {
    if (!r.accountCode) {
      accountingStore.showToast('Harap pilih akun COA untuk semua baris.', 'error')
      return
    }
  }

  try {
    cashBankStore.recordMultiLineVoucher({
      voucherType: voucherType.value,
      primaryAccountCode: primaryAccountCode.value,
      transactionDate: transactionDate.value,
      recipientOrPayer: recipientOrPayer.value.trim() || '-',
      description: description.value.trim(),
      referenceNo: referenceNo.value.trim() || undefined,
      businessUnit: businessUnit.value,
      lines: voucherRows.value.map(r => ({
        accountCode: r.accountCode,
        accountName: r.accountName,
        debit: Number(r.debit) || 0,
        credit: Number(r.credit) || 0,
        description: r.description.trim() || description.value.trim()
      }))
    })

    router.push('/accounting/cash-bank')
  } catch (err: any) {
    accountingStore.showToast(err.message || 'Gagal menyimpan voucher', 'error')
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
          <router-link to="/accounting/cash-bank" class="hover:text-zinc-900 transition flex items-center gap-1.5 text-zinc-500">
            <Folder class="w-4 h-4 text-zinc-400" />
            <span>Buku & Transaksi</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-500">Buku Kas & Bank</span>
          <span>/</span>
          <span class="font-mono font-bold text-zinc-900">
            {{ referenceNo || (voucherType === 'CASH_OUT' ? 'BKK-BARU' : voucherType === 'CASH_IN' ? 'BKM-BARU' : 'TRF-BARU') }}
          </span>
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
            v-if="!isBalanced && balanceDiff !== 0"
            type="button"
            @click="autoBalanceRows"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold rounded-xl text-xs border border-amber-200 transition active:scale-95 cursor-pointer"
            title="Seimbangkan otomatis selisih baris kas/bank"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Seimbangkan Selisih</span>
          </button>

          <button
            type="button"
            :disabled="!isBalanced"
            @click="handleSaveVoucher"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white font-bold rounded-xl text-xs shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Send class="w-3.5 h-3.5 text-emerald-400" />
            <span>Simpan & Posting Voucher</span>
          </button>

          <router-link
            to="/accounting/cash-bank"
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
        
        <!-- Section A: Informasi Voucher Kas & Bank -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Informasi Voucher Kas & Bank
          </h3>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Jenis Voucher *</label>
              <select
                :value="voucherType"
                @change="onVoucherTypeChange(($event.target as HTMLSelectElement).value as any)"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
              >
                <option value="CASH_OUT">BKK / BBK - Kas Keluar (Pengeluaran Operasional)</option>
                <option value="CASH_IN">BKM / BBM - Kas Masuk (Penerimaan Kasir/Lainnya)</option>
                <option value="TRANSFER">TRF - Mutasi Transfer Antar Kas/Bank</option>
              </select>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Buku Kas / Rekening Utama *</label>
              <SearchableSelect
                :model-value="primaryAccountCode"
                :options="primaryAccountOptions"
                placeholder="Pilih rekening kas/bank..."
                @update:model-value="(val: any) => onPrimaryAccountChange(String(val))"
              />
              <span v-if="selectedPrimaryAccount" class="text-[10px] text-zinc-500 mt-1 block">
                Saldo: <strong class="text-zinc-950 font-mono">Rp {{ formatIDR(selectedPrimaryAccount.currentBalance || 0) }}</strong>
              </span>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Tanggal Transaksi *</label>
              <input
                v-model="transactionDate"
                type="date"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Penerima / Pembayar *</label>
              <input
                v-model="recipientOrPayer"
                type="text"
                placeholder="e.g. Toko Buku, Petugas Kasir, Vendor"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">No. Referensi / Bukti</label>
              <input
                v-model="referenceNo"
                type="text"
                placeholder="e.g. BKK-2026-001, NOTA-12"
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

            <div class="md:col-span-2">
              <label class="block text-zinc-600 mb-1 font-medium">Catatan / Keterangan Memo Transaksi *</label>
              <input
                v-model="description"
                type="text"
                placeholder="Tuliskan maksud pengeluaran / penerimaan kas/bank ini..."
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-medium focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>
          </div>
        </div>

        <!-- Section B: Detail Baris Transaksi Kas & Bank -->
        <div>
          <div class="flex items-center justify-between mb-3.5">
            <div>
              <h3 class="text-sm font-bold text-zinc-900">
                Detail Transaksi Baris Akun (Multi-Line Voucher)
              </h3>
              <p class="text-xs text-zinc-500">
                Alokasi baris debet dan kredit kas/bank berpasangan seimbang.
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
                    v-for="(row, idx) in voucherRows"
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
                      <div class="relative">
                        <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                        <input
                          v-model.number="row.debit"
                          type="number"
                          min="0"
                          step="any"
                          class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-emerald-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                          placeholder="0"
                        />
                      </div>
                    </td>

                    <td class="px-4 py-2 text-right">
                      <div class="relative">
                        <span class="absolute left-2.5 top-1.5 text-zinc-400 font-mono text-xs">Rp</span>
                        <input
                          v-model.number="row.credit"
                          type="number"
                          min="0"
                          step="any"
                          class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                          placeholder="0"
                        />
                      </div>
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
                      Total Nilai Transaksi:
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

          <!-- Terbilang Banner -->
          <div v-if="totalDebit > 0" class="mt-3 p-3 bg-zinc-50 border border-zinc-200 rounded-xl text-xs text-zinc-700 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <span class="text-zinc-400 block text-[10px] uppercase font-bold">Terbilang:</span>
              <span class="font-bold text-zinc-900 italic"># {{ terbilang(totalDebit) }} Rupiah #</span>
            </div>
            <div class="sm:text-right">
              <span class="text-zinc-400 block text-[10px] uppercase font-bold">Rekening Utama:</span>
              <span class="font-mono font-bold text-zinc-900">{{ primaryAccountCode }} - {{ selectedPrimaryAccount?.name }}</span>
            </div>
          </div>
        </div>

        <!-- Section C: Ringkasan Nilai Voucher -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Ringkasan & Keseimbangan Voucher
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
                {{ isBalanced ? 'SEIMBANG (Rp 0)' : `Selisih Rp ${formatIDR(balanceDiff)}` }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. TAB 2: RECONCILIATION & STATUS -->
      <div v-else-if="activeTab === 'RECONCILIATION'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Status Rekonsiliasi Kas & Bank</h3>
          <p class="text-xs text-zinc-500">Kesesuaian saldo buku kas dan proyeksi mutasi rekening setelah transaksi diposting.</p>
        </div>

        <div class="space-y-3 text-xs">
          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span class="font-semibold text-zinc-900">Rekening Kas/Bank Terpilih</span>
            </div>
            <span class="font-bold text-zinc-900 font-mono">[{{ primaryAccountCode }}] {{ selectedPrimaryAccount?.name }}</span>
          </div>

          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 class="w-4 h-4 text-emerald-600" />
              <span class="font-semibold text-zinc-900">Saldo Buku Saat Ini</span>
            </div>
            <span class="font-bold text-zinc-900 font-mono">Rp {{ formatIDR(selectedPrimaryAccount?.currentBalance || 0) }}</span>
          </div>

          <div class="flex items-center justify-between p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
            <div class="flex items-center gap-2.5">
              <CheckCircle2 v-if="isBalanced" class="w-4 h-4 text-emerald-600" />
              <AlertCircle v-else class="w-4 h-4 text-amber-600" />
              <span class="font-semibold text-zinc-900">Status Keseimbangan Transaksi</span>
            </div>
            <span :class="['font-bold font-mono', isBalanced ? 'text-emerald-700' : 'text-amber-700']">
              {{ isBalanced ? 'Seimbang (Rp 0)' : `Selisih Rp ${formatIDR(balanceDiff)}` }}
            </span>
          </div>
        </div>
      </div>

      <!-- 5. TAB 3: AUDIT TRAIL & LOG -->
      <div v-else-if="activeTab === 'AUDIT'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Audit Trail & Log Otorisasi</h3>
          <p class="text-xs text-zinc-500">Informasi riwayat pembuatan voucher kas dan otorisasi posting.</p>
        </div>

        <div class="space-y-3 text-xs">
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Petugas Input Kasir / Finance:</span>
            <span class="font-bold text-zinc-900">Budi Santoso, SE, Ak., CA</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Waktu Transaksi:</span>
            <span class="font-bold text-zinc-900 font-mono">{{ new Date().toLocaleString('id-ID') }}</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Tipe Otorisasi:</span>
            <span class="font-bold text-zinc-900">Otomatis Terposting ke Buku Besar (GL)</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
