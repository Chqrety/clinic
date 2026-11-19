package usecase

import (
	"fmt"
	"log"
)

type WhatsAppAdapter interface {
	SendQueueTicket(phone, patientName, queueNumber, clinicUnit, trackingURL string) error
	SendMessage(phone, message string) error
}

type MockWhatsAppAdapter struct{}

func NewMockWhatsAppAdapter() *MockWhatsAppAdapter {
	return &MockWhatsAppAdapter{}
}

func (a *MockWhatsAppAdapter) SendQueueTicket(phone, patientName, queueNumber, clinicUnit, trackingURL string) error {
	msg := fmt.Sprintf(
		"[WHATSAPP GATEWAY -> %s] Halo %s, Tiket pendaftaran antrian Anda di %s berhasil diterbitkan dengan nomor: *%s*.\n\nKonfirmasi nomor antrean dapat dilakukan melalui Customer Service pada hari kunjungan. Nomor antrean juga dapat dipantau secara langsung melalui sistem ini pada halaman utama: %s",
		phone, patientName, clinicUnit, queueNumber, trackingURL,
	)
	log.Printf("%s", msg)
	return nil
}

func (a *MockWhatsAppAdapter) SendMessage(phone, message string) error {
	log.Printf("[WHATSAPP GATEWAY -> %s] %s", phone, message)
	return nil
}
