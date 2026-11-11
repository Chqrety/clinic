package domain

import (
	"context"
	"time"
)

type PrescriptionItemDetail struct {
	ID                int64  `json:"id"`
	PrescriptionID    int64  `json:"prescription_id"`
	DrugID            *int64 `json:"drug_id,omitempty"`
	ItemName          string `json:"item_name"`
	IsManual          bool   `json:"is_manual"`
	Quantity          int    `json:"quantity"`
	UnitPrice         int    `json:"unit_price"`
	Subtotal          int    `json:"subtotal"`
	DosageInstruction string `json:"dosage_instruction"`
}

type PharmacyPrescription struct {
	ID             int64                    `json:"id"`
	ConsultationID int64                    `json:"consultation_id"`
	QueueID        int64                    `json:"queue_id"`
	QueueNumber    string                   `json:"queue_number"`
	PatientName    string                   `json:"patient_name"`
	ClinicUnit     string                   `json:"clinic_unit"`
	DoctorName     string                   `json:"doctor_name"`
	Status         string                   `json:"status"` // RESERVED, PREPARING, DISPENSED, CANCELLED
	Notes          string                   `json:"notes,omitempty"`
	Items          []PrescriptionItemDetail `json:"items"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

type DrugBatch struct {
	ID             int64     `json:"id"`
	DrugID         int64     `json:"drug_id"`
	BatchNumber    string    `json:"batch_number"`
	ExpiredDate    string    `json:"expired_date"`
	PhysicalStock  int       `json:"physical_stock"`
	ReservedStock  int       `json:"reserved_stock"`
	AvailableStock int       `json:"available_stock"`
	IsExpiringSoon bool      `json:"is_expiring_soon"` // ED < 90 hari
	CreatedAt      time.Time `json:"created_at"`
}

type Drug struct {
	ID              int64       `json:"id"`
	ClinicUnit      string      `json:"clinic_unit"`
	Name            string      `json:"name"`
	KFACode         string      `json:"kfa_code,omitempty"`
	Unit            string      `json:"unit"`
	Category        string      `json:"category"`
	MinStockAlert   int         `json:"min_stock_alert"`
	BuyPrice        int         `json:"buy_price"`
	SellPrice       int         `json:"sell_price"`
	PhysicalStock   int         `json:"physical_stock"`
	ReservedStock   int         `json:"reserved_stock"`
	AvailableStock  int         `json:"available_stock"`
	IsLowStock      bool        `json:"is_low_stock"`
	HasExpiringSoon bool        `json:"has_expiring_soon"`
	Batches         []DrugBatch `json:"batches,omitempty"`
	IsActive        bool        `json:"is_active"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type CreateDrugInput struct {
	ClinicUnit    string `json:"clinic_unit"`
	Name          string `json:"name"`
	KFACode       string `json:"kfa_code"`
	Unit          string `json:"unit"`
	Category      string `json:"category"`
	MinStockAlert int    `json:"min_stock_alert"`
	BuyPrice      int    `json:"buy_price"`
	SellPrice     int    `json:"sell_price"`
}

type RestockInput struct {
	DrugID      int64  `json:"drug_id"`
	BatchNumber string `json:"batch_number"`
	ExpiredDate string `json:"expired_date"`
	Quantity    int    `json:"quantity"`
	BuyPrice    int    `json:"buy_price"`
}

type StockOpnameItem struct {
	DrugID      int64  `json:"drug_id"`
	DrugName    string `json:"drug_name"`
	BatchID     int64  `json:"batch_id"`
	BatchNumber string `json:"batch_number"`
	ExpiredDate string `json:"expired_date"`
	SystemStock int    `json:"system_stock"`
	ActualStock int    `json:"actual_stock"`
	Discrepancy int    `json:"discrepancy"` // actual - system
	Notes       string `json:"notes,omitempty"`
}

type StockOpnameLog struct {
	ID          int64     `json:"id"`
	DrugID      int64     `json:"drug_id"`
	DrugName    string    `json:"drug_name"`
	BatchID     *int64    `json:"batch_id,omitempty"`
	BatchNumber string    `json:"batch_number,omitempty"`
	SystemStock int       `json:"system_stock"`
	ActualStock int       `json:"actual_stock"`
	Discrepancy int       `json:"discrepancy"`
	AdjustedBy  string    `json:"adjusted_by"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type PharmacyRepository interface {
	ListIncomingPrescriptions(ctx context.Context, clinicUnit string) ([]*PharmacyPrescription, error)
	GetPrescriptionByID(ctx context.Context, id int64) (*PharmacyPrescription, error)
	UpdatePrescriptionStatus(ctx context.Context, id int64, status string) error
	DispensePrescriptionTx(ctx context.Context, id int64) error
	CancelPrescriptionTx(ctx context.Context, id int64) error
	ListDrugs(ctx context.Context, clinicUnit string) ([]*Drug, error)
	GetDrugByID(ctx context.Context, id int64) (*Drug, error)
	CreateDrug(ctx context.Context, in *CreateDrugInput) (*Drug, error)
	RestockBatch(ctx context.Context, in *RestockInput) error
	ApplyStockOpnameTx(ctx context.Context, items []StockOpnameItem, adjustedBy string) error
	ListOpnameLogs(ctx context.Context, limit int) ([]*StockOpnameLog, error)
}
