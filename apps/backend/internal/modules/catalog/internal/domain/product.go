package domain

// TODO:
//   type Status string   // draft | active | archived
//   type Product struct { ID; SellerID; CategoryID; Name; Slug; Description;
//                         Brand; Status; IsTaxable;
//                         RatingAvg float64; RatingCount, SoldCount int;
//                         PublishedAt; CreatedAt; UpdatedAt; DeletedAt }
//   func NewProduct(sellerID uuid.UUID, name string) (*Product, error)
//   func (p *Product) Publish(variantCount int) error
//         tolak kalau belum punya varian aktif, atau belum punya foto
//   func (p *Product) Archive()
//   func (p *Product) IsPurchasable() bool
//
// RatingAvg dan SoldCount adalah agregat yang diperbarui WORKER lewat event.
// JANGAN menghitungnya saat request — menghitung rata-rata dari jutaan ulasan
// di setiap pembukaan halaman produk adalah cara membunuh database.
//
// Produk dimiliki penjual. Dua penjual yang menjual barang sama punya dua baris
// berbeda. Katalog master bersama terdengar rapi tapi memunculkan pertanyaan
// siapa yang berhak mengubah nama dan fotonya; tunda sampai benar-benar perlu.
