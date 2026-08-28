package domain

// Abstraksi payment gateway. Didefinisikan di domain supaya usecase tidak
// bergantung ke bentuk Midtrans atau Xendit.
//
// TODO:
//   type ChargeRequest struct {
//       OrderNumber string; Amount money.Amount; Channel string
//       CustomerName, CustomerEmail, CustomerPhone string
//       Items []GatewayItem
//       ExpiresAt time.Time
//       IdempotencyKey string
//   }
//   type ChargeResult struct {
//       ProviderReference string
//       RedirectURL, VANumber, QRString string
//       Status string
//   }
//   type Gateway interface {
//       Charge(ctx, ChargeRequest) (ChargeResult, error)
//       Inquiry(ctx, providerReference string) (Status, error)
//       Refund(ctx, providerReference string, amount money.Amount, reason string) (string, error)
//       VerifyWebhook(rawBody []byte, headers map[string]string) (WebhookPayload, error)
//   }
//
// Inquiry ada karena webhook BISA HILANG. Jangan pernah mengandalkan webhook
// sebagai satu-satunya jalur — worker rekonsiliasi menanyakan status untuk
// pembayaran yang menggantung. Lihat ADR-008.
//
// VerifyWebhook memvalidasi signature. Webhook yang tidak diverifikasi berarti
// siapa pun yang menemukan URL-nya bisa menyatakan pesanan lunas. Ini bukan
// kemungkinan teoretis — endpoint webhook mudah ditebak.
