package domain

import "errors"

var (
	ErrInvoiceNotFound      = errors.New("invoice tidak ditemukan")
	ErrQueueNotFound        = errors.New("antrian pasien tidak ditemukan")
	ErrAlreadyPaid          = errors.New("tagihan ini sudah dibayar")
	ErrInvalidAmount        = errors.New("nominal pembayaran tidak valid")
	ErrInsufficientCash     = errors.New("uang tunai yang diterima kurang dari total tagihan")
	ErrVoucherExceedsTotal  = errors.New("diskon voucher tidak boleh melebihi total tagihan")
	ErrInvalidPaymentMethod = errors.New("metode pembayaran tidak valid")
	ErrAlreadyRefunded      = errors.New("invoice ini sudah pernah di-refund")
	ErrRefundReasonRequired = errors.New("rekening bank dan alasan refund wajib diisi")
	ErrPatientNameRequired  = errors.New("nama pasien walk-in wajib diisi")
	ErrPatientPhoneRequired = errors.New("nomor whatsapp walk-in wajib diisi")
)
