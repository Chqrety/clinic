-- Migration: Cash and Bank Management

CREATE TABLE IF NOT EXISTS cash_bank_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL, -- CASH, BANK
    bank_name VARCHAR(100),
    account_number VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cash_bank_accounts_branch ON cash_bank_accounts(branch_id);
CREATE INDEX idx_cash_bank_accounts_type ON cash_bank_accounts(type);

CREATE TABLE IF NOT EXISTS cash_bank_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cash_bank_account_id UUID NOT NULL REFERENCES cash_bank_accounts(id) ON DELETE RESTRICT,
    transaction_number VARCHAR(100) NOT NULL UNIQUE,
    transaction_type VARCHAR(50) NOT NULL, -- CASH_IN, CASH_OUT, TRANSFER, BANK_FEE
    amount NUMERIC(18, 2) NOT NULL,
    transaction_date DATE NOT NULL,
    counterpart_account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    description TEXT NOT NULL,
    reference_no VARCHAR(100),
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cash_bank_tx_account ON cash_bank_transactions(cash_bank_account_id);
CREATE INDEX idx_cash_bank_tx_date ON cash_bank_transactions(transaction_date);
CREATE INDEX idx_cash_bank_tx_journal ON cash_bank_transactions(journal_id);

