package catalog

// TODO:
//   type VariantInfo struct {
//       ID, ProductID, SellerID uuid.UUID
//       ProductName, VariantName, SKU, ImageURL string
//       UnitCode, BaseUnitCode string
//       ContentQuantity decimal        // isi per satuan: dus = 40 pcs
//       WeightGram int
//       IsActive bool
//   }
//   type Port interface {
//       GetVariant(ctx, variantID uuid.UUID) (VariantInfo, error)
//       GetVariants(ctx, ids []uuid.UUID) (map[uuid.UUID]VariantInfo, error)
//       IsPurchasable(ctx, variantID uuid.UUID) (bool, error)
//   }
//
// GetVariants versi batch WAJIB. Keranjang berisi 15 baris tidak boleh berarti
// 15 query, dan checkout membacanya lagi.
//
// ContentQuantity dan BaseUnitCode ikut di sini karena `inventory` dan `order`
// membutuhkannya untuk menghitung pengurangan stok. Satu dus terjual mengurangi
// stok sebanyak isinya, bukan satu. Lihat ADR-006.
