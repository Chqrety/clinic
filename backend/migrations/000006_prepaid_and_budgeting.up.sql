-- Migration 000006: Prepaid Expenses, Amortization, and Budgeting

-- 1. Prepaid Expenses Table
CREATE TABLE IF NOT EXISTS prepaid_expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cost_center_id UUID REFERENCES cost_centers(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'RENT', 'INSURANCE', 'SOFTWARE_LICENSE', 'MAINTENANCE'
    total_amount NUMERIC(18, 2) NOT NULL,
    prepaid_account_code VARCHAR(50) NOT NULL DEFAULT '1-1400',
    expense_account_code VARCHAR(50) NOT NULL DEFAULT '6-1200',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    period_months INT NOT NULL,
    monthly_amount NUMERIC(18, 2) NOT NULL,
    remaining_amount NUMERIC(18, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- 'ACTIVE', 'FULLY_AMORTIZED', 'TERMINATED'
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_prepaid_branch ON prepaid_expenses(branch_id);
CREATE INDEX IF NOT EXISTS idx_prepaid_status ON prepaid_expenses(status);

-- 2. Prepaid Amortizations Table
CREATE TABLE IF NOT EXISTS prepaid_amortizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    prepaid_expense_id UUID NOT NULL REFERENCES prepaid_expenses(id) ON DELETE CASCADE,
    period_year INT NOT NULL,
    period_month INT NOT NULL,
    amount NUMERIC(18, 2) NOT NULL,
    journal_id UUID REFERENCES journal_entries(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_amortization_expense ON prepaid_amortizations(prepaid_expense_id);

-- 3. Budgets Table
CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cost_center_id UUID REFERENCES cost_centers(id) ON DELETE RESTRICT,
    period_year INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- 'DRAFT', 'APPROVED', 'REVISED'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_budgets_branch_year ON budgets(branch_id, period_year);

-- 4. Budget Items Table
CREATE TABLE IF NOT EXISTS budget_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    account_code VARCHAR(50) NOT NULL,
    period_month INT NOT NULL,
    budgeted_amount NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    notes TEXT
);

CREATE INDEX IF NOT EXISTS idx_budget_items_budget ON budget_items(budget_id);

