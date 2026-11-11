package clinicalunit

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UnitCode identifies the dual-brand operational clinic units.
type UnitCode string

const (
	UnitKidinest UnitCode = "KIDINEST" // Kidinest (Poli Anak, Vaksinasi, Tumbuh Kembang)
	UnitSSC      UnitCode = "SSC"      // Semarang Skin Center (Estetika & Dermatologi)
)

// UnitType categorizes clinical focus.
type UnitType string

const (
	UnitTypePediatric          UnitType = "PEDIATRIC_CARE"
	UnitTypeAestheticDerm      UnitType = "AESTHETIC_DERMATOLOGY"
)

// PaymentTiming enforces unit-specific cashier billing flow.
type PaymentTiming string

const (
	// PaymentTimingPreService: Kidinest flow - Bayar pendaftaran & tindakan di awal sebelum nomor antrian dokter terbit.
	PaymentTimingPreService PaymentTiming = "PRE_SERVICE"
	// PaymentTimingPostService: SSC flow - Reservasi/booking jadwal Rp0, pembayaran di kasir setelah tindakan medis/treatment.
	PaymentTimingPostService PaymentTiming = "POST_SERVICE"
)

// ServiceCategory defines clinical services supported by Kidinest and Semarang Skin Center.
type ServiceCategory string

const (
	// Kidinest Categories (Anak & Tumbuh Kembang)
	CategoryVaccine          ServiceCategory = "VACCINE"           // Kategori A: Layanan Vaksinasi Anak
	CategoryPediatricIllness ServiceCategory = "PEDIATRIC_ILLNESS" // Kategori B: Pasien Sakit (Poli Umum Anak)
	CategoryGrowthMonitor    ServiceCategory = "GROWTH_MONITOR"    // Kategori C: Pemantauan Tumbuh Kembang (Denver II/KMS)
	CategoryGrowthPackage    ServiceCategory = "GROWTH_PACKAGE"    // Kategori D: Paket Vaksinasi & Tumbuh Kembang
	CategoryPediatricWalkIn  ServiceCategory = "PEDIATRIC_WALK_IN" // Kategori E: Walk-In / Datang Langsung

	// Semarang Skin Center (SSC) Categories (Estetika & Kulit)
	CategoryAestheticConsult ServiceCategory = "AESTHETIC_CONSULT" // Konsultasi Dokter Spesialis Kulit (Sp.DVE)
	CategoryFacialTreatment  ServiceCategory = "FACIAL_TREATMENT"  // Facial & Perawatan Kulit
	CategoryLaserPeeling     ServiceCategory = "LASER_PEELING"     // Tindakan Laser, Cauter, Chemical Peeling
	CategorySkincareRetail   ServiceCategory = "SKINCARE_RETAIL"   // Penjualan Produk Skincare & Kosmeceutical
)

// ClinicalUnit represents a clinic operational unit (Kidinest / Semarang Skin Center).
type ClinicalUnit struct {
	ID                  uuid.UUID         `json:"id"`
	Code                UnitCode          `json:"code"`
	Name                string            `json:"name"`
	Type                UnitType          `json:"type"`
	Tagline             string            `json:"tagline"`
	Address             string            `json:"address"`
	Phone               string            `json:"phone"`
	PaymentTiming       PaymentTiming     `json:"payment_timing"`
	DefaultCashAccount  string            `json:"default_cash_account"`  // e.g. 1-1111 (Kidinest) or 1-1112 (SSC)
	DefaultReceivableAcc string           `json:"default_receivable_acc"` // 1-1200
	SupportedCategories []ServiceCategory `json:"supported_categories"`
	IsActive            bool              `json:"is_active"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

var (
	ErrUnitNotFound       = errors.New("clinical unit not found")
	ErrInvalidCategory    = errors.New("invalid service category for unit")
	ErrPrePaymentRequired = errors.New("payment must be settled prior to service queue issuance")
)

// Static definitions for Kidinest and Semarang Skin Center
var (
	KidinestUnitID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	SSCUnitID      = uuid.MustParse("a0000000-0000-0000-0000-000000000002")

	KidinestUnit = ClinicalUnit{
		ID:                  KidinestUnitID,
		Code:                UnitKidinest,
		Name:                "Kidinest Pediatric & Child Development Center",
		Type:                UnitTypePediatric,
		Tagline:             "Klinik Spesialis Tumbuh Kembang, Vaksinasi & Kesehatan Anak",
		Address:             "Semarang, Jawa Tengah",
		Phone:               "024-8888001",
		PaymentTiming:       PaymentTimingPreService,
		DefaultCashAccount:  "1-1111", // Kas Kasir Kidinest
		DefaultReceivableAcc: "1-1200",
		SupportedCategories: []ServiceCategory{
			CategoryVaccine,
			CategoryPediatricIllness,
			CategoryGrowthMonitor,
			CategoryGrowthPackage,
			CategoryPediatricWalkIn,
		},
		IsActive: true,
	}

	SSCUnit = ClinicalUnit{
		ID:                  SSCUnitID,
		Code:                UnitSSC,
		Name:                "Semarang Skin Center (SSC)",
		Type:                UnitTypeAestheticDerm,
		Tagline:             "Pusat Dermatologi, Perawatan Kulit Medis & Estetika",
		Address:             "Semarang, Jawa Tengah",
		Phone:               "024-8888002",
		PaymentTiming:       PaymentTimingPostService,
		DefaultCashAccount:  "1-1112", // Kas Kasir Semarang Skin Center
		DefaultReceivableAcc: "1-1200",
		SupportedCategories: []ServiceCategory{
			CategoryAestheticConsult,
			CategoryFacialTreatment,
			CategoryLaserPeeling,
			CategorySkincareRetail,
		},
		IsActive: true,
	}
)

// GetUnitByCode returns clinical unit by code.
func GetUnitByCode(code UnitCode) (*ClinicalUnit, error) {
	switch code {
	case UnitKidinest:
		return &KidinestUnit, nil
	case UnitSSC:
		return &SSCUnit, nil
	default:
		return nil, ErrUnitNotFound
	}
}

// GetRevenueAccount maps unit and service category to the corresponding Chart of Accounts (COA) code.
func GetRevenueAccount(unitCode UnitCode, cat ServiceCategory) string {
	if unitCode == UnitKidinest {
		switch cat {
		case CategoryVaccine:
			return "4-1300" // Pendapatan Vaksinasi & Imunisasi Anak
		case CategoryPediatricIllness:
			return "4-1200" // Pendapatan Jasa Dokter Spesialis Anak (Sp.A)
		case CategoryGrowthMonitor, CategoryGrowthPackage:
			return "4-1400" // Pendapatan Pemantauan & Paket Tumbuh Kembang
		case CategoryPediatricWalkIn:
			return "4-1100" // Pendapatan Pendaftaran & Skrining Anak
		default:
			return "4-1100"
		}
	}

	if unitCode == UnitSSC {
		switch cat {
		case CategoryAestheticConsult:
			return "4-2100" // Pendapatan Konsultasi Dokter Kulit (Sp.DVE)
		case CategoryFacialTreatment:
			return "4-2200" // Pendapatan Treatment Facial & Estetika
		case CategoryLaserPeeling:
			return "4-2300" // Pendapatan Tindakan Laser, Cauter & Peeling
		case CategorySkincareRetail:
			return "4-2400" // Pendapatan Produk Skincare & Sunscreen
		default:
			return "4-2100"
		}
	}

	return "4-1100"
}
