# 13. Audit & Rencana Penyelarasan UI/UX Sistem Akuntansi Klinik
**Klinik Utama Semarang Medical Group (Kidinest Pediatric Care & Semarang Skin Center)**  
*Dokumen Standar Desain Akuntansi mengacu pada [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md) & Ergonomi Workstation Medis (`#cashier`)*

---

## 1. Analisis Kritis: Kesenjangan Desain Akuntansi vs Workstation Medis

Logika akuntansi (*double-entry*, *COA*, *cashier settlement sync*) sudah berjalan dengan benar secara fungsional. Namun dari aspek **UI/UX**, modul akuntansi saat ini masih berstatus **"alien"** di dalam sistem klinik:

```
┌──────────────────────────────────────────────┬──────────────────────────────────────────────┐
│ WORKSTATION MEDIS KLINIK (STANDAR DESIGN.MD) │ MODUL AKUNTANSI SAAT INI (DEFECT)            │
├──────────────────────────────────────────────┼──────────────────────────────────────────────┤
│ • Palet Zinc Netral Monokrom (zinc-950/900)  │ • Didominasi warna biru Bootstrap/SaaS murahan│
│ • Hairline Dividers 1px, zero card-clutter   │ • Card-box tebal mengambang di mana-mana     │
│ • Minimalist Clinical Header terbuka        │ • Header dibungkus kotak card berbayang      │
│ • 1-Row Compact Hairline KPI Ticker         │ • Floating KPI Card warna-warni terpisah    │
│ • Brand resmi SMG (Semarang Medical Group)   │ • Label generik "CLINIC ERP Integrated v2.4" │
│ • Tabular-nums disiplin pada data finansial │ • Format angka belum seragam                 │
│ • Modal rounded-3xl border zinc-200/80       │ • Modal rounded-xl tombol biru               │
└──────────────────────────────────────────────┴──────────────────────────────────────────────┘
```

---

## 2. Temuan Kritis UI/UX & Solusi Berdasarkan DESIGN.md

### [1.0] Shell & Layout Global (Sidebar, Navbar, Kanvas Kerja)
* **Kritikal 1.1 — Disorientasi Brand & Palet Warna (Blue SaaS vs SMG Zinc Editorial):**
  * *Temuan:* `AppSidebar.vue` dan `AppNavbar.vue` memakai warna slate-900 dengan aksen `blue-600`, gradient logo biru-indigo, dan label "CLINIC ERP". Saat berpindah dari kasir/owner ke akuntansi, pengguna merasa berpindah ke aplikasi luar.
  * *Solusi:*
    1. Samakan sidebar akuntansi dengan `AdminSidebar.vue`: `bg-zinc-950 text-zinc-100 border-zinc-800`.
    2. Header identitas SMG monokrom (`Semarang Medical Group • Finance & Accounting`).
    3. Hapus kelas warna biru AI/Bootstrap (`blue-600`, `blue-500`, dsb.); ganti dengan standar interaksi klinik `bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl shadow-2xs`.
    4. Business Unit Selector memakai semantic brand badges resmi:
       * Kidinest: `bg-amber-50 text-amber-900 border-amber-200/80`
       * SSC: `bg-rose-50 text-rose-900 border-rose-200/80`
       * Konsolidasi: `bg-zinc-100 text-zinc-800 border-zinc-200`
* **Kritikal 1.2 — Wadah Kanvas Kerja (Canvas Container):**
  * *Temuan:* Halaman akuntansi memakai padding sembarangan `p-6` tanpa batas kanvas, membuat layout melar tak berbatas di monitor desktop lebar.
  * *Solusi:* Terapkan kanvas ergonomis DESIGN.md Section 5.A:
    `py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]`.

---

### [2.0] Struktur Header & KPI Ticker
* **Kritikal 2.1 — Sindrom Kotak Header (Card Everywhere Anti-Pattern):**
  * *Temuan:* Setiap halaman (COA, Jurnal, Kas-Bank, Pajak) membungkus judul dalam kotak card putih tebal (`p-6 rounded-2xl border shadow-2xs`) dengan badge biru besar. Ini membuang ruang vertikal 120px berharga.
  * *Solusi:* Terapkan Minimalist Clinical Header terbuka:
    * Breadcrumb tipis: `text-xs font-semibold text-zinc-500 uppercase tracking-wider` (contoh: `Akuntansi & Keuangan • Semarang Medical Group`).
    * Judul H1 tegas: `text-2xl font-bold text-zinc-950 tracking-tight mt-0.5`.
    * Kanan atas: Duty badge staf yang login (`Petugas Bertugas`) + tombol `Sync` berputar.
* **Kritikal 2.2 — Visual Clutter pada KPI Card:**
  * *Temuan:* `DashboardView.vue` memajang 5 kartu KPI terpisah dengan icon kotak warna-warni (`KpiCard`).
  * *Solusi:* Ganti menjadi **1-Row Compact KPI Ticker Strip (Hairline Grid)** seperti di `#cashier`:
    `grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden`.

---

### [3.0] Bagan Akun (COA) & Saldo Awal
* **Kritikal 3.1 — Hierarki Visual COA Datar:**
  * *Temuan:* Tabel COA menampilkan akun Induk (Level 1–3) dan Akun Detail (Leaf Level 4) dengan baris yang hampir sama, sehingga mata akuntan cepat lelah membedakan induk akun dengan pos detail.
  * *Solusi:* 
    * Akun Induk (Level 1): Background `bg-zinc-100/70`, font bold `text-zinc-950`, tanpa tombol posting.
    * Akun Golongan (Level 2 & 3): Semi-bold dengan indentasi padding `pl-6` dan `pl-10`.
    * Akun Detail (Level 4 Leaf): Indentasi `pl-14`, nomor akun monospaced tipis, saldo `tabular-nums font-bold text-right`, aksi Edit & Nonaktifkan compact.
