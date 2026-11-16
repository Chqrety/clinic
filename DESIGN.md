# Design System & UI/UX Guidelines: Clinical Editorial (Anti-AI Slop)

Dokumen ini adalah standar desain resmi antarmuka Smart Clinic Management System (SCMS) untuk dua brand terpadu: **Kidinest Pediatric Care** dan **Semarang Skin Center (SSC)** di bawah **Klinik Utama Semarang Medical Group**.

---

## 1. Prinsip Utama: Anti "AI Slop" (Bukan Teks Raksasa & Bukan Card di Mana-Mana)

Pola umum template buatan AI memiliki kelemahan fatal yang membuat website terasa murah dan tidak manusiawi:
* ❌ **Giant Screaming Headlines:** Teks heading 72–96px hitam tebal yang memenuhi setengah layar dan berteriak kepada pasien.
* ❌ **Cards Everywhere (Sindrom Kotak di Mana-Mana):** Membungkus setiap judul, fitur, dokter, dan layanan ke dalam kotak kartu tebal (`border rounded-3xl p-6 shadow`) yang berderet kaku 4 kolom.
* ❌ **Warna Palsu / Tech-Noir:** Latar abu-abu semen dingin, gradasi ungu neon, atau nuansa korporat yang menakutkan bagi pasien.

### Standar Baru: **Editorial Clinical Architecture**
Mengadopsi pendekatan situs medis dan estetika kelas dunia (*Aesop*, *Dr. Barbara Sturm*, klinik spesialis dermatologi Tokyo/Seoul, serta feed resmi Instagram Kidinest & SSC):
1. **Tata Letak Garis Rambut (*1px Hairline Dividers*):** Pemisah antar konten menggunakan garis 1px tipis elegan atau whitespace proporsional, bukan kotak kartu tebal.
2. **Tipografi Berbobot Manusiawi (*Human-Scale Typography*):** Heading berukuran proporsional (32–44px di desktop), dipadukan dengan aksen **Serif Italic** yang anggun dan berwibawa.
3. **Format Daftar Editorial (*Editorial Index Row*):** Layanan disajikan dalam bentuk daftar berbaris rapi dengan nomor tipis (*01, 02, 03, 04*), nama treatment, dokter penanggung jawab, dan tombol reservasi langsung.

---

## 2. Karakteristik Visual Dua Brand (Dual-Brand Identity)

| Parameter | **Kidinest Pediatric Care** (`/kidinest`) | **Semarang Skin Center / SSC** (`/ssc`) |
| :--- | :--- | :--- |
| **Audiens** | Orang tua muda (*moms & dads*), bayi, dan balita. | Pasien dewasa dermatologi & estetika klinis. |
| **Nuansa Emosional** | Hangat, ceria, higienis, menenangkan (*nurturing*). | Anggun, mewah (*quiet luxury*), bersih, ilmiah medis. |
| **Simbol Brand** | **Maskot "Baby Bean"** (oval mustard dengan senyum inisial *kn*). | **Rosette Mandala** (kelopak bunga geometrik presisi). |
| **Palet Warna Utama** | • **Mustard Warm Gold:** `#E5A93C`<br>• **Forest Bottle Green:** `#0C3B36`<br>• **Warm Linen Cream:** `#FAF7F2` | • **Dusty Rose / Mauve:** `#9E606D`<br>• **Deep Plum Charcoal:** `#2D1E24`<br>• **Soft Silk Blush:** `#FAF5F5` |
| **Tipografi** | *Friendly Bold Rounded* + panah lengkung organik. | *Modern Sans* dipadu aksen **Serif Italic Elegan**. |
| **Elemen Grafis** | Panah lengkung (*organic arrows*), foto cutout anak. | Garis pembagi rambut (*hairline*), perbandingan *do's & don'ts*. |
| **Pita Footer** | Pita hijau botol solid `#0C3B36` teks putih. | Pita dusty rose / slate dengan identitas resmi SMG. |

---

## 3. Sistem Tipografi & Hirarki

Font keluarga:
* **Sans-Serif:** `Plus Jakarta Sans` (jelas, modern, highly legible).
* **Serif Accent:** `Playfair Display` / `Georgia` italic untuk aksen brand SSC (*"Advanced care for your skin"*).

