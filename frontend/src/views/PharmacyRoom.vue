<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  api,
  triggerBlobDownload,
  type PharmacyPrescription,
} from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useClinicStore } from '../stores/clinic'
import {
  Pill,
  CheckCircle2,
  RefreshCw,
  Clock,
  Sparkles,
  AlertCircle,
  RotateCcw,
  User,
  Stethoscope,
  ArrowRight,
  ShieldCheck,
  Send,
  Download,
  Upload,
  FileSpreadsheet,
  X,
  Loader2,
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

const prescriptions = ref<PharmacyPrescription[]>([])
const loading = ref(false)
const processingId = ref<number | null>(null)
const filterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const statusFilter = ref<'ALL' | 'RESERVED' | 'PREPARING' | 'DISPENSED'>('ALL')
const searchQuery = ref('')
const selectedRx = ref<PharmacyPrescription | null>(null)

function formatRupiah(val?: number) {
  return 'Rp ' + (val || 0).toLocaleString('id-ID')
}

function formatDate(val?: string) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return val
  }
}

async function fetchPrescriptions() {
  loading.value = true
  try {
    prescriptions.value = await api.getPharmacyPrescriptions(filterUnit.value)
    if (selectedRx.value) {
      const match = prescriptions.value.find((p) => p.id === selectedRx.value!.id)
      selectedRx.value = match || (prescriptions.value[0] || null)
    } else if (prescriptions.value.length > 0) {
      selectedRx.value = prescriptions.value[0]
    }
  } catch (err: any) {
    console.error('Gagal memuat resep:', err)
  } finally {
    loading.value = false
  }
}

// Invariant Actions
async function startPreparing(id: number) {
  processingId.value = id
  try {
    await api.preparePharmacyPrescription(id)
    clinicStore.playChime()
    await fetchPrescriptions()
  } catch (err: any) {
    alert(err.message || 'Gagal mengubah status racik')
  } finally {
    processingId.value = null
  }
}

async function dispense(id: number) {
  processingId.value = id
  try {
    await api.dispensePharmacyPrescription(id)
    clinicStore.playChime()
    await fetchPrescriptions()
  } catch (err: any) {
    alert(err.message || 'Gagal memproses penyerahan obat')
  } finally {
    processingId.value = null
  }
}

async function cancel(id: number) {
  if (!confirm('Yakin ingin membatalkan resep ini? Stok obat yang terkunci (reserved) akan dikembalikan ke apotek.')) {
    return
  }
  processingId.value = id
  try {
    await api.cancelPharmacyPrescription(id)
    await fetchPrescriptions()
  } catch (err: any) {
    alert(err.message || 'Gagal membatalkan resep')
  } finally {
    processingId.value = null
  }
}

function goToBilling(queueId: number) {
  clinicStore.setSelectedInvoiceQueue(queueId)
  clinicStore.setRole('cashier')
}

// Opname CSV State & Actions
const showOpnameModal = ref(false)
const opnameFile = ref<File | null>(null)
const opnamePreview = ref<Array<{ drugId: number; name: string; batch: string; system: number; actual: number }>>([])
const isImportingOpname = ref(false)
const isDownloadingOpname = ref(false)
const opnameSuccessMessage = ref('')
const opnameErrorMessage = ref('')

async function handleDownloadOpname() {
  isDownloadingOpname.value = true
  try {
    const unitParam = filterUnit.value === 'ALL' ? '' : filterUnit.value
    const blob = await api.exportPharmacyOpnameCSV(unitParam)
    const dateStr = new Date().toISOString().split('T')[0]
    const unitName = filterUnit.value === 'ALL' ? 'semua-unit' : filterUnit.value.toLowerCase()
    triggerBlobDownload(blob, `stok_opnam_${unitName}_${dateStr}.csv`)
  } catch (err: any) {
    alert(err.message || 'Gagal mengunduh lembar stok opnam CSV')
  } finally {
    isDownloadingOpname.value = false
  }
}

function onOpnameFileSelected(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  const file = target.files[0]
  if (!file.name.endsWith('.csv')) {
    alert('Harap pilih file dengan format .csv')
    return
  }
  opnameFile.value = file
  opnameSuccessMessage.value = ''
  opnameErrorMessage.value = ''

  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target?.result as string
    if (!text) return
    parseOpnamePreview(text)
  }
  reader.readAsText(file)
}

