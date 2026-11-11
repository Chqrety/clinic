<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  api,
  type PharmacyDrug,
  type PharmacyOpnameLog,
} from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useClinicStore } from '../stores/clinic'
import {
  Package,
  AlertTriangle,
  AlertCircle,
  Plus,
  RefreshCw,
  CheckCircle2,
  Download,
  Upload,
  Calendar,
  DollarSign,
  Search,
  FileSpreadsheet,
  X,
  FileText,
  Clock,
  ShieldCheck,
  ChevronDown,
  ChevronUp,
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

type InvTab = 'catalog' | 'opname' | 'audit'
const activeTab = ref<InvTab>('catalog')

const drugs = ref<PharmacyDrug[]>([])
const opnameLogs = ref<PharmacyOpnameLog[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const expandedDrugIds = ref<number[]>([])

// Create Drug Modal
const showCreateModal = ref(false)
const newDrugForm = ref({
  clinic_unit: 'KIDINEST' as 'KIDINEST' | 'SSC' | 'ALL',
  name: '',
  kfa_code: '',
  unit: 'Botol',
  category: 'Antipiretik Anak',
  min_stock_alert: 10,
  buy_price: 20000,
  sell_price: 35000,
})
const isSubmittingCreate = ref(false)

// Restock Batch Modal
const showRestockModal = ref(false)
const selectedDrugForRestock = ref<PharmacyDrug | null>(null)
const restockForm = ref({
  batch_number: '',
  expired_date: '',
  quantity: 50,
  buy_price: 0,
})
const isSubmittingRestock = ref(false)

// Opname Import State
const csvFile = ref<File | null>(null)
const csvPreviewLines = ref<
  Array<{
    drug_id: number
    name: string
    batch: string
    system: number
    actual: number
    discrepancy: number
    notes: string
  }>
>([])
const isImporting = ref(false)
const importSuccessMsg = ref('')

function formatRupiah(val?: number) {
  return 'Rp ' + (val || 0).toLocaleString('id-ID')
}

function formatDate(val?: string) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return val
  }
}

// Fetch Data
async function loadDrugs() {
  loading.value = true
  try {
    drugs.value = await api.getPharmacyDrugs(filterUnit.value)
  } catch (err: any) {
    console.error('Gagal memuat obat:', err)
  } finally {
    loading.value = false
  }
}

async function loadOpnameLogs() {
  try {
    opnameLogs.value = await api.getPharmacyOpnameLogs()
  } catch (err: any) {
    console.error('Gagal memuat log opnam:', err)
  }
}

// Filtered Drugs
const filteredDrugs = computed(() => {
  return drugs.value.filter((d) => {
    const matchUnit = filterUnit.value === 'ALL' || d.clinic_unit === filterUnit.value || d.clinic_unit === 'ALL'
    const q = searchQuery.value.toLowerCase().trim()
    const matchQuery = !q || d.name.toLowerCase().includes(q) || d.category.toLowerCase().includes(q)
    return matchUnit && matchQuery
  })
})

// Metrics
const totalSKU = computed(() => drugs.value.length)
const totalAssetValue = computed(() => {
  return drugs.value.reduce((acc, d) => acc + (d.physical_stock || 0) * (d.buy_price || 0), 0)
})
const totalExpiringSoon = computed(() => drugs.value.filter((d) => d.has_expiring_soon).length)
const totalLowStock = computed(() => drugs.value.filter((d) => d.is_low_stock).length)

function toggleExpand(id: number) {
  const idx = expandedDrugIds.value.indexOf(id)
  if (idx >= 0) {
    expandedDrugIds.value.splice(idx, 1)
  } else {
    expandedDrugIds.value.push(id)
  }
}

// Action: Create Drug
function openCreateModal() {
  const randomSuffix = Math.floor(1000 + Math.random() * 9000)
  newDrugForm.value = {
    clinic_unit: filterUnit.value === 'SSC' ? 'SSC' : 'KIDINEST',
    name: '',
    kfa_code: `KFA-${randomSuffix}`,
    unit: 'Botol',
    category: filterUnit.value === 'SSC' ? 'Skincare & Estetika' : 'Antipiretik Anak',
    min_stock_alert: 10,
    buy_price: 25000,
    sell_price: 45000,
  }
  showCreateModal.value = true
}

