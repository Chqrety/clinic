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

// Educational Benefits from Official Feed (Benefits of Baby Care)
const massageBenefits = [
  { title: 'Brain Growth', desc: 'Stimulasi sensori taktil untuk mengoptimalkan koneksi sinaps otak bayi.' },
  { title: 'Congestion Relief', desc: 'Membantu melegakan saluran napas si kecil saat batuk & pilek.' },
  { title: 'Improved Digestion', desc: 'Mengurangi kembung, kolik perut, dan melancarkan pencernaan.' },
  { title: 'Restful Sleep', desc: 'Merelaksasikan otot agar bayi tidur lebih lelap dan berkualitas.' },
  { title: 'Parent Bonding', desc: 'Mempererat ikatan emosional dan rasa aman antara orang tua dan anak.' },
  { title: 'Pain Relief', desc: 'Menenangkan rasa nyeri saat tumbuh gigi atau masa imunisasi.' }
]

// Core Pediatric Services (Hairline Editorial Index)
const servicesList = [
  {
    num: '01',
    title: 'Vaksinasi & Imunisasi Lengkap IDAI',
    desc: 'Vaksin wajib & tambahan dengan cold-chain medis terstandar internasional.',
    tag: 'Cold-Chain Certified',
    doctor: 'dr. Patricia Vanessa, Sp.A'
  },
  {
    num: '02',
    title: 'Skrining Milestone & Tumbuh Kembang',
    desc: 'Pemantauan grafik pertumbuhan fisik, motorik, bahasa, dan nutrisi gizi anak.',
    tag: 'Milestone IDAI',
    doctor: 'dr. Hendra Pratama, Sp.A'
  },
  {
    num: '03',
    title: 'Pemeriksaan Sakit, Infeksi & Alergi',
    desc: 'Diagnosis ramah anak untuk demam, flu, batuk, muntah, diare, dan dermatitis.',
    tag: 'Responsif & Lembut',
    doctor: 'dr. Patricia Vanessa, Sp.A'
  },
  {
    num: '04',
    title: 'Nebulisasi & Stimulasi Fisioterapi',
    desc: 'Terapi uap inhalasi dan stimulasi fisik di ruang periksa ramah anak yang ceria.',
    tag: 'Kid-Friendly Care',
    doctor: 'dr. Hendra Pratama, Sp.A'
  }
]

const faqs = [
  {
    q: 'Bagaimana cara mengambil nomor antrian online Kidinest?',
    a: 'Cukup klik tombol "Ambil Antrian Online", pilih tanggal dan dokter spesialis anak, lalu tiket digital langsung tersimpan di ponsel dan terhubung ke WhatsApp Anda.'
  },
  {
    q: 'Apakah Kidinest melayani imunisasi lengkap bayi baru lahir?',
    a: 'Ya, kami melayani seluruh paket imunisasi IDAI lengkap dengan rantai dingin (cold chain) vaksin yang terjaga ketat demi efektivitas dan keamanan buah hati.'
  },
  {
    q: 'Di mana alamat Kidinest Pediatric Care?',
    a: 'Kidinest berlokasi di dalam Klinik Utama Semarang Medical Group, Jl. Melati Utara No. 3, Semarang.'
  }
]
</script>