| Tingkat | Ukuran & Bobot | Karakteristik |
| :--- | :--- | :--- |
| **Display Heading** | `text-3xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight` | Skala proporsional (tidak berteriak, nyaman dibaca). |
| **Section Title** | `text-xl sm:text-2xl font-bold tracking-tight` | Jelas, tegas, mengarahkan fokus. |
| **Editorial Accent** | `font-serif italic font-normal` | Digunakan untuk 1–2 kata penegas kemewahan di SSC. |
| **Pill Badge** | `text-xs font-bold tracking-wide uppercase` | `rounded-full px-3.5 py-1` warna mustard atau dusty rose. |
| **Body & Penjelasan** | `text-sm sm:text-base font-normal leading-relaxed` | Spasi baris lega (*line-height 1.6*), warna abu gelap netral. |
| **Nominal & Finansial** | `tabular-nums font-bold` / `tabular-nums font-extrabold tracking-tight` | **WAJIB tabular-nums** (Plus Jakarta Sans). DILARANG `font-mono` untuk harga/uang. `font-mono` hanya untuk ID/NIK/Batch/No Antrian. |

---

## 4. Pola Komponen Baru (Bebas dari Penyakit Card Clutter)

### A. Editorial Index (Pengganti 4 Kotak Card Layanan)
Alih-alih membuat 4 kotak kartu identik, layanan klinik disajikan secara berbaris elegan:
```html
<div class="border-t border-b divide-y divide-zinc-200">
  <div class="py-5 flex items-center justify-between group hover:bg-zinc-50/50 transition-colors">
    <div class="flex items-center gap-6">
      <span class="font-mono text-xs font-bold text-zinc-400">01</span>
      <div>
        <h4 class="text-base font-bold text-zinc-950 group-hover:text-rose-800 transition-colors">
          Laser Skin Rejuvenation
        </h4>
        <p class="text-xs text-zinc-500">Pori, pigmentasi, melasma & peremajaan kolagen kulit</p>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <span class="text-xs font-semibold text-zinc-600 hidden sm:inline">dr. Dian Paramita, Sp.D.V.E</span>
      <button class="text-xs font-bold text-zinc-950 underline underline-offset-4 group-hover:text-rose-700">
        Reservasi →
      </button>
    </div>
  </div>
</div>
```

### B. Perbandingan Format Feed (Do's vs Don'ts / Two-Column Breath)
Mengadopsi format layout resmi feed SSC:
```html
<div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-start relative">
  <!-- Hairline divider di tengah (desktop) -->
  <div class="hidden md:block absolute top-0 bottom-0 left-1/2 w-px bg-zinc-200"></div>

  <!-- Kolom Kiri -->
  <div class="space-y-4 md:pr-8">
    <h3 class="text-lg font-bold text-zinc-950">Do's: Protokol Perawatan</h3>
    <ul class="space-y-2 text-sm text-zinc-600">
      <li class="flex items-start gap-2.5">
        <span class="text-zinc-400">•</span>
        <span>Gunakan tabir surya SPF 30+ 15 menit sebelum keluar ruangan.</span>
      </li>
    </ul>
  </div>

  <!-- Kolom Kanan -->
  <div class="space-y-4 md:pl-8">
    <h3 class="text-lg font-bold text-zinc-950">Don'ts: Kesalahan Umum</h3>
    <ul class="space-y-2 text-sm text-zinc-600">
      <li class="flex items-start gap-2.5">
        <span class="text-zinc-400">•</span>
        <span>Mengaplikasikan krim baru saat kulit masih meradang aktif.</span>
      </li>
    </ul>
  </div>
</div>
```

### C. Live Antrian Strip (Zero Fluff & Scalable Dropdown)
Widget live antrian dirancang menyatu harmonis dengan kanvas halaman, fokus murni pada prioritas pasien:
* **Dropdown Filter Dokter (Skalabel):** Menggantikan deretan tombol tab dengan `<select>` dropdown ramping di kanan atas. Siap menampung pertambahan dokter spesialis (5–20 dokter) tanpa merusak atau melebarkan layout.
* **Bebas Microcopy Redundan (Zero Fluff):**
  * ❌ *Dihapus:* `• Diperbarui Baru Saja` (sudah jelas live).
  * ❌ *Dihapus:* `Tiket terkirim via WhatsApp` (noise visual).
  * ❌ *Dihapus:* Tip bawah panjang `Panggilan otomatis dikirim via WhatsApp saat antrian mendekati 2 giliran`.
