package app

// TODO (dasbor penjual):
//   ListStock(ctx, outletID, cursor) — daftar stok satu outlet
//   AdjustStock(ctx, outletID, variantID, delta, note) — koreksi manual
//   SetLowStockThreshold
//   TransferStock(ctx, fromOutletID, toOutletID, variantID, qty) — antar outlet
//
// AdjustStock WAJIB menulis Movement dengan catatan dan pelakunya. Koreksi stok
// tanpa jejak adalah lubang yang akan dipakai untuk menutupi kehilangan barang.
//
// TransferStock adalah dua movement dalam satu transaksi: transfer_out di asal,
// transfer_in di tujuan. Jumlahnya harus sama; kalau tidak, ada barang yang
// menguap.
