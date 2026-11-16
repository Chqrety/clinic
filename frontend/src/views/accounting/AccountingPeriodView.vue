<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLedgerStore } from '@/stores/ledger'
import { useJournalsStore } from '@/stores/journals'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountingStore } from '@/stores/accounting'
import StatusBadge from '@/components/common/StatusBadge.vue'
import {
  CalendarCheck,
  Lock,
  Unlock,
  CheckCircle2,
  AlertCircle,
  FileCheck2,
  ShieldCheck,
  RotateCcw
} from 'lucide-vue-next'

const ledgerStore = useLedgerStore()
const journalsStore = useJournalsStore()
const transactionsStore = useTransactionsStore()
const accountingStore = useAccountingStore()

const periods = computed(() => ledgerStore.periods)
const activePeriod = computed(() => periods.value.find(p => p.month === accountingStore.fiscalMonth))

// Checklist Verification States
const unprocessedCount = computed(() =>
  transactionsStore.transactions.filter(t => t.journalStatus === 'UNPROCESSED').length
)

const unpostedCount = computed(() =>
  journalsStore.journals.filter(j => j.status !== 'POSTED' && j.status !== 'REVERSED').length
)

const hasDepreciation = computed(() =>
  journalsStore.journals.some(j => j.journalSource === 'FIXED_ASSET' && j.status === 'POSTED')
)

const isTrialBalanceBalanced = computed(() => {
  return ledgerStore.trialBalanceTotals.endingDebit === ledgerStore.trialBalanceTotals.endingCredit
})

const canExecuteClose = computed(() => {
  return (
    unprocessedCount.value === 0 &&
    unpostedCount.value === 0 &&
    hasDepreciation.value &&
    isTrialBalanceBalanced.value
  )
})

