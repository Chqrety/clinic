import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CashBankTransaction, CashBankTransactionType, BusinessUnit, Account } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'
import { accountingApi, type BackendCashBankAccount, type BackendCashBankTransaction } from '@/api/accounting'
import { formatIDR } from '@/utils/currency'

export const useCashBankStore = defineStore('cashbank', () => {
  const accountingStore = useAccountingStore()
  const ledgerStore = useLedgerStore()
  const journalsStore = useJournalsStore()

  // Transactions list
  // Real transactions fetched dynamically from PostgreSQL database via Go API
  const transactions = ref<CashBankTransaction[]>([])

  const isLoading = ref(false)
  const registeredCashBankAccounts = ref<BackendCashBankAccount[]>([])

  // Fetch Cash & Bank accounts & transactions directly from PostgreSQL database
  async function fetchCashBankData() {
    try {
      isLoading.value = true
      const [backendCba, backendTxs] = await Promise.all([
        accountingApi.getCashBankAccounts(),
        accountingApi.getCashBankTransactions({ limit: 200 })
      ])

      if (backendCba && backendCba.length > 0) {
        registeredCashBankAccounts.value = backendCba
      }

      if (backendTxs && backendTxs.length > 0) {
        transactions.value = backendTxs.map(t => {
          const cba = registeredCashBankAccounts.value.find(c => c.id === t.cash_bank_account_id)
          const accCode = t.cash_bank_account_code || cba?.account_code || '1102.001'
          const accName = t.cash_bank_account_name || cba?.name || 'Kas / Bank Operasional'

          return {
            id: t.id,
            transactionNumber: t.transaction_number,
            date: t.transaction_date,
            type: t.transaction_type,
            accountCode: accCode,
            accountName: accName,
            counterpartAccountCode: t.counterpart_code || '',
            counterpartAccountName: t.counterpart_name || 'Kontra Akun',
            amount: Number(t.amount),
            description: t.description,
            referenceNo: t.reference_no || '-',
            recipientOrPayer: t.counterpart_name || '-',
            businessUnit: 'ALL',
            journalId: t.journal_id,
            createdAt: t.created_at,
            createdBy: 'Rian Pratama, S.Ak.'
          }
        })
      }
    } catch (err: any) {
      console.error('Failed to fetch cash bank data from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load
  fetchCashBankData()

  // Get all active cash accounts (parentCode === '1101' or startsWith '1101')
  const cashAccounts = computed(() => {
    return ledgerStore.accounts.filter(a => {
      if (a.isHeader || a.isActive === false) return false
      if (accountingStore.businessUnit !== 'ALL' && a.businessUnit && a.businessUnit !== 'ALL' && a.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return a.parentCode === '1101' || a.code.startsWith('1101.') || a.code === '1101'
    })
  })

  // Get all active bank accounts (parentCode === '1102' or startsWith '1102' or '1103')
  const bankAccounts = computed(() => {
    return ledgerStore.accounts.filter(a => {
      if (a.isHeader || a.isActive === false) return false
      if (accountingStore.businessUnit !== 'ALL' && a.businessUnit && a.businessUnit !== 'ALL' && a.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return (
        a.parentCode === '1102' ||
        a.parentCode === '1103' ||
        a.code.startsWith('1102.') ||
        a.code.startsWith('1103.') ||
        a.code === '1102' ||
        a.code === '1103'
      )
    })
  })

  // All Cash and Bank accounts combined (Leaf accounts under 11)
  const allCashBankAccounts = computed(() => {
    return [...cashAccounts.value, ...bankAccounts.value]
  })

  // Action: Catat Kas/Bank Masuk (BKM / BBM)
  function recordCashIn(data: {
    accountCode: string
    counterpartAccountCode: string
    amount: number
    transactionDate: string
    recipientOrPayer: string
    description: string
    referenceNo?: string
    businessUnit?: BusinessUnit
  }) {
    const acc = ledgerStore.accounts.find(a => a.code === data.accountCode)
    const counterpartAcc = ledgerStore.accounts.find(a => a.code === data.counterpartAccountCode)

    const isBank = data.accountCode.startsWith('1102') || data.accountCode.startsWith('1103')
    const prefix = isBank ? 'BBM' : 'BKM'
    const seq = transactions.value.filter(t => t.type === 'CASH_IN').length + 1
    const txNumber = `${prefix}-${data.transactionDate.slice(0, 7)}-${String(seq).padStart(3, '0')}`

    const bu = data.businessUnit || accountingStore.businessUnit

    // Create journal entry automatically
    const journal = journalsStore.createManualJournal({
      journalDate: data.transactionDate,
      businessUnit: bu,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Penerimaan Kas/Bank [${txNumber}]: ${data.description}`,
      journalSource: 'MANUAL',
      sourceReference: txNumber,
      status: 'POSTED',
      lines: [
        {
          accountCode: data.accountCode,
          accountName: acc?.name || `Buku Kas/Bank (${data.accountCode})`,
          debit: data.amount,
          credit: 0,
          description: data.description
        },
        {
          accountCode: data.counterpartAccountCode,
          accountName: counterpartAcc?.name || `Kontra Akun (${data.counterpartAccountCode})`,
          debit: 0,
          credit: data.amount,
          description: `Diterima dari: ${data.recipientOrPayer}`
        }
      ]
    })

    const newTx: CashBankTransaction = {
      id: 'cb-' + Date.now(),
      transactionNumber: txNumber,
      date: data.transactionDate,
      type: 'CASH_IN',
      accountCode: data.accountCode,
      accountName: acc?.name || data.accountCode,
      counterpartAccountCode: data.counterpartAccountCode,
      counterpartAccountName: counterpartAcc?.name || data.counterpartAccountCode,
      amount: data.amount,
      description: data.description,
      referenceNo: data.referenceNo || '-',
      recipientOrPayer: data.recipientOrPayer,
      businessUnit: bu,
      journalId: journal.id,
      journalNumber: journal.journalNumber,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    transactions.value.unshift(newTx)
    accountingStore.showToast(`Penerimaan ${txNumber} berhasil dibukukan ke ${data.accountCode}!`, 'success')
    return newTx
  }

  // Action: Catat Kas/Bank Keluar (BKK / BBK)
  function recordCashOut(data: {
    accountCode: string
    counterpartAccountCode: string
    amount: number
    transactionDate: string
    recipientOrPayer: string
    description: string
    referenceNo?: string
    businessUnit?: BusinessUnit
  }) {
    const acc = ledgerStore.accounts.find(a => a.code === data.accountCode)
    const counterpartAcc = ledgerStore.accounts.find(a => a.code === data.counterpartAccountCode)

    const isBank = data.accountCode.startsWith('1102') || data.accountCode.startsWith('1103')
    const prefix = isBank ? 'BBK' : 'BKK'
    const seq = transactions.value.filter(t => t.type === 'CASH_OUT').length + 1
    const txNumber = `${prefix}-${data.transactionDate.slice(0, 7)}-${String(seq).padStart(3, '0')}`

    const bu = data.businessUnit || accountingStore.businessUnit

    // Create journal entry automatically
    const journal = journalsStore.createManualJournal({
      journalDate: data.transactionDate,
      businessUnit: bu,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Pengeluaran Kas/Bank [${txNumber}]: ${data.description}`,
      journalSource: 'MANUAL',
      sourceReference: txNumber,
      status: 'POSTED',
      lines: [
        {
          accountCode: data.counterpartAccountCode,
          accountName: counterpartAcc?.name || `Kontra Akun (${data.counterpartAccountCode})`,
          debit: data.amount,
          credit: 0,
          description: `Dibayarkan kepada: ${data.recipientOrPayer}`
        },
        {
          accountCode: data.accountCode,
          accountName: acc?.name || `Buku Kas/Bank (${data.accountCode})`,
          debit: 0,
          credit: data.amount,
          description: data.description
        }
      ]
    })

    const newTx: CashBankTransaction = {
      id: 'cb-' + Date.now(),
      transactionNumber: txNumber,
      date: data.transactionDate,
      type: 'CASH_OUT',
      accountCode: data.accountCode,
      accountName: acc?.name || data.accountCode,
      counterpartAccountCode: data.counterpartAccountCode,
      counterpartAccountName: counterpartAcc?.name || data.counterpartAccountCode,
      amount: data.amount,
      description: data.description,
      referenceNo: data.referenceNo || '-',
      recipientOrPayer: data.recipientOrPayer,
      businessUnit: bu,
      journalId: journal.id,
      journalNumber: journal.journalNumber,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    transactions.value.unshift(newTx)
    accountingStore.showToast(`Pengeluaran ${txNumber} berhasil dibukukan dari ${data.accountCode}!`, 'success')
    return newTx
  }

  // Action: Transfer Antar Rekening Kas & Bank (Mutasi / Pemindahbukuan)
  function recordTransfer(data: {
    sourceAccountCode: string
    targetAccountCode: string
    amount: number
    adminFee?: number
    transactionDate: string
    description: string
    referenceNo?: string
    businessUnit?: BusinessUnit
  }) {
    if (data.sourceAccountCode === data.targetAccountCode) {
      throw new Error('Rekening asal dan rekening tujuan tidak boleh sama')
    }

    const sourceAcc = ledgerStore.accounts.find(a => a.code === data.sourceAccountCode)
    const targetAcc = ledgerStore.accounts.find(a => a.code === data.targetAccountCode)
    const fee = data.adminFee || 0
    const totalOut = data.amount + fee

    const seq = transactions.value.filter(t => t.type === 'TRANSFER').length + 1
    const txNumber = `TRF-${data.transactionDate.slice(0, 7)}-${String(seq).padStart(3, '0')}`

    const bu = data.businessUnit || accountingStore.businessUnit

    const lines = [
      {
        accountCode: data.targetAccountCode,
        accountName: targetAcc?.name || `Rekening Penerima (${data.targetAccountCode})`,
        debit: data.amount,
        credit: 0,
        description: `Transfer masuk dari ${data.sourceAccountCode}: ${data.description}`
      }
    ]

    if (fee > 0) {
      lines.push({
        accountCode: '7103',
        accountName: 'BEBAN BIAYA ADMINISTRASI BANK',
        debit: fee,
        credit: 0,
        description: `Biaya transfer antar bank via ${data.sourceAccountCode}`
      })
    }

    lines.push({
      accountCode: data.sourceAccountCode,
      accountName: sourceAcc?.name || `Rekening Pengirim (${data.sourceAccountCode})`,
      debit: 0,
      credit: totalOut,
      description: `Transfer keluar ke ${data.targetAccountCode}: ${data.description}`
    })

    const journal = journalsStore.createManualJournal({
      journalDate: data.transactionDate,
      businessUnit: bu,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Transfer Antar Kas/Bank [${txNumber}]: ${data.description}`,
      journalSource: 'BANK_MUTATION',
      sourceReference: txNumber,
      status: 'POSTED',
      lines
    })

    const newTx: CashBankTransaction = {
      id: 'cb-' + Date.now(),
      transactionNumber: txNumber,
      date: data.transactionDate,
      type: 'TRANSFER',
      accountCode: data.sourceAccountCode,
      accountName: sourceAcc?.name || data.sourceAccountCode,
      counterpartAccountCode: data.targetAccountCode,
      counterpartAccountName: targetAcc?.name || data.targetAccountCode,
      amount: data.amount,
      adminFee: fee,
      description: data.description,
      referenceNo: data.referenceNo || '-',
      recipientOrPayer: `Pindah buku: ${sourceAcc?.name} -> ${targetAcc?.name}`,
      businessUnit: bu,
      journalId: journal.id,
      journalNumber: journal.journalNumber,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    transactions.value.unshift(newTx)
    accountingStore.showToast(`Mutasi transfer ${txNumber} berhasil dibukukan!`, 'success')
    return newTx
  }

  // Action: Catat Voucher Multi-Line Kas & Bank (BKM / BKK / TRF ala Jurnal Umum)
  function recordMultiLineVoucher(data: {
    voucherType: 'CASH_IN' | 'CASH_OUT' | 'TRANSFER'
    primaryAccountCode: string
    transactionDate: string
    recipientOrPayer: string
    description: string
    referenceNo?: string
    businessUnit?: BusinessUnit
    lines: {
      accountCode: string
      accountName?: string
      debit: number
      credit: number
      description?: string
    }[]
  }) {
    const totalDebit = data.lines.reduce((s, l) => s + (Number(l.debit) || 0), 0)
    const totalCredit = data.lines.reduce((s, l) => s + (Number(l.credit) || 0), 0)

    if (Math.abs(totalDebit - totalCredit) > 0.01) {
      throw new Error(`Total Debit (${totalDebit}) tidak seimbang dengan Total Kredit (${totalCredit})`)
    }

    const primaryAcc = ledgerStore.accounts.find(a => a.code === data.primaryAccountCode)
    const primaryName = primaryAcc?.name || data.primaryAccountCode
    const isBank = data.primaryAccountCode.startsWith('1102') || data.primaryAccountCode.startsWith('1103')

    let prefix = 'BKK'
    if (data.voucherType === 'CASH_IN') {
      prefix = isBank ? 'BBM' : 'BKM'
    } else if (data.voucherType === 'CASH_OUT') {
      prefix = isBank ? 'BBK' : 'BKK'
    } else {
      prefix = 'TRF'
    }

    const seq = transactions.value.filter(t => t.type === data.voucherType).length + 1
    const txNumber = `${prefix}-${data.transactionDate.slice(0, 7)}-${String(seq).padStart(3, '0')}`
    const bu = data.businessUnit || accountingStore.businessUnit

    // Calculate total flow amount on primary cash/bank account
    const primaryLines = data.lines.filter(l => l.accountCode === data.primaryAccountCode)
    let amount = 0
    if (data.voucherType === 'CASH_IN') {
      amount = primaryLines.reduce((s, l) => s + (Number(l.debit) || 0), 0) || totalDebit
    } else if (data.voucherType === 'CASH_OUT') {
      amount = primaryLines.reduce((s, l) => s + (Number(l.credit) || 0), 0) || totalCredit
    } else {
      amount = totalDebit
    }

    // Resolve counterpart account summary
    const counterpartLines = data.lines.filter(l => l.accountCode !== data.primaryAccountCode)
    let counterpartCode = 'MULTIPLE'
    let counterpartName = 'Multi-Akun Lawan'
    if (counterpartLines.length === 1) {
      counterpartCode = counterpartLines[0].accountCode
      counterpartName = counterpartLines[0].accountName || ledgerStore.accounts.find(a => a.code === counterpartCode)?.name || counterpartCode
    } else if (counterpartLines.length > 1) {
      counterpartCode = 'MULTIPLE'
      counterpartName = `${counterpartLines.length} Akun Lawan (${counterpartLines.map(c => c.accountCode).slice(0, 2).join(', ')}...)`
    }

    // Map lines with proper account names
    const resolvedLines = data.lines.map(l => {
      const acc = ledgerStore.accounts.find(a => a.code === l.accountCode)
      return {
        accountCode: l.accountCode,
        accountName: l.accountName || acc?.name || l.accountCode,
        debit: Number(l.debit) || 0,
        credit: Number(l.credit) || 0,
        description: l.description || data.description
      }
    })

    // Create journal entry automatically
    const journal = journalsStore.createManualJournal({
      journalDate: data.transactionDate,
      businessUnit: bu,
      branchId: 'BRANCH-MAIN',
      branchName: 'Kantor Pusat Semarang',
      description: `Voucher Kas/Bank [${txNumber}]: ${data.description}`,
      journalSource: 'MANUAL',
      sourceReference: txNumber,
      status: 'POSTED',
      lines: resolvedLines
    })

    const newTx: CashBankTransaction = {
      id: 'cb-' + Date.now(),
      transactionNumber: txNumber,
      date: data.transactionDate,
      type: data.voucherType,
      accountCode: data.primaryAccountCode,
      accountName: primaryName,
      counterpartAccountCode: counterpartCode,
      counterpartAccountName: counterpartName,
      amount,
      description: data.description,
      referenceNo: data.referenceNo || '-',
      recipientOrPayer: data.recipientOrPayer,
      businessUnit: bu,
      journalId: journal.id,
      journalNumber: journal.journalNumber,
      lines: resolvedLines,
      createdAt: new Date().toISOString(),
      createdBy: accountingStore.currentUserName
    }

    transactions.value.unshift(newTx)
    accountingStore.showToast(`Voucher ${txNumber} (${formatIDR(amount)}) berhasil dibukukan!`, 'success')
    return newTx
  }

  // Action: Buat Buku Kas / Rekening Bank Baru (Tambah ke COA & Database)
  async function createCashBankAccount(data: {
    category: 'CASH' | 'BANK'
    code: string
    name: string
    businessUnit: BusinessUnit
    openingBalance: number
    description?: string
  }) {
    const parentCode = data.category === 'CASH' ? '1101' : '1102'

    // Add account to COA in PostgreSQL database
    const acc = await ledgerStore.addAccount({
      code: data.code,
      name: data.name,
      type: 'ASSET',
      normalBalance: 'DEBIT',
      level: 4,
      businessUnit: data.businessUnit,
      description: data.description || (data.category === 'CASH' ? 'Buku Kas Operasional' : 'Rekening Bank Operasional'),
      parentCode
    })

    // Register in cash_bank_accounts table in database
    try {
      const cbaCode = (data.category === 'CASH' ? 'KAS-' : 'BANK-') + data.code.replace('.', '-')
      await accountingApi.createCashBankAccount({
        account_id: acc.id,
        code: cbaCode,
        name: data.name,
        type: data.category,
        bank_name: data.category === 'BANK' ? data.name : undefined,
        account_number: data.category === 'BANK' ? data.code : undefined
      })
      await fetchCashBankData()
    } catch (err: any) {
      console.warn('Note: Registered locally, backend CBA sync:', err.message)
    }

    // If opening balance > 0, create opening balance journal
    if (data.openingBalance > 0) {
      journalsStore.createManualJournal({
        journalDate: '2026-03-01',
        businessUnit: data.businessUnit,
        branchId: 'BRANCH-MAIN',
        branchName: 'Kantor Pusat Semarang',
        description: `Saldo Awal Pembukaan ${data.name} [${data.code}]`,
        journalSource: 'MANUAL',
        sourceReference: `OP-${data.code}`,
        status: 'POSTED',
        lines: [
          {
            accountCode: data.code,
            accountName: data.name,
            debit: data.openingBalance,
            credit: 0,
            description: `Saldo awal buku kas/bank ${data.code}`
          },
          {
            accountCode: '3-9999',
            accountName: 'Akun Penyeimbang Saldo Awal',
            debit: 0,
            credit: data.openingBalance,
            description: `Penyeimbang saldo awal pembukaan buku`
          }
        ]
      })
    }

    accountingStore.showToast(`Buku ${data.category === 'CASH' ? 'Kas' : 'Bank'} ${data.code} - ${data.name} berhasil disimpan ke database!`, 'success')
  }

  return {
    transactions,
    isLoading,
    registeredCashBankAccounts,
    cashAccounts,
    bankAccounts,
    allCashBankAccounts,
    fetchCashBankData,
    recordCashIn,
    recordCashOut,
    recordTransfer,
    recordMultiLineVoucher,
    createCashBankAccount
  }
})

