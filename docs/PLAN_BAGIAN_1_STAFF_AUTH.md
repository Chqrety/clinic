# Rencana Implementasi: Bagian 1 — Fondasi Akun Staf & Modul Atasan (Staff Management & Security)

Dokumen ini adalah rencana teknis pelaksanaan tahap pertama pengembangan **Smart Clinic Management System (SCMS)** sesuai filosofi **Ponytail** (solusi paling ringkas, tepat sasaran, zero-bloat) dan arsitektur modular Go + Vue 3.

---

## 1. Tujuan & Ruang Lingkup

### A. Sasaran
1. Membangun fondasi autentikasi internal berbasis **Nomor Induk Karyawan (NIK Karyawan / Employee ID)** dan **Password**.
2. Memberikan hak kelola penuh (*Role-Based Access Control / RBAC*) kepada **Atasan/Owner** untuk menambah staf, mengatur peran/unit klinik, dan melakukan reset password langsung dari dashboard web (tanpa keterlibatan IT/database manual).
3. Menyiapkan data benih (*seed*) realistis untuk kebutuhan demo presentasi di depan klien (16–19 November 2026).

### B. Batasan (Out of Scope untuk Bagian 1)
- Tidak ada registrasi mandiri (*self-service registration*) oleh staf publik. Semua akun dibuat oleh Atasan.
- Tidak ada Google OAuth / 2FA SMS / Email verification (sesuai kesepakatan: login cepat menggunakan kredensial internal).
- Tidak ada modul pasien, antrian, maupun resep (dikerjakan pada Bagian 2–5).

---

## 2. Struktur Data & Migrasi PostgreSQL

### A. Tabel `staff_users`
```sql
CREATE TABLE IF NOT EXISTS staff_users (
    id SERIAL PRIMARY KEY,
    employee_number VARCHAR(50) UNIQUE NOT NULL, -- NIK Karyawan resmi (misal: OWNER-01, DOC-KID-01)
    name VARCHAR(150) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN ('SUPERADMIN_OWNER', 'DOCTOR', 'GRO_CASHIER', 'PHARMACIST')),
    clinic_unit VARCHAR(20) NOT NULL CHECK (clinic_unit IN ('KIDINEST', 'SSC', 'ALL')),
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_staff_employee_number ON staff_users(employee_number);
CREATE INDEX IF NOT EXISTS idx_staff_role ON staff_users(role);
CREATE INDEX IF NOT EXISTS idx_staff_unit ON staff_users(clinic_unit);
```

### B. Seed Akun Default Demo (Password: `klinik123`)
| NIK Karyawan | Nama | Peran | Unit | Password Demo |
|---|---|---|---|---|
| `OWNER-01` | dr. Ratna Paramita, M.Kes | `SUPERADMIN_OWNER` | `ALL` | `klinik123` |
| `DOC-KID-01` | dr. Hendra Pratama, Sp.A | `DOCTOR` | `KIDINEST` | `klinik123` |
| `DOC-SSC-01` | dr. Nadia Kusuma, Sp.KK | `DOCTOR` | `SSC` | `klinik123` |
| `GRO-01` | Siti Sarah, A.Md | `GRO_CASHIER` | `ALL` | `klinik123` |
| `FARMASI-01` | Apt. Dimas Wicaksono, S.Farm | `PHARMACIST` | `ALL` | `klinik123` |

---

## 3. Desain Arsitektur Backend (Domain-Driven Design / DDD)

Backend menggunakan arsitektur **Domain-Driven Design (DDD)** yang modular dan terisolasi per bounded context:

