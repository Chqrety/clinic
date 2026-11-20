import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  OperationalTransaction,
  SourceModule,
  AccountingMappingRule,
  BusinessUnit
} from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useJournalsStore } from './journals'

export const useTransactionsStore = defineStore('transactions', () => {
  const accountingStore = useAccountingStore()
  const journalsStore = useJournalsStore()

  // Pre-configured Accounting Mapping Rules
  const mappingRules = ref<AccountingMappingRule[]>([
    {
      id: 'rule-1',
      eventKey: 'CASHIER_CLOSING_KIDINEST',
      eventName: 'Penutupan Kasir Harian Kidinest',
      businessUnit: 'KIDINEST',
      sourceModule: 'GRO_CASHIER',
      debitAccountCode: '1-1100',
      debitAccountName: 'Kas Operasional Kidinest',
      creditAccountCode: '4-1100',
      creditAccountName: 'Pendapatan Jasa Medis Kidinest',
      description: 'Mencatat penerimaan kasir shift harian klinik anak',
      isActive: true
    },
    {
      id: 'rule-2',
      eventKey: 'CASHIER_CLOSING_SSC',
      eventName: 'Penutupan Kasir Harian SSC',
      businessUnit: 'SSC',
      sourceModule: 'GRO_CASHIER',
      debitAccountCode: '1-1101',
      debitAccountName: 'Kas Operasional SSC',
      creditAccountCode: '4-1101',
      creditAccountName: 'Pendapatan Jasa Medis SSC',
      description: 'Mencatat penerimaan kasir shift harian klinik estetika',
      isActive: true
    },
    {
      id: 'rule-3',
      eventKey: 'PHARMACY_SALE',
      eventName: 'Penjualan Obat Resep Farmasi',
      businessUnit: 'ALL',
      sourceModule: 'PHARMACY',
      debitAccountCode: '1-1100',
      debitAccountName: 'Kas Operasional Kidinest',
      creditAccountCode: '4-1200',
      creditAccountName: 'Pendapatan Farmasi & Obat',
      cogsAccountCode: '5-1100',
      inventoryAccountCode: '1-1400',
      description: 'Pengakuan omset obat beserta beban pokok (HPP)',
      isActive: true
    },
    {
      id: 'rule-4',
      eventKey: 'DOCTOR_FEE_CALC',
      eventName: 'Kompensasi Jasa Medis Dokter + PPh 21',
      businessUnit: 'ALL',
      sourceModule: 'DOCTOR_SHARE',
      debitAccountCode: '5-1200',
      debitAccountName: 'Beban Jasa Medis Dokter',
      creditAccountCode: '2-1200',
      creditAccountName: 'Hutang Jasa Medis Dokter',
      taxAccountCode: '2-1300',
      description: 'Kewajiban honor dokter dikurangi potongan PPh 21 dokter',
      isActive: true
    },
    {
      id: 'rule-5',
      eventKey: 'PURCHASE_PHARMACY_INVOICE',
      eventName: 'Faktur Pembelian Obat dari Distributor',
      businessUnit: 'ALL',
      sourceModule: 'PURCHASING',
      debitAccountCode: '1-1400',
      debitAccountName: 'Persediaan Farmasi & Obat',
      creditAccountCode: '2-1100',
      creditAccountName: 'Hutang Usaha Distributor Farmasi',
      taxAccountCode: '1-1180',
      description: 'Penerimaan stok obat berjangka tempo dengan PPN Masukan',
      isActive: true
    },
    {
      id: 'rule-6',
      eventKey: 'FIXED_ASSET_DEPRECIATION',
      eventName: 'Penyusutan Bulanan Peralatan Medis',
      businessUnit: 'ALL',
      sourceModule: 'FIXED_ASSET',
      debitAccountCode: '6-1400',
      debitAccountName: 'Beban Penyusutan Aset Tetap',
      creditAccountCode: '1-2100',
      creditAccountName: 'Akumulasi Penyusutan Alat Medis',
      description: 'Beban amortisasi bulanan peralatan laser dan USG',
      isActive: true
    },
    {
      id: 'rule-7',
      eventKey: 'PAYROLL_SALARY',
      eventName: 'Penggajian Karyawan & Tenaga Medis',
      businessUnit: 'ALL',
      sourceModule: 'PAYROLL',
      debitAccountCode: '6-1100',
      debitAccountName: 'Gaji & Upah Karyawan',
      creditAccountCode: '1-1200',
      creditAccountName: 'Bank Mandiri Operasional',
      taxAccountCode: '2-1300',
      description: 'Pembayaran payroll bulanan staf klinik',
      isActive: true
    }
  ])

  // Operational Transactions (Empty in production until created by cashier/billing)
  const transactions = ref<OperationalTransaction[]>([])

  // Filtered list based on global unit
  const filteredTransactions = computed(() => {
    return transactions.value.filter(t => {
      if (accountingStore.businessUnit !== 'ALL' && t.businessUnit !== 'ALL' && t.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return true
    })
  })

  // Action: 1-click Generate Draft Journal from Operational Transaction
  function generateDraftJournal(txId: string) {
    const tx = transactions.value.find(t => t.id === txId)
    if (!tx) {
      accountingStore.showToast('Transaksi tidak ditemukan', 'error')
      return
    }

    if (tx.journalStatus !== 'UNPROCESSED') {
      accountingStore.showToast(`Transaksi ${tx.transactionNumber} sudah pernah dibuatkan jurnal draft (${tx.generatedJournalId})`, 'info')
      return
    }

    // Determine account line items based on transaction details & mapping
    const lines: any[] = []
    const txDate = tx.transactionDate.split(' ')[0]

    if (tx.sourceModule === 'GRO_CASHIER') {
      const isKidinest = tx.businessUnit === 'KIDINEST'
      const cashAcc = isKidinest
        ? { code: '1-1100', name: 'Kas Operasional Kidinest' }
        : { code: '1-1101', name: 'Kas Operasional SSC' }
      const revAcc = isKidinest
        ? { code: '4-1100', name: 'Pendapatan Jasa Medis Kidinest' }
        : { code: '4-1101', name: 'Pendapatan Jasa Medis SSC' }

      lines.push({
        accountId: `acc-${cashAcc.code}`,
        accountCode: cashAcc.code,
        accountName: cashAcc.name,
        debit: tx.amount,
        credit: 0,
        description: `Penerimaan Kasir - ${tx.referenceNumber}`
      })
      lines.push({
        accountId: `acc-${revAcc.code}`,
        accountCode: revAcc.code,
        accountName: revAcc.name,
        debit: 0,
        credit: tx.amount,
        description: `Pendapatan Pelayanan Pasien Klinik`
      })
    } else if (tx.sourceModule === 'DOCTOR_SHARE') {
      const feeBruto = tx.details.feeBruto || tx.amount
      const pph21 = tx.details.pph21Potongan || Math.round(feeBruto * 0.025)
      const feeNet = feeBruto - pph21

      lines.push({
        accountId: 'acc-51200',
        accountCode: '5-1200',
        accountName: 'Beban Jasa Medis Dokter',
        debit: feeBruto,
        credit: 0,
        description: `Jasa Medis ${tx.details.dokter}`
      })
      lines.push({
        accountId: 'acc-21300',
        accountCode: '2-1300',
        accountName: 'Hutang PPh 21 Dokter',
        debit: 0,
        credit: pph21,
        description: `Potongan PPh 21 Pasal 17 PMK 168/2023`
      })
      lines.push({
        accountId: 'acc-21200',
        accountCode: '2-1200',
        accountName: 'Hutang Jasa Medis Dokter',
        debit: 0,
        credit: feeNet,
        description: `Hutang Bersih ke Dokter`
      })
    } else if (tx.sourceModule === 'PURCHASING') {
      const dpp = tx.details.dpp || Math.round(tx.amount / 1.11)
      const ppn = tx.amount - dpp

      lines.push({
        accountId: 'acc-11400',
        accountCode: '1-1400',
        accountName: 'Persediaan Farmasi & Obat',
        debit: dpp,
        credit: 0,
        description: `Pengadaan Obat DPP - ${tx.details.supplier}`
      })
      lines.push({
        accountId: 'acc-11180',
        accountCode: '1-1180',
        accountName: 'PPN Masukan (Dapat Dikreditkan)',
        debit: ppn,
        credit: 0,
        description: `PPN Masukan 11% Faktur ${tx.details.fakturPajak}`
      })
      lines.push({
        accountId: 'acc-21100',
        accountCode: '2-1100',
        accountName: 'Hutang Usaha Distributor Farmasi',
        debit: 0,
        credit: tx.amount,
        description: `Hutang Dagang ${tx.details.supplier}`
      })
    } else {
      // Generic balanced 2-line entry
      const debitCode = tx.suggestedDebitAccount.split(' ')[0]
      const debitName = tx.suggestedDebitAccount.replace(`${debitCode} (`, '').replace(')', '')
      const creditCode = tx.suggestedCreditAccount.split(' ')[0]
      const creditName = tx.suggestedCreditAccount.replace(`${creditCode} (`, '').replace(')', '')

      lines.push({
        accountId: `acc-${debitCode}`,
        accountCode: debitCode,
        accountName: debitName,
        debit: tx.amount,
        credit: 0,
        description: tx.description
      })
      lines.push({
        accountId: `acc-${creditCode}`,
        accountCode: creditCode,
        accountName: creditName,
        debit: 0,
        credit: tx.amount,
        description: tx.description
      })
    }

    const createdJournal = journalsStore.createManualJournal({
      journalDate: txDate,
      businessUnit: tx.businessUnit,
      branchId: 'BRANCH-MAIN',
      branchName: tx.branchName,
      description: `[Auto-Generated dari ${tx.sourceModule}] ${tx.description}`,
      lines
    })

    // Update operational transaction
    tx.journalStatus = 'DRAFT_GENERATED'
    tx.generatedJournalId = createdJournal.journalNumber
    accountingStore.showToast(`Draft Jurnal ${createdJournal.journalNumber} berhasil dibuat dari transaksi ${tx.transactionNumber}!`, 'success')
  }

  return {
    transactions,
    mappingRules,
    filteredTransactions,
    generateDraftJournal
  }
})
