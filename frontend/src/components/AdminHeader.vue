<script setup lang="ts">
import { useClinicStore, type RoleType } from '../stores/clinic'
import { Stethoscope, Pill, CreditCard, Package, Receipt, ShieldCheck, LogOut, Volume2, VolumeX } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
}>()

const store = useClinicStore()

const roles: Array<{ key: RoleType; label: string; icon: any }> = [
  { key: 'doctor', label: 'Dokter', icon: Stethoscope },
  { key: 'pharmacy', label: 'Apotek', icon: Pill },
  { key: 'cashier', label: 'Kasir', icon: CreditCard },
  { key: 'inventory', label: 'Gudang', icon: Package },
  { key: 'billing', label: 'Billing', icon: Receipt },
  { key: 'satusehat', label: 'SATUSEHAT', icon: ShieldCheck },
]

function getStaffTitle(role: RoleType) {
  switch (role) {
    case 'doctor':
      return 'dr. Hendra Pratama, Sp.PD (Poli Umum)'
    case 'pharmacy':
      return 'apt. Annisa Fitri, S.Farm (Kepala Farmasi)'
    case 'cashier':
      return 'Budi Setiawan (Loket Pendaftaran)'
    case 'inventory':
      return 'Siti Nurhaliza (Gudang & Logistik)'
    case 'billing':
      return 'Keuangan & Akuntansi Klinik'
    case 'satusehat':
      return 'Tim Integrasi DTO Kemenkes'
    default:
      return 'Staf Klinik'
  }
}
</script>

<template>
  <header class="w-full bg-white border-b border-black/10 sticky top-0 z-40">
    <div class="max-w-7xl mx-auto px-6 sm:px-10 py-3 flex items-center justify-between gap-4">

      <!-- Left: Logo & Role Title -->
      <div class="flex items-center gap-3">
        <button
          @click="emit('navigate', 'landing')"
          class="font-extrabold text-base tracking-tight text-zinc-950 flex items-center gap-1.5"
        >
          <span class="w-6 h-6 rounded-full bg-zinc-950 text-white flex items-center justify-center text-xs font-black">K</span>
          <span>Klinik Sehat</span>
        </button>

        <span class="text-zinc-300 hidden sm:inline">/</span>

        <span class="text-xs font-bold text-zinc-700 bg-stone-100 px-3 py-1 rounded-full hidden sm:inline-block">
          {{ getStaffTitle(store.currentRole) }}
        </span>
      </div>

      <!-- Center: Role Tabs -->
      <nav class="flex items-center gap-1 overflow-x-auto no-scrollbar py-1">
        <button
          v-for="r in roles"
          :key="r.key"
          @click="store.setRole(r.key)"
          :class="[
            'flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold whitespace-nowrap transition-all duration-150 active:scale-95',
            store.currentRole === r.key
              ? 'bg-zinc-950 text-white shadow-sm'
              : 'text-zinc-600 hover:text-zinc-950 hover:bg-stone-100'
          ]"
        >
          <component :is="r.icon" class="w-3.5 h-3.5" />
          <span>{{ r.label }}</span>
        </button>
      </nav>

      <!-- Right: Controls & Logout -->
      <div class="flex items-center gap-2">
        <button
          @click="store.soundEnabled = !store.soundEnabled"
          :title="store.soundEnabled ? 'Suara Lonceng Aktif' : 'Suara Dimatikan'"
          class="p-2 rounded-full hover:bg-stone-100 text-zinc-500 hover:text-zinc-950 transition-colors"
        >
          <Volume2 v-if="store.soundEnabled" class="w-4 h-4 text-emerald-600" />
          <VolumeX v-else class="w-4 h-4 text-zinc-400" />
        </button>

        <button
          @click="emit('navigate', 'landing')"
          title="Keluar ke Halaman Depan"
          class="flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold text-rose-700 bg-rose-50 hover:bg-rose-100 transition-colors"
        >
          <LogOut class="w-3.5 h-3.5" />
          <span class="hidden sm:inline">Keluar</span>
        </button>
      </div>

    </div>
  </header>
</template>
