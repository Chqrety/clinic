package domain

import (
	"context"
	"strings"
	"time"
)

type ClinicUnit string

const (
	UnitKidinest ClinicUnit = "KIDINEST"
	UnitSSC      ClinicUnit = "SSC"
)

type QueueCategory string

const (
	CategoryVaksin    QueueCategory = "A_VAKSIN"
	CategorySakit     QueueCategory = "B_SAKIT"
	CategoryTumbang   QueueCategory = "C_TUMBANG"
	CategoryPaket     QueueCategory = "D_PAKET"
	CategoryWalkin    QueueCategory = "E_WALKIN"
	CategoryTreatment QueueCategory = "SSC_TREATMENT"
)

type PaymentMethod string

const (
	PaymentQRIS          PaymentMethod = "QRIS"
	PaymentBankTransfer  PaymentMethod = "BANK_TRANSFER"
	PaymentCash          PaymentMethod = "CASH"
	PaymentUnpaidBooking PaymentMethod = "UNPAID_BOOKING"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PENDING"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentRefunded PaymentStatus = "REFUNDED"
)

type QueueStatus string

const (
	QueueWaiting         QueueStatus = "WAITING"
	QueueCalling         QueueStatus = "CALLING"
	QueueInExam          QueueStatus = "IN_EXAM"
	QueuePharmacyWaiting QueueStatus = "PHARMACY_WAITING"
	QueueCompleted       QueueStatus = "COMPLETED"
	QueueCancelled       QueueStatus = "CANCELLED"
)

type Patient struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type QueueItem struct {
	ID              int64         `json:"id"`
	PatientID       int64         `json:"patient_id"`
	PatientName     string        `json:"patient_name"`
	PatientPhone    string        `json:"patient_phone,omitempty"`
	ClinicUnit      ClinicUnit    `json:"clinic_unit"`
	Category        QueueCategory `json:"category"`
	DoctorID        *int64        `json:"doctor_id,omitempty"`
	DoctorName      string        `json:"doctor_name"`
	QueueNumber     string        `json:"queue_number"`
	RegistrationFee int           `json:"registration_fee"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	QueueStatus     QueueStatus   `json:"queue_status"`
	PaidAt          *time.Time    `json:"paid_at,omitempty"`
	CalledAt        *time.Time    `json:"called_at,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
}

type RegisterInput struct {
	ClinicUnit    ClinicUnit    `json:"clinic_unit"`
	PatientName   string        `json:"patient_name"`
	PatientPhone  string        `json:"patient_phone"`
	Category      QueueCategory `json:"category"`
	DoctorID      *int64        `json:"doctor_id,omitempty"`
	DoctorName    string        `json:"doctor_name"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type RegisterResult struct {
	QueueID         int64         `json:"queue_id"`
	QueueNumber     string        `json:"queue_number"`
	ClinicUnit      ClinicUnit    `json:"clinic_unit"`
	Category        QueueCategory `json:"category"`
	DoctorName      string        `json:"doctor_name"`
	RegistrationFee int           `json:"registration_fee"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	QueueStatus     QueueStatus   `json:"queue_status"`
	QRISPayload     string        `json:"qris_payload,omitempty"`
	BankTransferAcc string        `json:"bank_transfer_account,omitempty"`
	Message         string        `json:"message"`
}

type MonitorResponse struct {
	CurrentlyCalling *QueueItem  `json:"currently_calling"`
	UpcomingQueue    []QueueItem `json:"upcoming_queue"`
	ClinicUnit       string      `json:"clinic_unit"`
}

func MaskPatientName(fullName string) string {
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) == 0 {
		return "Pasien"
	}
	if len(parts) == 1 {
		word := parts[0]
		if len(word) <= 2 {
			return word + "***"
		}
		return string([]rune(word)[:2]) + strings.Repeat("*", len([]rune(word))-2)
	}

	// First word intact, second word masked: "Budi P****"
	first := parts[0]
	second := parts[1]
	maskedSecond := string([]rune(second)[:1]) + strings.Repeat("*", 4)
	return first + " " + maskedSecond
}

type PatientRepository interface {
	FindOrCreateByPhone(ctx context.Context, name, phone string) (*Patient, error)
	GetPatientByID(ctx context.Context, id int64) (*Patient, error)
}

type QueueRepository interface {
	Create(ctx context.Context, q *QueueItem) error
	GetByID(ctx context.Context, id int64) (*QueueItem, error)
	UpdatePaymentSuccess(ctx context.Context, id int64, queueNumber string) (*QueueItem, error)
	UpdateQueueStatus(ctx context.Context, id int64, status QueueStatus) error
	GetCallingQueue(ctx context.Context, unit string) (*QueueItem, error)
	ListUpcomingQueues(ctx context.Context, unit string, limit int) ([]QueueItem, error)
	GenerateNextQueueNumber(ctx context.Context, date time.Time, unit ClinicUnit, category QueueCategory) (string, error)
}
