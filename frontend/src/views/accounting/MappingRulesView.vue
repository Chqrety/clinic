<script setup lang="ts">
import { computed } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountingStore } from '@/stores/accounting'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  Sliders,
  Sparkles,
  CheckCircle2,
  Building,
  Layers,
  ArrowRight
} from 'lucide-vue-next'

const transactionsStore = useTransactionsStore()
const accountingStore = useAccountingStore()

const rules = computed(() => transactionsStore.mappingRules)
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Rules Engine & Automation
          </span>
          <span class="text-xs text-slate-400">&bull; Event-Driven Accounting</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Aturan Pemetaan Akun (COA Mapping Rules)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Menghubungkan event operasional klinik (closing shift kasir, resep obat, faktur PO supplier, depresiasi) ke posisi Debit & Kredit secara otomatis.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <span class="px-3 py-1.5 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs font-bold flex items-center gap-1.5">
          <CheckCircle2 class="w-4 h-4 text-emerald-600" />
          <span>{{ rules.length }} Aturan Aktif</span>
        </span>
      </div>
    </div>

    <!-- Rules Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-slate-50/80 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase tracking-wider">
              <th class="px-4 py-3">Event Key & Modul</th>
              <th class="px-4 py-3">Nama Aturan Pemetaan</th>
              <th class="px-4 py-3 text-center">Unit</th>
              <th class="px-4 py-3">Posisi Akun Debit (Dr)</th>
              <th class="px-4 py-3">Posisi Akun Kredit (Cr)</th>
              <th class="px-4 py-3">Akun Terkait (Pajak/HPP)</th>
              <th class="px-4 py-3 text-center w-24">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <tr
              v-for="rule in rules"
              :key="rule.id"
              class="hover:bg-slate-50/80 transition-colors"
            >
              <td class="px-4 py-3.5">
                <span class="font-mono font-bold text-blue-700 block">{{ rule.eventKey }}</span>
                <span class="text-[10px] text-slate-500 font-semibold uppercase">{{ rule.sourceModule }}</span>
              </td>
              <td class="px-4 py-3.5">
                <p class="font-bold text-slate-900">{{ rule.eventName }}</p>
                <p class="text-[11px] text-slate-500">{{ rule.description }}</p>
              </td>
              <td class="px-4 py-3.5 text-center">
                <StatusBadge :status="rule.businessUnit" size="sm" />
              </td>
              <td class="px-4 py-3.5 font-medium text-emerald-800">
                <span class="font-mono font-bold">{{ rule.debitAccountCode }}</span>
                <span class="block text-[11px] text-slate-600">{{ rule.debitAccountName }}</span>
              </td>
              <td class="px-4 py-3.5 font-medium text-blue-800">
                <span class="font-mono font-bold">{{ rule.creditAccountCode }}</span>
                <span class="block text-[11px] text-slate-600">{{ rule.creditAccountName }}</span>
              </td>
              <td class="px-4 py-3.5 font-mono text-[11px] text-slate-600">
                <div v-if="rule.taxAccountCode">Pajak: {{ rule.taxAccountCode }}</div>
                <div v-if="rule.cogsAccountCode">HPP: {{ rule.cogsAccountCode }}</div>
                <div v-if="rule.inventoryAccountCode">Stok: {{ rule.inventoryAccountCode }}</div>
                <span v-if="!rule.taxAccountCode && !rule.cogsAccountCode" class="text-slate-400">-</span>
              </td>
              <td class="px-4 py-3.5 text-center">
                <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-700 border border-emerald-200">
                  Aktif
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