function handleLockPeriod() {
  if (activePeriod.value) {
    ledgerStore.lockPeriod(activePeriod.value.id)
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Period Closing & Governance
          </span>
          <span class="text-xs text-slate-400">&bull; Anti-Tamper Lock</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Tutup Buku Periode Akuntansi
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Proses penutupan buku bulanan mencegah pengeditan atau posting jurnal susulan pada periode yang telah diaudit dan dikunci.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          v-if="activePeriod && activePeriod.status === 'OPEN'"
          @click="handleLockPeriod"
          :disabled="!accountingStore.canLockPeriod || !canExecuteClose"
          class="flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-xl text-xs font-bold shadow-xs transition-colors"
        >
          <Lock class="w-4 h-4" />
          <span>Kunci Periode & Tutup Buku</span>
        </button>
      </div>
    </div>

    <!-- Active Period Closing Status Banner -->
    <div
      v-if="activePeriod"
      :class="[
        'p-5 rounded-2xl border flex flex-col md:flex-row md:items-center justify-between gap-4',
        activePeriod.status === 'LOCKED'
          ? 'bg-slate-900 text-white border-slate-800'
          : 'bg-blue-50 border-blue-200 text-blue-950'
      ]"
    >
      <div class="space-y-1">
        <div class="flex items-center gap-2">
          <Lock v-if="activePeriod.status === 'LOCKED'" class="w-5 h-5 text-amber-400" />
          <Unlock v-else class="w-5 h-5 text-blue-600" />
          <h2 class="text-base font-bold">
            Status Periode: {{ activePeriod.periodName }}
          </h2>
          <StatusBadge :status="activePeriod.status" size="sm" />
        </div>
        <p :class="['text-xs', activePeriod.status === 'LOCKED' ? 'text-slate-400' : 'text-blue-700']">
          {{ activePeriod.status === 'LOCKED'
            ? `Periode ini telah dikunci pada ${activePeriod.lockedAt} oleh ${activePeriod.lockedBy}. Jurnal tidak dapat diedit.`
            : 'Periode masih aktif dibuka. Lakukan checklist di bawah ini sebelum melaksanakan tutup buku.' }}
        </p>
      </div>

      <div class="flex items-center gap-3 text-xs">
        <span class="px-3 py-1.5 rounded-lg bg-white/10 font-mono font-bold">
          Rentang: {{ activePeriod.startDate }} s/d {{ activePeriod.endDate }}
        </span>
      </div>
    </div>

    <!-- 5-Point Closing Checklist -->
    <div class="bg-white rounded-2xl border border-slate-200/80 p-6 shadow-2xs space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div class="flex items-center gap-2">
          <ShieldCheck class="w-4 h-4 text-blue-600" />
          <h3 class="text-sm font-bold text-slate-900">Checklist Prasyarat Tutup Buku (Month-End Checklist)</h3>
        </div>
        <span class="text-xs font-semibold text-slate-500">5 Poin Verifikasi</span>
      </div>

      <div class="space-y-3 text-xs">
        <!-- Item 1 -->
        <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <CheckCircle2 v-if="unprocessedCount === 0" class="w-5 h-5 text-emerald-600 shrink-0" />
            <AlertCircle v-else class="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p class="font-bold text-slate-800">1. Semua Transaksi Operasional Telah Dibuatkan Draft Jurnal</p>
              <p class="text-[11px] text-slate-500">
                Memastikan tidak ada closing kasir, resep farmasi, atau faktur pembelian yang belum tercatat.
              </p>
            </div>
          </div>
          <span :class="['font-bold font-mono px-2.5 py-1 rounded-md text-[11px]', unprocessedCount === 0 ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700']">
            {{ unprocessedCount === 0 ? 'Lolos (0 Belum)' : `${unprocessedCount} Transaksi Menunggu` }}
          </span>
        </div>

        <!-- Item 2 -->
        <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <CheckCircle2 v-if="unpostedCount === 0" class="w-5 h-5 text-emerald-600 shrink-0" />
            <AlertCircle v-else class="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p class="font-bold text-slate-800">2. Seluruh Draft Jurnal Telah Diapprove & Diposting ke Buku Besar</p>
              <p class="text-[11px] text-slate-500">
                Tidak ada jurnal berstatus DRAFT, SUBMITTED, atau REVIEWED yang menggantung.
              </p>
            </div>
          </div>
          <span :class="['font-bold font-mono px-2.5 py-1 rounded-md text-[11px]', unpostedCount === 0 ? 'Lolos (Semua Posted)' : 'bg-amber-50 text-amber-700']">
            {{ unpostedCount === 0 ? 'Lolos (0 Unposted)' : `${unpostedCount} Jurnal Belum Diposting` }}
          </span>
        </div>

        <!-- Item 3 -->
        <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <CheckCircle2 v-if="hasDepreciation" class="w-5 h-5 text-emerald-600 shrink-0" />
            <AlertCircle v-else class="w-5 h-5 text-amber-500 shrink-0" />
            <div>
              <p class="font-bold text-slate-800">3. Jurnal Penyusutan Aset Tetap Medis Telah Diposting</p>
              <p class="text-[11px] text-slate-500">
                Depresiasi alat USG 4D Kidinest & Pico Laser SSC telah diamortisasi ke beban bulan ini.
              </p>
            </div>
          </div>
          <span :class="['font-bold font-mono px-2.5 py-1 rounded-md text-[11px]', hasDepreciation ? 'bg-emerald-50 text-emerald-700' : 'bg-amber-50 text-amber-700']">
            {{ hasDepreciation ? 'Sudah Diposting' : 'Belum Diposting' }}
          </span>
        </div>

        <!-- Item 4 -->
        <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <CheckCircle2 class="w-5 h-5 text-emerald-600 shrink-0" />
            <div>
              <p class="font-bold text-slate-800">4. Rekonsiliasi Kas Kasir & Rekening Bank Kliring Selesai</p>
              <p class="text-[11px] text-slate-500">
                Saldo kasir fisik dan mutasi rekening BCA / Mandiri kliring telah dicocokkan.
              </p>
            </div>
          </div>
          <span class="font-bold font-mono px-2.5 py-1 rounded-md text-[11px] bg-emerald-50 text-emerald-700">
            Rekonsiliasi Lengkap
          </span>
        </div>

        <!-- Item 5 -->
        <div class="p-3.5 rounded-xl border border-slate-200 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <CheckCircle2 v-if="isTrialBalanceBalanced" class="w-5 h-5 text-emerald-600 shrink-0" />
            <AlertCircle v-else class="w-5 h-5 text-rose-500 shrink-0" />
            <div>
              <p class="font-bold text-slate-800">5. Neraca Saldo (Trial Balance) Sempurna Seimbang</p>
              <p class="text-[11px] text-slate-500">
                Total debit dan total kredit buku besar bernilai sama persis tanpa deviasi pembulatan.
              </p>
            </div>
          </div>
          <span :class="['font-bold font-mono px-2.5 py-1 rounded-md text-[11px]', isTrialBalanceBalanced ? 'bg-emerald-50 text-emerald-700' : 'bg-rose-50 text-rose-700']">
            {{ isTrialBalanceBalanced ? 'Seimbang (Balance)' : 'Ada Selisih' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Periods History Table -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
      <div class="p-4 border-b border-slate-100 flex items-center justify-between">
        <h4 class="font-bold text-xs text-slate-900">Riwayat Periode Buku Tahunan 2026</h4>
        <span class="text-[11px] text-slate-500">Tahun Fiskal 2026</span>
      </div>

      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold text-slate-600 uppercase">
            <th class="px-4 py-3">Nama Periode</th>
            <th class="px-4 py-3">Tanggal Mulai</th>
            <th class="px-4 py-3">Tanggal Berakhir</th>
            <th class="px-4 py-3 text-center">Status</th>
            <th class="px-4 py-3">Waktu Penutupan</th>
            <th class="px-4 py-3">Petugas Penutup</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100">
          <tr v-for="p in periods" :key="p.id" class="hover:bg-slate-50/50">
            <td class="px-4 py-3 font-bold text-slate-900">{{ p.periodName }}</td>
            <td class="px-4 py-3 text-slate-600">{{ p.startDate }}</td>
            <td class="px-4 py-3 text-slate-600">{{ p.endDate }}</td>
            <td class="px-4 py-3 text-center">
              <StatusBadge :status="p.status" size="sm" />
            </td>
            <td class="px-4 py-3 text-slate-600">{{ p.lockedAt || '-' }}</td>
            <td class="px-4 py-3 text-slate-600">{{ p.lockedBy || '-' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

