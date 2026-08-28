package domain

// TODO:
//   type Stock struct { ID; OutletID; VariantID;
//                       QuantityOnHand decimal    // fisik ada di outlet
//                       QuantityReserved decimal  // sudah di-checkout, belum diambil
//                       LowStockThreshold decimal
//                       Version int
//                       UpdatedAt }
//   func (s *Stock) Available() decimal            // OnHand - Reserved
//   func (s *Stock) CanReserve(qty decimal) bool
//   func (s *Stock) Reserve(qty decimal) error     // Reserved += qty
//   func (s *Stock) Commit(qty decimal) error      // OnHand -= qty, Reserved -= qty
//   func (s *Stock) ReleaseReservation(qty decimal) error   // Reserved -= qty
//   func (s *Stock) Adjust(delta decimal) error
//   func (s *Stock) IsLow() bool
//
// Dua angka terpisah, bukan satu "tersedia". Dengan begitu kamu bisa menjawab
// "barangnya ada di rak atau sudah dipesan orang" — pertanyaan yang pasti
// muncul dari penjual.
//
// Commit mengurangi KEDUANYA: barang keluar dari rak dan reservasinya selesai.
// Salah satu saja terlewat berarti stok berbohong selamanya.
//
// Semua kuantitas dalam SATUAN DASAR. Lihat ADR-006.
