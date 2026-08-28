package payment

// TODO:
//   type Instruction struct {
//       PaymentID uuid.UUID
//       Method, Channel string
//       RedirectURL, VANumber, QRString string
//       ExpiresAt time.Time
//   }
//   type Port interface {
//       CreatePayment(ctx, orderID uuid.UUID, amount money.Amount,
//                     method, channel string, idempotencyKey string) (Instruction, error)
//       GetByOrder(ctx, orderID uuid.UUID) (Info, error)
//       RequestRefund(ctx, RefundInput) (uuid.UUID, error)
//   }
//
// CreatePayment dipanggil `order` saat checkout. Ia mengembalikan instruksi yang
// diteruskan ke pembeli — nomor virtual account, QR, atau URL redirect.
//
// Untuk COD, CreatePayment tetap dipanggil tapi tidak menghubungi gateway;
// ia hanya mencatat bahwa pembayaran akan terjadi saat barang diterima.
// Menyeragamkan jalurnya membuat `order` tidak perlu bercabang.
