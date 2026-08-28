package domain

// Tarif antar lokal.
//
// TODO:
//   type LocalTariff struct { BaseFare, PerKmFare money.Amount
//                             MaxDistanceKm float64
//                             MinOrderAmount money.Amount }
//   func (t LocalTariff) Calculate(distanceKm float64) (money.Amount, error)
//         di luar MaxDistanceKm -> error, opsi ini tidak ditawarkan
//         pembulatan jarak ke atas per km, tetapkan aturannya dan konsisten
//
//   func Haversine(lat1, lng1, lat2, lng2 float64) float64
//         jarak garis lurus dalam km
//
// PERINGATAN tentang Haversine: ia menghitung jarak GARIS LURUS, bukan jarak
// tempuh. Untuk kota padat, jarak sebenarnya bisa 1,3–1,6 kali lipat. Pilih
// salah satu dan sadari konsekuensinya:
//   - pakai Haversine dengan faktor pengali, sederhana tapi kasar
//   - panggil API routing, akurat tapi berbayar dan menambah latensi checkout
// Rekomendasi: mulai dengan Haversine + faktor, catat sebagai ADR, dan ganti
// kalau keluhan ongkir mulai muncul.
//
// MinOrderAmount ada karena mengantar pesanan Rp15.000 sejauh 8 km itu rugi.
