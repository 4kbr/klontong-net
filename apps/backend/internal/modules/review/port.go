package review

// TODO:
//   type Summary struct { ProductID uuid.UUID; RatingAvg float64; RatingCount int }
//   type Port interface {
//       ProductSummary(ctx, productID uuid.UUID) (Summary, error)
//       ProductSummaries(ctx, ids []uuid.UUID) (map[uuid.UUID]Summary, error)
//   }
//
// Meski catalog menyimpan agregatnya sendiri (diperbarui lewat event), port ini
// tetap berguna untuk halaman yang butuh angka paling mutakhir.
