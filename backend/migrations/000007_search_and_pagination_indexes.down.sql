-- 000007_search_and_pagination_indexes.down.sql
DROP INDEX IF EXISTS idx_accounts_type_active;
DROP INDEX IF EXISTS idx_accounts_level;
DROP INDEX IF EXISTS idx_accounts_lower_name;
DROP INDEX IF EXISTS idx_accounts_code_prefix;

DROP INDEX IF EXISTS idx_journal_entries_branch_date_id;
DROP INDEX IF EXISTS idx_journal_entries_status_type;
DROP INDEX IF EXISTS idx_journal_entries_lower_number;
DROP INDEX IF EXISTS idx_journal_entries_ref;

DROP INDEX IF EXISTS idx_invoices_branch_status_created;
DROP INDEX IF EXISTS idx_invoices_lower_patient_name;
DROP INDEX IF EXISTS idx_invoices_lower_number;
DROP INDEX IF EXISTS idx_invoices_type_status;

DROP INDEX IF EXISTS idx_cash_bank_tx_branch_date_id;
DROP INDEX IF EXISTS idx_cash_bank_tx_account_type;
DROP INDEX IF EXISTS idx_cash_bank_tx_lower_number;

DROP INDEX IF EXISTS idx_receivables_branch_status_due;
DROP INDEX IF EXISTS idx_receivables_lower_customer;

DROP INDEX IF EXISTS idx_payables_branch_status_due;
DROP INDEX IF EXISTS idx_payables_lower_vendor;

DROP INDEX IF EXISTS idx_doctor_share_rules_branch_active;
DROP INDEX IF EXISTS idx_doctor_settlements_branch_period;
DROP INDEX IF EXISTS idx_doctor_settlements_lower_doctor;

DROP INDEX IF EXISTS idx_cashier_shifts_branch_status_started;

DROP INDEX IF EXISTS idx_tax_withholdings_branch_period;
DROP INDEX IF EXISTS idx_tax_withholdings_lower_taxpayer;

DROP INDEX IF EXISTS idx_fixed_assets_branch_cat_status;
DROP INDEX IF EXISTS idx_fixed_assets_lower_name;

DROP INDEX IF EXISTS idx_prepaid_expenses_branch_status;
DROP INDEX IF EXISTS idx_prepaid_expenses_lower_name;

DROP INDEX IF EXISTS idx_budgets_branch_year_center;

DROP INDEX IF EXISTS idx_audit_logs_branch_created_id;
DROP INDEX IF EXISTS idx_audit_logs_entity_action;
