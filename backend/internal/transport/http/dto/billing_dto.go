package dto

import (
	"time"

	"clinic-api/internal/domain/billing"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type CreateInvoiceItemRequest struct {
	ItemType    billing.ItemType `json:"item_type"`
	ItemID      *uuid.UUID       `json:"item_id,omitempty"`
	Description string           `json:"description"`
	Quantity    int              `json:"quantity"`
	UnitPrice   decimal.Decimal  `json:"unit_price"`
}

type CreateInvoiceRequest struct {
	PatientID       *uuid.UUID                 `json:"patient_id,omitempty"`
	PatientName     string                     `json:"patient_name"`
	PatientNIK      string                     `json:"patient_nik"`
	QueueID         *int64                     `json:"queue_id,omitempty"`
	InvoiceType     billing.InvoiceType        `json:"invoice_type"`
	RegistrationFee decimal.Decimal            `json:"registration_fee"`
	ConsultationFee decimal.Decimal            `json:"consultation_fee"`
	MedicationFee   decimal.Decimal            `json:"medication_fee"`
	Items           []CreateInvoiceItemRequest `json:"items,omitempty"`
}

type PayInvoiceRequest struct {
	Amount            decimal.Decimal `json:"amount"`
	PaymentMethod     string          `json:"payment_method"` // CASH, BANK_TRANSFER, RECEIVABLE, QRIS, EDC
	UnitCode          string          `json:"unit_code,omitempty"` // 'KIDINEST' or 'SSC'
	ReferenceNo       string          `json:"reference_no,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	DoctorName        string          `json:"doctor_name,omitempty"`
	DoctorShareAmount decimal.Decimal `json:"doctor_share_amount,omitempty"`
	COGSAmount        decimal.Decimal `json:"cogs_amount,omitempty"`
}

type InvoiceItemResponse struct {
	ID          uuid.UUID        `json:"id"`
	ItemType    billing.ItemType `json:"item_type"`
	Description string           `json:"description"`
	Quantity    int              `json:"quantity"`
	UnitPrice   decimal.Decimal  `json:"unit_price"`
	TotalPrice  decimal.Decimal  `json:"total_price"`
}

type PaymentRecordResponse struct {
	ID            uuid.UUID       `json:"id"`
	PaymentNumber string          `json:"payment_number"`
	Amount        decimal.Decimal `json:"amount"`
	PaymentMethod string          `json:"payment_method"`
	ReferenceNo   string          `json:"reference_no,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	JournalID     *uuid.UUID      `json:"journal_id,omitempty"`
	PaidAt        time.Time       `json:"paid_at"`
}

type InvoiceResponse struct {
	ID              uuid.UUID               `json:"id"`
	BranchID        uuid.UUID               `json:"branch_id"`
	InvoiceNumber   string                  `json:"invoice_number"`
	PatientID       *uuid.UUID              `json:"patient_id,omitempty"`
	PatientName     string                  `json:"patient_name"`
	PatientNIK      string                  `json:"patient_nik"`
	QueueID         *int64                  `json:"queue_id,omitempty"`
	InvoiceType     billing.InvoiceType     `json:"invoice_type"`
	RegistrationFee decimal.Decimal         `json:"registration_fee"`
	ConsultationFee decimal.Decimal         `json:"consultation_fee"`
	MedicationFee   decimal.Decimal         `json:"medication_fee"`
	TotalAmount     decimal.Decimal         `json:"total_amount"`
	PaidAmount      decimal.Decimal         `json:"paid_amount"`
	PaymentStatus   billing.PaymentStatus   `json:"payment_status"`
	PaymentMethod   string                  `json:"payment_method,omitempty"`
	JournalID       *uuid.UUID              `json:"journal_id,omitempty"`
	Items           []InvoiceItemResponse   `json:"items,omitempty"`
	Payments        []PaymentRecordResponse `json:"payments,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
}

