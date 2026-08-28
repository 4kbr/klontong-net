package domain

// TODO: ErrOrderNotFound, ErrSuborderNotFound, ErrNotOrderOwner,
//       ErrInvalidTransition, ErrOrderNotCancellable, ErrCartEmpty,
//       ErrPriceChanged, ErrOutOfStock, ErrAddressRequired,
//       ErrNoFulfillmentAvailable, ErrTotalMismatch, ErrOrderExpired
//
// ErrPriceChanged dan ErrOutOfStock harus membawa DETAIL per baris: barang mana,
// harga lama berapa, harga baru berapa, tersedia berapa. Checkout yang gagal
// dengan pesan "ada yang berubah" memaksa pembeli menebak, dan biasanya mereka
// menyerah.
//
// ErrTotalMismatch adalah penjaga internal: kalau jumlah suborder tidak sama
// dengan total order, itu bug kita dan pesanan tidak boleh dibuat.
