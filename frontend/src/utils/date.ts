import dayjs from 'dayjs'
import 'dayjs/locale/id'

dayjs.locale('id')

/**
 * Formats date to DD/MM/YYYY
 */
export function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY')
}

/**
 * Formats date and time to DD/MM/YYYY HH:mm
 */
export function formatDateTime(date: string | Date | null | undefined): string {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY HH:mm')
}

/**
 * Returns human readable Indonesian date: "15 Maret 2026"
 */
export function formatDateLong(date: string | Date | null | undefined): string {
  if (!date) return '-'
  return dayjs(date).format('D MMMM YYYY')
}

/**
 * Returns month & year in Indonesian: "Maret 2026"
 */
export function getPeriodLabel(year: number, month: number): string {
  const d = dayjs().year(year).month(month - 1)
  return d.format('MMMM YYYY')
}

/**
 * Checks if a date string is valid
 */
export function isValidDate(d: string): boolean {
  return dayjs(d).isValid()
}
