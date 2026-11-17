<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  router.replace('/accounting/tax/spt-masa')
})
</script>

<template>
  <div class="flex items-center justify-center p-12 text-zinc-400 text-xs">
    <span>Memuat modul perpajakan klinik...</span>
  </div>
</template>
