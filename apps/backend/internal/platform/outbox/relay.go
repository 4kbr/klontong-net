package outbox

// Worker pemindah event dari tabel outbox ke bus. Dijalankan dari cmd/worker.
//
// TODO — Run(ctx), loop tiap tick:
//   BEGIN
//     SELECT ... WHERE published_at IS NULL ORDER BY created_at LIMIT $batch
//     FOR UPDATE SKIP LOCKED
//     publish tiap event
//     sukses -> published_at = now(); gagal -> attempts+1, last_error
//   COMMIT
//
// SKIP LOCKED membuatnya aman dijalankan banyak instance.
// At-least-once, jadi consumer harus idempoten.
//
// Sediakan metrik jumlah baris tertunda. Angka yang menanjak adalah sinyal
// paling awal ada consumer macet, dan sering satu-satunya.
