import openpyxl
import datetime
import uuid

EXCEL_PATH = '/Users/macbookpro/Documents/project-11/clinic/Daily Report GRO.xlsx'
OUTPUT_SQL = '/Users/macbookpro/Documents/project-11/clinic/backend/seeds/import_daily_report_gro.sql'

# Account IDs in PostgreSQL
ACCOUNTS = {
    '1-1110': '01542f52-866d-40d8-87c2-fd76688d053f', # Kas Kecil
    '1-1111': '2a3ad241-587d-4853-bdc2-615d72496b4e', # Kas Kasir Kidinest
    '1-1112': '69d4af2e-96d2-4856-bc88-da32eac1bc53', # Kas Kasir SSC
    '1-1113': '1a71d7f9-ae11-4b76-9a3b-87c4f34221f7', # Bank BCA
    '1-1114': '58f965af-60a2-4926-a6ea-68d36cc46be7', # Bank Mandiri
    '1-1120': '93a7af6f-818a-4b70-a5df-2ac3af7ba7ef', # Kliring Payment Gateway
    '1-1200': 'd9e6d198-d2cb-4233-abfe-373e79f99376', # Piutang Pasien
    '1-1310': '3b13a78c-4c76-4442-b0df-778261e1b933', # Persediaan Obat Anak
    '1-1320': '182b2210-5a12-4b7f-bb41-a95b436c6791', # Persediaan Vaksin
    '1-1330': '51d96060-3b9e-46a6-aad4-f427e8c8c868', # Persediaan Skincare
    '1-1340': 'c59cf3d3-c27a-4892-bdb3-32ffb971a935', # Persediaan BMHP
    '2-1110': '89da7d4b-4ebe-41bd-bd19-ce7ef20e3d80', # Hutang Supplier Obat
    '2-1115': 'e4d46584-be3d-4fc4-9340-470f9bdb0954', # Hutang Supplier Skincare
    '2-1120': 'f56c0c8e-3f9f-4afb-a034-a30b12c571a4', # Hutang Jasa Dokter Anak
    '2-1125': '7e5342bf-d77a-4c7e-bf46-6cd4a31aaf76', # Hutang Jasa Dokter Kulit
    '2-1130': '31e9948a-aa37-4e7c-b6f3-69b2de9b1061', # Hutang Gaji
    '2-1140': 'fb5c3156-3a4b-4a86-a896-ec354caef246', # Uang Muka Pasien / Deposit
    '4-1900': 'a0004190-0000-0000-0000-000000000000', # Potongan & Voucher Promosi
    '2-1500': 'b39b2593-66d0-41de-81ff-7b6c54faf290', # Hutang Pajak
    '2-1510': '3ea36fdf-ddb0-4590-8a91-43aa7bacd46c', # Kredit PPh 23
    '2-1520': '5b57fe55-38a3-4dd7-8030-eb83a63a4c17', # Uang Muka PPh 25
    '4-1100': '14004cd8-d83d-4954-bcc2-72d68e1d6999', # Pendapatan Registrasi Anak
    '4-1200': '44689ad5-00d4-487d-b1bb-98980bc4fdb8', # Pendapatan Dokter Sp.A
    '4-1300': '89dcca29-dfc6-4758-a18f-2ca40b0756cd', # Pendapatan Vaksin Anak
    '4-1400': '5d8b5f2f-8acb-4668-b54f-9d9dde5c75d5', # Pendapatan Tumbuh Kembang / Pijat
    '4-1500': '41a32972-5dba-4f6f-a175-4ce58b861e84', # Pendapatan Farmasi Anak
    '4-2100': '63585f66-7ed5-4e58-954a-689881e7e186', # Pendapatan Konsultasi Kulit
    '4-2200': 'eab0bf83-ae3b-4ea5-a320-3faa34488032', # Pendapatan Facial & Estetika
    '4-2300': '13f34cd9-557c-42ab-b0d8-09056fd2988f', # Pendapatan Laser & Cauter
    '4-2400': '6fc5cb16-ef68-437b-b0f5-c317dd40d40b', # Pendapatan Skincare SSC
    '5-1100': '4c5aad71-b560-4863-baed-f69fe56cf657', # HPP Obat
    '5-1200': '754de0e7-9cc3-4b0e-9955-2a28d7ce6999', # HPP Vaksin
    '5-1300': '380a8bf4-a426-41e7-b448-f905ed84def8', # HPP Skincare
    '5-1400': '080389a6-bbb9-438e-a898-f0ac1e1ced90', # HPP BMHP
    '6-1100': '4f29da55-b343-4441-a21a-9259ff0844b5', # Beban Honor Dokter Anak
    '6-1110': '47364f09-148f-4086-8eeb-1ace791ee795', # Beban Honor Dokter Kulit
    '6-1200': '7bba663e-deb3-4c7b-8b95-875645152602', # Beban Gaji Karyawan
    '6-1300': 'b6e05ccd-baf4-4c7c-9387-e68ee68eaa94', # Beban Listrik / Sewa
    '6-1400': '8e9e51ca-fede-4ebf-a2d8-6ca302ff5199'  # Beban Admin Bank
}

