package handler

import (
	"encoding/json"
	"net/http"
	"time"

	payableApp "clinic-api/internal/application/payable"
	"clinic-api/internal/domain/payable"
	"clinic-api/internal/domain/shared"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PayableHandler struct {
	payableSvc *payableApp.Service
}

func NewPayableHandler(payableSvc *payableApp.Service) *PayableHandler {
	return &PayableHandler{payableSvc: payableSvc}
}

func (h *PayableHandler) HandleCreatePayable(w http.ResponseWriter, r *http.Request) {
	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok || branchID == uuid.Nil {
		response.BadRequest(w, "Valid X-Branch-ID header is required", nil)
		return
	}

	var req dto.CreatePayableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	cmd := payableApp.CreatePayableCommand{
		BranchScoped:     shared.BranchScoped{BranchID: branchID},
		VendorType:       req.VendorType,
		VendorID:         req.VendorID,
		VendorName:       req.VendorName,
		InvoiceReference: req.InvoiceReference,
		TotalAmount:      req.TotalAmount,
		IssueDate:        req.IssueDate,
		DueDate:          req.DueDate,
		AccountCode:      req.AccountCode,
		Notes:            req.Notes,
	}

	p, err := h.payableSvc.CreatePayable(r.Context(), cmd)
	if err != nil {
		response.InternalServerError(w, "Failed to create payable: "+err.Error())
		return
	}

	response.Created(w, "Payable successfully recorded", p)
}

func (h *PayableHandler) HandlePayPayable(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	payableID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid payable UUID", nil)
		return
	}

	var req dto.PayPayableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	cmd := payableApp.PayPayableCommand{
		PayableID:           payableID,
		Amount:              req.Amount,
		PaymentDate:         req.PaymentDate,
		CashBankAccountCode: req.CashBankAccountCode,
		ReferenceNo:         req.ReferenceNo,
		Notes:               req.Notes,
	}

	payRecord, err := h.payableSvc.PayPayable(r.Context(), cmd)
	if err != nil {
		response.InternalServerError(w, "Failed to process payable payment: "+err.Error())
		return
	}

	response.OK(w, "Payable payment recorded and Cash Disbursement Journal posted", payRecord, nil)
}

func (h *PayableHandler) HandleGetPayableByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	payableID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid payable UUID", nil)
		return
	}

	p, err := h.payableSvc.GetPayableByID(r.Context(), payableID)
	if err != nil {
		response.NotFound(w, "Payable record not found: "+err.Error())
		return
	}

	response.OK(w, "Payable details retrieved", p, nil)
}

func (h *PayableHandler) HandleListPayables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"due_date":           "due_date",
			"issue_date":         "issue_date",
			"created_at":         "created_at",
			"payable_number":     "payable_number",
			"vendor_name":        "vendor_name",
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

	filter := payable.Filter{
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
		status := payable.Status(s)
		filter.Status = &status
	}
	if vt := q.Get("vendor_type"); vt != "" {
		vendorType := payable.VendorType(vt)
		filter.VendorType = &vendorType
	}

	list, total, err := h.payableSvc.ListPayables(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list payables: "+err.Error())
		return
	}

	response.OK(w, "Payables retrieved successfully", list, pg.ToMeta(total))
}

func (h *PayableHandler) HandleGetAgingReport(w http.ResponseWriter, r *http.Request) {
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

	aging, err := h.payableSvc.GetAgingReport(r.Context(), branchID, asOf)
	if err != nil {
		response.InternalServerError(w, "Failed to generate payable aging report: "+err.Error())
		return
	}

	response.OK(w, "Payable aging report calculated successfully", aging, nil)
}
