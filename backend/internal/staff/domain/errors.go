package domain

import "errors"

var (
	ErrStaffNotFound         = errors.New("data staf tidak ditemukan")
	ErrDuplicateEmployeeNum  = errors.New("nomor induk karyawan sudah digunakan")
	ErrInvalidRole           = errors.New("peran staf tidak valid")
	ErrInvalidClinicUnit     = errors.New("unit klinik tidak valid")
	ErrEmptyEmployeeNumber   = errors.New("nomor induk karyawan tidak boleh kosong")
	ErrEmptyName             = errors.New("nama staf tidak boleh kosong")
	ErrEmptyPassword         = errors.New("password tidak boleh kosong")
)
