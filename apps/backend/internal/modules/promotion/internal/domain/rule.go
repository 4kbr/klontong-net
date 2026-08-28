package domain

// TODO — aturan kelayakan:
//   type ApplyContext struct {
//       UserID uuid.UUID
//       Subtotals map[uuid.UUID]money.Amount   // per suborder
//       ShippingAmounts map[uuid.UUID]money.Amount
//       SellerIDs []uuid.UUID
//       At time.Time
//   }
//   func (v *Voucher) EligibleFor(c ApplyContext) (basis money.Amount, err error)
//
// Aturan yang harus ditegakkan:
//   - voucher penjual hanya berlaku untuk suborder penjual TERSEBUT, bukan
//     seluruh keranjang. Ini kesalahan paling umum di marketplace: pembeli
//     memakai voucher toko A dan mendapat diskon untuk barang toko B.
//   - MinOrderAmount dihitung dari basis yang relevan (subtotal penjual itu
//     untuk voucher penjual; total keseluruhan untuk voucher marketplace)
//   - free_shipping hanya mengurangi ongkir, dan tidak lebih dari ongkirnya
