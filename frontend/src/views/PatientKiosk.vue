<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { QrCode, Banknote, CheckCircle2, Clock, Sparkles, RefreshCw, ArrowRight } from 'lucide-vue-next'

const store = useClinicStore()

const nik = ref('3171012345670008')
const name = ref('Reza Ramadhan')
const phone = ref('081288776655')
const complaint = ref('Batuk pilek dan sakit tenggorokan')
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
  <div class="min-h-screen pt-24 pb-16 px-4 flex justify-center items-start">
    <div class="w-full max-w-md">

      <!-- View: Active Ticket -->
      <div v-if="activeTicket" class="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-300">
        <!-- Ticket Card (Boarding Pass style) -->
        <div class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] relative overflow-hidden">
          <div class="absolute -top-12 -right-12 w-32 h-32 bg-amber-200/50 rounded-full blur-2xl"></div>

          <!-- Header -->
          <div class="flex items-center justify-between pb-6 border-b border-dashed border-zinc-200">
            <div>
              <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200/60">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                {{ activeTicket.payment_status === 'PAID' ? 'Antrian Aktif' : 'Menunggu Kasir' }}
              </span>
              <h2 class="text-lg font-bold text-zinc-950 mt-1">Klinik Sehat Modern</h2>
            </div>
            <button
              @click="refreshStatus"
              :class="['p-2 rounded-full hover:bg-stone-100 text-zinc-500 transition-transform active:scale-90', loading && 'animate-spin']"
              title="Perbarui Status"
            >
              <RefreshCw class="w-4 h-4" />
            </button>
          </div>

          <!-- Queue Number Display (Brandlos Bold) -->
          <div class="py-8 text-center">
            <span class="text-xs font-semibold uppercase tracking-wider text-zinc-400">Nomor Antrian Anda</span>
            <div class="text-6xl sm:text-7xl font-extrabold tracking-tight text-zinc-950 my-2">
              {{ activeTicket.queue_number }}
            </div>
            <p class="text-sm font-medium text-zinc-600">
              <span v-if="activeTicket.payment_status === 'PENDING'" class="text-amber-600 font-semibold">
                Silakan lakukan pembayaran tunai di kasir untuk memanggil dokter
              </span>
              <span v-else-if="aheadCount > 0" class="text-emerald-700">
                Ada <strong>{{ aheadCount }}</strong> pasien di depan Anda
              </span>
              <span v-else class="text-blue-600 font-semibold">
                Giliran Anda! Segera menuju ke Ruang Periksa Dokter
              </span>
            </p>
          </div>

          <!-- Details Breakdown -->
          <div class="space-y-3 pt-6 border-t border-dashed border-zinc-200 text-sm">
            <div class="flex justify-between">
              <span class="text-zinc-500">Nama Pasien</span>
              <span class="font-semibold text-zinc-900">{{ activeTicket.patient_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-500">Keluhan</span>
              <span class="font-medium text-zinc-800 text-right max-w-[200px] truncate">{{ activeTicket.complaint }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-zinc-500">Metode Bayar</span>
              <span class="inline-flex items-center gap-1 font-semibold text-xs px-2.5 py-0.5 rounded-full bg-stone-100 text-zinc-800">
                {{ activeTicket.payment_method }}
                <span v-if="activeTicket.payment_status === 'PAID'" class="text-emerald-600 ml-1">✓ Lunas</span>
                <span v-else class="text-amber-600 ml-1">⏳ Pending</span>
              </span>
            </div>
            <div v-if="activeTicket.patient_ihs" class="flex justify-between items-center pt-2">
              <span class="text-zinc-500 text-xs">SATUSEHAT IHS</span>
              <span class="text-xs font-mono font-bold text-blue-700 bg-blue-50 px-2 py-0.5 rounded-full border border-blue-200/50">
                {{ activeTicket.patient_ihs }}
              </span>
            </div>
          </div>
        </div>

        <button
          @click="resetForm"
          class="w-full py-3.5 rounded-full text-xs font-bold text-zinc-500 hover:text-zinc-900 bg-white border border-black/5 hover:border-black/10 transition-all active:scale-[0.98]"
        >
          Daftar Pasien Lain / Ambil Antrian Baru
        </button>
      </div>

      <!-- View: Registration Form -->
      <div v-else class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] animate-in fade-in duration-200">
        <!-- Title & Intro -->
        <div class="mb-6">
          <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-blue-600 text-white mb-2 shadow-sm">
            <Sparkles class="w-3 h-3" /> Pendaftaran Mandiri
          </span>
          <h1 class="text-2xl font-extrabold text-zinc-950 tracking-tight">
            Ambil Antrian Online
          </h1>
          <p class="text-sm text-zinc-500 mt-1">
            Tanpa perlu chat WA manual. Masukkan data dan dapatkan nomor antrian instan.
          </p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Error alert -->
          <div v-if="errorMsg" class="p-3 rounded-2xl bg-rose-50 text-rose-700 text-xs font-medium border border-rose-200">
            {{ errorMsg }}
          </div>

          <!-- NIK -->
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
              NIK (Nomor Induk Kependudukan)
            </label>
            <input
              v-model="nik"
              type="text"
              maxlength="16"
              required
              placeholder="Contoh: 3171012345670008"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all"
            />
          </div>

          <!-- Nama Lengkap -->
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
              Nama Lengkap
            </label>
            <input
              v-model="name"
              type="text"
              required
              placeholder="Nama sesuai KTP"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all"
            />
          </div>

          <!-- No WhatsApp -->
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
              Nomor WhatsApp
            </label>
            <input
              v-model="phone"
              type="tel"
              required
              placeholder="Contoh: 081288776655"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all"
            />
          </div>

          <!-- Keluhan Singkat -->
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
              Keluhan Gejala
            </label>
            <textarea
              v-model="complaint"
              rows="2"
              required
              placeholder="Ceritakan keluhan yang dirasakan..."
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all resize-none"
            ></textarea>
          </div>

          <!-- Pilihan Pembayaran (FIFO Rule) -->
          <div class="pt-2">
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-2">
              Biaya Pendaftaran: <strong class="text-zinc-950 font-bold">Rp 35.000</strong>
            </label>
            <div class="grid grid-cols-2 gap-2.5">
              <!-- QRIS Option -->
              <button
                type="button"
                @click="paymentMethod = 'QRIS'"
                :class="[
                  'p-3.5 rounded-2xl border text-left flex flex-col justify-between transition-all active:scale-95',
                  paymentMethod === 'QRIS'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-sm'
                    : 'border-black/5 bg-[#F9F9F7] text-zinc-800 hover:border-black/15'
                ]"
              >
                <div class="flex items-center justify-between mb-2">
                  <QrCode class="w-5 h-5" />
                  <span v-if="paymentMethod === 'QRIS'" class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-amber-400 text-zinc-950">REKOMENDASI</span>
                </div>
                <div>
                  <div class="text-xs font-bold">QRIS / Instant</div>
                  <div :class="['text-[11px]', paymentMethod === 'QRIS' ? 'text-zinc-300' : 'text-zinc-500']">Antrian langsung aktif</div>
                </div>
              </button>

              <!-- CASH Option -->
              <button
                type="button"
                @click="paymentMethod = 'CASH'"
                :class="[
                  'p-3.5 rounded-2xl border text-left flex flex-col justify-between transition-all active:scale-95',
                  paymentMethod === 'CASH'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-sm'
                    : 'border-black/5 bg-[#F9F9F7] text-zinc-800 hover:border-black/15'
                ]"
              >
                <div class="flex items-center justify-between mb-2">
                  <Banknote class="w-5 h-5" />
                </div>
                <div>
                  <div class="text-xs font-bold">Tunai di Kasir</div>
                  <div :class="['text-[11px]', paymentMethod === 'CASH' ? 'text-zinc-300' : 'text-zinc-500']">Bayar saat tiba di klinik</div>
                </div>
              </button>
            </div>
          </div>

          <!-- Submit Button (Tactile Brandlos Button) -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full mt-4 py-4 rounded-full font-bold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-[0.98] transition-all duration-150 shadow-sm flex items-center justify-center gap-2 text-sm disabled:opacity-50"
          >
            <span v-if="loading">Memproses Pendaftaran...</span>
            <span v-else class="flex items-center gap-1.5">
              Konfirmasi & Ambil Nomor Antrian <ArrowRight class="w-4 h-4" />
            </span>
          </button>
        </form>
      </div>

    </div>
  </div>
</template>
