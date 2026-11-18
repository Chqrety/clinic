<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import {
  Volume2,
  VolumeX,
  Clock,
  ArrowLeft,
  Users,
  Stethoscope,
  Building2,
  Sparkles,
  Baby,
  ShieldCheck,
  CheckCircle2,
  Radio,
  DoorOpen,
} from 'lucide-vue-next'
import { announcePatientCall } from '../utils/audioAnnouncement'

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

interface DoctorRoomLive {
  doctor_name: string
  clinic_unit: string
  room_name: string
  status: string
  current_patient?: QueueItem
  waiting_count: number
  next_patient?: QueueItem
}

const selectedUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const soundEnabled = ref(true)
const userInteracted = ref(false)

const monitorData = ref<MonitorData>({
  currently_calling: null,
  upcoming_queue: [],
  clinic_unit: 'ALL',
})

const doctorRooms = ref<DoctorRoomLive[]>([
  {
    doctor_name: 'dr. Hendra Pratama, Sp.A',
    clinic_unit: 'KIDINEST',
    room_name: 'Poli Anak 1',
    status: 'READY',
    waiting_count: 0,
  },
  {
    doctor_name: 'dr. Patricia, Sp.A',
    clinic_unit: 'KIDINEST',
    room_name: 'Poli Anak 2',
    status: 'READY',
    waiting_count: 0,
  },
  {
    doctor_name: 'dr. Dian Paramita, Sp.D.V.E',
    clinic_unit: 'SSC',
    room_name: 'Poli Estetika 1',
    status: 'READY',
    waiting_count: 0,
  },
  {
    doctor_name: 'dr. Kevin Wijaya, Sp.D.V.E',
    clinic_unit: 'SSC',
    room_name: 'Poli Estetika 2',
    status: 'READY',
    waiting_count: 0,
  },
])

const currentTime = ref('')
const currentDate = ref('')
let clockInterval: any = null
let eventSource: EventSource | null = null
let pollInterval: any = null
const lastCalledQueueNumber = ref('')

function updateClock() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  currentDate.value = now.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

function resolveRoomName(item: QueueItem): string {
  if (item.clinic_unit === 'SSC') {
    return 'Ruang Tindakan Estetika 1'
  }
  if (item.doctor_name.includes('Patricia')) {
    return 'Ruang Poli Anak 2'
  }
  return 'Ruang Poli Anak 1'
}

async function triggerVoiceCall(item: QueueItem) {
  const roomName = resolveRoomName(item)
  await announcePatientCall(item.queue_number, roomName, item.doctor_name, soundEnabled.value)
}

