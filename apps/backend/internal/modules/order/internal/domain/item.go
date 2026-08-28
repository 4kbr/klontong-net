package domain

// TODO:
//   type Item struct {
//       ID; OrderID; SuborderID; VariantID; ProductID; SellerID
//       ProductName, VariantName, SKU, ImageURL string   // SALINAN
//       UnitCode string; ContentQuantity decimal          // SALINAN
//       Quantity decimal
//       UnitPriceAmount money.Amount                      // SALINAN
//       DiscountAmount, TotalAmount money.Amount
//       WeightGram int
//   }
//   func (i *Item) BaseQuantity() decimal    // Quantity * ContentQuantity
//   func (i *Item) TotalWeightGram() int
//
// SEMUA yang bertanda SALINAN adalah snapshot saat pesanan dibuat. Membuka
// pesanan tahun lalu harus menampilkan nama, harga, dan foto SAAT ITU.
// Produk bisa dihapus penjual dan harga bisa berubah; invoice tidak boleh
// ikut berubah. Lihat ADR-010.
//
// BaseQuantity dipakai untuk mengurangi stok: 2 dus isi 40 mengurangi 80.
