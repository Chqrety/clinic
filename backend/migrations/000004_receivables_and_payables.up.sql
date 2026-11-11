-- Migration: Account Receivable (Piutang) and Account Payable (Hutang)

-- 1. Receivables (Piutang Pasien, Asuransi, Korporat, BPJS)
CREATE TABLE IF NOT EXISTS receivables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    receivable_number VARCHAR(100) NOT NULL UNIQUE,
    customer_type VARCHAR(50) NOT NULL, -- PATIENT, INSURANCE, CORPORATE, BPJS
    customer_id UUID,
    customer_name VARCHAR(255) NOT NULL,
    invoice_reference VARCHAR(100),
    total_amount NUMERIC(18, 2) NOT NULL,
    outstanding_amount NUMERIC(18, 2) NOT NULL,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN', -- OPEN, PARTIAL, PAID, OVERDUE, WRITE_OFF
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_receivables_branch ON receivables(branch_id);
CREATE INDEX idx_receivables_status ON receivables(status);
CREATE INDEX idx_receivables_due ON receivables(due_date);

-- 2. Receivable Payments
CREATE TABLE IF NOT EXISTS receivable_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    receivable_id UUID NOT NULL REFERENCES receivables(id) ON DELETE RESTRICT,
    payment_number VARCHAR(100) NOT NULL UNIQUE,
    amount NUMERIC(18, 2) NOT NULL,
    payment_date DATE NOT NULL,
    cash_bank_account_code VARCHAR(50) NOT NULL, -- 1-1111, 1-1113
    reference_no VARCHAR(100),
    notes TEXT,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_receivable_pay_rec ON receivable_payments(receivable_id);
CREATE INDEX idx_receivable_pay_journal ON receivable_payments(journal_id);

-- 3. Payables (Hutang Supplier Obat, Alat Medis, Fee Dokter)
CREATE TABLE IF NOT EXISTS payables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    payable_number VARCHAR(100) NOT NULL UNIQUE,
    vendor_type VARCHAR(50) NOT NULL, -- SUPPLIER_MEDICINE, SUPPLIER_EQUIPMENT, DOCTOR, OPERATIONAL
    vendor_id UUID,
    vendor_name VARCHAR(255) NOT NULL,
    invoice_reference VARCHAR(100),
    total_amount NUMERIC(18, 2) NOT NULL,
    outstanding_amount NUMERIC(18, 2) NOT NULL,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN', -- OPEN, PARTIAL, PAID, OVERDUE
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payables_branch ON payables(branch_id);
CREATE INDEX idx_payables_status ON payables(status);
CREATE INDEX idx_payables_due ON payables(due_date);

-- 4. Payable Payments
CREATE TABLE IF NOT EXISTS payable_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payable_id UUID NOT NULL REFERENCES payables(id) ON DELETE RESTRICT,
    payment_number VARCHAR(100) NOT NULL UNIQUE,
    amount NUMERIC(18, 2) NOT NULL,
    payment_date DATE NOT NULL,
    cash_bank_account_code VARCHAR(50) NOT NULL, -- 1-1113 Bank BCA
    reference_no VARCHAR(100),
    notes TEXT,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payable_pay_payable ON payable_payments(payable_id);
CREATE INDEX idx_payable_pay_journal ON payable_payments(journal_id);

