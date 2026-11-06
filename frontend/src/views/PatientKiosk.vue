<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { QrCode, Banknote, CheckCircle2, Clock, RefreshCw, ArrowRight, ArrowLeft, Check, ShieldCheck } from 'lucide-vue-next'

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

async function checkAheadQueues(myNum: string) {
  try {
    const list = await api.getQueues()
    const index = list.findIndex(q => q.queue_number === myNum)
    aheadCount.value = index > 0 ? index : 0
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
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-900 py-10 px-4 flex flex-col justify-center items-center">

    <!-- Top Navigation Bar -->
    <div class="w-full max-w-lg mb-6 flex items-center justify-between">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-1.5 text-xs font-semibold text-zinc-500 hover:text-zinc-900 transition-colors"
      >
        <ArrowLeft class="w-3.5 h-3.5" /> Beranda Klinik
      </button>

      <span class="text-[11px] font-bold text-zinc-400 uppercase tracking-widest">
        Pelayanan Pasien Mandiri
      </span>
    </div>

    <div class="w-full max-w-lg">

      <!-- View 1: Active Ticket (Clean Minimalist Receipt, No AI Slop) -->
      <div v-if="activeTicket" class="space-y-4 animate-in fade-in duration-200">
        <div class="bg-white rounded-3xl p-8 border border-zinc-200/80 shadow-[0_2px_16px_rgba(0,0,0,0.03)]">

          <!-- Ticket Header -->
          <div class="flex items-center justify-between pb-6 border-b border-zinc-100">
            <div>
              <span class="text-xs font-bold uppercase tracking-wider text-zinc-400 block">Klinik Sehat Modern</span>
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

          <!-- Big Queue Number (Clean High-Contrast Typography) -->
          <div class="py-10 text-center">
            <span class="text-[11px] font-bold uppercase tracking-widest text-zinc-400">Nomor Giliran Anda</span>
            <div class="text-6xl sm:text-7xl font-black tracking-tight text-zinc-950 font-mono my-2">
              {{ activeTicket.queue_number }}
            </div>

            <!-- Soft Status Pill -->
            <div class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-semibold mt-2"
              :class="activeTicket.payment_status === 'PAID'
                ? 'bg-emerald-50 text-emerald-800 border border-emerald-200/70'
                : 'bg-stone-100 text-zinc-700 border border-zinc-200'"
            >
              <span class="w-1.5 h-1.5 rounded-full"
                :class="activeTicket.payment_status === 'PAID' ? 'bg-emerald-600 animate-pulse' : 'bg-zinc-400'"
              ></span>
              <span>
                {{ activeTicket.payment_status === 'PAID' ? 'Antrian Aktif Terjadwal' : 'Menunggu Konfirmasi Kasir' }}
              </span>
            </div>

            <!-- Wait Context -->
            <p class="text-xs text-zinc-500 mt-3 font-medium">
              <span v-if="activeTicket.payment_status === 'PENDING'">
                Silakan lakukan pembayaran tunai di loket pendaftaran agar dipanggil dokter.
              </span>
              <span v-else-if="aheadCount > 0">
                Terdapat <strong>{{ aheadCount }}</strong> pasien di depan Anda (Estimasi &plusmn; {{ aheadCount * 8 }} menit).
              </span>
              <span v-else class="text-zinc-900 font-bold">
                Giliran Anda telah tiba. Silakan langsung menuju ke Ruang Periksa Dokter.
              </span>
            </p>
          </div>

          <!-- Structured Details (Hairline Table, No Cluttered Cards) -->
          <div class="pt-6 border-t border-zinc-100 space-y-3 text-xs">
            <div class="flex justify-between items-center py-1">
              <span class="text-zinc-400">Nama Pasien</span>
              <span class="font-bold text-zinc-900">{{ activeTicket.patient_name }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">NIK Pasien</span>
              <span class="font-mono font-medium text-zinc-700">{{ activeTicket.patient_nik }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">Keluhan Gejala</span>
              <span class="font-medium text-zinc-800 text-right max-w-[220px] truncate">{{ activeTicket.complaint }}</span>
            </div>
            <div class="flex justify-between items-center py-1 border-t border-zinc-50">
              <span class="text-zinc-400">Metode Bayar</span>
              <span class="font-bold text-zinc-900">
                {{ activeTicket.payment_method }}
                <span v-if="activeTicket.payment_status === 'PAID'" class="text-emerald-700 font-semibold ml-1">(Lunas)</span>
                <span v-else class="text-zinc-500 font-normal ml-1">(Menunggu Kasir)</span>
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
          class="w-full py-3 rounded-full text-xs font-bold text-zinc-600 hover:text-zinc-950 bg-white border border-zinc-200/80 hover:border-zinc-300 transition-all active:scale-[0.99]"
        >
          Ambil Nomor Antrian Lain
        </button>
      </div>

      <!-- View 2: Elegant Registration Form (Calm, Airy, Harmonious) -->
      <div v-else class="bg-white rounded-3xl p-8 sm:p-10 border border-zinc-200/80 shadow-[0_2px_16px_rgba(0,0,0,0.03)] animate-in fade-in duration-200">

        <!-- Header -->
        <div class="mb-8 pb-6 border-b border-zinc-100">
          <span class="text-xs font-bold uppercase tracking-wider text-blue-700 block mb-1">
            Pendaftaran Mandiri
          </span>
          <h1 class="text-2xl font-extrabold text-zinc-950 tracking-tight">
            Ambil Nomor Antrian
          </h1>
          <p class="text-xs text-zinc-500 mt-1">
            Sistem antrian terhubung langsung dengan ruang dokter dan loket kasir.
          </p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-5">
          <!-- Error Alert -->
          <div v-if="errorMsg" class="p-3 rounded-2xl bg-rose-50 text-rose-800 text-xs font-medium border border-rose-200">
            {{ errorMsg }}
          </div>

          <!-- NIK -->
          <div>
            <label class="block text-xs font-bold text-zinc-700 mb-1.5">
              NIK (Nomor Induk Kependudukan)
            </label>
            <input
              v-model="nik"
              type="text"
              maxlength="16"
              required
              placeholder="16 digit NIK sesuai KTP"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200/80 focus:border-zinc-900 focus:bg-white text-zinc-900 text-xs font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal"
            />
          </div>

          <!-- Nama Lengkap -->
          <div>
            <label class="block text-xs font-bold text-zinc-700 mb-1.5">
              Nama Pasien
            </label>
            <input
              v-model="name"
              type="text"
              required
              placeholder="Nama lengkap pasien"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200/80 focus:border-zinc-900 focus:bg-white text-zinc-900 text-xs font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal"
            />
          </div>

          <!-- No Telepon / WA -->
          <div>
            <label class="block text-xs font-bold text-zinc-700 mb-1.5">
              Nomor WhatsApp / Ponsel
            </label>
            <input
              v-model="phone"
              type="tel"
              required
              placeholder="Contoh: 08123456789"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200/80 focus:border-zinc-900 focus:bg-white text-zinc-900 text-xs font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal"
            />
          </div>

          <!-- Keluhan Singkat -->
          <div>
            <label class="block text-xs font-bold text-zinc-700 mb-1.5">
              Keluhan yang Dirasakan
            </label>
            <textarea
              v-model="complaint"
              rows="2"
              required
              placeholder="Ceritakan singkat keluhan utama Anda..."
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200/80 focus:border-zinc-900 focus:bg-white text-zinc-900 text-xs font-normal outline-none transition-all resize-none placeholder:text-zinc-400"
            ></textarea>
          </div>

          <!-- Pilihan Pembayaran (Harmonious Selection, No Harsh Contrast) -->
          <div class="pt-2">
            <div class="flex items-baseline justify-between mb-2">
              <label class="text-xs font-bold text-zinc-700">Metode Pembayaran</label>
              <span class="text-xs text-zinc-500 font-medium">Biaya: <strong class="text-zinc-950 font-extrabold">Rp 35.000</strong></span>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <!-- QRIS Option -->
              <button
                type="button"
                @click="paymentMethod = 'QRIS'"
                :class="[
                  'p-3.5 rounded-2xl border text-left flex flex-col justify-between transition-all active:scale-[0.98]',
                  paymentMethod === 'QRIS'
                    ? 'border-zinc-900 bg-stone-50/90 shadow-sm'
                    : 'border-zinc-200/80 bg-white hover:border-zinc-300'
                ]"
              >
                <div class="flex items-center justify-between mb-3">
                  <QrCode class="w-4 h-4 text-zinc-800" />
                  <span v-if="paymentMethod === 'QRIS'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center">
                    <Check class="w-2.5 h-2.5" />
                  </span>
                </div>
                <div>
                  <div class="text-xs font-bold text-zinc-950">QRIS / Digital</div>
                  <div class="text-[11px] text-zinc-500 mt-0.5">Antrian langsung aktif</div>
                </div>
              </button>

              <!-- CASH Option -->
              <button
                type="button"
                @click="paymentMethod = 'CASH'"
                :class="[
                  'p-3.5 rounded-2xl border text-left flex flex-col justify-between transition-all active:scale-[0.98]',
                  paymentMethod === 'CASH'
                    ? 'border-zinc-900 bg-stone-50/90 shadow-sm'
                    : 'border-zinc-200/80 bg-white hover:border-zinc-300'
                ]"
              >
                <div class="flex items-center justify-between mb-3">
                  <Banknote class="w-4 h-4 text-zinc-800" />
                  <span v-if="paymentMethod === 'CASH'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center">
                    <Check class="w-2.5 h-2.5" />
                  </span>
                </div>
                <div>
                  <div class="text-xs font-bold text-zinc-950">Tunai di Kasir</div>
                  <div class="text-[11px] text-zinc-500 mt-0.5">Bayar saat tiba di loket</div>
                </div>
              </button>
            </div>
          </div>

          <!-- Submit Button (Sophisticated, Calm Charcoal, No Harsh Clashing Yellow) -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full mt-4 py-3.5 rounded-full font-bold bg-zinc-900 hover:bg-zinc-800 text-white active:scale-[0.98] transition-all text-xs shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <span v-if="loading">Menerbitkan Antrian...</span>
            <span v-else class="flex items-center gap-1.5">
              Daftar & Terbitkan Nomor Antrian <ArrowRight class="w-3.5 h-3.5" />
            </span>
          </button>
        </form>
      </div>

    </div>
  </div>
</template>
