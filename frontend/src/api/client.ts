const API_BASE = '/api'

export interface QueueItem {
  id: number
  patient_id: number
  patient_name: string
  patient_nik: string
  patient_phone: string
  patient_ihs?: string
  queue_number: string
  complaint: string
  payment_method: 'QRIS' | 'CASH'
  payment_status: 'PENDING' | 'PAID'
  queue_status: 'WAITING' | 'CALLING' | 'IN_EXAM' | 'COMPLETED' | 'CANCELLED'
  paid_at?: string
  created_at: string
}

export interface DrugBatch {
  id: number
  drug_id: number
  batch_number: string
  expired_date: string
  physical_stock: number
  reserved_stock: number
  available_stock: number
  buy_price: number
  is_expiring_soon: boolean
  created_at: string
}

export interface Drug {
  id: number
  name: string
  kfa_code: string
  unit: string
  category: string
  min_stock_alert: number
  sell_price: number
  physical_stock: number
  reserved_stock: number
  available_stock: number
  is_low_stock: boolean
  has_expiring: boolean
  batches?: DrugBatch[]
  created_at: string
}

export interface PrescriptionItem {
  id: number
  prescription_id: number
  drug_id: number
  drug_name: string
  batch_id?: number
  batch_number: string
  quantity: number
  dosage_instruction: string
  unit_price: number
  subtotal: number
}

export interface Prescription {
  id: number
  consultation_id: number
  queue_id: number
  patient_name: string
  patient_nik: string
  doctor_name: string
  status: 'RESERVED' | 'PREPARING' | 'DISPENSED' | 'CANCELLED'
  items: PrescriptionItem[]
  total_amount: number
  created_at: string
}

export interface Consultation {
  id: number
  queue_id: number
  doctor_name: string
  anamnesis: string
  diagnosis: string
  icd10_code: string
  status: 'DRAFT' | 'COMPLETED'
  prescription?: Prescription
  created_at: string
}

export interface Invoice {
  id: number
  queue_id: number
  invoice_number: string
  invoice_type: 'REGISTRATION' | 'PHARMACY'
  patient_name: string
  patient_nik: string
  registration_fee: number
  consultation_fee: number
  medication_fee: number
  total_amount: number
  payment_method: string
  payment_status: string
  items?: PrescriptionItem[]
  doctor_share: number
  clinic_share: number
  tax_share: number
  created_at: string
}

export interface SatusehatStats {
  is_ready: boolean
  org_id: string
  total_patients: number
  verified_patients: number
  mapped_drugs: number
  total_drugs: number
  synced_logs: number
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const token = localStorage.getItem('clinic_token')
  const authHeaders: Record<string, string> = {}
  if (token) {
    authHeaders['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...authHeaders,
      ...(options?.headers || {}),
    },
  })
  const json = await res.json()
  if (!json.success) {
    throw new Error(json.error || 'Terjadi kesalahan sistem')
  }
  return json.data as T
}

export interface PendingBill {
  queue_id: number
  queue_number: string
  patient_id: number
  patient_name: string
  patient_phone: string
  clinic_unit: 'KIDINEST' | 'SSC'
  category: string
  doctor_name: string
  subtotal_registration: number
  subtotal_action: number
  subtotal_medication: number
  total_due: number
  bill_type: string
  queue_status: string
  payment_status: string
  created_at: string
}

export interface CashierInvoice {
  id: number
  queue_id: number
  patient_id: number
  invoice_number: string
  patient_name: string
  patient_phone: string
  clinic_unit: string
  queue_number: string
  subtotal_registration: number
  subtotal_action: number
  subtotal_medication: number
  voucher_discount: number
  total_amount: number
  payment_method: 'CASH' | 'BANK_TRANSFER' | 'EDC' | 'QRIS'
  cash_received: number
  cash_change: number
  payment_reference: string
  cashier_staff_id?: number
  cashier_name: string
  status: 'PAID' | 'REFUNDED'
  refund_bank_name?: string
  refund_bank_account?: string
  refund_account_name?: string
  refund_reason?: string
  refunded_at?: string
  notes?: string
  created_at: string
}

export interface CashierShiftSummary {
  shift_date: string
  total_transactions: number
  total_cash: number
  total_transfer: number
  total_edc: number
  total_voucher: number
  total_refund: number
  net_total: number
}

export interface CashierShiftReport {
  id: number
  staff_id: number
  staff_name: string
  shift_date: string
  total_transactions: number
  total_cash: number
  total_transfer: number
  total_edc: number
  total_voucher: number
  total_refund: number
  net_total: number
  notes?: string
  sent_to_owner_at: string
  created_at: string
}

export interface DoctorQueuePatient {
  queue_id: number
  queue_number: string
  patient_id: number
  patient_name: string
  patient_phone?: string
  clinic_unit: 'KIDINEST' | 'SSC'
  category: string
  doctor_id?: number
  doctor_name: string
  registration_fee: number
  payment_status: string
  queue_status: 'WAITING' | 'CALLING' | 'IN_EXAM' | 'PHARMACY_WAITING' | 'COMPLETED'
  called_at?: string
  created_at: string
}

