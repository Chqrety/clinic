package payroll

import (
	"context"
	"fmt"
	"math"
	"time"

	accountingApp "clinic-api/internal/application/accounting"
	accountingDomain "clinic-api/internal/domain/accounting"
	"clinic-api/internal/domain/payroll"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	uow               UnitOfWork
	repo              payroll.Repository
	accountingService *accountingApp.Service
}

func NewService(uow UnitOfWork, repo payroll.Repository, accountingService *accountingApp.Service) *Service {
	return &Service{
		uow:               uow,
		repo:              repo,
		accountingService: accountingService,
	}
}

// ----------------------------------------------------------------------------
// Tax & Deduction Calculators (Indonesian Regulatory Standard)
// ----------------------------------------------------------------------------

func CalculateBPJSKesehatan(gross decimal.Decimal) decimal.Decimal {
	g := gross.InexactFloat64()
	capped := math.Min(g, 12000000)
	fee := math.Round(capped * 0.01)
	return decimal.FromFloat(fee)
}

func CalculateBPJSKetenagakerjaan(gross decimal.Decimal) decimal.Decimal {
	g := gross.InexactFloat64()
	fee := math.Round(g * 0.03) // 2% JHT + 1% JP
	return decimal.FromFloat(fee)
}

// CalculatePPh21TER calculates monthly withholding using PP 58/2023 (Kategori A, B, C)
func CalculatePPh21TER(gross decimal.Decimal, ptkpStatus string) decimal.Decimal {
	g := gross.InexactFloat64()
	if g <= 5400000 {
		return decimal.Zero
	}

	category := "A"
	switch ptkpStatus {
	case "TK/0", "TK/1", "K/0":
		category = "A"
	case "TK/2", "TK/3", "K/1", "K/2":
		category = "B"
	case "K/3":
		category = "C"
	default:
		category = "A"
	}

	var rate float64
	if category == "A" {
		if g <= 5400000 {
			rate = 0.0
		} else if g <= 5650000 {
			rate = 0.0025
		} else if g <= 5950000 {
			rate = 0.005
		} else if g <= 6300000 {
			rate = 0.0075
		} else if g <= 6750000 {
			rate = 0.01
		} else if g <= 7500000 {
			rate = 0.0125
		} else if g <= 8550000 {
			rate = 0.015
		} else if g <= 9650000 {
			rate = 0.0175
		} else if g <= 10050000 {
			rate = 0.02
		} else if g <= 10350000 {
			rate = 0.0225
		} else if g <= 10700000 {
			rate = 0.025
		} else if g <= 12500000 {
			rate = 0.03
		} else if g <= 15000000 {
			rate = 0.05
		} else {
			rate = 0.07
		}
	} else if category == "B" {
		if g <= 6200000 {
			rate = 0.0
		} else if g <= 6500000 {
			rate = 0.0025
		} else if g <= 6850000 {
			rate = 0.005
		} else if g <= 7300000 {
			rate = 0.0075
		} else if g <= 9200000 {
			rate = 0.01
		} else if g <= 10750000 {
			rate = 0.015
		} else if g <= 12500000 {
			rate = 0.02
		} else if g <= 15000000 {
			rate = 0.04
		} else {
			rate = 0.06
		}
	} else {
		// Category C
		if g <= 6600000 {
			rate = 0.0
		} else if g <= 6950000 {
			rate = 0.0025
		} else if g <= 7350000 {
			rate = 0.005
		} else if g <= 7800000 {
			rate = 0.0075
		} else if g <= 8850000 {
			rate = 0.01
		} else if g <= 10000000 {
			rate = 0.0125
		} else if g <= 15000000 {
			rate = 0.03
		} else {
			rate = 0.05
		}
	}

	pph := math.Round(g * rate)
	return decimal.FromFloat(pph)
}

func getMonthNameID(month int) string {
	names := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month >= 1 && month <= 12 {
		return names[month]
	}
	return "Bulan"
}

// ----------------------------------------------------------------------------
// Service Methods
// ----------------------------------------------------------------------------

func (s *Service) ListEmployees(ctx context.Context, businessUnit string, isActiveOnly bool) ([]payroll.Employee, error) {
	return s.repo.ListEmployees(ctx, businessUnit, isActiveOnly)
}

func (s *Service) GetEmployeeByID(ctx context.Context, id string) (*payroll.Employee, error) {
	return s.repo.GetEmployeeByID(ctx, id)
}

