# Analisis Desain Anti-AI Slop: Tipografi Angka, Harmoni Warna & Reduksi Badge (GRO / Kasir)

Dokumen ini membedah secara kritis mengapa antarmuka **GRO / Kasir** (`CashierView.vue`) saat ini masih terasa *"AI-generated / template AI"*, serta merumuskan standar visual baru yang berakar pada prinsip desain editorial klinis (*Nordic Modern / Clinical Editorial*) dan pedoman resmi di [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md).

---

## 1. Diagnosis Kritis: Mengapa Masih Kelihatan "AI Slop"?

Meskipun struktur tata letak sudah dipadatkan menjadi *master-detail 2-pane*, ada 3 elemen visual mikro yang secara instan memberikan impresi *"dibuat otomatis oleh generator AI"*:

### A. Tipografi Angka: Sindrom Monospace Mentah (`font-mono`)
* **Masalah:** Penerapan class `font-mono` bawaan Tailwind mengevaluasi font ke monospace sistem standar (`Courier New`, `Menlo`, atau font monospaced konsol Linux).
* **Kenapa terasa AI?**
  1. Font monospace mentah membuat angka nominal rupiah renggang (*kerning* lebar tidak wajar, spasi di sekitar angka `1` atau `Rp` terlihat bolong).
  2. Memberikan estetika *"terminal programmer / developer tool / hacker dashboard"*, bukan sistem manajemen klinik medis kelas atas atau POS kasir premium.
* **Solusi Manusiawi (Human-Crafted):**
  - Gunakan font Sans proporsional dengan fitur tabular (*Tabular Numbers* / `tabular-nums font-sans tracking-tight`).
  - Seluruh angka tetap sejajar secara vertikal (digit 0–9 memiliki lebar kolom yang sama persis sehingga kalkulasi rapi), namun bentuk glifnya tetap elegan, rapat, dan modern mengikuti font utama (`Plus Jakarta Sans`).

---

### B. Warna & Kontras: Kontras Kasar & Tabrakan Palet Tailwind Bawaan
* **Masalah:**
  1. **Kotak Hitam Jarang (The Inverted Black Box):** Pada Ticker Keuangan, 4 kolom pertama berlatar putih dengan border tipis, tetapi kolom ke-5 (*Net Setoran Shift*) tiba-tiba memakai `bg-zinc-950 text-white` pekat dengan teks hijau neon `text-emerald-400`. Ini adalah klise template AI modern (terlalu ekstrem, seperti widget crypto/SaaS).
  2. **Bukan Palet Brand Resmi:** Menggunakan warna default Tailwind seperti `text-purple-900`, `text-amber-400`, `text-rose-600`, dan `text-emerald-400`. Warna-warna ini tidak selaras dengan identitas brand di `DESIGN.md`:
     - **Kidinest:** *Mustard Warm Gold* (`#E5A93C`), *Forest Bottle Green* (`#0C3B36`), *Warm Linen* (`#FAF7F2`).
     - **SSC:** *Dusty Rose / Mauve* (`#9E606D`), *Deep Plum Charcoal* (`#2D1E24`), *Soft Silk Blush* (`#FAF5F5`).
  3. **Aksesibilitas Kontras (WCAG AA):** Teks label `text-zinc-400` berukuran `10px` pada latar putih sulit dibaca di bawah pencahayaan lampu neon klinik meja depan (*low legibility*).
* **Solusi Manusiawi:**
  - Satukan seluruh Ticker Keuangan dalam satu kanvas hangat yang tenang (*uniform background*), tanpa kotak hitam mendadak.
  - Net Setoran ditonjolkan dengan bobot tipografi (*weight & size*) dan warna tinta yang berwibawa (*Deep Slate/Charcoal* atau *Forest Green* tenang), bukan neon.
  - Naikkan kontras microcopy dari `text-zinc-400` ke `text-zinc-500` / `text-zinc-600` agar terbaca jelas oleh kasir.

---

### C. Badge Soup (Sindrom Permen Warna-Warni)
* **Masalah:**
  - Setiap baris antrian ditempeli 3 hingga 4 badge berbentuk kapsul pill tebal (`rounded-full px-2.5 py-0.5`):
    - Badge Unit: `KIDINEST` (kuning) / `SSC` (merah muda)
    - Badge Kategori: `POLI_ANAK` / `TREATMENT`
    - Badge Status Antrian: `WAITING` / `IN_CONSULTATION`
    - Badge Status Bayar: `UNPAID` / `PAID`
  - Hasilnya adalah visual *"confetti"* yang ramai dan melelahkan mata. Petugas kasir harus memindai puluhan kapsul warna-warni sebelum melihat nama pasien.
