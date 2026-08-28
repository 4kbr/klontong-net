package domain

// TODO:
//   type OrderRepository interface {
//       Create(ctx, *Order) error            // menyimpan order + suborder + item sekaligus
//       Update(ctx, *Order) error
//       FindByID(ctx, id) (*Order, error)
//       FindByNumber(ctx, number string) (*Order, error)
//       ListByBuyer(ctx, userID uuid.UUID, filter, cursor, limit) ([]*Order, string, error)
//       ListExpired(ctx, before time.Time, limit int) ([]*Order, error)
//       NextOrderNumber(ctx, t time.Time) (string, error)
//   }
//   type SuborderRepository interface {
//       Update(ctx, *Suborder) error
//       FindByID(ctx, id) (*Suborder, error)
//       ListByOrder(ctx, orderID uuid.UUID) ([]*Suborder, error)
//       ListBySeller(ctx, sellerID uuid.UUID, filter, cursor, limit) ([]*Suborder, string, error)
//       ListDeliveredBefore(ctx, before time.Time, limit int) ([]*Suborder, error)
//   }
//   type ItemRepository interface { ListByOrder/ListBySuborder/FindPurchased }
//   type StatusEventRepository interface { Create/ListByOrder }
//
// NextOrderNumber harus aman terhadap konkurensi. Pakai sequence Postgres, bukan
// SELECT MAX + 1 — yang kedua akan menghasilkan nomor ganda saat ramai.
//
// ListDeliveredBefore dipakai worker untuk menyelesaikan pesanan otomatis
// setelah masa komplain lewat.
