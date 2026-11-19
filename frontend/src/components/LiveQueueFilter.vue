<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Search,
  X,
  ArrowRight
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
        doctor_name: 'dr. Patricia Vanessa A., Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'IN_EXAM',
        created_at: '09:15',
      },
      {
        queue_number: 'K-007',
        patient_masked: 'Ananda M**** (8 bln)',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'CALLING',
        created_at: '09:20',
      },
      {
        queue_number: 'K-008',
        patient_masked: 'Ananda F**** (4 th)',
        doctor_name: 'dr. Patricia Vanessa A., Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 8,
        created_at: '09:30',
      },
      {
        queue_number: 'K-009',
        patient_masked: 'Ananda S**** (1 th)',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 15,
        created_at: '09:35',
      },
      {
        queue_number: 'K-010',
        patient_masked: 'Ananda D**** (5 th)',
        doctor_name: 'dr. Patricia Vanessa A., Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 25,
        created_at: '09:42',
      },
      {
        queue_number: 'K-011',
        patient_masked: 'Ananda Z**** (3 bln)',
        doctor_name: 'dr. Patricia Vanessa A., Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'WAITING',
        estimated_wait_minutes: 35,
        created_at: '09:50',
      },
      {
        queue_number: 'K-004',
        patient_masked: 'Ananda K****',
        doctor_name: 'dr. Patricia Vanessa A., Sp.A',
        clinic_unit: 'KIDINEST',
        status: 'PHARMACY',
        created_at: '08:45',
      },
      {
        queue_number: 'K-005',
        patient_masked: 'Ananda B****',
        doctor_name: 'dr. Hendra Pratama, Sp.A',
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
        doctor_name: 'dr. Dian Paramita, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'IN_EXAM',
        created_at: '10:05',
      },
      {
        queue_number: 'S-005',
        patient_masked: 'Tn. Raymond W****',
        doctor_name: 'dr. Kevin Wijaya, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'CALLING',
        created_at: '10:15',
      },
      {
        queue_number: 'S-006',
        patient_masked: 'Nn. Vanessa L****',
        doctor_name: 'dr. Dian Paramita, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 10,
        created_at: '10:22',
      },
      {
        queue_number: 'S-007',
        patient_masked: 'Ny. Helena T****',
        doctor_name: 'dr. Kevin Wijaya, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 20,
        created_at: '10:30',
      },
      {
        queue_number: 'S-008',
        patient_masked: 'Nn. Gracia M****',
        doctor_name: 'dr. Dian Paramita, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'WAITING',
        estimated_wait_minutes: 32,
        created_at: '10:45',
      },
      {
        queue_number: 'S-002',
        patient_masked: 'Ny. Jessica S****',
        doctor_name: 'dr. Dian Paramita, Sp.D.V.E',
        clinic_unit: 'SSC',
        status: 'PHARMACY',
        created_at: '09:20',
      },
      {
        queue_number: 'S-003',
        patient_masked: 'Tn. Arya B****',
        doctor_name: 'dr. Kevin Wijaya, Sp.D.V.E',
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
  <div class="w-full select-none">

    <!-- 1. Header & Minimalist Doctor Dropdown (Skalabel untuk Banyak Dokter) -->
    <div class="pb-4 border-b flex flex-col sm:flex-row sm:items-center justify-between gap-4"
      :class="brand === 'kidinest' ? 'border-[#E5A93C]/20' : 'border-[#9E606D]/15'"
    >
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span
            class="w-2 h-2 rounded-full animate-pulse"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C]' : 'bg-[#9E606D]'"
          ></span>
          <span class="text-[11px] font-bold uppercase tracking-wider opacity-70">
            Antrian Langsung
          </span>
        </div>
        <h3 class="text-xl sm:text-2xl font-bold tracking-tight">
          Nomor Antrian Sedang Berjalan
        </h3>
      </div>

      <!-- Dropdown Filter Dokter (Clean, Minimalist, Menampung Banyak Dokter) -->
      <div class="flex items-center gap-2 self-start sm:self-auto">
        <label for="doctor-select" class="text-xs font-semibold opacity-70 whitespace-nowrap">
          Pilih Dokter:
        </label>
        <div class="relative">
          <select
            id="doctor-select"
            v-model="selectedDoctorId"
            class="appearance-none text-xs font-bold pl-3 pr-8 py-2 rounded-lg border bg-white cursor-pointer focus:outline-hidden transition-all shadow-2xs"
            :class="
              brand === 'kidinest'
                ? 'border-[#0C3B36]/20 text-[#0C3B36] focus:border-[#0C3B36] focus:ring-1 focus:ring-[#0C3B36]'
                : 'border-[#9E606D]/25 text-[#2D1E24] focus:border-[#9E606D] focus:ring-1 focus:ring-[#9E606D]'
            "
          >
            <option value="ALL">Semua Dokter</option>
            <option
              v-for="doc in doctors"
              :key="doc.id"
              :value="doc.id"
            >
              {{ doc.name }} ({{ doc.room }})
            </option>
          </select>
          <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 opacity-50">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 2. Integrated Hairline Board (Prioritas Utama: Nomor Aktif + 3 Antrian Berikutnya) -->
    <div
      class="py-6 sm:py-8 border-b grid grid-cols-1 lg:grid-cols-12 gap-8 items-start"
      :class="brand === 'kidinest' ? 'border-[#E5A93C]/20' : 'border-[#9E606D]/15'"
    >

      <!-- Kiri: Nomor Antrian yang Sedang Berjalan (Col 7) -->
      <div class="lg:col-span-7 space-y-3">
        
        <!-- Status Label Strip Ringkas -->
        <div class="flex items-center gap-2">
          <span
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-md text-[11px] font-bold"
            :class="
              currentActiveQueue?.status === 'IN_EXAM'
                ? 'bg-emerald-50 text-emerald-800 border border-emerald-200/80'
                : 'bg-amber-50 text-amber-900 border border-amber-200/80'
            "
          >
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-600 animate-pulse"></span>
            {{
              currentActiveQueue?.status === 'IN_EXAM'
                ? 'Sedang Diperiksa'
                : 'Sedang Dipanggil'
            }}
          </span>
        </div>

        <!-- Nomor Antrian Besar Bersih & Nama Pasien -->
        <div class="flex flex-col sm:flex-row sm:items-baseline gap-2 sm:gap-6 pt-1">
          <div
            class="font-mono text-5xl sm:text-6xl font-black tracking-tight"
            :class="brand === 'kidinest' ? 'text-[#0C3B36]' : 'text-[#2D1E24]'"
          >
            {{ currentActiveQueue?.queue_number || '—' }}
          </div>
          <div>
            <div class="text-base sm:text-lg font-bold">
              {{ currentActiveQueue?.patient_masked || 'Tidak ada antrian' }}
            </div>
            <div class="text-xs opacity-75 mt-0.5">
              {{ currentActiveQueue?.doctor_name }}
            </div>
          </div>
        </div>
      </div>

      <!-- Kanan: 3 Antrian Berikutnya (Hairline Rows) (Col 5) -->
      <div
        class="lg:col-span-5 lg:pl-8 lg:border-l space-y-3"
        :class="brand === 'kidinest' ? 'lg:border-[#E5A93C]/20' : 'lg:border-[#9E606D]/15'"
      >
        <div class="flex items-center justify-between text-xs font-semibold opacity-70 pb-1">
          <span>Antrian Berikutnya</span>
          <span>{{ totalWaiting }} menunggu</span>
        </div>

        <!-- Hairline Rows -->
        <div class="divide-y border-t border-b"
          :class="brand === 'kidinest' ? 'divide-[#0C3B36]/10 border-[#0C3B36]/10' : 'divide-[#9E606D]/10 border-[#9E606D]/10'"
        >
          <div
            v-for="(item, idx) in nextThreeQueues"
            :key="item.queue_number"
            class="py-2.5 flex items-center justify-between text-xs"
          >
            <div class="flex items-center gap-3">
              <span
                class="font-mono font-bold px-2 py-0.5 rounded-md text-[11px]"
                :class="
                  idx === 0
                    ? brand === 'kidinest'
                      ? 'bg-[#E5A93C]/20 text-[#0C3B36]'
                      : 'bg-[#9E606D]/15 text-[#9E606D]'
                    : 'bg-zinc-100 text-zinc-700'
                "
              >
                {{ item.queue_number }}
              </span>
              <div>
                <span class="font-bold text-zinc-900">{{ item.patient_masked }}</span>
                <span class="text-[11px] opacity-60 ml-2">({{ item.doctor_name.split(',')[0] }})</span>
              </div>
            </div>

            <span class="text-[11px] font-semibold opacity-70">
              ~{{ item.estimated_wait_minutes }} mnt
            </span>
          </div>

          <div
            v-if="nextThreeQueues.length === 0"
            class="py-4 text-xs opacity-50 italic text-center"
          >
            Tidak ada antrian lain yang menunggu.
          </div>
        </div>

        <!-- Action: Open Drawer All Queues -->
        <div class="pt-1 flex items-center justify-between">
          <button
            @click="isAllQueuesModalOpen = true"
            class="text-xs font-bold underline underline-offset-4 hover:opacity-80 transition-opacity"
          >
            Lihat semua antrian ({{ filteredByDoctor.length }}) →
          </button>
        </div>

      </div>

    </div>



    <!-- Modal Drawer: Seluruh Antrian Hari Ini (Hairline Minimalist) -->
    <div
      v-if="isAllQueuesModalOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs animate-in fade-in duration-100"
    >
      <div
        class="bg-white rounded-2xl border border-zinc-200 shadow-xl w-full max-w-xl max-h-[80vh] flex flex-col overflow-hidden text-zinc-950"
      >
        <!-- Modal Header -->
        <div class="p-5 border-b border-zinc-100 flex items-center justify-between">
          <div>
            <h4 class="text-base font-bold">
              Daftar Lengkap Antrian Hari Ini
            </h4>
            <p class="text-xs text-zinc-500">
              Data real-time poli {{ brand === 'kidinest' ? 'Kidinest Pediatric' : 'Semarang Skin Center' }}
            </p>
          </div>
          <button
            @click="isAllQueuesModalOpen = false"
            class="p-1 rounded-lg hover:bg-zinc-100 text-zinc-500 transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Filter & Search (Clean Bar) -->
        <div class="p-3 bg-zinc-50/80 border-b border-zinc-100 flex items-center gap-2 text-xs">
          <div class="relative flex-1">
            <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-2.5" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Cari nomor atau nama pasien..."
              class="w-full pl-8 pr-3 py-1.5 bg-white rounded-lg border border-zinc-200 text-xs font-semibold focus:outline-hidden focus:border-zinc-900"
            />
          </div>
          <button
            @click="statusFilter = 'ALL'"
            class="px-2.5 py-1.5 rounded-md font-bold text-[11px]"
            :class="statusFilter === 'ALL' ? 'bg-zinc-900 text-white' : 'bg-white border text-zinc-600'"
          >
            Semua
          </button>
          <button
            @click="statusFilter = 'ACTIVE'"
            class="px-2.5 py-1.5 rounded-md font-bold text-[11px]"
            :class="statusFilter === 'ACTIVE' ? 'bg-zinc-900 text-white' : 'bg-white border text-zinc-600'"
          >
            Dipanggil
          </button>
          <button
            @click="statusFilter = 'WAITING'"
            class="px-2.5 py-1.5 rounded-md font-bold text-[11px]"
            :class="statusFilter === 'WAITING' ? 'bg-zinc-900 text-white' : 'bg-white border text-zinc-600'"
          >
            Menunggu
          </button>
        </div>

        <!-- Table List (Hairline, Bebas Kotak Card) -->
        <div class="p-4 overflow-y-auto divide-y divide-zinc-100 flex-1">
          <div
            v-for="item in modalFilteredQueues"
            :key="item.queue_number"
            class="py-3 flex items-center justify-between text-xs"
          >
            <div class="flex items-center gap-3">
              <span class="font-mono font-bold px-2 py-0.5 bg-zinc-100 rounded text-zinc-800 text-[11px]">
                {{ item.queue_number }}
              </span>
              <div>
                <div class="font-bold text-zinc-900">{{ item.patient_masked }}</div>
                <div class="text-[11px] text-zinc-500">{{ item.doctor_name.split(',')[0] }} • {{ item.created_at }}</div>
              </div>
            </div>

            <div>
              <span
                v-if="item.status === 'IN_EXAM' || item.status === 'CALLING'"
                class="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-100 text-emerald-800"
              >
                ● Diperiksa
              </span>
              <span
                v-else-if="item.status === 'WAITING'"
                class="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-100 text-amber-900"
              >
                Menunggu (~{{ item.estimated_wait_minutes || 10 }}m)
              </span>
              <span
                v-else
                class="px-2 py-0.5 rounded text-[10px] font-bold bg-zinc-100 text-zinc-600"
              >
                Selesai
              </span>
            </div>
          </div>

          <div
            v-if="modalFilteredQueues.length === 0"
            class="py-8 text-center text-xs text-zinc-400 italic"
          >
            Tidak ada nomor antrian yang cocok.
          </div>
        </div>

        <!-- Footer -->
        <div class="p-3 border-t border-zinc-100 bg-zinc-50/50 flex items-center justify-between text-xs text-zinc-500">
          <span>Total {{ modalFilteredQueues.length }} antrian</span>
          <button
            @click="isAllQueuesModalOpen = false"
            class="px-3 py-1.5 rounded-lg bg-zinc-900 text-white font-bold text-xs"
          >
            Tutup
          </button>
        </div>

      </div>
    </div>

  </div>
</template>
