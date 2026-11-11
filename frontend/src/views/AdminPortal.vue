<script setup lang="ts">
import { computed } from 'vue'
import { useClinicStore } from '../stores/clinic'
import AdminHeader from '../components/AdminHeader.vue'
import DoctorRoom from './DoctorRoom.vue'
import PharmacyRoom from './PharmacyRoom.vue'
import CashierView from './CashierView.vue'
import InventoryView from './InventoryView.vue'
import BillingView from './BillingView.vue'
import SatusehatView from './SatusehatView.vue'
import StaffAdmin from './StaffAdmin.vue'

const emit = defineEmits<{
  (e: 'navigate', view: 'landing' | 'kiosk' | 'login'): void
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
      return DoctorRoom
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-950 flex flex-col">
    <!-- Sleek Top Masthead for Authenticated Staff -->
    <AdminHeader @navigate="(v) => emit('navigate', v)" />

    <!-- Active Staff Workstation Content -->
    <main class="flex-1">
      <Transition name="fade" mode="out-in">
        <component :is="currentWorkstation" :key="store.currentRole" />
      </Transition>
    </main>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
