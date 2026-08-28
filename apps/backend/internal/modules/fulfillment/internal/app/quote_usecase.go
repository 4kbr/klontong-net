package app

// TODO — Quote(ctx, req QuoteRequest) ([]Option, error):
//   1. ambil outlet lewat seller.Port -> koordinat asal dan metode yang didukung
//   2. kumpulkan opsi:
//        pickup         kalau outlet mendukung; ongkir NOL
//        local_delivery kalau outlet mendukung, tujuan punya koordinat,
//                       jarak dalam radius, nilai pesanan di atas minimum
//        courier        kalau outlet mendukung; panggil CourierProvider.Rates
//   3. simpan semua opsi sebagai Quote dengan ExpiresAt
//   4. kembalikan; kalau KOSONG, itu error yang harus dijelaskan — pembeli perlu
//      tahu kenapa tidak ada satu pun cara mengirim barangnya
//
// Panggilan ke agregator kurir dibungkus timeout ketat. Kegagalannya TIDAK
// boleh menggagalkan seluruh quote — kembalikan opsi lain dan catat masalahnya.
//
// TODO — ValidateQuote(ctx, quoteID, at): dipanggil saat checkout.
//   Kedaluwarsa -> ErrQuoteExpired, dan checkout meminta pembeli memilih ulang.
