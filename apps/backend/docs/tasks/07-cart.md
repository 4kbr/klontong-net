# Fase 07 — Cart + Pemilihan Outlet

> Prasyarat: fase 05 + fase 06. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 7).
> ADR terkait: ADR-003 (tidak menahan stok dari keranjang), ADR-004 (tanpa harga di
> request), ADR-007 (satu suborder satu outlet).
> Kontrak (**mengikat**, ADR-015): [`paths/cart.yaml`](../../../../contracts/openapi/paths/cart.yaml).

## Tujuan

Keranjang dikelompokkan per penjual + pemilihan outlet. **Selesai kalau** keranjang 3
penjual terkelompok dengan benar.

## Aturan khusus fase ini

- Add-to-cart **TIDAK menahan stok** (ADR-003). Boleh memperingatkan stok menipis,
  **tidak pernah menolak** karena stok.
- `Item.UnitPriceAmount` hanya **cache tampilan**. Harga dihitung ulang setiap keranjang
  dibuka, dan lagi saat checkout (ADR-004).
- `GetCart` melakukan **4 batch call** (catalog / pricing / inventory / seller), bukan
  4N query.
- `GroupBySeller()` dipakai hampir di semua tampilan keranjang **dan** jadi dasar
  pemecahan suborder.
- `Merge` saat guest login: gabung, **jangan timpa**; varian+outlet sama → jumlahkan qty
  dengan batas per-baris.
- Satuan tak-pecahan menolak qty pecahan; qty ≤ 0 = hapus baris, bukan error.
- `SuggestOutlet`: stok cukup + outlet aktif + dukung ≥1 metode + **terdekat** dengan
  alamat pembeli; tidak ada → `ErrInsufficientStock`. Satu suborder = satu outlet
  (ADR-007).
- `ReassignOutlets` dipanggil ulang setiap alamat kirim berubah.

## Urutan kerja

### Domain
- [ ] `internal/modules/cart/internal/domain/cart.go` — `Status` = active|converted|
      abandoned; `Cart{UserID *uuid,SessionToken,Items}`; `NewCart`, `GroupBySeller()
      map[uuid][]Item`, `Merge(other)`, `TotalItemCount`.
- [ ] `.../domain/item.go` — `Item{SellerID,OutletID,VariantID,Quantity decimal,
      UnitPriceAmount(CACHE),PriceCheckedAt,Note}`; `SetQuantity(q,unitAllowsFraction,
      max)`, `IsPriceStale(now,ttl)`. `SellerID` didenormalisasi ke item.
- [ ] `.../domain/repository.go` — `CartRepository{FindActiveByUser,FindActiveBySession,
      Create,UpdateStatus,Delete}`, `ItemRepository{Upsert,UpdateQuantity,Delete,
      DeleteByCart,ListByCart,CountByCart}`.
- [ ] `.../domain/errors.go` — `ErrSellerCannotSell`, `ErrOutletInactive`,
      `ErrQuantityTooLarge`, `ErrFractionNotAllowed`, `ErrCartEmpty`, `ErrTooManyItems`.

### Migrasi
- [ ] `migrations/000007_create_cart.up.sql` / `.down.sql` — `cart_carts` (guest
      `session_token`, unique active per user), `cart_items` sesuai spec.

### Infra
- [ ] `.../infra/{cart_repository,item_repository,mapper}.go`.

### App
- [ ] `.../app/cart_usecase.go` — `GetCart`: (1) ambil item, (2) `catalog.GetVariants`
      batch, (3) `pricing.ResolveMany` batch, (4) `inventory.AvailableMany` batch, (5)
      `seller.GetMany` batch, (6) tandai baris bermasalah (harga berubah / stok menipis /
      diarsip / seller suspended), (7) kelompok per seller. Plus `AddItem/UpdateQuantity/
      RemoveItem/ClearCart`, `MergeGuestCart(sessionToken,userID)`. `AddItem` cek
      purchasable / can-sell / outlet aktif / qty wajar; **tidak menahan stok**.
- [ ] `.../app/outlet_selection_usecase.go` — `SuggestOutlet(sellerID,variantID,qty,
      buyerLat,buyerLng)`: `inventory.OutletsWithStock` → filter aktif + dukung ≥1 metode
      → terdekat bila ada koordinat → else outlet default seller → none →
      `ErrInsufficientStock`. `ReassignOutlets(cartID)`.
- [ ] `.../app/{service,dto,port_adapter,handlers}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,cart_handler,handler}.go` — `/api/v1`
      (butuh login untuk keranjang milik user; guest via `session_token`).

### Port
- [ ] `internal/modules/cart/port.go` — sempit: `GetActive(userID)`,
      `MarkConverted(cartID,orderID)`, `Clear(cartID)`. `Snapshot.Items` **tanpa harga**
      (variant/outlet/seller/qty/note saja).
- [ ] `events.go` — `EventItemAdded`, `EventCartConverted`, `EventCartAbandoned`.
- [ ] `module.go` — `Config{...,Catalog,Pricing,Inventory,Sellers}`; punya subscription
      (mis. `EventProductArchived` / `EventPriceChanged` → tandai baris).

### Wiring
- [ ] `internal/app/registry.go` — build `cart` setelah `pricing` + `inventory` +
      `seller`.

## Test wajib

- Keranjang 3 penjual → `GroupBySeller` menghasilkan 3 kelompok benar.
- `GetCart` = 4 batch call, bukan 4N (assert lewat port palsu penghitung).
- `MergeGuestCart`: varian+outlet sama → qty dijumlahkan, bukan ditimpa.
- `SuggestOutlet` memilih outlet terdekat yang stoknya cukup; semua kurang →
  `ErrInsufficientStock`.
- `AddItem` saat stok menipis → tetap masuk keranjang dengan flag, bukan error.

## Sengaja TIDAK dikerjakan di fase ini

- Ongkir (fase 8), voucher (fase 9), pembuatan order (fase 10).
- `cart.MarkConverted` dipanggil dari checkout (fase 10).
