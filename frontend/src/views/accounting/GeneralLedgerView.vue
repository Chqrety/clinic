<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import { exportToExcel, printReport } from '@/utils/export'
import JournalDetailModal from '@/components/journal/JournalDetailModal.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import type { JournalEntry } from '@/types/accounting'
import {
  FileText,
  Search,
  Download,
  Printer,
  Calendar,
  Eye,
  Building,
  Scale,
  ArrowLeft,
  BookOpen
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const accountingStore = useAccountingStore()

// Account alias mapping from reports (e.g. 11, 12, 13, 2101, etc.) to standard COA codes
const accountAliases: Record<string, string> = {
  '11': '1-1100',
  '1101': '1-1111',
  '1102': '1-1112',
  '1103': '1-1113',
  '1104': '1-1121',
  '1105': '1-1122',
  '1106': '1-1130',
  '12': '1-1200',
  '1201': '1-1201',
  '1202': '1-1202',
  '13': '1-1300',
  '1301': '1-1310',
  '1302': '1-1330',
  '1303': '1-1320',
  '1304': '1-1311',
  '14': '1-1400',
  '15': '1-1500',
  '16': '1-2000',
  '1601': '1-2010',
  '1602': '1-2200',
  '1603': '1-2100',
  '2101': '2-1410',
  '2102': '2-1420',
  '2103': '2-1430',
  '22': '2-1110',
  '23': '2-1000',
  '24': '2-1000',
  '25': '2-1501',
  '26': '2-1000',
  '31': '3-1100',
  '32': '3-1200',
  '41': '4-1000',
  '4101': '4-1110',
  '4102': '4-1120',
  '4103': '4-1130',
  '4104': '4-1140',
  '4106': '4-2110',
  '4107': '4-2120',
  '4108': '4-2130',
  '4109': '4-2140',
  '51': '5-0000',
  '5101': '5-1101',
  '5102': '5-1102',
  '5103': '5-1103',
  '5105': '5-1105',
  '5106': '5-1106',
  '5107': '5-1107',
  '5108': '5-1108',
  '52': '6-0000',
  '5201': '6-1101',
  '5202': '6-1101',
  '6101': '4-9200',
  '7103': '6-1201',
  '7104': '2-1504',
  '7107': '6-1201'
}

// Read query param if passed from BalanceSheet, ProfitLoss, or CashFlow (e.g. ?account=1-1111)
const initialAccount = (route.query.account as string) || '1-1111'
const resolvedInitial = accountAliases[initialAccount] || initialAccount
const selectedAccountCode = ref(resolvedInitial)
const searchQuery = ref('')
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')

watch(
  () => route.query.account,
  (newAcc) => {
    if (newAcc) {
      const codeStr = newAcc as string
      selectedAccountCode.value = accountAliases[codeStr] || codeStr
    }
  }
)

const selectedJournal = ref<JournalEntry | null>(null)
const showDetailModal = ref(false)

function inspectJournalByNumber(jNumber: string) {
  const j = journalsStore.journals.find(item => item.journalNumber === jNumber)
  if (j) {
    selectedJournal.value = j
    showDetailModal.value = true
    return
  }

  // If not found in store, create a dynamic balanced double-entry voucher
  const acc = activeAccount.value
  const entry = accountLedgerEntries.value.find(e => e.journalNumber === jNumber)
  if (acc && entry) {
    const isDebit = entry.debit > 0
    const amt = isDebit ? entry.debit : entry.credit

    // Corresponding balancing account
    const balancingCode = acc.type === 'ASSET' ? '4-1110' : acc.type === 'REVENUE' ? '1-1111' : acc.type === 'EXPENSE' || acc.type === 'COGS' ? '1-1121' : '1-1122'
    const balancingName = acc.type === 'ASSET' ? 'Pendapatan Layanan Medis' : acc.type === 'REVENUE' ? 'Kas Kasir Operasional' : 'Bank BCA Operasional'

    selectedJournal.value = {
      id: `dyn-${jNumber}`,
      journalNumber: jNumber,
      journalDate: entry.date,
      journalSource: entry.source as any,
      sourceReference: `REF-${jNumber.slice(-4)}`,
      businessUnit: (entry.businessUnit as any) || 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: entry.branchName,
      description: entry.description,
      totalDebit: amt,
      totalCredit: amt,
      status: 'POSTED',
      lines: [
        {
          id: `line-${jNumber}-1`,
          accountId: acc.id,
          accountCode: acc.code,
          accountName: acc.name,
          debit: entry.debit,
          credit: entry.credit,
          description: entry.description
        },
        {
          id: `line-${jNumber}-2`,
          accountId: `acc-bal-${balancingCode}`,
          accountCode: balancingCode,
          accountName: balancingName,
          debit: entry.credit,
          credit: entry.debit,
          description: `Lawas Akun Transaksi: ${entry.description}`
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: `${entry.date} 09:00:00`,
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: `${entry.date} 10:15:00`,
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: `${entry.date} 11:00:00`,
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: `${entry.date} 11:05:00`,
      isManual: false
    }
    showDetailModal.value = true
  }
}

const nonHeaderAccounts = computed(() => {
  return ledgerStore.accounts.filter(a => !a.isHeader)
})

const allAccountsGrouped = computed(() => {
  return ledgerStore.accounts
})

const activeAccount = computed(() => {
  const code = selectedAccountCode.value
  const target = accountAliases[code] || code
  const found = ledgerStore.accounts.find(a => a.code === target || a.code === code)
  return found || nonHeaderAccounts.value[0]
})

// Comprehensive Subsidiary Ledger Items for "Keseluruhan Transaksi"
interface SubsidiaryLedgerItem {
  id: string
  date: string
  journalNumber: string
  source: string
  description: string
  businessUnit: string
  branchName: string
  debit: number
  credit: number
  runningBalance: number
}

// Generate realistic comprehensive clinical transaction history for November 2026
function generateCompleteLedgerForAccount(acc: typeof ledgerStore.accounts[0]): SubsidiaryLedgerItem[] {
  const isDebit = acc.normalBalance === 'DEBIT'
  const baseBal = acc.currentBalance || 25000000
  let currentRunning = baseBal * 0.55

  // Account-specific comprehensive transaction templates
  const items: SubsidiaryLedgerItem[] = []

  const clinicDays = [
    { day: '01', j: '001', desc: `Penerimaan kasir awal bulan & posting reguler ${acc.name}`, mult: 0.08, dr: isDebit, unit: 'KIDINEST', src: 'CASHIER_CLOSING' },
    { day: '03', j: '004', desc: `Pemeriksaan dokter spesialis anak & vaksinasi pasien Kidinest`, mult: 0.06, dr: isDebit, unit: 'KIDINEST', src: 'OPERATIONAL_BILLING' },
    { day: '05', j: '008', desc: `Tindakan laser Nd:YAG Q-Switch & facial treatment estetika SSC`, mult: 0.10, dr: isDebit, unit: 'SSC', src: 'OPERATIONAL_BILLING' },
    { day: '07', j: '012', desc: `Pembelian perbekalan farmasi, vaksin & krim distributor PBF`, mult: 0.07, dr: !isDebit, unit: 'KIDINEST', src: 'PHARMACY_PURCHASE' },
    { day: '09', j: '015', desc: `Settlement EDC & pembayaran QRIS harian pasien`, mult: 0.09, dr: isDebit, unit: 'ALL', src: 'SETTLEMENT' },
    { day: '11', j: '018', desc: `Tagihan klaim asuransi kesehatan Sinarmas & Prudential`, mult: 0.08, dr: isDebit, unit: 'KIDINEST', src: 'OPERATIONAL_BILLING' },
    { day: '13', j: '022', desc: `Pemakaian bahan medis habis pakai (BMHP) & obat dispensing`, mult: 0.05, dr: !isDebit, unit: 'KIDINEST', src: 'OPERATIONAL_BILLING' },
    { day: '15', j: '025', desc: `Penerimaan deposit uang muka paket perawatan 10x Facial Glow`, mult: 0.06, dr: isDebit, unit: 'SSC', src: 'CASHIER_CLOSING' },
    { day: '17', j: '029', desc: `Pencairan klaim asuransi rawat jalan via transfer bank`, mult: 0.09, dr: isDebit, unit: 'KIDINEST', src: 'SETTLEMENT' },
    { day: '19', j: '032', desc: `Pengeluaran operasional petty cash & logistik klinik harian`, mult: 0.03, dr: !isDebit, unit: 'ALL', src: 'EXPENSE' },
    { day: '21', j: '035', desc: `Pengolahan & incinerasi limbah medis B3 oleh PT Medivest`, mult: 0.04, dr: !isDebit, unit: 'ALL', src: 'EXPENSE' },
    { day: '23', j: '038', desc: `Pelunasan hutang tempo distributor supplier PBF via Giro`, mult: 0.07, dr: !isDebit, unit: 'KIDINEST', src: 'PAYMENT' },
    { day: '25', j: '041', desc: `Akrual bagi hasil dokter spesialis (Sp.A & Sp.DVE) & PPh 21 Dokter`, mult: 0.12, dr: !isDebit, unit: 'ALL', src: 'DOCTOR_FEE' },
    { day: '27', j: '044', desc: `Pencairan honor dokter spesialis via transfer payroll Bank Mandiri`, mult: 0.10, dr: isDebit, unit: 'ALL', src: 'DOCTOR_FEE' },
    { day: '28', j: '047', desc: `Pembayaran gaji tenaga kesehatan, perawat & staf administrasi`, mult: 0.08, dr: !isDebit, unit: 'ALL', src: 'PAYROLL' },
    { day: '29', j: '050', desc: `Pembayaran rekening listrik PLN daya medis & internet WiFi ruko`, mult: 0.03, dr: !isDebit, unit: 'ALL', src: 'EXPENSE' },
    { day: '30', j: '053', desc: `Alokasi penyusutan aset tetap & peralatan medis laser/USG 4D`, mult: 0.04, dr: !isDebit, unit: 'ALL', src: 'DEPRECIATION' },
    { day: '30', j: '056', desc: `Rekonsiliasi tutup buku bulanan & selisih kas kasir akhir periode`, mult: 0.05, dr: isDebit, unit: 'ALL', src: 'CASHIER_CLOSING' }
  ]

  for (const step of clinicDays) {
    const amt = Math.round(baseBal * step.mult)
    const debit = step.dr ? amt : 0
    const credit = !step.dr ? amt : 0

    if (isDebit) {
      currentRunning += (debit - credit)
    } else {
      currentRunning += (credit - debit)
    }

    items.push({
      id: `gl-${acc.code}-${step.day}`,
      date: `2026-11-${step.day}`,
      journalNumber: `JV/2026/11/${step.j}`,
      source: step.src,
      description: step.desc,
      businessUnit: step.unit,
      branchName: step.unit === 'SSC' ? 'Cabang SSC Estetika' : 'Kebayoran Baru (Pusat)',
      debit,
      credit,
      runningBalance: currentRunning
    })
  }

  return items
}

const accountLedgerEntries = computed<SubsidiaryLedgerItem[]>(() => {
  const acc = activeAccount.value
  if (!acc) return []

  const isDebit = acc.normalBalance === 'DEBIT'

  // Match real entries in store
  const targetCode = acc.code
  const realEntries = ledgerStore.generalLedgerEntries.filter(
    e => e.accountCode === targetCode || (acc.isHeader && e.accountCode.startsWith(targetCode.slice(0, 3)))
  )

  if (realEntries.length >= 8) {
    let bal = (acc.currentBalance || 20000000) * 0.55
    return realEntries.map(e => {
      if (isDebit) {
        bal += (e.debit - e.credit)
      } else {
        bal += (e.credit - e.debit)
      }
      return {
        id: e.id,
        date: e.date,
        journalNumber: e.journalNumber,
        source: e.source,
        description: e.description,
        businessUnit: e.businessUnit,
        branchName: e.branchName,
        debit: e.debit,
        credit: e.credit,
        runningBalance: bal
      }
    })
  }

  // Generate full complete transactions for the month of November 2026
  return generateCompleteLedgerForAccount(acc)
})

const filteredLedgerEntries = computed(() => {
  if (!searchQuery.value.trim()) return accountLedgerEntries.value
  const q = searchQuery.value.toLowerCase().trim()
  return accountLedgerEntries.value.filter(e =>
    e.journalNumber.toLowerCase().includes(q) ||
    e.description.toLowerCase().includes(q) ||
    e.source.toLowerCase().includes(q) ||
    e.date.includes(q)
  )
})

const totalDebit = computed(() => filteredLedgerEntries.value.reduce((s, e) => s + e.debit, 0))
const totalCredit = computed(() => filteredLedgerEntries.value.reduce((s, e) => s + e.credit, 0))
const endingBalance = computed(() => {
  if (filteredLedgerEntries.value.length === 0) return activeAccount.value?.currentBalance || 0
  return filteredLedgerEntries.value[filteredLedgerEntries.value.length - 1].runningBalance
})

function handleExportExcel() {
  const acc = activeAccount.value
  const data = filteredLedgerEntries.value.map((e, idx) => ({
    'No': idx + 1,
    'Tanggal': e.date,
    'No Jurnal': e.journalNumber,
    'Sumber': e.source,
    'Uraian': e.description,
    'Unit': e.businessUnit,
    'Cabang': e.branchName,
    'Debit (IDR)': e.debit,
    'Kredit (IDR)': e.credit,
    'Saldo Berjalan (IDR)': e.runningBalance
  }))

  exportToExcel(data, `Buku_Besar_${acc?.code}_${acc?.name}_Nov2026`, 'Buku Besar')
  accountingStore.showToast(`Buku Besar Akun ${acc?.code} berhasil diekspor ke Excel`, 'success')
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header matching Image 2 & 4 with Dark ERP Look -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-4 text-white flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-xl">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-slate-400">
          <span class="hover:text-white cursor-pointer font-medium" @click="router.push('/reports/balance-sheet')">
            Laporan Neraca
          </span>
          <span>/</span>
          <span class="hover:text-white cursor-pointer font-medium" @click="router.push('/reports/profit-loss')">
            Laba Rugi
          </span>
          <span>/</span>
          <span class="text-blue-400 font-semibold">Buku Besar (General Ledger)</span>
        </div>
        <h1 class="text-base font-bold text-white tracking-wide flex items-center gap-2">
          <span>Rincian Buku Besar Akun: {{ activeAccount?.code }} - {{ activeAccount?.name }}</span>
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-3 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-[#1E293B] px-3 py-1.5 rounded-lg border border-slate-700">
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
          <span class="text-slate-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-white w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-xs transition-colors"
        >
          <Download class="w-3.5 h-3.5" />
          <span>Ekspor</span>
        </button>

        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white rounded-lg font-medium transition-colors"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Account Picker & Search Toolbar -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl p-3.5 flex flex-col md:flex-row items-center justify-between gap-3 text-xs text-slate-300">
      <div class="flex flex-wrap items-center gap-3 w-full md:w-auto">
        <span class="font-bold text-slate-400">Pilih Akun Buku Besar:</span>
        <select
          v-model="selectedAccountCode"
          class="px-3 py-1.5 bg-[#1E293B] border border-slate-700 rounded-lg font-bold text-xs text-white focus:ring-2 focus:ring-blue-500 focus:outline-none min-w-[340px] cursor-pointer"
        >
          <option
            v-for="acc in allAccountsGrouped"
            :key="acc.code"
            :value="acc.code"
          >
            {{ acc.code }} - {{ acc.name }} ({{ acc.normalBalance }})
          </option>
        </select>

        <!-- Search in Ledger Filter -->
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari transaksi / no jurnal..."
            class="pl-7 pr-3 py-1.5 bg-[#1E293B] border border-slate-700 rounded-lg text-xs text-white placeholder-slate-500 focus:ring-2 focus:ring-blue-500 focus:outline-none w-56"
          />
          <Search class="w-3.5 h-3.5 text-slate-400 absolute left-2 top-2.5 pointer-events-none" />
        </div>
      </div>

      <!-- Quick Metrics -->
      <div class="flex items-center gap-4 font-mono text-xs">
        <div>
          <span class="text-slate-500 block text-[10px] uppercase font-bold">Saldo Normal</span>
          <span class="font-bold text-blue-400">{{ activeAccount?.normalBalance }}</span>
        </div>
        <div class="h-6 w-px bg-slate-800"></div>
        <div>
          <span class="text-slate-500 block text-[10px] uppercase font-bold">Saldo Akhir</span>
          <span class="font-bold text-emerald-400 text-sm">{{ formatIDR(endingBalance) }}</span>
        </div>
      </div>
    </div>

    <!-- Ledger Table (High-Contrast Dark Look matching Images 2, 3, 4) -->
    <div class="bg-[#121826] border border-slate-800 rounded-xl overflow-hidden shadow-2xl text-slate-200">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-slate-800 bg-[#0F172A] text-[11px] font-bold text-slate-400">
              <th class="py-2.5 px-4 w-12 text-slate-500">#</th>
              <th class="py-2.5 px-4 w-28">Tanggal</th>
              <th class="py-2.5 px-4 w-36">No Jurnal</th>
              <th class="py-2.5 px-4 w-32">Sumber</th>
              <th class="py-2.5 px-4">Uraian / Memo Transaksi</th>
              <th class="py-2.5 px-4 text-center w-24">Unit</th>
              <th class="py-2.5 px-5 text-right w-44 text-emerald-400">Debet (Dr)</th>
              <th class="py-2.5 px-5 text-right w-44 text-rose-400">Kredit (Cr)</th>
              <th class="py-2.5 px-6 text-right w-48 font-bold text-white bg-slate-800/40">Saldo Berjalan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60 font-mono">
            <!-- Opening Balance Row -->
            <tr class="bg-[#141D2E] font-bold text-slate-300">
              <td class="py-2.5 px-4 text-slate-500 font-mono text-[11px]">0</td>
              <td class="py-2.5 px-4 text-slate-400">01/11/2026</td>
              <td colspan="4" class="py-2.5 px-4 italic font-sans text-slate-400">
                Saldo Awal Buku Besar Akun {{ activeAccount?.code }} (Periode November 2026)
              </td>
              <td class="py-2.5 px-5 text-right text-slate-500">-</td>
              <td class="py-2.5 px-5 text-right text-slate-500">-</td>
              <td class="py-2.5 px-6 text-right font-bold text-slate-100 bg-slate-800/60">
                {{ formatIDR(activeAccount?.currentBalance ? activeAccount.currentBalance * 0.55 : 0) }}
              </td>
            </tr>

            <!-- Transaction Entries -->
            <tr
              v-for="(entry, index) in filteredLedgerEntries"
              :key="entry.id"
              class="hover:bg-slate-800/40 transition-colors group"
            >
              <td class="py-2.5 px-4 text-slate-500 text-[11px] select-none font-mono">
                {{ index + 1 }}
              </td>
              <td class="py-2.5 px-4 text-slate-400 whitespace-nowrap">
                {{ formatDate(entry.date) }}
              </td>
              <td class="py-2.5 px-4 font-bold text-blue-400 whitespace-nowrap">
                <button
                  @click="inspectJournalByNumber(entry.journalNumber)"
                  class="hover:underline flex items-center gap-1"
                >
                  <span>{{ entry.journalNumber }}</span>
                  <Eye class="w-3 h-3 text-slate-500 group-hover:text-blue-400" />
                </button>
              </td>
              <td class="py-2.5 px-4 text-slate-300 font-sans font-semibold text-[11px]">
                {{ entry.source }}
              </td>
              <td class="py-2.5 px-4 font-sans text-slate-200 max-w-sm">
                {{ entry.description }}
                <span class="block text-[10px] text-slate-500 font-mono">{{ entry.branchName }}</span>
              </td>
              <td class="py-2.5 px-4 text-center">
                <StatusBadge :status="entry.businessUnit" size="sm" />
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-emerald-400">
                {{ entry.debit > 0 ? formatIDR(entry.debit, false) : '-' }}
              </td>
              <td class="py-2.5 px-5 text-right font-mono font-bold text-rose-400">
                {{ entry.credit > 0 ? formatIDR(entry.credit, false) : '-' }}
              </td>
              <td class="py-2.5 px-6 text-right font-mono font-bold text-white bg-slate-800/40">
                {{ formatIDR(entry.runningBalance) }}
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-[#0F172A] font-bold border-t-2 border-slate-700 text-slate-100">
              <td colspan="6" class="py-3 px-4 text-right uppercase tracking-wider text-slate-400 font-sans">
                TOTAL MUTASI &amp; SALDO AKHIR:
              </td>
              <td class="py-3 px-5 text-right font-mono text-emerald-400 font-bold">
                {{ formatIDR(totalDebit) }}
              </td>
              <td class="py-3 px-5 text-right font-mono text-rose-400 font-bold">
                {{ formatIDR(totalCredit) }}
              </td>
              <td class="py-3 px-6 text-right font-mono text-emerald-300 text-sm font-black bg-slate-800/80">
                {{ formatIDR(endingBalance) }}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>

      <!-- Bottom Note -->
      <div class="px-4 py-3 bg-[#0F172A] border-t border-slate-800 flex items-center justify-between text-xs text-slate-400">
        <span class="font-sans">
          Menampilkan <strong>{{ filteredLedgerEntries.length }}</strong> transaksi mutasi buku besar untuk akun {{ activeAccount?.code }} - {{ activeAccount?.name }}.
        </span>
        <button
          @click="router.push('/journals')"
          class="text-blue-400 hover:underline flex items-center gap-1 font-semibold"
        >
          <span>Lihat Seluruh Voucher Jurnal</span>
          <ArrowLeft class="w-3.5 h-3.5 rotate-180" />
        </button>
      </div>
    </div>

    <!-- Modal Drilldown Voucher Jurnal -->
    <JournalDetailModal
      :show="showDetailModal"
      :journal="selectedJournal"
      @close="showDetailModal = false; selectedJournal = null"
    />
  </div>
</template>
