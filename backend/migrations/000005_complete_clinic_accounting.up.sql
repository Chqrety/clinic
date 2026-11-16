-- Migration 000005: Complete Clinic ERP Accounting Engine
-- 1. Cost Centers (Kidinest, Semarang Skin Center, Shared Back-Office)
-- 2. Accounting Policies & Dimension Linking
-- 3. Doctor Sharing & Therapist Commission Engine
-- 4. Tax Management & e-Bupot 21/26 Engine
-- 5. Cashier Shift & Cash Control Engine
-- 6. Fixed Asset & Depreciation Engine
-- 7. Financial Audit Logs

-- 1. Cost Centers Table
CREATE TABLE IF NOT EXISTS cost_centers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE, -- 'UNIT-KIDINEST', 'UNIT-SSC', 'SHARED-BACKOFFICE'
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed Default Cost Centers
INSERT INTO cost_centers (code, name, description)
VALUES 
    ('UNIT-KIDINEST', 'Unit Poli Anak & Tumbuh Kembang Kidinest', 'Pusat Biaya Layanan Pediatrik, Vaksinasi & Farmasi Anak'),
    ('UNIT-SSC', 'Unit Semarang Skin Center (SSC)', 'Pusat Biaya Layanan Estetika, Kulit, Facial & Produk Skincare'),
    ('SHARED-BACKOFFICE', 'Shared Back-Office & Treasury', 'Pusat Biaya Kantor Pusat, Treasury Kas/Bank, & Administrasi Bersama')
ON CONFLICT (code) DO NOTHING;

