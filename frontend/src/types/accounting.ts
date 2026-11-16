// Core Accounting Type Definitions for Clinic ERP (Kidinest & SSC)

export type BusinessUnit = 'ALL' | 'KIDINEST' | 'SSC'

export type BranchCode = 'ALL' | 'BRANCH-MAIN' | 'BRANCH-JAKSEL' | 'BRANCH-SBY'

export type UserRole =
  | 'super_admin'
  | 'finance_manager'
  | 'accounting_supervisor'
  | 'accounting_staff'
  | 'tax_staff'
  | 'auditor'
  | 'management'

export type AccountType =
  | 'ASSET'
  | 'LIABILITY'
  | 'EQUITY'
  | 'REVENUE'
  | 'COGS'
  | 'EXPENSE'

export type NormalBalance = 'DEBIT' | 'CREDIT'

export interface Account {
  id: string
  code: string
  name: string
  type: AccountType
  normalBalance: NormalBalance
  level: number
  isHeader: boolean
  isActive: boolean
  businessUnit: 'ALL' | 'KIDINEST' | 'SSC'
  parentCode?: string
  description?: string
  currentBalance?: number
  taxCategory?: 'PPH_21' | 'PPH_23' | 'PPH_4_2' | 'PPN_OUT' | 'PPN_IN' | 'PPH_25' | 'PPH_29' | 'NON_TAX'
}

export type JournalStatus =
  | 'DRAFT'
  | 'SUBMITTED'
  | 'REVIEWED'
  | 'APPROVED'
  | 'POSTED'
  | 'REJECTED'
  | 'REVERSED'

export type JournalSource =
  | 'MANUAL'
  | 'CASHIER_CLOSING'
  | 'PATIENT_BILLING'
  | 'PATIENT_PAYMENT'
  | 'PHARMACY_SALES'
  | 'PHARMACY_PURCHASE'
  | 'INVENTORY_CONSUMPTION'
  | 'INVENTORY_ADJUSTMENT'
  | 'DOCTOR_FEE'
  | 'PAYROLL'
  | 'BANK_MUTATION'
  | 'FIXED_ASSET'
  | 'PREPAID_AMORTIZATION'
  | 'TAX_SETTLEMENT'

export interface JournalLine {
  id: string
  accountId: string
  accountCode: string
  accountName: string
  debit: number
  credit: number
  costCenterId?: string
  costCenterName?: string
  description: string
}

export interface JournalEntry {
  id: string
  journalNumber: string
  journalDate: string
  journalSource: JournalSource
  sourceReference: string
  businessUnit: BusinessUnit
  branchId: string
  branchName: string
  description: string
  totalDebit: number
  totalCredit: number
  status: JournalStatus
  lines: JournalLine[]
  createdBy: string
  createdByName: string
  createdAt: string
  reviewedBy?: string
  reviewedByName?: string
  reviewedAt?: string
  approvedBy?: string
  approvedByName?: string
  approvedAt?: string
  postedBy?: string
  postedByName?: string
  postedAt?: string
  reversalJournalId?: string
  reversalReason?: string
  isManual: boolean
  attachments?: string[]
}

// Operational Transaction Source System
export type SourceModule =
  | 'GRO_CASHIER'
  | 'PATIENT_BILLING'
  | 'PAYMENT_GATEWAY'
  | 'PHARMACY'
  | 'PURCHASING'
  | 'INVENTORY'
  | 'DOCTOR_SHARE'
  | 'PAYROLL'
  | 'BANK'
  | 'FIXED_ASSET'

export type MappingStatus = 'MAPPED' | 'UNMAPPED' | 'PARTIALLY_MAPPED'

export interface OperationalTransaction {
  id: string
  transactionNumber: string
  transactionDate: string
  sourceModule: SourceModule
  referenceNumber: string
  businessUnit: BusinessUnit
  branchName: string
  description: string
  amount: number
  mappingStatus: MappingStatus
  journalStatus: 'UNPROCESSED' | 'DRAFT_GENERATED' | 'POSTED'
  generatedJournalId?: string
  details: Record<string, any>
  suggestedDebitAccount: string
  suggestedCreditAccount: string
}

// Accounting Mapping Rule
export interface AccountingMappingRule {
  id: string
  eventKey: string
  eventName: string
  businessUnit: BusinessUnit
  sourceModule: SourceModule
  debitAccountCode: string
  debitAccountName: string
  creditAccountCode: string
  creditAccountName: string
  cogsAccountCode?: string
  inventoryAccountCode?: string
  taxAccountCode?: string
  description: string
  isActive: boolean
}

// General Ledger Transaction
export interface LedgerEntry {
  id: string
  date: string
  journalNumber: string
  journalId: string
  source: JournalSource
  description: string
  accountCode: string
  accountName: string
  businessUnit: BusinessUnit
  branchName: string
  debit: number
  credit: number
  runningBalance: number
}

// Trial Balance Line
export interface TrialBalanceRow {
  accountCode: string
  accountName: string
  accountType: AccountType
  openingDebit: number
  openingCredit: number
  periodDebit: number
  periodCredit: number
  endingDebit: number
  endingCredit: number
}

// Fiscal Period
export interface AccountingPeriod {
  id: string
  year: number
  month: number
  periodName: string
  startDate: string
  endDate: string
  status: 'OPEN' | 'CLOSING_IN_PROGRESS' | 'LOCKED'
  lockedAt?: string
  lockedBy?: string
  unpostedJournalCount: number
  pendingApprovalCount: number
}

// Tax Transaction
export interface TaxTransaction {
  id: string
  taxNumber: string
  taxDate: string
  businessUnit: BusinessUnit
  taxCategory: 'PPN_OUTPUT' | 'PPN_INPUT' | 'PPH_21' | 'PPH_23' | 'PPH_4_2' | 'PPH_25' | 'PPH_29'
  accountCode?: string
  journalId?: string
  journalNumber?: string
  counterpartyName: string
  counterpartyNPWP: string
  counterpartyNIK?: string
  objectCode: string
  grossAmount: number
  dppAmount: number
  taxRate: number
  taxAmount: number
  documentReference: string
  reconciliationStatus: 'RECONCILED' | 'UNRECONCILED'
  isReported: boolean
  periodMonth: number
  periodYear: number
}

// Financial Audit Log
export interface AuditLog {
  id: string
  timestamp: string
  userId: string
  userName: string
  userRole: string
  action: string
  module: string
  entityId: string
  entityReference: string
  beforeValue?: any
  afterValue?: any
  ipAddress: string
  userAgent: string
}

// Cash & Bank Management Types (BKM, BKK, Transfer Antar Rekening)
export type CashBankTransactionType = 'CASH_IN' | 'CASH_OUT' | 'TRANSFER'

export interface CashBankTransaction {
  id: string
  transactionNumber: string
  date: string
  type: CashBankTransactionType
  accountCode: string
  accountName: string
  counterpartAccountCode: string
  counterpartAccountName: string
  amount: number
  adminFee?: number
  description: string
  referenceNo?: string
  recipientOrPayer?: string
  businessUnit: BusinessUnit
  journalId?: string
  journalNumber?: string
  createdAt: string
  createdBy: string
}
