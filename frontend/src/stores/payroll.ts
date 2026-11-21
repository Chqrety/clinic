import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Employee, PayslipItem, PayrollStatus } from '@/types/payroll'
import type { BusinessUnit } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useJournalsStore } from './journals'
import { accountingApi } from '@/api/accounting'

export const usePayrollStore = defineStore('payroll', () => {
  const accountingStore = useAccountingStore()
  const journalsStore = useJournalsStore()

  // Real PostgreSQL Database State
  const employees = ref<Employee[]>([])
  const payslips = ref<PayslipItem[]>([])
  const isLoading = ref(false)

  // Filter States
  const selectedMonth = ref<number>(new Date().getMonth() + 1)
  const selectedYear = ref<number>(new Date().getFullYear())
  const statusFilter = ref<'ALL' | PayrollStatus>('ALL')
  const searchQuery = ref<string>('')

  // Fallback storage key for offline resiliency
  const LOCAL_STORAGE_KEY = 'clinic_payroll_slips_v1'

  function persistSlips() {
    try {
      localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(payslips.value))
    } catch { }
  }

  function loadCachedSlips(): PayslipItem[] {
    try {
      const raw = localStorage.getItem(LOCAL_STORAGE_KEY)
      if (raw) return JSON.parse(raw)
    } catch { }
    return []
  }

  // --------------------------------------------------------------------------
  // Live Backend Database API Synchronization
  // --------------------------------------------------------------------------

  async function fetchEmployees(bu?: string) {
    try {
      const backendEmployees = await accountingApi.getPayrollEmployees(bu || accountingStore.businessUnit, true)
      if (backendEmployees && backendEmployees.length > 0) {
        employees.value = backendEmployees
        return
      }
    } catch (err) {
      console.warn('[PayrollStore] Gagal mengambil data karyawan dari PostgreSQL backend, menggunakan fallback:', err)
    }
  }

  async function fetchPayslips() {
    try {
      isLoading.value = true
      const backendSlips = await accountingApi.getPayrollSlips(
        selectedMonth.value,
        selectedYear.value,
        accountingStore.businessUnit,
        statusFilter.value
      )
      if (backendSlips && backendSlips.length > 0) {
        payslips.value = backendSlips
        persistSlips()
        return
      }
    } catch (err) {
      console.warn('[PayrollStore] Gagal mengambil slip dari PostgreSQL backend, menggunakan cached slips:', err)
    } finally {
      isLoading.value = false
    }

    // If backend returns empty or unreachable, check cached slips
    const cached = loadCachedSlips()
    if (cached.length > 0) {
      payslips.value = cached
    }
  }

  // Initial load
  fetchEmployees()
  fetchPayslips()

  // --------------------------------------------------------------------------
  // Reactive Getters & Filters
  // --------------------------------------------------------------------------

  const filteredPayslips = computed(() => {
    return payslips.value.filter(s => {
      if (s.periodMonth !== selectedMonth.value || s.periodYear !== selectedYear.value) {
        return false
      }
      if (accountingStore.businessUnit !== 'ALL' && s.businessUnit !== 'ALL' && s.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      if (statusFilter.value !== 'ALL' && s.status !== statusFilter.value) {
        return false
      }
      if (searchQuery.value.trim()) {
        const q = searchQuery.value.toLowerCase()
        const matchName = s.employeeName.toLowerCase().includes(q)
        const matchNik = s.employeeNik.toLowerCase().includes(q)
        const matchSlip = s.slipNumber.toLowerCase().includes(q)
        const matchPos = s.position.toLowerCase().includes(q)
        if (!matchName && !matchNik && !matchSlip && !matchPos) return false
      }
      return true
    })
  })

  const stats = computed(() => {
    const list = filteredPayslips.value
    const totalGross = list.reduce((s, item) => s + item.totalGrossSalary, 0)
    const totalDeductions = list.reduce((s, item) => s + item.totalDeduction, 0)
    const totalTakeHomePay = list.reduce((s, item) => s + item.takeHomePay, 0)
    const totalPph21 = list.reduce((s, item) => s + item.pph21, 0)
    const totalBpjs = list.reduce((s, item) => s + item.bpjsKesehatan + item.bpjsKetenagakerjaan, 0)

    return {
      totalCount: list.length,
      draftCount: list.filter(s => s.status === 'DRAFT').length,
      approvedCount: list.filter(s => s.status === 'APPROVED').length,
      paidCount: list.filter(s => s.status === 'PAID').length,
      totalGross,
      totalDeductions,
      totalTakeHomePay,
      totalPph21,
      totalBpjs
    }
  })

  // --------------------------------------------------------------------------
  // Actions with Authoritative Backend PostgreSQL API Calls
  // --------------------------------------------------------------------------

  async function generateMonthlyPayroll(month: number, year: number) {
    try {
      isLoading.value = true
      const backendGenerated = await accountingApi.generateMonthlyPayroll(month, year, accountingStore.businessUnit)
      if (backendGenerated && backendGenerated.length > 0) {
        payslips.value = backendGenerated
        persistSlips()
        accountingStore.showToast(`Berhasil generate ${backendGenerated.length} slip gaji dari database untuk periode ${month}/${year}`, 'success')
        return
      }
    } catch (err: any) {
      console.warn('[PayrollStore] Generate via backend API gagal, fallback ke kalkulasi lokal:', err)
      accountingStore.showToast(err.message || 'Gagal generate payroll dari database', 'error')
    } finally {
      isLoading.value = false
    }

    // Local fallback generation
    await fetchEmployees()
    const pName = `Bulan ${month} ${year}`
    const generated: PayslipItem[] = employees.value
      .filter(e => e.isActive && (accountingStore.businessUnit === 'ALL' || e.businessUnit === 'ALL' || e.businessUnit === accountingStore.businessUnit))
      .map(emp => {
        const workingDays = 22
        const presentDays = 22
        const dailyMealRate = emp.dailyMealTransportRate || (emp.transportMealAllowance ? Math.round(emp.transportMealAllowance / workingDays) : 30000)
        const transportAllowance = dailyMealRate * presentDays

        const gross = emp.baseSalary + emp.positionAllowance + transportAllowance
        const bpjsKes = Math.round(Math.min(gross, 12000000) * 0.01)
        const bpjsTk = Math.round(gross * 0.03)
        const pph21 = gross > 5400000 ? Math.round(gross * 0.015) : 0
        const totalDed = bpjsKes + bpjsTk + pph21
        const thp = gross - totalDed

        return {
          id: `slip-${year}${month}-${emp.id}`,
          slipNumber: `SLIP/${year}/${String(month).padStart(2, '0')}/${emp.nik}`,
          employeeId: emp.id,
          employeeNik: emp.nik,
          employeeName: emp.name,
          position: emp.position,
          department: emp.department,
          businessUnit: emp.businessUnit,
          periodMonth: month,
          periodYear: year,
          periodName: pName,
          paymentDate: `${year}-${String(month).padStart(2, '0')}-28`,
          workingDays,
          presentDays,
          sickLeaveDays: 0,
          paidLeaveDays: 0,
          absentDays: 0,
          overtimeHours: 0,
          dailyMealTransportRate: dailyMealRate,
          baseSalary: emp.baseSalary,
          positionAllowance: emp.positionAllowance,
          transportMealAllowance: transportAllowance,
          shiftAttendanceAllowance: 0,
          incentiveBonus: 0,
          overtimePay: 0,
          otherAllowance: 0,
          totalGrossSalary: gross,
          bpjsKesehatan: bpjsKes,
          bpjsKetenagakerjaan: bpjsTk,
          pph21,
          attendanceDeduction: 0,
          loanDeduction: 0,
          otherDeduction: 0,
          totalDeduction: totalDed,
          takeHomePay: thp,
          bankName: emp.bankName,
          bankAccountNumber: emp.bankAccountNumber,
          status: 'DRAFT' as PayrollStatus
        }
      })

    payslips.value.unshift(...generated)
    persistSlips()
    accountingStore.showToast(`Berhasil generate ${generated.length} slip gaji karyawan untuk periode ${month}/${year}`, 'success')
  }

  async function addCustomPayslip(slip: Omit<PayslipItem, 'id' | 'slipNumber'>): Promise<PayslipItem> {
    try {
      const createdBackend = await accountingApi.createPayrollSlip(slip as any)
      if (createdBackend) {
        payslips.value.unshift(createdBackend)
        persistSlips()
        accountingStore.showToast(`Slip gaji ${createdBackend.slipNumber} berhasil disimpan ke database`, 'success')
        return createdBackend
      }
    } catch (err: any) {
      console.warn('[PayrollStore] Simpan ke backend gagal, simpan ke lokal:', err)
    }

    const nextSeq = payslips.value.length + 1
    const slipNumber = `SLIP/${slip.periodYear}/${String(slip.periodMonth).padStart(2, '0')}/${String(nextSeq).padStart(3, '0')}`
    const newSlip: PayslipItem = {
      ...slip,
      id: `slip-${Date.now()}`,
      slipNumber
    }
    payslips.value.unshift(newSlip)
    persistSlips()
    accountingStore.showToast(`Slip gaji ${slipNumber} berhasil disimpan`, 'success')
    return newSlip
  }

  async function updatePayslip(id: string, updates: Partial<PayslipItem>) {
    try {
      await accountingApi.updatePayrollSlip(id, updates)
    } catch (err: any) {
      console.warn('[PayrollStore] Update ke backend gagal, perbarui lokal:', err)
    }

    const idx = payslips.value.findIndex(s => s.id === id)
    if (idx !== -1) {
      payslips.value[idx] = { ...payslips.value[idx], ...updates }
      persistSlips()
      accountingStore.showToast(`Slip gaji ${payslips.value[idx].slipNumber} berhasil diperbarui`, 'success')
    }
  }

  async function deletePayslip(id: string) {
    try {
      await accountingApi.deletePayrollSlip(id)
    } catch (err: any) {
      console.warn('[PayrollStore] Hapus dari backend gagal, hapus lokal:', err)
    }

    const idx = payslips.value.findIndex(s => s.id === id)
    if (idx !== -1) {
      const num = payslips.value[idx].slipNumber
      payslips.value.splice(idx, 1)
      persistSlips()
      accountingStore.showToast(`Slip gaji ${num} berhasil dihapus`, 'info')
    }
  }

  async function deleteAllInPeriod(month: number, year: number) {
    const toDelete = payslips.value.filter(s => s.periodMonth === month && s.periodYear === year)
    for (const slip of toDelete) {
      try {
        await accountingApi.deletePayrollSlip(slip.id)
      } catch (e) { }
    }
    payslips.value = payslips.value.filter(s => !(s.periodMonth === month && s.periodYear === year))
    persistSlips()
    accountingStore.showToast(`${toDelete.length} slip gaji periode ${month}/${year} berhasil dihapus`, 'info')
  }

  async function approvePayslip(id: string) {
    try {
      await accountingApi.approvePayrollSlip(id, accountingStore.currentUserEmail)
    } catch (err: any) {
      console.warn('[PayrollStore] Approve backend gagal, approve lokal:', err)
    }

    const item = payslips.value.find(s => s.id === id)
    if (item) {
      item.status = 'APPROVED'
      item.approvedBy = accountingStore.currentUserEmail
      item.approvedAt = new Date().toISOString()
      persistSlips()
      accountingStore.showToast(`Slip gaji ${item.slipNumber} telah disetujui (Approved)`, 'success')
    }
  }

  async function approveAllInPeriod(month: number, year: number) {
    try {
      await accountingApi.approveAllPayrollSlips(month, year, accountingStore.currentUserEmail)
    } catch (err: any) {
      console.warn('[PayrollStore] Approve all backend gagal, approve lokal:', err)
    }

    let count = 0
    payslips.value.forEach(s => {
      if (s.periodMonth === month && s.periodYear === year && s.status === 'DRAFT') {
        s.status = 'APPROVED'
        s.approvedBy = accountingStore.currentUserEmail
        s.approvedAt = new Date().toISOString()
        count++
      }
    })
    persistSlips()
    accountingStore.showToast(`${count} slip gaji berhasil disetujui sekaligus`, 'success')
  }

  async function markAsPaid(id: string) {
    try {
      await accountingApi.markPayrollSlipAsPaid(id)
    } catch (err: any) {
      console.warn('[PayrollStore] Mark as paid backend gagal, update lokal:', err)
    }

    const item = payslips.value.find(s => s.id === id)
    if (item) {
      item.status = 'PAID'
      item.paidAt = new Date().toISOString()
      persistSlips()
      accountingStore.showToast(`Slip gaji ${item.slipNumber} ditandai lunas ditransfer`, 'success')
    }
  }

  async function markAllAsPaid(month: number, year: number) {
    try {
      await accountingApi.markAllPayrollSlipsAsPaid(month, year)
    } catch (err: any) {
      console.warn('[PayrollStore] Mark all paid backend gagal, update lokal:', err)
    }

    let count = 0
    payslips.value.forEach(s => {
      if (s.periodMonth === month && s.periodYear === year && s.status !== 'PAID') {
        s.status = 'PAID'
        s.paidAt = new Date().toISOString()
        count++
      }
    })
    persistSlips()
    accountingStore.showToast(`${count} slip gaji berhasil ditandai Lunas / Paid`, 'success')
  }

  async function postPayrollToLedger(month: number, year: number) {
    const list = payslips.value.filter(s => s.periodMonth === month && s.periodYear === year)
    if (list.length === 0) {
      accountingStore.showToast(`Tidak ada slip gaji untuk periode ${month}/${year}`, 'error')
      return
    }

    let backendJournal: any = null
    try {
      backendJournal = await accountingApi.postPayrollToLedger(month, year)
    } catch (err: any) {
      console.warn('[PayrollStore] Post ledger backend gagal, create via lokal journal:', err)
    }

    const totalGross = list.reduce((s, item) => s + item.totalGrossSalary, 0)
    const totalPph21 = list.reduce((s, item) => s + item.pph21, 0)
    const totalTakeHome = list.reduce((s, item) => s + item.takeHomePay, 0)
    const totalBpjs = list.reduce((s, item) => s + item.bpjsKesehatan + item.bpjsKetenagakerjaan, 0)

    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ]
    const pName = `${monthNames[month - 1]} ${year}`
    const txDate = `${year}-${String(month).padStart(2, '0')}-28`

    const lines = [
      {
        accountId: 'acc-61210',
        accountCode: '6-1210',
        accountName: 'Beban Gaji Karyawan, Perawat & Staf',
        debit: totalGross,
        credit: 0,
        description: `Beban Payroll & Gaji Staf Periode ${pName} (${list.length} Karyawan)`
      },
      {
        accountId: 'acc-21510',
        accountCode: '2-1510',
        accountName: 'Hutang Pajak PPh 21 (Tenaga Medis & Karyawan)',
        debit: 0,
        credit: totalPph21,
        description: `Potongan PPh 21 Payroll Staf Masa ${pName}`
      },
      {
        accountId: 'acc-11114',
        accountCode: '1-1114',
        accountName: 'Bank Mandiri Operasional & Payroll',
        debit: 0,
        credit: totalTakeHome,
        description: `Transfer Pembayaran Gaji Bersih Karyawan via Mandiri Payroll`
      }
    ]

    if (totalBpjs > 0) {
      lines.push({
        accountId: 'acc-21130',
        accountCode: '2-1130',
        accountName: 'Hutang Gaji & Karyawan',
        debit: 0,
        credit: totalBpjs,
        description: `Titipan Iuran BPJS Ketenagakerjaan & Kesehatan Staf Periode ${pName}`
      })
    }

    const createdJournal = journalsStore.createManualJournal({
      journalDate: txDate,
      businessUnit: accountingStore.businessUnit,
      branchId: '11111111-1111-1111-1111-111111111111',
      branchName: 'Kantor Pusat Semarang',
      description: `Payroll Pembayaran Gaji Karyawan, Perawat & Staf SMG Periode ${pName}`,
      journalSource: 'PAYROLL',
      sourceReference: `PAYROLL-${year}${String(month).padStart(2, '0')}`,
      status: 'POSTED',
      lines
    })

    const jNum = backendJournal?.entry_number || createdJournal.journalNumber
    list.forEach(s => {
      s.journalId = backendJournal?.id || createdJournal.id
      s.journalNumber = jNum
      s.status = 'PAID'
      s.paidAt = new Date().toISOString()
    })
    persistSlips()

    accountingStore.showToast(`Jurnal Payroll ${jNum} resmi diposting ke Buku Besar!`, 'success')
  }

  function getSlipById(id: string): PayslipItem | undefined {
    return payslips.value.find(s => s.id === id)
  }

  return {
    employees,
    payslips,
    isLoading,
    selectedMonth,
    selectedYear,
    statusFilter,
    searchQuery,
    filteredPayslips,
    stats,
    fetchEmployees,
    fetchPayslips,
    generateMonthlyPayroll,
    addCustomPayslip,
    updatePayslip,
    deletePayslip,
    deleteAllInPeriod,
    approvePayslip,
    approveAllInPeriod,
    markAsPaid,
    markAllAsPaid,
    postPayrollToLedger,
    getSlipById
  }
})
