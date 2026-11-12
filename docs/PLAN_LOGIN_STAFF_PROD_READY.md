# Rencana Kerja: Pemisahan Portal Login Staf (Production vs Hidden Demo Mode)

Dokumen ini merinci arsitektur transisi **Portal Login Staf** (`LoginView.vue`) dari prototipe demo menjadi **antarmuka siap produksi (Production-Ready)** dengan mode demo rahasia yang hanya dapat diakses melalui URL parameter khusus oleh tim developer.

---

## 1. Analisis & Identifikasi Masalah Saat Ini

| Aspek | Kondisi Saat Ini (Prototipe) | Masalah / Risiko Keamanan | Standar Produksi (Target) |
|---|---|---|---|
| **Layout Layar** | Split layout 12 kolom (7 kolom daftar akun demo, 5 kolom form) | Terlihat seperti aplikasi dummy; membingungkan staf asli di lapangan | Layout terpusat (*Centered Single-Card*, `max-w-md`) fokus pada autentikasi |
| **Kredensial Default** | `employeeNumber = 'OWNER-01'` & `password = 'klinik123'` ter-prefill | Kerentanan keamanan serius jika rilis ke server publik | Input NIK & Password wajib kosong (`""`) secara default |
| **Bocoran Akun Demo** | Daftar seluruh dokter, kasir, dan atasan dengan password terbuka | Membocorkan NIK staf dan kredensial internal klinik | Nol akun dummy terlihat di layar login produksi |
| **Kebutuhan Demo Klien** | Presentasi ke calon klien butuh login instan tanpa repot ketik NIK | Jika akun demo dihapus total sekarang, presentasi ke klien jadi lambat | Mode demo tetap ada, namun disembunyikan (*hidden*) di balik URL trigger |

---

## 2. Mekanisme Akses URL Hidden Demo

Mode demo **hanya** akan aktif jika browser mendeteksi parameter khusus pada URL:

### A. Pola URL:
* **Mode Normal / Produksi (Default):**
  * Kidinest: `/login-kidinest` atau `/login`
  * Semarang Skin Center: `/login-ssc`
  * *Hasil:* Tampilan form login bersih, input kosong, tidak ada tombol demo.
* **Mode Developer / Pitching Demo (Hidden):**
  * Kidinest: `/login-kidinest?demo=true` atau `/login-kidinest?dev=1` atau `/login-kidinest#demo`
  * Semarang Skin Center: `/login-ssc?demo=true` atau `/login-ssc?dev=1` atau `/login-ssc#demo`
  * *Hasil:* Banner penanda `[Developer Demo Mode]` aktif, panel akun demo 1-klik terbuka, dan input default terisi otomatis.

### B. Implementasi Kode Ponytail (Zero-Dependency):
```typescript
// frontend/src/views/LoginView.vue
const isDemoMode = ref(false)

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  isDemoMode.value = 
    urlParams.has('demo') || 
    urlParams.has('dev') || 
    window.location.hash.includes('demo')

  // Prefill hanya jika mode demo aktif
  if (isDemoMode.value) {
    employeeNumber.value = props.brand === 'kidinest' ? 'DOC-KID-01' : 'DOC-SSC-01'
    password.value = 'klinik123'
  } else {
    employeeNumber.value = ''
    password.value = ''
  }
})
```

---

## 3. Desain Tampilan UI

### A. Tampilan Produksi (Default - `isDemoMode = false`):
* **Kanvas:** Warm off-white (`#FAF7F2` untuk Kidinest, `#FAF5F5` untuk SSC).
* **Kartu Form:** `max-w-md mx-auto` terpusat vertikal dan horizontal.
* **Header Kartu:**
  * Pill badge unit medis (`Kidinest Pediatric Care` / `Semarang Skin Center`).
  * Judul: *"Portal Autentikasi Staf Medis"*.
  * Sub-judul: *"Silakan masukkan Nomor Induk Karyawan (NIK) dan kata sandi operasional Anda."*
* **Input Fields:**
  * Field 1: NIK Internal Karyawan (placeholder: *"Contoh: DOC-001 / ADM-002"*, font monospace tebal).
  * Field 2: Kata Sandi (type: password).
  * Tombol Submit: *"Masuk ke Sistem Workstation"*.
* **Footer Kartu:**
  * Bantuan: *"Kendala akses NIK? Hubungi Administrator / Superadmin Klinik."*

### B. Tampilan Mode Demo (Hidden - `isDemoMode = true`):
* Banner pita di atas: `[DEMO MODE AKTIF • Gunakan URL tanpa ?demo=true untuk tampilan produksi]`.
* Toggle panel atau split accordion di bawah form yang memuat tombol 1-klik `kidinestAccounts` dan `sscAccounts`.

---

## 4. Rencana Kerja Implementasi (Checklist)

- [ ] **Langkah 1:** Tambahkan deteksi query param `?demo=true` / `?dev=1` di `LoginView.vue`.
- [ ] **Langkah 2:** Kosongkan default value `employeeNumber` dan `password` pada mode non-demo.
- [ ] **Langkah 3:** Restrukturisasi layout `LoginView.vue` menjadi centered card saat non-demo, dan expand layout saat mode demo aktif.
- [ ] **Langkah 4:** Tambahkan penanda mode developer yang jelas agar presenter tahu sedang di mode apa.
- [ ] **Langkah 5:** Catat Technical Debt Ledger di `CONTEXT.md` mengenai file dan baris yang harus dihapus saat migrasi ke production hosting.
- [ ] **Langkah 6:** Uji `npm run build` dan pastikan kedua URL (`/login-kidinest` dan `/login-kidinest?demo=true`) bekerja tanpa error.
