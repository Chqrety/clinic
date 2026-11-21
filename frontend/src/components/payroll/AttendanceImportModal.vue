<script setup lang="ts">
import { ref, computed } from 'vue'
import * as XLSX from 'xlsx'
import type { Employee, PayslipItem, PayrollStatus } from '@/types/payroll'
import { formatIDR } from '@/utils/currency'
import { useAccountingStore } from '@/stores/accounting'
import { usePayrollStore } from '@/stores/payroll'
import {
  X,
  Upload,
  FileSpreadsheet,
  CheckCircle2,
  AlertTriangle,
  Download,
  ArrowRight,
  RefreshCw,
  Table,
  Check,
  AlertCircle
} from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
  employees: Employee[]
  selectedMonth: number
  selectedYear: number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'imported', count: number): void
}>()

const accountingStore = useAccountingStore()
const payrollStore = usePayrollStore()

// State
const fileInputRef = ref<HTMLInputElement | null>(null)
const fileName = ref('')
const sheetNames = ref<string[]>([])
const selectedSheet = ref('')
const rawRows = ref<any[][]>([])
const isProcessing = ref(false)
const isSaving = ref(false)

// Column mapping indices (-1 if not mapped)
const colNIK = ref<number>(-1)
const colName = ref<number>(-1)
const colHadir = ref<number>(-1)
const colAlpa = ref<number>(-1)
const colIzin = ref<number>(-1)
const colLembur = ref<number>(-1)

interface ParsedAttendanceRow {
  rowIndex: number
  rawNIK: string
  rawName: string
  hadir: number
  alpa: number
  izin: number
  lemburHours: number
  matchedEmployee: Employee | null
  isValid: boolean
  validationMessage: string
  calculatedValues?: {
    overtimePay: number
    transportMealAllowance: number
    attendanceDeduction: number
    totalGrossSalary: number
    totalDeduction: number
    takeHomePay: number
  }
}

const parsedRows = ref<ParsedAttendanceRow[]>([])

const headerOptions = computed(() => {
  if (rawRows.value.length === 0) return []
  const firstRow = rawRows.value[0] || []
  return firstRow.map((val, idx) => ({
    index: idx,
    label: `Kolom ${String.fromCharCode(65 + idx)}: ${String(val || '').trim() || '(Kosong)'}`
  }))
})

const summaryStats = computed(() => {
  const total = parsedRows.value.length
  const valid = parsedRows.value.filter(r => r.isValid).length
  const invalid = total - valid
  return { total, valid, invalid }
})

function triggerFileInput() {
  fileInputRef.value?.click()
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  processFile(file)
}

function handleFileDrop(e: DragEvent) {
  e.preventDefault()
  const file = e.dataTransfer?.files?.[0]
  if (!file) return
  processFile(file)
}

function processFile(file: File) {
  fileName.value = file.name
  isProcessing.value = true

  const reader = new FileReader()
  reader.onload = (evt) => {
    try {
      const data = new Uint8Array(evt.target?.result as ArrayBuffer)
      const workbook = XLSX.read(data, { type: 'array' })
      sheetNames.value = workbook.SheetNames
      if (workbook.SheetNames.length > 0) {
        selectedSheet.value = workbook.SheetNames[0]
        loadSheetData(workbook, selectedSheet.value)
      }
    } catch (err: any) {
      accountingStore.showToast('Gagal membaca file Excel: ' + (err.message || 'Format tidak didukung'), 'error')
    } finally {
      isProcessing.value = false
    }
  }
  reader.readAsArrayBuffer(file)
}

function handleSheetChange() {
  if (!fileName.value || !selectedSheet.value) return
  // re-read workbook if needed or user switched sheet
}

function loadSheetData(workbook: XLSX.WorkBook, sheetName: string) {
  const sheet = workbook.Sheets[sheetName]
  if (!sheet) return

  const rows: any[][] = XLSX.utils.sheet_to_json(sheet, { header: 1, defval: '' })
  if (rows.length === 0) {
    accountingStore.showToast('Sheet kosong, tidak ada data', 'error')
    return
  }

  // Find header row (row containing strings like 'nik', 'nama', 'hadir', etc.)
  let headerRowIndex = 0
  for (let i = 0; i < Math.min(rows.length, 10); i++) {
    const rowStr = rows[i].map(c => String(c).toLowerCase()).join(' ')
    if (rowStr.includes('nik') || rowStr.includes('nama') || rowStr.includes('hadir') || rowStr.includes('masuk')) {
      headerRowIndex = i
      break
    }
  }

  rawRows.value = rows.slice(headerRowIndex)
  autoDetectColumns()
  recalculateRows()
}

