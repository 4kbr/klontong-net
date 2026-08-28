package domain

// TODO: ErrPaymentNotFound, ErrPaymentAlreadySettled, ErrPaymentExpired,
//       ErrInvalidSignature, ErrAmountMismatch, ErrRefundExceedsPayment,
//       ErrUnsupportedChannel, ErrGatewayUnavailable
//
// ErrAmountMismatch: gateway melaporkan jumlah yang berbeda dari tagihan kita.
// JANGAN pernah menerima ini diam-diam. Catat sebagai insiden, tandai pembayaran
// untuk ditinjau manusia, dan jangan tandai pesanan lunas.
