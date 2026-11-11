<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Users,
  Clock,
  Search,
  X,
  Stethoscope,
  ChevronRight,
  Sparkles,
  CheckCircle2,
  AlertCircle,
  ArrowRight,
  Filter
} from 'lucide-vue-next'

export interface DoctorInfo {
  id: string
  name: string
  specialty: string
  room: string
  status: 'active' | 'break' | 'off'
}

export interface QueueEntry {
  queue_number: string
  patient_masked: string
  doctor_name: string
  clinic_unit: 'KIDINEST' | 'SSC'
  status: 'CALLING' | 'IN_EXAM' | 'WAITING' | 'PHARMACY' | 'COMPLETED'
  estimated_wait_minutes?: number
  created_at: string
}

const props = defineProps<{
  brand: 'kidinest' | 'ssc'
  doctors: DoctorInfo[]
}>()

const emit = defineEmits<{
  (e: 'register'): void
}>()

// State
const selectedDoctorId = ref<string>('ALL')
const isAllQueuesModalOpen = ref(false)
const searchQuery = ref('')
const statusFilter = ref<'ALL' | 'WAITING' | 'ACTIVE' | 'COMPLETED'>('ALL')

// Mock Live Queues Data for Realistic Experience
const mockQueues = ref<QueueEntry[]>([])

function initMockData() {
  if (props.brand === 'kidinest') {
    mockQueues.value = [
      {
        queue_number: 'K-006',
        patient_masked: 'Ananda R**** (2 th)',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'IN_EXAM',
        created_at: '09:15',
      },
      {
        queue_number: 'K-007',
        patient_masked: 'Ananda M**** (8 bln)',
        doctor_name: 'dr. Anita Ratna, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'CALLING',
        created_at: '09:20',
      },
      {
        queue_number: 'K-008',
        patient_masked: 'Ananda F**** (4 th)',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 8,
        created_at: '09:30',
      },
      {
        queue_number: 'K-009',
        patient_masked: 'Ananda S**** (1 th)',
        doctor_name: 'dr. Anita Ratna, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 15,
        created_at: '09:35',
      },
      {
        queue_number: 'K-010',
        patient_masked: 'Ananda D**** (5 th)',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 25,
        created_at: '09:42',
      },
      {
        queue_number: 'K-011',
        patient_masked: 'Ananda Z**** (3 bln)',
        doctor_name: 'dr. Anita Ratna, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 35,
        created_at: '09:50',
      },
      {
        queue_number: 'K-004',
        patient_masked: 'Ananda K****',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'PHARMACY',
        created_at: '08:45',
      },
      {
        queue_number: 'K-005',
        patient_masked: 'Ananda B****',
        doctor_name: 'dr. Anita Ratna, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'COMPLETED',
        created_at: '08:30',
      },
    ]
  } else {
    // SSC
    mockQueues.value = [
      {
        queue_number: 'S-004',
        patient_masked: 'Ny. Clarissa D****',
        doctor_name: 'dr. Dian Paramita, Sp.DV',
        clinic_unit: 'SSC',
        status: 'IN_EXAM',
        created_at: '10:05',
      },
      {
        queue_number: 'S-005',
        patient_masked: 'Tn. Raymond W****',
        doctor_name: 'dr. Kevin Wijaya, Sp.KK',
        clinic_unit: 'SSC',
        status: 'CALLING',
        created_at: '10:15',
      },
      {
        queue_number: 'S-006',
        patient_masked: 'Nn. Vanessa L****',
        doctor_name: 'dr. Dian Paramita, Sp.DV',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 10,
        created_at: '10:22',
      },
      {
        queue_number: 'S-007',
        patient_masked: 'Ny. Helena T****',
        doctor_name: 'dr. Kevin Wijaya, Sp.KK',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 20,
        created_at: '10:30',
      },
      {
        queue_number: 'S-008',
        patient_masked: 'Nn. Gracia M****',
        doctor_name: 'dr. Dian Paramita, Sp.DV',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 32,
        created_at: '10:45',
      },
      {
        queue_number: 'S-002',
        patient_masked: 'Ny. Jessica S****',
        doctor_name: 'dr. Dian Paramita, Sp.DV',
        clinic_unit: 'SSC',
        status: 'PHARMACY',
        created_at: '09:20',
      },
      {
        queue_number: 'S-003',
        patient_masked: 'Tn. Arya B****',
        doctor_name: 'dr. Kevin Wijaya, Sp.KK',
        clinic_unit: 'SSC',
        status: 'COMPLETED',
        created_at: '09:10',
      },
    ]
  }
}

