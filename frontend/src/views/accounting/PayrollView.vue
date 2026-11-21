<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { usePayrollStore } from '@/stores/payroll'
import { useAccountingStore } from '@/stores/accounting'
import { useLocaleStore } from '@/stores/locale'
import { formatIDR } from '@/utils/currency'
import { exportToExcel } from '@/utils/export'
import type { PayslipItem } from '@/types/payroll'
import PayslipDetailModal from '@/components/payroll/PayslipDetailModal.vue'
import PayslipFormModal from '@/components/payroll/PayslipFormModal.vue'
import AttendanceImportModal from '@/components/payroll/AttendanceImportModal.vue'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  Wallet,
  CreditCard,
  Plus,
  Printer,
  Download,
  FileSpreadsheet,
  Upload,
  Search,
  Filter,
  CheckCircle2,
  Clock,
  Sparkles,
  Building2,
  Users,
  Eye,
  Trash2,
  Edit,
  Send,
  BookOpen,
  FileCheck2,
  Check
} from 'lucide-vue-next'

const router = useRouter()
const payrollStore = usePayrollStore()
const accountingStore = useAccountingStore()
const localeStore = useLocaleStore()

const selectedSlip = ref<PayslipItem | null>(null)
const showDetailModal = ref(false)
const showFormModal = ref(false)
const showImportModal = ref(false)
const editingSlip = ref<PayslipItem | null>(null)
const showPrintRecap = ref(false)

// Entity separation filter (Master Prompt 2 & 9 F)
// CV SMG (Kidinest, SSC) vs CV SMC (Apotek) vs ALL
const selectedEntity = ref<'ALL' | 'SMG' | 'SMC'>('ALL')

onMounted(() => {
  payrollStore.fetchEmployees()
  payrollStore.fetchPayslips()
})

watch([() => payrollStore.selectedMonth, () => payrollStore.selectedYear, () => accountingStore.businessUnit], () => {
  payrollStore.fetchPayslips()
})

const entityFilteredPayslips = computed(() => {
  return payrollStore.filteredPayslips.filter(s => {
    if (selectedEntity.value === 'ALL') return true
    if (selectedEntity.value === 'SMC') {
      return s.businessUnit === 'APOTEK' || s.businessUnit === 'SMC' || s.department.toLowerCase().includes('farmasi') || s.department.toLowerCase().includes('apotek')
    }
    // SMG
    return s.businessUnit === 'KIDINEST' || s.businessUnit === 'SSC' || s.businessUnit === 'ALL' || s.businessUnit === 'SMG'
  })
})

function isClinicEntity(slip: PayslipItem): boolean {
  return slip.businessUnit === 'APOTEK' || slip.businessUnit === 'SMC' || slip.department.toLowerCase().includes('farmasi') || slip.department.toLowerCase().includes('apotek')
}

function getEntityName(slip: PayslipItem): string {
  return isClinicEntity(slip) ? 'CV Semarang Medical Clinic' : 'CV Semarang Medical Group'
}

function getUnitBadgeClass(unit: string): string {
  switch (unit) {
    case 'APOTEK':
    case 'SMC':
      return 'bg-blue-50 text-blue-800 border-blue-200'
    case 'SSC':
      return 'bg-rose-50 text-rose-800 border-rose-200'
    case 'KIDINEST':
      return 'bg-emerald-50 text-emerald-800 border-emerald-200'
    default:
      return 'bg-zinc-100 text-zinc-700 border-zinc-200'
  }
}

function getUnitLabel(unit: string): string {
  switch (unit) {
    case 'APOTEK':
    case 'SMC':
      return 'Apotek SMC'
    case 'SSC':
      return 'SSC Skin'
    case 'KIDINEST':
      return 'Kidinest'
    default:
      return 'Kantor Pusat'
  }
}

function handleDeleteSlip(slip: PayslipItem) {
  if (confirm(`Apakah Anda yakin ingin menghapus slip gaji ${slip.employeeName} (${slip.employeeNik}) periode ${slip.periodMonth}/${slip.periodYear}?`)) {
    payrollStore.deletePayslip(slip.id)
  }
}

function handleDeleteAllInPeriod() {
  if (confirm(`Hapus seluruh (${entityFilteredPayslips.value.length}) slip gaji periode ${payrollStore.selectedMonth}/${payrollStore.selectedYear}? Tindakan ini akan menghapus data pada periode aktif.`)) {
    payrollStore.deleteAllInPeriod(payrollStore.selectedMonth, payrollStore.selectedYear)
  }
}

