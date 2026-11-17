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
  FileText
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const cashBankStore = useCashBankStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()

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
    // Default BKK: Counterpart (Expense) on Debit, Cash/Bank on Credit
    voucherRows.value = [
      {
        accountCode: '5202',
        accountName: 'BEBAN UMUM DAN ADMINISTRASI',
        description: 'Pembelian konsumsi rapat & BHP operasional',
        debit: 500000,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || 'Kas Kecil',
        description: 'Pengeluaran kas kecil operasional',
        debit: 0,
        credit: 500000,
        unit: businessUnit.value
      }
    ]
  } else if (voucherType.value === 'CASH_IN') {
    // Default BKM: Cash/Bank on Debit, Counterpart (Revenue/Receivable) on Credit
    voucherRows.value = [
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || 'Kas Kasir',
        description: 'Penerimaan kas/bank masuk',
        debit: 1500000,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: '4101',
        accountName: 'Pendapatan Jasa Medis Kidinest',
        description: 'Penerimaan konsultasi medis spesialis',
        debit: 0,
        credit: 1500000,
        unit: businessUnit.value
      }
    ]
  } else {
    // Default Transfer: Destination Bank on Debit, Source Bank/Cash on Credit
    const destBank = availableCashBankAccounts.value.find(b => b.code !== primaryAccountCode.value) || availableCashBankAccounts.value[1]
    voucherRows.value = [
      {
        accountCode: destBank?.code || '1102.001',
        accountName: destBank?.name || 'Bank BCA Operasional',
        description: 'Penerimaan setoran dana transfer',
        debit: 2000000,
        credit: 0,
        unit: businessUnit.value
      },
      {
        accountCode: primaryAccountCode.value,
        accountName: primaryAcc?.name || 'Kas Kasir',
        description: 'Penyetoran mutasi kas/bank',
        debit: 0,
        credit: 2000000,
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
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-16 max-w-6xl mx-auto">
    <!-- Header with Back Button -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <router-link
            to="/accounting/cash-bank"
            class="text-xs font-semibold text-zinc-500 hover:text-zinc-900 transition flex items-center gap-1"
          >
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>Kembali ke Buku Kas & Bank</span>
          </router-link>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Multi-Line Voucher Entry</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-1">
          Pencatatan Voucher Kas & Bank (BKM / BKK / Transfer)
        </h1>
        <p class="text-xs text-zinc-500">
          Formulir jurnal multi-baris berimbang (gaya General Ledger). Satu voucher kas keluar dapat membiayai multi-pos beban sekaligus, dan penerimaan kasir dapat memuat potongan fee perbankan.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <router-link
          to="/accounting/cash-bank"
          class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition"
        >
          Batal
        </router-link>

        <button
          type="button"
          @click="handleSaveVoucher"
          class="inline-flex items-center gap-1.5 px-5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan & Posting Voucher</span>
        </button>
      </div>
    </header>

    <!-- Top Card: Voucher Type Selector & Primary Account Header -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-5">
      <!-- Voucher Type Tabs -->
      <div>
        <label class="block font-bold text-zinc-700 text-xs mb-1.5 uppercase tracking-wider">
          Jenis Voucher Kas & Bank:
        </label>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
          <button
            type="button"
            @click="onVoucherTypeChange('CASH_OUT')"
            class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl border text-xs font-bold transition cursor-pointer"
            :class="voucherType === 'CASH_OUT' ? 'bg-rose-50 text-rose-950 border-rose-300 ring-2 ring-rose-500/20' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
          >
            <div class="w-6 h-6 rounded-lg bg-rose-100 text-rose-700 flex items-center justify-center">
              <ArrowUpRight class="w-4 h-4" />
            </div>
            <div class="text-left">
              <span class="block">BKK / BBK (Kas Keluar)</span>
              <span class="block text-[10px] font-normal text-zinc-500">Biaya operasional multi-pos, vendor, dll</span>
            </div>
          </button>

          <button
            type="button"
            @click="onVoucherTypeChange('CASH_IN')"
            class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl border text-xs font-bold transition cursor-pointer"
            :class="voucherType === 'CASH_IN' ? 'bg-emerald-50 text-emerald-950 border-emerald-300 ring-2 ring-emerald-500/20' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
          >
            <div class="w-6 h-6 rounded-lg bg-emerald-100 text-emerald-700 flex items-center justify-center">
              <ArrowDownLeft class="w-4 h-4" />
            </div>
            <div class="text-left">
              <span class="block">BKM / BBM (Kas Masuk)</span>
              <span class="block text-[10px] font-normal text-zinc-500">Pendapatan kasir, klaim asuransi, dll</span>
            </div>
          </button>

          <button
            type="button"
            @click="onVoucherTypeChange('TRANSFER')"
            class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl border text-xs font-bold transition cursor-pointer"
            :class="voucherType === 'TRANSFER' ? 'bg-zinc-950 text-white border-zinc-950 ring-2 ring-zinc-500/20' : 'bg-white hover:bg-zinc-50 text-zinc-700 border-zinc-200'"
          >
            <div class="w-6 h-6 rounded-lg bg-zinc-800 text-white flex items-center justify-center">
              <ArrowLeftRight class="w-4 h-4" />
            </div>
            <div class="text-left">
              <span class="block">TRF (Mutasi Transfer)</span>
              <span class="block text-[10px] font-normal" :class="voucherType === 'TRANSFER' ? 'text-zinc-400' : 'text-zinc-500'">Pemindahbukuan antar rekening bank/kas</span>
            </div>
          </button>
        </div>
      </div>

      <!-- Voucher Metadata Fields -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
        <div>
          <label class="block font-bold text-zinc-700 mb-1">
            Buku Kas / Rekening Bank Utama:
          </label>
          <SearchableSelect
            :model-value="primaryAccountCode"
            :options="primaryAccountOptions"
            placeholder="Cari atau pilih rekening kas/bank..."
            @update:model-value="(val: any) => onPrimaryAccountChange(String(val))"
          />
          <span v-if="selectedPrimaryAccount" class="text-[10px] text-zinc-500 mt-1 block">
            Saldo Buku: <strong class="text-zinc-950">{{ formatIDR(selectedPrimaryAccount.currentBalance || 0) }}</strong>
          </span>
        </div>

        <div>
          <label class="block font-bold text-zinc-700 mb-1">Tanggal Transaksi:</label>
          <input
            v-model="transactionDate"
            type="date"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 font-bold focus:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-950 text-xs"
          />
        </div>

        <div>
          <label class="block font-bold text-zinc-700 mb-1">
            {{ voucherType === 'CASH_IN' ? 'Diterima Dari:' : 'Dibayarkan Kepada / Penerima:' }}
          </label>
          <input
            v-model="recipientOrPayer"
            type="text"
            placeholder="Contoh: Toko Berkah / Pasien Umum"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl text-zinc-950 focus:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-950 text-xs"
          />
        </div>

        <div>
          <label class="block font-bold text-zinc-700 mb-1">No. Referensi / Bukti:</label>
          <input
            v-model="referenceNo"
            type="text"
            placeholder="Contoh: NOTA-1102 atau TRF-BCA-992"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-mono text-zinc-950 focus:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-950 text-xs"
          />
        </div>
      </div>

      <!-- Keterangan Umum & Unit Bisnis -->
      <div class="grid grid-cols-1 sm:grid-cols-4 gap-4 text-xs pt-1">
        <div class="sm:col-span-1">
          <label class="block font-bold text-zinc-700 mb-1">Unit Bisnis Klinik:</label>
          <select
            v-model="businessUnit"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl font-bold text-zinc-950 focus:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-950 cursor-pointer text-xs"
          >
            <option value="KIDINEST">Kidinest</option>
            <option value="SSC">Semarang Skin Center</option>
            <option value="ALL">SMG (Seluruh Unit)</option>
          </select>
        </div>

        <div class="sm:col-span-3">
          <label class="block font-bold text-zinc-700 mb-1">Keterangan Memo Umum Voucher:</label>
          <input
            v-model="description"
            type="text"
            placeholder="Contoh: Pengeluaran operasional pembelian konsumsi dokter poli dan ATK medis"
            class="w-full px-3 py-2 bg-zinc-50 border border-zinc-300 rounded-xl text-zinc-950 focus:bg-white focus:outline-none focus:ring-2 focus:ring-zinc-950 text-xs font-medium"
          />
        </div>
      </div>
    </div>

    <!-- SPREADSHEET-STYLE MULTI-LINE VOUCHER ENTRY (Gaya Jurnal Umum Brayo) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <FileText class="w-4 h-4 text-zinc-800" />
            <h3 class="text-sm font-bold text-zinc-900">
              Rincian Alokasi Akun Rekening (Multi-Line Double Entry)
            </h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Tambahkan baris sebanyak yang diperlukan. Pastikan total debit seimbang dengan total kredit.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="autoBalanceRows"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-100 hover:bg-zinc-200 text-zinc-800 rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer"
            title="Otomatis seimbangkan baris kredit/debit dengan akun kas/bank"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Otomatis Seimbangkan</span>
          </button>

          <button
            type="button"
            @click="addRow"
            class="inline-flex items-center gap-1 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition active:scale-95 cursor-pointer shadow-xs"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>+ Tambah Baris</span>
          </button>
        </div>
      </div>

      <!-- Spreadsheet Table -->
      <div class="overflow-x-auto border border-zinc-200 rounded-xl">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="py-2.5 px-3 w-10 text-center">#</th>
              <th class="py-2.5 px-3 w-32">Unit</th>
              <th class="py-2.5 px-3 w-72">Akun Rekening (COA)</th>
              <th class="py-2.5 px-3">Uraian Transaksi Baris</th>
              <th class="py-2.5 px-3 w-40 text-right">Debit (Rp)</th>
              <th class="py-2.5 px-3 w-40 text-right">Kredit (Rp)</th>
              <th class="py-2.5 px-3 w-12 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="(row, idx) in voucherRows"
              :key="idx"
              class="hover:bg-zinc-50/70"
            >
              <td class="py-2 px-3 text-center text-zinc-400 font-mono font-bold">{{ idx + 1 }}</td>
              <td class="py-2 px-3">
                <select
                  v-model="row.unit"
                  class="w-full px-2 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
                >
                  <option value="KIDINEST">Kidinest</option>
                  <option value="SSC">SSC</option>
                  <option value="ALL">ALL</option>
                </select>
              </td>
              <td class="py-2 px-3 min-w-[260px]">
                <SearchableSelect
                  :model-value="row.accountCode"
                  :options="accountSelectOptions"
                  placeholder="Pilih atau cari akun COA..."
                  @update:model-value="(val: any) => onAccountSelect(idx, String(val))"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model="row.description"
                  type="text"
                  placeholder="Keterangan alokasi baris"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs text-zinc-800 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model.number="row.debit"
                  type="number"
                  min="0"
                  step="1000"
                  placeholder="0"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs tabular-nums font-bold text-right text-zinc-950 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3">
                <input
                  v-model.number="row.credit"
                  type="number"
                  min="0"
                  step="1000"
                  placeholder="0"
                  class="w-full px-2.5 py-1.5 bg-white border border-zinc-200 rounded-xl text-xs tabular-nums font-bold text-right text-zinc-950 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
                />
              </td>
              <td class="py-2 px-3 text-center">
                <button
                  type="button"
                  @click="removeRow(idx)"
                  class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition cursor-pointer"
                  title="Hapus baris ini"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          </tbody>
          <!-- Totals Footer -->
          <tfoot>
            <tr class="bg-zinc-50 border-t border-zinc-200 font-bold text-zinc-900">
              <td colspan="4" class="py-3 px-4">
                <div class="flex items-center justify-between">
                  <span class="text-xs uppercase tracking-wider text-zinc-500 font-bold">Total Nilai Voucher:</span>
                  <span v-if="isBalanced" class="text-xs text-emerald-700 font-bold flex items-center gap-1">
                    <CheckCircle2 class="w-4 h-4 text-emerald-600" />
                    <span>Seimbang (Balanced)</span>
                  </span>
                  <span v-else class="text-xs text-rose-600 font-bold flex items-center gap-1">
                    <AlertCircle class="w-4 h-4 text-rose-600" />
                    <span>Selisih: {{ formatIDR(balanceDiff) }}</span>
                  </span>
                </div>
              </td>
              <td class="py-3 px-3 text-right tabular-nums font-extrabold text-zinc-950 bg-zinc-100/60 text-sm">
                {{ formatIDR(totalDebit) }}
              </td>
              <td class="py-3 px-3 text-right tabular-nums font-extrabold text-zinc-950 bg-zinc-100/60 text-sm">
                {{ formatIDR(totalCredit) }}
              </td>
              <td></td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Terbilang Display -->
      <div v-if="totalDebit > 0" class="p-3 bg-zinc-50 border border-zinc-200 rounded-xl text-xs text-zinc-700 flex items-center justify-between">
        <div>
          <span class="text-zinc-400 block text-[10px] uppercase font-bold">Terbilang:</span>
          <span class="font-bold text-zinc-900 italic"># {{ terbilang(totalDebit) }} Rupiah #</span>
        </div>
        <div class="text-right">
          <span class="text-zinc-400 block text-[10px] uppercase font-bold">Akun Kas/Bank Utama:</span>
          <span class="font-mono font-bold text-zinc-900">{{ primaryAccountCode }} - {{ selectedPrimaryAccount?.name }}</span>
        </div>
      </div>
    </div>

    <!-- Bottom Actions Bar -->
    <div class="flex items-center justify-between pt-4 border-t border-zinc-200">
      <router-link
        to="/accounting/cash-bank"
        class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition"
      >
        Batal & Kembali
      </router-link>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="autoBalanceRows"
          class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-semibold border border-zinc-200 transition cursor-pointer"
        >
          Otomatis Seimbangkan
        </button>

        <button
          type="button"
          @click="handleSaveVoucher"
          class="inline-flex items-center gap-1.5 px-6 py-2.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-black shadow-xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan & Posting Voucher Kas/Bank</span>
        </button>
      </div>
    </div>
  </div>
</template>
