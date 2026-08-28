package app

// TODO — CreateReview(ctx, CreateReviewInput):
//   1. order.Port.HasPurchased -> pastikan pembeli memang membeli barang ini
//   2. pastikan suborder-nya sudah completed
//   3. pastikan belum pernah diulas (unique index sebagai penjaga terakhir)
//   4. WithinTx: simpan + outbox EventReviewPublished
//
// TODO: ListProductReviews (dengan filter rating), ListMyReviews, EditReview,
//       ReplyToReview, ReportReview.
//
// ListProductReviews perlu nama penulis — ambil BATCH lewat identity.Port,
// dan hormati IsAnonymous.
