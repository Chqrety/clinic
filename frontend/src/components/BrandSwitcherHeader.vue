<script setup lang="ts">
import { ref } from 'vue'
import {
  ArrowRight,
  ArrowUpRight,
  Menu,
  X,
  MapPin
} from 'lucide-vue-next'

const props = defineProps<{
  currentBrand: 'kidinest' | 'ssc'
}>()

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login'): void
}>()

const mobileMenuOpen = ref(false)

function scrollToSection(id: string) {
  mobileMenuOpen.value = false
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<template>
  <header
    class="sticky top-0 z-40 border-b backdrop-blur-md transition-colors select-none"
    :class="
      currentBrand === 'kidinest'
        ? 'bg-[#FAF7F2]/95 border-[#E5A93C]/20'
        : 'bg-[#FAF5F5]/95 border-[#9E606D]/20'
    "
  >
    <!-- Top Micro-Ribbon: Hidden di mobile, ramping di desktop (sm:block) -->
    <div
      class="hidden sm:block w-full py-1 text-[11px] font-medium border-b transition-colors"
      :class="
        currentBrand === 'kidinest'
          ? 'bg-[#E5A93C]/10 border-[#E5A93C]/20 text-[#0C3B36]/80'
          : 'bg-[#9E606D]/10 border-[#9E606D]/15 text-[#2D1E24]/80'
      "
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex items-center justify-between">
        <div class="flex items-center gap-1.5">
          <MapPin class="w-3 h-3 shrink-0 opacity-70" />
          <span>Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang</span>
        </div>
        <div class="flex items-center gap-3 text-[11px] opacity-75">
          <span>Poli Spesialis Buka Setiap Hari</span>
          <span>•</span>
          <span>WA: {{ currentBrand === 'kidinest' ? '0812-3456-7890' : '0811-9876-5432' }}</span>
        </div>
      </div>
    </div>

    <!-- Main Navigation Bar (Ramping 56px mobile, 64px desktop) -->
    <div class="max-w-7xl mx-auto px-4 sm:px-8">
      <div class="flex items-center justify-between h-14 sm:h-16">

        <!-- Left: Logo Bersih & Elegan (Bebas Box AI & Subtitle Panjang) -->
        <button
          @click="emit('navigate', currentBrand)"
          class="flex items-center gap-2.5 text-left group focus:outline-hidden"
        >
          <img
            v-if="currentBrand === 'kidinest'"
            src="/assets/kidinest/logo-maskot.png"
            alt="Kidinest"
            class="w-7 h-7 sm:w-8 sm:h-8 object-contain transition-transform group-hover:scale-105"
          />
          <img
            v-else
            src="/assets/ssc/logo-ssc.png"
            alt="SSC"
            class="w-7 h-7 sm:w-8 sm:h-8 object-contain transition-transform group-hover:scale-105"
          />

          <div class="flex items-baseline gap-1.5 leading-none">
            <span
              class="font-extrabold text-base sm:text-lg tracking-tight"
              :class="currentBrand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#2D1E24]'"
            >
              {{ currentBrand === 'kidinest' ? 'kidinest' : 'Semarang Skin Center' }}
            </span>
            <span
              v-if="currentBrand === 'kidinest'"
              class="text-[11px] font-semibold text-[#0C3B36]/60 tracking-normal hidden sm:inline"
            >
              Pediatric Care
            </span>
          </div>
        </button>

        <!-- Center: Clean Navigation Links -->
        <nav
          class="hidden md:flex items-center gap-7 text-xs font-semibold"
          :class="currentBrand === 'kidinest' ? 'text-[#0C3B36]/75' : 'text-[#2D1E24]/75'"
        >
          <button
            @click="scrollToSection('live-antrian')"
            class="hover:opacity-100 transition-opacity flex items-center gap-1.5 py-1"
          >
            <span
              class="w-1.5 h-1.5 rounded-full animate-pulse"
              :class="currentBrand === 'kidinest' ? 'bg-[#E5A93C]' : 'bg-[#9E606D]'"
            ></span>
            <span>Antrian Live</span>
          </button>
          <button
            @click="scrollToSection('layanan')"
            class="hover:opacity-100 transition-opacity py-1"
          >
            Layanan
          </button>
          <button
            @click="scrollToSection('jadwal-dokter')"
            class="hover:opacity-100 transition-opacity py-1"
          >
            Dokter Spesialis
          </button>
          <button
            @click="scrollToSection('lokasi')"
            class="hover:opacity-100 transition-opacity py-1"
          >
            Lokasi & Kontak
          </button>
        </nav>

        <!-- Right: Link Santai Sister Clinic + Single Primary CTA -->
        <div class="hidden sm:flex items-center gap-4">
          <!-- Link Santai Sister Clinic (Bukan Switcher Tab) -->
          <button
            v-if="currentBrand === 'kidinest'"
            @click="emit('navigate', 'ssc')"
            class="text-xs font-semibold text-[#0C3B36]/70 hover:text-[#0C3B36] flex items-center gap-1 transition-colors hover:underline underline-offset-4"
            title="Kunjungi Klinik Kulit Semarang Skin Center"
          >
            <span>Layanan Kulit (SSC)</span>
            <ArrowUpRight class="w-3.5 h-3.5 opacity-60" />
          </button>

          <button
            v-else
            @click="emit('navigate', 'kidinest')"
            class="text-xs font-semibold text-[#2D1E24]/70 hover:text-[#2D1E24] flex items-center gap-1 transition-colors hover:underline underline-offset-4"
            title="Kunjungi Kidinest Pediatric Care"
          >
            <span>Poli Anak (Kidinest)</span>
            <ArrowUpRight class="w-3.5 h-3.5 opacity-60" />
          </button>

          <!-- Single Primary CTA Button -->
          <button
            @click="emit('navigate', 'kiosk')"
            class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider text-white transition-all shadow-xs active:scale-95"
            :class="
              currentBrand === 'kidinest'
                ? 'bg-[#0C3B36] hover:bg-[#134E48]'
                : 'bg-[#9E606D] hover:bg-[#8B505D]'
            "
          >
            <span>Ambil Antrian</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Mobile Menu Toggle Button -->
        <div class="flex items-center gap-2 sm:hidden">
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="p-2 rounded-lg text-zinc-700 hover:bg-black/5"
            aria-label="Toggle Navigation"
          >
            <Menu v-if="!mobileMenuOpen" class="w-5 h-5" />
            <X v-else class="w-5 h-5" />
          </button>
        </div>

      </div>
    </div>

    <!-- Mobile Dropdown Menu -->
    <div
      v-if="mobileMenuOpen"
      class="sm:hidden border-t px-4 pt-3 pb-5 space-y-3 bg-white/98 shadow-lg animate-in fade-in duration-150"
      :class="currentBrand === 'kidinest' ? 'border-[#E5A93C]/20' : 'border-[#9E606D]/20'"
    >
      <nav class="space-y-1 text-sm font-semibold text-zinc-800">
        <button
          @click="scrollToSection('live-antrian')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-black/5 flex items-center justify-between"
        >
          <span>Antrian Live</span>
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
        </button>
        <button
          @click="scrollToSection('layanan')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-black/5"
        >
          Layanan
        </button>
        <button
          @click="scrollToSection('jadwal-dokter')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-black/5"
        >
          Dokter Spesialis
        </button>
        <button
          @click="scrollToSection('lokasi')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-black/5"
        >
          Lokasi & Kontak
        </button>
      </nav>

      <!-- Mobile Sister Clinic Link -->
      <div class="pt-2 border-t border-zinc-100 flex items-center justify-between text-xs px-3">
        <span class="text-zinc-500 font-medium">Fasilitas Lain:</span>
        <button
          v-if="currentBrand === 'kidinest'"
          @click="emit('navigate', 'ssc'); mobileMenuOpen = false"
          class="font-bold text-[#9E606D] flex items-center gap-1 hover:underline"
        >
          <span>Layanan Kulit (SSC)</span>
          <ArrowUpRight class="w-3.5 h-3.5" />
        </button>
        <button
          v-else
          @click="emit('navigate', 'kidinest'); mobileMenuOpen = false"
          class="font-bold text-[#0C3B36] flex items-center gap-1 hover:underline"
        >
          <span>Poli Anak (Kidinest)</span>
          <ArrowUpRight class="w-3.5 h-3.5" />
        </button>
      </div>

      <div class="pt-1">
        <button
          @click="emit('navigate', 'kiosk'); mobileMenuOpen = false"
          class="w-full py-2.5 rounded-xl font-bold text-xs uppercase tracking-wider text-center text-white shadow-xs"
          :class="currentBrand === 'kidinest' ? 'bg-[#0C3B36]' : 'bg-[#9E606D]'"
        >
          Ambil Antrian Online →
        </button>
      </div>
    </div>
  </header>
</template>
