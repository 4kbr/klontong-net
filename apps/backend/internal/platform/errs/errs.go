package errs

// Tipe error lintas modul.
//
// TODO — sentinel:
//   ErrNotFound, ErrConflict, ErrInvalidInput, ErrUnauthorized, ErrForbidden,
//   ErrTooManyRequests, ErrUpstream, ErrInternal
//
//   ErrUpstream untuk kegagalan payment gateway dan penyedia ongkir.
//   Membedakannya dari ErrInternal penting: yang satu salah kita, yang satu
//   salah pihak lain, dan penanganannya berbeda.
//
// TODO — error kaya:
//   type Error struct {
//       Kind error; Code string; Message string
//       Fields map[string]string
//       Retryable bool
//       cause error
//   }
//   func (e *Error) Error() string
//   func (e *Error) Unwrap() error   -> kembalikan e.Kind agar errors.Is bekerja
//
// TODO — konstruktor: NotFound, Conflict, Invalid, Unauthorized, Forbidden,
//        Upstream, Internal
//
// Kode error harus SPESIFIK dan stabil, karena frontend bercabang berdasarkan
// itu. `out_of_stock` berbeda penanganannya dari `price_changed`, dan keduanya
// berbeda dari `voucher_expired` — ketiganya muncul di layar checkout dengan
// pesan dan tombol yang berbeda.
