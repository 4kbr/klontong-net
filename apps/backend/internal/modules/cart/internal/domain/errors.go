package domain

// TODO: ErrCartNotFound, ErrItemNotFound, ErrVariantNotPurchasable,
//       ErrSellerCannotSell, ErrOutletInactive, ErrQuantityTooLarge,
//       ErrFractionNotAllowed, ErrCartEmpty, ErrTooManyItems
//
// ErrSellerCannotSell muncul saat penjual di-suspend sementara barangnya sudah
// ada di keranjang orang. Ini kondisi yang PASTI terjadi dan harus ditangani
// dengan pesan yang jelas, bukan 500.
