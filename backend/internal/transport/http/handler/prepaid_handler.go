package handler

import (
	"encoding/json"
	"net/http"

	prepaidApp "clinic-api/internal/application/prepaid"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PrepaidHandler struct {
	svc *prepaidApp.Service
}

func NewPrepaidHandler(svc *prepaidApp.Service) *PrepaidHandler {
	return &PrepaidHandler{svc: svc}
}

func (h *PrepaidHandler) HandleRegisterPrepaid(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterPrepaidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	p, journal, err := h.svc.RegisterPrepaid(r.Context(), prepaidApp.RegisterPrepaidCommand{
		BranchID:           req.BranchID,
		CostCenterID:       req.CostCenterID,
		Name:               req.Name,
		Category:           req.Category,
		TotalAmount:        req.TotalAmount,
		PrepaidAccountCode: req.PrepaidAccountCode,
		ExpenseAccountCode: req.ExpenseAccountCode,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		PaymentAccountCode: req.PaymentAccountCode,
		AutoPostPayment:    req.AutoPostPayment,
	})
	if err != nil {
		response.BadRequest(w, "Failed to register prepaid expense: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"prepaid": p,
		"journal": journal,
	}
	response.Created(w, "Prepaid expense registered successfully", data)
}

func (h *PrepaidHandler) HandleListActivePrepaid(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	list, err := h.svc.ListActivePrepaid(r.Context(), branchID)
	if err != nil {
		response.InternalServerError(w, "Failed to list prepaid expenses: "+err.Error())
		return
	}

	response.OK(w, "Active prepaid expenses retrieved successfully", list, nil)
}

func (h *PrepaidHandler) HandleGetPrepaidByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid prepaid ID", nil)
		return
	}

	p, err := h.svc.GetPrepaidByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Prepaid expense not found")
		return
	}

	response.OK(w, "Prepaid expense retrieved successfully", p, nil)
}

func (h *PrepaidHandler) HandleAmortize(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid prepaid ID", nil)
		return
	}

	var req dto.AmortizePrepaidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	amort, journal, err := h.svc.AmortizeMonthly(r.Context(), prepaidApp.AmortizeCommand{
		PrepaidID:   id,
		PeriodYear:  req.PeriodYear,
		PeriodMonth: req.PeriodMonth,
	})
	if err != nil {
		response.BadRequest(w, "Failed to amortize prepaid expense: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"amortization": amort,
		"journal":      journal,
	}
	response.OK(w, "Monthly amortization recorded and journal posted", data, nil)
}

func (h *PrepaidHandler) HandleRunBatchAmortization(w http.ResponseWriter, r *http.Request) {
	var req dto.RunBatchAmortizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	results, err := h.svc.RunMonthlyAmortizationBatch(r.Context(), branchID, req.PeriodYear, req.PeriodMonth)
	if err != nil {
		response.InternalServerError(w, "Failed to run batch amortization: "+err.Error())
		return
	}

	response.OK(w, "Batch amortization executed successfully", results, nil)
}

func (h *PrepaidHandler) HandleGetAmortizationHistory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid prepaid ID", nil)
		return
	}

	history, err := h.svc.GetAmortizationHistory(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, "Failed to get amortization history: "+err.Error())
		return
	}

	response.OK(w, "Amortization history retrieved successfully", history, nil)
}

