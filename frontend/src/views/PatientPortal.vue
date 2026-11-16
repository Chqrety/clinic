<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import {
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
  ChevronLeft,
  ShieldCheck,
  ExternalLink,
  RefreshCw,
  Search,
  UserCheck,
  UserPlus,
  Calendar as CalendarIcon,
  Stethoscope,
  RotateCcw,
  ReceiptText
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  brand?: 'kidinest' | 'ssc'
}>(), {
  brand: 'kidinest'
})

const emit = defineEmits<{
  (e: 'navigate', view: string): void
}>()

// Step wizard: 1 | 2 | 3 | 4 | 'payment' | 'ticket'
const currentStep = ref<number | 'payment' | 'ticket'>(1)

// Active unit computed strictly from brand context
const selectedUnit = computed<'KIDINEST' | 'SSC'>(() => {
  return props.brand === 'ssc' ? 'SSC' : 'KIDINEST'
})

interface TicketData {
  id?: number
  queue_id?: number
  doctor_id?: number
  queue_number: string
  clinic_unit: string
  category: string
  doctor_name: string
  appointment_date?: string
  appointment_time?: string
  patient_type?: string
  registration_fee: number
  payment_method: string
  payment_status: string
  queue_status: string
  payment_reference?: string
  payment_url?: string
  snap_token?: string
  client_key?: string
  qris_payload?: string
  bank_transfer_account?: string
}

// ==================== STEP 1: TIPE PASIEN ====================
const patientType = ref<'BARU' | 'LAMA'>('BARU')
const patientName = ref('Budi Pratama')
const patientPhone = ref('081234567890')
const lookupQuery = ref('')
const isSearchingPatient = ref(false)
const lookupFeedback = ref<{ success: boolean; message: string } | null>(null)
const foundActiveTicket = ref<TicketData | null>(null)

// Modal Lacak Tiket Mandiri
const showTrackModal = ref(false)
const trackQuery = ref('')
const isTracking = ref(false)
const trackError = ref('')

function openFoundTicket(ticket: TicketData) {
  activeTicket.value = ticket
  currentStep.value = 'ticket'
  showTrackModal.value = false
  try {
    const id = ticket.queue_id || ticket.id
    if (id) {
      localStorage.setItem('clinic_last_ticket_id', String(id))
    }
  } catch {}
}

async function handleTrackTicket() {
  const query = trackQuery.value.trim()
  if (!query) return
  isTracking.value = true
  trackError.value = ''
  try {
    const res = await fetch(`/api/public/patients/lookup?query=${encodeURIComponent(query)}`)
    const data = await res.json()
    if (res.ok && data.success && data.data && data.data.active_ticket) {
      openFoundTicket(data.data.active_ticket)
    } else if (res.ok && data.success && data.data && !data.data.active_ticket) {
      trackError.value = `Pasien terdaftar (${data.data.patient?.name || ''}), namun tidak memiliki tiket antrian aktif saat ini.`
    } else {
      trackError.value = 'Tiket atau data pasien tidak ditemukan dengan nomor tersebut.'
    }
  } catch {
    trackError.value = 'Gagal memeriksa tiket. Periksa koneksi Anda.'
  } finally {
    isTracking.value = false
  }
}

async function handleLookupPatient() {
  const query = lookupQuery.value.trim()
  if (!query) return

  isSearchingPatient.value = true
  lookupFeedback.value = null
  foundActiveTicket.value = null

  try {
    const res = await fetch(`/api/public/patients/lookup?query=${encodeURIComponent(query)}`)
    const data = await res.json()
    if (res.ok && data.success && data.data) {
      const patient = data.data.patient || data.data
      patientName.value = patient.name
      patientPhone.value = patient.phone
      patientType.value = 'LAMA'

      if (data.data.active_ticket) {
        foundActiveTicket.value = data.data.active_ticket
        lookupFeedback.value = {
          success: true,
          message: `Pasien ${patient.name} ditemukan! Memiliki tiket aktif.`
        }
      } else {
        foundActiveTicket.value = null
        lookupFeedback.value = {
          success: true,
          message: `Pasien terdaftar: ${patient.name} (${patient.phone})`
        }
      }
    } else {
      lookupFeedback.value = {
        success: false,
        message: 'Data pasien tidak ditemukan. Silakan isi form sebagai Pasien Baru.'
      }
    }
  } catch {
    lookupFeedback.value = {
      success: false,
      message: 'Gagal mencari data pasien. Periksa koneksi Anda.'
    }
  } finally {
    isSearchingPatient.value = false
  }
}

// ==================== STEP 2: KATEGORI LAYANAN ====================
const kidinestCategories = [
  { key: 'A_VAKSIN', label: 'A • Vaksinasi & Imunisasi Anak', desc: 'Imunisasi dasar & lanjutan rantai dingin IDAI' },
  { key: 'B_SAKIT', label: 'B • Konsultasi Pediatri & Anak Sakit', desc: 'Demam, batuk pilek, flu, diare, ruam infeksi' },
  { key: 'C_TUMBANG', label: 'C • Skrining Milestone & Nutrisi', desc: 'Pemantauan grafik tumbuh kembang & gizi anak' },
  { key: 'D_PAKET', label: 'D • Paket Vaksin + Konsultasi Tumbang', desc: 'Pemeriksaan pediatri menyeluruh' },
]

const sscCategories = [
  { key: 'SSC_FACIAL', label: 'Medical Aesthetic Facial & Infusion', desc: 'Pembersihan komedo steril & serum asam buah medis' },
  { key: 'SSC_LASER', label: 'Laser Rejuvenation & Pori Kulit', desc: 'Melasma, peremajaan kolagen kulit & tekstur halus' },
  { key: 'SSC_ACNE', label: 'Acne & Atrophic Scar Reconstruction', desc: 'Penanganan jerawat aktif & perbaikan bopeng bekas jerawat' },
  { key: 'SSC_SURGERY', label: 'Bedah Minor Kulit & Kauterisasi', desc: 'Biopsi lesi, kauter keratosis & pengangkatan kista milia' },
]

const currentCategories = computed(() => {
  return selectedUnit.value === 'KIDINEST' ? kidinestCategories : sscCategories
})

const selectedCategory = ref(props.brand === 'ssc' ? 'SSC_FACIAL' : 'A_VAKSIN')

// ==================== STEP 3: PILIH DOKTER & JADWAL (30-60 HARI) ====================
interface DoctorItem {
  id: number
  name: string
  role: string
  clinic_unit: string
  room_name?: string
}

interface AvailableSlot {
  time_slot: string
  start_time: string
  end_time: string
  total_quota: number
  booked: number
  available: number
  is_available: boolean
}

interface AvailableDate {
  date: string
  day_name: string
  day_of_week: number
  shift_name?: string
  is_available: boolean
  slots: AvailableSlot[]
}

