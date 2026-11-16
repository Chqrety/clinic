<script setup lang="ts">
import { computed } from 'vue'
import { formatIDR, compactIDR } from '@/utils/currency'

const props = withDefaults(
  defineProps<{
    title: string
    amount: number
    subtitle?: string
    trend?: 'up' | 'down' | 'neutral'
    trendValue?: string
    iconColor?: 'blue' | 'emerald' | 'amber' | 'indigo' | 'purple' | 'rose'
    compact?: boolean
  }>(),
  {
    iconColor: 'blue',
    compact: false
  }
)

const formattedAmount = computed(() => {
  return props.compact ? compactIDR(props.amount) : formatIDR(props.amount)
})

const colorClasses = computed(() => {
  switch (props.iconColor) {
    case 'emerald':
      return {
        bg: 'bg-emerald-50 text-emerald-600 border-emerald-100',
        ring: 'group-hover:border-emerald-300'
      }
    case 'amber':
      return {
        bg: 'bg-amber-50 text-amber-600 border-amber-100',
        ring: 'group-hover:border-amber-300'
      }
    case 'indigo':
      return {
        bg: 'bg-indigo-50 text-indigo-600 border-indigo-100',
        ring: 'group-hover:border-indigo-300'
      }
    case 'purple':
      return {
        bg: 'bg-purple-50 text-purple-600 border-purple-100',
        ring: 'group-hover:border-purple-300'
      }
    case 'rose':
      return {
        bg: 'bg-rose-50 text-rose-600 border-rose-100',
        ring: 'group-hover:border-rose-300'
      }
    case 'blue':
    default:
      return {
        bg: 'bg-blue-50 text-blue-600 border-blue-100',
        ring: 'group-hover:border-blue-300'
      }
  }
})
</script>

<template>
  <div
    class="group relative bg-white rounded-xl border border-slate-200/80 p-5 shadow-xs hover:shadow-md transition-all duration-200 flex flex-col justify-between"
  >
    <div class="flex items-start justify-between">
      <div class="space-y-1">
        <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">{{ title }}</p>
        <h3 class="text-2xl font-bold tracking-tight text-slate-900 font-mono">
          {{ formattedAmount }}
        </h3>
      </div>
      <div
        :class="[
          'w-10 h-10 rounded-lg flex items-center justify-center border shadow-xs',
          colorClasses.bg
        ]"
      >
        <slot name="icon" />
      </div>
    </div>

    <div class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs">
      <div class="flex items-center gap-1.5 font-medium">
        <span
          v-if="trendValue"
          :class="[
            'inline-flex items-center px-1.5 py-0.5 rounded text-[11px] font-semibold',
            trend === 'up'
              ? 'bg-emerald-50 text-emerald-700'
              : trend === 'down'
              ? 'bg-rose-50 text-rose-700'
              : 'bg-slate-100 text-slate-600'
          ]"
        >
          {{ trend === 'up' ? '▲' : trend === 'down' ? '▼' : '●' }} {{ trendValue }}
        </span>
        <span class="text-slate-500">{{ subtitle || 'vs bulan lalu' }}</span>
      </div>
      <slot name="action" />
    </div>
  </div>
</template>
