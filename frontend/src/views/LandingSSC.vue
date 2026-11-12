<script setup lang="ts">
import { ref } from 'vue'
import BrandSwitcherHeader from '../components/BrandSwitcherHeader.vue'
import LiveQueueFilter, { type DoctorInfo } from '../components/LiveQueueFilter.vue'
import {
  Sparkles,
  ShieldCheck,
  Clock,
  CheckCircle2,
  MapPin,
  Phone,
  Stethoscope,
  HeartPulse,
  ArrowRight,
  ChevronRight,
  Activity
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login'): void
}>()

const doctors: DoctorInfo[] = [
  {
    id: 'dian',
    name: 'dr. Dian Paramita, Sp.D.V.E',
    specialty: 'Spesialis Dermatologi, Venereologi & Estetika Medis',
    room: 'Ruang Laser & Konsultasi 1',
    status: 'active'
  },
  {
    id: 'kevin',
    name: 'dr. Kevin Wijaya, Sp.D.V.E',
    specialty: 'Spesialis Dermatologi • Bedah Minor & Scar Treatment',
    room: 'Ruang Tindakan Bedah Kulit 2',
    status: 'active'
  }
]

// Core Treatments in Editorial Index Table (Anti-Card Format)
const treatmentsIndex = [
  {
    num: '01',
    title: 'Laser Skin Rejuvenation & Pigmentation',
    focus: 'Melasma, Hiperpigmentasi, Pori Membesar & Peremajaan Kolagen',
    technology: 'Medical Laser FDA Approved',
    doctor: 'dr. Dian Paramita, Sp.D.V.E'
  },
  {
    num: '02',
    title: 'Acne & Atrophic Scar Reconstruction',
    focus: 'Penanganan Jerawat Aktif, Subcision Bopeng & RF Microneedling',
    technology: 'Multimodal Scar Protocol',
    doctor: 'dr. Kevin Wijaya, Sp.D.V.E'
  },
  {
    num: '03',
    title: 'Medical Aesthetic Facial & Deep Infusion',
    focus: 'Pembersihan Komedo Steril, Peeling Asam Buah Medis & Serum Aktif',
    technology: 'Dermatologist Certified',
    doctor: 'dr. Dian Paramita, Sp.D.V.E'
  },
  {
    num: '04',
    title: 'Bedah Minor & Dermoskopi Patologi Kulit',
    focus: 'Pemeriksaan Tahi Lalat, Biopsi Lesi & Kauterisasi Kista/Keratosis',
    technology: 'High-Precision Electrosurgery',
    doctor: 'dr. Kevin Wijaya, Sp.D.V.E'
  }
]

// Clinical Protocols (Official Feed Do's & Don'ts)
const clinicalDos = [
  'Gunakan tabir surya SPF 30+ berformula broad-spectrum minimal 15 menit sebelum beraktivitas di luar ruangan.',
  'Aplikasikan ulang (re-apply) tabir surya setiap 2–3 jam, terutama saat berkeringat atau terpapar sinar matahari langsung.',
  'Selalu lindungi area leher, dada atas, dan punggung tangan yang kerap terlewat dari perlindungan UV.',
  'Lakukan konsultasi berkala dengan dokter spesialis jika muncul bercak pigmentasi atau lesi baru pada kulit.'
]

const clinicalDonts = [
  'Jangan mengaplikasikan tabir surya sesaat setelah berada di bawah terik matahari.',
  'Jangan menggunakan produk eksfoliasi keras (scrub/retinoid tinggi) saat kulit sedang meradang atau terbakar matahari.',
  'Jangan memencet atau menusuk sendiri jerawat kistik karena berisiko memicu bopeng atrofi permanen.',
  'Jangan mengabaikan perubahan bentuk, warna, atau ukuran pada tahi lalat yang membesar secara asimetris.'
]

const faqs = [
  {
    q: 'Apakah semua tindakan di Semarang Skin Center ditangani langsung oleh dokter spesialis?',
    a: 'Benar. Seluruh konsultasi mendalam, penetapan resep racikan medis, dan tindakan berdaya tinggi (Laser, Botox, Subcision, Kauter) dilakukan langsung oleh Dokter Spesialis Dermatologi, Venereologi, dan Estetika (Sp.D.V.E) berizin resmi.'
  },
  {
    q: 'Bagaimana cara reservasi antrian konsultasi di SSC?',
    a: 'Anda dapat mengklik tombol "Reservasi Antrian Konsultasi", memilih dokter spesialis kulit yang dituju, dan memantau status giliran ruang tindakan secara live dari ponsel Anda.'
  },
  {
    q: 'Di mana lokasi Semarang Skin Center?',
    a: 'SSC berlokasi di Klinik Utama Semarang Medical Group, Jl. Melati Utara No. 3, Semarang.'
  }
]
</script>

