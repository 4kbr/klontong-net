# Fase 04 — Catalog

> Prasyarat: fase 03. Guide: [`../GUIDES.md`](../GUIDES.md) §2 (tahap 4), §6.
> ADR terkait: ADR-006 (stok satuan terkecil).
> Kontrak (**mengikat**, ADR-015): [`paths/catalog-public.yaml`](../../../../contracts/openapi/paths/catalog-public.yaml), [`seller-catalog.yaml`](../../../../contracts/openapi/paths/seller-catalog.yaml).

## Tujuan

Produk, varian, kategori, satuan. **Selesai kalau** bisa membuat produk dengan varian
dus / renceng / pcs.

## Aturan khusus fase ini

- **ADR-006**: varian simpan `unit_code` (satuan jual), `content_quantity` (isi),
  `base_unit_code` (satuan terkecil). Konversi lewat `Variant.ToBaseQuantity(qty) =
  qty * content_quantity`. `inventory` nanti hanya kenal satuan dasar.
- `content_quantity` wajib > 0; kalau `unit_code == base_unit_code` maka nilainya 1.
- `WeightGram` wajib > 0 (ongkir kurir dihitung dari berat).
- Satuan yang tidak boleh pecahan (dus, renceng) menolak kuantitas pecahan — ditegakkan
  di domain, bukan frontend.
- `Publish` menolak bila tidak ada varian aktif, tidak ada foto, atau harga belum diset.
- `RatingAvg` / `SoldCount` adalah agregat yang **diisi worker** (fase 13), tidak dihitung
  per request.
- Produk milik seller (belum ada master katalog bersama).
- Browse produk pakai area `/api/v1` dengan `OptionalAuth` — halaman produk & pencarian
  harus bisa dibuka tanpa login.

## Urutan kerja

### Domain
- [ ] `internal/modules/catalog/internal/domain/product.go` — `Product{SellerID,
      CategoryID,Slug,Status(draft|active|archived),IsTaxable,RatingAvg,RatingCount,
      SoldCount}`; `Publish(variantCount)`, `Archive`, `IsPurchasable`.
- [ ] `.../domain/variant.go` — `Variant{SKU,UnitCode(pcs|renceng|dus|kg),
      ContentQuantity decimal,BaseUnitCode,WeightGram,Length/Width/HeightCm,IsDefault,
      IsActive,Position}`; `NewVariant` (WeightGram>0, ContentQuantity>0, aturan ==1);
      `ToBaseQuantity(qty)`.
- [ ] `.../domain/category.go` — `Category{ParentID *uuid,Slug,Position,IsActive}`,
      `IsRoot`; simpan `path` (materialized path / ltree).
- [ ] `.../domain/unit.go` — `Unit{Code,Name,Symbol,IsDecimal}`; map `Units` (pcs/
      renceng/dus/kg/liter); `AllowsFraction`.
- [ ] `.../domain/repository.go` — `CategoryRepository` (.../ListTree),
      `ProductRepository` (.../Search(SearchQuery)/UpdateAggregates), `VariantRepository`
      (.../FindManyByID/CountActive), `ImageRepository` (Create/ListByProduct/Delete/
      Reorder).
- [ ] `.../domain/errors.go` — `ErrProductNotPurchasable`, `ErrProductHasNoVariant/
      NoImage`, `ErrSKUTaken`, `ErrInvalidUnit`, `ErrFractionNotAllowed`,
      `ErrWeightRequired`.

### Migrasi
- [ ] `migrations/000004_create_catalog.up.sql` / `.down.sql` — `catalog_categories`
      (tree/ltree), `catalog_products`, `catalog_variants`, `catalog_units`,
      `catalog_product_images` sesuai spec; unique index SKU; trigram index untuk search.

### Infra
- [ ] `.../infra/{product_repository,variant_repository,category_repository,
      image_repository,mapper}.go` — search mulai `ILIKE` + trigram (tsvector menyusul).

### App
- [ ] `.../app/product_usecase.go` — `Create/Update/Publish/Archive/ListMyProducts`;
      `PublishProduct` cek ≥1 varian aktif, ≥1 foto, harga diset.
- [ ] `.../app/variant_usecase.go` — `Create/Update/Deactivate/ListVariants`; nonaktif
      varian aktif terakhir → peringatan.
- [ ] `.../app/category_usecase.go` — `ListCategories` (tree, cached), `GetCategory`,
      admin CRUD.
- [ ] `.../app/image_usecase.go` — `RequestUploadURL` (presigned PUT), `AttachImage`,
      `DeleteImage`, `ReorderImages` (browser upload langsung ke object storage).
- [ ] `.../app/browse_usecase.go` — `ListProducts(BrowseQuery)` (filter kategori/harga/
      rating/seller/kota; sort newest/cheapest/best-selling/rating), `GetProductDetail
      (slug)` (butuh pricing + inventory, **batch per varian**), `SearchProducts`.
- [ ] `.../app/{service,dto,port_adapter,handlers}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,product_handler,variant_handler,
      browse_handler,image_handler,admin_handler,handler}.go` — browse di `/api/v1`
      (`OptionalAuth`), kelola produk di `/api/v1/seller`.

### Port
- [ ] `internal/modules/catalog/port.go` — `GetVariant`, `GetVariants` (**batch wajib**),
      `IsPurchasable`. `VariantInfo` bawa `UnitCode`, `BaseUnitCode`, `ContentQuantity`,
      `WeightGram`.
- [ ] `events.go` — `EventProductPublished/Archived/VariantUpdated/ProductDeleted`.
- [ ] `module.go` — `Config{...,Sellers seller.Port,Storage,Redis}`; punya subscription
      (agregat rating & sold, di-hook fase 13/15) + worker (refresh agregat, fase 13).

### Wiring
- [ ] `internal/app/registry.go` — build `catalog` setelah `seller` (terima
      `sellerMod.Port()`).

## Test wajib

- `Variant.ToBaseQuantity`: 2 dus isi 40 → 80; 1 renceng isi 10 → 10.
- Buat varian dus dengan `ContentQuantity` pecahan / `WeightGram=0` → ditolak.
- `PublishProduct` tanpa foto / tanpa harga → `ErrProductHasNoImage` / harga belum diset.
- `GetProductDetail` memanggil pricing & inventory secara batch (bukan N query).

## Sengaja TIDAK dikerjakan di fase ini

- Perhitungan harga / tier (fase 5).
- Stok & reservasi (fase 6).
- Worker refresh agregat rating/sold (fase 13) — sediakan `UpdateAggregates` di repo,
  jangan panggil per request.
