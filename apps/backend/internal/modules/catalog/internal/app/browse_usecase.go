package app

// TODO (pembeli):
//   ListProducts(ctx, BrowseQuery) — filter kategori, harga, rating, penjual,
//     kota; urut terbaru/termurah/terlaris/rating
//   GetProductDetail(ctx, slug) — produk + varian + foto + harga + ketersediaan
//   SearchProducts(ctx, keyword, ...)
//
// GetProductDetail perlu data dari `pricing` dan `inventory`. Ambil lewat port
// mereka secara BATCH untuk semua varian sekaligus, jangan per varian.
//
// Ketersediaan yang ditampilkan bergantung outlet mana yang melayani pembeli.
// Untuk halaman produk, cukup tampilkan "tersedia di N outlet" atau ketersediaan
// di outlet terdekat; penentuan outlet yang sebenarnya terjadi saat checkout.
