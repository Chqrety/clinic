// REST API Client for Clinic Accounting & Financial ERP
// Communicates directly with PostgreSQL-backed Go backend at /api/accounting/...

import type {
  Account,
  JournalEntry,
  CashBankTransaction,
  BusinessUnit
} from '@/types/accounting'
import type { Employee, PayslipItem } from '@/types/payroll'

const API_BASE = '/api/accounting'

interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
  meta?: {
    current_page: number
    per_page: number
    total_items: number
    total_pages: number
  }
  error?: string
}

function getValidAuthToken(): string {
  return localStorage.getItem('clinic_auth_token') || localStorage.getItem('clinic_token') || ''
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
  const token = getValidAuthToken()
  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      'X-Branch-ID': localStorage.getItem('clinic_branch_id') || 'ALL',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options?.headers || {}),
    },
  })

  if (res.status === 401) {
    localStorage.removeItem('clinic_auth_token')
    localStorage.removeItem('clinic_token')
    localStorage.removeItem('clinic_auth_user')
    let errMsg = 'Sesi login telah berakhir (401), silakan login kembali.'
    try {
      const errJson = await res.json()
      if (errJson.message) errMsg = errJson.message
      else if (errJson.error) errMsg = errJson.error
    } catch {}
    throw new Error(errMsg)
  }

  if (!res.ok) {
    let errMsg = `HTTP Error ${res.status} (${res.statusText})`
    try {
      const errJson = await res.json()
      if (errJson.message) errMsg = errJson.message
      else if (errJson.error) errMsg = errJson.error
    } catch {
      // ignore json parse error
    }
    throw new Error(errMsg)
  }

  const json: ApiResponse<T> = await res.json()
  return json
}

export interface BackendAccount {
  id: string
  code: string
  name: string
  type: 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'COGS' | 'EXPENSE'
  normal_balance: 'DEBIT' | 'CREDIT'
  parent_id?: string
  level: number
  is_header: boolean
  is_active: boolean
  description?: string
}

export interface BackendLedgerTransaction {
  date: string
  journal_id: string
  journal_number: string
  journal_type: string
  reference_type: string
  reference_id: string
  description: string
  debit: string | number
  credit: string | number
  running_balance: string | number
}

export interface BackendGeneralLedger {
  account_id: string
  account_code: string
  account_name: string
  account_type: string
  normal_balance: 'DEBIT' | 'CREDIT'
  opening_balance: string | number
  total_debit: string | number
  total_credit: string | number
  closing_balance: string | number
  transactions: BackendLedgerTransaction[] | null
}

export interface BackendCashBankAccount {
  id: string
  branch_id: string
  account_id: string
  account_code: string
  code: string
  name: string
  type: 'CASH' | 'BANK'
  bank_name?: string
  account_number?: string
  is_active: boolean
  created_at: string
}

export interface BackendCashBankTransaction {
  id: string
  branch_id: string
  cash_bank_account_id: string
  cash_bank_account_code?: string
  cash_bank_account_name?: string
  transaction_number: string
  transaction_type: 'CASH_IN' | 'CASH_OUT' | 'TRANSFER'
  amount: string | number
  transaction_date: string
  counterpart_account_id: string
  counterpart_code?: string
  counterpart_name?: string
  description: string
  reference_no?: string
  journal_id?: string
  created_at: string
}

export interface BackendJournalLine {
  id: string
  account_id: string
  account_code: string
  account_name: string
  debit: string | number
  credit: string | number
  description: string
  line_order: number
}

export interface BackendJournal {
  id: string
  branch_id: string
  journal_number: string
  journal_type: string
  transaction_date: string
  posting_date?: string
  status: string
  reference_type: string
  reference_id: string
  description: string
  total_debit: string | number
  total_credit: string | number
  lines?: BackendJournalLine[]
  created_at: string
}

export interface BackendTaxWithholding {
  id: string
  branch_id: string
  tax_type: string
  bupot_number: string
  taxpayer_nik?: string
  taxpayer_npwp?: string
  taxpayer_name: string
  object_code: string
  gross_amount: string | number
  dpp_amount: string | number
  tax_rate: string | number
  tax_amount: string | number
  period_year: number
  period_month: number
  is_reported: boolean
  journal_id?: string
  created_at?: string
}

export interface BackendTrialBalanceItem {
  account_id: string
  account_code: string
  account_name: string
  account_type: string
  normal_balance: 'DEBIT' | 'CREDIT'
  opening_debit: string | number
  opening_credit: string | number
  period_debit: string | number
  period_credit: string | number
  ending_debit: string | number
  ending_credit: string | number
}

