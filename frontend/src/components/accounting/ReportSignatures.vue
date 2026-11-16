<script setup lang="ts">
interface SignaturePerson {
  label: string
  name: string
  title: string
}

const defaultSignatures: SignaturePerson[] = [
  {
    label: 'Dibuat Oleh,',
    name: 'Rian Pratama, S.Ak.',
    title: 'Staf Akuntansi & Pajak'
  },
  {
    label: 'Direview Oleh,',
    name: 'Dewi Lestari, SE, Ak.',
    title: 'Accounting Supervisor'
  },
  {
    label: 'Disetujui Oleh,',
    name: 'Budi Santoso, SE, Ak., CA',
    title: 'Finance & Accounting Manager'
  }
]
</script>

<template>
  <div class="mt-8 pt-6 border-t border-slate-200/80 bg-white/50 rounded-2xl p-6 print:border-t-2 print:border-black print:p-0 print:mt-6">
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-8 text-center">
      <div
        v-for="(person, idx) in defaultSignatures"
        :key="idx"
        class="flex flex-col items-center justify-between min-h-[140px]"
      >
        <span class="text-xs font-bold text-slate-700 tracking-wide mb-14 print:text-black">
          {{ person.label }}
        </span>

        <div class="w-full max-w-[220px]">
          <div class="border-b border-slate-400 pb-1 w-full text-center print:border-black">
            <span class="text-xs font-bold text-slate-900 tracking-wide print:text-black">
              {{ person.name }}
            </span>
          </div>
          <p class="text-[11px] text-slate-500 mt-1 font-medium print:text-gray-700">
            {{ person.title }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

