package app

// Jalur paling rawan di sistem. Salah di sini berarti pesanan palsu ditandai
// lunas.
//
// TODO — HandleWebhook(ctx, rawBody []byte, headers map[string]string) error:
//   1. gateway.VerifyWebhook -> signature tidak sah, TOLAK dan catat.
//      Jangan pernah memproses payload yang belum diverifikasi.
//   2. WebhookRepository.Record -> duplikat, balas 200 dan berhenti.
//      Gateway mengirim ulang sampai mendapat 200; membalas error untuk
//      duplikat membuat mereka mengirim selamanya.
//   3. cari Payment lewat ProviderReference
//   4. COCOKKAN JUMLAHNYA. Berbeda -> ErrAmountMismatch, tandai untuk ditinjau
//      manusia, JANGAN tandai lunas.
//   5. WithinTx: perbarui status pembayaran + outbox EventPaymentSettled
//   6. balas 200 CEPAT. Pemrosesan lanjutan lewat event, bukan di sini —
//      gateway punya batas waktu dan akan mengirim ulang kalau kita lambat.
//
// Urutan status bisa TERBALIK: `settled` kadang datang sebelum `pending`.
// Tangani dengan memeriksa transisi, bukan dengan mengasumsikan urutan.
