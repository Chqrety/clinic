<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { api, type Invoice } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Receipt, Printer, Building2, UserCheck, Percent, RefreshCw, CheckCircle2, Pill, Ticket, Eye } from 'lucide-vue-next'

const store = useClinicStore()
const invoices = ref<Invoice[]>([])
const selectedInvoice = ref<Invoice | null>(null)
const viewMode = ref<'patient' | 'internal'>('patient')
const previewThermal = ref(false)
const typeFilter = ref<'ALL' | 'REGISTRATION' | 'PHARMACY'>('ALL')
const loading = ref(false)

function formatRupiah(val?: number) {
  return 'Rp ' + (val || 0).toLocaleString('id-ID')
}

function formatDate(val?: string) {
  if (!val) return '-'
  try {
    return new Date(val.replace(' ', 'T')).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  } catch {
    return val
  }
}

const filteredInvoices = computed(() => {
  if (typeFilter.value === 'ALL') return invoices.value
  return invoices.value.filter(i => (i.invoice_type || 'REGISTRATION') === typeFilter.value)
})

async function fetchInvoices() {
  loading.value = true
  try {
    invoices.value = await api.getInvoices()
    if (store.selectedInvoiceQueueId) {
      // Prioritize pharmacy invoice if arriving from pharmacy room
      const pharmMatch = invoices.value.find(i => i.queue_id === store.selectedInvoiceQueueId && i.invoice_type === 'PHARMACY')
      const anyMatch = invoices.value.find(i => i.queue_id === store.selectedInvoiceQueueId)
      if (pharmMatch) {
        selectedInvoice.value = pharmMatch
      } else if (anyMatch) {
        selectedInvoice.value = anyMatch
      } else if (invoices.value.length > 0 && !selectedInvoice.value) {
        selectedInvoice.value = invoices.value[0]
      }
    } else if (invoices.value.length > 0 && !selectedInvoice.value) {
      selectedInvoice.value = invoices.value[0]
    }
  } catch (e) {
    // ignore
  } finally {
    loading.value = false
  }
}

