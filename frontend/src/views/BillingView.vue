<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { api, type Invoice } from '../api/client'
import { useClinicStore } from '../stores/clinic'
import { Receipt, Printer, Building2, UserCheck, Percent, RefreshCw, CheckCircle2, Pill, Ticket, Eye, CreditCard, QrCode, Banknote } from 'lucide-vue-next'

const store = useClinicStore()
const invoices = ref<Invoice[]>([])
const selectedInvoice = ref<Invoice | null>(null)
const viewMode = ref<'patient' | 'internal'>('patient')
const previewThermal = ref(false)
const paymentFilter = ref<'ALL' | 'CASH' | 'EDC' | 'QRIS' | 'BANK_TRANSFER'>('ALL')
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
  if (paymentFilter.value === 'ALL') return invoices.value
  return invoices.value.filter(i => i.payment_method === paymentFilter.value)
})

const totalAmountSum = computed(() => invoices.value.reduce((acc, i) => acc + (i.total_amount || 0), 0))
const cashCount = computed(() => invoices.value.filter(i => i.payment_method === 'CASH').length)
const edcCount = computed(() => invoices.value.filter(i => i.payment_method === 'EDC').length)
const qrisCount = computed(() => invoices.value.filter(i => i.payment_method === 'QRIS').length)
const transferCount = computed(() => invoices.value.filter(i => i.payment_method === 'BANK_TRANSFER').length)

