<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { usePayrollStore } from '@/stores/payroll'
import { useAccountingStore } from '@/stores/accounting'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useLocaleStore } from '@/stores/locale'
import { formatIDR } from '@/utils/currency'
import { formatDate } from '@/utils/date'
import SearchableSelect from '@/components/common/SearchableSelect.vue'
import type { BusinessUnit, AccountType, NormalBalance } from '@/types/accounting'
import type { Employee, PayslipItem } from '@/types/payroll'
import {
  Folder,
  ArrowLeft,
  Calendar,
  Save,
  CheckCircle2,
  AlertTriangle,
  Building2,
  FileText,
  Plus,
  Trash2,
  Sparkles,
  BookOpen,
  Calculator,
  Clock,
  Printer,
  CreditCard,
  User,
  Search,
  ChevronDown,
  X,
  FileSpreadsheet,
  Check,
  Send,
  HelpCircle,
  ShieldCheck,
  Layers
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const payrollStore = usePayrollStore()
const accountingStore = useAccountingStore()
const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const localeStore = useLocaleStore()

// Active Tab matching Brayo SPT Masa / Jurnal layout
const activeTab = ref<'DETAILS' | 'RECONCILIATION' | 'AUDIT'>('DETAILS')
const isSubmitting = ref(false)

// Available COA accounts from ledger
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

// Header & Employee State
const slipId = ref<string>('')
const slipNumber = ref<string>('')
const periodMonth = ref<number>(payrollStore.selectedMonth || (new Date().getMonth() + 1))
const periodYear = ref<number>(payrollStore.selectedYear || new Date().getFullYear())
const paymentDate = ref<string>(new Date().toISOString().slice(0, 10))
const businessUnit = ref<BusinessUnit>(accountingStore.businessUnit || 'ALL')
const status = ref<'DRAFT' | 'APPROVED' | 'PAID'>('DRAFT')
const notes = ref<string>('')

// Combobox Employee Search State
const dropdownContainerRef = ref<HTMLElement | null>(null)
const isDropdownOpen = ref<boolean>(false)
const selectedEmployeeId = ref<string>('')
const searchInputText = ref<string>('')

const employeeName = ref<string>('')
const employeeNik = ref<string>('')
const position = ref<string>('')
const department = ref<string>('')
const bankName = ref<string>('Bank Mandiri')
const bankAccountNumber = ref<string>('')
const bankPaymentSourceCode = ref<string>('1-1114') // 1-1114 Bank Mandiri Operasional Payroll

// Attendance & Working Days
const workingDays = ref<number>(22)
const presentDays = ref<number>(22)
const sickLeaveDays = ref<number>(0)
const paidLeaveDays = ref<number>(0)
const absentDays = ref<number>(0)
const overtimeHours = ref<number>(0)
const dailyMealTransportRate = ref<number>(30000)

// Earnings (Komponen Gaji)
const baseSalary = ref<number>(0)
const positionAllowance = ref<number>(0)
const transportMealAllowance = ref<number>(0)
const shiftAttendanceAllowance = ref<number>(0)
const incentiveBonus = ref<number>(0)
const overtimePay = ref<number>(0)
const otherAllowance = ref<number>(0)

// Deductions (Potongan & Pajak)
const bpjsKesehatan = ref<number>(0)
const bpjsKetenagakerjaan = ref<number>(0)
const pph21 = ref<number>(0)
const attendanceDeduction = ref<number>(0)
const loanDeduction = ref<number>(0)
const otherDeduction = ref<number>(0)

// Multi-line Journal Rows (Double-Entry Bookkeeping)
interface JournalRowState {
  accountCode: string
  accountName: string
  accountType?: AccountType
  normalBalance?: NormalBalance
  unit: BusinessUnit
  description: string
  debit: number
  credit: number
}

const journalRows = ref<JournalRowState[]>([])

// Employee Pool from Live PostgreSQL Database
const employeePool = computed<Employee[]>(() => {
  return payrollStore.allAvailableEmployees || payrollStore.employees || []
})

const selectedEmployee = computed(() => {
  if (!selectedEmployeeId.value) return null
  return employeePool.value.find(e => e.id === selectedEmployeeId.value || e.nik === selectedEmployeeId.value) || null
})

const filteredEmployees = computed(() => {
  const pool = employeePool.value
  const q = searchInputText.value.trim().toLowerCase()
  if (!q || (selectedEmployee.value && searchInputText.value === `${selectedEmployee.value.name} (${selectedEmployee.value.nik})`)) {
    return pool
  }
  return pool.filter(emp =>
    (emp.name || '').toLowerCase().includes(q) ||
    (emp.nik || '').toLowerCase().includes(q) ||
    (emp.position || '').toLowerCase().includes(q) ||
    (emp.department || '').toLowerCase().includes(q) ||
    (emp.businessUnit || '').toLowerCase().includes(q)
  )
})

function selectEmployee(emp: Employee) {
  selectedEmployeeId.value = emp.id
  employeeName.value = emp.name
  employeeNik.value = emp.nik
  position.value = emp.position
  department.value = emp.department
  businessUnit.value = emp.businessUnit
  bankName.value = emp.bankName || 'Bank Mandiri'
  bankAccountNumber.value = emp.bankAccountNumber || ''

  baseSalary.value = emp.baseSalary || 0
  positionAllowance.value = emp.positionAllowance || 0
  workingDays.value = 22
  presentDays.value = 22
  sickLeaveDays.value = 0
  paidLeaveDays.value = 0
  absentDays.value = 0
  overtimeHours.value = 0

  dailyMealTransportRate.value = emp.dailyMealTransportRate || (emp.transportMealAllowance ? Math.round(emp.transportMealAllowance / 22) : 30000)
  transportMealAllowance.value = emp.transportMealAllowance || (dailyMealTransportRate.value * presentDays.value)
  shiftAttendanceAllowance.value = 0
  incentiveBonus.value = 0
  overtimePay.value = 0
  attendanceDeduction.value = 0
  loanDeduction.value = 0

  // Standard Indonesian payroll statutory rates
  bpjsKesehatan.value = Math.round(baseSalary.value * 0.01)
  bpjsKetenagakerjaan.value = Math.round(baseSalary.value * 0.03)

  // Standard TER PPh 21 estimate
  const estimatedGross = baseSalary.value + positionAllowance.value + transportMealAllowance.value
  pph21.value = estimatedGross > 5400000 ? Math.round(estimatedGross * 0.015) : 0

  searchInputText.value = `${emp.name} (${emp.nik})`
  isDropdownOpen.value = false

  generateSlipNumber()
  syncJournalRows()
}

function selectManualMode() {
  selectedEmployeeId.value = ''
  searchInputText.value = ''
  employeeName.value = ''
  employeeNik.value = ''
  position.value = ''
  department.value = ''
  baseSalary.value = 0
  positionAllowance.value = 0
  transportMealAllowance.value = 0
  dailyMealTransportRate.value = 30000
  shiftAttendanceAllowance.value = 0
  incentiveBonus.value = 0
  overtimePay.value = 0
  bpjsKesehatan.value = 0
  bpjsKetenagakerjaan.value = 0
  pph21.value = 0
  attendanceDeduction.value = 0
  loanDeduction.value = 0
  isDropdownOpen.value = false
  syncJournalRows()
}

function handleSearchInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  searchInputText.value = val
  isDropdownOpen.value = true
  if (!selectedEmployeeId.value) {
    employeeName.value = val
  }
}

