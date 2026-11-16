package whatsapp

import (
	"strings"
	"testing"
)

func TestNormalizeIndonesianPhone(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"081234567890", "6281234567890"},
		{"+6281234567890", "6281234567890"},
		{"6281234567890", "6281234567890"},
		{"0812-3456-7890", "6281234567890"},
		{"(0812) 3456 7890", "6281234567890"},
		{"81234567890", "6281234567890"},
		{"+62 819-0888-1122", "6281908881122"},
	}

	for _, tc := range testCases {
		got := NormalizeIndonesianPhone(tc.input)
		if got != tc.expected {
			t.Errorf("NormalizeIndonesianPhone(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestFormatTicketMessage(t *testing.T) {
	msg := FormatTicketMessage("Ahmad Dahlan", "A-005", "KIDINEST", "https://klinik.com/kiosk?ticket=10")

	if !strings.Contains(msg, "Ahmad Dahlan") {
		t.Errorf("expected message to contain patient name")
	}
	if !strings.Contains(msg, "A-005") {
		t.Errorf("expected message to contain queue number A-005")
	}
	if !strings.Contains(msg, "Kidinest Pediatric Care") {
		t.Errorf("expected message to contain Kidinest")
	}
	if !strings.Contains(msg, "https://klinik.com/kiosk?ticket=10") {
		t.Errorf("expected message to contain tracking URL")
	}
}

func TestService_MockMode(t *testing.T) {
	svc := NewService(Config{Provider: "mock"})
	if svc.IsLive() {
		t.Errorf("expected mock service to return IsLive() = false")
	}

	err := svc.SendQueueTicket("08123456789", "Budi", "A-001", "KIDINEST", "http://localhost")
	if err != nil {
		t.Errorf("unexpected error in mock mode: %v", err)
	}
}
