# Analisis & Rencana Redesign Frontend: Meja Depan & Kasir (GRO Workstation)
## Transformasi dari "Card-Clutter" Menuju Antarmuka *Nordic Clinical Minimalist*

Dokumen ini adalah hasil audit mendalam terhadap antarmuka peran **GRO / Kasir** ([CashierView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/CashierView.vue)), mengidentifikasi penumpukan elemen visual yang berlebihan (*UI bloat*), dan menyusun arsitektur desain baru yang tenang, ringkas, dan berkepadatan tinggi (*high-density efficiency*) untuk operasional harian klinik.

---

## 1. Masalah Utama: Sindrom "Card-itis" & Visual Fatigue

Antarmuka GRO saat ini secara fungsional lengkap, namun secara visual mengalami **sindrom umum template AI ("AI Slop")**:

```text
MASALAH STRUKTUR SAAT INI (KOTAK BERANAK):
┌─────────────────────────────────────────────────────────────────┐
│ Container Card (rounded-3xl border shadow-sm p-8)               │
│ ┌─────────────────────────────────────────────────────────────┐ │
│ │ Sub-Card Grid (rounded-3xl border p-5 shadow-sm)            │ │
│ │ ┌─────────────────────────────────────────────────────────┐ │ │
│ │ │ Inner Card Status (bg-white/90 rounded-2xl border p-3)   │ │ │
│ │ │ ┌─────────────────────────────────────────────────────┐ │ │ │
│ │ │ │ Pill Badges (rounded-full bg-zinc-100 font-mono)    │ │ │ │
│ │ │ └─────────────────────────────────────────────────────┘ │ │ │
│ │ └─────────────────────────────────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Empat Dosa Visual di Antarmuka Saat Ini:

1. **Kotak di dalam Kotak (*Nested Cards Everywhere*):**  
   Setiap modul, daftar, baris, bahkan badge dibungkus dalam sudut melengkung ekstrem (`rounded-3xl`, `rounded-2xl`) dengan bayangan dan border ganda. Ini memboroskan ruang layar hingga 35% hanya untuk *margin* dan *padding*.
2. **Teks Edukasi Berlebihan (*Instructional Slop*):**  
   Ada banyak paragraf penjelasan panjang seperti:  
   > *"Untuk pasien lansia, anak darurat, atau yang datang langsung tanpa smartphone. Nomor tiket antrian otomatis terkirim ke WhatsApp."*  
   > *"Pasien booking online tidak akan dipanggil ke ruang dokter sebelum melapor ke meja GRO dan di-check in."*  
   Staf klinik bekerja 7–8 jam sehari; mereka **tidak butuh membaca instruksi tutorial berulang kali**. Teks ini hanya menambah kebisingan visual (*cognitive load*).
3. **Inflasi Warna & Badge (*Badge Spam*):**  
   Dalam satu layar terdapat belasan pill badge warna-warni yang berdenyut (*pulsing dots*): amber, emerald, rose, zinc. Ketika semuanya mencolok, tidak ada lagi informasi penting yang menonjol.
4. **Alur Kerja Terpecah (*Tab Fragmentation*):**  
   Kasir harus berpindah tab antara `Tagihan Kasir (POS)` dan `Meja Depan (GRO)` untuk sekadar melihat apakah dokter sudah siap atau pasien reservasi sudah datang. Di jam sibuk, perpindahan tab ini memperlambat waktu pelayanan loket.

---

## 2. Filosofi Redesign: *Nordic Clinical Architecture*

Sesuai arahan [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md), meja depan klinik estetika dan spesialis anak harus memancarkan kesan **tenang, higienis, berwibawa, dan efisien**:

* **Garis Rambut 1px (*1px Hairline Architecture*):** Singkirkan kotak tebal bersudut bundar 24px. Gunakan garis pemisah tipis `border-zinc-200/80` yang bersih dan elegan (seperti lembar kerja medis profesional atau terminal modern).
* **Flat Surfaces & Warm Canvas:** Latar canvas off-white hangat (`#FAF7F2` untuk Kidinest / `#FAF5F5` untuk SSC) tanpa gradasi semen gelap.
* **Scanability dalam 1 Detik:** Data penting (Nomor Antrian, Nama Pasien, Status Dokter, Total Biaya) menggunakan tipografi tajam berbobot tegas, tanpa dekorasi berlebihan.
* **Micro-Copy Fungsional:** Ganti kalimat panjang dengan label ringkas 1–2 kata: `Reservasi Online`, `Ruang Dokter`, `Walk-In`.

