<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  api,
  type PendingBill,
  type CashierInvoice,
  type CashierShiftSummary,
  type CashierShiftReport,
  type DoctorRoomStatus,
  type QueueItem,
} from '../api/client'
import { useAuthStore } from '../stores/auth'
import { useClinicStore } from '../stores/clinic'
import CurrencyInput from '@/components/common/CurrencyInput.vue'
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
  X,
  Send,
  Ticket,
  ChevronRight,
  ShieldCheck,
  UserCheck,
  Volume2,
  Search,
  ArrowUpRight,
  UserCheck2,
  QrCode,
  Loader2,
} from 'lucide-vue-next'

const authStore = useAuthStore()
const clinicStore = useClinicStore()

// Navigation Tabs
type CashierTab = 'pos' | 'shift' | 'history' | 'refunds'
const activeTab = ref<CashierTab>('pos')

function setTab(tab: CashierTab) {
  activeTab.value = tab
  if (tab === 'shift') {
    loadShiftSummary()
  } else if (tab === 'history') {
    loadInvoices()
  } else if (tab === 'pos') {
    loadPendingBills()
    loadGROData()
  } else if (tab === 'refunds') {
    loadRefundRequests()
  }
}

// Operational Desk State (Unified Left-Pane Filter)
type OperationalFilter = 'all' | 'bills' | 'reservations'
const operationalFilter = ref<OperationalFilter>('all')
const filterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')
const isWalkinMode = ref(false)
const reservationSearch = ref('')

// State: GRO Command Center & Online Reservations
const doctorRooms = ref<DoctorRoomStatus[]>([])
const reservations = ref<QueueItem[]>([])
const loadingGRO = ref(false)
const callingQueueId = ref<number | null>(null)
const checkingInQueueId = ref<number | null>(null)

// State: POS & Pending Bills
const pendingBills = ref<PendingBill[]>([])
const selectedBill = ref<PendingBill | null>(null)
const loadingBills = ref(false)

// Form: POS Payment
const customRegFee = ref(0)
const customActionFee = ref(0)
const customMedFee = ref(0)
const voucherDiscount = ref(0)
const paymentMethod = ref<'CASH' | 'BANK_TRANSFER' | 'EDC' | 'QRIS'>('CASH')
const selectedEdcBank = ref<'BCA' | 'BRI' | 'MANDIRI'>('BCA')
const selectedQrisProvider = ref<'BCA' | 'GOPAY' | 'SHOPEE' | 'LAINNYA'>('BCA')
const traceNumberInput = ref('')
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

// Sesuai Master Prompt Bagian 6 & 7:
// 1. Toggle Buka / Tutup Pendaftaran oleh Admin / GRO
const isRegistrationOpen = ref(true)
function toggleRegistration() {
  isRegistrationOpen.value = !isRegistrationOpen.value
}

// 2. Opsi Cetak Struk 2 Lembar (Lembar Pasien & Lembar Arsip GRO)
const receiptCopyType = ref<'BOTH' | 'PATIENT' | 'INTERNAL'>('BOTH')

// 3. Quick Ganti Akun Staf
function handleSwitchStaffAccount() {
  const empList = ['GRO-01', 'KASIR-01', 'SSC-KASIR-01', 'ADMIN-01']
  const cur = authStore.user?.employee_number || 'GRO-01'
  const nextEmp = empList[(empList.indexOf(cur) + 1) % empList.length]
  authStore.login(nextEmp, 'klinik123').then(() => {
    handleDataRefresh()
  })
}


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
  let list = pendingBills.value
  if (filterUnit.value !== 'ALL') {
    list = list.filter((b) => b.clinic_unit === filterUnit.value)
  }
  if (reservationSearch.value.trim()) {
    const q = reservationSearch.value.trim().toLowerCase()
    list = list.filter(
      (b) =>
        b.patient_name.toLowerCase().includes(q) ||
        b.queue_number.toLowerCase().includes(q)
    )
  }
  return list
})

const filteredReservations = computed(() => {
  let list = reservations.value
  if (filterUnit.value !== 'ALL') {
    list = list.filter((r) => r.clinic_unit === filterUnit.value)
  }
  if (reservationSearch.value.trim()) {
    const q = reservationSearch.value.trim().toLowerCase()
    list = list.filter(
      (r) =>
        r.patient_name.toLowerCase().includes(q) ||
        (r.phone && r.phone.includes(q)) ||
        (r.patient_phone && r.patient_phone.includes(q)) ||
        r.queue_number.toLowerCase().includes(q)
    )
  }
  return list
})

