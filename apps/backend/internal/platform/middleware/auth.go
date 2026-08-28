package middleware

// TODO:
//   func Authenticate(verifier auth.TokenVerifier) func(http.Handler) http.Handler
//     Wajib login. Token invalid -> 401 dengan code yang MEMBEDAKAN
//     "token_expired" dari "token_invalid", supaya frontend tahu kapan harus
//     refresh dan kapan harus melempar ke halaman login.
//
//   func OptionalAuth(verifier) func(http.Handler) http.Handler
//     Untuk katalog dan detail produk: kalau ada token dipakai (agar bisa
//     menampilkan harga khusus dan status wishlist), kalau tidak ada tetap
//     lanjut sebagai tamu.
//
// Middleware ini menjawab "siapa kamu". "Boleh apa" diputuskan di usecase,
// karena butuh tahu toko atau pesanan mana yang sedang diakses.
