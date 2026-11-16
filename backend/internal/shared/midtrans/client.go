package midtrans

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	ServerKey       string
	ClientKey       string
	SignatureSecret string
	IsProduction    bool
}

type Client struct {
	cfg        Config
	httpClient *http.Client
	baseURL    string
	snapURL    string
}

func NewClient(cfg Config) *Client {
	isProd := cfg.IsProduction || strings.HasPrefix(strings.TrimSpace(cfg.ServerKey), "Mid-server-")
	baseURL := "https://api.sandbox.midtrans.com"
	snapURL := "https://app.sandbox.midtrans.com"
	if isProd {
		baseURL = "https://api.midtrans.com"
		snapURL = "https://app.midtrans.com"
	}
	cfg.IsProduction = isProd
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		snapURL: snapURL,
	}
}

func (c *Client) IsConfigured() bool {
	return strings.TrimSpace(c.cfg.ServerKey) != ""
}

func (c *Client) IsProduction() bool {
	return c.cfg.IsProduction
}

func (c *Client) ClientKey() string {
	return c.cfg.ClientKey
}

func (c *Client) authHeader() string {
	raw := strings.TrimSpace(c.cfg.ServerKey) + ":"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

// QRIS Action represents Midtrans action link (e.g. generate-qr-code)
type Action struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

// ChargeResponse represents Midtrans Core API Charge response
type ChargeResponse struct {
	StatusCode        string   `json:"status_code"`
	StatusMessage     string   `json:"status_message"`
	TransactionID     string   `json:"transaction_id"`
	OrderID           string   `json:"order_id"`
	GrossAmount       string   `json:"gross_amount"`
	Currency          string   `json:"currency"`
	PaymentType       string   `json:"payment_type"`
	TransactionTime   string   `json:"transaction_time"`
	TransactionStatus string   `json:"transaction_status"`
	Actions           []Action `json:"actions"`
	QRString          string   `json:"qr_string"`
	QRURL             string   `json:"-"`
	SimulatorURL      string   `json:"-"`
}

// SnapResponse represents Midtrans Snap token response
type SnapResponse struct {
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

// NotificationPayload represents webhook payload from Midtrans
type NotificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

// IsPaymentSuccess checks whether the transaction status denotes a paid invoice
func (p *NotificationPayload) IsPaymentSuccess() bool {
	if p.TransactionStatus == "settlement" {
		return true
	}
	if p.TransactionStatus == "capture" && p.FraudStatus == "accept" {
		return true
	}
	return false
}

func (c *Client) chargeWithPayload(ctx context.Context, reqBody map[string]interface{}) (*ChargeResponse, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal midtrans payload: %w", err)
	}

	endpoint := c.baseURL + "/v2/charge"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create midtrans request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans http call failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read midtrans response: %w", err)
	}

	log.Printf("[MIDTRANS] POST %s -> HTTP %d: %s", endpoint, resp.StatusCode, string(respBytes))

	var chargeResp ChargeResponse
	if err := json.Unmarshal(respBytes, &chargeResp); err != nil {
		return nil, fmt.Errorf("decode midtrans response (%s): %w", string(respBytes), err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("midtrans HTTP %d error (code %s): %s", resp.StatusCode, chargeResp.StatusCode, chargeResp.StatusMessage)
	}

	if chargeResp.StatusCode != "" && chargeResp.StatusCode != "200" && chargeResp.StatusCode != "201" {
		return nil, fmt.Errorf("midtrans API error (code %s): %s", chargeResp.StatusCode, chargeResp.StatusMessage)
	}

	for _, act := range chargeResp.Actions {
		if act.Name == "generate-qr-code-v2" {
			// Prioritas tertinggi: Format resmi QRIS terbaru (merchants-app.../qr-code)
			chargeResp.QRURL = act.URL
		} else if act.Name == "generate-qr-code" && chargeResp.QRURL == "" {
			chargeResp.QRURL = act.URL
		}
		if act.Name == "deeplink-redirect" {
			chargeResp.SimulatorURL = act.URL
		}
	}

	if chargeResp.QRURL == "" && chargeResp.QRString != "" {
		chargeResp.QRURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", url.QueryEscape(chargeResp.QRString))
	}

	return &chargeResp, nil
}

// CreateQRIS initiates a dynamic direct QR-Code charge via Midtrans Core API (/v2/charge).
// Fungsi ini mencoba 3 format payload Core API Midtrans secara berurutan agar kompatibel
// baik dengan GoPay QRIS, Other QRIS, maupun e-Wallet GoPay direct:
// 1. payment_type: "qris" dengan acquirer "gopay"
// 2. payment_type: "qris" standar (Other QRIS)
// 3. payment_type: "gopay" direct
func (c *Client) CreateQRIS(ctx context.Context, orderID string, grossAmount int, patientName string, patientPhone string) (*ChargeResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("midtrans credentials not configured")
	}

	custDetails := map[string]interface{}{
		"first_name": patientName,
		"phone":      patientPhone,
	}
	txDetails := map[string]interface{}{
		"order_id":     orderID,
		"gross_amount": grossAmount,
	}

	// 1. Coba payment_type: "qris" dengan acquirer: "gopay" (Format resmi GoPay QRIS)
	resp, err := c.chargeWithPayload(ctx, map[string]interface{}{
		"payment_type":        "qris",
		"transaction_details": txDetails,
		"customer_details":    custDetails,
		"qris": map[string]interface{}{
			"acquirer": "gopay",
		},
	})
	if err == nil && resp != nil && (resp.QRURL != "" || resp.QRString != "") {
		return resp, nil
	}

	// 2. Coba payment_type: "qris" dengan acquirer: "airpay shopee" (Channel Other QRIS di Midtrans)
	log.Printf("[MIDTRANS] Core API QRIS gopay acquirer gagal (%v), mencoba acquirer: 'airpay shopee'...", err)
	resp, err = c.chargeWithPayload(ctx, map[string]interface{}{
		"payment_type":        "qris",
		"transaction_details": txDetails,
		"customer_details":    custDetails,
		"qris": map[string]interface{}{
			"acquirer": "airpay shopee",
		},
	})
	if err == nil && resp != nil && (resp.QRURL != "" || resp.QRString != "") {
		return resp, nil
	}

	return nil, err
}