---

## 3. Analisis Per Komponen: SEBELUM vs SESUDAH

### A. Indikator Ruang Dokter (*Doctor Room Monitor*)

| Aspek | Saat Ini (Old Card Grid) | Desain Baru Minimalis (New Status Strip) |
| :--- | :--- | :--- |
| **Bentuk** | 4 Card besar bertingkat memakan 350px tinggi layar. | **1 Baris Horizontal Strip (*Compact Matrix*)** setinggi 72px. |
| **Elemen** | Kotak dalam kotak + 3 badge + tombol besar. | 4 kolom setara dipisahkan garis 1px vertikal. |
| **Informasi** | Nama dokter, spesialisasi, status, pasien di dalam, pasien tunggu, tombol panggil. | • Dot hijau/kuning tenang.<br>• Nama dokter & nomor poli.<br>• Pasien aktif (misal: `SSC-001 (12m)`).<br>• Tombol panggil minimalis di samping nama pasien berikutnya. |
| **Dampak** | Menghabiskan setengah layar; staf harus scroll ke bawah. | Menghemat 70% ruang vertikal, langsung terlihat sekilas di paling atas. |

---

### B. Daftar Reservasi Online & Check-In

| Aspek | Saat Ini (Old Nested Table) | Desain Baru Minimalis (Interactive Stream) |
| :--- | :--- | :--- |
| **Bentuk** | Kotak putih tebal `rounded-3xl` dengan tabel lebar berjarak longgar. | **Daftar Bersih Bergaris (*Hairline Data Rows*)** menyatu dengan kanvas. |
| **Teks** | Paragraf peringatan panjang tentang kewajiban konfirmasi kehadiran fisik. | Hanya header ringkas: `Reservasi Hari Ini • 3 Belum Tiba` + search input garis bawah (*borderless bottom line*). |
| **Aksi** | Tombol hijau tebal mencolok di tiap baris. | Tombol pill `Check-In` ringkas yang berubah centang saat berhasil. |

---

### C. Alur Kasir / POS & Billing

| Aspek | Saat Ini (Old POS View) | Desain Baru Minimalis (Split Workstation) |
| :--- | :--- | :--- |
| **Tata Letak** | Kolom kiri (Daftar Bill) & Kolom kanan (Payment Form) penuh kotak terisolasi. | **Panel Berdampingan (*Master-Detail Split*)** dipisahkan 1 garis pemisah elegan. |
| **Kalkulator** | Banyak input rupiah dengan background abu-abu tebal. | Angka nominal besar bersih (`tabular-nums`), perhitungan kembalian otomatis instan. |
| **Struk / Invoice** | Pratinjau invoice berada di bawah atau modal tersembunyi. | *Slide-over slip* ringkas di sisi kanan saat transaksi selesai. |

---

## 4. Konsep Antarmuka Baru: *The Unified GRO Workspace*

Alih-alih memisahkan tab Kasir dan tab GRO secara kaku, layar baru dirancang dengan layout **3 Zona Terpadu**:

