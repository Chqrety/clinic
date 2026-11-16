package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Provider string // "fonnte", "generic", "mock"
	APIURL   string
	APIToken string
}

type Service struct {
	cfg        Config
	httpClient *http.Client
}

func NewService(cfg Config) *Service {
	provider := strings.ToLower(strings.TrimSpace(cfg.Provider))
	if provider == "" {
		provider = "mock"
	}
	cfg.Provider = provider

	if cfg.APIURL == "" && provider == "fonnte" {
		cfg.APIURL = "https://api.fonnte.com/send"
	}

	return &Service{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 6 * time.Second,
		},
	}
}

func (s *Service) IsLive() bool {
	return strings.TrimSpace(s.cfg.APIToken) != "" && s.cfg.Provider != "mock"
}

// NormalizeIndonesianPhone converts any variation (08..., +628..., 628..., 0812-3456-7890) into standard international 628...
func NormalizeIndonesianPhone(phone string) string {
	re := regexp.MustCompile(`[^0-9]`)
	cleaned := re.ReplaceAllString(phone, "")

	if strings.HasPrefix(cleaned, "08") {
		return "62" + cleaned[1:]
	}
	if strings.HasPrefix(cleaned, "8") {
		return "62" + cleaned
	}
	if strings.HasPrefix(cleaned, "628") {
		return cleaned
	}
	return cleaned
}

// FormatTicketMessage produces standard friendly notification text
func FormatTicketMessage(patientName, queueNumber, clinicUnit, trackingURL string) string {
	unitName := "Kidinest Pediatric Care (Lantai 1)"
	if strings.ToUpper(clinicUnit) == "SSC" {
		unitName = "Semarang Skin Center (Lantai 2)"
	}

	return fmt.Sprintf(
		"Halo Bpk/Ibu *%s*,\n\n"+
			"Terima kasih telah mendaftar di *%s*.\n"+
			"Nomor Antrian Resmi Anda: *%s*\n\n"+
			"Pantau pergerakan antrian live secara paperless melalui ponsel Anda di tautan berikut:\n%s\n\n"+
			"Mohon bersiap 10-15 menit sebelum nomor antrian Anda dipanggil. Terima kasih atas kerja sama Anda.",
		patientName, unitName, queueNumber, trackingURL,
	)
}

// SendQueueTicket formats the ticket message and dispatches it asynchronously
func (s *Service) SendQueueTicket(phone, patientName, queueNumber, clinicUnit, trackingURL string) error {
	msg := FormatTicketMessage(patientName, queueNumber, clinicUnit, trackingURL)
	return s.SendMessage(phone, msg)
}

// SendMessage delivers WhatsApp message via configured gateway provider
func (s *Service) SendMessage(phone, message string) error {
	normalizedPhone := NormalizeIndonesianPhone(phone)
	if normalizedPhone == "" {
		return fmt.Errorf("nomor telepon kosong atau tidak valid")
	}

	if !s.IsLive() {
		log.Printf("[WHATSAPP MOCK -> %s]:\n%s", normalizedPhone, message)
		return nil
	}

	// Asynchronous dispatch via goroutine so caller response is non-blocking
	go func(targetPhone, textMsg string) {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		var req *http.Request
		var err error

		if s.cfg.Provider == "fonnte" || strings.Contains(s.cfg.APIURL, "fonnte.com") {
			payload := map[string]string{
				"target":  targetPhone,
				"message": textMsg,
			}
			jsonBytes, _ := json.Marshal(payload)
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.APIURL, bytes.NewBuffer(jsonBytes))
			if err == nil {
				req.Header.Set("Authorization", s.cfg.APIToken)
				req.Header.Set("Content-Type", "application/json")
			}
		} else {
			// Generic Webhook Provider
			payload := map[string]string{
				"phone":   targetPhone,
				"message": textMsg,
			}
			jsonBytes, _ := json.Marshal(payload)
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.APIURL, bytes.NewBuffer(jsonBytes))
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+s.cfg.APIToken)
				req.Header.Set("Content-Type", "application/json")
			}
		}

		if err != nil {
			log.Printf("[WHATSAPP ERROR] Failed creating request for %s: %v", targetPhone, err)
			return
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			log.Printf("[WHATSAPP ERROR] Delivery failed for %s: %v", targetPhone, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("[WHATSAPP ERROR] Gateway rejected (%d): %s", resp.StatusCode, string(body))
			return
		}

		log.Printf("[WHATSAPP SUCCESS] Ticket notification delivered to %s", targetPhone)
	}(normalizedPhone, message)

	return nil
}
