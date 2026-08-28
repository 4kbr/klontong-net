package domain

// TODO:
//   type Item struct { ID; CartID; SellerID; OutletID; VariantID;
//                      Quantity decimal
//                      UnitPriceAmount money.Amount   // CACHE, bukan kebenaran
//                      PriceCheckedAt time.Time
//                      Note string
//                      CreatedAt; UpdatedAt }
//   func (i *Item) SetQuantity(q decimal, unitAllowsFraction bool, max decimal) error
//         tolak pecahan untuk satuan yang tidak mengizinkan (0,5 dus tidak masuk akal)
//         tolak melebihi batas maksimum per baris
//         qty <= 0 berarti hapus baris, bukan error
//   func (i *Item) IsPriceStale(now time.Time, ttl time.Duration) bool
//
// UnitPriceAmount adalah CACHE untuk tampilan. Harga dihitung ulang setiap kali
// keranjang dibuka dan sekali lagi saat checkout. Kalau berubah, beri tahu
// pembeli dengan jelas — jangan diam-diam memakai angka lama, dan jangan juga
// diam-diam memakai angka baru. Lihat ADR-004.
//
// SellerID disimpan di item meski bisa didapat dari varian. Pengelompokan per
// penjual terjadi di setiap tampilan; tidak layak membayar join setiap kali.
