package app

// TODO — Login(ctx, LoginInput) (TokenPair, error):
//   - terima email ATAU nomor HP di field yang sama, deteksi bentuknya
//   - user tidak ada -> TETAP jalankan Compare terhadap hash dummy lalu balas
//     ErrInvalidCredential (cegah timing attack)
//   - user suspended -> tolak dengan pesan berbeda; ini bukan rahasia
//   - terbitkan access + refresh, simpan hash refresh beserta user agent & IP
//
// TODO: Logout, LogoutAll.
