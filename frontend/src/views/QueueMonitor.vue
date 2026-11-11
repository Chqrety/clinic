<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import {
  Volume2,
  VolumeX,
  Clock,
  Tv,
  ArrowLeft,
  Users,
  Stethoscope,
  Building2,
  Sparkles,
  Baby,
  ShieldCheck
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
}>()

interface QueueItem {
  id: number
  patient_name: string
  clinic_unit: string
  category: string
  doctor_name: string
  queue_number: string
  queue_status: string
}

interface MonitorData {
  currently_calling: QueueItem | null
  upcoming_queue: QueueItem[]
  clinic_unit: string
}

const selectedUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const monitorData = ref<MonitorData>({
  currently_calling: {
    id: 1,
    patient_name: 'Budi P****',
    clinic_unit: 'KIDINEST',
    category: 'A_VAKSIN',
    doctor_name: 'dr. Hendra Pratama, Sp.A',
    queue_number: 'A-001',
    queue_status: 'CALLING',
  },
  upcoming_queue: [
    {
      id: 2,
      patient_name: 'Clarissa W*****',
      clinic_unit: 'SSC',
      category: 'SSC_TREATMENT',
      doctor_name: 'dr. Nadia Kusuma, Sp.KK',
      queue_number: 'S-001',
      queue_status: 'WAITING',
    },
    {
      id: 3,
      patient_name: 'Ahmad F****',
      clinic_unit: 'KIDINEST',
      category: 'B_SAKIT',
      doctor_name: 'dr. Hendra Pratama, Sp.A',
      queue_number: 'B-001',
      queue_status: 'WAITING',
    },
    {
      id: 4,
      patient_name: 'Siti R*******',
      clinic_unit: 'KIDINEST',
      category: 'A_VAKSIN',
      doctor_name: 'dr. Hendra Pratama, Sp.A',
      queue_number: 'A-002',
      queue_status: 'WAITING',
    },
  ],
  clinic_unit: 'ALL',
})

const currentTime = ref('')
const soundEnabled = ref(true)
let clockInterval: any = null
let eventSource: EventSource | null = null
let pollInterval: any = null

function updateClock() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function playAudioChime() {
  if (!soundEnabled.value) return
  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.type = 'sine'
    osc.frequency.setValueAtTime(587.33, ctx.currentTime) // D5
    osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.15) // A5
    gain.gain.setValueAtTime(0.2, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.6)
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.start()
    osc.stop(ctx.currentTime + 0.6)
  } catch {
    // Autoplay policy fallback
  }
}

async function fetchMonitorData() {
  try {
    const unitParam = selectedUnit.value === 'ALL' ? '' : `?unit=${selectedUnit.value}`
    const res = await fetch(`/api/public/queue/monitor${unitParam}`)
    const data = await res.json()
    if (res.ok && data.success && data.data) {
      if (
        data.data.currently_calling &&
        monitorData.value.currently_calling?.queue_number !== data.data.currently_calling.queue_number
      ) {
        playAudioChime()
      }
      monitorData.value = data.data
    }
  } catch {
    // Keep local state
  }
}

function initSSE() {
  if (eventSource) {
    eventSource.close()
  }

  const unitParam = selectedUnit.value === 'ALL' ? '' : `?unit=${selectedUnit.value}`
  try {
    eventSource = new EventSource(`/api/public/queue/stream${unitParam}`)
    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data) {
          if (
            data.currently_calling &&
            monitorData.value.currently_calling?.queue_number !== data.currently_calling.queue_number
          ) {
            playAudioChime()
          }
          monitorData.value = data
        }
      } catch {
        // Ignored
      }
    }
    eventSource.onerror = () => {
      eventSource?.close()
      eventSource = null
    }
  } catch {
    // Fallback to polling
  }
}

watch(selectedUnit, () => {
  fetchMonitorData()
  initSSE()
})

onMounted(() => {
  updateClock()
  clockInterval = setInterval(updateClock, 1000)
  fetchMonitorData()
  initSSE()
  pollInterval = setInterval(fetchMonitorData, 4000)
})