-- 2. Link Cost Center to Journal Lines
ALTER TABLE journal_lines ADD COLUMN IF NOT EXISTS cost_center_id UUID REFERENCES cost_centers(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_journal_lines_cost_center ON journal_lines(cost_center_id);

-- 3. Accounting Policies Table
CREATE TABLE IF NOT EXISTS accounting_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    policy_key VARCHAR(100) NOT NULL UNIQUE,
    policy_value TEXT NOT NULL,
    data_type VARCHAR(20) NOT NULL DEFAULT 'STRING', -- 'STRING', 'BOOLEAN', 'NUMBER', 'JSON'
    description TEXT,
    effective_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Seed Core Accounting Policies
INSERT INTO accounting_policies (policy_key, policy_value, data_type, description)
VALUES 
    ('ACCOUNTING_BASIS', 'ACCRUAL', 'STRING', 'Basis akuntansi klinik: Murni Akrual'),
    ('INVENTORY_VALUATION_METHOD', 'FIFO', 'STRING', 'Metode valuasi persediaan obat & skincare: First In First Out'),
    ('DEPRECIATION_METHOD', 'STRAIGHT_LINE', 'STRING', 'Metode depresiasi aset tetap: Garis Lurus'),
    ('IS_PKP_TAXABLE', 'false', 'BOOLEAN', 'Status Pengusaha Kena Pajak klinik (Default false)'),
    ('CURRENCY', 'IDR', 'STRING', 'Mata uang fungsional pelaporan akuntansi klinik: Indonesian Rupiah (IDR)'),
    ('DEFAULT_DOCTOR_SHARE_PERCENT', '70.00', 'NUMBER', 'Persentase baku bagi hasil dokter spesialis anak (70%)'),
    ('RETAINED_EARNINGS_ACCOUNT', '3-1200', 'STRING', 'Akun penampung Laba Ditahan'),
    ('CURRENT_EARNINGS_ACCOUNT', '3-1300', 'STRING', 'Akun penampung Laba Bersih Periode Berjalan')
ON CONFLICT (policy_key) DO NOTHING;

-- 4. Doctor Share Rules & Monthly Settlements
CREATE TABLE IF NOT EXISTS doctor_share_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cost_center_id UUID REFERENCES cost_centers(id) ON DELETE RESTRICT,
    doctor_id UUID,
    doctor_name VARCHAR(255) NOT NULL,
    service_category VARCHAR(100) NOT NULL, -- 'CONSULTATION', 'VACCINE_INJECTION', 'GROWTH_PACKAGE', 'FACIAL_TREATMENT', 'LASER_PEELING'
    sharing_percentage NUMERIC(5, 2) NOT NULL DEFAULT 70.00,
    fixed_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS doctor_fee_settlements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    doctor_id UUID,
    doctor_name VARCHAR(255) NOT NULL,
    doctor_npwp VARCHAR(50),
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    gross_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    dpp_amount NUMERIC(18, 2) NOT NULL DEFAULT 0.00, -- 50% PMK 168
    tax_pph21 NUMERIC(18, 2) NOT NULL DEFAULT 0.00,  -- Progresif Pasal 17
    net_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT',     -- 'DRAFT', 'APPROVED', 'PAID'
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS doctor_fee_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    settlement_id UUID REFERENCES doctor_fee_settlements(id) ON DELETE CASCADE,
    invoice_id UUID REFERENCES invoices(id) ON DELETE RESTRICT,
    service_date DATE NOT NULL,
    patient_name VARCHAR(255) NOT NULL,
    service_description VARCHAR(255) NOT NULL,
    patient_fee NUMERIC(18, 2) NOT NULL,
    doctor_share_amount NUMERIC(18, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Tax Withholdings (clinic_pajak / e-Bupot 21/26)
CREATE TABLE IF NOT EXISTS tax_withholdings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    tax_type VARCHAR(50) NOT NULL, -- 'PPH21_DOCTOR', 'PPH21_EMPLOYEE', 'PPH23', 'PPH4_2', 'PPN_OUTPUT'
    bupot_number VARCHAR(100) NOT NULL UNIQUE,
    taxpayer_nik VARCHAR(50),
    taxpayer_npwp VARCHAR(50),
    taxpayer_name VARCHAR(255) NOT NULL,
    object_code VARCHAR(50) NOT NULL, -- '21-100-07' for Doctor
    gross_amount NUMERIC(18, 2) NOT NULL,
    dpp_amount NUMERIC(18, 2) NOT NULL,
    tax_rate NUMERIC(5, 2) NOT NULL,
    tax_amount NUMERIC(18, 2) NOT NULL,
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    is_reported BOOLEAN NOT NULL DEFAULT FALSE,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. Cashier Shifts & Physical Cash Count
CREATE TABLE IF NOT EXISTS cashier_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cost_center_id UUID REFERENCES cost_centers(id) ON DELETE RESTRICT,
    cashier_user_id UUID,
    cashier_name VARCHAR(255) NOT NULL,
    cash_account_code VARCHAR(50) NOT NULL, -- '1-1111' Kidinest or '1-1112' SSC
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP WITH TIME ZONE,
    opening_cash NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    system_cash_in NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    system_cash_out NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    expected_cash NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    actual_cash NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    difference_amount NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    difference_type VARCHAR(20) NOT NULL DEFAULT 'NONE', -- 'NONE', 'OVERAGE', 'SHORTAGE'
    reconciliation_journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    deposit_journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN', -- 'OPEN', 'CLOSED', 'DEPOSITED'
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cashier_shift_denominations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    shift_id UUID NOT NULL REFERENCES cashier_shifts(id) ON DELETE CASCADE,
    denomination_value INT NOT NULL, -- 100000, 50000, 20000, 10000, 5000, 2000, 1000, 500
    quantity INT NOT NULL DEFAULT 0,
    subtotal NUMERIC(18, 2) NOT NULL DEFAULT 0.00
);

-- 7. Fixed Assets & Depreciation Schedules
CREATE TABLE IF NOT EXISTS fixed_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cost_center_id UUID REFERENCES cost_centers(id) ON DELETE RESTRICT,
    asset_code VARCHAR(100) NOT NULL UNIQUE,
    asset_name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'MEDICAL_EQUIPMENT', 'AESTHETIC_MACHINE', 'RENOVATION', 'OFFICE_IT'
    acquisition_date DATE NOT NULL,
    acquisition_cost NUMERIC(18, 2) NOT NULL,
    salvage_value NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    useful_life_months INT NOT NULL DEFAULT 48,
    depreciation_method VARCHAR(50) NOT NULL DEFAULT 'STRAIGHT_LINE',
    asset_account_code VARCHAR(50) NOT NULL,   -- '1-2100' or '1-2200'
    accum_account_code VARCHAR(50) NOT NULL DEFAULT '1-2900',
    expense_account_code VARCHAR(50) NOT NULL DEFAULT '6-1300',
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'FULLY_DEPRECIATED', 'DISPOSED'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fixed_asset_depreciations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    asset_id UUID NOT NULL REFERENCES fixed_assets(id) ON DELETE RESTRICT,
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    depreciation_amount NUMERIC(18, 2) NOT NULL,
    accumulated_amount NUMERIC(18, 2) NOT NULL,
    book_value NUMERIC(18, 2) NOT NULL,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 8. Financial Audit Logs Table
CREATE TABLE IF NOT EXISTS financial_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    user_id UUID,
    action VARCHAR(100) NOT NULL, -- 'JOURNAL_POST', 'JOURNAL_REVERSE', 'PERIOD_CLOSE', 'POLICY_CHANGE', 'SHIFT_CLOSE'
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    before_payload JSONB,
    after_payload JSONB,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_action ON financial_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON financial_audit_logs(entity_type, entity_id);

