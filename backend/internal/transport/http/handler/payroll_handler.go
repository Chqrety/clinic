package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	payrollApp "clinic-api/internal/application/payroll"
	"clinic-api/internal/domain/payroll"
	"clinic-api/internal/transport/http/dto"
	"clinic-api/internal/transport/http/middleware"
	"clinic-api/pkg/response"

	"github.com/go-chi/chi/v5"
)

type PayrollHandler struct {
	svc *payrollApp.Service
}

func NewPayrollHandler(svc *payrollApp.Service) *PayrollHandler {
	return &PayrollHandler{svc: svc}
}

// ----------------------------------------------------------------------------
// Employee Endpoints
// ----------------------------------------------------------------------------

func (h *PayrollHandler) HandleListEmployees(w http.ResponseWriter, r *http.Request) {
	bu := r.URL.Query().Get("business_unit")
	activeOnlyStr := r.URL.Query().Get("active_only")
	activeOnly := activeOnlyStr == "true" || activeOnlyStr == "1"

	employees, err := h.svc.ListEmployees(r.Context(), bu, activeOnly)
	if err != nil {
		response.InternalServerError(w, "Gagal mengambil daftar karyawan: "+err.Error())
		return
	}
	response.OK(w, "Daftar karyawan berhasil dimuat", employees, nil)
}

func (h *PayrollHandler) HandleGetEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	emp, err := h.svc.GetEmployeeByID(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, "Gagal mengambil data karyawan: "+err.Error())
		return
	}
	if emp == nil {
		response.NotFound(w, "Karyawan tidak ditemukan")
		return
	}
	response.OK(w, "Data karyawan berhasil dimuat", emp, nil)
}

func (h *PayrollHandler) HandleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}

	emp := payroll.Employee{
		NIK:                    req.NIK,
		Name:                   req.Name,
		Position:               req.Position,
		Department:             req.Department,
		BusinessUnit:           req.BusinessUnit,
		EmploymentType:         req.EmploymentType,
		JoinDate:               req.JoinDate,
		BaseSalary:             req.BaseSalary,
		PositionAllowance:      req.PositionAllowance,
		TransportMealAllowance: req.TransportMealAllowance,
		BankName:               req.BankName,
		BankAccountNumber:      req.BankAccountNumber,
		NPWP:                   req.NPWP,
		PTKPStatus:             req.PTKPStatus,
		BPJSKesehatanNo:        req.BPJSKesehatanNo,
		BPJSKetenagakerjaanNo:  req.BPJSKetenagakerjaanNo,
		IsActive:               req.IsActive,
	}

	if err := h.svc.CreateEmployee(r.Context(), &emp); err != nil {
		response.InternalServerError(w, "Gagal membuat karyawan: "+err.Error())
		return
	}
	response.Created(w, "Karyawan berhasil didaftarkan", emp)
}

func (h *PayrollHandler) HandleUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}

	emp := payroll.Employee{
		ID:                     id,
		Name:                   req.Name,
		Position:               req.Position,
		Department:             req.Department,
		BusinessUnit:           req.BusinessUnit,
		EmploymentType:         req.EmploymentType,
		BaseSalary:             req.BaseSalary,
		PositionAllowance:      req.PositionAllowance,
		TransportMealAllowance: req.TransportMealAllowance,
		BankName:               req.BankName,
		BankAccountNumber:      req.BankAccountNumber,
		NPWP:                   req.NPWP,
		PTKPStatus:             req.PTKPStatus,
		BPJSKesehatanNo:        req.BPJSKesehatanNo,
		BPJSKetenagakerjaanNo:  req.BPJSKetenagakerjaanNo,
		IsActive:               req.IsActive,
	}

	if err := h.svc.UpdateEmployee(r.Context(), &emp); err != nil {
		response.InternalServerError(w, "Gagal mengupdate karyawan: "+err.Error())
		return
	}
	response.OK(w, "Karyawan berhasil diupdate", emp, nil)
}

// ----------------------------------------------------------------------------
// Payslip Endpoints
// ----------------------------------------------------------------------------

func (h *PayrollHandler) HandleListPayslips(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")
	bu := r.URL.Query().Get("business_unit")
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	slips, err := h.svc.ListPayslips(r.Context(), payroll.PayslipFilter{
		Month:        month,
		Year:         year,
		BusinessUnit: bu,
		Status:       status,
		Search:       search,
	})
	if err != nil {
		response.InternalServerError(w, "Gagal mengambil daftar slip gaji: "+err.Error())
		return
	}
	response.OK(w, "Daftar slip gaji berhasil dimuat", slips, nil)
}

