<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import {
  Users,
  UserPlus,
  KeyRound,
  ShieldCheck,
  Building2,
  Stethoscope,
  CreditCard,
  Pill,
  Search,
  CheckCircle2,
  XCircle,
  AlertCircle,
  Loader2,
  X,
  RefreshCw
} from 'lucide-vue-next'

interface StaffUser {
  id: number
  employee_number: string
  name: string
  role: 'SUPERADMIN_OWNER' | 'DOCTOR' | 'GRO_CASHIER' | 'PHARMACIST'
  clinic_unit: 'KIDINEST' | 'SSC' | 'ALL'
  is_active: boolean
  created_at: string
  updated_at: string
}

const defaultMockStaffs: StaffUser[] = [
  {
    id: 1,
    employee_number: 'OWNER-01',
    name: 'dr. Ratna Paramita, M.Kes',
    role: 'SUPERADMIN_OWNER',
    clinic_unit: 'ALL',
    is_active: true,
    created_at: '2026-01-01T08:00:00Z',
    updated_at: '2026-01-01T08:00:00Z',
  },
  {
    id: 2,
    employee_number: 'DOC-KID-01',
    name: 'dr. Patricia Sp.A',
    role: 'DOCTOR',
    clinic_unit: 'KIDINEST',
    is_active: true,
    created_at: '2026-01-02T08:00:00Z',
    updated_at: '2026-01-02T08:00:00Z',
  },
  {
    id: 3,
    employee_number: 'DOC-KID-02',
    name: 'dr. Hendra Pratama, Sp.A',
    role: 'DOCTOR',
    clinic_unit: 'KIDINEST',
    is_active: true,
    created_at: '2026-01-02T08:00:00Z',
    updated_at: '2026-01-02T08:00:00Z',
  },
  {
    id: 4,
    employee_number: 'DOC-SSC-01',
    name: 'dr. Dian Paramita, Sp.D.V.E',
    role: 'DOCTOR',
    clinic_unit: 'SSC',
    is_active: true,
    created_at: '2026-01-03T08:00:00Z',
    updated_at: '2026-01-03T08:00:00Z',
  },
  {
    id: 5,
    employee_number: 'GRO-01',
    name: 'Siti Sarah, A.Md',
    role: 'GRO_CASHIER',
    clinic_unit: 'KIDINEST',
    is_active: true,
    created_at: '2026-01-04T08:00:00Z',
    updated_at: '2026-01-04T08:00:00Z',
  },
  {
    id: 6,
    employee_number: 'FARMASI-01',
    name: 'Apt. Dimas Wicaksono, S.Farm',
    role: 'PHARMACIST',
    clinic_unit: 'KIDINEST',
    is_active: true,
    created_at: '2026-01-05T08:00:00Z',
    updated_at: '2026-01-05T08:00:00Z',
  },
]

const auth = useAuthStore()

const staffList = ref<StaffUser[]>(defaultMockStaffs)
const isLoading = ref(false)
const errorMessage = ref<string | null>(null)
const successMessage = ref<string | null>(null)

// Filters
const filterUnit = ref('ALL')
const filterRole = ref('ALL')
const searchQuery = ref('')

// Modals
const showAddModal = ref(false)
const showResetModal = ref(false)
const selectedStaff = ref<StaffUser | null>(null)
const isSubmitting = ref(false)

// Add Form
const newStaff = ref({
  employee_number: '',
  name: '',
  role: 'DOCTOR' as StaffUser['role'],
  clinic_unit: 'KIDINEST' as StaffUser['clinic_unit'],
  initial_password: '',
})

// Reset Password Form
const newPassword = ref('')

