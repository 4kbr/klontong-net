package app

// TODO — ExpireUnpaidOrders(ctx) (int, error): dipanggil worker.
//   1. ListExpired dengan FOR UPDATE SKIP LOCKED
//   2. untuk tiap pesanan: WithinTx
//        - order.Cancel("pembayaran kedaluwarsa")
//        - inventory.Release
//        - promotion.ReleaseRedemption — kuota voucher dikembalikan
//        - payment.MarkExpired
//        - outbox EventOrderExpired
//
// TODO — AutoCompleteDelivered(ctx) (int, error):
//   Suborder berstatus delivered yang sudah lewat masa komplain dipindahkan ke
//   completed. Inilah yang memicu pematangan dana penjual.
//
// TODO — AutoRejectUnconfirmed(ctx) (int, error):
//   Suborder yang tidak dikonfirmasi penjual dalam batas waktu ditolak otomatis.
//
// Ketiganya idempoten. Dijalankan dua kali tidak boleh melepas stok dua kali
// atau mengembalikan kuota voucher dua kali.
