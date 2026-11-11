package domain

import (
	"context"
	"time"
)

type DoctorQueuePatient struct {
	QueueID         int64      `json:"queue_id"`
	QueueNumber     string     `json:"queue_number"`
	PatientID       int64      `json:"patient_id"`
	PatientName     string     `json:"patient_name"`
	PatientPhone    string     `json:"patient_phone,omitempty"`
	ClinicUnit      string     `json:"clinic_unit"`
	Category        string     `json:"category"`
	DoctorID        *int64     `json:"doctor_id,omitempty"`
	DoctorName      string     `json:"doctor_name"`
	RegistrationFee int        `json:"registration_fee"`
	PaymentStatus   string     `json:"payment_status"`
	QueueStatus     string     `json:"queue_status"`
	CalledAt        *time.Time `json:"called_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CatalogDrug struct {
	ID             int64  `json:"id"`
	ClinicUnit     string `json:"clinic_unit"`
	Name           string `json:"name"`
	KFACode        string `json:"kfa_code,omitempty"`
	Unit           string `json:"unit"`
	Category       string `json:"category"`
	SellPrice      int    `json:"sell_price"`
	PhysicalStock  int    `json:"physical_stock"`
	ReservedStock  int    `json:"reserved_stock"`
	AvailableStock int    `json:"available_stock"`
	IsActive       bool   `json:"is_active"`
}

type PrescriptionItem struct {
	ID                int64     `json:"id"`
	PrescriptionID    int64     `json:"prescription_id"`
	DrugID            *int64    `json:"drug_id,omitempty"`
	ItemName          string    `json:"item_name"`
	IsManual          bool      `json:"is_manual"`
	Quantity          int       `json:"quantity"`
	UnitPrice         int       `json:"unit_price"`
	Subtotal          int       `json:"subtotal"`
	DosageInstruction string    `json:"dosage_instruction"`
	CreatedAt         time.Time `json:"created_at"`
}

type Prescription struct {
	ID             int64              `json:"id"`
	ConsultationID int64              `json:"consultation_id"`
	QueueID        int64              `json:"queue_id"`
	Status         string             `json:"status"` // RESERVED, PREPARING, DISPENSED, CANCELLED
	Notes          string             `json:"notes,omitempty"`
	Items          []PrescriptionItem `json:"items"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type Consultation struct {
	ID           int64         `json:"id"`
	QueueID      int64         `json:"queue_id"`
	DoctorID     *int64        `json:"doctor_id,omitempty"`
	DoctorName   string        `json:"doctor_name"`
	Anamnesis    string        `json:"anamnesis"`
	Diagnosis    string        `json:"diagnosis"`
	ICD10Code    string        `json:"icd10_code"`
	ActionNotes  string        `json:"action_notes,omitempty"`
	ActionFee    int           `json:"action_fee"`
	Status       string        `json:"status"` // DRAFT, COMPLETED
	Prescription *Prescription `json:"prescription,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type RxItemInput struct {
	DrugID            *int64 `json:"drug_id,omitempty"`
	ItemName          string `json:"item_name"`
	IsManual          bool   `json:"is_manual"`
	Quantity          int    `json:"quantity"`
	UnitPrice         int    `json:"unit_price"`
	DosageInstruction string `json:"dosage_instruction"`
}

type SaveConsultationInput struct {
	QueueID           int64         `json:"queue_id"`
	DoctorID          *int64        `json:"doctor_id,omitempty"`
	DoctorName        string        `json:"doctor_name"`
	Anamnesis         string        `json:"anamnesis"`
	Diagnosis         string        `json:"diagnosis"`
	ICD10Code         string        `json:"icd10_code"`
	ActionNotes       string        `json:"action_notes"`
	ActionFee         int           `json:"action_fee"`
	PrescriptionItems []RxItemInput `json:"prescription_items"`
}

type ConsultationRepository interface {
	ListDoctorQueues(ctx context.Context, doctorName, clinicUnit string) ([]*DoctorQueuePatient, error)
	GetQueuePatientByID(ctx context.Context, queueID int64) (*DoctorQueuePatient, error)
	UpdateQueueStatus(ctx context.Context, queueID int64, status string) error
	GetConsultationByQueueID(ctx context.Context, queueID int64) (*Consultation, error)
	SaveConsultationTx(ctx context.Context, c *Consultation, items []PrescriptionItem) error
	ListActiveDrugs(ctx context.Context, clinicUnit string) ([]*CatalogDrug, error)
	GetDrugByID(ctx context.Context, id int64) (*CatalogDrug, error)
}
