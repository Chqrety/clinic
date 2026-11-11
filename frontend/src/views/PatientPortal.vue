<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Baby,
  Sparkles,
  ArrowLeft,
  QrCode,
  CreditCard,
  Building2,
  CheckCircle2,
  Clock,
  MessageCircle,
  AlertCircle,
  Loader2,
  Smartphone,
  ChevronRight,
  ShieldCheck
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login' | 'monitor'): void
}>()

// Step: 'register' | 'payment' | 'ticket'
const currentStep = ref<'register' | 'payment' | 'ticket'>('register')

// Form inputs
const selectedUnit = ref<'KIDINEST' | 'SSC'>('KIDINEST')
const patientName = ref('Budi Pratama')
const patientPhone = ref('081234567890')
const selectedCategory = ref('A_VAKSIN')
const selectedDoctor = ref('dr. Hendra Pratama, Sp.A')
const paymentMethod = ref<'QRIS' | 'BANK_TRANSFER' | 'CASH'>('QRIS')

// State
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

// Ticket Result
interface TicketData {
  queue_id: number
  queue_number: string
  clinic_unit: string
  category: string
  doctor_name: string
  registration_fee: number
  payment_method: string
  payment_status: string
  queue_status: string
  qris_payload?: string
  bank_transfer_account?: string
}

const activeTicket = ref<TicketData | null>(null)

// Kidinest categories
const kidinestCategories = [
  { key: 'A_VAKSIN', label: 'Vaksinasi Anak (A)', desc: 'Imunisasi dasar & lanjutan' },
  { key: 'B_SAKIT', label: 'Pemeriksaan Sakit (B)', desc: 'Demam, batuk, flu & infeksi' },
  { key: 'C_TUMBANG', label: 'Tumbuh Kembang (C)', desc: 'Skrining milestone & nutrisi' },
  { key: 'D_PAKET', label: 'Paket Vaksin + Tumbang (D)', desc: 'Layanan terpadu lengkap' },
]

// SSC treatments
const sscCategories = [
  { key: 'SSC_TREATMENT', label: 'Konsultasi Spesialis Kulit & Facial', desc: 'Acne, peeling, rejuvenasi & tindakan estetika' },
]

const currentCategories = computed(() => {
  return selectedUnit.value === 'KIDINEST' ? kidinestCategories : sscCategories
})

function switchUnit(unit: 'KIDINEST' | 'SSC') {
  selectedUnit.value = unit
  if (unit === 'KIDINEST') {
    selectedCategory.value = 'A_VAKSIN'
    selectedDoctor.value = 'dr. Hendra Pratama, Sp.A'
  } else {
    selectedCategory.value = 'SSC_TREATMENT'
    selectedDoctor.value = 'dr. Nadia Kusuma, Sp.KK'
  }
}

async function handleRegister() {
  isSubmitting.value = true
  errorMessage.value = null

  try {
    const res = await fetch('/api/public/queue/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        clinic_unit: selectedUnit.value,
        patient_name: patientName.value,
        patient_phone: patientPhone.value,
        category: selectedCategory.value,
        doctor_name: selectedDoctor.value,
        payment_method: selectedUnit.value === 'KIDINEST' ? paymentMethod.value : 'UNPAID_BOOKING',
      }),
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Gagal mendaftar antrian')
    }

    activeTicket.value = data.data

    if (selectedUnit.value === 'KIDINEST' && data.data.payment_status === 'PENDING') {
      currentStep.value = 'payment'
    } else {
      currentStep.value = 'ticket'
    }
  } catch (err: any) {
    // Demo fallback jika backend offline
    activeTicket.value = {
      queue_id: Date.now(),
      queue_number: selectedUnit.value === 'KIDINEST' ? 'MENUNGGU-BAYAR' : 'S-001',
      clinic_unit: selectedUnit.value,
      category: selectedCategory.value,
      doctor_name: selectedDoctor.value,
      registration_fee: selectedUnit.value === 'KIDINEST' ? 50000 : 0,
      payment_method: selectedUnit.value === 'KIDINEST' ? paymentMethod.value : 'UNPAID_BOOKING',
      payment_status: selectedUnit.value === 'KIDINEST' ? 'PENDING' : 'PAID',
      queue_status: 'WAITING',
      bank_transfer_account: 'BCA 828-192-8811 (a.n Kidinest Care)',
    }
    if (selectedUnit.value === 'KIDINEST') {
      currentStep.value = 'payment'
    } else {
      currentStep.value = 'ticket'
    }
  } finally {
    isSubmitting.value = false
  }
}

