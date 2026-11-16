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
	PaymentEDC           PaymentMethod = "EDC"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PENDING"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentRefunded PaymentStatus = "REFUNDED"
)

type QueueStatus string

const (
	QueueBooked          QueueStatus = "BOOKED"
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

type PatientLookupResult struct {
	Patient      Patient    `json:"patient"`
	ActiveTicket *QueueItem `json:"active_ticket,omitempty"`
}

type QueueItem struct {
	ID                   int64         `json:"id"`
	PatientID            int64         `json:"patient_id"`
	PatientName          string        `json:"patient_name"`
	PatientPhone         string        `json:"patient_phone,omitempty"`
	PatientType          string        `json:"patient_type,omitempty"`
	ClinicUnit           ClinicUnit    `json:"clinic_unit"`
	Category             QueueCategory `json:"category"`
	DoctorID             *int64        `json:"doctor_id,omitempty"`
	DoctorName           string        `json:"doctor_name"`
	AppointmentDate      *string       `json:"appointment_date,omitempty"`
	AppointmentTime      *string       `json:"appointment_time,omitempty"`
	QueueNumber          string        `json:"queue_number"`
	RegistrationFee      int           `json:"registration_fee"`
	PaymentMethod        PaymentMethod `json:"payment_method"`
	PaymentStatus        PaymentStatus `json:"payment_status"`
	QueueStatus          QueueStatus   `json:"queue_status"`
	RescheduleCount      int           `json:"reschedule_count"`
	RefundStatus         *string       `json:"refund_status,omitempty"`
	RefundNotes          *string       `json:"refund_notes,omitempty"`
	RefundBank           *string       `json:"refund_bank,omitempty"`
	RefundAccountNumber  *string       `json:"refund_account_number,omitempty"`
	RefundAccountName    *string       `json:"refund_account_name,omitempty"`
	PaymentReference     string        `json:"payment_reference,omitempty"`
	PaymentURL           string        `json:"payment_url,omitempty"`
	PaidAt               *time.Time    `json:"paid_at,omitempty"`
	CalledAt             *time.Time    `json:"called_at,omitempty"`
	CompletedAt          *time.Time    `json:"completed_at,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
}

type RegisterInput struct {
	ClinicUnit      ClinicUnit    `json:"clinic_unit"`
	PatientName     string        `json:"patient_name"`
	PatientPhone    string        `json:"patient_phone"`
	PatientType     string        `json:"patient_type,omitempty"` // "BARU" | "LAMA"
	Category        QueueCategory `json:"category"`
	DoctorID        *int64        `json:"doctor_id,omitempty"`
	DoctorName      string        `json:"doctor_name"`
	AppointmentDate *string       `json:"appointment_date,omitempty"` // YYYY-MM-DD
	AppointmentTime *string       `json:"appointment_time,omitempty"` // HH:MM
	PaymentMethod   PaymentMethod `json:"payment_method"`
}

type RegisterResult struct {
	QueueID          int64         `json:"queue_id"`
	QueueNumber      string        `json:"queue_number"`
	ClinicUnit       ClinicUnit    `json:"clinic_unit"`
	Category         QueueCategory `json:"category"`
	DoctorName       string        `json:"doctor_name"`
	AppointmentDate  *string       `json:"appointment_date,omitempty"`
	AppointmentTime  *string       `json:"appointment_time,omitempty"`
	PatientType      string        `json:"patient_type,omitempty"`
	RegistrationFee  int           `json:"registration_fee"`
	PaymentMethod    PaymentMethod `json:"payment_method"`
	PaymentStatus    PaymentStatus `json:"payment_status"`
	QueueStatus      QueueStatus   `json:"queue_status"`
	PaymentReference string        `json:"payment_reference,omitempty"`
	PaymentURL       string        `json:"payment_url,omitempty"`
	SnapToken        string        `json:"snap_token,omitempty"`
	ClientKey        string        `json:"client_key,omitempty"`
	QRISPayload      string        `json:"qris_payload,omitempty"`
	BankTransferAcc  string        `json:"bank_transfer_account,omitempty"`
	Message          string        `json:"message"`
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

type DoctorRoomStatus struct {
	DoctorName     string      `json:"doctor_name"`
	ClinicUnit     ClinicUnit  `json:"clinic_unit"`
	RoomName       string      `json:"room_name"`
	Status         string      `json:"status"` // "READY" | "BUSY"
	CurrentPatient *QueueItem  `json:"current_patient,omitempty"`
	WaitingCount   int         `json:"waiting_count"`
	NextPatient    *QueueItem  `json:"next_patient,omitempty"`
}

type PatientRepository interface {
	FindOrCreateByPhone(ctx context.Context, name, phone string) (*Patient, error)
	GetPatientByID(ctx context.Context, id int64) (*Patient, error)
	LookupPatient(ctx context.Context, query string) (*PatientLookupResult, error)
}

type QueueRepository interface {
	Create(ctx context.Context, q *QueueItem) error
	GetByID(ctx context.Context, id int64) (*QueueItem, error)
	UpdatePaymentSuccess(ctx context.Context, id int64, queueNumber string) (*QueueItem, error)
	UpdatePaymentDetails(ctx context.Context, id int64, reference, paymentURL string) error
	GetByPaymentReference(ctx context.Context, reference string) (*QueueItem, error)
	UpdateQueueStatus(ctx context.Context, id int64, status QueueStatus) error
	GetCallingQueue(ctx context.Context, unit string) (*QueueItem, error)
	ListUpcomingQueues(ctx context.Context, unit string, limit int) ([]QueueItem, error)
	ListReservations(ctx context.Context, unit string) ([]QueueItem, error)
	ListActiveDoctorQueues(ctx context.Context) ([]QueueItem, error)
	GenerateNextQueueNumber(ctx context.Context, date time.Time, unit ClinicUnit, category QueueCategory) (string, error)
	Reschedule(ctx context.Context, id int64, newDate, newTime string) error
	RequestRefund(ctx context.Context, id int64, bank, accNumber, accName, notes string) error
	ProcessRefund(ctx context.Context, id int64, status string, notes string) error
	ListRefundRequests(ctx context.Context, unit string) ([]QueueItem, error)
}
