package app

// TODO — Validate(ctx, code string, actx ApplyContext) (Applied, error):
//   1. cari voucher, periksa aktif, periode, kuota global dan per pengguna
//   2. tentukan basis: untuk voucher penjual, hanya subtotal penjual itu
//   3. periksa MinOrderAmount terhadap basis
//   4. hitung diskon
//   5. BAGI ke suborder yang berhak dengan money.Distribute
//
// Langkah 5 adalah yang paling mudah salah. Voucher marketplace Rp10.000 untuk
// keranjang berisi tiga penjual dengan subtotal berbeda harus dibagi
// proporsional, dan jumlah tiga bagian itu harus PERSIS Rp10.000. Sisa
// pembulatan dialokasikan deterministik, tidak dibuang dan tidak dibulatkan
// ke atas semua.
//
// Validate tidak mengubah apa pun. Ia dipanggil berkali-kali saat pembeli
// mencoba kode di halaman checkout.
