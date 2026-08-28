package money

// TODO: aturan pembulatan, ditulis eksplisit di satu tempat.
//
//   func RoundDown(v int64, unit int64) int64
//   func RoundNearest(v int64, unit int64) int64
//
// Kapan dipakai:
//   - komisi marketplace  -> bulatkan ke bawah, marketplace tidak boleh
//                            mengambil lebih dari haknya karena pembulatan
//   - diskon              -> bulatkan ke bawah
//   - ongkir              -> ikuti angka dari penyedia apa adanya
//
// Tulis alasannya di komentar setiap kali menambah aturan baru. Pembulatan yang
// tidak terdokumentasi akan diubah orang berikutnya yang mengira itu bug.
