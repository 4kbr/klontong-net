package settlement

// TODO:
//   type Balance struct { Pending, Available money.Amount }
//   type Port interface {
//       SellerBalance(ctx, sellerID uuid.UUID) (Balance, error)
//   }
//
// Sengaja sangat sempit. Modul lain tidak berhak menulis ke buku besar —
// pencatatan hanya terjadi lewat event, di dalam modul ini. Kalau ada modul
// yang butuh "menambah saldo penjual", jawabannya bukan menambah method di
// sini, melainkan menerbitkan event yang tepat.
