<script setup lang="ts">
import { ref } from 'vue'
import { useClinicStore, type RoleType } from '../stores/clinic'
import { Stethoscope, Pill, CreditCard, Package, Receipt, ShieldCheck, ArrowLeft, KeyRound } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'admin'): void
}>()

const store = useClinicStore()
const username = ref('dr.hendra')
const password = ref('••••••••')

function loginAs(role: RoleType) {
  store.setRole(role)
  emit('navigate', 'admin')
}

function handleFormSubmit() {
  // Demo auto-login based on username
  if (username.value.includes('apotek')) {
    loginAs('pharmacy')
  } else if (username.value.includes('kasir')) {
    loginAs('cashier')
  } else {
    loginAs('doctor')
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col justify-between p-6 sm:p-12 selection:bg-amber-300">
    <!-- Top Back Link -->
    <div class="max-w-4xl mx-auto w-full">
      <button
        @click="emit('navigate', 'landing')"
        class="inline-flex items-center gap-2 text-xs font-bold text-zinc-500 hover:text-zinc-950 transition-colors"
      >
        <ArrowLeft class="w-3.5 h-3.5" /> Kembali ke Beranda
      </button>
    </div>

    <!-- Center Login Layout (Editorial Brandlos Minimal) -->
    <main class="max-w-4xl mx-auto w-full py-12">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-12 items-start">

        <!-- Left: Heading & Context (6 cols) -->
        <div class="md:col-span-6 space-y-4">
          <span class="inline-flex items-center gap-1.5 px-3.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-zinc-950 text-white shadow-sm">
            Portal Khusus Staf
          </span>
          <h1 class="text-4xl sm:text-5xl font-extrabold tracking-tight text-zinc-950 leading-tight">
            Masuk ke Sistem Klinik.
          </h1>
          <p class="text-sm text-zinc-600 leading-relaxed">
            Akses dibatasi hanya untuk tenaga medis, kasir, farmasi, dan manajemen klinik yang berwenang.
          </p>

          <!-- Pitch Quick Role Login (Demo Super Weapon) -->
          <div class="pt-6 border-t border-black/10">
            <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 mb-3">
              Akses Cepat Demo Presentasi:
            </div>
            <div class="grid grid-cols-2 gap-2">
              <button
                @click="loginAs('doctor')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <Stethoscope class="w-4 h-4 text-blue-600" />
                <span>Ruang Dokter</span>
              </button>

              <button
                @click="loginAs('pharmacy')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <Pill class="w-4 h-4 text-emerald-600" />
                <span>Apotek & Farmasi</span>
              </button>

              <button
                @click="loginAs('cashier')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <CreditCard class="w-4 h-4 text-amber-600" />
                <span>Kasir & Antrian</span>
              </button>

              <button
                @click="loginAs('inventory')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <Package class="w-4 h-4 text-purple-600" />
                <span>Gudang Obat</span>
              </button>

              <button
                @click="loginAs('billing')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <Receipt class="w-4 h-4 text-rose-600" />
                <span>Audit & Invoice</span>
              </button>

              <button
                @click="loginAs('satusehat')"
                class="p-3 rounded-2xl bg-white border border-black/10 hover:border-black/25 text-left text-xs font-bold text-zinc-900 transition-all hover:bg-stone-50 active:scale-95 shadow-sm flex items-center gap-2"
              >
                <ShieldCheck class="w-4 h-4 text-blue-600" />
                <span>SATUSEHAT</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Right: Clean Minimalist Form (6 cols) -->
        <div class="md:col-span-6 bg-white p-8 sm:p-10 rounded-3xl border border-black/10 shadow-[0_4px_24px_rgba(0,0,0,0.04)]">
          <h2 class="text-lg font-bold text-zinc-950 mb-1">Login Kredensial</h2>
          <p class="text-xs text-zinc-500 mb-6">Gunakan akun SIP atau nomor registrasi pegawai.</p>

          <form @submit.prevent="handleFormSubmit" class="space-y-4">
            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">
                Username / NIP
              </label>
              <input
                v-model="username"
                type="text"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none focus:bg-white focus:border-zinc-950 transition-all"
              />
            </div>

            <div>
              <label class="block text-[11px] font-bold uppercase tracking-wider text-zinc-500 mb-1.5">
                Kata Sandi
              </label>
              <input
                v-model="password"
                type="password"
                required
                class="w-full px-4 py-3 rounded-2xl bg-[#F9F9F7] border border-black/10 text-xs font-bold text-zinc-900 outline-none focus:bg-white focus:border-zinc-950 transition-all"
              />
            </div>

            <button
              type="submit"
              class="w-full mt-2 py-3.5 rounded-full font-bold bg-amber-400 hover:bg-amber-500 text-zinc-950 text-xs active:scale-[0.98] transition-all shadow-sm flex items-center justify-center gap-2"
            >
              <KeyRound class="w-3.5 h-3.5" />
              <span>Masuk ke Ruang Kerja</span>
            </button>
          </form>
        </div>

      </div>
    </main>

    <!-- Bottom Footer Note -->
    <div class="max-w-4xl mx-auto w-full text-center text-xs text-zinc-400 border-t border-black/5 pt-6">
      Sistem Manajemen Klinik Sehat Modern &copy; 2026. Terkoneksi jaringan lokal on-premise & Kemenkes SATUSEHAT.
    </div>
  </div>
</template>
