import type { BusinessUnit } from './accounting'

export type PayrollStatus = 'DRAFT' | 'APPROVED' | 'PAID'
export type EmploymentType = 'TETAP' | 'KONTRAK' | 'PARUH_WAKTU'

export interface Employee {
  id: string
  nik: string
  name: string
  position: string
  department: string
  businessUnit: BusinessUnit
  employmentType: EmploymentType
  joinDate: string
  baseSalary: number
  positionAllowance: number
  transportMealAllowance: number
  dailyMealTransportRate?: number
  bankName: string
  bankAccountNumber: string
  npwp?: string
  ptkpStatus: string
  bpjsKesehatanNo?: string
  bpjsKetenagakerjaanNo?: string
  isActive: boolean
}

export interface PayslipItem {
  id: string
  slipNumber: string
  employeeId: string
  employeeNik: string
  employeeName: string
  position: string
  department: string
  businessUnit: BusinessUnit
  periodMonth: number
  periodYear: number
  periodName: string
  paymentDate: string
  
  // Kehadiran / Hari Masuk (Attendance)
  workingDays: number
  presentDays: number
  sickLeaveDays: number
  paidLeaveDays: number
  absentDays: number
  overtimeHours: number
  dailyMealTransportRate: number

  // Earnings (Penerimaan)
  baseSalary: number
  positionAllowance: number
  transportMealAllowance: number
  shiftAttendanceAllowance: number
  incentiveBonus: number
  overtimePay: number
  totalGrossSalary: number
  
  // Deductions (Potongan)
  bpjsKesehatan: number
  bpjsKetenagakerjaan: number
  pph21: number
  attendanceDeduction: number
  loanDeduction: number
  totalDeduction: number
  
  // Net
  takeHomePay: number
  
  status: PayrollStatus
  bankName: string
  bankAccountNumber: string
  journalId?: string
  journalNumber?: string
  notes?: string
  approvedBy?: string
  approvedAt?: string
  paidAt?: string
}

export interface PayrollBatchSummary {
  periodMonth: number
  periodYear: number
  totalEmployees: number
  totalGrossSalary: number
  totalDeductions: number
  totalTakeHomePay: number
  totalPph21: number
  totalBpjs: number
  isJournalPosted: boolean
  journalNumber?: string
}