```text
┌────────────────────────────────────────────────────────────────────────┐
│ TOP BAR: Identitas Shift Staf • Jam • Switcher Unit • Tombol Tutup Shift │
├────────────────────────────────────────────────────────────────────────┤
│ ZONA 1 (70px): DOCK STATUS RUANG DOKTER (Compact 4-Column Strip)       │
│ [Poli 1: dr. Hendra • 🟡 B-001 (8m)]  [Poli 2: dr. Patricia • 🟢 SIAP]  │
│ [Poli 3: dr. Dian • 🟡 SSC-001 (14m)] [Poli 4: dr. Kevin • 🟢 SIAP]   │
├───────────────────────────────────┬────────────────────────────────────┤
│ ZONA 2 (KIRI): ANTREAN & RESERVASI│ ZONA 3 (KANAN): TRANSAKSI & DETAIL │
│                                   │                                    │
│ [Filter: Semua | Kasir | Booking] │ • Rincian Tagihan Pasien Terpilih  │
│                                   │ • Subtotal: Registrasi, Obat, Tind.│
│ Daftar 1: Siap Bayar di Kasir     │ • Diskon Voucher & Metode Bayar    │
│ • B-001 - Kenzo Alvaro            │                                    │
│ • SSC-001 - Nadia Saphira         │ [ Tombol Lunasi & Cetak Struk ]    │
│                                   │                                    │
│ Daftar 2: Reservasi Belum Tiba    │ ─── ATAU JIKA WALK-IN: ─────────── │
│ • S-019 - Rio Dewanto [Check-In]  │ Form cepat 3 field (Nama, WA, Poli)│
│ • S-020 - Marsha T.   [Check-In]  │ tanpa form bertingkat-tingkat.     │
└───────────────────────────────────┴────────────────────────────────────┘
```

### Keuntungan Arsitektur Ini:
1. **Zero-Tab Flipping:** GRO bisa melayani pembayaran kasir sambil melirik status dokter dan menyambut pasien reservasi yang baru datang dalam **satu layar tanpa ganti-ganti tab**.
2. **Tidak Ada Elemen Sia-Sia:** Semua border tebal, box shadow tebal, dan teks narasi dihapus.
3. **Efisiensi Gerakan Tangan:** Tombol Check-In, Panggil, dan Bayar ditempatkan di jalur gerak mouse yang natural.

---

## 5. Rencana Eksekusi Refactoring (Langkah Kerja Teknis)

Refactoring akan dilakukan langsung pada [CashierView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/CashierView.vue) dengan panduan langkah berikut:

### Langkah 1: Pembersihan Awal (*De-Slop*)
- Hapus semua paragraf tutorial panjang (`Untuk pasien lansia...`, `Pasien booking online tidak akan...`).
- Hilangkan pembungkus luar kartu `rounded-3xl border shadow-sm p-8` berganti ke kontainer kanvas minimalis bergaris rambut (`border-b border-zinc-200/80`).

### Langkah 2: Rekonstruksi Baris Status Dokter (*Doctor Room Strip*)
- Ubah 4 card besar menjadi satu *dock status horizontal* yang ringkas, menempati bagian atas halaman.
- Tampilkan nama dokter, status dot (hijau/kuning), pasien yang sedang di dalam, dan tombol mikro panggil jika dokter siap.

### Langkah 3: Integrasi Feed Antrian (Antrian Kasir + Reservasi Booking)
- Gabungkan daftar antrian kasir belum lunas dan reservasi online dalam satu panel kiri berkepadatan tinggi (*high-density list*) dengan segmented control pemilih:
  `[ Semua (6) ]` `[ Tagihan Kasir (4) ]` `[ Reservasi Belum Tiba (2) ]`.
- Baris reservasi memiliki tombol langsung `[ Check-In ]` tanpa harus masuk ke tab terpisah.

### Langkah 4: Simplifikasi Form Walk-In & Pembayaran POS
- Form Walk-in disederhanakan menjadi panel input cepat (Nama, No. WA, Poli, Dokter) yang bisa dibuka cepat tanpa memenuhi layar.
- Kalkulator pembayaran di panel kanan dibuat bersih dengan font monospace tegas untuk angka nominal rupiah.

---

*Hasil dari dokumen analisis ini siap dijadikan acuan langsung untuk mengeksekusi refactoring kode antarmuka GRO.*
