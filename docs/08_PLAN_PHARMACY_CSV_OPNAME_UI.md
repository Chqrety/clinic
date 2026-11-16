# 08. Rencana Antarmuka Ekspor & Impor CSV Stok Opnam Farmasi

## 1. Tujuan & Latar Belakang
Backend sistem sudah memiliki modul `ExportStockOpnameCSV` dan `ImportStockOpnameCSV` di [pharmacy_usecase.go](file:///home/alpi/Projects/clinic/backend/internal/pharmacy/usecase/pharmacy_usecase.go) beserta unit test yang lolos 100%. 

Plan ini bertujuan untuk menghubungkan modul tersebut ke antarmuka pengguna (**UI**) di [PharmacyRoom.vue](file:///home/alpi/Projects/clinic/frontend/src/views/PharmacyRoom.vue) dan [InventoryView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/InventoryView.vue) agar apoteker dapat:
1. Mengunduh spreadsheet CSV saldo fisik obat terkini untuk audit lemari farmasi.
2. Mengunggah kembali file CSV hasil hitung opnam untuk memperbarui stok fisik secara massal dalam 1 kali klik.

---

## 2. Invarian Domain & Keamanan Data Obat
1. **Preservasi Reserved Stock:** Stok obat yang sedang dipesan dokter (*reserved*) tidak boleh terhapus secara sepihak saat stok fisik di-update via CSV.
   $$\text{available\_stock} = \text{new\_physical\_stock} - \text{existing\_reserved\_stock}$$
2. **Validasi Format CSV Ketat:** Sistem menolak unggahan file jika kolom wajib (ID Obat, Nomor Batch, Jumlah Fisik Baru) tidak lengkap atau bernilai negatif.
3. **Pencatatan Riwayat Mutasi:** Setiap penyesuaian stok opnam mencatat log mutasi bersumber "STOK_OPNAME" untuk keperluan audit internal.

---

## 3. Langkah Implementasi Teknis

### Tahap 1: Client API di Frontend (`frontend/src/api/client.ts`)
* Tambahkan method pemanggil backend:
  1. `exportStockOpname(clinicUnit: string): Promise<Blob>`:
     - Request ke `GET /api/pharmacy/opname/export?clinic_unit={unit}` dengan header auth staf.
     - Mengembalikan blob file `.csv`.
  2. `importStockOpname(clinicUnit: string, file: File): Promise<ImportOpnameResult>`:
     - Request `POST /api/pharmacy/opname/import` dengan form-data multipart (`file`).
     - Mengembalikan ringkasan: jumlah item berhasil diupdate, jumlah item gagal/dilewati.

### Tahap 2: Komponen Aksi di Header Apotek (`PharmacyRoom.vue` & `InventoryView.vue`)
* Tambahkan tombol aksi di sebelah filter unit:
  * **Tombol "Unduh Template Opnam":** Memanggil `exportStockOpname()` dan memicu unduhan file browser otomatis dengan nama `stok_opname_{unit}_{YYYY-MM-DD}.csv`.
  * **Tombol "Unggah Hasil Opnam":** Membuka modal dialog unggah CSV.

### Tahap 3: Modal Dialog Unggah & Preview Audit
* Desain modal sederhana & bersih (Nordic Modern / Brandlos):
  1. Area dropzone / file picker yang hanya menerima format `.csv`.
  2. Tombol "Mulai Sinkronisasi Stok".
  3. Indikator loading spinner saat pemrosesan batch data.
  4. Kartu hasil eksekusi:
     - Badge hijau: Berhasil memperbarui $N$ jenis obat.
     - Tombol "Tutup & Muat Ulang Tabel".
  5. Panggil kembali `loadInventory()` seketika agar angka stok di layar apotek langsung ter-update.

---

## 4. Rencana Pengujian (Verification)
1. **Pengujian Ekspor:**
   - Klik tombol unduh -> verifikasi file CSV terunduh dengan header: `medicine_id, medicine_name, batch_number, expired_date, system_stock, physical_count, notes`.
2. **Pengujian Impor:**
   - Edit nilai `physical_count` pada salah satu obat di file CSV.
   - Unggah kembali file -> verifikasi stok fisik di database dan tabel UI ter-update sesuai angka baru.
   - Verifikasi obat berstatus *reserved* tetap mempertahankan jumlah reservasi yang benar.
