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
  patient_name: string
  patient_nik: string
  registration_fee: number
  consultation_fee: number
  medication_fee: number
  total_amount: number
  payment_method: string
  payment_status: string
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
  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(options?.headers || {}),
    },
  })
  const json = await res.json()
  if (!json.success) {
    throw new Error(json.error || 'Terjadi kesalahan sistem')
  }
  return json.data as T
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
  getInvoiceByQueue: (queueId: number) => request<Invoice>(`/invoices/queue/${queueId}`),

  // SATUSEHAT
  getSatusehatStats: () => request<SatusehatStats>('/satusehat/stats'),
  getSatusehatEncounter: (queueId: number) => request<any>(`/satusehat/encounter/${queueId}`),
  syncSatusehat: (queueId: number) => request<void>(`/satusehat/sync/${queueId}`, { method: 'POST' }),
}
