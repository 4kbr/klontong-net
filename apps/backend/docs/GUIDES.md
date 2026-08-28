# Panduan — backend

Resep lengkap ada di [`../../docs/GUIDES.md`](../../docs/GUIDES.md). File ini
memuat hal-hal yang khas backend dan sering dibutuhkan cepat.

---

## Menjalankan

```bash
make dev        # API, hot reload, :8080
make worker     # terminal lain — WAJIB
```

Tanpa worker: reservasi stok tidak lepas, pesanan tidak kedaluwarsa, dana tidak
cair. Ketiganya gagal secara senyap.

---

## Struktur direktori

```
cmd/
  api/        HTTP server
  worker/     pekerjaan latar (proses terpisah)
  migrate/    runner migrasi programatik
  seed/       data contoh
internal/
  platform/   config, db, http, auth, money, event, klien luar
  app/        composition root: registry, deps, runner worker
  modules/    15 modul bisnis
migrations/   SQL golang-migrate
```

---

## Isi platform

| Package | Peran |
|---|---|
| `config` | satu-satunya pembaca `os.Getenv`; validasi saat start |
| `money` | **`Amount` int64, `BasisPoints`, `Distribute`** — baca ini dulu |
| `postgres` | pool, `TxManager`, penerjemah error |
| `httpx` | router, envelope, error, decode, validasi, **idempotency** |
| `middleware` | requestid, log, recover, cors, auth, role, ratelimit |
| `auth` | JWT, context identitas, hashing password |
| `eventbus` + `outbox` | event lintas modul dengan jaminan konsistensi |
| `httpclient` | klien HTTP bersama dengan retry dan redaksi log |
| `storage` | object storage, presigned upload untuk foto produk |
| `id` | uuid + **nomor pesanan yang dibaca manusia** |

`money` adalah package yang paling menentukan di sini. Baca `money.go` dan
ADR-005 sebelum menulis apa pun yang menghitung.

---

## Pemetaan modul ke tabel

| Modul | Prefix tabel |
|---|---|
| `identity` | `identity_` |
| `customer` | `customer_` |
| `seller` | `seller_` |
| `catalog` | `catalog_` |
| `inventory` | `inventory_` |
| `pricing` | `pricing_` |
| `cart` | `cart_` |
| `order` | `order_` |
| `payment` | `payment_` |
| `fulfillment` | `fulfillment_` |
| `promotion` | `promotion_` |
| `settlement` | `settlement_` |
| `review` | `review_` |
| `notification` | `notification_` |
| `audit` | `audit_events` |

Tidak ada foreign key lintas modul. Relasi disimpan sebagai `uuid` biasa;
integritasnya dijaga di layer aplikasi.

---

## Urutan dependensi modul

```
identity
  ├── customer
  └── seller
        └── catalog
              ├── inventory
              └── pricing
                     └── cart ──▶ order ◀── promotion, fulfillment, payment
                                    │
                                    ├── settlement
                                    └── review
```

`order` bergantung ke paling banyak modul. **Tidak ada modul yang boleh
bergantung balik ke `order`** kecuali lewat event. Kalau muncul kebutuhan itu,
berhenti dan pikirkan ulang.

---

## Test yang wajib ada

| Test | Kenapa |
|---|---|
| Properti `money.Distribute` | selisih rupiah tidak bisa dijelaskan ke penjual |
| Checkout bersamaan stok terakhir | tepat satu harus berhasil |
| Idempotensi checkout & webhook | satu niat beli = satu pesanan |
| Refund parsial | satu penjual menolak, dua lainnya jalan |
| Buku besar seimbang | debit = kredit setelah rangkaian transaksi |
| Konversi satuan | 2 dus isi 40 mengurangi 80 pcs |

Semuanya tidak akan ketahuan dari klik manual.

```bash
make test              # unit, cepat
make test-integration  # butuh docker
make cover
```

---

## Migrasi

```bash
make migrate-create name=xxx
make migrate-up
make migrate-version
make migrate-force version=N   # saat dirty
```

Kolom uang `bigint`. Kolom kuantitas `numeric(14,3)`. Persentase `int` (basis
poin).
