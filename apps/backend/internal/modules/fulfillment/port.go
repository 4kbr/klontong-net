package fulfillment

// TODO:
//   type QuoteRequest struct {
//       OutletID uuid.UUID
//       DestLat, DestLng *float64
//       DestCity, DestPostalCode string
//       TotalWeightGram int
//       TotalValue money.Amount
//   }
//   type Option struct {
//       Method string              // local_delivery | courier | pickup
//       CourierCode, ServiceName string
//       Amount money.Amount
//       ETDMinDay, ETDMaxDay int
//       DistanceKm float64
//       QuoteID uuid.UUID
//       ExpiresAt time.Time
//   }
//   type Port interface {
//       Quote(ctx, req QuoteRequest) ([]Option, error)
//       ValidateQuote(ctx, quoteID uuid.UUID, at time.Time) (Option, error)
//       CreateShipment(ctx, CreateShipmentInput) (uuid.UUID, error)
//   }
//
// Quote dipanggil PER SUBORDER, karena tiap penjual mengirim dari outletnya
// sendiri. Satu pesanan berisi tiga penjual menghasilkan tiga daftar opsi
// dengan harga yang berbeda-beda.
//
// ValidateQuote dipanggil saat checkout untuk memastikan pilihan pembeli masih
// berlaku. Tarif kurir berubah; quote kedaluwarsa harus dihitung ulang, bukan
// dipakai begitu saja. Kalau tidak, ada selisih ongkir yang harus ditanggung
// seseorang.
//
// Opsi `pickup` hanya muncul kalau outlet mendukungnya. Opsi `local_delivery`
// hanya muncul kalau outlet mendukung DAN alamat tujuan punya koordinat DAN
// jaraknya dalam radius layanan.
