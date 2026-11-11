package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterConsultationRoutes(
	r chi.Router,
	h *ConsultationHandler,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route("/api/consultations", func(cr chi.Router) {
		cr.Use(authMiddleware)

		// 1. Doctor Queue FIFO & Calling
		cr.Get("/queues", h.HandleListQueues)
		cr.Post("/queues/{id}/call", h.HandleCallPatient)
		cr.Post("/queues/{id}/start", h.HandleStartExam)

		// 2. Consultation & Prescriptions
		cr.Get("/queue/{queueId}", h.HandleGetByQueue)
		cr.Post("", h.HandleSubmitConsultation)

		// 3. Catalog Drugs with Live Stock for E-Prescription
		cr.Get("/drugs", h.HandleListDrugs)
	})
}