function autoDetectColumns() {
  if (rawRows.value.length === 0) return
  const headers = rawRows.value[0].map(h => String(h).toLowerCase().trim())

  headers.forEach((h, idx) => {
    if (h.includes('nik') || h === 'id' || h.includes('nip') || h.includes('no_pegawai')) {
      colNIK.value = idx
    } else if (h.includes('nama') || h.includes('name') || h.includes('karyawan') || h.includes('pegawai')) {
      colName.value = idx
    } else if (h.includes('hadir') || h.includes('masuk') || h === 'hk' || h.includes('kehadiran') || h === 'h') {
      colHadir.value = idx
    } else if (h.includes('alpa') || h.includes('mangkir') || h === 'a' || h.includes('absen') || h.includes('tanpa_keterangan')) {
      colAlpa.value = idx
    } else if (h.includes('izin') || h.includes('cuti') || h.includes('sakit') || h.includes('i/s/c') || h.includes('leave')) {
      colIzin.value = idx
    } else if (h.includes('lembur') || h.includes('overtime') || h === 'ot' || h.includes('jam_lembur')) {
      colLembur.value = idx
    }
  })

  // Fallback if not detected
  if (colNIK.value === -1 && headers.length > 0) colNIK.value = 0
  if (colName.value === -1 && headers.length > 1) colName.value = 1
  if (colHadir.value === -1 && headers.length > 2) colHadir.value = 2
}

function recalculateRows() {
  if (rawRows.value.length <= 1) {
    parsedRows.value = []
    return
  }

  const dataRows = rawRows.value.slice(1)
  const results: ParsedAttendanceRow[] = []

  dataRows.forEach((row, idx) => {
    const rawNIK = colNIK.value >= 0 ? String(row[colNIK.value] || '').trim() : ''
    const rawName = colName.value >= 0 ? String(row[colName.value] || '').trim() : ''
    
    // Ignore completely empty rows
    if (!rawNIK && !rawName) return

    const hadir = colHadir.value >= 0 ? Math.max(0, parseInt(String(row[colHadir.value])) || 0) : 22
    const alpa = colAlpa.value >= 0 ? Math.max(0, parseInt(String(row[colAlpa.value])) || 0) : 0
    const izin = colIzin.value >= 0 ? Math.max(0, parseInt(String(row[colIzin.value])) || 0) : 0
    const lembur = colLembur.value >= 0 ? Math.max(0, parseFloat(String(row[colLembur.value]).replace(',', '.')) || 0) : 0

    // Match employee from DB
    let matchedEmp: Employee | null = null

    if (rawNIK) {
      const cleanNIK = rawNIK.toLowerCase().replace(/[^a-z0-9]/g, '')
      matchedEmp = props.employees.find(e => e.nik.toLowerCase().replace(/[^a-z0-9]/g, '') === cleanNIK) || null
    }

    if (!matchedEmp && rawName) {
      const cleanName = rawName.toLowerCase().trim()
      matchedEmp = props.employees.find(e => {
        const dbName = e.name.toLowerCase().trim()
        return dbName === cleanName || dbName.includes(cleanName) || cleanName.includes(dbName)
      }) || null
    }

    const isValid = matchedEmp !== null
    const validationMessage = matchedEmp
      ? 'Cocok dengan profil karyawan terdaftar'
      : 'NIK / Nama tidak ditemukan di database karyawan'

    let calculatedValues: any = undefined

    if (matchedEmp) {
      const baseSalary = matchedEmp.baseSalary || 0
      const positionAllowance = matchedEmp.positionAllowance || 0
      const dailyMealRate = matchedEmp.dailyMealTransportRate || (matchedEmp.transportMealAllowance ? Math.round(matchedEmp.transportMealAllowance / 22) : 30000)
      
      const transportMealAllowance = dailyMealRate * hadir
      const hourlyRate = ((baseSalary + positionAllowance) / 173) * 1.5
      const overtimePay = Math.round(lembur * hourlyRate)
      const attendanceDeduction = alpa > 0 ? Math.round(alpa * (baseSalary / 22)) : 0

      const gross = baseSalary + positionAllowance + transportMealAllowance + overtimePay
      const bpjsKes = Math.round(Math.min(gross, 12000000) * 0.01)
      const bpjsTk = Math.round(gross * 0.03)
      const pph21 = gross > 5400000 ? Math.round(gross * 0.015) : 0
      const totalDed = bpjsKes + bpjsTk + pph21 + attendanceDeduction
      const thp = gross - totalDed

      calculatedValues = {
        overtimePay,
        transportMealAllowance,
        attendanceDeduction,
        totalGrossSalary: gross,
        totalDeduction: totalDed,
        takeHomePay: thp
      }
    }

    results.push({
      rowIndex: idx + 2,
      rawNIK,
      rawName,
      hadir,
      alpa,
      izin,
      lemburHours: lembur,
      matchedEmployee: matchedEmp,
      isValid,
      validationMessage,
      calculatedValues
    })
  })

  parsedRows.value = results
}