onMounted(() => {
  initMockData()
})

// Filtered Queues by Selected Doctor
const filteredByDoctor = computed(() => {
  if (selectedDoctorId.value === 'ALL') {
    return mockQueues.value
  }
  const doc = props.doctors.find(d => d.id === selectedDoctorId.value)
  if (!doc) return mockQueues.value
  return mockQueues.value.filter(q => q.doctor_name.toLowerCase().includes(doc.name.toLowerCase().split(',')[0]))
})

// Currently Calling/Examining Queue
const currentActiveQueue = computed(() => {
  const calling = filteredByDoctor.value.find(q => q.status === 'CALLING')
  if (calling) return calling
  const inExam = filteredByDoctor.value.find(q => q.status === 'IN_EXAM')
  if (inExam) return inExam
  return filteredByDoctor.value[0] || null
})

// Next 3 upcoming queues in waiting status
const nextThreeQueues = computed(() => {
  const waitingList = filteredByDoctor.value.filter(q => q.status === 'WAITING')
  return waitingList.slice(0, 3)
})

// Total Waiting Count
const totalWaiting = computed(() => {
  return filteredByDoctor.value.filter(q => q.status === 'WAITING').length
})

// Full Modal Filter
const modalFilteredQueues = computed(() => {
  let list = filteredByDoctor.value

  if (statusFilter.value === 'WAITING') {
    list = list.filter(q => q.status === 'WAITING')
  } else if (statusFilter.value === 'ACTIVE') {
    list = list.filter(q => q.status === 'CALLING' || q.status === 'IN_EXAM')
  } else if (statusFilter.value === 'COMPLETED') {
    list = list.filter(q => q.status === 'COMPLETED' || q.status === 'PHARMACY')
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase().trim()
    list = list.filter(
      item =>
        item.queue_number.toLowerCase().includes(q) ||
        item.patient_masked.toLowerCase().includes(q) ||
        item.doctor_name.toLowerCase().includes(q)
    )
  }

  return list
})
</script>

