# Fase 06 — Inventory

> Prasyarat: fase 04. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 6), §11.
> ADR terkait: ADR-003 (tahan saat checkout), ADR-006 (satuan dasar), ADR-007.

## Tujuan

Stok per outlet, reservasi, mutasi. **Selesai kalau** stok per outlet benar dan konversi
satuan benar (beli 2 dus isi 40 → stok berkurang 80).

## Aturan khusus fase ini

- Semua kuantitas di `inventory.Port` dalam **satuan dasar**. **Pemanggil** yang
  mengonversi lewat `Variant.ToBaseQuantity()`. Jangan taruh konversi di sini.
- Siklus: ditahan saat **checkout** (ADR-003), di-commit saat **kirim**, dilepas saat
  **batal**. `Reserve`/`Commit`/`Release` jalan **di dalam transaksi milik pemanggil**
  (`ExecutorFrom`), tidak buka tx sendiri.
- `Stock.Commit` → `on_hand -= qty` **dan** `reserved -= qty`. `Stock.ReleaseReservation`
  → `reserved -= qty` saja, **tidak menulis Movement** (barang tak pernah keluar rak).
- **Setiap perubahan `on_hand` menulis `Movement`.** Tanpa kecuali.
- `GetForUpdate` mengunci baris dengan **urutan konsisten** `ORDER BY outlet_id,
  variant_id` sebelum `FOR UPDATE`. Urutan beda antar tx = deadlock.
- `AvailableMany` boleh cache Redis pendek untuk halaman produk — **tidak pernah** untuk
  checkout.
- `on_hand` dan `reserved` disimpan terpisah supaya penjual tahu "ada di rak tapi sudah
  dipesan".

## Urutan kerja

### Domain
- [ ] `internal/modules/inventory/internal/domain/stock.go` — `Stock{QuantityOnHand,
      QuantityReserved,LowStockThreshold,Version int}`; `Available()=OnHand-Reserved`;
      `CanReserve`, `Reserve`, `Commit` (kurangi keduanya), `ReleaseReservation`
      (reserved saja), `Adjust(delta)`, `IsLow`.
- [ ] `.../domain/reservation.go` — `ReservationStatus` = held|committed|released|
      expired; `Reservation{OrderID,SuborderID,OutletID,VariantID,Quantity,Status,
      ExpiresAt}`; `IsExpired`, `Commit` (hanya dari held), `Release` (hanya dari held).
- [ ] `.../domain/movement.go` — append-only; `MovementKind` = Purchase|Sale|Return|
      Adjustment|TransferIn|TransferOut|Damage|Expiry; `Movement{Kind,Quantity(+/-),
      BalanceAfter,ReferenceType,ReferenceID,Note,ActorID}`.
- [ ] `.../domain/opname.go` — `Opname{OutletID,Status,StartedBy}`,
      `OpnameItem{SystemQuantity,CountedQuantity,Difference}`; `Finish(items) ->
      []Movement` (hasil masuk sebagai `kind=adjustment`, bukan edit `on_hand` langsung).
- [ ] `.../domain/repository.go` — `StockRepository{Get,GetForUpdate(pairs)
      SELECT..FOR UPDATE urut (outlet_id,variant_id),GetMany,Upsert,OutletsWithStock,
      ListByOutlet}`, `ReservationRepository{CreateMany,ListByOrder,ListExpired(before,
      limit),UpdateStatus}`, `MovementRepository{Create,CreateMany,ListByVariant,
      SumByVariant}`, `OpnameRepository`.
- [ ] `.../domain/errors.go` — `ErrInsufficientStock` (bawa variant + requested +
      available), `ErrReservationExpired`, `ErrNegativeStock`, `ErrStockVersionConflict`.

### Migrasi
- [ ] `migrations/000005_create_inventory.up.sql` / `.down.sql` — `inventory_stocks`
      (`quantity_on_hand`/`quantity_reserved` `numeric(14,3)`, `version`),
      `inventory_reservations` (`expires_at`, status), `inventory_movements`
      (append-only), `inventory_opnames` + items, sesuai spec; index untuk `ListExpired`.

