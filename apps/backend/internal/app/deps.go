package app

// Dependency infrastruktur bersama yang dibangun di main.
//
// TODO:
//   type Deps struct {
//       Config config.Config
//       Logger *slog.Logger
//       Pool *pgxpool.Pool
//       Redis *redis.Client
//       Tx postgres.TxManager
//       Bus eventbus.Bus
//       Outbox outbox.Store
//       Tokens *auth.JWT
//       Hasher auth.Hasher
//       Storage storage.Storage
//       Mailer mailer.Mailer
//       Clock clock.Clock
//   }
//
// Hanya infrastruktur. Dependensi antar modul dirakit di registry.go.
