-- Seed Standard Chart of Accounts (COA) for Clinic ERP

INSERT INTO accounts (code, name, type, normal_balance, level, is_header, is_active, description)
VALUES
-- 1. ASSETS
('1-0000', 'AKTIVA / ASSETS', 'ASSET', 'DEBIT', 1, TRUE, TRUE, 'Induk Aset'),
('1-1000', 'Aset Lancar', 'ASSET', 'DEBIT', 2, TRUE, TRUE, 'Aset Likuid dan Lancar'),
('1-1110', 'Kas Kecil (Petty Cash)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Kecil Operasional'),
('1-1111', 'Kas Kasir', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Penerimaan Kasir Pasien'),
('1-1112', 'Kas Operasional', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Pengeluaran Operasional'),
('1-1113', 'Bank BCA', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Rekening Operasional BCA'),
('1-1114', 'Bank Mandiri', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Rekening Operasional Mandiri'),
('1-1120', 'Kliring Payment Gateway', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Akun Perantara Settlement Gateway [Deferred]'),
('1-1200', 'Piutang Pasien', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Piutang Pelayanan Pasien'),
('1-1210', 'Piutang Asuransi & Korporat', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Piutang Klaim Asuransi/Perusahaan'),
('1-1310', 'Persediaan Obat', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Persediaan Farmasi & Obat'),
('1-1320', 'Persediaan Alat Kesehatan & BMHP', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Bahan Medis Habis Pakai'),

-- 2. LIABILITIES
('2-0000', 'KEWAJIBAN / LIABILITIES', 'LIABILITY', 'CREDIT', 1, TRUE, TRUE, 'Induk Kewajiban'),
('2-1110', 'Hutang Dagang / Supplier Obat', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Hutang Pembelian Obat Supplier'),
('2-1120', 'Hutang Jasa Medis Dokter', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Bagi Hasil / Fee Dokter Belum Dibayar'),
('2-1130', 'Hutang Gaji & Karyawan', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Akrual Beban Gaji'),
('2-1140', 'Uang Muka Pasien / Deposit', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Titipan Uang Muka Pasien'),

-- 3. EQUITY
('3-0000', 'MODAL / EQUITY', 'EQUITY', 'CREDIT', 1, TRUE, TRUE, 'Induk Ekuitas'),
('3-1100', 'Modal Disetor', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Modal Pemilik Klinik'),
('3-1200', 'Laba Ditahan (Retained Earnings)', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Akumulasi Laba Periode Lalu'),
('3-1300', 'Laba Periode Berjalan', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Laba Bersih Tahun/Bulan Berjalan'),

-- 4. REVENUE
('4-0000', 'PENDAPATAN / REVENUE', 'REVENUE', 'CREDIT', 1, TRUE, TRUE, 'Induk Pendapatan'),
('4-1100', 'Pendapatan Pendaftaran & Admin', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Biaya Registrasi Pasien'),
('4-1200', 'Pendapatan Jasa Konsultasi & Dokter', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Biaya Konsultasi & Tindakan Dokter'),
('4-1300', 'Pendapatan Penjualan Obat Farmasi', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Penjualan Obat Resep dan Non-Resep'),
('4-1400', 'Pendapatan Tindakan Medis & Laboratorium', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Tindakan Keperawatan / Lab'),

-- 5. COST OF GOODS SOLD (COGS / HPP)
('5-0000', 'HARGA POKOK PENJUALAN (HPP)', 'COGS', 'DEBIT', 1, TRUE, TRUE, 'Induk HPP'),
('5-1100', 'HPP Obat & Farmasi', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Harga Pokok Obat Terjual'),
('5-1200', 'HPP Bahan Medis Habis Pakai (BMHP)', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Beban Penggunaan BMHP Tindakan'),

-- 6. EXPENSES (BEBAN OPERASIONAL)
('6-0000', 'BEBAN OPERASIONAL / EXPENSES', 'EXPENSE', 'DEBIT', 1, TRUE, TRUE, 'Induk Beban'),
('6-1100', 'Beban Gaji Karyawan', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Gaji Staf & Perawat'),
('6-1200', 'Beban Jasa / Honor Dokter', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Honor Dokter Tindakan & Konsultasi'),
('6-1300', 'Beban Listrik, Air & Telepon', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Utilitas Klinik'),
('6-1400', 'Beban Administrasi Bank & Clearing', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Biaya Admin Bank & Transaksi')
ON CONFLICT (code) DO NOTHING;

