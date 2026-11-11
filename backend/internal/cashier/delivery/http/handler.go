package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"clinic-api/internal/cashier/domain"
	"clinic-api/internal/cashier/usecase"
	"clinic-api/internal/shared/middleware"
	"clinic-api/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type CashierHandler struct {
	usecase *usecase.CashierUsecase
}

func NewCashierHandler(u *usecase.CashierUsecase) *CashierHandler {
	return &CashierHandler{usecase: u}
}

func getStaffInfo(r *http.Request) (int64, string) {
	var staffID int64
	var staffName string

	if idVal := r.Context().Value(middleware.UserIDKey); idVal != nil {
		if id, ok := idVal.(int64); ok {
			staffID = id
		}
	}
	if nameVal := r.Context().Value(middleware.UserNameKey); nameVal != nil {
		if name, ok := nameVal.(string); ok {
			staffName = name
		}
	}
	if staffName == "" {
		staffName = "Kasir Meja Depan"
	}
	return staffID, staffName
}

func (h *CashierHandler) HandleRegisterWalkIn(w http.ResponseWriter, r *http.Request) {
	var req domain.WalkInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload tidak valid")
		return
	}

	staffID, staffName := getStaffInfo(r)
	bill, inv, err := h.usecase.RegisterWalkIn(r.Context(), &req, staffID, staffName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPatientNameRequired), errors.Is(err, domain.ErrPatientPhoneRequired):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal registrasi walk-in: "+err.Error())
		}
		return
	}

	type resultData struct {
		Bill    *domain.PendingBill `json:"bill"`
		Invoice *domain.Invoice     `json:"invoice,omitempty"`
	}

	response.Success(w, http.StatusCreated, resultData{
		Bill:    bill,
		Invoice: inv,
	})
}

func (h *CashierHandler) HandleListPendingBills(w http.ResponseWriter, r *http.Request) {
	bills, err := h.usecase.ListPendingBills(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat daftar tagihan antrian: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, bills)
}

func (h *CashierHandler) HandlePayBill(w http.ResponseWriter, r *http.Request) {
	var req domain.PayBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload tidak valid")
		return
	}

	staffID, staffName := getStaffInfo(r)
	inv, err := h.usecase.PayBill(r.Context(), &req, staffID, staffName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrQueueNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInvalidAmount),
			errors.Is(err, domain.ErrInsufficientCash),
			errors.Is(err, domain.ErrVoucherExceedsTotal):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memproses pembayaran: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusCreated, inv)
}

func (h *CashierHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID invoice tidak valid")
		return
	}

	inv, err := h.usecase.GetInvoice(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrInvoiceNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal memuat invoice: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, inv)
}

func (h *CashierHandler) HandleGetInvoiceByQueue(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "queueId")
	qID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID antrian tidak valid")
		return
	}

	inv, err := h.usecase.GetInvoiceByQueue(r.Context(), qID)
	if err != nil {
		if errors.Is(err, domain.ErrInvoiceNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal memuat invoice antrian: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, inv)
}

func (h *CashierHandler) HandleListInvoices(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	invoices, err := h.usecase.ListInvoices(r.Context(), limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat invoice: "+err.Error())
		return
	}

	response.Success(w, http.StatusOK, invoices)
}

func (h *CashierHandler) HandleRefund(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID invoice tidak valid")
		return
	}

	var req domain.RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "format payload refund tidak valid")
		return
	}

	staffID, staffName := getStaffInfo(r)
	inv, err := h.usecase.ProcessRefund(r.Context(), id, &req, staffID, staffName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvoiceNotFound):
			response.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrAlreadyRefunded), errors.Is(err, domain.ErrRefundReasonRequired):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memproses refund: "+err.Error())
		}
		return
	}

	response.Success(w, http.StatusOK, inv)
}

func (h *CashierHandler) HandleGetShiftSummary(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	sum, err := h.usecase.GetShiftSummary(r.Context(), date)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil ringkasan shift: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, sum)
}

func (h *CashierHandler) HandleCloseShift(w http.ResponseWriter, r *http.Request) {
	var req domain.CloseShiftRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	staffID, staffName := getStaffInfo(r)
	report, err := h.usecase.CloseShift(r.Context(), staffID, staffName, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal menutup shift kasir: "+err.Error())
		return
	}
	response.Success(w, http.StatusCreated, report)
}

func (h *CashierHandler) HandleListShiftReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.usecase.ListShiftReports(r.Context(), 30)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil daftar laporan shift: "+err.Error())
		return
	}
	response.Success(w, http.StatusOK, reports)
}
