<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useClinicStore, type RoleType } from '../stores/clinic'
import {
  ShieldCheck,
  Stethoscope,
  Pill,
  CreditCard,
  Building2,
  ArrowLeft,
  KeyRound,
  AlertCircle,
  Loader2,
  Sparkles,
  Baby
} from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  brand?: 'kidinest' | 'ssc'
}>(), {
  brand: 'kidinest'
})

const emit = defineEmits<{
  (e: 'navigate', view: string): void
}>()

const auth = useAuthStore()
const clinicStore = useClinicStore()

const employeeNumber = ref('OWNER-01')
const password = ref('klinik123')
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

interface QuickAccount {
  empNum: string
  name: string
  title: string
  role: 'SUPERADMIN_OWNER' | 'DOCTOR' | 'GRO_CASHIER' | 'PHARMACIST'
  unit: string
  clinicRole: RoleType
  icon: any
}

// Quick demo accounts tailored by floor context
const kidinestAccounts: QuickAccount[] = [
  {
    empNum: 'OWNER-01',
    name: 'dr. Ratna Paramita, M.Kes',
    title: 'Owner & Atasan Klinik (Semua Unit)',
    role: 'SUPERADMIN_OWNER',
    unit: 'Semarang Medical Group',
    clinicRole: 'owner',
    icon: ShieldCheck,
  },
  {
    empNum: 'DOC-KID-01',
    name: 'dr. Patricia Vanessa A., Sp.A',
    title: 'Dokter Spesialis Anak • Poli 1',
    role: 'DOCTOR',
    unit: 'Kidinest (Lantai 1)',
    clinicRole: 'doctor',
    icon: Stethoscope,
  },
  {
    empNum: 'DOC-KID-02',
    name: 'dr. Hendra Pratama, Sp.A',
    title: 'Konsultan Tumbang & Nutrisi • Poli 2',
    role: 'DOCTOR',
    unit: 'Kidinest (Lantai 1)',
    clinicRole: 'doctor',
    icon: Stethoscope,
  },
  {
    empNum: 'GRO-01',
    name: 'Siti Sarah, A.Md',
    title: 'Kasir & Meja Depan GRO',
    role: 'GRO_CASHIER',
    unit: 'Kidinest (Lantai 1)',
    clinicRole: 'cashier',
    icon: CreditCard,
  },
  {
    empNum: 'FARMASI-01',
    name: 'Apt. Dimas Wicaksono, S.Farm',
    title: 'Apoteker Penanggung Jawab',
    role: 'PHARMACIST',
    unit: 'Farmasi Pediatri',
    clinicRole: 'pharmacy',
    icon: Pill,
  },
]

const sscAccounts: QuickAccount[] = [
  {
    empNum: 'OWNER-01',
    name: 'dr. Ratna Paramita, M.Kes',
    title: 'Owner & Atasan Klinik (Semua Unit)',
    role: 'SUPERADMIN_OWNER',
    unit: 'Semarang Medical Group',
    clinicRole: 'owner',
    icon: ShieldCheck,
  },
  {
    empNum: 'DOC-SSC-01',
    name: 'dr. Dian Paramita, Sp.D.V.E',
    title: 'Spesialis Dermatologi & Laser • Ruang 1',
    role: 'DOCTOR',
    unit: 'SSC (Lantai 2)',
    clinicRole: 'doctor',
    icon: Sparkles,
  },
  {
    empNum: 'DOC-SSC-02',
    name: 'dr. Kevin Wijaya, Sp.D.V.E',
    title: 'Spesialis Kulit & Bedah Minor • Ruang 2',
    role: 'DOCTOR',
    unit: 'SSC (Lantai 2)',
    clinicRole: 'doctor',
    icon: Building2,
  },
  {
    empNum: 'GRO-01',
    name: 'Amanda L., S.E',
    title: 'Kasir & Resepsionis VIP',
    role: 'GRO_CASHIER',
    unit: 'SSC (Lantai 2)',
    clinicRole: 'cashier',
    icon: CreditCard,
  },
  {
    empNum: 'FARMASI-01',
    name: 'Apt. Dimas Wicaksono, S.Farm',
    title: 'Apoteker & Resep Racikan Kulit',
    role: 'PHARMACIST',
    unit: 'Apotek Dermatologi',
    clinicRole: 'pharmacy',
    icon: Pill,
  },
]

