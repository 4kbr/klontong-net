package postgres

// TxManager. Komponen paling penting di platform.
//
// Di sistem ini transaksi bukan sekadar kerapian — checkout menulis ke pesanan,
// reservasi stok, penukaran voucher, dan outbox sekaligus. Kalau salah satu
// gagal dan yang lain terlanjur, kamu punya pesanan tanpa stok atau voucher
// terpakai tanpa pesanan.
//
// TODO:
//   type Executor interface {   // dipenuhi *pgxpool.Pool dan pgx.Tx
//       Exec(ctx, sql string, args ...any) (pgconn.CommandTag, error)
//       Query(ctx, sql string, args ...any) (pgx.Rows, error)
//       QueryRow(ctx, sql string, args ...any) pgx.Row
//   }
//
//   type TxManager interface {
//       WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
//       WithinTxOpts(ctx context.Context, opts pgx.TxOptions, fn func(ctx) error) error
//   }
//   func NewTxManager(pool *pgxpool.Pool) TxManager
//
// Cara kerja:
//   1. ctx sudah punya tx -> jalankan fn dengan itu, jangan buka baru
//   2. belum -> Begin, simpan tx ke ctx lewat key privat
//   3. fn error -> Rollback; sukses -> Commit; panic -> Rollback lalu re-panic
//
//   func ExecutorFrom(ctx, pool) Executor
//     Ambil tx dari ctx bila ada, kalau tidak kembalikan pool.
//     SETIAP method repository dimulai dengan baris ini.
//
// WithinTxOpts ada karena beberapa alur butuh isolasi lebih ketat.
// Pengurangan stok dan penukaran voucher berkuota rawan write skew di level
// READ COMMITTED. Pilihannya: SERIALIZABLE dengan penanganan retry pada error
// 40001, atau kunci baris eksplisit dengan SELECT ... FOR UPDATE.
// Rekomendasi: FOR UPDATE untuk stok (barisnya sedikit dan jelas), dan
// biarkan unique index yang menegakkan kuota voucher.
//
// Transaksi dibuka HANYA di layer app.
