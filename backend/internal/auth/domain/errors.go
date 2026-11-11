package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("nomor induk karyawan atau password salah")
	ErrStaffInactive      = errors.New("akun staf telah dinonaktifkan oleh atasan")
	ErrStaffNotFound      = errors.New("data staf tidak ditemukan")
)
