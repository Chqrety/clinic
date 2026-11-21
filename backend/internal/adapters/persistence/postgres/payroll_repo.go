package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"clinic-api/internal/domain/payroll"
	"clinic-api/pkg/decimal"

	"github.com/google/uuid"
)

type PostgresPayrollRepo struct {
	db *sql.DB
}

func NewPostgresPayrollRepo(db *sql.DB) *PostgresPayrollRepo {
	return &PostgresPayrollRepo{db: db}
}

// ----------------------------------------------------------------------------
// Employee Operations
// ----------------------------------------------------------------------------

func (r *PostgresPayrollRepo) ListEmployees(ctx context.Context, businessUnit string, isActiveOnly bool) ([]payroll.Employee, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT id, nik, name, position, department, business_unit, employment_type,
		       join_date::text, base_salary, position_allowance, transport_meal_allowance,
		       bank_name, bank_account_number, npwp, ptkp_status, bpjs_kesehatan_no,
		       bpjs_ketenagakerjaan_no, is_active, created_at, updated_at
		FROM payroll_employees
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if businessUnit != "" && businessUnit != "ALL" {
		query += fmt.Sprintf(" AND (business_unit = $%d OR business_unit = 'ALL')", argIdx)
		args = append(args, businessUnit)
		argIdx++
	}

	if isActiveOnly {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, true)
		argIdx++
	}

	query += " ORDER BY business_unit, name ASC"

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payroll employees: %w", err)
	}
	defer rows.Close()

	var employees []payroll.Employee
	for rows.Next() {
		var emp payroll.Employee
		var baseSalary, posAllowance, transAllowance float64
		err := rows.Scan(
			&emp.ID, &emp.NIK, &emp.Name, &emp.Position, &emp.Department, &emp.BusinessUnit, &emp.EmploymentType,
			&emp.JoinDate, &baseSalary, &posAllowance, &transAllowance,
			&emp.BankName, &emp.BankAccountNumber, &emp.NPWP, &emp.PTKPStatus, &emp.BPJSKesehatanNo,
			&emp.BPJSKetenagakerjaanNo, &emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan payroll employee: %w", err)
		}
		emp.BaseSalary = decimal.FromFloat(baseSalary)
		emp.PositionAllowance = decimal.FromFloat(posAllowance)
		emp.TransportMealAllowance = decimal.FromFloat(transAllowance)
		employees = append(employees, emp)
	}

	return employees, nil
}

