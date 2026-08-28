package config

// Konfigurasi. Dibaca SEKALI di main lalu diteruskan lewat parameter.
//
// TODO:
//   type Config struct { App; HTTP; DB; Redis; JWT; Payment; Shipping;
//                        Marketplace; Stock; Storage; SMTP; Outbox; Log }
//
//   type PaymentConfig struct {
//       Provider, ServerKey, ClientKey, WebhookSecret string
//       IsProduction bool
//       Expiry time.Duration
//   }
//   type MarketplaceConfig struct {
//       CommissionDefaultBPS int          // basis poin, 250 = 2,5%
//       SettlementHoldPeriod time.Duration
//       PayoutMinimum int64               // rupiah
//   }
//   type StockConfig struct { ReservationTTL time.Duration }
//
//   func Load() (Config, error)
//     - godotenv saat APP_ENV=development
//     - VALIDASI di akhir dan GAGALKAN start bila:
//         * JWT_SECRET < 32 byte
//         * PAYMENT_WEBHOOK_SECRET kosong sementara PAYMENT_PROVIDER terisi
//         * CommissionDefaultBPS di luar 0..10000
//         * StockConfig.ReservationTTL >= PaymentConfig.Expiry
//           (reservasi harus lepas SEBELUM pesanan kedaluwarsa, bukan sesudah —
//            kalau terbalik, ada jendela waktu stok tertahan untuk pesanan yang
//            sudah mati. Lihat ADR-003.)
//         * di produksi: DB_SSLMODE bukan disable dan PUBLIC_BASE_URL https
//
//   func (d DBConfig) DSN() string
//
// ATURAN: os.Getenv hanya dipanggil di package ini.
