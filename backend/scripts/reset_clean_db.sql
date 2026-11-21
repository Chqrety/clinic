-- =========================================================================
-- CLEAN DATABASE SCRIPT (ZERO STATE - TRANSACTIONS & DUMMY DATA PURGE)
-- =========================================================================
-- Skrip ini mengosongkan seluruh data transaksi, antrian, rekam medis,
-- obat dummy, dan jurnal akuntansi, namun TETAP MEMPERTAHANKAN:
-- 1. accounts (Chart of Accounts / COA)
-- 2. branches (Cabang / Unit Bisnis: Pusat, Kidinest, SSC)
-- 3. staff_users (Akun Staf & Dokter untuk keperluan login)
-- 4. cost_centers & accounting_policies (Struktur dasar akuntansi)
-- =========================================================================

BEGIN;

TRUNCATE TABLE 
    -- Jurnal & Transaksi Akuntansi
    journal_lines,
    journal_entries,
    cash_bank_transactions,
    tax_withholdings,
    receivable_payments,
    receivables,
    payable_payments,
    payables,
    fiscal_periods,
    budgets,
    budget_items,
    fixed_assets,
    depreciation_schedules,
    cash_bank_accounts,
    
    -- Operasional Klinik, Kasir, Resep & Antrian Pasien
    invoice_items,
    invoices,
    cashier_shift_reports,
    prescription_items,
    prescriptions,
    consultations,
    queues,
    daily_queue_counters,
    satusehat_logs,
    patients,
    
    -- Stok & Master Obat Dummy
    stock_opname_logs,
    drug_batches,
    drugs,
    
    -- Jadwal Dokter Demo
    doctor_schedule_overrides,
    doctor_schedules
RESTART IDENTITY CASCADE;

COMMIT;
