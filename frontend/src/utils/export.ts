import * as XLSX from 'xlsx'

/**
 * Exports JSON array of objects to an Excel (.xlsx) file
 */
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

/**
 * Exports simple array of objects to CSV file
 */
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

/**
 * Triggers browser print dialog for current view
 */
export function printReport() {
  window.print()
}