<template>
  <div class="min-h-screen bg-[#FAF7F2] text-[#0C3B36] font-['Plus_Jakarta_Sans',sans-serif] selection:bg-[#E5A93C]/30 flex flex-col">

    <!-- Sticky Dual-Brand Header -->
    <BrandSwitcherHeader
      current-brand="kidinest"
      @navigate="(v) => emit('navigate', v)"
    />

    <!-- Hero Section (Friendly Warmth, Scale Humanis, Inspirasi Feed Kidinest) -->
    <section class="w-full border-b border-[#E5A93C]/20 py-12 sm:py-16 bg-gradient-to-b from-[#FAF7F2] via-[#FAF7F2] to-[#F5EFE6]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Top Pill Tag (Mustard Badge Khas Feed) -->
        <div class="mb-5">
          <span class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-bold bg-[#E5A93C] text-[#0C3B36] shadow-2xs">
            <img src="/assets/kidinest/logo-maskot.png" alt="Baby Bean" class="w-4 h-4 object-contain" />
            <span>Kidinest Pediatric Care • Semarang Medical Group</span>
          </span>
        </div>

        <!-- Headline & Action Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-center">

          <!-- Left: Narasi Hangat & Proporsional (Col 7) -->
          <div class="lg:col-span-7 space-y-5">
            <h1 class="text-3xl sm:text-5xl font-extrabold text-[#0C3B36] tracking-tight leading-[1.15]">
              Pemeriksaan Buah Hati dengan Hangat & Tenang.
            </h1>
            <p class="text-sm sm:text-base text-[#0C3B36]/80 leading-relaxed max-w-xl font-medium">
              Pelayanan dokter spesialis anak bersama <strong>dr. Patricia Vanessa A., Sp.A</strong>. Pantau nomor antrian secara live dari ponsel Anda, hindari antrean panjang di ruang tunggu, dan dapatkan notifikasi WhatsApp saat giliran tiba.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-2">
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center justify-center gap-2.5 px-6 py-3.5 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-white font-bold text-xs sm:text-sm uppercase tracking-wider shadow-md shadow-[#0C3B36]/20 transition-all active:scale-95"
              >
                <span>Ambil Nomor Antrian Online</span>
                <ArrowRight class="w-4 h-4 text-[#E5A93C]" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-3.5 rounded-xl bg-white hover:bg-white/80 border border-[#E5A93C]/40 text-[#0C3B36] font-bold text-xs sm:text-sm transition-all shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-[#E5A93C] animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>

            <!-- Mini Trust Strip -->
            <div class="pt-3 flex items-center gap-6 text-xs text-[#0C3B36]/75 font-semibold">
              <div class="flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-[#E5A93C]" />
                <span>Dokter Spesialis Anak Berlisensi</span>
              </div>
              <div class="flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-[#E5A93C]" />
                <span>Notifikasi WhatsApp Paperless</span>
              </div>
            </div>
          </div>

          <!-- Right: Visual Card Inspirasi Feed Asli (Col 5) -->
          <div class="lg:col-span-5">
            <div class="p-6 rounded-3xl bg-white/90 border border-[#E5A93C]/30 shadow-sm relative overflow-hidden">
              
              <!-- Header Visual Feed -->
              <div class="flex items-center justify-between pb-4 border-b border-zinc-100">
                <div class="flex items-center gap-2">
                  <div class="w-7 h-7 rounded-lg bg-[#E5A93C]/20 flex items-center justify-center">
                    <img src="/assets/kidinest/logo-maskot.png" class="w-5 h-5 object-contain" />
                  </div>
                  <span class="font-bold text-xs text-[#0C3B36]">Milestone & Tumbuh Kembang</span>
                </div>
                <span class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase bg-[#0C3B36] text-white">
                  Panduan IDAI
                </span>
              </div>

              <!-- Interactive Milestone Showcase -->
              <div class="py-4 space-y-3">
                <div class="flex items-center justify-between p-3 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/25">
                  <div class="flex items-center gap-3">
                    <span class="w-7 h-7 rounded-full bg-[#E5A93C] text-[#0C3B36] font-extrabold text-xs flex items-center justify-center">
                      12m
                    </span>
                    <div>
                      <div class="text-xs font-bold text-zinc-950">Kosakata Pertama (1-5 Kata)</div>
                      <div class="text-[11px] text-zinc-500 font-medium">Mama, Papa, Daa, Minta</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-bold text-[#0C3B36]">Normal</span>
                </div>

                <div class="flex items-center justify-between p-3 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/25">
                  <div class="flex items-center gap-3">
                    <span class="w-7 h-7 rounded-full bg-[#0C3B36] text-white font-extrabold text-xs flex items-center justify-center">
                      24m
                    </span>
                    <div>
                      <div class="text-xs font-bold text-zinc-950">Kalimat 2 Kata (50-300 Kata)</div>
                      <div class="text-[11px] text-zinc-500 font-medium">Mau minum, Naik mobil</div>
                    </div>
                  </div>
                  <span class="text-[11px] font-bold text-[#0C3B36]">Aktif</span>
                </div>
              </div>

              <!-- Doctor Signature Info -->
              <div class="pt-3 border-t border-zinc-100 flex items-center justify-between text-xs text-[#0C3B36]">
                <div class="font-medium">
                  Konsultasi: <strong>dr. Patricia Vanessa A., Sp.A</strong>
                </div>
                <button
                  @click="emit('navigate', 'kiosk')"
                  class="text-xs font-bold text-[#E5A93C] hover:text-[#0C3B36] flex items-center gap-1"
                >
                  <span>Daftar →</span>
                </button>
              </div>

            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue Strip (Menyatu Alami dengan Halaman) -->
    <section id="live-antrian" class="w-full py-12 sm:py-16 border-b border-[#E5A93C]/20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="kidinest"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk')"
        />
      </div>
    </section>

    <!-- Section 3: Layanan Pediatri (Editorial Index Table - Bukan Card Berjejer!) -->
    <section id="layanan" class="w-full py-12 sm:py-16 border-b border-[#E5A93C]/20 bg-[#FAF7F2]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Section Heading Proporsional -->
        <div class="mb-8 flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div>
            <span class="text-xs font-bold uppercase tracking-wider text-[#0C3B36]/70">
              Layanan Spesialis
            </span>
            <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-1">
              Menu Pelayanan & Vaksinasi Buah Hati
            </h2>
          </div>
          <button
            @click="emit('navigate', 'kiosk')"
            class="text-xs font-bold uppercase tracking-wider text-[#0C3B36] underline hover:text-[#E5A93C] transition-colors"
          >
            Reservasi Semua Layanan →
          </button>
        </div>

        <!-- Editorial Index Table (Hairline Dividers, Bebas Card) -->
        <div class="border-t border-b border-[#0C3B36]/15 divide-y divide-[#0C3B36]/10">
          <div
            v-for="item in servicesList"
            :key="item.num"
            class="py-5 sm:py-6 flex flex-col md:flex-row md:items-center justify-between gap-4 group hover:bg-white/40 transition-colors px-2 sm:px-4 rounded-xl"
          >
            <div class="flex items-start sm:items-center gap-5">
              <span class="font-mono text-sm font-bold text-[#E5A93C] shrink-0 pt-0.5 sm:pt-0">
                {{ item.num }}
              </span>
              <div>
                <div class="flex items-center gap-3">
                  <h3 class="text-base sm:text-lg font-bold text-[#0C3B36] group-hover:text-[#134E48] transition-colors">
                    {{ item.title }}
                  </h3>
                  <span class="hidden sm:inline-block px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-[#E5A93C]/20 text-[#0C3B36]">
                    {{ item.tag }}
                  </span>
                </div>
                <p class="text-xs sm:text-sm text-[#0C3B36]/75 mt-0.5 font-medium">
                  {{ item.desc }}
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between sm:justify-end gap-5 shrink-0 pl-10 md:pl-0">
              <span class="text-xs font-semibold text-[#0C3B36]/70">
                {{ item.doctor }}
              </span>
              <button
                @click="emit('navigate', 'kiosk')"
                class="inline-flex items-center gap-1.5 text-xs font-bold text-[#0C3B36] hover:text-[#E5A93C] transition-colors"
              >
                <span>Daftar</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 4: Panduan Edukasi Bayi (Mengadopsi Layout Feed "Benefits for Babies") -->
    <section class="w-full py-12 sm:py-16 border-b border-[#E5A93C]/20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-10 text-center max-w-2xl mx-auto">
          <span class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold bg-[#E5A93C] text-[#0C3B36]">
            Tips Edukasi Dokter Spesialis
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-3">
            Manfaat Stimulasi & Terapi Lembut untuk Bayi
          </h2>
          <p class="text-xs sm:text-sm text-[#0C3B36]/70 mt-1 font-medium">
            Mendukung tumbuh kembang optimal dengan pendekatan medis yang menyenangkan di rumah dan klinik.
          </p>
        </div>

        <!-- 6 Benefit Points Berbaris Rapi (Bukan Card Tebal) -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
          <div
            v-for="(b, idx) in massageBenefits"
            :key="idx"
            class="p-5 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/20 flex items-start gap-3.5"
          >
            <div class="w-8 h-8 rounded-full bg-[#E5A93C]/30 text-[#0C3B36] font-bold text-xs flex items-center justify-center shrink-0">
              0{{ idx + 1 }}
            </div>
            <div>
              <h4 class="text-sm font-bold text-[#0C3B36]">{{ b.title }}</h4>
              <p class="text-xs text-[#0C3B36]/70 mt-1 leading-relaxed font-medium">
                {{ b.desc }}
              </p>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 5: Jadwal Dokter Spesialis Anak -->
    <section id="jadwal-dokter" class="w-full py-12 sm:py-16 border-b border-[#E5A93C]/20 bg-[#FAF7F2]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-8">
          <span class="text-xs font-bold uppercase tracking-wider text-[#0C3B36]/70">
            Jadwal Praktik
          </span>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-1">
            Dokter Spesialis Anak Kidinest
          </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Dokter 1: dr. Patricia Vanessa A., Sp.A -->
          <div class="p-6 rounded-3xl bg-white border border-[#E5A93C]/30 shadow-2xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-14 h-14 rounded-2xl bg-[#0C3B36] text-[#E5A93C] font-black text-xl flex items-center justify-center shadow-xs">
                  PV
                </div>
                <div>
                  <h3 class="text-lg font-bold text-[#0C3B36]">
                    dr. Patricia Vanessa A., Sp.A
                  </h3>
                  <p class="text-xs font-semibold text-[#E5A93C] mt-0.5">
                    Spesialis Ilmu Kesehatan Anak & Vaksinasi Pediatri
                  </p>
                  <p class="text-[11px] text-[#0C3B36]/60">
                    Ruang Poli Anak 1 • Semarang Medical Group
                  </p>
                </div>
              </div>

              <div class="py-3 border-t border-b border-zinc-100 space-y-1.5 text-xs text-[#0C3B36]/80 font-medium">
                <div class="flex justify-between">
                  <span>Jadwal Praktik Pagi:</span>
                  <span class="font-bold text-[#0C3B36]">Senin - Sabtu (08:30 - 13:00)</span>
                </div>
                <div class="flex justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-[#0C3B36]">Vaksinasi, Nutrisi & Infeksi Anak</span>
                </div>
              </div>
            </div>

            <div class="pt-5">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-white text-xs font-bold uppercase tracking-wider transition-all"
              >
                Daftar Konsultasi ke dr. Patricia →
              </button>
            </div>
          </div>

          <!-- Dokter 2: dr. Hendra Pratama, Sp.A -->
          <div class="p-6 rounded-3xl bg-white border border-[#E5A93C]/30 shadow-2xs flex flex-col justify-between">
            <div>
              <div class="flex items-center gap-4 mb-4">
                <div class="w-14 h-14 rounded-2xl bg-[#E5A93C] text-[#0C3B36] font-black text-xl flex items-center justify-center shadow-xs">
                  HP
                </div>
                <div>
                  <h3 class="text-lg font-bold text-[#0C3B36]">
                    dr. Hendra Pratama, Sp.A
                  </h3>
                  <p class="text-xs font-semibold text-[#0C3B36] mt-0.5">
                    Spesialis Anak • Konsultan Tumbuh Kembang & Alergi
                  </p>
                  <p class="text-[11px] text-[#0C3B36]/60">
                    Ruang Poli Anak 2 • Semarang Medical Group
                  </p>
                </div>
              </div>

              <div class="py-3 border-t border-b border-zinc-100 space-y-1.5 text-xs text-[#0C3B36]/80 font-medium">
                <div class="flex justify-between">
                  <span>Jadwal Praktik Sore:</span>
                  <span class="font-bold text-[#0C3B36]">Senin - Sabtu (14:00 - 18:00)</span>
                </div>
                <div class="flex justify-between">
                  <span>Fokus Penanganan:</span>
                  <span class="font-bold text-[#0C3B36]">Milestone Perkembangan & Asma</span>
                </div>
              </div>
            </div>

            <div class="pt-5">
              <button
                @click="emit('navigate', 'kiosk')"
                class="w-full py-3 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-white text-xs font-bold uppercase tracking-wider transition-all"
              >
                Daftar Konsultasi ke dr. Hendra →
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
              <span class="text-xs font-bold uppercase tracking-wider text-[#0C3B36]/70">
                Akses & Janji Temu
              </span>
              <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-1">
                Kidinest Pediatric Care
              </h2>
            </div>

            <div class="p-5 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/20 space-y-3 text-xs text-[#0C3B36]">
              <div class="flex items-start gap-2.5">
                <MapPin class="w-4 h-4 text-[#E5A93C] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">Alamat Resmi:</div>
                  <div class="opacity-80 mt-0.5">
                    Klinik Utama Semarang Medical Group<br />
                    Jl. Melati Utara No. 3, Kota Semarang, Jawa Tengah
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-2.5 border-t border-zinc-200/60 pt-2.5">
                <Clock class="w-4 h-4 text-[#E5A93C] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">Jam Layanan Poli:</div>
                  <div class="opacity-80 mt-0.5">
                    Pagi: 08:30 - 13:00 WIB • Sore: 14:00 - 18:00 WIB
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-2.5 border-t border-zinc-200/60 pt-2.5">
                <Phone class="w-4 h-4 text-[#E5A93C] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">WhatsApp Reservasi & Dokter:</div>
                  <div class="opacity-80 mt-0.5 font-mono">
                    +62 812-3456-7890
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right: FAQ (Col 7) -->
          <div class="lg:col-span-7 space-y-3">
            <h3 class="text-lg font-bold text-[#0C3B36] mb-3">
              Pertanyaan yang Sering Diajukan
            </h3>

            <div
              v-for="(faq, idx) in faqs"
              :key="idx"
              class="p-4 rounded-xl border border-zinc-200/80 bg-[#FAF7F2]/50"
            >
              <h4 class="text-xs sm:text-sm font-bold text-[#0C3B36] mb-1">
                {{ faq.q }}
              </h4>
              <p class="text-xs text-[#0C3B36]/75 leading-relaxed font-medium">
                {{ faq.a }}
              </p>
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Footer Solid Pita Hijau Botol Khas Feed Kidinest -->
    <footer class="w-full bg-[#0C3B36] text-white py-6 text-xs select-none">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 flex flex-col sm:flex-row items-center justify-between gap-3">
        <div class="flex items-center gap-2">
          <img src="/assets/kidinest/logo-maskot.png" class="w-5 h-5 object-contain" />
          <span class="font-bold">Kidinest Pediatric Care • Semarang Medical Group</span>
        </div>
        <div class="font-medium text-white/80">
          Jl. Melati Utara No 3, Semarang
        </div>
        <div class="flex items-center gap-4 text-white/70 font-semibold">
          <button @click="emit('navigate', 'kiosk')" class="hover:text-white">Portal Antrian</button>
          <span>•</span>
          <button @click="emit('navigate', 'ssc')" class="hover:text-white">Semarang Skin Center</button>
          <span>•</span>
          <button @click="emit('navigate', 'login')" class="hover:text-white">Login Staf</button>
        </div>
      </div>
    </footer>

  </div>
</template>