watch(
  () => store.selectedInvoiceQueueId,
  (newQId) => {
    if (newQId && invoices.value.length > 0) {
      const match = invoices.value.find(i => i.queue_id === newQId && i.invoice_type === 'PHARMACY') ||
                    invoices.value.find(i => i.queue_id === newQId)
      if (match) selectedInvoice.value = match
    }
  }
)

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
          Pemisahan transaksi resmi antara Biaya Pendaftaran/Antrian dan Tebus Resep Obat Apotek.
        </p>
      </div>

      <button
        @click="fetchInvoices"
        :disabled="loading"
        class="inline-flex items-center gap-1.5 px-4 py-2 rounded-full bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50 self-start sm:self-auto"
      >
        <RefreshCw :class="['w-3.5 h-3.5', loading && 'animate-spin']" />
        Segarkan Data
      </button>
    </div>

    <!-- Active Navigation Banner from Pharmacy -->
    <div
      v-if="store.selectedInvoiceQueueId && selectedInvoice?.queue_id === store.selectedInvoiceQueueId"
      class="mb-6 p-4 rounded-2xl bg-amber-50 border border-amber-200/70 flex items-center justify-between gap-3 text-xs"
    >
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-4 h-4 text-amber-700 flex-shrink-0" />
        <span class="text-amber-950 font-bold">
          Menampilkan transaksi dari Apotek: <strong>{{ selectedInvoice.patient_name }}</strong> ({{ selectedInvoice.invoice_number }})
        </span>
      </div>
      <button
        @click="store.setSelectedInvoiceQueue(null)"
        class="text-amber-800 hover:text-amber-950 text-[11px] font-bold underline"
      >
        Tutup Tanda
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">

      <!-- Left: Invoices List with Type Filter (5 cols) -->
      <div class="lg:col-span-5 space-y-3">
        <!-- Filter Pills -->
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-zinc-400 pl-1">
            Daftar Transaksi ({{ filteredInvoices.length }})
          </span>

          <div class="inline-flex p-0.5 bg-stone-200/70 rounded-full text-[11px] font-bold">
            <button
              @click="typeFilter = 'ALL'"
              :class="['px-2.5 py-1 rounded-full transition-all', typeFilter === 'ALL' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-600']"
            >
              Semua
            </button>
            <button
              @click="typeFilter = 'REGISTRATION'"
              :class="['px-2.5 py-1 rounded-full transition-all', typeFilter === 'REGISTRATION' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-600']"
            >
              Antrian
            </button>
            <button
              @click="typeFilter = 'PHARMACY'"
              :class="['px-2.5 py-1 rounded-full transition-all', typeFilter === 'PHARMACY' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-600']"
            >
              Obat
            </button>
          </div>
        </div>

        <div class="space-y-2 max-h-[600px] overflow-y-auto pr-1">
          <div
            v-for="inv in filteredInvoices"
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
              <div class="flex items-center gap-1.5">
                <!-- Type Pill -->
                <span
                  :class="[
                    'text-[10px] font-bold px-2 py-0.5 rounded-full',
                    inv.invoice_type === 'PHARMACY' ? 'bg-emerald-100 text-emerald-800' : 'bg-blue-100 text-blue-800'
                  ]"
                >
                  {{ inv.invoice_type === 'PHARMACY' ? 'Obat Apotek' : 'Pendaftaran' }}
                </span>
                <span class="text-[10px] font-bold px-1.5 py-0.5 rounded-full bg-stone-100 text-zinc-700">
                  {{ inv.payment_method }}
                </span>
              </div>
            </div>
            <div class="font-bold text-sm">{{ inv.patient_name }}</div>
            <div :class="['text-xs mt-1 font-semibold', selectedInvoice?.id === inv.id ? 'text-amber-300' : 'text-zinc-600']">
              Total: {{ formatRupiah(inv.total_amount) }}
            </div>
          </div>

          <div v-if="filteredInvoices.length === 0" class="p-8 text-center text-zinc-400 bg-white rounded-3xl border border-black/5 text-sm">
            Belum ada transaksi pada kategori ini.
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
              1. Format Pasien (Resmi)
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
              2. Audit Internal Klinik
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              v-if="viewMode === 'patient'"
              @click="previewThermal = !previewThermal"
              class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full border border-black/10 bg-white hover:bg-stone-50 text-xs font-bold text-zinc-700 transition-all active:scale-95 shadow-xs"
              :title="previewThermal ? 'Kembali ke tampilan web modern' : 'Lihat preview struk kertas kasir 80mm'"
            >
              <Eye v-if="!previewThermal" class="w-3.5 h-3.5 text-zinc-500" />
              <Receipt v-else class="w-3.5 h-3.5 text-zinc-500" />
              {{ previewThermal ? 'Tampilan Web Modern' : 'Preview Struk Fisik' }}
            </button>

            <button
              @click="printReceipt"
              class="inline-flex items-center gap-1.5 px-4 py-1.5 rounded-full bg-zinc-950 hover:bg-zinc-800 text-xs font-bold text-white transition-all active:scale-95 shadow-sm"
            >
              <Printer class="w-3.5 h-3.5" /> Cetak Struk
            </button>
          </div>
        </div>

        <!-- Mode 1: Format Pasien (Web Modern vs Struk Tradisional Cetak) -->
        <template v-if="viewMode === 'patient'">

          <!-- 1A: Web Modern Invoice (Konsisten dengan Desain Nordic / Brandlos Seluruh Halaman Web) -->
          <div
            v-if="!previewThermal"
            id="web-invoice"
            class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] space-y-6 print:hidden"
          >
            <!-- Header Invoice Web -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-black/5">
              <div>
                <div class="flex items-center gap-2.5">
                  <h3 class="font-extrabold text-xl text-zinc-950 tracking-tight">KLINIK SEHAT MODERN</h3>
                  <span
                    :class="[
                      'text-[10px] font-bold px-2.5 py-0.5 rounded-full border',
                      selectedInvoice.invoice_type === 'PHARMACY'
                        ? 'bg-emerald-50 text-emerald-800 border-emerald-200/60'
                        : 'bg-blue-50 text-blue-800 border-blue-200/60'
                    ]"
                  >
                    {{ selectedInvoice.invoice_type === 'PHARMACY' ? 'Resep Obat Apotek' : 'Pendaftaran Loket' }}
                  </span>
                </div>
                <p class="text-xs text-zinc-500 mt-1">
                  Jl. Kesehatan No. 88, Jakarta Pusat • Telp: (021) 555-0199 • Faskes Terakreditasi
                </p>
              </div>

              <div>
                <span class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                  <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" />
                  LUNAS • {{ selectedInvoice.payment_method }}
                </span>
              </div>
            </div>

            <!-- Metadata Cards -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 p-4 rounded-2xl bg-[#F9F9F7] border border-black/5 text-xs">
              <div>
                <span class="text-zinc-400 block text-[10px] uppercase font-bold tracking-wider">No. Invoice</span>
                <span class="font-mono font-bold text-zinc-900 mt-0.5 block">{{ selectedInvoice.invoice_number }}</span>
              </div>
              <div>
                <span class="text-zinc-400 block text-[10px] uppercase font-bold tracking-wider">Tanggal</span>
                <span class="font-bold text-zinc-900 mt-0.5 block">{{ formatDate(selectedInvoice.created_at) }}</span>
              </div>
              <div>
                <span class="text-zinc-400 block text-[10px] uppercase font-bold tracking-wider">Nama Pasien</span>
                <span class="font-bold text-zinc-900 mt-0.5 block truncate">{{ selectedInvoice.patient_name }}</span>
              </div>
              <div>
                <span class="text-zinc-400 block text-[10px] uppercase font-bold tracking-wider">NIK Pasien</span>
                <span class="font-bold text-zinc-900 mt-0.5 block">{{ selectedInvoice.patient_nik }}</span>
              </div>
            </div>

            <!-- Itemized Rincian Transaksi -->
            <div>
              <div class="flex items-center justify-between mb-3 pl-1">
                <span class="text-xs font-bold uppercase tracking-wider text-zinc-400">
                  Rincian Transaksi Layanan
                </span>
                <span class="text-xs text-zinc-500 font-medium">
                  Petugas: {{ selectedInvoice.invoice_type === 'PHARMACY' ? 'apt. Annisa Fitri' : 'Budi (Kasir)' }}
                </span>
              </div>

              <!-- Tabel Obat Apotek -->
              <div v-if="selectedInvoice.invoice_type === 'PHARMACY'" class="border border-black/5 rounded-2xl overflow-hidden divide-y divide-black/5">
                <div v-if="selectedInvoice.items && selectedInvoice.items.length > 0">
                  <div
                    v-for="it in selectedInvoice.items"
                    :key="it.id"
                    class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-stone-50/50 transition-colors"
                  >
                    <div class="space-y-1">
                      <div class="font-bold text-sm text-zinc-950">{{ it.drug_name }}</div>
                      <div class="flex flex-wrap items-center gap-2 text-xs text-zinc-500">
                        <span v-if="it.dosage_instruction" class="bg-stone-100 text-zinc-700 px-2 py-0.5 rounded-md text-[11px] font-medium">
                          {{ it.dosage_instruction }}
                        </span>
                        <span v-if="it.batch_number" class="text-zinc-400 font-mono text-[11px]">
                          Batch: {{ it.batch_number }}
                        </span>
                      </div>
                    </div>
                    <div class="flex items-center justify-between sm:justify-end gap-6 text-sm">
                      <span class="text-zinc-500 text-xs font-medium">
                        {{ it.quantity }} item &times; {{ formatRupiah(it.unit_price) }}
                      </span>
                      <span class="font-bold text-zinc-900 min-w-[90px] text-right">
                        {{ formatRupiah(it.subtotal) }}
                      </span>
                    </div>
                  </div>
                </div>
                <div v-else class="p-4 flex items-center justify-between text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Resep Obat Dokter</div>
                    <div class="text-xs text-zinc-500">Paket obat resep lengkap dari apotek</div>
                  </div>
                  <div class="font-bold text-zinc-900">{{ formatRupiah(selectedInvoice.medication_fee) }}</div>
                </div>
              </div>

              <!-- Tabel Antrian & Pendaftaran -->
              <div v-else class="border border-black/5 rounded-2xl overflow-hidden divide-y divide-black/5">
                <div class="p-4 flex items-center justify-between hover:bg-stone-50/50 transition-colors text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Biaya Administrasi & Tiket Antrian</div>
                    <div class="text-xs text-zinc-500">Registrasi sistem & pencatatan rekam medis terpadu</div>
                  </div>
                  <div class="font-bold text-zinc-900">{{ formatRupiah(selectedInvoice.registration_fee) }}</div>
                </div>
                <div v-if="selectedInvoice.consultation_fee > 0" class="p-4 flex items-center justify-between hover:bg-stone-50/50 transition-colors text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Jasa Konsultasi Dokter Umum</div>
                    <div class="text-xs text-zinc-500">Pemeriksaan klinis dan edukasi dokter</div>
                  </div>
                  <div class="font-bold text-zinc-900">{{ formatRupiah(selectedInvoice.consultation_fee) }}</div>
                </div>
              </div>
            </div>

            <!-- Financial Calculation Summary -->
            <div class="pt-4 border-t border-black/5 space-y-2">
              <div class="flex justify-between text-xs text-zinc-500">
                <span>Subtotal Layanan</span>
                <span class="font-medium text-zinc-800">{{ formatRupiah(selectedInvoice.total_amount) }}</span>
              </div>
              <div class="flex justify-between text-xs text-zinc-500">
                <span>PPN 11% (Sudah Termasuk)</span>
                <span class="font-medium text-zinc-800">{{ formatRupiah(selectedInvoice.tax_share) }}</span>
              </div>
              <div class="flex justify-between items-baseline pt-4 border-t border-black/10">
                <span class="text-xs font-bold uppercase tracking-wider text-zinc-400">Total Tagihan</span>
                <span class="text-3xl font-extrabold text-zinc-950 tracking-tight">
                  {{ formatRupiah(selectedInvoice.total_amount) }}
                </span>
              </div>
            </div>

            <div class="text-center pt-2 text-xs text-zinc-400 border-t border-dashed border-zinc-200">
              Bukti transaksi sah Klinik Sehat Modern • Disimpan sebagai tanda terima resmi
            </div>
          </div>

          <!-- 1B: Traditional Physical Thermal Receipt (Digunakan saat Cetak Printer & Dapat Di-preview) -->
          <div :class="['flex justify-center py-2', !previewThermal ? 'hidden print:flex' : '']">
            <div
              id="printable-receipt"
              class="w-full max-w-[340px] bg-white text-zinc-950 font-mono text-[11px] leading-tight p-6 sm:p-7 shadow-[0_12px_40px_rgba(0,0,0,0.08)] border border-black/10 relative"
            >
              <!-- Serrated Top Edge Decoration (Visible on screen preview, hidden on physical print) -->
              <div class="no-print absolute -top-2 left-0 right-0 h-2 overflow-hidden flex">
                <svg class="w-full h-2 text-white fill-current filter drop-shadow-[0_-2px_1px_rgba(0,0,0,0.05)]" viewBox="0 0 120 8" preserveAspectRatio="none">
                  <polygon points="0,8 3,0 6,8 9,0 12,8 15,0 18,8 21,0 24,8 27,0 30,8 33,0 36,8 39,0 42,8 45,0 48,8 51,0 54,8 57,0 60,8 63,0 66,8 69,0 72,8 75,0 78,8 81,0 84,8 87,0 90,8 93,0 96,8 99,0 102,8 105,0 108,8 111,0 114,8 117,0 120,8" />
                </svg>
              </div>

              <!-- Header Struk Tradisional -->
              <div class="text-center space-y-0.5 pb-2">
                <div class="font-black text-xs sm:text-sm tracking-tight uppercase">KLINIK SEHAT MODERN</div>
                <div class="text-[10px] text-zinc-600">Jl. Kesehatan No. 88, Jakarta Pusat</div>
                <div class="text-[10px] text-zinc-600">Telp: (021) 555-0199 | WA: 0812-8877-6655</div>
                <div class="text-[9px] text-zinc-400">Izin Faskes: 445/0912/DINKES/2024</div>
              </div>

              <!-- Dashed line -->
              <div class="border-b border-dashed border-zinc-400 my-2"></div>

              <!-- Meta Data Struk -->
              <div class="space-y-1 text-[10px]">
                <div class="flex justify-between">
                  <span>NO. STRUK</span>
                  <span class="font-bold font-mono">{{ selectedInvoice.invoice_number }}</span>
                </div>
                <div class="flex justify-between">
                  <span>TANGGAL</span>
                  <span>{{ formatDate(selectedInvoice.created_at) }}</span>
                </div>
                <div class="flex justify-between">
                  <span>PETUGAS</span>
                  <span>{{ selectedInvoice.invoice_type === 'PHARMACY' ? 'apt. Annisa Fitri' : 'Budi (Kasir)' }}</span>
                </div>
                <div class="flex justify-between">
                  <span>PASIEN</span>
                  <span class="font-bold uppercase truncate max-w-[170px] text-right">{{ selectedInvoice.patient_name }}</span>
                </div>
                <div class="flex justify-between">
                  <span>NIK</span>
                  <span>{{ selectedInvoice.patient_nik }}</span>
                </div>
                <div class="flex justify-between">
                  <span>LAYANAN</span>
                  <span class="font-bold">{{ selectedInvoice.invoice_type === 'PHARMACY' ? 'RESEP APOTEK' : 'PENDAFTARAN LOKET' }}</span>
                </div>
              </div>

              <!-- Dashed line -->
              <div class="border-b border-dashed border-zinc-400 my-2.5"></div>

              <!-- Daftar Item Barang / Jasa (Layout 2 Baris Klasik) -->
              <div class="space-y-2.5 text-[10px]">
                <!-- Rincian Obat Apotek -->
                <template v-if="selectedInvoice.invoice_type === 'PHARMACY'">
                  <div v-if="selectedInvoice.items && selectedInvoice.items.length > 0" class="space-y-2">
                    <div v-for="it in selectedInvoice.items" :key="it.id" class="space-y-0.5">
                      <div class="font-bold uppercase leading-tight">{{ it.drug_name }}</div>
                      <div class="flex justify-between text-zinc-700">
                        <span>{{ it.quantity }} x {{ formatRupiah(it.unit_price) }}</span>
                        <span class="font-bold">{{ formatRupiah(it.subtotal) }}</span>
                      </div>
                      <div class="text-[9px] text-zinc-500 italic">
                        &bull; Batch: {{ it.batch_number }} | {{ it.dosage_instruction }}
                      </div>
                    </div>
                  </div>
                  <div v-else class="space-y-0.5">
                    <div class="font-bold uppercase">OBAT RESEP DOKTER</div>
                    <div class="flex justify-between text-zinc-700">
                      <span>1 x {{ formatRupiah(selectedInvoice.medication_fee) }}</span>
                      <span class="font-bold">{{ formatRupiah(selectedInvoice.medication_fee) }}</span>
                    </div>
                  </div>
                </template>

                <!-- Rincian Pendaftaran & Antrian -->
                <template v-else>
                  <div class="space-y-0.5">
                    <div class="font-bold uppercase">BIAYA REGISTRASI & ANTRIAN</div>
                    <div class="flex justify-between text-zinc-700">
                      <span>1 x {{ formatRupiah(selectedInvoice.registration_fee) }}</span>
                      <span class="font-bold">{{ formatRupiah(selectedInvoice.registration_fee) }}</span>
                    </div>
                  </div>

                  <div v-if="selectedInvoice.consultation_fee > 0" class="space-y-0.5 pt-1">
                    <div class="font-bold uppercase">KONSULTASI DOKTER UMUM</div>
                    <div class="flex justify-between text-zinc-700">
                      <span>1 x {{ formatRupiah(selectedInvoice.consultation_fee) }}</span>
                      <span class="font-bold">{{ formatRupiah(selectedInvoice.consultation_fee) }}</span>
                    </div>
                  </div>
                </template>
              </div>

              <!-- Double dashed line separator -->
              <div class="border-b-2 border-dashed border-zinc-800 my-3"></div>

              <!-- Kalkulasi Biaya -->
              <div class="space-y-1 text-[10px]">
                <div class="flex justify-between">
                  <span>SUBTOTAL</span>
                  <span>{{ formatRupiah(selectedInvoice.total_amount) }}</span>
                </div>
                <div class="flex justify-between text-zinc-600">
                  <span>PPN 11% (Termasuk)</span>
                  <span>{{ formatRupiah(selectedInvoice.tax_share) }}</span>
                </div>
                <div class="border-b border-dotted border-zinc-400 my-1"></div>
                <div class="flex justify-between text-xs font-black pt-0.5">
                  <span>TOTAL AKHIR</span>
                  <span>{{ formatRupiah(selectedInvoice.total_amount) }}</span>
                </div>
                <div class="flex justify-between pt-0.5">
                  <span>BAYAR ({{ selectedInvoice.payment_method }})</span>
                  <span>{{ formatRupiah(selectedInvoice.total_amount) }}</span>
                </div>
                <div class="flex justify-between text-zinc-600">
                  <span>KEMBALIAN</span>
                  <span>Rp 0</span>
                </div>
              </div>

              <!-- Double dashed line separator -->
              <div class="border-b border-dashed border-zinc-400 my-3"></div>

              <!-- Traditional LUNAS Box Stamp -->
              <div class="text-center py-1">
                <div class="inline-block border border-zinc-900 px-3.5 py-0.5 text-[11px] font-black tracking-widest uppercase">
                  *** L U N A S ***
                </div>
              </div>

              <!-- Traditional Simulated Barcode -->
              <div class="pt-2.5 pb-1 text-center space-y-0.5 select-none">
                <div class="text-xl tracking-widest font-mono opacity-80 leading-none scale-y-125">
                  ||||| | |||| || ||| |||| | ||||| ||| ||| | ||
                </div>
                <div class="text-[9px] tracking-widest text-zinc-500 font-mono">
                  {{ selectedInvoice.invoice_number.replace(/-/g, '') }}
                </div>
              </div>

              <!-- Traditional Footer Message -->
              <div class="text-center text-[9px] text-zinc-500 space-y-0.5 pt-1 border-t border-dashed border-zinc-200 mt-2">
                <div class="font-bold text-zinc-800">SEMOGA LEKAS SEMBUH</div>
                <div>Terima kasih atas kunjungan Anda</div>
                <div>Simpan struk ini sebagai bukti pembayaran resmi</div>
                <div class="text-[8px] text-zinc-400 pt-0.5">POS Klinik Sehat Modern v2.5</div>
              </div>

              <!-- Serrated Bottom Edge Decoration (Visible on screen preview, hidden on physical print) -->
              <div class="no-print absolute -bottom-2 left-0 right-0 h-2 overflow-hidden flex">
                <svg class="w-full h-2 text-white fill-current filter drop-shadow-[0_2px_1px_rgba(0,0,0,0.05)]" viewBox="0 0 120 8" preserveAspectRatio="none">
                  <polygon points="0,0 3,8 6,0 9,8 12,0 15,8 18,0 21,8 24,0 27,8 30,0 33,8 36,0 39,8 42,0 45,8 48,0 51,8 54,0 57,8 60,0 63,8 66,0 69,8 72,0 75,8 78,0 81,8 84,0 87,8 90,0 93,8 96,0 99,8 102,0 105,8 108,0 111,8 114,0 117,8 120,0" />
                </svg>
              </div>
            </div>
          </div>

        </template>

        <!-- Mode 2: Internal Clinic Breakdown Rules -->
        <div
          v-else
          class="bg-white rounded-3xl p-8 border border-black/5 shadow-[0_8px_30px_rgba(0,0,0,0.06)] space-y-6"
        >
          <div>
            <span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold bg-purple-100 text-purple-900 mb-2">
              Laporan Akuntansi Internal
            </span>
            <h3 class="text-xl font-extrabold text-zinc-950">
              Alokasi Pendapatan: {{ selectedInvoice.invoice_type === 'PHARMACY' ? 'Pos Farmasi Apotek' : 'Pos Layanan & Pendaftaran' }}
            </h3>
            <p class="text-xs text-zinc-500">
              Rincian ini memisahkan margin bersih klinik, hak dokter, dan setoran pajak resmi.
            </p>
          </div>

          <!-- Cards breakdown -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="p-4 rounded-2xl bg-blue-50/60 border border-blue-200">
              <div class="flex items-center justify-between text-xs font-bold text-blue-800">
                <span>{{ selectedInvoice.invoice_type === 'PHARMACY' ? 'Jasa Racik' : 'Jasa Medis (70%)' }}</span>
                <UserCheck class="w-4 h-4 text-blue-600" />
              </div>
              <div class="text-xl font-extrabold text-blue-950 mt-2">
                {{ formatRupiah(selectedInvoice.doctor_share) }}
              </div>
              <div class="text-[10px] text-blue-700 mt-1">
                {{ selectedInvoice.invoice_type === 'PHARMACY' ? 'Apotek Mandiri' : 'Hak Dokter Jasa Medis' }}
              </div>
            </div>

            <div class="p-4 rounded-2xl bg-emerald-50/60 border border-emerald-200">
              <div class="flex items-center justify-between text-xs font-bold text-emerald-800">
                <span>Pendapatan Bersih Klinik</span>
                <Building2 class="w-4 h-4 text-emerald-600" />
              </div>
              <div class="text-xl font-extrabold text-emerald-950 mt-2">
                {{ formatRupiah(selectedInvoice.clinic_share) }}
              </div>
              <div class="text-[10px] text-emerald-700 mt-1">Operasional & sarana</div>
            </div>

            <div class="p-4 rounded-2xl bg-amber-50/60 border border-amber-200">
              <div class="flex items-center justify-between text-xs font-bold text-amber-800">
                <span>Pajak (PPN 11%)</span>
                <Percent class="w-4 h-4 text-amber-600" />
              </div>
              <div class="text-xl font-extrabold text-amber-950 mt-2">
                {{ formatRupiah(selectedInvoice.tax_share) }}
              </div>
              <div class="text-[10px] text-amber-700 mt-1">Setoran perpajakan resmi</div>
            </div>
          </div>

          <!-- Rules explanation -->
          <div class="p-4 rounded-2xl bg-[#F9F9F7] border border-black/5 text-xs text-zinc-600 space-y-1.5">
            <div class="font-bold text-zinc-900">Catatan Kebijakan Akuntansi:</div>
            <div v-if="selectedInvoice.invoice_type === 'PHARMACY'">
              • Transaksi Farmasi merupakan pos penjualan obat mandiri terpisah dari registrasi antrian.
              <br>• Pendapatan obat masuk ke margin sarana klinik setelah dipotong PPN 11%.
            </div>
            <div v-else>
              • Biaya Pendaftaran Rp 35.000 dibayar upfront di awal antrian untuk administrasi klinik.
              <br>• Pasien memegang struk antrian dan struk obat sebagai 2 bukti transaksi terpisah.
            </div>
          </div>
        </div>

      </div>

    </div>
  </div>
</template>

<style>
@media print {
  /* Sembunyikan seluruh UI luar dan navigasi aplikasi saat cetak */
  body {
    background: #ffffff !important;
  }
  body * {
    visibility: hidden;
  }
  #printable-receipt,
  #printable-receipt * {
    visibility: visible !important;
  }
  #printable-receipt {
    display: block !important;
    position: absolute;
    left: 50%;
    top: 0;
    transform: translateX(-50%);
    width: 80mm !important;
    max-width: 80mm !important;
    margin: 0 auto !important;
    padding: 16px 12px !important;
    box-shadow: none !important;
    border: none !important;
  }
  #web-invoice {
    display: none !important;
  }
  .no-print {
    display: none !important;
  }
}
</style>
