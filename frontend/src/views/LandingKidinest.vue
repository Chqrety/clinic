<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BrandSwitcherHeader from '../components/BrandSwitcherHeader.vue'
import LiveQueueFilter, { type DoctorInfo } from '../components/LiveQueueFilter.vue'
import {
  Clock,
  MapPin,
  Phone,
  ArrowRight,
  Calendar,
  CalendarDays,
  UserCheck,
  CheckCircle2,
  MessageCircle,
  Mail,
  CreditCard,
  Sparkles
} from 'lucide-vue-next'
import { api } from '../api/client'

const emit = defineEmits<{
  (e: 'navigate', view: string): void
}>()

interface DoctorScheduleSummary {
  doctor_id: number
  doctor_name: string
  role: string
  clinic_unit: string
  room_name: string
  specialty: string
  shift_name: string
  days_text: string
  time_text: string
  start_time: string
  end_time: string
  slot_minutes: number
  quota_per_slot: number
  is_practicing_today: boolean
  is_on_leave_today: boolean
  today_status: string
  weekly_schedules: any[]
  available_dates?: any[]
}

const doctorSchedules = ref<DoctorScheduleSummary[]>([])
const doctors = ref<DoctorInfo[]>([])
const isLoadingSchedules = ref(true)

function getDoctorInitials(name: string): string {
  const clean = name.replace(/^dr\.\s+/i, '').replace(/,\s*.*$/, '')
  const parts = clean.split(' ').filter(Boolean)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return (parts[0] ? parts[0].substring(0, 2) : 'DR').toUpperCase()
}

async function fetchDoctorSchedules() {
  isLoadingSchedules.value = true
  try {
    const data = await api.getPublicDoctorSchedulesSummary('KIDINEST')
    if (data && Array.isArray(data) && data.length > 0) {
      doctorSchedules.value = data
      // Sync doctors list for live queue filter component
      doctors.value = data.map((d: any) => ({
        id: String(d.doctor_id),
        name: d.doctor_name,
        specialty: d.specialty,
        room: '',
        status: (d.is_practicing_today ? 'active' : 'off') as 'active' | 'off' | 'break'
      }))
    }
  } catch (err) {
    console.warn('Gagal memuat jadwal dokter live:', err)
  } finally {
    isLoadingSchedules.value = false
  }
}

