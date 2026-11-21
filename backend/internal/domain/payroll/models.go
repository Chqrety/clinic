package payroll

import (
	"context"
	"time"

	"clinic-api/pkg/decimal"
)

type PayrollStatus string

const (
	StatusDraft    PayrollStatus = "DRAFT"
	StatusApproved PayrollStatus = "APPROVED"
	StatusPaid     PayrollStatus = "PAID"
)

type Employee struct {
	ID                     string          `json:"id"`
	NIK                    string          `json:"nik"`
	Name                   string          `json:"name"`
	Position               string          `json:"position"`
	Department             string          `json:"department"`
	BusinessUnit           string          `json:"businessUnit"` // KIDINEST, SSC, ALL, MINIMO
	EmploymentType         string          `json:"employmentType"` // TETAP, KONTRAK, HONORER, MAGANG
	JoinDate               string          `json:"joinDate"`
	BaseSalary             decimal.Decimal `json:"baseSalary"`
	PositionAllowance      decimal.Decimal `json:"positionAllowance"`
	TransportMealAllowance decimal.Decimal `json:"transportMealAllowance"`
	BankName               string          `json:"bankName"`
	BankAccountNumber      string          `json:"bankAccountNumber"`
	NPWP                   string          `json:"npwp"`
	PTKPStatus             string          `json:"ptkpStatus"`
	BPJSKesehatanNo        string          `json:"bpjsKesehatanNo"`
	BPJSKetenagakerjaanNo  string          `json:"bpjsKetenagakerjaanNo"`
	IsActive               bool            `json:"isActive"`
	CreatedAt              time.Time       `json:"createdAt"`
	UpdatedAt              time.Time       `json:"updatedAt"`
}

type Payslip struct {
	ID                        string          `json:"id"`
	SlipNumber                string          `json:"slipNumber"`
	EmployeeID                string          `json:"employeeId"`
	EmployeeNIK               string          `json:"employeeNik"`
	EmployeeName              string          `json:"employeeName"`
	Position                  string          `json:"position"`
	Department                string          `json:"department"`
	BusinessUnit              string          `json:"businessUnit"`
	PeriodMonth               int             `json:"periodMonth"`
	PeriodYear                int             `json:"periodYear"`
	PeriodName                string          `json:"periodName"`
	PaymentDate               string          `json:"paymentDate"`
	
	// Kehadiran / Hari Masuk (Attendance)
	WorkingDays               int             `json:"workingDays"`
	PresentDays               int             `json:"presentDays"`
	SickLeaveDays             int             `json:"sickLeaveDays"`
	PaidLeaveDays             int             `json:"paidLeaveDays"`
	AbsentDays                int             `json:"absentDays"`
	OvertimeHours             float64         `json:"overtimeHours"`
	DailyMealTransportRate    decimal.Decimal `json:"dailyMealTransportRate"`

	BaseSalary                decimal.Decimal `json:"baseSalary"`
	PositionAllowance         decimal.Decimal `json:"positionAllowance"`
	TransportMealAllowance    decimal.Decimal `json:"transportMealAllowance"`
	ShiftAttendanceAllowance  decimal.Decimal `json:"shiftAttendanceAllowance"`
	IncentiveBonus            decimal.Decimal `json:"incentiveBonus"`
	OvertimePay               decimal.Decimal `json:"overtimePay"`
	OtherAllowance            decimal.Decimal `json:"otherAllowance"`
	TotalGrossSalary          decimal.Decimal `json:"totalGrossSalary"`
	BPJSKesehatan             decimal.Decimal `json:"bpjsKesehatan"`
	BPJSKetenagakerjaan       decimal.Decimal `json:"bpjsKetenagakerjaan"`
	PPh21                     decimal.Decimal `json:"pph21"`
	AttendanceDeduction       decimal.Decimal `json:"attendanceDeduction"`
	LoanDeduction             decimal.Decimal `json:"loanDeduction"`
	OtherDeduction            decimal.Decimal `json:"otherDeduction"`
	TotalDeduction            decimal.Decimal `json:"totalDeduction"`
	TakeHomePay               decimal.Decimal `json:"takeHomePay"`
	BankName                  string          `json:"bankName"`
	BankAccountNumber         string          `json:"bankAccountNumber"`
	Status                    PayrollStatus   `json:"status"`
	ApprovedBy                *string         `json:"approvedBy,omitempty"`
	ApprovedAt                *time.Time      `json:"approvedAt,omitempty"`
	PaidAt                    *time.Time      `json:"paidAt,omitempty"`
	JournalID                 *string         `json:"journalId,omitempty"`
	Notes                     string          `json:"notes,omitempty"`
	CreatedAt                 time.Time       `json:"createdAt"`
	UpdatedAt                 time.Time       `json:"updatedAt"`
}

type PayslipFilter struct {
	Month        int
	Year         int
	BusinessUnit string
	Status       string
	Search       string
}

type Repository interface {
	// Employee operations
	ListEmployees(ctx context.Context, businessUnit string, isActiveOnly bool) ([]Employee, error)
	GetEmployeeByID(ctx context.Context, id string) (*Employee, error)
	CreateEmployee(ctx context.Context, emp *Employee) error
	UpdateEmployee(ctx context.Context, emp *Employee) error

	// Payslip operations
	ListPayslips(ctx context.Context, filter PayslipFilter) ([]Payslip, error)
	GetPayslipByID(ctx context.Context, id string) (*Payslip, error)
	CreatePayslip(ctx context.Context, slip *Payslip) error
	UpdatePayslip(ctx context.Context, slip *Payslip) error
	DeletePayslip(ctx context.Context, id string) error
	ApprovePayslip(ctx context.Context, id string, approvedBy string) error
	ApproveAllInPeriod(ctx context.Context, month, year int, approvedBy string) error
	MarkPayslipAsPaid(ctx context.Context, id string) error
	MarkAllAsPaidInPeriod(ctx context.Context, month, year int) error
	UpdatePayslipJournalID(ctx context.Context, id string, journalID string) error
}
