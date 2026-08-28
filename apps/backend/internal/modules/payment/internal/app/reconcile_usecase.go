package app

// TODO — ReconcilePending(ctx) (int, error): dipanggil worker.
//   Ambil pembayaran yang masih pending dan sudah cukup lama, lalu
//   gateway.Inquiry untuk mengetahui status sebenarnya.
//
// KENAPA INI WAJIB ADA: webhook hilang. Jaringan putus, server kita sedang
// deploy, gateway sedang bermasalah. Tanpa rekonsiliasi, ada pesanan yang sudah
// dibayar pembeli tapi selamanya berstatus menunggu pembayaran — dan yang
// menemukannya adalah pembeli yang marah, bukan monitoring kita.
//
// Jalankan juga rekonsiliasi harian yang membandingkan total transaksi kita
// dengan laporan settlement gateway. Selisih harus nol; kalau tidak, ada yang
// harus diselidiki hari itu juga.
