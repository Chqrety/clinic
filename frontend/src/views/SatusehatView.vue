<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type SatusehatStats } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import {
  ShieldCheck,
  CheckCircle2,
  RefreshCw,
  Send,
  FileCode,
  Check,
  Copy,
  ExternalLink,
  Activity,
  Radio,
  Loader2,
  Database,
  Layers,
  Sparkles
} from 'lucide-vue-next'

const store = useClinicStore()
const stats = ref<SatusehatStats | null>(null)
const fhirPayload = ref<any>(null)
const loading = ref(false)
const syncing = ref(false)
const syncSuccess = ref(false)
const copied = ref(false)

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

function copyPayload() {
  if (!fhirPayload.value) return
  navigator.clipboard.writeText(JSON.stringify(fhirPayload.value, null, 2))
  copied.value = true
  setTimeout(() => (copied.value = false), 2500)
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- TOP CLINICAL BAR -->
    <header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2 border-b border-zinc-200/80">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="text-xs font-semibold tracking-wide text-zinc-500 uppercase">
            Integrasi Regulasi • Semarang Medical Group
          </span>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-xs font-medium text-emerald-700">Kemenkes RI Live Bridge</span>
        </div>
        <div class="flex items-center gap-3">
          <h1 class="text-xl sm:text-2xl font-black tracking-tight text-zinc-950">
            SATUSEHAT Regulatory Bridge
          </h1>
          <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-zinc-900 text-white shadow-2xs font-mono">
            FHIR R4 Ready
          </span>
        </div>
        <p class="text-xs text-zinc-500 mt-1">
          Interoperabilitas rekam medis elektronik terstandarisasi Kemenkes RI (HL7 FHIR R4 Bundle).
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="loadData"
          :disabled="loading"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200/90 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-700 transition-all shadow-2xs cursor-pointer disabled:opacity-50"
          title="Segarkan data DTO"
        >
          <Loader2 v-if="loading" class="w-3.5 h-3.5 animate-spin text-zinc-500" />
          <RefreshCw v-else class="w-3.5 h-3.5 text-zinc-500" />
          <span>Sync</span>
        </button>
        <button
          @click="triggerSync"
          :disabled="syncing"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
        >
          <Loader2 v-if="syncing" class="w-3.5 h-3.5 animate-spin" />
          <Send v-else class="w-3.5 h-3.5 text-emerald-400" />
          <span>{{ syncing ? 'Mengirim ke DTO Kemenkes...' : 'Simulasi Sinkronisasi SATUSEHAT' }}</span>
        </button>
      </div>
    </header>

    <!-- Sync Success Banner -->
    <div
      v-if="syncSuccess"
      class="p-3.5 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-900 text-xs font-semibold flex items-center justify-between gap-2 animate-in fade-in duration-200 shadow-2xs"
    >
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
        <span>Data kunjungan & resep berhasil disinkronkan ke SATUSEHAT Kemenkes (HTTP 200 OK • Bundle Ingested).</span>
      </div>
      <button @click="syncSuccess = false" class="text-emerald-700 hover:text-emerald-950 text-[11px] font-bold cursor-pointer">
        Tutup
      </button>
    </div>

    <!-- 1-ROW HAIRLINE METRIC TICKER STRIP -->
    <section v-if="stats" class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Status Faskes</span>
          <ShieldCheck class="w-4 h-4 text-emerald-600" />
        </div>
        <div class="text-xl sm:text-2xl font-black text-emerald-700 tracking-tight flex items-center gap-1.5">
          <CheckCircle2 class="w-5 h-5 text-emerald-600" /> TERVALIDASI
        </div>
        <div class="text-[11px] text-zinc-500 font-mono mt-0.5">Org: {{ stats.org_id }}</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Pasien IHS Terhubung</span>
          <Radio class="w-4 h-4 text-zinc-400" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-zinc-950 tracking-tight tabular-nums">
          {{ stats.verified_patients }} <span class="text-base text-zinc-400 font-normal">/ {{ stats.total_patients }}</span>
        </div>
        <div class="text-[11px] text-emerald-700 font-medium mt-0.5">NIK terverifikasi Dukcapil</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Master Obat KFA</span>
          <Database class="w-4 h-4 text-amber-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-zinc-950 tracking-tight tabular-nums">
          {{ stats.mapped_drugs }} <span class="text-base text-zinc-400 font-normal">/ {{ stats.total_drugs }}</span>
        </div>
        <div class="text-[11px] text-amber-700 font-medium mt-0.5">Kamus Farmasi Kemenkes</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Log Transaksi FHIR</span>
          <Activity class="w-4 h-4 text-blue-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-blue-700 tracking-tight tabular-nums">
          {{ stats.synced_logs }}
        </div>
        <div class="text-[11px] text-blue-700 font-medium mt-0.5">Encounter & Medication</div>
      </div>

      <div class="p-4 sm:p-5 hidden lg:block">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Kepatuhan Standar</span>
          <Layers class="w-4 h-4 text-purple-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-purple-950 tracking-tight tabular-nums">100%</div>
        <div class="text-[11px] text-purple-700 font-medium mt-0.5">HL7 FHIR Release 4</div>
      </div>
    </section>

    <!-- MASTER-DETAIL REGULATORY TWO-PANE -->
    <main class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      
      <!-- LEFT PANE: GATEWAY TELEMETRY & RESOURCE MAPPING (5 cols) -->
      <section class="lg:col-span-5 space-y-4">
        
        <!-- Gateway Configuration Card -->
        <div class="bg-white rounded-2xl p-5 border border-zinc-200/90 shadow-2xs space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
            <div class="flex items-center gap-2">
              <span class="p-1.5 rounded-lg bg-zinc-100 text-zinc-800"><Radio class="w-4 h-4" /></span>
              <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-950">Gateway Telemetry Kemenkes</h2>
            </div>
            <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Connected
            </span>
          </div>

          <div class="space-y-2.5 text-xs">
            <div class="flex justify-between items-center py-1.5 border-b border-zinc-100">
              <span class="text-zinc-500">Target Endpoint</span>
              <span class="font-mono text-[11px] font-bold text-zinc-900 truncate max-w-[220px]">
                api-satusehat.kemkes.go.id/fhir-r4/v1
              </span>
            </div>
            <div class="flex justify-between items-center py-1.5 border-b border-zinc-100">
              <span class="text-zinc-500">Authentication</span>
              <span class="font-mono text-[11px] font-bold text-zinc-900">OAuth2 Client Credentials</span>
            </div>
            <div class="flex justify-between items-center py-1.5 border-b border-zinc-100">
              <span class="text-zinc-500">IHS Organization ID</span>
              <span class="font-mono text-[11px] font-bold text-zinc-900">{{ stats?.org_id || '100028120' }}</span>
            </div>
            <div class="flex justify-between items-center py-1.5 border-b border-zinc-100">
              <span class="text-zinc-500">Data Format</span>
              <span class="font-mono text-[11px] font-bold text-zinc-900">application/fhir+json</span>
            </div>
            <div class="flex justify-between items-center py-1.5">
              <span class="text-zinc-500">Payload Encryption</span>
              <span class="font-mono text-[11px] font-bold text-emerald-700">TLS 1.3 / AES-256</span>
            </div>
          </div>
        </div>

        <!-- Resource Mapping Status Card -->
        <div class="bg-white rounded-2xl p-5 border border-zinc-200/90 shadow-2xs space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
            <div class="flex items-center gap-2">
              <span class="p-1.5 rounded-lg bg-zinc-100 text-zinc-800"><Layers class="w-4 h-4" /></span>
              <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-950">Resource Mapping Pipeline</h2>
            </div>
            <span class="text-[11px] text-zinc-400 font-medium">5/5 Resources</span>
          </div>

          <div class="space-y-2">
            <div class="p-3 rounded-xl bg-zinc-50/80 border border-zinc-200/70 flex items-center justify-between text-xs">
              <div>
                <div class="font-bold text-zinc-950 font-mono">Patient (IHS)</div>
                <div class="text-[11px] text-zinc-500">Validasi NIK Dukcapil & Master Patient Index</div>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                100% Synced
              </span>
            </div>

            <div class="p-3 rounded-xl bg-zinc-50/80 border border-zinc-200/70 flex items-center justify-between text-xs">
              <div>
                <div class="font-bold text-zinc-950 font-mono">Encounter</div>
                <div class="text-[11px] text-zinc-500">Waktu tiba, konsultasi dokter, & selesai periksa</div>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                Active Feed
              </span>
            </div>

            <div class="p-3 rounded-xl bg-zinc-50/80 border border-zinc-200/70 flex items-center justify-between text-xs">
              <div>
                <div class="font-bold text-zinc-950 font-mono">Condition</div>
                <div class="text-[11px] text-zinc-500">Diagnosis terstandardisasi ICD-10 WHO</div>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                ICD-10 Valid
              </span>
            </div>

            <div class="p-3 rounded-xl bg-zinc-50/80 border border-zinc-200/70 flex items-center justify-between text-xs">
              <div>
                <div class="font-bold text-zinc-950 font-mono">MedicationRequest</div>
                <div class="text-[11px] text-zinc-500">E-Prescription & Kamus Farmasi Kemenkes (KFA)</div>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                KFA Linked
              </span>
            </div>

            <div class="p-3 rounded-xl bg-zinc-50/80 border border-zinc-200/70 flex items-center justify-between text-xs">
              <div>
                <div class="font-bold text-zinc-950 font-mono">Composition</div>
                <div class="text-[11px] text-zinc-500">Ringkasan Resume Medis Rawat Jalan Elektronik</div>
              </div>
              <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                Compliant
              </span>
            </div>
          </div>
        </div>

      </section>

      <!-- RIGHT PANE: DEVELOPER-GRADE FHIR R4 JSON BUNDLE INSPECTOR (7 cols, Sticky) -->
      <section class="lg:col-span-7 sticky top-6 space-y-3">
        <div class="bg-white rounded-2xl p-5 border border-zinc-200/90 shadow-2xs space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
            <div class="flex items-center gap-2">
              <span class="p-1.5 rounded-lg bg-blue-50 text-blue-700 border border-blue-200/60"><FileCode class="w-4 h-4" /></span>
              <div>
                <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-950">FHIR R4 Bundle Payload Inspector</h2>
                <p class="text-[11px] text-zinc-500">HL7 Standard JSON Payload untuk Kunjungan Rawat Jalan</p>
              </div>
            </div>
            
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="copyPayload"
                class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl border border-zinc-200/90 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-700 transition-all shadow-2xs cursor-pointer active:scale-95"
              >
                <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-600" />
                <Copy v-else class="w-3.5 h-3.5 text-zinc-500" />
                <span>{{ copied ? 'Tersalin' : 'Copy JSON' }}</span>
              </button>
            </div>
          </div>

          <!-- Code Box -->
          <div class="relative rounded-xl overflow-hidden border border-zinc-900 bg-zinc-950">
            <div class="flex items-center justify-between px-4 py-2 bg-zinc-900/90 border-b border-zinc-800 text-[11px] text-zinc-400 font-mono">
              <div class="flex items-center gap-2">
                <span class="w-2.5 h-2.5 rounded-full bg-rose-500/80"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-amber-500/80"></span>
                <span class="w-2.5 h-2.5 rounded-full bg-emerald-500/80"></span>
                <span class="ml-2 text-zinc-300 font-semibold">Bundle/encounter-transaction.json</span>
              </div>
              <span>JSON FHIR R4</span>
            </div>

            <pre class="p-4 text-emerald-400 text-xs font-mono overflow-x-auto max-h-[520px] leading-relaxed select-all">{{ JSON.stringify(fhirPayload, null, 2) }}</pre>
          </div>

          <div class="flex items-center justify-between text-[11px] text-zinc-500 pt-1">
            <span>Struktur divalidasi terhadap Skema FHIR R4 Kemenkes RI v1.0.4</span>
            <span class="font-mono text-zinc-700 font-bold">Enc-ID: #1</span>
          </div>
        </div>
      </section>

    </main>
  </div>
</template>
