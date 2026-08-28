package rest

// TODO:
//   Pembeli:
//     POST /api/v1/checkout/preview
//     POST /api/v1/checkout          (Idempotency-Key WAJIB)
//     GET  /api/v1/orders
//     GET  /api/v1/orders/{orderID}
//     POST /api/v1/orders/{orderID}/cancel
//     POST /api/v1/suborders/{suborderID}/confirm-received
//
//   Dasbor penjual:
//     GET  /api/v1/seller/orders
//     GET  /api/v1/seller/orders/{suborderID}
//     POST /api/v1/seller/orders/{suborderID}/confirm
//     POST /api/v1/seller/orders/{suborderID}/reject
//     POST /api/v1/seller/orders/{suborderID}/pack
//     POST /api/v1/seller/orders/{suborderID}/ship
//     POST /api/v1/seller/orders/{suborderID}/ready-for-pickup
//
//   Admin:
//     GET  /api/v1/admin/orders
//     POST /api/v1/admin/orders/{orderID}/cancel
//
// Perhatikan bahwa penjual mengakses SUBORDER, bukan order. Penjual tidak
// pernah melihat pesanan penjual lain di keranjang yang sama, dan URL-nya
// mencerminkan itu.