* **Kenapa terasa AI?**
  - Model AI secara default gemar menambahkan status badge untuk setiap properti JSON yang diterimanya, tanpa memikirkan redundansi informasi.
* **Solusi Manusiawi:**
  - **Prinsip Kontekstual (Zero Redundant Badges):** Jika pasien sudah berada di tab/bagian *"Menunggu Pembayaran"*, tidak perlu ada badge bertuliskan *"Belum Lunas"*.
  - **Ganti Pill Tebal dengan Subtle Typographic Marker:**
    - Identitas brand Kidinest vs SSC: Cukup garis aksen vertikal 2px di tepi kiri kartu/baris atau teks kecil bernada tenang (`KIDINEST` dalam mustard redup, `SSC` dalam dusty mauve).
    - Status ruangan/antrian: Cukup dot tenang 6px (`w-1.5 h-1.5 rounded-full`) yang berdampingan dengan teks netral, bukan kapsul besar yang mencolok.

---

## 2. Tabel Perbandingan: AI Slop vs Editorial Human Design

| Elemen | Pendekatan AI Slop Saat Ini | Standar Human-Crafted (Nordic Clinic) |
| :--- | :--- | :--- |
| **Font Angka** | `font-mono` default (lebar tidak rata, kerning bolong). | `font-sans tabular-nums font-semibold tracking-tight`. |
| **Ticker Finansial** | Kotak hitam kontras tinggi `bg-zinc-950` + teks neon `text-emerald-400`. | Kanvas putih/linen terpadu, penegasan pada ukuran angka & tinta deep charcoal `#18181b`. |
| **Badge Unit** | Kapsul tebal `rounded-full bg-amber-100 text-amber-900`. | Subtle hairline prefix / small uppercase muted label (`text-[10px] text-zinc-500 font-semibold tracking-wider`). |
| **Status Antrian** | Pill lonjong berwarna-warni di setiap baris. | Single 6px subtle status dot + teks status 11px yang tenang. |
| **Skema Warna** | Tailwind generic palette (purple, neon emerald, bright amber, pure black). | Dual-brand palette terkurasi (Warm mustard gold, dusty rose, deep plum, warm zinc). |
| **Kepadatan Informasi** | 4 badge per baris (nama tertutup badge). | Nama pasien & nomor antrian adalah pahlawan utama (*hierarchy first*). |

---

## 3. Rencana Perbaikan Teknis di `CashierView.vue`

1. **Definisi Utilitas Angka Tabular:**
   - Mengganti seluruh class `font-mono` untuk angka nominal dan nomor antrean menjadi:
     ```html
     class="font-sans font-bold tabular-nums tracking-tight text-zinc-900"
     ```
   - Nomor tiket antrean (misal `E-001`, `K-012`): Menggunakan `font-sans font-extrabold tabular-nums tracking-tighter` sehingga padat, tegas, dan modern.

2. **Refactor Ticker Keuangan (Harmoni & Anti-Neon):**
   - Menghilangkan `bg-zinc-950` pada kolom ke-5.
   - Menggunakan layout 5 kolom berlatar belakang sama (`bg-white` dengan border halus `border-zinc-200/90`).
   - Kolom *"Net Setoran Shift"* ditonjolkan dengan background subtle linen/warm tint (`bg-stone-50/80`) dengan angka berukuran lebih besar (`text-lg font-bold text-zinc-950`).

3. **Pembersihan Badge (De-Slopification):**
   - Hapus pill tebal pada list pending bills dan reservasi online.
   - Ganti identitas unit menjadi hairline tag minimalis:
     - Kidinest: `text-amber-800 bg-amber-50/60 border border-amber-200/60 text-[10px] px-1.5 py-0.5 rounded`
     - SSC: `text-rose-800 bg-rose-50/60 border border-rose-200/60 text-[10px] px-1.5 py-0.5 rounded`
   - Ganti status ruangan dokter dengan dot minimalis bersih (tanpa badge background tebal).

4. **Koreksi Kontras & Hierarki:**
   - Naikkan kontras teks bantu (`text-zinc-400` → `text-zinc-500` / `text-zinc-600`).
   - Pertegas judul item dan data finansial agar sekilas pandang langsung terbaca oleh kasir tanpa silau atau kebingungan.
