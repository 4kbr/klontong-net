# Fase 10 — Order: Preview + Place

> Prasyarat: fase 06, 07, 08, 09, dan fase 11 minimal (`payment.Port.CreatePayment` +
> `gateway_noop`). Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 10), §4, §5, §8,
> plus `../ARCHITECTURE.md` §6. ADR: 002, 004, 005, 007, 008, 009, 010.
> **File paling kompleks di proyek — tulis test lebih dulu** (`../GUIDES.md` §15).

## Tujuan

Checkout dua tahap + entity Order/Suborder/Item. **Selesai kalau** pesanan multi-penjual
terbentuk dengan suborder per penjual.

## Aturan khusus fase ini

- **Suborder adalah unit kerja** (ADR-002). Status Order induk **dihitung** dari suborder
  (`SyncStatusFromSuborders`), tidak pernah disetel langsung setelah pembayaran.
- Transisi status **hanya lewat `Suborder.Transition(to, now)`** — peta di `state.go`
  sebagai data, bukan `if` tersebar (ADR-009).
- Struct request checkout **tidak punya field harga/total** (ADR-004). Klien kirim
  **pilihan** (alamat, kurir, voucher, barang+qty). `PlaceOrder` **hitung ulang dari
  nol**; beda dengan angka yang ditampilkan klien → `ErrPriceChanged` + rincian.
- `POST /api/v1/checkout` **menolak** tanpa `Idempotency-Key` (ADR-008); di balik
  middleware `Idempotent`.
- Di `WithinTx`: `inventory.Reserve` → `promotion.Redeem` → buat Order+Suborder+Item
  (semua **snapshot**, ADR-010) → **bekukan `CommissionBPS`** tiap suborder dari
  `seller.Port` → `cart.MarkConverted` → `outbox.Save(EventOrderPlaced)`.
- **`payment.CreatePayment` dipanggil SETELAH commit** (atau sebagai pending sebelum tx,
  ditautkan setelah). **Tidak ada panggilan jaringan di dalam tx.**
- `inventory.Reserve` mengunci baris stok dalam urutan konsisten — inventory yang
  menangani; checkout hanya memanggil.
- `Order.RecalculateTotals`: `GrandTotal = Items − Discount + Shipping + Tax` **dan**
  harus = Σ semua suborder. Beda 1 rupiah → `ErrTotalMismatch` = **bug internal**, jangan
  buat order.
- `NextOrderNumber` pakai **sequence Postgres**, bukan `SELECT MAX(...)+1`.
- Preview **tidak menahan apa pun** dan boleh dipanggil berkali-kali.

## Urutan kerja

### Domain
- [ ] `internal/modules/order/internal/domain/state.go` — `suborderTransitions` map:
      `AwaitingConfirmation → {Confirmed,Rejected,Cancelled}`; `Confirmed →
      {Packing,Cancelled}`; `Packing → {Shipped,ReadyForPickup,Cancelled}`;
      `ReadyForPickup → {Delivered,Cancelled}`; `Shipped → {Delivered}` (**tidak bisa
      cancel setelah shipped**); `Delivered → {Completed}`; `Completed`/`Cancelled`/
      `Rejected` terminal. `CanTransition(from,to)`, `TerminalStatuses()`.
- [ ] `.../domain/suborder.go` — `SuborderStatus` consts; `Suborder{OrderID,Number,
      SellerID,OutletID,Status,FulfillmentMethod,Items,ItemsAmount/DiscountAmount/
      ShippingAmount/TaxAmount/TotalAmount,CommissionBPS(FROZEN),CommissionAmount,
      SellerEarningAmount,timestamps,CancelReason,RejectedReason}`; `Transition(to,now)`
      (satu-satunya jalan); `ComputeCommission()` (`CommissionAmount = discounted
      ItemsAmount × CommissionBPS`, bulat **ke bawah**; `SellerEarningAmount =
      (ItemsAmount − DiscountAmount) − CommissionAmount + ShippingAmount bila seller
      menanggung ongkir`); `RequiresShipment()` (pickup tidak).
  > **Keputusan terbuka**: komisi dari barang saja atau termasuk ongkir? Rekomendasi
  > `../DECISIONS.md`: barang saja. **Tulis ADR baru sebelum memilih**, jangan diam-diam.
- [ ] `.../domain/order.go` — `Status` = pending_payment|paid|partially_fulfilled|
      completed|cancelled|expired; `Order{Number,BuyerUserID,ItemsAmount,DiscountAmount,
      ShippingAmount,TaxAmount,GrandTotal,PaymentMethod,ShippingAddress AddressSnapshot,
      Suborders,timestamps,ExpiresAt,CancelReason}`; `RecalculateTotals()` (+cek = Σ
      suborder), `MarkPaid`, `Cancel(reason,now)`, `SyncStatusFromSuborders()`,
      `CanBeCancelledByBuyer()` (true selama belum ada suborder shipped).
- [ ] `.../domain/item.go` — `Item{OrderID,SuborderID,VariantID,ProductID,SellerID,
      ProductName/VariantName/SKU/ImageURL(COPY),UnitCode/ContentQuantity(COPY),Quantity,
      UnitPriceAmount(COPY),DiscountAmount,TotalAmount,WeightGram}`; `BaseQuantity() =
      Quantity × ContentQuantity`, `TotalWeightGram()`.
- [ ] `.../domain/address_snapshot.go` — `AddressSnapshot{recipient/geo,Latitude/
      Longitude *float64}`, `FromCustomerAddress(a customer.Address)`; disimpan jsonb.
