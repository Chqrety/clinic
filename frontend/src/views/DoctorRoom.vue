<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api, type QueueItem, type Drug, type Consultation } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Stethoscope, Send, Plus, Trash2, CheckCircle2, ShieldCheck, AlertCircle } from 'lucide-vue-next'

const store = useClinicStore()

const queues = ref<QueueItem[]>([])
const drugs = ref<Drug[]>([])
const selectedQueue = ref<QueueItem | null>(null)
const consultation = ref<Consultation | null>(null)

const doctorName = ref('dr. Hendra Pratama, Sp.PD')
const anamnesis = ref('')
const diagnosis = ref('')
const icd10Code = ref('J06.9')

interface RxFormItem {
  drug_id: number
  quantity: number
  dosage_instruction: string
}
const prescriptionItems = ref<RxFormItem[]>([
  { drug_id: 1, quantity: 10, dosage_instruction: '3x1 tablet sehari sesudah makan' },
])

const loading = ref(false)
const successMsg = ref('')
const errorMsg = ref('')

const icd10Presets = [
  { code: 'J06.9', name: 'ISPA / Infeksi Saluran Pernapasan Akut' },
  { code: 'K29.7', name: 'Gastritis / Maag Akut' },
  { code: 'R50.9', name: 'Febris / Demam Tanpa Sebab Spesifik' },
  { code: 'E11.9', name: 'Diabetes Melitus Tipe 2' },
  { code: 'I10', name: 'Hipertensi Esensial' },
]

async function loadData() {
  loading.value = true
  try {
    const [qList, dList] = await Promise.all([api.getQueues(), api.getDrugs()])
    queues.value = qList.filter(q => q.payment_status === 'PAID')
    drugs.value = dList

    if (queues.value.length > 0 && !selectedQueue.value) {
      selectPatient(queues.value[0])
    }
  } catch (e: any) {
    errorMsg.value = e.message
  } finally {
    loading.value = false
  }
}

async function selectPatient(q: QueueItem) {
  selectedQueue.value = q
  store.setSelectedDoctorQueue(q.id)
  successMsg.value = ''
  errorMsg.value = ''

  try {
    const c = await api.getConsultationByQueue(q.id)
    consultation.value = c
    anamnesis.value = c.anamnesis || `Pasien datang dengan keluhan: ${q.complaint}`
    diagnosis.value = c.diagnosis || ''
    icd10Code.value = c.icd10_code || 'J06.9'

    if (c.prescription && c.prescription.items.length > 0) {
      prescriptionItems.value = c.prescription.items.map(it => ({
        drug_id: it.drug_id,
        quantity: it.quantity,
        dosage_instruction: it.dosage_instruction,
      }))
    }
  } catch (e) {
    // ignore
  }
}

function addRxItem() {
  if (drugs.value.length > 0) {
    prescriptionItems.value.push({
      drug_id: drugs.value[0].id,
      quantity: 10,
      dosage_instruction: '3x1 sehari sesudah makan',
    })
  }
}

function removeRxItem(index: number) {
  prescriptionItems.value.splice(index, 1)
}

function getDrugStockInfo(drugId: number) {
  const d = drugs.value.find(item => item.id === drugId)
  if (!d) return { available: 0, unit: '' }
  return { available: d.available_stock, unit: d.unit, name: d.name }
}

