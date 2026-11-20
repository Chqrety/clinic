<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue'
import { useRouter } from 'vue-router'
import { useClinicStore, type RoleType } from '../stores/clinic'
import { useAuthStore } from '../stores/auth'
import {
  ShieldCheck,
  Stethoscope,
  Pill,
  CreditCard,
  Package,
  Receipt,
  Tv,
  LogOut,
  Volume2,
  VolumeX,
  ChevronDown,
  Building2,
  User,
  Sparkles,
  ExternalLink,
  TrendingUp
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login' | 'monitor'): void
}>()

const router = useRouter()
const store = useClinicStore()
const auth = useAuthStore()

const isRoleMenuOpen = ref(false)

// Sesuai Master Prompt Bagian 2 (Struktur Entitas):
// 1. CV Semarang Medical Group: Kidinest, SSC, GRO
// 2. CV Semarang Medical Clinic: Farmasi, Apoteker, Gudang Logistik
const activeEntityInfo = computed(() => {
  if (['pharmacy', 'inventory'].includes(store.currentRole)) {
    return {
      entity: 'CV Semarang Medical Clinic',
      unitDesc: 'Apotek & Logistik Medis',
      color: 'text-purple-400 bg-purple-500/10 border-purple-500/30'
    }
  }
  return {
    entity: 'CV Semarang Medical Group',
    unitDesc: 'Kidinest & SSC Care',
    color: 'text-emerald-400 bg-emerald-500/10 border-emerald-500/30'
  }
})


interface MenuRole {
  key: RoleType
  label: string
  sub: string
  icon: any
  badge?: string
  allowedRoles: string[]
}

const allRoles: MenuRole[] = [
  { key: 'cashier', label: 'Kasir & Pendaftaran', sub: 'Meja Depan & GRO', icon: CreditCard, allowedRoles: ['SUPERADMIN_OWNER', 'GRO_CASHIER', 'FINANCE'] },
  { key: 'doctor', label: 'Ruang Dokter', sub: 'E-Prescription & ICD-10', icon: Stethoscope, allowedRoles: ['SUPERADMIN_OWNER', 'DOCTOR'] },
  { key: 'pharmacy', label: 'Farmasi & Apotek', sub: 'Dispensing Resep', icon: Pill, allowedRoles: ['SUPERADMIN_OWNER', 'PHARMACIST'] },
  { key: 'inventory', label: 'Gudang & Logistik', sub: 'Stok Opname & Batch', icon: Package, allowedRoles: ['SUPERADMIN_OWNER', 'PHARMACIST'] },
  { key: 'billing', label: 'Billing & Tutup Shift', sub: 'Rekap Finansial', icon: Receipt, allowedRoles: ['SUPERADMIN_OWNER', 'GRO_CASHIER', 'FINANCE'] },
  { key: 'owner', label: 'Manajemen Staf', sub: 'Hak Akses Atasan', icon: ShieldCheck, allowedRoles: ['SUPERADMIN_OWNER', 'OWNER'] },
  { key: 'satusehat', label: 'SATUSEHAT Bridge', sub: 'Kemenkes RME', icon: Sparkles, allowedRoles: ['SUPERADMIN_OWNER'] },
]

const isSuperAdmin = computed(() => {
  const role = (auth.user?.role || 'SUPERADMIN_OWNER').toUpperCase()
  return role === 'SUPERADMIN_OWNER' || role === 'OWNER' || role === 'FINANCE' || auth.isOwner || !auth.user
})

// Hak akses khusus Dashboard Eksekutif Manajemen & Omzet
const canAccessManagement = computed(() => {
  const role = (auth.user?.role || '').toUpperCase()
  return isSuperAdmin.value || ['OWNER', 'GRO_CASHIER', 'FINANCE'].includes(role)
})

// Hak akses khusus ERP Akuntansi
const canAccessAccounting = computed(() => {
  const role = (auth.user?.role || '').toUpperCase()
  return isSuperAdmin.value || ['OWNER', 'FINANCE'].includes(role)
})

// Superadmin sees all 7 modules; non-owners see only their designated modules
const visibleRoles = computed(() => {
  if (isSuperAdmin.value) {
    return allRoles
  }
  const userRole = (auth.user?.role || '').toUpperCase()
  return allRoles.filter(r => r.allowedRoles.includes(userRole))
})

// Auto-align active role with user permissions
watchEffect(() => {
  if (visibleRoles.value.length > 0 && !visibleRoles.value.some(r => r.key === store.currentRole)) {
    store.setRole(visibleRoles.value[0].key)
  }
})

function getRoleDisplayName(role: RoleType): string {
  const found = allRoles.find(r => r.key === role)
  return found ? found.label : 'Staf'
}

