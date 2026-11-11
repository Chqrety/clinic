<script setup lang="ts">
import { ref } from 'vue'
import {
  Baby,
  Sparkles,
  ArrowRight,
  Menu,
  X,
  Clock,
  Phone,
  ShieldCheck,
  Stethoscope
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
        ? 'bg-[#FBFBF9]/95 border-emerald-950/10'
        : 'bg-[#FAFAFA]/95 border-zinc-200'
    "
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-8">
      <div class="flex items-center justify-between h-20">

        <!-- Left: Clinic Brand & Badge -->
        <div class="flex items-center gap-4">
          <button
            @click="emit('navigate', currentBrand)"
            class="flex items-center gap-3 text-left group focus:outline-hidden"
          >
            <!-- Brand Icon Mark -->
            <div
              class="w-10 h-10 rounded-xl flex items-center justify-center font-black text-lg transition-transform group-hover:scale-105 shadow-xs"
              :class="
                currentBrand === 'kidinest'
                  ? 'bg-[#0F4C3A] text-amber-300'
                  : 'bg-zinc-950 text-rose-300'
              "
            >
              <Baby v-if="currentBrand === 'kidinest'" class="w-5 h-5" />
              <Sparkles v-else class="w-5 h-5" />
            </div>

            <!-- Brand Name & Unit Subtitle -->
            <div>
              <div class="flex items-center gap-2">
                <span class="font-black text-xl tracking-tight text-zinc-950">
                  {{ currentBrand === 'kidinest' ? 'KIDINEST' : 'SEMARANG SKIN CENTER' }}
                </span>
                <span
                  class="px-2 py-0.5 text-[10px] font-extrabold uppercase tracking-wider rounded-md"
                  :class="
                    currentBrand === 'kidinest'
                      ? 'bg-emerald-100 text-emerald-800'
                      : 'bg-rose-100 text-rose-800'
                  "
                >
                  {{ currentBrand === 'kidinest' ? 'Poli Anak' : 'Estetika & Kulit' }}
                </span>
              </div>
              <p class="text-xs text-zinc-500 font-medium hidden sm:block">
                {{
                  currentBrand === 'kidinest'
                    ? 'Klinik Spesialis Anak & Tumbuh Kembang'
                    : 'Pusat Perawatan Dermatologi & Estetika Medis'
                }}
              </p>
            </div>
          </button>
        </div>

        <!-- Center: Navigation Anchor Links -->
        <nav class="hidden lg:flex items-center gap-7 text-sm font-bold text-zinc-700">
          <button
            @click="scrollToSection('live-antrian')"
            class="hover:text-zinc-950 transition-colors flex items-center gap-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
            Antrian Hari Ini
          </button>
          <button
            @click="scrollToSection('jadwal-dokter')"
            class="hover:text-zinc-950 transition-colors"
          >
            Jadwal Dokter
          </button>
          <button
            @click="scrollToSection('layanan')"
            class="hover:text-zinc-950 transition-colors"
          >
            Layanan
          </button>
          <button
            @click="scrollToSection('lokasi')"
            class="hover:text-zinc-950 transition-colors"
          >
            Lokasi & Kontak
          </button>
        </nav>

        <!-- Right: Brand Switcher Pill & CTA Button -->
        <div class="hidden sm:flex items-center gap-3">
          <!-- Switcher to Other Brand -->
          <button
            v-if="currentBrand === 'kidinest'"
            @click="emit('navigate', 'ssc')"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-zinc-300 hover:border-zinc-900 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all shadow-2xs"
            title="Kunjungi Semarang Skin Center"
          >
            <Sparkles class="w-3.5 h-3.5 text-rose-600" />
            <span>Ke Semarang Skin Center</span>
            <ArrowRight class="w-3 h-3 text-zinc-400" />
          </button>

          <button
            v-else
            @click="emit('navigate', 'kidinest')"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full border border-zinc-300 hover:border-zinc-900 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all shadow-2xs"
            title="Kunjungi Kidinest Anak"
          >
            <Baby class="w-3.5 h-3.5 text-emerald-700" />
            <span>Ke Kidinest Poli Anak</span>
            <ArrowRight class="w-3 h-3 text-zinc-400" />
          </button>

          <!-- Primary Action CTA -->
          <button
            @click="emit('navigate', 'kiosk')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-extrabold text-white transition-all shadow-md active:scale-95"
            :class="
              currentBrand === 'kidinest'
                ? 'bg-[#0F4C3A] hover:bg-[#0A3629]'
                : 'bg-zinc-950 hover:bg-zinc-800'
            "
          >
            <span>Daftar Antrian</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </div>

        <!-- Mobile Menu Toggle Button -->
        <div class="flex items-center gap-2 sm:hidden">
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="p-2 rounded-lg text-zinc-700 hover:text-zinc-950 hover:bg-zinc-100"
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
      class="sm:hidden border-t px-4 pt-4 pb-6 space-y-4 bg-white shadow-xl animate-in fade-in slide-in-from-top-2 duration-150"
      :class="currentBrand === 'kidinest' ? 'border-emerald-950/10' : 'border-zinc-200'"
    >
      <!-- Switcher in Mobile -->
      <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200 flex items-center justify-between">
        <span class="text-xs font-bold text-zinc-600">Pindah Poli Klinik:</span>
        <button
          v-if="currentBrand === 'kidinest'"
          @click="emit('navigate', 'ssc'); mobileMenuOpen = false"
          class="text-xs font-extrabold text-rose-700 flex items-center gap-1"
        >
          <span>Semarang Skin Center</span>
          <ArrowRight class="w-3 h-3" />
        </button>
        <button
          v-else
          @click="emit('navigate', 'kidinest'); mobileMenuOpen = false"
          class="text-xs font-extrabold text-emerald-800 flex items-center gap-1"
        >
          <span>Kidinest Poli Anak</span>
          <ArrowRight class="w-3 h-3" />
        </button>
      </div>

      <nav class="space-y-2 text-sm font-bold text-zinc-800">
        <button
          @click="scrollToSection('live-antrian')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100 flex items-center justify-between"
        >
          <span>Antrian Hari Ini</span>
          <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
        </button>
        <button
          @click="scrollToSection('jadwal-dokter')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Jadwal Dokter
        </button>
        <button
          @click="scrollToSection('layanan')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Layanan Poli
        </button>
        <button
          @click="scrollToSection('lokasi')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Lokasi & Jam Praktik
        </button>
      </nav>

      <div class="pt-2">
        <button
          @click="emit('navigate', 'kiosk'); mobileMenuOpen = false"
          class="w-full py-3 rounded-xl font-black text-sm text-center text-white shadow-md"
          :class="currentBrand === 'kidinest' ? 'bg-[#0F4C3A]' : 'bg-zinc-950'"
        >
          Daftar Nomor Antrian Online →
        </button>
      </div>
    </div>
  </header>
</template>