const quickAccounts = computed(() => {
  return props.brand === 'ssc' ? sscAccounts : kidinestAccounts
})

async function submitLogin(empNum: string, pass: string, targetClinicRole?: RoleType) {
  isSubmitting.value = true
  errorMessage.value = null

  // 1. Coba login ke API backend
  const success = await auth.login(empNum, pass)
  if (success && auth.user) {
    let roleTarget: RoleType = 'doctor'
    switch (auth.user.role) {
      case 'SUPERADMIN_OWNER':
        roleTarget = 'owner'
        break
      case 'DOCTOR':
        roleTarget = 'doctor'
        break
      case 'GRO_CASHIER':
        roleTarget = 'cashier'
        break
      case 'PHARMACIST':
        roleTarget = 'pharmacy'
        break
    }
    clinicStore.setRole(roleTarget)
    emit('navigate', 'admin')
  } else {
    // 2. Demo fallback
    if (auth.error) {
      errorMessage.value = auth.error
    }
    if (targetClinicRole) {
      clinicStore.setRole(targetClinicRole)
      emit('navigate', 'admin')
    }
  }
  isSubmitting.value = false
}

function selectQuickAccount(acc: QuickAccount) {
  employeeNumber.value = acc.empNum
  password.value = 'klinik123'
  submitLogin(acc.empNum, 'klinik123', acc.clinicRole)
}

function handleFormSubmit() {
  submitLogin(employeeNumber.value, password.value)
}
</script>

