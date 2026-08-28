package app

// TODO — handler event untuk memperbarui agregat:
//   OnReviewPublished -> hitung ulang rating produk
//   OnOrderCompleted  -> tambah sold_count
//   OnSellerSuspended -> sembunyikan produk penjual tersebut dari hasil browse
//
// Idempoten. Untuk agregat, cara paling aman adalah menghitung ULANG dari
// sumbernya (SELECT avg, count FROM reviews WHERE ...) alih-alih menambah satu.
// Penambahan inkremental akan menyimpang setiap kali ada event ganda.
