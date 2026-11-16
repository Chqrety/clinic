package tripay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestGenerateRequestSignature(t *testing.T) {
	cfg := Config{
		MerchantCode: "T12345",
		PrivateKey:   "secret-key-123",
	}
	client := NewClient(cfg)

	merchantRef := "KID-001"
	amount := 50000

	sig := client.GenerateRequestSignature(merchantRef, amount)

	// Manual calculation
	expectedPayload := "T12345KID-00150000"
	h := hmac.New(sha256.New, []byte("secret-key-123"))
	h.Write([]byte(expectedPayload))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if sig != expectedSig {
		t.Fatalf("expected signature %s, got %s", expectedSig, sig)
	}
}

func TestValidateCallbackSignature(t *testing.T) {
	cfg := Config{
		PrivateKey: "my-webhook-secret",
	}
	client := NewClient(cfg)

	rawBody := []byte(`{"reference":"DEV-T001","status":"PAID"}`)

	h := hmac.New(sha256.New, []byte("my-webhook-secret"))
	h.Write(rawBody)
	validSig := hex.EncodeToString(h.Sum(nil))

	// 1. Test Valid Signature
	if !client.ValidateCallbackSignature(rawBody, validSig) {
		t.Errorf("expected valid signature to pass verification")
	}

	// 2. Test Invalid/Tampered Signature
	if client.ValidateCallbackSignature(rawBody, "invalid-tampered-signature") {
		t.Errorf("expected invalid signature to fail verification")
	}

	// 3. Test Tampered Body
	tamperedBody := []byte(`{"reference":"DEV-T001","status":"PAID","amount":999999}`)
	if client.ValidateCallbackSignature(tamperedBody, validSig) {
		t.Errorf("expected tampered body to fail verification")
	}
}
