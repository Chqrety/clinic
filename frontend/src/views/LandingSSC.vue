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
  Flame,
  Zap,
  Award,
  ChevronRight,
  Baby
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
    specialty: 'Spesialis Kulit • Bedah Minor & Scar Treatment',
    room: 'Ruang Tindakan Bedah Kulit 2',
    status: 'active'
  }
]

const services = [
  {
    icon: Zap,
    title: 'Laser Skin Rejuvenation',
    desc: 'Teknologi laser medis mutakhir untuk mengatasi pori besar, pigmentasi, melasma, dan peremajaan tekstur kulit.',
    tag: 'Advanced Medical Laser'
  },
  {
    icon: Flame,
    title: 'Acne & Atrophic Scar Care',
    desc: 'Protokol klinis tuntas untuk jerawat meradang dan rekonstruksi bopeng bekas jerawat dengan subsisi & RF Microneedling.',
    tag: 'Clinical Protocol'
  },
  {
    icon: Sparkles,
    title: 'Medical Aesthetic Facial',
    desc: 'Pembersihan komedo mendalam secara steril dan infus serum medis aktif yang disesuaikan kondisi kulit Anda.',
    tag: 'Dermatologist Safe'
  },
  {
    icon: ShieldCheck,
    title: 'Bedah Minor & Patologi Kulit',
    desc: 'Pemeriksaan dermoskopi tahi lalat, biopsi lesi, dan elektrokauterisasi kista atau keratosis secara higienis.',
    tag: 'Surgical Precision'
  }
]

const faqs = [
  {
    q: 'Apakah seluruh perawatan di SSC dilakukan langsung oleh dokter spesialis?',
    a: 'Benar. Seluruh konsultasi mendalam, penetapan resep krim medis, serta tindakan berdaya tinggi (Laser, Botox, Filler, Cauter, Subsisi) dilakukan langsung oleh Dokter Spesialis Dermatologi, Venereologi, dan Estetika (Sp.D.V.E) berizin praktik resmi.'
  },
  {
    q: 'Bagaimana cara memantau antrian ruang konsultasi atau ruang tindakan?',
    a: 'Anda dapat mendaftar mandiri melalui tombol "Reservasi Antrian Konsultasi" untuk mendapatkan nomor antrian digital. Panggilan antrian terupdate secara real-time di layar ponsel Anda.'
  },
  {
    q: 'Di mana lokasi Semarang Skin Center (SSC)?',
    a: 'Semarang Skin Center berlokasi di Klinik Utama Semarang Medical Group, Jl. Melati Utara No. 3, Kota Semarang (lokasi yang sama dengan Kidinest Pediatric Care).'
  }
]
</script>

