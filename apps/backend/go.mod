module github.com/yourorg/klontong-net/apps/backend

go 1.23

// TODO: `go get` dependency berikut saat mulai implementasi.
//
//   github.com/go-chi/chi/v5                     router
//   github.com/go-chi/cors                       cors
//   github.com/jackc/pgx/v5                      driver + pool postgres
//   github.com/golang-migrate/migrate/v4          migrasi
//   github.com/google/uuid                        id
//   github.com/golang-jwt/jwt/v5                  access token
//   golang.org/x/crypto                           bcrypt
//   github.com/redis/go-redis/v9                  cache, rate limit, idempotensi
//   github.com/joho/godotenv                      .env saat dev
//   github.com/shopspring/decimal                 HANYA untuk perhitungan pajak/komisi
//                                                 berbasis persen; uang tetap disimpan
//                                                 sebagai int64. Lihat ADR-005.
//   github.com/stretchr/testify                   assertion
//   github.com/testcontainers/testcontainers-go   integration test
//
// Payment gateway dan agregator ongkir dipanggil lewat REST langsung dengan
// net/http, dibungkus di modulnya masing-masing. SDK mereka biasanya tipis dan
// mengikat kita ke bentuk mereka. Lihat ADR-014.
