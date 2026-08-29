# Fase 01 — Platform

> Prasyarat: tidak ada. Guide: [`../GUIDES.md`](../GUIDES.md) §0, §1, §2 (tahap 1).
> ADR terkait: ADR-005 (uang), ADR-003 (urutan TTL), ADR-008 (idempotensi).

## Tujuan

Fondasi bersama: config, uang, database + transaksi, HTTP, auth, event/outbox, error,
clock, id. **Selesai kalau** `/healthz` balas 200 dan property test `money.Distribute`
lulus.

## Aturan khusus fase ini

- `os.Getenv` **hanya** boleh dipanggil di package `config`.
- `math/big` **dilarang** di luar `platform/money` (dijaga depguard).
- `money.Distribute` menjamin Σ hasil = total, sisa dialokasikan deterministik. Sisa
  **tidak boleh dibuang**. Wajib property test.
- `config.Load()` **fail-fast validate**: `JWT_SECRET` ≥ 32 byte; webhook secret ada
  bila provider diset; `CommissionDefaultBPS` ∈ 0..10000; `STOCK_RESERVATION_TTL <
  PAYMENT_EXPIRY` (ADR-003); di production `DB_SSLMODE != disable` + `PUBLIC_BASE_URL`
  https.

## Urutan kerja

### 1. Uang (kerjakan pertama, test dulu)

- [ ] `internal/platform/money/money.go` — `type Amount int64`; `Add/Sub/Mul(int64)/
      IsZero/IsNegative/String`; `MarshalJSON` emit **angka**; `type BasisPoints int`;
      `Amount.ApplyBPS(bps)` = `a*bps/10000` dengan rounding eksplisit; `Distribute(total
      Amount, weights []Amount) []Amount`.
- [ ] `internal/platform/money/rounding.go` — `RoundDown`, `RoundNearest`; dokumentasikan
      alasan tiap aturan (komisi & diskon → ke bawah; ongkir → angka provider apa adanya).
- [ ] `internal/platform/money/money_test.go` (**baru**) — property test `Distribute`:
      input acak, Σ hasil = total, tidak ada bagian negatif.

### 2. Config & error

- [ ] `internal/platform/config/config.go` — `Config` + sub-config (`App/HTTP/DB/Redis/
      JWT/Payment/Shipping/Marketplace/Stock/Storage/SMTP/Outbox/Log`); `Load()` dengan
      godotenv saat `APP_ENV=development` + validasi fail-fast; `DBConfig.DSN()`.
- [ ] `internal/platform/errs/errs.go` — sentinel (`ErrNotFound/Conflict/InvalidInput/
      Unauthorized/Forbidden/TooManyRequests/Upstream/Internal`); `type Error` dengan
      `Kind/Code/Message/Fields/Retryable/cause`, `Error()`, `Unwrap()` → `Kind`;
      konstruktor `NotFound/Conflict/Invalid/...`. Kode error spesifik & stabil
      (`out_of_stock`, `price_changed`, `voucher_expired`).

### 3. Postgres

- [ ] `internal/platform/postgres/pool.go` — `NewPool(ctx, cfg)` (MaxConns/MinConns/
      lifetime/healthcheck + Ping fail-fast); `Healthy(ctx, pool)`.
- [ ] `internal/platform/postgres/tx.go` — `Executor` interface (dipenuhi `*pgxpool.Pool`
      & `pgx.Tx`); `TxManager` (`WithinTx`, `WithinTxOpts`); `NewTxManager(pool)`; reuse
      tx dari ctx atau Begin + simpan via private key; rollback on error/panic + re-panic;
      `ExecutorFrom(ctx, pool)`.
- [ ] `internal/platform/postgres/pgerr.go` — `Translate(err)`: `ErrNoRows`→NotFound;
      `23505`→Conflict(+constraint); `23503`→Conflict; `23514`→Invalid;
      `40001`→Conflict+Retryable; `55P03`→Conflict+Retryable.

### 4. HTTP

