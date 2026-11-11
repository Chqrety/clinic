package domain

import (
	"context"
	"time"
)

type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "CASH"
	PaymentMethodBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentMethodEDC          PaymentMethod = "EDC"
	PaymentMethodQRIS         PaymentMethod = "QRIS"
)

type InvoiceStatus string

const (
	InvoiceStatusPaid     InvoiceStatus = "PAID"
	InvoiceStatusRefunded InvoiceStatus = "REFUNDED"
)

type Invoice struct {
	ID                   int64         `json:"id"`
	QueueID              int64         `json:"queue_id"`
	PatientID            int64         `json:"patient_id"`
	InvoiceNumber        string        `json:"invoice_number"`
	PatientName          string        `json:"patient_name"`
	PatientPhone         string        `json:"patient_phone"`
	ClinicUnit           string        `json:"clinic_unit"`
	QueueNumber          string        `json:"queue_number"`
	SubtotalRegistration int           `json:"subtotal_registration"`
	SubtotalAction       int           `json:"subtotal_action"`
	SubtotalMedication   int           `json:"subtotal_medication"`
	VoucherDiscount      int           `json:"voucher_discount"`
	TotalAmount          int           `json:"total_amount"`
	PaymentMethod        PaymentMethod `json:"payment_method"`
	CashReceived         int           `json:"cash_received"`
	CashChange           int           `json:"cash_change"`
	PaymentReference     string        `json:"payment_reference"`
	CashierStaffID       int64         `json:"cashier_staff_id"`
	CashierName          string        `json:"cashier_name"`
	Status               InvoiceStatus `json:"status"`
	RefundBankName       string        `json:"refund_bank_name,omitempty"`
	RefundBankAccount    string        `json:"refund_bank_account,omitempty"`
	RefundAccountName    string        `json:"refund_account_name,omitempty"`
	RefundReason         string        `json:"refund_reason,omitempty"`
	RefundedAt           *time.Time    `json:"refunded_at,omitempty"`
	Notes                string        `json:"notes,omitempty"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
}

type CashierShiftReport struct {
	ID                int64     `json:"id"`
	StaffID           int64     `json:"staff_id"`
	StaffName         string    `json:"staff_name"`
	ShiftDate         string    `json:"shift_date"`
	TotalTransactions int       `json:"total_transactions"`
	TotalCash         int       `json:"total_cash"`
	TotalTransfer     int       `json:"total_transfer"`
	TotalEDC          int       `json:"total_edc"`
	TotalVoucher      int       `json:"total_voucher"`
	TotalRefund       int       `json:"total_refund"`
	NetTotal          int       `json:"net_total"`
	Notes             string    `json:"notes,omitempty"`
	SentToOwnerAt     time.Time `json:"sent_to_owner_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type ShiftSummary struct {
	ShiftDate         string `json:"shift_date"`
	TotalTransactions int    `json:"total_transactions"`
	TotalCash         int    `json:"total_cash"`
	TotalTransfer     int    `json:"total_transfer"`
	TotalEDC          int    `json:"total_edc"`
	TotalVoucher      int    `json:"total_voucher"`
	TotalRefund       int    `json:"total_refund"`
	NetTotal          int    `json:"net_total"`
}

type PendingBill struct {
	QueueID              int64     `json:"queue_id"`
	QueueNumber          string    `json:"queue_number"`
	PatientID            int64     `json:"patient_id"`
	PatientName          string    `json:"patient_name"`
	PatientPhone         string    `json:"patient_phone"`
	ClinicUnit           string    `json:"clinic_unit"`
	Category             string    `json:"category"`
	DoctorName           string    `json:"doctor_name"`
	SubtotalRegistration int       `json:"subtotal_registration"`
	SubtotalAction       int       `json:"subtotal_action"`
	SubtotalMedication   int       `json:"subtotal_medication"`
	TotalDue             int       `json:"total_due"`
	BillType             string    `json:"bill_type"` // REGISTRATION_ONLY, FULL_VISIT
	QueueStatus          string    `json:"queue_status"`
	PaymentStatus        string    `json:"payment_status"`
	CreatedAt            time.Time `json:"created_at"`
}

type WalkInRequest struct {
	PatientName     string `json:"patient_name"`
	PatientPhone    string `json:"patient_phone"`
	ClinicUnit      string `json:"clinic_unit"` // KIDINEST or SSC
	DoctorID        int64  `json:"doctor_id"`
	DoctorName      string `json:"doctor_name"`
	PaymentMethod   string `json:"payment_method"` // CASH, QRIS, BANK_TRANSFER
	IsPaidImmediate bool   `json:"is_paid_immediate"`
}

type PayBillRequest struct {
	QueueID              int64         `json:"queue_id"`
	SubtotalRegistration int           `json:"subtotal_registration"`
	SubtotalAction       int           `json:"subtotal_action"`
	SubtotalMedication   int           `json:"subtotal_medication"`
	VoucherDiscount      int           `json:"voucher_discount"`
	PaymentMethod        PaymentMethod `json:"payment_method"`
	CashReceived         int           `json:"cash_received"`
	PaymentReference     string        `json:"payment_reference"`
	Notes                string        `json:"notes"`
}

type RefundRequest struct {
	RefundBankName    string `json:"refund_bank_name"`
	RefundBankAccount string `json:"refund_bank_account"`
	RefundAccountName string `json:"refund_account_name"`
	RefundReason      string `json:"refund_reason"`
}

type CloseShiftRequest struct {
	Notes string `json:"notes"`
}

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id int64) (*Invoice, error)
	GetByQueueID(ctx context.Context, queueID int64) (*Invoice, error)
	List(ctx context.Context, limit, offset int) ([]*Invoice, error)
	UpdateRefund(ctx context.Context, id int64, bankName, bankAccount, accountName, reason string, refundedAt time.Time) error
	GetShiftSummary(ctx context.Context, date string) (*ShiftSummary, error)
	GenerateInvoiceNumber(ctx context.Context) (string, error)
}

type ShiftRepository interface {
	CreateReport(ctx context.Context, report *CashierShiftReport) error
	ListReports(ctx context.Context, limit int) ([]*CashierShiftReport, error)
}

type CashierQueueRepository interface {
	GetQueueByID(ctx context.Context, id int64) (*PendingBill, error)
	ListPendingBills(ctx context.Context) ([]*PendingBill, error)
	UpdatePaymentAndQueueStatus(ctx context.Context, queueID int64, paymentStatus, queueStatus string) error
	RegisterWalkIn(ctx context.Context, req *WalkInRequest, queueNumber string, fee int) (int64, string, error)
}
