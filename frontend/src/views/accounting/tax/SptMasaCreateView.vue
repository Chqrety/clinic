<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTaxStore, type SptMasaType } from '@/stores/tax'
import { useAccountingStore } from '@/stores/accounting'
import { formatIDR } from '@/utils/currency'
import type { BusinessUnit } from '@/types/accounting'
import {
  ArrowLeft,
  Calculator,
  Save,
  CheckCircle2,
  Building2,
  Calendar,
  Landmark,
  FileText,
  ShieldCheck,
  AlertCircle
} from 'lucide-vue-next'

const router = useRouter()
const taxStore = useTaxStore()
const accountingStore = useAccountingStore()

const masaMonth = ref<number>(accountingStore.fiscalMonth || 3)
const masaYear = ref<number>(accountingStore.fiscalYear || 2026)
const unit = ref<BusinessUnit>(accountingStore.businessUnit === 'SSC' ? 'SSC' : 'KIDINEST')
const taxType = ref<SptMasaType>('PPH_23')
const bankAccountCode = ref<string>('1102.001')
const grossAmount = ref<number>(150000000)
const dppAmount = ref<number>(150000000)
const taxRate = ref<string>('2%')
const taxAmount = ref<number>(3000000)
const ntpnNumber = ref<string>('')
const notes = ref<string>('')
const status = ref<'DRAFT' | 'PAID' | 'REPORTED'>('REPORTED')

function updateSuggestedTaxAmount() {
  const dpp = Number(dppAmount.value) || 0
  if (taxType.value === 'PPH_23') {
    taxRate.value = '2%'
    taxAmount.value = Math.round(dpp * 0.02)
  } else if (taxType.value === 'PPH_4_2') {
    taxRate.value = '10%'
    taxAmount.value = Math.round(dpp * 0.10)
  } else if (taxType.value === 'PPH_25') {
    taxRate.value = 'Angsuran Tetap'
    taxAmount.value = 13000000
  } else if (taxType.value === 'PPN') {
    taxRate.value = '11%'
    taxAmount.value = Math.round(dpp * 0.11)
  } else if (taxType.value === 'PPH_21') {
    taxRate.value = 'PMK 168 / TER'
    taxAmount.value = Math.round(dpp * 0.05 * 0.5)
  }
}

function handleSave() {
  if (dppAmount.value <= 0 || taxAmount.value <= 0) {
    accountingStore.showToast('Nilai DPP dan Pajak Terutang harus lebih dari Rp 0', 'error')
    return
  }

  taxStore.createSptMasa({
    masaMonth: Number(masaMonth.value),
    masaYear: Number(masaYear.value),
    taxType: taxType.value,
    bankAccountCode: bankAccountCode.value,
    grossAmount: Number(grossAmount.value),
    dppAmount: Number(dppAmount.value),
    taxRate: taxRate.value,
    taxAmount: Number(taxAmount.value),
    notes: `[Unit ${unit.value}] ${notes.value || ''}`.trim(),
    ntpnNumber: ntpnNumber.value.trim(),
    status: status.value
  })

  accountingStore.showToast(`SPT Masa ${taxType.value} Unit ${unit.value} Masa ${masaMonth.value}/${masaYear.value} berhasil dibuat!`, 'success')
  router.push('/accounting/tax/spt-masa')
}
</script>