async function submitCreateDrug() {
  if (!newDrugForm.value.name.trim()) {
    alert('Nama obat wajib diisi!')
    return
  }

  isSubmittingCreate.value = true
  try {
    await api.createPharmacyDrug(newDrugForm.value)
    clinicStore.playChime()
    showCreateModal.value = false
    await loadDrugs()
  } catch (err: any) {
    alert(err.message || 'Gagal membuat jenis obat baru')
  } finally {
    isSubmittingCreate.value = false
  }
}

// Action: Restock Batch
function openRestockModal(drug: PharmacyDrug) {
  selectedDrugForRestock.value = drug
  const yearNext = new Date().getFullYear() + 2
  restockForm.value = {
    batch_number: `BCH-${drug.name.substring(0, 3).toUpperCase()}-${Math.floor(10 + Math.random() * 90)}`,
    expired_date: `${yearNext}-12-31`,
    quantity: 50,
    buy_price: drug.buy_price || 20000,
  }
  showRestockModal.value = true
}

async function submitRestock() {
  if (!selectedDrugForRestock.value) return
  if (!restockForm.value.batch_number.trim() || !restockForm.value.expired_date.trim()) {
    alert('Nomor batch dan tanggal kadaluarsa wajib diisi!')
    return
  }

  isSubmittingRestock.value = true
  try {
    await api.restockPharmacyDrug(selectedDrugForRestock.value.id, restockForm.value)
    clinicStore.playChime()
    showRestockModal.value = false
    await loadDrugs()
  } catch (err: any) {
    alert(err.message || 'Gagal restock batch obat')
  } finally {
    isSubmittingRestock.value = false
  }
}

// Action: Download CSV Opname
function downloadOpnameCSV() {
  const token = localStorage.getItem('clinic_token')
  const unitParam = filterUnit.value === 'ALL' ? '' : filterUnit.value
  window.open(`/api/pharmacy/opname/export?unit=${unitParam}`, '_blank')
}

// Action: Handle File Upload & Parse CSV Preview
function onFileSelected(event: Event) {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length === 0) return
  csvFile.value = target.files[0]

  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target?.result as string
    if (!text) return
    parsePreviewCSV(text)
  }
  reader.readAsText(csvFile.value)
}

function parsePreviewCSV(text: string) {
  const lines = text.split('\n')
  const preview = []

  for (let i = 1; i < lines.length; i++) {
    const line = lines[i].trim()
    if (!line) continue
    const parts = line.split(',')
    if (parts.length >= 8) {
      const drug_id = parseInt(parts[0]) || 0
      const name = parts[1] || ''
      const batch = parts[5] || '-'
      const system = parseInt(parts[7]) || 0
      const actual = parts[8] ? parseInt(parts[8]) : system
      const notes = parts[9] || ''
      preview.push({
        drug_id,
        name,
        batch,
        system,
        actual,
        discrepancy: actual - system,
        notes,
      })
    }
  }
  csvPreviewLines.value = preview
}

async function submitImportCSV() {
  if (!csvFile.value) return
  isImporting.value = true
  importSuccessMsg.value = ''

  try {
    const reader = new FileReader()
    reader.onload = async (e) => {
      const content = e.target?.result as string
      try {
        const res = await api.importPharmacyOpnameCSV(content)
        clinicStore.playChime()
        importSuccessMsg.value = res.message || 'Stok fisik berhasil direkonsiliasi secara massal!'
        csvFile.value = null
        csvPreviewLines.value = []
        await Promise.all([loadDrugs(), loadOpnameLogs()])
      } catch (err: any) {
        alert(err.message || 'Gagal rekonsiliasi CSV')
      } finally {
        isImporting.value = false
      }
    }
    reader.readAsText(csvFile.value)
  } catch (err: any) {
    alert(err.message || 'Gagal memproses file')
    isImporting.value = false
  }
}

onMounted(() => {
  loadDrugs()
  loadOpnameLogs()
})
</script>

