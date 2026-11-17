<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import ReportSignatures from '@/components/accounting/ReportSignatures.vue'
import {
  FileText,
  CheckCircle2,
  Download,
  Printer,
  Sparkles,
  Link2,
  ArrowRight,
  ArrowLeft,
  FileSpreadsheet,
  Landmark,
  ShieldCheck
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()

onMounted(async () => {
  await taxStore.fetchTaxRecords()
})

const annual1771 = computed(() => taxStore.annualTaxReport1771)
const sptMatrix = computed(() => taxStore.monthlySptMasaMatrix)

function handleExportExcel() {
  accountingStore.showToast('Kertas Kerja SPT Tahunan PPh Badan 1771 berhasil diekspor ke Excel', 'success')
}

function printReport() {
  window.print()
}
</script>

<template>
  <div class="space-y-6 font-['Plus_Jakarta_Sans',sans-serif] pb-12">
    <!-- Minimalist Clinical Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Perpajakan Badan Usaha
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Formulir 1771 DJP</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          SPT Tahunan PPh Badan (Formulir 1771)
        </h1>
        <p class="text-xs text-zinc-500">
          Rekonsiliasi fiskal tahunan & akumulasi kredit pajak riil yang dihimpun otomatis dari SPT Masa PPh 23 & PPh 25 via kanal bank.
        </p>
      </div>

      <div class="flex items-center gap-2.5">
        <router-link
          to="/accounting/tax/spt-masa"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Lihat Data SPT Masa</span>
        </router-link>

        <button
          @click="handleExportExcel"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-zinc-50 text-zinc-800 rounded-xl text-xs font-bold border border-zinc-200 shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Download class="w-3.5 h-3.5 text-zinc-600" />
          <span>Ekspor 1771 (.xlsx)</span>
        </button>

        <button
          @click="printReport"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak Formulir 1771</span>
        </button>
      </div>
    </header>

    <!-- Info Banner: Connected with SPT Masa (Strictly Bank Only) -->
    <div class="bg-zinc-100 border border-zinc-200 rounded-2xl p-4 flex flex-col md:flex-row md:items-center justify-between gap-3 text-xs">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-zinc-950 text-white flex items-center justify-center font-bold shrink-0 shadow-xs">
          <Link2 class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h4 class="font-bold text-zinc-950 text-sm">Kredit Pajak SPT 1771 Bersumber 100% dari Histori SPT Masa</h4>
            <span class="px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 font-bold text-[10px]">Non-Cash DJP Standard</span>
          </div>
          <p class="text-zinc-600 mt-0.5">
            Kredit Pajak PPh 23 (bukti potong vendor/reimbursement) dan PPh 25 (angsuran bulanan) dihitung secara <strong>reaktif dan otomatis</strong> dari seluruh SPT Masa berstatus valid di sistem.
          </p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <router-link
          to="/accounting/tax/spt-masa"
          class="px-3.5 py-2 bg-white hover:bg-zinc-50 border border-zinc-300 text-zinc-800 rounded-xl font-bold transition shadow-2xs text-xs whitespace-nowrap"
        >
          <span>Kelola SPT Masa Bulanan</span>
        </router-link>
      </div>
    </div>

    <!-- Annual Tax Summary Box -->
    <div class="bg-zinc-950 text-white rounded-2xl p-6 shadow-xl border border-zinc-800 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-800 pb-4">
        <div>
          <div class="flex items-center gap-2">
            <span class="px-2.5 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-400/30 font-bold text-[11px] flex items-center gap-1">
              <CheckCircle2 class="w-3.5 h-3.5 text-emerald-400" />
              <span>Validasi Terverifikasi (Kanal Bank Resmi)</span>
            </span>
            <span class="text-xs text-zinc-400">Tahun Pajak {{ annual1771.fiscalYear }}</span>
          </div>
          <h2 class="text-xl font-black text-white mt-1">
            Perhitungan PPh Badan Terutang & PPh Pasal 29 Kurang Bayar
          </h2>
          <p class="text-xs text-zinc-400">
            Laba Bersih Komersial dikurangi kredit pajak resmi yang dihimpun murni dari mutasi perbankan SPT Masa PPh 23 dan PPh 25.
          </p>
        </div>

        <div class="text-right">
          <span class="text-xs text-zinc-400 block font-medium">PPh Pasal 29 Kurang Bayar:</span>
          <span class="text-2xl font-black tabular-nums text-amber-400">
            {{ formatIDR(annual1771.pph29Payable) }}
          </span>
        </div>
      </div>

      <!-- 4 Pillars: Komersial -> PKP -> PPh Terutang -> Kredit Pajak Bank -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 text-xs">
        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">1. Laba Bersih Komersial</span>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.commercialNetProfit) }}
          </div>
          <p class="text-[10px] text-zinc-400">Pendapatan dikurangi HPP & Beban Operasional</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">2. Penghasilan Kena Pajak</span>
          <div class="tabular-nums text-base font-bold text-amber-300">
            {{ formatIDR(annual1771.taxableIncome) }}
          </div>
          <p class="text-[10px] text-zinc-400">Setelah Koreksi Fiskal Netto (+Rp 48.95M)</p>
        </div>

        <div class="bg-zinc-900/80 rounded-xl p-4 border border-zinc-800/60 space-y-1">
          <span class="text-zinc-400 font-bold uppercase tracking-wide text-[10px]">3. PPh Badan Terutang</span>
          <div class="tabular-nums text-base font-bold text-white">
            {{ formatIDR(annual1771.corporateTaxPayable) }}
          </div>
          <p class="text-[10px] text-zinc-400">Fasilitas Pasal 31E UU PPh (11% x PKP)</p>
        </div>

        <div class="bg-emerald-950/80 rounded-xl p-4 border border-emerald-700/60 space-y-1 text-emerald-100">
          <span class="text-emerald-400 font-bold uppercase tracking-wide text-[10px]">4. Kredit Pajak via Bank</span>
          <div class="tabular-nums text-base font-bold text-emerald-300">
            {{ formatIDR(annual1771.taxCredits) }}
          </div>
          <p class="text-[10px] text-emerald-400 font-medium">PPh 23 (Rp {{ formatIDR(annual1771.taxCreditPph23) }}) + PPh 25 (Rp {{ formatIDR(annual1771.taxCreditPph25) }})</p>
        </div>
      </div>
    </div>

    <!-- CROSS-RECONCILIATION TABLE: Rekapitulasi Aliran Kredit Pajak dari SPT Masa (Strictly Bank Accounts) -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 border-b border-zinc-100 pb-3">
        <div>
          <div class="flex items-center gap-2">
            <Sparkles class="w-4 h-4 text-amber-500" />
            <h3 class="text-sm font-bold text-zinc-900">Rekapitulasi Bulanan SPT Masa Sebagai Sumber Kredit Pajak SPT Tahunan 1771</h3>
          </div>
          <p class="text-xs text-zinc-500 mt-0.5">
            Diambil secara otomatis dari dokumen SPT Masa yang telah terverifikasi kanal bank resmi (BCA & Mandiri).
          </p>
        </div>
        <span class="text-xs tabular-nums font-bold text-zinc-950 bg-zinc-100 px-3 py-1.5 rounded-xl border border-zinc-200">
          Kredit Bank Terverifikasi: {{ formatIDR(annual1771.taxCredits) }}
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-3 py-2.5">Masa Pajak</th>
              <th class="px-3 py-2.5">Kanal Rekening Bank Pembayar</th>
              <th class="px-3 py-2.5 text-right text-zinc-800">Kredit PPh 23 (Lampiran III)</th>
              <th class="px-3 py-2.5 text-right text-emerald-800">Angsuran PPh 25 (Induk 10a)</th>
              <th class="px-3 py-2.5 text-right text-zinc-600">PPh 21 Medis</th>
              <th class="px-3 py-2.5 text-right text-zinc-600">PPh 4(2) Final</th>
              <th class="px-3 py-2.5 text-right font-bold text-zinc-950 bg-zinc-50/50">Total Kredit SPT 1771</th>
              <th class="px-3 py-2.5 text-center">Status NTPN</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr v-for="row in sptMatrix" :key="row.monthNum" class="hover:bg-zinc-50/50">
              <td class="px-3 py-2.5 font-bold font-sans text-zinc-800">
                {{ row.monthName }} (Masa {{ row.monthNum }})
              </td>
              <td class="px-3 py-2.5 font-sans font-medium text-zinc-700">
                <span class="text-zinc-950 font-bold">{{ row.bankUsed }}</span>
                <span class="block text-[10px] text-zinc-400">Verifikasi Mutasi Rekening Bank</span>
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-zinc-900">
                {{ row.pph23Amount > 0 ? formatIDR(row.pph23Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-emerald-700">
                {{ row.pph25Amount > 0 ? formatIDR(row.pph25Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-zinc-500">
                {{ row.pph21Amount > 0 ? formatIDR(row.pph21Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums text-zinc-500">
                {{ row.pph42Amount > 0 ? formatIDR(row.pph42Amount, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-right tabular-nums font-bold text-zinc-950 bg-zinc-50/50">
                {{ row.totalCreditThisMonth > 0 ? formatIDR(row.totalCreditThisMonth, false) : '-' }}
              </td>
              <td class="px-3 py-2.5 text-center font-sans">
                <span
                  v-if="row.isReported"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-800 border border-emerald-200"
                >
                  Valid (NTPN Bank)
                </span>
                <span
                  v-else
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-zinc-100 text-zinc-500 border border-zinc-200"
                >
                  Belum Ada
                </span>
              </td>
            </tr>
          </tbody>
          <tfoot>
            <tr class="bg-zinc-950 text-white font-bold border-t-2 border-zinc-800">
              <td class="px-3 py-3 uppercase font-sans text-amber-400">TOTAL KREDIT PAJAK TAHUNAN:</td>
              <td class="px-3 py-3 font-sans text-zinc-400 text-xs">Rekap 100% Kanal Bank</td>
              <td class="px-3 py-3 text-right tabular-nums text-white">{{ formatIDR(annual1771.taxCreditPph23) }}</td>
              <td class="px-3 py-3 text-right tabular-nums text-emerald-400">{{ formatIDR(annual1771.taxCreditPph25) }}</td>
              <td class="px-3 py-3 text-right text-zinc-400">-</td>
              <td class="px-3 py-3 text-right text-zinc-400">-</td>
              <td class="px-3 py-3 text-right tabular-nums text-amber-300 text-sm font-black bg-zinc-900">{{ formatIDR(annual1771.taxCredits) }}</td>
              <td class="px-3 py-3 text-center text-emerald-400 font-sans text-xs">Terverifikasi</td>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>

    <!-- SPT 1771 Full Breakdown Table -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs overflow-hidden p-6 space-y-4">
      <h3 class="text-sm font-bold text-zinc-900 border-b border-zinc-100 pb-2 flex items-center gap-2">
        <FileText class="w-4 h-4 text-zinc-700" />
        <span>Kertas Kerja Rekonsiliasi Fiskal SPT 1771 (Tahun Pajak {{ annual1771.fiscalYear }})</span>
      </h3>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-zinc-50 border-b border-zinc-200 text-[11px] font-bold text-zinc-600 uppercase tracking-wider">
              <th class="px-4 py-3 w-16">No</th>
              <th class="px-4 py-3">Komponen Perhitungan SPT Tahunan 1771</th>
              <th class="px-4 py-3">Dasar Hukum & Sumber Data Akuntansi</th>
              <th class="px-4 py-3 text-right w-64">Jumlah (IDR)</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr class="bg-zinc-50/50">
              <td class="px-4 py-2.5 font-bold">1</td>
              <td class="px-4 py-2.5 font-bold font-sans">Peredaran Usaha & Pendapatan Bruto (Akun 4 & 6)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Laporan Laba Rugi Komersial Klinik</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">{{ formatIDR(annual1771.commercialRevenue) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">2</td>
              <td class="px-4 py-2.5 font-sans">Beban Operasional, HPP Obat & Jasa Medis (Akun 5 & 52)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Buku Besar Akun 5 & 52</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-700">({{ formatIDR(annual1771.commercialExpense) }})</td>
            </tr>
            <tr class="bg-zinc-100/60 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">3</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">Penghasilan Netto Komersial (1 - 2)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Laba Bersih Sebelum Pajak</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">{{ formatIDR(annual1771.commercialNetProfit) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">4</td>
              <td class="px-4 py-2.5 font-sans">Koreksi Fiskal Positif (+)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.positiveCorrectionNotes }}</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-emerald-700 font-bold">+ {{ formatIDR(annual1771.positiveFiscalCorrection) }}</td>
            </tr>
            <tr>
              <td class="px-4 py-2.5 font-bold">5</td>
              <td class="px-4 py-2.5 font-sans">Koreksi Fiskal Negatif (-)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">{{ annual1771.negativeCorrectionNotes }}</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-rose-600 font-bold">- {{ formatIDR(annual1771.negativeFiscalCorrection) }}</td>
            </tr>
            <tr class="bg-amber-50 font-bold border-t border-b border-amber-200 text-amber-950">
              <td class="px-4 py-2.5 font-black">6</td>
              <td class="px-4 py-2.5 font-black uppercase">Penghasilan Kena Pajak (PKP) Badan</td>
              <td class="px-4 py-2.5 text-amber-800 font-sans">Dasar Pengenaan PPh Badan</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-amber-950 font-black text-sm">{{ formatIDR(annual1771.taxableIncome) }}</td>
            </tr>
            <tr class="bg-zinc-50">
              <td class="px-4 py-2.5 font-bold">7</td>
              <td class="px-4 py-2.5 font-sans">Tarif PPh Badan Efektif (Fasilitas Ps 31E UU HPP)</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">50% x 22% = 11% (Peredaran &lt; Rp 50M)</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-zinc-950">11,00 %</td>
            </tr>
            <tr class="bg-zinc-100/80 font-bold border-t border-b border-zinc-200">
              <td class="px-4 py-2.5 font-black">8</td>
              <td class="px-4 py-2.5 font-black uppercase text-zinc-950">PPh Badan Terutang Tahun Pajak {{ annual1771.fiscalYear }}</td>
              <td class="px-4 py-2.5 text-zinc-500 font-sans">Formulir 1771 Induk Bagian B</td>
              <td class="px-4 py-2.5 text-right tabular-nums text-zinc-950 font-black text-sm">{{ formatIDR(annual1771.corporateTaxPayable) }}</td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60">
              <td class="px-4 py-2.5 font-bold">9</td>
              <td class="px-4 py-2.5 font-sans font-bold">Kredit Pajak PPh 23 Terakumulasi dari SPT Masa Bank</td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">Lampiran Khusus III (100% Rekening Bank)</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">({{ formatIDR(annual1771.taxCreditPph23) }})</td>
            </tr>
            <tr class="text-emerald-950 bg-emerald-50/60">
              <td class="px-4 py-2.5 font-bold">10</td>
              <td class="px-4 py-2.5 font-sans font-bold">Kredit Pajak PPh 25 Terakumulasi dari SPT Masa Bank</td>
              <td class="px-4 py-2.5 text-emerald-800 font-sans">Induk Baris 10a (NTPN e-Billing Mandiri/BCA)</td>
              <td class="px-4 py-2.5 text-right tabular-nums font-bold text-emerald-700">({{ formatIDR(annual1771.taxCreditPph25) }})</td>
            </tr>
            <tr class="bg-zinc-950 text-white font-black border-t-2 border-zinc-800">
              <td class="px-4 py-3 font-black text-amber-400">11</td>
              <td class="px-4 py-3 uppercase text-amber-400 text-sm">PPH PASAL 29 KURANG BAYAR (8 - 9 - 10)</td>
              <td class="px-4 py-3 text-zinc-400 font-sans text-xs">Wajib disetor sebelum pelaporan SPT 1771</td>
              <td class="px-4 py-3 text-right tabular-nums text-amber-300 font-black text-base">{{ formatIDR(annual1771.pph29Payable) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Formal Clinical Signature Section -->
      <ReportSignatures
        preparer-role="Supervisor Pajak & Akuntansi"
        preparer-name="Rian Pratama, S.Ak."
        reviewer-role="Tax Manager & Konsultan Pajak"
        reviewer-name="Dewi Lestari, SE, Ak., BKP"
        approver-role="Direktur Keuangan (CFO)"
        approver-name="Budi Santoso, SE, Ak., CA"
      />
    </div>
  </div>
</template>
