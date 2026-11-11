<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  api,
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
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

const prescriptions = ref<PharmacyPrescription[]>([])
const loading = ref(false)
const processingId = ref<number | null>(null)
const filterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')

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

// Counts
const countReserved = computed(() => prescriptions.value.filter((p) => p.status === 'RESERVED').length)
const countPreparing = computed(() => prescriptions.value.filter((p) => p.status === 'PREPARING').length)
const countDispensed = computed(() => prescriptions.value.filter((p) => p.status === 'DISPENSED').length)

const filteredPrescriptions = computed(() => {
  if (filterUnit.value === 'ALL') return prescriptions.value
  return prescriptions.value.filter((p) => p.clinic_unit === filterUnit.value)
})

onMounted(() => {
  fetchPrescriptions()
})
</script>

<template>
  <div class="py-8 px-4 sm:px-8 max-w-7xl mx-auto space-y-6">
    <!-- Header & Staff Bar -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-black/10">
      <div>
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-emerald-50 text-emerald-900 border border-emerald-200">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
            Apotek & Farmasi Terpadu
          </span>
          <span class="text-xs text-zinc-400 font-mono">SCMS Rx Engine v2.0</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-zinc-950 tracking-tight mt-1">
          Antrian Racik & Penyerahan Obat
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          E-resep masuk otomatis dari dokter. Penyerahan obat memotong stok fisik riil apotek dan merilis alokasi reserved.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="bg-white rounded-2xl px-4 py-2 border border-black/5 shadow-sm text-right">
          <div class="text-[10px] uppercase font-bold text-zinc-400">Apoteker Bertugas</div>
          <div class="text-xs font-bold text-zinc-900">{{ authStore.user?.name || 'Apt. Dimas Wicaksono, S.Farm' }}</div>
        </div>

        <button
          @click="fetchPrescriptions"
          :disabled="loading"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-2xl bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50"
        >
          <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
          Sync Resep
        </button>
      </div>
    </div>

    <!-- Quick Stats Cards & Unit Filter -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div class="grid grid-cols-3 gap-3 w-full sm:w-auto">
        <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm min-w-[130px]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Menunggu Racik</div>
          <div class="text-2xl font-black text-amber-600 font-mono mt-0.5">{{ countReserved }}</div>
          <div class="text-[10px] text-zinc-400">Stok Terkunci</div>
        </div>

        <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm min-w-[130px]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Sedang Diracik</div>
          <div class="text-2xl font-black text-blue-600 font-mono mt-0.5">{{ countPreparing }}</div>
          <div class="text-[10px] text-zinc-400">Di Meja Farmasi</div>
        </div>

        <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm min-w-[130px]">
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Selesai Diserahkan</div>
          <div class="text-2xl font-black text-emerald-600 font-mono mt-0.5">{{ countDispensed }}</div>
          <div class="text-[10px] text-zinc-400">Fisik Terpotong</div>
        </div>
      </div>

      <!-- Unit Filter Tabs -->
      <div class="flex items-center bg-white border border-black/5 shadow-sm rounded-full p-1 text-xs font-bold self-start sm:self-auto">
        <button
          @click="filterUnit = 'ALL'; fetchPrescriptions()"
          :class="['px-4 py-1.5 rounded-full transition-all', filterUnit === 'ALL' ? 'bg-zinc-950 text-white' : 'text-zinc-600 hover:text-zinc-900']"
        >
          Semua Unit
        </button>
        <button
          @click="filterUnit = 'KIDINEST'; fetchPrescriptions()"
          :class="['px-4 py-1.5 rounded-full transition-all', filterUnit === 'KIDINEST' ? 'bg-amber-500 text-white' : 'text-zinc-600 hover:text-zinc-900']"
        >
          Kidinest
        </button>
        <button
          @click="filterUnit = 'SSC'; fetchPrescriptions()"
          :class="['px-4 py-1.5 rounded-full transition-all', filterUnit === 'SSC' ? 'bg-rose-500 text-white' : 'text-zinc-600 hover:text-zinc-900']"
        >
          Semarang Skin Center
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredPrescriptions.length === 0 && !loading" class="bg-white rounded-3xl p-12 border border-black/5 text-center shadow-sm">
      <Pill class="w-12 h-12 mx-auto text-zinc-300 mb-3" />
      <h3 class="text-base font-extrabold text-zinc-800">Tidak Ada Antrian Resep Aktif</h3>
      <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
        Ketika dokter dari ruang periksa mengirimkan e-resep baru, resep akan langsung masuk ke layar apotek secara real-time.
      </p>
    </div>

    <!-- Prescriptions Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-5">
      <div
        v-for="rx in filteredPrescriptions"
        :key="rx.id"
        class="bg-white rounded-3xl p-6 border border-black/5 shadow-sm space-y-4 flex flex-col justify-between transition-all hover:shadow-md"
      >
        <!-- Top Info Header -->
        <div class="space-y-3">
          <div class="flex items-center justify-between pb-3 border-b border-black/5">
            <div class="flex items-center gap-2">
              <span class="px-3 py-1 rounded-full font-mono text-xs font-black bg-zinc-950 text-white">
                {{ rx.queue_number }}
              </span>
              <span
                :class="[
                  'text-[10px] font-bold px-2.5 py-0.5 rounded-full',
                  rx.clinic_unit === 'KIDINEST' ? 'bg-amber-100 text-amber-900' : 'bg-rose-100 text-rose-900',
                ]"
              >
                {{ rx.clinic_unit }}
              </span>
            </div>

            <!-- Status Badge -->
            <span
              :class="[
                'text-[10px] font-bold px-3 py-1 rounded-full',
                rx.status === 'RESERVED'
                  ? 'bg-amber-100 text-amber-950 animate-pulse'
                  : rx.status === 'PREPARING'
                  ? 'bg-blue-100 text-blue-950'
                  : rx.status === 'DISPENSED'
                  ? 'bg-emerald-100 text-emerald-950'
                  : 'bg-rose-100 text-rose-950',
              ]"
            >
              {{
                rx.status === 'RESERVED'
                  ? 'Menunggu Racik'
                  : rx.status === 'PREPARING'
                  ? 'Sedang Diracik'
                  : rx.status === 'DISPENSED'
                  ? 'Lunas & Diserahkan'
                  : 'Dibatalkan'
              }}
            </span>
          </div>

          <!-- Patient & Doctor Info -->
          <div>
            <h3 class="text-base font-extrabold text-zinc-950 tracking-tight">{{ rx.patient_name }}</h3>
            <p class="text-xs text-zinc-500 flex items-center gap-1.5 mt-0.5">
              <Stethoscope class="w-3.5 h-3.5 text-zinc-400" />
              Dokter Perujuk: {{ rx.doctor_name }} • {{ formatDate(rx.created_at) }}
            </p>
          </div>

          <!-- Items Table -->
          <div class="space-y-2 bg-zinc-50 rounded-2xl p-4 border border-black/5 text-xs">
            <div class="text-[10px] font-extrabold uppercase tracking-wider text-zinc-400">
              Daftar Obat & Aturan Pakai:
            </div>

            <div class="space-y-2 divide-y divide-black/5">
              <div v-for="it in rx.items" :key="it.id" class="pt-2 first:pt-0 space-y-1">
                <div class="flex items-center justify-between font-bold text-zinc-900">
                  <div class="flex items-center gap-1.5">
                    <span
                      v-if="it.is_manual"
                      class="px-1.5 py-0.5 rounded text-[9px] font-bold bg-amber-200 text-amber-950"
                    >
                      Racikan
                    </span>
                    <span>{{ it.item_name }}</span>
                  </div>
                  <span class="font-mono text-zinc-600">{{ it.quantity }}x</span>
                </div>

                <div class="text-[11px] text-zinc-500 font-mono flex items-center justify-between">
                  <span class="italic">Signa: {{ it.dosage_instruction || 'Sesuai anjuran dokter' }}</span>
                  <span>{{ formatRupiah(it.subtotal) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="pt-3 border-t border-black/5 flex items-center justify-between gap-2">
          <!-- Cancel Button -->
          <button
            v-if="rx.status === 'RESERVED' || rx.status === 'PREPARING'"
            @click="cancel(rx.id)"
            :disabled="processingId === rx.id"
            class="px-3 py-2 rounded-xl text-rose-700 hover:bg-rose-50 text-xs font-bold transition-all disabled:opacity-50"
          >
            Batalkan
          </button>

          <div class="flex items-center gap-2 ml-auto">
            <!-- Preparing Step -->
            <button
              v-if="rx.status === 'RESERVED'"
              @click="startPreparing(rx.id)"
              :disabled="processingId === rx.id"
              class="px-4 py-2 rounded-xl bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold transition-all active:scale-95 disabled:opacity-50 flex items-center gap-1.5 shadow-xs"
            >
              <Sparkles class="w-3.5 h-3.5" />
              Mulai Racik
            </button>

            <!-- Dispense Step -->
            <button
              v-if="rx.status === 'PREPARING' || rx.status === 'RESERVED'"
              @click="dispense(rx.id)"
              :disabled="processingId === rx.id"
              class="px-4 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-extrabold transition-all active:scale-95 disabled:opacity-50 flex items-center gap-1.5 shadow-xs"
            >
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
              Selesai & Serahkan Obat
            </button>

            <!-- Go to Billing (if completed) -->
            <button
              v-if="rx.status === 'DISPENSED'"
              @click="goToBilling(rx.queue_id)"
              class="px-4 py-2 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-zinc-900 text-xs font-bold transition-all flex items-center gap-1.5"
            >
              <span>Ke Meja Kasir</span>
              <ArrowRight class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
