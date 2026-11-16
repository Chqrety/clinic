package handler

import (
	"encoding/json"
	"net/http"

	accountingApp "clinic-api/internal/application/accounting"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

type ClosingHandler struct {
	svc *accountingApp.ClosingService
}

func NewClosingHandler(svc *accountingApp.ClosingService) *ClosingHandler {
	return &ClosingHandler{svc: svc}
}

func (h *ClosingHandler) HandleYearEndClosing(w http.ResponseWriter, r *http.Request) {
	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok || branchID == uuid.Nil {
		response.BadRequest(w, "Valid X-Branch-ID header is required", nil)
		return
	}

	var req dto.YearEndClosingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	result, err := h.svc.ExecuteYearEndClosing(r.Context(), branchID, req.FiscalYear, req.ClosedBy)
	if err != nil {
		response.BadRequest(w, "Year-end closing failed: "+err.Error(), nil)
		return
	}

	response.OK(w, "Fiscal year closed and closing journal posted to retained earnings successfully", result, nil)
}

