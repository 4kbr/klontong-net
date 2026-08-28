package domain

// TODO:
//   type ReservationStatus string   // held | committed | released | expired
//   type Reservation struct { ID; OrderID; SuborderID; OutletID; VariantID;
//                             Quantity decimal; Status; ExpiresAt;
//                             CreatedAt; UpdatedAt }
//   func NewReservation(...) *Reservation
//   func (r *Reservation) IsExpired(now time.Time) bool
//   func (r *Reservation) Commit() error    // hanya dari held
//   func (r *Reservation) Release() error   // hanya dari held
//
// KAPAN STOK DITAHAN: saat CHECKOUT, bukan saat masuk keranjang.
// Menahan sejak keranjang berarti barang hilang dari peredaran karena ada yang
// menimbunnya berminggu-minggu. Konsekuensinya, barang di keranjang bisa habis
// diambil orang lain — dan itu harus dikomunikasikan jelas di checkout, bukan
// jadi kejutan. Lihat ADR-003.
//
// Reservasi PUNYA masa berlaku dan worker yang melepasnya. Kalau worker tidak
// jalan, stok habis di layar padahal barangnya ada di rak, dan tidak ada satu
// pun pesan error yang muncul.
