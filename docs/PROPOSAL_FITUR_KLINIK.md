# PROPOSAL FITUR PENGEMBANGAN SISTEM INFORMASI KLINIK
## Kidinest (Poli Anak & Vaksinasi) & Semarang Skin Center (SSC)

Proposal ini memetakan seluruh kebutuhan pengembangan sistem berdasarkan **Peran Pengguna / Bagian Operasional Klinik**. Setiap bagian dirinci ke dalam format tabel 2 kolom yang mudah dipahami oleh pihak manajemen klinik.

---

## Bagian 1: Portal Pasien & Pendaftaran Mandiri (Self-Service QR)
Bagian ini digunakan langsung oleh pasien untuk mendaftar secara mandiri melalui ponsel (scan QR / tautan web) tanpa perlu antri di meja resepsionis.

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **1.1. Landing Page Unit Kidinest** | Halaman pendaftaran khusus poli anak & vaksinasi. Tampilan minimalis dan fokus (hanya logo klinik dan form pendaftaran), langsung terhubung ke metode pembayaran online. |
| **1.2. Landing Page Unit Semarang Skin Center (SSC)** | Halaman pendaftaran khusus layanan estetika/facial. Menggunakan sistem reservasi jadwal (*booking* tanpa bayar di awal/Rp0). Pembayaran dilakukan di kasir setelah tindakan. |
| **1.3. Form Pendaftaran Mandiri (Privasi Penuh / Tanpa NIK)** | Menjaga privasi penuh pelanggan (*privacy-first*): pendaftaran tidak meminta Nomor KTP / NIK. Pasien cukup mengisi Nama Lengkap Pasien/Anak dan Nomor WhatsApp aktif untuk konfirmasi tiket antrian. |
| **1.4. Pemilihan Kategori Antrian Kidinest** | Penomoran antrian otomatis berdasarkan jenis kedatangan:<br>• **Kategori A:** Layanan Vaksinasi Anak.<br>• **Kategori B:** Pasien Sakit (Poli Umum Anak).<br>• **Kategori C:** Pemantauan Tumbuh Kembang.<br>• **Kategori D:** Paket Vaksinasi & Tumbuh Kembang.<br>• **Kategori E:** Pasien Datang Langsung / Walk-In. |
| **1.5. Pemilihan Dokter Jaga** | Pasien dapat memilih dokter spesialis anak atau dokter tindakan yang diinginkan sesuai jadwal jaga aktif (tersedia opsi untuk 2–3 dokter). |
| **1.6. Pembayaran Online Pendaftaran** | Pembayaran biaya registrasi Kidinest via QRIS / transfer otomatis. Setelah dana diverifikasi, nomor tiket antrian resmi langsung terbit seketika. |
| **1.7. Notifikasi Antrian Otomatis via WhatsApp** | Tiket antrian, nama dokter yang dituju, estimasi waktu, serta tautan langsung untuk memantau antrian otomatis terkirim ke WhatsApp pasien. |
| **1.8. Layar Pantau Antrian Real-Time (Live Queue)** | Pasien dapat memantau pergerakan antrian langsung dari layar HP miliknya atau melalui TV monitor ruang tunggu klinik. |
| **1.9. Perlindungan Privasi Pasien (Sensor Nama)** | Nama pasien di layar monitor ruang tunggu publik disamarkan secara otomatis (contoh: *Budi P\*\*\*\** atau *Siti A\*\*\*\*) demi kerahasiaan data medis. |
| **1.10. Tampilan Antrian Berikutnya (Next Queue)** | Menampilkan nomor antrian yang akan dipanggil berikutnya (*preview next queue*), sehingga pasien berikutnya dapat bersiap. |
| **1.11. Pemilahan Antrian Berdasarkan Dokter** | Memungkinkan pasien dan staf memfilter tampilan monitor antrian per masing-masing dokter yang sedang praktek. |

---

