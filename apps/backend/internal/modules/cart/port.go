package cart

// TODO:
//   type Item struct {
//       VariantID, OutletID, SellerID uuid.UUID
//       Quantity decimal
//       Note string
//   }
//   type Snapshot struct { CartID uuid.UUID; Items []Item }
//
//   type Port interface {
//       GetActive(ctx, userID uuid.UUID) (Snapshot, error)
//       MarkConverted(ctx, cartID, orderID uuid.UUID) error
//       Clear(ctx, cartID uuid.UUID) error
//   }
//
// Sengaja sempit. `order` mengambil ISI keranjang lalu menghitung ulang
// semuanya sendiri — harga, ketersediaan, ongkir. Snapshot ini hanya membawa
// "apa yang dipilih pembeli", bukan angka apa pun.
//
// Tidak ada harga di Item. Itu disengaja: harga yang mengikat dihitung saat
// checkout, bukan diambil dari keranjang. Lihat ADR-004.
