# Rencana Implementasi: Alur Check-In Reservasi Online & Pusat Kendali Pemanggilan GRO
## Solusi Anti-Miskomunikasi Antara Pasien Online, Meja GRO, dan Ruang Dokter

Dokumen ini memuat spesifikasi alur kerja, arsitektur data, dan rancangan antarmuka untuk fitur **Konfirmasi Kehadiran (Check-In) Reservasi Online** serta **Pusat Kendali Antrian Meja Depan (GRO Command Center)** untuk **Kidinest Pediatric Care** dan **Semarang Skin Center (SSC)**.

---

## 1. Latar Belakang & Masalah (*Problem Statement*)

1. **Tingginya Pasien Reservasi Online:**  
   Banyak pasien mendaftar dan memilih jadwal dokter dari rumah. Namun, dokter di dalam ruang periksa tidak bisa melihat keluar ruang tunggu, sehingga tidak tahu apakah pasien sudah tiba di klinik atau masih dalam perjalanan.
2. **Kebutuhan Kontrol di Meja Depan (GRO):**  
   Manajemen klinik menghendaki kendali alur antrian dipegang oleh staf meja depan (*Guest Relations Officer* / Resepsionis).
3. **Risiko Miskomunikasi:**
   - Jika GRO memanggil pasien tanpa tahu kondisi di dalam, pasien bisa masuk saat dokter belum selesai membersihkan alat / menulis rekam medis.
   - Jika dokter memanggil langsung tanpa konfirmasi kedatangan, dokter akan sering memanggil "hantu" (pasien online yang terlambat atau batal datang).

---

## 2. Alur Kerja Terpadu (*End-to-End Workflow*)

```
[ PASIEN DI RUMAH ]
       │
       ▼ (1) Booking online: Pilih Poli, Dokter, Tanggal, & Jam Sesi
[ Status: BOOKED / RESERVED ] (Tidak muncul di daftar siap panggil dokter)
       │
       ▼ (2) Pasien tiba fisik di klinik, melapor ke meja GRO
[ MEJA RESEPSIONIS / GRO ]
       │
       ├─► GRO cari data / scan tiket -> Klik [ ✅ Check-In Hadir ]
       │   Status berubah menjadi: WAITING (Pasien resmi duduk di ruang tunggu)
       │
       ▼ (3) GRO memantau Kartu Status Dokter
┌──────────────────────────────────────────────────────────────┐
│ Jika Dokter 🟡 SEDANG MEMERIKSA : Tombol Panggil Terkunci   │
│ Jika Dokter 🟢 KOSONG / SIAP     : Tombol Panggil AKTIF      │
└──────────────────────────────────────────────────────────────┘
       │
       ▼ (4) GRO (atau Dokter) klik [ 📢 PANGGIL KE RUANGAN ]
[ Status: CALLING ] -> Bel berbunyi di TV Monitor Antrian
       │
       ▼ (5) Pasien masuk, Dokter klik [ Mulai Periksa ]
[ Status: IN_EXAM ] -> Ruangan otomatis terkunci (KUNING / SIBUK)
       │
       ▼ (6) Dokter selesai & kirim resep
[ Status: PHARMACY_WAITING ] -> Ruang Dokter OTOMATIS KEMBALI KOSONG (HIJAU / SIAP)
```

---

## 3. Spesifikasi Tampilan GRO Workstation (*Anti-Miskom UI*)

Fitur ini ditempatkan di workstation staf GRO / Kasir dengan 2 modul utama:

### Modul A: Quick Check-In Reservasi Online
Form pencarian cepat di bagian atas meja GRO:
* **Kolom Input Cepat:** Pencarian berbasis Nama Pasien, Nomor WhatsApp, atau Kode Booking (misal: `BK-0812`).
* **Daftar Booking Hari Ini:** Menampilkan pasien yang terjadwal hari ini dengan tab filter:
  - *Belum Tiba (`BOOKED`)*
  - *Sudah di Ruang Tunggu (`WAITING`)*
  - *Sedang Diperiksa (`IN_EXAM`)*
* **Tombol Aksi 1-Klik:** `[ ✅ Konfirmasi Hadir / Check-In ]`  
  *Efek:* Menerbitkan nomor antrian fisik (misal: `SSC-002`) dan langsung mengupdate daftar tunggu dokter.

---

