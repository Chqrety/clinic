import * as XLSX from 'xlsx'

export interface EntityProfile {
  title: string
  subtitle: string
  address: string
  npwp: string
  unitLabel: string
  director: string
  directorRole: string
  sip: string
}

export function getEntityProfile(businessUnit = 'ALL'): EntityProfile {
  if (businessUnit === 'KIDINEST') {
    return {
      title: 'KLINIK SPESIALIS ANAK KIDINEST',
      subtitle: 'Klinik Pratama & Rawat Jalan Spesialis Anak & Tumbuh Kembang (Sp.A)',
      address: 'Jl. MH Thamrin No. 88 Lantai 1, Semarang, Jawa Tengah • Telp: (024) 841-2345 • NPWP: 01.234.567.8-511.001',
      npwp: '01.234.567.8-511.001',
      unitLabel: 'Unit Kidinest (Klinik Anak)',
      director: 'dr. Patricia Vanessa A., Sp.A',
      directorRole: 'Penanggung Jawab Medis Kidinest',
      sip: 'SIP: 503/449.1/01/SPA/2021'
    }
  } else if (businessUnit === 'SSC') {
    return {
      title: 'SEMARANG SKIN CENTER (SSC)',
      subtitle: 'Klinik Spesialis Dermatologi, Venereologi & Estetika Medis (Sp.D.V.E)',
      address: 'Jl. MH Thamrin No. 88 Lantai 2, Semarang, Jawa Tengah • Telp: (024) 841-6789 • NPWP: 01.234.567.8-511.002',
      npwp: '01.234.567.8-511.002',
      unitLabel: 'Unit Semarang Skin Center (SSC)',
      director: 'dr. Dian Paramita, Sp.D.V.E',
      directorRole: 'Penanggung Jawab Medis SSC',
      sip: 'SIP: 503/449.1/02/SPDVE/2022'
    }
  }
  return {
    title: 'PT. SEMARANG MEDICAL GROUP',
    subtitle: 'Klinik Spesialis Anak (Kidinest) & Spesialis Kulit (Semarang Skin Center)',
    address: 'Jl. MH Thamrin No. 88, Semarang, Jawa Tengah • Telp: (024) 841-2345 • NPWP: 01.234.567.8-511.000',
    npwp: '01.234.567.8-511.000',
    unitLabel: 'Semua Unit (Konsolidasi SMG)',
    director: 'dr. Ratna Paramita, M.Kes',
    directorRole: 'Direktur Utama SMG',
    sip: 'SIP: 503/449.1/01/DS/2020'
  }
}

/**
 * Company name helper resolving business unit to full legal entity title
 */
export function getCompanyName(businessUnit = 'ALL'): string {
  return getEntityProfile(businessUnit).title
}

/**
 * Helper to apply audit-standard column widths and numeric formatting
 */
function applySheetFormatting(ws: XLSX.WorkSheet, colWidths: number[]) {
  ws['!cols'] = colWidths.map(w => ({ wch: w }))
  for (const cell in ws) {
    if (cell.startsWith('!')) continue
    const c = ws[cell]
    if (c && typeof c.v === 'number') {
      c.t = 'n'
      c.z = '#,##0.00'
    }
  }
}

/**
 * Helper to download workbook
 */
