import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Account, LedgerEntry, TrialBalanceRow, AccountingPeriod, BusinessUnit } from '@/types/accounting'
import { useJournalsStore } from './journals'
import { useAccountingStore } from './accounting'
import { accountingApi } from '@/api/accounting'

export const useLedgerStore = defineStore('ledger', () => {
  const journalsStore = useJournalsStore()
  const accountingStore = useAccountingStore()

  // Official Dual-Brand Clinic COA (Kidinest Pediatric & Semarang Skin Center)
  const accounts = ref<Account[]>([
    // ==========================================
    // 1. AKTIVA / ASSETS (DEBIT)
    // ==========================================
    { id: '1-0000', code: '1-0000', name: 'AKTIVA / ASSETS', type: 'ASSET', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '1-1000', code: '1-1000', name: 'Aset Lancar', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000', currentBalance: 0 },
    { id: '1-1110', code: '1-1110', name: 'Kas Kecil (Petty Cash Operasional)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1111', code: '1-1111', name: 'Kas Kasir Kidinest', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1112', code: '1-1112', name: 'Kas Kasir Semarang Skin Center (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1113', code: '1-1113', name: 'Bank BCA Operasional (Giro & EDC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1114', code: '1-1114', name: 'Bank Mandiri Operasional & Payroll', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1120', code: '1-1120', name: 'Kliring Payment Gateway (QRIS/Tripay/Midtrans)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },

    { id: '1-1200', code: '1-1200', name: 'Piutang Pasien & Klaim Asuransi', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },
    { id: '1-1210', code: '1-1210', name: 'Piutang Klaim BPJS Kesehatan', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 0 },

    { id: '1-1300', code: '1-1300', name: 'Persediaan Farmasi & Medis', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000', currentBalance: 0 },
    { id: '1-1310', code: '1-1310', name: 'Persediaan Obat & Sirup Anak (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1300', currentBalance: 0 },
    { id: '1-1320', code: '1-1320', name: 'Persediaan Vaksin & Imunisasi Anak (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1300', currentBalance: 0 },
    { id: '1-1330', code: '1-1330', name: 'Persediaan Produk Skincare & Sunscreen (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1300', currentBalance: 0 },
    { id: '1-1340', code: '1-1340', name: 'Persediaan BMHP Estetika & Facial (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1300', currentBalance: 0 },

    { id: '1-2000', code: '1-2000', name: 'Aset Tetap & Peralatan Medis', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000', currentBalance: 0 },
    { id: '1-2100', code: '1-2100', name: 'Peralatan Medis & USG 4D (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-2000', currentBalance: 0 },
    { id: '1-2200', code: '1-2200', name: 'Mesin Laser Nd:YAG & Alkes Estetika (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-2000', currentBalance: 0 },
    { id: '1-2300', code: '1-2300', name: 'Akumulasi Penyusutan Aset Tetap', type: 'ASSET', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-2000', currentBalance: 0 },

    // ==========================================
    // 2. KEWAJIBAN / LIABILITIES (CREDIT)
    // ==========================================
    { id: '2-0000', code: '2-0000', name: 'KEWAJIBAN / LIABILITIES', type: 'LIABILITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '2-1000', code: '2-1000', name: 'Hutang Lancar', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2-0000', currentBalance: 0 },
    { id: '2-1110', code: '2-1110', name: 'Hutang Supplier Obat & Vaksin (Kidinest)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1115', code: '2-1115', name: 'Hutang Supplier Skincare & BMHP (SSC)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1120', code: '2-1120', name: 'Hutang Jasa Medis Dokter Anak (Sp.A)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1125', code: '2-1125', name: 'Hutang Jasa Medis Dokter Kulit & Terapis (SSC)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1130', code: '2-1130', name: 'Hutang Gaji & Karyawan', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1140', code: '2-1140', name: 'Uang Muka Pasien / Deposit Tindakan', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1500', code: '2-1500', name: 'Hutang Pajak (PPh 21, 23, Final 4(2), PPN)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0, taxCategory: 'PPH_21' },
    { id: '2-1510', code: '2-1510', name: 'Kredit Pajak PPh 23 (Jasa Limbah & IT)', type: 'LIABILITY', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0, taxCategory: 'PPH_23' },
    { id: '2-1520', code: '2-1520', name: 'Uang Muka / Angsuran PPh 25 Badan', type: 'LIABILITY', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0 },
    { id: '2-1530', code: '2-1530', name: 'Hutang PPh 29 (Kurang Bayar SPT Tahunan Badan)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 0 },

    // ==========================================
    // 3. EKUITAS / EQUITY (CREDIT)
    // ==========================================
    { id: '3-0000', code: '3-0000', name: 'MODAL / EQUITY', type: 'EQUITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '3-1100', code: '3-1100', name: 'Modal Disetor Pemilik Klinik', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 0 },
    { id: '3-1200', code: '3-1200', name: 'Laba Ditahan (Retained Earnings)', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 0 },
    { id: '3-1300', code: '3-1300', name: 'Laba Periode Berjalan', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 0 },
    { id: '3-9999', code: '3-9999', name: 'Akun Penyeimbang Saldo Awal (Equity Balancing)', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 0 },

    // ==========================================
    // 4. PENDAPATAN / REVENUE (CREDIT)
    // ==========================================
    { id: '4-0000', code: '4-0000', name: 'PENDAPATAN / REVENUE', type: 'REVENUE', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '4-1000', code: '4-1000', name: 'PENDAPATAN UNIT KIDINEST (ANAK)', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-0000', currentBalance: 0 },
    { id: '4-1100', code: '4-1100', name: 'Pendapatan Pendaftaran & Skrining Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 0 },
    { id: '4-1200', code: '4-1200', name: 'Pendapatan Jasa Dokter Spesialis Anak (Sp.A)', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 0 },
    { id: '4-1300', code: '4-1300', name: 'Pendapatan Vaksinasi & Imunisasi Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 0 },
    { id: '4-1400', code: '4-1400', name: 'Pendapatan Pemantauan Tumbuh Kembang Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 0 },
    { id: '4-1500', code: '4-1500', name: 'Pendapatan Penjualan Obat & Racikan Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 0 },

    { id: '4-2000', code: '4-2000', name: 'PENDAPATAN SEMARANG SKIN CENTER (SSC)', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'SSC', parentCode: '4-0000', currentBalance: 0 },
    { id: '4-2100', code: '4-2100', name: 'Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 0 },
    { id: '4-2200', code: '4-2200', name: 'Pendapatan Treatment Facial & Estetika Medis', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 0 },
    { id: '4-2300', code: '4-2300', name: 'Pendapatan Tindakan Laser, Cauter & Peeling', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 0 },
    { id: '4-2400', code: '4-2400', name: 'Pendapatan Penjualan Produk Skincare & Sunscreen', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 0 },

    // ==========================================
    // 5. HARGA POKOK PENJUALAN / COGS (DEBIT)
    // ==========================================
    { id: '5-0000', code: '5-0000', name: 'HARGA POKOK PENJUALAN (HPP)', type: 'COGS', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '5-1100', code: '5-1100', name: 'HPP Obat & Racikan Anak (Kidinest)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '5-0000', currentBalance: 0 },
    { id: '5-1200', code: '5-1200', name: 'HPP Vaksin & Imunisasi Anak (Kidinest)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '5-0000', currentBalance: 0 },
    { id: '5-1300', code: '5-1300', name: 'HPP Produk Skincare & Homecare (SSC)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '5-0000', currentBalance: 0 },
    { id: '5-1400', code: '5-1400', name: 'HPP Bahan BMHP & Treatment Estetika (SSC)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '5-0000', currentBalance: 0 },

    // ==========================================
    // 6. BEBAN OPERASIONAL / EXPENSES (DEBIT)
    // ==========================================
    { id: '6-0000', code: '6-0000', name: 'BEBAN OPERASIONAL / EXPENSES', type: 'EXPENSE', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 0 },
    { id: '6-1100', code: '6-1100', name: 'Beban Honor Jasa Dokter Spesialis Anak (Sp.A)', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '6-0000', currentBalance: 0 },
    { id: '6-1110', code: '6-1110', name: 'Beban Honor Jasa Dokter Kulit & Terapis Estetika', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '6-0000', currentBalance: 0 },
    { id: '6-1200', code: '6-1200', name: 'Beban Gaji Karyawan, Perawat & Staf', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 0 },
    { id: '6-1300', code: '6-1300', name: 'Beban Listrik, Air & Telepon', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 0 },
    { id: '6-1400', code: '6-1400', name: 'Beban Administrasi Bank & EDC', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 0 },
    { id: '6-1500', code: '6-1500', name: 'Beban Pajak Penghasilan Badan (PPh Badan Berjalan)', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 0 }
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
      if (accountingStore.fiscalMonth && journal.journalDate) {
        const parts = journal.journalDate.split('-')
        if (parts.length >= 2) {
          const m = parseInt(parts[1], 10)
          if (m !== accountingStore.fiscalMonth) continue
        }
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
          runningBalance: 0
        })
      }
    }

    list.sort((a, b) => a.date.localeCompare(b.date))
    return list
  })

  // Function to get account ledger with running balance
  function getAccountLedger(accountCode: string) {
    const account = accounts.value.find(a => a.code === accountCode)
    if (!account) return { account: null, entries: [], openingBalance: 0, endingBalance: 0 }

    const isDebitNormal = account.normalBalance === 'DEBIT'
    const openingBalance = account.currentBalance || 0
    let currentBalance = openingBalance

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
      openingBalance,
      endingBalance: currentBalance
    }
  }

  // Compute Trial Balance Rows dynamically
  const trialBalanceRows = computed<TrialBalanceRow[]>(() => {
    const nonHeaderAccounts = accounts.value.filter(a => {
      if (a.isHeader) return false
      if (accountingStore.businessUnit !== 'ALL' && a.businessUnit && a.businessUnit !== 'ALL' && a.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return true
    })
    const rows: TrialBalanceRow[] = []

    for (const acc of nonHeaderAccounts) {
      const glForAcc = generalLedgerEntries.value.filter(e => e.accountCode === acc.code)
      const periodDebit = glForAcc.reduce((sum, e) => sum + e.debit, 0)
      const periodCredit = glForAcc.reduce((sum, e) => sum + e.credit, 0)

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

  const isLoading = ref(false)

  // Local persistence for custom edits/deactivations across refreshes
  const LOCAL_STORAGE_KEY = 'clinic_coa_overrides'

  function loadLocalOverrides() {
    try {
      const raw = localStorage.getItem(LOCAL_STORAGE_KEY)
      if (!raw) return
      const overrides: Record<string, { name?: string; isActive?: boolean; description?: string }> = JSON.parse(raw)
      accounts.value.forEach(a => {
        const o = overrides[a.code] || overrides[a.id]
        if (o) {
          if (o.isActive !== undefined) a.isActive = o.isActive
          if (o.name !== undefined) a.name = o.name
          if (o.description !== undefined) a.description = o.description
        }
      })
    } catch {}
  }

  function saveLocalOverrides() {
    try {
      const overrides: Record<string, { name?: string; isActive?: boolean; description?: string }> = {}
      accounts.value.forEach(a => {
        if (a.isActive === false || a.name || a.description) {
          overrides[a.code] = {
            name: a.name,
            isActive: a.isActive,
            description: a.description
          }
        }
      })
      localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(overrides))
    } catch {}
  }

  // Load initial overrides on store creation
  loadLocalOverrides()

  // Fetch Accounts directly from PostgreSQL database via Go API
  async function fetchAccounts() {
    try {
      isLoading.value = true
      const backendAccounts = await accountingApi.getAccounts({ limit: 250 })
      if (backendAccounts && backendAccounts.length > 0) {
        // Map id to code for parent lookup
        const idToCode = new Map<string, string>()
        backendAccounts.forEach(b => idToCode.set(b.id, b.code))

        // Cache existing calculated balances
        const balanceMap = new Map<string, number>()
        accounts.value.forEach(a => {
          if (a.currentBalance !== undefined) {
            balanceMap.set(a.code, a.currentBalance)
          }
        })

        accounts.value = backendAccounts.map(b => {
          const parentCode = b.parent_id ? idToCode.get(b.parent_id) : undefined
          let bu: 'ALL' | 'KIDINEST' | 'SSC' = 'ALL'
          if (b.name.includes('Kidinest') || b.code.endsWith('.001')) bu = 'KIDINEST'
          else if (b.name.includes('SSC') || b.name.includes('Skin Center') || b.code.endsWith('.002')) bu = 'SSC'

          return {
            id: b.id,
            code: b.code,
            name: b.name,
            type: b.type,
            normalBalance: b.normal_balance,
            level: b.level,
            isHeader: b.is_header,
            isActive: b.is_active,
            businessUnit: bu,
            parentCode,
            description: b.description,
            currentBalance: balanceMap.get(b.code) || 0
          }
        })

        // Re-apply local overrides after fetching backend accounts
        loadLocalOverrides()
      }
    } catch (err: any) {
      console.error('Failed to fetch accounts from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load from database
  fetchAccounts()

  async function addAccount(newAcc: Partial<Account>) {
    try {
      isLoading.value = true
      let parentId: string | undefined = undefined
      if (newAcc.parentCode) {
        const parent = accounts.value.find(a => a.code === newAcc.parentCode)
        if (parent) parentId = parent.id
      }

      const createdBackend = await accountingApi.createAccount({
        code: newAcc.code || '',
        name: newAcc.name || '',
        type: newAcc.type || 'EXPENSE',
        normal_balance: newAcc.normalBalance || 'DEBIT',
        parent_id: parentId,
        level: newAcc.level || 2,
        is_header: newAcc.isHeader || false,
        description: newAcc.description
      })

      const created: Account = {
        id: createdBackend.id,
        code: createdBackend.code,
        name: createdBackend.name,
        type: createdBackend.type,
        normalBalance: createdBackend.normal_balance,
        level: createdBackend.level,
        isHeader: createdBackend.is_header,
        isActive: createdBackend.is_active,
        businessUnit: newAcc.businessUnit || 'ALL',
        description: createdBackend.description,
        parentCode: newAcc.parentCode,
        currentBalance: 0
      }
      accounts.value.push(created)
      accountingStore.showToast(`Akun ${created.code} - ${created.name} berhasil disimpan ke database!`, 'success')
      return created
    } catch (err: any) {
      accountingStore.showToast(`Gagal menyimpan akun ke database: ${err.message}`, 'error')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function updateAccount(id: string, payload: { name?: string; isActive?: boolean; description?: string }) {
    try {
      isLoading.value = true
      let updatedBackend: any = null
      try {
        updatedBackend = await accountingApi.updateAccount(id, {
          name: payload.name,
          is_active: payload.isActive,
          description: payload.description
        })
      } catch {
        // Fallback local if offline/mock
      }

      const idx = accounts.value.findIndex(a => a.id === id || a.code === id || (updatedBackend && a.code === updatedBackend.code))
      if (idx !== -1) {
        if (payload.name) accounts.value[idx].name = payload.name
        if (payload.isActive !== undefined) accounts.value[idx].isActive = payload.isActive
        if (payload.description !== undefined) accounts.value[idx].description = payload.description
        saveLocalOverrides()
      }

      const actionName = payload.isActive === false ? 'dinonaktifkan' : payload.isActive === true ? 'diaktifkan kembali' : 'diperbarui'
      accountingStore.showToast(`Akun berhasil ${actionName}!`, 'success')
      return accounts.value[idx]
    } catch (err: any) {
      accountingStore.showToast(`Gagal memperbarui akun: ${err.message}`, 'error')
      throw err
    } finally {
      isLoading.value = false
    }
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

  // Save and post official Opening Balance Journal to PostgreSQL database
  async function saveOpeningBalance(payload: {
    cutOffDate: string
    referenceNo: string
    businessUnit: BusinessUnit
    description: string
    lines: Array<{
      accountCode: string
      accountName: string
      debit: number
      credit: number
      description?: string
    }>
  }) {
    const validLines = payload.lines.filter(l => (l.debit > 0 || l.credit > 0))
    if (validLines.length === 0) {
      throw new Error('Harap masukkan minimal satu akun dengan saldo awal debit atau kredit!')
    }

    const totalDebit = validLines.reduce((s, l) => s + (l.debit || 0), 0)
    const totalCredit = validLines.reduce((s, l) => s + (l.credit || 0), 0)

    if (Math.abs(totalDebit - totalCredit) > 0.01) {
      throw new Error(`Total Debit (Rp ${totalDebit.toLocaleString('id-ID')}) tidak seimbang dengan Total Kredit (Rp ${totalCredit.toLocaleString('id-ID')}). Selisih: Rp ${Math.abs(totalDebit - totalCredit).toLocaleString('id-ID')}. Harap seimbangkan terlebih dahulu!`)
    }

    try {
      isLoading.value = true

      // 1. Create Opening Balance Journal in PostgreSQL database
      const journalRes = await accountingApi.createJournal({
        journal_type: 'GENERAL',
        transaction_date: payload.cutOffDate,
        reference_type: 'OPENING_BALANCE',
        reference_id: payload.referenceNo || `OP-${payload.cutOffDate.replace(/-/g, '')}`,
        description: payload.description || `Saldo Awal Pembukuan Akuntansi Klinik per ${payload.cutOffDate}`,
        auto_post: true,
        lines: validLines.map(l => ({
          account_code: l.accountCode,
          debit: l.debit || 0,
          credit: l.credit || 0,
          description: l.description || `Saldo Awal ${l.accountCode} - ${l.accountName}`
        }))
      })

      // 2. Update local account balances
      const balanceMap = new Map<string, number>()
      validLines.forEach(l => {
        balanceMap.set(l.accountCode, l.debit > 0 ? l.debit : l.credit)
      })

      accounts.value.forEach(a => {
        if (balanceMap.has(a.code)) {
          a.currentBalance = balanceMap.get(a.code)
        }
      })

      // 3. Refresh journals store
      await journalsStore.fetchJournals()

      accountingStore.showToast(`Saldo Awal Akuntansi (${journalRes.journal_number}) berhasil disimpan & diposting ke database!`, 'success')
      return journalRes
    } catch (err: any) {
      accountingStore.showToast(`Gagal menyimpan Saldo Awal: ${err.message}`, 'error')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    accounts,
    periods,
    isLoading,
    generalLedgerEntries,
    trialBalanceRows,
    trialBalanceTotals,
    fetchAccounts,
    getAccountLedger,
    addAccount,
    updateAccount,
    lockPeriod,
    saveOpeningBalance
  }
})
