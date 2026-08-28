package app

// TODO: Available, AvailableMany, OutletsWithStock.
//
// AvailableMany adalah query terpanas di sistem — dipanggil setiap kali
// keranjang dibuka dan setiap kali halaman produk dimuat. Satu query dengan
// WHERE (outlet_id, variant_id) IN (...), bukan loop.
//
// Pertimbangkan cache pendek di Redis untuk halaman produk, TAPI JANGAN untuk
// checkout. Checkout harus membaca angka sebenarnya di dalam transaksi dengan
// kunci baris.
