package handler

import (
	"encoding/json"
	"net/http"

	assetApp "clinic-api/internal/application/asset"
	"clinic-api/internal/domain/asset"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/pagination"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AssetHandler struct {
	svc *assetApp.Service
}

func NewAssetHandler(svc *assetApp.Service) *AssetHandler {
	return &AssetHandler{svc: svc}
}

func (h *AssetHandler) HandleRegisterAsset(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if ok && branchID != uuid.Nil {
		req.BranchID = branchID
	}

	assetItem, journal, err := h.svc.RegisterAsset(r.Context(), assetApp.RegisterAssetCommand{
		BranchID:            req.BranchID,
		CostCenterID:        req.CostCenterID,
		AssetCode:           req.AssetCode,
		AssetName:           req.AssetName,
		Category:            req.Category,
		AcquisitionDate:     req.AcquisitionDate,
		AcquisitionCost:     req.AcquisitionCost,
		SalvageValue:        req.SalvageValue,
		UsefulLifeMonths:    req.UsefulLifeMonths,
		CreditAccountCode:   req.CreditAccountCode,
		AutoPostAcquisition: req.AutoPostAcquisition,
	})
	if err != nil {
		response.BadRequest(w, "Failed to register fixed asset: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"asset":   assetItem,
		"journal": journal,
	}
	response.Created(w, "Fixed asset registered successfully", data)
}

func (h *AssetHandler) HandleListAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pg, err := pagination.Parse(r, pagination.Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		AllowedSortColumns: map[string]string{
			"asset_code":       "asset_code",
			"asset_name":       "asset_name",
			"category":         "category",
			"acquisition_date": "acquisition_date",
			"acquisition_cost": "acquisition_cost",
			"status":           "status",
			"created_at":       "created_at",
		},
		DefaultSortBy:    "asset_code",
		DefaultSortOrder: "ASC",
	})
	if err != nil {
		response.BadRequest(w, "Invalid pagination parameters", err.Error())
		return
	}

	filter := asset.AssetFilter{
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
	if cat := q.Get("category"); cat != "" {
		filter.Category = cat
	}
	if st := q.Get("status"); st != "" {
		status := asset.AssetStatus(st)
		filter.Status = &status
	}
	if ccStr := q.Get("cost_center_id"); ccStr != "" {
		if ccID, err := uuid.Parse(ccStr); err == nil {
			filter.CostCenterID = &ccID
		}
	}

	assets, total, err := h.svc.ListAssets(ctx, filter)
	if err != nil {
		response.InternalServerError(w, "Failed to list assets: "+err.Error())
		return
	}

	response.OK(w, "Fixed assets retrieved successfully", assets, pg.ToMeta(total))
}

func (h *AssetHandler) HandleGetAssetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid asset ID", nil)
		return
	}

	assetItem, err := h.svc.GetAssetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Fixed asset not found")
		return
	}

	response.OK(w, "Fixed asset retrieved successfully", assetItem, nil)
}

func (h *AssetHandler) HandleDepreciateAsset(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid asset ID", nil)
		return
	}

	var req dto.DepreciateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	dep, journal, err := h.svc.DepreciateAsset(r.Context(), assetApp.DepreciateAssetCommand{
		AssetID:     id,
		PeriodYear:  req.PeriodYear,
		PeriodMonth: req.PeriodMonth,
	})
	if err != nil {
		response.BadRequest(w, "Failed to depreciate asset: "+err.Error(), nil)
		return
	}

	data := map[string]any{
		"depreciation": dep,
		"journal":      journal,
	}
	response.OK(w, "Fixed asset depreciation processed and journal posted", data, nil)
}

func (h *AssetHandler) HandleRunMonthlyDepreciation(w http.ResponseWriter, r *http.Request) {
	var req dto.RunMonthlyDepreciationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload: "+err.Error(), nil)
		return
	}

	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(r.Context()); ok && bID != uuid.Nil {
		branchID = &bID
	}

	results, err := h.svc.RunMonthlyDepreciationAll(r.Context(), branchID, req.PeriodYear, req.PeriodMonth)
	if err != nil {
		response.InternalServerError(w, "Failed to run monthly depreciation: "+err.Error())
		return
	}

	response.OK(w, "Monthly depreciation batch executed successfully", results, nil)
}

func (h *AssetHandler) HandleGetDepreciationHistory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid asset ID", nil)
		return
	}

	history, err := h.svc.GetDepreciationHistory(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, "Failed to get depreciation history: "+err.Error())
		return
	}

	response.OK(w, "Depreciation history retrieved successfully", history, nil)
}

