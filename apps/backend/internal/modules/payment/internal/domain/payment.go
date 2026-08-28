package domain

// TODO:
//   type Method string   // gateway | cod
//   type Status string   // pending | authorized | settled | failed | expired |
//                        // refunded | partially_refunded
//   type Payment struct { ID; OrderID; Method; Provider; ProviderReference;
//                         Amount money.Amount; Status; Channel;
//                         PaidAt, ExpiredAt *time.Time; FailedReason string
//                         IdempotencyKey string
//                         CreatedAt; UpdatedAt }
//   func (p *Payment) MarkSettled(at time.Time, providerRef string) error
//   func (p *Payment) MarkFailed(reason string) error
//   func (p *Payment) MarkExpired() error
//   func (p *Payment) CanRefund() bool
//
// MarkSettled harus IDEMPOTEN: dipanggil dua kali dengan referensi yang sama
// tidak melakukan apa-apa dan tidak error. Webhook gateway datang berkali-kali,
// dan sering tidak berurutan.
//
// Satu pesanan hanya boleh punya satu pembayaran aktif. Ditegakkan unique index
// parsial di database — tanpa itu, klik ganda pada tombol bayar menghasilkan
// dua tagihan.
//
// COD: status langsung `pending` dan baru `settled` saat kurir menyetorkan uang.
// Arah arus kasnya berbeda total dari gateway. Lihat ADR-012.
