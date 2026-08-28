package app

// TODO — RequestRefund(ctx, RefundInput) (uuid.UUID, error):
//   1. periksa jumlah: total refund tidak boleh melebihi pembayaran.
//      Hitung SumByPayment di dalam transaksi, bukan sebelumnya.
//   2. buat Refund status requested
//   3. worker yang memanggil gateway.Refund
//
// Refund parsial adalah kasus NORMAL di marketplace, bukan pengecualian:
// satu penjual menolak pesanan, dua lainnya jalan terus.
//
// Untuk COD, refund berarti uang belum pernah masuk ke kita — yang perlu
// dilakukan hanya membatalkan tagihan. Bedakan alurnya dengan jelas.
