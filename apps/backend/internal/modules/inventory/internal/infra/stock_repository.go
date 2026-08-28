package infra

// TODO: implementasi StockRepository.
//
// GetForUpdate:
//   SELECT ... FROM inventory_stocks
//   WHERE (outlet_id, variant_id) IN (...)
//   ORDER BY outlet_id, variant_id        <- urutan konsisten, cegah deadlock
//   FOR UPDATE
//
// Urutan di ORDER BY itu bukan kosmetik. Dua transaksi yang mengunci baris yang
// sama dalam urutan berbeda akan saling menunggu selamanya.
