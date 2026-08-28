package logger

// Wrapper di atas log/slog.
//
// TODO:
//   func New(cfg config.LogConfig) *slog.Logger    // json di produksi, text di dev
//   func FromContext(ctx) *slog.Logger
//   func WithContext(ctx, *slog.Logger) context.Context
//
// Log yang berkaitan dengan request membawa request_id dan user_id.
// Untuk alur pesanan, sertakan juga order_id dan suborder_id — menelusuri satu
// pesanan bermasalah lintas modul tanpa itu hampir mustahil.
//
// JANGAN pernah menulis ke log: server key gateway, nomor rekening lengkap,
// isi payload webhook mentah yang memuat data kartu.
