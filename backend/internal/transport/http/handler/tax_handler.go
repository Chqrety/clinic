package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	taxApp "clinic-api/internal/application/tax"
	"clinic-api/internal/domain/tax"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaxHandler struct {
	svc *taxApp.Service
}

func NewTaxHandler(svc *taxApp.Service) *TaxHandler {
	return &TaxHandler{svc: svc}
}

func (h *TaxHandler) HandleCreateWithholding(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWithholdingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	tw, err := h.svc.CreateWithholding(r.Context(), taxApp.CreateWithholdingCommand{
		BranchID:     req.BranchID,
		TaxType:      req.TaxType,
		TaxpayerNIK:  req.TaxpayerNIK,
		TaxpayerNPWP: req.TaxpayerNPWP,
		TaxpayerName: req.TaxpayerName,
		ObjectCode:   req.ObjectCode,
		GrossAmount:  req.GrossAmount,
		DPPAmount:    req.DPPAmount,
		TaxRate:      req.TaxRate,
		TaxAmount:    req.TaxAmount,
		PeriodYear:   req.PeriodYear,
		PeriodMonth:  req.PeriodMonth,
		JournalID:    req.JournalID,
	})
	if err != nil {
		response.InternalServerError(w, "Failed to create tax withholding: "+err.Error())
		return
	}

	response.Created(w, "Tax withholding created successfully", tw)
}

func (h *TaxHandler) HandleListWithholdings(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	var taxTypePtr *tax.TaxType
	if tStr := r.URL.Query().Get("tax_type"); tStr != "" {
		tt := tax.TaxType(tStr)
		taxTypePtr = &tt
	}

	withholdings, err := h.svc.ListWithholdings(r.Context(), branchID, year, month, taxTypePtr)
	if err != nil {
		response.InternalServerError(w, "Failed to list tax withholdings: "+err.Error())
		return
	}

	response.OK(w, "Tax withholdings retrieved successfully", withholdings, nil)
}

func (h *TaxHandler) HandleGetTaxSummary(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	summary, err := h.svc.GetTaxSummary(r.Context(), branchID, year, month)
	if err != nil {
		response.InternalServerError(w, "Failed to get tax summary: "+err.Error())
		return
	}

	response.OK(w, "Tax summary retrieved successfully", summary, nil)
}

func (h *TaxHandler) HandleExportEBupot21CSV(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))

	csvBytes, err := h.svc.ExportEBupot21CSV(r.Context(), branchID, year, month)
	if err != nil {
		response.InternalServerError(w, "Failed to export e-Bupot CSV: "+err.Error())
		return
	}

	filename := fmt.Sprintf("ebupot_21_%d_%02d.csv", year, month)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(csvBytes)
}

func (h *TaxHandler) HandleMarkAsReported(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid withholding ID", nil)
		return
	}

	if err := h.svc.MarkAsReported(r.Context(), id); err != nil {
		response.InternalServerError(w, "Failed to mark as reported: "+err.Error())
		return
	}

	response.OK(w, "Tax withholding marked as reported to DJP", nil, nil)
}

