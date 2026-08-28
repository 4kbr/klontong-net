package inventory

// TODO:
//   type Availability struct {
//       OutletID, VariantID uuid.UUID
//       Available decimal      // on_hand - reserved, dalam satuan DASAR
//   }
//   type ReserveItem struct { OutletID, VariantID uuid.UUID; BaseQuantity decimal }
//
//   type Port interface {
//       Available(ctx, outletID, variantID uuid.UUID) (decimal, error)
//       AvailableMany(ctx, pairs []OutletVariant) (map[OutletVariant]decimal, error)
//       OutletsWithStock(ctx, variantID uuid.UUID, minQty decimal) ([]uuid.UUID, error)
//
//       Reserve(ctx, orderID, suborderID uuid.UUID, items []ReserveItem) error
//       Commit(ctx, orderID uuid.UUID) error     // pesanan dibayar & dikirim
//       Release(ctx, orderID uuid.UUID) error    // pesanan batal
//   }
//
// SEMUA kuantitas di port ini dalam SATUAN DASAR varian, bukan satuan jual.
// Pemanggil yang mengonversi (lewat catalog.VariantInfo.ContentQuantity).
// Menaruh konversi di sini akan menyebarkan pengetahuan tentang satuan ke modul
// yang tidak perlu tahu.
//
// AvailableMany versi batch wajib: halaman keranjang memeriksa 15 baris
// sekaligus, dan checkout memeriksanya lagi.
//
// OutletsWithStock dipakai `order` untuk memilih outlet mana yang melayani
// pesanan. Lihat ADR-007.
//
// Reserve/Commit/Release dipanggil DI DALAM transaksi milik `order`. Ketiganya
// wajib memakai ExecutorFrom agar ikut transaksi itu.
