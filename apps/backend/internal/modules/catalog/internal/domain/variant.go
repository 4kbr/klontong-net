package domain

// Entity paling penting di modul ini, dan yang paling khas dagang klontong.
//
// TODO:
//   type Variant struct {
//       ID; ProductID; SKU; Name; Barcode
//       UnitCode string             // satuan jual: pcs | renceng | dus | kg
//       ContentQuantity decimal     // isi dalam satuan dasar: dus = 40
//       BaseUnitCode string         // satuan terkecil produk ini: pcs
//       WeightGram int
//       LengthCm, WidthCm, HeightCm int
//       IsDefault, IsActive bool
//       Position int
//       CreatedAt; UpdatedAt
//   }
//   func NewVariant(...) (*Variant, error)
//         WeightGram wajib > 0
//         ContentQuantity wajib > 0
//         kalau UnitCode == BaseUnitCode maka ContentQuantity harus 1
//
//   func (v *Variant) ToBaseQuantity(qty decimal) decimal   // qty * ContentQuantity
//   func (v *Variant) TotalWeightGram(qty decimal) int
//
// ToBaseQuantity adalah fungsi yang mengikat seluruh sistem stok. Pembeli
// memesan 2 dus; stok berkurang 80 pcs. Tanpa ini, stok tiga varian yang
// sebenarnya barang yang sama akan saling berbohong dan penjual akan menjual
// barang yang sudah habis. Lihat ADR-006.
//
// WeightGram wajib dan tidak boleh nol. Ongkir kurir dihitung dari berat total;
// varian tanpa berat menghasilkan ongkir nol dan kerugian yang baru ketahuan
// saat rekonsiliasi.
