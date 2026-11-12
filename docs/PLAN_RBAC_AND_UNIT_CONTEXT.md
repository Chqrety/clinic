# Rencana Kerja: Implementasi RBAC & Konteks Unit Klinik

Dokumen ini merinci arsitektur dan langkah implementasi **Role-Based Access Control (RBAC)** dan **Konteks Unit Klinik (Kidinest vs Semarang Skin Center)** pada antarmuka workstation staf (`AdminPortal.vue`, `AdminSidebar.vue`, dan workstation terkait) dengan backend PostgreSQL.

---

## 1. Analisis & Tujuan Arsitektur

Saat ini seluruh staf yang berhasil login (atau membuka `/admin`) dapat mengakses dan berganti ke semua modul pelayanan secara bebas tanpa pembatasan hak akses. 

### Target yang Hendak Dicapai:
1. **Superadmin Full Access (`SUPERADMIN_OWNER`):**
   * Sebagai representasi dr. Ratna Paramita, M.Kes (Atasan Klinik Utama), akun ini **memiliki akses ke seluruh 7 modul workstation** dan tetap dapat berganti peran secara bebas melalui Role Switcher untuk keperluan supervisi dan presentasi demo ke calon klien.
2. **Role-Based Workstation Filtering untuk Staf Operasional:**
   * Staf non-superadmin (Dokter, Kasir, Apoteker) hanya melihat menu dan modul yang relevan dengan tugas operasionalnya.
   * Menu sensitif *"Manajemen Staf"* (`StaffAdmin.vue`) hanya boleh diakses oleh `SUPERADMIN_OWNER`.
3. **Konteks Unit Klinik Otomatis (`clinic_unit`):**
   * Dokter Poli Anak (`KIDINEST`) langsung masuk ke antrian pasien anak, resep obat anak, dan diagnosa anak.
   * Dokter Estetika (`SSC`) langsung masuk ke antrian pasien estetika medis, tindakan laser, dan krim dermatologi.
   * Kasir dan Apoteker (`ALL`) dapat berpindah/memfilter antrian unit dengan jelas.
4. **Auth Guard & Manajemen Sesi:**
   * Memastikan halaman `/admin` terlindungi dari akses tanpa sesi token JWT yang sah.

---

## 2. Matriks Hak Akses (RBAC Matrix)

| Modul Pelayanan | Key di Frontend | SUPERADMIN_OWNER | DOCTOR | GRO_CASHIER | PHARMACIST |
|---|---|:---:|:---:|:---:|:---:|
| **Kasir & Pendaftaran** | `cashier` | ✅ | ❌ | ✅ | ❌ |
| **Ruang Dokter & E-Resep** | `doctor` | ✅ | ✅ | ❌ | ❌ |
| **Farmasi & Apotek** | `pharmacy` | ✅ | ❌ | ❌ | ✅ |
| **Gudang Obat & Stok Opname** | `inventory` | ✅ | ❌ | ❌ | ✅ |
| **Billing & Tutup Shift** | `billing` | ✅ | ❌ | ✅ | ❌ |
| **Manajemen Staf (Atasan)** | `owner` | ✅ | ❌ | ❌ | ❌ |
| **SATUSEHAT Bridge** | `satusehat` | ✅ | ❌ | ❌ | ❌ |
| **Dropdown Demo Switcher** | Role Switcher | ✅ | ❌ | ❌ | ❌ |

---

## 3. Konteks Unit Klinik pada Workstation Dokter (`DoctorRoom.vue`)

Saat login berhasil, objek user di `authStore.user` menyimpan data:
```typescript
{
  id: 2,
  employee_number: "DOC-KID-01",
  name: "dr. Hendra Pratama, Sp.A",
  role: "DOCTOR",
  clinic_unit: "KIDINEST" // atau "SSC" atau "ALL"
}
```

