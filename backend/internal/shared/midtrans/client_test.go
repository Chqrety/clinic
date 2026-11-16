package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestMidtransSignatureValidation(t *testing.T) {
	serverKey := "SB-Mid-server-testkey123"
	client := NewClient(Config{
		ServerKey:    serverKey,
		IsProduction: false,
	})

	orderID := "KID-42-1731500000"
	statusCode := "200"
	grossAmount := "50000.00"

	// SHA512(order_id + status_code + gross_amount + ServerKey)
	raw := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	expectedSig := hex.EncodeToString(hash[:])

	if !client.ValidateNotificationSignature(orderID, statusCode, grossAmount, expectedSig) {
		t.Errorf("expected valid signature, got invalid")
	}

	if client.ValidateNotificationSignature(orderID, statusCode, grossAmount, "invalidsignature123") {
		t.Errorf("expected invalid signature to fail")
	}
}

func TestMidtransPaymentSuccessStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		fraudStatus string
		expected    bool
	}{
		{"settlement is success", "settlement", "accept", true},
		{"settlement without fraudStatus is success", "settlement", "", true},
		{"capture with accept is success", "capture", "accept", true},
		{"capture with challenge is not success", "capture", "challenge", false},
		{"pending is not success", "pending", "", false},
		{"expire is not success", "expire", "", false},
		{"cancel is not success", "cancel", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NotificationPayload{
				TransactionStatus: tt.status,
				FraudStatus:       tt.fraudStatus,
			}
			if p.IsPaymentSuccess() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, p.IsPaymentSuccess())
			}
		})
	}
}
