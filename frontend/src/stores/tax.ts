import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TaxTransaction, BusinessUnit } from '@/types/accounting'
import { useAccountingStore } from './accounting'
import { exportToCSV } from '@/utils/export'

export interface DoctorTaxRecord {
  id: string
  doctorName: string
  specialty: string
  npwp: string
  nik: string
  businessUnit: BusinessUnit
  grossIncomeCurrentMonth: number
  grossIncomeCumulativeYear: number
  dppPercentage: number // 50% for non-pegawai dokter
  dppAmount: number
  taxBracketApplied: string
  taxRate: number
  pph21Amount: number
  netPayable: number
  ebupotNumber: string
}

export const useTaxStore = defineStore('tax', () => {
  const accountingStore = useAccountingStore()

  // Doctor Tax Calculations per PMK 168/2023
  const doctorTaxRecords = ref<DoctorTaxRecord[]>([
    {
      id: 'dt-1',
      doctorName: 'dr. Hendra Wijaya, Sp.A',
      specialty: 'Spesialis Anak (Pediatri)',
      npwp: '08.123.456.7-012.000',
      nik: '3171012304850001',
      businessUnit: 'KIDINEST',
      grossIncomeCurrentMonth: 38500000,
      grossIncomeCumulativeYear: 112000000,
      dppPercentage: 0.5,
      dppAmount: 19250000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%)',
      taxRate: 0.05,
      pph21Amount: 962500,
      netPayable: 37537500,
      ebupotNumber: 'BUPOT-21-202603-001'
    },
    {
      id: 'dt-2',
      doctorName: 'dr. Maya Sandria, Sp.KK, FINSDV',
      specialty: 'Spesialis Dermatologi & Venereologi',
      npwp: '09.345.678.9-015.000',
      nik: '3174055208880003',
      businessUnit: 'SSC',
      grossIncomeCurrentMonth: 64200000,
      grossIncomeCumulativeYear: 185000000,
      dppPercentage: 0.5,
      dppAmount: 32100000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%) & Lapisan 2 (15%)',
      taxRate: 0.05,
      pph21Amount: 1605000,
      netPayable: 62595000,
      ebupotNumber: 'BUPOT-21-202603-002'
    },
    {
      id: 'dt-3',
      doctorName: 'dr. Nadia Kirana, Sp.BP-RE',
      specialty: 'Spesialis Bedah Plastik & Rekonstruksi',
      npwp: '12.876.543.2-021.000',
      nik: '3276014407910005',
      businessUnit: 'SSC',
      grossIncomeCurrentMonth: 42300000,
      grossIncomeCumulativeYear: 98000000,
      dppPercentage: 0.5,
      dppAmount: 21150000,
      taxBracketApplied: 'Pasal 17 Lapisan 1 (5%)',
      taxRate: 0.05,
      pph21Amount: 1057500,
      netPayable: 41242500,
      ebupotNumber: 'BUPOT-21-202603-003'
    }
  ])

  // PPN Summary (PPN Masukan vs PPN Keluaran)
  const ppnSummary = ref({
    ppnOutputGross: 113636364, // DPP Penjualan Obat Non-Resep & Skincare
    ppnOutputTax: 12500000, // 11%
    ppnInputGross: 76363636, // DPP Faktur Pajak Masukan Pembelian Distributor
    ppnInputTax: 8400000, // 11%
    ppnPayableNet: 4100000 // Kurang Bayar (Output - Input)
  })

  // Tax transactions list for reconciliation
  const taxTransactions = ref<TaxTransaction[]>([
    {
      id: 'tax-01',
      taxNumber: '010.002-26.11029384',
      taxDate: '2026-03-01',
      businessUnit: 'SSC',
      taxCategory: 'PPN_OUTPUT',
      counterpartyName: 'Pasien Treatment Konsumen Akhir',
      counterpartyNPWP: '00.000.000.0-000.000',
      objectCode: 'PPN-KELUARAN-11',
      grossAmount: 48000000,
      dppAmount: 43243243,
      taxRate: 0.11,
      taxAmount: 4756757,
      documentReference: 'INV-SSC-202603-01',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 3,
      periodYear: 2026
    },
    {
      id: 'tax-02',
      taxNumber: '010.000-26.88392100',
      taxDate: '2026-03-02',
      businessUnit: 'KIDINEST',
      taxCategory: 'PPN_INPUT',
      counterpartyName: 'PT Kimia Farma Trading & Distribution',
      counterpartyNPWP: '01.001.605.3-051.000',
      objectCode: 'PPN-MASUKAN-11',
      grossAmount: 15540000,
      dppAmount: 14000000,
      taxRate: 0.11,
      taxAmount: 1540000,
      documentReference: 'PO-KF-20260301-44',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 3,
      periodYear: 2026
    },
    {
      id: 'tax-03',
      taxNumber: 'BUPOT-23-202603-01',
      taxDate: '2026-03-03',
      businessUnit: 'ALL',
      taxCategory: 'PPH_23',
      counterpartyName: 'PT Lasertech Medika Utama (Maintenance Mesin Laser)',
      counterpartyNPWP: '02.441.980.1-013.000',
      objectCode: '24-104-01 (Jasa Perawatan Alat Medik)',
      grossAmount: 12000000,
      dppAmount: 12000000,
      taxRate: 0.02,
      taxAmount: 240000,
      documentReference: 'INV-LSR-0081',
      reconciliationStatus: 'RECONCILED',
      isReported: false,
      periodMonth: 3,
      periodYear: 2026
    }
  ])

  // Filtered doctor tax records
  const filteredDoctorTax = computed(() => {
    return doctorTaxRecords.value.filter(d => {
      if (accountingStore.businessUnit !== 'ALL' && d.businessUnit !== accountingStore.businessUnit) {
        return false
      }
      return true
    })
  })

  // Action: Export DJP e-Bupot CSV Format
  function exportEbupotCSV() {
    const csvData = doctorTaxRecords.value.map(d => ({
      'KODE_OBJEK_PAJAK': '21-100-07',
      'NAMA_PENERIMA_PENGHASILAN': d.doctorName,
      'NPWP': d.npwp,
      'NIK': d.nik,
      'PENGHASILAN_BRUTO': d.grossIncomeCurrentMonth,
      'TARIF': (d.taxRate * 100) + '%',
      'DPP': d.dppAmount,
      'PPH_DIPOTONG': d.pph21Amount,
      'NO_BUKTI_POTONG': d.ebupotNumber,
      'TGL_BUKTI_POTONG': '31/03/2026'
    }))

    exportToCSV(csvData, `DJP_eBupot_21_${accountingStore.fiscalYear}_0${accountingStore.fiscalMonth}`)
    accountingStore.showToast('File e-Bupot PPh 21 format DJP berhasil diunduh!', 'success')
  }

  return {
    doctorTaxRecords,
    filteredDoctorTax,
    ppnSummary,
    taxTransactions,
    exportEbupotCSV
  }
})
