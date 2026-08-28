package infra

// TODO: implementasi OrderRepository.
//
// Create menyimpan order + suborder + item dalam satu transaksi. Pakai insert
// batch untuk item, bukan loop — pesanan berisi 30 baris bukan hal aneh.
//
// NextOrderNumber memakai sequence Postgres. SELECT MAX + 1 akan menghasilkan
// nomor ganda saat dua checkout terjadi bersamaan, dan nomor pesanan ganda
// adalah masalah yang baru ketahuan saat customer service kebingungan.