* **Prioritas Pasien yang Dipertahankan:**
  * Status badge minimalis: `● Sedang Diperiksa` / `● Sedang Dipanggil`.
  * Nomor antrian besar tajam (`font-mono font-black text-5xl sm:text-6xl`).
  * Nama pasien tersamar & nama dokter pemeriksa.
  * Tiga tiket antrian berikutnya (`K-008 ~8m`, `K-009 ~15m`) + tombol modal riwayat lengkap.
  * Satu aksi utama: `Ambil Antrian Online →`.

### D. Standar Header & Navbar (Sleek Single-Bar + Sister Clinic Link)
* **Tinggi Bar Proporsional:** `h-16` (64px) di desktop, `h-14` (56px) di mobile (menghindari header tebal 120px).
* **Top Micro-Ribbon (Auto-Hide di Mobile):**
  * Di desktop (`hidden sm:block`): Pita tipis `py-1 text-[11px]` berisi alamat resmi SMG dan kontak WA aktif.
  * Di mobile: Otomatis di-hide agar tidak memakan ruang vertikal layar HP.
* **Logo Bersih (Bebas Kotak AI & Subtitle Berjejal):**
  * Logo mark asli + teks brand tegas (`kidinest` / `Semarang Skin Center`).
  * Dihilangkan badge kotak, subtitle panjang, dan dekorasi yang membuat sempit.
* **Sister-Clinic Link (Bukan Switcher Tab):**
  * Jangan gunakan switcher tab `[ Kidinest | SSC ]` yang membingungkan kepemilikan halaman.
  * Gunakan **link teks santai** dengan icon panah keluar: `Layanan Kulit (SSC) ↗` pada Kidinest, dan `Poli Anak (Kidinest) ↗` pada SSC.
* **Hierarki Tombol:** Hanya ada **satu tombol CTA primer** yang solid di navbar: `[ Ambil Antrian → ]`.

### E. Standar Hero Section: Pure Typography & Hairline Split (Zero Foto Produk & Zero Card)
Menghindari jebakan klise AI (*floating card box* di sisi kanan hero):
* **Tanpa Foto Produk / Tanpa Card Widget:** Area kanan hero tidak diisi card widget AI ataupun foto produk yang dipaksakan.
* **Tata Letak Asimetris Bergaris Rambut (`lg:col-span-7` vs `lg:col-span-5` dengan `lg:border-l`):**
  * **Sisi Kiri:** Tag brand tenang, Display H1 berbobot manusiawi (*line-height 1.15*), narasi pengantar 2 baris, dan tombol aksi utama.
  * **Sisi Kanan:** Dibatasi garis pembagi 1px tipis vertikal, memuat pilar klinis murni berbasis tipografi:
    * Headline filosofi pelayanan klinik.
    * 4 baris hairline bernomor monospaced (`01–04`) yang merangkum fokus keunggulan medis.
    * Sign-off dokter spesialis penanggung jawab.
* **Strip Fakta Klinik 3-Kolom (Hairline Terbuka):**
  * Tepat di bawah narasi hero membentang 3 kolom terbuka tanpa kotak card:
    * `01 • Dokter Spesialis Jaga Hari Ini`
    * `02 • Sesi Layanan Praktik (Pagi & Sore)`
    * `03 • Lokasi & Fasilitas Resmi`
  * Berfungsi sebagai transisi alami yang langsung mengalirkan pandangan ke Papan Antrian di bawahnya.

---

## 5. Standar Antarmuka Admin & Workstation Minimalist (*Workstation Ergonomics*)

Semua modul internal klinik (Kasir, Dokter, Farmasi, Gudang, Billing, Owner, SATUSEHAT) **WAJIB** merujuk pada standar arsitektur dan ergonomi visual yang telah diterapkan pada **`#cashier` (`CashierView.vue`)**.

### A. Anatomical Anatomy Layar Workstation (Standar #cashier)
Setiap halaman workstation memiliki susunan visual yang konsisten:
```
┌────────────────────────────────────────────────────────────────────────┐
│ 1. Minimalist Clinical Header (Breadcrumb, Title, Duty Badge, Sync)     │
├────────────────────────────────────────────────────────────────────────┤
│ 2. Compact Financial / Operational Ticker Strip (1-Row Hairline Grid)  │
├────────────────────────────────────────────────────────────────────────┤
│ 3. Navigation Tabs & Segmented Filters (Pill Tabs + Monospace Counters) │
├────────────────────────────────────────────────────────────────────────┤
│ 4. Master-Detail Two-Pane Grid (lg:grid-cols-12 gap-6)                 │
│   ├── Left Pane (Queue / Catalog Feed, 5-7 Cols)                       │
│   └── Right Pane (Active Action / Execution Drawer, 5-7 Cols Sticky)   │
└────────────────────────────────────────────────────────────────────────┘
```