<template>
  <div
    class="min-h-screen flex flex-col justify-between p-6 sm:p-12 font-['Plus_Jakarta_Sans',sans-serif] select-none transition-colors"
    :class="brand === 'kidinest' ? 'bg-[#FAF7F2] text-[#0C3B36]' : 'bg-[#FAF5F5] text-[#2D1E24]'"
  >
    <!-- Top Back Link -->
    <div class="max-w-5xl mx-auto w-full">
      <button
        @click="emit('navigate', brand)"
        class="inline-flex items-center gap-2 text-xs font-bold opacity-70 hover:opacity-100 transition-opacity"
      >
        <ArrowLeft class="w-3.5 h-3.5" />
        <span>Kembali ke Beranda {{ brand === 'kidinest' ? 'Kidinest (Lantai 1)' : 'Semarang Skin Center (Lantai 2)' }}</span>
      </button>
    </div>

    <!-- Center Login Layout -->
    <main class="max-w-5xl mx-auto w-full py-8">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-10 items-start">

        <!-- Left: Context & Quick-Fill Demo (7 cols) -->
        <div class="md:col-span-7 space-y-5">
          <div
            class="inline-flex items-center gap-2 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider shadow-2xs"
            :class="brand === 'kidinest' ? 'bg-[#E5A93C] text-[#0C3B36]' : 'bg-[#9E606D]/15 text-[#9E606D] border border-[#9E606D]/25'"
          >
            <img
              v-if="brand === 'kidinest'"
              src="/assets/kidinest/logo-maskot.png"
              class="w-4 h-4 object-contain"
            />
            <img
              v-else
              src="/assets/ssc/logo-ssc.png"
              class="w-3.5 h-3.5 object-contain"
            />
            <span>{{ brand === 'kidinest' ? 'Kidinest Pediatric Care • Lantai 1' : 'Semarang Skin Center • Lantai 2' }}</span>
          </div>

          <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight leading-tight">
            Portal Masuk Staf & Dokter.
          </h1>
          <p class="text-xs sm:text-sm opacity-75 leading-relaxed font-medium">
            {{
              brand === 'kidinest'
                ? 'Autentikasi staf medis poli anak Lantai 1. Terintegrasi satu database terpusat di bawah supervisi Atasan Klinik Utama Semarang Medical Group.'
                : 'Autentikasi staf medis dermatologi & laser Lantai 2. Terintegrasi satu database terpusat di bawah supervisi Atasan Klinik Utama Semarang Medical Group.'
            }}
          </p>

          <!-- Quick Access Demo Accounts -->
          <div class="pt-4 border-t border-zinc-200/80 space-y-3">
            <div class="flex items-center justify-between text-[11px] font-bold uppercase tracking-wider opacity-60">
              <span>Akun Demo Cepat {{ brand === 'kidinest' ? 'Lantai 1' : 'Lantai 2' }} (1-Klik):</span>
              <span class="font-mono text-emerald-700">Password: klinik123</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <button
                v-for="acc in quickAccounts"
                :key="acc.empNum"
                @click="selectQuickAccount(acc)"
                class="p-3.5 rounded-2xl bg-white border text-left transition-all hover:bg-stone-50 active:scale-95 shadow-2xs group flex items-start gap-3"
                :class="brand === 'kidinest' ? 'border-[#E5A93C]/30 hover:border-[#0C3B36]' : 'border-[#9E606D]/30 hover:border-[#9E606D]'"
              >
                <div
                  class="p-2 rounded-xl transition-colors shrink-0"
                  :class="
                    brand === 'kidinest'
                      ? 'bg-[#E5A93C]/20 text-[#0C3B36] group-hover:bg-[#0C3B36] group-hover:text-white'
                      : 'bg-[#9E606D]/15 text-[#9E606D] group-hover:bg-[#9E606D] group-hover:text-white'
                  "
                >
                  <component :is="acc.icon" class="w-4 h-4" />
                </div>
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-extrabold truncate">{{ acc.name }}</div>
                  <div class="text-[10px] opacity-70 font-semibold">{{ acc.title }}</div>
                  <div class="flex items-center gap-1.5 mt-1">
                    <span class="font-mono text-[10px] font-bold bg-zinc-100 px-1.5 py-0.5 rounded text-zinc-700">
                      {{ acc.empNum }}
                    </span>
                    <span class="text-[9px] font-bold opacity-60">
                      • {{ acc.unit }}
                    </span>
                  </div>
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Right: Login Form (5 cols) -->
        <div
          class="md:col-span-5 bg-white p-7 sm:p-9 rounded-3xl border shadow-xs space-y-5"
          :class="brand === 'kidinest' ? 'border-[#E5A93C]/30' : 'border-[#9E606D]/25'"
        >
          <div>
            <h2 class="text-base font-extrabold">Form Kredensial Karyawan</h2>
            <p class="text-xs opacity-60 font-medium">Masukkan NIK Internal staf & kata sandi</p>
          </div>

          <!-- Error Alert -->
          <div
            v-if="errorMessage"
            class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2 animate-in fade-in duration-150"
          >
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <form @submit.prevent="handleFormSubmit" class="space-y-4">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-70 mb-1.5">
                Nomor Induk Karyawan (NIK Internal)
              </label>
              <input
                v-model="employeeNumber"
                type="text"
                required
                :placeholder="brand === 'kidinest' ? 'Misal: DOC-KID-01 atau OWNER-01' : 'Misal: DOC-SSC-01 atau OWNER-01'"
                class="w-full px-4 py-3 rounded-2xl border text-xs font-mono font-bold outline-none focus:bg-white transition-all"
                :class="
                  brand === 'kidinest'
                    ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]'
                    : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'
                "
              />
            </div>

            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider opacity-70 mb-1.5">
                Kata Sandi
              </label>
              <input
                v-model="password"
                type="password"
                required
                placeholder="Masukkan kata sandi..."
                class="w-full px-4 py-3 rounded-2xl border text-xs font-semibold outline-none focus:bg-white transition-all"
                :class="
                  brand === 'kidinest'
                    ? 'bg-[#FAF7F2] border-[#E5A93C]/30 text-[#0C3B36] focus:border-[#0C3B36]'
                    : 'bg-[#FAF5F5] border-[#9E606D]/30 text-[#2D1E24] focus:border-[#9E606D]'
                "
              />
            </div>

            <button
              type="submit"
              :disabled="isSubmitting"
              class="w-full mt-2 py-3 rounded-full font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
              :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
            >
              <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
              <KeyRound v-else class="w-3.5 h-3.5" />
              <span>Masuk ke Workstation {{ brand === 'kidinest' ? 'Lantai 1' : 'Lantai 2' }}</span>
            </button>
          </form>

          <div class="text-[11px] text-center opacity-60 font-medium pt-2">
            Lupa kata sandi? Hubungi Atasan / Owner klinik untuk reset kredensial NIK instan.
          </div>
        </div>

      </div>
    </main>

    <!-- Bottom Footer Note -->
    <div class="max-w-5xl mx-auto w-full text-center text-[11px] font-medium opacity-50 border-t border-zinc-200/60 pt-6">
      Satu Sistem Manajemen Klinik Terpadu &bull; Klinik Utama Semarang Medical Group &copy; 2026.
    </div>
  </div>
</template>