// CreateSnapTransaction creates a Midtrans Snap transaction token specifically for GoPay QRIS
func (c *Client) CreateSnapTransaction(ctx context.Context, orderID string, grossAmount int, patientName string, patientPhone string) (*SnapResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("midtrans credentials not configured")
	}

	reqBody := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": grossAmount,
		},
		"customer_details": map[string]interface{}{
			"first_name": patientName,
			"phone":      patientPhone,
		},
		"enabled_payments": []string{"gopay"},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal snap payload: %w", err)
	}

	endpoint := c.snapURL + "/snap/v1/transactions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create snap request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("snap http call failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read snap response: %w", err)
	}

	log.Printf("[MIDTRANS] Snap %s -> HTTP %d: %s", endpoint, resp.StatusCode, string(respBytes))

	var snapResp SnapResponse
	if err := json.Unmarshal(respBytes, &snapResp); err != nil {
		return nil, fmt.Errorf("decode snap response (%s): %w", string(respBytes), err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("snap API error: %v", snapResp.ErrorMessages)
	}

	return &snapResp, nil
}

var qrURLRegex = regexp.MustCompile(`https?://[^\s"']+/qr-code[^\s"']*`)

// FetchSnapGoPayQR melakukan handshake internal Snap untuk memilih metode GoPay dan mengekstrak
// link gambar QRIS resmi (misal: https://merchants-app.sbx.midtrans.com/v4/qris/gopay/.../qr-code).
func (c *Client) FetchSnapGoPayQR(ctx context.Context, snapToken string) (string, error) {
	if snapToken == "" {
		return "", fmt.Errorf("empty snap token")
	}

	payload := map[string]interface{}{
		"payment_type":   "gopay",
		"payment_method": "gopay",
		"token":          snapToken,
	}
	bodyBytes, _ := json.Marshal(payload)

	candidates := []struct {
		url    string
		method string
	}{
		{url: c.snapURL + "/snap/v1/transactions/" + snapToken + "/pay", method: http.MethodPost},
		{url: c.snapURL + "/snap/v1/transactions/" + snapToken + "/charge", method: http.MethodPost},
		{url: c.snapURL + "/snap/v1/pay", method: http.MethodPost},
		{url: c.snapURL + "/snap/v4/transactions/" + snapToken + "/pay", method: http.MethodPost},
		{url: c.snapURL + "/snap/v1/transactions/" + snapToken, method: http.MethodGet},
	}

	for _, cand := range candidates {
		var req *http.Request
		var err error
		if cand.method == http.MethodGet {
			req, err = http.NewRequestWithContext(ctx, http.MethodGet, cand.url, nil)
		} else {
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, cand.url, bytes.NewReader(bodyBytes))
		}
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("X-Source", "snap-js")
		if c.cfg.ClientKey != "" {
			req.Header.Set("Client-Key", c.cfg.ClientKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			continue
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		respStr := string(respBytes)
		log.Printf("[MIDTRANS] Handshake Snap (%s) -> HTTP %d: %s", cand.url, resp.StatusCode, respStr)

		// 1. Coba regex mencari URL yang mengarah ke endpoint qr-code
		if match := qrURLRegex.FindString(respStr); match != "" {
			return match, nil
		}

		// 2. Coba parse JSON untuk mencari action generate-qr-code-v2 / generate-qr-code
		var data map[string]interface{}
		if err := json.Unmarshal(respBytes, &data); err == nil {
			if acts, ok := data["actions"].([]interface{}); ok {
				for _, a := range acts {
					if actMap, ok := a.(map[string]interface{}); ok {
						name, _ := actMap["name"].(string)
						u, _ := actMap["url"].(string)
						if (name == "generate-qr-code-v2" || name == "generate-qr-code") && u != "" {
							return u, nil
						}
					}
				}
			}
			for _, k := range []string{"qr_code_url", "qr_url", "qr_code", "url"} {
				if val, ok := data[k].(string); ok && strings.Contains(val, "/qr-code") {
					return val, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no qr-code URL found in snap handshake")
}

// TransactionStatusResponse represents the response from Midtrans Get Status API
type TransactionStatusResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	Currency          string `json:"currency"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

// CheckTransactionStatus queries Midtrans GET /v2/{order_id}/status to check live payment settlement status
func (c *Client) CheckTransactionStatus(ctx context.Context, orderID string) (*TransactionStatusResponse, error) {
	if !c.IsConfigured() || strings.TrimSpace(orderID) == "" {
		return nil, fmt.Errorf("midtrans not configured or empty orderID")
	}

	endpoint := fmt.Sprintf("%s/v2/%s/status", c.baseURL, url.PathEscape(orderID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create status request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("status request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read status response: %w", err)
	}

	var statusResp TransactionStatusResponse
	if err := json.Unmarshal(respBytes, &statusResp); err != nil {
		return nil, fmt.Errorf("decode status response (%s): %w", string(respBytes), err)
	}

	return &statusResp, nil
}

// ValidateNotificationSignature validates SHA512(order_id + status_code + gross_amount + SecretKey)
func (c *Client) ValidateNotificationSignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	if !c.IsConfigured() {
		return false
	}
	secret := strings.TrimSpace(c.cfg.SignatureSecret)
	if secret == "" {
		secret = strings.TrimSpace(c.cfg.ServerKey)
	}
	raw := orderID + statusCode + grossAmount + secret
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])
	return strings.EqualFold(expected, strings.TrimSpace(signatureKey))
}
