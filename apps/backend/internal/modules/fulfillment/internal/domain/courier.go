package domain

// Abstraksi agregator kurir.
//
// TODO:
//   type RateRequest struct { OriginPostalCode, DestPostalCode string
//                             OriginLat, OriginLng, DestLat, DestLng float64
//                             WeightGram int; ItemValue money.Amount }
//   type Rate struct { CourierCode, ServiceName string; Amount money.Amount
//                      ETDMinDay, ETDMaxDay int }
//   type CourierProvider interface {
//       Rates(ctx, RateRequest) ([]Rate, error)
//       CreateBooking(ctx, BookingRequest) (trackingNumber string, err error)
//       Track(ctx, trackingNumber string) ([]TrackingEvent, error)
//   }
//
// Rates dipanggil saat checkout dan HARUS cepat. Beri timeout ketat (2–3 detik)
// dan tangani kegagalannya dengan anggun: kalau agregator sedang bermasalah,
// tampilkan opsi lain (pickup, antar lokal) alih-alih menggagalkan seluruh
// checkout. Kehilangan satu opsi lebih baik daripada kehilangan pesanan.
//
// Cache hasil Rates di Redis dengan kunci (asal, tujuan, berat dibulatkan).
// Tarif tidak berubah tiap menit, dan ini query yang paling sering dipanggil.
