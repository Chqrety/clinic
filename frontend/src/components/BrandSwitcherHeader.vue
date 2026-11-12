<script setup lang="ts">
import { ref } from 'vue'
import {
  ArrowRight,
  Menu,
  X,
  Clock,
  Phone,
  MapPin,
  Sparkles
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
    class="sticky top-0 z-40 border-b backdrop-blur-md transition-colors"
    :class="
      currentBrand === 'kidinest'
        ? 'bg-[#FAF7F2]/95 border-[#E5A93C]/20'
        : 'bg-[#FAF5F5]/95 border-[#9E606D]/20'
    "
  >
    <!-- Top ribbon bar: Single Clinical Group Notice -->
    <div
      class="w-full py-1.5 text-[11px] font-semibold border-b transition-colors"
      :class="
        currentBrand === 'kidinest'
          ? 'bg-[#E5A93C]/10 border-[#E5A93C]/20 text-[#0C3B36]'
          : 'bg-[#9E606D]/10 border-[#9E606D]/15 text-[#2D1E24]'
      "
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <MapPin class="w-3.5 h-3.5 shrink-0 opacity-70" />
          <span>Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang</span>
        </div>
        <div class="hidden sm:flex items-center gap-4 text-[11px] opacity-80">
          <span>Poli Spesialis Buka Setiap Hari</span>
          <span>•</span>
          <span>WA: 0812-3456-7890</span>
        </div>
      </div>
    </div>

    <!-- Main Navigation Bar -->
    <div class="max-w-7xl mx-auto px-4 sm:px-8">
      <div class="flex items-center justify-between h-20">

        <!-- Left: Logo Asli Brand -->
        <button
          @click="emit('navigate', currentBrand)"
          class="flex items-center gap-3.5 text-left group focus:outline-hidden"
        >
          <!-- Logo Image Container -->
          <div
            class="w-11 h-11 rounded-2xl flex items-center justify-center p-1 transition-transform group-hover:scale-105 shadow-2xs"
            :class="
              currentBrand === 'kidinest'
                ? 'bg-[#E5A93C]/15 border border-[#E5A93C]/30'
                : 'bg-[#9E606D] shadow-[#9E606D]/20'
            "
          >
            <img
              v-if="currentBrand === 'kidinest'"
              src="/assets/kidinest/logo-maskot.png"
              alt="Kidinest Mascot"
              class="w-8 h-8 object-contain"
            />
            <img
              v-else
              src="/assets/ssc/logo-ssc.png"
              alt="SSC Rosette Logo"
              class="w-9 h-9 object-contain"
            />
          </div>

          <!-- Brand Title & Tag -->
          <div>
            <div class="flex items-center gap-2">
              <span
                class="font-extrabold text-lg sm:text-xl tracking-tight"
                :class="currentBrand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#2D1E24]'"
              >
                {{ currentBrand === 'kidinest' ? 'kidinest' : 'SSC' }}
              </span>
              <span
                class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded-md"
                :class="
                  currentBrand === 'kidinest'
                    ? 'bg-[#E5A93C]/20 text-[#0C3B36]'
                    : 'bg-[#9E606D]/15 text-[#9E606D]'
                "
              >
                {{ currentBrand === 'kidinest' ? 'Pediatric Care' : 'Semarang Skin Center' }}
              </span>
            </div>
            <p
              class="text-xs font-medium hidden sm:block"
              :class="currentBrand === 'kidinest' ? 'text-[#0C3B36]/70' : 'text-[#2D1E24]/70'"
            >
              {{
                currentBrand === 'kidinest'
                  ? 'Pediatric • Vaccination • Tumbuh Kembang'
                  : 'Dermatology & Aesthetic Center'
              }}
            </p>
          </div>
        </button>

        <!-- Center: Hairline Navigation Links -->
        <nav
          class="hidden lg:flex items-center gap-8 text-sm font-semibold"
          :class="currentBrand === 'kidinest' ? 'text-[#0C3B36]/80' : 'text-[#2D1E24]/80'"
        >
          <button
            @click="scrollToSection('live-antrian')"
            class="hover:opacity-100 transition-opacity flex items-center gap-2"
          >
            <span
              class="w-2 h-2 rounded-full animate-pulse"
              :class="currentBrand === 'kidinest' ? 'bg-[#E5A93C]' : 'bg-[#9E606D]'"
            ></span>
            <span>Antrian Live</span>
          </button>
          <button
            @click="scrollToSection('layanan')"
            class="hover:opacity-100 transition-opacity"
          >
            {{ currentBrand === 'kidinest' ? 'Layanan & Edukasi' : 'Menu Perawatan' }}
          </button>
          <button
            @click="scrollToSection('jadwal-dokter')"
            class="hover:opacity-100 transition-opacity"
          >
            Dokter Spesialis
          </button>
          <button
            @click="scrollToSection('lokasi')"
            class="hover:opacity-100 transition-opacity"
          >
            Lokasi & Kontak
          </button>
        </nav>

        <!-- Right: Switcher Pill & CTA Button -->
        <div class="hidden sm:flex items-center gap-3">
          <!-- Switcher to Sibling Unit -->
          <button
            v-if="currentBrand === 'kidinest'"
            @click="emit('navigate', 'ssc')"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full border border-[#9E606D]/30 bg-white/70 hover:bg-white text-xs font-semibold text-[#9E606D] transition-all shadow-2xs"
            title="Kunjungi Semarang Skin Center"
          >
            <span>Ke Semarang Skin Center</span>
            <ArrowRight class="w-3 h-3 text-[#9E606D]" />
          </button>

          <button
            v-else
            @click="emit('navigate', 'kidinest')"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full border border-[#E5A93C]/40 bg-white/70 hover:bg-white text-xs font-semibold text-[#0C3B36] transition-all shadow-2xs"
            title="Kunjungi Kidinest Pediatric Care"
          >
            <span>Ke Kidinest Poli Anak</span>
            <ArrowRight class="w-3 h-3 text-[#E5A93C]" />
          </button>

          <!-- Primary CTA Button -->
          <button
            @click="emit('navigate', 'kiosk')"
            class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-bold uppercase tracking-wider text-white transition-all shadow-sm active:scale-95"
            :class="
              currentBrand === 'kidinest'
                ? 'bg-[#0C3B36] hover:bg-[#134E48]'
                : 'bg-[#9E606D] hover:bg-[#8B505D]'
            "
          >
            <span>Daftar Antrian</span>
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
            <Menu v-if="!mobileMenuOpen" class="w-6 h-6" />
            <X v-else class="w-6 h-6" />
          </button>
        </div>

      </div>
    </div>

    <!-- Mobile Dropdown Menu -->
    <div
      v-if="mobileMenuOpen"
      class="sm:hidden border-t px-4 pt-4 pb-6 space-y-4 bg-white/98 shadow-xl animate-in fade-in duration-150"
      :class="currentBrand === 'kidinest' ? 'border-[#E5A93C]/20' : 'border-[#9E606D]/20'"
    >
      <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200 flex items-center justify-between">
        <span class="text-xs font-bold text-zinc-600">Pindah Poli:</span>
        <button
          v-if="currentBrand === 'kidinest'"
          @click="emit('navigate', 'ssc'); mobileMenuOpen = false"
          class="text-xs font-bold text-[#9E606D] flex items-center gap-1"
        >
          <span>Semarang Skin Center →</span>
        </button>
        <button
          v-else
          @click="emit('navigate', 'kidinest'); mobileMenuOpen = false"
          class="text-xs font-bold text-[#0C3B36] flex items-center gap-1"
        >
          <span>Kidinest Poli Anak →</span>
        </button>
      </div>

      <nav class="space-y-2 text-sm font-semibold text-zinc-800">
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
          {{ currentBrand === 'kidinest' ? 'Layanan & Edukasi' : 'Menu Perawatan' }}
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

      <div class="pt-2">
        <button
          @click="emit('navigate', 'kiosk'); mobileMenuOpen = false"
          class="w-full py-3 rounded-xl font-bold text-xs uppercase tracking-wider text-center text-white shadow-md"
          :class="currentBrand === 'kidinest' ? 'bg-[#0C3B36]' : 'bg-[#9E606D]'"
        >
          Ambil Nomor Antrian Online →
        </button>
      </div>
    </div>
  </header>
</template>
