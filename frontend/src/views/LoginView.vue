<script setup lang="ts">
import { ref } from 'vue'
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
  CheckCircle2
} from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'admin'): void
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

const quickAccounts: QuickAccount[] = [
  {
    empNum: 'OWNER-01',
    name: 'dr. Ratna Paramita, M.Kes',
    title: 'Owner & Atasan Klinik',
    role: 'SUPERADMIN_OWNER',
    unit: 'Semua Unit',
    clinicRole: 'owner',
    icon: ShieldCheck,
  },
  {
    empNum: 'DOC-KID-01',
    name: 'dr. Hendra Pratama, Sp.A',
    title: 'Dokter Spesialis Anak',
    role: 'DOCTOR',
    unit: 'Kidinest',
    clinicRole: 'doctor',
    icon: Stethoscope,
  },
  {
    empNum: 'DOC-SSC-01',
    name: 'dr. Nadia Kusuma, Sp.KK',
    title: 'Dokter Spesialis Kulit',
    role: 'DOCTOR',
    unit: 'Semarang Skin Center',
    clinicRole: 'doctor',
    icon: Building2,
  },
  {
    empNum: 'GRO-01',
    name: 'Siti Sarah, A.Md',
    title: 'Meja Depan & Kasir GRO',
    role: 'GRO_CASHIER',
    unit: 'Semua Unit',
    clinicRole: 'cashier',
    icon: CreditCard,
  },
  {
    empNum: 'FARMASI-01',
    name: 'Apt. Dimas Wicaksono, S.Farm',
    title: 'Apoteker Penanggung Jawab',
    role: 'PHARMACIST',
    unit: 'Semua Unit',
    clinicRole: 'pharmacy',
    icon: Pill,
  },
]

async function submitLogin(empNum: string, pass: string, targetClinicRole?: RoleType) {
  isSubmitting.value = true
  errorMessage.value = null

  // 1. Coba login ke API backend
  const success = await auth.login(empNum, pass)
  if (success && auth.user) {
    // Mapping role backend ke workstation frontend
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
    // 2. Jika backend offline saat demo, fallback gracefully
    if (auth.error) {
      errorMessage.value = auth.error
    }
    // Jika tombol demo diklik dan backend offline, izinkan masuk mode demo
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
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col justify-between p-6 sm:p-12 selection:bg-amber-300">
    <!-- Top Back Link -->
    <div class="max-w-5xl mx-auto w-full">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-2 text-xs font-bold text-zinc-500 hover:text-zinc-950 transition-colors"
      >
        <ArrowLeft class="w-3.5 h-3.5" /> Kembali ke Beranda
      </button>
    </div>

    <!-- Center Login Layout -->
    <main class="max-w-5xl mx-auto w-full py-8">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-10 items-start">

        <!-- Left: Context & Quick-Fill Demo (7 cols) -->
        <div class="md:col-span-7 space-y-5">
          <div class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-xs">
            Smart Clinic Management System
          </div>
          <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight text-zinc-950 leading-tight">
            Masuk dengan NIK Karyawan.
          </h1>
          <p class="text-xs sm:text-sm text-zinc-600 leading-relaxed">
            Sistem autentikasi terpadu Kidinest dan Semarang Skin Center. Akun dan wewenang dikelola langsung oleh Atasan/Owner klinik.
          </p>

          <!-- Quick Access Pill Selector (Demo Ready) -->
          <div class="pt-4 border-t border-zinc-200/80 space-y-3">
            <div class="flex items-center justify-between text-[11px] font-bold uppercase tracking-wider text-zinc-500">
              <span>Akun Demo Cepat (1-Klik):</span>
              <span class="text-emerald-700 font-mono">Password: klinik123</span>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
              <button
                v-for="acc in quickAccounts"
                :key="acc.empNum"
                @click="selectQuickAccount(acc)"
                class="p-3.5 rounded-2xl bg-white border border-zinc-200/80 hover:border-zinc-950 text-left transition-all hover:bg-stone-50 active:scale-95 shadow-2xs group flex items-start gap-3"
              >
                <div class="p-2 rounded-xl bg-stone-100 group-hover:bg-zinc-950 group-hover:text-white transition-colors shrink-0">
                  <component :is="acc.icon" class="w-4 h-4" />
                </div>
                <div class="min-w-0 flex-1">
                  <div class="text-xs font-extrabold text-zinc-950 truncate">{{ acc.name }}</div>
                  <div class="text-[10px] text-zinc-500 font-semibold">{{ acc.title }}</div>
                  <div class="flex items-center gap-1.5 mt-1">
                    <span class="font-mono text-[10px] font-bold text-zinc-700 bg-stone-100 px-1.5 py-0.5 rounded">
                      {{ acc.empNum }}
                    </span>
                    <span class="text-[9px] font-bold text-zinc-500">
                      • {{ acc.unit }}
                    </span>
                  </div>
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Right: Login Form (5 cols) -->
        <div class="md:col-span-5 bg-white p-7 sm:p-9 rounded-3xl border border-zinc-200/80 shadow-xs space-y-5">
          <div>
            <h2 class="text-base font-extrabold text-zinc-950">Form Kredensial Staf</h2>
            <p class="text-xs text-zinc-500 font-medium">Masukkan Nomor Induk Karyawan & kata sandi</p>
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
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                Nomor Induk Karyawan (NIK Internal)
              </label>
              <input
                v-model="employeeNumber"
                type="text"
                required
                placeholder="Misal: OWNER-01 atau DOC-KID-01"
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200 text-xs font-mono font-bold text-zinc-950 outline-none focus:bg-white focus:border-zinc-950 transition-all"
              />
            </div>

            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-600 mb-1.5">
                Kata Sandi
              </label>
              <input
                v-model="password"
                type="password"
                required
                placeholder="Masukkan kata sandi..."
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-zinc-200 text-xs font-semibold text-zinc-950 outline-none focus:bg-white focus:border-zinc-950 transition-all"
              />
            </div>

            <button
              type="submit"
              :disabled="isSubmitting"
              class="w-full mt-2 py-3 rounded-full font-bold bg-zinc-950 hover:bg-zinc-800 text-white text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
              <KeyRound v-else class="w-3.5 h-3.5" />
              <span>Masuk ke Ruang Kerja</span>
            </button>
          </form>

          <div class="text-[11px] text-center text-zinc-500 font-medium pt-2">
            Lupa kata sandi? Hubungi Atasan / Owner klinik untuk melakukan reset password instan.
          </div>
        </div>

      </div>
    </main>

    <!-- Bottom Footer Note -->
    <div class="max-w-5xl mx-auto w-full text-center text-[11px] font-medium text-zinc-400 border-t border-zinc-200/60 pt-6">
      Sistem Manajemen Klinik Sehat Modern &copy; 2026. Kidinest & Semarang Skin Center.
    </div>
  </div>
</template>
