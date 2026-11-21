<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: number | string | null | undefined
    placeholder?: string
    disabled?: boolean
    readonly?: boolean
    allowNegative?: boolean
    allowEmpty?: boolean
    inputClass?: string
    prefix?: string
  }>(),
  {
    modelValue: 0,
    placeholder: '0',
    disabled: false,
    readonly: false,
    allowNegative: false,
    allowEmpty: false,
    inputClass: '',
    prefix: ''
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: number): void
  (e: 'input', value: number): void
  (e: 'change', value: number): void
  (e: 'blur', event: FocusEvent): void
  (e: 'focus', event: FocusEvent): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const displayValue = ref('')

function formatToIdrDots(val: number | string | null | undefined): string {
  if (val === null || val === undefined) {
    return props.allowEmpty ? '' : '0'
  }
  const numericVal = typeof val === 'string' ? parseRaw(val) : val
  if (isNaN(numericVal)) {
    return props.allowEmpty ? '' : '0'
  }
  if (numericVal === 0 && props.allowEmpty) return ''
  const isNeg = numericVal < 0
  const abs = Math.abs(Math.round(numericVal))
  const formatted = abs.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  return isNeg && props.allowNegative ? `-${formatted}` : formatted
}

function parseRaw(str: string): number {
  if (!str) return 0
  const isNeg = props.allowNegative && str.trim().startsWith('-')
  const digitsOnly = str.replace(/\D/g, '')
  if (!digitsOnly) return 0
  const num = parseInt(digitsOnly, 10)
  return isNeg ? -num : num
}

// Watch parent modelValue changes (e.g. programmatic recalculations or initial load)
watch(
  () => props.modelValue,
  (newVal) => {
    const numericNewVal = typeof newVal === 'string' ? parseRaw(newVal) : (newVal || 0)
    const curParsed = parseRaw(displayValue.value)
    if (numericNewVal !== curParsed || displayValue.value === '') {
      displayValue.value = formatToIdrDots(newVal)
    }
  },
  { immediate: true }
)

function onInput(event: Event) {
  const el = event.target as HTMLInputElement
  const rawText = el.value
  const oldCursor = el.selectionStart || 0

  // Count how many numeric digits occurred before the cursor position
  const digitsBeforeCursor = rawText.slice(0, oldCursor).replace(/\D/g, '').length

  const numericVal = parseRaw(rawText)
  emit('update:modelValue', numericVal)
  emit('input', numericVal)

  // Re-format display value with IDR dots
  const formatted = rawText === '' && props.allowEmpty ? '' : formatToIdrDots(numericVal)
  displayValue.value = formatted
  el.value = formatted

  // Recalculate and restore cursor position seamlessly
  let newCursor = 0
  let countedDigits = 0
  for (let i = 0; i < formatted.length; i++) {
    if (/\d/.test(formatted[i])) {
      countedDigits++
    }
    if (countedDigits >= digitsBeforeCursor) {
      newCursor = i + 1
      break
    }
  }
  if (countedDigits === 0) newCursor = formatted.length

  requestAnimationFrame(() => {
    if (el && el === document.activeElement) {
      el.setSelectionRange(newCursor, newCursor)
    }
  })
}

function onChange(event: Event) {
  const numericVal = parseRaw((event.target as HTMLInputElement).value)
  emit('change', numericVal)
}

function onBlur(event: FocusEvent) {
  const numericVal = parseRaw((event.target as HTMLInputElement).value)
  displayValue.value = formatToIdrDots(numericVal)
  emit('blur', event)
}

function onFocus(event: FocusEvent) {
  const target = event.target as HTMLInputElement
  if (target) {
    target.select()
  }
  emit('focus', event)
}
</script>

<template>
  <div v-if="prefix" class="relative w-full">
    <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-xs pointer-events-none select-none">
      {{ prefix }}
    </span>
    <input
      ref="inputRef"
      type="text"
      inputmode="numeric"
      :value="displayValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :class="['pl-8', inputClass]"
      @input="onInput"
      @change="onChange"
      @blur="onBlur"
      @focus="onFocus"
    />
  </div>
  <input
    v-else
    ref="inputRef"
    type="text"
    inputmode="numeric"
    :value="displayValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :readonly="readonly"
    :class="inputClass"
    @input="onInput"
    @change="onChange"
    @blur="onBlur"
    @focus="onFocus"
  />
</template>