```text
backend/internal/
├── shared/                       # Shared Kernel lintas modul
│   ├── database/                 # Postgres connection pool & auto-migrator
│   ├── response/                 # Standar JSON response wrapper ({ success, data, error })
│   └── middleware/               # Auth JWT guard & RequireRole("SUPERADMIN_OWNER")
├── auth/                         # Bounded Context: Autentikasi Staf
│   ├── domain/
│   │   ├── entity.go             # Entity StaffAuth, JWTClaims, AuthRepository interface
│   │   └── errors.go             # Domain errors (ErrInvalidCredentials, ErrStaffInactive)
│   ├── usecase/
│   │   └── auth_usecase.go       # Logika login, verifikasi bcrypt, generate JWT
│   ├── repository/
│   │   └── postgres_repo.go      # Query PostgreSQL: FindByEmployeeNumber
│   └── delivery/http/
│       ├── handler.go            # HTTP Handler login, logout, get me
│       └── routes.go             # Rute /api/auth/*
└── staff/                        # Bounded Context: Manajemen Data Staf (Atasan)
    ├── domain/
    │   ├── entity.go             # Entity StaffUser, Role enum, ClinicUnit enum, StaffRepository interface
    │   └── errors.go             # Domain errors (ErrStaffNotFound, ErrDuplicateEmployeeNum)
    ├── usecase/
    │   └── staff_usecase.go      # Logika tambah staf, update status, dan reset password
    ├── repository/
    │   └── postgres_repo.go      # Query PostgreSQL: CRUD staff_users
    └── delivery/http/
        ├── handler.go            # HTTP Handler CRUD staff & reset password
        └── routes.go             # Rute /api/admin/staff/*
```

### A. Rincian Endpoint API Autentikasi (`internal/auth/delivery/http`)
1. **`POST /api/auth/login`**
   - Payload:
     ```json
     { "employee_number": "OWNER-01", "password": "klinik123" }
     ```
   - Respon (200 OK):
     ```json
     {
       "success": true,
       "data": {
         "token": "jwt-token-string",
         "user": {
           "id": 1,
           "employee_number": "OWNER-01",
           "name": "dr. Ratna Paramita, M.Kes",
           "role": "SUPERADMIN_OWNER",
           "clinic_unit": "ALL"
         }
       }
     }
     ```
   - Validasi: Gagal jika NIK tidak terdaftar, password salah, atau `is_active == false`.
2. **`GET /api/auth/me`**
   - Header: `Authorization: Bearer <token>`
   - Mengambil identitas user yang sedang aktif.
3. **`POST /api/auth/logout`**
   - Menghapus session/token di sisi client.

### B. Rincian Endpoint Kelola Staf Atasan (`internal/staff/delivery/http`)
*Wajib melewati middleware `RequireRole("SUPERADMIN_OWNER")`.*

1. **`GET /api/admin/staff`**
   - Mengambil daftar semua staf klinik.
   - Filter opsional: `?unit=KIDINEST&role=DOCTOR&is_active=true`.
2. **`POST /api/admin/staff`**
   - Payload:
     ```json
     {
       "employee_number": "DOC-KID-02",
       "name": "dr. Anisa S., Sp.A",
       "role": "DOCTOR",
       "clinic_unit": "KIDINEST",
       "initial_password": "klinik123"
     }
     ```
3. **`PUT /api/admin/staff/{id}`**
   - Mengubah data staf (nama, role, unit, status aktif/non-aktif).
4. **`PUT /api/admin/staff/{id}/reset-password`**
   - Payload:
     ```json
     { "new_password": "passwordBaru123" }
     ```
   - Atasan langsung menetapkan password baru tanpa butuh email/link reset.

---

## 4. Desain Tampilan Frontend (Vue 3 + Tailwind CSS)

### A. State Store (`src/stores/auth.ts`)
- Menyimpan: `token`, `currentUser`, status autentikasi, helper `hasRole(roles)`.
- Sinkronisasi dengan `localStorage` agar tidak logout saat refresh halaman.

### B. Halaman Login Staf (`src/views/StaffLogin.vue`)
- Mengikuti **DESIGN.md**: background sand hangat (`bg-[#F9F9F7]`), kartu kontras putih (`bg-white rounded-3xl border border-zinc-200/80 shadow-sm`), pill badge status.
- Dilengkapi **"Demo Quick-Switch"**: Baris pill button di bawah form login (`Owner`, `Dokter Kidinest`, `Dokter SSC`, `Kasir GRO`, `Apoteker`) yang bila diklik otomatis mengisi NIK & Password untuk kelancaran presentasi klien.

