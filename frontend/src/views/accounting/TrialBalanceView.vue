<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel } from '@/utils/export'
import type { BusinessUnit, AccountType } from '@/types/accounting'
import {
  Scale,
  Download,
  CheckCircle2,
  AlertCircle,
  Search,
  Filter,
  Calendar,
  Building2,
  ExternalLink,
  Printer
} from 'lucide-vue-next'

const router = useRouter()
const ledgerStore = useLedgerStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await ledgerStore.fetchAccounts()
})

// Filter State
const startDate = ref<string>('2026-01-01')
const endDate = ref<string>(new Date().toISOString().slice(0, 10))
const selectedUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const selectedType = ref<string>('ALL')
const searchQuery = ref<string>('')
const hideZeroBalances = ref<boolean>(false)

// Trial Balance Data Calculation per Account
interface TrialBalanceCalculatedRow {
  code: string
  name: string
  type: AccountType
  normalBalance: 'DEBIT' | 'CREDIT'
  unit: BusinessUnit
  openingDebit: number
  openingCredit: number
  periodDebit: number
  periodCredit: number
  endingDebit: number
  endingCredit: number
}

const calculatedRows = computed<TrialBalanceCalculatedRow[]>(() => {
  const leafAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.isActive !== false)

  return leafAccounts.map(acc => {
    // 1. Initial / Opening balance
    const cur = acc.currentBalance || 0
    let openingDebit = 0
    let openingCredit = 0
    if (acc.normalBalance === 'DEBIT') {
      openingDebit = cur > 0 ? cur : 0
      openingCredit = cur < 0 ? Math.abs(cur) : 0
    } else {
      openingCredit = cur > 0 ? cur : 0
      openingDebit = cur < 0 ? Math.abs(cur) : 0
    }

    // 2. Period mutations from general ledger entries
    const entries = ledgerStore.generalLedgerEntries.filter(e => e.accountCode === acc.code)
    const periodDebit = entries.reduce((s, e) => s + (e.debit || 0), 0)
    const periodCredit = entries.reduce((s, e) => s + (e.credit || 0), 0)

    // 3. Ending balance
    let endingDebit = 0
    let endingCredit = 0

    if (acc.normalBalance === 'DEBIT') {
      const net = openingDebit - openingCredit + periodDebit - periodCredit
      if (net >= 0) endingDebit = net
      else endingCredit = Math.abs(net)
    } else {
      const net = openingCredit - openingDebit + periodCredit - periodDebit
      if (net >= 0) endingCredit = net
      else endingDebit = Math.abs(net)
    }

    return {
      code: acc.code,
      name: acc.name,
      type: acc.type,
      normalBalance: acc.normalBalance,
      unit: acc.businessUnit || 'ALL',
      openingDebit,
      openingCredit,
      periodDebit,
      periodCredit,
      endingDebit,
      endingCredit
    }
  })
})

// Filtered Rows
const filteredRows = computed(() => {
  return calculatedRows.value.filter(r => {
    if (selectedUnit.value !== 'ALL' && r.unit !== 'ALL' && r.unit !== selectedUnit.value) {
      return false
    }
    if (selectedType.value !== 'ALL' && r.type !== selectedType.value) {
      return false
    }
    if (hideZeroBalances.value) {
      const isZero =
        r.openingDebit === 0 &&
        r.openingCredit === 0 &&
        r.periodDebit === 0 &&
        r.periodCredit === 0 &&
        r.endingDebit === 0 &&
        r.endingCredit === 0
      if (isZero) return false
    }
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase()
      return r.code.toLowerCase().includes(q) || r.name.toLowerCase().includes(q)
    }
    return true
  })
})

// Totals
const totals = computed(() => {
  return calculatedRows.value.reduce(
    (acc, r) => {
      acc.openingDebit += r.openingDebit
      acc.openingCredit += r.openingCredit
      acc.periodDebit += r.periodDebit
      acc.periodCredit += r.periodCredit
      acc.endingDebit += r.endingDebit
      acc.endingCredit += r.endingCredit
      return acc
    },
    {
      openingDebit: 0,
      openingCredit: 0,
      periodDebit: 0,
      periodCredit: 0,
      endingDebit: 0,
      endingCredit: 0
    }
  )
})

const isBalanced = computed(() => {
  const openingDiff = Math.abs(totals.value.openingDebit - totals.value.openingCredit)
  const periodDiff = Math.abs(totals.value.periodDebit - totals.value.periodCredit)
  const endingDiff = Math.abs(totals.value.endingDebit - totals.value.endingCredit)
  return openingDiff < 1 && periodDiff < 1 && endingDiff < 1
})

function drillDownToLedger(accountCode: string) {
  router.push(`/accounting/ledger?account=${accountCode}`)
}

