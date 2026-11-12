<script setup lang="ts">
import { ref } from 'vue'
import BrandSwitcherHeader from '../components/BrandSwitcherHeader.vue'
import LiveQueueFilter, { type DoctorInfo } from '../components/LiveQueueFilter.vue'
import {
  Baby,
  Sparkles,
  ShieldCheck,
  Clock,
  CheckCircle2,
  MapPin,
  Phone,
  Stethoscope,
  HeartPulse,
  ArrowRight,
  Syringe,
  Smile,
  Activity,
  Calendar,
  ChevronRight,
  Heart
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login'): void
}>()

const doctors: DoctorInfo[] = [
  {
    id: 'patricia',
    name: 'dr. Patricia Vanessa A., Sp.A',
    specialty: 'Spesialis Anak & Vaksinasi Pediatri',
    room: 'Ruang Poli Anak 1',
    status: 'active'
  },
  {
    id: 'hendra',
    name: 'dr. Hendra Pratama, Sp.A',
    specialty: 'Konsultan Tumbuh Kembang & Nutrisi',
    room: 'Ruang Poli Anak 2',
    status: 'active'
  }
]

const services = [
  {
    icon: Syringe,
    title: 'Vaksinasi & Imunisasi Anak',
    desc: 'Jadwal imunisasi dasar dan lanjutan sesuai panduan IDAI dengan penyimpanan vaksin standar cold-chain medis.',
    tag: 'Standar IDAI'
  },
  {
    icon: Smile,
    title: 'Skrining Tumbuh Kembang',
    desc: 'Pemantauan berkala milestone motorik, sensorik, bahasa, dan grafik pertumbuhan bersama dokter spesialis anak.',
    tag: 'Deteksi Dini'
  },
  {
    icon: Activity,
    title: 'Pemeriksaan Sakit & Infeksi',
    desc: 'Penanganan demam, infeksi saluran napas, batuk, pilek, diare, dan alergi anak dengan resep obat yang aman.',
    tag: 'Responsif & Ramah'
  },
  {
    icon: HeartPulse,
    title: 'Nebulisasi & Fisioterapi Anak',
    desc: 'Terapi uap nebulizer yang nyaman dan stimulasi fisik di ruang periksa ceria yang tidak menakutkan bagi si kecil.',
    tag: 'Kid-Friendly'
  }
]

const faqs = [
  {
    q: 'Apakah pendaftaran antrian Kidinest bisa dilakukan online sebelum datang?',
    a: 'Sangat bisa! Bunda dan Ayah dapat menekan tombol "Ambil Antrian Online" dari ponsel, memilih dokter spesialis anak, dan tiket antrian beserta estimasi jam panggil akan langsung terbit.'
  },
  {
    q: 'Apakah Kidinest melayani vaksinasi bayi baru lahir hingga anak usia sekolah?',
    a: 'Ya, kami menyediakan vaksin lengkap mulai dari Hepatitis B, BCG, Polio, DTP, Hib, PCV, Rotavirus, Influenza, MR/MMR, hingga Cacar Air dan Japanese Encephalitis.'
  },
  {
    q: 'Di mana lokasi Kidinest Pediatric Care?',
    a: 'Kidinest berlokasi di dalam Klinik Utama Semarang Medical Group, Jl. Melati Utara No. 3, Semarang (dekat pusat kota dengan akses parkir luas dan nyaman untuk keluarga).'
  }
]
</script>