<template>
  <div class="w-full">

    <!-- Segmented Multi-Doctor Filter Header -->
    <div class="mb-6 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2 mb-1.5">
          <span
            class="w-2.5 h-2.5 rounded-full animate-ping"
            :class="brand === 'kidinest' ? 'bg-emerald-500' : 'bg-rose-500'"
          ></span>
          <span class="text-xs font-black uppercase tracking-wider text-zinc-500">
            Live Queue Board • Sistem Sinkronisasi Real-Time
          </span>
        </div>
        <h3 class="text-2xl sm:text-3xl font-extrabold text-zinc-950 tracking-tight">
          Pantau Nomor Antrian & Dokter Jaga
        </h3>
      </div>

      <!-- Segmented Doctor Pills -->
      <div class="flex items-center flex-wrap gap-1.5 p-1.5 rounded-2xl bg-zinc-200/70 border border-zinc-300/80">
        <button
          @click="selectedDoctorId = 'ALL'"
          class="px-3.5 py-1.5 rounded-xl text-xs font-extrabold transition-all duration-150"
          :class="
            selectedDoctorId === 'ALL'
              ? 'bg-zinc-950 text-white shadow-sm'
              : 'text-zinc-700 hover:text-zinc-950 hover:bg-white/60'
          "
        >
          Semua Dokter
        </button>

        <button
          v-for="doc in doctors"
          :key="doc.id"
          @click="selectedDoctorId = doc.id"
          class="px-3.5 py-1.5 rounded-xl text-xs font-extrabold transition-all duration-150 flex items-center gap-1.5"
          :class="
            selectedDoctorId === doc.id
              ? brand === 'kidinest'
                ? 'bg-[#0F4C3A] text-white shadow-sm'
                : 'bg-zinc-950 text-white shadow-sm'
              : 'text-zinc-700 hover:text-zinc-950 hover:bg-white/60'
          "
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-400"></span>
          <span>{{ doc.name }}</span>
          <span class="text-[10px] opacity-75 hidden xl:inline">({{ doc.room }})</span>
        </button>
      </div>
    </div>

    <!-- Main Grid: Active Card (Left) + Next 3 Preview Cards (Right) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-stretch">

      <!-- Active Calling Ticket Hero Box (Col 7) -->
      <div
        class="lg:col-span-7 rounded-3xl p-6 sm:p-8 flex flex-col justify-between border shadow-sm transition-all"
        :class="
          brand === 'kidinest'
            ? 'bg-white border-emerald-950/15'
            : 'bg-white border-zinc-200'
        "
      >
        <div>
          <!-- Header Tag -->
          <div class="flex items-center justify-between mb-4">
            <span
              class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-black uppercase tracking-wider"
              :class="
                currentActiveQueue?.status === 'IN_EXAM'
                  ? 'bg-emerald-100 text-emerald-900 border border-emerald-200'
                  : 'bg-amber-100 text-amber-900 border border-amber-200'
              "
            >
              <span class="w-2 h-2 rounded-full bg-emerald-600 animate-pulse"></span>
              {{
                currentActiveQueue?.status === 'IN_EXAM'
                  ? 'Sedang Dalam Pemeriksaan'
                  : 'Sedang Dipanggil ke Ruangan'
              }}
            </span>

            <span class="text-xs font-bold text-zinc-500">
              Update: Terkini
            </span>
          </div>

          <!-- Giant Number & Patient Display -->
          <div class="my-3">
            <div
              class="text-6xl sm:text-7xl font-black font-mono tracking-tighter"
              :class="brand === 'kidinest' ? 'text-[#0F4C3A]' : 'text-zinc-950'"
            >
              {{ currentActiveQueue?.queue_number || 'A-001' }}
            </div>
            <div class="text-lg sm:text-xl font-bold text-zinc-800 mt-1">
              {{ currentActiveQueue?.patient_masked || 'Pasien Mandiri' }}
            </div>
          </div>
        </div>

        <!-- Doctor & Room Details -->
        <div class="pt-6 border-t border-zinc-100 mt-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 rounded-xl flex items-center justify-center font-bold text-sm"
              :class="
                brand === 'kidinest'
                  ? 'bg-emerald-50 text-emerald-800 border border-emerald-200'
                  : 'bg-zinc-100 text-zinc-800 border border-zinc-200'
              "
            >
              <Stethoscope class="w-5 h-5" />
            </div>
            <div>
              <div class="text-xs text-zinc-500 font-semibold">Dokter Pemeriksa:</div>
              <div class="text-sm font-extrabold text-zinc-950">
                {{ currentActiveQueue?.doctor_name || 'dr. Jaga Klinik' }}
              </div>
            </div>
          </div>

          <!-- Quick Action to Register -->
          <button
            @click="emit('register')"
            class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl font-black text-xs uppercase tracking-wider text-white transition-all shadow-sm active:scale-95"
            :class="brand === 'kidinest' ? 'bg-[#0F4C3A] hover:bg-[#0A3629]' : 'bg-zinc-950 hover:bg-zinc-800'"
          >
            <span>Ambil Nomor Baru</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      <!-- Next 3 Waiting Queues Preview (Col 5) -->
      <div class="lg:col-span-5 flex flex-col justify-between space-y-4">

        <div class="flex items-center justify-between">
          <div class="text-xs font-black uppercase tracking-wider text-zinc-500 flex items-center gap-1.5">
            <Clock class="w-3.5 h-3.5 text-zinc-600" />
            <span>3 Antrian Berikutnya (Next In Line)</span>
          </div>
          <span class="text-xs font-bold text-zinc-600">
            {{ totalWaiting }} pasien menunggu
          </span>
        </div>

        <!-- 3 Cards Stack -->
        <div class="space-y-3">
          <div
            v-for="(item, idx) in nextThreeQueues"
            :key="item.queue_number"
            class="p-4 rounded-2xl border bg-white hover:border-zinc-400 transition-all shadow-2xs flex items-center justify-between"
            :class="idx === 0 ? 'border-zinc-300 ring-2 ring-zinc-950/5' : 'border-zinc-200'"
          >
            <div class="flex items-center gap-3.5">
              <div
                class="w-11 h-11 rounded-xl flex items-center justify-center font-black font-mono text-base"
                :class="
                  idx === 0
                    ? brand === 'kidinest'
                      ? 'bg-emerald-100 text-[#0F4C3A]'
                      : 'bg-zinc-900 text-white'
                    : 'bg-zinc-100 text-zinc-800'
                "
              >
                {{ item.queue_number }}
              </div>
              <div>
                <div class="text-xs font-extrabold text-zinc-900">
                  {{ item.patient_masked }}
                </div>
                <div class="text-[11px] text-zinc-500 font-medium">
                  {{ item.doctor_name.split(',')[0] }}
                </div>
              </div>
            </div>

            <div class="text-right">
              <span class="inline-flex items-center gap-1 text-[11px] font-bold text-amber-800 bg-amber-50 border border-amber-200 px-2.5 py-0.5 rounded-full">
                <Clock class="w-3 h-3" />
                ~{{ item.estimated_wait_minutes }} mnt
              </span>
            </div>
          </div>

          <div
            v-if="nextThreeQueues.length === 0"
            class="p-8 rounded-2xl border border-dashed border-zinc-300 text-center text-zinc-500 text-xs font-medium"
          >
            Tidak ada pasien lain dalam antrian tunggu saat ini.
          </div>
        </div>

        <!-- Action: Open Drawer All Queues -->
        <button
          @click="isAllQueuesModalOpen = true"
          class="w-full py-3 px-4 rounded-2xl border border-zinc-300 bg-white hover:bg-zinc-50 hover:border-zinc-900 text-zinc-900 text-xs font-black uppercase tracking-wider flex items-center justify-center gap-2 transition-all shadow-2xs group"
        >
          <span>Lihat Semua Antrian Hari Ini ({{ filteredByDoctor.length }})</span>
          <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
        </button>

      </div>

    </div>

    <!-- Notification Tip / WhatsApp Integration Notice -->
    <div class="mt-6 p-4 rounded-2xl bg-zinc-100/90 border border-zinc-200/80 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-xs text-zinc-700">
      <div class="flex items-center gap-2.5">
        <Sparkles class="w-4 h-4 text-amber-500 shrink-0" />
        <span class="font-medium">
          <strong>Tips Pasien:</strong> Anda tidak perlu berdiam di ruang tunggu. Sistem kami mengirim notifikasi WhatsApp otomatis saat antrian Anda berjarak 2 pasien.
        </span>
      </div>
      <button
        @click="emit('register')"
        class="font-extrabold text-zinc-950 underline hover:text-zinc-700 shrink-0"
      >
        Ambil Tiket Online →
      </button>
    </div>

    <!-- Modal Drawer: Lihat Semua Antrian Hari Ini -->
    <div
      v-if="isAllQueuesModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/60 backdrop-blur-xs animate-in fade-in duration-150"
    >
      <div
        class="bg-white rounded-3xl border border-zinc-200 shadow-2xl w-full max-w-2xl max-h-[85vh] flex flex-col overflow-hidden animate-in zoom-in-95 duration-150"
      >
        <!-- Modal Header -->
        <div class="p-6 border-b border-zinc-200 flex items-center justify-between">
          <div>
            <h4 class="text-xl font-black text-zinc-950">
              Daftar Seluruh Antrian Hari Ini
            </h4>
            <p class="text-xs text-zinc-500 font-medium">
              Data real-time poli {{ brand === 'kidinest' ? 'Kidinest' : 'Semarang Skin Center' }} • Nama disamarkan
            </p>
          </div>
          <button
            @click="isAllQueuesModalOpen = false"
            class="w-9 h-9 rounded-full bg-zinc-100 hover:bg-zinc-200 flex items-center justify-center text-zinc-600 transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Filter & Search Inside Modal -->
        <div class="p-4 bg-zinc-50 border-b border-zinc-200 flex flex-col sm:flex-row gap-3">
          <div class="relative flex-1">
            <Search class="w-4 h-4 text-zinc-400 absolute left-3 top-3" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari nomor tiket atau nama..."
              class="w-full pl-9 pr-3 py-2 bg-white rounded-xl border border-zinc-300 text-xs font-bold text-zinc-900 focus:outline-hidden focus:ring-2 focus:ring-zinc-950"
            />
          </div>

          <div class="flex items-center gap-1">
            <button
              @click="statusFilter = 'ALL'"
              class="px-2.5 py-1.5 rounded-lg text-[11px] font-bold"
              :class="statusFilter === 'ALL' ? 'bg-zinc-950 text-white' : 'bg-white border text-zinc-700'"
            >
              Semua
            </button>
            <button
              @click="statusFilter = 'ACTIVE'"
              class="px-2.5 py-1.5 rounded-lg text-[11px] font-bold"
              :class="statusFilter === 'ACTIVE' ? 'bg-zinc-950 text-white' : 'bg-white border text-zinc-700'"
            >
              Diperiksa
            </button>
            <button
              @click="statusFilter = 'WAITING'"
              class="px-2.5 py-1.5 rounded-lg text-[11px] font-bold"
              :class="statusFilter === 'WAITING' ? 'bg-zinc-950 text-white' : 'bg-white border text-zinc-700'"
            >
              Menunggu
            </button>
          </div>
        </div>

        <!-- Queues Table / List -->
        <div class="p-6 overflow-y-auto space-y-2.5 flex-1">
          <div
            v-for="item in modalFilteredQueues"
            :key="item.queue_number"
            class="p-3.5 rounded-2xl border border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100/60 flex items-center justify-between transition-colors"
          >
            <div class="flex items-center gap-3">
              <span class="font-mono font-black text-sm px-2.5 py-1 rounded-lg bg-zinc-200 text-zinc-900">
                {{ item.queue_number }}
              </span>
              <div>
                <div class="text-xs font-extrabold text-zinc-900">
                  {{ item.patient_masked }}
                </div>
                <div class="text-[11px] text-zinc-500 font-medium">
                  {{ item.doctor_name }} • Masuk: {{ item.created_at }}
                </div>
              </div>
            </div>

            <div>
              <span
                v-if="item.status === 'IN_EXAM' || item.status === 'CALLING'"
                class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-emerald-100 text-emerald-900 border border-emerald-300"
              >
                ● Sedang Diperiksa
              </span>
              <span
                v-else-if="item.status === 'WAITING'"
                class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-amber-100 text-amber-900 border border-amber-300"
              >
                Menunggu (~{{ item.estimated_wait_minutes || 10 }}m)
              </span>
              <span
                v-else-if="item.status === 'PHARMACY'"
                class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-sky-100 text-sky-900 border border-sky-300"
              >
                Farmasi / Kasir
              </span>
              <span
                v-else
                class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider bg-zinc-200 text-zinc-600"
              >
                Selesai
              </span>
            </div>
          </div>

          <div
            v-if="modalFilteredQueues.length === 0"
            class="text-center py-12 text-zinc-400 text-xs font-semibold"
          >
            Tidak ada data antrian yang cocok dengan kriteria pencarian.
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-4 border-t border-zinc-200 bg-zinc-50 flex items-center justify-between">
          <span class="text-xs text-zinc-500 font-medium">
            Total {{ modalFilteredQueues.length }} antrian tercatat
          </span>
          <button
            @click="isAllQueuesModalOpen = false"
            class="px-4 py-2 rounded-xl bg-zinc-950 text-white font-extrabold text-xs"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>

  </div>
</template>