ACCOUNT_NAMES = {
    '1-1110': 'Kas Kecil (Petty Cash Operasional)',
    '1-1111': 'Kas Kasir Kidinest',
    '1-1112': 'Kas Kasir Semarang Skin Center (SSC)',
    '1-1113': 'Bank BCA Operasional (Giro & EDC)',
    '1-1114': 'Bank Mandiri Operasional & Payroll',
    '1-1120': 'Kliring Payment Gateway (QRIS/EDC/Tripay)',
    '1-1200': 'Piutang Pasien & Klaim Asuransi',
    '2-1140': 'Uang Muka Pasien / Deposit Tindakan',
    '4-1900': 'Potongan Penjualan & Voucher Promosi Pasien',
    '4-1100': 'Pendapatan Pendaftaran & Skrining Anak',
    '4-1200': 'Pendapatan Jasa Dokter Spesialis Anak (Sp.A)',
    '4-1300': 'Pendapatan Vaksinasi & Imunisasi Anak',
    '4-1400': 'Pendapatan Pemantauan Tumbuh Kembang Anak',
    '4-1500': 'Pendapatan Penjualan Obat & Racikan Anak',
    '4-2100': 'Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)',
    '4-2200': 'Pendapatan Treatment Facial & Estetika Medis',
    '4-2300': 'Pendapatan Tindakan Laser, Cauter & Peeling',
    '4-2400': 'Pendapatan Penjualan Produk Skincare & Sunscreen'
}

BRANCH_KID = 'a0000000-0000-0000-0000-000000000001'
BRANCH_SSC = 'a0000000-0000-0000-0000-000000000002'
BRANCH_MAIN = '11111111-1111-1111-1111-111111111111'

CBA_KAS_KID = 'c0000000-0000-0000-0000-000000000002'
CBA_KAS_SSC = 'c0000000-0000-0000-0000-000000000003'
CBA_BCA = 'c0000000-0000-0000-0000-000000000004'
CBA_MANDIRI = 'c0000000-0000-0000-0000-000000000005'
CBA_GATEWAY = 'c0000000-0000-0000-0000-000000000006'

