<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type SatusehatStats } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { ShieldCheck, CheckCircle2, RefreshCw, Send, FileCode, Check } from 'lucide-vue-next'

const store = useClinicStore()
const stats = ref<SatusehatStats | null>(null)
const fhirPayload = ref<any>(null)
const loading = ref(false)
const syncing = ref(false)
const syncSuccess = ref(false)

async function loadData() {
  loading.value = true
  try {
    const [s, f] = await Promise.all([
      api.getSatusehatStats(),
      api.getSatusehatEncounter(1),
    ])
    stats.value = s
    fhirPayload.value = f
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

async function triggerSync() {
  syncing.value = true
  syncSuccess.value = false
  try {
    await api.syncSatusehat(1)
    store.playChime()
    syncSuccess.value = true
    await loadData()
  } catch (err: any) {
    alert(err.message || 'Gagal sinkronisasi SATUSEHAT')
  } finally {
    syncing.value = false
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
          Integrasi SATUSEHAT (FHIR R4 Ready)
        </h1>
        <p class="text-xs text-zinc-500 mt-1">
          Klinik terstandarisasi untuk regulasi interoperabilitas rekam medis Kemenkes RI.
        </p>
      </div>

      <button
        @click="triggerSync"
        :disabled="syncing"
        class="inline-flex items-center gap-2 px-6 py-3 rounded-full bg-amber-400 hover:bg-amber-500 text-zinc-950 text-xs font-extrabold shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <Send :class="['w-4 h-4', syncing && 'animate-spin']" />
        {{ syncing ? 'Mengirim ke Server DTO...' : 'Simulasi Sinkronisasi SATUSEHAT' }}
      </button>
    </div>

    <!-- Sync Success Banner -->
    <div v-if="syncSuccess" class="mb-6 p-4 rounded-3xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-sm font-semibold flex items-center gap-2 animate-in fade-in">
      <CheckCircle2 class="w-5 h-5 text-emerald-600" />
      Data kunjungan & resep berhasil disinkronkan ke SATUSEHAT Kemenkes (Status: 200 OK).
    </div>

    <!-- Stats Matrix -->
    <div v-if="stats" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
        <div class="text-xs font-bold text-zinc-400 uppercase tracking-wider">Status Faskes</div>
        <div class="text-lg font-extrabold text-emerald-700 mt-1 flex items-center gap-1.5">
          <CheckCircle2 class="w-4 h-4" /> TERVALIDASI
        </div>
        <div class="text-[11px] font-mono text-zinc-500 mt-0.5">Org: {{ stats.org_id }}</div>
      </div>

      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
        <div class="text-xs font-bold text-zinc-400 uppercase tracking-wider">Pasien IHS Terhubung</div>
        <div class="text-2xl font-extrabold text-zinc-950 mt-1">
          {{ stats.verified_patients }} / {{ stats.total_patients }}
        </div>
        <div class="text-[11px] text-zinc-500 mt-0.5">NIK terverifikasi Dukcapil</div>
      </div>

      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
        <div class="text-xs font-bold text-zinc-400 uppercase tracking-wider">Master Obat KFA</div>
        <div class="text-2xl font-extrabold text-zinc-950 mt-1">
          {{ stats.mapped_drugs }} / {{ stats.total_drugs }}
        </div>
        <div class="text-[11px] text-zinc-500 mt-0.5">Kamus Farmasi Kemenkes</div>
      </div>

      <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
        <div class="text-xs font-bold text-zinc-400 uppercase tracking-wider">Log Transaksi FHIR</div>
        <div class="text-2xl font-extrabold text-blue-700 mt-1">
          {{ stats.synced_logs }} Terkirim
        </div>
        <div class="text-[11px] text-zinc-500 mt-0.5">Encounter & Medication</div>
      </div>
    </div>

    <!-- FHIR Payload Inspector Card -->
    <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)]">
      <div class="flex items-center justify-between pb-4 border-b border-zinc-100 mb-4">
        <div class="flex items-center gap-2">
          <FileCode class="w-5 h-5 text-blue-600" />
          <h2 class="text-base font-bold text-zinc-950">JSON Payload FHIR R4 Bundle (HL7 Standard)</h2>
        </div>
        <span class="text-xs font-mono font-bold px-2.5 py-0.5 rounded-full bg-stone-100 text-zinc-700">
          Encounter & Condition
        </span>
      </div>

      <pre class="p-4 rounded-2xl bg-zinc-950 text-zinc-200 text-xs font-mono overflow-x-auto max-h-[400px] leading-relaxed">{{ JSON.stringify(fhirPayload, null, 2) }}</pre>
    </div>
  </div>
</template>
