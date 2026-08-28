package domain

// Unit kerja sesungguhnya di sistem ini. Baca ADR-002.
//
// TODO:
//   type SuborderStatus string
//   const (
//       SubAwaitingConfirmation SuborderStatus = "awaiting_confirmation"
//       SubConfirmed            SuborderStatus = "confirmed"
//       SubPacking              SuborderStatus = "packing"
//       SubReadyForPickup       SuborderStatus = "ready_for_pickup"
//       SubShipped              SuborderStatus = "shipped"
//       SubDelivered            SuborderStatus = "delivered"
//       SubCompleted            SuborderStatus = "completed"
//       SubCancelled            SuborderStatus = "cancelled"
//       SubRejected             SuborderStatus = "rejected"
//   )
//
//   type Suborder struct {
//       ID; OrderID; Number; SellerID; OutletID; Status
//       FulfillmentMethod string        // local_delivery | courier | pickup
//       Items []*Item
//       ItemsAmount, DiscountAmount, ShippingAmount, TaxAmount, TotalAmount money.Amount
//       CommissionBPS int               // DIBEKUKAN saat pesanan dibuat
//       CommissionAmount, SellerEarningAmount money.Amount
//       ConfirmedAt, ShippedAt, DeliveredAt, CompletedAt, CancelledAt *time.Time
//       CancelReason, RejectedReason string
//   }
//
//   func (s *Suborder) Transition(to SuborderStatus, now time.Time) error
//         Transisi HANYA lewat method ini. Peta transisi yang sah didefinisikan
//         di state.go. Menyetel Status langsung dari usecase adalah cara paling
//         cepat menciptakan pesanan yang statusnya mustahil.
//
//   func (s *Suborder) ComputeCommission() error
//         CommissionAmount = ItemsAmount setelah diskon, dikali CommissionBPS,
//         dibulatkan KE BAWAH.
//         SellerEarningAmount = (ItemsAmount - DiscountAmount) - CommissionAmount
//                               + ShippingAmount (kalau penjual yang menanggung kirim)
//
//         KEPUTUSAN YANG HARUS DIAMBIL: apakah komisi dihitung dari nilai barang
//         saja atau termasuk ongkir. Rekomendasi: barang saja — mengambil komisi
//         dari ongkir yang diteruskan ke kurir berarti penjual rugi. Catat
//         pilihannya sebagai ADR.
//
//   func (s *Suborder) RequiresShipment() bool    // pickup tidak butuh pengiriman
//
// CommissionBPS DIBEKUKAN di sini, bukan dibaca dari seller saat dibutuhkan.
// Marketplace yang mengubah komisi besok tidak boleh mengubah bagi hasil
// pesanan kemarin.
