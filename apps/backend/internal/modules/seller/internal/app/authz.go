package app

// TODO:
//   func (s *Service) requireSellerMember(ctx, sellerID uuid.UUID) (*domain.Member, error)
//   func (s *Service) requireSellerRole(ctx, sellerID uuid.UUID, min domain.MemberRole) error
//
// SETIAP usecase dasbor penjual diawali salah satunya. Peran `seller` di token
// hanya berarti orang ini punya toko — bukan berarti ia berhak mengubah toko
// yang sedang diakses. Tanpa pemeriksaan ini, penjual A bisa mengubah produk
// penjual B hanya dengan mengganti id di URL.
