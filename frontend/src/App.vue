<script setup lang="ts">
import { computed } from 'vue'
import { useClinicStore } from './stores/clinic'
import RoleSwitcher from './components/RoleSwitcher.vue'
import PatientKiosk from './views/PatientKiosk.vue'
import CashierView from './views/CashierView.vue'
import DoctorRoom from './views/DoctorRoom.vue'
import PharmacyRoom from './views/PharmacyRoom.vue'
import InventoryView from './views/InventoryView.vue'
import BillingView from './views/BillingView.vue'
import SatusehatView from './views/SatusehatView.vue'

const store = useClinicStore()

const currentComponent = computed(() => {
  switch (store.currentRole) {
    case 'patient':
      return PatientKiosk
    case 'cashier':
      return CashierView
    case 'doctor':
      return DoctorRoom
    case 'pharmacy':
      return PharmacyRoom
    case 'inventory':
      return InventoryView
    case 'billing':
      return BillingView
    case 'satusehat':
      return SatusehatView
    default:
      return PatientKiosk
  }
})
</script>

<template>
  <div class="min-h-screen bg-[#F9F9F7] text-zinc-900 selection:bg-amber-300">
    <!-- Floating Role Switcher for seamless pitch/demo navigation -->
    <RoleSwitcher />

    <!-- Main View with smooth transition -->
    <main>
      <Transition name="view" mode="out-in">
        <component :is="currentComponent" :key="store.currentRole" />
      </Transition>
    </main>
  </div>
</template>

<style scoped>
.view-enter-active,
.view-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.view-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.view-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
