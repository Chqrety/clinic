<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import ToastNotification from '@/components/common/ToastNotification.vue'

const isMobileSidebarOpen = ref(false)
const route = useRoute()

// Otomatis tutup drawer navigasi mobile saat rute berubah
watch(() => route.fullPath, () => {
  isMobileSidebarOpen.value = false
})
</script>

<template>
  <div class="h-screen w-screen flex overflow-hidden bg-slate-50 text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif] print:block print:h-auto print:w-full print:overflow-visible print:bg-white print:p-0">
    <!-- Persistent Left Editorial Zinc Sidebar (Responsive Drawer on Mobile) -->
    <AppSidebar
      :is-mobile-open="isMobileSidebarOpen"
      @close="isMobileSidebarOpen = false"
      class="print:hidden"
    />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 h-screen overflow-hidden print:block print:h-auto print:w-full print:overflow-visible">
      <!-- Persistent Top Clinical Navbar with Mobile Toggle -->
      <AppNavbar
        @toggle-sidebar="isMobileSidebarOpen = !isMobileSidebarOpen"
        class="print:hidden"
      />

      <!-- Active Route Page Content (Standard Clinical Ergonomics: 1600px Max-Width) -->
      <main class="flex-1 h-[calc(100vh-4rem)] overflow-y-auto overflow-x-hidden bg-slate-100/70 p-3 sm:p-5 lg:p-8 print:p-0 print:h-auto print:bg-white print:overflow-visible print:max-w-none">
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

