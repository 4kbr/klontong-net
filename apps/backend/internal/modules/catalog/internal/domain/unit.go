package domain

// TODO:
//   type Unit struct { Code, Name, Symbol string; IsDecimal bool }
//   var Units = map[string]Unit{ "pcs":..., "renceng":..., "dus":..., "kg":..., "liter":... }
//   func (u Unit) AllowsFraction() bool
//
// AllowsFraction penting: 0,5 kg beras masuk akal; 0,5 dus tidak. Validasi
// kuantitas di keranjang harus menolak pecahan untuk satuan yang tidak
// mengizinkannya, dan itu diputuskan di sini bukan di frontend.
