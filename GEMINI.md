# Antigravity & Gemini Project Directives

Pada setiap sesi kerja atau ketika menerima instruksi baru di repositori **clinic** ini, kamu **WAJIB** membaca dan mematuhi panduan pada dokumen-dokumen berikut sebelum menulis atau merevisi kode:

1. **[AGENT.md](file:///home/alpi/Projects/clinic/AGENT.md):** Standar penulisan kode (Go & Vue 3), filosofi Ponytail (anti over-engineering, solusi paling ringkas dan efektif), dan aturan ketat domain klinik (invarian stok obat, antrian FIFO, single invoice).
2. **[DESIGN.md](file:///home/alpi/Projects/clinic/DESIGN.md):** Panduan sistem desain dan UI/UX (estetika Brandlos / Nordic Modern, anti "AI Slop", warm off-white canvas, bold typography, pill badges, dan micro-interaction snappy).
3. **[CONTEXT.md](file:///home/alpi/Projects/clinic/CONTEXT.md):** Arsitektur sistem (Modular Monolith), definisi entity database, siklus status (state transitions), dan struktur direktori.
4. **[docs/PRD.md](file:///home/alpi/Projects/clinic/docs/PRD.md):** Spesifikasi fungsional fitur, skenario demo untuk presentasi klien, dan batasan cakupan (scope vs non-goals).

### Instruksi Kerja Utama:
* Pertahankan kesederhanaan solusi (prinsip Ponytail).
* Utamakan fungsionalitas yang berjalan mulus saat didemokan di depan klien (dilengkapi Role Switcher dan data *seed* realistis).
* Jangan menambahkan dependensi atau layer abstraksi yang tidak diminta.
* **Kemitraan Kritis & Objektif (Anti "Yes-Man" Agent):**
  * DILARANG sekadar mengiyakan (*yes-man*) atau selalu pro terhadap pendapat/rencana programmer jika keputusan tersebut keliru, kurang tepat, berisiko merusak arsitektur/keamanan, melanggar invarian klinik, atau over-engineering.
  * Jika programmer mengajukan ide atau solusi yang tidak tepat, kamu **WAJIB** bersikap kritis dan objektif: tunjukkan titik kesalahan atau risiko teknisnya secara lugas, lalu sertakan saran serta alternatif solusi terbaik yang lebih aman, ringkas, dan tepat guna.