## Bagian 2: Meja Depan & Kasir (GRO / Frontdesk)
Bagian ini digunakan oleh petugas resepsionis dan kasir untuk melayani pembayaran, mencetak bukti transaksi, dan menangani pasien darurat.

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **2.1. Bantuan Reservasi Pasien Datang Langsung (Walk-In / Darurat)** | Pasien yang datang langsung atau kondisi darurat tetap diproses melalui sistem reservasi online dengan dibantu oleh staf GRO di meja depan. Bebas kertas (*paperless* / tanpa struk fisik) karena nomor antrian otomatis terkirim langsung ke WhatsApp pasien dan tampil di layar monitor antrian. |
| **2.2. Manajemen Pembayaran Multi-Metode** | Kasir dapat memproses transaksi melalui berbagai metode: Tunai (*Cash*), Mesin EDC, Transfer Bank, hingga pemotongan kupon/Voucher Diskon dengan nominal fleksibel. |
| **2.3. Cetak Invoice & Struk Transaksi Fisik** | Mencetak bukti pembayaran resmi: struk thermal kasir untuk pasien di tempat atau invoice format PDF ber-kop klinik untuk klaim asuransi mandiri. |
| **2.4. Fitur Pelaporan Kasir Harian ke Pimpinan** | Tombol tutup shift (*end-of-shift*). Kasir cukup menekan satu tombol untuk mengirim rekapitulasi setoran fisik (tunai, transfer, voucher) langsung ke pimpinan klinik. |
| **2.5. Pencatatan Pengembalian Dana (Refund Manual)** | Fitur resmi jika ada pasien yang membatalkan tindakan berbayar. Kasir mencatat nomor rekening tujuan pengembalian, status transaksi menjadi batal (*refunded*), dan alokasi obat dibatalkan. |

---

## Bagian 3: Ruang Dokter (Pemeriksaan & E-Resep)
Bagian ini digunakan oleh dokter di ruang praktek untuk memanggil pasien, mencatat rekam medis, dan menerbitkan resep elektronik.

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **3.1. Antrian Pasien Adil (Sistem FIFO Ketat)** | Daftar antrian di ruang dokter otomatis memprioritaskan pasien yang menyelesaikan pembayaran pendaftaran lebih awal (*First In, First Out*), mencegah saling serobot antrian. |
| **3.2. Filter Pasien Khusus Dokter Jaga** | Setiap dokter hanya melihat antrian pasien yang memilih dirinya sendiri, lengkap dengan indikator keluhan dan kategori antrian (A, B, C, D, E [Walk-In], atau tindakan SSC). |
| **3.3. Penulisan E-Resep Cepat (Tabel, Harga & Input Manual)** | Dokter menulis resep dalam format tabel interaktif: dilengkapi *Dropdown Search* instan (ketik 2 huruf nama obat muncul), indikator stok fisik apotek, tampilan harga satuan obat secara transparan, serta opsi **input manual diketik** (fleksibel untuk obat racikan khusus atau item non-database beserta harga/nominalnya). |
| **3.4. Input Komponen Tindakan Medis Kidinest** | Form rincian biaya medis anak: input jasa dokter, suntik vaksinasi, dan daftar obat racikan maupun obat paten. |
| **3.5. Input Komponen Tindakan Medis SSC** | Form rincian tindakan estetika: input jasa dokter/terapis (facial, chemical peeling, dll) dan paket produk skincare klinik. |
| **3.6. Notifikasi Otomatis Resep ke Apotek** | Begitu dokter menekan tombol "Kirim Resep", sistem otomatis membunyikan tanda suara/chime di komputer apotek agar racikan obat bisa segera dipersiapkan. |

---

