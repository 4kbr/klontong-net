# Fase 05 — Pricing

> Prasyarat: fase 04. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 5), §10.
> ADR terkait: ADR-004 (harga dihitung server, jalur sama), ADR-005 (uang).

## Tujuan

Harga jual, tier grosir, harga per outlet. **Selesai kalau** harga per satuan dan tier
tampil benar.

## Aturan khusus fase ini

- **ADR-004**: jalur `Resolve` untuk tampilan produk dan untuk checkout **wajib kode yang
  sama**. Kalau beda, cepat atau lambat hasilnya beda dan pembeli yang menemukannya.
- Tier dihitung **per baris keranjang**, bukan dari total order.
- `ResolveTier`: ambil tier dengan `MinQuantity` **terbesar yang ≤ qty**. Tier
  divalidasi: kuantitas menaik, harga menurun.
- `Amount` > 0; `CompareAtAmount` bila diset harus > `Amount` (kalau tidak, coretan
  harga palsu).
- `CostAmount` hanya untuk laporan margin — **jangan pernah** masuk response yang dilihat
  pembeli.
- Uang `money.Amount` int64; persen tier tidak ada — tier menyimpan harga absolut.

## Urutan kerja

### Domain
- [ ] `internal/modules/pricing/internal/domain/price.go` — `Price{VariantID,OutletID
      *uuid,Amount,CompareAtAmount,CostAmount,Currency,StartsAt/EndsAt,IsActive,Tiers
      []QuantityTier}`; `NewPrice` (Amount>0, CompareAt>Amount bila diset);
      `IsEffectiveAt(t)`. `OutletID` nil = semua outlet.
- [ ] `.../domain/tier.go` — `QuantityTier{PriceID,MinQuantity decimal,Amount}`;
      `ResolveTier(base, tiers, qty) -> (price, tierUsed, minQtyForNext)`; aturan tier
      terbesar ≤ qty; validasi menaik-qty / menurun-harga.
- [ ] `.../domain/repository.go` — `PriceRepository{Upsert,FindByVariant,
      FindByVariantAndOutlet,FindManyEffective(pairs, at),Deactivate}`,
      `TierRepository{ReplaceForPrice,ListByPrice,ListByPrices}`. `FindManyEffective`
      ambil harga outlet-spesifik DAN umum dalam satu query.
- [ ] `.../domain/errors.go` — `ErrPriceNotSet`, `ErrCompareAtTooLow`,
      `ErrTierNotDescending`.

### Migrasi
- [ ] `migrations/000006_create_pricing.up.sql` / `.down.sql` — `pricing_prices`
      (amount `bigint`, per-outlet nullable), `pricing_tiers` (`min_quantity`
      `numeric(14,3)`, `amount` `bigint`) sesuai spec.

### Infra
- [ ] `.../infra/{price_repository,tier_repository,mapper}.go`.

### App
- [ ] `.../app/price_usecase.go` — seller dashboard: `SetPrice`, `SetTiers` (ganti semua
      tier sekaligus), `ListPrices`, `SchedulePrice`, `DeactivatePrice`.
- [ ] `.../app/resolve_usecase.go` — `Resolve` / `ResolveMany`: (1) harga efektif
      (outlet-spesifik dulu, lalu umum), (2) filter `IsEffectiveAt(now)`, (3) terapkan
      tier by quantity, (4) kembalikan `Resolved{UnitPrice,CompareAt,TierApplied,
      MinQuantityForNextTier}`. **Jalur yang sama** dipakai tampilan & invoice.
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,price_handler,handler}.go` — kelola
      harga di `/api/v1/seller`.

### Port
- [ ] `internal/modules/pricing/port.go` — `Resolve(Query{VariantID,OutletID,Quantity})`,
      `ResolveMany` (**batch wajib**).
- [ ] `events.go` — `EventPriceChanged` (consumer: cart menandai baris yang harganya
      berubah).
- [ ] `module.go` — `Config{...,Catalog catalog.Port}`.

### Wiring
- [ ] `internal/app/registry.go` — build `pricing` setelah `catalog` (terima
      `catalogMod.Port()`).

## Test wajib

- `ResolveTier`: qty di antara dua tier → ambil tier bawah; qty di bawah semua tier →
  harga dasar; `MinQuantityForNextTier` benar.
- Harga terjadwal (`StartsAt` masa depan) tidak dipakai oleh `Resolve(now)`.
- `ResolveMany` untuk 10 varian = 1 query, bukan 10.
- Harga outlet-spesifik menang atas harga umum untuk outlet itu.

## Sengaja TIDAK dikerjakan di fase ini

- Pemakaian `Resolve` di keranjang / checkout (fase 7 / 10) — cukup pastikan port siap.
- Diskon voucher (fase 9).
- Konsumen `EventPriceChanged` (fase 7).
