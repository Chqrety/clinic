package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	cashbankApp "clinic-api/internal/application/cashbank"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/internal/domain/cashbank"
	"clinic-api/internal/domain/shared"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CashBankHandler struct {
	cashbankService   *cashbankApp.Service
	accountingService *accountingApp.Service
}

func NewCashBankHandler(cashbankService *cashbankApp.Service, accountingService *accountingApp.Service) *CashBankHandler {
	return &CashBankHandler{
		cashbankService:   cashbankService,
		accountingService: accountingService,
	}
}

func (h *CashBankHandler) RegisterRoutes(r chi.Router) {
	r.Route("/cash-bank", func(cb chi.Router) {
		cb.Get("/accounts", h.HandleListAccounts)
		cb.Post("/accounts", h.HandleCreateAccount)
		cb.Post("/cash-in", h.HandleRecordCashIn)
		cb.Post("/cash-out", h.HandleRecordCashOut)
		cb.Post("/transfer", h.HandleTransfer)
		cb.Get("/transactions", h.HandleListTransactions)
	})

	r.Route("/periods", func(p chi.Router) {
		p.Post("/close", h.HandleClosePeriod)
		p.Post("/{id}/reopen", h.HandleReopenPeriod)
	})
}

func (h *CashBankHandler) HandleListAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var branchID *uuid.UUID
	if bID, ok := middleware.GetBranchID(ctx); ok {
		branchID = &bID
	}

	var accType *cashbank.AccountType
	if tStr := r.URL.Query().Get("type"); tStr != "" {
		t := cashbank.AccountType(tStr)
		accType = &t
	}

	accounts, err := h.cashbankService.ListAccounts(ctx, branchID, accType)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.CashBankAccountResponse, 0, len(accounts))
	for _, a := range accounts {
		resList = append(resList, dto.CashBankAccountResponse{
			ID:            a.ID,
			BranchID:      a.BranchID,
			AccountID:     a.AccountID,
			AccountCode:   a.AccountCode,
			Code:          a.Code,
			Name:          a.Name,
			Type:          a.Type,
			BankName:      a.BankName,
			AccountNumber: a.AccountNumber,
			IsActive:      a.IsActive,
			CreatedAt:     a.CreatedAt,
		})
	}

	response.OK(w, "Cash and bank accounts retrieved successfully", resList, nil)
}

func (h *CashBankHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCashBankAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New()
	}

	acc := &cashbank.CashBankAccount{
		BaseEntity:    shared.BaseEntity{ID: uuid.New()},
		BranchScoped:  shared.BranchScoped{BranchID: branchID},
		AccountID:     req.AccountID,
		Code:          req.Code,
		Name:          req.Name,
		Type:          req.Type,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		IsActive:      true,
	}

	if err := h.cashbankService.CreateAccount(r.Context(), acc); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.Created(w, "Cash/Bank account created successfully", dto.CashBankAccountResponse{
		ID:            acc.ID,
		BranchID:      acc.BranchID,
		AccountID:     acc.AccountID,
		Code:          acc.Code,
		Name:          acc.Name,
		Type:          acc.Type,
		BankName:      acc.BankName,
		AccountNumber: acc.AccountNumber,
		IsActive:      acc.IsActive,
		CreatedAt:     acc.CreatedAt,
	})
}

func (h *CashBankHandler) HandleRecordCashIn(w http.ResponseWriter, r *http.Request) {
	var req dto.RecordCashInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New()
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		txDate = time.Now()
	}

	cmd := cashbankApp.RecordCashInCommand{
		BranchID:             branchID,
		CashBankAccountID:    req.CashBankAccountID,
		Amount:               req.Amount,
		TransactionDate:      txDate,
		CounterpartAccountID: req.CounterpartAccountID,
		Description:          req.Description,
		ReferenceNo:          req.ReferenceNo,
	}

	txItem, journal, err := h.cashbankService.RecordCashIn(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	type ResultPayload struct {
		Transaction   dto.CashBankTransactionResponse `json:"transaction"`
		JournalNumber string                          `json:"journal_number,omitempty"`
	}

	var jNum string
	if journal != nil {
		jNum = journal.JournalNumber
	}

	response.Created(w, "Cash in recorded and journal posted successfully", ResultPayload{
		Transaction:   toTransactionResponse(txItem),
		JournalNumber: jNum,
	})
}

