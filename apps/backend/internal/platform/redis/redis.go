package redis

// TODO:
//   func New(cfg config.RedisConfig) (*redis.Client, error)   // dengan Ping
//
// Dipakai untuk: cache katalog, state rate limit, kunci idempotensi, kunci
// terdistribusi untuk worker.
//
// ATURAN: tidak ada data yang HANYA ada di Redis. Kalau Redis dikosongkan,
// sistem harus tetap benar — cuma lebih lambat sesaat. Perlindungan idempotensi
// yang sesungguhnya ada di unique index database.
