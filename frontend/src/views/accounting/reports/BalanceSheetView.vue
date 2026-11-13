<script setup lang="ts">
import { computed } from 'vue'
import { useReportsStore } from '@/stores/reports'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import { exportToExcel, printReport } from '@/utils/export'
import {
  Scale,
  Download,
  Printer,
  CheckCircle2,
  AlertCircle
} from 'lucide-vue-next'

const reportsStore = useReportsStore()
const accountingStore = useAccountingStore()

const bs = computed(() => reportsStore.balanceSheet)

function handleExportExcel() {
  const data = [
    { 'Kategori': 'ASET LANCAR', 'Akun': 'Kas Operasional (Kidinest & SSC)', 'Jumlah (IDR)': bs.value.kasOperasional },
    { 'Kategori': 'ASET LANCAR', 'Akun': 'Bank Operasional (BCA & Mandiri)', 'Jumlah (IDR)': bs.value.bankOperasional },
    { 'Kategori': 'ASET LANCAR', 'Akun': 'Piutang Pasien & Asuransi', 'Jumlah (IDR)': bs.value.piutangPasien },
    { 'Kategori': 'ASET LANCAR', 'Akun': 'Persediaan Farmasi & Obat', 'Jumlah (IDR)': bs.value.persediaan },
    { 'Kategori': 'ASET LANCAR', 'Akun': 'PPN Masukan (Dapat Dikreditkan)', 'Jumlah (IDR)': bs.value.ppnMasukan },
    { 'Kategori': 'SUBTOTAL', 'Akun': 'TOTAL ASET LANCAR', 'Jumlah (IDR)': bs.value.totalAsetLancar },
    { 'Kategori': 'ASET TETAP', 'Akun': 'Peralatan Medis & Laser (Biaya Perolehan)', 'Jumlah (IDR)': bs.value.asetTetapGross },
    { 'Kategori': 'ASET TETAP', 'Akun': 'Akumulasi Penyusutan Peralatan Medis', 'Jumlah (IDR)': -bs.value.akumulasiPenyusutan },
    { 'Kategori': 'SUBTOTAL', 'Akun': 'TOTAL ASET TETAP (NET)', 'Jumlah (IDR)': bs.value.totalAsetTetapNet },
    { 'Kategori': 'TOTAL', 'Akun': 'TOTAL ASET KESELURUHAN', 'Jumlah (IDR)': bs.value.totalAset },
    { 'Kategori': 'KEWAJIBAN', 'Akun': 'Hutang Usaha Distributor Farmasi', 'Jumlah (IDR)': bs.value.hutangDistributor },
    { 'Kategori': 'KEWAJIBAN', 'Akun': 'Hutang Jasa Medis Dokter', 'Jumlah (IDR)': bs.value.hutangDokter },
    { 'Kategori': 'KEWAJIBAN', 'Akun': 'Hutang Pajak PPh 21, 23 & PPN Keluaran', 'Jumlah (IDR)': bs.value.hutangPajak },
    { 'Kategori': 'SUBTOTAL', 'Akun': 'TOTAL KEWAJIBAN', 'Jumlah (IDR)': bs.value.totalKewajiban },
    { 'Kategori': 'EKUITAS', 'Akun': 'Modal Disetor Pendiri', 'Jumlah (IDR)': bs.value.modalDisetor },
    { 'Kategori': 'EKUITAS', 'Akun': 'Laba Ditahan (Retained Earnings)', 'Jumlah (IDR)': bs.value.labaDitahan },
    { 'Kategori': 'EKUITAS', 'Akun': 'Laba Periode Berjalan', 'Jumlah (IDR)': bs.value.labaTahunBerjalan },
    { 'Kategori': 'SUBTOTAL', 'Akun': 'TOTAL EKUITAS', 'Jumlah (IDR)': bs.value.totalEkuitas },
    { 'Kategori': 'TOTAL', 'Akun': 'TOTAL KEWAJIBAN & EKUITAS', 'Jumlah (IDR)': bs.value.totalKewajibanDanEkuitas }
  ]

  exportToExcel(data, `Laporan_Neraca_${accountingStore.businessUnit}_2026-03`, 'Neraca')
  accountingStore.showToast('Laporan Neraca berhasil diekspor ke Excel', 'success')
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs print:hidden">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold text-blue-600 uppercase tracking-wider bg-blue-50 px-2 py-0.5 rounded border border-blue-200">
            Laporan Posisi Keuangan
          </span>
          <span class="text-xs text-slate-400">&bull; Balance Sheet (Neraca)</span>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mt-1">
          Neraca Posisi Keuangan (Balance Sheet)
        </h1>
        <p class="text-xs text-slate-500 mt-0.5">
          Menampilkan aset, liabilitas, dan ekuitas dengan pembuktian formula fundamental: <strong>Aset = Kewajiban + Ekuitas</strong>.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          @click="printReport"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-xs font-semibold border border-slate-200 transition-colors"
        >
          <Printer class="w-4 h-4 text-slate-600" />
          <span>Cetak Neraca</span>
        </button>

        <button
          @click="handleExportExcel"
          class="flex items-center gap-1.5 px-3.5 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-xs font-bold shadow-xs transition-colors"
        >
          <Download class="w-4 h-4" />
          <span>Ekspor Excel (.xlsx)</span>
        </button>
      </div>
    </div>

    <!-- Presentation Sheet -->
    <div class="bg-white rounded-2xl border border-slate-200/80 shadow-2xs p-8 max-w-5xl mx-auto space-y-6">
      <!-- Title -->
      <div class="text-center pb-6 border-b border-slate-200 space-y-1">
        <h2 class="text-lg font-extrabold uppercase tracking-wide text-slate-900">
          KLINIK TERPADU KIDINEST & SKIN SLIMMING CLINIC (SSC)
        </h2>
        <h3 class="text-sm font-bold text-blue-600 uppercase tracking-wider">
          LAPORAN POSISI KEUANGAN (NERACA)
        </h3>
        <p class="text-xs text-slate-500">
          Per 31 Maret 2026 (Dinyatakan dalam Rupiah Indonesia)
        </p>
      </div>

      <!-- Balance Integrity Banner -->
      <div
        v-if="bs.isBalanced"
        class="p-4 bg-emerald-50 border border-emerald-200 rounded-xl flex items-center justify-between text-emerald-900 text-xs"
      >
        <div class="flex items-center gap-2.5">
          <CheckCircle2 class="w-5 h-5 text-emerald-600" />
          <div>
            <h4 class="font-bold">Neraca Seimbang Sempurna (Balanced)</h4>
            <p class="text-[11px] text-emerald-700">Total Aset = Total Kewajiban + Ekuitas (Selisih: Rp 0)</p>
          </div>
        </div>
        <span class="font-mono font-bold text-sm">{{ formatIDR(bs.totalAset) }}</span>
      </div>

      <!-- Two-Column Financial Layout: Left = Assets, Right = Liabilities & Equity -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 text-xs">
        <!-- LEFT COLUMN: ASET -->
        <div class="space-y-6">
          <div class="border-b-2 border-blue-600 pb-2">
            <h3 class="font-bold text-sm text-slate-900 uppercase tracking-wider">ASET (AKTIVA)</h3>
          </div>

          <!-- Aset Lancar -->
          <div class="space-y-3">
            <h4 class="font-bold text-blue-900 uppercase text-[11px]">Aset Lancar:</h4>
            <div class="space-y-2 pl-3">
              <div class="flex justify-between text-slate-700">
                <span>Kas Operasional Kasir (Kidinest & SSC)</span>
                <span class="font-mono">{{ formatIDR(bs.kasOperasional) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Bank Mandiri & BCA Operasional</span>
                <span class="font-mono">{{ formatIDR(bs.bankOperasional) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Piutang Pasien & Asuransi</span>
                <span class="font-mono">{{ formatIDR(bs.piutangPasien) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Persediaan Farmasi & Obat-obatan</span>
                <span class="font-mono">{{ formatIDR(bs.persediaan) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>PPN Masukan Dapat Dikreditkan</span>
                <span class="font-mono">{{ formatIDR(bs.ppnMasukan) }}</span>
              </div>
            </div>
            <!-- Subtotal Aset Lancar -->
            <div class="pt-2 border-t border-slate-200 flex justify-between font-bold text-slate-900 bg-slate-50 p-2 rounded">
              <span>Total Aset Lancar:</span>
              <span class="font-mono">{{ formatIDR(bs.totalAsetLancar) }}</span>
            </div>
          </div>

          <!-- Aset Tetap -->
          <div class="space-y-3">
            <h4 class="font-bold text-blue-900 uppercase text-[11px]">Aset Tetap & Peralatan Medis:</h4>
            <div class="space-y-2 pl-3">
              <div class="flex justify-between text-slate-700">
                <span>Peralatan Medis & Laser (Biaya Perolehan)</span>
                <span class="font-mono">{{ formatIDR(bs.asetTetapGross) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Akumulasi Penyusutan Peralatan</span>
                <span class="font-mono text-rose-700">({{ formatIDR(bs.akumulasiPenyusutan) }})</span>
              </div>
            </div>
            <!-- Subtotal Aset Tetap -->
            <div class="pt-2 border-t border-slate-200 flex justify-between font-bold text-slate-900 bg-slate-50 p-2 rounded">
              <span>Total Aset Tetap (Nilai Buku):</span>
              <span class="font-mono">{{ formatIDR(bs.totalAsetTetapNet) }}</span>
            </div>
          </div>

          <!-- Total Aset -->
          <div class="pt-4 border-t-2 border-slate-300 flex justify-between font-extrabold text-sm text-white bg-blue-900 p-3 rounded-xl shadow-xs">
            <span>TOTAL ASET (AKTIVA):</span>
            <span class="font-mono">{{ formatIDR(bs.totalAset) }}</span>
          </div>
        </div>

        <!-- RIGHT COLUMN: LIABILITIES & EQUITY -->
        <div class="space-y-6">
          <div class="border-b-2 border-rose-600 pb-2">
            <h3 class="font-bold text-sm text-slate-900 uppercase tracking-wider">KEWAJIBAN & EKUITAS</h3>
          </div>

          <!-- Kewajiban Jangka Pendek -->
          <div class="space-y-3">
            <h4 class="font-bold text-rose-900 uppercase text-[11px]">Kewajiban Jangka Pendek:</h4>
            <div class="space-y-2 pl-3">
              <div class="flex justify-between text-slate-700">
                <span>Hutang Usaha Distributor Farmasi</span>
                <span class="font-mono">{{ formatIDR(bs.hutangDistributor) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Hutang Jasa Medis Dokter (Share Medis)</span>
                <span class="font-mono">{{ formatIDR(bs.hutangDokter) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Hutang Pajak PPh 21, 23 & PPN Keluaran</span>
                <span class="font-mono">{{ formatIDR(bs.hutangPajak) }}</span>
              </div>
            </div>
            <!-- Subtotal Kewajiban -->
            <div class="pt-2 border-t border-slate-200 flex justify-between font-bold text-slate-900 bg-slate-50 p-2 rounded">
              <span>Total Kewajiban:</span>
              <span class="font-mono">{{ formatIDR(bs.totalKewajiban) }}</span>
            </div>
          </div>

          <!-- Ekuitas -->
          <div class="space-y-3">
            <h4 class="font-bold text-blue-900 uppercase text-[11px]">Ekuitas (Modal):</h4>
            <div class="space-y-2 pl-3">
              <div class="flex justify-between text-slate-700">
                <span>Modal Disetor Pendiri</span>
                <span class="font-mono">{{ formatIDR(bs.modalDisetor) }}</span>
              </div>
              <div class="flex justify-between text-slate-700">
                <span>Laba Ditahan (Retained Earnings)</span>
                <span class="font-mono">{{ formatIDR(bs.labaDitahan) }}</span>
              </div>
              <div class="flex justify-between text-slate-700 font-semibold text-emerald-800">
                <span>Laba Bersih Tahun Berjalan</span>
                <span class="font-mono">{{ formatIDR(bs.labaTahunBerjalan) }}</span>
              </div>
            </div>
            <!-- Subtotal Ekuitas -->
            <div class="pt-2 border-t border-slate-200 flex justify-between font-bold text-slate-900 bg-slate-50 p-2 rounded">
              <span>Total Ekuitas:</span>
              <span class="font-mono">{{ formatIDR(bs.totalEkuitas) }}</span>
            </div>
          </div>

          <!-- Total Kewajiban & Ekuitas -->
          <div class="pt-4 border-t-2 border-slate-300 flex justify-between font-extrabold text-sm text-white bg-slate-900 p-3 rounded-xl shadow-xs">
            <span>TOTAL KEWAJIBAN & EKUITAS:</span>
            <span class="font-mono">{{ formatIDR(bs.totalKewajibanDanEkuitas) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
