package domain

// TODO:
//   type MemberRole string   // owner | manager | staff
//   type Member struct { SellerID; UserID; Role; CreatedAt }
//   func (r MemberRole) CanManageProducts() bool
//   func (r MemberRole) CanManagePayout() bool   // hanya owner
//   func (r MemberRole) CanProcessOrders() bool
//
// Rekening pencairan hanya boleh diubah owner. Ini titik paling rawan
// penyalahgunaan kalau staf toko bisa mengubahnya.
