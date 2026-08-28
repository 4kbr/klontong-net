package app

// TODO:
//   GetCart(ctx) (CartView, error)
//     Ini bukan sekadar membaca tabel. Setiap kali keranjang dibuka:
//       1. ambil item
//       2. catalog.GetVariants BATCH  -> nama, foto, status
//       3. pricing.ResolveMany BATCH  -> harga terkini beserta tier
//       4. inventory.AvailableMany BATCH -> ketersediaan
//       5. seller.GetMany BATCH       -> nama toko, status
//       6. tandai baris bermasalah: harga berubah, stok kurang, produk
//          diarsipkan, penjual suspended
//       7. kelompokkan per penjual
//     Empat panggilan batch, bukan 4N. Ini endpoint yang dibuka berkali-kali
//     sebelum checkout dan performanya terasa langsung.
//
//   AddItem / UpdateQuantity / RemoveItem / ClearCart
//   MergeGuestCart(ctx, sessionToken string, userID uuid.UUID)
//
// AddItem memeriksa: varian bisa dibeli, penjual bisa berjualan, outlet aktif,
// kuantitas masuk akal untuk satuannya. TIDAK memeriksa stok secara mengikat —
// stok baru ditahan saat checkout. Boleh memberi peringatan kalau stok menipis,
// tapi jangan menolak; stok bisa bertambah sebelum pembeli checkout.
