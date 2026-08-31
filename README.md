# Klontong Net

Marketplace kebutuhan sehari-hari. Banyak penjual, tiap penjual punya beberapa
outlet dengan stok masing-masing, dan pembeli bisa memilih diantar kurir lokal,
dikirim ekspedisi, atau diambil sendiri di toko.

## Isi repo

```
apps/
  backend/     Go + chi + PostgreSQL. Seluruh sistem.
  frontend/    React + TypeScript + Vite. Storefront (pembeli) + dashboard (penjual, admin).
contracts/     OpenAPI 3.0.3. Kontrak API yang mengikat kedua app.
docs/          Arsitektur dan keputusan (lintas-app).
```

Backend dan frontend dikerjakan **paralel dan tidak saling menunggu**. Yang
menyambungkan keduanya adalah `contracts/`: backend mengimplementasikan apa yang
tertulis di sana, frontend men-generate tipenya dari sana dan mengembangkan seluruh
alurnya di atas mock. Kontrak **mengikat kedua sisi** — kalau implementasi perlu
menyimpang, kontraknya yang diubah, di PR yang sama. Baca ADR-015.

## Bentuk yang paling menentukan

Satu keranjang bisa berisi barang dari beberapa penjual sekaligus. Karena itu
**satu pesanan pembeli bukan satu unit pengiriman**:

```
Order (milik pembeli, satu pembayaran)
 ├── Suborder  penjual A, outlet A1, kurir ekspedisi
 ├── Suborder  penjual B, outlet B2, antar lokal
 └── Suborder  penjual C, outlet C1, ambil di toko
```

Tiap suborder punya ongkir, status, dan pencairan dananya sendiri. Ini keputusan
struktural terbesar di project dan menyentuh hampir semua modul. Baca ADR-002 di
[docs/DECISIONS.md](docs/DECISIONS.md) sebelum menulis kode yang berhubungan
dengan pesanan.

## Mulai

Backend:

```bash
make setup
make migrate
make dev          # API di :8080
make worker       # terminal lain — WAJIB
```

Tanpa worker: reservasi stok tidak pernah kedaluwarsa, pembayaran tidak
direkonsiliasi, dan pencairan ke penjual tidak berjalan. Semuanya bergejala sama —
sistem terlihat normal padahal tidak ada yang terjadi.

Frontend (tidak butuh backend jalan):

```bash
make fe-gen       # generate tipe TS dari kontrak
make mock         # mock server dari kontrak di :4010 — terminal lain
make fe-dev       # storefront + dashboard
```

Kontrak:

```bash
make contracts-check   # lint + bundle
make check-all         # gate lintas-app: kontrak + backend + frontend
```

## Dokumentasi

| Dokumen | Isi |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Peta modul, alur checkout, model uang & stok |
| [docs/DECISIONS.md](docs/DECISIONS.md) | Catatan keputusan teknis (ADR) |
| [contracts/README.md](contracts/README.md) | Kontrak API: aturan yang dikunci, cara memakai, cara mengubah |
| [apps/backend/docs/GUIDES.md](apps/backend/docs/GUIDES.md) | Panduan kerja backend: mulai, tambah modul, endpoint, usecase |
| [apps/backend/docs/TASKS.md](apps/backend/docs/TASKS.md) | Urutan 15 fase implementasi backend |
| [apps/frontend/docs/TASKS.md](apps/frontend/docs/TASKS.md) | Urutan 11 fase implementasi frontend |
| [docs/tasks/INTEGRASI.md](docs/tasks/INTEGRASI.md) | Fase lintas-app: lepas mock, sambung ke backend asli |
| [AGENTS.md](AGENTS.md) | Aturan untuk AI coding agent |
