# STANDAR ARSITEKTUR BACKEND, DEPENDENSI LIBRARY & KEBIJAKAN MATA UANG IDR
## Sistem ERP Akuntansi Multi-Unit: Kidinest (Poli Anak) & Semarang Skin Center (SSC)

---

### 1. Daftar Library & Dependensi Backend (Golang)

Backend dibangun menggunakan **Go (Golang)** murni dengan pendekatan *Idiomatic Go* dan prinsip **Domain-Driven Design (DDD) Hexagonal Architecture**. Seluruh dependensi dipilih secara ketat untuk menjamin keandalan transaksi finansial tanpa *magic reflection* atau ORM yang menyembunyikan query.

| Library / Package | Modul / Path | Versi | Peran & Justifikasi Arsitektur |
| :--- | :--- | :--- | :--- |
| **HTTP Router** | `github.com/go-chi/chi/v5` | `v5.5.0` | Router HTTP ringan, composable, dan 100% kompatibel dengan `net/http` standar Go. Digunakan untuk memisahkan route per sub-domain akuntansi secara bersih. |
| **CORS Middleware** | `github.com/go-chi/cors` | `v1.2.2` | Menangani *Cross-Origin Resource Sharing* untuk integrasi aman antara frontend web kasir/backoffice dan backend API. |
| **UUID Generator** | `github.com/google/uuid` | `v1.6.0` | Pembangkitan UUIDv4 sebagai Primary Key di seluruh entitas database untuk mendukung desentralisasi ID dan mencegah *enumeration attack*. |
| **PostgreSQL Driver** | `github.com/jackc/pgx/v5` | `v5.12.0` | Driver native PostgreSQL dengan performa tinggi, mendukung *connection pooling* (`puddle`), *transaction savepoints*, dan kepatuhan penuh tipe data SQL `NUMERIC`. |
| **Fixed-Point Decimal** | `github.com/shopspring/decimal` | `v1.6.0` | **Wajib untuk seluruh kalkulasi moneter**. Mencegah galat pembulatan biner *floating-point* (`float64`/`float32`). |
| **SQLite Driver (Testing)** | `modernc.org/sqlite` | `v1.61.1` | Driver SQLite murni (CGO-free) untuk menjalankan unit test dan skenario integrasi akuntansi hermetis dan super cepat. |

---

### 2. Standar Mata Uang: Indonesian Rupiah (IDR)

Sistem akuntansi klinik Kidinest & Semarang Skin Center secara ketat menetapkan **Indonesian Rupiah (IDR / Rp)** sebagai mata uang fungsional (*functional currency*) dan mata uang pelaporan (*reporting currency*).

#### 2.1 Konvensi Penamaan & Simbol
- **ISO 4217 Currency Code**: `IDR`
- **Simbol Resmi**: `Rp` (dengan satu spasi pemisah sebelum nominal, misal: `Rp 150.000,00`)
- **Tanda Pemisah Ribuan**: Titik (`.`)
- **Tanda Pemisah Desimal / Sen**: Koma (`,`)
- **Presisi Database**: `NUMERIC(18, 2)` (18 digit total dengan 2 angka di belakang koma untuk akurasi perhitungan pajak PPh 21, diskon, dan tarif pembagian hasil dokter).
- **Format JSON API**: Dikirimkan sebagai *string representation* dari desimal tetap untuk mencegah parser JSON di browser/frontend mengonversinya menjadi IEEE-754 float.

#### 2.2 Aturan Pembulatan (Rounding Policy)
- **Metode Pembulatan**: *Round Half Up* (pembulatan ke atas jika digit berikutnya $\ge 5$).
- **Perhitungan Kas Fisik (Cashier Shift)**: Kas kasir fisik tidak menggunakan pecahan sen. Pembayaran tunai dibulatkan ke kelipatan terdekat sesuai pecahan uang kartal Bank Indonesia (Rp 500, Rp 1.000, Rp 2.000, dst).
- **Perhitungan Pajak & Jasa Medis**: Pemotongan PPh 21 Dokter dan PPN tetap mempertahankan 2 digit desimal pada tingkat jurnal umum untuk mencegah selisih akumulatif pada Laporan Keuangan.

#### 2.3 Utilitas Paket `pkg/decimal`
Paket `clinic-api/pkg/decimal` menyediakan fungsi standar untuk memanipulasi dan memformat nilai uang IDR:

```go
package decimal

// Konstanta Mata Uang IDR
const (
    CurrencyCode   = "IDR"
    CurrencySymbol = "Rp"
    CurrencyName   = "Indonesian Rupiah"
)

// Contoh Penggunaan:
// d := decimal.FromInt(1500000)
// decimal.FormatIDR(d)         -> "Rp 1.500.000,00"
// decimal.FormatRupiahNoSen(d) -> "Rp 1.500.000"
```

