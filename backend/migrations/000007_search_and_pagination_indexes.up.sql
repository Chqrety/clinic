-- ==============================================================================
-- 000007_search_and_pagination_indexes.up.sql
-- Optimized composite B-Tree indexes for server-side pagination, sorting & search
-- ==============================================================================

-- 1. Accounts (Chart of Accounts)
CREATE INDEX IF NOT EXISTS idx_accounts_type_active ON accounts(type, is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_level ON accounts(level);
CREATE INDEX IF NOT EXISTS idx_accounts_lower_name ON accounts(LOWER(name));
CREATE INDEX IF NOT EXISTS idx_accounts_code_prefix ON accounts(code varchar_pattern_ops);

-- 2. Journal Entries & Ledgers (High Volume)
CREATE INDEX IF NOT EXISTS idx_journal_entries_branch_date_id ON journal_entries(branch_id, transaction_date DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_journal_entries_status_type ON journal_entries(status, journal_type);
CREATE INDEX IF NOT EXISTS idx_journal_entries_lower_number ON journal_entries(LOWER(journal_number));
CREATE INDEX IF NOT EXISTS idx_journal_entries_ref ON journal_entries(reference_type, reference_id);

-- 3. Invoices (Billing & POS)
CREATE INDEX IF NOT EXISTS idx_invoices_branch_status_created ON invoices(branch_id, payment_status, created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_invoices_lower_patient_name ON invoices(LOWER(patient_name));
CREATE INDEX IF NOT EXISTS idx_invoices_lower_number ON invoices(LOWER(invoice_number));
CREATE INDEX IF NOT EXISTS idx_invoices_type_status ON invoices(invoice_type, payment_status);

-- 4. Cash & Bank Transactions
CREATE INDEX IF NOT EXISTS idx_cash_bank_tx_branch_date_id ON cash_bank_transactions(branch_id, transaction_date DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_cash_bank_tx_account_type ON cash_bank_transactions(cash_bank_account_id, transaction_type);
CREATE INDEX IF NOT EXISTS idx_cash_bank_tx_lower_number ON cash_bank_transactions(LOWER(transaction_number));

-- 5. Receivables (AR)
CREATE INDEX IF NOT EXISTS idx_receivables_branch_status_due ON receivables(branch_id, status, due_date ASC, id ASC);
CREATE INDEX IF NOT EXISTS idx_receivables_lower_customer ON receivables(LOWER(customer_name));

-- 6. Payables (AP)
CREATE INDEX IF NOT EXISTS idx_payables_branch_status_due ON payables(branch_id, status, due_date ASC, id ASC);
CREATE INDEX IF NOT EXISTS idx_payables_lower_vendor ON payables(LOWER(vendor_name));

-- 7. Doctor Share Rules & Monthly Settlements
CREATE INDEX IF NOT EXISTS idx_doctor_share_rules_branch_active ON doctor_share_rules(branch_id, is_active, doctor_name);
CREATE INDEX IF NOT EXISTS idx_doctor_settlements_branch_period ON doctor_fee_settlements(branch_id, period_year, period_month, status);
CREATE INDEX IF NOT EXISTS idx_doctor_settlements_lower_doctor ON doctor_fee_settlements(LOWER(doctor_name));

-- 8. Treasury & Cashier Shifts
CREATE INDEX IF NOT EXISTS idx_cashier_shifts_branch_status_started ON cashier_shifts(branch_id, status, start_time DESC, id DESC);

-- 9. Tax Withholdings (e-Bupot)
CREATE INDEX IF NOT EXISTS idx_tax_withholdings_branch_period ON tax_withholdings(branch_id, period_year, period_month, tax_type);
CREATE INDEX IF NOT EXISTS idx_tax_withholdings_lower_taxpayer ON tax_withholdings(LOWER(taxpayer_name));

-- 10. Fixed Assets & Depreciation
CREATE INDEX IF NOT EXISTS idx_fixed_assets_branch_cat_status ON fixed_assets(branch_id, category, status, id DESC);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_lower_name ON fixed_assets(LOWER(asset_name));

-- 11. Prepaid Expenses
CREATE INDEX IF NOT EXISTS idx_prepaid_expenses_branch_status ON prepaid_expenses(branch_id, status, start_date);
CREATE INDEX IF NOT EXISTS idx_prepaid_expenses_lower_name ON prepaid_expenses(LOWER(name));

-- 12. Budgets
CREATE INDEX IF NOT EXISTS idx_budgets_branch_year_center ON budgets(branch_id, period_year, cost_center_id);

-- 13. Financial Audit Trails (Forensic Querying)
CREATE INDEX IF NOT EXISTS idx_audit_logs_branch_created_id ON financial_audit_logs(branch_id, created_at DESC, id DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_action ON financial_audit_logs(entity_type, action, created_at DESC);
