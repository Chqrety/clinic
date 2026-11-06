<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Prescription } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Pill, CheckCircle2, XCircle, RefreshCw, AlertTriangle, Clock } from 'lucide-vue-next'

const store = useClinicStore()
const prescriptions = ref<Prescription[]>([])
const loading = ref(false)
const processingId = ref<number | null>(null)

async function fetchPrescriptions() {
  loading.value = true
  try {
    prescriptions.value = await api.getPrescriptions()
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

async function dispense(id: number) {
  processingId.value = id
  try {
    await api.dispensePrescription(id)
    store.playChime()
    await fetchPrescriptions()
  } catch (err: any) {
    alert(err.message || 'Gagal memproses penyerahan obat')
  } finally {
    processingId.value = null
  }
}

async function cancel(id: number) {
  if (!confirm('Yakin ingin membatalkan resep ini? Reserved stock akan dikembalikan ke apotek.')) {
    return
  }
  processingId.value = id
  try {
    await api.cancelPrescription(id)
    await fetchPrescriptions()
  } catch (err: any) {
    alert(err.message || 'Gagal membatalkan resep')
  } finally {
    processingId.value = null
  }
}

onMounted(() => {
  fetchPrescriptions()
})
</script>

<template>
  <div class="min-h-screen pt-24 pb-16 px-4 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8">
      <div>
        <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-600 text-white mb-2 shadow-sm">
          Apotek & Farmasi
        </span>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Antrian Racik E-Resep
        </h1>
        <p class="text-sm text-zinc-500 mt-1">
          Resep elektronik masuk otomatis dari dokter tanpa kertas/disposisi manual.
        </p>
      </div>
      <button
        @click="fetchPrescriptions"
        :disabled="loading"
        class="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
        Segarkan Resep
      </button>
    </div>

    <!-- Active Prescription Cards List -->
    <div class="space-y-4">
      <div
        v-for="rx in prescriptions"
        :key="rx.id"
        class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] hover:shadow-[0_8px_24px_rgba(0,0,0,0.05)] transition-all"
      >
        <!-- Card Top Bar -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-zinc-100 gap-2">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-stone-100 flex items-center justify-center text-zinc-800 font-extrabold text-sm">
              #{{ rx.id }}
            </div>
            <div>
              <div class="text-base font-extrabold text-zinc-950">{{ rx.patient_name }}</div>
              <div class="text-xs text-zinc-400">NIK: {{ rx.patient_nik }} • Dokter: {{ rx.doctor_name }}</div>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <span
              :class="[
                'inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-bold',
                rx.status === 'RESERVED' && 'bg-amber-100 text-amber-800 border border-amber-200',
                rx.status === 'DISPENSED' && 'bg-emerald-100 text-emerald-800 border border-emerald-200',
                rx.status === 'CANCELLED' && 'bg-rose-100 text-rose-800'
              ]"
            >
              <span v-if="rx.status === 'RESERVED'" class="w-2 h-2 rounded-full bg-amber-500 animate-ping"></span>
              {{ rx.status === 'RESERVED' ? 'Perlu Disiapkan (Stok Reserved)' : rx.status }}
            </span>
          </div>
        </div>

        <!-- Medicine Items List -->
        <div class="py-4 space-y-2.5">
          <div
            v-for="it in rx.items"
            :key="it.id"
            class="flex items-center justify-between p-3 rounded-2xl bg-[#F9F9F7] text-sm"
          >
            <div class="flex items-center gap-3">
              <Pill class="w-4 h-4 text-emerald-600 flex-shrink-0" />
              <div>
                <span class="font-bold text-zinc-900">{{ it.drug_name }}</span>
                <span class="text-xs text-zinc-400 ml-2">Batch: {{ it.batch_number }}</span>
                <div class="text-xs font-medium text-zinc-500">{{ it.dosage_instruction }}</div>
              </div>
            </div>
            <div class="text-right">
              <div class="font-extrabold text-zinc-950">{{ it.quantity }} item</div>
              <div class="text-xs text-zinc-400">Rp {{ it.subtotal.toLocaleString('id-ID') }}</div>
            </div>
          </div>
        </div>

        <!-- Card Footer Actions -->
        <div class="pt-4 border-t border-zinc-100 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div class="text-sm font-semibold text-zinc-600">
            Total Harga Obat: <strong class="text-zinc-950 font-extrabold">Rp {{ rx.total_amount.toLocaleString('id-ID') }}</strong>
          </div>

          <div v-if="rx.status === 'RESERVED'" class="flex items-center gap-2">
            <button
              @click="cancel(rx.id)"
              :disabled="processingId === rx.id"
              class="px-4 py-2.5 rounded-full border border-black/10 hover:border-black/20 text-xs font-bold text-rose-600 hover:bg-rose-50 transition-all active:scale-95"
            >
              Batalkan
            </button>
            <button
              @click="dispense(rx.id)"
              :disabled="processingId === rx.id"
              class="px-6 py-2.5 rounded-full bg-amber-400 hover:bg-amber-500 text-zinc-950 text-xs font-extrabold transition-all active:scale-95 shadow-sm flex items-center gap-1.5"
            >
              <CheckCircle2 class="w-4 h-4" />
              <span>{{ processingId === rx.id ? 'Memotong Stok...' : 'Selesai & Serahkan Obat' }}</span>
            </button>
          </div>

          <div v-else-if="rx.status === 'DISPENSED'" class="text-xs font-bold text-emerald-700 flex items-center gap-1">
            <CheckCircle2 class="w-4 h-4" /> Obat telah diserahkan, stok fisik apotek terpotong
          </div>
        </div>
      </div>

      <div v-if="prescriptions.length === 0" class="p-12 text-center text-zinc-400 bg-white rounded-3xl border border-black/5 text-sm">
        Belum ada e-resep masuk dari dokter.
      </div>
    </div>
  </div>
</template>
