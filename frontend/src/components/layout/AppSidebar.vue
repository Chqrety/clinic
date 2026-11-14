<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useJournalsStore } from '@/stores/journals'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountingStore } from '@/stores/accounting'
import {
  LayoutDashboard,
  Layers,
  BookOpen,
  FileCheck2,
  FileText,
  Calculator,
  Sliders,
  ShieldAlert,
  Building2,
  Receipt,
  FileSpreadsheet,
  Coins,
  History,
  Scale,
  CalendarCheck,
  BookmarkCheck
} from 'lucide-vue-next'

const route = useRoute()
const journalsStore = useJournalsStore()
const transactionsStore = useTransactionsStore()
const accountingStore = useAccountingStore()

// Badge counts
const unprocessedTxCount = computed(() => {
  return transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').length
})

const pendingApprovalCount = computed(() => {
  return journalsStore.journals.filter(j => j.status === 'REVIEWED' || j.status === 'SUBMITTED').length
})

interface NavItem {
  name: string
  to: string
  icon: any
  badge?: number
  badgeColor?: string
}

interface NavGroup {
  groupName: string
  items: NavItem[]
}

const navGroups = computed<NavGroup[]>(() => [
  {
    groupName: 'DASHBOARD',
    items: [
      { name: 'Ikhtisar Keuangan', to: '/dashboard', icon: LayoutDashboard }
    ]
  },
  {
    groupName: 'SUMBER TRANSAKSI',
    items: [
      {
        name: 'Transaksi Operasional',
        to: '/transactions',
        icon: Layers,
        badge: unprocessedTxCount.value,
        badgeColor: 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
      }
    ]
  },
  {
    groupName: 'MANAJEMEN JURNAL',
    items: [
      { name: 'Daftar Jurnal Umum', to: '/journals', icon: BookOpen },
      {
        name: 'Antrian Review & Approval',
        to: '/journals/approvals',
        icon: FileCheck2,
        badge: pendingApprovalCount.value,
        badgeColor: 'bg-sky-500/20 text-sky-400 border border-sky-500/30'
      }
    ]
  },
  {
    groupName: 'BUKU BESAR & COA',
    items: [
      { name: 'Buku Besar (Ledger)', to: '/ledger', icon: FileText },
      { name: 'Buku Kas & Bank', to: '/cash-bank', icon: Coins },
      { name: 'Neraca Saldo (Trial Balance)', to: '/trial-balance', icon: Scale },
      { name: 'Bagan Akun (COA)', to: '/chart-of-accounts', icon: FileSpreadsheet },
      { name: 'Saldo Awal Akuntansi', to: '/opening-balance', icon: BookmarkCheck },
      { name: 'Tutup Buku Periode', to: '/periods', icon: CalendarCheck }
    ]
  },
  {
    groupName: 'LAPORAN KEUANGAN',
    items: [
      { name: 'Laba Rugi (P&L)', to: '/reports/profit-loss', icon: Receipt },
      { name: 'Neraca (Balance Sheet)', to: '/reports/balance-sheet', icon: Scale },
      { name: 'Arus Kas (Cash Flow)', to: '/reports/cash-flow', icon: Coins },
      { name: 'Umur Piutang & Hutang', to: '/reports/aging', icon: History }
    ]
  },
  {
    groupName: 'PAJAK & KEPATUHAN',
    items: [
      { name: 'PPh 21 Dokter & PPN', to: '/tax', icon: Calculator }
    ]
  },
  {
    groupName: 'MASTER & AUDIT',
    items: [
      { name: 'Aturan Pemetaan Akun', to: '/mapping-rules', icon: Sliders },
      { name: 'Jejak Audit Sistem', to: '/audit-trail', icon: ShieldAlert }
    ]
  }
])

function isCurrentRoute(to: string): boolean {
  if (to === '/dashboard' && (route.path === '/' || route.path === '/dashboard')) return true
  if (to !== '/dashboard' && route.path.startsWith(to)) return true
  return false
}
</script>

<template>
  <aside
    class="w-64 bg-slate-900 border-r border-slate-800 flex flex-col justify-between shrink-0 select-none text-slate-300 min-h-screen"
  >
    <!-- Brand Header -->
    <div class="px-5 py-5 border-b border-slate-800/80 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-blue-600 to-indigo-600 flex items-center justify-center text-white shadow-lg shadow-blue-500/20 font-bold text-lg">
          <Building2 class="w-5 h-5 text-white" />
        </div>
        <div>
          <h1 class="text-sm font-bold text-white tracking-wide">CLINIC ERP</h1>
          <p class="text-[11px] font-medium text-slate-400">Integrated Accounting</p>
        </div>
      </div>
      <span class="text-[10px] font-semibold uppercase px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/20">
        v2.4
      </span>
    </div>

    <!-- Active Unit Banner Pill -->
    <div class="px-4 py-2.5 bg-slate-950/60 border-b border-slate-800/60 flex items-center justify-between text-xs">
      <span class="text-slate-400 text-[11px] font-medium">Entitas Aktif:</span>
      <span
        v-if="accountingStore.businessUnit === 'ALL'"
        class="font-semibold text-slate-200"
      >
        Kidinest + SSC (Konsolidasi)
      </span>
      <span
        v-else-if="accountingStore.businessUnit === 'KIDINEST'"
        class="font-semibold text-blue-400"
      >
        Kidinest (Klinik Anak)
      </span>
      <span
        v-else
        class="font-semibold text-rose-400"
      >
        SSC (Skin & Slimming)
      </span>
    </div>

    <!-- Navigation Menu Items -->
    <div class="flex-1 overflow-y-auto px-3 py-4 space-y-6 text-xs scrollbar-thin scrollbar-thumb-slate-800">
      <div v-for="group in navGroups" :key="group.groupName" class="space-y-1">
        <h4 class="px-3 text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
          {{ group.groupName }}
        </h4>

        <router-link
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          :class="[
            'flex items-center justify-between px-3 py-2 rounded-lg font-medium transition-all duration-150 group',
            isCurrentRoute(item.to)
              ? 'bg-blue-600 text-white shadow-sm'
              : 'text-slate-300 hover:text-white hover:bg-slate-800/80'
          ]"
        >
          <div class="flex items-center gap-2.5 truncate">
            <component
              :is="item.icon"
              :class="[
                'w-4 h-4 shrink-0 transition-colors',
                isCurrentRoute(item.to) ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'
              ]"
            />
            <span class="truncate">{{ item.name }}</span>
          </div>

          <span
            v-if="item.badge && item.badge > 0"
            :class="[
              'px-1.5 py-0.2 rounded-full text-[10px] font-bold shrink-0',
              item.badgeColor || 'bg-blue-500/20 text-blue-400'
            ]"
          >
            {{ item.badge }}
          </span>
        </router-link>
      </div>
    </div>

    <!-- Bottom User Profile & Switcher -->
    <div class="p-4 border-t border-slate-800/80 bg-slate-950/40 space-y-3">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-xs text-blue-400">
          {{ accountingStore.currentUserName.charAt(0) }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-white truncate">
            {{ accountingStore.currentUserName }}
          </p>
          <p class="text-[10px] text-slate-400 truncate capitalize">
            {{ accountingStore.currentUserRole.replace('_', ' ') }}
          </p>
        </div>
      </div>
    </div>
  </aside>
</template>
