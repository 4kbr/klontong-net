package app

// TODO — RequestPayout(ctx, sellerID uuid.UUID, amount money.Amount):
//   1. periksa rekening pencairan sudah diatur dan tervalidasi
//   2. periksa saldo available mencukupi — hitung dari buku besar, di dalam
//      transaksi
//   3. periksa tidak ada pencairan yang sedang berjalan (ExistsInProgress)
//   4. periksa di atas PayoutMinimum
//   5. buat Payout status requested dengan IdempotencyKey
//   6. RecordPayout ke buku besar
//
// TODO — ProcessPayouts(ctx) (int, error): worker.
//   Untuk tiap payout requested: panggil Disburser.Transfer dengan
//   IdempotencyKey milik payout tersebut.
//   Sukses -> MarkPaid. Gagal -> MarkFailed dan BALIKKAN jurnalnya.
//
//   Worker ini WAJIB memakai kunci terdistribusi. Dua instance yang memproses
//   payout yang sama berarti transfer ganda, dan idempotency key di sisi
//   penyedia adalah pertahanan terakhir, bukan yang pertama.
//
// TODO — HandleDisburseCallback / ReconcilePayouts:
//   Sama seperti pembayaran, status transfer bisa datang terlambat atau tidak
//   datang sama sekali. Sediakan rekonsiliasi lewat Inquiry.
