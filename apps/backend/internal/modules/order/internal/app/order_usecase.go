package app

// TODO (pembeli):
//   ListMyOrders(ctx, filter, cursor)
//   GetMyOrder(ctx, orderID) — beserta seluruh suborder dan statusnya
//   CancelOrder(ctx, orderID, reason) — hanya kalau belum ada yang dikirim
//   ConfirmReceived(ctx, suborderID) — pembeli menyatakan barang diterima,
//     memindahkan suborder ke completed lebih cepat dari batas otomatis
//
// GetMyOrder menampilkan satu pesanan dengan beberapa suborder yang statusnya
// bisa berbeda-beda. Sampaikan itu dengan jujur: "1 dari 3 toko sudah dikirim"
// lebih berguna daripada satu status tunggal yang menyederhanakan kenyataan.
