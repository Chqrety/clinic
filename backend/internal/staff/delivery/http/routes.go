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
	r.Group(func(sr chi.Router) {
		sr.Use(authMw)
		sr.Use(ownerMw)

		sr.Get("/api/admin/staff", h.List)
		sr.Post("/api/admin/staff", h.Create)
		sr.Route("/api/admin/staff", func(sub chi.Router) {
			sub.Get("/", h.List)
			sub.Post("/", h.Create)
			sub.Get("/{id}", h.GetByID)
			sub.Put("/{id}", h.Update)
			sub.Put("/{id}/reset-password", h.ResetPassword)
		})
	})
}