func (s *Service) CreateEmployee(ctx context.Context, emp *payroll.Employee) error {
	return s.repo.CreateEmployee(ctx, emp)
}

func (s *Service) UpdateEmployee(ctx context.Context, emp *payroll.Employee) error {
	return s.repo.UpdateEmployee(ctx, emp)
}

func (s *Service) ListPayslips(ctx context.Context, filter payroll.PayslipFilter) ([]payroll.Payslip, error) {
	return s.repo.ListPayslips(ctx, filter)
}

func (s *Service) GetPayslipByID(ctx context.Context, id string) (*payroll.Payslip, error) {
	return s.repo.GetPayslipByID(ctx, id)
}

func (s *Service) CreatePayslip(ctx context.Context, slip *payroll.Payslip) error {
	if slip.WorkingDays <= 0 {
		slip.WorkingDays = 22
	}
	if slip.PresentDays == 0 && slip.AbsentDays == 0 && slip.SickLeaveDays == 0 && slip.PaidLeaveDays == 0 {
		slip.PresentDays = slip.WorkingDays
	}
	if slip.DailyMealTransportRate.IsPositive() && slip.PresentDays > 0 {
		slip.TransportMealAllowance = slip.DailyMealTransportRate.Mul(decimal.FromInt(int64(slip.PresentDays)))
	}
	if slip.AbsentDays > 0 && slip.AttendanceDeduction.IsZero() && slip.WorkingDays > 0 {
		deduction := math.Round(float64(slip.AbsentDays) * (slip.BaseSalary.InexactFloat64() / float64(slip.WorkingDays)))
		slip.AttendanceDeduction = decimal.FromFloat(deduction)
	}
	if slip.OvertimeHours > 0 && slip.OvertimePay.IsZero() {
		hourlyRate := (slip.BaseSalary.Add(slip.PositionAllowance).InexactFloat64() / 173.0) * 1.5
		slip.OvertimePay = decimal.FromFloat(math.Round(slip.OvertimeHours * hourlyRate))
	}

	// Auto calculate total gross, deductions and net if not set
	gross := slip.BaseSalary.Add(slip.PositionAllowance).
		Add(slip.TransportMealAllowance).
		Add(slip.ShiftAttendanceAllowance).
		Add(slip.IncentiveBonus).
		Add(slip.OvertimePay).
		Add(slip.OtherAllowance)
	slip.TotalGrossSalary = gross

	if slip.BPJSKesehatan.IsZero() {
		slip.BPJSKesehatan = CalculateBPJSKesehatan(gross)
	}
	if slip.BPJSKetenagakerjaan.IsZero() {
		slip.BPJSKetenagakerjaan = CalculateBPJSKetenagakerjaan(gross)
	}
	if slip.PPh21.IsZero() {
		slip.PPh21 = CalculatePPh21TER(gross, "TK/0")
	}

	ded := slip.BPJSKesehatan.Add(slip.BPJSKetenagakerjaan).
		Add(slip.PPh21).
		Add(slip.AttendanceDeduction).
		Add(slip.LoanDeduction).
		Add(slip.OtherDeduction)
	slip.TotalDeduction = ded
	slip.TakeHomePay = gross.Sub(ded)

	return s.repo.CreatePayslip(ctx, slip)
}

func (s *Service) UpdatePayslip(ctx context.Context, slip *payroll.Payslip) error {
	if slip.WorkingDays <= 0 {
		slip.WorkingDays = 22
	}
	if slip.DailyMealTransportRate.IsPositive() && slip.PresentDays > 0 {
		slip.TransportMealAllowance = slip.DailyMealTransportRate.Mul(decimal.FromInt(int64(slip.PresentDays)))
	}
	if slip.AbsentDays > 0 && slip.AttendanceDeduction.IsZero() && slip.WorkingDays > 0 {
		deduction := math.Round(float64(slip.AbsentDays) * (slip.BaseSalary.InexactFloat64() / float64(slip.WorkingDays)))
		slip.AttendanceDeduction = decimal.FromFloat(deduction)
	}
	if slip.OvertimeHours > 0 && slip.OvertimePay.IsZero() {
		hourlyRate := (slip.BaseSalary.Add(slip.PositionAllowance).InexactFloat64() / 173.0) * 1.5
		slip.OvertimePay = decimal.FromFloat(math.Round(slip.OvertimeHours * hourlyRate))
	}

	gross := slip.BaseSalary.Add(slip.PositionAllowance).
		Add(slip.TransportMealAllowance).
		Add(slip.ShiftAttendanceAllowance).
		Add(slip.IncentiveBonus).
		Add(slip.OvertimePay).
		Add(slip.OtherAllowance)
	slip.TotalGrossSalary = gross

	ded := slip.BPJSKesehatan.Add(slip.BPJSKetenagakerjaan).
		Add(slip.PPh21).
		Add(slip.AttendanceDeduction).
		Add(slip.LoanDeduction).
		Add(slip.OtherDeduction)
	slip.TotalDeduction = ded
	slip.TakeHomePay = gross.Sub(ded)

	return s.repo.UpdatePayslip(ctx, slip)
}

