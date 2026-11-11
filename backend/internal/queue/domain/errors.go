package domain

import "errors"

var (
	ErrQueueNotFound           = errors.New("nomor antrian tidak ditemukan")
	ErrPatientNotFound         = errors.New("data pasien tidak ditemukan")
	ErrPaymentAlreadyProcessed = errors.New("pembayaran antrian ini sudah diproses")
	ErrInvalidClinicUnit       = errors.New("unit klinik tidak valid (KIDINEST / SSC)")
	ErrInvalidCategory         = errors.New("kategori layanan antrian tidak valid")
	ErrInvalidPaymentMethod    = errors.New("metode pembayaran tidak valid")
	ErrEmptyPatientName        = errors.New("nama pasien tidak boleh kosong")
	ErrEmptyPatientPhone       = errors.New("nomor WhatsApp pasien tidak boleh kosong")
)
