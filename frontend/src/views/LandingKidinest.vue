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
  ChevronRight
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login'): void
}>()

const doctors: DoctorInfo[] = [
  {
    id: 'hendra',
    name: 'dr. Hendra Pratama, Sp.A',
    specialty: 'Spesialis Anak & Pediatri Umum',
    room: 'Ruang Poli Anak 1',
    status: 'active'
  },
  {
    id: 'anita',
    name: 'dr. Anita Ratna, Sp.A',
    specialty: 'Konsultan Tumbuh Kembang & Nutrisi Anak',
    room: 'Ruang Poli Anak 2',
    status: 'active'
  }
]

const services = [
  {
    icon: Syringe,
    title: 'Vaksinasi & Imunisasi Lengkap',
    desc: 'Vaksinasi wajib & tambahan sesuai rekomendasi IDAI dengan rantai dingin (cold-chain) terjamin.',
    tag: 'Cold-Chain Certified'
  },
  {
    icon: Smile,
    title: 'Skrining Tumbuh Kembang',
    desc: 'Pemantauan milestone motorik, sensorik, dan wicara bersama dokter spesialis anak secara berkala.',
    tag: 'Milestone IDAI'
  },
  {
    icon: Activity,
    title: 'Pemeriksaan Sakit & Infeksi',
    desc: 'Penanganan cepat demam, flu, batuk, pilek, diare, dan gangguan pencernaan anak dengan resep presisi.',
    tag: 'Fast Response'
  },
  {
    icon: HeartPulse,
    title: 'Fisioterapi & Nebulisasi Anak',
    desc: 'Terapi inhalasi uap modern dan stimulasi postur tubuh anak di ruang ramah anak yang menyenangkan.',
    tag: 'Kid-Friendly Room'
  }
]

const faqs = [
  {
    q: 'Apakah bisa mendaftar antrian secara online dari rumah?',
    a: 'Bisa! Anda cukup menekan tombol "Ambil Nomor Antrian Online", pilih dokter dan jam kunjungan, lalu tiket antrian akan otomatis masuk ke ponsel Anda via WhatsApp.'
  },
  {
    q: 'Bagaimana jika saya datang terlambat saat nomor dipanggil?',
    a: 'Sistem Kidinest memberikan jeda toleransi 2 nomor berikutnya tanpa menghanguskan tiket. Anda akan dipanggil kembali pada giliran jeda dokter.'
  },
  {
    q: 'Apakah resep obat Kidinest langsung terhubung ke Apotek?',
    a: 'Ya, resep elektronik terhubung langsung ke Apotek Farmasi internal. Pasien tidak perlu membawa kertas resep, cukup tunggu obat dipanggil di kasir.'
  }
]

function handleDirectRegister(docName?: string) {
  emit('navigate', 'kiosk')
}
</script>