function clearEmployeeSelection() {
  selectedEmployeeId.value = ''
  searchInputText.value = ''
  selectManualMode()
  isDropdownOpen.value = true
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownContainerRef.value && !dropdownContainerRef.value.contains(event.target as Node)) {
    isDropdownOpen.value = false
    if (selectedEmployee.value) {
      searchInputText.value = `${selectedEmployee.value.name} (${selectedEmployee.value.nik})`
    }
  }
}

// Calculations
const totalGrossSalary = computed(() => {
  return (
    (Number(baseSalary.value) || 0) +
    (Number(positionAllowance.value) || 0) +
    (Number(transportMealAllowance.value) || 0) +
    (Number(shiftAttendanceAllowance.value) || 0) +
    (Number(incentiveBonus.value) || 0) +
    (Number(overtimePay.value) || 0) +
    (Number(otherAllowance.value) || 0)
  )
})

const totalDeduction = computed(() => {
  return (
    (Number(bpjsKesehatan.value) || 0) +
    (Number(bpjsKetenagakerjaan.value) || 0) +
    (Number(pph21.value) || 0) +
    (Number(attendanceDeduction.value) || 0) +
    (Number(loanDeduction.value) || 0) +
    (Number(otherDeduction.value) || 0)
  )
})

const takeHomePay = computed(() => {
  return Math.max(0, totalGrossSalary.value - totalDeduction.value)
})

// Auto update Uang Makan based on presence
watch([presentDays, dailyMealTransportRate], ([pres, rate]) => {
  transportMealAllowance.value = Math.round((Number(pres) || 0) * (Number(rate) || 0))
  syncJournalRows()
})

// Auto update Potongan Alpa
watch([absentDays, workingDays, baseSalary], ([abs, wDays, base]) => {
  const numAbs = Number(abs) || 0
  const numWDays = Number(wDays) || 22
  const numBase = Number(base) || 0
  if (numAbs > 0 && numWDays > 0) {
    attendanceDeduction.value = Math.round(numAbs * (numBase / numWDays))
  } else {
    attendanceDeduction.value = 0
  }
  syncJournalRows()
})