export interface DoctorCatalogDrug {
  id: number
  clinic_unit: string
  name: string
  kfa_code?: string
  unit: string
  category: string
  sell_price: number
  physical_stock: number
  reserved_stock: number
  available_stock: number
  is_active: boolean
}

export interface DoctorPrescriptionItem {
  id?: number
  drug_id?: number
  item_name: string
  is_manual: boolean
  quantity: number
  unit_price: number
  subtotal: number
  dosage_instruction: string
}

export interface DoctorConsultation {
  id: number
  queue_id: number
  doctor_id?: number
  doctor_name: string
  anamnesis: string
  diagnosis: string
  icd10_code: string
  action_notes?: string
  action_fee: number
  status: string
  prescription?: {
    id: number
    status: string
    notes?: string
    items: DoctorPrescriptionItem[]
  }
}

export interface PharmacyPrescriptionItem {
  id: number
  prescription_id: number
  drug_id?: number
  item_name: string
  is_manual: boolean
  quantity: number
  unit_price: number
  subtotal: number
  dosage_instruction: string
}

export interface PharmacyPrescription {
  id: number
  consultation_id: number
  queue_id: number
  queue_number: string
  patient_name: string
  clinic_unit: string
  doctor_name: string
  status: 'RESERVED' | 'PREPARING' | 'DISPENSED' | 'CANCELLED'
  notes?: string
  items: PharmacyPrescriptionItem[]
  created_at: string
  updated_at: string
}

export interface PharmacyDrugBatch {
  id: number
  drug_id: number
  batch_number: string
  expired_date: string
  physical_stock: number
  reserved_stock: number
  available_stock: number
  is_expiring_soon: boolean
  created_at: string
}

