package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterPublicQueueRoutes(r chi.Router, h *QueueHandler) {
	r.Route("/api/public/queue", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/verify-payment", h.VerifyPayment)
		r.Get("/ticket/{id}", h.GetTicket)
		r.Get("/monitor", h.GetMonitor)
		r.Get("/stream", h.StreamMonitor)
	})
}
