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
  Stethoscope,
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
    class="sticky top-0 z-40 border-b backdrop-blur-md transition-colors"
    :class="
      currentBrand === 'kidinest'
        ? 'bg-[#FFFDF9]/95 border-sky-100 shadow-2xs'
        : 'bg-[#FFFFFF]/95 border-rose-100/80 shadow-2xs'
    "
  >
    <!-- Sub-header top ribbon: Single Clinic Family notice -->
    <div
      class="w-full py-1 text-[11px] font-bold text-center border-b transition-colors"
      :class="
        currentBrand === 'kidinest'
          ? 'bg-sky-50/70 border-sky-100/70 text-sky-800'
          : 'bg-rose-50/60 border-rose-100/70 text-rose-900'
      "
    >
      <div class="max-w-7xl mx-auto px-4 flex items-center justify-between">
        <div class="flex items-center gap-1.5">
          <MapPin class="w-3 h-3 text-zinc-500" />
          <span>Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang</span>
        </div>
        <div class="hidden sm:flex items-center gap-3">
          <span>Buka Setiap Hari Praktik</span>
          <span>•</span>
          <span>WhatsApp Reservasi: 0812-3456-7890</span>
        </div>
      </div>
    </div>

    <!-- Main Navigation Bar -->
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
              class="w-11 h-11 rounded-2xl flex items-center justify-center font-black text-lg transition-transform group-hover:scale-105 shadow-sm"
              :class="
                currentBrand === 'kidinest'
                  ? 'bg-gradient-to-br from-sky-400 to-sky-600 text-white shadow-sky-200'
                  : 'bg-gradient-to-br from-rose-500 to-zinc-900 text-white shadow-rose-200'
              "
            >
              <Baby v-if="currentBrand === 'kidinest'" class="w-6 h-6" />
              <Sparkles v-else class="w-5 h-5" />
            </div>

            <!-- Brand Name & Unit Subtitle -->
            <div>
              <div class="flex items-center gap-2">
                <span class="font-extrabold text-xl tracking-tight text-zinc-950">
                  {{ currentBrand === 'kidinest' ? 'KIDINEST' : 'SEMARANG SKIN CENTER' }}
                </span>
                <span
                  class="px-2 py-0.5 text-[10px] font-black uppercase tracking-wider rounded-md"
                  :class="
                    currentBrand === 'kidinest'
                      ? 'bg-sky-100 text-sky-800'
                      : 'bg-rose-100 text-rose-800'
                  "
                >
                  {{ currentBrand === 'kidinest' ? 'Pediatric Care' : 'Aesthetic & Skin' }}
                </span>
              </div>
              <p class="text-xs text-zinc-500 font-semibold hidden sm:block">
                {{
                  currentBrand === 'kidinest'
                    ? 'Spesialis Anak & Vaksinasi • dr. Patricia Vanessa A., Sp.A'
                    : 'Dermatology & Aesthetic Center • Advanced Care for Your Skin'
                }}
              </p>
            </div>
          </button>
        </div>

        <!-- Center: Navigation Anchor Links -->
        <nav class="hidden lg:flex items-center gap-7 text-sm font-bold text-zinc-600">
          <button
            @click="scrollToSection('live-antrian')"
            class="hover:text-zinc-950 transition-colors flex items-center gap-1.5"
          >
            <span
              class="w-2 h-2 rounded-full animate-pulse"
              :class="currentBrand === 'kidinest' ? 'bg-sky-500' : 'bg-rose-500'"
            ></span>
            Antrian Live
          </button>
          <button
            @click="scrollToSection('jadwal-dokter')"
            class="hover:text-zinc-950 transition-colors"
          >
            Dokter Spesialis
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
          <!-- Switcher to Other Brand in Semarang Medical Group -->
          <button
            v-if="currentBrand === 'kidinest'"
            @click="emit('navigate', 'ssc')"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full border border-rose-200 hover:border-rose-400 bg-rose-50/70 hover:bg-rose-100 text-xs font-bold text-rose-900 transition-all shadow-2xs"
            title="Kunjungi Semarang Skin Center"
          >
            <Sparkles class="w-3.5 h-3.5 text-rose-600" />
            <span>Poli Kulit (SSC)</span>
            <ArrowRight class="w-3 h-3 text-rose-400" />
          </button>

          <button
            v-else
            @click="emit('navigate', 'kidinest')"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full border border-sky-200 hover:border-sky-400 bg-sky-50/70 hover:bg-sky-100 text-xs font-bold text-sky-900 transition-all shadow-2xs"
            title="Kunjungi Kidinest Pediatric Care"
          >
            <Baby class="w-3.5 h-3.5 text-sky-600" />
            <span>Poli Anak (Kidinest)</span>
            <ArrowRight class="w-3 h-3 text-sky-400" />
          </button>

          <!-- Primary Action CTA -->
          <button
            @click="emit('navigate', 'kiosk')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-extrabold text-white transition-all shadow-md active:scale-95"
            :class="
              currentBrand === 'kidinest'
                ? 'bg-sky-600 hover:bg-sky-700 shadow-sky-200'
                : 'bg-zinc-950 hover:bg-zinc-800 shadow-zinc-300'
            "
          >
            <span>Ambil Antrian</span>
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
      :class="currentBrand === 'kidinest' ? 'border-sky-100' : 'border-rose-100'"
    >
      <!-- Switcher in Mobile -->
      <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200 flex items-center justify-between">
        <span class="text-xs font-bold text-zinc-600">Pindah Layanan Poli:</span>
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
          class="text-xs font-extrabold text-sky-700 flex items-center gap-1"
        >
          <span>Kidinest Pediatric Care</span>
          <ArrowRight class="w-3 h-3" />
        </button>
      </div>

      <nav class="space-y-2 text-sm font-bold text-zinc-800">
        <button
          @click="scrollToSection('live-antrian')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100 flex items-center justify-between"
        >
          <span>Antrian Live</span>
          <span
            class="w-2 h-2 rounded-full"
            :class="currentBrand === 'kidinest' ? 'bg-sky-500' : 'bg-rose-500'"
          ></span>
        </button>
        <button
          @click="scrollToSection('jadwal-dokter')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Dokter Spesialis
        </button>
        <button
          @click="scrollToSection('layanan')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Layanan & Vaksinasi
        </button>
        <button
          @click="scrollToSection('lokasi')"
          class="w-full text-left py-2 px-3 rounded-lg hover:bg-zinc-100"
        >
          Lokasi & Kontak
        </button>
      </nav>

      <div class="pt-2">
        <button
          @click="emit('navigate', 'kiosk'); mobileMenuOpen = false"
          class="w-full py-3 rounded-xl font-black text-sm text-center text-white shadow-md"
          :class="currentBrand === 'kidinest' ? 'bg-sky-600' : 'bg-zinc-950'"
        >
          Ambil Nomor Antrian Online →
        </button>
      </div>
    </div>
  </header>
</template>
