package billing

import (
	"time"

	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type InvoiceType string

const (
	InvoiceTypeRegistration InvoiceType = "REGISTRATION"
	InvoiceTypeConsultation InvoiceType = "CONSULTATION"
	InvoiceTypePharmacy     InvoiceType = "PHARMACY"
	InvoiceTypeCombined     InvoiceType = "COMBINED"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid    PaymentStatus = "UNPAID"
	PaymentStatusPartial   PaymentStatus = "PARTIAL"
	PaymentStatusPaid      PaymentStatus = "PAID"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

// Invoice represents a clinic sales invoice for patient services and pharmacy.
type Invoice struct {
	shared.BaseEntity
	shared.BranchScoped
	InvoiceNumber   string          `json:"invoice_number"`
	PatientID       *uuid.UUID      `json:"patient_id,omitempty"`
	PatientName     string          `json:"patient_name"`
	PatientNIK      string          `json:"patient_nik"`
	QueueID         *int64          `json:"queue_id,omitempty"`
	InvoiceType     InvoiceType     `json:"invoice_type"`
	RegistrationFee decimal.Decimal `json:"registration_fee"`
	ConsultationFee decimal.Decimal `json:"consultation_fee"`
	MedicationFee   decimal.Decimal `json:"medication_fee"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	PaidAmount      decimal.Decimal `json:"paid_amount"`
	PaymentStatus   PaymentStatus   `json:"payment_status"`
	PaymentMethod   string          `json:"payment_method,omitempty"` // CASH, BANK_TRANSFER, RECEIVABLE
	JournalID       *uuid.UUID      `json:"journal_id,omitempty"`     // Audit drill-down link
	Items           []InvoiceItem   `json:"items,omitempty"`
	Payments        []PaymentRecord `json:"payments,omitempty"`
}

type ItemType string

const (
	ItemTypeRegistration ItemType = "REGISTRATION"
	ItemTypeConsultation ItemType = "CONSULTATION"
	ItemTypeMedication   ItemType = "MEDICATION"
	ItemTypeProcedure    ItemType = "PROCEDURE"
)

// InvoiceItem represents detailed items inside an invoice.
type InvoiceItem struct {
	ID          uuid.UUID       `json:"id"`
	InvoiceID   uuid.UUID       `json:"invoice_id"`
	ItemType    ItemType        `json:"item_type"`
	ItemID      *uuid.UUID      `json:"item_id,omitempty"` // ID produk obat atau tindakan
	Description string          `json:"description"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	TotalPrice  decimal.Decimal `json:"total_price"`
}

// PaymentRecord represents a verified payment receipt allocated to an invoice.
type PaymentRecord struct {
	ID            uuid.UUID       `json:"id"`
	InvoiceID     uuid.UUID       `json:"invoice_id"`
	PaymentNumber string          `json:"payment_number"`
	Amount        decimal.Decimal `json:"amount"`
	PaymentMethod string          `json:"payment_method"` // CASH, BANK_TRANSFER, RECEIVABLE
	ReferenceNo   string          `json:"reference_no,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	ReceivedBy    *uuid.UUID      `json:"received_by,omitempty"`
	JournalID     *uuid.UUID      `json:"journal_id,omitempty"` // Linked accounting journal
	PaidAt        time.Time       `json:"paid_at"`
}