1. **Wadah Utama (Canvas Container):**
   * Kelas wajib: `py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 font-['Plus_Jakarta_Sans',sans-serif]`.
   * Ruang pandang lega tanpa batasan sempit `max-w-7xl`, menghindari horizontal scroll pada monitor desktop klinik (1080p ke atas).

2. **Header Klinis Minimalis (Top Bar):**
   * **Breadcrumb Institusi:** `text-xs font-semibold text-zinc-500 uppercase tracking-wider` dengan titik pembatas `text-zinc-300` (contoh: `Meja Depan & Kasir • Semarang Medical Group`).
   * **Judul Utama:** `text-2xl font-bold text-zinc-950 tracking-tight mt-0.5`.
   * **Duty Badge:** Sisi kanan menampilkan `Petugas Bertugas` (font 10px uppercase zinc-400) dan nama user yang sedang login dalam `text-xs font-bold text-zinc-900`.
   * **Tombol Aksi Sinkronisasi:** Tombol `Sync` dengan ikon `RefreshCw` berotasi saat loading (`border border-zinc-200 bg-white hover:bg-zinc-50 text-xs font-bold text-zinc-800 shadow-2xs rounded-xl`).

3. **1-Row Compact KPI Ticker Strip (Hairline Grid):**
   * Menggantikan kotak-kotak card besar berbayang tebal. Menggunakan satu baris bersambung dengan garis pembagi 1px:
   * Kelas pembungkus: `grid grid-cols-2 sm:grid-cols-4 lg:grid-cols-5 divide-y sm:divide-y-0 sm:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden`.
   * Tiap sel metrik:
     * Padding proporsional: `p-3.5`.
     * Label: `text-[10px] font-bold uppercase tracking-wider text-zinc-500`.
     * Nilai: `text-base font-bold tabular-nums tracking-tight text-zinc-900 mt-0.5`.
     * Sel summary/total akhir: `bg-stone-50/90 border-l border-zinc-200/90 text-lg font-extrabold text-zinc-950`.

4. **Navigasi Tab Sub-Modul:**
   * Tombol tab bergaya pill netral:
     * Aktif: `bg-zinc-950 text-white shadow-2xs`.
     * Tidak aktif: `text-zinc-600 hover:text-zinc-950 hover:bg-zinc-100`.
   * Badge angka counter: `px-1.5 py-0.2 rounded-md text-[10px] tabular-nums font-bold` (`bg-zinc-800 text-zinc-200` saat aktif, `bg-stone-200 text-zinc-800` saat pasif).

5. **Tata Letak Master-Detail (Two-Pane Workspace):**
   * Grid desktop: `grid grid-cols-1 lg:grid-cols-12 gap-6 items-start`.
   * **Left Pane (5 atau 7 kolom):** Daftar antrean atau katalog barang, dilengkapi pencarian dan filter unit (`ALL`, `KIDINEST`, `SSC`).
   * **Right Pane (7 atau 5 kolom):** Panel eksekusi tindakan (input resep, bayar tagihan, telaah obat, atau audit) dengan posisi `sticky top-6`.

---

## 6. Standar Komponen & Status Medis (*Design Tokens*)

### A. Palet Warna Status Medis (Standardized Semantic Tokens)
Hindari pemakaian warna-warna neon acak. Gunakan palet desaturated clinical:
* **Ready / Selesai / Lunas:** Emerald/Zinc tint (`bg-emerald-50 text-emerald-900 border-emerald-200` atau bullet `bg-emerald-500`).
* **In-Progress / Sedang Diperiksa / Racik:** Amber tint (`bg-amber-50 text-amber-950 border-amber-200/90`, border aktif `border-l-2 border-l-amber-500`, bullet `bg-amber-500`).
* **Waiting / Standby / Antre:** Zinc/Stone tint (`bg-zinc-100 text-zinc-800 border-zinc-200`).
* **Alert / Batal / Refund / Overdue:** Rose tint (`bg-rose-50 text-rose-900 border-rose-200`, teks `text-rose-700`).
* **Brand Kidinest Unit Tag:** `bg-amber-50 text-amber-900 border-amber-200/80`.
* **Brand SSC Unit Tag:** `bg-rose-50 text-rose-900 border-rose-200/80`.

