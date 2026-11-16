import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { BusinessUnit, BranchCode, UserRole } from '@/types/accounting'

export interface AppNotification {
  id: string
  title: string
  message: string
  type: 'info' | 'warning' | 'success' | 'danger'
  timestamp: string
  read: boolean
  link?: string
}

export const useAccountingStore = defineStore('accounting', () => {
  // Global Filters
  const businessUnit = ref<BusinessUnit>('ALL')
  const branchCode = ref<BranchCode>('ALL')
  const fiscalYear = ref<number>(2026)
  const fiscalMonth = ref<number>(3) // March 2026

  // User Authentication & Role Context
  const currentUserRole = ref<UserRole>('finance_manager')
  const currentUserName = ref<string>('Budi Santoso, SE, Ak., CA')
  const currentUserEmail = ref<string>('budi.santoso@kidinest-ssc.id')

  // Notification Feed
  const notifications = ref<AppNotification[]>([
    {
      id: 'notif-1',
      title: 'Perlu Approval Jurnal',
      message: '2 draft jurnal penutupan kasir shift malam memerlukan approval Anda.',
      type: 'warning',
      timestamp: '10 menit yang lalu',
      read: false,
      link: '/journals?status=REVIEWED'
    },
    {
      id: 'notif-2',
      title: 'Transaksi Kasir Baru',
      message: 'Closing kasir SSC Cabang Kebayoran telah masuk dari modul kasir.',
      type: 'info',
      timestamp: '25 menit yang lalu',
      read: false,
      link: '/transactions'
    },
    {
      id: 'notif-3',
      title: 'Tutup Buku Periode',
      message: 'Periode Februari 2026 telah berhasil dikunci oleh Super Admin.',
      type: 'success',
      timestamp: 'Kemarin',
      read: true,
      link: '/periods'
    }
  ])

  // System toast
  const toastMessage = ref<{ text: string; type: 'success' | 'error' | 'info' } | null>(null)
  let toastTimer: any = null

  function showToast(text: string, type: 'success' | 'error' | 'info' = 'success') {
    if (toastTimer) clearTimeout(toastTimer)
    toastMessage.value = { text, type }
    toastTimer = setTimeout(() => {
      toastMessage.value = null
    }, 4000)
  }

  // Branch names mapping
  const branchOptions = [
    { code: 'ALL', name: 'Semua Cabang' },
    { code: 'BRANCH-MAIN', name: 'Kebayoran Baru (Pusat)' },
    { code: 'BRANCH-JAKSEL', name: 'Pondok Indah' },
    { code: 'BRANCH-SBY', name: 'Surabaya Barat' }
  ]

  const unitOptions = [
    { code: 'ALL', name: 'Semua Unit (Konsolidasi)', badge: 'ALL' },
    { code: 'KIDINEST', name: 'Kidinest (Klinik Anak)', badge: 'KID' },
    { code: 'SSC', name: 'SSC (Skin & Slimming)', badge: 'SSC' }
  ]

  const roleOptions: { role: UserRole; label: string }[] = [
    { role: 'finance_manager', label: 'Finance Manager (Approval & Posting)' },
    { role: 'accounting_supervisor', label: 'Accounting Supervisor (Reviewer)' },
    { role: 'accounting_staff', label: 'Accounting Staff (Maker / Draft)' },
    { role: 'tax_staff', label: 'Tax Specialist (Pajak PPh & PPN)' },
    { role: 'auditor', label: 'Auditor (Read-Only / Inspection)' },
    { role: 'management', label: 'Executive Management (Laporan)' },
    { role: 'super_admin', label: 'Super Admin (Full Akses)' }
  ]

  const unreadNotificationCount = computed(() => notifications.value.filter(n => !n.read).length)

  function markAllNotificationsRead() {
    notifications.value.forEach(n => {
      n.read = true
    })
  }

  function setBusinessUnit(unit: BusinessUnit) {
    businessUnit.value = unit
  }

  function setBranch(branch: BranchCode) {
    branchCode.value = branch
  }

  function setRole(role: UserRole) {
    currentUserRole.value = role
    switch (role) {
      case 'finance_manager':
        currentUserName.value = 'Budi Santoso, SE, Ak., CA'
        break
      case 'accounting_supervisor':
        currentUserName.value = 'Dewi Lestari, SE, Ak.'
        break
      case 'accounting_staff':
        currentUserName.value = 'Rian Pratama, S.Ak.'
        break
      case 'tax_staff':
        currentUserName.value = 'Siti Nurhaliza, B.K.P.'
        break
      case 'auditor':
        currentUserName.value = 'KAP Tanubrata & Rekan'
        break
      case 'management':
        currentUserName.value = 'dr. Hendra Wijaya, Sp.A (Direktur)'
        break
      case 'super_admin':
        currentUserName.value = 'Administrator Sistem'
        break
    }
  }

  // Permission helpers
  const canCreateDraft = computed(() =>
    ['super_admin', 'finance_manager', 'accounting_supervisor', 'accounting_staff'].includes(currentUserRole.value)
  )

  const canReviewJournal = computed(() =>
    ['super_admin', 'finance_manager', 'accounting_supervisor'].includes(currentUserRole.value)
  )

  const canApproveJournal = computed(() =>
    ['super_admin', 'finance_manager'].includes(currentUserRole.value)
  )

  const canPostJournal = computed(() =>
    ['super_admin', 'finance_manager'].includes(currentUserRole.value)
  )

  const canReverseJournal = computed(() =>
    ['super_admin', 'finance_manager'].includes(currentUserRole.value)
  )

  const canLockPeriod = computed(() =>
    ['super_admin', 'finance_manager'].includes(currentUserRole.value)
  )

  return {
    businessUnit,
    branchCode,
    fiscalYear,
    fiscalMonth,
    currentUserRole,
    currentUserName,
    currentUserEmail,
    notifications,
    unreadNotificationCount,
    toastMessage,
    branchOptions,
    unitOptions,
    roleOptions,
    showToast,
    markAllNotificationsRead,
    setBusinessUnit,
    setBranch,
    setRole,
    canCreateDraft,
    canReviewJournal,
    canApproveJournal,
    canPostJournal,
    canReverseJournal,
    canLockPeriod
  }
})
