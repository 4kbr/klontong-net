package domain

// TODO: ErrPriceNotFound, ErrPriceNotSet, ErrInvalidPrice, ErrInvalidTier,
//       ErrCompareAtTooLow, ErrTierNotDescending
//
// ErrPriceNotSet muncul saat varian aktif tapi belum punya harga. Produk
// seperti itu tidak boleh bisa dipublikasikan — tapi periksa lagi saat checkout,
// karena harga bisa dinonaktifkan setelah produk tayang.
