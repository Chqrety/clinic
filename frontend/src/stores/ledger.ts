import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Account, LedgerEntry, TrialBalanceRow, AccountingPeriod } from '@/types/accounting'
import { useJournalsStore } from './journals'
import { useAccountingStore } from './accounting'
import { accountingApi } from '@/api/accounting'

export const useLedgerStore = defineStore('ledger', () => {
  const journalsStore = useJournalsStore()
  const accountingStore = useAccountingStore()

  // Standard Chart of Accounts (COA) with 4-Tier Indexing Structure & Tax Integration (BCS)
  const accounts = ref<Account[]>([
    // ==========================================
    // 1 - ASET - BCS (Rp 26.461.564.553,84 Dr)
    // ==========================================
    { id: '1', code: '1', name: 'ASET - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 26461564553.84 },
    
    // 11 - KAS DAN BANK - BCS (Rp 1.277.929.020,88 Dr)
    { id: '11', code: '11', name: 'KAS DAN BANK - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 1277929020.88 },
    { id: '1101', code: '1101', name: 'Kas Kasir & Operasional - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '11', currentBalance: 88600000 },
    { id: '1101.001', code: '1101.001', name: 'Kas Kasir Kidinest', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1101', currentBalance: 45200000 },
    { id: '1101.002', code: '1101.002', name: 'Kas Kasir Semarang Skin Center (SSC)', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1101', currentBalance: 38400000 },
    { id: '1101.003', code: '1101.003', name: 'Kas Kecil (Petty Cash Operasional)', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1101', currentBalance: 5000000 },
    { id: '1102', code: '1102', name: 'Bank Operasional & Kliring - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '11', currentBalance: 1189329020.88 },
    { id: '1102.001', code: '1102.001', name: 'Bank BCA Operasional (Giro & EDC)', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1102', currentBalance: 680964487.55 },
    { id: '1102.002', code: '1102.002', name: 'Bank Mandiri Operasional & Payroll', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1102', currentBalance: 493864533.33 },
    { id: '1102.003', code: '1102.003', name: 'Settlement QRIS & Payment Gateway Transit', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1102', currentBalance: 14500000 },

    // 12 - PIUTANG - BCS (Rp 1.568.434.065,00 Dr)
    { id: '12', code: '12', name: 'PIUTANG - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 1568434065 },
    { id: '1201', code: '1201', name: 'Piutang Pasien Perawatan Rawat Jalan & Treatment', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '12', currentBalance: 485600000 },
    { id: '1202', code: '1202', name: 'Piutang Klaim Asuransi & Korporat', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '12', currentBalance: 1097834065 },
    { id: '1202.001', code: '1202.001', name: 'Piutang Asuransi Sinarmas', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1202', currentBalance: 420000000 },
    { id: '1202.002', code: '1202.002', name: 'Piutang Asuransi Prudential', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '1202', currentBalance: 465000000 },
    { id: '1202.003', code: '1202.003', name: 'Piutang MCU Eksekutif Korporat', type: 'ASSET', normalBalance: 'DEBIT', level: 4, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '1202', currentBalance: 212834065 },
    { id: '1203', code: '1203', name: 'Cadangan Kerugian Piutang Tak Tertagih', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '12', currentBalance: -15000000 },

    // 13 - Persediaan - BCS (Rp 19.335.748.296,88 Dr)
    { id: '13', code: '13', name: 'Persediaan - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 19335748296.88 },
    { id: '1301', code: '1301', name: 'Persediaan Obat Resep & Paten Farmasi Kidinest', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '13', currentBalance: 4885000000 },
    { id: '1302', code: '1302', name: 'Persediaan Produk Skincare & Kosmetika SSC', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '13', currentBalance: 6442748296.88 },
    { id: '1303', code: '1303', name: 'Persediaan Bahan Medis Habis Pakai (BMHP)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '13', currentBalance: 3208000000 },
    { id: '1304', code: '1304', name: 'Persediaan Vaksin & Serum Imunisasi Anak', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '13', currentBalance: 4800000000 },

    // 14 - INVESTASI - BCS (Rp 3.234.375.045,00 Dr)
    { id: '14', code: '14', name: 'INVESTASI - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 3234375045 },
    { id: '1401', code: '1401', name: 'Investasi Tanah & Bangunan Cabang Baru', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '14', currentBalance: 2500000000 },
    { id: '1402', code: '1402', name: 'Investasi Pengembangan Riset Formulasi Skincare', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '14', currentBalance: 734375045 },

    // 15 - BIAYA DIBAYAR DI MUKA - BCS (Rp 2.250.000,00 Dr)
    { id: '15', code: '15', name: 'BIAYA DIBAYAR DI MUKA - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 2250000 },
    { id: '1501', code: '1501', name: 'Sewa Ruko & Fasilitas Dibayar di Muka', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '15', currentBalance: 2250000 },
    { id: '1504', code: '1504', name: 'Uang Muka PPh 25 & Pajak Masukan (Kredit Pajak)', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '15', currentBalance: 0, taxCategory: 'PPH_25' },

    // 16 - AKTIVA TETAP - BCS (Rp 1.042.828.126,08 Dr)
    { id: '16', code: '16', name: 'AKTIVA TETAP - BCS', type: 'ASSET', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 1042828126.08 },
    { id: '1601', code: '1601', name: 'Bangunan & Renovasi Interior Klinik', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '16', currentBalance: 350000000 },
    { id: '1602', code: '1602', name: 'Mesin Laser Nd:YAG & Alkes Estetika SSC', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '16', currentBalance: 500000000 },
    { id: '1603', code: '1603', name: 'Peralatan USG 4D & Poli Anak Kidinest', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '16', currentBalance: 250000000 },
    { id: '1604', code: '1604', name: 'Inventaris Komputer, POS Kasir & Jaringan', type: 'ASSET', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '16', currentBalance: 97828126.08 },
    { id: '1606', code: '1606', name: 'Akumulasi Penyusutan Aktiva Tetap', type: 'ASSET', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '16', currentBalance: 155000000 },
    { id: '1999', code: '1999', name: 'Temporary Opening - BCS', type: 'ASSET', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '1', currentBalance: 0 },

    // ==========================================
    // 2 - KEWAJIBAN - BCS (Rp 9.962.073.223,05 Cr)
    // ==========================================
    { id: '2', code: '2', name: 'KEWAJIBAN - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 9962073223.05 },
    
    // 21 - UANG MUKA & TITIPAN - BCS (Rp 2.147.099.360,00 Cr)
    { id: '21', code: '21', name: 'UANG MUKA & TITIPAN - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 2147099360 },
    { id: '2101', code: '2101', name: 'UANG TANDA JADI', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '21', currentBalance: 35500000 },
    { id: '2101.001', code: '2101.001', name: 'Uang Muka Reservasi Poli Spesialis Anak Kidinest', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '2101', currentBalance: 15500000 },
    { id: '2101.002', code: '2101.002', name: 'Down Payment Treatment Laser & Estetika SSC', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2101', currentBalance: 20000000 },
    { id: '2102', code: '2102', name: 'UANG MUKA PENJUALAN', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '21', currentBalance: 1904190360 },
    { id: '2102.001', code: '2102.001', name: 'Deposit Pasien Paket Facial & Skincare SSC', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '2102', currentBalance: 980190360 },
    { id: '2102.002', code: '2102.002', name: 'Deposit Pasien Terapi Tumbuh Kembang Kidinest', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'KIDINEST', parentCode: '2102', currentBalance: 924000000 },
    { id: '2103', code: '2103', name: 'TITIPAN BIAYA', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '21', currentBalance: 207409000 },

    // 22 - HUTANG DAGANG - BCS (Rp 4.300.000,00 Cr)
    { id: '22', code: '22', name: 'HUTANG DAGANG - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 4300000 },
    { id: '2201', code: '2201', name: 'Hutang PBF Distributor Farmasi & Obat', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '22', currentBalance: 2800000 },
    { id: '2202', code: '2202', name: 'Hutang Supplier Kosmetika & Skincare SSC', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'SSC', parentCode: '22', currentBalance: 1500000 },

    // 23 - HUTANG BANK - BCS (Rp 6.245.750.462,00 Cr)
    { id: '23', code: '23', name: 'HUTANG BANK - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 6245750462 },

    // 24 - HUTANG PIHAK KETIGA - BCS (Rp 1.363.725.333,01 Cr)
    { id: '24', code: '24', name: 'HUTANG PIHAK KETIGA - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 1363725333.01 },

    // 25 - HUTANG PAJAK - BCS (Rp 39.249.649,00 Cr)
    // Directly integrated with Tax Engine!
    { id: '25', code: '25', name: 'HUTANG PAJAK - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 39249649 },
    { id: '2501', code: '2501', name: 'Hutang PPh 21 (Dokter & Karyawan)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 22374594, taxCategory: 'PPH_21' },
    { id: '2501.001', code: '2501.001', name: 'PPh 21 Dokter Tenaga Medis Ahli (PMK 168/2023)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2501', currentBalance: 14874594, taxCategory: 'PPH_21' },
    { id: '2501.002', code: '2501.002', name: 'PPh 21 Karyawan & Staf Tetap (TER PP 58)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 4, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2501', currentBalance: 7500000, taxCategory: 'PPH_21' },
    { id: '2502', code: '2502', name: 'Hutang PPh 23 (Jasa Medis, B3 & Maintenance Alkes)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 2400000, taxCategory: 'PPH_23' },
    { id: '2503', code: '2503', name: 'Hutang PPh Final Pasal 4(2) (Sewa Gedung Ruko)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 5000000, taxCategory: 'PPH_4_2' },
    { id: '2504', code: '2504', name: 'Hutang PPN Keluaran (Penjualan Kosmetika & Skincare)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 9475055, taxCategory: 'PPN_OUT' },
    { id: '2505', code: '2505', name: 'Hutang PPh 25 (Angsuran Pajak Penghasilan Badan)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 0, taxCategory: 'PPH_25' },
    { id: '2506', code: '2506', name: 'Hutang PPh 29 (SPT Tahunan PPh Badan 1771)', type: 'LIABILITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '25', currentBalance: 0, taxCategory: 'PPH_29' },

    // 26 - HUTANG PEMBEBANAN HPP - BCS (Rp 161.948.419,04 Cr)
    { id: '26', code: '26', name: 'HUTANG PEMBEBANAN HPP - BCS', type: 'LIABILITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '2', currentBalance: 161948419.04 },

    // ==========================================
    // 3 - EKUITAS - BCS (Rp 13.864.303.991,93 Cr)
    // ==========================================
    { id: '3', code: '3', name: 'EKUITAS - BCS', type: 'EQUITY', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 13864303991.93 },
    
    // 31 - MODAL - BCS (Rp 2.367.382.942,15 Cr)
    { id: '31', code: '31', name: 'MODAL - BCS', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '3', currentBalance: 2367382942.15 },
    { id: '3101', code: '3101', name: 'Modal Disetor Pemilik Klinik Kidinest & SSC', type: 'EQUITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '31', currentBalance: 2367382942.15 },

    // 32 - LABA AKUMULASI - BCS (Rp 11.496.921.049,78 Cr)
    { id: '32', code: '32', name: 'LABA AKUMULASI - BCS', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '3', currentBalance: 11496921049.78 },
    { id: '3201', code: '3201', name: 'Laba Ditahan (Retained Earnings)', type: 'EQUITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '32', currentBalance: 9818921049.78 },
    { id: '3202', code: '3202', name: 'Laba Periode Berjalan', type: 'EQUITY', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '32', currentBalance: 1678000000.00 },
    { id: '3999', code: '3999', name: 'Opening Balance Equity - BCS', type: 'EQUITY', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '3', currentBalance: 0 },

    // ==========================================
    // 4 - PENDAPATAN - BCS (Rp 15.860.537.440,00 Cr)
    // ==========================================
    { id: '4', code: '4', name: 'PENDAPATAN - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 15860537440 },
    { id: '41', code: '41', name: 'PENDAPATAN USAHA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '4', currentBalance: 15860537440 },
    { id: '4101', code: '4101', name: 'PENJUALAN RUMAH / TINDAKAN UTAMA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '41', currentBalance: 15860537440 },
    { id: '4102', code: '4102', name: 'PENJUALAN RUKO & KIOS - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '41', currentBalance: 0 },
    { id: '4103', code: '4103', name: 'PENJUALAN KAVLING / KELEBIHAN TANAH - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '41', currentBalance: 0 },
    { id: '4999', code: '4999', name: 'Opening Balance Income - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '4', currentBalance: 0 },

    // ==========================================
    // 5 - BEBAN - BCS (Rp 12.284.590.157,71 Dr)
    // ==========================================
    { id: '5', code: '5', name: 'BEBAN - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 12284590157.71 },
    
    // 51 - HARGA POKOK PENJUALAN - BCS (Rp 9.807.287.775,40 Dr)
    { id: '51', code: '51', name: 'HARGA POKOK PENJUALAN - BCS', type: 'COGS', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '5', currentBalance: 9807287775.40 },
    { id: '5102', code: '5102', name: 'HPP PEROLEHAN TANAH / BAHAN MEDIS UTAMA', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 4415886188 },
    { id: '5104', code: '5104', name: 'HPP LEGALITAS & PERIZINAN', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 226425353.98 },
    { id: '5105', code: '5105', name: 'HPP COR JALAN UTAMA / PRASARANA', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 34474242.20 },
    { id: '5107', code: '5107', name: 'HPP PONDASI TURAP', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 22110104.61 },
    { id: '5109', code: '5109', name: 'HPP SALURAN', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 119211301.60 },
    { id: '5110', code: '5110', name: 'HPP JARINGAN LISTRIK', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 115185189 },
    { id: '5111', code: '5111', name: 'HPP JARINGAN AIR', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 113318172.69 },
    { id: '5114', code: '5114', name: 'HPP INFRASTRUKTUR/KANSTIN', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 245434412.10 },
    { id: '5119', code: '5119', name: 'HPP RUMAH / TINDAKAN UNIT', type: 'COGS', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '51', currentBalance: 4515242811.22 },

    // 52 - BEBAN OPERASIONAL - BCS (Rp 2.477.302.382,31 Dr)
    { id: '52', code: '52', name: 'BEBAN OPERASIONAL - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '5', currentBalance: 2477302382.31 },
    { id: '5201', code: '5201', name: 'BIAYA MARKETING & PROMOSI', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '52', currentBalance: 443637768 },
    { id: '5202', code: '5202', name: 'BEBAN UMUM DAN ADMINISTRASI (Gaji, Listrik, Sewa, Limbah)', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '52', currentBalance: 2033664614.31 },
    { id: '5998', code: '5998', name: 'Opening Balance Expense - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '5', currentBalance: 0 },
    { id: '5999', code: '5999', name: 'Penyesuaian Stok - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '5', currentBalance: 0 },

    // ==========================================
    // 6 - PENDAPATAN NON USAHA - BCS (Rp 38.481.777,73 Cr)
    // ==========================================
    { id: '6', code: '6', name: 'PENDAPATAN NON USAHA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 38481777.73 },
    { id: '61', code: '61', name: 'PENDAPATAN LAIN-LAIN NON USAHA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '6', currentBalance: 38481777.73 },
    { id: '6101', code: '6101', name: 'PENDAPATAN BUNGA BANK - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 9545162.73 },
    { id: '6102', code: '6102', name: 'PENDAPATAN PEMBATALAN UNIT - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 5000000 },
    { id: '6103', code: '6103', name: 'PENDAPATAN PINDAH KAVLING - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 0 },
    { id: '6104', code: '6104', name: 'LABA PENJUALAN AKTIVA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 0 },
    { id: '6105', code: '6105', name: 'SELISIH PEMBULATAN - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 0 },
    { id: '6106', code: '6106', name: 'PENDAPATAN LAIN-LAIN NON USAHA - BCS', type: 'REVENUE', normalBalance: 'CREDIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '61', currentBalance: 23936615 },

    // ==========================================
    // 7 - BEBAN NON USAHA - BCS (Rp 872.980.221,16 Dr)
    // ==========================================
    { id: '7', code: '7', name: 'BEBAN NON USAHA - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 1, isHeader: true, isActive: true, businessUnit: 'ALL', currentBalance: 872980221.16 },
    { id: '71', code: '71', name: 'BIAYA LAIN-LAIN NON USAHA - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 2, isHeader: true, isActive: true, businessUnit: 'ALL', parentCode: '7', currentBalance: 872980221.16 },
    { id: '7101', code: '7101', name: 'BEBAN PROVISI KYG KPL - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 0 },
    { id: '7102', code: '7102', name: 'BUNGA KYG KPL - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 242190000 },
    { id: '7103', code: '7103', name: 'BEBAN BIAYA ADMINISTRASI BANK - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 47574598 },
    { id: '7104', code: '7104', name: 'PAJAK JASA GIRO - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 1777489.15, taxCategory: 'PPH_4_2' },
    { id: '7105', code: '7105', name: 'KERUGIAN PENJUALAN AKTIVA - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 0 },
    { id: '7106', code: '7106', name: 'SELISIH PEMBULATAN & PENYESUAIAN - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 256 },
    { id: '7107', code: '7107', name: 'BEBAN BUNGA PINJAMAN LAIN - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 581437878.01 },
    { id: '7108', code: '7108', name: 'AKUN PENGHAPUSAN - BCS', type: 'EXPENSE', normalBalance: 'DEBIT', level: 3, isHeader: false, isActive: true, businessUnit: 'ALL', parentCode: '71', currentBalance: 0 }
  ])

  // Accounting Periods
  const periods = ref<AccountingPeriod[]>([
    {
      id: 'p-2026-01',
      year: 2026,
      month: 1,
      periodName: 'Januari 2026',
      startDate: '2026-01-01',
      endDate: '2026-01-31',
      status: 'LOCKED',
      lockedAt: '2026-02-05 18:00:00',
      lockedBy: 'budi.santoso@kidinest-ssc.id',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    },
    {
      id: 'p-2026-02',
      year: 2026,
      month: 2,
      periodName: 'Februari 2026',
      startDate: '2026-02-01',
      endDate: '2026-02-28',
      status: 'LOCKED',
      lockedAt: '2026-03-05 17:30:00',
      lockedBy: 'budi.santoso@kidinest-ssc.id',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    },
    {
      id: 'p-2026-03',
      year: 2026,
      month: 3,
      periodName: 'Maret 2026 (Aktif)',
      startDate: '2026-03-01',
      endDate: '2026-03-31',
      status: 'OPEN',
      unpostedJournalCount: 3,
      pendingApprovalCount: 2
    },
    {
      id: 'p-2026-04',
      year: 2026,
      month: 4,
      periodName: 'April 2026',
      startDate: '2026-04-01',
      endDate: '2026-04-30',
      status: 'OPEN',
      unpostedJournalCount: 0,
      pendingApprovalCount: 0
    }
  ])

  // Compute General Ledger Entries directly from POSTED Journals
  const generalLedgerEntries = computed<LedgerEntry[]>(() => {
    const list: LedgerEntry[] = []
    const postedJournals = journalsStore.journals.filter(j => j.status === 'POSTED')

    for (const journal of postedJournals) {
      if (accountingStore.businessUnit !== 'ALL' && journal.businessUnit !== 'ALL' && journal.businessUnit !== accountingStore.businessUnit) {
        continue
      }

      for (const line of journal.lines) {
        list.push({
          id: `gle-${journal.id}-${line.id}`,
          date: journal.journalDate,
          journalNumber: journal.journalNumber,
          journalId: journal.id,
          source: journal.journalSource,
          description: line.description || journal.description,
          accountCode: line.accountCode,
          accountName: line.accountName,
          businessUnit: journal.businessUnit,
          branchName: journal.branchName,
          debit: line.debit,
          credit: line.credit,
          runningBalance: 0
        })
      }
    }

    list.sort((a, b) => a.date.localeCompare(b.date))
    return list
  })

  // Function to get account ledger with running balance
  function getAccountLedger(accountCode: string) {
    const account = accounts.value.find(a => a.code === accountCode)
    if (!account) return { account: null, entries: [], openingBalance: 0, endingBalance: 0 }

    const isDebitNormal = account.normalBalance === 'DEBIT'
    const openingBalance = account.currentBalance || 0
    let currentBalance = openingBalance

    const entries = generalLedgerEntries.value
      .filter(e => e.accountCode === accountCode)
      .map(entry => {
        if (isDebitNormal) {
          currentBalance += (entry.debit - entry.credit)
        } else {
          currentBalance += (entry.credit - entry.debit)
        }
        return {
          ...entry,
          runningBalance: currentBalance
        }
      })

    return {
      account,
      entries,
      openingBalance,
      endingBalance: currentBalance
    }
  }

  // Compute Trial Balance Rows dynamically
  const trialBalanceRows = computed<TrialBalanceRow[]>(() => {
    const nonHeaderAccounts = accounts.value.filter(a => !a.isHeader)
    const rows: TrialBalanceRow[] = []

    for (const acc of nonHeaderAccounts) {
      const glForAcc = generalLedgerEntries.value.filter(e => e.accountCode === acc.code)
      const periodDebit = glForAcc.reduce((sum, e) => sum + e.debit, 0)
      const periodCredit = glForAcc.reduce((sum, e) => sum + e.credit, 0)

      const openingDebit = acc.normalBalance === 'DEBIT' ? (acc.currentBalance || 0) : 0
      const openingCredit = acc.normalBalance === 'CREDIT' ? (acc.currentBalance || 0) : 0

      let endingDebit = 0
      let endingCredit = 0

      if (acc.normalBalance === 'DEBIT') {
        const net = openingDebit + periodDebit - periodCredit
        if (net >= 0) endingDebit = net
        else endingCredit = Math.abs(net)
      } else {
        const net = openingCredit + periodCredit - periodDebit
        if (net >= 0) endingCredit = net
        else endingDebit = Math.abs(net)
      }

      rows.push({
        accountCode: acc.code,
        accountName: acc.name,
        accountType: acc.type,
        openingDebit,
        openingCredit,
        periodDebit,
        periodCredit,
        endingDebit,
        endingCredit
      })
    }

    return rows
  })

  // Trial Balance totals
  const trialBalanceTotals = computed(() => {
    return trialBalanceRows.value.reduce(
      (acc, row) => ({
        openingDebit: acc.openingDebit + row.openingDebit,
        openingCredit: acc.openingCredit + row.openingCredit,
        periodDebit: acc.periodDebit + row.periodDebit,
        periodCredit: acc.periodCredit + row.periodCredit,
        endingDebit: acc.endingDebit + row.endingDebit,
        endingCredit: acc.endingCredit + row.endingCredit
      }),
      {
        openingDebit: 0,
        openingCredit: 0,
        periodDebit: 0,
        periodCredit: 0,
        endingDebit: 0,
        endingCredit: 0
      }
    )
  })

  const isLoading = ref(false)

  // Fetch Accounts directly from PostgreSQL database via Go API
  async function fetchAccounts() {
    try {
      isLoading.value = true
      const backendAccounts = await accountingApi.getAccounts({ limit: 250 })
      if (backendAccounts && backendAccounts.length > 0) {
        // Map id to code for parent lookup
        const idToCode = new Map<string, string>()
        backendAccounts.forEach(b => idToCode.set(b.id, b.code))

        // Cache existing calculated balances
        const balanceMap = new Map<string, number>()
        accounts.value.forEach(a => {
          if (a.currentBalance !== undefined) {
            balanceMap.set(a.code, a.currentBalance)
          }
        })

        accounts.value = backendAccounts.map(b => {
          const parentCode = b.parent_id ? idToCode.get(b.parent_id) : undefined
          let bu: 'ALL' | 'KIDINEST' | 'SSC' = 'ALL'
          if (b.name.includes('Kidinest') || b.code.endsWith('.001')) bu = 'KIDINEST'
          else if (b.name.includes('SSC') || b.name.includes('Skin Center') || b.code.endsWith('.002')) bu = 'SSC'

          return {
            id: b.id,
            code: b.code,
            name: b.name,
            type: b.type,
            normalBalance: b.normal_balance,
            level: b.level,
            isHeader: b.is_header,
            isActive: b.is_active,
            businessUnit: bu,
            parentCode,
            description: b.description,
            currentBalance: balanceMap.get(b.code) || 0
          }
        })
      }
    } catch (err: any) {
      console.error('Failed to fetch accounts from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load from database
  fetchAccounts()

  async function addAccount(newAcc: Partial<Account>) {
    try {
      isLoading.value = true
      let parentId: string | undefined = undefined
      if (newAcc.parentCode) {
        const parent = accounts.value.find(a => a.code === newAcc.parentCode)
        if (parent) parentId = parent.id
      }

      const createdBackend = await accountingApi.createAccount({
        code: newAcc.code || '',
        name: newAcc.name || '',
        type: newAcc.type || 'EXPENSE',
        normal_balance: newAcc.normalBalance || 'DEBIT',
        parent_id: parentId,
        level: newAcc.level || 2,
        is_header: newAcc.isHeader || false,
        description: newAcc.description
      })

      const created: Account = {
        id: createdBackend.id,
        code: createdBackend.code,
        name: createdBackend.name,
        type: createdBackend.type,
        normalBalance: createdBackend.normal_balance,
        level: createdBackend.level,
        isHeader: createdBackend.is_header,
        isActive: createdBackend.is_active,
        businessUnit: newAcc.businessUnit || 'ALL',
        description: createdBackend.description,
        parentCode: newAcc.parentCode,
        currentBalance: 0
      }
      accounts.value.push(created)
      accountingStore.showToast(`Akun ${created.code} - ${created.name} berhasil disimpan ke database!`, 'success')
      return created
    } catch (err: any) {
      accountingStore.showToast(`Gagal menyimpan akun ke database: ${err.message}`, 'error')
      throw err
    } finally {
      isLoading.value = false
    }
  }

  function lockPeriod(periodId: string) {
    const p = periods.value.find(item => item.id === periodId)
    if (p) {
      p.status = 'LOCKED'
      p.lockedAt = new Date().toISOString()
      p.lockedBy = accountingStore.currentUserEmail
      accountingStore.showToast(`Periode ${p.periodName} berhasil ditutup & dikunci!`, 'success')
    }
  }

  return {
    accounts,
    periods,
    isLoading,
    generalLedgerEntries,
    trialBalanceRows,
    trialBalanceTotals,
    fetchAccounts,
    getAccountLedger,
    addAccount,
    lockPeriod
  }
})