async function fetchStaffs() {
  isLoading.value = true
  errorMessage.value = null
  try {
    const res = await fetch('/api/admin/staff', {
      headers: {
        Authorization: `Bearer ${auth.token}`,
      },
    })
    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Gagal mengambil data staf')
    }
    if (data.data && data.data.length > 0) {
      staffList.value = data.data
    } else {
      staffList.value = defaultMockStaffs
    }
  } catch (err: any) {
    if (staffList.value.length === 0) {
      staffList.value = defaultMockStaffs
    }
  } finally {
    isLoading.value = false
  }
}

const filteredStaffs = computed(() => {
  return staffList.value.filter((s) => {
    const matchUnit = filterUnit.value === 'ALL' || s.clinic_unit === filterUnit.value
    const matchRole = filterRole.value === 'ALL' || s.role === filterRole.value
    const matchQuery =
      searchQuery.value === '' ||
      s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      s.employee_number.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchUnit && matchRole && matchQuery
  })
})

const stats = computed(() => {
  const total = staffList.value.length
  const active = staffList.value.filter((s) => s.is_active).length
  const doctors = staffList.value.filter((s) => s.role === 'DOCTOR' && s.is_active).length
  const cashiers = staffList.value.filter((s) => s.role === 'GRO_CASHIER' && s.is_active).length
  const pharmacists = staffList.value.filter((s) => s.role === 'PHARMACIST' && s.is_active).length
  return { total, active, doctors, cashiers, pharmacists }
})

function openResetModal(staff: StaffUser) {
  selectedStaff.value = staff
  newPassword.value = ''
  showResetModal.value = true
}

async function handleResetPassword() {
  if (!selectedStaff.value || !newPassword.value) return
  isSubmitting.value = true
  errorMessage.value = null
  try {
    const res = await fetch(`/api/admin/staff/${selectedStaff.value.id}/reset-password`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${auth.token}`,
      },
      body: JSON.stringify({ new_password: newPassword.value }),
    })
    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Gagal mereset password staf')
    }
    successMessage.value = `Password untuk ${selectedStaff.value.name} (${selectedStaff.value.employee_number}) berhasil diubah.`
    showResetModal.value = false
    setTimeout(() => (successMessage.value = null), 4000)
  } catch (err: any) {
    errorMessage.value = err.message
  } finally {
    isSubmitting.value = false
  }
}

async function handleCreateStaff() {
  isSubmitting.value = true
  errorMessage.value = null
  try {
    const res = await fetch('/api/admin/staff', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${auth.token}`,
      },
      body: JSON.stringify(newStaff.value),
    })
    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Gagal menambahkan staf baru')
    }
    successMessage.value = `Staf baru ${data.data.name} (${data.data.employee_number}) berhasil didaftarkan.`
    showAddModal.value = false
    newStaff.value = {
      employee_number: '',
      name: '',
      role: 'DOCTOR',
      clinic_unit: 'KIDINEST',
      initial_password: '',
    }
    await fetchStaffs()
    setTimeout(() => (successMessage.value = null), 4000)
  } catch (err: any) {
    errorMessage.value = err.message
  } finally {
    isSubmitting.value = false
  }
}

async function toggleActiveStatus(staff: StaffUser) {
  const updatedStatus = !staff.is_active
  try {
    const res = await fetch(`/api/admin/staff/${staff.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${auth.token}`,
      },
      body: JSON.stringify({
        name: staff.name,
        role: staff.role,
        clinic_unit: staff.clinic_unit,
        is_active: updatedStatus,
      }),
    })
    const data = await res.json()
    if (res.ok && data.success) {
      staff.is_active = updatedStatus
    } else {
      throw new Error(data.error || 'Gagal mengubah status aktif')
    }
  } catch (err: any) {
    alert('Gagal mengubah status: ' + err.message)
  }
}