const entityStats = computed(() => {
  const list = entityFilteredPayslips.value
  const totalGross = list.reduce((s, item) => s + item.totalGrossSalary, 0)
  const totalDeductions = list.reduce((s, item) => s + item.totalDeduction, 0)
  const totalTakeHomePay = list.reduce((s, item) => s + item.takeHomePay, 0)
  const totalPph21 = list.reduce((s, item) => s + item.pph21, 0)
  const totalBpjs = list.reduce((s, item) => s + item.bpjsKesehatan + item.bpjsKetenagakerjaan, 0)
  const totalBaseSalary = list.reduce((s, item) => s + item.baseSalary, 0)
  const totalAllowances = list.reduce((s, item) => s + (item.positionAllowance + item.transportMealAllowance), 0)
  const totalOvertimePay = list.reduce((s, item) => s + (item.overtimePay || 0), 0)
  const totalBonus = list.reduce((s, item) => s + ((item.shiftAttendanceAllowance || 0) + (item.incentiveBonus || 0)), 0)

  return {
    totalCount: list.length,
    draftCount: list.filter(s => s.status === 'DRAFT').length,
    approvedCount: list.filter(s => s.status === 'APPROVED').length,
    paidCount: list.filter(s => s.status === 'PAID').length,
    totalGross,
    totalDeductions,
    totalTakeHomePay,
    totalPph21,
    totalBpjs,
    totalBaseSalary,
    totalAllowances,
    totalOvertimePay,
    totalBonus
  }
})

function viewPayslip(slip: PayslipItem) {
  selectedSlip.value = slip
  showDetailModal.value = true
}

function openEditModal(slip: PayslipItem) {
  editingSlip.value = slip
  showFormModal.value = true
}

function openCreateModal() {
  editingSlip.value = null
  showFormModal.value = true
}

function handleGeneratePayroll() {
  payrollStore.generateMonthlyPayroll(payrollStore.selectedMonth, payrollStore.selectedYear)
}

function handlePostToLedger() {
  if (confirm(`Posting jurnal payroll periode ${payrollStore.selectedMonth}/${payrollStore.selectedYear} ke Buku Besar?`)) {
    payrollStore.postPayrollToLedger(payrollStore.selectedMonth, payrollStore.selectedYear)
  }
}

function handlePrintRecap() {
  window.print()
}

