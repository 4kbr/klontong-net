package httpx

// TODO:
//   type Validatable interface { Validate() error }
//   func DecodeAndValidate[T Validatable](r) (T, error)
//
// Validasi ditulis tangan di tiap struct request, hanya soal BENTUK.
// Aturan bisnis ("stok tidak cukup", "voucher tidak berlaku untuk penjual ini")
// ada di domain atau usecase — bukan di sini.
