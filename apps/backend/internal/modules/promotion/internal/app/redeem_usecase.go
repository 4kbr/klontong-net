package app

// TODO — Redeem(ctx, voucherID, userID, orderID, perSuborder) error:
//   Dipanggil DI DALAM transaksi checkout.
//   1. FindByCodeForUpdate -> kunci baris voucher
//   2. periksa ulang kuota; keadaan bisa berubah sejak Validate
//   3. IncrementUsed
//   4. buat Redemption per suborder
//
// TODO — ReleaseRedemption(ctx, orderID): saat pesanan batal atau kedaluwarsa.
//   DecrementUsed dan hapus redemption. Idempoten — dipanggil dua kali tidak
//   boleh mengembalikan kuota dua kali.
