package domain

import "errors"

var (
	ErrQueueNotFound           = errors.New("antrian pasien tidak ditemukan")
	ErrConsultationNotFound    = errors.New("data konsultasi belum ditemukan")
	ErrDoctorNameRequired      = errors.New("nama dokter pemeriksa wajib diisi")
	ErrDiagnosisRequired       = errors.New("diagnosa dokter wajib diisi")
	ErrEmptyPrescriptionItems  = errors.New("item resep obat tidak boleh kosong")
	ErrInvalidQuantity         = errors.New("jumlah obat harus lebih dari 0")
	ErrInsufficientStock       = errors.New("stok obat apotek tidak mencukupi untuk di-reserve")
	ErrDrugNotFound            = errors.New("obat katalog tidak ditemukan")
	ErrPatientNotReadyForExam  = errors.New("pasien belum terverifikasi untuk diperiksa dokter")
)
