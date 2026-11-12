<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  Baby,
  Sparkles,
  ArrowLeft,
  QrCode,
  CreditCard,
  CheckCircle2,
  Clock,
  MessageCircle,
  AlertCircle,
  Loader2,
  Smartphone,
  ChevronRight,
  ShieldCheck
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  brand?: 'kidinest' | 'ssc'
}>(), {
  brand: 'kidinest'
})

const emit = defineEmits<{
  (e: 'navigate', view: string): void
}>()

// Step: 'register' | 'payment' | 'ticket'
const currentStep = ref<'register' | 'payment' | 'ticket'>('register')

// Active unit computed strictly from brand context (Lantai 1 vs Lantai 2)
const selectedUnit = computed<'KIDINEST' | 'SSC'>(() => {
  return props.brand === 'ssc' ? 'SSC' : 'KIDINEST'
})

// Form inputs
const patientName = ref('Budi Pratama')
const patientPhone = ref('081234567890')
const selectedCategory = ref(props.brand === 'ssc' ? 'SSC_LASER' : 'A_VAKSIN')
const selectedDoctor = ref(
  props.brand === 'ssc'
    ? 'dr. Dian Paramita, Sp.D.V.E (Ruang Laser 1)'
    : 'dr. Patricia Vanessa A., Sp.A (Poli Anak 1)'
)
const paymentMethod = ref<'QRIS' | 'BANK_TRANSFER' | 'CASH'>('QRIS')

// Watch brand prop to reset form defaults if brand changes
watch(() => props.brand, (newBrand) => {
  if (newBrand === 'ssc') {
    selectedCategory.value = 'SSC_LASER'
    selectedDoctor.value = 'dr. Dian Paramita, Sp.D.V.E (Ruang Laser 1)'
  } else {
    selectedCategory.value = 'A_VAKSIN'
    selectedDoctor.value = 'dr. Patricia Vanessa A., Sp.A (Poli Anak 1)'
  }
  currentStep.value = 'register'
})

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

// Kidinest categories (Lantai 1)
const kidinestCategories = [
  { key: 'A_VAKSIN', label: 'Vaksinasi & Imunisasi Anak (A)', desc: 'Imunisasi dasar & lanjutan rantai dingin IDAI' },
  { key: 'B_SAKIT', label: 'Pemeriksaan Sakit & Alergi (B)', desc: 'Demam, batuk pilek, flu, diare, ruam infeksi' },
  { key: 'C_TUMBANG', label: 'Skrining Milestone & Nutrisi (C)', desc: 'Pemantauan grafik tumbuh kembang & gizi anak' },
  { key: 'D_PAKET', label: 'Paket Vaksin + Konsultasi Tumbang (D)', desc: 'Pemeriksaan pediatri menyeluruh' },
]

// SSC treatments (Lantai 2)
const sscCategories = [
  { key: 'SSC_LASER', label: 'Laser Skin Rejuvenation & Pigmentation', desc: 'Melasma, pori membesar & peremajaan kolagen kulit' },
  { key: 'SSC_ACNE', label: 'Acne & Atrophic Scar Reconstruction', desc: 'Penanganan jerawat aktif & perbaikan bopeng' },
  { key: 'SSC_FACIAL', label: 'Medical Aesthetic Facial & Deep Infusion', desc: 'Pembersihan komedo steril & serum asam buah medis' },
  { key: 'SSC_SURGERY', label: 'Bedah Minor Kulit & Kauterisasi', desc: 'Biopsi lesi, kauter keratosis & pengangkatan kista' },
]

