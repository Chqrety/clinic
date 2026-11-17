<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { ChevronDown, X } from 'lucide-vue-next'

export interface SelectOption {
  value: string | number
  label: string
  sublabel?: string
}

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    options: SelectOption[]
    placeholder?: string
    disabled?: boolean
    inputClass?: string
  }>(),
  {
    modelValue: '',
    placeholder: 'Pilih atau ketik untuk mencari...',
    disabled: false,
    inputClass: ''
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: string | number): void
  (e: 'change', val: string | number): void
}>()

const isOpen = ref(false)
const searchInput = ref('')
const containerRef = ref<HTMLElement | null>(null)

// Find currently selected option
const selectedOption = computed(() => {
  return props.options.find(opt => String(opt.value) === String(props.modelValue))
})

// Sync display text with modelValue
watch(
  () => props.modelValue,
  () => {
    if (selectedOption.value) {
      searchInput.value = selectedOption.value.label
    } else if (!props.modelValue) {
      searchInput.value = ''
    }
  },
  { immediate: true }
)

watch(
  () => props.options,
  () => {
    if (selectedOption.value) {
      searchInput.value = selectedOption.value.label
    }
  }
)

const filteredOptions = computed(() => {
  if (!searchInput.value.trim() || (selectedOption.value && searchInput.value === selectedOption.value.label)) {
    return props.options
  }
  const q = searchInput.value.toLowerCase()
  return props.options.filter(opt =>
    opt.label.toLowerCase().includes(q) ||
    String(opt.value).toLowerCase().includes(q) ||
    (opt.sublabel && opt.sublabel.toLowerCase().includes(q))
  )
})

function handleInput(e: Event) {
  const target = e.target as HTMLInputElement
  searchInput.value = target.value
  isOpen.value = true

  // If user cleared the input
  if (!target.value.trim()) {
    emit('update:modelValue', '')
    emit('change', '')
    return
  }

  // Check for exact value match
  const exact = props.options.find(
    opt => String(opt.value).toLowerCase() === target.value.trim().toLowerCase() ||
           opt.label.toLowerCase() === target.value.trim().toLowerCase()
  )
  if (exact) {
    emit('update:modelValue', exact.value)
    emit('change', exact.value)
  }
}

function selectOption(opt: SelectOption) {
  emit('update:modelValue', opt.value)
  emit('change', opt.value)
  searchInput.value = opt.label
  isOpen.value = false
}

function handleClear() {
  emit('update:modelValue', '')
  emit('change', '')
  searchInput.value = ''
  isOpen.value = false
}

function toggleOpen() {
  if (props.disabled) return
  isOpen.value = !isOpen.value
}

function handleClickOutside(event: MouseEvent) {
  if (containerRef.value && !containerRef.value.contains(event.target as Node)) {
    isOpen.value = false
    if (selectedOption.value) {
      searchInput.value = selectedOption.value.label
    } else if (!props.modelValue) {
      searchInput.value = ''
    }
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="containerRef" class="relative w-full select-none text-xs">
    <div class="relative flex items-center">
      <input
        type="text"
        :value="searchInput"
        :placeholder="placeholder"
        :disabled="disabled"
        @focus="isOpen = true"
        @input="handleInput"
        :class="[
          'w-full pl-3 pr-14 py-2 bg-white border border-zinc-200 rounded-xl text-xs font-medium text-zinc-900 focus:bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950 focus:border-zinc-950 transition-all',
          disabled ? 'bg-zinc-100 text-zinc-400 cursor-not-allowed' : 'cursor-text',
          inputClass
        ]"
      />
      <div class="absolute right-2 flex items-center gap-1">
        <button
          v-if="modelValue && !disabled"
          type="button"
          @click.stop="handleClear"
          class="p-0.5 text-zinc-400 hover:text-zinc-700 rounded transition-colors cursor-pointer"
          title="Hapus pilihan"
        >
          <X class="w-3.5 h-3.5" />
        </button>
        <button
          type="button"
          @click.stop="toggleOpen"
          :disabled="disabled"
          class="p-0.5 text-zinc-400 hover:text-zinc-700 rounded transition-colors cursor-pointer"
        >
          <ChevronDown :class="['w-3.5 h-3.5 transition-transform duration-150', isOpen ? 'rotate-180' : '']" />
        </button>
      </div>
    </div>

    <!-- Dropdown Options Menu -->
    <div
      v-if="isOpen && !disabled"
      class="absolute left-0 right-0 top-full mt-1 max-h-60 overflow-y-auto bg-white rounded-xl border border-zinc-200 shadow-xl z-50 divide-y divide-zinc-50"
    >
      <div
        v-for="opt in filteredOptions"
        :key="opt.value"
        @click="selectOption(opt)"
        :class="[
          'px-3 py-2 cursor-pointer transition-colors flex items-center justify-between text-xs',
          String(modelValue) === String(opt.value)
            ? 'bg-zinc-100 font-bold text-zinc-950'
            : 'hover:bg-zinc-50 text-zinc-700'
        ]"
      >
        <div class="truncate">
          <span class="block truncate">{{ opt.label }}</span>
          <span v-if="opt.sublabel" class="block text-[10px] text-zinc-400 truncate mt-0.5">{{ opt.sublabel }}</span>
        </div>
        <span
          v-if="String(modelValue) === String(opt.value)"
          class="w-1.5 h-1.5 rounded-full bg-zinc-950 shrink-0 ml-2"
        ></span>
      </div>

      <div v-if="filteredOptions.length === 0" class="p-3 text-center text-zinc-400 text-xs italic">
        Tidak ditemukan hasil yang cocok
      </div>
    </div>
  </div>
</template>
