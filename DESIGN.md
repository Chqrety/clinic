# Design System & UI/UX Guidelines
## Inspired by Modern Nordic / Brandlos Aesthetic

---

## 1. Visual Philosophy: Anti "AI Slop"
Sistem klinik ini **TIDAK** menggunakan template generik khas AI (latar ungu neon, gradien berlebihan, font kaku, atau card berantakan tanpa hirarki). 

Aplikasi ini mengadopsi estetika **Brandlos (Dutch / Nordic Modern)**:
* **Hangat & Ramah Pasien:** Menghilangkan kesan dingin dan menakutkan dari rumah sakit konvensional.
* **Tipografi Berani & Percaya Diri:** Angka dan status langsung terbaca jelas tanpa perlu memicingkan mata.
* **Minimalis & Bersih:** Ruang bernapas (*whitespace*) lapang, memisahkan konten penting dengan tegas.

---

## 2. Color Palette & Token

| Token | Nilai / Tailwind Class | Peruntukan |
| :--- | :--- | :--- |
| **Canvas Background** | `bg-[#F9F9F7]` atau `bg-stone-50` | Latar utama (Warm Off-White, ramah di mata) |
| **Surface / Card** | `bg-white` | Kartu konten, modal, container |
| **Borders** | `border-black/5` atau `border-zinc-200/70` | Garis batas tipis elegan (bukan abu-abu tebal) |
| **Text Primary** | `text-zinc-950` atau `text-black` | Judul, angka antrian, teks utama |
| **Text Secondary** | `text-zinc-500` | Label, metadata, timestamp, deskripsi |
| **Brand Accent (Vibrant Blue)** | `bg-blue-600` / `text-white` | Pill badge kategori, link aktif, highlight |
| **Action Primary (Warm Amber)** | `bg-amber-400` / `hover:bg-amber-500` | Tombol CTA utama, tombol bayar, konfirmasi |
| **Healthcare Green** | `bg-emerald-500` / `text-emerald-700` | Status Lunas, Resep Siap, Stok Aman |
| **Alert / Expired** | `bg-rose-500` / `text-rose-700` | Peringatan ED < 3 bulan, Stok Kosong |

---

## 3. Tipografi & Skala Hirarki
Font keluarga: **Plus Jakarta Sans** atau **Inter** (sans-serif modern).

* **Display / Queue Number:** `text-4xl` hingga `text-6xl`, `font-extrabold`, `tracking-tight`, `text-zinc-950`.
* **Section Heading:** `text-2xl`, `font-bold`, `tracking-tight`.
* **Pill Badges:** `text-xs`, `font-semibold`, `tracking-wide`, `uppercase`.
* **Body / Label:** `text-sm`, `font-normal`, `text-zinc-600`.

---

## 4. Komponen Khas (Brandlos Style Elements)

### A. Pill Badges
Badge melengkung penuh untuk status dan tag:
```html
<!-- Kategori / Role -->
<span class="inline-flex items-center px-3.5 py-1 rounded-full text-xs font-semibold bg-blue-600 text-white shadow-sm">
  Poliklinik Umum
</span>

<!-- Status Live -->
<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200/60">
  <span class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
  Antrian Berjalan
</span>
```

### B. Rounded Surface Cards
Card dengan sudut membulat lebar dan border super halus:
```html
<div class="bg-white rounded-3xl p-6 border border-black/5 shadow-[0_2px_12px_rgba(0,0,0,0.03)] hover:shadow-[0_8px_24px_rgba(0,0,0,0.05)] transition-all duration-200">
  ...
</div>
```

### C. Buttons with Tactile Feel
Tombol aksi dengan efek tekan (*active scale*):
```html
<button class="px-6 py-3 rounded-full font-semibold bg-amber-400 text-zinc-950 hover:bg-amber-500 active:scale-[0.98] transition-all duration-150 shadow-sm flex items-center gap-2">
  Panggil Pasien Selanjutnya →
</button>
```

---

## 5. Micro-Interactions & Animasi Halus
Semua animasi harus **cepat dan halus** (150ms – 250ms), jangan lambat yang membuat kerja staf klinik terhambat:
* Gunakan Vue `<Transition>` dan `<TransitionGroup>` pada:
  * Masuk/keluarnya tiket antrian.
  * Baris resep saat dokter menambah obat.
  * Kartu resep baru di apotek yang meluncur masuk (*slide up* halus).
* Modal/dialog menggunakan latar belakang buram (`backdrop-blur-sm bg-black/20`).

---

## 6. Layout Ergonomics Menurut Pengguna

1. **Pasien (Mobile-First):**
   * Mengikuti gaya Brandlos 100%: Judul besar, sapaan personal hangat (*"Halo, selamat datang di Klinik Sehat"*), kartu tiket antrian digital yang menyerupai boarding pass / kartu bank digital, tombol aksi besar ramah jempol.
2. **Dokter & Apotek (Workstation Ergonomics):**
   * Tetap menggunakan warna dasar warm sand dan pill badge, namun dengan **densitas data yang efisien**.
   * Daftar obat, rekam medis, dan checklist racikan tersaji dalam card/tabel rapi yang mudah di-scan mata dan minim scrolling berlebihan.
3. **Demo Role Switcher:**
   * Bar mengambang (*floating pill bar*) di bagian atas tengah layar:  
     `[ Pasien ] [ Kasir ] [ Dokter ] [ Apotek ] [ Gudang ] [ SATUSEHAT ]`  
     Memudahkan presenter mendemonstrasikan end-to-end flow hanya dengan 1 klik.
