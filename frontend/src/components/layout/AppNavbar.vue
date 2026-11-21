<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useAuthStore } from '@/stores/auth'
import { useClinicStore } from '@/stores/clinic'
import { useLocaleStore } from '@/stores/locale'
import type { BusinessUnit } from '@/types/accounting'
import {
  Bell,
  Calendar,
  ChevronDown,
  LogOut,
  ShieldCheck,
  Stethoscope,
  Settings,
  Menu
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
}>()

const router = useRouter()
const store = useAccountingStore()
const authStore = useAuthStore()
const localeStore = useLocaleStore()

const showNotifications = ref(false)
const showMobileUnitMenu = ref(false)

async function handleLogout() {
  window.location.hash = ''
  await authStore.logout()
  const clinicStore = useClinicStore()
  clinicStore.setRole('patient')
  store.showToast('Berhasil logout dari sistem akuntansi', 'info')
  router.replace('/login-kidinest')
}
</script>

<template>
  <header
    class="h-16 bg-white border-b border-zinc-200/90 px-3 sm:px-6 flex items-center justify-between sticky top-0 z-30 shadow-2xs gap-2 sm:gap-4 select-none"
  >
    <!-- Left: Mobile Menu Toggle, Global Entity (Unit Bisnis) & Fiscal Controls -->
    <div class="flex items-center gap-2 sm:gap-3 min-w-0">
      <!-- Mobile Sidebar Toggle Button -->
      <button
        @click="emit('toggle-sidebar')"
        class="lg:hidden p-2 rounded-xl text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 transition-colors cursor-pointer shrink-0"
        title="Buka Navigasi Akuntansi"
      >
        <Menu class="w-5 h-5" />
      </button>

      <!-- Mobile Business Unit Selector (Clean Pill with Dropdown, Never Cut Off) -->
      <div class="relative sm:hidden shrink-0">
        <button
          @click="showMobileUnitMenu = !showMobileUnitMenu"
          class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border text-xs font-bold transition-all shadow-2xs cursor-pointer"
          :class="
            store.businessUnit === 'KIDINEST'
              ? 'bg-blue-50 text-blue-900 border-blue-200 hover:bg-blue-100'
              : store.businessUnit === 'SSC'
              ? 'bg-rose-50 text-rose-900 border-rose-200 hover:bg-rose-100'
              : store.businessUnit === 'MINIMO'
              ? 'bg-amber-50 text-amber-900 border-amber-200 hover:bg-amber-100'
              : 'bg-zinc-900 text-white border-zinc-900 hover:bg-zinc-800'
          "
        >
          <span
            class="w-1.5 h-1.5 rounded-full shrink-0"
            :class="
              store.businessUnit === 'KIDINEST' ? 'bg-blue-500' :
              store.businessUnit === 'SSC' ? 'bg-rose-500' :
              store.businessUnit === 'MINIMO' ? 'bg-amber-500' : 'bg-emerald-400'
            "
          ></span>
          <span class="truncate max-w-[90px]">{{ store.businessUnit === 'ALL' ? 'Konsolidasi' : store.businessUnit }}</span>
          <ChevronDown class="w-3.5 h-3.5 opacity-70 shrink-0" />
        </button>

        <!-- Mobile Unit Dropdown Backdrop -->
        <div
          v-if="showMobileUnitMenu"
          class="fixed inset-0 z-40"
          @click="showMobileUnitMenu = false"
        ></div>

        <!-- Mobile Unit Dropdown Menu -->
        <div
          v-if="showMobileUnitMenu"
          class="absolute left-0 mt-1.5 w-48 bg-white rounded-xl shadow-xl border border-zinc-200/90 p-1.5 z-50 text-xs animate-in fade-in zoom-in-95 duration-100"
        >
          <div class="px-2 py-1 text-[10px] font-bold text-zinc-400 uppercase tracking-wider">{{ localeStore.t('nav.selectEntity') }}</div>
          <button
            v-for="unit in store.unitOptions"
            :key="unit.code"
            @click="store.setBusinessUnit(unit.code as BusinessUnit); showMobileUnitMenu = false"
            class="w-full flex items-center justify-between px-2.5 py-2 rounded-lg font-semibold transition-colors cursor-pointer text-left"
            :class="store.businessUnit === unit.code ? 'bg-zinc-100 text-zinc-950 font-bold' : 'text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900'"
          >
            <div class="flex items-center gap-2">
              <span
                class="w-2 h-2 rounded-full"
                :class="
                  unit.code === 'KIDINEST' ? 'bg-blue-500' :
                  unit.code === 'SSC' ? 'bg-rose-500' :
                  unit.code === 'MINIMO' ? 'bg-amber-500' : 'bg-zinc-900'
                "
              ></span>
              <span>{{ unit.name }}</span>
            </div>
            <span v-if="store.businessUnit === unit.code" class="text-zinc-950 font-black">✓</span>
          </button>
        </div>
      </div>

      <!-- Desktop/Tablet Business Unit Segmented Tabs (Fully visible without clipping) -->
      <div class="hidden sm:flex items-center bg-zinc-100 p-0.5 rounded-xl border border-zinc-200/90 shrink-0">
        <button
          v-for="unit in store.unitOptions"
          :key="unit.code"
          @click="store.setBusinessUnit(unit.code as BusinessUnit)"
          :class="[
            'px-2.5 lg:px-3 py-1 sm:py-1.5 rounded-lg text-[11px] sm:text-xs font-bold transition-all duration-150 shrink-0 cursor-pointer',
            store.businessUnit === unit.code
              ? unit.code === 'KIDINEST'
                ? 'bg-blue-100 text-blue-950 border border-blue-300 shadow-2xs'
                : unit.code === 'SSC'
                ? 'bg-rose-100 text-rose-950 border border-rose-300 shadow-2xs'
                : unit.code === 'MINIMO'
                ? 'bg-amber-100 text-amber-950 border border-amber-300 shadow-2xs'
                : 'bg-zinc-950 text-white shadow-2xs'
              : 'text-zinc-600 hover:text-zinc-950'
          ]"
        >
          {{ unit.code === 'ALL' ? (localeStore.isEnglish ? 'Consolidated' : 'Konsolidasi') : unit.code }}
        </button>
      </div>

      <!-- Fiscal Period Selector (Hidden on screens below lg to avoid navbar crowding) -->
      <div class="relative hidden lg:flex items-center shrink-0">
        <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-zinc-400">
          <Calendar class="w-3.5 h-3.5" />
        </div>
        <select
          :value="store.fiscalMonth"
          @change="(e: any) => store.setFiscalMonth(Number(e.target.value))"
          class="pl-8 pr-7 py-1.5 bg-zinc-50 border border-zinc-200 text-zinc-900 text-xs rounded-xl font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none appearance-none cursor-pointer"
        >
          <option :value="11">Nov 2026 (Berjalan)</option>
          <option :value="10">Okt 2026</option>
          <option :value="9">Sep 2026</option>
          <option :value="0">YTD 2026</option>
          <option :value="3">Mar 2026</option>
          <option :value="2">Feb 2026</option>
          <option :value="1">Jan 2026</option>
        </select>
        <ChevronDown class="w-3.5 h-3.5 text-zinc-400 absolute right-2 pointer-events-none" />
      </div>
    </div>

    <!-- Right: Authenticated User Profile & Actions -->
    <div class="flex items-center gap-1.5 sm:gap-2 shrink-0">
      <!-- Settings Shortcut (Pengaturan Bahasa & Regional di Halaman Setting Saja) -->
      <router-link
        to="/accounting/settings"
        class="flex items-center gap-1.5 p-2 sm:px-3 sm:py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 hover:text-zinc-950 text-xs font-bold transition-colors shadow-2xs cursor-pointer shrink-0"
        :title="localeStore.t('nav.languageSettings')"
      >
        <Settings class="w-4 h-4 text-zinc-600 shrink-0" />
        <span class="hidden sm:inline">{{ localeStore.t('common.settings') }}</span>
      </router-link>

      <!-- Authenticated User Identity (XL Desktop only, mobile & tablet see profile in drawer) -->
      <div class="hidden xl:flex items-center gap-2 px-3 py-1.5 rounded-xl border border-zinc-200 bg-zinc-50 text-xs shrink-0">
        <ShieldCheck class="w-4 h-4 text-emerald-600 shrink-0" />
        <div class="text-left">
          <span class="block font-bold text-zinc-950 leading-none truncate max-w-[200px] 2xl:max-w-none">
            {{ authStore.user?.name || store.currentUserName }}
          </span>
          <span class="text-[9px] text-zinc-500 font-semibold tracking-wide uppercase">
            Finance & Accounting
          </span>
        </div>
      </div>

      <!-- Logout Button (Tablet & Desktop only, mobile accesses from drawer) -->
      <button
        @click="handleLogout"
        class="hidden sm:flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-rose-50 hover:border-rose-200 text-zinc-700 hover:text-rose-700 text-xs font-bold transition-colors cursor-pointer shadow-2xs shrink-0"
        :title="localeStore.t('nav.logout')"
      >
        <LogOut class="w-3.5 h-3.5 shrink-0" />
        <span class="hidden md:inline">{{ localeStore.t('nav.logout') }}</span>
      </button>

      <!-- Notifications Bell -->
      <div class="relative">
        <button
          @click="showNotifications = !showNotifications"
          class="relative p-2 rounded-xl text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 transition-colors cursor-pointer"
          :aria-label="localeStore.t('nav.notifications')"
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
          class="absolute right-0 mt-2 w-80 bg-white rounded-2xl shadow-xl border border-zinc-200/90 p-3.5 z-50 text-xs"
        >
          <div class="flex items-center justify-between pb-2 border-b border-zinc-100 mb-2">
            <h4 class="font-bold text-zinc-950">{{ localeStore.t('nav.notifications') }}</h4>
            <button
              @click="store.markAllNotificationsRead()"
              class="text-[11px] text-zinc-600 hover:text-zinc-950 font-semibold hover:underline"
            >
              {{ localeStore.t('nav.markAllRead') }}
            </button>
          </div>

          <div class="space-y-2 max-h-72 overflow-y-auto">
            <div
              v-for="notif in store.notifications"
              :key="notif.id"
              class="p-2.5 rounded-xl bg-zinc-50 hover:bg-zinc-100/80 transition-colors border border-zinc-100 space-y-1"
            >
              <div class="flex items-center justify-between">
                <span class="font-bold text-zinc-900">{{ notif.title }}</span>
                <span class="text-[10px] text-zinc-400">{{ notif.timestamp }}</span>
              </div>
              <p class="text-zinc-600 text-[11px] leading-relaxed">{{ notif.message }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>