- [ ] `internal/platform/httpx/router.go` — `NewRouter(deps)`: middleware global
      (RequestID, RealIP, Recoverer, Logger, Timeout, CORS); **4 area**
      (`/api/v1/*` publik-sebagian, `/api/v1/seller/*`, `/api/v1/admin/*`, `/webhook/*`);
      `/healthz`, `/readyz`. Jangan pasang satu auth middleware yang bercabang di dalam.
- [ ] `internal/platform/httpx/response.go` — envelope sukses `{data,meta}` / gagal
      `{error:{code,message,fields}}`; `JSON/OK/Created/Accepted/NoContent/Paginated`.
- [ ] `internal/platform/httpx/error.go` — `Error(w,r,err)`: satu-satunya tempat pilih
      status; 5xx log penuh + pesan generik; jangan teruskan teks error gateway ke buyer.
- [ ] `internal/platform/httpx/request.go` — `DecodeJSON[T]` (tolak non-JSON,
      MaxBytesReader, DisallowUnknownFields, error → `errs.Invalid`); `URLParamUUID`,
      `QueryInt/Bool/String`, `Pagination` (keyset cursor).
- [ ] `internal/platform/httpx/validate.go` — `Validatable`, `DecodeAndValidate[T]`.
- [ ] `internal/platform/httpx/idempotency.go` — `Idempotent(store, ttl)` middleware;
      `IdempotencyStore` (`Begin` → replay/acquired, `Finish`); Redis untuk cepat, unique
      index DB jaminan sesungguhnya.
- [ ] `internal/platform/middleware/*.go` — `requestid`, `logging` (skip healthz),
      `recover` (re-panic `http.ErrAbortHandler`), `cors` (hanya `/api/v1`, origin
      eksplisit), `auth` (`Authenticate`, `OptionalAuth`; bedakan `token_expired` vs
      `token_invalid`), `role` (`RequireRole`), `ratelimit` (per area).

### 5. Auth

- [ ] `internal/platform/auth/context.go` — `Identity{UserID,Roles}`, `Has`,
      `WithIdentity/IdentityFrom/MustUserID` (private key type).
- [ ] `internal/platform/auth/jwt.go` — `Claims`; `TokenIssuer`/`TokenVerifier`;
      `NewJWT(cfg)`; access token pendek (15m); `Verify` **cek signing method eksplisit**
      (blokir `alg: none`).
- [ ] `internal/platform/auth/password.go` — `Hasher` (`Hash`/`Compare`);
      `NewBcryptHasher(cost)` (12 prod, 4 test); `Compare` tetap jalan lawan dummy hash
      saat user tidak ada (constant-time).

### 6. Event & outbox

- [ ] `internal/platform/eventbus/event.go` — `Event{ID,Type,AggregateType,AggregateID,
      Payload,OccurredAt}`; `New(...)`, `Decode[T]`.
- [ ] `internal/platform/eventbus/bus.go` — `Handler`, `Publisher`, `Subscriber`, `Bus`.
- [ ] `internal/platform/eventbus/memory.go` — `MemoryBus` sinkron; error satu handler
      di-log tapi tidak menghentikan yang lain; `Subscribe` hanya dipanggil saat startup.
- [ ] `internal/platform/outbox/model.go` — `Record` ↔ `outbox_events` ↔ `eventbus.Event`.
- [ ] `internal/platform/outbox/store.go` — `Store.Save(ctx, ...events)` **wajib**
      `postgres.ExecutorFrom(ctx, pool)` supaya baris outbox ikut tx berjalan.
- [ ] `internal/platform/outbox/relay.go` — `Run(ctx)`: per tick `BEGIN; SELECT ... WHERE
      published_at IS NULL ORDER BY created_at LIMIT $batch FOR UPDATE SKIP LOCKED;
      publish; sukses→published_at, gagal→attempts+1,last_error; COMMIT`. Metrik pending.

### 7. Sisa platform

- [ ] `internal/platform/clock/clock.go` — `Clock` interface, `New()`, `Fixed{T}`.
- [ ] `internal/platform/id/id.go` — `New()` (uuid v7 bila ada), `Parse`;
      `OrderNumber(t, seq)` → `KN-20260827-000123`, `SuborderNumber(orderNumber, idx)`.
