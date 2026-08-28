package app

// TODO (dasbor penjual):
//   ListSellerOrders(ctx, filter, cursor)
//   GetSellerOrder(ctx, suborderID)
//   ConfirmSuborder(ctx, suborderID)      -> awaiting_confirmation -> confirmed
//   RejectSuborder(ctx, suborderID, reason) -> melepas stok, memicu refund parsial
//   StartPacking(ctx, suborderID)
//   MarkShipped(ctx, suborderID, ShipInput) -> membuat shipment, commit stok
//   MarkReadyForPickup(ctx, suborderID)    -> menerbitkan kode ambil
//
// SEMUA transisi lewat suborder.Transition, tidak ada yang menyetel Status
// langsung.
//
// RejectSuborder adalah alur yang paling mudah salah: satu penjual menolak,
// dua penjual lain tetap berjalan. Yang harus terjadi:
//   - stok suborder ITU saja yang dilepas
//   - refund PARSIAL sebesar nilai suborder itu
//   - order induk TIDAK batal
//   - status order induk dihitung ulang dari anak-anaknya
// Menolak satu suborder lalu membatalkan seluruh pesanan adalah bug yang akan
// membuat dua penjual lain kehilangan penjualan tanpa sebab.
//
// Beri batas waktu konfirmasi. Suborder yang tidak dikonfirmasi dalam N jam
// ditolak otomatis oleh worker — pembeli tidak boleh menunggu tanpa kepastian.
