package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	budgetApp "clinic-api/internal/application/budget"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BudgetHandler struct {
	svc *budgetApp.Service
}

func NewBudgetHandler(svc *budgetApp.Service) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

func (h *BudgetHandler) HandleCreateBudget(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	var items []budgetApp.BudgetItemInput
	for _, it := range req.Items {
		items = append(items, budgetApp.BudgetItemInput{
			AccountCode:    it.AccountCode,
			PeriodMonth:    it.PeriodMonth,
			BudgetedAmount: it.BudgetedAmount,
			Notes:          it.Notes,
		})
	}

	b, err := h.svc.CreateBudget(r.Context(), budgetApp.CreateBudgetCommand{
		BranchID:     req.BranchID,
		CostCenterID: req.CostCenterID,
		PeriodYear:   req.PeriodYear,
		Name:         req.Name,
		Items:        items,
	})
	if err != nil {
		response.BadRequest(w, "Failed to create budget: "+err.Error(), nil)
		return
	}

	response.Created(w, "Budget created successfully", b)
}

func (h *BudgetHandler) HandleApproveBudget(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid budget ID", nil)
		return
	}

	b, err := h.svc.ApproveBudget(r.Context(), id)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Budget approved successfully", b, nil)
}

func (h *BudgetHandler) HandleGetVarianceReport(w http.ResponseWriter, r *http.Request) {
	costCenterStr := r.URL.Query().Get("cost_center_id")
	costCenterID, err := uuid.Parse(costCenterStr)
	if err != nil {
		response.BadRequest(w, "Query parameter 'cost_center_id' must be a valid UUID", nil)
		return
	}

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	if year == 0 || month == 0 {
		response.BadRequest(w, "Query parameters 'year' and 'month' are required", nil)
		return
	}

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	report, err := h.svc.GetBudgetVarianceReport(r.Context(), branchID, costCenterID, year, month)
	if err != nil {
		response.BadRequest(w, "Failed to calculate budget variance report: "+err.Error(), nil)
		return
	}

	response.OK(w, "Budget variance report retrieved successfully", report, nil)
}

func (h *BudgetHandler) HandleGetBudgetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid budget ID", nil)
		return
	}

	b, err := h.svc.GetBudgetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Budget not found")
		return
	}

	response.OK(w, "Budget retrieved successfully", b, nil)
}

func (h *BudgetHandler) HandleListBudgets(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	list, err := h.svc.ListBudgets(r.Context(), branchID, year)
	if err != nil {
		response.InternalServerError(w, "Failed to list budgets: "+err.Error())
		return
	}

	response.OK(w, "Budgets retrieved successfully", list, nil)
}

