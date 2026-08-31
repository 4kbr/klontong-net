# TASKS — Urutan Implementasi Backend

Backend ini **scaffold**: semua `.go` hanya `package` + komentar TODO kontrak, semua
`migrations/*.up.sql` hanya spec prosa tanpa DDL. Dokumen ini menerjemahkan urutan 15
tahap resmi di [`GUIDES.md` §2](GUIDES.md) jadi checklist per-file yang bisa langsung
dikerjakan.

## Cara pakai

1. Kerjakan **berurutan**. Tiap file fase = satu atau beberapa PR/commit fokus.
2. **Jangan lompat.** Khususnya jangan mulai fase 10 (`order`) sebelum fase 6
   (`inventory`) dan 9 (`promotion`) selesai — checkout memanggil keduanya.
3. Tiap fase punya file di [`tasks/`](tasks/) dengan template sama: Tujuan → Aturan
   khusus → Urutan kerja (domain → migrasi → repo → usecase → transport → wiring) →
   Test wajib → Sengaja tidak dikerjakan.
4. Selesai satu fase, jalankan **Definition of Done** di bawah sebelum lanjut.

## 15 Fase

| # | Fase | Modul / paket | Selesai kalau | Butuh fase |
|---|---|---|---|---|
| [01](tasks/01-platform.md) | Platform | `internal/platform/*`, `cmd/api`, `cmd/migrate`, migrasi `000015` (outbox) | `/healthz` 200, property test `money.Distribute` lulus | — |
| [02](tasks/02-identity-customer.md) | Identity + Customer | `modules/identity`, `modules/customer` | bisa daftar, login, tambah alamat | 01 |
| [03](tasks/03-seller-outlet.md) | Seller + Outlet | `modules/seller` | bisa buka toko + tambah outlet berkoordinat | 02 |
| [04](tasks/04-catalog.md) | Catalog | `modules/catalog` | bisa buat produk dengan varian dus/renceng/pcs | 03 |
| [05](tasks/05-pricing.md) | Pricing | `modules/pricing` | harga per satuan + tier grosir tampil benar | 04 |
| [06](tasks/06-inventory.md) | Inventory | `modules/inventory` | stok per outlet benar, konversi satuan benar | 04 |
| [07](tasks/07-cart.md) | Cart + pemilihan outlet | `modules/cart` | keranjang 3 penjual terkelompok benar | 05, 06 |
| [08](tasks/08-fulfillment.md) | Fulfillment | `modules/fulfillment` | ongkir per suborder tampil, pickup gratis | 03 |
| [09](tasks/09-promotion.md) | Promotion | `modules/promotion` | voucher terbagi ke suborder, jumlahnya persis | 03 |
| [10](tasks/10-order-checkout.md) | Order: preview + place | `modules/order` | pesanan multi-penjual terbentuk dengan suborder | 06, 07, 08, 09, 11(minimal) |
| [11](tasks/11-payment.md) | Payment + webhook + rekonsiliasi | `modules/payment` | pesanan jadi `paid` lewat webhook palsu | 01 |
| [12](tasks/12-suborder-lifecycle.md) | Siklus suborder | `modules/order` (usecase) | stok ter-commit saat kirim | 10, 06 |
| [13](tasks/13-outbox-audit-notification.md) | Outbox + Audit + Notification | `cmd/worker`, `modules/audit`, `modules/notification` | setiap aksi tercatat, notifikasi terkirim | 01, 10 |
| [14](tasks/14-settlement.md) | Settlement | `modules/settlement` | rekonsiliasi selisihnya nol | 10, 11, 13 |
| [15](tasks/15-review.md) | Review | `modules/review` | ulasan hanya dari pembelian yang selesai | 12 |

Fase 1, 9, 14, dan checkout di fase 10 **menyentuh uang → tulis test lebih dulu**
([`GUIDES.md` §15](GUIDES.md)).

## Fase integrasi lintas-app

Setelah fase 12, jalur buyer sudah utuh dari sisi backend. Saat frontend juga sudah
sampai fase 05, jalankan [`docs/tasks/INTEGRASI.md`](../../../docs/tasks/INTEGRASI.md)
— lepas mock, sambungkan storefront ke API asli, dan cocokkan respons dengan kontrak.

Fase itu **tidak memblokir**. Kalau frontend belum siap, backend lanjut ke fase 13.

## Aturan global (berlaku di semua fase)

Ringkas dari [`../AGENTS.md`](../AGENTS.md) dan [`../../../AGENTS.md`](../../../AGENTS.md).
Kalau ragu, baca ulang di sana.

- **Uang**: `money.Amount` (int64 rupiah). **Tidak ada `float`, tidak pernah**, di jalur
  uang. Persentase = basis poin `int` (250 = 2,5%). Pembagian pakai `money.Distribute`.
  Sisa pembulatan tidak dibuang. Komisi dibulatkan ke bawah. JSON kirim angka.
- **Transaksi** dibuka **hanya di layer `app`** lewat `s.tx.WithinTx`. **Tidak ada
  panggilan jaringan di dalam transaksi database** — gateway/kurir/transfer dipanggil di
  luar tx atau diserahkan ke worker.
