# 09. Rencana Integrasi WhatsApp Gateway Notification

## 1. Tujuan & Latar Belakang
Mengirimkan notifikasi tiket antrian digital dan tautan pemantau langsung (*live monitor URL*) ke nomor WhatsApp pasien secara instan saat pendaftaran (online maupun walk-in via GRO) berhasil diproses:
* Pasien tidak perlu menunggu di ruang tunggu yang padat; mereka bisa menunggu di mobil atau area sekitar sambil memantau antrian dari HP.
* Mengurangi beban pertanyaan pasien ke kasir/resepsionis mengenai sisa nomor antrian di depannya.
* Mewujudkan sistem antrian nirkertas (*paperless*) tanpa printer tiket fisik.

---

## 2. Invarian Domain & Keandalan Sistem
1. **Asynchronous Non-Blocking Execution:** Pengiriman pesan WhatsApp wajib dijalankan di background goroutine (`go waService.SendQueueTicket(...)`). Kegagalan koneksi atau latensi provider WA tidak boleh memperlambat respons pembuatan tiket antrian ke pasien.
2. **Graceful Degradation:** Jika kuota WA habis atau provider sedang gangguan, sistem mencatat error log tanpa membatalkan transaksi tiket antrian.
3. **Format Nomor Indonesia:** Normalisasi otomatis format nomor telepon pasien (mengubah awalan `08xx`, `+628xx`, atau `628xx` menjadi standar format internasional `628xx` sebelum dikirim ke API WA).

---

## 3. Langkah Implementasi Teknis

### Tahap 1: Konfigurasi Environment (`.env`)
* Tambahkan kredensial provider WhatsApp Gateway (kompatibel dengan Fonnte / Wablas / Waha / generic webhook):
  ```env
  WA_GATEWAY_PROVIDER=fonnte # pilihan: mock, fonnte, wablas, generic
  WA_GATEWAY_API_URL=https://api.fonnte.com/send
  WA_GATEWAY_API_TOKEN=your_secret_token
  CLINIC_PUBLIC_WEB_URL=https://antrian.klinikkita.com
  ```

### Tahap 2: Backend WhatsApp Service Adapter (`backend/internal/shared/whatsapp`)
* Definisikan Interface Service:
  ```go
  type WhatsAppService interface {
      SendQueueTicket(ctx context.Context, payload QueueTicketMessage) error
  }
  ```
* Implementasikan Adapter:
  1. `MockWhatsAppService`: Mencatat log payload ke terminal (aktif otomatis jika token kosong saat development lokal).
  2. `FonnteWhatsAppAdapter`: Mengirim HTTP POST JSON/Form ke API provider dengan header Authorization.
* Normalizer Nomor Telepon:
  ```go
  func NormalizeIndonesianPhone(phone string) string
  ```

### Tahap 3: Format Template Pesan Tiket Antrian
Pesan WhatsApp disusun dengan format bersih, jelas, dan ramah keluarga:
```text
Halo Bpk/Ibu {{.PatientName}},

Terima kasih telah mendaftar di {{.ClinicName}}. Berikut adalah detail tiket antrian Anda:

Nomor Antrian: *{{.QueueNumber}}*
Kategori: {{.CategoryLabel}}
Dokter Tujuan: {{.DoctorName}}
Estimasi Panggilan: {{.EstimatedTime}} WIB

Pantau pergerakan antrian Anda secara langsung tanpa perlu berdiri di ruang tunggu melalui tautan berikut:
{{.LiveQueueURL}}

Mohon hadir di klinik 10-15 menit sebelum nomor Anda dipanggil. Terima kasih atas kerja sama Anda.
```

### Tahap 4: Integrasi Pemanggilan di Usecase Antrian
* Di [queue_usecase.go](file:///home/alpi/Projects/clinic/backend/internal/queue/usecase/queue_usecase.go):
  * Pada method `CreateQueue` (untuk pendaftaran pasien mandiri setelah bayar lunas) dan `CreateWalkInQueue` (untuk bantuan walk-in oleh GRO meja depan).
  * Panggil goroutine `go uc.waService.SendQueueTicket(context.Background(), msg)`.

---

## 4. Rencana Pengujian (Verification)
1. **Unit Test:**
   - Test normalisasi format nomor HP (`08123...` -> `628123...`).
   - Test template message formatter menghasilkan string lengkap tanpa nil pointer.
2. **Integration Test:**
   - Daftarkan pasien baru dengan nomor WhatsApp uji coba.
   - Verifikasi payload HTTP request berhasil dikirim ke provider WA.
   - Verifikasi tiket antrian tetap sukses tersimpan di database jika request WA diputus/gagal.
