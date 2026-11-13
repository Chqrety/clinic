import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CashBankTransaction, CashBankTransactionType, BusinessUnit, Account } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'

export const useCashBankStore = defineStore('cashbank', () => {
  const accountingStore = useAccountingStore()
  const ledgerStore = useLedgerStore()
  const journalsStore = useJournalsStore()

  // Transactions list
  const transactions = ref<CashBankTransaction[]>([
    {
      id: 'cb-1',
      transactionNumber: 'BKM-2026-03-001',
      date: '2026-03-01',
      type: 'CASH_IN',
      accountCode: '1101.003',
      accountName: 'Kas Kecil (Petty Cash Operasional)',
      counterpartAccountCode: '1102.001',
      counterpartAccountName: 'Bank BCA Operasional (Giro & EDC)',
      amount: 3000000,
      description: 'Pengisian awal dana kas kecil operasional klinik bulan Maret',
      referenceNo: 'CHQ-BCA-88912',
      recipientOrPayer: 'Bank BCA KCU Pemuda',
      businessUnit: 'ALL',
      journalId: 'jv-101',
      journalNumber: 'JV-2026-03-001',
      createdAt: '2026-03-01 09:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-2',
      transactionNumber: 'BKK-2026-03-001',
      date: '2026-03-02',
      type: 'CASH_OUT',
      accountCode: '1101.003',
      accountName: 'Kas Kecil (Petty Cash Operasional)',
      counterpartAccountCode: '5202',
      counterpartAccountName: 'BEBAN UMUM DAN ADMINISTRASI',
      amount: 450000,
      description: 'Pembelian 10 galon air mineral & konsumsi rapat koordinasi dokter poli',
      referenceNo: 'NOTA-SB-3312',
      recipientOrPayer: 'Toko Sumber Berkah',
      businessUnit: 'KIDINEST',
      journalId: 'jv-105',
      journalNumber: 'JV-2026-03-005',
      createdAt: '2026-03-02 11:30:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-3',
      transactionNumber: 'BKK-2026-03-002',
      date: '2026-03-03',
      type: 'CASH_OUT',
      accountCode: '1101.003',
      accountName: 'Kas Kecil (Petty Cash Operasional)',
      counterpartAccountCode: '5202',
      counterpartAccountName: 'BEBAN UMUM DAN ADMINISTRASI',
      amount: 275000,
      description: 'Pembelian materai 10.000 & map formulir rekam medis',
      referenceNo: 'KW-GRM-901',
      recipientOrPayer: 'Gramedia Pandanaran',
      businessUnit: 'ALL',
      journalId: 'jv-106',
      journalNumber: 'JV-2026-03-006',
      createdAt: '2026-03-03 14:15:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-4',
      transactionNumber: 'BBM-2026-03-001',
      date: '2026-03-03',
      type: 'CASH_IN',
      accountCode: '1102.001',
      accountName: 'Bank BCA Operasional (Giro & EDC)',
      counterpartAccountCode: '1202.001',
      counterpartAccountName: 'Piutang Asuransi Sinarmas',
      amount: 45000000,
      description: 'Penerimaan klaim reimbursement asuransi Sinarmas periode rawat jalan',
      referenceNo: 'TRF-SINAR-88391',
      recipientOrPayer: 'PT Asuransi Sinarmas',
      businessUnit: 'KIDINEST',
      journalId: 'jv-107',
      journalNumber: 'JV-2026-03-007',
      createdAt: '2026-03-03 15:40:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-5',
      transactionNumber: 'BBK-2026-03-001',
      date: '2026-03-04',
      type: 'CASH_OUT',
      accountCode: '1102.001',
      accountName: 'Bank BCA Operasional (Giro & EDC)',
      counterpartAccountCode: '1501',
      counterpartAccountName: 'Sewa Ruko & Fasilitas Dibayar di Muka',
      amount: 25000000,
      description: 'Pembayaran termin perpanjangan sewa gedung ruko klinik lantai 2',
      referenceNo: 'TRF-SEWA-2026',
      recipientOrPayer: 'Ibu Hj. Siti Nurhaliza (Pemilik Gedung)',
      businessUnit: 'ALL',
      journalId: 'jv-108',
      journalNumber: 'JV-2026-03-008',
      createdAt: '2026-03-04 10:20:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-6',
      transactionNumber: 'TRF-2026-03-001',
      date: '2026-03-04',
      type: 'TRANSFER',
      accountCode: '1101.001',
      accountName: 'Kas Kasir Kidinest',
      counterpartAccountCode: '1102.001',
      counterpartAccountName: 'Bank BCA Operasional (Giro & EDC)',
      amount: 10000000,
      description: 'Setoran tunai penerimaan kasir Kidinest shift pagi ke rekening Giro BCA',
      referenceNo: 'SETORAN-BCA-0019',
      recipientOrPayer: 'Kasir Utama Kidinest',
      businessUnit: 'KIDINEST',
      journalId: 'jv-109',
      journalNumber: 'JV-2026-03-009',
      createdAt: '2026-03-04 13:00:00',
      createdBy: 'Rian Pratama, S.Ak.'
    },
    {
      id: 'cb-7',
      transactionNumber: 'BBM-2026-03-002',
      date: '2026-03-05',
      type: 'CASH_IN',
      accountCode: '1102.002',
      accountName: 'Bank Mandiri Operasional & Payroll',
      counterpartAccountCode: '6101',
      counterpartAccountName: 'PENDAPATAN BUNGA BANK',
      amount: 1250000,
      description: 'Pendapatan jasa giro bunga bank rekening Mandiri Payroll',
      referenceNo: 'BUNGA-MANDIRI-MAR',
      recipientOrPayer: 'Bank Mandiri (Persero) Tbk',
      businessUnit: 'ALL',
      journalId: 'jv-110',
      journalNumber: 'JV-2026-03-010',
      createdAt: '2026-03-05 16:30:00',
      createdBy: 'Rian Pratama, S.Ak.'
    }
  ])

  // Get all cash accounts (parentCode === '1101' or startsWith '1101')
  const cashAccounts = computed(() => {
    return ledgerStore.accounts.filter(a => {
      if (a.isHeader) return false
      return a.parentCode === '1101' || a.code.startsWith('1101.') || a.code === '1101'
    })
  })

  // Get all bank accounts (parentCode === '1102' or startsWith '1102' or '1103')
  const bankAccounts = computed(() => {
    return ledgerStore.accounts.filter(a => {
      if (a.isHeader) return false
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

  // Action: Buat Buku Kas / Rekening Bank Baru (Tambah ke COA)
  function createCashBankAccount(data: {
    category: 'CASH' | 'BANK'
    code: string
    name: string
    businessUnit: BusinessUnit
    openingBalance: number
    description?: string
  }) {
    const parentCode = data.category === 'CASH' ? '1101' : '1102'

    // Add account to COA
    ledgerStore.addAccount({
      code: data.code,
      name: data.name,
      type: 'ASSET',
      normalBalance: 'DEBIT',
      level: 4,
      businessUnit: data.businessUnit,
      description: data.description || (data.category === 'CASH' ? 'Buku Kas Operasional' : 'Rekening Bank Operasional'),
      parentCode
    })

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
            accountCode: '1999',
            accountName: 'Temporary Opening - BCS',
            debit: 0,
            credit: data.openingBalance,
            description: `Penyeimbang saldo awal pembukaan buku`
          }
        ]
      })
    }

    accountingStore.showToast(`Buku ${data.category === 'CASH' ? 'Kas' : 'Bank'} ${data.code} - ${data.name} berhasil dibuat!`, 'success')
  }

  return {
    transactions,
    cashAccounts,
    bankAccounts,
    allCashBankAccounts,
    recordCashIn,
    recordCashOut,
    recordTransfer,
    createCashBankAccount
  }
})
