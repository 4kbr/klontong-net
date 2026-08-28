package app

// TODO:
//   func (s *Service) requireOrderOwner(ctx, orderID uuid.UUID) (*domain.Order, error)
//   func (s *Service) requireSuborderSeller(ctx, suborderID uuid.UUID) (*domain.Suborder, error)
//         penjual hanya boleh menyentuh suborder miliknya sendiri
//   func (s *Service) requireAdmin(ctx) error
//
// SETIAP usecase publik diawali salah satunya. Tanpa ini, mengganti id di URL
// cukup untuk membaca pesanan orang lain — jenis kebocoran yang paling sering
// ditemukan di sistem e-commerce.