async function fetchInvoices() {
  loading.value = true
  try {
    invoices.value = await api.getInvoices()
    if (store.selectedInvoiceQueueId) {
      const anyMatch = invoices.value.find(i => i.queue_id === store.selectedInvoiceQueueId)
      if (anyMatch) {
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
      const match = invoices.value.find(i => i.queue_id === newQId)
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
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- TOP BAR: MINIMALIST CLINICAL HEADER -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Billing Kasir & Rekap Keuangan
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Rekap Faktur & Pembagian Hasil
        </h1>
      </div>

      <div class="flex items-center gap-2.5">
        <button
          v-if="selectedInvoice"
          type="button"
          @click="printReceipt"
          title="Cetak struk fisik transaksi kasir"
          class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white text-xs font-bold transition-all active:scale-95 shadow-2xs cursor-pointer"
        >
          <Printer class="w-3.5 h-3.5" />
          <span>Cetak Struk</span>
        </button>

        <button
          type="button"
          @click="fetchInvoices"
          :disabled="loading"
          title="Segarkan daftar invoice transaksi"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', loading && 'animate-spin']" />
          <span>Sync</span>
        </button>
      </div>
    </header>

    <!-- 1-ROW COMPACT OPERATIONAL TICKER STRIP (Hairline Grid) -->
    <div
      class="grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Total Faktur Kasir</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ invoices.length }} Faktur
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Tunai (CASH)</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ cashCount }} Transaksi
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Mesin EDC Bank</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ edcCount }} Transaksi
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">QRIS Dinamis / Statis</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-700 mt-0.5">
          {{ qrisCount }} Transaksi
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex items-center justify-between">
        <div>
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Total Penerimaan</div>
          <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
            {{ formatRupiah(totalAmountSum) }}
          </div>
        </div>
        <Receipt class="w-4 h-4 text-zinc-400" />
      </div>
    </div>

    <!-- Active Navigation Banner from Cashier / Pharmacy -->
    <div
      v-if="store.selectedInvoiceQueueId && selectedInvoice?.queue_id === store.selectedInvoiceQueueId"
      class="p-3.5 rounded-xl border border-amber-200 bg-amber-50/80 flex items-center justify-between gap-3 text-xs"
    >
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-4 h-4 text-amber-700 shrink-0" />
        <span class="text-amber-950 font-bold">
          Menampilkan transaksi antrian terpilih: <strong>{{ selectedInvoice.patient_name }}</strong> ({{ selectedInvoice.invoice_number }})
        </span>
      </div>
      <button
        type="button"
        @click="store.selectedInvoiceQueueId = null"
        class="text-amber-900 hover:text-amber-950 font-semibold underline cursor-pointer"
      >
        Tutup
      </button>
    </div>

    <!-- MAIN MASTER-DETAIL WORKSPACE (2 PANES) -->
    <main class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">

      <!-- LEFT PANE: INVOICES LIST (5 cols) -->
      <section class="lg:col-span-5 space-y-3">
        <!-- Filter Tabs (Kanal Pembayaran) -->
        <div class="flex items-center justify-between gap-2">
          <div class="flex flex-wrap items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              type="button"
              @click="paymentFilter = 'ALL'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                paymentFilter === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Semua ({{ invoices.length }})
            </button>
            <button
              type="button"
              @click="paymentFilter = 'CASH'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                paymentFilter === 'CASH' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Tunai ({{ cashCount }})
            </button>
            <button
              type="button"
              @click="paymentFilter = 'EDC'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                paymentFilter === 'EDC' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              EDC ({{ edcCount }})
            </button>
            <button
              type="button"
              @click="paymentFilter = 'QRIS'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                paymentFilter === 'QRIS' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              QRIS ({{ qrisCount }})
            </button>
            <button
              type="button"
              @click="paymentFilter = 'BANK_TRANSFER'"
              :class="[
                'px-2.5 py-1 rounded-lg transition-all cursor-pointer',
                paymentFilter === 'BANK_TRANSFER' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              Transfer ({{ transferCount }})
            </button>
          </div>

          <span class="text-[11px] font-medium text-zinc-400">
            {{ filteredInvoices.length }} Faktur
          </span>
        </div>

        <!-- Invoices List Container -->
        <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden">
          <div class="p-3.5 border-b border-zinc-100 bg-zinc-50/60 flex items-center justify-between">
            <span class="text-xs font-bold text-zinc-900">Daftar Transaksi Kasir</span>
            <span class="text-[10px] text-zinc-400 font-medium">Single Invoice</span>
          </div>

          <div v-if="loading" class="py-12 text-center text-xs text-zinc-400">
            <RefreshCw class="w-4 h-4 mx-auto animate-spin mb-2 text-zinc-400" />
            Memuat daftar invoice...
          </div>

          <div v-else-if="filteredInvoices.length === 0" class="py-12 text-center text-xs text-zinc-400">
            <CheckCircle2 class="w-6 h-6 mx-auto text-emerald-500 mb-1.5 opacity-80" />
            Belum ada transaksi pada kanal ini
          </div>

          <div v-else class="divide-y divide-zinc-100 max-h-[640px] overflow-y-auto">
            <div
              v-for="inv in filteredInvoices"
              :key="inv.id"
              @click="selectedInvoice = inv"
              :class="[
                'p-3.5 transition-colors cursor-pointer space-y-1.5',
                selectedInvoice?.id === inv.id
                  ? 'bg-zinc-900 text-white'
                  : 'hover:bg-zinc-50/70 text-zinc-900',
              ]"
            >
              <!-- Row Line 1 -->
              <div class="flex items-center justify-between gap-2">
                <span class="font-mono text-xs font-bold">{{ inv.invoice_number }}</span>
                <div class="flex items-center gap-1.5">
                  <span
                    :class="[
                      'text-[10px] font-bold px-1.5 py-0.5 rounded border',
                      selectedInvoice?.id === inv.id
                        ? 'bg-zinc-800 text-zinc-200 border-zinc-700'
                        : inv.payment_method === 'EDC'
                          ? 'bg-blue-50 text-blue-900 border-blue-200'
                          : inv.payment_method === 'QRIS'
                            ? 'bg-amber-50 text-amber-900 border-amber-200'
                            : inv.payment_method === 'CASH'
                              ? 'bg-emerald-50 text-emerald-900 border-emerald-200'
                              : 'bg-zinc-100 text-zinc-800 border-zinc-200',
                    ]"
                  >
                    {{ inv.payment_method }}
                  </span>
                  <span
                    v-if="inv.payment_reference"
                    :class="[
                      'text-[9px] font-mono px-1.5 py-0.5 rounded border',
                      selectedInvoice?.id === inv.id
                        ? 'bg-zinc-800 text-amber-300 border-zinc-700'
                        : 'bg-stone-100 text-stone-700 border-stone-200'
                    ]"
                    :title="'Trace / Ref: ' + inv.payment_reference"
                  >
                    #{{ inv.payment_reference }}
                  </span>
                </div>
              </div>

              <!-- Row Line 2 -->
              <div>
                <div class="font-bold text-sm leading-snug">{{ inv.patient_name }}</div>
                <div
                  :class="[
                    'text-xs mt-0.5 flex items-center justify-between',
                    selectedInvoice?.id === inv.id ? 'text-zinc-400' : 'text-zinc-500'
                  ]"
                >
                  <span>{{ formatDate(inv.created_at) }}</span>
                  <span class="text-sm font-extrabold tabular-nums tracking-tight" :class="selectedInvoice?.id === inv.id ? 'text-amber-300' : 'text-zinc-950'">
                    {{ formatRupiah(inv.total_amount) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- RIGHT PANE: INVOICE PREVIEW (7 cols, Sticky) -->
      <section v-if="selectedInvoice" class="lg:col-span-7 sticky top-6 space-y-3">
        
        <!-- Mode Switcher Pill Tabs -->
        <div class="flex items-center justify-between">
          <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              type="button"
              @click="viewMode = 'patient'"
              :class="[
                'px-3 py-1 rounded-lg transition-all cursor-pointer',
                viewMode === 'patient' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              1. Lembar Faktur Pasien
            </button>
            <button
              type="button"
              @click="viewMode = 'internal'"
              :class="[
                'px-3 py-1 rounded-lg transition-all cursor-pointer',
                viewMode === 'internal' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-600 hover:text-zinc-950'
              ]"
            >
              2. Audit Internal Klinik
            </button>
          </div>

          <button
            v-if="viewMode === 'patient'"
            type="button"
            @click="previewThermal = !previewThermal"
            class="inline-flex items-center gap-1.5 px-3 py-1 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all shadow-2xs cursor-pointer"
          >
            <Eye v-if="!previewThermal" class="w-3.5 h-3.5 text-zinc-500" />
            <Receipt v-else class="w-3.5 h-3.5 text-zinc-500" />
            <span>{{ previewThermal ? 'Format Modern' : 'Preview Struk Fisik 80mm' }}</span>
          </button>
        </div>

        <!-- Mode 1: Format Pasien (Web Modern vs Struk Tradisional Cetak) -->
        <template v-if="viewMode === 'patient'">

          <!-- 1A: Web Modern Invoice (Konsisten dengan Desain Nordic / Brandlos Seluruh Halaman Web) -->
          <div
            v-if="!previewThermal"
            id="web-invoice"
            class="bg-white rounded-2xl p-6 sm:p-7 border border-zinc-200/90 shadow-2xs space-y-6 print:hidden"
          >
            <!-- Header Invoice Web -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-zinc-200/80">
              <div>
                <div class="flex items-center gap-2.5">
                  <h3 class="font-extrabold text-lg text-zinc-950 tracking-tight">SEMARANG MEDICAL GROUP</h3>
                  <span class="text-[10px] font-bold px-2.5 py-0.5 rounded-full border bg-stone-50 text-stone-800 border-stone-200">
                    Layanan Terpadu
                  </span>
                </div>
                <p class="text-xs text-zinc-500 mt-1">
                  Jl. Pandanaran No. 45, Semarang • Telp: (024) 841-2345 • Faskes Terakreditasi Paripurna
                </p>
              </div>

              <div class="flex flex-col sm:items-end gap-1.5">
                <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-800 border border-emerald-200/60">
                  <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" />
                  LUNAS • {{ selectedInvoice.payment_method }}
                </span>
                <span v-if="selectedInvoice.payment_reference" class="text-[11px] font-mono text-zinc-500">
                  No. Trace / Ref: <strong class="text-zinc-900">{{ selectedInvoice.payment_reference }}</strong>
                </span>
              </div>
            </div>

            <!-- Metadata Cards -->
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 p-4 rounded-xl bg-zinc-50/80 border border-zinc-200/80 text-xs">
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
                <span class="font-mono font-bold text-zinc-900 mt-0.5 block">{{ selectedInvoice.patient_nik }}</span>
              </div>
            </div>

            <!-- Itemized Rincian Transaksi -->
            <div>
              <div class="flex items-center justify-between mb-3 pl-1">
                <span class="text-xs font-bold uppercase tracking-wider text-zinc-400">
                  Rincian Transaksi Layanan
                </span>
                <span class="text-xs text-zinc-500 font-medium">
                  Kasir POS Semarang Medical Group
                </span>
              </div>

              <!-- Tabel Itemized Single Invoice -->
              <div class="border border-zinc-200/90 rounded-xl overflow-hidden divide-y divide-zinc-200/70">
                <!-- Baris Pendaftaran -->
                <div v-if="selectedInvoice.registration_fee > 0" class="p-4 flex items-center justify-between hover:bg-zinc-50/60 transition-colors text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Biaya Administrasi & Antrian</div>
                    <div class="text-xs text-zinc-500">Registrasi sistem & pencatatan rekam medis terpadu</div>
                  </div>
                  <div class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.registration_fee) }}</div>
                </div>

                <!-- Baris Jasa Dokter -->
                <div v-if="selectedInvoice.consultation_fee > 0" class="p-4 flex items-center justify-between hover:bg-zinc-50/60 transition-colors text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Jasa Konsultasi & Tindakan Dokter</div>
                    <div class="text-xs text-zinc-500">Pemeriksaan klinis, edukasi, dan tindakan medis</div>
                  </div>
                  <div class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.consultation_fee) }}</div>
                </div>

                <!-- Baris Obat / Resep Apotek -->
                <template v-if="selectedInvoice.items && selectedInvoice.items.length > 0">
                  <div
                    v-for="it in selectedInvoice.items"
                    :key="it.id"
                    class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3 hover:bg-zinc-50/60 transition-colors"
                  >
                    <div class="space-y-1">
                      <div class="font-bold text-sm text-zinc-950">{{ it.drug_name }}</div>
                      <div class="flex flex-wrap items-center gap-2 text-xs text-zinc-500">
                        <span v-if="it.dosage_instruction" class="bg-zinc-100 text-zinc-700 px-2 py-0.5 rounded-md text-[11px] font-medium">
                          {{ it.dosage_instruction }}
                        </span>
                        <span v-if="it.batch_number" class="text-zinc-400 font-mono text-[11px]">
                          Batch: {{ it.batch_number }}
                        </span>
                      </div>
                    </div>
                    <div class="flex items-center justify-between sm:justify-end gap-6 text-sm">
                      <span class="text-zinc-500 text-xs font-medium tabular-nums">
                        {{ it.quantity }} item &times; {{ formatRupiah(it.unit_price) }}
                      </span>
                      <span class="font-bold tabular-nums text-zinc-900 min-w-[90px] text-right">
                        {{ formatRupiah(it.subtotal) }}
                      </span>
                    </div>
                  </div>
                </template>
                <div v-else-if="selectedInvoice.medication_fee > 0" class="p-4 flex items-center justify-between text-sm">
                  <div>
                    <div class="font-bold text-zinc-900">Resep Obat Apotek</div>
                    <div class="text-xs text-zinc-500">Paket obat resep lengkap dari apotek</div>
                  </div>
                  <div class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.medication_fee) }}</div>
                </div>
              </div>
            </div>

            <!-- Financial Calculation Summary -->
            <div class="pt-4 border-t border-zinc-200/80 space-y-2">
              <div class="flex justify-between text-xs text-zinc-500">
                <span>Subtotal Layanan</span>
                <span class="font-bold tabular-nums text-zinc-800">{{ formatRupiah(selectedInvoice.total_amount) }}</span>
              </div>
              <div class="flex justify-between text-xs text-zinc-500">
                <span>PPN 11% (Sudah Termasuk)</span>
                <span class="font-bold tabular-nums text-zinc-800">{{ formatRupiah(selectedInvoice.tax_share) }}</span>
              </div>
              <div class="flex justify-between items-baseline pt-4 border-t border-zinc-200">
                <span class="text-xs font-bold uppercase tracking-wider text-zinc-400">Total Tagihan</span>
                <span class="text-2xl sm:text-3xl font-extrabold tabular-nums tracking-tight text-zinc-950">
                  {{ formatRupiah(selectedInvoice.total_amount) }}
                </span>
              </div>
            </div>

            <div class="text-center pt-2 text-xs text-zinc-400 border-t border-dashed border-zinc-200">
              Bukti transaksi sah Semarang Medical Group • Disimpan sebagai tanda terima resmi
            </div>
          </div>

          <!-- 1B: Traditional Physical Thermal Receipt (Digunakan saat Cetak Printer & Dapat Di-preview) -->
          <div :class="['flex justify-center py-2', !previewThermal ? 'hidden print:flex' : '']">
            <div
              id="printable-receipt"
              class="w-full max-w-[340px] bg-white text-zinc-950 font-mono text-[11px] leading-tight p-6 sm:p-7 shadow-[0_12px_40px_rgba(0,0,0,0.08)] border border-black/10 relative"
            >
              <!-- Serrated Top Edge Decoration -->
              <div class="no-print absolute -top-2 left-0 right-0 h-2 overflow-hidden flex">
                <svg class="w-full h-2 text-white fill-current filter drop-shadow-[0_-2px_1px_rgba(0,0,0,0.05)]" viewBox="0 0 120 8" preserveAspectRatio="none">
                  <polygon points="0,8 3,0 6,8 9,0 12,8 15,0 18,8 21,0 24,8 27,0 30,8 33,0 36,8 39,0 42,8 45,0 48,8 51,0 54,8 57,0 60,8 63,0 66,8 69,0 72,8 75,0 78,8 81,0 84,8 87,0 90,8 93,0 96,8 99,0 102,8 105,0 108,8 111,0 114,8 117,0 120,8" />
                </svg>
              </div>

              <!-- Header Struk Tradisional -->
              <div class="text-center space-y-0.5 pb-2">
                <div class="font-black text-xs sm:text-sm tracking-tight uppercase">SEMARANG MEDICAL GROUP</div>
                <div class="text-[10px] text-zinc-600">Jl. Pandanaran No. 45, Semarang</div>
                <div class="text-[10px] text-zinc-600">Telp: (024) 841-2345 | WA: 0812-8877-6655</div>
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
                  <span>METODE</span>
                  <span class="font-bold">{{ selectedInvoice.payment_method }}</span>
                </div>
                <div v-if="selectedInvoice.payment_reference" class="flex justify-between">
                  <span>NO. TRACE / REF</span>
                  <span class="font-mono font-bold">{{ selectedInvoice.payment_reference }}</span>
                </div>
                <div class="flex justify-between">
                  <span>PASIEN</span>
                  <span class="font-bold uppercase truncate max-w-[170px] text-right">{{ selectedInvoice.patient_name }}</span>
                </div>
                <div class="flex justify-between">
                  <span>NIK</span>
                  <span>{{ selectedInvoice.patient_nik }}</span>
                </div>
              </div>

              <!-- Dashed line -->
              <div class="border-b border-dashed border-zinc-400 my-2.5"></div>

              <!-- Daftar Item Barang / Jasa -->
              <div class="space-y-2.5 text-[10px]">
                <!-- Registrasi -->
                <div v-if="selectedInvoice.registration_fee > 0" class="space-y-0.5">
                  <div class="font-bold uppercase">BIAYA REGISTRASI & ANTRIAN</div>
                  <div class="flex justify-between text-zinc-700">
                    <span>1 x {{ formatRupiah(selectedInvoice.registration_fee) }}</span>
                    <span class="font-bold">{{ formatRupiah(selectedInvoice.registration_fee) }}</span>
                  </div>
                </div>

                <!-- Konsultasi Dokter -->
                <div v-if="selectedInvoice.consultation_fee > 0" class="space-y-0.5 pt-1">
                  <div class="font-bold uppercase">KONSULTASI & TINDAKAN DOKTER</div>
                  <div class="flex justify-between text-zinc-700">
                    <span>1 x {{ formatRupiah(selectedInvoice.consultation_fee) }}</span>
                    <span class="font-bold">{{ formatRupiah(selectedInvoice.consultation_fee) }}</span>
                  </div>
                </div>

                <!-- Obat -->
                <template v-if="selectedInvoice.items && selectedInvoice.items.length > 0">
                  <div v-for="it in selectedInvoice.items" :key="it.id" class="space-y-0.5 pt-1">
                    <div class="font-bold uppercase leading-tight">{{ it.drug_name }}</div>
                    <div class="flex justify-between text-zinc-700">
                      <span>{{ it.quantity }} x {{ formatRupiah(it.unit_price) }}</span>
                      <span class="font-bold">{{ formatRupiah(it.subtotal) }}</span>
                    </div>
                  </div>
                </template>
                <div v-else-if="selectedInvoice.medication_fee > 0" class="space-y-0.5 pt-1">
                  <div class="font-bold uppercase">OBAT RESEP DOKTER</div>
                  <div class="flex justify-between text-zinc-700">
                    <span>1 x {{ formatRupiah(selectedInvoice.medication_fee) }}</span>
                    <span class="font-bold">{{ formatRupiah(selectedInvoice.medication_fee) }}</span>
                  </div>
                </div>
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
                  <span>STATUS</span>
                  <span class="font-bold">LUNAS</span>
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
                <div class="text-[8px] text-zinc-400 pt-0.5">POS Semarang Medical Group v2.5</div>
              </div>

              <!-- Serrated Bottom Edge Decoration -->
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
          class="bg-white rounded-2xl p-6 sm:p-7 border border-zinc-200/90 shadow-2xs space-y-6"
        >
          <div>
            <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-purple-50 text-purple-800 border border-purple-200/60 mb-2">
              Laporan Akuntansi Internal
            </span>
            <h3 class="text-lg font-extrabold text-zinc-950 tracking-tight">
              Alokasi Pendapatan & Bagi Hasil
            </h3>
            <p class="text-xs text-zinc-500 mt-0.5">
              Rincian ini memisahkan hak dokter medis (70% jasa konsultasi), margin bersih klinik, dan setoran pajak resmi.
            </p>
          </div>

          <!-- Cards breakdown: Hairline borders -->
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="p-4 rounded-xl bg-white border border-zinc-200/90 shadow-2xs">
              <div class="flex items-center justify-between text-xs font-bold text-zinc-700">
                <span>Jasa Medis (70%)</span>
                <UserCheck class="w-4 h-4 text-blue-600" />
              </div>
              <div class="text-xl font-extrabold tabular-nums tracking-tight text-zinc-950 mt-2">
                {{ formatRupiah(selectedInvoice.doctor_share) }}
              </div>
              <div class="text-[11px] text-zinc-500 mt-1">Hak Dokter Pemeriksa</div>
            </div>

            <div class="p-4 rounded-xl bg-white border border-zinc-200/90 shadow-2xs">
              <div class="flex items-center justify-between text-xs font-bold text-zinc-700">
                <span>Pendapatan Bersih Klinik</span>
                <Building2 class="w-4 h-4 text-emerald-600" />
              </div>
              <div class="text-xl font-extrabold tabular-nums tracking-tight text-zinc-950 mt-2">
                {{ formatRupiah(selectedInvoice.clinic_share) }}
              </div>
              <div class="text-[11px] text-zinc-500 mt-1">Operasional & sarana klinik</div>
            </div>

            <div class="p-4 rounded-xl bg-white border border-zinc-200/90 shadow-2xs">
              <div class="flex items-center justify-between text-xs font-bold text-zinc-700">
                <span>Pajak (PPN 11%)</span>
                <Percent class="w-4 h-4 text-amber-600" />
              </div>
              <div class="text-xl font-extrabold tabular-nums tracking-tight text-zinc-950 mt-2">
                {{ formatRupiah(selectedInvoice.tax_share) }}
              </div>
              <div class="text-[11px] text-zinc-500 mt-1">Setoran perpajakan resmi</div>
            </div>
          </div>

          <!-- Tabel Ringkasan Komponen Internal -->
          <div class="border border-zinc-200/90 rounded-xl overflow-hidden divide-y divide-zinc-200/70 text-xs">
            <div class="p-3 bg-zinc-50 font-bold text-zinc-700 flex justify-between">
              <span>Komponen Tagihan</span>
              <span>Nominal Masuk</span>
            </div>
            <div class="p-3 flex justify-between">
              <span class="text-zinc-600">Subtotal Registrasi</span>
              <span class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.registration_fee) }}</span>
            </div>
            <div class="p-3 flex justify-between">
              <span class="text-zinc-600">Subtotal Jasa Tindakan Dokter</span>
              <span class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.consultation_fee) }}</span>
            </div>
            <div class="p-3 flex justify-between">
              <span class="text-zinc-600">Subtotal Obat & Resep Farmasi</span>
              <span class="font-bold tabular-nums text-zinc-900">{{ formatRupiah(selectedInvoice.medication_fee) }}</span>
            </div>
            <div class="p-3 bg-zinc-50/70 font-extrabold text-zinc-950 flex justify-between text-sm">
              <span>Total Diterima Kasir</span>
              <span class="tabular-nums">{{ formatRupiah(selectedInvoice.total_amount) }}</span>
            </div>
          </div>
        </div>

      </section>

    </main>

  </div>
</template>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #printable-receipt,
  #printable-receipt * {
    visibility: visible;
  }
  #printable-receipt {
    position: absolute;
    left: 0;
    top: 0;
    width: 80mm !important;
    max-width: 80mm !important;
    margin: 0 !important;
    padding: 4mm !important;
    box-shadow: none !important;
    border: none !important;
  }
  .no-print {
    display: none !important;
  }
}
</style>