export interface BackendTrialBalance {
  branch_id?: string
  start_date: string
  end_date: string
  total_opening_debit: string | number
  total_opening_credit: string | number
  total_period_debit: string | number
  total_period_credit: string | number
  total_ending_debit: string | number
  total_ending_credit: string | number
  is_balanced: boolean
  items: BackendTrialBalanceItem[]
}

export const accountingApi = {
  // =========================================================================
  // 1. Chart of Accounts (COA)
  // =========================================================================
  async getAccounts(params?: { search?: string; limit?: number; offset?: number }): Promise<BackendAccount[]> {
    const q = new URLSearchParams()
    if (params?.search) q.set('search', params.search)
    q.set('limit', String(params?.limit || 200))
    if (params?.offset) q.set('offset', String(params.offset))

    const res = await request<BackendAccount[]>(`/accounts?${q.toString()}`)
    return res.data || []
  },

  async createAccount(data: {
    code: string
    name: string
    type: string
    normal_balance: string
    parent_id?: string
    level: number
    is_header: boolean
    description?: string
  }): Promise<BackendAccount> {
    const res = await request<BackendAccount>('/accounts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async updateAccount(id: string, data: {
    name?: string
    is_active?: boolean
    description?: string
  }): Promise<BackendAccount> {
    const res = await request<BackendAccount>(`/accounts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
    return res.data
  },

  // =========================================================================
  // 2. Cash & Bank Accounts & Transactions
  // =========================================================================
  async getCashBankAccounts(): Promise<BackendCashBankAccount[]> {
    const res = await request<BackendCashBankAccount[]>('/cash/accounts?limit=100')
    return res.data || []
  },

  async createCashBankAccount(data: {
    account_id: string
    code: string
    name: string
    type: 'CASH' | 'BANK'
    bank_name?: string
    account_number?: string
  }): Promise<BackendCashBankAccount> {
    const res = await request<BackendCashBankAccount>('/cash/accounts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async getCashBankTransactions(params?: {
    search?: string
    limit?: number
    offset?: number
    type?: string
  }): Promise<BackendCashBankTransaction[]> {
    const q = new URLSearchParams()
    if (params?.search) q.set('search', params.search)
    q.set('limit', String(params?.limit || 200))
    if (params?.offset) q.set('offset', String(params.offset))
    if (params?.type) q.set('type', params.type)

    const res = await request<BackendCashBankTransaction[]>(`/cash/transactions?${q.toString()}`)
    return res.data || []
  },

  async recordCashIn(data: {
    cash_bank_account_id: string
    amount: number
    transaction_date: string
    counterpart_account_id: string
    description: string
    reference_no?: string
  }): Promise<{ transaction: BackendCashBankTransaction; journal_number?: string }> {
    const res = await request<{ transaction: BackendCashBankTransaction; journal_number?: string }>('/cash/in', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async recordCashOut(data: {
    cash_bank_account_id: string
    amount: number
    transaction_date: string
    counterpart_account_id: string
    description: string
    reference_no?: string
  }): Promise<{ transaction: BackendCashBankTransaction; journal_number?: string }> {
    const res = await request<{ transaction: BackendCashBankTransaction; journal_number?: string }>('/cash/out', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async transferCashBank(data: {
    source_account_id: string
    target_account_id: string
    amount: number
    transaction_date: string
    description: string
    reference_no?: string
  }): Promise<{ transaction: BackendCashBankTransaction; journal_number?: string }> {
    const res = await request<{ transaction: BackendCashBankTransaction; journal_number?: string }>('/cash/transfer', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  // =========================================================================
  // 3. Journals
  // =========================================================================
  async getJournals(params?: {
    search?: string
    limit?: number
    offset?: number
    type?: string
    status?: string
  }): Promise<BackendJournal[]> {
    const q = new URLSearchParams()
    if (params?.search) q.set('search', params.search)
    q.set('limit', String(params?.limit || 100))
    if (params?.offset) q.set('offset', String(params.offset))
    if (params?.type) q.set('type', params.type)
    if (params?.status) q.set('status', params.status)

    const res = await request<BackendJournal[]>(`/journals?${q.toString()}`)
    return res.data || []
  },

  async getJournalByID(id: string): Promise<BackendJournal> {
    const res = await request<BackendJournal>(`/journals/${id}`)
    return res.data
  },

  async createJournal(data: {
    journal_type: string
    transaction_date: string
    reference_type: string
    reference_id: string
    description: string
    auto_post?: boolean
    lines: Array<{
      account_id?: string
      account_code?: string
      debit: number
      credit: number
      description?: string
    }>
  }): Promise<BackendJournal> {
    const res = await request<BackendJournal>('/journals', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async postJournal(id: string): Promise<BackendJournal> {
    const res = await request<BackendJournal>(`/journals/${id}/post`, {
      method: 'POST',
    })
    return res.data
  },

  // =========================================================================
  // 4. Taxes & Withholdings (e-Bupot & SPT Masa/Tahunan)
  // =========================================================================
  async getTaxWithholdings(params?: {
    year?: number
    month?: number
    tax_type?: string
    is_reported?: boolean
    limit?: number
  }): Promise<BackendTaxWithholding[]> {
    const q = new URLSearchParams()
    if (params?.year) q.set('year', String(params.year))
    if (params?.month) q.set('month', String(params.month))
    if (params?.tax_type) q.set('tax_type', params.tax_type)
    if (params?.is_reported !== undefined) q.set('is_reported', String(params.is_reported))
    q.set('limit', String(params?.limit || 100))

    const res = await request<BackendTaxWithholding[]>(`/taxes/withholdings?${q.toString()}`)
    return res.data || []
  },

  async createTaxWithholding(data: {
    tax_type: string
    taxpayer_nik?: string
    taxpayer_npwp?: string
    taxpayer_name: string
    object_code: string
    gross_amount: number
    dpp_amount: number
    tax_rate: number
    tax_amount: number
    period_year: number
    period_month: number
    journal_id?: string
  }): Promise<BackendTaxWithholding> {
    const res = await request<BackendTaxWithholding>('/taxes/withholdings', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return res.data
  },

  async getTaxSummary(year?: number, month?: number): Promise<any> {
    const q = new URLSearchParams()
    if (year) q.set('year', String(year))
    if (month) q.set('month', String(month))
    const res = await request<any>(`/taxes/summary?${q.toString()}`)
    return res.data
  },

  async markTaxAsReported(id: string): Promise<void> {
    await request<void>(`/taxes/withholdings/${id}/report`, {
      method: 'POST',
    })
  },

  // =========================================================================
  // 5. Financial Reports & Statements
  // =========================================================================
  async getGeneralLedger(
    accountCode: string,
    startDate?: string,
    endDate?: string
  ): Promise<BackendGeneralLedger | null> {
    const q = new URLSearchParams()
    q.set('account_code', accountCode)
    if (startDate) q.set('start_date', startDate)
    if (endDate) q.set('end_date', endDate)
    const res = await request<BackendGeneralLedger>(`/reports/general-ledger?${q.toString()}`)
    return res.data
  },

  async getTrialBalance(startDate?: string, endDate?: string): Promise<BackendTrialBalance | null> {
    const q = new URLSearchParams()
    if (startDate) q.set('start_date', startDate)
    if (endDate) q.set('end_date', endDate)
    const res = await request<BackendTrialBalance>(`/reports/trial-balance?${q.toString()}`)
    return res.data
  },

  async getIncomeStatement(startDate?: string, endDate?: string): Promise<any> {
    const q = new URLSearchParams()
    if (startDate) q.set('start_date', startDate)
    if (endDate) q.set('end_date', endDate)
    const res = await request<any>(`/reports/income-statement?${q.toString()}`)
    return res.data
  },

  async getBalanceSheet(asOfDate?: string): Promise<any> {
    const q = new URLSearchParams()
    if (asOfDate) q.set('as_of_date', asOfDate)
    const res = await request<any>(`/reports/balance-sheet?${q.toString()}`)
    return res.data
  },

  async getCashFlowStatement(startDate?: string, endDate?: string): Promise<any> {
    const q = new URLSearchParams()
    if (startDate) q.set('start_date', startDate)
    if (endDate) q.set('end_date', endDate)
    const res = await request<any>(`/reports/cash-flow?${q.toString()}`)
    return res.data
  },

  async getReceivableAging(asOfDate?: string): Promise<any> {
    const q = new URLSearchParams()
    if (asOfDate) q.set('as_of_date', asOfDate)
    const res = await request<any>(`/receivables/aging?${q.toString()}`)
    return res.data
  },

  async getPayableAging(asOfDate?: string): Promise<any> {
    const q = new URLSearchParams()
    if (asOfDate) q.set('as_of_date', asOfDate)
    const res = await request<any>(`/payables/aging?${q.toString()}`)
    return res.data
  },

  // =========================================================================
  // 6. Staff Payroll & Payslips
  // =========================================================================
  async getPayrollEmployees(businessUnit?: string, activeOnly = true): Promise<Employee[]> {
    const q = new URLSearchParams()
    if (businessUnit && businessUnit !== 'ALL') q.set('business_unit', businessUnit)
    if (activeOnly) q.set('active_only', 'true')
    const res = await request<any[]>(`/payroll/employees?${q.toString()}`)
    const list = res.data || []
    return list.map(emp => ({
      ...emp,
      baseSalary: Number(emp.baseSalary) || 0,
      positionAllowance: Number(emp.positionAllowance) || 0,
      transportMealAllowance: Number(emp.transportMealAllowance) || 0,
      dailyMealTransportRate: Number(emp.dailyMealTransportRate) || (emp.transportMealAllowance ? Math.round(Number(emp.transportMealAllowance) / 22) : 0),
    }))
  },

  async createPayrollEmployee(emp: Partial<Employee>): Promise<Employee> {
    const res = await request<Employee>('/payroll/employees', {
      method: 'POST',
      body: JSON.stringify(emp),
    })
    return res.data
  },

  async updatePayrollEmployee(id: string, emp: Partial<Employee>): Promise<Employee> {
    const res = await request<Employee>(`/payroll/employees/${id}`, {
      method: 'PUT',
      body: JSON.stringify(emp),
    })
    return res.data
  },

  async getPayrollSlips(month?: number, year?: number, businessUnit?: string, status?: string): Promise<PayslipItem[]> {
    const q = new URLSearchParams()
    if (month) q.set('month', String(month))
    if (year) q.set('year', String(year))
    if (businessUnit && businessUnit !== 'ALL') q.set('business_unit', businessUnit)
    if (status && status !== 'ALL') q.set('status', status)
    const res = await request<any[]>(`/payroll/slips?${q.toString()}`)
    const list = res.data || []
    return list.map(slip => ({
      ...slip,
      baseSalary: Number(slip.baseSalary) || 0,
      positionAllowance: Number(slip.positionAllowance) || 0,
      transportMealAllowance: Number(slip.transportMealAllowance) || 0,
      shiftAttendanceAllowance: Number(slip.shiftAttendanceAllowance) || 0,
      incentiveBonus: Number(slip.incentiveBonus) || 0,
      overtimePay: Number(slip.overtimePay) || 0,
      otherAllowance: Number(slip.otherAllowance) || 0,
      totalGrossSalary: Number(slip.totalGrossSalary) || 0,
      bpjsKesehatan: Number(slip.bpjsKesehatan) || 0,
      bpjsKetenagakerjaan: Number(slip.bpjsKetenagakerjaan) || 0,
      pph21: Number(slip.pph21) || 0,
      attendanceDeduction: Number(slip.attendanceDeduction) || 0,
      loanDeduction: Number(slip.loanDeduction) || 0,
      otherDeduction: Number(slip.otherDeduction) || 0,
      totalDeduction: Number(slip.totalDeduction) || 0,
      takeHomePay: Number(slip.takeHomePay) || 0,
      dailyMealTransportRate: Number(slip.dailyMealTransportRate) || 0,
    }))
  },

  async createPayrollSlip(slip: Partial<PayslipItem>): Promise<PayslipItem> {
    const res = await request<PayslipItem>('/payroll/slips', {
      method: 'POST',
      body: JSON.stringify(slip),
    })
    return res.data
  },

  async updatePayrollSlip(id: string, slip: Partial<PayslipItem>): Promise<PayslipItem> {
    const res = await request<PayslipItem>(`/payroll/slips/${id}`, {
      method: 'PUT',
      body: JSON.stringify(slip),
    })
    return res.data
  },

  async deletePayrollSlip(id: string): Promise<void> {
    await request<void>(`/payroll/slips/${id}`, {
      method: 'DELETE',
    })
  },

  async approvePayrollSlip(id: string, approvedBy = 'Finance Manager'): Promise<void> {
    await request<void>(`/payroll/slips/${id}/approve`, {
      method: 'POST',
      body: JSON.stringify({ approvedBy }),
    })
  },

  async approveAllPayrollSlips(month: number, year: number, approvedBy = 'Finance Manager'): Promise<void> {
    await request<void>('/payroll/slips/approve-all', {
      method: 'POST',
      body: JSON.stringify({ month, year, approvedBy }),
    })
  },

  async markPayrollSlipAsPaid(id: string): Promise<void> {
    await request<void>(`/payroll/slips/${id}/pay`, {
      method: 'POST',
    })
  },

  async markAllPayrollSlipsAsPaid(month: number, year: number): Promise<void> {
    await request<void>('/payroll/slips/pay-all', {
      method: 'POST',
      body: JSON.stringify({ month, year }),
    })
  },

  async postPayrollToLedger(month: number, year: number): Promise<any> {
    const res = await request<any>('/payroll/post-ledger', {
      method: 'POST',
      body: JSON.stringify({ month, year }),
    })
    return res.data
  }
}