// Sync Journal Rows (Double-Entry Bookkeeping)
function syncJournalRows() {
  const gross = totalGrossSalary.value
  const tax = Number(pph21.value) || 0
  const bpjs = (Number(bpjsKesehatan.value) || 0) + (Number(bpjsKetenagakerjaan.value) || 0)
  const loan = Number(loanDeduction.value) || 0
  const thp = takeHomePay.value

  const bankNameStr = bankName.value || 'Bank Mandiri'
  const empNameStr = employeeName.value || 'Karyawan'
  const periodStr = `${periodMonth.value}/${periodYear.value}`

  const newRows: JournalRowState[] = [
    // 1. Debet Beban Gaji Bruto
    {
      accountCode: '6-1200',
      accountName: 'Beban Gaji Karyawan, Perawat & Staf',
      unit: businessUnit.value,
      description: `Beban Gaji Bruto Periode ${periodStr} - ${empNameStr}`,
      debit: gross,
      credit: 0
    },
    // 2. Kredit Hutang Pajak PPh 21
    {
      accountCode: '2-1510',
      accountName: 'Hutang Pajak PPh 21 Karyawan (TER PP 58)',
      unit: businessUnit.value,
      description: `Potongan PPh 21 Karyawan Masa ${periodStr} - ${empNameStr}`,
      debit: 0,
      credit: tax
    },
    // 3. Kredit Hutang Iuran BPJS
    {
      accountCode: '2-1130',
      accountName: 'Hutang Iuran BPJS (Kesehatan 1% & TK 3%)',
      unit: businessUnit.value,
      description: `Potongan Iuran BPJS Kesehatan & Ketenagakerjaan Masa ${periodStr} - ${empNameStr}`,
      debit: 0,
      credit: bpjs
    }
  ]

  // If there is loan / kasbon deduction
  if (loan > 0) {
    newRows.push({
      accountCode: '1-1170',
      accountName: 'Piutang Karyawan / Potongan Kasbon',
      unit: businessUnit.value,
      description: `Pengembalian Pinjaman / Kasbon Gaji - ${empNameStr}`,
      debit: 0,
      credit: loan
    })
  }

  // 4. Kredit Pengeluaran Kas/Bank THP
  newRows.push({
    accountCode: bankPaymentSourceCode.value,
    accountName: bankPaymentSourceCode.value === '1-1113' ? 'Bank BCA Operasional' : bankPaymentSourceCode.value === '1-1111' ? 'Kas Operasional Kasir' : 'Bank Mandiri Operasional Payroll',
    unit: businessUnit.value,
    description: `Transfer Take Home Pay (Gaji Bersih) via ${bankNameStr} - ${empNameStr}`,
    debit: 0,
    credit: thp
  })

  journalRows.value = newRows
}

watch([totalGrossSalary, totalDeduction, bankPaymentSourceCode, businessUnit], () => {
  syncJournalRows()
})

// Journal Balance Metrics
const totalDebit = computed(() => {
  return journalRows.value.reduce((sum, r) => sum + (Number(r.debit) || 0), 0)
})

const totalCredit = computed(() => {
  return journalRows.value.reduce((sum, r) => sum + (Number(r.credit) || 0), 0)
})

const difference = computed(() => {
  return totalDebit.value - totalCredit.value
})

const isBalanced = computed(() => {
  return Math.abs(difference.value) < 1 && totalDebit.value > 0
})

function addJournalRow() {
  journalRows.value.push({
    accountCode: '',
    accountName: '',
    unit: businessUnit.value,
    description: `Penyesuaian Payroll - ${employeeName.value || 'Karyawan'}`,
    debit: 0,
    credit: 0
  })
}

function removeJournalRow(idx: number) {
  journalRows.value.splice(idx, 1)
}

function onAccountSelect(index: number, code: string) {
  const acc = availableAccounts.value.find(a => a.code === code)
  if (acc && journalRows.value[index]) {
    journalRows.value[index].accountCode = acc.code
    journalRows.value[index].accountName = acc.name
    journalRows.value[index].accountType = acc.type
    journalRows.value[index].normalBalance = acc.normalBalance
  }
}

function autoBalanceJournal() {
  const diff = difference.value
  if (diff === 0 || journalRows.value.length === 0) return

  const lastRow = journalRows.value[journalRows.value.length - 1]
  if (diff > 0) {
    // Debet > Kredit, tambahkan ke kredit baris terakhir (kas/bank)
    lastRow.credit = (Number(lastRow.credit) || 0) + diff
  } else {
    // Kredit > Debet, tambahkan ke debet baris pertama (beban)
    journalRows.value[0].debit = (Number(journalRows.value[0].debit) || 0) + Math.abs(diff)
  }
  accountingStore.showToast('Jurnal penggajian berhasil diseimbangkan otomatis (Rp 0 Selisih)', 'success')
}

function generateSlipNumber() {
  const m = String(periodMonth.value).padStart(2, '0')
  const y = periodYear.value
  const nik = employeeNik.value || 'EMP'
  slipNumber.value = `SLIP/${y}/${m}/${nik}`
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside)
  await Promise.all([
    ledgerStore.fetchAccounts(),
    payrollStore.fetchEmployees(),
    payrollStore.fetchPayslips()
  ])

  // Check if edit mode via query param id
  const editId = route.query.id as string
  if (editId) {
    const slip = payrollStore.payslips.find(s => s.id === editId)
    if (slip) {
      loadSlipData(slip)
      return
    }
  }

  generateSlipNumber()
  syncJournalRows()
})