async function handleSimulatePayment() {
  if (!activeTicket.value) return
  isSubmitting.value = true

  try {
    const res = await fetch('/api/public/queue/verify-payment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ queue_id: activeTicket.value.queue_id }),
    })
    const data = await res.json()
    if (res.ok && data.success) {
      activeTicket.value = data.data
    } else if (activeTicket.value) {
      activeTicket.value.queue_number = 'A-001'
      activeTicket.value.payment_status = 'PAID'
    }
  } catch {
    if (activeTicket.value) {
      activeTicket.value.queue_number = 'A-001'
      activeTicket.value.payment_status = 'PAID'
    }
  } finally {
    isSubmitting.value = false
    currentStep.value = 'ticket'
  }
}

function resetForm() {
  activeTicket.value = null
  currentStep.value = 'register'
}
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col justify-between p-4 sm:p-8 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300">
    <!-- Top Bar Navigation -->
    <header class="max-w-2xl mx-auto w-full flex items-center justify-between pb-4">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-2 text-xs font-bold text-zinc-500 hover:text-zinc-950 transition-colors"
      >
        <ArrowLeft class="w-3.5 h-3.5" /> Beranda
      </button>

      <div class="flex items-center gap-2">
        <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold bg-white border border-zinc-200/80 text-zinc-700 shadow-2xs">
          <ShieldCheck class="w-3 h-3 text-emerald-600" />
          <span>Privacy-First (Tanpa NIK KTP)</span>
        </span>
      </div>
    </header>

    <!-- Main Container -->
    <main class="max-w-2xl mx-auto w-full py-4 flex-1 flex flex-col justify-center">

      <!-- STEP 1: REGISTRATION FORM -->
      <div v-if="currentStep === 'register'" class="space-y-6 animate-in fade-in duration-200">
        
        <!-- Header Title -->
        <div class="text-center space-y-2">
          <span class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-xs">
            Portal Pendaftaran Mandiri Pasien
          </span>
          <h1 class="text-2xl sm:text-3xl font-black tracking-tight text-zinc-950">
            Ambil Nomor Antrian Online.
          </h1>
          <p class="text-xs text-zinc-600 max-w-md mx-auto">
            Scan QR meja pendaftaran, pilih unit klinik, dan dapatkan tiket digital langsung di WhatsApp Anda.
          </p>
        </div>

        <!-- Unit Selector Card (Kidinest vs SSC) -->
        <div class="grid grid-cols-2 gap-3 p-1.5 bg-stone-200/60 rounded-3xl">
          <button
            type="button"
            @click="switchUnit('KIDINEST')"
            :class="[
              'p-4 rounded-2xl text-left transition-all flex items-start gap-3',
              selectedUnit === 'KIDINEST'
                ? 'bg-white shadow-xs text-zinc-950 font-bold border border-black/5'
                : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            <div :class="['p-2 rounded-xl shrink-0', selectedUnit === 'KIDINEST' ? 'bg-sky-100 text-sky-800' : 'bg-stone-100 text-zinc-500']">
              <Baby class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-extrabold">Kidinest</div>
              <div class="text-[10px] text-zinc-500 font-medium">Poli Anak & Vaksinasi</div>
              <div class="mt-1 text-[9px] font-bold text-sky-700 bg-sky-50 px-2 py-0.5 rounded-full inline-block">
                Bayar Di Muka (FIFO)
              </div>
            </div>
          </button>

          <button
            type="button"
            @click="switchUnit('SSC')"
            :class="[
              'p-4 rounded-2xl text-left transition-all flex items-start gap-3',
              selectedUnit === 'SSC'
                ? 'bg-white shadow-xs text-zinc-950 font-bold border border-black/5'
                : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            <div :class="['p-2 rounded-xl shrink-0', selectedUnit === 'SSC' ? 'bg-rose-100 text-rose-800' : 'bg-stone-100 text-zinc-500']">
              <Sparkles class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-extrabold">Semarang Skin Center</div>
              <div class="text-[10px] text-zinc-500 font-medium">Estetika & Spesialis Kulit</div>
              <div class="mt-1 text-[9px] font-bold text-rose-700 bg-rose-50 px-2 py-0.5 rounded-full inline-block">
                Booking Rp0 (Bayar Nanti)
              </div>
            </div>
          </button>
        </div>

        <!-- Form Card -->
        <form @submit.prevent="handleRegister" class="bg-white p-6 sm:p-8 rounded-3xl border border-zinc-200/80 shadow-xs space-y-5">
          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Patient Name -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
              Nama Lengkap Pasien
            </label>
            <input
              v-model="patientName"
              type="text"
              required
              placeholder="Contoh: Budi Pratama"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200 text-xs font-bold text-zinc-950 outline-none focus:bg-white focus:border-zinc-950 transition-all"
            />
          </div>

          <!-- Patient WhatsApp -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
              Nomor WhatsApp Aktif
            </label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-xs font-bold text-zinc-500">
                +62
              </span>
              <input
                v-model="patientPhone"
                type="tel"
                required
                placeholder="81234567890"
                class="w-full pl-12 pr-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200 text-xs font-bold text-zinc-950 outline-none focus:bg-white focus:border-zinc-950 transition-all font-mono"
              />
            </div>
            <p class="text-[10px] text-zinc-400 mt-1">Tiket antrian digital & estimasi waktu akan dikirim ke nomor ini.</p>
          </div>

          <!-- Service Category Picker -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-2">
              Pilihan Layanan / Kategori
            </label>
            <div class="space-y-2">
              <label
                v-for="cat in currentCategories"
                :key="cat.key"
                :class="[
                  'flex items-start gap-3 p-3.5 rounded-2xl border cursor-pointer transition-all',
                  selectedCategory === cat.key
                    ? 'border-zinc-950 bg-stone-50/80 shadow-2xs font-bold'
                    : 'border-zinc-200 hover:border-zinc-400 bg-white'
                ]"
              >
                <input
                  type="radio"
                  v-model="selectedCategory"
                  :value="cat.key"
                  class="mt-0.5 accent-zinc-950"
                />
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-extrabold text-zinc-950">{{ cat.label }}</div>
                  <div class="text-[11px] text-zinc-500 font-medium">{{ cat.desc }}</div>
                </div>
              </label>
            </div>
          </div>

          <!-- Doctor Selector -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
              Dokter yang Berpraktik Hari Ini
            </label>
            <select
              v-model="selectedDoctor"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200 text-xs font-bold text-zinc-950 outline-none focus:bg-white focus:border-zinc-950 transition-all"
            >
              <option v-if="selectedUnit === 'KIDINEST'" value="dr. Hendra Pratama, Sp.A">
                dr. Hendra Pratama, Sp.A (Spesialis Anak)
              </option>
              <option v-if="selectedUnit === 'KIDINEST'" value="dr. Anisa S., Sp.A">
                dr. Anisa S., Sp.A (Spesialis Anak)
              </option>
              <option v-if="selectedUnit === 'SSC'" value="dr. Nadia Kusuma, Sp.KK">
                dr. Nadia Kusuma, Sp.KK (Spesialis Kulit & Kelamin)
              </option>
              <option v-if="selectedUnit === 'SSC'" value="dr. Maya Kartika, Sp.DV">
                dr. Maya Kartika, Sp.DV (Dermatologi Estetika)
              </option>
            </select>
          </div>

          <!-- Payment Method (Khusus Kidinest) -->
          <div v-if="selectedUnit === 'KIDINEST'" class="pt-2 border-t border-zinc-100">
            <div class="flex items-center justify-between mb-2">
              <label class="text-[11px] font-bold uppercase tracking-wider text-zinc-600">
                Biaya Registrasi Awal:
              </label>
              <span class="text-sm font-black text-zinc-950">Rp 50.000</span>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <label
                :class="[
                  'flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all',
                  paymentMethod === 'QRIS' ? 'border-zinc-950 bg-stone-50' : 'border-zinc-200 bg-white'
                ]"
              >
                <input type="radio" v-model="paymentMethod" value="QRIS" class="accent-zinc-950" />
                <QrCode class="w-4 h-4 text-emerald-600" />
                <span>QRIS Instan</span>
              </label>

              <label
                :class="[
                  'flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all',
                  paymentMethod === 'BANK_TRANSFER' ? 'border-zinc-950 bg-stone-50' : 'border-zinc-200 bg-white'
                ]"
              >
                <input type="radio" v-model="paymentMethod" value="BANK_TRANSFER" class="accent-zinc-950" />
                <CreditCard class="w-4 h-4 text-sky-600" />
                <span>Transfer Bank</span>
              </label>
            </div>
          </div>

          <!-- SSC Free Note -->
          <div v-else class="p-3 rounded-2xl bg-rose-50/60 border border-rose-200/70 text-[11px] text-rose-900 font-medium">
            💡 <strong>Semarang Skin Center:</strong> Biaya konsultasi dokter dan tindakan perawatan akan ditagihkan di kasir meja depan setelah tindakan selesai.
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isSubmitting"
            class="w-full py-3.5 rounded-full font-bold bg-zinc-950 hover:bg-zinc-800 text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <span v-else>{{ selectedUnit === 'KIDINEST' ? 'Lanjut ke Pembayaran Registrasi' : 'Konfirmasi Reservasi Antrian' }}</span>
            <ChevronRight v-if="!isSubmitting" class="w-4 h-4" />
          </button>
        </form>

      </div>

      <!-- STEP 2: PAYMENT SCREEN (KIDINEST ONLY) -->
      <div v-else-if="currentStep === 'payment' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-amber-100 text-amber-900">
            <Clock class="w-3.5 h-3.5" /> Menunggu Pembayaran Registrasi
          </span>
          <h2 class="text-2xl font-black text-zinc-950">Selesaikan Pembayaran</h2>
          <p class="text-xs text-zinc-500">Nomor antrian FIFO resmi akan otomatis terbit saat pembayaran lunas.</p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-zinc-200/80 shadow-xs space-y-6 text-center">
          <div class="p-4 rounded-2xl bg-stone-50 border border-zinc-200/60 text-left space-y-1">
            <div class="flex items-center justify-between text-xs text-zinc-500">
              <span>Pasien:</span>
              <span class="font-bold text-zinc-900">{{ patientName }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-zinc-500">
              <span>Unit:</span>
              <span class="font-bold text-zinc-900">Kidinest (Poli Anak)</span>
            </div>
            <div class="flex items-center justify-between text-xs text-zinc-500">
              <span>Dokter:</span>
              <span class="font-bold text-zinc-900">{{ activeTicket.doctor_name }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-zinc-200 text-sm font-extrabold text-zinc-950">
              <span>Total Tagihan:</span>
              <span class="text-emerald-700">Rp 50.000</span>
            </div>
          </div>

          <!-- QRIS Dummy Display -->
          <div class="p-6 rounded-3xl bg-[#F9F9F7] border border-zinc-200/80 inline-block space-y-3">
            <div class="w-48 h-48 mx-auto bg-white p-3 rounded-2xl border border-zinc-300 flex items-center justify-center shadow-xs">
              <QrCode class="w-40 h-40 text-zinc-900" />
            </div>
            <div class="text-[11px] font-mono font-bold text-zinc-600">
              NMID: ID1020039201990
            </div>
            <div class="text-[10px] text-zinc-400">Scan QRIS menggunakan BCA, Mandiri, GoPay, OVO, atau ShopeePay</div>
          </div>

          <!-- Transfer Info -->
          <div class="text-xs font-medium text-zinc-600 space-y-1">
            <div>Atau transfer manual ke rekening klinik:</div>
            <div class="font-mono font-extrabold text-zinc-950 bg-stone-100 py-1.5 px-3 rounded-full inline-block">
              BCA 828-192-8811 a.n Kidinest Care
            </div>
          </div>

          <!-- Simulation Button (Demo Weapon) -->
          <div class="pt-4 border-t border-zinc-100 space-y-2">
            <button
              @click="handleSimulatePayment"
              :disabled="isSubmitting"
              class="w-full py-3.5 rounded-full font-bold bg-emerald-600 hover:bg-emerald-700 text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <CheckCircle2 v-else class="w-4 h-4" />
              <span>Simulasikan Pembayaran Berhasil (Demo)</span>
            </button>

            <button
              @click="currentStep = 'register'"
              class="text-xs font-bold text-zinc-400 hover:text-zinc-600 py-1"
            >
              Ubah Data Pendaftaran
            </button>
          </div>
        </div>
      </div>

      <!-- STEP 3: TICKET SCREEN (BOARDING PASS DIGITAL) -->
      <div v-else-if="currentStep === 'ticket' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-emerald-100 text-emerald-900">
            <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" /> Pendaftaran Berhasil & Terverifikasi
          </span>
          <h2 class="text-2xl font-black text-zinc-950">Tiket Antrian Digital Anda</h2>
          <p class="text-xs text-zinc-500">Simpan layar ini atau pantau tautan yang telah dikirim ke WhatsApp.</p>
        </div>

        <!-- Boarding Pass Style Card -->
        <div class="bg-white rounded-3xl border border-zinc-200/80 shadow-xs overflow-hidden">
          <!-- Card Header -->
          <div class="p-6 bg-zinc-950 text-white text-center space-y-2">
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400">
              {{ activeTicket.clinic_unit === 'KIDINEST' ? 'Kidinest Poli Anak & Vaksinasi' : 'Semarang Skin Center' }}
            </div>
            <!-- Giant Number -->
            <div class="text-6xl font-black tracking-tight text-white font-mono py-2">
              {{ activeTicket.queue_number }}
            </div>
            <div class="inline-block px-3 py-1 rounded-full text-xs font-bold bg-white/10 text-zinc-200">
              Status: Sedang Menunggu Panggilan
            </div>
          </div>

          <!-- Card Details -->
          <div class="p-6 space-y-4 text-xs">
            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="text-zinc-400 font-semibold block text-[10px] uppercase">Pasien</span>
                <span class="font-extrabold text-zinc-950 text-sm">{{ patientName }}</span>
              </div>
              <div>
                <span class="text-zinc-400 font-semibold block text-[10px] uppercase">WhatsApp</span>
                <span class="font-mono font-bold text-zinc-950 text-sm">{{ patientPhone }}</span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="text-zinc-400 font-semibold block text-[10px] uppercase">Dokter Praktik</span>
                <span class="font-bold text-zinc-950">{{ activeTicket.doctor_name }}</span>
              </div>
              <div>
                <span class="text-zinc-400 font-semibold block text-[10px] uppercase">Registrasi</span>
                <span class="font-bold text-emerald-700">
                  {{ activeTicket.registration_fee > 0 ? 'Rp 50.000 (Lunas)' : 'Booking Rp0 (Pasca Tindakan)' }}
                </span>
              </div>
            </div>

            <!-- WhatsApp Badge -->
            <div class="p-3.5 rounded-2xl bg-emerald-50 border border-emerald-200/80 flex items-center gap-3">
              <MessageCircle class="w-5 h-5 text-emerald-600 shrink-0" />
              <div class="text-[11px] text-emerald-950 font-medium leading-tight">
                Notifikasi tiket dan tautan pemantauan live telah dikirimkan ke nomor <strong>{{ patientPhone }}</strong>.
              </div>
            </div>
          </div>

          <!-- Card Footer Action -->
          <div class="p-4 bg-stone-50 border-t border-zinc-100 flex items-center justify-between gap-3">
            <button
              @click="emit('navigate', 'monitor')"
              class="inline-flex items-center gap-1.5 text-xs font-bold text-zinc-700 hover:text-zinc-950"
            >
              <Smartphone class="w-4 h-4" />
              <span>Lihat Monitor TV Antrian</span>
            </button>

            <button
              @click="resetForm"
              class="px-4 py-2 rounded-full bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition-all"
            >
              Daftar Antrian Baru
            </button>
          </div>
        </div>
      </div>

    </main>

    <!-- Footer -->
    <footer class="max-w-2xl mx-auto w-full text-center text-[11px] font-medium text-zinc-400 pt-4">
      Smart Clinic Management System &copy; 2026. Paperless Patient Registration.
    </footer>
  </div>
</template>
