# Fase 15 — Review

> Prasyarat: fase 12 (suborder bisa `completed`). Guide:
> [`../GUIDES.md`](../GUIDES.md) §2 (tahap 15).
> Kontrak (**mengikat**, ADR-015): [`paths/reviews.yaml`](../../../../contracts/openapi/paths/reviews.yaml), [`seller-reviews.yaml`](../../../../contracts/openapi/paths/seller-reviews.yaml).

## Tujuan

Ulasan produk & penjual. **Selesai kalau** ulasan hanya bisa dari pembelian yang
**selesai** (suborder `completed`), satu ulasan per barang yang dibeli.

## Aturan khusus fase ini

- Verifikasi pembelian lewat **`order.Port.HasPurchased(userID, variantID) ->
  (orderItemID, ok)`** — **jangan percaya `order_item_id` dari klien**.
- Ulasan terikat ke **`order_item_id`** dan **unik** — satu ulasan per barang yang
  dibeli. Unique index sebagai penjaga terakhir.
- Review boleh dibuat setelah suborder **`completed`**, bukan setelah **`paid`**.
- `Rating` 1..5; komentar boleh kosong.
- `Edit` hanya dalam jendela waktu tertentu (`ErrEditWindowExpired`).
- Satu `Reply` per review, hanya oleh seller pemilik.
- Nama penulis di listing diambil **batch** via `identity.Port`, menghormati
  `IsAnonymous`.
- `catalog` mengonsumsi `EventReviewPublished` untuk memperbarui `RatingAvg` /
  `RatingCount` produk (agregat, lewat worker fase 13 — bukan per request).

## Urutan kerja

### Domain
- [ ] `internal/modules/review/internal/domain/review.go` — `Review{OrderItemID,
      ProductID,VariantID,SellerID,UserID,Rating int,Comment,Images []string,IsAnonymous,
      Status}`; `NewReview(orderItemID,rating,comment)` (rating 1..5); `Edit(rating,
      comment,within,now)`.
- [ ] `.../domain/reply.go` — `Reply{ReviewID,SellerID,Body}` — satu per review.
- [ ] `.../domain/repository.go` — `ReviewRepository{Create,Update,FindByID,
      FindByOrderItem,ListByProduct,ListBySeller,ListByUser,SummaryByProduct,
      SummariesByProducts}`, `ReplyRepository{Upsert,FindByReview}`.
- [ ] `.../domain/errors.go` — `ErrAlreadyReviewed`, `ErrNotPurchased`,
      `ErrOrderNotCompleted`, `ErrEditWindowExpired`, `ErrNotReviewAuthor`.

### Migrasi
- [ ] `migrations/000013_create_review.up.sql` / `.down.sql` — `review_reviews`
      (`order_item_id` **unique**, `rating` check 1..5), `review_replies` sesuai spec.

### Infra
- [ ] `.../infra/{review_repository,reply_repository,mapper}.go` —
      `SummariesByProducts` **batch**.

### App
- [ ] `.../app/review_usecase.go` — `CreateReview`: (1) `order.Port.HasPurchased`, (2)
      suborder `completed`, (3) belum pernah direview (unique index penjaga terakhir), (4)
      `WithinTx` simpan + outbox `EventReviewPublished`. Plus `ListProductReviews` (filter
      rating; nama penulis batch via `identity.Port`, hormati `IsAnonymous`),
      `ListMyReviews`, `EditReview`, `ReplyToReview`, `ReportReview`.
- [ ] `.../app/{service,dto,port_adapter}.go`.

### Transport
- [ ] `.../transport/rest/{routes,request,response,review_handler,handler}.go` — tulis
      ulasan di `/api/v1` (butuh login); baca ulasan produk publik (`OptionalAuth`);
      reply di `/api/v1/seller`.

### Port
- [ ] `internal/modules/review/port.go` — `ProductSummary(productID)`,
      `ProductSummaries` (**batch**). `Summary{RatingAvg,RatingCount}`.
- [ ] `events.go` — `EventReviewPublished`, `EventReviewHidden`.
- [ ] `module.go` — `Config{...,Orders order.Port,Users identity.Port,Storage}`.

### Wiring
- [ ] `internal/app/registry.go` — build `review` setelah `order` (terima
      `orderMod.Port()`).
- [ ] Pastikan `catalog` subscribe `EventReviewPublished` (fase 13) untuk update agregat
      rating.

## Test wajib

- `CreateReview` untuk barang yang tidak dibeli user → `ErrNotPurchased` (bukan percaya
  `order_item_id` klien).
- Review saat suborder masih `paid`/`delivered` (belum `completed`) →
  `ErrOrderNotCompleted`.
- Review kedua untuk `order_item_id` sama → `ErrAlreadyReviewed` (unique index).
- `ListProductReviews` mengambil nama penulis batch; review `IsAnonymous` tidak
  membocorkan nama.
- `EditReview` setelah jendela waktu → `ErrEditWindowExpired`.

## Sengaja TIDAK dikerjakan di fase ini

- Perhitungan agregat rating di produk (fase 13, worker katalog) — cukup terbitkan
  `EventReviewPublished`.
- Moderasi ulasan / alur `ReportReview` lanjutan (admin) — sediakan endpoint, proses
  moderasi menyusul.
