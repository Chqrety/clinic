package handler

import (
	"encoding/json"
	"net/http"
	"time"

	treasuryApp "clinic-api/internal/application/treasury"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TreasuryHandler struct {
	svc *treasuryApp.Service
}

func NewTreasuryHandler(svc *treasuryApp.Service) *TreasuryHandler {
	return &TreasuryHandler{svc: svc}
}

func (h *TreasuryHandler) HandleOpenShift(w http.ResponseWriter, r *http.Request) {
	var req dto.OpenShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	shift, err := h.svc.OpenShift(r.Context(), treasuryApp.OpenShiftCommand{
		BranchID:        req.BranchID,
		CostCenterID:    req.CostCenterID,
		CashierUserID:   req.CashierUserID,
		CashierName:     req.CashierName,
		CashAccountCode: req.CashAccountCode,
		OpeningCash:     req.OpeningCash,
		Notes:           req.Notes,
	})
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Cashier shift opened successfully", shift)
}

func (h *TreasuryHandler) HandleGetActiveShift(w http.ResponseWriter, r *http.Request) {
	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok || branchID == uuid.Nil {
		response.BadRequest(w, "Valid X-Branch-ID header is required", nil)
		return
	}

	cashAccount := r.URL.Query().Get("cash_account")
	if cashAccount == "" {
		cashAccount = "1-1111" // Default Kidinest Kasir
	}

	shift, err := h.svc.GetActiveShift(r.Context(), branchID, cashAccount)
	if err != nil {
		response.NotFound(w, "No active shift found: "+err.Error())
		return
	}

	response.OK(w, "Active cashier shift retrieved successfully", shift, nil)
}

func (h *TreasuryHandler) HandleCloseShift(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid shift ID", nil)
		return
	}

	var req dto.CloseShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	var denoms []treasuryApp.DenominationInput
	for _, d := range req.Denominations {
		denoms = append(denoms, treasuryApp.DenominationInput{
			DenominationValue: d.DenominationValue,
			Quantity:          d.Quantity,
		})
	}

	shift, journal, err := h.svc.CloseShift(r.Context(), treasuryApp.CloseShiftCommand{
		ShiftID:       id,
		SystemCashIn:  req.SystemCashIn,
		SystemCashOut: req.SystemCashOut,
		Denominations: denoms,
		Notes:         req.Notes,
	})
	if err != nil {
		response.BadRequest(w, "Failed to close shift: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"shift":   shift,
		"journal": journal,
	}
	response.OK(w, "Cashier shift closed and physical cash reconciled successfully", data, nil)
}

func (h *TreasuryHandler) HandleDepositCash(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid shift ID", nil)
		return
	}

	var req dto.DepositCashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	shift, journal, err := h.svc.DepositCashToBank(r.Context(), treasuryApp.DepositCommand{
		ShiftID:        id,
		DepositAmount:  req.DepositAmount,
		BankAccCode:    req.BankAccCode,
		DepositedBy:    req.DepositedBy,
		DepositSlipRef: req.DepositSlipRef,
	})
	if err != nil {
		response.BadRequest(w, "Failed to deposit shift cash: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"shift":   shift,
		"journal": journal,
	}
	response.OK(w, "Shift cash deposited to bank successfully", data, nil)
}

func (h *TreasuryHandler) HandleGetShift(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid shift ID", nil)
		return
	}

	shift, err := h.svc.GetShiftByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Cashier shift not found")
		return
	}

	response.OK(w, "Cashier shift retrieved successfully", shift, nil)
}

func (h *TreasuryHandler) HandleListShifts(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	var startPtr, endPtr *time.Time
	if s := r.URL.Query().Get("start_date"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			startPtr = &t
		}
	}
	if e := r.URL.Query().Get("end_date"); e != "" {
		if t, err := time.Parse("2006-01-02", e); err == nil {
			endPtr = &t
		}
	}

	shifts, err := h.svc.ListShifts(r.Context(), branchID, startPtr, endPtr)
	if err != nil {
		response.InternalServerError(w, "Failed to list shifts: "+err.Error())
		return
	}

	response.OK(w, "Cashier shifts retrieved successfully", shifts, nil)
}

