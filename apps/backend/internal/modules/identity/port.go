package identity

// TODO:
//   type UserInfo struct { ID uuid.UUID; Email, Phone, FullName, AvatarURL string; Roles []string }
//   type Port interface {
//       Get(ctx, id uuid.UUID) (UserInfo, error)
//       GetMany(ctx, ids []uuid.UUID) (map[uuid.UUID]UserInfo, error)
//       HasRole(ctx, id uuid.UUID, role string) (bool, error)
//       GrantRole(ctx, id uuid.UUID, role string) error
//   }
//
// GetMany versi batch WAJIB ada. Daftar ulasan produk menampilkan nama penulis;
// tanpa batch itu jadi 20 query untuk satu halaman.
//
// GrantRole ada karena pendaftaran penjual harus menambahkan peran `seller` ke
// akun yang sudah ada. Pembeli yang membuka toko tetap satu akun, bukan dua.
