<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem, type Invoice } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { QrCode, Banknote, RefreshCw, ArrowRight, ArrowLeft, Check, ShieldCheck, Receipt, Printer, X } from 'lucide-vue-next'

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
const patientInvoice = ref<Invoice | null>(null)
const showInvoiceModal = ref(false)

async function viewPatientInvoice() {
  if (!activeTicket.value) return
  loading.value = true
  try {
    patientInvoice.value = await api.getInvoiceByQueue(activeTicket.value.id)
    showInvoiceModal.value = true
  } catch (err: any) {
    alert(err.message || 'Invoice belum tersedia')
  } finally {
    loading.value = false
  }
}

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

    <!-- Top Navigation (Fully Responsive Header) -->
    <header class="w-full max-w-6xl mx-auto px-4 sm:px-8 lg:px-10 pt-6 sm:pt-8 pb-3 sm:pb-4 flex items-center justify-between gap-3">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-1.5 text-xs font-bold text-zinc-500 hover:text-zinc-950 transition-colors"
      >
        <ArrowLeft class="w-3.5 h-3.5 sm:w-4 sm:h-4 flex-shrink-0" />
        <span>Beranda</span>
      </button>

      <span class="text-[10px] sm:text-xs font-extrabold uppercase tracking-widest text-zinc-400 truncate">
        Pendaftaran Pasien Mandiri
      </span>
    </header>

    <!-- Main Content Area -->
    <main class="flex-1 max-w-6xl mx-auto w-full px-4 sm:px-8 lg:px-10 py-4 sm:py-8">

      <!-- ========================================== -->
      <!-- VIEW 1: ACTIVE DIGITAL TICKET PASS        -->
      <!-- ========================================== -->
      <div v-if="activeTicket" class="max-w-md mx-auto space-y-4 sm:space-y-6 animate-in fade-in duration-200">
        <!-- Physical Boarding Slip / Appointment Receipt -->
        <div class="bg-white rounded-3xl p-6 sm:p-8 border border-black/10 shadow-[0_4px_24px_rgba(0,0,0,0.04)]">

          <!-- Ticket Header -->
          <div class="flex items-center justify-between pb-4 sm:pb-5 border-b border-zinc-100 gap-2">
            <div>
              <span class="text-[10px] sm:text-[11px] font-extrabold uppercase tracking-widest text-zinc-400 block">Klinik Sehat Modern</span>
              <h2 class="text-sm sm:text-base font-extrabold text-zinc-950 mt-0.5">Tiket Antrian Digital</h2>
            </div>
            <button
              @click="refreshStatus"
              :class="['p-2 rounded-full hover:bg-stone-100 text-zinc-500 transition-transform active:scale-90', loading && 'animate-spin']"
              title="Perbarui Status Antrian"
            >
              <RefreshCw class="w-4 h-4" />
            </button>
          </div>

          <!-- Queue Number Display -->
          <div class="py-6 sm:py-8 text-center">
            <span class="text-[10px] sm:text-xs font-bold uppercase tracking-widest text-zinc-400">Nomor Giliran Anda</span>
            <div class="text-5xl sm:text-6xl md:text-7xl font-black tracking-tight text-zinc-950 font-mono my-2 sm:my-3">
              {{ activeTicket.queue_number }}
            </div>

            <!-- Status Indicator Pill -->
            <div
              class="inline-flex items-center gap-1.5 px-3 py-1 sm:px-3.5 sm:py-1.5 rounded-full text-[11px] sm:text-xs font-bold mt-1"
              :class="[
                activeTicket.queue_status === 'COMPLETED'
                  ? 'bg-blue-50 text-blue-800 border border-blue-200/60'
                  : (activeTicket.payment_status === 'PAID'
                      ? 'bg-emerald-50 text-emerald-800 border border-emerald-200/60'
                      : 'bg-stone-100 text-zinc-700 border border-zinc-200')
              ]"
            >
              <span
                class="w-2 h-2 rounded-full flex-shrink-0"
                :class="[
                  activeTicket.queue_status === 'COMPLETED'
                    ? 'bg-blue-600'
                    : (activeTicket.payment_status === 'PAID' ? 'bg-emerald-600 animate-pulse' : 'bg-zinc-400')
                ]"
              ></span>
              <span class="truncate">
                {{
                  activeTicket.queue_status === 'COMPLETED'
                    ? 'Pelayanan Medis Selesai'
                    : (activeTicket.payment_status === 'PAID' ? 'Antrian Aktif Terjadwal' : 'Menunggu Kasir')
                }}
              </span>
            </div>

            <!-- Wait Notice -->
            <p class="text-xs text-zinc-600 mt-3 font-medium leading-relaxed px-2">
              <span v-if="activeTicket.queue_status === 'COMPLETED'" class="text-zinc-950 font-bold">
                Pemeriksaan dokter dan penyerahan obat di apotek telah selesai. Terima kasih atas kunjungan Anda.
              </span>
              <span v-else-if="activeTicket.payment_status === 'PENDING'">
                Silakan lakukan pembayaran di kasir loket untuk mengaktifkan giliran panggilan dokter.
              </span>
              <span v-else-if="aheadCount > 0">
                Ada <strong>{{ aheadCount }} pasien</strong> di depan Anda. Harap menunggu di ruang tunggu.
              </span>
              <span v-else class="text-zinc-950 font-extrabold">
                Giliran Anda! Silakan langsung menuju ke Ruang Periksa Dokter.
              </span>
            </p>
          </div>

          <!-- Structured Details Rows -->
          <div class="pt-4 sm:pt-6 border-t border-zinc-100 space-y-2.5 sm:space-y-3 text-xs">
            <div class="flex justify-between items-center py-0.5">
              <span class="text-zinc-400">Nama Pasien</span>
              <span class="font-extrabold text-zinc-900 truncate max-w-[180px]">{{ activeTicket.patient_name }}</span>
            </div>
            <div class="flex justify-between items-center py-0.5 border-t border-zinc-50">
              <span class="text-zinc-400">NIK KTP</span>
              <span class="font-mono font-semibold text-zinc-700">{{ activeTicket.patient_nik }}</span>
            </div>
            <div class="flex justify-between items-center py-0.5 border-t border-zinc-50">
              <span class="text-zinc-400">Keluhan</span>
              <span class="font-medium text-zinc-800 text-right truncate max-w-[180px] sm:max-w-[220px]">{{ activeTicket.complaint }}</span>
            </div>
            <div class="flex justify-between items-center py-0.5 border-t border-zinc-50">
              <span class="text-zinc-400">Pembayaran</span>
              <span class="font-bold text-zinc-900">
                {{ activeTicket.payment_method }}
                <span v-if="activeTicket.payment_status === 'PAID'" class="text-emerald-700 font-semibold ml-1">(Lunas)</span>
                <span v-else class="text-zinc-500 font-normal ml-1">(Pending)</span>
              </span>
            </div>
            <div v-if="activeTicket.patient_ihs" class="flex justify-between items-center py-0.5 border-t border-zinc-50">
              <span class="text-zinc-400 flex items-center gap-1">
                <ShieldCheck class="w-3.5 h-3.5 text-blue-600 flex-shrink-0" /> SATUSEHAT
              </span>
              <span class="font-mono font-bold text-zinc-900">{{ activeTicket.patient_ihs }}</span>
            </div>
          </div>

          <!-- Action for Completed Queue: Digital Invoice Pass -->
          <div v-if="activeTicket.queue_status === 'COMPLETED'" class="pt-4 border-t border-zinc-100">
            <button
              @click="viewPatientInvoice"
              :disabled="loading"
              class="w-full py-3 rounded-full text-xs font-extrabold bg-blue-600 hover:bg-blue-700 text-white transition-all active:scale-95 shadow-sm flex items-center justify-center gap-2"
            >
              <Receipt class="w-3.5 h-3.5" />
              <span>Lihat Salinan Invoice Digital Pasien</span>
            </button>
          </div>
        </div>

        <button
          @click="resetForm"
          class="w-full py-3.5 rounded-full text-xs font-bold text-zinc-600 hover:text-zinc-950 bg-white border border-black/10 hover:border-black/20 transition-all active:scale-95 shadow-sm"
        >
          Daftar Pasien Lain / Ambil Antrian Baru
        </button>
      </div>

      <!-- ========================================== -->
      <!-- VIEW 2: FULL CANVAS FORM (NO CARD WRAPPER) -->
      <!-- ========================================== -->
      <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-14 items-start animate-in fade-in duration-200">

        <!-- Left Column: Context, Narrative & Live Pulse (5 cols) -->
        <div class="lg:col-span-5 space-y-4 sm:space-y-6 lg:sticky lg:top-24">
          <div>
            <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] sm:text-xs font-bold uppercase tracking-wider bg-blue-600 text-white mb-2 sm:mb-3 shadow-sm">
              Loket Mandiri
            </span>
            <h1 class="text-3xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight text-zinc-950 leading-tight">
              Ambil nomor antrian tanpa antri di loket.
            </h1>
            <p class="text-xs sm:text-sm text-zinc-600 leading-relaxed mt-2 sm:mt-3">
              Cukup masukkan NIK dan keluhan Anda. Nomor antrian FIFO resmi terbit seketika dan tersinkron ke ruang periksa dokter.
            </p>
          </div>

          <!-- Live Clinic Pulse (Responsive Grid on Mobile vs Vertical on Desktop) -->
          <div class="py-4 sm:py-6 border-t border-b border-black/10 grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-1 gap-3 sm:gap-4 lg:space-y-3">
            <div class="flex sm:flex-col lg:flex-row sm:items-start lg:items-center justify-between text-xs gap-1">
              <span class="text-zinc-500 font-medium text-[11px] sm:text-xs">Nomor Dipanggil:</span>
              <span class="font-mono font-black text-sm sm:text-base text-zinc-950">{{ activeCallingNum }}</span>
            </div>
            <div class="flex sm:flex-col lg:flex-row sm:items-start lg:items-center justify-between text-xs gap-1">
              <span class="text-zinc-500 font-medium text-[11px] sm:text-xs">Estimasi Waktu:</span>
              <span class="font-bold text-zinc-900">&plusmn; 8 menit / pasien</span>
            </div>
            <div class="flex sm:flex-col lg:flex-row sm:items-start lg:items-center justify-between text-xs gap-1">
              <span class="text-zinc-500 font-medium text-[11px] sm:text-xs">Dokter Siaga:</span>
              <span class="font-bold text-zinc-900 flex items-center gap-1 text-[11px] sm:text-xs">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 flex-shrink-0"></span>
                dr. Hendra Pratama
              </span>
            </div>
          </div>

          <div class="hidden lg:block text-xs text-zinc-400">
            Terhubung otomatis dengan Kamus Farmasi (KFA) dan platform SATUSEHAT Kemenkes RI.
          </div>
        </div>

        <!-- Right Column: Registration Form (Direct Canvas, 7 cols) -->
        <div class="lg:col-span-7">
          <form @submit.prevent="handleSubmit" class="space-y-4 sm:space-y-5">

            <!-- Error Notification -->
            <div v-if="errorMsg" class="p-3.5 sm:p-4 rounded-2xl bg-rose-50 text-rose-800 text-xs font-bold border border-rose-200">
              {{ errorMsg }}
            </div>

            <!-- NIK Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                1. NIK (Nomor Induk Kependudukan)
              </label>
              <input
                v-model="nik"
                type="text"
                maxlength="16"
                required
                placeholder="16 digit sesuai e-KTP"
                class="w-full px-4 py-3 sm:py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-base sm:text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- Nama Lengkap Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                2. Nama Lengkap Pasien
              </label>
              <input
                v-model="name"
                type="text"
                required
                placeholder="Nama lengkap sesuai KTP"
                class="w-full px-4 py-3 sm:py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-base sm:text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- WhatsApp Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                3. Nomor WhatsApp / Ponsel
              </label>
              <input
                v-model="phone"
                type="tel"
                required
                placeholder="0812xxxxxxxx"
                class="w-full px-4 py-3 sm:py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-base sm:text-sm font-semibold outline-none transition-all placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              />
            </div>

            <!-- Keluhan Field -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                4. Keluhan yang Dirasakan
              </label>
              <textarea
                v-model="complaint"
                rows="2"
                required
                placeholder="Tuliskan keluhan atau gejala sakit yang dirasakan..."
                class="w-full px-4 py-3 sm:py-3.5 rounded-2xl bg-white border border-black/10 focus:border-zinc-950 text-zinc-950 text-base sm:text-sm font-medium outline-none transition-all resize-none placeholder:text-zinc-400 placeholder:font-normal shadow-sm"
              ></textarea>
            </div>

            <!-- Payment Method Field -->
            <div class="pt-1 sm:pt-2">
              <div class="flex items-baseline justify-between mb-2">
                <label class="text-xs font-bold uppercase tracking-wider text-zinc-600">
                  5. Biaya & Pilihan Bayar
                </label>
                <span class="text-xs font-semibold text-zinc-500">
                  Biaya: <strong class="text-zinc-950 font-black text-sm">Rp 35.000</strong>
                </span>
              </div>

              <!-- Payment Grid (Stacked on Mobile, 2 Columns on Tablet/Desktop) -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 sm:gap-3">
                <!-- QRIS Option -->
                <button
                  type="button"
                  @click="paymentMethod = 'QRIS'"
                  :class="[
                    'p-3.5 sm:p-4 rounded-2xl text-left flex items-start justify-between transition-all active:scale-[0.99] shadow-sm',
                    paymentMethod === 'QRIS'
                      ? 'bg-white border-2 border-zinc-950'
                      : 'bg-white/60 border border-black/10 hover:border-black/25 text-zinc-600'
                  ]"
                >
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <QrCode class="w-4 h-4 text-zinc-950 flex-shrink-0" />
                      <span class="font-extrabold text-xs text-zinc-950">QRIS / Digital</span>
                    </div>
                    <div class="text-[11px] text-zinc-500">
                      Antrian langsung aktif seketika
                    </div>
                  </div>
                  <div v-if="paymentMethod === 'QRIS'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center text-[10px] flex-shrink-0">
                    ✓
                  </div>
                </button>

                <!-- CASH Option -->
                <button
                  type="button"
                  @click="paymentMethod = 'CASH'"
                  :class="[
                    'p-3.5 sm:p-4 rounded-2xl text-left flex items-start justify-between transition-all active:scale-[0.99] shadow-sm',
                    paymentMethod === 'CASH'
                      ? 'bg-white border-2 border-zinc-950'
                      : 'bg-white/60 border border-black/10 hover:border-black/25 text-zinc-600'
                  ]"
                >
                  <div class="space-y-1">
                    <div class="flex items-center gap-2">
                      <Banknote class="w-4 h-4 text-zinc-950 flex-shrink-0" />
                      <span class="font-extrabold text-xs text-zinc-950">Tunai di Kasir</span>
                    </div>
                    <div class="text-[11px] text-zinc-500">
                      Bayar saat tiba di loket klinik
                    </div>
                  </div>
                  <div v-if="paymentMethod === 'CASH'" class="w-4 h-4 rounded-full bg-zinc-950 text-white flex items-center justify-center text-[10px] flex-shrink-0">
                    ✓
                  </div>
                </button>
              </div>
            </div>

            <!-- Submit Button (Full-width, Touch Friendly on Mobile) -->
            <div class="pt-2 sm:pt-4">
              <button
                type="submit"
                :disabled="loading"
                class="w-full min-h-[50px] py-3.5 sm:py-4 rounded-full font-extrabold bg-zinc-950 hover:bg-zinc-800 text-white active:scale-[0.98] transition-all text-xs sm:text-sm shadow-md flex items-center justify-center gap-2 disabled:opacity-50"
              >
                <span v-if="loading">Menerbitkan Nomor Antrian...</span>
                <span v-else class="flex items-center gap-2">
                  <span>Konfirmasi & Ambil Antrian</span>
                  <ArrowRight class="w-4 h-4" />
                </span>
              </button>
            </div>

          </form>
        </div>

      </div>

    </main>

    <!-- Digital Invoice Modal for Patient -->
    <div
      v-if="showInvoiceModal && patientInvoice"
      class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 sm:p-8 space-y-6 shadow-2xl animate-in fade-in zoom-in-95 duration-150">
        <!-- Receipt Header -->
        <div class="flex items-start justify-between pb-4 border-b border-dashed border-zinc-200">
          <div>
            <span class="text-[10px] font-extrabold uppercase tracking-widest text-zinc-400 block">Salinan Resmi</span>
            <h3 class="text-base font-extrabold text-zinc-950 mt-0.5">KLINIK SEHAT MODERN</h3>
            <p class="text-[11px] text-zinc-500">Jl. Kesehatan No. 88 — Telepon: (021) 555-0199</p>
          </div>
          <button
            @click="showInvoiceModal = false"
            class="p-2 rounded-full hover:bg-stone-100 text-zinc-400 hover:text-zinc-950 transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Receipt Meta -->
        <div class="grid grid-cols-2 text-xs gap-y-1 text-zinc-600">
          <div>No. Invoice: <strong class="text-zinc-900 font-mono">{{ patientInvoice.invoice_number }}</strong></div>
          <div class="text-right">Tgl: {{ new Date(patientInvoice.created_at).toLocaleDateString('id-ID') }}</div>
          <div>Pasien: <strong class="text-zinc-900">{{ patientInvoice.patient_name }}</strong></div>
          <div class="text-right">NIK: {{ patientInvoice.patient_nik }}</div>
        </div>

        <!-- Line Items -->
        <div class="border-t border-b border-zinc-100 py-3 space-y-2 text-xs">
          <div class="flex justify-between">
            <span class="text-zinc-600">Biaya Administrasi & Pendaftaran</span>
            <span class="font-bold text-zinc-900">Rp {{ patientInvoice.registration_fee.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-zinc-600">Jasa Konsultasi Dokter Umum</span>
            <span class="font-bold text-zinc-900">Rp {{ patientInvoice.consultation_fee.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-zinc-600">Total Biaya Resep & Obat</span>
            <span class="font-bold text-zinc-900">Rp {{ patientInvoice.medication_fee.toLocaleString('id-ID') }}</span>
          </div>
        </div>

        <!-- Total -->
        <div class="flex justify-between items-baseline pt-1">
          <div>
            <span class="text-xs font-bold uppercase tracking-wider text-zinc-500 block">Total Tagihan Lunas</span>
            <span class="text-[11px] text-emerald-700 font-bold">Via {{ patientInvoice.payment_method }} (Sudah Lunas)</span>
          </div>
          <span class="text-2xl font-black text-zinc-950">
            Rp {{ patientInvoice.total_amount.toLocaleString('id-ID') }}
          </span>
        </div>

        <!-- Close / Print -->
        <div class="pt-2 flex gap-2">
          <button
            @click="showInvoiceModal = false"
            class="flex-1 py-3 rounded-full text-xs font-bold bg-zinc-950 hover:bg-zinc-800 text-white transition-all active:scale-95"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