### B. Kotak Kontainer Card & Hairline Dividers
* Seluruh kontainer data menggunakan: `border border-zinc-200/90 rounded-2xl bg-white shadow-2xs`.
* DILARANG menggunakan bayangan besar seperti `shadow-xl` atau `shadow-2xl` pada card halaman statis (hanya boleh pada dialog modal melayang).
* Pemisah baris item di dalam list: `divide-y divide-zinc-100` dengan baris interaktif `hover:bg-zinc-50/70 transition-colors`.

### C. Modal Dialog Standar
* Backdrop: `fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4`.
* Modal Box: `w-full max-w-lg bg-white rounded-3xl border border-zinc-200/80 shadow-2xl overflow-hidden`.
* Header Modal: `p-5 border-b border-zinc-100 flex items-center justify-between` (atau dark header `px-6 py-4 bg-zinc-950 text-white`).
* Footer Aksi: `p-4 border-t border-zinc-100 bg-zinc-50/70 flex items-center justify-end gap-2.5`.

---

## 7. Standar Desain Modul Akuntansi & Keuangan (Finance & Accounting)

Modul Akuntansi (`/accounting/*`) diselaraskan penuh dengan standar workstation medis SMG guna mengeliminasi estetika "AI Slop" dan blue SaaS generic:

1. **Kanvas & Navigasi:**
   - Kanvas kerja: `py-6 px-4 sm:px-8 max-w-[1600px] mx-auto space-y-5 bg-slate-50 text-zinc-950 font-['Plus_Jakarta_Sans',sans-serif]`.
   - Sidebar: Monokrom `bg-zinc-950 text-zinc-300 border-r border-zinc-800` dengan identitas resmi `SMG Finance & Accounting`.
   - Dual-brand pills pada Navbar: Kidinest (`bg-amber-50 text-amber-900 border-amber-200/80`), SSC (`bg-rose-50 text-rose-900 border-rose-200/80`), Konsolidasi (`bg-zinc-100 text-zinc-800 border-zinc-200`).

2. **Minimalist Clinical Header Akuntansi:**
   - Breadcrumb: `text-[11px] font-bold text-zinc-400 tracking-wider uppercase` (`Akuntansi & Keuangan • Semarang Medical Group`).
   - Judul H1: `text-2xl font-black text-zinc-950 tracking-tight`.
   - Action buttons: Menggunakan `bg-zinc-950 hover:bg-zinc-800 text-white rounded-xl shadow-xs`.

3. **1-Row Compact Hairline KPI Ticker:**
   - Diterapkan pada Dashboard, Buku Kas & Bank, dan Manajemen Pajak:
   - `grid grid-cols-2 md:grid-cols-4 divide-y md:divide-y-0 md:divide-x divide-zinc-200/90 border border-zinc-200/90 rounded-2xl bg-white shadow-2xs overflow-hidden`.
   - Seluruh nominal moneter diformat dengan `tabular-nums font-bold text-zinc-950`.

4. **Hierarki Bagan Akun (COA):**
   - Level 1 (Akun Induk): `bg-zinc-900 text-white font-bold tracking-wider uppercase` atau `bg-zinc-100 font-bold`.
   - Level 2–3 (Golongan/Sub-golongan): Indentasi bertahap `pl-6` dan `pl-10`.
   - Level 4 (Leaf / Akun Posting): Indentasi `pl-14`, nomor akun `font-mono`, saldo `tabular-nums font-bold text-right`.

5. **Sticky Bottom Balance Dock:**
   - Pada lembar kerja panjang (seperti Saldo Awal dan Jurnal Penyesuaian), pasang dock melayang di bawah:
   - `fixed bottom-4 left-68 right-8 z-30 bg-white/95 backdrop-blur-md rounded-2xl border border-zinc-200/90 p-4 shadow-xl flex items-center justify-between`.
   - Menampilkan total debit vs kredit secara instan dengan indikator status seimbang (Emerald) atau selisih (Rose).

6. **Ergonomi Input Jurnal:**
   - Mode entri cepat (*Direct Fast-Entry*) langsung pada baris tabel dengan auto-balancing baris penyeimbang kredit.
   - Form modal standar: Wadah `rounded-3xl border border-zinc-200/80 shadow-2xl`, tombol solid `zinc-950`, dan preview mutasi jurnal otomatis.
