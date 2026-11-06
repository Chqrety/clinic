package domain

import "time"

type Patient struct {
	ID        int64     `json:"id"`
	NIK       string    `json:"nik"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BirthDate string    `json:"birth_date"`
	Gender    string    `json:"gender"`
	IHSNumber string    `json:"ihs_number"`
	CreatedAt time.Time `json:"created_at"`
}

type QueueItem struct {
	ID            int64      `json:"id"`
	PatientID     int64      `json:"patient_id"`
	PatientName   string     `json:"patient_name"`
	PatientNIK    string     `json:"patient_nik"`
	PatientPhone  string     `json:"patient_phone"`
	PatientIHS    string     `json:"patient_ihs"`
	QueueNumber   string     `json:"queue_number"`
	Complaint     string     `json:"complaint"`
	PaymentMethod string     `json:"payment_method"` // "QRIS" / "CASH"
	PaymentStatus string     `json:"payment_status"` // "PENDING" / "PAID"
	QueueStatus   string     `json:"queue_status"`   // "WAITING", "CALLING", "IN_EXAM", "COMPLETED", "CANCELLED"
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type DrugBatch struct {
	ID             int64     `json:"id"`
	DrugID         int64     `json:"drug_id"`
	BatchNumber    string    `json:"batch_number"`
	ExpiredDate    string    `json:"expired_date"`
	PhysicalStock  int       `json:"physical_stock"`
	ReservedStock  int       `json:"reserved_stock"`
	AvailableStock int       `json:"available_stock"`
	BuyPrice       int       `json:"buy_price"`
	IsExpiringSoon bool      `json:"is_expiring_soon"` // < 3 bulan
	CreatedAt      time.Time `json:"created_at"`
}

type Drug struct {
	ID             int64       `json:"id"`
	Name           string      `json:"name"`
	KFACode        string      `json:"kfa_code"`
	Unit           string      `json:"unit"`
	Category       string      `json:"category"`
	MinStockAlert  int         `json:"min_stock_alert"`
	SellPrice      int         `json:"sell_price"`
	PhysicalStock  int         `json:"physical_stock"`
	ReservedStock  int         `json:"reserved_stock"`
	AvailableStock int         `json:"available_stock"`
	IsLowStock     bool        `json:"is_low_stock"`
	HasExpiring    bool        `json:"has_expiring"`
	Batches        []DrugBatch `json:"batches,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}

type Consultation struct {
	ID           int64         `json:"id"`
	QueueID      int64         `json:"queue_id"`
	DoctorName   string        `json:"doctor_name"`
	Anamnesis    string        `json:"anamnesis"`
	Diagnosis    string        `json:"diagnosis"`
	ICD10Code    string        `json:"icd10_code"`
	Status       string        `json:"status"` // "DRAFT" / "COMPLETED"
	Prescription *Prescription `json:"prescription,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
}

type PrescriptionItem struct {
	ID                int64  `json:"id"`
	PrescriptionID    int64  `json:"prescription_id"`
	DrugID            int64  `json:"drug_id"`
	DrugName          string `json:"drug_name"`
	BatchID           *int64 `json:"batch_id"`
	BatchNumber       string `json:"batch_number"`
	Quantity          int    `json:"quantity"`
	DosageInstruction string `json:"dosage_instruction"`
	UnitPrice         int    `json:"unit_price"`
	Subtotal          int    `json:"subtotal"`
}

type Prescription struct {
	ID             int64              `json:"id"`
	ConsultationID int64              `json:"consultation_id"`
	QueueID        int64              `json:"queue_id"`
	PatientName    string             `json:"patient_name"`
	PatientNIK     string             `json:"patient_nik"`
	DoctorName     string             `json:"doctor_name"`
	Status         string             `json:"status"` // "RESERVED", "PREPARING", "DISPENSED", "CANCELLED"
	Items          []PrescriptionItem `json:"items"`
	TotalAmount    int                `json:"total_amount"`
	CreatedAt      time.Time          `json:"created_at"`
}

type Invoice struct {
	ID              int64     `json:"id"`
	QueueID         int64     `json:"queue_id"`
	InvoiceNumber   string    `json:"invoice_number"`
	PatientName     string    `json:"patient_name"`
	PatientNIK      string    `json:"patient_nik"`
	RegistrationFee int       `json:"registration_fee"`
	ConsultationFee int       `json:"consultation_fee"`
	MedicationFee   int       `json:"medication_fee"`
	TotalAmount     int       `json:"total_amount"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentStatus   string    `json:"payment_status"`
	// Internal Breakdown (Audit / Klinik)
	DoctorShare int       `json:"doctor_share"` // 70% jasa konsultasi
	ClinicShare int       `json:"clinic_share"` // Registrasi + 30% konsultasi + margin obat
	TaxShare    int       `json:"tax_share"`    // PPN / Pajak
	CreatedAt   time.Time `json:"created_at"`
}

type SatusehatLog struct {
	ID           int64      `json:"id"`
	ResourceType string     `json:"resource_type"`
	EntityID     int64      `json:"entity_id"`
	Status       string     `json:"status"` // "READY" / "SYNCED"
	FHIRPayload  string     `json:"fhir_payload"`
	SyncedAt     *time.Time `json:"synced_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// JSON API Response envelope
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}
