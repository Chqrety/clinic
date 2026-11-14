// REST API Client for Clinic Accounting & Financial ERP
// Communicates directly with PostgreSQL-backed Go backend at /api/v1/...

import type {
  Account,
  JournalEntry,
  CashBankTransaction,
  BusinessUnit
} from '@/types/accounting'

const API_BASE = '/api/v1'

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

async function request<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      'X-Branch-ID': '11111111-1111-1111-1111-111111111111',
      ...(options?.headers || {}),
    },
  })

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
  }
}
