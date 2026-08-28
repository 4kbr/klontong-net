package outbox

// TODO:
//   type Store interface { Save(ctx context.Context, events ...eventbus.Event) error }
//   func NewStore(pool *pgxpool.Pool) Store
//
// Save WAJIB memakai postgres.ExecutorFrom(ctx, pool) supaya baris outbox ikut
// transaksi yang sedang berjalan. Itu seluruh inti pola ini.
//
// Contoh: checkout menulis pesanan, reservasi stok, dan event OrderPlaced dalam
// satu transaksi. Kalau reservasi gagal karena stok habis, event-nya ikut
// batal — tidak ada notifikasi "pesanan berhasil" untuk pesanan yang tidak jadi.
