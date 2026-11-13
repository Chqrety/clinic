<script setup lang="ts">
import { ref } from 'vue'
import { useAccountingStore } from '@/stores/accounting'
import type { BusinessUnit, BranchCode, UserRole } from '@/types/accounting'
import {
  Bell,
  Building,
  Calendar,
  Check,
  UserCheck,
  Layers,
  ChevronDown
} from 'lucide-vue-next'

const store = useAccountingStore()
const showNotifications = ref(false)
const showRoleMenu = ref(false)

function onUnitChange(e: Event) {
  const target = e.target as HTMLSelectElement
  store.setBusinessUnit(target.value as BusinessUnit)
}

function onBranchChange(e: Event) {
  const target = e.target as HTMLSelectElement
  store.setBranch(target.value as BranchCode)
}

function selectRole(role: UserRole) {
  store.setRole(role)
  showRoleMenu.value = false
  store.showToast(`Role aktif berganti ke: ${role.replace('_', ' ').toUpperCase()}`, 'info')
}
</script>

<template>
  <header
    class="h-16 bg-white border-b border-slate-200/90 px-6 flex items-center justify-between sticky top-0 z-30 shadow-2xs"
  >
    <!-- Left: Global Entity & Branch Controls -->
    <div class="flex items-center gap-3">
      <!-- Business Unit Selector -->
      <div class="flex items-center bg-slate-100 p-0.5 rounded-lg border border-slate-200">
        <button
          v-for="unit in store.unitOptions"
          :key="unit.code"
          @click="store.setBusinessUnit(unit.code as BusinessUnit)"
          :class="[
            'px-3 py-1.5 rounded-md text-xs font-semibold transition-all duration-150',
            store.businessUnit === unit.code
              ? 'bg-white text-blue-700 shadow-xs border border-slate-200/60'
              : 'text-slate-600 hover:text-slate-900'
          ]"
        >
          {{ unit.badge }}: {{ unit.code === 'ALL' ? 'Konsolidasi' : unit.code }}
        </button>
      </div>

      <!-- Branch Dropdown -->
      <div class="relative flex items-center">
        <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
          <Building class="w-3.5 h-3.5" />
        </div>
        <select
          :value="store.branchCode"
          @change="onBranchChange"
          class="pl-8 pr-7 py-1.5 bg-slate-50 border border-slate-200 text-slate-800 text-xs rounded-lg font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none appearance-none cursor-pointer"
        >
          <option
            v-for="branch in store.branchOptions"
            :key="branch.code"
            :value="branch.code"
          >
            {{ branch.name }}
          </option>
        </select>
        <ChevronDown class="w-3.5 h-3.5 text-slate-400 absolute right-2 pointer-events-none" />
      </div>

      <!-- Fiscal Period Selector -->
      <div class="relative flex items-center">
        <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-slate-400">
          <Calendar class="w-3.5 h-3.5" />
        </div>
        <select
          :value="store.fiscalMonth"
          @change="(e: any) => store.fiscalMonth = Number(e.target.value)"
          class="pl-8 pr-7 py-1.5 bg-slate-50 border border-slate-200 text-slate-800 text-xs rounded-lg font-medium focus:ring-2 focus:ring-blue-500 focus:outline-none appearance-none cursor-pointer"
        >
          <option :value="3">Maret 2026 (Open)</option>
          <option :value="2">Februari 2026 (Locked)</option>
          <option :value="1">Januari 2026 (Locked)</option>
        </select>
        <ChevronDown class="w-3.5 h-3.5 text-slate-400 absolute right-2 pointer-events-none" />
      </div>
    </div>

    <!-- Right: Role Switcher & Notifications -->
    <div class="flex items-center gap-3">
      <!-- Role Simulation Selector -->
      <div class="relative">
        <button
          @click="showRoleMenu = !showRoleMenu"
          class="flex items-center gap-2 px-3 py-1.5 rounded-lg border border-slate-200 bg-slate-50 hover:bg-slate-100 text-xs font-medium text-slate-700 transition-colors"
        >
          <UserCheck class="w-3.5 h-3.5 text-blue-600" />
          <span class="max-w-[130px] truncate">{{ store.roleOptions.find(r => r.role === store.currentUserRole)?.label.split(' ')[0] }}</span>
          <ChevronDown class="w-3 h-3 text-slate-400" />
        </button>

        <!-- Role Dropdown Menu -->
        <div
          v-if="showRoleMenu"
          class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-xl border border-slate-200 py-1.5 z-50 text-xs font-normal"
        >
          <div class="px-3 py-2 border-b border-slate-100">
            <p class="font-bold text-slate-900">Simulasi Hak Akses (Role)</p>
            <p class="text-[11px] text-slate-500">Uji coba permission approval, posting & audit</p>
          </div>
          <button
            v-for="r in store.roleOptions"
            :key="r.role"
            @click="selectRole(r.role)"
            class="w-full text-left px-3 py-2 hover:bg-blue-50 flex items-center justify-between text-slate-700 hover:text-blue-700 transition-colors"
          >
            <span :class="{'font-bold text-blue-700': store.currentUserRole === r.role}">
              {{ r.label }}
            </span>
            <Check v-if="store.currentUserRole === r.role" class="w-4 h-4 text-blue-600" />
          </button>
        </div>
      </div>

      <!-- Notifications Bell -->
      <div class="relative">
        <button
          @click="showNotifications = !showNotifications"
          class="relative p-2 rounded-lg text-slate-600 hover:text-slate-900 hover:bg-slate-100 transition-colors"
          aria-label="Notifikasi"
        >
          <Bell class="w-4 h-4" />
          <span
            v-if="store.unreadNotificationCount > 0"
            class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white"
          ></span>
        </button>

        <!-- Notifications Dropdown -->
        <div
          v-if="showNotifications"
          class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-xl border border-slate-200 p-3 z-50 text-xs"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100 mb-2">
            <h4 class="font-bold text-slate-900">Pemberitahuan Sistem</h4>
            <button
              @click="store.markAllNotificationsRead()"
              class="text-[11px] text-blue-600 hover:underline"
            >
              Tandai Dibaca
            </button>
          </div>

          <div class="space-y-2 max-h-72 overflow-y-auto">
            <div
              v-for="notif in store.notifications"
              :key="notif.id"
              class="p-2.5 rounded-lg bg-slate-50 hover:bg-slate-100 transition-colors border border-slate-100 space-y-1"
            >
              <div class="flex items-center justify-between">
                <span class="font-bold text-slate-800">{{ notif.title }}</span>
                <span class="text-[10px] text-slate-400">{{ notif.timestamp }}</span>
              </div>
              <p class="text-slate-600 text-[11px] leading-relaxed">{{ notif.message }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
