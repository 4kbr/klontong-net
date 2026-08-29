# Panduan Kerja

Resep praktis. Setiap resep menyebut **file apa yang dibuat** dan **file apa yang
diubah**, karena itu bagian yang paling sering terlewat.

0. [Peta backend](#0-peta-backend)
1. [Mulai dari nol](#1-mulai-dari-nol)
2. [Urutan pengerjaan](#2-urutan-pengerjaan)
3. [Menjalankan sehari-hari](#3-menjalankan-sehari-hari)
4. [Menambah endpoint](#4-menambah-endpoint)
5. [Menambah usecase](#5-menambah-usecase)
6. [Menambah entity + tabel](#6-menambah-entity--tabel)
7. [Menambah modul](#7-menambah-modul)
8. [Menambah event](#8-menambah-event)
9. [Menambah pekerjaan worker](#9-menambah-pekerjaan-worker)
10. [Bekerja dengan uang](#10-bekerja-dengan-uang)
11. [Bekerja dengan stok](#11-bekerja-dengan-stok)
12. [Menambah penyedia luar](#12-menambah-penyedia-luar)
13. [Menulis migrasi](#13-menulis-migrasi)
14. [Menambah konfigurasi](#14-menambah-konfigurasi)
15. [Menulis test](#15-menulis-test)
16. [Troubleshooting](#16-troubleshooting)

---

## 0. Peta backend

Referensi cepat sebelum masuk ke resep.

### Struktur direktori

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

### Isi platform

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

### Pemetaan modul ke tabel

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

### Urutan dependensi modul

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

## 1. Mulai dari nol

```bash
make setup       # tool, .env, docker
make migrate
make dev         # API di :8080
```

Di terminal terpisah:

```bash
make worker      # WAJIB
make tunnel      # kalau menguji webhook pembayaran
```

Cek: `curl localhost:8080/healthz` balas 200.

**Worker itu wajib, bukan opsional.** Tanpa worker: reservasi stok tidak pernah
lepas, pesanan tidak pernah kedaluwarsa, dana penjual tidak pernah matang, dan
pencairan tidak pernah jalan. Semuanya bergejala sama — sistem terlihat normal
padahal tidak ada yang terjadi. Ini penyebab kebingungan nomor satu saat
development.

Untuk menguji alur pembayaran tanpa tunnel dan akun sandbox, pakai
`gateway_noop.go` — implementasi palsu yang bisa dipicu manual. Kerjakan itu
lebih awal; tanpanya menguji pesanan jadi menyebalkan dan orang berhenti
mengetesnya.

---

## 2. Urutan pengerjaan

Berurutan. Setiap tahap menghasilkan sesuatu yang bisa dicoba.

| # | Tahap | Selesai kalau |
|---|---|---|
| 1 | platform: config, db, TxManager, errs, httpx, **money** | `/healthz` hijau, test `money.Distribute` lulus |
| 2 | `identity` + `customer` | bisa daftar, login, tambah alamat |
| 3 | `seller` + outlet | bisa buka toko dan menambah outlet berkoordinat |
| 4 | `catalog`: produk, varian, satuan | bisa membuat produk dengan varian dus/renceng/pcs |
| 5 | `pricing`: harga + tier grosir | harga per satuan dan tier tampil benar |
| 6 | `inventory`: stok, reservasi, mutasi | stok per outlet benar, konversi satuan benar |
| 7 | `cart` + pemilihan outlet | keranjang 3 penjual terkelompok dengan benar |
| 8 | `fulfillment`: quote 3 metode | ongkir per suborder tampil, pickup gratis |
| 9 | `promotion`: voucher + pembagian diskon | voucher terbagi ke suborder, jumlahnya persis |
| 10 | `order`: preview + place | pesanan multi-penjual terbentuk dengan suborder |
| 11 | `payment` + webhook + rekonsiliasi | pesanan jadi `paid` lewat webhook palsu |
| 12 | siklus suborder: konfirmasi → kirim → terima | stok ter-commit saat kirim |
| 13 | outbox + `audit` + `notification` | setiap aksi tercatat, notifikasi terkirim |
| 14 | `settlement`: buku besar + pencairan | rekonsiliasi selisihnya nol |
| 15 | `review` | ulasan hanya dari pembelian yang selesai |

**Jangan lompat ke tahap 10 sebelum 6 dan 9 selesai.** Checkout memanggil
keduanya, dan menulis checkout di atas stok atau diskon yang belum benar berarti
men-debug dua hal sekaligus.

Tahap 1 dan 14 layak dikerjakan dengan test lebih dulu. Keduanya menyentuh uang,
dan bug di sana mahal.

---

## 3. Menjalankan sehari-hari

```bash
make up          # infra
make dev         # API
make worker      # terminal lain
```

Kalau menguji webhook pembayaran sungguhan, `make tunnel` lalu daftarkan URL-nya
di dashboard gateway. Untuk pengembangan biasa, cukup gateway palsu.

---

## 4. Menambah endpoint

Contoh: `POST /api/v1/seller/orders/{suborderID}/hold`.

**Buat / ubah, berurutan:**

1. `internal/transport/rest/request.go` — struct + `Validate()`
2. `internal/transport/rest/<area>_handler.go` — method handler
3. `internal/transport/rest/routes.go` — daftarkan rute
4. `internal/transport/rest/response.go` — kalau bentuknya baru

**Tentukan dulu areanya:** `/api/v1` (pembeli), `/api/v1/seller`, `/api/v1/admin`,
atau `/webhook`. Salah area berarti salah model akses.

**Checklist:**

- [ ] Rute di area yang benar dengan middleware yang benar
- [ ] `userID` dari context, **tidak pernah** dari body
- [ ] Kepemilikan diperiksa di usecase, bukan cukup mengandalkan peran
- [ ] Error lewat `httpx.Error`
- [ ] Tidak ada aturan bisnis di handler
- [ ] Endpoint yang menciptakan sesuatu bernilai uang: di balik middleware
      `Idempotent`

---

## 5. Menambah usecase

**Buat** `internal/app/<nama>_usecase.go`.
**Ubah** `dto.go`, `domain/*.go`, `domain/repository.go`, `infra/*_repository.go`,
`events.go` bila menerbitkan event.

**Kerangka, urutannya selalu sama:**

```go
func (s *Service) HoldSuborder(ctx context.Context, in HoldInput) (*domain.Suborder, error) {
    // 1. IZIN — selalu paling awal
    sub, err := s.requireSuborderSeller(ctx, in.SuborderID)
    if err != nil { return nil, err }

    // 2. VALIDASI BISNIS (transisi status, dsb)

    // 3. TRANSAKSI — semua tulis, termasuk outbox
    err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
        // ...
        // return s.outbox.Save(ctx, evt)
        return nil
    })
    ...
}
```

**Checklist:**

- [ ] Izin di baris pertama
- [ ] Semua penulisan dalam satu `WithinTx`
- [ ] `outbox.Save` **di dalam** transaksi itu
- [ ] **Tidak ada panggilan jaringan di dalam transaksi** — gateway, kurir, dan
      penyedia transfer dipanggil di luar atau diserahkan ke worker
- [ ] Transisi status lewat `Transition()`, bukan menyetel field langsung
- [ ] `s.clock.Now()`, bukan `time.Now()`
- [ ] Uang memakai `money.Amount`, tidak ada float

---

## 6. Menambah entity + tabel

**Buat:**

1. `migrations/0000XX_create_<modul>_<nama>.up.sql` + `.down.sql`
2. `internal/modules/<modul>/internal/domain/<entity>.go`
3. `internal/modules/<modul>/internal/infra/<entity>_repository.go`
4. `internal/modules/<modul>/internal/app/<entity>_usecase.go`
5. `internal/modules/<modul>/internal/transport/rest/<entity>_handler.go`

**Ubah:**

6. `domain/repository.go` — interface baru
7. `domain/errors.go` — error khusus
8. `app/service.go` — field repo + parameter konstruktor
9. `infra/mapper.go` — fungsi scan
10. `transport/rest/routes.go`, `request.go`, `response.go`
11. **`module.go`** — rakit repository baru dan teruskan ke `NewService`

Nomor 11 paling sering terlupa. Gejalanya: nil pointer panic saat endpoint
pertama dipanggil.

---

## 7. Menambah modul

```
internal/modules/<nama>/
├── module.go
├── port.go
├── events.go
└── internal/
    ├── domain/     { entity.go, errors.go, repository.go }
    ├── app/        { service.go, *_usecase.go, handlers.go, dto.go, port_adapter.go }
    ├── infra/      { *_repository.go, mapper.go }
    └── transport/rest/ { routes.go, handler.go, request.go, response.go }
```

**Ubah:** `migrations/` (prefix `<nama>_`), `internal/app/registry.go`,
`../../../docs/ARCHITECTURE.md` tabel modul, `../../../docs/DECISIONS.md` ADR yang
menjelaskan kenapa ini modul terpisah.

**Sebelum menulis kode, jawab tiga pertanyaan:**

1. **Tabel apa yang dimiliki modul ini sendirian?** Kalau jawabannya "tidak ada,
   dia hanya membaca tabel order" — ini bukan modul, ini usecase di `order`.
2. **Apa isi `port.go`-nya?** Kosong dan tidak ada yang memanggil berarti ia
   konsumen murni. Itu bagus, konfirmasikan saja.
3. **Apakah ia butuh transaksi bersama modul lain?** Kalau ya, batasnya salah —
   gabungkan.

---

## 8. Menambah event

**Ubah `internal/modules/<modul>/events.go`:**

```go
const EventSuborderHeld = "order.suborder.held"

type SuborderHeldPayload struct {
    OrderID, SuborderID, SellerID, OutletID uuid.UUID
    Reason string
    At     time.Time
}
```

**Ubah usecase yang menerbitkannya**, di dalam `WithinTx`:

```go
evt, err := eventbus.New(order.EventSuborderHeld, "suborder", sub.ID, payload)
if err != nil { return err }
return s.outbox.Save(ctx, evt)
```

**Aturan payload:**

- Event suborder **wajib** memuat `SellerID` dan `OutletID`. Tanpa itu setiap
  consumer harus query balik, dan kopling yang mau dihindari kembali lewat pintu
  belakang.
- Sertakan cukup data agar consumer tidak perlu memanggil balik.
- Payload adalah kontrak. Menambah field aman; menghapus atau mengubah tipe tidak.

Penamaan: `<modul>.<agregat>.<aksi lampau>`.

**Untuk consumer**, ubah `internal/app/handlers.go` dan `module.go`
(`RegisterSubscriptions`), lalu pastikan itu dipanggil di `registry.go`. Kalau
lupa memanggilnya, **tidak ada error apa pun** — handlernya cuma tidak pernah
jalan.

**Idempotensi wajib.** Di modul `settlement` ini bukan soal kerapian: event yang
diproses dua kali berarti penjual dibayar dua kali. Pakai unique index pada
`event_id` dan `ON CONFLICT DO NOTHING`.

**Perhatikan levelnya.** Peristiwa penting terjadi di suborder, bukan order:

| Consumer | Mendengarkan | Bukan |
|---|---|---|
| `inventory` commit stok | `SuborderShipped` | `OrderPaid` |
| `settlement` matangkan dana | `SuborderCompleted` | `OrderPaid` |

---

## 9. Menambah pekerjaan worker

**Ubah:**

1. `internal/modules/<modul>/internal/app/<nama>_usecase.go` — logikanya
2. `internal/modules/<modul>/module.go` — `RegisterWorkers`
3. `internal/app/registry.go` — pastikan dipanggil

**Aturan:**

- Berebut baris → `FOR UPDATE SKIP LOCKED`
- Tabel besar → per batch dengan jeda, bukan satu operasi raksasa
- **Menyentuh uang → kunci terdistribusi di Redis + idempotency key di database.**
  Redis bisa hilang; unique index tidak.
- Idempoten: dijalankan dua kali tidak boleh melepas stok dua kali atau membayar
  penjual dua kali.

Beri metrik untuk setiap pekerjaan. Semua kegagalan worker di sistem ini bersifat
senyap — tidak ada yang tahu sampai ada yang mengeluh.

---

## 10. Bekerja dengan uang

Baca ADR-005 dulu. Ringkasnya:

- **`money.Amount` (int64 rupiah).** Tidak ada `float`, tidak pernah.
- **Persentase sebagai basis poin** `int`. 250 = 2,5%.
- **Pembagian pakai `money.Distribute`**, yang menjamin jumlah hasilnya persis
  sama dengan total.
- **Pembulatan komisi ke bawah.** Marketplace tidak boleh mengambil lebih dari
  haknya karena pembulatan.
- **JSON mengirim angka**, bukan string terformat.

Kalau kamu menulis `float64` di dekat perhitungan uang, berhenti. Satu saja sudah
cukup untuk membuat invoice yang jumlahnya tidak cocok, dan bug seperti itu baru
ketahuan saat rekonsiliasi akhir bulan.

Setiap kali menambah perhitungan yang membagi uang, tanya: **kalau ada sisa
pembulatan, ke mana perginya?** Jawaban "dibuang" selalu salah.

---

## 11. Bekerja dengan stok

Baca ADR-003 dan ADR-006. Ringkasnya:

- Semua kuantitas di `inventory.Port` dalam **satuan dasar**. Pemanggil yang
  mengonversi lewat `Variant.ToBaseQuantity()`.
- Stok ditahan saat **checkout**, di-commit saat **kirim**, dilepas saat **batal**.
- `Reserve`/`Commit`/`Release` dipanggil **di dalam transaksi milik pemanggil**,
  memakai `ExecutorFrom`.
- Kunci baris stok dalam **urutan konsisten** (`ORDER BY outlet_id, variant_id`)
  sebelum `FOR UPDATE`. Urutan berbeda antar transaksi = deadlock.
- **Setiap perubahan `on_hand` menulis `Movement`.** Tanpa kecuali. Itu yang
  menjawab "kok stok saya berkurang 12".

---

## 12. Menambah penyedia luar

Contoh: gateway pembayaran kedua.

**Ubah:**

1. `internal/modules/payment/internal/domain/gateway.go` — kalau kontraknya perlu
   diperluas
2. **Buat** `internal/modules/payment/internal/infra/gateway_<nama>.go`
3. `internal/platform/config/config.go` — kredensial baru
4. `.env.example`
5. `internal/app/registry.go` — pilih implementasi berdasarkan config

**Wajib ada di setiap integrasi:**

- Timeout ketat, terutama untuk yang dipanggil di jalur checkout
- Penerjemahan error ke error domain kita dengan `Retryable` diisi benar
- Verifikasi signature webhook dengan perbandingan **waktu-konstan**
- Idempotency key untuk operasi yang menciptakan sesuatu
- Implementasi palsu untuk development

**Jangan** memanggil penyedia dari modul lain. Semua akses lewat interface domain
di modul yang bersangkutan.

---

## 13. Menulis migrasi

```bash
make migrate-create name=add_hold_to_suborders
```

`make migrate-up` adalah alias `make migrate`.

**Aturan:**

- Selalu tulis `.down.sql` yang benar-benar bekerja.
- Satu migrasi = satu perubahan logis.
- Migrasi yang sudah di-commit tidak boleh diedit. Buat yang baru.
- Prefix tabel sesuai modul pemiliknya.
- Buat index bersamaan dengan tabelnya.
- Kolom uang bertipe `bigint`, **bukan** `numeric` atau `float`.
- Kolom kuantitas bertipe `numeric(14,3)` — kuantitas boleh pecahan, uang tidak.

**Perubahan berbahaya:**

| Perubahan | Cara aman |
|---|---|
| Tambah kolom `NOT NULL` | nullable → isi data → set `NOT NULL` |
| Ganti nama kolom | kolom baru → tulis ganda → migrasi → hapus lama |
| Index di tabel besar | `CREATE INDEX CONCURRENTLY` (di luar transaksi migrate) |
| Hapus kolom | pastikan tidak ada kode yang menyebutnya, deploy, baru hapus |

**Dirty:**

```bash
make migrate-version
make migrate-force version=5
make migrate
```

---

## 14. Menambah konfigurasi

Ubah keempatnya: `.env.example`, `internal/platform/config/config.go` (field +
baca + **validasi**), `../../../docker-compose.dev.yml` bila perlu, tempat pemakaian
(terima lewat `Config`, jangan `os.Getenv`).

Kalau nilai yang salah membuat sistem tidak aman atau tidak masuk akal,
**validasi saat start**. Contoh yang sudah ada: `STOCK_RESERVATION_TTL` harus
lebih pendek dari `PAYMENT_EXPIRY` — kalau terbalik, ada jendela waktu stok
tertahan untuk pesanan yang sudah mati.

---

## 15. Menulis test

```bash
make test              # unit, cepat
make test-integration  # butuh docker
make cover             # laporan coverage
```

**Yang wajib ada di sistem ini:**

**Test properti `money.Distribute`.** Input acak, jumlah hasil selalu sama dengan
total, tidak ada bagian negatif. Ini fungsi yang membagi voucher ke suborder dan
diskon ke tiap barang; salah di sini berarti selisih rupiah yang tidak bisa
dijelaskan.

**Checkout bersamaan untuk stok terakhir.** Dua goroutine, stok tinggal satu,
tepat satu berhasil.

**Idempotensi.** Checkout kunci sama dua kali → satu pesanan. Webhook sama dua
kali → satu perubahan. Event settlement dua kali → satu jurnal.

**Refund parsial.** Satu penjual menolak, dua lainnya jalan, jumlah refund tepat,
order induk tidak batal.

**Buku besar seimbang.** Setelah rangkaian transaksi, debit sama dengan kredit
dan saldo tiap akun sesuai harapan.

**Konversi satuan.** Beli 2 dus isi 40, stok berkurang 80. Beli 1 renceng isi 10,
berkurang 10. Ketiga varian membaca stok yang sama.

Repository dites dengan Postgres asli lewat testcontainers, diberi tag build atau
dilewati dengan `testing.Short()` supaya `make test` tetap cepat.

---

## 16. Troubleshooting

**Stok habis di layar padahal barangnya ada**
Worker pelepas reservasi tidak jalan. Cek `inventory_reservations` yang berstatus
`held` dengan `expires_at` sudah lewat.

**Pesanan tidak pernah kedaluwarsa / dana tidak pernah cair**
Worker tidak jalan. Ini penyebab paling umum saat development.

**Pesanan sudah dibayar tapi statusnya masih menunggu**
Webhook hilang. Jalankan rekonsiliasi. Kalau sering terjadi, periksa apakah
endpoint webhook bisa dijangkau dari internet dan apakah signature-nya lolos.

**Nil pointer saat endpoint baru dipanggil**
Repository belum dirakit di `module.go` dan diteruskan ke `NewService`.

**Event handler tidak pernah jalan**
`RegisterSubscriptions` lupa dipanggil di `registry.go`. Cek juga `outbox_events`:
kalau baris menumpuk dengan `published_at IS NULL`, relay tidak jalan.

**Deadlock saat checkout ramai**
Baris stok dikunci dalam urutan berbeda antar transaksi. Pastikan
`ORDER BY outlet_id, variant_id` sebelum `FOR UPDATE`.

**Total pesanan tidak sama dengan jumlah suborder**
Sisa pembulatan dibuang di suatu tempat. Cari pembagian yang tidak memakai
`money.Distribute`.

**Penjual protes bagi hasilnya tidak sesuai**
Periksa `owner_type` voucher yang dipakai. Voucher marketplace dipotong dari
komisi kita; voucher penjual dipotong dari pendapatan penjual. Kalau tertukar,
pembukuan salah.

**Rekonsiliasi selisihnya tidak nol**
Ada jurnal yang tidak seimbang atau ada pergerakan uang yang tidak dicatat.
Selidiki hari itu juga — selisih kecil yang dibiarkan akan jadi besar.

**Ulasan muncul dari orang yang tidak membeli**
`order.Port.HasPurchased` tidak dipanggil, atau `order_item_id` dipercaya dari
klien tanpa verifikasi.

**`import cycle not allowed`**
Dua modul saling mengimpor. Salah satu arah harus jadi event. Kalau `order` yang
diimpor balik, hampir pasti itu yang salah.

**`use of internal package not allowed`**
Mencoba mengimpor isi modul lain. Bekerja sesuai rancangan — tambahkan method di
`port.go` modul tujuan.
