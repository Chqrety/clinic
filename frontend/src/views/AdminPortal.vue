<script setup lang="ts">
import { computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useClinicStore } from '../stores/clinic'
import { useAuthStore } from '../stores/auth'
import AdminSidebar from '../components/AdminSidebar.vue'
import DoctorRoom from './DoctorRoom.vue'
import PharmacyRoom from './PharmacyRoom.vue'
import CashierView from './CashierView.vue'
import InventoryView from './InventoryView.vue'
import BillingView from './BillingView.vue'
import SatusehatView from './SatusehatView.vue'
import StaffAdmin from './StaffAdmin.vue'

const emit = defineEmits<{
  (e: 'navigate', view: 'kidinest' | 'ssc' | 'kiosk' | 'login' | 'monitor'): void
}>()

const store = useClinicStore()
const auth = useAuthStore()
const router = useRouter()

const validRoles = ['owner', 'doctor', 'pharmacy', 'cashier', 'inventory', 'billing', 'satusehat']

function getRoleFromHash(): string | null {
  const hash = window.location.hash.replace(/^#\/?/, '').toLowerCase()
  if (hash === 'management' || hash === 'management-report') {
    router.push('/management')
    return null
  }
  return validRoles.includes(hash) ? hash : null
}

function syncRoleFromHash() {
  const roleFromHash = getRoleFromHash()
  if (roleFromHash && store.currentRole !== roleFromHash) {
    store.setRole(roleFromHash as any)
  }
}

// Inisialisasi awal role secara langsung saat setup agar F5 tidak reset ke kasir
const initialRole = getRoleFromHash()
if (initialRole) {
  store.setRole(initialRole as any)
}

watch(() => store.currentRole, (newRole) => {
  if (getRoleFromHash() !== newRole) {
    window.location.hash = newRole
  }
})

onMounted(async () => {
  if (!auth.isAuthenticated) {
    const ok = await auth.fetchMe()
    if (!ok) {
      router.replace('/login-kidinest')
      return
    }
  }
  if (!window.location.hash) {
    window.location.hash = store.currentRole || 'cashier'
  }
  window.addEventListener('hashchange', syncRoleFromHash)
})

onUnmounted(() => {
  window.removeEventListener('hashchange', syncRoleFromHash)
})

const currentWorkstation = computed(() => {
  switch (store.currentRole) {
    case 'owner':
      return StaffAdmin
    case 'doctor':
      return DoctorRoom
    case 'pharmacy':
      return PharmacyRoom
    case 'cashier':
      return CashierView
    case 'inventory':
      return InventoryView
    case 'billing':
      return BillingView
    case 'satusehat':
      return SatusehatView
    default:
      return CashierView
  }
})
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-slate-50 text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- Fixed Minimalist Left Sidebar -->
    <AdminSidebar @navigate="(v) => emit('navigate', v)" />

    <!-- Main Workstation View Area -->
    <main class="flex-1 h-screen overflow-y-auto bg-slate-100/70 p-4 sm:p-6 lg:p-8">
      <div class="max-w-[1600px] mx-auto w-full">
        <Transition name="fade" mode="out-in">
          <component :is="currentWorkstation" :key="store.currentRole" />
        </Transition>
      </div>
    </main>

  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.12s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