- **Repository** dibuka `postgres.ExecutorFrom(ctx, r.pool)`, ditutup
  `postgres.Translate(err)`. Repo tidak pernah membuka tx dan tidak memanggil repo lain.
- **Baris diperebutkan** (stok, voucher, payout): `FOR UPDATE` dengan **urutan
  konsisten** (`ORDER BY outlet_id, variant_id` untuk stok). Urutan beda = deadlock.
- **Status** hanya berubah lewat `Transition()` / state machine. Status Order induk
  **dihitung** dari suborder (`SyncStatusFromSuborders`), tidak pernah disetel langsung.
- **Event** ditulis lewat `outbox.Save(ctx, evt)` **di dalam transaksi yang sama** dengan
  perubahan datanya. Tidak pernah `bus.Publish` dari usecase. Event suborder **wajib**
  memuat `SellerID` dan `OutletID`.
- **Idempotensi**: penjaga yang mengikat selalu **unique index di database**, bukan
  Redis. Kunci dibuat sekali saat objek dibuat. `POST /checkout` menolak request tanpa
  `Idempotency-Key`. Webhook duplikat dibalas 200.
- **Otorisasi**: setiap usecase publik diawali fungsi `require*` dari `authz.go`.
  Kepemilikan diperiksa di usecase, bukan cukup peran.
- **Buku besar**: tidak ada kolom saldo yang di-UPDATE. Setiap jurnal wajib seimbang;
  gagalkan transaksi kalau tidak.
- **Batas modul**: isi `internal/` modul lain tidak boleh diimpor. Butuh data modul
  lain? Tambah method di `port.go` modul tujuan. Tidak ada modul bergantung balik ke
  `order` kecuali lewat event. Method port yang dipanggil dalam loop **wajib punya versi
  batch**.
- **Kontrak API**: [`contracts/openapi`](../../../contracts/openapi) **mengikat**
  (ADR-015). Path, verb, nama field, dan enum diambil dari sana — tiap file fase
  menunjuk `paths/*.yaml`-nya. Yang dikunci kontrak dan sering tidak tertulis di sisi
  Go: field JSON `snake_case`; waktu RFC3339 UTC; uang `integer` int64 (persen basis
  poin `integer`); kuantitas `number`; envelope sukses `{data,meta}`, error
  `{error:{code,message,fields?,retryable?}}` dengan `code` `lower_snake_case`;
  paginasi keyset `?limit` + `?cursor` dengan `meta.next_cursor` + `meta.has_more`;
  `Authorization: Bearer <jwt>` (akses TTL 15 menit); `X-Request-ID` di echo semua
  respons; `Retry-After` pada 429. Perlu menyimpang? Ubah kontrak di **PR yang sama**,
  `make -C contracts ci` hijau — jangan diam-diam.
- **Konvensi**: `context.Context` parameter pertama. `clock.Now()` bukan `time.Now()`.
  Error dibungkus `errs.*`. SQL sebagai konstanta string. Struct response HTTP terpisah
  dari entity domain. Komentar Bahasa Indonesia. Package tidak pernah `utils`/`helpers`/
  `common`/`shared`.
- **Scaffold**: kerjakan **hanya yang diminta fase**. TODO tetangga yang tidak diminta
  → sebutkan, jangan dikerjakan.

## Definition of Done per fase

- [ ] `go build ./...` hijau
- [ ] `make lint` bersih (tanpa `//nolint` baru)
- [ ] `make test` hijau
- [ ] Migrasi fase: `make migrate-up` lalu `make migrate-down` jalan tanpa error
- [ ] Test wajib fase (lihat bagian "Test wajib" di file fase) ada dan lulus
- [ ] Modul baru dirakit di `internal/app/registry.go`; `RegisterSubscriptions` /
      `RegisterWorkers` (bila ada) dipanggil dari `registry.go`
- [ ] Path, verb, nama field, dan enum fase ini cocok dengan
      [`contracts/openapi`](../../../contracts/openapi). Beda → ubah kontrak di PR yang
      sama, `make -C contracts ci` hijau (ADR-015)
- [ ] Laporan singkat: file apa yang diisi, keputusan yang diambil, TODO yang sengaja
      ditinggalkan

## Catatan scaffold yang perlu dibereskan awal

- Ada **dua** `Makefile`. Root punya `setup`, `up`, `down`, `logs`, `dev`, `worker`,
  `tunnel`, `migrate`, `migrate-create`, `seed`, `test`, `lint`, `check`.
  `apps/backend/Makefile` punya `migrate-up`/`migrate-down`/`migrate-force`/
  `migrate-version`, `run`, `build`, `fmt`, `cover`, `tools`. Perintah di `GUIDES.md`
  §1 dan §3 (`make setup`, `make up`, `make migrate`, `make tunnel`) dijalankan **dari
  root** — semuanya sudah ada, tidak perlu ditambah. `Makefile` root juga dipakai
  frontend dan `contracts/`; jangan menambah target backend di sana.
- Semua `migrations/*.up.sql` hanya komentar spec. Tiap fase menuliskan DDL-nya sesuai
  spec yang sudah ada di file itu. Kolom uang `bigint`, kolom kuantitas `numeric(14,3)`.
- Semua `migrations/*.down.sql` hanya `-- TODO` — tulis drop yang benar-benar jalan,
  urutan terbalik.
