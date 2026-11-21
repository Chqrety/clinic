// Indonesian number to words (Terbilang) converter
export function terbilang(n: number): string {
  if (isNaN(n) || n === 0) return 'Nol Rupiah'
  
  const satuan = ['', 'Satu', 'Dua', 'Tiga', 'Empat', 'Lima', 'Enam', 'Tujuh', 'Delapan', 'Sembilan', 'Sepuluh', 'Sebelas']
  
  function convert(num: number): string {
    num = Math.floor(Math.abs(num))
    if (num < 12) return satuan[num]
    if (num < 20) return convert(num - 10) + ' Belas'
    if (num < 100) return convert(Math.floor(num / 10)) + ' Puluh ' + convert(num % 10)
    if (num < 200) return 'Seratus ' + convert(num - 100)
    if (num < 1000) return convert(Math.floor(num / 100)) + ' Ratus ' + convert(num % 100)
    if (num < 2000) return 'Seribu ' + convert(num - 1000)
    if (num < 1000000) return convert(Math.floor(num / 1000)) + ' Ribu ' + convert(num % 1000)
    if (num < 1000000000) return convert(Math.floor(num / 1000000)) + ' Juta ' + convert(num % 1000000)
    if (num < 1000000000000) return convert(Math.floor(num / 1000000000)) + ' Miliar ' + convert(num % 1000000000)
    return convert(Math.floor(num / 1000000000000)) + ' Triliun ' + convert(num % 1000000000000)
  }

  const result = convert(n).replace(/\s+/g, ' ').trim()
  return (n < 0 ? 'Minus ' : '') + result + ' Rupiah'
}
