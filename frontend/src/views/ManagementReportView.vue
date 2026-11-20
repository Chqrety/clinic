<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  TrendingUp,
  BarChart3,
  Calendar,
  CreditCard,
  Receipt,
  Pill,
  Users,
  ArrowDownRight,
  ArrowUpRight,
  AlertCircle,
  Printer,
  FileText,
  CheckCircle2,
  Clock,
  Sparkles,
  RefreshCw,
  Plus,
  Trash2,
  DollarSign,
  Package,
  Layers,
  ChevronRight,
  Search,
  Building2,
  ArrowLeft
} from 'lucide-vue-next'
import { api, type PharmacyDrug, type Invoice, type QueueItem } from '../api/client'

const router = useRouter()

// ==================== REAL DATABASE STATE ====================
const realInvoices = ref<Invoice[]>([])
const realQueues = ref<QueueItem[]>([])
const realDrugs = ref<PharmacyDrug[]>([])
const loadingData = ref(false)

// Helper date generator (default to today or latest invoice date)
function getTodayString(): string {
  const dt = new Date()
  const y = dt.getFullYear()
  const m = String(dt.getMonth() + 1).padStart(2, '0')
  const d = String(dt.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

// ==================== CONTROLS & FILTER ====================
const selectedDate = ref(getTodayString())
const selectedUnit = ref<'KIDINEST' | 'SSC' | 'MINIMO'>('KIDINEST')
const updateTime = ref('23.59')
const groReporterName = ref('Staf Kasir GRO')
const showPatientDetailModal = ref(false)
const searchDispensed = ref('')

const googleReviewNote = ref('')
const shortageAmount = ref(0)
const packageUsage = ref('-')

// ==================== FETCH REAL DATABASE DATA ====================
async function fetchRealDatabaseData() {
  loadingData.value = true
  try {
    const [invRes, queueRes, drugRes] = await Promise.allSettled([
      api.getInvoices(),
      api.getQueues(),
      api.getPharmacyDrugs(selectedUnit.value === 'SSC' ? 'SSC' : 'KIDINEST')
    ])

    if (invRes.status === 'fulfilled' && Array.isArray(invRes.value)) {
      realInvoices.value = invRes.value
      // Jika pada tanggal terpilih belum ada transaksi, tapi di database ada transaksi,
      // kita set tanggal otomatis ke tanggal invoice terbaru agar data riil langsung terlihat
      if (realInvoices.value.length > 0) {
        const hasCurrentDate = realInvoices.value.some(inv => (inv.created_at || '').substring(0, 10) === selectedDate.value)
        if (!hasCurrentDate) {
          const latestInv = realInvoices.value[0]
          if (latestInv && latestInv.created_at) {
            selectedDate.value = latestInv.created_at.substring(0, 10)
          }
        }
      }
    } else {
      realInvoices.value = []
    }

    if (queueRes.status === 'fulfilled' && Array.isArray(queueRes.value)) {
      realQueues.value = queueRes.value
    } else {
      realQueues.value = []
    }

    if (drugRes.status === 'fulfilled' && Array.isArray(drugRes.value)) {
      realDrugs.value = drugRes.value
    } else {
      realDrugs.value = []
    }
  } catch (err) {
    console.warn('Gagal memuat data dari database:', err)
  } finally {
    loadingData.value = false
  }
}

// Filter Invoices by Date & Unit
const dailyInvoices = computed(() => {
  return realInvoices.value.filter(inv => {
    const invDate = (inv.created_at || '').substring(0, 10)
    const matchDate = invDate === selectedDate.value
    const matchUnit = !inv.clinic_unit || inv.clinic_unit.toUpperCase() === selectedUnit.value.toUpperCase()
    return matchDate && matchUnit
  })
})

// Filter Queues by Date & Unit
const dailyQueues = computed(() => {
  return realQueues.value.filter(q => {
    const qDate = (q.created_at || '').substring(0, 10)
    const matchDate = qDate === selectedDate.value
    const matchUnit = !q.clinic_unit || q.clinic_unit.toUpperCase() === selectedUnit.value.toUpperCase()
    return matchDate && matchUnit
  })
})

// Reporter GRO
watch(
  [dailyQueues, selectedUnit],
  ([queues, unit]) => {
    if (queues.length > 0 && queues[0].doctor_name) {
      groReporterName.value = `GRO ${unit}`
    } else {
      groReporterName.value = `GRO ${unit}`
    }
  },
  { immediate: true }
)

// ==================== FINANCIAL CALCULATIONS (100% REAL) ====================
// Omzet Sistem: Total invoice yang telah terbayar (PAID) pada tanggal terpilih
const systemRevenue = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID')
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

// Tagihan yang masih pending (Nominal Menggantung)
const pendingInvoices = computed(() => {
  return dailyInvoices.value.filter(inv => (inv.payment_status || '').toUpperCase() === 'PENDING')
})

const totalPendingSuspense = computed(() => {
  return pendingInvoices.value.reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const pendingItems = computed(() => {
  return pendingInvoices.value.map(inv => ({
    id: String(inv.id),
    amount: -(Number(inv.total_amount) || 0),
    description: `Kekurangan ${inv.invoice_number} (${inv.patient_name})`
  }))
})

// Omzet Manual = Transaksi riil masuk kasir
const manualRevenue = computed(() => {
  return systemRevenue.value
})

const diffRevenue = computed(() => {
  return totalPendingSuspense.value
})

// Rincian Metode Pembayaran (Real Database Breakdown)
const edcBri = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('BRI'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const edcBca = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('BCA'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const edcCimb = computed<number | null>(() => {
  const sum = dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('CIMB'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
  return sum > 0 ? sum : null
})

const totalDebit = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('DEBIT'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const totalCc = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('CC'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const totalQris = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && (inv.payment_method || '').toUpperCase().includes('QRIS'))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const cashTotal = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && ['CASH', 'TUNAI'].includes((inv.payment_method || '').toUpperCase()))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const transferTotal = computed(() => {
  return dailyInvoices.value
    .filter(inv => (inv.payment_status || 'PAID').toUpperCase() === 'PAID' && ['TRANSFER', 'TF', 'BANK_TRANSFER'].includes((inv.payment_method || '').toUpperCase()))
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

const kidinestPoint = ref('-')

// Arus Pasien Riil dari Antrian Database
const totalRealArrival = computed(() => {
  return dailyQueues.value.filter(q => q.queue_status !== 'CANCELLED').length
})

const totalAppointment = computed(() => {
  return dailyQueues.value.filter(q => {
    const c = (q.complaint || '').toLowerCase()
    const num = (q.queue_number || '').toUpperCase()
    return c.includes('booking') || c.includes('reservasi') || num.startsWith('B')
  }).length
})

const totalWalkIn = computed(() => {
  return Math.max(0, totalRealArrival.value - totalAppointment.value)
})

const cancelledQueues = computed(() => {
  return dailyQueues.value.filter(q => q.queue_status === 'CANCELLED')
})

const totalCancelReschedule = computed(() => cancelledQueues.value.length)

const rescheduleList = computed(() => {
  return cancelledQueues.value.map(q => ({
    id: String(q.id),
    targetDate: `${q.patient_name} (${q.queue_number}) - Dibatalkan`
  }))
})

// Obat / Layanan yang Keluar (Riil dari Invoices Database)
const realDispensedItems = computed(() => {
  return dailyInvoices.value.map((inv, idx) => ({
    row: idx + 1,
    patient_name: inv.patient_name || 'Pasien',
    mr: inv.patient_nik || inv.invoice_number,
    item_name: inv.invoice_type === 'PHARMACY' ? 'Obat & Resep Medis' : 'Registrasi & Pemeriksaan Poliklinik',
    val: Number(inv.total_amount) || 0,
    gro: 'GRO ' + selectedUnit.value,
    nurse: '-'
  }))
})

const filteredDispensedItems = computed(() => {
  if (!searchDispensed.value) return realDispensedItems.value
  const q = searchDispensed.value.toLowerCase()
  return realDispensedItems.value.filter(
    it =>
      it.item_name.toLowerCase().includes(q) ||
      it.patient_name.toLowerCase().includes(q) ||
      it.mr.toLowerCase().includes(q) ||
      it.gro.toLowerCase().includes(q)
  )
})

const totalDispensedAmount = computed(() => {
  return realDispensedItems.value.reduce((acc, it) => acc + (it.val || 0), 0)
})

// Grand Total Omzet MTD (Bulan Berjalan Riil dari Database)
const mtdGrandTotal = computed(() => {
  const currentMonth = selectedDate.value.substring(0, 7) // 'YYYY-MM'
  return realInvoices.value
    .filter(inv => (inv.created_at || '').substring(0, 7) === currentMonth && (inv.payment_status || 'PAID').toUpperCase() === 'PAID')
    .reduce((acc, inv) => acc + (Number(inv.total_amount) || 0), 0)
})

// Transaksi Kas Keluar Operasional Kasir (Petty Cash BKK)
interface OperationalOutflow {
  voucherNo: string
  time: string
  category: string
  description: string
  recipient: string
  amount: number
}

const operationalOutflows = ref<OperationalOutflow[]>([])

const totalCashOutflow = computed(() => {
  return operationalOutflows.value.reduce((acc, it) => acc + it.amount, 0)
})

const netDailyCashflow = computed(() => {
  return systemRevenue.value - totalCashOutflow.value
})

// ==================== FORMATTERS ====================
function formatRupiah(val?: number | null) {
  if (val === null || val === undefined) return 'Rp -'
  return 'Rp' + Math.round(val).toLocaleString('id-ID')
}

function formatRupiahPlain(val?: number | null) {
  if (val === null || val === undefined) return '0'
  return Math.round(val).toLocaleString('id-ID')
}

const formattedDayDate = computed(() => {
  try {
    const parts = selectedDate.value.split('-')
    const dt = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]))
    const dayNames = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
    const monthNames = [
      'Januari',
      'Februari',
      'Maret',
      'April',
      'Mei',
      'Juni',
      'Juli',
      'Agustus',
      'September',
      'Oktober',
      'November',
      'Desember'
    ]
    const dayStr = dayNames[dt.getDay()]
    const dateStr = `${dt.getDate()} ${monthNames[dt.getMonth()]} ${dt.getFullYear()}`
    return { dayStr, dateStr, full: `${dayStr}/${dateStr}` }
  } catch {
    return { dayStr: 'Hari Ini', dateStr: selectedDate.value, full: selectedDate.value }
  }
})

function printReport() {
  window.print()
}

onMounted(() => {
  fetchRealDatabaseData()
})

watch(selectedUnit, () => {
  fetchRealDatabaseData()
})
</script>

<template>
  <div class="min-h-screen bg-slate-100/80 text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] pb-24">
    
    <!-- Top Navigation Bar Khusus Dashboard Manajemen -->
    <nav class="bg-white border-b border-zinc-200 sticky top-0 z-30 shadow-2xs px-4 sm:px-8 py-3.5">
      <div class="max-w-[1400px] mx-auto flex flex-wrap items-center justify-between gap-4">
        
        <!-- Brand & Back -->
        <div class="flex items-center gap-3">
          <button
            @click="router.push('/admin')"
            class="p-2 rounded-xl border border-zinc-200 hover:bg-zinc-100 text-zinc-600 transition-colors cursor-pointer"
            title="Kembali ke Workstation Admin / Kasir"
          >
            <ArrowLeft class="w-4 h-4" />
          </button>
          <div>
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
              <span class="text-xs font-black uppercase tracking-wider text-[#0C3B36]">
                {{ selectedUnit === 'SSC' ? 'Semarang Skin Center' : 'Kidinest Pediatric' }}
              </span>
              <span class="text-[10px] px-2 py-0.5 rounded-full font-bold bg-emerald-50 text-emerald-800 border border-emerald-200">
                Dashboard Manajemen Eksekutif
              </span>
            </div>
            <h1 class="text-base sm:text-lg font-black text-zinc-900 leading-tight">
              Laporan Penjualan & Omzet Harian
            </h1>
          </div>
        </div>

        <!-- Filter Controls -->
        <div class="flex flex-wrap items-center gap-2 sm:gap-3">
          
          <!-- Unit Switcher -->
          <div class="flex items-center bg-zinc-100 p-1 rounded-xl text-xs font-bold">
            <button
              @click="selectedUnit = 'KIDINEST'"
              class="px-3 py-1.5 rounded-lg transition-all cursor-pointer"
              :class="selectedUnit === 'KIDINEST' ? 'bg-[#0C3B36] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
            >
              Kidinest
            </button>
            <button
              @click="selectedUnit = 'SSC'"
              class="px-3 py-1.5 rounded-lg transition-all cursor-pointer"
              :class="selectedUnit === 'SSC' ? 'bg-[#9E606D] text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
            >
              SSC
            </button>
            <button
              @click="selectedUnit = 'MINIMO'"
              class="px-3 py-1.5 rounded-lg transition-all cursor-pointer"
              :class="selectedUnit === 'MINIMO' ? 'bg-amber-800 text-white shadow-2xs' : 'text-zinc-600 hover:text-zinc-900'"
            >
              Minimo
            </button>
          </div>

          <!-- Date Picker -->
          <div class="flex items-center gap-1.5 bg-zinc-50 border border-zinc-200 px-3 py-1.5 rounded-xl">
            <Calendar class="w-3.5 h-3.5 text-zinc-500" />
            <input
              type="date"
              v-model="selectedDate"
              class="bg-transparent text-xs font-bold text-zinc-900 focus:outline-hidden cursor-pointer"
            />
          </div>

          <!-- Refresh Database Button -->
          <button
            @click="fetchRealDatabaseData"
            :disabled="loadingData"
            class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 text-xs font-bold shadow-2xs transition-colors cursor-pointer disabled:opacity-50"
            title="Muat Ulang Data Database Riil"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingData }" />
            <span>{{ loadingData ? 'Memuat...' : 'Refresh Database' }}</span>
          </button>

          <button
            @click="printReport"
            class="inline-flex items-center gap-2 px-3.5 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 text-xs font-bold shadow-2xs transition-colors cursor-pointer"
            title="Cetak Laporan Penjualan"
          >
            <Printer class="w-4 h-4" />
            <span>Cetak Laporan</span>
          </button>
        </div>

      </div>
    </nav>

    <!-- Main Container Dashboard -->
    <main class="max-w-[1400px] mx-auto px-4 sm:px-8 pt-6 space-y-6">

      <!-- Status Banner Database (Real-time Feedback) -->
      <div
        v-if="dailyInvoices.length === 0 && dailyQueues.length === 0"
        class="p-4 rounded-2xl bg-amber-50/90 border border-amber-200 text-amber-900 text-xs flex items-center justify-between gap-3 shadow-2xs"
      >
        <div class="flex items-center gap-2.5">
          <AlertCircle class="w-4 h-4 text-amber-600 shrink-0" />
          <span>
            <strong>Database Bersih:</strong> Tidak ada transaksi pada tanggal <strong>{{ formattedDayDate.full }}</strong> untuk unit <strong>{{ selectedUnit }}</strong>. Seluruh angka omzet dan pasien menampilkan nilai riil dari database (0).
          </span>
        </div>
        <button
          @click="fetchRealDatabaseData"
          class="text-xs font-bold text-amber-800 underline hover:text-amber-950 shrink-0 cursor-pointer"
        >
          Cek Ulang Database
        </button>
      </div>

      <!-- ========================================================================= -->
      <!-- DOKUMEN LAPORAN RESMI (TAMPILAN SESUAI FORMAT CONTOH LAPORAN GRO) -->
      <!-- ========================================================================= -->
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-zinc-200 shadow-sm space-y-8">
        
        <!-- Header Dokumen Laporan -->
        <div class="border-b border-zinc-100 pb-6 flex flex-col md:flex-row md:items-start justify-between gap-4">
          <div class="space-y-1">
            <div class="text-xs font-semibold text-zinc-400">Selamat malam, Berikut saya sampaikan,</div>
            <h2 class="text-2xl font-black text-zinc-900 tracking-tight">
              Laporan Penjualan Harian
            </h2>
            <div class="text-base font-black text-[#0C3B36] uppercase tracking-wide">
              {{ selectedUnit === 'SSC' ? 'Semarang Skin Center' : 'Kidinest Pediatric' }}
            </div>
          </div>

          <div class="bg-zinc-50 border border-zinc-200/80 p-4 rounded-2xl text-xs space-y-1.5 md:text-right min-w-[240px]">
            <div>
              <span class="text-zinc-400">Hari / Tanggal :</span>
              <strong class="text-zinc-900 ml-1 font-mono">{{ formattedDayDate.full }}</strong>
            </div>
            <div class="flex items-center md:justify-end gap-1.5">
              <span class="text-zinc-400">Update Jam :</span>
              <input
                type="text"
                v-model="updateTime"
                class="w-16 bg-white px-2 py-0.5 border border-zinc-200 rounded text-xs font-bold text-zinc-900 text-center font-mono"
              />
            </div>
            <div>
              <span class="text-zinc-400">Pelapor GRO :</span>
              <strong class="text-[#0C3B36] ml-1">{{ groReporterName }}</strong>
            </div>
          </div>
        </div>

        <!-- ========================================================================= -->
        <!-- 1. OMZET & REKONSILIASI SISTEM -->
        <!-- ========================================================================= -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold uppercase tracking-wider text-zinc-400 flex items-center gap-2">
              <DollarSign class="w-4 h-4 text-emerald-600" />
              1. Omzet & Rekonsiliasi Sistem
            </h3>
            <span class="text-[11px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-800 border border-emerald-200">
              Modern Rekonsiliasi
            </span>
          </div>

          <!-- 3 Highlight Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            
            <!-- Omzet Manual -->
            <div class="p-5 rounded-2xl bg-zinc-50 border border-zinc-200">
              <div class="text-xs font-bold text-zinc-500 uppercase">Omzet Manual</div>
              <div class="text-2xl font-black text-zinc-900 font-mono mt-1">
                {{ formatRupiah(manualRevenue) }}
              </div>
              <div class="text-[11px] text-zinc-400 mt-1">Pencatatan kasir harian meja GRO</div>
            </div>

            <!-- Omzet Sistem -->
            <div class="p-5 rounded-2xl bg-emerald-50/50 border border-emerald-200">
              <div class="flex items-center justify-between text-xs font-bold text-emerald-900 uppercase">
                <span>Omzet Sistem</span>
                <span class="w-2 h-2 rounded-full bg-emerald-600"></span>
              </div>
              <div class="text-2xl font-black text-[#0C3B36] font-mono mt-1">
                {{ formatRupiah(systemRevenue) }}
              </div>
              <div class="text-[11px] text-emerald-700 mt-1">Dihitung otomatis dari seluruh transaksi sistem</div>
            </div>

            <!-- Selisih -->
            <div class="p-5 rounded-2xl bg-amber-50/50 border border-amber-200">
              <div class="text-xs font-bold text-amber-900 uppercase">
                Selisih Manual dan Sistem
              </div>
              <div class="text-2xl font-black text-amber-700 font-mono mt-1">
                Rp {{ formatRupiahPlain(diffRevenue) }}
              </div>
              <div class="text-[11px] text-amber-800 mt-1">Sesuai total nominal menggantung</div>
            </div>

          </div>

          <!-- Rincian Nominal Menggantung Box -->
          <div class="p-4 rounded-2xl bg-red-50/40 border border-red-200/80 space-y-2">
            <div class="flex items-center justify-between text-xs font-bold text-red-900">
              <span class="flex items-center gap-1.5">
                <AlertCircle class="w-4 h-4 text-red-600" />
                Nominal Menggantung (Pending) :
              </span>
              <span class="font-mono text-red-700">{{ formatRupiah(totalPendingSuspense) }}</span>
            </div>

            <div v-if="pendingItems.length === 0" class="pl-6 text-xs text-zinc-500 italic">
              - (Nihil - Seluruh tagihan lunas atau belum ada transaksi pending)
            </div>
            <div v-else class="space-y-1.5 pl-6 text-xs text-red-800">
              <div v-for="item in pendingItems" :key="item.id" class="flex items-center justify-between">
                <span>
                  {{ item.amount < 0 ? '-' : '+' }}Rp{{ formatRupiahPlain(Math.abs(item.amount)) }} ({{ item.description }})
                </span>
                <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-red-100 text-red-900">Pending</span>
              </div>
            </div>
          </div>
        </div>

        <!-- ========================================================================= -->
        <!-- 2. METODE PEMBAYARAN (RINCIAN TRANSAKSI MASUK) -->
        <!-- ========================================================================= -->
        <div class="space-y-4 pt-2">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold uppercase tracking-wider text-zinc-400 flex items-center gap-2">
              <CreditCard class="w-4 h-4 text-blue-600" />
              2. Metode Pembayaran :
            </h3>
            <span class="text-xs font-bold text-zinc-500 font-mono">
              Total Masuk: {{ formatRupiah(systemRevenue) }}
            </span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 text-xs">
            
            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">EDC BRI</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(edcBri) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">EDC CIMB</span>
              <strong class="font-mono text-zinc-900">{{ edcCimb ? formatRupiah(edcCimb) : 'Rp -' }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">EDC BCA</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(edcBca) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">Total Debit</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(totalDebit) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">Total CC</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(totalCc) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">Total Qris</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiahPlain(totalQris) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">Tunai</span>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(cashTotal) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <div>
                <span class="font-semibold text-zinc-700">Total transfer</span>
                <span class="block text-[10px] text-zinc-400">(payment booking)</span>
              </div>
              <strong class="font-mono text-zinc-900">{{ formatRupiah(transferTotal) }}</strong>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 flex items-center justify-between bg-zinc-50/50">
              <span class="font-semibold text-zinc-700">Kidinest POINT</span>
              <strong class="font-mono text-zinc-500">{{ kidinestPoint }}</strong>
            </div>

          </div>
        </div>

        <!-- ========================================================================= -->
        <!-- 3. DETAIL PASIEN MASUK (KUNJUNGAN & RESCHEDULE) -->
        <!-- ========================================================================= -->
        <div class="space-y-4 pt-2">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-bold uppercase tracking-wider text-zinc-400 flex items-center gap-2">
              <Users class="w-4 h-4 text-purple-600" />
              3. Detail Pasien Masuk & Kunjungan :
            </h3>
            <button
              @click="showPatientDetailModal = !showPatientDetailModal"
              class="text-xs font-bold text-[#0C3B36] hover:underline cursor-pointer"
            >
              {{ showPatientDetailModal ? 'Sembunyikan Rincian Pasien ▲' : 'Lihat ' + realDispensedItems.length + ' Pasien Nyata ▼' }}
            </button>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-5 gap-3 text-xs">
            
            <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50">
              <div class="text-zinc-500">Total Appointment</div>
              <div class="text-xl font-black text-zinc-900 font-mono mt-0.5">{{ totalAppointment }}</div>
            </div>

            <div class="p-3.5 rounded-xl border border-emerald-200 bg-emerald-50/40">
              <div class="text-emerald-800 font-bold">Kedatangan Real</div>
              <div class="text-xl font-black text-[#0C3B36] font-mono mt-0.5">{{ totalRealArrival }}</div>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50">
              <div class="text-zinc-500">Walk in</div>
              <div class="text-xl font-black text-zinc-900 font-mono mt-0.5">{{ totalWalkIn }}</div>
            </div>

            <div class="p-3.5 rounded-xl border border-amber-200 bg-amber-50/40">
              <div class="text-amber-900 font-bold">Cancel / Reschedule</div>
              <div class="text-xl font-black text-amber-800 font-mono mt-0.5">{{ totalCancelReschedule }}</div>
            </div>

            <div class="p-3.5 rounded-xl border border-zinc-200 bg-zinc-50 col-span-2 sm:col-span-1">
              <div class="text-zinc-500">Pemakaian Paket</div>
              <div class="text-xl font-black text-zinc-500 font-mono mt-0.5">{{ packageUsage }}</div>
            </div>

          </div>

          <!-- Reschedule Breakdown Box -->
          <div class="p-4 rounded-xl bg-amber-50/50 border border-amber-200 text-xs space-y-1">
            <div class="font-bold text-amber-950">Rincian Cancel / Reschedule :</div>
            <div v-if="rescheduleList.length === 0" class="text-zinc-500 italic pl-4">
              * (Nihil - Tidak ada antrian batal / reschedule pada tanggal ini)
            </div>
            <div v-else v-for="r in rescheduleList" :key="r.id" class="text-amber-900 pl-4">
              * {{ r.targetDate }}
            </div>
          </div>

          <!-- Collapsible Real Patient List Table -->
          <div v-if="showPatientDetailModal" class="border border-zinc-200 rounded-2xl overflow-hidden animate-in fade-in">
            <div class="p-3 bg-zinc-50 border-b border-zinc-200 flex items-center justify-between text-xs">
              <span class="font-bold text-zinc-800">Daftar Pasien Nyata yang Masuk Hari Ini ({{ realDispensedItems.length }} Item)</span>
              <div class="relative w-56">
                <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-2.5 top-2" />
                <input
                  type="text"
                  v-model="searchDispensed"
                  placeholder="Cari pasien / obat..."
                  class="w-full pl-8 pr-2 py-1 bg-white border border-zinc-200 rounded-lg text-xs"
                />
              </div>
            </div>

            <div class="max-h-[300px] overflow-y-auto">
              <table class="w-full text-left text-xs">
                <thead class="bg-zinc-50 text-zinc-400 uppercase text-[10px] sticky top-0">
                  <tr class="border-b border-zinc-200">
                    <th class="py-2 px-3">No MR</th>
                    <th class="py-2 px-3">Nama Pasien</th>
                    <th class="py-2 px-3">Layanan / Obat / Vaksin</th>
                    <th class="py-2 px-3">GRO</th>
                    <th class="py-2 px-3 text-right">Nilai</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-zinc-100 text-zinc-700">
                  <tr v-if="filteredDispensedItems.length === 0">
                    <td colspan="5" class="py-6 text-center text-zinc-400 italic">
                      Belum ada transaksi pasien pada tanggal ini (Database bersih).
                    </td>
                  </tr>
                  <tr v-for="(it, idx) in filteredDispensedItems" :key="idx" class="hover:bg-zinc-50/60">
                    <td class="py-2 px-3 font-mono text-zinc-400">{{ it.mr || '-' }}</td>
                    <td class="py-2 px-3 font-bold text-zinc-900">{{ it.patient_name }}</td>
                    <td class="py-2 px-3">{{ it.item_name }}</td>
                    <td class="py-2 px-3 text-zinc-500">{{ it.gro || '-' }}</td>
                    <td class="py-2 px-3 text-right font-mono font-bold text-[#0C3B36]">{{ formatRupiah(it.val) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- ========================================================================= -->
        <!-- 4. GRAND TOTAL OMZET & EVALUASI -->
        <!-- ========================================================================= -->
        <div class="space-y-4 pt-2">
          <h3 class="text-sm font-bold uppercase tracking-wider text-zinc-400 flex items-center gap-2">
            <CheckCircle2 class="w-4 h-4 text-emerald-600" />
            4. Grand Total Omzet & Penutup Laporan :
          </h3>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            
            <div class="p-5 rounded-2xl bg-zinc-900 text-white space-y-1">
              <div class="text-xs text-zinc-400 uppercase font-bold">Grand Total Omzet (Bulan Berjalan)</div>
              <div class="text-2xl sm:text-3xl font-black font-mono text-emerald-400">
                {{ formatRupiah(mtdGrandTotal) }}
              </div>
              <div class="text-xs text-zinc-400 pt-1">
                Kekurangan Omzet : <strong class="text-white">{{ shortageAmount }}</strong>
              </div>
            </div>

            <div class="p-5 rounded-2xl bg-zinc-50 border border-zinc-200 space-y-2 text-xs">
              <div class="font-bold text-zinc-900">Info Google Review :</div>
              <div class="text-zinc-600">
                Total Update hari ini : <span class="font-bold text-emerald-700">{{ googleReviewNote || 'Belum ada review baru hari ini' }}</span>
              </div>
              <div class="pt-3 border-t border-zinc-200 flex items-center justify-between text-zinc-500">
                <span>Tanda Tangan Pelapor :</span>
                <strong class="text-[#0C3B36]">{{ groReporterName }}</strong>
              </div>
            </div>

          </div>
        </div>

      </div>

      <!-- ========================================================================= -->
      <!-- BAGIAN DETAIL: OBAT YANG DIBELI DAN KELUAR (MUTASI FARMASI) -->
      <!-- ========================================================================= -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        
        <!-- Panel 1: Obat yang Keluar (Resep & Penjualan Pasien Hari Ini) -->
        <div class="bg-white rounded-3xl p-6 border border-zinc-200 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-black text-zinc-900 flex items-center gap-2">
                <Pill class="w-4 h-4 text-red-600" />
                Obat yang Keluar (Resep Pasien Hari Ini)
              </h4>
              <p class="text-[11px] text-zinc-500">Ditebus pasien untuk vaksinasi, sirup, dan racikan medis</p>
            </div>
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-red-100 text-red-800">
              {{ realDispensedItems.length }} Item Keluar
            </span>
          </div>

          <div class="max-h-[360px] overflow-y-auto border border-zinc-100 rounded-2xl">
            <table class="w-full text-left text-xs">
              <thead class="bg-zinc-50 text-zinc-400 uppercase text-[10px] sticky top-0">
                <tr class="border-b border-zinc-200">
                  <th class="py-2.5 px-3">Nama Obat / Tindakan</th>
                  <th class="py-2.5 px-3">Pasien</th>
                  <th class="py-2.5 px-3 text-right">Nilai</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-zinc-100 text-zinc-700">
                <tr v-if="realDispensedItems.length === 0">
                  <td colspan="3" class="py-6 text-center text-zinc-400 italic">
                    Belum ada resep / obat keluar pada tanggal ini.
                  </td>
                </tr>
                <tr v-for="(it, idx) in realDispensedItems.slice(0, 50)" :key="idx" class="hover:bg-zinc-50/50">
                  <td class="py-2 px-3 font-bold text-zinc-900">
                    {{ it.item_name }}
                  </td>
                  <td class="py-2 px-3 text-zinc-600 text-[11px]">{{ it.patient_name }}</td>
                  <td class="py-2 px-3 text-right font-mono font-bold text-zinc-900">{{ formatRupiah(it.val) }}</td>
                </tr>
              </tbody>
              <tfoot class="bg-zinc-50 font-bold border-t border-zinc-200">
                <tr>
                  <td colspan="2" class="py-2.5 px-3">Total Nilai Obat & Tindakan Keluar</td>
                  <td class="py-2.5 px-3 text-right font-mono text-[#0C3B36]">{{ formatRupiah(totalDispensedAmount) }}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>

        <!-- Panel 2: Obat yang Dibeli (Restock Masuk dari PBF / Supplier) -->
        <div class="bg-white rounded-3xl p-6 border border-zinc-200 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-black text-zinc-900 flex items-center gap-2">
                <Package class="w-4 h-4 text-emerald-600" />
                Obat yang Dibeli (Restock Masuk PBF)
              </h4>
              <p class="text-[11px] text-zinc-500">Master database obat riil dari distributor farmasi</p>
            </div>
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
              {{ realDrugs.length }} Obat Terdaftar
            </span>
          </div>

          <div class="max-h-[360px] overflow-y-auto border border-zinc-100 rounded-2xl">
            <table class="w-full text-left text-xs">
              <thead class="bg-zinc-50 text-zinc-400 uppercase text-[10px] sticky top-0">
                <tr class="border-b border-zinc-200">
                  <th class="py-2.5 px-3">Nama Obat</th>
                  <th class="py-2.5 px-3">Batch & Exp</th>
                  <th class="py-2.5 px-3 text-center">Stok</th>
                  <th class="py-2.5 px-3 text-right">Harga Beli PBF</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-zinc-100 text-zinc-700">
                <tr v-if="realDrugs.length === 0">
                  <td colspan="4" class="py-6 text-center text-zinc-400 italic">
                    Master obat belum tersedia di database.
                  </td>
                </tr>
                <tr v-for="d in realDrugs.slice(0, 50)" :key="d.id" class="hover:bg-zinc-50/50">
                  <td class="py-2 px-3 font-bold text-zinc-900">
                    {{ d.name }}
                    <span class="block text-[10px] text-zinc-400 font-normal">{{ d.kfa_code || d.category }}</span>
                  </td>
                  <td class="py-2 px-3 font-mono text-[11px] text-zinc-600">
                    <span v-if="d.batches && d.batches.length > 0">
                      {{ d.batches[0].batch_number }}
                    </span>
                    <span v-else class="text-zinc-400">-</span>
                  </td>
                  <td class="py-2 px-3 text-center font-bold font-mono">{{ d.physical_stock }}</td>
                  <td class="py-2 px-3 text-right font-mono text-zinc-800">{{ formatRupiah(d.buy_price || d.sell_price) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

      </div>

      <!-- ========================================================================= -->
      <!-- BAGIAN DETAIL: TRANSAKSI MASUK DAN KELUAR (ARUS KAS OPERASIONAL) -->
      <!-- ========================================================================= -->
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-zinc-200 shadow-sm space-y-6">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="text-base font-black text-zinc-900 flex items-center gap-2">
              <Receipt class="w-4 h-4 text-emerald-600" />
              Arus Kas Operasional Kasir (Transaksi Masuk vs Keluar)
            </h4>
            <p class="text-xs text-zinc-500">Penerimaan kas harian dikurangi pengeluaran operasional kas kecil kasir</p>
          </div>
          <span class="text-xs font-bold text-emerald-800 bg-emerald-50 px-3 py-1 rounded-full border border-emerald-200">
            Kas Bersih: {{ formatRupiah(netDailyCashflow) }}
          </span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="p-4 rounded-2xl bg-emerald-50/50 border border-emerald-200">
            <div class="text-[11px] font-bold text-emerald-800 uppercase flex items-center gap-1">
              <ArrowUpRight class="w-3.5 h-3.5" /> Total Transaksi Masuk
            </div>
            <div class="text-2xl font-black text-[#0C3B36] font-mono mt-1">{{ formatRupiah(systemRevenue) }}</div>
            <div class="text-[10px] text-emerald-700 mt-1">{{ dailyInvoices.length }} transaksi pasien & tindakan</div>
          </div>

          <div class="p-4 rounded-2xl bg-red-50/50 border border-red-200">
            <div class="text-[11px] font-bold text-red-800 uppercase flex items-center gap-1">
              <ArrowDownRight class="w-3.5 h-3.5" /> Total Transaksi Keluar (BKK)
            </div>
            <div class="text-2xl font-black text-red-600 font-mono mt-1">{{ formatRupiah(totalCashOutflow) }}</div>
            <div class="text-[10px] text-red-700 mt-1">Kas kecil operasional harian meja kasir</div>
          </div>

          <div class="p-4 rounded-2xl bg-zinc-900 text-white">
            <div class="text-[11px] font-bold text-zinc-400 uppercase flex items-center gap-1">
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" /> Sisa Kas Bersih Harian
            </div>
            <div class="text-2xl font-black text-emerald-400 font-mono mt-1">{{ formatRupiah(netDailyCashflow) }}</div>
            <div class="text-[10px] text-zinc-400 mt-1">Siap disetor / di-settle kasir</div>
          </div>
        </div>

        <!-- Rincian Voucher Pengeluaran Kasir -->
        <div class="border border-zinc-100 rounded-2xl overflow-hidden">
          <table class="w-full text-left text-xs">
            <thead class="bg-zinc-50 text-zinc-400 uppercase text-[10px]">
              <tr class="border-b border-zinc-200">
                <th class="py-2.5 px-3">No Voucher</th>
                <th class="py-2.5 px-3">Jam</th>
                <th class="py-2.5 px-3">Kategori Beban</th>
                <th class="py-2.5 px-3">Keterangan</th>
                <th class="py-2.5 px-3 text-right">Nominal</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100 text-zinc-700">
              <tr v-if="operationalOutflows.length === 0">
                <td colspan="5" class="py-6 text-center text-zinc-400 italic">
                  Tidak ada transaksi pengeluaran kas operasional kasir (BKK) pada tanggal ini.
                </td>
              </tr>
              <tr v-for="it in operationalOutflows" :key="it.voucherNo">
                <td class="py-2.5 px-3 font-mono font-bold text-zinc-900">{{ it.voucherNo }}</td>
                <td class="py-2.5 px-3 text-zinc-400 font-mono">{{ it.time }}</td>
                <td class="py-2.5 px-3 font-bold text-zinc-800">{{ it.category }}</td>
                <td class="py-2.5 px-3 text-zinc-600">{{ it.description }}</td>
                <td class="py-2.5 px-3 text-right font-mono font-bold text-red-600">-{{ formatRupiah(it.amount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </main>

  </div>
</template>

<style scoped>
@media print {
  nav {
    display: none !important;
  }
  body, main {
    background: white !important;
    padding: 0 !important;
  }
}
</style>
