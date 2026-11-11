package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h *AuthHandler, authMw func(http.Handler) http.Handler) {
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(authMw)
			r.Get("/me", h.Me)
		})
	})
}
