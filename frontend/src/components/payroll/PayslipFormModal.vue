<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { usePayrollStore } from '@/stores/payroll'
import { useLocaleStore } from '@/stores/locale'
import { formatIDR } from '@/utils/currency'
import type { PayslipItem, Employee } from '@/types/payroll'
import type { BusinessUnit } from '@/types/accounting'
import {
  X,
  Plus,
  Save,
  User,
  Calculator,
  Building2,
  Calendar,
  AlertCircle,
  Search,
  ChevronDown,
  Check
} from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
  editSlip?: PayslipItem | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const payrollStore = usePayrollStore()
const localeStore = useLocaleStore()

const dropdownContainerRef = ref<HTMLElement | null>(null)
const isDropdownOpen = ref<boolean>(false)
const selectedEmployeeId = ref<string>('')
const searchInputText = ref<string>('')
const employeeName = ref<string>('')
const employeeNik = ref<string>('')
const position = ref<string>('')
const department = ref<string>('')
const businessUnit = ref<BusinessUnit>('KIDINEST')
const periodMonth = ref<number>(new Date().getMonth() + 1)
const periodYear = ref<number>(new Date().getFullYear())
const bankName = ref<string>('Bank Mandiri')
const bankAccountNumber = ref<string>('')

const selectedEmployee = computed(() => {
  if (!selectedEmployeeId.value) return null
  return payrollStore.employees.find(e => e.id === selectedEmployeeId.value) || null
})

// Display all employees when dropdown is open, or filter in real time when typing
const filteredEmployees = computed(() => {
  const q = searchInputText.value.trim().toLowerCase()
  if (!q || (selectedEmployee.value && searchInputText.value === `${selectedEmployee.value.name} (${selectedEmployee.value.nik})`)) {
    return payrollStore.employees
  }
  return payrollStore.employees.filter(emp =>
    emp.name.toLowerCase().includes(q) ||
    emp.nik.toLowerCase().includes(q) ||
    emp.position.toLowerCase().includes(q) ||
    emp.department.toLowerCase().includes(q) ||
    emp.businessUnit.toLowerCase().includes(q)
  )
})

function selectEmployee(emp: Employee) {
  selectedEmployeeId.value = emp.id
  employeeName.value = emp.name
  employeeNik.value = emp.nik
  position.value = emp.position
  department.value = emp.department
  businessUnit.value = emp.businessUnit
  bankName.value = emp.bankName
  bankAccountNumber.value = emp.bankAccountNumber
  baseSalary.value = emp.baseSalary
  positionAllowance.value = emp.positionAllowance
  
  // Kehadiran defaults
  workingDays.value = 22
  presentDays.value = 22
  sickLeaveDays.value = 0
  paidLeaveDays.value = 0
  absentDays.value = 0
  overtimeHours.value = 0
  dailyMealTransportRate.value = emp.dailyMealTransportRate || (emp.transportMealAllowance ? Math.round(emp.transportMealAllowance / 22) : 30000)
  transportMealAllowance.value = Math.round(presentDays.value * dailyMealTransportRate.value)
  attendanceDeduction.value = 0
  overtimePay.value = 0

  bpjsKesehatan.value = Math.round(emp.baseSalary * 0.01)
  bpjsKetenagakerjaan.value = Math.round(emp.baseSalary * 0.03)

  searchInputText.value = `${emp.name} (${emp.nik})`
  isDropdownOpen.value = false
}

function selectManualMode() {
  selectedEmployeeId.value = ''
  searchInputText.value = ''
  isDropdownOpen.value = false
}

function handleInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  searchInputText.value = val
  isDropdownOpen.value = true
  if (!selectedEmployeeId.value) {
    employeeName.value = val
  }
}

function toggleDropdown() {
  isDropdownOpen.value = !isDropdownOpen.value
}

