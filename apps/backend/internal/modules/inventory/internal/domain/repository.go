package domain

// TODO:
//   type StockRepository interface {
//       Get(ctx, outletID, variantID) (*Stock, error)
//       GetForUpdate(ctx, pairs []OutletVariant) ([]*Stock, error)   // SELECT ... FOR UPDATE
//       GetMany(ctx, pairs []OutletVariant) (map[OutletVariant]*Stock, error)
//       Upsert(ctx, *Stock) error
//       OutletsWithStock(ctx, variantID uuid.UUID, minQty decimal) ([]uuid.UUID, error)
//       ListByOutlet(ctx, outletID uuid.UUID, cursor string, limit int) ([]*Stock, string, error)
//   }
//   type ReservationRepository interface {
//       CreateMany/ListByOrder/ListExpired(ctx, before time.Time, limit int)/UpdateStatus
//   }
//   type MovementRepository interface { Create/CreateMany/ListByVariant/SumByVariant }
//   type OpnameRepository interface { Create/Get/AddItems/Finish/ListByOutlet }
//
// GetForUpdate mengunci baris stok dalam urutan yang KONSISTEN (urutkan
// berdasarkan outlet_id, variant_id sebelum mengunci). Mengunci dalam urutan
// berbeda antar transaksi adalah cara paling umum menciptakan deadlock, dan
// checkout dua pembeli yang membeli barang sama persis akan menemukannya.
