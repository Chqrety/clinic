<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  api,
  triggerBlobDownload,
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
async function downloadOpnameCSV() {
  try {
    const unitParam = filterUnit.value === 'ALL' ? '' : filterUnit.value
    const blob = await api.exportPharmacyOpnameCSV(unitParam)
    const dateStr = new Date().toISOString().split('T')[0]
    const unitName = filterUnit.value === 'ALL' ? 'semua-unit' : filterUnit.value.toLowerCase()
    triggerBlobDownload(blob, `stok_opnam_${unitName}_${dateStr}.csv`)
  } catch (err: any) {
    alert(err.message || 'Gagal mengunduh lembar stok opnam CSV')
  }
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
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- TOP BAR: MINIMALIST CLINICAL HEADER -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Gudang Logistik & Farmasi
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Inventory Catalog & Batch FEFO Management
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="openCreateModal"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition-all active:scale-95 cursor-pointer shadow-2xs"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Tambah Obat Baru</span>
        </button>

        <button
          type="button"
          @click="loadDrugs(); loadOpnameLogs()"
          :disabled="loading"
          title="Segarkan katalog obat dan log audit"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', loading && 'animate-spin']" />
          <span>Sync</span>
        </button>
      </div>
    </header>

    <!-- 1-ROW COMPACT OPERATIONAL TICKER STRIP (Hairline Grid) -->
    <div
      class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Total SKU Obat</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ totalSKU }} Jenis
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Stok Menipis (Alert)</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-amber-700 mt-0.5">
          {{ totalLowStock }} Obat
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Mendekati ED (&lt; 90 Hari)</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-rose-700 mt-0.5">
          {{ totalExpiringSoon }} Batch
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Unit Terpilih</div>
        <div class="text-base font-bold tracking-tight text-zinc-700 mt-0.5">
          {{ filterUnit === 'ALL' ? 'Semua Unit' : filterUnit }}
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex items-center justify-between">
        <div>
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Nilai Aset Fisik Riil</div>
          <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
            {{ formatRupiah(totalAssetValue) }}
          </div>
        </div>
        <Package class="w-4 h-4 text-zinc-400" />
      </div>
    </div>

    <!-- NAVIGATION TABS -->
    <nav class="flex items-center gap-2 border-b border-zinc-200 pb-2">
      <button
        type="button"
        @click="activeTab = 'catalog'"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'catalog'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <Package class="w-3.5 h-3.5" />
        <span>Katalog & Batch FEFO</span>
        <span
          :class="[
            'px-1.5 py-0.2 rounded-md text-[10px] tabular-nums font-bold',
            activeTab === 'catalog' ? 'bg-zinc-800 text-zinc-200' : 'bg-stone-200 text-zinc-800'
          ]"
        >
          {{ totalSKU }}
        </span>
      </button>

      <button
        type="button"
        @click="activeTab = 'opname'"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'opname'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <FileSpreadsheet class="w-3.5 h-3.5" />
        <span>Rekonsiliasi Stok Opnam (CSV)</span>
      </button>

      <button
        type="button"
        @click="activeTab = 'audit'"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'audit'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <Clock class="w-3.5 h-3.5" />
        <span>Riwayat Audit Selisih</span>
        <span
          v-if="opnameLogs.length > 0"
          :class="[
            'px-1.5 py-0.2 rounded-md text-[10px] tabular-nums font-bold',
            activeTab === 'audit' ? 'bg-zinc-800 text-zinc-200' : 'bg-stone-200 text-zinc-800'
          ]"
        >
          {{ opnameLogs.length }}
        </span>
      </button>
    </nav>

    <!-- ==================== TAB 1: MASTER OBAT & BATCH FEFO ==================== -->
    <div v-if="activeTab === 'catalog'" class="space-y-3">
      <!-- Search & Unit Filter Controls -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="relative flex-1 max-w-sm">
          <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari nama obat, KFA, atau kategori..."
            class="w-full pl-8 pr-3 py-1.5 rounded-xl bg-white border border-zinc-200 text-xs text-zinc-900 placeholder:text-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-950"
          />
        </div>

        <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
          <button
            type="button"
            @click="filterUnit = 'ALL'; loadDrugs()"
            :class="[
              'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
              filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Semua
          </button>
          <button
            type="button"
            @click="filterUnit = 'KIDINEST'; loadDrugs()"
            :class="[
              'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
              filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            Kidinest
          </button>
          <button
            type="button"
            @click="filterUnit = 'SSC'; loadDrugs()"
            :class="[
              'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
              filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
            ]"
          >
            SSC
          </button>
        </div>
      </div>

      <!-- High-Density Data Table -->
      <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-zinc-50/70 border-b border-zinc-200/90 text-zinc-500 font-bold uppercase tracking-wider text-[10px]">
              <tr>
                <th class="py-3 px-4">Nama Obat / SKU</th>
                <th class="py-3 px-3">Unit</th>
                <th class="py-3 px-3">Harga Beli / Jual</th>
                <th class="py-3 px-3">Stok Fisik</th>
                <th class="py-3 px-3">Stok Dipesan</th>
                <th class="py-3 px-3">Stok Tersedia</th>
                <th class="py-3 px-3">Status FEFO</th>
                <th class="py-3 px-4 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <template v-for="drug in filteredDrugs" :key="drug.id">
                <tr class="hover:bg-zinc-50/70 transition-colors">
                  <td class="py-3 px-4">
                    <div class="font-bold text-zinc-950 flex items-center gap-1.5">
                      <span>{{ drug.name }}</span>
                      <button
                        v-if="drug.batches && drug.batches.length > 0"
                        type="button"
                        @click="toggleExpand(drug.id)"
                        class="text-zinc-400 hover:text-zinc-700 p-0.5 cursor-pointer"
                        title="Tampilkan rincian batch FEFO"
                      >
                        <ChevronDown v-if="!expandedDrugIds.includes(drug.id)" class="w-3.5 h-3.5" />
                        <ChevronUp v-else class="w-3.5 h-3.5" />
                      </button>
                    </div>
                    <div class="text-[11px] text-zinc-400 font-mono">
                      {{ drug.kfa_code || 'Tanpa KFA' }} • {{ drug.category }} • {{ drug.unit }}
                    </div>
                  </td>

                  <td class="py-3 px-3">
                    <span
                      :class="[
                        'px-1.5 py-0.2 rounded border text-[10px] font-semibold',
                        drug.clinic_unit === 'KIDINEST'
                          ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                          : drug.clinic_unit === 'SSC'
                          ? 'bg-rose-50 text-rose-900 border-rose-200/80'
                          : 'bg-zinc-50 text-zinc-800 border-zinc-200',
                      ]"
                    >
                      {{ drug.clinic_unit }}
                    </span>
                  </td>

                  <td class="py-3 px-3 tabular-nums text-zinc-700">
                    <div class="font-bold text-zinc-900">{{ formatRupiah(drug.sell_price) }}</div>
                    <div class="text-[10px] text-zinc-400">Modal: {{ formatRupiah(drug.buy_price) }}</div>
                  </td>

                  <td class="py-3 px-3 tabular-nums font-bold text-zinc-900">
                    {{ drug.physical_stock }} {{ drug.unit }}
                  </td>

                  <td class="py-3 px-3 tabular-nums font-bold text-amber-700">
                    {{ drug.reserved_stock }} {{ drug.unit }}
                  </td>

                  <td class="py-3 px-3 tabular-nums font-extrabold text-emerald-700 text-sm">
                    {{ drug.available_stock }} {{ drug.unit }}
                  </td>

                  <td class="py-3 px-3">
                    <span
                      v-if="drug.has_expiring_soon"
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-rose-50 text-rose-900 border border-rose-200 flex items-center gap-1 w-fit"
                    >
                      <AlertTriangle class="w-3 h-3 text-rose-600" />
                      <span>ED &lt; 90 Hari</span>
                    </span>
                    <span
                      v-else-if="drug.is_low_stock"
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-50 text-amber-900 border border-amber-200 flex items-center gap-1 w-fit"
                    >
                      <AlertCircle class="w-3 h-3 text-amber-600" />
                      <span>Stok Kritis</span>
                    </span>
                    <span
                      v-else
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-900 border border-emerald-200 flex items-center gap-1 w-fit"
                    >
                      <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                      <span>Aman</span>
                    </span>
                  </td>

                  <td class="py-3 px-4 text-right">
                    <button
                      type="button"
                      @click="openRestockModal(drug)"
                      class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800 text-[11px] font-bold transition-all shadow-2xs cursor-pointer"
                    >
                      <Plus class="w-3 h-3 text-emerald-600" />
                      <span>Restock</span>
                    </button>
                  </td>
                </tr>

                <!-- Accordion: Sub-row for Drug Batches (FEFO Tracking) -->
                <tr v-if="expandedDrugIds.includes(drug.id)" class="bg-zinc-50/50">
                  <td colspan="8" class="p-3 pl-8">
                    <div class="space-y-1.5">
                      <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">
                        Rincian Batch Masuk (FEFO Priority):
                      </div>
                      <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
                        <div
                          v-for="b in drug.batches"
                          :key="b.id"
                          class="p-2.5 rounded-xl bg-white border border-zinc-200 text-[11px] space-y-1 shadow-2xs"
                        >
                          <div class="flex items-center justify-between font-bold">
                            <span class="font-mono text-zinc-900">{{ b.batch_number }}</span>
                            <span
                              :class="[
                                'text-[9px] px-1.5 py-0.2 rounded',
                                b.is_expiring_soon ? 'bg-rose-50 text-rose-900 border border-rose-200 font-bold' : 'text-zinc-500',
                              ]"
                            >
                              ED: {{ b.expired_date }}
                            </span>
                          </div>
                          <div class="text-[10px] text-zinc-500 flex justify-between font-mono">
                            <span>Fisik: {{ b.physical_stock }}</span>
                            <span>Dipesan: {{ b.reserved_stock }}</span>
                            <span class="font-bold text-emerald-700">Tersedia: {{ b.available_stock }}</span>
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
    </div>

    <!-- ==================== TAB 2: REKONSILIASI STOK OPNAM (CSV) ==================== -->
    <div v-if="activeTab === 'opname'" class="max-w-4xl mx-auto space-y-4">
      <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-6 space-y-5">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2 py-0.5 rounded-md text-[10px] font-bold bg-zinc-100 text-zinc-800">
              Audit Stok Fisik
            </span>
            <span class="text-xs text-zinc-400 font-mono">Excel & CSV Reconciler</span>
          </div>
          <h2 class="text-xl font-bold text-zinc-950 tracking-tight mt-1">
            Ekspor & Impor Massal Stok Opnam Fisik
          </h2>
          <p class="text-xs text-zinc-500 mt-0.5">
            Unduh lembar kerja rekonsiliasi lemari apotek ke file spreadsheet CSV, catat hasil hitung fisik nyata, lalu unggah kembali untuk penyesuaian massal.
          </p>
        </div>

        <!-- 3-Step Workflow -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 space-y-2">
            <div class="w-5 h-5 rounded-md bg-zinc-950 text-white flex items-center justify-center text-[10px] font-black">
              1
            </div>
            <div class="font-bold text-xs text-zinc-900">Unduh Format Lembar Audit</div>
            <p class="text-[11px] text-zinc-500 leading-relaxed">Download data seluruh batch apotek lengkap dengan saldo sistem.</p>
            <button
              type="button"
              @click="downloadOpnameCSV"
              class="w-full py-1.5 px-3 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all flex items-center justify-center gap-1.5 shadow-2xs cursor-pointer"
            >
              <Download class="w-3.5 h-3.5 text-zinc-700" />
              <span>Download CSV</span>
            </button>
          </div>

          <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 space-y-2">
            <div class="w-5 h-5 rounded-md bg-zinc-950 text-white flex items-center justify-center text-[10px] font-black">
              2
            </div>
            <div class="font-bold text-xs text-zinc-900">Hitung Fisik Lemari Apotek</div>
            <p class="text-[11px] text-zinc-500 leading-relaxed">
              Isi kolom <code class="bg-zinc-200 px-1 rounded text-zinc-900 font-mono text-[10px]">Stok_Fisik_Nyata</code> pada file spreadsheet dengan hasil hitung nyata.
            </p>
            <div class="text-[10px] font-mono text-zinc-400 pt-1">Format: CSV UTF-8</div>
          </div>

          <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50/70 space-y-2">
            <div class="w-5 h-5 rounded-md bg-zinc-950 text-white flex items-center justify-center text-[10px] font-black">
              3
            </div>
            <div class="font-bold text-xs text-zinc-900">Unggah Hasil Hitung Fisik</div>
            <p class="text-[11px] text-zinc-500 leading-relaxed">Sistem otomatis menghitung selisih dan memperbarui saldo riil apotek.</p>
            <label class="w-full py-1.5 px-3 rounded-lg bg-zinc-950 hover:bg-zinc-800 text-xs font-bold text-white transition-all flex items-center justify-center gap-1.5 shadow-2xs cursor-pointer">
              <Upload class="w-3.5 h-3.5 text-emerald-400" />
              <span>Pilih File CSV</span>
              <input type="file" accept=".csv" @change="onFileSelected" class="hidden" />
            </label>
          </div>
        </div>

        <!-- CSV Preview Table before applying -->
        <div v-if="csvPreviewLines.length > 0" class="space-y-3 pt-2">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-bold uppercase tracking-wider text-zinc-950">
              Pratinjau Hasil Perbandingan Selisih ({{ csvPreviewLines.length }} Item)
            </h3>
            <span class="text-[11px] text-zinc-400 font-mono">File: {{ csvFile?.name }}</span>
          </div>

          <div class="overflow-x-auto max-h-60 border border-zinc-200 rounded-xl">
            <table class="w-full text-left text-xs">
              <thead class="bg-zinc-50 text-[10px] uppercase font-bold text-zinc-500 sticky top-0 border-b border-zinc-200">
                <tr>
                  <th class="p-2.5">Obat</th>
                  <th class="p-2.5">Batch</th>
                  <th class="p-2.5">Stok Sistem</th>
                  <th class="p-2.5">Stok Fisik</th>
                  <th class="p-2.5">Selisih</th>
                  <th class="p-2.5">Catatan</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-zinc-100">
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
                          ? 'bg-rose-50 text-rose-900 border border-rose-200'
                          : it.discrepancy > 0
                          ? 'bg-blue-50 text-blue-900 border border-blue-200'
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
            type="button"
            @click="submitImportCSV"
            :disabled="isImporting"
            class="w-full py-3 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2 cursor-pointer"
          >
            <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
            <span>{{ isImporting ? 'Menyesuaikan Saldo Fisik...' : 'Terapkan Penyesuaian Saldo Stok Opnam Massal' }}</span>
          </button>
        </div>

        <div v-if="importSuccessMsg" class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-950 flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span>{{ importSuccessMsg }}</span>
        </div>
      </div>
    </div>

    <!-- ==================== TAB 3: RIWAYAT AUDIT SELISIH ==================== -->
    <div v-if="activeTab === 'audit'" class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs p-5 sm:p-6 space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
        <div>
          <h2 class="text-sm font-bold text-zinc-950">Log Audit Rekonsiliasi Stok Opnam</h2>
          <p class="text-xs text-zinc-500 mt-0.5">Catatan riwayat setiap penyesuaian selisih fisik obat oleh petugas apotek.</p>
        </div>

        <button
          type="button"
          @click="loadOpnameLogs"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 text-xs font-bold transition-all shadow-2xs cursor-pointer"
        >
          <RefreshCw class="w-3.5 h-3.5" />
          <span>Segarkan</span>
        </button>
      </div>

      <div v-if="opnameLogs.length === 0" class="py-12 text-center text-xs text-zinc-400">
        Belum ada riwayat audit stok opnam.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-zinc-50/70 border-b border-zinc-200/90 text-zinc-500 font-bold uppercase tracking-wider text-[10px]">
            <tr>
              <th class="py-2.5 px-3">Waktu</th>
              <th class="py-2.5 px-3">Nama Obat</th>
              <th class="py-2.5 px-3">Batch</th>
              <th class="py-2.5 px-3">Stok Sistem</th>
              <th class="py-2.5 px-3">Stok Nyata</th>
              <th class="py-2.5 px-3">Selisih</th>
              <th class="py-2.5 px-3">Pemeriksa</th>
              <th class="py-2.5 px-3">Catatan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr v-for="log in opnameLogs" :key="log.id" class="hover:bg-zinc-50">
              <td class="py-3 px-3 font-mono text-zinc-500">{{ formatDate(log.created_at) }}</td>
              <td class="py-3 px-3 font-bold text-zinc-900">{{ log.drug_name }}</td>
              <td class="py-3 px-3 font-mono text-zinc-600">{{ log.batch_number || '-' }}</td>
              <td class="py-3 px-3 font-mono">{{ log.system_stock }}</td>
              <td class="py-3 px-3 font-mono font-bold text-zinc-900">{{ log.actual_stock }}</td>
              <td class="py-3 px-3 font-mono font-bold">
                <span
                  :class="[
                    'px-2 py-0.5 rounded text-[10px]',
                    log.discrepancy < 0
                      ? 'bg-rose-50 text-rose-900 border border-rose-200'
                      : log.discrepancy > 0
                      ? 'bg-blue-50 text-blue-900 border border-blue-200'
                      : 'text-zinc-500',
                  ]"
                >
                  {{ log.discrepancy > 0 ? '+' : '' }}{{ log.discrepancy }}
                </span>
              </td>
              <td class="py-3 px-3 font-bold text-zinc-800">{{ log.adjusted_by }}</td>
              <td class="py-3 px-3 text-zinc-500 italic text-[11px]">{{ log.notes || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL 1: TAMBAH OBAT BARU -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 space-y-4 shadow-2xl border border-zinc-200/80">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">Tambah Jenis Obat Baru</h3>
            <p class="text-xs text-zinc-400">Master katalog obat klinik</p>
          </div>
          <button type="button" @click="showCreateModal = false" class="text-zinc-400 hover:text-zinc-700 cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Unit Klinik:</label>
            <select
              v-model="newDrugForm.clinic_unit"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-medium text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
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
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-medium text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
            />
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Satuan Kemasan:</label>
              <select
                v-model="newDrugForm.unit"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-medium text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
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
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-right text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Beli / Modal (Rp):</label>
              <input
                v-model.number="newDrugForm.buy_price"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-right text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Jual Pasien (Rp):</label>
              <input
                v-model.number="newDrugForm.sell_price"
                type="number"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-right font-bold text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2 border-t border-zinc-100">
          <button
            type="button"
            @click="submitCreateDrug"
            :disabled="isSubmittingCreate"
            class="flex-1 py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
          >
            {{ isSubmittingCreate ? 'Menyimpan...' : 'Simpan Obat Baru' }}
          </button>
          <button
            type="button"
            @click="showCreateModal = false"
            class="px-4 py-2.5 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs cursor-pointer"
          >
            Batal
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 2: RESTOCK BATCH BARU -->
    <div
      v-if="showRestockModal && selectedDrugForRestock"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 space-y-4 shadow-2xl border border-zinc-200/80">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div>
            <h3 class="text-sm font-bold text-zinc-950">Restock Batch Masuk</h3>
            <p class="text-xs text-zinc-400">{{ selectedDrugForRestock.name }}</p>
          </div>
          <button type="button" @click="showRestockModal = false" class="text-zinc-400 hover:text-zinc-700 cursor-pointer">
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
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono font-medium text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
            />
          </div>

          <div class="space-y-1">
            <label class="font-bold text-zinc-700">Tanggal Kadaluarsa (Expired Date):</label>
            <input
              v-model="restockForm.expired_date"
              type="date"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
            />
          </div>

          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Jumlah Unit Masuk:</label>
              <input
                v-model.number="restockForm.quantity"
                type="number"
                min="1"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-center font-bold text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div class="space-y-1">
              <label class="font-bold text-zinc-700">Harga Beli Per Unit (Rp):</label>
              <input
                v-model.number="restockForm.buy_price"
                type="number"
                min="0"
                class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-zinc-200 font-mono text-right text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2 border-t border-zinc-100">
          <button
            type="button"
            @click="submitRestock"
            :disabled="isSubmittingRestock"
            class="flex-1 py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
          >
            {{ isSubmittingRestock ? 'Menyimpan...' : 'Konfirmasi Masuk Gudang' }}
          </button>
          <button
            type="button"
            @click="showRestockModal = false"
            class="px-4 py-2.5 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs cursor-pointer"
          >
            Batal
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
