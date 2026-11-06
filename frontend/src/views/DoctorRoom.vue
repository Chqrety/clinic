<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem, type Drug, type Consultation } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Stethoscope, Send, Plus, Trash2, CheckCircle2, ShieldCheck, AlertCircle, RefreshCw } from 'lucide-vue-next'

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
    anamnesis.value = c.anamnesis || `Pasien mengeluh: ${q.complaint}`
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
      dosage_instruction: '3x1 tablet sesudah makan',
    })
  }
}

function removeRxItem(index: number) {
  prescriptionItems.value.splice(index, 1)
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
    successMsg.value = 'E-Resep terkirim ke Apotek & stok obat di-reserve!'
    store.playChime()
    await loadData()
  } catch (err: any) {
    errorMsg.value = err.message || 'Gagal menyimpan resep'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="py-8 px-6 sm:px-10 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-4 mb-6 pb-6 border-b border-black/10">
      <div>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Ruang Periksa Dokter & E-Resep
        </h1>
        <p class="text-xs text-zinc-500 mt-1">
          Input diagnosa ICD-10 dan racik resep elektronik. Stok ter-reserve otomatis ke sistem farmasi.
        </p>
      </div>

      <button
        @click="loadData"
        :disabled="loading"
        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
        Segarkan Pasien
      </button>
    </div>

    <!-- Feedback Toasts -->
    <div v-if="successMsg" class="mb-6 p-4 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs font-bold flex items-center gap-2">
      <CheckCircle2 class="w-4 h-4 text-emerald-600" />
      {{ successMsg }}
    </div>
    <div v-if="errorMsg" class="mb-6 p-4 rounded-2xl bg-rose-50 border border-rose-200 text-rose-800 text-xs font-bold flex items-center gap-2">
      <AlertCircle class="w-4 h-4 text-rose-600" />
      {{ errorMsg }}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">

      <!-- Left Column: Patient Waiting List (4 cols) -->
      <div class="lg:col-span-4 space-y-3">
        <div class="flex items-center justify-between pl-1">
          <span class="text-xs font-bold uppercase tracking-wider text-zinc-400">
            Antrian Pasien Menunggu
          </span>
          <span class="text-xs font-extrabold text-zinc-950 bg-stone-200/70 px-2 py-0.5 rounded-full">
            {{ queues.length }}
          </span>
        </div>

        <div class="space-y-2 max-h-[640px] overflow-y-auto pr-1">
          <div
            v-for="q in queues"
            :key="q.id"
            @click="selectPatient(q)"
            :class="[
              'p-4 rounded-2xl border transition-all cursor-pointer active:scale-[0.99]',
              selectedQueue?.id === q.id
                ? 'bg-zinc-950 text-white border-zinc-950 shadow-sm'
                : 'bg-white text-zinc-900 border-black/5 hover:border-black/15'
            ]"
          >
            <div class="flex items-center justify-between mb-1">
              <span class="text-base font-extrabold font-mono">{{ q.queue_number }}</span>
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
            <div :class="['text-xs line-clamp-1 mt-0.5', selectedQueue?.id === q.id ? 'text-zinc-400' : 'text-zinc-500']">
              "{{ q.complaint }}"
            </div>
          </div>

          <div v-if="queues.length === 0" class="p-8 text-center text-zinc-400 bg-white rounded-2xl border border-black/5 text-xs">
            Belum ada pasien yang menunggu di ruang periksa.
          </div>
        </div>
      </div>

      <!-- Right Column: Unified Doctor Workstation (8 cols) -->
      <div v-if="selectedQueue" class="lg:col-span-8 bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.02)] space-y-6">

        <!-- Patient Header Strip -->
        <div class="flex flex-col sm:flex-row sm:items-baseline justify-between pb-4 border-b border-zinc-100 gap-2">
          <div>
            <div class="flex items-center gap-2">
              <span class="text-2xl font-black font-mono text-zinc-950">{{ selectedQueue.queue_number }}</span>
              <h2 class="text-xl font-extrabold text-zinc-900">{{ selectedQueue.patient_name }}</h2>
            </div>
            <div class="text-xs text-zinc-500 mt-0.5">
              NIK: {{ selectedQueue.patient_nik }} • Keluhan Awal: "{{ selectedQueue.complaint }}"
            </div>
          </div>

          <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-blue-50 text-blue-700 text-xs font-bold border border-blue-200/50 self-start sm:self-auto">
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>IHS: {{ selectedQueue.patient_ihs || 'P01928374611' }}</span>
          </div>
        </div>

        <!-- Examination Form -->
        <div class="space-y-4">
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-1">
              Catatan Anamnesis / Pemeriksaan Fisik
            </label>
            <textarea
              v-model="anamnesis"
              rows="3"
              class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 focus:border-zinc-950 focus:bg-white text-zinc-900 text-xs font-medium outline-none transition-all resize-none"
              placeholder="Catatan tensi, suhu, dan keluhan mendalam..."
            ></textarea>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-1">
                Diagnosa Klinis
              </label>
              <input
                v-model="diagnosis"
                type="text"
                placeholder="Contoh: ISPA Akut"
                class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 focus:border-zinc-950 focus:bg-white text-zinc-900 text-xs font-bold outline-none transition-all"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-1">
                Kode ICD-10 (Kemenkes)
              </label>
              <select
                v-model="icd10Code"
                class="w-full px-4 py-2.5 rounded-2xl bg-[#F9F9F7] border border-black/10 focus:border-zinc-950 focus:bg-white text-zinc-900 text-xs font-bold outline-none transition-all"
              >
                <option v-for="item in icd10Presets" :key="item.code" :value="item.code">
                  {{ item.code }} — {{ item.name }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- E-Prescription Builder (Direct Rows, Anti Card-in-Card) -->
        <div class="pt-4 border-t border-zinc-100 space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-sm font-extrabold text-zinc-950">Resep Obat Elektronik</h3>
              <p class="text-[11px] text-zinc-500">Stok obat otomatis di-reserve saat resep dikirim ke apotek.</p>
            </div>
            <button
              type="button"
              @click="addRxItem"
              class="inline-flex items-center gap-1 px-3 py-1.5 rounded-full bg-stone-100 hover:bg-stone-200 text-xs font-bold text-zinc-800 transition-all active:scale-95"
            >
              <Plus class="w-3.5 h-3.5" /> Tambah Obat
            </button>
          </div>

          <div class="space-y-2">
            <div
              v-for="(item, index) in prescriptionItems"
              :key="index"
              class="p-3 rounded-2xl bg-[#F9F9F7] border border-black/5 flex flex-col md:flex-row md:items-center gap-3"
            >
              <div class="flex-1">
                <select
                  v-model="item.drug_id"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-bold text-zinc-900 outline-none"
                >
                  <option v-for="d in drugs" :key="d.id" :value="d.id">
                    {{ d.name }} (Tersedia: {{ d.available_stock }} {{ d.unit }})
                  </option>
                </select>
              </div>

              <div class="w-24">
                <input
                  v-model.number="item.quantity"
                  type="number"
                  min="1"
                  placeholder="Jumlah"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-bold text-zinc-900 outline-none text-center"
                />
              </div>

              <div class="flex-1">
                <input
                  v-model="item.dosage_instruction"
                  type="text"
                  placeholder="Aturan pakai: 3x1 tablet sesudah makan"
                  class="w-full px-3 py-2 rounded-xl bg-white border border-black/10 text-xs font-medium text-zinc-900 outline-none"
                />
              </div>

              <button
                type="button"
                @click="removeRxItem(index)"
                class="self-end md:self-center p-2 text-rose-500 hover:text-rose-700 hover:bg-rose-50 rounded-xl transition-colors"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            @click="handleSaveConsultation"
            :disabled="loading"
            class="w-full mt-4 py-3.5 rounded-full font-bold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 text-xs disabled:opacity-50"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Kirim E-Resep ke Apotek & Selesaikan Sesi</span>
          </button>
        </div>

      </div>

    </div>
  </div>
</template>
