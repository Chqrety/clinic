-- BCS (PT. Brayo Cipta Sejahtera) Complete Chart of Accounts & Seed Data
BEGIN;

INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b6c9c6b5-cc48-5d71-be85-c12755b1a66e', '1', 'ASET - BCS', 'ASSET', 'DEBIT', NULL, 1, TRUE, TRUE, 'Akun BCS 1 - ASET - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('a4496c5b-ecfb-58fd-9300-107c5e872797', '2', 'KEWAJIBAN - BCS', 'LIABILITY', 'CREDIT', NULL, 1, TRUE, TRUE, 'Akun BCS 2 - KEWAJIBAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6b3c2bf0-3006-55cd-896a-7fa628cc2adc', '3', 'EKUITAS - BCS', 'EQUITY', 'CREDIT', NULL, 1, TRUE, TRUE, 'Akun BCS 3 - EKUITAS - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('75c860f0-af71-560f-a1b1-68278c21580a', '4', 'PENDAPATAN - BCS', 'REVENUE', 'CREDIT', NULL, 1, TRUE, TRUE, 'Akun BCS 4 - PENDAPATAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('9067d9e5-2d7f-56f9-8ca3-10bdbbdf9333', '5', 'BEBAN - BCS', 'EXPENSE', 'DEBIT', NULL, 1, TRUE, TRUE, 'Akun BCS 5 - BEBAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('5034e238-7cac-50c1-beed-c72489f058bb', '6', 'PENDAPATAN NON USAHA - BCS', 'REVENUE', 'CREDIT', NULL, 1, TRUE, TRUE, 'Akun BCS 6 - PENDAPATAN NON USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('4872ccc5-255c-55f5-88d6-fb6375ccfe78', '7', 'BEBAN NON USAHA - BCS', 'EXPENSE', 'DEBIT', NULL, 1, TRUE, TRUE, 'Akun BCS 7 - BEBAN NON USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('c1d26d6a-28ec-552b-8490-e1602d2a6bb3', '11', 'KAS DAN BANK - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 11 - KAS DAN BANK - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('7ba1ce85-8756-5faf-80ae-9678066843f0', '12', 'PIUTANG - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 12 - PIUTANG - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f8eca542-2839-5611-b2e4-db51a67ab302', '13', 'Persediaan - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 13 - Persediaan - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ccf4a25e-dd98-521e-97e4-6cdfa2f7cdce', '14', 'INVESTASI - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 14 - INVESTASI - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b35ed619-b476-5504-94fc-d4bcdb146124', '15', 'BIAYA DIBAYAR DI MUKA - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 15 - BIAYA DIBAYAR DI MUKA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f1023803-2263-51ff-b92d-1dd57dd98e59', '16', 'AKTIVA TETAP - BCS', 'ASSET', 'DEBIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, TRUE, TRUE, 'Akun BCS 16 - AKTIVA TETAP - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('8506907a-a44f-51bd-aece-0fce8bfd2bc1', '1999', 'Temporary Opening - BCS', 'ASSET', 'CREDIT', 'b6c9c6b5-cc48-5d71-be85-c12755b1a66e', 2, FALSE, TRUE, 'Akun BCS 1999 - Temporary Opening - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('861638dc-3b2e-51e3-8e27-58127d12e8da', '21', 'UANG MUKA & TITIPAN - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, TRUE, TRUE, 'Akun BCS 21 - UANG MUKA & TITIPAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6d68b524-f668-549c-9232-f0778ba753d8', '22', 'HUTANG DAGANG - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, TRUE, TRUE, 'Akun BCS 22 - HUTANG DAGANG - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('88cd3013-4c91-59f6-bda6-bace4bfa2e55', '23', 'HUTANG BANK - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, FALSE, TRUE, 'Akun BCS 23 - HUTANG BANK - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('7504f921-4132-5989-8f7b-669ac8d6c180', '24', 'HUTANG PIHAK KETIGA - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, FALSE, TRUE, 'Akun BCS 24 - HUTANG PIHAK KETIGA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('770fc4a9-a6cc-5bd6-a72b-5829deb25a07', '25', 'HUTANG PAJAK - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, TRUE, TRUE, 'Akun BCS 25 - HUTANG PAJAK - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('4f693476-efd1-5e80-a0bb-646923f5420e', '26', 'HUTANG PEMBEBANAN HPP - BCS', 'LIABILITY', 'CREDIT', 'a4496c5b-ecfb-58fd-9300-107c5e872797', 2, FALSE, TRUE, 'Akun BCS 26 - HUTANG PEMBEBANAN HPP - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('9a0a7754-db1a-5c27-b9b9-6b68d4420f07', '31', 'MODAL - BCS', 'EQUITY', 'CREDIT', '6b3c2bf0-3006-55cd-896a-7fa628cc2adc', 2, TRUE, TRUE, 'Akun BCS 31 - MODAL - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('d07bd578-5157-5251-9214-453bb1533665', '32', 'LABA AKUMULASI - BCS', 'EQUITY', 'CREDIT', '6b3c2bf0-3006-55cd-896a-7fa628cc2adc', 2, TRUE, TRUE, 'Akun BCS 32 - LABA AKUMULASI - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ebb755aa-6878-5c0a-94b3-054e08d32759', '3999', 'Opening Balance Equity - BCS', 'EQUITY', 'CREDIT', '6b3c2bf0-3006-55cd-896a-7fa628cc2adc', 2, FALSE, TRUE, 'Akun BCS 3999 - Opening Balance Equity - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('386c856c-cb2f-5bf7-aa95-a4ac87503f0c', '41', 'PENDAPATAN USAHA - BCS', 'REVENUE', 'CREDIT', '75c860f0-af71-560f-a1b1-68278c21580a', 2, TRUE, TRUE, 'Akun BCS 41 - PENDAPATAN USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('baa44c98-5ecb-5705-8e7b-10db8683b967', '4999', 'Opening Balance Income - BCS', 'REVENUE', 'CREDIT', '75c860f0-af71-560f-a1b1-68278c21580a', 2, FALSE, TRUE, 'Akun BCS 4999 - Opening Balance Income - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('03223443-b03b-53c4-8793-51401c5b30f4', '51', 'HARGA POKOK PENJUALAN - BCS', 'COGS', 'DEBIT', '9067d9e5-2d7f-56f9-8ca3-10bdbbdf9333', 2, TRUE, TRUE, 'Akun BCS 51 - HARGA POKOK PENJUALAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('04d8f0e8-282d-5f33-b258-11c6d9867282', '52', 'BEBAN OPERASIONAL - BCS', 'EXPENSE', 'DEBIT', '9067d9e5-2d7f-56f9-8ca3-10bdbbdf9333', 2, TRUE, TRUE, 'Akun BCS 52 - BEBAN OPERASIONAL - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('885b1eaf-27a0-5b79-a3b4-6d50ef25b9e5', '5998', 'Opening Balance Expense - BCS', 'EXPENSE', 'DEBIT', '9067d9e5-2d7f-56f9-8ca3-10bdbbdf9333', 2, FALSE, TRUE, 'Akun BCS 5998 - Opening Balance Expense - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('fe8be95c-3317-5184-b16f-c4063718d825', '5999', 'Penyesuaian Stok - BCS', 'EXPENSE', 'DEBIT', '9067d9e5-2d7f-56f9-8ca3-10bdbbdf9333', 2, FALSE, TRUE, 'Akun BCS 5999 - Penyesuaian Stok - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', '61', 'PENDAPATAN LAIN-LAIN NON USAHA - BCS', 'REVENUE', 'CREDIT', '5034e238-7cac-50c1-beed-c72489f058bb', 2, TRUE, TRUE, 'Akun BCS 61 - PENDAPATAN LAIN-LAIN NON USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('312816a3-6854-5c19-ae6c-fe11ab01b00f', '71', 'BIAYA LAIN-LAIN NON USAHA - BCS', 'EXPENSE', 'DEBIT', '4872ccc5-255c-55f5-88d6-fb6375ccfe78', 2, TRUE, TRUE, 'Akun BCS 71 - BIAYA LAIN-LAIN NON USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f131a934-6728-5c7d-abf5-3a7360c89a55', '1101', 'Kas Kasir & Operasional - BCS', 'ASSET', 'DEBIT', 'c1d26d6a-28ec-552b-8490-e1602d2a6bb3', 3, TRUE, TRUE, 'Akun BCS 1101 - Kas Kasir & Operasional - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('31f47ef8-7b6b-507b-beb8-ddfe04ec7169', '1102', 'Bank Operasional & Kliring - BCS', 'ASSET', 'DEBIT', 'c1d26d6a-28ec-552b-8490-e1602d2a6bb3', 3, TRUE, TRUE, 'Akun BCS 1102 - Bank Operasional & Kliring - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('cd77184b-224d-5b41-8981-8bcc078aa19c', '1201', 'Piutang Pasien Perawatan Rawat Jalan & Treatment', 'ASSET', 'DEBIT', '7ba1ce85-8756-5faf-80ae-9678066843f0', 3, FALSE, TRUE, 'Akun BCS 1201 - Piutang Pasien Perawatan Rawat Jalan & Treatment')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('876554c3-32f7-5f86-8119-723b3303594c', '1202', 'Piutang Klaim Asuransi & Korporat', 'ASSET', 'DEBIT', '7ba1ce85-8756-5faf-80ae-9678066843f0', 3, TRUE, TRUE, 'Akun BCS 1202 - Piutang Klaim Asuransi & Korporat')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('19ee0051-d7e7-535a-959c-a0b0b2ba48cd', '1203', 'Cadangan Kerugian Piutang Tak Tertagih', 'ASSET', 'DEBIT', '7ba1ce85-8756-5faf-80ae-9678066843f0', 3, FALSE, TRUE, 'Akun BCS 1203 - Cadangan Kerugian Piutang Tak Tertagih')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6a77bf6f-4baf-53be-b47d-b4438eb35c51', '1301', 'Persediaan Obat Resep & Paten Farmasi Kidinest', 'ASSET', 'DEBIT', 'f8eca542-2839-5611-b2e4-db51a67ab302', 3, FALSE, TRUE, 'Akun BCS 1301 - Persediaan Obat Resep & Paten Farmasi Kidinest')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('435daefa-ec7d-5740-bae9-b2b6dcd2d767', '1302', 'Persediaan Produk Skincare & Kosmetika SSC', 'ASSET', 'DEBIT', 'f8eca542-2839-5611-b2e4-db51a67ab302', 3, FALSE, TRUE, 'Akun BCS 1302 - Persediaan Produk Skincare & Kosmetika SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('d9543bee-8590-5793-a50b-6cba71957d03', '1303', 'Persediaan Bahan Medis Habis Pakai (BMHP)', 'ASSET', 'DEBIT', 'f8eca542-2839-5611-b2e4-db51a67ab302', 3, FALSE, TRUE, 'Akun BCS 1303 - Persediaan Bahan Medis Habis Pakai (BMHP)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('63a3bbee-d65e-5fc7-9052-5735d5f3fd7e', '1304', 'Persediaan Vaksin & Serum Imunisasi Anak', 'ASSET', 'DEBIT', 'f8eca542-2839-5611-b2e4-db51a67ab302', 3, FALSE, TRUE, 'Akun BCS 1304 - Persediaan Vaksin & Serum Imunisasi Anak')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6a0f58a8-a062-5f41-9659-67f5eed7e2ba', '1401', 'Investasi Tanah & Bangunan Cabang Baru', 'ASSET', 'DEBIT', 'ccf4a25e-dd98-521e-97e4-6cdfa2f7cdce', 3, FALSE, TRUE, 'Akun BCS 1401 - Investasi Tanah & Bangunan Cabang Baru')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('8362486f-5ac5-5520-a0cd-7827d2bda155', '1402', 'Investasi Pengembangan Riset Formulasi Skincare', 'ASSET', 'DEBIT', 'ccf4a25e-dd98-521e-97e4-6cdfa2f7cdce', 3, FALSE, TRUE, 'Akun BCS 1402 - Investasi Pengembangan Riset Formulasi Skincare')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('90bdd086-b704-5fcf-a0fb-8cf87c9575e8', '1501', 'Sewa Ruko & Fasilitas Dibayar di Muka', 'ASSET', 'DEBIT', 'b35ed619-b476-5504-94fc-d4bcdb146124', 3, FALSE, TRUE, 'Akun BCS 1501 - Sewa Ruko & Fasilitas Dibayar di Muka')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('3b410d71-b34b-5313-a0d6-2fded1a0fb1e', '1504', 'Uang Muka PPh 25 & Pajak Masukan (Kredit Pajak)', 'ASSET', 'DEBIT', 'b35ed619-b476-5504-94fc-d4bcdb146124', 3, FALSE, TRUE, 'Akun BCS 1504 - Uang Muka PPh 25 & Pajak Masukan (Kredit Pajak)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('855182f8-8b7f-5845-9d0d-33a12c1b08bc', '1601', 'Bangunan & Renovasi Interior Klinik', 'ASSET', 'DEBIT', 'f1023803-2263-51ff-b92d-1dd57dd98e59', 3, FALSE, TRUE, 'Akun BCS 1601 - Bangunan & Renovasi Interior Klinik')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('8948a658-fac8-5111-8159-e872e224450f', '1602', 'Mesin Laser Nd:YAG & Alkes Estetika SSC', 'ASSET', 'DEBIT', 'f1023803-2263-51ff-b92d-1dd57dd98e59', 3, FALSE, TRUE, 'Akun BCS 1602 - Mesin Laser Nd:YAG & Alkes Estetika SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('8b050433-4630-5c2b-b979-450c1f5bbccb', '1603', 'Peralatan USG 4D & Poli Anak Kidinest', 'ASSET', 'DEBIT', 'f1023803-2263-51ff-b92d-1dd57dd98e59', 3, FALSE, TRUE, 'Akun BCS 1603 - Peralatan USG 4D & Poli Anak Kidinest')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('4e485584-050e-57b8-ad9f-1355f60f663c', '1604', 'Inventaris Komputer, POS Kasir & Jaringan', 'ASSET', 'DEBIT', 'f1023803-2263-51ff-b92d-1dd57dd98e59', 3, FALSE, TRUE, 'Akun BCS 1604 - Inventaris Komputer, POS Kasir & Jaringan')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('1b8cf641-6e28-5ea0-aeeb-8ed50db9bb8f', '1606', 'Akumulasi Penyusutan Aktiva Tetap', 'ASSET', 'CREDIT', 'f1023803-2263-51ff-b92d-1dd57dd98e59', 3, FALSE, TRUE, 'Akun BCS 1606 - Akumulasi Penyusutan Aktiva Tetap')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('4f31f30a-cb1c-5b5c-904a-c482cc330fb5', '2101', 'UANG TANDA JADI', 'LIABILITY', 'CREDIT', '861638dc-3b2e-51e3-8e27-58127d12e8da', 3, TRUE, TRUE, 'Akun BCS 2101 - UANG TANDA JADI')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('0819ad8c-3476-5e21-a502-1a5ae5c3fbb3', '2102', 'UANG MUKA PENJUALAN', 'LIABILITY', 'CREDIT', '861638dc-3b2e-51e3-8e27-58127d12e8da', 3, TRUE, TRUE, 'Akun BCS 2102 - UANG MUKA PENJUALAN')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('5bf83956-e38f-5ccd-9253-efbb98edf14d', '2103', 'TITIPAN BIAYA', 'LIABILITY', 'CREDIT', '861638dc-3b2e-51e3-8e27-58127d12e8da', 3, FALSE, TRUE, 'Akun BCS 2103 - TITIPAN BIAYA')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('df80c80e-93f1-5212-a2d6-1c44d2ea15ab', '2201', 'Hutang PBF Distributor Farmasi & Obat', 'LIABILITY', 'CREDIT', '6d68b524-f668-549c-9232-f0778ba753d8', 3, FALSE, TRUE, 'Akun BCS 2201 - Hutang PBF Distributor Farmasi & Obat')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('5fa35d2c-6e1a-51a7-84a8-35513f271a3d', '2202', 'Hutang Supplier Kosmetika & Skincare SSC', 'LIABILITY', 'CREDIT', '6d68b524-f668-549c-9232-f0778ba753d8', 3, FALSE, TRUE, 'Akun BCS 2202 - Hutang Supplier Kosmetika & Skincare SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('333eb5dd-f1a6-5620-ab25-c6f41c6944da', '2501', 'Hutang PPh 21 (Dokter & Karyawan)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, TRUE, TRUE, 'Akun BCS 2501 - Hutang PPh 21 (Dokter & Karyawan)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('8432219e-cbd6-5d56-9cb5-92a18d865ac6', '2502', 'Hutang PPh 23 (Jasa Medis, B3 & Maintenance Alkes)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, FALSE, TRUE, 'Akun BCS 2502 - Hutang PPh 23 (Jasa Medis, B3 & Maintenance Alkes)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('2d3d3080-d0c2-5eff-a0ef-d228c4654574', '2503', 'Hutang PPh Final Pasal 4(2) (Sewa Gedung Ruko)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, FALSE, TRUE, 'Akun BCS 2503 - Hutang PPh Final Pasal 4(2) (Sewa Gedung Ruko)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('d1f658f5-3255-5256-b544-6b7bcf23b005', '2504', 'Hutang PPN Keluaran (Penjualan Kosmetika & Skincare)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, FALSE, TRUE, 'Akun BCS 2504 - Hutang PPN Keluaran (Penjualan Kosmetika & Skincare)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ffa7306a-c5f2-5b0e-9849-5f23babcbdef', '2505', 'Hutang PPh 25 (Angsuran Pajak Penghasilan Badan)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, FALSE, TRUE, 'Akun BCS 2505 - Hutang PPh 25 (Angsuran Pajak Penghasilan Badan)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b9b2d121-98b3-524a-9a0f-98b7a9c6f1eb', '2506', 'Hutang PPh 29 (SPT Tahunan PPh Badan 1771)', 'LIABILITY', 'CREDIT', '770fc4a9-a6cc-5bd6-a72b-5829deb25a07', 3, FALSE, TRUE, 'Akun BCS 2506 - Hutang PPh 29 (SPT Tahunan PPh Badan 1771)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b1bf6f92-79ad-57b8-bdb1-45ba4cf1e2d5', '3101', 'Modal Disetor Pemilik Klinik Kidinest & SSC', 'EQUITY', 'CREDIT', '9a0a7754-db1a-5c27-b9b9-6b68d4420f07', 3, FALSE, TRUE, 'Akun BCS 3101 - Modal Disetor Pemilik Klinik Kidinest & SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f815afe5-dc0f-5f48-b97d-4f723566b3af', '3201', 'Laba Ditahan (Retained Earnings)', 'EQUITY', 'CREDIT', 'd07bd578-5157-5251-9214-453bb1533665', 3, FALSE, TRUE, 'Akun BCS 3201 - Laba Ditahan (Retained Earnings)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('708643d0-00da-5cb3-ace6-4f70ba6f6139', '3202', 'Laba Periode Berjalan', 'EQUITY', 'CREDIT', 'd07bd578-5157-5251-9214-453bb1533665', 3, FALSE, TRUE, 'Akun BCS 3202 - Laba Periode Berjalan')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('20c26efe-0653-5d7a-9e7b-97612352f3b7', '4101', 'PENJUALAN RUMAH / TINDAKAN UTAMA - BCS', 'REVENUE', 'CREDIT', '386c856c-cb2f-5bf7-aa95-a4ac87503f0c', 3, FALSE, TRUE, 'Akun BCS 4101 - PENJUALAN RUMAH / TINDAKAN UTAMA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ab0b8384-bc6a-5eeb-8c17-32705d4f2016', '4102', 'PENJUALAN RUKO & KIOS - BCS', 'REVENUE', 'CREDIT', '386c856c-cb2f-5bf7-aa95-a4ac87503f0c', 3, FALSE, TRUE, 'Akun BCS 4102 - PENJUALAN RUKO & KIOS - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('a509ac29-9abd-5157-b265-e1cfff956a2a', '4103', 'PENJUALAN KAVLING / KELEBIHAN TANAH - BCS', 'REVENUE', 'CREDIT', '386c856c-cb2f-5bf7-aa95-a4ac87503f0c', 3, FALSE, TRUE, 'Akun BCS 4103 - PENJUALAN KAVLING / KELEBIHAN TANAH - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('c6ecfb4c-b4c0-550f-a1c6-6ed322a11199', '5102', 'HPP PEROLEHAN TANAH / BAHAN MEDIS UTAMA', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5102 - HPP PEROLEHAN TANAH / BAHAN MEDIS UTAMA')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('5d187282-4189-57aa-9bc8-12b56063f11d', '5104', 'HPP LEGALITAS & PERIZINAN', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5104 - HPP LEGALITAS & PERIZINAN')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('735a6854-2363-525c-a632-6de9e149cea7', '5105', 'HPP COR JALAN UTAMA / PRASARANA', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5105 - HPP COR JALAN UTAMA / PRASARANA')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('137224e6-6f84-5a0b-8361-0f9d4b40452f', '5107', 'HPP PONDASI TURAP', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5107 - HPP PONDASI TURAP')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('4570c781-9a89-586f-bf58-ff0349591fda', '5109', 'HPP SALURAN', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5109 - HPP SALURAN')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('17386582-b829-5d4b-982d-b68a90159fa9', '5110', 'HPP JARINGAN LISTRIK', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5110 - HPP JARINGAN LISTRIK')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('9779baa6-cb1b-5fe0-823d-0f7faa6b08ab', '5111', 'HPP JARINGAN AIR', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5111 - HPP JARINGAN AIR')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6a702cc2-80ec-5d2f-b6c7-65e86d02eac3', '5114', 'HPP INFRASTRUKTUR/KANSTIN', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5114 - HPP INFRASTRUKTUR/KANSTIN')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('72b65371-d2d2-532b-8de3-e9720599034f', '5119', 'HPP RUMAH / TINDAKAN UNIT', 'COGS', 'DEBIT', '03223443-b03b-53c4-8793-51401c5b30f4', 3, FALSE, TRUE, 'Akun BCS 5119 - HPP RUMAH / TINDAKAN UNIT')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('450f8e28-2105-5567-9091-8f319d18385f', '5201', 'BIAYA MARKETING & PROMOSI', 'EXPENSE', 'DEBIT', '04d8f0e8-282d-5f33-b258-11c6d9867282', 3, FALSE, TRUE, 'Akun BCS 5201 - BIAYA MARKETING & PROMOSI')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('5d49dfeb-c8e9-5cff-a211-c53e072dc7bc', '5202', 'BEBAN UMUM DAN ADMINISTRASI (Gaji, Listrik, Sewa, Limbah)', 'EXPENSE', 'DEBIT', '04d8f0e8-282d-5f33-b258-11c6d9867282', 3, FALSE, TRUE, 'Akun BCS 5202 - BEBAN UMUM DAN ADMINISTRASI (Gaji, Listrik, Sewa, Limbah)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f9421b35-56a5-5171-9a83-6630ca20a118', '6101', 'PENDAPATAN BUNGA BANK - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6101 - PENDAPATAN BUNGA BANK - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b8794dfe-b3f8-50ad-8b3d-7de0ac0c1c47', '6102', 'PENDAPATAN PEMBATALAN UNIT - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6102 - PENDAPATAN PEMBATALAN UNIT - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('149b0016-ae4d-55fa-9faa-5ef6b16a931d', '6103', 'PENDAPATAN PINDAH KAVLING - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6103 - PENDAPATAN PINDAH KAVLING - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('354ce71f-a879-5a02-bff0-d46042755c7d', '6104', 'LABA PENJUALAN AKTIVA - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6104 - LABA PENJUALAN AKTIVA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('6c61b429-e285-52ce-a466-8db9ad2fdc39', '6105', 'SELISIH PEMBULATAN - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6105 - SELISIH PEMBULATAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('15578450-a020-5641-b634-4e863e2f5808', '6106', 'PENDAPATAN LAIN-LAIN NON USAHA - BCS', 'REVENUE', 'CREDIT', '9d99d2e7-51e4-50e9-88bb-c556c65d2bfd', 3, FALSE, TRUE, 'Akun BCS 6106 - PENDAPATAN LAIN-LAIN NON USAHA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b1f3f390-9d93-530c-9f14-16d77e8a50e5', '7101', 'BEBAN PROVISI KYG KPL - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7101 - BEBAN PROVISI KYG KPL - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('67f8a8c6-4eb2-5631-a3ce-7df6dda5d8f8', '7102', 'BUNGA KYG KPL - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7102 - BUNGA KYG KPL - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('485744e8-4bb9-5d31-8b39-b14663b109b2', '7103', 'BEBAN BIAYA ADMINISTRASI BANK - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7103 - BEBAN BIAYA ADMINISTRASI BANK - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('2909ab43-e432-57d9-95d7-4d06ac705119', '7104', 'PAJAK JASA GIRO - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7104 - PAJAK JASA GIRO - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('082781f2-82a2-5269-b994-e214ebf4f023', '7105', 'KERUGIAN PENJUALAN AKTIVA - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7105 - KERUGIAN PENJUALAN AKTIVA - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('e871745d-5e56-5532-8a73-f6831e039041', '7106', 'SELISIH PEMBULATAN & PENYESUAIAN - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7106 - SELISIH PEMBULATAN & PENYESUAIAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('eeacfd28-05d8-5e40-afdd-7c26ac8bd6ba', '7107', 'BEBAN BUNGA PINJAMAN LAIN - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7107 - BEBAN BUNGA PINJAMAN LAIN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('2296bcd5-1bbc-5dc8-bc9b-cf2a0fbd5542', '7108', 'AKUN PENGHAPUSAN - BCS', 'EXPENSE', 'DEBIT', '312816a3-6854-5c19-ae6c-fe11ab01b00f', 3, FALSE, TRUE, 'Akun BCS 7108 - AKUN PENGHAPUSAN - BCS')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('77dd7a90-8ae3-59b8-a65d-d76f38628738', '1101.001', 'Kas Kasir Kidinest', 'ASSET', 'DEBIT', 'f131a934-6728-5c7d-abf5-3a7360c89a55', 4, FALSE, TRUE, 'Akun BCS 1101.001 - Kas Kasir Kidinest')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('444d0bee-6bab-5a11-b703-dc216563e8d5', '1101.002', 'Kas Kasir Semarang Skin Center (SSC)', 'ASSET', 'DEBIT', 'f131a934-6728-5c7d-abf5-3a7360c89a55', 4, FALSE, TRUE, 'Akun BCS 1101.002 - Kas Kasir Semarang Skin Center (SSC)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ccb4cf53-4b16-5cd6-b5c6-600d2724dc50', '1101.003', 'Kas Kecil (Petty Cash Operasional)', 'ASSET', 'DEBIT', 'f131a934-6728-5c7d-abf5-3a7360c89a55', 4, FALSE, TRUE, 'Akun BCS 1101.003 - Kas Kecil (Petty Cash Operasional)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('7ddc4922-b987-5268-a704-d53880f59d75', '1102.001', 'Bank BCA Operasional (Giro & EDC)', 'ASSET', 'DEBIT', '31f47ef8-7b6b-507b-beb8-ddfe04ec7169', 4, FALSE, TRUE, 'Akun BCS 1102.001 - Bank BCA Operasional (Giro & EDC)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('0033a4e5-ae85-5a17-a004-4ba0f8b26f15', '1102.002', 'Bank Mandiri Operasional & Payroll', 'ASSET', 'DEBIT', '31f47ef8-7b6b-507b-beb8-ddfe04ec7169', 4, FALSE, TRUE, 'Akun BCS 1102.002 - Bank Mandiri Operasional & Payroll')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('cf3a16d6-659e-5911-b80d-bf3bf527daba', '1102.003', 'Settlement QRIS & Payment Gateway Transit', 'ASSET', 'DEBIT', '31f47ef8-7b6b-507b-beb8-ddfe04ec7169', 4, FALSE, TRUE, 'Akun BCS 1102.003 - Settlement QRIS & Payment Gateway Transit')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('ce371684-0182-5a22-bb10-708330d81262', '1202.001', 'Piutang Asuransi Sinarmas', 'ASSET', 'DEBIT', '876554c3-32f7-5f86-8119-723b3303594c', 4, FALSE, TRUE, 'Akun BCS 1202.001 - Piutang Asuransi Sinarmas')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('f9161547-385b-5c6c-82fe-68f0ea6667f6', '1202.002', 'Piutang Asuransi Prudential', 'ASSET', 'DEBIT', '876554c3-32f7-5f86-8119-723b3303594c', 4, FALSE, TRUE, 'Akun BCS 1202.002 - Piutang Asuransi Prudential')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('fc4ec513-47d1-5e3f-befb-322b4f2ff194', '1202.003', 'Piutang MCU Eksekutif Korporat', 'ASSET', 'DEBIT', '876554c3-32f7-5f86-8119-723b3303594c', 4, FALSE, TRUE, 'Akun BCS 1202.003 - Piutang MCU Eksekutif Korporat')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('38fa59e7-6f95-5bcd-9bdc-982a7879a3d9', '2101.001', 'Uang Muka Reservasi Poli Spesialis Anak Kidinest', 'LIABILITY', 'CREDIT', '4f31f30a-cb1c-5b5c-904a-c482cc330fb5', 4, FALSE, TRUE, 'Akun BCS 2101.001 - Uang Muka Reservasi Poli Spesialis Anak Kidinest')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('b6092a83-23ad-5309-b50c-6910557275e3', '2101.002', 'Down Payment Treatment Laser & Estetika SSC', 'LIABILITY', 'CREDIT', '4f31f30a-cb1c-5b5c-904a-c482cc330fb5', 4, FALSE, TRUE, 'Akun BCS 2101.002 - Down Payment Treatment Laser & Estetika SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('d71d245c-b2d1-519c-b1e2-c8d58177af37', '2102.001', 'Deposit Pasien Paket Facial & Skincare SSC', 'LIABILITY', 'CREDIT', '0819ad8c-3476-5e21-a502-1a5ae5c3fbb3', 4, FALSE, TRUE, 'Akun BCS 2102.001 - Deposit Pasien Paket Facial & Skincare SSC')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('e7e906e2-d618-5f09-ba17-1eba2126bd76', '2102.002', 'Deposit Pasien Terapi Tumbuh Kembang Kidinest', 'LIABILITY', 'CREDIT', '0819ad8c-3476-5e21-a502-1a5ae5c3fbb3', 4, FALSE, TRUE, 'Akun BCS 2102.002 - Deposit Pasien Terapi Tumbuh Kembang Kidinest')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('96cfe025-5cfa-515b-af3e-5d16b14c730a', '2501.001', 'PPh 21 Dokter Tenaga Medis Ahli (PMK 168/2023)', 'LIABILITY', 'CREDIT', '333eb5dd-f1a6-5620-ab25-c6f41c6944da', 4, FALSE, TRUE, 'Akun BCS 2501.001 - PPh 21 Dokter Tenaga Medis Ahli (PMK 168/2023)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO accounts (id, code, name, type, normal_balance, parent_id, level, is_header, is_active, description)
VALUES ('345c619a-559a-5fdd-a9ca-ade2673cc7b4', '2501.002', 'PPh 21 Karyawan & Staf Tetap (TER PP 58)', 'LIABILITY', 'CREDIT', '333eb5dd-f1a6-5620-ab25-c6f41c6944da', 4, FALSE, TRUE, 'Akun BCS 2501.002 - PPh 21 Karyawan & Staf Tetap (TER PP 58)')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    parent_id = EXCLUDED.parent_id,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    is_active = EXCLUDED.is_active,
    description = EXCLUDED.description;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('e0f4f5e5-ee90-5219-bdb9-a1381db31d54', '11111111-1111-1111-1111-111111111111', '77dd7a90-8ae3-59b8-a65d-d76f38628738', 'KAS-KIDINEST', 'Kas Kasir Kidinest', 'CASH', '', '', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('c4df7b77-567c-50ee-8c86-04043c9cfa43', '11111111-1111-1111-1111-111111111111', '444d0bee-6bab-5a11-b703-dc216563e8d5', 'KAS-SSC', 'Kas Kasir Semarang Skin Center (SSC)', 'CASH', '', '', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('2198fd22-3e43-5522-9585-79eee15393f0', '11111111-1111-1111-1111-111111111111', 'ccb4cf53-4b16-5cd6-b5c6-600d2724dc50', 'PETTY-CASH', 'Kas Kecil (Petty Cash Operasional)', 'CASH', '', '', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('7e739988-4ebd-502c-bdb5-c5f1cf456fe6', '11111111-1111-1111-1111-111111111111', '7ddc4922-b987-5268-a704-d53880f59d75', 'BCA-OPS', 'Bank BCA Operasional (Giro & EDC)', 'BANK', 'Bank Central Asia (BCA)', '8830-112-990', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('55d86b49-805b-5710-8668-94940e9aa34a', '11111111-1111-1111-1111-111111111111', '0033a4e5-ae85-5a17-a004-4ba0f8b26f15', 'MANDIRI-OPS', 'Bank Mandiri Operasional & Payroll', 'BANK', 'Bank Mandiri', '137-00-998811-2', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_accounts (id, branch_id, account_id, code, name, type, bank_name, account_number, is_active)
VALUES ('cd073a70-4871-5ae2-87fc-c35bdcdba668', '11111111-1111-1111-1111-111111111111', 'cf3a16d6-659e-5911-b80d-bf3bf527daba', 'QRIS-TRANSIT', 'Settlement QRIS & Payment Gateway Transit', 'BANK', 'Payment Gateway Transit', 'ID-QRIS-99210', TRUE)
ON CONFLICT (code) DO UPDATE SET
    account_id = EXCLUDED.account_id,
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    bank_name = EXCLUDED.bank_name,
    account_number = EXCLUDED.account_number,
    is_active = EXCLUDED.is_active;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    '01cf72b4-d144-5c0f-9d25-ed1d582353ee', '11111111-1111-1111-1111-111111111111', '2198fd22-3e43-5522-9585-79eee15393f0', 'BKM-2026-03-001', 'CASH_IN', 3000000,
    '2026-03-01', '7ddc4922-b987-5268-a704-d53880f59d75', 'Pengisian awal dana kas kecil operasional klinik bulan Maret', 'CHQ-BCA-88912'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    'e26bf461-64ee-5899-8a54-f4d51a2dd040', '11111111-1111-1111-1111-111111111111', '2198fd22-3e43-5522-9585-79eee15393f0', 'BKK-2026-03-001', 'CASH_OUT', 450000,
    '2026-03-02', '5d49dfeb-c8e9-5cff-a211-c53e072dc7bc', 'Pembelian 10 galon air mineral & konsumsi rapat koordinasi dokter poli', 'NOTA-SB-3312'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    '9239e10e-3859-5791-bfb2-8a544dc77198', '11111111-1111-1111-1111-111111111111', '2198fd22-3e43-5522-9585-79eee15393f0', 'BKK-2026-03-002', 'CASH_OUT', 275000,
    '2026-03-03', '5d49dfeb-c8e9-5cff-a211-c53e072dc7bc', 'Pembelian materai 10.000 & map formulir rekam medis', 'KW-GRM-901'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    'b85e2878-49fd-5d3a-a97b-9ff367a9efb6', '11111111-1111-1111-1111-111111111111', '7e739988-4ebd-502c-bdb5-c5f1cf456fe6', 'BBM-2026-03-001', 'CASH_IN', 45000000,
    '2026-03-03', 'ce371684-0182-5a22-bb10-708330d81262', 'Penerimaan klaim asuransi Sinarmas pasien rawat jalan poli anak', 'TRF-SINAR-8902'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    '0dc03135-7695-51b0-b646-1a6c1f58118c', '11111111-1111-1111-1111-111111111111', '7e739988-4ebd-502c-bdb5-c5f1cf456fe6', 'BBK-2026-03-001', 'CASH_OUT', 25000000,
    '2026-03-04', 'df80c80e-93f1-5212-a2d6-1c44d2ea15ab', 'Pelunasan faktur supplier PT Kimia Farma pengadaan vaksin hepatitis B & DPT', 'INV-KF-99120'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    '2aba8284-df88-58d2-96a7-d498bd14af56', '11111111-1111-1111-1111-111111111111', '7e739988-4ebd-502c-bdb5-c5f1cf456fe6', 'TRF-2026-03-001', 'TRANSFER', 10000000,
    '2026-03-05', '0033a4e5-ae85-5a17-a004-4ba0f8b26f15', 'Transfer likuiditas antar bank untuk persiapan payroll gaji perawat', 'TRF-INTERN-01'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO cash_bank_transactions (
    id, branch_id, cash_bank_account_id, transaction_number, transaction_type, amount,
    transaction_date, counterpart_account_id, description, reference_no
) VALUES (
    '65316d92-9749-500c-bbcf-62ae9c936c89', '11111111-1111-1111-1111-111111111111', '7e739988-4ebd-502c-bdb5-c5f1cf456fe6', 'BBM-2026-03-002', 'CASH_IN', 1250000,
    '2026-03-06', 'f9421b35-56a5-5171-9a83-6630ca20a118', 'Pendapatan jasa giro rekening koran operasional BCA bulan Februari 2026', 'GIRO-BCA-FEB26'
) ON CONFLICT (transaction_number) DO UPDATE SET
    amount = EXCLUDED.amount,
    description = EXCLUDED.description,
    reference_no = EXCLUDED.reference_no;
INSERT INTO journal_entries (
    id, branch_id, journal_number, journal_type, transaction_date, posting_date, status,
    reference_type, reference_id, total_debit, total_credit, description
) VALUES (
    'f7771332-2d92-5d66-b524-f905bad28260', '11111111-1111-1111-1111-111111111111', 'JV-2026-03-001', 'GENERAL', '2026-03-01', '2026-03-01 12:00:00+07', 'POSTED',
    'CASH_BANK', 'BKM-2026-03-001', 3000000, 3000000, 'Pengisian awal dana kas kecil operasional klinik bulan Maret'
) ON CONFLICT (journal_number) DO UPDATE SET
    total_debit = EXCLUDED.total_debit,
    total_credit = EXCLUDED.total_credit,
    description = EXCLUDED.description;
DELETE FROM journal_lines WHERE journal_entry_id = 'f7771332-2d92-5d66-b524-f905bad28260';
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    'a80bcc98-f5de-505f-a4c9-77e36f7fed56', 'f7771332-2d92-5d66-b524-f905bad28260', 'ccb4cf53-4b16-5cd6-b5c6-600d2724dc50', 3000000, 0, 'Debet Kas Kecil Operasional', 1
);
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    'e8046cdb-ab68-5146-ba3c-667ecff901c4', 'f7771332-2d92-5d66-b524-f905bad28260', '7ddc4922-b987-5268-a704-d53880f59d75', 0, 3000000, 'Kredit Bank BCA Operasional', 2
);
INSERT INTO journal_entries (
    id, branch_id, journal_number, journal_type, transaction_date, posting_date, status,
    reference_type, reference_id, total_debit, total_credit, description
) VALUES (
    '8e0f3862-0b25-51f2-89d4-9fc75a60f0c5', '11111111-1111-1111-1111-111111111111', 'JV-2026-03-002', 'GENERAL', '2026-03-02', '2026-03-02 12:00:00+07', 'POSTED',
    'RECEIVABLE', 'BBM-2026-03-001', 45000000, 45000000, 'Penerimaan klaim asuransi Sinarmas pasien poli anak'
) ON CONFLICT (journal_number) DO UPDATE SET
    total_debit = EXCLUDED.total_debit,
    total_credit = EXCLUDED.total_credit,
    description = EXCLUDED.description;
DELETE FROM journal_lines WHERE journal_entry_id = '8e0f3862-0b25-51f2-89d4-9fc75a60f0c5';
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    '1338cd02-a783-5078-922f-63397ba7b05d', '8e0f3862-0b25-51f2-89d4-9fc75a60f0c5', '7ddc4922-b987-5268-a704-d53880f59d75', 45000000, 0, 'Debet Bank BCA Pelunasan Asuransi', 1
);
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    '332659e4-3fef-5177-8d81-389c150eac73', '8e0f3862-0b25-51f2-89d4-9fc75a60f0c5', 'ce371684-0182-5a22-bb10-708330d81262', 0, 45000000, 'Kredit Piutang Asuransi Sinarmas', 2
);
INSERT INTO journal_entries (
    id, branch_id, journal_number, journal_type, transaction_date, posting_date, status,
    reference_type, reference_id, total_debit, total_credit, description
) VALUES (
    'aec4cba4-069e-5700-a63f-6ef584c7d371', '11111111-1111-1111-1111-111111111111', 'JV-2026-03-003', 'GENERAL', '2026-03-04', '2026-03-04 12:00:00+07', 'POSTED',
    'PAYABLE', 'BBK-2026-03-001', 25000000, 25000000, 'Pelunasan faktur supplier PT Kimia Farma pengadaan vaksin hepatitis B & DPT'
) ON CONFLICT (journal_number) DO UPDATE SET
    total_debit = EXCLUDED.total_debit,
    total_credit = EXCLUDED.total_credit,
    description = EXCLUDED.description;
DELETE FROM journal_lines WHERE journal_entry_id = 'aec4cba4-069e-5700-a63f-6ef584c7d371';
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    'efab8a91-a2b8-56ac-a9e7-0c2bd40794ba', 'aec4cba4-069e-5700-a63f-6ef584c7d371', 'df80c80e-93f1-5212-a2d6-1c44d2ea15ab', 25000000, 0, 'Debet Hutang Supplier Obat & Farmasi', 1
);
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    '1f1e60c6-b47e-592e-a582-1dbfa8affb90', 'aec4cba4-069e-5700-a63f-6ef584c7d371', '7ddc4922-b987-5268-a704-d53880f59d75', 0, 25000000, 'Kredit Bank BCA Pembayaran Supplier', 2
);
INSERT INTO journal_entries (
    id, branch_id, journal_number, journal_type, transaction_date, posting_date, status,
    reference_type, reference_id, total_debit, total_credit, description
) VALUES (
    '8a9635b6-cf86-5cc1-b7f2-d93c20a8b7c6', '11111111-1111-1111-1111-111111111111', 'JV-2026-03-004', 'GENERAL', '2026-03-05', '2026-03-05 12:00:00+07', 'POSTED',
    'TRANSFER', 'TRF-2026-03-001', 10000000, 10000000, 'Transfer likuiditas antar bank BCA ke Bank Mandiri untuk payroll'
) ON CONFLICT (journal_number) DO UPDATE SET
    total_debit = EXCLUDED.total_debit,
    total_credit = EXCLUDED.total_credit,
    description = EXCLUDED.description;
DELETE FROM journal_lines WHERE journal_entry_id = '8a9635b6-cf86-5cc1-b7f2-d93c20a8b7c6';
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    'b79c6943-a318-5eea-8c0c-ba3f73fed52d', '8a9635b6-cf86-5cc1-b7f2-d93c20a8b7c6', '0033a4e5-ae85-5a17-a004-4ba0f8b26f15', 10000000, 0, 'Debet Bank Mandiri Operasional', 1
);
INSERT INTO journal_lines (
    id, journal_entry_id, account_id, debit, credit, description, line_order
) VALUES (
    'd9c31dd7-2964-5da0-9bb2-92ded563a46b', '8a9635b6-cf86-5cc1-b7f2-d93c20a8b7c6', '7ddc4922-b987-5268-a704-d53880f59d75', 0, 10000000, 'Kredit Bank BCA Operasional', 2
);
INSERT INTO tax_withholdings (
    id, branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name,
    object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported
) VALUES (
    'aa5b878f-0ff0-5575-bb94-6143859e7a40', '11111111-1111-1111-1111-111111111111', 'PPH23', 'BP-23-202601-9901', '3374011202800001', '01.234.567.8-501.000', 'PT Wastec International (Pengolahan Limbah Medis B3)',
    '28-100-02', 120000000, 120000000, 0.02, 2400000, 2026, 1, TRUE
) ON CONFLICT (bupot_number) DO UPDATE SET
    taxpayer_name = EXCLUDED.taxpayer_name,
    gross_amount = EXCLUDED.gross_amount,
    dpp_amount = EXCLUDED.dpp_amount,
    tax_amount = EXCLUDED.tax_amount,
    is_reported = EXCLUDED.is_reported;
INSERT INTO tax_withholdings (
    id, branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name,
    object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported
) VALUES (
    'd5831282-f493-5605-a1c7-f3a5d8f186bb', '11111111-1111-1111-1111-111111111111', 'PPH23', 'BP-23-202602-9902', '3374011202800002', '01.234.567.8-501.001', 'PT Laserindo Medika (Jasa Kalibrasi Alat Medis Laser SSC)',
    '28-100-02', 150000000, 150000000, 0.02, 3000000, 2026, 2, TRUE
) ON CONFLICT (bupot_number) DO UPDATE SET
    taxpayer_name = EXCLUDED.taxpayer_name,
    gross_amount = EXCLUDED.gross_amount,
    dpp_amount = EXCLUDED.dpp_amount,
    tax_amount = EXCLUDED.tax_amount,
    is_reported = EXCLUDED.is_reported;
INSERT INTO tax_withholdings (
    id, branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name,
    object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported
) VALUES (
    'f4b059c5-d3af-5902-a984-9d68a1a4b36c', '11111111-1111-1111-1111-111111111111', 'PPH23', 'BP-23-202603-9903', '3374011202800003', '01.234.567.8-501.002', 'PT Cloud Medika Asia (Pemeliharaan SIMKlinik & Server Cloud)',
    '28-100-02', 100000000, 100000000, 0.02, 2000000, 2026, 3, TRUE
) ON CONFLICT (bupot_number) DO UPDATE SET
    taxpayer_name = EXCLUDED.taxpayer_name,
    gross_amount = EXCLUDED.gross_amount,
    dpp_amount = EXCLUDED.dpp_amount,
    tax_amount = EXCLUDED.tax_amount,
    is_reported = EXCLUDED.is_reported;
INSERT INTO tax_withholdings (
    id, branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name,
    object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported
) VALUES (
    'a4d9286b-c69b-5645-b6ac-c081c5464157', '11111111-1111-1111-1111-111111111111', 'PPH21_DOCTOR', 'BP-21-202601-01', '3374012304850001', '07.123.456.7-501.000', 'dr. Aris Munandar, Sp.A (K)',
    '21-100-07', 45000000, 22500000, 0.05, 1125000, 2026, 1, TRUE
) ON CONFLICT (bupot_number) DO UPDATE SET
    taxpayer_name = EXCLUDED.taxpayer_name,
    gross_amount = EXCLUDED.gross_amount,
    dpp_amount = EXCLUDED.dpp_amount,
    tax_amount = EXCLUDED.tax_amount,
    is_reported = EXCLUDED.is_reported;
INSERT INTO tax_withholdings (
    id, branch_id, tax_type, bupot_number, taxpayer_nik, taxpayer_npwp, taxpayer_name,
    object_code, gross_amount, dpp_amount, tax_rate, tax_amount, period_year, period_month, is_reported
) VALUES (
    '3ad8a43f-3b2e-5cdf-beae-ffb2afd3ace9', '11111111-1111-1111-1111-111111111111', 'PPH21_DOCTOR', 'BP-21-202601-02', '3374015607900002', '08.234.567.8-501.000', 'dr. Citra Lestari, Sp.D.V.E (SSC)',
    '21-100-07', 58000000, 29000000, 0.05, 1450000, 2026, 1, TRUE
) ON CONFLICT (bupot_number) DO UPDATE SET
    taxpayer_name = EXCLUDED.taxpayer_name,
    gross_amount = EXCLUDED.gross_amount,
    dpp_amount = EXCLUDED.dpp_amount,
    tax_amount = EXCLUDED.tax_amount,
    is_reported = EXCLUDED.is_reported;
COMMIT;
