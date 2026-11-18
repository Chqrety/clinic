<script setup lang="ts">
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import ToastNotification from '@/components/common/ToastNotification.vue'
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-slate-50 text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] print:block print:h-auto print:w-full print:overflow-visible print:bg-white print:p-0">
    <!-- Persistent Left Editorial Zinc Sidebar (Suppressed during print) -->
    <AppSidebar class="print:hidden" />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 h-screen overflow-hidden print:block print:h-auto print:w-full print:overflow-visible">
      <!-- Persistent Top Clinical Navbar (Suppressed during print) -->
      <AppNavbar class="print:hidden" />

      <!-- Active Route Page Content (Standard Clinical Ergonomics: 1600px Max-Width) -->
      <main class="flex-1 h-[calc(100vh-4rem)] overflow-y-auto bg-slate-100/70 p-4 sm:p-6 lg:p-8 print:p-0 print:h-auto print:bg-white print:overflow-visible print:max-w-none">
        <div class="max-w-[1600px] mx-auto w-full print:max-w-none print:w-full print:p-0">
          <router-view v-slot="{ Component }">
            <Transition name="fade" mode="out-in">
              <component :is="Component" />
            </Transition>
          </router-view>
        </div>
      </main>
    </div>

    <!-- Floating Global Toast Notifications (Suppressed during print) -->
    <ToastNotification class="print:hidden" />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>

