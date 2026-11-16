package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	docApp "clinic-api/internal/application/doctorshare"
	"clinic-api/internal/domain/doctorshare"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DoctorShareHandler struct {
	svc *docApp.Service
}

func NewDoctorShareHandler(svc *docApp.Service) *DoctorShareHandler {
	return &DoctorShareHandler{svc: svc}
}

func (h *DoctorShareHandler) HandleCreateRule(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDoctorShareRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	rule, err := h.svc.CreateRule(r.Context(), docApp.CreateRuleCommand{
		BranchID:          req.BranchID,
		CostCenterID:      req.CostCenterID,
		DoctorID:          req.DoctorID,
		DoctorName:        req.DoctorName,
		ServiceCategory:   req.ServiceCategory,
		SharingPercentage: req.SharingPercentage,
		FixedFee:          req.FixedFee,
		IsActive:          req.IsActive,
	})
	if err != nil {
		response.InternalServerError(w, "Failed to create doctor share rule: "+err.Error())
		return
	}

	response.Created(w, "Doctor share rule created successfully", rule)
}

func (h *DoctorShareHandler) HandleListRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"doctor_name":        "doctor_name",
			"service_category":   "service_category",
			"sharing_percentage": "sharing_percentage",
			"fixed_fee":          "fixed_fee",
			"is_active":          "is_active",
			"created_at":         "created_at",
		},
		DefaultSortBy:    "doctor_name",
		DefaultSortOrder: "ASC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := doctorshare.RuleFilter{
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

	rules, total, err := h.svc.ListRules(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list doctor share rules: "+err.Error())
		return
	}

	response.OK(w, "Doctor share rules retrieved successfully", rules, pg.ToMeta(total))
}

func (h *DoctorShareHandler) HandleGenerateSettlement(w http.ResponseWriter, r *http.Request) {
	var req dto.GenerateDoctorSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	var items []docApp.FeeItemInput
	for _, it := range req.Items {
		items = append(items, docApp.FeeItemInput{
			InvoiceID:          it.InvoiceID,
			ServiceDate:        it.ServiceDate,
			PatientName:        it.PatientName,
			ServiceDescription: it.ServiceDescription,
			PatientFee:         it.PatientFee,
			DoctorShareAmount:  it.DoctorShareAmount,
		})
	}

	settlement, err := h.svc.GenerateSettlement(r.Context(), docApp.GenerateSettlementCommand{
		BranchID:    req.BranchID,
		DoctorID:    req.DoctorID,
		DoctorName:  req.DoctorName,
		DoctorNPWP:  req.DoctorNPWP,
		TaxpayerNIK: req.TaxpayerNIK,
		PeriodYear:  req.PeriodYear,
		PeriodMonth: req.PeriodMonth,
		Items:       items,
	})
	if err != nil {
		response.InternalServerError(w, "Failed to generate settlement: "+err.Error())
		return
	}

	response.Created(w, "Doctor fee settlement generated successfully", settlement)
}

func (h *DoctorShareHandler) HandleApproveSettlement(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid settlement ID", nil)
		return
	}

	settlement, err := h.svc.ApproveSettlement(r.Context(), id)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Doctor fee settlement approved successfully", settlement, nil)
}

func (h *DoctorShareHandler) HandlePaySettlement(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid settlement ID", nil)
		return
	}

	var req dto.PayDoctorSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	settlement, journal, err := h.svc.PaySettlement(r.Context(), docApp.PaySettlementCommand{
		SettlementID: id,
		BankAccCode:  req.BankAccCode,
		UnitCode:     req.UnitCode,
		PaidBy:       req.PaidBy,
	})
	if err != nil {
		response.BadRequest(w, "Failed to pay doctor fee settlement: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"settlement": settlement,
		"journal":    journal,
	}
	response.OK(w, "Doctor fee settlement paid and journal posted successfully", data, nil)
}

func (h *DoctorShareHandler) HandleGetSettlement(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid settlement ID", nil)
		return
	}

	settlement, err := h.svc.GetSettlementByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Doctor fee settlement not found")
		return
	}

	response.OK(w, "Doctor fee settlement retrieved successfully", settlement, nil)
}

func (h *DoctorShareHandler) HandleListSettlements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"doctor_name":  "doctor_name",
			"period_year":  "period_year",
			"period_month": "period_month",
			"gross_fee":    "gross_fee",
			"net_fee":      "net_fee",
			"status":       "status",
			"created_at":   "created_at",
		},
		DefaultSortBy:    "period_year",
		DefaultSortOrder: "DESC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := doctorshare.SettlementFilter{
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
	if s := q.Get("status"); s != "" {
		st := doctorshare.SettlementStatus(s)
		filter.Status = &st
	}
	if docIDStr := q.Get("doctor_id"); docIDStr != "" {
		if docID, err := uuid.Parse(docIDStr); err == nil {
			filter.DoctorID = &docID
		}
	}

	settlements, total, err := h.svc.ListSettlements(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list doctor settlements: "+err.Error())
		return
	}

	response.OK(w, "Doctor fee settlements retrieved successfully", settlements, pg.ToMeta(total))
}