function handleExportExcel() {
  const data = entityFilteredPayslips.value.map((s, idx) => ({
    'No': idx + 1,
    'No Slip': s.slipNumber,
    'NIK': s.employeeNik,
    'Nama Karyawan': s.employeeName,
    'Jabatan': s.position,
    'Departemen': s.department,
    'Unit Bisnis': s.businessUnit,
    'Entitas Perusahaan': (s.businessUnit === 'APOTEK' || s.department.toLowerCase().includes('farmasi')) ? 'CV Semarang Medical Clinic' : 'CV Semarang Medical Group',
    'Periode': s.periodName,
    'Hari Kerja': s.workingDays || 22,
    'Hari Hadir': s.presentDays !== undefined ? s.presentDays : (s.workingDays || 22),
    'Izin / Cuti': s.paidLeaveDays || 0,
    'Sakit': s.sickLeaveDays || 0,
    'Alpa (Mangkir)': s.absentDays || 0,
    'Jam Lembur': s.overtimeHours || 0,
    'Nilai Lembur (IDR)': s.overtimePay || 0,
    'Gaji Pokok (IDR)': s.baseSalary,
    'Tunjangan Jabatan (IDR)': s.positionAllowance,
    'Tunjangan Makan & Transport (IDR)': s.transportMealAllowance,
    'Shift & Insentif (IDR)': (s.shiftAttendanceAllowance || 0) + (s.incentiveBonus || 0),
    'Total Pendapatan Bruto (IDR)': s.totalGrossSalary,
    'BPJS Kesehatan (IDR)': s.bpjsKesehatan,
    'BPJS Ketenagakerjaan (IDR)': s.bpjsKetenagakerjaan,
    'PPh Pasal 21 (IDR)': s.pph21,
    'Potongan Alpa & Kasbon (IDR)': (s.attendanceDeduction || 0) + (s.loanDeduction || 0),
    'Total Potongan (IDR)': s.totalDeduction,
    'Gaji Bersih / THP (IDR)': s.takeHomePay,
    'Bank': s.bankName,
    'No Rekening': s.bankAccountNumber,
    'Status': s.status
  }))

  const entityTag = selectedEntity.value === 'ALL' ? 'KONSOLIDASI' : selectedEntity.value
  const filename = `Rekap_Payroll_${entityTag}_${payrollStore.selectedYear}_${payrollStore.selectedMonth}`
  exportToExcel(data, filename, 'Rekap Penggajian')
  accountingStore.showToast('Data rekap penggajian berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- Printable Period Payroll Summary Sheet (Shown only on Print if showPrintRecap is active) -->
    <div id="rekap-print-area" class="hidden print:block text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] p-8 space-y-6">
      <div class="flex items-center justify-between pb-4 border-b-2 border-zinc-950">
        <div>
          <h1 class="text-xl font-black tracking-tight">
            {{ selectedEntity === 'SMC' ? 'CV SEMARANG MEDICAL CLINIC' : 'CV SEMARANG MEDICAL GROUP' }}
          </h1>
          <p class="text-xs font-bold text-zinc-600">
            {{ selectedEntity === 'SMC' ? 'Instalasi Farmasi, Apotek & Logistik Medis' : 'Kidinest Pediatric Care & Semarang Skin Center' }}
          </p>
          <p class="text-[10px] text-zinc-400">Jl. Melati Utara No. 3, Kota Semarang • Telp: (024) 8412345</p>
        </div>
        <div class="text-right">
          <div class="text-[10px] font-bold uppercase text-zinc-500">REKAPITULASI PENGGAJIAN</div>
          <div class="text-base font-black">PERIODE: {{ payrollStore.selectedMonth }}/{{ payrollStore.selectedYear }}</div>
          <div class="text-xs text-zinc-600 font-mono">Dicetak: {{ new Date().toLocaleDateString('id-ID') }}</div>
        </div>
      </div>

      <table class="w-full text-left text-[10px] border-collapse border border-zinc-300">
        <thead class="bg-zinc-100 font-bold uppercase text-center border-b border-zinc-300">
          <tr>
            <th class="border border-zinc-300 p-1.5">No</th>
            <th class="border border-zinc-300 p-1.5 text-left">NIK & Nama Karyawan</th>
            <th class="border border-zinc-300 p-1.5 text-left">Jabatan</th>
            <th class="border border-zinc-300 p-1.5">Hadir</th>
            <th class="border border-zinc-300 p-1.5">Alpa</th>
            <th class="border border-zinc-300 p-1.5">Lembur</th>
            <th class="border border-zinc-300 p-1.5 text-right">Gaji Pokok</th>
            <th class="border border-zinc-300 p-1.5 text-right">Tunjangan</th>
            <th class="border border-zinc-300 p-1.5 text-right">Bruto</th>
            <th class="border border-zinc-300 p-1.5 text-right">Potongan</th>
            <th class="border border-zinc-300 p-1.5 text-right">THP Bersih</th>
            <th class="border border-zinc-300 p-1.5">Status</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(slip, idx) in entityFilteredPayslips" :key="slip.id" class="border-b border-zinc-200">
            <td class="border border-zinc-300 p-1 text-center">{{ idx + 1 }}</td>
            <td class="border border-zinc-300 p-1 font-semibold">{{ slip.employeeName }}<br/><span class="font-mono text-[9px] text-zinc-500">{{ slip.employeeNik }}</span></td>
            <td class="border border-zinc-300 p-1">{{ slip.position }}</td>
            <td class="border border-zinc-300 p-1 text-center font-bold">{{ slip.presentDays || 22 }}</td>
            <td class="border border-zinc-300 p-1 text-center">{{ slip.absentDays || 0 }}</td>
            <td class="border border-zinc-300 p-1 text-center">{{ slip.overtimeHours ? `+${slip.overtimeHours}j` : '—' }}</td>
            <td class="border border-zinc-300 p-1 text-right font-mono">{{ formatIDR(slip.baseSalary) }}</td>
            <td class="border border-zinc-300 p-1 text-right font-mono">{{ formatIDR(slip.positionAllowance + slip.transportMealAllowance + (slip.overtimePay || 0)) }}</td>
            <td class="border border-zinc-300 p-1 text-right font-mono font-bold">{{ formatIDR(slip.totalGrossSalary) }}</td>
            <td class="border border-zinc-300 p-1 text-right font-mono text-rose-700">-{{ formatIDR(slip.totalDeduction) }}</td>
            <td class="border border-zinc-300 p-1 text-right font-mono font-black">{{ formatIDR(slip.takeHomePay) }}</td>
            <td class="border border-zinc-300 p-1 text-center font-bold text-[9px]">{{ slip.status }}</td>
          </tr>
        </tbody>
        <tfoot class="bg-zinc-100 font-bold border-t-2 border-zinc-950">
          <tr>
            <td colspan="6" class="border border-zinc-300 p-2 text-right">TOTAL KESELURUHAN:</td>
            <td colspan="2" class="border border-zinc-300 p-2 text-right font-mono">{{ formatIDR(entityStats.totalGross) }}</td>
            <td class="border border-zinc-300 p-2 text-right font-mono font-bold">{{ formatIDR(entityStats.totalGross) }}</td>
            <td class="border border-zinc-300 p-2 text-right font-mono text-rose-700">-{{ formatIDR(entityStats.totalDeductions) }}</td>
            <td class="border border-zinc-300 p-2 text-right font-mono font-black">{{ formatIDR(entityStats.totalTakeHomePay) }}</td>
            <td class="border border-zinc-300 p-2"></td>
          </tr>
        </tfoot>
      </table>

      <!-- Signatures -->
      <div class="grid grid-cols-3 gap-6 pt-10 text-center text-xs">
        <div>
          <div class="text-zinc-500 mb-14">Dibuat Oleh (HR & Personalia):</div>
          <div class="font-bold underline">Wahyu Pratama, S.Psi</div>
          <div class="text-[10px] text-zinc-400">Staf HR & Payroll</div>
        </div>
        <div>
          <div class="text-zinc-500 mb-14">Diperiksa (Finance & Tax):</div>
          <div class="font-bold underline">Ratna Paramita, S.E.</div>
          <div class="text-[10px] text-zinc-400">Accounting & Finance Lead</div>
        </div>
        <div>
          <div class="text-zinc-500 mb-14">Disetujui Oleh:</div>
          <div class="font-bold underline">Direktur Operasional</div>
          <div class="text-[10px] text-zinc-400">{{ selectedEntity === 'SMC' ? 'CV Semarang Medical Clinic' : 'CV Semarang Medical Group' }}</div>
        </div>
      </div>
    </div>

    <!-- Screen UI Header -->
    <header class="flex flex-col lg:flex-row lg:items-center justify-between gap-4 pb-4 border-b border-zinc-200 print:hidden">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Manajemen Personalia & Akuntansi Penggajian
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs font-bold text-emerald-800">Semarang Medical Group & Clinic</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Slip Gaji & Penggajian Karyawan
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Import data absensi Excel, generate formulir data grid Jurnal Umum, posting buku besar, dan cetak slip individual / rekap.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-wrap items-center gap-2.5">
        <!-- Button: Import Data Absensi Excel (Master Prompt 9 A, 9 B) -->
        <button
          @click="showImportModal = true"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-emerald-700 hover:bg-emerald-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Upload class="w-3.5 h-3.5" />
          <span>Import Absensi Excel</span>
        </button>

        <!-- Button: Auto-Generate Monthly Fallback -->
        <button
          @click="handleGeneratePayroll"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-900 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Sparkles class="w-3.5 h-3.5 text-amber-300" />
          <span>Generate Otomatis</span>
        </button>

        <!-- Button: Tambah Slip Manual -->
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5 text-zinc-600" />
          <span>Input Manual</span>
        </button>

        <!-- Button: Post to Ledger -->
        <button
          @click="handlePostToLedger"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
          title="Posting Beban Gaji & Kasbon ke Jurnal Umum Buku Besar"
        >
          <BookOpen class="w-3.5 h-3.5 text-indigo-600" />
          <span>Posting Buku Besar</span>
        </button>

        <!-- Button: Cetak Rekap PDF -->
        <button
          @click="handlePrintRecap"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5 text-zinc-600" />
          <span>Cetak Rekap</span>
        </button>

        <!-- Button: Export Excel -->
        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-emerald-600" />
          <span>Excel</span>
        </button>
      </div>
    </header>

    <!-- Entity Separation Bar & Filters (Master Prompt 2 & 9 F) -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-4 bg-white rounded-2xl border border-zinc-200/90 shadow-2xs print:hidden">
      <!-- Entity Pills -->
      <div class="flex items-center gap-1.5 bg-zinc-100 p-1 rounded-xl">
        <button
          type="button"
          @click="selectedEntity = 'ALL'"
          class="px-3 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer"
          :class="selectedEntity === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
        >
          Semua Entitas
        </button>
        <button
          type="button"
          @click="selectedEntity = 'SMG'"
          class="px-3 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer"
          :class="selectedEntity === 'SMG' ? 'bg-[#0C3B36] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
        >
          CV Semarang Medical Group (Kidinest & SSC)
        </button>
        <button
          type="button"
          @click="selectedEntity = 'SMC'"
          class="px-3 py-1.5 rounded-lg text-xs font-bold transition cursor-pointer"
          :class="selectedEntity === 'SMC' ? 'bg-[#1E3A8A] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
        >
          CV Semarang Medical Clinic (Apotek)
        </button>
      </div>

      <!-- Quick Summary in Header -->
      <div class="text-xs text-zinc-500 font-medium">
        Menampilkan <strong>{{ entityStats.totalCount }}</strong> karyawan aktif terdaftar
      </div>
    </div>

    <!-- 1-Row Compact Hairline KPI Ticker Strip -->
    <div class="grid grid-cols-2 sm:grid-cols-5 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden print:hidden">
      <div class="p-3.5 border-r border-b sm:border-b-0 sm:border-r border-zinc-200/90">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">
          TOTAL BEBAN GAJI BRUTO
        </div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(entityStats.totalGross) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">
          {{ entityStats.totalCount }} Karyawan Aktif
        </div>
      </div>

      <div class="p-3.5 border-b sm:border-b-0 sm:border-r border-zinc-200/90">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">
          TOTAL POTONGAN
        </div>
        <div class="text-base font-bold tabular-nums tracking-tight text-rose-600 mt-0.5">
          {{ formatIDR(entityStats.totalDeductions) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">
          BPJS, PPh 21 & Alpa
        </div>
      </div>

      <div class="p-3.5 border-r border-b sm:border-b-0 sm:border-r border-zinc-200/90">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">
          PPH 21 TERUTANG
        </div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(entityStats.totalPph21) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">
          Pasal 17 / TER PP 58
        </div>
      </div>

      <div class="p-3.5 border-b sm:border-b-0 sm:border-r border-zinc-200/90">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">
          IURAN BPJS (KES & TK)
        </div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatIDR(entityStats.totalBpjs) }}
        </div>
        <div class="text-[10px] text-zinc-400 font-medium mt-0.5">
          JHT, JP & BPJS Kes
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-emerald-50/70 border-t sm:border-t-0 border-zinc-200/90 flex flex-col justify-center">
        <div class="text-[10px] font-bold uppercase tracking-wider text-emerald-900">
          TOTAL GAJI BERSIH (NET THP)
        </div>
        <div class="text-lg font-black tabular-nums tracking-tight text-emerald-950 mt-0.5">
          {{ formatIDR(entityStats.totalTakeHomePay) }}
        </div>
        <div class="text-[10px] text-emerald-700 font-bold mt-0.5">
          Siap Ditransfer via Payroll
        </div>
      </div>
    </div>

    <!-- Filter & Action Controls Strip -->
    <div class="p-4 bg-white rounded-2xl border border-zinc-200/90 shadow-2xs space-y-3 print:hidden">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <!-- Search Input -->
        <div class="relative flex-1 max-w-sm">
          <Search class="w-4 h-4 text-zinc-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            v-model="payrollStore.searchQuery"
            type="text"
            class="w-full pl-9 pr-4 py-2 border border-zinc-200 rounded-xl text-xs focus:ring-2 focus:ring-zinc-900 focus:outline-hidden transition-all bg-zinc-50/50 hover:bg-white"
            placeholder="Cari nama karyawan, NIK, atau jabatan..."
          />
        </div>

        <!-- Dropdown Filters & Batch Actions -->
        <div class="flex flex-wrap items-center gap-2">
          <!-- Month Selector -->
          <select
            v-model="payrollStore.selectedMonth"
            class="px-3 py-2 border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-800 bg-white focus:ring-2 focus:ring-zinc-900 focus:outline-hidden cursor-pointer"
          >
            <option :value="1">Januari</option>
            <option :value="2">Februari</option>
            <option :value="3">Maret</option>
            <option :value="4">April</option>
            <option :value="5">Mei</option>
            <option :value="6">Juni</option>
            <option :value="7">Juli</option>
            <option :value="8">Agustus</option>
            <option :value="9">September</option>
            <option :value="10">Oktober</option>
            <option :value="11">November</option>
            <option :value="12">Desember</option>
          </select>

          <!-- Year Selector -->
          <select
            v-model="payrollStore.selectedYear"
            class="px-3 py-2 border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-800 bg-white focus:ring-2 focus:ring-zinc-900 focus:outline-hidden cursor-pointer"
          >
            <option :value="2025">2025</option>
            <option :value="2026">2026</option>
            <option :value="2027">2027</option>
          </select>

          <!-- Status Selector -->
          <select
            v-model="payrollStore.statusFilter"
            class="px-3 py-2 border border-zinc-200 rounded-xl text-xs font-semibold text-zinc-800 bg-white focus:ring-2 focus:ring-zinc-900 focus:outline-hidden cursor-pointer"
          >
            <option value="ALL">Semua Status</option>
            <option value="DRAFT">DRAFT</option>
            <option value="APPROVED">APPROVED (Disetujui)</option>
            <option value="PAID">PAID (Ditransfer)</option>
          </select>

          <!-- Batch Action: Approve All -->
          <button
            @click="payrollStore.approveAllInPeriod(payrollStore.selectedMonth, payrollStore.selectedYear)"
            class="px-3 py-2 bg-amber-50 hover:bg-amber-100 text-amber-900 border border-amber-200 rounded-xl text-xs font-bold transition-colors cursor-pointer"
            title="Setujui seluruh draft slip gaji pada bulan aktif"
          >
            Setujui Semua
          </button>

          <!-- Batch Action: Mark All as Paid -->
          <button
            @click="payrollStore.markAllAsPaid(payrollStore.selectedMonth, payrollStore.selectedYear)"
            class="px-3 py-2 bg-emerald-50 hover:bg-emerald-100 text-emerald-800 border border-emerald-200 rounded-xl text-xs font-bold transition-colors cursor-pointer"
            title="Tandai seluruh slip gaji telah lunas ditransfer"
          >
            Tandai Lunas Semua
          </button>

          <!-- Batch Action: Hapus Semua -->
          <button
            v-if="entityFilteredPayslips.length > 0"
            @click="handleDeleteAllInPeriod"
            class="px-3 py-2 bg-rose-50 hover:bg-rose-100 text-rose-800 border border-rose-200 rounded-xl text-xs font-bold transition-colors cursor-pointer"
            title="Hapus seluruh slip gaji pada periode aktif"
          >
            Hapus Semua
          </button>
        </div>
      </div>
    </div>

    <!-- FORMULIR SLIP GAJI MIRIP TABEL JURNAL UMUM (DATA GRID 16 KOLOM - MASTER PROMPT 9 C) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden print:hidden">
      <!-- Empty State -->
      <div
        v-if="entityFilteredPayslips.length === 0"
        class="py-16 text-center space-y-3 px-4"
      >
        <div class="w-12 h-12 rounded-2xl bg-zinc-100 flex items-center justify-center text-zinc-400 mx-auto">
          <Wallet class="w-6 h-6" />
        </div>
        <h3 class="text-sm font-bold text-zinc-900">
          Belum ada data slip gaji untuk periode {{ payrollStore.selectedMonth }}/{{ payrollStore.selectedYear }}
        </h3>
        <p class="text-xs text-zinc-500 max-w-sm mx-auto">
          Unggah file Excel absensi untuk membaca otomatis hari hadir dan alpa, atau klik tombol di bawah untuk membuat draft otomatis.
        </p>
        <div class="flex items-center justify-center gap-2 pt-1">
          <button
            @click="showImportModal = true"
            class="inline-flex items-center gap-1.5 px-4 py-2 bg-emerald-700 hover:bg-emerald-800 text-white rounded-xl text-xs font-bold transition-all shadow-2xs cursor-pointer active:scale-95"
          >
            <Upload class="w-3.5 h-3.5" />
            <span>Import Absensi Excel</span>
          </button>
          <button
            @click="handleGeneratePayroll"
            class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-900 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition-all shadow-2xs cursor-pointer active:scale-95"
          >
            <Sparkles class="w-3.5 h-3.5 text-amber-300" />
            <span>Generate Otomatis</span>
          </button>
        </div>
      </div>

      <!-- Table Header Status Strip with Scroll Hint -->
      <div v-else class="px-4 py-2.5 bg-zinc-50/90 border-b border-zinc-200 flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-xs text-zinc-600">
        <div class="flex items-center gap-2">
          <span class="inline-block w-2 h-2 rounded-full bg-emerald-500"></span>
          <span class="font-semibold text-zinc-800">Formulir Rekap Slip Gaji (Format Jurnal Umum 16 Kolom)</span>
          <span class="text-zinc-400">|</span>
          <span class="text-zinc-500">{{ entityFilteredPayslips.length }} baris data</span>
        </div>
        <div class="flex items-center gap-1.5 text-zinc-400 text-[11px] font-medium">
          <span>Geser horizontal untuk melihat rincian gaji & potongan</span>
          <span class="text-zinc-500 font-bold">→</span>
        </div>
      </div>

      <!-- High-Density Data Grid (Jurnal Umum Style) -->
      <div v-if="entityFilteredPayslips.length > 0" class="overflow-x-auto relative">
        <table class="w-full min-w-[2050px] text-left text-xs border-collapse">
          <!-- Table Header: Jurnal Umum Look -->
          <thead class="bg-zinc-100 border-b-2 border-zinc-300 text-zinc-800 uppercase text-[10px] tracking-wider font-extrabold sticky top-0 z-10 shadow-xs">
            <tr>
              <th class="px-3 py-3 border-r border-zinc-200 text-center w-12 min-w-[48px] whitespace-nowrap">No</th>
              <th class="px-3 py-3 border-r border-zinc-200 w-[145px] min-w-[145px] whitespace-nowrap">ID / NIK</th>
              <th class="px-4 py-3 border-r border-zinc-200 w-[220px] min-w-[220px] whitespace-nowrap">Nama Karyawan</th>
              <th class="px-3 py-3 border-r border-zinc-200 w-[170px] min-w-[170px] whitespace-nowrap">Entitas & Unit</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center w-[85px] min-w-[85px] whitespace-nowrap">Periode</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center bg-emerald-50/80 text-emerald-950 w-[80px] min-w-[80px] whitespace-nowrap">Hadir</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center bg-rose-50/80 text-rose-950 w-[75px] min-w-[75px] whitespace-nowrap">Alpa</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center w-[80px] min-w-[80px] whitespace-nowrap">Izin/Cuti</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center w-[80px] min-w-[80px] whitespace-nowrap">Lembur</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right w-[125px] min-w-[125px] whitespace-nowrap">Nilai Lembur</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right w-[135px] min-w-[135px] whitespace-nowrap">Gaji Pokok</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right w-[135px] min-w-[135px] whitespace-nowrap">Tunjangan</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right text-rose-700 w-[130px] min-w-[130px] whitespace-nowrap">Potongan</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right text-amber-800 w-[130px] min-w-[130px] whitespace-nowrap">Bonus/Insentif</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right font-black w-[140px] min-w-[140px] whitespace-nowrap">Gross</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-right bg-emerald-100/70 font-black text-emerald-950 w-[155px] min-w-[155px] whitespace-nowrap">Gaji Bersih (THP)</th>
              <th class="px-3 py-3 border-r border-zinc-200 text-center w-[110px] min-w-[110px] whitespace-nowrap">Status</th>
              <th class="px-3 py-3 text-center w-[130px] min-w-[130px] whitespace-nowrap">Aksi</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-zinc-200 font-medium">
            <tr
              v-for="(slip, idx) in entityFilteredPayslips"
              :key="slip.id"
              class="hover:bg-zinc-50/90 transition-colors"
              :class="idx % 2 === 0 ? 'bg-white' : 'bg-zinc-50/40'"
            >
              <!-- 0. Index -->
              <td class="px-3 py-2.5 text-center border-r border-zinc-100 text-zinc-400 font-mono text-[11px] whitespace-nowrap">
                {{ idx + 1 }}
              </td>

              <!-- 1. ID Karyawan / NIK -->
              <td class="px-3 py-2.5 border-r border-zinc-100 font-mono text-[11px] font-bold text-zinc-900 whitespace-nowrap">
                {{ slip.employeeNik }}
              </td>

              <!-- 2. Nama Karyawan -->
              <td class="px-4 py-2.5 border-r border-zinc-100 whitespace-nowrap">
                <span class="font-bold text-zinc-950 block leading-tight">{{ slip.employeeName }}</span>
                <span class="text-[10px] text-zinc-500 block mt-0.5 truncate max-w-[200px]">{{ slip.position }}</span>
              </td>

              <!-- 2b. Entitas & Unit Bisnis -->
              <td class="px-3 py-2.5 border-r border-zinc-100 whitespace-nowrap">
                <div class="text-[11px] font-bold" :class="isClinicEntity(slip) ? 'text-blue-900' : 'text-emerald-950'">
                  {{ isClinicEntity(slip) ? 'CV SM Clinic' : 'CV SM Group' }}
                </div>
                <div class="inline-flex items-center gap-1 mt-0.5">
                  <span
                    class="px-1.5 py-0.5 rounded text-[9px] font-bold uppercase tracking-wider border"
                    :class="getUnitBadgeClass(slip.businessUnit)"
                  >
                    {{ getUnitLabel(slip.businessUnit) }}
                  </span>
                  <span class="text-[9px] text-zinc-400 font-mono">• {{ slip.department }}</span>
                </div>
              </td>

              <!-- 3. Periode Gaji -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center font-mono text-[11px] text-zinc-600 whitespace-nowrap">
                {{ slip.periodMonth }}/{{ slip.periodYear }}
              </td>

              <!-- 4. Jumlah Hari Hadir -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center font-bold bg-emerald-50/40 text-emerald-800 whitespace-nowrap">
                {{ slip.presentDays !== undefined ? slip.presentDays : (slip.workingDays || 22) }} hr
              </td>

              <!-- 5. Jumlah Hari Alpa -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center font-bold whitespace-nowrap" :class="(slip.absentDays || 0) > 0 ? 'text-rose-600 bg-rose-50/50 font-black' : 'text-zinc-400'">
                {{ slip.absentDays || 0 }} hr
              </td>

              <!-- 6. Jumlah Hari Izin / Sakit / Cuti -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center text-zinc-600 whitespace-nowrap">
                {{ (slip.paidLeaveDays || 0) + (slip.sickLeaveDays || 0) }} hr
              </td>

              <!-- 7. Jumlah Jam Lembur -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center font-bold whitespace-nowrap" :class="(slip.overtimeHours || 0) > 0 ? 'text-indigo-600' : 'text-zinc-400'">
                {{ slip.overtimeHours ? `+${slip.overtimeHours}j` : '—' }}
              </td>

              <!-- 8. Nilai Lembur -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono text-zinc-700 tabular-nums whitespace-nowrap">
                {{ formatIDR(slip.overtimePay || 0) }}
              </td>

              <!-- 9. Gaji Pokok -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono text-zinc-950 font-bold tabular-nums whitespace-nowrap">
                {{ formatIDR(slip.baseSalary) }}
              </td>

              <!-- 10. Tunjangan (Jabatan, Makan, Transport) -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono text-emerald-800 font-semibold tabular-nums whitespace-nowrap">
                +{{ formatIDR(slip.positionAllowance + slip.transportMealAllowance) }}
              </td>

              <!-- 11. Potongan (BPJS, Kasbon, Alpa) -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono text-rose-600 font-semibold tabular-nums whitespace-nowrap">
                -{{ formatIDR(slip.totalDeduction) }}
              </td>

              <!-- 12. Bonus / Insentif -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono text-amber-800 tabular-nums whitespace-nowrap">
                {{ formatIDR((slip.shiftAttendanceAllowance || 0) + (slip.incentiveBonus || 0)) }}
              </td>

              <!-- 13. Total Pendapatan (Gross) -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono font-black text-zinc-950 tabular-nums whitespace-nowrap">
                {{ formatIDR(slip.totalGrossSalary) }}
              </td>

              <!-- 14. Gaji Bersih (Take Home Pay / THP) -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-right font-mono font-black text-[13px] text-emerald-950 bg-emerald-50/60 tabular-nums whitespace-nowrap">
                {{ formatIDR(slip.takeHomePay) }}
              </td>

              <!-- 15. Status (Draft, Disetujui, Dibayar) -->
              <td class="px-3 py-2.5 border-r border-zinc-100 text-center whitespace-nowrap">
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold border tracking-wide"
                  :class="{
                    'bg-zinc-100 text-zinc-700 border-zinc-200': slip.status === 'DRAFT',
                    'bg-amber-50 text-amber-800 border-amber-200': slip.status === 'APPROVED',
                    'bg-emerald-50 text-emerald-800 border-emerald-200': slip.status === 'PAID'
                  }"
                >
                  {{ slip.status }}
                </span>
              </td>

              <!-- 16. Aksi (Detail, Edit, Cetak, Hapus) -->
              <td class="px-3 py-2.5 text-center whitespace-nowrap">
                <div class="flex items-center justify-center gap-1">
                  <!-- Detail / Print Slip -->
                  <button
                    @click="viewPayslip(slip)"
                    class="p-1.5 text-zinc-700 hover:text-zinc-950 hover:bg-zinc-100 rounded-lg transition-colors cursor-pointer"
                    title="Lihat & Cetak Slip PDF"
                  >
                    <Eye class="w-3.5 h-3.5" />
                  </button>

                  <!-- Edit Nilai Komponen -->
                  <button
                    @click="openEditModal(slip)"
                    class="p-1.5 text-zinc-500 hover:text-zinc-800 hover:bg-zinc-100 rounded-lg transition-colors cursor-pointer"
                    title="Ubah Nilai Komponen Gaji"
                  >
                    <Edit class="w-3.5 h-3.5" />
                  </button>

                  <!-- Approve button (if draft) -->
                  <button
                    v-if="slip.status === 'DRAFT'"
                    @click="payrollStore.approvePayslip(slip.id)"
                    class="p-1.5 text-amber-600 hover:text-amber-800 hover:bg-amber-50 rounded-lg transition-colors cursor-pointer"
                    title="Setujui Slip Gaji"
                  >
                    <CheckCircle2 class="w-3.5 h-3.5" />
                  </button>

                  <!-- Mark Paid button (if approved) -->
                  <button
                    v-if="slip.status === 'APPROVED'"
                    @click="payrollStore.markAsPaid(slip.id)"
                    class="p-1.5 text-emerald-600 hover:text-emerald-800 hover:bg-emerald-50 rounded-lg transition-colors cursor-pointer"
                    title="Tandai Telah Ditransfer (Paid)"
                  >
                    <Send class="w-3.5 h-3.5" />
                  </button>

                  <!-- Delete (Always available with confirmation) -->
                  <button
                    @click="handleDeleteSlip(slip)"
                    class="p-1.5 text-zinc-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                    :title="`Hapus Slip Gaji ${slip.employeeName}`"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>

          <!-- Table Footer: Totals Row (Jurnal Umum Style) -->
          <tfoot class="bg-zinc-100 border-t-2 border-zinc-400 font-extrabold text-xs text-zinc-900 sticky bottom-0 z-10 shadow-xs">
            <tr>
              <td colspan="9" class="px-4 py-3 border-r border-zinc-200 text-right uppercase tracking-wider text-[11px] whitespace-nowrap">
                TOTAL KESELURUHAN ({{ entityStats.totalCount }} KARYAWAN):
              </td>
              <!-- Total Lembur -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums whitespace-nowrap">
                {{ formatIDR(entityStats.totalOvertimePay) }}
              </td>
              <!-- Total Gaji Pokok -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums whitespace-nowrap">
                {{ formatIDR(entityStats.totalBaseSalary) }}
              </td>
              <!-- Total Tunjangan -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums text-emerald-800 whitespace-nowrap">
                +{{ formatIDR(entityStats.totalAllowances) }}
              </td>
              <!-- Total Potongan -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums text-rose-600 whitespace-nowrap">
                -{{ formatIDR(entityStats.totalDeductions) }}
              </td>
              <!-- Total Bonus -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums text-amber-800 whitespace-nowrap">
                {{ formatIDR(entityStats.totalBonus) }}
              </td>
              <!-- Total Gross -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums font-black whitespace-nowrap">
                {{ formatIDR(entityStats.totalGross) }}
              </td>
              <!-- Total THP -->
              <td class="px-3 py-3 border-r border-zinc-200 text-right font-mono tabular-nums font-black text-[13px] text-emerald-950 bg-emerald-100/70 whitespace-nowrap">
                {{ formatIDR(entityStats.totalTakeHomePay) }}
              </td>
              <!-- Status & Aksi -->
              <td colspan="2" class="px-3 py-3 text-center text-[10px] text-zinc-500 font-semibold whitespace-nowrap">
                {{ entityStats.paidCount }} Paid • {{ entityStats.approvedCount }} Approved
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- Modals -->
    <PayslipDetailModal
      :show="showDetailModal"
      :slip="selectedSlip"
      @close="showDetailModal = false; selectedSlip = null"
    />

    <PayslipFormModal
      :show="showFormModal"
      :edit-slip="editingSlip"
      @close="showFormModal = false; editingSlip = null"
    />

    <AttendanceImportModal
      :show="showImportModal"
      :employees="payrollStore.employees"
      :selected-month="payrollStore.selectedMonth"
      :selected-year="payrollStore.selectedYear"
      @close="showImportModal = false"
      @imported="payrollStore.fetchPayslips()"
    />

  </div>
</template>
