<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
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
  Lock,
  Terminal,
  ChevronRight
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

// State Form (Default kosong untuk Production)
const employeeNumber = ref('')
const password = ref('')
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

// Mode Developer / Pitching Demo (Hidden via URL: ?demo=true atau ?dev=1 atau #demo)
const isDemoMode = ref(false)

function checkDemoMode() {
  const urlParams = new URLSearchParams(window.location.search)
  isDemoMode.value =
    urlParams.has('demo') ||
    urlParams.has('dev') ||
    window.location.hash.includes('demo')

  if (isDemoMode.value) {
    employeeNumber.value = props.brand === 'kidinest' ? 'DOC-KID-01' : 'DOC-SSC-01'
    password.value = 'klinik123'
  } else {
    employeeNumber.value = ''
    password.value = ''
  }
}

onMounted(() => {
  checkDemoMode()
})

watch(() => props.brand, () => {
  if (isDemoMode.value) {
    employeeNumber.value = props.brand === 'kidinest' ? 'DOC-KID-01' : 'DOC-SSC-01'
  }
})

function toggleDemoMode() {
  isDemoMode.value = !isDemoMode.value
  const url = new URL(window.location.href)
  if (!isDemoMode.value) {
    employeeNumber.value = ''
    password.value = ''
    url.searchParams.delete('demo')
    url.searchParams.delete('dev')
    window.history.replaceState({}, '', url.pathname)
  } else {
    employeeNumber.value = props.brand === 'kidinest' ? 'DOC-KID-01' : 'DOC-SSC-01'
    password.value = 'klinik123'
    url.searchParams.set('demo', 'true')
    window.history.replaceState({}, '', url.pathname + '?' + url.searchParams.toString())
  }
}

interface QuickAccount {
  empNum: string
  name: string
  title: string
  role: 'SUPERADMIN_OWNER' | 'DOCTOR' | 'GRO_CASHIER' | 'PHARMACIST'
  unit: string
  clinicRole: RoleType
  icon: any
}

