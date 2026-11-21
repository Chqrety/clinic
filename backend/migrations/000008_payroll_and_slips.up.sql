-- Migration: 000008_payroll_and_slips.up.sql
-- Table: payroll_employees (Master Data Karyawan & Staf Klinik)
CREATE TABLE IF NOT EXISTS payroll_employees (
    id VARCHAR(50) PRIMARY KEY,
    nik VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    position VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    business_unit VARCHAR(20) NOT NULL CHECK (business_unit IN ('KIDINEST', 'SSC', 'ALL', 'MINIMO', 'APOTEK', 'SMC', 'SMG')),
    employment_type VARCHAR(20) NOT NULL CHECK (employment_type IN ('TETAP', 'KONTRAK', 'HONORER', 'MAGANG')),
    join_date DATE NOT NULL DEFAULT CURRENT_DATE,
    base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    position_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    transport_meal_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    bank_name VARCHAR(50) NOT NULL DEFAULT 'Bank Mandiri',
    bank_account_number VARCHAR(50) NOT NULL DEFAULT '',
    npwp VARCHAR(50) NOT NULL DEFAULT '',
    ptkp_status VARCHAR(10) NOT NULL DEFAULT 'TK/0',
    bpjs_kesehatan_no VARCHAR(50) NOT NULL DEFAULT '',
    bpjs_ketenagakerjaan_no VARCHAR(50) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payroll_emp_unit ON payroll_employees(business_unit);
CREATE INDEX IF NOT EXISTS idx_payroll_emp_active ON payroll_employees(is_active);

-- Table: payroll_slips (Transaksi Slip Gaji Bulanan Karyawan)
CREATE TABLE IF NOT EXISTS payroll_slips (
    id VARCHAR(50) PRIMARY KEY,
    slip_number VARCHAR(50) UNIQUE NOT NULL,
    employee_id VARCHAR(50) NOT NULL REFERENCES payroll_employees(id) ON DELETE RESTRICT,
    employee_nik VARCHAR(50) NOT NULL,
    employee_name VARCHAR(150) NOT NULL,
    position VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    business_unit VARCHAR(20) NOT NULL,
    period_month INT NOT NULL CHECK (period_month BETWEEN 1 AND 12),
    period_year INT NOT NULL CHECK (period_year >= 2020),
    period_name VARCHAR(50) NOT NULL,
    payment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    
    -- Kehadiran / Hari Masuk (Attendance)
    working_days INT NOT NULL DEFAULT 22,
    present_days INT NOT NULL DEFAULT 22,
    sick_leave_days INT NOT NULL DEFAULT 0,
    paid_leave_days INT NOT NULL DEFAULT 0,
    absent_days INT NOT NULL DEFAULT 0,
    overtime_hours NUMERIC(5,2) NOT NULL DEFAULT 0,
    daily_meal_transport_rate NUMERIC(15,2) NOT NULL DEFAULT 30000,

    -- Earnings (Penerimaan)
    base_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    position_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    transport_meal_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    shift_attendance_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    incentive_bonus NUMERIC(15,2) NOT NULL DEFAULT 0,
    overtime_pay NUMERIC(15,2) NOT NULL DEFAULT 0,
    other_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_gross_salary NUMERIC(15,2) NOT NULL DEFAULT 0,
    
    -- Deductions (Potongan)
    bpjs_kesehatan NUMERIC(15,2) NOT NULL DEFAULT 0,
    bpjs_ketenagakerjaan NUMERIC(15,2) NOT NULL DEFAULT 0,
    pph21 NUMERIC(15,2) NOT NULL DEFAULT 0,
    attendance_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
    loan_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
    other_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
    total_deduction NUMERIC(15,2) NOT NULL DEFAULT 0,
    
    -- Net Take Home Pay
    take_home_pay NUMERIC(15,2) NOT NULL DEFAULT 0,
    
    -- Payment & Audit
    bank_name VARCHAR(50) NOT NULL DEFAULT 'Bank Mandiri',
    bank_account_number VARCHAR(50) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'APPROVED', 'PAID')),
    approved_by VARCHAR(100),
    approved_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    journal_id VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_payroll_slips_period ON payroll_slips(period_year, period_month);
CREATE INDEX IF NOT EXISTS idx_payroll_slips_emp ON payroll_slips(employee_id);
CREATE INDEX IF NOT EXISTS idx_payroll_slips_unit ON payroll_slips(business_unit);
CREATE INDEX IF NOT EXISTS idx_payroll_slips_status ON payroll_slips(status);
