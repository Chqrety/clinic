# Standar Backend Pagination, Filtering, Searching, & Sorting
## ERP Klinik (Akuntansi, Medis, Kas/Bank, Perpajakan, Inventaris & Audit)

Dokumen ini mendefinisikan standar teknis implementasi server-side pagination, multi-column search, whitelist sorting, authorization scoping, dan filtering yang berlaku seragam untuk seluruh modul backend.

---

## 1. Prinsip & Ketentuan Arsitektur

1. **Server-Side Exclusivity**: Seluruh endpoint list memproses pagination, search, dan filter langsung di PostgreSQL/SQLite via parameterized query (`$1, $2, ...` / `?`). Dilarang mengambil seluruh record ke frontend atau layer transport untuk pemotongan memori.
2. **Deterministic Tie-Breaker**: Sorting multi-kolom selalu ditambahkan tie-breaker unik (misalnya `id DESC` atau `id ASC`) untuk menjamin konsistensi urutan pagination saat entri memiliki timestamp atau nominal identik.
3. **Pencarian Aman (Safe Search)**: Input `search` dibatasi pada kolom-kolom tekstual relevan (misalnya nomor faktur, nama pasien/vendor/dokter, NIK, NPWP) menggunakan `LOWER(...) LIKE $X` yang didukung indeks database B-Tree.
4. **Whitelist Sorting**: Kolom pengurutan (`sort_by`) dan arah urutan (`sort_order`) divalidasi ketat terhadap dictionary kolom terdaftar per endpoint. Upaya injeksi atau kolom tidak terdaftar akan otomatis ditolak dengan status HTTP 400 Bad Request.
5. **Otorisasi Terpusat**: Parameter `branch_id` dan `business_unit_id` divalidasi silang terhadap `UserSession` token pengguna yang login (`middleware.ValidateUnitAuthorization`). Permintaan lintas cabang tanpa izin SuperAdmin akan ditolak dengan status HTTP 403 Forbidden.
6. **Konsistensi Total Count**: Query `COUNT(*)` dieksekusi dengan klausa `WHERE` yang 100% identik dengan query data utama sebelum aplikasi `LIMIT` dan `OFFSET`.
7. **Indeks Database Terverifikasi**: Indeks komposit B-Tree diterapkan pada kolom filter dan sort utama (misalnya `(branch_id, transaction_date DESC, id DESC)`), diverifikasi menggunakan `EXPLAIN ANALYZE` dengan runtime `< 0.5 ms`.

---

## 2. Parameter Standar Query String

| Parameter | Tipe | Default | Batasan / Validasi | Keterangan |
| :--- | :--- | :--- | :--- | :--- |
| `page` | integer | `1` | Minimal `1` | Nomor halaman saat ini |
| `limit` | integer | `20` | Default `20`, Maksimal `100` | Jumlah baris per halaman |
| `search` | string | `""` | Maksimal 100 karakter | Kata kunci pencarian nama, kode, nomor referensi |
| `sort_by` | string | Variatif per modul | Whitelist per endpoint | Nama field pengurutan |
| `sort_order` | string | `DESC` atau `ASC` | `ASC` / `DESC` (case-insensitive) | Arah pengurutan |
| `start_date` | string | `nil` | Format ISO `YYYY-MM-DD` | Filter awal periode (00:00:00) |
| `end_date` | string | `nil` | Format ISO `YYYY-MM-DD` | Filter akhir periode (23:59:59) |
| `status` | string | `nil` | Enum status modul terkait | Filter status (misal `POSTED`, `PAID`, `ACTIVE`) |
| `branch_id` | UUID | Otomatis dari Session/Header | UUID v4 valid | Identitas cabang klinik |
| `business_unit_id` | UUID | `nil` | UUID v4 valid | Identitas unit/cost center (Kidinest, SSC, dll) |

---

## 3. Format Response Standar