function parseOpnamePreview(text: string) {
  const lines = text.split('\n')
  const preview = []
  for (let i = 1; i < lines.length; i++) {
    const line = lines[i].trim()
    if (!line) continue
    const parts = line.split(',')
    if (parts.length >= 8) {
      const drugId = parseInt(parts[0]) || 0
      const name = parts[1] || ''
      const batch = parts[5] || '-'
      const system = parseInt(parts[7]) || 0
      const actual = parts[8] ? parseInt(parts[8]) : system
      preview.push({ drugId, name, batch, system, actual })
    }
  }
  opnamePreview.value = preview
}

async function handleExecuteOpname() {
  if (!opnameFile.value) return
  isImportingOpname.value = true
  opnameErrorMessage.value = ''

  const reader = new FileReader()
  reader.onload = async (e) => {
    const content = e.target?.result as string
    try {
      const res = await api.importPharmacyOpnameCSV(content)
      clinicStore.playChime()
      opnameSuccessMessage.value = res.data?.message || res.message || 'Stok fisik berhasil direkonsiliasi secara massal!'
      opnameFile.value = null
      opnamePreview.value = []
      await fetchPrescriptions()
    } catch (err: any) {
      opnameErrorMessage.value = err.message || 'Gagal sinkronisasi CSV opnam'
    } finally {
      isImportingOpname.value = false
    }
  }
  reader.readAsText(opnameFile.value)
}

function closeOpnameModal() {
  showOpnameModal.value = false
  opnameFile.value = null
  opnamePreview.value = []
  opnameSuccessMessage.value = ''
  opnameErrorMessage.value = ''
}

// Counts
const countReserved = computed(() => prescriptions.value.filter((p) => p.status === 'RESERVED').length)
const countPreparing = computed(() => prescriptions.value.filter((p) => p.status === 'PREPARING').length)
const countDispensed = computed(() => prescriptions.value.filter((p) => p.status === 'DISPENSED').length)

const totalRxValue = computed(() => {
  return prescriptions.value.reduce((total, rx) => {
    const rxSum = rx.items?.reduce((sum, item) => sum + (item.subtotal || 0), 0) || 0
    return total + rxSum
  }, 0)
})

