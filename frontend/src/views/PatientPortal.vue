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
  RefreshCw,
  Search,
  Calendar as CalendarIcon,
  Stethoscope,
  RotateCcw,
  ReceiptText,
  FileCheck2,
  Building2,
  Info
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  brand?: 'kidinest' | 'ssc'
}>(), {
  brand: 'kidinest'
})

const emit = defineEmits<{
  (e: 'navigate', view: string): void
}>()

// Step wizard: 1 (Identitas) | 2 (Layanan) | 3 (Tanggal & Jadwal) | 4 (Review) | 'payment' | 'ticket'
const currentStep = ref<number | 'payment' | 'ticket'>(1)

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
  qris_payload?: string
}

// ==================== STEP 1: IDENTITAS PASIEN (SESUAI MASTER PROMPT 3 A) ====================
// Status pasien: Pasien Baru atau Pasien Lama (pilihan langsung di formulir)
const patientType = ref<'BARU' | 'LAMA'>('BARU')
const patientName = ref('Budi Pratama')
const patientPhone = ref('081234567890')
const birthDate = ref('2024-05-10')

// Hitung umur otomatis berdasarkan tanggal lahir
const calculatedAge = computed(() => {
  if (!birthDate.value) return ''
  const birth = new Date(birthDate.value)
  const now = new Date()
  if (isNaN(birth.getTime())) return ''

  let years = now.getFullYear() - birth.getFullYear()
  let months = now.getMonth() - birth.getMonth()
  let days = now.getDate() - birth.getDate()
  if (days < 0) {
    months -= 1
  }
  if (months < 0) {
    years -= 1
    months += 12
  }

  if (years > 0) {
    return `${years} thn ${months} bln`
  }
  if (months > 0) {
    return `${months} bulan`
  }
  return 'Bayi (< 1 bln)'
})

const todayStr = computed(() => {
  return new Date().toISOString().split('T')[0]
})

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

interface ServiceCategory {
  key: string
  code: string
  label: string
  price: number
}

const kidinestCategories: ServiceCategory[] = [
  { key: 'A_SAKIT', code: 'A', label: 'A. Konsultasi Anak Sakit', price: 200000 },
  { key: 'B_VAKSIN', code: 'B', label: 'B. Vaksinasi dan Imunisasi', price: 200000 },
  { key: 'C_TUMBANG', code: 'C', label: 'C. Tumbuh Kembang', price: 200000 },
  { key: 'D_PAKET', code: 'D', label: 'D. Vaksinasi dan Konsultasi Tumbuh Kembang', price: 250000 },
]

const selectedCategory = ref<string>('A_SAKIT')

const selectedService = computed(() => {
  return kidinestCategories.find(c => c.key === selectedCategory.value) || kidinestCategories[0]
})

// Biaya Layanan Akhir (Master Prompt 3 C & 3 G)
const registrationFee = computed<number>(() => selectedService.value.price)
const REGISTRATION_FEE = computed<number>(() => selectedService.value.price)

