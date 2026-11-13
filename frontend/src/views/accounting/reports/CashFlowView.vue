<script setup lang="ts">
import { computed } from 'vue'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import {
  Coins,
  Download,
  Printer,
  ArrowUpRight,
  ArrowDownLeft,
  Wallet
} from 'lucide-vue-next'

const accountingStore = useAccountingStore()

// Operating Cash Flow lines
const cashFromPatients = 492000000 // Inflow
const paymentToDistributors = -48500000 // Outflow
const paymentToDoctors = -138000000 // Outflow
const paymentToEmployees = -68000000 // Outflow
const paymentForRentAndUtilities = -34800000 // Outflow
const paymentForTaxes = -14850000 // Outflow

const netOperatingCash = computed(() => {
  return (
    cashFromPatients +
    paymentToDistributors +
    paymentToDoctors +
    paymentToEmployees +
    paymentForRentAndUtilities +
    paymentForTaxes
  )
})

// Investing Cash Flow lines
const purchaseMedicalEquipment = -25000000 // e.g. Probe USG baru
const netInvestingCash = computed(() => purchaseMedicalEquipment)

// Financing Cash Flow lines
const ownerDrawings = -10000000 // Prive
const netFinancingCash = computed(() => ownerDrawings)

// Net Change
const netCashChange = computed(() => {
  return netOperatingCash.value + netInvestingCash.value + netFinancingCash.value
})

const beginningCash = 354200000
const endingCash = computed(() => beginningCash + netCashChange.value)

