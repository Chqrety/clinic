package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"clinic-api/internal/pharmacy/domain"
	"clinic-api/internal/pharmacy/usecase"
	"clinic-api/internal/shared/middleware"
	"clinic-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type PharmacyHandler struct {
	usecase *usecase.PharmacyUsecase
}

func NewPharmacyHandler(u *usecase.PharmacyUsecase) *PharmacyHandler {
	return &PharmacyHandler{usecase: u}
}

func getStaffName(r *http.Request) string {
	if nameVal := r.Context().Value(middleware.UserNameKey); nameVal != nil {
		if name, ok := nameVal.(string); ok && name != "" {
			return name
		}
	}
	return "Apoteker Bertugas"
}

func (h *PharmacyHandler) HandleListPrescriptions(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	list, err := h.usecase.ListIncomingPrescriptions(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat resep masuk: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, list)
}

func (h *PharmacyHandler) HandleGetPrescription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID resep tidak valid")
		return
	}

	rx, err := h.usecase.GetPrescription(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPrescriptionNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal memuat resep: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, rx)
}

func (h *PharmacyHandler) HandleStartPreparing(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID resep tidak valid")
		return
	}

	if err := h.usecase.StartPreparing(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"id": id, "status": "PREPARING"})
}

func (h *PharmacyHandler) HandleDispense(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID resep tidak valid")
		return
	}

	if err := h.usecase.Dispense(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, domain.ErrPrescriptionNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrAlreadyDispensed), errors.Is(err, domain.ErrAlreadyCancelled):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memproses penyerahan obat: "+err.Error())
		}
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"id": id, "status": "DISPENSED"})
}

func (h *PharmacyHandler) HandleCancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID resep tidak valid")
		return
	}

	if err := h.usecase.Cancel(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"id": id, "status": "CANCELLED"})
}

func (h *PharmacyHandler) HandleListDrugs(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	drugs, err := h.usecase.ListDrugs(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat daftar obat: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, drugs)
}

func (h *PharmacyHandler) HandleCreateDrug(w http.ResponseWriter, r *http.Request) {
	var in domain.CreateDrugInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload obat tidak valid")
		return
	}

	drug, err := h.usecase.CreateDrug(r.Context(), &in)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, drug)
}

func (h *PharmacyHandler) HandleRestock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	drugID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID obat tidak valid")
		return
	}

	var in domain.RestockInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload restock tidak valid")
		return
	}
	in.DrugID = drugID

	if err := h.usecase.Restock(r.Context(), &in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]any{"drug_id": drugID, "message": "stok batch berhasil ditambahkan"})
}

func (h *PharmacyHandler) HandleExportOpnameCSV(w http.ResponseWriter, r *http.Request) {
	unit := r.URL.Query().Get("unit")
	csvStr, err := h.usecase.ExportStockOpnameCSV(r.Context(), unit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal ekspor lembar stok opnam: "+err.Error())
		return
	}

	filename := fmt.Sprintf("stok-opnam-%s.csv", time.Now().Format("20060102-150405"))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(csvStr))
}

func (h *PharmacyHandler) HandleImportOpnameCSV(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		response.Error(w, http.StatusBadRequest, "konten CSV kosong atau tidak terbaca")
		return
	}

	adjustedBy := getStaffName(r)
	count, err := h.usecase.ImportStockOpnameCSV(r.Context(), string(bodyBytes), adjustedBy)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("Berhasil melakukan rekonsiliasi %d item stok opnam", count),
		"count":   count,
	})
}

func (h *PharmacyHandler) HandleListOpnameLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.usecase.ListOpnameLogs(r.Context(), 50)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat log audit opnam: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, logs)
}
