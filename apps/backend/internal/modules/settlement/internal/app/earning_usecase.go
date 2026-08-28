package app

// TODO — MatureEarnings(ctx) (int, error): worker harian.
//   ListMatured dengan FOR UPDATE SKIP LOCKED, lalu untuk masing-masing:
//   WithinTx { earning.Mature + RecordEarningMatured }
//
// Idempoten. Earning yang sudah matured dilewati, bukan dicatat ulang.
//
// TODO: ListSellerEarnings, GetSellerBalance.
