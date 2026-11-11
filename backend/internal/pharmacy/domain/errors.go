package domain

import "errors"

var (
	ErrPrescriptionNotFound   = errors.New("resep tidak ditemukan")
	ErrDrugNotFound           = errors.New("obat tidak ditemukan")
	ErrAlreadyDispensed       = errors.New("resep sudah pernah diserahkan")
	ErrAlreadyCancelled       = errors.New("resep sudah dibatalkan")
	ErrInvalidQuantity        = errors.New("kuantitas obat harus lebih dari 0")
	ErrInvalidBatchNumber     = errors.New("nomor batch obat wajib diisi")
	ErrInvalidExpiredDate     = errors.New("format tanggal kadaluarsa tidak valid (gunakan YYYY-MM-DD)")
	ErrEmptyDrugName          = errors.New("nama obat wajib diisi")
	ErrEmptyOpnameItems       = errors.New("tidak ada item stok opnam yang valid untuk diproses")
	ErrInsufficientStock      = errors.New("stok fisik tidak mencukupi untuk pemotongan")
)
