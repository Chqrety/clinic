import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAccountingStore } from './accounting'
import { useLedgerStore } from './ledger'
import { useJournalsStore } from './journals'
import { accountingApi } from '@/api/accounting'

export interface ProfitLossItem {
  code: string
  name: string
  kidinestAmount: number
  sscAmount: number
  totalAmount: number
  percentage: number
}

export interface AgingBucket {
  id: string
  entityName: string
  invoiceNumber: string
  invoiceDate: string
  dueDate: string
  current: number
  days1to30: number
  days31to60: number
  days61to90: number
  over90: number
  total: number
  status: 'CURRENT' | 'OVERDUE'
}

export const useReportsStore = defineStore('reports', () => {
  const accountingStore = useAccountingStore()
  const ledgerStore = useLedgerStore()
  const journalsStore = useJournalsStore()

  // Dynamically compute Revenue items from Buku Besar (General Ledger & Accounts)
  const revenueItems = computed<ProfitLossItem[]>(() => {
    const revAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.type === 'REVENUE')
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    const items: ProfitLossItem[] = []

    for (const acc of revAccounts) {
      let kidinest = 0
      let ssc = 0

      // Opening balance
      if (acc.currentBalance) {
        const ob = acc.normalBalance === 'CREDIT' ? acc.currentBalance : -acc.currentBalance
        if (acc.businessUnit === 'KIDINEST' || acc.code.startsWith('4-1')) kidinest += ob
        else if (acc.businessUnit === 'SSC' || acc.code.startsWith('4-2')) ssc += ob
        else {
          kidinest += ob / 2
          ssc += ob / 2
        }
      }

      // Posted journal transactions
      for (const j of postedJournals) {
        if (accountingStore.fiscalMonth && accountingStore.fiscalMonth !== 0 && j.journalDate) {
          const str = String(j.journalDate).trim()
          const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
          const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
          const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
          if (m !== null && m !== accountingStore.fiscalMonth) continue
        }

        for (const line of j.lines) {
          if (line.accountCode === acc.code) {
            const net = line.credit - line.debit
            const isK = j.businessUnit === 'KIDINEST' || acc.businessUnit === 'KIDINEST' || acc.code.startsWith('4-1')
            const isS = j.businessUnit === 'SSC' || acc.businessUnit === 'SSC' || acc.code.startsWith('4-2')
            if (isK) kidinest += net
            else if (isS) ssc += net
            else {
              kidinest += net / 2
              ssc += net / 2
            }
          }
        }
      }

      const total = kidinest + ssc
      if (total > 0 || kidinest > 0 || ssc > 0) {
        items.push({
          code: acc.code,
          name: acc.name,
          kidinestAmount: Math.max(0, kidinest),
          sscAmount: Math.max(0, ssc),
          totalAmount: Math.max(0, total),
          percentage: 0
        })
      }
    }

    const totalRev = items.reduce((s, i) => s + i.totalAmount, 0)
    items.forEach(i => {
      i.percentage = totalRev > 0 ? Math.round((i.totalAmount / totalRev) * 1000) / 10 : 0
    })

    return items
  })

  // Dynamically compute COGS items from Buku Besar (code 5-xxxx)
  const cogsItems = computed<ProfitLossItem[]>(() => {
    const cogsAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.type === 'COGS')
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    const items: ProfitLossItem[] = []

    for (const acc of cogsAccounts) {
      let kidinest = 0
      let ssc = 0

      if (acc.currentBalance) {
        const ob = acc.normalBalance === 'DEBIT' ? acc.currentBalance : -acc.currentBalance
        if (acc.businessUnit === 'KIDINEST' || acc.code.startsWith('5-11')) kidinest += ob
        else if (acc.businessUnit === 'SSC' || acc.code.startsWith('5-12')) ssc += ob
        else {
          kidinest += ob / 2
          ssc += ob / 2
        }
      }

      for (const j of postedJournals) {
        if (accountingStore.fiscalMonth && accountingStore.fiscalMonth !== 0 && j.journalDate) {
          const str = String(j.journalDate).trim()
          const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
          const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
          const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
          if (m !== null && m !== accountingStore.fiscalMonth) continue
        }

        for (const line of j.lines) {
          if (line.accountCode === acc.code) {
            const net = line.debit - line.credit
            const isK = j.businessUnit === 'KIDINEST' || acc.businessUnit === 'KIDINEST' || acc.code.startsWith('5-11')
            const isS = j.businessUnit === 'SSC' || acc.businessUnit === 'SSC' || acc.code.startsWith('5-12')
            if (isK) kidinest += net
            else if (isS) ssc += net
            else {
              kidinest += net / 2
              ssc += net / 2
            }
          }
        }
      }

      const total = kidinest + ssc
      if (total > 0 || kidinest > 0 || ssc > 0) {
        items.push({
          code: acc.code,
          name: acc.name,
          kidinestAmount: Math.max(0, kidinest),
          sscAmount: Math.max(0, ssc),
          totalAmount: Math.max(0, total),
          percentage: 0
        })
      }
    }

    const totalCogs = items.reduce((s, i) => s + i.totalAmount, 0)
    items.forEach(i => {
      i.percentage = totalCogs > 0 ? Math.round((i.totalAmount / totalCogs) * 1000) / 10 : 0
    })

    return items
  })

  // Dynamically compute Operational Expense items from Buku Besar (code 6-xxxx)
  const operationalExpenseItems = computed<ProfitLossItem[]>(() => {
    const expAccounts = ledgerStore.accounts.filter(a => !a.isHeader && a.type === 'EXPENSE')
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    const items: ProfitLossItem[] = []

    for (const acc of expAccounts) {
      let kidinest = 0
      let ssc = 0

      if (acc.currentBalance) {
        const ob = acc.normalBalance === 'DEBIT' ? acc.currentBalance : -acc.currentBalance
        if (acc.businessUnit === 'KIDINEST' || acc.code === '6-1110') kidinest += ob
        else if (acc.businessUnit === 'SSC' || acc.code === '6-1120') ssc += ob
        else {
          kidinest += ob / 2
          ssc += ob / 2
        }
      }

      for (const j of postedJournals) {
        if (accountingStore.fiscalMonth && accountingStore.fiscalMonth !== 0 && j.journalDate) {
          const str = String(j.journalDate).trim()
          const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
          const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
          const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
          if (m !== null && m !== accountingStore.fiscalMonth) continue
        }

        for (const line of j.lines) {
          if (line.accountCode === acc.code) {
            const net = line.debit - line.credit
            const isK = j.businessUnit === 'KIDINEST' || acc.businessUnit === 'KIDINEST' || acc.code === '6-1110'
            const isS = j.businessUnit === 'SSC' || acc.businessUnit === 'SSC' || acc.code === '6-1120'
            if (isK) kidinest += net
            else if (isS) ssc += net
            else {
              kidinest += net / 2
              ssc += net / 2
            }
          }
        }
      }

      const total = kidinest + ssc
      if (total > 0 || kidinest > 0 || ssc > 0) {
        items.push({
          code: acc.code,
          name: acc.name,
          kidinestAmount: Math.max(0, kidinest),
          sscAmount: Math.max(0, ssc),
          totalAmount: Math.max(0, total),
          percentage: 0
        })
      }
    }

    const totalOpex = items.reduce((s, i) => s + i.totalAmount, 0)
    items.forEach(i => {
      i.percentage = totalOpex > 0 ? Math.round((i.totalAmount / totalOpex) * 1000) / 10 : 0
    })

    return items
  })

  // Computed P&L Totals strictly derived from Buku Besar (General Ledger)
  const profitLossSummary = computed(() => {
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    let kidinestRev = 0
    let sscRev = 0
    let kidinestCogs = 0
    let sscCogs = 0
    let kidinestOpex = 0
    let sscOpex = 0

    // 1. Account opening balances
    for (const acc of ledgerStore.accounts) {
      if (acc.isHeader) continue
      const bal = acc.currentBalance || 0
      if (bal === 0) continue

      const isK = acc.businessUnit === 'KIDINEST' || acc.code.startsWith('4-1') || acc.code.startsWith('5-11') || acc.code === '6-1110'
      const isS = acc.businessUnit === 'SSC' || acc.code.startsWith('4-2') || acc.code.startsWith('5-12') || acc.code === '6-1120'

      if (acc.type === 'REVENUE') {
        const net = acc.normalBalance === 'CREDIT' ? bal : -bal
        if (isK) kidinestRev += net
        else if (isS) sscRev += net
        else { kidinestRev += net / 2; sscRev += net / 2 }
      } else if (acc.type === 'COGS') {
        const net = acc.normalBalance === 'DEBIT' ? bal : -bal
        if (isK) kidinestCogs += net
        else if (isS) sscCogs += net
        else { kidinestCogs += net / 2; sscCogs += net / 2 }
      } else if (acc.type === 'EXPENSE') {
        const net = acc.normalBalance === 'DEBIT' ? bal : -bal
        if (isK) kidinestOpex += net
        else if (isS) sscOpex += net
        else { kidinestOpex += net / 2; sscOpex += net / 2 }
      }
    }

    // 2. All posted journal entries from Buku Besar
    for (const j of postedJournals) {
      if (accountingStore.fiscalMonth && accountingStore.fiscalMonth !== 0 && j.journalDate) {
        const str = String(j.journalDate).trim()
        const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
        const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
        const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
        if (m !== null && m !== accountingStore.fiscalMonth) continue
      }

      for (const line of j.lines) {
        const code = line.accountCode || ''
        const acc = ledgerStore.accounts.find(a => a.code === code)
        const accType = acc?.type || (
          code.startsWith('4') ? 'REVENUE' :
          code.startsWith('5') ? 'COGS' :
          code.startsWith('6') ? 'EXPENSE' : undefined
        )

        if (!accType) continue

        const isK = j.businessUnit === 'KIDINEST' || acc?.businessUnit === 'KIDINEST' || code.startsWith('4-1') || code.startsWith('5-11') || code === '6-1110'
        const isS = j.businessUnit === 'SSC' || acc?.businessUnit === 'SSC' || code.startsWith('4-2') || code.startsWith('5-12') || code === '6-1120'

        if (accType === 'REVENUE') {
          const net = line.credit - line.debit
          if (isK) kidinestRev += net
          else if (isS) sscRev += net
          else { kidinestRev += net / 2; sscRev += net / 2 }
        } else if (accType === 'COGS') {
          const net = line.debit - line.credit
          if (isK) kidinestCogs += net
          else if (isS) sscCogs += net
          else { kidinestCogs += net / 2; sscCogs += net / 2 }
        } else if (accType === 'EXPENSE') {
          const net = line.debit - line.credit
          if (isK) kidinestOpex += net
          else if (isS) sscOpex += net
          else { kidinestOpex += net / 2; sscOpex += net / 2 }
        }
      }
    }

    const totalRev = kidinestRev + sscRev
    const totalCogs = kidinestCogs + sscCogs
    const grossProfit = totalRev - totalCogs
    const kidinestGrossProfit = kidinestRev - kidinestCogs
    const sscGrossProfit = sscRev - sscCogs

    const totalOpex = kidinestOpex + sscOpex
    const netProfit = grossProfit - totalOpex
    const kidinestNetProfit = kidinestGrossProfit - kidinestOpex
    const sscNetProfit = sscGrossProfit - sscOpex

    const netMarginPercent = totalRev > 0 ? (netProfit / totalRev) * 100 : 0

    return {
      totalRev: Math.max(0, totalRev),
      kidinestRev: Math.max(0, kidinestRev),
      sscRev: Math.max(0, sscRev),
      totalCogs: Math.max(0, totalCogs),
      kidinestCogs: Math.max(0, kidinestCogs),
      sscCogs: Math.max(0, sscCogs),
      grossProfit,
      kidinestGrossProfit,
      sscGrossProfit,
      totalOpex: Math.max(0, totalOpex),
      kidinestOpex: Math.max(0, kidinestOpex),
      sscOpex: Math.max(0, sscOpex),
      netProfit,
      kidinestNetProfit,
      sscNetProfit,
      netMarginPercent: Math.round(netMarginPercent * 10) / 10
    }
  })

  // Balance Sheet Data dynamically calculated from Buku Besar (General Ledger balances)
  const balanceSheet = computed(() => {
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    const balances = new Map<string, number>()

    // 1. Initial balances from accounts
    for (const acc of ledgerStore.accounts) {
      if (!acc.isHeader) {
        balances.set(acc.code, acc.currentBalance || 0)
      }
    }

    // 2. Add debit/credit from posted journals
    for (const j of postedJournals) {
      if (accountingStore.businessUnit !== 'ALL' && j.businessUnit !== 'ALL' && j.businessUnit !== accountingStore.businessUnit) {
        continue
      }
      for (const line of j.lines) {
        const acc = ledgerStore.accounts.find(a => a.code === line.accountCode)
        const current = balances.get(line.accountCode) || 0
        const isDebitNormal = acc ? acc.normalBalance === 'DEBIT' : !line.accountCode.startsWith('2') && !line.accountCode.startsWith('3') && !line.accountCode.startsWith('4')

        if (isDebitNormal) {
          balances.set(line.accountCode, current + line.debit - line.credit)
        } else {
          balances.set(line.accountCode, current + line.credit - line.debit)
        }
      }
    }

    function sumAccounts(predicate: (code: string, acc?: any) => boolean): number {
      let sum = 0
      for (const [code, amount] of balances.entries()) {
        const acc = ledgerStore.accounts.find(a => a.code === code)
        if (predicate(code, acc)) {
          sum += amount
        }
      }
      return sum
    }

    // Kas Operasional & Bank Likuid
    const kasOperasional = sumAccounts(code => code === '1-1110' || code === '1-1111' || code === '1-1112' || (code.startsWith('1-111') && !code.startsWith('1-1113') && !code.startsWith('1-1114')))
    const bankOperasional = sumAccounts(code => code === '1-1113' || code === '1-1114' || code === '1-1120' || code.startsWith('1-112'))

    // Piutang Pasien & Klaim
    const piutangPasien = sumAccounts(code => code.startsWith('1-12'))

    // Persediaan
    const persediaan = sumAccounts(code => code.startsWith('1-13'))

    // Aset Tetap
    const asetTetapGross = sumAccounts(code => code.startsWith('1-2') && !code.includes('2300'))
    const akumulasiPenyusutan = sumAccounts(code => code.includes('2300'))
    const totalAsetTetapNet = Math.max(0, asetTetapGross - akumulasiPenyusutan)
    const totalAsetLancar = kasOperasional + bankOperasional + piutangPasien + persediaan
    const totalAset = totalAsetLancar + totalAsetTetapNet

    // Hutang Usaha / Distributor
    const hutangDistributor = sumAccounts(code => code === '2-1110' || code === '2-1115' || code.startsWith('2-111'))

    // Hutang Jasa Medis & Gaji
    const hutangDokter = sumAccounts(code => code === '2-1120' || code === '2-1125' || code === '2-1130' || code.startsWith('2-112') || code.startsWith('2-113'))

    // Hutang Pajak
    const hutangPajak = sumAccounts(code => code.startsWith('2-15'))
    const totalKewajiban = hutangDistributor + hutangDokter + hutangPajak

    // Equity
    const modalDisetor = sumAccounts(code => code === '3-1100' || code.startsWith('3-11'))
    const labaDitahan = sumAccounts(code => code === '3-1200' || code.startsWith('3-12'))
    const labaTahunBerjalan = profitLossSummary.value.netProfit
    const totalEkuitas = modalDisetor + labaDitahan + labaTahunBerjalan
    const totalKewajibanDanEkuitas = totalKewajiban + totalEkuitas
    const balanceDifference = totalAset - totalKewajibanDanEkuitas

    return {
      kasOperasional: Math.max(0, kasOperasional),
      bankOperasional: Math.max(0, bankOperasional),
      piutangPasien: Math.max(0, piutangPasien),
      persediaan: Math.max(0, persediaan),
      ppnMasukan: 0,
      totalAsetLancar: Math.max(0, totalAsetLancar),
      asetTetapGross: Math.max(0, asetTetapGross),
      akumulasiPenyusutan: Math.max(0, akumulasiPenyusutan),
      totalAsetTetapNet: Math.max(0, totalAsetTetapNet),
      totalAset: Math.max(0, totalAset),
      hutangDistributor: Math.max(0, hutangDistributor),
      hutangDokter: Math.max(0, hutangDokter),
      hutangPajak: Math.max(0, hutangPajak),
      totalKewajiban: Math.max(0, totalKewajiban),
      modalDisetor: Math.max(0, modalDisetor),
      labaDitahan: Math.max(0, labaDitahan),
      labaTahunBerjalan,
      totalEkuitas,
      totalKewajibanDanEkuitas,
      isBalanced: Math.abs(balanceDifference) < 1,
      balanceDifference
    }
  })

  // AR Aging (Live Database from PostgreSQL)
  const arAgingList = ref<AgingBucket[]>([])
  const apAgingList = ref<AgingBucket[]>([])
  const isLoadingAging = ref(false)

  async function fetchAgingReports(asOfDate?: string) {
    isLoadingAging.value = true
    try {
      const [arRes, apRes] = await Promise.all([
        accountingApi.getReceivableAging(asOfDate),
        accountingApi.getPayableAging(asOfDate)
      ])

      if (arRes && arRes.buckets && Array.isArray(arRes.buckets) && arRes.buckets.length > 0) {
        arAgingList.value = arRes.buckets.map((b: any, idx: number) => {
          const cur = Number(b.current) || 0
          const d30 = Number(b.days_1_30) || 0
          const d60 = Number(b.days_31_60) || 0
          const d90 = Number(b.days_61_90) || 0
          const dOver = Number(b.days_over_90) || 0
          const tot = Number(b.total_outstanding) || (cur + d30 + d60 + d90 + dOver)
          return {
            id: `ar-${idx + 1}`,
            entityName: b.customer_name || 'Pasien / Rekanan',
            invoiceNumber: '-',
            invoiceDate: asOfDate || new Date().toISOString().slice(0, 10),
            dueDate: asOfDate || new Date().toISOString().slice(0, 10),
            current: cur,
            days1to30: d30,
            days31to60: d60,
            days61to90: d90,
            over90: dOver,
            total: tot,
            status: (d30 + d60 + d90 + dOver) > 0 ? 'OVERDUE' : 'CURRENT'
          } as AgingBucket
        })
      } else {
        arAgingList.value = []
      }

      if (apRes && apRes.buckets && Array.isArray(apRes.buckets) && apRes.buckets.length > 0) {
        apAgingList.value = apRes.buckets.map((b: any, idx: number) => {
          const cur = Number(b.current) || 0
          const d30 = Number(b.days_1_30) || 0
          const d60 = Number(b.days_31_60) || 0
          const d90 = Number(b.days_61_90) || 0
          const dOver = Number(b.days_over_90) || 0
          const tot = Number(b.total_outstanding) || (cur + d30 + d60 + d90 + dOver)
          return {
            id: `ap-${idx + 1}`,
            entityName: b.vendor_name || 'Distributor / Supplier',
            invoiceNumber: '-',
            invoiceDate: asOfDate || new Date().toISOString().slice(0, 10),
            dueDate: asOfDate || new Date().toISOString().slice(0, 10),
            current: cur,
            days1to30: d30,
            days31to60: d60,
            days61to90: d90,
            over90: dOver,
            total: tot,
            status: (d30 + d60 + d90 + dOver) > 0 ? 'OVERDUE' : 'CURRENT'
          } as AgingBucket
        })
      } else {
        apAgingList.value = []
      }
    } catch (err) {
      console.warn('[ReportsStore] Gagal mengambil laporan aging dari backend:', err)
      arAgingList.value = []
      apAgingList.value = []
    } finally {
      isLoadingAging.value = false
    }
  }

  return {
    revenueItems,
    cogsItems,
    operationalExpenseItems,
    profitLossSummary,
    balanceSheet,
    arAgingList,
    apAgingList,
    isLoadingAging,
    fetchAgingReports
  }
})
