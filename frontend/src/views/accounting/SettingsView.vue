<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useLocaleStore, type Locale } from '@/stores/locale'
import { useAccountingStore } from '@/stores/accounting'
import {
  Globe,
  Languages,
  CheckCircle2,
  Calendar,
  Coins,
  FileSpreadsheet,
  BookOpen,
  Scale,
  TrendingUp,
  Receipt,
  ArrowRight,
  Sparkles,
  ShieldCheck,
  Save
} from 'lucide-vue-next'

const router = useRouter()
const localeStore = useLocaleStore()
const accountingStore = useAccountingStore()

const tempLocale = ref<Locale>(localeStore.currentLocale)

const sampleTerminology = [
  { id: 'coa', idTerm: 'Bagan Akun (COA)', enTerm: 'Chart of Accounts (COA)', icon: FileSpreadsheet },
  { id: 'ob', idTerm: 'Saldo Awal Akuntansi', enTerm: 'Accounting Opening Balance', icon: BookOpen },
  { id: 'gl', idTerm: 'Buku Besar (Ledger)', enTerm: 'General Ledger', icon: BookOpen },
  { id: 'tb', idTerm: 'Neraca Saldo (Trial Balance)', enTerm: 'Trial Balance', icon: Scale },
  { id: 'pl', idTerm: 'Laporan Laba Rugi', enTerm: 'Profit & Loss (Income Statement)', icon: TrendingUp },
  { id: 'bs', idTerm: 'Laporan Neraca', enTerm: 'Balance Sheet / Financial Position', icon: Scale },
  { id: 'cf', idTerm: 'Laporan Arus Kas', enTerm: 'Statement of Cash Flows', icon: Coins },
  { id: 'tax', idTerm: 'SPT Masa (Pajak Bulanan)', enTerm: 'Monthly Tax Returns (SPT Masa)', icon: Receipt }
]

function handleSelectLocale(loc: Locale) {
  tempLocale.value = loc
  localeStore.setLocale(loc)
  accountingStore.showToast(localeStore.t('settings.toastSuccess'), 'success')
}

