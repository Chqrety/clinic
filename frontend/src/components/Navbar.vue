<script setup lang="ts">
import { ArrowRight, Stethoscope } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
}>()
</script>

<template>
  <header class="w-full max-w-7xl mx-auto px-6 sm:px-10 pt-8 pb-6 flex items-center justify-between">
    <!-- Brand Mark (Brandlos style) -->
    <button
      @click="emit('navigate', 'landing')"
      class="flex items-center gap-2.5 text-left group"
    >
      <div class="w-8 h-8 rounded-full bg-zinc-950 text-white flex items-center justify-center font-extrabold text-sm group-hover:scale-105 transition-transform">
        K
      </div>
      <div>
        <span class="font-extrabold text-lg tracking-tight text-zinc-950 block leading-tight">
          Klinik Sehat<span class="text-amber-500 font-black">.</span>
        </span>
        <span class="text-[10px] uppercase font-bold tracking-widest text-zinc-400 block -mt-0.5">
          Modern Healthcare
        </span>
      </div>
    </button>

    <!-- Nav Links & Staff CTA -->
    <div class="flex items-center gap-4 sm:gap-8">
      <nav class="hidden md:flex items-center gap-6 text-sm font-semibold text-zinc-600">
        <a href="#keunggulan" class="hover:text-zinc-950 transition-colors">Alur Pelayanan</a>
        <a href="#live-queue" class="hover:text-zinc-950 transition-colors">Antrian Live</a>
        <a href="#lokasi" class="hover:text-zinc-950 transition-colors">Lokasi & Kontak</a>
      </nav>

      <div class="flex items-center gap-2.5">
        <button
          @click="emit('navigate', 'kiosk')"
          class="hidden sm:inline-flex items-center gap-1.5 px-5 py-2.5 rounded-full text-xs font-bold text-zinc-900 bg-white border border-black/10 hover:border-black/20 hover:bg-stone-50 active:scale-95 transition-all shadow-sm"
        >
          Daftar Pasien
        </button>

        <button
          @click="emit('navigate', 'login')"
          class="inline-flex items-center gap-1.5 px-6 py-2.5 rounded-full text-xs font-extrabold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-95 transition-all shadow-sm"
        >
          <Stethoscope class="w-3.5 h-3.5" />
          <span>Portal Staf</span>
        </button>
      </div>
    </div>
  </header>
</template>
