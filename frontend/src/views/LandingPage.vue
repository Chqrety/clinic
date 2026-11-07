<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import Navbar from '../components/Navbar.vue'
import { ArrowRight, Sparkles, ShieldCheck, Clock, CheckCircle2, MapPin, Phone, Stethoscope, HeartPulse } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
}>()

const activeCallingQueue = ref<QueueItem | null>(null)
const waitingCount = ref(0)
const loadingQueue = ref(false)

async function loadQueueSnapshot() {
  loadingQueue.value = true
  try {
    const list = await api.getQueues()
    const calling = list.find(q => q.queue_status === 'CALLING')
    activeCallingQueue.value = calling || list.find(q => q.queue_status === 'IN_EXAM') || null
    waitingCount.value = list.filter(q => q.queue_status === 'WAITING' && q.payment_status === 'PAID').length
  } catch (e) {
    // fallback
  } finally {
    loadingQueue.value = false
  }
}

onMounted(() => {
  loadQueueSnapshot()
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col selection:bg-amber-300">
    <Navbar @navigate="(v) => emit('navigate', v)" />

    <!-- Hero Section (Brandlos Editorial Asymmetric) -->
    <main class="flex-1 max-w-7xl mx-auto px-6 sm:px-10 pt-8 sm:pt-16 pb-20 w-full">

      <!-- Category Pill (Brandlos Warm Aesthetic - No Jarring Blue) -->
      <div class="mb-6">
        <span class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-xs">
          <HeartPulse class="w-3.5 h-3.5 text-amber-400" />
          Klinik Sehat Modern • Layanan Dokter Umum & Apotek
        </span>
      </div>

      <!-- Headline & Asymmetric Body (Brandlos Look) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start mb-16">
        <!-- Giant Heading (7-8 cols) -->
        <div class="lg:col-span-8">
          <h1 class="text-5xl sm:text-7xl lg:text-8xl font-extrabold tracking-tighter text-zinc-950 leading-[0.98]">
            Berobat lebih tenang, pantau antrian dari ponsel.
          </h1>
        </div>

        <!-- Editorial Description (4-5 cols) -->
        <div class="lg:col-span-4 lg:pt-3 space-y-6">
          <p class="text-base sm:text-lg text-zinc-700 leading-relaxed font-normal">
            Kesehatan Anda adalah prioritas kami. Nikmati pelayanan dokter umum yang ramah tanpa perlu berdesakan di ruang tunggu. Cek nomor antrian yang sedang berjalan secara live, ambil tiket karcis mandiri, dan tebus obat di apotek tanpa repot membawa kertas resep.
          </p>

          <div class="pt-2">
            <button
              @click="emit('navigate', 'kiosk')"
              class="group inline-flex items-center gap-2 text-base font-extrabold text-zinc-950 hover:text-zinc-700 transition-colors"
            >
              <span class="w-7 h-7 rounded-full bg-amber-400 group-hover:bg-zinc-950 group-hover:text-white flex items-center justify-center text-xs transition-colors">
                →
              </span>
              <span class="underline underline-offset-4 decoration-zinc-300 group-hover:decoration-zinc-950 transition-colors">
                Ambil nomor antrian pasien sekarang!
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- Live Interactive Ticker Hook (Direct on Canvas, Anti-Card) -->
      <div id="live-queue" class="pt-8 pb-12 border-t border-b border-black/10 my-12">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1 flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
              Sedang Diperiksa Dokter
            </div>
            <div class="text-4xl sm:text-5xl font-extrabold text-zinc-950 font-mono tracking-tight">
              {{ activeCallingQueue?.queue_number || 'A-001' }}
            </div>
            <div class="text-xs text-zinc-500 mt-1 font-medium">
              {{ activeCallingQueue ? activeCallingQueue.patient_name : 'dr. Hendra Pratama, Sp.PD' }} • Ruang 1
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Pasien Sedang Menunggu
            </div>
            <div class="text-4xl sm:text-5xl font-extrabold text-emerald-700 tracking-tight">
              {{ waitingCount }} <span class="text-lg font-normal text-zinc-500">Orang</span>
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              Estimasi tunggu &plusmn; 8 menit per orang
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Layanan Apotek
            </div>
            <div class="text-xl font-extrabold text-zinc-950 mt-2 flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
              Obat Langsung Siap
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              Racikan & sediaan siap usai periksa
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Jam Pelayanan Klinik
            </div>
            <div class="text-xl font-extrabold text-zinc-950 mt-2 flex items-center gap-1.5">
              <Clock class="w-5 h-5 text-zinc-700" />
              08.00 – 21.00 WIB
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              Buka Setiap Hari (Senin – Minggu)
            </div>
          </div>
        </div>
      </div>

      <!-- 3 Key Value Editorial Rows (No Card-Hell) -->
      <section id="keunggulan" class="py-16 space-y-16">
        <div>
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            01 / Pendaftaran Cepat
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              Daftar dari ponsel, langsung dapat nomor antrian.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Tidak perlu lagi berdiri berlama-lama di loket pendaftaran. Cukup masukkan NIK dan nama, pilih bayar QRIS instan atau tunai di kasir. Anda bisa duduk santai di ruang tunggu sambil memantau giliran periksa secara transparan.
            </p>
          </div>
        </div>

        <div class="border-t border-black/10 pt-16">
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            02 / Resep Otomatis ke Apotek
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              Resep langsung diterima apotek, tidak perlu bawa kertas.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Selesai konsultasi dengan dokter, resep otomatis terkirim ke meja apoteker. Anda tidak perlu repot menenteng kertas resep atau takut salah baca tulisan tangan dokter. Obat disiapkan dengan teliti dan langsung dipanggil saat selesai diracik.
            </p>
          </div>
        </div>

        <div class="border-t border-black/10 pt-16">
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            03 / Biaya Transparan & Struk Resmi
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              Bukti pembayaran resmi & transparan, tanpa biaya tersembunyi.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Pasien menerima bukti transaksi resmi yang merinci biaya pemeriksaan dokter dan obat secara jelas. Setiap nominal tercatat rapi sehingga Anda tahu persis apa yang Anda bayarkan tanpa ada pungutan membingungkan.
            </p>
          </div>
        </div>
      </section>

      <!-- Visit Block & Contact Information -->
      <section id="lokasi" class="border-t border-black/10 pt-16 pb-12">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 p-8 rounded-3xl bg-white border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.02)]">
          <div class="space-y-2">
            <div class="flex items-center gap-2 text-zinc-950 font-bold text-sm">
              <MapPin class="w-4 h-4 text-amber-500 flex-shrink-0" />
              <span>Lokasi Klinik</span>
            </div>
            <p class="text-xs text-zinc-600 leading-relaxed">
              Jl. Kesehatan No. 88, Jakarta Pusat<br />
              <span class="text-zinc-400">Akses mudah, dekat halte & area parkir luas</span>
            </p>
          </div>

          <div class="space-y-2">
            <div class="flex items-center gap-2 text-zinc-950 font-bold text-sm">
              <Clock class="w-4 h-4 text-zinc-700 flex-shrink-0" />
              <span>Jam Operasional</span>
            </div>
            <p class="text-xs text-zinc-600 leading-relaxed">
              Senin – Minggu: 08.00 – 21.00 WIB<br />
              <span class="text-zinc-400">Poli Umum & Instalasi Farmasi Buka Setiap Hari</span>
            </p>
          </div>

          <div class="space-y-2">
            <div class="flex items-center gap-2 text-zinc-950 font-bold text-sm">
              <Phone class="w-4 h-4 text-emerald-600 flex-shrink-0" />
              <span>Kontak & WhatsApp Pasien</span>
            </div>
            <p class="text-xs text-zinc-600 leading-relaxed">
              Telp: (021) 555-0199<br />
              WhatsApp: 0812-8877-6655
            </p>
          </div>
        </div>
      </section>

      <!-- Bottom High-Hook CTA Section (Brandlos Yellow Banner or Dark Contrast) -->
      <section class="pt-8 pb-16 flex flex-col sm:flex-row sm:items-end justify-between gap-8">
        <div>
          <h2 class="text-4xl sm:text-6xl font-extrabold tracking-tighter text-zinc-950 leading-tight">
            Pelayanan kesehatan yang ramah & cepat untuk Anda.
          </h2>
          <p class="text-base text-zinc-600 mt-2">
            Daftarkan pemeriksaan Anda atau anggota keluarga sekarang tanpa perlu antri berdesakan.
          </p>
        </div>

        <div class="flex flex-col sm:flex-row gap-3">
          <button
            @click="emit('navigate', 'kiosk')"
            class="px-8 py-4 rounded-full font-extrabold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-95 transition-all text-sm shadow-sm"
          >
            Ambil Antrian Pasien →
          </button>
          <button
            @click="emit('navigate', 'login')"
            class="px-8 py-4 rounded-full font-bold bg-zinc-950 text-white hover:bg-zinc-800 active:scale-95 transition-all text-sm shadow-sm"
          >
            Portal Petugas Medis
          </button>
        </div>
      </section>

    </main>
  </div>
</template>