### Modul B: Kartu Status Ruang Poli Dokter (*Doctor Room Monitor*)
GRO dapat melihat ketersediaan tiap dokter secara *real-time* tanpa perlu meninggalkan meja loket:

#### 1. Saat Dokter Siap (Kosong):
```text
┌──────────────────────────────────────────────────────────────────┐
│ RUANG POLI 1 • dr. Dian Paramita, Sp.D.V.E (SSC)                │
│ Status: 🟢 SIAP / RUANGAN KOSONG                                 │
│ Pasien Terakhir: SSC-001 (Selesai 14:15 WIB)                    │
│                                                                  │
│ Antrian Terdepan di Ruang Tunggu:                                │
│ ▶ [SSC-002] Nadia Saphira — Menunggu 8 menit                    │
│                                                                  │
│ [ 📢 PANGGIL KE RUANG DOKTER ]     [ ⏭️ Toleransi / Lewati ]     │
└──────────────────────────────────────────────────────────────────┘
```

#### 2. Saat Dokter Sedang Memeriksa Pasien (Sibuk):
```text
┌──────────────────────────────────────────────────────────────────┐
│ RUANG POLI 1 • dr. Dian Paramita, Sp.D.V.E (SSC)                │
│ Status: 🟡 SEDANG MEMERIKSA PASIEN (Berjalan 12 menit)           │
│ Pasien di Dalam: [SSC-002] Nadia Saphira                        │
│                                                                  │
│ Pasien Berikutnya di Ruang Tunggu:                               │
│ • [SSC-003] Jessica Mila                                        │
│                                                                  │
│ [ 🔒 RUANG SIBUK — MENUNGGU DOKTER SELESAI ]                     │
└──────────────────────────────────────────────────────────────────┘
```
*(Tombol panggil otomatis dinonaktifkan sehingga GRO tidak bisa salah memanggil).*

---

## 4. Penyesuaian Teknis & Database (Prinsip Ponytail — Tanpa Tabel Baru)

Kita memanfaatkan struktur tabel dan *state machine* yang sudah ada agar tidak menambah dependensi rumit:

### A. State Transition Antrian:
1. `BOOKED`: Pasien reservasi online yang belum tiba di lokasi.
2. `WAITING`: Pasien fisik sudah berada di ruang tunggu (lolos check-in GRO / walk-in).
3. `CALLING`: Pasien sedang dipanggil via TV Monitor.
4. `IN_EXAM`: Pasien berada di dalam ruang dokter (ruangan dokter terdeteksi sibuk).
5. `PHARMACY_WAITING`: Dokter selesai periksa, ruangan dokter terdeteksi kosong kembali.

### B. Endpoint API Tambahan / Penyesuaian:
1. `GET /api/queues/reservations?date=today`  
   Mengambil seluruh daftar reservasi online yang belum check-in.
2. `POST /api/queues/{id}/checkin`  
   Mengubah status dari `BOOKED` ke `WAITING` dan menetapkan nomor antrian resmi FIFO.
3. `POST /api/queues/{id}/call`  
   Endpoint panggil terpadu (dapat dipicu oleh GRO maupun Dokter) yang mengirim sinyal WebSocket / SSE ke TV Monitor.
4. `GET /api/consultations/rooms/status`  
   Menghitung status tiap dokter: jika ada antrian berstatus `IN_EXAM` untuk dokter tersebut maka `BUSY`, jika tidak ada maka `READY`.

---

## 5. Rencana Eksekusi Bertahap

* **Fase 1 (Backend & State):**
  - Pastikan query daftar tunggu dokter hanya menampilkan antrian berstatus `WAITING` (pasien yang sudah hadir).
  - Tambahkan endpoint konfirmasi check-in di backend.
* **Fase 2 (Frontend GRO Workstation):**
  - Tambahkan tab *"Check-In Reservasi"* dan kartu *"Status Ruangan Dokter"* di antarmuka staf meja depan (`CashierView.vue` / GRO portal).
* **Fase 3 (TV Monitor & Audio):**
  - Pastikan event panggil dari GRO membunyikan lonceng dan menampilkan nomor panggilan di TV Monitor secara realtime.

---

### Aturan Batasan Operasional (*Business Invariants*):
- Pasien yang belum di-check-in oleh GRO **dilarang keras** masuk ke daftar panggil dokter.
- Pasien yang dipanggil 3 kali tanpa respon dialihkan ke antrian toleransi (*Skipped Queue*) dan dapat dipulihkan kembali oleh GRO saat pasien melapor.
