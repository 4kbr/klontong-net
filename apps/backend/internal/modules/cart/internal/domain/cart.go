package domain

// TODO:
//   type Status string   // active | converted | abandoned
//   type Cart struct { ID; UserID *uuid.UUID; SessionToken string; Status;
//                      Items []Item; CreatedAt; UpdatedAt }
//   func NewCart(userID *uuid.UUID, sessionToken string) *Cart
//   func (c *Cart) GroupBySeller() map[uuid.UUID][]Item
//   func (c *Cart) Merge(other *Cart) error
//   func (c *Cart) TotalItemCount() decimal
//
// GroupBySeller adalah operasi yang dipakai di hampir setiap tampilan keranjang
// DAN yang menentukan pembagian suborder saat checkout. Menaruhnya di domain
// membuat aturan pengelompokannya satu tempat, bukan tersebar di handler.
//
// Merge dipakai saat tamu login: keranjang tamu digabung dengan keranjang yang
// sudah ada, bukan menimpanya. Baris yang sama (varian + outlet) kuantitasnya
// dijumlahkan, dengan batas maksimum per baris.