func (h *CashBankHandler) HandleRecordCashOut(w http.ResponseWriter, r *http.Request) {
	var req dto.RecordCashOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New()
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		txDate = time.Now()
	}

	cmd := cashbankApp.RecordCashOutCommand{
		BranchID:             branchID,
		CashBankAccountID:    req.CashBankAccountID,
		Amount:               req.Amount,
		TransactionDate:      txDate,
		CounterpartAccountID: req.CounterpartAccountID,
		Description:          req.Description,
		ReferenceNo:          req.ReferenceNo,
	}

	txItem, journal, err := h.cashbankService.RecordCashOut(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	type ResultPayload struct {
		Transaction   dto.CashBankTransactionResponse `json:"transaction"`
		JournalNumber string                          `json:"journal_number,omitempty"`
	}

	var jNum string
	if journal != nil {
		jNum = journal.JournalNumber
	}

	response.Created(w, "Cash out recorded and journal posted successfully", ResultPayload{
		Transaction:   toTransactionResponse(txItem),
		JournalNumber: jNum,
	})
}

func (h *CashBankHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	var req dto.TransferCashBankRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New()
	}

	txDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		txDate = time.Now()
	}

	cmd := cashbankApp.TransferCommand{
		BranchID:        branchID,
		SourceAccountID: req.SourceAccountID,
		TargetAccountID: req.TargetAccountID,
		Amount:          req.Amount,
		TransactionDate: txDate,
		Description:     req.Description,
		ReferenceNo:     req.ReferenceNo,
	}

	txItem, journal, err := h.cashbankService.TransferCashBank(r.Context(), cmd)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	type ResultPayload struct {
		Transaction   dto.CashBankTransactionResponse `json:"transaction"`
		JournalNumber string                          `json:"journal_number,omitempty"`
	}

	var jNum string
	if journal != nil {
		jNum = journal.JournalNumber
	}

	response.Created(w, "Transfer processed and journal posted successfully", ResultPayload{
		Transaction:   toTransactionResponse(txItem),
		JournalNumber: jNum,
	})
}

func (h *CashBankHandler) HandleListTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var filter cashbank.TransactionFilter

	if bID, ok := middleware.GetBranchID(ctx); ok {
		filter.BranchID = &bID
	}
	if accIDStr := r.URL.Query().Get("account_id"); accIDStr != "" {
		if id, err := uuid.Parse(accIDStr); err == nil {
			filter.CashBankAccountID = &id
		}
	}
	if tTypeStr := r.URL.Query().Get("type"); tTypeStr != "" {
		t := cashbank.TransactionType(tTypeStr)
		filter.TransactionType = &t
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

	txs, totalCount, err := h.cashbankService.ListTransactions(ctx, filter)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	resList := make([]dto.CashBankTransactionResponse, 0, len(txs))
	for _, t := range txs {
		resList = append(resList, toTransactionResponse(&t))
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))
	meta := response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
	}

	response.OK(w, "Transactions retrieved successfully", resList, meta)
}

func (h *CashBankHandler) HandleClosePeriod(w http.ResponseWriter, r *http.Request) {
	var req dto.ClosePeriodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid JSON payload", err.Error())
		return
	}

	if req.Year <= 0 || req.Month < 1 || req.Month > 12 {
		response.BadRequest(w, "Valid year and month (1-12) are required", nil)
		return
	}

	branchID, ok := middleware.GetBranchID(r.Context())
	if !ok {
		branchID = uuid.New()
	}

	period, err := h.accountingService.ClosePeriod(r.Context(), branchID, req.Year, req.Month, nil, req.Notes)
	if err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Fiscal period successfully closed and locked", period, nil)
}

func (h *CashBankHandler) HandleReopenPeriod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	periodID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid period UUID", nil)
		return
	}

	var req dto.ReopenPeriodRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		req.Reason = "Otorisasi pembukaan kembali periode oleh manajemen"
	}

	if err := h.accountingService.ReopenPeriod(r.Context(), periodID, req.Reason); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	response.OK(w, "Fiscal period successfully reopened", nil, nil)
}

func toTransactionResponse(t *cashbank.CashBankTransaction) dto.CashBankTransactionResponse {
	return dto.CashBankTransactionResponse{
		ID:                   t.ID,
		BranchID:             t.BranchID,
		CashBankAccountID:    t.CashBankAccountID,
		CashBankAccountCode:  t.CashBankAccountCode,
		CashBankAccountName:  t.CashBankAccountName,
		TransactionNumber:    t.TransactionNumber,
		TransactionType:      t.TransactionType,
		Amount:               t.Amount,
		TransactionDate:      t.TransactionDate.Format("2006-01-02"),
		CounterpartAccountID: t.CounterpartAccountID,
		CounterpartCode:      t.CounterpartCode,
		CounterpartName:      t.CounterpartName,
		Description:          t.Description,
		ReferenceNo:          t.ReferenceNo,
		JournalID:            t.JournalID,
		CreatedAt:            t.CreatedAt,
	}
}

