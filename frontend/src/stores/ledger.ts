import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Account, LedgerEntry, TrialBalanceRow, AccountingPeriod } from '@/types/accounting'
import { useJournalsStore } from './journals'
import { useAccountingStore } from './accounting'

export const useLedgerStore = defineStore('ledger', () => {
  const journalsStore = useJournalsStore()
  const accountingStore = useAccountingStore()

  // Standard Indonesian Chart of Accounts for Clinic ERP (Blueprint Kidinest & SSC)
  const accounts = ref<Account[]>([
    // 1-0000 ASET / AKTIVA
    { id: '1-0000', code: '1-0000', name: 'ASET (AKTIVA)', type: 'ASSET', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '1-1000', code: '1-1000', name: 'Aset Lancar', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000' },
    { id: '1-1100', code: '1-1100', name: 'Kas & Setara Kas', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-1000' },
    { id: '1-1111', code: '1-1111', name: 'Kas Kasir Kidinest (Poli Anak)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1100', currentBalance: 45200000 },
    { id: '1-1112', code: '1-1112', name: 'Kas Kasir Semarang Skin Center (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1100', currentBalance: 38400000 },
    { id: '1-1113', code: '1-1113', name: 'Kas Kecil Operasional (Petty Cash)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1100', currentBalance: 5000000 },
    { id: '1-1121', code: '1-1121', name: 'Bank BCA Operasional (EDC & Giro)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1100', currentBalance: 242000000 },
    { id: '1-1122', code: '1-1122', name: 'Bank Mandiri Operasional', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1100', currentBalance: 185600000 },
    { id: '1-1130', code: '1-1130', name: 'Kliring Payment Gateway (QRIS & EDC Transit)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1100', currentBalance: 14500000 },

    // Piutang
    { id: '1-1200', code: '1-1200', name: 'Piutang Usaha', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-1000' },
    { id: '1-1201', code: '1-1201', name: 'Piutang Pasien Perawatan', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1200', currentBalance: 21400000 },
    { id: '1-1202', code: '1-1202', name: 'Piutang Klaim Asuransi & Korporat', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1200', currentBalance: 42800000 },

    // Persediaan (FIFO)
    { id: '1-1300', code: '1-1300', name: 'Persediaan Medis & Non-Medis (FIFO)', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-1000' },
    { id: '1-1310', code: '1-1310', name: 'Persediaan Obat & Sirup Anak (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1300', currentBalance: 38200000 },
    { id: '1-1311', code: '1-1311', name: 'Persediaan Vaksin & Imunisasi Anak (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-1300', currentBalance: 24500000 },
    { id: '1-1320', code: '1-1320', name: 'Persediaan Bahan Medis Habis Pakai (BMHP)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1300', currentBalance: 12400000 },
    { id: '1-1330', code: '1-1330', name: 'Persediaan Produk Skincare & Kosmetika SSC', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-1300', currentBalance: 23400000 },
    { id: '1-1180', code: '1-1180', name: 'PPN Masukan (Dapat Dikreditkan)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-1000', currentBalance: 8400000 },

    // Aset Tetap
    { id: '1-2000', code: '1-2000', name: 'Aset Tetap & Peralatan Medis', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1-0000' },
    { id: '1-2100', code: '1-2100', name: 'Peralatan Medis Poli Anak & USG 4D (Kidinest)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1-2000', currentBalance: 250000000 },
    { id: '1-2200', code: '1-2200', name: 'Peralatan Estetika Medis & Laser Nd:YAG (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1-2000', currentBalance: 500000000 },
    { id: '1-2900', code: '1-2900', name: 'Akumulasi Penyusutan Aset Tetap', type: 'ASSET', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1-2000', currentBalance: 125000000 },

    // 2-0000 KEWAJIBAN / LIABILITAS
    { id: '2-0000', code: '2-0000', name: 'KEWAJIBAN (LIABILITAS)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '2-1000', code: '2-1000', name: 'Kewajiban Jangka Pendek', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2-0000' },
    { id: '2-1110', code: '2-1110', name: 'Hutang Supplier Farmasi & PBF', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 38200000 },
    { id: '2-1120', code: '2-1120', name: 'Hutang Supplier Skincare & Kosmetika SSC', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2-1000', currentBalance: 20000000 },
    { id: '2-1210', code: '2-1210', name: 'Hutang Bagi Hasil Dokter Spesialis Anak (Kidinest)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '2-1000', currentBalance: 18400000 },
    { id: '2-1220', code: '2-1220', name: 'Hutang Bagi Hasil Dokter Estetika & Kulit (SSC)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2-1000', currentBalance: 23700000 },
    { id: '2-1501', code: '2-1501', name: 'Hutang PPh 21 Dokter Tenaga Medis Ahli (PMK 168)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 7850000 },
    { id: '2-1502', code: '2-1502', name: 'Hutang PPh 23 (Limbah B3 & Pemeliharaan)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 2400000 },
    { id: '2-1504', code: '2-1504', name: 'Hutang PPN Keluaran', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2-1000', currentBalance: 12500000 },

    // 3-0000 EKUITAS (MODAL)
    { id: '3-0000', code: '3-0000', name: 'EKUITAS (MODAL)', type: 'EQUITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '3-1100', code: '3-1100', name: 'Modal Disetor Pemilik Klinik', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 800000000 },
    { id: '3-1200', code: '3-1200', name: 'Laba Ditahan (Retained Earnings)', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 215450000 },
    { id: '3-1300', code: '3-1300', name: 'Laba Bersih Periode Berjalan', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3-0000', currentBalance: 167800000 },

    // 4-0000 PENDAPATAN
    { id: '4-0000', code: '4-0000', name: 'PENDAPATAN (REVENUE)', type: 'REVENUE', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '4-1000', code: '4-1000', name: 'PENDAPATAN UNIT KIDINEST (ANAK)', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-0000' },
    { id: '4-1110', code: '4-1110', name: 'Pendapatan Jasa Konsultasi Spesialis Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 115000000 },
    { id: '4-1120', code: '4-1120', name: 'Pendapatan Tindakan Vaksinasi & Imunisasi', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 32000000 },
    { id: '4-1130', code: '4-1130', name: 'Pendapatan Terapi & Tumbuh Kembang Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 18000000 },
    { id: '4-1140', code: '4-1140', name: 'Pendapatan Penjualan Obat Resep Farmasi Anak', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '4-1000', currentBalance: 62000000 },

    { id: '4-2000', code: '4-2000', name: 'PENDAPATAN UNIT SEMARANG SKIN CENTER (SSC)', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'SSC', parentCode: '4-0000' },
    { id: '4-2110', code: '4-2110', name: 'Pendapatan Konsultasi Dokter Estetika & Kulit (Sp.DVE)', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 65000000 },
    { id: '4-2120', code: '4-2120', name: 'Pendapatan Tindakan Medis Laser & Cauter SSC', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 145000000 },
    { id: '4-2130', code: '4-2130', name: 'Pendapatan Tindakan Facial & Treatment Terapis', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 35000000 },
    { id: '4-2140', code: '4-2140', name: 'Pendapatan Penjualan Produk Skincare & Kosmetika', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '4-2000', currentBalance: 68000000 },

    // 5-0000 BEBAN POKOK PENDAPATAN (HPP & BEBAN MEDIS)
    { id: '5-0000', code: '5-0000', name: 'BEBAN POKOK PENDAPATAN (HPP)', type: 'COGS', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '5-1101', code: '5-1101', name: 'HPP Obat Resep & Paten Farmasi Kidinest', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '5-0000', currentBalance: 24000000 },
    { id: '5-1102', code: '5-1102', name: 'HPP Vaksin & Serum Imunisasi Kidinest', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '5-0000', currentBalance: 8000000 },
    { id: '5-1103', code: '5-1103', name: 'HPP Produk Skincare & Kosmetika Terjual SSC', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '5-0000', currentBalance: 16500000 },
    { id: '5-1105', code: '5-1105', name: 'Beban Pemakaian Bahan Medis Habis Pakai (BMHP)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '5-0000', currentBalance: 9800000 },
    { id: '5-1106', code: '5-1106', name: 'Beban Bagi Hasil Dokter Spesialis Anak (Kidinest)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '5-0000', currentBalance: 58000000 },
    { id: '5-1107', code: '5-1107', name: 'Beban Bagi Hasil Dokter Estetika & Kulit (SSC)', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '5-0000', currentBalance: 78500000 },
    { id: '5-1108', code: '5-1108', name: 'Beban Komisi Tindakan Beautician / Terapis SSC', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '5-0000', currentBalance: 8500000 },

    // 6-0000 BEBAN OPERASIONAL & UMUM
    { id: '6-0000', code: '6-0000', name: 'BEBAN OPERASIONAL & UMUM', type: 'EXPENSE', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL' },
    { id: '6-1101', code: '6-1101', name: 'Beban Gaji Staf, Perawat, Apoteker & Kasir', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 68000000 },
    { id: '6-1102', code: '6-1102', name: 'Beban Listrik PLN, Air PDAM & Internet', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 9800000 },
    { id: '6-1103', code: '6-1103', name: 'Beban Pengolahan Limbah Medis B3', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 4200000 },
    { id: '6-1104', code: '6-1104', name: 'Beban Kalibrasi & Pemeliharaan Alat Medis/Laser', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 5600000 },
    { id: '6-1105', code: '6-1105', name: 'Beban Sewa Gedung & Ruko Klinik', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 25000000 },
    { id: '6-1107', code: '6-1107', name: 'Beban Penyusutan Aktiva Tetap Klinik', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 10416667 },
    { id: '6-1201', code: '6-1201', name: 'Beban Administrasi Bank & MDR EDC Merchant', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '6-0000', currentBalance: 450000 }
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