<template>
  <div class="py-8 px-4 sm:px-8 max-w-7xl mx-auto space-y-6">
    <!-- Top Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-black/10">
      <div>
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-purple-50 text-purple-900 border border-purple-200">
            <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse"></span>
            Gudang Farmasi & Stok Opnam
          </span>
          <span class="text-xs text-zinc-400 font-mono">SCMS Inventory v2.0</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-zinc-950 tracking-tight mt-1">
          Inventori Obat & Rekonsiliasi Excel
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Manajemen stok fisik FEFO, deteksi batch kadaluarsa segera, dan ekspor/impor massal stok opnam via file CSV/Excel.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-sm transition-all active:scale-95"
        >
          <Plus class="w-3.5 h-3.5" />
          Tambah Obat Baru
        </button>

        <button
          @click="loadDrugs(); loadOpnameLogs()"
          :disabled="loading"
          class="p-2.5 rounded-2xl bg-white border border-black/10 hover:border-black/20 text-zinc-700 transition-all active:scale-95 shadow-sm"
        >
          <RefreshCw :class="['w-4 h-4', loading && 'animate-spin']" />
        </button>
      </div>
    </div>

    <!-- Quick Stat Cards -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Total SKU Obat</span>
          <Package class="w-4 h-4 text-zinc-400" />
        </div>
        <div class="text-2xl font-black text-zinc-950 mt-1 font-mono">{{ totalSKU }} Jenis</div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Katalog aktif apotek</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Nilai Aset Fisik</span>
          <DollarSign class="w-4 h-4 text-emerald-600" />
        </div>
        <div class="text-xl font-black text-emerald-600 mt-1 font-mono">
          {{ formatRupiah(totalAssetValue) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Berdasarkan harga beli faktur</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Expired &lt; 90 Hari</span>
          <AlertTriangle class="w-4 h-4 text-rose-600" />
        </div>
        <div class="text-2xl font-black text-rose-600 mt-1 font-mono">{{ totalExpiringSoon }} Batch</div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Peringatan FEFO kritis</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Stok Menipis</span>
          <AlertCircle class="w-4 h-4 text-amber-500" />
        </div>
        <div class="text-2xl font-black text-amber-500 mt-1 font-mono">{{ totalLowStock }} Obat</div>
        <div class="text-[10px] text-zinc-400 mt-0.5">&lt; batas minimum alert</div>
      </div>
    </div>

    <!-- Navigation Pills -->
    <div class="flex items-center gap-2 border-b border-black/5 pb-2">
      <button
        @click="activeTab = 'catalog'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'catalog'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <Package class="w-4 h-4" />
        Katalog & Batch FEFO
      </button>

      <button
        @click="activeTab = 'opname'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'opname'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <FileSpreadsheet class="w-4 h-4" />
        Rekonsiliasi Stok Opnam (Excel/CSV)
      </button>

      <button
        @click="activeTab = 'audit'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'audit'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <Clock class="w-4 h-4" />
        Riwayat Audit Selisih
      </button>
    </div>

    <!-- TAB 1: MASTER OBAT & BATCH FEFO -->
    <div v-if="activeTab === 'catalog'" class="bg-white rounded-3xl p-6 sm:p-7 border border-black/5 shadow-sm space-y-4">
      <!-- Search & Unit Filter Controls -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-black/5">
        <div class="relative flex-1 max-w-sm">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3.5 top-3" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama obat atau kategori..."
            class="w-full pl-9 pr-4 py-2 rounded-xl bg-zinc-50 border border-black/10 text-xs focus:ring-2 focus:ring-zinc-950"
          />
        </div>

        <div class="flex items-center bg-zinc-100 rounded-full p-1 text-xs font-bold">
          <button
            @click="filterUnit = 'ALL'; loadDrugs()"
            :class="['px-3 py-1 rounded-full transition-all', filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
          >
            Semua
          </button>
          <button
            @click="filterUnit = 'KIDINEST'; loadDrugs()"
            :class="['px-3 py-1 rounded-full transition-all', filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
          >
            Kidinest
          </button>
          <button
            @click="filterUnit = 'SSC'; loadDrugs()"
            :class="['px-3 py-1 rounded-full transition-all', filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
          >
            SSC
          </button>
        </div>
      </div>

      <!-- Drugs Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="border-b border-black/5 text-zinc-400 font-semibold uppercase tracking-wider text-[10px]">
              <th class="py-2.5">Nama Obat / SKU</th>
              <th class="py-2.5">Unit</th>
              <th class="py-2.5">Harga Beli / Jual</th>
              <th class="py-2.5">Stok Fisik</th>
              <th class="py-2.5">Stok Dipesan</th>
              <th class="py-2.5">Stok Tersedia</th>
              <th class="py-2.5">Status FEFO</th>
              <th class="py-2.5 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-black/5">
            <template v-for="drug in filteredDrugs" :key="drug.id">
              <tr class="hover:bg-zinc-50/80 transition-colors">
                <td class="py-3.5">
                  <div class="font-extrabold text-zinc-950 flex items-center gap-1.5">
                    <span>{{ drug.name }}</span>
                    <button
                      v-if="drug.batches && drug.batches.length > 0"
                      @click="toggleExpand(drug.id)"
                      class="text-zinc-400 hover:text-zinc-700 p-0.5"
                    >
                      <ChevronDown v-if="!expandedDrugIds.includes(drug.id)" class="w-3.5 h-3.5" />
                      <ChevronUp v-else class="w-3.5 h-3.5" />
                    </button>
                  </div>
                  <div class="text-[11px] text-zinc-400 font-mono">
                    {{ drug.kfa_code || 'Tanpa KFA' }} • {{ drug.category }} • {{ drug.unit }}
                  </div>
                </td>

                <td class="py-3.5">
                  <span
                    :class="[
                      'px-2 py-0.5 rounded-full text-[10px] font-bold',
                      drug.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-100 text-amber-900'
                        : drug.clinic_unit === 'SSC'
                        ? 'bg-rose-100 text-rose-900'
                        : 'bg-zinc-100 text-zinc-800',
                    ]"
                  >
                    {{ drug.clinic_unit }}
                  </span>
                </td>

                <td class="py-3.5 font-mono text-zinc-700">
                  <div class="font-bold">{{ formatRupiah(drug.sell_price) }}</div>
                  <div class="text-[10px] text-zinc-400">Modal: {{ formatRupiah(drug.buy_price) }}</div>
                </td>

                <td class="py-3.5 font-mono font-bold text-zinc-900">
                  {{ drug.physical_stock }} {{ drug.unit }}
                </td>

                <td class="py-3.5 font-mono font-bold text-amber-600">
                  {{ drug.reserved_stock }} {{ drug.unit }}
                </td>

                <td class="py-3.5 font-mono font-black text-emerald-600 text-sm">
                  {{ drug.available_stock }} {{ drug.unit }}
                </td>

                <td class="py-3.5">
                  <span
                    v-if="drug.has_expiring_soon"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-rose-100 text-rose-900 flex items-center gap-1 w-fit"
                  >
                    <AlertTriangle class="w-3 h-3" /> ED &lt; 90 Hari
                  </span>
                  <span
                    v-else-if="drug.is_low_stock"
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 text-amber-900 flex items-center gap-1 w-fit"
                  >
                    <AlertCircle class="w-3 h-3" /> Stok Kritis
                  </span>
                  <span
                    v-else
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-900 flex items-center gap-1 w-fit"
                  >
                    <CheckCircle2 class="w-3 h-3" /> Aman
                  </span>
                </td>

                <td class="py-3.5 text-right">
                  <button
                    @click="openRestockModal(drug)"
                    class="inline-flex items-center gap-1 px-3 py-1 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-zinc-800 text-[11px] font-bold transition-all shadow-xs"
                  >
                    <Plus class="w-3 h-3 text-emerald-600" />
                    Restock
                  </button>
                </td>
              </tr>

              <!-- Accordion: Sub-row for Drug Batches (FEFO Tracking) -->
              <tr v-if="expandedDrugIds.includes(drug.id)" class="bg-zinc-50/60">
                <td colspan="8" class="p-3 pl-6">
                  <div class="space-y-1.5">
                    <div class="text-[10px] font-extrabold uppercase tracking-wider text-zinc-400">
                      Rincian Batch FEFO Masuk:
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
                      <div
                        v-for="b in drug.batches"
                        :key="b.id"
                        class="p-2.5 rounded-xl bg-white border border-black/5 text-[11px] space-y-1"
                      >
                        <div class="flex items-center justify-between font-bold">
                          <span class="font-mono text-zinc-900">{{ b.batch_number }}</span>
                          <span
                            :class="[
                              'text-[9px] px-1.5 py-0.2 rounded',
                              b.is_expiring_soon ? 'bg-rose-100 text-rose-900 font-bold' : 'text-zinc-500',
                            ]"
                          >
                            ED: {{ b.expired_date }}
                          </span>
                        </div>
                        <div class="text-[10px] text-zinc-500 flex justify-between font-mono">
                          <span>Fisik: {{ b.physical_stock }}</span>
                          <span>Dipesan: {{ b.reserved_stock }}</span>
                          <span class="font-bold text-emerald-600">Sisa: {{ b.available_stock }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 2: REKONSILIASI STOK OPNAM (EXCEL/CSV) -->
    <div v-if="activeTab === 'opname'" class="max-w-4xl mx-auto space-y-6">
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-sm space-y-6">
        <div>
          <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-blue-50 text-blue-900 border border-blue-200">
            Fitur 4.4 & 4.5 • Anti-Selisih Fisik
          </span>
          <h2 class="text-xl sm:text-2xl font-black text-zinc-950 mt-1">
            Ekspor & Impor Massal Stok Opnam Fisik
          </h2>
          <p class="text-xs text-zinc-500 mt-1">
            Download lembar kerja rekonsiliasi lemari apotek ke file spreadsheet CSV, catat hasil hitung fisik nyata, lalu unggah kembali untuk penyesuaian massal.
          </p>
        </div>

        <!-- 3-Step Visual Workflow -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div class="p-4 rounded-2xl bg-zinc-50 border border-black/5 space-y-2">
            <div class="w-6 h-6 rounded-full bg-zinc-950 text-white flex items-center justify-center text-xs font-black">
              1
            </div>
            <div class="font-extrabold text-xs text-zinc-900">Unduh Format Lembar Audit</div>
            <p class="text-[11px] text-zinc-400">Download data seluruh batch apotek lengkap dengan saldo sistem.</p>
            <button
              @click="downloadOpnameCSV"
              class="w-full py-2 px-3 rounded-xl bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 transition-all flex items-center justify-center gap-1.5 shadow-xs"
            >
              <Download class="w-3.5 h-3.5 text-blue-600" />
              Download CSV
            </button>
          </div>

          <div class="p-4 rounded-2xl bg-zinc-50 border border-black/5 space-y-2">
            <div class="w-6 h-6 rounded-full bg-zinc-950 text-white flex items-center justify-center text-xs font-black">
              2
            </div>
            <div class="font-extrabold text-xs text-zinc-900">Hitung Fisik Lemari Apotek</div>
            <p class="text-[11px] text-zinc-400">
              Isi kolom <code class="bg-zinc-200 px-1 rounded text-zinc-900 font-mono text-[10px]">Stok_Fisik_Nyata</code> pada file spreadsheet dengan hasil hitung nyata.
            </p>
            <div class="text-[10px] font-mono text-zinc-400 pt-2">Kolom format: CSV UTF-8</div>
          </div>

          <div class="p-4 rounded-2xl bg-zinc-50 border border-black/5 space-y-2">
            <div class="w-6 h-6 rounded-full bg-zinc-950 text-white flex items-center justify-center text-xs font-black">
              3
            </div>
            <div class="font-extrabold text-xs text-zinc-900">Unggah Hasil Hitung Fisik</div>
            <p class="text-[11px] text-zinc-400">Sistem otomatis menghitung selisih dan memperbarui saldo riil apotek.</p>
            <label class="w-full py-2 px-3 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-xs font-bold text-white transition-all flex items-center justify-center gap-1.5 shadow-xs cursor-pointer">
              <Upload class="w-3.5 h-3.5 text-emerald-400" />
              Pilih File CSV
              <input type="file" accept=".csv" @change="onFileSelected" class="hidden" />
            </label>
          </div>
        </div>

        <!-- CSV Preview Table before applying -->
        <div v-if="csvPreviewLines.length > 0" class="space-y-4 pt-2">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-extrabold uppercase tracking-wider text-zinc-950">
              Pratinjau Hasil Perbandingan Selisih ({{ csvPreviewLines.length }} Item)
            </h3>
            <span class="text-[11px] text-zinc-400 font-mono">File: {{ csvFile?.name }}</span>
          </div>

          <div class="overflow-x-auto max-h-60 border border-black/5 rounded-2xl">
            <table class="w-full text-left text-xs">
              <thead class="bg-zinc-50 text-[10px] uppercase font-bold text-zinc-400 sticky top-0">
                <tr>
                  <th class="p-2.5">Obat</th>
                  <th class="p-2.5">Batch</th>
                  <th class="p-2.5">Stok Sistem</th>
                  <th class="p-2.5">Stok Fisik</th>
                  <th class="p-2.5">Selisih</th>
                  <th class="p-2.5">Catatan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-black/5">
                <tr v-for="(it, idx) in csvPreviewLines" :key="idx" class="hover:bg-zinc-50">
                  <td class="p-2.5 font-bold text-zinc-900">{{ it.name }}</td>
                  <td class="p-2.5 font-mono text-zinc-600">{{ it.batch }}</td>
                  <td class="p-2.5 font-mono">{{ it.system }}</td>
                  <td class="p-2.5 font-mono font-bold text-zinc-900">{{ it.actual }}</td>
                  <td class="p-2.5 font-mono font-bold">
                    <span
                      :class="[
                        'px-2 py-0.5 rounded text-[10px]',
                        it.discrepancy < 0
                          ? 'bg-rose-100 text-rose-900'
                          : it.discrepancy > 0
                          ? 'bg-blue-100 text-blue-900'
                          : 'text-zinc-500',
                      ]"
                    >
                      {{ it.discrepancy > 0 ? '+' : '' }}{{ it.discrepancy }}
                    </span>
                  </td>
                  <td class="p-2.5 text-zinc-500 italic text-[11px]">{{ it.notes || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <button
            @click="submitImportCSV"
            :disabled="isImporting"
            class="w-full py-3.5 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-xs shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <CheckCircle2 class="w-4 h-4 text-emerald-400" />
            {{ isImporting ? 'Menyesuaikan Saldo Fisik...' : 'Terapkan Penyesuaian Saldo Stok Opnam Massal' }}
          </button>
        </div>

        <div v-if="importSuccessMsg" class="p-4 rounded-2xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-950 flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span>{{ importSuccessMsg }}</span>
        </div>
      </div>
    </div>

    <!-- TAB 3: RIWAYAT AUDIT SELISIH (AUDIT TRAIL) -->
    <div v-if="activeTab === 'audit'" class="bg-white rounded-3xl p-6 sm:p-7 border border-black/5 shadow-sm space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-black/5">
        <div>
          <h2 class="text-base font-extrabold text-zinc-950">Log Audit Rekonsiliasi Stok Opnam</h2>
          <p class="text-xs text-zinc-500 mt-0.5">Catatan riwayat setiap penyesuaian selisih fisik obat oleh petugas apotek.</p>
        </div>

        <button
          @click="loadOpnameLogs"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-zinc-100 text-zinc-700 text-xs font-bold hover:bg-zinc-200 transition-all"
        >
          <RefreshCw class="w-3.5 h-3.5" />
          Segarkan
        </button>
      </div>

      <div v-if="opnameLogs.length === 0" class="py-12 text-center text-xs text-zinc-400">
        Belum ada riwayat audit stok opnam.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="border-b border-black/5 text-zinc-400 font-semibold uppercase tracking-wider text-[10px]">
              <th class="py-2.5">Waktu</th>
              <th class="py-2.5">Nama Obat</th>
              <th class="py-2.5">Batch</th>
              <th class="py-2.5">Stok Sistem</th>
              <th class="py-2.5">Stok Nyata</th>
              <th class="py-2.5">Selisih</th>
              <th class="py-2.5">Pemeriksa</th>
              <th class="py-2.5">Catatan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-black/5">
            <tr v-for="log in opnameLogs" :key="log.id" class="hover:bg-zinc-50">
              <td class="py-3 font-mono text-zinc-500">{{ formatDate(log.created_at) }}</td>
              <td class="py-3 font-bold text-zinc-900">{{ log.drug_name }}</td>
              <td class="py-3 font-mono text-zinc-600">{{ log.batch_number || '-' }}</td>
              <td class="py-3 font-mono">{{ log.system_stock }}</td>
              <td class="py-3 font-mono font-bold text-zinc-900">{{ log.actual_stock }}</td>
              <td class="py-3 font-mono font-bold">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px]',
                    log.discrepancy < 0
                      ? 'bg-rose-100 text-rose-900'
                      : log.discrepancy > 0
                      ? 'bg-blue-100 text-blue-900'
                      : 'text-zinc-500',
                  ]"
                >
                  {{ log.discrepancy > 0 ? '+' : '' }}{{ log.discrepancy }}
                </span>
              </td>
              <td class="py-3 font-bold text-zinc-800">{{ log.adjusted_by }}</td>
              <td class="py-3 text-zinc-500 italic text-[11px]">{{ log.notes || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL 1: TAMBAH OBAT BARU -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 space-y-4 shadow-2xl border border-black/10">
        <div class="flex items-center justify-between pb-3 border-b border-black/10">
          <div>
            <h3 class="text-sm font-extrabold text-zinc-950">Tambah Jenis Obat Baru</h3>
            <p class="text-xs text-zinc-400">Master katalog obat klinik</p>
          </div>
          <button @click="showCreateModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Unit Klinik:</label>
            <select
              v-model="newDrugForm.clinic_unit"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-medium"
            >
              <option value="KIDINEST">Kidinest (Poli Anak & Vaksinasi)</option>
              <option value="SSC">Semarang Skin Center (Estetika Kulit)</option>
              <option value="ALL">Semua Unit (Umum)</option>
            </select>
          </div>

          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Nama Obat & Dosis:</label>
            <input
              v-model="newDrugForm.name"
              type="text"
              placeholder="Contoh: Ibuprofen Sirup 100mg/5ml"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-medium"
            />
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Satuan Kemasan:</label>
              <select
                v-model="newDrugForm.unit"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-medium"
              >
                <option value="Botol">Botol</option>
                <option value="Strip">Strip</option>
                <option value="Tablet">Tablet</option>
                <option value="Tube">Tube</option>
                <option value="Vial">Vial / Ampul</option>
                <option value="Pot">Pot</option>
              </select>
            </div>

            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Batas Min Alert:</label>
              <input
                v-model.number="newDrugForm.min_stock_alert"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono text-right"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Beli / Modal (Rp):</label>
              <input
                v-model.number="newDrugForm.buy_price"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono text-right"
              />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Jual Pasien (Rp):</label>
              <input
                v-model.number="newDrugForm.sell_price"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono text-right font-bold"
              />
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            @click="submitCreateDrug"
            :disabled="isSubmittingCreate"
            class="flex-1 py-3 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-xs shadow-md transition-all active:scale-95 disabled:opacity-50"
          >
            {{ isSubmittingCreate ? 'Menyimpan...' : 'Simpan Obat Baru' }}
          </button>
          <button
            @click="showCreateModal = false"
            class="px-4 py-3 rounded-2xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs"
          >
            Batal
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 2: RESTOCK BATCH BARU -->
    <div
      v-if="showRestockModal && selectedDrugForRestock"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 space-y-4 shadow-2xl border border-black/10">
        <div class="flex items-center justify-between pb-3 border-b border-black/10">
          <div>
            <h3 class="text-sm font-extrabold text-zinc-950">Restock Batch Masuk</h3>
            <p class="text-xs text-zinc-400">{{ selectedDrugForRestock.name }}</p>
          </div>
          <button @click="showRestockModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Nomor Batch Faktur PBF:</label>
            <input
              v-model="restockForm.batch_number"
              type="text"
              placeholder="Contoh: BCH-HEX-202611"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono font-medium"
            />
          </div>

          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Tanggal Kadaluarsa (Expired Date):</label>
            <input
              v-model="restockForm.expired_date"
              type="date"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono"
            />
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Jumlah Unit Masuk:</label>
              <input
                v-model.number="restockForm.quantity"
                type="number"
                min="1"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono text-center font-bold"
              />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Beli Per Unit (Rp):</label>
              <input
                v-model.number="restockForm.buy_price"
                type="number"
                min="0"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 font-mono text-right"
              />
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            @click="submitRestock"
            :disabled="isSubmittingRestock"
            class="flex-1 py-3 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-xs shadow-md transition-all active:scale-95 disabled:opacity-50"
          >
            {{ isSubmittingRestock ? 'Menyimpan...' : 'Konfirmasi Masuk Gudang' }}
          </button>
          <button
            @click="showRestockModal = false"
            class="px-4 py-3 rounded-2xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs"
          >
            Batal
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
