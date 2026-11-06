<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type QueueItem } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { CheckCircle2, Clock, Banknote, Users, RefreshCw, Volume2 } from 'lucide-vue-next'

const store = useClinicStore()
const queues = ref<QueueItem[]>([])
const loading = ref(false)
const processingId = ref<number | null>(null)

async function fetchQueues() {
  loading.value = true
  try {
    queues.value = await api.getQueues()
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

async function confirmCash(id: number) {
  processingId.value = id
  try {
    await api.payQueue(id)
    store.playChime()
    await fetchQueues()
  } catch (err: any) {
    alert(err.message || 'Gagal konfirmasi pembayaran')
  } finally {
    processingId.value = null
  }
}

async function updateStatus(id: number, status: string) {
  processingId.value = id
  try {
    await api.updateQueueStatus(id, status)
    store.playChime()
    await fetchQueues()
  } catch (err: any) {
    alert(err.message || 'Gagal mengubah status antrian')
  } finally {
    processingId.value = null
  }
}

onMounted(() => {
  fetchQueues()
})
</script>

<template>
  <div class="py-8 px-6 sm:px-10 max-w-7xl mx-auto">
    <!-- Header & Action Bar -->
    <div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-4 mb-6 pb-6 border-b border-black/10">
      <div>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Pusat Antrian & Pembayaran Loket
        </h1>
        <p class="text-xs text-zinc-500 mt-1">
          Monitor FIFO live antrian, terima setoran tunai, dan panggil pasien ke loket.
        </p>
      </div>
      <button
        @click="fetchQueues"
        :disabled="loading"
        class="inline-flex items-center gap-2 px-5 py-2.5 rounded-full bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
        Segarkan Data
      </button>
    </div>

    <!-- Quick Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] flex items-center justify-between">
        <div>
          <div class="text-xs font-semibold uppercase tracking-wider text-zinc-400">Total Pasien Hari Ini</div>
          <div class="text-3xl font-extrabold text-zinc-950 mt-1">{{ queues.length }}</div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-stone-100 flex items-center justify-center text-zinc-700">
          <Users class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] flex items-center justify-between">
        <div>
          <div class="text-xs font-semibold uppercase tracking-wider text-amber-600">Perlu Bayar Tunai</div>
          <div class="text-3xl font-extrabold text-amber-600 mt-1">
            {{ queues.filter(q => q.payment_status === 'PENDING').length }}
          </div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-amber-50 flex items-center justify-center text-amber-700 border border-amber-200/60">
          <Banknote class="w-5 h-5" />
        </div>
      </div>

      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] flex items-center justify-between">
        <div>
          <div class="text-xs font-semibold uppercase tracking-wider text-emerald-600">Lunas / Antri Aktif</div>
          <div class="text-3xl font-extrabold text-emerald-700 mt-1">
            {{ queues.filter(q => q.payment_status === 'PAID').length }}
          </div>
        </div>
        <div class="w-10 h-10 rounded-2xl bg-emerald-50 flex items-center justify-center text-emerald-700 border border-emerald-200/60">
          <CheckCircle2 class="w-5 h-5" />
        </div>
      </div>
    </div>

    <!-- Section 1: Pending Cash Confirmation Alert -->
    <div v-if="queues.some(q => q.payment_status === 'PENDING')" class="mb-8">
      <div class="flex items-center gap-2 mb-3">
        <span class="w-2.5 h-2.5 rounded-full bg-amber-500 animate-ping"></span>
        <h2 class="text-base font-bold text-zinc-950">Menunggu Pembayaran Tunai di Loket</h2>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div
          v-for="q in queues.filter(q => q.payment_status === 'PENDING')"
          :key="q.id"
          class="bg-amber-50/70 border border-amber-200 rounded-3xl p-5 flex items-center justify-between shadow-sm"
        >
          <div>
            <div class="flex items-center gap-2">
              <span class="text-xl font-extrabold text-zinc-950">{{ q.queue_number }}</span>
              <span class="text-xs px-2.5 py-0.5 rounded-full bg-amber-200/70 text-amber-900 font-bold">
                Tagihan: Rp 35.000
              </span>
            </div>
            <div class="text-sm font-bold text-zinc-900 mt-1">{{ q.patient_name }}</div>
            <div class="text-xs text-zinc-500">NIK: {{ q.patient_nik }} • {{ q.patient_phone }}</div>
            <div class="text-xs text-zinc-600 italic mt-1">"{{ q.complaint }}"</div>
          </div>

          <button
            @click="confirmCash(q.id)"
            :disabled="processingId === q.id"
            class="px-4 py-2.5 rounded-full bg-zinc-950 text-white hover:bg-zinc-800 text-xs font-bold transition-transform active:scale-95 shadow-sm whitespace-nowrap"
          >
            {{ processingId === q.id ? 'Memproses...' : '✓ Terima Tunai' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Section 2: Complete Queue Table (FIFO) -->
    <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] overflow-hidden">
      <h2 class="text-lg font-bold text-zinc-950 mb-4">Daftar Antrian Lengkap (FIFO)</h2>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-sm">
          <thead>
            <tr class="border-b border-zinc-100 text-[11px] font-bold uppercase tracking-wider text-zinc-400">
              <th class="pb-3 pl-2">Nomor</th>
              <th class="pb-3">Pasien</th>
              <th class="pb-3">Keluhan</th>
              <th class="pb-3">Status Bayar</th>
              <th class="pb-3">Status Antrian</th>
              <th class="pb-3 text-right pr-2">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-50">
            <tr
              v-for="q in queues"
              :key="q.id"
              class="hover:bg-stone-50/60 transition-colors"
            >
              <td class="py-4 pl-2 font-mono font-extrabold text-base text-zinc-950">
                {{ q.queue_number }}
              </td>
              <td class="py-4">
                <div class="font-bold text-zinc-900">{{ q.patient_name }}</div>
                <div class="text-xs text-zinc-400">NIK: {{ q.patient_nik }}</div>
              </td>
              <td class="py-4 text-zinc-600 max-w-xs truncate">
                {{ q.complaint || '-' }}
              </td>
              <td class="py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold',
                    q.payment_status === 'PAID'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200/60'
                      : 'bg-amber-50 text-amber-700 border border-amber-200/60'
                  ]"
                >
                  {{ q.payment_status === 'PAID' ? '✓ Lunas' : '⏳ Pending' }}
                  <span class="text-[10px] opacity-75">({{ q.payment_method }})</span>
                </span>
              </td>
              <td class="py-4">
                <span
                  :class="[
                    'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold',
                    q.queue_status === 'CALLING' && 'bg-blue-600 text-white animate-pulse',
                    q.queue_status === 'IN_EXAM' && 'bg-purple-100 text-purple-800',
                    q.queue_status === 'WAITING' && 'bg-stone-100 text-zinc-700',
                    q.queue_status === 'COMPLETED' && 'bg-emerald-100 text-emerald-800',
                    q.queue_status === 'CANCELLED' && 'bg-rose-100 text-rose-800'
                  ]"
                >
                  {{ q.queue_status }}
                </span>
              </td>
              <td class="py-4 text-right pr-2 space-x-1 whitespace-nowrap">
                <button
                  v-if="q.queue_status === 'WAITING'"
                  @click="updateStatus(q.id, 'CALLING')"
                  class="px-3 py-1.5 rounded-full bg-zinc-950 text-white hover:bg-zinc-800 text-xs font-bold active:scale-95 transition-all"
                >
                  Panggil
                </button>
                <button
                  v-if="q.queue_status === 'CALLING'"
                  @click="updateStatus(q.id, 'IN_EXAM')"
                  class="px-3 py-1.5 rounded-full bg-blue-600 text-white hover:bg-blue-700 text-xs font-bold active:scale-95 transition-all"
                >
                  Masuk Dokter
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