function loadSlipData(slip: PayslipItem) {
  slipId.value = slip.id
  slipNumber.value = slip.slipNumber
  selectedEmployeeId.value = slip.employeeId
  employeeName.value = slip.employeeName
  employeeNik.value = slip.employeeNik
  position.value = slip.position
  department.value = slip.department
  businessUnit.value = slip.businessUnit
  periodMonth.value = slip.periodMonth
  periodYear.value = slip.periodYear
  paymentDate.value = slip.paymentDate || new Date().toISOString().slice(0, 10)
  bankName.value = slip.bankName || 'Bank Mandiri'
  bankAccountNumber.value = slip.bankAccountNumber || ''
  status.value = slip.status

  workingDays.value = slip.workingDays || 22
  presentDays.value = slip.presentDays !== undefined ? slip.presentDays : 22
  sickLeaveDays.value = slip.sickLeaveDays || 0
  paidLeaveDays.value = slip.paidLeaveDays || 0
  absentDays.value = slip.absentDays || 0
  overtimeHours.value = slip.overtimeHours || 0
  dailyMealTransportRate.value = slip.dailyMealTransportRate || 30000

  baseSalary.value = slip.baseSalary || 0
  positionAllowance.value = slip.positionAllowance || 0
  transportMealAllowance.value = slip.transportMealAllowance || 0
  shiftAttendanceAllowance.value = slip.shiftAttendanceAllowance || 0
  incentiveBonus.value = slip.incentiveBonus || 0
  overtimePay.value = slip.overtimePay || 0
  otherAllowance.value = slip.otherAllowance || 0

  bpjsKesehatan.value = slip.bpjsKesehatan || 0
  bpjsKetenagakerjaan.value = slip.bpjsKetenagakerjaan || 0
  pph21.value = slip.pph21 || 0
  attendanceDeduction.value = slip.attendanceDeduction || 0
  loanDeduction.value = slip.loanDeduction || 0
  otherDeduction.value = slip.otherDeduction || 0

  searchInputText.value = `${slip.employeeName} (${slip.employeeNik})`
  syncJournalRows()
}

