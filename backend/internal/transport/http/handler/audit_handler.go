package handler

import (
	"net/http"
	"strconv"

	auditApp "clinic-api/internal/application/audit"
	"clinic-api/internal/domain/audit"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/google/uuid"
)

type AuditHandler struct {
	svc *auditApp.Service
}

func NewAuditHandler(svc *auditApp.Service) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) HandleListAuditLogs(w http.ResponseWriter, r *http.Request) {
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	entityType := r.URL.Query().Get("entity_type")
	var actionPtr *audit.ActionType
	if aStr := r.URL.Query().Get("action"); aStr != "" {
		act := audit.ActionType(aStr)
		actionPtr = &act
	}

	logs, count, err := h.svc.ListLogs(r.Context(), branchID, actionPtr, entityType, limit, offset)
	if err != nil {
		response.InternalServerError(w, "Failed to list audit logs: "+err.Error())
		return
	}

	meta := map[string]any{
		"total":  count,
		"limit":  limit,
		"offset": offset,
	}
	response.OK(w, "Financial audit logs retrieved successfully", logs, meta)
}

