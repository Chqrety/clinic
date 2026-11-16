# 10. Rencana Redesign UI Pemanggilan Pasien & Layar Monitor Antrian (Brandlos / Nordic Modern)

## 1. Latar Belakang & Analisis Masalah Saat Ini
Sistem kasir meja depan (*Frontdesk Cashier*) telah berhasil direstrukturisasi dengan estetika **Nordic Modern / Brandlos** (tampilan bersih, warm off-white canvas, hairline border tipis, pill badge kontras tinggi, dan tanpa AI-slop).

Namun, antarmuka pemanggilan pasien saat ini:
1. **Layar TV Display Ruang Tunggu ([QueueMonitor.vue](file:///home/alpi/Projects/clinic/frontend/src/views/QueueMonitor.vue)):**
   - Masih menggunakan tata letak prototipe awal yang kurang kontras saat ditampilkan di layar TV TV/monitor dinding dari jarak 5–10 meter.
   - Belum menampilkan status ruangan dokter secara simultan (Ruang Poli Anak 1, Poli Anak 2, dan Ruang Laser SSC).
   - Suara pemanggilan baru berupa bunyi beep oscillator sederhana tanpa suara pemanggilan suara (Text-to-Speech / TTS Indonesia).
2. **Kontrol Pemanggilan Dokter & GRO ([DoctorRoom.vue](file:///home/alpi/Projects/clinic/frontend/src/views/DoctorRoom.vue) & [CashierView.vue](file:///home/alpi/Projects/clinic/frontend/src/views/CashierView.vue)):**
   - Indikator feedback saat pasien sedang berstatus `CALLING` perlu diselaraskan agar dokter dan petugas meja depan melihat status ruangan yang sinkron secara real-time.

---

## 2. Prinsip Desain & Visual Hierarchy (Selaras UI Kasir)

Sesuai [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md) dan standar anti AI-slop:
* **Canvas & Surfaces:** Latar belakang warm off-white (`bg-[#FAF8F5]` / `bg-zinc-50`), permukaan kartu putih murni (`bg-white`) dengan hairline border `border-zinc-200/80` tanpa drop-shadow tebal.
* **Tipografi & Skala Pembacaan TV (Legibility):**
  - Nomor antrian yang sedang dipanggil menggunakan ukuran raksasa (`text-8xl` s/d `text-9xl`) dengan font tabular mono bold, terlihat tajam dari jarak jauh.
  - Nama pasien disamarkan otomatis (*privacy-first*, contoh: `Budi P****`) dengan label proteksi privasi medis.
* **Multi-Room Status Grid:** Menampilkan kartu status tiap ruangan periksa dokter yang sedang praktek:
  - Ruang 1: dr. Patricia Vanessa A., Sp.A (Poli Anak) — Status: `SEDANG DIPERIKSA`
  - Ruang 2: dr. Hendra Pratama, Sp.A (Tumbuh Kembang) — Status: `MEMANGGIL A-002`
  - Ruang 3: dr. Dian Paramita, Sp.D.V.E (Ruang Tindakan SSC) — Status: `SIAP / MENUNGGU`
* **Pill Badges & Warna Status Bersih:**
  - `MEMANGGIL`: Badge emerald/zinc kontras dengan animasi pulse tenang (*quiet ping*).
  - `SEDANG PERIKSA`: Badge amber lembut.
  - `MENUNGGU`: Badge stone/slate netral.
  - Tanpa gradien ungu/pink neon murah (Anti AI-Slop).

---

## 3. Fitur Pemanggilan Suara Pintar (Voice Synthesizer / TTS)
* **Audio Cue Ganda:**
  1. *Lonceng Pembuka (Harmonic Chime):* Bunyi 2-nada lembut (587 Hz $\rightarrow$ 880 Hz).
  2. *Voice Speech Synthesizer (Bahasa Indonesia):* Menggunakan Web Speech API bawaan browser (`window.speechSynthesis` dengan `lang: "id-ID"`) untuk melafalkan:
     > *"Nomor antrian, A, nol, nol, satu. Silakan menuju Ruang Poli Anak Satu."*
* **Toleransi Panggilan & Limitasi:** Tombol "Panggil Ulang" di ruang dokter dan kasir memicu kembali suara pemanggilan hingga batas maksimal 3 kali panggilan sebelum dialihkan ke antrian toleransi.

---

## 4. Langkah Implementasi Teknis

### Tahap 1: Redesign Layar TV Monitor Ruang Tunggu (`QueueMonitor.vue`)
1. **Header TV Editorial:**
   - Brand logo Kidinest & Semarang Skin Center.
   - Jam digital font mono besar dengan detik real-time (`HH:mm:ss WIB`).
   - Kontrol suara (Mute/Unmute) dan filter unit (Semua Layanan / Kidinest / SSC).
2. **Hero Calling Panel (Sisi Kiri / 7 Kolom):**
   - Nomor antrian raksasa dengan animasi highlight transisi saat berganti pasien.
   - Nama pasien tersensor, nama dokter pemeriksa, dan nomor ruangan tujuan.
   - Running text informasi klinik di bagian bawah.
3. **Multi-Room & Antrian Berikutnya (Sisi Kanan / 5 Kolom):**
   - Mini-card status live per ruangan dokter aktif hari ini.
   - Daftar 5 antrian berikutnya yang wajib bersiap di depan poliklinik.

### Tahap 2: Voice Announcement Engine di Frontend
* Buat helper fungsi suara di `QueueMonitor.vue` atau `frontend/src/utils/audioAnnouncement.ts`:
  - `announceQueue(queueNumber: string, roomName: string, doctorName: string)`
  - Deteksi ketersediaan voice Bahasa Indonesia (`id-ID`) pada browser (Chrome/Edge/Safari).
  - Fallback ke chime audio jika browser tidak mendukung Text-to-Speech lokal.

### Tahap 3: Penyelarasan Tombol Panggil di Ruang Dokter (`DoctorRoom.vue`)
* Perbarui kartu pasien aktif di ruang dokter:
  - Tampilkan tombol aksi terpadu: **"Panggil Pasien"** (saat pasien status WAITING), **"Panggil Ulang"** (saat CALLING), dan **"Mulai Periksa"** (masuk ke status IN_EXAM).
  - Indikator badge status live yang selaras dengan layar monitor antrian.

### Tahap 4: Penyelarasan Kartu Ruangan GRO di Kasir (`CashierView.vue`)
* Perbarui kartu "Ruang Dokter & Antrian" di tab GRO Meja Depan:
  - Tambahkan counter berapa kali pasien telah dipanggil.
  - Tampilan visual selaras dengan kartu display monitor.

---

## 5. Rencana Pengujian (Verification)
1. **Uji Tampilan Resolusi TV:**
   - Verifikasi layout di resolusi 1920x1080 (FHD) dan 3840x2160 (4K TV): teks nomor antrian terbaca tajam dan tidak terjadi layout shift atau text wrap berantakan.
2. **Uji Pemanggilan Suara:**
   - Trigger pemanggilan antrian dari Ruang Dokter -> Layar monitor membunyikan lonceng dan menyuarakan panggilan via TTS.
3. **Uji Privasi Data Pasien:**
   - Memastikan tidak ada nama lengkap tanpa sensor atau nomor telepon yang terpapar di layar publik monitor.
4. **Build Verification:**
   - Jalankan `npm run build` di frontend untuk memastikan tidak ada kesalahan TypeScript dan template compile.