func (r *PostgresPayrollRepo) GetEmployeeByID(ctx context.Context, id string) (*payroll.Employee, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT id, nik, name, position, department, business_unit, employment_type,
		       join_date::text, base_salary, position_allowance, transport_meal_allowance,
		       bank_name, bank_account_number, npwp, ptkp_status, bpjs_kesehatan_no,
		       bpjs_ketenagakerjaan_no, is_active, created_at, updated_at
		FROM payroll_employees
		WHERE id = $1
	`
	var emp payroll.Employee
	var baseSalary, posAllowance, transAllowance float64
	err := exec.QueryRowContext(ctx, query, id).Scan(
		&emp.ID, &emp.NIK, &emp.Name, &emp.Position, &emp.Department, &emp.BusinessUnit, &emp.EmploymentType,
		&emp.JoinDate, &baseSalary, &posAllowance, &transAllowance,
		&emp.BankName, &emp.BankAccountNumber, &emp.NPWP, &emp.PTKPStatus, &emp.BPJSKesehatanNo,
		&emp.BPJSKetenagakerjaanNo, &emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get payroll employee: %w", err)
	}

	emp.BaseSalary = decimal.FromFloat(baseSalary)
	emp.PositionAllowance = decimal.FromFloat(posAllowance)
	emp.TransportMealAllowance = decimal.FromFloat(transAllowance)
	return &emp, nil
}

func (r *PostgresPayrollRepo) CreateEmployee(ctx context.Context, emp *payroll.Employee) error {
	exec := GetExecutor(ctx, r.db)
	if emp.ID == "" {
		emp.ID = fmt.Sprintf("emp-%s", uuid.New().String()[:8])
	}
	if emp.JoinDate == "" {
		emp.JoinDate = time.Now().Format("2006-01-02")
	}

	query := `
		INSERT INTO payroll_employees (
			id, nik, name, position, department, business_unit, employment_type,
			join_date, base_salary, position_allowance, transport_meal_allowance,
			bank_name, bank_account_number, npwp, ptkp_status, bpjs_kesehatan_no,
			bpjs_ketenagakerjaan_no, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8::date, $9, $10, $11,
			$12, $13, $14, $15, $16,
			$17, $18, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
	`
	_, err := exec.ExecContext(ctx, query,
		emp.ID, emp.NIK, emp.Name, emp.Position, emp.Department, emp.BusinessUnit, emp.EmploymentType,
		emp.JoinDate, emp.BaseSalary.InexactFloat64(), emp.PositionAllowance.InexactFloat64(), emp.TransportMealAllowance.InexactFloat64(),
		emp.BankName, emp.BankAccountNumber, emp.NPWP, emp.PTKPStatus, emp.BPJSKesehatanNo,
		emp.BPJSKetenagakerjaanNo, emp.IsActive,
	)
	if err != nil {
		return fmt.Errorf("insert payroll employee: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) UpdateEmployee(ctx context.Context, emp *payroll.Employee) error {
	exec := GetExecutor(ctx, r.db)

	query := `
		UPDATE payroll_employees SET
			name = $2, position = $3, department = $4, business_unit = $5,
			employment_type = $6, base_salary = $7, position_allowance = $8,
			transport_meal_allowance = $9, bank_name = $10, bank_account_number = $11,
			npwp = $12, ptkp_status = $13, bpjs_kesehatan_no = $14,
			bpjs_ketenagakerjaan_no = $15, is_active = $16, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := exec.ExecContext(ctx, query,
		emp.ID, emp.Name, emp.Position, emp.Department, emp.BusinessUnit,
		emp.EmploymentType, emp.BaseSalary.InexactFloat64(), emp.PositionAllowance.InexactFloat64(),
		emp.TransportMealAllowance.InexactFloat64(), emp.BankName, emp.BankAccountNumber,
		emp.NPWP, emp.PTKPStatus, emp.BPJSKesehatanNo,
		emp.BPJSKetenagakerjaanNo, emp.IsActive,
	)
	if err != nil {
		return fmt.Errorf("update payroll employee: %w", err)
	}
	return nil
}

// ----------------------------------------------------------------------------
// Payslip Operations
// ----------------------------------------------------------------------------

func (r *PostgresPayrollRepo) ListPayslips(ctx context.Context, filter payroll.PayslipFilter) ([]payroll.Payslip, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT id, slip_number, employee_id, employee_nik, employee_name, position,
		       department, business_unit, period_month, period_year, period_name,
		       payment_date::text, working_days, present_days, sick_leave_days, paid_leave_days,
		       absent_days, overtime_hours, daily_meal_transport_rate,
		       base_salary, position_allowance, transport_meal_allowance,
		       shift_attendance_allowance, incentive_bonus, overtime_pay, other_allowance,
		       total_gross_salary, bpjs_kesehatan, bpjs_ketenagakerjaan, pph21,
		       attendance_deduction, loan_deduction, other_deduction, total_deduction,
		       take_home_pay, bank_name, bank_account_number, status, approved_by,
		       approved_at, paid_at, journal_id, notes, created_at, updated_at
		FROM payroll_slips
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if filter.Month > 0 {
		query += fmt.Sprintf(" AND period_month = $%d", argIdx)
		args = append(args, filter.Month)
		argIdx++
	}

	if filter.Year > 0 {
		query += fmt.Sprintf(" AND period_year = $%d", argIdx)
		args = append(args, filter.Year)
		argIdx++
	}

	if filter.BusinessUnit != "" && filter.BusinessUnit != "ALL" {
		query += fmt.Sprintf(" AND (business_unit = $%d OR business_unit = 'ALL')", argIdx)
		args = append(args, filter.BusinessUnit)
		argIdx++
	}

	if filter.Status != "" && filter.Status != "ALL" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}

	if filter.Search != "" {
		s := "%" + strings.ToLower(filter.Search) + "%"
		query += fmt.Sprintf(" AND (LOWER(employee_name) LIKE $%d OR LOWER(employee_nik) LIKE $%d OR LOWER(slip_number) LIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, s)
		argIdx++
	}

	query += " ORDER BY business_unit, employee_name ASC"

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payroll slips: %w", err)
	}
	defer rows.Close()

	var slips []payroll.Payslip
	for rows.Next() {
		var s payroll.Payslip
		var workDays, presDays, sickDays, paidLeaveDays, absDays int
		var otHours, dailyRate float64
		var baseSal, posAllow, transAllow, shiftAllow, incBonus, otPay, othAllow, gross float64
		var bpjsKes, bpjsTk, pph, attDed, loanDed, othDed, totDed, thp float64
		var approvedBy, journalID, notes sql.NullString
		var approvedAt, paidAt *time.Time

		err := rows.Scan(
			&s.ID, &s.SlipNumber, &s.EmployeeID, &s.EmployeeNIK, &s.EmployeeName, &s.Position,
			&s.Department, &s.BusinessUnit, &s.PeriodMonth, &s.PeriodYear, &s.PeriodName,
			&s.PaymentDate, &workDays, &presDays, &sickDays, &paidLeaveDays,
			&absDays, &otHours, &dailyRate,
			&baseSal, &posAllow, &transAllow,
			&shiftAllow, &incBonus, &otPay, &othAllow,
			&gross, &bpjsKes, &bpjsTk, &pph,
			&attDed, &loanDed, &othDed, &totDed,
			&thp, &s.BankName, &s.BankAccountNumber, &s.Status, &approvedBy,
			&approvedAt, &paidAt, &journalID, &notes, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan payroll slip: %w", err)
		}

		s.WorkingDays = workDays
		s.PresentDays = presDays
		s.SickLeaveDays = sickDays
		s.PaidLeaveDays = paidLeaveDays
		s.AbsentDays = absDays
		s.OvertimeHours = otHours
		s.DailyMealTransportRate = decimal.FromFloat(dailyRate)

		s.BaseSalary = decimal.FromFloat(baseSal)
		s.PositionAllowance = decimal.FromFloat(posAllow)
		s.TransportMealAllowance = decimal.FromFloat(transAllow)
		s.ShiftAttendanceAllowance = decimal.FromFloat(shiftAllow)
		s.IncentiveBonus = decimal.FromFloat(incBonus)
		s.OvertimePay = decimal.FromFloat(otPay)
		s.OtherAllowance = decimal.FromFloat(othAllow)
		s.TotalGrossSalary = decimal.FromFloat(gross)

		s.BPJSKesehatan = decimal.FromFloat(bpjsKes)
		s.BPJSKetenagakerjaan = decimal.FromFloat(bpjsTk)
		s.PPh21 = decimal.FromFloat(pph)
		s.AttendanceDeduction = decimal.FromFloat(attDed)
		s.LoanDeduction = decimal.FromFloat(loanDed)
		s.OtherDeduction = decimal.FromFloat(othDed)
		s.TotalDeduction = decimal.FromFloat(totDed)
		s.TakeHomePay = decimal.FromFloat(thp)

		if approvedBy.Valid {
			s.ApprovedBy = &approvedBy.String
		}
		s.ApprovedAt = approvedAt
		s.PaidAt = paidAt
		if journalID.Valid {
			s.JournalID = &journalID.String
		}
		if notes.Valid {
			s.Notes = notes.String
		}

		slips = append(slips, s)
	}

	return slips, nil
}

