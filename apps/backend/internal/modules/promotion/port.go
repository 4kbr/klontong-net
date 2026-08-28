package promotion

// TODO:
//   type Applied struct {
//       VoucherID uuid.UUID
//       Code string
//       OwnerType string             // marketplace | seller
//       OwnerSellerID *uuid.UUID
//       AppliesTo string             // items | shipping
//       TotalDiscount money.Amount
//       PerSuborder map[uuid.UUID]money.Amount   // sudah dibagi
//   }
//   type Port interface {
//       Validate(ctx, code string, ctx ApplyContext) (Applied, error)
//       Redeem(ctx, voucherID, userID, orderID uuid.UUID,
//              perSuborder map[uuid.UUID]money.Amount) error
//       ReleaseRedemption(ctx, orderID uuid.UUID) error
//   }
//
// PerSuborder sudah TERBAGI saat dikembalikan. Pembagian memakai
// money.Distribute agar jumlahnya persis sama dengan TotalDiscount — sisa
// pembulatan dialokasikan secara deterministik, tidak dibuang.
//
// OwnerType menentukan SIAPA YANG MENANGGUNG diskon:
//   marketplace -> dipotong dari komisi kita
//   seller      -> dipotong dari pendapatan penjual
// Kalau ini tidak dibedakan, pembukuan settlement akan salah dan baru ketahuan
// saat penjual protes bagi hasilnya tidak sesuai.
//
// Redeem dipanggil DI DALAM transaksi checkout. Kuota dipotong di sana, bukan
// sebelumnya — kalau tidak, dua pembeli bisa memakai voucher terakhir yang sama.