def classify_kidinest(item, consult):
    s = str(item or '').lower().strip()
    c = str(consult or '').lower().strip()
    if c == 'v' or 'konsul' in s or 'dokter' in s or s == 'konsultasi' or s == 'periksa' or s == 'dr':
        return '4-1200'
    if any(vx in s for vx in ['vaksin', 'vaccin', 'infanrix', 'hexa', 'rotarix', 'rotavin', 'rotateq', 'prevenar', 'synflorix', 'hexaxim', 'mmr', 'varisela', 'varivax', 'typhim', 'havrix', 'avaxim', 'vaxigrip', 'fluarix', 'engerix', 'bcg', 'dpt', 'polio', 'campak', 'je', 'imojev', 'gardasil', 'bio farma', 'tetrabioko', 'pentabio']):
        return '4-1300'
    if any(tk in s for tk in ['pijat', 'massage', 'spa', 'tumbuh kembang', 'fisioterapi', 'terapi', 'nebulizer', 'nebul', 'suction', 'ear wash', 'tindik', 'irigasi']):
        return '4-1400'
    if any(lb in s for lb in ['lab', 'darah', 'urine', 'feces', 'feses', 'cbc', 'led', 'cito', 'swab', 'crp', 'ns1', 'skrining', 'admin', 'registrasi', 'buku']):
        return '4-1100'
    if s:
        return '4-1500'
    return '4-1200'

def classify_ssc(item, consult, dr):
    s = str(item or '').lower().strip()
    c = str(consult or '').lower().strip()
    d = str(dr or '').lower().strip()
    if c == 'v' or 'konsul' in s or 'consultation' in s or ('dr.' in d and not s):
        return '4-2100'
    if any(lx in s for lx in ['laser', 'nd:yag', 'cauter', 'peeling', 'tca', 'prp', 'botox', 'filler', 'skin booster', 'injection', 'injeksi']):
        return '4-2300'
    if any(fx in s for fx in ['facial', 'oxy', 'brightening', 'therapy', 'totok', 'infusion', 'massage', 'treatment', 'masker', 'anastesi', 'anasthesi', 'hydra']):
        return '4-2200'
    if s:
        return '4-2400'
    return '4-2100'

def get_debit_account(unit, p_name):
    p = p_name.upper()
    if 'TUNAI' in p:
        if unit == 'SSC':
            return '1-1112', CBA_KAS_SSC, 'Kas Kasir SSC'
        else:
            return '1-1111', CBA_KAS_KID, 'Kas Kasir Kidinest'
    elif any(bca in p for bca in ['BCA', 'TRANSFER']):
        return '1-1113', CBA_BCA, 'Bank BCA Operasional'
    elif 'CREDIT VOUCHER' in p or 'VOUCHER' in p:
        return '4-1900', None, 'Potongan Penjualan & Voucher Promosi Pasien'
    else:
        return '1-1120', CBA_GATEWAY, 'Kliring Payment Gateway & EDC'