* **Kritikal 3.2 — Visibilitas Keseimbangan Saldo Awal:**
  * *Temuan:* Pada `OpeningBalanceView.vue`, indikator keseimbangan (Total Debit vs Kredit) berada di bawah tabel panjang, mengharuskan scroll jauh untuk mengecek.
  * *Solusi:* Pasang **Sticky Footer Bar** di bagian bawah viewport:
    * Total Debit vs Total Kredit langsung terlihat setiap saat.
    * Badge pill: `● SEIMBANG (Selisih: Rp 0)` (Emerald) atau `● TIDAK SEIMBANG` (Rose).
    * Tombol Simpan disable otomatis jika selisih tidak nol.

---

### [4.0] Entri Jurnal Umum (Spreadsheet Ergonomics)
* **Kritikal 4.1 — Kecepatan Input Jurnal Terlalu Lambat (Click-Heavy UX):**
  * *Temuan:* Input jurnal umum membutuhkan banyak klik mouse untuk memilih akun dan berpindah antar kolom.
  * *Solusi:*
    * Mode entri tabel cepat (*Direct Table Fast Entry*).
    * Kolom Akun didukung smart autocomplete.
    * Alur navigasi keyboard: tombol `Tab` mengalir dari Akun -> Keterangan -> Debit -> Kredit -> Baris Baru.
    * Auto-balancing baris kedua (mengisi nominal kredit penyeimbang otomatis).

---

### [5.0] Standar Dialog Modal & Tipografi Finansial
* **Kritikal 5.1 — Inkonsistensi Dialog Modal:**
  * *Temuan:* Modal Tambah/Edit Akun dan Detail Jurnal memakai radius `rounded-xl` / `rounded-lg` dengan tombol biru.
  * *Solusi:* Standar DESIGN.md Section 6.C:
    * Backdrop: `fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4`.
    * Container: `max-w-xl bg-white rounded-3xl border border-zinc-200/80 shadow-2xl overflow-hidden`.
    * Header: `p-5 border-b border-zinc-100 flex items-center justify-between`.
    * Footer: `p-4 border-t border-zinc-100 bg-zinc-50/70 flex items-center justify-end gap-2.5`.
* **Kritikal 5.2 — Disiplin Angka Finansial (DESIGN.md Rule 49):**
  * *Temuan:* Kolom angka rupiah masih ada yang proporsional atau rata kiri.
  * *Solusi:* **WAJIB tabular-nums** (Plus Jakarta Sans) dengan alignment `text-right` untuk seluruh nilai mata uang. Monospace (`font-mono`) hanya untuk nomor kode akun, nomor jurnal, dan referensi transaksi.

---

## 3. Rencana Eksekusi Bertahap (Numbered Action Plan)

| No | Tahap Pekerjaan | Target File | Status | Cakupan Output |
| :---: | :--- | :--- | :---: | :--- |
| **1** | **Refactor Shell & Global Navigation** | `frontend/src/components/layout/AppSidebar.vue`<br>`frontend/src/components/layout/AppNavbar.vue`<br>`frontend/src/layouts/AccountingLayout.vue` | **SELESAI** | Standar sidebar `zinc-950` monokrom SMG. Brand tag Kidinest (Amber) & SSC (Rose). Hilangkan semua warna biru Bootstrap. Standar kanvas `max-w-[1600px]`. |
| **2** | **Refactor Header & 1-Row KPI Ticker** | `frontend/src/views/accounting/DashboardView.vue` | **SELESAI** | Minimalist Clinical Header terbuka + 1-row hairline KPI ticker summary menggantikan floating cards. |
| **3** | **Hierarki Visual COA & Sticky Saldo Awal** | `frontend/src/views/accounting/ChartOfAccountsView.vue`<br>`frontend/src/views/accounting/OpeningBalanceView.vue` | **SELESAI** | Indentasi visual COA Level 1–4. Sticky bottom balance bar pada Saldo Awal. |
| **4** | **Ergonomi Spreadsheet Jurnal Umum** | `frontend/src/views/accounting/JournalListView.vue`<br>`JournalDetailModal.vue`<br>`ManualJournalModal.vue` | **SELESAI** | Direct fast-entry table dengan auto-balance credit, autocomplete akun, dan balance badge real-time. Modal standardized ke `rounded-3xl`. |
| **5** | **Polish Buku Kas-Bank, Pajak & Modal Dialog** | `frontend/src/views/accounting/CashBankBookView.vue`<br>`frontend/src/views/accounting/TaxView.vue`<br>Modal Components | **SELESAI** | Standardisasi modal dialog `rounded-3xl`, tombol aksi `bg-zinc-950`, 1-row compact KPI strip, dan format angka `tabular-nums font-bold text-right`. |

---
*Status Implementasi: Seluruh tahap 1 s.d. 5 telah selesai 100% dan terverifikasi bersih (`npm run build` 0 error, `go test ./...` pass).*
*Dokumen ini merupakan acuan resmi tahap penyelarasan antarmuka sistem akuntansi SCMS Semarang Medical Group.*