### Infra
- [ ] `.../infra/{stock_repository,reservation_repository,movement_repository,
      opname_repository,mapper}.go` — `GetForUpdate` lock berurutan.

### App
- [ ] `.../app/reservation_usecase.go` (**paling kritis**) — `Reserve(orderID,suborderID,
      items)`: (1) urutkan item by (outlet_id,variant_id), (2) `GetForUpdate` semua baris
      sekaligus, (3) cek `CanReserve` tiap → `ErrInsufficientStock` berdetail, (4)
      `Reserve` + simpan, (5) buat `Reservation` `ExpiresAt = now + TTL`. Jalan **di tx
      pemanggil**. `Commit(orderID)`: held→committed, `on_hand`&`reserved` turun, tulis
      `Movement kind=sale`. `Release(orderID)`: held→released, `reserved` turun, **tanpa
      Movement**. `ReleaseExpired()` worker: `ListExpired` + `FOR UPDATE SKIP LOCKED`,
      lepas per batch.
- [ ] `.../app/availability_usecase.go` — `Available`, `AvailableMany` (query terpanas,
      `WHERE (outlet_id,variant_id) IN (...)`, tanpa loop), `OutletsWithStock`. Cache
      Redis pendek OK untuk halaman produk, **tidak untuk checkout**.
- [ ] `.../app/stock_usecase.go` — seller: `ListStock`, `AdjustStock(delta,note)` (tulis
      Movement + note + actor), `SetLowStockThreshold`, `TransferStock` (dua Movement
      transfer_out/transfer_in dalam satu tx, jumlah sama).
- [ ] `.../app/opname_usecase.go` — `StartOpname`, `SubmitCount`, `FinishOpname`
      (Movement adjustment per selisih, satu tx, resumable), `ListOpnames`.
- [ ] `.../app/{service,dto,port_adapter,handlers}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,stock_handler,opname_handler,
      handler}.go` — di `/api/v1/seller`.

### Port
- [ ] `internal/modules/inventory/port.go` — `Available`, `AvailableMany` (batch),
      `OutletsWithStock(variantID,minQty)`, `Reserve(orderID,suborderID,items)`,
      `Commit(orderID)`, `Release(orderID)`. Semua kuantitas **satuan dasar**.
- [ ] `events.go` — `EventStockReserved/Committed/Released/Adjusted/LowStock/OutOfStock`.
- [ ] `module.go` — `Config{...,Catalog catalog.Port,ReservationTTL time.Duration}`;
      subscription + worker pelepas reservasi (didaftarkan di fase 13, tapi usecase
      `ReleaseExpired` ditulis di sini).

### Wiring
- [ ] `internal/app/registry.go` — build `inventory` setelah `catalog` (terima
      `catalogMod.Port()`).

## Test wajib

- **Konversi satuan**: beli 2 dus isi 40 → stok turun 80; 1 renceng isi 10 → turun 10;
  tiga varian barang sama membaca stok yang sama.
- **Checkout bersamaan stok terakhir**: dua goroutine `Reserve` variant stok = 1 → tepat
  satu berhasil, satu `ErrInsufficientStock`.
- `Release` tidak menulis Movement; `Commit` menulis Movement `kind=sale` dengan
  `BalanceAfter` benar.
- Idempotensi: `ReleaseExpired` dijalankan dua kali tidak melepas reservasi dua kali.
- `GetForUpdate` mengunci dalam urutan `(outlet_id,variant_id)`.

## Sengaja TIDAK dikerjakan di fase ini

- Pemanggilan `Reserve`/`Commit`/`Release` dari checkout & siklus suborder (fase 10 /
  12) — hanya pastikan port + usecase siap dan teruji.
- Pendaftaran worker di `cmd/worker` (fase 13).
- Konsumen `EventLowStock` di notification (fase 13).
