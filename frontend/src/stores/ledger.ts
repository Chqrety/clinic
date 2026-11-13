import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Account, LedgerEntry, TrialBalanceRow, AccountingPeriod } from '@/types/accounting'
import { useJournalsStore } from './journals'
import { useAccountingStore } from './accounting'

export const useLedgerStore = defineStore('ledger', () => {
  const journalsStore = useJournalsStore()
  const accountingStore = useAccountingStore()

  // Standard Indonesian Chart of Accounts for Clinic ERP
  const accounts = ref<Account[]>([
    // 1-0000 ASET
    { id: '1-0000', code: '1-0000', name: 'ASET', type: 'ASSET', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '1-1000', code: '1-1000', name: 'Aset Lancar', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000' },
    { id: '1-1100', code: '1-1100', name: 'Kas Operasional Kidinest', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1000', currentBalance: 45200000 },
    { id: '1-1101', code: '1-1101', name: 'Kas Operasional SSC', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1000', currentBalance: 38400000 },
    { id: '1-1200', code: '1-1200', name: 'Bank Mandiri Operasional', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 185600000 },
    { id: '1-1201', code: '1-1201', name: 'Bank BCA Operasional (EDC & QRIS)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 242000000 },
    { id: '1-1300', code: '1-1300', name: 'Piutang Pasien & Asuransi', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 64200000 },
    { id: '1-1400', code: '1-1400', name: 'Persediaan Farmasi & Obat', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 98500000 },
    { id: '1-1180', code: '1-1180', name: 'PPN Masukan (Dapat Dikreditkan)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 8400000 },
    { id: '1-2000', code: '1-2000', name: 'Aset Tetap', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000' },
    { id: '1-2010', code: '1-2010', name: 'Peralatan Medis & Laser', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-2000', currentBalance: 750000000 },
    { id: '1-2100', code: '1-2100', name: 'Akumulasi Penyusutan Alat Medis', type: 'ASSET', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-2000', currentBalance: 125000000 },

    // 2-0000 KEWAJIBAN
    { id: '2-0000', code: '2-0000', name: 'KEWAJIBAN', type: 'LIABILITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '2-1000', code: '2-1000', name: 'Kewajiban Jangka Pendek', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2-0000' },
    { id: '2-1100', code: '2-1100', name: 'Hutang Usaha Distributor Farmasi', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 58200000 },
    { id: '2-1200', code: '2-1200', name: 'Hutang Jasa Medis Dokter', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 42100000 },
    { id: '2-1300', code: '2-1300', name: 'Hutang PPh 21 Dokter & Staf', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 7850000 },
    { id: '2-1310', code: '2-1310', name: 'Hutang PPh 23 / Final', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 2400000 },
    { id: '2-1320', code: '2-1320', name: 'Hutang PPN Keluaran', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 12500000 },

    // 3-0000 EKUITAS
    { id: '3-0000', code: '3-0000', name: 'EKUITAS', type: 'EQUITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '3-1000', code: '3-1000', name: 'Modal Disetor Pendiri', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 800000000 },
    { id: '3-2000', code: '3-2000', name: 'Laba Ditahan', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 215450000 },
    { id: '3-3000', code: '3-3000', name: 'Laba Periode Berjalan', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 167800000 },

    // 4-0000 PENDAPATAN
    { id: '4-0000', code: '4-0000', name: 'PENDAPATAN', type: 'REVENUE', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '4-1100', code: '4-1100', name: 'Pendapatan Jasa Medis Kidinest', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-0000', currentBalance: 165000000 },
    { id: '4-1101', code: '4-1101', name: 'Pendapatan Jasa Medis SSC', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-0000', currentBalance: 245000000 },
    { id: '4-1200', code: '4-1200', name: 'Pendapatan Farmasi & Obat', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '4-0000', currentBalance: 82000000 },
    { id: '4-1300', code: '4-1300', name: 'Pendapatan Skincare & Produk Slimming', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-0000', currentBalance: 48000000 },

    // 5-0000 BEBAN POKOK PENDAPATAN (HPP)
    { id: '5-0000', code: '5-0000', name: 'BEBAN POKOK PENDAPATAN (HPP)', type: 'COGS', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '5-1100', code: '5-1100', name: 'HPP Farmasi & Obat-obatan', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '5-0000', currentBalance: 48500000 },
    { id: '5-1200', code: '5-1200', name: 'Beban Jasa Medis Dokter', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '5-0000', currentBalance: 145000000 },

    // 6-0000 BEBAN OPERASIONAL
    { id: '6-0000', code: '6-0000', name: 'BEBAN OPERASIONAL', type: 'EXPENSE', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '6-1100', code: '6-1100', name: 'Gaji & Upah Karyawan', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 68000000 },
    { id: '6-1200', code: '6-1200', name: 'Beban Sewa Gedung Klinik', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 25000000 },
    { id: '6-1300', code: '6-1300', name: 'Beban Listrik, Air & Internet', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 9800000 },
    { id: '6-1400', code: '6-1400', name: 'Beban Penyusutan Aset Tetap', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 10416667 },
    { id: '6-1500', code: '6-1500', name: 'Beban Perlengkapan & BHP Medis', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 14800000 },
    { id: '6-1900', code: '6-1900', name: 'Beban Administrasi Bank', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 450000 }
  ])

  // Accounting Periods
  const periods = ref<AccountingPeriod[]>([
    {
      id: 'p-2026-01',
      year: 2026,
      month: 1,
      periodName: 'Januari 2026',
      startDate: '2026-01-01',
      endDate: '2026-01-31',
      status: 'LOCKED',
      lockedAt: '2026-02-05 18:00:00',
      lockedBy: 'budi.santoso@kidinest-ssc.id',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    },
    {
      id: 'p-2026-02',
      year: 2026,
      month: 2,
      periodName: 'Februari 2026',
      startDate: '2026-02-01',
      endDate: '2026-02-28',
      status: 'LOCKED',
      lockedAt: '2026-03-05 17:30:00',
      lockedBy: 'budi.santoso@kidinest-ssc.id',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    },
    {
      id: 'p-2026-03',
      year: 2026,
      month: 3,
      periodName: 'Maret 2026 (Aktif)',
      startDate: '2026-03-01',
      endDate: '2026-03-31',
      status: 'OPEN',
      unpostedJournalCount: 3,
      pendingApprovalCount: 2
    },
    {
      id: 'p-2026-04',
      year: 2026,
      month: 4,
      periodName: 'April 2026',
      startDate: '2026-04-01',
      endDate: '2026-04-30',
      status: 'OPEN',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    }
  ])

  // Compute General Ledger Entries directly from POSTED Journals
  const generalLedgerEntries = computed<LedgerEntry[]>(() => {
    const list: LedgerEntry[] = []
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    for (const journal of postedJournals) {
      if (accountingStore.businessUnit !== 'ALL' && journal.businessUnit !== 'ALL' && journal.businessUnit !== accountingStore.businessUnit) {
        continue
      }

      for (const line of journal.lines) {
        list.push({
          id: `gle-${journal.id}-${line.id}`,
          date: journal.journalDate,
          journalNumber: journal.journalNumber,
          journalId: journal.id,
          source: journal.journalSource,
          description: line.description || journal.description,
          accountCode: line.accountCode,
          accountName: line.accountName,
          businessUnit: journal.businessUnit,
          branchName: journal.branchName,
          debit: line.debit,
          credit: line.credit,
          runningBalance: 0 // Will be computed per account
        })
      }
    }

    // Sort by date ascending
    list.sort((a, b) => a.date.localeCompare(b.date))
    return list
  })

  // Function to get account ledger with running balance
  function getAccountLedger(accountCode: string) {
    const account = accounts.value.find(a => a.code === accountCode)
    if (!account) return { account: null, entries: [], openingBalance: 0, endingBalance: 0 }

    const isDebitNormal = account.normalBalance === 'DEBIT'
    let currentBalance = 0

    const entries = generalLedgerEntries.value
      .filter(e => e.accountCode === accountCode)
      .map(entry => {
        if (isDebitNormal) {
          currentBalance += (entry.debit - entry.credit)
        } else {
          currentBalance += (entry.credit - entry.debit)
        }
        return {
          ...entry,
          runningBalance: currentBalance
        }
      })

    return {
      account,
      entries,
      openingBalance: 0,
      endingBalance: currentBalance
    }
  }

  // Compute Trial Balance Rows dynamically
  const trialBalanceRows = computed<TrialBalanceRow[]>(() => {
    const nonHeaderAccounts = accounts.value.filter(a => !a.isHeader)
    const rows: TrialBalanceRow[] = []

    for (const acc of nonHeaderAccounts) {
      // Find movements in posted journals
      const glForAcc = generalLedgerEntries.value.filter(e => e.accountCode === acc.code)
      const periodDebit = glForAcc.reduce((sum, e) => sum + e.debit, 0)
      const periodCredit = glForAcc.reduce((sum, e) => sum + e.credit, 0)

      // Initial opening balances for realism
      const openingDebit = acc.normalBalance === 'DEBIT' ? (acc.currentBalance || 0) : 0
      const openingCredit = acc.normalBalance === 'CREDIT' ? (acc.currentBalance || 0) : 0

      let endingDebit = 0
      let endingCredit = 0

      if (acc.normalBalance === 'DEBIT') {
        const net = openingDebit + periodDebit - periodCredit
        if (net >= 0) endingDebit = net
        else endingCredit = Math.abs(net)
      } else {
        const net = openingCredit + periodCredit - periodDebit
        if (net >= 0) endingCredit = net
        else endingDebit = Math.abs(net)
      }

      rows.push({
        accountCode: acc.code,
        accountName: acc.name,
        accountType: acc.type,
        openingDebit,
        openingCredit,
        periodDebit,
        periodCredit,
        endingDebit,
        endingCredit
      })
    }

    return rows
  })

  // Trial Balance totals
  const trialBalanceTotals = computed(() => {
    return trialBalanceRows.value.reduce(
      (acc, row) => ({
        openingDebit: acc.openingDebit + row.openingDebit,
        openingCredit: acc.openingCredit + row.openingCredit,
        periodDebit: acc.periodDebit + row.periodDebit,
        periodCredit: acc.periodCredit + row.periodCredit,
        endingDebit: acc.endingDebit + row.endingDebit,
        endingCredit: acc.endingCredit + row.endingCredit
      }),
      {
        openingDebit: 0,
        openingCredit: 0,
        periodDebit: 0,
        periodCredit: 0,
        endingDebit: 0,
        endingCredit: 0
      }
    )
  })

  function addAccount(newAcc: Partial<Account>) {
    const created: Account = {
      id: 'acc-' + Date.now(),
      code: newAcc.code || '',
      name: newAcc.name || '',
      type: newAcc.type || 'EXPENSE',
      normalBalance: newAcc.normalBalance || 'DEBIT',
      level: newAcc.level || 2,
      isHeader: false,
      isActive: true,
      businessUnit: newAcc.businessUnit || 'ALL',
      description: newAcc.description,
      parentCode: newAcc.parentCode
    }
    accounts.value.push(created)
    accountingStore.showToast(`Akun ${created.code} - ${created.name} berhasil ditambahkan`, 'success')
  }

  function lockPeriod(periodId: string) {
    const p = periods.value.find(item => item.id === periodId)
    if (p) {
      p.status = 'LOCKED'
      p.lockedAt = new Date().toISOString()
      p.lockedBy = accountingStore.currentUserEmail
      accountingStore.showToast(`Periode ${p.periodName} berhasil ditutup & dikunci!`, 'success')
    }
  }

  return {
    accounts,
    periods,
    generalLedgerEntries,
    trialBalanceRows,
    trialBalanceTotals,
    getAccountLedger,
    addAccount,
    lockPeriod
  }
})
