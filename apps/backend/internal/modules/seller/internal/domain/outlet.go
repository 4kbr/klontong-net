package domain

// TODO:
//   type Outlet struct { ID; SellerID; Name; Phone;
//                        Province, City, District, PostalCode, AddressLine string;
//                        Latitude, Longitude float64;
//                        IsActive bool;
//                        SupportsPickup, SupportsLocalDelivery, SupportsCourier bool;
//                        OperatingHours OperatingHours;
//                        CreatedAt; UpdatedAt }
//   func NewOutlet(sellerID uuid.UUID, name string, lat, lng float64) (*Outlet, error)
//         koordinat WAJIB dan harus masuk akal (lat -90..90, lng -180..180)
//   func (o *Outlet) SupportsMethod(method string) bool
//   func (o *Outlet) IsOpenAt(t time.Time) bool
//
// type OperatingHours: jam buka per hari. Dipakai untuk memutuskan apakah
// ambil di toko dan antar lokal boleh ditawarkan sekarang. Outlet yang tutup
// tetap boleh menerima pesanan kurir — barangnya dikirim besok.
//
// Setiap penjual harus punya minimal satu outlet aktif untuk bisa berjualan.
// Tanpa outlet, tidak ada tempat stok disimpan dan tidak ada titik asal ongkir.
