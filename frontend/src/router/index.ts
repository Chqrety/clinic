import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AccountingLayout from '@/layouts/AccountingLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: AccountingLayout,
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/accounting/DashboardView.vue'),
        meta: { title: 'Ikhtisar Keuangan' }
      },
      {
        path: 'transactions',
        name: 'TransactionSources',
        component: () => import('@/views/accounting/TransactionSourcesView.vue'),
        meta: { title: 'Sumber Transaksi Operasional' }
      },
      {
        path: 'journals',
        name: 'JournalList',
        component: () => import('@/views/accounting/JournalListView.vue'),
        meta: { title: 'Daftar Jurnal Umum' }
      },
      {
        path: 'journals/approvals',
        name: 'JournalApprovals',
        component: () => import('@/views/accounting/JournalApprovalsView.vue'),
        meta: { title: 'Review & Approval Jurnal' }
      },
      {
        path: 'ledger',
        name: 'GeneralLedger',
        component: () => import('@/views/accounting/GeneralLedgerView.vue'),
        meta: { title: 'Buku Besar (General Ledger)' }
      },
      {
        path: 'cash-bank',
        name: 'CashBankBook',
        component: () => import('@/views/accounting/CashBankBookView.vue'),
        meta: { title: 'Buku Kas & Bank' }
      },
      {
        path: 'trial-balance',
        name: 'TrialBalance',
        component: () => import('@/views/accounting/TrialBalanceView.vue'),
        meta: { title: 'Neraca Saldo' }
      },
      {
        path: 'chart-of-accounts',
        name: 'ChartOfAccounts',
        component: () => import('@/views/accounting/ChartOfAccountsView.vue'),
        meta: { title: 'Bagan Akun (COA)' }
      },
      {
        path: 'periods',
        name: 'AccountingPeriods',
        component: () => import('@/views/accounting/AccountingPeriodView.vue'),
        meta: { title: 'Tutup Buku Periode' }
      },
      {
        path: 'reports/profit-loss',
        name: 'ProfitLossReport',
        component: () => import('@/views/accounting/reports/ProfitLossView.vue'),
        meta: { title: 'Laporan Laba Rugi' }
      },
      {
        path: 'reports/balance-sheet',
        name: 'BalanceSheetReport',
        component: () => import('@/views/accounting/reports/BalanceSheetView.vue'),
        meta: { title: 'Laporan Neraca' }
      },
      {
        path: 'reports/cash-flow',
        name: 'CashFlowReport',
        component: () => import('@/views/accounting/reports/CashFlowView.vue'),
        meta: { title: 'Laporan Arus Kas' }
      },
      {
        path: 'reports/aging',
        name: 'AgingReport',
        component: () => import('@/views/accounting/reports/AgingView.vue'),
        meta: { title: 'Laporan Umur Piutang & Hutang' }
      },
      {
        path: 'tax',
        name: 'TaxManagement',
        component: () => import('@/views/accounting/TaxView.vue'),
        meta: { title: 'Manajemen Pajak' }
      },
      {
        path: 'mapping-rules',
        name: 'MappingRules',
        component: () => import('@/views/accounting/MappingRulesView.vue'),
        meta: { title: 'Aturan Pemetaan Akun' }
      },
      {
        path: 'audit-trail',
        name: 'AuditTrail',
        component: () => import('@/views/accounting/AuditTrailView.vue'),
        meta: { title: 'Jejak Audit Sistem' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to, _from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} - Sistem Akuntansi Klinik (Kidinest & SSC)`
  }
  next()
})

export default router