async function fetchMonitorData() {
  try {
    const unitParam = selectedUnit.value === 'ALL' ? '' : `?unit=${selectedUnit.value}`
    const res = await fetch(`/api/public/queue/monitor${unitParam}`)
    const data = await res.json()
    if (res.ok && data.success && data.data) {
      const calling = data.data.currently_calling
      if (calling) {
        if (calling.queue_number !== lastCalledQueueNumber.value || calling.queue_status === 'CALLING') {
          if (calling.queue_number !== lastCalledQueueNumber.value) {
            triggerVoiceCall(calling)
            lastCalledQueueNumber.value = calling.queue_number
          }
        }
      } else {
        lastCalledQueueNumber.value = ''
      }
      monitorData.value = data.data
    }
  } catch {
    // Keep local fallback
  }

  // Fetch live doctor rooms
  try {
    const rRes = await fetch('/api/public/queue/rooms/status')
    const rData = await rRes.json()
    if (rRes.ok && rData.success && rData.data) {
      doctorRooms.value = rData.data
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
          if (data.currently_calling) {
            const calling = data.currently_calling
            if (calling.queue_number !== lastCalledQueueNumber.value) {
              triggerVoiceCall(calling)
              lastCalledQueueNumber.value = calling.queue_number
            }
          } else {
            lastCalledQueueNumber.value = ''
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

function enableAudioInteraction() {
  userInteracted.value = true
  soundEnabled.value = true
  if (monitorData.value.currently_calling) {
    triggerVoiceCall(monitorData.value.currently_calling)
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
  pollInterval = setInterval(fetchMonitorData, 3000)

  // Listen for browser click to unlock audio
  const unlockAudio = () => {
    userInteracted.value = true
    window.removeEventListener('click', unlockAudio)
  }
  window.addEventListener('click', unlockAudio)
})

onUnmounted(() => {
  if (clockInterval) clearInterval(clockInterval)
  if (pollInterval) clearInterval(pollInterval)
  if (eventSource) eventSource.close()
})
</script>

<template>
  <div class="min-h-screen bg-[#FAF8F5] text-zinc-950 flex flex-col justify-between p-4 sm:p-8 lg:p-10 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300">
    
    <!-- TOP HEADER BAR (Minimal Nordic Editorial) -->
    <header class="w-full flex flex-col md:flex-row md:items-center justify-between gap-4 pb-5 border-b border-zinc-200/80">
      
      <!-- Brand & Status Tag -->
      <div class="flex items-center gap-3.5">
        <button
          @click="emit('navigate', 'landing')"
          class="p-2.5 rounded-2xl bg-white border border-zinc-200/90 hover:bg-zinc-100 text-zinc-600 hover:text-zinc-950 transition-colors shadow-2xs cursor-pointer"
          title="Kembali ke Beranda"
        >
          <ArrowLeft class="w-4 h-4" />
        </button>

        <div>
          <div class="flex items-center gap-2">
            <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-emerald-50 text-emerald-900 border border-emerald-200/70">
              <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
              Live Queue Display
            </span>
            <span class="text-xs text-zinc-400">•</span>
            <span class="text-xs font-semibold text-zinc-500 font-mono">Semarang Medical Group</span>
          </div>
          <h1 class="text-lg sm:text-2xl font-black text-zinc-950 tracking-tight mt-0.5">
            Kidinest Pediatric Care & Semarang Skin Center
          </h1>
        </div>
      </div>

      <!-- Unit Filter Pills -->
      <div class="flex items-center p-1 bg-stone-200/60 rounded-full text-xs font-bold text-zinc-600 self-start md:self-auto shadow-2xs">
        <button
          @click="selectedUnit = 'ALL'"
          :class="['px-4 py-1.5 rounded-full transition-all cursor-pointer', selectedUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          Semua Layanan
        </button>
        <button
          @click="selectedUnit = 'KIDINEST'"
          :class="['px-4 py-1.5 rounded-full transition-all cursor-pointer', selectedUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          Kidinest (Poli Anak)
        </button>
        <button
          @click="selectedUnit = 'SSC'"
          :class="['px-4 py-1.5 rounded-full transition-all cursor-pointer', selectedUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-xs' : 'hover:text-zinc-950']"
        >
          SSC (Estetika Kulit)
        </button>
      </div>

      <!-- Audio Toggle & Clock -->
      <div class="flex items-center gap-3 self-end md:self-auto">
        <button
          @click="soundEnabled = !soundEnabled; userInteracted = true"
          :title="soundEnabled ? 'Suara Pemanggilan Aktif' : 'Suara Senyap'"
          class="p-2.5 rounded-2xl bg-white border border-zinc-200/90 text-zinc-600 hover:text-zinc-950 transition-colors shadow-2xs cursor-pointer flex items-center gap-2"
        >
          <Volume2 v-if="soundEnabled" class="w-4 h-4 text-emerald-600" />
          <VolumeX v-else class="w-4 h-4 text-zinc-400" />
          <span class="text-[11px] font-bold hidden sm:inline">{{ soundEnabled ? 'Suara Aktif' : 'Mute' }}</span>
        </button>

        <div class="px-4 py-2 rounded-2xl bg-white border border-zinc-200/90 shadow-2xs text-right">
          <div class="font-mono text-base font-black text-zinc-950 leading-tight">{{ currentTime }}</div>
          <div class="text-[9px] font-bold uppercase tracking-wider text-zinc-400">{{ currentDate }}</div>
        </div>
      </div>
    </header>

    <!-- UNMUTED NOTIFICATION BANNER (if browser audio blocked) -->
    <div
      v-if="!userInteracted && soundEnabled"
      @click="enableAudioInteraction"
      class="mt-3 p-3 rounded-2xl bg-amber-50 border border-amber-200/80 text-xs text-amber-900 font-semibold flex items-center justify-between cursor-pointer hover:bg-amber-100/70 transition-all shadow-2xs"
    >
      <div class="flex items-center gap-2">
        <Volume2 class="w-4 h-4 text-amber-700 animate-bounce" />
        <span>Klik di sini sekali untuk mengaktifkan pemanggilan suara otomatis (Audio Synthesizer).</span>
      </div>
      <span class="px-2.5 py-1 rounded-xl bg-amber-200/60 text-[10px] font-bold uppercase">Aktifkan Suara</span>
    </div>

    <!-- MAIN MONITOR DISPLAY GRID (12 Cols) -->
    <main class="flex-1 py-6 grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 items-stretch">

      <!-- HERO LEFT: GIANT ACTIVE CALL CARD (7 Cols) -->
      <div class="lg:col-span-7 bg-white p-6 sm:p-10 rounded-3xl border border-zinc-200/90 shadow-xs flex flex-col justify-between text-center relative overflow-hidden min-h-[500px]">
        
        <!-- Status Pill Indicator -->
        <div class="flex items-center justify-between w-full">
          <span class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-xs">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
            Sedang Dipanggil
          </span>

          <span
            v-if="monitorData.currently_calling"
            :class="[
              'px-3 py-1 rounded-full text-xs font-extrabold border',
              monitorData.currently_calling.clinic_unit === 'KIDINEST'
                ? 'bg-amber-50 text-amber-900 border-amber-200'
                : 'bg-rose-50 text-rose-900 border-rose-200'
            ]"
          >
            {{ monitorData.currently_calling.clinic_unit === 'KIDINEST' ? 'Lantai 1 • Poli Anak' : 'Lantai 2 • Estetika Medis' }}
          </span>
        </div>

        <!-- Giant Queue Number -->
        <div v-if="monitorData.currently_calling" class="py-6 sm:py-8 space-y-4">
          <div class="text-7xl sm:text-9xl md:text-[10rem] font-black font-mono tracking-tight text-zinc-950 leading-none select-all transition-transform">
            {{ monitorData.currently_calling.queue_number }}
          </div>

          <!-- Patient Name (Masked for Privacy) -->
          <div class="space-y-1.5">
            <div class="text-2xl sm:text-3xl font-extrabold text-zinc-900 tracking-tight">
              {{ monitorData.currently_calling.patient_name }}
            </div>
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-stone-100/80 text-[11px] font-semibold text-zinc-600">
              <ShieldCheck class="w-3.5 h-3.5 text-emerald-600" />
              <span>Privasi Terjaga • Nama belakang disamarkan</span>
            </div>
          </div>
        </div>

        <!-- Empty / Standby State -->
        <div v-else class="py-16 text-zinc-400 space-y-3">
          <Users class="w-14 h-14 mx-auto text-zinc-300" />
          <h2 class="text-lg font-bold text-zinc-700">Ruang Pemeriksaan Siap Melayani</h2>
          <p class="text-xs max-w-sm mx-auto text-zinc-500">
            Pasien dengan nomor antrian berikutnya harap bersiap di ruang tunggu depan poliklinik.
          </p>
        </div>

        <!-- Destination Room & Doctor Footer Card -->
        <div
          v-if="monitorData.currently_calling"
          class="p-4 sm:p-5 rounded-2xl bg-stone-50/80 border border-zinc-200/90 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-left"
        >
          <div class="flex items-center gap-3.5">
            <div class="p-3 rounded-2xl bg-white border border-zinc-200 text-zinc-900 shadow-2xs shrink-0">
              <Stethoscope class="w-6 h-6 text-[#0C3B36]" />
            </div>
            <div>
              <div class="text-xs text-zinc-400 font-semibold uppercase tracking-wider">Dokter Pemeriksa</div>
              <div class="text-sm sm:text-base font-extrabold text-zinc-950">
                {{ monitorData.currently_calling.doctor_name }}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-2 self-start sm:self-auto">
            <div class="px-4 py-2 rounded-2xl bg-zinc-950 text-white font-mono font-black text-xs sm:text-sm tracking-wide shadow-xs flex items-center gap-2">
              <DoorOpen class="w-4 h-4 text-emerald-400" />
              <span>{{ resolveRoomName(monitorData.currently_calling).toUpperCase() }}</span>
            </div>
          </div>
        </div>

        <!-- Standby Hint -->
        <div v-else class="p-3 rounded-2xl bg-stone-50 border border-zinc-200/80 text-xs text-zinc-500">
          💡 TV Monitor diperbarui otomatis secara real-time via Server-Sent Events.
        </div>
      </div>

      <!-- RIGHT SIDE: MULTI-ROOM STATUS & NEXT QUEUES (5 Cols) -->
      <div class="lg:col-span-5 flex flex-col gap-6">

        <!-- 1. LIVE DOCTOR ROOMS STATUS -->
        <div class="bg-white p-5 sm:p-6 rounded-3xl border border-zinc-200/90 shadow-xs space-y-3">
          <div class="flex items-center justify-between pb-2.5 border-b border-zinc-100">
            <div class="flex items-center gap-2">
              <Building2 class="w-4 h-4 text-zinc-600" />
              <h2 class="text-xs font-black uppercase tracking-wider text-zinc-900">Status Ruang Periksa</h2>
            </div>
            <span class="text-[10px] font-bold text-zinc-400">Live Status</span>
          </div>

          <div class="space-y-2.5">
            <div
              v-for="(room, idx) in doctorRooms"
              :key="idx"
              class="p-3 rounded-2xl border border-zinc-200/80 bg-[#FAF8F5] flex items-center justify-between gap-2.5"
            >
              <div class="min-w-0">
                <div class="flex items-center gap-1.5">
                  <span class="text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-md bg-white border border-zinc-200 text-zinc-800">
                    {{ room.room_name }}
                  </span>
                  <span
                    :class="[
                      'w-2 h-2 rounded-full',
                      room.status === 'BUSY' ? 'bg-amber-500 animate-pulse' : 'bg-emerald-500'
                    ]"
                  ></span>
                </div>
                <div class="text-xs font-extrabold text-zinc-900 truncate mt-1">
                  {{ room.doctor_name }}
                </div>
              </div>

              <div class="text-right shrink-0">
                <span
                  :class="[
                    'px-2.5 py-1 rounded-full text-[10px] font-bold border font-mono',
                    room.status === 'BUSY'
                      ? 'bg-amber-50 text-amber-900 border-amber-200'
                      : 'bg-emerald-50 text-emerald-900 border-emerald-200'
                  ]"
                >
                  {{ room.status === 'BUSY' ? 'Sedang Periksa' : 'Siap / Kosong' }}
                </span>
                <div class="text-[10px] text-zinc-400 mt-0.5 font-medium">
                  {{ room.waiting_count }} Pasien Antri
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. UPCOMING QUEUE LIST -->
        <div class="bg-white p-5 sm:p-6 rounded-3xl border border-zinc-200/90 shadow-xs flex-1 flex flex-col justify-between space-y-4">
          <div>
            <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
              <div>
                <h2 class="text-xs font-black uppercase tracking-wider text-zinc-900">Antrian Berikutnya</h2>
                <p class="text-[10px] text-zinc-400 font-medium">Harap bersiap di depan pintu ruang periksa</p>
              </div>
              <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-stone-100 text-zinc-700 font-mono">
                {{ monitorData.upcoming_queue.length }} Menunggu
              </span>
            </div>

            <!-- List Items -->
            <div v-if="monitorData.upcoming_queue.length === 0" class="py-8 text-center text-zinc-400 space-y-1">
              <Clock class="w-7 h-7 mx-auto text-zinc-300" />
              <p class="text-xs font-bold text-zinc-600">Tidak ada antrian menunggu</p>
            </div>

            <div v-else class="divide-y divide-zinc-100 mt-1">
              <div
                v-for="(item, idx) in monitorData.upcoming_queue"
                :key="item.id"
                class="py-3 flex items-center justify-between gap-3 hover:bg-stone-50/60 rounded-xl px-1.5 transition-colors"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <span class="w-6 h-6 rounded-full bg-stone-100 text-zinc-600 flex items-center justify-center text-xs font-mono font-bold shrink-0">
                    {{ idx + 1 }}
                  </span>
                  <div class="min-w-0">
                    <div class="font-mono text-base font-black text-zinc-950 leading-tight">
                      {{ item.queue_number }}
                    </div>
                    <div class="text-xs font-bold text-zinc-800 truncate">
                      {{ item.patient_name }}
                    </div>
                  </div>
                </div>

                <div class="text-right shrink-0">
                  <span
                    :class="[
                      'px-2 py-0.5 rounded-full text-[9px] font-extrabold border uppercase tracking-wider',
                      item.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-50 text-amber-900 border-amber-200'
                        : 'bg-rose-50 text-rose-900 border-rose-200'
                    ]"
                  >
                    {{ item.clinic_unit === 'KIDINEST' ? 'Kidinest' : 'SSC' }}
                  </span>
                  <div class="text-[10px] text-zinc-400 font-medium mt-0.5">Menunggu Panggilan</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Bottom Warning / Policy -->
          <div class="p-3 rounded-2xl bg-[#FAF8F5] border border-zinc-200/80 text-[11px] text-zinc-500 font-medium text-center">
            💡 Pasien yang telah dipanggil 3 kali tanpa kehadiran akan dialihkan ke antrian toleransi meja kasir.
          </div>
        </div>

      </div>

    </main>

    <!-- TV FOOTER TICKER -->
    <footer class="w-full pt-4 border-t border-zinc-200/70 flex flex-col sm:flex-row items-center justify-between gap-2 text-xs text-zinc-400">
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
        <span>Sistem Antrian Berkeadilan (FIFO) • Paperless Smart Clinic Management</span>
      </div>
      <div class="font-mono text-[11px] text-zinc-500 font-bold">
        Kidinest Pediatric Care & Semarang Skin Center
      </div>
    </footer>

  </div>
</template>
