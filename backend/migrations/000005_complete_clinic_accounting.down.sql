-- Down Migration 000005
DROP TABLE IF EXISTS financial_audit_logs;
DROP TABLE IF EXISTS fixed_asset_depreciations;
DROP TABLE IF EXISTS fixed_assets;
DROP TABLE IF EXISTS cashier_shift_denominations;
DROP TABLE IF EXISTS cashier_shifts;
DROP TABLE IF EXISTS tax_withholdings;
DROP TABLE IF EXISTS doctor_fee_items;
DROP TABLE IF EXISTS doctor_fee_settlements;
DROP TABLE IF EXISTS doctor_share_rules;
DROP TABLE IF EXISTS accounting_policies;
ALTER TABLE journal_lines DROP COLUMN IF EXISTS cost_center_id;
DROP TABLE IF EXISTS cost_centers;