function handleSave() {
  localeStore.setLocale(tempLocale.value)
  accountingStore.showToast(localeStore.t('settings.toastSuccess'), 'success')
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] max-w-5xl mx-auto pb-12">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
          <router-link to="/accounting/dashboard" class="hover:text-zinc-950 transition">
            {{ localeStore.t('nav.dashboardGroup') }}
          </router-link>
          <span>/</span>
          <span class="text-zinc-950">{{ localeStore.t('nav.settingsGroup') }}</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight flex items-center gap-2.5">
          <Languages class="w-6 h-6 text-blue-600" />
          <span>{{ localeStore.t('settings.title') }}</span>
        </h1>
        <p class="text-xs text-zinc-500 max-w-2xl">
          {{ localeStore.t('settings.subtitle') }}
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="handleSave"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>{{ localeStore.t('common.save') }}</span>
        </button>
      </div>
    </header>

    <!-- Section 1: Interactive Language Cards -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2 text-zinc-900 font-bold text-sm">
          <Globe class="w-4 h-4 text-blue-600" />
          <span>{{ localeStore.t('settings.langSectionTitle') }}</span>
        </div>
        <p class="text-xs text-zinc-500 mt-0.5">
          {{ localeStore.t('settings.langSectionDesc') }}
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
        <!-- Bahasa Indonesia Card -->
        <div
          @click="handleSelectLocale('id')"
          :class="[
            'p-5 rounded-2xl border-2 transition-all cursor-pointer relative flex flex-col justify-between space-y-4 group',
            localeStore.currentLocale === 'id'
              ? 'border-blue-600 bg-blue-50/20 shadow-md ring-2 ring-blue-500/10'
              : 'border-zinc-200 bg-zinc-50/50 hover:border-zinc-300 hover:bg-white'
          ]"
        >
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <span class="text-2xl">🇮🇩</span>
                <div>
                  <h3 class="font-bold text-sm text-zinc-950 flex items-center gap-2">
                    <span>{{ localeStore.t('settings.idName') }}</span>
                    <span
                      v-if="localeStore.currentLocale === 'id'"
                      class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-blue-100 text-blue-800 border border-blue-200"
                    >
                      Aktif
                    </span>
                  </h3>
                  <span class="text-[11px] font-semibold text-zinc-500">ID / Indonesia</span>
                </div>
              </div>

              <div
                :class="[
                  'w-5 h-5 rounded-full border flex items-center justify-center transition-colors',
                  localeStore.currentLocale === 'id'
                    ? 'border-blue-600 bg-blue-600 text-white'
                    : 'border-zinc-300 bg-white group-hover:border-zinc-400'
                ]"
              >
                <CheckCircle2 v-if="localeStore.currentLocale === 'id'" class="w-3.5 h-3.5" />
              </div>
            </div>

            <p class="text-xs text-zinc-600 leading-relaxed">
              {{ localeStore.t('settings.idDesc') }}
            </p>
          </div>

          <div class="pt-3 border-t border-zinc-200/80 grid grid-cols-2 gap-2 text-[11px]">
            <div class="bg-white/80 p-2 rounded-lg border border-zinc-200">
              <span class="text-zinc-400 block text-[9px] font-bold uppercase">Format Angka</span>
              <span class="font-mono font-bold text-zinc-800">Rp 15.000.000</span>
            </div>
            <div class="bg-white/80 p-2 rounded-lg border border-zinc-200">
              <span class="text-zinc-400 block text-[9px] font-bold uppercase">Format Tanggal</span>
              <span class="font-mono font-bold text-zinc-800">17/11/2026</span>
            </div>
          </div>
        </div>

        <!-- English Card -->
        <div
          @click="handleSelectLocale('en')"
          :class="[
            'p-5 rounded-2xl border-2 transition-all cursor-pointer relative flex flex-col justify-between space-y-4 group',
            localeStore.currentLocale === 'en'
              ? 'border-blue-600 bg-blue-50/20 shadow-md ring-2 ring-blue-500/10'
              : 'border-zinc-200 bg-zinc-50/50 hover:border-zinc-300 hover:bg-white'
          ]"
        >
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <span class="text-2xl">🇬🇧</span>
                <div>
                  <h3 class="font-bold text-sm text-zinc-950 flex items-center gap-2">
                    <span>{{ localeStore.t('settings.enName') }}</span>
                    <span
                      v-if="localeStore.currentLocale === 'en'"
                      class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-blue-100 text-blue-800 border border-blue-200"
                    >
                      Active
                    </span>
                  </h3>
                  <span class="text-[11px] font-semibold text-zinc-500">EN / International</span>
                </div>
              </div>

              <div
                :class="[
                  'w-5 h-5 rounded-full border flex items-center justify-center transition-colors',
                  localeStore.currentLocale === 'en'
                    ? 'border-blue-600 bg-blue-600 text-white'
                    : 'border-zinc-300 bg-white group-hover:border-zinc-400'
                ]"
              >
                <CheckCircle2 v-if="localeStore.currentLocale === 'en'" class="w-3.5 h-3.5" />
              </div>
            </div>

            <p class="text-xs text-zinc-600 leading-relaxed">
              {{ localeStore.t('settings.enDesc') }}
            </p>
          </div>

          <div class="pt-3 border-t border-zinc-200/80 grid grid-cols-2 gap-2 text-[11px]">
            <div class="bg-white/80 p-2 rounded-lg border border-zinc-200">
              <span class="text-zinc-400 block text-[9px] font-bold uppercase">Number Format</span>
              <span class="font-mono font-bold text-zinc-800">IDR 15,000,000</span>
            </div>
            <div class="bg-white/80 p-2 rounded-lg border border-zinc-200">
              <span class="text-zinc-400 block text-[9px] font-bold uppercase">Date Format</span>
              <span class="font-mono font-bold text-zinc-800">2026-11-17</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Section 2: Terminology Comparison Matrix -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-4">
      <div class="border-b border-zinc-100 pb-3">
        <div class="flex items-center gap-2 text-zinc-900 font-bold text-sm">
          <BookOpen class="w-4 h-4 text-emerald-600" />
          <span>{{ localeStore.t('settings.previewTitle') }}</span>
        </div>
        <p class="text-xs text-zinc-500 mt-0.5">
          {{ localeStore.t('settings.previewDesc') }}
        </p>
      </div>

      <div class="border border-zinc-200 rounded-xl overflow-hidden shadow-2xs">
        <table class="w-full text-left text-xs">
          <thead class="bg-zinc-50 border-b border-zinc-200 text-[10px] font-bold text-zinc-600 uppercase tracking-wider">
            <tr>
              <th class="py-2.5 px-4 w-12 text-center">#</th>
              <th class="py-2.5 px-4 w-1/2">
                <span class="inline-flex items-center gap-1.5">
                  <span>🇮🇩</span>
                  <span>{{ localeStore.t('settings.indonesianTerm') }}</span>
                </span>
              </th>
              <th class="py-2.5 px-4 w-1/2">
                <span class="inline-flex items-center gap-1.5">
                  <span>🇬🇧</span>
                  <span>{{ localeStore.t('settings.englishTerm') }}</span>
                </span>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr
              v-for="(term, idx) in sampleTerminology"
              :key="term.id"
              class="hover:bg-zinc-50/70 transition-colors"
            >
              <td class="py-2.5 px-4 text-center font-mono text-zinc-400 text-[11px]">{{ idx + 1 }}</td>
              <td class="py-2.5 px-4 font-semibold text-zinc-900 flex items-center gap-2">
                <component :is="term.icon" class="w-3.5 h-3.5 text-zinc-500" />
                <span>{{ term.idTerm }}</span>
              </td>
              <td class="py-2.5 px-4 font-semibold text-blue-700 font-mono">
                {{ term.enTerm }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Section 3: Quick Navigation to Test Reports -->
    <div class="bg-zinc-900 rounded-2xl p-6 text-white flex flex-col md:flex-row items-start md:items-center justify-between gap-4 shadow-xl">
      <div class="space-y-1">
        <span class="text-[10px] font-bold uppercase tracking-wider text-blue-400 flex items-center gap-1.5">
          <Sparkles class="w-3.5 h-3.5" />
          <span>Verifikasi Tampilan Laporan</span>
        </span>
        <h3 class="text-base font-bold text-white">
          Ingin melihat tampilan Buku Besar & Laporan Keuangan dalam bahasa yang dipilih?
        </h3>
        <p class="text-xs text-zinc-400">
          Seluruh navigasi sidebar, judul laporan, dan ekspor excel otomatis menyesuaikan preferensi bahasa aktif.
        </p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <router-link
          to="/accounting/ledger"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-xl text-xs font-bold transition shadow-xs flex items-center gap-1.5 cursor-pointer"
        >
          <span>Buka Buku Besar</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </router-link>
        <router-link
          to="/accounting/reports/profit-loss"
          class="px-4 py-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 hover:text-white rounded-xl text-xs font-bold transition border border-zinc-700"
        >
          <span>Buka Laba Rugi</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

