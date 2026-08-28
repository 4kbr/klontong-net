package app

// TODO — SyncTracking(ctx) (int, error): dipanggil worker.
//   Untuk setiap pengiriman kurir yang masih berjalan, panggil Track dan simpan
//   event baru. Saat status menjadi delivered, publish EventShipmentDelivered.
//
// Jangan menarik terlalu sering — agregator punya batas kuota. Interval yang
// masuk akal: 30 menit untuk yang baru dikirim, lebih jarang untuk yang sudah
// lama berjalan.
//
// Simpan event tracking apa adanya. Saat pembeli protes barang belum sampai
// padahal berstatus terkirim, riwayat inilah buktinya.
