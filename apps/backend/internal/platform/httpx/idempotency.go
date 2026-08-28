package httpx

// Kunci idempotensi untuk endpoint yang menciptakan sesuatu bernilai uang.
//
// TODO:
//   func Idempotent(store IdempotencyStore, ttl time.Duration) func(http.Handler) http.Handler
//     1. baca header Idempotency-Key; wajib untuk POST /checkout dan sejenisnya
//     2. kalau kunci sudah pernah dipakai, kembalikan response yang TERSIMPAN,
//        jangan jalankan handler lagi
//     3. kalau belum, jalankan handler, simpan response, lalu balas
//     4. permintaan dengan kunci sama yang masih BERJALAN dibalas 409
//
//   type IdempotencyStore interface {
//       Begin(ctx, key string, ttl time.Duration) (existing *StoredResponse, acquired bool, err error)
//       Finish(ctx, key string, resp StoredResponse) error
//   }
//
// Kenapa ini ada: pembeli menekan tombol bayar dua kali, jaringan timeout lalu
// aplikasi mengulang, atau tombol back di browser. Tanpa idempotensi, satu niat
// beli jadi dua pesanan dan dua tagihan. Lihat ADR-008.
//
// Kunci disimpan di Redis untuk kecepatan, TAPI perlindungan sesungguhnya ada
// di unique index database (nomor pesanan, idempotency_key pembayaran). Redis
// boleh hilang; database tidak.