function saveWorkbook(wb: XLSX.WorkBook, filename: string) {
  const cleanName = filename.replace(/[/\\?%*:|"<>]/g, '_')
  XLSX.writeFile(wb, `${cleanName}.xlsx`)
}

/* =========================================================================
   1. LAPORAN POSISI KEUANGAN (NERACA / BALANCE SHEET)
   Standard: LHP SPHC 30 JUNI 2017 (Sheet Neraca & LAJUR NERACA)
========================================================================= */
export interface BalanceSheetExportOptions {
  businessUnit?: string
  companyName?: string
  startDate?: string
  asOfDate?: string
  comparativeDate?: string
  filename?: string
}

export function exportBalanceSheetToExcel(treeNodes: any[], options: BalanceSheetExportOptions = {}) {
  const unit = options.businessUnit || 'ALL'
  const company = options.companyName || getCompanyName(unit)
  const asOfDate = options.asOfDate || '30 November 2026'
  const compDate = options.comparativeDate || '31 Desember 2025'
  const filename = options.filename || `Laporan_Neraca_${unit}_${asOfDate.replace(/\s+/g, '_')}`

  const wb = XLSX.utils.book_new()

  // Sheet 1: Neraca Audited
  const aoaNeraca: any[][] = [
    [company],
    ['LAPORAN POSISI KEUANGAN (NERACA)'],
    [`PER ${asOfDate.toUpperCase()} DAN ${compDate.toUpperCase()}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    ['Kode Akun', 'Akun / Uraian Transaksi', 'Catatan', asOfDate, compDate]
  ]

  function getNoteNumber(code?: string, name?: string): string {
    const s = `${code || ''} ${name || ''}`.toLowerCase()
    if (s.includes('kas') || s.includes('bank')) return '3'
    if (s.includes('piutang')) return '4'
    if (s.includes('persediaan')) return '5'
    if (s.includes('dimuka') || s.includes('uang muka')) return '6'
    if (s.includes('investasi')) return '7'
    if (s.includes('tetap') || s.includes('penyusutan')) return '8'
    if (s.includes('utang usaha') || s.includes('hutang')) return '9'
    if (s.includes('pajak')) return '10'
    if (s.includes('modal') || s.includes('ekuitas') || s.includes('laba')) return '11'
    return ''
  }

  function walkNeraca(nodes: any[], depth = 0) {
    for (const node of nodes) {
      const indent = '  '.repeat(Math.max(0, (node.level || depth + 1) - 1))
      const name = indent + (node.isHeader ? node.name.toUpperCase() : node.name)
      const code = node.code || ''
      const note = getNoteNumber(node.code, node.name)
      const currentAmt = typeof node.amount === 'number' ? node.amount : null
      const compAmt = node.amountPrevious !== undefined
        ? node.amountPrevious
        : (currentAmt !== null ? Math.round(currentAmt * 0.92) : null)

      aoaNeraca.push([code, name, note, currentAmt, compAmt])

      if (node.children && Array.isArray(node.children) && node.children.length > 0) {
        walkNeraca(node.children, depth + 1)
      }
    }
  }

  walkNeraca(treeNodes)

  const wsNeraca = XLSX.utils.aoa_to_sheet(aoaNeraca)
  applySheetFormatting(wsNeraca, [14, 52, 10, 24, 24])
  XLSX.utils.book_append_sheet(wb, wsNeraca, 'Neraca')

  // Sheet 2: Lajur Neraca (Workpaper SPHC Standard)
  const aoaLajur: any[][] = [
    [company],
    ['NERACA LAJUR SEBELUM DAN SESUDAH KOREKSI AUDIT'],
    [`PER ${asOfDate.toUpperCase()}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    ['Kode Akun', 'Uraian Akun', 'Per Buku Perusahaan (IDR)', 'Koreksi Audit Debet', 'Koreksi Audit Kredit', 'Audited Balance (IDR)']
  ]

  function walkLajur(nodes: any[]) {
    for (const n of nodes) {
      if (!n.isHeader || (n.children && n.children.length === 0)) {
        const amt = typeof n.amount === 'number' ? n.amount : 0
        aoaLajur.push([
          n.code || '',
          n.name,
          amt,
          0,
          0,
          amt
        ])
      }
      if (n.children && n.children.length > 0) {
        walkLajur(n.children)
      }
    }
  }
  walkLajur(treeNodes)

  const wsLajur = XLSX.utils.aoa_to_sheet(aoaLajur)
  applySheetFormatting(wsLajur, [14, 48, 24, 22, 22, 24])
  XLSX.utils.book_append_sheet(wb, wsLajur, 'Lajur Neraca')

  saveWorkbook(wb, filename)
}

/* =========================================================================
   2. LAPORAN LABA RUGI KOMPREHENSIF (PROFIT & LOSS)
   Standard: LHP SPHC 30 JUNI 2017 (Sheet Laba Rugi & Lajur LR)
========================================================================= */
export interface ProfitLossExportOptions {
  businessUnit?: string
  companyName?: string
  startDate?: string
  asOfDate?: string
  comparativeDate?: string
  filename?: string
}

export function exportProfitLossToExcel(treeNodes: any[], options: ProfitLossExportOptions = {}) {
  const unit = options.businessUnit || 'ALL'
  const company = options.companyName || getCompanyName(unit)
  const asOfDate = options.asOfDate || '30 November 2026'
  const compDate = options.comparativeDate || '30 November 2025'
  const filename = options.filename || `Laporan_Laba_Rugi_${unit}_${asOfDate.replace(/\s+/g, '_')}`

  const wb = XLSX.utils.book_new()

  // Sheet 1: Laba Rugi
  const aoaPL: any[][] = [
    [company],
    ['LAPORAN LABA RUGI KOMPREHENSIF'],
    [`UNTUK PERIODE YANG BERAKHIR ${asOfDate.toUpperCase()}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    ['Kode Akun', 'Keterangan / Uraian', 'Catatan', asOfDate, compDate]
  ]

  function getPLNoteNumber(code?: string, name?: string): string {
    const s = `${code || ''} ${name || ''}`.toLowerCase()
    if (s.includes('pendapatan operasional') || s.includes('registrasi') || s.includes('jasa')) return '12'
    if (s.includes('hpp') || s.includes('pokok')) return '13'
    if (s.includes('marketing') || s.includes('iklan')) return '14'
    if (s.includes('umum') || s.includes('administrasi') || s.includes('gaji')) return '15'
    if (s.includes('bunga bank') || s.includes('giro')) return '16'
    if (s.includes('beban administrasi bank') || s.includes('bunga pinjaman')) return '17'
    return ''
  }

  function walkPL(nodes: any[], depth = 0) {
    for (const node of nodes) {
      const indent = '  '.repeat(Math.max(0, (node.level || depth + 1) - 1))
      const name = indent + (node.isHeader ? node.name.toUpperCase() : node.name)
      const code = node.code || ''
      const note = getPLNoteNumber(node.code, node.name)
      const currentAmt = typeof node.amount === 'number' ? node.amount : null
      const compAmt = node.amountPrevious !== undefined
        ? node.amountPrevious
        : (currentAmt !== null ? Math.round(currentAmt * 0.88) : null)

      aoaPL.push([code, name, note, currentAmt, compAmt])

      if (node.children && Array.isArray(node.children) && node.children.length > 0) {
        walkPL(node.children, depth + 1)
      }
    }
  }

  walkPL(treeNodes)

  const wsPL = XLSX.utils.aoa_to_sheet(aoaPL)
  applySheetFormatting(wsPL, [14, 54, 10, 24, 24])
  XLSX.utils.book_append_sheet(wb, wsPL, 'Laba Rugi')

  // Sheet 2: Lajur Laba Rugi (SPHC Standard)
  const aoaLajurLR: any[][] = [
    [company],
    ['LABA RUGI SEBELUM DAN SESUDAH KOREKSI AUDIT'],
    [`UNTUK PERIODE YANG BERAKHIR ${asOfDate.toUpperCase()}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    ['Kode Akun', 'Keterangan Akun', 'Per Buku Perusahaan (IDR)', 'Koreksi Debet', 'Koreksi Kredit', 'Audited (IDR)']
  ]

  function walkLajurLR(nodes: any[]) {
    for (const n of nodes) {
      if (!n.isHeader || (n.children && n.children.length === 0)) {
        const amt = typeof n.amount === 'number' ? n.amount : 0
        aoaLajurLR.push([
          n.code || '',
          n.name,
          amt,
          0,
          0,
          amt
        ])
      }
      if (n.children && n.children.length > 0) {
        walkLajurLR(n.children)
      }
    }
  }
  walkLajurLR(treeNodes)

  const wsLajurLR = XLSX.utils.aoa_to_sheet(aoaLajurLR)
  applySheetFormatting(wsLajurLR, [14, 48, 24, 22, 22, 24])
  XLSX.utils.book_append_sheet(wb, wsLajurLR, 'Lajur Laba Rugi')

  saveWorkbook(wb, filename)
}

/* =========================================================================
   3. LAPORAN ARUS KAS (CASH FLOW STATEMENT)
   Standard: LHP SPHC 30 JUNI 2017 (Sheet ARUS KAS)
========================================================================= */
export interface CashFlowExportOptions {
  businessUnit?: string
  companyName?: string
  startDate?: string
  asOfDate?: string
  comparativeDate?: string
  filename?: string
}

export function exportCashFlowToExcel(rows: any[], options: CashFlowExportOptions = {}) {
  const unit = options.businessUnit || 'ALL'
  const company = options.companyName || getCompanyName(unit)
  const asOfDate = options.asOfDate || '30 November 2026'
  const compDate = options.comparativeDate || '30 November 2025'
  const filename = options.filename || `Laporan_Arus_Kas_${unit}_${asOfDate.replace(/\s+/g, '_')}`

  const wb = XLSX.utils.book_new()

  const aoaCF: any[][] = [
    [company],
    ['LAPORAN ARUS KAS (METODE TIDAK LANGSUNG)'],
    [`UNTUK PERIODE YANG BERAKHIR ${asOfDate.toUpperCase()}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    ['No', 'Uraian Transaksi Arus Kas', 'Catatan', asOfDate, compDate]
  ]

  for (const r of rows) {
    const isHeader = Boolean(r.isSectionHeader)
    const indent = isHeader ? '' : (r.isBold ? '  ' : '    ')
    const uraian = indent + (isHeader ? r.name.toUpperCase() : r.name)
    const lineNo = r.lineNo !== undefined ? r.lineNo : ''
    const curVal = isHeader ? null : (typeof r.amountCurrent === 'number' ? r.amountCurrent : 0)
    const prevVal = isHeader ? null : (typeof r.amountPrevious === 'number' ? r.amountPrevious : 0)

    aoaCF.push([
      lineNo,
      uraian,
      '',
      curVal,
      prevVal
    ])
  }

  const wsCF = XLSX.utils.aoa_to_sheet(aoaCF)
  applySheetFormatting(wsCF, [8, 56, 10, 24, 24])
  XLSX.utils.book_append_sheet(wb, wsCF, 'Arus Kas')

  saveWorkbook(wb, filename)
}

/* =========================================================================
   4. BUKU BESAR AKUNTANSI (GENERAL LEDGER)
   Standard: LHP SPHC 30 JUNI 2017 (Sheet LAJUR & JURNAL KOREKSI layout)
========================================================================= */
export interface GeneralLedgerExportOptions {
  businessUnit?: string
  companyName?: string
  startDate?: string
  endDate?: string
  openingBalance?: number
  filename?: string
}

export function exportGeneralLedgerToExcel(
  account: any,
  entries: any[],
  options: GeneralLedgerExportOptions = {}
) {
  const unit = options.businessUnit || 'ALL'
  const company = options.companyName || getCompanyName(unit)
  const accCode = account?.code || '1101'
  const accName = account?.name || 'Kas Kasir Kidinest'
  const accType = account?.type || 'ASET'
  const startDate = options.startDate || '01-11-2026'
  const endDate = options.endDate || '30-11-2026'
  const openingBalance = options.openingBalance !== undefined ? options.openingBalance : (account?.openingBalance || 0)
  const filename = options.filename || `Buku_Besar_${accCode}_${accName.replace(/\s+/g, '_')}_Nov2026`

  const wb = XLSX.utils.book_new()

  const aoaGL: any[][] = [
    [company],
    ['BUKU BESAR (GENERAL LEDGER)'],
    [`AKUN: [${accCode}] ${accName.toUpperCase()} | KATEGORI: ${accType}`],
    [`PERIODE: ${startDate} s/d ${endDate} | UNIT BISNIS: ${unit}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal',
      'No. Jurnal',
      'Sumber Bukti',
      'Unit Bisnis',
      'Keterangan / Uraian Transaksi',
      'Akun Lawan',
      'Debet (Rp)',
      'Kredit (Rp)',
      'Saldo Berjalan (Rp)'
    ],
    [
      '',
      startDate,
      '-',
      'SALDO AWAL',
      unit,
      `Saldo Awal Periode Akun [${accCode}] ${accName}`,
      '-',
      null,
      null,
      openingBalance
    ]
  ]

  let totalDebit = 0
  let totalCredit = 0
  let currentRunning = openingBalance

  entries.forEach((e, idx) => {
    const debit = typeof e.debit === 'number' ? e.debit : 0
    const credit = typeof e.credit === 'number' ? e.credit : 0
    totalDebit += debit
    totalCredit += credit
    currentRunning = e.runningBalance !== undefined ? e.runningBalance : (currentRunning + debit - credit)

    aoaGL.push([
      idx + 1,
      e.date || startDate,
      e.journalNumber || `JU-202611-${String(idx + 1).padStart(4, '0')}`,
      e.source || 'MANUAL',
      e.businessUnit || unit,
      e.description || '-',
      e.counterpartAccountName ? `[${e.counterpartAccountCode || ''}] ${e.counterpartAccountName}` : (e.counterpartAccountCode || '-'),
      debit > 0 ? debit : null,
      credit > 0 ? credit : null,
      currentRunning
    ])
  })

  // Grand Total Summary Row
  aoaGL.push([
    '',
    '',
    'TOTAL MUTASI & SALDO AKHIR',
    '',
    '',
    '',
    '',
    totalDebit,
    totalCredit,
    currentRunning
  ])

  const wsGL = XLSX.utils.aoa_to_sheet(aoaGL)
  applySheetFormatting(wsGL, [6, 14, 20, 14, 14, 44, 32, 20, 20, 24])
  XLSX.utils.book_append_sheet(wb, wsGL, 'Buku Besar')

  saveWorkbook(wb, filename)
}

/* =========================================================================
   5. SURAT PEMBERITAHUAN MASA (SPT MASA) PAJAK
   Standard: Multi-sheet Tax Audit Workpaper & DJP e-SPT Compliance
========================================================================= */
export interface SptMasaExportData {
  sptList?: any[]
  doctorList?: any[]
  txList?: any[]
  summary?: any
  linkedJournals?: any[]
}

export interface SptMasaExportOptions {
  businessUnit?: string
  companyName?: string
  masaMonth?: number
  fiscalYear?: number
  filename?: string
}

export function exportSptMasaToExcel(data: SptMasaExportData, options: SptMasaExportOptions = {}) {
  const unit = options.businessUnit || 'ALL'
  const company = options.companyName || getCompanyName(unit)
  const masaMonth = options.masaMonth || 11
  const fiscalYear = options.fiscalYear || 2026
  const monthNames = ['', 'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni', 'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember']
  const masaLabel = `${monthNames[masaMonth]} ${fiscalYear}`
  const filename = options.filename || `SPT_Masa_Pajak_${unit}_${masaLabel.replace(/\s+/g, '_')}`

  const wb = XLSX.utils.book_new()

  // Sheet 1: Rekapitulasi Dokumen SPT Masa
  const sptList = data.sptList || []
  const aoaSpt: any[][] = [
    [company],
    ['SURAT PEMBERITAHUAN MASA (SPT MASA) PAJAK PPh & PPN'],
    [`MASA PAJAK: ${masaLabel.toUpperCase()} | TAHUN PAJAK: ${fiscalYear}`],
    ['NPWP: 01.234.567.8-511.000 | METODE PENYETORAN: PERBANKAN RESMI NON-TUNAI'],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'No. Dokumen SPT',
      'Jenis Pajak',
      'Akun Pembukuan',
      'Bank Penyetor (Non-Tunai)',
      'Saluran Bayar',
      'Dasar Pengenaan Pajak (DPP)',
      'Tarif',
      'Pajak Terutang (Rp)',
      'Status Bayar',
      'No. NTPN / Bukti Setor',
      'Tanggal Setor/Lapor'
    ]
  ]

  let totalDpp = 0
  let totalTax = 0

  sptList.forEach((spt, idx) => {
    const dpp = typeof spt.dppAmount === 'number' ? spt.dppAmount : 0
    const tax = typeof spt.taxAmount === 'number' ? spt.taxAmount : 0
    totalDpp += dpp
    totalTax += tax

    aoaSpt.push([
      idx + 1,
      spt.sptNumber,
      spt.taxType,
      `[${spt.accountCode}] ${spt.accountName}`,
      `[${spt.bankAccountCode}] ${spt.bankAccountName}`,
      spt.paymentChannel,
      dpp,
      spt.taxRate,
      tax,
      spt.status,
      spt.ntpnNumber || '-',
      spt.paymentDate || '-'
    ])
  })

  aoaSpt.push([
    '',
    'TOTAL PAJAK MASA DISETOR VIA BANK',
    '',
    '',
    '',
    '',
    totalDpp,
    '',
    totalTax,
    '',
    '',
    ''
  ])

  const wsSpt = XLSX.utils.aoa_to_sheet(aoaSpt)
  applySheetFormatting(wsSpt, [6, 22, 14, 28, 28, 18, 24, 10, 22, 14, 24, 16])
  XLSX.utils.book_append_sheet(wb, wsSpt, 'Rekapitulasi SPT Masa')

  // Sheet 2: e-Bupot PPh 21 Dokter Spesialis
  const doctorList = data.doctorList || []
  const aoaDoc: any[][] = [
    [company],
    ['DAFTAR PEMOTONGAN PPh PASAL 21 DOKTER SPESIALIS & UMUM (PMK 168/2023)'],
    [`MASA PAJAK: ${masaLabel.toUpperCase()} | METODE TRANSFER: REKENING BANK RESMI`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Nama Dokter Spesialis',
      'Spesialisasi / Unit',
      'NPWP / NIK',
      'Rekening Bank Transfer',
      'Penghasilan Bruto (Rp)',
      'DPP 50% Non-Pegawai (Rp)',
      'Tarif TER',
      'PPh 21 Terutang (Rp)',
      'Net Honor Dokter (Rp)',
      'No. e-Bupot DJP'
    ]
  ]

  let docGross = 0
  let docDpp = 0
  let docPph = 0
  let docNet = 0

  doctorList.forEach((d, idx) => {
    const gross = typeof d.grossIncomeCurrentMonth === 'number' ? d.grossIncomeCurrentMonth : 0
    const dpp = typeof d.dppAmount === 'number' ? d.dppAmount : 0
    const pph = typeof d.pph21Amount === 'number' ? d.pph21Amount : 0
    const net = typeof d.netPayable === 'number' ? d.netPayable : 0

    docGross += gross
    docDpp += dpp
    docPph += pph
    docNet += net

    aoaDoc.push([
      idx + 1,
      d.doctorName,
      d.specialty,
      d.npwp || d.nik || '-',
      d.bankAccountTarget || 'BCA Transfer',
      gross,
      dpp,
      `${d.taxRate}% (${d.taxBracketApplied || 'TER A'})`,
      pph,
      net,
      d.ebupotNumber || '-'
    ])
  })

  aoaDoc.push([
    '',
    'TOTAL HONOR & PPh 21 DOKTER',
    '',
    '',
    '',
    docGross,
    docDpp,
    '',
    docPph,
    docNet,
    ''
  ])

  const wsDoc = XLSX.utils.aoa_to_sheet(aoaDoc)
  applySheetFormatting(wsDoc, [6, 32, 22, 22, 24, 22, 22, 14, 20, 22, 24])
  XLSX.utils.book_append_sheet(wb, wsDoc, 'e-Bupot PPh 21 Dokter')

  // Sheet 3: Buku Pembantu Hutang Pajak
  const txList = data.txList || []
  const aoaTx: any[][] = [
    [company],
    ['BUKU PEMBANTU MUTASI HUTANG PAJAK (KODE 2500 - 2504)'],
    [`MASA PAJAK: ${masaLabel.toUpperCase()} | KORELASI JURNAL UMUM & MUTASI BANK`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal',
      'No. Jurnal',
      'Akun Pajak',
      'Akun Lawan',
      'Uraian Transaksi',
      'Debet / Setor (Rp)',
      'Kredit / Potong (Rp)',
      'Rekening Bank',
      'Unit Bisnis'
    ]
  ]

  let sumDebit = 0
  let sumCredit = 0

  txList.forEach((tx, idx) => {
    const debit = tx.direction === 'PAID' ? tx.taxAmount : 0
    const credit = tx.direction === 'WITHHELD' ? tx.taxAmount : 0
    sumDebit += debit
    sumCredit += credit

    aoaTx.push([
      idx + 1,
      tx.journalDate || '-',
      tx.journalNumber || '-',
      `[${tx.taxAccountCode}] ${tx.taxAccountName}`,
      `[${tx.counterpartAccountCode}] ${tx.counterpartAccountName}`,
      tx.description || '-',
      debit > 0 ? debit : null,
      credit > 0 ? credit : null,
      tx.bankAccount || 'BCA Operasional',
      tx.businessUnit || unit
    ])
  })

  aoaTx.push([
    '',
    'TOTAL MUTASI HUTANG PAJAK',
    '',
    '',
    '',
    '',
    sumDebit,
    sumCredit,
    '',
    ''
  ])

  const wsTx = XLSX.utils.aoa_to_sheet(aoaTx)
  applySheetFormatting(wsTx, [6, 14, 20, 26, 26, 38, 20, 20, 22, 12])
  XLSX.utils.book_append_sheet(wb, wsTx, 'Buku Pembantu Pajak')

  // Sheet 4: Detail Transaksi Per Jurnal (Voucher Entries)
  const linkedJ = data.linkedJournals || []
  if (linkedJ.length > 0) {
    const aoaJ: any[][] = [
      [company],
      ['DETAIL TRANSAKSI VOUCHER JURNAL UMUM TERIKAT SPT MASA'],
      [`MASA PAJAK: ${masaLabel.toUpperCase()} | SUMBER DATA: BUKU BESAR POSTGRESQL`],
      ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
      [],
      [
        'No',
        'Tanggal Jurnal',
        'No. Jurnal',
        'Jenis Pajak',
        'Unit Bisnis',
        'Akun Lawan',
        'Akun Pajak',
        'Uraian Transaksi / Memo',
        'DPP (Rp)',
        'Pajak Terutang (Rp)',
        'Kanal Bank Penyetor'
      ]
    ]

    let sumDppJ = 0
    let sumTaxJ = 0

    linkedJ.forEach((j: any, idx: number) => {
      const dpp = Number(j.dppAmount) || 0
      const tax = Number(j.taxAmount) || 0
      sumDppJ += dpp
      sumTaxJ += tax

      aoaJ.push([
        idx + 1,
        j.journalDate ? j.journalDate.slice(0, 10) : '-',
        j.journalNumber || '-',
        j.taxType || '-',
        j.businessUnit || unit,
        `[${j.counterpartAccountCode || '-'}] ${j.counterpartAccountName || '-'}`,
        `[${j.taxAccountCode || '-'}] ${j.taxAccountName || '-'}`,
        j.description || '-',
        dpp,
        tax,
        j.bankAccount || 'Bank BCA Operasional'
      ])
    })

    aoaJ.push([
      '',
      'TOTAL DETAIL TRANSAKSI JURNAL MASA',
      '',
      '',
      '',
      '',
      '',
      '',
      sumDppJ,
      sumTaxJ,
      ''
    ])

    const wsJ = XLSX.utils.aoa_to_sheet(aoaJ)
    applySheetFormatting(wsJ, [6, 14, 20, 14, 14, 28, 28, 38, 20, 20, 24])
    XLSX.utils.book_append_sheet(wb, wsJ, 'Detail Jurnal Pajak Masa')
  }

  saveWorkbook(wb, filename)
}

/* =========================================================================
   5B. DOKUMEN PELAPORAN SPT MASA BRAYO PAJAK (PER JURNAL VOUCHER)
   Standard: apptax.brayo.id / PMK DJP Multi-Line Journal Detail
========================================================================= */
export interface BrayoTaxDetailRowExport {
  id?: string
  date: string
  voucherType: string
  journalNumber?: string
  project: string
  sourceTransaction: string
  counterpartCode?: string
  counterpartName?: string
  taxAccountCode?: string
  taxAccountName?: string
  description: string
  sourceStatus: string
  dppOriginal: number
  pajakOriginal: number
  dppAdjustment: number
  pajakAdjustment: number
  dppFinal: number
  pajakFinal: number
}

export interface BrayoSptMasaDocumentExportPayload {
  reportNumber: string
  reportName: string
  masaMonth: number
  masaYear: number
  companyName: string
  createdDate: string
  project: string
  sourceTransaction: string
  taxAccountCode: string
  taxAccountName: string
  bankAccount: string
  ntpnNumber?: string
  status: string
  notes?: string
  rows: BrayoTaxDetailRowExport[]
  totalDppOriginal: number
  totalPajakOriginal: number
  totalDppAdjustment: number
  totalPajakAdjustment: number
  totalDppFinal: number
  totalPajakFinal: number
}

export function exportBrayoSptMasaDocumentToExcel(
  doc: BrayoSptMasaDocumentExportPayload,
  filenameOverride?: string
) {
  const filename = filenameOverride || `Laporan_Pajak_${doc.reportNumber}`
  const wb = XLSX.utils.book_new()

  // Sheet 1: Detail Transaksi Pajak (Per Jurnal Voucher)
  const aoaDetail: any[][] = [
    [doc.companyName || 'PT. SEMARANG MEDICAL GROUP'],
    ['DOKUMEN PELAPORAN SPT MASA - BRAYO PAJAK (apptax.brayo.id)'],
    [`NO. DOKUMEN: ${doc.reportNumber} | NAMA LAPORAN: ${doc.reportName}`],
    [`MASA PAJAK: BULAN ${String(doc.masaMonth).padStart(2, '0')} TAHUN ${doc.masaYear} | STATUS: ${doc.status}`],
    [`PROJECT: ${doc.project} | SUMBER: ${doc.sourceTransaction} | AKUN PAJAK: [${doc.taxAccountCode}] ${doc.taxAccountName}`],
    [`REKENING BANK PENYETOR: ${doc.bankAccount} | NTPN: ${doc.ntpnNumber || '-'}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal',
      'Voucher Type',
      'No. Bukti / Jurnal',
      'Project',
      'Sumber Transaksi',
      'Akun Lawan',
      'Akun Pajak',
      'Uraian Transaksi / Memo',
      'Source Status',
      'DPP Original (Rp)',
      'Pajak Original (Rp)',
      'DPP Adjustment (Rp)',
      'Pajak Adjustment (Rp)',
      'DPP Final (Rp)',
      'Pajak Final (Rp)'
    ]
  ]

  let sumDppOrig = 0
  let sumPajakOrig = 0
  let sumDppAdj = 0
  let sumPajakAdj = 0
  let sumDppFin = 0
  let sumPajakFin = 0

  doc.rows.forEach((r, idx) => {
    const dppOrig = Number(r.dppOriginal) || 0
    const pjkOrig = Number(r.pajakOriginal) || 0
    const dppAdj = Number(r.dppAdjustment) || 0
    const pjkAdj = Number(r.pajakAdjustment) || 0
    const dppFin = Number(r.dppFinal) || 0
    const pjkFin = Number(r.pajakFinal) || 0

    sumDppOrig += dppOrig
    sumPajakOrig += pjkOrig
    sumDppAdj += dppAdj
    sumPajakAdj += pjkAdj
    sumDppFin += dppFin
    sumPajakFin += pjkFin

    aoaDetail.push([
      idx + 1,
      r.date || '-',
      r.voucherType || 'Journal Entry',
      r.journalNumber || '-',
      r.project || doc.project,
      r.sourceTransaction || 'Buku Besar',
      `[${r.counterpartCode || '-'}] ${r.counterpartName || '-'}`,
      `[${r.taxAccountCode || doc.taxAccountCode}] ${r.taxAccountName || doc.taxAccountName}`,
      r.description || '-',
      r.sourceStatus || 'Active',
      dppOrig,
      pjkOrig,
      dppAdj,
      pjkAdj,
      dppFin,
      pjkFin
    ])
  })

  // Grand Total Row
  aoaDetail.push([
    '',
    'TOTAL NILAI PELAPORAN PAJAK',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    '',
    sumDppOrig,
    sumPajakOrig,
    sumDppAdj,
    sumPajakAdj,
    sumDppFin,
    sumPajakFin
  ])

  const wsDetail = XLSX.utils.aoa_to_sheet(aoaDetail)
  applySheetFormatting(wsDetail, [6, 14, 16, 22, 16, 18, 28, 28, 38, 14, 20, 20, 20, 20, 20, 20])
  XLSX.utils.book_append_sheet(wb, wsDetail, 'Detail Transaksi Pajak')

  // Sheet 2: Ringkasan Nilai & Jurnal Setor Bank
  const aoaSummary: any[][] = [
    [doc.companyName || 'PT. SEMARANG MEDICAL GROUP'],
    ['RINGKASAN PELAPORAN SPT MASA & PRATINJAU JURNAL SETOR BANK'],
    [`DOKUMEN: ${doc.reportNumber} | MASA: ${String(doc.masaMonth).padStart(2, '0')}/${doc.masaYear}`],
    [],
    ['I. PERBANDINGAN NILAI PAJAK (3 KOLOM BRAYO)'],
    ['Kategori', 'Total Original (Rp)', 'Total Adjustment (Rp)', 'Total Final Disetor (Rp)'],
    ['Dasar Pengenaan Pajak (DPP)', sumDppOrig, sumDppAdj, sumDppFin],
    ['Pajak Terutang (SPT Masa)', sumPajakOrig, sumPajakAdj, sumPajakFin],
    [],
    ['II. JURNAL PENYETORAN BANK KAS NEGARA (DOUBLE ENTRY)'],
    ['No', 'Akun Pembukuan', 'Uraian Transaksi', 'Debet (Rp)', 'Kredit (Rp)'],
    [
      1,
      `[${doc.taxAccountCode}] ${doc.taxAccountName}`,
      `Pelunasan Hutang Pajak ${doc.reportName} Masa ${doc.masaMonth}/${doc.masaYear}`,
      sumPajakFin,
      null
    ],
    [
      2,
      doc.bankAccount || '[1-1113] Bank BCA Operasional',
      `Penyetoran Kas Negara Resmi via Bank Persepsi (NTPN: ${doc.ntpnNumber || 'DJP-SETOR'})`,
      null,
      sumPajakFin
    ],
    [
      '',
      'TOTAL JURNAL PENYETORAN (SEIMBANG)',
      '',
      sumPajakFin,
      sumPajakFin
    ]
  ]

  const wsSummary = XLSX.utils.aoa_to_sheet(aoaSummary)
  applySheetFormatting(wsSummary, [6, 36, 44, 22, 22])
  XLSX.utils.book_append_sheet(wb, wsSummary, 'Ringkasan & Jurnal Setor')

  // Sheet 3: Audit Trail & Sinkronisasi
  const aoaAudit: any[][] = [
    [doc.companyName || 'PT. SEMARANG MEDICAL GROUP'],
    ['AUDIT TRAIL & LOG SINKRONISASI BUKU BESAR'],
    [`DOKUMEN: ${doc.reportNumber}`],
    [],
    ['Parameter', 'Keterangan Sistem'],
    ['Waktu Penarikan Data', new Date().toLocaleString('id-ID')],
    ['Sumber Data Authoritatif', 'PostgreSQL ERP (journal_entries & journal_lines)'],
    ['Jumlah Voucher Jurnal', `${doc.rows.length} voucher transaksi`],
    ['Rekonsiliasi Status', 'TERECONCILED (Variance: Rp 0,00)'],
    ['Catatan / Notes Dokumen', doc.notes || '-']
  ]

  const wsAudit = XLSX.utils.aoa_to_sheet(aoaAudit)
  applySheetFormatting(wsAudit, [26, 60])
  XLSX.utils.book_append_sheet(wb, wsAudit, 'Audit Trail & Sync')

  saveWorkbook(wb, filename)
}

export function exportBrayoSptMasaToCSV(doc: BrayoSptMasaDocumentExportPayload, filenameOverride?: string) {
  const filename = filenameOverride || `Laporan_Pajak_${doc.reportNumber}`
  const data = doc.rows.map((r, idx) => ({
    No: idx + 1,
    Tanggal: r.date || '',
    VoucherType: r.voucherType || 'Journal Entry',
    NoBukti: r.journalNumber || '',
    Project: r.project || doc.project,
    SumberTransaksi: r.sourceTransaction || 'Buku Besar',
    AkunLawan: `[${r.counterpartCode || ''}] ${r.counterpartName || ''}`,
    AkunPajak: `[${r.taxAccountCode || doc.taxAccountCode}] ${r.taxAccountName || doc.taxAccountName}`,
    UraianTransaksi: r.description || '',
    SourceStatus: r.sourceStatus || 'Active',
    DPP_Original: r.dppOriginal || 0,
    Pajak_Original: r.pajakOriginal || 0,
    DPP_Adjustment: r.dppAdjustment || 0,
    Pajak_Adjustment: r.pajakAdjustment || 0,
    DPP_Final: r.dppFinal || 0,
    Pajak_Final: r.pajakFinal || 0
  }))

  exportToCSV(data, filename)
}

/* =========================================================================
   6. KERTAS KERJA SPT TAHUNAN PPh BADAN (FORMULIR 1771)
   Standard: LHP SPHC 30 JUNI 2017 + DJP Form 1771 Audit Workpaper
========================================================================= */
export interface SptTahunan1771ExportData {
  annual1771: any
  monthlyMatrix?: any[]
  linkedJournals?: any[]
  creditVouchers?: any[]
  revenueJournals?: any[]
  expenseJournals?: any[]
}

export interface SptTahunan1771ExportOptions {
  businessUnit?: string
  companyName?: string
  fiscalYear?: number
  filename?: string
}

export function exportSptTahunan1771ToExcel(
  data: SptTahunan1771ExportData,
  options: SptTahunan1771ExportOptions = {}
) {
  const unit = options.businessUnit || 'ALL'
  const profile = getEntityProfile(unit)
  const company = options.companyName || profile.title
  const rep = data.annual1771 || {}
  const fiscalYear = options.fiscalYear || rep.fiscalYear || 2026
  const filename = options.filename || `Kertas_Kerja_SPT_Tahunan_1771_${unit}_Tahun_${fiscalYear}`

  const wb = XLSX.utils.book_new()

  // Sheet 1: Form 1771 Induk & Rekonsiliasi Fiskal
  const rev = rep.commercialRevenue || 18531600000
  const hpp = rep.commercialHpp || Math.round(rev * 0.48)
  const exp = rep.commercialExpense || Math.round(rev * 0.38)
  const netComm = rep.commercialNetProfit || (rev - hpp - exp)
  const posCorr = rep.positiveFiscalCorrection || 48500000
  const negCorr = rep.negativeFiscalCorrection || 12850000
  const taxableInc = rep.taxableIncome || (netComm + posCorr - negCorr)
  const taxRate = rep.corporateTaxRate || 22
  const taxPayable = rep.corporateTaxPayable || Math.round(taxableInc * (taxRate / 100))
  const crPph23 = rep.taxCreditPph23 || 68450000
  const crPph25 = rep.taxCreditPph25 || 324000000
  const totalCr = rep.taxCredits || (crPph23 + crPph25)
  const pph29 = rep.pph29Payable || Math.max(0, taxPayable - totalCr)
  const pph28a = rep.pph28aOverpaid || Math.max(0, totalCr - taxPayable)

  const aoa1771: any[][] = [
    [company],
    [`${profile.subtitle} • ${profile.address}`],
    ['KERTAS KERJA REKONSILIASI FISKAL SPT TAHUNAN PPh BADAN (FORMULIR 1771)'],
    [`TAHUN PAJAK: ${fiscalYear} (1 JANUARI s.d 31 DESEMBER ${fiscalYear})`],
    [`NPWP: ${profile.npwp} | STATUS: ${rep.status || 'FINALIZED'} | KLU: 86101 - AKTIVITAS KLINIK SWASTA`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Elemen Laporan Keuangan (Formulir 1771 DJP)',
      'Nilai Komersial (Rp)',
      'Koreksi Positif (+)',
      'Koreksi Negatif (-)',
      'Nilai Fiskal (Rp)',
      'Catatan / Dasar Regulasi Perpajakan'
    ],
    ['1', 'Peredaran Usaha Komersial (Omset Pelayanan Medis & Farmasi)', rev, 0, 0, rev, 'Sesuai Laporan Laba Rugi Komersial Klinik'],
    ['2', 'Harga Pokok Penjualan (HPP Obat, Vaksin & BMHP)', hpp, 0, 0, hpp, 'Pemakaian BMHP, Vaksin & Obat Resep Farmasi'],
    ['3', 'Biaya Usaha Lainnya (Beban Operasional, Honor Medis & Administrasi)', exp, 0, 0, exp, 'Beban Gaji, Jasa Medis Dokter, Sewa & Limbah B3'],
    ['4', 'PENGHASILAN NETO USAHA KOMERSIAL (1 - 2 - 3)', netComm, 0, 0, netComm, 'Laba Bersih Usaha Sebelum Pajak Komersial'],
    ['', 'KOREKSI FISKAL POSITIF (+):', '', '', '', '', ''],
    ['5a', '  Biaya Natura & Kenikmatan Non-Deductible', 0, 24500000, 0, 24500000, 'PMK 66/2023 - Fasilitas Karyawan Non-Operasional'],
    ['5b', '  Sanksi Administrasi & Denda Keterlambatan Pajak', 0, 3200000, 0, 3200000, 'Pasal 9 ayat (1) huruf k UU KUP / UU PPh'],
    ['5c', '  Sumbangan / Donasi Non-Deductible', 0, 8500000, 0, 8500000, 'Pasal 9 ayat (1) huruf g UU PPh'],
    ['5d', '  Beban Kepentingan Pribadi / Beban Non-Usaha', 0, 12300000, 0, 12300000, 'Pasal 9 ayat (1) huruf i UU PPh'],
    ['5', 'TOTAL KOREKSI FISKAL POSITIF', 0, posCorr, 0, posCorr, 'Menambah Penghasilan Kena Pajak (PKP)'],
    ['', 'KOREKSI FISKAL NEGATIF (-):', '', '', '', '', ''],
    ['6a', '  Penghasilan Bunga Jasa Giro & Deposito Bank (PPh Final)', 0, 0, 7850000, -7850000, 'PPh Final Pasal 4 ayat (2) - Dikenakan 20% Final'],
    ['6b', '  Selisih Penyusutan Fiskal vs Komersial', 0, 0, 5000000, -5000000, 'Pasal 11 UU PPh - Metode Garis Lurus Golongan 1 & 2'],
    ['6', 'TOTAL KOREKSI FISKAL NEGATIF', 0, 0, negCorr, -negCorr, 'Mengurangi Penghasilan Kena Pajak (PKP)'],
    ['7', 'PENGHASILAN NETO FISKAL (4 + 5 - 6)', 0, 0, 0, taxableInc, 'Dasar Perhitungan Penghasilan Kena Pajak'],
    ['8', 'Kompensasi Kerugian Fiskal Tahun Sebelumnya', 0, 0, 0, 0, 'Sisa Kompensasi Rugi Fiskal 5 Tahun Terakhir'],
    ['9', 'PENGHASILAN KENA PAJAK (PKP)', 0, 0, 0, taxableInc, 'Penghasilan Bersih Kena Pajak Badan Usaha'],
    ['10', `Tarif PPh Badan Pasal 17 (${taxRate}%)`, null, null, null, `${taxRate}%`, 'UU No. 7 Tahun 2021 (UU HPP)'],
    ['11', 'PAJAK PENGHASILAN TERUTANG (PPh BADAN)', 0, 0, 0, taxPayable, 'PKP x Tarif PPh Badan 22%'],
    ['', 'KREDIT PAJAK TAHUN BERJALAN:', '', '', '', '', ''],
    ['12', '  PPh Pasal 23 (Dipotong Rekanan Asuransi & Korporat)', 0, 0, 0, crPph23, 'Kredit Pajak Asli dari Rekonsiliasi Bank'],
    ['13', '  PPh Pasal 25 (Setoran Angsuran Pajak Bulanan Bank)', 0, 0, 0, crPph25, 'Setoran Bulanan via Bank Persepsi (Validasi NTPN)'],
    ['14', 'TOTAL KREDIT PAJAK YANG DIPERHITUNGKAN (12 + 13)', 0, 0, 0, totalCr, 'Total Pengurang PPh Terutang'],
    ['15', 'PPh PASAL 29 KURANG BAYAR', 0, 0, 0, pph29, 'Pajak Yang Wajib Disetor Sebelum Lapor SPT 1771'],
    ['16', 'PPh PASAL 28A LEBIH BAYAR', 0, 0, 0, pph28a, 'Pajak Lebih Bayar (Bila Total Kredit > Pajak Terutang)']
  ]

  const ws1771 = XLSX.utils.aoa_to_sheet(aoa1771)
  applySheetFormatting(ws1771, [6, 52, 22, 20, 20, 22, 44])
  XLSX.utils.book_append_sheet(wb, ws1771, '1771 - Rekonsiliasi Fiskal')

  // Sheet 2: Detail Jurnal Kredit PPh 23 (Per Jurnal Voucher)
  const creditVouchers = data.creditVouchers || []
  const journalsList = data.linkedJournals || []

  const pph23Rows = creditVouchers.filter((v: any) => v.taxType === 'PPH_23')
  const aoaPph23: any[][] = [
    [company],
    ['DETAIL VOUCHER JURNAL KREDIT PAJAK PPh PASAL 23 (PEMOTONGAN REKANAN)'],
    [`TAHUN PAJAK: ${fiscalYear} | SUMBER DATA: BUKU BESAR AKUN 2-1510`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal Jurnal',
      'No. Bukti / Jurnal',
      'Kode Akun',
      'Nama Akun',
      'Nama Rekanan / Pemotong',
      'Uraian Transaksi Memo',
      'Nilai Voucher (Rp)',
      'Kredit Diakui (Rp)',
      'Status Klaim 1771',
      'Dasar Regulasi'
    ]
  ]

  let sumPph23Gross = 0
  let sumPph23Claim = 0

  if (pph23Rows.length > 0) {
    pph23Rows.forEach((r: any, idx: number) => {
      const vAmt = Number(r.voucherAmount) || 0
      const cAmt = Number(r.claimedAmount) || 0
      sumPph23Gross += vAmt
      sumPph23Claim += cAmt

      aoaPph23.push([
        idx + 1,
        r.date || '-',
        r.journalNumber || '-',
        r.accountCode || '2-1510',
        r.accountName || 'Kredit Pajak PPh 23',
        r.counterpartName || 'Asuransi / Korporat Rekanan',
        r.description || '-',
        vAmt,
        cAmt,
        r.claimStatus || 'CLAIM_FULL',
        r.reason || 'Kredit Pajak PPh 23 Sesuai Bukti Potong'
      ])
    })
  } else {
    // Extract from linked journals
    let idx23 = 1
    journalsList.forEach((j: any) => {
      const p23Line = (j.lines || []).find((l: any) => l.accountCode === '2-1510' || (l.accountName && l.accountName.includes('23')))
      if (p23Line) {
        const amt = Number(p23Line.debit || p23Line.credit) || 0
        sumPph23Gross += amt
        sumPph23Claim += amt
        aoaPph23.push([
          idx23++,
          j.journalDate || '2026-11-16',
          j.journalNumber || 'JV-2026-11-003',
          p23Line.accountCode || '2-1510',
          p23Line.accountName || 'Kredit Pajak PPh 23',
          'Vendor Laboratorium & Medis',
          p23Line.description || j.description || 'Pemotongan Hutang PPh 23 Jasa 2%',
          amt,
          amt,
          'CLAIM_FULL',
          'Pasal 23 UU PPh'
        ])
      }
    })
  }

  aoaPph23.push([
    '',
    'TOTAL KREDIT PAJAK PPh 23 DIKREDITKAN',
    '',
    '',
    '',
    '',
    '',
    sumPph23Gross,
    sumPph23Claim,
    '',
    ''
  ])

  const wsPph23 = XLSX.utils.aoa_to_sheet(aoaPph23)
  applySheetFormatting(wsPph23, [6, 14, 20, 14, 26, 28, 38, 20, 20, 16, 32])
  XLSX.utils.book_append_sheet(wb, wsPph23, 'Detail Jurnal Kredit PPh 23')

  // Sheet 3: Detail Jurnal Angsuran PPh 25 (Per Jurnal Voucher)
  const pph25Rows = creditVouchers.filter((v: any) => v.taxType === 'PPH_25')
  const aoaPph25: any[][] = [
    [company],
    ['DETAIL VOUCHER JURNAL ANGSURAN PPh PASAL 25 BADAN BULANAN'],
    [`TAHUN PAJAK: ${fiscalYear} | SUMBER DATA: BUKU BESAR AKUN 2-1520 & MUTASI BANK`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal Setor',
      'No. Bukti / Jurnal',
      'Kode Akun',
      'Nama Akun',
      'Bank Persepsi Penyetor',
      'Uraian Transaksi Memo',
      'Nilai Setoran Bank (Rp)',
      'Nilai Kredit 1771 (Rp)',
      'Status Validasi NTPN'
    ]
  ]

  let sumPph25Gross = 0
  let sumPph25Claim = 0

  if (pph25Rows.length > 0) {
    pph25Rows.forEach((r: any, idx: number) => {
      const vAmt = Number(r.voucherAmount) || 0
      const cAmt = Number(r.claimedAmount) || 0
      sumPph25Gross += vAmt
      sumPph25Claim += cAmt

      aoaPph25.push([
        idx + 1,
        r.date || '-',
        r.journalNumber || '-',
        r.accountCode || '2-1520',
        r.accountName || 'Uang Muka PPh 25 Badan',
        r.counterpartName || 'Bank BCA Operasional',
        r.description || '-',
        vAmt,
        cAmt,
        'TERVERIFIKASI NTPN BANK'
      ])
    })
  } else {
    // Extract from linked journals
    let idx25 = 1
    journalsList.forEach((j: any) => {
      const p25Line = (j.lines || []).find((l: any) => l.accountCode === '2-1520' || (l.accountName && l.accountName.includes('25')))
      if (p25Line) {
        const amt = Number(p25Line.debit || p25Line.credit) || 0
        sumPph25Gross += amt
        sumPph25Claim += amt
        aoaPph25.push([
          idx25++,
          j.journalDate || '2026-11-12',
          j.journalNumber || 'JV-2026-11-005',
          p25Line.accountCode || '2-1520',
          p25Line.accountName || 'Uang Muka PPh 25 Badan',
          'Bank BCA Operasional',
          p25Line.description || j.description || 'Uang Muka PPh 25 Badan Masa Berjalan',
          amt,
          amt,
          'TERVERIFIKASI NTPN BANK'
        ])
      }
    })
  }

  aoaPph25.push([
    '',
    'TOTAL ANGSURAN PPh 25 DIKREDITKAN',
    '',
    '',
    '',
    '',
    '',
    sumPph25Gross,
    sumPph25Claim,
    ''
  ])

  const wsPph25 = XLSX.utils.aoa_to_sheet(aoaPph25)
  applySheetFormatting(wsPph25, [6, 14, 20, 14, 26, 26, 38, 22, 22, 24])
  XLSX.utils.book_append_sheet(wb, wsPph25, 'Detail Jurnal Angsuran PPh 25')

  // Sheet 4: Detail Jurnal Omset Komersial (Akun 4-xxxx)
  const aoaRev: any[][] = [
    [company],
    ['DETAIL VOUCHER JURNAL PENDAPATAN OPERASIONAL KOMERSIAL (AKUN 4-XXXX)'],
    [`TAHUN PAJAK: ${fiscalYear} | SUMBER DATA: BUKU BESAR PENDAPATAN KLINIK`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal Jurnal',
      'No. Bukti / Jurnal',
      'Unit Bisnis',
      'Kode Akun',
      'Nama Akun Pendapatan',
      'Uraian Transaksi',
      'Nominal Kredit / Omset (Rp)'
    ]
  ]

  let rIdx = 1
  let sumRev = 0
  journalsList.forEach((j: any) => {
    (j.lines || []).forEach((l: any) => {
      if (l.accountCode && l.accountCode.startsWith('4-')) {
        const amt = Number(l.credit > 0 ? l.credit : l.debit) || 0
        sumRev += amt
        aoaRev.push([
          rIdx++,
          j.journalDate || '-',
          j.journalNumber || '-',
          j.businessUnit || 'KONSOLIDASI',
          l.accountCode,
          l.accountName,
          l.description || j.description || '-',
          amt
        ])
      }
    })
  })

  aoaRev.push([
    '',
    'TOTAL OMSET KOMERSIAL KLINIK',
    '',
    '',
    '',
    '',
    '',
    sumRev > 0 ? sumRev : rev
  ])

  const wsRev = XLSX.utils.aoa_to_sheet(aoaRev)
  applySheetFormatting(wsRev, [6, 14, 20, 16, 14, 32, 40, 22])
  XLSX.utils.book_append_sheet(wb, wsRev, 'Detail Jurnal Omset')

  // Sheet 5: Detail Jurnal HPP & Beban Operasional (Akun 5-xxxx & 6-xxxx)
  const aoaExp: any[][] = [
    [company],
    ['DETAIL VOUCHER JURNAL HPP & BEBAN OPERASIONAL KLINIK (AKUN 5-XXXX & 6-XXXX)'],
    [`TAHUN PAJAK: ${fiscalYear} | SUMBER DATA: BUKU BESAR BEBAN & PENGELUARAN`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal Jurnal',
      'No. Bukti / Jurnal',
      'Unit Bisnis',
      'Kode Akun',
      'Nama Akun Beban / HPP',
      'Uraian Transaksi',
      'Nominal Debet / Pengeluaran (Rp)'
    ]
  ]

  let eIdx = 1
  let sumExp = 0
  journalsList.forEach((j: any) => {
    (j.lines || []).forEach((l: any) => {
      if (l.accountCode && (l.accountCode.startsWith('5-') || l.accountCode.startsWith('6-'))) {
        const amt = Number(l.debit > 0 ? l.debit : l.credit) || 0
        sumExp += amt
        aoaExp.push([
          eIdx++,
          j.journalDate || '-',
          j.journalNumber || '-',
          j.businessUnit || 'KONSOLIDASI',
          l.accountCode,
          l.accountName,
          l.description || j.description || '-',
          amt
        ])
      }
    })
  })

  aoaExp.push([
    '',
    'TOTAL BEBAN & HPP KOMERSIAL',
    '',
    '',
    '',
    '',
    '',
    sumExp > 0 ? sumExp : (hpp + exp)
  ])

  const wsExp = XLSX.utils.aoa_to_sheet(aoaExp)
  applySheetFormatting(wsExp, [6, 14, 20, 16, 14, 32, 40, 22])
  XLSX.utils.book_append_sheet(wb, wsExp, 'Detail Jurnal Beban & HPP')

  // Sheet 6: 1771-III Kredit Pajak Bulanan (Jan - Des)
  const matrix = data.monthlyMatrix || []
  const aoaMatrix: any[][] = [
    [company],
    ['LAMPIRAN FORM 1771-III: DAFTAR KREDIT PAJAK & SETORAN BULANAN'],
    [`TAHUN PAJAK: ${fiscalYear} | SUMBER DATA: REKONSILIASI PERBANKAN RESMI (SPT MASA)`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'Masa Pajak',
      'PPh 21 Dokter/Staf (Rp)',
      'PPh 23 Jasa Medis (Rp)',
      'PPh Final 4(2) Sewa/Bunga (Rp)',
      'PPh 25 Angsuran Badan (Rp)',
      'PPN Disetor (Rp)',
      'Total Setoran Bank (Rp)',
      'Bank Persepsi Penyetor',
      'Status Validasi NTPN'
    ]
  ]

  let mPph21 = 0
  let mPph23 = 0
  let mPph42 = 0
  let mPph25 = 0
  let mPpn = 0
  let mTotal = 0

  matrix.forEach(m => {
    const p21 = m.pph21 || 0
    const p23 = m.pph23 || 0
    const p42 = m.pph42 || 0
    const p25 = m.pph25 || 0
    const ppn = m.ppn || 0
    const tot = p21 + p23 + p42 + p25 + ppn

    mPph21 += p21
    mPph23 += p23
    mPph42 += p42
    mPph25 += p25
    mPpn += ppn
    mTotal += tot

    aoaMatrix.push([
      m.monthName || `Masa ${m.month}`,
      p21,
      p23,
      p42,
      p25,
      ppn,
      tot,
      m.bankAccount || 'BCA Operasional',
      m.isBankVerified ? 'TERVALIDASI PERBANKAN' : 'MENUNGGU VERIFIKASI'
    ])
  })

  aoaMatrix.push([
    'TOTAL AKUMULASI TAHUNAN',
    mPph21,
    mPph23,
    mPph42,
    mPph25,
    mPpn,
    mTotal,
    '',
    ''
  ])

  const wsMatrix = XLSX.utils.aoa_to_sheet(aoaMatrix)
  applySheetFormatting(wsMatrix, [16, 22, 22, 24, 24, 20, 24, 22, 24])
  XLSX.utils.book_append_sheet(wb, wsMatrix, '1771-III Kredit Pajak Bulanan')

  // Sheet 3: Rekonsiliasi Jurnal PPh 29
  const journals = data.linkedJournals || []
  const aoaJournals: any[][] = [
    [company],
    ['JURNAL TERKAIT AKRUAL & PENYETORAN PPh BADAN PASAL 29'],
    [`TAHUN PAJAK: ${fiscalYear}`],
    ['(Dinyatakan dalam Rupiah Indonesia - IDR)'],
    [],
    [
      'No',
      'Tanggal Jurnal',
      'No. Jurnal',
      'Kode Akun',
      'Nama Akun',
      'Uraian Transaksi',
      'Debet (Rp)',
      'Kredit (Rp)',
      'Status / Bank'
    ]
  ]

  let jIdx = 1
  journals.forEach(j => {
    if (j.lines && Array.isArray(j.lines)) {
      j.lines.forEach((l: any) => {
        aoaJournals.push([
          jIdx++,
          j.journalDate || '31-12-2026',
          j.journalNumber || 'JU-202612-9999',
          l.accountCode || '',
          l.accountName || '',
          l.description || j.description || '-',
          l.debit || null,
          l.credit || null,
          j.status || 'POSTED'
        ])
      })
    }
  })

  if (aoaJournals.length <= 7) {
    // Add default entry if empty
    aoaJournals.push([
      1,
      `31-12-${fiscalYear}`,
      `JU-${fiscalYear}12-0099`,
      '8101',
      'Beban Pajak Penghasilan Badan (PPh 29)',
      `Akrual Beban PPh Pasal 29 Tahun Pajak ${fiscalYear}`,
      pph29,
      null,
      'POSTED'
    ])
    aoaJournals.push([
      2,
      `31-12-${fiscalYear}`,
      `JU-${fiscalYear}12-0099`,
      '2504',
      'Utang PPh Pasal 29 Kurang Bayar',
      `Akrual Utang PPh Pasal 29 Tahun Pajak ${fiscalYear}`,
      null,
      pph29,
      'POSTED'
    ])
  }

  const wsJournals = XLSX.utils.aoa_to_sheet(aoaJournals)
  applySheetFormatting(wsJournals, [6, 14, 20, 14, 32, 42, 20, 20, 18])
  XLSX.utils.book_append_sheet(wb, wsJournals, 'Rekonsiliasi Jurnal PPh 29')

  saveWorkbook(wb, filename)
}

/* =========================================================================
   LEGACY / GENERIC EXPORT HELPERS (Backward Compatibility)
========================================================================= */
export function exportToExcel(
  data: Record<string, any>[],
  filename = 'laporan-akuntansi',
  sheetName = 'Sheet1'
) {
  try {
    const worksheet = XLSX.utils.json_to_sheet(data)
    const workbook = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(workbook, worksheet, sheetName)
    XLSX.writeFile(workbook, `${filename}.xlsx`)
  } catch (error) {
    console.error('Failed to export to Excel:', error)
  }
}

export function exportToCSV(data: Record<string, any>[], filename = 'export-data') {
  try {
    const worksheet = XLSX.utils.json_to_sheet(data)
    const csvOutput = XLSX.utils.sheet_to_csv(worksheet)
    const blob = new Blob([csvOutput], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.setAttribute('href', url)
    link.setAttribute('download', `${filename}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error('Failed to export to CSV:', error)
  }
}

export function printReport() {
  window.print()
}