func (s *Service) DeletePayslip(ctx context.Context, id string) error {
	return s.repo.DeletePayslip(ctx, id)
}

func (s *Service) ApprovePayslip(ctx context.Context, id string, approvedBy string) error {
	return s.repo.ApprovePayslip(ctx, id, approvedBy)
}

func (s *Service) ApproveAllInPeriod(ctx context.Context, month, year int, approvedBy string) error {
	return s.repo.ApproveAllInPeriod(ctx, month, year, approvedBy)
}

func (s *Service) MarkPayslipAsPaid(ctx context.Context, id string) error {
	return s.repo.MarkPayslipAsPaid(ctx, id)
}

func (s *Service) MarkAllAsPaidInPeriod(ctx context.Context, month, year int) error {
	return s.repo.MarkAllAsPaidInPeriod(ctx, month, year)
}

// GenerateMonthlyPayroll generates and saves payslips for all active employees for the given period
func (s *Service) GenerateMonthlyPayroll(ctx context.Context, month, year int, businessUnit string) ([]payroll.Payslip, error) {
	employees, err := s.repo.ListEmployees(ctx, businessUnit, true)
	if err != nil {
		return nil, fmt.Errorf("list active employees: %w", err)
	}

	periodName := fmt.Sprintf("%s %d", getMonthNameID(month), year)
	var generated []payroll.Payslip

	for _, emp := range employees {
		workingDays := 22
		presentDays := 22
		dailyMealRate := decimal.FromFloat(30000)
		if emp.TransportMealAllowance.IsPositive() {
			dailyMealRate = decimal.FromFloat(math.Round(emp.TransportMealAllowance.InexactFloat64() / float64(workingDays)))
		}
		transportAllowance := dailyMealRate.Mul(decimal.FromInt(int64(presentDays)))

		gross := emp.BaseSalary.Add(emp.PositionAllowance).Add(transportAllowance)
		bpjsKes := CalculateBPJSKesehatan(gross)
		bpjsTk := CalculateBPJSKetenagakerjaan(gross)
		pph21 := CalculatePPh21TER(gross, emp.PTKPStatus)
		totalDed := bpjsKes.Add(bpjsTk).Add(pph21)
		thp := gross.Sub(totalDed)

		slipNum := fmt.Sprintf("SLIP/%04d/%02d/%s", year, month, emp.NIK)
		slipID := fmt.Sprintf("slip-%04d%02d-%s", year, month, emp.ID)

		slip := payroll.Payslip{
			ID:                     slipID,
			SlipNumber:             slipNum,
			EmployeeID:             emp.ID,
			EmployeeNIK:            emp.NIK,
			EmployeeName:           emp.Name,
			Position:               emp.Position,
			Department:             emp.Department,
			BusinessUnit:           emp.BusinessUnit,
			PeriodMonth:            month,
			PeriodYear:             year,
			PeriodName:             periodName,
			PaymentDate:            fmt.Sprintf("%04d-%02d-28", year, month),
			WorkingDays:            workingDays,
			PresentDays:            presentDays,
			SickLeaveDays:          0,
			PaidLeaveDays:          0,
			AbsentDays:             0,
			OvertimeHours:          0,
			DailyMealTransportRate: dailyMealRate,
			BaseSalary:             emp.BaseSalary,
			PositionAllowance:      emp.PositionAllowance,
			TransportMealAllowance: transportAllowance,
			TotalGrossSalary:       gross,
			BPJSKesehatan:          bpjsKes,
			BPJSKetenagakerjaan:    bpjsTk,
			PPh21:                  pph21,
			TotalDeduction:         totalDed,
			TakeHomePay:            thp,
			BankName:               emp.BankName,
			BankAccountNumber:      emp.BankAccountNumber,
			Status:                 payroll.StatusDraft,
		}

		if err := s.repo.CreatePayslip(ctx, &slip); err != nil {
			return nil, fmt.Errorf("save payslip for %s: %w", emp.Name, err)
		}
		generated = append(generated, slip)
	}

	return generated, nil
}