## Bagian 4: Farmasi & Pengelolaan Gudang Obat (Apotek)
Bagian ini digunakan oleh staf farmasi untuk meracik obat, memantau ketersediaan stok, dan melakukan audit berkala (*stok opnam*).

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **4.1. Layar Antrian Resep Masuk Real-Time** | Petugas apotek menerima resep dari ruang dokter secara langsung dan memprosesnya dengan alur status: *Diterima*, *Diracik*, dan *Siap Diserahkan*. |
| **4.2. Penguncian & Pemotongan Stok Otomatis (Anti Selisih)** | Obat otomatis berstatus "Dipesan" (*Reserved*) saat dokter membuat resep, dan terpotong permanen saat obat diserahkan ke pasien. Menghilangkan risiko selisih hitung fisik. |
| **4.3. Master Data Obat & Manajemen Harga** | Database master obat klinik: pencatatan nama, satuan (tablet, botol, ampul, tube), harga dasar, harga jual, nomor batch, dan tanggal kadaluarsa (*expired date*). |
| **4.4. Fitur Ekspor Data Stok ke Excel/CSV** | Apoteker dapat mengunduh seluruh data saldo obat ke dalam file Excel hanya dengan 1 klik untuk persiapan audit fisik lemari obat (*stok opnam*). |
| **4.5. Fitur Impor Data Stok dari Excel/CSV** | Memperbarui jumlah stok fisik hasil hitung opnam secara massal dengan mengunggah kembali file Excel ke sistem tanpa perlu input manual satu per satu. |
| **4.6. Modul Migrasi & Cadangan Data Sistem Lama** | Fasilitas impor data riwayat pasien dan daftar obat dari aplikasi lama klinik ke dalam sistem baru agar riwayat klinik tidak hilang. |

---

## Bagian 5: Keuangan, Akuntansi & Perpajakan (Backoffice)
Bagian ini digunakan oleh bagian keuangan dan manajemen klinik untuk memantau pendapatan bersih dan menyusun pelaporan pajak.

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **5.1. Pemisahan Akun Pendapatan Otomatis** | Sistem secara otomatis memecah 1 tagihan pasien ke dalam pos akuntansi yang tepat: Akun Jasa Medis Dokter, Akun Penjualan Obat, Akun Vaksinasi, dan Akun Treatment Estetika. |
| **5.2. Laporan Rekapitulasi Pajak Transaksi** | Mencatat dan merekapitulasi beban pajak (PPN/PPh) atas penjualan obat dan jasa tindakan untuk memudahkan akuntan menyusun SPT pajak bulanan. |
| **5.3. Rekonsiliasi Tutup Buku Kasir Harian** | Memverifikasi kesesuaian antara uang fisik yang disetor kasir dengan data transaksi digital yang tercatat di sistem secara otomatis. |
| **5.4. Laporan Rekapitulasi Honor Jasa Dokter** | Menghitung otomatis pembagian bagi hasil / jasa medis dokter berdasarkan jumlah tindakan dan pasien yang telah selesai dilayani. |

---

## Bagian 6: Keamanan Data & Manajemen Staf
Bagian pengelolaan akses untuk memastikan privasi data pasien dan keamanan operasional klinik.

| Item Fitur | Keterangan & Manfaat Fitur bagi Klinik |
| :--- | :--- |
| **6.1. Login Staf Menggunakan Nomor Induk Karyawan** | Seluruh staf (kasir, dokter, perawat, apoteker) masuk ke sistem menggunakan Nomor Karyawan resmi klinik dan kata sandi. Tidak bergantung pada akun email Google pribadi staf. |
| **6.2. Pembatasan Hak Akses Berjenjang (Role-Based Access)** | Setiap staf hanya dapat mengakses modul kerjanya masing-masing: kasir tidak dapat melihat riwayat medis, dokter fokus pada pasien, apotek fokus pada resep, dan pimpinan memegang akses laporan penuh. |
| **6.3. Modul Pengelolaan Staf oleh Atasan (Super Admin / Pimpinan)** | Menu eksklusif bagi Atasan/Pimpinan untuk mendaftarkan dan mengelola akun staf: input data karyawan baru, penerbitan Nomor Karyawan resmi, penentuan peran/hak akses (kasir, dokter, farmasi), serta kontrol reset kata sandi staf. |

---

## Rencana Waktu & Tahapan Pelaksanaan (Timeline)

| Tanggal | Target Capaian Kerja |
| :--- | :--- |
| **16 November 2026** | **Kunjungan & Sesi Demo di Klinik:** Uji coba portal QR pasien (Kidinest & SSC), antrian real-time, layar kasir, cetak invoice, dan penulisan e-resep dokter. |
| **Pertengahan Desember 2026** | **Peluncuran Penuh (Go-Live):** Seluruh modul sistem (termasuk modul stok opnam Excel, integrasi akuntansi, dan pajak) aktif digunakan operasional klinik. |