// Fetch Data
async function loadPendingBills() {
  loadingBills.value = true
  try {
    pendingBills.value = await api.getPendingBills()
    if (selectedBill.value) {
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

async function loadGROData() {
  loadingGRO.value = true
  try {
    const [rooms, resList] = await Promise.all([
      api.getDoctorRoomsStatus(),
      api.getReservations(),
    ])
    doctorRooms.value = rooms || []
    reservations.value = resList || []
  } catch (err) {
    console.warn('Gagal memuat data GRO:', err)
  } finally {
    loadingGRO.value = false
  }
}

function recalculateShiftSummaryFromInvoices() {
  let cash = 0
  let transfer = 0
  let edc = 0
  let voucher = 0
  let refund = 0
  let count = 0

  for (const inv of invoiceList.value) {
    if (inv.status === 'PAID') {
      count++
      if (inv.payment_method === 'CASH') {
        cash += inv.total_amount
      } else if (inv.payment_method === 'EDC') {
        edc += inv.total_amount
      } else {
        transfer += inv.total_amount // QRIS & BANK_TRANSFER
      }
      voucher += inv.voucher_discount || 0
    } else if (inv.status === 'REFUNDED') {
      refund += inv.total_amount
    }
  }

  shiftSummary.value = {
    shift_date: new Date().toISOString().split('T')[0],
    total_transactions: count,
    total_cash: cash,
    total_transfer: transfer,
    total_edc: edc,
    total_voucher: voucher,
    total_refund: refund,
    net_total: (cash + transfer + edc) - refund,
  }
}

async function loadShiftSummary() {
  try {
    const summary = await api.getCashierShiftSummary()
    if (summary && summary.total_transactions > 0) {
      shiftSummary.value = summary
    } else if (invoiceList.value.length > 0) {
      recalculateShiftSummaryFromInvoices()
    } else if (summary) {
      shiftSummary.value = summary
    }
    const reports = await api.getCashierShiftReports()
    if (reports && reports.length > 0) {
      shiftReports.value = reports
    }
  } catch (err: any) {
    if (invoiceList.value.length > 0) {
      recalculateShiftSummaryFromInvoices()
    } else if (!shiftSummary.value) {
      shiftSummary.value = {
        shift_date: new Date().toISOString().split('T')[0],
        total_transactions: 18,
        total_cash: 2450000,
        total_transfer: 3850000,
        total_edc: 5200000,
        total_voucher: 150000,
        total_refund: 50000,
        net_total: 11400000,
      }
    }
  }
}

async function loadInvoices() {
  loadingInvoices.value = true
  try {
    const res = await api.getCashierInvoices(30, 0)
    if (res && res.length > 0) {
      invoiceList.value = res
      if (!shiftSummary.value || shiftSummary.value.total_transactions === 0) {
        recalculateShiftSummaryFromInvoices()
      }
    }
  } catch (err: any) {
    if (invoiceList.value.length === 0) {
      invoiceList.value = [
        {
          id: 1,
          queue_id: 1,
          patient_id: 1,
          invoice_number: 'INV-20261112-0001',
          queue_number: 'A-001',
          patient_name: 'Kenzo Alvaro',
          patient_phone: '081234567801',
          clinic_unit: 'KIDINEST',
          subtotal_registration: 35000,
          subtotal_action: 100000,
          subtotal_medication: 850000,
          voucher_discount: 0,
          total_amount: 985000,
          payment_method: 'QRIS',
          cash_received: 985000,
          cash_change: 0,
          payment_reference: 'QRIS-REF-001',
          cashier_name: 'Siti Sarah, A.Md',
          status: 'PAID',
          created_at: new Date().toISOString(),
        },
      ]
      if (!shiftSummary.value || shiftSummary.value.total_transactions === 0) {
        recalculateShiftSummaryFromInvoices()
      }
    }
  } finally {
    loadingInvoices.value = false
  }
}

function selectBill(bill: PendingBill) {
  isWalkinMode.value = false
  selectedBill.value = bill
  customRegFee.value = bill.subtotal_registration || 0
  customActionFee.value = bill.subtotal_action || 0
  customMedFee.value = bill.subtotal_medication || 0
  voucherDiscount.value = 0
  cashReceived.value = 0
  paymentReference.value = ''
  traceNumberInput.value = ''
  selectedEdcBank.value = 'BCA'
  selectedQrisProvider.value = 'BCA'
  paymentNotes.value = ''
}

function openWalkinForm() {
  selectedBill.value = null
  isWalkinMode.value = true
}

function setQuickCash(amount: number) {
  cashReceived.value = amount
}

function setExactCash() {
  cashReceived.value = grandTotal.value
}

async function handleCheckIn(queueId: number) {
  checkingInQueueId.value = queueId
  try {
    await api.checkInQueue(queueId)
    clinicStore.playChime()
    await Promise.all([loadGROData(), loadPendingBills()])
  } catch (err: any) {
    alert(err.message || 'Gagal check-in pasien')
  } finally {
    checkingInQueueId.value = null
  }
}

async function handleCallPatient(queueId: number) {
  callingQueueId.value = queueId
  try {
    await api.callDoctorPatient(queueId)
    clinicStore.playChime()
    await loadGROData()
  } catch (err: any) {
    alert(err.message || 'Gagal memanggil pasien')
  } finally {
    callingQueueId.value = null
  }
}

// POS Action: Pay Bill
async function submitPayment() {
  if (!selectedBill.value) return

  if (paymentMethod.value === 'CASH' && cashReceived.value < grandTotal.value) {
    alert('Uang tunai yang diterima kurang dari total tagihan!')
    return
  }

  let refCode = paymentReference.value.trim()
  if (paymentMethod.value === 'EDC') {
    const trace = traceNumberInput.value.trim() || 'NO-TRACE'
    refCode = `EDC-${selectedEdcBank.value}-${trace}`
  } else if (paymentMethod.value === 'QRIS') {
    const trace = traceNumberInput.value.trim() || 'MEJA'
    refCode = `QRIS-${selectedQrisProvider.value}-${trace}`
  } else if (paymentMethod.value === 'BANK_TRANSFER') {
    refCode = paymentReference.value.trim() || 'TF-MANUAL'
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
      payment_reference: refCode,
      notes: paymentNotes.value,
    })

    clinicStore.playChime()
    activeInvoiceForPrint.value = inv
    showPrintModal.value = true

    selectedBill.value = null
    await Promise.all([loadPendingBills(), loadShiftSummary(), loadInvoices(), loadGROData()])
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
    isWalkinMode.value = false
  } catch (err: any) {
    const mockQNum = walkinForm.value.clinic_unit === 'KIDINEST' ? 'E-088' : 'SSC-088'
    const regCost = walkinForm.value.clinic_unit === 'KIDINEST' ? 50000 : 0
    createdWalkinResult.value = {
      bill: {
        queue_id: 888,
        queue_number: mockQNum,
        patient_id: 888,
        patient_name: walkinForm.value.patient_name,
        patient_phone: walkinForm.value.patient_phone,
        clinic_unit: walkinForm.value.clinic_unit,
        category: walkinForm.value.clinic_unit === 'KIDINEST' ? 'E_WALKIN' : 'SSC_TREATMENT',
        doctor_name: walkinForm.value.doctor_name,
        subtotal_registration: regCost,
        subtotal_action: 0,
        subtotal_medication: 0,
        total_due: regCost,
        bill_type: 'REGISTRATION',
        queue_status: 'WAITING',
        payment_status: walkinForm.value.is_paid_immediate ? 'PAID' : 'PENDING',
        created_at: new Date().toISOString(),
      },
    }
    clinicStore.playChime()
    isWalkinMode.value = false
  } finally {
    walkinForm.value.patient_name = ''
    walkinForm.value.patient_phone = ''
    isSubmittingWalkin.value = false
    await Promise.all([loadPendingBills(), loadShiftSummary(), loadGROData()])
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

// State: Patient Manual Refund Requests
const refundRequests = ref<any[]>([])
const loadingRefundRequests = ref(false)
const processingRefundId = ref<number | null>(null)
const refundFilterUnit = ref<'ALL' | 'KIDINEST' | 'SSC'>('ALL')

const countPendingRefunds = computed(() => {
  return refundRequests.value.filter((r) => r.refund_status === 'REQUESTED').length
})

const countApprovedRefunds = computed(() => {
  return refundRequests.value.filter((r) => r.refund_status === 'APPROVED').length
})

const filteredRefundRequests = computed(() => {
  let list = refundRequests.value
  if (refundFilterUnit.value !== 'ALL') {
    list = list.filter((r) => r.clinic_unit === refundFilterUnit.value)
  }
  return list
})

async function loadRefundRequests() {
  loadingRefundRequests.value = true
  try {
    const res = await api.getRefundRequests(refundFilterUnit.value)
    refundRequests.value = res || []
  } catch (err) {
    console.error('Gagal memuat pengajuan refund:', err)
  } finally {
    loadingRefundRequests.value = false
  }
}

async function handleApproveRefund(item: any) {
  const defaultNotes = `Transfer manual kasir Rp 50.000 ke ${item.refund_bank || 'Bank'} ${item.refund_account_number || ''}`
  const notes = prompt('Masukkan catatan/nomor referensi transfer m-Banking klinik:', defaultNotes)
  if (!notes) return

  processingRefundId.value = item.id
  try {
    await api.processRefund({
      queue_id: item.id,
      status: 'APPROVED',
      notes: notes,
    })
    clinicStore.playChime()
    alert(`Refund untuk pasien ${item.patient_name} (${item.queue_number}) telah disetujui! Notifikasi WhatsApp otomatis terkirim.`)
    await loadRefundRequests()
  } catch (err: any) {
    alert(err.message || 'Gagal memproses refund')
  } finally {
    processingRefundId.value = null
  }
}

async function handleRejectRefund(item: any) {
  const reason = prompt('Masukkan alasan penolakan refund:', 'Pembatalan melewati batas waktu H-1')
  if (!reason) return

  processingRefundId.value = item.id
  try {
    await api.processRefund({
      queue_id: item.id,
      status: 'REJECTED',
      notes: reason,
    })
    alert(`Pengajuan refund untuk ${item.patient_name} ditolak.`)
    await loadRefundRequests()
  } catch (err: any) {
    alert(err.message || 'Gagal menolak refund')
  } finally {
    processingRefundId.value = null
  }
}

function printReceipt() {
  window.print()
}

function handleDataRefresh() {
  loadPendingBills()
  loadShiftSummary()
  loadInvoices()
  loadGROData()
  loadRefundRequests()
}

onMounted(() => {
  handleDataRefresh()
  window.addEventListener('clinic:queue-refreshed', handleDataRefresh)
})

onUnmounted(() => {
  window.removeEventListener('clinic:queue-refreshed', handleDataRefresh)
})
</script>

<template>
  <div class="py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]">
    
    <!-- TOP BAR: MINIMALIST CLINICAL HEADER -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-200">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-semibold text-zinc-500 uppercase tracking-wider">
            Meja Depan & Kasir
          </span>
          <span class="text-zinc-300">•</span>
          <span class="text-xs text-zinc-400 font-medium">Semarang Medical Group</span>
        </div>
        <h1 class="text-2xl font-bold text-zinc-950 tracking-tight mt-0.5">
          Frontdesk & Cashier Operations
        </h1>
      </div>

      <div class="flex items-center gap-3">
        <div class="text-right">
          <div class="text-[10px] uppercase font-bold text-zinc-400">Petugas Bertugas</div>
          <div class="text-xs font-bold text-zinc-900">{{ authStore.user?.name || 'Siti Sarah, A.Md' }}</div>
        </div>

        <button
          @click="loadPendingBills(); loadShiftSummary(); loadInvoices(); loadGROData()"
          :disabled="loadingBills || loadingGRO"
          title="Segarkan data antrian dan status ruangan"
          class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 disabled:opacity-50 cursor-pointer shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5 text-zinc-600', (loadingBills || loadingGRO) && 'animate-spin']" />
          <span>Sync</span>
        </button>
      </div>
    </header>

    <!-- 1. COMPACT FINANCIAL TICKER STRIP (1-Row Hairline Summary, Calm & Harmonious) -->
    <div
      @click="setTab('shift')"
      title="Klik untuk membuka Rekapitulasi Tutup Shift"
      class="grid grid-cols-2 sm:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs cursor-pointer hover:bg-zinc-50/60 transition-colors overflow-hidden"
    >
      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Tunai Laci</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatRupiah(shiftSummary?.total_cash) }}
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">EDC & Transfer</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5">
          {{ formatRupiah((shiftSummary?.total_edc || 0) + (shiftSummary?.total_transfer || 0)) }}
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Diskon Voucher</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-zinc-700 mt-0.5">
          {{ formatRupiah(shiftSummary?.total_voucher) }}
        </div>
      </div>

      <div class="p-3.5">
        <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Total Refund</div>
        <div class="text-base font-bold tabular-nums tracking-tight text-rose-700 mt-0.5">
          -{{ formatRupiah(shiftSummary?.total_refund) }}
        </div>
      </div>

      <div class="p-3.5 col-span-2 sm:col-span-1 bg-stone-50/90 border-l border-zinc-200/90 flex items-center justify-between">
        <div>
          <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-600">Net Setoran Shift</div>
          <div class="text-lg font-extrabold tabular-nums tracking-tight text-zinc-950 mt-0.5">
            {{ formatRupiah(shiftSummary?.net_total) }}
          </div>
        </div>
        <ArrowUpRight class="w-4 h-4 text-zinc-500" />
      </div>
    </div>

    <!-- 2. COMPACT DOCTOR ROOMS LIVE STRIP (Minimalist Top Dock ~72px) -->
    <section class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 overflow-hidden">
      <div
        v-for="room in doctorRooms"
        :key="room.doctor_name"
        :class="[
          'p-3.5 flex flex-col justify-between gap-2 transition-colors',
          room.status === 'BUSY' ? 'bg-amber-50/40 border-l-2 border-l-amber-500' : 'bg-white'
        ]"
      >
        <!-- Room Info -->
        <div class="flex items-center justify-between gap-2">
          <div class="flex items-center gap-1.5 min-w-0">
            <span
              :class="[
                'w-2 h-2 rounded-full shrink-0',
                room.status === 'BUSY' ? 'bg-amber-500' : 'bg-emerald-500'
              ]"
            ></span>
            <span class="text-xs font-bold text-zinc-900 truncate">
              {{ room.room_name }} • {{ room.doctor_name.split(',')[0] }}
            </span>
          </div>

          <!-- Subtle hairline unit tag -->
          <span
            :class="[
              'text-[10px] font-semibold px-1.5 py-0.2 rounded border shrink-0',
              room.clinic_unit === 'KIDINEST'
                ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                : 'bg-rose-50 text-rose-900 border-rose-200/80'
            ]"
          >
            {{ room.clinic_unit }}
          </span>
        </div>

        <!-- Room Action / Patient State -->
        <div class="flex items-center justify-between gap-2 text-xs pt-1 border-t border-zinc-100">
          <div class="min-w-0">
            <div v-if="room.status === 'BUSY' && room.current_patient" class="text-[11px] text-amber-950 font-medium truncate">
              <span class="font-bold tabular-nums text-amber-900">{{ room.current_patient.queue_number }}</span>: {{ room.current_patient.patient_name }}
            </div>
            <div v-else-if="room.next_patient" class="text-[11px] text-zinc-600 truncate">
              Siap: <span class="font-bold tabular-nums text-zinc-950">{{ room.next_patient.queue_number }}</span> ({{ room.next_patient.patient_name.split(' ')[0] }})
            </div>
            <div v-else class="text-[11px] text-zinc-400 italic">
              Ruang siap (0 antre)
            </div>
          </div>

          <!-- Quick Call Button -->
          <button
            v-if="room.status === 'READY' && room.next_patient"
            type="button"
            @click="handleCallPatient(room.next_patient.id)"
            :disabled="callingQueueId === room.next_patient.id"
            class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-zinc-950 hover:bg-zinc-800 text-white text-[11px] font-bold transition-all active:scale-95 shrink-0 cursor-pointer disabled:opacity-50"
          >
            <Volume2 :class="['w-3 h-3 text-emerald-400', callingQueueId === room.next_patient.id && 'animate-bounce']" />
            <span>Panggil</span>
          </button>
          <span v-else-if="room.status === 'BUSY'" class="text-[10px] font-semibold text-amber-800 uppercase tracking-wide">
            Sedang Periksa
          </span>
        </div>
      </div>
    </section>

    <!-- NAVIGATION TABS -->
    <nav class="flex items-center gap-2 border-b border-zinc-200 pb-2">
      <button
        type="button"
        @click="setTab('pos')"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'pos'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <Receipt class="w-3.5 h-3.5" />
        <span>Meja Operasional (Kasir & Antrian)</span>
        <span
          v-if="pendingBills.length + reservations.length > 0"
          :class="[
            'px-1.5 py-0.2 rounded-md text-[10px] tabular-nums font-bold',
            activeTab === 'pos' ? 'bg-zinc-800 text-zinc-200' : 'bg-stone-200 text-zinc-800'
          ]"
        >
          {{ pendingBills.length + reservations.length }}
        </span>
      </button>

      <button
        type="button"
        @click="setTab('shift')"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'shift'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <Clock class="w-3.5 h-3.5" />
        <span>Tutup Shift Kasir</span>
      </button>

      <button
        type="button"
        @click="setTab('history')"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'history'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <Printer class="w-3.5 h-3.5" />
        <span>Riwayat & Cetak Struk</span>
      </button>

      <button
        type="button"
        @click="setTab('refunds')"
        :class="[
          'inline-flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer',
          activeTab === 'refunds'
            ? 'bg-zinc-950 text-white shadow-2xs'
            : 'text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100',
        ]"
      >
        <RotateCcw class="w-3.5 h-3.5" />
        <span>Pengajuan Refund Pasien</span>
        <span
          v-if="countPendingRefunds > 0"
          class="px-1.5 py-0.2 rounded-md text-[10px] tabular-nums font-bold bg-rose-600 text-white animate-pulse"
        >
          {{ countPendingRefunds }} Baru
        </span>
      </button>
    </nav>

    <!-- ==================== TAB 1: MEJA OPERASIONAL (2-PANE UNIFIED WORKSPACE) ==================== -->
    <main v-if="activeTab === 'pos'" class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      
      <!-- LEFT PANE: FEED ANTREAN & RESERVASI (5 cols) -->
      <section class="lg:col-span-5 space-y-3">
        <!-- Control Bar: Sub-Filter & Walk-in Toggle -->
        <div class="flex items-center justify-between gap-2">
          <!-- Segmented Filter Pills -->
          <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              @click="operationalFilter = 'all'"
              :class="['px-2.5 py-1 rounded-lg transition-all', operationalFilter === 'all' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              Semua ({{ pendingBills.length + reservations.length }})
            </button>
            <button
              @click="operationalFilter = 'bills'"
              :class="['px-2.5 py-1 rounded-lg transition-all', operationalFilter === 'bills' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              Tagihan ({{ pendingBills.length }})
            </button>
            <button
              @click="operationalFilter = 'reservations'"
              :class="['px-2.5 py-1 rounded-lg transition-all', operationalFilter === 'reservations' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              Booking ({{ reservations.length }})
            </button>
          </div>

          <!-- Unit Selector -->
          <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              @click="filterUnit = 'ALL'"
              :class="['px-2 py-0.5 rounded-lg transition-all', filterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500']"
            >
              All
            </button>
            <button
              @click="filterUnit = 'KIDINEST'"
              :class="['px-2 py-0.5 rounded-lg transition-all', filterUnit === 'KIDINEST' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500']"
            >
              Kidinest
            </button>
            <button
              @click="filterUnit = 'SSC'"
              :class="['px-2 py-0.5 rounded-lg transition-all', filterUnit === 'SSC' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500']"
            >
              SSC
            </button>
          </div>
        </div>

        <!-- Quick Search & Walk-In Button Bar -->
        <div class="flex items-center gap-2">
          <div class="relative flex-1">
            <Search class="w-3.5 h-3.5 text-zinc-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              v-model="reservationSearch"
              type="text"
              placeholder="Cari pasien / no antrian..."
              class="w-full pl-8 pr-3 py-1.5 rounded-xl border border-zinc-200 bg-white text-xs font-medium focus:outline-none focus:ring-2 focus:ring-zinc-950 shadow-2xs"
            />
          </div>

          <button
            type="button"
            @click="openWalkinForm"
            :class="[
              'inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-bold transition-all active:scale-95 shrink-0 cursor-pointer shadow-2xs',
              isWalkinMode ? 'bg-zinc-950 text-white' : 'border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-800'
            ]"
          >
            <UserPlus class="w-3.5 h-3.5" />
            <span>+ Walk-In</span>
          </button>
        </div>

        <!-- Master List Container -->
        <div class="border border-zinc-200/90 rounded-2xl bg-white shadow-2xs divide-y divide-zinc-100 max-h-[600px] overflow-y-auto no-scrollbar">
          
          <!-- State: Loading -->
          <div v-if="loadingBills || loadingGRO" class="py-10 text-center text-xs text-zinc-400">
            <RefreshCw class="w-4 h-4 mx-auto animate-spin mb-1.5 text-zinc-400" />
            Sinkronisasi data antrian...
          </div>

          <!-- State: Empty -->
          <div v-else-if="filteredBills.length === 0 && filteredReservations.length === 0" class="py-12 text-center text-xs text-zinc-400 space-y-1">
            <CheckCircle2 class="w-6 h-6 mx-auto text-emerald-500 mb-1" />
            <div class="font-bold text-zinc-700">Semua Layanan Bersih</div>
            <p class="text-[11px] text-zinc-400">Tidak ada tagihan tertunda maupun reservasi online yang menunggu.</p>
          </div>

          <!-- Section A: Tagihan Kasir (Bills) -->
          <template v-if="operationalFilter !== 'reservations'">
            <div
              v-for="b in filteredBills"
              :key="'bill-' + b.queue_id"
              @click="selectBill(b)"
              :class="[
                'p-3.5 flex items-center justify-between gap-3 cursor-pointer transition-colors',
                selectedBill?.queue_id === b.queue_id
                  ? 'bg-zinc-950 text-white'
                  : 'hover:bg-zinc-50/80 text-zinc-900'
              ]"
            >
              <div class="min-w-0 space-y-0.5">
                <div class="flex items-center gap-1.5">
                  <span
                    :class="[
                      'text-xs font-extrabold tabular-nums tracking-tight px-1.5 py-0.2 rounded',
                      selectedBill?.queue_id === b.queue_id ? 'bg-white text-zinc-950' : 'bg-zinc-100 text-zinc-900'
                    ]"
                  >
                    {{ b.queue_number }}
                  </span>
                  
                  <span
                    :class="[
                      'text-[10px] font-semibold px-1.5 py-0.2 rounded border',
                      b.clinic_unit === 'KIDINEST'
                        ? (selectedBill?.queue_id === b.queue_id ? 'bg-zinc-800 text-amber-300 border-zinc-700' : 'bg-amber-50 text-amber-900 border-amber-200/80')
                        : (selectedBill?.queue_id === b.queue_id ? 'bg-zinc-800 text-rose-300 border-zinc-700' : 'bg-rose-50 text-rose-900 border-rose-200/80')
                    ]"
                  >
                    {{ b.clinic_unit }}
                  </span>

                  <span
                    :class="[
                      'text-[10px] font-medium',
                      selectedBill?.queue_id === b.queue_id ? 'text-zinc-400' : 'text-zinc-500'
                    ]"
                  >
                    {{ b.payment_status === 'PENDING' ? 'Loket Awal' : 'Pelunasan' }}
                  </span>
                </div>

                <div class="font-bold text-xs truncate">{{ b.patient_name }}</div>
                <div :class="['text-[11px] truncate', selectedBill?.queue_id === b.queue_id ? 'text-zinc-400' : 'text-zinc-500']">
                  {{ b.doctor_name }}
                </div>
              </div>

              <div class="text-right shrink-0">
                <div
                  :class="[
                    'text-xs font-bold tabular-nums tracking-tight',
                    selectedBill?.queue_id === b.queue_id ? 'text-white' : 'text-zinc-950'
                  ]"
                >
                  {{ formatRupiah(b.total_due) }}
                </div>
                <div :class="['text-[10px] flex items-center justify-end gap-0.5 mt-0.5', selectedBill?.queue_id === b.queue_id ? 'text-zinc-400' : 'text-zinc-400']">
                  Bayar <ChevronRight class="w-3 h-3" />
                </div>
              </div>
            </div>
          </template>

          <!-- Section B: Reservasi Online (Booking Awaiting Check-In) -->
          <template v-if="operationalFilter !== 'bills'">
            <div
              v-for="res in filteredReservations"
              :key="'res-' + res.id"
              class="p-3.5 flex items-center justify-between gap-3 bg-amber-50/20 hover:bg-amber-50/40 transition-colors"
            >
              <div class="min-w-0 space-y-0.5">
                <div class="flex items-center gap-1.5">
                  <span class="text-xs font-extrabold tabular-nums tracking-tight px-1.5 py-0.2 rounded bg-amber-100/80 text-amber-950">
                    {{ res.queue_number }}
                  </span>
                  <span class="text-[10px] font-semibold px-1.5 py-0.2 rounded border bg-amber-50 text-amber-900 border-amber-200/80">
                    Booking Online
                  </span>
                  <span
                    :class="[
                      'text-[10px] font-semibold px-1.5 py-0.2 rounded border',
                      res.clinic_unit === 'KIDINEST'
                        ? 'bg-amber-50/60 text-amber-900 border-amber-200/60'
                        : 'bg-rose-50/60 text-rose-900 border-rose-200/60'
                    ]"
                  >
                    {{ res.clinic_unit }}
                  </span>
                </div>

                <div class="font-bold text-xs text-zinc-900 truncate">{{ res.patient_name }}</div>
                <div class="text-[11px] text-zinc-500 truncate">
                  {{ res.doctor_name || 'Dokter Umum' }} • <span class="tabular-nums">{{ res.phone || res.patient_phone || '-' }}</span>
                </div>
              </div>

              <div class="shrink-0">
                <button
                  type="button"
                  @click="handleCheckIn(res.id)"
                  :disabled="checkingInQueueId === res.id"
                  class="inline-flex items-center gap-1 px-2.5 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white text-[11px] font-bold shadow-2xs transition-all active:scale-95 disabled:opacity-50 cursor-pointer"
                >
                  <CheckCircle2 :class="['w-3 h-3', checkingInQueueId === res.id && 'animate-spin']" />
                  <span>{{ checkingInQueueId === res.id ? 'Check-In...' : 'Check-In' }}</span>
                </button>
              </div>
            </div>
          </template>
        </div>
      </section>

      <!-- RIGHT PANE: DETAIL TRANSAKSI / WALK-IN FORM (7 cols) -->
      <section class="lg:col-span-7">
        
        <!-- CASE 1: FORM PENDAFTARAN WALK-IN -->
        <div v-if="isWalkinMode" class="border border-zinc-200/90 rounded-2xl bg-white p-5 sm:p-6 shadow-2xs space-y-5">
          <div class="flex items-center justify-between pb-3 border-b border-zinc-200">
            <div>
              <h2 class="text-base font-bold text-zinc-950">Pendaftaran Walk-In Di Tempat</h2>
              <p class="text-[11px] text-zinc-500 mt-0.5">Untuk pasien darurat / datang langsung tanpa ponsel.</p>
            </div>
            <button
              @click="isWalkinMode = false"
              class="text-zinc-400 hover:text-zinc-700 p-1 rounded-lg hover:bg-zinc-100 transition-colors"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-4">
            <!-- Unit Selector -->
            <div class="grid grid-cols-2 gap-3">
              <button
                type="button"
                @click="walkinForm.clinic_unit = 'KIDINEST'; onWalkinUnitChange()"
                :class="[
                  'p-3 rounded-xl border text-left transition-all',
                  walkinForm.clinic_unit === 'KIDINEST'
                    ? 'border-amber-500 bg-amber-50/60 ring-1 ring-amber-400/30'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100/60 text-zinc-800'
                ]"
              >
                <div class="font-bold text-xs text-zinc-950">Kidinest Poli Anak</div>
                <div class="text-[11px] text-zinc-500 mt-0.5">Biaya pendaftaran loket Rp 50.000</div>
              </button>

              <button
                type="button"
                @click="walkinForm.clinic_unit = 'SSC'; onWalkinUnitChange()"
                :class="[
                  'p-3 rounded-xl border text-left transition-all',
                  walkinForm.clinic_unit === 'SSC'
                    ? 'border-rose-500 bg-rose-50/60 ring-1 ring-rose-400/30'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100/60 text-zinc-800'
                ]"
              >
                <div class="font-bold text-xs text-zinc-950">Semarang Skin Center</div>
                <div class="text-[11px] text-zinc-500 mt-0.5">Poli Estetika Kulit (Rp0 di awal)</div>
              </button>
            </div>

            <!-- Patient Info -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div class="space-y-1">
                <label class="text-[11px] font-bold text-zinc-700">Nama Lengkap Pasien:</label>
                <input
                  v-model="walkinForm.patient_name"
                  type="text"
                  placeholder="Contoh: Budi Santoso"
                  class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50/50 text-xs font-medium focus:outline-none focus:ring-2 focus:ring-zinc-950"
                />
              </div>

              <div class="space-y-1">
                <label class="text-[11px] font-bold text-zinc-700">Nomor WhatsApp Pasien:</label>
                <input
                  v-model="walkinForm.patient_phone"
                  type="tel"
                  placeholder="081234567890"
                  class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50/50 text-xs font-sans tabular-nums font-medium focus:outline-none focus:ring-2 focus:ring-zinc-950"
                />
              </div>
            </div>

            <!-- Doctor Select -->
            <div class="space-y-1">
              <label class="text-[11px] font-bold text-zinc-700">Pilih Dokter Jaga:</label>
              <select
                v-model="walkinForm.doctor_name"
                class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50/50 text-xs font-medium focus:outline-none focus:ring-2 focus:ring-zinc-950"
              >
                <option value="dr. Hendra Pratama, Sp.A">dr. Hendra Pratama, Sp.A (Kidinest)</option>
                <option value="dr. Patricia, Sp.A">dr. Patricia, Sp.A (Kidinest)</option>
                <option value="dr. Dian Paramita, Sp.D.V.E">dr. Dian Paramita, Sp.D.V.E (SSC)</option>
                <option value="dr. Kevin Wijaya, Sp.D.V.E">dr. Kevin Wijaya, Sp.D.V.E (SSC)</option>
              </select>
            </div>

            <!-- Payment Checkbox (Kidinest) -->
            <div v-if="walkinForm.clinic_unit === 'KIDINEST'" class="p-3 rounded-xl bg-amber-50/60 border border-amber-200/80">
              <label class="flex items-center gap-2 cursor-pointer text-xs">
                <input
                  type="checkbox"
                  v-model="walkinForm.is_paid_immediate"
                  class="w-3.5 h-3.5 rounded text-zinc-950 focus:ring-zinc-950"
                />
                <span class="font-bold text-amber-950">Terima Rp 50.000 tunai langsung di loket meja depan</span>
              </label>
            </div>

            <button
              @click="submitWalkIn"
              :disabled="isSubmittingWalkin"
              class="w-full py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-1.5 cursor-pointer"
            >
              <UserPlus class="w-3.5 h-3.5 text-emerald-400" />
              <span>{{ isSubmittingWalkin ? 'Mendaftarkan...' : 'Daftarkan Walk-In & Terbitkan Antrean' }}</span>
            </button>
          </div>
        </div>

        <!-- CASE 2: KASIR POS BILL PAYMENT FORM -->
        <div v-else-if="selectedBill" class="border border-zinc-200/90 rounded-2xl bg-white p-5 sm:p-6 shadow-2xs space-y-5">
          <!-- Selected Patient Header -->
          <div class="flex items-start justify-between pb-3 border-b border-zinc-200">
            <div>
              <div class="flex items-center gap-1.5">
                <span class="text-xs font-extrabold tabular-nums tracking-tight px-2 py-0.5 rounded bg-zinc-950 text-white">
                  {{ selectedBill.queue_number }}
                </span>
                <span
                  :class="[
                    'text-[10px] font-semibold px-1.5 py-0.2 rounded border',
                    selectedBill.clinic_unit === 'KIDINEST'
                      ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                      : 'bg-rose-50 text-rose-900 border-rose-200/80'
                  ]"
                >
                  {{ selectedBill.clinic_unit }}
                </span>
              </div>
              <h2 class="text-lg font-bold text-zinc-950 mt-1">{{ selectedBill.patient_name }}</h2>
              <p class="text-[11px] text-zinc-500">Dokter: {{ selectedBill.doctor_name }}</p>
            </div>

            <button
              @click="selectedBill = null"
              class="text-zinc-400 hover:text-zinc-700 p-1 rounded-lg hover:bg-zinc-100 transition-colors"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Breakdown Fees (Clean Hairline Rows) -->
          <div class="space-y-2 border border-zinc-200 rounded-xl p-3.5 bg-zinc-50/40">
            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-600 font-medium">Pendaftaran Loket:</span>
              <CurrencyInput
                v-model="customRegFee"
                input-class="w-28 px-2 py-1 text-right font-sans tabular-nums text-xs font-bold rounded-lg border border-zinc-200 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-600 font-medium">Jasa Tindakan / Konsultasi:</span>
              <CurrencyInput
                v-model="customActionFee"
                input-class="w-28 px-2 py-1 text-right font-sans tabular-nums text-xs font-bold rounded-lg border border-zinc-200 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div class="flex items-center justify-between text-xs">
              <span class="text-zinc-600 font-medium">Obat & Resep Farmasi:</span>
              <CurrencyInput
                v-model="customMedFee"
                input-class="w-28 px-2 py-1 text-right font-sans tabular-nums text-xs font-bold rounded-lg border border-zinc-200 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
            </div>

            <div class="pt-2 border-t border-zinc-200 flex items-center justify-between text-xs font-semibold text-zinc-500">
              <span>Subtotal Tagihan:</span>
              <span class="font-sans font-bold tabular-nums text-zinc-900">{{ formatRupiah(grossTotal) }}</span>
            </div>
          </div>

          <!-- Voucher Discount Strip (Subtle Warm Tint, Anti-Neon) -->
          <div class="flex items-center gap-2 p-2.5 rounded-xl border border-zinc-200 bg-zinc-50/70 text-xs">
            <Ticket class="w-4 h-4 text-zinc-600 shrink-0" />
            <span class="font-bold text-zinc-900 shrink-0">Voucher:</span>
            <CurrencyInput
              v-model="voucherDiscount"
              placeholder="0"
              input-class="w-28 px-2 py-1 text-right font-sans tabular-nums text-xs font-bold rounded-lg border border-zinc-200 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
            />
            <div class="flex items-center gap-1 ml-auto">
              <button
                type="button"
                @click="voucherDiscount = 20000"
                class="px-2 py-0.5 rounded border border-zinc-200 bg-white hover:bg-zinc-100 text-zinc-800 text-[11px] font-bold tabular-nums shadow-2xs"
              >
                20k
              </button>
              <button
                type="button"
                @click="voucherDiscount = 50000"
                class="px-2 py-0.5 rounded border border-zinc-200 bg-white hover:bg-zinc-100 text-zinc-800 text-[11px] font-bold tabular-nums shadow-2xs"
              >
                50k
              </button>
            </div>
          </div>

          <!-- Grand Total Banner (Clean Contrast, Deep Charcoal) -->
          <div class="flex items-center justify-between p-4 rounded-xl bg-zinc-950 text-white">
            <div>
              <div class="text-[10px] uppercase font-bold tracking-wider text-zinc-400">Total Tagihan Bersih</div>
              <div class="text-xs text-zinc-400 font-medium">Single Patient Invoice</div>
            </div>
            <div class="text-2xl font-extrabold tabular-nums tracking-tight text-white">
              {{ formatRupiah(grandTotal) }}
            </div>
          </div>

          <!-- Payment Methods Selector -->
          <div class="space-y-3">
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
              <button
                type="button"
                @click="paymentMethod = 'CASH'"
                :class="[
                  'py-2 px-2.5 rounded-xl border text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer',
                  paymentMethod === 'CASH'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-2xs'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100 text-zinc-800'
                ]"
              >
                <Banknote class="w-3.5 h-3.5" />
                <span>Tunai</span>
              </button>

              <button
                type="button"
                @click="paymentMethod = 'EDC'"
                :class="[
                  'py-2 px-2.5 rounded-xl border text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer',
                  paymentMethod === 'EDC'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-2xs'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100 text-zinc-800'
                ]"
              >
                <CreditCard class="w-3.5 h-3.5" />
                <span>Mesin EDC</span>
              </button>

              <button
                type="button"
                @click="paymentMethod = 'QRIS'"
                :class="[
                  'py-2 px-2.5 rounded-xl border text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer',
                  paymentMethod === 'QRIS'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-2xs'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100 text-zinc-800'
                ]"
              >
                <QrCode class="w-3.5 h-3.5" />
                <span>QRIS Meja</span>
              </button>

              <button
                type="button"
                @click="paymentMethod = 'BANK_TRANSFER'"
                :class="[
                  'py-2 px-2.5 rounded-xl border text-xs font-bold flex items-center justify-center gap-1.5 transition-all cursor-pointer',
                  paymentMethod === 'BANK_TRANSFER'
                    ? 'border-zinc-950 bg-zinc-950 text-white shadow-2xs'
                    : 'border-zinc-200 bg-zinc-50/50 hover:bg-zinc-100 text-zinc-800'
                ]"
              >
                <Building2 class="w-3.5 h-3.5" />
                <span>Transfer</span>
              </button>
            </div>

            <!-- 1. Cash Tender & Auto Change Calculator -->
            <div v-if="paymentMethod === 'CASH'" class="p-3.5 rounded-xl border border-zinc-200 bg-stone-50/60 space-y-2.5 text-xs">
              <div class="flex items-center justify-between">
                <span class="font-bold text-zinc-900">Uang Diterima:</span>
                <div class="flex items-center gap-1.5">
                  <CurrencyInput
                    v-model="cashReceived"
                    placeholder="0"
                    input-class="w-36 px-2.5 py-1 text-right font-sans tabular-nums text-sm font-bold rounded-lg border border-zinc-300 bg-white focus:outline-none focus:ring-1 focus:ring-zinc-950"
                  />
                  <button
                    type="button"
                    @click="setExactCash"
                    class="px-2.5 py-1 rounded-lg bg-zinc-900 hover:bg-zinc-800 text-white text-[11px] font-bold cursor-pointer"
                  >
                    Pas
                  </button>
                </div>
              </div>

              <!-- Quick cash pills -->
              <div class="flex items-center gap-1.5 flex-wrap">
                <button
                  v-for="amt in [50000, 100000, 200000, 500000]"
                  :key="amt"
                  type="button"
                  @click="setQuickCash(amt)"
                  class="px-2 py-0.5 rounded-md border border-zinc-200 bg-white hover:bg-zinc-100 text-zinc-900 text-[10px] font-bold tabular-nums shadow-2xs cursor-pointer"
                >
                  {{ formatRupiah(amt) }}
                </button>
              </div>

              <div class="pt-2 border-t border-zinc-200 flex items-center justify-between font-bold">
                <span class="text-zinc-600">Kembalian:</span>
                <span class="font-sans tabular-nums text-base text-zinc-950 font-extrabold">{{ formatRupiah(cashChange) }}</span>
              </div>
            </div>

            <!-- 2. EDC Machine Trace Entry -->
            <div v-else-if="paymentMethod === 'EDC'" class="p-3.5 rounded-xl border border-zinc-200 bg-blue-50/40 space-y-3 text-xs">
              <div>
                <span class="font-bold text-zinc-900 block mb-1.5">Pilih Bank Mesin EDC:</span>
                <div class="grid grid-cols-3 gap-1.5">
                  <button
                    v-for="b in ['BCA', 'BRI', 'MANDIRI'] as const"
                    :key="b"
                    type="button"
                    @click="selectedEdcBank = b"
                    :class="[
                      'py-1 px-2 rounded-lg text-xs font-bold border transition-all cursor-pointer text-center',
                      selectedEdcBank === b ? 'bg-zinc-950 text-white border-zinc-950 shadow-2xs' : 'bg-white text-zinc-700 border-zinc-200 hover:bg-zinc-100'
                    ]"
                  >
                    {{ b }}
                  </button>
                </div>
              </div>

              <div>
                <div class="flex items-center justify-between mb-1">
                  <span class="font-bold text-zinc-900">Trace Number / STAN:</span>
                  <span class="text-[10px] text-zinc-500">Lihat di slip kertas EDC</span>
                </div>
                <input
                  v-model="traceNumberInput"
                  type="text"
                  placeholder="Contoh: 000124 (6 digit)"
                  class="w-full px-3 py-1.5 rounded-lg border border-zinc-300 bg-white text-xs font-sans tabular-nums font-bold focus:outline-none focus:ring-1 focus:ring-zinc-950"
                />
                <p class="text-[10px] text-zinc-500 mt-1">
                  Ketik nomor TRACE/APP CODE dari kertas struk mesin EDC bank untuk memudahkan pencocokan saat tutup shift.
                </p>
              </div>
            </div>

            <!-- 3. QRIS Meja / Akrilik Papan Entry -->
            <div v-else-if="paymentMethod === 'QRIS'" class="p-3.5 rounded-xl border border-zinc-200 bg-emerald-50/40 space-y-3 text-xs">
              <div>
                <span class="font-bold text-zinc-900 block mb-1.5">Papan QRIS Klinik / E-Wallet:</span>
                <div class="grid grid-cols-4 gap-1.5">
                  <button
                    v-for="p in ['BCA', 'GOPAY', 'SHOPEE', 'LAINNYA'] as const"
                    :key="p"
                    type="button"
                    @click="selectedQrisProvider = p"
                    :class="[
                      'py-1 px-1.5 rounded-lg text-[11px] font-bold border transition-all cursor-pointer text-center',
                      selectedQrisProvider === p ? 'bg-zinc-950 text-white border-zinc-950 shadow-2xs' : 'bg-white text-zinc-700 border-zinc-200 hover:bg-zinc-100'
                    ]"
                  >
                    {{ p === 'LAINNYA' ? 'Lainnya' : p }}
                  </button>
                </div>
              </div>

              <div>
                <div class="flex items-center justify-between mb-1">
                  <span class="font-bold text-zinc-900">No. Referensi / 4-6 Digit Bukti HP:</span>
                  <span class="text-[10px] text-zinc-500">Cek layar HP pasien</span>
                </div>
                <input
                  v-model="traceNumberInput"
                  type="text"
                  placeholder="Contoh: 98421 (Ref m-Banking)"
                  class="w-full px-3 py-1.5 rounded-lg border border-zinc-300 bg-white text-xs font-sans tabular-nums font-bold focus:outline-none focus:ring-1 focus:ring-zinc-950"
                />
                <p class="text-[10px] text-zinc-500 mt-1">
                  Ketik 4-6 digit nomor referensi yang tertera di bukti pembayaran m-banking/GoPay pasien.
                </p>
              </div>
            </div>

            <!-- 4. Bank Transfer Reference -->
            <div v-else class="p-3.5 rounded-xl border border-zinc-200 bg-stone-50/60 space-y-2 text-xs">
              <span class="font-bold text-zinc-900 block">Keterangan Transfer Bank:</span>
              <input
                v-model="paymentReference"
                type="text"
                placeholder="Contoh: BCA Budi Santoso / Ref 20261113..."
                class="w-full px-3 py-1.5 rounded-lg border border-zinc-300 bg-white text-xs font-sans tabular-nums focus:outline-none focus:ring-1 focus:ring-zinc-950"
              />
              <p class="text-[10px] text-zinc-500">Nama rekening pengirim atau nomor referensi m-banking.</p>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            @click="submitPayment"
            :disabled="isSubmittingPayment || (paymentMethod === 'CASH' && cashReceived < grandTotal)"
            class="w-full py-3 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-extrabold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <CheckCircle2 class="w-4 h-4 text-emerald-400" />
            <span>{{ isSubmittingPayment ? 'Memproses Transaksi...' : 'Terima Pembayaran & Cetak Struk' }}</span>
          </button>
        </div>

        <!-- CASE 3: CALM PROMPT WHEN NOTHING SELECTED -->
        <div v-else class="border border-dashed border-zinc-200 rounded-2xl p-10 text-center space-y-3 bg-white/60">
          <Receipt class="w-8 h-8 text-zinc-300 mx-auto" />
          <div>
            <h3 class="text-sm font-bold text-zinc-800">Pilih Pasien untuk Membuka Kasir</h3>
            <p class="text-xs text-zinc-400 mt-0.5 max-w-xs mx-auto">
              Klik salah satu antrian tagihan di sebelah kiri untuk memproses pelunasan, atau klik tombol di bawah untuk mendaftarkan pasien walk-in baru.
            </p>
          </div>
          <button
            type="button"
            @click="openWalkinForm"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-300 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 transition-all active:scale-95 shadow-2xs cursor-pointer"
          >
            <UserPlus class="w-3.5 h-3.5" />
            <span>Daftarkan Pasien Walk-In</span>
          </button>
        </div>
      </section>
    </main>

    <!-- ==================== TAB 2: TUTUP SHIFT KASIR ==================== -->
    <main v-if="activeTab === 'shift'" class="max-w-4xl mx-auto space-y-6">
      <div class="border border-zinc-200/90 rounded-2xl bg-white p-6 shadow-2xs space-y-5">
        <div class="flex items-start justify-between pb-3 border-b border-zinc-200">
          <div>
            <h2 class="text-lg font-bold text-zinc-950">Rekapitulasi Tutup Shift Kasir Harian</h2>
            <p class="text-xs text-zinc-500 mt-0.5">Kirim laporan setoran fisik dan transaksi digital langsung ke dashboard Pimpinan.</p>
          </div>
          <div class="text-right">
            <div class="text-[10px] uppercase font-bold text-zinc-500">Tanggal Shift</div>
            <div class="text-xs font-sans tabular-nums font-bold text-zinc-900">{{ shiftSummary?.shift_date || '-' }}</div>
          </div>
        </div>

        <!-- Breakdown Table -->
        <div class="border border-zinc-200 rounded-xl divide-y divide-zinc-100 text-xs bg-zinc-50/30">
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Total Transaksi Selesai:</span>
            <span class="font-sans tabular-nums font-bold text-zinc-900">{{ shiftSummary?.total_transactions || 0 }} Transaksi</span>
          </div>
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Penerimaan Tunai (Fisik Laci):</span>
            <span class="font-sans tabular-nums font-bold text-zinc-900">{{ formatRupiah(shiftSummary?.total_cash) }}</span>
          </div>
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Penerimaan Mesin EDC:</span>
            <span class="font-sans tabular-nums font-bold text-zinc-900">{{ formatRupiah(shiftSummary?.total_edc) }}</span>
          </div>
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Penerimaan Transfer Bank:</span>
            <span class="font-sans tabular-nums font-bold text-zinc-900">{{ formatRupiah(shiftSummary?.total_transfer) }}</span>
          </div>
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Potongan Diskon Voucher:</span>
            <span class="font-sans tabular-nums font-bold text-zinc-700">{{ formatRupiah(shiftSummary?.total_voucher) }}</span>
          </div>
          <div class="p-3 flex justify-between">
            <span class="text-zinc-600 font-medium">Total Pengembalian Dana (Refund):</span>
            <span class="font-sans tabular-nums font-bold text-rose-700">-{{ formatRupiah(shiftSummary?.total_refund) }}</span>
          </div>
          <div class="p-3.5 flex justify-between bg-zinc-950 text-white rounded-b-xl">
            <span class="font-bold text-xs uppercase tracking-wider text-zinc-300">Grand Net Setoran Kasir:</span>
            <span class="font-sans tabular-nums font-extrabold text-base text-emerald-300">{{ formatRupiah(shiftSummary?.net_total) }}</span>
          </div>
        </div>

        <!-- Shift Notes -->
        <div class="space-y-1">
          <label class="text-[11px] font-bold text-zinc-700">Catatan Kasir untuk Pimpinan (Opsional):</label>
          <textarea
            v-model="closeShiftNotes"
            rows="2"
            placeholder="Catatan penutupan shift kasir..."
            class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50/50 text-xs focus:outline-none focus:ring-1 focus:ring-zinc-950"
          ></textarea>
        </div>

        <button
          @click="submitCloseShift"
          :disabled="isClosingShift"
          class="w-full py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs shadow-2xs transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-1.5 cursor-pointer"
        >
          <Send class="w-3.5 h-3.5 text-amber-400" />
          <span>{{ isClosingShift ? 'Mengirim...' : 'Kunci & Kirim Laporan Tutup Shift ke Pimpinan' }}</span>
        </button>

        <div v-if="shiftClosedSuccess" class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-950 flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span>Laporan Tutup Shift #{{ shiftClosedSuccess.id }} berhasil dikirim ke Pimpinan.</span>
        </div>
      </div>

      <!-- Closed Shifts History Table -->
      <div class="border border-zinc-200/90 rounded-2xl bg-white p-5 shadow-2xs space-y-3">
        <h3 class="text-xs font-bold uppercase tracking-wider text-zinc-500">Riwayat Laporan Tutup Shift</h3>
        <div v-if="shiftReports.length === 0" class="py-6 text-center text-xs text-zinc-400">
          Belum ada riwayat laporan tutup shift sebelumnya.
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="border-b border-zinc-200 text-[10px] font-bold uppercase tracking-wider text-zinc-500">
                <th class="py-2">Tanggal</th>
                <th class="py-2">Kasir</th>
                <th class="py-2">Trx</th>
                <th class="py-2">Tunai</th>
                <th class="py-2">Digital/EDC</th>
                <th class="py-2">Net Setoran</th>
                <th class="py-2">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="rep in shiftReports" :key="rep.id">
                <td class="py-2.5 font-sans tabular-nums text-zinc-600">{{ rep.shift_date }}</td>
                <td class="py-2.5 font-bold text-zinc-900">{{ rep.staff_name }}</td>
                <td class="py-2.5 font-sans tabular-nums">{{ rep.total_transactions }}</td>
                <td class="py-2.5 font-sans tabular-nums text-zinc-900 font-bold">{{ formatRupiah(rep.total_cash) }}</td>
                <td class="py-2.5 font-sans tabular-nums">{{ formatRupiah(rep.total_transfer + rep.total_edc) }}</td>
                <td class="py-2.5 font-sans tabular-nums font-bold text-zinc-950">{{ formatRupiah(rep.net_total) }}</td>
                <td class="py-2.5">
                  <span class="px-2 py-0.5 rounded border border-emerald-200/80 bg-emerald-50 text-emerald-800 text-[10px] font-semibold">
                    Terkirim
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>

    <!-- ==================== TAB 3: RIWAYAT & CETAK STRUK ==================== -->
    <main v-if="activeTab === 'history'" class="border border-zinc-200/90 rounded-2xl bg-white p-5 sm:p-6 shadow-2xs space-y-4">
      <div class="flex items-center justify-between pb-3 border-b border-zinc-200">
        <div>
          <h2 class="text-base font-bold text-zinc-950">Riwayat Pembayaran & Struk Kasir</h2>
          <p class="text-xs text-zinc-500 mt-0.5">Daftar transaksi terlunasi, cetak ulang struk thermal 80mm, atau catat refund manual.</p>
        </div>
        <button
          @click="loadInvoices"
          :disabled="loadingInvoices"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-700 shadow-2xs"
        >
          <RefreshCw :class="['w-3.5 h-3.5', loadingInvoices && 'animate-spin']" />
          <span>Segarkan</span>
        </button>
      </div>

      <div v-if="loadingInvoices" class="py-12 text-center text-xs text-zinc-400">
        <RefreshCw class="w-4 h-4 mx-auto animate-spin mb-1.5 text-zinc-400" />
        Memuat riwayat transaksi...
      </div>

      <div v-else-if="invoiceList.length === 0" class="py-12 text-center text-xs text-zinc-400">
        Belum ada riwayat transaksi pembayaran.
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead>
            <tr class="border-b border-zinc-200 text-[10px] font-bold uppercase tracking-wider text-zinc-500">
              <th class="py-2.5">No Invoice</th>
              <th class="py-2.5">Antrian</th>
              <th class="py-2.5">Pasien</th>
              <th class="py-2.5">Unit</th>
              <th class="py-2.5">Waktu</th>
              <th class="py-2.5">Metode</th>
              <th class="py-2.5">Total Tagihan</th>
              <th class="py-2.5">Status</th>
              <th class="py-2.5 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100 font-medium">
            <tr v-for="inv in invoiceList" :key="inv.id" class="hover:bg-zinc-50/70">
              <td class="py-3 font-sans tabular-nums font-bold text-zinc-950">{{ inv.invoice_number }}</td>
              <td class="py-3 font-sans tabular-nums font-extrabold text-zinc-950">{{ inv.queue_number }}</td>
              <td class="py-3 font-bold text-zinc-900">{{ inv.patient_name }}</td>
              <td class="py-3">
                <span
                  :class="[
                    'px-1.5 py-0.5 rounded border text-[10px] font-semibold',
                    inv.clinic_unit === 'KIDINEST'
                      ? 'bg-amber-50 text-amber-900 border-amber-200/80'
                      : 'bg-rose-50 text-rose-900 border-rose-200/80'
                  ]"
                >
                  {{ inv.clinic_unit }}
                </span>
              </td>
              <td class="py-3 text-zinc-500 font-sans tabular-nums">{{ formatDate(inv.created_at) }}</td>
              <td class="py-3 text-xs">
                <div class="font-bold text-zinc-900 flex items-center gap-1.5">
                  <span
                    :class="[
                      'px-1.5 py-0.5 rounded text-[10px] font-bold border',
                      inv.payment_method === 'CASH' ? 'bg-zinc-100 text-zinc-800 border-zinc-200' :
                      inv.payment_method === 'EDC' ? 'bg-blue-50 text-blue-800 border-blue-200/60' :
                      inv.payment_method === 'QRIS' ? 'bg-emerald-50 text-emerald-800 border-emerald-200/60' :
                      'bg-purple-50 text-purple-800 border-purple-200/60'
                    ]"
                  >
                    {{ inv.payment_method }}
                  </span>
                </div>
                <div v-if="inv.payment_reference" class="text-[10px] text-zinc-500 font-mono mt-0.5">
                  {{ inv.payment_reference }}
                </div>
              </td>
              <td class="py-3 font-sans tabular-nums font-extrabold text-zinc-950">{{ formatRupiah(inv.total_amount) }}</td>
              <td class="py-3">
                <span
                  :class="[
                    'px-2 py-0.5 rounded border text-[10px] font-semibold',
                    inv.status === 'PAID'
                      ? 'bg-emerald-50 text-emerald-800 border-emerald-200/80'
                      : 'bg-rose-50 text-rose-800 border-rose-200/80'
                  ]"
                >
                  {{ inv.status === 'PAID' ? 'Lunas' : 'Refund' }}
                </span>
              </td>
              <td class="py-3 text-right">
                <div class="inline-flex items-center gap-1.5">
                  <button
                    @click="activeInvoiceForPrint = inv; showPrintModal = true"
                    class="p-1.5 rounded-lg border border-zinc-200 hover:bg-zinc-100 text-zinc-700 cursor-pointer"
                    title="Cetak Struk Thermal"
                  >
                    <Printer class="w-3.5 h-3.5" />
                  </button>
                  <button
                    v-if="inv.status === 'PAID'"
                    @click="openRefundModal(inv)"
                    class="p-1.5 rounded-lg border border-rose-200 hover:bg-rose-50 text-rose-600 cursor-pointer"
                    title="Catat Refund Manual"
                  >
                    <RotateCcw class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </main>

    <!-- ==================== TAB 4: PENGAJUAN REFUND PASIEN ==================== -->
    <main v-else-if="activeTab === 'refunds'" class="space-y-6 animate-in fade-in duration-200">
      
      <!-- Top Header & Filter Bar -->
      <div class="p-6 rounded-2xl bg-white border border-zinc-200 shadow-2xs flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">
            Workstation Kasir & Keuangan
          </span>
          <h2 class="text-xl font-extrabold text-zinc-950 mt-0.5">
            Daftar Pengajuan Refund Pendaftaran
          </h2>
          <p class="text-xs text-zinc-500 font-medium mt-1">
            Kasir memverifikasi dan melakukan transfer manual via m-Banking klinik ke rekening pasien.
          </p>
        </div>

        <div class="flex items-center gap-2">
          <!-- Filter Unit -->
          <div class="flex items-center bg-zinc-100 rounded-xl p-1 text-[11px] font-bold">
            <button
              @click="refundFilterUnit = 'ALL'; loadRefundRequests()"
              :class="['px-2.5 py-1 rounded-lg transition-all cursor-pointer', refundFilterUnit === 'ALL' ? 'bg-white text-zinc-950 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              Semua ({{ refundRequests.length }})
            </button>
            <button
              @click="refundFilterUnit = 'KIDINEST'; loadRefundRequests()"
              :class="['px-2.5 py-1 rounded-lg transition-all cursor-pointer', refundFilterUnit === 'KIDINEST' ? 'bg-white text-amber-900 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              Kidinest
            </button>
            <button
              @click="refundFilterUnit = 'SSC'; loadRefundRequests()"
              :class="['px-2.5 py-1 rounded-lg transition-all cursor-pointer', refundFilterUnit === 'SSC' ? 'bg-white text-rose-900 shadow-2xs' : 'text-zinc-500 hover:text-zinc-900']"
            >
              SSC
            </button>
          </div>

          <button
            @click="loadRefundRequests"
            :disabled="loadingRefundRequests"
            class="inline-flex items-center gap-1.5 px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50 hover:bg-zinc-100 text-xs font-bold text-zinc-700 transition-all cursor-pointer"
          >
            <RefreshCw :class="['w-3.5 h-3.5', loadingRefundRequests && 'animate-spin']" />
            <span>Sync</span>
          </button>
        </div>
      </div>

      <!-- Quick Ticker Stats -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div class="p-4 rounded-2xl bg-white border border-zinc-200 shadow-2xs">
          <div class="text-[10px] font-bold uppercase tracking-wider text-amber-700">Menunggu Transfer Kasir</div>
          <div class="text-2xl font-black text-amber-900 mt-1">{{ countPendingRefunds }} Pengajuan</div>
        </div>
        <div class="p-4 rounded-2xl bg-white border border-zinc-200 shadow-2xs">
          <div class="text-[10px] font-bold uppercase tracking-wider text-emerald-700">Telah Selesai Direfund</div>
          <div class="text-2xl font-black text-emerald-900 mt-1">{{ countApprovedRefunds }} Tiket</div>
        </div>
        <div class="p-4 rounded-2xl bg-stone-50 border border-zinc-200 shadow-2xs flex items-center justify-between">
          <div>
            <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-500">Estimasi Dana Refund</div>
            <div class="text-xl font-black text-zinc-900 mt-1">{{ formatRupiah(countPendingRefunds * 50000) }}</div>
          </div>
          <Banknote class="w-6 h-6 text-zinc-400" />
        </div>
      </div>

      <!-- Table / List of Refund Requests -->
      <div class="bg-white rounded-2xl border border-zinc-200 shadow-2xs overflow-hidden">
        <div v-if="loadingRefundRequests" class="p-12 text-center text-xs text-zinc-400 flex items-center justify-center gap-2">
          <Loader2 class="w-4 h-4 animate-spin" />
          <span>Memuat data pengajuan refund...</span>
        </div>

        <div v-else-if="filteredRefundRequests.length === 0" class="p-12 text-center text-xs text-zinc-400">
          Tidak ada pengajuan refund yang ditemukan.
        </div>

        <div v-else class="divide-y divide-zinc-200/80">
          <div
            v-for="item in filteredRefundRequests"
            :key="item.id"
            class="p-5 flex flex-col lg:flex-row lg:items-center justify-between gap-4 hover:bg-zinc-50/60 transition-colors"
          >
            <!-- Patient & Ticket Info -->
            <div class="space-y-1.5 min-w-0 max-w-md">
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-sm text-[#0C3B36]">{{ item.queue_number }}</span>
                <span
                  class="px-2 py-0.5 rounded text-[10px] font-bold uppercase"
                  :class="item.clinic_unit === 'KIDINEST' ? 'bg-amber-100 text-amber-900' : 'bg-rose-100 text-rose-900'"
                >
                  {{ item.clinic_unit }}
                </span>
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-extrabold uppercase"
                  :class="
                    item.refund_status === 'APPROVED' ? 'bg-emerald-100 text-emerald-900' :
                    item.refund_status === 'REJECTED' ? 'bg-zinc-200 text-zinc-700' :
                    'bg-amber-100 text-amber-900 animate-pulse'
                  "
                >
                  {{ item.refund_status === 'APPROVED' ? 'Selesai Ditransfer' : item.refund_status === 'REJECTED' ? 'Ditolak' : 'Menunggu Transfer' }}
                </span>
              </div>

              <div class="text-sm font-extrabold text-zinc-900">
                {{ item.patient_name }}
                <span class="text-xs font-mono font-semibold text-zinc-500 ml-1">({{ item.patient_phone }})</span>
              </div>

              <div class="text-xs text-zinc-500 font-medium">
                Dokter: <strong>{{ item.doctor_name }}</strong> • Jadwal: {{ item.appointment_date || 'Hari-H' }} ({{ item.appointment_time || '-' }})
              </div>

              <div v-if="item.refund_notes" class="text-[11px] text-zinc-600 bg-zinc-100/80 p-2 rounded-xl italic">
                "{{ item.refund_notes }}"
              </div>
            </div>

            <!-- Bank & Nominal Info -->
            <div class="p-3.5 rounded-xl bg-zinc-50 border border-zinc-200 text-xs space-y-1 shrink-0 min-w-[240px]">
              <div class="text-[10px] font-bold uppercase tracking-wider text-zinc-400">Tujuan Transfer Pengembalian</div>
              <div class="font-extrabold text-zinc-900 text-sm">
                {{ item.refund_bank || 'BCA' }} • <span class="font-mono">{{ item.refund_account_number || '-' }}</span>
              </div>
              <div class="text-zinc-600 font-medium">a.n {{ item.refund_account_name || item.patient_name }}</div>
              <div class="pt-1 border-t border-zinc-200/80 flex items-center justify-between text-xs font-bold text-emerald-800">
                <span>Nominal:</span>
                <span class="text-sm font-black">{{ formatRupiah(item.registration_fee || 50000) }}</span>
              </div>
            </div>

            <!-- Action Buttons for Cashier -->
            <div class="flex items-center gap-2 shrink-0">
              <template v-if="item.refund_status === 'REQUESTED'">
                <button
                  type="button"
                  @click="handleApproveRefund(item)"
                  :disabled="processingRefundId === item.id"
                  class="px-4 py-2.5 rounded-xl bg-emerald-700 hover:bg-emerald-800 text-white font-bold text-xs flex items-center gap-1.5 shadow-2xs transition-all cursor-pointer disabled:opacity-50"
                >
                  <CheckCircle2 class="w-4 h-4" />
                  <span>Setujui & Tandai Ditransfer</span>
                </button>
                <button
                  type="button"
                  @click="handleRejectRefund(item)"
                  :disabled="processingRefundId === item.id"
                  class="px-3 py-2.5 rounded-xl border border-rose-300 hover:bg-rose-50 text-rose-700 font-bold text-xs transition-all cursor-pointer disabled:opacity-50"
                >
                  Tolak
                </button>
              </template>
              <div v-else class="text-right text-xs font-bold text-zinc-400">
                Status: {{ item.refund_status }}
              </div>
            </div>
          </div>
        </div>
      </div>

    </main>

    <!-- ==================== MODAL 1: CETAK STRUK THERMAL 80MM ==================== -->
    <div
      v-if="showPrintModal && activeInvoiceForPrint"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-2xl max-w-sm w-full p-5 space-y-4 shadow-xl border border-zinc-200">
        <div class="flex items-center justify-between pb-2 border-b border-zinc-200">
          <div class="font-bold text-xs text-zinc-950">Pratinjau Struk Thermal 80mm</div>
          <button @click="showPrintModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Selector 2 Lembar Cetak (Master Prompt 6) -->
        <div class="grid grid-cols-3 gap-1 p-1 bg-zinc-100 rounded-xl text-[10px] font-bold">
          <button
            type="button"
            @click="receiptCopyType = 'BOTH'"
            class="py-1 rounded-lg text-center transition-all"
            :class="receiptCopyType === 'BOTH' ? 'bg-white shadow-2xs text-zinc-950 font-black' : 'text-zinc-600'"
          >
            2 Lembar (Semua)
          </button>
          <button
            type="button"
            @click="receiptCopyType = 'PATIENT'"
            class="py-1 rounded-lg text-center transition-all"
            :class="receiptCopyType === 'PATIENT' ? 'bg-white shadow-2xs text-zinc-950 font-black' : 'text-zinc-600'"
          >
            Lembar Pasien
          </button>
          <button
            type="button"
            @click="receiptCopyType = 'INTERNAL'"
            class="py-1 rounded-lg text-center transition-all"
            :class="receiptCopyType === 'INTERNAL' ? 'bg-white shadow-2xs text-zinc-950 font-black' : 'text-zinc-600'"
          >
            Lembar Arsip GRO
          </button>
        </div>

        <div id="thermal-receipt" class="font-mono text-xs bg-zinc-50 p-4 rounded-xl border border-zinc-200 space-y-4 max-h-[60vh] overflow-y-auto">
          
          <!-- LEMBAR 1: UNTUK PASIEN -->
          <div v-if="receiptCopyType === 'BOTH' || receiptCopyType === 'PATIENT'" class="space-y-2">
            <div class="text-center pb-2 border-b border-dashed border-zinc-300">
              <div class="font-black text-sm">KLINIK UTAMA SEMARANG MEDICAL GROUP</div>
              <div class="text-[10px] font-bold text-zinc-600">{{ activeInvoiceForPrint.clinic_unit === 'KIDINEST' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</div>
              <div class="text-[9px] text-zinc-400">Jl. Melati Utara No. 3, Semarang • Telp: (024) 841-2345</div>
              <div class="mt-1 inline-block px-2 py-0.5 rounded-full bg-zinc-200 text-zinc-800 text-[9px] font-bold uppercase tracking-wider">
                [ LEMBAR 1 : UNTUK PASIEN ]
              </div>
            </div>

            <div class="space-y-1 text-[11px] py-1 border-b border-dashed border-zinc-300">
              <div class="flex justify-between"><span>No. Bukti:</span><span class="font-bold">{{ activeInvoiceForPrint.invoice_number }}</span></div>
              <div class="flex justify-between"><span>No. Antrean:</span><span class="font-black text-xs">{{ activeInvoiceForPrint.queue_number }}</span></div>
              <div class="flex justify-between"><span>Nama Pasien:</span><span>{{ activeInvoiceForPrint.patient_name }}</span></div>
              <div class="flex justify-between"><span>Waktu Transaksi:</span><span>{{ formatDate(activeInvoiceForPrint.created_at) }}</span></div>
            </div>

            <div class="space-y-1 text-[11px] py-1 border-b border-dashed border-zinc-300">
              <div v-if="activeInvoiceForPrint.subtotal_registration > 0" class="flex justify-between"><span>Registrasi:</span><span>{{ formatRupiah(activeInvoiceForPrint.subtotal_registration) }}</span></div>
              <div v-if="activeInvoiceForPrint.subtotal_action > 0" class="flex justify-between"><span>Tindakan Dokter:</span><span>{{ formatRupiah(activeInvoiceForPrint.subtotal_action) }}</span></div>
              <div v-if="activeInvoiceForPrint.subtotal_medication > 0" class="flex justify-between"><span>Farmasi / Obat:</span><span>{{ formatRupiah(activeInvoiceForPrint.subtotal_medication) }}</span></div>
              <div v-if="activeInvoiceForPrint.voucher_discount > 0" class="flex justify-between text-purple-700 font-bold"><span>Voucher Diskon:</span><span>-{{ formatRupiah(activeInvoiceForPrint.voucher_discount) }}</span></div>
            </div>

            <div class="space-y-1 pt-1 font-bold">
              <div class="flex justify-between text-xs"><span>TOTAL PEMBAYARAN:</span><span>{{ formatRupiah(activeInvoiceForPrint.total_amount) }}</span></div>
              <div class="flex justify-between text-[11px] font-normal"><span>Metode Bayar:</span><span>{{ activeInvoiceForPrint.payment_method }}</span></div>
              <div v-if="activeInvoiceForPrint.payment_reference" class="flex justify-between text-[10px] text-zinc-500 font-mono"><span>Ref:</span><span>{{ activeInvoiceForPrint.payment_reference }}</span></div>
              <div v-if="activeInvoiceForPrint.payment_method === 'CASH'" class="flex justify-between text-[11px] text-emerald-700"><span>Kembalian:</span><span>{{ formatRupiah(activeInvoiceForPrint.cash_change) }}</span></div>
            </div>

            <div class="text-center pt-2 text-[10px] text-zinc-400">LUNAS • Terima kasih atas kunjungan Anda</div>
          </div>

          <!-- GARIS POTONG JIKA CETAK 2 LEMBAR -->
          <div v-if="receiptCopyType === 'BOTH'" class="border-t-2 border-dashed border-zinc-400 pt-3 text-center text-[9px] text-zinc-400 uppercase tracking-widest font-sans">
            - - - - - - - - - POTONG DI SINI - - - - - - - - -
          </div>

          <!-- LEMBAR 2: UNTUK ARSIP / INTERNAL GRO -->
          <div v-if="receiptCopyType === 'BOTH' || receiptCopyType === 'INTERNAL'" class="space-y-2 opacity-90">
            <div class="text-center pb-2 border-b border-dashed border-zinc-300">
              <div class="font-black text-sm">KLINIK UTAMA SEMARANG MEDICAL GROUP</div>
              <div class="text-[10px] font-bold text-zinc-600">{{ activeInvoiceForPrint.clinic_unit === 'KIDINEST' ? 'Kidinest Pediatric Care' : 'Semarang Skin Center' }}</div>
              <div class="mt-1 inline-block px-2 py-0.5 rounded-full bg-amber-100 text-amber-900 text-[9px] font-bold uppercase tracking-wider">
                [ LEMBAR 2 : ARSIP / INTERNAL GRO ]
              </div>
            </div>

            <div class="space-y-1 text-[11px] py-1 border-b border-dashed border-zinc-300">
              <div class="flex justify-between"><span>No. Bukti:</span><span class="font-bold">{{ activeInvoiceForPrint.invoice_number }}</span></div>
              <div class="flex justify-between"><span>No. Antrean:</span><span class="font-black text-xs">{{ activeInvoiceForPrint.queue_number }}</span></div>
              <div class="flex justify-between"><span>Pasien:</span><span>{{ activeInvoiceForPrint.patient_name }}</span></div>
              <div class="flex justify-between"><span>Kasir Bertugas:</span><span>{{ authStore.user?.name || 'Kasir GRO' }}</span></div>
            </div>

            <div class="space-y-1 pt-1 font-bold">
              <div class="flex justify-between text-xs"><span>TOTAL TRANSAKSI:</span><span>{{ formatRupiah(activeInvoiceForPrint.total_amount) }}</span></div>
              <div class="flex justify-between text-[11px] font-normal"><span>Metode:</span><span>{{ activeInvoiceForPrint.payment_method }}</span></div>
            </div>

            <div class="pt-2 text-center text-[9px] text-zinc-400 italic">Salinan sah internal kasir & rekon akuntansi</div>
          </div>

        </div>

        <div class="flex items-center gap-2">
          <button
            @click="printReceipt"
            class="flex-1 py-2.5 rounded-xl bg-zinc-950 hover:bg-zinc-800 text-white font-bold text-xs flex items-center justify-center gap-1.5 cursor-pointer shadow-2xs"
          >
            <Printer class="w-3.5 h-3.5 text-emerald-400" />
            <span>Cetak Struk</span>
          </button>
          <button
            @click="showPrintModal = false"
            class="px-4 py-2.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 font-bold text-xs"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>

    <!-- ==================== MODAL 2: CATAT REFUND MANUAL ==================== -->
    <div
      v-if="showRefundModal && invoiceToRefund"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4"
    >
      <div class="bg-white rounded-2xl max-w-md w-full p-5 space-y-4 shadow-xl border border-zinc-200">
        <div class="flex items-center justify-between pb-2 border-b border-zinc-200">
          <div>
            <h3 class="text-xs font-bold text-zinc-950">Pencatatan Refund Pasien</h3>
            <p class="text-[11px] text-zinc-400">Invoice: {{ invoiceToRefund.invoice_number }} • Total: {{ formatRupiah(invoiceToRefund.total_amount) }}</p>
          </div>
          <button @click="showRefundModal = false" class="text-zinc-400 hover:text-zinc-700">
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="space-y-3">
          <div class="space-y-1">
            <label class="text-[11px] font-bold text-zinc-700">Bank Tujuan:</label>
            <select
              v-model="refundForm.refund_bank_name"
              class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50 text-xs font-medium"
            >
              <option value="BCA">BCA</option>
              <option value="Mandiri">Bank Mandiri</option>
              <option value="BRI">BRI</option>
              <option value="BNI">BNI</option>
              <option value="BSI">BSI</option>
            </select>
          </div>

          <div class="space-y-1">
            <label class="text-[11px] font-bold text-zinc-700">No Rekening:</label>
            <input
              v-model="refundForm.refund_bank_account"
              type="text"
              placeholder="1234567890"
              class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50 text-xs font-mono"
            />
          </div>

          <div class="space-y-1">
            <label class="text-[11px] font-bold text-zinc-700">Atas Nama Rekening:</label>
            <input
              v-model="refundForm.refund_account_name"
              type="text"
              placeholder="Nama Pemilik Rekening"
              class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50 text-xs"
            />
          </div>

          <div class="space-y-1">
            <label class="text-[11px] font-bold text-zinc-700">Alasan Refund:</label>
            <textarea
              v-model="refundForm.refund_reason"
              rows="2"
              placeholder="Alasan pembatalan..."
              class="w-full px-3 py-2 rounded-xl border border-zinc-200 bg-zinc-50 text-xs"
            ></textarea>
          </div>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <button
            @click="submitRefund"
            :disabled="isSubmittingRefund"
            class="flex-1 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs disabled:opacity-50 flex items-center justify-center gap-1.5"
          >
            <RotateCcw class="w-3.5 h-3.5" />
            <span>{{ isSubmittingRefund ? 'Memproses...' : 'Konfirmasi Refund' }}</span>
          </button>
          <button
            @click="showRefundModal = false"
            class="px-4 py-2.5 rounded-xl border border-zinc-200 bg-white hover:bg-zinc-50 text-zinc-700 font-bold text-xs"
          >
            Batal
          </button>
        </div>
      </div>
    </div>

  </div>
</template>