export interface PharmacyDrug {
  id: number
  clinic_unit: string
  name: string
  kfa_code?: string
  unit: string
  category: string
  min_stock_alert: number
  buy_price: number
  sell_price: number
  physical_stock: number
  reserved_stock: number
  available_stock: number
  is_low_stock: boolean
  has_expiring_soon: boolean
  batches?: PharmacyDrugBatch[]
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface PharmacyOpnameLog {
  id: number
  drug_id: number
  drug_name: string
  batch_id?: number
  batch_number?: string
  system_stock: number
  actual_stock: number
  discrepancy: number
  adjusted_by: string
  notes?: string
  created_at: string
}

export const api = {
  // Queue
  getQueues: () => request<QueueItem[]>('/queues'),
  createQueue: (data: {
    nik: string
    name: string
    phone: string
    birth_date?: string
    gender?: string
    complaint: string
    payment_method: 'QRIS' | 'CASH'
  }) => request<QueueItem>('/queues', { method: 'POST', body: JSON.stringify(data) }),
  payQueue: (id: number) => request<void>(`/queues/${id}/pay`, { method: 'POST' }),
  updateQueueStatus: (id: number, status: string) =>
    request<void>(`/queues/${id}/status`, { method: 'POST', body: JSON.stringify({ status }) }),

  // Pharmacy
  getDrugs: () => request<Drug[]>('/drugs'),
  getDrugDetail: (id: number) => request<Drug>(`/drugs/${id}`),
  createDrug: (data: {
    name: string
    kfa_code: string
    unit: string
    category: string
    min_stock_alert: number
    sell_price: number
  }) => request<Drug>('/drugs', { method: 'POST', body: JSON.stringify(data) }),
  restockDrug: (id: number, data: { batch_number: string; expired_date: string; quantity: number; buy_price: number }) =>
    request<void>(`/drugs/${id}/restock`, { method: 'POST', body: JSON.stringify(data) }),

  // Consultation & Prescriptions
  getConsultationByQueue: (queueId: number) => request<Consultation>(`/consultations/queue/${queueId}`),
  saveConsultation: (data: {
    queue_id: number
    doctor_name: string
    anamnesis: string
    diagnosis: string
    icd10_code: string
    prescription: Array<{ drug_id: number; quantity: number; dosage_instruction: string }>
  }) => request<Consultation>('/consultations', { method: 'POST', body: JSON.stringify(data) }),
  getPrescriptions: () => request<Prescription[]>('/prescriptions'),
  dispensePrescription: (id: number) => request<void>(`/prescriptions/${id}/dispense`, { method: 'POST' }),
  cancelPrescription: (id: number) => request<void>(`/prescriptions/${id}/cancel`, { method: 'POST' }),

  // Billing
  getInvoices: () => request<Invoice[]>('/invoices'),
  getInvoiceByQueue: (queueId: number) => request<Invoice[]>(`/invoices/queue/${queueId}`),

  // Cashier & GRO Bounded Context
  getPendingBills: () => request<PendingBill[]>('/cashier/pending-bills'),
  payBill: (data: {
    queue_id: number
    subtotal_registration: number
    subtotal_action: number
    subtotal_medication: number
    voucher_discount: number
    payment_method: 'CASH' | 'BANK_TRANSFER' | 'EDC'
    cash_received: number
    payment_reference?: string
    notes?: string
  }) => request<CashierInvoice>('/cashier/invoices', { method: 'POST', body: JSON.stringify(data) }),
  registerWalkIn: (data: {
    patient_name: string
    patient_phone: string
    clinic_unit: 'KIDINEST' | 'SSC'
    doctor_id?: number
    doctor_name: string
    payment_method: 'CASH' | 'QRIS' | 'BANK_TRANSFER'
    is_paid_immediate: boolean
  }) => request<{ bill: PendingBill; invoice?: CashierInvoice }>('/cashier/walkin', { method: 'POST', body: JSON.stringify(data) }),
  getCashierInvoices: (limit = 50, offset = 0) =>
    request<CashierInvoice[]>(`/cashier/invoices?limit=${limit}&offset=${offset}`),
  getCashierInvoice: (id: number) => request<CashierInvoice>(`/cashier/invoices/${id}`),
  refundCashierInvoice: (
    id: number,
    data: {
      refund_bank_name: string
      refund_bank_account: string
      refund_account_name: string
      refund_reason: string
    }
  ) => request<CashierInvoice>(`/cashier/invoices/${id}/refund`, { method: 'POST', body: JSON.stringify(data) }),
  getCashierShiftSummary: (date = '') =>
    request<CashierShiftSummary>(`/cashier/shift-summary?date=${date}`),
  closeCashierShift: (notes = '') =>
    request<CashierShiftReport>('/cashier/close-shift', { method: 'POST', body: JSON.stringify({ notes }) }),
  getCashierShiftReports: () => request<CashierShiftReport[]>('/cashier/shift-reports'),

  // Doctor Consultation & E-Prescription Bounded Context
  getDoctorQueues: (doctorName = '', unit = '') =>
    request<DoctorQueuePatient[]>(
      `/consultations/queues?doctor_name=${encodeURIComponent(doctorName)}&unit=${encodeURIComponent(unit)}`
    ),
  callDoctorPatient: (queueId: number) =>
    request<any>(`/consultations/queues/${queueId}/call`, { method: 'POST' }),
  startDoctorExam: (queueId: number) =>
    request<any>(`/consultations/queues/${queueId}/start`, { method: 'POST' }),
  getDoctorDrugs: (unit = '') =>
    request<DoctorCatalogDrug[]>(`/consultations/drugs?unit=${encodeURIComponent(unit)}`),
  getDoctorConsultationByQueue: (queueId: number) =>
    request<DoctorConsultation | null>(`/consultations/queue/${queueId}`),
  submitDoctorConsultation: (data: {
    queue_id: number
    doctor_id?: number
    doctor_name: string
    anamnesis: string
    diagnosis: string
    icd10_code: string
    action_notes?: string
    action_fee: number
    prescription_items: Array<{
      drug_id?: number
      item_name: string
      is_manual: boolean
      quantity: number
      unit_price: number
      dosage_instruction: string
    }>
  }) => request<DoctorConsultation>('/consultations', { method: 'POST', body: JSON.stringify(data) }),

  // Pharmacy & Stock Opname Bounded Context
  getPharmacyPrescriptions: (unit = '') =>
    request<PharmacyPrescription[]>(`/pharmacy/prescriptions?unit=${encodeURIComponent(unit)}`),
  preparePharmacyPrescription: (id: number) =>
    request<any>(`/pharmacy/prescriptions/${id}/status`, { method: 'POST' }),
  dispensePharmacyPrescription: (id: number) =>
    request<any>(`/pharmacy/prescriptions/${id}/dispense`, { method: 'POST' }),
  cancelPharmacyPrescription: (id: number) =>
    request<any>(`/pharmacy/prescriptions/${id}/cancel`, { method: 'POST' }),
  getPharmacyDrugs: (unit = '') =>
    request<PharmacyDrug[]>(`/pharmacy/drugs?unit=${encodeURIComponent(unit)}`),
  createPharmacyDrug: (data: {
    clinic_unit: string
    name: string
    kfa_code?: string
    unit: string
    category: string
    min_stock_alert: number
    buy_price: number
    sell_price: number
  }) => request<PharmacyDrug>('/pharmacy/drugs', { method: 'POST', body: JSON.stringify(data) }),
  restockPharmacyDrug: (
    drugId: number,
    data: { batch_number: string; expired_date: string; quantity: number; buy_price: number }
  ) => request<any>(`/pharmacy/drugs/${drugId}/restock`, { method: 'POST', body: JSON.stringify(data) }),
  getPharmacyOpnameLogs: () => request<PharmacyOpnameLog[]>('/pharmacy/opname/logs'),
  importPharmacyOpnameCSV: (csvText: string) =>
    request<any>('/pharmacy/opname/import', {
      method: 'POST',
      body: csvText,
      headers: { 'Content-Type': 'text/plain' },
    }),

  // SATUSEHAT
  getSatusehatStats: () => request<SatusehatStats>('/satusehat/stats'),
  getSatusehatEncounter: (queueId: number) => request<any>(`/satusehat/encounter/${queueId}`),
  syncSatusehat: (queueId: number) => request<void>(`/satusehat/sync/${queueId}`, { method: 'POST' }),
}