async function handleSave(postToLedger = false) {
  if (!employeeName.value.trim()) {
    accountingStore.showToast('Nama karyawan wajib diisi!', 'error')
    return
  }

  isSubmitting.value = true
  try {
    const slipPayload: Partial<PayslipItem> = {
      slipNumber: slipNumber.value,
      employeeId: selectedEmployeeId.value || `emp-${Date.now()}`,
      employeeNik: employeeNik.value || `SMG-${Date.now().toString().slice(-4)}`,
      employeeName: employeeName.value,
      position: position.value || 'Staf Medis / Umum',
      department: department.value || 'Operasional',
      businessUnit: businessUnit.value,
      periodMonth: periodMonth.value,
      periodYear: periodYear.value,
      periodName: `Bulan ${periodMonth.value} ${periodYear.value}`,
      paymentDate: paymentDate.value,
      workingDays: workingDays.value,
      presentDays: presentDays.value,
      sickLeaveDays: sickLeaveDays.value,
      paidLeaveDays: paidLeaveDays.value,
      absentDays: absentDays.value,
      overtimeHours: overtimeHours.value,
      dailyMealTransportRate: dailyMealTransportRate.value,
      baseSalary: baseSalary.value,
      positionAllowance: positionAllowance.value,
      transportMealAllowance: transportMealAllowance.value,
      shiftAttendanceAllowance: shiftAttendanceAllowance.value,
      incentiveBonus: incentiveBonus.value,
      overtimePay: overtimePay.value,
      otherAllowance: otherAllowance.value,
      totalGrossSalary: totalGrossSalary.value,
      bpjsKesehatan: bpjsKesehatan.value,
      bpjsKetenagakerjaan: bpjsKetenagakerjaan.value,
      pph21: pph21.value,
      attendanceDeduction: attendanceDeduction.value,
      loanDeduction: loanDeduction.value,
      otherDeduction: otherDeduction.value,
      totalDeduction: totalDeduction.value,
      takeHomePay: takeHomePay.value,
      bankName: bankName.value,
      bankAccountNumber: bankAccountNumber.value,
      status: postToLedger ? 'APPROVED' : status.value,
      notes: notes.value
    }

    if (slipId.value) {
      await payrollStore.updatePayslip(slipId.value, slipPayload)
    } else {
      await payrollStore.addCustomPayslip(slipPayload as any)
    }

    // If Post to Ledger is checked, also create the balanced General Journal entry
    if (postToLedger && isBalanced.value) {
      const journalLines = journalRows.value.map(row => ({
        accountCode: row.accountCode,
        accountName: row.accountName,
        debit: Number(row.debit) || 0,
        credit: Number(row.credit) || 0,
        description: row.description
      }))

      journalsStore.createManualJournal({
        journalDate: paymentDate.value,
        businessUnit: businessUnit.value,
        branchId: 'BRANCH-MAIN',
        branchName: 'Semarang Pusat',
        description: `Beban Gaji & Payroll Periode ${periodMonth.value}/${periodYear.value} - ${employeeName.value} (Ref: ${slipNumber.value})`,
        status: 'POSTED',
        journalSource: 'MANUAL',
        sourceReference: slipNumber.value,
        lines: journalLines
      })

      accountingStore.showToast(`Slip gaji & Jurnal Akrual Payroll (${slipNumber.value}) berhasil diposting ke Buku Besar!`, 'success')
    }

    router.push('/accounting/payroll')
  } catch (err: any) {
    console.error('Gagal menyimpan slip gaji:', err)
    accountingStore.showToast(err.message || 'Gagal menyimpan slip gaji', 'error')
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
          <router-link to="/accounting/payroll" class="hover:text-zinc-900 transition flex items-center gap-1.5 text-zinc-500">
            <Folder class="w-4 h-4 text-zinc-400" />
            <span>Penggajian</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-500">Slip Gaji</span>
          <span>/</span>
          <span class="font-mono font-bold text-zinc-900">{{ slipNumber || 'SLIP-BARU' }}</span>
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
            @click="autoBalanceJournal"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-amber-50 hover:bg-amber-100 text-amber-800 font-bold rounded-xl text-xs border border-amber-200 transition active:scale-95 cursor-pointer"
            title="Seimbangkan debet dan kredit secara otomatis"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-600" />
            <span>Seimbangkan Jurnal</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting"
            @click="handleSave(false)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-white hover:bg-zinc-50 disabled:opacity-50 text-zinc-800 font-bold rounded-xl text-xs border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Save class="w-3.5 h-3.5 text-zinc-600" />
            <span>Simpan Draft</span>
          </button>

          <button
            type="button"
            :disabled="isSubmitting || !isBalanced"
            @click="handleSave(true)"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 disabled:opacity-50 text-white font-bold rounded-xl text-xs shadow-2xs transition active:scale-95 cursor-pointer"
          >
            <Send class="w-3.5 h-3.5 text-emerald-400" />
            <span>Simpan & Posting Jurnal</span>
          </button>

          <router-link
            to="/accounting/payroll"
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

      <!-- 3. TAB 1: DETAILS (Inside single unified sheet) -->
      <div v-if="activeTab === 'DETAILS'" class="space-y-7">
        
        <!-- Section A: Informasi Karyawan & Periode -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Informasi Periode & Karyawan
          </h3>

          <!-- Combobox Employee Search Bar -->
          <div ref="dropdownContainerRef" class="relative mb-4">
            <label class="block text-zinc-600 mb-1 font-medium text-xs">
              Pilih / Ketik Karyawan (Database PostgreSQL) *
            </label>
            <div class="relative flex items-center">
              <Search class="w-4 h-4 text-zinc-400 absolute left-3.5 pointer-events-none" />
              <input
                type="text"
                :value="searchInputText"
                @input="handleSearchInput"
                @focus="isDropdownOpen = true"
                class="w-full pl-10 pr-20 py-2 bg-zinc-50 border border-zinc-200 rounded-xl focus:bg-white focus:outline-none focus:border-zinc-400 text-xs font-semibold text-zinc-900 transition-all placeholder:text-zinc-400"
                placeholder="-- Ketik nama atau NIK karyawan (cth: Fitri Handayani / SMG-KID) --"
              />

              <div class="absolute right-2.5 flex items-center gap-1">
                <button
                  v-if="searchInputText"
                  type="button"
                  @click.stop="clearEmployeeSelection"
                  class="p-1 text-zinc-400 hover:text-zinc-600 rounded-md transition cursor-pointer"
                  title="Reset pilihan"
                >
                  <X class="w-3.5 h-3.5" />
                </button>

                <button
                  type="button"
                  @click.stop="isDropdownOpen = !isDropdownOpen"
                  class="p-1 text-zinc-500 hover:text-zinc-800 rounded-md transition cursor-pointer"
                >
                  <ChevronDown class="w-4 h-4 transition-transform duration-150" :class="isDropdownOpen ? 'rotate-180' : ''" />
                </button>
              </div>
            </div>

            <!-- Dropdown Menu Results -->
            <div
              v-if="isDropdownOpen"
              class="absolute left-0 right-0 top-full mt-1.5 bg-white border border-zinc-200 rounded-xl shadow-2xl z-40 max-h-64 overflow-y-auto divide-y divide-zinc-100 animate-in fade-in zoom-in-95 duration-100 text-xs"
            >
              <div
                @click="selectManualMode"
                class="px-4 py-2.5 hover:bg-zinc-50 cursor-pointer transition flex items-center justify-between"
                :class="!selectedEmployeeId ? 'bg-zinc-50 font-bold text-zinc-950' : 'text-zinc-700'"
              >
                <div class="flex items-center gap-2">
                  <User class="w-4 h-4 text-zinc-400" />
                  <div>
                    <div class="font-bold">-- Input Karyawan Baru / Ketik Manual --</div>
                    <div class="text-[10px] text-zinc-400 font-normal">Isi nama, NIK, dan jabatan secara manual di formulir bawah</div>
                  </div>
                </div>
                <span class="text-[10px] px-2 py-0.5 rounded bg-zinc-200 text-zinc-700 font-bold">Manual</span>
              </div>

              <div
                v-for="emp in filteredEmployees"
                :key="emp.id"
                @click="selectEmployee(emp)"
                class="px-4 py-2.5 hover:bg-zinc-50 cursor-pointer transition flex items-center justify-between group"
                :class="selectedEmployeeId === emp.id ? 'bg-emerald-50/70 font-bold' : ''"
              >
                <div>
                  <div class="font-bold text-zinc-900 group-hover:text-zinc-950 flex items-center gap-1.5">
                    <span>{{ emp.name }}</span>
                    <Check v-if="selectedEmployeeId === emp.id" class="w-3.5 h-3.5 text-emerald-600 inline" />
                  </div>
                  <div class="text-[11px] text-zinc-500 font-mono mt-0.5">
                    {{ emp.nik }} • {{ emp.position }} ({{ emp.department }})
                  </div>
                </div>

                <div class="text-right shrink-0 ml-4">
                  <span
                    class="px-2 py-0.5 rounded text-[9px] font-bold uppercase tracking-wider border block"
                    :class="emp.businessUnit === 'APOTEK' ? 'bg-blue-50 text-blue-800 border-blue-200' : emp.businessUnit === 'SSC' ? 'bg-rose-50 text-rose-800 border-rose-200' : 'bg-emerald-50 text-emerald-800 border-emerald-200'"
                  >
                    {{ emp.businessUnit === 'APOTEK' ? 'CV SMC • Apotek' : emp.businessUnit === 'SSC' ? 'CV SMG • SSC' : 'CV SMG • Kidinest' }}
                  </span>
                  <span class="text-[10px] font-mono text-zinc-500 block mt-0.5">Rp {{ formatIDR(emp.baseSalary) }}</span>
                </div>
              </div>

              <div v-if="filteredEmployees.length === 0" class="p-4 text-center text-zinc-400 text-xs italic">
                Tidak ada karyawan yang cocok dengan pencarian "{{ searchInputText }}". Anda dapat mengisi data manual di bawah.
              </div>
            </div>
          </div>

          <!-- Metadata Form Grid (4 Cols) -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Nama Lengkap Karyawan *</label>
              <input
                v-model="employeeName"
                type="text"
                required
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white transition"
                placeholder="cth: Fitri Handayani"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">NIK / ID Karyawan *</label>
              <input
                v-model="employeeNik"
                type="text"
                required
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white transition"
                placeholder="cth: SMG-KID-2024-001"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Jabatan & Penugasan *</label>
              <input
                v-model="position"
                type="text"
                required
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-semibold focus:outline-none focus:border-zinc-400 focus:bg-white transition"
                placeholder="cth: Kasir Shift GRO / Perawat Poli Anak"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Entitas & Unit Bisnis *</label>
              <select
                v-model="businessUnit"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-semibold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
              >
                <optgroup label="CV SEMARANG MEDICAL GROUP (SMG)">
                  <option value="KIDINEST">CV SMG - Kidinest (Poli Spesialis Anak)</option>
                  <option value="SSC">CV SMG - Semarang Skin Center (SSC)</option>
                  <option value="ALL">CV SMG - Kantor Pusat / Konsolidasi</option>
                </optgroup>
                <optgroup label="CV SEMARANG MEDICAL CLINIC (SMC)">
                  <option value="APOTEK">CV SMC - Instalasi Farmasi & Apotek</option>
                </optgroup>
              </select>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Bulan & Tahun Penggajian *</label>
              <div class="grid grid-cols-2 gap-2">
                <select
                  v-model.number="periodMonth"
                  class="w-full px-2.5 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-semibold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer"
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
                <input
                  v-model.number="periodYear"
                  type="number"
                  class="w-full px-2.5 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white text-center"
                />
              </div>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Tanggal Pembayaran *</label>
              <input
                v-model="paymentDate"
                type="date"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Rekening Kas/Bank Payroll *</label>
              <select
                v-model="bankPaymentSourceCode"
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-semibold text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white cursor-pointer transition"
              >
                <option value="1-1114">[1-1114] Bank Mandiri Payroll Operasional</option>
                <option value="1-1113">[1-1113] Bank BCA Operasional Klinik</option>
                <option value="1-1111">[1-1111] Kas Operasional Kasir Klinik</option>
              </select>
            </div>

            <div>
              <label class="block text-zinc-600 mb-1 font-medium">Rekening Karyawan *</label>
              <div class="grid grid-cols-5 gap-1.5">
                <input
                  v-model="bankName"
                  type="text"
                  class="col-span-2 px-2 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-medium text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white"
                  placeholder="Bank"
                />
                <input
                  v-model="bankAccountNumber"
                  type="text"
                  class="col-span-3 px-2 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 focus:outline-none focus:border-zinc-400 focus:bg-white"
                  placeholder="No. Rekening"
                />
              </div>
            </div>

            <div class="md:col-span-2 lg:col-span-4">
              <label class="block text-zinc-600 mb-1 font-medium">Catatan / Keterangan Memo Slip Gaji</label>
              <input
                v-model="notes"
                type="text"
                placeholder="cth: Penggajian rutin bulanan staf medis & operasional Semarang Medical Group..."
                class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 font-medium focus:outline-none focus:border-zinc-400 focus:bg-white transition"
              />
            </div>
          </div>
        </div>

        <!-- Section B: Rekapitulasi Kehadiran & Nilai Komponen Gaji -->
        <div>
          <div class="flex items-center justify-between mb-3.5">
            <div>
              <h3 class="text-sm font-bold text-zinc-900">
                Komponen Gaji & Rekap Kehadiran
              </h3>
              <p class="text-xs text-zinc-500">
                Data kehadiran dan rincian penerimaan serta potongan yang membentuk Take Home Pay.
              </p>
            </div>
            <span class="text-xs font-mono font-bold text-emerald-800 bg-emerald-50 px-2.5 py-1 rounded-lg border border-emerald-200">
              Hadir: {{ presentDays }} / {{ workingDays }} Hari
            </span>
          </div>

          <!-- Attendance Inputs -->
          <div class="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-7 gap-3 text-xs mb-4">
            <div>
              <label class="block text-zinc-500 font-medium mb-1">Standar Hari Kerja</label>
              <input
                v-model.number="workingDays"
                type="number"
                class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-center text-zinc-900 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-emerald-700 font-bold mb-1">Hari Hadir *</label>
              <input
                v-model.number="presentDays"
                type="number"
                class="w-full px-3 py-1.5 bg-emerald-50/50 border border-emerald-200 rounded-xl font-mono font-bold text-center text-emerald-950 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-zinc-500 font-medium mb-1">Sakit (Surat Dokter)</label>
              <input
                v-model.number="sickLeaveDays"
                type="number"
                class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-center text-zinc-900 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-zinc-500 font-medium mb-1">Izin Resmi / Cuti</label>
              <input
                v-model.number="paidLeaveDays"
                type="number"
                class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-center text-zinc-900 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-rose-600 font-bold mb-1">Alpa / Mangkir</label>
              <input
                v-model.number="absentDays"
                type="number"
                class="w-full px-3 py-1.5 bg-rose-50/50 border border-rose-200 rounded-xl font-mono font-bold text-center text-rose-900 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-indigo-600 font-bold mb-1">Jam Lembur (OT)</label>
              <input
                v-model.number="overtimeHours"
                type="number"
                step="0.5"
                class="w-full px-3 py-1.5 bg-indigo-50/50 border border-indigo-200 rounded-xl font-mono font-bold text-center text-indigo-900 focus:bg-white focus:outline-none"
              />
            </div>

            <div>
              <label class="block text-zinc-500 font-medium mb-1">Tarif Makan/Hari</label>
              <input
                v-model.number="dailyMealTransportRate"
                type="number"
                class="w-full px-3 py-1.5 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-semibold text-right text-zinc-900 focus:bg-white focus:outline-none"
              />
            </div>
          </div>

          <!-- Components Breakdown (3 Cols) -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-5 text-xs bg-zinc-50/60 border border-zinc-200 rounded-xl p-4">
            <!-- Col 1: Earnings -->
            <div class="space-y-2.5">
              <div class="flex items-center justify-between border-b border-zinc-200 pb-1.5 font-bold text-zinc-900">
                <span>A. KOMPONEN PENERIMAAN</span>
                <span class="font-mono text-emerald-700">Rp {{ formatIDR(totalGrossSalary) }}</span>
              </div>
              <div class="space-y-1.5 text-zinc-600 text-[11px]">
                <div class="flex items-center justify-between">
                  <span>Gaji Pokok:</span>
                  <input v-model.number="baseSalary" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right font-bold text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Tunjangan Jabatan:</span>
                  <input v-model.number="positionAllowance" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Uang Makan & Transport:</span>
                  <input v-model.number="transportMealAllowance" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Tunjangan Shift / Hadir:</span>
                  <input v-model.number="shiftAttendanceAllowance" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Insentif / Bonus Kinerja:</span>
                  <input v-model.number="incentiveBonus" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Upah Lembur:</span>
                  <input v-model.number="overtimePay" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
              </div>
            </div>

            <!-- Col 2: Deductions -->
            <div class="space-y-2.5">
              <div class="flex items-center justify-between border-b border-zinc-200 pb-1.5 font-bold text-zinc-900">
                <span>B. KOMPONEN POTONGAN</span>
                <span class="font-mono text-rose-700">Rp {{ formatIDR(totalDeduction) }}</span>
              </div>
              <div class="space-y-1.5 text-zinc-600 text-[11px]">
                <div class="flex items-center justify-between">
                  <span>BPJS Kesehatan (1%):</span>
                  <input v-model.number="bpjsKesehatan" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>BPJS TK (3%):</span>
                  <input v-model.number="bpjsKetenagakerjaan" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>PPh 21 (TER PP 58):</span>
                  <input v-model.number="pph21" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Potongan Mangkir/Alpa:</span>
                  <input v-model.number="attendanceDeduction" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
                <div class="flex items-center justify-between">
                  <span>Potongan Kasbon:</span>
                  <input v-model.number="loanDeduction" type="number" class="w-28 px-2 py-1 bg-white border border-zinc-200 rounded-lg font-mono text-right text-zinc-900 focus:outline-none" />
                </div>
              </div>
            </div>

            <!-- Col 3: Take Home Pay Summary -->
            <div class="space-y-2.5 bg-white p-3.5 rounded-xl border border-zinc-200 shadow-2xs flex flex-col justify-between">
              <div>
                <span class="text-zinc-500 block text-[10px] mb-1 font-bold">TOTAL TAKE HOME PAY (GAJI BERSIH)</span>
                <div class="font-mono text-2xl font-black text-emerald-700 tracking-tight">
                  Rp {{ formatIDR(takeHomePay) }}
                </div>
                <p class="text-[10px] text-zinc-400 mt-1">
                  Transfer via {{ bankPaymentSourceCode === '1-1113' ? 'Bank BCA' : 'Bank Mandiri' }} ke {{ bankName }} ({{ bankAccountNumber || '-' }})
                </p>
              </div>

              <div class="pt-2 border-t border-zinc-100 space-y-1 text-[11px]">
                <div class="flex justify-between text-zinc-500">
                  <span>Total Bruto:</span>
                  <span class="font-mono font-bold text-zinc-900">Rp {{ formatIDR(totalGrossSalary) }}</span>
                </div>
                <div class="flex justify-between text-zinc-500">
                  <span>Total Potongan:</span>
                  <span class="font-mono font-bold text-rose-600">-Rp {{ formatIDR(totalDeduction) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Section C: Detail Transaksi Jurnal Penggajian (Double-Entry Bookkeeping) -->
        <div>
          <div class="flex items-center justify-between mb-3.5">
            <div>
              <h3 class="text-sm font-bold text-zinc-900">
                Detail Transaksi Jurnal Penggajian ({{ journalRows.length }} Baris Akuntansi)
              </h3>
              <p class="text-xs text-zinc-500">
                Posisi akun akrual beban gaji (Debet) berpasangan dengan hutang pajak/BPJS & pengeluaran kas/bank THP (Kredit).
              </p>
            </div>

            <button
              type="button"
              @click="addJournalRow"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition shadow-2xs cursor-pointer active:scale-95"
            >
              <Plus class="w-3.5 h-3.5" />
              <span>Tambah Baris Jurnal</span>
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
                    v-for="(row, idx) in journalRows"
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
                        <option value="APOTEK">Apotek</option>
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
                          step="1000"
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
                          step="1000"
                          class="w-full pl-8 pr-2.5 py-1.5 bg-zinc-50 border border-zinc-200 rounded-lg text-right font-mono font-bold text-rose-700 focus:outline-none focus:border-zinc-400 focus:bg-white"
                          placeholder="0"
                        />
                      </div>
                    </td>

                    <td class="px-3 py-2 text-center">
                      <button
                        type="button"
                        @click="removeJournalRow(idx)"
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
                      Total Jurnal Penggajian:
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

        <!-- Section D: Ringkasan & Keseimbangan Jurnal Gaji -->
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-3.5">
            Ringkasan & Keseimbangan Jurnal Penggajian
          </h3>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs font-mono">
            <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block mb-1">TOTAL DEBET (BEBAN GAJI)</span>
              <span class="text-base font-extrabold text-emerald-700">Rp {{ formatIDR(totalDebit) }}</span>
            </div>

            <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl">
              <span class="text-[11px] font-sans font-medium text-zinc-500 block mb-1">TOTAL KREDIT (POTONGAN + THP)</span>
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

      <!-- 4. TAB 2: RECONCILIATION & STATUTORY COMPLIANCE -->
      <div v-else-if="activeTab === 'RECONCILIATION'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Kepatuhan Statutori & Rekonsiliasi Pajak PPh 21 / BPJS</h3>
          <p class="text-xs text-zinc-500">Validasi aturan pemotongan pajak penghasilan tenaga kerja dan jaminan sosial ketenagakerjaan.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
          <div class="p-4 rounded-xl bg-zinc-50 border border-zinc-200 space-y-2.5">
            <h4 class="font-bold text-zinc-900 flex items-center gap-1.5">
              <ShieldCheck class="w-4 h-4 text-emerald-600" />
              <span>Pajak Penghasilan Pasal 21 (TER PP 58)</span>
            </h4>
            <ul class="list-disc list-inside space-y-1 text-zinc-600 leading-relaxed text-[11px]">
              <li>Penghasilan Bruto Karyawan: <strong>Rp {{ formatIDR(totalGrossSalary) }}</strong></li>
              <li>Status PTKP Karyawan: <strong>{{ selectedEmployee?.ptkpStatus || 'TK/0' }}</strong></li>
              <li>Potongan PPh 21 Ditahan: <strong>Rp {{ formatIDR(pph21) }}</strong></li>
              <li>Lawan Akun Buku Besar: <strong>[2-1510] Hutang Pajak PPh 21 Karyawan</strong></li>
              <li>Pelaporan: Terkoneksi otomatis dengan SPT Masa Pajak Bulanan</li>
            </ul>
          </div>

          <div class="p-4 rounded-xl bg-zinc-50 border border-zinc-200 space-y-2.5">
            <h4 class="font-bold text-zinc-900 flex items-center gap-1.5">
              <ShieldCheck class="w-4 h-4 text-blue-600" />
              <span>Iuran BPJS Kesehatan & Ketenagakerjaan</span>
            </h4>
            <ul class="list-disc list-inside space-y-1 text-zinc-600 leading-relaxed text-[11px]">
              <li>BPJS Kesehatan (1%): <strong>Rp {{ formatIDR(bpjsKesehatan) }}</strong></li>
              <li>BPJS Ketenagakerjaan (3%): <strong>Rp {{ formatIDR(bpjsKetenagakerjaan) }}</strong></li>
              <li>Total Iuran Potongan BPJS: <strong>Rp {{ formatIDR(bpjsKesehatan + bpjsKetenagakerjaan) }}</strong></li>
              <li>Lawan Akun Buku Besar: <strong>[2-1130] Hutang Iuran BPJS Karyawan</strong></li>
              <li>Penyetoran: Kolektif via virtual account per tanggal 10</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 5. TAB 3: AUDIT TRAIL -->
      <div v-else-if="activeTab === 'AUDIT'" class="space-y-6">
        <div>
          <h3 class="text-sm font-bold text-zinc-900 mb-1">Audit Trail & Jejak Otentikasi Penggajian</h3>
          <p class="text-xs text-zinc-500">Informasi riwayat pembuatan voucher penggajian dan otorisasi posting.</p>
        </div>

        <div class="space-y-3 text-xs">
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Petugas Pembuat:</span>
            <span class="font-bold text-zinc-900">Budi Santoso, SE, Ak., CA (Finance & Accounting)</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Tanggal Pembayaran:</span>
            <span class="font-bold text-zinc-900 font-mono">{{ paymentDate }}</span>
          </div>
          <div class="p-3.5 bg-zinc-50 border border-zinc-200 rounded-xl flex items-center justify-between">
            <span class="text-zinc-600 font-medium">Rekening Sumber Payroll:</span>
            <span class="font-bold text-zinc-900 font-mono">[{{ bankPaymentSourceCode }}] {{ bankPaymentSourceCode === '1-1113' ? 'Bank BCA' : 'Bank Mandiri Payroll' }}</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>
