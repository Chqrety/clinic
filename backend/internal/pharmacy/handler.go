package pharmacy

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clinic-api/internal/domain"
	"clinic-api/pkg/pagination"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, s *Service) {
	r.Get("/drugs", s.HandleGetDrugs)
	r.Post("/drugs", s.HandleCreateDrug)
	r.Get("/drugs/{id}", s.HandleGetDrugDetail)
	r.Post("/drugs/{id}/restock", s.HandleRestock)
}

func (s *Service) HandleGetDrugs(w http.ResponseWriter, r *http.Request) {
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"name":       "d.name",
			"category":   "d.category",
			"sell_price": "d.sell_price",
			"created_at": "d.created_at",
		},
		DefaultSortBy:    "name",
		DefaultSortOrder: "ASC",
	})
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid pagination parameters: "+err.Error())
		return
	}

	filter := DrugFilter{
		Search:    pg.Search,
		Category:  r.URL.Query().Get("category"),
		SortBy:    pg.SortBy,
		SortOrder: pg.SortOrder,
		Limit:     pg.Limit,
		Offset:    pg.Offset,
	}

	drugs, total, err := s.GetDrugsPaginated(filter)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	meta := pg.ToMeta(total)
	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success:    true,
		Data:       drugs,
		Meta:       meta,
		Pagination: meta,
	})
}

func (s *Service) HandleGetDrugDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid drug ID")
		return
	}

	drug, err := s.GetDrugDetail(id)
	if err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, domain.ApiResponse{Success: true, Data: drug})
}

func (s *Service) HandleCreateDrug(w http.ResponseWriter, r *http.Request) {
	var req CreateDrugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	created, err := s.CreateDrug(req)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, domain.ApiResponse{
		Success: true,
		Message: "Jenis obat baru berhasil ditambahkan",
		Data:    created,
	})
}

func (s *Service) HandleRestock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid drug ID")
		return
	}

	var req RestockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.Restock(id, req); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, domain.ApiResponse{
		Success: true,
		Message: "Restock obat berhasil dicatat ke gudang apotek",
	})
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, domain.ApiResponse{Success: false, Error: msg})
}

