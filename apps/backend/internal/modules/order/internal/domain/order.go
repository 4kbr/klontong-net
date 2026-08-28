package domain

// TODO:
//   type Status string
//   const (
//       StatusPendingPayment      Status = "pending_payment"
//       StatusPaid                Status = "paid"
//       StatusPartiallyFulfilled  Status = "partially_fulfilled"
//       StatusCompleted           Status = "completed"
//       StatusCancelled           Status = "cancelled"
//       StatusExpired             Status = "expired"
//   )
//
//   type Order struct {
//       ID; Number; BuyerUserID; Status
//       ItemsAmount, DiscountAmount, ShippingAmount, TaxAmount, GrandTotal money.Amount
//       PaymentMethod string
//       ShippingAddress AddressSnapshot     // SALINAN, bukan referensi
//       BuyerNote string
//       Suborders []*Suborder
//       PlacedAt, PaidAt, CompletedAt, CancelledAt, ExpiresAt *time.Time
//       CancelReason string
//   }
//
//   func (o *Order) RecalculateTotals() error
//         GrandTotal = ItemsAmount - DiscountAmount + ShippingAmount + TaxAmount
//         dan harus SAMA dengan jumlah total seluruh suborder. Periksa itu
//         secara eksplisit dan kembalikan error kalau tidak cocok — selisih
//         satu rupiah di sini berarti pembukuan tidak akan pernah balance.
//
//   func (o *Order) MarkPaid(now time.Time) error
//   func (o *Order) Cancel(reason string, now time.Time) error
//   func (o *Order) SyncStatusFromSuborders() Status
//         Status order INDUK adalah TURUNAN dari status suborder:
//           semua suborder selesai        -> completed
//           semua dibatalkan/ditolak      -> cancelled
//           sebagian selesai              -> partially_fulfilled
//         Jangan pernah menyetel status order langsung setelah pembayaran —
//         hitung dari anak-anaknya. Lihat ADR-002.
//
//   func (o *Order) CanBeCancelledByBuyer() bool
//         Pembeli boleh membatalkan selama BELUM ADA suborder yang dikirim.
//         Setelah ada yang dikirim, pembatalan jadi urusan per suborder.