func (r *PostgresPayrollRepo) GetPayslipByID(ctx context.Context, id string) (*payroll.Payslip, error) {
	exec := GetExecutor(ctx, r.db)

	query := `
		SELECT id, slip_number, employee_id, employee_nik, employee_name, position,
		       department, business_unit, period_month, period_year, period_name,
		       payment_date::text, working_days, present_days, sick_leave_days, paid_leave_days,
		       absent_days, overtime_hours, daily_meal_transport_rate,
		       base_salary, position_allowance, transport_meal_allowance,
		       shift_attendance_allowance, incentive_bonus, overtime_pay, other_allowance,
		       total_gross_salary, bpjs_kesehatan, bpjs_ketenagakerjaan, pph21,
		       attendance_deduction, loan_deduction, other_deduction, total_deduction,
		       take_home_pay, bank_name, bank_account_number, status, approved_by,
		       approved_at, paid_at, journal_id, notes, created_at, updated_at
		FROM payroll_slips
		WHERE id = $1
	`
	var s payroll.Payslip
	var workDays, presDays, sickDays, paidLeaveDays, absDays int
	var otHours, dailyRate float64
	var baseSal, posAllow, transAllow, shiftAllow, incBonus, otPay, othAllow, gross float64
	var bpjsKes, bpjsTk, pph, attDed, loanDed, othDed, totDed, thp float64
	var approvedBy, journalID, notes sql.NullString
	var approvedAt, paidAt *time.Time

	err := exec.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.SlipNumber, &s.EmployeeID, &s.EmployeeNIK, &s.EmployeeName, &s.Position,
		&s.Department, &s.BusinessUnit, &s.PeriodMonth, &s.PeriodYear, &s.PeriodName,
		&s.PaymentDate, &workDays, &presDays, &sickDays, &paidLeaveDays,
		&absDays, &otHours, &dailyRate,
		&baseSal, &posAllow, &transAllow,
		&shiftAllow, &incBonus, &otPay, &othAllow,
		&gross, &bpjsKes, &bpjsTk, &pph,
		&attDed, &loanDed, &othDed, &totDed,
		&thp, &s.BankName, &s.BankAccountNumber, &s.Status, &approvedBy,
		&approvedAt, &paidAt, &journalID, &notes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get payroll slip: %w", err)
	}

	s.WorkingDays = workDays
	s.PresentDays = presDays
	s.SickLeaveDays = sickDays
	s.PaidLeaveDays = paidLeaveDays
	s.AbsentDays = absDays
	s.OvertimeHours = otHours
	s.DailyMealTransportRate = decimal.FromFloat(dailyRate)

	s.BaseSalary = decimal.FromFloat(baseSal)
	s.PositionAllowance = decimal.FromFloat(posAllow)
	s.TransportMealAllowance = decimal.FromFloat(transAllow)
	s.ShiftAttendanceAllowance = decimal.FromFloat(shiftAllow)
	s.IncentiveBonus = decimal.FromFloat(incBonus)
	s.OvertimePay = decimal.FromFloat(otPay)
	s.OtherAllowance = decimal.FromFloat(othAllow)
	s.TotalGrossSalary = decimal.FromFloat(gross)

	s.BPJSKesehatan = decimal.FromFloat(bpjsKes)
	s.BPJSKetenagakerjaan = decimal.FromFloat(bpjsTk)
	s.PPh21 = decimal.FromFloat(pph)
	s.AttendanceDeduction = decimal.FromFloat(attDed)
	s.LoanDeduction = decimal.FromFloat(loanDed)
	s.OtherDeduction = decimal.FromFloat(othDed)
	s.TotalDeduction = decimal.FromFloat(totDed)
	s.TakeHomePay = decimal.FromFloat(thp)

	if approvedBy.Valid {
		s.ApprovedBy = &approvedBy.String
	}
	s.ApprovedAt = approvedAt
	s.PaidAt = paidAt
	if journalID.Valid {
		s.JournalID = &journalID.String
	}
	if notes.Valid {
		s.Notes = notes.String
	}

	return &s, nil
}

