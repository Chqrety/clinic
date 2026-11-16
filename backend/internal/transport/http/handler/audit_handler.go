package handler

import (
	"net/http"

	auditApp "clinic-api/internal/application/audit"
	"clinic-api/internal/domain/audit"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
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
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 50,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"created_at":  "created_at",
			"action":      "action",
			"entity_type": "entity_type",
			"entity_id":   "entity_id",
		},
		DefaultSortBy:    "created_at",
		DefaultSortOrder: "DESC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := audit.AuditFilter{
		StartDate:  pg.StartDate,
		EndDate:    pg.EndDate,
		Search:     pg.Search,
		SortBy:     pg.SortBy,
		SortOrder:  pg.SortOrder,
		Limit:      pg.Limit,
		Offset:     pg.Offset,
		EntityType: r.URL.Query().Get("entity_type"),
	}

	if pg.BranchID != nil {
		filter.BranchID = pg.BranchID
	} else if bID, ok := middleware.GetBranchID(ctx); ok && bID != uuid.Nil {
		filter.BranchID = &bID
	}

	if aStr := r.URL.Query().Get("action"); aStr != "" {
		act := audit.ActionType(aStr)
		filter.Action = &act
	}

	logs, count, err := h.svc.ListLogs(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list audit logs: "+err.Error())
		return
	}

	response.OK(w, "Financial audit logs retrieved successfully", logs, pg.ToMeta(count))
}