- [ ] `internal/platform/httpclient/client.go` — `New(name, timeout, log)` (redaksi
      `Authorization`, forward request id); `Do[T]` (≥400 → `errs.Upstream` +
      `Retryable` benar); `Retry(ctx, attempts, fn)` backoff+jitter, hormati `Retry-After`.
- [ ] `internal/platform/storage/storage.go` — `Storage` interface (`Put/PresignGet/
      PresignPut/Delete`), `NewS3(cfg)`. Dokumen verifikasi seller **tidak pernah publik**.
- [ ] `internal/platform/logger/logger.go` — `New(cfg)` (json prod, text dev),
      `FromContext/WithContext`. Jangan log kunci gateway / nomor bank / data kartu.
- [ ] `internal/platform/mailer/mailer.go` — `Mail`, `Mailer.Send`, `NewSMTP(cfg)`
      (dev → MailHog :1025). Email tidak pernah dikirim di jalur request.
- [ ] `internal/platform/redis/redis.go` — `New(cfg)` + Ping.
- [ ] `internal/platform/server/server.go` — `New(cfg, h, log)`, `Run(ctx)` graceful
      shutdown.
- [ ] `internal/platform/pagination/cursor.go` — `Cursor{SortValue,ID}`, base64-JSON
      encode/decode, `Page[T]{Items,NextCursor,HasMore}` (fetch n+1).

### 8. Entry point & wiring dasar

- [ ] `cmd/api/main.go` — urutan: `config.Load` → `logger.New` → ctx SIGINT/SIGTERM →
      `postgres.NewPool`+Ping → `redis.New` → build `app.Deps` → `app.NewRegistry(deps)`
      → `httpx.NewRouter` + `registry.MountRoutes` → `server.New(...).Run(ctx)`.
      Worker **tidak** dijalankan di sini (ADR-013).
- [ ] `cmd/migrate/main.go` — runner migrasi programatik; flag `-direction`, `-steps`,
      `-force`; baca `DATABASE_URL` via `config.Load`; exit non-zero saat gagal.
- [ ] `internal/app/deps.go` — `type Deps struct { Config, Logger, Pool, Redis, Tx, Bus,
      Outbox, Tokens, Hasher, Storage, Mailer, Clock }` (infra saja).
- [ ] `internal/app/registry.go` — kerangka `NewRegistry(deps)`, `MountRoutes(mux)`,
      `RegisterWorkers(runner)`, `Close()` (modul diisi mulai fase 2).
- [ ] `internal/app/worker.go` — `Job{Name,Interval,Run}`, `Runner` (`Register`,
      `Start(ctx)`): tiap job ticker sendiri, panic di-recover, error di-log dengan nama
      job, stop saat ctx batal.
- [ ] `migrations/000015_create_audit_and_outbox.up.sql` / `.down.sql` — **bagian
      `outbox_events` saja** (audit_events menyusul di fase 13): kolom sesuai spec, index
      `WHERE published_at IS NULL`.
- [ ] `Makefile` — tambahkan target tipis `setup` / `up` / `migrate` (alias
      `migrate-up`) / `tunnel` yang dirujuk `GUIDES.md`, atau perbaiki rujukannya.

## Test wajib

- Property test `money.Distribute` (Σ = total, tidak ada bagian negatif) — [`GUIDES.md`
  §15](../GUIDES.md).
- `config.Load` menolak `STOCK_RESERVATION_TTL >= PAYMENT_EXPIRY`.
- `auth.jwt.Verify` menolak token `alg: none`.
- `postgres.Translate` memetakan `23505` → `errs.ErrConflict`.

## Sengaja TIDAK dikerjakan di fase ini

- Implementasi modul bisnis apa pun (`internal/modules/*`).
- `cmd/worker` (rangka runner boleh, daftar job di fase 13).
- `gateway_midtrans.go` / disburser nyata (fase 11 / 14).
- Bagian `audit_events` di migrasi 000015 (fase 13).
