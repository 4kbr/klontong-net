package payment

// TODO:
//   EventPaymentCreated, EventPaymentSettled, EventPaymentFailed,
//   EventPaymentExpired, EventRefundCompleted
//
// EventPaymentSettled adalah pemicu terpenting di sistem: dari sinilah pesanan
// menjadi `paid` dan penjual mulai bekerja. Payload wajib memuat OrderID dan
// jumlah yang benar-benar diterima.
