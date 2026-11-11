package handler

import (
	"encoding/json"
	"net/http"
	"time"

	receivableApp "clinic-api/internal/application/receivable"
	"clinic-api/internal/domain/receivable"
	"clinic-api/internal/domain/shared"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReceivableHandler struct {
	receivableSvc *receivableApp.Service
}

func NewReceivableHandler(receivableSvc *receivableApp.Service) *ReceivableHandler {
	return &ReceivableHandler{receivableSvc: receivableSvc}
}

func (h *ReceivableHandler) HandleCreateReceivable(w http.ResponseWriter, r *http.Request) {
	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok || branchID == uuid.Nil {
		response.BadRequest(w, "Valid X-Branch-ID header is required", nil)
		return
	}

	var req dto.CreateReceivableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	cmd := receivableApp.CreateReceivableCommand{
		BranchScoped:     shared.BranchScoped{BranchID: branchID},
		CustomerType:     req.CustomerType,
		CustomerID:       req.CustomerID,
		CustomerName:     req.CustomerName,
		InvoiceReference: req.InvoiceReference,
		TotalAmount:      req.TotalAmount,
		IssueDate:        req.IssueDate,
		DueDate:          req.DueDate,
		AccountCode:      req.AccountCode,
		Notes:            req.Notes,
	}

	rec, err := h.receivableSvc.CreateReceivable(r.Context(), cmd)
	if err != nil {
		response.InternalServerError(w, "Failed to create receivable: "+err.Error())
		return
	}

	response.Created(w, "Receivable successfully recorded", rec)
}

func (h *ReceivableHandler) HandlePayReceivable(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	recID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid receivable UUID", nil)
		return
	}

	var req dto.PayReceivableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	cmd := receivableApp.PayReceivableCommand{
		ReceivableID:        recID,
		Amount:              req.Amount,
		PaymentDate:         req.PaymentDate,
		CashBankAccountCode: req.CashBankAccountCode,
		ReferenceNo:         req.ReferenceNo,
		Notes:               req.Notes,
	}

	payRecord, err := h.receivableSvc.PayReceivable(r.Context(), cmd)
	if err != nil {
		response.InternalServerError(w, "Failed to process receivable payment: "+err.Error())
		return
	}

	response.OK(w, "Receivable payment recorded and Cash Receipt Journal posted", payRecord, nil)
}

func (h *ReceivableHandler) HandleGetReceivableByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	recID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid receivable UUID", nil)
		return
	}

	rec, err := h.receivableSvc.GetReceivableByID(r.Context(), recID)
	if err != nil {
		response.NotFound(w, "Receivable record not found: "+err.Error())
		return
	}

	response.OK(w, "Receivable details retrieved", rec, nil)
}

func (h *ReceivableHandler) HandleListReceivables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"due_date":           "due_date",
			"issue_date":         "issue_date",
			"created_at":         "created_at",
			"receivable_number":  "receivable_number",
			"customer_name":      "customer_name",
			"total_amount":       "total_amount",
			"outstanding_amount": "outstanding_amount",
			"status":             "status",
		},
		DefaultSortBy:    "due_date",
		DefaultSortOrder: "ASC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := receivable.Filter{
		StartDate: pg.StartDate,
		EndDate:   pg.EndDate,
		Search:    pg.Search,
		SortBy:    pg.SortBy,
		SortOrder: pg.SortOrder,
		Limit:     pg.Limit,
		Offset:    pg.Offset,
	}

	if pg.BranchID != nil {
		filter.BranchID = pg.BranchID
	} else if bID, ok := middleware.GetBranchID(ctx); ok {
		filter.BranchID = &bID
	}

	q := r.URL.Query()
	if s := q.Get("status"); s != "" {
		status := receivable.Status(s)
		filter.Status = &status
	}
	if ct := q.Get("customer_type"); ct != "" {
		customerType := receivable.CustomerType(ct)
		filter.CustomerType = &customerType
	}

	list, total, err := h.receivableSvc.ListReceivables(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list receivables: "+err.Error())
		return
	}

	response.OK(w, "Receivables retrieved successfully", list, pg.ToMeta(total))
}

func (h *ReceivableHandler) HandleGetAgingReport(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok {
		branchID = &bID
	}

	asOf := time.Now()
	if d := r.URL.Query().Get("as_of_date"); d != "" {
		if parsed, err := time.Parse("2006-01-02", d); err == nil {
			asOf = parsed
		}
	}

	aging, err := h.receivableSvc.GetAgingReport(r.Context(), branchID, asOf)
	if err != nil {
		response.InternalServerError(w, "Failed to generate receivable aging report: "+err.Error())
		return
	}

	response.OK(w, "Receivable aging report calculated successfully", aging, nil)
}
