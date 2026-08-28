package main

// Background worker. Proses terpisah dari API.
//
// TODO — pekerjaan yang dijalankan, masing-masing goroutine dengan ctx yang sama:
//
//   1. outbox relay
//        Memungut outbox_events, publish ke event bus.
//        FOR UPDATE SKIP LOCKED, aman dijalankan banyak instance.
//
//   2. pelepasan reservasi stok kedaluwarsa
//        Reservasi yang lewat STOCK_RESERVATION_TTL dan belum dibayar dilepas,
//        stoknya kembali tersedia. Kalau ini tidak jalan, stok habis di layar
//        padahal barangnya ada — dan tidak ada error apa pun yang muncul.
//
//   3. pembatalan pesanan yang tidak dibayar
//        Pesanan menunggu bayar yang melewati PAYMENT_EXPIRY dibatalkan.
//
//   4. rekonsiliasi pembayaran
//        Menanyakan status ke gateway untuk pembayaran yang menggantung.
//        Webhook bisa hilang; jangan hanya mengandalkannya. Lihat ADR-008.
//
//   5. pematangan dana (settlement maturation)
//        Suborder yang selesai dan sudah lewat SETTLEMENT_HOLD_PERIOD
//        dipindahkan dari saldo tertahan ke saldo tersedia penjual.
//
//   6. pemrosesan pencairan (payout)
//        Membuat batch pencairan untuk penjual yang saldonya melewati
//        PAYOUT_MINIMUM, lalu mengirim perintah transfer.
//
//   7. sinkronisasi status pengiriman
//        Menarik tracking dari kurir untuk pengiriman yang masih berjalan.
//
//   8. penyegaran materialized view / cache katalog
//        Rating produk, jumlah terjual, dan agregat lain yang mahal dihitung
//        saat request.
//
// Setiap pekerjaan menghormati ctx.Done() dan berhenti rapi.
// Pekerjaan yang menyentuh uang WAJIB idempoten — dijalankan dua kali tidak
// boleh membayar penjual dua kali.

func main() {
	// TODO
}
