# AI Agent Guidelines & Coding Standards

This repository follows the **Ponytail (Lazy Senior Developer)** philosophy: write the minimum amount of robust code that completely solves the problem. No premature abstractions, no unnecessary dependencies, no boilerplate for imaginary future requirements.

---

## 1. Core Operating Principles

1. **Simplicity First (YAGNI):** Build what was requested. If a single function or component does the job, don't build three layers of abstraction around it.
2. **Modular Monolith Discipline:**
   * Backend code belongs to its domain package (`queue`, `pharmacy`, `consultation`, `billing`, `satusehat`). Keep cross-package coupling clean and explicit.
   * Frontend views reflect user roles (`PatientKiosk`, `CashierView`, `DoctorRoom`, `PharmacyRoom`, `InventoryView`).
3. **Design System & UI Fidelity (Anti "AI Slop"):**
   * Setiap kali membuat atau merevisi tampilan antarmuka (UI/UX), **WAJIB** membaca dan mengacu pada [DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md).
   * Terapkan estetika **Brandlos / Nordic Modern**: warm sand background (`bg-[#F9F9F7]`), tipografi sans-serif tebal dan berani (`text-zinc-950`), pill badges bulat penuh (`rounded-full`), kartu membulat halus (`rounded-3xl`), dan sentuhan micro-interaction cepat (150-250ms).
4. **Demo & Presentation Readiness:**
   * Provide a sticky top bar or floating **Role Switcher** in the UI to allow seamless switching between Patient, Cashier, Doctor, Pharmacist, and Admin views during client demos.
   * Include pre-seeded realistic clinic data (common Indonesian drugs, ICD-10 codes, NIK samples, batches with upcoming expiry).

---

## 2. Backend Guidelines (Go)

* **Folder Architecture (Package by Feature):** DILARANG menggunakan Global Layered Architecture (`controllers/`, `services/`, `repositories/` global). Gunakan pengelompokan modular domain di `internal/<domain>/` (misal `internal/pharmacy/` berisi `handler.go`, `service.go`, `repository.go`).
* **Code Style:** Idiomatic Go. Return errors explicitly (`if err != nil`). Avoid deep nesting.
* **HTTP Handlers:** Simple REST handlers returning standard JSON structures:
  ```json
  { "success": true, "data": { ... }, "error": null }
  ```
* **Database:** SQLite using modern lightweight drivers. Always enforce foreign keys and transactions when handling financial data or inventory mutations (e.g. `reserved_stock` -> `physical_stock`).
* **Concurrency & Inventory Safety:** When reserving or dispensing drugs, use database transactions or atomic updates (`UPDATE drug_batches SET reserved_stock = reserved_stock + ? WHERE available_stock >= ?`) to prevent race conditions.

---

## 3. Frontend Guidelines (Vue 3 + Tailwind CSS)

* **Folder Architecture (Role/View-Based):** DILARANG menggunakan Atomic Design (`atoms/`, `molecules/`, `organisms/`). Gunakan `src/views/` per peran pengguna (`PatientKiosk.vue`, `DoctorRoom.vue`, dll.) dan `src/components/ui/` hanya untuk komponen yang benar-benar dipakai berulang.
* **Composition API:** Use `<script setup lang="ts">` exclusively.
* **Component Granularity:** Keep components cohesive. Do not extract one-time micro-components unless they are genuinely reused across multiple views.
* **State Management:** Use Pinia stores for shared global state (e.g. current active user role, notification bell, live queue updates). Local view state stays in `ref()` / `reactive()`.
* **Mobile-First for Patient Kiosk:** Ensure `PatientKiosk.vue` looks and feels like a native mobile web app (responsive, touch-friendly QR/buttons).
* **Desktop-First for Staff Rooms:** Doctor, Pharmacy, Cashier, and Inventory views should maximize screen real-estate with clear status badges and keyboard accessibility.

---

## 4. Domain Invariants (Never Break These)

1. **Available Stock Formula:**  
   `available_stock = physical_stock - reserved_stock`. An order/prescription must never reserve more than `available_stock`.
2. **Queue FIFO Rule:**  
   Queue numbers must be assigned sequentially based on successful payment completion timestamp (`PAID`), not registration draft time.
3. **Single Invoice for Patients:**  
   The patient-facing receipt must display a single consolidated breakdown. Internal margin/tax allocations must be kept strictly within management reporting views.
4. **SATUSEHAT Compliance:**  
   Every encounter and prescription must format its data compatible with Kemenkes FHIR R4 structure (`Encounter`, `Condition`, `MedicationRequest`).