// PostPayrollToLedger creates a balanced double-entry journal entry for the payroll period
func (s *Service) PostPayrollToLedger(ctx context.Context, month, year int, branchIDStr string, postedBy string) (*accountingDomain.JournalEntry, error) {
	slips, err := s.repo.ListPayslips(ctx, payroll.PayslipFilter{
		Month: month,
		Year:  year,
	})
	if err != nil {
		return nil, fmt.Errorf("list slips for ledger post: %w", err)
	}
	if len(slips) == 0 {
		return nil, fmt.Errorf("tidak ada data slip gaji untuk periode %d/%02d", year, month)
	}

	var totalGross, totalPph, totalTakeHome, totalBpjs decimal.Decimal
	for _, slip := range slips {
		totalGross = totalGross.Add(slip.TotalGrossSalary)
		totalPph = totalPph.Add(slip.PPh21)
		totalTakeHome = totalTakeHome.Add(slip.TakeHomePay)
		totalBpjs = totalBpjs.Add(slip.BPJSKesehatan).Add(slip.BPJSKetenagakerjaan).Add(slip.AttendanceDeduction).Add(slip.LoanDeduction).Add(slip.OtherDeduction)
	}

	branchID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if branchIDStr != "" {
		if parsed, err := uuid.Parse(branchIDStr); err == nil && parsed != uuid.Nil {
			branchID = parsed
		}
	}

	periodName := fmt.Sprintf("%s %d", getMonthNameID(month), year)
	now := time.Now()

	// Double-entry balancing lines:
	// Debit: 6-1210 Beban Gaji Karyawan, Perawat & Staf
	// Credit: 2-1510 Hutang Pajak PPh 21
	// Credit: 1-1114 Bank Mandiri Operasional & Payroll (Take Home Pay)
	// Credit: 2-1130 Hutang Gaji & Karyawan / BPJS
	lines := []accountingApp.CreateLineCommand{
		{
			AccountCode: "6-1210",
			Debit:       totalGross,
			Credit:      decimal.Zero,
			Description: fmt.Sprintf("Beban Gaji Karyawan & Staf Medis Periode %s", periodName),
		},
		{
			AccountCode: "2-1510",
			Debit:       decimal.Zero,
			Credit:      totalPph,
			Description: fmt.Sprintf("Potongan PPh 21 Karyawan & Tenaga Medis Periode %s", periodName),
		},
		{
			AccountCode: "1-1114",
			Debit:       decimal.Zero,
			Credit:      totalTakeHome,
			Description: fmt.Sprintf("Disbursement Mandiri Payroll (Take Home Pay) Periode %s", periodName),
		},
	}

	if totalBpjs.IsPositive() {
		lines = append(lines, accountingApp.CreateLineCommand{
			AccountCode: "2-1130",
			Debit:       decimal.Zero,
			Credit:      totalBpjs,
			Description: fmt.Sprintf("Hutang Iuran BPJS Ketenagakerjaan & Kesehatan Periode %s", periodName),
		})
	}

	var journalEntry *accountingDomain.JournalEntry
	err = s.uow.Do(ctx, func(txCtx context.Context) error {
		jEntry, err := s.accountingService.CreateJournal(txCtx, accountingApp.CreateJournalCommand{
			BranchID:        branchID,
			JournalType:     accountingDomain.JournalTypeGeneral,
			TransactionDate: now,
			ReferenceType:   "PAYROLL_BATCH",
			ReferenceID:     fmt.Sprintf("PAY-%04d%02d", year, month),
			Description:     fmt.Sprintf("Payroll Gaji Staf & Karyawan Semarang Medical Group Periode %s", periodName),
			AutoPost:        true,
			Lines:           lines,
		})
		if err != nil {
			return fmt.Errorf("create payroll journal: %w", err)
		}
		journalEntry = jEntry

		// Update all slips in period with the journal ID
		for _, slip := range slips {
			_ = s.repo.UpdatePayslipJournalID(txCtx, slip.ID, jEntry.ID.String())
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("post payroll to ledger uow: %w", err)
	}

	return journalEntry, nil
}
