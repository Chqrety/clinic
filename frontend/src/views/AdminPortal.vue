<script setup lang="ts">
import { computed } from 'vue'
import { useClinicStore } from '../stores/clinic'
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