const filteredPrescriptions = computed(() => {
  let list = prescriptions.value
  if (filterUnit.value !== 'ALL') {
    list = list.filter((p) => p.clinic_unit === filterUnit.value)
  }
  if (statusFilter.value !== 'ALL') {
    list = list.filter((p) => p.status === statusFilter.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(
      (p) =>
        p.patient_name.toLowerCase().includes(q) ||
        p.queue_number.toLowerCase().includes(q) ||
        p.doctor_name.toLowerCase().includes(q)
    )
  }
  return list
})

function selectPrescription(rx: PharmacyPrescription) {
  selectedRx.value = rx
}

onMounted(() => {
  fetchPrescriptions()
})
</script>

<template>
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- TOP BAR: MINIMALIST CLINICAL HEADER -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Instalasi Farmasi & Apotek
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Prescription Dispensing & Drug Verification
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <div class="text-right hidden sm:block">
          <div class="text-[10px] uppercase font-bold text-zinc-400">Apoteker Bertugas</div>
          <div class="text-xs font-bold text-zinc-900">{{ authStore.user?.name || 'Apt. Dimas Wicaksono, S.Farm' }}</div>
        </div>

        <button
          @click="handleDownloadOpname"
          :disabled="isDownloadingOpname"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
          title="Unduh Lembar Kerja Saldo Stok CSV"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span class="hidden md:inline">Unduh Template Opnam</span>
          <span class="md:hidden">Unduh</span>
        </button>

        <button
          @click="showOpnameModal = true"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-purple-200 bg-purple-50 hover:bg-purple-100 text-xs font-bold text-purple-900 transition-all active:scale-95 cursor-pointer shadow-2xs"
          title="Unggah Hasil Hitung Fisik Lemari Farmasi"
        >
          <Upload class="w-3.5 h-3.5 text-purple-700" />
          <span class="hidden md:inline">Unggah CSV Opnam</span>
          <span class="md:hidden">Unggah</span>
        </button>

        <button
          @click="fetchPrescriptions"
          :disabled="loading"
          title="Segarkan daftar resep masuk"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', loading && 'animate-spin']" />
          <span>Sync</span>
        </button>
      </div>
    </header>

    <!-- 1-ROW COMPACT OPERATIONAL TICKER STRIP (Hairline Grid) -->
    <div
      class="grid grid-cols-2 sm:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Menunggu Racik</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-amber-700 mt-0.5">
          {{ countReserved }} Resep
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Sedang Diracik</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-blue-700 mt-0.5">
          {{ countPreparing }} Resep
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Telah Diserahkan</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-emerald-700 mt-0.5">
          {{ countDispensed }} Resep
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Total Resep Hari Ini</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-700 mt-0.5">
          {{ prescriptions.length }} Pasien
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex items-center justify-between">
        <div>
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Total Valuasi Obat</div>
          <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
            {{ formatRupiah(totalRxValue) }}
          </div>
        </div>
        <Pill class="w-4 h-4 text-zinc-400" />
      </div>
    </div>

    <!-- MAIN MASTER-DETAIL WORKSPACE (2 PANES) -->
    <main class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      
      <!-- LEFT PANE: INCOMING PRESCRIPTIONS FEED (5 cols) -->
      <section class="lg:col-span-5 space-y-3">
        <!-- Control Bar: Unit Filter & Status Filter -->
        <div class="space-y-2">
          <div class="flex items-center justify-between gap-2">
            <!-- Unit Filter Pills -->
            <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
              <button
                type="button"
                @click="filterUnit = 'ALL'; fetchPrescriptions()"
                :class="[
                  'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                  filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                Semua
              </button>
              <button
                type="button"
                @click="filterUnit = 'KIDINEST'; fetchPrescriptions()"
                :class="[
                  'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                  filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                Kidinest
              </button>
              <button
                type="button"
                @click="filterUnit = 'SSC'; fetchPrescriptions()"
                :class="[
                  'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                  filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                SSC
              </button>
            </div>

            <!-- Status Filter -->
            <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
              <button
                type="button"
                @click="statusFilter = 'ALL'"
                :class="[
                  'px-2 py-1 rounded-lg transition-all cursor-pointer',
                  statusFilter === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                Semua
              </button>
              <button
                type="button"
                @click="statusFilter = 'RESERVED'"
                :class="[
                  'px-2 py-1 rounded-lg transition-all cursor-pointer',
                  statusFilter === 'RESERVED' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                Antre ({{ countReserved }})
              </button>
              <button
                type="button"
                @click="statusFilter = 'PREPARING'"
                :class="[
                  'px-2 py-1 rounded-lg transition-all cursor-pointer',
                  statusFilter === 'PREPARING' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
                ]"
              >
                Racik ({{ countPreparing }})
              </button>
            </div>
          </div>

          <!-- Search Input -->
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari pasien, no. antrean, atau nama dokter..."
              class="w-full pl-3 pr-8 py-2 rounded-xl bg-white border border-zinc-200 text-xs text-zinc-900 placeholder:text-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-950"
            />
            <button
              v-if="searchQuery"
              @click="searchQuery = ''"
              class="absolute right-2.5 top-1/2 -translate-y-1/2 text-zinc-400 hover:text-zinc-700"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        <!-- Feed Container -->
        <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
          <div class="p-3.5 border-b border-zinc-100 bg-zinc-50/60 flex items-center justify-between">
            <div class="flex items-center gap-1.5">
              <span class="text-xs font-bold text-zinc-900">Antrean Masuk Resep</span>
              <span class="text-[10px] font-bold px-1.5 py-0.2 rounded-md bg-zinc-200 text-zinc-800 tabular-nums">
                {{ filteredPrescriptions.length }}
              </span>
            </div>
            <span class="text-[10px] text-zinc-400 font-medium">Auto FIFO</span>
          </div>

          <div v-if="loading" class="py-12 text-center text-xs text-zinc-400">
            <RefreshCw class="w-4 h-4 mx-auto animate-spin mb-2 text-zinc-400" />
            Memuat daftar resep...
          </div>

          <div v-else-if="filteredPrescriptions.length === 0" class="py-12 text-center text-xs text-zinc-400">
            <CheckCircle2 class="w-6 h-6 mx-auto text-emerald-500 mb-1.5 opacity-80" />
            Tidak ada antrean resep dalam filter ini.
          </div>

          <div v-else class="divide-y divide-zinc-100 max-h-[640px] overflow-y-auto">
            <div
              v-for="rx in filteredPrescriptions"
              :key="rx.id"
              @click="selectPrescription(rx)"
              :class="[
                'p-3.5 transition-colors cursor-pointer space-y-2',
                selectedRx?.id === rx.id
                  ? 'bg-zinc-900 text-white'
                  : 'hover:bg-zinc-50/70 text-zinc-900',
              ]"
            >
              <!-- Row Header -->
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-1.5">
                  <span
                    :class="[
                      'font-mono text-xs font-extrabold px-2 py-0.5 rounded-md tabular-nums',
                      selectedRx?.id === rx.id
                        ? 'bg-zinc-800 text-white'
                        : 'bg-zinc-100 text-zinc-900'
                    ]"
                  >
                    {{ rx.queue_number }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-semibold px-1.5 py-0.2 rounded border',
                      rx.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                        : 'bg-rose-50 text-rose-900 border-rose-200/80',
                    ]"
                  >
                    {{ rx.clinic_unit }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-mono',
                      selectedRx?.id === rx.id ? 'text-zinc-400' : 'text-zinc-500'
                    ]"
                  >
                    {{ (rx.items || []).length }} item obat
                  </span>
                </div>

                <!-- Status Badge -->
                <span
                  :class="[
                    'text-[10px] font-bold px-2 py-0.5 rounded-full',
                    rx.status === 'RESERVED'
                      ? 'bg-amber-400 text-zinc-950 animate-pulse'
                      : rx.status === 'PREPARING'
                      ? 'bg-blue-400 text-zinc-950'
                      : rx.status === 'DISPENSED'
                      ? 'bg-emerald-400 text-zinc-950'
                      : 'bg-zinc-200 text-zinc-700',
                  ]"
                >
                  {{
                    rx.status === 'RESERVED'
                      ? 'Menunggu Racik'
                      : rx.status === 'PREPARING'
                      ? 'Sedang Diracik'
                      : rx.status === 'DISPENSED'
                      ? 'Diserahkan'
                      : 'Batal'
                  }}
                </span>
              </div>

              <!-- Patient & Doctor -->
              <div>
                <div class="font-bold text-sm tracking-tight leading-snug">{{ rx.patient_name }}</div>
                <div
                  :class="[
                    'text-xs mt-0.5 flex items-center justify-between',
                    selectedRx?.id === rx.id ? 'text-zinc-400' : 'text-zinc-500'
                  ]"
                >
                  <span>{{ rx.doctor_name.split(',')[0] }}</span>
                  <span class="font-mono text-[11px]">{{ formatDate(rx.created_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- RIGHT PANE: ACTIVE DISPENSING & VERIFICATION WORKSTATION (7 cols, Sticky) -->
      <section class="lg:col-span-7 sticky top-6">
        
        <!-- Empty State -->
        <div
          v-if="!selectedRx"
          class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-12 text-center"
        >
          <Pill class="w-10 h-10 mx-auto text-zinc-300 mb-3" />
          <h3 class="text-sm font-bold text-zinc-800">Pilih Resep untuk Mulai Menyiapkan Obat</h3>
          <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
            Klik nomor resep pada panel antrean sebelah kiri untuk memeriksa rincian obat, telaah etiket, dan menyerahkan obat ke pasien.
          </p>
        </div>

        <!-- Selected Prescription Card -->
        <div v-else class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-5 sm:p-6 space-y-5">
          
          <!-- Prescription Header -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-zinc-100">
            <div>
              <div class="flex items-center gap-2">
                <span class="px-2.5 py-0.5 rounded-md font-mono text-xs font-black bg-zinc-950 text-white">
                  {{ selectedRx.queue_number }}
                </span>
                <span class="text-xs font-semibold px-2 py-0.5 rounded border border-zinc-200 bg-zinc-50 text-zinc-800">
                  {{ selectedRx.clinic_unit }}
                </span>
                <span class="text-xs text-zinc-400 font-mono">
                  ID: #RX-{{ selectedRx.id }}
                </span>
              </div>
              <h2 class="text-xl font-bold text-zinc-950 tracking-tight mt-1.5">
                {{ selectedRx.patient_name }}
              </h2>
              <div class="text-xs text-zinc-500 mt-0.5">
                Peresep: {{ selectedRx.doctor_name }} • {{ formatDate(selectedRx.created_at) }}
              </div>
            </div>

            <!-- Status Tag -->
            <span
              :class="[
                'text-xs font-bold px-3 py-1 rounded-full self-start sm:self-auto',
                selectedRx.status === 'RESERVED'
                  ? 'bg-amber-100 text-amber-950 border border-amber-200'
                  : selectedRx.status === 'PREPARING'
                  ? 'bg-blue-100 text-blue-950 border border-blue-200'
                  : selectedRx.status === 'DISPENSED'
                  ? 'bg-emerald-100 text-emerald-950 border border-emerald-200'
                  : 'bg-rose-100 text-rose-950 border border-rose-200',
              ]"
            >
              {{
                selectedRx.status === 'RESERVED'
                  ? 'Menunggu Racik (Stok Terkunci)'
                  : selectedRx.status === 'PREPARING'
                  ? 'Sedang Diracik di Meja Farmasi'
                  : selectedRx.status === 'DISPENSED'
                  ? 'Obat Telah Diserahkan ke Pasien'
                  : 'Dibatalkan'
              }}
            </span>
          </div>

          <!-- 7 Benar Obat Verification Checklist (Clinical Ergonomics) -->
          <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 space-y-2">
            <div class="flex items-center justify-between">
              <span class="text-[11px] font-bold uppercase tracking-wider text-zinc-700 flex items-center gap-1.5">
                <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
                <span>Protokol Telaah Resep (7 Benar Obat Farmasi)</span>
              </span>
              <span class="text-[10px] font-medium text-emerald-700">Verifikasi Aman</span>
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-1.5 text-[11px] text-zinc-600">
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Pasien</span>
              </div>
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Obat</span>
              </div>
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Dosis</span>
              </div>
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Rute</span>
              </div>
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Waktu</span>
              </div>
              <div class="flex items-center gap-1.5">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Informasi</span>
              </div>
              <div class="flex items-center gap-1.5 col-span-2">
                <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                <span>Benar Dokumentasi & FEFO</span>
              </div>
            </div>
          </div>

          <!-- Items Breakdown Table -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <h3 class="text-xs font-bold uppercase tracking-wider text-zinc-950">
                Daftar Obat & Aturan Pakai (Etiket)
              </h3>
              <span class="text-[11px] font-mono text-zinc-400">
                Total: {{ (selectedRx.items || []).length }} Formula
              </span>
            </div>

            <div class="border border-zinc-200 rounded-xl divide-y divide-zinc-100 overflow-hidden bg-white">
              <div
                v-for="it in selectedRx.items"
                :key="it.id"
                class="p-3.5 space-y-1.5 hover:bg-zinc-50/50 transition-colors"
              >
                <div class="flex items-center justify-between text-xs">
                  <div class="flex items-center gap-2">
                    <span
                      v-if="it.is_manual"
                      class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-amber-100 text-amber-950"
                    >
                      Racikan
                    </span>
                    <span class="font-bold text-zinc-950">{{ it.item_name }}</span>
                  </div>
                  <div class="font-bold tabular-nums text-zinc-900">
                    {{ it.quantity }}x &bull; {{ formatRupiah(it.subtotal) }}
                  </div>
                </div>

                <div class="text-[11px] text-zinc-600 bg-zinc-50 p-2 rounded-lg border border-zinc-100 flex items-center justify-between">
                  <div>
                    <span class="font-semibold text-zinc-700">Signa / Aturan: </span>
                    <span class="italic">{{ it.dosage_instruction || 'Sesuai anjuran dokter' }}</span>
                  </div>
                  <div class="tabular-nums text-zinc-400 text-[10px] font-medium">
                    @ {{ formatRupiah(it.unit_price) }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Total Price Summary Strip -->
            <div class="p-3 rounded-xl border border-zinc-200 bg-stone-50/80 flex items-center justify-between text-xs">
              <span class="font-bold text-zinc-700">Total Biaya Obat:</span>
              <span class="font-extrabold tabular-nums text-sm text-zinc-950">
                {{ formatRupiah(selectedRx.items?.reduce((acc, it) => acc + (it.subtotal || 0), 0)) }}
              </span>
            </div>
          </div>

          <!-- Dispensing Actions -->
          <div class="pt-3 border-t border-zinc-100 flex flex-col sm:flex-row items-center justify-between gap-2.5">
            <!-- Cancel Button -->
            <button
              v-if="selectedRx.status === 'RESERVED' || selectedRx.status === 'PREPARING'"
              type="button"
              @click="cancel(selectedRx.id)"
              :disabled="processingId === selectedRx.id"
              class="px-3.5 py-2 rounded-xl text-rose-700 hover:bg-rose-50 border border-rose-200 text-xs font-bold transition-all disabled:opacity-50 cursor-pointer w-full sm:w-auto"
            >
              Batalkan Resep
            </button>
            <div v-else></div>

            <div class="flex items-center gap-2 w-full sm:w-auto justify-end">
              <!-- Preparing Step -->
              <button
                v-if="selectedRx.status === 'RESERVED'"
                type="button"
                @click="startPreparing(selectedRx.id)"
                :disabled="processingId === selectedRx.id"
                class="px-4 py-2.5 rounded-xl border border-blue-200 bg-blue-50 hover:bg-blue-100 text-blue-900 text-xs font-bold transition-all active:scale-95 disabled:opacity-50 flex items-center gap-1.5 shadow-2xs cursor-pointer w-full sm:w-auto justify-center"
              >
                <Sparkles class="w-3.5 h-3.5 text-blue-700" />
                <span>Mulai Racik</span>
              </button>

              <!-- Dispense Step -->
              <button
                v-if="selectedRx.status === 'PREPARING' || selectedRx.status === 'RESERVED'"
                type="button"
                @click="dispense(selectedRx.id)"
                :disabled="processingId === selectedRx.id"
                class="px-5 py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition-all active:scale-95 disabled:opacity-50 flex items-center gap-1.5 shadow-2xs cursor-pointer w-full sm:w-auto justify-center"
              >
                <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
                <span>Selesai & Serahkan Obat</span>
              </button>

              <!-- Go to Billing (if completed) -->
              <button
                v-if="selectedRx.status === 'DISPENSED'"
                type="button"
                @click="goToBilling(selectedRx.queue_id)"
                class="px-4 py-2.5 rounded-xl border border-zinc-200 bg-zinc-100 hover:bg-zinc-200 text-zinc-900 text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer w-full sm:w-auto justify-center shadow-2xs"
              >
                <span>Buka Meja Kasir</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Modal: Unggah & Rekonsiliasi CSV Stok Opnam (Refined #cashier standard) -->
    <div
      v-if="showOpnameModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs animate-in fade-in duration-150"
    >
      <div class="bg-white w-full max-w-xl rounded-3xl p-6 shadow-2xl border border-zinc-200/80 space-y-4">
        <div class="flex items-start justify-between pb-3 border-b border-zinc-100">
          <div>
            <div class="flex items-center gap-2">
              <span class="px-2 py-0.5 rounded-md text-[10px] font-bold bg-purple-100 text-purple-900">
                Audit Stok Fisik Farmasi
              </span>
              <span class="text-xs text-zinc-400 font-mono">FEFO & Reserved Safe</span>
            </div>
            <h2 class="text-lg font-bold text-zinc-950 mt-1">Unggah Hasil Stok Opnam (CSV)</h2>
          </div>
          <button
            type="button"
            @click="closeOpnameModal"
            class="p-1 rounded-full text-zinc-400 hover:text-zinc-700 hover:bg-zinc-100 transition-colors cursor-pointer"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="text-xs text-zinc-600 leading-relaxed bg-zinc-50 p-3 rounded-xl border border-zinc-200">
          💡 <strong>Petunjuk:</strong> Unduh template CSV saldo terlebih dahulu, catat hasil hitung lemari nyata pada kolom <code>Stok_Fisik_Nyata</code>, lalu unggah kembali ke sini.
          <div class="text-[11px] text-emerald-800 font-medium mt-1">
            ✓ Invarian Aman: Stok obat yang sedang dipesan dokter (reserved) tidak akan hilang saldonya.
          </div>
        </div>

        <!-- File Picker Dropzone -->
        <div class="space-y-1.5">
          <label class="block text-[11px] font-bold text-zinc-700 uppercase tracking-wider">File Spreadsheet (.CSV)</label>
          <div class="relative flex items-center justify-center border-2 border-dashed border-zinc-300 hover:border-zinc-400 rounded-2xl p-6 transition-colors bg-white">
            <input
              type="file"
              accept=".csv"
              @change="onOpnameFileSelected"
              class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            />
            <div class="text-center space-y-1">
              <FileSpreadsheet class="w-7 h-7 mx-auto text-zinc-600" />
              <div v-if="opnameFile" class="text-xs font-bold text-zinc-900">
                {{ opnameFile.name }} ({{ (opnameFile.size / 1024).toFixed(1) }} KB)
              </div>
              <div v-else class="text-xs text-zinc-500">
                <span class="font-bold text-zinc-950">Klik untuk memilih file</span> atau tarik file CSV ke sini
              </div>
            </div>
          </div>
        </div>

        <!-- Preview Rows -->
        <div v-if="opnamePreview.length > 0" class="space-y-2">
          <div class="flex items-center justify-between text-xs font-bold text-zinc-700">
            <span>Preview Rekonsiliasi (Total: {{ opnamePreview.length }} item)</span>
            <span class="text-[11px] font-normal text-zinc-400">5 baris teratas</span>
          </div>
          <div class="max-h-36 overflow-y-auto border border-zinc-200 rounded-xl text-[11px] divide-y divide-zinc-100">
            <div
              v-for="(item, idx) in opnamePreview.slice(0, 5)"
              :key="idx"
              class="p-2.5 flex items-center justify-between bg-white"
            >
              <div class="min-w-0 pr-2">
                <div class="font-bold text-zinc-900 truncate">{{ item.name }}</div>
                <div class="text-zinc-400 text-[10px]">Batch: {{ item.batch }}</div>
              </div>
              <div class="text-right shrink-0">
                <div class="font-mono">
                  Sistem: <span class="font-bold">{{ item.system }}</span> → Nyata: <span class="font-bold text-zinc-950">{{ item.actual }}</span>
                </div>
                <div
                  :class="[
                    'text-[10px] font-bold font-mono',
                    item.actual === item.system ? 'text-zinc-400' : item.actual > item.system ? 'text-emerald-600' : 'text-rose-600'
                  ]"
                >
                  Selisih: {{ item.actual - item.system > 0 ? '+' : '' }}{{ item.actual - item.system }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Success Notification -->
        <div v-if="opnameSuccessMessage" class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-900 flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span class="font-semibold">{{ opnameSuccessMessage }}</span>
        </div>

        <!-- Error Notification -->
        <div v-if="opnameErrorMessage" class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-xs text-rose-900 flex items-center gap-2">
          <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
          <span class="font-semibold">{{ opnameErrorMessage }}</span>
        </div>

        <!-- Modal Actions -->
        <div class="flex items-center justify-end gap-2.5 pt-3 border-t border-zinc-100">
          <button
            type="button"
            @click="closeOpnameModal"
            class="px-4 py-2 rounded-xl text-xs font-bold text-zinc-600 hover:bg-zinc-100 transition-colors cursor-pointer"
          >
            {{ opnameSuccessMessage ? 'Selesai & Tutup' : 'Batal' }}
          </button>
          <button
            v-if="!opnameSuccessMessage"
            type="button"
            @click="handleExecuteOpname"
            :disabled="!opnameFile || isImportingOpname"
            class="px-4 py-2 rounded-xl text-xs font-bold bg-zinc-950 hover:bg-zinc-800 text-white transition-all shadow-2xs flex items-center gap-2 disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isImportingOpname" class="w-3.5 h-3.5 animate-spin" />
            <CheckCircle2 v-else class="w-3.5 h-3.5" />
            <span>Mulai Sinkronisasi Stok</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