def parse_excel():
    print(f'Loading {EXCEL_PATH}...')
    wb = openpyxl.load_workbook(EXCEL_PATH, data_only=True)

    units_cfg = [
        {
            'name': 'KIDINEST',
            'sheet': 'KIDINEST',
            'branch_id': BRANCH_KID,
            'prefix': 'CR-KID',
            'p_cols': [('BRI CC', 13), ('BRI DEBIT', 14), ('BCA CC', 15), ('BCA DEBIT', 16), ('Tunai', 17), ('Qris BCA', 18), ('Qris BRI', 19), ('Qris', 20), ('Tf BCA', 21), ('Credit Voucher', 22)],
            'v_col': 23,
            'item_col': 12,
            'consult_col': 8,
            'dr_col': None,
            'classifier': lambda item, con, dr: classify_kidinest(item, con)
        },
        {
            'name': 'SSC',
            'sheet': 'SSC',
            'branch_id': BRANCH_SSC,
            'prefix': 'CR-SSC',
            'p_cols': [('BRI CC', 15), ('BRI DEBIT', 16), ('BCA CC', 17), ('BCA DEBIT', 18), ('QR BCA', 19), ('Tunai', 20), ('CC CIMB', 21), ('Qris', 22), ('Credit Voucher', 23), ('Transfer', 24)],
            'v_col': 25,
            'item_col': 14,
            'consult_col': 8,
            'dr_col': 9,
            'classifier': lambda item, con, dr: classify_ssc(item, con, dr)
        },
        {
            'name': 'MINIMO',
            'sheet': 'MINIMO',
            'branch_id': BRANCH_KID,
            'prefix': 'CR-MNM',
            'p_cols': [('BRI CC', 13), ('BRI DEBIT', 14), ('BCA CC', 15), ('BCA DEBIT', 16), ('Tunai', 17), ('Qris BRI', 18), ('Qris BCA', 19), ('Qris', 20), ('Credit Voucher', 21)],
            'v_col': 22,
            'item_col': 12,
            'consult_col': 8,
            'dr_col': None,
            'classifier': lambda item, con, dr: '4-1400'
        }
    ]

    daily_journals = []

    for cfg in units_cfg:
        ws = wb[cfg['sheet']]
        current_date = None
        unit_daily = {}

        for r in range(4, ws.max_row + 1):
            d_val = ws.cell(r, 1).value
            if isinstance(d_val, (datetime.datetime, datetime.date)):
                year = 2026 if d_val.year in [2026, 2029] else 2026
                current_date = datetime.date(year, d_val.month, d_val.day)
            elif d_val and str(d_val).strip() not in ['', '#VALUE!']:
                s = str(d_val).strip()
                if '13/10' in s:
                    current_date = datetime.date(2026, 10, 13)
            elif d_val == '#VALUE!' and cfg['name'] == 'MINIMO' and r == 46:
                current_date = datetime.date(2026, 10, 8)

            name_val = ws.cell(r, 3).value
            item_val = ws.cell(r, cfg['item_col'] + 1).value
            mr_val = ws.cell(r, 2).value
            consult_val = ws.cell(r, cfg['consult_col'] + 1).value
            dr_val = ws.cell(r, cfg['dr_col'] + 1).value if cfg['dr_col'] else None

            if not name_val and not item_val and not mr_val and not consult_val:
                continue

            row_pay = {}
            row_total = 0
            for p_name, col_idx in cfg['p_cols']:
                c_val = ws.cell(r, col_idx + 1).value
                if isinstance(c_val, (int, float)) and c_val > 0:
                    row_pay[p_name] = c_val
                    row_total += c_val

            val_cell = ws.cell(r, cfg['v_col'] + 1).value
            if isinstance(val_cell, (int, float)) and val_cell > 0 and row_total == 0:
                row_total = val_cell
                row_pay['Tunai'] = val_cell

            if current_date and row_total > 0:
                d_str = current_date.strftime('%Y-%m-%d')
                if d_str not in unit_daily:
                    unit_daily[d_str] = {
                        'date': current_date,
                        'unit': cfg['name'],
                        'branch_id': cfg['branch_id'],
                        'prefix': cfg['prefix'],
                        'total': 0,
                        'tx_count': 0,
                        'patients': set(),
                        'debits': {},
                        'cba_map': {},
                        'credits': {}
                    }

                day_rec = unit_daily[d_str]
                day_rec['total'] += row_total
                day_rec['tx_count'] += 1
                if name_val and str(name_val).strip():
                    day_rec['patients'].add(str(name_val).strip())

                for p_name, amt in row_pay.items():
                    d_acc, cba_id, label = get_debit_account(cfg['name'], p_name)
                    day_rec['debits'][d_acc] = day_rec['debits'].get(d_acc, 0) + amt
                    if cba_id:
                        day_rec['cba_map'][d_acc] = (cba_id, label)

                rev_acc = cfg['classifier'](item_val, consult_val, dr_val)
                day_rec['credits'][rev_acc] = day_rec['credits'].get(rev_acc, 0) + row_total

        for d_str, day_rec in sorted(unit_daily.items()):
            tot_dr = sum(day_rec['debits'].values())
            tot_cr = sum(day_rec['credits'].values())
            diff = round(tot_dr - tot_cr, 2)
            if diff != 0:
                first_rev = list(day_rec['credits'].keys())[0]
                day_rec['credits'][first_rev] += diff
            daily_journals.append(day_rec)

    print(f'Total Daily Journals aggregated: {len(daily_journals)}')
    return daily_journals

