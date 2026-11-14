# 📘 DOKUMENTASI SISTEM LENGKAP: BACKEND & FRONTEND
## Sistem Informasi Manajemen Klinik & ERP Akuntansi Terpadu (Smart Clinic OS - BCS)
**PT. Brayo Cipta Sejahtera (Klinik Pratama Kidinest & Samarinda Skin Center)**

Dokumentasi teknis lengkap menyeluruh untuk seluruh sistem backend (Go + PostgreSQL) dan frontend (Vue 3 + Tailwind CSS + Pinia) telah disusun secara terperinci.

Silakan membaca dokumen lengkap pada:
👉 **[DOKUMENTASI SISTEM LENGKAP (docs/SYSTEM_DOCUMENTATION.md)](file:///Users/macbookpro/Documents/project-11/clinic/docs/SYSTEM_DOCUMENTATION.md)**

---

### Ringkasan Cakupan Dokumentasi:
1. **Ringkasan Eksekutif & Arsitektur Global:** Diagram Mermaid interaksi frontend, Vite proxy gateway, backend Go monolith, dan database PostgreSQL 14.
2. **Arsitektur Backend (Go Modular Monolith):** Struktur direktori, aplikasi services, clean architecture, dan multi-branch context (`X-Branch-ID`).
3. **Arsitektur Database & Skema PostgreSQL:** Relasi tabel ERD (`accounts`, `cash_bank_accounts`, `cash_bank_transactions`, `journal_entries`, `journal_lines`, `tax_withholdings`, `billing_invoices`).
4. **Bagan Akun Standar BCS (Chart of Accounts - 106 Akun):** Struktur 4-tingkat (Aset, Kewajiban, Ekuitas, Pendapatan, HPP/Beban Pokok, Pendapatan Non-Usaha, Beban Non-Usaha).
5. **Katalog Lengkap REST API Endpoints:** Dokumentasi detail endpoint `/api/v1/accounts`, `/api/v1/journals`, `/api/v1/cash`, `/api/v1/bank`, `/api/v1/taxes`, `/api/v1/reports`, `/api/v1/billing`, `/api/v1/receivables`, `/api/v1/payables`, `/api/v1/doctor-shares`, `/api/v1/treasury`, `/api/v1/periods`, dan `/api/v1/audit-logs`.
6. **Arsitektur Frontend (Vue 3 + Pinia + Tailwind CSS):** Struktur komponen, Pinia stores (`ledger.ts`, `cashbank.ts`, `journals.ts`, `tax.ts`, `reports.ts`), routing, dan typed API client.
7. **Modul Finansial & Aturan Bisnis Khusus:**
   - Pembuatan mutasi kas & bank mandiri (Direct In, Direct Out, Transfer kas/bank).
   - Pemisahan SPT Masa dan SPT Tahunan Badan 1771.
   - Penarikan otomatis kredit pajak SPT Tahunan dari data SPT Masa (PPh 23 & PPh 25).
   - Filter ketat basis pajak **khusus Bank saja (Bukan Kas Tunai)** sesuai regulasi DJP.
   - Pembuatan Saldo Awal Akuntansi (Opening Balance Engine), validator seimbang debit=kredit, dan 1-klik auto-balance ke ekuitas (`3999 Opening Balance Equity`).
   - Persistensi data murni ke database PostgreSQL (bebas data hardcoded).
8. **Panduan Operasional & Menjalankan Sistem (Native Local Run):** Cara menjalankan PostgreSQL, Go Backend, Vue Frontend dev & build secara native tanpa Docker.
9. **Jejak Perbaikan, Troubleshooting & Solusi Teknis:** Solusi FK constraint `branch_id`, pencegahan tabrakan nomor jurnal (`journal_number`), dan penegakan isolasi docs dari router UI.

