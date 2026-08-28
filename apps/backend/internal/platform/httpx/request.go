package httpx

// TODO:
//   func DecodeJSON[T any](r) (T, error)
//     tolak content-type selain json, MaxBytesReader, DisallowUnknownFields,
//     error decode -> errs.Invalid bukan 500
//   func URLParamUUID(r, name) (uuid.UUID, error)
//   func QueryInt / QueryBool / QueryString
//   func Pagination(r) (limit int, cursor string)
//
// Daftar produk dan pesanan memakai cursor keyset, bukan offset. Katalog
// berubah terus; offset menghasilkan barang terlewat atau muncul dua kali saat
// pembeli menggulir.