| Fungsi | Input | Output Contoh | Kegunaan |
| :--- | :--- | :--- | :--- |
| `FromInt(int64)` | `150000` | `Decimal(150000.00)` | Konstruksi dari integer |
| `FromFloat(float64)` | `25000.50` | `Decimal(25000.50)` | Konstruksi aman dari float via string |
| `FromString(string)` | `"120000000.00"` | `Decimal(120000000.00)` | Parsing dari payload JSON/Database |
| `RoundCurrency(Decimal)` | `Decimal(100.556)` | `Decimal(100.56)` | Pembulatan standar 2 desimal |
| `IsBalanced(debit, credit)` | `debit, credit` | `bool` | Validasi keseimbangan $\sum \text{Dr} = \sum \text{Cr}$ |
| `FormatIDR(Decimal)` | `Decimal(1500000)` | `"Rp 1.500.000,00"` | Tampilan resmi kwitansi & laporan |
| `FormatRupiahNoSen(Decimal)` | `Decimal(1500000)` | `"Rp 1.500.000"` | Tampilan kasir POS & display ringkas |

---

### 3. Struktur Modul & Pemisahan Tanggung Jawab (Hexagonal Architecture)

```
clinic/backend/
├── cmd/
│   ├── api/main.go            # Entry point REST API server
│   ├── worker/main.go         # Asynchronous worker runner
│   └── migrate/main.go        # Database migration CLI runner
├── internal/
│   ├── domain/                # Inti Logika Bisnis & Entity (Tanpa Framework)
│   │   ├── accounting/        # COA, Journal Entry, Fiscal Period, General Ledger
│   │   ├── clinicalunit/      # Profil Kidinest vs Semarang Skin Center
│   │   ├── billing/           # Invoicing, POS Line Items, Payment Records
│   │   ├── cashbank/          # Akun Kas, Rekening Bank BCA, BKM, BKK, Transfer
│   │   ├── doctorshare/       # Aturan Bagi Hasil, Settlement, Formula PMK 168
│   │   ├── treasury/          # Cashier Shift, Denominasi Kas Fisik, Rekonsiliasi
│   │   ├── tax/               # Withholding, PPh 21, PPh 23, e-Bupot CSV
│   │   ├── asset/             # Fixed Asset, Depresiasi Garis Lurus
│   │   ├── prepaid/           # Biaya Dibayar di Muka (Sewa Gedung), Amortisasi
│   │   ├── budget/            # Anggaran Cost Center, Analisis Varians
│   │   ├── audit/             # Jejak Audit Finansial
│   │   ├── receivable/        # Piutang Pasien & Klaim Asuransi, Aging
│   │   ├── payable/           # Hutang Supplier PBF/Skincare, Aging
│   │   └── shared/            # BaseEntity, BranchScoped, AuditFields
│   ├── application/           # Orkestrator Use Case & Business Workflows
│   │   ├── accounting/        # Journal posting, reversal, year-end closing
│   │   ├── billing/           # Atomic 3-cluster invoice settlement
│   │   ├── doctorshare/       # Approval, settlement payroll, JKK auto-post
│   │   ├── treasury/          # Open shift, close with count, overage/shortage
│   │   ├── prepaid/           # Capitalization, monthly amortization schedules
│   │   ├── budget/            # Budget creation, approval, variance report
│   │   ├── tax/               # Tax withholding recording, DJP e-Bupot export
│   │   ├── asset/             # Asset registration, batch straight-line depreciation
│   │   └── reporting/         # GL, Trial Balance, Income Statement, Balance Sheet, Cash Flow
│   ├── adapters/              # Secondary Driving Adapters
│   │   ├── persistence/postgres/ # PostgreSQL implementations via UnitOfWork
│   │   └── payment/           # Payment gateway processors (Cash, BCA, QRIS)
│   ├── transport/http/        # Primary Driving Transport
│   │   ├── dto/               # Data Transfer Objects & JSON request/response
│   │   ├── handler/           # HTTP controllers
│   │   ├── routes/            # Rute per sub-domain terpisah
│   │   └── middleware/        # BranchID context, CORS, logging, recovery
│   └── bootstrap/
│       ├── app.go             # Application lifecycles
│       └── container.go       # Pure manual dependency injection
└── pkg/
    ├── decimal/               # Standard fixed-point decimal & IDR currency math
    ├── database/              # DB connection factory
    └── response/              # Standard JSON API envelope (OK, Created, BadRequest)
```

---

### 4. Kebijakan Anti-Floating Point & Integritas Saldo

1. **Dilarang Keras Menggunakan Tipe `float32` atau `float64` untuk Uang**:
   Seluruh nilai uang wajib menggunakan `decimal.Decimal` dari paket `clinic-api/pkg/decimal`.
2. **Saldo Akun Bersifat Imutabel & Turunan (Derived Balance)**:
   Tidak ada kolom `current_balance` pada tabel akun yang dapat dimutasi secara bebas melalui query `UPDATE`. Seluruh saldo akun (Buku Besar, Neraca Saldo, Laba Rugi, dan Neraca) dihitung secara matematis dari baris transaksi jurnal (`journal_lines`) yang berstatus **`POSTED`**.
3. **Double-Entry Bookkeeping**:
   Setiap transaksi wajib memenuhi syarat $\sum \text{Debit} = \sum \text{Credit} > 0$. Jika terdapat selisih sekecil 0,01 rupiah, transaksi dibatalkan (*rollback*) dan melempar `ErrUnbalancedJournal`.
4. **Pemisahan Unit Usaha (Multi-Unit Segment)**:
   Setiap baris jurnal mencantumkan `cost_center_id` (`UNIT-KIDINEST`, `UNIT-SSC`, `SHARED-BACKOFFICE`) sehingga performa keuangan klinik anak dan klinik estetika dapat dilaporkan secara terpisah maupun terkonsolidasi.

