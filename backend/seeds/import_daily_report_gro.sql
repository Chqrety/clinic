-- ==============================================================================
-- REAL CLINIC ACCOUNTING TRANSACTIONS IMPORT (DAILY REPORT GRO)
-- Covers 2 Months: Oktober & November 2026 (+ Soft-Launch Akhir September)
-- Generates Balanced Double-Entry Journal Entries, Journal Lines & Cash/Bank Tx
-- ==============================================================================

BEGIN;

-- 1. Register Kliring Payment Gateway & EDC Account if not present

    INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
    VALUES ('c0000000-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 'GATEWAY-1120', 'Kliring Payment Gateway & EDC (QRIS/Tripay)', 'BANK', 'Settlement EDC & QRIS', 'EDC-QRIS-001', TRUE)
    ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, is_active = TRUE;
    
-- 2. Ensure Fiscal Periods exist for Sept, Oct, and Nov 2026

    INSERT INTO fiscal_periods (id, branch_id, year, month, status)
    VALUES
        ('f0000000-0000-0000-0000-000000000009', '11111111-1111-1111-1111-111111111111', 2026, 9, 'CLOSED'),
        ('f0000000-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111111', 2026, 10, 'CLOSED'),
        ('f0000000-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111111', 2026, 11, 'OPEN')
    ON CONFLICT (branch_id, year, month) DO UPDATE SET status = EXCLUDED.status;
    
-- 3. Clean up existing automated cashier closing journals to allow clean idempotent re-seed

    DELETE FROM cash_bank_transactions WHERE transaction_number LIKE 'BKM-KID-%' OR transaction_number LIKE 'BKM-SSC-%' OR transaction_number LIKE 'BKM-MNM-%';
    DELETE FROM journal_entries WHERE journal_number LIKE 'CR-KID-%' OR journal_number LIKE 'CR-SSC-%' OR journal_number LIKE 'CR-MNM-%';
    -- Also remove old arbitrary dummy receipts on Nov 13-14
    DELETE FROM journal_entries WHERE journal_number IN ('JV-2026-11-007', 'JV-2026-11-008', 'JV-2026-11-017', 'JV-2026-11-018', 'JV-2026-11-019', 'JV-2026-11-020', 'JV-2026-11-021');
    
-- 4. Insert Daily Cashier Revenue Journals and Lines from Daily Report GRO

    -- Journal CR-KID-20260924 (KIDINEST - 20260924): Rp 2,824,396.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('b2651754-7717-571e-ad2e-14ae8d1a90eb', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260924', 'CASH_RECEIPT', '2026-09-24', '2026-09-24 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260924', 2824396.00, 2824396.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 24-09-2026 Keysha Aurore Delika, Afkar Muhammad Alby, Hasseva Rayyan Pratyaksa, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4ec83b65-7400-53fe-aeda-f15c2b837bc1', 'b2651754-7717-571e-ad2e-14ae8d1a90eb', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2574396.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260924-11120', 'CASH_IN', 2574396.00, '2026-09-24', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 24/09/2026', 'GRO-KIDINEST-20260924', 'b2651754-7717-571e-ad2e-14ae8d1a90eb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9c0c15e4-68f7-5100-86f5-301f40c95bfd', 'b2651754-7717-571e-ad2e-14ae8d1a90eb', '2a3ad241-587d-4853-bdc2-615d72496b4e', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20260924-11111', 'CASH_IN', 250000.00, '2026-09-24', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 24/09/2026', 'GRO-KIDINEST-20260924', 'b2651754-7717-571e-ad2e-14ae8d1a90eb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c6856176-12ec-5369-a778-8d245774c236', 'b2651754-7717-571e-ad2e-14ae8d1a90eb', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9ae6a04c-ec46-5a1e-95ed-0eff0948dd98', 'b2651754-7717-571e-ad2e-14ae8d1a90eb', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 414396.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cd17c7d9-e4fc-5313-b05f-d9559dc53073', 'b2651754-7717-571e-ad2e-14ae8d1a90eb', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 910000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20260925 (KIDINEST - 20260925): Rp 1,750,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('55b3fad0-1919-568e-901b-da27d7933bd2', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260925', 'CASH_RECEIPT', '2026-09-25', '2026-09-25 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260925', 1750000.00, 1750000.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 25-09-2026 Alayya Rafa Azura, Ananta Faaz Auda, Bianca Gracella Handojo, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('59c6b318-b312-5c57-9a1a-f7b92f311bd3', '55b3fad0-1919-568e-901b-da27d7933bd2', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 500000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260925-11120', 'CASH_IN', 500000.00, '2026-09-25', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 25/09/2026', 'GRO-KIDINEST-20260925', '55b3fad0-1919-568e-901b-da27d7933bd2')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6f15a71b-a179-593b-bc3b-838bb66c1e3a', '55b3fad0-1919-568e-901b-da27d7933bd2', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20260925-11113', 'CASH_IN', 1250000.00, '2026-09-25', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 25/09/2026', 'GRO-KIDINEST-20260925', '55b3fad0-1919-568e-901b-da27d7933bd2')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ff99beaf-9945-5c1c-90c2-700ce428a634', '55b3fad0-1919-568e-901b-da27d7933bd2', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20260926 (KIDINEST - 20260926): Rp 10,044,112.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('aab6d936-63f4-5ace-aedd-3f08a9f8544e', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260926', 'CASH_RECEIPT', '2026-09-26', '2026-09-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260926', 10044112.00, 10044112.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 26-09-2026 Michaela Kalashnikova, Daphne Fumiko Wijaya, Hayyina Yafa Nasyitha, dkk. (total 17 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ed52bb4c-4e9a-50c1-a942-41d8d86ef5d7', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 8866112.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260926-11120', 'CASH_IN', 8866112.00, '2026-09-26', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 26/09/2026', 'GRO-KIDINEST-20260926', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4e96d877-ad19-5172-b477-5007f74daf53', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '2a3ad241-587d-4853-bdc2-615d72496b4e', 428000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20260926-11111', 'CASH_IN', 428000.00, '2026-09-26', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 26/09/2026', 'GRO-KIDINEST-20260926', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e25f3782-0531-58ff-900c-bfa44446edf8', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 750000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20260926-11113', 'CASH_IN', 750000.00, '2026-09-26', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 26/09/2026', 'GRO-KIDINEST-20260926', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1904df12-ea25-56c7-8a10-90891b491f15', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4125000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1858d536-b075-59cf-8705-6b4686b9df03', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3454112.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f7cef3fc-0d75-58d6-9a0f-695ca08a427b', 'aab6d936-63f4-5ace-aedd-3f08a9f8544e', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2465000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20260927 (KIDINEST - 20260927): Rp 500,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('bcc04240-e07c-582f-bc75-31beaa81684a', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260927', 'CASH_RECEIPT', '2026-09-27', '2026-09-27 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260927', 500000.00, 500000.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 27-09-2026 Andrew Elnathan Winoto, Audrey Eliana Winoto (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('44dab09e-2b07-5a51-92b1-5f4e4198ae03', 'bcc04240-e07c-582f-bc75-31beaa81684a', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 500000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260927-11120', 'CASH_IN', 500000.00, '2026-09-27', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 27/09/2026', 'GRO-KIDINEST-20260927', 'bcc04240-e07c-582f-bc75-31beaa81684a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('068d9520-b97e-59e8-88ca-e4fb3e3580da', 'bcc04240-e07c-582f-bc75-31beaa81684a', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20260928 (KIDINEST - 20260928): Rp 6,417,671.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c89d6fc3-709a-5229-9871-44484d88a487', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260928', 'CASH_RECEIPT', '2026-09-28', '2026-09-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260928', 6417671.00, 6417671.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 28-09-2026 Muhammad Syaghghaf Arham Elzayya, Sean Byantara Nursholeh, Raphael Gregorio Hadi Projo, dkk. (total 16 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('85002aff-cce0-56e8-938e-b5805c15cfcc', 'c89d6fc3-709a-5229-9871-44484d88a487', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5592584.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260928-11120', 'CASH_IN', 5592584.00, '2026-09-28', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 28/09/2026', 'GRO-KIDINEST-20260928', 'c89d6fc3-709a-5229-9871-44484d88a487')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5f7e3bf6-dccb-5a56-b972-a8c96bb76c95', 'c89d6fc3-709a-5229-9871-44484d88a487', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 683250.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20260928-11113', 'CASH_IN', 683250.00, '2026-09-28', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 28/09/2026', 'GRO-KIDINEST-20260928', 'c89d6fc3-709a-5229-9871-44484d88a487')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9faafbd4-cb03-5fe2-91dc-4a616f2e2d57', 'c89d6fc3-709a-5229-9871-44484d88a487', '2a3ad241-587d-4853-bdc2-615d72496b4e', 141837.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20260928-11111', 'CASH_IN', 141837.00, '2026-09-28', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 28/09/2026', 'GRO-KIDINEST-20260928', 'c89d6fc3-709a-5229-9871-44484d88a487')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('811e8d8b-1943-5a13-b78f-71cd0181ae60', 'c89d6fc3-709a-5229-9871-44484d88a487', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 3750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f516e910-8458-5f0c-8291-206247f8ad5d', 'c89d6fc3-709a-5229-9871-44484d88a487', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2667671.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20260929 (KIDINEST - 20260929): Rp 4,539,683.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('8ae94fd1-9d90-55ed-8db9-916f407e3826', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260929', 'CASH_RECEIPT', '2026-09-29', '2026-09-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260929', 4539683.00, 4539683.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 29-09-2026 Medin Kala Senja, Gwen Nona Kagumi, Chrysander Erlano Arsyanendra, dkk. (total 9 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('07cc52a8-446f-5537-a352-ce0de39dda4c', '8ae94fd1-9d90-55ed-8db9-916f407e3826', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2880674.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20260929-11113', 'CASH_IN', 2880674.00, '2026-09-29', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 29/09/2026', 'GRO-KIDINEST-20260929', '8ae94fd1-9d90-55ed-8db9-916f407e3826')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('417f02c6-6fee-5daa-80d7-791675528eba', '8ae94fd1-9d90-55ed-8db9-916f407e3826', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1659009.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260929-11120', 'CASH_IN', 1659009.00, '2026-09-29', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 29/09/2026', 'GRO-KIDINEST-20260929', '8ae94fd1-9d90-55ed-8db9-916f407e3826')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c0d66201-f040-529f-8c57-9db00e6505c7', '8ae94fd1-9d90-55ed-8db9-916f407e3826', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dd4301d8-1864-5181-99fe-cf6ca71a9343', '8ae94fd1-9d90-55ed-8db9-916f407e3826', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 979043.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6670c205-88b8-54c3-abc9-4a0e0971d502', '8ae94fd1-9d90-55ed-8db9-916f407e3826', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1310640.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20260930 (KIDINEST - 20260930): Rp 3,441,353.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ce47ec6a-5647-55e0-9151-696cc175d6cb', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20260930', 'CASH_RECEIPT', '2026-09-30', '2026-09-30 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20260930', 3441353.00, 3441353.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 30-09-2026 Mousa Arkan Dewantara, Adhiraja Radhitya Kalandra, Gabriella Evangeline Salim, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2bef0c41-3a98-5a61-bc17-50ffc5fc3a40', 'ce47ec6a-5647-55e0-9151-696cc175d6cb', '2a3ad241-587d-4853-bdc2-615d72496b4e', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20260930-11111', 'CASH_IN', 250000.00, '2026-09-30', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 30/09/2026', 'GRO-KIDINEST-20260930', 'ce47ec6a-5647-55e0-9151-696cc175d6cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7abf11e2-aa0b-5696-a2a2-e62c07e3de4b', 'ce47ec6a-5647-55e0-9151-696cc175d6cb', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2941353.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20260930-11120', 'CASH_IN', 2941353.00, '2026-09-30', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 30/09/2026', 'GRO-KIDINEST-20260930', 'ce47ec6a-5647-55e0-9151-696cc175d6cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3dd1cd37-ac7c-5d6d-afce-c1a84ba7dafc', 'ce47ec6a-5647-55e0-9151-696cc175d6cb', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20260930-11113', 'CASH_IN', 250000.00, '2026-09-30', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 30/09/2026', 'GRO-KIDINEST-20260930', 'ce47ec6a-5647-55e0-9151-696cc175d6cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('125f0dca-f7f0-52b2-bf66-1dc0f0bfd85c', 'ce47ec6a-5647-55e0-9151-696cc175d6cb', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e5feadff-0454-5e3a-99cb-6e8118644f6e', 'ce47ec6a-5647-55e0-9151-696cc175d6cb', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 941353.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20261001 (KIDINEST - 20261001): Rp 6,598,111.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('94426d5b-a839-5223-b6f3-321cdc5a3622', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261001', 'CASH_RECEIPT', '2026-10-01', '2026-10-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261001', 6598111.00, 6598111.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 01-10-2026 Daffi Raffasa Alfatih, Alesha Nareswari Pambudi, Medina Jannah Arafah, dkk. (total 11 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2f441648-6a7a-5786-81c3-b681a57e71dc', '94426d5b-a839-5223-b6f3-321cdc5a3622', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6167111.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261001-11120', 'CASH_IN', 6167111.00, '2026-10-01', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 01/10/2026', 'GRO-KIDINEST-20261001', '94426d5b-a839-5223-b6f3-321cdc5a3622')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8c541c1a-03ff-5a4a-a6b1-a21f8c82c42b', '94426d5b-a839-5223-b6f3-321cdc5a3622', '2a3ad241-587d-4853-bdc2-615d72496b4e', 431000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261001-11111', 'CASH_IN', 431000.00, '2026-10-01', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 01/10/2026', 'GRO-KIDINEST-20261001', '94426d5b-a839-5223-b6f3-321cdc5a3622')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bb1e5ea9-31a7-5f81-aeb4-1c4d76728638', '94426d5b-a839-5223-b6f3-321cdc5a3622', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2546753.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c9fe728d-ea85-564a-9555-f821e446a7c0', '94426d5b-a839-5223-b6f3-321cdc5a3622', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2976358.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d688c036-8fcb-577b-9b6c-a2d534c88125', '94426d5b-a839-5223-b6f3-321cdc5a3622', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 650000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e6cb3b1c-030c-5218-aaa3-fb320556d421', '94426d5b-a839-5223-b6f3-321cdc5a3622', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 425000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261002 (KIDINEST - 20261002): Rp 8,142,313.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('f9864c95-4a9b-57bb-aaf7-c3d20154331c', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261002', 'CASH_RECEIPT', '2026-10-02', '2026-10-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261002', 8142313.00, 8142313.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 02-10-2026 Kaysan Muhammad Andiansyah, Syahril Muadz Zayn, Azzea Zenith Sadiya, dkk. (total 8 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('19366fcc-a0b9-5546-9fd0-2bb4f5f24ffe', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3826543.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261002-11120', 'CASH_IN', 3826543.00, '2026-10-02', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 02/10/2026', 'GRO-KIDINEST-20261002', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a938f8a0-fb4b-5072-a73a-0abbfdf1a8e9', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '2a3ad241-587d-4853-bdc2-615d72496b4e', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261002-11111', 'CASH_IN', 250000.00, '2026-10-02', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 02/10/2026', 'GRO-KIDINEST-20261002', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5b5e1925-970a-527e-aeff-7b5d7960a18c', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4065770.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261002-11113', 'CASH_IN', 4065770.00, '2026-10-02', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 02/10/2026', 'GRO-KIDINEST-20261002', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dd6268f5-38d3-5747-a509-b646f752a21d', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4ff991e0-e278-5c18-8618-3a366032a8e6', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 4410000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c2aaa7b5-06dd-5c3a-841a-cc3a44051f2d', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1582313.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e2ac7862-8165-5cff-8667-3a68688b6dbe', 'f9864c95-4a9b-57bb-aaf7-c3d20154331c', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261003 (KIDINEST - 20261003): Rp 9,586,730.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0272abbe-9c8a-526d-837e-4ddebcb24c57', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261003', 'CASH_RECEIPT', '2026-10-03', '2026-10-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261003', 9586730.00, 9586730.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 03-10-2026 Rusdan Azmi Alfathan, Muhammad Aiden Pranadipa, Haziqul Riffat Nafian, dkk. (total 17 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('571088c5-c7e1-5cb7-82b7-3272e2ba5a77', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4984726.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261003-11120', 'CASH_IN', 4984726.00, '2026-10-03', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 03/10/2026', 'GRO-KIDINEST-20261003', '0272abbe-9c8a-526d-837e-4ddebcb24c57')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fddec1be-7fff-5cb5-b5ef-f71f2a1164d9', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3852004.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261003-11113', 'CASH_IN', 3852004.00, '2026-10-03', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 03/10/2026', 'GRO-KIDINEST-20261003', '0272abbe-9c8a-526d-837e-4ddebcb24c57')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3c651d14-1043-57e3-aff9-6fe189f3eeb3', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '2a3ad241-587d-4853-bdc2-615d72496b4e', 750000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261003-11111', 'CASH_IN', 750000.00, '2026-10-03', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 03/10/2026', 'GRO-KIDINEST-20261003', '0272abbe-9c8a-526d-837e-4ddebcb24c57')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('06444f23-3f47-5b72-9afe-700aa9f198ca', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c100a466-65be-512a-9527-d8c02ff178c0', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3410000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('78ca2077-c5e6-5258-9779-e109534794a1', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 33000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d28299da-2439-54a1-9630-40034e9573ca', '0272abbe-9c8a-526d-837e-4ddebcb24c57', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1893730.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20261004 (KIDINEST - 20261004): Rp 3,104,311.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7327ff7b-6a7d-54ac-b432-993e4fa6145c', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261004', 'CASH_RECEIPT', '2026-10-04', '2026-10-04 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261004', 3104311.00, 3104311.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 04-10-2026 Dairo Fiden Natawidjaja, Joshua Ethanuel Wirawan, Reyner Nathanuel Wirawan, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('80586d69-b0d7-5600-8a3a-5e77e018b070', '7327ff7b-6a7d-54ac-b432-993e4fa6145c', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 673075.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261004-11113', 'CASH_IN', 673075.00, '2026-10-04', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 04/10/2026', 'GRO-KIDINEST-20261004', '7327ff7b-6a7d-54ac-b432-993e4fa6145c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e77c841b-d6b2-5f6b-bfbc-e3fe71fbd3e0', '7327ff7b-6a7d-54ac-b432-993e4fa6145c', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2431236.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261004-11120', 'CASH_IN', 2431236.00, '2026-10-04', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 04/10/2026', 'GRO-KIDINEST-20261004', '7327ff7b-6a7d-54ac-b432-993e4fa6145c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('87fdccda-9e78-5086-ba0d-e07bb7b78665', '7327ff7b-6a7d-54ac-b432-993e4fa6145c', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('549cce19-b245-522e-8222-19e5d1a3bde4', '7327ff7b-6a7d-54ac-b432-993e4fa6145c', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 784311.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('62e87d0f-1720-5710-ab4f-b99ad7df2970', '7327ff7b-6a7d-54ac-b432-993e4fa6145c', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 570000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261005 (KIDINEST - 20261005): Rp 13,522,194.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c23cbf35-19fc-5cbd-8667-58c30696e279', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261005', 'CASH_RECEIPT', '2026-10-05', '2026-10-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261005', 13522194.00, 13522194.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 05-10-2026 Mahakertamada Bhre Dace, Arshaka Putra Arthura, Revalina Savrinadeya Almahyra, dkk. (total 17 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f829739b-780f-56df-bce0-2f213749d242', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1102641.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261005-11111', 'CASH_IN', 1102641.00, '2026-10-05', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 05/10/2026', 'GRO-KIDINEST-20261005', 'c23cbf35-19fc-5cbd-8667-58c30696e279')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('53d15cf4-54bf-53a0-b111-37a0e9003699', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 9783016.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261005-11120', 'CASH_IN', 9783016.00, '2026-10-05', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 05/10/2026', 'GRO-KIDINEST-20261005', 'c23cbf35-19fc-5cbd-8667-58c30696e279')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b051fbe1-7b11-5be3-89a9-d85ddcdc14ca', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2636537.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261005-11113', 'CASH_IN', 2636537.00, '2026-10-05', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 05/10/2026', 'GRO-KIDINEST-20261005', 'c23cbf35-19fc-5cbd-8667-58c30696e279')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b194983a-4bde-53b4-94d8-58c49e54c969', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('57efab23-b047-54a7-8e84-4492a690bbed', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 835000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('43f90790-5839-5522-88a5-c00d7e448358', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3687194.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d136cbab-8d09-5fa0-b942-f0be7eb9bb02', 'c23cbf35-19fc-5cbd-8667-58c30696e279', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5000000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261006 (KIDINEST - 20261006): Rp 9,250,147.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('4804f747-ada8-51b6-8aab-07d693936df7', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261006', 'CASH_RECEIPT', '2026-10-06', '2026-10-06 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261006', 9250147.00, 9250147.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 06-10-2026 Wilson Noel Suswanto, Abigael Tiffany, Garphly Walter Ginlead Purba, dkk. (total 15 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('09ca2ccc-6909-591e-bbf4-b4d81e200f5e', '4804f747-ada8-51b6-8aab-07d693936df7', '2a3ad241-587d-4853-bdc2-615d72496b4e', 701539.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261006-11111', 'CASH_IN', 701539.00, '2026-10-06', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 06/10/2026', 'GRO-KIDINEST-20261006', '4804f747-ada8-51b6-8aab-07d693936df7')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('35f2ec81-30df-51a6-833e-22f4ad8d664f', '4804f747-ada8-51b6-8aab-07d693936df7', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 7549510.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261006-11120', 'CASH_IN', 7549510.00, '2026-10-06', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 06/10/2026', 'GRO-KIDINEST-20261006', '4804f747-ada8-51b6-8aab-07d693936df7')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('85b65361-36cb-5da0-b23b-06185d2528d6', '4804f747-ada8-51b6-8aab-07d693936df7', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 999098.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261006-11113', 'CASH_IN', 999098.00, '2026-10-06', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 06/10/2026', 'GRO-KIDINEST-20261006', '4804f747-ada8-51b6-8aab-07d693936df7')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d25e8324-5731-5c84-9166-39f9ffdb0b3c', '4804f747-ada8-51b6-8aab-07d693936df7', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 3500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e41304f7-ed11-59c3-825c-d84105da68a9', '4804f747-ada8-51b6-8aab-07d693936df7', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2735147.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ff8753c7-ac0a-5147-bfd2-249691204c71', '4804f747-ada8-51b6-8aab-07d693936df7', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2815000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('84b88979-e372-5cbd-be38-89357ce3d304', '4804f747-ada8-51b6-8aab-07d693936df7', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 200000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261007 (KIDINEST - 20261007): Rp 4,362,583.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('65c891b2-b42e-5abf-a762-80efe38d6764', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261007', 'CASH_RECEIPT', '2026-10-07', '2026-10-07 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261007', 4362583.00, 4362583.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 07-10-2026 Charlotte Eloise Guritno, Arthur Louis Guritno, Tanayya Nadlyne Staraska Viansyah, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a6945bfa-8862-5f4d-b993-c3b8a917c78f', '65c891b2-b42e-5abf-a762-80efe38d6764', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1000000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261007-11111', 'CASH_IN', 1000000.00, '2026-10-07', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 07/10/2026', 'GRO-KIDINEST-20261007', '65c891b2-b42e-5abf-a762-80efe38d6764')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5c1f44ef-7e3b-5ee0-bf0f-d22162be2c67', '65c891b2-b42e-5abf-a762-80efe38d6764', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3112583.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261007-11120', 'CASH_IN', 3112583.00, '2026-10-07', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 07/10/2026', 'GRO-KIDINEST-20261007', '65c891b2-b42e-5abf-a762-80efe38d6764')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6969be40-5d8d-5960-b6fe-1b9f0c49ad60', '65c891b2-b42e-5abf-a762-80efe38d6764', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261007-11113', 'CASH_IN', 250000.00, '2026-10-07', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 07/10/2026', 'GRO-KIDINEST-20261007', '65c891b2-b42e-5abf-a762-80efe38d6764')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e90995a5-e365-5dfc-b2eb-35a0f4c209fa', '65c891b2-b42e-5abf-a762-80efe38d6764', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3686583.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('30480721-9879-5bf1-bf4c-27dc8337654c', '65c891b2-b42e-5abf-a762-80efe38d6764', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 676000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261008 (KIDINEST - 20261008): Rp 13,428,810.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('6b4368d4-1abd-5da7-9f79-ea804bf27e5d', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261008', 'CASH_RECEIPT', '2026-10-08', '2026-10-08 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261008', 13428810.00, 13428810.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 08-10-2026 Jihan Afifah S, Yarra Gjin Rosette Chandra, Odilia Claryne Auberta, dkk. (total 15 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c04c2a96-aea7-5192-98b5-4f4434cfb9c9', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 10651953.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261008-11120', 'CASH_IN', 10651953.00, '2026-10-08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 08/10/2026', 'GRO-KIDINEST-20261008', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d745bc62-4ebb-5306-a816-66ed169a26a6', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1488102.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261008-11111', 'CASH_IN', 1488102.00, '2026-10-08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 08/10/2026', 'GRO-KIDINEST-20261008', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c5e35af9-d5c7-5f21-a1a3-9c9b3e42515e', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1288755.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261008-11113', 'CASH_IN', 1288755.00, '2026-10-08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 08/10/2026', 'GRO-KIDINEST-20261008', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('caae1fff-2858-5f16-b3ca-bf205f1bc358', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 3500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8a7080bd-1f7c-528c-91f9-6905ca9f30e2', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 8150000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('92d43779-ff25-5f28-a7e1-9d9b01f0fb2b', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1617810.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('83d2dacd-27ba-53b0-bb4c-88cd23a43365', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 11000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('101e672e-2a92-5201-84cf-34934ac742a5', '6b4368d4-1abd-5da7-9f79-ea804bf27e5d', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261009 (KIDINEST - 20261009): Rp 3,125,803.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('28b7a077-b0ab-55bb-bd36-ebd884a9c502', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261009', 'CASH_RECEIPT', '2026-10-09', '2026-10-09 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261009', 3125803.00, 3125803.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 09-10-2026 Arsy Medina Bhatarajiwa, Seranala Aizalea Savana, Mason Agung Pramono, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('16ad276c-b415-5b8b-860d-0a1e9fac0db0', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1074420.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261009-11120', 'CASH_IN', 1074420.00, '2026-10-09', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 09/10/2026', 'GRO-KIDINEST-20261009', '28b7a077-b0ab-55bb-bd36-ebd884a9c502')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c97f3f87-0b3e-54c0-9fbd-aa8e954e1137', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1301383.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261009-11111', 'CASH_IN', 1301383.00, '2026-10-09', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 09/10/2026', 'GRO-KIDINEST-20261009', '28b7a077-b0ab-55bb-bd36-ebd884a9c502')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('59b795bd-6dab-5bd0-85a9-6abeca200cc2', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 750000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261009-11113', 'CASH_IN', 750000.00, '2026-10-09', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 09/10/2026', 'GRO-KIDINEST-20261009', '28b7a077-b0ab-55bb-bd36-ebd884a9c502')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('902c9d39-f126-5fd0-ba06-2eb47c850c23', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('387b72a5-6e19-5536-a9aa-3a5b601f10bf', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 535640.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ab17ed55-ba0c-5209-b0cd-87ea26dffbe1', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 745163.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('305f3fa3-f65e-559c-b6bc-d3879ba808f5', '28b7a077-b0ab-55bb-bd36-ebd884a9c502', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 345000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261010 (KIDINEST - 20261010): Rp 89,574.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7e37d218-9c0c-5c71-8772-958cba8ee01e', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261010', 'CASH_RECEIPT', '2026-10-10', '2026-10-10 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261010', 89574.00, 89574.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 10-10-2026 Natayana Haqqi Malik (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9880f9c2-70ec-59b8-afe8-dc1798064e7a', '7e37d218-9c0c-5c71-8772-958cba8ee01e', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 89574.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261010-11120', 'CASH_IN', 89574.00, '2026-10-10', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 10/10/2026', 'GRO-KIDINEST-20261010', '7e37d218-9c0c-5c71-8772-958cba8ee01e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e04f1f87-9835-5d82-866f-9d588e31b15e', '7e37d218-9c0c-5c71-8772-958cba8ee01e', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 89574.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20261011 (KIDINEST - 20261011): Rp 5,877,449.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('bb4748d0-efa1-5028-814d-265fc24b6705', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261011', 'CASH_RECEIPT', '2026-10-11', '2026-10-11 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261011', 5877449.00, 5877449.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 11-10-2026 Issaura Maryam Elora, Arshila Shadira Andarfa, Hartley Ardianto, dkk. (total 11 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('eeca9465-1bfe-5949-836f-0181c5754974', 'bb4748d0-efa1-5028-814d-265fc24b6705', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5342449.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261011-11120', 'CASH_IN', 5342449.00, '2026-10-11', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 11/10/2026', 'GRO-KIDINEST-20261011', 'bb4748d0-efa1-5028-814d-265fc24b6705')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e016bf2c-b41f-5da1-a158-fd94c3bb3717', 'bb4748d0-efa1-5028-814d-265fc24b6705', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 535000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261011-11113', 'CASH_IN', 535000.00, '2026-10-11', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 11/10/2026', 'GRO-KIDINEST-20261011', 'bb4748d0-efa1-5028-814d-265fc24b6705')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('38a1c5d2-1d06-55c2-bf7b-2dd1bd4eab89', 'bb4748d0-efa1-5028-814d-265fc24b6705', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9a1a9fce-0a8d-5570-9b29-6de1a7d1a943', 'bb4748d0-efa1-5028-814d-265fc24b6705', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2196000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('77f97752-486d-5c92-9d71-084a1fcc4f59', 'bb4748d0-efa1-5028-814d-265fc24b6705', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1420449.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0137baec-eb22-5d46-a0a6-5f7d791c2269', 'bb4748d0-efa1-5028-814d-265fc24b6705', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 11000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261012 (KIDINEST - 20261012): Rp 7,834,420.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('bfebc698-fc9b-56b6-a654-67141a38d21e', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261012', 'CASH_RECEIPT', '2026-10-12', '2026-10-12 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261012', 7834420.00, 7834420.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 12-10-2026 Arkananta Zaydan Manaf Rahmani, Arcelio Lanang Rayyanka, Elio Matthew Hamdani, dkk. (total 14 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c332a5a2-d975-5168-be82-eaa29a49c091', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5914468.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261012-11120', 'CASH_IN', 5914468.00, '2026-10-12', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 12/10/2026', 'GRO-KIDINEST-20261012', 'bfebc698-fc9b-56b6-a654-67141a38d21e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('25841fdb-6595-5e3d-9950-8d8995646c5d', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1545000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261012-11113', 'CASH_IN', 1545000.00, '2026-10-12', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 12/10/2026', 'GRO-KIDINEST-20261012', 'bfebc698-fc9b-56b6-a654-67141a38d21e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e33c3061-2325-5a31-96ec-ff1989025436', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '2a3ad241-587d-4853-bdc2-615d72496b4e', 374952.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261012-11111', 'CASH_IN', 374952.00, '2026-10-12', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 12/10/2026', 'GRO-KIDINEST-20261012', 'bfebc698-fc9b-56b6-a654-67141a38d21e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3eb3c043-3c19-513f-9030-65659e231b08', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 3250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('84be7e43-f34e-5969-83ed-23f37ee2bd52', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2710000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cf1db616-73b7-532f-b1c2-c6d55de23cda', 'bfebc698-fc9b-56b6-a654-67141a38d21e', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1874420.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20261013 (KIDINEST - 20261013): Rp 11,128,082.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261013', 'CASH_RECEIPT', '2026-10-13', '2026-10-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261013', 11128082.00, 11128082.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 13-10-2026 Odilia Claryne Auberta, Martha Evelyn Wibowo, Henna Maheswari Nur Khafidoh, dkk. (total 18 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('72bcf5ba-a1d9-53a3-bf04-eb1d8ecea8bf', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2141861.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261013-11113', 'CASH_IN', 2141861.00, '2026-10-13', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 13/10/2026', 'GRO-KIDINEST-20261013', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6de20857-7af5-5241-a2ab-995919def3d6', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '2a3ad241-587d-4853-bdc2-615d72496b4e', 3881015.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261013-11111', 'CASH_IN', 3881015.00, '2026-10-13', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 13/10/2026', 'GRO-KIDINEST-20261013', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('147cb1b7-7e50-5940-bb2e-1e8d00d0dbdf', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5105206.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261013-11120', 'CASH_IN', 5105206.00, '2026-10-13', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 13/10/2026', 'GRO-KIDINEST-20261013', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c5737857-79ec-5a24-b4de-8da0e4949457', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1570000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2dc05d1f-cff1-56ea-b286-fcb50653d496', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8e90540c-cc4b-5e43-86a4-d8344a759fa4', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 5058082.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('988a4acc-45e9-5c30-a893-f2d4fba18267', '95f2365a-0c99-5ba8-b6e8-1641da1b6a7d', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 500000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261014 (KIDINEST - 20261014): Rp 3,954,916.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261014', 'CASH_RECEIPT', '2026-10-14', '2026-10-14 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261014', 3954916.00, 3954916.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 14-10-2026 Henna Maheswari Nur Khafidoh, Elnathan Radeva Adzika, Alfarghani Averroes, dkk. (total 11 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6865a63b-1b12-5b2c-bbee-c7eeb887af65', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', '2a3ad241-587d-4853-bdc2-615d72496b4e', 554336.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261014-11111', 'CASH_IN', 554336.00, '2026-10-14', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 14/10/2026', 'GRO-KIDINEST-20261014', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6483132d-2fd9-581b-b0ff-96be95b51d3b', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2802446.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261014-11120', 'CASH_IN', 2802446.00, '2026-10-14', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 14/10/2026', 'GRO-KIDINEST-20261014', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0d8013a2-8509-5da0-b77f-9b3111e0a4d7', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 598134.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261014-11113', 'CASH_IN', 598134.00, '2026-10-14', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 14/10/2026', 'GRO-KIDINEST-20261014', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('89751dd8-6c3c-54c0-8166-175f6d39517f', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1454916.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5eaef18d-379f-5b41-9e31-29ff0e8f3698', '412e6d55-5a7c-54c4-8ada-4a6cd13a04b3', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261015 (KIDINEST - 20261015): Rp 10,856,884.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261015', 'CASH_RECEIPT', '2026-10-15', '2026-10-15 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261015', 10856884.00, 10856884.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 15-10-2026 Gwyneth Stevelyn, Maverick Bradley Handoko, Elmira Arunika Yasmine, dkk. (total 19 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6a9d5ac3-182b-597c-a5c4-2a10aa2004bf', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 9024599.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261015-11120', 'CASH_IN', 9024599.00, '2026-10-15', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 15/10/2026', 'GRO-KIDINEST-20261015', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6e16f79b-3b32-5dfb-a638-9906e1c3c4c2', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1554713.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261015-11111', 'CASH_IN', 1554713.00, '2026-10-15', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 15/10/2026', 'GRO-KIDINEST-20261015', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('86a82c05-dd68-52f8-b4af-3c0350940851', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 277572.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261015-11113', 'CASH_IN', 277572.00, '2026-10-15', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 15/10/2026', 'GRO-KIDINEST-20261015', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('43a7b69c-4904-52c2-ad08-5492df50e8d3', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3511884.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f5f9ac4f-068f-53b7-af0c-ce4b353e3942', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('73e65b70-ca2b-5ebb-8be2-58cb18b9387b', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3100000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('57d241d6-b35d-57d8-ac21-db9c8722e4c5', 'ac5e2185-4919-5a55-a6d7-e8acc9b2b94c', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 245000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261016 (KIDINEST - 20261016): Rp 4,792,520.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('bb8c2a2b-5020-503f-8ece-c2c70ca19f56', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261016', 'CASH_RECEIPT', '2026-10-16', '2026-10-16 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261016', 4792520.00, 4792520.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 16-10-2026 M. Arshaka Zayyan Alvarendra, Rafasya Kalandra Irghi, Ganendra Atharrizqi Y, dkk. (total 8 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ec47f499-e0f1-5341-826e-e7f547d31785', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2945111.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261016-11120', 'CASH_IN', 2945111.00, '2026-10-16', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 16/10/2026', 'GRO-KIDINEST-20261016', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('031a354b-f74d-5f68-b65d-349f73781260', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1847409.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261016-11111', 'CASH_IN', 1847409.00, '2026-10-16', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 16/10/2026', 'GRO-KIDINEST-20261016', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('89c29e5f-af54-5138-9b0b-b45c43e84f51', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6af94b7d-2945-5b9e-9dcc-2eec659e0d9e', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1797520.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f0cc1bd1-5152-58ec-aac2-ed82008bdf16', 'bb8c2a2b-5020-503f-8ece-c2c70ca19f56', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1245000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261017 (KIDINEST - 20261017): Rp 11,664,684.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('680cd4fe-0bc7-5ce6-a984-86d571c83ddb', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261017', 'CASH_RECEIPT', '2026-10-17', '2026-10-17 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261017', 11664684.00, 11664684.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 17-10-2026 Sierra Nadlyne, Zalif Ashif Mustafa, Kendrick Oliver Candra, dkk. (total 18 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fbae22ed-c96f-5a73-afb1-f660fdabafa9', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3575230.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261017-11113', 'CASH_IN', 3575230.00, '2026-10-17', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 17/10/2026', 'GRO-KIDINEST-20261017', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('89e51d65-6279-5aaf-bade-6b4559e7ebfc', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 8089454.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261017-11120', 'CASH_IN', 8089454.00, '2026-10-17', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 17/10/2026', 'GRO-KIDINEST-20261017', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('50fcb790-658c-5c5a-a9f6-0890fe3fd8a9', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f065aa68-2d33-5505-b069-ff8ded407b0f', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3994684.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4ff66fb7-6888-52d2-92b5-db09c91d0e94', '680cd4fe-0bc7-5ce6-a984-86d571c83ddb', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3420000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261018 (KIDINEST - 20261018): Rp 5,723,569.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('20904538-71ae-586a-be5a-df85780e4eeb', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261018', 'CASH_RECEIPT', '2026-10-18', '2026-10-18 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261018', 5723569.00, 5723569.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 18-10-2026 Michaela Kalashnikova, Keenan Athaya Bian, Ibrahim Alfa Romeo, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1f1efaf2-878f-5425-b1fe-d4a666455854', '20904538-71ae-586a-be5a-df85780e4eeb', '2a3ad241-587d-4853-bdc2-615d72496b4e', 538266.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261018-11111', 'CASH_IN', 538266.00, '2026-10-18', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 18/10/2026', 'GRO-KIDINEST-20261018', '20904538-71ae-586a-be5a-df85780e4eeb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ac3d8dd7-3a96-58c1-9cb2-8e62e2eae34b', '20904538-71ae-586a-be5a-df85780e4eeb', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4149119.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261018-11120', 'CASH_IN', 4149119.00, '2026-10-18', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 18/10/2026', 'GRO-KIDINEST-20261018', '20904538-71ae-586a-be5a-df85780e4eeb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d4dac010-4bd6-538c-aa64-2af78dd22edb', '20904538-71ae-586a-be5a-df85780e4eeb', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1036184.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261018-11113', 'CASH_IN', 1036184.00, '2026-10-18', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 18/10/2026', 'GRO-KIDINEST-20261018', '20904538-71ae-586a-be5a-df85780e4eeb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('83295a64-0b0c-512f-86f6-83cac7facfcf', '20904538-71ae-586a-be5a-df85780e4eeb', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2dcef190-4fef-598a-9993-564305cc73e0', '20904538-71ae-586a-be5a-df85780e4eeb', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2353569.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c8d6f745-5bfc-595a-86be-dd43ebf1eeb4', '20904538-71ae-586a-be5a-df85780e4eeb', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1450000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('55ad4634-1bbc-5034-9844-cafe1e328177', '20904538-71ae-586a-be5a-df85780e4eeb', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 170000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261019 (KIDINEST - 20261019): Rp 9,378,227.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3d5280d2-cfed-5d40-90c5-cba37c7e15f1', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261019', 'CASH_RECEIPT', '2026-10-19', '2026-10-19 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261019', 9378227.00, 9378227.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 19-10-2026 Gwen Nona Kagumi, Mathilda Caressa Anvienka, Celestial Shaula Amaryllis, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ac37c5ec-c596-5fe6-b9f7-7d67e92de32a', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6420227.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261019-11120', 'CASH_IN', 6420227.00, '2026-10-19', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 19/10/2026', 'GRO-KIDINEST-20261019', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9bf5692c-0506-5182-ba5e-68e5ccc1048a', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2958000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261019-11113', 'CASH_IN', 2958000.00, '2026-10-19', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 19/10/2026', 'GRO-KIDINEST-20261019', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('aec60f00-aa3e-572f-87a4-e8e94d31052e', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fc4e61d7-1e49-5bc0-bf72-841405605598', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 728727.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('603ad375-ade3-5b10-bb4d-4f54e4d88b23', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 6283000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bfb14188-8397-507d-8646-efe14fd22406', '3d5280d2-cfed-5d40-90c5-cba37c7e15f1', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 116500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261020 (KIDINEST - 20261020): Rp 13,244,853.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c8bdabb6-6910-587c-8a92-ab53c1162dbd', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261020', 'CASH_RECEIPT', '2026-10-20', '2026-10-20 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261020', 13244853.00, 13244853.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 20-10-2026 Saqueena Damar Bagisha, Migunani Saka Janitra, Kaelzio Gavin Dhanendra, dkk. (total 14 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('19f04380-376e-593f-a1b6-e70a4a1cc5a2', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 9800180.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261020-11120', 'CASH_IN', 9800180.00, '2026-10-20', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 20/10/2026', 'GRO-KIDINEST-20261020', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('56acd1db-04e4-517d-941a-b871708be81b', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1587950.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261020-11111', 'CASH_IN', 1587950.00, '2026-10-20', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 20/10/2026', 'GRO-KIDINEST-20261020', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ba0f1e44-65ed-5cd8-8196-aa082287cf67', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1856723.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261020-11113', 'CASH_IN', 1856723.00, '2026-10-20', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 20/10/2026', 'GRO-KIDINEST-20261020', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b762cd6e-e768-5552-b931-b3ce13d250f7', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bd7e966b-c949-5d38-8478-cac77a105e7b', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5271000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1585476d-4875-5bcd-8c62-a7afbd76e564', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 5588853.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5e377cc1-efc6-5670-934a-5e84e217fcd4', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 60000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9913895e-e26f-5faf-a898-4ac5985fe682', 'c8bdabb6-6910-587c-8a92-ab53c1162dbd', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 325000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261021 (KIDINEST - 20261021): Rp 1,426,824.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('8241b37a-e30b-5070-a7d0-bc71b8757fb5', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261021', 'CASH_RECEIPT', '2026-10-21', '2026-10-21 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261021', 1426824.00, 1426824.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 21-10-2026 Javier Harris Susanto, Rayyan Manggala Jezan (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9eb11958-a291-518e-bcb9-1e1954816533', '8241b37a-e30b-5070-a7d0-bc71b8757fb5', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1426824.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261021-11113', 'CASH_IN', 1426824.00, '2026-10-21', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 21/10/2026', 'GRO-KIDINEST-20261021', '8241b37a-e30b-5070-a7d0-bc71b8757fb5')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7a16d6cd-da76-598d-a112-327967ecaa56', '8241b37a-e30b-5070-a7d0-bc71b8757fb5', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 760000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('62583047-744f-52df-8235-00ad4a33193d', '8241b37a-e30b-5070-a7d0-bc71b8757fb5', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9ab6c4f8-eaf4-571b-b926-04a9baee17a3', '8241b37a-e30b-5070-a7d0-bc71b8757fb5', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 266824.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d44948dd-e242-58cc-887f-7fab6053fbfa', '8241b37a-e30b-5070-a7d0-bc71b8757fb5', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261022 (KIDINEST - 20261022): Rp 13,227,689.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('79f2498a-58df-5ff2-862b-a933a3430b28', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261022', 'CASH_RECEIPT', '2026-10-22', '2026-10-22 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261022', 13227689.00, 13227689.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 22-10-2026 Mahira Hameeda Shreya Eldyawan, Raffael Hidayat, Kamal Meliadi Rizqi Putta, dkk. (total 19 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fd51a388-90e9-5a88-8d30-15acbac4c4bb', '79f2498a-58df-5ff2-862b-a933a3430b28', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 6451250.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261022-11113', 'CASH_IN', 6451250.00, '2026-10-22', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 22/10/2026', 'GRO-KIDINEST-20261022', '79f2498a-58df-5ff2-862b-a933a3430b28')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1271f4a0-3fc1-5c4d-8678-162f5a32c9c7', '79f2498a-58df-5ff2-862b-a933a3430b28', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6526439.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261022-11120', 'CASH_IN', 6526439.00, '2026-10-22', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 22/10/2026', 'GRO-KIDINEST-20261022', '79f2498a-58df-5ff2-862b-a933a3430b28')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3e5d16f4-99d1-5009-bf3e-68db99c091a9', '79f2498a-58df-5ff2-862b-a933a3430b28', '2a3ad241-587d-4853-bdc2-615d72496b4e', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261022-11111', 'CASH_IN', 250000.00, '2026-10-22', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 22/10/2026', 'GRO-KIDINEST-20261022', '79f2498a-58df-5ff2-862b-a933a3430b28')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f5839bac-a858-5d04-9900-c8bc20ce8197', '79f2498a-58df-5ff2-862b-a933a3430b28', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4125000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3024ba4d-19a5-5336-bd83-2444410049b5', '79f2498a-58df-5ff2-862b-a933a3430b28', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 5057689.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e49be890-c9bc-5268-b3ce-d8a22c15bc7d', '79f2498a-58df-5ff2-862b-a933a3430b28', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3945000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('53ab2fc5-e6b7-50db-9a86-2f47b3c35825', '79f2498a-58df-5ff2-862b-a933a3430b28', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 100000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261023 (KIDINEST - 20261023): Rp 4,644,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('2e35b911-ca23-54c9-886c-0cbae556198a', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261023', 'CASH_RECEIPT', '2026-10-23', '2026-10-23 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261023', 4644500.00, 4644500.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 23-10-2026 Arsyakayla Kikanaya R, Esmeralda Shereen Rhaenyra, Hisham Gani Manggala, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b73867bc-4bac-5559-ac36-601978c9352f', '2e35b911-ca23-54c9-886c-0cbae556198a', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1345000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261023-11111', 'CASH_IN', 1345000.00, '2026-10-23', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 23/10/2026', 'GRO-KIDINEST-20261023', '2e35b911-ca23-54c9-886c-0cbae556198a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ba34debb-4c79-52b1-91d6-2a6dd3faeab1', '2e35b911-ca23-54c9-886c-0cbae556198a', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2210500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261023-11113', 'CASH_IN', 2210500.00, '2026-10-23', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 23/10/2026', 'GRO-KIDINEST-20261023', '2e35b911-ca23-54c9-886c-0cbae556198a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7f0936e2-321f-51a0-94f0-cf721d933a04', '2e35b911-ca23-54c9-886c-0cbae556198a', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1089000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261023-11120', 'CASH_IN', 1089000.00, '2026-10-23', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 23/10/2026', 'GRO-KIDINEST-20261023', '2e35b911-ca23-54c9-886c-0cbae556198a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ef1c7966-53e6-5407-88c9-a0c746d1f69f', '2e35b911-ca23-54c9-886c-0cbae556198a', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2258000.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('18d01815-d2cd-537e-aa9c-62d28bab7ef9', '2e35b911-ca23-54c9-886c-0cbae556198a', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2220000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6ba2d4b1-64a0-58a6-8b48-016384e1f21f', '2e35b911-ca23-54c9-886c-0cbae556198a', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 166500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261024 (KIDINEST - 20261024): Rp 13,685,914.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0c888d0c-3143-508b-90fb-cd8fd0be6370', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261024', 'CASH_RECEIPT', '2026-10-24', '2026-10-24 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261024', 13685914.00, 13685914.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 24-10-2026 Kalana Jennaira Wicaksana, Henna Maheswari Nur Khafidoh, Maximilian Axelson Santoso, dkk. (total 17 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3377f967-fdf1-51b1-9a0f-34818ad6ac27', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6843500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261024-11120', 'CASH_IN', 6843500.00, '2026-10-24', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 24/10/2026', 'GRO-KIDINEST-20261024', '0c888d0c-3143-508b-90fb-cd8fd0be6370')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5400b369-62b2-5ab7-9ddc-dd5001c48474', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 6342414.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261024-11113', 'CASH_IN', 6342414.00, '2026-10-24', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 24/10/2026', 'GRO-KIDINEST-20261024', '0c888d0c-3143-508b-90fb-cd8fd0be6370')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c765d8ce-04d8-5f7f-8409-f57c7f73d78a', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '2a3ad241-587d-4853-bdc2-615d72496b4e', 500000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261024-11111', 'CASH_IN', 500000.00, '2026-10-24', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 24/10/2026', 'GRO-KIDINEST-20261024', '0c888d0c-3143-508b-90fb-cd8fd0be6370')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2b27d4da-d4f4-54ef-bb62-cf373b309b2f', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5270000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ee4d5c67-e3ce-5915-aff8-d0b84687cb29', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 7059414.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7b168744-5d01-5300-a969-f69ec874cfea', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 606500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ad5a3e70-7c99-590d-9dbc-5af00f67cc1d', '0c888d0c-3143-508b-90fb-cd8fd0be6370', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261025 (KIDINEST - 20261025): Rp 7,612,916.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('14c523ff-cd3a-57fe-ae9e-776ea12f0f20', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261025', 'CASH_RECEIPT', '2026-10-25', '2026-10-25 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261025', 7612916.00, 7612916.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 25-10-2026 Arkananta Zaydan Manaf Rahmani, Issaura Maryam Elora, Maureen Kana Wijaya, dkk. (total 13 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('774fe28d-7d63-5481-ad7d-9813095d8d17', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5494028.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261025-11113', 'CASH_IN', 5494028.00, '2026-10-25', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 25/10/2026', 'GRO-KIDINEST-20261025', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('77f62e47-5897-5fa5-924c-fd75d3257889', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1776888.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261025-11120', 'CASH_IN', 1776888.00, '2026-10-25', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 25/10/2026', 'GRO-KIDINEST-20261025', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('551de003-dfcc-58ca-af47-8f4f078b5eb6', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '2a3ad241-587d-4853-bdc2-615d72496b4e', 342000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261025-11111', 'CASH_IN', 342000.00, '2026-10-25', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 25/10/2026', 'GRO-KIDINEST-20261025', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('996342ec-b7d8-53a4-b0f1-e65477d259af', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 3050000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4a483f9c-9c31-50f2-946b-9fe5652adc9d', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3135000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c989a3c5-ba8e-548c-925b-33399938c6d0', '14c523ff-cd3a-57fe-ae9e-776ea12f0f20', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1427916.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    -- Journal CR-KID-20261026 (KIDINEST - 20261026): Rp 5,601,503.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3cc3e591-9454-5232-b8d1-3a38cdb81474', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261026', 'CASH_RECEIPT', '2026-10-26', '2026-10-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261026', 5601503.00, 5601503.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 26-10-2026 Rainzo Anandra Al Ghaniy, Muhammad Rayyan Ar Rana, Maverick Seymour, dkk. (total 9 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ead144a2-af40-54b4-9e41-4136f620b00c', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2494000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261026-11120', 'CASH_IN', 2494000.00, '2026-10-26', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 26/10/2026', 'GRO-KIDINEST-20261026', '3cc3e591-9454-5232-b8d1-3a38cdb81474')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1969f604-c3be-587b-b146-657611b7cf41', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2857503.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261026-11113', 'CASH_IN', 2857503.00, '2026-10-26', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 26/10/2026', 'GRO-KIDINEST-20261026', '3cc3e591-9454-5232-b8d1-3a38cdb81474')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('53b8e2e1-00aa-5fe0-9c5c-dc9e78bf75b2', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '2a3ad241-587d-4853-bdc2-615d72496b4e', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261026-11111', 'CASH_IN', 250000.00, '2026-10-26', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 26/10/2026', 'GRO-KIDINEST-20261026', '3cc3e591-9454-5232-b8d1-3a38cdb81474')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('50bc56db-7fe2-53d0-9248-65fc40cbe863', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2386503.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b5ce1a88-49f5-5dd7-b5b5-7faf7f015aff', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1235000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d8e4a0aa-4ab0-5f92-a876-88e90e77e8c2', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d27940c5-b03a-50c3-8822-082800b50611', '3cc3e591-9454-5232-b8d1-3a38cdb81474', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 1230000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261027 (KIDINEST - 20261027): Rp 9,497,255.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('2adc8172-8312-5707-aa85-7eca535e9d8d', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261027', 'CASH_RECEIPT', '2026-10-27', '2026-10-27 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261027', 9497255.00, 9497255.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 27-10-2026 Aryasatya Sagara Nararya, Rasyid Ibrahim Setyawan, Muhammad Elman Arrasyi, dkk. (total 14 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('89173b3f-a0fc-555d-91e9-14f3de48b86d', '2adc8172-8312-5707-aa85-7eca535e9d8d', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5136755.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261027-11113', 'CASH_IN', 5136755.00, '2026-10-27', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 27/10/2026', 'GRO-KIDINEST-20261027', '2adc8172-8312-5707-aa85-7eca535e9d8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a61f15a6-50b2-5dd2-a916-10243b83eb0e', '2adc8172-8312-5707-aa85-7eca535e9d8d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3992500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261027-11120', 'CASH_IN', 3992500.00, '2026-10-27', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 27/10/2026', 'GRO-KIDINEST-20261027', '2adc8172-8312-5707-aa85-7eca535e9d8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2572c4b3-42fd-530f-8848-808d59459ab6', '2adc8172-8312-5707-aa85-7eca535e9d8d', '2a3ad241-587d-4853-bdc2-615d72496b4e', 368000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261027-11111', 'CASH_IN', 368000.00, '2026-10-27', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 27/10/2026', 'GRO-KIDINEST-20261027', '2adc8172-8312-5707-aa85-7eca535e9d8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ff41b485-27b6-59ab-90ff-6d5a64cbb51c', '2adc8172-8312-5707-aa85-7eca535e9d8d', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 5339255.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6d35f257-94bd-580f-9f32-316108b1bd38', '2adc8172-8312-5707-aa85-7eca535e9d8d', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 4158000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261028 (KIDINEST - 20261028): Rp 4,421,211.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ce4e156f-ad8a-5555-8420-b2c16824a258', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261028', 'CASH_RECEIPT', '2026-10-28', '2026-10-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261028', 4421211.00, 4421211.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 28-10-2026 Migunani Aleena Irghi, Abigael Tiffany, Arzen Baraka, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3ef82170-5fec-5b3b-8d96-522f60a00eff', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2066083.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261028-11120', 'CASH_IN', 2066083.00, '2026-10-28', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 28/10/2026', 'GRO-KIDINEST-20261028', 'ce4e156f-ad8a-5555-8420-b2c16824a258')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('49ebc333-4255-56a4-9d99-4f63ced44b1b', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2155128.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261028-11113', 'CASH_IN', 2155128.00, '2026-10-28', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 28/10/2026', 'GRO-KIDINEST-20261028', 'ce4e156f-ad8a-5555-8420-b2c16824a258')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2e519951-755a-5533-a440-e109751e951e', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '2a3ad241-587d-4853-bdc2-615d72496b4e', 200000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261028-11111', 'CASH_IN', 200000.00, '2026-10-28', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 28/10/2026', 'GRO-KIDINEST-20261028', 'ce4e156f-ad8a-5555-8420-b2c16824a258')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fddff03f-a52d-5fdb-8b24-abebeec69cb5', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3979711.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('de23bc6e-a776-5476-ae5c-26f9295b2df5', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 300000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3aa5900c-ec9f-5c19-a41f-6742596d77aa', 'ce4e156f-ad8a-5555-8420-b2c16824a258', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 141500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261029 (KIDINEST - 20261029): Rp 10,801,210.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ff4231af-5a6d-59ad-8b46-13007c663cef', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261029', 'CASH_RECEIPT', '2026-10-29', '2026-10-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261029', 10801210.00, 10801210.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 29-10-2026 Alfonsus Felix Hartanto, Olivia Jasmine Ferigo, Haziqul Riffat Nafian, dkk. (total 13 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f389f37c-fecb-54e4-8783-d200773d2853', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6909352.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261029-11120', 'CASH_IN', 6909352.00, '2026-10-29', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 29/10/2026', 'GRO-KIDINEST-20261029', 'ff4231af-5a6d-59ad-8b46-13007c663cef')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4bdf2f64-e8cd-5196-906c-575025b5b0cf', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3891858.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261029-11113', 'CASH_IN', 3891858.00, '2026-10-29', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 29/10/2026', 'GRO-KIDINEST-20261029', 'ff4231af-5a6d-59ad-8b46-13007c663cef')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c1b8d4cc-7aa9-570e-a349-4f67ba630b65', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 4931210.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('00a02a36-98fb-5f2f-823d-188adef768d7', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 4720000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bbdaafa3-5eb3-5237-8640-2231ad53c02d', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('69d89909-ae09-5812-9ea8-c9206f94946a', 'ff4231af-5a6d-59ad-8b46-13007c663cef', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261030 (KIDINEST - 20261030): Rp 5,706,452.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7f3e6fbc-30df-53ef-b575-9f8f5c81e43e', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261030', 'CASH_RECEIPT', '2026-10-30', '2026-10-30 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261030', 5706452.00, 5706452.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 30-10-2026 Eurico Max Yuarta, Kanawa Sih Rahadian, Rado Akarsana, dkk. (total 9 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('21df4c18-2b02-51ef-8894-3775f124ca87', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4141068.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261030-11113', 'CASH_IN', 4141068.00, '2026-10-30', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 30/10/2026', 'GRO-KIDINEST-20261030', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('519da250-abb9-5b3f-ab4d-08fa82eb5c96', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1565384.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261030-11120', 'CASH_IN', 1565384.00, '2026-10-30', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 30/10/2026', 'GRO-KIDINEST-20261030', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('833872b1-2711-515f-bd1b-394a87ed87b4', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3574452.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a4f383c5-68c1-5187-a185-3d3cb2898ce0', '7f3e6fbc-30df-53ef-b575-9f8f5c81e43e', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2132000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261031 (KIDINEST - 20261031): Rp 13,345,724.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('5c62fea3-350b-59a8-af58-672d3517834f', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261031', 'CASH_RECEIPT', '2026-10-31', '2026-10-31 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261031', 13345724.00, 13345724.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 31-10-2026 Hayyina Yafa Nasyitha, Lakeisha AninditaBaskoro, Humaira Rainsari Putri Negara, dkk. (total 15 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('286fb8b6-4403-5e2c-b8c4-dcec6e1d0620', '5c62fea3-350b-59a8-af58-672d3517834f', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3159305.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261031-11113', 'CASH_IN', 3159305.00, '2026-10-31', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 31/10/2026', 'GRO-KIDINEST-20261031', '5c62fea3-350b-59a8-af58-672d3517834f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('37b803ca-71da-5b2c-9e05-3adc8bb9cb35', '5c62fea3-350b-59a8-af58-672d3517834f', '2a3ad241-587d-4853-bdc2-615d72496b4e', 5378168.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261031-11111', 'CASH_IN', 5378168.00, '2026-10-31', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 31/10/2026', 'GRO-KIDINEST-20261031', '5c62fea3-350b-59a8-af58-672d3517834f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('151769d6-a299-55de-8dad-6c5f539c70e1', '5c62fea3-350b-59a8-af58-672d3517834f', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4808251.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261031-11120', 'CASH_IN', 4808251.00, '2026-10-31', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 31/10/2026', 'GRO-KIDINEST-20261031', '5c62fea3-350b-59a8-af58-672d3517834f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cf98ea2a-a5cc-5758-ad5c-f0fa3afb9902', '5c62fea3-350b-59a8-af58-672d3517834f', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 6780224.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5983dc7f-69a4-5d7d-9ed8-63c3e9c425d9', '5c62fea3-350b-59a8-af58-672d3517834f', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5607000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a57faa9e-42dd-57d9-9ffb-75cffabdb3b4', '5c62fea3-350b-59a8-af58-672d3517834f', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 22000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('abe16894-ce04-5e65-b46f-af3d3d99e472', '5c62fea3-350b-59a8-af58-672d3517834f', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 750000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4aca9b4c-622c-528c-845e-5dc2bd89d094', '5c62fea3-350b-59a8-af58-672d3517834f', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 186500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261101 (KIDINEST - 20261101): Rp 4,108,216.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('67d6640b-b6f7-58a4-bc16-7dff74a12595', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261101', 'CASH_RECEIPT', '2026-11-01', '2026-11-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261101', 4108216.00, 4108216.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 01-11-2026 Faruq Abdul Hafiz, Amran Ramadhan Arfianto, Bhre Isa Baihaqi, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f1a6e879-9eb2-570d-b5da-5afc9e09ca64', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '2a3ad241-587d-4853-bdc2-615d72496b4e', 535000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261101-11111', 'CASH_IN', 535000.00, '2026-11-01', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 01/11/2026', 'GRO-KIDINEST-20261101', '67d6640b-b6f7-58a4-bc16-7dff74a12595')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4f5751ac-5fa9-523d-911c-9a0c82efddc1', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3573216.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261101-11113', 'CASH_IN', 3573216.00, '2026-11-01', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 01/11/2026', 'GRO-KIDINEST-20261101', '67d6640b-b6f7-58a4-bc16-7dff74a12595')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ffb61cf4-2d46-5cf6-9968-c6210e648abf', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ab6aa959-d068-5375-adee-92d6f79d5cd0', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1120000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9fc55ccf-e105-5654-8337-f98db39ddf49', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2688216.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7b78c443-6a13-52bb-88b8-272e45792c3d', '67d6640b-b6f7-58a4-bc16-7dff74a12595', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 50000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261102 (KIDINEST - 20261102): Rp 9,020,034.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('66adf702-fe43-5464-9c5a-7232cc56f6cf', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261102', 'CASH_RECEIPT', '2026-11-02', '2026-11-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261102', 9020034.00, 9020034.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 02-11-2026 Deru Imanuel Mawadipta, Arkananta Zaydan Manaf Rahmani, Mahakertamada Bhre Dace, dkk. (total 13 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('61c10f7c-afc8-5d4a-b545-48b6dc52afc1', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5427154.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261102-11113', 'CASH_IN', 5427154.00, '2026-11-02', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 02/11/2026', 'GRO-KIDINEST-20261102', '66adf702-fe43-5464-9c5a-7232cc56f6cf')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9dd271b6-2fd2-5cf0-8159-36b7853b3cdb', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3592880.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261102-11120', 'CASH_IN', 3592880.00, '2026-11-02', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 02/11/2026', 'GRO-KIDINEST-20261102', '66adf702-fe43-5464-9c5a-7232cc56f6cf')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('790bdf65-6c00-56b7-8c41-bd31013be8df', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5430000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0c320fea-d2b4-5530-99b4-6d3066fbc14c', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2529034.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('511e9026-91c4-5c96-956c-711ca8c03aa2', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 11000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('584cb6b9-921d-5aa7-8b46-ecc9f3bd7db4', '66adf702-fe43-5464-9c5a-7232cc56f6cf', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1050000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261103 (KIDINEST - 20261103): Rp 11,862,752.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('4b9b8820-d460-5c0d-b80e-14909c5854ba', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261103', 'CASH_RECEIPT', '2026-11-03', '2026-11-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261103', 11862752.00, 11862752.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 03-11-2026 Alfonsus Felix Hartanto, Bazla Ahnaf Muhammad Fawwaz, Gheziya Aelgytha Naco, dkk. (total 18 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('77520acc-2193-5d54-9f50-8540c876d15a', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6368444.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261103-11120', 'CASH_IN', 6368444.00, '2026-11-03', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 03/11/2026', 'GRO-KIDINEST-20261103', '4b9b8820-d460-5c0d-b80e-14909c5854ba')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('484b15fe-8997-5475-84d6-a1673be7961c', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4288720.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261103-11113', 'CASH_IN', 4288720.00, '2026-11-03', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 03/11/2026', 'GRO-KIDINEST-20261103', '4b9b8820-d460-5c0d-b80e-14909c5854ba')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('91da842c-c871-5238-bf13-7742d0371a80', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '2a3ad241-587d-4853-bdc2-615d72496b4e', 1205588.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261103-11111', 'CASH_IN', 1205588.00, '2026-11-03', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 03/11/2026', 'GRO-KIDINEST-20261103', '4b9b8820-d460-5c0d-b80e-14909c5854ba')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5394dbbc-4c5a-501b-b4d8-5445f68a9480', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 8352752.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('40e77078-d672-532c-91a6-5e304d8658a0', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 600000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2106fd3f-511a-58ff-9d3f-17475d130017', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2170000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d2bf6537-5b0b-559d-9724-c67deba4224c', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 652500.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6fbc2c7f-e47b-5f5a-820e-9b1ab0b022cc', '4b9b8820-d460-5c0d-b80e-14909c5854ba', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 87500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261104 (KIDINEST - 20261104): Rp 5,510,504.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9b0c0976-fe67-5789-9209-0dccb8a65603', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261104', 'CASH_RECEIPT', '2026-11-04', '2026-11-04 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261104', 5510504.00, 5510504.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 04-11-2026 Gabriella Nicole Anggara, Evelyn Mikhayla Nayara, Arasya Adna Khalif, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('04128a5e-6dce-570c-9327-72a9a30548f6', '9b0c0976-fe67-5789-9209-0dccb8a65603', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3552500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261104-11120', 'CASH_IN', 3552500.00, '2026-11-04', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 04/11/2026', 'GRO-KIDINEST-20261104', '9b0c0976-fe67-5789-9209-0dccb8a65603')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e12bc7f9-7320-5912-b123-493f31728395', '9b0c0976-fe67-5789-9209-0dccb8a65603', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1958004.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261104-11113', 'CASH_IN', 1958004.00, '2026-11-04', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 04/11/2026', 'GRO-KIDINEST-20261104', '9b0c0976-fe67-5789-9209-0dccb8a65603')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f3e34e1f-50e6-5123-acfc-136ef8c698c4', '9b0c0976-fe67-5789-9209-0dccb8a65603', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 975000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d19d531e-21a0-540e-b2ce-43272e41cd5e', '9b0c0976-fe67-5789-9209-0dccb8a65603', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 50000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7e23b2b2-d0c4-5c74-9368-2f2237be5172', '9b0c0976-fe67-5789-9209-0dccb8a65603', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1575000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0cfd3627-8f15-5606-939c-41f30542e997', '9b0c0976-fe67-5789-9209-0dccb8a65603', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2867004.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b9074c5f-e95a-5183-85d7-6a8dbed8c9ac', '9b0c0976-fe67-5789-9209-0dccb8a65603', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 43500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261105 (KIDINEST - 20261105): Rp 15,356,385.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('664c330f-fece-5951-b731-e2ca853c6cd8', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261105', 'CASH_RECEIPT', '2026-11-05', '2026-11-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261105', 15356385.00, 15356385.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 05-11-2026 Arasya Adna Khalif, Khalif Danendra Sadha Ferdyansah, Reyhan Sulaiman Rahmandra, dkk. (total 27 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fdf644e8-a58f-5d37-8cc4-b652998be99e', '664c330f-fece-5951-b731-e2ca853c6cd8', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 8110307.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261105-11113', 'CASH_IN', 8110307.00, '2026-11-05', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 05/11/2026', 'GRO-KIDINEST-20261105', '664c330f-fece-5951-b731-e2ca853c6cd8')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('611b9aaf-f8fb-57f2-80a2-2b178acd1b57', '664c330f-fece-5951-b731-e2ca853c6cd8', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 7140578.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261105-11120', 'CASH_IN', 7140578.00, '2026-11-05', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 05/11/2026', 'GRO-KIDINEST-20261105', '664c330f-fece-5951-b731-e2ca853c6cd8')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5d30faf9-2be9-595e-80d2-c680be9aa0b7', '664c330f-fece-5951-b731-e2ca853c6cd8', '2a3ad241-587d-4853-bdc2-615d72496b4e', 105500.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261105-11111', 'CASH_IN', 105500.00, '2026-11-05', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 05/11/2026', 'GRO-KIDINEST-20261105', '664c330f-fece-5951-b731-e2ca853c6cd8')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b6028005-407b-5916-b554-0975f13f0231', '664c330f-fece-5951-b731-e2ca853c6cd8', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 4000885.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b42982d4-6255-53d7-aca2-1ec700014164', '664c330f-fece-5951-b731-e2ca853c6cd8', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2100000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('41c7083e-a117-58ec-b506-48a9668e6bf2', '664c330f-fece-5951-b731-e2ca853c6cd8', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 9228000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c8ad28f8-2ff7-5bac-95f8-e2817d6be692', '664c330f-fece-5951-b731-e2ca853c6cd8', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 27500.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261106 (KIDINEST - 20261106): Rp 7,131,467.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0c2462a8-651b-5450-9194-b6651ef36326', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261106', 'CASH_RECEIPT', '2026-11-06', '2026-11-06 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261106', 7131467.00, 7131467.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 06-11-2026 Javier Harris Susanto, Ryuga Ghafi Allendra, Celestial Shaula Amaryllis, dkk. (total 12 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a407326d-d12f-5db0-93d2-249c8848daaa', '0c2462a8-651b-5450-9194-b6651ef36326', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2923404.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261106-11120', 'CASH_IN', 2923404.00, '2026-11-06', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 06/11/2026', 'GRO-KIDINEST-20261106', '0c2462a8-651b-5450-9194-b6651ef36326')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3acb8103-469c-5382-bd65-af5d9664fdc6', '0c2462a8-651b-5450-9194-b6651ef36326', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4208063.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261106-11113', 'CASH_IN', 4208063.00, '2026-11-06', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 06/11/2026', 'GRO-KIDINEST-20261106', '0c2462a8-651b-5450-9194-b6651ef36326')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9d53726f-038a-578d-8ce1-ab08d6374ae4', '0c2462a8-651b-5450-9194-b6651ef36326', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2441467.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5fdb391f-9d3f-5d95-928a-e7abd3382e3f', '0c2462a8-651b-5450-9194-b6651ef36326', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3690000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4a31a37f-e41e-5c7e-8e25-f90df181dbec', '0c2462a8-651b-5450-9194-b6651ef36326', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261107 (KIDINEST - 20261107): Rp 15,779,329.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261107', 'CASH_RECEIPT', '2026-11-07', '2026-11-07 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261107', 15779329.00, 15779329.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 07-11-2026 Syahril Muadz Zayn, Noah Aaron Ferrizia Wiratama, Bhre Isa Baihaqi, dkk. (total 28 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b8d7b6da-6556-55f7-9c69-3c20a60eb079', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 9533655.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261107-11113', 'CASH_IN', 9533655.00, '2026-11-07', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 07/11/2026', 'GRO-KIDINEST-20261107', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('10dee5f2-0379-59eb-b3d9-79df79b80d78', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5794138.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261107-11120', 'CASH_IN', 5794138.00, '2026-11-07', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 07/11/2026', 'GRO-KIDINEST-20261107', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('57d8bd93-a376-5d6f-8f2c-81b30272d706', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '2a3ad241-587d-4853-bdc2-615d72496b4e', 451536.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261107-11111', 'CASH_IN', 451536.00, '2026-11-07', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 07/11/2026', 'GRO-KIDINEST-20261107', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('071de9e8-ba75-5d00-adbc-a76e6ee025ad', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 7535000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3139d56c-b244-5faa-820d-7487e9f2d266', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 4075829.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6e2f976a-817e-5a16-83d0-4c1f2bbd754a', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 4051000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a0a8c4cf-b1c7-5588-aa9a-ae1d8b50cccc', 'c9f5a8e1-befb-5bc0-bae7-22f39f8753cb', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 117500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261108 (KIDINEST - 20261108): Rp 10,011,954.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261108', 'CASH_RECEIPT', '2026-11-08', '2026-11-08 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261108', 10011954.00, 10011954.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 08-11-2026 Aldavin Nichel Ahmad, Seranala Aizalea Savana, Aleeya Jennaira Munaf, dkk. (total 17 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('66b0165a-1bab-58c8-b57c-686ec359655e', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 8114370.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261108-11113', 'CASH_IN', 8114370.00, '2026-11-08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 08/11/2026', 'GRO-KIDINEST-20261108', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5a857ead-2023-5f2b-9da6-29265871e80e', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1897584.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261108-11120', 'CASH_IN', 1897584.00, '2026-11-08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 08/11/2026', 'GRO-KIDINEST-20261108', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('778302a7-89c3-5649-832c-6249b46268d8', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 693000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fe3a00c1-8da1-5b1f-8ac0-c3cc5008b143', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 7100954.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9f340934-0023-5651-a40f-d222a574a9e3', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 271000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2aca15c4-400e-5dd1-9f00-72d402cb8b01', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 1936000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('01c1aaf5-077e-510e-962a-44b122da257e', 'ff0827b6-fdbe-594f-a9ac-14f53e60d5dd', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 11000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261109 (KIDINEST - 20261109): Rp 11,109,391.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9f08019c-6fbd-5f12-b52f-d4d0e1039385', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261109', 'CASH_RECEIPT', '2026-11-09', '2026-11-09 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261109', 11109391.00, 11109391.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 09-11-2026 Grace Vincentia Supraptono, Devanka Ilaan Altair, Mikayla Rosalind Onggabean, dkk. (total 25 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b8cf0b2f-fed3-5e0a-a8e0-7d46fab9b53c', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5833300.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261109-11113', 'CASH_IN', 5833300.00, '2026-11-09', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 09/11/2026', 'GRO-KIDINEST-20261109', '9f08019c-6fbd-5f12-b52f-d4d0e1039385')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4305d16f-8bc7-5734-a320-e8441967bf46', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4687588.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261109-11120', 'CASH_IN', 4687588.00, '2026-11-09', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 09/11/2026', 'GRO-KIDINEST-20261109', '9f08019c-6fbd-5f12-b52f-d4d0e1039385')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4bf6bbca-040d-5520-8051-db338e3bbda4', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '2a3ad241-587d-4853-bdc2-615d72496b4e', 588503.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261109-11111', 'CASH_IN', 588503.00, '2026-11-09', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 09/11/2026', 'GRO-KIDINEST-20261109', '9f08019c-6fbd-5f12-b52f-d4d0e1039385')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0fcefd6c-005c-5c6a-a234-374651dba451', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 4566000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('51dcc583-e6c0-56ea-af5b-776b76ac9022', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2757391.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2ff79e63-9f8d-549f-98b3-2755024a68f3', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 536000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('98c10961-d6e6-5cfb-be1e-629895e1a3d1', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2fa58749-381d-535d-8a01-191ddefc13cf', '9f08019c-6fbd-5f12-b52f-d4d0e1039385', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 750000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261110 (KIDINEST - 20261110): Rp 11,300,652.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('295abf7b-06fc-59ff-92a6-4e058ecd97af', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261110', 'CASH_RECEIPT', '2026-11-10', '2026-11-10 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261110', 11300652.00, 11300652.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 10-11-2026 Sierra Nadlyne, Odilia Claryne Auberta, Kiora Elena Maristella, dkk. (total 20 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('054fff0d-ce92-5321-bfdf-f9103546c67c', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2510000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261110-11120', 'CASH_IN', 2510000.00, '2026-11-10', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 10/11/2026', 'GRO-KIDINEST-20261110', '295abf7b-06fc-59ff-92a6-4e058ecd97af')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1a302ca2-060e-5116-a07f-dbac432b2ff1', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 7990038.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261110-11113', 'CASH_IN', 7990038.00, '2026-11-10', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 10/11/2026', 'GRO-KIDINEST-20261110', '295abf7b-06fc-59ff-92a6-4e058ecd97af')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('457a2fce-8240-5019-b3dd-afaa97db11f1', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '2a3ad241-587d-4853-bdc2-615d72496b4e', 800614.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261110-11111', 'CASH_IN', 800614.00, '2026-11-10', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 10/11/2026', 'GRO-KIDINEST-20261110', '295abf7b-06fc-59ff-92a6-4e058ecd97af')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('130a5f93-d39f-5a48-a6a0-d2b47f1ee72d', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 6185000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f476acfc-d37b-5cee-8cb2-e961b280fb24', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2925152.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ac861130-b28d-581b-b6ee-b7440ed41c3e', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 640500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('553ad669-0cd0-5446-a1c9-375ef977f900', '295abf7b-06fc-59ff-92a6-4e058ecd97af', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1550000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261111 (KIDINEST - 20261111): Rp 5,606,432.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0f82b621-1bc7-522c-9b0a-dc33556deb32', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261111', 'CASH_RECEIPT', '2026-11-11', '2026-11-11 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261111', 5606432.00, 5606432.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 11-11-2026 Shafaluna Adzqiara Savrinadeya, Kavi Arashaquille Dharmawan, Danisya Erudita Alfatiha, dkk. (total 15 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d077c24f-2670-5edf-829e-e3ad45e55ec6', '0f82b621-1bc7-522c-9b0a-dc33556deb32', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1016224.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261111-11120', 'CASH_IN', 1016224.00, '2026-11-11', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 11/11/2026', 'GRO-KIDINEST-20261111', '0f82b621-1bc7-522c-9b0a-dc33556deb32')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5e7e694c-2958-5093-a951-86b1a5836bac', '0f82b621-1bc7-522c-9b0a-dc33556deb32', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4590208.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261111-11113', 'CASH_IN', 4590208.00, '2026-11-11', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 11/11/2026', 'GRO-KIDINEST-20261111', '0f82b621-1bc7-522c-9b0a-dc33556deb32')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('36278963-99c9-5cc7-a1b1-e47e9b991844', '0f82b621-1bc7-522c-9b0a-dc33556deb32', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 2251432.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c4449b42-b3f0-5950-a3b4-4517e105a1e8', '0f82b621-1bc7-522c-9b0a-dc33556deb32', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3355000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-KID-20261112 (KIDINEST - 20261112): Rp 12,950,196.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('8b5ad27f-5314-52c7-a4cb-6619ea04f797', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261112', 'CASH_RECEIPT', '2026-11-12', '2026-11-12 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261112', 12950196.00, 12950196.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 12-11-2026 Hisyam Tsaqif Saputro, Yoona Valerie, Qiao Jinchuan, dkk. (total 28 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('78c82725-6460-587a-b88d-fab0607b5fff', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 9558059.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261112-11113', 'CASH_IN', 9558059.00, '2026-11-12', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 12/11/2026', 'GRO-KIDINEST-20261112', '8b5ad27f-5314-52c7-a4cb-6619ea04f797')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5c42e4ec-6f4b-585d-bc98-20cccf6f9e02', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '2a3ad241-587d-4853-bdc2-615d72496b4e', 9300.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261112-11111', 'CASH_IN', 9300.00, '2026-11-12', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 12/11/2026', 'GRO-KIDINEST-20261112', '8b5ad27f-5314-52c7-a4cb-6619ea04f797')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ee8be4ac-8580-5cd6-a92f-b9533a7c1997', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3382837.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261112-11120', 'CASH_IN', 3382837.00, '2026-11-12', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 12/11/2026', 'GRO-KIDINEST-20261112', '8b5ad27f-5314-52c7-a4cb-6619ea04f797')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fa1bd0d2-0f66-52b1-bbca-f5fc1e7a6342', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 3602000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fb293812-7eb1-530e-a944-9d89f9324ed6', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 6347696.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('83aa420d-3acd-50a7-90bb-fdfa9c951f58', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dd14153f-c5dc-5687-a020-2c9818de263a', '8b5ad27f-5314-52c7-a4cb-6619ea04f797', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2825000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261113 (KIDINEST - 20261113): Rp 5,798,319.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('a5590473-9fad-5ee1-bebf-69d6144a7a08', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261113', 'CASH_RECEIPT', '2026-11-13', '2026-11-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261113', 5798319.00, 5798319.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 13-11-2026 Yovella Meilody Puimera, Gianna Abigail Nugroho, Seranala Aizalea Savana, dkk. (total 15 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('06653719-3822-5c3f-8453-516eb6bc9f32', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3810732.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261113-11113', 'CASH_IN', 3810732.00, '2026-11-13', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 13/11/2026', 'GRO-KIDINEST-20261113', 'a5590473-9fad-5ee1-bebf-69d6144a7a08')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fe3f1b74-faca-5ff2-b6f7-4dadee88c3ff', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1846503.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261113-11120', 'CASH_IN', 1846503.00, '2026-11-13', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 13/11/2026', 'GRO-KIDINEST-20261113', 'a5590473-9fad-5ee1-bebf-69d6144a7a08')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1df96ba4-e9bc-5d01-985e-d5b2dda6f36b', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '2a3ad241-587d-4853-bdc2-615d72496b4e', 141084.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261113-11111', 'CASH_IN', 141084.00, '2026-11-13', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 13/11/2026', 'GRO-KIDINEST-20261113', 'a5590473-9fad-5ee1-bebf-69d6144a7a08')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9f1ec05c-6699-5421-89a1-029d85ebf75b', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3163319.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f5dde175-a3eb-5727-bd07-68f9073e89aa', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 2385000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dd04a7f8-2be4-58b5-8b76-70b20059492c', 'a5590473-9fad-5ee1-bebf-69d6144a7a08', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261114 (KIDINEST - 20261114): Rp 9,619,793.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3b6ed001-436d-582b-9558-90f2a5733eda', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261114', 'CASH_RECEIPT', '2026-11-14', '2026-11-14 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261114', 9619793.00, 9619793.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 14-11-2026 Khaira Najwa Azzahra, Benjamin Nathanael Sarwono, Migunani Saka Janitra, dkk. (total 21 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d6de7dbb-0112-5554-969c-e0754a57e2b5', '3b6ed001-436d-582b-9558-90f2a5733eda', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5526284.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261114-11120', 'CASH_IN', 5526284.00, '2026-11-14', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 14/11/2026', 'GRO-KIDINEST-20261114', '3b6ed001-436d-582b-9558-90f2a5733eda')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8ef9ab7a-1846-56db-9371-2dbc026960bf', '3b6ed001-436d-582b-9558-90f2a5733eda', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4093509.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261114-11113', 'CASH_IN', 4093509.00, '2026-11-14', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 14/11/2026', 'GRO-KIDINEST-20261114', '3b6ed001-436d-582b-9558-90f2a5733eda')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cd7b7cb9-d753-5840-a7d0-bf19a786ae13', '3b6ed001-436d-582b-9558-90f2a5733eda', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 4940000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7004d439-7d85-5cf8-a47a-a31a860564f5', '3b6ed001-436d-582b-9558-90f2a5733eda', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f639e262-a7c3-538a-b3b2-c44d5bfeea8a', '3b6ed001-436d-582b-9558-90f2a5733eda', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3139293.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('be6ce04d-e541-58fe-a0c5-fa7ccecf80fb', '3b6ed001-436d-582b-9558-90f2a5733eda', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 40500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261115 (KIDINEST - 20261115): Rp 3,664,940.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0f50cb69-062f-51b6-9944-c1dd0b7d9c85', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261115', 'CASH_RECEIPT', '2026-11-15', '2026-11-15 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261115', 3664940.00, 3664940.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 15-11-2026 Arsy Medina Bhatarajiwa, Mahakertamada Bhre Dace, Mareea Maryam Rilora, dkk. (total 12 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e90a1bcf-d262-5f5d-80a7-a6ecc0beb449', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1297136.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261115-11120', 'CASH_IN', 1297136.00, '2026-11-15', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 15/11/2026', 'GRO-KIDINEST-20261115', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('14f26967-03ca-5d12-aa00-6a48332ac215', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2367804.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261115-11113', 'CASH_IN', 2367804.00, '2026-11-15', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 15/11/2026', 'GRO-KIDINEST-20261115', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2840669a-a296-5842-abee-9e3d4a998909', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1199440.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('66e9a913-5bdd-5542-9916-16bb6f896e15', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 44000.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7e3f4800-0b09-50dd-8258-5a42c1370174', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 800000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('91654fb2-aa23-56c5-a3cb-3c08d82b4015', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 121500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('265472b8-c0b1-5919-989e-14f7de428ef7', '0f50cb69-062f-51b6-9944-c1dd0b7d9c85', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    -- Journal CR-KID-20261116 (KIDINEST - 20261116): Rp 13,607,071.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('f457a707-fab3-5f25-b705-58b94440d503', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261116', 'CASH_RECEIPT', '2026-11-16', '2026-11-16 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261116', 13607071.00, 13607071.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 16-11-2026 Camellia Khaira Arunika, Gavin Razelle Prakoso, Alara Farzana Setiawan, dkk. (total 24 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('eee48d5f-a6da-5bca-a0b2-efd40bdb730d', 'f457a707-fab3-5f25-b705-58b94440d503', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 9588916.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261116-11113', 'CASH_IN', 9588916.00, '2026-11-16', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 16/11/2026', 'GRO-KIDINEST-20261116', 'f457a707-fab3-5f25-b705-58b94440d503')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('27a993ef-72ff-5cbe-ab6c-d5473d072021', 'f457a707-fab3-5f25-b705-58b94440d503', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4018155.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261116-11120', 'CASH_IN', 4018155.00, '2026-11-16', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 16/11/2026', 'GRO-KIDINEST-20261116', 'f457a707-fab3-5f25-b705-58b94440d503')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7aabe63c-35d1-51eb-8158-628cad9773a0', 'f457a707-fab3-5f25-b705-58b94440d503', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 6193071.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5216c5a9-0da0-5c3f-94c4-6a48cfb669f2', 'f457a707-fab3-5f25-b705-58b94440d503', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 5240000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('01cf3a34-c7a9-59dd-8a79-d3c17af83fdb', 'f457a707-fab3-5f25-b705-58b94440d503', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 2000000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4d407d1a-0613-5ae2-aecc-db882c1a1097', 'f457a707-fab3-5f25-b705-58b94440d503', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 174000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-KID-20261117 (KIDINEST - 20261117): Rp 12,163,675.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ba19b845-5171-5923-aaa9-86c77b2d9bec', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261117', 'CASH_RECEIPT', '2026-11-17', '2026-11-17 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261117', 12163675.00, 12163675.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 17-11-2026 Aliyya Khairina Rahadi, Svetalana Mayra Angel Feisya, Sultan Nael Ahmed, dkk. (total 28 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('125cf143-b404-57d3-a4ea-1ba4bb67929b', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 6891228.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261117-11113', 'CASH_IN', 6891228.00, '2026-11-17', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 17/11/2026', 'GRO-KIDINEST-20261117', 'ba19b845-5171-5923-aaa9-86c77b2d9bec')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cd8bc2ab-638a-5f6b-b592-da0a7c8e6193', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4496164.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-KIDINEST-20261117-11120', 'CASH_IN', 4496164.00, '2026-11-17', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kliring Payment Gateway & EDC) - 17/11/2026', 'GRO-KIDINEST-20261117', 'ba19b845-5171-5923-aaa9-86c77b2d9bec')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9734e346-99d0-52bd-891e-af32c30b737b', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '2a3ad241-587d-4853-bdc2-615d72496b4e', 776283.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-KIDINEST-20261117-11111', 'CASH_IN', 776283.00, '2026-11-17', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 'Penerimaan Kasir Harian KIDINEST (Kas Kasir Kidinest) - 17/11/2026', 'GRO-KIDINEST-20261117', 'ba19b845-5171-5923-aaa9-86c77b2d9bec')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e3d8766b-b953-54a8-9203-0f015e3de697', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 7325000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('55dfe055-1924-5ebc-a7f7-b5d8472f8663', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 3322175.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6a550ac9-8908-5333-beb3-cc9bb77e8da9', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '44689ad5-00d4-487d-b1bb-98980bc4fdb8', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Jasa Dokter Spesialis Anak (Sp.A)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('054ee7ed-5c9d-5407-87b3-24d20c864f4b', 'ba19b845-5171-5923-aaa9-86c77b2d9bec', '14004cd8-d83d-4954-bcc2-72d68e1d6999', 0.00, 16500.00, 'Pendapatan Operasional Pendapatan Pendaftaran & Skrining Anak');
            

    -- Journal CR-KID-20261118 (KIDINEST - 20261118): Rp 2,440,531.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('027cadc2-743b-5f7f-9d19-666e7e1b3d5e', 'a0000000-0000-0000-0000-000000000001', 'CR-KID-20261118', 'CASH_RECEIPT', '2026-11-18', '2026-11-18 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-KIDINEST-20261118', 2440531.00, 2440531.00, 'Penerimaan Kasir Harian [KIDINEST] - Tanggal 18-11-2026 Naadhira Qalesya Wijaya, Oryza Sativa Florencia Zaen, Kaizel Sagara Halvar, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('30c26d02-91ef-53c9-8b0e-7bb8f02a93af', '027cadc2-743b-5f7f-9d19-666e7e1b3d5e', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2440531.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-KIDINEST-20261118-11113', 'CASH_IN', 2440531.00, '2026-11-18', '41a32972-5dba-4f6f-a175-4ce58b861e84', 'Penerimaan Kasir Harian KIDINEST (Bank BCA Operasional) - 18/11/2026', 'GRO-KIDINEST-20261118', '027cadc2-743b-5f7f-9d19-666e7e1b3d5e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('762c7e8e-44f2-587d-93e2-226dff7355b0', '027cadc2-743b-5f7f-9d19-666e7e1b3d5e', '41a32972-5dba-4f6f-a175-4ce58b861e84', 0.00, 1840531.00, 'Pendapatan Operasional Pendapatan Penjualan Obat & Racikan Anak');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('51352064-0748-537f-a61b-3f2e1b36e812', '027cadc2-743b-5f7f-9d19-666e7e1b3d5e', '89dcca29-dfc6-4758-a18f-2ca40b0756cd', 0.00, 600000.00, 'Pendapatan Operasional Pendapatan Vaksinasi & Imunisasi Anak');
            

    -- Journal CR-SSC-20260925 (SSC - 20260925): Rp 250,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('6ecd5580-3d38-52a1-af15-6380c0833a1d', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20260925', 'CASH_RECEIPT', '2026-09-25', '2026-09-25 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20260925', 250000.00, 250000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 25-09-2026 Oky Silvia (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('09b60230-2e40-5843-81c6-7d93432aa00b', '6ecd5580-3d38-52a1-af15-6380c0833a1d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20260925-11120', 'CASH_IN', 250000.00, '2026-09-25', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 25/09/2026', 'GRO-SSC-20260925', '6ecd5580-3d38-52a1-af15-6380c0833a1d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('57be738a-950c-5b51-b0d9-d5e944c022a8', '6ecd5580-3d38-52a1-af15-6380c0833a1d', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20260926 (SSC - 20260926): Rp 2,000,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20260926', 'CASH_RECEIPT', '2026-09-26', '2026-09-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20260926', 2000000.00, 2000000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 26-09-2026 Retno Sulistyowati, Fanny Listiani, Vivi Handayani, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bf660502-337c-5161-b61c-7f94c0a57aa1', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1000000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20260926-11120', 'CASH_IN', 1000000.00, '2026-09-26', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 26/09/2026', 'GRO-SSC-20260926', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('58690837-6027-55cf-a304-0be7a16c07e7', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 750000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20260926-11112', 'CASH_IN', 750000.00, '2026-09-26', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 26/09/2026', 'GRO-SSC-20260926', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('661bb368-411b-5f73-be4b-14287c2c1a72', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20260926-11113', 'CASH_IN', 250000.00, '2026-09-26', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 26/09/2026', 'GRO-SSC-20260926', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dac1e6ab-4a59-50cd-90c4-75de5ebf5bde', '4dcd3d35-c2b8-56e6-8bdd-d3e89ff9ff35', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2000000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20260928 (SSC - 20260928): Rp 1,500,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('52916fa2-fb57-5718-8f9e-319be2dcd199', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20260928', 'CASH_RECEIPT', '2026-09-28', '2026-09-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20260928', 1500000.00, 1500000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 28-09-2026 Marcella Davina Sulistyani, Fitri Dwi Aprilia, Olivia Elena H, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('658d7990-fc9b-5234-83db-f406e0b07316', '52916fa2-fb57-5718-8f9e-319be2dcd199', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 750000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20260928-11120', 'CASH_IN', 750000.00, '2026-09-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 28/09/2026', 'GRO-SSC-20260928', '52916fa2-fb57-5718-8f9e-319be2dcd199')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('286cb146-a00a-5e3d-aada-a0369157eb63', '52916fa2-fb57-5718-8f9e-319be2dcd199', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 500000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20260928-11113', 'CASH_IN', 500000.00, '2026-09-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 28/09/2026', 'GRO-SSC-20260928', '52916fa2-fb57-5718-8f9e-319be2dcd199')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('efb2ff5d-e3cd-58e2-bc01-86fdad751536', '52916fa2-fb57-5718-8f9e-319be2dcd199', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20260928-11112', 'CASH_IN', 250000.00, '2026-09-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 28/09/2026', 'GRO-SSC-20260928', '52916fa2-fb57-5718-8f9e-319be2dcd199')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ac9d11e7-a4d2-54f4-ab54-acb8192aa639', '52916fa2-fb57-5718-8f9e-319be2dcd199', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1250000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('aef4496a-912c-5feb-8234-b0f63475951b', '52916fa2-fb57-5718-8f9e-319be2dcd199', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20260929 (SSC - 20260929): Rp 950,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('f8ab3832-88b5-5274-8db8-6003fc0f8956', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20260929', 'CASH_RECEIPT', '2026-09-29', '2026-09-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20260929', 950000.00, 950000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 29-09-2026 Rachel Wansa, Vivian Eunike, Novylia Susanti Hadi, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('aee60d42-ff92-5eee-a6a8-37175300a030', 'f8ab3832-88b5-5274-8db8-6003fc0f8956', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 700000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20260929-11113', 'CASH_IN', 700000.00, '2026-09-29', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 29/09/2026', 'GRO-SSC-20260929', 'f8ab3832-88b5-5274-8db8-6003fc0f8956')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f3410a56-ffad-5b08-b705-b2ea460a4893', 'f8ab3832-88b5-5274-8db8-6003fc0f8956', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20260929-11112', 'CASH_IN', 250000.00, '2026-09-29', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 29/09/2026', 'GRO-SSC-20260929', 'f8ab3832-88b5-5274-8db8-6003fc0f8956')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('566ce6bf-8491-5690-9215-9802dc4c3617', 'f8ab3832-88b5-5274-8db8-6003fc0f8956', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 950000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20260930 (SSC - 20260930): Rp 5,137,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3d801087-44e9-5667-bf8f-fec9bb74f908', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20260930', 'CASH_RECEIPT', '2026-09-30', '2026-09-30 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20260930', 5137000.00, 5137000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 30-09-2026 Jessica Amara Wijaya, Cecilia Lely Prima, Kwik Moey Tjen, dkk. (total 8 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('897d55b2-78dc-58bf-a7b3-da6b9911c391', '3d801087-44e9-5667-bf8f-fec9bb74f908', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 612000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20260930-11112', 'CASH_IN', 612000.00, '2026-09-30', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 30/09/2026', 'GRO-SSC-20260930', '3d801087-44e9-5667-bf8f-fec9bb74f908')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0d4ac7d9-2fb6-5027-9b05-3691ff6ca1ab', '3d801087-44e9-5667-bf8f-fec9bb74f908', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4275000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20260930-11120', 'CASH_IN', 4275000.00, '2026-09-30', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 30/09/2026', 'GRO-SSC-20260930', '3d801087-44e9-5667-bf8f-fec9bb74f908')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8fe7194a-cfd5-5235-ab87-dfd88067d853', '3d801087-44e9-5667-bf8f-fec9bb74f908', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20260930-11113', 'CASH_IN', 250000.00, '2026-09-30', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 30/09/2026', 'GRO-SSC-20260930', '3d801087-44e9-5667-bf8f-fec9bb74f908')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c89beb6d-1fae-5f31-aee1-22e5558584de', '3d801087-44e9-5667-bf8f-fec9bb74f908', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2712000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('335b3267-773b-5467-9375-5349f7d24305', '3d801087-44e9-5667-bf8f-fec9bb74f908', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 125000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('72a72cb0-b19a-5831-a6d2-74f6576de683', '3d801087-44e9-5667-bf8f-fec9bb74f908', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 2300000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261001 (SSC - 20261001): Rp 10,155,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('521e6bcf-4b2d-5700-aa1d-ef6be372d200', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261001', 'CASH_RECEIPT', '2026-10-01', '2026-10-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261001', 10155000.00, 10155000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 01-10-2026 Silvia Soegiar S, Tjan Puput Wulandari H, Siti Aryani, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('00f9a79f-1d34-597d-8cf4-e7d91b47c452', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1080000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261001-11120', 'CASH_IN', 1080000.00, '2026-10-01', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 01/10/2026', 'GRO-SSC-20261001', '521e6bcf-4b2d-5700-aa1d-ef6be372d200')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('67cadfbd-0133-5163-a114-028087739455', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 500000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261001-11112', 'CASH_IN', 500000.00, '2026-10-01', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 01/10/2026', 'GRO-SSC-20261001', '521e6bcf-4b2d-5700-aa1d-ef6be372d200')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('09da15e2-8a5d-5b3d-82ab-652665a795f7', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 8575000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261001-11113', 'CASH_IN', 8575000.00, '2026-10-01', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 01/10/2026', 'GRO-SSC-20261001', '521e6bcf-4b2d-5700-aa1d-ef6be372d200')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7a9108cd-cd68-51c2-8493-01f7c347b0a8', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1950000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('39ad9b44-f437-5a3d-ac66-e3c1f8e6a576', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7971f91b-0a6a-5cbd-93b7-82cfddf36a3a', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 230000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a6127a07-d7f2-589a-b2ef-0a62ee82247c', '521e6bcf-4b2d-5700-aa1d-ef6be372d200', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 7725000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261002 (SSC - 20261002): Rp 2,150,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('91f5fcea-f658-5e2f-951a-6499a83b1509', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261002', 'CASH_RECEIPT', '2026-10-02', '2026-10-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261002', 2150000.00, 2150000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 02-10-2026 Felita, Fatiha Asti Amalia, Fanny Listiani, dkk. (total 9 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('85fc48e2-dc47-5d1c-9387-6bb2d84a84cc', '91f5fcea-f658-5e2f-951a-6499a83b1509', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1700000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261002-11120', 'CASH_IN', 1700000.00, '2026-10-02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 02/10/2026', 'GRO-SSC-20261002', '91f5fcea-f658-5e2f-951a-6499a83b1509')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3e089329-5bc2-5e52-b7ec-aa1c18c8084d', '91f5fcea-f658-5e2f-951a-6499a83b1509', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 250000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261002-11112', 'CASH_IN', 250000.00, '2026-10-02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 02/10/2026', 'GRO-SSC-20261002', '91f5fcea-f658-5e2f-951a-6499a83b1509')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('69afe27b-cc82-55c8-bee1-2e69d1518fb2', '91f5fcea-f658-5e2f-951a-6499a83b1509', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 200000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261002-11113', 'CASH_IN', 200000.00, '2026-10-02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 02/10/2026', 'GRO-SSC-20261002', '91f5fcea-f658-5e2f-951a-6499a83b1509')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('51d8c48d-7c03-54dc-82b7-5b7485275461', '91f5fcea-f658-5e2f-951a-6499a83b1509', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2150000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261003 (SSC - 20261003): Rp 275,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9a4b2d0e-bf1e-5104-bc5d-7c56c8db3183', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261003', 'CASH_RECEIPT', '2026-10-03', '2026-10-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261003', 275000.00, 275000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 03-10-2026 Dzakkiyatul Ulya Alfiannas Tito (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('960de270-26f1-557a-ab36-eb5fd68dcb12', '9a4b2d0e-bf1e-5104-bc5d-7c56c8db3183', 'fb5c3156-3a4b-4a86-a896-ec354caef246', 275000.00, 0.00, 'Penerimaan Kas/Bank via Uang Muka Pasien / Deposit Tindakan');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('786935b9-f5a1-564d-86ad-bf3c6eb36a87', '9a4b2d0e-bf1e-5104-bc5d-7c56c8db3183', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 275000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261005 (SSC - 20261005): Rp 922,883.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261005', 'CASH_RECEIPT', '2026-10-05', '2026-10-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261005', 922883.00, 922883.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 05-10-2026 Hendra, Diasa Ayu Raharni, Ambar Prihartini (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a26a1065-c41b-53cd-b86b-518a13cd3e7e', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 550000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261005-11120', 'CASH_IN', 550000.00, '2026-10-05', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 05/10/2026', 'GRO-SSC-20261005', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('abbd5d10-33c8-57ab-9d15-37212e4889eb', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 372883.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261005-11113', 'CASH_IN', 372883.00, '2026-10-05', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 05/10/2026', 'GRO-SSC-20261005', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d2c00db5-8d4e-51a0-b1ed-d232bcbe9f06', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 550000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('88f36261-8a46-5350-a07f-edb294286f14', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 204895.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6bd2bef3-372a-5d6f-b352-4042064eb061', 'ce6bc527-1de3-530c-9b73-ce3a2e5a85c6', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 167988.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261006 (SSC - 20261006): Rp 4,741,398.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('14b96c44-71cb-51e4-8a10-eeb849615409', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261006', 'CASH_RECEIPT', '2026-10-06', '2026-10-06 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261006', 4741398.00, 4741398.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 06-10-2026 Ranu Hamukti Kusuma, Tennesha Nelafrinanda, Novalia Kuntardjo, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('095dd73e-e5d2-57bf-9006-063d961333f3', '14b96c44-71cb-51e4-8a10-eeb849615409', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1975000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261006-11113', 'CASH_IN', 1975000.00, '2026-10-06', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 06/10/2026', 'GRO-SSC-20261006', '14b96c44-71cb-51e4-8a10-eeb849615409')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c6c726e4-ec8a-5656-81d7-6cae47075cec', '14b96c44-71cb-51e4-8a10-eeb849615409', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2766398.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261006-11120', 'CASH_IN', 2766398.00, '2026-10-06', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 06/10/2026', 'GRO-SSC-20261006', '14b96c44-71cb-51e4-8a10-eeb849615409')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('39115647-50ba-55ae-9843-60669e8d388c', '14b96c44-71cb-51e4-8a10-eeb849615409', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 450000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6c8807dc-f7cb-5b05-9bb9-826e0a221268', '14b96c44-71cb-51e4-8a10-eeb849615409', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 4155000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bac0938a-43bf-512e-a556-e91b3f8c4890', '14b96c44-71cb-51e4-8a10-eeb849615409', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 136398.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261008 (SSC - 20261008): Rp 5,380,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('89c83cce-3827-53cc-935c-0088e7e437e6', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261008', 'CASH_RECEIPT', '2026-10-08', '2026-10-08 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261008', 5380000.00, 5380000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 08-10-2026 Lie Lee Hong (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('530e8258-7157-5869-b9c4-abdc63632140', '89c83cce-3827-53cc-935c-0088e7e437e6', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5045000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261008-11113', 'CASH_IN', 5045000.00, '2026-10-08', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 08/10/2026', 'GRO-SSC-20261008', '89c83cce-3827-53cc-935c-0088e7e437e6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('daca393c-4ca9-5bf1-83ef-8d4268c9f523', '89c83cce-3827-53cc-935c-0088e7e437e6', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 335000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261008-11112', 'CASH_IN', 335000.00, '2026-10-08', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 08/10/2026', 'GRO-SSC-20261008', '89c83cce-3827-53cc-935c-0088e7e437e6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('53fe1944-c377-5bb7-bb05-aa3df43ec484', '89c83cce-3827-53cc-935c-0088e7e437e6', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 225000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e26ab3a5-047f-592c-bad9-81e8ec6e606a', '89c83cce-3827-53cc-935c-0088e7e437e6', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 5030000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('28e5d0ea-16a8-56a8-9cb3-69af111d6801', '89c83cce-3827-53cc-935c-0088e7e437e6', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 125000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261009 (SSC - 20261009): Rp 550,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('5a381c75-cb3d-59b5-8779-0074963f4278', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261009', 'CASH_RECEIPT', '2026-10-09', '2026-10-09 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261009', 550000.00, 550000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 09-10-2026 Rony Hendarto, Dian Nurul (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('39a0ff1d-43f6-5a30-8c24-86b00b103cc0', '5a381c75-cb3d-59b5-8779-0074963f4278', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 550000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261009-11120', 'CASH_IN', 550000.00, '2026-10-09', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 09/10/2026', 'GRO-SSC-20261009', '5a381c75-cb3d-59b5-8779-0074963f4278')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('406d92ab-e1ba-55da-abda-f011d0d701b8', '5a381c75-cb3d-59b5-8779-0074963f4278', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 550000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261010 (SSC - 20261010): Rp 787,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9d89deb5-163d-5e9e-bda4-c9a7fc2c679b', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261010', 'CASH_RECEIPT', '2026-10-10', '2026-10-10 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261010', 787500.00, 787500.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 10-10-2026 Kirana Ramadani Putri Susilawatu, Yanti Sri Susilawati, Sifa Adila, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f4cdc992-2ff7-57a5-8da6-fa27025ee9df', '9d89deb5-163d-5e9e-bda4-c9a7fc2c679b', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 425000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261010-11120', 'CASH_IN', 425000.00, '2026-10-10', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 10/10/2026', 'GRO-SSC-20261010', '9d89deb5-163d-5e9e-bda4-c9a7fc2c679b')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fcc2f58e-2824-57ab-a731-0e004fede9cd', '9d89deb5-163d-5e9e-bda4-c9a7fc2c679b', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 362500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261010-11113', 'CASH_IN', 362500.00, '2026-10-10', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 10/10/2026', 'GRO-SSC-20261010', '9d89deb5-163d-5e9e-bda4-c9a7fc2c679b')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c7fd2618-d5a5-57ce-bc1e-ef62fcb8b7e8', '9d89deb5-163d-5e9e-bda4-c9a7fc2c679b', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 787500.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261011 (SSC - 20261011): Rp 6,275,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('4711efee-340d-5db7-9945-8b05bce8c52f', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261011', 'CASH_RECEIPT', '2026-10-11', '2026-10-11 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261011', 6275000.00, 6275000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 11-10-2026 Rima, Arshila Shadira Andarfa, Yenny Permatasari Santoso, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2854b89d-61ee-5ba0-84db-e69f99c6d28e', '4711efee-340d-5db7-9945-8b05bce8c52f', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 325000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261011-11120', 'CASH_IN', 325000.00, '2026-10-11', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 11/10/2026', 'GRO-SSC-20261011', '4711efee-340d-5db7-9945-8b05bce8c52f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cf32039c-092f-5d6e-9eb4-850126ad57b0', '4711efee-340d-5db7-9945-8b05bce8c52f', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 5950000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261011-11112', 'CASH_IN', 5950000.00, '2026-10-11', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 11/10/2026', 'GRO-SSC-20261011', '4711efee-340d-5db7-9945-8b05bce8c52f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a5f41d4f-eff8-5b21-b5a2-89787898402c', '4711efee-340d-5db7-9945-8b05bce8c52f', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 550000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2b665db3-cac3-5fbe-af86-bac7ee077117', '4711efee-340d-5db7-9945-8b05bce8c52f', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 5525000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e443d0e0-af24-5cda-9056-fd66ac40ff3a', '4711efee-340d-5db7-9945-8b05bce8c52f', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4260a5e8-a796-54bc-a8e5-7eadcbba98b0', '4711efee-340d-5db7-9945-8b05bce8c52f', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 50000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261012 (SSC - 20261012): Rp 3,196,765.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('288d3199-7d11-5b37-b71b-991f8f285c96', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261012', 'CASH_RECEIPT', '2026-10-12', '2026-10-12 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261012', 3196765.00, 3196765.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 12-10-2026 Listriani Gunawan, Adriana Eka, Devina Yuliastanti Setianto, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e1df9db1-ecf7-5f1c-941b-8f24e115fb2f', '288d3199-7d11-5b37-b71b-991f8f285c96', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 198350.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261012-11112', 'CASH_IN', 198350.00, '2026-10-12', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 12/10/2026', 'GRO-SSC-20261012', '288d3199-7d11-5b37-b71b-991f8f285c96')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('201dfa1a-536f-5309-81e0-442f57abb8e5', '288d3199-7d11-5b37-b71b-991f8f285c96', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2842000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261012-11113', 'CASH_IN', 2842000.00, '2026-10-12', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 12/10/2026', 'GRO-SSC-20261012', '288d3199-7d11-5b37-b71b-991f8f285c96')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3d25c597-a218-5836-be5c-482bad32a34b', '288d3199-7d11-5b37-b71b-991f8f285c96', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 156415.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261012-11120', 'CASH_IN', 156415.00, '2026-10-12', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 12/10/2026', 'GRO-SSC-20261012', '288d3199-7d11-5b37-b71b-991f8f285c96')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3d7bbb50-dd03-5947-952c-858c6fa8fc28', '288d3199-7d11-5b37-b71b-991f8f285c96', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 373350.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d1f5b22b-0583-5bdc-8553-cc00349b3485', '288d3199-7d11-5b37-b71b-991f8f285c96', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2277000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('aee979b7-3ba0-5724-bd8f-ca6ac5607e6a', '288d3199-7d11-5b37-b71b-991f8f285c96', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 546415.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261013 (SSC - 20261013): Rp 2,140,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('b2906cd0-470e-54f2-8f9f-61c602acbebd', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261013', 'CASH_RECEIPT', '2026-10-13', '2026-10-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261013', 2140000.00, 2140000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 13-10-2026 Ade Ayu Surya (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('77010644-ac91-52fe-9d15-7437d91bc8a9', 'b2906cd0-470e-54f2-8f9f-61c602acbebd', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2140000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261013-11120', 'CASH_IN', 2140000.00, '2026-10-13', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 13/10/2026', 'GRO-SSC-20261013', 'b2906cd0-470e-54f2-8f9f-61c602acbebd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e91a1730-40cc-5963-922e-59a06e1d4720', 'b2906cd0-470e-54f2-8f9f-61c602acbebd', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 200000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2fcce390-1281-598c-bdb3-499c73866cc9', 'b2906cd0-470e-54f2-8f9f-61c602acbebd', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1850000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ad85df37-a988-584d-aa4f-eda6a15dfed5', 'b2906cd0-470e-54f2-8f9f-61c602acbebd', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 90000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261014 (SSC - 20261014): Rp 1,313,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('13edc370-8cbd-5c8c-a54d-8f5a13df1d18', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261014', 'CASH_RECEIPT', '2026-10-14', '2026-10-14 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261014', 1313000.00, 1313000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 14-10-2026 Upik Nur Aprilia, Mahakertamada Bhre Dace, Stephanie Emmanuel P S, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bac1995b-03c3-5974-97cf-b9a2e999a396', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1113000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261014-11120', 'CASH_IN', 1113000.00, '2026-10-14', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 14/10/2026', 'GRO-SSC-20261014', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('de63ddf0-3940-5a90-b7f6-0c50c4868dfe', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 200000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261014-11112', 'CASH_IN', 200000.00, '2026-10-14', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 14/10/2026', 'GRO-SSC-20261014', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4a27c217-dca2-5af7-b4b6-4cbf3b312e77', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1025000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5e22610c-95f7-5919-9a84-19c8217b0496', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fa893a58-a6ba-5060-81e8-e2647babf7a9', '13edc370-8cbd-5c8c-a54d-8f5a13df1d18', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 138000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261015 (SSC - 20261015): Rp 3,798,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('b1aeb1b2-484a-525f-8a8b-888c1a43e3e2', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261015', 'CASH_RECEIPT', '2026-10-15', '2026-10-15 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261015', 3798000.00, 3798000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 15-10-2026 Jesslyn Felicia Chendra, Feronika Ardella Meirizka, Salsabila Nabita, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('76acf6bf-33e4-545f-ad2d-3d6e90019ab7', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2270000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261015-11120', 'CASH_IN', 2270000.00, '2026-10-15', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 15/10/2026', 'GRO-SSC-20261015', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('611f3514-cef3-56a2-beae-2bae4421fd2d', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1528000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261015-11113', 'CASH_IN', 1528000.00, '2026-10-15', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 15/10/2026', 'GRO-SSC-20261015', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f72e5de3-e9e0-5cb3-9414-ab7bcbac1e2b', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 3390000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('34622838-e442-5fbf-87cc-818417b23e4c', 'b1aeb1b2-484a-525f-8a8b-888c1a43e3e2', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 408000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261016 (SSC - 20261016): Rp 2,068,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261016', 'CASH_RECEIPT', '2026-10-16', '2026-10-16 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261016', 2068500.00, 2068500.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 16-10-2026 Yoanita, Caesar Ayudhea, Ogie Hartantyo, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f82857cb-1bab-5513-900a-df7f370cb625', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 665000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261016-11113', 'CASH_IN', 665000.00, '2026-10-16', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 16/10/2026', 'GRO-SSC-20261016', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cb1ee82b-47ed-5644-afb0-84c545c5f3fd', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1403500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261016-11120', 'CASH_IN', 1403500.00, '2026-10-16', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 16/10/2026', 'GRO-SSC-20261016', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e0717c73-fb3b-54eb-83dd-d16464628b3f', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1805000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('471b689d-031b-5d75-82af-9e9357ef5caf', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 113500.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7f2e40b2-6c89-5f5d-ad64-c1084bb46a84', '1f15d401-cb79-52d6-ba2f-dfaed4bcdd03', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261017 (SSC - 20261017): Rp 3,169,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0e79e255-2941-5622-b954-a0421859ce8d', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261017', 'CASH_RECEIPT', '2026-10-17', '2026-10-17 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261017', 3169500.00, 3169500.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 17-10-2026 Susana, Grace Kimendora Gautama, Liem Yohanes Kevin Shanjaya, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a1607ffd-b1cf-52b0-b1e8-e7055ef220fd', '0e79e255-2941-5622-b954-a0421859ce8d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1604500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261017-11120', 'CASH_IN', 1604500.00, '2026-10-17', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 17/10/2026', 'GRO-SSC-20261017', '0e79e255-2941-5622-b954-a0421859ce8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('695bd9bf-02ab-5e63-af1c-9a9783875675', '0e79e255-2941-5622-b954-a0421859ce8d', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 220000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261017-11112', 'CASH_IN', 220000.00, '2026-10-17', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 17/10/2026', 'GRO-SSC-20261017', '0e79e255-2941-5622-b954-a0421859ce8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('35516f52-24c2-5c00-90b8-33c3192e6f78', '0e79e255-2941-5622-b954-a0421859ce8d', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1345000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261017-11113', 'CASH_IN', 1345000.00, '2026-10-17', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 17/10/2026', 'GRO-SSC-20261017', '0e79e255-2941-5622-b954-a0421859ce8d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2a7c4682-9969-5a1b-98e3-5935de4693fe', '0e79e255-2941-5622-b954-a0421859ce8d', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1589500.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b9ba74de-6fa5-5fee-b290-30d81cf686ae', '0e79e255-2941-5622-b954-a0421859ce8d', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 800000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cc3b1d2e-c412-5f0b-91ed-74f10a28b597', '0e79e255-2941-5622-b954-a0421859ce8d', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 200000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('10e7529b-8422-5dd9-bdd3-d90a85cb5e33', '0e79e255-2941-5622-b954-a0421859ce8d', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 580000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261018 (SSC - 20261018): Rp 1,390,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('499c2c10-44dd-5750-baaa-23b38c186db1', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261018', 'CASH_RECEIPT', '2026-10-18', '2026-10-18 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261018', 1390000.00, 1390000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 18-10-2026 Ong San Lien, Ina Ratnawati, Vina Santoso, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1db58e86-4074-58af-b3de-9d31dde831cf', '499c2c10-44dd-5750-baaa-23b38c186db1', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 665000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261018-11120', 'CASH_IN', 665000.00, '2026-10-18', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 18/10/2026', 'GRO-SSC-20261018', '499c2c10-44dd-5750-baaa-23b38c186db1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('14df8db0-80bf-5e8e-a5a9-34553db98612', '499c2c10-44dd-5750-baaa-23b38c186db1', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 725000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261018-11113', 'CASH_IN', 725000.00, '2026-10-18', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 18/10/2026', 'GRO-SSC-20261018', '499c2c10-44dd-5750-baaa-23b38c186db1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6145b004-bbb0-5f6a-bab5-34b25c48b739', '499c2c10-44dd-5750-baaa-23b38c186db1', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1390000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261019 (SSC - 20261019): Rp 522,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('155eff45-ff6f-5c97-8c46-ecaf0f426bc4', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261019', 'CASH_RECEIPT', '2026-10-19', '2026-10-19 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261019', 522000.00, 522000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 19-10-2026 Fakhrana Nabila, Yuhi Syaula (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('fb834929-e256-59a6-a86f-12c3c7ffebd5', '155eff45-ff6f-5c97-8c46-ecaf0f426bc4', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 522000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261019-11120', 'CASH_IN', 522000.00, '2026-10-19', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 19/10/2026', 'GRO-SSC-20261019', '155eff45-ff6f-5c97-8c46-ecaf0f426bc4')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('80b16729-b544-53a1-86bc-1ae4d8fcebfc', '155eff45-ff6f-5c97-8c46-ecaf0f426bc4', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 522000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261020 (SSC - 20261020): Rp 376,090.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7ae43d13-5ba7-5de5-8cb2-9fec03cc435c', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261020', 'CASH_RECEIPT', '2026-10-20', '2026-10-20 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261020', 376090.00, 376090.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 20-10-2026 Riska Dian Nadilla (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('491ccdab-7e0a-57a8-bab1-39e527af16a1', '7ae43d13-5ba7-5de5-8cb2-9fec03cc435c', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 376090.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261020-11112', 'CASH_IN', 376090.00, '2026-10-20', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 20/10/2026', 'GRO-SSC-20261020', '7ae43d13-5ba7-5de5-8cb2-9fec03cc435c')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5a7778f6-ff95-59a1-9f9f-b694d3e8b85e', '7ae43d13-5ba7-5de5-8cb2-9fec03cc435c', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 176090.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cac5c9ab-4d12-5065-9602-5a26db7595e1', '7ae43d13-5ba7-5de5-8cb2-9fec03cc435c', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 200000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261021 (SSC - 20261021): Rp 5,060,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3a4b0606-df0a-5fad-b8a7-610066c63e58', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261021', 'CASH_RECEIPT', '2026-10-21', '2026-10-21 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261021', 5060000.00, 5060000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 21-10-2026 Lia Kartikasari, Lie Lee Hong, Marcella Davina Sulistyani, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c1dbc113-c1f4-5780-aeb9-e04c66c71cd8', '3a4b0606-df0a-5fad-b8a7-610066c63e58', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1140000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261021-11120', 'CASH_IN', 1140000.00, '2026-10-21', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 21/10/2026', 'GRO-SSC-20261021', '3a4b0606-df0a-5fad-b8a7-610066c63e58')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('92363de0-2ebc-56c4-875b-c57a8d40bec1', '3a4b0606-df0a-5fad-b8a7-610066c63e58', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 3465000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261021-11113', 'CASH_IN', 3465000.00, '2026-10-21', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 21/10/2026', 'GRO-SSC-20261021', '3a4b0606-df0a-5fad-b8a7-610066c63e58')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f36bebdb-f383-5b4b-b398-dffdf2c053a9', '3a4b0606-df0a-5fad-b8a7-610066c63e58', 'fb5c3156-3a4b-4a86-a896-ec354caef246', 455000.00, 0.00, 'Penerimaan Kas/Bank via Uang Muka Pasien / Deposit Tindakan');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('257c6ca6-db1b-54a1-849f-ce9d933957b4', '3a4b0606-df0a-5fad-b8a7-610066c63e58', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2260000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9c0ced81-d37b-583a-8bf2-fd6a5106b72c', '3a4b0606-df0a-5fad-b8a7-610066c63e58', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 2800000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261022 (SSC - 20261022): Rp 6,775,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ad24cc3d-f23f-542f-98bb-3a95b54f8b02', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261022', 'CASH_RECEIPT', '2026-10-22', '2026-10-22 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261022', 6775000.00, 6775000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 22-10-2026 Hugo Hendarto, Nani Kristyawati, Dzakkiyatul Ulya Alfiannas Tito, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dbab374c-0748-5132-b9e5-688a294e7126', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4550000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261022-11120', 'CASH_IN', 4550000.00, '2026-10-22', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 22/10/2026', 'GRO-SSC-20261022', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('062107cf-11ab-5600-ba30-ad4ecc11da6d', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', 'fb5c3156-3a4b-4a86-a896-ec354caef246', 275000.00, 0.00, 'Penerimaan Kas/Bank via Uang Muka Pasien / Deposit Tindakan');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('637d6ec8-44a3-5444-bfca-7c2390c5644c', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 1950000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261022-11112', 'CASH_IN', 1950000.00, '2026-10-22', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 22/10/2026', 'GRO-SSC-20261022', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('44c735d3-db40-5cba-b683-9710815d30a2', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 825000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1e721e96-dd42-5425-bcee-55d6ec61e98f', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('33cde322-4b6d-58e2-8891-d113f628e51f', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 1600000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('80201788-f123-5311-b968-00ba7b1044da', 'ad24cc3d-f23f-542f-98bb-3a95b54f8b02', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 4000000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261023 (SSC - 20261023): Rp 665,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('14a33f32-5eb3-53cd-87eb-cfc2196828f3', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261023', 'CASH_RECEIPT', '2026-10-23', '2026-10-23 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261023', 665000.00, 665000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 23-10-2026 Amanda Pramesty (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('31e8b56b-b59b-5de1-b03c-924b14729df3', '14a33f32-5eb3-53cd-87eb-cfc2196828f3', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 665000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261023-11120', 'CASH_IN', 665000.00, '2026-10-23', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 23/10/2026', 'GRO-SSC-20261023', '14a33f32-5eb3-53cd-87eb-cfc2196828f3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('62ac70d8-bdef-568e-9d25-f2c4f432aea6', '14a33f32-5eb3-53cd-87eb-cfc2196828f3', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 665000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261024 (SSC - 20261024): Rp 2,055,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('5ec6d60d-4ee1-5457-a761-665efcaab1ab', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261024', 'CASH_RECEIPT', '2026-10-24', '2026-10-24 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261024', 2055000.00, 2055000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 24-10-2026 Dimas Bayu Suharno, Melisa Christinawai, Roseh Hartati, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b29be842-678d-500c-9c69-4f764a780c95', '5ec6d60d-4ee1-5457-a761-665efcaab1ab', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2055000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261024-11120', 'CASH_IN', 2055000.00, '2026-10-24', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 24/10/2026', 'GRO-SSC-20261024', '5ec6d60d-4ee1-5457-a761-665efcaab1ab')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5b89f547-e0fc-5913-863d-2fdfc057b822', '5ec6d60d-4ee1-5457-a761-665efcaab1ab', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2055000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261025 (SSC - 20261025): Rp 3,224,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c996785a-cf77-5190-ab3c-f495d90cd9e6', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261025', 'CASH_RECEIPT', '2026-10-25', '2026-10-25 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261025', 3224000.00, 3224000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 25-10-2026 Rizki Intan Wahyu Utami, D Nina Sartini, Lina Maulida, dkk. (total 8 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a2704c9a-3cf3-5c77-b576-31afe2722247', 'c996785a-cf77-5190-ab3c-f495d90cd9e6', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2974000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261025-11120', 'CASH_IN', 2974000.00, '2026-10-25', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 25/10/2026', 'GRO-SSC-20261025', 'c996785a-cf77-5190-ab3c-f495d90cd9e6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('90b1d632-95f5-52b5-b0ca-4c5babbee0f1', 'c996785a-cf77-5190-ab3c-f495d90cd9e6', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261025-11113', 'CASH_IN', 250000.00, '2026-10-25', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 25/10/2026', 'GRO-SSC-20261025', 'c996785a-cf77-5190-ab3c-f495d90cd9e6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('bad9361c-2165-500c-9f99-22629243b2a7', 'c996785a-cf77-5190-ab3c-f495d90cd9e6', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 3174000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2990a75e-39ba-51c2-91bc-393a91181fa2', 'c996785a-cf77-5190-ab3c-f495d90cd9e6', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 50000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261026 (SSC - 20261026): Rp 2,990,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0f9b5838-9514-5cce-b6ab-30d280be4cb4', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261026', 'CASH_RECEIPT', '2026-10-26', '2026-10-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261026', 2990000.00, 2990000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 26-10-2026 Niniek Lucianawati, Axel (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c0c89283-691b-53b0-a213-752f2e03b50a', '0f9b5838-9514-5cce-b6ab-30d280be4cb4', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1140000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261026-11120', 'CASH_IN', 1140000.00, '2026-10-26', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 26/10/2026', 'GRO-SSC-20261026', '0f9b5838-9514-5cce-b6ab-30d280be4cb4')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1f440019-58be-59ca-8301-8be8d8623d23', '0f9b5838-9514-5cce-b6ab-30d280be4cb4', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1850000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261026-11113', 'CASH_IN', 1850000.00, '2026-10-26', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 26/10/2026', 'GRO-SSC-20261026', '0f9b5838-9514-5cce-b6ab-30d280be4cb4')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('be50de31-95e3-53ff-9791-a243697d97e2', '0f9b5838-9514-5cce-b6ab-30d280be4cb4', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2990000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261027 (SSC - 20261027): Rp 7,045,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d769e824-62fe-55c3-a3c2-f5054455b8d3', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261027', 'CASH_RECEIPT', '2026-10-27', '2026-10-27 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261027', 7045000.00, 7045000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 27-10-2026 Ririn Puryanti, Jonathan Puguh Suhadi, Isnia Pratiwi, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('450c4099-b38d-5a86-ba5d-78fd3a441071', 'd769e824-62fe-55c3-a3c2-f5054455b8d3', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3014000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261027-11120', 'CASH_IN', 3014000.00, '2026-10-27', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 27/10/2026', 'GRO-SSC-20261027', 'd769e824-62fe-55c3-a3c2-f5054455b8d3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('839ace6f-b8ce-5489-8448-22a0d1dfb72a', 'd769e824-62fe-55c3-a3c2-f5054455b8d3', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 4031000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261027-11113', 'CASH_IN', 4031000.00, '2026-10-27', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 27/10/2026', 'GRO-SSC-20261027', 'd769e824-62fe-55c3-a3c2-f5054455b8d3')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6ee1729c-59e9-526c-a85d-f9bd5886f57d', 'd769e824-62fe-55c3-a3c2-f5054455b8d3', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2824000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8b7bd6c2-9788-58cb-a878-6b529e08a209', 'd769e824-62fe-55c3-a3c2-f5054455b8d3', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 2421000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1cc72d11-7cc9-5b2e-93ab-28a0754e2c87', 'd769e824-62fe-55c3-a3c2-f5054455b8d3', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 1800000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261028 (SSC - 20261028): Rp 4,620,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d9643b1e-82ce-570f-bf0c-339d248744e1', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261028', 'CASH_RECEIPT', '2026-10-28', '2026-10-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261028', 4620000.00, 4620000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 28-10-2026 Elsa Pithaloka Putri, Mindalena, Nabila Ananda Putri, dkk. (total 9 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('57e5f03a-144c-5807-8811-b5d45eb73d36', 'd9643b1e-82ce-570f-bf0c-339d248744e1', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3945000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261028-11120', 'CASH_IN', 3945000.00, '2026-10-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 28/10/2026', 'GRO-SSC-20261028', 'd9643b1e-82ce-570f-bf0c-339d248744e1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('53604d4f-bd19-552c-b9a7-85f5a2271850', 'd9643b1e-82ce-570f-bf0c-339d248744e1', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 275000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261028-11113', 'CASH_IN', 275000.00, '2026-10-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 28/10/2026', 'GRO-SSC-20261028', 'd9643b1e-82ce-570f-bf0c-339d248744e1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c5d2f98b-ab81-5726-a87e-f09ce7b79edf', 'd9643b1e-82ce-570f-bf0c-339d248744e1', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 400000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261028-11112', 'CASH_IN', 400000.00, '2026-10-28', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 28/10/2026', 'GRO-SSC-20261028', 'd9643b1e-82ce-570f-bf0c-339d248744e1')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('da0e05e8-f6d8-58ee-bfb1-1b8b13759980', 'd9643b1e-82ce-570f-bf0c-339d248744e1', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2615000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8842d7ef-cf52-5f35-9ca8-7cf4b9bdce16', 'd9643b1e-82ce-570f-bf0c-339d248744e1', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 150000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('19b6eaa5-192f-5b78-8b0a-de1787029f6c', 'd9643b1e-82ce-570f-bf0c-339d248744e1', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 1855000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261029 (SSC - 20261029): Rp 5,925,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261029', 'CASH_RECEIPT', '2026-10-29', '2026-10-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261029', 5925000.00, 5925000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 29-10-2026 Cynthia Carolina Santoso, Thea Permata, Sisilia Astuti, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cc227d8b-d0f3-5240-baaa-2bb4f0063cce', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5675000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261029-11120', 'CASH_IN', 5675000.00, '2026-10-29', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 29/10/2026', 'GRO-SSC-20261029', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8fb0a1d4-3d0a-5836-82ba-9bb7ef870163', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261029-11113', 'CASH_IN', 250000.00, '2026-10-29', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 29/10/2026', 'GRO-SSC-20261029', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d3ac964f-8440-56dd-9356-edf1cfc94ed5', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 4825000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('79ef235e-298f-5a21-9aba-c46d3d5c2dc7', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 900000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a36ca79d-8af0-5039-8dd9-8b4683bb51fa', 'd0d1d2b7-0e6a-521a-9467-ab15e5c62b2d', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 200000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261030 (SSC - 20261030): Rp 6,491,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3fd2073f-5512-5a6f-a4c8-0faebd0b1920', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261030', 'CASH_RECEIPT', '2026-10-30', '2026-10-30 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261030', 6491500.00, 6491500.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 30-10-2026 Feronika Ardella Meirizka, Isowedha Widya Dewi, Fadia Qistina Sidqi, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('dbbef95e-547f-5a84-b562-9690dd601f7f', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 5741500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261030-11120', 'CASH_IN', 5741500.00, '2026-10-30', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 30/10/2026', 'GRO-SSC-20261030', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('47f2cea5-0367-5ab1-9833-ebdcac3978eb', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 750000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261030-11112', 'CASH_IN', 750000.00, '2026-10-30', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 30/10/2026', 'GRO-SSC-20261030', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b8afbbea-b7bd-500a-bcc0-5bc89261ec02', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 4457000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e3bc363b-578b-5b1f-9317-6a23aeba0815', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 1634500.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a7050a02-3cd6-555f-9714-90a7c4c540ab', '3fd2073f-5512-5a6f-a4c8-0faebd0b1920', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 400000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261031 (SSC - 20261031): Rp 13,960,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('a128619b-2948-5933-8140-b8d870eb4619', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261031', 'CASH_RECEIPT', '2026-10-31', '2026-10-31 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261031', 13960000.00, 13960000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 31-10-2026 Marcia Ayu, Gladys Johan, Safa Maharashtr (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b6171931-365d-5c3e-8100-9eff369e3f07', 'a128619b-2948-5933-8140-b8d870eb4619', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 8820000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261031-11120', 'CASH_IN', 8820000.00, '2026-10-31', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 31/10/2026', 'GRO-SSC-20261031', 'a128619b-2948-5933-8140-b8d870eb4619')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7cb56266-591e-5c76-bb08-8a117dab4488', 'a128619b-2948-5933-8140-b8d870eb4619', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 5140000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261031-11113', 'CASH_IN', 5140000.00, '2026-10-31', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 31/10/2026', 'GRO-SSC-20261031', 'a128619b-2948-5933-8140-b8d870eb4619')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6ba29547-8629-5223-b171-ea54d04406b3', 'a128619b-2948-5933-8140-b8d870eb4619', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1120000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b57d28cd-4366-5680-9584-b797d3dbef17', 'a128619b-2948-5933-8140-b8d870eb4619', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 10325000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f1d2274d-ebf4-50a5-a942-6be9c3aa60bb', 'a128619b-2948-5933-8140-b8d870eb4619', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 2515000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261101 (SSC - 20261101): Rp 1,140,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('eb832e18-7414-5b31-a0e4-663c32126721', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261101', 'CASH_RECEIPT', '2026-11-01', '2026-11-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261101', 1140000.00, 1140000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 01-11-2026 Natalia Giofanny, Titis Dyah Hapsari, Aurelia Adistya Muljono (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7801c14e-9567-5f27-ac59-6525dd2a90ac', 'eb832e18-7414-5b31-a0e4-663c32126721', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1140000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261101-11120', 'CASH_IN', 1140000.00, '2026-11-01', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 01/11/2026', 'GRO-SSC-20261101', 'eb832e18-7414-5b31-a0e4-663c32126721')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('49db9e6d-04ad-5885-8647-ef26fae98401', 'eb832e18-7414-5b31-a0e4-663c32126721', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1140000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261102 (SSC - 20261102): Rp 4,530,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d7ec900d-b444-538d-9b94-efd0f47cef29', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261102', 'CASH_RECEIPT', '2026-11-02', '2026-11-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261102', 4530000.00, 4530000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 02-11-2026 Fellen, Diasa Ayu Raharni, Jessica Pratiknyo, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0e348f96-c203-5a29-8efb-55299e8f561e', 'd7ec900d-b444-538d-9b94-efd0f47cef29', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3730000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261102-11120', 'CASH_IN', 3730000.00, '2026-11-02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 02/11/2026', 'GRO-SSC-20261102', 'd7ec900d-b444-538d-9b94-efd0f47cef29')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('305b6e54-5454-5d3d-9fcb-a962c374680e', 'd7ec900d-b444-538d-9b94-efd0f47cef29', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 800000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261102-11113', 'CASH_IN', 800000.00, '2026-11-02', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 02/11/2026', 'GRO-SSC-20261102', 'd7ec900d-b444-538d-9b94-efd0f47cef29')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7dff649f-20ac-5f93-afbf-f1d1ff733f72', 'd7ec900d-b444-538d-9b94-efd0f47cef29', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 3030000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b1bbe5ba-1655-5c01-8b34-f43ec61e355b', 'd7ec900d-b444-538d-9b94-efd0f47cef29', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 1500000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261103 (SSC - 20261103): Rp 6,023,750.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0355b9c5-e431-5891-a8e2-8c474601a41b', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261103', 'CASH_RECEIPT', '2026-11-03', '2026-11-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261103', 6023750.00, 6023750.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 03-11-2026 Cemara Disa, Yuanita, Angela Valensia Siejaya, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6611f1cd-7bb0-5450-b08b-2c1859d69587', '0355b9c5-e431-5891-a8e2-8c474601a41b', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3550000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261103-11120', 'CASH_IN', 3550000.00, '2026-11-03', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 03/11/2026', 'GRO-SSC-20261103', '0355b9c5-e431-5891-a8e2-8c474601a41b')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('064cf707-d4c5-5974-9a1f-6feeda55a09f', '0355b9c5-e431-5891-a8e2-8c474601a41b', '69d4af2e-96d2-4856-bc88-da32eac1bc53', 2473750.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Semarang Skin Center (SSC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003', 'BKM-SSC-20261103-11112', 'CASH_IN', 2473750.00, '2026-11-03', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kas Kasir SSC) - 03/11/2026', 'GRO-SSC-20261103', '0355b9c5-e431-5891-a8e2-8c474601a41b')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5347c74d-2d1c-5080-8890-f37d639a2b6e', '0355b9c5-e431-5891-a8e2-8c474601a41b', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 4005000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('654b7dbc-6f01-5223-b062-decf927ca607', '0355b9c5-e431-5891-a8e2-8c474601a41b', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 2018750.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261104 (SSC - 20261104): Rp 1,280,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7a21a1cd-c209-5984-b7e2-e158cf0473fa', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261104', 'CASH_RECEIPT', '2026-11-04', '2026-11-04 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261104', 1280000.00, 1280000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 04-11-2026 Yunita Arumningsih, Soekarno Bastian, Christina Dewi Soesanto, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('808508f8-121b-527d-b87b-a79020c6fc12', '7a21a1cd-c209-5984-b7e2-e158cf0473fa', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 825000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261104-11120', 'CASH_IN', 825000.00, '2026-11-04', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 04/11/2026', 'GRO-SSC-20261104', '7a21a1cd-c209-5984-b7e2-e158cf0473fa')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('27421872-5f13-5ee2-bb2d-d18c7061e938', '7a21a1cd-c209-5984-b7e2-e158cf0473fa', 'fb5c3156-3a4b-4a86-a896-ec354caef246', 455000.00, 0.00, 'Penerimaan Kas/Bank via Uang Muka Pasien / Deposit Tindakan');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('51af6a4c-954e-5402-abbf-2499cc9244de', '7a21a1cd-c209-5984-b7e2-e158cf0473fa', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1280000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261105 (SSC - 20261105): Rp 1,105,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('e5ab9013-06a3-56f4-8edf-3c13f476ec33', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261105', 'CASH_RECEIPT', '2026-11-05', '2026-11-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261105', 1105000.00, 1105000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 05-11-2026 Tresniawati, Lie Lee Hong, Iin Viriany, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8dedbd76-c071-5e00-b374-a9c8935f970d', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 830000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261105-11113', 'CASH_IN', 830000.00, '2026-11-05', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 05/11/2026', 'GRO-SSC-20261105', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1f97bc3c-fdd2-5903-a4f6-567e4d2f1628', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 275000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261105-11120', 'CASH_IN', 275000.00, '2026-11-05', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 05/11/2026', 'GRO-SSC-20261105', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3a5714ff-530c-56c5-b24e-9f10cdce16e1', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1005000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a0ae7e18-98a2-5148-9e86-575f546a644b', 'e5ab9013-06a3-56f4-8edf-3c13f476ec33', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 100000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261106 (SSC - 20261106): Rp 2,550,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('fed08b23-3a94-512f-8ed5-8c8dc35388a5', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261106', 'CASH_RECEIPT', '2026-11-06', '2026-11-06 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261106', 2550000.00, 2550000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 06-11-2026 Lenawati Djojo Santoso, Melita Tetriani Darmawan, Kezia Dwi Sutjipto, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ce5b01e9-e0fa-5253-9283-a886735d7972', 'fed08b23-3a94-512f-8ed5-8c8dc35388a5', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2550000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261106-11120', 'CASH_IN', 2550000.00, '2026-11-06', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 06/11/2026', 'GRO-SSC-20261106', 'fed08b23-3a94-512f-8ed5-8c8dc35388a5')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ee9bfb2c-a51f-56e9-8327-2be96de51a47', 'fed08b23-3a94-512f-8ed5-8c8dc35388a5', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 2450000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('146db5d9-4ecd-5443-8701-5c5145da34b1', 'fed08b23-3a94-512f-8ed5-8c8dc35388a5', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 100000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261107 (SSC - 20261107): Rp 31,002,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('edb6c914-d56a-5230-a85b-01b19998c455', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261107', 'CASH_RECEIPT', '2026-11-07', '2026-11-07 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261107', 31002000.00, 31002000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 07-11-2026 Sudibyo Singgih, Erdina Adiwati, Aldiba Dodiek Tasan Wartono, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('63a59b66-3bcc-5002-9902-f8449cb0d7cc', 'edb6c914-d56a-5230-a85b-01b19998c455', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 29602000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261107-11120', 'CASH_IN', 29602000.00, '2026-11-07', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 07/11/2026', 'GRO-SSC-20261107', 'edb6c914-d56a-5230-a85b-01b19998c455')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('711e3242-29b3-5475-8c9a-f83e7ecfc9bc', 'edb6c914-d56a-5230-a85b-01b19998c455', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1400000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261107-11113', 'CASH_IN', 1400000.00, '2026-11-07', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 07/11/2026', 'GRO-SSC-20261107', 'edb6c914-d56a-5230-a85b-01b19998c455')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('61f7ec57-ef00-58eb-bd89-505307c1b417', 'edb6c914-d56a-5230-a85b-01b19998c455', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1705000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('46eadb0c-28c0-5ccc-9351-5bd0a1eac6e8', 'edb6c914-d56a-5230-a85b-01b19998c455', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 12437000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('376958b0-b2a5-5eb5-bf33-eecc84a511b7', 'edb6c914-d56a-5230-a85b-01b19998c455', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 16560000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8e4aa7cd-fa36-5de9-a137-92a126d05a2d', 'edb6c914-d56a-5230-a85b-01b19998c455', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 300000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    -- Journal CR-SSC-20261108 (SSC - 20261108): Rp 346,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('fcf3f267-b2f0-5826-8dcc-74d0dd567108', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261108', 'CASH_RECEIPT', '2026-11-08', '2026-11-08 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261108', 346000.00, 346000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 08-11-2026 Defandryan Elzan Alfarizky (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9f1218ad-2aee-580f-9736-9e1234e01ccf', 'fcf3f267-b2f0-5826-8dcc-74d0dd567108', 'fb5c3156-3a4b-4a86-a896-ec354caef246', 346000.00, 0.00, 'Penerimaan Kas/Bank via Uang Muka Pasien / Deposit Tindakan');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e74d96d6-1a28-52a7-97dd-9380865da906', 'fcf3f267-b2f0-5826-8dcc-74d0dd567108', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 100000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c1bd4dfc-6fc1-536c-b7a9-3f23d94ae6e1', 'fcf3f267-b2f0-5826-8dcc-74d0dd567108', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 246000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261109 (SSC - 20261109): Rp 1,630,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ea8bc6b3-0b44-5a41-b9a7-64d956f1a0be', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261109', 'CASH_RECEIPT', '2026-11-09', '2026-11-09 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261109', 1630000.00, 1630000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 09-11-2026 Lidiawati Surjanto, Tony Prajogo, Carmenita Budihartono, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4bdda223-c618-5d33-941d-0199df91758f', 'ea8bc6b3-0b44-5a41-b9a7-64d956f1a0be', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1630000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261109-11113', 'CASH_IN', 1630000.00, '2026-11-09', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 09/11/2026', 'GRO-SSC-20261109', 'ea8bc6b3-0b44-5a41-b9a7-64d956f1a0be')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1ce95c50-d89e-5f7d-bbe1-9c3b868484fd', 'ea8bc6b3-0b44-5a41-b9a7-64d956f1a0be', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1510000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b9e52f3d-82b9-5b67-9ba7-150351b6db81', 'ea8bc6b3-0b44-5a41-b9a7-64d956f1a0be', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 120000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261110 (SSC - 20261110): Rp 622,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('a4278a76-d70b-51e7-bdc2-8d58148cb5f9', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261110', 'CASH_RECEIPT', '2026-11-10', '2026-11-10 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261110', 622500.00, 622500.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 10-11-2026 Sally Dwi Purnamasari (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7fae573e-3111-5a43-a0d9-fd76e8e8c4de', 'a4278a76-d70b-51e7-bdc2-8d58148cb5f9', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 622500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261110-11120', 'CASH_IN', 622500.00, '2026-11-10', '63585f66-7ed5-4e58-954a-689881e7e186', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 10/11/2026', 'GRO-SSC-20261110', 'a4278a76-d70b-51e7-bdc2-8d58148cb5f9')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6a4576cf-e82b-586a-8247-c8bd55866563', 'a4278a76-d70b-51e7-bdc2-8d58148cb5f9', '63585f66-7ed5-4e58-954a-689881e7e186', 0.00, 100000.00, 'Pendapatan Operasional Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8ac08873-c28b-5ecc-8756-cf6c1478ca20', 'a4278a76-d70b-51e7-bdc2-8d58148cb5f9', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 225000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('40c7b350-5256-598f-b0ce-4455716e56eb', 'a4278a76-d70b-51e7-bdc2-8d58148cb5f9', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 297500.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261111 (SSC - 20261111): Rp 2,717,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('e08fb827-ab27-5d5b-a467-5731e3268e33', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261111', 'CASH_RECEIPT', '2026-11-11', '2026-11-11 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261111', 2717000.00, 2717000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 11-11-2026 Adelia Bayu Isfandiari, Devina Wirawan (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a101359a-dc45-5d6f-8097-192648ef31e9', 'e08fb827-ab27-5d5b-a467-5731e3268e33', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 2717000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261111-11120', 'CASH_IN', 2717000.00, '2026-11-11', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 11/11/2026', 'GRO-SSC-20261111', 'e08fb827-ab27-5d5b-a467-5731e3268e33')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('066c495d-0083-56d4-b141-22f65c69b218', 'e08fb827-ab27-5d5b-a467-5731e3268e33', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 1237000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2bee1499-9b55-5171-ab8e-c2486f5b1c4d', 'e08fb827-ab27-5d5b-a467-5731e3268e33', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 1440000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6b150f6d-2fb1-56d9-80eb-7f084d34efcb', 'e08fb827-ab27-5d5b-a467-5731e3268e33', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 40000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    -- Journal CR-SSC-20261112 (SSC - 20261112): Rp 1,050,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7bed86db-a583-53b9-80d4-125b2b98443d', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261112', 'CASH_RECEIPT', '2026-11-12', '2026-11-12 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261112', 1050000.00, 1050000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 12-11-2026 Cynthia Carolina Santoso, Agus Sutrisno, Roni Setyaji, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d9c64acc-4c05-5402-ac65-2fe64d96a90f', '7bed86db-a583-53b9-80d4-125b2b98443d', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1050000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261112-11120', 'CASH_IN', 1050000.00, '2026-11-12', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 12/11/2026', 'GRO-SSC-20261112', '7bed86db-a583-53b9-80d4-125b2b98443d')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a1b3b8fe-7f1b-502f-9a12-8c9af99ceb8a', '7bed86db-a583-53b9-80d4-125b2b98443d', '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', 0.00, 250000.00, 'Pendapatan Operasional Pendapatan Penjualan Produk Skincare & Sunscreen');
            

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('992560c6-439f-538b-9778-f17a7f59cb3e', '7bed86db-a583-53b9-80d4-125b2b98443d', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 800000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261113 (SSC - 20261113): Rp 455,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9cf611bf-6f5b-57ad-9d63-c1c22e310676', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261113', 'CASH_RECEIPT', '2026-11-13', '2026-11-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261113', 455000.00, 455000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 13-11-2026 Devina Chandra (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d5c25de7-1bf3-58df-8c1c-64c1c3f772ef', '9cf611bf-6f5b-57ad-9d63-c1c22e310676', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 455000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261113-11113', 'CASH_IN', 455000.00, '2026-11-13', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 13/11/2026', 'GRO-SSC-20261113', '9cf611bf-6f5b-57ad-9d63-c1c22e310676')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7500a0f9-4bc4-5056-a847-d5d90338c449', '9cf611bf-6f5b-57ad-9d63-c1c22e310676', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 455000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261114 (SSC - 20261114): Rp 6,900,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('3025eb29-49a0-5bb5-8029-195189467ca4', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261114', 'CASH_RECEIPT', '2026-11-14', '2026-11-14 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261114', 6900000.00, 6900000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 14-11-2026 Erdina Adiwati (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('30be9eed-a6a6-5454-9d8b-9e5d54449eef', '3025eb29-49a0-5bb5-8029-195189467ca4', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 6900000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261114-11120', 'CASH_IN', 6900000.00, '2026-11-14', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 14/11/2026', 'GRO-SSC-20261114', '3025eb29-49a0-5bb5-8029-195189467ca4')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('16a07117-9a1c-57e4-9956-a8ce7ddb2055', '3025eb29-49a0-5bb5-8029-195189467ca4', '13f34cd9-557c-42ab-b0d8-09056fd2988f', 0.00, 6900000.00, 'Pendapatan Operasional Pendapatan Tindakan Laser, Cauter & Peeling');
            

    -- Journal CR-SSC-20261115 (SSC - 20261115): Rp 275,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('88d94798-2b4c-580a-9c7e-95a457f5812a', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261115', 'CASH_RECEIPT', '2026-11-15', '2026-11-15 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261115', 275000.00, 275000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 15-11-2026 Jo Ananta Yoshua Surya (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9a8fc297-6dd4-55ce-aab6-585113307f1c', '88d94798-2b4c-580a-9c7e-95a457f5812a', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 275000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000004', 'BKM-SSC-20261115-11113', 'CASH_IN', 275000.00, '2026-11-15', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Bank BCA Operasional) - 15/11/2026', 'GRO-SSC-20261115', '88d94798-2b4c-580a-9c7e-95a457f5812a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('74717a1f-b2e5-56ce-b1f5-d4af4757a41a', '88d94798-2b4c-580a-9c7e-95a457f5812a', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 275000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261116 (SSC - 20261116): Rp 4,295,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('1f704272-2ea5-5d93-887e-3aded17dcb94', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261116', 'CASH_RECEIPT', '2026-11-16', '2026-11-16 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261116', 4295000.00, 4295000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 16-11-2026 Upik Nur Aprilia, Susi Herawati, Noviandi Hutagaol, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2e350d5d-347f-58c0-8d9c-750860ae33aa', '1f704272-2ea5-5d93-887e-3aded17dcb94', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 4295000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261116-11120', 'CASH_IN', 4295000.00, '2026-11-16', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 16/11/2026', 'GRO-SSC-20261116', '1f704272-2ea5-5d93-887e-3aded17dcb94')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('15d1bc1a-4799-5403-a449-ebf3d3df3c06', '1f704272-2ea5-5d93-887e-3aded17dcb94', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 4295000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261117 (SSC - 20261117): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d2eee44b-3b6e-54f5-9918-a8410315e1df', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261117', 'CASH_RECEIPT', '2026-11-17', '2026-11-17 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261117', 175000.00, 175000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 17-11-2026 Diasa Ayu Raharni (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8d67907f-3486-55d3-89be-8ceecc940558', 'd2eee44b-3b6e-54f5-9918-a8410315e1df', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261117-11120', 'CASH_IN', 175000.00, '2026-11-17', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 17/11/2026', 'GRO-SSC-20261117', 'd2eee44b-3b6e-54f5-9918-a8410315e1df')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('41227c61-2f1f-5ae9-a41d-0bb53c0daf95', 'd2eee44b-3b6e-54f5-9918-a8410315e1df', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-SSC-20261118 (SSC - 20261118): Rp 665,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ac52dab0-eff9-5864-82b5-53defe3fc436', 'a0000000-0000-0000-0000-000000000002', 'CR-SSC-20261118', 'CASH_RECEIPT', '2026-11-18', '2026-11-18 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-SSC-20261118', 665000.00, 665000.00, 'Penerimaan Kasir Harian [SSC] - Tanggal 18-11-2026 Dessy Setyaningrum (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ec4d0e5c-3711-5c61-8919-b3cd5e2475cf', 'ac52dab0-eff9-5864-82b5-53defe3fc436', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 665000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000006', 'BKM-SSC-20261118-11120', 'CASH_IN', 665000.00, '2026-11-18', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 'Penerimaan Kasir Harian SSC (Kliring Payment Gateway & EDC) - 18/11/2026', 'GRO-SSC-20261118', 'ac52dab0-eff9-5864-82b5-53defe3fc436')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('efaaedfe-eb35-596e-898b-2130b853eb74', 'ac52dab0-eff9-5864-82b5-53defe3fc436', 'eab0bf83-ae3b-4ea5-a320-3faa34488032', 0.00, 665000.00, 'Pendapatan Operasional Pendapatan Treatment Facial & Estetika Medis');
            

    -- Journal CR-MNM-20260926 (MINIMO - 20260926): Rp 87,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('6557fb7d-f01c-5008-abae-5516ac44a8b6', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20260926', 'CASH_RECEIPT', '2026-09-26', '2026-09-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20260926', 87500.00, 87500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 26-09-2026 Clarissa Olivia (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('944e07c1-90e7-5cd2-a13c-5425659dbbd4', '6557fb7d-f01c-5008-abae-5516ac44a8b6', '2a3ad241-587d-4853-bdc2-615d72496b4e', 87500.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20260926-11111', 'CASH_IN', 87500.00, '2026-09-26', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 26/09/2026', 'GRO-MINIMO-20260926', '6557fb7d-f01c-5008-abae-5516ac44a8b6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('efebc285-6fa6-5c45-9b78-98e132d903ec', '6557fb7d-f01c-5008-abae-5516ac44a8b6', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 87500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20260928 (MINIMO - 20260928): Rp 375,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('08450f37-ced7-5d58-919e-097521625747', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20260928', 'CASH_RECEIPT', '2026-09-28', '2026-09-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20260928', 375000.00, 375000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 28-09-2026 Elon Savio Budiono, Ethan Filbert Budiono, Nael Zevran Arkhairym (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('591fee07-512b-5dfb-a739-8290c89bc6fe', '08450f37-ced7-5d58-919e-097521625747', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 375000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20260928-11120', 'CASH_IN', 375000.00, '2026-09-28', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 28/09/2026', 'GRO-MINIMO-20260928', '08450f37-ced7-5d58-919e-097521625747')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b5f94291-e55d-59dd-af26-87940583b7d2', '08450f37-ced7-5d58-919e-097521625747', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 375000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20260929 (MINIMO - 20260929): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('5b9a7162-f00a-5267-9d2a-f5467288db4a', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20260929', 'CASH_RECEIPT', '2026-09-29', '2026-09-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20260929', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 29-09-2026 Ellora Joyceline Handoyo (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2ebdec2a-ee8b-569d-b855-d7f9049b9089', '5b9a7162-f00a-5267-9d2a-f5467288db4a', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 87500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20260929-11120', 'CASH_IN', 87500.00, '2026-09-29', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 29/09/2026', 'GRO-MINIMO-20260929', '5b9a7162-f00a-5267-9d2a-f5467288db4a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('46f2fbba-6a41-5611-a670-c1bc818468bf', '5b9a7162-f00a-5267-9d2a-f5467288db4a', '2a3ad241-587d-4853-bdc2-615d72496b4e', 87500.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20260929-11111', 'CASH_IN', 87500.00, '2026-09-29', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 29/09/2026', 'GRO-MINIMO-20260929', '5b9a7162-f00a-5267-9d2a-f5467288db4a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ecaed656-d85c-572b-8361-a80544f7c389', '5b9a7162-f00a-5267-9d2a-f5467288db4a', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20260930 (MINIMO - 20260930): Rp 185,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7f25a5c0-010b-55b1-ba0e-e36d7e88ad73', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20260930', 'CASH_RECEIPT', '2026-09-30', '2026-09-30 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20260930', 185000.00, 185000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 30-09-2026 Bianca Gracella Handojo, Fabiola Celina Susanto (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('185cd7dc-3553-5ef4-a069-40c244cc0c23', '7f25a5c0-010b-55b1-ba0e-e36d7e88ad73', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 185000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20260930-11120', 'CASH_IN', 185000.00, '2026-09-30', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 30/09/2026', 'GRO-MINIMO-20260930', '7f25a5c0-010b-55b1-ba0e-e36d7e88ad73')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('114a743a-00b3-57b1-af8b-02e514f06b88', '7f25a5c0-010b-55b1-ba0e-e36d7e88ad73', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 185000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261001 (MINIMO - 20261001): Rp 925,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('07e704e8-22a9-5da4-bb87-0dc0a170ddf0', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261001', 'CASH_RECEIPT', '2026-10-01', '2026-10-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261001', 925000.00, 925000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 01-10-2026 Audrey Darielle, Levin Arlo Lesmana (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('db5dde6d-36b4-5931-8650-842986b7dcce', '07e704e8-22a9-5da4-bb87-0dc0a170ddf0', '2a3ad241-587d-4853-bdc2-615d72496b4e', 925000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261001-11111', 'CASH_IN', 925000.00, '2026-10-01', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 01/10/2026', 'GRO-MINIMO-20261001', '07e704e8-22a9-5da4-bb87-0dc0a170ddf0')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3ac6b690-83fa-5412-8ea8-07d3ad90efa5', '07e704e8-22a9-5da4-bb87-0dc0a170ddf0', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 925000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261002 (MINIMO - 20261002): Rp 3,010,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('f4d8528b-a70e-5095-bd1a-08a111f3bfcd', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261002', 'CASH_RECEIPT', '2026-10-02', '2026-10-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261002', 3010000.00, 3010000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 02-10-2026 Odilia Claryne Auberta, Kaysan Muhammad Andiansyah, Marshall Calyson Sadara, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('78694043-70b0-5c53-8b50-32fd8889b765', 'f4d8528b-a70e-5095-bd1a-08a111f3bfcd', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1610000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261002-11120', 'CASH_IN', 1610000.00, '2026-10-02', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 02/10/2026', 'GRO-MINIMO-20261002', 'f4d8528b-a70e-5095-bd1a-08a111f3bfcd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d056013d-00c4-59ca-89aa-be820a478a49', 'f4d8528b-a70e-5095-bd1a-08a111f3bfcd', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1400000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261002-11113', 'CASH_IN', 1400000.00, '2026-10-02', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 02/10/2026', 'GRO-MINIMO-20261002', 'f4d8528b-a70e-5095-bd1a-08a111f3bfcd')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('38c949ac-4779-5eb8-8634-9c793140bd4f', 'f4d8528b-a70e-5095-bd1a-08a111f3bfcd', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 3010000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261003 (MINIMO - 20261003): Rp 5,162,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('54147fe5-c4aa-5bf6-a708-fab598343330', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261003', 'CASH_RECEIPT', '2026-10-03', '2026-10-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261003', 5162500.00, 5162500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 03-10-2026 Daphne Fumiko Wijaya, Syra Aurelia, Chester Power Hartono, dkk. (total 10 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1e6b60ee-b0a4-58b2-b3d7-0414e38dcac6', '54147fe5-c4aa-5bf6-a708-fab598343330', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 3012500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261003-11120', 'CASH_IN', 3012500.00, '2026-10-03', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 03/10/2026', 'GRO-MINIMO-20261003', '54147fe5-c4aa-5bf6-a708-fab598343330')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9d4906ba-c5cd-5b13-9591-309e937c97ed', '54147fe5-c4aa-5bf6-a708-fab598343330', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2150000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261003-11113', 'CASH_IN', 2150000.00, '2026-10-03', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 03/10/2026', 'GRO-MINIMO-20261003', '54147fe5-c4aa-5bf6-a708-fab598343330')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d76b9939-3582-5077-934b-94d2219610e2', '54147fe5-c4aa-5bf6-a708-fab598343330', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 5162500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261005 (MINIMO - 20261005): Rp 2,425,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('a85217c0-51a4-59d3-9bc0-fa39105f6e28', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261005', 'CASH_RECEIPT', '2026-10-05', '2026-10-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261005', 2425000.00, 2425000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 05-10-2026 Rumaisha Safa Rahmindra, Rayyanza Salman Rahmandra, Immanuel David Wijaya, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2603d26e-3207-5d52-9747-46b267a85b69', 'a85217c0-51a4-59d3-9bc0-fa39105f6e28', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 2250000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261005-11113', 'CASH_IN', 2250000.00, '2026-10-05', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 05/10/2026', 'GRO-MINIMO-20261005', 'a85217c0-51a4-59d3-9bc0-fa39105f6e28')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d63d22c7-535c-56bd-ad98-9b376d384672', 'a85217c0-51a4-59d3-9bc0-fa39105f6e28', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261005-11120', 'CASH_IN', 175000.00, '2026-10-05', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 05/10/2026', 'GRO-MINIMO-20261005', 'a85217c0-51a4-59d3-9bc0-fa39105f6e28')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4839cd2c-bf3b-50d8-83ca-d70bdf45902b', 'a85217c0-51a4-59d3-9bc0-fa39105f6e28', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 2425000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261006 (MINIMO - 20261006): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('42e18707-feba-5f79-8b27-abf3f8bb1da7', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261006', 'CASH_RECEIPT', '2026-10-06', '2026-10-06 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261006', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 06-10-2026 Najib Muhammad Husain (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('441aca49-115e-51a2-8075-557450e5668a', '42e18707-feba-5f79-8b27-abf3f8bb1da7', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261006-11120', 'CASH_IN', 175000.00, '2026-10-06', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 06/10/2026', 'GRO-MINIMO-20261006', '42e18707-feba-5f79-8b27-abf3f8bb1da7')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('675fe82c-4fcf-54ab-93a6-1c5a76476b50', '42e18707-feba-5f79-8b27-abf3f8bb1da7', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261007 (MINIMO - 20261007): Rp 350,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('5358205c-db3a-5acc-ae27-03707a008dce', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261007', 'CASH_RECEIPT', '2026-10-07', '2026-10-07 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261007', 350000.00, 350000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 07-10-2026 Samantha Abigail Saputra, Kyra Thalassa Shanum (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d05bc373-61a4-5563-9674-90dab45bbb65', '5358205c-db3a-5acc-ae27-03707a008dce', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261007-11120', 'CASH_IN', 350000.00, '2026-10-07', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 07/10/2026', 'GRO-MINIMO-20261007', '5358205c-db3a-5acc-ae27-03707a008dce')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4510a2fc-4ccf-5ba3-8fd2-8f87b330ef9e', '5358205c-db3a-5acc-ae27-03707a008dce', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261008 (MINIMO - 20261008): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('dc21986c-66b3-57fb-8dc6-4731b54922b9', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261008', 'CASH_RECEIPT', '2026-10-08', '2026-10-08 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261008', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 08-10-2026 Muhammad Beryl Rasendria (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ac96b5e3-b500-5062-a145-2f91e3e08536', 'dc21986c-66b3-57fb-8dc6-4731b54922b9', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261008-11120', 'CASH_IN', 175000.00, '2026-10-08', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 08/10/2026', 'GRO-MINIMO-20261008', 'dc21986c-66b3-57fb-8dc6-4731b54922b9')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('de2944c2-be39-5aca-8d4d-0998bac74380', 'dc21986c-66b3-57fb-8dc6-4731b54922b9', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261009 (MINIMO - 20261009): Rp 350,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0784800f-6969-5736-8b53-c0d5006e0970', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261009', 'CASH_RECEIPT', '2026-10-09', '2026-10-09 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261009', 350000.00, 350000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 09-10-2026 Elmarshad Rais Rafazka, Aldrich Ramos Sinaga (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1758f8cb-f2c6-53a5-b646-413eb9d54a50', '0784800f-6969-5736-8b53-c0d5006e0970', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261009-11120', 'CASH_IN', 350000.00, '2026-10-09', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 09/10/2026', 'GRO-MINIMO-20261009', '0784800f-6969-5736-8b53-c0d5006e0970')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4836201f-51c5-5401-b6d6-8c6be2873211', '0784800f-6969-5736-8b53-c0d5006e0970', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261010 (MINIMO - 20261010): Rp 525,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('1bce06d1-520c-56ae-aec9-39f7ddb58581', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261010', 'CASH_RECEIPT', '2026-10-10', '2026-10-10 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261010', 525000.00, 525000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 10-10-2026 Elyandra Vellantia (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3faf0349-50b6-51f2-b965-57c1be7583d0', '1bce06d1-520c-56ae-aec9-39f7ddb58581', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 525000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261010-11113', 'CASH_IN', 525000.00, '2026-10-10', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 10/10/2026', 'GRO-MINIMO-20261010', '1bce06d1-520c-56ae-aec9-39f7ddb58581')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f1c5b763-10f0-5f79-84d5-81617f43440f', '1bce06d1-520c-56ae-aec9-39f7ddb58581', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 525000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261013 (MINIMO - 20261013): Rp 875,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('aba18708-6035-5de2-9a5a-fb7371229951', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261013', 'CASH_RECEIPT', '2026-10-13', '2026-10-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261013', 875000.00, 875000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 13-10-2026 Qiao Jinchuan, Shaziva Zerryn Khumaira, Benjamin Nathanael Sarwono, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b1ecc06f-bb5a-58bb-ac4c-b4f3ea48f5da', 'aba18708-6035-5de2-9a5a-fb7371229951', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261013-11120', 'CASH_IN', 350000.00, '2026-10-13', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 13/10/2026', 'GRO-MINIMO-20261013', 'aba18708-6035-5de2-9a5a-fb7371229951')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e10fe5d8-5f07-5fb7-8d1b-b7bfc6447e01', 'aba18708-6035-5de2-9a5a-fb7371229951', '2a3ad241-587d-4853-bdc2-615d72496b4e', 525000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261013-11111', 'CASH_IN', 525000.00, '2026-10-13', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 13/10/2026', 'GRO-MINIMO-20261013', 'aba18708-6035-5de2-9a5a-fb7371229951')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('75965054-b875-5471-9a7c-10fac197762a', 'aba18708-6035-5de2-9a5a-fb7371229951', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 875000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261014 (MINIMO - 20261014): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('508d5288-503f-56eb-a8d4-76a826548d9e', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261014', 'CASH_RECEIPT', '2026-10-14', '2026-10-14 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261014', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 14-10-2026 Xian Fei Danli (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('914f7e31-357a-57c3-87d5-507b6197f528', '508d5288-503f-56eb-a8d4-76a826548d9e', '2a3ad241-587d-4853-bdc2-615d72496b4e', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261014-11111', 'CASH_IN', 175000.00, '2026-10-14', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 14/10/2026', 'GRO-MINIMO-20261014', '508d5288-503f-56eb-a8d4-76a826548d9e')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cb815183-864f-554f-a74e-58728e81f08d', '508d5288-503f-56eb-a8d4-76a826548d9e', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261015 (MINIMO - 20261015): Rp 350,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('cb539c02-3d63-59bf-9c78-e004dc19108b', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261015', 'CASH_RECEIPT', '2026-10-15', '2026-10-15 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261015', 350000.00, 350000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 15-10-2026 Bennett Karuna Hartaji, Jennifer Eliora Saputro (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c5fd974b-32e8-596b-94b6-43ede3ba9eb6', 'cb539c02-3d63-59bf-9c78-e004dc19108b', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261015-11120', 'CASH_IN', 350000.00, '2026-10-15', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 15/10/2026', 'GRO-MINIMO-20261015', 'cb539c02-3d63-59bf-9c78-e004dc19108b')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c136c24b-d0ce-5baa-bd9a-9bde2d8beb67', 'cb539c02-3d63-59bf-9c78-e004dc19108b', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261016 (MINIMO - 20261016): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('72bba667-b276-5fc1-b870-4c55662c6d95', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261016', 'CASH_RECEIPT', '2026-10-16', '2026-10-16 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261016', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 16-10-2026 Reyner Nathanuel Wirawan, Joshua Ethanuel Wirawan (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('2a01aef2-83d8-5777-9ebf-b707bb20b3e7', '72bba667-b276-5fc1-b870-4c55662c6d95', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261016-11120', 'CASH_IN', 175000.00, '2026-10-16', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 16/10/2026', 'GRO-MINIMO-20261016', '72bba667-b276-5fc1-b870-4c55662c6d95')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d07c0a96-cb97-581d-9a44-4f6466395506', '72bba667-b276-5fc1-b870-4c55662c6d95', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261019 (MINIMO - 20261019): Rp 875,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('544ee0d3-5533-5ff7-9903-cd1789819154', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261019', 'CASH_RECEIPT', '2026-10-19', '2026-10-19 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261019', 875000.00, 875000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 19-10-2026 Kenny Arsa Maleek (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('3d2177c7-6d65-51a0-bb54-5790e90f3ce1', '544ee0d3-5533-5ff7-9903-cd1789819154', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 875000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261019-11120', 'CASH_IN', 875000.00, '2026-10-19', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 19/10/2026', 'GRO-MINIMO-20261019', '544ee0d3-5533-5ff7-9903-cd1789819154')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('23c220b0-397f-529f-a538-d0f57d8ad1f3', '544ee0d3-5533-5ff7-9903-cd1789819154', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 875000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261021 (MINIMO - 20261021): Rp 612,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('0a0c6b6d-d361-5ad0-953a-9b1900a5a10a', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261021', 'CASH_RECEIPT', '2026-10-21', '2026-10-21 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261021', 612500.00, 612500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 21-10-2026 Samantha Abigail Saputra, Bernadette Clousia (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c947aec0-a2a0-5984-88c6-79ea622a5518', '0a0c6b6d-d361-5ad0-953a-9b1900a5a10a', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 612500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261021-11113', 'CASH_IN', 612500.00, '2026-10-21', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 21/10/2026', 'GRO-MINIMO-20261021', '0a0c6b6d-d361-5ad0-953a-9b1900a5a10a')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b8834752-1624-5984-b1f8-0111e606f232', '0a0c6b6d-d361-5ad0-953a-9b1900a5a10a', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 612500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261022 (MINIMO - 20261022): Rp 1,400,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('ce5d07e9-6def-5a28-8654-7d70f0514354', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261022', 'CASH_RECEIPT', '2026-10-22', '2026-10-22 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261022', 1400000.00, 1400000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 22-10-2026 Maheera Naradynandri Ansara, Ellora Lavanya Shaqeela (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7e34ff4b-65cd-5b15-a611-2312f5f1030d', 'ce5d07e9-6def-5a28-8654-7d70f0514354', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1400000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261022-11120', 'CASH_IN', 1400000.00, '2026-10-22', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 22/10/2026', 'GRO-MINIMO-20261022', 'ce5d07e9-6def-5a28-8654-7d70f0514354')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4ad691ec-b412-51ed-9252-c47656d5d4c4', 'ce5d07e9-6def-5a28-8654-7d70f0514354', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1400000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261023 (MINIMO - 20261023): Rp 262,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('6f1313b5-0539-5a26-91f3-e28b280ade8f', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261023', 'CASH_RECEIPT', '2026-10-23', '2026-10-23 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261023', 262500.00, 262500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 23-10-2026 Salmaputri Novitasari, Maureen Kana Wijaya, Syahdanputro Nugroho (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c7927e03-d930-539b-bcec-21b8305cc99c', '6f1313b5-0539-5a26-91f3-e28b280ade8f', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 87500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261023-11113', 'CASH_IN', 87500.00, '2026-10-23', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 23/10/2026', 'GRO-MINIMO-20261023', '6f1313b5-0539-5a26-91f3-e28b280ade8f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('cfb1164f-2737-54b6-9b41-18a232f4d6e9', '6f1313b5-0539-5a26-91f3-e28b280ade8f', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261023-11120', 'CASH_IN', 175000.00, '2026-10-23', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 23/10/2026', 'GRO-MINIMO-20261023', '6f1313b5-0539-5a26-91f3-e28b280ade8f')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('271d12d0-99d0-5a70-9856-bd046f65af41', '6f1313b5-0539-5a26-91f3-e28b280ade8f', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 262500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261024 (MINIMO - 20261024): Rp 700,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('8c418fce-9f93-5282-a9cf-0ad1ef02f8bf', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261024', 'CASH_RECEIPT', '2026-10-24', '2026-10-24 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261024', 700000.00, 700000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 24-10-2026 Jeremiah Russel, Kendrick Oliver Candra, Yuna Eleora Russlyn, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('d04c5876-f02a-5f26-9507-7c2c4de4f027', '8c418fce-9f93-5282-a9cf-0ad1ef02f8bf', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 700000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261024-11120', 'CASH_IN', 700000.00, '2026-10-24', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 24/10/2026', 'GRO-MINIMO-20261024', '8c418fce-9f93-5282-a9cf-0ad1ef02f8bf')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f3982f18-471c-5e53-9961-2538ca7d634c', '8c418fce-9f93-5282-a9cf-0ad1ef02f8bf', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 700000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261025 (MINIMO - 20261025): Rp 350,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('432076f5-7127-53d5-b5fc-81ccc6b13859', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261025', 'CASH_RECEIPT', '2026-10-25', '2026-10-25 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261025', 350000.00, 350000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 25-10-2026 Nael Zevran Arkhairym (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('31f6c00f-89bf-55f1-82d9-1ef6dafb0121', '432076f5-7127-53d5-b5fc-81ccc6b13859', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261025-11120', 'CASH_IN', 350000.00, '2026-10-25', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 25/10/2026', 'GRO-MINIMO-20261025', '432076f5-7127-53d5-b5fc-81ccc6b13859')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('9a8b8f81-d8d0-5f96-94fb-b89e0d8a371c', '432076f5-7127-53d5-b5fc-81ccc6b13859', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261026 (MINIMO - 20261026): Rp 1,060,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('e99f0958-e81d-577c-b9b0-0ed158ba57e8', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261026', 'CASH_RECEIPT', '2026-10-26', '2026-10-26 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261026', 1060000.00, 1060000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 26-10-2026 Ghazzal Rivo (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('8cb5e6e9-af9d-5228-ba67-ae5ce3adc88c', 'e99f0958-e81d-577c-b9b0-0ed158ba57e8', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1060000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261026-11113', 'CASH_IN', 1060000.00, '2026-10-26', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 26/10/2026', 'GRO-MINIMO-20261026', 'e99f0958-e81d-577c-b9b0-0ed158ba57e8')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('180c0eee-8ee4-50d4-8f5a-de8ab03ffbd2', 'e99f0958-e81d-577c-b9b0-0ed158ba57e8', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1060000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261027 (MINIMO - 20261027): Rp 1,400,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('a32d6821-4a0e-5a9c-b574-078e47cba966', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261027', 'CASH_RECEIPT', '2026-10-27', '2026-10-27 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261027', 1400000.00, 1400000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 27-10-2026 Narendra Sagara Aryaputra, Royce Wibowo, Bima Adhikara Pradana, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('07a13712-fd06-50fd-bb2c-1a816b58cf8c', 'a32d6821-4a0e-5a9c-b574-078e47cba966', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261027-11120', 'CASH_IN', 350000.00, '2026-10-27', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 27/10/2026', 'GRO-MINIMO-20261027', 'a32d6821-4a0e-5a9c-b574-078e47cba966')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('615d1656-a004-58d5-9c7c-bc622ce781e7', 'a32d6821-4a0e-5a9c-b574-078e47cba966', '2a3ad241-587d-4853-bdc2-615d72496b4e', 875000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261027-11111', 'CASH_IN', 875000.00, '2026-10-27', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 27/10/2026', 'GRO-MINIMO-20261027', 'a32d6821-4a0e-5a9c-b574-078e47cba966')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a1de1f3c-d154-59ac-a71d-077fee04c29a', 'a32d6821-4a0e-5a9c-b574-078e47cba966', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 175000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261027-11113', 'CASH_IN', 175000.00, '2026-10-27', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 27/10/2026', 'GRO-MINIMO-20261027', 'a32d6821-4a0e-5a9c-b574-078e47cba966')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('55c02d7d-9e7b-5bab-a96a-1ea7340cfea2', 'a32d6821-4a0e-5a9c-b574-078e47cba966', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1400000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261028 (MINIMO - 20261028): Rp 1,312,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('88be1409-9674-5fbb-937b-8844f030e986', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261028', 'CASH_RECEIPT', '2026-10-28', '2026-10-28 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261028', 1312500.00, 1312500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 28-10-2026 Dave Maverick Gunawan, Loui Asher Gunawan, Alysha Wibowo (3 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('25b12488-1d98-5d45-9f1d-b2e5be361def', '88be1409-9674-5fbb-937b-8844f030e986', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 525000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261028-11120', 'CASH_IN', 525000.00, '2026-10-28', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 28/10/2026', 'GRO-MINIMO-20261028', '88be1409-9674-5fbb-937b-8844f030e986')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7ca5d512-fced-5fb1-82b0-1c40a3946a0e', '88be1409-9674-5fbb-937b-8844f030e986', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 787500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261028-11113', 'CASH_IN', 787500.00, '2026-10-28', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 28/10/2026', 'GRO-MINIMO-20261028', '88be1409-9674-5fbb-937b-8844f030e986')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('0366e2b0-e6a5-5d50-8131-03cd91926ef4', '88be1409-9674-5fbb-937b-8844f030e986', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1312500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261029 (MINIMO - 20261029): Rp 2,100,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('9c227cbe-b93d-5c2e-9c58-f6a1a464f999', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261029', 'CASH_RECEIPT', '2026-10-29', '2026-10-29 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261029', 2100000.00, 2100000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 29-10-2026 Kyra Thalassa Shanum, Jaler CendekiaBashara Siddiq, Seraphina Meera Zafira, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7cc04d5b-fe6e-5614-96fe-56e1f7ef9b42', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 875000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261029-11113', 'CASH_IN', 875000.00, '2026-10-29', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 29/10/2026', 'GRO-MINIMO-20261029', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('acefd358-34ad-51b5-a506-8d0327e79678', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1050000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261029-11120', 'CASH_IN', 1050000.00, '2026-10-29', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 29/10/2026', 'GRO-MINIMO-20261029', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('4d91e000-1e1b-57fa-b9b6-e95fbec2451e', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999', '2a3ad241-587d-4853-bdc2-615d72496b4e', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261029-11111', 'CASH_IN', 175000.00, '2026-10-29', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 29/10/2026', 'GRO-MINIMO-20261029', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('92198df6-de8a-53c2-88e8-c264d68a2b62', '9c227cbe-b93d-5c2e-9c58-f6a1a464f999', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 2100000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261031 (MINIMO - 20261031): Rp 1,312,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('1335da66-d0b8-54bd-93bd-fdea68656cf5', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261031', 'CASH_RECEIPT', '2026-10-31', '2026-10-31 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261031', 1312500.00, 1312500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 31-10-2026 Elion Satya Theophilus, Eloisa Stella Valencia Santoso, Denaya Malaika Sidqi, dkk. (total 6 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('7da8b0b7-58da-5e38-8909-b414993323ad', '1335da66-d0b8-54bd-93bd-fdea68656cf5', '2a3ad241-587d-4853-bdc2-615d72496b4e', 525000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261031-11111', 'CASH_IN', 525000.00, '2026-10-31', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 31/10/2026', 'GRO-MINIMO-20261031', '1335da66-d0b8-54bd-93bd-fdea68656cf5')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('6f64a171-661a-56f2-92d7-f6d3f51eaf57', '1335da66-d0b8-54bd-93bd-fdea68656cf5', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 437500.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261031-11113', 'CASH_IN', 437500.00, '2026-10-31', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 31/10/2026', 'GRO-MINIMO-20261031', '1335da66-d0b8-54bd-93bd-fdea68656cf5')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('093b05a8-aa89-5e72-a8e4-c3fb9ad476cb', '1335da66-d0b8-54bd-93bd-fdea68656cf5', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 350000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261031-11120', 'CASH_IN', 350000.00, '2026-10-31', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 31/10/2026', 'GRO-MINIMO-20261031', '1335da66-d0b8-54bd-93bd-fdea68656cf5')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('c07e4cc4-b4b2-5095-b3cf-87779fbf1db3', '1335da66-d0b8-54bd-93bd-fdea68656cf5', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1312500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261101 (MINIMO - 20261101): Rp 1,225,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('c1f4d808-3341-556e-adbc-afbb2ed16d42', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261101', 'CASH_RECEIPT', '2026-11-01', '2026-11-01 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261101', 1225000.00, 1225000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 01-11-2026 Habib Satria, Kinara Anindhita Nafisa, Hamizan Sakha, dkk. (total 4 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('a1755ced-0a6e-5e6a-97a6-f1a5b4230f1e', 'c1f4d808-3341-556e-adbc-afbb2ed16d42', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 525000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261101-11113', 'CASH_IN', 525000.00, '2026-11-01', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 01/11/2026', 'GRO-MINIMO-20261101', 'c1f4d808-3341-556e-adbc-afbb2ed16d42')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('ae78048f-481d-5695-acbb-baa3631a7e64', 'c1f4d808-3341-556e-adbc-afbb2ed16d42', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 700000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261101-11120', 'CASH_IN', 700000.00, '2026-11-01', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 01/11/2026', 'GRO-MINIMO-20261101', 'c1f4d808-3341-556e-adbc-afbb2ed16d42')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('991789fb-8820-5529-be42-efaa3cd668c9', 'c1f4d808-3341-556e-adbc-afbb2ed16d42', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1225000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261102 (MINIMO - 20261102): Rp 612,500.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('7dbb05a0-6eb7-5d3b-a89a-3ac4629d9bf0', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261102', 'CASH_RECEIPT', '2026-11-02', '2026-11-02 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261102', 612500.00, 612500.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 02-11-2026 Arshila Shadira Andarfa, Elsie Claire Tjahjadi, Phoenix Trygve Hadirahardja, dkk. (total 5 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('b2c1d927-2091-550d-866d-8376e647aea3', '7dbb05a0-6eb7-5d3b-a89a-3ac4629d9bf0', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 612500.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261102-11120', 'CASH_IN', 612500.00, '2026-11-02', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 02/11/2026', 'GRO-MINIMO-20261102', '7dbb05a0-6eb7-5d3b-a89a-3ac4629d9bf0')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('658927c8-e6b2-5d76-8519-48e72adb66b2', '7dbb05a0-6eb7-5d3b-a89a-3ac4629d9bf0', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 612500.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261103 (MINIMO - 20261103): Rp 2,450,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('95e15e9d-fe02-5c1a-8430-d5b03cdd2405', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261103', 'CASH_RECEIPT', '2026-11-03', '2026-11-03 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261103', 2450000.00, 2450000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 03-11-2026 Mahakertamada Bhre Dace, Alayya Rafa Azura, Ayrton Sidney Maverick, dkk. (total 7 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('61c1b20d-75ad-5759-8fa4-05f4f1e2c863', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405', '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', 1750000.00, 0.00, 'Penerimaan Kas/Bank via Kliring Payment Gateway (QRIS/EDC/Tripay)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006', 'BKM-MINIMO-20261103-11120', 'CASH_IN', 1750000.00, '2026-11-03', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kliring Payment Gateway & EDC) - 03/11/2026', 'GRO-MINIMO-20261103', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('128859f7-49fe-5156-9c19-af175dad5acc', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405', '2a3ad241-587d-4853-bdc2-615d72496b4e', 175000.00, 0.00, 'Penerimaan Kas/Bank via Kas Kasir Kidinest');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000002', 'BKM-MINIMO-20261103-11111', 'CASH_IN', 175000.00, '2026-11-03', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Kas Kasir Kidinest) - 03/11/2026', 'GRO-MINIMO-20261103', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('5fd0f9c5-0a09-50eb-a691-64273d84a054', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 525000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261103-11113', 'CASH_IN', 525000.00, '2026-11-03', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 03/11/2026', 'GRO-MINIMO-20261103', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('de42a6d0-8e4d-573f-943f-6882d45cd21c', '95e15e9d-fe02-5c1a-8430-d5b03cdd2405', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 2450000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261105 (MINIMO - 20261105): Rp 350,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('2def39d5-1af9-52ed-be6b-eee7062dd980', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261105', 'CASH_RECEIPT', '2026-11-05', '2026-11-05 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261105', 350000.00, 350000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 05-11-2026 Andrea Renatta Desmon, Louisa Adrianna Desmon (2 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f39a86eb-90f0-5711-ac7a-3cc4804703d9', '2def39d5-1af9-52ed-be6b-eee7062dd980', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 350000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261105-11113', 'CASH_IN', 350000.00, '2026-11-05', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 05/11/2026', 'GRO-MINIMO-20261105', '2def39d5-1af9-52ed-be6b-eee7062dd980')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f8cf32ac-285b-5606-b2a5-49406859d4f5', '2def39d5-1af9-52ed-be6b-eee7062dd980', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 350000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261113 (MINIMO - 20261113): Rp 1,600,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('268bdde1-bc80-54b0-95eb-54baa0252c42', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261113', 'CASH_RECEIPT', '2026-11-13', '2026-11-13 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261113', 1600000.00, 1600000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 13-11-2026 Mikayla Rosalind Onggabean (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('eb84a42f-def1-5553-bed0-cc3cc11dca91', '268bdde1-bc80-54b0-95eb-54baa0252c42', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 1600000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261113-11113', 'CASH_IN', 1600000.00, '2026-11-13', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 13/11/2026', 'GRO-MINIMO-20261113', '268bdde1-bc80-54b0-95eb-54baa0252c42')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('e6e8865d-20c0-50d9-b70a-ab638e800c69', '268bdde1-bc80-54b0-95eb-54baa0252c42', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 1600000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            

    -- Journal CR-MNM-20261117 (MINIMO - 20261117): Rp 175,000.00
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('b07f1616-39d7-5503-99b9-50cc6c1057f6', 'a0000000-0000-0000-0000-000000000001', 'CR-MNM-20261117', 'CASH_RECEIPT', '2026-11-17', '2026-11-17 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', 'GRO-MINIMO-20261117', 175000.00, 175000.00, 'Penerimaan Kasir Harian [MINIMO] - Tanggal 17-11-2026 Samantha Abigail Saputra (1 pasien) - Ref: Daily Report GRO');
        

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('1842a9f8-b1e3-5a9c-9f1b-fe85c1c8c305', 'b07f1616-39d7-5503-99b9-50cc6c1057f6', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 175000.00, 0.00, 'Penerimaan Kas/Bank via Bank BCA Operasional (Giro & EDC)');
            

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKM-MINIMO-20261117-11113', 'CASH_IN', 175000.00, '2026-11-17', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 'Penerimaan Kasir Harian MINIMO (Bank BCA Operasional) - 17/11/2026', 'GRO-MINIMO-20261117', 'b07f1616-39d7-5503-99b9-50cc6c1057f6')
    ON CONFLICT (transaction_number) DO NOTHING;
                

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('f49d3428-fbd4-5ab4-9b30-c3e05c8fc68d', 'b07f1616-39d7-5503-99b9-50cc6c1057f6', '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', 0.00, 175000.00, 'Pendapatan Operasional Pendapatan Pemantauan Tumbuh Kembang Anak');
            
-- 5. Insert Balanced Operational Overhead & Expenses for Oktober 2026

    DELETE FROM journal_entries WHERE journal_number LIKE 'JV-2026-10-%';
    DELETE FROM cash_bank_transactions WHERE transaction_number LIKE 'BKK-202610-%';
    DELETE FROM tax_withholdings WHERE bupot_number LIKE 'BUPOT-%-202610-%';

    -- J10-1: Honor dr. Hendra Wijaya, Sp.A & PPh 21 Dokter Masa Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000001', 'a0000000-0000-0000-0000-000000000001', 'JV-2026-10-001', 'GENERAL', '2026-10-25', '2026-10-25 17:00:00+07', 'POSTED', 'MANUAL', 'BUPOT-21-202610-001', 46500000.00, 46500000.00, 'Pembayaran Honor & Pemotongan PPh 21 Dokter Spesialis Anak (dr. Hendra Wijaya, Sp.A) Masa Oktober 2026');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000001', 'd0000000-0000-0000-0000-100000000001', '4f29da55-b343-4441-a21a-9259ff0844b5', 46500000.00, 0.00, 'Beban Honor Jasa Dokter Spesialis Anak (Sp.A)'),
    ('e0000000-0000-0000-0000-100000000002', 'd0000000-0000-0000-0000-100000000001', 'b39b2593-66d0-41de-81ff-7b6c54faf290', 0.00, 1162500.00, 'Withholding PPh 21 Dokter 5% atas DPP 50%'),
    ('e0000000-0000-0000-0000-100000000003', 'd0000000-0000-0000-0000-100000000001', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 0.00, 45337500.00, 'Transfer Bank BCA Pembayaran Net Honor Dokter');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000004', 'BKK-202610-001', 'CASH_OUT', 45337500.00, '2026-10-25', '4f29da55-b343-4441-a21a-9259ff0844b5', 'Pembayaran Net Honor dr. Hendra Wijaya, Sp.A Masa Oktober', 'BUPOT-21-202610-001', 'd0000000-0000-0000-0000-100000000001');

    -- J10-2: Honor dr. Kevin Wijaya, Sp.DVE Masa Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000002', 'a0000000-0000-0000-0000-000000000002', 'JV-2026-10-002', 'GENERAL', '2026-10-25', '2026-10-25 17:00:00+07', 'POSTED', 'MANUAL', 'BUPOT-21-202610-002', 33000000.00, 33000000.00, 'Pembayaran Honor & Pemotongan PPh 21 Dokter Spesialis Kulit (dr. Kevin Wijaya, Sp.DVE) Masa Oktober 2026');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000004', 'd0000000-0000-0000-0000-100000000002', '47364f09-148f-4086-8eeb-1ace791ee795', 33000000.00, 0.00, 'Beban Honor Jasa Dokter Kulit & Terapis Estetika'),
    ('e0000000-0000-0000-0000-100000000005', 'd0000000-0000-0000-0000-100000000002', 'b39b2593-66d0-41de-81ff-7b6c54faf290', 0.00, 825000.00, 'Withholding PPh 21 Dokter 5% atas DPP 50%'),
    ('e0000000-0000-0000-0000-100000000006', 'd0000000-0000-0000-0000-100000000002', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 0.00, 32175000.00, 'Transfer Bank BCA Pembayaran Net Honor Dokter');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000004', 'BKK-202610-002', 'CASH_OUT', 32175000.00, '2026-10-25', '47364f09-148f-4086-8eeb-1ace791ee795', 'Pembayaran Net Honor dr. Kevin Wijaya, Sp.DVE Masa Oktober', 'BUPOT-21-202610-002', 'd0000000-0000-0000-0000-100000000002');

    -- J10-3: Limbah B3 Masa Oktober (PPh 23 2%)
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000003', '11111111-1111-1111-1111-111111111111', 'JV-2026-10-003', 'GENERAL', '2026-10-20', '2026-10-20 10:00:00+07', 'POSTED', 'MANUAL', 'BUPOT-23-202610-003', 18000000.00, 18000000.00, 'Biaya Pengolahan Limbah Medis B3 Paripurna PT Bio Medika Bersih (PPh 23 2%) Masa Oktober 2026');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000007', 'd0000000-0000-0000-0000-100000000003', 'b6e05ccd-baf4-4c7c-9387-e68ee68eaa94', 18000000.00, 0.00, 'Beban Pengolahan Limbah Medis B3'),
    ('e0000000-0000-0000-0000-100000000008', 'd0000000-0000-0000-0000-100000000003', '3ea36fdf-ddb0-4590-8a91-43aa7bacd46c', 0.00, 360000.00, 'Pemotongan Hutang PPh 23 Jasa 2%'),
    ('e0000000-0000-0000-0000-100000000009', 'd0000000-0000-0000-0000-100000000003', '58f965af-60a2-4926-a6ea-68d36cc46be7', 0.00, 17640000.00, 'Transfer Bank Mandiri Pelunasan Vendor B3');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000005', 'BKK-202610-003', 'CASH_OUT', 17640000.00, '2026-10-20', 'b6e05ccd-baf4-4c7c-9387-e68ee68eaa94', 'Pelunasan Vendor Limbah Medis B3 Oktober', 'INV-B3-202610', 'd0000000-0000-0000-0000-100000000003');

    -- J10-4: Sewa Gedung Klinik Jl. Cipto No. 100 Masa Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000004', '11111111-1111-1111-1111-111111111111', 'JV-2026-10-004', 'GENERAL', '2026-10-10', '2026-10-10 11:00:00+07', 'POSTED', 'MANUAL', 'BUPOT-42-202610-004', 60000000.00, 60000000.00, 'Sewa Gedung Klinik Operasional Jl. Cipto No. 100 Periode Oktober 2026 (PPh Final 4(2) 10%)');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000010', 'd0000000-0000-0000-0000-100000000004', 'b6e05ccd-baf4-4c7c-9387-e68ee68eaa94', 60000000.00, 0.00, 'Beban Sewa Gedung Klinik Operasional'),
    ('e0000000-0000-0000-0000-100000000011', 'd0000000-0000-0000-0000-100000000004', 'b39b2593-66d0-41de-81ff-7b6c54faf290', 0.00, 6000000.00, 'Hutang PPh Final Pasal 4(2) Sewa Gedung 10%'),
    ('e0000000-0000-0000-0000-100000000012', 'd0000000-0000-0000-0000-100000000004', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 0.00, 54000000.00, 'Transfer Giro BCA Pelunasan Sewa Gedung');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000004', 'BKK-202610-004', 'CASH_OUT', 54000000.00, '2026-10-10', 'b6e05ccd-baf4-4c7c-9387-e68ee68eaa94', 'Pelunasan Sewa Gedung Klinik Jl. Cipto Oktober', 'SEWA-202610', 'd0000000-0000-0000-0000-100000000004');

    -- J10-5: Angsuran PPh 25 Badan Masa Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000005', '11111111-1111-1111-1111-111111111111', 'JV-2026-10-005', 'GENERAL', '2026-10-12', '2026-10-12 14:00:00+07', 'POSTED', 'MANUAL', 'NTPN-MDR-20261012-25', 4500000.00, 4500000.00, 'Pembayaran Angsuran Pajak PPh 25 Badan Masa Oktober 2026 via Mandiri NTPN');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000013', 'd0000000-0000-0000-0000-100000000005', '5b57fe55-38a3-4dd7-8030-eb83a63a4c17', 4500000.00, 0.00, 'Uang Muka PPh 25 Badan Berjalan'),
    ('e0000000-0000-0000-0000-100000000014', 'd0000000-0000-0000-0000-100000000005', '58f965af-60a2-4926-a6ea-68d36cc46be7', 0.00, 4500000.00, 'Bank Mandiri Pelunasan NTPN Pajak');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000005', 'BKK-202610-005', 'CASH_OUT', 4500000.00, '2026-10-12', '5b57fe55-38a3-4dd7-8030-eb83a63a4c17', 'Angsuran Pajak PPh 25 Badan Kas Negara Oktober', 'NTPN-MDR-20261012-25', 'd0000000-0000-0000-0000-100000000005');

    -- J10-6: Restock Obat Kimia Farma Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000006', 'a0000000-0000-0000-0000-000000000001', 'JV-2026-10-006', 'GENERAL', '2026-10-07', '2026-10-07 10:00:00+07', 'POSTED', 'MANUAL', 'INV-SUPP-KF-7712', 28500000.00, 28500000.00, 'Pembelian Persediaan Obat Farmasi Anak Kidinest PT Kimia Farma Masa Oktober');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000015', 'd0000000-0000-0000-0000-100000000006', '3b13a78c-4c76-4442-b0df-778261e1b933', 28500000.00, 0.00, 'Persediaan Obat & Sirup Anak (Kidinest)'),
    ('e0000000-0000-0000-0000-100000000016', 'd0000000-0000-0000-0000-100000000006', '89da7d4b-4ebe-41bd-bd19-ce7ef20e3d80', 0.00, 28500000.00, 'Hutang Supplier Obat & Vaksin (Kidinest)');

    -- J10-7: Restock Vaksin Bio Farma Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000007', 'a0000000-0000-0000-0000-000000000001', 'JV-2026-10-007', 'GENERAL', '2026-10-09', '2026-10-09 10:00:00+07', 'POSTED', 'MANUAL', 'INV-SUPP-BF-9931', 35000000.00, 35000000.00, 'Pembelian Persediaan Vaksin Hexaxim & Rotavirus PT Bio Farma via BCA Masa Oktober');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000017', 'd0000000-0000-0000-0000-100000000007', '182b2210-5a12-4b7f-bb41-a95b436c6791', 35000000.00, 0.00, 'Persediaan Vaksin & Imunisasi (Kidinest)'),
    ('e0000000-0000-0000-0000-100000000018', 'd0000000-0000-0000-0000-100000000007', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 0.00, 35000000.00, 'Bank BCA Operasional');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('a0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000004', 'BKK-202610-007', 'CASH_OUT', 35000000.00, '2026-10-09', '182b2210-5a12-4b7f-bb41-a95b436c6791', 'Pelunasan Vaksin Bio Farma via BCA Oktober', 'INV-SUPP-BF-9931', 'd0000000-0000-0000-0000-100000000007');

    -- J10-8: Payroll Gaji Karyawan & Perawat Masa Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000008', '11111111-1111-1111-1111-111111111111', 'JV-2026-10-008', 'GENERAL', '2026-10-28', '2026-10-28 17:00:00+07', 'POSTED', 'MANUAL', 'PAYROLL-202610', 44000000.00, 44000000.00, 'Payroll Pembayaran Beban Gaji Karyawan, Perawat & Staf via Bank Mandiri Masa Oktober 2026');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000019', 'd0000000-0000-0000-0000-100000000008', '7bba663e-deb3-4c7b-8b95-875645152602', 44000000.00, 0.00, 'Beban Gaji Karyawan, Perawat & Staf'),
    ('e0000000-0000-0000-0000-100000000020', 'd0000000-0000-0000-0000-100000000008', '58f965af-60a2-4926-a6ea-68d36cc46be7', 0.00, 44000000.00, 'Bank Mandiri Payroll Karyawan');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000005', 'BKK-202610-008', 'CASH_OUT', 44000000.00, '2026-10-28', '7bba663e-deb3-4c7b-8b95-875645152602', 'Payroll Gaji Karyawan & Perawat Oktober', 'PAYROLL-202610', 'd0000000-0000-0000-0000-100000000008');

    -- J10-9: Admin Bank & EDC BCA Oktober
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('d0000000-0000-0000-0000-100000000009', '11111111-1111-1111-1111-111111111111', 'JV-2026-10-009', 'GENERAL', '2026-10-31', '2026-10-31 23:59:00+07', 'POSTED', 'MANUAL', 'ADM-BCA-202610', 750000.00, 750000.00, 'Beban Administrasi Bank & EDC BCA Periode Oktober 2026');

    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description) VALUES
    ('e0000000-0000-0000-0000-100000000021', 'd0000000-0000-0000-0000-100000000009', '8e9e51ca-fede-4ebf-a2d8-6ca302ff5199', 750000.00, 0.00, 'Beban Administrasi Bank & EDC'),
    ('e0000000-0000-0000-0000-100000000022', 'd0000000-0000-0000-0000-100000000009', '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', 0.00, 750000.00, 'Bank BCA Operasional');

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('11111111-1111-1111-1111-111111111111', 'c0000000-0000-0000-0000-000000000004', 'BKK-202610-009', 'CASH_OUT', 750000.00, '2026-10-31', '8e9e51ca-fede-4ebf-a2d8-6ca302ff5199', 'Biaya Administrasi & Sewa EDC Bank BCA Oktober', 'ADM-BCA-202610', 'd0000000-0000-0000-0000-100000000009');

    -- Seed Tax Withholdings for October (e-Bupot 21/23/4(2))
    INSERT INTO tax_withholdings (branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name, object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported, journal_id)
    VALUES
    ('a0000000-0000-0000-0000-000000000001', 'PPH21_DOCTOR', 'BUPOT-21-202610-001', '3374011203850001', '81.234.567.8-511.000', 'dr. Hendra Wijaya, Sp.A', '21-100-07', 46500000.00, 23250000.00, 0.05, 1162500.00, 2026, 10, TRUE, 'd0000000-0000-0000-0000-100000000001'),
    ('a0000000-0000-0000-0000-000000000002', 'PPH21_DOCTOR', 'BUPOT-21-202610-002', '3374021508880002', '82.345.678.9-511.000', 'dr. Kevin Wijaya, Sp.DVE', '21-100-07', 33000000.00, 16500000.00, 0.05, 825000.00, 2026, 10, TRUE, 'd0000000-0000-0000-0000-100000000002'),
    ('11111111-1111-1111-1111-111111111111', 'PPH23', 'BUPOT-23-202610-003', NULL, '01.345.678.9-503.000', 'PT Bio Medika Bersih (Limbah B3)', '28-100-02', 18000000.00, 18000000.00, 0.02, 360000.00, 2026, 10, TRUE, 'd0000000-0000-0000-0000-100000000003'),
    ('11111111-1111-1111-1111-111111111111', 'PPH4_2', 'BUPOT-42-202610-004', '3374032001700003', '71.456.789.0-511.000', 'H. Agus Suryanto (Pemilik Gedung)', '28-402-01', 60000000.00, 60000000.00, 0.10, 6000000.00, 2026, 10, TRUE, 'd0000000-0000-0000-0000-100000000004');
    
COMMIT;
