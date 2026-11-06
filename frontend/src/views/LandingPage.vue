<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import Navbar from '../components/Navbar.vue'
import { ArrowRight, QrCode, Sparkles, ShieldCheck, Clock, CheckCircle2 } from 'lucide-vue-next'

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

      <!-- Category Pill -->
      <div class="mb-6">
        <span class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider bg-blue-600 text-white shadow-sm">
          <Sparkles class="w-3.5 h-3.5" />
          The Modern Clinic Platform
        </span>
      </div>

      <!-- Headline & Asymmetric Body (Brandlos Look) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12 items-start mb-16">
        <!-- Giant Heading (7 cols) -->
        <div class="lg:col-span-8">
          <h1 class="text-5xl sm:text-7xl lg:text-8xl font-extrabold tracking-tighter text-zinc-950 leading-[0.98]">
            Karena antri berdesakan sudah kuno.
          </h1>
        </div>

        <!-- Editorial Description (5 cols) -->
        <div class="lg:col-span-4 lg:pt-3 space-y-6">
          <p class="text-base sm:text-lg text-zinc-700 leading-relaxed font-normal">
            Pasien berobat untuk lekas sembuh, bukan untuk lelah menunggu balasan 100 chat WhatsApp semalaman. Cukup scan QR dari ponsel: nomor antrian FIFO terbit otomatis, resep dokter langsung terhubung ke apotek tanpa kertas, dan stok obat terpotong presisi.
          </p>

          <div class="pt-2">
            <button
              @click="emit('navigate', 'kiosk')"
              class="group inline-flex items-center gap-2 text-base font-extrabold text-zinc-950 hover:text-blue-600 transition-colors"
            >
              <span class="w-7 h-7 rounded-full bg-amber-400 group-hover:bg-blue-600 group-hover:text-white flex items-center justify-center text-xs transition-colors">
                →
              </span>
              <span class="underline underline-offset-4 decoration-zinc-300 group-hover:decoration-blue-600 transition-colors">
                Ambil antrian mandiri sekarang!
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- Live Interactive Ticker Hook (Direct on Canvas, Anti-Card) -->
      <div id="live-queue" class="pt-8 pb-12 border-t border-b border-black/10 my-12">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Nomor Dipanggil Saat Ini
            </div>
            <div class="text-4xl sm:text-5xl font-extrabold text-zinc-950 font-mono tracking-tight">
              {{ activeCallingQueue?.queue_number || 'A-001' }}
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              {{ activeCallingQueue ? activeCallingQueue.patient_name : 'dr. Hendra Pratama, Sp.PD' }}
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Antrian Menunggu (FIFO)
            </div>
            <div class="text-4xl sm:text-5xl font-extrabold text-emerald-700 tracking-tight">
              {{ waitingCount }} <span class="text-lg font-normal text-zinc-500">Pasien</span>
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              Estimasi tunggu &plusmn; 8 menit/pasien
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Apotek & E-Resep
            </div>
            <div class="text-xl font-extrabold text-zinc-950 mt-2 flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
              Stok Ter-reserve Realtime
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              Nol disposisi kertas resep fisik
            </div>
          </div>

          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-1">
              Standar Kemenkes RI
            </div>
            <div class="text-xl font-extrabold text-blue-700 mt-2 flex items-center gap-1.5">
              <ShieldCheck class="w-5 h-5" />
              SATUSEHAT Ready
            </div>
            <div class="text-xs text-zinc-500 mt-1">
              FHIR R4 Diagnostic & KFA mapping
            </div>
          </div>
        </div>
      </div>

      <!-- 3 Key Value Editorial Rows (No Card-Hell) -->
      <section id="keunggulan" class="py-16 space-y-16">
        <div>
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            01 / Antrian & Pendaftaran
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              Pasien scan QR, bayar instan, nomor antrian langsung terbit.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Tidak ada staf yang harus membaca 100 chat WhatsApp semalaman saat klinik baru buka. Pasien mandiri memilih metode pembayaran: QRIS instan otomatis aktif, atau tunai saat tiba di kasir. Aturan FIFO adil bagi siapa pun.
            </p>
          </div>
        </div>

        <div class="border-t border-black/10 pt-16">
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            02 / Komunikasi Dokter & Apotek
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              E-Resep langsung muncul di layar apoteker dengan stok ter-reserve.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Dokter tidak lagi menulis di kertas yang rawan salah dibaca. Saat resep dikirim, stok obat di rak apotek langsung berstatus "reserved" agar tidak diserobot. Ketika apoteker selesai meracik, stok fisik resmi terpotong dan notifikasi siap diserahkan.
            </p>
          </div>
        </div>

        <div class="border-t border-black/10 pt-16">
          <span class="text-xs font-bold uppercase tracking-widest text-zinc-400 block mb-2">
            03 / Transparansi Billing
          </span>
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-baseline">
            <h2 class="lg:col-span-6 text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950">
              Pasien hanya terima 1 invoice sederhana. Klinik memegang audit internal.
            </h2>
            <p class="lg:col-span-6 text-base text-zinc-600 leading-relaxed">
              Pasien tidak dibingungkan oleh hitungan komisi atau pajak yang rumit: struk mereka ringkas berisi pendaftaran, konsultasi, dan obat. Di sisi internal, manajemen klinik langsung memegang rincian bagi hasil jasa dokter (70%), margin klinik, dan setoran pajak.
            </p>
          </div>
        </div>
      </section>

      <!-- Bottom High-Hook CTA Section (Brandlos Yellow Banner or Dark Contrast) -->
      <section class="mt-8 pt-16 pb-20 border-t border-black/10 flex flex-col sm:flex-row sm:items-end justify-between gap-8">
        <div>
          <h2 class="text-4xl sm:text-6xl font-extrabold tracking-tighter text-zinc-950 leading-tight">
            Mulai pelayanan klinik hari ini.
          </h2>
          <p class="text-base text-zinc-600 mt-2">
            Uji coba simulasi end-to-end dari pendaftaran pasien hingga penyerahan obat di apotek.
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
            Masuk Portal Staf Klinik
          </button>
        </div>
      </section>

    </main>
  </div>
</template>
