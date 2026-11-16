-- Migration: Billing, Invoices, and Payment Records

CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    invoice_number VARCHAR(100) NOT NULL UNIQUE,
    patient_id UUID,
    patient_name VARCHAR(255) NOT NULL,
    patient_nik VARCHAR(50),
    queue_id BIGINT,
    invoice_type VARCHAR(50) NOT NULL, -- REGISTRATION, CONSULTATION, PHARMACY, COMBINED
    registration_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    consultation_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    medication_fee NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    total_amount NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    paid_amount NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    payment_status VARCHAR(50) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PARTIAL, PAID, CANCELLED
    payment_method VARCHAR(50), -- CASH, BANK_TRANSFER, RECEIVABLE
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoices_branch ON invoices(branch_id);
CREATE INDEX idx_invoices_status ON invoices(payment_status);
CREATE INDEX idx_invoices_queue ON invoices(queue_id);
CREATE INDEX idx_invoices_journal ON invoices(journal_id);

CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    item_type VARCHAR(50) NOT NULL, -- REGISTRATION, SERVICE, MEDICINE, PROCEDURE
    item_id UUID,
    description TEXT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    total_price NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_items_invoice ON invoice_items(invoice_id);

CREATE TABLE IF NOT EXISTS payment_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    payment_number VARCHAR(100) NOT NULL UNIQUE,
    amount NUMERIC(18, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- CASH, BANK_TRANSFER, RECEIVABLE
    reference_no VARCHAR(100),
    notes TEXT,
    received_by UUID,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    paid_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_records_invoice ON payment_records(invoice_id);
CREATE INDEX idx_payment_records_journal ON payment_records(journal_id);