### Mekanisme Pengikatan Otomatis:
1. **Inisialisasi Dokter Aktif:**
   * Jika user yang login adalah `DOCTOR`, `currentDoctor` otomatis diikat ke profil dokter yang login (nama & `clinic_unit`).
   * Tombol ganti dokter manual di ruang periksa disembunyikan untuk staf dokter biasa, namun tetap ditampilkan jika login sebagai `SUPERADMIN_OWNER`.
2. **Filter Antrian & Resep:**
   * Pemanggilan `api.getDoctorQueues(doctorName, unit)` otomatis menggunakan nama dan unit dari sesi dokter aktif.
   * Katalog obat e-resep otomatis memfilter obat unit yang sesuai (`KIDINEST` untuk anak, `SSC` untuk estetika).
   * ICD-10 presets otomatis menyesuaikan (Batuk/ISPA/Demam/Vaksin untuk Kidinest, Acne/Melasma/Facial untuk SSC).

---

## 4. Perlindungan Rute (Auth Guard) & Sesi

### A. Auth Guard di `App.vue`:
```typescript
// Saat rute adalah 'admin'
if (currentRoute.value === 'admin') {
  if (!auth.isAuthenticated) {
    // Arahkan ke login jika belum ada sesi aktif
    navigateTo('login-kidinest')
  } else {
    // Verifikasi validitas token ke backend
    auth.fetchMe().then((valid) => {
      if (!valid) navigateTo('login-kidinest')
    })
  }
}
```

### B. Mekanisme Logout:
* Tombol logout di `AdminSidebar.vue` memanggil `auth.logout()`.
* Endpoint `POST /api/auth/logout` dipanggil untuk membersihkan cookie HttpOnly backend.
* LocalStorage dibersihkan (`clinic_auth_token`, `clinic_token`, `clinic_auth_user`).
* Mengarahkan user kembali ke landing page atau portal login unit asal.

---

## 5. Rencana Tahapan Eksekusi

- [ ] **Langkah 1: Pembaruan `AdminSidebar.vue`**
  * Tambahkan metadata `allowedRoles` pada definisi daftar menu `roles[]`.
  * Tampilkan hanya modul yang diizinkan untuk role yang sedang login.
  * Hanya tampilkan dropdown demo switcher untuk `SUPERADMIN_OWNER`.
- [ ] **Langkah 2: Sinkronisasi Default Workstation di `LoginView.vue`**
  * Saat login berhasil, set `clinicStore.currentRole` sesuai role utama user:
    * `SUPERADMIN_OWNER` -> default `'owner'`
    * `DOCTOR` -> default `'doctor'`
    * `GRO_CASHIER` -> default `'cashier'`
    * `PHARMACIST` -> default `'pharmacy'`
- [ ] **Langkah 3: Pengikatan Konteks Unit di `DoctorRoom.vue`**
  * Ikat `currentDoctor` ke `auth.user` jika role adalah `DOCTOR`.
  * Kunci antrian dan katalog obat sesuai `auth.user.clinic_unit`.
- [ ] **Langkah 4: Auth Guard di `App.vue` & Verifikasi Sesi**
  * Redirect ke login jika membuka `/admin` tanpa token valid.
- [ ] **Langkah 5: Uji End-to-End dengan 5 Akun Seed PostgreSQL**
  * Login sebagai Owner (`OWNER-01`) -> Verifikasi semua 7 modul dan role switcher tampil.
  * Login sebagai Dokter Anak (`DOC-KID-01`) -> Verifikasi hanya Ruang Dokter & antrian Kidinest.
  * Login sebagai Dokter Kulit (`DOC-SSC-01`) -> Verifikasi hanya Ruang Dokter & antrian SSC.
  * Login sebagai Kasir (`GRO-01`) -> Verifikasi modul Kasir & Billing.
  * Login sebagai Apoteker (`FARMASI-01`) -> Verifikasi modul Farmasi & Gudang Obat.
