# Breakdown Analisis Desain Visual Kidinest Pediatric Care

Dokumen ini membedah elemen desain asli dari aset visual resmi [@kidinest.id](https://www.instagram.com/kidinest.id/) (berdasarkan feed Instagram dan file maskot yang diunggah) sebagai acuan presisi pembuatan landing page `/kidinest`.

---

## 1. Bedah Aset Visual Resmi

| Referensi | Gambar | Karakteristik Kunci |
| :--- | :--- | :--- |
| **Foto 1 (Feed)** | *Benefits of Massage For Babies* | • Latar kanvas linen rajut krem hangat (`#FAF7F2`).<br>• Tipografi judul hijau botol tua (`#134E48`) membulat tebal (*friendly bold rounded*).<br>• Pill badge kuning mustard (`#E5A93C`) dengan teks hijau tua: *"Helpful Healing Tips & Techniques"*.<br>• Foto bayi asli berbaring dalam balutan rajutan pastel.<br>• Panah melengkung tipis organik (*hand-drawn organic arrows*) mengarah ke 6 poin manfaat. |
| **Foto 2 (Feed)** | *How Many Words Can My Baby Say?* | • Latar belakang pola gelombang organik warna mustard hangat (`#E5A93C`).<br>• Judul hijau botol dengan **efek sticker outline putih tebal**.<br>• Foto cutout wajah balita ekspresif di tengah.<br>• Kartu indikator milestone usia (12 bulan, 18 bulan, 24 bulan, 36 bulan) berujung oval membulat.<br>• Pita footer horizontal hijau botol pekat (`#0C3B36`) berisi alamat: *"Jl. Melati Utara No 3, Semarang"*. |
| **Foto 3 (Maskot)** | *Maskot Resmi Kidinest ("Baby Bean")* | • Karakter oval kuning mustard hangat (`#E5A93C`) miring riang.<br>• Sehelai rambut ikal spiral hijau botol di atas kepala.<br>• Wajah tersenyum yang cerdik dibentuk dari inisial **"k"** (mata kiri kedip) dan **"n"** (mata kanan), dengan garis kurva senyum lebar (*Kidinest = kn*). |

---

## 2. Palet Warna Resmi Kidinest (Design Tokens)

| Token Warna | Nilai HEX | Peruntukan UI |
| :--- | :--- | :--- |
| **Kidinest Mustard Gold** | `#E5A93C` / `#D99B26` | Warna maskot, badge highlight, aksen gelombang hero, tombol sekunder. |
| **Kidinest Forest Green** | `#0C3B36` / `#134E48` | Warna teks heading, garis kontur maskot, tombol primer, pita footer. |
| **Kidinest Warm Cream Canvas** | `#FAF7F2` / `#FFFDF9` | Latar belakang seluruh halaman (hangat, menenangkan, bukan abu-abu dingin). |
| **Kidinest Light Mustard Tint** | `#FEF3C7` / `#FDE68A` | Latar kartu edukasi, tag usia, gelombang latar belakang sekunder. |
| **Kidinest Mint Soft** | `#D1FAE5` / `#A7F3D0` | Tag status aktif, variasi karakter maskot sekunder. |
| **Sticker White Stroke** | `#FFFFFF` | Outline tebal di sekitar teks judul dan foto cutout anak. |

---

## 3. Tipografi & Gaya Teks

1. **Heading (Judul Utama & Section):**
   * Gaya *Friendly Rounded Bold* (geometrik membulat, tebal, tidak kaku).
   * Warna wajib: Hijau Botol Tua (`#0C3B36`).
   * Pada latar belakang berwarna mustard/gelombang, gunakan efek teks outline putih (*stroke* / shadow tebal) seperti pada feed Instagram.
2. **Pill Badges:**
   * Sudut membulat penuh (`rounded-full`).
   * Dua variasi resmi:
     * Latar Mustard Gold (`#E5A93C`) dengan teks Hijau Tua (`#0C3B36`).
     * Latar Hijau Tua (`#0C3B36`) dengan teks Putih (`#FFFFFF`).
3. **Pita Alamat Bawah (Footer Ribbon):**
   * Pita hijau botol solid `#0C3B36` dengan teks putih bersih: *"Jl. Melati Utara No 3, Semarang"*.

---

## 4. Elemen Grafis Khas (Brand Signatures)

1. **Maskot "Baby Bean" (Wajah 'kn'):**
   * Digunakan pada Logo Header, favicon, dan hiasan interaktif section anak.
2. **Panah Lengkung Organik (*Curved Organic Arrows*):**
   * Garis panah tipis melengkung yang menghubungkan foto anak dengan informasi manfaat layanan atau milestone perkembangan.
3. **Cutout Foto Bayi/Anak Nyata:**
   * Foto anak nyata beresolusi tinggi tanpa latar klinik rumah sakit yang kaku, melainkan bernuansa hangat di atas kain rajut atau latar gelombang pastel.

---

## 5. Rencana Penerapan pada `LandingKidinest.vue`

1. **Header Navigation:**
   * Menampilkan Maskot Resmi Kidinest (`/assets/kidinest/logo-maskot.png`) + Logotype `kidinest` dengan tipografi asli.
   * Switcher ke Semarang Skin Center (SSC) + Tombol CTA *"Ambil Antrian Online"*.
2. **Hero Section:**
   * Latar gelombang organik Mustard Warm (`#E5A93C`) dan Linen Cream (`#FAF7F2`).
   * Headline: *"Pemeriksaan Buah Hati dengan Hangat & Tenang"* dengan nuansa Hijau Botol tebal.
   * Visual kartu dokter: dr. Patricia Vanessa A., Sp.A diapit elemen maskot Baby Bean dan panah penunjuk fitur (Imunisasi Lengkap, Konsultasi Sakit, Antrian Live WhatsApp).
3. **Widget Live Antrian (Poli 1 dr. Patricia & Poli 2 dr. Hendra):**
   * Selaras dengan token warna Mustard & Forest Green, pill filter dokter yang responsif.
4. **Section Edukasi & Layanan (Model Feed Massage & Words):**
   * 4 kartu layanan utama menggunakan pola kartu feed: ikon ramah anak, panah indikator, dan poin panduan orang tua.
5. **Footer Bar:**
   * Mengadopsi pita hijau botol pekat resmi dengan teks *"Jl. Melati Utara No 3, Semarang"*.