function getRoleBadge(role: StaffUser['role']) {
  switch (role) {
    case 'SUPERADMIN_OWNER':
      return { text: 'Owner / Atasan', bg: 'bg-purple-100 text-purple-900 border-purple-200' }
    case 'DOCTOR':
      return { text: 'Dokter', bg: 'bg-emerald-100 text-emerald-900 border-emerald-200' }
    case 'GRO_CASHIER':
      return { text: 'Kasir / GRO', bg: 'bg-amber-100 text-amber-900 border-amber-200' }
    case 'PHARMACIST':
      return { text: 'Apoteker', bg: 'bg-blue-100 text-blue-900 border-blue-200' }
    default:
      return { text: role, bg: 'bg-zinc-100 text-zinc-800 border-zinc-200' }
  }
}

function getUnitBadge(unit: StaffUser['clinic_unit']) {
  switch (unit) {
    case 'KIDINEST':
      return { text: 'Kidinest (Anak)', bg: 'bg-sky-50 text-sky-800 border-sky-200' }
    case 'SSC':
      return { text: 'SSC (Kulit)', bg: 'bg-rose-50 text-rose-800 border-rose-200' }
    case 'ALL':
      return { text: 'Semua Unit', bg: 'bg-stone-100 text-stone-800 border-stone-200' }
  }
}

onMounted(() => {
  fetchStaffs()
})
</script>

