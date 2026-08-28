package app

// TODO — CreatePayment(ctx, orderID, amount, method, channel, idempotencyKey):
//   1. cek idempotency key; sudah ada -> kembalikan hasil lama
//   2. method == cod -> buat Payment status pending, TANPA memanggil gateway,
//      kembalikan instruksi "bayar saat barang diterima"
//   3. method == gateway -> gateway.Charge, simpan ProviderReference
//   4. simpan Payment dengan ExpiredAt = now + Expiry
//
// Panggilan ke gateway TIDAK di dalam transaksi milik checkout. Ia dilakukan
// setelah commit, atau dengan pola dua langkah: simpan pending, panggil gateway,
// perbarui referensi. Memanggil pihak luar di dalam transaksi berarti menahan
// koneksi database selama mereka berpikir.
//
// TODO: GetByOrder, CancelPayment, MarkCODCollected.
//
// MarkCODCollected dipanggil saat kurir menyetorkan uang. Ini yang membuat
// pembayaran COD berstatus settled, dan waktunya bisa berhari-hari setelah
// barang diterima.
