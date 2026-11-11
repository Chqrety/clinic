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
	"clinic-api/pkg/pagination"
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
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"created_at":    "created_at",
			"taxpayer_name": "taxpayer_name",
			"gross_amount":  "gross_amount",
			"tax_amount":    "tax_amount",
			"period_year":   "period_year",
			"period_month":  "period_month",
			"bupot_number":  "bupot_number",
		},
		DefaultSortBy:    "created_at",
		DefaultSortOrder: "DESC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := tax.TaxFilter{
		Search:    pg.Search,
		SortBy:    pg.SortBy,
		SortOrder: pg.SortOrder,
		Limit:     pg.Limit,
		Offset:    pg.Offset,
	}

	if pg.BranchID != nil {
		filter.BranchID = pg.BranchID
	} else if bID, ok := middleware.GetBranchID(ctx); ok && bID != uuid.Nil {
		filter.BranchID = &bID
	}

	q := r.URL.Query()
	if y, err := strconv.Atoi(q.Get("year")); err == nil && y > 0 {
		filter.Year = y
	}
	if m, err := strconv.Atoi(q.Get("month")); err == nil && m > 0 {
		filter.Month = m
	}
	if tStr := q.Get("tax_type"); tStr != "" {
		tt := tax.TaxType(tStr)
		filter.TaxType = &tt
	}
	if repStr := q.Get("is_reported"); repStr != "" {
		rep := repStr == "true" || repStr == "1"
		filter.IsReported = &rep
	}

	withholdings, total, err := h.svc.ListWithholdings(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list tax withholdings: "+err.Error())
		return
	}

	response.OK(w, "Tax withholdings retrieved successfully", withholdings, pg.ToMeta(total))
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