function handleDownloadTemplate() {
  const templateData = props.employees.map(e => ({
    'NIK': e.nik,
    'Nama Karyawan': e.name,
    'Departemen': e.department,
    'Hari Hadir': 22,
    'Alpa (Mangkir)': 0,
    'Izin / Cuti': 0,
    'Jam Lembur': 0
  }))

  const ws = XLSX.utils.json_to_sheet(templateData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Data_Absensi')

  // Auto-width
  ws['!cols'] = [
    { wch: 20 },
    { wch: 30 },
    { wch: 22 },
    { wch: 12 },
    { wch: 14 },
    { wch: 12 },
    { wch: 12 }
  ]

  XLSX.writeFile(wb, `Template_Absensi_SMG_${props.selectedYear}_${props.selectedMonth}.xlsx`)
  accountingStore.showToast('Template Excel absensi berhasil diunduh', 'success')
}

async function handleProcessAndCreateSlips() {
  const validRows = parsedRows.value.filter(r => r.isValid && r.matchedEmployee)
  if (validRows.length === 0) {
    accountingStore.showToast('Tidak ada data karyawan valid yang dapat diproses', 'error')
    return
  }

  isSaving.value = true
  let successCount = 0

  const monthNames = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
  const pName = `${monthNames[props.selectedMonth - 1]} ${props.selectedYear}`

  try {
    for (const r of validRows) {
      const emp = r.matchedEmployee!
      const calc = r.calculatedValues!

      const slipData: Omit<PayslipItem, 'id' | 'slipNumber'> = {
        employeeId: emp.id,
        employeeNik: emp.nik,
        employeeName: emp.name,
        position: emp.position,
        department: emp.department,
        businessUnit: emp.businessUnit,
        periodMonth: props.selectedMonth,
        periodYear: props.selectedYear,
        periodName: pName,
        paymentDate: `${props.selectedYear}-${String(props.selectedMonth).padStart(2, '0')}-28`,

        // Direct Ingestion of Attendance without Manual Recalculation
        workingDays: 22,
        presentDays: r.hadir,
        absentDays: r.alpa,
        sickLeaveDays: 0,
        paidLeaveDays: r.izin,
        overtimeHours: r.lemburHours,
        dailyMealTransportRate: emp.dailyMealTransportRate || 30000,

        baseSalary: emp.baseSalary,
        positionAllowance: emp.positionAllowance,
        transportMealAllowance: calc.transportMealAllowance,
        shiftAttendanceAllowance: 0,
        incentiveBonus: 0,
        overtimePay: calc.overtimePay,
        totalGrossSalary: calc.totalGrossSalary,

        bpjsKesehatan: Math.round(Math.min(calc.totalGrossSalary, 12000000) * 0.01),
        bpjsKetenagakerjaan: Math.round(calc.totalGrossSalary * 0.03),
        pph21: calc.totalGrossSalary > 5400000 ? Math.round(calc.totalGrossSalary * 0.015) : 0,
        attendanceDeduction: calc.attendanceDeduction,
        loanDeduction: 0,
        totalDeduction: calc.totalDeduction,

        takeHomePay: calc.takeHomePay,
        status: 'DRAFT' as PayrollStatus,
        bankName: emp.bankName,
        bankAccountNumber: emp.bankAccountNumber
      }

      await payrollStore.addCustomPayslip(slipData)
      successCount++
    }

    accountingStore.showToast(`Sukses membuat ${successCount} draft slip gaji dari data absensi Excel!`, 'success')
    emit('imported', successCount)
    emit('close')
  } catch (err: any) {
    accountingStore.showToast('Gagal memproses sebagian slip gaji: ' + (err.message || ''), 'error')
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 bg-black/60 backdrop-blur-xs overflow-y-auto"
  >
    <div
      class="bg-white rounded-3xl border border-zinc-200 shadow-2xl w-full max-w-4xl my-auto overflow-hidden flex flex-col max-h-[90vh]"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-zinc-200 bg-zinc-50/80">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-emerald-700 text-white flex items-center justify-center shadow-xs">
            <FileSpreadsheet class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-extrabold text-zinc-950 tracking-tight">
              Import Data Absensi Excel (.xlsx / .xls)
            </h3>
            <p class="text-xs text-zinc-500">
              Otomatis membaca hari hadir, alpa, izin, dan lembur untuk periode {{ selectedMonth }}/{{ selectedYear }}
            </p>
          </div>
        </div>

        <button
          @click="emit('close')"
          class="p-2 text-zinc-400 hover:text-zinc-700 hover:bg-zinc-200/60 rounded-xl transition cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto space-y-6 flex-1 text-xs">
        
        <!-- Step 1: Upload Dropzone & Actions -->
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 p-4 rounded-2xl bg-[#FAF7F2] border border-amber-200/70">
          <div class="space-y-1">
            <div class="font-bold text-zinc-900 flex items-center gap-2">
              <Download class="w-4 h-4 text-emerald-700" />
              <span>Belum memiliki template Excel absensi?</span>
            </div>
            <p class="text-zinc-600">
              Unduh format resmi yang telah terisi daftar NIK dan nama karyawan aktif untuk diisi oleh bagian HR/Personalia.
            </p>
          </div>

          <button
            type="button"
            @click="handleDownloadTemplate"
            class="px-4 py-2.5 rounded-xl border border-emerald-700 text-emerald-800 bg-white hover:bg-emerald-50 font-bold transition flex items-center justify-center gap-2 shrink-0 cursor-pointer shadow-2xs"
          >
            <Download class="w-4 h-4" />
            <span>Download Format Template</span>
          </button>
        </div>

        <!-- File Upload Area -->
        <div
          @dragover.prevent
          @drop="handleFileDrop"
          @click="triggerFileInput"
          class="border-2 border-dashed rounded-3xl p-6 sm:p-8 text-center cursor-pointer transition-all hover:bg-zinc-50/80"
          :class="fileName ? 'border-emerald-500 bg-emerald-50/20' : 'border-zinc-300 bg-white'"
        >
          <input
            ref="fileInputRef"
            type="file"
            accept=".xlsx,.xls"
            class="hidden"
            @change="handleFileSelect"
          />

          <div v-if="!fileName" class="space-y-2">
            <div class="w-12 h-12 rounded-2xl bg-zinc-100 text-zinc-500 mx-auto flex items-center justify-center">
              <Upload class="w-6 h-6" />
            </div>
            <div class="font-bold text-sm text-zinc-800">
              Pilih atau seret file Excel absensi ke sini (.xlsx atau .xls)
            </div>
            <p class="text-zinc-500 text-[11px]">
              Sistem akan membaca kolom NIK, Nama, Hari Hadir, Alpa, Izin/Cuti, dan Jam Lembur.
            </p>
          </div>

          <div v-else class="flex items-center justify-center gap-3">
            <FileSpreadsheet class="w-6 h-6 text-emerald-700 shrink-0" />
            <div class="text-left">
              <div class="font-bold text-zinc-900 text-sm">{{ fileName }}</div>
              <div class="text-[11px] text-zinc-500">
                {{ rawRows.length }} baris terdeteksi • Klik untuk mengganti file
              </div>
            </div>
          </div>
        </div>

        <!-- Sheet Selector & Column Mapping (If file loaded) -->
        <div v-if="rawRows.length > 0" class="space-y-4 pt-2 border-t border-zinc-100">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <h4 class="font-extrabold text-zinc-900 flex items-center gap-2">
              <Table class="w-4 h-4 text-emerald-700" />
              <span>Pemetaan Kolom Data Absensi (Auto-Detection)</span>
            </h4>

            <div v-if="sheetNames.length > 1" class="flex items-center gap-2">
              <label class="text-[11px] font-bold text-zinc-600">Pilih Sheet:</label>
              <select
                v-model="selectedSheet"
                @change="handleSheetChange"
                class="px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-semibold"
              >
                <option v-for="s in sheetNames" :key="s" :value="s">{{ s }}</option>
              </select>
            </div>
          </div>

          <!-- Mapping Grid -->
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 p-4 rounded-2xl bg-zinc-50 border border-zinc-200">
            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom NIK / ID:</label>
              <select
                v-model.number="colNIK"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Tidak Dipetakan --</option>
                <option v-for="h in headerOptions" :key="'nik-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>

            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom Nama Karyawan:</label>
              <select
                v-model.number="colName"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Tidak Dipetakan --</option>
                <option v-for="h in headerOptions" :key="'name-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>

            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom Hari Hadir:</label>
              <select
                v-model.number="colHadir"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Default 22 Hari --</option>
                <option v-for="h in headerOptions" :key="'hadir-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>

            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom Alpa / Mangkir:</label>
              <select
                v-model.number="colAlpa"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Default 0 --</option>
                <option v-for="h in headerOptions" :key="'alpa-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>

            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom Izin / Sakit / Cuti:</label>
              <select
                v-model.number="colIzin"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Default 0 --</option>
                <option v-for="h in headerOptions" :key="'izin-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>

            <div>
              <label class="block text-[10px] font-bold uppercase text-zinc-500 mb-1">Kolom Jam Lembur:</label>
              <select
                v-model.number="colLembur"
                @change="recalculateRows"
                class="w-full px-2.5 py-1.5 rounded-lg border border-zinc-300 text-xs font-medium bg-white"
              >
                <option :value="-1">-- Default 0 --</option>
                <option v-for="h in headerOptions" :key="'lembur-' + h.index" :value="h.index">{{ h.label }}</option>
              </select>
            </div>
          </div>

          <!-- Summary Badges -->
          <div class="grid grid-cols-3 gap-3">
            <div class="p-3 rounded-xl bg-zinc-100 border border-zinc-200 text-center">
              <span class="text-[10px] font-bold uppercase text-zinc-500 block">Total Baris File</span>
              <span class="font-black text-lg text-zinc-900">{{ summaryStats.total }}</span>
            </div>

            <div class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-center">
              <span class="text-[10px] font-bold uppercase text-emerald-800 block">Cocok di Database (Valid)</span>
              <span class="font-black text-lg text-emerald-700">{{ summaryStats.valid }}</span>
            </div>

            <div class="p-3 rounded-xl bg-amber-50 border border-amber-200 text-center">
              <span class="text-[10px] font-bold uppercase text-amber-900 block">Peringatan / Tidak Cocok</span>
              <span class="font-black text-lg text-amber-700">{{ summaryStats.invalid }}</span>
            </div>
          </div>

          <!-- Preview Table -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <h5 class="font-bold text-zinc-800">Preview Data Hasil Pembacaan Excel:</h5>
              <span class="text-[10px] text-zinc-400">Nilai lembur & tunjangan dihitung otomatis</span>
            </div>

            <div class="overflow-x-auto rounded-2xl border border-zinc-200 max-h-60 overflow-y-auto">
              <table class="w-full text-left text-[11px]">
                <thead class="bg-zinc-50 border-b border-zinc-200 uppercase tracking-wider text-[10px] font-bold text-zinc-600 sticky top-0">
                  <tr>
                    <th class="px-3 py-2 text-center">Status</th>
                    <th class="px-3 py-2">NIK File</th>
                    <th class="px-3 py-2">Nama Karyawan</th>
                    <th class="px-3 py-2 text-center">Hadir</th>
                    <th class="px-3 py-2 text-center">Alpa</th>
                    <th class="px-3 py-2 text-center">Izin</th>
                    <th class="px-3 py-2 text-center">Lembur</th>
                    <th class="px-3 py-2 text-right">Nilai Lembur</th>
                    <th class="px-3 py-2 text-right">Estimasi THP</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-zinc-100">
                  <tr
                    v-for="row in parsedRows"
                    :key="row.rowIndex"
                    class="hover:bg-zinc-50/60"
                    :class="{ 'bg-amber-50/30': !row.isValid }"
                  >
                    <!-- Status -->
                    <td class="px-3 py-2 text-center">
                      <span
                        v-if="row.isValid"
                        class="inline-flex items-center gap-1 text-[10px] font-bold text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-md"
                        title="Cocok di database karyawan"
                      >
                        <Check class="w-3 h-3" /> Valid
                      </span>
                      <span
                        v-else
                        class="inline-flex items-center gap-1 text-[10px] font-bold text-amber-800 bg-amber-100 px-2 py-0.5 rounded-md"
                        :title="row.validationMessage"
                      >
                        <AlertTriangle class="w-3 h-3 text-amber-600" /> Lewati
                      </span>
                    </td>

                    <td class="px-3 py-2 font-mono font-bold text-zinc-900">{{ row.rawNIK || '—' }}</td>
                    <td class="px-3 py-2">
                      <div class="font-bold text-zinc-900">{{ row.matchedEmployee?.name || row.rawName }}</div>
                      <div class="text-[10px] text-zinc-500" v-if="row.matchedEmployee">
                        {{ row.matchedEmployee.position }} ({{ row.matchedEmployee.businessUnit }})
                      </div>
                      <div class="text-[10px] text-rose-600 font-semibold" v-else>
                        {{ row.validationMessage }}
                      </div>
                    </td>

                    <td class="px-3 py-2 text-center font-bold text-emerald-800">{{ row.hadir }} hr</td>
                    <td class="px-3 py-2 text-center font-bold" :class="row.alpa > 0 ? 'text-rose-600' : 'text-zinc-500'">
                      {{ row.alpa }} hr
                    </td>
                    <td class="px-3 py-2 text-center text-zinc-600">{{ row.izin }} hr</td>
                    <td class="px-3 py-2 text-center font-bold text-indigo-700">
                      {{ row.lemburHours > 0 ? `+${row.lemburHours}j` : '—' }}
                    </td>
                    <td class="px-3 py-2 text-right font-mono tabular-nums text-zinc-800">
                      {{ row.calculatedValues ? formatIDR(row.calculatedValues.overtimePay) : '—' }}
                    </td>
                    <td class="px-3 py-2 text-right font-mono font-bold tabular-nums text-emerald-900">
                      {{ row.calculatedValues ? formatIDR(row.calculatedValues.takeHomePay) : '—' }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

      </div>

      <!-- Modal Footer -->
      <div class="px-6 py-4 border-t border-zinc-200 bg-zinc-50/50 flex flex-col sm:flex-row items-center justify-between gap-3">
        <div class="text-xs text-zinc-500">
          <span v-if="summaryStats.valid > 0">
            Siap membuat <strong>{{ summaryStats.valid }}</strong> slip gaji untuk periode {{ selectedMonth }}/{{ selectedYear }}.
          </span>
          <span v-else>
            Unggah file Excel untuk memproses slip gaji.
          </span>
        </div>

        <div class="flex items-center gap-2.5 w-full sm:w-auto">
          <button
            type="button"
            @click="emit('close')"
            class="flex-1 sm:flex-none px-4 py-2.5 rounded-xl border border-zinc-200 text-zinc-700 hover:bg-zinc-100 font-bold transition cursor-pointer text-xs"
          >
            Batal
          </button>

          <button
            type="button"
            @click="handleProcessAndCreateSlips"
            :disabled="summaryStats.valid === 0 || isSaving"
            class="flex-1 sm:flex-none px-5 py-2.5 rounded-xl bg-emerald-700 hover:bg-emerald-800 text-white font-bold transition flex items-center justify-center gap-2 shadow-sm disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer text-xs"
          >
            <RefreshCw v-if="isSaving" class="w-4 h-4 animate-spin" />
            <span v-else>Proses & Buat Slip Gaji ({{ summaryStats.valid }} Karyawan)</span>
            <ArrowRight v-if="!isSaving" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
