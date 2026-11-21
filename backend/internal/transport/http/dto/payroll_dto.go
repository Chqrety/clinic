package dto

import (
	"clinic-api/internal/domain/payroll"
	"clinic-api/pkg/decimal"
)

type CreateEmployeeRequest struct {
	NIK                    string          `json:"nik"`
	Name                   string          `json:"name"`
	Position               string          `json:"position"`
	Department             string          `json:"department"`
	BusinessUnit           string          `json:"businessUnit"`
	EmploymentType         string          `json:"employmentType"`
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
}

type UpdateEmployeeRequest struct {
	Name                   string          `json:"name"`
	Position               string          `json:"position"`
	Department             string          `json:"department"`
	BusinessUnit           string          `json:"businessUnit"`
	EmploymentType         string          `json:"employmentType"`
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
}

type GeneratePayrollRequest struct {
	Month        int    `json:"month"`
	Year         int    `json:"year"`
	BusinessUnit string `json:"businessUnit,omitempty"`
}

type ApprovePayslipRequest struct {
	ApprovedBy string `json:"approvedBy"`
}

type ApproveAllPayslipsRequest struct {
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	ApprovedBy string `json:"approvedBy"`
}

type MarkAllPaidRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type PostPayrollLedgerRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type PayslipResponse struct {
	payroll.Payslip
}
