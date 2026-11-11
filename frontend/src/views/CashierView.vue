<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  api,
  type PendingBill,
  type CashierInvoice,
  type CashierShiftSummary,
  type CashierShiftReport,
} from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useClinicStore } from '../stores/clinic'
import {
  Banknote,
  Receipt,
  UserPlus,
  Clock,
  Printer,
  RotateCcw,
  RefreshCw,
  CreditCard,
  Building2,
  CheckCircle2,
  AlertCircle,
  X,
  Send,
  Sparkles,
  Ticket,
  ChevronRight,
  ShieldCheck,
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

// Navigation Tabs
type CashierTab = 'pos' | 'walkin' | 'shift' | 'history'
const activeTab = ref<CashierTab>('pos')

// State: POS & Pending Bills
const pendingBills = ref<PendingBill[]>([])
const selectedBill = ref<PendingBill | null>(null)
const loadingBills = ref(false)
const filterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')

// Form: POS Payment
const customRegFee = ref(0)
const customActionFee = ref(0)
const customMedFee = ref(0)
const voucherDiscount = ref(0)
const paymentMethod = ref<'CASH' | 'BANK_TRANSFER' | 'EDC'>('CASH')
const cashReceived = ref<number>(0)
const paymentReference = ref('')
const paymentNotes = ref('')
const isSubmittingPayment = ref(false)

// State: Walk-in Registration
const walkinForm = ref({
  patient_name: '',
  patient_phone: '',
  clinic_unit: 'KIDINEST' as 'KIDINEST' | 'SSC',
  doctor_id: 2,
  doctor_name: 'dr. Hendra Pratama, Sp.A',
  payment_method: 'CASH' as 'CASH' | 'QRIS' | 'BANK_TRANSFER',
  is_paid_immediate: true,
})
const isSubmittingWalkin = ref(false)
const createdWalkinResult = ref<{ bill: PendingBill; invoice?: CashierInvoice } | null>(null)

// State: Shift Report
const shiftSummary = ref<CashierShiftSummary | null>(null)
const shiftReports = ref<CashierShiftReport[]>([])
const closeShiftNotes = ref('')
const isClosingShift = ref(false)
const shiftClosedSuccess = ref<CashierShiftReport | null>(null)

// State: Invoices History & Printing
const invoiceList = ref<CashierInvoice[]>([])
const loadingInvoices = ref(false)
const activeInvoiceForPrint = ref<CashierInvoice | null>(null)
const showPrintModal = ref(false)

// State: Refund Modal
const showRefundModal = ref(false)
const invoiceToRefund = ref<CashierInvoice | null>(null)
const refundForm = ref({
  refund_bank_name: 'BCA',
  refund_bank_account: '',
  refund_account_name: '',
  refund_reason: '',
})
const isSubmittingRefund = ref(false)

// Formatters
function formatRupiah(val?: number) {
  return 'Rp ' + (val || 0).toLocaleString('id-ID')
}

function formatDate(val?: string) {
  if (!val) return '-'
  try {
    return new Date(val).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return val
  }
}

// Calculations for POS
const grossTotal = computed(() => {
  return (customRegFee.value || 0) + (customActionFee.value || 0) + (customMedFee.value || 0)
})

const grandTotal = computed(() => {
  const net = grossTotal.value - (voucherDiscount.value || 0)
  return net > 0 ? net : 0
})

const cashChange = computed(() => {
  if (paymentMethod.value !== 'CASH') return 0
  const diff = (cashReceived.value || 0) - grandTotal.value
  return diff > 0 ? diff : 0
})

const filteredBills = computed(() => {
  if (filterUnit.value === 'ALL') return pendingBills.value
  return pendingBills.value.filter((b) => b.clinic_unit === filterUnit.value)
})

// Fetch Data
async function loadPendingBills() {
  loadingBills.value = true
  try {
    pendingBills.value = await api.getPendingBills()
    if (selectedBill.value) {
      // Refresh current selection
      const updated = pendingBills.value.find((b) => b.queue_id === selectedBill.value?.queue_id)
      if (updated) selectBill(updated)
      else selectedBill.value = null
    }
  } catch (err: any) {
    console.error('Gagal memuat pending bills:', err)
  } finally {
    loadingBills.value = false
  }
}

async function loadShiftSummary() {
  try {
    shiftSummary.value = await api.getCashierShiftSummary()
    shiftReports.value = await api.getCashierShiftReports()
  } catch (err: any) {
    console.error('Gagal memuat rekap shift:', err)
  }
}

async function loadInvoices() {
  loadingInvoices.value = true
  try {
    invoiceList.value = await api.getCashierInvoices(30, 0)
  } catch (err: any) {
    console.error('Gagal memuat riwayat invoice:', err)
  } finally {
    loadingInvoices.value = false
  }
}

function selectBill(bill: PendingBill) {
  selectedBill.value = bill
  customRegFee.value = bill.subtotal_registration || 0
  customActionFee.value = bill.subtotal_action || 0
  customMedFee.value = bill.subtotal_medication || 0
  voucherDiscount.value = 0
  cashReceived.value = 0
  paymentReference.value = ''
  paymentNotes.value = ''
}

function setQuickCash(amount: number) {
  cashReceived.value = amount
}

function setExactCash() {
  cashReceived.value = grandTotal.value
}

// POS Action: Pay Bill
async function submitPayment() {
  if (!selectedBill.value) return

  if (paymentMethod.value === 'CASH' && cashReceived.value < grandTotal.value) {
    alert('Uang tunai yang diterima kurang dari total tagihan!')
    return
  }

  isSubmittingPayment.value = true
  try {
    const inv = await api.payBill({
      queue_id: selectedBill.value.queue_id,
      subtotal_registration: customRegFee.value,
      subtotal_action: customActionFee.value,
      subtotal_medication: customMedFee.value,
      voucher_discount: voucherDiscount.value,
      payment_method: paymentMethod.value,
      cash_received: cashReceived.value,
      payment_reference: paymentReference.value,
      notes: paymentNotes.value,
    })

    clinicStore.playChime()
    activeInvoiceForPrint.value = inv
    showPrintModal.value = true

    // Refresh views
    selectedBill.value = null
    await Promise.all([loadPendingBills(), loadShiftSummary(), loadInvoices()])
  } catch (err: any) {
    alert(err.message || 'Gagal memproses pembayaran')
  } finally {
    isSubmittingPayment.value = false
  }
}

// Walk-in Registration
function onWalkinUnitChange() {
  if (walkinForm.value.clinic_unit === 'KIDINEST') {
    walkinForm.value.doctor_id = 2
    walkinForm.value.doctor_name = 'dr. Hendra Pratama, Sp.A'
    walkinForm.value.is_paid_immediate = true
  } else {
    walkinForm.value.doctor_id = 3
    walkinForm.value.doctor_name = 'dr. Nadia Kusuma, Sp.KK'
    walkinForm.value.is_paid_immediate = false
  }
}

async function submitWalkIn() {
  if (!walkinForm.value.patient_name || !walkinForm.value.patient_phone) {
    alert('Nama dan Nomor WhatsApp pasien wajib diisi!')
    return
  }

  isSubmittingWalkin.value = true
  try {
    const res = await api.registerWalkIn({
      patient_name: walkinForm.value.patient_name,
      patient_phone: walkinForm.value.patient_phone,
      clinic_unit: walkinForm.value.clinic_unit,
      doctor_id: walkinForm.value.doctor_id,
      doctor_name: walkinForm.value.doctor_name,
      payment_method: walkinForm.value.payment_method,
      is_paid_immediate: walkinForm.value.is_paid_immediate,
    })

    createdWalkinResult.value = res
    clinicStore.playChime()

    // Reset Form
    walkinForm.value.patient_name = ''
    walkinForm.value.patient_phone = ''

    await Promise.all([loadPendingBills(), loadShiftSummary()])
  } catch (err: any) {
    alert(err.message || 'Gagal registrasi pasien walk-in')
  } finally {
    isSubmittingWalkin.value = false
  }
}

// Close Shift
async function submitCloseShift() {
  if (!confirm('Apakah Anda yakin ingin menutup shift kasir saat ini dan mengirim rekap ke Pimpinan?')) {
    return
  }

  isClosingShift.value = true
  try {
    const report = await api.closeCashierShift(closeShiftNotes.value)
    shiftClosedSuccess.value = report
    closeShiftNotes.value = ''
    clinicStore.playChime()
    await loadShiftSummary()
  } catch (err: any) {
    alert(err.message || 'Gagal menutup shift')
  } finally {
    isClosingShift.value = false
  }
}

// Refund Flow
function openRefundModal(inv: CashierInvoice) {
  invoiceToRefund.value = inv
  refundForm.value = {
    refund_bank_name: 'BCA',
    refund_bank_account: '',
    refund_account_name: inv.patient_name,
    refund_reason: '',
  }
  showRefundModal.value = true
}

async function submitRefund() {
  if (!invoiceToRefund.value) return
  if (!refundForm.value.refund_bank_account || !refundForm.value.refund_reason) {
    alert('Nomor rekening dan alasan refund wajib diisi!')
    return
  }

  isSubmittingRefund.value = true
  try {
    await api.refundCashierInvoice(invoiceToRefund.value.id, refundForm.value)
    alert(`Refund untuk invoice ${invoiceToRefund.value.invoice_number} berhasil diproses.`)
    showRefundModal.value = false
    invoiceToRefund.value = null
    await Promise.all([loadInvoices(), loadPendingBills(), loadShiftSummary()])
  } catch (err: any) {
    alert(err.message || 'Gagal memproses refund')
  } finally {
    isSubmittingRefund.value = false
  }
}

function printReceipt() {
  window.print()
}

onMounted(() => {
  loadPendingBills()
  loadShiftSummary()
  loadInvoices()
})
</script>

<template>
  <div class="py-8 px-4 sm:px-8 max-w-7xl mx-auto space-y-6">
    <!-- Top Bar: Header & Active Staff -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-black/10">
      <div>
        <div class="flex items-center gap-2">
          <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-amber-50 text-amber-900 border border-amber-200">
            <span class="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span>
            Workstation Kasir & GRO
          </span>
          <span class="text-xs text-zinc-400 font-mono">SCMS POS v2.0</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-extrabold text-zinc-950 tracking-tight mt-1">
          Meja Depan & Kasir Terpadu
        </h1>
        <p class="text-xs text-zinc-500 mt-0.5">
          Pelayanan pembayaran multi-metode, bantuan walk-in paperless, dan rekapitulasi tutup shift kasir.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="bg-white rounded-2xl px-4 py-2 border border-black/5 shadow-sm text-right">
          <div class="text-[10px] uppercase font-bold text-zinc-400">Kasir Bertugas</div>
          <div class="text-xs font-bold text-zinc-900">{{ authStore.user?.name || 'Siti Sarah (GRO)' }}</div>
        </div>

        <button
          @click="loadPendingBills(); loadShiftSummary(); loadInvoices()"
          :disabled="loadingBills"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-2xl bg-white border border-black/10 hover:border-black/20 text-xs font-bold text-zinc-800 shadow-sm transition-all active:scale-95 disabled:opacity-50"
        >
          <RefreshCw :class="['w-3.5 h-3.5', loadingBills && 'animate-spin']" />
          Sync Data
        </button>
      </div>
    </div>

    <!-- Shift Financial Summary Banner -->
    <div class="grid grid-cols-2 sm:grid-cols-5 gap-3">
      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Setoran Tunai</span>
          <Banknote class="w-3.5 h-3.5 text-emerald-600" />
        </div>
        <div class="text-lg font-black text-zinc-950 mt-1 font-mono">
          {{ formatRupiah(shiftSummary?.total_cash) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Fisik di laci loket</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>EDC & Transfer</span>
          <CreditCard class="w-3.5 h-3.5 text-sky-600" />
        </div>
        <div class="text-lg font-black text-zinc-950 mt-1 font-mono">
          {{ formatRupiah((shiftSummary?.total_edc || 0) + (shiftSummary?.total_transfer || 0)) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Rekening penampung</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Potongan Voucher</span>
          <Ticket class="w-3.5 h-3.5 text-purple-600" />
        </div>
        <div class="text-lg font-black text-purple-950 mt-1 font-mono">
          {{ formatRupiah(shiftSummary?.total_voucher) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Diskon promosi</div>
      </div>

      <div class="bg-white rounded-2xl p-4 border border-black/5 shadow-sm">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Total Refund</span>
          <RotateCcw class="w-3.5 h-3.5 text-rose-600" />
        </div>
        <div class="text-lg font-black text-rose-600 mt-1 font-mono">
          -{{ formatRupiah(shiftSummary?.total_refund) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">Retur pembatalan</div>
      </div>

      <div class="col-span-2 sm:col-span-1 bg-zinc-950 text-white rounded-2xl p-4 shadow-sm flex flex-col justify-between">
        <div class="text-[11px] font-bold uppercase tracking-wider text-zinc-400 flex items-center justify-between">
          <span>Net Setoran Kasir</span>
          <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
        </div>
        <div class="text-xl font-black text-emerald-400 font-mono mt-1">
          {{ formatRupiah(shiftSummary?.net_total) }}
        </div>
        <div class="text-[10px] text-zinc-400 mt-0.5">{{ shiftSummary?.total_transactions || 0 }} transaksi shift</div>
      </div>
    </div>

    <!-- Navigation Pills -->
    <div class="flex items-center gap-2 border-b border-black/5 pb-2">
      <button
        @click="activeTab = 'pos'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'pos'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <Receipt class="w-4 h-4" />
        Meja Kasir & POS
        <span
          v-if="pendingBills.length > 0"
          :class="[
            'px-2 py-0.5 rounded-full text-[10px] font-bold font-mono',
            activeTab === 'pos' ? 'bg-amber-400 text-zinc-950' : 'bg-amber-100 text-amber-900',
          ]"
        >
          {{ pendingBills.length }}
        </span>
      </button>

      <button
        @click="activeTab = 'walkin'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'walkin'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <UserPlus class="w-4 h-4" />
        Bantuan Walk-In (GRO)
      </button>

      <button
        @click="activeTab = 'shift'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'shift'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <Clock class="w-4 h-4" />
        Tutup Shift Kasir
      </button>

      <button
        @click="activeTab = 'history'"
        :class="[
          'inline-flex items-center gap-2 px-5 py-2.5 rounded-full text-xs font-bold transition-all',
          activeTab === 'history'
            ? 'bg-zinc-950 text-white shadow-sm'
            : 'bg-white hover:bg-zinc-100 text-zinc-700 border border-black/5',
        ]"
      >
        <Printer class="w-4 h-4" />
        Riwayat & Cetak Struk
      </button>
    </div>

    <!-- TAB 1: MEJA KASIR & POS -->
    <div v-if="activeTab === 'pos'" class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left Column: Pending Bills List (5 cols) -->
      <div class="lg:col-span-5 space-y-4">
        <div class="bg-white rounded-3xl p-5 border border-black/5 shadow-sm space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-sm font-extrabold text-zinc-950 uppercase tracking-wider">
              Antrian Menunggu Pembayaran
            </h2>
            <!-- Filter Unit -->
            <div class="flex items-center bg-zinc-100 rounded-full p-0.5 text-[11px] font-bold">
              <button
                @click="filterUnit = 'ALL'"
                :class="['px-2.5 py-1 rounded-full transition-all', filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
              >
                Semua
              </button>
              <button
                @click="filterUnit = 'KIDINEST'"
                :class="['px-2.5 py-1 rounded-full transition-all', filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
              >
                Kidinest
              </button>
              <button
                @click="filterUnit = 'SSC'"
                :class="['px-2.5 py-1 rounded-full transition-all', filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-xs' : 'text-zinc-500']"
              >
                SSC
              </button>
            </div>
          </div>

          <div v-if="loadingBills" class="py-12 text-center text-xs text-zinc-400">
            <RefreshCw class="w-5 h-5 mx-auto animate-spin mb-2 text-zinc-400" />
            Memuat antrian pembayaran...
          </div>

          <div v-else-if="filteredBills.length === 0" class="py-12 text-center text-xs text-zinc-400">
            <CheckCircle2 class="w-8 h-8 mx-auto text-emerald-500 mb-2" />
            Tidak ada antrian tertunda. Semua pembayaran lunas!
          </div>

          <div v-else class="space-y-2.5 max-h-[580px] overflow-y-auto pr-1">
            <div
              v-for="b in filteredBills"
              :key="b.queue_id"
              @click="selectBill(b)"
              :class="[
                'p-4 rounded-2xl border transition-all cursor-pointer flex items-center justify-between',
                selectedBill?.queue_id === b.queue_id
                  ? 'bg-zinc-950 text-white border-zinc-950 shadow-md ring-2 ring-zinc-900/10'
                  : 'bg-zinc-50/60 hover:bg-zinc-100 text-zinc-900 border-black/5',
              ]"
            >
              <div class="space-y-1">
                <div class="flex items-center gap-2">
                  <span
                    :class="[
                      'px-2.5 py-0.5 rounded-full font-mono text-xs font-black',
                      selectedBill?.queue_id === b.queue_id ? 'bg-white text-zinc-950' : 'bg-zinc-950 text-white',
                    ]"
                  >
                    {{ b.queue_number }}
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-bold px-2 py-0.5 rounded-full',
                      b.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-100 text-amber-900'
                        : 'bg-rose-100 text-rose-900',
                    ]"
                  >
                    {{ b.clinic_unit }}
                  </span>
                  <span
                    v-if="b.payment_status === 'PENDING'"
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-rose-100 text-rose-800"
                  >
                    Tunai Registrasi
                  </span>
                  <span
                    v-else
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-blue-100 text-blue-800"
                  >
                    Pelunasan Tindakan
                  </span>
                </div>

                <div class="font-extrabold text-sm tracking-tight">{{ b.patient_name }}</div>
                <div :class="['text-xs', selectedBill?.queue_id === b.queue_id ? 'text-zinc-400' : 'text-zinc-500']">
                  Dokter: {{ b.doctor_name }}
                </div>
              </div>

              <div class="text-right space-y-1">
                <div
                  :class="[
                    'font-mono font-black text-sm',
                    selectedBill?.queue_id === b.queue_id ? 'text-emerald-400' : 'text-zinc-950',
                  ]"
                >
                  {{ formatRupiah(b.total_due) }}
                </div>
                <div class="text-[10px] text-zinc-400 flex items-center justify-end gap-1">
                  Pilih <ChevronRight class="w-3 h-3" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: POS Billing Register & Calculations (7 cols) -->
      <div class="lg:col-span-7">
        <div v-if="!selectedBill" class="bg-white rounded-3xl p-12 border border-black/5 text-center shadow-sm">
          <Receipt class="w-12 h-12 mx-auto text-zinc-300 mb-3" />
          <h3 class="text-base font-extrabold text-zinc-800">Pilih Pasien untuk Membuka Kasir</h3>
          <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
            Klik salah satu antrian pasien di kolom sebelah kiri untuk menginput rincian biaya, voucher diskon, dan metode pembayaran.
          </p>
        </div>

        <div v-else class="bg-white rounded-3xl p-6 sm:p-7 border border-black/5 shadow-sm space-y-6">
          <!-- Selected Patient Banner -->
          <div class="flex items-start justify-between pb-4 border-b border-black/5">
            <div>
              <div class="flex items-center gap-2">
                <span class="px-3 py-1 rounded-full font-mono text-sm font-black bg-zinc-950 text-white">
                  {{ selectedBill.queue_number }}
                </span>
                <span class="text-xs font-bold px-2.5 py-0.5 rounded-full bg-zinc-100 text-zinc-700">
                  {{ selectedBill.clinic_unit }}
                </span>
              </div>
              <h2 class="text-xl font-black text-zinc-950 mt-2">{{ selectedBill.patient_name }}</h2>
              <p class="text-xs text-zinc-500">
                WhatsApp: {{ selectedBill.patient_phone || '-' }} • Dokter: {{ selectedBill.doctor_name }}
              </p>
            </div>

            <button
              @click="selectedBill = null"
              class="text-zinc-400 hover:text-zinc-700 p-1 rounded-full hover:bg-zinc-100 transition-all"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Bill Items Breakdown -->
          <div class="space-y-3">
            <div class="text-xs font-extrabold uppercase tracking-wider text-zinc-400">
              Rincian Item Tagihan Single Invoice
            </div>

            <div class="space-y-2.5 bg-zinc-50 rounded-2xl p-4 border border-black/5">
              <!-- Item 1: Registration Fee -->
              <div class="flex items-center justify-between text-xs">
                <span class="font-medium text-zinc-700">Biaya Pendaftaran Loket:</span>
                <input
                  v-model.number="customRegFee"
                  type="number"
                  min="0"
                  class="w-32 px-3 py-1.5 text-right font-mono text-xs font-bold rounded-xl bg-white border border-black/10 focus:ring-2 focus:ring-zinc-950"
                />
              </div>

              <!-- Item 2: Doctor Action / Consultation Fee -->
              <div class="flex items-center justify-between text-xs">
                <div>
                  <div class="font-medium text-zinc-700">Jasa Tindakan Medis / Konsultasi Dokter:</div>
                  <div class="text-[10px] text-zinc-400">Vaksinasi, pemeriksaan dokter, atau facial/laser SSC</div>
                </div>
                <input
                  v-model.number="customActionFee"
                  type="number"
                  min="0"
                  class="w-32 px-3 py-1.5 text-right font-mono text-xs font-bold rounded-xl bg-white border border-black/10 focus:ring-2 focus:ring-zinc-950"
                />
              </div>

              <!-- Item 3: Pharmacy / Medication Fee -->
              <div class="flex items-center justify-between text-xs">
                <div>
                  <div class="font-medium text-zinc-700">Biaya Obat & Resep Farmasi:</div>
                  <div class="text-[10px] text-zinc-400">Obat paten, racikan, atau produk skincare SSC</div>
                </div>
                <input
                  v-model.number="customMedFee"
                  type="number"
                  min="0"
                  class="w-32 px-3 py-1.5 text-right font-mono text-xs font-bold rounded-xl bg-white border border-black/10 focus:ring-2 focus:ring-zinc-950"
                />
              </div>

              <div class="pt-2 border-t border-black/5 flex items-center justify-between text-xs font-bold text-zinc-600">
                <span>Subtotal Tagihan Kotor:</span>
                <span class="font-mono">{{ formatRupiah(grossTotal) }}</span>
              </div>
            </div>
          </div>

          <!-- Voucher Discount Section -->
          <div class="bg-purple-50/50 rounded-2xl p-4 border border-purple-100 space-y-2">
            <div class="flex items-center justify-between text-xs font-bold text-purple-950">
              <span class="flex items-center gap-1.5">
                <Ticket class="w-4 h-4 text-purple-600" />
                Voucher Diskon Promosi
              </span>
              <span class="text-[11px] text-purple-700 font-normal">Potongan langsung tagihan</span>
            </div>
            <div class="flex items-center gap-2">
              <input
                v-model.number="voucherDiscount"
                type="number"
                min="0"
                :max="grossTotal"
                placeholder="Nominal potongan (Rp)"
                class="flex-1 px-3 py-2 text-right font-mono text-xs font-bold rounded-xl bg-white border border-purple-200 focus:ring-2 focus:ring-purple-950"
              />
              <button
                @click="voucherDiscount = 20000"
                type="button"
                class="px-3 py-2 rounded-xl bg-purple-100 hover:bg-purple-200 text-purple-900 text-xs font-bold"
              >
                Voucher 20k
              </button>
              <button
                @click="voucherDiscount = 50000"
                type="button"
                class="px-3 py-2 rounded-xl bg-purple-100 hover:bg-purple-200 text-purple-900 text-xs font-bold"
              >
                Voucher 50k
              </button>
            </div>
          </div>

          <!-- Grand Total Banner -->
          <div class="bg-zinc-950 text-white rounded-2xl p-5 flex items-center justify-between">
            <div>
              <div class="text-[11px] uppercase font-bold text-zinc-400">Total Tagihan Bersih</div>
              <div class="text-xs text-zinc-400 mt-0.5">Sudah termasuk potongan voucher</div>
            </div>
            <div class="text-2xl sm:text-3xl font-black font-mono text-emerald-400">
              {{ formatRupiah(grandTotal) }}
            </div>
          </div>

          <!-- Payment Method Selector -->
          <div class="space-y-3">
            <label class="text-xs font-extrabold uppercase tracking-wider text-zinc-400">
              Pilih Metode Pembayaran
            </label>
            <div class="grid grid-cols-3 gap-2">
              <button
                type="button"
                @click="paymentMethod = 'CASH'"
                :class="[
                  'p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all',
                  paymentMethod === 'CASH'
                    ? 'bg-zinc-950 text-white border-zinc-950 shadow-sm'
                    : 'bg-zinc-50 hover:bg-zinc-100 text-zinc-800 border-black/5',
                ]"
              >
                <Banknote class="w-4 h-4" />
                Tunai (Cash)
              </button>

              <button
                type="button"
                @click="paymentMethod = 'EDC'"
                :class="[
                  'p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all',
                  paymentMethod === 'EDC'
                    ? 'bg-zinc-950 text-white border-zinc-950 shadow-sm'
                    : 'bg-zinc-50 hover:bg-zinc-100 text-zinc-800 border-black/5',
                ]"
              >
                <CreditCard class="w-4 h-4" />
                Mesin EDC
              </button>

              <button
                type="button"
                @click="paymentMethod = 'BANK_TRANSFER'"
                :class="[
                  'p-3 rounded-2xl border text-xs font-bold flex flex-col items-center gap-1.5 transition-all',
                  paymentMethod === 'BANK_TRANSFER'
                    ? 'bg-zinc-950 text-white border-zinc-950 shadow-sm'
                    : 'bg-zinc-50 hover:bg-zinc-100 text-zinc-800 border-black/5',
                ]"
              >
                <Building2 class="w-4 h-4" />
                Transfer Bank
              </button>
            </div>

            <!-- Cash Tender & Change Calculator -->
            <div v-if="paymentMethod === 'CASH'" class="bg-emerald-50/50 rounded-2xl p-4 border border-emerald-100 space-y-3">
              <div class="flex items-center justify-between text-xs font-bold text-emerald-950">
                <span>Uang Diterima dari Pasien:</span>
                <span class="font-normal text-[11px] text-emerald-800">Kalkulator Kembalian Otomatis</span>
              </div>

              <div class="flex items-center gap-2">
                <input
                  v-model.number="cashReceived"
                  type="number"
                  placeholder="Jumlah uang tunai (Rp)"
                  class="flex-1 px-3 py-2 text-right font-mono text-base font-extrabold rounded-xl bg-white border border-emerald-300 focus:ring-2 focus:ring-emerald-900"
                />
                <button
                  type="button"
                  @click="setExactCash"
                  class="px-3 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-bold transition-all"
                >
                  Uang Pas
                </button>
              </div>

              <!-- Quick cash pills -->
              <div class="flex items-center gap-1.5 flex-wrap">
                <button
                  type="button"
                  @click="setQuickCash(50000)"
                  class="px-2.5 py-1 rounded-full bg-white border border-emerald-200 text-emerald-950 text-[11px] font-bold"
                >
                  Rp 50.000
                </button>
                <button
                  type="button"
                  @click="setQuickCash(100000)"
                  class="px-2.5 py-1 rounded-full bg-white border border-emerald-200 text-emerald-950 text-[11px] font-bold"
                >
                  Rp 100.000
                </button>
                <button
                  type="button"
                  @click="setQuickCash(200000)"
                  class="px-2.5 py-1 rounded-full bg-white border border-emerald-200 text-emerald-950 text-[11px] font-bold"
                >
                  Rp 200.000
                </button>
              </div>

              <div class="pt-2 border-t border-emerald-200/60 flex items-center justify-between">
                <span class="text-xs font-bold text-emerald-900">Uang Kembalian Pasien:</span>
                <span class="text-lg font-black font-mono text-emerald-700">
                  {{ formatRupiah(cashChange) }}
                </span>
              </div>
            </div>

            <!-- Non-cash Reference -->
            <div v-else class="space-y-2">
              <label class="text-xs font-bold text-zinc-700">
                Nomor Referensi Transaksi (Auth Code / No Rekening):
              </label>
              <input
                v-model="paymentReference"
                type="text"
                placeholder="Contoh: EDC-BCA-987123 atau TRF-MANDIRI-445"
                class="w-full px-3 py-2 text-xs rounded-xl bg-zinc-50 border border-black/10 focus:ring-2 focus:ring-zinc-950"
              />
            </div>
          </div>

          <!-- Submit Payment Button -->
          <div class="pt-2 flex items-center justify-between gap-3">
            <button
              @click="submitPayment"
              :disabled="isSubmittingPayment || (paymentMethod === 'CASH' && cashReceived < grandTotal)"
              class="flex-1 py-3.5 px-6 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-sm shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <CheckCircle2 class="w-4 h-4 text-emerald-400" />
              {{ isSubmittingPayment ? 'Memproses Transaksi...' : 'Terima Pembayaran & Cetak Struk' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: BANTUAN REGISTRASI WALK-IN (GRO) -->
    <div v-if="activeTab === 'walkin'" class="max-w-3xl mx-auto space-y-6">
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-sm space-y-6">
        <div>
          <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-zinc-100 text-zinc-800">
            Resepsionis Meja Depan (Paperless)
          </span>
          <h2 class="text-xl sm:text-2xl font-black text-zinc-950 mt-1">
            Bantuan Pendaftaran Pasien Walk-In
          </h2>
          <p class="text-xs text-zinc-500 mt-1">
            Untuk pasien lansia, anak darurat, atau yang datang langsung tanpa smartphone. Nomor tiket antrian otomatis terkirim ke WhatsApp.
          </p>
        </div>

        <div class="space-y-4">
          <!-- Unit Selector -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-zinc-700">Pilih Unit Klinik:</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                @click="walkinForm.clinic_unit = 'KIDINEST'; onWalkinUnitChange()"
                :class="[
                  'p-4 rounded-2xl border text-left transition-all',
                  walkinForm.clinic_unit === 'KIDINEST'
                    ? 'bg-amber-500 text-white border-amber-500 shadow-md ring-2 ring-amber-400/20'
                    : 'bg-zinc-50 hover:bg-zinc-100 text-zinc-800 border-black/5',
                ]"
              >
                <div class="font-extrabold text-sm">Kidinest</div>
                <div :class="['text-xs mt-0.5', walkinForm.clinic_unit === 'KIDINEST' ? 'text-amber-100' : 'text-zinc-500']">
                  Poli Anak, Vaksinasi & Tumbang (Biaya Loket Rp 50.000)
                </div>
              </button>

              <button
                type="button"
                @click="walkinForm.clinic_unit = 'SSC'; onWalkinUnitChange()"
                :class="[
                  'p-4 rounded-2xl border text-left transition-all',
                  walkinForm.clinic_unit === 'SSC'
                    ? 'bg-rose-500 text-white border-rose-500 shadow-md ring-2 ring-rose-400/20'
                    : 'bg-zinc-50 hover:bg-zinc-100 text-zinc-800 border-black/5',
                ]"
              >
                <div class="font-extrabold text-sm">Semarang Skin Center</div>
                <div :class="['text-xs mt-0.5', walkinForm.clinic_unit === 'SSC' ? 'text-rose-100' : 'text-zinc-500']">
                  Estetika Kulit & Facial (Booking Jadwal Rp0 di Awal)
                </div>
              </button>
            </div>
          </div>

          <!-- Patient Fields -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-zinc-700">Nama Lengkap Pasien / Anak:</label>
              <input
                v-model="walkinForm.patient_name"
                type="text"
                placeholder="Contoh: Ananda Rayyan Pratama"
                class="w-full px-4 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs font-medium focus:ring-2 focus:ring-zinc-950"
              />
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-bold text-zinc-700">Nomor WhatsApp Pasien / Orang Tua:</label>
              <input
                v-model="walkinForm.patient_phone"
                type="tel"
                placeholder="081234567890"
                class="w-full px-4 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs font-medium font-mono focus:ring-2 focus:ring-zinc-950"
              />
            </div>
          </div>

          <!-- Doctor Picker -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-zinc-700">Dokter Jaga Bertugas:</label>
            <select
              v-model="walkinForm.doctor_name"
              class="w-full px-4 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs font-medium focus:ring-2 focus:ring-zinc-950"
            >
              <option value="dr. Hendra Pratama, Sp.A">dr. Hendra Pratama, Sp.A (Spesialis Anak - Kidinest)</option>
              <option value="dr. Nadia Kusuma, Sp.KK">dr. Nadia Kusuma, Sp.KK (Spesialis Kulit - SSC)</option>
            </select>
          </div>

          <!-- Immediate Payment Checkbox (For Kidinest Cash walk-in) -->
          <div v-if="walkinForm.clinic_unit === 'KIDINEST'" class="bg-amber-50 rounded-2xl p-4 border border-amber-200 space-y-2">
            <label class="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                v-model="walkinForm.is_paid_immediate"
                class="w-4 h-4 rounded text-zinc-950 focus:ring-zinc-950"
              />
              <div>
                <div class="text-xs font-extrabold text-amber-950">
                  Terima Setoran Tunai Rp 50.000 Langsung di Loket Meja Depan
                </div>
                <div class="text-[11px] text-amber-800">
                  Nomor antrian FIFO langsung berstatus Lunas (PAID) dan siap dipanggil dokter.
                </div>
              </div>
            </label>
          </div>

          <button
            @click="submitWalkIn"
            :disabled="isSubmittingWalkin"
            class="w-full py-3.5 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-sm shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <UserPlus class="w-4 h-4 text-emerald-400" />
            {{ isSubmittingWalkin ? 'Mendaftarkan Pasien...' : 'Daftarkan Walk-In & Terbitkan Nomor Antrian' }}
          </button>
        </div>
      </div>

      <!-- Result Card: Walk-in Registered -->
      <div v-if="createdWalkinResult" class="bg-emerald-50 rounded-3xl p-6 border border-emerald-200 space-y-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="w-5 h-5 text-emerald-600" />
            <h3 class="font-extrabold text-emerald-950 text-sm">Pendaftaran Walk-In Berhasil!</h3>
          </div>
          <button @click="createdWalkinResult = null" class="text-emerald-700 hover:text-emerald-950">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="bg-white rounded-2xl p-5 border border-emerald-100 flex items-center justify-between">
          <div>
            <div class="text-[10px] uppercase font-bold text-zinc-400">Nomor Antrian Resmi</div>
            <div class="text-3xl font-black font-mono text-zinc-950 mt-0.5">
              {{ createdWalkinResult.bill.queue_number }}
            </div>
            <div class="text-xs font-bold text-zinc-700 mt-1">
              {{ createdWalkinResult.bill.patient_name }} • {{ createdWalkinResult.bill.clinic_unit }}
            </div>
          </div>

          <div class="text-right">
            <span
              :class="[
                'px-3 py-1 rounded-full text-xs font-bold',
                createdWalkinResult.bill.payment_status === 'PAID'
                  ? 'bg-emerald-100 text-emerald-900'
                  : 'bg-rose-100 text-rose-900',
              ]"
            >
              {{ createdWalkinResult.bill.payment_status === 'PAID' ? 'LUNAS (Aktif Menunggu)' : 'Pending Kasir' }}
            </span>
            <div class="text-[10px] text-zinc-400 mt-2">
              Tiket paperless terkirim ke WhatsApp
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: TUTUP SHIFT KASIR (END-OF-SHIFT) -->
    <div v-if="activeTab === 'shift'" class="max-w-4xl mx-auto space-y-6">
      <div class="bg-white rounded-3xl p-6 sm:p-8 border border-black/5 shadow-sm space-y-6">
        <div class="flex items-start justify-between">
          <div>
            <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-amber-50 text-amber-900 border border-amber-200">
              Laporan Pertanggungjawaban Loket
            </span>
            <h2 class="text-xl sm:text-2xl font-black text-zinc-950 mt-1">
              Rekapitulasi Tutup Shift Kasir Harian
            </h2>
            <p class="text-xs text-zinc-500 mt-1">
              Kirim rekapitulasi setoran fisik laci dan transaksi digital langsung ke dashboard Pimpinan klinik.
            </p>
          </div>

          <div class="text-right">
            <div class="text-[10px] uppercase font-bold text-zinc-400">Tanggal Shift</div>
            <div class="text-xs font-mono font-bold text-zinc-900">{{ shiftSummary?.shift_date || '-' }}</div>
          </div>
        </div>

        <!-- Breakdown Grid -->
        <div class="bg-zinc-50 rounded-2xl p-5 border border-black/5 space-y-3">
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Total Transaksi Selesai:</span>
            <span class="font-mono font-bold text-zinc-900">{{ shiftSummary?.total_transactions || 0 }} Transaksi</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Penerimaan Tunai (Uang Fisik Kasir):</span>
            <span class="font-mono font-bold text-emerald-600">{{ formatRupiah(shiftSummary?.total_cash) }}</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Penerimaan Kartu EDC Loket:</span>
            <span class="font-mono font-bold text-sky-600">{{ formatRupiah(shiftSummary?.total_edc) }}</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Penerimaan Transfer Bank Rekening:</span>
            <span class="font-mono font-bold text-indigo-600">{{ formatRupiah(shiftSummary?.total_transfer) }}</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Total Potongan Diskon Voucher:</span>
            <span class="font-mono font-bold text-purple-600">{{ formatRupiah(shiftSummary?.total_voucher) }}</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-zinc-600">Pengembalian Dana Pasien (Refund):</span>
            <span class="font-mono font-bold text-rose-600">-{{ formatRupiah(shiftSummary?.total_refund) }}</span>
          </div>
          <div class="pt-3 border-t border-black/10 flex items-center justify-between">
            <span class="font-extrabold text-sm text-zinc-950">Grand Net Setoran Kasir ke Pimpinan:</span>
            <span class="text-xl font-black font-mono text-zinc-950">{{ formatRupiah(shiftSummary?.net_total) }}</span>
          </div>
        </div>

        <!-- Shift Notes -->
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-zinc-700">Catatan Kasir untuk Pimpinan (Opsional):</label>
          <textarea
            v-model="closeShiftNotes"
            rows="2"
            placeholder="Contoh: Shift pagi berjalan tertib. Uang fisik di laci cocok Rp 350.000."
            class="w-full px-4 py-2.5 rounded-xl bg-zinc-50 border border-black/10 text-xs focus:ring-2 focus:ring-zinc-950"
          ></textarea>
        </div>

        <button
          @click="submitCloseShift"
          :disabled="isClosingShift"
          class="w-full py-3.5 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-sm shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2"
        >
          <Send class="w-4 h-4 text-amber-400" />
          {{ isClosingShift ? 'Mengirim Rekapitulasi...' : 'Kunci & Kirim Laporan Tutup Shift ke Pimpinan' }}
        </button>

        <!-- Success Alert -->
        <div v-if="shiftClosedSuccess" class="bg-emerald-50 rounded-2xl p-4 border border-emerald-200 text-xs text-emerald-950 space-y-1">
          <div class="font-bold flex items-center gap-1.5">
            <CheckCircle2 class="w-4 h-4 text-emerald-600" />
            Laporan Tutup Shift Berhasil Disimpan & Dikirim!
          </div>
          <div>Laporan ID: #{{ shiftClosedSuccess.id }} • Waktu Kirim: {{ formatDate(shiftClosedSuccess.sent_to_owner_at) }}</div>
        </div>
      </div>

      <!-- Closed Shifts History Table -->
      <div class="bg-white rounded-3xl p-6 border border-black/5 shadow-sm space-y-4">
        <h3 class="text-sm font-extrabold text-zinc-950 uppercase tracking-wider">
          Riwayat Laporan Tutup Shift Kasir
        </h3>

        <div v-if="shiftReports.length === 0" class="py-8 text-center text-xs text-zinc-400">
          Belum ada riwayat laporan tutup shift sebelumnya.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="border-b border-black/5 text-zinc-400 font-semibold uppercase tracking-wider text-[10px]">
                <th class="py-2.5">Tanggal</th>
                <th class="py-2.5">Kasir</th>
                <th class="py-2.5">Transaksi</th>
                <th class="py-2.5">Tunai</th>
                <th class="py-2.5">Transfer/EDC</th>
                <th class="py-2.5">Net Total</th>
                <th class="py-2.5">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-black/5">
              <tr v-for="rep in shiftReports" :key="rep.id" class="hover:bg-zinc-50">
                <td class="py-3 font-mono text-zinc-600">{{ rep.shift_date }}</td>
                <td class="py-3 font-bold text-zinc-900">{{ rep.staff_name }}</td>
                <td class="py-3 font-mono text-zinc-700">{{ rep.total_transactions }}</td>
                <td class="py-3 font-mono text-emerald-600 font-bold">{{ formatRupiah(rep.total_cash) }}</td>
                <td class="py-3 font-mono text-zinc-700">{{ formatRupiah(rep.total_transfer + rep.total_edc) }}</td>
                <td class="py-3 font-mono font-black text-zinc-950">{{ formatRupiah(rep.net_total) }}</td>
                <td class="py-3">
                  <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-900">
                    Terkirim ke Pimpinan
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB 4: RIWAYAT & CETAK STRUK -->
    <div v-if="activeTab === 'history'" class="bg-white rounded-3xl p-6 sm:p-7 border border-black/5 shadow-sm space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-black/5">
        <div>
          <h2 class="text-base font-extrabold text-zinc-950">Riwayat Pembayaran & Struk Kasir</h2>
          <p class="text-xs text-zinc-500 mt-0.5">Daftar transaksi terlunasi, cetak ulang struk thermal 80mm, atau catat refund manual.</p>
        </div>

        <button
          @click="loadInvoices"
          :disabled="loadingInvoices"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-zinc-100 text-zinc-700 text-xs font-bold hover:bg-zinc-200 transition-all"
        >
          <RefreshCw :class="['w-3.5 h-3.5', loadingInvoices && 'animate-spin']" />
          Segarkan
        </button>
      </div>

      <div v-if="invoiceList.length === 0" class="py-12 text-center text-xs text-zinc-400">
        Belum ada riwayat transaksi pembayaran kasir.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="border-b border-black/5 text-zinc-400 font-semibold uppercase tracking-wider text-[10px]">
              <th class="py-2.5">No Invoice</th>
              <th class="py-2.5">Antrian</th>
              <th class="py-2.5">Pasien</th>
              <th class="py-2.5">Unit</th>
              <th class="py-2.5">Metode</th>
              <th class="py-2.5">Total Bayar</th>
              <th class="py-2.5">Status</th>
              <th class="py-2.5 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-black/5">
            <tr v-for="inv in invoiceList" :key="inv.id" class="hover:bg-zinc-50/80">
              <td class="py-3 font-mono font-bold text-zinc-900">{{ inv.invoice_number }}</td>
              <td class="py-3 font-mono font-extrabold text-zinc-700">{{ inv.queue_number }}</td>
              <td class="py-3 font-bold text-zinc-900">{{ inv.patient_name }}</td>
              <td class="py-3">
                <span
                  :class="[
                    'px-2 py-0.5 rounded-full text-[10px] font-bold',
                    inv.clinic_unit === 'KIDINEST' ? 'bg-amber-100 text-amber-900' : 'bg-rose-100 text-rose-900',
                  ]"
                >
                  {{ inv.clinic_unit }}
                </span>
              </td>
              <td class="py-3 font-bold text-zinc-600">{{ inv.payment_method }}</td>
              <td class="py-3 font-mono font-black text-zinc-950">{{ formatRupiah(inv.total_amount) }}</td>
              <td class="py-3">
                <span
                  :class="[
                    'px-2.5 py-0.5 rounded-full text-[10px] font-bold',
                    inv.status === 'PAID' ? 'bg-emerald-100 text-emerald-900' : 'bg-rose-100 text-rose-900',
                  ]"
                >
                  {{ inv.status === 'PAID' ? 'LUNAS' : 'REFUNDED' }}
                </span>
              </td>
              <td class="py-3 text-right space-x-2">
                <button
                  @click="activeInvoiceForPrint = inv; showPrintModal = true"
                  class="inline-flex items-center gap-1 px-3 py-1 rounded-xl bg-zinc-100 hover:bg-zinc-200 text-zinc-800 text-[11px] font-bold transition-all"
                >
                  <Printer class="w-3 h-3" />
                  Struk 80mm
                </button>

                <button
                  v-if="inv.status === 'PAID'"
                  @click="openRefundModal(inv)"
                  class="inline-flex items-center gap-1 px-3 py-1 rounded-xl bg-rose-50 hover:bg-rose-100 text-rose-800 text-[11px] font-bold transition-all"
                >
                  <RotateCcw class="w-3 h-3" />
                  Refund
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL 1: CETAK STRUK THERMAL 80MM / PDF INVOICE -->
    <div
      v-if="showPrintModal && activeInvoiceForPrint"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-sm w-full p-6 space-y-4 shadow-2xl border border-black/10">
        <div class="flex items-center justify-between pb-3 border-b border-black/10">
          <div class="flex items-center gap-1.5 text-xs font-bold text-zinc-600">
            <Printer class="w-4 h-4 text-zinc-800" />
            Cetak Struk Thermal 80mm
          </div>
          <button @click="showPrintModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Thermal Receipt Format (80mm Width) -->
        <div id="thermal-receipt" class="bg-zinc-50 rounded-2xl p-4 font-mono text-xs text-zinc-900 space-y-3 border border-black/5">
          <div class="text-center space-y-1 pb-2 border-b border-dashed border-black/20">
            <div class="font-extrabold text-sm tracking-tight">SMART CLINIC MEDICAL</div>
            <div class="text-[10px] text-zinc-500">
              {{ activeInvoiceForPrint.clinic_unit === 'KIDINEST' ? 'Unit Kidinest (Poli Anak & Vaksinasi)' : 'Semarang Skin Center (SSC)' }}
            </div>
            <div class="text-[10px] text-zinc-400">Jl. Pemuda No. 88, Semarang • 024-87654321</div>
          </div>

          <div class="space-y-1 text-[11px]">
            <div class="flex justify-between">
              <span class="text-zinc-500">No Invoice:</span>
              <span class="font-bold">{{ activeInvoiceForPrint.invoice_number }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-500">No Antrian:</span>
              <span class="font-black text-sm">{{ activeInvoiceForPrint.queue_number }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-500">Pasien:</span>
              <span class="font-bold">{{ activeInvoiceForPrint.patient_name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-500">Waktu:</span>
              <span>{{ formatDate(activeInvoiceForPrint.created_at) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-zinc-500">Kasir:</span>
              <span>{{ activeInvoiceForPrint.cashier_name }}</span>
            </div>
          </div>

          <div class="pt-2 border-t border-dashed border-black/20 space-y-1 text-[11px]">
            <div v-if="activeInvoiceForPrint.subtotal_registration > 0" class="flex justify-between">
              <span>Registrasi Loket:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.subtotal_registration) }}</span>
            </div>
            <div v-if="activeInvoiceForPrint.subtotal_action > 0" class="flex justify-between">
              <span>Jasa Tindakan/Konsultasi:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.subtotal_action) }}</span>
            </div>
            <div v-if="activeInvoiceForPrint.subtotal_medication > 0" class="flex justify-between">
              <span>Obat & Farmasi:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.subtotal_medication) }}</span>
            </div>
            <div v-if="activeInvoiceForPrint.voucher_discount > 0" class="flex justify-between text-purple-700 font-bold">
              <span>Diskon Voucher:</span>
              <span>-{{ formatRupiah(activeInvoiceForPrint.voucher_discount) }}</span>
            </div>
          </div>

          <div class="pt-2 border-t border-black/20 space-y-1 text-xs">
            <div class="flex justify-between font-black text-sm">
              <span>TOTAL BAYAR:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.total_amount) }}</span>
            </div>
            <div class="flex justify-between text-[11px] text-zinc-600">
              <span>Metode Bayar:</span>
              <span>{{ activeInvoiceForPrint.payment_method }}</span>
            </div>
            <div v-if="activeInvoiceForPrint.payment_method === 'CASH'" class="flex justify-between text-[11px]">
              <span>Uang Diterima:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.cash_received) }}</span>
            </div>
            <div v-if="activeInvoiceForPrint.payment_method === 'CASH'" class="flex justify-between text-[11px] font-bold text-emerald-700">
              <span>Kembalian:</span>
              <span>{{ formatRupiah(activeInvoiceForPrint.cash_change) }}</span>
            </div>
          </div>

          <div class="pt-3 border-t border-dashed border-black/20 text-center space-y-1 text-[10px] text-zinc-500">
            <div class="font-bold text-zinc-800">STATUS: LUNAS (TERVERIFIKASI)</div>
            <div>Semoga lekas sembuh & terima kasih atas kunjungan Anda.</div>
            <div class="font-mono text-[9px] text-zinc-400">Paperless Healthcare System</div>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            @click="printReceipt"
            class="flex-1 py-3 rounded-2xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-xs shadow-md transition-all active:scale-95 flex items-center justify-center gap-1.5"
          >
            <Printer class="w-3.5 h-3.5 text-emerald-400" />
            Cetak Struk Sekarang
          </button>
          <button
            @click="showPrintModal = false"
            class="px-4 py-3 rounded-2xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL 2: CATAT REFUND MANUAL -->
    <div
      v-if="showRefundModal && invoiceToRefund"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-3xl max-w-md w-full p-6 space-y-4 shadow-2xl border border-black/10">
        <div class="flex items-center justify-between pb-3 border-b border-black/10">
          <div>
            <h3 class="text-sm font-extrabold text-zinc-950">Pencatatan Refund Pasien (Manual)</h3>
            <p class="text-xs text-zinc-400">Invoice: {{ invoiceToRefund.invoice_number }} • Total: {{ formatRupiah(invoiceToRefund.total_amount) }}</p>
          </div>
          <button @click="showRefundModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3">
          <div class="space-y-1">
            <label class="text-xs font-bold text-zinc-700">Nama Bank Tujuan:</label>
            <select
              v-model="refundForm.refund_bank_name"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 text-xs font-medium focus:ring-2 focus:ring-zinc-950"
            >
              <option value="BCA">BCA (Bank Central Asia)</option>
              <option value="Mandiri">Bank Mandiri</option>
              <option value="BRI">Bank Rakyat Indonesia (BRI)</option>
              <option value="BNI">Bank Negara Indonesia (BNI)</option>
              <option value="BSI">Bank Syariah Indonesia (BSI)</option>
            </select>
          </div>

          <div class="space-y-1">
            <label class="text-xs font-bold text-zinc-700">Nomor Rekening Pengembalian:</label>
            <input
              v-model="refundForm.refund_bank_account"
              type="text"
              placeholder="Contoh: 1234567890"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 text-xs font-mono font-medium focus:ring-2 focus:ring-zinc-950"
            />
          </div>

          <div class="space-y-1">
            <label class="text-xs font-bold text-zinc-700">Atas Nama Pemilik Rekening:</label>
            <input
              v-model="refundForm.refund_account_name"
              type="text"
              placeholder="Contoh: Budi Santoso"
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 text-xs font-medium focus:ring-2 focus:ring-zinc-950"
            />
          </div>

          <div class="space-y-1">
            <label class="text-xs font-bold text-zinc-700">Alasan Pembatalan / Refund:</label>
            <textarea
              v-model="refundForm.refund_reason"
              rows="2"
              placeholder="Contoh: Pasien mendadak demam tinggi harus segera dirujuk ke IGD rumah sakit."
              class="w-full px-3 py-2 rounded-xl bg-zinc-50 border border-black/10 text-xs focus:ring-2 focus:ring-zinc-950"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            @click="submitRefund"
            :disabled="isSubmittingRefund"
            class="flex-1 py-3 rounded-2xl bg-rose-600 hover:bg-rose-700 text-white font-extrabold text-xs shadow-md transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-1.5"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            {{ isSubmittingRefund ? 'Memproses Refund...' : 'Konfirmasi Refund Dana' }}
          </button>
          <button
            @click="showRefundModal = false"
            class="px-4 py-3 rounded-2xl bg-zinc-100 hover:bg-zinc-200 text-zinc-700 font-bold text-xs"
          >
            Batal
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
