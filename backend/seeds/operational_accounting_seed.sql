-- Operational Accounting Data Seed for Clinic ERP (Kidinest & SSC)
-- Inserts Cash & Bank Master, Initial Journals, Transactions, and Tax Withholdings directly into PostgreSQL

DO $$
DECLARE
    v_branch_main UUID := '11111111-1111-1111-1111-111111111111';
    v_branch_kidinest UUID := 'a0000000-0000-0000-0000-000000000001';
    v_branch_ssc UUID := 'a0000000-0000-0000-0000-000000000002';
    
    -- Account IDs
    acc_petty UUID;
    acc_kas_kid UUID;
    acc_kas_ssc UUID;
    acc_bca UUID;
    acc_mandiri UUID;
    acc_piutang UUID;
    acc_persediaan_kid UUID;
    acc_persediaan_ssc UUID;
    acc_persediaan_bmhp UUID;
    acc_hutang_supp_kid UUID;
    acc_hutang_supp_ssc UUID;
    acc_hutang_dr_kid UUID;
    acc_hutang_dr_ssc UUID;
    acc_pajak_hutang UUID;
    acc_kredit_pph23 UUID;
    acc_um_pph25 UUID;
    acc_pendapatan_reg_kid UUID;
    acc_pendapatan_dr_kid UUID;
    acc_pendapatan_facial_ssc UUID;
    acc_pendapatan_laser_ssc UUID;
    acc_beban_honor_kid UUID;
    acc_beban_honor_ssc UUID;
    acc_beban_listrik UUID;

    -- CBA IDs
    cba_petty UUID := 'c0000000-0000-0000-0000-000000000001';
    cba_kas_kid UUID := 'c0000000-0000-0000-0000-000000000002';
    cba_kas_ssc UUID := 'c0000000-0000-0000-0000-000000000003';
    cba_bca UUID := 'c0000000-0000-0000-0000-000000000004';
    cba_mandiri UUID := 'c0000000-0000-0000-0000-000000000005';

    -- Journal IDs
    j1 UUID := 'd0000000-0000-0000-0000-000000000001';
    j2 UUID := 'd0000000-0000-0000-0000-000000000002';
    j3 UUID := 'd0000000-0000-0000-0000-000000000003';
    j4 UUID := 'd0000000-0000-0000-0000-000000000004';
    j5 UUID := 'd0000000-0000-0000-0000-000000000005';
    j6 UUID := 'd0000000-0000-0000-0000-000000000006';
    j7 UUID := 'd0000000-0000-0000-0000-000000000007';
    j8 UUID := 'd0000000-0000-0000-0000-000000000008';
    j9 UUID := 'd0000000-0000-0000-0000-000000000009';
