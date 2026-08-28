package domain

// TODO:
//   type QuantityTier struct { ID; PriceID; MinQuantity decimal; Amount money.Amount }
//   func ResolveTier(base money.Amount, tiers []QuantityTier, qty decimal) (money.Amount, bool, decimal)
//         kembalikan: harga terpilih, apakah tier dipakai, kuantitas minimum
//         untuk tier berikutnya
//
// Aturan: pilih tier dengan MinQuantity terbesar yang masih <= qty.
// Tier harus divalidasi menaik secara kuantitas dan MENURUN secara harga —
// tier yang lebih banyak tapi lebih mahal adalah kesalahan input, dan menolaknya
// saat disimpan jauh lebih baik daripada membiarkan pembeli menemukannya.
//
// Tier dihitung per BARIS keranjang, bukan dari total pesanan. Membeli 5 sabun
// dan 5 sampo tidak membuat keduanya dapat harga grosir 10 pcs.
//
// Ini ciri khas dagang klontong dan harus ada sejak awal.
