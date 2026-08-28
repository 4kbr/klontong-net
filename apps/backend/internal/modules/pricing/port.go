package pricing

// TODO:
//   type Resolved struct {
//       VariantID uuid.UUID
//       UnitPrice money.Amount      // harga satuan setelah tier diterapkan
//       CompareAt money.Amount      // harga coret, 0 kalau tidak ada
//       TierApplied bool
//       MinQuantityForNextTier decimal   // untuk pesan "beli 3 lagi, hemat Rp2.000"
//   }
//   type Query struct { VariantID, OutletID uuid.UUID; Quantity decimal }
//
//   type Port interface {
//       Resolve(ctx, q Query) (Resolved, error)
//       ResolveMany(ctx, qs []Query) (map[uuid.UUID]Resolved, error)
//   }
//
// Harga bergantung pada KUANTITAS, karena ada tier grosir. Jadi Resolve butuh
// tahu berapa yang dibeli, bukan hanya varian mana. Ini yang membedakannya dari
// "ambil harga produk" biasa.
//
// MinQuantityForNextTier bukan sekadar hiasan — memberi tahu pembeli bahwa
// menambah sedikit lagi akan lebih murah adalah cara menaikkan nilai keranjang,
// dan datanya sudah ada di sini.
//
// ResolveMany versi batch WAJIB. Checkout menghitung ulang harga seluruh
// keranjang, dan itu tidak boleh jadi N query.
