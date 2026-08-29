# Fase 12 — Siklus Suborder

> Prasyarat: fase 10 + fase 06. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 12), §5.
> ADR: 002 (suborder unit kerja, refund parsial), 009 (state machine, delivered ≠
> completed).

## Tujuan

Siklus suborder: konfirmasi → packing → kirim → terima, plus worker kedaluwarsa/auto.
**Selesai kalau** stok ter-commit saat kirim dan satu penjual menolak tidak membatalkan
pesanan penjual lain.

## Aturan khusus fase ini

- Semua perubahan status lewat `suborder.Transition(to, now)`; setiap transisi menulis
  `order_status_events`.
- `RejectSuborder` (mudah salah): lepas **hanya stok suborder itu** (`inventory.Release`
  parsial / per-suborder), **refund parsial** nilai suborder itu, **order induk TIDAK
  dibatalkan**, lalu `order.SyncStatusFromSuborders()`.
- `inventory.Commit` dipanggil saat **`MarkShipped`** (bukan saat `OrderPaid`).
- `delivered` ≠ `completed`: `delivered` = barang sampai; `completed` = tidak ada
  sengketa lagi, dana boleh cair. `Shipped` **tidak bisa** di-cancel — hanya retur (alur
  terpisah, di luar scope).
- `CancelOrder` oleh pembeli hanya bila **belum ada suborder yang shipped**.
- Worker `expiry_usecase` semuanya **idempoten** + `FOR UPDATE SKIP LOCKED`: dijalankan
  dua kali tidak melepas stok / mengembalikan kuota voucher dua kali.
- Semua panggilan `inventory.*` / `promotion.*` di dalam `WithinTx`; panggilan
  `payment` / `fulfillment` yang menyentuh jaringan di luar tx atau lewat worker.
- Event suborder yang diterbitkan **wajib** memuat `SellerID` + `OutletID`.

## Urutan kerja

### App — `order/internal/app/`

- [ ] `suborder_usecase.go` (seller dashboard, tiap usecase diawali
      `requireSuborderSeller`):
  - `ListSellerOrders`, `GetSellerOrder`.
  - `ConfirmSuborder` (awaiting → confirmed).
  - `RejectSuborder(reason)` — `Transition(Rejected)` + `inventory.Release` suborder itu
    + `payment.RequestRefund` parsial (SuborderID) + `SyncStatusFromSuborders` +
    outbox `EventSuborderRejected`. **Order induk tidak dibatalkan.**
  - `StartPacking` (confirmed → packing).
  - `MarkShipped(ShipInput)` — `fulfillment.CreateShipment` (di luar / sebelum tx),
    lalu `WithinTx`: `Transition(Shipped)` + `inventory.Commit` (tulis Movement
    `kind=sale`) + outbox `EventSuborderShipped`.
  - `MarkReadyForPickup` — terbitkan kode pickup via `fulfillment`.
- [ ] `order_usecase.go` (buyer, diawali `requireOrderOwner`):
  - `ListMyOrders`, `GetMyOrder` (tampilkan semua suborder dengan status berbeda — "1
    dari 3 toko sudah dikirim", jujur).
  - `CancelOrder(reason)` — hanya bila `CanBeCancelledByBuyer()`; `WithinTx`:
    `order.Cancel` + `inventory.Release` semua + `promotion.ReleaseRedemption` +
    `payment` batalkan + outbox `EventOrderCancelled`.
  - `ConfirmReceived(suborderID)` — buyer memindah suborder ke `Completed` lebih awal.
- [ ] `expiry_usecase.go` (worker, semua idempoten, `FOR UPDATE SKIP LOCKED`):
  - `ExpireUnpaidOrders` — `ListExpired`; per order `WithinTx`: `order.Cancel("payment
    expired")` + `inventory.Release` + `promotion.ReleaseRedemption` +
    `payment.MarkExpired` + outbox `EventOrderExpired`.
  - `AutoCompleteDelivered` — `ListDeliveredBefore(window)` → `Transition(Completed)` +
    outbox `EventSuborderCompleted` (memicu pematangan dana seller).
  - `AutoRejectUnconfirmed` — suborder awaiting melewati deadline → `Transition(Rejected)`
    (perlakuan sama seperti `RejectSuborder`: lepas stok + refund parsial).
- [ ] `handlers.go` — bila ada event internal yang order sendiri konsumsi (mis.
      `EventPaymentSettled` → `order.MarkPaid` + `SyncStatusFromSuborders` + outbox
      `EventOrderPaid`). Idempoten.

### Migrasi
- [ ] Pastikan `order_status_events` (dari migrasi `000008`, fase 10) menyimpan tiap
      transisi: `from`, `to`, `actor`, `reason`, `at`. Bila perlu kolom tambahan → migrasi
      baru `make migrate-create name=...`, jangan edit `000008`.

### Transport
- [ ] `order/internal/transport/rest/seller_handler.go` — endpoint confirm/reject/pack/
      ship/ready di `/api/v1/seller/suborders/{id}/...`.
- [ ] `order/internal/transport/rest/order_handler.go` — `cancel` + `confirm-received` di
      `/api/v1`.
- [ ] Perbarui `routes.go`, `request.go` (`Validate()`), `response.go`.

### Wiring
- [ ] `order/module.go` `RegisterWorkers` — daftarkan `ExpireUnpaidOrders`,
      `AutoCompleteDelivered`, `AutoRejectUnconfirmed`.
- [ ] `internal/app/registry.go` — pastikan `RegisterWorkers` order dipanggil; worker
      dijalankan dari `cmd/worker` di fase 13.
- [ ] `order/module.go` `RegisterSubscriptions` — subscribe `EventPaymentSettled` bila
      `order` yang menangani transisi ke `paid`.

## Test wajib

- **Refund parsial** (`../GUIDES.md` §15): 1 penjual menolak, 2 lain jalan terus, jumlah
  refund tepat, order induk **tidak** batal, status induk jadi `partially_fulfilled`.
- `MarkShipped` men-`Commit` stok (Movement `kind=sale` tertulis) dan menutup jalur
  cancel.
- `CancelOrder` setelah salah satu suborder shipped → `ErrOrderNotCancellable`.
- **Idempotensi worker**: `ExpireUnpaidOrders` jalan dua kali → stok dilepas sekali,
  kuota voucher dikembalikan sekali.
- Transisi ilegal (`Shipped → Confirmed`) → `ErrInvalidTransition`.
- `SyncStatusFromSuborders`: semua suborder completed → order `completed`; semua
  cancelled/rejected → order `cancelled`.

## Sengaja TIDAK dikerjakan di fase ini

- Alur retur barang setelah `Shipped` (belum ada di state machine — butuh desain + ADR).
- Jurnal buku besar untuk refund parsial (fase 14).
- Menjalankan worker (fase 13 merakit `cmd/worker`).
