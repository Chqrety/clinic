<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { QrCode, Banknote, Clock, RefreshCw, ArrowRight, ArrowLeft, Check, ShieldCheck, Stethoscope } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
}>()

const store = useClinicStore()

const nik = ref('3171012345670008')
const name = ref('Reza Ramadhan')
const phone = ref('081288776655')
const complaint = ref('Demam dan batuk pilek 2 hari')
const paymentMethod = ref<'QRIS' | 'CASH'>('QRIS')

const loading = ref(false)
const errorMsg = ref('')
const activeTicket = ref<QueueItem | null>(store.activePatientQueue)
const aheadCount = ref(0)
const activeCallingNum = ref('A-001')

async function checkAheadQueues(myNum: string) {
  try {
    const list = await api.getQueues()
    const index = list.findIndex(q => q.queue_number === myNum)
    aheadCount.value = index > 0 ? index : 0
    const calling = list.find(q => q.queue_status === 'CALLING')
    if (calling) {
      activeCallingNum.value = calling.queue_number
    }
  } catch (e) {
    // ignore
  }
}

async function handleSubmit() {
  if (!nik.value || !name.value || !complaint.value) {
    errorMsg.value = 'Mohon lengkapi NIK, Nama, dan Keluhan Anda'
    return
  }

  loading.value = true
  errorMsg.value = ''
  try {
    const result = await api.createQueue({
      nik: nik.value,
      name: name.value,
      phone: phone.value,
      complaint: complaint.value,
      payment_method: paymentMethod.value,
    })
    activeTicket.value = result
    store.setActivePatientQueue(result)
    await checkAheadQueues(result.queue_number)
    store.playChime()
  } catch (err: any) {
    errorMsg.value = err.message || 'Gagal mengambil nomor antrian'
  } finally {
    loading.value = false
  }
}