function handleExportExcel() {
  const data = filteredRows.value.map(r => ({
    'Kode Akun': r.code,
    'Nama Akun': r.name,
    'Kategori': r.type,
    'Saldo Normal': r.normalBalance,
    'Unit': r.unit,
    'Saldo Awal Debet': r.openingDebit,
    'Saldo Awal Kredit': r.openingCredit,
    'Mutasi Periode Debet': r.periodDebit,
    'Mutasi Periode Kredit': r.periodCredit,
    'Saldo Akhir Debet': r.endingDebit,
    'Saldo Akhir Kredit': r.endingCredit
  }))

  exportToExcel(data, `Neraca_Saldo_Trial_Balance_BCS_${selectedUnit.value}`, 'Trial Balance')
  accountingStore.showToast('Neraca Saldo berhasil diekspor ke Excel (.xlsx)', 'success')
}

function handlePrint() {
  window.print()
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-24">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Laporan Keuangan & Kontrol
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Format 6-Kolom Standar Brayo</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Neraca Saldo (Trial Balance)
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Verifikasi keseimbangan pembukuan dari Saldo Awal, Mutasi Periode Berjalan, hingga Saldo Akhir seluruh akun COA.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>

        <button
          @click="handlePrint"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5 text-zinc-600" />
          <span>Cetak</span>
        </button>
      </div>
    </header>

    <!-- Filter & Control Bar -->
    <div class="bg-white p-4 rounded-2xl border border-zinc-200/90 shadow-2xs grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-3 text-xs">
      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Dari Tanggal:</span>
        </label>
        <input
          v-model="startDate"
          type="date"
          class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-medium text-zinc-900 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-zinc-500" />
          <span>Sampai Tanggal:</span>
        </label>
        <input
          v-model="endDate"
          type="date"
          class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-medium text-zinc-900 focus:outline-none"
        />
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Building2 class="w-3.5 h-3.5 text-zinc-500" />
          <span>Unit Bisnis:</span>
        </label>
        <select
          v-model="selectedUnit"
          class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Semua Unit (Konsolidasi)</option>
          <option value="KIDINEST">Unit Kidinest (Poli Anak)</option>
          <option value="SSC">Unit Semarang Skin Center (SSC)</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Filter class="w-3.5 h-3.5 text-zinc-500" />
          <span>Kategori Akun:</span>
        </label>
        <select
          v-model="selectedType"
          class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:outline-none cursor-pointer"
        >
          <option value="ALL">Semua Kategori (1 - 7)</option>
          <option value="ASSET">1 - Aset (Aktiva)</option>
          <option value="LIABILITY">2 - Kewajiban (Hutang)</option>
          <option value="EQUITY">3 - Ekuitas (Modal)</option>
          <option value="REVENUE">4 - Pendapatan</option>
          <option value="COGS">5 - Beban Pokok (HPP)</option>
          <option value="EXPENSE">6/7 - Beban Operasional</option>
        </select>
      </div>

      <div>
        <label class="block font-bold text-zinc-700 mb-1 flex items-center gap-1">
          <Search class="w-3.5 h-3.5 text-zinc-500" />
          <span>Cari Akun:</span>
        </label>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Kode / nama akun..."
          class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:outline-none"
        />
      </div>
    </div>

    <!-- Toggle Zero Balance & Balance Verification Status Banner -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <label class="inline-flex items-center gap-2 cursor-pointer text-xs font-bold text-zinc-700 select-none">
        <input
          v-model="hideZeroBalances"
          type="checkbox"
          class="w-4 h-4 rounded text-zinc-950 focus:ring-0 cursor-pointer"
        />
        <span>Sembunyikan akun bersaldo nol (Hide Zero Balances)</span>
      </label>

      <!-- Status Badge -->
      <div
        :class="[
          'px-4 py-2 rounded-xl border flex items-center gap-2.5 text-xs font-bold shadow-2xs',
          isBalanced
            ? 'bg-emerald-50 border-emerald-200 text-emerald-900'
            : 'bg-rose-50 border-rose-200 text-rose-900'
        ]"
      >
        <CheckCircle2 v-if="isBalanced" class="w-4 h-4 text-emerald-600" />
        <AlertCircle v-else class="w-4 h-4 text-rose-600" />
        <span>
          {{ isBalanced ? '✓ NERACA SALDO SEIMBANG (BALANCED)' : '⚠ TERDAPAT SELISIH KESEIMBANGAN' }}
        </span>
        <span class="font-mono text-zinc-600 font-normal">
          (Saldo Akhir: Rp {{ formatIDR(totals.endingDebit) }})
        </span>
      </div>
    </div>

    <!-- 6-Column Standard Table (Format Brayo & PSAK) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-100/80 border-b border-zinc-300 text-[11px] font-bold text-zinc-700 uppercase tracking-wider">
              <th rowspan="2" class="px-4 py-2.5 border-r border-zinc-200 w-32">Kode Akun</th>
              <th rowspan="2" class="px-4 py-2.5 border-r border-zinc-200 min-w-[240px]">Nama Akun (COA)</th>
              <th colspan="2" class="px-3 py-2 text-center border-r border-zinc-200 bg-zinc-50 font-black">
                1. Saldo Awal
              </th>
              <th colspan="2" class="px-3 py-2 text-center border-r border-zinc-200 bg-blue-50/60 font-black text-blue-950">
                2. Mutasi Periode
              </th>
              <th colspan="2" class="px-3 py-2 text-center bg-emerald-50/60 font-black text-emerald-950">
                3. Saldo Akhir
              </th>
            </tr>
            <tr class="bg-zinc-50 border-b border-zinc-300 text-[10px] font-bold text-zinc-600 uppercase">
              <th class="px-3 py-1.5 text-right w-32 border-r border-zinc-200">Debet (Rp)</th>
              <th class="px-3 py-1.5 text-right w-32 border-r border-zinc-200">Kredit (Rp)</th>
              <th class="px-3 py-1.5 text-right w-32 border-r border-zinc-200 text-emerald-800">Debet (Rp)</th>
              <th class="px-3 py-1.5 text-right w-32 border-r border-zinc-200 text-rose-800">Kredit (Rp)</th>
              <th class="px-3 py-1.5 text-right w-36 border-r border-zinc-200 text-zinc-950 font-black">Debet (Rp)</th>
              <th class="px-3 py-1.5 text-right w-36 text-zinc-950 font-black">Kredit (Rp)</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="row in filteredRows"
              :key="row.code"
              class="hover:bg-zinc-50/80 transition-colors group cursor-pointer"
              @click="drillDownToLedger(row.code)"
              title="Klik untuk membuka Buku Besar akun ini"
            >
              <!-- Account Code -->
              <td class="px-4 py-2.5 font-mono font-bold text-zinc-900 border-r border-zinc-100 flex items-center justify-between">
                <span>{{ row.code }}</span>
                <ExternalLink class="w-3 h-3 text-zinc-400 opacity-0 group-hover:opacity-100 transition" />
              </td>

              <!-- Account Name & Type -->
              <td class="px-4 py-2.5 font-medium text-zinc-900 border-r border-zinc-100">
                <span>{{ row.name }}</span>
                <span class="text-[10px] text-zinc-400 block font-normal">
                  {{ row.type }} • Normal: {{ row.normalBalance }}
                </span>
              </td>

              <!-- Opening Balance -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-zinc-100 text-zinc-700">
                {{ row.openingDebit > 0 ? formatIDR(row.openingDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono border-r border-zinc-100 text-zinc-700">
                {{ row.openingCredit > 0 ? formatIDR(row.openingCredit, false) : '-' }}
              </td>

              <!-- Period Movements -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-zinc-100 text-emerald-700 font-semibold">
                {{ row.periodDebit > 0 ? formatIDR(row.periodDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono border-r border-zinc-100 text-rose-700 font-semibold">
                {{ row.periodCredit > 0 ? formatIDR(row.periodCredit, false) : '-' }}
              </td>

              <!-- Ending Balance -->
              <td class="px-3 py-2.5 text-right font-mono border-r border-zinc-100 font-extrabold text-zinc-950 bg-zinc-50/40">
                {{ row.endingDebit > 0 ? formatIDR(row.endingDebit, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right font-mono font-extrabold text-zinc-950 bg-zinc-50/40">
                {{ row.endingCredit > 0 ? formatIDR(row.endingCredit, false) : '-' }}
              </td>
            </tr>
          </tbody>

          <!-- Totals Footer -->
          <tfoot class="bg-zinc-100 font-bold border-t-2 border-zinc-300 text-xs">
            <tr>
              <td colspan="2" class="px-4 py-3 text-right text-zinc-900 uppercase tracking-wider font-black">
                TOTAL KESEIMBANGAN:
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-zinc-200">
                {{ formatIDR(totals.openingDebit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-zinc-200">
                {{ formatIDR(totals.openingCredit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-zinc-200 text-emerald-800">
                {{ formatIDR(totals.periodDebit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono border-r border-zinc-200 text-rose-800">
                {{ formatIDR(totals.periodCredit, false) }}
              </td>
              <td class="px-3 py-3 text-right font-mono font-black text-sm border-r border-zinc-200 text-zinc-950">
                Rp {{ formatIDR(totals.endingDebit) }}
              </td>
              <td class="px-3 py-3 text-right font-mono font-black text-sm text-zinc-950">
                Rp {{ formatIDR(totals.endingCredit) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>
