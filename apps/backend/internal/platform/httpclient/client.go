package httpclient

// Klien HTTP bersama untuk memanggil layanan luar.
//
// TODO:
//   func New(name string, timeout time.Duration, log *slog.Logger) *http.Client
//     - Transport dengan batas koneksi yang masuk akal
//     - RoundTripper pembungkus: catat durasi & status, teruskan request id,
//       dan REDAKSI header Authorization sebelum log
//
//   func Do[T any](ctx, client, req) (T, error)
//     status >= 400 -> errs.Upstream dengan Retryable diisi benar:
//       429 dan 5xx -> retryable; 4xx lain -> tidak
//
//   func Retry(ctx, attempts int, fn func() error) error
//     backoff eksponensial dengan jitter, hormati Retry-After
//
// Selalu http.NewRequestWithContext. Request tanpa context tidak bisa
// dibatalkan dan menahan goroutine saat shutdown.
//
// PERINGATAN untuk pemanggilan yang menyentuh uang: retry pada request yang
// MENCIPTAKAN sesuatu (buat transaksi, kirim transfer) hanya aman kalau
// requestnya membawa idempotency key. Tanpa itu, timeout lalu retry bisa
// berarti dua transaksi. Bedakan dengan jelas mana yang aman diretry.
