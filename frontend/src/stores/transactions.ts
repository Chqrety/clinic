import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  OperationalTransaction,
  SourceModule,
  AccountingMappingRule,
  BusinessUnit
} from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { useJournalsStore } from './journals'

export const useTransactionsStore = defineStore('transactions', () => {
  const accountingStore = useAccountingStore()
  const journalsStore = useJournalsStore()

  // Pre-configured Accounting Mapping Rules
  const mappingRules = ref<AccountingMappingRule[]>([
    {
      id: 'rule-1',
      eventKey: 'CASHIER_CLOSING_KIDINEST',
      eventName: 'Penutupan Kasir Harian Kidinest',
      businessUnit: 'KIDINEST',
      sourceModule: 'GRO_CASHIER',
      debitAccountCode: '1-1100',
      debitAccountName: 'Kas Operasional Kidinest',
      creditAccountCode: '4-1100',
      creditAccountName: 'Pendapatan Jasa Medis Kidinest',
      description: 'Mencatat penerimaan kasir shift harian klinik anak',
      isActive: true
    },
    {
      id: 'rule-2',
      eventKey: 'CASHIER_CLOSING_SSC',
      eventName: 'Penutupan Kasir Harian SSC',
      businessUnit: 'SSC',
      sourceModule: 'GRO_CASHIER',
      debitAccountCode: '1-1101',
      debitAccountName: 'Kas Operasional SSC',
      creditAccountCode: '4-1101',
      creditAccountName: 'Pendapatan Jasa Medis SSC',
      description: 'Mencatat penerimaan kasir shift harian klinik estetika',
      isActive: true
    },
    {
      id: 'rule-3',
      eventKey: 'PHARMACY_SALE',
      eventName: 'Penjualan Obat Resep Farmasi',
      businessUnit: 'ALL',
      sourceModule: 'PHARMACY',
      debitAccountCode: '1-1100',
      debitAccountName: 'Kas Operasional Kidinest',
      creditAccountCode: '4-1200',
      creditAccountName: 'Pendapatan Farmasi & Obat',
      cogsAccountCode: '5-1100',
      inventoryAccountCode: '1-1400',
      description: 'Pengakuan omset obat beserta beban pokok (HPP)',
      isActive: true
    },
    {
      id: 'rule-4',
      eventKey: 'DOCTOR_FEE_CALC',
      eventName: 'Kompensasi Jasa Medis Dokter + PPh 21',
      businessUnit: 'ALL',
      sourceModule: 'DOCTOR_SHARE',
      debitAccountCode: '5-1200',
      debitAccountName: 'Beban Jasa Medis Dokter',
      creditAccountCode: '2-1200',
      creditAccountName: 'Hutang Jasa Medis Dokter',
      taxAccountCode: '2-1300',
      description: 'Kewajiban honor dokter dikurangi potongan PPh 21 dokter',
      isActive: true
    },
    {
      id: 'rule-5',
      eventKey: 'PURCHASE_PHARMACY_INVOICE',
      eventName: 'Faktur Pembelian Obat dari Distributor',
      businessUnit: 'ALL',
      sourceModule: 'PURCHASING',
      debitAccountCode: '1-1400',
      debitAccountName: 'Persediaan Farmasi & Obat',
      creditAccountCode: '2-1100',
      creditAccountName: 'Hutang Usaha Distributor Farmasi',
      taxAccountCode: '1-1180',
      description: 'Penerimaan stok obat berjangka tempo dengan PPN Masukan',
      isActive: true
    },
    {
      id: 'rule-6',
      eventKey: 'FIXED_ASSET_DEPRECIATION',
      eventName: 'Penyusutan Bulanan Peralatan Medis',
      businessUnit: 'ALL',
      sourceModule: 'FIXED_ASSET',
      debitAccountCode: '6-1400',
      debitAccountName: 'Beban Penyusutan Aset Tetap',
      creditAccountCode: '1-2100',
      creditAccountName: 'Akumulasi Penyusutan Alat Medis',
      description: 'Beban amortisasi bulanan peralatan laser dan USG',
      isActive: true
    },
    {
      id: 'rule-7',
      eventKey: 'PAYROLL_SALARY',
      eventName: 'Penggajian Karyawan & Tenaga Medis',
      businessUnit: 'ALL',
      sourceModule: 'PAYROLL',
      debitAccountCode: '6-1100',
      debitAccountName: 'Gaji & Upah Karyawan',
      creditAccountCode: '1-1200',
      creditAccountName: 'Bank Mandiri Operasional',
      taxAccountCode: '2-1300',
      description: 'Pembayaran payroll bulanan staf klinik',
      isActive: true
    }
  ])

  // Mock Operational Transactions from External Systems
  const transactions = ref<OperationalTransaction[]>([
    {
      id: 'tx-001',
      transactionNumber: 'TX-CSH-20260304-001',
      transactionDate: '2026-03-04 14:00:00',
      sourceModule: 'GRO_CASHIER',
      referenceNumber: 'SHIFT-KID-01',
      businessUnit: 'KIDINEST',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Closing Kasir Shift 1 Kidinest - 12 Pasien Poli Spesialis Anak',
      amount: 9450000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '1-1100 (Kas Operasional Kidinest)',
      suggestedCreditAccount: '4-1100 (Pendapatan Jasa Medis Kidinest)',
      details: {
        totalPasien: 12,
        tunai: 4200000,
        edcCard: 5250000,
        kasir: 'Fitri Handayani'
      }
    },
    {
      id: 'tx-002',
      transactionNumber: 'TX-CSH-20260304-002',
      transactionDate: '2026-03-04 15:30:00',
      sourceModule: 'GRO_CASHIER',
      referenceNumber: 'SHIFT-SSC-01',
      businessUnit: 'SSC',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Closing Kasir Shift Siang SSC - Treatment Laser & Skin Rejuvenation',
      amount: 16800000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '1-1101 (Kas Operasional SSC)',
      suggestedCreditAccount: '4-1101 (Pendapatan Jasa Medis SSC)',
      details: {
        totalTreatment: 8,
        edcCard: 14500000,
        qris: 2300000,
        kasir: 'Nadia Safitri'
      }
    },
    {
      id: 'tx-003',
      transactionNumber: 'TX-PHARM-20260304-08',
      transactionDate: '2026-03-04 11:20:00',
      sourceModule: 'PHARMACY',
      referenceNumber: 'RX-KD-8819',
      businessUnit: 'KIDINEST',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Dispensing Obat Resep Anak (Antibiotik, Puyer Paracetamol & Multivitamin)',
      amount: 1250000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '1-1100 (Kas Operasional Kidinest)',
      suggestedCreditAccount: '4-1200 (Pendapatan Farmasi & Obat)',
      details: {
        noResep: 'RX-KD-8819',
        namaPasien: 'Ananda Kevin (4 Th)',
        dokter: 'dr. Hendra Wijaya, Sp.A',
        estimasiHPP: 680000
      }
    },
    {
      id: 'tx-004',
      transactionNumber: 'TX-PO-20260303-15',
      transactionDate: '2026-03-03 16:00:00',
      sourceModule: 'PURCHASING',
      referenceNumber: 'INV-KALBE-9901',
      businessUnit: 'ALL',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Faktur Pengadaan Suplemen & Vaksin dari PT Kalbe Farma Tbk (Termin 30 Hari)',
      amount: 22200000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '1-1400 (Persediaan Farmasi & Obat)',
      suggestedCreditAccount: '2-1100 (Hutang Usaha Distributor Farmasi)',
      details: {
        supplier: 'PT Kalbe Farma Tbk',
        fakturPajak: '010.002-26.88219034',
        dpp: 20000000,
        ppn11: 2200000,
        jatuhTempo: '2026-04-02'
      }
    },
    {
      id: 'tx-005',
      transactionNumber: 'TX-DOC-20260304-03',
      transactionDate: '2026-03-04 17:00:00',
      sourceModule: 'DOCTOR_SHARE',
      referenceNumber: 'FEE-DR-HENDRA-01',
      businessUnit: 'KIDINEST',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Kalkulasi Fee Jasa Medis dr. Hendra Wijaya, Sp.A (14 Pasien Konsultasi)',
      amount: 5600000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '5-1200 (Beban Jasa Medis Dokter)',
      suggestedCreditAccount: '2-1200 (Hutang Jasa Medis Dokter)',
      details: {
        dokter: 'dr. Hendra Wijaya, Sp.A',
        npwp: '08.123.456.7-012.000',
        feeBruto: 5600000,
        pph21Tarif: 0.05,
        pph21Potongan: 140000,
        feeNetto: 5460000
      }
    },
    {
      id: 'tx-006',
      transactionNumber: 'TX-PAY-20260301-01',
      transactionDate: '2026-03-01 09:00:00',
      sourceModule: 'PAYROLL',
      referenceNumber: 'PR-202603-KID',
      businessUnit: 'KIDINEST',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penggajian Staf Perawat & Administrasi Kidinest Maret 2026',
      amount: 34500000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '6-1100 (Gaji & Upah Karyawan)',
      suggestedCreditAccount: '1-1200 (Bank Mandiri Operasional)',
      details: {
        jumlahPegawai: 8,
        gajiPokok: 30000000,
        tunjangan: 4500000
      }
    },
    {
      id: 'tx-007',
      transactionNumber: 'TX-FA-20260301-02',
      transactionDate: '2026-03-01 08:00:00',
      sourceModule: 'FIXED_ASSET',
      referenceNumber: 'FA-SSC-LASER-01',
      businessUnit: 'SSC',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Penyusutan Bulanan Mesin Pico Laser Aesthetic SSC',
      amount: 6250000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '6-1400 (Beban Penyusutan Aset Tetap)',
      suggestedCreditAccount: '1-2100 (Akumulasi Penyusutan Alat Medis)',
      details: {
        assetName: 'PicoSure Laser Aesthetic Workstation',
        hargaPerolehan: 300000000,
        masaManfaatBulan: 48,
        metode: 'Garis Lurus'
      }
    },
    {
      id: 'tx-008',
      transactionNumber: 'TX-BNK-20260304-05',
      transactionDate: '2026-03-04 16:30:00',
      sourceModule: 'BANK',
      referenceNumber: 'SETOR-TUNAI-MANDIRI-01',
      businessUnit: 'ALL',
      branchName: 'Kebayoran Baru (Pusat)',
      description: 'Setoran Tunai Hasil Penerimaan Kasir ke Rekening Bank Mandiri Operasional',
      amount: 15000000,
      mappingStatus: 'MAPPED',
      journalStatus: 'UNPROCESSED',
      suggestedDebitAccount: '1-1200 (Bank Mandiri Operasional)',
      suggestedCreditAccount: '1-1100 (Kas Operasional Kidinest)',
      details: {
        bankTujuan: 'Bank Mandiri (123-00-9876543-1)',
        penyetor: 'Fitri Handayani',
        validasiSlip: 'SLIP-MDR-20260304-77'
      }
    }
  ])

  // Filtered list based on global unit
  const filteredTransactions = computed(() => {
    return transactions.value.filter(t => {
      if (accountingStore.businessUnit !== 'ALL' && t.businessUnit !== 'ALL' && t.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return true
    })
  })

  // Action: 1-click Generate Draft Journal from Operational Transaction
  function generateDraftJournal(txId: string) {
    const tx = transactions.value.find(t => t.id === txId)
    if (!tx) {
      accountingStore.showToast('Transaksi tidak ditemukan', 'error')
      return
    }

    if (tx.journalStatus !== 'UNPROCESSED') {
      accountingStore.showToast(`Transaksi ${tx.transactionNumber} sudah pernah dibuatkan jurnal draft (${tx.generatedJournalId})`, 'info')
      return
    }

    // Determine account line items based on transaction details & mapping
    const lines: any[] = []
    const txDate = tx.transactionDate.split(' ')[0]

    if (tx.sourceModule === 'GRO_CASHIER') {
      const isKidinest = tx.businessUnit === 'KIDINEST'
      const cashAcc = isKidinest
        ? { code: '1-1100', name: 'Kas Operasional Kidinest' }
        : { code: '1-1101', name: 'Kas Operasional SSC' }
      const revAcc = isKidinest
        ? { code: '4-1100', name: 'Pendapatan Jasa Medis Kidinest' }
        : { code: '4-1101', name: 'Pendapatan Jasa Medis SSC' }

      lines.push({
        accountId: `acc-${cashAcc.code}`,
        accountCode: cashAcc.code,
        accountName: cashAcc.name,
        debit: tx.amount,
        credit: 0,
        description: `Penerimaan Kasir - ${tx.referenceNumber}`
      })
      lines.push({
        accountId: `acc-${revAcc.code}`,
        accountCode: revAcc.code,
        accountName: revAcc.name,
        debit: 0,
        credit: tx.amount,
        description: `Pendapatan Pelayanan Pasien Klinik`
      })
    } else if (tx.sourceModule === 'DOCTOR_SHARE') {
      const feeBruto = tx.details.feeBruto || tx.amount
      const pph21 = tx.details.pph21Potongan || Math.round(feeBruto * 0.025)
      const feeNet = feeBruto - pph21

      lines.push({
        accountId: 'acc-51200',
        accountCode: '5-1200',
        accountName: 'Beban Jasa Medis Dokter',
        debit: feeBruto,
        credit: 0,
        description: `Jasa Medis ${tx.details.dokter}`
      })
      lines.push({
        accountId: 'acc-21300',
        accountCode: '2-1300',
        accountName: 'Hutang PPh 21 Dokter',
        debit: 0,
        credit: pph21,
        description: `Potongan PPh 21 Pasal 17 PMK 168/2023`
      })
      lines.push({
        accountId: 'acc-21200',
        accountCode: '2-1200',
        accountName: 'Hutang Jasa Medis Dokter',
        debit: 0,
        credit: feeNet,
        description: `Hutang Bersih ke Dokter`
      })
    } else if (tx.sourceModule === 'PURCHASING') {
      const dpp = tx.details.dpp || Math.round(tx.amount / 1.11)
      const ppn = tx.amount - dpp

      lines.push({
        accountId: 'acc-11400',
        accountCode: '1-1400',
        accountName: 'Persediaan Farmasi & Obat',
        debit: dpp,
        credit: 0,
        description: `Pengadaan Obat DPP - ${tx.details.supplier}`
      })
      lines.push({
        accountId: 'acc-11180',
        accountCode: '1-1180',
        accountName: 'PPN Masukan (Dapat Dikreditkan)',
        debit: ppn,
        credit: 0,
        description: `PPN Masukan 11% Faktur ${tx.details.fakturPajak}`
      })
      lines.push({
        accountId: 'acc-21100',
        accountCode: '2-1100',
        accountName: 'Hutang Usaha Distributor Farmasi',
        debit: 0,
        credit: tx.amount,
        description: `Hutang Dagang ${tx.details.supplier}`
      })
    } else {
      // Generic balanced 2-line entry
      const debitCode = tx.suggestedDebitAccount.split(' ')[0]
      const debitName = tx.suggestedDebitAccount.replace(`${debitCode} (`, '').replace(')', '')
      const creditCode = tx.suggestedCreditAccount.split(' ')[0]
      const creditName = tx.suggestedCreditAccount.replace(`${creditCode} (`, '').replace(')', '')

      lines.push({
        accountId: `acc-${debitCode}`,
        accountCode: debitCode,
        accountName: debitName,
        debit: tx.amount,
        credit: 0,
        description: tx.description
      })
      lines.push({
        accountId: `acc-${creditCode}`,
        accountCode: creditCode,
        accountName: creditName,
        debit: 0,
        credit: tx.amount,
        description: tx.description
      })
    }

    const createdJournal = journalsStore.createManualJournal({
      journalDate: txDate,
      businessUnit: tx.businessUnit,
      branchId: 'BRANCH-MAIN',
      branchName: tx.branchName,
      description: `[Auto-Generated dari ${tx.sourceModule}] ${tx.description}`,
      lines
    })

    // Update operational transaction
    tx.journalStatus = 'DRAFT_GENERATED'
    tx.generatedJournalId = createdJournal.journalNumber
    accountingStore.showToast(`Draft Jurnal ${createdJournal.journalNumber} berhasil dibuat dari transaksi ${tx.transactionNumber}!`, 'success')
  }

  return {
    transactions,
    mappingRules,
    filteredTransactions,
    generateDraftJournal
  }
})
