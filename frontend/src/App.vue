<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClinicStore } from './stores/clinic'
import LandingKidinest from './views/LandingKidinest.vue'
import LandingSSC from './views/LandingSSC.vue'
import PatientPortal from './views/PatientPortal.vue'
import LoginView from './views/LoginView.vue'
import AdminPortal from './views/AdminPortal.vue'
import QueueMonitor from './views/QueueMonitor.vue'
import { Baby, Sparkles, Tv, Shield, QrCode } from 'lucide-vue-next'

export type AppRoute = 'kidinest' | 'ssc' | 'kiosk' | 'login' | 'admin' | 'monitor' | 'landing'
const currentRoute = ref<AppRoute>('kidinest')
const showDemoNav = ref(false)

const store = useClinicStore()

function syncRouteFromPath() {
  const path = window.location.pathname.toLowerCase().replace(/^\/+/, '')
  if (path === 'ssc') {
    currentRoute.value = 'ssc'
  } else if (path === 'kiosk') {
    currentRoute.value = 'kiosk'
  } else if (path === 'login') {
    currentRoute.value = 'login'
  } else if (path === 'admin') {
    currentRoute.value = 'admin'
  } else if (path === 'monitor') {
    currentRoute.value = 'monitor'
  } else {
    // Default to kidinest
    currentRoute.value = 'kidinest'
  }
}

function navigateTo(route: AppRoute) {
  const targetRoute = route === 'landing' ? 'kidinest' : route
  currentRoute.value = targetRoute
  const targetPath = '/' + targetRoute
  if (window.location.pathname !== targetPath) {
    window.history.pushState({}, '', targetPath)
  }
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  syncRouteFromPath()
  window.addEventListener('popstate', syncRouteFromPath)
})

const activeComponent = computed(() => {
  switch (currentRoute.value) {
    case 'kidinest':
      return LandingKidinest
    case 'ssc':
      return LandingSSC
    case 'kiosk':
      return PatientPortal
    case 'login':
      return LoginView
    case 'admin':
      return AdminPortal
    case 'monitor':
      return QueueMonitor
    default:
      return LandingKidinest
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] selection:bg-amber-300">

    <!-- Active View with Snappy Transition -->
    <Transition name="page" mode="out-in">
      <component
        :is="activeComponent"
        :key="currentRoute"
        @navigate="navigateTo"
      />
    </Transition>

    <!-- Discreet Presenter Quick Jump Pill (Bottom-Right) -->
    <div class="fixed bottom-4 right-4 z-50">
      <div
        v-if="showDemoNav"
        class="mb-2 p-2 bg-white/95 backdrop-blur-md rounded-2xl border border-black/10 shadow-2xl space-y-1 text-xs font-bold text-zinc-800 animate-in fade-in slide-in-from-bottom-2 duration-150 min-w-[220px]"
      >
        <div class="px-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400">
          Dual Brand Landing
        </div>
        <button
          @click="navigateTo('kidinest'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Baby class="w-3.5 h-3.5 text-emerald-700" />
            <span>1. Kidinest (Poli Anak)</span>
          </div>
          <span v-if="currentRoute === 'kidinest'" class="text-emerald-600 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('ssc'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Sparkles class="w-3.5 h-3.5 text-rose-600" />
            <span>2. Semarang Skin Center</span>
          </div>
          <span v-if="currentRoute === 'ssc'" class="text-rose-600 font-bold">●</span>
        </button>

        <div class="px-2 pt-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400 border-t border-zinc-100">
          Workstation & Portal
        </div>

        <button
          @click="navigateTo('kiosk'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <QrCode class="w-3.5 h-3.5 text-zinc-600" />
            <span>3. Pasien (Kiosk HP)</span>
          </div>
          <span v-if="currentRoute === 'kiosk'" class="text-emerald-600">●</span>
        </button>

        <button
          @click="navigateTo('login'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-zinc-600" />
            <span>4. Login Staf</span>
          </div>
          <span v-if="currentRoute === 'login'" class="text-emerald-600">●</span>
        </button>

        <button
          @click="navigateTo('admin'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-amber-600" />
            <span>5. Admin Portal Staf</span>
          </div>
          <span v-if="currentRoute === 'admin'" class="text-emerald-600">●</span>
        </button>

        <button
          @click="navigateTo('monitor'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <Tv class="w-3.5 h-3.5 text-sky-600" />
            <span>6. Monitor TV Antrian</span>
          </div>
          <span v-if="currentRoute === 'monitor'" class="text-emerald-600">●</span>
        </button>
      </div>

      <button
        @click="showDemoNav = !showDemoNav"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-zinc-950/85 hover:bg-zinc-950 text-white text-xs font-extrabold shadow-xl backdrop-blur-sm transition-all active:scale-95 border border-white/20"
      >
        <Sparkles class="w-3.5 h-3.5 text-amber-400" />
        <span>Demo Switcher ({{ currentRoute }})</span>
      </button>
    </div>

  </div>
</template>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
