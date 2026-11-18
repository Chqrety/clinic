<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  Zap,
  Flame,
  RotateCcw,
  Trash2,
  X,
  CheckCircle2,
  AlertCircle,
  Loader2,
  Database,
  ShieldAlert,
  Sparkles
} from 'lucide-vue-next'
import { useClinicStore } from '../stores/clinic'

const emit = defineEmits<{
  (e: 'seeded'): void
}>()

const store = useClinicStore()

const isOpen = ref(false)
const isLoading = ref(false)
const activeAction = ref<string | null>(null)
const statusMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)

// Handle Keyboard Shortcut: Ctrl + Shift + G or Cmd + Shift + G
function handleKeyDown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && (e.key === 'G' || e.key === 'g')) {
    e.preventDefault()
    isOpen.value = !isOpen.value
  } else if (e.key === 'Escape' && isOpen.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)

  // Auto-open if query parameter ?godmode=1 exists
  const params = new URLSearchParams(window.location.search)
  if (params.get('godmode') === '1' || params.get('godmode') === 'true') {
    isOpen.value = true
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})

async function executeSeed(endpoint: string, actionKey: string) {
  isLoading.value = true
  activeAction.value = actionKey
  statusMessage.value = null

  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    const data = await res.json()
    if (!res.ok || !data.success) {
      throw new Error(data.error || 'Operasi seeder gagal dijalankan')
    }

    statusMessage.value = {
      type: 'success',
      text: data.message || 'Operasi berhasil!',
    }

    // Broadcast refresh event to active workstations
    window.dispatchEvent(new CustomEvent('clinic:queue-refreshed'))
    emit('seeded')

    // Auto-clear message after 4s
    setTimeout(() => {
      if (statusMessage.value?.type === 'success') {
        statusMessage.value = null
      }
    }, 4000)
  } catch (err: any) {
    statusMessage.value = {
      type: 'error',
      text: err.message || 'Terjadi kesalahan saat memanggil seeder backend.',
    }
  } finally {
    isLoading.value = false
    activeAction.value = null
  }
}
</script>

