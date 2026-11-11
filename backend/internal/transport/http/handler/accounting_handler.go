package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	app "clinic-api/internal/application/accounting"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	domain "clinic-api/internal/domain/accounting"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AccountingHandler struct {
	service *app.Service
}

func NewAccountingHandler(service *app.Service) *AccountingHandler {
	return &AccountingHandler{service: service}
}

func (h *AccountingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(acc chi.Router) {
		acc.Get("/", h.HandleListAccounts)
		acc.Post("/", h.HandleCreateAccount)
		acc.Get("/{id}", h.HandleGetAccountByID)
	})

	r.Route("/journals", func(j chi.Router) {
		j.Get("/", h.HandleListJournals)
		j.Post("/", h.HandleCreateJournal)
		j.Get("/{id}", h.HandleGetJournalByID)
		j.Post("/{id}/post", h.HandlePostJournal)
		j.Post("/{id}/reverse", h.HandleReverseJournal)
	})
}

func (h *AccountingHandler) HandleListAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var filter domain.AccountFilter

	if accType := r.URL.Query().Get("type"); accType != "" {
		t := domain.AccountType(accType)
		filter.Type = &t
	}
	if activeStr := r.URL.Query().Get("is_active"); activeStr != "" {
		active := activeStr == "true"
		filter.IsActive = &active
	}

	accounts, err := h.service.GetAccounts(ctx, filter)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.AccountResponse, 0, len(accounts))
	for _, a := range accounts {
		resList = append(resList, dto.AccountResponse{
			ID:            a.ID,
			Code:          a.Code,
			Name:          a.Name,
			Type:          a.Type,
			NormalBalance: a.NormalBalance,
			ParentID:      a.ParentID,
			Level:         a.Level,
			IsHeader:      a.IsHeader,
			IsActive:      a.IsActive,
			Description:   a.Description,
		})
	}

	response.OK(w, "Accounts retrieved successfully", resList, nil)
}

func (h *AccountingHandler) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid account ID UUID", nil)
		return
	}

	acc, err := h.service.GetAccountByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Account not found")
		return
	}

	response.OK(w, "Account retrieved successfully", dto.AccountResponse{
		ID:            acc.ID,
		Code:          acc.Code,
		Name:          acc.Name,
		Type:          acc.Type,
		NormalBalance: acc.NormalBalance,
		ParentID:      acc.ParentID,
		Level:         acc.Level,
		IsHeader:      acc.IsHeader,
		IsActive:      acc.IsActive,
		Description:   acc.Description,
	}, nil)
}

func (h *AccountingHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	if req.Code == "" || req.Name == "" || req.Type == "" || req.NormalBalance == "" {
		response.BadRequest(w, "Fields code, name, type, and normal_balance are required", nil)
		return
	}

	acc := &domain.Account{
		Code:          req.Code,
		Name:          req.Name,
		Type:          req.Type,
		NormalBalance: req.NormalBalance,
		ParentID:      req.ParentID,
		Level:         req.Level,
		IsHeader:      req.IsHeader,
		IsActive:      true,
		Description:   req.Description,
	}

	if err := h.service.CreateAccount(r.Context(), acc); err != nil {
		if strings.Contains(err.Error(), "accounts_code_key") || strings.Contains(err.Error(), "23505") {
			response.BadRequest(w, fmt.Sprintf("Akun dengan kode '%s' sudah terdaftar dalam Chart of Accounts (COA). Gunakan kode akun lain yang unik.", req.Code), nil)
			return
		}
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Account created successfully", dto.AccountResponse{
		ID:            acc.ID,
		Code:          acc.Code,
		Name:          acc.Name,
		Type:          acc.Type,
		NormalBalance: acc.NormalBalance,
		ParentID:      acc.ParentID,
		Level:         acc.Level,
		IsHeader:      acc.IsHeader,
		IsActive:      acc.IsActive,
		Description:   acc.Description,
	})
}

func (h *AccountingHandler) HandleListJournals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var filter domain.JournalFilter

	if branchID, ok := middleware.GetBranchID(ctx); ok {
		filter.BranchID = &branchID
	}

	if jType := r.URL.Query().Get("type"); jType != "" {
		t := domain.JournalType(jType)
		filter.Type = &t
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := domain.JournalStatus(status)
		filter.Status = &s
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	filter.Limit = limit

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	filter.Offset = (page - 1) * limit

	journals, totalCount, err := h.service.ListJournals(ctx, filter)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.JournalResponse, 0, len(journals))
	for _, j := range journals {
		resList = append(resList, toJournalResponse(&j))
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))
	meta := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
	}

	response.OK(w, "Journals retrieved successfully", resList, meta)
}

func (h *AccountingHandler) HandleGetJournalByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid journal ID UUID", nil)
		return
	}

	entry, err := h.service.GetJournalByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Journal entry not found")
		return
	}

	response.OK(w, "Journal entry retrieved successfully", toJournalResponse(entry), nil)
}

func (h *AccountingHandler) HandleCreateJournal(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New() // Fallback if single-branch or dev
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		txDate = time.Now()
	}

	lines := make([]app.CreateLineCommand, 0, len(req.Lines))
	for _, l := range req.Lines {
		lines = append(lines, app.CreateLineCommand{
			AccountID:   l.AccountID,
			AccountCode: l.AccountCode,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Description: l.Description,
		})
	}

	cmd := app.CreateJournalCommand{
		BranchID:        branchID,
		JournalType:     req.JournalType,
		TransactionDate: txDate,
		ReferenceType:   req.ReferenceType,
		ReferenceID:     req.ReferenceID,
		Description:     req.Description,
		AutoPost:        req.AutoPost,
		Lines:           lines,
	}

	entry, err := h.service.CreateJournal(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Journal entry created successfully", toJournalResponse(entry))
}

func (h *AccountingHandler) HandlePostJournal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid journal ID UUID", nil)
		return
	}

	entry, err := h.service.PostJournal(r.Context(), id, nil)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Journal entry posted successfully to General Ledger", toJournalResponse(entry), nil)
}

func (h *AccountingHandler) HandleReverseJournal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid journal ID UUID", nil)
		return
	}

	var req dto.ReverseJournalRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "Koreksi manual"
	}

	revEntry, err := h.service.ReverseJournal(r.Context(), id, req.Reason, nil)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Reversal journal created successfully", toJournalResponse(revEntry))
}

func toJournalResponse(j *domain.JournalEntry) dto.JournalResponse {
	lineResponses := make([]dto.JournalLineResponse, 0, len(j.Lines))
	for _, l := range j.Lines {
		lineResponses = append(lineResponses, dto.JournalLineResponse{
			ID:          l.ID,
			AccountID:   l.AccountID,
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Debit:       l.Debit,
			Credit:      l.Credit,
			Description: l.Description,
			LineOrder:   l.LineOrder,
		})
	}

	return dto.JournalResponse{
		ID:              j.ID,
		BranchID:        j.BranchID,
		JournalNumber:   j.JournalNumber,
		JournalType:     j.JournalType,
		TransactionDate: j.TransactionDate.Format("2006-01-02"),
		PostingDate:     j.PostingDate,
		Status:          j.Status,
		ReferenceType:   j.ReferenceType,
		ReferenceID:     j.ReferenceID,
		Description:     j.Description,
		TotalDebit:      j.TotalDebit,
		TotalCredit:     j.TotalCredit,
		PostedBy:        j.PostedBy,
		CreatedBy:       j.CreatedBy,
		Lines:           lineResponses,
		CreatedAt:       j.CreatedAt,
	}
}