<template>
  <div class="min-h-screen bg-[#FFFDF9] text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-sky-200 flex flex-col">

    <!-- Sticky Dual-Brand Header -->
    <BrandSwitcherHeader
      current-brand="kidinest"
      @navigate="(v) => emit('navigate', v)"
    />

    <!-- Hero Section (Friendly, Warm & Action-Oriented) -->
    <section class="w-full border-b border-sky-100 py-12 sm:py-20 bg-gradient-to-b from-[#FFFDF9] to-[#F0F9FF]/60">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Category Pill Badge -->
        <div class="mb-6">
          <span class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-sky-100 text-sky-900 border border-sky-200 shadow-2xs">
            <Baby class="w-4 h-4 text-sky-600" />
            Kidinest Pediatric Care • Semarang Medical Group
          </span>
        </div>

        <!-- Headline & Direct Description Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center mb-10">

          <!-- Headline & Primary Action (Col 7) -->
          <div class="lg:col-span-7 space-y-6">
            <h1 class="text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-zinc-950 leading-[1.08]">
              Pemeriksaan Buah Hati dengan Hangat & Tenang.
            </h1>
            <p class="text-base sm:text-lg text-zinc-600 leading-relaxed max-w-2xl font-medium">
              Pelayanan dokter spesialis anak ramah keluarga bersama <strong>dr. Patricia Vanessa A., Sp.A</strong>. Pantau antrian secara live dari ponsel, hindari ruang tunggu yang padat, dan dapatkan notifikasi WhatsApp saat giliran tiba.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-2">
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center justify-center gap-2.5 px-6 py-4 rounded-2xl bg-sky-600 hover:bg-sky-700 text-white font-extrabold text-sm sm:text-base shadow-lg shadow-sky-200 transition-all active:scale-95"
              >
                <span>Ambil Nomor Antrian Online</span>
                <ArrowRight class="w-4 h-4 text-sky-200" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-4 rounded-2xl bg-white hover:bg-sky-50 border border-sky-200 text-sky-900 font-bold text-sm sm:text-base transition-all shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-sky-500 animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>
          </div>

          <!-- Quick Stats / Feature Card (Col 5) -->
          <div class="lg:col-span-5">
            <div class="p-6 sm:p-8 rounded-3xl bg-white border border-sky-100 shadow-md shadow-sky-100/50 space-y-5">
              <div class="flex items-center justify-between border-b border-sky-50 pb-4">
                <span class="text-xs font-black uppercase tracking-wider text-zinc-400">Poli Anak Hari Ini</span>
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-sky-50 text-sky-800 border border-sky-200">
                  <span class="w-1.5 h-1.5 rounded-full bg-sky-600 animate-pulse"></span>
                  Layanan Buka
                </span>
              </div>

              <!-- Doctors Row -->
              <div class="space-y-3">
                <div class="flex items-center justify-between p-3.5 rounded-2xl bg-sky-50/50 border border-sky-100">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-sky-500 text-white flex items-center justify-center font-bold text-sm shadow-xs">
                      PV
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-950">dr. Patricia Vanessa A., Sp.A</div>
                      <div class="text-[11px] text-zinc-500">Poli Anak 1 • Praktik Pagi (08:30 - 13:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-sky-700 bg-white border border-sky-200 px-2.5 py-1 rounded-lg">
                    Praktik
                  </span>
                </div>

                <div class="flex items-center justify-between p-3.5 rounded-2xl bg-amber-50/50 border border-amber-100">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-amber-400 text-amber-950 flex items-center justify-center font-bold text-sm shadow-xs">
                      HP
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-950">dr. Hendra Pratama, Sp.A</div>
                      <div class="text-[11px] text-zinc-500">Poli Anak 2 • Praktik Sore (14:00 - 18:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-amber-800 bg-white border border-amber-200 px-2.5 py-1 rounded-lg">
                    Sore
                  </span>
                </div>
              </div>

              <div class="pt-2 text-xs text-zinc-600 flex items-center justify-between font-semibold border-t border-sky-50">
                <span>Rata-rata Waktu Tunggu:</span>
                <span class="text-sky-950 font-black">~10 - 15 Menit / Anak</span>
              </div>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue & Multi-Doctor Filter Widget -->
    <section id="live-antrian" class="w-full py-16 sm:py-20 border-b border-sky-100 bg-[#F8FAFC]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="kidinest"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk')"
        />
      </div>
    </section>

    <!-- Section 3: Dokter Spesialis Anak & Profil dr. Patricia -->
    <section id="jadwal-dokter" class="w-full py-16 sm:py-20 border-b border-sky-100 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10 text-center sm:text-left">
          <span class="text-xs font-black uppercase tracking-wider text-sky-800 bg-sky-50 px-3 py-1 rounded-full border border-sky-200">
            Tim Dokter Pediatri
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Dokter Spesialis Anak Kidinest
          </h2>
          <p class="text-sm sm:text-base text-zinc-600 font-medium mt-1">
            Menghadirkan komunikasi yang menenangkan bagi orang tua dan penanganan yang lembut bagi anak.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Doctor 1: dr. Patricia Vanessa A., Sp.A (Instagram Profil Riil) -->
          <div class="p-6 sm:p-8 rounded-3xl border border-sky-200 bg-[#F0F9FF]/40 hover:border-sky-400 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-sky-600 text-white font-black text-2xl flex items-center justify-center shadow-md shadow-sky-200">
                  PV
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Patricia Vanessa A., Sp.A
                  </h3>
                  <p class="text-xs font-bold text-sky-700">
                    Spesialis Ilmu Kesehatan Anak & Vaksinasi Pediatri
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Klinik Utama Semarang Medical Group • Ruang Poli 1
                  </p>
                </div>
              </div>

              <div class="space-y-2.5 py-4 border-t border-b border-sky-100 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Sabtu (08:30 - 13:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Keahlian Klinis:</span>
                  <span class="font-bold text-zinc-950">Imunisasi Lengkap, Nutrisi Bayi & Infeksi Anak</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3.5 rounded-xl bg-sky-600 hover:bg-sky-700 text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Patricia</span>
                <ArrowRight class="w-3.5 h-3.5 text-sky-200" />
              </button>
            </div>
          </div>

          <!-- Doctor 2: dr. Hendra Pratama, Sp.A -->
          <div class="p-6 sm:p-8 rounded-3xl border border-zinc-200 bg-white hover:border-sky-300 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-amber-500 text-white font-black text-2xl flex items-center justify-center shadow-md shadow-amber-200">
                  HP
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Hendra Pratama, Sp.A
                  </h3>
                  <p class="text-xs font-bold text-amber-700">
                    Spesialis Anak • Tumbuh Kembang & Alergi
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Klinik Utama Semarang Medical Group • Ruang Poli 2
                  </p>
                </div>
              </div>

              <div class="space-y-2.5 py-4 border-t border-b border-zinc-100 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Sabtu (14:00 - 18:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Keahlian Klinis:</span>
                  <span class="font-bold text-zinc-950">Milestone Tumbuh Kembang, GTM & Asma</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Hendra</span>
                <ArrowRight class="w-3.5 h-3.5 text-zinc-400" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 4: Layanan Unggulan -->
    <section id="layanan" class="w-full py-16 sm:py-20 border-b border-sky-100 bg-[#FFFDF9]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10">
          <span class="text-xs font-black uppercase tracking-wider text-sky-800 bg-sky-50 px-3 py-1 rounded-full border border-sky-200">
            Layanan Pediatri
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Layanan Kesehatan Lengkap untuk Buah Hati
          </h2>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <div
            v-for="(svc, idx) in services"
            :key="idx"
            class="p-6 rounded-3xl bg-white border border-sky-100 hover:border-sky-300 hover:shadow-md transition-all shadow-2xs flex flex-col justify-between"
          >
            <div>
              <div class="w-12 h-12 rounded-2xl bg-sky-50 text-sky-600 flex items-center justify-center mb-5">
                <component :is="svc.icon" class="w-6 h-6" />
              </div>
              <span class="text-[10px] font-black uppercase tracking-wider text-sky-800 bg-sky-50 px-2.5 py-1 rounded-md">
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
                class="text-xs font-extrabold text-sky-700 hover:text-sky-900 flex items-center gap-1 group"
              >
                <span>Daftar Layanan</span>
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
                Alamat & Kontak Resmi
              </span>
              <h2 class="text-3xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Kidinest Pediatric Care
              </h2>
            </div>

            <div class="p-6 rounded-3xl bg-[#F0F9FF]/40 border border-sky-100 space-y-4 text-xs">
              <div class="flex items-start gap-3">
                <MapPin class="w-5 h-5 text-sky-600 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Lokasi Klinik:</div>
                  <div class="text-zinc-600 mt-0.5">
                    Klinik Utama Semarang Medical Group<br />
                    Jl. Melati Utara No. 3, Kota Semarang, Jawa Tengah
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-sky-100 pt-3">
                <Clock class="w-5 h-5 text-sky-600 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Jam Praktik Poli Anak:</div>
                  <div class="text-zinc-600 mt-0.5 space-y-0.5">
                    <div>Poli Pagi: 08:30 - 13:00 WIB (dr. Patricia Vanessa, Sp.A)</div>
                    <div>Poli Sore: 14:00 - 18:00 WIB (dr. Hendra Pratama, Sp.A)</div>
                    <div class="text-[11px] text-sky-700 font-bold">Minggu & Hari Libur: Dengan Perjanjian</div>
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-sky-100 pt-3">
                <Phone class="w-5 h-5 text-sky-600 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">WhatsApp Care & Info Dokter:</div>
                  <div class="text-zinc-600 mt-0.5 font-mono">
                    +62 812-3456-7890 (Customer Service Kidinest)
                  </div>
                </div>
              </div>
            </div>

            <!-- Switch to SSC Banner -->
            <div class="p-5 rounded-3xl bg-zinc-950 text-white flex items-center justify-between">
              <div>
                <div class="text-xs font-bold text-zinc-400">Butuh konsultasi kulit & estetika?</div>
                <div class="text-sm font-black text-rose-300">Semarang Skin Center (SSC)</div>
              </div>
              <button
                @click="emit('navigate', 'ssc')"
                class="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/20 text-white text-xs font-bold transition-colors"
              >
                Kunjungi SSC →
              </button>
            </div>
          </div>

          <!-- Right: FAQ Section (Col 7) -->
          <div class="lg:col-span-7 space-y-4">
            <div class="mb-4">
              <span class="text-xs font-black uppercase tracking-wider text-sky-800 bg-sky-50 px-3 py-1 rounded-full border border-sky-200">
                Informasi Penting
              </span>
              <h3 class="text-2xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Pertanyaan yang Sering Diajukan
              </h3>
            </div>

            <div class="space-y-3">
              <div
                v-for="(faq, idx) in faqs"
                :key="idx"
                class="p-5 rounded-2xl border border-sky-100 bg-[#F0F9FF]/20"
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
    <footer class="w-full border-t border-sky-100 py-8 bg-[#FFFDF9] text-xs text-zinc-500">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          © 2026 <strong>Kidinest Pediatric Care</strong> • Bagian dari Semarang Medical Group.
        </div>
        <div class="flex items-center gap-4 font-bold text-zinc-700">
          <button @click="emit('navigate', 'kiosk')" class="hover:text-zinc-950">Portal Pasien</button>
          <span>•</span>
          <button @click="emit('navigate', 'ssc')" class="hover:text-zinc-950">Semarang Skin Center</button>
          <span>•</span>
          <button @click="emit('navigate', 'login')" class="hover:text-zinc-950">Login Staf Medis</button>
        </div>
      </div>
    </footer>

  </div>
</template>
