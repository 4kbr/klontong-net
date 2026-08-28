package domain

// Buku besar stok. Append-only.
//
// TODO:
//   type MovementKind string
//   const ( KindPurchase, KindSale, KindReturn, KindAdjustment,
//           KindTransferIn, KindTransferOut, KindDamage, KindExpiry )
//   type Movement struct { ID; OutletID; VariantID; Kind;
//                          Quantity decimal        // positif menambah, negatif mengurangi
//                          BalanceAfter decimal
//                          ReferenceType string; ReferenceID uuid.UUID
//                          Note string; ActorID uuid.UUID; CreatedAt }
//
// SETIAP perubahan stok menulis satu baris di sini, tanpa kecuali. Ini yang
// menjawab "kok stok saya berkurang 12" — pertanyaan yang pasti datang.
//
// QuantityOnHand di tabel stok adalah RINGKASAN yang bisa dihitung ulang dari
// movements. Kalau keduanya tidak cocok, movements yang benar. Sediakan perintah
// rekonsiliasi sejak awal, bukan setelah ada laporan.
//
// KindExpiry penting untuk dagang klontong: barang kedaluwarsa harus keluar dari
// stok dengan sebab yang jelas, bukan disamarkan sebagai adjustment.
