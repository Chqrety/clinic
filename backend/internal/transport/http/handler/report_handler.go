package handler

import (
	"net/http"
	"time"

	app "clinic-api/internal/application/accounting"
	"clinic-api/internal/transport/http/middleware"
	domain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReportHandler struct {
	service *app.ReportService
}

func NewReportHandler(service *app.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) RegisterRoutes(r chi.Router) {
	r.Route("/reports", func(rep chi.Router) {
		rep.Get("/general-ledger", h.HandleGetGeneralLedger)
		rep.Get("/trial-balance", h.HandleGetTrialBalance)
		rep.Get("/income-statement", h.HandleGetIncomeStatement)
		rep.Get("/balance-sheet", h.HandleGetBalanceSheet)
		rep.Get("/cash-flow", h.HandleGetCashFlowStatement)
	})
}

func (h *ReportHandler) HandleGetGeneralLedger(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountCode := r.URL.Query().Get("account_code")
	accountIDStr := r.URL.Query().Get("account_id")

	var accountID *uuid.UUID
	if accountIDStr != "" {
		if id, err := uuid.Parse(accountIDStr); err == nil {
			accountID = &id
		}
	}

	if accountCode == "" && accountID == nil {
		response.BadRequest(w, "Query parameter 'account_code' or 'account_id' is required", nil)
		return
	}

	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	filter := domain.GeneralLedgerFilter{
		BranchID:    branchID,
		AccountID:   accountID,
		AccountCode: accountCode,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	ledger, err := h.service.GetGeneralLedger(ctx, filter)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "General Ledger retrieved successfully", ledger, nil)
}

func (h *ReportHandler) HandleGetTrialBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	tb, err := h.service.GetTrialBalance(ctx, branchID, startDate, endDate)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Trial Balance retrieved successfully", tb, nil)
}

func (h *ReportHandler) HandleGetIncomeStatement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	is, err := h.service.GetIncomeStatement(ctx, branchID, startDate, endDate)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Income Statement retrieved successfully", is, nil)
}

func (h *ReportHandler) HandleGetBalanceSheet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	asOfDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("as_of_date"))

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	bs, err := h.service.GetBalanceSheet(ctx, branchID, asOfDate)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Balance Sheet retrieved successfully", bs, nil)
}

func (h *ReportHandler) HandleGetCashFlowStatement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	cf, err := h.service.GetCashFlowStatement(ctx, branchID, startDate, endDate)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Cash Flow Statement retrieved successfully", cf, nil)
}