BEGIN
    -- Fetch account IDs
    SELECT id INTO acc_petty FROM accounts WHERE code = '1-1110';
    SELECT id INTO acc_kas_kid FROM accounts WHERE code = '1-1111';
    SELECT id INTO acc_kas_ssc FROM accounts WHERE code = '1-1112';
    SELECT id INTO acc_bca FROM accounts WHERE code = '1-1113';
    SELECT id INTO acc_mandiri FROM accounts WHERE code = '1-1114';
    SELECT id INTO acc_piutang FROM accounts WHERE code = '1-1200';
    SELECT id INTO acc_persediaan_kid FROM accounts WHERE code = '1-1310';
    SELECT id INTO acc_persediaan_ssc FROM accounts WHERE code = '1-1330';
    SELECT id INTO acc_persediaan_bmhp FROM accounts WHERE code = '1-1340';
    SELECT id INTO acc_hutang_supp_kid FROM accounts WHERE code = '2-1110';
    SELECT id INTO acc_hutang_supp_ssc FROM accounts WHERE code = '2-1115';
    SELECT id INTO acc_hutang_dr_kid FROM accounts WHERE code = '2-1120';
    SELECT id INTO acc_hutang_dr_ssc FROM accounts WHERE code = '2-1125';
    SELECT id INTO acc_pajak_hutang FROM accounts WHERE code = '2-1500';
    SELECT id INTO acc_kredit_pph23 FROM accounts WHERE code = '2-1510';
    SELECT id INTO acc_um_pph25 FROM accounts WHERE code = '2-1520';
    SELECT id INTO acc_pendapatan_reg_kid FROM accounts WHERE code = '4-1100';
    SELECT id INTO acc_pendapatan_dr_kid FROM accounts WHERE code = '4-1200';
    SELECT id INTO acc_pendapatan_facial_ssc FROM accounts WHERE code = '4-2200';
    SELECT id INTO acc_pendapatan_laser_ssc FROM accounts WHERE code = '4-2300';
    SELECT id INTO acc_beban_honor_kid FROM accounts WHERE code = '6-1100';
    SELECT id INTO acc_beban_honor_ssc FROM accounts WHERE code = '6-1110';
    SELECT id INTO acc_beban_listrik FROM accounts WHERE code = '6-1300';

    -- 1. SEED CASH & BANK ACCOUNTS (Master Buku Kas & Bank)
    INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
    VALUES
        (cba_petty, v_branch_main, acc_petty, 'KAS-1110', 'Kas Kecil (Petty Cash Operasional)', 'CASH', NULL, NULL, TRUE),
        (cba_kas_kid, v_branch_kidinest, acc_kas_kid, 'KAS-1111', 'Kas Kasir Kidinest', 'CASH', NULL, NULL, TRUE),
        (cba_kas_ssc, v_branch_ssc, acc_kas_ssc, 'KAS-1112', 'Kas Kasir Semarang Skin Center (SSC)', 'CASH', NULL, NULL, TRUE),
        (cba_bca, v_branch_main, acc_bca, 'BANK-1113', 'Bank BCA Operasional (Giro & EDC)', 'BANK', 'PT Bank Central Asia Tbk', '035-888-9900', TRUE),
        (cba_mandiri, v_branch_main, acc_mandiri, 'BANK-1114', 'Bank Mandiri Operasional & Payroll', 'BANK', 'PT Bank Mandiri (Persero) Tbk', '136-00-1122334', TRUE)
    ON CONFLICT (code) DO UPDATE SET
        name = EXCLUDED.name,
        bank_name = EXCLUDED.bank_name,
        account_number = EXCLUDED.account_number;

    -- 2. SEED FISCAL PERIOD
    INSERT INTO fiscal_periods (id, branch_id, year, month, status)
    VALUES 
        ('f0000000-0000-0000-0000-000000000001', v_branch_main, 2026, 11, 'OPEN')
    ON CONFLICT (branch_id, year, month) DO NOTHING;

    -- 3. SEED BALANCED JOURNAL ENTRIES FOR NOVEMBER 2026
    -- J1: Honor dr. Hendra Wijaya, Sp.A & PPh 21 Dokter
    DELETE FROM journal_entries WHERE journal_number IN (
        'JV-2026-11-001', 'JV-2026-11-002', 'JV-2026-11-003', 'JV-2026-11-004',
        'JV-2026-11-005', 'JV-2026-11-006', 'JV-2026-11-007', 'JV-2026-11-008', 'JV-2026-11-009'
    );

    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES
    (j1, v_branch_kidinest, 'JV-2026-11-001', 'GENERAL', '2026-11-15', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'BUPOT-21-202611-001', 48500000.00, 48500000.00, 'Pembayaran Honor & Pemotongan PPh 21 Dokter Spesialis Anak (dr. Hendra Wijaya, Sp.A)'),
    (j2, v_branch_ssc, 'JV-2026-11-002', 'GENERAL', '2026-11-15', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'BUPOT-21-202611-002', 35000000.00, 35000000.00, 'Pembayaran Honor & Pemotongan PPh 21 Dokter Spesialis Kulit & Kelamin (dr. Kevin Wijaya, Sp.DVE)'),
    (j3, v_branch_main, 'JV-2026-11-003', 'GENERAL', '2026-11-16', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'BUPOT-23-202611-003', 18500000.00, 18500000.00, 'Biaya Pengolahan Limbah Medis B3 Paripurna PT Bio Medika Bersih (PPh 23 2%)'),
    (j4, v_branch_main, 'JV-2026-11-004', 'GENERAL', '2026-11-10', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'BUPOT-42-202611-004', 60000000.00, 60000000.00, 'Sewa Gedung Klinik Operasional Jl. Cipto No. 100 Periode Berjalan (PPh Final 4(2) 10%)'),
    (j5, v_branch_main, 'JV-2026-11-005', 'GENERAL', '2026-11-12', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'NTPN-MDR-20261112-25', 4500000.00, 4500000.00, 'Pembayaran Angsuran Pajak PPh 25 Badan Masa November 2026 via Mandiri NTPN'),
    (j6, v_branch_ssc, 'JV-2026-11-006', 'GENERAL', '2026-11-14', CURRENT_TIMESTAMP, 'POSTED', 'MANUAL', 'INV-SUPP-SSC-8891', 24420000.00, 24420000.00, 'Pembelian BMHP & Bahan Medis Estetika Treatment Laser (PPN Masukan 11%)'),
    (j7, v_branch_kidinest, 'JV-2026-11-007', 'CASH_RECEIPT', '2026-11-14', CURRENT_TIMESTAMP, 'POSTED', 'INVOICE', 'INV-KID-202611-01', 8750000.00, 8750000.00, 'Penerimaan Kasir Kidinest: Pendaftaran & Jasa Konsultasi Spesialis Anak'),
    (j8, v_branch_ssc, 'JV-2026-11-008', 'CASH_RECEIPT', '2026-11-14', CURRENT_TIMESTAMP, 'POSTED', 'INVOICE', 'INV-SSC-202611-01', 14200000.00, 14200000.00, 'Penerimaan Kasir SSC: Tindakan Laser Rejuvenation & Facial Estetika'),
    (j9, v_branch_main, 'JV-2026-11-009', 'BANK', '2026-11-15', CURRENT_TIMESTAMP, 'POSTED', 'TRANSFER', 'SETORAN-20261115', 22950000.00, 22950000.00, 'Setoran Penerimaan Kas Kasir Kidinest & SSC ke Rekening Giro Bank BCA');

    -- Insert Journal Lines for J1
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j1, acc_beban_honor_kid, 48500000.00, 0.00, 'Beban Honor Jasa Dokter Spesialis Anak (Sp.A)', 1),
    (j1, acc_pajak_hutang, 0.00, 1212500.00, 'Withholding PPh 21 Dokter 5% atas DPP 50%', 2),
    (j1, acc_bca, 0.00, 47287500.00, 'Transfer Bank BCA Pembayaran Net Honor Dokter', 3);

    -- Insert Journal Lines for J2
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j2, acc_beban_honor_ssc, 35000000.00, 0.00, 'Beban Honor Jasa Dokter Kulit & Terapis Estetika', 1),
    (j2, acc_pajak_hutang, 0.00, 875000.00, 'Withholding PPh 21 Dokter 5% atas DPP 50%', 2),
    (j2, acc_bca, 0.00, 34125000.00, 'Transfer Bank BCA Pembayaran Net Honor Dokter', 3);

    -- Insert Journal Lines for J3
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j3, acc_beban_listrik, 18500000.00, 0.00, 'Beban Pengolahan Limbah Medis B3', 1),
    (j3, acc_kredit_pph23, 0.00, 370000.00, 'Pemotongan Hutang PPh 23 Jasa 2%', 2),
    (j3, acc_mandiri, 0.00, 18130000.00, 'Transfer Bank Mandiri Pelunasan Vendor B3', 3);

    -- Insert Journal Lines for J4
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j4, acc_beban_listrik, 60000000.00, 0.00, 'Beban Sewa Gedung Klinik Operasional', 1),
    (j4, acc_pajak_hutang, 0.00, 6000000.00, 'Hutang PPh Final Pasal 4(2) Sewa Gedung 10%', 2),
    (j4, acc_bca, 0.00, 54000000.00, 'Transfer Giro BCA Pelunasan Sewa Gedung', 3);

    -- Insert Journal Lines for J5
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j5, acc_um_pph25, 4500000.00, 0.00, 'Uang Muka PPh 25 Badan Masa Berjalan', 1),
    (j5, acc_mandiri, 0.00, 4500000.00, 'Penyetoran Kas Negara via Bank Mandiri NTPN', 2);

    -- Insert Journal Lines for J6
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j6, acc_persediaan_bmhp, 22000000.00, 0.00, 'Persediaan BMHP Estetika & Facial (SSC)', 1),
    (j6, acc_pajak_hutang, 2420000.00, 0.00, 'PPN Masukan 11% Faktur Pajak Masukan Supplier', 2),
    (j6, acc_hutang_supp_ssc, 0.00, 24420000.00, 'Hutang Supplier Skincare & BMHP (SSC)', 3);

    -- Insert Journal Lines for J7
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j7, acc_kas_kid, 8750000.00, 0.00, 'Kas Penerimaan Kasir Unit Kidinest', 1),
    (j7, acc_pendapatan_reg_kid, 0.00, 1250000.00, 'Pendapatan Pendaftaran & Skrining Anak', 2),
    (j7, acc_pendapatan_dr_kid, 0.00, 7500000.00, 'Pendapatan Jasa Dokter Spesialis Anak (Sp.A)', 3);

    -- Insert Journal Lines for J8
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j8, acc_kas_ssc, 14200000.00, 0.00, 'Kas Penerimaan Kasir Unit SSC', 1),
    (j8, acc_pendapatan_facial_ssc, 0.00, 6200000.00, 'Pendapatan Treatment Facial & Estetika Medis', 2),
    (j8, acc_pendapatan_laser_ssc, 0.00, 8000000.00, 'Pendapatan Tindakan Laser, Cauter & Peeling', 3);

    -- Insert Journal Lines for J9
    INSERT INTO journal_lines (journal_entry_id, account_id, debit, credit, description, line_order) VALUES
    (j9, acc_bca, 22950000.00, 0.00, 'Bank BCA Operasional (Giro & EDC)', 1),
    (j9, acc_kas_kid, 0.00, 8750000.00, 'Kas Kasir Kidinest disetor ke Bank', 2),
    (j9, acc_kas_ssc, 0.00, 14200000.00, 'Kas Kasir SSC disetor ke Bank', 3);

    -- 4. SEED CASH & BANK TRANSACTIONS
    DELETE FROM cash_bank_transactions WHERE transaction_number IN (
        'BKK-202611-001', 'BKK-202611-002', 'BKK-202611-003', 'BKK-202611-004',
        'BKK-202611-005', 'BKM-202611-001', 'BKM-202611-002', 'TRF-202611-001'
    );

    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES
    (v_branch_main, cba_bca, 'BKK-202611-001', 'CASH_OUT', 47287500.00, '2026-11-15', acc_beban_honor_kid, 'Pembayaran Net Honor dr. Hendra Wijaya, Sp.A', 'BUPOT-21-202611-001', j1),
    (v_branch_main, cba_bca, 'BKK-202611-002', 'CASH_OUT', 34125000.00, '2026-11-15', acc_beban_honor_ssc, 'Pembayaran Net Honor dr. Kevin Wijaya, Sp.DVE', 'BUPOT-21-202611-002', j2),
    (v_branch_main, cba_mandiri, 'BKK-202611-003', 'CASH_OUT', 18130000.00, '2026-11-16', acc_beban_listrik, 'Pelunasan Pengolahan Limbah Medis B3', 'INV-B3-202611', j3),
    (v_branch_main, cba_bca, 'BKK-202611-004', 'CASH_OUT', 54000000.00, '2026-11-10', acc_beban_listrik, 'Pelunasan Sewa Gedung Klinik Jl. Cipto', 'SEWA-202611', j4),
    (v_branch_main, cba_mandiri, 'BKK-202611-005', 'CASH_OUT', 4500000.00, '2026-11-12', acc_um_pph25, 'Angsuran Pajak PPh 25 Badan Kas Negara', 'NTPN-MDR-20261112-25', j5),
    (v_branch_kidinest, cba_kas_kid, 'BKM-202611-001', 'CASH_IN', 8750000.00, '2026-11-14', acc_pendapatan_dr_kid, 'Penerimaan Kasir Pelayanan Anak Kidinest', 'INV-KID-202611-01', j7),
    (v_branch_ssc, cba_kas_ssc, 'BKM-202611-002', 'CASH_IN', 14200000.00, '2026-11-14', acc_pendapatan_laser_ssc, 'Penerimaan Kasir Treatment Laser SSC', 'INV-SSC-202611-01', j8),
    (v_branch_main, cba_bca, 'TRF-202611-001', 'CASH_IN', 22950000.00, '2026-11-15', acc_kas_kid, 'Penyetoran Kas Kasir Kidinest & SSC ke BCA Giro', 'SETORAN-20261115', j9);

    -- 5. SEED TAX WITHHOLDINGS (e-Bupot 21/23/4(2)/25)
    DELETE FROM tax_withholdings WHERE bupot_number IN (
        'BUPOT-21-202611-001', 'BUPOT-21-202611-002', 'BUPOT-23-202611-003', 'BUPOT-42-202611-004'
    );

    INSERT INTO tax_withholdings (branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name, object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported, journal_id)
    VALUES
    (v_branch_kidinest, 'PPH21_DOCTOR', 'BUPOT-21-202611-001', '3374011203850001', '81.234.567.8-511.000', 'dr. Hendra Wijaya, Sp.A', '21-100-07', 48500000.00, 24250000.00, 0.05, 1212500.00, 2026, 11, TRUE, j1),
    (v_branch_ssc, 'PPH21_DOCTOR', 'BUPOT-21-202611-002', '3374021508880002', '82.345.678.9-511.000', 'dr. Kevin Wijaya, Sp.DVE', '21-100-07', 35000000.00, 17500000.00, 0.05, 875000.00, 2026, 11, TRUE, j2),
    (v_branch_main, 'PPH23', 'BUPOT-23-202611-003', NULL, '01.345.678.9-503.000', 'PT Bio Medika Bersih (Limbah B3)', '28-100-02', 18500000.00, 18500000.00, 0.02, 370000.00, 2026, 11, TRUE, j3),
    (v_branch_main, 'PPH4_2', 'BUPOT-42-202611-004', '3374032001700003', '71.456.789.0-511.000', 'H. Agus Suryanto (Pemilik Gedung)', '28-402-01', 60000000.00, 60000000.00, 0.10, 6000000.00, 2026, 11, TRUE, j4);

END $$;
