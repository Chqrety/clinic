import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AccountingLayout from '@/layouts/AccountingLayout.vue'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  // ==================== CLINIC OPERATIONAL & PUBLIC ROUTES ====================
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/LandingKidinest.vue'),
    meta: { title: 'Kidinest Pediatric Care' }
  },
  {
    path: '/kidinest',
    name: 'Kidinest',
    component: () => import('@/views/LandingKidinest.vue'),
    meta: { title: 'Kidinest Pediatric Care' }
  },
  {
    path: '/ssc',
    name: 'SSC',
    component: () => import('@/views/LandingSSC.vue'),
    meta: { title: 'Semarang Skin Center' }
  },
  {
    path: '/kiosk',
    name: 'PatientPortal',
    component: () => import('@/views/PatientPortal.vue'),
    props: { brand: 'kidinest' },
    meta: { title: 'Reservasi & Pendaftaran Pasien' }
  },
  {
    path: '/kiosk-kidinest',
    name: 'PatientPortalKidinest',
    component: () => import('@/views/PatientPortal.vue'),
    props: { brand: 'kidinest' },
    meta: { title: 'Reservasi Pasien Kidinest' }
  },
  {
    path: '/kiosk-ssc',
    name: 'PatientPortalSSC',
    component: () => import('@/views/PatientPortal.vue'),
    props: { brand: 'ssc' },
    meta: { title: 'Reservasi Pasien Semarang Skin Center' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { title: 'Portal Staf Medis' }
  },
  {
    path: '/login-kidinest',
    name: 'LoginKidinest',
    component: () => import('@/views/LoginView.vue'),
    props: { defaultBrand: 'kidinest' },
    meta: { title: 'Login Staf Kidinest' }
  },
  {
    path: '/login-ssc',
    name: 'LoginSSC',
    component: () => import('@/views/LoginView.vue'),
    props: { defaultBrand: 'ssc' },
    meta: { title: 'Login Staf SSC' }
  },
  {
    path: '/admin',
    name: 'AdminPortal',
    component: () => import('@/views/AdminPortal.vue'),
    meta: { title: 'Workstation Staf Medis' }
  },
  {
    path: '/monitor',
    name: 'QueueMonitor',
    component: () => import('@/views/QueueMonitor.vue'),
    meta: { title: 'TV Display Antrian' }
  },

  // ==================== ISOLATED ACCOUNTING & FINANCIAL ERP ROUTES ====================
  {
    path: '/accounting',
    component: AccountingLayout,
    children: [
      {
        path: '',
        redirect: '/accounting/dashboard'
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
        path: 'opening-balance',
        name: 'OpeningBalance',
        component: () => import('@/views/accounting/OpeningBalanceView.vue'),
        meta: { title: 'Pembuatan Saldo Awal Akuntansi' }
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
        redirect: '/accounting/tax/spt-masa'
      },
      {
        path: 'tax/spt-masa',
        name: 'SptMasaManagement',
        component: () => import('@/views/accounting/tax/SptMasaView.vue'),
        meta: { title: 'SPT Masa (Pajak Bulanan)' }
      },
      {
        path: 'tax/spt-masa/create',
        name: 'SptMasaCreate',
        component: () => import('@/views/accounting/tax/SptMasaCreateView.vue'),
        meta: { title: 'Formulir SPT Masa Baru' }
      },
      {
        path: 'tax/spt-tahunan',
        name: 'SptTahunanManagement',
        component: () => import('@/views/accounting/tax/SptTahunanView.vue'),
        meta: { title: 'SPT Tahunan PPh Badan (1771)' }
      },
      {
        path: 'tax/spt-tahunan/edit',
        name: 'SptTahunanEdit',
        component: () => import('@/views/accounting/tax/SptTahunanEditView.vue'),
        meta: { title: 'Formulir Penyesuaian SPT 1771' }
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
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach(async (to: any, _from: any, next: any) => {
  if (to.meta?.title) {
    document.title = `${to.meta.title} - Sistem Klinik & Akuntansi`
  }

  // Route Guard: Lindungi seluruh rute /accounting/* wajib login role keuangan atau owner
  if (to.path.startsWith('/accounting')) {
    const authStore = useAuthStore()

    if (!authStore.isAuthenticated) {
      const ok = await authStore.fetchMe()
      if (!ok) {
        return next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
      }
    }

    // Role yang memiliki akses ke modul accounting: FINANCE, SUPERADMIN_OWNER, OWNER
    const allowedRoles = ['FINANCE', 'SUPERADMIN_OWNER', 'OWNER']
    const currentRole = (authStore.user?.role || '').toUpperCase()
    if (authStore.user && !allowedRoles.includes(currentRole)) {
      return next({ path: '/admin' })
    }
  }

  next()
})

export default router
