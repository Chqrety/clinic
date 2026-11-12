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

export type AppRoute =
  | 'kidinest'
  | 'ssc'
  | 'kiosk'
  | 'kiosk-kidinest'
  | 'kiosk-ssc'
  | 'login'
  | 'login-kidinest'
  | 'login-ssc'
  | 'admin'
  | 'monitor'
  | 'landing'

const currentRoute = ref<AppRoute>('kidinest')
const showDemoNav = ref(false)

const store = useClinicStore()

// Current Brand context (kidinest = Lantai 1, ssc = Lantai 2)
const activeBrand = computed<'kidinest' | 'ssc'>(() => {
  if (
    currentRoute.value === 'ssc' ||
    currentRoute.value === 'kiosk-ssc' ||
    currentRoute.value === 'login-ssc'
  ) {
    return 'ssc'
  }
  return 'kidinest'
})

function syncRouteFromPath() {
  const path = window.location.pathname.toLowerCase().replace(/^\/+/, '')
  if (path === 'ssc') {
    currentRoute.value = 'ssc'
  } else if (path === 'kiosk-ssc' || path === 'kiosk/ssc') {
    currentRoute.value = 'kiosk-ssc'
  } else if (path === 'kiosk-kidinest' || path === 'kiosk/kidinest' || path === 'kiosk') {
    currentRoute.value = 'kiosk-kidinest'
  } else if (path === 'login-ssc' || path === 'login/ssc') {
    currentRoute.value = 'login-ssc'
  } else if (path === 'login-kidinest' || path === 'login/kidinest' || path === 'login') {
    currentRoute.value = 'login-kidinest'
  } else if (path === 'admin') {
    currentRoute.value = 'admin'
  } else if (path === 'monitor') {
    currentRoute.value = 'monitor'
  } else {
    currentRoute.value = 'kidinest'
  }
}

function navigateTo(route: string) {
  let targetRoute: AppRoute = 'kidinest'
  if (route === 'landing') {
    targetRoute = 'kidinest'
  } else {
    targetRoute = route as AppRoute
  }
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
    case 'kiosk-kidinest':
    case 'kiosk-ssc':
      return PatientPortal
    case 'login':
    case 'login-kidinest':
    case 'login-ssc':
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
        :brand="activeBrand"
        @navigate="navigateTo"
      />
    </Transition>

    <!-- Discreet Presenter Quick Jump Pill (Bottom-Right) -->
    <div class="fixed bottom-4 right-4 z-50">
      <div
        v-if="showDemoNav"
        class="mb-2 p-2 bg-white/95 backdrop-blur-md rounded-2xl border border-black/10 shadow-2xl space-y-1 text-xs font-bold text-zinc-800 animate-in fade-in slide-in-from-bottom-2 duration-150 min-w-[240px]"
      >
        <div class="px-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400">
          Lantai 1 • Kidinest
        </div>
        <button
          @click="navigateTo('kidinest'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Baby class="w-3.5 h-3.5 text-emerald-700" />
            <span>1. Landing Kidinest</span>
          </div>
          <span v-if="currentRoute === 'kidinest'" class="text-emerald-600 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('kiosk-kidinest'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'kiosk-kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <QrCode class="w-3.5 h-3.5 text-emerald-700" />
            <span>2. Kiosk Antrian (Lt. 1)</span>
          </div>
          <span v-if="currentRoute === 'kiosk-kidinest'" class="text-emerald-600 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('login-kidinest'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'login-kidinest' ? 'bg-emerald-50 text-emerald-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-emerald-700" />
            <span>3. Login Staf (Lt. 1)</span>
          </div>
          <span v-if="currentRoute === 'login-kidinest'" class="text-emerald-600 font-bold">●</span>
        </button>

        <div class="px-2 pt-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400 border-t border-zinc-100">
          Lantai 2 • Semarang Skin Center
        </div>

        <button
          @click="navigateTo('ssc'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Sparkles class="w-3.5 h-3.5 text-rose-600" />
            <span>4. Landing SSC</span>
          </div>
          <span v-if="currentRoute === 'ssc'" class="text-rose-600 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('kiosk-ssc'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'kiosk-ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <QrCode class="w-3.5 h-3.5 text-rose-600" />
            <span>5. Kiosk Antrian (Lt. 2)</span>
          </div>
          <span v-if="currentRoute === 'kiosk-ssc'" class="text-rose-600 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('login-ssc'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'login-ssc' ? 'bg-rose-50 text-rose-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-rose-600" />
            <span>6. Login Staf (Lt. 2)</span>
          </div>
          <span v-if="currentRoute === 'login-ssc'" class="text-rose-600 font-bold">●</span>
        </button>

        <div class="px-2 pt-2 py-1 text-[10px] font-black uppercase tracking-wider text-zinc-400 border-t border-zinc-100">
          Shared Medical Workstation
        </div>

        <button
          @click="navigateTo('admin'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'admin' ? 'bg-zinc-100 text-zinc-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Shield class="w-3.5 h-3.5 text-zinc-900" />
            <span>7. Admin Workstation</span>
          </div>
          <span v-if="currentRoute === 'admin'" class="text-zinc-900 font-bold">●</span>
        </button>

        <button
          @click="navigateTo('monitor'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
          :class="currentRoute === 'monitor' ? 'bg-zinc-100 text-zinc-950' : ''"
        >
          <div class="flex items-center gap-2">
            <Tv class="w-3.5 h-3.5 text-zinc-900" />
            <span>8. TV Monitor Antrian</span>
          </div>
          <span v-if="currentRoute === 'monitor'" class="text-zinc-900 font-bold">●</span>
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
