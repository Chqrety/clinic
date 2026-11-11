package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterStaffRoutes(
	r chi.Router,
	h *StaffHandler,
	authMw func(http.Handler) http.Handler,
	ownerMw func(http.Handler) http.Handler,
) {
	r.Route("/api/admin/staff", func(r chi.Router) {
		r.Use(authMw)
		r.Use(ownerMw)

		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Put("/{id}/reset-password", h.ResetPassword)
	})
}
