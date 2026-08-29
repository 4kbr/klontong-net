# Klontong Net

Marketplace kebutuhan sehari-hari. Banyak penjual, tiap penjual punya beberapa
outlet dengan stok masing-masing, dan pembeli bisa memilih diantar kurir lokal,
dikirim ekspedisi, atau diambil sendiri di toko.

## Isi repo

```
apps/
  backend/     Go + chi + PostgreSQL. Seluruh sistem.
  frontend/    Belum dikerjakan.
docs/          Arsitektur dan keputusan (lintas-app).
```

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

```bash
make setup
make migrate
make dev          # API di :8080
make worker       # terminal lain — WAJIB
```

Tanpa worker: reservasi stok tidak pernah kedaluwarsa, pembayaran tidak
direkonsiliasi, dan pencairan ke penjual tidak berjalan. Semuanya bergejala sama —
sistem terlihat normal padahal tidak ada yang terjadi.

## Dokumentasi

| Dokumen | Isi |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Peta modul, alur checkout, model uang & stok |
| [docs/DECISIONS.md](docs/DECISIONS.md) | Catatan keputusan teknis (ADR) |
| [apps/backend/docs/GUIDES.md](apps/backend/docs/GUIDES.md) | Panduan kerja backend: mulai, tambah modul, endpoint, usecase |
| [AGENTS.md](AGENTS.md) | Aturan untuk AI coding agent |