const doctorList = ref<DoctorItem[]>([])
const selectedDoctorId = ref<number | null>(null)
const selectedDoctorName = ref('')
const isLoadingDoctors = ref(false)

const availableDates = ref<AvailableDate[]>([])
const selectedDateStr = ref<string>('')
const selectedTimeSlot = ref<string>('')
const isLoadingSchedules = ref(false)

async function fetchDoctors() {
  isLoadingDoctors.value = true
  try {
    const res = await fetch(`/api/public/doctors?unit=${selectedUnit.value}`)
    const data = await res.json()
    if (res.ok && data.success && Array.isArray(data.data) && data.data.length > 0) {
      doctorList.value = data.data
      if (!selectedDoctorId.value || !doctorList.value.some(d => d.id === selectedDoctorId.value)) {
        selectedDoctorId.value = doctorList.value[0].id
        selectedDoctorName.value = doctorList.value[0].name
      }
    } else {
      // Fallback data demo
      if (selectedUnit.value === 'KIDINEST') {
        doctorList.value = [
          { id: 1, name: 'dr. Patricia Vanessa A., Sp.A', role: 'DOCTOR', clinic_unit: 'KIDINEST', room_name: 'Poli Anak 1' },
          { id: 2, name: 'dr. Hendra Pratama, Sp.A', role: 'DOCTOR', clinic_unit: 'KIDINEST', room_name: 'Poli Anak 2' }
        ]
      } else {
        doctorList.value = [
          { id: 3, name: 'dr. Dian Paramita, Sp.D.V.E', role: 'DOCTOR', clinic_unit: 'SSC', room_name: 'Poli Estetika 1' },
          { id: 4, name: 'dr. Kevin Wijaya, Sp.D.V.E', role: 'DOCTOR', clinic_unit: 'SSC', room_name: 'Poli Estetika 2' }
        ]
      }
      selectedDoctorId.value = doctorList.value[0].id
      selectedDoctorName.value = doctorList.value[0].name
    }
  } catch {
    // Fallback offline
    selectedDoctorName.value = selectedUnit.value === 'KIDINEST' ? 'dr. Patricia Vanessa A., Sp.A' : 'dr. Dian Paramita, Sp.D.V.E'
  } finally {
    isLoadingDoctors.value = false
    if (selectedDoctorId.value) {
      fetchDoctorSchedules(selectedDoctorId.value)
    }
  }
}

async function fetchDoctorSchedules(doctorId: number) {
  isLoadingSchedules.value = true
  availableDates.value = []
  selectedDateStr.value = ''
  selectedTimeSlot.value = ''

  try {
    const res = await fetch(`/api/public/doctors/${doctorId}/schedules?days=45`)
    const data = await res.json()
    if (res.ok && data.success && Array.isArray(data.data) && data.data.length > 0) {
      availableDates.value = data.data
      // Pilih tanggal pertama yang tersedia
      const firstAvail = availableDates.value.find(d => d.is_available && d.slots.some(s => s.is_available))
      if (firstAvail) {
        selectedDateStr.value = firstAvail.date
        const firstSlot = firstAvail.slots.find(s => s.is_available)
        if (firstSlot) {
          selectedTimeSlot.value = firstSlot.start_time
        }
      }
    } else {
      generateMockSchedules()
    }
  } catch {
    generateMockSchedules()
  } finally {
    isLoadingSchedules.value = false
  }
}

function generateMockSchedules() {
  const list: AvailableDate[] = []
  const today = new Date()
  const dayNames = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']

  for (let i = 0; i < 30; i++) {
    const d = new Date()
    d.setDate(today.getDate() + i)
    const dateStr = d.toISOString().split('T')[0]
    const dow = d.getDay()
    const isAvail = dow !== 0 // Minggu libur

    const slots: AvailableSlot[] = []
    if (isAvail) {
      const times = ['08:30', '09:00', '09:30', '10:00', '10:30', '11:00', '11:30']
      times.forEach(t => {
        slots.push({
          time_slot: `${t} - ...`,
          start_time: t,
          end_time: '',
          total_quota: 1,
          booked: 0,
          available: 1,
          is_available: true
        })
      })
    }

    list.push({
      date: dateStr,
      day_name: dayNames[dow],
      day_of_week: dow,
      shift_name: 'PAGI',
      is_available: isAvail,
      slots: slots
    })
  }
  availableDates.value = list
  const firstAvail = list.find(d => d.is_available)
  if (firstAvail) {
    selectedDateStr.value = firstAvail.date
    if (firstAvail.slots[0]) selectedTimeSlot.value = firstAvail.slots[0].start_time
  }
}

function selectDoctor(doc: DoctorItem) {
  selectedDoctorId.value = doc.id
  selectedDoctorName.value = doc.name
  fetchDoctorSchedules(doc.id)
}

const currentSelectedDateSlots = computed(() => {
  const d = availableDates.value.find(item => item.date === selectedDateStr.value)
  return d ? d.slots : []
})

// ==================== STEP 4: KONFIRMASI & PEMBAYARAN ====================
const paymentMethod = ref<'QRIS' | 'BANK_TRANSFER' | 'CASH'>('QRIS')
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

const activeTicket = ref<TicketData | null>(null)