### C. Halaman Kelola Staf Atasan (`src/views/StaffAdmin.vue`)
- Layout desktop bersih dengan header ringkasan total staf, dokter aktif, kasir, dan apoteker.
- Tabel staf interaktif: NIK, Nama, Role (badge berwarna), Unit Klinik, Status (Aktif / Nonaktif toggle).
- Aksi:
  - Tombol **"Tambah Karyawan Baru"** (Modal dialog form).
  - Tombol **"Reset Password"** pada baris staf (Modal dialog input password baru).
  - Toggle Aktif/Nonaktif akun instan.

---

## 5. Rencana Langkah Kerja Eksekusi (Checklist DDD)

1. [x] **Step 1: Shared Kernel & Database PostgreSQL**
   - Tambahkan driver PostgreSQL (`github.com/jackc/pgx/v5/stdlib`) dan JWT (`golang-jwt/jwt/v5`) ke `backend/go.mod`.
   - Buat `internal/shared/database/` untuk koneksi Postgres + auto-migrate tabel `staff_users` & seed default akun demo.
   - Buat `internal/shared/response/` (JSON wrapper) dan `internal/shared/middleware/` (JWT verification & Role guard).
2. [x] **Step 2: Bounded Context Auth (`internal/auth/`)**
   - Buat `domain/entity.go` & `domain/errors.go` (definisi StaffAuth & interface AuthRepository).
   - Buat `repository/postgres_repo.go` (query PostgreSQL `staff_users`).
   - Buat `usecase/auth_usecase.go` (verifikasi password bcrypt, cek `is_active`, generate JWT token).
   - Buat `delivery/http/handler.go` & `routes.go` (endpoint `/api/auth/login`, `/api/auth/me`, `/api/auth/logout`).
   - Tulis unit test `usecase/auth_usecase_test.go`.
3. [x] **Step 3: Bounded Context Staff Management (`internal/staff/`)**
   - Buat `domain/entity.go` (StaffUser, Role, ClinicUnit) & `StaffRepository` interface.
   - Buat `repository/postgres_repo.go` (CRUD query `staff_users`).
   - Buat `usecase/staff_usecase.go` (ListStaff, CreateStaff dengan hash password, UpdateStaff, ResetPassword).
   - Buat `delivery/http/handler.go` & `routes.go` (rute `/api/admin/staff`).
   - Tulis unit test `usecase/staff_usecase_test.go`.
4. [x] **Step 4: Frontend Pinia Auth Store & API Client**
   - Buat `frontend/src/stores/auth.ts` (state token, currentUser, role helpers) dan sync localStorage.
5. [x] **Step 5: Frontend Views (StaffLogin & StaffAdmin)**
   - Buat `src/views/LoginView.vue` dengan Brandlos styling + baris pill button demo quick-fill (Owner, Dokter Kidinest, Dokter SSC, Kasir, Apotek).
   - Buat `src/views/StaffAdmin.vue` dengan tabel staf, modal tambah staf, dan modal reset password oleh Atasan.
   - Daftarkan ke `src/views/AdminPortal.vue` dan `AdminHeader.vue`.
6. [x] **Step 6: Verifikasi & Uji Alur End-to-End**
   - Unit test auth & staff usecase lulus 100%.
   - Build backend binary lulus.
   - Build frontend `vue-tsc && vite build` lulus 0 errors.

---

## 6. Skenario Pengujian & Akseptansi
- [x] Login berhasil mengembalikan token JWT dan objek profil user.
- [x] User dengan role `GRO_CASHIER` atau `DOCTOR` yang mencoba akses `/api/admin/staff` mendapatkan status `403 Forbidden`.
- [x] Akun yang diset `is_active = false` tidak dapat login (`401 Unauthorized`).
- [x] Reset password oleh Atasan langsung berlaku seketika pada login berikutnya.