function clearEmployeeSelection() {
  selectedEmployeeId.value = ''
  searchInputText.value = ''
  isDropdownOpen.value = false
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownContainerRef.value && !dropdownContainerRef.value.contains(event.target as Node)) {
    isDropdownOpen.value = false
    if (selectedEmployee.value) {
      searchInputText.value = `${selectedEmployee.value.name} (${selectedEmployee.value.nik})`
    }
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

// Attendance & Working Days (Hari Masuk / Kehadiran)
const workingDays = ref<number>(22)
const presentDays = ref<number>(22)
const sickLeaveDays = ref<number>(0)
const paidLeaveDays = ref<number>(0)
const absentDays = ref<number>(0)
const overtimeHours = ref<number>(0)
const dailyMealTransportRate = ref<number>(30000)

// Earnings
const baseSalary = ref<number>(4500000)
const positionAllowance = ref<number>(500000)
const transportMealAllowance = ref<number>(660000)
const shiftAttendanceAllowance = ref<number>(0)
const incentiveBonus = ref<number>(0)
const overtimePay = ref<number>(0)

// Deductions
const bpjsKesehatan = ref<number>(45000)
const bpjsKetenagakerjaan = ref<number>(135000)
const pph21 = ref<number>(0)
const attendanceDeduction = ref<number>(0)
const loanDeduction = ref<number>(0)
const notes = ref<string>('')

// When employee selected, populate defaults
watch(selectedEmployeeId, (newId) => {
  if (!newId) return
  const emp = payrollStore.employees.find(e => e.id === newId)
  if (emp) {
    employeeName.value = emp.name
    employeeNik.value = emp.nik
    position.value = emp.position
    department.value = emp.department
    businessUnit.value = emp.businessUnit
    bankName.value = emp.bankName
    bankAccountNumber.value = emp.bankAccountNumber
    baseSalary.value = emp.baseSalary
    positionAllowance.value = emp.positionAllowance
    
    // Kehadiran defaults
    workingDays.value = 22
    presentDays.value = 22
    sickLeaveDays.value = 0
    paidLeaveDays.value = 0
    absentDays.value = 0
    overtimeHours.value = 0
    dailyMealTransportRate.value = emp.dailyMealTransportRate || (emp.transportMealAllowance ? Math.round(emp.transportMealAllowance / 22) : 30000)
    transportMealAllowance.value = Math.round(presentDays.value * dailyMealTransportRate.value)
    attendanceDeduction.value = 0
    overtimePay.value = 0

    bpjsKesehatan.value = Math.round(emp.baseSalary * 0.01)
    bpjsKetenagakerjaan.value = Math.round(emp.baseSalary * 0.03)
  }
})

// Auto recalculate Uang Makan & Transport based on Present Days & Rate
watch([presentDays, dailyMealTransportRate], ([pres, rate]) => {
  transportMealAllowance.value = Math.round((Number(pres) || 0) * (Number(rate) || 0))
})

// Auto recalculate Attendance Deduction if employee is absent / mangkir
watch([absentDays, workingDays, baseSalary], ([abs, wDays, base]) => {
  const numAbs = Number(abs) || 0
  const numWDays = Number(wDays) || 22
  const numBase = Number(base) || 0
  if (numAbs > 0 && numWDays > 0) {
    attendanceDeduction.value = Math.round(numAbs * (numBase / numWDays))
  } else {
    attendanceDeduction.value = 0
  }
})

// Auto recalculate Overtime Pay based on Overtime Hours
watch([overtimeHours, baseSalary, positionAllowance], ([ot, base, pos]) => {
  const numOt = Number(ot) || 0
  if (numOt > 0) {
    const hourly = ((Number(base || 0) + Number(pos || 0)) / 173) * 1.5
    overtimePay.value = Math.round(numOt * hourly)
  } else {
    overtimePay.value = 0
  }
})

// Auto-recalculate PPh 21 when gross changes
const totalGrossSalary = computed(() => {
  return (
    Number(baseSalary.value || 0) +
    Number(positionAllowance.value || 0) +
    Number(transportMealAllowance.value || 0) +
    Number(shiftAttendanceAllowance.value || 0) +
    Number(incentiveBonus.value || 0) +
    Number(overtimePay.value || 0)
  )
})

watch(totalGrossSalary, (gross) => {
  if (gross <= 5400000) pph21.value = 0
  else if (gross <= 5650000) pph21.value = Math.round(gross * 0.0025)
  else if (gross <= 5950000) pph21.value = Math.round(gross * 0.005)
  else if (gross <= 6300000) pph21.value = Math.round(gross * 0.0075)
  else if (gross <= 6750000) pph21.value = Math.round(gross * 0.01)
  else if (gross <= 7500000) pph21.value = Math.round(gross * 0.0125)
  else pph21.value = Math.round(gross * 0.015)
}, { immediate: true })

const totalDeduction = computed(() => {
  return (
    Number(bpjsKesehatan.value || 0) +
    Number(bpjsKetenagakerjaan.value || 0) +
    Number(pph21.value || 0) +
    Number(attendanceDeduction.value || 0) +
    Number(loanDeduction.value || 0)
  )
})

const takeHomePay = computed(() => {
  return Math.max(0, totalGrossSalary.value - totalDeduction.value)
})

// Watch editSlip props
watch(() => props.editSlip, (item) => {
  if (item) {
    selectedEmployeeId.value = item.employeeId
    searchInputText.value = `${item.employeeName} (${item.employeeNik})`
    employeeName.value = item.employeeName
    employeeNik.value = item.employeeNik
    position.value = item.position
    department.value = item.department
    businessUnit.value = item.businessUnit
    periodMonth.value = item.periodMonth
    periodYear.value = item.periodYear
    bankName.value = item.bankName
    bankAccountNumber.value = item.bankAccountNumber

    // Attendance
    workingDays.value = item.workingDays || 22
    presentDays.value = item.presentDays !== undefined ? item.presentDays : 22
    sickLeaveDays.value = item.sickLeaveDays || 0
    paidLeaveDays.value = item.paidLeaveDays || 0
    absentDays.value = item.absentDays || 0
    overtimeHours.value = item.overtimeHours || 0
    dailyMealTransportRate.value = item.dailyMealTransportRate || 30000

    baseSalary.value = item.baseSalary
    positionAllowance.value = item.positionAllowance
    transportMealAllowance.value = item.transportMealAllowance
    shiftAttendanceAllowance.value = item.shiftAttendanceAllowance
    incentiveBonus.value = item.incentiveBonus
    overtimePay.value = item.overtimePay
    bpjsKesehatan.value = item.bpjsKesehatan
    bpjsKetenagakerjaan.value = item.bpjsKetenagakerjaan
    pph21.value = item.pph21
    attendanceDeduction.value = item.attendanceDeduction
    loanDeduction.value = item.loanDeduction
    notes.value = item.notes || ''
  } else {
    selectedEmployeeId.value = ''
    searchInputText.value = ''
    employeeName.value = ''
    employeeNik.value = ''
    position.value = ''
    department.value = ''
    businessUnit.value = 'KIDINEST'
    workingDays.value = 22
    presentDays.value = 22
    sickLeaveDays.value = 0
    paidLeaveDays.value = 0
    absentDays.value = 0
    overtimeHours.value = 0
    dailyMealTransportRate.value = 30000
    baseSalary.value = 4500000
    positionAllowance.value = 500000
    transportMealAllowance.value = 660000
    shiftAttendanceAllowance.value = 0
    incentiveBonus.value = 0
    overtimePay.value = 0
    bpjsKesehatan.value = 45000
    bpjsKetenagakerjaan.value = 135000
    attendanceDeduction.value = 0
    loanDeduction.value = 0
    notes.value = ''
  }
}, { immediate: true })

function handleSave() {
  if (!employeeName.value.trim() || !employeeNik.value.trim()) {
    alert('Nama dan NIK karyawan wajib diisi')
    return
  }

  const monthNames = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
  const pName = `${monthNames[periodMonth.value - 1]} ${periodYear.value}`
  const pDate = `${periodYear.value}-${String(periodMonth.value).padStart(2, '0')}-28`

  const payload = {
    employeeId: selectedEmployeeId.value || `emp-${Date.now()}`,
    employeeNik: employeeNik.value,
    employeeName: employeeName.value,
    position: position.value,
    department: department.value,
    businessUnit: businessUnit.value,
    periodMonth: periodMonth.value,
    periodYear: periodYear.value,
    periodName: pName,
    paymentDate: pDate,
    workingDays: Number(workingDays.value) || 22,
    presentDays: Number(presentDays.value) || 0,
    sickLeaveDays: Number(sickLeaveDays.value) || 0,
    paidLeaveDays: Number(paidLeaveDays.value) || 0,
    absentDays: Number(absentDays.value) || 0,
    overtimeHours: Number(overtimeHours.value) || 0,
    dailyMealTransportRate: Number(dailyMealTransportRate.value) || 30000,
    baseSalary: baseSalary.value,
    positionAllowance: positionAllowance.value,
    transportMealAllowance: transportMealAllowance.value,
    shiftAttendanceAllowance: shiftAttendanceAllowance.value,
    incentiveBonus: incentiveBonus.value,
    overtimePay: overtimePay.value,
    totalGrossSalary: totalGrossSalary.value,
    bpjsKesehatan: bpjsKesehatan.value,
    bpjsKetenagakerjaan: bpjsKetenagakerjaan.value,
    pph21: pph21.value,
    attendanceDeduction: attendanceDeduction.value,
    loanDeduction: loanDeduction.value,
    totalDeduction: totalDeduction.value,
    takeHomePay: takeHomePay.value,
    status: (props.editSlip?.status || 'DRAFT') as any,
    bankName: bankName.value,
    bankAccountNumber: bankAccountNumber.value,
    notes: notes.value
  }

  if (props.editSlip) {
    payrollStore.updatePayslip(props.editSlip.id, payload)
  } else {
    payrollStore.addCustomPayslip(payload)
  }

  emit('saved')
  emit('close')
}
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs overflow-y-auto"
  >
    <div class="bg-white rounded-2xl border border-zinc-200 shadow-2xl w-full max-w-2xl my-auto overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-zinc-200 bg-zinc-50/80">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-zinc-950 text-white flex items-center justify-center font-bold text-xs shadow-xs">
            <Calculator class="w-4 h-4" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-zinc-950">
              {{ editSlip ? (localeStore.isEnglish ? 'Edit Staff Payslip' : 'Edit Slip Gaji Karyawan') : (localeStore.isEnglish ? 'Create New Payslip' : 'Tambah Slip Gaji Karyawan') }}
            </h3>
            <p class="text-xs text-zinc-500">Semarang Medical Group (Kidinest & SSC)</p>
          </div>
        </div>

        <button
          @click="emit('close')"
          class="p-1.5 text-zinc-400 hover:text-zinc-700 hover:bg-zinc-200/60 rounded-lg transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Form Body -->
      <form @submit.prevent="handleSave" class="p-6 space-y-5 max-h-[75vh] overflow-y-auto text-xs font-['Plus_Jakarta_Sans',sans-serif]">
        <!-- Employee Searchable Dropdown (Bisa Ketik & Tetap Muncul Dropdown) -->
        <div ref="dropdownContainerRef" class="relative">
          <div class="flex items-center justify-between mb-1">
            <label class="block font-bold text-zinc-700">
              Pilih / Ketik Karyawan *
            </label>
            <span class="text-[10px] text-zinc-500 font-medium">Bisa klik dropdown atau ketik nama/NIK untuk mencari</span>
          </div>

          <!-- Combobox Input Bar -->
          <div class="relative flex items-center">
            <Search class="w-4 h-4 text-zinc-400 absolute left-3 pointer-events-none" />
            <input
              type="text"
              :value="searchInputText"
              @input="handleInput"
              @focus="isDropdownOpen = true"
              class="w-full pl-9 pr-16 py-2.5 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 focus:outline-hidden text-xs font-medium bg-white hover:border-zinc-400 transition-colors cursor-text"
              placeholder="-- Pilih karyawan dari dropdown atau ketik nama/NIK --"
            />
            
            <div class="absolute right-2 flex items-center gap-1">
              <!-- Clear Button -->
              <button
                v-if="searchInputText"
                type="button"
                @click.stop="clearEmployeeSelection"
                class="p-1 text-zinc-400 hover:text-zinc-600 rounded-md transition-colors cursor-pointer"
                title="Hapus / Reset Pilihan"
              >
                <X class="w-3.5 h-3.5" />
              </button>

              <!-- Dropdown Trigger Chevron -->
              <button
                type="button"
                @click.stop="toggleDropdown"
                class="p-1 text-zinc-500 hover:text-zinc-900 hover:bg-zinc-100 rounded-md transition-colors cursor-pointer"
                title="Buka / Tutup Daftar Karyawan"
              >
                <ChevronDown
                  class="w-4 h-4 transition-transform duration-150"
                  :class="isDropdownOpen ? 'rotate-180' : ''"
                />
              </button>
            </div>
          </div>

          <!-- Dropdown Menu List -->
          <div
            v-if="isDropdownOpen"
            class="absolute left-0 right-0 top-full mt-1.5 bg-white border border-zinc-200 rounded-xl shadow-2xl z-40 max-h-64 overflow-y-auto divide-y divide-zinc-100 animate-in fade-in zoom-in-95 duration-100"
          >
            <!-- Option 0: Manual Input -->
            <div
              @click="selectManualMode"
              class="px-3.5 py-2.5 hover:bg-zinc-50 cursor-pointer transition-colors flex items-center justify-between text-xs"
              :class="!selectedEmployeeId ? 'bg-zinc-50/80 font-bold text-zinc-900' : 'text-zinc-600'"
            >
              <div class="flex items-center gap-2">
                <User class="w-4 h-4 text-zinc-400" />
                <div>
                  <div class="font-bold text-zinc-800">-- Input Karyawan Baru / Ketik Manual --</div>
                  <div class="text-[10px] text-zinc-400">Isi nama, NIK, dan jabatan secara manual di formulir bawah</div>
                </div>
              </div>
              <span v-if="!selectedEmployeeId" class="text-zinc-500 text-[10px] font-bold">Manual</span>
            </div>

            <!-- Options: Filtered Employees -->
            <div
              v-for="emp in filteredEmployees"
              :key="emp.id"
              @click="selectEmployee(emp)"
              class="px-3.5 py-2.5 hover:bg-zinc-50 cursor-pointer transition-colors flex items-center justify-between text-xs group"
              :class="selectedEmployeeId === emp.id ? 'bg-emerald-50/60 font-bold' : ''"
            >
              <div>
                <div class="font-bold text-zinc-900 group-hover:text-zinc-950 flex items-center gap-1.5">
                  <span>{{ emp.name }}</span>
                  <Check v-if="selectedEmployeeId === emp.id" class="w-3.5 h-3.5 text-emerald-600 inline" />
                </div>
                <div class="text-[10px] text-zinc-500 font-mono mt-0.5">
                  {{ emp.nik }} • {{ emp.position }} ({{ emp.department }})
                </div>
              </div>

              <!-- Entity & Unit Badge in dropdown -->
              <div class="text-right shrink-0 ml-3">
                <span
                  class="px-2 py-0.5 rounded text-[9px] font-bold uppercase tracking-wider border block"
                  :class="emp.businessUnit === 'APOTEK' ? 'bg-blue-50 text-blue-800 border-blue-200' : emp.businessUnit === 'SSC' ? 'bg-rose-50 text-rose-800 border-rose-200' : 'bg-emerald-50 text-emerald-800 border-emerald-200'"
                >
                  {{ emp.businessUnit === 'APOTEK' ? 'CV SMC • Apotek' : emp.businessUnit === 'SSC' ? 'CV SMG • SSC' : emp.businessUnit === 'KIDINEST' ? 'CV SMG • Kidinest' : 'CV SMG • Pusat' }}
                </span>
                <span class="text-[9px] font-mono text-zinc-400 block mt-0.5">{{ formatIDR(emp.baseSalary) }}</span>
              </div>
            </div>

            <!-- Empty result -->
            <div v-if="filteredEmployees.length === 0" class="p-4 text-center text-zinc-400 text-xs italic">
              Tidak ada karyawan yang cocok dengan pencarian "{{ searchInputText }}". Anda dapat mengisi manual formulir di bawah.
            </div>
          </div>
        </div>

        <!-- Personal & Period Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">{{ localeStore.isEnglish ? 'Employee Name' : 'Nama Lengkap Karyawan' }} *</label>
            <input
              v-model="employeeName"
              type="text"
              required
              class="w-full px-3 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 focus:outline-hidden text-xs font-medium"
              placeholder="cth: Fitri Handayani"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">NIK / ID Karyawan *</label>
            <input
              v-model="employeeNik"
              type="text"
              required
              class="w-full px-3 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 focus:outline-hidden text-xs font-mono font-medium"
              placeholder="cth: SMG-KID-2024-001"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">{{ localeStore.isEnglish ? 'Position / Role' : 'Jabatan / Tugas' }} *</label>
            <input
              v-model="position"
              type="text"
              required
              class="w-full px-3 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 focus:outline-hidden text-xs font-medium"
              placeholder="cth: Kasir Shift GRO / Perawat Poli Anak"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Entitas Perusahaan & Unit Bisnis *</label>
            <select
              v-model="businessUnit"
              class="w-full px-3 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 focus:outline-hidden text-xs font-medium bg-white"
            >
              <optgroup label="CV SEMARANG MEDICAL GROUP (SMG)">
                <option value="KIDINEST">CV SMG - Kidinest (Poli Spesialis Anak)</option>
                <option value="SSC">CV SMG - Semarang Skin Center (SSC)</option>
                <option value="ALL">CV SMG - Kantor Pusat / Manajemen Konsolidasi</option>
              </optgroup>
              <optgroup label="CV SEMARANG MEDICAL CLINIC (SMC)">
                <option value="APOTEK">CV SMC - Instalasi Farmasi & Apotek</option>
              </optgroup>
            </select>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">{{ localeStore.isEnglish ? 'Payroll Month & Year' : 'Bulan & Tahun Penggajian' }} *</label>
            <div class="grid grid-cols-2 gap-2">
              <select
                v-model="periodMonth"
                class="w-full px-2.5 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 text-xs font-medium"
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
              <input
                v-model.number="periodYear"
                type="number"
                class="w-full px-3 py-2 border border-zinc-300 rounded-xl focus:ring-2 focus:ring-zinc-900 text-xs font-mono font-medium"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">{{ localeStore.isEnglish ? 'Bank & Account No' : 'Bank & Nomor Rekening' }}</label>
            <div class="grid grid-cols-2 gap-2">
              <input
                v-model="bankName"
                type="text"
                class="w-full px-3 py-2 border border-zinc-300 rounded-xl text-xs font-medium"
                placeholder="Bank Mandiri"
              />
              <input
                v-model="bankAccountNumber"
                type="text"
                class="w-full px-3 py-2 border border-zinc-300 rounded-xl text-xs font-mono font-medium"
                placeholder="136-00-..."
              />
            </div>
          </div>
        </div>

        <!-- Section: Attendance & Kehadiran / Hari Masuk Kerja -->
        <div class="p-4 rounded-xl border border-sky-200 bg-sky-50/40 space-y-3">
          <div class="font-bold text-xs text-sky-950 uppercase tracking-wider flex items-center justify-between">
            <span>{{ localeStore.isEnglish ? 'WORKING DAYS & ATTENDANCE RECORD' : 'REKAPITULASI HARI MASUK & KEHADIRAN KERJA' }}</span>
            <span class="text-[11px] font-semibold text-sky-800">
              {{ localeStore.isEnglish ? 'Present' : 'Masuk' }}: {{ presentDays }} / {{ workingDays }} {{ localeStore.isEnglish ? 'Days' : 'Hari' }}
            </span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">
                {{ localeStore.isEnglish ? 'Standard Working Days' : 'Standar Hari Kerja (Bulan)' }} *
              </label>
              <input
                v-model.number="workingDays"
                type="number"
                min="1"
                max="31"
                class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-mono font-medium bg-white"
              />
            </div>

            <div>
              <label class="block font-semibold text-emerald-800 mb-0.5">
                {{ localeStore.isEnglish ? 'Present / In-Office Days' : 'Hari Masuk / Hadir Kerja' }} *
              </label>
              <input
                v-model.number="presentDays"
                type="number"
                min="0"
                max="31"
                class="w-full px-3 py-1.5 border border-emerald-300 bg-white rounded-lg text-xs font-mono font-bold text-emerald-900"
              />
              <span class="text-[10px] text-zinc-500 mt-0.5 block">Mempengaruhi uang makan</span>
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">
                {{ localeStore.isEnglish ? 'Meal & Transp. Rate/Day' : 'Tarif Makan & Trp / Hari' }}
              </label>
              <input
                v-model.number="dailyMealTransportRate"
                type="number"
                step="5000"
                class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-mono font-medium bg-white"
              />
              <span class="text-[10px] text-zinc-500 mt-0.5 block">Standar: Rp 30.000/hr</span>
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">
                {{ localeStore.isEnglish ? 'Paid Leave / Permit (Days)' : 'Izin Resmi / Cuti (Hari)' }}
              </label>
              <input
                v-model.number="paidLeaveDays"
                type="number"
                min="0"
                class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-mono font-medium bg-white"
              />
            </div>

            <div>
              <label class="block font-semibold text-amber-800 mb-0.5">
                {{ localeStore.isEnglish ? 'Sick Leave (Doctor Note)' : 'Sakit - Surat Dokter (Hari)' }}
              </label>
              <input
                v-model.number="sickLeaveDays"
                type="number"
                min="0"
                class="w-full px-3 py-1.5 border border-amber-300 rounded-lg text-xs font-mono font-medium bg-white"
              />
            </div>

            <div>
              <label class="block font-semibold text-rose-800 mb-0.5">
                {{ localeStore.isEnglish ? 'Absent / Unexcused (Days)' : 'Alpa / Mangkir (Hari)' }}
              </label>
              <input
                v-model.number="absentDays"
                type="number"
                min="0"
                class="w-full px-3 py-1.5 border border-rose-300 rounded-lg text-xs font-mono font-bold text-rose-900 bg-white"
              />
              <span class="text-[10px] text-rose-600 mt-0.5 block">Potong gaji proporsional</span>
            </div>

            <div>
              <label class="block font-semibold text-indigo-800 mb-0.5">
                {{ localeStore.isEnglish ? 'Overtime Hours' : 'Jam Kerja Lembur (Jam)' }}
              </label>
              <input
                v-model.number="overtimeHours"
                type="number"
                step="0.5"
                min="0"
                class="w-full px-3 py-1.5 border border-indigo-300 rounded-lg text-xs font-mono font-medium bg-white"
              />
              <span class="text-[10px] text-zinc-500 mt-0.5 block">Formula 1.5× upah/jam</span>
            </div>
          </div>
        </div>

        <!-- Section Earnings (Penerimaan) -->
        <div class="p-4 rounded-xl border border-emerald-200 bg-emerald-50/30 space-y-3">
          <div class="font-bold text-xs text-emerald-950 uppercase tracking-wider flex items-center justify-between">
            <span>A. {{ localeStore.isEnglish ? 'EARNINGS & ALLOWANCES' : 'KOMPONEN PENERIMAAN & TUNJANGAN' }}</span>
            <span class="tabular-nums font-mono text-emerald-800">{{ formatIDR(totalGrossSalary) }}</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Basic Salary' : 'Gaji Pokok' }}</label>
              <input v-model.number="baseSalary" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Position Allowance' : 'Tunjangan Jabatan' }}</label>
              <input v-model.number="positionAllowance" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Meal & Transport' : 'Uang Makan & Transport' }}</label>
              <input v-model.number="transportMealAllowance" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Shift & Attendance' : 'Tunjangan Shift Malam/Hadir' }}</label>
              <input v-model.number="shiftAttendanceAllowance" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Incentive / Bonus' : 'Insentif / Bonus Kinerja' }}</label>
              <input v-model.number="incentiveBonus" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Overtime Pay' : 'Uang Lembur (Overtime)' }}</label>
              <input v-model.number="overtimePay" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>
          </div>
        </div>

        <!-- Section Deductions (Potongan) -->
        <div class="p-4 rounded-xl border border-rose-200 bg-rose-50/30 space-y-3">
          <div class="font-bold text-xs text-rose-950 uppercase tracking-wider flex items-center justify-between">
            <span>B. {{ localeStore.isEnglish ? 'DEDUCTIONS & WITHHOLDINGS' : 'KOMPONEN POTONGAN & PAJAK' }}</span>
            <span class="tabular-nums font-mono text-rose-800">{{ formatIDR(totalDeduction) }}</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">BPJS Kesehatan (1%)</label>
              <input v-model.number="bpjsKesehatan" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">BPJS TK (JHT & JP 3%)</label>
              <input v-model.number="bpjsKetenagakerjaan" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">PPh Pasal 21 (TER)</label>
              <input v-model.number="pph21" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Attendance Penalty' : 'Potongan Terlambat / Absen' }}</label>
              <input v-model.number="attendanceDeduction" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Loan / Kasbon' : 'Cicilan Kasbon / Pinjaman' }}</label>
              <input v-model.number="loanDeduction" type="number" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" />
            </div>

            <div>
              <label class="block font-semibold text-zinc-700 mb-0.5">{{ localeStore.isEnglish ? 'Notes' : 'Catatan Khusus' }}</label>
              <input v-model="notes" type="text" class="w-full px-3 py-1.5 border border-zinc-300 rounded-lg text-xs font-medium" placeholder="cth: Lembur audit" />
            </div>
          </div>
        </div>

        <!-- Take Home Pay Preview Strip -->
        <div class="p-3.5 rounded-xl bg-zinc-950 text-white flex items-center justify-between">
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-emerald-400 block">
              {{ localeStore.isEnglish ? 'NET TAKE HOME PAY' : 'GAJI BERSIH DITERIMA' }}
            </span>
            <span class="text-xs text-zinc-400">Total Penghasilan dikurangi Seluruh Potongan</span>
          </div>
          <span class="text-xl font-black tabular-nums text-emerald-400">
            {{ formatIDR(takeHomePay) }}
          </span>
        </div>

        <!-- Footer Actions -->
        <div class="flex items-center justify-end gap-2.5 pt-3 border-t border-zinc-200">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 bg-white hover:bg-zinc-100 text-zinc-700 border border-zinc-200 rounded-xl text-xs font-bold transition-colors cursor-pointer"
          >
            {{ localeStore.isEnglish ? 'Cancel' : 'Batal' }}
          </button>

          <button
            type="submit"
            class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95 cursor-pointer"
          >
            <Save class="w-3.5 h-3.5" />
            <span>{{ editSlip ? (localeStore.isEnglish ? 'Update Payslip' : 'Simpan Perubahan') : (localeStore.isEnglish ? 'Save Payslip' : 'Simpan Slip Gaji') }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
