# Fase 03 — Storefront katalog (guest)

> Prasyarat: 01 | Ref: kontrak `paths/catalog-public.yaml`,
> `paths/sellers-public.yaml`, `paths/reviews.yaml#/productReviews`;
> `docs/ARCHITECTURE.md` §5 (OptionalAuth, SEO), ADR-006 (satuan).

## Tujuan

Permukaan belanja publik yang bisa dibuka **tanpa login**: beranda, kategori,
telusur/pencarian produk dengan filter + paginasi keyset, detail produk (varian
per satuan + tier grosir + indikator stok), dan halaman toko publik. **Selesai
kalau** tamu bisa menjelajah seluruh katalog dari mock, tiap halaman punya
`<title>/<meta>` dan produk punya JSON-LD.

## Aturan khusus fase ini

- Semua rute fase ini **publik** (`security: []`); jangan pasang `RequireAuth`.
- **Campur slug/id** (catatan kontrak): produk diakses via `slug`
  (`GET /api/v1/products/{slug}`) tapi review via `productId`
  (`GET /api/v1/products/{productId}/reviews`). Simpan `slug` **dan** `id` di
  tipe hasil fetch; jangan asumsikan sama.
- Harga ditampilkan **apa adanya dari server** per varian/tier. **Tidak menghitung
  "hemat X%" sendiri** kecuali server menyediakannya; kalau dihitung klien, beri
  label perkiraan dan pakai `bpsToPercent`.
- Satuan non-pecahan (dus, renceng) — tampilkan `content_quantity` & `base_unit`
  ("1 dus = 40 pcs"); pemilih qty menolak pecahan untuk satuan itu.
- Indikator stok: tampilkan "stok menipis"/"habis" dari flag server; **jangan**
  hitung ketersediaan sendiri. Stok bisa berubah — ini tampilan, bukan janji
  (ADR-003 dijelaskan penuh di Fase 04).
- Filter & pencarian masuk **URL query** (shareable, back-button benar).

## Urutan kerja

### Tipe & API (packages/api)
- [ ] `src/endpoints/catalog.api.ts` — `listCategories`, `browseProducts`
      (query: `limit,cursor,sort,category,seller,city,price_min,price_max,rating_min`),
      `searchProducts` (`/products/search`), `getProductBySlug` (`/products/{slug}`).
- [ ] `src/endpoints/seller-public.api.ts` — `getSeller` (`/sellers/{slug}`), `listSellerProducts` (`/sellers/{slug}/products`).
- [ ] `src/endpoints/review.api.ts` (parsial) — `listProductReviews` (`/products/{productId}/reviews`).
- [ ] `src/hooks/` — `useCategories`, `useProducts` (`useInfiniteQuery`), `useProductSearch`, `useProduct`, `useSellerPublic`, `useProductReviews`.
- [ ] `src/schemas/catalog.ts` — Zod parse ringan pada boundary detail produk (varian, tier) untuk menjaga asumsi UI.

### State (Zustand / Query)
- [ ] `useCatalogFilters()` — hook sinkron `URLSearchParams` ↔ state filter (tanpa store global).
- [ ] `queryKeys.products(filterKey)`, `queryKeys.product(slug)`, `queryKeys.categories`.

### Komponen (packages/ui)
- [ ] `ProductCard`, `PriceRange`/`PriceLabel`, `RatingStars`, `CategoryTree`,
      `FilterSidebar` (harga, rating, kota), `SortSelect`, `VariantPicker`,
      `QuantityStepper` (hormati satuan non-pecahan), `StockBadge`, `ImageGallery`,
      `Breadcrumbs`, `LoadMore` (keyset).

### Halaman & rute (storefront)
- [ ] `/` beranda — kategori unggulan + produk terbaru/populer.
- [ ] `/kategori` & `/kategori/:slug` — pohon kategori + hasil `browseProducts`.
- [ ] `/produk` — telusur + `FilterSidebar` + `SortSelect` + `LoadMore`.
- [ ] `/cari?q=` — `searchProducts`.
- [ ] `/produk/:slug` — galeri, `VariantPicker`, harga & tier, `StockBadge`,
      deskripsi, blok ulasan (ringkas + link), tombol "tambah ke keranjang"
      (stub → aktif di Fase 04). `PageHead` + JSON-LD `Product`/`Offer`.
- [ ] `/toko/:slug` — profil toko + produknya.

### Wiring
- [ ] MSW handlers `catalog` + `seller-public` + `product-reviews` dengan fixtures:
      ≥3 kategori bertingkat, ≥3 toko, produk dengan varian pcs/renceng/dus, sebagian stok menipis/habis.
- [ ] `PageHead` default di `RootLayout`; per-rute override.

## Test wajib

- `browseProducts` paginasi: `LoadMore` memakai `meta.next_cursor`; tidak ada duplikat; berhenti saat `has_more=false`.
- Filter harga/rating → tercermin di URL; reload mempertahankan hasil.
- Detail produk: varian dus menampilkan "1 dus = 40 pcs"; `QuantityStepper` menolak `0,5` untuk dus, izinkan untuk kg.
- Produk `slug` ≠ `id`: blok ulasan memanggil `/products/{id}/reviews`, bukan slug.
- Halaman produk merender `<title>` spesifik + `<script type="application/ld+json">` valid.
- Semua rute fase ini terbuka tanpa sesi (tidak ada redirect login).

## Sengaja TIDAK dikerjakan

- Tambah ke keranjang secara nyata, ongkir, checkout — Fase 04.
- Tulis ulasan — Fase 05.
- Prerender/SSR untuk SEO — Fase 10.