<template>
  <div>
    <!-- Floating Bottom-Right Trigger Pill (Placed next to Demo Switcher, completely clears sidebar logout on the left) -->
    <div id="god-mode-pill" class="fixed bottom-4 right-54 z-50 select-none print:hidden">
      <button
        @click="isOpen = true"
        title="Klik atau tekan Ctrl+Shift+G untuk membuka God Mode Seeder"
        class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-zinc-950/85 hover:bg-zinc-950 text-zinc-300 hover:text-white text-[11px] font-bold shadow-xl backdrop-blur-sm transition-all border border-zinc-700/60 active:scale-95 group"
      >
        <Zap class="w-3.5 h-3.5 text-amber-400 group-hover:scale-110 transition-transform" />
        <span>GOD MODE</span>
        <span class="text-[9px] text-zinc-400 font-mono bg-zinc-800/80 px-1.5 py-0.5 rounded border border-zinc-700/50">
          Ctrl+Shift+G
        </span>
      </button>
    </div>

    <!-- Backdrop & Modal Dialog -->
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-zinc-950/75 backdrop-blur-xs select-none animate-in fade-in duration-150"
      @click.self="isOpen = false"
    >
      <div
        class="w-full max-w-xl bg-zinc-900 border border-zinc-700/80 rounded-2xl shadow-2xl text-zinc-100 overflow-hidden animate-in zoom-in-95 duration-150 font-['Plus_Jakarta_Sans',sans-serif]"
      >
        <!-- Modal Header -->
        <div class="p-5 border-b border-zinc-800 bg-zinc-950/60 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-amber-400/10 border border-amber-400/20 text-amber-400 flex items-center justify-center">
              <Zap class="w-5 h-5" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-black text-sm text-white tracking-tight">God Mode Playground</h3>
                <span class="px-2 py-0.5 rounded-md text-[9px] font-extrabold uppercase bg-emerald-950/80 text-emerald-300 border border-emerald-800/60">
                  Dev &bull; Demo Active
                </span>
              </div>
              <p class="text-xs text-zinc-400 mt-0.5">
                Kontrol instan seeder transaksi dan pengujian skenario klinik.
              </p>
            </div>
          </div>

          <button
            @click="isOpen = false"
            class="p-1.5 rounded-lg text-zinc-400 hover:text-white hover:bg-zinc-800 transition-colors"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Alert Notification -->
        <div
          v-if="statusMessage"
          class="mx-5 mt-4 p-3 rounded-xl text-xs flex items-center gap-2.5 transition-all"
          :class="
            statusMessage.type === 'success'
              ? 'bg-emerald-950/60 border border-emerald-800 text-emerald-200'
              : 'bg-rose-950/60 border border-rose-800 text-rose-200'
          "
        >
          <CheckCircle2 v-if="statusMessage.type === 'success'" class="w-4 h-4 text-emerald-400 shrink-0" />
          <AlertCircle v-else class="w-4 h-4 text-rose-400 shrink-0" />
          <div class="flex-1 font-medium">{{ statusMessage.text }}</div>
          <button @click="statusMessage = null" class="text-zinc-400 hover:text-white">
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Action Cards Body -->
        <div class="p-5 space-y-3">
          <!-- Action 1: Rush Hour 50 Patients (Primary Highlight) -->
          <div class="p-4 rounded-xl bg-zinc-800/60 border border-zinc-700/80 hover:border-amber-400/40 transition-all space-y-3 group">
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-lg bg-amber-500/15 text-amber-400 flex items-center justify-center shrink-0">
                  <Flame class="w-4 h-4" />
                </div>
                <div>
                  <div class="text-xs font-bold text-white flex items-center gap-2">
                    <span>Simulate Rush Hour (50 Antrian)</span>
                    <span class="px-1.5 py-0.2 rounded text-[9px] font-black uppercase bg-amber-400/10 text-amber-300 border border-amber-400/30">
                      Rekomendasi Demo
                    </span>
                  </div>
                  <div class="text-[11px] text-zinc-400 mt-0.5">
                    30 Kidinest &bull; 20 SSC: 22 selesai bayar (billing omzet), 8 tunggu obat di apotek, 4 ruang dokter, 2 panggilan TV, 14 ruang tunggu.
                  </div>
                </div>
              </div>
            </div>

            <button
              @click="executeSeed('/api/dev/seeds/rush-hour', 'rush')"
              :disabled="isLoading"
              class="w-full py-2.5 px-3 rounded-lg bg-amber-400 hover:bg-amber-300 disabled:opacity-50 text-zinc-950 font-extrabold text-xs flex items-center justify-center gap-2 transition-all shadow-md active:scale-[0.99]"
            >
              <Loader2 v-if="isLoading && activeAction === 'rush'" class="w-3.5 h-3.5 animate-spin" />
              <Sparkles v-else class="w-3.5 h-3.5" />
              <span>Muat 50 Antrian Jam Sibuk Sekarang</span>
            </button>
          </div>

          <!-- Action 2: Reset Initial Demo (~8 Patients) -->
          <div class="p-4 rounded-xl bg-zinc-800/40 border border-zinc-700/60 hover:border-zinc-600 transition-all space-y-3">
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-lg bg-sky-500/15 text-sky-400 flex items-center justify-center shrink-0">
                  <RotateCcw class="w-4 h-4" />
                </div>
                <div>
                  <div class="text-xs font-bold text-white">
                    Reset Demo Awal (8 Antrian Bersih)
                  </div>
                  <div class="text-[11px] text-zinc-400 mt-0.5">
                    Membersihkan riwayat testing dan mengembalikan antrian awal yang tenang untuk walkthrough bertahap.
                  </div>
                </div>
              </div>
            </div>

            <button
              @click="executeSeed('/api/dev/seeds/reset', 'reset')"
              :disabled="isLoading"
              class="w-full py-2 px-3 rounded-lg bg-zinc-800 hover:bg-zinc-700 disabled:opacity-50 text-zinc-200 hover:text-white font-bold text-xs flex items-center justify-center gap-2 transition-all border border-zinc-700"
            >
              <Loader2 v-if="isLoading && activeAction === 'reset'" class="w-3.5 h-3.5 animate-spin" />
              <RotateCcw v-else class="w-3.5 h-3.5" />
              <span>Reset ke 8 Antrian Awal</span>
            </button>
          </div>

          <!-- Action 3: Wipe All Transactions (Zero State) -->
          <div class="p-4 rounded-xl bg-zinc-800/40 border border-zinc-700/60 hover:border-rose-900/40 transition-all space-y-3">
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-2.5">
                <div class="w-8 h-8 rounded-lg bg-rose-500/15 text-rose-400 flex items-center justify-center shrink-0">
                  <Trash2 class="w-4 h-4" />
                </div>
                <div>
                  <div class="text-xs font-bold text-white flex items-center gap-2">
                    <span>Kosongkan Seluruh Transaksi</span>
                    <span class="text-[9px] text-zinc-400">Jam 07:00 Pagi</span>
                  </div>
                  <div class="text-[11px] text-zinc-400 mt-0.5">
                    Nolkan seluruh antrian, resep, dan invoice. Master akun staf dan master obat tetap aman.
                  </div>
                </div>
              </div>
            </div>

            <button
              @click="executeSeed('/api/dev/seeds/wipe-transactions', 'wipe')"
              :disabled="isLoading"
              class="w-full py-2 px-3 rounded-lg bg-rose-950/40 hover:bg-rose-900/60 disabled:opacity-50 text-rose-300 font-bold text-xs flex items-center justify-center gap-2 transition-all border border-rose-800/60"
            >
              <Loader2 v-if="isLoading && activeAction === 'wipe'" class="w-3.5 h-3.5 animate-spin" />
              <Trash2 v-else class="w-3.5 h-3.5" />
              <span>Kosongkan Transaksi (Zero State)</span>
            </button>
          </div>
        </div>

        <!-- Modal Footer: Security Notice -->
        <div class="p-3.5 border-t border-zinc-800 bg-zinc-950/60 flex items-center justify-between text-[10px] text-zinc-400">
          <div class="flex items-center gap-1.5">
            <ShieldAlert class="w-3.5 h-3.5 text-zinc-400" />
            <span>Rute seeder hanya aktif pada server Dev/Demo (Otomatis 404 di Production).</span>
          </div>
          <span class="font-mono text-zinc-400">Esc untuk tutup</span>
        </div>
      </div>
    </div>
  </div>
</template>