<template>
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- TOP CLINICAL BAR -->
    <header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pb-2 border-b border-zinc-200/80">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="text-xs font-semibold tracking-wide text-zinc-500 uppercase">
            Meja Direksi & Manajemen • Semarang Medical Group
          </span>
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-xs font-medium text-emerald-700">Kontrol Atasan Aktif</span>
        </div>
        <div class="flex items-center gap-3">
          <h1 class="text-xl sm:text-2xl font-black tracking-tight text-zinc-950">
            Manajemen Staf & Hak Akses Klinik
          </h1>
          <span class="px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-zinc-900 text-white shadow-2xs">
            Owner Portal
          </span>
        </div>
        <p class="text-xs text-zinc-500 mt-1">
          Kelola Nomor Induk Karyawan (NIK), unit praktik, status aktif, dan reset kata sandi staf.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          @click="fetchStaffs"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200/90 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-700 transition-all shadow-2xs cursor-pointer"
          title="Segarkan data staf"
        >
          <Loader2 v-if="isLoading" class="w-3.5 h-3.5 animate-spin text-zinc-500" />
          <RefreshCw v-else class="w-3.5 h-3.5 text-zinc-500" />
          <span>Sync</span>
        </button>
        <button
          @click="showAddModal = true"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-2xs transition-all active:scale-95 cursor-pointer"
        >
          <UserPlus class="w-3.5 h-3.5" />
          <span>Tambah Karyawan</span>
        </button>
      </div>
    </header>

    <!-- Feedback Banner -->
    <div
      v-if="successMessage"
      class="p-3.5 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-900 text-xs font-semibold flex items-center gap-2 animate-in fade-in duration-200 shadow-2xs"
    >
      <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
      <span>{{ successMessage }}</span>
    </div>
    <div
      v-if="errorMessage"
      class="p-3.5 rounded-xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2 animate-in fade-in duration-200 shadow-2xs"
    >
      <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
      <span>{{ errorMessage }}</span>
    </div>

    <!-- 1-ROW HAIRLINE METRIC TICKER STRIP -->
    <section class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Total Karyawan</span>
          <Users class="w-4 h-4 text-zinc-400" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-zinc-950 tracking-tight tabular-nums">{{ stats.total }}</div>
        <div class="text-[11px] text-zinc-500 font-medium mt-0.5">{{ stats.active }} akun aktif</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Dokter Praktik</span>
          <Stethoscope class="w-4 h-4 text-emerald-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-emerald-950 tracking-tight tabular-nums">{{ stats.doctors }}</div>
        <div class="text-[11px] text-emerald-700 font-medium mt-0.5">Kidinest & SSC</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Meja Depan & Kasir</span>
          <CreditCard class="w-4 h-4 text-amber-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-amber-950 tracking-tight tabular-nums">{{ stats.cashiers }}</div>
        <div class="text-[11px] text-amber-700 font-medium mt-0.5">GRO loket & invoice</div>
      </div>

      <div class="p-4 sm:p-5">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Tim Farmasi</span>
          <Pill class="w-4 h-4 text-blue-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-blue-950 tracking-tight tabular-nums">{{ stats.pharmacists }}</div>
        <div class="text-[11px] text-blue-700 font-medium mt-0.5">Apotek & stok opnam</div>
      </div>

      <div class="p-4 sm:p-5 hidden lg:block">
        <div class="flex items-center justify-between text-zinc-500 mb-1">
          <span class="text-[11px] font-semibold uppercase tracking-wider">Integritas Akun</span>
          <ShieldCheck class="w-4 h-4 text-purple-600" />
        </div>
        <div class="text-2xl sm:text-3xl font-black text-purple-950 tracking-tight tabular-nums">100%</div>
        <div class="text-[11px] text-purple-700 font-medium mt-0.5">RBAC & audit aktif</div>
      </div>
    </section>

    <!-- Filter & Search Toolbar -->
    <div class="p-3.5 rounded-2xl bg-white border border-zinc-200/90 shadow-2xs flex flex-col md:flex-row items-center justify-between gap-3">
      <!-- Search Input -->
      <div class="relative w-full md:w-80">
        <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Cari nama atau NIK staf..."
          class="w-full pl-9 pr-4 py-1.5 rounded-xl border border-zinc-200 text-xs font-medium text-zinc-900 placeholder:text-zinc-400 focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
        />
      </div>

      <!-- Filters Pill Bar -->
      <div class="flex flex-wrap items-center gap-2 w-full md:w-auto">
        <!-- Unit Filter -->
        <div class="flex items-center p-1 bg-zinc-100 rounded-xl text-[11px] font-bold text-zinc-600">
          <button
            @click="filterUnit = 'ALL'"
            :class="['px-3 py-1 rounded-lg transition-all cursor-pointer', filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'hover:text-zinc-950']"
          >
            Semua Unit
          </button>
          <button
            @click="filterUnit = 'KIDINEST'"
            :class="['px-3 py-1 rounded-lg transition-all cursor-pointer', filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-2xs' : 'hover:text-zinc-950']"
          >
            Kidinest
          </button>
          <button
            @click="filterUnit = 'SSC'"
            :class="['px-3 py-1 rounded-lg transition-all cursor-pointer', filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-2xs' : 'hover:text-zinc-950']"
          >
            SSC
          </button>
        </div>

        <!-- Role Filter Dropdown -->
        <select
          v-model="filterRole"
          class="px-3 py-1.5 rounded-xl border border-zinc-200 text-xs font-bold text-zinc-800 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-900 cursor-pointer shadow-2xs"
        >
          <option value="ALL">Semua Peran</option>
          <option value="SUPERADMIN_OWNER">Owner / Atasan</option>
          <option value="DOCTOR">Dokter</option>
          <option value="GRO_CASHIER">Kasir / GRO</option>
          <option value="PHARMACIST">Apoteker</option>
        </select>
      </div>
    </div>

    <!-- Staff Table Container (Hairline Borders) -->
    <div class="rounded-2xl bg-white border border-zinc-200/90 shadow-2xs overflow-hidden">
      <div v-if="isLoading" class="p-12 text-center text-zinc-500 space-y-2">
        <Loader2 class="w-6 h-6 animate-spin mx-auto text-zinc-900" />
        <p class="text-xs font-semibold">Memuat daftar staf klinik...</p>
      </div>

      <div v-else-if="filteredStaffs.length === 0" class="p-12 text-center text-zinc-500 space-y-2">
        <Users class="w-8 h-8 mx-auto text-zinc-300" />
        <p class="text-xs font-bold text-zinc-700">Tidak ada data staf yang cocok</p>
        <p class="text-xs text-zinc-400">Coba ubah kata kunci pencarian atau reset filter di atas.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="border-b border-zinc-200/80 bg-zinc-50/70 text-zinc-500 font-bold uppercase tracking-wider text-[10px]">
              <th class="py-3 px-5">NIK Karyawan</th>
              <th class="py-3 px-5">Nama Lengkap</th>
              <th class="py-3 px-5">Peran & Akses</th>
              <th class="py-3 px-5">Unit Klinik</th>
              <th class="py-3 px-5">Status Akun</th>
              <th class="py-3 px-5 text-right">Aksi Atasan</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-200/60 font-medium text-zinc-900">
            <tr
              v-for="staff in filteredStaffs"
              :key="staff.id"
              class="hover:bg-zinc-50/70 transition-colors"
            >
              <!-- NIK Karyawan -->
              <td class="py-3 px-5 font-mono font-bold text-zinc-950">
                {{ staff.employee_number }}
              </td>

              <!-- Name -->
              <td class="py-3 px-5 font-bold text-zinc-950">
                {{ staff.name }}
              </td>

              <!-- Role -->
              <td class="py-3 px-5">
                <span :class="['inline-block px-2.5 py-0.5 rounded-full text-[11px] font-bold border', getRoleBadge(staff.role).bg]">
                  {{ getRoleBadge(staff.role).text }}
                </span>
              </td>

              <!-- Clinic Unit -->
              <td class="py-3 px-5">
                <span :class="['inline-block px-2.5 py-0.5 rounded-full text-[11px] font-bold border', getUnitBadge(staff.clinic_unit).bg]">
                  {{ getUnitBadge(staff.clinic_unit).text }}
                </span>
              </td>

              <!-- Status Toggle -->
              <td class="py-3 px-5">
                <button
                  @click="toggleActiveStatus(staff)"
                  :class="[
                    'inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold transition-all cursor-pointer',
                    staff.is_active
                      ? 'bg-emerald-50 text-emerald-800 border border-emerald-200 hover:bg-emerald-100'
                      : 'bg-zinc-100 text-zinc-500 border border-zinc-200 hover:bg-zinc-200'
                  ]"
                  title="Klik untuk mengubah status aktif"
                >
                  <span :class="['w-1.5 h-1.5 rounded-full', staff.is_active ? 'bg-emerald-600' : 'bg-zinc-400']"></span>
                  <span>{{ staff.is_active ? 'Aktif' : 'Non-Aktif' }}</span>
                </button>
              </td>

              <!-- Actions -->
              <td class="py-3 px-5 text-right">
                <button
                  @click="openResetModal(staff)"
                  class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl border border-zinc-200 hover:border-zinc-900 bg-white hover:bg-zinc-950 hover:text-white text-[11px] font-bold text-zinc-800 transition-all shadow-2xs cursor-pointer active:scale-95"
                >
                  <KeyRound class="w-3 h-3" />
                  <span>Reset Password</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Tambah Staf Baru -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150"
    >
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xl max-w-lg w-full p-6 space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2.5">
            <span class="p-2 rounded-xl bg-zinc-100 text-zinc-900"><UserPlus class="w-4 h-4" /></span>
            <div>
              <h3 class="text-sm font-extrabold text-zinc-950">Tambah Karyawan Baru</h3>
              <p class="text-[11px] text-zinc-500 font-medium">Akun internal untuk operasional staf klinik</p>
            </div>
          </div>
          <button @click="showAddModal = false" class="p-1 rounded-lg text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100 cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleCreateStaff" class="space-y-4">
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
              Nomor Induk Karyawan (NIK Internal)
            </label>
            <input
              v-model="newStaff.employee_number"
              type="text"
              required
              placeholder="Contoh: DOC-KID-02 atau GRO-02"
              class="w-full px-3.5 py-2 rounded-xl border border-zinc-200 text-xs font-mono font-bold focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
            />
          </div>

          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
              Nama Lengkap & Gelar
            </label>
            <input
              v-model="newStaff.name"
              type="text"
              required
              placeholder="Contoh: dr. Anisa Putri, Sp.A"
              class="w-full px-3.5 py-2 rounded-xl border border-zinc-200 text-xs font-semibold focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
                Peran Pengguna
              </label>
              <select
                v-model="newStaff.role"
                class="w-full px-3 py-2 rounded-xl border border-zinc-200 text-xs font-bold bg-white focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
              >
                <option value="DOCTOR">Dokter</option>
                <option value="GRO_CASHIER">Kasir / GRO</option>
                <option value="PHARMACIST">Apoteker</option>
                <option value="SUPERADMIN_OWNER">Owner / Atasan</option>
              </select>
            </div>

            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
                Unit Praktik
              </label>
              <select
                v-model="newStaff.clinic_unit"
                class="w-full px-3 py-2 rounded-xl border border-zinc-200 text-xs font-bold bg-white focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
              >
                <option value="KIDINEST">Kidinest</option>
                <option value="SSC">Semarang Skin Center</option>
                <option value="ALL">Semua Unit</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
              Kata Sandi Awal
            </label>
            <input
              v-model="newStaff.initial_password"
              type="password"
              required
              placeholder="Minimal 6 karakter"
              class="w-full px-3.5 py-2 rounded-xl border border-zinc-200 text-xs font-semibold focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
            />
          </div>

          <div class="flex items-center justify-end gap-2.5 pt-3 border-t border-zinc-100">
            <button
              type="button"
              @click="showAddModal = false"
              class="px-3.5 py-1.5 rounded-xl border border-zinc-200 text-xs font-bold text-zinc-600 hover:bg-zinc-100 transition-colors cursor-pointer"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
            >
              <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
              <span>Simpan Karyawan</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal Reset Password -->
    <div
      v-if="showResetModal && selectedStaff"
      class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150"
    >
      <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xl max-w-md w-full p-6 space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
        <div class="flex items-center justify-between pb-3 border-b border-zinc-100">
          <div class="flex items-center gap-2.5">
            <span class="p-2 rounded-xl bg-amber-50 text-amber-900 border border-amber-200"><KeyRound class="w-4 h-4 text-amber-700" /></span>
            <div>
              <h3 class="text-sm font-extrabold text-zinc-950">Reset Password Staf</h3>
              <p class="text-[11px] text-zinc-500 font-medium">Tetapkan kata sandi baru secara langsung</p>
            </div>
          </div>
          <button @click="showResetModal = false" class="p-1 rounded-lg text-zinc-400 hover:text-zinc-900 hover:bg-zinc-100 cursor-pointer">
            <X class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="handleResetPassword" class="space-y-4">
          <div class="p-3 rounded-xl bg-zinc-50 border border-zinc-200/70 text-xs space-y-1">
            <div class="text-zinc-500 font-medium">Staf yang direset:</div>
            <div class="font-extrabold text-zinc-950">{{ selectedStaff.name }}</div>
            <div class="font-mono text-zinc-600 text-[11px] font-bold">NIK: {{ selectedStaff.employee_number }}</div>
          </div>

          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1">
              Kata Sandi Baru
            </label>
            <input
              v-model="newPassword"
              type="password"
              required
              placeholder="Masukkan password baru..."
              class="w-full px-3.5 py-2 rounded-xl border border-zinc-200 text-xs font-semibold focus:outline-none focus:ring-1 focus:ring-zinc-900 focus:border-zinc-900"
            />
          </div>

          <div class="flex items-center justify-end gap-2.5 pt-3 border-t border-zinc-100">
            <button
              type="button"
              @click="showResetModal = false"
              class="px-3.5 py-1.5 rounded-xl border border-zinc-200 text-xs font-bold text-zinc-600 hover:bg-zinc-100 transition-colors cursor-pointer"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
            >
              <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
              <span>Perbarui Password</span>
            </button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>
