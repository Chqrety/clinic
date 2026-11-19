<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportCashFlowToExcel, exportToExcel, printReport } from '@/utils/export'
import {
  Download,
  Printer,
  BookOpen,
  Layers,
  Calculator
} from 'lucide-vue-next'

const router = useRouter()
const accountingStore = useAccountingStore()

// State
const activeTab = ref<'depreciation' | 'cashflow'>('cashflow')
const startDate = ref('01-11-2026')
const endDate = ref('30-11-2026')

interface CashFlowRow {
  lineNo: number
  name: string
  amountCurrent: number
  amountPrevious: number
  isSectionHeader?: boolean
  isBold?: boolean
  isLink?: boolean
  accountCode?: string
}

// Cash Flow Statement matching Image 3 and Clinic ERP Context
const cashFlowRows = ref<CashFlowRow[]>([
  { lineNo: 1, name: 'ARUS KAS DARI AKTIVITAS OPERASIONAL:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 2, name: 'Penerimaan Kasir Harian dari Pasien (Tunai, QRIS & EDC Bank)', amountCurrent: 660108120.00, amountPrevious: 175000000.00, isLink: true, accountCode: '1-1111' },
  { lineNo: 3, name: 'Pembayaran Pengadaan Obat, Vaksin & Skincare Distributor', amountCurrent: -81420000.00, amountPrevious: -45000000.00, isLink: true, accountCode: '2-1110' },
  { lineNo: 4, name: 'Pembayaran Honor Jasa Medis Dokter Spesialis & Terapis', amountCurrent: -163000000.00, amountPrevious: -65000000.00, isLink: true, accountCode: '6-1100' },
  { lineNo: 5, name: 'Pembayaran Gaji Karyawan, Perawat & Apoteker', amountCurrent: -89000000.00, amountPrevious: -44000000.00, isLink: true, accountCode: '6-1200' },
  { lineNo: 6, name: 'Pembayaran Operasional (Listrik, Sewa Ruko & Limbah Medis B3)', amountCurrent: -156500000.00, amountPrevious: -42000000.00, isLink: true, accountCode: '6-1300' },
  { lineNo: 7, name: 'Pembayaran Setoran Pajak Masa (PPh 21 Dokter, PPh 23, PPh 4(2))', amountCurrent: -9147500.00, amountPrevious: -2420000.00, isLink: true, accountCode: '2-1500' },
  { lineNo: 8, name: 'Beban Administrasi Bank & Sewa Mesin EDC (BCA/BRI)', amountCurrent: -1600000.00, amountPrevious: -750000.00, isLink: true, accountCode: '6-1400' },
  { lineNo: 9, name: 'Arus Kas Bersih yang Diperoleh dari Aktivitas Operasional', amountCurrent: 159440620.00, amountPrevious: -24170000.00, isBold: true },

  { lineNo: 11, name: 'ARUS KAS DARI AKTIVITAS INVESTASI:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 12, name: 'Pengadaan Alat Medis USG Pediatri & Mesin Facial Estetika', amountCurrent: -35000000.00, amountPrevious: 0.00, isLink: true, accountCode: '1-2100' },
  { lineNo: 13, name: 'Kas Bersih Digunakan untuk Aktivitas Investasi', amountCurrent: -35000000.00, amountPrevious: 0.00, isBold: true },

  { lineNo: 15, name: 'ARUS KAS DARI AKTIVITAS PENDANAAN:', amountCurrent: 0, amountPrevious: 0, isSectionHeader: true, isBold: true },
  { lineNo: 16, name: 'Penyetoran Tambahan Modal Kerja Pemilik Klinik', amountCurrent: 50000000.00, amountPrevious: 100000000.00, isLink: true, accountCode: '3-1100' },
  { lineNo: 17, name: 'Kas Bersih Diperoleh dari Aktivitas Pendanaan', amountCurrent: 50000000.00, amountPrevious: 100000000.00, isBold: true },

  { lineNo: 19, name: 'KENAIKAN (PENURUNAN) BERSIH KAS DAN SETARA KAS', amountCurrent: 174440620.00, amountPrevious: 75830000.00, isBold: true },
  { lineNo: 20, name: 'KAS DAN SETARA KAS PADA AWAL PERIODE', amountCurrent: 105000000.00, amountPrevious: 29170000.00, isBold: true },
  { lineNo: 21, name: 'KAS DAN SETARA KAS PADA AKHIR PERIODE (KAS, BCA, BRI, KLIRING)', amountCurrent: 279440620.00, amountPrevious: 105000000.00, isBold: true }
])

function openLedgerForAccount(accountCode?: string) {
  if (accountCode) {
    router.push({ path: '/accounting/ledger', query: { account: accountCode } })
  } else {
    router.push('/accounting/ledger')
  }
}

function handleExportExcel() {
  exportCashFlowToExcel(cashFlowRows.value, {
    businessUnit: accountingStore.businessUnit,
    asOfDate: endDate.value,
    startDate: startDate.value
  })
  accountingStore.showToast('Laporan Arus Kas Audited (Format LHP SPHC) berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-4 font-['Plus_Jakarta_Sans',sans-serif]">
    <!-- Top Header Tabs matching Clinical Light Design -->
    <div class="flex items-center gap-2 border-b border-zinc-200 pb-2">
      <button
        @click="activeTab = 'depreciation'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-xl transition-all cursor-pointer',
          activeTab === 'depreciation'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'
        ]"
      >
        Penyusutan dan Amortisasi
      </button>
      <button
        @click="activeTab = 'cashflow'"
        :class="[
          'px-4 py-2 text-xs font-bold rounded-xl transition-all cursor-pointer',
          activeTab === 'cashflow'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100'
        ]"
      >
        Arus Kas
      </button>
    </div>

    <!-- Header matching Clinical Light Design -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl p-4 sm:p-5 text-zinc-900 flex flex-col md:flex-row md:items-center justify-between gap-4 shadow-2xs">
      <div class="space-y-0.5">
        <div class="flex items-center gap-2 text-xs text-zinc-400 font-medium">
          <span class="hover:text-zinc-700 cursor-pointer">Accounting</span>
          <span>/</span>
          <span class="text-zinc-900 font-bold">Laporan Arus Kas</span>
        </div>
        <h1 class="text-base sm:text-lg font-bold text-zinc-950 tracking-tight">
          {{ accountingStore.businessUnit === 'ALL' ? 'KLINIK TERPADU (KIDINEST & SSC)' : accountingStore.businessUnit === 'KIDINEST' ? 'KLINIK SPESIALIS ANAK KIDINEST' : 'SEMARANG SKIN CENTER (SSC)' }}
        </h1>
      </div>

      <div class="flex flex-wrap items-center gap-2.5 text-xs">
        <!-- Date Range -->
        <div class="flex items-center gap-1.5 bg-zinc-50 px-3 py-1.5 rounded-xl border border-zinc-200 text-zinc-700 font-medium">
          <span class="text-zinc-500">Periode:</span>
          <input
            v-model="startDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
          <span class="text-zinc-400">s/d</span>
          <input
            v-model="endDate"
            type="text"
            class="bg-transparent font-mono font-bold text-zinc-900 w-24 text-center focus:outline-none"
          />
        </div>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor</span>
        </button>

        <button
          @click="printReport"
          class="inline-flex items-center gap-1.5 px-3.5 py-1.5 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak</span>
        </button>
      </div>
    </div>

    <!-- Tree Table Container (Clean Light Theme) -->
    <div class="bg-white border border-zinc-200/90 rounded-2xl overflow-hidden shadow-2xs text-zinc-800">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs select-none">
          <thead>
            <tr class="border-b border-zinc-200 bg-zinc-50/90 text-[11px] font-bold text-zinc-500 uppercase tracking-wider">
              <th class="py-2.5 px-4 w-16 text-zinc-400">#</th>
              <th class="py-2.5 px-4">Uraian</th>
              <th class="py-2.5 px-6 text-right w-56">November 2026</th>
              <th class="py-2.5 px-6 text-right w-56 text-zinc-400">November 2025</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 font-mono">
            <tr
              v-for="row in cashFlowRows"
              :key="row.lineNo"
              :class="[
                'transition-colors group',
                row.isSectionHeader
                  ? 'bg-zinc-100/90 font-extrabold text-zinc-950 text-xs uppercase tracking-wide border-b border-zinc-200'
                  : row.isBold
                  ? 'bg-zinc-50/80 font-bold text-zinc-900 border-b border-zinc-100'
                  : 'hover:bg-zinc-50/80 text-zinc-700'
              ]"
            >
              <!-- Line Number -->
              <td class="py-2 px-4 text-zinc-400 text-[11px] select-none font-mono">
                {{ row.lineNo }}
              </td>

              <!-- Uraian -->
              <td class="py-2 px-4 font-sans">
                <div class="flex items-center gap-1.5" :class="{ 'pl-4': !row.isSectionHeader && !row.isBold }">
                  <!-- Clickable Blue Link Style -->
                  <span
                    v-if="row.isLink"
                    @click="row.accountCode && openLedgerForAccount(row.accountCode)"
                    class="text-blue-600 hover:underline cursor-pointer font-medium flex items-center gap-1"
                  >
                    <span>{{ row.name }}</span>
                    <BookOpen v-if="row.accountCode" class="w-3 h-3 text-zinc-400 group-hover:text-blue-600 inline" />
                  </span>
                  <span
                    v-else
                    :class="[
                      row.isSectionHeader ? 'font-black uppercase tracking-wider text-zinc-950' : '',
                      row.isBold && !row.isSectionHeader ? 'font-bold text-zinc-900' : ''
                    ]"
                  >
                    {{ row.name }}
                  </span>
                </div>
              </td>

              <!-- November 2026 -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight',
                  row.isSectionHeader ? 'text-transparent' : '',
                  row.isBold ? 'font-bold text-zinc-950' : 'font-medium text-zinc-800'
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountCurrent) }}
              </td>

              <!-- November 2025 -->
              <td
                :class="[
                  'py-2 px-6 text-right font-mono tracking-tight text-zinc-400',
                  row.isSectionHeader ? 'text-transparent' : ''
                ]"
              >
                {{ row.isSectionHeader ? '' : formatIDR(row.amountPrevious) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Bottom Status Note (Clean Light Theme) -->
      <div class="px-4 py-3 bg-zinc-50/80 border-t border-zinc-200 flex items-center justify-between text-[11px] text-zinc-500">
        <span class="font-mono">
          Perbandingan arus kas tahun berjalan vs periode sebelumnya.
        </span>
        <span class="text-emerald-700 font-bold font-mono bg-emerald-50 px-2 py-0.5 rounded-md border border-emerald-200">
          Kenaikan Kas Bersih: {{ formatIDR(106261500) }}
        </span>
      </div>
    </div>
  </div>
</template>
