package postgres

// TODO:
//   func NewPool(ctx, cfg config.DBConfig) (*pgxpool.Pool, error)
//     ParseConfig -> MaxConns/MinConns/MaxConnLifetime/HealthCheckPeriod ->
//     NewWithConfig -> Ping (fail fast)
//   func Healthy(ctx, pool) error
//
// Satu pool untuk seluruh aplikasi.
