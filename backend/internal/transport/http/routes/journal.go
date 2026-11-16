package routes

import (
	"clinic-api/internal/transport/http/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterJournalRoutes registers routes for General & Special Journals.
func RegisterJournalRoutes(r chi.Router, h *handler.AccountingHandler) {
	if h == nil {
		return
	}
	r.Route("/journals", func(j chi.Router) {
		j.Get("/", h.HandleListJournals)
		j.Post("/", h.HandleCreateJournal)
		j.Get("/{id}", h.HandleGetJournalByID)
		j.Post("/{id}/post", h.HandlePostJournal)
		j.Post("/{id}/reverse", h.HandleReverseJournal)
	})
}