function handleExportExcel() {
  const data = [
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Penerimaan Kasir & Pasien Klinik', 'Nominal (IDR)': cashFromPatients },
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Pembayaran Obat ke Distributor Farmasi', 'Nominal (IDR)': paymentToDistributors },
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Pembayaran Jasa Medis Dokter', 'Nominal (IDR)': paymentToDoctors },
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Pembayaran Gaji Pegawai & Perawat', 'Nominal (IDR)': paymentToEmployees },
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Pembayaran Sewa & Utilitas Operasional', 'Nominal (IDR)': paymentForRentAndUtilities },
    { 'Aktivitas': 'ARUS KAS OPERASIONAL', 'Keterangan': 'Pembayaran Pajak PPh & PPN', 'Nominal (IDR)': paymentForTaxes },
    { 'Aktivitas': 'SUBTOTAL', 'Keterangan': 'Kas Bersih dari Aktivitas Operasi', 'Nominal (IDR)': netOperatingCash.value },
    { 'Aktivitas': 'ARUS KAS INVESTASI', 'Keterangan': 'Peremajaan & Aksesoris Alat Medis', 'Nominal (IDR)': purchaseMedicalEquipment },
    { 'Aktivitas': 'SUBTOTAL', 'Keterangan': 'Kas Bersih untuk Aktivitas Investasi', 'Nominal (IDR)': netInvestingCash.value },
    { 'Aktivitas': 'ARUS KAS PENDANAAN', 'Keterangan': 'Prive Pemilik Klinik', 'Nominal (IDR)': ownerDrawings },
    { 'Aktivitas': 'SUBTOTAL', 'Keterangan': 'Kas Bersih untuk Aktivitas Pendanaan', 'Nominal (IDR)': netFinancingCash.value },
    { 'Aktivitas': 'TOTAL', 'Keterangan': 'Kenaikan Bersih Kas & Bank', 'Nominal (IDR)': netCashChange.value },
    { 'Aktivitas': 'SALDO', 'Keterangan': 'Saldo Kas Awal Periode', 'Nominal (IDR)': beginningCash },
    { 'Aktivitas': 'SALDO', 'Keterangan': 'Saldo Kas Akhir Periode', 'Nominal (IDR)': endingCash.value }
  ]

  exportToExcel(data, `Laporan_Arus_Kas_${accountingStore.businessUnit}_2026-03`, 'Arus Kas')
  accountingStore.showToast('Laporan Arus Kas berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs print:hidden">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Statement of Cash Flows
          </span>
          <span class="text-xs text-slate-400">&bull; Metode Langsung (Direct Method)</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Laporan Arus Kas (Cash Flow)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Melacak aliran kas masuk dan kas keluar dari aktivitas operasi klinik, investasi peralatan medik, dan aktivitas pendanaan.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-xs font-semibold border border-slate-200 transition-colors"
        >
          <Printer class="w-4 h-4 text-slate-600" />
          <span>Cetak Laporan</span>
        </button>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-4 h-4" />
          <span>Ekspor Excel</span>
        </button>
      </div>
    </div>

    <!-- Presentation Sheet -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs p-8 max-w-4xl mx-auto space-y-6 text-xs text-slate-800">
      <div class="text-center pb-6 border-b border-slate-200 space-y-1">
        <h2 class="text-lg font-extrabold uppercase tracking-wide text-slate-900">
          KLINIK TERPADU KIDINEST & SKIN SLIMMING CLINIC (SSC)
        </h2>
        <h3 class="text-sm font-bold text-blue-600 uppercase tracking-wider">
          LAPORAN ARUS KAS (STATEMENT OF CASH FLOWS)
        </h3>
        <p class="text-xs text-slate-500">
          Untuk Periode yang Berakhir pada 31 Maret 2026 (Dinyatakan dalam Rupiah Indonesia)
        </p>
      </div>

      <!-- 1. Operating Activities -->
      <div class="space-y-3">
        <div class="font-bold text-sm text-slate-900 border-b border-slate-200 pb-1.5 flex items-center justify-between">
          <span>ARUS KAS DARI AKTIVITAS OPERASI</span>
          <span class="text-[11px] font-normal text-slate-500">Metode Langsung</span>
        </div>

        <div class="space-y-2 pl-4">
          <div class="flex justify-between">
            <span>Penerimaan kas dari pasien & klaim asuransi</span>
            <span class="font-mono text-emerald-700 font-semibold">{{ formatIDR(cashFromPatients) }}</span>
          </div>
          <div class="flex justify-between">
            <span>Pembayaran kas ke distributor farmasi & obat</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(paymentToDistributors)) }})</span>
          </div>
          <div class="flex justify-between">
            <span>Pembayaran kompensasi & jasa medis dokter</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(paymentToDoctors)) }})</span>
          </div>
          <div class="flex justify-between">
            <span>Pembayaran gaji, upah & tunjangan karyawan</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(paymentToEmployees)) }})</span>
          </div>
          <div class="flex justify-between">
            <span>Pembayaran sewa gedung, utilitas listrik & air</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(paymentForRentAndUtilities)) }})</span>
          </div>
          <div class="flex justify-between">
            <span>Pembayaran setoran pajak penghasilan & PPN</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(paymentForTaxes)) }})</span>
          </div>
        </div>

        <div class="pt-2 border-t border-slate-200 flex justify-between font-bold bg-blue-50/60 p-2.5 rounded-lg text-blue-950">
          <span>Arus Kas Bersih yang Diperoleh dari Aktivitas Operasi:</span>
          <span class="font-mono font-bold">{{ formatIDR(netOperatingCash) }}</span>
        </div>
      </div>

      <!-- 2. Investing Activities -->
      <div class="space-y-3 pt-3">
        <div class="font-bold text-sm text-slate-900 border-b border-slate-200 pb-1.5 flex items-center justify-between">
          <span>ARUS KAS DARI AKTIVITAS INVESTASI</span>
        </div>

        <div class="space-y-2 pl-4">
          <div class="flex justify-between">
            <span>Pembelian aksesoris & perlengkapan alat medis</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(purchaseMedicalEquipment)) }})</span>
          </div>
        </div>

        <div class="pt-2 border-t border-slate-200 flex justify-between font-bold bg-slate-50 p-2.5 rounded-lg text-slate-900">
          <span>Arus Kas Bersih yang Digunakan untuk Aktivitas Investasi:</span>
          <span class="font-mono font-bold text-rose-800">({{ formatIDR(Math.abs(netInvestingCash)) }})</span>
        </div>
      </div>

      <!-- 3. Financing Activities -->
      <div class="space-y-3 pt-3">
        <div class="font-bold text-sm text-slate-900 border-b border-slate-200 pb-1.5 flex items-center justify-between">
          <span>ARUS KAS DARI AKTIVITAS PENDANAAN</span>
        </div>

        <div class="space-y-2 pl-4">
          <div class="flex justify-between">
            <span>Prive / Penarikan dana pemilik klinik</span>
            <span class="font-mono text-rose-700">({{ formatIDR(Math.abs(ownerDrawings)) }})</span>
          </div>
        </div>

        <div class="pt-2 border-t border-slate-200 flex justify-between font-bold bg-slate-50 p-2.5 rounded-lg text-slate-900">
          <span>Arus Kas Bersih yang Digunakan untuk Aktivitas Pendanaan:</span>
          <span class="font-mono font-bold text-rose-800">({{ formatIDR(Math.abs(netFinancingCash)) }})</span>
        </div>
      </div>

      <!-- Summary & Reconciliation -->
      <div class="pt-6 border-t-2 border-slate-300 space-y-2.5 bg-slate-50 p-4 rounded-xl">
        <div class="flex justify-between font-bold text-slate-900 text-sm">
          <span>KENAIKAN BERSIH KAS & SETARA KAS:</span>
          <span class="font-mono text-emerald-700">{{ formatIDR(netCashChange) }}</span>
        </div>
        <div class="flex justify-between text-slate-600">
          <span>Saldo Kas & Setara Kas pada Awal Periode:</span>
          <span class="font-mono font-semibold">{{ formatIDR(beginningCash) }}</span>
        </div>
        <div class="flex justify-between font-extrabold text-sm text-blue-900 pt-2 border-t border-slate-300">
          <span>SALDO KAS & SETARA KAS PADA AKHIR PERIODE:</span>
          <span class="font-mono font-extrabold">{{ formatIDR(endingCash) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
