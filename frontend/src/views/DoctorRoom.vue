<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import {
  api,
  type DoctorQueuePatient,
  type DoctorCatalogDrug,
  type DoctorConsultation,
} from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useClinicStore } from '../stores/clinic'
import {
  Stethoscope,
  Send,
  Plus,
  Trash2,
  CheckCircle2,
  AlertCircle,
  RefreshCw,
  Clock,
  Volume2,
  FileText,
  Pill,
  Sparkles,
  Search,
  User,
  ShieldCheck,
  Building2,
  ChevronRight,
  Activity,
  Calendar,
  CalendarOff,
  Users,
  Clock3,
  Ban,
  Pencil,
  Loader2,
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

// Doctor Session State
const availableDoctors = [
  { id: 2, name: 'dr. Patricia Vanessa A., Sp.A', unit: 'KIDINEST', title: 'Dokter Spesialis Anak' },
  { id: 3, name: 'dr. Hendra Pratama, Sp.A', unit: 'KIDINEST', title: 'Konsultan Tumbang & Nutrisi' },
  { id: 4, name: 'dr. Dian Paramita, Sp.D.V.E', unit: 'SSC', title: 'Spesialis Dermatologi & Laser' },
  { id: 5, name: 'dr. Kevin Wijaya, Sp.D.V.E', unit: 'SSC', title: 'Spesialis Kulit & Bedah Minor' },
]

function getInitialDoctor() {
  if (authStore.user?.role === 'DOCTOR' && authStore.user.name) {
    const userUnit = authStore.user.clinic_unit || 'KIDINEST'
    const userName = authStore.user.name.toLowerCase()
    const found = availableDoctors.find((d) => {
      const docName = d.name.toLowerCase()
      const parts = docName.replace('dr.', '').trim().split(' ')
      return parts.some((p) => p.length > 3 && userName.includes(p.toLowerCase()))
    })
    if (found) return found
    return {
      id: userUnit === 'KIDINEST' ? 2 : 4,
      name: authStore.user.name,
      unit: userUnit,
      title: userUnit === 'KIDINEST' ? 'Dokter Spesialis Anak' : 'Dokter Spesialis Kulit & Estetika',
    }
  }
  return availableDoctors[0]
}

const currentDoctor = ref(getInitialDoctor())

// Tab: 'EXAMINATION' (Pemeriksaan Hari Ini) | 'SCHEDULE' (Jadwal Praktek & Cuti)
const activeDoctorTab = ref<'EXAMINATION' | 'SCHEDULE'>('EXAMINATION')

// Schedule & Leave State
interface WeeklyShift {
  id?: number
  day_of_week: number
  shift_name: string
  start_time: string
  end_time: string
  slot_minutes: number
  quota_per_slot: number
  is_active: boolean
}

interface DoctorLeave {
  id: number
  doctor_id: number
  override_date: string
  is_cancelled: boolean
  reason: string
}

const weeklySchedules = ref<WeeklyShift[]>([])
const doctorLeaves = ref<DoctorLeave[]>([])
const futureReservations = ref<any[]>([])
const loadingSchedules = ref(false)

const newLeaveDate = ref('')
const newLeaveReason = ref('')
const isSubmittingLeave = ref(false)
const leaveActionMsg = ref('')

const dayNamesMap: Record<number, string> = {
  0: 'Minggu',
  1: 'Senin',
  2: 'Selasa',
  3: 'Rabu',
  4: 'Kamis',
  5: 'Jumat',
  6: 'Sabtu',
}

const editingShift = ref<WeeklyShift | null>(null)
const isSavingShift = ref(false)
const shiftActionMsg = ref('')

function startEditShift(shift: WeeklyShift) {
  editingShift.value = { ...shift }
  shiftActionMsg.value = ''
}

async function handleSaveShift() {
  if (!editingShift.value || !currentDoctor.value.id) return
  isSavingShift.value = true
  shiftActionMsg.value = ''
  try {
    await api.updateDoctorSchedule({
      id: editingShift.value.id,
      doctor_id: currentDoctor.value.id,
      doctor_name: currentDoctor.value.name,
      clinic_unit: currentDoctor.value.unit,
      day_of_week: editingShift.value.day_of_week,
      shift_name: editingShift.value.shift_name || 'PAGI',
      start_time: editingShift.value.start_time,
      end_time: editingShift.value.end_time,
      slot_minutes: editingShift.value.slot_minutes || 30,
      quota_per_slot: editingShift.value.quota_per_slot || 1,
      is_active: editingShift.value.is_active,
    })
    shiftActionMsg.value = `Shift ${dayNamesMap[editingShift.value.day_of_week]} berhasil diperbarui!`
    editingShift.value = null
    await loadDoctorScheduleData()
  } catch (err: any) {
    shiftActionMsg.value = err.message || 'Gagal menyimpan perubahan shift'
  } finally {
    isSavingShift.value = false
  }
}

async function loadDoctorScheduleData() {
  if (!currentDoctor.value.id) return
  loadingSchedules.value = true
  leaveActionMsg.value = ''
  try {
    const [schedRes, leaveRes, resvRes] = await Promise.all([
      api.getDoctorWeeklySchedules(currentDoctor.value.id),
      api.getDoctorLeaves(currentDoctor.value.id),
      api.getReservations(currentDoctor.value.unit),
    ])
    weeklySchedules.value = schedRes || []
    doctorLeaves.value = leaveRes || []
    futureReservations.value = (resvRes || []).filter((r: any) =>
      !r.doctor_name || r.doctor_name.toLowerCase().includes(currentDoctor.value.name.toLowerCase().split(' ')[1] || '')
    )
  } catch (err) {
    console.error('Gagal memuat jadwal & cuti dokter:', err)
  } finally {
    loadingSchedules.value = false
  }
}

async function handleAddLeave() {
  if (!currentDoctor.value.id || !newLeaveDate.value) return
  isSubmittingLeave.value = true
  leaveActionMsg.value = ''
  try {
    await api.addDoctorLeave({
      doctor_id: currentDoctor.value.id,
      date: newLeaveDate.value,
      reason: newLeaveReason.value || 'Cuti Dokter',
    })
    leaveActionMsg.value = `Jadwal praktek tanggal ${newLeaveDate.value} berhasil ditutup!`
    newLeaveDate.value = ''
    newLeaveReason.value = ''
    await loadDoctorScheduleData()
  } catch (err: any) {
    leaveActionMsg.value = err.message || 'Gagal menyimpan cuti'
  } finally {
    isSubmittingLeave.value = false
  }
}

async function handleCancelLeave(date: string) {
  if (!currentDoctor.value.id || !date) return
  try {
    await api.cancelDoctorLeave({
      doctor_id: currentDoctor.value.id,
      date: date,
    })
    await loadDoctorScheduleData()
  } catch (err: any) {
    alert(err.message || 'Gagal membatalkan cuti')
  }
}

// State: Queues & Active Patient
const queues = ref<DoctorQueuePatient[]>([])
const activeQueue = ref<DoctorQueuePatient | null>(null)
const loadingQueues = ref(false)

// State: Catalog Drugs
const catalogDrugs = ref<DoctorCatalogDrug[]>([])
const loadingDrugs = ref(false)

// State: Consultation Form
const anamnesis = ref('')
const diagnosis = ref('')
const icd10Code = ref('J06.9')
const actionNotes = ref('')
const actionFee = ref<number>(75000)

interface RxRow {
  is_manual: boolean
  drug_id?: number
  item_name: string
  quantity: number
  unit_price: number
  dosage_instruction: string
}

const prescriptionRows = ref<RxRow[]>([])

const isSubmitting = ref(false)
const submitSuccessMsg = ref('')
const submitErrorMsg = ref('')

// Presets ICD-10
const countWaiting = computed(() => (queues.value || []).filter((q) => q.queue_status === 'WAITING').length)
const countCalling = computed(() => (queues.value || []).filter((q) => q.queue_status === 'CALLING').length)
const countInExam = computed(() => (queues.value || []).filter((q) => q.queue_status === 'IN_EXAM').length)
const queueFilter = ref<'ALL' | 'WAITING' | 'CALLING' | 'IN_EXAM'>('ALL')

const filteredQueues = computed(() => {
  const qList = queues.value || []
  if (queueFilter.value === 'ALL') return qList
  return qList.filter((q) => q.queue_status === queueFilter.value)
})

const icd10PresetsKidinest = [
  { code: 'J06.9', name: 'ISPA / Batuk Pilek Akut' },
  { code: 'R50.9', name: 'Febris / Demam Tanpa Sebab' },
  { code: 'A09', name: 'Diare Akut & Gastroenteritis' },
  { code: 'Z23', name: 'Imunisasi / Vaksinasi Rutin Anak' },
  { code: 'Z00.1', name: 'Pemeriksaan Rutin Tumbuh Kembang' },
]

const icd10PresetsSSC = [
  { code: 'L70.0', name: 'Acne Vulgaris / Jerawat Meradang' },
  { code: 'L20.9', name: 'Dermatitis Atopik / Eksem Kulit' },
  { code: 'L81.1', name: 'Melasma / Flek Hitam Wajah' },
  { code: 'L29.9', name: 'Pruritus / Gatal Alergi Kulit' },
  { code: 'Z41.1', name: 'Tindakan Facial Estetika / Peeling' },
]

const currentPresets = computed(() => {
  return currentDoctor.value.unit === 'SSC' ? icd10PresetsSSC : icd10PresetsKidinest
})

function formatRupiah(val?: number) {
  return 'Rp ' + (val || 0).toLocaleString('id-ID')
}

// Total prescription estimate
const totalRxCost = computed(() => {
  return prescriptionRows.value.reduce((acc, row) => acc + (row.quantity || 0) * (row.unit_price || 0), 0)
})

// Fetch Queues & Drugs
async function loadDoctorQueues() {
  loadingQueues.value = true
  try {
    const res = await api.getDoctorQueues(currentDoctor.value.name, currentDoctor.value.unit)
    queues.value = Array.isArray(res) ? res : []
    if (!activeQueue.value && queues.value.length > 0) {
      selectPatient(queues.value[0])
    }
  } catch (err: any) {
    console.error('Gagal memuat antrian dokter:', err)
    queues.value = []
  } finally {
    loadingQueues.value = false
  }
}

async function loadDrugs() {
  loadingDrugs.value = true
  try {
    const res = await api.getDoctorDrugs(currentDoctor.value.unit)
    catalogDrugs.value = Array.isArray(res) ? res : []
  } catch (err: any) {
    console.error('Gagal memuat katalog obat:', err)
    catalogDrugs.value = []
  } finally {
    loadingDrugs.value = false
  }
}

function switchDoctor(doc: (typeof availableDoctors)[0]) {
  currentDoctor.value = doc
  activeQueue.value = null
  prescriptionRows.value = []
  if (doc.unit === 'SSC') {
    actionFee.value = 150000
    icd10Code.value = 'L70.0'
  } else {
    actionFee.value = 75000
    icd10Code.value = 'J06.9'
  }
  loadDoctorQueues()
  loadDrugs()
  loadDoctorScheduleData()
}

// Select Patient
async function selectPatient(p: DoctorQueuePatient) {
  activeQueue.value = p
  submitSuccessMsg.value = ''
  submitErrorMsg.value = ''

  // Defaults per unit
  if (p.clinic_unit === 'SSC') {
    actionFee.value = 150000
    anamnesis.value = 'Pasien datang untuk perawatan jerawat & peremajaan kulit wajah.'
    diagnosis.value = 'Acne Vulgaris Grade 2'
    icd10Code.value = 'L70.0'
    actionNotes.value = 'Prosedur Deep Cleansing Acne + Chemical Peeling Glycolic 20%'
  } else {
    actionFee.value = 75000
    anamnesis.value = 'Pasien anak datang dengan keluhan batuk berdahak dan demam 2 hari.'
    diagnosis.value = 'ISPA / Infeksi Saluran Pernapasan Akut'
    icd10Code.value = 'J06.9'
    actionNotes.value = 'Pemeriksaan fisik toraks & edukasi kompres hangat'
  }

  // Check if existing consultation exists
  try {
    const existing = await api.getDoctorConsultationByQueue(p.queue_id)
    if (existing) {
      anamnesis.value = existing.anamnesis || anamnesis.value
      diagnosis.value = existing.diagnosis || diagnosis.value
      icd10Code.value = existing.icd10_code || icd10Code.value
      actionNotes.value = existing.action_notes || actionNotes.value
      actionFee.value = existing.action_fee || actionFee.value
      if (existing.prescription && existing.prescription.items.length > 0) {
        prescriptionRows.value = existing.prescription.items.map((it) => ({
          is_manual: it.is_manual,
          drug_id: it.drug_id,
          item_name: it.item_name,
          quantity: it.quantity,
          unit_price: it.unit_price,
          dosage_instruction: it.dosage_instruction,
        }))
        return
      }
    }
  } catch {
    // Ignore if not created yet
  }

  // Prepopulate standard prescription row
  if (prescriptionRows.value.length === 0 && catalogDrugs.value.length > 0) {
    const firstDrug = catalogDrugs.value[0]
    prescriptionRows.value = [
      {
        is_manual: false,
        drug_id: firstDrug.id,
        item_name: firstDrug.name,
        quantity: 1,
        unit_price: firstDrug.sell_price,
        dosage_instruction: p.clinic_unit === 'SSC' ? 'Oles tipis merata pagi & malam hari' : '3x 0.6 ml sehari sesudah makan',
      },
    ]
  }
}

// Call Patient
async function callPatient(p: DoctorQueuePatient) {
  try {
    await api.callDoctorPatient(p.queue_id)
    p.queue_status = 'CALLING'
    clinicStore.playChime()
    await loadDoctorQueues()
  } catch (err: any) {
    alert(err.message || 'Gagal memanggil pasien')
  }
}

// Start Examination
async function startExam(p: DoctorQueuePatient) {
  try {
    await api.startDoctorExam(p.queue_id)
    p.queue_status = 'IN_EXAM'
    await loadDoctorQueues()
  } catch (err: any) {
    alert(err.message || 'Gagal memulai pemeriksaan')
  }
}

// Add Rows to Prescription
function addCatalogDrugRow() {
  if (catalogDrugs.value.length === 0) return
  const d = catalogDrugs.value[0]
  prescriptionRows.value.push({
    is_manual: false,
    drug_id: d.id,
    item_name: d.name,
    quantity: 1,
    unit_price: d.sell_price,
    dosage_instruction: '3x sehari sesudah makan',
  })
}

function addManualRacikanRow() {
  prescriptionRows.value.push({
    is_manual: true,
    item_name: currentDoctor.value.unit === 'SSC' ? 'Krim Racikan Salep Acne Custom 10g' : 'Puyer Racikan Demam Batuk (10 Bungkus)',
    quantity: 1,
    unit_price: currentDoctor.value.unit === 'SSC' ? 45000 : 35000,
    dosage_instruction: currentDoctor.value.unit === 'SSC' ? 'Oles tipis malam hari sebelum tidur' : '3x 1 bungkus sesudah makan',
  })
}

function onDrugSelected(row: RxRow, drugId: number) {
  const match = catalogDrugs.value.find((d) => d.id === drugId)
  if (match) {
    row.drug_id = match.id
    row.item_name = match.name
    row.unit_price = match.sell_price
  }
}

function removePrescriptionRow(index: number) {
  prescriptionRows.value.splice(index, 1)
}

// Submit Consultation & Send Rx
async function submitConsultation() {
  if (!activeQueue.value) return
  if (!diagnosis.value) {
    alert('Diagnosa dokter wajib diisi!')
    return
  }

  isSubmitting.value = true
  submitSuccessMsg.value = ''
  submitErrorMsg.value = ''

  try {
    await api.submitDoctorConsultation({
      queue_id: activeQueue.value.queue_id,
      doctor_name: currentDoctor.value.name,
      anamnesis: anamnesis.value,
      diagnosis: diagnosis.value,
      icd10_code: icd10Code.value,
      action_notes: actionNotes.value,
      action_fee: actionFee.value,
      prescription_items: prescriptionRows.value.map((r) => ({
        drug_id: r.drug_id,
        item_name: r.item_name,
        is_manual: r.is_manual,
        quantity: r.quantity,
        unit_price: r.unit_price,
        dosage_instruction: r.dosage_instruction,
      })),
    })

    clinicStore.playChime()
    submitSuccessMsg.value = `Konsultasi selesai & resep untuk ${activeQueue.value.patient_name} berhasil dikirim ke Apotek.`

    // Refresh queues & drugs
    await Promise.all([loadDoctorQueues(), loadDrugs()])

    // Select next patient if available
    if (queues.value && queues.value.length > 0) {
      selectPatient(queues.value[0])
    } else {
      activeQueue.value = null
    }
  } catch (err: any) {
    submitErrorMsg.value = err.message || 'Gagal menyimpan konsultasi'
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  loadDoctorQueues()
  loadDrugs()
  loadDoctorScheduleData()
})

watch(activeDoctorTab, (tab) => {
  if (tab === 'SCHEDULE') {
    loadDoctorScheduleData()
  }
})
</script>

<template>
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- TOP BAR: MINIMALIST CLINICAL HEADER -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Ruang Periksa Dokter
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Clinical Examination & E-Prescription
        </h1>
      </div>

      <div class="flex items-center gap-3">
        <!-- Quick Doctor Switcher: Demo mode for Superadmin / Dev -->
        <div v-if="authStore.isOwner || !authStore.user" class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
          <button
            v-for="doc in availableDoctors"
            :key="doc.name"
            @click="switchDoctor(doc)"
            :class="[
              'px-2.5 py-1.5 rounded-lg transition-all cursor-pointer flex items-center gap-1.5',
              currentDoctor.name === doc.name
                ? 'bg-white text-zinc-950 shadow-2xs'
                : 'text-zinc-600 hover:text-zinc-950',
            ]"
          >
            <span
              :class="[
                'w-1.5 h-1.5 rounded-full',
                doc.unit === 'KIDINEST' ? 'bg-amber-500' : 'bg-rose-500',
              ]"
            ></span>
            <span class="truncate max-w-[110px]">{{ doc.name.split(',')[0] }}</span>
          </button>
        </div>

        <!-- Static Doctor Profile Badge (When logged in as Doctor) -->
        <div v-else class="text-right">
          <div class="text-[10px] uppercase font-bold text-zinc-400">Dokter Spesialis</div>
          <div class="text-xs font-bold text-zinc-900">{{ currentDoctor.name }}</div>
        </div>

        <button
          @click="loadDoctorQueues(); loadDrugs()"
          :disabled="loadingQueues || loadingDrugs"
          title="Segarkan antrian pasien dan katalog obat"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', (loadingQueues || loadingDrugs) && 'animate-spin']" />
          <span>Sync</span>
        </button>
      </div>
    </header>

    <!-- WORKSPACE TABS -->
    <div class="flex items-center gap-2 border-b border-zinc-200 pb-2">
      <button
        type="button"
        @click="activeDoctorTab = 'EXAMINATION'"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer"
        :class="activeDoctorTab === 'EXAMINATION' ? 'bg-[#0C3B36] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'"
      >
        <Stethoscope class="w-4 h-4" />
        <span>Ruang Konsultasi & Resep</span>
      </button>

      <button
        type="button"
        @click="activeDoctorTab = 'SCHEDULE'"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all cursor-pointer"
        :class="activeDoctorTab === 'SCHEDULE' ? 'bg-[#0C3B36] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'"
      >
        <Calendar class="w-4 h-4" />
        <span>Jadwal Praktik & Cuti Dokter</span>
        <span v-if="doctorLeaves.length > 0" class="px-1.5 py-0.5 rounded-full text-[10px] bg-rose-500 text-white font-mono">
          {{ doctorLeaves.length }} Cuti
        </span>
      </button>
    </div>

    <!-- TAB 1: EXAMINATION & E-PRESCRIPTION -->
    <div v-if="activeDoctorTab === 'EXAMINATION'" class="space-y-5">
      <!-- 1-ROW COMPACT OPERATIONAL TICKER STRIP (Hairline Grid) -->
      <div
        class="grid grid-cols-2 sm:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
      >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Antre Menunggu</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ countWaiting }} Pasien
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Dalam Panggilan</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-amber-700 mt-0.5">
          {{ countCalling }} Pasien
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Sedang Diperiksa</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-emerald-700 mt-0.5">
          {{ countInExam }} Pasien
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Katalog Obat Aktif</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-700 mt-0.5">
          {{ catalogDrugs.length }} Formula
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex items-center justify-between">
        <div>
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Estimasi Total Kasus</div>
          <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
            {{ formatRupiah((activeQueue ? actionFee : 0) + totalRxCost) }}
          </div>
        </div>
        <Activity class="w-4 h-4 text-zinc-400" />
      </div>
    </div>

    <!-- MAIN MASTER-DETAIL WORKSPACE (2 PANES) -->
    <main class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      
      <!-- LEFT PANE: FIFO PATIENT QUEUE (5 cols) -->
      <section class="lg:col-span-5 space-y-3">
        <!-- Control Bar: Filter Pills -->
        <div class="flex items-center justify-between gap-2">
          <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              type="button"
              @click="queueFilter = 'ALL'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                queueFilter === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Semua ({{ queues.length }})
            </button>
            <button
              type="button"
              @click="queueFilter = 'WAITING'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                queueFilter === 'WAITING' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Menunggu ({{ countWaiting }})
            </button>
            <button
              type="button"
              @click="queueFilter = 'CALLING'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                queueFilter === 'CALLING' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Dipanggil ({{ countCalling }})
            </button>
          </div>

          <span class="text-[11px] font-medium text-zinc-400">
            {{ currentDoctor.unit }}
          </span>
        </div>

        <!-- Patient List Container -->
        <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
          <div class="p-3.5 border-b border-zinc-100 bg-zinc-50/60 flex items-center justify-between">
            <div class="flex items-center gap-1.5">
              <span class="text-xs font-bold text-zinc-900">Antrean Pasien FIFO</span>
              <span class="text-[10px] font-bold px-1.5 py-0.2 rounded-md bg-zinc-200 text-zinc-800 tabular-nums">
                {{ filteredQueues.length }}
              </span>
            </div>
            <span class="text-[10px] text-zinc-400 font-medium">Urutan Kedatangan</span>
          </div>

          <div v-if="loadingQueues" class="py-12 text-center text-xs text-zinc-400">
            <RefreshCw class="w-4 h-4 mx-auto animate-spin mb-2 text-zinc-400" />
            Memuat daftar antrean...
          </div>

          <div v-else-if="filteredQueues.length === 0" class="py-12 text-center text-xs text-zinc-400">
            <CheckCircle2 class="w-6 h-6 mx-auto text-emerald-500 mb-1.5 opacity-80" />
            Tidak ada antrean dalam filter ini.
          </div>

          <div v-else class="divide-y divide-zinc-100 max-h-[640px] overflow-y-auto">
            <div
              v-for="p in filteredQueues"
              :key="p.queue_id"
              @click="selectPatient(p)"
              :class="[
                'p-3.5 transition-colors cursor-pointer space-y-2',
                activeQueue?.queue_id === p.queue_id
                  ? 'bg-zinc-900 text-white'
                  : 'hover:bg-zinc-50/70 text-zinc-900',
              ]"
            >
              <!-- Row Header: Queue No & Unit Badges -->
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-1.5">
                  <span
                    :class="[
                      'font-mono text-xs font-extrabold px-2 py-0.5 rounded-md tabular-nums',
                      activeQueue?.queue_id === p.queue_id
                        ? 'bg-zinc-800 text-white'
                        : 'bg-zinc-100 text-zinc-900'
                    ]"
                  >
                    {{ p.queue_number }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-semibold px-1.5 py-0.2 rounded border',
                      p.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                        : 'bg-rose-50 text-rose-900 border-rose-200/80',
                    ]"
                  >
                    {{ p.clinic_unit }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-medium px-1.5 py-0.2 rounded',
                      activeQueue?.queue_id === p.queue_id ? 'bg-zinc-800 text-zinc-300' : 'bg-zinc-100 text-zinc-600'
                    ]"
                  >
                    {{ p.category }}
                  </span>
                </div>

                <!-- Status Badge -->
                <span
                  :class="[
                    'text-[10px] font-bold px-2 py-0.5 rounded-full',
                    p.queue_status === 'CALLING'
                      ? 'bg-amber-400 text-zinc-950 animate-pulse'
                      : p.queue_status === 'IN_EXAM'
                      ? 'bg-emerald-400 text-zinc-950'
                      : activeQueue?.queue_id === p.queue_id
                      ? 'bg-zinc-800 text-zinc-300'
                      : 'bg-zinc-100 text-zinc-600',
                  ]"
                >
                  {{ p.queue_status }}
                </span>
              </div>

              <!-- Patient Info -->
              <div>
                <div class="font-bold text-sm tracking-tight leading-snug">{{ p.patient_name }}</div>
                <div
                  :class="[
                    'text-xs mt-0.5 flex items-center justify-between',
                    activeQueue?.queue_id === p.queue_id ? 'text-zinc-400' : 'text-zinc-500'
                  ]"
                >
                  <span>{{ p.doctor_name.split(',')[0] }}</span>
                  <span class="font-mono text-[11px]">{{ p.patient_phone || '-' }}</span>
                </div>
              </div>

              <!-- Action Buttons for Queue Row -->
              <div
                class="pt-2 border-t flex items-center justify-end gap-1.5"
                :class="activeQueue?.queue_id === p.queue_id ? 'border-zinc-800' : 'border-zinc-100'"
                @click.stop
              >
                <button
                  v-if="p.queue_status !== 'CALLING' && p.queue_status !== 'IN_EXAM'"
                  type="button"
                  @click="callPatient(p)"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg border border-amber-200 bg-amber-50 hover:bg-amber-100 text-amber-900 text-[11px] font-bold transition-all active:scale-95 cursor-pointer"
                >
                  <Volume2 class="w-3 h-3 text-amber-700" />
                  <span>Panggil</span>
                </button>

                <button
                  v-if="p.queue_status === 'CALLING'"
                  type="button"
                  @click="startExam(p)"
                  class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-[11px] font-bold transition-all active:scale-95 cursor-pointer shadow-2xs"
                >
                  <Stethoscope class="w-3 h-3" />
                  <span>Mulai Periksa</span>
                </button>

                <span v-if="p.queue_status === 'IN_EXAM'" class="text-[10px] font-bold text-emerald-400">
                  ● Sedang Diperiksa
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- RIGHT PANE: CLINICAL CONSULTATION & FAST E-PRESCRIPTION (7 cols, Sticky) -->
      <section class="lg:col-span-7 sticky top-6">
        
        <!-- Empty State -->
        <div
          v-if="!activeQueue"
          class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-12 text-center"
        >
          <Stethoscope class="w-10 h-10 mx-auto text-zinc-300 mb-3" />
          <h3 class="text-sm font-bold text-zinc-800">Pilih Pasien untuk Memulai Konsultasi</h3>
          <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
            Klik salah satu kartu antrean pada panel kiri untuk mengisi anamnesis, input diagnosa ICD-10, tindakan medis, dan peresepan obat terpadu.
          </p>
        </div>

        <!-- Active Patient Form -->
        <div v-else class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-5 sm:p-6 space-y-5">
          
          <!-- Patient Top Card -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-zinc-100">
            <div>
              <div class="flex items-center gap-2">
                <span class="px-2.5 py-0.5 rounded-md font-mono text-xs font-black bg-zinc-950 text-white">
                  {{ activeQueue.queue_number }}
                </span>
                <span class="text-xs font-semibold px-2 py-0.5 rounded border border-zinc-200 bg-zinc-50 text-zinc-800">
                  {{ activeQueue.clinic_unit }}
                </span>
                <span class="text-xs text-zinc-400 font-medium">
                  {{ activeQueue.category }}
                </span>
              </div>
              <h2 class="text-xl font-bold text-zinc-950 tracking-tight mt-1.5">
                {{ activeQueue.patient_name }}
              </h2>
              <div class="text-xs text-zinc-500 mt-0.5">
                Pemeriksa: {{ currentDoctor.name }}
              </div>
            </div>

            <!-- Calling & Exam Controls -->
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="callPatient(activeQueue)"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 shadow-2xs transition-all active:scale-95 cursor-pointer"
                title="Bunyikan panggilan suara di layar monitor ruang tunggu"
              >
                <Volume2 class="w-3.5 h-3.5 text-zinc-600" />
                <span>{{ activeQueue.queue_status === 'CALLING' ? 'Panggil Ulang' : 'Panggil' }}</span>
              </button>

              <button
                v-if="activeQueue.queue_status !== 'IN_EXAM'"
                type="button"
                @click="startExam(activeQueue)"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
              >
                <Stethoscope class="w-3.5 h-3.5 text-emerald-400" />
                <span>Mulai Periksa</span>
              </button>

              <span
                v-else
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-emerald-50 text-emerald-900 border border-emerald-200 text-xs font-bold"
              >
                <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                <span>Sedang Diperiksa</span>
              </span>
            </div>
          </div>

          <!-- Section 1: Anamnesis / Keluhan -->
          <div class="space-y-1.5">
            <label class="text-[11px] font-bold uppercase tracking-wider text-zinc-500">
              Anamnesis & Riwayat Keluhan Klinis
            </label>
            <textarea
              v-model="anamnesis"
              rows="2"
              placeholder="Ketik keluhan pasien, durasi gejala, riwayat alergi obat..."
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 text-xs text-zinc-900 focus:bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950 transition-all placeholder:text-zinc-400"
            ></textarea>
          </div>

          <!-- Section 2: Diagnosa & ICD-10 -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-[11px] font-bold uppercase tracking-wider text-zinc-500">
                Diagnosa Klinis & ICD-10 Kemenkes
              </label>
              <span class="text-[10px] text-zinc-400">Standar SATUSEHAT</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
              <input
                v-model="diagnosis"
                type="text"
                placeholder="Diagnosa utama dokter"
                class="sm:col-span-2 px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 text-xs font-bold text-zinc-900 focus:bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
              <input
                v-model="icd10Code"
                type="text"
                placeholder="Kode ICD-10"
                class="px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 text-xs font-mono font-bold text-center text-zinc-900 focus:bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <!-- Quick ICD-10 Chips -->
            <div class="flex items-center gap-1.5 flex-wrap pt-0.5">
              <span class="text-[10px] font-bold text-zinc-400">Pilihan Cepat:</span>
              <button
                v-for="chip in currentPresets"
                :key="chip.code"
                type="button"
                @click="diagnosis = chip.name; icd10Code = chip.code"
                class="px-2 py-0.5 rounded-md text-[10px] font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 transition-all cursor-pointer"
              >
                {{ chip.code }} - {{ chip.name.split('/')[0] }}
              </button>
            </div>
          </div>

          <!-- Section 3: Tindakan Medis & Jasa Dokter -->
          <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 space-y-2.5">
            <div class="flex items-center justify-between">
              <label class="text-[11px] font-bold uppercase tracking-wider text-zinc-700">
                Tindakan Medis & Biaya Jasa
              </label>
              <span class="text-[10px] text-zinc-400">Diteruskan ke Invoice Kasir</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
              <div class="sm:col-span-2 space-y-1">
                <label class="text-[10px] font-bold text-zinc-500">Uraian Tindakan / Prosedur:</label>
                <input
                  v-model="actionNotes"
                  type="text"
                  placeholder="Contoh: Injeksi Vaksin DTP / Peeling Acne"
                  class="w-full px-3 py-1.5 rounded-lg bg-white border border-zinc-200 text-xs text-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-950"
                />
              </div>

              <div class="space-y-1">
                <label class="text-[10px] font-bold text-zinc-500">Tarif Tindakan (Rp):</label>
                <input
                  v-model.number="actionFee"
                  type="number"
                  min="0"
                  class="w-full px-3 py-1.5 rounded-lg bg-white border border-zinc-200 text-xs font-mono font-bold text-right text-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-950"
                />
              </div>
            </div>
          </div>

          <!-- Section 4: Fast E-Prescription Table -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-xs font-bold uppercase tracking-wider text-zinc-950 flex items-center gap-1.5">
                  <Pill class="w-3.5 h-3.5 text-zinc-700" />
                  <span>Tabel Peresepan Elektronik (E-Resep)</span>
                </h3>
                <p class="text-[10px] text-zinc-400">Obat katalog memotong stok otomatis saat diserahkan apotek</p>
              </div>

              <div class="flex items-center gap-2">
                <button
                  type="button"
                  @click="addCatalogDrugRow"
                  class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800 text-xs font-bold transition-all shadow-2xs cursor-pointer"
                >
                  <Plus class="w-3.5 h-3.5" />
                  <span>Obat Katalog</span>
                </button>

                <button
                  type="button"
                  @click="addManualRacikanRow"
                  class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800 text-xs font-bold transition-all shadow-2xs cursor-pointer"
                >
                  <Sparkles class="w-3.5 h-3.5 text-amber-600" />
                  <span>Racikan Manual</span>
                </button>
              </div>
            </div>

            <!-- Rx Rows -->
            <div v-if="prescriptionRows.length === 0" class="py-6 text-center text-xs text-zinc-400 bg-zinc-50 rounded-xl border border-dashed border-zinc-200">
              Belum ada obat yang diresepkan. Klik tombol di atas untuk menambah item.
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="(row, idx) in prescriptionRows"
                :key="idx"
                class="p-3 rounded-xl border border-zinc-200/90 bg-zinc-50/50 space-y-2"
              >
                <!-- Row Line 1: Type Badge, Name & Delete -->
                <div class="flex items-center justify-between gap-2">
                  <div class="flex items-center gap-2 flex-1">
                    <span
                      :class="[
                        'px-1.5 py-0.2 rounded text-[10px] font-bold shrink-0',
                        row.is_manual ? 'bg-amber-100 text-amber-900' : 'bg-zinc-200 text-zinc-800',
                      ]"
                    >
                      {{ row.is_manual ? 'Racik' : 'Katalog' }}
                    </span>

                    <div v-if="!row.is_manual" class="flex-1">
                      <select
                        :value="row.drug_id"
                        @change="(e: any) => onDrugSelected(row, parseInt(e.target.value))"
                        class="w-full px-2.5 py-1 rounded-lg bg-white border border-zinc-200 text-xs font-bold text-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-950"
                      >
                        <option v-for="d in catalogDrugs" :key="d.id" :value="d.id">
                          {{ d.name }} (Stok: {{ d.available_stock }} {{ d.unit }}) - {{ formatRupiah(d.sell_price) }}
                        </option>
                      </select>
                    </div>

                    <div v-else class="flex-1">
                      <input
                        v-model="row.item_name"
                        type="text"
                        placeholder="Nama obat racikan puyer / sirup / salep..."
                        class="w-full px-2.5 py-1 rounded-lg bg-white border border-zinc-200 text-xs font-bold text-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-950"
                      />
                    </div>
                  </div>

                  <button
                    type="button"
                    @click="removePrescriptionRow(idx)"
                    class="p-1 rounded-lg text-zinc-400 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>

                <!-- Row Line 2: Quantity, Price, Dosage & Subtotal -->
                <div class="grid grid-cols-1 sm:grid-cols-12 gap-2 items-center text-xs">
                  <div class="sm:col-span-2">
                    <label class="text-[10px] text-zinc-500 font-bold">Qty:</label>
                    <input
                      v-model.number="row.quantity"
                      type="number"
                      min="1"
                      class="w-full px-2 py-1 text-center font-mono font-bold rounded-lg bg-white border border-zinc-200 text-xs text-zinc-900"
                    />
                  </div>

                  <div class="sm:col-span-3">
                    <label class="text-[10px] text-zinc-500 font-bold">Harga Satuan:</label>
                    <input
                      v-model.number="row.unit_price"
                      type="number"
                      min="0"
                      :disabled="!row.is_manual"
                      class="w-full px-2 py-1 text-right font-mono font-bold rounded-lg bg-white border border-zinc-200 text-xs text-zinc-900 disabled:bg-zinc-100"
                    />
                  </div>

                  <div class="sm:col-span-5">
                    <label class="text-[10px] text-zinc-500 font-bold">Aturan Pakai (Signa):</label>
                    <input
                      v-model="row.dosage_instruction"
                      type="text"
                      placeholder="Contoh: 3x1 sendok takar sesudah makan"
                      class="w-full px-2 py-1 rounded-lg bg-white border border-zinc-200 text-xs text-zinc-900 placeholder:text-zinc-400"
                    />
                  </div>

                  <div class="sm:col-span-2 text-right">
                    <div class="text-[10px] text-zinc-500 font-bold">Subtotal:</div>
                    <div class="font-bold tabular-nums text-zinc-900">
                      {{ formatRupiah(row.quantity * row.unit_price) }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Rx Subtotal Strip -->
              <div class="p-3 rounded-xl border border-zinc-200 bg-stone-50/80 flex items-center justify-between text-xs">
                <span class="font-bold text-zinc-700">Estimasi Biaya Obat Apotek:</span>
                <span class="font-extrabold tabular-nums text-sm text-zinc-950">{{ formatRupiah(totalRxCost) }}</span>
              </div>
            </div>
          </div>

          <!-- Feedback Messages -->
          <div v-if="submitSuccessMsg" class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-950 flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
            <span>{{ submitSuccessMsg }}</span>
          </div>

          <div v-if="submitErrorMsg" class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-xs text-rose-950 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ submitErrorMsg }}</span>
          </div>

          <!-- Final Submit Button -->
          <div class="pt-2">
            <button
              type="button"
              @click="submitConsultation"
              :disabled="isSubmitting"
              class="w-full py-3.5 px-5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2 cursor-pointer"
            >
              <Send class="w-3.5 h-3.5 text-emerald-400" />
              <span>{{ isSubmitting ? 'Mengirim E-Resep...' : 'Kirim E-Resep ke Apotek & Selesai Periksa' }}</span>
            </button>
          </div>
        </div>
      </section>
    </main>
    </div>

    <!-- TAB 2: SCHEDULE & LEAVE MANAGEMENT -->
    <div v-else-if="activeDoctorTab === 'SCHEDULE'" class="space-y-6">
      
      <!-- Doctor Summary Header Banner -->
      <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">
            Pengaturan Jadwal Dokter Mandiri
          </span>
          <h2 class="text-xl font-extrabold text-zinc-950 mt-0.5">
            {{ currentDoctor.name }}
          </h2>
          <p class="text-xs text-zinc-500 font-medium mt-1">
            Unit: <strong>{{ currentDoctor.unit === 'KIDINEST' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</strong> • 
            Status: <span class="text-emerald-700 font-bold">Praktik Aktif</span>
          </p>
        </div>

        <button
          @click="loadDoctorScheduleData"
          :disabled="loadingSchedules"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl border border-zinc-200 bg-zinc-50 hover:bg-zinc-100 text-xs font-bold text-zinc-700 transition-all cursor-pointer self-start sm:self-auto"
        >
          <RefreshCw :class="['w-3.5 h-3.5', loadingSchedules && 'animate-spin']" />
          <span>Refresh Jadwal</span>
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        
        <!-- LEFT COLUMN: WEEKLY SHIFTS & LEAVE FORM (7 cols) -->
        <div class="lg:col-span-7 space-y-6">
          
          <!-- 1. Jadwal Shift Mingguan -->
          <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <Clock3 class="w-4 h-4 text-zinc-600" />
                <h3 class="text-sm font-extrabold text-zinc-950">Shift Praktik Mingguan</h3>
              </div>
              <span class="text-[11px] font-bold text-zinc-400 font-mono">Slot 30 Menit</span>
            </div>

            <div v-if="shiftActionMsg" class="p-3 rounded-xl text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200 flex items-center gap-2">
              <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
              <span>{{ shiftActionMsg }}</span>
            </div>

            <div v-if="loadingSchedules" class="py-8 text-center text-xs text-zinc-400 flex items-center justify-center gap-2">
              <Loader2 class="w-4 h-4 animate-spin" />
              <span>Memuat jadwal dokter...</span>
            </div>

            <div v-else-if="weeklySchedules.length === 0" class="py-6 text-center text-xs text-zinc-400 bg-zinc-50 rounded-xl border border-dashed border-zinc-200">
              Belum ada konfigurasi shift mingguan terdaftar.
            </div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <div
                v-for="shift in weeklySchedules"
                :key="shift.id || shift.day_of_week"
                class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/60 space-y-2 transition-all hover:bg-zinc-50"
              >
                <div class="flex items-center justify-between">
                  <span class="text-xs font-black text-zinc-900">{{ dayNamesMap[shift.day_of_week] || 'Hari' }}</span>
                  <div class="flex items-center gap-1.5">
                    <span
                      class="px-2 py-0.5 rounded text-[10px] font-bold uppercase"
                      :class="shift.is_active ? 'bg-emerald-100 text-emerald-900' : 'bg-zinc-200 text-zinc-600'"
                    >
                      {{ shift.is_active ? shift.shift_name : 'Libur' }}
                    </span>
                    <button
                      type="button"
                      @click="startEditShift(shift)"
                      class="inline-flex items-center gap-1 px-2 py-0.5 rounded-lg text-[10px] font-bold bg-white border border-zinc-200 text-zinc-700 hover:bg-zinc-100 hover:text-zinc-950 transition-colors shadow-2xs cursor-pointer"
                      title="Ubah Jam Praktik"
                    >
                      <Pencil class="w-2.5 h-2.5 text-zinc-500" />
                      <span>Ubah</span>
                    </button>
                  </div>
                </div>
                <div class="text-xs font-mono font-bold" :class="shift.is_active ? 'text-[#0C3B36]' : 'text-zinc-400 line-through'">
                  {{ shift.start_time }} - {{ shift.end_time }} WIB
                </div>
                <div class="flex items-center justify-between text-[11px] text-zinc-500 pt-1 border-t border-zinc-200/60 font-medium">
                  <span>Durasi: {{ shift.slot_minutes }} mnt</span>
                  <span>Kuota: {{ shift.quota_per_slot }} pasien/slot</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 2. Form Ajukan Cuti Dokter / Tutup Jadwal -->
          <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs space-y-4">
            <div class="flex items-center gap-2">
              <CalendarOff class="w-4 h-4 text-rose-600" />
              <h3 class="text-sm font-extrabold text-zinc-950">Ajukan Cuti / Tutup Jadwal Praktik</h3>
            </div>
            <p class="text-xs text-zinc-500 leading-relaxed">
              Menutup jadwal pada tanggal tertentu akan otomatis menonaktifkan slot booking bagi pasien di Patient Portal.
            </p>

            <form @submit.prevent="handleAddLeave" class="space-y-3.5 pt-1">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
                    Tanggal Cuti / Libur:
                  </label>
                  <input
                    type="date"
                    v-model="newLeaveDate"
                    required
                    class="w-full px-3 py-2.5 rounded-xl border border-zinc-200 text-xs font-bold outline-none focus:border-zinc-950 bg-white"
                  />
                </div>

                <div>
                  <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
                    Alasan Cuti:
                  </label>
                  <input
                    type="text"
                    v-model="newLeaveReason"
                    placeholder="Contoh: Seminar IDAI / Izin Cuti"
                    required
                    class="w-full px-3 py-2.5 rounded-xl border border-zinc-200 text-xs font-bold outline-none focus:border-zinc-950 bg-white"
                  />
                </div>
              </div>

              <div v-if="leaveActionMsg" class="p-3 rounded-xl bg-amber-50 border border-amber-200 text-amber-950 text-xs font-semibold flex items-center gap-2">
                <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
                <span>{{ leaveActionMsg }}</span>
              </div>

              <button
                type="submit"
                :disabled="isSubmittingLeave || !newLeaveDate"
                class="py-3 px-5 rounded-xl bg-rose-700 hover:bg-rose-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2 cursor-pointer"
              >
                <Loader2 v-if="isSubmittingLeave" class="w-3.5 h-3.5 animate-spin" />
                <Ban v-else class="w-3.5 h-3.5" />
                <span>Tutup Jadwal Tanggal Ini</span>
              </button>
            </form>
          </div>

        </div>

        <!-- RIGHT COLUMN: LEAVE HISTORY & UPCOMING RESERVATIONS (5 cols) -->
        <div class="lg:col-span-5 space-y-6">
          
          <!-- Riwayat Cuti Aktif -->
          <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <Ban class="w-4 h-4 text-rose-600" />
                <h3 class="text-sm font-extrabold text-zinc-950">Daftar Cuti / Libur Aktif</h3>
              </div>
              <span class="text-[11px] font-bold text-zinc-400 font-mono">{{ doctorLeaves.length }} Tanggal</span>
            </div>

            <div v-if="doctorLeaves.length === 0" class="py-6 text-center text-xs text-zinc-400 bg-zinc-50 rounded-xl border border-dashed border-zinc-200">
              Tidak ada pengajuan cuti aktif. Jadwal praktik berjalan normal.
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="leave in doctorLeaves"
                :key="leave.id"
                class="p-3 rounded-xl border border-rose-200 bg-rose-50/50 flex items-center justify-between gap-3 text-xs"
              >
                <div>
                  <div class="font-extrabold text-rose-950">{{ leave.override_date }}</div>
                  <div class="text-[11px] text-zinc-600">{{ leave.reason || 'Cuti Dokter' }}</div>
                </div>

                <button
                  type="button"
                  @click="handleCancelLeave(leave.override_date)"
                  class="px-2.5 py-1.5 rounded-lg border border-rose-300 hover:bg-white text-[11px] font-bold text-rose-700 transition-all cursor-pointer"
                >
                  Buka Kembali
                </button>
              </div>
            </div>
          </div>

          <!-- Reservasi Pasien Mendatang -->
          <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <Users class="w-4 h-4 text-[#0C3B36]" />
                <h3 class="text-sm font-extrabold text-zinc-950">Booking Pasien Masa Depan</h3>
              </div>
              <span class="text-[11px] font-bold text-zinc-400 font-mono">{{ futureReservations.length }} Pasien</span>
            </div>

            <div v-if="futureReservations.length === 0" class="py-6 text-center text-xs text-zinc-400 bg-zinc-50 rounded-xl border border-dashed border-zinc-200">
              Belum ada pasien yang melakukan booking masa depan untuk dokter ini.
            </div>

            <div v-else class="space-y-2 max-h-[380px] overflow-y-auto pr-1">
              <div
                v-for="resv in futureReservations"
                :key="resv.id"
                class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/60 space-y-1.5 text-xs"
              >
                <div class="flex items-center justify-between">
                  <span class="font-mono font-bold text-[#0C3B36]">{{ resv.queue_number }}</span>
                  <span class="px-2 py-0.5 rounded text-[10px] font-extrabold bg-emerald-100 text-emerald-900">
                    {{ resv.appointment_date || 'Terjadwal' }}
                  </span>
                </div>
                <div class="font-bold text-zinc-900">{{ resv.patient_name }}</div>
                <div class="flex items-center justify-between text-[11px] text-zinc-500">
                  <span>Slot: {{ resv.appointment_time || '09:00 WIB' }}</span>
                  <span>{{ resv.patient_phone || '-' }}</span>
                </div>
              </div>
            </div>
          </div>

        </div>

      </div>

    </div>

    <!-- MODAL EDIT SHIFT PRAKTIK MINGGUAN -->
    <div v-if="editingShift" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div class="bg-white rounded-3xl p-6 max-w-sm w-full space-y-4 shadow-xl animate-in fade-in">
        <div class="flex items-center justify-between">
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Atur Jam Praktik</span>
            <h3 class="text-base font-extrabold text-zinc-900">Hari {{ dayNamesMap[editingShift.day_of_week] || 'Shift' }}</h3>
          </div>
          <button @click="editingShift = null" class="text-zinc-400 hover:text-zinc-700 text-sm font-bold">✕</button>
        </div>

        <form @submit.prevent="handleSaveShift" class="space-y-3.5">
          <div>
            <label class="block text-[11px] font-bold text-zinc-700 mb-1">Nama / Waktu Shift:</label>
            <select v-model="editingShift.shift_name" class="w-full px-3 py-2 border rounded-xl text-xs font-bold">
              <option value="PAGI">PAGI</option>
              <option value="SIANG">SIANG</option>
              <option value="SORE">SORE</option>
              <option value="MALAM">MALAM</option>
            </select>
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-[11px] font-bold text-zinc-700 mb-1">Jam Mulai:</label>
              <input type="time" v-model="editingShift.start_time" required class="w-full px-3 py-2 border rounded-xl text-xs font-bold font-mono" />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-zinc-700 mb-1">Jam Selesai:</label>
              <input type="time" v-model="editingShift.end_time" required class="w-full px-3 py-2 border rounded-xl text-xs font-bold font-mono" />
            </div>
          </div>

          <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200 flex items-center justify-between">
            <label class="text-xs font-bold text-zinc-800 cursor-pointer" for="shiftActiveToggle">
              Status Praktik Hari Ini:
            </label>
            <div class="flex items-center gap-2">
              <input
                id="shiftActiveToggle"
                type="checkbox"
                v-model="editingShift.is_active"
                class="w-4 h-4 accent-[#0C3B36] rounded cursor-pointer"
              />
              <span class="text-xs font-bold" :class="editingShift.is_active ? 'text-emerald-700' : 'text-zinc-400'">
                {{ editingShift.is_active ? 'Buka Praktik' : 'Libur' }}
              </span>
            </div>
          </div>

          <div class="flex gap-2 pt-2 border-t border-zinc-100">
            <button type="button" @click="editingShift = null" class="flex-1 py-2.5 rounded-xl bg-zinc-100 font-bold text-xs text-zinc-700 hover:bg-zinc-200">
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSavingShift || !editingShift.start_time || !editingShift.end_time"
              class="flex-1 py-2.5 rounded-xl bg-[#0C3B36] text-white font-bold text-xs hover:bg-[#134E48] transition-all disabled:opacity-50 flex items-center justify-center gap-1.5"
            >
              <Loader2 v-if="isSavingShift" class="w-3.5 h-3.5 animate-spin" />
              <span>{{ isSavingShift ? 'Menyimpan...' : 'Simpan Shift' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
