package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPublicQueueRoutes(r chi.Router, h *QueueHandler) {
	r.Route("/api/public/queue", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/verify-payment", h.VerifyPayment)
		r.Get("/ticket/{id}", h.GetTicket)
		r.Get("/monitor", h.GetMonitor)
		r.Get("/rooms/status", h.GetDoctorRoomsStatus)
		r.Get("/stream", h.StreamMonitor)
		r.Get("/patients/lookup", h.LookupPatient)
		r.Post("/reschedule", h.Reschedule)
		r.Post("/refund/request", h.RequestRefund)
	})
	r.Get("/api/public/patients/lookup", h.LookupPatient)
	r.Post("/api/webhooks/tripay", h.HandleTripayWebhook)
	r.Post("/api/webhooks/midtrans", h.HandleMidtransWebhook)
}

func RegisterStaffQueueRoutes(r chi.Router, h *QueueHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Group(func(gr chi.Router) {
		gr.Use(authMiddleware)
		gr.Route("/api/queues", func(qr chi.Router) {
			qr.Get("/reservations", h.ListReservations)
			qr.Get("/refunds", h.ListRefundRequests)
			qr.Post("/{id}/checkin", h.CheckIn)
			qr.Get("/rooms/status", h.GetDoctorRoomsStatus)
			qr.Post("/refund/process", h.ProcessRefund)
		})
	})
}
