<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type Invoice } from '../api/client'
import { Receipt, Printer, Eye, Building2, UserCheck, Percent } from 'lucide-vue-next'

const invoices = ref<Invoice[]>([])
const selectedInvoice = ref<Invoice | null>(null)
const viewMode = ref<'patient' | 'internal'>('patient')
const loading = ref(false)

async function fetchInvoices() {
  loading.value = true
  try {
    invoices.value = await api.getInvoices()
    if (invoices.value.length > 0 && !selectedInvoice.value) {
      selectedInvoice.value = invoices.value[0]
    }
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

function printReceipt() {
  window.print()
}

onMounted(() => {
  fetchInvoices()
})
</script>

<template>
  <div class="py-8 px-6 sm:px-10 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-4 mb-6 pb-6 border-b border-black/10">
      <div>
        <h1 class="text-3xl font-extrabold text-zinc-950 tracking-tight">
          Billing & Sistem Invoice Terpadu
        </h1>
        <p class="text-xs text-zinc-500 mt-1">
          Pasien menerima 1 invoice sederhana. Klinik memegang audit internal (Jasa Medis, Margin Klinik & Pajak).
        </p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

      <!-- Left: Invoices List (5 cols) -->
      <div class="lg:col-span-5 space-y-3">
        <h2 class="text-xs font-bold uppercase tracking-wider text-zinc-400 pl-1">
          Daftar Transaksi Selesai ({{ invoices.length }})
        </h2>

        <div class="space-y-2 max-h-[600px] overflow-y-auto pr-1">
          <div
            v-for="inv in invoices"
            :key="inv.id"
            @click="selectedInvoice = inv"
            :class="[
              'p-4 rounded-3xl border transition-all cursor-pointer active:scale-[0.99]',
              selectedInvoice?.id === inv.id
                ? 'bg-zinc-950 text-white border-zinc-950 shadow-md'
                : 'bg-white text-zinc-900 border-black/5 hover:border-black/15 shadow-[0_2px_12px_rgba(0,0,0,0.02)]'
            ]"
          >
            <div class="flex items-center justify-between mb-1">
              <span class="font-mono text-xs font-bold">{{ inv.invoice_number }}</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800">
                LUNAS ({{ inv.payment_method }})
              </span>
            </div>
            <div class="font-bold text-sm">{{ inv.patient_name }}</div>
            <div :class="['text-xs mt-1 font-semibold', selectedInvoice?.id === inv.id ? 'text-amber-300' : 'text-zinc-600']">
              Total: Rp {{ inv.total_amount.toLocaleString('id-ID') }}
            </div>
          </div>

          <div v-if="invoices.length === 0" class="p-8 text-center text-zinc-400 bg-white rounded-3xl border border-black/5 text-sm">
            Belum ada transaksi invoice selesai.
          </div>
        </div>
      </div>

      <!-- Right: Invoice Preview (7 cols) -->
      <div v-if="selectedInvoice" class="lg:col-span-7">

        <!-- Mode Switcher Pill Tabs -->
        <div class="flex items-center justify-between mb-4">
          <div class="inline-flex p-1 bg-white rounded-full border border-black/10 shadow-sm">
            <button
              @click="viewMode = 'patient'"
              :class="[
                'px-4 py-1.5 rounded-full text-xs font-bold transition-all active:scale-95',
                viewMode === 'patient'
                  ? 'bg-zinc-950 text-white shadow-sm'
                  : 'text-zinc-500 hover:text-zinc-900'
              ]"
            >
              1. Versi Pasien (Single Invoice)
            </button>
            <button
              @click="viewMode = 'internal'"
              :class="[
                'px-4 py-1.5 rounded-full text-xs font-bold transition-all active:scale-95',
                viewMode === 'internal'
                  ? 'bg-zinc-950 text-white shadow-sm'
                  : 'text-zinc-500 hover:text-zinc-900'
              ]"
            >
              2. Audit Internal (Aturan Klinik)
            </button>
          </div>

          <button
            @click="printReceipt"
            class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full bg-stone-100 hover:bg-stone-200 text-xs font-bold text-zinc-800 transition-all active:scale-95"
          >
            <Printer class="w-3.5 h-3.5" /> Cetak Struk
          </button>
        </div>

        <!-- Mode 1: Patient-facing Single Receipt -->
        <div
          v-if="viewMode === 'patient'"
          class="bg-white rounded-3xl p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] relative overflow-hidden space-y-6"
        >
          <div class="text-center pb-6 border-b border-dashed border-zinc-200">
            <h2 class="text-xl font-extrabold text-zinc-950">KLINIK SEHAT MODERN</h2>
            <p class="text-xs text-zinc-500 mt-0.5">Jl. Kesehatan No. 88 — Telepon: (021) 555-0199</p>
            <div class="mt-3 inline-block px-3 py-1 rounded-full bg-emerald-50 text-emerald-700 text-xs font-bold border border-emerald-200/60">
              BUKTI PEMBAYARAN LUNAS
            </div>
          </div>

          <!-- Metadata -->
          <div class="grid grid-cols-2 text-xs gap-y-2 text-zinc-600">
            <div>No. Invoice: <strong class="text-zinc-900 font-mono">{{ selectedInvoice.invoice_number }}</strong></div>
            <div class="text-right">Tanggal: {{ new Date(selectedInvoice.created_at).toLocaleDateString('id-ID') }}</div>
            <div>Pasien: <strong class="text-zinc-900">{{ selectedInvoice.patient_name }}</strong></div>
            <div class="text-right">NIK: {{ selectedInvoice.patient_nik }}</div>
          </div>

          <!-- Line items -->
          <div class="border-t border-b border-zinc-100 py-4 space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-zinc-600">Biaya Administrasi & Pendaftaran</span>
              <span class="font-bold text-zinc-900">Rp {{ selectedInvoice.registration_fee.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-600">Biaya Konsultasi Dokter Umum</span>
              <span class="font-bold text-zinc-900">Rp {{ selectedInvoice.consultation_fee.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-600">Total Biaya Resep & Farmasi</span>
              <span class="font-bold text-zinc-900">Rp {{ selectedInvoice.medication_fee.toLocaleString('id-ID') }}</span>
            </div>
          </div>

          <!-- Total -->
          <div class="flex justify-between items-baseline pt-2">
            <span class="text-sm font-bold uppercase tracking-wider text-zinc-500">Total Tagihan</span>
            <span class="text-3xl font-extrabold text-zinc-950">
              Rp {{ selectedInvoice.total_amount.toLocaleString('id-ID') }}
            </span>
          </div>

          <div class="text-center pt-4 text-[11px] text-zinc-400">
            Terima kasih atas kunjungan Anda. Semoga lekas sembuh.
          </div>
        </div>

        <!-- Mode 2: Internal Clinic Breakdown Rules -->
        <div
          v-else
          class="bg-white rounded-3xl p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] space-y-6"
        >
          <div>
            <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-900 mb-2">
              Laporan Manajemen Internal
            </span>
            <h3 class="text-xl font-extrabold text-zinc-950">
              Pemisahan Pendapatan & Aturan Klinik
            </h3>
            <p class="text-xs text-zinc-500">
              Rincian ini hanya dapat diakses staf berwenang / akuntansi klinik.
            </p>
          </div>

          <!-- Cards breakdown -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="p-4 rounded-2xl bg-blue-50/60 border border-blue-200">
              <div class="flex items-center justify-between text-xs font-bold text-blue-800">
                <span>Jasa Dokter (70%)</span>
                <UserCheck class="w-4 h-4 text-blue-600" />
              </div>
              <div class="text-xl font-extrabold text-blue-950 mt-2">
                Rp {{ selectedInvoice.doctor_share.toLocaleString('id-ID') }}
              </div>
              <div class="text-[10px] text-blue-700 mt-1">Hak bagi hasil medis</div>
            </div>

            <div class="p-4 rounded-2xl bg-emerald-50/60 border border-emerald-200">
              <div class="flex items-center justify-between text-xs font-bold text-emerald-800">
                <span>Pendapatan Klinik</span>
                <Building2 class="w-4 h-4 text-emerald-600" />
              </div>
              <div class="text-xl font-extrabold text-emerald-950 mt-2">
                Rp {{ selectedInvoice.clinic_share.toLocaleString('id-ID') }}
              </div>
              <div class="text-[10px] text-emerald-700 mt-1">Operasional & margin obat</div>
            </div>

            <div class="p-4 rounded-2xl bg-amber-50/60 border border-amber-200">
              <div class="flex items-center justify-between text-xs font-bold text-amber-800">
                <span>Pajak (PPN 11%)</span>
                <Percent class="w-4 h-4 text-amber-600" />
              </div>
              <div class="text-xl font-extrabold text-amber-950 mt-2">
                Rp {{ selectedInvoice.tax_share.toLocaleString('id-ID') }}
              </div>
              <div class="text-[10px] text-amber-700 mt-1">Setoran perpajakan resmi</div>
            </div>
          </div>

          <!-- Rules explanation -->
          <div class="p-4 rounded-2xl bg-[#F9F9F7] border border-black/5 text-xs text-zinc-600 space-y-1.5">
            <div class="font-bold text-zinc-900">Catatan Kebijakan Klinik:</div>
            <div>• Jasa dokter dihitung 70% dari nilai konsultasi tanpa potongan pendaftaran.</div>
            <div>• Margin obat dan biaya administrasi 100% masuk ke pos sarana klinik.</div>
            <div>• Pasien hanya menerima lembar tagihan konsolidasi tanpa rincian bagi hasil.</div>
          </div>
        </div>

      </div>

    </div>
  </div>
</template>