<template>
  <div class="min-h-screen bg-[#FFFFFF] text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-rose-200 flex flex-col">

    <!-- Sticky Dual-Brand Header -->
    <BrandSwitcherHeader
      current-brand="ssc"
      @navigate="(v) => emit('navigate', v)"
    />

    <!-- Hero Section (Clean Luxury Clinical, Radiant & Action-Focused) -->
    <section class="w-full border-b border-rose-100/70 py-12 sm:py-20 bg-gradient-to-b from-[#FFFFFF] to-[#FFF1F2]/30">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Category Pill Badge -->
        <div class="mb-6">
          <span class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-rose-50 text-rose-900 border border-rose-200 shadow-2xs">
            <Sparkles class="w-4 h-4 text-rose-600" />
            SSC - Semarang Skin Center • Advanced care for your skin
          </span>
        </div>

        <!-- Headline & Direct Description Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center mb-10">

          <!-- Headline & Primary Action (Col 7) -->
          <div class="lg:col-span-7 space-y-6">
            <h1 class="text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-zinc-950 leading-[1.08]">
              Advanced Care for Your Skin.
            </h1>
            <p class="text-base sm:text-lg text-zinc-600 leading-relaxed max-w-2xl font-medium">
              Pusat dermatologi medis dan estetika klinis di Semarang. Ditangani langsung oleh dokter spesialis kulit dengan teknologi laser modern, formula presisi, dan sistem antrian teratur tanpa waktu tunggu berlebih.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-2">
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center justify-center gap-2.5 px-6 py-4 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-sm sm:text-base shadow-lg shadow-zinc-300 transition-all active:scale-95"
              >
                <span>Reservasi Antrian Konsultasi</span>
                <ArrowRight class="w-4 h-4 text-rose-400" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-4 rounded-2xl bg-white hover:bg-rose-50 border border-rose-200 text-rose-950 font-bold text-sm sm:text-base transition-all shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-rose-500 animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>
          </div>

          <!-- Quick Stats / Feature Card (Col 5) -->
          <div class="lg:col-span-5">
            <div class="p-6 sm:p-8 rounded-3xl bg-white border border-rose-100 shadow-md shadow-rose-100/40 space-y-5">
              <div class="flex items-center justify-between border-b border-rose-50 pb-4">
                <span class="text-xs font-black uppercase tracking-wider text-zinc-400">Poli Kulit Hari Ini</span>
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-rose-50 text-rose-800 border border-rose-200">
                  <span class="w-1.5 h-1.5 rounded-full bg-rose-600 animate-pulse"></span>
                  Layanan Aktif
                </span>
              </div>

              <!-- Doctors Row -->
              <div class="space-y-3">
                <div class="flex items-center justify-between p-3.5 rounded-2xl bg-rose-50/40 border border-rose-100">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-zinc-950 text-rose-300 flex items-center justify-center font-bold text-sm shadow-xs">
                      DP
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-950">dr. Dian Paramita, Sp.D.V.E</div>
                      <div class="text-[11px] text-zinc-500">Poli Estetika & Laser • Pagi (09:00 - 14:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-rose-800 bg-white border border-rose-200 px-2.5 py-1 rounded-lg">
                    Praktik
                  </span>
                </div>

                <div class="flex items-center justify-between p-3.5 rounded-2xl bg-zinc-50 border border-zinc-200">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-zinc-800 text-zinc-200 flex items-center justify-center font-bold text-sm shadow-xs">
                      KW
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-950">dr. Kevin Wijaya, Sp.D.V.E</div>
                      <div class="text-[11px] text-zinc-500">Poli Bedah & Tindakan • Sore (14:00 - 19:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-zinc-700 bg-white border border-zinc-200 px-2.5 py-1 rounded-lg">
                    Sore
                  </span>
                </div>
              </div>

              <div class="pt-2 text-xs text-zinc-600 flex items-center justify-between font-semibold border-t border-rose-50">
                <span>Rata-rata Konsultasi Dokter:</span>
                <span class="text-zinc-950 font-black">~15 - 20 Menit / Pasien</span>
              </div>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue & Multi-Doctor Filter Widget -->
    <section id="live-antrian" class="w-full py-16 sm:py-20 border-b border-rose-100/70 bg-[#FAFAFA]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="ssc"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk')"
        />
      </div>
    </section>

    <!-- Section 3: Dokter Spesialis Kulit & Estetika -->
    <section id="jadwal-dokter" class="w-full py-16 sm:py-20 border-b border-zinc-200/80 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10 text-center sm:text-left">
          <span class="text-xs font-black uppercase tracking-wider text-rose-800 bg-rose-50 px-3 py-1 rounded-full border border-rose-200">
            Dermatologist Team
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Dokter Spesialis Dermatologi SSC
          </h2>
          <p class="text-sm sm:text-base text-zinc-600 font-medium mt-1">
            Penanganan klinis berlisensi resmi oleh dokter spesialis kulit dan estetika bersertifikasi.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Doctor 1: dr. Dian Paramita -->
          <div class="p-6 sm:p-8 rounded-3xl border border-rose-200/80 bg-[#FFF1F2]/30 hover:border-rose-400 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-zinc-950 text-rose-300 font-black text-2xl flex items-center justify-center shadow-md shadow-zinc-200">
                  DP
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Dian Paramita, Sp.D.V.E
                  </h3>
                  <p class="text-xs font-bold text-rose-700">
                    Spesialis Dermatologi, Venereologi & Estetika Medis
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Klinik Utama Semarang Medical Group • Ruang Laser 1
                  </p>
                </div>
              </div>

              <div class="space-y-2.5 py-4 border-t border-b border-rose-100 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Jumat (09:00 - 14:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-zinc-950">Melasma, Rejuvenasi Laser, Acne & Skin Glow</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Dian</span>
                <ArrowRight class="w-3.5 h-3.5 text-rose-300" />
              </button>
            </div>
          </div>

          <!-- Doctor 2: dr. Kevin Wijaya -->
          <div class="p-6 sm:p-8 rounded-3xl border border-zinc-200 bg-white hover:border-rose-300 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-zinc-800 text-zinc-100 font-black text-2xl flex items-center justify-center shadow-md shadow-zinc-200">
                  KW
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Kevin Wijaya, Sp.D.V.E
                  </h3>
                  <p class="text-xs font-bold text-zinc-800">
                    Spesialis Dermatologi • Bedah Minor & Scar Treatment
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Klinik Utama Semarang Medical Group • Ruang Tindakan 2
                  </p>
                </div>
              </div>

              <div class="space-y-2.5 py-4 border-t border-b border-zinc-100 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Sabtu (14:00 - 19:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-zinc-950">Bopeng Acne Scar, Kista Kulit, Cauter & Keloid</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Kevin</span>
                <ArrowRight class="w-3.5 h-3.5 text-zinc-400" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 4: Treatment & Layanan Unggulan -->
    <section id="layanan" class="w-full py-16 sm:py-20 border-b border-zinc-200/80 bg-[#FAFAFA]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10">
          <span class="text-xs font-black uppercase tracking-wider text-rose-800 bg-rose-50 px-3 py-1 rounded-full border border-rose-200">
            Clinical Aesthetics
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Protokol Dermatologi Berteknologi Mutakhir
          </h2>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <div
            v-for="(svc, idx) in services"
            :key="idx"
            class="p-6 rounded-3xl bg-white border border-zinc-200 hover:border-rose-300 hover:shadow-md transition-all shadow-2xs flex flex-col justify-between"
          >
            <div>
              <div class="w-12 h-12 rounded-2xl bg-rose-50 text-rose-800 flex items-center justify-center mb-5">
                <component :is="svc.icon" class="w-6 h-6" />
              </div>
              <span class="text-[10px] font-black uppercase tracking-wider text-rose-800 bg-rose-50 px-2.5 py-1 rounded-md">
                {{ svc.tag }}
              </span>
              <h3 class="text-base font-extrabold text-zinc-950 mt-3 mb-2">
                {{ svc.title }}
              </h3>
              <p class="text-xs text-zinc-600 leading-relaxed font-medium">
                {{ svc.desc }}
              </p>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="text-xs font-extrabold text-zinc-950 hover:text-rose-700 flex items-center gap-1 group"
              >
                <span>Reservasi Treatment</span>
                <ChevronRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 5: Lokasi, Jam Buka & FAQ -->
    <section id="lokasi" class="w-full py-16 sm:py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">

          <!-- Left: Location & Hours (Col 5) -->
          <div class="lg:col-span-5 space-y-6">
            <div>
              <span class="text-xs font-black uppercase tracking-wider text-zinc-400">
                Lokasi & Janji Temu
              </span>
              <h2 class="text-3xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Semarang Skin Center (SSC)
              </h2>
            </div>

            <div class="p-6 rounded-3xl bg-zinc-50 border border-zinc-200 space-y-4 text-xs">
              <div class="flex items-start gap-3">
                <MapPin class="w-5 h-5 text-rose-700 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Alamat Klinik SSC:</div>
                  <div class="text-zinc-600 mt-0.5">
                    Klinik Utama Semarang Medical Group<br />
                    Jl. Melati Utara No. 3, Kota Semarang, Jawa Tengah
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-zinc-200 pt-3">
                <Clock class="w-5 h-5 text-rose-700 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Jam Praktik Dokter:</div>
                  <div class="text-zinc-600 mt-0.5 space-y-0.5">
                    <div>Poli Estetika: 09:00 - 14:00 WIB (dr. Dian Paramita)</div>
                    <div>Poli Bedah Kulit: 14:00 - 19:00 WIB (dr. Kevin Wijaya)</div>
                    <div class="text-[11px] text-zinc-700 font-bold">Konsultasi Minggu: Dengan Perjanjian</div>
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-zinc-200 pt-3">
                <Phone class="w-5 h-5 text-rose-700 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">WhatsApp Care & Reservasi:</div>
                  <div class="text-zinc-600 mt-0.5 font-mono">
                    +62 811-9876-5432 (Layanan Pasien VIP SSC)
                  </div>
                </div>
              </div>
            </div>

            <!-- Switch to Kidinest Banner -->
            <div class="p-5 rounded-3xl bg-sky-950 text-white flex items-center justify-between">
              <div>
                <div class="text-xs font-bold text-sky-300">Mencari poli anak & vaksinasi?</div>
                <div class="text-sm font-black text-white">Kidinest Pediatric Care</div>
              </div>
              <button
                @click="emit('navigate', 'kidinest')"
                class="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/20 text-white text-xs font-bold transition-colors"
              >
                Lihat Kidinest →
              </button>
            </div>
          </div>

          <!-- Right: FAQ Section (Col 7) -->
          <div class="lg:col-span-7 space-y-4">
            <div class="mb-4">
              <span class="text-xs font-black uppercase tracking-wider text-rose-800 bg-rose-50 px-3 py-1 rounded-full border border-rose-200">
                Informasi Pasien
              </span>
              <h3 class="text-2xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Pertanyaan Seputar Perawatan Kulit di SSC
              </h3>
            </div>

            <div class="space-y-3">
              <div
                v-for="(faq, idx) in faqs"
                :key="idx"
                class="p-5 rounded-2xl border border-zinc-200 bg-zinc-50/60"
              >
                <h4 class="text-sm font-extrabold text-zinc-950 mb-2">
                  {{ faq.q }}
                </h4>
                <p class="text-xs text-zinc-600 leading-relaxed font-medium">
                  {{ faq.a }}
                </p>
              </div>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Footer -->
    <footer class="w-full border-t border-zinc-200 py-8 bg-white text-xs text-zinc-500">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          © 2026 <strong>SSC - Semarang Skin Center</strong> • Bagian dari Semarang Medical Group.
        </div>
        <div class="flex items-center gap-4 font-bold text-zinc-700">
          <button @click="emit('navigate', 'kiosk')" class="hover:text-zinc-950">Portal Reservasi Pasien</button>
          <span>•</span>
          <button @click="emit('navigate', 'kidinest')" class="hover:text-zinc-950">Kidinest Pediatric Care</button>
          <span>•</span>
          <button @click="emit('navigate', 'login')" class="hover:text-zinc-950">Login Staf Medis</button>
        </div>
      </div>
    </footer>

  </div>
</template>