const resolvedQRImage = computed(() => {
  const url = activeTicket.value?.payment_url
  if (url && (url.includes('/qr-code') || url.endsWith('.png') || url.endsWith('.jpg') || url.includes('qrserver') || url.includes('tripay'))) {
    return url
  }
  if (activeTicket.value?.qris_payload) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(activeTicket.value.qris_payload)}`
  }
  return null
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
        patient_type: patientType.value,
        category: selectedCategory.value,
        doctor_id: selectedDoctorId.value,
        doctor_name: selectedDoctorName.value,
        appointment_date: selectedDateStr.value || undefined,
        appointment_time: selectedTimeSlot.value || undefined,
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
      startPaymentPolling(data.data.queue_id)
    } else {
      currentStep.value = 'ticket'
      try {
        localStorage.setItem('clinic_last_ticket_id', String(data.data.queue_id || data.data.id))
      } catch {}
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Terjadi kesalahan sistem'
  } finally {
    isSubmitting.value = false
  }
}

// Polling status pembayaran
let pollTimer: any = null
function startPaymentPolling(queueId: number) {
  stopPaymentPolling()
  pollTimer = setInterval(async () => {
    try {
      const res = await fetch(`/api/public/queue/ticket/${queueId}`)
      const data = await res.json()
      if (res.ok && data.success && data.data) {
        if (data.data.payment_status === 'PAID') {
          activeTicket.value = data.data
          stopPaymentPolling()
          currentStep.value = 'ticket'
          try {
            localStorage.setItem('clinic_last_ticket_id', String(queueId))
          } catch {}
        }
      }
    } catch {
      // polling error silent
    }
  }, 3000)
}

function stopPaymentPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

const isCheckingPayment = ref(false)
const checkPaymentMessage = ref<string | null>(null)

async function handleManualCheckPayment() {
  if (!activeTicket.value || isCheckingPayment.value) return
  isCheckingPayment.value = true
  checkPaymentMessage.value = null

  try {
    const res = await fetch(`/api/public/queue/ticket/${activeTicket.value.queue_id}`)
    const data = await res.json()
    if (res.ok && data.success && data.data) {
      if (data.data.payment_status === 'PAID') {
        activeTicket.value = data.data
        stopPaymentPolling()
        currentStep.value = 'ticket'
        return
      }
    }
    checkPaymentMessage.value = 'Pembayaran belum terverifikasi oleh Midtrans. Pastikan pembayaran di simulator/e-wallet sudah berhasil.'
    setTimeout(() => { checkPaymentMessage.value = null }, 4000)
  } catch {
    checkPaymentMessage.value = 'Gagal memeriksa status pembayaran. Coba beberapa saat lagi.'
    setTimeout(() => { checkPaymentMessage.value = null }, 4000)
  } finally {
    isCheckingPayment.value = false
  }
}

async function handleSimulatePayment() {
  stopPaymentPolling()
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
      activeTicket.value.queue_number = activeTicket.value.appointment_date ? `BK-KID-${activeTicket.value.queue_id}` : 'A-001'
      activeTicket.value.payment_status = 'PAID'
    }
  } catch {
    if (activeTicket.value) {
      activeTicket.value.queue_number = activeTicket.value.appointment_date ? `BK-KID-${activeTicket.value.queue_id}` : 'A-001'
      activeTicket.value.payment_status = 'PAID'
    }
  } finally {
    isSubmitting.value = false
    currentStep.value = 'ticket'
  }
}

// ==================== RESCHEDULE & REFUND MODALS ====================
const showRescheduleModal = ref(false)
const reschedAvailableDates = ref<AvailableDate[]>([])
const reschedSelectedDateStr = ref<string>('')
const reschedSelectedTimeSlot = ref<string>('')
const isLoadingReschedSchedules = ref(false)
const isSubmittingResched = ref(false)

const currentReschedSlots = computed(() => {
  const d = reschedAvailableDates.value.find(item => item.date === reschedSelectedDateStr.value)
  return d ? d.slots : []
})

async function openRescheduleModal() {
  if (!activeTicket.value) return
  showRescheduleModal.value = true
  isLoadingReschedSchedules.value = true
  reschedAvailableDates.value = []
  reschedSelectedDateStr.value = ''
  reschedSelectedTimeSlot.value = ''

  // Ambil ID dokter dari tiket, fallback ke selectedDoctorId atau 1
  const doctorId = activeTicket.value.doctor_id || selectedDoctorId.value || 1
  try {
    const res = await fetch(`/api/public/doctors/${doctorId}/schedules?days=45`)
    const data = await res.json()
    if (res.ok && data.success && Array.isArray(data.data) && data.data.length > 0) {
      reschedAvailableDates.value = data.data
      const firstAvail = reschedAvailableDates.value.find(d => d.is_available && d.slots.some(s => s.is_available))
      if (firstAvail) {
        reschedSelectedDateStr.value = firstAvail.date
        const firstSlot = firstAvail.slots.find(s => s.is_available)
        if (firstSlot) {
          reschedSelectedTimeSlot.value = firstSlot.start_time
        }
      }
    }
  } catch (err) {
    console.error('Gagal memuat jadwal dokter untuk reschedule:', err)
  } finally {
    isLoadingReschedSchedules.value = false
  }
}

async function handleConfirmReschedule() {
  if (!activeTicket.value || !reschedSelectedDateStr.value || !reschedSelectedTimeSlot.value) return
  const queueId = activeTicket.value.queue_id || activeTicket.value.id
  if (!queueId) return

  isSubmittingResched.value = true
  try {
    const res = await fetch('/api/public/queue/reschedule', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        queue_id: queueId,
        new_date: reschedSelectedDateStr.value,
        new_time: reschedSelectedTimeSlot.value
      })
    })
    const data = await res.json()
    if (res.ok && data.success) {
      activeTicket.value = data.data
      showRescheduleModal.value = false
      alert(`Jadwal kunjungan berhasil diubah ke ${reschedSelectedDateStr.value} jam ${reschedSelectedTimeSlot.value}!`)
    } else {
      alert(data.error || 'Gagal mengubah jadwal')
    }
  } catch {
    alert('Koneksi bermasalah')
  } finally {
    isSubmittingResched.value = false
  }
}

const showRefundModal = ref(false)
const refundBank = ref('BCA')
const refundAccNum = ref('')
const refundAccName = ref('')
const refundNotes = ref('')
const isSubmittingRefund = ref(false)

async function handleConfirmRefund() {
  if (!activeTicket.value || !refundAccNum.value || !refundAccName.value) return
  isSubmittingRefund.value = true
  try {
    const res = await fetch('/api/public/queue/refund/request', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        queue_id: activeTicket.value.queue_id,
        bank: refundBank.value,
        account_number: refundAccNum.value,
        account_name: refundAccName.value,
        notes: refundNotes.value
      })
    })
    const data = await res.json()
    if (res.ok && data.success) {
      activeTicket.value = data.data
      showRefundModal.value = false
      alert('Pengajuan refund berhasil dicatat! Kasir kami akan mentransfer manual dan mengonfirmasi via WhatsApp.')
    } else {
      alert(data.error || 'Gagal mengajukan refund')
    }
  } catch {
    alert('Koneksi bermasalah')
  } finally {
    isSubmittingRefund.value = false
  }
}

function resetForm() {
  stopPaymentPolling()
  activeTicket.value = null
  currentStep.value = 1
}

watch(() => props.brand, () => {
  currentStep.value = 1
  fetchDoctors()
})

onMounted(async () => {
  fetchDoctors()

  // Auto-restore tiket dari URL (?ticket=...) atau localStorage jika ada
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const ticketParam = urlParams.get('ticket')
    const savedTicketId = ticketParam || localStorage.getItem('clinic_last_ticket_id')
    if (savedTicketId) {
      const res = await fetch(`/api/public/queue/ticket/${savedTicketId}`)
      const data = await res.json()
      if (res.ok && data.success && data.data) {
        if (ticketParam || ['BOOKED', 'WAITING', 'CALLING'].includes(data.data.queue_status)) {
          activeTicket.value = data.data
          currentStep.value = 'ticket'
        }
      }
    }
  } catch (err) {
    console.warn('Gagal memulihkan status tiket:', err)
  }
})

onUnmounted(() => {
  stopPaymentPolling()
})
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
        <button
          type="button"
          @click="showTrackModal = true"
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold border transition-colors cursor-pointer"
          :class="brand === 'kidinest' ? 'bg-[#0C3B36]/5 hover:bg-[#0C3B36]/10 text-[#0C3B36] border-[#0C3B36]/20' : 'bg-[#9E606D]/5 hover:bg-[#9E606D]/10 text-[#9E606D] border-[#9E606D]/20'"
        >
          <Search class="w-3 h-3" />
          <span>Lacak Tiket</span>
        </button>
        <span
          class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold border shadow-2xs"
          :class="brand === 'kidinest' ? 'bg-white border-[#E5A93C]/30 text-[#0C3B36]' : 'bg-white border-[#9E606D]/30 text-[#9E606D]'"
        >
          <ShieldCheck class="w-3 h-3 text-emerald-600" />
          <span>Faskes Resmi</span>
        </span>
      </div>
    </header>

    <!-- Main Container -->
    <main class="max-w-2xl mx-auto w-full py-2 flex-1 flex flex-col justify-center">

      <!-- STEP PROGRESS INDICATOR (1 -> 2 -> 3 -> 4) -->
      <div v-if="typeof currentStep === 'number'" class="mb-6 space-y-2">
        <div class="flex items-center justify-between px-2 text-[11px] font-extrabold uppercase tracking-wider">
          <span :class="currentStep >= 1 ? (brand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#9E606D]') : 'opacity-40'">
            1. Pasien
          </span>
          <span :class="currentStep >= 2 ? (brand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#9E606D]') : 'opacity-40'">
            2. Layanan
          </span>
          <span :class="currentStep >= 3 ? (brand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#9E606D]') : 'opacity-40'">
            3. Jadwal Dokter
          </span>
          <span :class="currentStep >= 4 ? (brand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#9E606D]') : 'opacity-40'">
            4. Konfirmasi
          </span>
        </div>
        <div class="w-full bg-zinc-200/80 h-1.5 rounded-full overflow-hidden">
          <div
            class="h-full transition-all duration-300 rounded-full"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36]' : 'bg-[#9E606D]'"
            :style="{ width: `${(currentStep / 4) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- ==================== STEP 1: PASIEN LAMA VS BARU ==================== -->
      <div v-if="currentStep === 1" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <span>Tahap 1 • Identitas Pasien</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            {{ brand === 'kidinest' ? 'Pendaftaran Pasien Anak' : 'Pendaftaran Konsultasi Estetika' }}
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Apakah pasien sudah pernah terdaftar di klinik sebelumnya?
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border shadow-xs space-y-5"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/25' : 'border-[#9E606D]/20'"
        >
          <!-- Switcher: Pasien Baru vs Pasien Lama -->
          <div class="grid grid-cols-2 gap-2 p-1.5 bg-zinc-100/80 rounded-2xl">
            <button
              type="button"
              @click="patientType = 'BARU'"
              class="py-2.5 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-2"
              :class="patientType === 'BARU' ? 'bg-white shadow-2xs font-extrabold text-[#0C3B36]' : 'text-zinc-600 hover:text-zinc-900'"
            >
              <UserPlus class="w-4 h-4" />
              <span>Pasien Baru</span>
            </button>
            <button
              type="button"
              @click="patientType = 'LAMA'"
              class="py-2.5 rounded-xl text-xs font-bold transition-all flex items-center justify-center gap-2"
              :class="patientType === 'LAMA' ? 'bg-white shadow-2xs font-extrabold text-[#0C3B36]' : 'text-zinc-600 hover:text-zinc-900'"
            >
              <UserCheck class="w-4 h-4" />
              <span>Pasien Lama</span>
            </button>
          </div>

          <!-- Bagian Pencarian Pasien Lama -->
          <div v-if="patientType === 'LAMA'" class="p-4 rounded-2xl bg-amber-50/60 border border-amber-200/60 space-y-3">
            <label class="block text-[11px] font-bold uppercase tracking-wider text-amber-900">
              Cari Berdasarkan No. HP WhatsApp atau Nama:
            </label>
            <div class="flex gap-2">
              <input
                v-model="lookupQuery"
                type="text"
                placeholder="Contoh: 081234567890 atau Budi"
                @keyup.enter="handleLookupPatient"
                class="flex-1 px-4 py-2.5 rounded-xl border border-amber-300/80 bg-white text-xs font-bold outline-none"
              />
              <button
                type="button"
                @click="handleLookupPatient"
                :disabled="isSearchingPatient"
                class="px-4 py-2.5 rounded-xl bg-[#0C3B36] text-white text-xs font-bold flex items-center gap-1.5 hover:bg-[#134E48] transition-all disabled:opacity-50"
              >
                <Loader2 v-if="isSearchingPatient" class="w-3.5 h-3.5 animate-spin" />
                <Search v-else class="w-3.5 h-3.5" />
                <span>Cari</span>
              </button>
            </div>

            <div v-if="lookupFeedback" class="text-xs font-semibold flex items-center gap-2"
              :class="lookupFeedback.success ? 'text-emerald-700' : 'text-rose-700'"
            >
              <CheckCircle2 v-if="lookupFeedback.success" class="w-4 h-4 shrink-0 text-emerald-600" />
              <AlertCircle v-else class="w-4 h-4 shrink-0 text-rose-600" />
              <span>{{ lookupFeedback.message }}</span>
            </div>

            <!-- Kartu Tiket Aktif Ditemukan -->
            <div
              v-if="foundActiveTicket"
              class="p-3.5 rounded-xl bg-white border border-emerald-300 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-3 animate-in fade-in"
            >
              <div class="space-y-0.5">
                <div class="flex items-center gap-2">
                  <span class="text-xs font-extrabold text-emerald-900">
                    Tiket Aktif: {{ foundActiveTicket.queue_number }}
                  </span>
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-emerald-100 text-emerald-800">
                    {{ foundActiveTicket.queue_status }}
                  </span>
                </div>
                <div class="text-[11px] text-zinc-600">
                  Jadwal: <span class="font-bold text-zinc-800">{{ foundActiveTicket.appointment_date || 'Hari Ini' }} • {{ foundActiveTicket.appointment_time || '-' }}</span> ({{ foundActiveTicket.doctor_name || 'Dokter Jaga' }})
                </div>
              </div>
              <button
                type="button"
                @click="openFoundTicket(foundActiveTicket)"
                class="px-4 py-2 rounded-xl bg-emerald-700 hover:bg-emerald-800 text-white text-xs font-bold shadow-sm transition-all flex items-center justify-center gap-1.5 shrink-0 cursor-pointer"
              >
                <span>Buka Tiket Saya</span>
                <ChevronRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Form Input Pasien -->
          <div class="space-y-4">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
                {{ brand === 'kidinest' ? 'Nama Pasien Anak / Bayi' : 'Nama Lengkap Pasien' }}
              </label>
              <input
                v-model="patientName"
                type="text"
                required
                :placeholder="brand === 'kidinest' ? 'Contoh: Ananda Mikaela' : 'Contoh: Amanda Clarissa'"
                class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all"
                :class="brand === 'kidinest' ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]' : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'"
              />
            </div>

            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
                Nomor WhatsApp Aktif
              </label>
              <input
                v-model="patientPhone"
                type="tel"
                required
                placeholder="Contoh: 081234567890"
                class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all font-mono"
                :class="brand === 'kidinest' ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]' : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'"
              />
              <p class="text-[10px] opacity-60 mt-1">Tiket digital, kode booking & panggilan giliran otomatis dikirim via WhatsApp.</p>
            </div>
          </div>

          <button
            type="button"
            @click="currentStep = 2"
            :disabled="!patientName.trim() || !patientPhone.trim()"
            class="w-full py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-40"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
          >
            <span>Lanjut: Pilih Kategori Layanan</span>
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- ==================== STEP 2: PILIH KATEGORI LAYANAN ==================== -->
      <div v-else-if="currentStep === 2" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <span>Tahap 2 • Kategori Layanan</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            {{ brand === 'kidinest' ? 'Pilih Kategori Poli Anak' : 'Pilih Perawatan Estetika' }}
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Tentukan jenis penanganan yang dibutuhkan pasien untuk penjadwalan tepat sasaran.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border shadow-xs space-y-4"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/25' : 'border-[#9E606D]/20'"
        >
          <div class="space-y-2.5">
            <label
              v-for="cat in currentCategories"
              :key="cat.key"
              class="flex items-start gap-3 p-4 rounded-2xl border cursor-pointer transition-all"
              :class="selectedCategory === cat.key
                ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#FAF7F2] shadow-2xs font-bold' : 'border-[#9E606D] bg-[#FAF5F5] shadow-2xs font-bold')
                : 'border-zinc-200 hover:border-zinc-300 bg-white'"
            >
              <input
                type="radio"
                v-model="selectedCategory"
                :value="cat.key"
                class="mt-1"
                :class="brand === 'kidinest' ? 'accent-[#0C3B36]' : 'accent-[#9E606D]'"
              />
              <div class="min-w-0 flex-1">
                <div class="text-xs font-extrabold">{{ cat.label }}</div>
                <div class="text-[11px] opacity-70 font-medium mt-0.5">{{ cat.desc }}</div>
              </div>
            </label>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 1"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="currentStep = 3"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2"
              :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
            >
              <span>Lanjut: Pilih Dokter & Jadwal</span>
              <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== STEP 3: PILIH DOKTER & KALENDER JADWAL (30-60 HARI) ==================== -->
      <div v-else-if="currentStep === 3" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <span>Tahap 3 • Dokter & Jadwal Kunjungan</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Pilih Dokter & Slot Waktu
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Reservasi tersedia hingga 45–60 hari ke depan sesuai shift dokter spesialis.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border shadow-xs space-y-6"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/25' : 'border-[#9E606D]/20'"
        >
          <!-- 1. Pilihan Dokter -->
          <div class="space-y-2">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              1. Pilih Dokter Spesialis:
            </label>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <div
                v-for="doc in doctorList"
                :key="doc.id"
                @click="selectDoctor(doc)"
                class="p-3.5 rounded-2xl border cursor-pointer transition-all flex items-center gap-3"
                :class="selectedDoctorId === doc.id
                  ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#FAF7F2] shadow-2xs ring-1 ring-[#0C3B36]' : 'border-[#9E606D] bg-[#FAF5F5] shadow-2xs ring-1 ring-[#9E606D]')
                  : 'border-zinc-200 hover:border-zinc-300 bg-white'"
              >
                <div class="w-9 h-9 rounded-full bg-zinc-100 flex items-center justify-center shrink-0">
                  <Stethoscope class="w-5 h-5 opacity-70" />
                </div>
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-extrabold truncate">{{ doc.name }}</div>
                  <div class="text-[10px] opacity-60 font-semibold">{{ doc.room_name || 'Praktik Spesialis' }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 2. Kalender Horizontal (Tanggal Kunjungan) -->
          <div class="space-y-2 pt-2 border-t border-zinc-100">
            <div class="flex items-center justify-between">
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
                2. Pilih Tanggal Kunjungan:
              </label>
              <span v-if="isLoadingSchedules" class="text-[11px] font-bold text-amber-700 flex items-center gap-1">
                <Loader2 class="w-3 h-3 animate-spin" /> Memuat jadwal...
              </span>
            </div>

            <!-- Horizontal Scroll Date Picker -->
            <div class="flex gap-2 overflow-x-auto pb-2 pt-1 no-scrollbar scroll-smooth">
              <button
                v-for="item in availableDates"
                :key="item.date"
                type="button"
                :disabled="!item.is_available"
                @click="selectedDateStr = item.date"
                class="shrink-0 w-24 p-2.5 rounded-2xl border text-center transition-all disabled:opacity-40 disabled:cursor-not-allowed"
                :class="selectedDateStr === item.date
                  ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#0C3B36] text-white shadow-2xs' : 'border-[#9E606D] bg-[#9E606D] text-white shadow-2xs')
                  : (item.is_available ? 'border-zinc-200 bg-[#FAF7F2] hover:border-zinc-300 text-zinc-800' : 'border-zinc-100 bg-zinc-50 text-zinc-400')"
              >
                <div class="text-[10px] font-bold uppercase tracking-wider">{{ item.day_name }}</div>
                <div class="text-sm font-black my-0.5">{{ item.date.slice(8) }} {{ item.date.slice(5, 7) === '11' ? 'Nov' : 'Des' }}</div>
                <div class="text-[9px] font-medium opacity-80">
                  {{ item.is_available ? (item.shift_name || 'Buka') : 'Libur' }}
                </div>
              </button>
            </div>
          </div>

          <!-- 3. Slot Waktu Jam Praktik -->
          <div class="space-y-2 pt-2 border-t border-zinc-100">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              3. Pilih Jam Slot Kunjungan:
            </label>

            <div v-if="currentSelectedDateSlots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 gap-2">
              <button
                v-for="slot in currentSelectedDateSlots"
                :key="slot.start_time"
                type="button"
                :disabled="!slot.is_available"
                @click="selectedTimeSlot = slot.start_time"
                class="p-2.5 rounded-xl border text-center transition-all disabled:opacity-40 disabled:cursor-not-allowed"
                :class="selectedTimeSlot === slot.start_time
                  ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#0C3B36] text-white font-extrabold shadow-2xs' : 'border-[#9E606D] bg-[#9E606D] text-white font-extrabold shadow-2xs')
                  : (slot.is_available ? 'border-zinc-200 bg-white hover:border-zinc-300 text-zinc-800' : 'border-zinc-100 bg-zinc-50 text-zinc-400')"
              >
                <div class="text-xs font-mono font-bold">{{ slot.start_time }} WIB</div>
                <div class="text-[10px] opacity-75 mt-0.5">
                  {{ slot.is_available ? 'Slot Tersedia' : 'Penuh' }}
                </div>
              </button>
            </div>
            <div v-else class="p-4 rounded-xl bg-zinc-50 text-center text-xs opacity-60">
              Pilih tanggal praktik yang tersedia terlebih dahulu untuk melihat slot jam.
            </div>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 2"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="currentStep = 4"
              :disabled="!selectedDoctorName || !selectedDateStr || !selectedTimeSlot"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-40"
              :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
            >
              <span>Lanjut: Review Konfirmasi</span>
              <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== STEP 4: KONFIRMASI RINGKASAN & BAYAR ==================== -->
      <div v-else-if="currentStep === 4" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <span>Tahap 4 • Review & Konfirmasi</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Konfirmasi Pendaftaran
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Periksa kembali detail pasien, jadwal dokter, dan metode pembayaran sebelum finalisasi.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border shadow-xs space-y-5"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/25' : 'border-[#9E606D]/20'"
        >
          <!-- Alert Error jika ada -->
          <div v-if="errorMessage" class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Ringkasan Tiket Nordic Modern -->
          <div class="p-4 rounded-2xl bg-[#FAF7F2] border border-zinc-200/80 space-y-3 text-xs">
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Unit Faskes:</span>
              <span class="font-extrabold text-[#0C3B36]">{{ brand === 'kidinest' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Pasien:</span>
              <span class="font-bold">{{ patientName }} ({{ patientType }})</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Nomor WhatsApp:</span>
              <span class="font-mono font-bold">{{ patientPhone }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Dokter Spesialis:</span>
              <span class="font-bold">{{ selectedDoctorName }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Jadwal Praktik:</span>
              <span class="font-bold text-emerald-800">{{ selectedDateStr }} • Jam {{ selectedTimeSlot }} WIB</span>
            </div>
            <div class="flex items-center justify-between pt-1 text-sm font-black">
              <span>Biaya Registrasi Awal:</span>
              <span :class="brand === 'kidinest' ? 'text-emerald-700' : 'text-zinc-700'">
                {{ brand === 'kidinest' ? 'Rp 50.000 (QRIS)' : 'Rp 0 (Bayar di Kasir Pasca Tindakan)' }}
              </span>
            </div>
          </div>

          <!-- Pilihan Pembayaran Kidinest -->
          <div v-if="brand === 'kidinest'" class="space-y-2">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              Metode Pembayaran Resmi:
            </label>
            <div class="grid grid-cols-2 gap-2">
              <label class="flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'QRIS' ? 'border-[#0C3B36] bg-[#FAF7F2]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="QRIS" class="accent-[#0C3B36]" />
                <QrCode class="w-4 h-4 text-emerald-600" />
                <span>QRIS Midtrans</span>
              </label>
              <label class="flex items-center gap-2 p-3 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'BANK_TRANSFER' ? 'border-[#0C3B36] bg-[#FAF7F2]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="BANK_TRANSFER" class="accent-[#0C3B36]" />
                <CreditCard class="w-4 h-4 text-sky-600" />
                <span>Transfer BCA</span>
              </label>
            </div>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 3"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="handleRegister"
              :disabled="isSubmitting"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
              :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <span v-else>{{ brand === 'kidinest' ? 'Konfirmasi & Bayar Rp 50.000' : 'Konfirmasi Reservasi Jadwal (Rp 0)' }}</span>
              <ChevronRight v-if="!isSubmitting" class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== PAYMENT SCREEN (MIDTRANS QRIS REAL) ==================== -->
      <div v-else-if="currentStep === 'payment' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-amber-100 text-amber-900">
            <Clock class="w-3.5 h-3.5" /> Menunggu Pembayaran Registrasi
          </span>
          <h2 class="text-2xl font-black text-[#0C3B36]">Selesaikan Pembayaran QRIS</h2>
          <p class="text-xs opacity-75">Nomor antrian resmi / kode booking akan aktif saat pembayaran diverifikasi.</p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-6 text-center">
          <div class="p-4 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/20 text-left space-y-1 text-xs">
            <div class="flex items-center justify-between opacity-70">
              <span>Pasien Anak:</span>
              <span class="font-bold text-[#0C3B36]">{{ patientName }}</span>
            </div>
            <div class="flex items-center justify-between opacity-70">
              <span>Dokter & Tanggal:</span>
              <span class="font-bold text-[#0C3B36]">{{ activeTicket.doctor_name }} • {{ activeTicket.appointment_date || selectedDateStr }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-[#E5A93C]/20 text-sm font-extrabold text-[#0C3B36]">
              <span>Total Tagihan:</span>
              <span class="text-emerald-700">Rp 50.000</span>
            </div>
          </div>

          <!-- GAMBAR QRIS RESMI MIDTRANS -->
          <div class="p-6 rounded-3xl bg-[#FAF7F2] border border-[#E5A93C]/20 inline-block space-y-3 max-w-sm mx-auto">
            <div class="w-60 h-60 mx-auto bg-white p-3 rounded-2xl border border-zinc-200 flex items-center justify-center shadow-xs overflow-hidden">
              <img
                v-if="resolvedQRImage"
                :src="resolvedQRImage"
                alt="QRIS Pembayaran Resmi Midtrans"
                class="w-56 h-56 object-contain"
              />
              <QrCode v-else class="w-44 h-44 text-[#0C3B36]" />
            </div>

            <div class="text-[11px] font-mono font-bold text-zinc-600">
              {{ activeTicket.payment_reference ? `Ref: ${activeTicket.payment_reference}` : 'NMID: ID1020039201990 (Kidinest Care)' }}
            </div>
            <div class="text-[10px] text-zinc-400">Scan QRIS BCA, Mandiri, GoPay, OVO, ShopeePay & Semua Bank</div>
            <div class="inline-flex items-center gap-1.5 text-[10px] text-emerald-700 font-semibold bg-emerald-50 px-2.5 py-1 rounded-full border border-emerald-200/50">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping"></span>
              Menunggu notifikasi pembayaran otomatis
            </div>

            <div v-if="resolvedQRImage" class="pt-1">
              <a
                href="https://simulator.sandbox.midtrans.com/openapi/qris/index"
                target="_blank"
                rel="noopener noreferrer"
                class="text-[11px] font-semibold text-emerald-700 hover:underline inline-flex items-center gap-1"
              >
                <span>Buka Midtrans Sandbox Simulator</span>
                <ExternalLink class="w-3 h-3" />
              </a>
            </div>
          </div>

          <div class="pt-4 border-t border-zinc-100 space-y-2">
            <button
              @click="handleManualCheckPayment"
              :disabled="isCheckingPayment || isSubmitting"
              class="w-full py-3.5 rounded-2xl font-bold bg-[#0C3B36] hover:bg-[#134E48] text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <Loader2 v-if="isCheckingPayment" class="w-4 h-4 animate-spin" />
              <RefreshCw v-else class="w-4 h-4" />
              <span>Cek Status Pembayaran</span>
            </button>

            <div v-if="checkPaymentMessage" class="p-2.5 rounded-xl bg-amber-50 border border-amber-200 text-amber-900 text-[11px] font-semibold text-center animate-in fade-in">
              {{ checkPaymentMessage }}
            </div>

            <button
              @click="handleSimulatePayment"
              :disabled="isSubmitting || isCheckingPayment"
              class="w-full py-3 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs active:scale-[0.98] transition-all flex items-center justify-center gap-2"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <CheckCircle2 v-else class="w-4 h-4 text-emerald-600" />
              <span>Simulasikan Pembayaran Berhasil (Demo)</span>
            </button>

            <button
              @click="currentStep = 4"
              class="text-xs font-bold opacity-60 hover:opacity-100 py-1"
            >
              Ubah Data Pendaftaran
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== TICKET SCREEN (BOARDING PASS DIGITAL) ==================== -->
      <div v-else-if="currentStep === 'ticket' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-emerald-100 text-emerald-900">
            <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" /> Reservasi Berhasil & Terkonfirmasi
          </span>
          <h2 class="text-2xl font-black">Tiket Digital Kunjungan</h2>
          <p class="text-xs opacity-75">Simpan layar ini atau konfirmasi tiket telah dikirim ke nomor WhatsApp Anda.</p>
        </div>

        <!-- Boarding Pass Card -->
        <div class="bg-white rounded-3xl border shadow-xs overflow-hidden"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/30' : 'border-[#9E606D]/25'"
        >
          <div class="p-6 text-white text-center space-y-2"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36]' : 'bg-[#9E606D]'"
          >
            <div class="text-[11px] font-bold uppercase tracking-wider text-white/70">
              {{ brand === 'kidinest' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}
            </div>
            <!-- Nomor Antrian / Kode Booking -->
            <div class="text-5xl sm:text-6xl font-black tracking-tight text-white font-mono py-2">
              {{ activeTicket.queue_number }}
            </div>
            <div class="inline-block px-3 py-1 rounded-full text-xs font-bold bg-white/15 text-white">
              Status: {{ activeTicket.appointment_date ? 'Jadwal Terbooking (Check-in Hari H)' : 'Sedang Menunggu Panggilan' }}
            </div>
          </div>

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
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Jadwal Praktik</span>
                <span class="font-bold text-emerald-800">
                  {{ activeTicket.appointment_date ? `${activeTicket.appointment_date} (${activeTicket.appointment_time || '-'})` : 'Hari Ini' }}
                </span>
              </div>
            </div>

            <!-- Fitur Tambahan: Reschedule & Refund Buttons -->
            <div class="flex items-center gap-2 pt-1">
              <button
                type="button"
                @click="openRescheduleModal"
                class="flex-1 py-2.5 px-3 rounded-xl border border-zinc-300 text-zinc-700 font-bold hover:bg-zinc-50 flex items-center justify-center gap-1.5 transition-all text-[11px]"
              >
                <RotateCcw class="w-3.5 h-3.5" />
                <span>Reschedule Jadwal</span>
              </button>

              <button
                v-if="brand === 'kidinest'"
                type="button"
                @click="showRefundModal = true"
                class="flex-1 py-2.5 px-3 rounded-xl border border-rose-300 text-rose-700 font-bold hover:bg-rose-50 flex items-center justify-center gap-1.5 transition-all text-[11px]"
              >
                <ReceiptText class="w-3.5 h-3.5" />
                <span>Ajukan Refund</span>
              </button>
            </div>

            <div class="p-3.5 rounded-2xl bg-emerald-50 border border-emerald-200/80 flex items-center gap-3">
              <MessageCircle class="w-5 h-5 text-emerald-600 shrink-0" />
              <div class="text-[11px] text-emerald-950 font-medium leading-tight">
                Notifikasi dan tracking live tiket terkirim ke WhatsApp <strong>{{ patientPhone }}</strong>.
              </div>
            </div>
          </div>

          <div class="p-4 bg-zinc-50 border-t border-zinc-100 flex items-center justify-between gap-3">
            <button
              @click="emit('navigate', 'monitor')"
              class="inline-flex items-center gap-1.5 text-xs font-bold opacity-75 hover:opacity-100"
            >
              <Smartphone class="w-4 h-4" />
              <span>Lihat TV Monitor</span>
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

    <!-- MODAL RESCHEDULE JADWAL (TERINTEGRASI JADWAL DOKTER) -->
    <div v-if="showRescheduleModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-lg w-full space-y-4 shadow-xl animate-in fade-in">
        <div>
          <h3 class="text-base font-extrabold text-zinc-900">Ubah Jadwal Kunjungan</h3>
          <p class="text-xs text-zinc-600">
            Pilih tanggal & slot praktik yang tersedia bersama <span class="font-bold text-zinc-800">{{ activeTicket?.doctor_name || 'Dokter Spesialis' }}</span>:
          </p>
        </div>

        <!-- Loading Sched State -->
        <div v-if="isLoadingReschedSchedules" class="py-8 flex flex-col items-center justify-center gap-2 text-zinc-500">
          <Loader2 class="w-6 h-6 animate-spin text-[#0C3B36]" />
          <span class="text-xs font-semibold">Memuat jadwal dokter...</span>
        </div>

        <div v-else class="space-y-4">
          <!-- 1. Horizontal Date Picker -->
          <div class="space-y-1.5">
            <label class="block text-[11px] font-bold text-zinc-700">1. Pilih Tanggal Baru:</label>
            <div class="flex gap-2 overflow-x-auto pb-2 pt-1 no-scrollbar scroll-smooth">
              <button
                v-for="item in reschedAvailableDates"
                :key="item.date"
                type="button"
                :disabled="!item.is_available"
                @click="reschedSelectedDateStr = item.date; reschedSelectedTimeSlot = (item.slots.find(s => s.is_available)?.start_time || '')"
                class="shrink-0 w-20 p-2 rounded-2xl border text-center transition-all disabled:opacity-35 disabled:cursor-not-allowed"
                :class="reschedSelectedDateStr === item.date
                  ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#0C3B36] text-white shadow-2xs' : 'border-[#9E606D] bg-[#9E606D] text-white shadow-2xs')
                  : (item.is_available ? 'border-zinc-200 bg-zinc-50 hover:border-zinc-300 text-zinc-800' : 'border-zinc-100 bg-zinc-100 text-zinc-400')"
              >
                <div class="text-[9px] font-bold uppercase tracking-wider">{{ item.day_name }}</div>
                <div class="text-xs font-black my-0.5">{{ item.date.slice(8) }} {{ item.date.slice(5, 7) === '11' ? 'Nov' : 'Des' }}</div>
                <div class="text-[8px] font-medium opacity-80">
                  {{ item.is_available ? (item.shift_name || 'Buka') : 'Libur' }}
                </div>
              </button>
            </div>
          </div>

          <!-- 2. Slot Jam Picker (1 slot = 1 pasien) -->
          <div class="space-y-1.5 pt-1 border-t border-zinc-100">
            <label class="block text-[11px] font-bold text-zinc-700">2. Pilih Jam Slot Kunjungan (1 Pasien per Slot):</label>
            <div v-if="currentReschedSlots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 gap-2 max-h-48 overflow-y-auto no-scrollbar">
              <button
                v-for="slot in currentReschedSlots"
                :key="slot.start_time"
                type="button"
                :disabled="!slot.is_available"
                @click="reschedSelectedTimeSlot = slot.start_time"
                class="p-2 rounded-xl border text-left transition-all text-xs font-bold disabled:opacity-35 disabled:cursor-not-allowed"
                :class="reschedSelectedTimeSlot === slot.start_time
                  ? (brand === 'kidinest' ? 'border-[#0C3B36] bg-[#0C3B36] text-white' : 'border-[#9E606D] bg-[#9E606D] text-white')
                  : (slot.is_available ? 'border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800' : 'border-zinc-100 bg-zinc-50 text-zinc-400')"
              >
                <div class="flex items-center justify-between">
                  <span class="font-mono text-xs">{{ slot.start_time }} - {{ slot.end_time }}</span>
                  <span class="text-[9px] font-bold" :class="slot.is_available ? (reschedSelectedTimeSlot === slot.start_time ? 'text-white' : 'text-emerald-700') : 'text-zinc-400'">
                    {{ slot.is_available ? 'Tersedia' : 'Penuh' }}
                  </span>
                </div>
              </button>
            </div>
            <div v-else class="text-xs text-zinc-400 italic py-2">
              Tidak ada slot praktik yang tersedia di tanggal ini.
            </div>
          </div>
        </div>

        <div class="flex gap-2 pt-2 border-t border-zinc-100">
          <button @click="showRescheduleModal = false" class="flex-1 py-2.5 rounded-xl bg-zinc-100 font-bold text-xs text-zinc-700 hover:bg-zinc-200 transition-colors">
            Batal
          </button>
          <button
            @click="handleConfirmReschedule"
            :disabled="isSubmittingResched || !reschedSelectedDateStr || !reschedSelectedTimeSlot"
            class="flex-1 py-2.5 rounded-xl bg-[#0C3B36] text-white font-bold text-xs hover:bg-[#134E48] transition-all disabled:opacity-50 flex items-center justify-center gap-1.5"
          >
            <Loader2 v-if="isSubmittingResched" class="w-3.5 h-3.5 animate-spin" />
            <span>{{ isSubmittingResched ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL REFUND MANUAL KASIR -->
    <div v-if="showRefundModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-md w-full space-y-4 shadow-xl animate-in fade-in">
        <h3 class="text-base font-extrabold text-zinc-900">Pengajuan Refund Registrasi</h3>
        <p class="text-xs text-zinc-600">
          Kasir klinik akan melakukan transfer manual ke rekening Anda setelah verifikasi:
        </p>
        <div class="space-y-2.5 text-xs">
          <div>
            <label class="block text-[10px] font-bold uppercase text-zinc-600 mb-1">Nama Bank:</label>
            <select v-model="refundBank" class="w-full px-3 py-2 border rounded-xl font-bold">
              <option value="BCA">BCA</option>
              <option value="Mandiri">Bank Mandiri</option>
              <option value="BNI">BNI</option>
              <option value="BRI">BRI</option>
              <option value="BSI">BSI Syariah</option>
              <option value="GoPay/OVO">E-Wallet (GoPay/OVO)</option>
            </select>
          </div>
          <div>
            <label class="block text-[10px] font-bold uppercase text-zinc-600 mb-1">Nomor Rekening / HP E-Wallet:</label>
            <input v-model="refundAccNum" type="text" placeholder="Contoh: 8281928811" class="w-full px-3 py-2 border rounded-xl font-bold font-mono" />
          </div>
          <div>
            <label class="block text-[10px] font-bold uppercase text-zinc-600 mb-1">Nama Pemilik Rekening:</label>
            <input v-model="refundAccName" type="text" placeholder="Contoh: Budi Pratama" class="w-full px-3 py-2 border rounded-xl font-bold" />
          </div>
          <div>
            <label class="block text-[10px] font-bold uppercase text-zinc-600 mb-1">Alasan Pembatalan:</label>
            <textarea v-model="refundNotes" rows="2" placeholder="Contoh: Pasien berhalangan hadir ke luar kota" class="w-full px-3 py-2 border rounded-xl"></textarea>
          </div>
        </div>
        <div class="flex gap-2 pt-2">
          <button @click="showRefundModal = false" class="flex-1 py-2.5 rounded-xl bg-zinc-100 font-bold text-xs text-zinc-700">
            Batal
          </button>
          <button
            @click="handleConfirmRefund"
            :disabled="isSubmittingRefund || !refundAccNum || !refundAccName"
            class="flex-1 py-2.5 rounded-xl bg-rose-700 text-white font-bold text-xs disabled:opacity-50"
          >
            {{ isSubmittingRefund ? 'Mengirim...' : 'Kirim Pengajuan' }}
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL LACAK TIKET MANDIRI -->
    <div v-if="showTrackModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-sm w-full space-y-4 shadow-xl animate-in fade-in">
        <div class="flex items-center justify-between">
          <h3 class="text-base font-extrabold text-zinc-900">Lacak Tiket Saya</h3>
          <button @click="showTrackModal = false" class="text-zinc-400 hover:text-zinc-700 text-sm font-bold">✕</button>
        </div>
        <p class="text-xs text-zinc-600">
          Masukkan nomor WhatsApp yang Anda gunakan saat mendaftar:
        </p>
        <div class="space-y-3">
          <input
            v-model="trackQuery"
            type="text"
            placeholder="Contoh: 081234567890"
            @keyup.enter="handleTrackTicket"
            class="w-full px-4 py-2.5 rounded-xl border border-zinc-200 text-xs font-bold font-mono outline-none focus:border-zinc-800"
          />
          <div v-if="trackError" class="text-xs text-rose-600 font-semibold flex items-center gap-1.5">
            <AlertCircle class="w-3.5 h-3.5 shrink-0" />
            <span>{{ trackError }}</span>
          </div>
          <button
            type="button"
            @click="handleTrackTicket"
            :disabled="isTracking || !trackQuery.trim()"
            class="w-full py-2.5 rounded-xl bg-[#0C3B36] text-white text-xs font-bold flex items-center justify-center gap-1.5 hover:bg-[#134E48] transition-all disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isTracking" class="w-3.5 h-3.5 animate-spin" />
            <Search v-else class="w-3.5 h-3.5" />
            <span>Cari Tiket</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="max-w-2xl mx-auto w-full text-center text-[11px] font-medium opacity-60 pt-4">
      Klinik Utama Semarang Medical Group • Jl. Melati Utara No. 3, Semarang
    </footer>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
