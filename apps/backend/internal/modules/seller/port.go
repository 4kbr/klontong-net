package seller

// TODO:
//   type Status string   // pending | verified | suspended | closed
//   type Info struct { ID uuid.UUID; Name, Slug string; Status Status;
//                      CommissionBPS int; LogoURL string }
//   type Outlet struct {
//       ID, SellerID uuid.UUID
//       Name, City, AddressLine string
//       Latitude, Longitude float64
//       IsActive bool
//       SupportsPickup, SupportsLocalDelivery, SupportsCourier bool
//   }
//
//   type Port interface {
//       Get(ctx, sellerID uuid.UUID) (Info, error)
//       GetMany(ctx, ids []uuid.UUID) (map[uuid.UUID]Info, error)
//       CanSell(ctx, sellerID uuid.UUID) (bool, error)
//       CommissionBPS(ctx, sellerID uuid.UUID) (int, error)
//       GetOutlet(ctx, outletID uuid.UUID) (Outlet, error)
//       ListOutlets(ctx, sellerID uuid.UUID) ([]Outlet, error)
//       IsMember(ctx, sellerID, userID uuid.UUID) (bool, error)
//   }
//
// CanSell dipanggil sebelum barang penjual boleh masuk keranjang dan sebelum
// checkout. Penjual yang di-suspend di tengah proses belanja harus ketahuan di
// checkout, bukan setelah pembeli membayar.
//
// CommissionBPS mengembalikan komisi efektif: milik penjual kalau ada, kalau
// tidak default marketplace. Pemanggil tidak perlu tahu aturan itu.
