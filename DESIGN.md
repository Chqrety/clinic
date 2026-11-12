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

## 5. Standar Antarmuka Admin Minimalist (*Workstation Ergonomics*)

Untuk operasional internal staf (Meja Kasir, Dokter, Apotek, Gudang):
1. **Sidebar Kiri Tetap (Fixed 240px):**
   * Latar gelap minimalis (`bg-zinc-950`), icon rapi, label ringkas.
   * Modul cepat: Kasir & GRO, Ruang Dokter, Farmasi Resep, Gudang Stok, Billing Shift, SATUSEHAT.
   * Role Switcher instan di bagian profil atas untuk kemudahan demo klien.
2. **Kepadatan Data Tinggi (*High-Density Table*):**
   * Font 12–13px, padding 8–10px, angka monospaced tabular.
   * Tanpa grafik dekoratif yang tidak perlu; utamakan status antrian, nama pasien, dan aksi klinis langsung.