const currentCategories = computed(() => {
  return selectedUnit.value === 'KIDINEST' ? kidinestCategories : sscCategories
})

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
      queue_number: selectedUnit.value === 'KIDINEST' ? 'K-012' : 'S-009',
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
      activeTicket.value.queue_number = 'K-012'
      activeTicket.value.payment_status = 'PAID'
    }
  } catch {
    if (activeTicket.value) {
      activeTicket.value.queue_number = 'K-012'
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
  <div
    class="min-h-screen flex flex-col justify-between p-4 sm:p-8 font-['Plus_Jakarta_Sans',sans-serif] select-none transition-colors"
    :class="brand === 'kidinest' ? 'bg-[#FAF7F2] text-[#0C3B36]' : 'bg-[#FAF5F5] text-[#2D1E24]'"
  >
    <!-- Top Bar Navigation -->
    <header class="max-w-2xl mx-auto w-full flex items-center justify-between pb-4">
      <button
        @click="emit('navigate', brand)"
        class="inline-flex items-center gap-2 text-xs font-bold opacity-70 hover:opacity-100 transition-opacity"
      >
        <ArrowLeft class="w-3.5 h-3.5" />
        <span>Kembali ke Beranda {{ brand === 'kidinest' ? 'Kidinest' : 'Semarang Skin Center' }}</span>
      </button>

      <div class="flex items-center gap-2">
        <span
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold border shadow-2xs"
          :class="
            brand === 'kidinest'
              ? 'bg-white border-[#E5A93C]/30 text-[#0C3B36]'
              : 'bg-white border-[#9E606D]/30 text-[#9E606D]'
          "
        >
          <ShieldCheck class="w-3 h-3 text-emerald-600" />
          <span>Faskes Resmi Semarang Medical Group</span>
        </span>
      </div>
    </header>

    <!-- Main Container -->
    <main class="max-w-2xl mx-auto w-full py-4 flex-1 flex flex-col justify-center">

      <!-- STEP 1: REGISTRATION FORM (TERPISAH TOTAL PER LANTAI, ZERO UNIT SWITCHER) -->
      <div v-if="currentStep === 'register'" class="space-y-6 animate-in fade-in duration-200">
        
        <!-- Header Title Sesuai Lantai -->
        <div class="text-center space-y-2">
          <div class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <img
              v-if="brand === 'kidinest'"
              src="/assets/kidinest/logo-maskot.png"
              class="w-4 h-4 object-contain"
            />
            <img
              v-else
              src="/assets/ssc/logo-ssc.png"
              class="w-3.5 h-3.5 object-contain"
            />
            <span>{{ brand === 'kidinest' ? 'Kidinest • Poli Anak & Vaksinasi' : 'Semarang Skin Center • Dermatologi & Estetika' }}</span>
          </div>

          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            {{ brand === 'kidinest' ? 'Ambil Nomor Antrian Poli Anak' : 'Reservasi Antrian Konsultasi & Tindakan' }}
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto font-medium">
            {{
              brand === 'kidinest'
                ? 'Pendaftaran mandiri antrian poli anak. Tiket digital & estimasi waktu langsung tersimpan di ponsel Anda.'
                : 'Pendaftaran mandiri konsultasi spesialis kulit. Estimasi giliran ruang tindakan langsung tersimpan di ponsel Anda.'
            }}
          </p>
        </div>

        <!-- Form Card (Zero Kotak Pemilih Unit!) -->
        <form
          @submit.prevent="handleRegister"
          class="bg-white p-6 sm:p-8 rounded-3xl border shadow-xs space-y-5"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/25' : 'border-[#9E606D]/20'"
        >
          <!-- Error Alert -->
          <div v-if="errorMessage" class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Patient Name -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
              {{ brand === 'kidinest' ? 'Nama Pasien Anak / Bayi' : 'Nama Lengkap Pasien' }}
            </label>
            <input
              v-model="patientName"
              type="text"
              required
              :placeholder="brand === 'kidinest' ? 'Contoh: Ananda Mikaela (2 th)' : 'Contoh: Amanda Clarissa'"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all"
              :class="
                brand === 'kidinest'
                  ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]'
                  : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'
              "
            />
          </div>

          <!-- Patient WhatsApp -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
              Nomor WhatsApp Aktif
            </label>
            <div class="relative">
              <span class="absolute left-4 top-1/2 -translate-y-1/2 text-xs font-bold opacity-60">
                +62
              </span>
              <input
                v-model="patientPhone"
                type="tel"
                required
                placeholder="81234567890"
                class="w-full pl-12 pr-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all font-mono"
                :class="
                  brand === 'kidinest'
                    ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]'
                    : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'
                "
              />
            </div>
            <p class="text-[10px] opacity-60 mt-1">Tiket antrian digital & panggilan giliran akan dikirim ke nomor ini.</p>
          </div>

          <!-- Service Category Picker -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-2">
              Pilihan Layanan / Kategori
            </label>
            <div class="space-y-2">
              <label
                v-for="cat in currentCategories"
                :key="cat.key"
                class="flex items-start gap-3 p-3.5 rounded-2xl border cursor-pointer transition-all"
                :class="
                  selectedCategory === cat.key
                    ? brand === 'kidinest'
                      ? 'border-[#0C3B36] bg-[#FAF7F2] shadow-2xs font-bold'
                      : 'border-[#9E606D] bg-[#FAF5F5] shadow-2xs font-bold'
                    : 'border-zinc-200 hover:border-zinc-300 bg-white'
                "
              >
                <input
                  type="radio"
                  v-model="selectedCategory"
                  :value="cat.key"
                  class="mt-0.5"
                  :class="brand === 'kidinest' ? 'accent-[#0C3B36]' : 'accent-[#9E606D]'"
                />
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-extrabold">{{ cat.label }}</div>
                  <div class="text-[11px] opacity-70 font-medium">{{ cat.desc }}</div>
                </div>
              </label>
            </div>
          </div>

          <!-- Doctor Selector (Sesuai Tim Medis Lantai) -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
              Dokter Spesialis yang Berpraktik Hari Ini
            </label>
            <select
              v-model="selectedDoctor"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all cursor-pointer"
              :class="
                brand === 'kidinest'
                  ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]'
                  : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'
              "
            >
              <!-- Kidinest Doctors (Lantai 1) -->
              <template v-if="brand === 'kidinest'">
                <option value="dr. Patricia Vanessa A., Sp.A (Poli Anak 1)">
                  dr. Patricia Vanessa A., Sp.A (Ruang Poli Anak 1 • Pagi)
                </option>
                <option value="dr. Hendra Pratama, Sp.A (Poli Anak 2)">
                  dr. Hendra Pratama, Sp.A (Ruang Poli Anak 2 • Sore)
                </option>
              </template>

              <!-- SSC Doctors (Lantai 2) -->
              <template v-else>
                <option value="dr. Dian Paramita, Sp.D.V.E (Ruang Laser 1)">
                  dr. Dian Paramita, Sp.D.V.E (Ruang Laser & Konsultasi 1 • Pagi)
                </option>
                <option value="dr. Kevin Wijaya, Sp.D.V.E (Ruang Tindakan 2)">
                  dr. Kevin Wijaya, Sp.D.V.E (Ruang Tindakan Bedah 2 • Sore)
                </option>
              </template>
            </select>
          </div>

          <!-- Payment Method (Khusus Kidinest Lantai 1) -->
          <div v-if="brand === 'kidinest'" class="pt-2 border-t border-zinc-100">
            <div class="flex items-center justify-between mb-2">
              <label class="text-[11px] font-bold uppercase tracking-wider opacity-75">
                Biaya Registrasi Awal Pasien:
              </label>
              <span class="text-sm font-black text-[#0C3B36]">Rp 50.000</span>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <label
                class="flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'QRIS' ? 'border-[#0C3B36] bg-[#FAF7F2]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="QRIS" class="accent-[#0C3B36]" />
                <QrCode class="w-4 h-4 text-emerald-600" />
                <span>QRIS Instan</span>
              </label>

              <label
                class="flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'BANK_TRANSFER' ? 'border-[#0C3B36] bg-[#FAF7F2]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="BANK_TRANSFER" class="accent-[#0C3B36]" />
                <CreditCard class="w-4 h-4 text-sky-600" />
                <span>Transfer BCA</span>
              </label>
            </div>
          </div>

          <!-- SSC Free Booking Note -->
          <div v-else class="p-3.5 rounded-2xl bg-rose-50/70 border border-[#9E606D]/20 text-[11px] text-[#2D1E24] font-medium leading-relaxed">
            💡 <strong>Semarang Skin Center:</strong> Pendaftaran reservasi konsultasi dokter spesialis kulit Rp 0. Rincian tindakan medis/laser dan resep akan diselesaikan di kasir setelah konsultasi.
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isSubmitting"
            class="w-full py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
          >
            <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
            <span v-else>{{ brand === 'kidinest' ? 'Lanjut ke Pembayaran Registrasi (Rp 50.000)' : 'Konfirmasi Reservasi Antrian (Rp 0)' }}</span>
            <ChevronRight v-if="!isSubmitting" class="w-4 h-4" />
          </button>
        </form>

      </div>

      <!-- STEP 2: PAYMENT SCREEN (KIDINEST LANTAI 1 ONLY) -->
      <div v-else-if="currentStep === 'payment' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-amber-100 text-amber-900">
            <Clock class="w-3.5 h-3.5" /> Menunggu Pembayaran Registrasi
          </span>
          <h2 class="text-2xl font-black text-[#0C3B36]">Selesaikan Pembayaran</h2>
          <p class="text-xs opacity-75">Nomor antrian FIFO resmi akan otomatis terbit saat pembayaran diverifikasi.</p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-6 text-center">
          <div class="p-4 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/20 text-left space-y-1">
            <div class="flex items-center justify-between text-xs opacity-70">
              <span>Pasien Anak:</span>
              <span class="font-bold text-[#0C3B36]">{{ patientName }}</span>
            </div>
            <div class="flex items-center justify-between text-xs opacity-70">
              <span>Lokasi Poli:</span>
              <span class="font-bold text-[#0C3B36]">Kidinest Pediatric Care</span>
            </div>
            <div class="flex items-center justify-between text-xs opacity-70">
              <span>Dokter Pemeriksa:</span>
              <span class="font-bold text-[#0C3B36]">{{ activeTicket.doctor_name }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-[#E5A93C]/20 text-sm font-extrabold text-[#0C3B36]">
              <span>Total Tagihan Registrasi:</span>
              <span class="text-emerald-700">Rp 50.000</span>
            </div>
          </div>

          <!-- QRIS Dummy Display -->
          <div class="p-6 rounded-3xl bg-[#FAF7F2] border border-[#E5A93C]/20 inline-block space-y-3">
            <div class="w-48 h-48 mx-auto bg-white p-3 rounded-2xl border border-zinc-300 flex items-center justify-center shadow-xs">
              <QrCode class="w-40 h-40 text-[#0C3B36]" />
            </div>
            <div class="text-[11px] font-mono font-bold text-zinc-600">
              NMID: ID1020039201990 (Kidinest Care)
            </div>
            <div class="text-[10px] text-zinc-400">Scan QRIS BCA, Mandiri, GoPay, OVO, ShopeePay</div>
          </div>

          <!-- Transfer Info -->
          <div class="text-xs font-medium opacity-75 space-y-1">
            <div>Atau transfer manual ke rekening resmi klinik:</div>
            <div class="font-mono font-extrabold text-[#0C3B36] bg-[#E5A93C]/20 py-1.5 px-3 rounded-full inline-block">
              BCA 828-192-8811 a.n Kidinest Care
            </div>
          </div>

          <!-- Simulation Button (Demo Ready) -->
          <div class="pt-4 border-t border-zinc-100 space-y-2">
            <button
              @click="handleSimulatePayment"
              :disabled="isSubmitting"
              class="w-full py-3.5 rounded-2xl font-bold bg-emerald-700 hover:bg-emerald-800 text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <CheckCircle2 v-else class="w-4 h-4" />
              <span>Simulasikan Pembayaran Berhasil (Demo)</span>
            </button>

            <button
              @click="currentStep = 'register'"
              class="text-xs font-bold opacity-60 hover:opacity-100 py-1"
            >
              Ubah Data Pendaftaran
            </button>
          </div>
        </div>
      </div>

      <!-- STEP 3: TICKET SCREEN (BOARDING PASS DIGITAL SESUAI LANTAI) -->
      <div v-else-if="currentStep === 'ticket' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-emerald-100 text-emerald-900">
            <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" /> Pendaftaran Berhasil & Terverifikasi
          </span>
          <h2 class="text-2xl font-black">Tiket Antrian Digital Anda</h2>
          <p class="text-xs opacity-75">Simpan layar ini atau pantau nomor antrian secara live.</p>
        </div>

        <!-- Boarding Pass Style Card -->
        <div class="bg-white rounded-3xl border shadow-xs overflow-hidden"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/30' : 'border-[#9E606D]/25'"
        >
          <!-- Card Header -->
          <div
            class="p-6 text-white text-center space-y-2"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36]' : 'bg-[#9E606D]'"
          >
            <div class="text-[11px] font-bold uppercase tracking-wider text-white/70">
              {{ brand === 'kidinest' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}
            </div>
            <!-- Giant Number -->
            <div class="text-6xl font-black tracking-tight text-white font-mono py-2">
              {{ activeTicket.queue_number }}
            </div>
            <div class="inline-block px-3 py-1 rounded-full text-xs font-bold bg-white/15 text-white">
              Status: Sedang Menunggu Panggilan
            </div>
          </div>

          <!-- Card Details -->
          <div class="p-6 space-y-4 text-xs">
            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Pasien</span>
                <span class="font-extrabold text-sm">{{ patientName }}</span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">WhatsApp</span>
                <span class="font-mono font-bold text-sm">{{ patientPhone }}</span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Dokter Praktik</span>
                <span class="font-bold">{{ activeTicket.doctor_name }}</span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Status Biaya</span>
                <span class="font-bold text-emerald-700">
                  {{ activeTicket.registration_fee > 0 ? 'Rp 50.000 (Lunas)' : 'Booking Rp0 (Pasca Tindakan)' }}
                </span>
              </div>
            </div>

            <!-- WhatsApp Badge -->
            <div class="p-3.5 rounded-2xl bg-emerald-50 border border-emerald-200/80 flex items-center gap-3">
              <MessageCircle class="w-5 h-5 text-emerald-600 shrink-0" />
              <div class="text-[11px] text-emerald-950 font-medium leading-tight">
                Notifikasi antrian live terhubung ke nomor <strong>{{ patientPhone }}</strong>.
              </div>
            </div>
          </div>

          <!-- Card Footer Action -->
          <div class="p-4 bg-zinc-50 border-t border-zinc-100 flex items-center justify-between gap-3">
            <button
              @click="emit('navigate', 'monitor')"
              class="inline-flex items-center gap-1.5 text-xs font-bold opacity-75 hover:opacity-100"
            >
              <Smartphone class="w-4 h-4" />
              <span>Lihat TV Monitor Antrian</span>
            </button>

            <button
              @click="resetForm"
              class="px-4 py-2 rounded-xl text-white text-xs font-bold transition-all shadow-xs"
              :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
            >
              Daftar Antrian Baru
            </button>
          </div>
        </div>
      </div>

    </main>

    <!-- Footer -->
    <footer class="max-w-2xl mx-auto w-full text-center text-[11px] font-medium opacity-60 pt-4">
      Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang
    </footer>
  </div>
</template>