<template>
  <div class="min-h-screen bg-[#FAF5F5] text-[#2D1E24] font-['Plus_Jakarta_Sans',sans-serif] selection:bg-[#9E606D]/20 flex flex-col">

    <!-- Sticky Dual-Brand Header -->
    <BrandSwitcherHeader
      current-brand="ssc"
      @navigate="(v) => emit('navigate', v)"
    />

    <!-- Hero Section (Clinical Editorial Luxury, Anti-Card, Skala Manusiawi) -->
    <section class="w-full border-b border-[#9E606D]/15 py-12 sm:py-16 bg-gradient-to-b from-[#FAF5F5] via-[#FAF5F5] to-[#F5ECEE]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Top Pill Tag (Subtle Rose Badge) -->
        <div class="mb-5">
          <span class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-bold bg-[#9E606D]/10 text-[#9E606D] border border-[#9E606D]/25">
            <img src="/assets/ssc/logo-ssc.png" alt="SSC Rosette" class="w-3.5 h-3.5 object-contain" />
            <span>SSC • Semarang Skin Center • Semarang Medical Group</span>
          </span>
        </div>

        <!-- Headline & Action Grid (Hairline Editorial Layout) -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">

          <!-- Left Column (Narasi & Tindakan Utama) -->
          <div class="lg:col-span-7 space-y-5">
            <h1 class="text-3xl sm:text-5xl font-extrabold text-[#2D1E24] tracking-tight leading-[1.15]">
              Advanced Care <span class="font-['Playfair_Display',serif] italic font-normal text-[#9E606D]">for Your Skin.</span>
            </h1>
            <p class="text-sm sm:text-base text-[#2D1E24]/80 leading-relaxed max-w-xl font-medium">
              Pusat dermatologi medis dan estetika klinis terpercaya di Semarang. Ditangani langsung oleh dokter spesialis kulit dengan teknologi laser medis modern, protokol presisi, dan antrian ruang tindakan transparan.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-2">
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center justify-center gap-2.5 px-6 py-3.5 rounded-xl bg-[#9E606D] hover:bg-[#8B505D] text-white font-bold text-xs sm:text-sm uppercase tracking-wider shadow-md shadow-[#9E606D]/25 transition-all active:scale-95"
              >
                <span>Reservasi Antrian Konsultasi</span>
                <ArrowRight class="w-4 h-4 text-white/80" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-3.5 rounded-xl bg-white hover:bg-white/80 border border-[#9E606D]/30 text-[#2D1E24] font-bold text-xs sm:text-sm transition-all shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-[#9E606D] animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>

            <!-- Mini Trust Strip -->
            <div class="pt-3 flex items-center gap-6 text-xs text-[#2D1E24]/75 font-semibold">
              <div class="flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-[#9E606D]" />
                <span>Dokter Spesialis Kulit Sp.D.V.E</span>
              </div>
              <div class="flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-[#9E606D]" />
                <span>Medical Laser FDA Certified</span>
              </div>
            </div>
          </div>

          <!-- Right Column: Live Clinical Strip (Hairline Separator) -->
          <div class="lg:col-span-5">
            <div class="p-6 rounded-3xl bg-white/90 border border-[#9E606D]/20 shadow-sm relative space-y-4">
              
              <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
                <div class="flex items-center gap-2">
                  <div class="w-6 h-6 rounded-full bg-[#9E606D] flex items-center justify-center p-1">
                    <img src="/assets/ssc/logo-ssc.png" class="w-4 h-4 object-contain" />
                  </div>
                  <span class="font-bold text-xs text-[#2D1E24]">Layanan Poli Kulit Hari Ini</span>
                </div>
                <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-rose-50 text-[#9E606D] border border-[#9E606D]/20">
                  Konsultasi Aktif
                </span>
              </div>

              <!-- Doctors Today Quick Rows (Hairline) -->
              <div class="space-y-2.5">
                <div class="flex items-center justify-between p-3 rounded-2xl bg-[#FAF5F5] border border-[#9E606D]/15">
                  <div class="flex items-center gap-3">
                    <span class="w-8 h-8 rounded-xl bg-[#9E606D]/15 text-[#9E606D] font-bold text-xs flex items-center justify-center">
                      DP
                    </span>
                    <div>
                      <div class="text-xs font-bold text-[#2D1E24]">dr. Dian Paramita, Sp.D.V.E</div>
                      <div class="text-[11px] text-[#2D1E24]/60">Poli Laser & Estetika (09:00 - 14:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-semibold text-[#9E606D]">Praktik Pagi</span>
                </div>

                <div class="flex items-center justify-between p-3 rounded-2xl bg-[#FAF5F5] border border-[#9E606D]/15">
                  <div class="flex items-center gap-3">
                    <span class="w-8 h-8 rounded-xl bg-[#2D1E24]/10 text-[#2D1E24] font-bold text-xs flex items-center justify-center">
                      KW
                    </span>
                    <div>
                      <div class="text-xs font-bold text-[#2D1E24]">dr. Kevin Wijaya, Sp.D.V.E</div>
                      <div class="text-[11px] text-[#2D1E24]/60">Poli Bedah & Tindakan (14:00 - 19:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-semibold text-[#2D1E24]/70">Praktik Sore</span>
                </div>
              </div>

              <div class="pt-2 border-t border-zinc-100 flex items-center justify-between text-xs">
                <span class="text-[#2D1E24]/70">Waktu Konsultasi Mendalam:</span>
                <span class="font-bold text-[#2D1E24]">~15 - 20 Menit / Pasien</span>
              </div>

            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue Strip (Menyatu Alami dengan Kanvas) -->
    <section id="live-antrian" class="w-full py-12 sm:py-16 border-b border-[#9E606D]/15 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="ssc"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk')"
        />
      </div>
    </section>

    <!-- Section 3: Menu Perawatan Medis (Editorial Index Table - Bebas Sindrom Card Berjejer!) -->
    <section id="layanan" class="w-full py-12 sm:py-16 border-b border-[#9E606D]/15 bg-[#FAF5F5]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Heading Proporsional -->
        <div class="mb-8 flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div>
            <span class="text-xs font-bold uppercase tracking-wider text-[#9E606D]">
              Medical Dermatology Index
            </span>
            <h2 class="text-2xl sm:text-3xl font-extrabold text-[#2D1E24] tracking-tight mt-1">
              Protokol & Perawatan Klinis Kulit
            </h2>
          </div>
          <button
            @click="emit('navigate', 'kiosk')"
            class="text-xs font-bold uppercase tracking-wider text-[#9E606D] underline hover:text-[#2D1E24] transition-colors"
          >
            Reservasi Konsultasi Treatment →
          </button>
        </div>

        <!-- Editorial Index Table (Hairline Borders 1px) -->
        <div class="border-t border-b border-[#9E606D]/20 divide-y divide-[#9E606D]/15">
          <div
            v-for="item in treatmentsIndex"
            :key="item.num"
            class="py-5 sm:py-6 flex flex-col md:flex-row md:items-center justify-between gap-4 group hover:bg-white/50 transition-colors px-2 sm:px-4 rounded-xl"
          >
            <div class="flex items-start sm:items-center gap-5">
              <span class="font-mono text-sm font-bold text-[#9E606D] shrink-0 pt-0.5 sm:pt-0">
                {{ item.num }}
              </span>
              <div>
                <div class="flex items-center gap-3">
                  <h3 class="text-base sm:text-lg font-bold text-[#2D1E24] group-hover:text-[#9E606D] transition-colors">
                    {{ item.title }}
                  </h3>
                  <span class="hidden sm:inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-[#9E606D]/10 text-[#9E606D]">
                    {{ item.technology }}
                  </span>
                </div>
                <p class="text-xs sm:text-sm text-[#2D1E24]/75 mt-0.5 font-medium">
                  {{ item.focus }}
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between sm:justify-end gap-5 shrink-0 pl-10 md:pl-0">
              <span class="text-xs font-semibold text-[#2D1E24]/70">
                {{ item.doctor }}
              </span>
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center gap-1.5 text-xs font-bold text-[#9E606D] hover:text-[#2D1E24] transition-colors"
              >
                <span>Reservasi</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 4: Protokol Klinis Resmi (Do's vs Don'ts Format - Mengadopsi Layout Feed SPF) -->
    <section class="w-full py-12 sm:py-16 border-b border-[#9E606D]/15 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10 text-center max-w-2xl mx-auto">
          <span class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold bg-rose-50 text-[#9E606D] border border-[#9E606D]/20">
            Dermatologist Educational Protocol
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-[#2D1E24] tracking-tight mt-3">
            Tips for <span class="font-['Playfair_Display',serif] italic font-normal text-[#9E606D]">Healthy & Protected Skin</span>
          </h2>
          <p class="text-xs sm:text-sm text-[#2D1E24]/70 mt-1 font-medium">
            Panduan harian menjaga integritas pelindung kulit (*skin barrier*) sesuai rekomendasi dokter spesialis kulit SSC.
          </p>
        </div>

        <!-- Two-Column Breath Layout (Hairline Center Divider) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12 relative max-w-5xl mx-auto items-start">
          
          <!-- Hairline Divider di tengah desktop -->
          <div class="hidden md:block absolute top-0 bottom-0 left-1/2 w-px bg-[#9E606D]/15"></div>

          <!-- Kolom Kiri: DO'S -->
          <div class="space-y-4 md:pr-6">
            <div class="flex items-center gap-2 pb-2 border-b border-[#9E606D]/20">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-600"></span>
              <h3 class="text-base font-bold text-[#2D1E24] tracking-wide uppercase">
                Do's: Protokol yang Dianjurkan
              </h3>
            </div>

            <ul class="space-y-3.5 text-xs sm:text-sm text-[#2D1E24]/80">
              <li
                v-for="(item, idx) in clinicalDos"
                :key="idx"
                class="flex items-start gap-3"
              >
                <span class="font-bold text-[#9E606D]">•</span>
                <span class="leading-relaxed font-medium">{{ item }}</span>
              </li>
            </ul>
          </div>

          <!-- Kolom Kanan: DON'TS -->
          <div class="space-y-4 md:pl-6">
            <div class="flex items-center gap-2 pb-2 border-b border-[#9E606D]/20">
              <span class="w-2.5 h-2.5 rounded-full bg-rose-600"></span>
              <h3 class="text-base font-bold text-[#2D1E24] tracking-wide uppercase">
                Don'ts: Kebiasaan yang Dihindari
              </h3>
            </div>

            <ul class="space-y-3.5 text-xs sm:text-sm text-[#2D1E24]/80">
              <li
                v-for="(item, idx) in clinicalDonts"
                :key="idx"
                class="flex items-start gap-3"
              >
                <span class="font-bold text-rose-600">•</span>
                <span class="leading-relaxed font-medium">{{ item }}</span>
              </li>
            </ul>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 5: Dokter Spesialis Kulit SSC -->
    <section id="jadwal-dokter" class="w-full py-12 sm:py-16 border-b border-[#9E606D]/15 bg-[#FAF5F5]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-8">
          <span class="text-xs font-bold uppercase tracking-wider text-[#9E606D]">
            Dermatologist Faculty
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-[#2D1E24] tracking-tight mt-1">
            Dokter Spesialis Dermatologi & Estetika SSC
          </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Dokter 1: dr. Dian Paramita -->
          <div class="p-6 rounded-3xl bg-white border border-[#9E606D]/20 shadow-2xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-14 h-14 rounded-2xl bg-[#9E606D] text-white font-black text-xl flex items-center justify-center shadow-xs">
                  DP
                </div>
                <div>
                  <h3 class="text-lg font-bold text-[#2D1E24]">
                    dr. Dian Paramita, Sp.D.V.E
                  </h3>
                  <p class="text-xs font-semibold text-[#9E606D] mt-0.5">
                    Spesialis Dermatologi, Venereologi & Estetika Medis
                  </p>
                  <p class="text-[11px] text-[#2D1E24]/60">
                    Ruang Laser & Konsultasi 1 • Semarang Medical Group
                  </p>
                </div>
              </div>

              <div class="py-3 border-t border-b border-zinc-100 space-y-1.5 text-xs text-[#2D1E24]/80 font-medium">
                <div class="flex justify-between">
                  <span>Jadwal Praktik Pagi:</span>
                  <span class="font-bold text-[#2D1E24]">Senin - Jumat (09:00 - 14:00)</span>
                </div>
                <div class="flex justify-between">
                  <span>Fokus Klinis:</span>
                  <span class="font-bold text-[#2D1E24]">Laser Rejuvenation, Melasma & Anti-Aging</span>
                </div>
              </div>
            </div>

            <div class="pt-5">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#9E606D] hover:bg-[#8B505D] text-white text-xs font-bold uppercase tracking-wider transition-all"
              >
                Daftar Konsultasi ke dr. Dian →
              </button>
            </div>
          </div>

          <!-- Dokter 2: dr. Kevin Wijaya -->
          <div class="p-6 rounded-3xl bg-white border border-[#9E606D]/20 shadow-2xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-14 h-14 rounded-2xl bg-[#2D1E24] text-white font-black text-xl flex items-center justify-center shadow-xs">
                  KW
                </div>
                <div>
                  <h3 class="text-lg font-bold text-[#2D1E24]">
                    dr. Kevin Wijaya, Sp.D.V.E
                  </h3>
                  <p class="text-xs font-semibold text-[#2D1E24] mt-0.5">
                    Spesialis Dermatologi • Bedah Minor & Scar Treatment
                  </p>
                  <p class="text-[11px] text-[#2D1E24]/60">
                    Ruang Tindakan Bedah Kulit 2 • Semarang Medical Group
                  </p>
                </div>
              </div>

              <div class="py-3 border-t border-b border-zinc-100 space-y-1.5 text-xs text-[#2D1E24]/80 font-medium">
                <div class="flex justify-between">
                  <span>Jadwal Praktik Sore:</span>
                  <span class="font-bold text-[#2D1E24]">Senin - Sabtu (14:00 - 19:00)</span>
                </div>
                <div class="flex justify-between">
                  <span>Fokus Klinis:</span>
                  <span class="font-bold text-[#2D1E24]">Acne Scar Bopeng, Cauter Kista & Keloid</span>
                </div>
              </div>
            </div>

            <div class="pt-5">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#2D1E24] hover:bg-zinc-800 text-white text-xs font-bold uppercase tracking-wider transition-all"
              >
                Daftar Konsultasi ke dr. Kevin →
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 6: Lokasi & FAQ -->
    <section id="lokasi" class="w-full py-12 sm:py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
          <!-- Left: Alamat & Kontak (Col 5) -->
          <div class="lg:col-span-5 space-y-5">
            <div>
              <span class="text-xs font-bold uppercase tracking-wider text-[#9E606D]">
                Akses & Reservasi VIP
              </span>
              <h2 class="text-2xl sm:text-3xl font-extrabold text-[#2D1E24] tracking-tight mt-1">
                Semarang Skin Center (SSC)
              </h2>
            </div>

            <div class="p-5 rounded-2xl bg-[#FAF5F5] border border-[#9E606D]/20 space-y-3 text-xs text-[#2D1E24]">
              <div class="flex items-start gap-2.5">
                <MapPin class="w-4 h-4 text-[#9E606D] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">Alamat Resmi:</div>
                  <div class="opacity-80 mt-0.5">
                    Klinik Utama Semarang Medical Group<br />
                    Jl. Melati Utara No. 3, Kota Semarang, Jawa Tengah
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-2.5 border-t border-zinc-200/60 pt-2.5">
                <Clock class="w-4 h-4 text-[#9E606D] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">Jam Layanan Poli Kulit:</div>
                  <div class="opacity-80 mt-0.5">
                    Pagi: 09:00 - 14:00 WIB • Sore: 14:00 - 19:00 WIB
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-2.5 border-t border-zinc-200/60 pt-2.5">
                <Phone class="w-4 h-4 text-[#9E606D] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">WhatsApp Care & Reservasi:</div>
                  <div class="opacity-80 mt-0.5 font-mono">
                    +62 811-9876-5432
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right: FAQ (Col 7) -->
          <div class="lg:col-span-7 space-y-3">
            <h3 class="text-lg font-bold text-[#2D1E24] mb-3">
              Pertanyaan yang Sering Diajukan
            </h3>

            <div
              v-for="(faq, idx) in faqs"
              :key="idx"
              class="p-4 rounded-xl border border-zinc-200/80 bg-[#FAF5F5]/60"
            >
              <h4 class="text-xs sm:text-sm font-bold text-[#2D1E24] mb-1">
                {{ faq.q }}
              </h4>
              <p class="text-xs text-[#2D1E24]/75 leading-relaxed font-medium">
                {{ faq.a }}
              </p>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Footer Solid Pita Dusty Rose Khas SSC -->
    <footer class="w-full bg-[#9E606D] text-white py-6 text-xs select-none">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex flex-col sm:flex-row items-center justify-between gap-3">
        <div class="flex items-center gap-2">
          <img src="/assets/ssc/logo-ssc.png" class="w-5 h-5 object-contain" />
          <span class="font-bold">SSC - Semarang Skin Center • Advanced care for your skin</span>
        </div>
        <div class="font-medium text-white/90">
          Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang
        </div>
        <div class="flex items-center gap-4 text-white/80 font-semibold">
          <button @click="emit('navigate', 'kiosk')" class="hover:text-white">Portal Reservasi</button>
          <span>•</span>
          <button @click="emit('navigate', 'kidinest')" class="hover:text-white">Kidinest Pediatric</button>
          <span>•</span>
          <button @click="emit('navigate', 'login')" class="hover:text-white">Login Staf</button>
        </div>
      </div>
    </footer>

  </div>
</template>
