package app

// Inti modul ini. Satu peristiwa sering menghasilkan DUA notifikasi berbeda.
//
// TODO:
//   OnOrderPlaced        -> pembeli: "pesanan dibuat, selesaikan pembayaran"
//   OnOrderPaid          -> pembeli: "pembayaran diterima"
//                           penjual: "ada pesanan baru, siapkan barang"
//   OnSuborderConfirmed  -> pembeli: "toko X memproses pesananmu"
//   OnSuborderRejected   -> pembeli: "toko X menolak, dana dikembalikan"
//   OnSuborderShipped    -> pembeli: nomor resi atau kode ambil
//   OnSuborderDelivered  -> pembeli: "barang sampai, konfirmasi penerimaan"
//   OnOrderExpired       -> pembeli: "pesanan dibatalkan karena belum dibayar"
//   OnLowStock           -> penjual: "stok X menipis di outlet Y"
//   OnPayoutCompleted    -> penjual: "dana sudah ditransfer"
//   OnReviewPublished    -> penjual: "ada ulasan baru"
//
// "Pesanan dibayar" berarti hal berbeda bagi pembeli dan penjual. Satu event,
// dua pesan, dua penerima.
//
// IDEMPOTEN WAJIB. Notifikasi ganda adalah keluhan yang cepat membuat orang
// mematikan notifikasi sama sekali.
//
// Pengiriman email TIDAK dilakukan di handler. Handler menyimpan notifikasi;
// worker yang mengirim. SMTP lambat tidak boleh memperlambat relay outbox.