function handleSelectRole(role: RoleType) {
  store.setRole(role)
  isRoleMenuOpen.value = false
}

async function handleLogout() {
  window.location.hash = ''
  await auth.logout()
  store.setRole('patient')
  emit('navigate', 'login')
  router.replace('/login-kidinest')
}
</script>

<template>
  <aside class="w-64 bg-zinc-950 text-zinc-100 flex flex-col justify-between border-r border-zinc-800 shrink-0 select-none">
    
    <!-- Top Part: Clinic Brand & Staff Profile -->
    <div>
      <!-- Brand Header -->
      <div class="p-4 border-b border-zinc-800/80 flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-lg bg-white text-zinc-950 flex items-center justify-center font-black text-sm shadow-xs">
            SMG
          </div>
          <div>
            <div class="font-extrabold text-xs tracking-tight text-white leading-tight">
              {{ activeEntityInfo.entity }}
            </div>
            <div class="text-[10px] font-medium flex items-center gap-1.5 mt-0.5">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
              <span class="text-zinc-400">{{ activeEntityInfo.unitDesc }}</span>
            </div>
          </div>
        </div>

        <button
          @click="emit('navigate', 'kidinest')"
          title="Lihat Landing Pasien"
          class="p-1.5 rounded-lg text-zinc-400 hover:text-white hover:bg-zinc-800 transition-colors"
        >
          <ExternalLink class="w-3.5 h-3.5" />
        </button>
      </div>

      <!-- Staff Profile & Role Switcher Dropdown (Superadmin gets switcher, others get clean profile) -->
      <div class="p-3 border-b border-zinc-800/80 bg-zinc-900/60">
        <div class="relative">
          <!-- Superadmin Interactive Switcher -->
          <button
            v-if="isSuperAdmin"
            @click="isRoleMenuOpen = !isRoleMenuOpen"
            class="w-full p-2 rounded-xl bg-zinc-800/80 hover:bg-zinc-800 text-left transition-all border border-zinc-700/60 flex items-center justify-between group"
          >
            <div class="flex items-center gap-2.5 overflow-hidden">
              <div class="w-7 h-7 rounded-lg bg-zinc-700 text-zinc-200 flex items-center justify-center shrink-0">
                <User class="w-4 h-4" />
              </div>
              <div class="truncate">
                <div class="text-xs font-bold text-white truncate">
                  {{ auth.user?.name || 'dr. Ratna Paramita, M.Kes' }}
                </div>
                <div class="text-[10px] text-amber-400 font-semibold truncate flex items-center gap-1">
                  <span>Atasan &bull; {{ getRoleDisplayName(store.currentRole) }}</span>
                </div>
              </div>
            </div>
            <ChevronDown class="w-3.5 h-3.5 text-zinc-400 group-hover:text-white shrink-0 transition-transform" :class="isRoleMenuOpen ? 'rotate-180' : ''" />
          </button>

          <!-- Non-Superadmin Fixed Profile (Clean & Professional) -->
          <div
            v-else
            class="w-full p-2 rounded-xl bg-zinc-800/80 text-left border border-zinc-700/60 flex items-center justify-between"
          >
            <div class="flex items-center gap-2.5 overflow-hidden">
              <div class="w-7 h-7 rounded-lg bg-zinc-700 text-zinc-200 flex items-center justify-center shrink-0">
                <User class="w-4 h-4" />
              </div>
              <div class="truncate">
                <div class="text-xs font-bold text-white truncate">
                  {{ auth.user?.name || 'Staf Medis' }}
                </div>
                <div class="text-[10px] text-emerald-400 font-semibold truncate flex items-center gap-1">
                  <span>{{ auth.user?.role || 'STAF' }} &bull; {{ auth.user?.clinic_unit || 'UNIT' }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Dropdown Switcher Menu (Superadmin Only) -->
          <div
            v-if="isRoleMenuOpen && isSuperAdmin"
            class="absolute top-full left-0 right-0 mt-1.5 p-1.5 bg-zinc-900 border border-zinc-700 rounded-xl shadow-2xl z-50 space-y-1 animate-in fade-in duration-100"
          >
            <div class="px-2 py-1 text-[9px] font-black uppercase tracking-wider text-zinc-400">
              Ganti Peran Staf (Demo Superadmin)
            </div>
            <button
              v-for="r in allRoles"
              :key="r.key"
              @click="handleSelectRole(r.key)"
              class="w-full px-2.5 py-1.5 rounded-lg text-xs font-bold text-left flex items-center justify-between transition-colors"
              :class="store.currentRole === r.key ? 'bg-white text-zinc-950 font-black' : 'text-zinc-300 hover:bg-zinc-800 hover:text-white'"
            >
              <div class="flex items-center gap-2">
                <component :is="r.icon" class="w-3.5 h-3.5" />
                <span>{{ r.label }}</span>
              </div>
              <span v-if="store.currentRole === r.key" class="text-[10px]">●</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Nav Workstation Menu -->
      <nav class="p-3 space-y-1 overflow-y-auto max-h-[calc(100vh-280px)]">
        <div class="px-3 pt-2 pb-1 text-[10px] font-black uppercase tracking-wider text-zinc-400">
          Modul Pelayanan
        </div>

        <button
          v-for="item in visibleRoles"
          :key="item.key"
          @click="store.setRole(item.key)"
          class="w-full px-3 py-2.5 rounded-xl text-left flex items-center gap-3 transition-all text-xs font-bold group"
          :class="
            store.currentRole === item.key
              ? 'bg-zinc-100 text-zinc-950 font-extrabold shadow-sm'
              : 'text-zinc-300 hover:bg-zinc-900 hover:text-white'
          "
        >
          <component
            :is="item.icon"
            class="w-4 h-4 shrink-0 transition-transform group-hover:scale-105"
            :class="store.currentRole === item.key ? 'text-zinc-950' : 'text-zinc-400'"
          />
          <div class="flex-1 truncate">
            <div class="truncate">{{ item.label }}</div>
            <div
              class="text-[10px] font-medium truncate"
              :class="store.currentRole === item.key ? 'text-zinc-600' : 'text-zinc-400'"
            >
              {{ item.sub }}
            </div>
          </div>
        </button>

        <!-- Section: Portal Eksekutif & Finansial (Akses Terpisah) -->
        <div v-if="canAccessManagement || canAccessAccounting" class="pt-3 pb-1">
          <div class="px-3 pt-2 pb-1 text-[10px] font-black uppercase tracking-wider text-zinc-500">
            Portal Khusus
          </div>

          <!-- Standalone Dashboard: Manajemen & Omzet -->
          <router-link
            v-if="canAccessManagement"
            to="/management"
            class="w-full px-3 py-2.5 rounded-xl text-left flex items-center gap-3 transition-all text-xs font-bold text-emerald-300 hover:text-white bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 group mt-1"
          >
            <TrendingUp class="w-4 h-4 shrink-0 text-emerald-400 group-hover:scale-105 transition-transform" />
            <div class="flex-1 truncate">
              <div class="truncate">Manajemen & Omzet</div>
              <div class="text-[10px] font-medium text-emerald-400/80 truncate">
                Laporan Penjualan Harian
              </div>
            </div>
            <ExternalLink class="w-3.5 h-3.5 text-emerald-400/80" />
          </router-link>
        </div>
      </nav>
    </div>

    <!-- Bottom Controls: Sound Alert, TV Monitor, Logout -->
    <div class="p-3 border-t border-zinc-800 bg-zinc-950 space-y-1.5">
      <!-- Quick Action: TV Monitor -->
      <button
        @click="emit('navigate', 'monitor')"
        class="w-full px-3 py-2 rounded-xl text-xs font-bold text-zinc-300 hover:text-white hover:bg-zinc-900 flex items-center justify-between transition-colors border border-zinc-800/80"
      >
        <div class="flex items-center gap-2">
          <Tv class="w-4 h-4 text-sky-400" />
          <span>Monitor TV Antrian</span>
        </div>
        <span class="text-[10px] text-zinc-400">Layar TV</span>
      </button>

      <!-- Toggle Sound Chime -->
      <button
        @click="store.soundEnabled = !store.soundEnabled"
        class="w-full px-3 py-2 rounded-xl text-xs font-bold text-zinc-300 hover:text-white hover:bg-zinc-900 flex items-center justify-between transition-colors"
      >
        <div class="flex items-center gap-2">
          <Volume2 v-if="store.soundEnabled" class="w-4 h-4 text-emerald-400" />
          <VolumeX v-else class="w-4 h-4 text-zinc-400" />
          <span>Suara Lonceng</span>
        </div>
        <span
          class="text-[10px] font-black uppercase px-2 py-0.5 rounded-md"
          :class="store.soundEnabled ? 'bg-emerald-950 text-emerald-300' : 'bg-zinc-800 text-zinc-400'"
        >
          {{ store.soundEnabled ? 'Aktif' : 'Mute' }}
        </span>
      </button>

      <!-- Logout Button -->
      <button
        @click="handleLogout"
        class="w-full px-3 py-2 rounded-xl text-xs font-bold text-rose-400 hover:text-rose-300 hover:bg-rose-950/40 flex items-center gap-2 transition-colors"
      >
        <LogOut class="w-4 h-4" />
        <span>Keluar Sistem</span>
      </button>
    </div>

  </aside>
</template>