func (r *PostgresPayrollRepo) CreatePayslip(ctx context.Context, slip *payroll.Payslip) error {
	exec := GetExecutor(ctx, r.db)
	if slip.ID == "" {
		slip.ID = fmt.Sprintf("slip-%s", uuid.New().String()[:8])
	}
	if slip.SlipNumber == "" {
		slip.SlipNumber = fmt.Sprintf("SLIP/%04d/%02d/%s", slip.PeriodYear, slip.PeriodMonth, strings.ToUpper(uuid.New().String()[:6]))
	}
	if slip.PaymentDate == "" {
		slip.PaymentDate = time.Now().Format("2006-01-02")
	}

	query := `
		INSERT INTO payroll_slips (
			id, slip_number, employee_id, employee_nik, employee_name, position,
			department, business_unit, period_month, period_year, period_name,
			payment_date, working_days, present_days, sick_leave_days, paid_leave_days,
			absent_days, overtime_hours, daily_meal_transport_rate,
			base_salary, position_allowance, transport_meal_allowance,
			shift_attendance_allowance, incentive_bonus, overtime_pay, other_allowance,
			total_gross_salary, bpjs_kesehatan, bpjs_ketenagakerjaan, pph21,
			attendance_deduction, loan_deduction, other_deduction, total_deduction,
			take_home_pay, bank_name, bank_account_number, status, approved_by,
			approved_at, paid_at, journal_id, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12::date, $13, $14, $15, $16,
			$17, $18, $19,
			$20, $21, $22,
			$23, $24, $25, $26,
			$27, $28, $29, $30,
			$31, $32, $33, $34,
			$35, $36, $37, $38, $39,
			$40, $41, $42, $43, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		ON CONFLICT (slip_number) DO UPDATE SET
			present_days = EXCLUDED.present_days,
			absent_days = EXCLUDED.absent_days,
			overtime_hours = EXCLUDED.overtime_hours,
			base_salary = EXCLUDED.base_salary,
			transport_meal_allowance = EXCLUDED.transport_meal_allowance,
			total_gross_salary = EXCLUDED.total_gross_salary,
			total_deduction = EXCLUDED.total_deduction,
			take_home_pay = EXCLUDED.take_home_pay,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := exec.ExecContext(ctx, query,
		slip.ID, slip.SlipNumber, slip.EmployeeID, slip.EmployeeNIK, slip.EmployeeName, slip.Position,
		slip.Department, slip.BusinessUnit, slip.PeriodMonth, slip.PeriodYear, slip.PeriodName,
		slip.PaymentDate, slip.WorkingDays, slip.PresentDays, slip.SickLeaveDays, slip.PaidLeaveDays,
		slip.AbsentDays, slip.OvertimeHours, slip.DailyMealTransportRate.InexactFloat64(),
		slip.BaseSalary.InexactFloat64(), slip.PositionAllowance.InexactFloat64(), slip.TransportMealAllowance.InexactFloat64(),
		slip.ShiftAttendanceAllowance.InexactFloat64(), slip.IncentiveBonus.InexactFloat64(), slip.OvertimePay.InexactFloat64(), slip.OtherAllowance.InexactFloat64(),
		slip.TotalGrossSalary.InexactFloat64(), slip.BPJSKesehatan.InexactFloat64(), slip.BPJSKetenagakerjaan.InexactFloat64(), slip.PPh21.InexactFloat64(),
		slip.AttendanceDeduction.InexactFloat64(), slip.LoanDeduction.InexactFloat64(), slip.OtherDeduction.InexactFloat64(), slip.TotalDeduction.InexactFloat64(),
		slip.TakeHomePay.InexactFloat64(), slip.BankName, slip.BankAccountNumber, string(slip.Status), slip.ApprovedBy,
		slip.ApprovedAt, slip.PaidAt, slip.JournalID, slip.Notes,
	)
	if err != nil {
		return fmt.Errorf("insert payroll slip: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) UpdatePayslip(ctx context.Context, slip *payroll.Payslip) error {
	exec := GetExecutor(ctx, r.db)

	query := `
		UPDATE payroll_slips SET
			working_days = $2, present_days = $3, sick_leave_days = $4, paid_leave_days = $5,
			absent_days = $6, overtime_hours = $7, daily_meal_transport_rate = $8,
			base_salary = $9, position_allowance = $10, transport_meal_allowance = $11,
			shift_attendance_allowance = $12, incentive_bonus = $13, overtime_pay = $14, other_allowance = $15,
			total_gross_salary = $16, bpjs_kesehatan = $17, bpjs_ketenagakerjaan = $18, pph21 = $19,
			attendance_deduction = $20, loan_deduction = $21, other_deduction = $22, total_deduction = $23,
			take_home_pay = $24, bank_name = $25, bank_account_number = $26, notes = $27, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := exec.ExecContext(ctx, query,
		slip.ID, slip.WorkingDays, slip.PresentDays, slip.SickLeaveDays, slip.PaidLeaveDays,
		slip.AbsentDays, slip.OvertimeHours, slip.DailyMealTransportRate.InexactFloat64(),
		slip.BaseSalary.InexactFloat64(), slip.PositionAllowance.InexactFloat64(), slip.TransportMealAllowance.InexactFloat64(),
		slip.ShiftAttendanceAllowance.InexactFloat64(), slip.IncentiveBonus.InexactFloat64(), slip.OvertimePay.InexactFloat64(), slip.OtherAllowance.InexactFloat64(),
		slip.TotalGrossSalary.InexactFloat64(), slip.BPJSKesehatan.InexactFloat64(), slip.BPJSKetenagakerjaan.InexactFloat64(), slip.PPh21.InexactFloat64(),
		slip.AttendanceDeduction.InexactFloat64(), slip.LoanDeduction.InexactFloat64(), slip.OtherDeduction.InexactFloat64(), slip.TotalDeduction.InexactFloat64(),
		slip.TakeHomePay.InexactFloat64(), slip.BankName, slip.BankAccountNumber, slip.Notes,
	)
	if err != nil {
		return fmt.Errorf("update payroll slip: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) DeletePayslip(ctx context.Context, id string) error {
	exec := GetExecutor(ctx, r.db)
	query := "DELETE FROM payroll_slips WHERE id = $1"
	res, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete payroll slip: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("slip tidak ditemukan")
	}
	return nil
}

func (r *PostgresPayrollRepo) ApprovePayslip(ctx context.Context, id string, approvedBy string) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payroll_slips SET
			status = 'APPROVED', approved_by = $2, approved_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'DRAFT'
	`
	_, err := exec.ExecContext(ctx, query, id, approvedBy)
	if err != nil {
		return fmt.Errorf("approve payroll slip: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) ApproveAllInPeriod(ctx context.Context, month, year int, approvedBy string) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payroll_slips SET
			status = 'APPROVED', approved_by = $3, approved_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE period_month = $1 AND period_year = $2 AND status = 'DRAFT'
	`
	_, err := exec.ExecContext(ctx, query, month, year, approvedBy)
	if err != nil {
		return fmt.Errorf("approve all payroll slips in period: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) MarkPayslipAsPaid(ctx context.Context, id string) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payroll_slips SET
			status = 'PAID', paid_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark payroll slip paid: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) MarkAllAsPaidInPeriod(ctx context.Context, month, year int) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payroll_slips SET
			status = 'PAID', paid_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE period_month = $1 AND period_year = $2
	`
	_, err := exec.ExecContext(ctx, query, month, year)
	if err != nil {
		return fmt.Errorf("mark all payroll slips paid in period: %w", err)
	}
	return nil
}

func (r *PostgresPayrollRepo) UpdatePayslipJournalID(ctx context.Context, id string, journalID string) error {
	exec := GetExecutor(ctx, r.db)
	query := `
		UPDATE payroll_slips SET
			journal_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := exec.ExecContext(ctx, query, id, journalID)
	if err != nil {
		return fmt.Errorf("update payroll slip journal id: %w", err)
	}
	return nil
}
