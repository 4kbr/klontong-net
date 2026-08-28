package app

// TODO — SendPending(ctx) (int, error): worker.
//   Ambil notifikasi yang belum terkirim, susun dari template, kirim lewat
//   channel yang sesuai preferensi penerima.
//   Gagal -> catat alasan, coba lagi dengan backoff, menyerah setelah batas.
