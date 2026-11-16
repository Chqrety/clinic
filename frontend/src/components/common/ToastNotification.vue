<script setup lang="ts">
import { useAccountingStore } from '@/stores/accounting'
import { CheckCircle2, AlertCircle, Info, X } from 'lucide-vue-next'

const store = useAccountingStore()
</script>

<template>
  <Transition
    enter-active-class="transform ease-out duration-300 transition"
    enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
    enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
    leave-active-class="transition ease-in duration-100"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="store.toastMessage"
      class="fixed bottom-6 right-6 z-50 max-w-md w-full bg-slate-900 text-white rounded-xl shadow-2xl p-4 border border-slate-700/80 flex items-start gap-3 backdrop-blur-md"
    >
      <div class="shrink-0 mt-0.5">
        <CheckCircle2
          v-if="store.toastMessage.type === 'success'"
          class="w-5 h-5 text-emerald-400"
        />
        <AlertCircle
          v-else-if="store.toastMessage.type === 'error'"
          class="w-5 h-5 text-rose-400"
        />
        <Info
          v-else
          class="w-5 h-5 text-sky-400"
        />
      </div>

      <div class="flex-1 text-sm leading-relaxed text-slate-100">
        {{ store.toastMessage.text }}
      </div>

      <button
        @click="store.toastMessage = null"
        class="text-slate-400 hover:text-white transition-colors"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  </Transition>
</template>