// ==================== STEP 3: PILIH DOKTER & TANGGAL KUNJUNGAN (SESUAI JADWAL DOKTER TERDAFTAR) ====================
// Ambil langsung dari jadwal dan nama dokter yang terdaftar di sistem (tanpa ruang/room)
export interface DoctorScheduleSummary {
  doctor_id: number
  doctor_name: string
  role: string
  clinic_unit: string
  room_name?: string
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
const isLoadingDoctors = ref(true)

async function fetchDoctorSchedules() {
  isLoadingDoctors.value = true
  try {
    const res = await fetch(`/api/public/doctors/schedules/summary?unit=${selectedUnit.value}`)
    const data = await res.json()
    if (res.ok && data.success && Array.isArray(data.data) && data.data.length > 0) {
      doctorSchedules.value = data.data
      if (!selectedDoctorId.value || !data.data.some((d: any) => d.doctor_id === selectedDoctorId.value)) {
        selectedDoctorId.value = data.data[0].doctor_id
      }
    }
  } catch (err) {
    console.error('Gagal mengambil data jadwal dokter dari backend:', err)
  } finally {
    isLoadingDoctors.value = false
  }
}

watch(() => props.brand, () => {
  fetchDoctorSchedules()
})

onMounted(() => {
  fetchDoctorSchedules()
})

function getDoctorInitials(name: string): string {
  const clean = name.replace(/^dr\.\s+/i, '').replace(/,\s*.*$/, '')
  const parts = clean.split(' ').filter(Boolean)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return (parts[0] ? parts[0].substring(0, 2) : 'DR').toUpperCase()
}

function parseAvailableDaysOfWeek(daysText: string): number[] {
  const lower = (daysText || '').toLowerCase()
  if (lower.includes('sabtu')) {
    return [1, 2, 3, 4, 5, 6]
  }
  return [1, 2, 3, 4, 5]
}

interface DoctorOption {
  id: number
  name: string
  specialty: string
  shift: string
  hours: string
  availableDays: string
  availableDaysOfWeek: number[]
  avatarBg: string
  avatarText: string
  initials: string
}

const availableDoctors = computed<DoctorOption[]>(() => {
  return doctorSchedules.value.map((doc, idx) => {
    let daysOfWeek = parseAvailableDaysOfWeek(doc.days_text || '')
    if (Array.isArray(doc.weekly_schedules) && doc.weekly_schedules.length > 0) {
      daysOfWeek = Array.from(new Set(doc.weekly_schedules.map((w: any) => w.day_of_week)))
    }

    return {
      id: doc.doctor_id,
      name: doc.doctor_name,
      specialty: doc.specialty,
      shift: doc.shift_name ? `Praktik ${doc.shift_name.charAt(0).toUpperCase() + doc.shift_name.slice(1).toLowerCase()}` : 'Praktik',
      hours: doc.time_text || `${doc.start_time} – ${doc.end_time} WIB`,
      availableDays: doc.days_text || 'Senin – Jumat',
      availableDaysOfWeek: daysOfWeek,
      avatarBg: idx % 2 === 0 ? 'bg-[#0C3B36]' : 'bg-[#E5A93C]',
      avatarText: idx % 2 === 0 ? 'text-[#E5A93C]' : 'text-[#0C3B36]',
      initials: getDoctorInitials(doc.doctor_name)
    }
  })
})

const selectedDoctorId = ref<number>(0)

// Dokter bertugas disesuaikan secara otomatis berdasarkan tanggal kunjungan terpilih (Master Prompt 3 D)
const assignedDoctor = computed<DoctorOption>(() => {
  if (availableDoctors.value.length === 0) {
    return {
      id: 0,
      name: isLoadingDoctors.value ? 'Memuat data dokter...' : 'Dokter Spesialis Anak',
      specialty: 'Spesialis Anak (Pediatri)',
      shift: 'Praktik Pagi / Sore',
      hours: '08:30 – 19:30 WIB',
      availableDays: 'Senin – Sabtu',
      availableDaysOfWeek: [1, 2, 3, 4, 5, 6],
      avatarBg: 'bg-[#0C3B36]',
      avatarText: 'text-[#E5A93C]',
      initials: 'DR'
    }
  }
  const selectedD = new Date(selectedDateStr.value)
  const dow = isNaN(selectedD.getTime()) ? 1 : selectedD.getDay()
  const match = availableDoctors.value.find(d => d.availableDaysOfWeek.includes(dow))
  return match || availableDoctors.value[0]
})

const selectedDoctor = computed<DoctorOption>(() => assignedDoctor.value)

function getDefaultVisitDate(): string {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  for (let i = 0; i < 7; i++) {
    const dow = tomorrow.getDay()
    if (dow !== 0) {
      return tomorrow.toISOString().split('T')[0]
    }
    tomorrow.setDate(tomorrow.getDate() + 1)
  }
  return tomorrow.toISOString().split('T')[0]
}

const selectedDateStr = ref<string>(getDefaultVisitDate())

// Generator tanggal 30 hari ke depan
interface VisitDateOption {
  date: string
  dayName: string
  displayDate: string
  dayOfWeek: number
  isSunday: boolean
}

const next30Days = computed<VisitDateOption[]>(() => {
  const days: VisitDateOption[] = []
  const today = new Date()
  const dayNames = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']

  for (let i = 0; i < 30; i++) {
    const d = new Date()
    d.setDate(today.getDate() + i)
    const dow = d.getDay()
    const dateStr = d.toISOString().split('T')[0]
    days.push({
      date: dateStr,
      dayName: dayNames[dow],
      displayDate: `${d.getDate()} ${monthNames[d.getMonth()]}`,
      dayOfWeek: dow,
      isSunday: dow === 0
    })
  }
  return days
})

// Geser / Scroll Kontainer Tanggal (Dukungan Tombol, Drag Mouse, Touch, dan Wheel)
const dateScrollRef = ref<HTMLElement | null>(null)

function scrollDateChips(direction: 'left' | 'right') {
  if (!dateScrollRef.value) return
  const scrollAmount = 280
  dateScrollRef.value.scrollBy({
    left: direction === 'left' ? -scrollAmount : scrollAmount,
    behavior: 'smooth'
  })
}

function handleDateWheel(e: WheelEvent) {
  if (!dateScrollRef.value) return
  if (Math.abs(e.deltaY) > Math.abs(e.deltaX)) {
    e.preventDefault()
    dateScrollRef.value.scrollLeft += e.deltaY * 1.2
  }
}

let isDraggingDate = false
let startDragX = 0
let scrollLeftStart = 0

function onDateMouseDown(e: MouseEvent) {
  if (!dateScrollRef.value) return
  isDraggingDate = true
  startDragX = e.pageX - dateScrollRef.value.offsetLeft
  scrollLeftStart = dateScrollRef.value.scrollLeft
}

function onDateMouseMove(e: MouseEvent) {
  if (!isDraggingDate || !dateScrollRef.value) return
  e.preventDefault()
  const x = e.pageX - dateScrollRef.value.offsetLeft
  const walk = (x - startDragX) * 1.5
  dateScrollRef.value.scrollLeft = scrollLeftStart - walk
}

function onDateMouseUp() {
  isDraggingDate = false
}

function onDateMouseLeave() {
  isDraggingDate = false
}

// Tabel Jadwal Dokter pada tanggal terpilih (READ-ONLY, tidak menampilkan kuota, bebas ruang)
interface DoctorScheduleRow {
  doctorName: string
  specialty: string
  shift: string
}

const doctorScheduleTable = computed<DoctorScheduleRow[]>(() => {
  const selectedD = new Date(selectedDateStr.value)
  const dow = isNaN(selectedD.getTime()) ? 1 : selectedD.getDay()

  if (dow === 0) {
    return []
  }

  return availableDoctors.value
    .filter(doc => doc.availableDaysOfWeek.includes(dow))
    .map(doc => ({
      doctorName: doc.name,
      specialty: doc.specialty,
      shift: `${doc.shift} (${doc.hours})`
    }))
})

// Pengecekan apakah ada dokter spesialis yang berpraktik pada tanggal terpilih
const isDoctorAvailableOnDate = computed(() => {
  const selectedD = new Date(selectedDateStr.value)
  const dow = isNaN(selectedD.getTime()) ? 1 : selectedD.getDay()
  if (dow === 0) return false
  return availableDoctors.value.some(doc => doc.availableDaysOfWeek.includes(dow))
})

// Status Ketersediaan Layanan (Master Prompt 3 E: Tersedia, Hampir Penuh, Penuh, Tidak Melayani)
type CapacityStatus = 'TERSEDIA' | 'HAMPIR_PENUH' | 'PENUH' | 'TIDAK_MELAYANI'

const capacityStatus = computed<{ status: CapacityStatus; label: string; color: string; desc: string }>(() => {
  const selectedD = new Date(selectedDateStr.value)
  const dow = isNaN(selectedD.getTime()) ? 1 : selectedD.getDay()

  // Hari Minggu: Klinik Libur
  if (dow === 0) {
    return {
      status: 'TIDAK_MELAYANI',
      label: 'Tidak Melayani Layanan Tersebut (Klinik Libur Hari Minggu)',
      color: 'bg-zinc-100 text-zinc-600 border-zinc-300',
      desc: 'Klinik Kidinest tidak beroperasi pada hari Minggu. Silakan pilih hari kunjungan Senin s.d. Sabtu.'
    }
  }

  // Jika tidak ada dokter yang melayani pada hari tersebut
  if (!isDoctorAvailableOnDate.value) {
    return {
      status: 'TIDAK_MELAYANI',
      label: 'Tidak Melayani Layanan Tersebut pada Tanggal Ini',
      color: 'bg-amber-100 text-amber-900 border-amber-300',
      desc: `Tidak ada dokter spesialis yang berpraktik pada tanggal terpilih. Silakan pilih tanggal kunjungan lain.`
    }
  }

  return {
    status: 'TERSEDIA',
    label: 'Tersedia',
    color: 'bg-emerald-100 text-emerald-800 border-emerald-300',
    desc: `Layanan ${selectedService.value.label} tersedia pada tanggal ${formatDateDDMMYYYY(selectedDateStr.value)}.`
  }
})

// Checklist Persetujuan Wajib Sebelum Konfirmasi (Master Prompt 3 F)
const agreePrivacy = ref(false)
const agreeTerms = ref(false)
const agreePayment = ref(false)

const allAgreementsChecked = computed(() => {
  return agreePrivacy.value && agreeTerms.value && agreePayment.value
})

// ==================== STEP 4: KONFIRMASI & METODE PEMBAYARAN (SESUAI 3 G, 3 H) ====================
// Metode pembayaran: QRIS & Virtual Account BCA (Midtrans dihilangkan dari tampilan pasien)
const paymentMethod = ref<'QRIS' | 'BANK_TRANSFER'>('QRIS')
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)
const activeTicket = ref<TicketData | null>(null)

// Format tanggal praktis DD-MM-YYYY
function formatDateDDMMYYYY(dateStr?: string): string {
  if (!dateStr) return ''
  const parts = dateStr.split('-')
  if (parts.length === 3 && parts[0].length === 4) {
    return `${parts[2]}-${parts[1]}-${parts[0]}`
  }
  try {
    const d = new Date(dateStr)
    const dd = String(d.getDate()).padStart(2, '0')
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const yyyy = d.getFullYear()
    return `${dd}-${mm}-${yyyy}`
  } catch {
    return dateStr
  }
}

