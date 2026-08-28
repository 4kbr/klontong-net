package domain

// TODO:
//   type Quote struct { ID; CartID, SuborderID *uuid.UUID; SellerID; OutletID;
//                       Method; CourierCode, CourierService string
//                       Amount money.Amount; ETDMinDay, ETDMaxDay int
//                       DistanceKm float64; ExpiresAt time.Time; CreatedAt }
//   func (q *Quote) IsValid(now time.Time) bool
//
// Quote punya masa berlaku (mis. 30 menit). Pembeli yang membuka checkout lalu
// pergi makan siang dan kembali harus melihat ongkir dihitung ulang, bukan
// memakai angka yang sudah basi.
