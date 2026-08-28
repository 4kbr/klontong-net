package app

// Runner pekerjaan latar, dipakai cmd/worker.
//
// TODO:
//   type Job struct { Name string; Interval time.Duration; Run func(ctx) error }
//   type Runner struct { jobs []Job; log *slog.Logger; redis *redis.Client }
//   func (r *Runner) Register(j Job)
//   func (r *Runner) Start(ctx context.Context) error
//
// Setiap job:
//   - ticker sendiri
//   - panic di-recover supaya satu job rusak tidak menjatuhkan yang lain
//   - error dicatat dengan nama job
//   - berhenti saat ctx dibatalkan
//
// UNTUK JOB YANG MENYENTUH UANG (pencairan, pematangan dana), tambahkan kunci
// terdistribusi di Redis supaya hanya satu instance yang menjalankannya pada
// satu waktu. Job yang mengurangi stok atau memindahkan saldo, kalau berjalan
// bersamaan di dua instance, menghasilkan kerusakan yang mahal diperbaiki.
//
// Kunci Redis saja tidak cukup sebagai jaminan — pasangkan dengan idempotency
// key di database. Redis bisa hilang; unique index tidak.