onMounted(() => {
  fetchDoctorSchedules()
})

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

        <!-- Headline & Typographic Focus (Murni Tipografi & Layout Bersih) -->
        <div class="max-w-3xl space-y-5">
          <h1 class="text-3xl sm:text-5xl font-extrabold text-[#0C3B36] tracking-tight leading-[1.15]">
            Pemeriksaan Buah Hati yang Tenang, Cermat, dan Hangat.
          </h1>
          <p class="text-sm sm:text-base text-[#0C3B36]/80 leading-relaxed font-medium max-w-2xl">
            Layanan dokter spesialis anak berstandar IDAI bersama <strong>dr. Patricia Vanessa A., Sp.A</strong> dan <strong>dr. Hendra Pratama, Sp.A</strong>. Pantau nomor antrian secara langsung dari ponsel Anda tanpa harus menunggu lama di ruang tunggu.
          </p>

          <!-- Dual Action Buttons -->
          <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-1">
            <button
              @click="emit('navigate', 'kiosk-kidinest')"
              class="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-white font-bold text-xs sm:text-sm uppercase tracking-wider shadow-sm transition-all active:scale-95 cursor-pointer"
            >
              <span>Pendaftaran Online</span>
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

      </div>
    </section>

    <!-- Section: Jadwal Dokter (Dibawah Hero Card - Real-Time dari Inputan Sistem Dokter) -->
    <section id="jadwal-dokter-hero" class="w-full py-8 sm:py-12 border-b border-[#E5A93C]/20 bg-gradient-to-b from-[#F5EFE6] to-[#FAF7F2]">
      <div class="max-w-7xl mx-auto px-4 sm:px-8">

        <!-- Top Header -->
        <div class="mb-6">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-[11px] font-bold bg-[#0C3B36] text-[#E5A93C] mb-2 shadow-2xs">
            <span class="w-2 h-2 rounded-full bg-[#E5A93C] animate-pulse"></span>
            <span>Jadwal Dokter Terkini • Inputan Resmi Dokter</span>
          </div>
          <h2 class="text-xl sm:text-3xl font-extrabold text-[#0C3B36] tracking-tight">
            Jadwal Praktik Spesialis Anak Kidinest
          </h2>
          <p class="text-xs sm:text-sm text-[#0C3B36]/75 mt-1 font-medium">
            Data sinkron otomatis langsung dari pengaturan shift kerja dan status cuti dokter di sistem.
          </p>
        </div>

        <!-- Schedule Cards Grid (Live Dynamic Cards) -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
          <div
            v-for="doc in doctorSchedules"
            :key="'card-' + doc.doctor_id"
            class="bg-white rounded-2xl border border-[#E5A93C]/30 p-5 sm:p-6 shadow-xs hover:border-[#E5A93C] hover:shadow-sm transition-all flex flex-col justify-between group"
          >
            <!-- Card Top: Avatar, Name, Specialty, Room -->
            <div class="space-y-4">
              <div class="flex items-start justify-between gap-3">
                <div class="flex items-center gap-3.5">
                  <div
                    class="w-12 h-12 rounded-xl flex items-center justify-center font-black text-base shrink-0 shadow-2xs"
                    :class="doc.doctor_id === 2 ? 'bg-[#0C3B36] text-[#E5A93C]' : 'bg-[#E5A93C] text-[#0C3B36]'"
                  >
                    {{ getDoctorInitials(doc.doctor_name) }}
                  </div>
                  <div>
                    <h3 class="text-base sm:text-lg font-bold text-[#0C3B36] group-hover:text-[#134E48] transition-colors">
                      {{ doc.doctor_name }}
                    </h3>
                    <p class="text-xs font-semibold text-[#0C3B36]/75 mt-0.5">
                      {{ doc.specialty }}
                    </p>
                  </div>
                </div>

                <!-- Shift Badge -->
                <span
                  class="px-2.5 py-1 rounded-md text-[10px] font-bold uppercase tracking-wider shrink-0"
                  :class="doc.shift_name === 'PAGI' ? 'bg-[#E5A93C]/20 text-[#0C3B36] border border-[#E5A93C]/40' : 'bg-[#0C3B36]/10 text-[#0C3B36] border border-[#0C3B36]/20'"
                >
                  {{ doc.shift_name || 'REGULER' }}
                </span>
              </div>

              <!-- Information Grid -->
              <div class="grid grid-cols-2 gap-3 pt-3 border-t border-[#0C3B36]/10">
                <!-- Days Range -->
                <div class="space-y-0.5">
                  <span class="text-[10px] font-bold uppercase tracking-wider text-[#0C3B36]/60">Hari Praktik</span>
                  <div class="text-xs font-bold text-[#0C3B36] flex items-center gap-1.5">
                    <Calendar class="w-3.5 h-3.5 text-[#E5A93C] shrink-0" />
                    <span>{{ doc.days_text }}</span>
                  </div>
                </div>

                <!-- Hours -->
                <div class="col-span-2 space-y-0.5 pt-1">
                  <span class="text-[10px] font-bold uppercase tracking-wider text-[#0C3B36]/60">Jam Layanan Konsultasi</span>
                  <div class="text-xs sm:text-sm font-extrabold text-[#0C3B36] flex items-center gap-1.5 font-mono">
                    <Clock class="w-3.5 h-3.5 text-[#E5A93C] shrink-0" />
                    <span>{{ doc.time_text }}</span>
                  </div>
                </div>
              </div>
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

        <!-- Hairline Directory Rows (Dynamic Data dari Inputan Dokter) -->
        <div class="border-t border-b border-[#0C3B36]/15 divide-y divide-[#0C3B36]/10">
          
          <div
            v-for="doc in doctorSchedules"
            :key="'dir-' + doc.doctor_id"
            class="py-6 flex flex-col md:flex-row md:items-center justify-between gap-5 group hover:bg-white/40 transition-colors px-2 sm:px-4 rounded-xl"
          >
            <div class="flex items-start sm:items-center gap-4">
              <div
                class="w-12 h-12 rounded-xl flex items-center justify-center font-black text-lg shrink-0 shadow-2xs"
                :class="doc.doctor_id === 2 ? 'bg-[#0C3B36] text-[#E5A93C]' : 'bg-[#E5A93C] text-[#0C3B36]'"
              >
                {{ getDoctorInitials(doc.doctor_name) }}
              </div>
              <div>
                <div class="flex items-center gap-2.5">
                  <h3 class="text-lg font-bold text-[#0C3B36]">
                    {{ doc.doctor_name }}
                  </h3>
                  <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-[#E5A93C]/20 text-[#0C3B36]">
                    {{ doc.shift_name ? 'Praktik ' + doc.shift_name : 'Praktik' }}
                  </span>
                </div>
                <p class="text-xs text-[#0C3B36]/75 mt-0.5 font-medium">
                  {{ doc.specialty }} • <span class="font-bold text-[#0C3B36]">{{ doc.today_status }}</span>
                </p>
              </div>
            </div>

            <div class="flex items-center justify-between sm:justify-end gap-6 shrink-0 pl-16 md:pl-0">
              <div class="text-right text-xs">
                <div class="font-bold text-[#0C3B36]">{{ doc.days_text }}</div>
                <div class="text-[11px] text-[#0C3B36]/70">{{ doc.time_text }}</div>
              </div>
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
                Semarang Medical Group
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

              <div class="pt-3 border-t border-zinc-100 space-y-2">
                <div class="font-bold text-[#0C3B36] flex items-center gap-2">
                  <Phone class="w-4 h-4 text-[#E5A93C] shrink-0" />
                  <span>Kontak Resmi Semarang Medical Group:</span>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
                  <div class="p-2.5 rounded-xl bg-[#FAF7F2] border border-[#0C3B36]/10 flex items-center gap-2">
                    <MessageCircle class="w-3.5 h-3.5 text-emerald-600 shrink-0" />
                    <div>
                      <div class="text-[10px] uppercase font-bold text-[#0C3B36]/60">WhatsApp Resmi</div>
                      <div class="font-mono font-bold text-[#0C3B36]">+62 812-3456-7890</div>
                    </div>
                  </div>
                  <div class="p-2.5 rounded-xl bg-[#FAF7F2] border border-[#0C3B36]/10 flex items-center gap-2">
                    <Phone class="w-3.5 h-3.5 text-[#0C3B36]/70 shrink-0" />
                    <div>
                      <div class="text-[10px] uppercase font-bold text-[#0C3B36]/60">Nomor Kantor</div>
                      <div class="font-mono text-[#0C3B36]/75 italic text-xs">(Menyusul)</div>
                    </div>
                  </div>
                  <div class="p-2.5 rounded-xl bg-[#FAF7F2] border border-[#0C3B36]/10 flex items-center gap-2">
                    <Mail class="w-3.5 h-3.5 text-[#0C3B36]/70 shrink-0" />
                    <div>
                      <div class="text-[10px] uppercase font-bold text-[#0C3B36]/60">Email Resmi</div>
                      <div class="font-mono text-[#0C3B36]/75 italic text-xs">(Menyusul)</div>
                    </div>
                  </div>
                  <div class="p-2.5 rounded-xl bg-[#FAF7F2] border border-[#0C3B36]/10 flex items-center gap-2">
                    <CreditCard class="w-3.5 h-3.5 text-[#0C3B36]/70 shrink-0" />
                    <div>
                      <div class="text-[10px] uppercase font-bold text-[#0C3B36]/60">Nomor Rekening Klinik</div>
                      <div class="font-mono text-[#0C3B36]/75 italic text-xs">(Menyusul)</div>
                    </div>
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

    <!-- Footer Solid Hijau Botol Khas Kidinest & Semarang Medical Group -->
    <footer class="w-full bg-[#0C3B36] text-white py-8 text-xs select-none border-t border-[#E5A93C]/20">
      <div class="max-w-7xl mx-auto px-4 sm:px-8 space-y-5">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <!-- Brand & Address -->
          <div class="space-y-1">
            <div class="flex items-center gap-2.5">
              <img src="/assets/kidinest/logo-maskot.png" class="w-5 h-5 object-contain" />
              <span class="font-extrabold text-sm tracking-tight text-white">Semarang Medical Group</span>
              <span class="text-white/40">•</span>
              <span class="text-white/80 font-bold">Kidinest Pediatric Care</span>
            </div>
            <p class="text-white/70 text-xs">
              Klinik Spesialis Anak & Tumbuh Kembang • Jl. Melati Utara No. 3, Kota Semarang
            </p>
          </div>

          <!-- Nomor Telepon Kantor, WhatsApp, Email, & Rekening Berdampingan Rapi -->
          <div class="flex flex-wrap items-center gap-3 bg-white/10 px-4 py-2.5 rounded-2xl border border-white/15 text-xs">
            <div class="flex items-center gap-1.5">
              <MessageCircle class="w-3.5 h-3.5 text-emerald-400" />
              <span class="text-white/70 text-[11px]">WhatsApp:</span>
              <strong class="font-mono text-white text-xs">+62 812-3456-7890</strong>
            </div>
            <span class="text-white/30 hidden sm:inline">•</span>
            <div class="flex items-center gap-1.5">
              <Phone class="w-3.5 h-3.5 text-[#E5A93C]" />
              <span class="text-white/70 text-[11px]">Nomor Kantor:</span>
              <span class="italic text-white/80 text-xs">(Menyusul)</span>
            </div>
            <span class="text-white/30 hidden sm:inline">•</span>
            <div class="flex items-center gap-1.5">
              <Mail class="w-3.5 h-3.5 text-white/70" />
              <span class="text-white/70 text-[11px]">Email:</span>
              <span class="italic text-white/80 text-xs">(Menyusul)</span>
            </div>
            <span class="text-white/30 hidden sm:inline">•</span>
            <div class="flex items-center gap-1.5">
              <CreditCard class="w-3.5 h-3.5 text-amber-300" />
              <span class="text-white/70 text-[11px]">No. Rekening:</span>
              <span class="italic text-white/80 text-xs">(Menyusul)</span>
            </div>
          </div>
        </div>

        <div class="pt-4 border-t border-white/10 flex flex-col sm:flex-row items-center justify-between gap-3 text-white/60 text-[11px]">
          <div>© 2026 CV Semarang Medical Group. Seluruh hak cipta dilindungi.</div>
          <div class="flex items-center gap-4 font-semibold text-white/80">
            <button @click="emit('navigate', 'kiosk-kidinest')" class="hover:text-white transition-colors cursor-pointer">Pendaftaran Online</button>
            <span>•</span>
            <button @click="emit('navigate', 'ssc')" class="hover:text-white transition-colors cursor-pointer">Semarang Skin Center</button>
            <span>•</span>
            <button @click="emit('navigate', 'login-kidinest')" class="hover:text-white transition-colors cursor-pointer">Portal Staf</button>
          </div>
        </div>
      </div>
    </footer>

  </div>
</template>
