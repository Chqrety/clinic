<script setup lang="ts">
import { ref } from 'vue'
import BrandSwitcherHeader from '../components/BrandSwitcherHeader.vue'
import LiveQueueFilter, { type DoctorInfo } from '../components/LiveQueueFilter.vue'
import {
  Clock,
  MapPin,
  Phone,
  ArrowRight
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: string): void
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
    a: 'Cukup klik tombol "Ambil Antrian Online", pilih dokter dan jam kunjungan, lalu tiket antrian digital langsung tersimpan di ponsel Anda.'
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

    <!-- Hero Section: Clean Editorial Healthcare (Bebas Kotak AI) -->
    <section class="w-full border-b border-[#E5A93C]/20 pt-12 sm:pt-16 pb-10 sm:pb-12 bg-gradient-to-b from-[#FAF7F2] to-[#F5EFE6]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Top Pill Tag (Mustard Badge Khas Feed) -->
        <div class="mb-4">
          <span class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-bold bg-[#E5A93C] text-[#0C3B36]">
            <img src="/assets/kidinest/logo-maskot.png" alt="Baby Bean" class="w-4 h-4 object-contain" />
            <span>Klinik Spesialis Anak • Semarang Medical Group</span>
          </span>
        </div>

        <!-- Headline & Typographic Focus Grid (Murni Tipografi & Layout Hairline, Bebas Foto & Bebas Card) -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start">

          <!-- Left: Narasi & Tindakan Utama (Col 7) -->
          <div class="lg:col-span-7 space-y-5">
            <h1 class="text-3xl sm:text-5xl font-extrabold text-[#0C3B36] tracking-tight leading-[1.15]">
              Pemeriksaan Buah Hati yang Tenang, Cermat, dan Hangat.
            </h1>
            <p class="text-sm sm:text-base text-[#0C3B36]/80 leading-relaxed font-medium max-w-xl">
              Layanan dokter spesialis anak berstandar IDAI bersama <strong>dr. Patricia Vanessa A., Sp.A</strong> dan <strong>dr. Hendra Pratama, Sp.A</strong>. Pantau nomor antrian secara langsung dari ponsel Anda tanpa harus menunggu lama di ruang tunggu.
            </p>

            <!-- Dual Action Buttons -->
            <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-1">
              <button
                @click="emit('navigate', 'kiosk-kidinest')"
                class="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-white font-bold text-xs sm:text-sm uppercase tracking-wider shadow-sm transition-all active:scale-95"
              >
                <span>Ambil Nomor Antrian Online</span>
                <ArrowRight class="w-4 h-4 text-[#E5A93C]" />
              </button>

              <a
                href="#live-antrian"
                class="inline-flex items-center justify-center gap-2 px-5 py-3 rounded-xl bg-white/80 hover:bg-white border border-[#E5A93C]/30 text-[#0C3B36] font-bold text-xs sm:text-sm transition-all shadow-2xs"
              >
                <span class="w-2 h-2 rounded-full bg-[#E5A93C] animate-pulse"></span>
                <span>Pantau Antrian Hari Ini</span>
              </a>
            </div>
          </div>

          <!-- Right: Pure Typographic Pillar (Col 5 - Hairline Divider, Zero Card, Zero Foto) -->
          <div class="lg:col-span-5 lg:pl-8 lg:border-l border-[#0C3B36]/15 space-y-4">
            
            <div class="space-y-1">
              <span class="text-[11px] font-bold uppercase tracking-wider text-[#E5A93C]">
                Standar Pelayanan Pediatri IDAI
              </span>
              <h3 class="text-base sm:text-lg font-bold text-[#0C3B36] leading-snug">
                Kesehatan dan kenyamanan si kecil adalah prioritas utama setiap konsultasi.
              </h3>
            </div>

            <!-- 4 Hairline Focus Rows -->
            <div class="border-t border-b border-[#0C3B36]/10 divide-y divide-[#0C3B36]/10 text-xs text-[#0C3B36]">
              <div class="py-2.5 flex items-baseline justify-between gap-4">
                <span class="font-mono text-[11px] font-bold text-[#E5A93C]">01</span>
                <span class="font-bold flex-1">Vaksinasi Terjaga Cold-Chain</span>
                <span class="text-[11px] opacity-70">Standar IDAI</span>
              </div>
              <div class="py-2.5 flex items-baseline justify-between gap-4">
                <span class="font-mono text-[11px] font-bold text-[#E5A93C]">02</span>
                <span class="font-bold flex-1">Skrining Milestone & Nutrisi</span>
                <span class="text-[11px] opacity-70">Tumbuh Kembang</span>
              </div>
              <div class="py-2.5 flex items-baseline justify-between gap-4">
                <span class="font-mono text-[11px] font-bold text-[#E5A93C]">03</span>
                <span class="font-bold flex-1">Pemeriksaan Ramah & Lembut</span>
                <span class="text-[11px] opacity-70">Bebas Trauma</span>
              </div>
              <div class="py-2.5 flex items-baseline justify-between gap-4">
                <span class="font-mono text-[11px] font-bold text-[#E5A93C]">04</span>
                <span class="font-bold flex-1">Antrian Real-Time di Ponsel</span>
                <span class="text-[11px] opacity-70">Bebas Antre Lama</span>
              </div>
            </div>

            <div class="flex items-center justify-between text-xs text-[#0C3B36]/75 pt-0.5">
              <span class="text-[11px]">Konsultasi bersama:</span>
              <span class="font-bold text-[#0C3B36]">dr. Patricia, Sp.A & dr. Hendra, Sp.A</span>
            </div>

          </div>

        </div>

        <!-- Clinic Facts Strip (Hairline 3-Column - Bukan Card Kotak!) -->
        <div class="mt-10 sm:mt-12 pt-6 border-t border-[#0C3B36]/15 grid grid-cols-1 sm:grid-cols-3 gap-6 text-xs text-[#0C3B36]">
          <div class="space-y-1">
            <div class="font-bold text-[11px] uppercase tracking-wider text-[#0C3B36]/60 flex items-center gap-1.5">
              <span>01 • Dokter Jaga Hari Ini</span>
            </div>
            <div class="font-bold text-sm text-[#0C3B36]">
              dr. Patricia Vanessa, Sp.A & dr. Hendra, Sp.A
            </div>
            <div class="text-[#0C3B36]/70 text-[11px]">
              Spesialis Anak & Tumbuh Kembang IDAI
            </div>
          </div>

          <div class="space-y-1 sm:border-l sm:border-[#0C3B36]/10 sm:pl-6">
            <div class="font-bold text-[11px] uppercase tracking-wider text-[#0C3B36]/60 flex items-center gap-1.5">
              <span>02 • Sesi Layanan Praktik</span>
            </div>
            <div class="font-bold text-sm text-[#0C3B36]">
              Pagi 08:30 – 13:00 • Sore 14:00 – 18:00
            </div>
            <div class="text-[#0C3B36]/70 text-[11px]">
              Senin – Sabtu • Reservasi Online Terbuka
            </div>
          </div>

          <div class="space-y-1 sm:border-l sm:border-[#0C3B36]/10 sm:pl-6">
            <div class="font-bold text-[11px] uppercase tracking-wider text-[#0C3B36]/60 flex items-center gap-1.5">
              <span>03 • Fasilitas Poli Anak</span>
            </div>
            <div class="font-bold text-sm text-[#0C3B36]">
              Klinik Utama Semarang Medical Group
            </div>
            <div class="text-[#0C3B36]/70 text-[11px]">
              Jl. Melati Utara No. 3, Kota Semarang
            </div>
          </div>
        </div>

      </div>
    </section>

    <!-- Section 2: Live Queue Strip (Menyatu Alami dengan Halaman, Tanpa Kotak-Kotak AI) -->
    <section id="live-antrian" class="w-full py-10 sm:py-14 border-b border-[#E5A93C]/20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">
        <LiveQueueFilter
          brand="kidinest"
          :doctors="doctors"
          @register="emit('navigate', 'kiosk-kidinest')"
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
            @click="emit('navigate', 'kiosk-kidinest')"
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
                @click="emit('navigate', 'kiosk-kidinest')"
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

        <!-- 6 Benefit Points Berbaris Rapi (Minimalist Hairline Grid, Tanpa Card Tebal) -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8">
          <div
            v-for="(b, idx) in massageBenefits"
            :key="idx"
            class="py-3 border-b border-[#E5A93C]/20 flex items-start gap-3.5"
          >
            <div class="w-7 h-7 rounded-full bg-[#E5A93C]/25 text-[#0C3B36] font-bold text-xs flex items-center justify-center shrink-0">
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

    <!-- Section 5: Jadwal Dokter Spesialis Anak (Editorial Directory - Tanpa Kotak Card Tebal!) -->
    <section id="jadwal-dokter" class="w-full py-12 sm:py-16 border-b border-[#E5A93C]/20 bg-[#FAF7F2]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="mb-8 flex flex-col sm:flex-row sm:items-end justify-between gap-4">
          <div>
            <span class="text-xs font-bold uppercase tracking-wider text-[#0C3B36]/70">
              Direktori Dokter
            </span>
            <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-1">
              Jadwal Praktik Spesialis Anak Kidinest
            </h2>
          </div>
          <span class="text-xs text-[#0C3B36]/70 font-medium">
            Klinik Utama Semarang Medical Group
          </span>
        </div>

        <!-- Hairline Directory Rows (Bebas Kotak Card!) -->
        <div class="border-t border-b border-[#0C3B36]/15 divide-y divide-[#0C3B36]/10">
          
          <!-- Dokter 1: dr. Patricia Vanessa A., Sp.A -->
          <div class="py-6 flex flex-col md:flex-row md:items-center justify-between gap-5 group hover:bg-white/40 transition-colors px-2 sm:px-4 rounded-xl">
            <div class="flex items-start sm:items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-[#0C3B36] text-[#E5A93C] font-black text-lg flex items-center justify-center shrink-0">
                PV
              </div>
              <div>
                <div class="flex items-center gap-2.5">
                  <h3 class="text-lg font-bold text-[#0C3B36]">
                    dr. Patricia Vanessa A., Sp.A
                  </h3>
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-[#E5A93C]/20 text-[#0C3B36]">
                    Poli 1 • Pagi
                  </span>
                </div>
                <p class="text-xs text-[#0C3B36]/75 mt-0.5 font-medium">
                  Spesialis Ilmu Kesehatan Anak & Vaksinasi Pediatri • Fokus: Imunisasi Lengkap, Nutrisi & Infeksi Anak
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between sm:justify-end gap-6 shrink-0 pl-16 md:pl-0">
              <div class="text-right text-xs">
                <div class="font-bold text-[#0C3B36]">Senin – Sabtu</div>
                <div class="text-[11px] text-[#0C3B36]/70">08:30 – 13:00 WIB</div>
              </div>
              <button
                @click="emit('navigate', 'kiosk-kidinest')"
                class="px-4 py-2 rounded-lg bg-[#0C3B36] text-white text-xs font-bold uppercase tracking-wider hover:bg-[#134E48] transition-colors"
              >
                Daftar →
              </button>
            </div>
          </div>

          <!-- Dokter 2: dr. Hendra Pratama, Sp.A -->
          <div class="py-6 flex flex-col md:flex-row md:items-center justify-between gap-5 group hover:bg-white/40 transition-colors px-2 sm:px-4 rounded-xl">
            <div class="flex items-start sm:items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-[#E5A93C] text-[#0C3B36] font-black text-lg flex items-center justify-center shrink-0">
                HP
              </div>
              <div>
                <div class="flex items-center gap-2.5">
                  <h3 class="text-lg font-bold text-[#0C3B36]">
                    dr. Hendra Pratama, Sp.A
                  </h3>
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-zinc-200 text-zinc-800">
                    Poli 2 • Sore
                  </span>
                </div>
                <p class="text-xs text-[#0C3B36]/75 mt-0.5 font-medium">
                  Spesialis Anak • Konsultan Tumbuh Kembang, Nutrisi & Alergi Anak
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between sm:justify-end gap-6 shrink-0 pl-16 md:pl-0">
              <div class="text-right text-xs">
                <div class="font-bold text-[#0C3B36]">Senin – Sabtu</div>
                <div class="text-[11px] text-[#0C3B36]/70">14:00 – 18:00 WIB</div>
              </div>
              <button
                @click="emit('navigate', 'kiosk-kidinest')"
                class="px-4 py-2 rounded-lg bg-zinc-950 text-white text-xs font-bold uppercase tracking-wider hover:bg-zinc-800 transition-colors"
              >
                Daftar →
              </button>
            </div>
          </div>

        </div>

      </div>
    </section>

    <!-- Section 6: Lokasi & FAQ (Hairline Clean Layout - Bebas Card Tebal) -->
    <section id="lokasi" class="w-full py-12 sm:py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-10 items-start">
          
          <!-- Left: Alamat & Kontak (Col 5) -->
          <div class="lg:col-span-5 space-y-4">
            <div>
              <span class="text-xs font-bold uppercase tracking-wider text-[#0C3B36]/70">
                Akses & Janji Temu
              </span>
              <h2 class="text-2xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight mt-1">
                Kidinest Pediatric Care
              </h2>
            </div>

            <div class="space-y-3 text-xs text-[#0C3B36] pt-1">
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

              <div class="flex items-start gap-2.5 pt-2 border-t border-zinc-100">
                <Clock class="w-4 h-4 text-[#E5A93C] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">Jam Layanan Poli:</div>
                  <div class="opacity-80 mt-0.5">
                    Pagi: 08:30 – 13:00 WIB • Sore: 14:00 – 18:00 WIB
                  </div>
                </div>
              </div>

              <div class="flex items-start gap-2.5 pt-2 border-t border-zinc-100">
                <Phone class="w-4 h-4 text-[#E5A93C] shrink-0 mt-0.5" />
                <div>
                  <div class="font-bold">WhatsApp Reservasi:</div>
                  <div class="opacity-80 mt-0.5 font-mono">
                    +62 812-3456-7890
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Right: FAQ (Hairline Clean Dividers, Bebas Kotak Card) (Col 7) -->
          <div class="lg:col-span-7 space-y-3">
            <h3 class="text-lg font-bold text-[#0C3B36] mb-3">
              Pertanyaan yang Sering Diajukan
            </h3>

            <div class="border-t border-b border-zinc-200/80 divide-y divide-zinc-200/60">
              <div
                v-for="(faq, idx) in faqs"
                :key="idx"
                class="py-3.5"
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
          <button @click="emit('navigate', 'kiosk-kidinest')" class="hover:text-white">Kiosk Antrian Mandiri</button>
          <span>•</span>
          <button @click="emit('navigate', 'ssc')" class="hover:text-white">Semarang Skin Center</button>
          <span>•</span>
          <button @click="emit('navigate', 'login-kidinest')" class="hover:text-white">Login Staf</button>
        </div>
      </div>
    </footer>

  </div>
</template>
