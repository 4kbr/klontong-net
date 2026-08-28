package order

// Event modul ini adalah yang paling banyak consumernya di seluruh sistem.
//
// TODO:
//   const (
//       EventOrderPlaced         = "order.order.placed"
//       EventOrderPaid           = "order.order.paid"
//       EventOrderCancelled      = "order.order.cancelled"
//       EventOrderExpired        = "order.order.expired"
//       EventOrderCompleted      = "order.order.completed"
//
//       EventSuborderConfirmed   = "order.suborder.confirmed"
//       EventSuborderRejected    = "order.suborder.rejected"
//       EventSuborderShipped     = "order.suborder.shipped"
//       EventSuborderDelivered   = "order.suborder.delivered"
//       EventSuborderCompleted   = "order.suborder.completed"
//       EventSuborderCancelled   = "order.suborder.cancelled"
//   )
//
// Perhatikan bahwa peristiwa penting terjadi di level SUBORDER, bukan order.
// Consumer yang salah mendengarkan level akan bereaksi terlalu dini atau tidak
// bereaksi sama sekali:
//   - inventory commit stok saat SuborderShipped, bukan OrderPaid
//   - settlement mematangkan dana saat SuborderCompleted
//   - notification mengabari penjual saat OrderPaid, dan mengabari pembeli
//     saat tiap suborder berubah
//
// Payload WAJIB memuat SellerID dan OutletID untuk event suborder. Tanpa itu
// setiap consumer harus query balik, dan kopling yang mau dihindari kembali
// lewat pintu belakang.