- [ ] `.../domain/repository.go` — `OrderRepository{Create(order+suborder+item),Update,
      FindByID,FindByNumber,ListByBuyer(filter,cursor,limit),ListExpired(before,limit),
      NextOrderNumber(t)}`, `SuborderRepository{Update,FindByID,ListByOrder,ListBySeller
      (filter,cursor,limit),ListDeliveredBefore(before,limit)}`, `ItemRepository
      {ListByOrder,ListBySuborder,FindPurchased}`, `StatusEventRepository{Create,
      ListByOrder}`.
- [ ] `.../domain/errors.go` — `ErrNotOrderOwner`, `ErrInvalidTransition`,
      `ErrOrderNotCancellable`, `ErrCartEmpty`, `ErrPriceChanged` (detail per baris),
      `ErrOutOfStock` (detail per baris), `ErrAddressRequired`, `ErrNoFulfillmentAvailable`,
      `ErrTotalMismatch`, `ErrOrderExpired`.

### Migrasi
- [ ] `migrations/000008_create_order.up.sql` / `.down.sql` — `order_orders`
      (`order_number`, status enum, kolom amount `bigint`, `shipping_address` jsonb),
      `order_suborders` (`commission_bps`, `commission_amount`, status enum),
      `order_items` (kolom snapshot), `order_status_events` sesuai spec. **Sequence** untuk
      nomor order.

### Infra
- [ ] `.../infra/{order_repository,suborder_repository,item_repository,
      status_event_repository,mapper}.go` — `Create` menulis order+suborder+item dalam
      satu tx pemanggil; `NextOrderNumber` pakai sequence.

### App
- [ ] `.../app/authz.go` — `requireOrderOwner(orderID)`, `requireSuborderSeller
      (suborderID)`, `requireAdmin`.
- [ ] `.../app/checkout_usecase.go` — **tulis test dulu**.
  - `PreviewCheckout(PreviewInput) -> CheckoutPreview`: (1) `cart.Port` ambil keranjang,
    (2) `customer.Port` ambil alamat, (3) kelompok per seller → calon suborder, (4) per
    kelompok: pilih outlet, `pricing.ResolveMany`, `inventory.AvailableMany`, (5) opsi
    kirim per kelompok via `fulfillment.Port`, (6) voucher via `promotion.Port`, split
    diskon per suborder (`money.Distribute`), (7) total + rincian per suborder. **Tidak
    menahan apa pun.**
  - `PlaceOrder(PlaceOrderInput) -> *domain.Order`: butuh `Idempotency-Key`. (1) hitung
    ulang semuanya, percaya **hanya pilihan klien**; (2) bandingkan dengan angka klien →
    `ErrPriceChanged` berdetail bila beda; (3) `WithinTx`: `inventory.Reserve` semua →
    `promotion.Redeem` → buat Order+Suborder+Item snapshot → freeze `CommissionBPS` →
    `cart.MarkConverted` → outbox `EventOrderPlaced`; COMMIT. (4) `payment.CreatePayment`
    / tandai COD **setelah commit**.
- [ ] `.../app/{service,dto,port_adapter,handlers}.go`.
- [ ] `.../app/order_usecase.go`, `suborder_usecase.go`, `expiry_usecase.go` — **kerangka
      saja di fase ini**; isi penuh di fase 12.

### Transport
- [ ] `.../transport/rest/{routes,request,response,checkout_handler,order_handler,
      seller_handler,admin_handler,handler}.go` — `POST /api/v1/checkout/preview` dan
      `POST /api/v1/checkout` (di balik `Idempotent`, `Authenticate`). Request checkout
      **tanpa field harga/total**.

### Port
- [ ] `internal/modules/order/port.go` — `GetOrder`, `GetSuborder`,
      `HasPurchased(userID,variantID) -> (orderItemID, ok)`.
- [ ] `events.go` — order-level: `order.order.placed/paid/cancelled/expired/completed`;
      suborder-level: `order.suborder.confirmed/rejected/shipped/delivered/completed/
      cancelled`. **Payload event suborder wajib `SellerID` + `OutletID`.**
- [ ] `module.go` — `Config` dengan port: Cart, Catalog, Pricing, Inventory, Promotion,
      Fulfillment, Payment, Customers, Sellers; `RegisterSubscriptions` +
      `RegisterWorkers` (expiry, diisi fase 12).

### Wiring
- [ ] `internal/app/registry.go` — build `order` setelah semua dependensinya; **tidak ada
      modul yang bergantung balik ke `order`** kecuali lewat event.

## Test wajib

- **Checkout menghitung ulang**: klien kirim total lebih murah → `ErrPriceChanged`, order
  tidak dibuat.
- **Idempotensi**: `POST /checkout` dengan `Idempotency-Key` sama dua kali → satu order.
- **Tanpa `Idempotency-Key`** → ditolak.
- **Multi-penjual**: keranjang 3 seller → 1 Order + 3 Suborder, masing-masing punya
  ongkir & `CommissionBPS` beku sendiri.
- `Order.RecalculateTotals`: Σ suborder ≠ grand total → `ErrTotalMismatch`, order tidak
  dibuat.
- **No network call dalam tx**: assert `payment.CreatePayment` dipanggil setelah commit
  (port palsu merekam urutan).
- Reserve gagal (stok kurang) → seluruh checkout batal, tidak ada order/redemption.

## Sengaja TIDAK dikerjakan di fase ini

- Siklus suborder confirm/ship/receive + worker expiry (fase 12) — kerangka file saja.
- Implementasi `payment` penuh (fase 11) — cukup `CreatePayment` + `gateway_noop`.
- Handler event di audit/notification/settlement (fase 13 / 14).
