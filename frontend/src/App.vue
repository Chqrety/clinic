<script setup lang="ts">
import { ref, computed } from 'vue'
import { useClinicStore } from './stores/clinic'
import LandingPage from './views/LandingPage.vue'
import PatientPortal from './views/PatientPortal.vue'
import LoginView from './views/LoginView.vue'
import AdminPortal from './views/AdminPortal.vue'
import QueueMonitor from './views/QueueMonitor.vue'
import { Compass, Sparkles, Tv } from 'lucide-vue-next'

type AppRoute = 'landing' | 'kiosk' | 'login' | 'admin' | 'monitor'
const currentRoute = ref<AppRoute>('landing')
const showDemoNav = ref(false)

const store = useClinicStore()

function navigateTo(route: AppRoute) {
  currentRoute.value = route
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const activeComponent = computed(() => {
  switch (currentRoute.value) {
    case 'landing':
      return LandingPage
    case 'kiosk':
      return PatientPortal
    case 'login':
      return LoginView
    case 'admin':
      return AdminPortal
    case 'monitor':
      return QueueMonitor
    default:
      return LandingPage
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
      <div v-if="showDemoNav" class="mb-2 p-2 bg-white/95 backdrop-blur-md rounded-2xl border border-black/10 shadow-xl space-y-1 text-xs font-bold text-zinc-800 animate-in fade-in slide-in-from-bottom-2 duration-150">
        <button
          @click="navigateTo('landing'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <span>1. Landing Page (Brandlos)</span>
          <span v-if="currentRoute === 'landing'" class="text-emerald-600">●</span>
        </button>
        <button
          @click="navigateTo('kiosk'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <span>2. Pasien (Kiosk HP)</span>
          <span v-if="currentRoute === 'kiosk'" class="text-emerald-600">●</span>
        </button>
        <button
          @click="navigateTo('login'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <span>3. Login Staf</span>
          <span v-if="currentRoute === 'login'" class="text-emerald-600">●</span>
        </button>
        <button
          @click="navigateTo('admin'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <span>4. Admin Portal Staf</span>
          <span v-if="currentRoute === 'admin'" class="text-emerald-600">●</span>
        </button>
        <button
          @click="navigateTo('monitor'); showDemoNav = false"
          class="w-full text-left px-3 py-1.5 rounded-xl hover:bg-stone-100 flex items-center justify-between"
        >
          <span>5. Monitor TV Antrian</span>
          <span v-if="currentRoute === 'monitor'" class="text-emerald-600">●</span>
        </button>
      </div>

      <button
        @click="showDemoNav = !showDemoNav"
        class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full bg-zinc-950/80 hover:bg-zinc-950 text-white text-[11px] font-bold shadow-lg backdrop-blur-sm transition-all active:scale-95"
      >
        <Sparkles class="w-3 h-3 text-amber-400" />
        <span>Demo Switcher</span>
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