// Quick demo accounts
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
    unit: 'Kidinest (Poli Anak)',
    clinicRole: 'doctor',
    icon: Stethoscope,
  },
  {
    empNum: 'DOC-KID-02',
    name: 'dr. Hendra Pratama, Sp.A',
    title: 'Konsultan Tumbang & Nutrisi • Poli 2',
    role: 'DOCTOR',
    unit: 'Kidinest (Poli Anak)',
    clinicRole: 'doctor',
    icon: Stethoscope,
  },
  {
    empNum: 'GRO-01',
    name: 'Siti Sarah, A.Md',
    title: 'Kasir & Meja Depan GRO',
    role: 'GRO_CASHIER',
    unit: 'Kidinest (Poli Anak)',
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
    unit: 'Semarang Skin Center',
    clinicRole: 'doctor',
    icon: Sparkles,
  },
  {
    empNum: 'DOC-SSC-02',
    name: 'dr. Kevin Wijaya, Sp.D.V.E',
    title: 'Spesialis Kulit & Bedah Minor • Ruang 2',
    role: 'DOCTOR',
    unit: 'Semarang Skin Center',
    clinicRole: 'doctor',
    icon: Building2,
  },
  {
    empNum: 'GRO-01',
    name: 'Amanda L., S.E',
    title: 'Kasir & Resepsionis VIP',
    role: 'GRO_CASHIER',
    unit: 'Semarang Skin Center',
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
      const allAccs = [...kidinestAccounts, ...sscAccounts]
      const matched = allAccs.find(a => a.empNum === empNum) ||
                      allAccs.find(a => a.clinicRole === targetClinicRole) ||
                      sscAccounts[0]

      auth.setAuth('demo-token-' + matched.empNum, {
        id: 1,
        employee_number: matched.empNum,
        name: matched.name,
        role: matched.role as any,
        clinic_unit: (matched.unit.includes('SSC') || matched.unit.includes('Skin') ? 'SSC' : matched.unit.includes('Kidinest') ? 'KIDINEST' : 'ALL') as any,
      })
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
    class="min-h-screen flex flex-col justify-between p-4 sm:p-8 font-['Plus_Jakarta_Sans',sans-serif] select-none transition-colors"
    :class="brand === 'kidinest' ? 'bg-[#FAF7F2] text-[#0C3B36]' : 'bg-[#FAF5F5] text-[#2D1E24]'"
  >
    <!-- Top Bar: Back Link & Demo Status Ribbon -->
    <div class="max-w-5xl mx-auto w-full space-y-3">
      <div class="flex items-center justify-between">
        <button
          @click="emit('navigate', brand)"
          class="inline-flex items-center gap-2 text-xs font-bold opacity-70 hover:opacity-100 transition-opacity"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Kembali ke Beranda {{ brand === 'kidinest' ? 'Kidinest' : 'Semarang Skin Center' }}</span>
        </button>

        <!-- Discreet Developer Mode Indicator -->
        <div class="flex items-center gap-2">
          <span
            v-if="isDemoMode"
            class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 text-amber-900 border border-amber-300/60 shadow-2xs"
          >
            <Sparkles class="w-3 h-3 text-amber-600" />
            <span>Developer / Pitch Mode Active</span>
          </span>
          <button
            @click="toggleDemoMode"
            class="text-[10px] font-mono opacity-40 hover:opacity-100 transition-opacity flex items-center gap-1"
            title="Klik untuk beralih mode demo / prod"
          >
            <Terminal class="w-3 h-3" />
            <span>{{ isDemoMode ? 'Exit ?demo' : 'Dev ?demo' }}</span>
          </button>
        </div>
      </div>

      <!-- Developer Demo Ribbon (Shown only in Demo Mode) -->
      <div
        v-if="isDemoMode"
        class="p-2.5 px-4 rounded-2xl bg-amber-500/10 border border-amber-500/25 text-amber-950 text-xs font-semibold flex items-center justify-between gap-3 animate-in fade-in slide-in-from-top-1 duration-150"
      >
        <div class="flex items-center gap-2">
          <Sparkles class="w-4 h-4 text-amber-600 shrink-0" />
          <span class="text-[11px] leading-tight">
            <strong>Mode Presentasi Klien:</strong> Menampilkan tombol 1-klik akun demo. Kunjungi URL tanpa <code>?demo</code> untuk tampilan produksi nyata.
          </span>
        </div>
        <button
          @click="toggleDemoMode"
          class="text-[10px] font-extrabold uppercase tracking-wider text-amber-900 hover:underline shrink-0"
        >
          Tutup Demo
        </button>
      </div>
    </div>

    <!-- MAIN BODY: ADAPTS BASED ON isDemoMode -->
    <main class="w-full flex-1 flex items-center justify-center py-6 sm:py-10">

      <!-- ======================================================== -->
      <!-- A. PRODUCTION VIEW (DEFAULT: CENTERED CLEAN CARD)       -->
      <!-- ======================================================== -->
      <div
        v-if="!isDemoMode"
        class="w-full max-w-md bg-white p-7 sm:p-9 rounded-3xl border shadow-xs space-y-6 animate-in fade-in zoom-in-95 duration-200"
        :class="brand === 'kidinest' ? 'border-[#E5A93C]/30' : 'border-[#9E606D]/25'"
      >
        <!-- Clinic Unit Badge & Header -->
        <div class="text-center space-y-3">
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
            <span>{{ brand === 'kidinest' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</span>
          </div>

          <h1 class="text-2xl font-extrabold tracking-tight">
            Portal Masuk Staf Medis
          </h1>
          <p class="text-xs opacity-75 max-w-xs mx-auto leading-relaxed font-medium">
            {{
              brand === 'kidinest'
                ? 'Masukkan NIK Internal dan kata sandi operasional Poli Anak & Vaksinasi.'
                : 'Masukkan NIK Internal dan kata sandi operasional Dermatologi & Estetika.'
            }}
          </p>
        </div>

        <!-- Error Alert -->
        <div
          v-if="errorMessage"
          class="p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-900 text-xs font-semibold flex items-center gap-2 animate-in fade-in duration-150"
        >
          <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
          <span>{{ errorMessage }}</span>
        </div>

        <!-- Production Form -->
        <form @submit.prevent="handleFormSubmit" class="space-y-4">
          <div>
            <label class="block text-[11px] font-bold uppercase tracking-wider opacity-70 mb-1.5">
              Nomor Induk Karyawan (NIK)
            </label>
            <input
              v-model="employeeNumber"
              type="text"
              required
              :placeholder="brand === 'kidinest' ? 'Contoh: DOC-KID-01 atau GRO-01' : 'Contoh: DOC-SSC-01 atau GRO-01'"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-mono font-bold uppercase tracking-wider outline-none focus:bg-white transition-all placeholder:normal-case placeholder:font-normal placeholder:tracking-normal placeholder:opacity-50"
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
              placeholder="••••••••••••"
              class="w-full px-4 py-3 rounded-2xl border text-xs font-semibold outline-none focus:bg-white transition-all placeholder:opacity-50"
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
            class="w-full mt-2 py-3.5 rounded-full font-bold text-white text-xs uppercase tracking-wider active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2 disabled:opacity-50"
            :class="brand === 'kidinest' ? 'bg-[#0C3B36] hover:bg-[#134E48]' : 'bg-[#9E606D] hover:bg-[#8B505D]'"
          >
            <Loader2 v-if="isSubmitting" class="w-3.5 h-3.5 animate-spin" />
            <KeyRound v-else class="w-3.5 h-3.5" />
            <span>Masuk ke Workstation</span>
          </button>
        </form>

        <div class="text-[11px] text-center opacity-60 font-medium pt-2 border-t border-zinc-100">
          Kendala akses NIK? Hubungi Atasan / Superadmin Klinik Utama Semarang Medical Group.
        </div>
      </div>

      <!-- ======================================================== -->
      <!-- B. DEVELOPER DEMO VIEW (EXPANDED DUAL PANEL / 1-CLICK)  -->
      <!-- ======================================================== -->
      <div v-else class="max-w-5xl mx-auto w-full animate-in fade-in duration-200">
        <div class="grid grid-cols-1 md:grid-cols-12 gap-8 items-start">

          <!-- Left: Quick-Fill Demo Accounts (7 cols) -->
          <div class="md:col-span-7 space-y-4">
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
              <span>{{ brand === 'kidinest' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</span>
            </div>

            <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight leading-tight">
              Portal Staf & Dokter (Demo Mode).
            </h1>
            <p class="text-xs sm:text-sm opacity-75 leading-relaxed font-medium">
              Mode presentasi aktif dengan akun siap pakai. Klik salah satu kartu staf di bawah ini untuk mengisi kredensial dan masuk ke workstation secara instan:
            </p>

            <!-- Quick Access Demo Accounts -->
            <div class="pt-3 border-t border-zinc-200/80 space-y-2.5">
              <div class="flex items-center justify-between text-[11px] font-bold uppercase tracking-wider opacity-60">
                <span>Pilih Akun Demo ({{ brand === 'kidinest' ? 'Kidinest' : 'SSC' }}):</span>
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
                  <ChevronRight class="w-3.5 h-3.5 opacity-30 group-hover:opacity-100 transition-opacity shrink-0 mt-1" />
                </button>
              </div>
            </div>
          </div>

          <!-- Right: Form Box with Demo Values (5 cols) -->
          <div
            class="md:col-span-5 bg-white p-7 sm:p-9 rounded-3xl border shadow-xs space-y-5"
            :class="brand === 'kidinest' ? 'border-[#E5A93C]/30' : 'border-[#9E606D]/25'"
          >
            <div>
              <h2 class="text-base font-extrabold">Form Kredensial</h2>
              <p class="text-xs opacity-60 font-medium">Terisi otomatis sesuai akun yang dipilih</p>
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
                  Nomor Induk Karyawan (NIK)
                </label>
                <input
                  v-model="employeeNumber"
                  type="text"
                  required
                  placeholder="Contoh: DOC-KID-01"
                  class="w-full px-4 py-3 rounded-2xl border text-xs font-mono font-bold outline-none focus:bg-white transition-all uppercase"
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
                <span>Masuk ke Workstation</span>
              </button>
            </form>

            <div class="text-[11px] text-center opacity-60 font-medium pt-2 border-t border-zinc-100">
              Demo Mode aktif. Buka tanpa <code>?demo</code> untuk melihat form polos produksi.
            </div>
          </div>

        </div>
      </div>

    </main>

    <!-- Bottom Footer Note -->
    <div class="max-w-5xl mx-auto w-full text-center text-[11px] font-medium opacity-50 border-t border-zinc-200/60 pt-4">
      Satu Sistem Manajemen Klinik Terpadu &bull; Klinik Utama Semarang Medical Group &copy; 2026.
    </div>
  </div>
</template>
