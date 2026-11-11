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
		"[WHATSAPP GATEWAY -> %s] Halo %s, Tiket antrian Anda di %s adalah: *%s*. Pantau panggilan antrian live secara paperless di tautan berikut: %s",
		phone, patientName, clinicUnit, queueNumber, trackingURL,
	)
	log.Printf("%s", msg)
	return nil
}

func (a *MockWhatsAppAdapter) SendMessage(phone, message string) error {
	log.Printf("[WHATSAPP GATEWAY -> %s] %s", phone, message)
	return nil
}
