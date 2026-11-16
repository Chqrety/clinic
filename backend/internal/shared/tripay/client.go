package tripay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	APIKey       string
	PrivateKey   string
	MerchantCode string
	IsSandbox    bool
	CallbackURL  string
	ReturnURL    string
}

type Client struct {
	cfg        Config
	httpClient *http.Client
	baseURL    string
}

func NewClient(cfg Config) *Client {
	baseURL := "https://tripay.co.id/api"
	if cfg.IsSandbox {
		baseURL = "https://tripay.co.id/api-sandbox"
	}
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) IsConfigured() bool {
	return strings.TrimSpace(c.cfg.APIKey) != "" &&
		strings.TrimSpace(c.cfg.PrivateKey) != "" &&
		strings.TrimSpace(c.cfg.MerchantCode) != ""
}

type OrderItem struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreateTransactionRequest struct {
	Method        string      `json:"method"`
	MerchantRef   string      `json:"merchant_ref"`
	Amount        int         `json:"amount"`
	CustomerName  string      `json:"customer_name"`
	CustomerEmail string      `json:"customer_email"`
	CustomerPhone string      `json:"customer_phone"`
	OrderItems    []OrderItem `json:"order_items"`
	CallbackURL   string      `json:"callback_url,omitempty"`
	ReturnURL     string      `json:"return_url,omitempty"`
	ExpiredTime   int64       `json:"expired_time,omitempty"`
	Signature     string      `json:"signature"`
}

type TransactionResponse struct {
	Reference      string `json:"reference"`
	MerchantRef    string `json:"merchant_ref"`
	PaymentMethod  string `json:"payment_method"`
	PaymentName    string `json:"payment_name"`
	CustomerName   string `json:"customer_name"`
	Amount         int    `json:"amount"`
	TotalFee       int    `json:"total_fee"`
	AmountReceived int    `json:"amount_received"`
	CheckoutURL    string `json:"checkout_url"`
	Status         string `json:"status"`
	QRString       string `json:"qr_string"`
	QRURL          string `json:"qr_url"`
	ExpiredTime    int64  `json:"expired_time"`
}

type apiResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    TransactionResponse `json:"data"`
}

type CallbackPayload struct {
	Reference   string `json:"reference"`
	MerchantRef string `json:"merchant_ref"`
	Method      string `json:"payment_method"`
	Status      string `json:"status"` // "PAID", "EXPIRED", "FAILED"
	TotalAmount int    `json:"total_amount"`
	PaidAt      int64  `json:"paid_at"`
	Note        string `json:"note"`
}

// GenerateRequestSignature: HMAC-SHA256(merchant_code + merchant_ref + amount, private_key)
func (c *Client) GenerateRequestSignature(merchantRef string, amount int) string {
	payload := fmt.Sprintf("%s%s%d", c.cfg.MerchantCode, merchantRef, amount)
	h := hmac.New(sha256.New, []byte(c.cfg.PrivateKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateCallbackSignature: checks if X-Callback-Signature matches HMAC-SHA256(raw_body, private_key)
func (c *Client) ValidateCallbackSignature(rawBody []byte, callbackSignature string) bool {
	if strings.TrimSpace(c.cfg.PrivateKey) == "" {
		return false
	}
	h := hmac.New(sha256.New, []byte(c.cfg.PrivateKey))
	h.Write(rawBody)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(callbackSignature))
}

// CreateClosedQRIS: creates a new QRIS transaction on Tripay
func (c *Client) CreateClosedQRIS(ctx context.Context, merchantRef string, amount int, custName, custPhone string) (*TransactionResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("tripay credentials not configured")
	}

	sig := c.GenerateRequestSignature(merchantRef, amount)
	email := fmt.Sprintf("%s@pasien.klinik.local", custPhone)
	if custPhone == "" {
		email = "pasien@klinik.local"
	}

	reqBody := CreateTransactionRequest{
		Method:        "QRIS2", // Tripay dynamic QRIS
		MerchantRef:   merchantRef,
		Amount:        amount,
		CustomerName:  custName,
		CustomerEmail: email,
		CustomerPhone: custPhone,
		OrderItems: []OrderItem{
			{
				Name:     "Registrasi & Antrian Pasien",
				Price:    amount,
				Quantity: 1,
			},
		},
		CallbackURL: c.cfg.CallbackURL,
		ReturnURL:   c.cfg.ReturnURL,
		ExpiredTime: time.Now().Add(2 * time.Hour).Unix(),
		Signature:   sig,
	}

	bodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/transaction/create", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tripay http call failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tripay response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("decode tripay response (%s): %w", string(respBytes), err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("tripay API error: %s", apiResp.Message)
	}

	return &apiResp.Data, nil
}