async function handleSaveConsultation() {
  if (!selectedQueue.value) return

  loading.value = true
  successMsg.value = ''
  errorMsg.value = ''

  try {
    const saved = await api.saveConsultation({
      queue_id: selectedQueue.value.id,
      doctor_name: doctorName.value,
      anamnesis: anamnesis.value,
      diagnosis: diagnosis.value,
      icd10_code: icd10Code.value,
      prescription: prescriptionItems.value,
    })
    consultation.value = saved
    successMsg.value = 'E-Resep berhasil dikirim ke Apotek & stok obat otomatis di-reserve!'
    store.playChime()
    await loadData() // refresh stock levels
  } catch (err: any) {
    errorMsg.value = err.message || 'Gagal menyimpan rekam medis / e-resep'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 px-4 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-blue-600 text-white mb-2 shadow-sm">
          Ruang Dokter
        </span>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Pemeriksaan & E-Resep Real-Time
        </h1>
        <p class="text-sm text-zinc-500 mt-0.5">
          Resep terhubung langsung ke Apotek. Stok obat langsung ter-reserve saat dokter menyimpan.
        </p>
      </div>
    </div>

    <!-- Alert Messages -->
    <div v-if="successMsg" class="mb-6 p-4 rounded-3xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-sm font-semibold flex items-center gap-2">
      <CheckCircle2 class="w-5 h-5 text-emerald-600" />
      {{ successMsg }}
    </div>
    <div v-if="errorMsg" class="mb-6 p-4 rounded-3xl bg-rose-50 border border-rose-200 text-rose-800 text-sm font-semibold flex items-center gap-2">
      <AlertCircle class="w-5 h-5 text-rose-600" />
      {{ errorMsg }}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

      <!-- Left Column: Patient Queue List (4 cols) -->
      <div class="lg:col-span-4 space-y-3">
        <h2 class="text-sm font-bold uppercase tracking-wider text-zinc-400 pl-1">
          Daftar Pasien Menunggu ({{ queues.length }})
        </h2>

        <div class="space-y-2 max-h-[600px] overflow-y-auto pr-1">
          <div
            v-for="q in queues"
            :key="q.id"
            @click="selectPatient(q)"
            :class="[
              'p-4 rounded-3xl border transition-all cursor-pointer active:scale-[0.99]',
              selectedQueue?.id === q.id
                ? 'bg-zinc-950 text-white border-zinc-950 shadow-md'
                : 'bg-white text-zinc-900 border-black/5 hover:border-black/15 shadow-[0_2px_12px_rgba(0,0,0,0.02)]'
            ]"
          >
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-lg font-extrabold font-mono">{{ q.queue_number }}</span>
              <span
                :class="[
                  'text-[10px] font-bold px-2 py-0.5 rounded-full',
                  selectedQueue?.id === q.id ? 'bg-zinc-800 text-zinc-200' : 'bg-stone-100 text-zinc-700'
                ]"
              >
                {{ q.queue_status }}
              </span>
            </div>
            <div class="font-bold text-sm">{{ q.patient_name }}</div>
            <div :class="['text-xs line-clamp-1 mt-1', selectedQueue?.id === q.id ? 'text-zinc-400' : 'text-zinc-500']">
              "{{ q.complaint }}"
            </div>
          </div>

          <div v-if="queues.length === 0" class="p-8 text-center text-zinc-400 bg-white rounded-3xl border border-black/5 text-sm">
            Tidak ada pasien yang mengantri saat ini.
          </div>
        </div>
      </div>

      <!-- Right Column: Consultation & E-Prescription Form (8 cols) -->
      <div v-if="selectedQueue" class="lg:col-span-8 space-y-6">

        <!-- Patient Header Card -->
        <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div class="flex items-center gap-2">
              <span class="text-2xl font-extrabold text-zinc-950 font-mono">{{ selectedQueue.queue_number }}</span>
              <h2 class="text-xl font-bold text-zinc-900">{{ selectedQueue.patient_name }}</h2>
            </div>
            <div class="text-xs text-zinc-500 mt-1 flex flex-wrap gap-x-4">
              <span>NIK: <strong>{{ selectedQueue.patient_nik }}</strong></span>
              <span>HP: <strong>{{ selectedQueue.patient_phone }}</strong></span>
            </div>
          </div>

          <!-- SATUSEHAT Patient Verification Badge -->
          <div class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-blue-50 text-blue-700 border border-blue-200/60 text-xs font-semibold self-start sm:self-auto">
            <ShieldCheck class="w-4 h-4 text-blue-600" />
            <span>IHS: {{ selectedQueue.patient_ihs || 'P01928374611' }}</span>
          </div>
        </div>

        <!-- Examination Form -->
        <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] space-y-4">
          <h3 class="text-base font-bold text-zinc-950">Catatan Medis & Diagnosa (ICD-10)</h3>

          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
              Anamnesis / Pemeriksaan Fisik
            </label>
            <textarea
              v-model="anamnesis"
              rows="3"
              class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all resize-none"
              placeholder="Tuliskan catatan keluhan, tekanan darah, suhu badan..."
            ></textarea>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
                Diagnosa Dokter
              </label>
              <input
                v-model="diagnosis"
                type="text"
                placeholder="Contoh: ISPA akut"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold uppercase tracking-wider text-zinc-500 mb-1.5">
                Kode ICD-10 (Standar Kemenkes)
              </label>
              <select
                v-model="icd10Code"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/5 focus:border-zinc-900 focus:bg-white text-zinc-900 text-sm font-medium outline-none transition-all"
              >
                <option v-for="item in icd10Presets" :key="item.code" :value="item.code">
                  {{ item.code }} — {{ item.name }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- E-Prescription Builder (Killer Feature) -->
        <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-base font-bold text-zinc-950">Resep Elektronik (E-Prescription)</h3>
              <p class="text-xs text-zinc-500 mt-0.5">
                Apotek akan menerima resep ini secara real-time. Stok di-reserve otomatis.
              </p>
            </div>
            <button
              type="button"
              @click="addRxItem"
              class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-full bg-stone-100 hover:bg-stone-200 text-xs font-bold text-zinc-800 transition-all active:scale-95"
            >
              <Plus class="w-3.5 h-3.5" /> Tambah Obat
            </button>
          </div>

          <div class="space-y-3">
            <div
              v-for="(item, index) in prescriptionItems"
              :key="index"
              class="p-4 rounded-2xl bg-[#F9F9F7] border border-black/5 flex flex-col md:flex-row md:items-center gap-3"
            >
              <!-- Drug Select -->
              <div class="flex-1">
                <label class="block text-[11px] font-bold text-zinc-500 mb-1">Pilih Obat</label>
                <select
                  v-model="item.drug_id"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-bold text-zinc-900 outline-none"
                >
                  <option v-for="d in drugs" :key="d.id" :value="d.id">
                    {{ d.name }} (Stok Tersedia: {{ d.available_stock }} {{ d.unit }})
                  </option>
                </select>
              </div>

              <!-- Quantity -->
              <div class="w-24">
                <label class="block text-[11px] font-bold text-zinc-500 mb-1">Jumlah</label>
                <input
                  v-model.number="item.quantity"
                  type="number"
                  min="1"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-bold text-zinc-900 outline-none"
                />
              </div>

              <!-- Dosage Instruction -->
              <div class="flex-1">
                <label class="block text-[11px] font-bold text-zinc-500 mb-1">Aturan Pakai (Signa)</label>
                <input
                  v-model="item.dosage_instruction"
                  type="text"
                  placeholder="3x1 tablet sesudah makan"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-medium text-zinc-900 outline-none"
                />
              </div>

              <!-- Delete button -->
              <button
                type="button"
                @click="removeRxItem(index)"
                class="self-end md:self-center mt-2 md:mt-4 p-2 text-rose-500 hover:text-rose-700 hover:bg-rose-50 rounded-xl transition-colors"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            @click="handleSaveConsultation"
            :disabled="loading"
            class="w-full mt-4 py-4 rounded-full font-bold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 text-sm disabled:opacity-50"
          >
            <Send class="w-4 h-4" />
            <span>Kirim E-Resep ke Apotek & Selesaikan Periksa</span>
          </button>
        </div>

      </div>

    </div>
  </div>
</template>
