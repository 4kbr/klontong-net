package app

// TODO — handler event:
//   OnOrderPaid       -> tidak melakukan apa-apa pada stok (sudah direservasi
//                        saat checkout); ada di sini sebagai penanda eksplisit
//   OnSuborderShipped -> Commit reservasi suborder tersebut
//   OnOrderCancelled  -> Release reservasi
//   OnSuborderRejected-> Release reservasi suborder itu saja, bukan seluruh order
//
// Perhatikan yang terakhir: satu penjual menolak pesanan tidak boleh melepas
// stok penjual lain di pesanan yang sama.
//
// Idempoten. Reservasi yang sudah committed tidak boleh di-commit lagi.
