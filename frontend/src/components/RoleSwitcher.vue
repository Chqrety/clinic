<script setup lang="ts">
import { useClinicStore, type RoleType } from '../stores/clinic'
import {
  Smartphone,
  CreditCard,
  Stethoscope,
  Pill,
  Package,
  Receipt,
  ShieldCheck,
  Volume2,
  VolumeX,
} from 'lucide-vue-next'

const store = useClinicStore()

const roles: Array<{ key: RoleType; label: string; icon: any }> = [
  { key: 'patient', label: 'Pasien (HP)', icon: Smartphone },
  { key: 'cashier', label: 'Kasir & Antrian', icon: CreditCard },
  { key: 'doctor', label: 'Ruang Dokter', icon: Stethoscope },
  { key: 'pharmacy', label: 'Apotek (Racik)', icon: Pill },
  { key: 'inventory', label: 'Gudang Stok', icon: Package },
  { key: 'billing', label: 'Invoice', icon: Receipt },
  { key: 'satusehat', label: 'SATUSEHAT', icon: ShieldCheck },
]
</script>

<template>
  <header class="fixed top-4 left-1/2 -translate-x-1/2 z-50 w-full max-w-4xl px-4 pointer-events-none">
    <div
      class="pointer-events-auto flex items-center justify-between gap-2 p-1.5 bg-white/90 backdrop-blur-md rounded-full border border-black/10 shadow-[0_4px_24px_rgba(0,0,0,0.08)]"
    >
      <!-- Brand Pill -->
      <div class="hidden sm:flex items-center gap-2 pl-3 pr-2 text-xs font-bold tracking-tight text-zinc-900">
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
        <span>KLINIK SEHAT</span>
      </div>

      <!-- Navigation Roles -->
      <nav class="flex items-center gap-1 overflow-x-auto no-scrollbar py-0.5">
        <button
          v-for="r in roles"
          :key="r.key"
          @click="store.setRole(r.key)"
          :class="[
            'flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all duration-150 active:scale-95',
            store.currentRole === r.key
              ? 'bg-zinc-950 text-white shadow-sm'
              : 'text-zinc-600 hover:text-zinc-950 hover:bg-stone-100'
          ]"
        >
          <component :is="r.icon" class="w-3.5 h-3.5" />
          <span>{{ r.label }}</span>
        </button>
      </nav>

      <!-- Audio Mute Toggle for Demo -->
      <button
        @click="store.soundEnabled = !store.soundEnabled"
        :title="store.soundEnabled ? 'Suara Notifikasi Aktif' : 'Suara Dimatikan'"
        class="hidden md:flex items-center justify-center w-8 h-8 rounded-full text-zinc-500 hover:text-zinc-950 hover:bg-stone-100 transition-colors"
      >
        <Volume2 v-if="store.soundEnabled" class="w-4 h-4 text-emerald-600" />
        <VolumeX v-else class="w-4 h-4 text-zinc-400" />
      </button>
    </div>
  </header>
</template>
