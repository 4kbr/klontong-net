package app

// TODO — Register(ctx, RegisterInput) (TokenPair, error):
//   1. validasi kekuatan password
//   2. normalisasi email dan nomor HP
//   3. cek keberadaan; race tetap mungkin, andalkan unique index sebagai
//      penjaga terakhir dan terjemahkan 23505 jadi ErrEmailTaken
//   4. hash password, buat user dengan peran `buyer`
//   5. WithinTx: simpan user + buat verifikasi email + outbox event
//   6. terbitkan token
//
// Verifikasi email TIDAK memblokir belanja. Memaksa verifikasi sebelum boleh
// menambah ke keranjang adalah cara kehilangan pembeli. Batasi hanya untuk aksi
// yang perlu: menarik dana, membuka toko.