async function refreshStatus() {
  if (!activeTicket.value) return
  loading.value = true
  try {
    const list = await api.getQueues()
    const found = list.find(q => q.id === activeTicket.value!.id)
    if (found) {
      activeTicket.value = found
      store.setActivePatientQueue(found)
      await checkAheadQueues(found.queue_number)
    }
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

function resetForm() {
  activeTicket.value = null
  store.setActivePatientQueue(null)
}

onMounted(() => {
  if (activeTicket.value) {
    checkAheadQueues(activeTicket.value.queue_number)
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 selection:bg-amber-300 flex flex-col">

    <!-- Top Navigation -->
    <header class="w-full max-w-6xl mx-auto px-6 sm:px-10 pt-8 pb-4 flex items-center justify-between">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-2 text-xs font-bold text-zinc-500 hover:text-zinc-950 transition-colors"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Kembali ke Beranda</span>
      </button>

      <span class="text-[11px] font-bold uppercase tracking-widest text-zinc-400">
        Loket Mandiri Pasien
      </span>
    </header>

    <!-- Main Container -->
    <main class="flex-1 max-w-6xl mx-auto w-full px-6 sm:px-10 py-6 sm:py-10">

      <!-- ========================================== -->
      <!-- VIEW 1: ACTIVE DIGITAL TICKET PASS        -->
      <!-- ========================================== -->
      <div v-if="activeTicket" class="max-w-lg mx-auto space-y-6 animate-in fade-in duration-200">
        <!-- Physical Boarding Slip / Appointment Receipt -->
        <div class="bg-white rounded-3xl p-8 sm:p-10 border border-black/10 shadow-[0_4px_24px_rgba(0,0,0,0.04)] relative">

          <!-- Ticket Header -->
          <div class="flex items-center justify-between pb-6 border-b border-zinc-100">
            <div>
              <span class="text-[11px] font-extrabold uppercase tracking-widest text-zinc-400 block">Klinik Sehat Modern</span>
              <h2 class="text-base font-extrabold text-zinc-950 mt-0.5">Tiket Antrian Digital</h2>
            </div>
            <button
              @click="refreshStatus"
              :class="['p-2 rounded-full hover:bg-stone-100 text-zinc-500 transition-transform active:scale-90', loading && 'animate-spin']"
              title="Perbarui Status Antrian"
            >
              <RefreshCw class="w-4 h-4" />
            </button>
          </div>

          <!-- Queue Number Highlight -->
          <div class="py-10 text-center">
            <span class="text-xs font-bold uppercase tracking-widest text-zinc-400">Nomor Panggilan Anda</span>
            <div class="text-6xl sm:text-7xl font-black tracking-tight text-zinc-950 font-mono my-2">
              {{ activeTicket.queue_number }}
            </div>

            <!-- Status Indicator -->
            <div
              class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold mt-2"
              :class="activeTicket.payment_status === 'PAID'
                ? 'bg-emerald-50 text-emerald-800 border border-emerald-200/60'
                : 'bg-stone-100 text-zinc-700 border border-zinc-200'"
            >
              <span
                class="w-2 h-2 rounded-full"
                :class="activeTicket.payment_status === 'PAID' ? 'bg-emerald-600 animate-pulse' : 'bg-zinc-400'"
              ></span>
              <span>
                {{ activeTicket.payment_status === 'PAID' ? 'Antrian Aktif Terverifikasi' : 'Menunggu Konfirmasi Kasir' }}
              </span>
            </div>

            <!-- Wait Notice -->
            <p class="text-xs text-zinc-600 mt-3 font-medium">
              <span v-if="activeTicket.payment_status === 'PENDING'">
                Silakan lakukan pembayaran di kasir untuk mengaktifkan giliran panggilan dokter.
              </span>
              <span v-else-if="aheadCount > 0">
                Ada <strong>{{ aheadCount }} pasien</strong> di depan Anda. Harap menunggu di ruang tunggu.
              </span>
              <span v-else class="text-zinc-950 font-extrabold">
                Giliran Anda! Silakan langsung menuju ke Ruang Periksa Dokter.
              </span>
            </p>
          </div>

          <!-- Details Rows -->
          <div class="pt-6 border-t border-zinc-100 space-y-3 text-xs">
            <div class="flex justify-between items-center py-1">
              <span class="text-zinc-400">Nama Pasien</span>
              <span class="font-extrabold text-zinc-900">{{ activeTicket.patient_name }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">NIK KTP</span>
              <span class="font-mono font-semibold text-zinc-700">{{ activeTicket.patient_nik }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">Keluhan Gejala</span>
              <span class="font-medium text-zinc-800 text-right max-w-[220px] truncate">{{ activeTicket.complaint }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">Pembayaran</span>
              <span class="font-bold text-zinc-900">
                {{ activeTicket.payment_method }}
                <span v-if="activeTicket.payment_status === 'PAID'" class="text-emerald-700 font-semibold ml-1">(Lunas)</span>
                <span v-else class="text-zinc-500 font-normal ml-1">(Pending Loket)</span>
              </span>
            </div>
            <div v-if="activeTicket.patient_ihs" class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400 flex items-center gap-1">
                <ShieldCheck class="w-3.5 h-3.5 text-blue-600" /> SATUSEHAT IHS
              </span>
              <span class="font-mono font-bold text-zinc-900">{{ activeTicket.patient_ihs }}</span>
            </div>
          </div>
        </div>

        <button
          @click="resetForm"
          class="w-full py-3.5 rounded-full text-xs font-bold text-zinc-500 hover:text-zinc-950 bg-white border border-black/10 hover:border-black/20 transition-all active:scale-95 shadow-sm"
        >
          Daftar Pasien Lain / Ambil Antrian Baru
        </button>
      </div>

      <!-- ========================================== -->
      <!-- VIEW 2: FULL CANVAS FORM (NO CARD WRAPPER) -->
      <!-- ========================================== -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-16 items-start animate-in fade-in duration-200">

        <!-- Left Column: Context, Narrative & Live Pulse (5 cols) -->
        <div class="lg:col-span-5 space-y-6 lg:sticky lg:top-24">
          <div>
            <span class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-blue-600 text-white mb-3">
              Pendaftaran Mandiri
            </span>
            <h1 class="text-4xl sm:text-5xl font-extrabold tracking-tight text-zinc-950 leading-tight">
              Ambil nomor antrian tanpa antri di loket.
            </h1>
            <p class="text-sm text-zinc-600 leading-relaxed mt-3">
              Cukup masukkan NIK dan keluhan Anda. Nomor antrian FIFO resmi terbit seketika dan langsung tersinkron ke ruang periksa dokter.
            </p>
          </div>

          <!-- Live Clinic Pulse Strip (Direct on Canvas) -->
          <div class="py-6 border-t border-b border-black/10 space-y-4">
            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-500 font-medium">Nomor Dipanggil Saat Ini:</span>
              <span class="font-mono font-black text-base text-zinc-950">{{ activeCallingNum }}</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-500 font-medium">Estimasi Waktu Tunggu:</span>
              <span class="font-bold text-zinc-900">&plusmn; 8 menit per pasien</span>
            </div>
            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-500 font-medium">Dokter Siaga:</span>
              <span class="font-bold text-zinc-900 flex items-center gap-1">
                <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                dr. Hendra Pratama, Sp.PD
              </span>
            </div>
          </div>

          <div class="text-xs text-zinc-400">
            Terhubung otomatis dengan Kamus Farmasi (KFA) dan platform SATUSEHAT Kemenkes RI.
          </div>
        </div>

        <!-- Right Column: Registration Form (Direct Canvas, 7 cols) -->
        <div class="lg:col-span-7">
          <form @submit.prevent="handleSubmit" class="space-y-6">

            <!-- Error Notification -->
            <div v-if="errorMsg" class="p-4 rounded-2xl bg-rose-50 text-rose-800 text-xs font-bold border border-rose-200">
              {{ errorMsg }}
            </div>

            <!-- NIK Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-2">
                1. NIK (Nomor Induk Kependudukan)
              </label>
              <input
                v-model="nik"
                type="text"
                maxlength="16"
                required
                placeholder="16 digit sesuai e-KTP"
                class="w-full px-5 py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- Nama Lengkap Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-2">
                2. Nama Lengkap Pasien
              </label>
              <input
                v-model="name"
                type="text"
                required
                placeholder="Nama lengkap sesuai KTP"
                class="w-full px-5 py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- WhatsApp Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-2">
                3. Nomor WhatsApp / Ponsel
              </label>
              <input
                v-model="phone"
                type="tel"
                required
                placeholder="0812xxxxxxxx"
                class="w-full px-5 py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- Keluhan Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-2">
                4. Keluhan yang Dirasakan
              </label>
              <textarea
                v-model="complaint"
                rows="3"
                required
                placeholder="Tuliskan keluhan atau gejala sakit yang dirasakan..."
                class="w-full px-5 py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-sm font-medium outline-none transition-all resize-none placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              ></textarea>
            </div>

            <!-- Payment Method Field -->
            <div class="pt-2">
              <div class="flex items-baseline justify-between mb-2.5">
                <label class="text-xs font-bold uppercase tracking-wider text-zinc-600">
                  5. Biaya Pendaftaran & Pilihan Bayar
                </label>
                <span class="text-xs font-semibold text-zinc-500">
                  Total: <strong class="text-zinc-950 font-black text-sm">Rp 35.000</strong>
                </span>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <!-- QRIS Option -->
                <button
                  type="button"
                  @click="paymentMethod = 'QRIS'"
                  :class="[
                    'p-4 rounded-2xl text-left flex items-start justify-between transition-all active:scale-[0.99] shadow-sm',
                    paymentMethod === 'QRIS'
                      ? 'bg-white border-2 border-zinc-950'
                      : 'bg-white/60 border border-black/10 hover:border-black/25 text-zinc-600'
                  ]"
                >
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <QrCode class="w-4 h-4 text-zinc-950" />
                      <span class="font-extrabold text-xs text-zinc-950">QRIS / Instant</span>
                    </div>
                    <div class="text-[11px] text-zinc-500">
                      Antrian otomatis aktif seketika
                    </div>
                  </div>
                  <div v-if="paymentMethod === 'QRIS'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center text-[10px]">
                    ✓
                  </div>
                </button>

                <!-- CASH Option -->
                <button
                  type="button"
                  @click="paymentMethod = 'CASH'"
                  :class="[
                    'p-4 rounded-2xl text-left flex items-start justify-between transition-all active:scale-[0.99] shadow-sm',
                    paymentMethod === 'CASH'
                      ? 'bg-white border-2 border-zinc-950'
                      : 'bg-white/60 border border-black/10 hover:border-black/25 text-zinc-600'
                  ]"
                >
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <Banknote class="w-4 h-4 text-zinc-950" />
                      <span class="font-extrabold text-xs text-zinc-950">Tunai di Kasir</span>
                    </div>
                    <div class="text-[11px] text-zinc-500">
                      Bayar saat tiba di loket klinik
                    </div>
                  </div>
                  <div v-if="paymentMethod === 'CASH'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center text-[10px]">
                    ✓
                  </div>
                </button>
              </div>
            </div>

            <!-- Submit Button (Direct on Canvas, Bold & Taktil) -->
            <div class="pt-4">
              <button
                type="submit"
                :disabled="loading"
                class="w-full py-4 rounded-full font-extrabold bg-zinc-950 hover:bg-zinc-800 text-white active:scale-[0.98] transition-all text-sm shadow-md flex items-center justify-center gap-2 disabled:opacity-50"
              >
                <span v-if="loading">Menerbitkan Nomor Antrian...</span>
                <span v-else class="flex items-center gap-2">
                  <span>Konfirmasi & Ambil Nomor Antrian</span>
                  <ArrowRight class="w-4 h-4" />
                </span>
              </button>
            </div>

          </form>
        </div>

      </div>

    </main>
  </div>
</template>