onUnmounted(() => {
  if (clockInterval) clearInterval(clockInterval)
  if (pollInterval) clearInterval(pollInterval)
  if (eventSource) eventSource.close()
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col justify-between p-6 sm:p-10 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300">
    <!-- Top Header Display -->
    <header class="w-full flex items-center justify-between pb-6 border-b border-zinc-200/80">
      <!-- Left: Clinic Brand & Back -->
      <div class="flex items-center gap-4">
        <button
          @click="emit('navigate', 'landing')"
          class="p-2 rounded-2xl bg-white border border-zinc-200/80 hover:bg-stone-100 text-zinc-500 hover:text-zinc-950 transition-colors shadow-2xs"
          title="Kembali ke Beranda"
        >
          <ArrowLeft class="w-4 h-4" />
        </button>

        <div>
          <div class="flex items-center gap-2">
            <span class="w-2.5 h-2.5 rounded-full bg-emerald-500 animate-pulse"></span>
            <span class="text-[11px] font-bold uppercase tracking-wider text-emerald-800">Layar Antrian Live Ruang Tunggu</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-black text-zinc-950 tracking-tight">
            Kidinest & Semarang Skin Center
          </h1>
        </div>
      </div>

      <!-- Center: TV Filter Pill Bar -->
      <div class="hidden md:flex items-center p-1 bg-stone-200/60 rounded-full text-xs font-bold text-zinc-600">
        <button
          @click="selectedUnit = 'ALL'"
          :class="['px-4 py-1.5 rounded-full transition-all', selectedUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          Semua Layanan
        </button>
        <button
          @click="selectedUnit = 'KIDINEST'"
          :class="['px-4 py-1.5 rounded-full transition-all', selectedUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          Kidinest (Poli Anak)
        </button>
        <button
          @click="selectedUnit = 'SSC'"
          :class="['px-4 py-1.5 rounded-full transition-all', selectedUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          SSC (Estetika Kulit)
        </button>
      </div>

      <!-- Right: Clock & Audio -->
      <div class="flex items-center gap-3">
        <button
          @click="soundEnabled = !soundEnabled"
          :title="soundEnabled ? 'Lonceng Pemanggilan Aktif' : 'Lonceng Senyap'"
          class="p-2.5 rounded-2xl bg-white border border-zinc-200/80 text-zinc-600 hover:text-zinc-950 transition-colors shadow-2xs"
        >
          <Volume2 v-if="soundEnabled" class="w-4 h-4 text-emerald-600" />
          <VolumeX v-else class="w-4 h-4 text-zinc-400" />
        </button>

        <div class="px-4 py-2 rounded-2xl bg-white border border-zinc-200/80 shadow-2xs text-right">
          <div class="font-mono text-base font-black text-zinc-950">{{ currentTime }}</div>
          <div class="text-[9px] font-bold uppercase tracking-wider text-zinc-400">WIB • Online</div>
        </div>
      </div>
    </header>

    <!-- Main Live Screen (12 cols) -->
    <main class="flex-1 py-8 grid grid-cols-1 md:grid-cols-12 gap-8 items-stretch">

      <!-- HERO LEFT: CURRENTLY CALLING (7 cols) -->
      <div class="md:col-span-7 bg-white p-8 sm:p-12 rounded-3xl border border-zinc-200/80 shadow-xs flex flex-col justify-between text-center relative overflow-hidden">
        <!-- Top Status Indicator -->
        <div class="flex items-center justify-center gap-2">
          <span class="inline-flex items-center gap-2 px-4 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-xs">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
            Sedang Dipanggil
          </span>
        </div>

        <!-- Giant Queue Number -->
        <div v-if="monitorData.currently_calling" class="py-6 space-y-4">
          <div class="text-7xl sm:text-9xl font-black font-mono tracking-tight text-zinc-950 leading-none">
            {{ monitorData.currently_calling.queue_number }}
          </div>

          <!-- Masked Name -->
          <div class="space-y-1">
            <div class="text-xl sm:text-2xl font-extrabold text-zinc-800">
              {{ monitorData.currently_calling.patient_name }}
            </div>
            <div class="inline-flex items-center gap-1.5 text-xs font-semibold text-zinc-500">
              <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
              <span>Identitas disamarkan untuk privasi</span>
            </div>
          </div>
        </div>

        <div v-else class="py-16 text-zinc-400 space-y-2">
          <Users class="w-12 h-12 mx-auto text-zinc-300" />
          <div class="text-base font-bold text-zinc-600">Belum Ada Panggilan Aktif</div>
          <div class="text-xs">Pasien berikutnya harap bersiap di ruang tunggu.</div>
        </div>

        <!-- Doctor & Room Banner -->
        <div v-if="monitorData.currently_calling" class="p-4 sm:p-5 rounded-2xl bg-stone-50 border border-zinc-200/80 flex items-center justify-between text-left">
          <div class="flex items-center gap-3">
            <div class="p-2.5 rounded-xl bg-white border border-zinc-200 text-zinc-900 shrink-0">
              <Stethoscope class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs font-extrabold text-zinc-950">
                {{ monitorData.currently_calling.doctor_name }}
              </div>
              <div class="text-[11px] text-zinc-500 font-medium">
                {{ monitorData.currently_calling.clinic_unit === 'KIDINEST' ? 'Ruang Pemeriksaan Poli Anak' : 'Ruang Tindakan Estetika Kulit' }}
              </div>
            </div>
          </div>

          <span class="px-3 py-1 rounded-full text-xs font-black bg-zinc-950 text-white font-mono">
            RUANG 1
          </span>
        </div>
      </div>

      <!-- UPCOMING QUEUE LIST (5 cols) -->
      <div class="md:col-span-5 bg-white p-6 sm:p-8 rounded-3xl border border-zinc-200/80 shadow-xs flex flex-col justify-between space-y-4">
        <div>
          <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
            <div>
              <h2 class="text-sm font-extrabold text-zinc-950">Antrian Berikutnya</h2>
              <p class="text-[10px] text-zinc-400 font-medium">Harap bersiap di dekat pintu ruang periksa</p>
            </div>
            <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-stone-100 text-zinc-700">
              {{ monitorData.upcoming_queue.length }} Menunggu
            </span>
          </div>

          <!-- List Items -->
          <div v-if="monitorData.upcoming_queue.length === 0" class="py-12 text-center text-zinc-400 space-y-1">
            <Clock class="w-8 h-8 mx-auto text-zinc-300" />
            <p class="text-xs font-bold text-zinc-600">Tidak ada antrian menunggu</p>
          </div>

          <div v-else class="divide-y divide-zinc-100 mt-2">
            <div
              v-for="(item, idx) in monitorData.upcoming_queue"
              :key="item.id"
              class="py-3 flex items-center justify-between gap-3"
            >
              <div class="flex items-center gap-3">
                <span class="w-6 h-6 rounded-full bg-stone-100 text-zinc-500 flex items-center justify-center text-xs font-black">
                  {{ idx + 1 }}
                </span>
                <div>
                  <div class="font-mono text-base font-black text-zinc-950 leading-tight">
                    {{ item.queue_number }}
                  </div>
                  <div class="text-[11px] font-bold text-zinc-700">
                    {{ item.patient_name }}
                  </div>
                </div>
              </div>

              <div class="text-right">
                <span :class="[
                  'px-2 py-0.5 rounded-full text-[9px] font-bold border',
                  item.clinic_unit === 'KIDINEST'
                    ? 'bg-sky-50 text-sky-800 border-sky-200'
                    : 'bg-rose-50 text-rose-800 border-rose-200'
                ]">
                  {{ item.clinic_unit === 'KIDINEST' ? 'Kidinest' : 'SSC' }}
                </span>
                <div class="text-[9px] text-zinc-400 font-medium mt-0.5">Menunggu</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Bottom Running Info -->
        <div class="p-3.5 rounded-2xl bg-[#F9F9F7] border border-zinc-200/80 text-[10px] text-zinc-500 font-medium text-center">
          💡 Pasien yang telah dipanggil 3 kali tanpa konfirmasi akan dialihkan ke antrian toleransi.
        </div>
      </div>

    </main>

    <!-- TV Footer Ticker -->
    <footer class="w-full pt-4 border-t border-zinc-200/60 flex items-center justify-between text-xs text-zinc-400">
      <div>Sistem Antrian Berkeadilan (FIFO) • Paperless Smart Clinic Management</div>
      <div class="font-mono text-[11px] text-zinc-500 font-bold">Kidinest &amp; Semarang Skin Center</div>
    </footer>
  </div>
</template>
