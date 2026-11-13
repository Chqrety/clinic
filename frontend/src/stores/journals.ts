import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { JournalEntry, JournalStatus, BusinessUnit, JournalSource } from '@/types/accounting'
import { useAccountingStore } from './accounting'

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
          accountId: 'acc-1100',
          accountCode: '1-1100',
          accountName: 'Kas Operasional Kidinest',
          debit: 12500000,
          credit: 0,
          description: 'Penerimaan Kas Kasir Kidinest'
        },
        {
          id: 'jl-101-2',
          accountId: 'acc-1201',
          accountCode: '1-1201',
          accountName: 'Bank BCA Operasional',
          debit: 6000000,
          credit: 0,
          description: 'Penerimaan EDC & QRIS BCA'
        },
        {
          id: 'jl-101-3',
          accountId: 'acc-4100',
          accountCode: '4-1100',
          accountName: 'Pendapatan Jasa Medis Kidinest',
          debit: 0,
          credit: 14000000,
          description: 'Konsultasi Spesialis Anak & Tumbuh Kembang'
        },
        {
          id: 'jl-101-4',
          accountId: 'acc-4200',
          accountCode: '4-1200',
          accountName: 'Pendapatan Farmasi & Obat',
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
          accountId: 'acc-1101',
          accountCode: '1-1101',
          accountName: 'Kas Operasional SSC',
          debit: 7800000,
          credit: 0,
          description: 'Penerimaan Tunai Kasir SSC'
        },
        {
          id: 'jl-102-2',
          accountId: 'acc-1201',
          accountCode: '1-1201',
          accountName: 'Bank BCA Operasional',
          debit: 20000000,
          credit: 0,
          description: 'Penerimaan Gesek EDC Kartu Kredit & Debit'
        },
        {
          id: 'jl-102-3',
          accountId: 'acc-4101',
          accountCode: '4-1101',
          accountName: 'Pendapatan Jasa Medis SSC',
          debit: 0,
          credit: 18000000,
          description: 'Tindakan Laser & Aesthetic Dokter Sp.KK'
        },
        {
          id: 'jl-102-4',
          accountId: 'acc-4102',
          accountCode: '4-1300',
          accountName: 'Pendapatan Skincare & Produk Slimming',
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
    }
  ])

  // Filtered journals based on global filters
  const filteredJournals = computed(() => {
    return journals.value.filter(j => {
      if (accountingStore.businessUnit !== 'ALL' && j.businessUnit !== 'ALL' && j.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      if (accountingStore.branchCode !== 'ALL' && j.branchId !== accountingStore.branchCode) {
        return false
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
    branchId: string
    branchName: string
    description: string
    lines: {
      accountId: string
      accountCode: string
      accountName: string
      debit: number
      credit: number
      description: string
    }[]
  }): JournalEntry {
    const totalDebit = data.lines.reduce((s, l) => s + (l.debit || 0), 0)
    const totalCredit = data.lines.reduce((s, l) => s + (l.credit || 0), 0)

    if (totalDebit !== totalCredit) {
      throw new Error(`Total Debit (Rp ${totalDebit}) tidak seimbang dengan Total Kredit (Rp ${totalCredit})`)
    }

    const nextSeq = journals.value.length + 1
    const seqStr = String(nextSeq).padStart(3, '0')
    const jNumber = `JV-2026-03-${seqStr}`

    const newJournal: JournalEntry = {
      id: 'jv-' + Date.now(),
      journalNumber: jNumber,
      journalDate: data.journalDate,
      journalSource: 'MANUAL',
      sourceReference: `MAN-${Date.now().toString().slice(-6)}`,
      businessUnit: data.businessUnit,
      branchId: data.branchId,
      branchName: data.branchName,
      description: data.description,
      totalDebit,
      totalCredit,
      status: 'DRAFT',
      lines: data.lines.map((l, idx) => ({
        id: `jl-${idx}-${Date.now()}`,
        ...l
      })),
      createdBy: accountingStore.currentUserEmail,
      createdByName: accountingStore.currentUserName,
      createdAt: new Date().toISOString(),
      isManual: true
    }

    journals.value.unshift(newJournal)
    accountingStore.showToast(`Draft jurnal manual ${jNumber} berhasil dibuat!`, 'success')
    return newJournal
  }

  return {
    journals,
    filteredJournals,
    stats,
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
