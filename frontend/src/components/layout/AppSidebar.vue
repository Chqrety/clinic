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
  badge?: number | string
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
      { name: 'Ikhtisar Keuangan', to: '/accounting/dashboard', icon: LayoutDashboard }
    ]
  },
  {
    groupName: 'MODUL UTAMA (SIAP PAKAI)',
    items: [
      { name: 'Bagan Akun (COA)', to: '/accounting/chart-of-accounts', icon: FileSpreadsheet },
      { name: 'Saldo Awal Akuntansi', to: '/accounting/opening-balance', icon: BookmarkCheck },
      { name: 'Jurnal Umum (Tabel)', to: '/accounting/journals', icon: BookOpen },
      { name: 'Buku Kas & Bank', to: '/accounting/cash-bank', icon: Coins },
      { name: 'Buku Besar (Ledger)', to: '/accounting/ledger', icon: FileText },
      { name: 'SPT Masa (Bulanan)', to: '/accounting/tax/spt-masa', icon: Calculator },
      { name: 'SPT Tahunan (1771)', to: '/accounting/tax/spt-tahunan', icon: Receipt }
    ]
  },
  {
    groupName: 'LAPORAN KEUANGAN & AUDIT',
    items: [
      { name: 'Laba Rugi (P&L)', to: '/accounting/reports/profit-loss', icon: FileSpreadsheet },
      { name: 'Laporan Neraca', to: '/accounting/reports/balance-sheet', icon: Scale },
      { name: 'Arus Kas (Cash Flow)', to: '/accounting/reports/cash-flow', icon: Coins },
      { name: 'Umur Piutang (Aging)', to: '/accounting/reports/aging', icon: History },
      {
        name: 'Transaksi Kasir Otomatis',
        to: '/accounting/transactions',
        icon: Layers
      },
      {
        name: 'Approval Bertingkat',
        to: '/accounting/journals/approvals',
        icon: FileCheck2
      },
      {
        name: 'Tutup Buku Periode',
        to: '/accounting/periods',
        icon: CalendarCheck
      },
      {
        name: 'Jejak Audit Sistem',
        to: '/accounting/audit-trail',
        icon: ShieldAlert
      }
    ]
  }
])

function isCurrentRoute(to: string): boolean {
  if (to === '/dashboard' && (route.path === '/' || route.path === '/dashboard')) return true
  if (to === '/accounting/tax/spt-masa' && route.path.startsWith('/accounting/tax/spt-masa')) return true
  if (to !== '/dashboard' && route.path === to) return true
  if (to !== '/dashboard' && route.path.startsWith(to + '/')) return true
  return false
}
</script>

<template>
  <aside
    class="w-64 bg-zinc-950 border-r border-zinc-800 flex flex-col justify-between shrink-0 select-none text-zinc-300 min-h-screen"
  >
    <!-- Brand Header (Official SMG Clinical Identity) -->
    <div class="p-4 border-b border-zinc-800/80 flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-lg bg-white text-zinc-950 flex items-center justify-center font-black text-sm shadow-xs">
          SMG
        </div>
        <div>
          <div class="text-xs font-bold text-zinc-100 leading-tight">Semarang Medical Group</div>
          <div class="text-[10px] text-zinc-400 font-medium">Finance & Accounting</div>
        </div>
      </div>
      <span class="text-[9px] font-bold uppercase tracking-wider px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-300 border border-zinc-700">
        PROD
      </span>
    </div>

    <!-- Active Unit Banner Pill (Semantic Dual-Brand Indicators) -->
    <div class="px-4 py-2.5 bg-zinc-900/60 border-b border-zinc-800/60 flex items-center justify-between text-xs">
      <span class="text-zinc-400 text-[11px] font-medium">Entitas Aktif:</span>
      <span
        v-if="accountingStore.businessUnit === 'ALL'"
        class="font-semibold text-zinc-200 inline-flex items-center gap-1.5"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-zinc-400"></span>
        Konsolidasi
      </span>
      <span
        v-else-if="accountingStore.businessUnit === 'KIDINEST'"
        class="font-bold text-amber-400 inline-flex items-center gap-1.5"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
        Kidinest (Anak)
      </span>
      <span
        v-else
        class="font-bold text-rose-400 inline-flex items-center gap-1.5"
      >
        <span class="w-1.5 h-1.5 rounded-full bg-rose-400"></span>
        SSC (Dermatologi)
      </span>
    </div>

    <!-- Navigation Menu Items -->
    <div class="flex-1 overflow-y-auto px-3 py-4 space-y-6 text-xs scrollbar-thin scrollbar-thumb-zinc-800">
      <div v-for="group in navGroups" :key="group.groupName" class="space-y-1">
        <h4 class="px-3 text-[10px] font-bold text-zinc-500 uppercase tracking-wider mb-1.5">
          {{ group.groupName }}
        </h4>

        <router-link
          v-for="item in group.items"
          :key="item.to"
          :to="item.to"
          :class="[
            'flex items-center justify-between px-3 py-2 rounded-xl font-medium transition-all duration-150 group',
            isCurrentRoute(item.to)
              ? 'bg-zinc-800 text-white font-bold shadow-2xs'
              : 'text-zinc-400 hover:text-white hover:bg-zinc-900/80'
          ]"
        >
          <div class="flex items-center gap-2.5 truncate">
            <component
              :is="item.icon"
              :class="[
                'w-4 h-4 shrink-0 transition-colors',
                isCurrentRoute(item.to) ? 'text-white' : 'text-zinc-400 group-hover:text-zinc-200'
              ]"
            />
            <span class="truncate">{{ item.name }}</span>
          </div>

          <span
            v-if="item.badge && (typeof item.badge === 'string' || item.badge > 0)"
            :class="[
              'px-1.5 py-0.5 rounded text-[10px] font-bold shrink-0',
              item.badgeColor || 'bg-zinc-800 text-zinc-300 border border-zinc-700/60'
            ]"
          >
            {{ item.badge }}
          </span>
        </router-link>
      </div>
    </div>

    <!-- Bottom User Profile & Switcher -->
    <div class="p-4 border-t border-zinc-800/80 bg-zinc-950 space-y-3">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-zinc-800 border border-zinc-700 flex items-center justify-center font-bold text-xs text-zinc-200">
          {{ (accountingStore.currentUserName || 'F').charAt(0) }}
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold text-zinc-100 truncate">
            {{ accountingStore.currentUserName }}
          </p>
          <p class="text-[10px] text-zinc-400 truncate capitalize">
            Finance & Accounting Manager
          </p>
        </div>
      </div>
    </div>
  </aside>
</template>
