package app

// Usecase paling kritikal di modul ini.
//
// TODO — Reserve(ctx, orderID, suborderID uuid.UUID, items []ReserveItem) error:
//   1. URUTKAN items berdasarkan (outlet_id, variant_id) — cegah deadlock
//   2. GetForUpdate untuk semua baris stok sekaligus
//   3. periksa CanReserve untuk setiap item; kalau ada yang gagal, kembalikan
//      ErrInsufficientStock yang menyebut varian mana dan berapa yang tersedia
//   4. Reserve pada tiap stok, simpan
//   5. buat baris Reservation dengan ExpiresAt = now + TTL
//
//   Dipanggil DI DALAM transaksi milik `order`. Tidak membuka transaksi sendiri.
//
// TODO — Commit(ctx, orderID): ubah reservasi held -> committed, stok
//   OnHand dan Reserved dikurangi, tulis Movement kind=sale.
//
// TODO — Release(ctx, orderID): reservasi held -> released, Reserved dikurangi.
//   TIDAK menulis Movement — barangnya tidak pernah keluar dari rak.
//
// TODO — ReleaseExpired(ctx) (int, error): dipanggil worker.
//   ListExpired dengan FOR UPDATE SKIP LOCKED, lalu Release per batch.
//   Kalau ini tidak jalan, stok tertahan selamanya untuk pesanan yang sudah mati.
