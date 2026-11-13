<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
  size?: 'sm' | 'md'
}>()

const badgeConfig = computed(() => {
  switch (props.status.toUpperCase()) {
    // Journal Statuses
    case 'POSTED':
      return {
        label: 'Posted',
        bg: 'bg-emerald-50 text-emerald-700 border-emerald-200 ring-emerald-500/10',
        dot: 'bg-emerald-500'
      }
    case 'APPROVED':
      return {
        label: 'Approved',
        bg: 'bg-teal-50 text-teal-700 border-teal-200 ring-teal-500/10',
        dot: 'bg-teal-500'
      }
    case 'REVIEWED':
      return {
        label: 'Reviewed',
        bg: 'bg-sky-50 text-sky-700 border-sky-200 ring-sky-500/10',
        dot: 'bg-sky-500'
      }
    case 'SUBMITTED':
      return {
        label: 'Submitted',
        bg: 'bg-indigo-50 text-indigo-700 border-indigo-200 ring-indigo-500/10',
        dot: 'bg-indigo-500'
      }
    case 'DRAFT':
      return {
        label: 'Draft',
        bg: 'bg-amber-50 text-amber-700 border-amber-200 ring-amber-500/10',
        dot: 'bg-amber-500'
      }
    case 'REJECTED':
      return {
        label: 'Rejected',
        bg: 'bg-rose-50 text-rose-700 border-rose-200 ring-rose-500/10',
        dot: 'bg-rose-500'
      }
    case 'REVERSED':
      return {
        label: 'Reversed',
        bg: 'bg-purple-50 text-purple-700 border-purple-200 ring-purple-500/10',
        dot: 'bg-purple-500'
      }

    // Mapping Statuses
    case 'MAPPED':
      return {
        label: 'Terpetakan',
        bg: 'bg-emerald-50 text-emerald-700 border-emerald-200',
        dot: 'bg-emerald-500'
      }
    case 'UNMAPPED':
      return {
        label: 'Belum Terpetakan',
        bg: 'bg-rose-50 text-rose-700 border-rose-200',
        dot: 'bg-rose-500'
      }

    // Period Statuses
    case 'OPEN':
      return {
        label: 'Buku Terbuka',
        bg: 'bg-emerald-50 text-emerald-700 border-emerald-200',
        dot: 'bg-emerald-500'
      }
    case 'LOCKED':
      return {
        label: 'Terkunci',
        bg: 'bg-slate-100 text-slate-700 border-slate-200',
        dot: 'bg-slate-500'
      }

    // Units
    case 'KIDINEST':
      return {
        label: 'Kidinest',
        bg: 'bg-blue-50 text-blue-700 border-blue-200',
        dot: 'bg-blue-500'
      }
    case 'SSC':
      return {
        label: 'SSC',
        bg: 'bg-rose-50 text-rose-700 border-rose-200',
        dot: 'bg-rose-500'
      }
    case 'ALL':
      return {
        label: 'Konsolidasi',
        bg: 'bg-slate-100 text-slate-700 border-slate-300',
        dot: 'bg-slate-600'
      }

    default:
      return {
        label: props.status,
        bg: 'bg-zinc-100 text-zinc-700 border-zinc-200',
        dot: 'bg-zinc-400'
      }
  }
})
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 font-medium rounded-full border ring-1 ring-inset transition-colors',
      badgeConfig.bg,
      size === 'sm' ? 'px-2 py-0.5 text-[11px]' : 'px-2.5 py-1 text-xs'
    ]"
  >
    <span :class="['w-1.5 h-1.5 rounded-full shrink-0', badgeConfig.dot]"></span>
    <span>{{ badgeConfig.label }}</span>
  </span>
</template>
