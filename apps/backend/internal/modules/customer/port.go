package customer

// TODO:
//   type Address struct {
//       ID uuid.UUID
//       RecipientName, RecipientPhone string
//       Province, City, District, Village, PostalCode, AddressLine, Notes string
//       Latitude, Longitude *float64
//   }
//   type Port interface {
//       GetAddress(ctx, userID, addressID uuid.UUID) (Address, error)
//       DefaultAddress(ctx, userID uuid.UUID) (Address, error)
//   }
//
// Dipakai `order` saat checkout. Perhatikan bahwa GetAddress menerima userID
// juga — alamat hanya boleh diambil oleh pemiliknya, dan pemeriksaan itu ada di
// sini, bukan diserahkan ke pemanggil.
//
// Koordinat bertipe pointer karena bisa kosong. Alamat tanpa koordinat TIDAK
// BOLEH ditawari antar lokal; tarifnya dihitung per km dan tanpa titik tujuan
// perhitungannya mustahil.
