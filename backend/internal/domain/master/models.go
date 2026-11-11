package master

import (
	"time"

	"clinic-api/internal/domain/shared"

	"github.com/google/uuid"
)

// Branch represents clinic branch entity for multi-branch isolation.
type Branch struct {
	shared.BaseEntity
	Code     string `json:"code"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"is_active"`
}

type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleAdmin      UserRole = "ADMIN"
	RoleDoctor     UserRole = "DOCTOR"
	RoleNurse      UserRole = "NURSE"
	RolePharmacist UserRole = "PHARMACIST"
	RoleCashier    UserRole = "CASHIER"
	RoleAccountant UserRole = "ACCOUNTANT"
)

// User represents system user with RBAC and branch access.
type User struct {
	shared.BaseEntity
	BranchID     uuid.UUID `json:"branch_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
}

// Patient represents master patient record.
type Patient struct {
	shared.BaseEntity
	MedicalRecordNumber string     `json:"medical_record_number"` // No RM
	NIK                 string     `json:"nik"`
	FullName            string     `json:"full_name"`
	DateOfBirth         time.Time  `json:"date_of_birth"`
	Gender              string     `json:"gender"`
	Phone               string     `json:"phone"`
	Address             string     `json:"address"`
	InsuranceType       string     `json:"insurance_type"` // UMUM, BPJS, ASURANSI_SWASTA
	InsuranceNumber     string     `json:"insurance_number,omitempty"`
	IsActive            bool       `json:"is_active"`
}

// Supplier represents medicine and medical supplies vendor.
type Supplier struct {
	shared.BaseEntity
	Code             string `json:"code"`
	Name             string `json:"name"`
	ContactPerson    string `json:"contact_person"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	PaymentTermsDays int    `json:"payment_terms_days"` // e.g. 30 days credit term (Hutang AP)
	TaxNumber        string `json:"tax_number,omitempty"`
	IsActive         bool   `json:"is_active"`
}

