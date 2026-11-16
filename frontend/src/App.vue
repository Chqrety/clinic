<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import GodModeModal from './components/GodModeModal.vue'
import { Sparkles, Tv, Shield, QrCode } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const showDemoNav = ref(false)

function navigateTo(path: string) {
  showDemoNav.value = false
  router.push(path)
}
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300">
    <!-- Vue Router View: Handles both Clinic and Isolated Accounting -->
    <router-view />

    <!-- Discreet Presenter Quick Jump Pill (Bottom-Right, hidden on accounting) -->
    <div v-if="!route.path.startsWith('/accounting')" class="fixed bottom-4 right-4 z-50">
      <div
        v-if="showDemoNav"
        class="mb-2 p-2 bg-white/95 backdrop-blur-md rounded-2xl border border-black/10 shadow-2xl space-y-1 text-xs font-bold text-zinc-800 animate-in fade-in slide-in-from-bottom-2 duration-150 min-w-[240px]"
      >
        <div class="px-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400">
          Kidinest Pediatric Care
        </div>
        <button
          @click="navigateTo('/')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/' || route.path === '/kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <span>1. Beranda Kidinest</span>
        </button>

        <button
          @click="navigateTo('/kiosk-kidinest')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/kiosk-kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <QrCode class="w-3.5 h-3.5 text-emerald-700" />
            <span>2. Kiosk Reservasi</span>
          </div>
        </button>

        <button
          @click="navigateTo('/login-kidinest')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/login-kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-emerald-700" />
            <span>3. Login Staf Kidinest</span>
          </div>
        </button>

        <div class="px-2 pt-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400 border-t border-zinc-100">
          Semarang Skin Center
        </div>

        <button
          @click="navigateTo('/ssc')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Sparkles class="w-3.5 h-3.5 text-rose-600" />
            <span>4. Beranda SSC</span>
          </div>
        </button>

        <button
          @click="navigateTo('/kiosk-ssc')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/kiosk-ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <QrCode class="w-3.5 h-3.5 text-rose-600" />
            <span>5. Kiosk Reservasi SSC</span>
          </div>
        </button>

        <button
          @click="navigateTo('/login-ssc')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/login-ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-rose-600" />
            <span>6. Login Staf SSC</span>
          </div>
        </button>

        <div class="px-2 pt-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400 border-t border-zinc-100">
          Shared Medical Workstation
        </div>

        <button
          @click="navigateTo('/admin')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/admin' ? 'bg-zinc-100 text-zinc-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-zinc-900" />
            <span>7. Admin Workstation</span>
          </div>
        </button>

        <button
          @click="navigateTo('/monitor')"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="route.path === '/monitor' ? 'bg-zinc-100 text-zinc-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Tv class="w-3.5 h-3.5 text-zinc-900" />
            <span>8. TV Monitor Antrian</span>
          </div>
        </button>
      </div>

      <button
        @click="showDemoNav = !showDemoNav"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-zinc-950/85 hover:bg-zinc-950 text-white text-xs font-extrabold shadow-xl backdrop-blur-sm transition-all active:scale-95 border border-white/20"
      >
        <Sparkles class="w-3.5 h-3.5 text-amber-400" />
        <span>Demo Switcher</span>
      </button>
    </div>

    <!-- Dev & Demo God Mode Modal -->
    <GodModeModal />
  </div>
</template>
