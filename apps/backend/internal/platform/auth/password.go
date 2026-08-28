package auth

// TODO:
//   type Hasher interface { Hash(plain string) (string, error); Compare(hash, plain string) error }
//   func NewBcryptHasher(cost int) Hasher    // 12 produksi, 4 test
//
// Compare HARUS tetap berjalan walau user tidak ditemukan (bandingkan dengan
// hash dummy), supaya waktu respons tidak membocorkan email mana yang terdaftar.