Seluruh endpoint list mengembalikan envelope JSON standar:

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": [
    { "...": "..." }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_items": 142,
    "total_pages": 8
  }
}
```

---

## 4. Katalog Endpoint & Implementasi per Modul

### 4.1. Accounting & COA (`GET /api/v1/accounts`)
* **Sort Whitelist**: `code`, `name`, `type`, `created_at` (Default: `code ASC, id ASC`)
* **Search Fields**: `code`, `name`
* **Filter Khusus**: `type` (`ASSET`, `LIABILITY`, `EQUITY`, `REVENUE`, `EXPENSE`, dll), `is_active`

**Contoh Request:**
```http
GET /api/v1/accounts?page=1&limit=10&search=Kas&type=ASSET&sort_by=code&sort_order=asc HTTP/1.1
Host: localhost:8080
X-Branch-ID: 11111111-1111-1111-1111-111111111111
```

**Contoh Response:**
```json
{
  "success": true,
  "message": "Accounts retrieved successfully",
  "data": [
    {
      "id": "a1b2c3d4-0000-0000-0000-000000000001",
      "code": "1-1111",
      "name": "Kasir Utama Kidinest",
      "type": "ASSET",
      "normal_balance": "DEBIT",
      "is_active": true
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

---

### 4.2. Jurnal Umum & Posting (`GET /api/v1/journals`)
* **Sort Whitelist**: `transaction_date`, `journal_number`, `created_at` (Default: `transaction_date DESC, id DESC`)
* **Search Fields**: `journal_number`, `description`, `reference_id`
* **Filter Khusus**: `status` (`DRAFT`, `POSTED`, `REVERSED`), `start_date`, `end_date`, `branch_id`

**Contoh Request:**
```http
GET /api/v1/journals?page=1&limit=20&start_date=2026-01-01&end_date=2026-12-31&status=POSTED&sort_by=transaction_date&sort_order=desc HTTP/1.1
Host: localhost:8080
X-Branch-ID: 11111111-1111-1111-1111-111111111111
```

---

### 4.3. Kas & Bank
* **Akun Kas/Bank**: `GET /api/v1/cashbank/accounts`
  * **Sort Whitelist**: `code`, `name`, `type`, `created_at`
  * **Search**: `code`, `name`, `bank_name`, `account_number`
* **Mutasi Transaksi**: `GET /api/v1/cashbank/transactions`
  * **Sort Whitelist**: `transaction_date`, `transaction_number`, `amount`, `created_at` (Tie-breaker: `id DESC`)
  * **Search**: `transaction_number`, `description`, `reference_number`
  * **Filter Khusus**: `type` (`CASH_IN`, `CASH_OUT`, `TRANSFER`), `cash_bank_account_id`

---

### 4.4. Billing & Tagihan Faktur (`GET /api/v1/invoices`)
* **Sort Whitelist**: `created_at`, `invoice_number`, `patient_name`, `total_amount`, `paid_amount`, `payment_status`
* **Search Fields**: `invoice_number`, `patient_name`, `patient_nik`
* **Filter Khusus**: `status` (`UNPAID`, `PARTIAL`, `PAID`, `CANCELLED`), `type` (`REGULAR`, `PACKAGE`, `BPJS`), `start_date`, `end_date`

---

### 4.5. Piutang Pasien / Asuransi / Korporasi (`GET /api/v1/receivables`)
* **Sort Whitelist**: `due_date`, `issue_date`, `created_at`, `receivable_number`, `customer_name`, `total_amount`, `outstanding_amount`, `status`
* **Search Fields**: `customer_name`, `invoice_reference`, `receivable_number`
* **Filter Khusus**: `customer_type` (`PATIENT`, `INSURANCE`, `CORPORATE`, `BPJS`), `status` (`OPEN`, `PARTIAL`, `PAID`, `OVERDUE`)

---

### 4.6. Hutang Vendor & Dokter (`GET /api/v1/payables`)
* **Sort Whitelist**: `due_date`, `issue_date`, `created_at`, `payable_number`, `vendor_name`, `total_amount`, `outstanding_amount`, `status`
* **Search Fields**: `vendor_name`, `invoice_reference`, `payable_number`
* **Filter Khusus**: `vendor_type` (`SUPPLIER_MEDICINE`, `SUPPLIER_EQUIPMENT`, `DOCTOR`, `OPERATIONAL`), `status`

---

### 4.7. Jasa Medis Dokter (Doctor Share)
* **Aturan Bagi Hasil**: `GET /api/v1/doctorshare/rules`
  * **Sort Whitelist**: `doctor_name`, `service_category`, `sharing_percentage`, `fixed_fee`, `is_active`, `created_at`
  * **Search**: `doctor_name`, `service_category`
* **Settlement Pembayaran Fee Dokter**: `GET /api/v1/doctorshare/settlements`
  * **Sort Whitelist**: `period_year`, `period_month`, `doctor_name`, `gross_fee`, `net_fee`, `status`, `created_at`
  * **Search**: `doctor_name`, `doctor_npwp`
  * **Filter**: `year`, `month`, `status` (`DRAFT`, `APPROVED`, `PAID`), `doctor_id`

---

### 4.8. Treasury & Shift Kasir (`GET /api/v1/treasury/shifts`)
* **Sort Whitelist**: `start_time`, `end_time`, `cashier_name`, `opening_cash`, `actual_cash`, `status`, `created_at`
* **Search Fields**: `cashier_name`, `cash_account_code`, `notes`
* **Filter Khusus**: `status` (`OPEN`, `CLOSED`, `DEPOSITED`), `cash_account_code`, `cost_center_id`

---

### 4.9. Perpajakan / Bukti Potong (`GET /api/v1/tax/withholdings`)
* **Sort Whitelist**: `created_at`, `taxpayer_name`, `gross_amount`, `tax_amount`, `period_year`, `period_month`, `bupot_number`
* **Search Fields**: `taxpayer_name`, `taxpayer_nik`, `taxpayer_npwp`, `bupot_number`
* **Filter Khusus**: `year`, `month`, `tax_type` (`PPH21_DOCTOR`, `PPH23`, `PPH4_2`), `is_reported`

---

### 4.10. Aset Tetap & Amortisasi Biaya Dibayar Dimuka
* **Fixed Assets**: `GET /api/v1/assets`
  * **Sort Whitelist**: `asset_code`, `asset_name`, `category`, `acquisition_date`, `acquisition_cost`, `status`, `created_at`
  * **Search**: `asset_code`, `asset_name`, `category`
* **Prepaid Expenses**: `GET /api/v1/prepaid`
  * **Sort Whitelist**: `name`, `category`, `total_amount`, `remaining_amount`, `start_date`, `end_date`, `status`, `created_at`
  * **Search**: `name`, `prepaid_account_code`, `expense_account_code`

---

### 4.11. Anggaran / Budgeting (`GET /api/v1/budgets`)
* **Sort Whitelist**: `period_year`, `name`, `status`, `created_at`
* **Search Fields**: `name`
* **Filter Khusus**: `year`, `status` (`DRAFT`, `APPROVED`, `REVISED`), `cost_center_id`

---

### 4.12. Audit Log Finansial (`GET /api/v1/audit/logs`)
* **Sort Whitelist**: `created_at`, `action`, `entity_type`, `entity_id`
* **Search Fields**: `entity_id`, `ip_address`, `user_agent`
* **Filter Khusus**: `action`, `entity_type`, `start_date`, `end_date`

---

### 4.13. Farmasi & Antrean
* **Obat Apotek**: `GET /api/v1/drugs`
  * **Sort Whitelist**: `name`, `category`, `sell_price`, `created_at`
  * **Search**: `name`, `kfa_code`
* **Antrean Pasien**: `GET /api/v1/queues`
  * **Sort Whitelist**: `id`, `queue_number`, `patient_name`, `queue_status`, `created_at`
  * **Search**: `name`, `nik`, `queue_number`

