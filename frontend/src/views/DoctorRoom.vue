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
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

// Doctor Session State
const availableDoctors = [
  { name: 'dr. Hendra Pratama, Sp.A', unit: 'KIDINEST', title: 'Dokter Spesialis Anak' },
  { name: 'dr. Nadia Kusuma, Sp.KK', unit: 'SSC', title: 'Dokter Spesialis Kulit & Estetika' },
]

const currentDoctor = ref(
  authStore.user?.role === 'DOCTOR' && authStore.user.name
    ? availableDoctors.find((d) => d.name.includes(authStore.user!.name.split(' ')[1] || '')) || availableDoctors[0]
    : availableDoctors[0]
)

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
    queues.value = await api.getDoctorQueues(currentDoctor.value.name, currentDoctor.value.unit)
    if (!activeQueue.value && queues.value.length > 0) {
      selectPatient(queues.value[0])
    }
  } catch (err: any) {
    console.error('Gagal memuat antrian dokter:', err)
  } finally {
    loadingQueues.value = false
  }
}

async function loadDrugs() {
  loadingDrugs.value = true
  try {
    catalogDrugs.value = await api.getDoctorDrugs(currentDoctor.value.unit)
  } catch (err: any) {
    console.error('Gagal memuat katalog obat:', err)
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
    if (queues.value.length > 0) {
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
})
</script>

<template>
  <div class="py-8 px-4 sm:px-8 max-w-7xl mx-auto space-y-6">
    <!-- Top Bar: Doctor Identity & Fast Switcher -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-black/10">
      <div>
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-blue-50 text-blue-900 border border-blue-200">
            <span class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></span>
            Ruang Periksa Dokter
          </span>
          <span class="text-xs text-zinc-400 font-mono">SCMS Clinical Room</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-zinc-950 tracking-tight mt-1">
          Pemeriksaan & E-Resep Cepat
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Antrian FIFO berkeadilan, input diagnosa ICD-10, tindakan klinis, dan peresepan obat racikan manual maupun katalog.
        </p>
      </div>

      <!-- Quick Doctor Switcher (For Demo & Multi-Unit) -->
      <div class="flex items-center gap-2 bg-white rounded-2xl p-1.5 border border-black/5 shadow-sm">
        <button
          v-for="doc in availableDoctors"
          :key="doc.name"
          @click="switchDoctor(doc)"
          :class="[
            'px-3.5 py-2 rounded-xl text-xs font-bold transition-all text-left flex items-center gap-2',
            currentDoctor.name === doc.name
              ? 'bg-zinc-950 text-white shadow-sm'
              : 'hover:bg-zinc-100 text-zinc-600',
          ]"
        >
          <div
            :class="[
              'w-2 h-2 rounded-full',
              doc.unit === 'KIDINEST' ? 'bg-amber-400' : 'bg-rose-400',
            ]"
          ></div>
          <div>
            <div class="text-[11px] leading-tight font-extrabold">{{ doc.name }}</div>
            <div :class="['text-[9px] font-normal', currentDoctor.name === doc.name ? 'text-zinc-400' : 'text-zinc-500']">
              {{ doc.unit === 'KIDINEST' ? 'Poli Anak' : 'Estetika Kulit' }}
            </div>
          </div>
        </button>
      </div>
    </div>

    <!-- Main Workspace Layout: 2 Columns -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- LEFT COLUMN: FIFO QUEUE LIST (5 cols) -->
      <div class="lg:col-span-5 space-y-4">
        <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <h2 class="text-xs font-black uppercase tracking-wider text-zinc-950">
                Antrian Menunggu FIFO
              </h2>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold bg-zinc-100 text-zinc-700">
                {{ queues.length }} Pasien
              </span>
            </div>

            <button
              @click="loadDoctorQueues"
              :disabled="loadingQueues"
              class="p-1.5 rounded-xl hover:bg-zinc-100 text-zinc-500 transition-all"
              title="Segarkan antrian"
            >
              <RefreshCw :class="['w-3.5 h-3.5', loadingQueues && 'animate-spin']" />
            </button>
          </div>

          <div v-if="loadingQueues" class="py-12 text-center text-xs text-zinc-400">
            <RefreshCw class="w-5 h-5 mx-auto animate-spin mb-2 text-zinc-400" />
            Memuat daftar antrian...
          </div>

          <div v-else-if="queues.length === 0" class="py-12 text-center text-xs text-zinc-400">
            <CheckCircle2 class="w-8 h-8 mx-auto text-emerald-500 mb-2" />
            Semua antrian dokter telah diperiksa!
          </div>

          <div v-else class="space-y-2.5 max-h-[640px] overflow-y-auto pr-1">
            <div
              v-for="p in queues"
              :key="p.queue_id"
              @click="selectPatient(p)"
              :class="[
                'p-4 rounded-2xl border transition-all cursor-pointer space-y-2',
                activeQueue?.queue_id === p.queue_id
                  ? 'bg-zinc-950 text-white border-zinc-950 shadow-md ring-2 ring-zinc-900/10'
                  : 'bg-zinc-50/70 hover:bg-zinc-100 text-zinc-900 border-black/5',
              ]"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span
                    :class="[
                      'px-2.5 py-0.5 rounded-full font-mono text-xs font-black',
                      activeQueue?.queue_id === p.queue_id ? 'bg-white text-zinc-950' : 'bg-zinc-950 text-white',
                    ]"
                  >
                    {{ p.queue_number }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-bold px-2 py-0.5 rounded-full',
                      p.clinic_unit === 'KIDINEST' ? 'bg-amber-100 text-amber-900' : 'bg-rose-100 text-rose-900',
                    ]"
                  >
                    {{ p.clinic_unit }}
                  </span>
                  <span
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-stone-200 text-stone-800"
                  >
                    {{ p.category }}
                  </span>
                </div>

                <span
                  :class="[
                    'text-[10px] font-bold px-2.5 py-0.5 rounded-full',
                    p.queue_status === 'CALLING'
                      ? 'bg-amber-400 text-zinc-950 animate-pulse'
                      : p.queue_status === 'IN_EXAM'
                      ? 'bg-emerald-400 text-zinc-950'
                      : 'bg-zinc-200 text-zinc-700',
                  ]"
                >
                  {{ p.queue_status }}
                </span>
              </div>

              <div>
                <div class="font-extrabold text-sm tracking-tight">{{ p.patient_name }}</div>
                <div :class="['text-xs mt-0.5', activeQueue?.queue_id === p.queue_id ? 'text-zinc-400' : 'text-zinc-500']">
                  Dokter: {{ p.doctor_name }}
                </div>
              </div>

              <!-- Quick Action Buttons for Doctor -->
              <div class="pt-2 border-t border-black/5 flex items-center justify-end gap-1.5" @click.stop>
                <button
                  v-if="p.queue_status !== 'CALLING' && p.queue_status !== 'IN_EXAM'"
                  @click="callPatient(p)"
                  class="inline-flex items-center gap-1 px-3 py-1 rounded-xl bg-amber-50 hover:bg-amber-100 text-amber-900 text-[11px] font-bold transition-all"
                >
                  <Volume2 class="w-3 h-3 text-amber-600" />
                  Panggil Pasien
                </button>

                <button
                  v-if="p.queue_status === 'CALLING'"
                  @click="startExam(p)"
                  class="inline-flex items-center gap-1 px-3 py-1 rounded-xl bg-emerald-500 hover:bg-emerald-600 text-white text-[11px] font-bold transition-all shadow-xs"
                >
                  <Stethoscope class="w-3 h-3" />
                  Mulai Periksa
                </button>

                <span v-if="p.queue_status === 'IN_EXAM'" class="text-[10px] font-bold text-emerald-400">
                  Sedang Diperiksa
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT COLUMN: CONSULTATION & FAST E-PRESCRIPTION (7 cols) -->
      <div class="lg:col-span-7">
        <div v-if="!activeQueue" class="bg-white rounded-3xl p-12 border border-black/5 text-center shadow-sm">
          <Stethoscope class="w-12 h-12 mx-auto text-zinc-300 mb-3" />
          <h3 class="text-base font-extrabold text-zinc-800">Pilih Pasien untuk Memulai Konsultasi</h3>
          <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
            Klik nomor antrian di kolom kiri untuk membuka form rekam medis, anamnesis, input tindakan klinis, dan penulisan e-resep cepat.
          </p>
        </div>

        <div v-else class="bg-white rounded-3xl p-6 sm:p-7 border border-black/5 shadow-sm space-y-6">
          <!-- Active Patient Header Banner -->
          <div class="flex items-start justify-between pb-4 border-b border-black/5">
            <div>
              <div class="flex items-center gap-2">
                <span class="px-3 py-1 rounded-full font-mono text-sm font-black bg-zinc-950 text-white">
                  {{ activeQueue.queue_number }}
                </span>
                <span class="text-xs font-bold px-2.5 py-0.5 rounded-full bg-zinc-100 text-zinc-700">
                  {{ activeQueue.clinic_unit }}
                </span>
                <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-stone-100 text-stone-700">
                  {{ activeQueue.category }}
                </span>
              </div>
              <h2 class="text-xl font-black text-zinc-950 mt-2">{{ activeQueue.patient_name }}</h2>
              <p class="text-xs text-zinc-500">
                Pemeriksa: {{ currentDoctor.name }} • Status: {{ activeQueue.queue_status }}
              </p>
            </div>

            <!-- Call & Exam Status Controls -->
            <div class="flex items-center gap-2">
              <button
                @click="callPatient(activeQueue)"
                class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-xs font-bold text-zinc-800 transition-all"
              >
                <Volume2 class="w-3.5 h-3.5" />
                Panggil Ulang
              </button>
            </div>
          </div>

          <!-- Section 1: Anamnesis & Keluhan -->
          <div class="space-y-1.5">
            <label class="text-xs font-extrabold uppercase tracking-wider text-zinc-500">
              Anamnesis / Keluhan Klinis Pasien:
            </label>
            <textarea
              v-model="anamnesis"
              rows="2"
              placeholder="Contoh: Batuk pilek sejak 2 hari, demam naik turun malam hari..."
              class="w-full px-4 py-2.5 rounded-2xl bg-zinc-50 border border-black/10 text-xs focus:ring-2 focus:ring-zinc-950 transition-all"
            ></textarea>
          </div>

          <!-- Section 2: Diagnosa & ICD-10 Presets -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <label class="text-xs font-extrabold uppercase tracking-wider text-zinc-500">
                Diagnosa Dokter & ICD-10:
              </label>
              <span class="text-[10px] text-zinc-400">Kode standar medis Kemenkes</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
              <input
                v-model="diagnosis"
                type="text"
                placeholder="Nama diagnosa klinis"
                class="sm:col-span-2 px-4 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs font-bold focus:ring-2 focus:ring-zinc-950"
              />
              <input
                v-model="icd10Code"
                type="text"
                placeholder="Kode ICD-10"
                class="px-3 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs font-mono font-bold text-center focus:ring-2 focus:ring-zinc-950"
              />
            </div>

            <!-- Quick ICD-10 Chips -->
            <div class="flex items-center gap-1.5 flex-wrap pt-1">
              <span class="text-[10px] font-bold text-zinc-400">Pilihan Cepat:</span>
              <button
                v-for="chip in currentPresets"
                :key="chip.code"
                type="button"
                @click="diagnosis = chip.name; icd10Code = chip.code"
                class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 hover:bg-zinc-200 text-zinc-700 transition-all"
              >
                {{ chip.code }} - {{ chip.name.split('/')[0] }}
              </button>
            </div>
          </div>

          <!-- Section 3: Tindakan Medis & Jasa Dokter / Prosedur -->
          <div class="bg-zinc-50/80 rounded-2xl p-4 border border-black/5 space-y-3">
            <div class="flex items-center justify-between">
              <label class="text-xs font-extrabold uppercase tracking-wider text-zinc-600">
                Tindakan Medis & Biaya Jasa:
              </label>
              <span class="text-[10px] text-zinc-400">Masuk ke billing invoice kasir</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div class="sm:col-span-2 space-y-1">
                <label class="text-[11px] font-bold text-zinc-700">Catatan Tindakan Medis / Prosedur:</label>
                <input
                  v-model="actionNotes"
                  type="text"
                  placeholder="Contoh: Injeksi Vaksin DTP / Peeling Acne"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs focus:ring-2 focus:ring-zinc-950"
                />
              </div>

              <div class="space-y-1">
                <label class="text-[11px] font-bold text-zinc-700">Biaya Jasa Tindakan (Rp):</label>
                <input
                  v-model.number="actionFee"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-mono font-bold text-right focus:ring-2 focus:ring-zinc-950"
                />
              </div>
            </div>
          </div>

          <!-- Section 4: Fast E-Prescription Table -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-xs font-extrabold uppercase tracking-wider text-zinc-950 flex items-center gap-1.5">
                  <Pill class="w-3.5 h-3.5 text-blue-600" />
                  Tabel Peresepan Elektronik (E-Resep Cepat)
                </h3>
                <p class="text-[11px] text-zinc-400">Mendukung obat katalog apotek & racikan manual puyer/salep</p>
              </div>

              <div class="flex items-center gap-2">
                <button
                  type="button"
                  @click="addCatalogDrugRow"
                  class="inline-flex items-center gap-1 px-3 py-1.5 rounded-xl bg-blue-50 hover:bg-blue-100 text-blue-900 text-xs font-bold transition-all"
                >
                  <Plus class="w-3.5 h-3.5" />
                  Obat Katalog
                </button>

                <button
                  type="button"
                  @click="addManualRacikanRow"
                  class="inline-flex items-center gap-1 px-3 py-1.5 rounded-xl bg-amber-50 hover:bg-amber-100 text-amber-900 text-xs font-bold transition-all"
                >
                  <Sparkles class="w-3.5 h-3.5" />
                  Racikan Manual
                </button>
              </div>
            </div>

            <!-- Rx Items Table -->
            <div v-if="prescriptionRows.length === 0" class="py-6 text-center text-xs text-zinc-400 bg-zinc-50 rounded-2xl border border-dashed border-black/10">
              Belum ada obat yang ditambahkan ke resep. Klik salah satu tombol di atas untuk menambah item.
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="(row, idx) in prescriptionRows"
                :key="idx"
                class="p-3.5 rounded-2xl bg-zinc-50 border border-black/5 space-y-2.5"
              >
                <div class="flex items-center justify-between gap-2">
                  <div class="flex items-center gap-2 flex-1">
                    <span
                      :class="[
                        'px-2 py-0.5 rounded-full text-[9px] font-bold',
                        row.is_manual ? 'bg-amber-100 text-amber-900' : 'bg-blue-100 text-blue-900',
                      ]"
                    >
                      {{ row.is_manual ? 'Racikan Manual' : 'Katalog Apotek' }}
                    </span>

                    <!-- Drug Name Selector or Manual Text Input -->
                    <div v-if="!row.is_manual" class="flex-1">
                      <select
                        :value="row.drug_id"
                        @change="(e: any) => onDrugSelected(row, parseInt(e.target.value))"
                        class="w-full px-3 py-1.5 rounded-xl bg-white border border-black/10 text-xs font-bold focus:ring-2 focus:ring-zinc-950"
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
                        placeholder="Ketik nama obat racikan puyer / sirup..."
                        class="w-full px-3 py-1.5 rounded-xl bg-white border border-black/10 text-xs font-bold focus:ring-2 focus:ring-zinc-950"
                      />
                    </div>
                  </div>

                  <button
                    type="button"
                    @click="removePrescriptionRow(idx)"
                    class="p-1 rounded-lg text-zinc-400 hover:text-rose-600 hover:bg-rose-50 transition-all"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>

                <!-- Row 2: Quantity, Price, Subtotal, and Dosage Instruction (Signa) -->
                <div class="grid grid-cols-1 sm:grid-cols-12 gap-2 items-center text-xs">
                  <div class="sm:col-span-2">
                    <label class="text-[10px] text-zinc-400 font-bold">Jumlah:</label>
                    <input
                      v-model.number="row.quantity"
                      type="number"
                      min="1"
                      class="w-full px-2.5 py-1 text-center font-mono font-bold rounded-lg bg-white border border-black/10"
                    />
                  </div>

                  <div class="sm:col-span-3">
                    <label class="text-[10px] text-zinc-400 font-bold">Harga Satuan:</label>
                    <input
                      v-model.number="row.unit_price"
                      type="number"
                      min="0"
                      :disabled="!row.is_manual"
                      class="w-full px-2.5 py-1 text-right font-mono font-bold rounded-lg bg-white border border-black/10 disabled:bg-zinc-100"
                    />
                  </div>

                  <div class="sm:col-span-5">
                    <label class="text-[10px] text-zinc-400 font-bold">Aturan Pakai (Signa):</label>
                    <input
                      v-model="row.dosage_instruction"
                      type="text"
                      placeholder="Contoh: 3x1 sendok takar sesudah makan"
                      class="w-full px-2.5 py-1 rounded-lg bg-white border border-black/10"
                    />
                  </div>

                  <div class="sm:col-span-2 text-right">
                    <div class="text-[10px] text-zinc-400 font-bold">Subtotal:</div>
                    <div class="font-mono font-extrabold text-zinc-900">
                      {{ formatRupiah(row.quantity * row.unit_price) }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Prescription Summary Banner -->
              <div class="bg-blue-50/60 rounded-2xl p-3.5 border border-blue-100 flex items-center justify-between text-xs">
                <span class="font-bold text-blue-950">Estimasi Total Biaya Obat Apotek:</span>
                <span class="font-mono font-black text-sm text-blue-900">{{ formatRupiah(totalRxCost) }}</span>
              </div>
            </div>
          </div>

          <!-- Feedback Banners -->
          <div v-if="submitSuccessMsg" class="p-4 rounded-2xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-950 flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
            <span>{{ submitSuccessMsg }}</span>
          </div>

          <div v-if="submitErrorMsg" class="p-4 rounded-2xl bg-rose-50 border border-rose-200 text-xs text-rose-950 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ submitErrorMsg }}</span>
          </div>

          <!-- Submit Button -->
          <div class="pt-2">
            <button
              @click="submitConsultation"
              :disabled="isSubmitting"
              class="w-full py-4 px-6 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-sm shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <Send class="w-4 h-4 text-emerald-400" />
              {{ isSubmitting ? 'Memproses E-Resep & Konsultasi...' : 'Kirim E-Resep ke Apotek & Selesai Periksa' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