<template>
  <div class="space-y-5 font-['Plus_Jakarta_Sans',sans-serif] max-w-4xl mx-auto pb-12">
    <!-- Breadcrumb & Back Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
          <router-link to="/accounting/tax/spt-masa" class="hover:text-zinc-950 transition flex items-center gap-1">
            <ArrowLeft class="w-3.5 h-3.5" />
            <span>SPT Masa</span>
          </router-link>
          <span>/</span>
          <span class="text-zinc-950">Formulir Rekam Baru</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight">
          Perekaman Dokumen SPT Masa Baru
        </h1>
        <p class="text-xs text-zinc-500">
          Mencatat kepatuhan pajak bulanan klinik (PPh 21, PPh 23, PPh 4(2), PPh 25, & PPN) terverifikasi kanal perbankan resmi.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <router-link
          to="/accounting/tax/spt-masa"
          class="px-4 py-2 bg-white hover:bg-zinc-50 text-zinc-700 rounded-xl text-xs font-bold border border-zinc-200 transition"
        >
          Batal
        </router-link>
        <button
          @click="handleSave"
          class="inline-flex items-center gap-1.5 px-4 py-2 bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl text-xs font-bold shadow-2xs transition active:scale-95 cursor-pointer"
        >
          <Save class="w-3.5 h-3.5" />
          <span>Simpan Dokumen SPT Masa</span>
        </button>
      </div>
    </header>

    <!-- Form Container -->
    <div class="bg-white rounded-2xl border border-zinc-200/90 shadow-2xs p-6 space-y-6">
      <!-- Section 1: Periode & Klasifikasi -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Calendar class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            1. Periode Masa & Unit Bisnis
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Masa Pajak (Bulan)</label>
            <select
              v-model.number="masaMonth"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option :value="1">01 - Januari</option>
              <option :value="2">02 - Februari</option>
              <option :value="3">03 - Maret</option>
              <option :value="4">04 - April</option>
              <option :value="5">05 - Mei</option>
              <option :value="6">06 - Juni</option>
              <option :value="7">07 - Juli</option>
              <option :value="8">08 - Agustus</option>
              <option :value="9">09 - September</option>
              <option :value="10">10 - Oktober</option>
              <option :value="11">11 - November</option>
              <option :value="12">12 - Desember</option>
            </select>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Tahun Pajak</label>
            <input
              v-model.number="masaYear"
              type="number"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold font-mono text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Unit Bisnis Klinik</label>
            <select
              v-model="unit"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option value="KIDINEST">Unit Kidinest (Klinik Anak)</option>
              <option value="SSC">Unit SSC (Semarang Skin Center)</option>
              <option value="ALL">Klinik Terpadu (Konsolidasi)</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Section 2: Jenis Pajak & Perhitungan -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Calculator class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            2. Jenis Pajak & Nilai Terutang
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-xs">
          <div>
            <label class="block font-bold text-zinc-700 mb-1">Jenis SPT Masa</label>
            <select
              v-model="taxType"
              @change="updateSuggestedTaxAmount"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option value="PPH_21">PPh Pasal 21 (Tenaga Medis Dokter & Karyawan)</option>
              <option value="PPH_23">PPh Pasal 23 (Jasa Medis, Laboratorium, Limbah B3)</option>
              <option value="PPH_4_2">PPh Final Pasal 4(2) (Sewa Gedung / Bunga Giro)</option>
              <option value="PPH_25">PPh Pasal 25 (Angsuran Pajak Badan Bulanan)</option>
              <option value="PPN">PPN (Pajak Pertambahan Nilai)</option>
            </select>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Tarif Pajak Efektif / Persentase</label>
            <input
              v-model="taxRate"
              type="text"
              placeholder="Contoh: 2%, 10%, 11%, TER"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Nilai Bruto / Objek Transaksi (Rp)</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-xs">Rp</span>
              <input
                v-model.number="grossAmount"
                type="number"
                class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              />
            </div>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">DPP (Dasar Pengenaan Pajak) (Rp)</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400 font-mono text-xs">Rp</span>
              <input
                v-model.number="dppAmount"
                @input="updateSuggestedTaxAmount"
                type="number"
                class="w-full pl-9 pr-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
              />
            </div>
          </div>
        </div>

        <!-- Calculated Tax Amount Banner -->
        <div class="p-4 rounded-xl bg-zinc-900 text-white flex items-center justify-between">
          <div>
            <span class="text-[10px] font-bold text-zinc-400 uppercase tracking-wider block">Total Pajak Terutang (SPT Masa)</span>
            <span class="text-lg font-bold tabular-nums text-amber-300">
              Rp {{ formatIDR(taxAmount) }}
            </span>
          </div>
          <div class="w-48">
            <label class="block text-[10px] font-semibold text-zinc-400 mb-1">Koreksi Manual Nominal (Rp):</label>
            <input
              v-model.number="taxAmount"
              type="number"
              class="w-full px-2.5 py-1 bg-zinc-800 border border-zinc-700 rounded-lg text-xs font-mono font-bold text-white focus:outline-none"
            />
          </div>
        </div>
      </div>

      <!-- Section 3: Kanal Bank & NTPN -->
      <div class="space-y-4">
        <div class="flex items-center gap-2 border-b border-zinc-100 pb-2">
          <Landmark class="w-4 h-4 text-zinc-700" />
          <h3 class="text-xs font-bold text-zinc-900 uppercase tracking-wider">
            3. Kanal Perbankan & NTPN (Strictly Non-Cash)
          </h3>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs">
          <div class="sm:col-span-2">
            <label class="block font-bold text-zinc-700 mb-1">Rekening Bank Penyetor Resmi</label>
            <select
              v-model="bankAccountCode"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option
                v-for="b in taxStore.availableBankAccounts"
                :key="b.code"
                :value="b.code"
              >
                {{ b.code }} - {{ b.name }} (Saldo: Rp {{ formatIDR(b.currentBalance) }})
              </option>
            </select>
            <span class="text-[10px] text-zinc-400 mt-0.5 block">
              Sesuai PMK DJP: Penyetoran SPT Masa wajib melalui rekening bank operasional, bukan kas tunai.
            </span>
          </div>

          <div>
            <label class="block font-bold text-zinc-700 mb-1">Status Pelaporan</label>
            <select
              v-model="status"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-bold text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none cursor-pointer"
            >
              <option value="DRAFT">DRAFT (Belum Bayar)</option>
              <option value="PAID">PAID (Sudah Terdebit Bank)</option>
              <option value="REPORTED">REPORTED (Lengkap NTPN / BPE)</option>
            </select>
          </div>

          <div class="sm:col-span-2">
            <label class="block font-bold text-zinc-700 mb-1">Nomor NTPN / Kode Billing DJP</label>
            <input
              v-model="ntpnNumber"
              type="text"
              placeholder="Contoh: NTPN-BCA-20260328-9901 atau e-Billing DJP"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl font-mono text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>

          <div class="sm:col-span-3">
            <label class="block font-bold text-zinc-700 mb-1">Keterangan / Uraian Transaksi</label>
            <input
              v-model="notes"
              type="text"
              placeholder="Contoh: Pemotongan PPh 23 Jasa Pengelolaan Limbah Medis B3 Periode Maret"
              class="w-full px-3 py-2 bg-zinc-50 border border-zinc-200 rounded-xl text-zinc-900 focus:ring-1 focus:ring-zinc-950 focus:outline-none"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
