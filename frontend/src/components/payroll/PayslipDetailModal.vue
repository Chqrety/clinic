<script setup lang="ts">
import { computed } from 'vue'
import type { PayslipItem } from '@/types/payroll'
import { formatIDR } from '@/utils/currency'
import { terbilang } from '@/utils/terbilang'
import { useLocaleStore } from '@/stores/locale'
import {
  X,
  Printer,
  Building2,
  CheckCircle2,
  Clock,
  CreditCard,
  User,
  ShieldCheck,
  FileCheck
} from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
  slip: PayslipItem | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const localeStore = useLocaleStore()

const terbilangText = computed(() => {
  if (!props.slip) return ''
  return terbilang(props.slip.takeHomePay)
})

function handlePrint() {
  window.print()
}
</script>

<template>
  <div
    v-if="show && slip"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-5 bg-black/60 backdrop-blur-xs overflow-y-auto print:p-0 print:bg-white print:static"
  >
    <!-- Modal Card / Printable Sheet -->
    <div
      class="bg-white rounded-2xl border border-zinc-200 shadow-2xl w-full max-w-2xl my-auto overflow-hidden print:border-none print:shadow-none print:rounded-none print:max-w-none print:w-full"
    >
      <!-- Action Bar (Hidden on print) -->
      <div class="flex items-center justify-between px-5 py-3.5 border-b border-zinc-200 bg-zinc-50/80 print:hidden">
        <div class="flex items-center gap-2">
          <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
          <span class="text-xs font-bold text-zinc-700 uppercase tracking-wider">
            {{ localeStore.isEnglish ? 'Official Employee Payslip' : 'Slip Gaji Resmi Karyawan' }}
          </span>
          <span class="text-xs font-mono font-semibold text-zinc-500">({{ slip.slipNumber }})</span>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="handlePrint"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold transition-all shadow-2xs active:scale-95 cursor-pointer"
          >
            <Printer class="w-3.5 h-3.5" />
            <span>{{ localeStore.isEnglish ? 'Print Payslip' : 'Cetak Slip Gaji' }}</span>
          </button>

          <button
            @click="emit('close')"
            class="p-1.5 text-zinc-400 hover:text-zinc-700 hover:bg-zinc-200/60 rounded-lg transition-colors cursor-pointer"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- Printable Payslip Document -->
      <div id="printable-payslip" class="p-6 sm:p-8 space-y-6 text-zinc-900 font-['Plus_Jakarta_Sans',sans-serif]">
        <!-- Header Klinik SMG / SMC -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b-2 border-zinc-900">
          <div class="flex items-center gap-3.5">
            <div
              class="w-12 h-12 rounded-xl text-white flex items-center justify-center font-black text-xl tracking-tight shrink-0 shadow-xs"
              :class="(slip.businessUnit === 'APOTEK' || slip.department.toLowerCase().includes('farmasi')) ? 'bg-[#1E3A8A]' : 'bg-[#0C3B36]'"
            >
              {{ (slip.businessUnit === 'APOTEK' || slip.department.toLowerCase().includes('farmasi')) ? 'SMC' : 'SMG' }}
            </div>
            <div>
              <h2 class="text-base font-extrabold text-zinc-950 tracking-tight leading-snug">
                {{ (slip.businessUnit === 'APOTEK' || slip.department.toLowerCase().includes('farmasi')) ? 'CV SEMARANG MEDICAL CLINIC' : 'CV SEMARANG MEDICAL GROUP' }}
              </h2>
              <p class="text-xs font-semibold text-zinc-600">
                {{ (slip.businessUnit === 'APOTEK' || slip.department.toLowerCase().includes('farmasi')) ? 'Instalasi Farmasi, Apotek & Logistik Medis' : 'Kidinest Pediatric Care & Semarang Skin Center' }}
              </p>
              <p class="text-[10px] text-zinc-400 mt-0.5">
                Jl. Melati Utara No. 3, Kota Semarang • Telp: (024) 8412345 • payroll@smg-clinic.id
              </p>
            </div>
          </div>

          <div class="text-right sm:border-l sm:border-zinc-200 sm:pl-5">
            <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">
              {{ localeStore.isEnglish ? 'DOCUMENT TYPE' : 'JENIS DOKUMEN' }}
            </div>
            <div class="text-sm font-black text-zinc-900 tracking-tight">
              SLIP GAJI KARYAWAN
            </div>
            <div class="text-xs font-bold text-emerald-700 mt-0.5">
              Periode: {{ slip.periodName }}
            </div>
          </div>
        </div>

        <!-- Meta Karyawan -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 p-3.5 rounded-xl bg-zinc-50 border border-zinc-200/80 text-xs">
          <div>
            <span class="text-[10px] font-bold uppercase text-zinc-400 block">{{ localeStore.isEnglish ? 'EMPLOYEE NAME' : 'NAMA KARYAWAN' }}</span>
            <span class="font-bold text-zinc-900 text-xs">{{ slip.employeeName }}</span>
            <span class="font-mono text-[10px] text-zinc-500 block">{{ slip.employeeNik }}</span>
          </div>

          <div>
            <span class="text-[10px] font-bold uppercase text-zinc-400 block">{{ localeStore.isEnglish ? 'POSITION' : 'JABATAN & TUGAS' }}</span>
            <span class="font-semibold text-zinc-800 text-xs block truncate">{{ slip.position }}</span>
            <span class="text-[10px] text-zinc-500 block">{{ slip.department }}</span>
          </div>

          <div>
            <span class="text-[10px] font-bold uppercase text-zinc-400 block">ENTITAS & UNIT</span>
            <span class="font-bold text-xs block" :class="(slip.businessUnit === 'APOTEK' || slip.businessUnit === 'SMC' || slip.department.toLowerCase().includes('farmasi')) ? 'text-blue-900' : 'text-emerald-950'">
              {{ (slip.businessUnit === 'APOTEK' || slip.businessUnit === 'SMC' || slip.department.toLowerCase().includes('farmasi')) ? 'CV Semarang Medical Clinic' : 'CV Semarang Medical Group' }}
            </span>
            <span
              class="inline-block mt-0.5 px-1.5 py-0.2 rounded text-[9px] font-bold uppercase tracking-wider border"
              :class="slip.businessUnit === 'APOTEK' || slip.businessUnit === 'SMC' ? 'bg-blue-50 text-blue-800 border-blue-200' : slip.businessUnit === 'SSC' ? 'bg-rose-50 text-rose-800 border-rose-200' : slip.businessUnit === 'KIDINEST' ? 'bg-emerald-50 text-emerald-800 border-emerald-200' : 'bg-zinc-100 text-zinc-700 border-zinc-200'"
            >
              {{ slip.businessUnit === 'APOTEK' || slip.businessUnit === 'SMC' ? 'Instalasi Farmasi / Apotek' : slip.businessUnit === 'SSC' ? 'Semarang Skin Center (SSC)' : slip.businessUnit === 'KIDINEST' ? 'Kidinest Pediatric Care' : 'Kantor Pusat SMG' }}
            </span>
          </div>

          <div>
            <span class="text-[10px] font-bold uppercase text-zinc-400 block">{{ localeStore.isEnglish ? 'TRANSFER ACCOUNT' : 'REKENING TRANSFER' }}</span>
            <span class="font-medium text-zinc-800 text-xs block">{{ slip.bankName }}</span>
            <span class="font-mono text-[11px] text-zinc-600">{{ slip.bankAccountNumber }}</span>
          </div>
        </div>

        <!-- Rekapitulasi Kehadiran & Hari Masuk Kerja -->
        <div class="p-3.5 rounded-xl bg-sky-50/60 border border-sky-200/80 text-xs">
          <div class="flex items-center justify-between mb-2">
            <span class="text-[10px] font-bold uppercase tracking-wider text-sky-950">
              {{ localeStore.isEnglish ? 'ATTENDANCE & WORKING DAYS RECORD' : 'REKAPITULASI KEHADIRAN & HARI MASUK KERJA' }}
            </span>
            <span class="text-[11px] font-semibold text-sky-900">
              {{ localeStore.isEnglish ? 'Standard' : 'Standar Periode' }}: {{ slip.workingDays || 22 }} {{ localeStore.isEnglish ? 'Days' : 'Hari' }}
            </span>
          </div>
          <div class="grid grid-cols-3 sm:grid-cols-6 gap-2 text-center">
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Present Days' : 'Hari Masuk' }}</span>
              <span class="font-black text-sm text-emerald-700">{{ slip.presentDays !== undefined ? slip.presentDays : (slip.workingDays || 22) }} hr</span>
            </div>
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Paid Leave' : 'Izin / Cuti' }}</span>
              <span class="font-bold text-sm text-zinc-800">{{ slip.paidLeaveDays || 0 }} hr</span>
            </div>
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Sick Leave' : 'Sakit (Dokter)' }}</span>
              <span class="font-bold text-sm text-amber-700">{{ slip.sickLeaveDays || 0 }} hr</span>
            </div>
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Absent' : 'Alpa / Mangkir' }}</span>
              <span :class="(slip.absentDays || 0) > 0 ? 'text-rose-600 font-black' : 'text-zinc-600 font-bold'" class="text-sm">
                {{ slip.absentDays || 0 }} hr
              </span>
            </div>
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Overtime' : 'Jam Lembur' }}</span>
              <span class="font-bold text-sm text-indigo-700">{{ slip.overtimeHours || 0 }} jam</span>
            </div>
            <div class="bg-white p-2 rounded-lg border border-sky-100 shadow-2xs">
              <span class="text-[10px] text-zinc-500 block">{{ localeStore.isEnglish ? 'Daily Meal Rate' : 'Tarif Uang Makan' }}</span>
              <span class="font-bold text-xs text-zinc-800">{{ formatIDR(slip.dailyMealTransportRate || 30000) }}/hr</span>
            </div>
          </div>
        </div>

        <!-- 2 Kolom: Penerimaan (Earnings) & Potongan (Deductions) -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
          <!-- Kolom Penerimaan -->
          <div class="border border-emerald-200 rounded-xl overflow-hidden bg-emerald-50/20">
            <div class="px-3.5 py-2 bg-emerald-600 text-white flex items-center justify-between">
              <span class="text-xs font-bold uppercase tracking-wider">A. {{ localeStore.isEnglish ? 'EARNINGS / PENERIMAAN' : 'PENERIMAAN (EARNINGS)' }}</span>
              <span class="text-[10px] font-bold bg-emerald-700 px-1.5 py-0.5 rounded">IDR</span>
            </div>

            <div class="p-3.5 space-y-2 text-xs divide-y divide-emerald-100">
              <div class="flex justify-between pt-1">
                <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Basic Salary' : 'Gaji Pokok' }}</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.baseSalary) }}</span>
              </div>

              <div class="flex justify-between pt-1.5">
                <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Position Allowance' : 'Tunjangan Jabatan' }}</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.positionAllowance) }}</span>
              </div>

              <div class="flex justify-between pt-1.5">
                <div>
                  <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Meal & Transport Allowance' : 'Tunjangan Uang Makan & Transport' }}</span>
                  <span class="text-[10px] text-zinc-400 block" v-if="slip.presentDays !== undefined">
                    ({{ slip.presentDays }} hari hadir × {{ formatIDR(slip.dailyMealTransportRate || 30000) }})
                  </span>
                </div>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.transportMealAllowance) }}</span>
              </div>

              <div v-if="slip.shiftAttendanceAllowance > 0" class="flex justify-between pt-1.5">
                <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Shift & Attendance' : 'Tunjangan Shift / Kehadiran' }}</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.shiftAttendanceAllowance) }}</span>
              </div>

              <div v-if="slip.incentiveBonus > 0" class="flex justify-between pt-1.5">
                <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Incentive & Bonus' : 'Insentif Pelayanan & Bonus' }}</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.incentiveBonus) }}</span>
              </div>

              <div v-if="slip.overtimePay > 0" class="flex justify-between pt-1.5">
                <div>
                  <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Overtime Pay' : 'Uang Lembur (Overtime)' }}</span>
                  <span class="text-[10px] text-zinc-400 block" v-if="slip.overtimeHours">
                    ({{ slip.overtimeHours }} jam lembur)
                  </span>
                </div>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.overtimePay) }}</span>
              </div>

              <div class="flex justify-between pt-2 border-t-2 border-emerald-300 font-extrabold text-emerald-950 bg-emerald-100/50 px-2 py-1.5 rounded-lg mt-2">
                <span>{{ localeStore.isEnglish ? 'TOTAL GROSS' : 'TOTAL PENGHASILAN BRUTO' }}</span>
                <span class="tabular-nums">{{ formatIDR(slip.totalGrossSalary) }}</span>
              </div>
            </div>
          </div>

          <!-- Kolom Potongan -->
          <div class="border border-rose-200 rounded-xl overflow-hidden bg-rose-50/20">
            <div class="px-3.5 py-2 bg-rose-600 text-white flex items-center justify-between">
              <span class="text-xs font-bold uppercase tracking-wider">B. {{ localeStore.isEnglish ? 'DEDUCTIONS / POTONGAN' : 'POTONGAN (DEDUCTIONS)' }}</span>
              <span class="text-[10px] font-bold bg-rose-700 px-1.5 py-0.5 rounded">IDR</span>
            </div>

            <div class="p-3.5 space-y-2 text-xs divide-y divide-rose-100">
              <div class="flex justify-between pt-1">
                <span class="text-zinc-600">BPJS Kesehatan (1%)</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.bpjsKesehatan) }}</span>
              </div>

              <div class="flex justify-between pt-1.5">
                <span class="text-zinc-600">BPJS Ketenagakerjaan (3%)</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.bpjsKetenagakerjaan) }}</span>
              </div>

              <div class="flex justify-between pt-1.5">
                <span class="text-zinc-600">PPh Pasal 21 (TER Pajak)</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.pph21) }}</span>
              </div>

              <div v-if="slip.attendanceDeduction > 0" class="flex justify-between pt-1.5">
                <div>
                  <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Late / Absence Penalty' : 'Potongan Mangkir / Absensi' }}</span>
                  <span class="text-[10px] text-rose-500 block" v-if="slip.absentDays">
                    ({{ slip.absentDays }} hari alpa)
                  </span>
                </div>
                <span class="font-bold tabular-nums text-rose-700">-{{ formatIDR(slip.attendanceDeduction) }}</span>
              </div>

              <div v-if="slip.loanDeduction > 0" class="flex justify-between pt-1.5">
                <span class="text-zinc-600">{{ localeStore.isEnglish ? 'Loan / Cashbon Installment' : 'Cicilan Kasbon / Pinjaman' }}</span>
                <span class="font-bold tabular-nums text-zinc-900">{{ formatIDR(slip.loanDeduction) }}</span>
              </div>

              <div class="flex justify-between pt-2 border-t-2 border-rose-300 font-extrabold text-rose-950 bg-rose-100/50 px-2 py-1.5 rounded-lg mt-2">
                <span>{{ localeStore.isEnglish ? 'TOTAL DEDUCTIONS' : 'TOTAL POTONGAN' }}</span>
                <span class="tabular-nums">{{ formatIDR(slip.totalDeduction) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Highlight Take Home Pay -->
        <div class="p-4 rounded-xl bg-zinc-950 text-white flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-sm print:bg-zinc-100 print:text-zinc-950 print:border print:border-zinc-300">
          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-emerald-400 print:text-emerald-700">
              {{ localeStore.isEnglish ? 'TAKE HOME PAY (NET SALARY TRANSFERRED)' : 'GAJI BERSIH DITERIMA (TAKE HOME PAY)' }}
            </div>
            <div class="text-[11px] italic text-zinc-300 print:text-zinc-600 mt-0.5">
              Terbilang: "{{ terbilangText }}"
            </div>
          </div>

          <div class="text-2xl font-black tabular-nums tracking-tight text-white print:text-zinc-950 shrink-0">
            {{ formatIDR(slip.takeHomePay) }}
          </div>
        </div>

        <!-- Tanda Tangan & Footer Pengesahan -->
        <div class="pt-6 border-t border-zinc-200 grid grid-cols-2 gap-8 text-center text-xs">
          <div>
            <p class="text-[11px] text-zinc-500">{{ localeStore.isEnglish ? 'Prepared & Authorized by' : 'Dibuat & Diserahkan Oleh' }},</p>
            <p class="font-bold text-zinc-900 mt-0.5">Semarang Medical Group</p>
            <div class="h-16 flex items-end justify-center">
              <span class="text-[10px] text-zinc-400 font-mono italic">[Tanda Tangan Digital Tersertifikasi]</span>
            </div>
            <p class="font-bold text-zinc-950 border-t border-zinc-300 pt-1 inline-block min-w-[160px]">
              Budi Santoso, SE, Ak., CA
            </p>
            <p class="text-[10px] text-zinc-500">Finance & Accounting Manager</p>
          </div>

          <div>
            <p class="text-[11px] text-zinc-500">{{ localeStore.isEnglish ? 'Received by Employee' : 'Diterima Oleh Karyawan' }},</p>
            <p class="font-bold text-zinc-900 mt-0.5">{{ slip.employeeName }}</p>
            <div class="h-16 flex items-end justify-center">
              <span class="text-[10px] text-zinc-400 font-mono italic">[Tanda Tangan Karyawan]</span>
            </div>
            <p class="font-bold text-zinc-950 border-t border-zinc-300 pt-1 inline-block min-w-[160px]">
              {{ slip.employeeName }}
            </p>
            <p class="text-[10px] text-zinc-500">NIK: {{ slip.employeeNik }}</p>
          </div>
        </div>

        <!-- Catatan Rahasia -->
        <p class="text-[10px] text-zinc-400 text-center italic pt-2">
          * Dokumen slip gaji ini bersifat RAHASIA (Strictly Confidential). Diterbitkan secara resmi melalui Sistem Akuntansi & Payroll SMG.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #printable-payslip,
  #printable-payslip * {
    visibility: visible;
  }
  #printable-payslip {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    margin: 0;
    padding: 20px;
  }
}
</style>