<template>
  <div class="min-h-screen bg-[#FBFBF9] text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300 flex flex-col">

    <!-- Sticky Dual-Brand Header -->
    <BrandSwitcherHeader
      current-brand="kidinest"
      @navigate="(v) => emit('navigate', v)"
    />

    <!-- Hero Section (Action-Oriented, Anti-Slop, Clean Nordic Warmth) -->
    <section class="w-full border-b border-emerald-950/10 py-12 sm:py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Category Pill Badge -->
        <div class="mb-6">
          <span class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-[#0F4C3A] text-amber-300 shadow-xs">
            <Baby class="w-4 h-4 text-amber-300" />
            Poli Spesialis Anak & Tumbuh Kembang • Semarang
          </span>
        </div>

        <!-- Headline & Direct Description Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center mb-10">

          <!-- Headline & Primary Action (Col 7) -->
          <div class="lg:col-span-7 space-y-6">
            <h1 class="text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-zinc-950 leading-[1.05]">
              Pemeriksaan Buah Hati Tanpa Antri Lama.
            </h1>
            <p class="text-base sm:text-lg text-zinc-700 leading-relaxed max-w-2xl font-medium">
              Pantau antrian dokter secara live dari ponsel Anda, ambil tiket mandiri tanpa ribet, dan dapatkan notifikasi WhatsApp otomatis sebelum giliran masuk. Ramah anak, aman, dan tanpa desak-desakan.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-2">
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center justify-center gap-2.5 px-6 py-3.5 rounded-2xl bg-[#0F4C3A] hover:bg-[#0A3629] text-white font-extrabold text-sm sm:text-base shadow-lg transition-all active:scale-95"
              >
                <span>Ambil Nomor Antrian Online</span>
                <ArrowRight class="w-4 h-4 text-amber-300" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-3.5 rounded-2xl bg-white hover:bg-zinc-100 border border-zinc-300 text-zinc-800 font-bold text-sm sm:text-base transition-all"
              >
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>
          </div>

          <!-- Quick Stats / Feature Visual (Col 5) -->
          <div class="lg:col-span-5">
            <div class="p-6 sm:p-8 rounded-3xl bg-white border border-emerald-950/10 shadow-sm space-y-5">
              <div class="flex items-center justify-between border-b border-zinc-100 pb-4">
                <span class="text-xs font-black uppercase tracking-wider text-zinc-400">Jadwal Poli Hari Ini</span>
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-100 text-emerald-800">
                  <span class="w-1.5 h-1.5 rounded-full bg-emerald-600 animate-pulse"></span>
                  Poli Buka
                </span>
              </div>

              <!-- Doctors Today Quick Row -->
              <div class="space-y-3">
                <div class="flex items-center justify-between p-3 rounded-xl bg-zinc-50 border border-zinc-200/70">
                  <div class="flex items-center gap-3">
                    <div class="w-9 h-9 rounded-lg bg-emerald-100 text-emerald-800 flex items-center justify-center font-bold text-xs">
                      1
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-900">dr. Hendra Pratama, Sp.A</div>
                      <div class="text-[11px] text-zinc-500">Poli Anak 1 • Pagi (08:00 - 13:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-emerald-700 bg-emerald-50 px-2 py-0.5 rounded-md">
                    Aktif
                  </span>
                </div>

                <div class="flex items-center justify-between p-3 rounded-xl bg-zinc-50 border border-zinc-200/70">
                  <div class="flex items-center gap-3">
                    <div class="w-9 h-9 rounded-lg bg-emerald-100 text-emerald-800 flex items-center justify-center font-bold text-xs">
                      2
                    </div>
                    <div>
                      <div class="text-xs font-extrabold text-zinc-900">dr. Anita Ratna, Sp.A</div>
                      <div class="text-[11px] text-zinc-500">Poli Anak 2 • Siang (13:00 - 18:00)</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-black text-amber-800 bg-amber-50 px-2 py-0.5 rounded-md">
                    Siang
                  </span>
                </div>
              </div>

              <div class="pt-2 text-xs text-zinc-600 flex items-center justify-between font-semibold">
                <span>Rata-rata Waktu Tunggu:</span>
                <span class="text-zinc-950 font-black">~12 Menit / Pasien</span>
              </div>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue & Multi-Doctor Filter Widget -->
    <section id="live-antrian" class="w-full py-16 sm:py-20 border-b border-emerald-950/10 bg-[#F6F6F2]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="kidinest"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk')"
        />
      </div>
    </section>

    <!-- Section 3: Dokter Spesialis Kidinest & Jadwal -->
    <section id="jadwal-dokter" class="w-full py-16 sm:py-20 border-b border-emerald-950/10 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10 text-center sm:text-left">
          <span class="text-xs font-black uppercase tracking-wider text-emerald-800 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200">
            Tim Medis Pediatri
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Dokter Spesialis Anak Kami
          </h2>
          <p class="text-sm sm:text-base text-zinc-600 font-medium mt-1">
            Ditangani oleh dokter spesialis berpengalaman dengan pendekatan ramah anak yang membuat si kecil nyaman.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <!-- Doctor 1 -->
          <div class="p-6 sm:p-8 rounded-3xl border border-zinc-200 bg-[#FBFBF9] hover:border-emerald-700/40 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-[#0F4C3A] text-amber-300 font-black text-2xl flex items-center justify-center shadow-xs">
                  HP
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Hendra Pratama, Sp.A
                  </h3>
                  <p class="text-xs font-bold text-emerald-800">
                    Spesialis Ilmu Kesehatan Anak & Vaksinasi
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Ruang Poli Anak 1 • SIP: 446.1/082/SIP/2021
                  </p>
                </div>
              </div>

              <div class="space-y-2 py-4 border-t border-b border-zinc-200 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Sabtu (08:00 - 13:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-zinc-950">Imunisasi, Batuk/Flu, Alergi Anak</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#0F4C3A] hover:bg-[#0A3629] text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Hendra</span>
                <ArrowRight class="w-3.5 h-3.5 text-amber-300" />
              </button>
            </div>
          </div>

          <!-- Doctor 2 -->
          <div class="p-6 sm:p-8 rounded-3xl border border-zinc-200 bg-[#FBFBF9] hover:border-emerald-700/40 transition-all shadow-xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-16 h-16 rounded-2xl bg-emerald-900 text-emerald-200 font-black text-2xl flex items-center justify-center shadow-xs">
                  AR
                </div>
                <div>
                  <h3 class="text-xl font-extrabold text-zinc-950">
                    dr. Anita Ratna, Sp.A
                  </h3>
                  <p class="text-xs font-bold text-emerald-800">
                    Konsultan Tumbuh Kembang, Nutrisi & Fisioterapi
                  </p>
                  <p class="text-[11px] text-zinc-500 font-medium">
                    Ruang Poli Anak 2 • SIP: 446.1/119/SIP/2022
                  </p>
                </div>
              </div>

              <div class="space-y-2 py-4 border-t border-b border-zinc-200 text-xs text-zinc-700 font-medium">
                <div class="flex items-center justify-between">
                  <span>Jadwal Praktik:</span>
                  <span class="font-bold text-zinc-950">Senin - Sabtu (13:00 - 18:00)</span>
                </div>
                <div class="flex items-center justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-zinc-950">Milestone Perkembangan, GTM & Berat Badan</span>
                </div>
              </div>
            </div>

            <div class="pt-6">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#0F4C3A] hover:bg-[#0A3629] text-white text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-sm"
              >
                <span>Daftar Konsultasi ke dr. Anita</span>
                <ArrowRight class="w-3.5 h-3.5 text-amber-300" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 4: Layanan Unggulan -->
    <section id="layanan" class="w-full py-16 sm:py-20 border-b border-emerald-950/10 bg-[#FBFBF9]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10">
          <span class="text-xs font-black uppercase tracking-wider text-emerald-800 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200">
            Layanan Terpadu
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-zinc-950 tracking-tight mt-3">
            Perawatan Menyeluruh untuk Tumbuh Kembang Si Kecil
          </h2>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <div
            v-for="(svc, idx) in services"
            :key="idx"
            class="p-6 rounded-3xl bg-white border border-zinc-200 hover:border-zinc-400 transition-all shadow-2xs flex flex-col justify-between"
          >
            <div>
              <div class="w-12 h-12 rounded-2xl bg-emerald-50 text-emerald-800 flex items-center justify-center mb-5">
                <component :is="svc.icon" class="w-6 h-6" />
              </div>
              <span class="text-[10px] font-black uppercase tracking-wider text-amber-800 bg-amber-50 px-2 py-0.5 rounded-md">
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
                class="text-xs font-extrabold text-emerald-800 hover:text-emerald-950 flex items-center gap-1 group"
              >
                <span>Daftar Layanan</span>
                <ChevronRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 5: Lokasi, Jam Buka Poli & FAQ -->
    <section id="lokasi" class="w-full py-16 sm:py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">

          <!-- Left: Location & Hours (Col 5) -->
          <div class="lg:col-span-5 space-y-6">
            <div>
              <span class="text-xs font-black uppercase tracking-wider text-zinc-500">
                Informasi Kunjungan
              </span>
              <h2 class="text-3xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Klinik Kidinest Semarang
              </h2>
            </div>

            <div class="p-6 rounded-3xl bg-[#FBFBF9] border border-zinc-200 space-y-4 text-xs">
              <div class="flex items-start gap-3">
                <MapPin class="w-5 h-5 text-emerald-800 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Alamat Klinik:</div>
                  <div class="text-zinc-600 mt-0.5">
                    Jl. Pandanaran No. 88, Mugassari, Semarang Selatan, Kota Semarang
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-zinc-200 pt-3">
                <Clock class="w-5 h-5 text-emerald-800 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">Jam Operasional Poli:</div>
                  <div class="text-zinc-600 mt-0.5 space-y-0.5">
                    <div>Poli Pagi: 08:00 - 13:00 WIB</div>
                    <div>Poli Sore: 13:00 - 18:00 WIB</div>
                    <div class="text-[11px] text-emerald-700 font-bold">Minggu & Libur Nasional: Dengan Perjanjian</div>
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-3 border-t border-zinc-200 pt-3">
                <Phone class="w-5 h-5 text-emerald-800 shrink-0 mt-0.5" />
                <div>
                  <div class="font-extrabold text-zinc-950 text-sm">WhatsApp Care & Konsultasi:</div>
                  <div class="text-zinc-600 mt-0.5">
                    +62 812-3456-7890 (Customer Care 24 Jam)
                  </div>
                </div>
              </div>
            </div>

            <!-- Switch to SSC Banner -->
            <div class="p-5 rounded-2xl bg-zinc-950 text-white flex items-center justify-between">
              <div>
                <div class="text-xs font-bold text-zinc-400">Mencari perawatan kulit & estetika?</div>
                <div class="text-sm font-black text-rose-300">Semarang Skin Center (SSC)</div>
              </div>
              <button
                @click="emit('navigate', 'ssc')"
                class="px-3.5 py-1.5 rounded-xl bg-white/10 hover:bg-white/20 text-white text-xs font-bold transition-colors"
              >
                Lihat SSC →
              </button>
            </div>
          </div>

          <!-- Right: FAQ Section (Col 7) -->
          <div class="lg:col-span-7 space-y-4">
            <div class="mb-4">
              <span class="text-xs font-black uppercase tracking-wider text-emerald-800 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200">
                Tanya Jawab Pasien
              </span>
              <h3 class="text-2xl font-extrabold text-zinc-950 tracking-tight mt-2">
                Pertanyaan Seputar Antrian Kidinest
              </h3>
            </div>

            <div class="space-y-3">
              <div
                v-for="(faq, idx) in faqs"
                :key="idx"
                class="p-5 rounded-2xl border border-zinc-200 bg-[#FBFBF9]"
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
    <footer class="w-full border-t border-zinc-200 py-8 bg-[#FBFBF9] text-xs text-zinc-500">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          © 2026 <strong>Kidinest Poli Anak</strong> • Smart Clinic Management System. Seluruh hak cipta dilindungi.
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
