package postgres

// TODO: func Translate(err error) error
//   pgx.ErrNoRows        -> errs.ErrNotFound
//   23505 unique         -> errs.ErrConflict (sertakan nama constraint —
//                           di sistem ini nama constraint sering menjelaskan
//                           persis apa yang bentrok: voucher terpakai, pesanan
//                           ganda, pembayaran dobel)
//   23503 foreign key    -> errs.ErrConflict
//   23514 check          -> errs.ErrInvalidInput
//   40001 serialization  -> errs.ErrConflict dengan Retryable = true
//   55P03 lock_not_available -> errs.ErrConflict, retryable
//
// Repository memanggil ini sebelum mengembalikan error.