def generate_sql(daily_journals):
    lines = []
    lines.append('-- ==============================================================================')
    lines.append('-- REAL CLINIC ACCOUNTING TRANSACTIONS IMPORT (DAILY REPORT GRO)')
    lines.append('-- Covers 2 Months: Oktober & November 2026 (+ Soft-Launch Akhir September)')
    lines.append('-- Generates Balanced Double-Entry Journal Entries, Journal Lines & Cash/Bank Tx')
    lines.append('-- ==============================================================================')
    lines.append('')
    lines.append('BEGIN;')
    lines.append('')

    # 1. Ensure Kliring Gateway exists in cash_bank_accounts
    lines.append('-- 1. Register Kliring Payment Gateway & EDC Account if not present')
    lines.append(f"""
    INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
    VALUES ('{CBA_GATEWAY}', '{BRANCH_MAIN}', '{ACCOUNTS['1-1120']}', 'GATEWAY-1120', 'Kliring Payment Gateway & EDC (QRIS/Tripay)', 'BANK', 'Settlement EDC & QRIS', 'EDC-QRIS-001', TRUE)
    ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, is_active = TRUE;
    """)

    # 2. Ensure Fiscal Periods exist for 2026-09, 2026-10, 2026-11
    lines.append('-- 2. Ensure Fiscal Periods exist for Sept, Oct, and Nov 2026')
    lines.append(f"""
    INSERT INTO fiscal_periods (id, branch_id, year, month, status)
    VALUES
        ('f0000000-0000-0000-0000-000000000009', '{BRANCH_MAIN}', 2026, 9, 'CLOSED'),
        ('f0000000-0000-0000-0000-000000000010', '{BRANCH_MAIN}', 2026, 10, 'CLOSED'),
        ('f0000000-0000-0000-0000-000000000011', '{BRANCH_MAIN}', 2026, 11, 'OPEN')
    ON CONFLICT (branch_id, year, month) DO UPDATE SET status = EXCLUDED.status;
    """)

    # 3. Clean up previously generated cashier receipts
    lines.append('-- 3. Clean up existing automated cashier closing journals to allow clean idempotent re-seed')
    lines.append("""
    DELETE FROM cash_bank_transactions WHERE transaction_number LIKE 'BKM-KID-%' OR transaction_number LIKE 'BKM-SSC-%' OR transaction_number LIKE 'BKM-MNM-%';
    DELETE FROM journal_entries WHERE journal_number LIKE 'CR-KID-%' OR journal_number LIKE 'CR-SSC-%' OR journal_number LIKE 'CR-MNM-%';
    -- Also remove old arbitrary dummy receipts on Nov 13-14
    DELETE FROM journal_entries WHERE journal_number IN ('JV-2026-11-007', 'JV-2026-11-008', 'JV-2026-11-017', 'JV-2026-11-018', 'JV-2026-11-019', 'JV-2026-11-020', 'JV-2026-11-021');
    """)

    # 4. Generate daily cashier journals
    lines.append('-- 4. Insert Daily Cashier Revenue Journals and Lines from Daily Report GRO')

    total_rev = 0
    tx_no = 0

    for j in daily_journals:
        d = j['date']
        d_compact = d.strftime('%Y%m%d')
        j_num = f"{j['prefix']}-{d_compact}"
        ref_id = f"GRO-{j['unit']}-{d_compact}"
        amount = sum(j['debits'].values())
        total_rev += amount
        tx_no += 1

        pts = [p.replace("'", "''") for p in list(j['patients'])[:3]]
        pts_str = ', '.join(pts)
        if len(j['patients']) > 3:
            pts_str += f", dkk. (total {len(j['patients'])} pasien)"
        elif pts:
            pts_str += f" ({len(j['patients'])} pasien)"
        else:
            pts_str = f"({j['tx_count']} transaksi)"

        desc = f"Penerimaan Kasir Harian [{j['unit']}] - Tanggal {d.strftime('%d-%m-%Y')} {pts_str} - Ref: Daily Report GRO"

        j_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, j_num))

        lines.append(f"""
    -- Journal {j_num} ({j['unit']} - {d_compact}): Rp {amount:,.2f}
    INSERT INTO journal_entries (id, branch_id, journal_number, journal_type, transaction_date, posting_date, status, reference_type, reference_id, total_debit, total_credit, description)
    VALUES ('{j_id}', '{j['branch_id']}', '{j_num}', 'CASH_RECEIPT', '{d.strftime('%Y-%m-%d')}', '{d.strftime('%Y-%m-%d')} 21:00:00+07', 'POSTED', 'CASHIER_CLOSING', '{ref_id}', {amount:.2f}, {amount:.2f}, '{desc}');
        """)

        # Journal Lines
        # Debits (Cash, Bank, Kliring Gateway, Voucher)
        for d_acc, d_amt in j['debits'].items():
            if d_amt <= 0: continue
            line_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{j_num}-DR-{d_acc}"))
            acc_name = ACCOUNT_NAMES.get(d_acc, d_acc).replace("'", "''")
            acc_id = ACCOUNTS[d_acc]
            line_desc = f"Penerimaan Kas/Bank via {acc_name}"
            lines.append(f"""
    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('{line_id}', '{j_id}', '{acc_id}', {d_amt:.2f}, 0.00, '{line_desc}');
            """)

            # Also create cash_bank_transaction if linked to cash_bank_account
            if d_acc in j['cba_map']:
                cba_id, cba_label = j['cba_map'][d_acc]
                cb_tx_num = f"BKM-{j['unit']}-{d_compact}-{d_acc.replace('-', '')}"
                cb_desc = f"Penerimaan Kasir Harian {j['unit']} ({cba_label}) - {d.strftime('%d/%m/%Y')}".replace("'", "''")
                lines.append(f"""
    INSERT INTO cash_bank_transactions (branch_id, cash_bank_account_id, transaction_number, transaction_type, amount, transaction_date, counterpart_account_id, description, reference_no, journal_id)
    VALUES ('{j['branch_id']}', '{cba_id}', '{cb_tx_num}', 'CASH_IN', {d_amt:.2f}, '{d.strftime('%Y-%m-%d')}', '{ACCOUNTS[list(j['credits'].keys())[0]]}', '{cb_desc}', '{ref_id}', '{j_id}')
    ON CONFLICT (transaction_number) DO NOTHING;
                """)

        # Credits (Revenue accounts)
        for c_acc, c_amt in j['credits'].items():
            if c_amt <= 0: continue
            line_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, f"{j_num}-CR-{c_acc}"))
            acc_name = ACCOUNT_NAMES.get(c_acc, c_acc).replace("'", "''")
            acc_id = ACCOUNTS[c_acc]
            line_desc = f"Pendapatan Operasional {acc_name}"
            lines.append(f"""
    INSERT INTO journal_lines (id, journal_entry_id, account_id, debit, credit, description)
    VALUES ('{line_id}', '{j_id}', '{acc_id}', 0.00, {c_amt:.2f}, '{line_desc}');
            """)

    # 5. Insert October Operational Overhead & Expenses
    lines.append('-- 5. Insert Balanced Operational Overhead & Expenses for Oktober 2026')
    lines.append("""
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
    """)

    lines.append('COMMIT;')
    lines.append('')

    with open(OUTPUT_SQL, 'w', encoding='utf-8') as f:
        f.write('\n'.join(lines))

    print(f'SQL written successfully to {OUTPUT_SQL} ({len(lines)} lines).')
    print(f'Total Daily Revenue Imported: Rp {total_rev:,.2f}')

if __name__ == '__main__':
    daily_journals = parse_excel()
    generate_sql(daily_journals)
