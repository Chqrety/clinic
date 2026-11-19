<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { useAuthStore } from '@/stores/auth'
import { useLocaleStore } from '@/stores/locale'
import type { BusinessUnit } from '@/types/accounting'
import {
  Bell,
  Calendar,
  ChevronDown,
  LogOut,
  ShieldCheck,
  Stethoscope,
  Languages,
  Settings
} from 'lucide-vue-next'

const router = useRouter()
const store = useAccountingStore()
const authStore = useAuthStore()
const localeStore = useLocaleStore()

const showNotifications = ref(false)

async function handleLogout() {
  await authStore.logout()
  store.showToast('Berhasil logout dari sistem akuntansi', 'info')
  router.push('/login')
}
</script>

<template>
  <header
    class="h-16 bg-white border-b border-zinc-200/90 px-6 flex items-center justify-between sticky top-0 z-30 shadow-2xs"
  >
    <!-- Left: Global Entity (Unit Bisnis) & Fiscal Controls -->
    <div class="flex items-center gap-3">
      <!-- Business Unit Selector (Semantic Dual-Brand Tokens) -->
      <div class="flex items-center bg-zinc-100 p-0.5 rounded-xl border border-zinc-200/90">
        <button
          v-for="unit in store.unitOptions"
          :key="unit.code"
          @click="store.setBusinessUnit(unit.code as BusinessUnit)"
          :class="[
            'px-3 py-1.5 rounded-lg text-xs font-bold transition-all duration-150',
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
          {{ unit.badge }}: {{ unit.code === 'ALL' ? 'Konsolidasi' : unit.code }}
        </button>
      </div>

      <!-- Fiscal Period Selector -->
      <div class="relative flex items-center">
        <div class="absolute inset-y-0 left-0 pl-2.5 flex items-center pointer-events-none text-zinc-400">
          <Calendar class="w-3.5 h-3.5" />
        </div>
        <select
          :value="store.fiscalMonth"
          @change="(e: any) => store.setFiscalMonth(Number(e.target.value))"
          class="pl-8 pr-7 py-1.5 bg-zinc-50 border border-zinc-200 text-zinc-900 text-xs rounded-xl font-medium focus:ring-1 focus:ring-zinc-950 focus:outline-none appearance-none cursor-pointer"
        >
          <option :value="11">November 2026 (Aktif / Berjalan)</option>
          <option :value="10">Oktober 2026 (Lengkap)</option>
          <option :value="9">September 2026 (Soft Launch)</option>
          <option :value="0">Semua Periode (YTD 2026)</option>
          <option :value="3">Maret 2026</option>
          <option :value="2">Februari 2026</option>
          <option :value="1">Januari 2026</option>
        </select>
        <ChevronDown class="w-3.5 h-3.5 text-zinc-400 absolute right-2 pointer-events-none" />
      </div>
    </div>

    <!-- Right: Authenticated User Profile & Logout -->
    <div class="flex items-center gap-3">
      <!-- Shortcut to Clinic Admin Workstation -->
      <router-link
        to="/admin"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800 text-xs font-bold transition-colors shadow-2xs"
        title="Buka Workstation Medis Klinik"
      >
        <Stethoscope class="w-3.5 h-3.5 text-zinc-700" />
        <span>Workstation Medis</span>
      </router-link>

      <!-- Language Switcher Pill -->
      <button
        @click="localeStore.toggleLocale()"
        class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl border border-zinc-200 bg-zinc-50 hover:bg-zinc-100 text-zinc-800 text-xs font-bold transition shadow-2xs cursor-pointer"
        :title="localeStore.isIndonesian ? 'Ganti ke English (US)' : 'Switch to Bahasa Indonesia'"
      >
        <Languages class="w-3.5 h-3.5 text-blue-600" />
        <span class="font-mono text-[11px]">{{ localeStore.isIndonesian ? '🇮🇩 ID' : '🇬🇧 EN' }}</span>
      </button>

      <!-- Settings Shortcut -->
      <router-link
        to="/accounting/settings"
        class="p-2 rounded-xl text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 transition-colors"
        :title="localeStore.t('nav.languageSettings')"
      >
        <Settings class="w-4 h-4" />
      </router-link>

      <!-- Authenticated User Identity (Production Locked) -->
      <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl border border-zinc-200 bg-zinc-50 text-xs">
        <ShieldCheck class="w-4 h-4 text-emerald-600" />
        <div class="text-left">
          <span class="block font-bold text-zinc-950 leading-none">
            {{ authStore.user?.name || store.currentUserName }}
          </span>
          <span class="text-[10px] text-zinc-500 font-semibold tracking-wide uppercase">
            Finance & Accounting Manager
          </span>
        </div>
      </div>

      <!-- Logout Button -->
      <button
        @click="handleLogout"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-rose-50 hover:border-rose-200 text-zinc-700 hover:text-rose-700 text-xs font-bold transition-colors cursor-pointer shadow-2xs"
        title="Logout dari Sistem Akuntansi"
      >
        <LogOut class="w-3.5 h-3.5" />
        <span>Keluar</span>
      </button>

      <!-- Notifications Bell -->
      <div class="relative">
        <button
          @click="showNotifications = !showNotifications"
          class="relative p-2 rounded-xl text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100 transition-colors"
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
          class="absolute right-0 mt-2 w-80 bg-white rounded-2xl shadow-xl border border-zinc-200/90 p-3.5 z-50 text-xs"
        >
          <div class="flex items-center justify-between pb-2 border-b border-zinc-100 mb-2">
            <h4 class="font-bold text-zinc-950">Pemberitahuan Sistem</h4>
            <button
              @click="store.markAllNotificationsRead()"
              class="text-[11px] text-zinc-600 hover:text-zinc-950 font-semibold hover:underline"
            >
              Tandai Dibaca
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
