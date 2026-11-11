package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	app "clinic-api/internal/application/billing"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	domain "clinic-api/internal/domain/billing"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BillingHandler struct {
	service *app.Service
}

func NewBillingHandler(service *app.Service) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/invoices", func(inv chi.Router) {
		inv.Get("/", h.HandleListInvoices)
		inv.Post("/", h.HandleCreateInvoice)
		inv.Get("/{id}", h.HandleGetInvoiceByID)
		inv.Post("/{id}/pay", h.HandlePayInvoice)
		inv.Get("/queue/{queueId}", h.HandleGetByQueueID)
	})
}

func (h *BillingHandler) HandleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New() // Fallback if single branch
	}

	items := make([]app.CreateItemCommand, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, app.CreateItemCommand{
			ItemType:    it.ItemType,
			ItemID:      it.ItemID,
			Description: it.Description,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
		})
	}

	cmd := app.CreateInvoiceCommand{
		BranchID:        branchID,
		PatientID:       req.PatientID,
		PatientName:     req.PatientName,
		PatientNIK:      req.PatientNIK,
		QueueID:         req.QueueID,
		InvoiceType:     req.InvoiceType,
		RegistrationFee: req.RegistrationFee,
		ConsultationFee: req.ConsultationFee,
		MedicationFee:   req.MedicationFee,
		Items:           items,
	}

	inv, err := h.service.CreateInvoice(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Invoice created successfully", toInvoiceResponse(inv))
}

func (h *BillingHandler) HandlePayInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid invoice UUID", nil)
		return
	}

	var req dto.PayInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "CASH"
	}

	cmd := app.PayInvoiceCommand{
		InvoiceID:         invoiceID,
		Amount:            req.Amount,
		PaymentMethod:     req.PaymentMethod,
		UnitCode:          req.UnitCode,
		ReferenceNo:       req.ReferenceNo,
		Notes:             req.Notes,
		DoctorName:        req.DoctorName,
		DoctorShareAmount: req.DoctorShareAmount,
		COGSAmount:        req.COGSAmount,
	}

	inv, journalEntry, err := h.service.PayInvoice(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	type PayResponseData struct {
		Invoice       dto.InvoiceResponse `json:"invoice"`
		JournalNumber string              `json:"journal_number,omitempty"`
		JournalID     *uuid.UUID          `json:"journal_id,omitempty"`
	}

	var jNum string
	var jID *uuid.UUID
	if journalEntry != nil {
		jNum = journalEntry.JournalNumber
		jID = &journalEntry.ID
	}

	response.OK(w, "Payment processed and accounting journal posted successfully", PayResponseData{
		Invoice:       toInvoiceResponse(inv),
		JournalNumber: jNum,
		JournalID:     jID,
	}, nil)
}

func (h *BillingHandler) HandleGetInvoiceByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	invoiceID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid invoice UUID", nil)
		return
	}

	inv, err := h.service.GetInvoiceByID(r.Context(), invoiceID)
	if err != nil {
		response.NotFound(w, "Invoice not found")
		return
	}

	response.OK(w, "Invoice retrieved successfully", toInvoiceResponse(inv), nil)
}

func (h *BillingHandler) HandleGetByQueueID(w http.ResponseWriter, r *http.Request) {
	qStr := chi.URLParam(r, "queueId")
	qID, err := strconv.ParseInt(qStr, 10, 64)
	if err != nil {
		response.BadRequest(w, "Invalid queue ID integer", nil)
		return
	}

	invoices, err := h.service.GetInvoicesByQueue(r.Context(), qID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.InvoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		resList = append(resList, toInvoiceResponse(&inv))
	}

	response.OK(w, "Invoices retrieved successfully", resList, nil)
}

func (h *BillingHandler) HandleListInvoices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"created_at":     "created_at",
			"invoice_number": "invoice_number",
			"patient_name":   "patient_name",
			"total_amount":   "total_amount",
			"paid_amount":    "paid_amount",
			"payment_status": "payment_status",
		},
		DefaultSortBy:    "created_at",
		DefaultSortOrder: "DESC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	var filter domain.InvoiceFilter
	if pg.BranchID != nil {
		filter.BranchID = pg.BranchID
	} else if branchID, ok := middleware.GetBranchID(ctx); ok {
		filter.BranchID = &branchID
	}

	if status := r.URL.Query().Get("status"); status != "" {
		st := domain.PaymentStatus(status)
		filter.PaymentStatus = &st
	}
	if invType := r.URL.Query().Get("type"); invType != "" {
		t := domain.InvoiceType(invType)
		filter.InvoiceType = &t
	}

	filter.StartDate = pg.StartDate
	filter.EndDate = pg.EndDate
	filter.Search = pg.Search
	filter.SortBy = pg.SortBy
	filter.SortOrder = pg.SortOrder
	filter.Limit = pg.Limit
	filter.Offset = pg.Offset

	invoices, totalCount, err := h.service.ListInvoices(ctx, filter)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.InvoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		resList = append(resList, toInvoiceResponse(&inv))
	}

	response.OK(w, "Invoices retrieved successfully", resList, pg.ToMeta(totalCount))
}

func toInvoiceResponse(inv *domain.Invoice) dto.InvoiceResponse {
	itemResponses := make([]dto.InvoiceItemResponse, 0, len(inv.Items))
	for _, it := range inv.Items {
		itemResponses = append(itemResponses, dto.InvoiceItemResponse{
			ID:          it.ID,
			ItemType:    it.ItemType,
			Description: it.Description,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			TotalPrice:  it.TotalPrice,
		})
	}

	paymentResponses := make([]dto.PaymentRecordResponse, 0, len(inv.Payments))
	for _, p := range inv.Payments {
		paymentResponses = append(paymentResponses, dto.PaymentRecordResponse{
			ID:            p.ID,
			PaymentNumber: p.PaymentNumber,
			Amount:        p.Amount,
			PaymentMethod: p.PaymentMethod,
			ReferenceNo:   p.ReferenceNo,
			Notes:         p.Notes,
			JournalID:     p.JournalID,
			PaidAt:        p.PaidAt,
		})
	}

	return dto.InvoiceResponse{
		ID:              inv.ID,
		BranchID:        inv.BranchID,
		InvoiceNumber:   inv.InvoiceNumber,
		PatientID:       inv.PatientID,
		PatientName:     inv.PatientName,
		PatientNIK:      inv.PatientNIK,
		QueueID:         inv.QueueID,
		InvoiceType:     inv.InvoiceType,
		RegistrationFee: inv.RegistrationFee,
		ConsultationFee: inv.ConsultationFee,
		MedicationFee:   inv.MedicationFee,
		TotalAmount:     inv.TotalAmount,
		PaidAmount:      inv.PaidAmount,
		PaymentStatus:   inv.PaymentStatus,
		PaymentMethod:   inv.PaymentMethod,
		JournalID:       inv.JournalID,
		Items:           itemResponses,
		Payments:        paymentResponses,
		CreatedAt:       inv.CreatedAt,
	}
}

