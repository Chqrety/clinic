package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, Envelope{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Envelope{
		Success: false,
		Data:    nil,
		Error:   message,
	})
}
