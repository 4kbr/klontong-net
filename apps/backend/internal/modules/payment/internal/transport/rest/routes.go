package rest

// TODO:
//   Pembeli:
//     GET  /api/v1/orders/{orderID}/payment
//     POST /api/v1/orders/{orderID}/payment/retry
//
//   Webhook (TANPA sesi, TANPA CORS, verifikasi signature):
//     POST /webhook/payment/{provider}
//
//   Admin:
//     GET  /api/v1/admin/payments
//     POST /api/v1/admin/payments/{paymentID}/reconcile
//     POST /api/v1/admin/refunds
//
// Endpoint webhook tetap kena rate limit. Gateway pun bisa salah dan membanjiri.