const resolvedQRImage = computed(() => {
  const url = activeTicket.value?.payment_url
  if (url && (url.includes('/qr-code') || url.endsWith('.png') || url.endsWith('.jpg') || url.includes('qrserver') || url.includes('tripay'))) {
    return url
  }
  if (activeTicket.value?.qris_payload) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=${encodeURIComponent(activeTicket.value.qris_payload)}`
  }
  const fee = activeTicket.value?.registration_fee || registrationFee.value
  return `https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=00020101021226680016ID.CO.QRIS.WWW01189360000201100000020215KIDINEST0303UMI51440014ID.CO.QRIS5204581253033605406${fee}5802ID5914KIDINEST_SMG6008SEMARANG6304`
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
        birth_date: birthDate.value,
        category: selectedCategory.value,
        doctor_id: assignedDoctor.value.id,
        doctor_name: assignedDoctor.value.name,
        appointment_date: selectedDateStr.value || undefined,
        payment_method: paymentMethod.value,
        registration_fee: registrationFee.value,
      }),
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Gagal mendaftar antrian')
    }

    activeTicket.value = data.data

    if (data.data.payment_status === 'PENDING') {
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
      } else {
        checkPaymentMessage.value = 'Pembayaran belum terdeteksi. Silakan selesaikan scan QRIS atau transfer BCA.'
      }
    } else {
      checkPaymentMessage.value = 'Belum menerima konfirmasi pelunasan dari jaringan pembayaran.'
    }
  } catch {
    checkPaymentMessage.value = 'Gagal memeriksa status ke gateway. Coba beberapa saat lagi.'
  } finally {
    isCheckingPayment.value = false
  }
}

async function handleSimulatePayment() {
  if (!activeTicket.value || isSubmitting.value) return
  isSubmitting.value = true

  try {
    const res = await fetch('/api/public/queue/verify-payment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        queue_id: activeTicket.value.queue_id || activeTicket.value.id
      }),
    })
    const data = await res.json()
    if (res.ok && data.success && data.data) {
      activeTicket.value = data.data
      stopPaymentPolling()
      currentStep.value = 'ticket'
    } else {
      alert(data.error || 'Gagal verifikasi pembayaran')
    }
  } catch {
    alert('Terjadi kesalahan jaringan')
  } finally {
    isSubmitting.value = false
  }
}

// Tombol Selesai pada tiket digital (Master Prompt 4 C)
function handleFinishTicket() {
  resetForm()
  emit('navigate', 'kidinest')
}

function resetForm() {
  stopPaymentPolling()
  currentStep.value = 1
  patientType.value = 'BARU'
  patientName.value = 'Budi Pratama'
  patientPhone.value = '081234567890'
  birthDate.value = '2024-05-10'
  selectedCategory.value = 'A_SAKIT'
  selectedDateStr.value = getDefaultVisitDate()
  agreePrivacy.value = false
  agreeTerms.value = false
  agreePayment.value = false
  paymentMethod.value = 'QRIS'
  activeTicket.value = null
  errorMessage.value = null
  checkPaymentMessage.value = null
}

// Reschedule Modal
const showRescheduleModal = ref(false)
const newRescheduleDate = ref('')
const isRescheduling = ref(false)

function openRescheduleModal() {
  newRescheduleDate.value = selectedDateStr.value
  showRescheduleModal.value = true
}

async function handleReschedule() {
  if (!activeTicket.value || !newRescheduleDate.value) return
  isRescheduling.value = true
  try {
    const res = await fetch('/api/public/queue/reschedule', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        queue_id: activeTicket.value.queue_id || activeTicket.value.id,
        new_date: newRescheduleDate.value,
        new_time: '09:00',
      }),
    })
    const data = await res.json()
    if (res.ok && data.success) {
      alert(data.message || 'Jadwal kunjungan berhasil diubah!')
      if (activeTicket.value) {
        activeTicket.value.appointment_date = newRescheduleDate.value
      }
      showRescheduleModal.value = false
    } else {
      alert(data.error || 'Gagal mengubah jadwal')
    }
  } catch {
    alert('Terjadi kesalahan koneksi saat reschedule')
  } finally {
    isRescheduling.value = false
  }
}

// Refund Modal
const showRefundModal = ref(false)
const refundBank = ref('BCA')
const refundAccNumber = ref('')
const refundAccName = ref('')
const refundNotes = ref('')
const isSubmittingRefund = ref(false)

async function handleRequestRefund() {
  if (!activeTicket.value || !refundAccNumber.value || !refundAccName.value) {
    alert('Harap lengkapi nomor rekening dan nama pemilik rekening')
    return
  }
  isSubmittingRefund.value = true
  try {
    const res = await fetch('/api/public/queue/refund', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        queue_id: activeTicket.value.queue_id || activeTicket.value.id,
        bank_name: refundBank.value,
        account_number: refundAccNumber.value,
        account_name: refundAccName.value,
        notes: refundNotes.value || 'Permohonan pembatalan mandiri via portal',
      }),
    })
    const data = await res.json()
    if (res.ok && data.success) {
      alert(data.message || 'Pengajuan refund berhasil dicatat. Petugas kasir GRO akan memverifikasi transfer dana.')
      showRefundModal.value = false
      resetForm()
    } else {
      alert(data.error || 'Gagal mengajukan refund')
    }
  } catch {
    alert('Terjadi kesalahan jaringan')
  } finally {
    isSubmittingRefund.value = false
  }
}

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const ticketId = urlParams.get('ticket') || localStorage.getItem('clinic_last_ticket_id')
  if (ticketId) {
    fetch(`/api/public/queue/ticket/${ticketId}`)
      .then(res => res.json())
      .then(data => {
        if (data.success && data.data) {
          activeTicket.value = data.data
          patientName.value = data.data.patient_name || patientName.value
          patientPhone.value = data.data.patient_phone || patientPhone.value
          currentStep.value = 'ticket'
        }
      })
      .catch(() => {})
  }
})

onUnmounted(() => {
  stopPaymentPolling()
})
</script>

<template>
  <div class="min-h-screen flex flex-col justify-between p-4 sm:p-8 font-['Plus_Jakarta_Sans',sans-serif] select-none bg-[#FAF7F2] text-[#0C3B36]">
    
    <!-- Top Bar Navigation -->
    <header class="max-w-2xl mx-auto w-full flex items-center justify-between pb-4">
      <button
        @click="emit('navigate', 'kidinest')"
        class="inline-flex items-center gap-2 text-xs font-bold opacity-75 hover:opacity-100 transition-opacity"
      >
        <ArrowLeft class="w-3.5 h-3.5" />
        <span>Kembali ke Beranda Kidinest</span>
      </button>

      <div class="flex items-center gap-2">
        <button
          type="button"
          @click="showTrackModal = true"
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold border transition-colors cursor-pointer bg-[#0C3B36]/5 hover:bg-[#0C3B36]/10 text-[#0C3B36] border-[#0C3B36]/20"
        >
          <Search class="w-3 h-3" />
          <span>Lacak Tiket</span>
        </button>
        <span class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-bold border bg-white border-[#E5A93C]/30 text-[#0C3B36] shadow-2xs">
          <ShieldCheck class="w-3 h-3 text-emerald-600" />
          <span>Semarang Medical Group</span>
        </span>
      </div>
    </header>

    <!-- Main Form Container -->
    <main class="max-w-2xl mx-auto w-full py-2 flex-1 flex flex-col justify-center">

      <!-- SSC BRAND GATE: PUBLIC ONLINE REGISTRATION DISABLED (MASTER PROMPT 3 B) -->
      <div v-if="props.brand === 'ssc'" class="bg-white p-8 rounded-3xl border border-zinc-200 shadow-xs text-center space-y-5 animate-in fade-in">
        <div class="w-16 h-16 mx-auto rounded-3xl bg-amber-100 text-amber-900 flex items-center justify-center shadow-2xs">
          <Building2 class="w-8 h-8" />
        </div>
        <div class="space-y-1.5">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold bg-amber-100 text-amber-900">
            <span>Akses Publik Dibatasi</span>
          </div>
          <h2 class="text-2xl font-extrabold text-[#2D1E24]">Pendaftaran Online SSC Belum Dibuka</h2>
          <p class="text-xs text-zinc-600 max-w-md mx-auto leading-relaxed">
            Halaman pendaftaran online publik untuk Semarang Skin Center (SSC) belum diaktifkan sesuai kebijakan Semarang Medical Group. Pendaftaran pasien SSC dilayani langsung di klinik atau melalui modul Admin / GRO internal.
          </p>
        </div>
        <div class="pt-3 flex flex-col sm:flex-row items-center justify-center gap-3">
          <button
            type="button"
            @click="emit('navigate', 'ssc')"
            class="w-full sm:w-auto px-5 py-3 rounded-2xl border border-zinc-200 text-xs font-bold text-zinc-700 hover:bg-zinc-50 transition cursor-pointer"
          >
            Kembali ke Beranda SSC
          </button>
          <button
            type="button"
            @click="emit('navigate', 'admin')"
            class="w-full sm:w-auto px-5 py-3 rounded-2xl bg-[#2D1E24] text-white text-xs font-bold hover:bg-[#3E2A32] transition cursor-pointer shadow-sm"
          >
            Akses Halaman Admin SSC
          </button>
        </div>
      </div>

      <template v-else>
      <!-- STEP PROGRESS INDICATOR (1 -> 2 -> 3 -> 4) -->
      <div v-if="typeof currentStep === 'number'" class="mb-6 space-y-2">
        <div class="flex items-center justify-between px-2 text-[11px] font-extrabold uppercase tracking-wider">
          <span :class="currentStep >= 1 ? 'text-[#0C3B36]' : 'opacity-40'">
            1. Pasien
          </span>
          <span :class="currentStep >= 2 ? 'text-[#0C3B36]' : 'opacity-40'">
            2. Layanan
          </span>
          <span :class="currentStep >= 3 ? 'text-[#0C3B36]' : 'opacity-40'">
            3. Jadwal Dokter
          </span>
          <span :class="currentStep >= 4 ? 'text-[#0C3B36]' : 'opacity-40'">
            4. Review & Bayar
          </span>
        </div>
        <div class="w-full bg-zinc-200/80 h-1.5 rounded-full overflow-hidden">
          <div
            class="h-full bg-[#0C3B36] transition-all duration-300 rounded-full"
            :style="{ width: `${(currentStep / 4) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- ==================== STEP 1: IDENTITAS PASIEN (MASTER PROMPT 3 A) ==================== -->
      <div v-if="currentStep === 1" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs bg-[#E5A93C] text-[#0C3B36]">
            <span>Tahap 1 • Identitas Pasien</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Pendaftaran Online Pasien Anak
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Lengkapi data diri pasien anak/bayi untuk penerbitan nomor antrean resmi.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-5">
          
          <!-- Status Pasien: Pasien Baru atau Pasien Lama (Pilihan Langsung di Form - Master Prompt 3 A) -->
          <div class="space-y-1.5">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              Status Pasien *
            </label>
            <div class="grid grid-cols-2 gap-3">
              <label
                class="flex items-center justify-center gap-2 p-3.5 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="patientType === 'BARU'
                  ? 'border-[#0C3B36] bg-[#FAF7F2] text-[#0C3B36] ring-2 ring-[#0C3B36]/30 shadow-2xs font-black'
                  : 'border-zinc-200 bg-white text-zinc-600 hover:border-zinc-300'"
              >
                <input type="radio" v-model="patientType" value="BARU" class="accent-[#0C3B36]" />
                <span>Pasien Baru</span>
              </label>

              <label
                class="flex items-center justify-center gap-2 p-3.5 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="patientType === 'LAMA'
                  ? 'border-[#0C3B36] bg-[#FAF7F2] text-[#0C3B36] ring-2 ring-[#0C3B36]/30 shadow-2xs font-black'
                  : 'border-zinc-200 bg-white text-zinc-600 hover:border-zinc-300'"
              >
                <input type="radio" v-model="patientType" value="LAMA" class="accent-[#0C3B36]" />
                <span>Pasien Lama</span>
              </label>
            </div>
          </div>

          <!-- Kolom: Nama Pasien (Anak/Bayi) -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
              Nama Pasien (Anak / Bayi) *
            </label>
            <input
              v-model="patientName"
              type="text"
              required
              placeholder="Contoh: Ananda Mikaela Pratama"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]"
            />
          </div>

          <!-- Kolom: Tanggal Lahir (Umur Dihitung Otomatis) -->
          <div>
            <div class="flex items-center justify-between mb-1.5">
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
                Tanggal Lahir Pasien *
              </label>
              <span v-if="calculatedAge" class="text-[11px] font-black text-emerald-800 bg-emerald-50 px-2.5 py-0.5 rounded-full border border-emerald-200">
                Umur: {{ calculatedAge }}
              </span>
            </div>
            <input
              v-model="birthDate"
              type="date"
              required
              :max="todayStr"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36] font-mono"
            />
          </div>

          <!-- Kolom: Nomor WhatsApp Aktif -->
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75 mb-1.5">
              Nomor WhatsApp / HP Aktif *
            </label>
            <input
              v-model="patientPhone"
              type="tel"
              required
              placeholder="Contoh: 081234567890"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-bold outline-none focus:bg-white transition-all bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36] font-mono"
            />
            <p class="text-[10px] opacity-60 mt-1">
              Tiket digital resmi dan update antrean akan dikirimkan ke nomor WhatsApp ini.
            </p>
          </div>

          <button
            type="button"
            @click="currentStep = 2"
            :disabled="!patientName.trim() || !patientPhone.trim() || !birthDate"
            class="w-full py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-40 bg-[#0C3B36] hover:bg-[#134E48] cursor-pointer"
          >
            <span>Lanjut: Pilih Layanan</span>
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
      </div>

      <!-- ==================== STEP 2: PILIHAN LAYANAN (MASTER PROMPT 3 C) ==================== -->
      <div v-else-if="currentStep === 2" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs bg-[#E5A93C] text-[#0C3B36]">
            <span>Tahap 2 • Pilihan Layanan</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Pilihan Layanan Kidinest
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Pilih jenis penanganan pediatri yang dibutuhkan si kecil.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-4">
          <!-- 4 Layanan Resmi Tanpa Keterangan Tambahan di Bawah Nama Layanan (Master Prompt 3 C) -->
          <div class="space-y-3">
            <label
              v-for="cat in kidinestCategories"
              :key="cat.key"
              class="flex items-center justify-between p-4 rounded-2xl border cursor-pointer transition-all"
              :class="selectedCategory === cat.key
                ? 'border-[#0C3B36] bg-[#FAF7F2] shadow-2xs font-extrabold ring-1 ring-[#0C3B36]'
                : 'border-zinc-200 hover:border-zinc-300 bg-white'"
            >
              <div class="flex items-center gap-3">
                <input
                  type="radio"
                  v-model="selectedCategory"
                  :value="cat.key"
                  class="accent-[#0C3B36] w-4 h-4"
                />
                <span class="text-xs font-bold text-[#0C3B36]">{{ cat.label }}</span>
              </div>
            </label>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 1"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5 cursor-pointer"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="currentStep = 3"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 bg-[#0C3B36] hover:bg-[#134E48] cursor-pointer"
            >
              <span>Lanjut: Pilih Tanggal Kunjungan</span>
              <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== STEP 3: PILIH TANGGAL KUNJUNGAN & JADWAL DOKTER (MASTER PROMPT 3 D, 3 E) ==================== -->
      <div v-else-if="currentStep === 3" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs bg-[#E5A93C] text-[#0C3B36]">
            <span>Tahap 3 • Tanggal Kunjungan</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Pilih Tanggal Kunjungan
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Pasien cukup memilih tanggal rencana kedatangan. Jadwal dokter yang bertugas ditampilkan secara informatif di bawah.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-6">
          
          <!-- 1. Pemilihan Tanggal Kunjungan (Tanpa Memilih Jam / Sesi) -->
          <div class="space-y-2">
            <div class="flex items-center justify-between gap-2">
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
                Pilih Tanggal Kunjungan:
              </label>

              <!-- Kontrol Geser Tanggal (Tombol Panah & Indikator Tanggal DD-MM-YYYY) -->
              <div class="flex items-center gap-2">
                <span class="text-[11px] font-bold text-zinc-600 hidden sm:inline">
                  Terpilih: <strong class="text-[#0C3B36] font-mono">{{ formatDateDDMMYYYY(selectedDateStr) }}</strong>
                </span>

                <div class="flex items-center gap-1">
                  <button
                    type="button"
                    @click="scrollDateChips('left')"
                    class="p-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-100 text-zinc-700 transition cursor-pointer active:scale-95 shadow-2xs"
                    title="Geser Tanggal ke Kiri"
                  >
                    <ChevronLeft class="w-3.5 h-3.5" />
                  </button>
                  <button
                    type="button"
                    @click="scrollDateChips('right')"
                    class="p-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-100 text-zinc-700 transition cursor-pointer active:scale-95 shadow-2xs"
                    title="Geser Tanggal ke Kanan"
                  >
                    <ChevronRight class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Horizontal Date Picker Chips -->
            <div
              ref="dateScrollRef"
              @wheel.prevent="handleDateWheel"
              @mousedown="onDateMouseDown"
              @mousemove="onDateMouseMove"
              @mouseup="onDateMouseUp"
              @mouseleave="onDateMouseLeave"
              class="flex gap-2 overflow-x-auto pb-3 pt-1 scroll-smooth select-none cursor-grab active:cursor-grabbing date-chips-container"
            >
              <button
                v-for="item in next30Days"
                :key="item.date"
                type="button"
                @click="selectedDateStr = item.date"
                class="shrink-0 w-24 p-2.5 rounded-2xl border text-center transition-all cursor-pointer"
                :class="selectedDateStr === item.date
                  ? 'border-[#0C3B36] bg-[#0C3B36] text-white shadow-xs font-extrabold ring-2 ring-[#0C3B36]/30'
                  : (item.isSunday ? 'border-zinc-200 bg-zinc-100 text-zinc-400 opacity-60' : 'border-zinc-200 bg-[#FAF7F2] hover:border-zinc-300 text-zinc-800')"
              >
                <div class="text-[10px] font-bold uppercase tracking-wider">{{ item.dayName }}</div>
                <div class="text-sm font-black my-0.5">{{ item.displayDate }}</div>
                <div class="text-[9px] font-medium" :class="selectedDateStr === item.date ? 'text-[#E5A93C]' : 'opacity-80'">
                  {{ item.isSunday ? 'Libur' : 'Buka' }}
                </div>
              </button>
            </div>

            <div class="flex items-center justify-between text-[10px] text-zinc-500 pt-0.5">
              <span class="sm:hidden font-mono font-bold text-[#0C3B36]">Terpilih: {{ formatDateDDMMYYYY(selectedDateStr) }}</span>
              <span class="opacity-70">💡 Geser tombol ‹ › atau usap kartu tanggal untuk memilih hari lain</span>
            </div>
          </div>

          <!-- 2. Status Ketersediaan Layanan (Master Prompt 3 E) -->
          <div class="p-4 rounded-2xl border transition-all" :class="capacityStatus.color">
            <div class="flex items-center justify-between gap-2">
              <div class="flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full"
                  :class="capacityStatus.status === 'TERSEDIA' ? 'bg-emerald-600' : (capacityStatus.status === 'HAMPIR_PENUH' ? 'bg-amber-600' : 'bg-rose-600')"
                ></span>
                <span class="text-xs font-black uppercase tracking-wider">
                  Status Ketersediaan: {{ capacityStatus.label }}
                </span>
              </div>
              <span class="text-[11px] font-bold opacity-80">{{ selectedService.label }}</span>
            </div>
            <p class="text-xs font-medium mt-1 leading-relaxed opacity-90">
              {{ capacityStatus.desc }}
            </p>
          </div>

          <!-- 3. Tabel Jadwal Dokter pada Tanggal Tersebut (Read-Only, Informatif, Tanpa Kuota) -->
          <div class="space-y-2.5 pt-2 border-t border-zinc-100">
            <div class="flex items-center justify-between">
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
                Jadwal Praktik Dokter Spesialis pada Tanggal Tersebut (Informatif):
              </label>
              <span class="text-[10px] text-zinc-400 font-medium">Read-Only</span>
            </div>

            <div v-if="doctorScheduleTable.length === 0" class="p-6 rounded-2xl bg-zinc-50 border border-zinc-200 text-center text-xs text-zinc-500">
              Klinik libur pada tanggal ini. Tidak ada dokter spesialis yang bertugas. Silakan pilih tanggal lain.
            </div>

            <div v-else class="overflow-x-auto rounded-2xl border border-zinc-200 bg-white">
              <table class="w-full text-left text-xs">
                <thead class="bg-[#FAF7F2] border-b border-zinc-200 text-[10px] font-extrabold uppercase tracking-wider text-[#0C3B36]/80">
                  <tr>
                    <th class="py-2.5 px-3">Nama Dokter Spesialis</th>
                    <th class="py-2.5 px-3">Jadwal Praktik</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-zinc-100">
                  <tr
                    v-for="(doc, idx) in doctorScheduleTable"
                    :key="idx"
                    class="transition-colors hover:bg-zinc-50/80"
                  >
                    <td class="py-3 px-3">
                      <div class="font-extrabold text-[#0C3B36]">
                        {{ doc.doctorName }}
                      </div>
                      <div class="text-[10px] text-zinc-500 font-normal">{{ doc.specialty }}</div>
                    </td>
                    <td class="py-3 px-3 font-semibold text-zinc-800">
                      {{ doc.shift }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- 4. Persetujuan Sebelum Konfirmasi (Master Prompt 3 F) -->
          <div class="space-y-2.5 pt-3 border-t border-zinc-100">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              Persetujuan Pendaftaran (Wajib Dicentang):
            </label>

            <div class="space-y-2 text-xs">
              <label class="flex items-start gap-2.5 p-2.5 rounded-xl border border-zinc-200 bg-[#FAF7F2]/60 cursor-pointer hover:bg-[#FAF7F2]">
                <input type="checkbox" v-model="agreePrivacy" class="mt-0.5 w-4 h-4 accent-[#0C3B36]" />
                <span class="font-medium text-zinc-800">
                  Saya menyetujui <strong>Kebijakan Privasi</strong> dan perlindungan data rekam medis pasien.
                </span>
              </label>

              <label class="flex items-start gap-2.5 p-2.5 rounded-xl border border-zinc-200 bg-[#FAF7F2]/60 cursor-pointer hover:bg-[#FAF7F2]">
                <input type="checkbox" v-model="agreeTerms" class="mt-0.5 w-4 h-4 accent-[#0C3B36]" />
                <span class="font-medium text-zinc-800">
                  Saya menyetujui <strong>Ketentuan Layanan</strong> dan tata tertib klinik Semarang Medical Group.
                </span>
              </label>

              <label class="flex items-start gap-2.5 p-2.5 rounded-xl border border-zinc-200 bg-[#FAF7F2]/60 cursor-pointer hover:bg-[#FAF7F2]">
                <input type="checkbox" v-model="agreePayment" class="mt-0.5 w-4 h-4 accent-[#0C3B36]" />
                <span class="font-medium text-zinc-800">
                  Saya bersedia melakukan <strong>pembayaran administrasi</strong> sesuai ketentuan layanan.
                </span>
              </label>
            </div>
          </div>

          <!-- Tombol Lanjut ke Review & Konfirmasi -->
          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 2"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5 cursor-pointer"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="currentStep = 4"
              :disabled="!allAgreementsChecked || !isDoctorAvailableOnDate || capacityStatus.status === 'TIDAK_MELAYANI' || capacityStatus.status === 'PENUH'"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-40 disabled:cursor-not-allowed bg-[#0C3B36] hover:bg-[#134E48] cursor-pointer"
            >
              <span>Review & Konfirmasi</span>
              <ChevronRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== STEP 4: KONFIRMASI RINGKASAN & METODE BAYAR (MASTER PROMPT 3 G, 3 H) ==================== -->
      <div v-else-if="currentStep === 4" class="space-y-5 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <div class="inline-flex items-center gap-2 px-3 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs bg-[#E5A93C] text-[#0C3B36]">
            <span>Tahap 4 • Review & Konfirmasi</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight">
            Konfirmasi Pendaftaran
          </h1>
          <p class="text-xs opacity-75 max-w-md mx-auto">
            Periksa detail pasien dan jadwal kunjungan sebelum menerbitkan nomor antrean.
          </p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-5">
          <!-- Alert Error jika ada -->
          <div v-if="errorMessage" class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- Ringkasan Tiket (Format Tanggal DD-MM-YYYY) -->
          <div class="p-4 rounded-2xl bg-[#FAF7F2] border border-zinc-200/80 space-y-3 text-xs">
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Unit Faskes:</span>
              <span class="font-extrabold text-[#0C3B36]">Kidinest Pediatric Care</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Dokter Spesialis:</span>
              <span class="font-extrabold text-[#0C3B36]">{{ selectedDoctor.name }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Sesi Praktik:</span>
              <span class="font-bold text-zinc-800">{{ selectedDoctor.shift }} ({{ selectedDoctor.hours }})</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Nama Pasien:</span>
              <span class="font-bold">{{ patientName }} ({{ patientType === 'BARU' ? 'Pasien Baru' : 'Pasien Lama' }})</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Tanggal Lahir & Umur:</span>
              <span class="font-bold">{{ formatDateDDMMYYYY(birthDate) }} ({{ calculatedAge }})</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Nomor WhatsApp:</span>
              <span class="font-mono font-bold">{{ patientPhone }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Layanan:</span>
              <span class="font-bold text-[#0C3B36]">{{ selectedService.label }}</span>
            </div>
            <div class="flex items-center justify-between pb-2 border-b border-zinc-200/60">
              <span class="opacity-70">Tanggal Kunjungan:</span>
              <span class="font-bold text-emerald-800 font-mono">{{ formatDateDDMMYYYY(selectedDateStr) }}</span>
            </div>
            <div class="flex items-center justify-between pt-1 text-sm font-black">
              <span>Biaya Administrasi Reservasi:</span>
              <span class="text-emerald-700 font-mono">Rp {{ REGISTRATION_FEE.toLocaleString('id-ID') }}</span>
            </div>
          </div>

          <!-- Informasi Pembayaran Wajib (Master Prompt 3 G) -->
          <div class="p-3.5 rounded-2xl bg-amber-50/70 border border-amber-200/80 flex items-center gap-2.5">
            <Info class="w-4 h-4 text-amber-700 shrink-0" />
            <p class="text-xs font-bold text-amber-950">
              Pasien tidak dikenakan biaya administrasi tambahan pada saat kedatangan.
            </p>
          </div>

          <!-- Pilihan Metode Pembayaran Resmi (QRIS & Virtual Account BCA saja) -->
          <div class="space-y-2">
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-75">
              Pilih Metode Pembayaran:
            </label>
            <div class="grid grid-cols-2 gap-2">
              <label class="flex items-center gap-2 p-3.5 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'QRIS' ? 'border-[#0C3B36] bg-[#FAF7F2] ring-1 ring-[#0C3B36]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="QRIS" class="accent-[#0C3B36]" />
                <QrCode class="w-4 h-4 text-emerald-600" />
                <span>QRIS</span>
              </label>
              <label class="flex items-center gap-2 p-3.5 rounded-2xl border cursor-pointer text-xs font-bold transition-all"
                :class="paymentMethod === 'BANK_TRANSFER' ? 'border-[#0C3B36] bg-[#FAF7F2] ring-1 ring-[#0C3B36]' : 'border-zinc-200 bg-white'"
              >
                <input type="radio" v-model="paymentMethod" value="BANK_TRANSFER" class="accent-[#0C3B36]" />
                <CreditCard class="w-4 h-4 text-sky-600" />
                <span>Virtual Account BCA</span>
              </label>
            </div>
          </div>

          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="currentStep = 3"
              class="py-3.5 px-5 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs transition-all flex items-center gap-1.5 cursor-pointer"
            >
              <ChevronLeft class="w-4 h-4" />
              <span>Kembali</span>
            </button>
            <button
              type="button"
              @click="handleRegister"
              :disabled="isSubmitting"
              class="flex-1 py-3.5 rounded-2xl font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50 bg-[#0C3B36] hover:bg-[#134E48] cursor-pointer"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <span v-else>Konfirmasi & Bayar Rp {{ REGISTRATION_FEE.toLocaleString('id-ID') }}</span>
              <ChevronRight v-if="!isSubmitting" class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== PAYMENT SCREEN (MASTER PROMPT 3 G) ==================== -->
      <div v-else-if="currentStep === 'payment' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-amber-100 text-amber-900">
            <Clock class="w-3.5 h-3.5" /> Menunggu Pembayaran Registrasi
          </span>
          <h2 class="text-2xl font-black text-[#0C3B36]">
            {{ paymentMethod === 'QRIS' ? 'Selesaikan Pembayaran QRIS' : 'Selesaikan Pembayaran Virtual Account BCA' }}
          </h2>
          <p class="text-xs opacity-75">Nomor antrean resmi aktif secara otomatis setelah pembayaran diverifikasi.</p>
        </div>

        <div class="bg-white p-6 sm:p-8 rounded-3xl border border-[#E5A93C]/25 shadow-xs space-y-6 text-center">
          
          <!-- Banner Teks Wajib -->
          <div class="p-3.5 rounded-2xl bg-amber-50 border border-amber-200 text-xs font-bold text-amber-950 text-center">
            Pasien tidak dikenakan biaya administrasi tambahan pada saat kedatangan.
          </div>

          <div class="p-4 rounded-2xl bg-[#FAF7F2] border border-[#E5A93C]/20 text-left space-y-1.5 text-xs">
            <div class="flex items-center justify-between opacity-70">
              <span>Pasien Anak:</span>
              <span class="font-bold text-[#0C3B36]">{{ patientName }}</span>
            </div>
            <div class="flex items-center justify-between opacity-70">
              <span>Dokter Spesialis:</span>
              <span class="font-bold text-[#0C3B36]">{{ selectedDoctor.name }}</span>
            </div>
            <div class="flex items-center justify-between opacity-70">
              <span>Layanan:</span>
              <span class="font-bold text-[#0C3B36]">{{ selectedService.label }}</span>
            </div>
            <div class="flex items-center justify-between opacity-70">
              <span>Tanggal Kunjungan:</span>
              <span class="font-bold text-[#0C3B36] font-mono">{{ formatDateDDMMYYYY(activeTicket.appointment_date || selectedDateStr) }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-[#E5A93C]/20 text-sm font-extrabold text-[#0C3B36]">
              <span>Biaya Administrasi Reservasi:</span>
              <span class="text-emerald-700 font-mono">Rp {{ (activeTicket.registration_fee || REGISTRATION_FEE).toLocaleString('id-ID') }}</span>
            </div>
          </div>

          <!-- QRIS atau Virtual Account Display -->
          <div v-if="paymentMethod === 'QRIS'" class="p-6 rounded-3xl bg-[#FAF7F2] border border-[#E5A93C]/20 inline-block space-y-3 max-w-sm mx-auto">
            <div class="w-60 h-60 mx-auto bg-white p-3 rounded-2xl border border-zinc-200 flex items-center justify-center shadow-xs overflow-hidden">
              <img
                v-if="resolvedQRImage"
                :src="resolvedQRImage"
                alt="QRIS Pembayaran Resmi"
                class="w-56 h-56 object-contain"
              />
              <QrCode v-else class="w-44 h-44 text-[#0C3B36]" />
            </div>

            <div class="text-[11px] font-mono font-bold text-zinc-600">
              NMID: ID1020039201990 (Kidinest Pediatric Care)
            </div>
            <div class="text-[10px] text-zinc-400">Scan melalui BCA mobile, Livin, GoPay, OVO, ShopeePay & Semua Aplikasi QRIS</div>
            <div class="inline-flex items-center gap-1.5 text-[10px] text-emerald-700 font-semibold bg-emerald-50 px-2.5 py-1 rounded-full border border-emerald-200/50">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-ping"></span>
              Menunggu konfirmasi pembayaran real-time
            </div>
          </div>

          <div v-else class="p-6 rounded-3xl bg-[#FAF7F2] border border-[#E5A93C]/20 max-w-sm mx-auto space-y-3 text-left">
            <div class="flex items-center gap-2">
              <CreditCard class="w-5 h-5 text-sky-600" />
              <span class="text-xs font-bold text-zinc-800">Virtual Account BCA</span>
            </div>
            <div class="p-3 bg-white rounded-xl border border-zinc-200 font-mono text-base font-black text-zinc-900 tracking-wider">
              8277 0812 3456 7890
            </div>
            <div class="text-[11px] text-zinc-500">
              Nama Rekening: <strong>SMG KIDINEST CARE</strong><br />
              Nominal tepat: <strong>Rp {{ (activeTicket.registration_fee || REGISTRATION_FEE).toLocaleString('id-ID') }}</strong>
            </div>
          </div>

          <div class="pt-4 border-t border-zinc-100 space-y-2">
            <button
              @click="handleManualCheckPayment"
              :disabled="isCheckingPayment || isSubmitting"
              class="w-full py-3.5 rounded-2xl font-bold bg-[#0C3B36] hover:bg-[#134E48] text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
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
              class="w-full py-3 rounded-2xl font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 text-xs active:scale-[0.98] transition-all flex items-center justify-center gap-2 cursor-pointer"
            >
              <Loader2 v-if="isSubmitting" class="w-4 h-4 animate-spin" />
              <CheckCircle2 v-else class="w-4 h-4 text-emerald-600" />
              <span>Simulasikan Pembayaran Berhasil (Demo)</span>
            </button>
          </div>
        </div>
      </div>

      <!-- ==================== STEP TICKET: BOARDING PASS DIGITAL (MASTER PROMPT 4 A, 4 C) ==================== -->
      <div v-else-if="currentStep === 'ticket' && activeTicket" class="space-y-6 animate-in fade-in duration-200">
        <div class="text-center space-y-1">
          <span class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-[11px] font-bold bg-emerald-100 text-emerald-900">
            <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" /> Pendaftaran Berhasil & Terkonfirmasi
          </span>
          <h2 class="text-2xl font-black">Tiket Digital Kunjungan</h2>
          <p class="text-xs opacity-75">Tunjukkan tiket ini saat kedatangan di klinik.</p>
        </div>

        <!-- Boarding Pass Card -->
        <div class="bg-white rounded-3xl border border-[#E5A93C]/30 shadow-xs overflow-hidden">
          <div class="p-6 text-white text-center space-y-2 bg-[#0C3B36]">
            <div class="text-[11px] font-bold uppercase tracking-wider text-white/70">
              Kidinest Pediatric Care • Semarang Medical Group
            </div>
            
            <!-- Format Nomor Antrean: 2 angka + kode huruf (08-A, 09-C, 10-C) -->
            <div class="text-5xl sm:text-6xl font-black tracking-tight text-white font-mono py-2">
              {{ activeTicket.queue_number }}
            </div>

            <div class="inline-block px-3 py-1 rounded-full text-xs font-bold bg-white/15 text-white">
              Status: {{ activeTicket.payment_status === 'PAID' ? 'Terkonfirmasi Lunas' : 'Menunggu Pelunasan' }}
            </div>
          </div>

          <div class="p-6 space-y-4 text-xs">
            
            <!-- Keterangan Wajib Master Prompt 4 A -->
            <div class="p-3.5 rounded-2xl bg-amber-50 border border-amber-200/80 text-amber-950 text-[11px] font-medium leading-relaxed">
              Konfirmasi nomor antrean dapat dilakukan melalui Customer Service pada hari kunjungan. Nomor antrean juga dapat dipantau secara langsung melalui sistem ini pada halaman utama.
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Pasien</span>
                <span class="font-extrabold text-sm">{{ patientName }} ({{ calculatedAge }})</span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">WhatsApp</span>
                <span class="font-mono font-bold text-sm">{{ patientPhone }}</span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Dokter Spesialis</span>
                <span class="font-extrabold text-sm text-[#0C3B36]">{{ activeTicket.doctor_name || selectedDoctor.name }}</span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Sesi Praktik</span>
                <span class="font-bold text-sm text-zinc-800">{{ selectedDoctor.shift }} ({{ selectedDoctor.hours }})</span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Layanan</span>
                <span class="font-bold">{{ selectedService.label }}</span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Tanggal Praktik</span>
                <span class="font-bold text-emerald-800 font-mono text-sm">
                  {{ formatDateDDMMYYYY(activeTicket.appointment_date || selectedDateStr) }}
                </span>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4 pb-4 border-b border-zinc-100">
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Biaya Administrasi</span>
                <span class="font-bold text-emerald-800 font-mono text-sm">
                  Rp {{ (activeTicket.registration_fee || REGISTRATION_FEE).toLocaleString('id-ID') }}
                </span>
              </div>
              <div>
                <span class="opacity-60 font-semibold block text-[10px] uppercase">Status Pembayaran</span>
                <span class="inline-flex items-center gap-1 text-xs font-bold text-emerald-800 bg-emerald-50 px-2 py-0.5 rounded-md mt-0.5">
                  <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" />
                  <span>Lunas</span>
                </span>
              </div>
            </div>

            <!-- Reschedule & Refund Actions -->
            <div class="flex items-center gap-2 pt-1">
              <button
                type="button"
                @click="openRescheduleModal"
                class="flex-1 py-2.5 px-3 rounded-xl border border-zinc-300 text-zinc-700 font-bold hover:bg-zinc-50 flex items-center justify-center gap-1.5 transition-all text-[11px] cursor-pointer"
              >
                <RotateCcw class="w-3.5 h-3.5" />
                <span>Reschedule Jadwal</span>
              </button>

              <button
                type="button"
                @click="showRefundModal = true"
                class="flex-1 py-2.5 px-3 rounded-xl border border-rose-300 text-rose-700 font-bold hover:bg-rose-50 flex items-center justify-center gap-1.5 transition-all text-[11px] cursor-pointer"
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

          <!-- Bottom Footer Bar dengan Tombol "Selesai" (Master Prompt 4 C) -->
          <div class="p-4 bg-zinc-50 border-t border-zinc-100 flex items-center justify-between gap-3">
            <button
              @click="emit('navigate', 'monitor')"
              class="inline-flex items-center gap-1.5 text-xs font-bold opacity-75 hover:opacity-100 cursor-pointer"
            >
              <Smartphone class="w-4 h-4" />
              <span>Pantau Monitor TV</span>
            </button>
            
            <button
              @click="handleFinishTicket"
              class="px-6 py-2.5 rounded-xl text-white text-xs font-extrabold uppercase tracking-wider transition-all shadow-xs bg-[#0C3B36] hover:bg-[#134E48] cursor-pointer"
            >
              Selesai
            </button>
          </div>
        </div>
      </div>
      </template>

    </main>

    <!-- MODAL RESCHEDULE JADWAL -->
    <div v-if="showRescheduleModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-lg w-full space-y-4 shadow-xl animate-in fade-in">
        <div>
          <h3 class="text-base font-extrabold text-zinc-900">Ubah Tanggal Kunjungan</h3>
          <p class="text-xs text-zinc-500 mt-0.5">Pilih tanggal baru untuk jadwal kunjungan Anda.</p>
        </div>

        <div class="space-y-2">
          <label class="text-[11px] font-bold text-zinc-600 uppercase">Tanggal Kunjungan Baru:</label>
          <input
            v-model="newRescheduleDate"
            type="date"
            class="w-full px-3 py-2.5 rounded-xl border border-zinc-300 text-xs font-bold font-mono"
          />
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            @click="showRescheduleModal = false"
            class="flex-1 py-2.5 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-xs font-bold text-zinc-700"
          >
            Batal
          </button>
          <button
            type="button"
            @click="handleReschedule"
            :disabled="isRescheduling"
            class="flex-1 py-2.5 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-xs font-bold text-white flex items-center justify-center gap-1.5"
          >
            <Loader2 v-if="isRescheduling" class="w-3.5 h-3.5 animate-spin" />
            <span>Simpan Jadwal Baru</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL REFUND -->
    <div v-if="showRefundModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-md w-full space-y-4 shadow-xl animate-in fade-in">
        <div>
          <h3 class="text-base font-extrabold text-zinc-900">Pengajuan Refund Registrasi</h3>
          <p class="text-xs text-zinc-500 mt-0.5">Pengembalian dana registrasi akan diproses manual oleh Kasir GRO.</p>
        </div>

        <div class="space-y-3 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Bank Tujuan:</label>
            <select v-model="refundBank" class="w-full px-3 py-2 rounded-xl border border-zinc-300 text-xs font-bold">
              <option value="BCA">BCA (Bank Central Asia)</option>
              <option value="MANDIRI">Bank Mandiri</option>
              <option value="BRI">BRI (Bank Rakyat Indonesia)</option>
              <option value="BNI">BNI (Bank Negara Indonesia)</option>
            </select>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Nomor Rekening:</label>
            <input v-model="refundAccNumber" type="text" placeholder="Contoh: 1234567890" class="w-full px-3 py-2 rounded-xl border border-zinc-300 font-mono text-xs font-bold" />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Nama Pemilik Rekening:</label>
            <input v-model="refundAccName" type="text" placeholder="Sesuai buku tabungan" class="w-full px-3 py-2 rounded-xl border border-zinc-300 text-xs font-bold" />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Alasan Pembatalan:</label>
            <textarea v-model="refundNotes" rows="2" placeholder="Contoh: Pasien berhalangan hadir..." class="w-full px-3 py-2 rounded-xl border border-zinc-300 text-xs"></textarea>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            @click="showRefundModal = false"
            class="flex-1 py-2.5 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-xs font-bold text-zinc-700"
          >
            Batal
          </button>
          <button
            type="button"
            @click="handleRequestRefund"
            :disabled="isSubmittingRefund"
            class="flex-1 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-xs font-bold text-white flex items-center justify-center gap-1.5"
          >
            <Loader2 v-if="isSubmittingRefund" class="w-3.5 h-3.5 animate-spin" />
            <span>Kirim Pengajuan</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL LACAK TIKET MANDIRI -->
    <div v-if="showTrackModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-md w-full space-y-4 shadow-xl animate-in fade-in">
        <div>
          <h3 class="text-base font-extrabold text-zinc-900">Lacak Tiket Antrean</h3>
          <p class="text-xs text-zinc-500 mt-0.5">Masukkan Nomor WhatsApp atau Nomor Antrean untuk melihat tiket Anda.</p>
        </div>

        <div class="space-y-2">
          <input
            v-model="trackQuery"
            type="text"
            placeholder="Nomor WhatsApp (misal 081234567890)"
            @keyup.enter="handleTrackTicket"
            class="w-full px-3.5 py-2.5 rounded-xl border border-zinc-300 text-xs font-bold outline-none"
          />
          <div v-if="trackError" class="text-xs text-rose-600 font-semibold">
            {{ trackError }}
          </div>
        </div>

        <div class="flex items-center gap-2 pt-1">
          <button
            type="button"
            @click="showTrackModal = false"
            class="flex-1 py-2.5 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-xs font-bold text-zinc-700"
          >
            Tutup
          </button>
          <button
            type="button"
            @click="handleTrackTicket"
            :disabled="isTracking"
            class="flex-1 py-2.5 rounded-xl bg-[#0C3B36] hover:bg-[#134E48] text-xs font-bold text-white flex items-center justify-center gap-1.5"
          >
            <Loader2 v-if="isTracking" class="w-3.5 h-3.5 animate-spin" />
            <span>Cari Tiket</span>
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.date-chips-container::-webkit-scrollbar {
  height: 4px;
}
.date-chips-container::-webkit-scrollbar-track {
  background: #f1ede6;
  border-radius: 9999px;
}
.date-chips-container::-webkit-scrollbar-thumb {
  background: #0C3B36;
  opacity: 0.6;
  border-radius: 9999px;
}
.date-chips-container::-webkit-scrollbar-thumb:hover {
  background: #E5A93C;
}
</style>
