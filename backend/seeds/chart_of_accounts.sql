-- Chart of Accounts (COA) for Dual-Brand Clinic:
-- 1. Kidinest (Pediatrik, Vaksinasi & Tumbuh Kembang)
-- 2. Semarang Skin Center (SSC - Estetika & Dermatologi)

INSERT INTO accounts (code, name, type, normal_balance, level, is_header, is_active, description)
VALUES
-- ==========================================
-- 1. AKTIVA / ASSETS (DEBIT)
-- ==========================================
('1-0000', 'AKTIVA / ASSETS', 'ASSET', 'DEBIT', 1, TRUE, TRUE, 'Induk Aset'),
('1-1000', 'Aset Lancar', 'ASSET', 'DEBIT', 2, TRUE, TRUE, 'Aset Likuid dan Lancar'),
('1-1110', 'Kas Kecil (Petty Cash Operasional)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Kecil Pengeluaran Harian'),
('1-1111', 'Kas Kasir Kidinest', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Penerimaan Kasir Unit Kidinest (Anak)'),
('1-1112', 'Kas Kasir Semarang Skin Center (SSC)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Kas Penerimaan Kasir Unit SSC (Estetika)'),
('1-1113', 'Bank BCA Operasional', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Rekening Operasional BCA Klinik'),
('1-1114', 'Bank Mandiri Operasional', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Rekening Operasional Mandiri Klinik'),
('1-1120', 'Kliring Payment Gateway', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Akun Perantara Settlement Gateway [Deferred]'),

('1-1200', 'Piutang Pasien & Klaim Asuransi', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Piutang Layanan Medis Belum Tertagih'),
('1-1210', 'Piutang Klaim BPJS Kesehatan', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Piutang Klaim BPJS'),

-- Persediaan Terpisah per Unit
('1-1300', 'Persediaan Farmasi & Medis', 'ASSET', 'DEBIT', 2, TRUE, TRUE, 'Induk Persediaan Medis'),
('1-1310', 'Persediaan Obat & Sirup Anak (Kidinest)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Stok Obat Racikan dan Paten Pediatrik'),
('1-1320', 'Persediaan Vaksin & Imunisasi (Kidinest)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Stok Vaksin Wajib & Tambahan Anak'),
('1-1330', 'Persediaan Produk Skincare & Sunscreen (SSC)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Stok Krim Wajah, Serum, dan Produk Homecare SSC'),
('1-1340', 'Persediaan BMHP Estetika & Facial (SSC)', 'ASSET', 'DEBIT', 3, FALSE, TRUE, 'Bahan Medis Habis Pakai Treatment Estetika'),

-- ==========================================
-- 2. KEWAJIBAN / LIABILITIES (CREDIT)
-- ==========================================
('2-0000', 'KEWAJIBAN / LIABILITIES', 'LIABILITY', 'CREDIT', 1, TRUE, TRUE, 'Induk Kewajiban'),
('2-1000', 'Hutang Lancar', 'LIABILITY', 'CREDIT', 2, TRUE, TRUE, 'Kewajiban Jangka Pendek'),
('2-1110', 'Hutang Supplier Obat & Vaksin (Kidinest)', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Hutang Pembelian Obat/Vaksin Pediatrik'),
('2-1115', 'Hutang Supplier Skincare & BMHP (SSC)', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Hutang Pembelian Produk Skincare SSC'),
('2-1120', 'Hutang Jasa Medis Dokter Anak (Sp.A)', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Bagi Hasil Honor Dokter Spesialis Anak Kidinest'),
('2-1125', 'Hutang Jasa Medis Dokter Kulit & Terapis (SSC)', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Bagi Hasil Honor Dokter Kulit & Terapis SSC'),
('2-1130', 'Hutang Gaji & Karyawan', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Akrual Beban Gaji Karyawan'),
('2-1140', 'Uang Muka Pasien / Deposit', 'LIABILITY', 'CREDIT', 3, FALSE, TRUE, 'Deposit & Uang Muka Tindakan Pasien'),

-- ==========================================
-- 3. EKUITAS / EQUITY (CREDIT)
-- ==========================================
('3-0000', 'MODAL / EQUITY', 'EQUITY', 'CREDIT', 1, TRUE, TRUE, 'Induk Ekuitas'),
('3-1100', 'Modal Disetor', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Modal Pemilik Klinik'),
('3-1200', 'Laba Ditahan (Retained Earnings)', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Akumulasi Laba Periode Lalu'),
('3-1300', 'Laba Periode Berjalan', 'EQUITY', 'CREDIT', 2, FALSE, TRUE, 'Laba Bersih Tahun/Bulan Berjalan'),

-- ==========================================
-- 4. PENDAPATAN / REVENUE (CREDIT)
-- ==========================================
('4-0000', 'PENDAPATAN / REVENUE', 'REVENUE', 'CREDIT', 1, TRUE, TRUE, 'Induk Pendapatan'),

-- Unit Kidinest (Pediatrik & Tumbuh Kembang)
('4-1000', 'PENDAPATAN UNIT KIDINEST (ANAK)', 'REVENUE', 'CREDIT', 2, TRUE, TRUE, 'Induk Pendapatan Unit Kidinest'),
('4-1100', 'Pendapatan Pendaftaran & Skrining Anak', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Biaya Registrasi Poli Anak & Skrining'),
('4-1200', 'Pendapatan Jasa Dokter Spesialis Anak (Sp.A)', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Konsultasi Dokter Spesialis Anak'),
('4-1300', 'Pendapatan Vaksinasi & Imunisasi Anak', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Layanan Suntik Vaksin Anak (Wajib/Pilihan)'),
('4-1400', 'Pendapatan Pemantauan Tumbuh Kembang Anak', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Skrining Denver II & Paket Tumbuh Kembang'),
('4-1500', 'Pendapatan Penjualan Obat & Racikan Anak', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Resep Obat Sirup & Racikan Farmasi Anak'),

-- Unit Semarang Skin Center (SSC - Estetika & Dermatologi)
('4-2000', 'PENDAPATAN SEMARANG SKIN CENTER (SSC)', 'REVENUE', 'CREDIT', 2, TRUE, TRUE, 'Induk Pendapatan Unit SSC'),
('4-2100', 'Pendapatan Konsultasi Dokter Spesialis Kulit (Sp.DVE)', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Konsultasi Dokter Kulit SSC'),
('4-2200', 'Pendapatan Treatment Facial & Estetika Medis', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Perawatan Facial, Microdermabrasi, dll'),
('4-2300', 'Pendapatan Tindakan Laser, Cauter & Peeling', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Tindakan Medis Kulit & Laser'),
('4-2400', 'Pendapatan Penjualan Produk Skincare & Sunscreen', 'REVENUE', 'CREDIT', 3, FALSE, TRUE, 'Penjualan Produk Kosmetik & Skincare SSC'),

-- ==========================================
-- 5. HARGA POKOK PENJUALAN / COGS (DEBIT)
-- ==========================================
('5-0000', 'HARGA POKOK PENJUALAN (HPP)', 'COGS', 'DEBIT', 1, TRUE, TRUE, 'Induk HPP'),
('5-1100', 'HPP Obat & Racikan Anak (Kidinest)', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Beban Pokok Obat Farmasi Anak'),
('5-1200', 'HPP Vaksin & Imunisasi Anak (Kidinest)', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Beban Pokok Ampul Vaksin Anak'),
('5-1300', 'HPP Produk Skincare & Homecare (SSC)', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Beban Pokok Produk Krim Skincare SSC'),
('5-1400', 'HPP Bahan BMHP & Treatment Estetika (SSC)', 'COGS', 'DEBIT', 2, FALSE, TRUE, 'Bahan Habis Pakai Tindakan Estetika'),

-- ==========================================
-- 6. BEBAN OPERASIONAL / EXPENSES (DEBIT)
-- ==========================================
('6-0000', 'BEBAN OPERASIONAL / EXPENSES', 'EXPENSE', 'DEBIT', 1, TRUE, TRUE, 'Induk Beban'),
('6-1100', 'Beban Honor Jasa Dokter Spesialis Anak (Sp.A)', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Bagi Hasil Tindakan & Konsul Dokter Anak'),
('6-1110', 'Beban Honor Jasa Dokter Kulit & Terapis Estetika', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Bagi Hasil Dokter Spesialis Kulit & Terapis SSC'),
('6-1200', 'Beban Gaji Karyawan, Perawat & Staf', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Gaji Perawat dan Karyawan Backoffice'),
('6-1300', 'Beban Listrik, Air & Telepon', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Utilitas Fasilitas Klinik'),
('6-1400', 'Beban Administrasi Bank & EDC', 'EXPENSE', 'DEBIT', 2, FALSE, TRUE, 'Biaya Admin Bank & Transaksi')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    normal_balance = EXCLUDED.normal_balance,
    level = EXCLUDED.level,
    is_header = EXCLUDED.is_header,
    description = EXCLUDED.description;
