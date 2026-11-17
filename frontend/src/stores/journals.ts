import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { JournalEntry, JournalStatus, BusinessUnit, JournalSource } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { accountingApi, type BackendJournal } from '@/api/accounting'

export const useJournalsStore = defineStore('journals', () => {
  const accountingStore = useAccountingStore()

  const journals = ref<JournalEntry[]>([
    {
      id: 'jv-101',
      journalNumber: 'JV-2026-03-001',
      journalDate: '2026-03-01',
      journalSource: 'CASHIER_CLOSING',
      sourceReference: 'CSH-KD-20260301-01',
      businessUnit: 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penutupan Kasir Shift 1 & 2 Kidinest - 1 Maret 2026',
      totalDebit: 18500000,
      totalCredit: 18500000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-101-1',
          accountId: 'acc-1-1111',
          accountCode: '1-1111',
          accountName: 'Kas Kasir Kidinest',
          debit: 12500000,
          credit: 0,
          description: 'Penerimaan Kas Kasir Kidinest'
        },
        {
          id: 'jl-101-2',
          accountId: 'acc-1-1113',
          accountCode: '1-1113',
          accountName: 'Bank BCA Operasional',
          debit: 6000000,
          credit: 0,
          description: 'Penerimaan EDC & QRIS BCA'
        },
        {
          id: 'jl-101-3',
          accountId: 'acc-4-1200',
          accountCode: '4-1200',
          accountName: 'Pendapatan Jasa Dokter Spesialis Anak (Sp.A)',
          debit: 0,
          credit: 14000000,
          description: 'Konsultasi Spesialis Anak & Tumbuh Kembang'
        },
        {
          id: 'jl-101-4',
          accountId: 'acc-4-1500',
          accountCode: '4-1500',
          accountName: 'Pendapatan Penjualan Obat & Racikan Anak',
          debit: 0,
          credit: 4500000,
          description: 'Penjualan Obat Resep Pasien Anak'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-01 21:00:00',
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: '2026-03-02 08:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-02 09:15:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-02 09:16:00',
      isManual: false
    },
    {
      id: 'jv-102',
      journalNumber: 'JV-2026-03-002',
      journalDate: '2026-03-01',
      journalSource: 'CASHIER_CLOSING',
      sourceReference: 'CSH-SSC-20260301-01',
      businessUnit: 'SSC',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penutupan Kasir Shift Harian Skin & Slimming Clinic (SSC)',
      totalDebit: 27800000,
      totalCredit: 27800000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-102-1',
          accountId: 'acc-1-1112',
          accountCode: '1-1112',
          accountName: 'Kas Kasir Semarang Skin Center (SSC)',
          debit: 7800000,
          credit: 0,
          description: 'Penerimaan Tunai Kasir SSC'
        },
        {
          id: 'jl-102-2',
          accountId: 'acc-1-1113',
          accountCode: '1-1113',
          accountName: 'Bank BCA Operasional',
          debit: 20000000,
          credit: 0,
          description: 'Penerimaan Gesek EDC Kartu Kredit & Debit'
        },
        {
          id: 'jl-102-3',
          accountId: 'acc-4-2300',
          accountCode: '4-2300',
          accountName: 'Pendapatan Tindakan Laser, Cauter & Peeling',
          debit: 0,
          credit: 18000000,
          description: 'Tindakan Laser & Aesthetic Dokter Sp.KK'
        },
        {
          id: 'jl-102-4',
          accountId: 'acc-4-2400',
          accountCode: '4-2400',
          accountName: 'Pendapatan Penjualan Produk Skincare & Sunscreen',
          debit: 0,
          credit: 9800000,
          description: 'Penjualan Krim & Suplemen'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-01 21:30:00',
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: '2026-03-02 08:35:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-02 09:18:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-02 09:18:30',
      isManual: false
    },
    {
      id: 'jv-103',
      journalNumber: 'JV-2026-03-003',
      journalDate: '2026-03-02',
      journalSource: 'PHARMACY_PURCHASE',
      sourceReference: 'PO-KF-20260301-44',
      businessUnit: 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pembelian Obat & Vaksin Anak dari PT Kimia Farma Trading',
      totalDebit: 15540000,
      totalCredit: 15540000,
      status: 'APPROVED',
      lines: [
        {
          id: 'jl-103-1',
          accountId: 'acc-1400',
          accountCode: '1-1400',
          accountName: 'Persediaan Farmasi & Obat',
          debit: 14000000,
          credit: 0,
          description: 'Restock Vaksin Hexaxim & Rotarix'
        },
        {
          id: 'jl-103-2',
          accountId: 'acc-1180',
          accountCode: '1-1180',
          accountName: 'PPN Masukan (Dapat Dikreditkan)',
          debit: 1540000,
          credit: 0,
          description: 'PPN 11% Faktur Pajak PT Kimia Farma'
        },
        {
          id: 'jl-103-3',
          accountId: 'acc-2100',
          accountCode: '2-1100',
          accountName: 'Hutang Usaha Distributor Farmasi',
          debit: 0,
          credit: 15540000,
          description: 'Hutang Dagang Termin 30 Hari'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-02 10:15:00',
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: '2026-03-02 11:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-02 11:30:00',
      isManual: false
    },
    {
      id: 'jv-104',
      journalNumber: 'JV-2026-03-004',
      journalDate: '2026-03-02',
      journalSource: 'DOCTOR_FEE',
      sourceReference: 'DF-20260302-09',
      businessUnit: 'SSC',
      branchId: 'BRANCH-JAKSEL',
      branchName: 'Pondok Indah',
      description: 'Perhitungan Jasa Medis dr. Maya Sandria, Sp.KK + PPh 21 Dokter',
      totalDebit: 12000000,
      totalCredit: 12000000,
      status: 'REVIEWED',
      lines: [
        {
          id: 'jl-104-1',
          accountId: 'acc-5200',
          accountCode: '5-1200',
          accountName: 'Beban Jasa Medis Dokter',
          debit: 12000000,
          credit: 0,
          description: 'Share Dokter 60% dari Tindakan Aesthetic'
        },
        {
          id: 'jl-104-2',
          accountId: 'acc-2130',
          accountCode: '2-1300',
          accountName: 'Hutang PPh 21 Dokter',
          debit: 0,
          credit: 300000,
          description: 'Potongan PPh 21 (5% x 50% x Rp 12.000.000)'
        },
        {
          id: 'jl-104-3',
          accountId: 'acc-2120',
          accountCode: '2-1200',
          accountName: 'Hutang Jasa Medis Dokter',
          debit: 0,
          credit: 11700000,
          description: 'Net Honor Dokter Terutang'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-02 14:00:00',
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: '2026-03-02 15:20:00',
      isManual: false
    },
    {
      id: 'jv-105',
      journalNumber: 'JV-2026-03-005',
      journalDate: '2026-03-03',
      journalSource: 'FIXED_ASSET',
      sourceReference: 'DEP-202603-01',
      businessUnit: 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penyusutan Bulanan Alat Medis USG 4D Pediatric',
      totalDebit: 4166667,
      totalCredit: 4166667,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-105-1',
          accountId: 'acc-6400',
          accountCode: '6-1400',
          accountName: 'Beban Penyusutan Aset Tetap',
          debit: 4166667,
          credit: 0,
          description: 'Depresiasi Garis Lurus USG 4D (5 Tahun)'
        },
        {
          id: 'jl-105-2',
          accountId: 'acc-1210',
          accountCode: '1-2100',
          accountName: 'Akumulasi Penyusutan Alat Medis',
          debit: 0,
          credit: 4166667,
          description: 'Akumulasi Depresiasi USG 4D'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-03 09:00:00',
      reviewedBy: 'dewi.lestari',
      reviewedByName: 'Dewi Lestari, SE, Ak.',
      reviewedAt: '2026-03-03 09:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-03 10:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-03 10:01:00',
      isManual: false
    },
    {
      id: 'jv-106',
      journalNumber: 'JV-2026-03-006',
      journalDate: '2026-03-03',
      journalSource: 'INVENTORY_CONSUMPTION',
      sourceReference: 'INV-USE-20260303-12',
      businessUnit: 'SSC',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pemakaian BHP Medis Jarum & Kassa Ruang Tindakan Estetika',
      totalDebit: 1850000,
      totalCredit: 1850000,
      status: 'SUBMITTED',
      lines: [
        {
          id: 'jl-106-1',
          accountId: 'acc-6500',
          accountCode: '6-1500',
          accountName: 'Beban Perlengkapan & BHP Medis',
          debit: 1850000,
          credit: 0,
          description: 'Pemakaian BHP Tindakan Pasien SSC'
        },
        {
          id: 'jl-106-2',
          accountId: 'acc-1400',
          accountCode: '1-1400',
          accountName: 'Persediaan Farmasi & Obat',
          debit: 0,
          credit: 1850000,
          description: 'Pengurangan Stok BHP Medis'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-03 16:45:00',
      isManual: false
    },
    {
      id: 'jv-107',
      journalNumber: 'JV-2026-03-007',
      journalDate: '2026-03-04',
      journalSource: 'MANUAL',
      sourceReference: 'MAN-ADJ-202603-01',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penyesuaian Biaya Administrasi Bank Mandiri & Bunga Tabungan',
      totalDebit: 75000,
      totalCredit: 75000,
      status: 'DRAFT',
      lines: [
        {
          id: 'jl-107-1',
          accountId: 'acc-6900',
          accountCode: '6-1900',
          accountName: 'Beban Administrasi Bank',
          debit: 75000,
          credit: 0,
          description: 'Admin Bulanan Mandiri'
        },
        {
          id: 'jl-107-2',
          accountId: 'acc-1200',
          accountCode: '1-1200',
          accountName: 'Bank Mandiri Operasional',
          debit: 0,
          credit: 75000,
          description: 'Potongan Langsung Buku Rekening'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-04 10:10:00',
      isManual: true
    },
    {
      id: 'jv-108',
      journalNumber: 'JV-2026-03-008',
      journalDate: '2026-03-10',
      journalSource: 'MANUAL',
      sourceReference: 'INV-WST-202603',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pembayaran Jasa Pemusnahan Limbah Medis B3 PT Wastec International + Potongan PPh 23',
      totalDebit: 45000000,
      totalCredit: 45000000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-108-1',
          accountId: 'acc-5202',
          accountCode: '5202',
          accountName: 'Beban Pengolahan Limbah Medis B3',
          debit: 45000000,
          credit: 0,
          description: 'Invoice Limbah B3 PT Wastec Periode Maret'
        },
        {
          id: 'jl-108-2',
          accountId: 'acc-2502',
          accountCode: '2502',
          accountName: 'PPh 23 Jasa Teknik & Limbah B3',
          debit: 0,
          credit: 900000,
          description: 'Withholding Tax PPh 23 (Tarif 2% x Rp 45.000.000)'
        },
        {
          id: 'jl-108-3',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 44100000,
          description: 'Transfer Pembayaran via Bank BCA Netto'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-10 11:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-10 13:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-10 13:05:00',
      isManual: false
    },
    {
      id: 'jv-109',
      journalNumber: 'JV-2026-03-009',
      journalDate: '2026-03-12',
      journalSource: 'MANUAL',
      sourceReference: 'SEWA-RUKO-202603',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pembayaran Sewa Gedung Klinik Hj. Siti Nurhaliza + Pemotongan PPh Final Pasal 4(2)',
      totalDebit: 50000000,
      totalCredit: 50000000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-109-1',
          accountId: 'acc-5201',
          accountCode: '5201',
          accountName: 'Beban Sewa Gedung & Ruang Medis',
          debit: 50000000,
          credit: 0,
          description: 'Sewa Gedung Klinik Pratama Periode Maret'
        },
        {
          id: 'jl-109-2',
          accountId: 'acc-2503',
          accountCode: '2503',
          accountName: 'PPh Final Pasal 4(2) Persewaan Bangunan',
          debit: 0,
          credit: 5000000,
          description: 'Potongan PPh Final 4(2) 10% x Rp 50.000.000'
        },
        {
          id: 'jl-109-3',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 45000000,
          description: 'Transfer Pembayaran via Bank BCA Netto'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-12 10:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-12 11:15:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-12 11:20:00',
      isManual: false
    },
    {
      id: 'jv-110',
      journalNumber: 'JV-2026-03-010',
      journalDate: '2026-03-15',
      journalSource: 'MANUAL',
      sourceReference: 'NTPN-BCA-PPH25-03',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Setoran Angsuran Pajak Badan PPh 25 Masa Maret 2026 via Bank BCA',
      totalDebit: 13000000,
      totalCredit: 13000000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-110-1',
          accountId: 'acc-1504',
          accountCode: '1504',
          accountName: 'Pajak Dibayar di Muka PPh 25 (Kredit Pajak Badan)',
          debit: 13000000,
          credit: 0,
          description: 'Angsuran PPh 25 Bulanan Klinik Masa Maret'
        },
        {
          id: 'jl-110-2',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 13000000,
          description: 'Debet Rekening BCA NTPN-BCA-20260315-9921'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-15 09:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-15 10:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-15 10:05:00',
      isManual: false
    },
    {
      id: 'jv-111',
      journalNumber: 'JV-2026-03-011',
      journalDate: '2026-03-20',
      journalSource: 'MANUAL',
      sourceReference: 'NTPN-MDR-PPH21-03',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penyetoran Pajak PPh 21 Dokter & Karyawan Masa Maret via Bank Mandiri',
      totalDebit: 22362500,
      totalCredit: 22362500,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-111-1',
          accountId: 'acc-2501-1',
          accountCode: '2501.001',
          accountName: 'PPh 21 Dokter Tenaga Ahli (PMK 168/2023)',
          debit: 14862500,
          credit: 0,
          description: 'Penyetoran SPT Masa PPh 21 Dokter Masa Maret'
        },
        {
          id: 'jl-111-2',
          accountId: 'acc-2501-2',
          accountCode: '2501.002',
          accountName: 'PPh 21 Karyawan & Staf Tetap (TER PP 58)',
          debit: 7500000,
          credit: 0,
          description: 'Penyetoran SPT Masa PPh 21 Pegawai Tetap Masa Maret'
        },
        {
          id: 'jl-111-3',
          accountId: 'acc-1102-2',
          accountCode: '1102.002',
          accountName: 'Bank Mandiri Payroll & Pajak',
          debit: 0,
          credit: 22362500,
          description: 'Debet Payroll Mandiri NTPN-MDR-PPH21-03'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-03-20 14:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-03-20 14:30:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-03-20 14:35:00',
      isManual: false
    },
    {
      id: 'jv-202611-01',
      journalNumber: 'JV-2026-11-001',
      journalDate: '2026-11-15',
      journalSource: 'MANUAL',
      sourceReference: 'BUPOT-21-202611-001',
      businessUnit: 'KIDINEST',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pembayaran Honor & Pemotongan PPh 21 Dokter Spesialis Anak (dr. Hendra Wijaya, Sp.A)',
      totalDebit: 48500000,
      totalCredit: 48500000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-1',
          accountId: 'acc-5106',
          accountCode: '5106',
          accountName: 'Beban Bagi Hasil Dokter Spesialis Anak (Kidinest)',
          debit: 48500000,
          credit: 0,
          description: 'Honor Jasa Medis Konsultasi & Tindakan Dokter Sp.A'
        },
        {
          id: 'jl-2611-2',
          accountId: 'acc-2501-1',
          accountCode: '2501.001',
          accountName: 'PPh 21 Dokter Tenaga Ahli (PMK 168/2023)',
          debit: 0,
          credit: 1212500,
          description: 'Withholding Tax PPh 21 Dokter 5% atas DPP 50%'
        },
        {
          id: 'jl-2611-3',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 47287500,
          description: 'Transfer Honor Bersih via Bank BCA TRF-BCA-DOC-202611-01'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-15 11:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-15 13:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-15 13:05:00',
      isManual: false
    },
    {
      id: 'jv-202611-02',
      journalNumber: 'JV-2026-11-002',
      journalDate: '2026-11-18',
      journalSource: 'MANUAL',
      sourceReference: 'BUPOT-21-202611-002',
      businessUnit: 'SSC',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pembayaran Honor Medis & Pemotongan PPh 21 Dokter Spesialis Kulit (dr. Maya Sandria, Sp.KK)',
      totalDebit: 82500000,
      totalCredit: 82500000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-4',
          accountId: 'acc-5107',
          accountCode: '5107',
          accountName: 'Beban Bagi Hasil Dokter Estetika & Kulit SSC',
          debit: 82500000,
          credit: 0,
          description: 'Jasa Medis Tindakan Laser & Aesthetic Dokter Kulit'
        },
        {
          id: 'jl-2611-5',
          accountId: 'acc-2501-1',
          accountCode: '2501.001',
          accountName: 'PPh 21 Dokter Tenaga Ahli (PMK 168/2023)',
          debit: 0,
          credit: 2062500,
          description: 'Pemotongan PPh 21 Tenaga Ahli Dokter Sp.KK'
        },
        {
          id: 'jl-2611-6',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 80437500,
          description: 'Transfer Bank BCA TRF-BCA-DOC-202611-02'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-18 10:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-18 11:30:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-18 11:35:00',
      isManual: false
    },
    {
      id: 'jv-202611-03',
      journalNumber: 'JV-2026-11-003',
      journalDate: '2026-11-20',
      journalSource: 'MANUAL',
      sourceReference: 'INV-WST-202611',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Pengolahan & Pemusnahan Limbah Medis B3 PT Wastec International - PPh 23 (2%)',
      totalDebit: 47500000,
      totalCredit: 47500000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-7',
          accountId: 'acc-5202',
          accountCode: '5202',
          accountName: 'Beban Umum & Operasional Limbah Medis B3',
          debit: 47500000,
          credit: 0,
          description: 'Biaya Pengolahan & Pemusnahan Limbah B3 Klinis'
        },
        {
          id: 'jl-2611-8',
          accountId: 'acc-2502',
          accountCode: '2502',
          accountName: 'PPh 23 (Jasa Limbah B3 & Maintenance Alkes)',
          debit: 0,
          credit: 950000,
          description: 'Withholding Tax PPh 23 (2%) atas Jasa Teknik & Limbah'
        },
        {
          id: 'jl-2611-9',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 46550000,
          description: 'Transfer Pembayaran Vendor Limbah Medis BCA'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-20 09:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-20 10:15:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-20 10:20:00',
      isManual: false
    },
    {
      id: 'jv-202611-04',
      journalNumber: 'JV-2026-11-004',
      journalDate: '2026-11-22',
      journalSource: 'MANUAL',
      sourceReference: 'SEWA-RUKO-202611',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Sewa Gedung Fasilitas Operasional Klinik - PPh Final Pasal 4(2) (10%)',
      totalDebit: 17500000,
      totalCredit: 17500000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-10',
          accountId: 'acc-5202',
          accountCode: '5202',
          accountName: 'Beban Umum & Sewa Gedung Klinik',
          debit: 17500000,
          credit: 0,
          description: 'Sewa Fasilitas Ruang Rawat & Operasional'
        },
        {
          id: 'jl-2611-11',
          accountId: 'acc-2503',
          accountCode: '2503',
          accountName: 'PPh Final Pasal 4(2) (Sewa Ruko & Fasilitas)',
          debit: 0,
          credit: 1750000,
          description: 'Withholding Tax PPh Final 4(2) Sewa Gedung 10%'
        },
        {
          id: 'jl-2611-12',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 15750000,
          description: 'Transfer Pembayaran Sewa via Bank BCA'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-22 14:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-22 15:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-22 15:05:00',
      isManual: false
    },
    {
      id: 'jv-202611-05',
      journalNumber: 'JV-2026-11-005',
      journalDate: '2026-11-25',
      journalSource: 'MANUAL',
      sourceReference: 'NTPN-BCA-PPH25-11',
      businessUnit: 'ALL',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Setoran Angsuran Pajak Badan PPh Pasal 25 Masa November 2026 via Bank BCA',
      totalDebit: 27000000,
      totalCredit: 27000000,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-13',
          accountId: 'acc-1504',
          accountCode: '1504',
          accountName: 'Pajak Dibayar di Muka PPh 25 (Kredit Pajak Badan)',
          debit: 27000000,
          credit: 0,
          description: 'Setoran Angsuran PPh 25 Badan Bulanan via Bank BCA NTPN-BCA-20261125-9925'
        },
        {
          id: 'jl-2611-14',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 0,
          credit: 27000000,
          description: 'Debet Giro BCA e-Billing DJP'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-25 10:00:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-25 11:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-25 11:05:00',
      isManual: false
    },
    {
      id: 'jv-202611-06',
      journalNumber: 'JV-2026-11-006',
      journalDate: '2026-11-28',
      journalSource: 'MANUAL',
      sourceReference: 'FP-PPN-202611',
      businessUnit: 'SSC',
      branchId: 'BRANCH-MAIN',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Faktur Pajak PPN Keluaran 11% Penjualan Produk Skincare & Treatment Non-Medis SSC Masa November',
      totalDebit: 95611919,
      totalCredit: 95611919,
      status: 'POSTED',
      lines: [
        {
          id: 'jl-2611-15',
          accountId: 'acc-1102-1',
          accountCode: '1102.001',
          accountName: 'Bank BCA Operasional (Giro & EDC)',
          debit: 95611919,
          credit: 0,
          description: 'Penerimaan Penjualan Skincare + PPN EDC BCA'
        },
        {
          id: 'jl-2611-16',
          accountId: 'acc-4109',
          accountCode: '4109',
          accountName: 'Pendapatan Penjualan Produk Skincare & Kosmetika SSC',
          debit: 0,
          credit: 86136864,
          description: 'DPP Penjualan Produk Skincare Terkena PPN'
        },
        {
          id: 'jl-2611-17',
          accountId: 'acc-2504',
          accountCode: '2504',
          accountName: 'PPN Keluaran (Penjualan Kosmetika & Skincare)',
          debit: 0,
          credit: 9475055,
          description: 'PPN Keluaran 11% Terutang Masa November 2026'
        }
      ],
      createdBy: 'rian.pratama',
      createdByName: 'Rian Pratama, S.Ak.',
      createdAt: '2026-11-28 16:30:00',
      approvedBy: 'budi.santoso',
      approvedByName: 'Budi Santoso, SE, Ak., CA',
      approvedAt: '2026-11-28 17:00:00',
      postedBy: 'budi.santoso',
      postedByName: 'Budi Santoso, SE, Ak., CA',
      postedAt: '2026-11-28 17:05:00',
      isManual: false
    }
  ])

  // Filtered journals based on global filters (Business Unit & Fiscal Month)
  const filteredJournals = computed(() => {
    return journals.value.filter(j => {
      if (accountingStore.businessUnit !== 'ALL' && j.businessUnit !== 'ALL' && j.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      if (accountingStore.fiscalMonth && j.journalDate) {
        const str = String(j.journalDate).trim()
        const isoMatch = str.match(/^(\d{4})[-/](\d{1,2})[-/]/)
        const dmyMatch = str.match(/^\d{1,2}[-/](\d{1,2})[-/]\d{4}/)
        const m = isoMatch ? parseInt(isoMatch[2], 10) : (dmyMatch ? parseInt(dmyMatch[1], 10) : null)
        if (m !== null && m !== accountingStore.fiscalMonth) return false
      }
      return true
    })
  })

  // Statistics
  const stats = computed(() => {
    const list = filteredJournals.value
    return {
      totalCount: list.length,
      draftCount: list.filter(j => j.status === 'DRAFT').length,
      submittedCount: list.filter(j => j.status === 'SUBMITTED').length,
      reviewedCount: list.filter(j => j.status === 'REVIEWED').length,
      approvedCount: list.filter(j => j.status === 'APPROVED').length,
      postedCount: list.filter(j => j.status === 'POSTED').length,
      reversedCount: list.filter(j => j.status === 'REVERSED').length,
      totalPostedDebit: list
        .filter(j => j.status === 'POSTED')
        .reduce((sum, j) => sum + j.totalDebit, 0)
    }
  })

  // Actions
  function getJournalById(id: string): JournalEntry | undefined {
    return journals.value.find(j => j.id === id)
  }

  function submitDraft(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'DRAFT') {
      j.status = 'SUBMITTED'
      accountingStore.showToast(`Jurnal ${j.journalNumber} berhasil diajukan untuk review`, 'success')
    }
  }

  function reviewJournal(id: string) {
    const j = getJournalById(id)
    if (j && (j.status === 'SUBMITTED' || j.status === 'DRAFT')) {
      j.status = 'REVIEWED'
      j.reviewedBy = accountingStore.currentUserEmail
      j.reviewedByName = accountingStore.currentUserName
      j.reviewedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} telah direview`, 'success')
    }
  }

  function approveJournal(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'REVIEWED') {
      j.status = 'APPROVED'
      j.approvedBy = accountingStore.currentUserEmail
      j.approvedByName = accountingStore.currentUserName
      j.approvedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} telah disetujui (Approved)`, 'success')
    }
  }

  function postJournal(id: string) {
    const j = getJournalById(id)
    if (j && j.status === 'APPROVED') {
      j.status = 'POSTED'
      j.postedBy = accountingStore.currentUserEmail
      j.postedByName = accountingStore.currentUserName
      j.postedAt = new Date().toISOString()
      accountingStore.showToast(`Jurnal ${j.journalNumber} resmi diposting ke Buku Besar!`, 'success')
    }
  }

  function rejectJournal(id: string, reason: string) {
    const j = getJournalById(id)
    if (j) {
      j.status = 'REJECTED'
      j.reversalReason = reason
      accountingStore.showToast(`Jurnal ${j.journalNumber} ditolak. Alasan: ${reason}`, 'info')
    }
  }

  function reverseJournal(id: string, reason: string) {
    const original = getJournalById(id)
    if (!original || original.status !== 'POSTED') {
      accountingStore.showToast('Hanya jurnal yang sudah berstatus POSTED yang dapat dibalik', 'error')
      return
    }

    const revNumber = `REV-${original.journalNumber}`
    const revJournal: JournalEntry = {
      id: 'jv-rev-' + Date.now(),
      journalNumber: revNumber,
      journalDate: new Date().toISOString().split('T')[0],
      journalSource: original.journalSource,
      sourceReference: `Reversal of ${original.journalNumber}`,
      businessUnit: original.businessUnit,
      branchId: original.branchId,
      branchName: original.branchName,
      description: `[JURNAL PEMBALIK] ${original.description} - Alasan: ${reason}`,
      totalDebit: original.totalCredit,
      totalCredit: original.totalDebit,
      status: 'POSTED',
      lines: original.lines.map((l, idx) => ({
        id: `jl-rev-${idx}-${Date.now()}`,
        accountId: l.accountId,
        accountCode: l.accountCode,
        accountName: l.accountName,
        // Swap debit and credit for reversal
        debit: l.credit,
        credit: l.debit,
        description: `Pembalikan: ${l.description}`
      })),
      createdBy: accountingStore.currentUserEmail,
      createdByName: accountingStore.currentUserName,
      createdAt: new Date().toISOString(),
      postedBy: accountingStore.currentUserEmail,
      postedByName: accountingStore.currentUserName,
      postedAt: new Date().toISOString(),
      reversalReason: reason,
      isManual: true
    }

    original.status = 'REVERSED'
    original.reversalJournalId = revJournal.id
    original.reversalReason = reason

    journals.value.unshift(revJournal)
    accountingStore.showToast(`Jurnal pembalik ${revNumber} berhasil diposting`, 'success')
  }

  function createManualJournal(data: {
    journalDate: string
    businessUnit: BusinessUnit
    branchId?: string
    branchName?: string
    description: string
    journalSource?: JournalSource
    sourceReference?: string
    status?: JournalStatus
    lines: {
      accountId?: string
      accountCode: string
      accountName: string
      debit: number
      credit: number
      description?: string
    }[]
  }): JournalEntry {
    const totalDebit = data.lines.reduce((s, l) => s + (l.debit || 0), 0)
    const totalCredit = data.lines.reduce((s, l) => s + (l.credit || 0), 0)

    if (Math.abs(totalDebit - totalCredit) > 0.01) {
      throw new Error(`Total Debit (Rp ${totalDebit}) tidak seimbang dengan Total Kredit (Rp ${totalCredit})`)
    }

    const nextSeq = journals.value.length + 1
    const seqStr = String(nextSeq).padStart(3, '0')
    const jNumber = `JV-2026-03-${seqStr}`
    const st = data.status || 'DRAFT'

    const newJournal: JournalEntry = {
      id: 'jv-' + Date.now(),
      journalNumber: jNumber,
      journalDate: data.journalDate,
      journalSource: data.journalSource || 'MANUAL',
      sourceReference: data.sourceReference || `MAN-${Date.now().toString().slice(-6)}`,
      businessUnit: data.businessUnit,
      branchId: data.branchId || 'BRANCH-MAIN',
      branchName: data.branchName || 'Kantor Pusat Semarang',
      description: data.description,
      totalDebit,
      totalCredit,
      status: st,
      lines: data.lines.map((l, idx) => ({
        id: `jl-${idx}-${Date.now()}`,
        accountId: l.accountId || `acc-${l.accountCode}`,
        accountCode: l.accountCode,
        accountName: l.accountName,
        debit: l.debit || 0,
        credit: l.credit || 0,
        description: l.description || data.description
      })),
      createdBy: accountingStore.currentUserEmail,
      createdByName: accountingStore.currentUserName,
      createdAt: new Date().toISOString(),
      postedBy: st === 'POSTED' ? accountingStore.currentUserEmail : undefined,
      postedByName: st === 'POSTED' ? accountingStore.currentUserName : undefined,
      postedAt: st === 'POSTED' ? new Date().toISOString() : undefined,
      isManual: true
    }

    journals.value.unshift(newJournal)
    accountingStore.showToast(`Jurnal ${jNumber} berhasil dibuat (${st})!`, 'success')

    // Persist to PostgreSQL database asynchronously
    accountingApi.createJournal({
      journal_type: 'GENERAL',
      transaction_date: data.journalDate,
      reference_type: data.journalSource || 'MANUAL',
      reference_id: data.sourceReference || jNumber,
      description: data.description,
      auto_post: st === 'POSTED',
      lines: data.lines.map(l => ({
        account_code: l.accountCode,
        debit: l.debit || 0,
        credit: l.credit || 0,
        description: l.description || data.description
      }))
    }).catch(err => {
      console.warn('Note: Journal created in memory, DB sync response:', err.message)
    })

    return newJournal
  }

  const isLoading = ref(false)

  // Fetch journals directly from PostgreSQL database via Go API
  async function fetchJournals() {
    try {
      isLoading.value = true
      const backendJournals = await accountingApi.getJournals({ limit: 100 })
      if (backendJournals && backendJournals.length > 0) {
        journals.value = backendJournals.map(j => {
          let bu: BusinessUnit = 'ALL'
          const descUpper = (j.description || '').toUpperCase()
          if (descUpper.includes('KIDINEST')) bu = 'KIDINEST'
          else if (descUpper.includes('SSC') || descUpper.includes('SEMARANG SKIN CENTER')) bu = 'SSC'

          return {
            id: j.id,
            journalNumber: j.journal_number,
            journalDate: j.transaction_date,
            journalSource: (j.reference_type as JournalSource) || 'MANUAL',
            sourceReference: j.reference_id || j.journal_number,
            businessUnit: bu,
            branchId: j.branch_id || 'BRANCH-MAIN',
            branchName: 'Kantor Pusat Semarang',
            description: j.description,
            totalDebit: Number(j.total_debit),
            totalCredit: Number(j.total_credit),
            status: (j.status as JournalStatus) || 'POSTED',
            lines: (j.lines || []).map((l, idx) => ({
              id: l.id || `jl-${idx}`,
              accountId: l.account_id,
              accountCode: l.account_code,
              accountName: l.account_name,
              debit: Number(l.debit),
              credit: Number(l.credit),
              description: l.description || j.description
            })),
            createdBy: accountingStore.currentUserEmail,
            createdByName: accountingStore.currentUserName,
            createdAt: j.created_at || new Date().toISOString(),
            isManual: true
          }
        })
      }
    } catch (err: any) {
      console.error('Failed to fetch journals from database:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Initial load
  fetchJournals()

  return {
    journals,
    isLoading,
    filteredJournals,
    stats,
    fetchJournals,
    getJournalById,
    submitDraft,
    reviewJournal,
    approveJournal,
    postJournal,
    rejectJournal,
    reverseJournal,
    createManualJournal
  }
})
