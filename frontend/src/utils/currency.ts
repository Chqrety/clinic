// Indonesian Rupiah & Financial Number Utilities

/**
 * Formats a number to Indonesian Rupiah currency string
 * e.g. 15000000 -> "Rp 15.000.000"
 */
export function formatIDR(amount: number | null | undefined, includeSymbol = true): string {
  if (amount === null || amount === undefined || isNaN(amount)) {
    return includeSymbol ? 'Rp 0' : '0'
  }

  const isNegative = amount < 0
  const absAmount = Math.abs(Math.round(amount))
  const formatted = absAmount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.')

  if (includeSymbol) {
    return isNegative ? `(Rp ${formatted})` : `Rp ${formatted}`
  }
  return isNegative ? `(${formatted})` : formatted
}

/**
 * Formats a number to a compact Rupiah string for KPI cards and charts
 * e.g. 15000000 -> "Rp 15,0 Jt", 2500000000 -> "Rp 2,5 M"
 */
export function compactIDR(amount: number): string {
  if (amount === 0) return 'Rp 0'
  const isNegative = amount < 0
  const abs = Math.abs(amount)

  let formatted = ''
  if (abs >= 1_000_000_000) {
    formatted = (abs / 1_000_000_000).toFixed(1).replace('.', ',') + ' M'
  } else if (abs >= 1_000_000) {
    formatted = (abs / 1_000_000).toFixed(1).replace('.', ',') + ' Jt'
  } else if (abs >= 1_000) {
    formatted = (abs / 1_000).toFixed(0) + ' Rb'
  } else {
    formatted = abs.toString()
  }

  return isNegative ? `(Rp ${formatted})` : `Rp ${formatted}`
}

/**
 * Parses an IDR string or formatted number back to standard numeric value
 */
export function parseIDR(str: string): number {
  if (!str) return 0
  const clean = str.replace(/[^0-9,-]/g, '').replace(',', '.')
  const num = parseFloat(clean)
  return isNaN(num) ? 0 : num
}

/**
 * Calculates percentage with safety against division by zero
 */
export function calculatePercentage(part: number, total: number): number {
  if (!total || total === 0) return 0
  return Math.round((part / total) * 1000) / 10
}