func (h *PayrollHandler) HandleGetPayslip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	slip, err := h.svc.GetPayslipByID(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, "Gagal mengambil slip gaji: "+err.Error())
		return
	}
	if slip == nil {
		response.NotFound(w, "Slip gaji tidak ditemukan")
		return
	}
	response.OK(w, "Slip gaji berhasil dimuat", slip, nil)
}

func (h *PayrollHandler) HandleGenerateMonthlyPayroll(w http.ResponseWriter, r *http.Request) {
	var req dto.GeneratePayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}
	if req.Month < 1 || req.Month > 12 || req.Year < 2020 {
		response.BadRequest(w, "Bulan atau tahun periode tidak valid", nil)
		return
	}

	slips, err := h.svc.GenerateMonthlyPayroll(r.Context(), req.Month, req.Year, req.BusinessUnit)
	if err != nil {
		response.InternalServerError(w, "Gagal membuat slip gaji bulanan: "+err.Error())
		return
	}
	response.OK(w, "Slip gaji bulanan berhasil di-generate secara otomatis", slips, nil)
}

func (h *PayrollHandler) HandleCreatePayslip(w http.ResponseWriter, r *http.Request) {
	var slip payroll.Payslip
	if err := json.NewDecoder(r.Body).Decode(&slip); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}
	slip.Status = payroll.StatusDraft

	if err := h.svc.CreatePayslip(r.Context(), &slip); err != nil {
		response.InternalServerError(w, "Gagal membuat slip gaji: "+err.Error())
		return
	}
	response.Created(w, "Slip gaji berhasil dibuat", slip)
}

func (h *PayrollHandler) HandleUpdatePayslip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var slip payroll.Payslip
	if err := json.NewDecoder(r.Body).Decode(&slip); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}
	slip.ID = id

	if err := h.svc.UpdatePayslip(r.Context(), &slip); err != nil {
		response.InternalServerError(w, "Gagal mengupdate slip gaji: "+err.Error())
		return
	}
	response.OK(w, "Slip gaji berhasil diupdate", slip, nil)
}

func (h *PayrollHandler) HandleDeletePayslip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeletePayslip(r.Context(), id); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}
	response.OK(w, "Slip gaji berhasil dihapus", nil, nil)
}

func (h *PayrollHandler) HandleApprovePayslip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req dto.ApprovePayslipRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.ApprovedBy == "" {
		req.ApprovedBy = "Finance & Accounting Manager"
	}

	if err := h.svc.ApprovePayslip(r.Context(), id, req.ApprovedBy); err != nil {
		response.InternalServerError(w, "Gagal menyetujui slip gaji: "+err.Error())
		return
	}
	response.OK(w, "Slip gaji berhasil disetujui", nil, nil)
}

func (h *PayrollHandler) HandleApproveAllPayslips(w http.ResponseWriter, r *http.Request) {
	var req dto.ApproveAllPayslipsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}
	if req.ApprovedBy == "" {
		req.ApprovedBy = "Finance & Accounting Manager"
	}

	if err := h.svc.ApproveAllInPeriod(r.Context(), req.Month, req.Year, req.ApprovedBy); err != nil {
		response.InternalServerError(w, "Gagal menyetujui slip gaji batch: "+err.Error())
		return
	}
	response.OK(w, "Seluruh slip gaji periode berhasil disetujui", nil, nil)
}

func (h *PayrollHandler) HandlePayPayslip(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.MarkPayslipAsPaid(r.Context(), id); err != nil {
		response.InternalServerError(w, "Gagal menandai lunas slip gaji: "+err.Error())
		return
	}
	response.OK(w, "Slip gaji berhasil ditandai LUNAS (Disbursed)", nil, nil)
}

func (h *PayrollHandler) HandlePayAllPayslips(w http.ResponseWriter, r *http.Request) {
	var req dto.MarkAllPaidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}

	if err := h.svc.MarkAllAsPaidInPeriod(r.Context(), req.Month, req.Year); err != nil {
		response.InternalServerError(w, "Gagal menandai lunas batch: "+err.Error())
		return
	}
	response.OK(w, "Seluruh slip gaji periode berhasil ditandai LUNAS", nil, nil)
}

func (h *PayrollHandler) HandlePostPayrollLedger(w http.ResponseWriter, r *http.Request) {
	var req dto.PostPayrollLedgerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Payload JSON tidak valid: "+err.Error(), nil)
		return
	}

	branchID, _ := middleware.GetBranchID(r.Context())
	branchIDStr := branchID.String()

	jEntry, err := h.svc.PostPayrollToLedger(r.Context(), req.Month, req.Year, branchIDStr, "Finance Manager")
	if err != nil {
		response.InternalServerError(w, "Gagal memposting payroll ke buku besar: "+err.Error())
		return
	}
	response.OK(w, "Payroll berhasil diposting ke Buku Besar secara resmi", jEntry, nil)
}
